package application_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// parseDevcontainerFeatures extracts the `features` map from a
// rendered devcontainer.json (after stripping JSONC comments).
// Returns nil when the key is absent — matching the "darf fehlen"
// contract analogous to forwardPorts.
func parseDevcontainerFeatures(t *testing.T, body []byte) map[string]map[string]any {
	t.Helper()
	stripped := application.StripJSONCForTest(body)
	if !json.Valid(stripped) {
		t.Fatalf("rendered devcontainer.json is not valid JSON after stripJSONC:\n%s",
			stripped)
	}
	var shape struct {
		Features map[string]map[string]any `json:"features"`
	}
	if err := json.Unmarshal(stripped, &shape); err != nil {
		t.Fatalf("unmarshal devcontainer.json features: %v", err)
	}
	return shape.Features
}

// seedUBootYAMLWithFeatures writes a u-boot.yaml fixture into the
// generator's test base directory. The devcontainer body is appended
// verbatim — callers supply the full `devcontainer:` subtree.
func seedUBootYAMLWithFeatures(t *testing.T, fs *fakeFS, devcontainerBody string) {
	t.Helper()
	body := "schemaVersion: 1\nproject:\n  name: t-uboot-gen\nservices:\n  postgres:\n    enabled: true\n" + devcontainerBody
	if err := fs.WriteFile(filepath.Join(generateTestBaseDir, "u-boot.yaml"),
		[]byte(body), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
}

// TestCollectDevcontainerFeatures_Projection pins the T3 collector
// contract: enabled filter + catalogue lookup + source override +
// version fallback + alphabetical sort by Source. Driven by YAML
// fixtures so the YAML codec path is exercised too.
func TestCollectDevcontainerFeatures_Projection(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		yaml string
		want []application.DevcontainerFeatureDataForTest
	}{
		{
			name: "no devcontainer key",
			yaml: "schemaVersion: 1\nproject:\n  name: demo\n",
			want: nil,
		},
		{
			name: "empty features map",
			yaml: "schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n",
			want: nil,
		},
		{
			name: "enabled=nil is skipped",
			yaml: `schemaVersion: 1
project:
  name: demo
devcontainer:
  features:
    node: {}
`,
			want: nil,
		},
		{
			name: "enabled=false is skipped",
			yaml: `schemaVersion: 1
project:
  name: demo
devcontainer:
  features:
    node:
      enabled: false
`,
			want: nil,
		},
		{
			name: "catalogue lookup, defaultVersion applied",
			yaml: `schemaVersion: 1
project:
  name: demo
devcontainer:
  features:
    node:
      enabled: true
`,
			want: []application.DevcontainerFeatureDataForTest{
				{Source: "ghcr.io/devcontainers/features/node", Version: "1"},
			},
		},
		{
			name: "version override beats defaultVersion",
			yaml: `schemaVersion: 1
project:
  name: demo
devcontainer:
  features:
    java:
      enabled: true
      version: "21"
`,
			want: []application.DevcontainerFeatureDataForTest{
				{Source: "ghcr.io/devcontainers/features/java", Version: "21"},
			},
		},
		{
			name: "source override + version override",
			yaml: `schemaVersion: 1
project:
  name: demo
devcontainer:
  features:
    custom-rust:
      enabled: true
      source: https://ghcr.io/orgX/features/custom-rust
      version: "2"
`,
			want: []application.DevcontainerFeatureDataForTest{
				{Source: "https://ghcr.io/orgX/features/custom-rust", Version: "2"},
			},
		},
		{
			name: "source override without version defaults to 1",
			yaml: `schemaVersion: 1
project:
  name: demo
devcontainer:
  features:
    custom:
      enabled: true
      source: https://example.test/features/custom
`,
			want: []application.DevcontainerFeatureDataForTest{
				{Source: "https://example.test/features/custom", Version: "1"},
			},
		},
		{
			name: "unknown catalogue name without source is skipped",
			yaml: `schemaVersion: 1
project:
  name: demo
devcontainer:
  features:
    not-in-catalogue:
      enabled: true
`,
			// T3 silently skips; T4 will raise Exit-Code-10 via
			// allowlist enforcement.
			want: nil,
		},
		{
			name: "alphabetical sort by Source",
			yaml: `schemaVersion: 1
project:
  name: demo
devcontainer:
  features:
    node:
      enabled: true
    git:
      enabled: true
    java:
      enabled: true
`,
			want: []application.DevcontainerFeatureDataForTest{
				{Source: "ghcr.io/devcontainers/features/git", Version: "1"},
				{Source: "ghcr.io/devcontainers/features/java", Version: "1"},
				{Source: "ghcr.io/devcontainers/features/node", Version: "1"},
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := application.CollectDevcontainerFeaturesForTest(t, []byte(tc.yaml))
			if len(got) == 0 && len(tc.want) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("collectDevcontainerFeatures =\n  %+v\nwant\n  %+v", got, tc.want)
			}
		})
	}
}

