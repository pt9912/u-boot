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

// fakeGenerateUseCase records the last GenerateRequest and returns
// the configured response/error.
type fakeGenerateUseCase struct {
	called  bool
	lastReq driving.GenerateRequest
	resp    driving.GenerateResponse
	err     error
}

func (f *fakeGenerateUseCase) Generate(_ context.Context, req driving.GenerateRequest) (driving.GenerateResponse, error) {
	f.called = true
	f.lastReq = req
	return f.resp, f.err
}

// fakeConfigUseCase records every Get / Set / Show invocation and
// returns the configured response/error. Three independent fields
// per method so a single test can wire (e.g.) a happy Get + an
// error Set without one path interfering with the other.
type fakeConfigUseCase struct {
	getCalled  bool
	getReq     driving.ConfigGetRequest
	getResp    driving.ConfigGetResponse
	getErr     error
	setCalled  bool
	setReq     driving.ConfigSetRequest
	setResp    driving.ConfigSetResponse
	setErr     error
	showCalled bool
	showReq    driving.ConfigShowRequest
	showResp   driving.ConfigShowResponse
	showErr    error
}

func (f *fakeConfigUseCase) Get(_ context.Context, req driving.ConfigGetRequest) (driving.ConfigGetResponse, error) {
	f.getCalled = true
	f.getReq = req
	return f.getResp, f.getErr
}

func (f *fakeConfigUseCase) Set(_ context.Context, req driving.ConfigSetRequest) (driving.ConfigSetResponse, error) {
	f.setCalled = true
	f.setReq = req
	return f.setResp, f.setErr
}

func (f *fakeConfigUseCase) Show(_ context.Context, req driving.ConfigShowRequest) (driving.ConfigShowResponse, error) {
	f.showCalled = true
	f.showReq = req
	return f.showResp, f.showErr
}

// fakeTemplateListUseCase records every List invocation and returns
// the configured response/error. The zero value returns an empty
// catalog without error — matching the LH-FA-TPL-004 empty-state
// path so existing tests that do not wire a template list still
// build (slice-v1-template-list T3).
type fakeTemplateListUseCase struct {
	called  bool
	lastReq driving.TemplateListRequest
	resp    driving.TemplateListResponse
	err     error
}

func (f *fakeTemplateListUseCase) List(_ context.Context, req driving.TemplateListRequest) (driving.TemplateListResponse, error) {
	f.called = true
	f.lastReq = req
	return f.resp, f.err
}

// fakeRemoveServiceUseCase records every Remove invocation and
// returns the configured response/error. Zero value answers a
// generic ErrServiceUnsupported (slice-v1-add-remove T4): a freshly
// constructed instance signals "use case wired but no expectations
// set", which is the default test-helper shape.
type fakeRemoveServiceUseCase struct {
	called  bool
	lastReq driving.RemoveServiceRequest
	resp    driving.RemoveServiceResponse
	err     error
}

func (f *fakeRemoveServiceUseCase) Remove(_ context.Context, req driving.RemoveServiceRequest) (driving.RemoveServiceResponse, error) {
	f.called = true
	f.lastReq = req
	return f.resp, f.err
}

// fakeLogsUseCase records the last LogsRequest and returns the
// configured response/error. Mirrors the shape of the other
// fake use cases (slice-v1-logs T3).
type fakeLogsUseCase struct {
	called  bool
	lastReq driving.LogsRequest
	resp    driving.LogsResponse
	err     error
}

func (f *fakeLogsUseCase) Logs(_ context.Context, req driving.LogsRequest) (driving.LogsResponse, error) {
	f.called = true
	f.lastReq = req
	return f.resp, f.err
}

func newApp(uc driving.InitProjectUseCase, opts ...cli.Option) *cli.App {
	return cli.New("0.0.0-test", uc, &fakeDoctorUseCase{}, &fakeAddServiceUseCase{}, &fakeUpUseCase{}, &fakeDownUseCase{}, &fakeGenerateUseCase{}, &fakeConfigUseCase{}, &fakeTemplateListUseCase{}, &fakeRemoveServiceUseCase{}, &fakeLogsUseCase{}, opts...)
}

// newAppWithDoctor is newApp's variant for doctor-focused tests; the
// caller can wire a fake DoctorUseCase explicitly.
func newAppWithDoctor(uc driving.InitProjectUseCase, doctorUC driving.DoctorUseCase, opts ...cli.Option) *cli.App {
	return cli.New("0.0.0-test", uc, doctorUC, &fakeAddServiceUseCase{}, &fakeUpUseCase{}, &fakeDownUseCase{}, &fakeGenerateUseCase{}, &fakeConfigUseCase{}, &fakeTemplateListUseCase{}, &fakeRemoveServiceUseCase{}, &fakeLogsUseCase{}, opts...)
}

