package application_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// sandboxCfg is the subset of `devcontainer.sandbox` the cases vary.
type sandboxCfg struct {
	runtime string
	egress  bool
}

func relaxationFlags(ws []driving.WarningEntry) (flagged, unflagged []driving.WarningEntry) {
	for _, w := range ws {
		if w.SecurityRelaxation {
			flagged = append(flagged, w)
		} else {
			unflagged = append(unflagged, w)
		}
	}
	return flagged, unflagged
}

// LH-FA-CLI-005 / LH-FA-DEV-006/-007/-008: the "security relaxation"
// notices and the egress "has no effect" notice (no restriction arises)
// carry SecurityRelaxation (they stay visible under --quiet); the
// `nestedRuntime` "has no effect" notices stay unflagged. The flag sits
// on the entry, so the notices that share the code LH-FA-DEV-007 differ.
func TestSandboxRelaxationWarnings_FlagIsOnTheEntryNotTheCode(t *testing.T) {
	t.Parallel()
	podman := sandboxCfg{runtime: "podman"}
	docker := sandboxCfg{runtime: "docker"}
	egress := sandboxCfg{egress: true}

	cases := []struct {
		name        string
		dc          sandboxCfg
		sandbox     bool
		wantFlagged int
		wantPlain   int
		flaggedText string
	}{
		{"podman-sandbox: six relaxations", podman, true, 6, 0, "security relaxation"},
		{"docker-sandbox: --privileged", docker, true, 1, 0, "security relaxation"},
		{"egress-sandbox: NET_ADMIN", egress, true, 1, 0, "security relaxation"},
		{"podman-without-sandbox: has no effect", podman, false, 0, 1, ""},
		{"docker-without-sandbox: has no effect", docker, false, 0, 1, ""},
		{"egress-without-sandbox: has no effect, flagged", egress, false, 1, 0, "egress.enabled has no effect"},
		{"no nested runtime: nothing", sandboxCfg{}, true, 0, 0, ""},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			flagged, plain := relaxationFlags(application.SandboxRelaxationWarningsForTest(tc.dc.runtime, "", tc.dc.egress, tc.sandbox))
			if len(flagged) != tc.wantFlagged || len(plain) != tc.wantPlain {
				t.Fatalf("flagged=%d plain=%d, want %d/%d", len(flagged), len(plain), tc.wantFlagged, tc.wantPlain)
			}
			for _, w := range flagged {
				if !strings.Contains(w.Message, tc.flaggedText) {
					t.Errorf("flagged entry lacks %q: %q", tc.flaggedText, w.Message)
				}
			}
			for _, w := range plain {
				if !strings.Contains(w.Message, "has no effect") {
					t.Errorf("unflagged entry is not a no-effect notice: %q", w.Message)
				}
			}
		})
	}
}

// Same code, different flag: LH-FA-DEV-007 carries both a relaxation
// and a "has no effect" notice, so a Code-based filter is wrong.
func TestSandboxRelaxationWarnings_SameCodeBothFlags(t *testing.T) {
	t.Parallel()
	flagged, _ := relaxationFlags(application.SandboxRelaxationWarningsForTest("docker", "", false, true))
	_, plain := relaxationFlags(application.SandboxRelaxationWarningsForTest("docker", "", false, false))
	if len(flagged) != 1 || len(plain) != 1 || flagged[0].Code != plain[0].Code {
		t.Fatalf("want one flagged and one plain entry with the same code, got %+v / %+v", flagged, plain)
	}
}

// LH-FA-DEV-007 degradations (no daemon, `vfs` storage, no egress
// restriction) arise at container start and in `u-boot doctor`; no
// render-time warning describes one, so the policy value
// (`warn`/`fail`) changes the generated script, never the set of
// render-time warnings. Their disclosure is the
// [driving.WarningEntry.DegradationFallback] flag (LH-FA-CLI-005
// exception 2: stays visible under `--quiet`); the `doctor` items
// carry it (TestDoctor_SandboxFallbackFlagIsOnTheItem).
func TestSandboxRelaxationWarnings_PolicyChangesNoRenderTimeWarning(t *testing.T) {
	t.Parallel()
	for _, policy := range []string{"", "warn", "fail"} {
		got := application.SandboxRelaxationWarningsForTest("podman", policy, true, true)
		flagged, plain := relaxationFlags(got)
		if len(flagged) != 7 || len(plain) != 0 {
			t.Errorf("policy %q: flagged=%d plain=%d, want 7/0 (6 podman + 1 egress): %+v", policy, len(flagged), len(plain), got)
		}
		for _, w := range got {
			if strings.Contains(w.Message, "vfs") || strings.Contains(w.Message, "fall back") || strings.Contains(w.Message, "without daemon") {
				t.Errorf("policy %q: degradation text in render-time warning: %q", policy, w.Message)
			}
		}
	}
}

const featuresYAMLHead = "devcontainer:\n  enabled: true\n  profile: sandbox\n  features:\n"

