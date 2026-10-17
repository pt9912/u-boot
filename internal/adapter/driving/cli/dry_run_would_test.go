package cli_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

type wouldCase struct {
	name   string
	dryRun string
	write  string
}

// wouldCases pins, per subcommand, the dry-run summary (LH-FA-CLI-007)
// and the summary of an actual write.
func wouldCases() []wouldCase {
	return []wouldCase{
		{"generate", "Would generate changelog (CHANGELOG.md).", "Generated changelog (CHANGELOG.md)."},
		{"remove", "Would remove service \"postgres\".\n\nWould change:\n", "Removed service \"postgres\".\n\nChanged:\n"},
		{"config set", "config: Would set project.name old → new.", "config: project.name old → new."},
		{"config set many", "config: Would set devcontainer.enabled false → true.", "config: devcontainer.enabled false → true."},
	}
}

func wouldSubcommand(t *testing.T, name string) previewSubcommand {
	t.Helper()
	for _, sc := range previewSubcommands() {
		if sc.name == name {
			return sc
		}
	}
	t.Fatalf("no preview subcommand %q", name)
	return previewSubcommand{}
}

// Happy: --dry-run (with and without --diff, with and without --quiet)
// says what would happen, and not what happened.
func TestDryRun_SummarySaysWould(t *testing.T) {
	t.Parallel()
	for _, c := range wouldCases() {
		for _, opts := range [][]string{{"--dry-run"}, {"--dry-run", "--diff"}} {
			for _, quiet := range []bool{false, true} {
				c, opts, quiet := c, opts, quiet
				t.Run(c.name+"/"+strings.Join(opts, " ")+map[bool]string{true: "/quiet"}[quiet], func(t *testing.T) {
					t.Parallel()
					sc := wouldSubcommand(t, c.name)
					args := append([]string{}, sc.cmd...)
					if quiet {
						args = append([]string{"--quiet"}, args...)
					}
					out, _, err := quietRun(t, sc.app(t), append(args, opts...)...)
					if err != nil {
						t.Fatalf("Execute: %v", err)
					}
					if !strings.Contains(out, c.dryRun) {
						t.Errorf("stdout lacks %q:\n%s", c.dryRun, out)
					}
					if strings.Contains(out, c.write) {
						t.Errorf("stdout carries the past form %q:\n%s", c.write, out)
					}
				})
			}
		}
	}
}

// Boundary: without --dry-run (no option, or --diff alone) the past
// form stays and "Would" does not appear.
func TestDryRun_ActualWriteKeepsPastForm(t *testing.T) {
	t.Parallel()
	for _, c := range wouldCases() {
		for _, opts := range [][]string{nil, {"--diff"}} {
			c, opts := c, opts
			t.Run(c.name+"/"+strings.Join(opts, " "), func(t *testing.T) {
				t.Parallel()
				sc := wouldSubcommand(t, c.name)
				out, _, err := quietRun(t, sc.app(t), append(append([]string{}, sc.cmd...), opts...)...)
				if err != nil {
					t.Fatalf("Execute: %v", err)
				}
				if !strings.Contains(out, c.write) {
					t.Errorf("stdout lacks %q:\n%s", c.write, out)
				}
				if strings.Contains(out, "Would") {
					t.Errorf("stdout says Would without --dry-run:\n%s", out)
				}
			})
		}
	}
}

// Negative: --dry-run --json carries no Would text (the envelope is
// the contract), and an error stays an error with --dry-run.
func TestDryRun_JSONAndErrorsUnchanged(t *testing.T) {
	t.Parallel()
	for _, c := range wouldCases() {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			sc := wouldSubcommand(t, c.name)
			out, _, err := quietRun(t, sc.app(t), append(append([]string{"--json"}, sc.cmd...), "--dry-run")...)
			if err != nil {
				t.Fatalf("--json: %v", err)
			}
			if strings.Contains(out, strings.SplitN(c.dryRun, "\n", 2)[0]) || !strings.Contains(out, `"dryRun":true`) {
				t.Errorf("--dry-run --json output changed:\n%s", out)
			}
			app, want := sc.failApp(t)
			_, _, err = quietRun(t, app, append(append([]string{}, sc.cmd...), "--dry-run")...)
			if !errors.Is(err, want) {
				t.Errorf("err = %v, want %v", err, want)
			}
		})
	}
}

