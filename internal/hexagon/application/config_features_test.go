package application_test

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// seedConfigUbootYAMLWithDevcontainer writes a project fixture
// that has a devcontainer subtree seeded with an existing allowlist
// + an enabled feature. T4 ConfigService tests drive their Set
// pipelines against this fixture.
func seedConfigUbootYAMLWithDevcontainer(t *testing.T, fs *fakeFS, body string) {
	t.Helper()
	if err := fs.WriteFile(filepath.Join(configTestBaseDir, "u-boot.yaml"),
		[]byte(body), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
}

const fixtureUBootYAMLDevcontainerSeed = `schemaVersion: 1
project:
  name: t-uboot-config
devcontainer:
  enabled: true
  featureSources:
    allow:
      - https://ghcr.io/orgX/features/custom-rust
`

// TestConfigSet_FeatureSourcesAllow_Append pins the LH-FA-DEV-003
// list-path append + dedupe contract: setting the path with a new
// URL adds it to the existing list; re-setting with the same URL
// is a NoOp; comma-separated values split per Spec §718.
func TestConfigSet_FeatureSourcesAllow_Append(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAMLWithDevcontainer(t, fs, fixtureUBootYAMLDevcontainerSeed)

	resp, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "devcontainer.featureSources.allow"),
		Value:   "https://example.test/x,https://example.test/y",
	})
	if err != nil {
		t.Fatalf("Set: %v", err)
	}
	wantOld := "https://ghcr.io/orgX/features/custom-rust"
	wantNew := "https://ghcr.io/orgX/features/custom-rust,https://example.test/x,https://example.test/y"
	if resp.OldValue != wantOld {
		t.Errorf("OldValue = %q, want %q", resp.OldValue, wantOld)
	}
	if resp.NewValue != wantNew {
		t.Errorf("NewValue = %q, want %q", resp.NewValue, wantNew)
	}

	// Second Set with the same URLs → NoOp.
	resp2, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "devcontainer.featureSources.allow"),
		Value:   "https://example.test/x",
	})
	if err != nil {
		t.Fatalf("second Set: %v", err)
	}
	if resp2.OldValue != resp2.NewValue {
		t.Errorf("second Set should be NoOp; got OldValue=%q NewValue=%q", resp2.OldValue, resp2.NewValue)
	}
}

// TestConfigSet_FeatureSourcesAllow_InvalidURL pins the validation
// failure path: a malformed entry in the new value rejects with
// ErrConfigValueInvalid (Code 10).
func TestConfigSet_FeatureSourcesAllow_InvalidURL(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAMLWithDevcontainer(t, fs, fixtureUBootYAMLDevcontainerSeed)

	_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "devcontainer.featureSources.allow"),
		Value:   "not-a-url",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, driving.ErrConfigValueInvalid) {
		t.Errorf("err = %v, want wrap of ErrConfigValueInvalid", err)
	}
}

// TestConfigSet_FeatureSourcesAllow_FlagMergesWithPositional pins
// the cumulative semantics of AllowExternalFeatureSources alongside
// the positional Value (Spec §718 "Multi-Flag-Vorkommen kumulieren").
func TestConfigSet_FeatureSourcesAllow_FlagMergesWithPositional(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAMLWithDevcontainer(t, fs, fixtureUBootYAMLDevcontainerSeed)

	resp, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir:                     configTestBaseDir,
		Path:                        mustConfigPath(t, "devcontainer.featureSources.allow"),
		Value:                       "https://positional.test/a",
		AllowExternalFeatureSources: []string{"https://flag.test/b", "https://flag.test/c"},
	})
	if err != nil {
		t.Fatalf("Set: %v", err)
	}
	for _, want := range []string{"positional.test/a", "flag.test/b", "flag.test/c"} {
		if !strings.Contains(resp.NewValue, want) {
			t.Errorf("NewValue = %q, expected to contain %q", resp.NewValue, want)
		}
	}
}