// LH-FA-DEV-006 / LH-FA-CLI-005: an active built-in feature that adds
// container options is disclosed in the sandbox profile with one warning per
// option (the capability and the seccomp profile of `go`), each flagged with
// SecurityRelaxation so that `--quiet` keeps it. `init: true` of the feature
// is no relaxation and stays out. A `version` override keeps the feature
// active, a `source:` override makes it another feature; without the sandbox
// profile, with the feature switched off or unset nothing is disclosed.
func TestSandboxRelaxationWarnings_FeatureOptionsAreDisclosed(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		yaml    string
		sandbox bool
		want    int
	}{
		{"go in the sandbox profile", featuresYAMLHead + "    go:\n      enabled: true\n", true, 2},
		{"go with a version override", featuresYAMLHead + "    go:\n      enabled: true\n      version: \"2\"\n", true, 2},
		{"go next to node", featuresYAMLHead + "    go:\n      enabled: true\n    node:\n      enabled: true\n", true, 2},
		{"go without the sandbox profile", featuresYAMLHead + "    go:\n      enabled: true\n", false, 0},
		{"go switched off", featuresYAMLHead + "    go:\n      enabled: false\n", true, 0},
		{"go unset", featuresYAMLHead + "    go: {}\n", true, 0},
		{"go with a source override is another feature", featuresYAMLHead + "    go:\n      enabled: true\n      source: oci://registry.example.test/features/go\n", true, 0},
		{"features without options", featuresYAMLHead + "    node:\n      enabled: true\n    kubectl-helm:\n      enabled: true\n", true, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var got []driving.WarningEntry
			for _, w := range application.SandboxRelaxationWarningsForYAMLTest(t, []byte(tc.yaml), tc.sandbox) {
				if strings.Contains(w.Message, "security relaxation for the devcontainer feature") {
					got = append(got, w)
				}
			}
			if len(got) != tc.want {
				t.Fatalf("warnings = %d, want %d: %+v", len(got), tc.want, got)
			}
			var options []string
			for _, w := range got {
				if w.Code != "LH-FA-DEV-006" || w.Level != "warn" || !w.SecurityRelaxation {
					t.Errorf("warning must be LH-FA-DEV-006, warn and flagged SecurityRelaxation: %+v", w)
				}
				if !strings.Contains(w.Message, "feature go:") {
					t.Errorf("warning does not name the feature: %q", w.Message)
				}
				if strings.Contains(w.Message, "init") {
					t.Errorf("warning discloses `init: true`, which is no relaxation: %q", w.Message)
				}
				options = append(options, w.Message)
			}
			if tc.want == 2 {
				joined := strings.Join(options, "\n")
				if !strings.Contains(joined, "--cap-add=SYS_PTRACE") || !strings.Contains(joined, "--security-opt=seccomp=unconfined") {
					t.Errorf("warnings lack one of the two options:\n%s", joined)
				}
			}
		})
	}
}

// LH-FA-DEV-006: a feature from an external source gets one warning in the
// sandbox profile, per enabled entry with a `source:` override, in name
// order, flagged with SecurityRelaxation: u-boot does not read the metadata
// of such a feature and the profile promises nothing for it. No warning
// without the sandbox profile, for a switched-off or unset entry, or for a
// catalogue key without `source:`.
func TestSandboxRelaxationWarnings_ExternalFeatureEachGetsOne(t *testing.T) {
	t.Parallel()
	const two = "    zeta:\n      enabled: true\n      source: https://example.test/zeta\n" +
		"    alpha:\n      enabled: true\n      source: oci://registry.example.test/alpha\n"
	cases := []struct {
		name    string
		yaml    string
		sandbox bool
		want    []string
	}{
		{"two external features, in name order", featuresYAMLHead + two, true, []string{"alpha", "zeta"}},
		{"without the sandbox profile", featuresYAMLHead + two, false, nil},
		{"switched off", featuresYAMLHead + "    alpha:\n      enabled: false\n      source: oci://registry.example.test/alpha\n", true, nil},
		{"enabled unset", featuresYAMLHead + "    alpha:\n      source: oci://registry.example.test/alpha\n", true, nil},
		{"catalogue key without source", featuresYAMLHead + "    node:\n      enabled: true\n", true, nil},
		{"catalogue key with a source", featuresYAMLHead + "    node:\n      enabled: true\n      source: oci://registry.example.test/node\n", true, []string{"node"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var names []string
			for _, w := range application.SandboxRelaxationWarningsForYAMLTest(t, []byte(tc.yaml), tc.sandbox) {
				if !strings.Contains(w.Message, "comes from the external source") {
					continue
				}
				if w.Code != "LH-FA-DEV-006" || w.Level != "warn" || !w.SecurityRelaxation {
					t.Errorf("warning must be LH-FA-DEV-006, warn and flagged SecurityRelaxation: %+v", w)
				}
				if !strings.Contains(w.Message, "does not read its metadata") {
					t.Errorf("warning does not say that u-boot reads no metadata: %q", w.Message)
				}
				names = append(names, strings.Fields(strings.TrimPrefix(w.Message, "devcontainer feature "))[0])
			}
			if !reflect.DeepEqual(names, tc.want) {
				t.Errorf("features warned about = %v, want %v", names, tc.want)
			}
		})
	}
}
