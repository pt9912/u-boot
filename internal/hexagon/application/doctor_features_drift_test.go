package application_test

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
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

// TestFormatDriftMessage_S5 pins the message-structure contract
// of `formatDriftMessage` (Review-Followup S5). The broader drift
// tests assert message content with `strings.Contains` for
// resilience; this table-driven unit-test pins the exact format
// so a refactor of the message-builder cannot silently change the
// broader test surface without one local failure here.
func TestFormatDriftMessage_S5(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name        string
		case1       []string
		case2a      []string
		case2b      []string
		jsonPresent bool
		want        string
	}{
		{
			name:        "case1 only, json present",
			case1:       []string{"k1:1"},
			jsonPresent: true,
			want:        "enabled feature(s) missing in devcontainer.json: k1:1. Run `u-boot generate devcontainer` to resync.",
		},
		{
			name:        "case1 only, json absent",
			case1:       []string{"k1:1", "k2:2"},
			jsonPresent: false,
			want:        "enabled feature(s) missing in devcontainer.json (file absent): k1:1, k2:2. Run `u-boot generate devcontainer` to resync.",
		},
		{
			name:        "case2a only",
			case2a:      []string{"old:1"},
			jsonPresent: true,
			want:        "disabled/unset feature(s) still present in devcontainer.json: old:1. Run `u-boot generate devcontainer` to resync.",
		},
		{
			name:        "case2b only (triggers hand-edit hint)",
			case2b:      []string{"foreign:7"},
			jsonPresent: true,
			want:        "devcontainer.json key(s) without a u-boot.yaml pendant (hand-edit or stale render): foreign:7. Run `u-boot generate devcontainer` to resync; reconcile hand-edited keys in u-boot.yaml or remove them from devcontainer.json.",
		},
		{
			name:        "case1 + case2a + case2b (all three)",
			case1:       []string{"new:1"},
			case2a:      []string{"old:1"},
			case2b:      []string{"foreign:7"},
			jsonPresent: true,
			want:        "enabled feature(s) missing in devcontainer.json: new:1; disabled/unset feature(s) still present in devcontainer.json: old:1; devcontainer.json key(s) without a u-boot.yaml pendant (hand-edit or stale render): foreign:7. Run `u-boot generate devcontainer` to resync; reconcile hand-edited keys in u-boot.yaml or remove them from devcontainer.json.",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := application.FormatDriftMessageForTest(tc.case1, tc.case2a, tc.case2b, tc.jsonPresent)
			if got != tc.want {
				t.Errorf("formatDriftMessage mismatch\n  got:  %q\n  want: %q", got, tc.want)
			}
		})
	}
}

// TestDoctor_FeaturesDrift_DisabledPreservesCase2a is the
// Review-Followup S1 anti-regression pin: even if `projectFeatureEntry`
// (renderer-side) gains an `enabled`-filter in the future, the
// drift detector must continue to project disabled entries via
// the shared `projectFeatureSourceVersion`-core so that disabled
// + still-in-JSON keeps surfacing as Case 2a (NOT Case 2b).
// Belongs together with the Refactor in `c2ff32f`+S1.
func TestDoctor_FeaturesDrift_DisabledPreservesCase2a(t *testing.T) {
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
		t.Fatalf("Severity = %v, want Warn", d.Severity)
	}
	// Must be Case 2a (disabled/unset), NOT Case 2b (hand-edit).
	if !strings.Contains(d.Message, "disabled/unset") {
		t.Errorf("Message lacks Case-2a indicator (disabled/unset): %q", d.Message)
	}
	if strings.Contains(d.Message, "hand-edit") {
		t.Errorf("Message has Case-2b indicator (hand-edit) — Review-Followup S1 anti-regression: disabled entry must stay projected: %q", d.Message)
	}
}

