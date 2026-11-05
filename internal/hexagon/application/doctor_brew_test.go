package application_test

import (
	"context"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// LH-FA-DEV-013: the doctor check `devcontainer.features.allowlist` reports a
// formula of a tap that is not in `devcontainer.brewTaps.allow` as an error,
// with the text of `generate devcontainer` and without a check of its own.

const doctorBrewHead = "schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n"

func doctorBrewFinding(t *testing.T, yaml string) domain.Diagnostic {
	t.Helper()
	svc, fs, _, _, _ := newDoctorService(t)
	seedDoctorUbootYAMLFeatures(t, fs, yaml)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	return findDiagnostic(t, resp.Report.Items, "devcontainer.features.allowlist")
}

// D1: an active key with an unapproved tap is an error that names the tap, the
// formula and the command; an approved tap, a core formula, a switched-off key
// and a key without `enabled` give no finding about taps; the entry in
// `featureSources.allow` approves no tap.
func TestDoctor_FeaturesAllowlist_BrewTapApproval(t *testing.T) {
	t.Parallel()
	const entry = "  features:\n    brew:\n      enabled: true\n      packages:\n        - ripgrep\n        - pt9912/u-boot/u-boot\n"
	d := doctorBrewFinding(t, doctorBrewHead+entry)
	if d.Severity != domain.SeverityError {
		t.Fatalf("Severity = %v, want Error; message = %q", d.Severity, d.Message)
	}
	for _, want := range []string{
		`the tap "pt9912/u-boot"`, `"pt9912/u-boot/u-boot"`, "`u-boot config set devcontainer.brewTaps.allow pt9912/u-boot`", "--yes does not substitute",
	} {
		if !strings.Contains(d.Message, want) {
			t.Errorf("Message lacks %q: %s", want, d.Message)
		}
	}
	for name, yaml := range map[string]string{
		"approved":             doctorBrewHead + "  brewTaps:\n    allow:\n      - pt9912/u-boot\n" + entry,
		"core formula only":    doctorBrewHead + "  features:\n    brew:\n      enabled: true\n      packages:\n        - ripgrep\n",
		"switched off":         doctorBrewHead + strings.Replace(entry, "enabled: true", "enabled: false", 1),
		"homebrew alone":       doctorBrewHead + "  features:\n    homebrew:\n      enabled: true\n",
		"empty list":           doctorBrewHead + "  features:\n    brew:\n      enabled: true\n",
		"no devcontainer tree": "schemaVersion: 1\nproject:\n  name: demo\n",
	} {
		if d := doctorBrewFinding(t, yaml); d.Severity == domain.SeverityError {
			t.Errorf("%s: Severity = Error, want no finding about taps: %s", name, d.Message)
		}
	}
	d = doctorBrewFinding(t, doctorBrewHead+"  featureSources:\n    allow:\n      - pt9912/u-boot\n"+entry)
	if d.Severity != domain.SeverityError {
		t.Errorf("the entry of featureSources.allow must approve no tap: Severity = %v, %s", d.Severity, d.Message)
	}
}

// D2: the finding stands next to the findings of the feature sources in one
// message, and the conflict with `homebrew` is no finding of the check.
func TestDoctor_FeaturesAllowlist_BrewFindingNextToTheOthers(t *testing.T) {
	t.Parallel()
	d := doctorBrewFinding(t, doctorBrewHead+
		"  features:\n    rogue:\n      enabled: true\n      source: https://uninvited.test/feature\n"+
		"    brew:\n      enabled: true\n      packages:\n        - a/b/c\n")
	if d.Severity != domain.SeverityError || !strings.Contains(d.Message, "rogue") || !strings.Contains(d.Message, `the tap "a/b"`) {
		t.Errorf("want both findings in one error message: %v %s", d.Severity, d.Message)
	}
	d = doctorBrewFinding(t, doctorBrewHead+"  features:\n    brew:\n      enabled: true\n    homebrew:\n      enabled: true\n")
	if d.Severity == domain.SeverityError {
		t.Errorf("the conflict is a rule of generate and of the write path, not a doctor finding: %s", d.Message)
	}
}

// D4: the key renders nothing under `features`, so the drift check stays
// silent for it with and without a `features` map in devcontainer.json, and
// the allowlist check finds nothing about it.
func TestDoctor_BrewKey_NoDriftFinding(t *testing.T) {
	t.Parallel()
	const yaml = doctorBrewHead + "  features:\n    brew:\n      enabled: true\n      packages:\n        - ripgrep\n"
	for name, jsonBody := range map[string]string{
		"features empty": `{"features": {}}`,
		"no features":    `{"name": "demo"}`,
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs, _, _, _ := newDoctorService(t)
			driftSeed(t, fs, yaml, jsonBody)
			resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
			if err != nil {
				t.Fatalf("Check: %v", err)
			}
			for _, id := range []string{"devcontainer.features.drift", "devcontainer.features.allowlist"} {
				if d := findDiagnostic(t, resp.Report.Items, id); d.Severity != domain.SeverityOK {
					t.Errorf("%s: Severity = %v, want OK; message = %q", id, d.Severity, d.Message)
				}
			}
		})
	}
}

// D3: a package list outside the form is an error of `uboot.yaml.valid` (the
// read-path validation), without a check of its own.
func TestDoctor_UbootYaml_ErrorOnBadBrewEntries(t *testing.T) {
	t.Parallel()
	for name, body := range map[string]string{
		"bad formula":  "  features:\n    brew:\n      packages:\n        - rip grep\n",
		"list at node": "  features:\n    node:\n      packages:\n        - ripgrep\n",
		"bad tap":      "  brewTaps:\n    allow:\n      - pt9912\n",
	} {
		svc, fs, _, _, _ := newDoctorService(t)
		seedDoctorUbootYAMLFeatures(t, fs, doctorBrewHead+body)
		resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
		if err != nil {
			t.Fatalf("Check: %v", err)
		}
		d := findDiagnostic(t, resp.Report.Items, "uboot.yaml.valid")
		if d.Severity != domain.SeverityError || !strings.Contains(d.Message, "devcontainer schema invalid") {
			t.Errorf("%s: Severity = %v, message = %q, want the schema error", name, d.Severity, d.Message)
		}
	}
}
