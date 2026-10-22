package cli_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/adapter/driving/cli/jsontestutil"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

const testContainerID = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

// fakeDevcontainerUseCase records the requests and plays the use case: it
// raises the configured disclosure through OnWarning, writes progress, and
// returns the configured response or error.
type fakeDevcontainerUseCase struct {
	upCalls   int
	upReq     driving.DevcontainerUpRequest
	upResp    driving.DevcontainerUpResponse
	upErr     error
	disclose  []driving.WarningEntry
	progress  string
	downCalls int
	downReq   driving.DevcontainerDownRequest
	downResp  driving.DevcontainerDownResponse
	downErr   error
}

func (f *fakeDevcontainerUseCase) Up(_ context.Context, req driving.DevcontainerUpRequest) (driving.DevcontainerUpResponse, error) {
	f.upCalls++
	f.upReq = req
	for _, w := range f.disclose {
		if req.OnWarning != nil {
			req.OnWarning(w)
		}
	}
	if f.progress != "" && req.ProgressSink != nil {
		_, _ = io.WriteString(req.ProgressSink, f.progress)
	}
	return f.upResp, f.upErr
}

func (f *fakeDevcontainerUseCase) Down(_ context.Context, req driving.DevcontainerDownRequest) (driving.DevcontainerDownResponse, error) {
	f.downCalls++
	f.downReq = req
	return f.downResp, f.downErr
}

func newAppWithDevcontainer(uc driving.DevcontainerUseCase, opts ...cli.Option) *cli.App {
	return newApp(&fakeInitUseCase{}, append([]cli.Option{cli.WithDevcontainerUseCase(uc), cli.WithGetwd(sandboxGetwd)}, opts...)...)
}

func socketDisclosure() driving.WarningEntry {
	return driving.WarningEntry{
		Code: "LH-FA-DEV-010", Level: "warn", SecurityRelaxation: true,
		Message: "this command mounts the host container-engine socket into a short-lived helper container (root-equivalent access)",
	}
}

func fallbackNote() driving.WarningEntry {
	return driving.WarningEntry{
		Code: "LH-FA-DEV-007", Level: "warn", DegradationFallback: true, Subject: ".devcontainer/devcontainer.json",
		Message: "state of the devcontainer container at its last start (2026-10-17T18:52:54Z): nested Podman uses vfs storage",
	}
}

func startedUseCase() *fakeDevcontainerUseCase {
	return &fakeDevcontainerUseCase{
		upResp:   driving.DevcontainerUpResponse{ContainerID: testContainerID, RemoteUser: "vscode", RemoteWorkspaceFolder: "/workspaces/demo", Warnings: []driving.WarningEntry{fallbackNote()}},
		disclose: []driving.WarningEntry{socketDisclosure()},
		progress: "Step 1/2 : FROM base\n",
	}
}

// ADR-0028 decision 3: bare `devcontainer` ends with exit 2, also with
// `--json`, and writes no envelope.
func TestExecute_Devcontainer_BareIsUsageError(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{{"devcontainer"}, {"devcontainer", "--json"}, {"--json", "devcontainer"}} {
		out, _, err := quietRun(t, newAppWithDevcontainer(&fakeDevcontainerUseCase{}), args...)
		if !errors.Is(err, cli.ErrDevcontainerSubcommandRequired) || cli.ExitCode(err) != 2 || out != "" {
			t.Errorf("%v: err=%v exit=%d stdout=%q, want exit 2 and no stdout", args, err, cli.ExitCode(err), out)
		}
	}
	_, _, err := quietRun(t, newAppWithDevcontainer(&fakeDevcontainerUseCase{}), "devcontainer", "frobnicate")
	if cli.ExitCode(err) != 2 {
		t.Errorf("unknown subcommand: exit=%d, want 2 (%v)", cli.ExitCode(err), err)
	}
}

