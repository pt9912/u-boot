package cli_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/adapter/driving/cli/jsontestutil"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// shellRun runs `u-boot <args>` against the fake use case with a standard
// input that is a buffer and a probe that says "no terminal", so no test
// depends on the standard input of the test process.
func shellRun(t *testing.T, uc *fakeDevcontainerUseCase, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	app := newAppWithDevcontainer(uc,
		cli.WithStdin(strings.NewReader("")), cli.WithTerminalProbe(func(io.Reader) bool { return false }))
	return quietRun(t, app, args...)
}

// sessionUseCase is the fake of a start with a disclosure, progress and a
// fallback warning, followed by a session that writes one line to each stream.
func sessionUseCase() *fakeDevcontainerUseCase {
	uc := startedUseCase()
	uc.sessionOut, uc.sessionErr = "line of the session on stdout\n", "line of the session on stderr\n"
	return uc
}

func nonEmptyLines(s string) []string {
	var out []string
	for _, line := range strings.Split(s, "\n") {
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

// ADR-0031 decision 2: bare `devcontainer` and the help of the group name the
// third subcommand.
func TestExecute_Devcontainer_NamesShell(t *testing.T) {
	t.Parallel()
	if !strings.Contains(cli.ErrDevcontainerSubcommandRequired.Error(), "u-boot devcontainer shell") {
		t.Errorf("the error of bare `devcontainer` does not name shell: %v", cli.ErrDevcontainerSubcommandRequired)
	}
	out, errOut, err := quietRun(t, newAppWithDevcontainer(&fakeDevcontainerUseCase{}), "devcontainer", "--help")
	if err != nil {
		t.Fatal(err)
	}
	if text := strings.Join(strings.Fields(out+errOut), " "); !strings.Contains(text, "u-boot devcontainer shell start the devcontainer and open a shell in it") {
		t.Errorf("the help of the group lacks the line of shell:\n%s", text)
	}
}

// ADR-0031 decisions 1 to 3, welle-devcontainer-shell-und-zweites-up §5
// question 3: the help of `shell` says what the command does, carries no
// requirement, ADR or specification identifier (not in the line of the exit
// codes either) and says nothing about signals, the end of the session or a
// login shell, which no test of this command pins.
func TestExecute_DevcontainerShell_Help(t *testing.T) {
	t.Parallel()
	out, errOut, err := quietRun(t, newAppWithDevcontainer(&fakeDevcontainerUseCase{}), "devcontainer", "shell", "--help")
	if err != nil {
		t.Fatalf("devcontainer shell --help: %v", err)
	}
	text := strings.Join(strings.Fields(out+errOut), " ")
	for _, id := range []string{"LH-", "ADR-", "SPEC-", "ARC-"} {
		if strings.Contains(text, id) {
			t.Errorf("the help of shell carries %q:\n%s", id, text)
		}
	}
	for _, part := range []string{
		`as "u-boot devcontainer up" does without --recreate, and open an interactive shell in the container`,
		"for the duration of the start (root-equivalent access to the engine). The command says so before it starts the helper, also under --quiet.",
		"The shell runs without a helper container and without a socket that this command mounts.",
		"goes to standard error; standard output belongs to the shell. --quiet holds the progress back.",
		"With a terminal on standard input the shell is attached to it; otherwise it reads its commands from standard input.",
		`A start script of the container that reported a failure ends the command with exit 11 and opens no shell, also when the container was already running; "u-boot devcontainer up --recreate" creates the container anew.`,
		"Once the shell has begun, the command ends with exit 0 whatever the exit status of the shell was;",
		`use "docker exec".`,
		"--json, --dry-run, --diff and arguments are rejected with exit 2. --yes and --no-interactive have no effect: the command asks for no confirmation.",
		"Exit codes: 0 the shell has begun and ended",
		"12 any other failure of the build or the start; the check before the shell opens fails (user, workspace folder or shell missing in the container)",
	} {
		if !strings.Contains(text, part) {
			t.Errorf("the help of shell lacks %q:\n%s", part, text)
		}
	}
	for _, word := range []string{"login", "signal", "kill", "Ctrl", "SIGINT", "closes the terminal"} {
		if strings.Contains(strings.ToLower(text), strings.ToLower(word)) {
			t.Errorf("the help of shell says %q, which no test of this command pins:\n%s", word, text)
		}
	}
}

// ADR-0031 decision 2: `--dry-run`, `--diff`, `--json` and an argument end
// with exit 2 before the use case runs, in each sentinel, and `--json` makes
// the rejection an envelope with `subcommand: shell`. Without `--json` nothing
// is written to either stream.
func TestExecute_DevcontainerShell_RejectsFlagsAndArguments(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		args []string
		want error // nil: a Cobra error for the argument
	}{
		{"dry-run", []string{"--dry-run"}, cli.ErrDevcontainerPreviewNotApplicable},
		{"diff", []string{"--diff"}, cli.ErrDevcontainerPreviewNotApplicable},
		{"json", []string{"--json"}, cli.ErrDevcontainerShellJSONNotApplicable},
		{"extra argument", []string{"extra"}, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			uc := sessionUseCase()
			out, errOut, err := shellRun(t, uc, append([]string{"devcontainer", "shell"}, tc.args...)...)
			if cli.ExitCode(err) != 2 || uc.shellCalls != 0 || (tc.want != nil && !errors.Is(err, tc.want)) {
				t.Errorf("err=%v exit=%d calls=%d, want exit 2 and no call of the use case", err, cli.ExitCode(err), uc.shellCalls)
			}
			if tc.name != "json" && (out != "" || errOut != "") {
				t.Errorf("stdout=%q stderr=%q, want nothing written", out, errOut)
			}
			if errOut != "" {
				t.Errorf("stderr=%q, want nothing: u-boot prints nothing before the rejection", errOut)
			}
			// With `--json` the rejection is the envelope, and it names shell.
			jsonArgs := append([]string{"devcontainer", "shell", "--json"}, tc.args...)
			jout, _, jerr := shellRun(t, sessionUseCase(), jsonArgs...)
			if cli.ExitCode(jerr) != 2 {
				t.Errorf("--json %v: exit=%d, want 2 (%v)", tc.args, cli.ExitCode(jerr), jerr)
			}
			jsontestutil.AssertMinimalEnvelope(t, []byte(jout),
				jsontestutil.WithCommand("devcontainer"), jsontestutil.WithSubcommand("shell"),
				jsontestutil.WithExitCode(2), jsontestutil.WithExpectedCodes("LH-FA-CLI-006"))
		})
	}
	t.Run("json before the command", func(t *testing.T) {
		t.Parallel()
		out, _, err := shellRun(t, sessionUseCase(), "--json", "devcontainer", "shell")
		if !errors.Is(err, cli.ErrDevcontainerShellJSONNotApplicable) || cli.ExitCode(err) != 2 {
			t.Errorf("err=%v exit=%d", err, cli.ExitCode(err))
		}
		jsontestutil.AssertMinimalEnvelope(t, []byte(out), jsontestutil.WithCommand("devcontainer"), jsontestutil.WithSubcommand("shell"), jsontestutil.WithExitCode(2))
	})
	t.Run("the rejection wins over the missing use case", func(t *testing.T) {
		t.Parallel()
		_, _, err := quietRun(t, newApp(&fakeInitUseCase{}, cli.WithGetwd(sandboxGetwd)), "devcontainer", "shell", "--json")
		if !errors.Is(err, cli.ErrDevcontainerShellJSONNotApplicable) || cli.ExitCode(err) != 2 {
			t.Errorf("err=%v exit=%d, want the rejection of --json (exit 2), not a wiring error", err, cli.ExitCode(err))
		}
	})
}

