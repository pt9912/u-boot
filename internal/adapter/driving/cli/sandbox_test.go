package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

func sandboxGetwd() (string, error) { return "/tmp/x/demo", nil }

// LH-FA-DEV-006: `init --devcontainer --sandbox` is forwarded to the
// use case.
func TestExecute_Init_SandboxFlagForwarded(t *testing.T) {
	uc := &fakeInitUseCase{}
	var stdout, stderr bytes.Buffer
	err := newApp(uc, cli.WithGetwd(sandboxGetwd)).Execute(context.Background(),
		[]string{"init", "--devcontainer", "--sandbox"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !uc.lastReq.Sandbox || !uc.lastReq.Devcontainer {
		t.Errorf("request = %+v, want Sandbox && Devcontainer", uc.lastReq)
	}
}

// `--sandbox` without `--devcontainer` (use-case error) → exit 10.
func TestExecute_Init_SandboxWithoutDevcontainer_Code10(t *testing.T) {
	uc := &fakeInitUseCase{err: fmt.Errorf("%w: --sandbox requires --devcontainer", domain.ErrInvalidSandboxSetting)}
	var stdout, stderr bytes.Buffer
	err := newApp(uc, cli.WithGetwd(sandboxGetwd)).Execute(context.Background(),
		[]string{"init", "--sandbox"}, &stdout, &stderr)
	if got := cli.ExitCode(err); got != 10 {
		t.Errorf("ExitCode = %d, want 10 (err=%v)", got, err)
	}
}

// Warnings (no remote) surface in human and JSON output, exit 0.
func TestExecute_Init_SandboxWarnings(t *testing.T) {
	warn := driving.WarningEntry{Code: "LH-FA-DEV-006", Level: "warn", Message: "no git remote 'origin' found"}
	uc := &fakeInitUseCase{resp: driving.InitProjectResponse{Warnings: []driving.WarningEntry{warn}}}

	var stdout, stderr bytes.Buffer
	if err := newApp(uc, cli.WithGetwd(sandboxGetwd)).Execute(context.Background(),
		[]string{"init", "--devcontainer", "--sandbox"}, &stdout, &stderr); err != nil {
		t.Fatalf("human: %v", err)
	}
	if !strings.Contains(stdout.String(), "Warning (LH-FA-DEV-006): no git remote") {
		t.Errorf("human output lacks warning:\n%s", stdout.String())
	}

	stdout.Reset()
	if err := newApp(uc, cli.WithGetwd(sandboxGetwd)).Execute(context.Background(),
		[]string{"--json", "init", "--devcontainer", "--sandbox"}, &stdout, &stderr); err != nil {
		t.Fatalf("json: %v", err)
	}
	var env struct {
		Status      string `json:"status"`
		Diagnostics []struct{ Level, Code string }
	}
	if err := json.Unmarshal(stdout.Bytes(), &env); err != nil {
		t.Fatalf("envelope: %v\n%s", err, stdout.String())
	}
	if len(env.Diagnostics) != 1 || env.Diagnostics[0].Level != "warn" || env.Diagnostics[0].Code != "LH-FA-DEV-006" {
		t.Errorf("diagnostics = %+v, want one warn LH-FA-DEV-006", env.Diagnostics)
	}
}

func TestExecute_Generate_SandboxFlagForwardedAndWarns(t *testing.T) {
	warn := driving.WarningEntry{Code: "LH-FA-DEV-006", Level: "warn", Message: "no git remote 'origin' found"}
	uc := &fakeGenerateUseCase{resp: driving.GenerateResponse{
		Artifact: domain.ArtifactDevcontainer, Action: driving.GenerateActionCreated,
		Changed: []string{".devcontainer/devcontainer.json"}, Warnings: []driving.WarningEntry{warn}}}
	var stdout, stderr bytes.Buffer
	if err := newAppWithGenerate(uc, cli.WithGetwd(sandboxGetwd)).Execute(context.Background(),
		[]string{"generate", "devcontainer", "--sandbox"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !uc.lastReq.Sandbox {
		t.Errorf("Sandbox not forwarded: %+v", uc.lastReq)
	}
	if !strings.Contains(stdout.String(), "Warning (LH-FA-DEV-006)") {
		t.Errorf("output lacks warning:\n%s", stdout.String())
	}

	stdout.Reset()
	if err := newAppWithGenerate(uc, cli.WithGetwd(sandboxGetwd)).Execute(context.Background(),
		[]string{"--json", "generate", "devcontainer", "--sandbox"}, &stdout, &stderr); err != nil {
		t.Fatalf("json: %v", err)
	}
	if !strings.Contains(stdout.String(), `"code":"LH-FA-DEV-006"`) && !strings.Contains(stdout.String(), `"code": "LH-FA-DEV-006"`) {
		t.Errorf("JSON lacks warn diagnostic:\n%s", stdout.String())
	}
}

// `generate <other> --sandbox` is a usage error (exit 2), like the
// allow-external flag; the use case is not called.
func TestExecute_Generate_SandboxOnOtherArtifact_Code2(t *testing.T) {
	uc := &fakeGenerateUseCase{}
	var stdout, stderr bytes.Buffer
	err := newAppWithGenerate(uc, cli.WithGetwd(sandboxGetwd)).Execute(context.Background(),
		[]string{"generate", "changelog", "--sandbox"}, &stdout, &stderr)
	if got := cli.ExitCode(err); got != 2 {
		t.Errorf("ExitCode = %d, want 2 (err=%v)", got, err)
	}
	if uc.called {
		t.Errorf("use case must not be called")
	}
}

// Origin URL with credentials (use-case error) → exit 10.
func TestExecute_Generate_SandboxUnsafeOrigin_Code10(t *testing.T) {
	uc := &fakeGenerateUseCase{err: fmt.Errorf("%w: origin embeds credentials", domain.ErrInvalidSandboxSetting)}
	var stdout, stderr bytes.Buffer
	err := newAppWithGenerate(uc, cli.WithGetwd(sandboxGetwd)).Execute(context.Background(),
		[]string{"generate", "devcontainer", "--sandbox"}, &stdout, &stderr)
	if !errors.Is(err, domain.ErrInvalidSandboxSetting) || cli.ExitCode(err) != 10 {
		t.Errorf("err = %v, ExitCode = %d, want sandbox sentinel / 10", err, cli.ExitCode(err))
	}
}
