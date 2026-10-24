package application_test

import (
	"bytes"
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// LH-FA-DEV-003 / LH-NFA-SEC-004: `generate devcontainer` renders the source of
// an enabled feature only when the effective allowlist contains it
// byte-equal. The effective allowlist is `devcontainer.featureSources.allow`
// plus the `--allow-external-feature-sources` entries of the same call.

const (
	freigabeSource = "https://example.test/features/evil"

	// freigabeEntryYAML enables one external feature and carries no allowlist.
	freigabeEntryYAML = "devcontainer:\n  enabled: true\n  features:\n    evil:\n      enabled: true\n      source: " + freigabeSource + "\n"
	// freigabeListedYAML lists the source verbatim.
	freigabeListedYAML = "devcontainer:\n  enabled: true\n  featureSources:\n    allow:\n      - " + freigabeSource + "\n  features:\n    evil:\n      enabled: true\n      source: " + freigabeSource + "\n"
)

func generateWithFlag(svc *application.GenerateService, mode driving.PreviewMode, flag ...string) (driving.GenerateResponse, error) {
	return svc.Generate(context.Background(), driving.GenerateRequest{
		BaseDir: generateTestBaseDir, Artifact: domain.ArtifactDevcontainer,
		PreviewMode: mode, AllowExternalFeatureSources: flag,
	})
}

func freigabeYAMLPath() string { return filepath.Join(generateTestBaseDir, "u-boot.yaml") }

// T1: an enabled entry whose source is outside the effective allowlist ends in
// ErrConfigValueInvalid and changes nothing, with and without `--dry-run`.
func TestGenerateDevcontainer_FreigabeRejectsUnlistedSource(t *testing.T) {
	t.Parallel()
	for _, mode := range []driving.PreviewMode{driving.PreviewNone, driving.PreviewDryRun} {
		svc, fs := newGenerateService(t)
		seedUBootYAMLWithFeatures(t, fs, freigabeEntryYAML)
		before, _ := fs.ReadFile(freigabeYAMLPath())
		writes, mkdirs := len(fs.writes), len(fs.mkdirs)

		_, err := generateWithFlag(svc, mode)
		if !errors.Is(err, driving.ErrConfigValueInvalid) {
			t.Fatalf("mode %v: err = %v, want ErrConfigValueInvalid", mode, err)
		}
		if len(fs.writes) != writes || len(fs.mkdirs) != mkdirs {
			t.Errorf("mode %v: rejected call wrote: writes %v, mkdirs %v", mode, fs.writes[writes:], fs.mkdirs[mkdirs:])
		}
		for p := range fs.files {
			if strings.Contains(p, ".devcontainer") {
				t.Errorf("mode %v: rejected call left %q", mode, p)
			}
		}
		if after, _ := fs.ReadFile(freigabeYAMLPath()); !bytes.Equal(before, after) {
			t.Errorf("mode %v: u-boot.yaml changed:\n%s", mode, after)
		}
	}
}

// T3 boundary: a flag entry for another source does not free this one.
func TestGenerateDevcontainer_FreigabeRejectsAlsoWithUnrelatedFlag(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, freigabeEntryYAML)
	_, err := generateWithFlag(svc, driving.PreviewNone, "https://example.test/features/other")
	if !errors.Is(err, driving.ErrConfigValueInvalid) {
		t.Fatalf("err = %v, want ErrConfigValueInvalid (flag names another source)", err)
	}
}

// T2: the source stands byte-equal in the list: the key is rendered.
func TestGenerateDevcontainer_FreigabeListedSourceRenders(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, freigabeListedYAML)
	if _, err := generateWithFlag(svc, driving.PreviewNone); err != nil {
		t.Fatalf("generate: %v", err)
	}
	feats, _ := devcontainerJSONMap(t, fs)["features"].(map[string]any)
	if _, ok := feats[freigabeSource+":1"]; !ok {
		t.Errorf("listed source not rendered: %v", feats)
	}
}

// T3: the flag of the same call frees the source; afterwards the list carries
// it, and the next call without the flag rests on the list.
func TestGenerateDevcontainer_FreigabeByFlagOfSameCall(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, freigabeEntryYAML)
	if _, err := generateWithFlag(svc, driving.PreviewNone, freigabeSource); err != nil {
		t.Fatalf("generate with flag: %v", err)
	}
	body, _ := fs.ReadFile(freigabeYAMLPath())
	if !strings.Contains(string(body), "- "+freigabeSource) {
		t.Errorf("flag source not appended to the list:\n%s", body)
	}
	if _, err := generateWithFlag(svc, driving.PreviewNone); err != nil {
		t.Errorf("second call without flag: %v", err)
	}
}

// T3: under `--dry-run` the flag frees the call as well. That the preview
// leaves the file untouched is the capture file system of the wiring
// (cmd/uboot/main_test.go).
func TestGenerateDevcontainer_FreigabeByFlagDryRun(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, freigabeEntryYAML)
	if _, err := generateWithFlag(svc, driving.PreviewDryRun, freigabeSource); err != nil {
		t.Fatalf("dry-run with flag: %v", err)
	}
}

