package cli_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// fakeInitUseCase records the last InitProjectRequest and returns
// the configured response/error. Plus a fake getwd hook lives in
// cli/init.go and is exposed through the test-only export in
// export_test.go.
type fakeInitUseCase struct {
	called  bool
	lastReq driving.InitProjectRequest
	resp    driving.InitProjectResponse
	err     error
}

func (f *fakeInitUseCase) Init(_ context.Context, req driving.InitProjectRequest) (driving.InitProjectResponse, error) {
	f.called = true
	f.lastReq = req
	return f.resp, f.err
}

// fakeDoctorUseCase records the last DoctorRequest and returns the
// configured response/error. Default zero-value yields an empty
// report (no issues) — init-focused tests use it as a stub.
type fakeDoctorUseCase struct {
	called  bool
	lastReq driving.DoctorRequest
	resp    driving.DoctorResponse
	err     error
}

func (f *fakeDoctorUseCase) Check(_ context.Context, req driving.DoctorRequest) (driving.DoctorResponse, error) {
	f.called = true
	f.lastReq = req
	return f.resp, f.err
}

// fakeAddServiceUseCase records the last AddServiceRequest and
// returns the configured response/error. Init/doctor-focused tests
// use the zero value as a stub.
type fakeAddServiceUseCase struct {
	called  bool
	lastReq driving.AddServiceRequest
	resp    driving.AddServiceResponse
	err     error
}

func (f *fakeAddServiceUseCase) Add(_ context.Context, req driving.AddServiceRequest) (driving.AddServiceResponse, error) {
	f.called = true
	f.lastReq = req
	return f.resp, f.err
}

// fakeUpUseCase records the last UpRequest and returns the
// configured response/error. Zero-value yields a stabilized
// empty-services response — a noop stub for unrelated tests.
type fakeUpUseCase struct {
	called  bool
	lastReq driving.UpRequest
	resp    driving.UpResponse
	err     error
}

func (f *fakeUpUseCase) Up(_ context.Context, req driving.UpRequest) (driving.UpResponse, error) {
	f.called = true
	f.lastReq = req
	return f.resp, f.err
}

// fakeDownUseCase records the last DownRequest and returns the
// configured response/error.
type fakeDownUseCase struct {
	called  bool
	lastReq driving.DownRequest
	resp    driving.DownResponse
	err     error
}

func (f *fakeDownUseCase) Down(_ context.Context, req driving.DownRequest) (driving.DownResponse, error) {
	f.called = true
	f.lastReq = req
	return f.resp, f.err
}

func newApp(uc driving.InitProjectUseCase, opts ...cli.Option) *cli.App {
	return cli.New("0.0.0-test", uc, &fakeDoctorUseCase{}, &fakeAddServiceUseCase{}, &fakeUpUseCase{}, &fakeDownUseCase{}, opts...)
}

// newAppWithDoctor is newApp's variant for doctor-focused tests; the
// caller can wire a fake DoctorUseCase explicitly.
func newAppWithDoctor(uc driving.InitProjectUseCase, doctorUC driving.DoctorUseCase, opts ...cli.Option) *cli.App {
	return cli.New("0.0.0-test", uc, doctorUC, &fakeAddServiceUseCase{}, &fakeUpUseCase{}, &fakeDownUseCase{}, opts...)
}

// newAppWithAdd is newApp's variant for add-focused tests.
func newAppWithAdd(uc driving.AddServiceUseCase, opts ...cli.Option) *cli.App {
	return cli.New("0.0.0-test", &fakeInitUseCase{}, &fakeDoctorUseCase{}, uc, &fakeUpUseCase{}, &fakeDownUseCase{}, opts...)
}

// newAppWithUp is newApp's variant for `u-boot up`-focused tests.
func newAppWithUp(uc driving.UpUseCase, opts ...cli.Option) *cli.App {
	return cli.New("0.0.0-test", &fakeInitUseCase{}, &fakeDoctorUseCase{}, &fakeAddServiceUseCase{}, uc, &fakeDownUseCase{}, opts...)
}

// newAppWithDown is newApp's variant for `u-boot down`-focused tests.
func newAppWithDown(uc driving.DownUseCase, opts ...cli.Option) *cli.App {
	return cli.New("0.0.0-test", &fakeInitUseCase{}, &fakeDoctorUseCase{}, &fakeAddServiceUseCase{}, &fakeUpUseCase{}, uc, opts...)
}

func mustProjectName(t *testing.T, raw string) domain.ProjectName {
	t.Helper()
	name, err := domain.NewProjectName(raw)
	if err != nil {
		t.Fatalf("NewProjectName(%q): %v", raw, err)
	}
	return name
}

func TestExecute_Version(t *testing.T) {
	var stdout, stderr bytes.Buffer
	uc := &fakeInitUseCase{}
	err := newApp(uc).Execute(context.Background(), []string{"--version"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Execute --version: %v", err)
	}
	if !strings.Contains(stdout.String(), "0.0.0-test") {
		t.Errorf("--version stdout missing version; got %q", stdout.String())
	}
	if uc.called {
		t.Errorf("--version triggered the use-case")
	}
}

func TestExecute_HelpListsInit(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := newApp(&fakeInitUseCase{}).Execute(context.Background(), []string{"--help"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Execute --help: %v", err)
	}
	if !strings.Contains(stdout.String(), "init") {
		t.Errorf("--help stdout does not list `init`; got %q", stdout.String())
	}
}

func TestExecute_UnknownCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := newApp(&fakeInitUseCase{}).Execute(context.Background(), []string{"frobnicate"}, &stdout, &stderr)
	if err == nil {
		t.Fatalf("Execute frobnicate: expected error")
	}
	if got := cli.ExitCode(err); got != 2 {
		t.Errorf("ExitCode(unknown command) = %d, want 2", got)
	}
}

func TestExecute_UnknownFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := newApp(&fakeInitUseCase{}).Execute(context.Background(), []string{"init", "--no-such-flag"}, &stdout, &stderr)
	if err == nil {
		t.Fatalf("Execute --no-such-flag: expected error")
	}
	if got := cli.ExitCode(err); got != 2 {
		t.Errorf("ExitCode(unknown flag) = %d, want 2", got)
	}
}