// Happy/Boundary per generate action: each changing action reads
// "Would ..." under --dry-run and keeps its past form otherwise; NoOp
// reads the same in both modes.
func TestDryRun_GenerateActions(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name   string
		action driving.GenerateAction
		dryRun string
		write  string
	}{
		{"created", driving.GenerateActionCreated, "Would generate changelog (CHANGELOG.md).", "Generated changelog (CHANGELOG.md)."},
		{"updated-block", driving.GenerateActionUpdatedBlock, "Would update changelog managed block (CHANGELOG.md).", "Updated changelog managed block (CHANGELOG.md)."},
		{"repaired-manual", driving.GenerateActionRepairedManual, "Would repair changelog structure (CHANGELOG.md).", "Repaired changelog structure (CHANGELOG.md)."},
		{"no-op", driving.GenerateActionNoOp, "changelog already up to date; no changes.", "changelog already up to date; no changes."},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			app := func() *cli.App {
				uc := generateFixture()
				uc.resp.Action = c.action
				return newAppWithGenerate(uc, cli.WithGetwd(sandboxGetwd))
			}
			dry, _, err := quietRun(t, app(), "generate", "changelog", "--dry-run")
			if err != nil || !strings.Contains(dry, c.dryRun) {
				t.Errorf("--dry-run: err=%v stdout=%q, want %q", err, dry, c.dryRun)
			}
			write, _, err := quietRun(t, app(), "generate", "changelog")
			if err != nil || !strings.Contains(write, c.write) || strings.Contains(write, "Would") {
				t.Errorf("write: err=%v stdout=%q, want %q", err, write, c.write)
			}
		})
	}
}

// Drift guard: every GenerateAction value that the port names has a
// summary shape. A changing action reads "Would ..." under --dry-run;
// a value that reaches the defensive default branch (no dry-run form)
// fails here.
func TestDryRun_GenerateEveryActionHasADryRunForm(t *testing.T) {
	t.Parallel()
	for a := driving.GenerateAction(0); a.String() != "unknown"; a++ {
		a := a
		t.Run(a.String(), func(t *testing.T) {
			t.Parallel()
			app := func() *cli.App {
				uc := generateFixture()
				uc.resp.Action = a
				return newAppWithGenerate(uc, cli.WithGetwd(sandboxGetwd))
			}
			dry, _, err := quietRun(t, app(), "generate", "changelog", "--dry-run")
			if err != nil {
				t.Fatalf("--dry-run: %v", err)
			}
			write, _, err := quietRun(t, app(), "generate", "changelog")
			if err != nil {
				t.Fatalf("write: %v", err)
			}
			if a == driving.GenerateActionNoOp {
				if dry != write || !strings.Contains(dry, "no changes") {
					t.Errorf("no-op differs: dry=%q write=%q", dry, write)
				}
				return
			}
			if !strings.HasPrefix(dry, "Would ") || strings.Contains(write, "Would") {
				t.Errorf("action %s: dry=%q write=%q, want a Would form under --dry-run only", a, dry, write)
			}
		})
	}
}

// Boundary: `config set` with several pairs, one unchanged and one
// changed, writes both lines under --dry-run: the no-op line as in a
// write, the changed one as "Would set".
func TestDryRun_ConfigSetManyMixedPairs(t *testing.T) {
	t.Parallel()
	app := func(t *testing.T) *cli.App {
		return newAppWithConfigStub(&fakeConfigUseCase{
			setManyResp: driving.ConfigSetManyResponse{
				Entries: []driving.ConfigSetEntry{
					{Path: mustCfgPath(t, "project.name"), OldValue: "same", NewValue: "same"},
					{Path: mustCfgPath(t, "devcontainer.enabled"), OldValue: "false", NewValue: "true"},
				},
				PlannedFiles: previewPlanned(),
			},
		})
	}
	args := []string{"config", "set", "project.name", "same", "devcontainer.enabled", "true"}
	dry, _, err := quietRun(t, app(t), append(args, "--dry-run")...)
	if err != nil {
		t.Fatalf("--dry-run: %v", err)
	}
	for _, want := range []string{
		"config: project.name already same; no changes.\n",
		"config: Would set devcontainer.enabled false → true.\n",
	} {
		if !strings.Contains(dry, want) {
			t.Errorf("--dry-run lacks %q:\n%s", want, dry)
		}
	}
	write, _, err := quietRun(t, app(t), args...)
	if err != nil {
		t.Fatalf("write: %v", err)
	}
	if !strings.Contains(write, "config: devcontainer.enabled false → true.\n") || strings.Contains(write, "Would") {
		t.Errorf("write output wrong:\n%s", write)
	}
}
