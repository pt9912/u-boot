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

// LH-FA-CLI-005 / LH-FA-DEV-006/-007/-008: `--quiet` shows errors
// only; the disclosure of a sandbox security relaxation stays visible
// in init, generate and up. The distinction hangs on the entry
// (SecurityRelaxation), not on its code: the fixtures give a
// relaxation notice and a "has no effect" notice the same code.

func quietRelaxation() driving.WarningEntry {
	return driving.WarningEntry{
		Code: "LH-FA-DEV-007", Level: "warn", SecurityRelaxation: true,
		Message: "security relaxation for nested Podman: --cap-add=SYS_ADMIN",
	}
}

func quietNoEffect() driving.WarningEntry {
	return driving.WarningEntry{
		Code: "LH-FA-DEV-007", Level: "warn",
		Message: "devcontainer.sandbox.nestedRuntime=podman has no effect without the sandbox profile",
	}
}

func quietNoRemote() driving.WarningEntry {
	return driving.WarningEntry{
		Code: "LH-FA-DEV-006", Level: "warn",
		Message: "sandbox profile: no git remote 'origin' found",
	}
}

func quietRecreate() driving.WarningEntry {
	return driving.WarningEntry{
		Code: "LH-FA-UP-003", Level: "warn",
		Message: "devcontainer configuration changed; recreate the container",
	}
}

func quietAllWarnings() []driving.WarningEntry {
	return []driving.WarningEntry{quietRelaxation(), quietNoEffect(), quietNoRemote(), quietRecreate()}
}

func assertOnlyRelaxationVisible(t *testing.T, out string) {
	t.Helper()
	if !strings.Contains(out, "Warning (LH-FA-DEV-007): security relaxation for nested Podman") {
		t.Errorf("--quiet hid the security-relaxation warning:\n%s", out)
	}
	for _, gone := range []string{"has no effect", "no git remote", "recreate the container"} {
		if strings.Contains(out, gone) {
			t.Errorf("--quiet kept a non-relaxation warning (%q):\n%s", gone, out)
		}
	}
}

