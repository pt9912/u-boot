package application_test

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
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
