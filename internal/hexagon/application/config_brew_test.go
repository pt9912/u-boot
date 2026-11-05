package application_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// `config set` and `config get` of the two lists of the Homebrew block
// (LH-FA-DEV-013): `devcontainer.features.brew.packages` and
// `devcontainer.brewTaps.allow`. Every rule of the read path of `generate
// devcontainer` holds on the write path too (BEO-APP/pruefung-nur-beim-lesen):
// the form of the entries, the placement of the list, the approval of the taps
// and the conflict with the OCI feature `homebrew`.

const brewConfigHead = "schemaVersion: 1\nproject:\n  name: t-uboot-config\ndevcontainer:\n  enabled: true\n"

func setBrew(t *testing.T, svc driving.ConfigUseCase, path, value string, mode driving.PreviewMode) (driving.ConfigSetResponse, error) {
	t.Helper()
	return svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir, Path: mustConfigPath(t, path), Value: value, PreviewMode: mode})
}

// P1 (Listenpfad hängt an): the list path appends and drops duplicates in
// first-seen order, the second call with known entries is a NoOp, and the
// comma-separated value splits per entry.
func TestConfigSet_BrewPackages_AppendsLikeTheOtherLists(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAMLWithDevcontainer(t, fs, brewConfigHead)
	const path = "devcontainer.features.brew.packages"
	resp, err := setBrew(t, svc, path, "ripgrep, jq", driving.PreviewNone)
	if err != nil {
		t.Fatalf("Set: %v", err)
	}
	if resp.OldValue != "" || resp.NewValue != "ripgrep,jq" {
		t.Errorf("first Set old=%q new=%q, want \"\" and ripgrep,jq", resp.OldValue, resp.NewValue)
	}
	resp, err = setBrew(t, svc, path, "bat,jq,ripgrep", driving.PreviewNone)
	if err != nil {
		t.Fatalf("second Set: %v", err)
	}
	if resp.OldValue != "ripgrep,jq" || resp.NewValue != "ripgrep,jq,bat" {
		t.Errorf("second Set old=%q new=%q, want the append of bat only", resp.OldValue, resp.NewValue)
	}
	writes := len(fs.writes)
	resp, err = setBrew(t, svc, path, "jq", driving.PreviewNone)
	if err != nil || resp.OldValue != resp.NewValue || len(fs.writes) != writes {
		t.Errorf("a known entry must be a NoOp: resp=%+v err=%v writes %d -> %d", resp, err, writes, len(fs.writes))
	}
	got, err := svc.Get(context.Background(), driving.ConfigGetRequest{BaseDir: configTestBaseDir, Path: mustConfigPath(t, path)})
	if err != nil || got.Value != "ripgrep,jq,bat" {
		t.Errorf("Get = %q, %v, want ripgrep,jq,bat", got.Value, err)
	}
	if y := readYAML(t, fs); !strings.Contains(y, "packages:") || !strings.Contains(y, "- bat") {
		t.Errorf("u-boot.yaml lacks the list:\n%s", y)
	}
}

// P2 (Get of the unset list): NotSet with the config-set hint, like the other
// list paths; an empty list is unset.
func TestConfigGet_BrewLists_NotSet(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAMLWithDevcontainer(t, fs, brewConfigHead+"  features:\n    brew:\n      enabled: true\n      packages: []\n")
	for _, p := range []string{"devcontainer.features.brew.packages", "devcontainer.brewTaps.allow", "devcontainer.features.node.packages"} {
		_, err := svc.Get(context.Background(), driving.ConfigGetRequest{BaseDir: configTestBaseDir, Path: mustConfigPath(t, p)})
		if !errors.Is(err, driving.ErrConfigValueNotSet) {
			t.Errorf("Get %s: err = %v, want ErrConfigValueNotSet", p, err)
		}
	}
}