// TestConfigSet_FeatureSource_AllowlistEnforced pins the
// LH-FA-DEV-003 / LH-NFA-SEC-004 enforcement: setting
// `devcontainer.features.<name>.source` to a URL not in
// `featureSources.allow` fails with ErrConfigValueInvalid (Code 10)
// regardless of how the user runs the call. `--yes` does not
// substitute.
func TestConfigSet_FeatureSource_AllowlistEnforced(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAMLWithDevcontainer(t, fs, fixtureUBootYAMLDevcontainerSeed)

	_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "devcontainer.features.custom.source"),
		Value:   "https://uninvited.test/feature",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, driving.ErrConfigValueInvalid) {
		t.Errorf("err = %v, want wrap of ErrConfigValueInvalid", err)
	}
	if !strings.Contains(err.Error(), "not in devcontainer.featureSources.allow") {
		t.Errorf("err = %v, want allowlist-hint message", err)
	}
}

// TestConfigSet_FeatureSource_AllowlistOK pins the happy path —
// setting source to an URL that IS in the allowlist succeeds.
func TestConfigSet_FeatureSource_AllowlistOK(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAMLWithDevcontainer(t, fs, fixtureUBootYAMLDevcontainerSeed)

	resp, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "devcontainer.features.custom-rust.source"),
		Value:   "https://ghcr.io/orgX/features/custom-rust",
	})
	if err != nil {
		t.Fatalf("Set: %v", err)
	}
	if resp.NewValue != "https://ghcr.io/orgX/features/custom-rust" {
		t.Errorf("NewValue = %q", resp.NewValue)
	}
}

// TestConfigSet_FeatureEnabled pins the bool-scalar path: setting
// `devcontainer.features.<name>.enabled true` writes the bool with
// no allowlist or catalogue check (T3 handles unknown-name skip;
// T5 doctor surfaces the warn).
func TestConfigSet_FeatureEnabled(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAMLWithDevcontainer(t, fs, fixtureUBootYAMLDevcontainerSeed)

	resp, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "devcontainer.features.node.enabled"),
		Value:   "true",
	})
	if err != nil {
		t.Fatalf("Set: %v", err)
	}
	if resp.NewValue != "true" {
		t.Errorf("NewValue = %q, want true", resp.NewValue)
	}
}

// TestConfigSet_FeatureVersion_EmptyRejected pins that empty
// version pins are rejected (sentinel-empty would be ambiguous
// against the "missing key, use catalogue default" semantics).
func TestConfigSet_FeatureVersion_EmptyRejected(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAMLWithDevcontainer(t, fs, fixtureUBootYAMLDevcontainerSeed)

	_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "devcontainer.features.java.version"),
		Value:   "",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, driving.ErrConfigValueInvalid) {
		t.Errorf("err = %v, want ErrConfigValueInvalid", err)
	}
}

// TestConfigGet_FeatureSourcesAllow_NotSet pins the NotSet path:
// getting the allow list on a project without one returns
// ErrConfigValueNotSet so the CLI surfaces a clean message.
func TestConfigGet_FeatureSourcesAllow_NotSet(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAML(t, fs) // no devcontainer subtree

	_, err := svc.Get(context.Background(), driving.ConfigGetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "devcontainer.featureSources.allow"),
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, driving.ErrConfigValueNotSet) {
		t.Errorf("err = %v, want ErrConfigValueNotSet", err)
	}
}

// TestConfigGet_FeatureSourcesAllow_CommaJoined pins the canonical
// comma-joined Get format (mirrors Set input format for round-trip).
func TestConfigGet_FeatureSourcesAllow_CommaJoined(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAMLWithDevcontainer(t, fs, `schemaVersion: 1
project:
  name: t-uboot-config
devcontainer:
  enabled: true
  featureSources:
    allow:
      - https://a.test/x
      - https://b.test/y
`)
	resp, err := svc.Get(context.Background(), driving.ConfigGetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "devcontainer.featureSources.allow"),
	})
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	want := "https://a.test/x,https://b.test/y"
	if resp.Value != want {
		t.Errorf("Value = %q, want %q", resp.Value, want)
	}
}

// helper var to silence unused-import linter if test list shrinks
var _ = application.ErrInvalidFeatureSource

