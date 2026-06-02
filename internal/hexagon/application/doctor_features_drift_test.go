package application_test

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// seedDoctorDevcontainerJSON writes a `.devcontainer/devcontainer.json`
// fixture into the doctor test BaseDir. JSONC is permitted; the
// Drift-Checker runs `stripJSONC` before parsing.
func seedDoctorDevcontainerJSON(t *testing.T, fs *fakeFS, body string) {
	t.Helper()
	if err := fs.WriteFile(filepath.Join(doctorBaseDir, ".devcontainer", "devcontainer.json"),
		[]byte(body), 0o644); err != nil {
		t.Fatalf("seed devcontainer.json: %v", err)
	}
}

// driftSeed wraps the two common fixture writes (u-boot.yaml +
// devcontainer.json) so the drift-tests stay concise.
func driftSeed(t *testing.T, fs *fakeFS, yaml, json string) {
	t.Helper()
	seedDoctorUbootYAMLFeatures(t, fs, yaml)
	if json != "" {
		seedDoctorDevcontainerJSON(t, fs, json)
	}
}

// TestDoctor_FeaturesDrift_OKWhenNothingConfigured pins the
// nichts-konfiguriert-Skip: u-boot.yaml without features + no
// devcontainer.json → OK with skip message.
func TestDoctor_FeaturesDrift_OKWhenNothingConfigured(t *testing.T) {
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
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.features.drift")
	if d.Severity != domain.SeverityOK {
		t.Errorf("Severity = %v, want OK", d.Severity)
	}
	if !strings.Contains(d.Message, "no devcontainer features configured anywhere") {
		t.Errorf("Message = %q, want skip-explanation", d.Message)
	}
}

// TestDoctor_FeaturesDrift_OKWhenInSync pins the no-drift happy
// path: cfg activates `node`, JSON contains exactly that key.
func TestDoctor_FeaturesDrift_OKWhenInSync(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	driftSeed(t, fs, `schemaVersion: 1
project:
  name: demo
devcontainer:
  enabled: true
  features:
    node:
      enabled: true
`, `{
  "name": "demo",
  "features": {
    "ghcr.io/devcontainers/features/node:1": {}
  }
}`)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.features.drift")
	if d.Severity != domain.SeverityOK {
		t.Errorf("Severity = %v (in-sync); message = %q", d.Severity, d.Message)
	}
}

// TestDoctor_FeaturesDrift_Case1_FeatureMissingInJSON pins the
// classic Case-1 drift: feature activated, JSON has no such key.
func TestDoctor_FeaturesDrift_Case1_FeatureMissingInJSON(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	driftSeed(t, fs, `schemaVersion: 1
project:
  name: demo
devcontainer:
  features:
    node:
      enabled: true
`, `{
  "name": "demo",
  "features": {}
}`)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.features.drift")
	if d.Severity != domain.SeverityWarn {
		t.Errorf("Severity = %v, want Warn (Case 1)", d.Severity)
	}
	if !strings.Contains(d.Message, "ghcr.io/devcontainers/features/node:1") {
		t.Errorf("Message lacks the missing render-key: %q", d.Message)
	}
	if !strings.Contains(d.Message, "generate devcontainer") {
		t.Errorf("Message lacks the repair-hint: %q", d.Message)
	}
}

// TestDoctor_FeaturesDrift_Case1_FileMissing pins the file-fehlt-
// Disziplin: cfg activates a feature, devcontainer.json is absent
// → still Case 1 (not OK).
func TestDoctor_FeaturesDrift_Case1_FileMissing(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	seedDoctorUbootYAMLFeatures(t, fs, `schemaVersion: 1
project:
  name: demo
devcontainer:
  features:
    java:
      enabled: true
      version: "21"
`)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.features.drift")
	if d.Severity != domain.SeverityWarn {
		t.Errorf("Severity = %v, want Warn (file absent)", d.Severity)
	}
	if !strings.Contains(d.Message, "file absent") {
		t.Errorf("Message lacks file-absent indicator: %q", d.Message)
	}
	if !strings.Contains(d.Message, "ghcr.io/devcontainers/features/java:21") {
		t.Errorf("Message lacks expected key: %q", d.Message)
	}
}