// TestDoctor_FeaturesDrift_S2_OrphanWithMatchingJSONKey pins the
// Review-Followup S2 case: an orphan-activation (enabled=true, no
// source override, name not in catalogue) plus a JSON-key that
// matches no projected entry results in two parallel diagnostics
// — Parent-T5 allowlist surfaces "orphan", Drift surfaces Case 2b
// "hand-edit". Both classifications are correct and intentional;
// they describe different concerns.
func TestDoctor_FeaturesDrift_S2_OrphanWithMatchingJSONKey(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	driftSeed(t, fs, `schemaVersion: 1
project:
  name: demo
devcontainer:
  features:
    not-a-real-feature:
      enabled: true
`, `{
  "name": "demo",
  "features": {
    "not-a-real-feature:1": {}
  }
}`)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	// Drift: Case 2b (orphan name → not in knownProjectableKeys
	// → JSON key is "hand-edit" from drift detector's view).
	drift := findDiagnostic(t, resp.Report.Items, "devcontainer.features.drift")
	if drift.Severity != domain.SeverityWarn {
		t.Errorf("drift Severity = %v, want Warn (Case 2b)", drift.Severity)
	}
	if !strings.Contains(drift.Message, "hand-edit") {
		t.Errorf("drift Message lacks Case-2b indicator: %q", drift.Message)
	}
	// Allowlist (Parent-T5 Teil A): orphan-activation warn.
	allow := findDiagnostic(t, resp.Report.Items, "devcontainer.features.allowlist")
	if allow.Severity != domain.SeverityWarn {
		t.Errorf("allowlist Severity = %v, want Warn (orphan-activation)", allow.Severity)
	}
	if !strings.Contains(allow.Message, "not-a-real-feature") {
		t.Errorf("allowlist Message does not name the orphan: %q", allow.Message)
	}
}

// TestDoctor_FeaturesDrift_S3_AllowlistViolationKeepsDriftOK pins
// the Review-Followup S3 cross-check: when a feature has
// `enabled=true` + `source=<not-in-allow>` AND the JSON contains
// the matching render-key, the drift check stays OK (the key IS
// in expectedKeys, no Case 1) while Parent-T5-Teil-A surfaces an
// Error. Pins the orthogonal classification between the two
// devcontainer.features.* checks.
func TestDoctor_FeaturesDrift_S3_AllowlistViolationKeepsDriftOK(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	driftSeed(t, fs, `schemaVersion: 1
project:
  name: demo
devcontainer:
  features:
    rogue:
      enabled: true
      source: https://uninvited.test/feature
`, `{
  "name": "demo",
  "features": {
    "https://uninvited.test/feature:1": {}
  }
}`)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	// Drift: OK — render-key is in expectedKeys + JSON.
	drift := findDiagnostic(t, resp.Report.Items, "devcontainer.features.drift")
	if drift.Severity != domain.SeverityOK {
		t.Errorf("drift Severity = %v, want OK (key present in both expectedKeys and JSON); message = %q",
			drift.Severity, drift.Message)
	}
	// Allowlist: Error — source URL is not in allowlist.
	allow := findDiagnostic(t, resp.Report.Items, "devcontainer.features.allowlist")
	if allow.Severity != domain.SeverityError {
		t.Errorf("allowlist Severity = %v, want Error (LH-FA-DEV-003 violation); message = %q",
			allow.Severity, allow.Message)
	}
}

// TestDoctor_FeaturesDrift_S6_ExplicitEmptyMapWithJSONKeys pins
// the Review-Followup S6 case (slice-plan §AK "kein Skip; Case 2b
// kann feuern"): cfg has `features: {}` (explicit empty map, not
// nil) and JSON carries feature keys — the drift detector must
// fire Case 2b, not skip. Distinguishes "user removed all
// features but never regenerated" from "user never opted into
// features at all".
func TestDoctor_FeaturesDrift_S6_ExplicitEmptyMapWithJSONKeys(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	driftSeed(t, fs, `schemaVersion: 1
project:
  name: demo
devcontainer:
  features: {}
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
		t.Errorf("Severity = %v, want Warn (Case 2b on explicit-empty-cfg + JSON-keys)", d.Severity)
	}
	if !strings.Contains(d.Message, "ghcr.io/devcontainers/features/node:1") {
		t.Errorf("Message lacks the stray JSON key: %q", d.Message)
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