// Both sentinels of exit class 2 are in `ExitCode`.
func TestExitCode_DevcontainerShellUsageSentinels(t *testing.T) {
	t.Parallel()
	for _, err := range []error{
		cli.ErrDevcontainerShellJSONNotApplicable, cli.ErrDevcontainerPreviewNotApplicable, cli.ErrDevcontainerSubcommandRequired,
		fmt.Errorf("wrapped: %w", cli.ErrDevcontainerShellJSONNotApplicable),
	} {
		if got := cli.ExitCode(err); got != 2 {
			t.Errorf("ExitCode(%v) = %d, want 2", err, got)
		}
	}
}

// ADR-0031 decision 2: everything u-boot itself writes goes to stderr — the
// disclosure first, the progress, the fallback warning — in this order, and the
// session's own lines reach the streams of the command; stdout holds nothing
// else, and no success line follows. The use case gets the working directory,
// a progress sink and the context of the command; the command ends with 0 once
// the use case returns nil.
func TestExecute_DevcontainerShell_HumanOutput(t *testing.T) {
	t.Parallel()
	type ctxKey struct{}
	uc := sessionUseCase()
	uc.ctxKey = ctxKey{}
	app := newAppWithDevcontainer(uc, cli.WithStdin(strings.NewReader("")), cli.WithTerminalProbe(func(io.Reader) bool { return false }))
	var o, e strings.Builder
	err := app.Execute(context.WithValue(context.Background(), ctxKey{}, "marker"), []string{"devcontainer", "shell"}, &o, &e)
	if err != nil || cli.ExitCode(err) != 0 {
		t.Fatalf("err=%v", err)
	}
	if o.String() != "line of the session on stdout\n" {
		t.Errorf("stdout = %q, want the line of the session and nothing of u-boot", o.String())
	}
	lines := nonEmptyLines(e.String())
	if len(lines) != 4 ||
		!strings.HasPrefix(lines[0], "Warning (LH-FA-DEV-010): this command mounts the host container-engine socket") ||
		lines[1] != "Step 1/2 : FROM base" ||
		!strings.HasPrefix(lines[2], "Warning (LH-FA-DEV-007): state of the devcontainer container at its last start") ||
		lines[3] != "line of the session on stderr" {
		t.Errorf("stderr lines (disclosure, progress, fallback, session):\n%s", e.String())
	}
	if strings.Contains(o.String()+e.String(), "Devcontainer ") {
		t.Errorf("a result line was written:\nstdout=%q\nstderr=%q", o.String(), e.String())
	}
	req := uc.shellReq
	if uc.shellCalls != 1 || req.BaseDir != "/tmp/x/demo" || req.ProgressSink == nil || req.OnWarning == nil || req.Terminal {
		t.Errorf("request = %+v after %d calls", req, uc.shellCalls)
	}
	if uc.ctxValue != "marker" {
		t.Errorf("the use case did not get the context of the command")
	}
}

