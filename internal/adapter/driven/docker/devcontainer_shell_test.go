package docker_test

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/u-boot/internal/adapter/driven/docker"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

const (
	shellTestID        = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	shellTestUser      = "vscode"
	shellTestWorkspace = "/workspaces/project"
)

// Messages of the docker client (Docker 29.9.2) that end `docker exec` before
// the session begins; measured in the Docker tests of this package.
const (
	msgNotRunning = "Error response from daemon: container " + shellTestID + " is not running"
	msgGone       = "Error response from daemon: No such container: " + shellTestID
	msgUserGone   = "Error response from daemon: unable to find user vscode: no matching entries in passwd file"
	msgEngineDown = "Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?"
	msgNoTerminal = "cannot attach stdin to a TTY-enabled container because stdin is not a terminal"
	msgOCI        = `OCI runtime exec failed: exec failed: unable to start container process: chdir to cwd ("/workspaces/project") set in config.json failed: no such file or directory`
	// Engine 28.0.4 reports a missing user on the stream of the session, with
	// status 126 and without the prefix of the daemon.
	msgUserGoneOnStdout = "unable to find user vscode: no matching entries in passwd file"
	// A `docker` command inside the session prints this about another container.
	msgOtherContainer = "Error response from daemon: No such container: some-other-container"
)

// shellOptions returns the options of a session without terminal and without
// streams in a valid container; edit changes them.
func shellOptions(edit func(*driven.DevcontainerShellOptions)) driven.DevcontainerShellOptions {
	opts := driven.DevcontainerShellOptions{ContainerID: shellTestID, User: shellTestUser, Workspace: shellTestWorkspace}
	if edit != nil {
		edit(&opts)
	}
	return opts
}

// execCalls returns the recorded `docker exec` argument lists in order.
func (f *fakeDocker) execCalls(t *testing.T) [][]string {
	t.Helper()
	var execs [][]string
	for _, c := range f.calls(t) {
		if len(c) > 0 && c[0] == "exec" {
			execs = append(execs, c)
		}
	}
	return execs
}

// sessionCall returns the `docker exec` of the session (the one with `-i`),
// or nil.
func (f *fakeDocker) sessionCall(t *testing.T) []string {
	t.Helper()
	for _, c := range f.execCalls(t) {
		if contains(c, "-i") {
			return c
		}
	}
	return nil
}

// ADR-0031 Entscheidung 1, Fitness Function: the argument list of the
// session. `-t` only with a terminal, user and workspace of the caller, the
// variables by name, the shell last and without an argument.
func TestShellArgs(t *testing.T) {
	id, user, ws := shellTestID, shellTestUser, shellTestWorkspace
	terminal := func(o *driven.DevcontainerShellOptions) { o.Terminal = true }
	cases := []struct {
		name  string
		edit  func(*driven.DevcontainerShellOptions)
		env   []string
		shell string
		want  []string
	}{
		{"no terminal, bash", nil, nil, "bash", []string{"exec", "-i", "-u", user, "-w", ws, id, "bash"}},
		{"no terminal, sh", nil, nil, "sh", []string{"exec", "-i", "-u", user, "-w", ws, id, "sh"}},
		{"terminal", terminal, nil, "bash", []string{"exec", "-i", "-t", "-u", user, "-w", ws, id, "bash"}},
		{"variables by name", nil, []string{"GIT_TOKEN", "OTHER"}, "bash",
			[]string{"exec", "-i", "-u", user, "-w", ws, "-e", "GIT_TOKEN", "-e", "OTHER", id, "bash"}},
		{"terminal with a variable", terminal, []string{"GIT_TOKEN"}, "sh",
			[]string{"exec", "-i", "-t", "-u", user, "-w", ws, "-e", "GIT_TOKEN", id, "sh"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := docker.ShellArgsForTest(shellOptions(tc.edit), tc.env, tc.shell)
			if strings.Join(got, "\x00") != strings.Join(tc.want, "\x00") {
				t.Errorf("args =\n%q\nwant\n%q", got, tc.want)
			}
			for _, arg := range got {
				switch {
				case arg == "-v", arg == "--volume", arg == "--mount", arg == "--privileged", arg == "--network",
					arg == "-d", arg == "--detach", strings.HasPrefix(arg, "--cap-add"), strings.HasPrefix(arg, "--user"):
					t.Errorf("argument %q is not part of the session contract", arg)
				case strings.Contains(arg, "="):
					t.Errorf("argument %q carries a value", arg)
				}
			}
			if last := got[len(got)-1]; last != tc.shell {
				t.Errorf("the shell is not the last argument: %q", got)
			}
		})
	}
}

