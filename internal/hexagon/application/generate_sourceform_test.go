package application_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// LH-FA-DEV-003 (lastenheft 0.33.0): `generate devcontainer` rejects the
// source form of an enabled feature that the Dev-Container-CLI does not
// process. Measured 2026-11-03 with the CLI of the repos
// (docs/messungen/slice-v1-katalog-features-quellen-und-freigabe/): the
// registry-reference form is processed, `oci://…:<tag>` and
// `https://…:<tag>` are not. The rejection stands before the allowlist
// check and before any plan and write; the message carries the Abhilfe
// (the processed registry-reference form).

// sourceFormMeasuredForms returns the two scheme'd forms the M2
// measurement probed; both use a tag that exists, so the failure is the
// form, not the reference.
func sourceFormMeasuredForms() []string {
	return []string{
		"oci://ghcr.io/devcontainers/features/go",
		"https://ghcr.io/devcontainers/features/go",
	}
}

// The form rejection fires even when the source is listed byte-equal: the
// allowlist way out would not help, so the check stands before it. Nothing
// is written, also under `--dry-run`.
func TestGenerateDevcontainer_SourceFormRejectsMeasuredForms(t *testing.T) {
	t.Parallel()
	for _, src := range sourceFormMeasuredForms() {
		for _, mode := range []driving.PreviewMode{driving.PreviewNone, driving.PreviewDryRun} {
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, "devcontainer:\n  enabled: true\n  featureSources:\n    allow:\n      - "+src+"\n  features:\n    evil:\n      enabled: true\n      source: "+src+"\n")
			before, _ := fs.ReadFile(freigabeYAMLPath())
			writes, mkdirs := len(fs.writes), len(fs.mkdirs)

			_, err := generateWithFlag(svc, mode)
			if !errors.Is(err, driving.ErrConfigValueInvalid) {
				t.Fatalf("%s mode %v: err = %v, want ErrConfigValueInvalid", src, mode, err)
			}
			if len(fs.writes) != writes || len(fs.mkdirs) != mkdirs {
				t.Errorf("%s mode %v: rejected call wrote: writes %v, mkdirs %v", src, mode, fs.writes[writes:], fs.mkdirs[mkdirs:])
			}
			for p := range fs.files {
				if strings.Contains(p, ".devcontainer") {
					t.Errorf("%s mode %v: rejected call left %q", src, mode, p)
				}
			}
			if after, _ := fs.ReadFile(freigabeYAMLPath()); !bytes.Equal(before, after) {
				t.Errorf("%s mode %v: u-boot.yaml changed:\n%s", src, mode, after)
			}
		}
	}
}

// The form message comes before the freigabe message: for an unprocessable
// source that is also unlisted, the error names the registry-reference
// Abhilfe and not the allow-listing way out.
func TestGenerateDevcontainer_SourceFormBeforeFreigabe(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, "devcontainer:\n  enabled: true\n  features:\n    evil:\n      enabled: true\n      source: "+sourceFormMeasuredForms()[0]+"\n")
	_, err := generateWithFlag(svc, driving.PreviewNone)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, driving.ErrConfigValueInvalid) {
		t.Fatalf("err = %v, want ErrConfigValueInvalid", err)
	}
	msg := err.Error()
	if !strings.Contains(msg, "has a form the Dev-Container-CLI does not process") {
		t.Errorf("message lacks the form sentence: %s", msg)
	}
	if strings.Contains(msg, "is not in devcontainer.featureSources.allow") {
		t.Errorf("message carries the freigabe way out instead: %s", msg)
	}
}

