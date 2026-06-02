package application_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"reflect"
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

// TestGenerateDevcontainer_AllowExternalFeatureSources_Append pins
// the LH-FA-DEV-003 / Spec §715 flag-wiring on `generate
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
