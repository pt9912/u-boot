package cli_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// LH-FA-CLI-005 (Lastenheft 0.9.0): under `--quiet` the modifying
// subcommands (init, generate, add, remove, config set) write only
// errors; their success and summary output is not written. The tests
// pair each silent success with the default-output boundary and an
// error that stays visible.

func quietRun(t *testing.T, app *cli.App, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	var o, e bytes.Buffer
	err = app.Execute(context.Background(), args, &o, &e)
	return o.String(), e.String(), err
}

func initFixture(t *testing.T) *fakeInitUseCase {
	t.Helper()
	return &fakeInitUseCase{resp: driving.InitProjectResponse{
		Project: domain.NewProject(mustProjectName(t, "myproj")),
		Created: []string{"u-boot.yaml", "compose.yaml"},
		Backups: []driving.BackupAction{{Original: "README.md", Backup: "/tmp/x/demo/README.md.bak"}},
	}}
}

func TestExecute_Init_Quiet_SuccessSummaryIsSilent(t *testing.T) {
	t.Parallel()
	out, errOut, err := quietRun(t, newApp(initFixture(t), cli.WithGetwd(sandboxGetwd)), "--quiet", "init")
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if out != "" || errOut != "" {
		t.Errorf("init --quiet wrote stdout=%q stderr=%q, want nothing", out, errOut)
	}
}

// Boundary: the same response without --quiet prints the summary.
func TestExecute_Init_NoQuiet_PrintsSummary(t *testing.T) {
	t.Parallel()
	out, _, err := quietRun(t, newApp(initFixture(t), cli.WithGetwd(sandboxGetwd)), "init")
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	for _, want := range []string{"Initialized u-boot project", "u-boot.yaml", "Backups:"} {
		if !strings.Contains(out, want) {
			t.Errorf("default init output lacks %q:\n%s", want, out)
		}
	}
}

// Negative: an error stays an error under --quiet (exit code carried).
func TestExecute_Init_Quiet_ErrorStays(t *testing.T) {
	t.Parallel()
	uc := &fakeInitUseCase{err: driving.ErrProjectExists}
	_, _, err := quietRun(t, newApp(uc, cli.WithGetwd(sandboxGetwd)), "--quiet", "init")
	if !errors.Is(err, driving.ErrProjectExists) {
		t.Fatalf("err = %v, want ErrProjectExists", err)
	}
	if got := cli.ExitCode(err); got != 10 {
		t.Errorf("ExitCode = %d, want 10", got)
	}
}

// LH-FA-CLI-005 exception 4 (LH-FA-CLI-005A, LH-FA-INIT-005): the
// `--force` path summary is the application's pre-write progress
// event. `--quiet` must not request its silencing; only `--json`
// does (the envelope on stdout must stay parseable).
func TestExecute_Init_Quiet_DoesNotSilenceForcePathSummary(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name        string
		args        []string
		wantSilence bool
	}{
		{"quiet", []string{"--quiet", "init", "--force", "--backup"}, false},
		{"default", []string{"init", "--force", "--backup"}, false},
		{"quiet json", []string{"--quiet", "--json", "init", "--force", "--backup"}, true},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			uc := initFixture(t)
			if _, _, err := quietRun(t, newApp(uc, cli.WithGetwd(sandboxGetwd)), tc.args...); err != nil {
				t.Fatalf("Execute: %v", err)
			}
			if uc.lastReq.SilenceProgress != tc.wantSilence {
				t.Errorf("SilenceProgress = %v, want %v", uc.lastReq.SilenceProgress, tc.wantSilence)
			}
			if !uc.lastReq.Force {
				t.Error("--force did not reach the request")
			}
		})
	}
}

func generateFixture() *fakeGenerateUseCase {
	return &fakeGenerateUseCase{resp: driving.GenerateResponse{
		Artifact: domain.ArtifactChangelog, Action: driving.GenerateActionCreated,
		Changed: []string{"CHANGELOG.md"},
	}}
}

func TestExecute_Generate_Quiet_SuccessSummaryIsSilent(t *testing.T) {
	t.Parallel()
	out, errOut, err := quietRun(t, newAppWithGenerate(generateFixture(), cli.WithGetwd(sandboxGetwd)), "--quiet", "generate", "changelog")
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if out != "" || errOut != "" {
		t.Errorf("generate --quiet wrote stdout=%q stderr=%q, want nothing", out, errOut)
	}
}