// TestGenerateDevcontainer_FeaturesAbsent_NoFeaturesKey pins the
// backwards-compatibility contract: a project without
// devcontainer.features renders devcontainer.json byte-equivalent
// to the pre-T3 shape (no `features` key at all). Critical for
// existing M5-/M7-test fixtures that don't seed feature config.
func TestGenerateDevcontainer_FeaturesAbsent_NoFeaturesKey(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLPostgres(t, fs)
	seedGenerateComposeYAML(t, fs, composeYAMLPostgres)

	if _, err := generateDevcontainer(t, svc); err != nil {
		t.Fatalf("generate devcontainer: %v", err)
	}
	body, err := fs.ReadFile(devcontainerJSONPath())
	if err != nil {
		t.Fatalf("read devcontainer.json: %v", err)
	}
	stripped := application.StripJSONCForTest(body)
	if !json.Valid(stripped) {
		t.Fatalf("rendered devcontainer.json invalid JSONC: %s", stripped)
	}
	if bytes.Contains(stripped, []byte(`"features"`)) {
		t.Errorf("features key emitted when no features configured; body:\n%s", stripped)
	}
}

// TestGenerateDevcontainer_FeaturesRendered_KeysSortedAndShape pins
// the LH-FA-DEV-003 happy path: catalogue features + a source-
// override entry compose a `features` map keyed by `<source>:<version>`,
// alphabetically sorted by Source.
func TestGenerateDevcontainer_FeaturesRendered_KeysSortedAndShape(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateComposeYAML(t, fs, composeYAMLPostgres)
	seedUBootYAMLWithFeatures(t, fs, `devcontainer:
  enabled: true
  featureSources:
    allow:
      - https://ghcr.io/orgX/features/custom-rust
  features:
    node:
      enabled: true
    java:
      enabled: true
      version: "21"
    custom-rust:
      enabled: true
      source: https://ghcr.io/orgX/features/custom-rust
`)

	if _, err := generateDevcontainer(t, svc); err != nil {
		t.Fatalf("generate: %v", err)
	}
	body, err := fs.ReadFile(devcontainerJSONPath())
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	features := parseDevcontainerFeatures(t, body)

	// Expected keys (composed by the renderer):
	wantKeys := []string{
		"ghcr.io/devcontainers/features/java:21",
		"ghcr.io/devcontainers/features/node:1",
		"https://ghcr.io/orgX/features/custom-rust:1",
	}
	gotKeys := make([]string, 0, len(features))
	for k := range features {
		gotKeys = append(gotKeys, k)
	}
	if len(gotKeys) != len(wantKeys) {
		t.Fatalf("feature count = %d (keys = %v), want %d (%v)",
			len(gotKeys), gotKeys, len(wantKeys), wantKeys)
	}
	// Verify each expected key exists; the JSON parser loses order
	// so explicit map-membership is the right shape-pin. The source-
	// ordering pin is covered by the byte-equal idempotency test
	// below: a re-render with shuffled YAML map iteration must
	// produce byte-identical JSON.
	for _, wk := range wantKeys {
		if _, ok := features[wk]; !ok {
			t.Errorf("feature key %q missing from rendered map %v", wk, gotKeys)
		}
	}
}