// ADR-0031 Entscheidung 1, Shell: the probe is non-interactive and carries
// no variables.
func TestProbeArgs(t *testing.T) {
	got := docker.ProbeArgsForTest(shellOptions(func(o *driven.DevcontainerShellOptions) {
		o.Terminal = true
		o.EnvNames = []string{"GIT_TOKEN"}
	}), "bash")
	want := []string{"exec", "-u", shellTestUser, "-w", shellTestWorkspace, shellTestID, "bash", "-c", "true"}
	if strings.Join(got, "\x00") != strings.Join(want, "\x00") {
		t.Errorf("probe args =\n%q\nwant\n%q", got, want)
	}
}

// ADR-0031 Entscheidung 1: probe first, then the session with the streams of
// the caller. The launcher writes nothing of its own to standard output.
func TestLauncher_Shell_ProbesThenOpensTheSession(t *testing.T) {
	f := newFakeDocker(t)
	f.write(t, "session-stdout", "from-session\n")
	f.write(t, "session-stderr", "warn-session\n")
	f.flag(t, "session-copy-stdin")
	var stdout, stderr bytes.Buffer
	err := f.launcher(nil).Shell(context.Background(), shellOptions(func(o *driven.DevcontainerShellOptions) {
		o.Stdin = strings.NewReader("echo hi\n")
		o.Stdout, o.Stderr = &stdout, &stderr
	}))
	if err != nil {
		t.Fatalf("Shell: %v", err)
	}
	execs := f.execCalls(t)
	if len(execs) != 2 {
		t.Fatalf("docker exec ran %d times, want the probe and the session: %q", len(execs), execs)
	}
	base := "exec -u vscode -w /workspaces/project " + shellTestID
	if got := strings.Join(execs[0], " "); got != base+" bash -c true" {
		t.Errorf("probe = %q", got)
	}
	if got := strings.Join(execs[1], " "); got != "exec -i -u vscode -w /workspaces/project "+shellTestID+" bash" {
		t.Errorf("session = %q", got)
	}
	if stdout.String() != "from-session\n" || stderr.String() != "warn-session\n" {
		t.Errorf("stdout = %q, stderr = %q; want what the session wrote, nothing else", stdout.String(), stderr.String())
	}
	copied, err := os.ReadFile(filepath.Join(f.dir, "stdin.copy"))
	if err != nil || string(copied) != "echo hi\n" {
		t.Errorf("standard input of the session = %q (%v), want the stream of the caller", copied, err)
	}
}

// ADR-0031 Entscheidung 1, Shell: `sh` when `bash` does not run.
func TestLauncher_Shell_ChoosesShWhenBashIsMissing(t *testing.T) {
	f := newFakeDocker(t)
	f.flag(t, "bash-missing")
	if err := f.launcher(nil).Shell(context.Background(), shellOptions(nil)); err != nil {
		t.Fatalf("Shell: %v", err)
	}
	var shells []string
	for _, c := range f.execCalls(t) {
		shells = append(shells, strings.Join(c, " "))
	}
	want := []string{
		"exec -u vscode -w /workspaces/project " + shellTestID + " bash -c true",
		"exec -u vscode -w /workspaces/project " + shellTestID + " sh -c true",
		"exec -i -u vscode -w /workspaces/project " + shellTestID + " sh",
	}
	if strings.Join(shells, "\n") != strings.Join(want, "\n") {
		t.Errorf("docker exec calls =\n%s\nwant\n%s", strings.Join(shells, "\n"), strings.Join(want, "\n"))
	}
}