func TestExecute_Generate_NoQuiet_PrintsSummary(t *testing.T) {
	t.Parallel()
	out, _, err := quietRun(t, newAppWithGenerate(generateFixture(), cli.WithGetwd(sandboxGetwd)), "generate", "changelog")
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !strings.Contains(out, "Generated changelog (CHANGELOG.md).") {
		t.Errorf("default generate output lacks the summary:\n%s", out)
	}
}

// Negative: the NoOp summary line is a status message too.
func TestExecute_Generate_Quiet_NoOpIsSilent(t *testing.T) {
	t.Parallel()
	uc := &fakeGenerateUseCase{resp: driving.GenerateResponse{Artifact: domain.ArtifactChangelog, Action: driving.GenerateActionNoOp}}
	out, _, err := quietRun(t, newAppWithGenerate(uc, cli.WithGetwd(sandboxGetwd)), "--quiet", "generate", "changelog")
	if err != nil || out != "" {
		t.Errorf("err=%v stdout=%q, want silent success", err, out)
	}
}

func TestExecute_Generate_Quiet_ErrorStays(t *testing.T) {
	t.Parallel()
	uc := &fakeGenerateUseCase{err: driving.ErrProjectNotInitialized}
	_, _, err := quietRun(t, newAppWithGenerate(uc, cli.WithGetwd(sandboxGetwd)), "--quiet", "generate", "changelog")
	if !errors.Is(err, driving.ErrProjectNotInitialized) {
		t.Fatalf("err = %v, want ErrProjectNotInitialized", err)
	}
}

func addFixture(t *testing.T, prior, state domain.ServiceState, changed []string) *fakeAddServiceUseCase {
	t.Helper()
	return &fakeAddServiceUseCase{resp: driving.AddServiceResponse{
		ServiceName: mustServiceNameCLI(t, "postgres"), PriorState: prior, State: state, Changed: changed,
	}}
}

// add: the three summary shapes (added, repaired, already active) are
// all silent under --quiet.
func TestExecute_Add_Quiet_AllSummaryShapesAreSilent(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name         string
		prior, state domain.ServiceState
		changed      []string
	}{
		{"added", domain.ServiceStateUnregistered, domain.ServiceStateActive, []string{"compose.yaml"}},
		{"repaired", domain.ServiceStateActive, domain.ServiceStateActive, []string{"compose.yaml"}},
		{"already-active", domain.ServiceStateActive, domain.ServiceStateActive, nil},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			uc := addFixture(t, tc.prior, tc.state, tc.changed)
			out, errOut, err := quietRun(t, newAppWithAdd(uc, cli.WithGetwd(sandboxGetwd)), "--quiet", "add", "postgres")
			if err != nil {
				t.Fatalf("Execute: %v", err)
			}
			if !uc.called {
				t.Fatal("add use case was not called")
			}
			if out != "" || errOut != "" {
				t.Errorf("add --quiet wrote stdout=%q stderr=%q, want nothing", out, errOut)
			}
		})
	}
}

func TestExecute_Add_NoQuiet_PrintsSummary(t *testing.T) {
	t.Parallel()
	uc := addFixture(t, domain.ServiceStateUnregistered, domain.ServiceStateActive, []string{"compose.yaml"})
	out, _, err := quietRun(t, newAppWithAdd(uc, cli.WithGetwd(sandboxGetwd)), "add", "postgres")
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !strings.Contains(out, `Added service "postgres".`) {
		t.Errorf("default add output lacks the summary:\n%s", out)
	}
}

func TestExecute_Add_Quiet_ErrorStays(t *testing.T) {
	t.Parallel()
	uc := &fakeAddServiceUseCase{err: driving.ErrServiceUnsupported}
	_, _, err := quietRun(t, newAppWithAdd(uc, cli.WithGetwd(sandboxGetwd)), "--quiet", "add", "postgres")
	if !errors.Is(err, driving.ErrServiceUnsupported) {
		t.Fatalf("err = %v, want ErrServiceUnsupported", err)
	}
	if got := cli.ExitCode(err); got != 10 {
		t.Errorf("ExitCode = %d, want 10", got)
	}
}