// TestGenerateDevcontainer_SchemaValidation_RejectsBadAllowEntry is
// the Audit-Followup A1 wiring pin for the Generate-use-case:
// a hand-edited u-boot.yaml with an invalid `featureSources.allow`
// entry (e.g. `ftp://`) must fail the generate with
// ErrGenerateManualConflict (LH-FA-CLI-006 Exit-Code 10) — before
// the fix, the bad URL silently propagated into devcontainer.json.
func TestGenerateDevcontainer_SchemaValidation_RejectsBadAllowEntry(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateComposeYAML(t, fs, composeYAMLPostgres)
	if err := fs.WriteFile(filepath.Join(generateTestBaseDir, "u-boot.yaml"),
		[]byte(`schemaVersion: 1
project:
  name: t-uboot-gen
devcontainer:
  featureSources:
    allow:
      - ftp://bad.test/feature
`), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
	_, err := svc.Generate(context.Background(), driving.GenerateRequest{
		BaseDir:  generateTestBaseDir,
		Artifact: domain.ArtifactDevcontainer,
	})
	if err == nil {
		t.Fatalf("expected error for bad allowlist entry, got nil")
	}
	if !errors.Is(err, driving.ErrGenerateManualConflict) {
		t.Errorf("err = %v, want wrap of ErrGenerateManualConflict (Exit-Code 10)", err)
	}
}

// TestGenerateDevcontainer_SchemaValidation_RejectsBadFeatureSource
// pins the same wiring for the `features.<name>.source`-override
// path.
func TestGenerateDevcontainer_SchemaValidation_RejectsBadFeatureSource(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateComposeYAML(t, fs, composeYAMLPostgres)
	if err := fs.WriteFile(filepath.Join(generateTestBaseDir, "u-boot.yaml"),
		[]byte(`schemaVersion: 1
project:
  name: t-uboot-gen
devcontainer:
  features:
    custom:
      enabled: true
      source: ftp://bad.test/feature
`), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
	_, err := svc.Generate(context.Background(), driving.GenerateRequest{
		BaseDir:  generateTestBaseDir,
		Artifact: domain.ArtifactDevcontainer,
	})
	if err == nil {
		t.Fatalf("expected error for bad features.<name>.source, got nil")
	}
	if !errors.Is(err, driving.ErrGenerateManualConflict) {
		t.Errorf("err = %v, want wrap of ErrGenerateManualConflict", err)
	}
}

// TestGenerateDevcontainer_FeaturesIdempotent_SameSourceTwoVersions
// is the Audit-Followup A2 anti-regression pin: two enabled
// feature entries that share a Source but carry different
// Versions must produce a deterministic alphabetical render across
// 30 repeated runs. Before the fix, `collectDevcontainerFeatures`
// sorted by Source alone — Go's map-iteration randomness flipped
// the relative order of same-Source/different-Version entries and
// turned every other `generate devcontainer` into a spurious
// UpdatedBlock action, breaking LH-FA-DEV-005 idempotency.
//
// The fixture uses a `source:`-override (the only realistic way
// to get two map entries with identical Source: catalogue-only
// activation is one entry per source by construction).
func TestGenerateDevcontainer_FeaturesIdempotent_SameSourceTwoVersions(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateComposeYAML(t, fs, composeYAMLPostgres)
	seedUBootYAMLWithFeatures(t, fs, `devcontainer:
  enabled: true
  featureSources:
    allow:
      - https://example.test/feat
  features:
    primary:
      enabled: true
      source: https://example.test/feat
      version: "1"
    secondary:
      enabled: true
      source: https://example.test/feat
      version: "2"
`)

	// First Generate writes the file.
	if _, err := generateDevcontainer(t, svc); err != nil {
		t.Fatalf("first run: %v", err)
	}
	writesAfterFirst := len(fs.writtenPaths())

	// 30 follow-up generates must all be NoOp regardless of the
	// underlying map-iteration order.
	for i := 0; i < 30; i++ {
		resp, err := generateDevcontainer(t, svc)
		if err != nil {
			t.Fatalf("run %d: %v", i+2, err)
		}
		if resp.Action != driving.GenerateActionNoOp {
			t.Errorf("run %d Action = %v, want NoOp (sort instability flips block content)",
				i+2, resp.Action)
		}
	}
	if delta := len(fs.writtenPaths()) - writesAfterFirst; delta != 0 {
		t.Errorf("30 follow-up generates produced %d extra WriteFile call(s), want 0",
			delta)
	}
}

// TestGenerateDevcontainer_AllowExternalFeatureSources_Append pins
// the LH-FA-DEV-003 / LH-FA-DEV-003 flag-wiring on `generate
// devcontainer`: invoking the use case with a non-empty
// AllowExternalFeatureSources list appends the URLs to
// `devcontainer.featureSources.allow` (with silent-dedupe and
// format validation) before the actual generate runs.
func TestGenerateDevcontainer_AllowExternalFeatureSources_Append(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateComposeYAML(t, fs, composeYAMLPostgres)
	// u-boot.yaml starts with an existing allowlist entry.
	body := `schemaVersion: 1
project:
  name: t-uboot-gen
services:
  postgres:
    enabled: true
devcontainer:
  enabled: true
  featureSources:
    allow:
      - https://a.test/x
`
	if err := fs.WriteFile(filepath.Join(generateTestBaseDir, "u-boot.yaml"),
		[]byte(body), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}

	resp, err := svc.Generate(context.Background(), driving.GenerateRequest{
		BaseDir:                     generateTestBaseDir,
		Artifact:                    domain.ArtifactDevcontainer,
		AllowExternalFeatureSources: []string{"https://b.test/y", "https://a.test/x"},
	})
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if resp.Action == driving.GenerateActionNoOp {
		t.Errorf("Action = NoOp; want a write because the allowlist was extended")
	}

	// Verify u-boot.yaml carries the merged + deduped list.
	updated, err := fs.ReadFile(filepath.Join(generateTestBaseDir, "u-boot.yaml"))
	if err != nil {
		t.Fatalf("read u-boot.yaml: %v", err)
	}
	for _, want := range []string{"https://a.test/x", "https://b.test/y"} {
		if !strings.Contains(string(updated), want) {
			t.Errorf("u-boot.yaml does not contain %q after generate with flag\nbody:\n%s",
				want, updated)
		}
	}
	// `a.test/x` must appear exactly once (silent-dedupe).
	if c := strings.Count(string(updated), "https://a.test/x"); c != 1 {
		t.Errorf("a.test/x count = %d, want 1 (silent-dedupe)", c)
	}
}

// TestGenerateDevcontainer_AllowExternalFeatureSources_NoOp pins
// that a re-run with the same flag values does not rewrite
// u-boot.yaml — equalAllowLists short-circuits the merge step.
func TestGenerateDevcontainer_AllowExternalFeatureSources_NoOp(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateComposeYAML(t, fs, composeYAMLPostgres)
	body := `schemaVersion: 1
project:
  name: t-uboot-gen
services:
  postgres:
    enabled: true
devcontainer:
  enabled: true
  featureSources:
    allow:
      - https://a.test/x
`
	if err := fs.WriteFile(filepath.Join(generateTestBaseDir, "u-boot.yaml"),
		[]byte(body), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}

	if _, err := svc.Generate(context.Background(), driving.GenerateRequest{
		BaseDir:                     generateTestBaseDir,
		Artifact:                    domain.ArtifactDevcontainer,
		AllowExternalFeatureSources: []string{"https://a.test/x"}, // already present
	}); err != nil {
		t.Fatalf("generate: %v", err)
	}
	// u-boot.yaml should not have been rewritten; the original
	// body (with its YAML formatting) should be byte-equal.
	after, err := fs.ReadFile(filepath.Join(generateTestBaseDir, "u-boot.yaml"))
	if err != nil {
		t.Fatalf("read u-boot.yaml: %v", err)
	}
	if string(after) != body {
		t.Errorf("u-boot.yaml rewritten on NoOp allowlist update\nbefore:\n%s\nafter:\n%s",
			body, after)
	}
}

// TestGenerateDevcontainer_AllowExternalFeatureSources_AtomicOnConflict
// pins the slice-v1-devcontainer-features Review-Followup R2 fix:
// when the devcontainer plan-and-execute aborts with
// ErrGenerateManualConflict (e.g. user-managed devcontainer.json
// without an init block), the `--allow-external-feature-sources`
// flag values must NOT be persisted to u-boot.yaml — otherwise the
// flag-fed allowlist extension survives a failed run and the marshal-
// rewrite also strips u-boot.yaml comments unrecoverably. This test
// pins that u-boot.yaml stays byte-identical when the generate
// aborts.
func TestGenerateDevcontainer_AllowExternalFeatureSources_AtomicOnConflict(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	originalYAML := `schemaVersion: 1
project:
  name: t-uboot-gen
# user comment that must survive across failed runs
services:
  postgres:
    enabled: true
devcontainer:
  enabled: true
`
	if err := fs.WriteFile(filepath.Join(generateTestBaseDir, "u-boot.yaml"),
		[]byte(originalYAML), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
	seedGenerateComposeYAML(t, fs, composeYAMLPostgres)
	// Seed devcontainer.json with content but no init block →
	// plan classifies as Manual-Conflict.
	if err := fs.WriteFile(devcontainerJSONPath(),
		[]byte(`{"name": "user-managed"}`), 0o644); err != nil {
		t.Fatalf("seed devcontainer.json: %v", err)
	}

	_, err := svc.Generate(context.Background(), driving.GenerateRequest{
		BaseDir:                     generateTestBaseDir,
		Artifact:                    domain.ArtifactDevcontainer,
		AllowExternalFeatureSources: []string{"https://example.test/never-applied"},
	})
	if !errors.Is(err, driving.ErrGenerateManualConflict) {
		t.Fatalf("err = %v, want wrap of ErrGenerateManualConflict", err)
	}

	// u-boot.yaml must be byte-identical: neither the allowlist
	// entry nor the comment loss should have happened.
	after, readErr := fs.ReadFile(filepath.Join(generateTestBaseDir, "u-boot.yaml"))
	if readErr != nil {
		t.Fatalf("read u-boot.yaml: %v", readErr)
	}
	if string(after) != originalYAML {
		t.Errorf("u-boot.yaml mutated despite generate conflict\nbefore:\n%s\nafter:\n%s",
			originalYAML, after)
	}
}

// TestGenerateDevcontainer_AllowExternalFeatureSources_InvalidURL
// pins that bad flag input fails the generate (before any artefact
// write) and surfaces the LH-FA-DEV-003 invalid-source sentinel.
func TestGenerateDevcontainer_AllowExternalFeatureSources_InvalidURL(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLPostgres(t, fs)
	seedGenerateComposeYAML(t, fs, composeYAMLPostgres)

	_, err := svc.Generate(context.Background(), driving.GenerateRequest{
		BaseDir:                     generateTestBaseDir,
		Artifact:                    domain.ArtifactDevcontainer,
		AllowExternalFeatureSources: []string{"not-a-url"},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, application.ErrInvalidFeatureSource) {
		t.Errorf("err = %v, want wrap of ErrInvalidFeatureSource", err)
	}
}

// TestGenerateDevcontainer_FeaturesIdempotent pins that a re-run
// with features configured is a NoOp (no WriteFile, byte-equal
// re-render). Critical: map iteration shuffles in Go would
// otherwise flip the JSON byte order across calls and trigger
// spurious ReplaceBlock-actions.
func TestGenerateDevcontainer_FeaturesIdempotent(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateComposeYAML(t, fs, composeYAMLPostgres)
	seedUBootYAMLWithFeatures(t, fs, `devcontainer:
  features:
    node:
      enabled: true
    git:
      enabled: true
    java:
      enabled: true
`)

	if _, err := generateDevcontainer(t, svc); err != nil {
		t.Fatalf("first run: %v", err)
	}
	writesAfterFirst := len(fs.writtenPaths())

	resp, err := generateDevcontainer(t, svc)
	if err != nil {
		t.Fatalf("second run: %v", err)
	}
	if resp.Action != driving.GenerateActionNoOp {
		t.Errorf("second run Action = %v, want NoOp", resp.Action)
	}
	if delta := len(fs.writtenPaths()) - writesAfterFirst; delta != 0 {
		t.Errorf("second run produced %d WriteFile call(s), want 0", delta)
	}
}

// Catalogue keys homebrew and sdkman: homebrew renders the community
// feature and is reported as third party (LH-NFA-SEC-004); sdkman renders
// the official java feature; java and sdkman together render it once.
func TestGenerateDevcontainer_FeaturesHomebrewSdkman(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		yaml      string
		wantKeys  []string
		wantWarns int
	}{
		"homebrew": {"devcontainer:\n  enabled: true\n  features:\n    homebrew:\n      enabled: true\n",
			[]string{"ghcr.io/meaningful-ooo/devcontainer-features/homebrew:2"}, 1},
		"sdkman": {"devcontainer:\n  enabled: true\n  features:\n    sdkman:\n      enabled: true\n",
			[]string{"ghcr.io/devcontainers/features/java:1"}, 0},
		"python": {"devcontainer:\n  enabled: true\n  features:\n    python:\n      enabled: true\n",
			[]string{"ghcr.io/devcontainers/features/python:1"}, 0},
		"sshd": {"devcontainer:\n  enabled: true\n  features:\n    sshd:\n      enabled: true\n",
			[]string{"ghcr.io/devcontainers/features/sshd:1"}, 0},
		"java and sdkman": {"devcontainer:\n  enabled: true\n  features:\n    java:\n      enabled: true\n    sdkman:\n      enabled: true\n",
			[]string{"ghcr.io/devcontainers/features/java:1"}, 0},
		"both": {"devcontainer:\n  enabled: true\n  features:\n    homebrew:\n      enabled: true\n    sdkman:\n      enabled: true\n",
			[]string{"ghcr.io/devcontainers/features/java:1", "ghcr.io/meaningful-ooo/devcontainer-features/homebrew:2"}, 1},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, tc.yaml)
			resp, err := generateSandbox(svc, false)
			if err != nil {
				t.Fatalf("generate: %v", err)
			}
			feats, _ := devcontainerJSONMap(t, fs)["features"].(map[string]any)
			var got []string
			for k := range feats {
				got = append(got, k)
			}
			sort.Strings(got)
			if !reflect.DeepEqual(got, tc.wantKeys) {
				t.Errorf("features = %v, want %v", got, tc.wantKeys)
			}
			warns := 0
			for _, w := range resp.Warnings {
				if w.Code == "LH-NFA-SEC-004" {
					warns++
					if !strings.Contains(w.Message, "third-party publisher meaningful-ooo") {
						t.Errorf("unexpected third-party warning: %+v", w)
					}
				}
			}
			if warns != tc.wantWarns {
				t.Errorf("third-party warnings = %d, want %d: %+v", warns, tc.wantWarns, resp.Warnings)
			}
			if r2, err := generateSandbox(svc, false); err != nil || r2.Action != driving.GenerateActionNoOp {
				t.Errorf("second run: action=%v err=%v, want NoOp", r2.Action, err)
			}
		})
	}
}