func TestExecute_InitWithName(t *testing.T) {
	// fake getwd so the test does not depend on host pwd.
	getwd := func() (string, error) { return "/tmp/x/demo", nil }

	uc := &fakeInitUseCase{
		resp: driving.InitProjectResponse{
			Project: domain.NewProject(mustProjectName(t, "my-service")),
			Created: []string{"docker/", "u-boot.yaml"},
		},
	}
	var stdout, stderr bytes.Buffer
	err := newApp(uc, cli.WithGetwd(getwd)).Execute(context.Background(), []string{"init", "my-service"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Execute init: %v", err)
	}

	if !uc.called {
		t.Fatalf("init did not call use-case")
	}
	if uc.lastReq.Name != "my-service" {
		t.Errorf("init Name = %q, want %q", uc.lastReq.Name, "my-service")
	}
	if uc.lastReq.BaseDir != "/tmp/x/demo" {
		t.Errorf("init BaseDir = %q, want %q", uc.lastReq.BaseDir, "/tmp/x/demo")
	}
	if uc.lastReq.SkipGit {
		t.Errorf("init SkipGit = true, want false")
	}
	out := stdout.String()
	for _, want := range []string{"my-service", "docker/", "u-boot.yaml"} {
		if !strings.Contains(out, want) {
			t.Errorf("init stdout missing %q; got %q", want, out)
		}
	}
}

func TestExecute_InitDerivedName(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }

	uc := &fakeInitUseCase{
		resp: driving.InitProjectResponse{
			Project: domain.NewProject(mustProjectName(t, "demo")),
		},
	}
	var stdout, stderr bytes.Buffer
	err := newApp(uc, cli.WithGetwd(getwd)).Execute(context.Background(), []string{"init"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Execute init: %v", err)
	}
	if uc.lastReq.Name != "" {
		t.Errorf("init Name = %q, want empty (let application derive)", uc.lastReq.Name)
	}
}

func TestExecute_InitNoGitFlag(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }

	uc := &fakeInitUseCase{
		resp: driving.InitProjectResponse{
			Project: domain.NewProject(mustProjectName(t, "demo")),
		},
	}
	var stdout, stderr bytes.Buffer
	err := newApp(uc, cli.WithGetwd(getwd)).Execute(context.Background(), []string{"init", "--no-git"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Execute init --no-git: %v", err)
	}
	if !uc.lastReq.SkipGit {
		t.Errorf("init SkipGit = false, want true (--no-git was passed)")
	}
}

func TestExecute_InitUseCaseErrorPropagates(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }

	uc := &fakeInitUseCase{
		err: driving.ErrProjectExists,
	}
	var stdout, stderr bytes.Buffer
	err := newApp(uc, cli.WithGetwd(getwd)).Execute(context.Background(), []string{"init"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrProjectExists) {
		t.Errorf("Execute init: error %v does not wrap ErrProjectExists", err)
	}
	if got := cli.ExitCode(err); got != 10 {
		t.Errorf("ExitCode(ErrProjectExists) = %d, want 10", got)
	}
}

func TestExecute_InitTooManyArgs(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }

	var stdout, stderr bytes.Buffer
	err := newApp(&fakeInitUseCase{}, cli.WithGetwd(getwd)).Execute(context.Background(), []string{"init", "a", "b"}, &stdout, &stderr)
	if err == nil {
		t.Fatalf("init with two positional args: expected error")
	}
	if got := cli.ExitCode(err); got != 2 {
		t.Errorf("ExitCode(too many args) = %d, want 2", got)
	}
}

// --- T4c: flag wiring + LH-FA-CLI-005A conflict detection ---

func TestExecute_InitForceAndBackupFlags_PassThrough(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }

	uc := &fakeInitUseCase{
		resp: driving.InitProjectResponse{
			Project: domain.NewProject(mustProjectName(t, "demo")),
		},
	}
	var stdout, stderr bytes.Buffer
	err := newApp(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(),
		[]string{"init", "--force", "--backup"},
		&stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !uc.lastReq.Force {
		t.Errorf("Force = false, want true (--force was passed)")
	}
	if !uc.lastReq.Backup {
		t.Errorf("Backup = false, want true (--backup was passed)")
	}
}

func TestExecute_InitAssumeExistingFlag_PassThrough(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }

	uc := &fakeInitUseCase{
		resp: driving.InitProjectResponse{
			Project: domain.NewProject(mustProjectName(t, "demo")),
		},
	}
	var stdout, stderr bytes.Buffer
	err := newApp(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(),
		[]string{"init", "--assume-existing"},
		&stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !uc.lastReq.AssumeExisting {
		t.Errorf("AssumeExisting = false, want true")
	}
}

func TestExecute_InitAssumeExistingFlag_NotGlobal(t *testing.T) {
	// Why: LH-FA-CLI-005A §238 — --assume-existing is init-only.
	// Putting it on the root command must fail with a usage error.
	var stdout, stderr bytes.Buffer
	err := newApp(&fakeInitUseCase{}).Execute(
		context.Background(),
		[]string{"--assume-existing"},
		&stdout, &stderr,
	)
	if err == nil {
		t.Fatalf("Execute --assume-existing on root: expected error")
	}
	if got := cli.ExitCode(err); got != 2 {
		t.Errorf("ExitCode(root --assume-existing) = %d, want 2", got)
	}
}