// F1 (Form, Schreibweg): an entry outside the form ends with
// ErrConfigValueInvalid (exit code 10), nothing is written, in every preview
// mode; a valid entry in the same value does not get through.
func TestConfigSet_BrewLists_FormIsCheckedOnTheWritePath(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ path, value string }{
		{"devcontainer.features.brew.packages", "rip grep"},
		{"devcontainer.features.brew.packages", "ripgrep,bad;entry"},
		{"devcontainer.features.brew.packages", "-HEAD"},
		{"devcontainer.features.brew.packages", "pt9912/u-boot"},
		{"devcontainer.features.brew.packages", "a/b/c/d"},
		{"devcontainer.features.brew.packages", "$(id)"},
		{"devcontainer.features.brew.packages", "ripgrep,,jq"},
		{"devcontainer.brewTaps.allow", "pt9912"},
		{"devcontainer.brewTaps.allow", "pt9912/u-boot/u-boot"},
		{"devcontainer.brewTaps.allow", "pt9912/u boot"},
		{"devcontainer.brewTaps.allow", "pt9912/u-boot,bad"},
	} {
		for _, mode := range []driving.PreviewMode{driving.PreviewNone, driving.PreviewDryRun, driving.PreviewAndApply} {
			svc, fs := newConfigService(t)
			seedConfigUbootYAMLWithDevcontainer(t, fs, brewConfigHead)
			before, writes := readYAML(t, fs), len(fs.writes)
			_, err := setBrew(t, svc, tc.path, tc.value, mode)
			if !errors.Is(err, driving.ErrConfigValueInvalid) {
				t.Fatalf("%s=%q mode %v: err = %v, want ErrConfigValueInvalid", tc.path, tc.value, mode, err)
			}
			if readYAML(t, fs) != before || len(fs.writes) != writes {
				t.Errorf("%s=%q mode %v: u-boot.yaml changed by a rejected Set", tc.path, tc.value, mode)
			}
		}
	}
}

// F2 (Platzierung, Schreibweg): the package list belongs to the key `brew`
// without a `source:`; any other key, also `homebrew`, and `brew` with a source
// end with ErrConfigValueInvalid and a message that names `brew`; the file stays
// byte-identical.
func TestConfigSet_BrewPackages_PlacementIsCheckedOnTheWritePath(t *testing.T) {
	t.Parallel()
	const external = "ghcr.io/example/features/brew"
	for name, tc := range map[string]struct{ path, yaml string }{
		"node":              {"devcontainer.features.node.packages", brewConfigHead + "  features:\n    node:\n      enabled: true\n"},
		"homebrew":          {"devcontainer.features.homebrew.packages", brewConfigHead},
		"brew with source":  {"devcontainer.features.brew.packages", brewConfigHead + "  featureSources:\n    allow:\n      - " + external + "\n  features:\n    brew:\n      enabled: true\n      source: " + external + "\n"},
		"unregistered node": {"devcontainer.features.node.packages", brewConfigHead},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newConfigService(t)
			seedConfigUbootYAMLWithDevcontainer(t, fs, tc.yaml)
			before, writes := readYAML(t, fs), len(fs.writes)
			_, err := setBrew(t, svc, tc.path, "ripgrep", driving.PreviewNone)
			if !errors.Is(err, driving.ErrConfigValueInvalid) || !strings.Contains(err.Error(), "brew") {
				t.Fatalf("err = %v, want ErrConfigValueInvalid naming brew", err)
			}
			if readYAML(t, fs) != before || len(fs.writes) != writes {
				t.Errorf("u-boot.yaml changed by a rejected Set:\n%s", readYAML(t, fs))
			}
		})
	}
}