// An external source override is not a catalogue feature: no third-party
// warning for it (the allowlist step already is the consent).
func TestGenerateDevcontainer_FeaturesHomebrew_ExternalOverrideNoWarning(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, "devcontainer:\n  enabled: true\n  featureSources:\n    allow:\n      - https://ghcr.io/stu-bell/devcontainer-features/homebrew\n  features:\n    homebrew:\n      enabled: true\n      source: https://ghcr.io/stu-bell/devcontainer-features/homebrew\n      version: \"1\"\n")
	resp, err := generateSandbox(svc, false)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	for _, w := range resp.Warnings {
		if w.Code == "LH-NFA-SEC-004" {
			t.Errorf("unexpected third-party warning for an explicit source override: %+v", w)
		}
	}
	feats, _ := devcontainerJSONMap(t, fs)["features"].(map[string]any)
	if _, ok := feats["https://ghcr.io/stu-bell/devcontainer-features/homebrew:1"]; !ok {
		t.Errorf("override not rendered: %v", feats)
	}
}

// The warning names the version that is rendered, also with a version
// override on the catalogue key.
func TestGenerateDevcontainer_FeaturesHomebrew_WarningNamesConfiguredVersion(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, "devcontainer:\n  enabled: true\n  features:\n    homebrew:\n      enabled: true\n      version: \"3\"\n")
	resp, err := generateSandbox(svc, false)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	feats, _ := devcontainerJSONMap(t, fs)["features"].(map[string]any)
	if _, ok := feats["ghcr.io/meaningful-ooo/devcontainer-features/homebrew:3"]; !ok {
		t.Fatalf("version override not rendered: %v", feats)
	}
	found := false
	for _, w := range resp.Warnings {
		if w.Code == "LH-NFA-SEC-004" && strings.Contains(w.Message, "ghcr.io/meaningful-ooo/devcontainer-features/homebrew:3") {
			found = true
		}
	}
	if !found {
		t.Errorf("warning must name homebrew:3: %+v", resp.Warnings)
	}
}

