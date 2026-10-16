package cli_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// LH-FA-CLI-005 (Lastenheft 0.10.0), LH-FA-CLI-007, LH-FA-CLI-008:
// under --quiet the output of --dry-run and --diff is the result of
// the call and stays; the success message that closes an actual write
// does not, also not with --diff alone. The matrix below is the rule
// table of slice-v1-cli-quiet-vorschau §1, run for every subcommand
// that takes the preview options.

const previewDiffMarker = "--- compose.yaml (modify)"

func previewPlanned() []driving.PlannedFile {
	return []driving.PlannedFile{{
		Path: "compose.yaml", Action: "modify",
		OldContent: []byte("a: 1\n"), NewContent: []byte("a: 2\n"),
	}}
}

// previewSubcommand builds a fresh app (with a fake use case that
// reports a planned change) for one subcommand with preview options.
type previewSubcommand struct {
	name    string
	cmd     []string // subcommand and its operands, without options
	summary string   // summary text that the dry-run plan carries
	app     func(t *testing.T) *cli.App
	failApp func(t *testing.T) (*cli.App, error)
}

func previewSubcommands() []previewSubcommand {
	return []previewSubcommand{
		{
			name: "init", cmd: []string{"init"}, summary: "Would initialize u-boot project",
			app: func(t *testing.T) *cli.App {
				uc := initFixture(t)
				uc.resp.PlannedFiles = previewPlanned()
				return newApp(uc, cli.WithGetwd(sandboxGetwd))
			},
			failApp: func(t *testing.T) (*cli.App, error) {
				return newApp(&fakeInitUseCase{err: driving.ErrProjectExists}, cli.WithGetwd(sandboxGetwd)), driving.ErrProjectExists
			},
		},
		{
			name: "generate", cmd: []string{"generate", "changelog"}, summary: "Generated changelog (CHANGELOG.md).",
			app: func(t *testing.T) *cli.App {
				uc := generateFixture()
				uc.resp.PlannedFiles = previewPlanned()
				return newAppWithGenerate(uc, cli.WithGetwd(sandboxGetwd))
			},
			failApp: func(t *testing.T) (*cli.App, error) {
				return newAppWithGenerate(&fakeGenerateUseCase{err: driving.ErrProjectNotInitialized}, cli.WithGetwd(sandboxGetwd)), driving.ErrProjectNotInitialized
			},
		},
		{
			name: "add", cmd: []string{"add", "postgres"}, summary: `Would add service "postgres".`,
			app: func(t *testing.T) *cli.App {
				uc := addFixture(t, domain.ServiceStateUnregistered, domain.ServiceStateActive, []string{"compose.yaml"})
				uc.resp.PlannedFiles = previewPlanned()
				return newAppWithAdd(uc, cli.WithGetwd(sandboxGetwd))
			},
			failApp: func(t *testing.T) (*cli.App, error) {
				return newAppWithAdd(&fakeAddServiceUseCase{err: driving.ErrServiceUnsupported}, cli.WithGetwd(sandboxGetwd)), driving.ErrServiceUnsupported
			},
		},
		{
			name: "remove", cmd: []string{"remove", "postgres"}, summary: `Removed service "postgres".`,
			app: func(t *testing.T) *cli.App {
				uc := removeFixture(t)
				uc.resp.PlannedFiles = previewPlanned()
				return newAppWithRemove(uc, cli.WithGetwd(sandboxGetwd))
			},
			failApp: func(t *testing.T) (*cli.App, error) {
				return newAppWithRemove(&fakeRemoveServiceUseCase{err: driving.ErrServiceUnregistered}, cli.WithGetwd(sandboxGetwd)), driving.ErrServiceUnregistered
			},
		},
		{
			name: "config set", cmd: []string{"config", "set", "project.name", "new"}, summary: "config: project.name old → new.",
			app: func(t *testing.T) *cli.App {
				return newAppWithConfigStub(&fakeConfigUseCase{
					setResp: driving.ConfigSetResponse{
						Path: mustCfgPath(t, "project.name"), OldValue: "old", NewValue: "new",
						PlannedFiles: previewPlanned(),
					},
				})
			},
			failApp: func(t *testing.T) (*cli.App, error) {
				return newAppWithConfigStub(&fakeConfigUseCase{setErr: driving.ErrConfigPathUnknown}), driving.ErrConfigPathUnknown
			},
		},
		{
			name: "config set many", cmd: []string{"config", "set", "project.name", "new", "devcontainer.enabled", "true"},
			summary: "config: devcontainer.enabled false → true.",
			app: func(t *testing.T) *cli.App {
				return newAppWithConfigStub(&fakeConfigUseCase{
					setManyResp: driving.ConfigSetManyResponse{
						Entries: []driving.ConfigSetEntry{
							{Path: mustCfgPath(t, "project.name"), OldValue: "old", NewValue: "new"},
							{Path: mustCfgPath(t, "devcontainer.enabled"), OldValue: "false", NewValue: "true"},
						},
						PlannedFiles: previewPlanned(),
					},
				})
			},
			failApp: func(t *testing.T) (*cli.App, error) {
				return newAppWithConfigStub(&fakeConfigUseCase{setManyErr: driving.ErrConfigPathUnknown}), driving.ErrConfigPathUnknown
			},
		},
	}
}

