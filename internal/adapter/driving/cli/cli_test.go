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

func TestExitCode_Mapping(t *testing.T) {
	if cli.ExitCode(nil) != 0 {
		t.Errorf("ExitCode(nil) = %d, want 0", cli.ExitCode(nil))
	}
	if cli.ExitCode(errors.New("boom")) != 1 {
		t.Errorf("ExitCode(generic) = %d, want 1", cli.ExitCode(errors.New("boom")))
	}
	if cli.ExitCode(driving.ErrProjectExists) != 10 {
		t.Errorf("ExitCode(ErrProjectExists) = %d, want 10", cli.ExitCode(driving.ErrProjectExists))
	}
	// Pinned usage-error prefixes (see isUsageError):
	for _, prefix := range []string{
		"unknown command",
		"unknown flag",
		"flag needs an argument",
		"invalid argument",
		"requires at least 1 arg(s)",
		"accepts 0 arg(s)",
	} {
		err := errors.New(prefix + " — sample")
		if got := cli.ExitCode(err); got != 2 {
			t.Errorf("ExitCode(%q) = %d, want 2", prefix, got)
		}
	}
}