// java pinned to another version next to sdkman at the default render two
// keys; each key is a distinct feature request.
func TestGenerateDevcontainer_FeaturesJavaVersionAndSdkman_RenderTwoKeys(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, "devcontainer:\n  enabled: true\n  features:\n    java:\n      enabled: true\n      version: \"21\"\n    sdkman:\n      enabled: true\n")
	if _, err := generateSandbox(svc, false); err != nil {
		t.Fatalf("generate: %v", err)
	}
	feats, _ := devcontainerJSONMap(t, fs)["features"].(map[string]any)
	var got []string
	for k := range feats {
		got = append(got, k)
	}
	sort.Strings(got)
	want := []string{"ghcr.io/devcontainers/features/java:1", "ghcr.io/devcontainers/features/java:21"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("features = %v, want %v", got, want)
	}
}

// sshdLoginWarnings returns the LH-FA-DEV-006 login-service entries of a
// generate response.
func sshdLoginWarnings(ws []driving.WarningEntry) []driving.WarningEntry {
	var out []driving.WarningEntry
	for _, w := range ws {
		if w.Code == "LH-FA-DEV-006" && strings.Contains(w.Message, "login service sshd") {
			out = append(out, w)
		}
	}
	return out
}

// LH-FA-DEV-006 / LH-FA-CLI-005: an active `sshd` feature in the sandbox
// profile is disclosed with a warning that carries SecurityRelaxation (it
// stays visible under `--quiet`); without the profile, without the feature
// or with the feature switched off there is none. A version override keeps
// the feature (the warning is about the source).
func TestGenerateDevcontainer_SshdLoginWarning(t *testing.T) {
	t.Parallel()
	const on = "    sshd:\n      enabled: true\n"
	cases := map[string]struct {
		yaml    string
		sandbox bool
		want    int
	}{
		"sandbox profile with sshd":    {"devcontainer:\n  enabled: true\n  profile: sandbox\n  features:\n" + on, false, 1},
		"sandbox flag with sshd":       {"devcontainer:\n  enabled: true\n  features:\n" + on, true, 1},
		"sandbox with sshd version":    {"devcontainer:\n  enabled: true\n  profile: sandbox\n  features:\n    sshd:\n      enabled: true\n      version: \"2\"\n", false, 1},
		"default profile with sshd":    {"devcontainer:\n  enabled: true\n  features:\n" + on, false, 0},
		"sandbox without sshd":         {"devcontainer:\n  enabled: true\n  profile: sandbox\n  features:\n    node:\n      enabled: true\n", false, 0},
		"sandbox with sshd switch off": {"devcontainer:\n  enabled: true\n  profile: sandbox\n  features:\n    sshd:\n      enabled: false\n", false, 0},
		"sandbox with sshd unset":      {"devcontainer:\n  enabled: true\n  profile: sandbox\n  features:\n    sshd: {}\n", false, 0},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, tc.yaml)
			resp, err := generateSandbox(svc, tc.sandbox)
			if err != nil {
				t.Fatalf("generate: %v", err)
			}
			got := sshdLoginWarnings(resp.Warnings)
			if len(got) != tc.want {
				t.Fatalf("sshd login warnings = %d, want %d: %+v", len(got), tc.want, resp.Warnings)
			}
			for _, w := range got {
				if !w.SecurityRelaxation || w.Level != "warn" {
					t.Errorf("sshd warning must be a warn entry with SecurityRelaxation: %+v", w)
				}
				if !strings.Contains(w.Message, "port 2222") || !strings.Contains(w.Message, "sudo passwd") ||
					!strings.Contains(w.Message, "not an isolation boundary") {
					t.Errorf("sshd warning lacks port, password step or boundary statement: %q", w.Message)
				}
				// The measured state (sshd -T): password and key login on,
				// root login permitted; "password login only" would understate it.
				if !strings.Contains(w.Message, "key login") || !strings.Contains(w.Message, "root login") ||
					strings.Contains(w.Message, "password login only") || strings.Contains(w.Message, "only") {
					t.Errorf("sshd warning must name key login and root login and not claim password-only: %q", w.Message)
				}
			}
			body, _ := fs.ReadFile(devcontainerJSONPath())
			if strings.Contains(string(body), "2222") {
				t.Errorf("port 2222 must not be forwarded (forwardPorts stems from services only):\n%s", body)
			}
		})
	}
}