// newAppWithAdd is newApp's variant for add-focused tests.
func newAppWithAdd(uc driving.AddServiceUseCase, opts ...cli.Option) *cli.App {
	return cli.New("0.0.0-test", &fakeInitUseCase{}, &fakeDoctorUseCase{}, uc, &fakeUpUseCase{}, &fakeDownUseCase{}, &fakeGenerateUseCase{}, &fakeConfigUseCase{}, &fakeTemplateListUseCase{}, &fakeRemoveServiceUseCase{}, &fakeLogsUseCase{}, opts...)
}

// newAppWithUp is newApp's variant for `u-boot up`-focused tests.
func newAppWithUp(uc driving.UpUseCase, opts ...cli.Option) *cli.App {
	return cli.New("0.0.0-test", &fakeInitUseCase{}, &fakeDoctorUseCase{}, &fakeAddServiceUseCase{}, uc, &fakeDownUseCase{}, &fakeGenerateUseCase{}, &fakeConfigUseCase{}, &fakeTemplateListUseCase{}, &fakeRemoveServiceUseCase{}, &fakeLogsUseCase{}, opts...)
}

// newAppWithDown is newApp's variant for `u-boot down`-focused tests.
func newAppWithDown(uc driving.DownUseCase, opts ...cli.Option) *cli.App {
	return cli.New("0.0.0-test", &fakeInitUseCase{}, &fakeDoctorUseCase{}, &fakeAddServiceUseCase{}, &fakeUpUseCase{}, uc, &fakeGenerateUseCase{}, &fakeConfigUseCase{}, &fakeTemplateListUseCase{}, &fakeRemoveServiceUseCase{}, &fakeLogsUseCase{}, opts...)
}

// newAppWithGenerate is newApp's variant for `u-boot generate`-focused
// tests.
func newAppWithGenerate(uc driving.GenerateUseCase, opts ...cli.Option) *cli.App {
	return cli.New("0.0.0-test", &fakeInitUseCase{}, &fakeDoctorUseCase{}, &fakeAddServiceUseCase{}, &fakeUpUseCase{}, &fakeDownUseCase{}, uc, &fakeConfigUseCase{}, &fakeTemplateListUseCase{}, &fakeRemoveServiceUseCase{}, &fakeLogsUseCase{}, opts...)
}

// newAppWithConfig is newApp's variant for `u-boot config`-focused
// tests.
func newAppWithConfig(uc driving.ConfigUseCase, opts ...cli.Option) *cli.App {
	return cli.New("0.0.0-test", &fakeInitUseCase{}, &fakeDoctorUseCase{}, &fakeAddServiceUseCase{}, &fakeUpUseCase{}, &fakeDownUseCase{}, &fakeGenerateUseCase{}, uc, &fakeTemplateListUseCase{}, &fakeRemoveServiceUseCase{}, &fakeLogsUseCase{}, opts...)
}

// newAppWithTemplateList is newApp's variant for
// `u-boot template list`-focused tests (slice-v1-template-list T3).
func newAppWithTemplateList(uc driving.TemplateListUseCase, opts ...cli.Option) *cli.App {
	return cli.New("0.0.0-test", &fakeInitUseCase{}, &fakeDoctorUseCase{}, &fakeAddServiceUseCase{}, &fakeUpUseCase{}, &fakeDownUseCase{}, &fakeGenerateUseCase{}, &fakeConfigUseCase{}, uc, &fakeRemoveServiceUseCase{}, &fakeLogsUseCase{}, opts...)
}

// newAppWithRemove is newApp's variant for `u-boot remove`-focused
// tests (slice-v1-add-remove T4).
func newAppWithRemove(uc driving.RemoveServiceUseCase, opts ...cli.Option) *cli.App {
	return cli.New("0.0.0-test", &fakeInitUseCase{}, &fakeDoctorUseCase{}, &fakeAddServiceUseCase{}, &fakeUpUseCase{}, &fakeDownUseCase{}, &fakeGenerateUseCase{}, &fakeConfigUseCase{}, &fakeTemplateListUseCase{}, uc, &fakeLogsUseCase{}, opts...)
}

