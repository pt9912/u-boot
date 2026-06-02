package application_test

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// seedDoctorUbootYAMLFeatures writes a u-boot.yaml fixture into the
// doctor test BaseDir. Callers pass the YAML body verbatim — the
// helper centralises the file path so the eight T5 tests below stay
// concise.
func seedDoctorUbootYAMLFeatures(t *testing.T, fs *fakeFS, body string) {
	t.Helper()
	if err := fs.WriteFile(filepath.Join(doctorBaseDir, "u-boot.yaml"),
		[]byte(body), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
}

// TestDoctor_UbootYaml_ErrorOnBadFeatureSource is the Audit-
// Followup A1 wiring pin: a hand-edited u-boot.yaml with an
// invalid `featureSources.allow` entry (ftp:// — unsupported
// scheme) surfaces as Error severity on the
// `uboot.yaml.valid` check (LH-FA-DEV-003 / Spec §1353 → Exit-10
// when the user surface goes through the CLI; in-Doctor we just
// classify it as Error so the user sees the consolidated
// u-boot.yaml-validity report).
func TestDoctor_UbootYaml_ErrorOnBadFeatureSource(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	seedDoctorUbootYAMLFeatures(t, fs, `schemaVersion: 1
project:
  name: demo
devcontainer:
  featureSources:
    allow:
      - ftp://bad.test/feature
`)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "uboot.yaml.valid")
	if d.Severity != domain.SeverityError {
		t.Errorf("Severity = %v, want Error (LH-FA-DEV-003 schema)", d.Severity)
	}
	if !strings.Contains(d.Message, "devcontainer schema invalid") {
		t.Errorf("Message lacks devcontainer-schema indicator: %q", d.Message)
	}
	if !strings.Contains(d.Hint, "LH-FA-DEV-003") {
		t.Errorf("Hint lacks LH-FA-DEV-003 reference: %q", d.Hint)
	}
}

// TestDoctor_FeaturesAllowlist_OKWhenNoUbootYaml pins the
// skip-on-missing-config branch: without u-boot.yaml the check
// keeps quiet (primary file-presence diagnostics live in
// `uboot.yaml.valid`).
func TestDoctor_FeaturesAllowlist_OKWhenNoUbootYaml(t *testing.T) {
	t.Parallel()
	svc, _, _, _, _ := newDoctorService(t)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.features.allowlist")
	if d.Severity != domain.SeverityOK {
		t.Errorf("Severity = %v, want OK (u-boot.yaml absent → skip)", d.Severity)
	}
}

// TestDoctor_FeaturesAllowlist_OKWhenNoFeatures pins
// `spec/lastenheft.md:2394` (LH-AK-005 erwartetes Ergebnis):
// `u-boot doctor` enthält keinen `error` zu `devcontainer`-
// Konfiguration oder Feature-Quellen. A fresh `init --devcontainer`
// produces a u-boot.yaml with `devcontainer.enabled: true` and no
// `features:` map → the allowlist check must classify as OK with
// a skip-explanation. The full init→doctor end-to-end pin lives in
// the T6 acceptance test (TestLHFADEV003_*); this check-level
// test is the unit-shaped equivalent.
func TestDoctor_FeaturesAllowlist_OKWhenNoFeatures(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	seedDoctorUbootYAMLFeatures(t, fs, `schemaVersion: 1
project:
  name: demo
devcontainer:
  enabled: true
`)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.features.allowlist")
	if d.Severity != domain.SeverityOK {
		t.Errorf("Severity = %v, want OK (no features configured)", d.Severity)
	}
	if !strings.Contains(d.Message, "No devcontainer features") {
		t.Errorf("Message = %q, want skip-explanation", d.Message)
	}
}

// TestDoctor_FeaturesAllowlist_OKWhenCatalogued pins the happy
// path: a catalogued feature with enabled=true, no source override,
// passes without warning.
func TestDoctor_FeaturesAllowlist_OKWhenCatalogued(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	seedDoctorUbootYAMLFeatures(t, fs, `schemaVersion: 1
project:
  name: demo
devcontainer:
  enabled: true
  features:
    node:
      enabled: true
    java:
      enabled: true
      version: "21"
`)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.features.allowlist")
	if d.Severity != domain.SeverityOK {
		t.Errorf("Severity = %v, want OK (catalogued); message = %q", d.Severity, d.Message)
	}
}