// T1 (LH-FA-DEV-013, Tap nur nach Freigabe, Schreibweg): with the key `brew`
// enabled, a formula of a tap outside `devcontainer.brewTaps.allow` ends with
// ErrConfigValueInvalid and leaves the file byte-identical; the message names
// the command that approves the tap. The approval is byte-equal; with the tap
// approved the same call succeeds. A core formula needs none.
func TestConfigSet_BrewPackages_TapNeedsApprovalOnTheWritePath(t *testing.T) {
	t.Parallel()
	const enabled = brewConfigHead + "  features:\n    brew:\n      enabled: true\n"
	svc, fs := newConfigService(t)
	seedConfigUbootYAMLWithDevcontainer(t, fs, enabled)
	before, writes := readYAML(t, fs), len(fs.writes)
	const path = "devcontainer.features.brew.packages"
	for _, mode := range []driving.PreviewMode{driving.PreviewNone, driving.PreviewDryRun, driving.PreviewAndApply} {
		_, err := setBrew(t, svc, path, "ripgrep,pt9912/u-boot/u-boot", mode)
		if !errors.Is(err, driving.ErrConfigValueInvalid) {
			t.Fatalf("mode %v: err = %v, want ErrConfigValueInvalid", mode, err)
		}
		for _, part := range []string{`the tap "pt9912/u-boot"`, "`u-boot config set devcontainer.brewTaps.allow pt9912/u-boot`", "--yes does not substitute"} {
			if !strings.Contains(err.Error(), part) {
				t.Errorf("mode %v: message lacks %q: %v", mode, part, err)
			}
		}
		if readYAML(t, fs) != before || len(fs.writes) != writes {
			t.Errorf("mode %v: u-boot.yaml changed by a rejected Set:\n%s", mode, readYAML(t, fs))
		}
	}
	// The list of the feature sources approves no tap, a case variant neither.
	for _, other := range []string{"devcontainer.featureSources.allow", "devcontainer.brewTaps.allow"} {
		value := "pt9912/u-boot"
		if strings.Contains(other, "brewTaps") {
			value = "Pt9912/u-boot"
		}
		if _, err := setBrew(t, svc, other, value, driving.PreviewNone); err != nil {
			t.Fatalf("Set %s: %v", other, err)
		}
		if _, err := setBrew(t, svc, path, "pt9912/u-boot/u-boot", driving.PreviewNone); !errors.Is(err, driving.ErrConfigValueInvalid) {
			t.Fatalf("after %s=%s: err = %v, want the rejection to stand", other, value, err)
		}
	}
	if _, err := setBrew(t, svc, "devcontainer.brewTaps.allow", "pt9912/u-boot", driving.PreviewNone); err != nil {
		t.Fatalf("approve the tap: %v", err)
	}
	resp, err := setBrew(t, svc, path, "ripgrep,pt9912/u-boot/u-boot", driving.PreviewNone)
	if err != nil || resp.NewValue != "ripgrep,pt9912/u-boot/u-boot" {
		t.Fatalf("approved tap: resp=%+v err=%v", resp, err)
	}
}