// ADR-0032 decision 5: the help of `devcontainer up` says what the generated
// configuration mounts: the socket, if a feature that does so is enabled, such
// as `docker-cli`, which the sandbox profile rejects. The sentence is the
// sense of the disclosure `up` raises before the start.
func TestExecute_DevcontainerUp_HelpSaysWhatTheGeneratedConfigurationMounts(t *testing.T) {
	t.Parallel()
	out, errOut, err := quietRun(t, newAppWithDevcontainer(&fakeDevcontainerUseCase{}), "devcontainer", "up", "--help")
	if err != nil {
		t.Fatalf("devcontainer up --help: %v", err)
	}
	text := strings.Join(strings.Fields(out+errOut), " ")
	for _, part := range []string{
		"The devcontainer configuration that u-boot generates mounts the socket if a feature that does so is enabled, " +
			"such as the built-in feature docker-cli (the sandbox profile rejects that feature);",
		"a hand-written .devcontainer/devcontainer.json can mount it, and its initializeCommand runs with this access",
	} {
		if !strings.Contains(text, part) {
			t.Errorf("help lacks %q:\n%s", part, text)
		}
	}
	if strings.Contains(text, "does not mount the socket") {
		t.Errorf("help keeps the sentence that the generated configuration does not mount the socket:\n%s", text)
	}
}

// ADR-0030 decisions 1 and 2: the help of `devcontainer up` says that a start
// script's failure ends the command with exit 11 also for an earlier start of
// an already running container, names `--recreate`, and names the result line
// and the JSON field of an already running container together with the
// exception of `--recreate` (the line says "started", the field is false). The
// behaviour behind each sentence has its test (the application tests of `Up`,
// among them recreate_skips_the_query, and
// TestExecute_DevcontainerUp_AlreadyRunning).
func TestExecute_DevcontainerUp_HelpSaysWhenExit11AndAlreadyRunning(t *testing.T) {
	t.Parallel()
	out, errOut, err := quietRun(t, newAppWithDevcontainer(&fakeDevcontainerUseCase{}), "devcontainer", "up", "--help")
	if err != nil {
		t.Fatalf("devcontainer up --help: %v", err)
	}
	text := strings.Join(strings.Fields(out+errOut), " ")
	for _, part := range []string{
		"A start script that reported a failure ends the command with exit 11, also when the container was already running and this call did not run the script;",
		`"u-boot devcontainer up --recreate" creates the container anew.`,
		`If the container was already running and the command runs without --recreate, the result line says "already running" and the JSON data carry alreadyRunning: true; ` +
			`with --recreate the line says "started" and alreadyRunning is false.`,
		"script of the container aborted for an environment reason (also in an earlier start of the container, if u-boot.yaml configures the capability)",
	} {
		if !strings.Contains(text, part) {
			t.Errorf("help lacks %q:\n%s", part, text)
		}
	}
}

// ADR-0028 decision 3: neither subcommand writes a project file, so
// `--dry-run` and `--diff` end with exit 2 before the use case runs.
func TestExecute_Devcontainer_RejectsPreviewFlagsAndArguments(t *testing.T) {
	t.Parallel()
	for _, sub := range []string{"up", "down"} {
		for _, flag := range []string{"--dry-run", "--diff"} {
			uc := startedUseCase()
			out, errOut, err := quietRun(t, newAppWithDevcontainer(uc), "devcontainer", sub, flag)
			if !errors.Is(err, cli.ErrDevcontainerPreviewNotApplicable) || cli.ExitCode(err) != 2 || uc.upCalls+uc.downCalls != 0 || out != "" || errOut != "" {
				t.Errorf("%s %s: err=%v exit=%d calls=%d stdout=%q stderr=%q", sub, flag, err, cli.ExitCode(err), uc.upCalls+uc.downCalls, out, errOut)
			}
			out, _, err = quietRun(t, newAppWithDevcontainer(uc), "devcontainer", sub, flag, "--json")
			if cli.ExitCode(err) != 2 {
				t.Errorf("%s %s --json: exit=%d, want 2", sub, flag, cli.ExitCode(err))
			}
			jsontestutil.AssertMinimalEnvelope(t, []byte(out),
				jsontestutil.WithCommand("devcontainer"), jsontestutil.WithSubcommand(sub),
				jsontestutil.WithExitCode(2), jsontestutil.WithExpectedCodes("LH-FA-CLI-006"))
		}
		out, _, err := quietRun(t, newAppWithDevcontainer(startedUseCase()), "devcontainer", sub, "extra", "--json")
		if cli.ExitCode(err) != 2 {
			t.Errorf("%s extra argument: exit=%d, want 2", sub, cli.ExitCode(err))
		}
		jsontestutil.AssertMinimalEnvelope(t, []byte(out), jsontestutil.WithCommand("devcontainer"), jsontestutil.WithSubcommand(sub), jsontestutil.WithExitCode(2))
	}
}