// ADR-0031 Entscheidung 3: a failed probe ends with ErrDevcontainerRuntime
// and no session.
func TestLauncher_Shell_FailedProbeIsRuntimeFailureWithoutSession(t *testing.T) {
	cases := []struct {
		name  string
		flags []string
		want  []string
	}{
		{"neither bash nor sh", []string{"bash-missing", "sh-missing"}, []string{`bash: OCI runtime exec failed`, `sh: OCI runtime exec failed`}},
		{"user missing", []string{"user-missing"}, []string{"bash: Error response from daemon: unable to find user", "sh: Error response from daemon: unable to find user"}},
		{"workspace missing", []string{"workspace-missing"}, []string{"chdir to cwd"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newFakeDocker(t)
			f.flag(t, tc.flags...)
			var stdout bytes.Buffer
			err := f.launcher(nil).Shell(context.Background(), shellOptions(func(o *driven.DevcontainerShellOptions) { o.Stdout = &stdout }))
			if !errors.Is(err, driven.ErrDevcontainerRuntime) || errors.Is(err, driven.ErrDockerUnavailable) {
				t.Fatalf("error = %v, want ErrDevcontainerRuntime", err)
			}
			for _, text := range tc.want {
				if !strings.Contains(err.Error(), text) {
					t.Errorf("error %q lacks %q", err, text)
				}
			}
			if f.sessionCall(t) != nil {
				t.Errorf("a session started although the probe failed: %q", f.sessionCall(t))
			}
			if got := len(f.execCalls(t)); got != 2 {
				t.Errorf("docker exec ran %d times, want the probe of bash and of sh", got)
			}
			if stdout.Len() != 0 {
				t.Errorf("standard output = %q, want nothing", stdout.String())
			}
		})
	}
}

// ADR-0031 Entscheidung 3, Grenze: what happens in the session does not
// change the result. Every exit status of the shell is no error, also the
// status 1 of a failed last command, 126, 127 and 137, and also when the
// session writes the same text a client error has (guards: the same exit
// status and the same message in another context are no error).
func TestLauncher_Shell_SessionEndIsNoError(t *testing.T) {
	terminal := func(o *driven.DevcontainerShellOptions) { o.Terminal = true }
	cases := []struct {
		name   string
		exit   string
		stdout string
		stderr string
		flags  []string
		edit   func(*driven.DevcontainerShellOptions)
	}{
		{name: "status 0", exit: "0"},
		{name: "status 0, the message names this container", exit: "0", stderr: msgNotRunning + "\n"},
		{name: "status 1, silent", exit: "1"},
		{name: "status 2", exit: "2"},
		{name: "status 126", exit: "126"},
		{name: "status 127, silent", exit: "127"},
		{name: "status 130", exit: "130"},
		{name: "status 137, the container ended", exit: "137"},
		{name: "status 1, a docker command in the session failed about another container", exit: "1", stderr: msgOtherContainer + "\n"},
		{name: "status 1, the message names this container but the session wrote to standard output", exit: "1", stdout: "output\n", stderr: msgNotRunning + "\n"},
		{name: "status 1, other text before the message", exit: "1", stderr: "warning\n" + msgNotRunning + "\n"},
		{name: "status 1, the engine message while the engine answers", exit: "1", stderr: msgEngineDown + "\n"},
		{name: "status 1, the user message of another user", exit: "1", stderr: "Error response from daemon: unable to find user other: no matching entries in passwd file\n"},
		{name: "status 127, the OCI line after other output", exit: "127", stdout: "output\n" + msgOCI + "\n"},
		{name: "status 127, the OCI line on standard error", exit: "127", stderr: msgOCI + "\n"},
		{name: "status 127, the OCI line and other text on standard error", exit: "127", stdout: msgOCI + "\n", stderr: "warning\n"},
		{name: "status 127, other text on standard output", exit: "127", stdout: "bash: foo: command not found\n"},
		{name: "status 1, the OCI line as the only output", exit: "1", stdout: msgOCI + "\n"},
		{name: "status 126, the missing-user line of another user", exit: "126", stdout: "unable to find user other: no matching entries in passwd file\r\n"},
		{name: "status 1, the missing-user line as the only output", exit: "1", stdout: msgUserGoneOnStdout + "\r\n"},
		{name: "status 126, the missing-user line and text on standard error", exit: "126", stdout: msgUserGoneOnStdout + "\r\n", stderr: "warning\n"},
		{name: "status 126, other output before the missing-user line", exit: "126", stdout: "output\n" + msgUserGoneOnStdout + "\n"},
		{name: "status 126, other output after the missing-user line", exit: "126", stdout: msgUserGoneOnStdout + "\nmore output\n"},
		{name: "terminal, status 126, the missing-user line", exit: "126", stdout: msgUserGoneOnStdout + "\r\n", edit: terminal},
		{name: "status 126, the missing-user line of a user whose name begins with the session user", exit: "126", stdout: "unable to find user vscode2: no matching entries in passwd file\r\n"},
		{name: "status 126, the missing-user line of the session user for a session as vscode2", exit: "126", stdout: msgUserGoneOnStdout + "\r\n",
			edit: func(o *driven.DevcontainerShellOptions) { o.User = "vscode2" }},
		{name: "status 1, the daemon message about a user whose name begins with the session user", exit: "1",
			stderr: "Error response from daemon: unable to find user vscode2: no matching entries in passwd file\n"},
		{name: "status 126, the missing-user line followed by more blank output than the head holds", exit: "126",
			stdout: msgUserGoneOnStdout + "\n" + strings.Repeat("\n", 5000)},
		{name: "terminal, status 127, the OCI line", exit: "127", stdout: msgOCI + "\n", edit: terminal},
		{name: "status 1, the terminal message without a terminal", exit: "1", stderr: msgNoTerminal + "\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newFakeDocker(t)
			f.write(t, "session-exit", tc.exit)
			if tc.stdout != "" {
				f.write(t, "session-stdout", tc.stdout)
			}
			if tc.stderr != "" {
				f.write(t, "session-stderr", tc.stderr)
			}
			f.flag(t, tc.flags...)
			var stdout, stderr bytes.Buffer
			err := f.launcher(nil).Shell(context.Background(), shellOptions(func(o *driven.DevcontainerShellOptions) {
				if tc.edit != nil {
					tc.edit(o)
				}
				o.Stdout, o.Stderr = &stdout, &stderr
			}))
			if err != nil {
				t.Errorf("error = %v, want none after the session began", err)
			}
			if stdout.String() != tc.stdout || stderr.String() != tc.stderr {
				t.Errorf("stdout = %q, stderr = %q; want the session output unchanged", stdout.String(), stderr.String())
			}
		})
	}
}