func TestExecute_Init_Quiet_KeepsOnlyRelaxationWarnings(t *testing.T) {
	t.Parallel()
	uc := &fakeInitUseCase{resp: driving.InitProjectResponse{Warnings: quietAllWarnings()}}
	var stdout, stderr bytes.Buffer
	if err := newApp(uc, cli.WithGetwd(sandboxGetwd)).Execute(context.Background(),
		[]string{"--quiet", "init", "--devcontainer", "--sandbox"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	assertOnlyRelaxationVisible(t, stdout.String())
}

// Boundary: without --quiet every warning stays (default behaviour).
func TestExecute_Init_NoQuiet_KeepsAllWarnings(t *testing.T) {
	t.Parallel()
	uc := &fakeInitUseCase{resp: driving.InitProjectResponse{Warnings: quietAllWarnings()}}
	var stdout, stderr bytes.Buffer
	if err := newApp(uc, cli.WithGetwd(sandboxGetwd)).Execute(context.Background(),
		[]string{"init", "--devcontainer", "--sandbox"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	for _, want := range []string{"security relaxation", "has no effect", "no git remote", "recreate the container"} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("default output lacks %q:\n%s", want, stdout.String())
		}
	}
}

// Negative: no relaxation, `--quiet` prints no warning line at all.
func TestExecute_Init_Quiet_WithoutRelaxationNoWarningLines(t *testing.T) {
	t.Parallel()
	uc := &fakeInitUseCase{resp: driving.InitProjectResponse{Warnings: []driving.WarningEntry{quietNoEffect(), quietNoRemote()}}}
	var stdout, stderr bytes.Buffer
	if err := newApp(uc, cli.WithGetwd(sandboxGetwd)).Execute(context.Background(),
		[]string{"--quiet", "init", "--devcontainer", "--sandbox"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if strings.Contains(stdout.String(), "Warning (") {
		t.Errorf("--quiet printed a non-relaxation warning:\n%s", stdout.String())
	}
}

func TestExecute_Generate_Quiet_KeepsOnlyRelaxationWarnings(t *testing.T) {
	t.Parallel()
	uc := &fakeGenerateUseCase{resp: driving.GenerateResponse{
		Artifact: domain.ArtifactDevcontainer, Action: driving.GenerateActionCreated,
		Changed:  []string{".devcontainer/devcontainer.json"},
		Warnings: quietAllWarnings(),
	}}
	var stdout, stderr bytes.Buffer
	if err := newAppWithGenerate(uc, cli.WithGetwd(sandboxGetwd)).Execute(context.Background(),
		[]string{"--quiet", "generate", "devcontainer", "--sandbox"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	assertOnlyRelaxationVisible(t, stdout.String())
}

func TestExecute_Generate_Quiet_WithoutRelaxationNoWarningLines(t *testing.T) {
	t.Parallel()
	uc := &fakeGenerateUseCase{resp: driving.GenerateResponse{
		Artifact: domain.ArtifactDevcontainer, Action: driving.GenerateActionCreated,
		Warnings: []driving.WarningEntry{quietNoEffect(), quietNoRemote()},
	}}
	var stdout, stderr bytes.Buffer
	if err := newAppWithGenerate(uc, cli.WithGetwd(sandboxGetwd)).Execute(context.Background(),
		[]string{"--quiet", "generate", "devcontainer"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if strings.Contains(stdout.String(), "Warning (") {
		t.Errorf("--quiet printed a non-relaxation warning:\n%s", stdout.String())
	}
}

// up: status table and diagnostics stay suppressed; of the warnings
// only the relaxation notice is written (the recreate warning,
// LH-FA-UP-003, is dropped).
func TestExecute_Up_Quiet_KeepsOnlyRelaxationWarnings(t *testing.T) {
	t.Parallel()
	uc := &fakeUpUseCase{resp: driving.UpResponse{
		Result: domain.UpResult{
			Stabilized: true,
			Services: []domain.ServiceStatus{
				{Name: "postgres", ContainerStatus: domain.StateRunning, Port: "5432:5432", Healthcheck: "healthy"},
			},
			Diagnostics: []domain.Diagnostic{{ID: "up.port.x.0", Severity: domain.SeverityWarn, Message: "port-diagnostic"}},
		},
		Warnings: quietAllWarnings(),
	}}
	var stdout, stderr bytes.Buffer
	if err := newAppWithUp(uc, cli.WithGetwd(func() (string, error) { return "/tmp/proj", nil })).Execute(
		context.Background(), []string{"--quiet", "up"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	out := stdout.String()
	assertOnlyRelaxationVisible(t, out)
	for _, gone := range []string{"postgres", "port-diagnostic"} {
		if strings.Contains(out, gone) {
			t.Errorf("--quiet kept status output %q:\n%s", gone, out)
		}
	}
}

// Negative: only a recreate warning, `up --quiet` stays empty.
func TestExecute_Up_Quiet_RecreateWarningOnlyIsSilent(t *testing.T) {
	t.Parallel()
	uc := &fakeUpUseCase{resp: driving.UpResponse{
		Result:   domain.UpResult{Stabilized: true},
		Warnings: []driving.WarningEntry{quietRecreate()},
	}}
	var stdout, stderr bytes.Buffer
	if err := newAppWithUp(uc, cli.WithGetwd(func() (string, error) { return "/tmp/proj", nil })).Execute(
		context.Background(), []string{"--quiet", "up"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if stdout.Len() != 0 {
		t.Errorf("--quiet printed the recreate warning: %q", stdout.String())
	}
}

// Boundary: without --quiet `up` keeps the recreate warning.
func TestExecute_Up_NoQuiet_KeepsRecreateWarning(t *testing.T) {
	t.Parallel()
	uc := &fakeUpUseCase{resp: driving.UpResponse{
		Result:   domain.UpResult{Stabilized: true},
		Warnings: []driving.WarningEntry{quietRecreate()},
	}}
	var stdout, stderr bytes.Buffer
	if err := newAppWithUp(uc, cli.WithGetwd(func() (string, error) { return "/tmp/proj", nil })).Execute(
		context.Background(), []string{"up"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !strings.Contains(stdout.String(), "recreate the container") {
		t.Errorf("default up lost the recreate warning: %q", stdout.String())
	}
}

// `--quiet --json` is a no-op for the warnings: every entry stays in
// diagnostics[], the SecurityRelaxation flag is not part of the wire form.
func TestExecute_QuietJSON_KeepsAllWarningsAsDiagnostics(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		run  func(args []string, stdout, stderr *bytes.Buffer) error
		args []string
	}{
		{"init", func(a []string, o, e *bytes.Buffer) error {
			uc := &fakeInitUseCase{resp: driving.InitProjectResponse{Warnings: quietAllWarnings()}}
			return newApp(uc, cli.WithGetwd(sandboxGetwd)).Execute(context.Background(), a, o, e)
		}, []string{"--quiet", "--json", "init", "--devcontainer", "--sandbox"}},
		{"generate", func(a []string, o, e *bytes.Buffer) error {
			uc := &fakeGenerateUseCase{resp: driving.GenerateResponse{Artifact: domain.ArtifactDevcontainer, Action: driving.GenerateActionCreated, Warnings: quietAllWarnings()}}
			return newAppWithGenerate(uc, cli.WithGetwd(sandboxGetwd)).Execute(context.Background(), a, o, e)
		}, []string{"--quiet", "--json", "generate", "devcontainer", "--sandbox"}},
		{"up", func(a []string, o, e *bytes.Buffer) error {
			uc := &fakeUpUseCase{resp: driving.UpResponse{Result: domain.UpResult{Stabilized: true}, Warnings: quietAllWarnings()}}
			return newAppWithUp(uc, cli.WithGetwd(sandboxGetwd)).Execute(context.Background(), a, o, e)
		}, []string{"--quiet", "--json", "up"}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var stdout, stderr bytes.Buffer
			if err := tc.run(tc.args, &stdout, &stderr); err != nil {
				t.Fatalf("Execute: %v", err)
			}
			out := stdout.String()
			if n := strings.Count(out, `"level":"warn"`) + strings.Count(out, `"level": "warn"`); n != 4 {
				t.Errorf("diagnostics warn entries = %d, want 4:\n%s", n, out)
			}
			if strings.Contains(strings.ToLower(out), "securityrelaxation") {
				t.Errorf("SecurityRelaxation leaked into the wire form:\n%s", out)
			}
		})
	}
}

// down --quiet: no success message, no volume list, also with --volumes.
func TestExecute_Down_Quiet_IsSilent(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		args []string
		resp driving.DownResponse
	}{
		{"stop", []string{"--quiet", "down"}, driving.DownResponse{}},
		{"volumes", []string{"--quiet", "down", "--volumes", "--yes"}, driving.DownResponse{RemovedVolumes: true, RemovedVolumeNames: []string{"pgdata"}}},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			uc := &fakeDownUseCase{resp: tc.resp}
			var stdout, stderr bytes.Buffer
			if err := newAppWithDown(uc, cli.WithGetwd(func() (string, error) { return "/tmp/proj", nil })).Execute(
				context.Background(), tc.args, &stdout, &stderr); err != nil {
				t.Fatalf("Execute: %v", err)
			}
			if stdout.Len() != 0 {
				t.Errorf("down --quiet printed: %q", stdout.String())
			}
		})
	}
}

// down --quiet holds the Compose stream back: success prints nothing
// on stderr, a failure prints the Compose lines (they carry the failure
// details) together with the error. Without --quiet the stream is live;
// --json leaves it as it is.
func TestExecute_Down_Quiet_ComposeStream(t *testing.T) {
	t.Parallel()
	const compose = " Container demo-db-1  Stopping\n Container demo-db-1  Error\n"
	failure := errors.New("docker compose down failed (exit 1)")
	for _, tc := range []struct {
		name       string
		args       []string
		err        error
		wantStream bool
	}{
		{"quiet-success", []string{"--quiet", "down"}, nil, false},
		{"quiet-failure", []string{"--quiet", "down"}, failure, true},
		{"loud-success", []string{"down"}, nil, true},
		{"loud-failure", []string{"down"}, failure, true},
		{"quiet-json-success", []string{"--quiet", "--json", "down"}, nil, true},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			uc := &fakeDownUseCase{compose: compose, err: tc.err}
			var stdout, stderr bytes.Buffer
			err := newAppWithDown(uc, cli.WithGetwd(func() (string, error) { return "/tmp/proj", nil })).Execute(
				context.Background(), tc.args, &stdout, &stderr)
			if (err != nil) != (tc.err != nil) {
				t.Fatalf("Execute error = %v, want error: %v", err, tc.err != nil)
			}
			if got := strings.Contains(stderr.String(), "Container demo-db-1  Error"); got != tc.wantStream {
				t.Errorf("compose stream on stderr = %v, want %v:\n%s", got, tc.wantStream, stderr.String())
			}
			if tc.err != nil && !errors.Is(err, tc.err) {
				t.Errorf("the failure is not returned: %v", err)
			}
		})
	}
}

// LH-FA-CLI-005 exception 2 / LH-FA-DEV-007: the three degradation
// fallbacks (container without daemon, `vfs` storage, no egress
// restriction) stay visible in the command output under `--quiet`,
// each one on its own. The flag (DegradationFallback) sits on the
// entry: the fixtures carry the code of the unflagged notices.
func quietFallbacks() map[string]driving.WarningEntry {
	return map[string]driving.WarningEntry{
		"no-daemon": {Code: "LH-FA-DEV-007", Level: "warn", DegradationFallback: true,
			Message: "nested Docker daemon is not ready; container runs without daemon"},
		"vfs": {Code: "LH-FA-DEV-007", Level: "warn", DegradationFallback: true,
			Message: "/dev/fuse is not available; nested Podman falls back to vfs storage"},
		"no-egress": {Code: "LH-FA-DEV-008", Level: "warn", DegradationFallback: true,
			Message: "NET_ADMIN is not grantable; the egress restriction does not apply"},
	}
}

func TestExecute_Quiet_KeepsEachDegradationFallback(t *testing.T) {
	t.Parallel()
	runners := map[string]func(ws []driving.WarningEntry, stdout, stderr *bytes.Buffer) error{
		"init": func(ws []driving.WarningEntry, o, e *bytes.Buffer) error {
			uc := &fakeInitUseCase{resp: driving.InitProjectResponse{Warnings: ws}}
			return newApp(uc, cli.WithGetwd(sandboxGetwd)).Execute(context.Background(), []string{"--quiet", "init", "--devcontainer", "--sandbox"}, o, e)
		},
		"generate": func(ws []driving.WarningEntry, o, e *bytes.Buffer) error {
			uc := &fakeGenerateUseCase{resp: driving.GenerateResponse{Artifact: domain.ArtifactDevcontainer, Action: driving.GenerateActionCreated, Warnings: ws}}
			return newAppWithGenerate(uc, cli.WithGetwd(sandboxGetwd)).Execute(context.Background(), []string{"--quiet", "generate", "devcontainer", "--sandbox"}, o, e)
		},
		"up": func(ws []driving.WarningEntry, o, e *bytes.Buffer) error {
			uc := &fakeUpUseCase{resp: driving.UpResponse{Result: domain.UpResult{Stabilized: true}, Warnings: ws}}
			return newAppWithUp(uc, cli.WithGetwd(sandboxGetwd)).Execute(context.Background(), []string{"--quiet", "up"}, o, e)
		},
	}
	for cmdName, run := range runners {
		for fbName, fb := range quietFallbacks() {
			cmdName, run, fbName, fb := cmdName, run, fbName, fb
			t.Run(cmdName+"/"+fbName, func(t *testing.T) {
				t.Parallel()
				// The unflagged notice sits next to the fallback and must go.
				ws := []driving.WarningEntry{quietNoEffect(), fb}
				var stdout, stderr bytes.Buffer
				if err := run(ws, &stdout, &stderr); err != nil {
					t.Fatalf("Execute: %v", err)
				}
				out := stdout.String()
				if !strings.Contains(out, "Warning ("+fb.Code+"): "+fb.Message) {
					t.Errorf("--quiet hid the %s fallback:\n%s", fbName, out)
				}
				if strings.Contains(out, "has no effect") {
					t.Errorf("--quiet kept the unflagged notice:\n%s", out)
				}
			})
		}
	}
}