// T4: a disabled entry and an entry without `enabled` render nothing, so an
// unlisted source there is no rejection (the way out "disable").
func TestGenerateDevcontainer_FreigabeIgnoresInactiveEntries(t *testing.T) {
	t.Parallel()
	for name, enabledLine := range map[string]string{"disabled": "      enabled: false\n", "absent": ""} {
		svc, fs := newGenerateService(t)
		seedUBootYAMLWithFeatures(t, fs, "devcontainer:\n  enabled: true\n  features:\n    evil:\n"+enabledLine+"      source: "+freigabeSource+"\n")
		if _, err := generateWithFlag(svc, driving.PreviewNone); err != nil {
			t.Fatalf("%s: generate: %v", name, err)
		}
		feats, _ := devcontainerJSONMap(t, fs)["features"].(map[string]any)
		if _, ok := feats[freigabeSource+":1"]; ok {
			t.Errorf("%s: inactive entry rendered: %v", name, feats)
		}
	}
}

// T6: two violations at once (unlisted source and `docker-cli` in the sandbox
// profile): the allowlist rejection comes first and carries ErrConfigValueInvalid,
// not ErrInvalidSandboxSetting.
func TestGenerateDevcontainer_FreigabeBeforeSandboxRejection(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, "devcontainer:\n  enabled: true\n  profile: sandbox\n  features:\n    docker-cli:\n      enabled: true\n    evil:\n      enabled: true\n      source: "+freigabeSource+"\n")
	_, err := generateWithFlag(svc, driving.PreviewNone)
	if !errors.Is(err, driving.ErrConfigValueInvalid) {
		t.Fatalf("err = %v, want ErrConfigValueInvalid", err)
	}
	if errors.Is(err, domain.ErrInvalidSandboxSetting) {
		t.Errorf("err carries the sandbox rejection: %v", err)
	}
}

// T8: the message names feature, source, the ways out in this order (allow,
// then disable) and the byte-equal comparison; each part has its own pin.
func TestGenerateDevcontainer_FreigabeMessageParts(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, freigabeEntryYAML)
	_, err := generateWithFlag(svc, driving.PreviewNone)
	if err == nil {
		t.Fatal("expected error")
	}
	parts := []struct{ name, want string }{
		{"feature", "devcontainer.features.evil: external source"},
		{"source", `external source "` + freigabeSource + `" is not in`},
		{"allowlist", "is not in devcontainer.featureSources.allow;"},
		{"abhilfe 1 (config set)", "`u-boot config set devcontainer.featureSources.allow " + freigabeSource + "`"},
		{"abhilfe 1 (flag)", "`--allow-external-feature-sources " + freigabeSource + "`"},
		{"abhilfe 2 (disable)", "`u-boot config set devcontainer.features.evil.enabled false`"},
		{"byte-equal", "Comparison is byte-equal; trailing slashes and host case must match"},
	}
	msg := err.Error()
	last := -1
	for _, p := range parts {
		at := strings.Index(msg, p.want)
		if at < 0 {
			t.Errorf("%s: message lacks %q: %s", p.name, p.want, msg)
			continue
		}
		if strings.HasPrefix(p.name, "abhilfe") || p.name == "byte-equal" {
			if at < last {
				t.Errorf("%s: %q stands before an earlier part: %s", p.name, p.want, msg)
			}
			last = at
		}
	}
}

// T5: generate, `config set …source` and `doctor` judge the same cases the
// same way: one predicate behind all three.
func TestFeatureSourceAllow_ThreeCallersAgree(t *testing.T) {
	t.Parallel()
	const base = "https://Example.test/features/x"
	cases := []struct {
		name  string
		src   string // the source of the feature entry
		allow string // lines under `devcontainer:` for the allowlist, may be empty
		want  bool   // source allowed
	}{
		{"same", base, "  featureSources:\n    allow:\n      - " + base + "\n", true},
		{"trailing slash in list", base, "  featureSources:\n    allow:\n      - " + base + "/\n", false},
		{"trailing slash in source", base + "/", "  featureSources:\n    allow:\n      - " + base + "\n", false},
		{"host case", base, "  featureSources:\n    allow:\n      - https://example.test/features/x\n", false},
		{"empty list", base, "  featureSources:\n    allow: []\n", false},
		{"no block", base, "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			src := tc.src
			head := "devcontainer:\n  enabled: true\n" + tc.allow

			gsvc, gfs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, gfs, head+"  features:\n    x:\n      enabled: true\n      source: "+src+"\n")
			_, gerr := generateWithFlag(gsvc, driving.PreviewNone)
			if got := gerr == nil; got != tc.want {
				t.Errorf("generate allowed = %v (err %v), want %v", got, gerr, tc.want)
			}

			csvc, cfs := newConfigService(t)
			seedConfigUbootYAMLWithDevcontainer(t, cfs, "schemaVersion: 1\nproject:\n  name: t-uboot-config\n"+head)
			_, cerr := csvc.Set(context.Background(), driving.ConfigSetRequest{
				BaseDir: configTestBaseDir,
				Path:    mustConfigPath(t, "devcontainer.features.x.source"),
				Value:   src,
			})
			if got := cerr == nil; got != tc.want {
				t.Errorf("config set allowed = %v (err %v), want %v", got, cerr, tc.want)
			}

			dsvc, dfs, _, _, _ := newDoctorService(t)
			seedDoctorUbootYAMLFeatures(t, dfs, "schemaVersion: 1\nproject:\n  name: demo\n"+head+"  features:\n    x:\n      enabled: true\n      source: "+src+"\n")
			resp, derr := dsvc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
			if derr != nil {
				t.Fatalf("doctor: %v", derr)
			}
			d := findDiagnostic(t, resp.Report.Items, "devcontainer.features.allowlist")
			if got := d.Severity != domain.SeverityError; got != tc.want {
				t.Errorf("doctor allowed = %v (%v: %s), want %v", got, d.Severity, d.Message, tc.want)
			}
		})
	}
}