// T2 (boundary): while `brew` is not enabled the list is checked for its form
// only; enabling the key asks the approval, so the pipeline cannot reach a
// state `generate devcontainer` rejects by the order of the two calls.
func TestConfigSet_BrewEnabling_AsksTheApprovalAndTheConflict(t *testing.T) {
	t.Parallel()
	const listed = brewConfigHead + "  features:\n    brew:\n      packages:\n        - pt9912/u-boot/u-boot\n"
	t.Run("list before the enabling", func(t *testing.T) {
		t.Parallel()
		svc, fs := newConfigService(t)
		seedConfigUbootYAMLWithDevcontainer(t, fs, brewConfigHead)
		if _, err := setBrew(t, svc, "devcontainer.features.brew.packages", "pt9912/u-boot/u-boot", driving.PreviewNone); err != nil {
			t.Fatalf("an inactive entry takes the list: %v", err)
		}
		before, writes := readYAML(t, fs), len(fs.writes)
		_, err := setBrew(t, svc, "devcontainer.features.brew.enabled", "true", driving.PreviewNone)
		if !errors.Is(err, driving.ErrConfigValueInvalid) || !strings.Contains(err.Error(), "devcontainer.brewTaps.allow") {
			t.Fatalf("enabling: err = %v, want ErrConfigValueInvalid naming the approval list", err)
		}
		if readYAML(t, fs) != before || len(fs.writes) != writes {
			t.Errorf("u-boot.yaml changed by a rejected enabling")
		}
		if _, err := setBrew(t, svc, "devcontainer.brewTaps.allow", "pt9912/u-boot", driving.PreviewNone); err != nil {
			t.Fatal(err)
		}
		if _, err := setBrew(t, svc, "devcontainer.features.brew.enabled", "true", driving.PreviewNone); err != nil {
			t.Fatalf("enabling after the approval: %v", err)
		}
	})
	t.Run("disabling is always possible", func(t *testing.T) {
		t.Parallel()
		svc, fs := newConfigService(t)
		seedConfigUbootYAMLWithDevcontainer(t, fs, strings.Replace(listed, "    brew:\n", "    brew:\n      enabled: true\n", 1))
		if _, err := setBrew(t, svc, "devcontainer.features.brew.enabled", "false", driving.PreviewNone); err != nil {
			t.Fatalf("a file in a state generate rejects can still switch the key off: %v", err)
		}
	})
	// A file with `brew` active, its tap not approved and `homebrew` active is
	// a state `generate devcontainer` rejects twice over; switching off either
	// key stays possible, because the written value decides and the sibling
	// key's violation is not this write's.
	t.Run("disabling at either key of a rejected file", func(t *testing.T) {
		t.Parallel()
		const rejected = brewConfigHead + "  features:\n    brew:\n      enabled: true\n      packages:\n        - pt9912/u-boot/u-boot\n    homebrew:\n      enabled: true\n"
		for _, path := range []string{"devcontainer.features.homebrew.enabled", "devcontainer.features.brew.enabled"} {
			// The key's own block, with the value `false` written into it.
			key := strings.TrimSuffix(strings.TrimPrefix(path, "devcontainer.features."), ".enabled")
			want := "    " + key + ":\n      enabled: false\n"
			svc, fs := newConfigService(t)
			seedConfigUbootYAMLWithDevcontainer(t, fs, rejected)
			if _, err := setBrew(t, svc, path, "false", driving.PreviewNone); err != nil {
				t.Errorf("Set %s false: %v", path, err)
			} else if y := readYAML(t, fs); !strings.Contains(y, want) {
				t.Errorf("Set %s false: u-boot.yaml lacks %q:\n%s", path, want, y)
			}
			svc, fs = newConfigService(t)
			seedConfigUbootYAMLWithDevcontainer(t, fs, rejected)
			item := driving.ConfigSetItem{Path: mustConfigPath(t, path), Value: "false"}
			if _, err := svc.SetMany(context.Background(), driving.ConfigSetManyRequest{BaseDir: configTestBaseDir, Items: []driving.ConfigSetItem{item}}); err != nil {
				t.Errorf("SetMany %s false: %v", path, err)
			} else if y := readYAML(t, fs); !strings.Contains(y, want) {
				t.Errorf("SetMany %s false: u-boot.yaml lacks %q:\n%s", path, want, y)
			}
		}
	})
	t.Run("conflict with homebrew in both directions", func(t *testing.T) {
		t.Parallel()
		for name, tc := range map[string]struct{ yaml, path string }{
			"enable homebrew next to brew": {brewConfigHead + "  features:\n    brew:\n      enabled: true\n    homebrew:\n      enabled: false\n", "devcontainer.features.homebrew.enabled"},
			"enable brew next to homebrew": {brewConfigHead + "  features:\n    brew:\n      enabled: false\n    homebrew:\n      enabled: true\n", "devcontainer.features.brew.enabled"},
		} {
			svc, fs := newConfigService(t)
			seedConfigUbootYAMLWithDevcontainer(t, fs, tc.yaml)
			before, writes := readYAML(t, fs), len(fs.writes)
			_, err := setBrew(t, svc, tc.path, "true", driving.PreviewNone)
			if !errors.Is(err, driving.ErrConfigValueInvalid) {
				t.Fatalf("%s: err = %v, want ErrConfigValueInvalid", name, err)
			}
			for _, part := range []string{"devcontainer.features.brew", "devcontainer.features.homebrew"} {
				if !strings.Contains(err.Error(), part) {
					t.Errorf("%s: message lacks %q: %v", name, part, err)
				}
			}
			if readYAML(t, fs) != before || len(fs.writes) != writes {
				t.Errorf("%s: u-boot.yaml changed by a rejected Set", name)
			}
		}
	})
	t.Run("an enabling next to the other key off is fine", func(t *testing.T) {
		t.Parallel()
		svc, fs := newConfigService(t)
		seedConfigUbootYAMLWithDevcontainer(t, fs, brewConfigHead+"  features:\n    homebrew:\n      enabled: false\n")
		if _, err := setBrew(t, svc, "devcontainer.features.brew.enabled", "true", driving.PreviewNone); err != nil {
			t.Fatalf("enabling brew: %v", err)
		}
	})
}