// LH-FA-CLI-005 (Ausnahmen 1 und 2, Satz zur Shell): under `--quiet` the
// session's lines stay unchanged in the streams of the command, and u-boot
// itself prints only the exceptions — the disclosure of the socket and the
// fallback warning — and no progress. Guard: the same session without
// `--quiet` gives the same session lines and the progress as well.
func TestExecute_DevcontainerShell_Quiet(t *testing.T) {
	t.Parallel()
	session := func(args ...string) (string, []string) {
		out, errOut, err := shellRun(t, sessionUseCase(), append(args, "devcontainer", "shell")...)
		if err != nil {
			t.Fatal(err)
		}
		return out, nonEmptyLines(errOut)
	}
	quietOut, quietErr := session("--quiet")
	loudOut, loudErr := session()
	if quietOut != "line of the session on stdout\n" || quietOut != loudOut {
		t.Errorf("session stdout: quiet %q, default %q, want the same line", quietOut, loudOut)
	}
	if len(quietErr) != 3 ||
		!strings.HasPrefix(quietErr[0], "Warning (LH-FA-DEV-010):") ||
		!strings.HasPrefix(quietErr[1], "Warning (LH-FA-DEV-007):") ||
		quietErr[2] != "line of the session on stderr" {
		t.Errorf("--quiet stderr = %q, want the disclosure, the fallback warning and the session line", quietErr)
	}
	if len(loudErr) != 4 || loudErr[1] != "Step 1/2 : FROM base" || loudErr[3] != "line of the session on stderr" {
		t.Errorf("default stderr = %q, want the progress and the session line as well", loudErr)
	}
	t.Run("a notice that is neither relaxation nor fallback is dropped", func(t *testing.T) {
		t.Parallel()
		uc := sessionUseCase()
		uc.upResp.Warnings = append(uc.upResp.Warnings, driving.WarningEntry{Code: "LH-FA-DEV-010", Level: "warn", Message: "a plain notice"})
		_, errOut, err := shellRun(t, uc, "--quiet", "devcontainer", "shell")
		if err != nil || strings.Contains(errOut, "a plain notice") {
			t.Errorf("err=%v stderr=%q", err, errOut)
		}
		_, errOut, _ = shellRun(t, uc, "devcontainer", "shell")
		if !strings.Contains(errOut, "a plain notice") {
			t.Errorf("without --quiet the notice is shown: %q", errOut)
		}
	})
	t.Run("the progress of a failed call is written after all", func(t *testing.T) {
		t.Parallel()
		uc := startedUseCase()
		uc.shellErr = fmt.Errorf("devcontainer up failed (exit status 1: boom): %w", driven.ErrDevcontainerRuntime)
		out, errOut, err := shellRun(t, uc, "--quiet", "devcontainer", "shell")
		if cli.ExitCode(err) != 12 || out != "" || !strings.Contains(errOut, "Step 1/2 : FROM base\n") || !strings.Contains(errOut, "Warning (LH-FA-DEV-010):") {
			t.Errorf("exit=%d stdout=%q stderr=%q: want exit 12, the held progress and the disclosure", cli.ExitCode(err), out, errOut)
		}
	})
}