func removeFixture(t *testing.T) *fakeRemoveServiceUseCase {
	t.Helper()
	return &fakeRemoveServiceUseCase{resp: driving.RemoveServiceResponse{
		ServiceName: mustServiceNameCLI(t, "postgres"),
		PriorState:  domain.ServiceStateActive, State: domain.ServiceStateDeactivated,
		Changed: []string{"compose.yaml", "u-boot.yaml"},
	}}
}

func purgeWarning() driving.WarningEntry {
	return driving.WarningEntry{Code: "LH-FA-ADD-007", Level: "warn", Subject: "postgres",
		Message: "--purge: volume \"pgdata\" could not be removed (in use)"}
}

// `remove` takes --quiet like every other subcommand (the root
// persistent flag); success summary, removed-volumes line and the
// purge warnings are not written.
func TestExecute_Remove_Quiet_SummaryAndVolumesLineAreSilent(t *testing.T) {
	t.Parallel()
	uc := removeFixture(t)
	uc.resp.PurgeAttempted = true
	uc.resp.PurgedVolumes = []string{"pgdata"}
	out, errOut, err := quietRun(t, newAppWithRemove(uc, cli.WithGetwd(sandboxGetwd)), "--quiet", "--yes", "remove", "postgres", "--purge")
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !uc.called || !uc.lastReq.Purge {
		t.Fatalf("remove use case called=%v purge=%v", uc.called, uc.lastReq.Purge)
	}
	if out != "" || errOut != "" {
		t.Errorf("remove --quiet wrote stdout=%q stderr=%q, want nothing", out, errOut)
	}
}

// LH-FA-ADD-007 / LH-FA-CLI-005 (explicit setting, pinned): the purge
// warnings are neither an error nor one of the four exceptions, so
// they fall under "all other notices are dropped" — the real-removal
// warning of the purge path and the deferred-volumes warning of the
// non-automated path both vanish under --quiet.
func TestExecute_Remove_Quiet_PurgeWarningsAreDropped(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		prep func(r *driving.RemoveServiceResponse)
	}{
		{"real purge: per-volume warning", func(r *driving.RemoveServiceResponse) {
			r.PurgeAttempted = true
			r.Warnings = []driving.WarningEntry{purgeWarning()}
		}},
		{"deferred purge: not-automated warning", func(r *driving.RemoveServiceResponse) {
			r.VolumesPurged = false
			r.Warnings = []driving.WarningEntry{purgeWarning()}
		}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			uc := removeFixture(t)
			tc.prep(&uc.resp)
			// Boundary first: without --quiet the warning is written.
			_, defErr, err := quietRun(t, newAppWithRemove(uc, cli.WithGetwd(sandboxGetwd)), "--yes", "remove", "postgres", "--purge")
			if err != nil || !strings.Contains(defErr, "WARNING:") {
				t.Fatalf("default run: err=%v stderr=%q, want a WARNING", err, defErr)
			}
			out, errOut, err := quietRun(t, newAppWithRemove(uc, cli.WithGetwd(sandboxGetwd)), "--quiet", "--yes", "remove", "postgres", "--purge")
			if err != nil {
				t.Fatalf("Execute: %v", err)
			}
			if out != "" || strings.Contains(errOut, "WARNING") || errOut != "" {
				t.Errorf("--quiet kept a purge warning: stdout=%q stderr=%q", out, errOut)
			}
		})
	}
}

func TestExecute_Remove_Quiet_ErrorStays(t *testing.T) {
	t.Parallel()
	uc := &fakeRemoveServiceUseCase{err: driving.ErrServiceUnregistered}
	_, _, err := quietRun(t, newAppWithRemove(uc, cli.WithGetwd(sandboxGetwd)), "--quiet", "remove", "postgres")
	if !errors.Is(err, driving.ErrServiceUnregistered) {
		t.Fatalf("err = %v, want ErrServiceUnregistered", err)
	}
	if got := cli.ExitCode(err); got != 10 {
		t.Errorf("ExitCode = %d, want 10", got)
	}
}

