package application_test

import (
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

// LH-FA-CLI-005 / LH-FA-DEV-006/-007/-008: exactly the "security
// relaxation" notices carry SecurityRelaxation (they stay visible
// under --quiet); the flag sits on the entry, so the "has no effect"
// notices that share the code LH-FA-DEV-007 stay unflagged.
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
	}{
		{"podman-sandbox: five relaxations", podman, true, 5, 0},
		{"docker-sandbox: --privileged", docker, true, 1, 0},
		{"egress-sandbox: NET_ADMIN", egress, true, 1, 0},
		{"podman-without-sandbox: has no effect", podman, false, 0, 1},
		{"docker-without-sandbox: has no effect", docker, false, 0, 1},
		{"egress-without-sandbox: has no effect", egress, false, 0, 1},
		{"no nested runtime: nothing", sandboxCfg{}, true, 0, 0},
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
				if !strings.Contains(w.Message, "security relaxation") {
					t.Errorf("flagged entry is not a relaxation notice: %q", w.Message)
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
		if len(flagged) != 6 || len(plain) != 0 {
			t.Errorf("policy %q: flagged=%d plain=%d, want 6/0 (5 podman + 1 egress): %+v", policy, len(flagged), len(plain), got)
		}
		for _, w := range got {
			if strings.Contains(w.Message, "vfs") || strings.Contains(w.Message, "fall back") || strings.Contains(w.Message, "without daemon") {
				t.Errorf("policy %q: degradation text in render-time warning: %q", policy, w.Message)
			}
		}
	}
}