// TestDoctor_FeaturesAllowlist_OKWhenSourceInAllow pins the happy
// path for external features: source override is in the allowlist
// → OK.
func TestDoctor_FeaturesAllowlist_OKWhenSourceInAllow(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	seedDoctorUbootYAMLFeatures(t, fs, `schemaVersion: 1
project:
  name: demo
devcontainer:
  enabled: true
  featureSources:
    allow:
      - https://ghcr.io/orgX/features/custom-rust
  features:
    custom-rust:
      enabled: true
      source: https://ghcr.io/orgX/features/custom-rust
`)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.features.allowlist")
	if d.Severity != domain.SeverityOK {
		t.Errorf("Severity = %v, want OK (source in allowlist); message = %q", d.Severity, d.Message)
	}
}

// TestDoctor_FeaturesAllowlist_ErrorWhenSourceNotInAllow pins the
// LH-FA-DEV-003 Spec §720 violation: source override is set, but
// not in `featureSources.allow` → Error with repair hint that
// names both the URL and the LH-NFA-SEC-004 `--yes`-not-sufficient
// clause.
func TestDoctor_FeaturesAllowlist_ErrorWhenSourceNotInAllow(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	seedDoctorUbootYAMLFeatures(t, fs, `schemaVersion: 1
project:
  name: demo
devcontainer:
  enabled: true
  features:
    rogue:
      enabled: true
      source: https://uninvited.test/feature
`)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.features.allowlist")
	if d.Severity != domain.SeverityError {
		t.Errorf("Severity = %v, want Error; message = %q", d.Severity, d.Message)
	}
	for _, want := range []string{
		"rogue",
		"https://uninvited.test/feature",
		"LH-FA-DEV-003",
		"LH-NFA-SEC-004",
	} {
		if !strings.Contains(d.Message, want) {
			t.Errorf("Message missing %q\n  full: %s", want, d.Message)
		}
	}
}

// TestDoctor_FeaturesAllowlist_WarnOnOrphanActivation pins the
// orphan-activation branch (feature name not in catalogue + no
// source override). Renderer skips silently (T3 contract); doctor
// surfaces a warn so users notice the typo / missing override.
func TestDoctor_FeaturesAllowlist_WarnOnOrphanActivation(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	seedDoctorUbootYAMLFeatures(t, fs, `schemaVersion: 1
project:
  name: demo
devcontainer:
  enabled: true
  features:
    nde:
      enabled: true
`)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.features.allowlist")
	if d.Severity != domain.SeverityWarn {
		t.Errorf("Severity = %v, want Warn (orphan); message = %q", d.Severity, d.Message)
	}
	if !strings.Contains(d.Message, "nde") {
		t.Errorf("Message does not name the orphan feature `nde`: %q", d.Message)
	}
	if !strings.Contains(d.Message, "renderer skips") {
		t.Errorf("Message does not explain renderer skip: %q", d.Message)
	}
}

// TestDoctor_FeaturesAllowlist_WarnOnEnabledKeyMissing pins the
// LH-FA-ADD-005 §893 enabled-key-missing convention extended to
// devcontainer features: explicit `true`/`false` is required.
func TestDoctor_FeaturesAllowlist_WarnOnEnabledKeyMissing(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	seedDoctorUbootYAMLFeatures(t, fs, `schemaVersion: 1
project:
  name: demo
devcontainer:
  enabled: true
  features:
    node: {}
`)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.features.allowlist")
	if d.Severity != domain.SeverityWarn {
		t.Errorf("Severity = %v, want Warn (enabled missing); message = %q", d.Severity, d.Message)
	}
	if !strings.Contains(d.Message, "node") || !strings.Contains(d.Message, "enabled:") {
		t.Errorf("Message lacks the missing-enabled-key hint: %q", d.Message)
	}
}

// TestDoctor_FeaturesAllowlist_ErrorBeatsWarn pins the worst-
// severity-wins contract: when both an allowlist violation AND an
// orphan exist in the same project, the call surfaces Error (not
// Warn). The user fixes the LH-FA-DEV-003 violation first; the
// orphan stays visible on the next run.
func TestDoctor_FeaturesAllowlist_ErrorBeatsWarn(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	seedDoctorUbootYAMLFeatures(t, fs, `schemaVersion: 1
project:
  name: demo
devcontainer:
  enabled: true
  features:
    nde:
      enabled: true
    rogue:
      enabled: true
      source: https://uninvited.test/feature
`)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.features.allowlist")
	if d.Severity != domain.SeverityError {
		t.Errorf("Severity = %v, want Error (combination); message = %q", d.Severity, d.Message)
	}
}