// ADR-0031 decision 1: the standard input of the command is the one the
// request carries, and the terminal flag is the answer of the probe to that
// reader — both answers pass through. The default probe says "no terminal" for
// a buffer, a pipe and the null device, and "terminal" for a pseudo-terminal.
func TestExecute_DevcontainerShell_Terminal(t *testing.T) {
	t.Parallel()
	t.Run("the probe is asked about the standard input and its answer passes through", func(t *testing.T) {
		t.Parallel()
		stdin := strings.NewReader("ls\n")
		for _, answer := range []bool{true, false} {
			var asked io.Reader
			uc := startedUseCase()
			app := newAppWithDevcontainer(uc, cli.WithStdin(stdin), cli.WithTerminalProbe(func(r io.Reader) bool { asked = r; return answer }))
			if _, _, err := quietRun(t, app, "devcontainer", "shell"); err != nil {
				t.Fatal(err)
			}
			if asked != io.Reader(stdin) || uc.shellReq.Stdin != io.Reader(stdin) || uc.shellReq.Terminal != answer {
				t.Errorf("answer %v: probe asked %v, request stdin %v terminal %v", answer, asked, uc.shellReq.Stdin, uc.shellReq.Terminal)
			}
		}
	})
	t.Run("the default probe", func(t *testing.T) {
		t.Parallel()
		pipeR, pipeW, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = pipeR.Close(); _ = pipeW.Close() })
		null, err := os.Open(os.DevNull)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = null.Close() })
		cases := []struct {
			name  string
			stdin io.Reader
			want  bool
		}{
			{"buffer", strings.NewReader("ls\n"), false},
			{"pipe", pipeR, false},
			{"null device", null, false},
		}
		if ptmx, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0); err == nil {
			t.Cleanup(func() { _ = ptmx.Close() })
			cases = append(cases, struct {
				name  string
				stdin io.Reader
				want  bool
			}{"pseudo-terminal", ptmx, true})
		} else {
			t.Logf("no /dev/ptmx (%v): the pseudo-terminal case is covered by the Docker test under `script`", err)
		}
		for _, tc := range cases {
			uc := startedUseCase()
			if _, _, err := quietRun(t, newAppWithDevcontainer(uc, cli.WithStdin(tc.stdin)), "devcontainer", "shell"); err != nil {
				t.Fatal(err)
			}
			if uc.shellReq.Terminal != tc.want {
				t.Errorf("%s: terminal = %v, want %v", tc.name, uc.shellReq.Terminal, tc.want)
			}
		}
	})
}