// The rule table: per option combination, is the plan summary and the
// diff written under --quiet. The last row is the unchanged baseline.
func TestExecute_Quiet_PreviewMatrix(t *testing.T) {
	t.Parallel()
	rows := []struct {
		name        string
		opts        []string
		wantSummary bool
		wantDiff    bool
	}{
		{"dry-run", []string{"--dry-run"}, true, false},
		{"dry-run diff", []string{"--dry-run", "--diff"}, true, true},
		{"diff without dry-run", []string{"--diff"}, false, true},
		{"no preview option", nil, false, false},
	}
	for _, sc := range previewSubcommands() {
		for _, row := range rows {
			sc, row := sc, row
			t.Run(sc.name+"/"+row.name, func(t *testing.T) {
				t.Parallel()
				args := func(quiet bool) []string {
					var a []string
					if quiet {
						a = append(a, "--quiet")
					}
					a = append(a, sc.cmd...)
					return append(a, row.opts...)
				}
				out, errOut, err := quietRun(t, sc.app(t), args(true)...)
				if err != nil {
					t.Fatalf("--quiet: %v", err)
				}
				if got := strings.Contains(out, sc.summary); got != row.wantSummary {
					t.Errorf("summary %q present=%v, want %v:\n%s", sc.summary, got, row.wantSummary, out)
				}
				if got := strings.Contains(out, previewDiffMarker); got != row.wantDiff {
					t.Errorf("diff present=%v, want %v:\n%s", got, row.wantDiff, out)
				}
				if errOut != "" {
					t.Errorf("stderr = %q, want empty", errOut)
				}
				// A dry-run plan is the same text with and without --quiet.
				if len(row.opts) > 0 && row.opts[0] == "--dry-run" {
					plain, _, err := quietRun(t, sc.app(t), args(false)...)
					if err != nil || plain != out {
						t.Errorf("--quiet plan differs from the default output (err=%v):\n quiet=%q\n plain=%q", err, out, plain)
					}
				}
			})
		}
	}
}