// Negative: a purge that the gate refuses stays an error under --quiet.
func TestExecute_Remove_Quiet_PurgeRefusedStaysAnError(t *testing.T) {
	t.Parallel()
	uc := &fakeRemoveServiceUseCase{err: driving.ErrConfirmationRequired}
	_, _, err := quietRun(t, newAppWithRemove(uc, cli.WithGetwd(sandboxGetwd)), "--quiet", "--no-interactive", "remove", "postgres", "--purge")
	if !errors.Is(err, driving.ErrConfirmationRequired) || cli.ExitCode(err) != 10 {
		t.Errorf("err = %v (exit %d), want ErrConfirmationRequired / 10", err, cli.ExitCode(err))
	}
}

// config set changes the project: its summary line is a success
// message and is not written under --quiet (single and multi form);
// the default output keeps it.
func TestExecute_ConfigSet_Quiet_SummaryIsSilent(t *testing.T) {
	t.Parallel()
	uc := func() *fakeConfigUseCase {
		return &fakeConfigUseCase{
			setResp: driving.ConfigSetResponse{Path: mustCfgPath(t, "project.name"), OldValue: "old", NewValue: "new"},
			setManyResp: driving.ConfigSetManyResponse{Entries: []driving.ConfigSetEntry{
				{Path: mustCfgPath(t, "project.name"), OldValue: "old", NewValue: "new"},
				{Path: mustCfgPath(t, "devcontainer.enabled"), OldValue: "false", NewValue: "true"},
			}},
		}
	}
	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{"single", []string{"config", "set", "project.name", "new"}, "config: project.name old → new."},
		{"many", []string{"config", "set", "project.name", "new", "devcontainer.enabled", "true"}, "config: devcontainer.enabled false → true."},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			out, _, err := quietRun(t, newAppWithConfigStub(uc()), tc.args...)
			if err != nil || !strings.Contains(out, tc.want) {
				t.Fatalf("default: err=%v stdout=%q, want %q", err, out, tc.want)
			}
			qOut, qErr, err := quietRun(t, newAppWithConfigStub(uc()), append([]string{"--quiet"}, tc.args...)...)
			if err != nil {
				t.Fatalf("--quiet: %v", err)
			}
			if qOut != "" || qErr != "" {
				t.Errorf("config set --quiet wrote stdout=%q stderr=%q, want nothing", qOut, qErr)
			}
		})
	}
}

func TestExecute_ConfigSet_Quiet_ErrorStays(t *testing.T) {
	t.Parallel()
	uc := &fakeConfigUseCase{setErr: driving.ErrConfigPathUnknown}
	_, _, err := quietRun(t, newAppWithConfigStub(uc), "--quiet", "config", "set", "project.name", "x")
	if !errors.Is(err, driving.ErrConfigPathUnknown) {
		t.Fatalf("err = %v, want ErrConfigPathUnknown", err)
	}
}

// LH-FA-CLI-005 (negative test against over-suppression): the query
// commands keep their output under --quiet — it is the result of the
// command.
func TestExecute_Quiet_QueryCommandsKeepTheirOutput(t *testing.T) {
	t.Parallel()
	cfg := func() *fakeConfigUseCase {
		return &fakeConfigUseCase{
			getResp:  driving.ConfigGetResponse{Path: mustCfgPath(t, "project.name"), Value: "demo"},
			listResp: driving.ConfigListResponse{Entries: []driving.ConfigPathValue{entry(t, "project.name", "demo")}},
		}
	}
	cases := []struct {
		name string
		run  func(args []string) (string, error)
		args []string
		want string
	}{
		{"config get", func(a []string) (string, error) {
			o, _, err := quietRun(t, newAppWithConfigStub(cfg()), a...)
			return o, err
		}, []string{"config", "get", "project.name"}, "demo\n"},
		{"config list", func(a []string) (string, error) {
			o, _, err := quietRun(t, newAppWithConfigStub(cfg()), a...)
			return o, err
		}, []string{"config", "list"}, "project.name=demo\n"},
		{"template list", func(a []string) (string, error) {
			uc := &fakeTemplateListUseCase{resp: driving.TemplateListResponse{Templates: []domain.TemplateMetadata{
				{Name: "basic", Description: "Minimal skeleton.", Version: "0.1.0"}}}}
			o, _, err := quietRun(t, newAppWithTemplateList(uc), a...)
			return o, err
		}, []string{"template", "list"}, "basic"},
		{"logs", func(a []string) (string, error) {
			uc := &fakeLogsUseCase{onLogs: func(req driving.LogsRequest) error {
				_, err := req.OutputSink.Write([]byte("postgres  | ready\n"))
				return err
			}}
			o, _, err := quietRun(t, newAppWithLogsStub(uc), a...)
			return o, err
		}, []string{"logs"}, "postgres  | ready\n"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			plain, err := tc.run(tc.args)
			if err != nil {
				t.Fatalf("default: %v", err)
			}
			quiet, err := tc.run(append([]string{"--quiet"}, tc.args...))
			if err != nil {
				t.Fatalf("--quiet: %v", err)
			}
			if !strings.Contains(quiet, tc.want) {
				t.Errorf("--quiet output %q lacks %q", quiet, tc.want)
			}
			if quiet != plain {
				t.Errorf("--quiet output differs from the default:\n quiet=%q\n plain=%q", quiet, plain)
			}
		})
	}
}