// The check of the standard input: a terminal is a character device that is
// not the null device; everything else is none.
func TestIsTerminal(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	regular, err := os.Create(dir + "/file")
	if err != nil {
		t.Fatal(err)
	}
	pipeR, pipeW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	null, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	closed, err := os.Create(dir + "/closed")
	if err != nil {
		t.Fatal(err)
	}
	_ = closed.Close()
	t.Cleanup(func() { _ = regular.Close(); _ = pipeR.Close(); _ = pipeW.Close(); _ = null.Close() })

	cases := []struct {
		name string
		in   io.Reader
		want bool
	}{
		{"nil reader", nil, false},
		{"buffer", strings.NewReader("x"), false},
		{"regular file", regular, false},
		{"pipe", pipeR, false},
		{"null device", null, false},
		{"closed file", closed, false},
		{"nil file", (*os.File)(nil), false},
	}
	if ptmx, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0); err == nil {
		t.Cleanup(func() { _ = ptmx.Close() })
		cases = append(cases, struct {
			name string
			in   io.Reader
			want bool
		}{"pseudo-terminal master", ptmx, true})
	} else {
		t.Logf("no /dev/ptmx (%v): the positive case is covered by the Docker test under `script`", err)
	}
	for _, tc := range cases {
		if got := cli.IsTerminalForTest(tc.in); got != tc.want {
			t.Errorf("%s: isTerminal = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// Risk "--yes und --no-interactive bei shell" (plan §6, welle-devcontainer-
// shell-und-zweites-up §3 point 3 (vi)): both global flags change nothing; the
// command asks for no confirmation, so the session opens as without them.
func TestExecute_DevcontainerShell_YesAndNoInteractiveHaveNoEffect(t *testing.T) {
	t.Parallel()
	for _, flag := range []string{"--yes", "--no-interactive"} {
		t.Run(flag, func(t *testing.T) {
			t.Parallel()
			uc := sessionUseCase()
			out, errOut, err := shellRun(t, uc, "devcontainer", "shell", flag)
			plain := sessionUseCase()
			plainOut, plainErr, plainErrV := shellRun(t, plain, "devcontainer", "shell")
			if err != nil || plainErrV != nil {
				t.Fatalf("err=%v plain=%v, want both to open the session", err, plainErrV)
			}
			if uc.shellCalls != 1 || out != plainOut || errOut != plainErr ||
				uc.shellReq.BaseDir != plain.shellReq.BaseDir || uc.shellReq.Terminal != plain.shellReq.Terminal {
				t.Errorf("%s changed the call: calls=%d stdout=%q/%q stderr=%q/%q", flag, uc.shellCalls, out, plainOut, errOut, plainErr)
			}
		})
	}
}

// LH-FA-CLI-006 / ADR-0031 decision 3: every sentinel of the start and of the
// session has its exit code, nothing is written to stdout, and the path of the
// project stays out of the message. The classes: 10 (no configuration, folder,
// variable name), 11 (Docker, remote engine, start script), 12 (runtime, which
// is also the failed check before the session), 14, and 1 for the rest. The
// class of the abort of a start script wins over the runtime class.
func TestExecute_DevcontainerShell_SentinelsAndExitCodes(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		err  error
		exit int
	}{
		{"docker unavailable", fmt.Errorf("docker binary not on PATH: %w", driven.ErrDockerUnavailable), 11},
		{"engine not local", driven.ErrContainerEngineNotLocal, 11},
		{"start script failed", fmt.Errorf("%w: /tmp/x/demo boom", driving.ErrDevcontainerStartScriptFailed), 11},
		{"runtime", fmt.Errorf("devcontainer up failed: %w", driven.ErrDevcontainerRuntime), 12},
		{"the check before the session fails", fmt.Errorf("devcontainer shell: no shell runs as user %q in %q of container x (bash: boom): %w",
			"vscode", "/tmp/x/demo/w", driven.ErrDevcontainerRuntime), 12},
		{"no configuration", fmt.Errorf("%w (looked for x)", driving.ErrDevcontainerConfigMissing), 10},
		{"folder", driving.ErrDevcontainerProjectFolderInvalid, 10},
		{"variable name", driving.ErrDevcontainerEnvNameInvalid, 10},
		{"file system", fmt.Errorf("read /tmp/x/demo/.devcontainer/devcontainer.json: %w", driving.ErrDevcontainerFileSystem), 14},
		{"other", errors.New("boom"), 1},
		{"start script and runtime", fmt.Errorf("%w: %w", driving.ErrDevcontainerStartScriptFailed, driven.ErrDevcontainerRuntime), 11},
		{"an aborted start", context.Canceled, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := cli.ExitCode(tc.err); got != tc.exit {
				t.Fatalf("ExitCode(%v) = %d, want %d", tc.err, got, tc.exit)
			}
			uc := &fakeDevcontainerUseCase{shellErr: tc.err}
			out, _, err := shellRun(t, uc, "devcontainer", "shell")
			if !errors.Is(err, tc.err) || cli.ExitCode(err) != tc.exit || out != "" {
				t.Fatalf("err=%v exit=%d stdout=%q, want %v, exit %d and no stdout", err, cli.ExitCode(err), out, tc.err, tc.exit)
			}
			if strings.Contains(err.Error(), "/tmp/x/demo") {
				t.Errorf("the project path leaks into the error: %v", err)
			}
		})
	}
}
