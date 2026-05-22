package cli_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
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

func newApp(uc driving.InitProjectUseCase, opts ...cli.Option) *cli.App {
	return cli.New("0.0.0-test", uc, opts...)
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

func TestExecute_InitAssumeExisting_EmitsStderrNote(t *testing.T) {
	// Why: review finding #5 — silent NoOp would mislead the user.
	// In M3 the flag is accepted but has no behavioural effect; the
	// CLI emits a one-line note on stderr so the inactivity is
	// visible. Use case still runs (flag is forward-compatible).
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
	if !strings.Contains(stderr.String(), "--assume-existing has no effect in M3") {
		t.Errorf("stderr missing M3 NoOp note: %q", stderr.String())
	}
	if !uc.called {
		t.Errorf("--assume-existing note must not block the use-case")
	}
}

func TestExecute_NoAssumeExisting_NoStderrNote(t *testing.T) {
	// Why: defensive — the M3-NoOp note must NOT fire when the flag
	// is absent (would be obnoxious noise on every plain init).
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
		{"ErrBackupTooLarge (fs)", driving.ErrBackupTooLarge, 14},
		{"wrapped ErrBackupSuffixExhausted", fmt.Errorf("ctx: %w", driving.ErrBackupSuffixExhausted), 14},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := cli.ExitCode(tc.err); got != tc.want {
				t.Errorf("ExitCode(%v) = %d, want %d", tc.err, got, tc.want)
			}
		})
	}
}