func TestExecute_YesAndNoInteractive_Conflict(t *testing.T) {
	// Why: LH-FA-CLI-005A §235 — `--yes` and `--no-interactive` are
	// mutually exclusive. The conflict surfaces via the CLI sentinel
	// (not the use-case) → exit code 2.
	getwd := func() (string, error) { return "/tmp/x/demo", nil }

	uc := &fakeInitUseCase{}
	var stdout, stderr bytes.Buffer
	err := newApp(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(),
		[]string{"--yes", "--no-interactive", "init"},
		&stdout, &stderr,
	)
	if err == nil {
		t.Fatalf("Execute --yes --no-interactive: expected error")
	}
	if !errors.Is(err, cli.ErrConflictingModeFlags) {
		t.Errorf("Execute --yes --no-interactive: error %v does not wrap ErrConflictingModeFlags", err)
	}
	if got := cli.ExitCode(err); got != 2 {
		t.Errorf("ExitCode(ErrConflictingModeFlags) = %d, want 2", got)
	}
	if uc.called {
		t.Errorf("conflict check must short-circuit before the use-case runs")
	}
}

func TestExecute_YesAlone_OnDeterministicPath_NoEffect(t *testing.T) {
	// Why: LH-FA-CLI-005A §247 — `--yes` on a deterministic path is
	// a no-op. The use-case still runs with its plain request; no
	// conflict check fires.
	getwd := func() (string, error) { return "/tmp/x/demo", nil }

	uc := &fakeInitUseCase{
		resp: driving.InitProjectResponse{
			Project: domain.NewProject(mustProjectName(t, "demo")),
		},
	}
	var stdout, stderr bytes.Buffer
	err := newApp(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(),
		[]string{"--yes", "init"},
		&stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !uc.called {
		t.Errorf("--yes on init must still invoke the use-case")
	}
}

func TestExecute_NoInteractiveAlone_OnDeterministicPath_NoEffect(t *testing.T) {
	// Mirror of TestExecute_YesAlone_OnDeterministicPath_NoEffect.
	getwd := func() (string, error) { return "/tmp/x/demo", nil }

	uc := &fakeInitUseCase{
		resp: driving.InitProjectResponse{
			Project: domain.NewProject(mustProjectName(t, "demo")),
		},
	}
	var stdout, stderr bytes.Buffer
	err := newApp(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(),
		[]string{"--no-interactive", "init"},
		&stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !uc.called {
		t.Errorf("--no-interactive on init must still invoke the use-case")
	}
}

func TestExecute_InitBackupsAppearInSummary(t *testing.T) {
	// Why: the LH-FA-INIT-005 §609 affected-paths line is emitted by
	// the application layer via the progress writer; the CLI's
	// printInitSummary additionally lists the resulting backup
	// actions so the user can see where their originals went.
	getwd := func() (string, error) { return "/tmp/x/demo", nil }

	uc := &fakeInitUseCase{
		resp: driving.InitProjectResponse{
			Project: domain.NewProject(mustProjectName(t, "demo")),
			Created: []string{"u-boot.yaml"},
			Backups: []driving.BackupAction{
				{Original: ".env.example", Backup: "/tmp/x/demo/.env.example.bak"},
			},
		},
	}
	var stdout, stderr bytes.Buffer
	err := newApp(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(),
		[]string{"init", "--backup"},
		&stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	out := stdout.String()
	if !strings.Contains(out, "Backups:") {
		t.Errorf("output missing Backups section: %q", out)
	}
	if !strings.Contains(out, ".env.example → /tmp/x/demo/.env.example.bak") {
		t.Errorf("output missing backup-action line: %q", out)
	}
}

func TestExecute_InitNoBackups_NoBackupsSection(t *testing.T) {
	// Why: defensive — empty backups list must not render an empty
	// "Backups:" header.
	getwd := func() (string, error) { return "/tmp/x/demo", nil }

	uc := &fakeInitUseCase{
		resp: driving.InitProjectResponse{
			Project: domain.NewProject(mustProjectName(t, "demo")),
			Created: []string{"u-boot.yaml"},
		},
	}
	var stdout, stderr bytes.Buffer
	err := newApp(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(),
		[]string{"init"},
		&stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if strings.Contains(stdout.String(), "Backups:") {
		t.Errorf("fresh init should not render a Backups section: %q", stdout.String())
	}
}

func TestExecute_FlagsDoNotLeakAcrossInvocations(t *testing.T) {
	// Why: review finding #1 — App.yes / App.noInteractive are bound
	// to PersistentFlags by pointer. Cobra rebuilds the root per
	// Execute call, BoolVar(&a.yes, "yes", false, …) writes the
	// default (false) into the bound variable on each registration.
	// Pin that contract so a future migration to a long-lived cmd
	// tree cannot silently leak state.
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeInitUseCase{
		resp: driving.InitProjectResponse{
			Project: domain.NewProject(mustProjectName(t, "demo")),
		},
	}
	app := newApp(uc, cli.WithGetwd(getwd))

	// First call: --yes init. Should not error.
	var out1, err1 bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--yes", "init"}, &out1, &err1); err != nil {
		t.Fatalf("Execute(--yes init): %v", err)
	}

	// Second call: --no-interactive init (no --yes). If the prior
	// --yes leaked, the conflict check fires; with proper
	// per-Execute defaulting it does not.
	var out2, err2 bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--no-interactive", "init"}, &out2, &err2); err != nil {
		t.Errorf("Execute(--no-interactive init): unexpected conflict — flags leaked from prior call: %v", err)
	}
}

func TestExecute_YesAndNoInteractive_Conflict_SubcommandThenFlag(t *testing.T) {
	// Why: review finding #2 — both flag orderings (root-then-sub
	// and sub-then-flag) must trip the conflict check. Cobra
	// inherits PersistentFlags into the subcommand FlagSet, but
	// pin it executable so a flag-restructure can't slip past.
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeInitUseCase{}
	var stdout, stderr bytes.Buffer
	err := newApp(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(),
		[]string{"init", "--yes", "--no-interactive"},
		&stdout, &stderr,
	)
	if !errors.Is(err, cli.ErrConflictingModeFlags) {
		t.Errorf("Execute init --yes --no-interactive: error %v does not wrap ErrConflictingModeFlags", err)
	}
}

func TestExecute_NoInteractive_WithForceBackup_ReInit_Succeeds(t *testing.T) {
	// Why: review finding #10 — production CI scenario. The
	// deterministic re-init path (--no-interactive --force --backup
	// on an existing project) must succeed and surface the backup
	// list to the user. Mirrors the manual smoketest in the T4c
	// commit message, executable now.
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeInitUseCase{
		resp: driving.InitProjectResponse{
			Project: domain.NewProject(mustProjectName(t, "demo")),
			Created: []string{"u-boot.yaml", "compose.yaml"},
			Backups: []driving.BackupAction{
				{Original: "compose.yaml", Backup: "/tmp/x/demo/compose.yaml.bak"},
			},
		},
	}
	var stdout, stderr bytes.Buffer
	err := newApp(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(),
		[]string{"--no-interactive", "init", "--force", "--backup"},
		&stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !uc.lastReq.Force || !uc.lastReq.Backup {
		t.Errorf("flags not propagated: Force=%v Backup=%v", uc.lastReq.Force, uc.lastReq.Backup)
	}
	if !strings.Contains(stdout.String(), "Backups:") {
		t.Errorf("CI re-init output missing Backups section: %q", stdout.String())
	}
}

func TestExecute_TooManyArgs_WithYesFlag_StillUsageError(t *testing.T) {
	// Why: review finding #11 — mode flags must not change positional
	// validation. `--yes init a b` must still trip the cobra
	// MaximumNArgs guard with exit-code 2.
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	var stdout, stderr bytes.Buffer
	err := newApp(&fakeInitUseCase{}, cli.WithGetwd(getwd)).Execute(
		context.Background(),
		[]string{"--yes", "init", "a", "b"},
		&stdout, &stderr,
	)
	if err == nil {
		t.Fatalf("init with two positional args: expected error")
	}
	if got := cli.ExitCode(err); got != 2 {
		t.Errorf("ExitCode(too many args with --yes) = %d, want 2", got)
	}
}

func TestExecute_InitAssumeExisting_NoLongerEmitsM3Note(t *testing.T) {
	// Why: the M3 stderr note ("--assume-existing has no effect in M3")
	// was removed with the M4 soft-existing-detection slice — the flag
	// is now load-bearing. Regression-guard against accidentally
	// re-adding the note (which would be obnoxious on every M4+ run).
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeInitUseCase{
		resp: driving.InitProjectResponse{
			Project: domain.NewProject(mustProjectName(t, "demo")),
		},
	}
	var stdout, stderr bytes.Buffer
	err := newApp(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(),
		[]string{"init", "--assume-existing"},
		&stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if strings.Contains(stderr.String(), "no effect in M3") {
		t.Errorf("stderr still contains the removed M3 NoOp note: %q", stderr.String())
	}
	if !uc.called {
		t.Errorf("--assume-existing must not block the use-case")
	}
	if !uc.lastReq.AssumeExisting {
		t.Errorf("AssumeExisting = false, want true (flag must pass through)")
	}
}

func TestExecute_NoInteractive_PassThrough(t *testing.T) {
	// Why: M4 soft-detection — --no-interactive must propagate into
	// req.NoInteractive so the service skips the prompt path
	// (LH-FA-INIT-004 §247).
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeInitUseCase{
		resp: driving.InitProjectResponse{
			Project: domain.NewProject(mustProjectName(t, "demo")),
		},
	}
	var stdout, stderr bytes.Buffer
	err := newApp(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(),
		[]string{"--no-interactive", "init"},
		&stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !uc.lastReq.NoInteractive {
		t.Errorf("NoInteractive = false in propagated request, want true")
	}
}

func TestExecute_NoAssumeExisting_NoStderrNote(t *testing.T) {
	// Why: defensive — no stderr note should fire on a plain init
	// (the old M3 NoOp note was the only producer; this test keeps
	// the contract even if a future slice adds new emit paths).
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeInitUseCase{
		resp: driving.InitProjectResponse{
			Project: domain.NewProject(mustProjectName(t, "demo")),
		},
	}
	var stdout, stderr bytes.Buffer
	err := newApp(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(),
		[]string{"init"},
		&stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if stderr.Len() != 0 {
		t.Errorf("unexpected stderr on plain init: %q", stderr.String())
	}
}

func TestExitCode_BaseMappings(t *testing.T) {
	// Tests against the LH-FA-CLI-006-sentinel mappings; the actual
	// usage-error-string detection is covered by the integration
	// tests TestExecute_UnknownCommand / TestExecute_UnknownFlag /
	// TestExecute_InitTooManyArgs which drive the real Cobra path
	// (and therefore catch a Cobra-string drift on upgrade).
	cases := []struct {
		name string
		err  error
		want int
	}{
		{"nil", nil, 0},
		{"generic", errors.New("boom"), 1},
		{"ErrProjectExists (validation)", driving.ErrProjectExists, 10},
		{"ErrBaseDirMissing (validation)", driving.ErrBaseDirMissing, 10},
		{"ErrInvalidProjectName (validation)", domain.ErrInvalidProjectName, 10},
		{"ErrBackupUnsupportedKind (validation)", driving.ErrBackupUnsupportedKind, 10},
		{"ErrForceRequiresBackup (validation)", driving.ErrForceRequiresBackup, 10},
		{"ErrFileExists (validation)", driving.ErrFileExists, 10},
		{"wrapped ErrProjectExists", fmt.Errorf("ctx: %w", driving.ErrProjectExists), 10},
		{"ErrConflictingModeFlags (usage)", cli.ErrConflictingModeFlags, 2},
		{"wrapped ErrConflictingModeFlags", fmt.Errorf("ctx: %w", cli.ErrConflictingModeFlags), 2},
		{"ErrBackupSuffixExhausted (fs)", driving.ErrBackupSuffixExhausted, 14},
		{"ErrBackupSourceMissing (fs)", driving.ErrBackupSourceMissing, 14},
		{"wrapped ErrBackupSuffixExhausted", fmt.Errorf("ctx: %w", driving.ErrBackupSuffixExhausted), 14},
		{"ErrDoctorFailures (doctor)", cli.ErrDoctorFailures, 11},
		{"wrapped ErrDoctorFailures", fmt.Errorf("ctx: %w", cli.ErrDoctorFailures), 11},
		{"ErrProjectNotInitialized (add)", driving.ErrProjectNotInitialized, 10},
		{"ErrServiceUnsupported (add)", driving.ErrServiceUnsupported, 10},
		{"ErrServiceInconsistent (add)", driving.ErrServiceInconsistent, 10},
		{"ErrInvalidServiceName (add)", domain.ErrInvalidServiceName, 10},
		{"wrapped ErrServiceUnsupported", fmt.Errorf("ctx: %w", driving.ErrServiceUnsupported), 10},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := cli.ExitCode(tc.err); got != tc.want {
				t.Errorf("ExitCode(%v) = %d, want %d", tc.err, got, tc.want)
			}
		})
	}
}

// ----------------------------------------------------------------------------
// `u-boot doctor` subcommand
// ----------------------------------------------------------------------------

func TestExecute_Doctor_NoIssues_ExitOK(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x", nil }
	doctorUC := &fakeDoctorUseCase{
		resp: driving.DoctorResponse{
			Report: domain.DiagnosticReport{Items: []domain.Diagnostic{
				{ID: "fs.write-permissions", Severity: domain.SeverityOK, Message: "BaseDir is writable."},
			}},
		},
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithDoctor(&fakeInitUseCase{}, doctorUC, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"doctor"}, &stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if cli.ExitCode(err) != 0 {
		t.Errorf("ExitCode(nil-error doctor) = %d, want 0", cli.ExitCode(err))
	}
	if !doctorUC.called {
		t.Errorf("DoctorUseCase.Check not invoked")
	}
	if !strings.Contains(stdout.String(), "BaseDir is writable") {
		t.Errorf("stdout missing diagnostic body: %q", stdout.String())
	}
}

func TestExecute_Doctor_ErrorReport_Exit11(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x", nil }
	doctorUC := &fakeDoctorUseCase{
		resp: driving.DoctorResponse{
			Report: domain.DiagnosticReport{Items: []domain.Diagnostic{
				{ID: "docker.reachable", Severity: domain.SeverityError, Message: "docker daemon is not reachable.", Hint: "Start Docker."},
				{ID: "fs.write-permissions", Severity: domain.SeverityOK, Message: "BaseDir is writable."},
			}},
		},
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithDoctor(&fakeInitUseCase{}, doctorUC, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"doctor"}, &stdout, &stderr,
	)
	if !errors.Is(err, cli.ErrDoctorFailures) {
		t.Fatalf("err = %v, want wrapped ErrDoctorFailures", err)
	}
	if got := cli.ExitCode(err); got != 11 {
		t.Errorf("ExitCode = %d, want 11", got)
	}
	out := stdout.String()
	if !strings.Contains(out, "docker daemon is not reachable") {
		t.Errorf("stdout missing error item: %q", out)
	}
	if !strings.Contains(out, "Start Docker") {
		t.Errorf("stdout missing error hint: %q", out)
	}
}

func TestExecute_Doctor_WarnNonStrict_Exit0(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x", nil }
	doctorUC := &fakeDoctorUseCase{
		resp: driving.DoctorResponse{
			Report: domain.DiagnosticReport{Items: []domain.Diagnostic{
				{ID: "uboot.yaml.valid", Severity: domain.SeverityWarn, Message: "u-boot.yaml not present."},
			}},
		},
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithDoctor(&fakeInitUseCase{}, doctorUC, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"doctor"}, &stdout, &stderr,
	)
	if err != nil {
		t.Errorf("err = %v, want nil (warn without --strict)", err)
	}
	if got := cli.ExitCode(err); got != 0 {
		t.Errorf("ExitCode = %d, want 0", got)
	}
}

func TestExecute_Doctor_WarnStrict_Exit11(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x", nil }
	doctorUC := &fakeDoctorUseCase{
		resp: driving.DoctorResponse{
			Report: domain.DiagnosticReport{Items: []domain.Diagnostic{
				{ID: "uboot.yaml.valid", Severity: domain.SeverityWarn, Message: "u-boot.yaml not present."},
			}},
		},
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithDoctor(&fakeInitUseCase{}, doctorUC, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"doctor", "--strict"}, &stdout, &stderr,
	)
	if !errors.Is(err, cli.ErrDoctorFailures) {
		t.Fatalf("err = %v, want wrapped ErrDoctorFailures with --strict", err)
	}
	if got := cli.ExitCode(err); got != 11 {
		t.Errorf("ExitCode = %d, want 11", got)
	}
}

func TestExecute_Doctor_Quiet_HidesOKItems(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x", nil }
	doctorUC := &fakeDoctorUseCase{
		resp: driving.DoctorResponse{
			Report: domain.DiagnosticReport{Items: []domain.Diagnostic{
				{ID: "fs.write-permissions", Severity: domain.SeverityOK, Message: "BaseDir is writable."},
				{ID: "uboot.yaml.valid", Severity: domain.SeverityWarn, Message: "u-boot.yaml not present."},
			}},
		},
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithDoctor(&fakeInitUseCase{}, doctorUC, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"--quiet", "doctor"}, &stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	out := stdout.String()
	if strings.Contains(out, "BaseDir is writable") {
		t.Errorf("--quiet did not hide OK item: %q", out)
	}
	if !strings.Contains(out, "u-boot.yaml not present") {
		t.Errorf("--quiet hid the warn item it should keep: %q", out)
	}
	// Summary still includes all counts even in --quiet mode.
	if !strings.Contains(out, "0 error, 1 warn, 1 ok") {
		t.Errorf("summary line missing or wrong: %q", out)
	}
}

func TestExecute_Doctor_SortedByIssuesFirst(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x", nil }
	doctorUC := &fakeDoctorUseCase{
		resp: driving.DoctorResponse{
			Report: domain.DiagnosticReport{Items: []domain.Diagnostic{
				{ID: "a.ok", Severity: domain.SeverityOK, Message: "ok-msg"},
				{ID: "b.error", Severity: domain.SeverityError, Message: "err-msg"},
				{ID: "c.warn", Severity: domain.SeverityWarn, Message: "warn-msg"},
			}},
		},
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithDoctor(&fakeInitUseCase{}, doctorUC, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"doctor"}, &stdout, &stderr,
	)
	if !errors.Is(err, cli.ErrDoctorFailures) {
		t.Fatalf("err = %v, want ErrDoctorFailures (the report has an error)", err)
	}
	out := stdout.String()
	// Error (b.error) must come before warn (c.warn) must come before ok (a.ok).
	bErr := strings.Index(out, "err-msg")
	cWarn := strings.Index(out, "warn-msg")
	aOK := strings.Index(out, "ok-msg")
	if bErr >= cWarn || cWarn >= aOK {
		t.Errorf("rendered order wrong (err=%d warn=%d ok=%d):\n%s", bErr, cWarn, aOK, out)
	}
}

func TestExecute_Doctor_TooManyArgs_Exit2(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x", nil }
	var stdout, stderr bytes.Buffer
	err := newAppWithDoctor(&fakeInitUseCase{}, &fakeDoctorUseCase{}, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"doctor", "extra-arg"}, &stdout, &stderr,
	)
	if err == nil {
		t.Fatal("expected usage error for doctor with positional arg")
	}
	if got := cli.ExitCode(err); got != 2 {
		t.Errorf("ExitCode = %d, want 2", got)
	}
}

// --- M5-T6: `u-boot add <service>` ---------------------------------

func TestExecute_Add_HappyPathRegister(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeAddServiceUseCase{
		resp: driving.AddServiceResponse{
			ServiceName: mustServiceName(t, "postgres"),
			PriorState:  domain.ServiceStateUnregistered,
			State:       domain.ServiceStateActive,
			Changed:     []string{"u-boot.yaml", "compose.yaml", ".env.example"},
		},
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithAdd(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"add", "postgres"}, &stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !uc.called {
		t.Fatalf("use-case not called")
	}
	if uc.lastReq.BaseDir != "/tmp/x/demo" {
		t.Errorf("BaseDir = %q, want /tmp/x/demo", uc.lastReq.BaseDir)
	}
	if uc.lastReq.ServiceName.String() != "postgres" {
		t.Errorf("ServiceName = %q, want postgres", uc.lastReq.ServiceName.String())
	}
	out := stdout.String()
	for _, want := range []string{
		`Added service "postgres".`,
		"  - u-boot.yaml",
		"  - compose.yaml",
		"  - .env.example",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("stdout missing %q; got:\n%s", want, out)
		}
	}
}

func TestExecute_Add_AlreadyActiveNoOp(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeAddServiceUseCase{
		resp: driving.AddServiceResponse{
			ServiceName: mustServiceName(t, "postgres"),
			PriorState:  domain.ServiceStateActive,
			State:       domain.ServiceStateActive,
			Changed:     nil,
		},
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithAdd(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"add", "postgres"}, &stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	out := stdout.String()
	if !strings.Contains(out, `Service "postgres" is already active; no changes.`) {
		t.Errorf("missing no-op summary; got:\n%s", out)
	}
}

func TestExecute_Add_RepairArtifactsSummary(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeAddServiceUseCase{
		resp: driving.AddServiceResponse{
			ServiceName: mustServiceName(t, "postgres"),
			PriorState:  domain.ServiceStateActive,
			State:       domain.ServiceStateActive,
			Changed:     []string{".env.example"},
		},
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithAdd(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"add", "postgres"}, &stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	out := stdout.String()
	for _, want := range []string{
		`Repaired service "postgres" artefacts.`,
		"  - .env.example",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("stdout missing %q; got:\n%s", want, out)
		}
	}
}

func TestExecute_Add_InvalidServiceName_Code10(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeAddServiceUseCase{}
	var stdout, stderr bytes.Buffer
	err := newAppWithAdd(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"add", "INVALID NAME"}, &stdout, &stderr,
	)
	if err == nil {
		t.Fatal("expected ErrInvalidServiceName, got nil")
	}
	if !errors.Is(err, domain.ErrInvalidServiceName) {
		t.Errorf("err = %v, want ErrInvalidServiceName", err)
	}
	if got := cli.ExitCode(err); got != 10 {
		t.Errorf("ExitCode = %d, want 10", got)
	}
	if uc.called {
		t.Error("use-case must not run when service name fails validation")
	}
}

func TestExecute_Add_UnsupportedService_Code10(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeAddServiceUseCase{
		err: fmt.Errorf("%w: redis", driving.ErrServiceUnsupported),
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithAdd(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"add", "redis"}, &stdout, &stderr,
	)
	if !errors.Is(err, driving.ErrServiceUnsupported) {
		t.Fatalf("err = %v, want ErrServiceUnsupported", err)
	}
	if got := cli.ExitCode(err); got != 10 {
		t.Errorf("ExitCode = %d, want 10", got)
	}
}

func TestExecute_Add_ProjectNotInitialized_Code10(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeAddServiceUseCase{
		err: fmt.Errorf("%w: u-boot.yaml missing", driving.ErrProjectNotInitialized),
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithAdd(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"add", "postgres"}, &stdout, &stderr,
	)
	if !errors.Is(err, driving.ErrProjectNotInitialized) {
		t.Fatalf("err = %v, want ErrProjectNotInitialized", err)
	}
	if got := cli.ExitCode(err); got != 10 {
		t.Errorf("ExitCode = %d, want 10", got)
	}
}

func TestExecute_Add_ServiceInconsistent_Code10(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeAddServiceUseCase{
		err: fmt.Errorf("%w: orphan block", driving.ErrServiceInconsistent),
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithAdd(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"add", "postgres"}, &stdout, &stderr,
	)
	if !errors.Is(err, driving.ErrServiceInconsistent) {
		t.Fatalf("err = %v, want ErrServiceInconsistent", err)
	}
	if got := cli.ExitCode(err); got != 10 {
		t.Errorf("ExitCode = %d, want 10", got)
	}
}

func TestExecute_Add_NoArgs_Code2(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := newAppWithAdd(&fakeAddServiceUseCase{}).Execute(
		context.Background(), []string{"add"}, &stdout, &stderr,
	)
	if err == nil {
		t.Fatal("expected usage error for `add` without positional arg")
	}
	if got := cli.ExitCode(err); got != 2 {
		t.Errorf("ExitCode = %d, want 2", got)
	}
}

func TestExecute_Add_TooManyArgs_Code2(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := newAppWithAdd(&fakeAddServiceUseCase{}).Execute(
		context.Background(), []string{"add", "postgres", "extra"}, &stdout, &stderr,
	)
	if err == nil {
		t.Fatal("expected usage error for too many positional args")
	}
	if got := cli.ExitCode(err); got != 2 {
		t.Errorf("ExitCode = %d, want 2", got)
	}
}

func TestExecute_Add_ConflictingModeFlags_Code2(t *testing.T) {
	uc := &fakeAddServiceUseCase{}
	var stdout, stderr bytes.Buffer
	err := newAppWithAdd(uc, cli.WithGetwd(func() (string, error) { return "/tmp/x", nil })).Execute(
		context.Background(), []string{"--yes", "--no-interactive", "add", "postgres"}, &stdout, &stderr,
	)
	if !errors.Is(err, cli.ErrConflictingModeFlags) {
		t.Fatalf("err = %v, want ErrConflictingModeFlags", err)
	}
	if got := cli.ExitCode(err); got != 2 {
		t.Errorf("ExitCode = %d, want 2", got)
	}
	if uc.called {
		t.Error("use-case must not run when mode-flag check fails")
	}
}

func TestExecute_Add_GetwdFailure_Wrapped(t *testing.T) {
	uc := &fakeAddServiceUseCase{}
	var stdout, stderr bytes.Buffer
	err := newAppWithAdd(uc, cli.WithGetwd(func() (string, error) {
		return "", errors.New("getwd boom")
	})).Execute(
		context.Background(), []string{"add", "postgres"}, &stdout, &stderr,
	)
	if err == nil {
		t.Fatal("expected wrapped getwd error")
	}
	if !strings.Contains(err.Error(), "getwd boom") {
		t.Errorf("err = %v, want it to wrap getwd boom", err)
	}
	if uc.called {
		t.Error("use-case must not run when getwd fails")
	}
}

// mustServiceName is the test helper analogue of mustProjectName.
func mustServiceName(t *testing.T, raw string) domain.ServiceName {
	t.Helper()
	name, err := domain.NewServiceName(raw)
	if err != nil {
		t.Fatalf("NewServiceName(%q): %v", raw, err)
	}
	return name
}

// --- M6-T7 up subcommand pin tests ---

func TestExecute_Up_HappyPath_RendersStatusTable(t *testing.T) {
	t.Parallel()
	uc := &fakeUpUseCase{
		resp: driving.UpResponse{Result: domain.UpResult{
			Stabilized: true,
			Services: []domain.ServiceStatus{
				{Name: "postgres", ContainerStatus: domain.StateRunning, Port: "5432:5432", Healthcheck: "healthy"},
			},
		}},
	}
	getwd := func() (string, error) { return "/tmp/proj", nil }
	var stdout, stderr bytes.Buffer
	err := newAppWithUp(uc, cli.WithGetwd(getwd)).Execute(context.Background(), []string{"up"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("up: %v", err)
	}
	if !uc.called {
		t.Error("use-case not called")
	}
	// Default --timeout=60 → 60s.
	if uc.lastReq.Timeout != 60*time.Second {
		t.Errorf("Timeout = %v, want 60s (default)", uc.lastReq.Timeout)
	}
	out := stdout.String()
	if !strings.Contains(out, "postgres") || !strings.Contains(out, "healthy") {
		t.Errorf("expected status table to render service row, got: %q", out)
	}
}

func TestExecute_Up_Timeout0_PassesFireAndForget(t *testing.T) {
	t.Parallel()
	uc := &fakeUpUseCase{
		resp: driving.UpResponse{Result: domain.UpResult{
			Stabilized: false,
			Diagnostics: []domain.Diagnostic{
				{ID: "up.fire-and-forget", Severity: domain.SeverityInfo, Message: "started"},
			},
		}},
	}
	getwd := func() (string, error) { return "/tmp/proj", nil }
	var stdout, stderr bytes.Buffer
	err := newAppWithUp(uc, cli.WithGetwd(getwd)).Execute(context.Background(), []string{"up", "--timeout=0"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("up --timeout=0: %v", err)
	}
	if uc.lastReq.Timeout != 0 {
		t.Errorf("Timeout = %v, want 0 (fire-and-forget)", uc.lastReq.Timeout)
	}
	out := stdout.String()
	// Fire-and-forget: no status table.
	if strings.Contains(out, "SERVICE") {
		t.Errorf("status table leaked into --timeout=0 output: %q", out)
	}
	// But the info diagnostic should be shown.
	if !strings.Contains(out, "started") {
		t.Errorf("fire-and-forget info diagnostic missing: %q", out)
	}
}

func TestExecute_Up_NegativeTimeout_ReturnsExitCode2(t *testing.T) {
	t.Parallel()
	uc := &fakeUpUseCase{}
	getwd := func() (string, error) { return "/tmp/proj", nil }
	var stdout, stderr bytes.Buffer
	err := newAppWithUp(uc, cli.WithGetwd(getwd)).Execute(context.Background(), []string{"up", "--timeout=-1"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error for --timeout=-1")
	}
	if !errors.Is(err, cli.ErrInvalidTimeout) {
		t.Errorf("expected ErrInvalidTimeout, got: %v", err)
	}
	if got := cli.ExitCode(err); got != 2 {
		t.Errorf("ExitCode = %d, want 2", got)
	}
	if uc.called {
		t.Error("use-case must not run with invalid timeout")
	}
}

func TestExecute_Up_ProjectNotInitialized_ReturnsExitCode10(t *testing.T) {
	t.Parallel()
	uc := &fakeUpUseCase{err: driving.ErrProjectNotInitialized}
	getwd := func() (string, error) { return "/tmp/proj", nil }
	var stdout, stderr bytes.Buffer
	err := newAppWithUp(uc, cli.WithGetwd(getwd)).Execute(context.Background(), []string{"up"}, &stdout, &stderr)
	if got := cli.ExitCode(err); got != 10 {
		t.Errorf("ExitCode = %d, want 10 (got err: %v)", got, err)
	}
}

// --- M6-T7 down subcommand pin tests ---

func TestExecute_Down_NoVolumes_HappyPath(t *testing.T) {
	t.Parallel()
	uc := &fakeDownUseCase{resp: driving.DownResponse{RemovedVolumes: false}}
	getwd := func() (string, error) { return "/tmp/proj", nil }
	var stdout, stderr bytes.Buffer
	err := newAppWithDown(uc, cli.WithGetwd(getwd)).Execute(context.Background(), []string{"down"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("down: %v", err)
	}
	if uc.lastReq.RemoveVolumes {
		t.Errorf("RemoveVolumes = true, want false")
	}
	if !strings.Contains(stdout.String(), "environment stopped") {
		t.Errorf("expected success message, got: %q", stdout.String())
	}
}

func TestExecute_Down_VolumesYes_BypassesConfirm(t *testing.T) {
	t.Parallel()
	uc := &fakeDownUseCase{resp: driving.DownResponse{RemovedVolumes: true}}
	getwd := func() (string, error) { return "/tmp/proj", nil }
	var stdout, stderr bytes.Buffer
	err := newAppWithDown(uc, cli.WithGetwd(getwd)).Execute(context.Background(), []string{"down", "--volumes", "--yes"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("down --volumes --yes: %v", err)
	}
	if !uc.lastReq.RemoveVolumes {
		t.Errorf("RemoveVolumes = false, want true")
	}
	if !uc.lastReq.AssumeYes {
		t.Errorf("AssumeYes = false, want true")
	}
	if !strings.Contains(stdout.String(), "volumes removed") {
		t.Errorf("expected volumes-removed message, got: %q", stdout.String())
	}
}

func TestExecute_Down_VolumesConfirmRefused_ReturnsExitCode10(t *testing.T) {
	t.Parallel()
	uc := &fakeDownUseCase{err: driving.ErrConfirmationRequired}
	getwd := func() (string, error) { return "/tmp/proj", nil }
	var stdout, stderr bytes.Buffer
	err := newAppWithDown(uc, cli.WithGetwd(getwd)).Execute(context.Background(), []string{"down", "--volumes"}, &stdout, &stderr)
	if got := cli.ExitCode(err); got != 10 {
		t.Errorf("ExitCode = %d, want 10 (got err: %v)", got, err)
	}
}

func TestExecute_Down_YesAndNoInteractive_ReturnsExitCode2(t *testing.T) {
	t.Parallel()
	// §235 mutual exclusion fires before the use case.
	uc := &fakeDownUseCase{}
	getwd := func() (string, error) { return "/tmp/proj", nil }
	var stdout, stderr bytes.Buffer
	err := newAppWithDown(uc, cli.WithGetwd(getwd)).Execute(context.Background(), []string{"down", "--volumes", "--yes", "--no-interactive"}, &stdout, &stderr)
	if !errors.Is(err, cli.ErrConflictingModeFlags) {
		t.Errorf("expected ErrConflictingModeFlags, got: %v", err)
	}
	if got := cli.ExitCode(err); got != 2 {
		t.Errorf("ExitCode = %d, want 2", got)
	}
	if uc.called {
		t.Error("use-case must not run on flag-conflict")
	}
}

// --- ExitCode pin tests for the new M6 sentinels ---

func TestExitCode_M6Sentinels(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		err  error
		want int
	}{
		{"driven.ErrDockerUnavailable", driven.ErrDockerUnavailable, 11},
		{"driven.ErrComposeRuntime", driven.ErrComposeRuntime, 12},
		{"driving.ErrStabilizationTimeout", driving.ErrStabilizationTimeout, 12},
		{"driving.ErrComposeFileMissing", driving.ErrComposeFileMissing, 10},
		{"driving.ErrConfirmationRequired", driving.ErrConfirmationRequired, 10},
		{"cli.ErrInvalidTimeout", cli.ErrInvalidTimeout, 2},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := cli.ExitCode(tc.err); got != tc.want {
				t.Errorf("ExitCode(%v) = %d, want %d", tc.err, got, tc.want)
			}
			// Pin sentinel-chain survival: wrap and re-check.
			wrapped := fmt.Errorf("cli: %w", tc.err)
			if got := cli.ExitCode(wrapped); got != tc.want {
				t.Errorf("ExitCode(wrap(%v)) = %d, want %d", tc.err, got, tc.want)
			}
		})
	}
}