// LH-FA-DEV-010: `up` passes the working directory and `--recreate`, shows
// the disclosure before anything else, the success line, then the fallback
// warnings; progress goes to stderr.
func TestExecute_DevcontainerUp_HumanOutput(t *testing.T) {
	t.Parallel()
	uc := startedUseCase()
	out, errOut, err := quietRun(t, newAppWithDevcontainer(uc), "devcontainer", "up", "--recreate")
	if err != nil {
		t.Fatal(err)
	}
	if uc.upCalls != 1 || uc.upReq.BaseDir != "/tmp/x/demo" || !uc.upReq.Recreate {
		t.Errorf("request = %+v after %d calls", uc.upReq, uc.upCalls)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 ||
		!strings.HasPrefix(lines[0], "Warning (LH-FA-DEV-010): this command mounts the host container-engine socket") ||
		lines[1] != "Devcontainer started: container 0123456789ab, remote user vscode, workspace /workspaces/demo" ||
		!strings.HasPrefix(lines[2], "Warning (LH-FA-DEV-007): state of the devcontainer container at its last start") {
		t.Errorf("stdout lines (disclosure, result, fallback):\n%s", out)
	}
	if errOut != "Step 1/2 : FROM base\n" {
		t.Errorf("stderr = %q, want the progress", errOut)
	}
	uc2 := startedUseCase()
	if _, _, err := quietRun(t, newAppWithDevcontainer(uc2), "devcontainer", "up"); err != nil || uc2.upReq.Recreate {
		t.Errorf("without --recreate: err=%v recreate=%v", err, uc2.upReq.Recreate)
	}
}

// LH-FA-CLI-005 exceptions 1 and 2: under `--quiet` the disclosure and the
// fallback warning stay, the success line and the progress go; the progress
// of a failed call is written after all.
func TestExecute_DevcontainerUp_Quiet(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		out, errOut, err := quietRun(t, newAppWithDevcontainer(startedUseCase()), "--quiet", "devcontainer", "up")
		if err != nil {
			t.Fatal(err)
		}
		if errOut != "" || strings.Contains(out, "Devcontainer started") ||
			!strings.Contains(out, "Warning (LH-FA-DEV-010):") || !strings.Contains(out, "Warning (LH-FA-DEV-007):") {
			t.Errorf("stdout=%q stderr=%q: want both warnings, no result line, no progress", out, errOut)
		}
	})
	t.Run("a notice that is neither relaxation nor fallback is dropped", func(t *testing.T) {
		t.Parallel()
		uc := startedUseCase()
		uc.upResp.Warnings = append(uc.upResp.Warnings, driving.WarningEntry{Code: "LH-FA-DEV-010", Level: "warn", Message: "a plain notice"})
		out, _, err := quietRun(t, newAppWithDevcontainer(uc), "--quiet", "devcontainer", "up")
		if err != nil || strings.Contains(out, "a plain notice") {
			t.Errorf("err=%v stdout=%q", err, out)
		}
	})
	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		uc := startedUseCase()
		uc.upErr = fmt.Errorf("devcontainer up failed (exit status 1: boom): %w", driven.ErrDevcontainerRuntime)
		out, errOut, err := quietRun(t, newAppWithDevcontainer(uc), "--quiet", "devcontainer", "up")
		if cli.ExitCode(err) != 12 || errOut != "Step 1/2 : FROM base\n" || !strings.Contains(out, "Warning (LH-FA-DEV-010):") {
			t.Errorf("err=%v exit=%d stdout=%q stderr=%q: want exit 12, the held progress and the disclosure", err, cli.ExitCode(err), out, errOut)
		}
	})
}

