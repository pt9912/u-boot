package application_test

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
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
// is a NoOp; comma-separated values split per LH-FA-DEV-003.
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
// the positional Value (LH-FA-DEV-003 "Multi-Flag-Vorkommen kumulieren").
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

// LH-FA-DEV-006, the way out of the rejection: in a project with the
// sandbox profile and an active `docker-cli`, `generate devcontainer` ends
// with the business error, `config set …docker-cli.enabled false` still
// succeeds (the rejection is not a rule of the schema check that `config`
// runs when it loads u-boot.yaml), and `generate devcontainer` succeeds
// afterwards.
func TestConfigSet_SandboxProfileWithDockerCLI_DisablingStillWorks(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAMLWithDevcontainer(t, fs, "schemaVersion: 1\nproject:\n  name: t-uboot-config\n"+
		"devcontainer:\n  enabled: true\n  profile: sandbox\n  features:\n    docker-cli:\n      enabled: true\n")
	gen := application.NewGenerateService(fs, &fakeYAML{}, nil)
	genReq := driving.GenerateRequest{BaseDir: configTestBaseDir, Artifact: domain.ArtifactDevcontainer}

	if _, err := gen.Generate(context.Background(), genReq); !errors.Is(err, domain.ErrInvalidSandboxSetting) {
		t.Fatalf("generate before the way out: err = %v, want domain.ErrInvalidSandboxSetting", err)
	}
	resp, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "devcontainer.features.docker-cli.enabled"),
		Value:   "false",
	})
	if err != nil {
		t.Fatalf("config set …docker-cli.enabled false: %v", err)
	}
	if resp.OldValue != "true" || resp.NewValue != "false" {
		t.Errorf("config set: %q -> %q, want true -> false", resp.OldValue, resp.NewValue)
	}
	if _, err := gen.Generate(context.Background(), genReq); err != nil {
		t.Errorf("generate after the way out: %v", err)
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

// TestConfigSet_FeatureEnabled_OrphanWarning pins the slice-v1-
// devcontainer-features Review-Followup R3 fix: activating a
// feature that is neither in the catalogue nor carries a `source:`
// override surfaces an Info-log line. The check is non-fatal — the
// write still succeeds — but the user sees a clear hint that
// `generate devcontainer` will skip the entry until T5 doctor
// flags it.
func TestConfigSet_FeatureEnabled_OrphanWarning(t *testing.T) {
	t.Parallel()
	fs := newFakeFS()
	fs.markDirExists(configTestBaseDir)
	y := &fakeYAML{}
	logger := &fakeLogger{}
	svc := application.NewConfigService(fs, y, logger)
	seedConfigUbootYAML(t, fs)

	_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "devcontainer.features.unknown-thing.enabled"),
		Value:   "true",
	})
	if err != nil {
		t.Fatalf("Set: %v", err)
	}
	found := false
	for _, entry := range logger.entries {
		if entry.Level == "INFO" && strings.Contains(entry.Msg, "orphan feature activation") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected orphan-feature INFO log; got entries: %+v", logger.entries)
	}
}

// TestConfigSet_FeatureEnabled_DeactivationNoWarning pins that
// `config set <orphan>.enabled false` does NOT trigger the orphan
// warning. The Info-line is reserved for the "user activates an
// unknown feature" case; deactivations are routine cleanups.
func TestConfigSet_FeatureEnabled_DeactivationNoWarning(t *testing.T) {
	t.Parallel()
	fs := newFakeFS()
	fs.markDirExists(configTestBaseDir)
	y := &fakeYAML{}
	logger := &fakeLogger{}
	svc := application.NewConfigService(fs, y, logger)
	// Fixture already has an orphan with enabled=true so the
	// deactivation has something to flip without auto-creating.
	seedConfigUbootYAMLWithDevcontainer(t, fs, `schemaVersion: 1
project:
  name: t-uboot-config
devcontainer:
  features:
    unknown-thing:
      enabled: true
`)
	if _, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "devcontainer.features.unknown-thing.enabled"),
		Value:   "false",
	}); err != nil {
		t.Fatalf("Set: %v", err)
	}
	for _, entry := range logger.entries {
		if entry.Level == "INFO" && strings.Contains(entry.Msg, "orphan feature activation") {
			t.Errorf("did not expect orphan warning on deactivation; entry: %+v", entry)
		}
	}
}