// TestDoctor_FeaturesDrift_Case2a_DisabledStillInJSON pins
// Case 2a: the user toggled enabled: false but never regenerated.
// The render-key still matches a known cfg-entry, so the message
// names the disabled-drift hint, not the hand-edit hint.
func TestDoctor_FeaturesDrift_Case2a_DisabledStillInJSON(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	driftSeed(t, fs, `schemaVersion: 1
project:
  name: demo
devcontainer:
  features:
    node:
      enabled: false
`, `{
  "name": "demo",
  "features": {
    "ghcr.io/devcontainers/features/node:1": {}
  }
}`)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.features.drift")
	if d.Severity != domain.SeverityWarn {
		t.Errorf("Severity = %v, want Warn (Case 2a)", d.Severity)
	}
	if !strings.Contains(d.Message, "disabled/unset") {
		t.Errorf("Message lacks Case-2a indicator: %q", d.Message)
	}
	if strings.Contains(d.Message, "hand-edit") {
		t.Errorf("Message has Case-2b indicator (hand-edit) where only Case 2a should fire: %q", d.Message)
	}
}

// TestDoctor_FeaturesDrift_Case2b_HandEditUnknownKey pins
// Case 2b: JSON has a feature key that has no cfg pendant at
// all (hand-edit or stale render).
func TestDoctor_FeaturesDrift_Case2b_HandEditUnknownKey(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	driftSeed(t, fs, `schemaVersion: 1
project:
  name: demo
`, `{
  "name": "demo",
  "features": {
    "ghcr.io/devcontainers/features/foreign:7": {}
  }
}`)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.features.drift")
	if d.Severity != domain.SeverityWarn {
		t.Errorf("Severity = %v, want Warn (Case 2b)", d.Severity)
	}
	if !strings.Contains(d.Message, "hand-edit") {
		t.Errorf("Message lacks Case-2b indicator (hand-edit): %q", d.Message)
	}
	if !strings.Contains(d.Message, "ghcr.io/devcontainers/features/foreign:7") {
		t.Errorf("Message lacks the foreign key: %q", d.Message)
	}
}

// TestDoctor_FeaturesDrift_SkipOnParseError pins that a malformed
// devcontainer.json skips with a hint pointing at
// `devcontainer.json.valid`. The primary validity-severity is the
// other check's job.
func TestDoctor_FeaturesDrift_SkipOnParseError(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	driftSeed(t, fs, `schemaVersion: 1
project:
  name: demo
devcontainer:
  features:
    node:
      enabled: true
`, `{not-json`)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.features.drift")
	if d.Severity != domain.SeverityOK {
		t.Errorf("Severity = %v, want OK (parse-error skip)", d.Severity)
	}
	if !strings.Contains(d.Message, "unparseable devcontainer.json") {
		t.Errorf("Message lacks parse-error explanation: %q", d.Message)
	}
}

// TestDoctor_FeaturesDrift_NilVsEmptyFeaturesMap pins the
// nil-vs-explicitly-empty distinction. Both fixtures result in a
// cfg without enabled features; the Skip-Disziplin keeps both
// at OK because the JSON has no features-section either.
// The distinction matters because the Skip-Bedingung must not
// confuse `nil` (no devcontainer key at all) with `features: {}`
// (explicit empty map after user removed everything).
func TestDoctor_FeaturesDrift_NilVsEmptyFeaturesMap(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		yaml string
	}{
		{
			name: "nil features (no devcontainer key)",
			yaml: `schemaVersion: 1
project:
  name: demo
`,
		},
		{
			name: "nil features (devcontainer enabled, no features key)",
			yaml: `schemaVersion: 1
project:
  name: demo
devcontainer:
  enabled: true
`,
		},
		{
			name: "explicit empty map",
			yaml: `schemaVersion: 1
project:
  name: demo
devcontainer:
  features: {}
`,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			svc, fs, _, _, _ := newDoctorService(t)
			seedDoctorUbootYAMLFeatures(t, fs, tc.yaml)
			resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
			if err != nil {
				t.Fatalf("Check: %v", err)
			}
			d := findDiagnostic(t, resp.Report.Items, "devcontainer.features.drift")
			if d.Severity != domain.SeverityOK {
				t.Errorf("Severity = %v, want OK", d.Severity)
			}
		})
	}
}