// newAppWithLogs is newApp's variant for `u-boot logs`-focused
// tests (slice-v1-logs T3). Mirrors the shape of the other
// "with"-helpers.
func newAppWithLogs(uc driving.LogsUseCase, opts ...cli.Option) *cli.App {
	return cli.New("0.0.0-test", &fakeInitUseCase{}, &fakeDoctorUseCase{}, &fakeAddServiceUseCase{}, &fakeUpUseCase{}, &fakeDownUseCase{}, &fakeGenerateUseCase{}, &fakeConfigUseCase{}, &fakeTemplateListUseCase{}, &fakeRemoveServiceUseCase{}, uc, opts...)
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

// TestExecute_Init_DevcontainerFlag_Propagates pins the MVP-Closure
// T1 flag wiring: `u-boot init --devcontainer` sets the
// InitProjectRequest.Devcontainer field that the application
// service branches on. Mirrors the --no-git propagation test.
func TestExecute_Init_DevcontainerFlag_Propagates(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeInitUseCase{
		resp: driving.InitProjectResponse{
			Project: domain.NewProject(mustProjectName(t, "demo")),
		},
	}
	var stdout, stderr bytes.Buffer
	err := newApp(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"init", "--devcontainer"}, &stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute init --devcontainer: %v", err)
	}
	if !uc.lastReq.Devcontainer {
		t.Errorf("init Devcontainer = false, want true (--devcontainer was passed)")
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
		// M7-T6: generate sentinels. ErrArtifactUnknown lives in
		// isUsageError (code 2) by spec mandate (§LH-FA-GEN-001) —
		// distinct from `add <unknown-service>` which maps to 10.
		// ErrGenerateManualConflict joins the code-10 cohort;
		// ErrGenerateFileSystem joins the code-14 isFilesystemError
		// list (slice plan T6 DoD pin).
		{"ErrArtifactUnknown (usage)", driving.ErrArtifactUnknown, 2},
		{"wrapped ErrArtifactUnknown", fmt.Errorf("ctx: %w", driving.ErrArtifactUnknown), 2},
		{"ErrGenerateManualConflict (validation)", driving.ErrGenerateManualConflict, 10},
		{"wrapped ErrGenerateManualConflict", fmt.Errorf("ctx: %w", driving.ErrGenerateManualConflict), 10},
		{"ErrGenerateFileSystem (fs)", driving.ErrGenerateFileSystem, 14},
		{"wrapped ErrGenerateFileSystem", fmt.Errorf("ctx: %w", driving.ErrGenerateFileSystem), 14},
		// M8-T5: config sentinels.
		{"ErrConfigPathUnknown (validation)", driving.ErrConfigPathUnknown, 10},
		{"ErrConfigValueInvalid (validation)", driving.ErrConfigValueInvalid, 10},
		{"ErrConfigSchemaInvalid (validation)", driving.ErrConfigSchemaInvalid, 10},
		{"ErrConfigValueNotSet (validation)", driving.ErrConfigValueNotSet, 10},
		{"wrapped ErrConfigPathUnknown", fmt.Errorf("ctx: %w", driving.ErrConfigPathUnknown), 10},
		{"ErrConfigFileSystem (fs)", driving.ErrConfigFileSystem, 14},
		{"wrapped ErrConfigFileSystem", fmt.Errorf("ctx: %w", driving.ErrConfigFileSystem), 14},
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

func TestExecute_Add_WithDepsFlag_PlumbedToRequest(t *testing.T) {
	// LH-FA-ADD-006: pin that --with-deps reaches AddServiceRequest.
	// The use case sees req.WithDeps == true; --yes / --no-interactive
	// stay false in this scenario.
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeAddServiceUseCase{
		resp: driving.AddServiceResponse{
			ServiceName: mustServiceName(t, "postgres"),
			PriorState:  domain.ServiceStateUnregistered,
			State:       domain.ServiceStateActive,
			Changed:     []string{"u-boot.yaml"},
		},
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithAdd(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"add", "--with-deps", "postgres"}, &stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !uc.lastReq.WithDeps {
		t.Errorf("req.WithDeps = false, want true (--with-deps plumbed through)")
	}
	if uc.lastReq.Yes {
		t.Errorf("req.Yes = true, want false (no --yes given)")
	}
	if uc.lastReq.NoInteractive {
		t.Errorf("req.NoInteractive = true, want false (no --no-interactive given)")
	}
}

func TestExecute_Add_YesFlag_PlumbedToRequest(t *testing.T) {
	// Pin that the root --yes flag reaches AddServiceRequest.Yes so
	// the application layer's four-mode dispatch can promote it to
	// autoInstall (LH-FA-ADD-006).
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeAddServiceUseCase{
		resp: driving.AddServiceResponse{
			ServiceName: mustServiceName(t, "postgres"),
			PriorState:  domain.ServiceStateUnregistered,
			State:       domain.ServiceStateActive,
			Changed:     []string{"u-boot.yaml"},
		},
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithAdd(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"--yes", "add", "postgres"}, &stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !uc.lastReq.Yes {
		t.Errorf("req.Yes = false, want true (--yes plumbed)")
	}
	if uc.lastReq.WithDeps {
		t.Errorf("req.WithDeps = true, want false")
	}
}

func TestExecute_Add_NoInteractiveFlag_PlumbedToRequest(t *testing.T) {
	// Pin that --no-interactive reaches AddServiceRequest so the
	// dispatch can take the fail-fast arm when deps are missing.
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeAddServiceUseCase{
		resp: driving.AddServiceResponse{
			ServiceName: mustServiceName(t, "postgres"),
			PriorState:  domain.ServiceStateUnregistered,
			State:       domain.ServiceStateActive,
			Changed:     []string{"u-boot.yaml"},
		},
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithAdd(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"--no-interactive", "add", "postgres"}, &stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !uc.lastReq.NoInteractive {
		t.Errorf("req.NoInteractive = false, want true")
	}
}

func TestExecute_Add_DependenciesRequired_Code10(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeAddServiceUseCase{
		err: fmt.Errorf("%w: %q requires [postgres] which is/are not registered — add them first or rerun with --with-deps",
			driving.ErrDependenciesRequired, "keycloak"),
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithAdd(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"add", "keycloak"}, &stdout, &stderr,
	)
	if !errors.Is(err, driving.ErrDependenciesRequired) {
		t.Fatalf("err = %v, want wrap of ErrDependenciesRequired", err)
	}
	if got := cli.ExitCode(err); got != 10 {
		t.Errorf("ExitCode = %d, want 10", got)
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

func TestExecute_Up_Quiet_SuppressesStatusTableAndDiagnostics(t *testing.T) {
	t.Parallel()
	// Slice §T6 binding contract: `up --quiet` must suppress BOTH
	// the status table AND the diagnostic section. T7-review (post-
	// 6d9aa88) found that the initial implementation only
	// suppressed diagnostics; this pin makes the regression
	// impossible.
	uc := &fakeUpUseCase{
		resp: driving.UpResponse{Result: domain.UpResult{
			Stabilized: true,
			Services: []domain.ServiceStatus{
				{Name: "postgres", ContainerStatus: domain.StateRunning, Port: "5432:5432", Healthcheck: "healthy"},
			},
			Diagnostics: []domain.Diagnostic{
				{ID: "up.port.x.0", Severity: domain.SeverityWarn, Message: "would-show-without-quiet"},
			},
		}},
	}
	getwd := func() (string, error) { return "/tmp/proj", nil }
	var stdout, stderr bytes.Buffer
	err := newAppWithUp(uc, cli.WithGetwd(getwd)).Execute(context.Background(), []string{"--quiet", "up"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("up --quiet: %v", err)
	}
	if stdout.Len() != 0 {
		t.Errorf("--quiet should produce empty stdout (no table, no diagnostics), got: %q", stdout.String())
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

func TestExecute_Down_RootYesBeforeSubcmd_BypassesConfirm(t *testing.T) {
	t.Parallel()
	// M6-closure-review fix #2: spec §237 lists `u-boot down
	// --volumes` among the commands governed by the persistent
	// --yes / --no-interactive root flags. Pin that `u-boot --yes
	// down --volumes` (root flag BEFORE subcommand) behaves
	// identically to `u-boot down --volumes --yes` (root flag
	// AFTER subcommand) — Cobra propagates persistent flags either
	// way, and the down subcommand must read a.yes, not a local
	// --yes flag.
	uc := &fakeDownUseCase{resp: driving.DownResponse{RemovedVolumes: true}}
	getwd := func() (string, error) { return "/tmp/proj", nil }
	var stdout, stderr bytes.Buffer
	err := newAppWithDown(uc, cli.WithGetwd(getwd)).Execute(context.Background(), []string{"--yes", "down", "--volumes"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("--yes down --volumes: %v", err)
	}
	if !uc.lastReq.AssumeYes {
		t.Errorf("AssumeYes = false; root --yes did not flow through to DownRequest")
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

// ----------------------------------------------------------------------------
// `u-boot generate` subcommand (M7-T6)
// ----------------------------------------------------------------------------

// TestExitCode_GenerateFileSystemError_MapsTo14 is the explicit pin
// the slice plan T6 DoD asks for: ErrGenerateFileSystem joins the
// isFilesystemError list and surfaces as exit code 14. The shared
// TestExitCode_BaseMappings table already covers this case, but a
// stand-alone test keeps the spec-anchored intent visible by name
// so future refactors do not silently drop the mapping.
func TestExitCode_GenerateFileSystemError_MapsTo14(t *testing.T) {
	t.Parallel()
	if got := cli.ExitCode(driving.ErrGenerateFileSystem); got != 14 {
		t.Errorf("ExitCode(ErrGenerateFileSystem) = %d, want 14", got)
	}
	wrapped := fmt.Errorf("write %q: %w", ".env.example", driving.ErrGenerateFileSystem)
	if got := cli.ExitCode(wrapped); got != 14 {
		t.Errorf("ExitCode(wrap(ErrGenerateFileSystem)) = %d, want 14", got)
	}
}

func TestExecute_Generate_HappyPathCreated(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeGenerateUseCase{
		resp: driving.GenerateResponse{
			Artifact: domain.ArtifactEnvExample,
			Action:   driving.GenerateActionCreated,
			Changed:  []string{".env.example"},
		},
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithGenerate(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"generate", "env-example"}, &stdout, &stderr,
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
	if uc.lastReq.Artifact != domain.ArtifactEnvExample {
		t.Errorf("Artifact = %v, want ArtifactEnvExample", uc.lastReq.Artifact)
	}
	out := stdout.String()
	if !strings.Contains(out, "Generated env-example (.env.example).") {
		t.Errorf("missing Created summary; got:\n%s", out)
	}
}

func TestExecute_Generate_UpdatedBlockSummary(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeGenerateUseCase{
		resp: driving.GenerateResponse{
			Artifact: domain.ArtifactReadme,
			Action:   driving.GenerateActionUpdatedBlock,
			Changed:  []string{"README.md"},
		},
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithGenerate(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"generate", "readme"}, &stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !strings.Contains(stdout.String(), "Updated readme managed block (README.md).") {
		t.Errorf("missing UpdatedBlock summary; got:\n%s", stdout.String())
	}
}

func TestExecute_Generate_NoOpSummary(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeGenerateUseCase{
		resp: driving.GenerateResponse{
			Artifact: domain.ArtifactChangelog,
			Action:   driving.GenerateActionNoOp,
		},
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithGenerate(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"generate", "changelog"}, &stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !strings.Contains(stdout.String(), "changelog already up to date; no changes.") {
		t.Errorf("missing NoOp summary; got:\n%s", stdout.String())
	}
}

func TestExecute_Generate_RepairedManualSummary(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeGenerateUseCase{
		resp: driving.GenerateResponse{
			Artifact: domain.ArtifactChangelog,
			Action:   driving.GenerateActionRepairedManual,
			Changed:  []string{"CHANGELOG.md"},
		},
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithGenerate(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"generate", "changelog"}, &stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !strings.Contains(stdout.String(), "Repaired changelog structure (CHANGELOG.md).") {
		t.Errorf("missing RepairedManual summary; got:\n%s", stdout.String())
	}
}

func TestExecute_Generate_UnknownArtifact_Code2(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeGenerateUseCase{}
	var stdout, stderr bytes.Buffer
	err := newAppWithGenerate(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"generate", "dockerfile"}, &stdout, &stderr,
	)
	if err == nil {
		t.Fatalf("expected error for unknown artifact")
	}
	if !errors.Is(err, driving.ErrArtifactUnknown) {
		t.Errorf("err does not wrap ErrArtifactUnknown: %v", err)
	}
	if got := cli.ExitCode(err); got != 2 {
		t.Errorf("ExitCode = %d, want 2 (LH-FA-GEN-001 mandates code 2 for unknown artefact)", got)
	}
	if uc.called {
		t.Errorf("use-case should not have been called on validation failure")
	}
}

func TestExecute_Generate_NoArgs_Code2(t *testing.T) {
	uc := &fakeGenerateUseCase{}
	var stdout, stderr bytes.Buffer
	err := newAppWithGenerate(uc).Execute(
		context.Background(), []string{"generate"}, &stdout, &stderr,
	)
	if err == nil {
		t.Fatalf("expected error for missing positional argument")
	}
	if got := cli.ExitCode(err); got != 2 {
		t.Errorf("ExitCode = %d, want 2 (Cobra ExactArgs(1) miss)", got)
	}
}

func TestExecute_Generate_ProjectNotInitialized_Code10(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeGenerateUseCase{
		err: fmt.Errorf("u-boot.yaml missing: %w", driving.ErrProjectNotInitialized),
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithGenerate(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"generate", "env-example"}, &stdout, &stderr,
	)
	if err == nil {
		t.Fatalf("expected error from use-case")
	}
	if got := cli.ExitCode(err); got != 10 {
		t.Errorf("ExitCode = %d, want 10", got)
	}
}

func TestExecute_Generate_ManualConflict_Code10(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeGenerateUseCase{
		err: fmt.Errorf("no init block: %w", driving.ErrGenerateManualConflict),
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithGenerate(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"generate", "readme"}, &stdout, &stderr,
	)
	if err == nil {
		t.Fatalf("expected error from use-case")
	}
	if got := cli.ExitCode(err); got != 10 {
		t.Errorf("ExitCode = %d, want 10", got)
	}
}

func TestExecute_Generate_HelpListsFourArtifacts(t *testing.T) {
	uc := &fakeGenerateUseCase{}
	var stdout, stderr bytes.Buffer
	err := newAppWithGenerate(uc).Execute(
		context.Background(), []string{"generate", "--help"}, &stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute --help: %v", err)
	}
	out := stdout.String()
	for _, want := range []string{"changelog", "readme", "env-example", "devcontainer"} {
		if !strings.Contains(out, want) {
			t.Errorf("--help missing artefact %q; got:\n%s", want, out)
		}
	}
}

// ----------------------------------------------------------------------------
// `u-boot config` subcommand (M8-T5)
// ----------------------------------------------------------------------------

// TestExecute_Config_Show_PrintsBodyByteIdentical pins the §D5
// contract: Show writes the use-case Body to stdout without
// re-formatting.
func TestExecute_Config_Show_PrintsBodyByteIdentical(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	body := []byte("# u-boot project config\nschemaVersion: 1\nproject:\n  name: demo  # display\n")
	uc := &fakeConfigUseCase{showResp: driving.ConfigShowResponse{Body: body}}
	var stdout, stderr bytes.Buffer
	err := newAppWithConfig(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"config"}, &stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !uc.showCalled {
		t.Fatal("Show was not invoked")
	}
	if uc.showReq.BaseDir != "/tmp/x/demo" {
		t.Errorf("BaseDir = %q, want /tmp/x/demo", uc.showReq.BaseDir)
	}
	if !bytes.Equal(stdout.Bytes(), body) {
		t.Errorf("stdout differs from Body:\n got:%q\nwant:%q", stdout.Bytes(), body)
	}
}

// TestExecute_Config_Get_PrintsScalarWithTrailingNewline pins §D4:
// bare scalar + exactly one trailing newline.
func TestExecute_Config_Get_PrintsScalarWithTrailingNewline(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeConfigUseCase{getResp: driving.ConfigGetResponse{Value: "demo"}}
	var stdout, stderr bytes.Buffer
	err := newAppWithConfig(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"config", "get", "project.name"}, &stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if got, want := stdout.String(), "demo\n"; got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
	if uc.getReq.Path.String() != "project.name" {
		t.Errorf("Path = %v, want project.name", uc.getReq.Path)
	}
}

func TestExecute_Config_Set_ChangedSummary(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeConfigUseCase{
		setResp: driving.ConfigSetResponse{
			Path:     mustConfigPathInTest(t, "project.name"),
			OldValue: "demo",
			NewValue: "renamed",
		},
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithConfig(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"config", "set", "project.name", "renamed"}, &stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !strings.Contains(stdout.String(), "config: project.name demo → renamed.") {
		t.Errorf("changed-summary missing; got:\n%s", stdout.String())
	}
}

func TestExecute_Config_Set_NoOpSummary(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeConfigUseCase{
		setResp: driving.ConfigSetResponse{
			Path:     mustConfigPathInTest(t, "project.name"),
			OldValue: "demo",
			NewValue: "demo",
		},
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithConfig(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"config", "set", "project.name", "demo"}, &stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !strings.Contains(stdout.String(), "already demo; no changes.") {
		t.Errorf("noop-summary missing; got:\n%s", stdout.String())
	}
}

// TestExecute_Config_Set_UnsetOld_RendersUnsetMarker pins the
// summary-line edge case: first-time write of an optional field
// where the empty OldValue would otherwise render as a bare
// arrow with no left side.
func TestExecute_Config_Set_UnsetOld_RendersUnsetMarker(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeConfigUseCase{
		setResp: driving.ConfigSetResponse{
			Path:     mustConfigPathInTest(t, "devcontainer.enabled"),
			OldValue: "",
			NewValue: "true",
		},
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithConfig(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"config", "set", "devcontainer.enabled", "true"}, &stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !strings.Contains(stdout.String(), "config: devcontainer.enabled (unset) → true.") {
		t.Errorf("unset-marker missing; got:\n%s", stdout.String())
	}
}

func TestExecute_Config_Get_UnknownPath_Code10(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeConfigUseCase{}
	var stdout, stderr bytes.Buffer
	err := newAppWithConfig(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"config", "get", "totally.unknown"}, &stdout, &stderr,
	)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, driving.ErrConfigPathUnknown) {
		t.Errorf("err = %v, want wrap of ErrConfigPathUnknown", err)
	}
	if got := cli.ExitCode(err); got != 10 {
		t.Errorf("ExitCode = %d, want 10", got)
	}
	if uc.getCalled {
		t.Errorf("use case should not have been called on path-validation failure")
	}
}

func TestExecute_Config_Set_TooFewArgs_Code2(t *testing.T) {
	uc := &fakeConfigUseCase{}
	var stdout, stderr bytes.Buffer
	err := newAppWithConfig(uc).Execute(
		context.Background(), []string{"config", "set", "project.name"}, &stdout, &stderr,
	)
	if err == nil {
		t.Fatalf("expected error for missing value arg")
	}
	if got := cli.ExitCode(err); got != 2 {
		t.Errorf("ExitCode = %d, want 2 (Cobra ExactArgs(2) miss)", got)
	}
}

func TestExecute_Config_Get_NoArgs_Code2(t *testing.T) {
	uc := &fakeConfigUseCase{}
	var stdout, stderr bytes.Buffer
	err := newAppWithConfig(uc).Execute(
		context.Background(), []string{"config", "get"}, &stdout, &stderr,
	)
	if err == nil {
		t.Fatalf("expected error for missing path arg")
	}
	if got := cli.ExitCode(err); got != 2 {
		t.Errorf("ExitCode = %d, want 2 (Cobra ExactArgs(1) miss)", got)
	}
}

func TestExecute_Config_UseCaseError_Code10(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeConfigUseCase{
		getErr: fmt.Errorf("u-boot.yaml missing: %w", driving.ErrProjectNotInitialized),
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithConfig(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"config", "get", "project.name"}, &stdout, &stderr,
	)
	if err == nil {
		t.Fatalf("expected error from use case")
	}
	if got := cli.ExitCode(err); got != 10 {
		t.Errorf("ExitCode = %d, want 10", got)
	}
}

// mustConfigPathInTest mirrors the application-test-package helper
// of the same name (lives in config_test.go there). The CLI tests
// can not reach that package's helper, so this is a local copy.
func mustConfigPathInTest(t *testing.T, raw string) domain.ConfigPath {
	t.Helper()
	p, err := domain.NewConfigPath(raw)
	if err != nil {
		t.Fatalf("NewConfigPath(%q): %v", raw, err)
	}
	return p
}

// TestExecute_Init_AllowExternalWithoutDevcontainer_Code10 pins the
// slice-v1-devcontainer-features Review-Followup R1 fix: the
// LH-FA-DEV-003 `--allow-external-feature-sources requires
// --devcontainer` rejection must map to exit-code 10 per Spec §720,
// not the default-1 fallback. The sentinel was moved from
// `application` to `domain` so this adapter could include it in
// [cli.isValidationError]; the test pins the wiring.
func TestExecute_Init_AllowExternalWithoutDevcontainer_Code10(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeInitUseCase{
		err: fmt.Errorf("--allow-external-feature-sources requires --devcontainer (Spec §714): %w", domain.ErrInvalidFeatureSource),
	}
	var stdout, stderr bytes.Buffer
	err := newApp(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(),
		[]string{"init", "--allow-external-feature-sources", "https://example.test/x"},
		&stdout, &stderr,
	)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrInvalidFeatureSource) {
		t.Errorf("err = %v, want wrap of domain.ErrInvalidFeatureSource", err)
	}
	if got := cli.ExitCode(err); got != 10 {
		t.Errorf("ExitCode = %d, want 10 (Spec §720)", got)
	}
}

// TestExecute_Generate_InvalidAllowFlagURL_Code10 pins that a
// malformed URL on `generate devcontainer --allow-external-feature-
// sources <bad>` rejects with exit-code 10 (Spec §1353).
func TestExecute_Generate_InvalidAllowFlagURL_Code10(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeGenerateUseCase{
		err: fmt.Errorf("generate devcontainer: --allow-external-feature-sources: %w", domain.ErrInvalidFeatureSource),
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithGenerate(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(),
		[]string{"generate", "devcontainer", "--allow-external-feature-sources", "not-a-url"},
		&stdout, &stderr,
	)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrInvalidFeatureSource) {
		t.Errorf("err = %v, want wrap of domain.ErrInvalidFeatureSource", err)
	}
	if got := cli.ExitCode(err); got != 10 {
		t.Errorf("ExitCode = %d, want 10 (Spec §1353)", got)
	}
}

// --- slice-v1-logs T3: CLI Logs-Subkommando -----------------------

func TestExecute_Logs_NoArgs_PropagatesEmptyService(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeLogsUseCase{}
	var stdout, stderr bytes.Buffer
	err := newAppWithLogs(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"logs"}, &stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute logs: %v", err)
	}
	if !uc.called {
		t.Fatal("LogsUseCase.Logs not called")
	}
	if uc.lastReq.Service != "" {
		t.Errorf("Service = %q, want empty (T0-(a) Compose-Default)", uc.lastReq.Service)
	}
	if uc.lastReq.Follow {
		t.Errorf("Follow = true, want false (default)")
	}
	if uc.lastReq.Tail != "" {
		t.Errorf("Tail = %q, want empty (default → Use-Case normalises to \"all\")", uc.lastReq.Tail)
	}
	if uc.lastReq.OutputSink == nil {
		t.Errorf("OutputSink is nil; want cmd.OutOrStdout() forwarded")
	}
}

func TestExecute_Logs_ServiceArg_PropagatesService(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeLogsUseCase{}
	var stdout, stderr bytes.Buffer
	err := newAppWithLogs(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"logs", "postgres"}, &stdout, &stderr,
	)
	if err != nil {
		t.Fatalf("Execute logs postgres: %v", err)
	}
	if uc.lastReq.Service != "postgres" {
		t.Errorf("Service = %q, want \"postgres\"", uc.lastReq.Service)
	}
}

func TestExecute_Logs_FollowFlag_PropagatesTrue(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeLogsUseCase{}
	var stdout, stderr bytes.Buffer
	if err := newAppWithLogs(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"logs", "--follow"}, &stdout, &stderr,
	); err != nil {
		t.Fatalf("Execute logs --follow: %v", err)
	}
	if !uc.lastReq.Follow {
		t.Errorf("Follow = false, want true")
	}
}

func TestExecute_Logs_TailFlag_PropagatesValue(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeLogsUseCase{}
	var stdout, stderr bytes.Buffer
	if err := newAppWithLogs(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"logs", "--tail", "100"}, &stdout, &stderr,
	); err != nil {
		t.Fatalf("Execute logs --tail 100: %v", err)
	}
	if uc.lastReq.Tail != "100" {
		t.Errorf("Tail = %q, want \"100\"", uc.lastReq.Tail)
	}
}

// TestExecute_Logs_InvalidTail_Negative_Code2 pins T0-(c) + §AK:
// negative integers map to Exit-Code 2 via ErrInvalidLogsTail —
// the Use-Case is never called.
func TestExecute_Logs_InvalidTail_Negative_Code2(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeLogsUseCase{}
	var stdout, stderr bytes.Buffer
	err := newAppWithLogs(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"logs", "--tail", "-1"}, &stdout, &stderr,
	)
	if err == nil {
		t.Fatalf("expected error for negative --tail, got nil")
	}
	if !errors.Is(err, cli.ErrInvalidLogsTail) {
		t.Errorf("err = %v, want wrap of ErrInvalidLogsTail", err)
	}
	if got := cli.ExitCode(err); got != 2 {
		t.Errorf("ExitCode = %d, want 2 (Stage-1-Validation)", got)
	}
	if uc.called {
		t.Error("LogsUseCase.Logs called despite invalid --tail; Stage-1 should reject before dispatch")
	}
}

func TestExecute_Logs_InvalidTail_NonNumeric_Code2(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeLogsUseCase{}
	var stdout, stderr bytes.Buffer
	err := newAppWithLogs(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"logs", "--tail", "all"}, &stdout, &stderr,
	)
	if err == nil {
		t.Fatalf("expected error for non-numeric --tail (\"all\" is internal, not user-input)")
	}
	if !errors.Is(err, cli.ErrInvalidLogsTail) {
		t.Errorf("err = %v, want wrap of ErrInvalidLogsTail", err)
	}
	if got := cli.ExitCode(err); got != 2 {
		t.Errorf("ExitCode = %d, want 2", got)
	}
}

// TestExecute_Logs_InvalidServiceName_Code10 pins T0-(b):
// regex-only validation in the CLI; format failure → Exit-10
// via isServiceValidationError. The "all" reserved-word would
// pass the regex (lowercase, valid characters), so we use an
// uppercase name instead which definitely violates the regex.
func TestExecute_Logs_InvalidServiceName_Code10(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeLogsUseCase{}
	var stdout, stderr bytes.Buffer
	err := newAppWithLogs(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"logs", "Postgres"}, &stdout, &stderr,
	)
	if err == nil {
		t.Fatalf("expected error for invalid service-name format")
	}
	if !errors.Is(err, domain.ErrInvalidServiceName) {
		t.Errorf("err = %v, want wrap of domain.ErrInvalidServiceName", err)
	}
	if got := cli.ExitCode(err); got != 10 {
		t.Errorf("ExitCode = %d, want 10", got)
	}
	if uc.called {
		t.Error("LogsUseCase.Logs called despite invalid service-name; Stage-1 should reject")
	}
}

// TestExecute_Logs_TooManyArgs_Code2 pins MaximumNArgs(1) — Cobra
// rejects 2+ positional arguments with a usage-error → Exit-2.
func TestExecute_Logs_TooManyArgs_Code2(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeLogsUseCase{}
	var stdout, stderr bytes.Buffer
	err := newAppWithLogs(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"logs", "postgres", "keycloak"}, &stdout, &stderr,
	)
	if err == nil {
		t.Fatalf("expected error for two positional args")
	}
	if got := cli.ExitCode(err); got != 2 {
		t.Errorf("ExitCode = %d, want 2 (Cobra usage error)", got)
	}
}

// TestExecute_Logs_UseCaseError_ComposeRuntime_Code12 pins the
// 12-exit-code path from the application service through the
// CLI mapping.
func TestExecute_Logs_UseCaseError_ComposeRuntime_Code12(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeLogsUseCase{
		err: fmt.Errorf("logs service: ComposeLogs: %w", driven.ErrComposeRuntime),
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithLogs(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"logs"}, &stdout, &stderr,
	)
	if err == nil {
		t.Fatal("expected error from use case")
	}
	if got := cli.ExitCode(err); got != 12 {
		t.Errorf("ExitCode = %d, want 12 (ErrComposeRuntime)", got)
	}
}

// TestExecute_Logs_UseCaseError_DockerUnavailable_Code11 pins the
// 11-exit-code path.
func TestExecute_Logs_UseCaseError_DockerUnavailable_Code11(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeLogsUseCase{
		err: fmt.Errorf("logs service: ComposeLogs: %w", driven.ErrDockerUnavailable),
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithLogs(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"logs"}, &stdout, &stderr,
	)
	if err == nil {
		t.Fatal("expected error from use case")
	}
	if got := cli.ExitCode(err); got != 11 {
		t.Errorf("ExitCode = %d, want 11 (ErrDockerUnavailable)", got)
	}
}

// TestExecute_Logs_UseCaseError_ProjectNotInitialized_Code10 pins
// the 10-exit-code path for missing u-boot.yaml.
func TestExecute_Logs_UseCaseError_ProjectNotInitialized_Code10(t *testing.T) {
	getwd := func() (string, error) { return "/tmp/x/demo", nil }
	uc := &fakeLogsUseCase{
		err: fmt.Errorf("logs service: u-boot.yaml absent: %w", driving.ErrProjectNotInitialized),
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithLogs(uc, cli.WithGetwd(getwd)).Execute(
		context.Background(), []string{"logs"}, &stdout, &stderr,
	)
	if err == nil {
		t.Fatal("expected error from use case")
	}
	if got := cli.ExitCode(err); got != 10 {
		t.Errorf("ExitCode = %d, want 10 (ErrProjectNotInitialized)", got)
	}
}