// LH-FA-CLI-006, ADR-0031 Entscheidung 3: an error of the docker client
// before the session began is ErrDevcontainerRuntime. The exit status cannot
// tell it from the status of a shell, the message of the client can.
func TestLauncher_Shell_ClientErrorBeforeSessionIsRuntimeFailure(t *testing.T) {
	terminal := func(o *driven.DevcontainerShellOptions) { o.Terminal = true }
	cases := []struct {
		name   string
		exit   string
		stdout string
		stderr string
		flags  []string
		edit   func(*driven.DevcontainerShellOptions)
	}{
		{name: "the container is not running", exit: "1", stderr: msgNotRunning + "\n"},
		{name: "the container does not exist", exit: "1", stderr: msgGone + "\n"},
		{name: "the user does not exist after the probe", exit: "1", stderr: msgUserGone + "\n"},
		{name: "the engine does not answer after the probe", exit: "1", stderr: msgEngineDown + "\n", flags: []string{"daemon-down-late"}},
		{name: "the workspace does not exist after the probe", exit: "127", stdout: msgOCI + "\n"},
		{name: "the user does not exist after the probe (engine 28.0.4: status 126 on standard output)", exit: "126", stdout: msgUserGoneOnStdout + "\r\n"},
		{name: "the user vscode2 does not exist after the probe (engine 28.0.4)", exit: "126", stdout: "unable to find user vscode2: no matching entries in passwd file\r\n",
			edit: func(o *driven.DevcontainerShellOptions) { o.User = "vscode2" }},
		{name: "terminal, the container is not running", exit: "1", stderr: msgNotRunning + "\n", edit: terminal},
		{name: "terminal, standard input is no terminal", exit: "1", stderr: msgNoTerminal + "\n", edit: terminal},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newFakeDocker(t)
			f.write(t, "session-exit", tc.exit)
			if tc.stdout != "" {
				f.write(t, "session-stdout", tc.stdout)
			}
			f.write(t, "session-stderr", tc.stderr)
			f.flag(t, tc.flags...)
			err := f.launcher(nil).Shell(context.Background(), shellOptions(tc.edit))
			if !errors.Is(err, driven.ErrDevcontainerRuntime) || errors.Is(err, driven.ErrDockerUnavailable) {
				t.Fatalf("error = %v, want ErrDevcontainerRuntime", err)
			}
			message := strings.TrimSpace(tc.stderr + tc.stdout)
			if !strings.Contains(err.Error(), "before the session began") || !strings.Contains(err.Error(), message) {
				t.Errorf("error %q lacks the message of the client %q", err, message)
			}
		})
	}
}