// LH-FA-CLI-005 / LH-NFA-USE-004: `--quiet --json` is unchanged for
// add, remove and config set — the stdout is byte-identical to
// `--json` and every warning stays in diagnostics[] (no new wire
// field; the exception flags are not serialised).
func TestExecute_QuietJSON_AddRemoveUnchanged(t *testing.T) {
	t.Parallel()
	fb := quietFallbacks()["vfs"]
	cases := []struct {
		name string
		run  func(args []string) (string, error)
		args []string
		warn int
	}{
		{"add", func(a []string) (string, error) {
			uc := addFixture(t, domain.ServiceStateUnregistered, domain.ServiceStateActive, []string{"compose.yaml"})
			o, _, err := quietRun(t, newAppWithAdd(uc, cli.WithGetwd(sandboxGetwd)), a...)
			return o, err
		}, []string{"add", "postgres", "--json"}, 0},
		{"remove", func(a []string) (string, error) {
			uc := removeFixture(t)
			uc.resp.Warnings = []driving.WarningEntry{purgeWarning(), fb}
			o, _, err := quietRun(t, newAppWithRemove(uc, cli.WithGetwd(sandboxGetwd)), a...)
			return o, err
		}, []string{"remove", "postgres", "--json"}, 2},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			plain, err := tc.run(tc.args)
			if err != nil {
				t.Fatalf("--json: %v", err)
			}
			quiet, err := tc.run(append([]string{"--quiet"}, tc.args...))
			if err != nil {
				t.Fatalf("--quiet --json: %v", err)
			}
			if quiet != plain {
				t.Errorf("--quiet --json differs from --json:\n quiet=%s\n plain=%s", quiet, plain)
			}
			if n := strings.Count(quiet, `"level":"warn"`); n != tc.warn {
				t.Errorf("warn diagnostics = %d, want %d:\n%s", n, tc.warn, quiet)
			}
			if low := strings.ToLower(quiet); strings.Contains(low, "securityrelaxation") || strings.Contains(low, "degradationfallback") {
				t.Errorf("exception flag leaked into the wire form:\n%s", quiet)
			}
		})
	}
}

// The fallback flag stays off the wire for init, generate and up too
// (the existing count test uses relaxation-only fixtures).
func TestExecute_QuietJSON_FallbackFlagNotInWireForm(t *testing.T) {
	t.Parallel()
	ws := []driving.WarningEntry{quietFallbacks()["vfs"]}
	var stdout, stderr bytes.Buffer
	uc := &fakeInitUseCase{resp: driving.InitProjectResponse{Warnings: ws}}
	if err := newApp(uc, cli.WithGetwd(sandboxGetwd)).Execute(context.Background(),
		[]string{"--quiet", "--json", "init"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	out := stdout.String()
	if !strings.Contains(out, "falls back to vfs storage") {
		t.Errorf("fallback warning missing from diagnostics[]:\n%s", out)
	}
	if strings.Contains(strings.ToLower(out), "degradationfallback") {
		t.Errorf("DegradationFallback leaked into the wire form:\n%s", out)
	}
}
