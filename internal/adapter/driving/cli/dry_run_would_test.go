package cli_test

import (
	"errors"
	"strings"
	"testing"
)

// wouldCases pins, per subcommand, the dry-run summary (LH-FA-CLI-007)
// and the summary of an actual write.
var wouldCases = []struct {
	name   string
	dryRun string
	write  string
}{
	{"generate", "Would generate changelog (CHANGELOG.md).", "Generated changelog (CHANGELOG.md)."},
	{"remove", "Would remove service \"postgres\".\n\nWould change:\n", "Removed service \"postgres\".\n\nChanged:\n"},
	{"config set", "config: Would set project.name old → new.", "config: project.name old → new."},
	{"config set many", "config: Would set devcontainer.enabled false → true.", "config: devcontainer.enabled false → true."},
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
	for _, c := range wouldCases {
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
	for _, c := range wouldCases {
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
	for _, c := range wouldCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			sc := wouldSubcommand(t, c.name)
			out, _, err := quietRun(t, sc.app(t), append(append([]string{"--json"}, sc.cmd...), "--dry-run")...)
			if err != nil {
				t.Fatalf("--json: %v", err)
			}
			if strings.Contains(out, "Would") || !strings.Contains(out, `"dryRun":true`) {
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