// ADR-0031 Entscheidung 4: without a terminal, a cancelled context
// interrupts the client and ends the session with the context error.
func TestLauncher_Shell_CancelWithoutTerminalEndsWithTheContextError(t *testing.T) {
	f := newFakeDocker(t)
	f.flag(t, "session-hang")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- f.launcher(nil).Shell(ctx, shellOptions(nil)) }()
	waitForSession(t, f)
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) || errors.Is(err, driven.ErrDevcontainerRuntime) {
			t.Errorf("error = %v, want the context error", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("Shell did not return after cancel: the client was not interrupted (the fake session runs 30 s)")
	}
}

// ADR-0034 Entscheidung 1 (Entscheidung 4 von ADR-0031): without a terminal a
// cancelled context sends SIGINT to the client, not SIGKILL. The fake client
// records the signal it can catch; a kill leaves no record.
func TestLauncher_Shell_CancelWithoutTerminalSendsSIGINTToTheClient(t *testing.T) {
	f := newFakeDocker(t)
	f.flag(t, "session-trap")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- f.launcher(nil).Shell(ctx, shellOptions(nil)) }()
	ready := filepath.Join(f.dir, "trap-ready")
	deadline := time.Now().Add(20 * time.Second)
	for {
		if _, err := os.Stat(ready); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("the fake client did not install its signal handlers")
		}
		time.Sleep(20 * time.Millisecond)
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Errorf("error = %v, want the context error", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("Shell did not return after cancel")
	}
	got, _ := os.ReadFile(filepath.Join(f.dir, "signal"))
	if strings.TrimSpace(string(got)) != "INT" {
		t.Errorf("the client received %q, want the signal INT", got)
	}
}

// ADR-0031 Entscheidung 4: with a terminal the context is not observed once
// the session has begun; the session runs to its end.
func TestLauncher_Shell_TerminalSessionIgnoresTheContext(t *testing.T) {
	f := newFakeDocker(t)
	f.flag(t, "session-delay")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var stdout bytes.Buffer
	done := make(chan error, 1)
	go func() {
		done <- f.launcher(nil).Shell(ctx, shellOptions(func(o *driven.DevcontainerShellOptions) {
			o.Terminal = true
			o.Stdout = &stdout
		}))
	}()
	waitForSession(t, f)
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("error = %v, want none: the session ran to its end", err)
		}
		if got := stdout.String(); got != "session-done\n" {
			t.Errorf("standard output = %q, want the output of the finished session", got)
		}
	case <-time.After(30 * time.Second):
		t.Fatal("Shell did not return")
	}
}

// ADR-0031 Entscheidung 4: with a terminal the client gets the standard input
// and output of the caller as they are, because it looks for the terminal on
// them; without a terminal the launcher records standard output on its way
// (guard: the same probe sees a pipe there).
func TestLauncher_Shell_TerminalHandsTheStandardStreamsToTheClient(t *testing.T) {
	dir := t.TempDir()
	stdin, err := os.Create(filepath.Join(dir, "stdin"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = stdin.Close() })
	stdout, err := os.Create(filepath.Join(dir, "stdout"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = stdout.Close() })
	for _, terminal := range []bool{true, false} {
		f := newFakeDocker(t)
		f.flag(t, "session-record-fds")
		err := f.launcher(nil).Shell(context.Background(), shellOptions(func(o *driven.DevcontainerShellOptions) {
			o.Terminal = terminal
			o.Stdin, o.Stdout = stdin, stdout
		}))
		if err != nil {
			t.Fatalf("Terminal=%v: Shell: %v", terminal, err)
		}
		fd0, _ := os.ReadFile(filepath.Join(f.dir, "fd0.link"))
		fd1, _ := os.ReadFile(filepath.Join(f.dir, "fd1.link"))
		gotIn, gotOut := strings.TrimSpace(string(fd0)), strings.TrimSpace(string(fd1))
		if gotIn != stdin.Name() {
			t.Errorf("Terminal=%v: standard input of the client = %q, want the file of the caller %q", terminal, gotIn, stdin.Name())
		}
		if terminal && gotOut != stdout.Name() {
			t.Errorf("Terminal: standard output of the client = %q, want the file of the caller %q", gotOut, stdout.Name())
		}
		if !terminal && !strings.HasPrefix(gotOut, "pipe:") {
			t.Errorf("no Terminal: standard output of the client = %q, want a pipe; the probe cannot tell the two apart", gotOut)
		}
	}
}