// T8 of the freigabeliste pattern, applied to the form rejection: the
// message names feature, source, the form sentence and the Abhilfe, in
// this order; each part has its own pin.
func TestGenerateDevcontainer_SourceFormMessageParts(t *testing.T) {
	t.Parallel()
	const src = "https://ghcr.io/devcontainers/features/go"
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, "devcontainer:\n  enabled: true\n  featureSources:\n    allow:\n      - "+src+"\n  features:\n    evil:\n      enabled: true\n      source: "+src+"\n")
	_, err := generateWithFlag(svc, driving.PreviewNone)
	if err == nil {
		t.Fatal("expected error")
	}
	msg := err.Error()
	parts := []string{
		"devcontainer.features.evil: external source",
		`external source "` + src + `" has a form`,
		"the Dev-Container-CLI does not process",
		"declare the source as the registry reference host/path (the processed form)",
		"e.g. ghcr.io/devcontainers/features/go",
	}
	last := -1
	for _, p := range parts {
		at := strings.Index(msg, p)
		if at < 0 {
			t.Errorf("message lacks %q: %s", p, msg)
			continue
		}
		if at < last {
			t.Errorf("%q stands before an earlier part: %s", p, msg)
		}
		last = at
	}
	if strings.Contains(msg, "enabled false") {
		t.Errorf("message carries the disable way out: %s", msg)
	}
}

// Several unprocessable entries: one error, entries sorted by feature name
// and joined by "; " — the same message contract as
// TestGenerateDevcontainer_FreigabeSeveralUnlistedSortedByName.
func TestGenerateDevcontainer_SourceFormSeveralSortedByName(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, "devcontainer:\n  enabled: true\n  features:\n    zeta:\n      enabled: true\n      source: https://example.test/features/z\n    alpha:\n      enabled: true\n      source: oci://example.test/features/a\n")
	_, err := generateWithFlag(svc, driving.PreviewNone)
	if !errors.Is(err, driving.ErrConfigValueInvalid) {
		t.Fatalf("err = %v, want ErrConfigValueInvalid", err)
	}
	msg := err.Error()
	a := strings.Index(msg, "devcontainer.features.alpha: external source")
	z := strings.Index(msg, "; devcontainer.features.zeta: external source")
	if a < 0 || z < 0 || a > z {
		t.Errorf("entries not sorted by name and joined by \"; \": %s", msg)
	}
}

// A deactivated entry renders nothing, so its unprocessable source is no
// rejection (the same iteration contract as
// TestGenerateDevcontainer_FreigabeIgnoresInactiveEntries).
func TestGenerateDevcontainer_SourceFormIgnoresInactiveEntries(t *testing.T) {
	t.Parallel()
	for name, enabledLine := range map[string]string{"disabled": "      enabled: false\n", "absent": ""} {
		svc, fs := newGenerateService(t)
		seedUBootYAMLWithFeatures(t, fs, "devcontainer:\n  enabled: true\n  features:\n    evil:\n"+enabledLine+"      source: "+sourceFormMeasuredForms()[0]+"\n")
		if _, err := generateWithFlag(svc, driving.PreviewNone); err != nil {
			t.Fatalf("%s: generate: %v", name, err)
		}
		feats, _ := devcontainerJSONMap(t, fs)["features"].(map[string]any)
		if _, ok := feats[sourceFormMeasuredForms()[0]+":1"]; ok {
			t.Errorf("%s: inactive entry rendered: %v", name, feats)
		}
	}
}

// The Abhilfe of the message works: a source declared as the registry
// reference and listed byte-equal renders its feature key
// (`<source>:<version>`), the form the Dev-Container-CLI resolves.
func TestGenerateDevcontainer_SourceFormRegistryReferenceRenders(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, "devcontainer:\n  enabled: true\n  featureSources:\n    allow:\n      - ghcr.io/orgX/features/custom-rust\n  features:\n    custom-rust:\n      enabled: true\n      source: ghcr.io/orgX/features/custom-rust\n")
	if _, err := generateWithFlag(svc, driving.PreviewNone); err != nil {
		t.Fatalf("generate with registry-reference source: %v", err)
	}
	feats, _ := devcontainerJSONMap(t, fs)["features"].(map[string]any)
	if _, ok := feats["ghcr.io/orgX/features/custom-rust:1"]; !ok {
		t.Errorf("registry-reference source not rendered: %v", feats)
	}
}