// Negative against over-correction: after an actual write the success
// message stays dropped under --quiet, with and without --diff; only
// the diff of --diff may appear.
func TestExecute_Quiet_ActualWriteStaysSilent(t *testing.T) {
	t.Parallel()
	for _, sc := range previewSubcommands() {
		sc := sc
		t.Run(sc.name, func(t *testing.T) {
			t.Parallel()
			// Boundary: the default run writes the success message.
			plain, _, err := quietRun(t, sc.app(t), sc.cmd...)
			if err != nil || plain == "" {
				t.Fatalf("default: err=%v stdout=%q, want a success message", err, plain)
			}
			successLine := strings.SplitN(plain, "\n", 2)[0]
			out, _, err := quietRun(t, sc.app(t), append([]string{"--quiet"}, sc.cmd...)...)
			if err != nil || out != "" {
				t.Errorf("--quiet: err=%v stdout=%q, want silence", err, out)
			}
			diffOnly, _, err := quietRun(t, sc.app(t), append(append([]string{"--quiet"}, sc.cmd...), "--diff")...)
			if err != nil {
				t.Fatalf("--quiet --diff: %v", err)
			}
			if !strings.Contains(diffOnly, previewDiffMarker) {
				t.Errorf("--quiet --diff lacks the diff:\n%s", diffOnly)
			}
			if strings.Contains(diffOnly, successLine) {
				t.Errorf("--quiet --diff wrote the success message %q:\n%s", successLine, diffOnly)
			}
		})
	}
}

// Negative: an error stays visible under --quiet with every preview
// option (the exit code and error text are carried by the error).
func TestExecute_Quiet_PreviewErrorStays(t *testing.T) {
	t.Parallel()
	for _, sc := range previewSubcommands() {
		for _, opts := range [][]string{{"--dry-run"}, {"--diff"}, {"--dry-run", "--diff"}} {
			sc, opts := sc, opts
			t.Run(sc.name+"/"+strings.Join(opts, " "), func(t *testing.T) {
				t.Parallel()
				app, want := sc.failApp(t)
				args := append(append([]string{"--quiet"}, sc.cmd...), opts...)
				_, _, err := quietRun(t, app, args...)
				if !errors.Is(err, want) {
					t.Fatalf("err = %v, want %v", err, want)
				}
			})
		}
	}
}

// Negative: --quiet --json is byte-identical to --json for the
// preview options (the envelope carries the plan either way).
func TestExecute_Quiet_PreviewJSONUnchanged(t *testing.T) {
	t.Parallel()
	for _, sc := range previewSubcommands() {
		sc := sc
		t.Run(sc.name, func(t *testing.T) {
			t.Parallel()
			args := func(quiet bool) []string {
				var a []string
				if quiet {
					a = append(a, "--quiet")
				}
				a = append(a, "--json")
				a = append(a, sc.cmd...)
				return append(a, "--dry-run", "--diff")
			}
			plain, _, err := quietRun(t, sc.app(t), args(false)...)
			if err != nil {
				t.Fatalf("--json: %v", err)
			}
			quiet, _, err := quietRun(t, sc.app(t), args(true)...)
			if err != nil || quiet != plain || plain == "" {
				t.Errorf("--quiet --json differs (err=%v):\n quiet=%s\n plain=%s", err, quiet, plain)
			}
		})
	}
}

// Boundary: a no-op plan keeps its "no changes" line under
// `--dry-run --quiet` (add: already active).
func TestExecute_Add_Quiet_DryRunNoOpPlanStays(t *testing.T) {
	t.Parallel()
	uc := addFixture(t, domain.ServiceStateActive, domain.ServiceStateActive, nil)
	out, _, err := quietRun(t, newAppWithAdd(uc, cli.WithGetwd(sandboxGetwd)), "--quiet", "add", "postgres", "--dry-run")
	if err != nil || !strings.Contains(out, `Service "postgres" is already active; no changes.`) {
		t.Errorf("err=%v stdout=%q, want the no-op plan line", err, out)
	}
}

// LH-FA-CLI-005: the dry-run plan of init is the result, but its
// notices stay dropped except the exception warnings.
func TestExecute_Init_Quiet_DryRunPlanKeepsOnlyExceptionWarnings(t *testing.T) {
	t.Parallel()
	uc := initFixture(t)
	uc.resp.Warnings = quietAllWarnings()
	out, _, err := quietRun(t, newApp(uc, cli.WithGetwd(sandboxGetwd)), "--quiet", "init", "--dry-run")
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !strings.Contains(out, "Would initialize u-boot project") {
		t.Errorf("plan missing:\n%s", out)
	}
	assertOnlyRelaxationVisible(t, out)
}