// The third-party notice of another catalogue feature stays an
// unflagged entry next to the flagged sshd warning, so `--quiet` drops
// only the former.
func TestGenerateDevcontainer_SshdLoginWarning_ThirdPartyStaysUnflagged(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, "devcontainer:\n  enabled: true\n  profile: sandbox\n  features:\n    sshd:\n      enabled: true\n    homebrew:\n      enabled: true\n")
	resp, err := generateSandbox(svc, false)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	third, login := 0, 0
	for _, w := range resp.Warnings {
		switch {
		case w.Code == "LH-NFA-SEC-004":
			third++
			if w.SecurityRelaxation {
				t.Errorf("third-party notice must not carry SecurityRelaxation: %+v", w)
			}
		case strings.Contains(w.Message, "login service sshd"):
			login++
			if !w.SecurityRelaxation {
				t.Errorf("login-service warning lost SecurityRelaxation: %+v", w)
			}
		}
	}
	if third != 1 || login != 1 {
		t.Errorf("third-party=%d login=%d, want 1/1: %+v", third, login, resp.Warnings)
	}
	if r2, err := generateSandbox(svc, false); err != nil || r2.Action != driving.GenerateActionNoOp {
		t.Errorf("second run: action=%v err=%v, want NoOp", r2.Action, err)
	}
}

const claudeCodeKey = "ghcr.io/anthropics/devcontainer-features/claude-code"

// claudeCodeThirdPartyWarnings returns the LH-NFA-SEC-004 notices that
// name the Anthropic publisher.
func claudeCodeThirdPartyWarnings(ws []driving.WarningEntry) []driving.WarningEntry {
	var out []driving.WarningEntry
	for _, w := range ws {
		if w.Code == "LH-NFA-SEC-004" && strings.Contains(w.Message, "third-party publisher anthropics") {
			out = append(out, w)
		}
	}
	return out
}