// S1 (several pairs in one call): the pair of the tap list stands before the
// pair of the formulas; the reverse order, and a pair that breaks a rule,
// end with ErrConfigValueInvalid and leave the file byte-identical (the call is
// atomic).
func TestConfigSetMany_BrewTapListBeforePackages(t *testing.T) {
	t.Parallel()
	const seed = brewConfigHead + "  features:\n    brew:\n      enabled: true\n"
	tap := driving.ConfigSetItem{Path: mustConfigPath(t, "devcontainer.brewTaps.allow"), Value: "pt9912/u-boot"}
	pkg := driving.ConfigSetItem{Path: mustConfigPath(t, "devcontainer.features.brew.packages"), Value: "ripgrep,pt9912/u-boot/u-boot"}
	svc, fs := newConfigService(t)
	seedConfigUbootYAMLWithDevcontainer(t, fs, seed)
	before, writes := readYAML(t, fs), len(fs.writes)
	if _, err := svc.SetMany(context.Background(), driving.ConfigSetManyRequest{BaseDir: configTestBaseDir, Items: []driving.ConfigSetItem{pkg, tap}}); !errors.Is(err, driving.ErrConfigValueInvalid) {
		t.Fatalf("packages before the tap list: err = %v, want ErrConfigValueInvalid", err)
	}
	if readYAML(t, fs) != before || len(fs.writes) != writes {
		t.Errorf("u-boot.yaml changed by a rejected SetMany:\n%s", readYAML(t, fs))
	}
	bad := driving.ConfigSetItem{Path: mustConfigPath(t, "devcontainer.features.brew.packages"), Value: "bad entry"}
	if _, err := svc.SetMany(context.Background(), driving.ConfigSetManyRequest{BaseDir: configTestBaseDir, Items: []driving.ConfigSetItem{tap, bad}}); !errors.Is(err, driving.ErrConfigValueInvalid) {
		t.Fatalf("a bad entry after a good pair: err = %v, want ErrConfigValueInvalid", err)
	}
	if readYAML(t, fs) != before || len(fs.writes) != writes {
		t.Errorf("u-boot.yaml changed by a rejected SetMany (atomic):\n%s", readYAML(t, fs))
	}
	if _, err := svc.SetMany(context.Background(), driving.ConfigSetManyRequest{BaseDir: configTestBaseDir, Items: []driving.ConfigSetItem{tap, pkg}}); err != nil {
		t.Fatalf("tap list before packages: %v", err)
	}
	y := readYAML(t, fs)
	for _, want := range []string{"brewTaps:", "- pt9912/u-boot", "- pt9912/u-boot/u-boot", "- ripgrep"} {
		if !strings.Contains(y, want) {
			t.Errorf("u-boot.yaml lacks %q:\n%s", want, y)
		}
	}
}