// LH-NFA-USE-004: `--json` writes exactly one envelope; the disclosure is a
// `diagnostics[]` entry with the text of the console; the status is `warn`;
// `--quiet --json` equals `--json`; progress is not written.
func TestExecute_DevcontainerUp_JSON(t *testing.T) {
	t.Parallel()
	out, errOut, err := quietRun(t, newAppWithDevcontainer(startedUseCase()), "devcontainer", "up", "--json")
	if err != nil || errOut != "" {
		t.Fatalf("err=%v stderr=%q: progress must not be written under --json", err, errOut)
	}
	jsontestutil.AssertMinimalEnvelope(t, []byte(out),
		jsontestutil.WithCommand("devcontainer"), jsontestutil.WithSubcommand("up"), jsontestutil.WithExitCode(0),
		jsontestutil.WithExpectedCodes("LH-FA-DEV-010", "LH-FA-DEV-007"))
	var env struct {
		Status      string `json:"status"`
		Diagnostics []struct {
			Level, Code, Message, File string
		} `json:"diagnostics"`
		Data struct {
			ContainerID           *string `json:"containerId"`
			RemoteUser            *string `json:"remoteUser"`
			RemoteWorkspaceFolder *string `json:"remoteWorkspaceFolder"`
			AlreadyRunning        *bool   `json:"alreadyRunning"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(out), &env); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	if env.Status != "warn" || len(env.Diagnostics) != 2 {
		t.Fatalf("status=%q diagnostics=%+v, want warn with the disclosure and the fallback", env.Status, env.Diagnostics)
	}
	d := env.Diagnostics[0]
	if d.Level != "warn" || d.Code != "LH-FA-DEV-010" || d.Message != socketDisclosure().Message || d.File != "" {
		t.Errorf("disclosure diagnostic = %+v", d)
	}
	if f := env.Diagnostics[1]; f.Code != "LH-FA-DEV-007" || f.File != ".devcontainer/devcontainer.json" {
		t.Errorf("fallback diagnostic = %+v", f)
	}
	if env.Data.ContainerID == nil || *env.Data.ContainerID != testContainerID || env.Data.RemoteUser == nil || *env.Data.RemoteUser != "vscode" ||
		env.Data.RemoteWorkspaceFolder == nil || *env.Data.RemoteWorkspaceFolder != "/workspaces/demo" {
		t.Errorf("data = %+v", env.Data)
	}
	// ADR-0030 decision 2: the field is always there, false for a container
	// that was started by the call.
	if env.Data.AlreadyRunning == nil || *env.Data.AlreadyRunning {
		t.Errorf("data.alreadyRunning = %v, want present and false", env.Data.AlreadyRunning)
	}
	quiet, _, err := quietRun(t, newAppWithDevcontainer(startedUseCase()), "--quiet", "devcontainer", "up", "--json")
	if err != nil || quiet != out {
		t.Errorf("--quiet --json differs from --json (err=%v):\n%s\n%s", err, quiet, out)
	}
	// Console and JSON carry the same text.
	console, _, _ := quietRun(t, newAppWithDevcontainer(startedUseCase()), "devcontainer", "up")
	if !strings.Contains(console, "Warning (LH-FA-DEV-010): "+d.Message+"\n") {
		t.Errorf("the console text differs from the JSON message %q:\n%s", d.Message, console)
	}
}

// ADR-0030 decision 2, LH-FA-CLI-007: a container that was already running
// changes the result line and the field `alreadyRunning`, not the exit code,
// the envelope status or the warnings; `--quiet` writes no result line, as for
// a started container.
func TestExecute_DevcontainerUp_AlreadyRunning(t *testing.T) {
	t.Parallel()
	running := func() *fakeDevcontainerUseCase {
		uc := startedUseCase()
		uc.upResp.AlreadyRunning = true
		return uc
	}
	t.Run("human", func(t *testing.T) {
		t.Parallel()
		out, _, err := quietRun(t, newAppWithDevcontainer(running()), "devcontainer", "up")
		if err != nil {
			t.Fatal(err)
		}
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) != 3 ||
			!strings.HasPrefix(lines[0], "Warning (LH-FA-DEV-010):") ||
			lines[1] != "Devcontainer already running: container 0123456789ab, remote user vscode, workspace /workspaces/demo" ||
			!strings.HasPrefix(lines[2], "Warning (LH-FA-DEV-007):") {
			t.Errorf("stdout lines (disclosure, result, fallback):\n%s", out)
		}
		if strings.Contains(out, "Devcontainer started") {
			t.Errorf("an already running container is not reported as started:\n%s", out)
		}
	})
	t.Run("quiet", func(t *testing.T) {
		t.Parallel()
		out, _, err := quietRun(t, newAppWithDevcontainer(running()), "--quiet", "devcontainer", "up")
		if err != nil || strings.Contains(out, "Devcontainer ") {
			t.Errorf("err=%v stdout=%q, want no result line", err, out)
		}
	})
	t.Run("json", func(t *testing.T) {
		t.Parallel()
		out, _, err := quietRun(t, newAppWithDevcontainer(running()), "devcontainer", "up", "--json")
		if err != nil {
			t.Fatal(err)
		}
		jsontestutil.AssertMinimalEnvelope(t, []byte(out),
			jsontestutil.WithCommand("devcontainer"), jsontestutil.WithSubcommand("up"), jsontestutil.WithExitCode(0),
			jsontestutil.WithExpectedCodes("LH-FA-DEV-010", "LH-FA-DEV-007"),
			jsontestutil.WithDataKeyPresent("alreadyRunning", true), jsontestutil.WithDataKeyPresent("containerId", testContainerID))
		if !strings.Contains(out, `"status":"warn"`) {
			t.Errorf("the status stays warn (the disclosure):\n%s", out)
		}
	})
}

// ADR-0028 decision 4: a failed `up` keeps the disclosure in the envelope when
// the use case had raised it, as a `diagnostics[]` entry before the error.
func TestExecute_DevcontainerUp_JSONErrorKeepsTheDisclosure(t *testing.T) {
	t.Parallel()
	uc := startedUseCase()
	uc.upErr = fmt.Errorf("devcontainer up failed: %w", driven.ErrDevcontainerRuntime)
	out, _, err := quietRun(t, newAppWithDevcontainer(uc), "devcontainer", "up", "--json")
	if cli.ExitCode(err) != 12 {
		t.Fatalf("exit=%d, want 12 (%v)", cli.ExitCode(err), err)
	}
	jsontestutil.AssertMinimalEnvelope(t, []byte(out),
		jsontestutil.WithCommand("devcontainer"), jsontestutil.WithSubcommand("up"), jsontestutil.WithExitCode(12),
		jsontestutil.WithExpectedCodes("LH-FA-DEV-010"), jsontestutil.WithDataKeyAbsent("containerId"))
	var env struct {
		Status      string `json:"status"`
		Diagnostics []struct{ Level, Code string }
	}
	if err := json.Unmarshal([]byte(out), &env); err != nil {
		t.Fatal(err)
	}
	if env.Status != "error" || len(env.Diagnostics) != 2 || env.Diagnostics[0].Level != "warn" || env.Diagnostics[1].Level != "error" {
		t.Errorf("status=%q diagnostics=%+v, want the warn disclosure, then the error", env.Status, env.Diagnostics)
	}
	// A failure before the disclosure carries no disclosure.
	uc2 := &fakeDevcontainerUseCase{upErr: fmt.Errorf("%w: looked for x", driving.ErrDevcontainerConfigMissing)}
	out, _, _ = quietRun(t, newAppWithDevcontainer(uc2), "devcontainer", "up", "--json")
	if strings.Contains(out, `"warn"`) {
		t.Errorf("a failure before the helper starts must not disclose the socket:\n%s", out)
	}
}

// LH-FA-CLI-006 / ADR-0028 decision 3 / ADR-0029 decision 1: every sentinel
// has its exit code in both classifiers — `ExitCode` and the diagnostic
// mapping of each subcommand — and the path of the project stays out of the
// message.
func TestExecute_Devcontainer_SentinelsPerSubcommand(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		err  error
		exit int
		code string
		subs []string
	}{
		{"docker unavailable", fmt.Errorf("docker binary not on PATH: %w", driven.ErrDockerUnavailable), 11, "LH-FA-DEV-010", []string{"up", "down"}},
		{"engine not local", driven.ErrContainerEngineNotLocal, 11, "LH-FA-DEV-010", []string{"up", "down"}},
		{"start script failed", fmt.Errorf("%w: /tmp/x/demo boom", driving.ErrDevcontainerStartScriptFailed), 11, "LH-FA-DEV-010", []string{"up"}},
		{"runtime", fmt.Errorf("docker stop failed: %w", driven.ErrDevcontainerRuntime), 12, "LH-FA-DEV-010", []string{"up", "down"}},
		{"no configuration", fmt.Errorf("%w (looked for x)", driving.ErrDevcontainerConfigMissing), 10, "LH-FA-DEV-010", []string{"up"}},
		{"folder", driving.ErrDevcontainerProjectFolderInvalid, 10, "LH-FA-DEV-010", []string{"up"}},
		{"variable name", driving.ErrDevcontainerEnvNameInvalid, 10, "LH-FA-DEV-010", []string{"up"}},
		{"file system", fmt.Errorf("read /tmp/x/demo/.devcontainer/devcontainer.json: %w", driving.ErrDevcontainerFileSystem), 14, "LH-NFA-REL-003", []string{"up"}},
		{"other", errors.New("boom"), 1, "LH-FA-CLI-006", []string{"up", "down"}},
		// The classes do not overlap: the abort of a start script wins over
		// the runtime class if an error carries both.
		{"start script and runtime", fmt.Errorf("%w: %w", driving.ErrDevcontainerStartScriptFailed, driven.ErrDevcontainerRuntime), 11, "LH-FA-DEV-010", []string{"up"}},
	}
	for _, tc := range cases {
		for _, sub := range tc.subs {
			t.Run(tc.name+"/"+sub, func(t *testing.T) {
				t.Parallel()
				if got := cli.ExitCode(tc.err); got != tc.exit {
					t.Errorf("ExitCode(%v) = %d, want %d", tc.err, got, tc.exit)
				}
				run := func(extra ...string) (string, error) {
					uc := &fakeDevcontainerUseCase{upErr: tc.err, downErr: tc.err}
					out, _, err := quietRun(t, newAppWithDevcontainer(uc), append([]string{"devcontainer", sub}, extra...)...)
					return out, err
				}
				out, err := run()
				if !errors.Is(err, tc.err) || cli.ExitCode(err) != tc.exit || out != "" {
					t.Fatalf("human: err=%v exit=%d stdout=%q", err, cli.ExitCode(err), out)
				}
				out, err = run("--json")
				if cli.ExitCode(err) != tc.exit {
					t.Fatalf("json: exit=%d, want %d", cli.ExitCode(err), tc.exit)
				}
				jsontestutil.AssertMinimalEnvelope(t, []byte(out),
					jsontestutil.WithCommand("devcontainer"), jsontestutil.WithSubcommand(sub),
					jsontestutil.WithExitCode(tc.exit), jsontestutil.WithExpectedCodes(tc.code))
				if strings.Contains(out, "/tmp/x/demo") {
					t.Errorf("the project path leaks into the diagnostic:\n%s", out)
				}
			})
		}
	}
}

// A command without its use case fails with a wiring error instead of doing
// nothing.
func TestExecute_Devcontainer_NotWired(t *testing.T) {
	t.Parallel()
	for _, sub := range []string{"up", "down"} {
		_, _, err := quietRun(t, newApp(&fakeInitUseCase{}, cli.WithGetwd(sandboxGetwd)), "devcontainer", sub)
		if err == nil || !strings.Contains(err.Error(), "not wired") || cli.ExitCode(err) != 1 {
			t.Errorf("%s: err=%v exit=%d", sub, err, cli.ExitCode(err))
		}
	}
}

// ADR-0028 decision 3: `down` names what happened; `--quiet` writes nothing;
// the envelope carries `containerId` (empty without a container) and
// `stopped`.
func TestExecute_DevcontainerDown(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		resp    driving.DevcontainerDownResponse
		human   string
		id      string
		stopped bool
	}{
		{"stopped", driving.DevcontainerDownResponse{ContainerID: testContainerID, Stopped: true}, "Devcontainer stopped: container 0123456789ab (not removed)\n", testContainerID, true},
		{"already stopped", driving.DevcontainerDownResponse{ContainerID: testContainerID}, "Devcontainer container 0123456789ab is not running; nothing to stop\n", testContainerID, false},
		{"no container", driving.DevcontainerDownResponse{}, "No devcontainer container found for this project; nothing to stop\n", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			newUC := func() *fakeDevcontainerUseCase { return &fakeDevcontainerUseCase{downResp: tc.resp} }
			uc := newUC()
			out, errOut, err := quietRun(t, newAppWithDevcontainer(uc), "devcontainer", "down")
			if err != nil || out != tc.human || errOut != "" || uc.downReq.BaseDir != "/tmp/x/demo" {
				t.Errorf("human: err=%v stdout=%q stderr=%q req=%+v", err, out, errOut, uc.downReq)
			}
			out, _, err = quietRun(t, newAppWithDevcontainer(newUC()), "--quiet", "devcontainer", "down")
			if err != nil || out != "" {
				t.Errorf("--quiet: err=%v stdout=%q, want nothing", err, out)
			}
			jsonOut, _, err := quietRun(t, newAppWithDevcontainer(newUC()), "devcontainer", "down", "--json")
			if err != nil {
				t.Fatal(err)
			}
			jsontestutil.AssertMinimalEnvelope(t, []byte(jsonOut),
				jsontestutil.WithCommand("devcontainer"), jsontestutil.WithSubcommand("down"), jsontestutil.WithExitCode(0),
				jsontestutil.WithDataKeyPresent("containerId", tc.id), jsontestutil.WithDataKeyPresent("stopped", tc.stopped))
			if !strings.Contains(jsonOut, `"status":"ok"`) || !strings.Contains(jsonOut, `"diagnostics":[]`) {
				t.Errorf("down raises no diagnostic (no helper container, no socket):\n%s", jsonOut)
			}
			quietJSON, _, _ := quietRun(t, newAppWithDevcontainer(newUC()), "--quiet", "devcontainer", "down", "--json")
			if quietJSON != jsonOut {
				t.Errorf("--quiet --json differs from --json:\n%s\n%s", quietJSON, jsonOut)
			}
		})
	}
}