// LH-FA-DEV-003 / ADR-0020: `claude-code` renders the Anthropic feature at
// the measured tag `1` (a `version` override renders as given) and is
// reported as a third-party feature without SecurityRelaxation (the notice
// is dropped under `--quiet`); the second run is a NoOp.
func TestGenerateDevcontainer_ClaudeCodeRender(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		yaml string
		want string
	}{
		"default tag":      {"devcontainer:\n  enabled: true\n  features:\n    claude-code:\n      enabled: true\n", claudeCodeKey + ":1"},
		"version override": {"devcontainer:\n  enabled: true\n  features:\n    claude-code:\n      enabled: true\n      version: \"1.0\"\n", claudeCodeKey + ":1.0"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, tc.yaml)
			resp, err := generateSandbox(svc, false)
			if err != nil {
				t.Fatalf("generate: %v", err)
			}
			body, _ := fs.ReadFile(devcontainerJSONPath())
			feats := parseDevcontainerFeatures(t, body)
			if _, ok := feats[tc.want]; !ok || len(feats) != 1 {
				t.Errorf("features = %v, want exactly %q", feats, tc.want)
			}
			third := claudeCodeThirdPartyWarnings(resp.Warnings)
			if len(third) != 1 {
				t.Fatalf("third-party notices for anthropics = %d, want 1: %+v", len(third), resp.Warnings)
			}
			if third[0].SecurityRelaxation || !strings.Contains(third[0].Message, tc.want) {
				t.Errorf("notice must name %q and carry no SecurityRelaxation: %+v", tc.want, third[0])
			}
			if r2, err := generateSandbox(svc, false); err != nil || r2.Action != driving.GenerateActionNoOp {
				t.Errorf("second run: action=%v err=%v, want NoOp", r2.Action, err)
			}
		})
	}
}

// `claude-code` installs through npm; on the base image without Node.js its
// own install fails (measured), so `generate` tells the user to enable the
// `node` feature: one unflagged notice when claude-code is active and node
// is not; none with node, with claude-code switched off, or when
// claude-code comes from an own source.
func TestGenerateDevcontainer_ClaudeCodeNodeNotice(t *testing.T) {
	t.Parallel()
	const on = "    claude-code:\n      enabled: true\n"
	cases := map[string]struct {
		yaml string
		want int
	}{
		"claude-code alone":        {"devcontainer:\n  enabled: true\n  features:\n" + on, 1},
		"with node":                {"devcontainer:\n  enabled: true\n  features:\n" + on + "    node:\n      enabled: true\n", 0},
		"node switched off":        {"devcontainer:\n  enabled: true\n  features:\n" + on + "    node:\n      enabled: false\n", 1},
		"claude-code switched off": {"devcontainer:\n  enabled: true\n  features:\n    claude-code:\n      enabled: false\n", 0},
		"own source":               {"devcontainer:\n  enabled: true\n  featureSources:\n    allow:\n      - https://ghcr.io/me/features/claude-code\n  features:\n    claude-code:\n      enabled: true\n      source: https://ghcr.io/me/features/claude-code\n", 0},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, tc.yaml)
			resp, err := generateSandbox(svc, false)
			if err != nil {
				t.Fatalf("generate: %v", err)
			}
			got := 0
			for _, w := range resp.Warnings {
				if w.Code == "LH-FA-DEV-003" && strings.Contains(w.Message, "claude-code") {
					got++
					if w.SecurityRelaxation || w.Level != "warn" || !strings.Contains(w.Message, "devcontainer.features.node.enabled") {
						t.Errorf("node notice must be an unflagged warn naming the node feature: %+v", w)
					}
				}
			}
			if got != tc.want {
				t.Errorf("node notices = %d, want %d: %+v", got, tc.want, resp.Warnings)
			}
		})
	}
}

// LH-FA-DEV-009: with `claude-code` active the sandbox `remoteEnv` carries
// the reference of the host variable next to GIT_TOKEN, never a value.
// Without the feature the block stays byte-identical to the GIT_TOKEN-only
// form; without the sandbox profile there is no `remoteEnv`.
func TestGenerateDevcontainer_ClaudeCodeCredentialReference(t *testing.T) {
	t.Parallel()
	const gitOnly = "  \"remoteEnv\": {\n    \"GIT_TOKEN\": \"${localEnv:GIT_TOKEN}\"\n  },\n"
	const withAgent = "  \"remoteEnv\": {\n    \"GIT_TOKEN\": \"${localEnv:GIT_TOKEN}\",\n    \"ANTHROPIC_API_KEY\": \"${localEnv:ANTHROPIC_API_KEY}\"\n  },\n"
	const on = "    claude-code:\n      enabled: true\n"
	cases := map[string]struct {
		yaml string
		want string // exact remoteEnv block, "" = no remoteEnv key
	}{
		"sandbox with claude-code":      {"devcontainer:\n  enabled: true\n  profile: sandbox\n  features:\n" + on, withAgent},
		"sandbox with claude-code, off": {"devcontainer:\n  enabled: true\n  profile: sandbox\n  features:\n    claude-code:\n      enabled: false\n", gitOnly},
		"sandbox without the feature":   {"devcontainer:\n  enabled: true\n  profile: sandbox\n  features:\n    node:\n      enabled: true\n", gitOnly},
		"sandbox with version override": {"devcontainer:\n  enabled: true\n  profile: sandbox\n  features:\n    claude-code:\n      enabled: true\n      version: \"1.0\"\n", withAgent},
		"default profile":               {"devcontainer:\n  enabled: true\n  features:\n" + on, ""},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, tc.yaml)
			if _, err := generateSandbox(svc, false); err != nil {
				t.Fatalf("generate: %v", err)
			}
			raw, _ := fs.ReadFile(devcontainerJSONPath())
			body := string(raw)
			if tc.want == "" {
				if strings.Contains(body, "remoteEnv") || strings.Contains(body, "ANTHROPIC") {
					t.Errorf("no remoteEnv expected:\n%s", body)
				}
			} else if !strings.Contains(body, tc.want) {
				t.Errorf("remoteEnv block differs, want:\n%s\ngot:\n%s", tc.want, body)
			}
			if strings.Contains(body, "sk-ant") {
				t.Errorf("a key value must never be rendered:\n%s", body)
			}
			m := devcontainerJSONMap(t, fs)
			if env, ok := m["remoteEnv"].(map[string]any); ok {
				for k, v := range env {
					if s, _ := v.(string); !strings.HasPrefix(s, "${localEnv:") || !strings.HasSuffix(s, "}") {
						t.Errorf("remoteEnv[%s] = %q is not a ${localEnv:...} reference", k, s)
					}
				}
			}
		})
	}
}