// waitForSession blocks until the fake docker has seen the session.
func waitForSession(t *testing.T, f *fakeDocker) {
	t.Helper()
	deadline := time.Now().Add(20 * time.Second)
	for f.sessionCall(t) == nil {
		if time.Now().After(deadline) {
			t.Fatal("the session did not start")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// LH-FA-DEV-009: the value of a passed variable is on no command line, the
// probe carries no variable, and a name that is not set on the host is not
// passed.
func TestLauncher_Shell_PassesVariablesByNameOnly(t *testing.T) {
	f := newFakeDocker(t)
	const secret = "placeholder-token-value"
	env := map[string]string{"GIT_TOKEN": secret, "UNLISTED": "placeholder-unlisted"}
	err := f.launcher(env).Shell(context.Background(), shellOptions(func(o *driven.DevcontainerShellOptions) {
		o.EnvNames = []string{"GIT_TOKEN", "NOT_SET"}
	}))
	if err != nil {
		t.Fatalf("Shell: %v", err)
	}
	joined := strings.Join(f.flatCalls(t), "\n")
	for _, leaked := range []string{secret, "placeholder-unlisted", "NOT_SET"} {
		if strings.Contains(joined, leaked) {
			t.Errorf("%q reached a command line:\n%s", leaked, joined)
		}
	}
	var passed []string
	session := f.sessionCall(t)
	for i, arg := range session {
		if arg == "-e" {
			passed = append(passed, session[i+1])
		}
	}
	if strings.Join(passed, ",") != "GIT_TOKEN" {
		t.Errorf("passed variables = %v, want [GIT_TOKEN]", passed)
	}
	for _, c := range f.execCalls(t) {
		if !contains(c, "-i") && contains(c, "-e") {
			t.Errorf("the probe carries a variable: %q", c)
		}
	}
}

// ADR-0031, akzeptiertes Negativ 1: the session is not masked; what the
// shell prints reaches the caller as it is.
func TestLauncher_Shell_MasksNothingInTheSession(t *testing.T) {
	f := newFakeDocker(t)
	const secret = "placeholder-token-value"
	f.write(t, "session-stdout", "token="+secret+"\n")
	var stdout bytes.Buffer
	err := f.launcher(map[string]string{"GIT_TOKEN": secret}).Shell(context.Background(), shellOptions(func(o *driven.DevcontainerShellOptions) {
		o.EnvNames = []string{"GIT_TOKEN"}
		o.Stdout = &stdout
	}))
	if err != nil {
		t.Fatalf("Shell: %v", err)
	}
	if stdout.String() != "token="+secret+"\n" {
		t.Errorf("standard output = %q, want the text of the session unchanged", stdout.String())
	}
}

// The launcher refuses an invalid argument before any docker call.
func TestLauncher_Shell_RejectsUnsafeInput(t *testing.T) {
	cases := []struct {
		name string
		edit func(*driven.DevcontainerShellOptions)
	}{
		{"container ID empty", func(o *driven.DevcontainerShellOptions) { o.ContainerID = "" }},
		{"container ID too short", func(o *driven.DevcontainerShellOptions) { o.ContainerID = strings.Repeat("a", 11) }},
		{"container ID is an option", func(o *driven.DevcontainerShellOptions) { o.ContainerID = "--all" }},
		{"container ID upper case", func(o *driven.DevcontainerShellOptions) { o.ContainerID = strings.Repeat("A", 12) }},
		{"user empty", func(o *driven.DevcontainerShellOptions) { o.User = "" }},
		{"user is an option", func(o *driven.DevcontainerShellOptions) { o.User = "--privileged" }},
		{"workspace relative", func(o *driven.DevcontainerShellOptions) { o.Workspace = "workspaces/project" }},
		{"workspace empty", func(o *driven.DevcontainerShellOptions) { o.Workspace = "" }},
		{"variable name with a value", func(o *driven.DevcontainerShellOptions) { o.EnvNames = []string{"A=x"} }},
		{"variable name empty", func(o *driven.DevcontainerShellOptions) { o.EnvNames = []string{""} }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newFakeDocker(t)
			err := f.launcher(map[string]string{"A": "x"}).Shell(context.Background(), shellOptions(tc.edit))
			if !errors.Is(err, driven.ErrDevcontainerRuntime) {
				t.Errorf("error = %v, want ErrDevcontainerRuntime", err)
			}
			if calls := f.flatCalls(t); len(calls) != 0 {
				t.Errorf("docker ran for rejected input: %v", calls)
			}
		})
	}
}