// S2 (several pairs, the conflict): the pair that switches a key off stands
// before the pair that enables the other one; the reverse order ends with
// ErrConfigValueInvalid and leaves the file byte-identical.
func TestConfigSetMany_BrewConflictOrder(t *testing.T) {
	t.Parallel()
	const seed = brewConfigHead + "  features:\n    homebrew:\n      enabled: true\n"
	off := driving.ConfigSetItem{Path: mustConfigPath(t, "devcontainer.features.homebrew.enabled"), Value: "false"}
	on := driving.ConfigSetItem{Path: mustConfigPath(t, "devcontainer.features.brew.enabled"), Value: "true"}
	svc, fs := newConfigService(t)
	seedConfigUbootYAMLWithDevcontainer(t, fs, seed)
	before, writes := readYAML(t, fs), len(fs.writes)
	if _, err := svc.SetMany(context.Background(), driving.ConfigSetManyRequest{BaseDir: configTestBaseDir, Items: []driving.ConfigSetItem{on, off}}); !errors.Is(err, driving.ErrConfigValueInvalid) {
		t.Fatalf("enabling before switching off: err = %v, want ErrConfigValueInvalid", err)
	}
	if readYAML(t, fs) != before || len(fs.writes) != writes {
		t.Errorf("u-boot.yaml changed by a rejected SetMany:\n%s", readYAML(t, fs))
	}
	if _, err := svc.SetMany(context.Background(), driving.ConfigSetManyRequest{BaseDir: configTestBaseDir, Items: []driving.ConfigSetItem{off, on}}); err != nil {
		t.Fatalf("switching off before enabling: %v", err)
	}
}

// R1 (read path of config): a file with a package list at another key, an entry
// outside the form or a bad tap is rejected on load by every call, `config get`
// included, with ErrConfigSchemaInvalid (exit code 10); the file stays as it is.
func TestConfig_BrewFileOutsideTheFormIsRejectedOnLoad(t *testing.T) {
	t.Parallel()
	for name, body := range map[string]string{
		"bad formula":         brewConfigHead + "  features:\n    brew:\n      packages:\n        - rip grep\n",
		"list at node":        brewConfigHead + "  features:\n    node:\n      packages:\n        - ripgrep\n",
		"bad tap in the list": brewConfigHead + "  brewTaps:\n    allow:\n      - pt9912\n",
	} {
		svc, fs := newConfigService(t)
		seedConfigUbootYAMLWithDevcontainer(t, fs, body)
		before := readYAML(t, fs)
		_, err := svc.Get(context.Background(), driving.ConfigGetRequest{BaseDir: configTestBaseDir, Path: mustConfigPath(t, "devcontainer.enabled")})
		if !errors.Is(err, driving.ErrConfigSchemaInvalid) {
			t.Errorf("%s: Get err = %v, want ErrConfigSchemaInvalid", name, err)
		}
		_, err = setBrew(t, svc, "devcontainer.enabled", "false", driving.PreviewNone)
		if !errors.Is(err, driving.ErrConfigSchemaInvalid) {
			t.Errorf("%s: Set err = %v, want ErrConfigSchemaInvalid", name, err)
		}
		if readYAML(t, fs) != before {
			t.Errorf("%s: u-boot.yaml changed", name)
		}
	}
}

// L1 (`config list`): the package list and the tap list appear with their
// comma-joined values, in the sorted path order.
func TestConfigList_ShowsTheBrewLists(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAMLWithDevcontainer(t, fs, brewConfigHead+"  brewTaps:\n    allow:\n      - pt9912/u-boot\n  features:\n    brew:\n      enabled: true\n      packages:\n        - ripgrep\n        - pt9912/u-boot/u-boot\n")
	resp, err := svc.List(context.Background(), driving.ConfigListRequest{BaseDir: configTestBaseDir})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	got := map[string]string{}
	for _, e := range resp.Entries {
		got[e.Path.String()] = e.Value
	}
	if got["devcontainer.features.brew.packages"] != "ripgrep,pt9912/u-boot/u-boot" || got["devcontainer.brewTaps.allow"] != "pt9912/u-boot" {
		t.Errorf("List = %v, want both brew lists", got)
	}
}

// W1 (path form): `<feature>.packages` is a placeholder path for every
// feature name; the application layer places it.
func TestConfigPath_PackagesParsesForEveryFeatureName(t *testing.T) {
	t.Parallel()
	for _, raw := range []string{"devcontainer.features.node.packages", "devcontainer.features.brew.packages"} {
		p, err := domain.NewConfigPath(raw)
		if err != nil || p.Kind != domain.ConfigDevcontainerFeaturePackages || !p.WriteAllowed {
			t.Errorf("NewConfigPath(%q) = %+v, %v", raw, p, err)
		}
	}
}