// egressHostList reads the baked allowlist of the egress start script.
func egressHostList(t *testing.T, fs *fakeFS) []string {
	t.Helper()
	script, err := fs.ReadFile(egressInitPath())
	if err != nil {
		t.Fatalf("read egress script: %v", err)
	}
	for _, line := range strings.Split(string(script), "\n") {
		if rest, ok := strings.CutPrefix(line, `hosts="`); ok {
			return strings.Fields(strings.TrimSuffix(rest, `"`))
		}
	}
	t.Fatalf("no hosts= line in the egress script:\n%s", script)
	return nil
}

// Welle-Closure (golden): `sshd` and `claude-code` (with `node`) together.
// Every render key and the agent reference appear exactly once, the second
// run is a NoOp, and with the egress restriction the allowlist is the
// sorted, duplicate-free union (registry.npmjs.org comes from `node` and
// `claude-code` and stays one entry; `sshd` adds none; optional hosts of
// the agent are absent).
func TestGenerateDevcontainer_SshdAndClaudeCodeTogether(t *testing.T) {
	t.Parallel()
	const features = "  features:\n    sshd:\n      enabled: true\n    claude-code:\n      enabled: true\n    node:\n      enabled: true\n"
	cases := map[string]struct {
		yaml      string
		sandbox   bool
		wantLogin int
		wantHosts []string
	}{
		"default profile": {"devcontainer:\n  enabled: true\n" + features, false, 0, nil},
		"sandbox":         {"devcontainer:\n  enabled: true\n  profile: sandbox\n" + features, false, 1, nil},
		"sandbox flag":    {"devcontainer:\n  enabled: true\n" + features, true, 1, nil},
		"sandbox egress": {"devcontainer:\n  enabled: true\n  profile: sandbox\n" + features + "  sandbox:\n    egress:\n      enabled: true\n", false, 1, []string{
			"api.anthropic.com", "api.github.com", "codeload.github.com", "deb.debian.org",
			"github.com", "objects.githubusercontent.com", "raw.githubusercontent.com",
			"registry.npmjs.org", "security.debian.org"}},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, tc.yaml)
			resp, err := generateSandbox(svc, tc.sandbox)
			if err != nil {
				t.Fatalf("generate: %v", err)
			}
			raw, _ := fs.ReadFile(devcontainerJSONPath())
			body := string(raw)
			for _, key := range []string{
				"ghcr.io/devcontainers/features/sshd:1", claudeCodeKey + ":1", "ghcr.io/devcontainers/features/node:1"} {
				if n := strings.Count(body, `"`+key+`"`); n != 1 {
					t.Errorf("render key %q appears %d times, want 1:\n%s", key, n, body)
				}
			}
			if n := len(parseDevcontainerFeatures(t, raw)); n != 3 {
				t.Errorf("features = %d keys, want 3", n)
			}
			wantAgent := 0
			if tc.sandbox || strings.Contains(tc.yaml, "profile: sandbox") {
				wantAgent = 1
			}
			if n := strings.Count(body, `"ANTHROPIC_API_KEY"`); n != wantAgent {
				t.Errorf("ANTHROPIC_API_KEY reference appears %d times, want %d", n, wantAgent)
			}
			if got := len(sshdLoginWarnings(resp.Warnings)); got != tc.wantLogin {
				t.Errorf("login warnings = %d, want %d: %+v", got, tc.wantLogin, resp.Warnings)
			}
			if got := len(claudeCodeThirdPartyWarnings(resp.Warnings)); got != 1 {
				t.Errorf("anthropics notices = %d, want 1: %+v", got, resp.Warnings)
			}
			for _, w := range resp.Warnings {
				if w.Code == "LH-FA-DEV-003" && strings.Contains(w.Message, "claude-code") {
					t.Errorf("node is enabled, no node notice expected: %+v", w)
				}
			}
			if tc.wantHosts != nil {
				hosts := egressHostList(t, fs)
				if !sort.StringsAreSorted(hosts) || !reflect.DeepEqual(hosts, tc.wantHosts) {
					t.Errorf("egress hosts =\n%v\nwant (sorted, no duplicates)\n%v", hosts, tc.wantHosts)
				}
			}
			if r2, err := generateSandbox(svc, tc.sandbox); err != nil || r2.Action != driving.GenerateActionNoOp {
				t.Errorf("second run: action=%v err=%v, want NoOp", r2.Action, err)
			}
		})
	}
}