// TestConfigSet_FeatureEnabled_CataloguedNoWarning pins that a
// catalogued name does NOT trigger the orphan warning.
func TestConfigSet_FeatureEnabled_CataloguedNoWarning(t *testing.T) {
	t.Parallel()
	fs := newFakeFS()
	fs.markDirExists(configTestBaseDir)
	y := &fakeYAML{}
	logger := &fakeLogger{}
	svc := application.NewConfigService(fs, y, logger)
	seedConfigUbootYAML(t, fs)

	if _, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "devcontainer.features.node.enabled"),
		Value:   "true",
	}); err != nil {
		t.Fatalf("Set: %v", err)
	}
	for _, entry := range logger.entries {
		if entry.Level == "INFO" && strings.Contains(entry.Msg, "orphan feature activation") {
			t.Errorf("did not expect orphan warning for catalogued `node`; entry: %+v", entry)
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

// seedForgeVersionWithoutSource is a u-boot.yaml whose `gh` or `glab` entry
// carries a `version:` without `source:`, the shape every command of the
// config family and `doctor` and `generate devcontainer` reject.
func seedForgeVersionWithoutSource(t *testing.T, fs *fakeFS, key string) []byte {
	t.Helper()
	body := []byte("schemaVersion: 1\nproject:\n  name: t-uboot-config\n" +
		"devcontainer:\n  enabled: true\n  features:\n    " + key + ":\n      enabled: true\n      version: \"9.9.9\"\n")
	seedConfigUbootYAMLWithDevcontainer(t, fs, string(body))
	return body
}

// A `version:` without `source:` on `gh` or `glab` makes every further
// `config set` end with ErrConfigSchemaInvalid (exit code 10) and leaves the
// file byte-identical; the way out is to remove the `version:` entry by hand.
func TestConfigSet_ForgeVersionWithoutSource_RejectsEverySet(t *testing.T) {
	t.Parallel()
	for _, key := range []string{"gh", "glab"} {
		for _, set := range []struct{ path, value string }{
			{"project.name", "other-name"},
			{"devcontainer.features." + key + ".enabled", "false"},
			{"devcontainer.features." + key + ".version", "8.8.8"},
			{"devcontainer.featureSources.allow", "https://example.test/x"},
		} {
			t.Run(key+"/"+set.path, func(t *testing.T) {
				t.Parallel()
				svc, fs := newConfigService(t)
				before := seedForgeVersionWithoutSource(t, fs, key)
				_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
					BaseDir: configTestBaseDir,
					Path:    mustConfigPath(t, set.path),
					Value:   set.value,
				})
				if !errors.Is(err, driving.ErrConfigSchemaInvalid) {
					t.Fatalf("Set error = %v, want a wrap of ErrConfigSchemaInvalid (exit code 10)", err)
				}
				after, _ := fs.ReadFile(filepath.Join(configTestBaseDir, "u-boot.yaml"))
				if string(after) != string(before) {
					t.Errorf("u-boot.yaml changed by a rejected Set:\n%s", after)
				}
			})
		}
	}
}

// `config set …<feature>.version <value>` follows `generate devcontainer`: a
// Dockerfile-block feature (`ollama`, `gh`, `glab`) takes a `version` only with
// a `source`. Without it the Set ends with ErrConfigValueInvalid (exit code
// 10, not the file-broken ErrConfigSchemaInvalid), the file stays
// byte-identical and nothing is written, in every preview mode.
func TestConfigSet_DockerfileBlockVersionWithoutSource_Rejected(t *testing.T) {
	t.Parallel()
	for _, key := range []string{"ollama", "gh", "glab"} {
		for _, mode := range []struct {
			name string
			mode driving.PreviewMode
		}{{"apply", driving.PreviewNone}, {"dry-run", driving.PreviewDryRun}} {
			t.Run(key+"/"+mode.name, func(t *testing.T) {
				t.Parallel()
				svc, fs := newConfigService(t)
				seedConfigUbootYAMLWithDevcontainer(t, fs, "schemaVersion: 1\nproject:\n  name: t-uboot-config\n"+
					"devcontainer:\n  enabled: true\n  features:\n    "+key+":\n      enabled: true\n")
				before, writes := readYAML(t, fs), len(fs.writes)
				_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
					BaseDir:     configTestBaseDir,
					Path:        mustConfigPath(t, "devcontainer.features."+key+".version"),
					Value:       "2.0.0",
					PreviewMode: mode.mode,
				})
				if !errors.Is(err, driving.ErrConfigValueInvalid) {
					t.Fatalf("Set error = %v, want a wrap of ErrConfigValueInvalid", err)
				}
				if errors.Is(err, driving.ErrConfigSchemaInvalid) {
					t.Errorf("Set error = %v, must not claim a broken file", err)
				}
				if !strings.Contains(err.Error(), "devcontainer.features."+key+".source") {
					t.Errorf("Set error = %v, want the source path as the way out", err)
				}
				if readYAML(t, fs) != before || len(fs.writes) != writes {
					t.Errorf("u-boot.yaml changed by a rejected Set:\n%s", readYAML(t, fs))
				}
			})
		}
	}
}

// With an allowlisted `source` in the entry the version is set; a feature
// without a Dockerfile block (`java`) keeps its `version` without a source.
func TestConfigSet_VersionWithSourceOrWithoutDockerfileBlock_Succeeds(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAMLWithDevcontainer(t, fs, "schemaVersion: 1\nproject:\n  name: t-uboot-config\n"+
		"devcontainer:\n  enabled: true\n  featureSources:\n    allow:\n      - https://example.test/features/gh\n"+
		"  features:\n    gh:\n      enabled: true\n      source: https://example.test/features/gh\n")
	for _, set := range []struct{ path, value string }{
		{"devcontainer.features.gh.version", "2.0.0"},
		{"devcontainer.features.java.version", "21"},
	} {
		resp, err := svc.Set(context.Background(), driving.ConfigSetRequest{
			BaseDir: configTestBaseDir, Path: mustConfigPath(t, set.path), Value: set.value})
		if err != nil {
			t.Fatalf("Set %s: %v", set.path, err)
		}
		if resp.NewValue != set.value {
			t.Errorf("Set %s: NewValue = %q, want %q", set.path, resp.NewValue, set.value)
		}
	}
	if y := readYAML(t, fs); !strings.Contains(y, "2.0.0") || !strings.Contains(y, "21") {
		t.Errorf("u-boot.yaml does not carry both versions:\n%s", y)
	}
}