// TestConfigGet_FeatureEntry_HappyPaths pins the Get-side projection
// for the three scalar feature kinds, against a fixture that has the
// matching entry already set.
func TestConfigGet_FeatureEntry_HappyPaths(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAMLWithDevcontainer(t, fs, `schemaVersion: 1
project:
  name: t-uboot-config
devcontainer:
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
	cases := []struct {
		path string
		want string
	}{
		{"devcontainer.features.node.enabled", "true"},
		{"devcontainer.features.java.version", "21"},
		{"devcontainer.features.custom-rust.source", "https://ghcr.io/orgX/features/custom-rust"},
	}
	for _, tc := range cases {
		resp, err := svc.Get(context.Background(), driving.ConfigGetRequest{
			BaseDir: configTestBaseDir,
			Path:    mustConfigPath(t, tc.path),
		})
		if err != nil {
			t.Errorf("Get(%q): %v", tc.path, err)
			continue
		}
		if resp.Value != tc.want {
			t.Errorf("Get(%q) = %q, want %q", tc.path, resp.Value, tc.want)
		}
	}
}

// TestConfigGet_FeatureEntry_RegisteredButLeafEmpty pins the
// per-leaf NotSet branch: a feature entry that exists but has the
// requested leaf (.source/.version/.enabled) absent should still
// surface NotSet so the user gets a hint how to populate it.
func TestConfigGet_FeatureEntry_RegisteredButLeafEmpty(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	// `node` has only `enabled: true` — no source, no version.
	seedConfigUbootYAMLWithDevcontainer(t, fs, `schemaVersion: 1
project:
  name: t-uboot-config
devcontainer:
  features:
    node:
      enabled: true
`)
	cases := []string{
		"devcontainer.features.node.source",
		"devcontainer.features.node.version",
	}
	for _, raw := range cases {
		_, err := svc.Get(context.Background(), driving.ConfigGetRequest{
			BaseDir: configTestBaseDir,
			Path:    mustConfigPath(t, raw),
		})
		if err == nil {
			t.Errorf("Get(%q): expected NotSet error, got nil", raw)
			continue
		}
		if !errors.Is(err, driving.ErrConfigValueNotSet) {
			t.Errorf("Get(%q): err %v does not wrap ErrConfigValueNotSet", raw, err)
		}
	}
}

// TestConfigGet_FeatureEntry_NotSet pins the NotSet path for each
// feature scalar leaf — a project without the entry returns
// ErrConfigValueNotSet so the CLI surfaces the LH-FA-CLI-006 code-12
// path.
func TestConfigGet_FeatureEntry_NotSet(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAML(t, fs) // no devcontainer.features

	for _, raw := range []string{
		"devcontainer.features.node.enabled",
		"devcontainer.features.java.source",
		"devcontainer.features.go.version",
	} {
		_, err := svc.Get(context.Background(), driving.ConfigGetRequest{
			BaseDir: configTestBaseDir,
			Path:    mustConfigPath(t, raw),
		})
		if err == nil {
			t.Errorf("Get(%q): expected error, got nil", raw)
			continue
		}
		if !errors.Is(err, driving.ErrConfigValueNotSet) {
			t.Errorf("Get(%q): err %v does not wrap ErrConfigValueNotSet", raw, err)
		}
	}
}

// TestConfigSet_FeatureSourcesAllow_EmptyElement pins the
// comma-list parser rejects an empty entry between commas
// ("https://a,,https://b"). Without this, a typo would silently
// drop an element.
func TestConfigSet_FeatureSourcesAllow_EmptyElement(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAMLWithDevcontainer(t, fs, fixtureUBootYAMLDevcontainerSeed)

	_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "devcontainer.featureSources.allow"),
		Value:   "https://a.test/x,,https://b.test/y",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, driving.ErrConfigValueInvalid) {
		t.Errorf("err = %v, want ErrConfigValueInvalid", err)
	}
}

// TestConfigSet_FeatureSource_FeatureMissing pins the branch where
// the user sets `<name>.source` against a non-existent feature
// entry. PatchScalar creates the intermediate map; revalidate
// still sees the entry (post-patch), so the allowlist enforcement
// runs against a now-existing feature with a non-empty source.
// With an empty allowlist + a URL not in the allowlist → error.
func TestConfigSet_FeatureSource_FeatureMissing(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	// No devcontainer.featureSources.allow at all.
	seedConfigUbootYAML(t, fs)

	_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "devcontainer.features.custom.source"),
		Value:   "https://anywhere.test/x",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, driving.ErrConfigValueInvalid) {
		t.Errorf("err = %v, want ErrConfigValueInvalid (allowlist enforcement)", err)
	}
}
