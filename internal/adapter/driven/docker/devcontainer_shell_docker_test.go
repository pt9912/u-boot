//go:build docker

// Integration tests for [docker.Launcher.Shell] against the Docker engine of
// the runner (`make test-docker`): user, workspace and variables of the
// session, the probe, the end of the client, the terminal path under
// `script` and the errors of the docker client before the session begins.
// The containers are plain containers of an image that sleeps, started by the
// test and removed with it; the session does not depend on the start of the
// devcontainer. A missing docker binary or a missing `script` fails the test,
// it does not skip it.

package docker_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/pt9912/u-boot/internal/adapter/driven/docker"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

const (
	shellImageNoBash  = "busybox:1.36"
	shellWorkspaceDir = "/workspaces/ubl"
	shellHelperEnv    = "UBL_SHELL_HELPER_MODE"
)

// TestMain turns the test binary into the helper of the terminal tests when
// the environment names a mode: the helper opens one session with the
// standard streams it was given (under `script` they are a pseudo-terminal)
// and writes its result to a file. Otherwise it runs the tests.
func TestMain(m *testing.M) {
	if mode := os.Getenv(shellHelperEnv); mode != "" {
		os.Exit(runShellHelper(mode))
	}
	os.Exit(m.Run())
}

func runShellHelper(mode string) int {
	err := docker.NewLauncher().Shell(context.Background(), driven.DevcontainerShellOptions{
		ContainerID: os.Getenv("UBL_SHELL_HELPER_ID"),
		User:        os.Getenv("UBL_SHELL_HELPER_USER"),
		Workspace:   os.Getenv("UBL_SHELL_HELPER_WORKSPACE"),
		Stdin:       os.Stdin,
		Stdout:      os.Stdout,
		Stderr:      os.Stderr,
		Terminal:    mode == "terminal",
	})
	result := "ok"
	if err != nil {
		result = "error: " + err.Error()
	}
	if werr := os.WriteFile(os.Getenv("UBL_SHELL_HELPER_RESULT"), []byte(result), 0o644); werr != nil || err != nil {
		return 1
	}
	return 0
}

// newShellContainer starts a container of the image that sleeps and removes
// it with the test. It returns the full container ID.
func newShellContainer(t *testing.T, image string) string {
	t.Helper()
	requireDocker(t)
	name := "ubl-shelltest-" + randomHex(t)
	t.Cleanup(func() { _, _ = tryDocker("rm", "-f", name) })
	return mustDocker(t, "run", "-d", "--name", name, image, "sleep", "3600")
}

// newVscodeContainer starts a container of the devcontainer base image with
// the workspace directory owned by `vscode`.
func newVscodeContainer(t *testing.T) string {
	t.Helper()
	id := newShellContainer(t, devBaseImage)
	mustDocker(t, "exec", "-u", "root", id, "sh", "-c", "mkdir -p "+shellWorkspaceDir+" && chown vscode:vscode "+shellWorkspaceDir)
	return id
}

func vscodeShell(id string, edit func(*driven.DevcontainerShellOptions)) driven.DevcontainerShellOptions {
	opts := driven.DevcontainerShellOptions{ContainerID: id, User: "vscode", Workspace: shellWorkspaceDir}
	if edit != nil {
		edit(&opts)
	}
	return opts
}

// runShell opens one session of the launcher with the script on its standard
// input and returns what the session wrote.
func runShell(t *testing.T, l *docker.Launcher, opts driven.DevcontainerShellOptions, script string) (stdout, stderr string, err error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), launcherTestBudget)
	defer cancel()
	var out, errBuf bytes.Buffer
	opts.Stdin = strings.NewReader(script)
	opts.Stdout, opts.Stderr = &out, &errBuf
	err = l.Shell(ctx, opts)
	return out.String(), errBuf.String(), err
}

// containerProcs lists the PIDs of the processes called comm in the container
// (PID 1 excluded).
func containerProcs(t *testing.T, id, comm string) []string {
	t.Helper()
	script := `for p in /proc/[0-9]*; do i=${p#/proc/}; [ "$i" = 1 ] && continue; [ "$(cat $p/comm 2>/dev/null)" = ` + comm + ` ] && echo $i; done; true`
	return strings.Fields(mustDocker(t, "exec", "-u", "root", id, "sh", "-c", script))
}

// waitProcs polls until the number of processes called comm satisfies done or
// the time is up, and returns the last list.
func waitProcs(t *testing.T, id, comm string, within time.Duration, done func(n int) bool) []string {
	t.Helper()
	deadline := time.Now().Add(within)
	for {
		procs := containerProcs(t, id, comm)
		if done(len(procs)) || time.Now().After(deadline) {
			return procs
		}
		time.Sleep(250 * time.Millisecond)
	}
}

// ADR-0031 Entscheidung 1 (Sitzungsphase, Umgebung, Kein Socket während der
// Sitzung): the session runs as the given user in the given workspace, a
// variable that is set on the host arrives with its value, one that is not
// set is absent, the value is on no command line, and no helper container
// exists during the session. Guards: the same `docker exec` without `-e`
// does not carry the variable, `docker exec -e` leaves an unset variable out
// (not empty), the command-line scan finds a value in an argument, and the
// leftover check finds a container with the helper prefix.
func TestShell_Docker_NoTerminal_UserWorkspaceAndEnvironment(t *testing.T) {
	id := newVscodeContainer(t)
	suffix := randomHex(t)
	setName, unsetName := "UBL_SHELL_SET_"+suffix, "UBL_SHELL_UNSET_"+suffix
	value := "placeholder-shell-" + suffix
	t.Setenv(setName, value)
	_ = os.Unsetenv(unsetName)

	script := "id -un\npwd\nprintenv " + setName + "\nprintenv " + unsetName + " || echo unset-absent\necho \"bash=${BASH_VERSION:+yes}\"\nsleep 3\n"
	scanCtx, stopScan := context.WithCancel(context.Background())
	defer stopScan()
	cmdHits := scanDuring(scanCtx, func() []string { return cmdlineHolders(value) })
	helperHits := scanDuring(scanCtx, func() []string { return helperLeftovers(t) })

	l := testLauncher(t)
	stdout, stderr, err := runShell(t, l, vscodeShell(id, func(o *driven.DevcontainerShellOptions) {
		o.EnvNames = []string{setName, unsetName}
	}), script)
	stopScan()
	if err != nil {
		t.Fatalf("Shell: %v\nstderr: %s", err, stderr)
	}
	if want := "vscode\n" + shellWorkspaceDir + "\n" + value + "\nunset-absent\nbash=yes\n"; stdout != want {
		t.Errorf("standard output =\n%q\nwant\n%q", stdout, want)
	}
	if stderr != "" {
		t.Errorf("standard error = %q, want nothing", stderr)
	}
	if got := <-cmdHits; len(got) != 0 {
		t.Errorf("the value was on a command line: %q", got)
	}
	if got := <-helperHits; len(got) != 0 {
		t.Errorf("a helper container existed during the session: %v", got)
	}

	t.Run("guard: docker exec without -e does not carry the variable", func(t *testing.T) {
		out, err := tryDocker("exec", "-u", "vscode", "-w", shellWorkspaceDir, id, "printenv", setName)
		if err == nil || out != "" {
			t.Errorf("the variable is in the container without `-e`: %q (%v); the test proves nothing", out, err)
		}
	})
	t.Run("guard: docker exec -e leaves an unset variable out", func(t *testing.T) {
		out, err := tryDocker("exec", "-e", unsetName, "-u", "vscode", id, "printenv", unsetName)
		if err == nil {
			t.Errorf("`docker exec -e %s` set a variable that is unset on the host (value %q): the ADR assumption is wrong", unsetName, out)
		}
	})
	t.Run("guard: the command-line scan finds a value in an argument", func(t *testing.T) {
		guardValue := "placeholder-guard-" + suffix
		cmd := exec.Command("docker", "exec", "-i", "-e", "GUARD="+guardValue, "-u", "vscode", id, "sleep", "20")
		cmd.Stdin = strings.NewReader("")
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		defer func() { _ = cmd.Process.Kill(); _ = cmd.Wait() }()
		deadline := time.Now().Add(20 * time.Second)
		for len(cmdlineHolders(guardValue)) == 0 {
			if time.Now().After(deadline) {
				t.Fatal("the command-line scan did not find a value passed as `-e NAME=value`")
			}
			time.Sleep(20 * time.Millisecond)
		}
	})
	t.Run("guard: the leftover check finds a container with the helper prefix", func(t *testing.T) {
		name := helperPrefix(t) + "guard-" + suffix
		t.Cleanup(func() { removeContainers("name=" + name) })
		mustDocker(t, "run", "-d", "--name", name, shellImageNoBash, "sleep", "30")
		if left := helperLeftovers(t); len(left) != 1 {
			t.Errorf("the leftover check found %v, want the one container started by hand", left)
		}
	})
}

// ADR-0031 Entscheidung 1 (Shell) und 3: the probe finds user, workspace and
// shell, or the session does not start.
func TestShell_Docker_Probe(t *testing.T) {
	const marker = "/tmp/ubl-session-ran"
	noSession := func(t *testing.T, id string) {
		t.Helper()
		if out, err := tryDocker("exec", "-u", "root", id, "ls", marker); err == nil {
			t.Errorf("the session ran although the probe failed: %q", out)
		}
	}
	failing := func(t *testing.T, id string, edit func(*driven.DevcontainerShellOptions), wantText string) {
		t.Helper()
		stdout, _, err := runShell(t, testLauncher(t), vscodeShell(id, edit), "touch "+marker+"\n")
		if !errors.Is(err, driven.ErrDevcontainerRuntime) || errors.Is(err, driven.ErrDockerUnavailable) {
			t.Fatalf("error = %v, want ErrDevcontainerRuntime", err)
		}
		if !strings.Contains(err.Error(), wantText) {
			t.Errorf("error %q lacks %q", err, wantText)
		}
		if stdout != "" {
			t.Errorf("standard output = %q, want nothing", stdout)
		}
		noSession(t, id)
	}

	t.Run("user missing", func(t *testing.T) {
		id := newVscodeContainer(t)
		failing(t, id, func(o *driven.DevcontainerShellOptions) { o.User = "nosuchuser" }, "unable to find user nosuchuser")
	})
	t.Run("workspace missing", func(t *testing.T) {
		id := newVscodeContainer(t)
		failing(t, id, func(o *driven.DevcontainerShellOptions) { o.Workspace = "/workspaces/nonexistent" }, "chdir to cwd")
	})
	t.Run("container not running", func(t *testing.T) {
		id := newVscodeContainer(t)
		mustDocker(t, "kill", id)
		failing(t, id, nil, "is not running")
	})
	t.Run("image without bash chooses sh", func(t *testing.T) {
		id := newShellContainer(t, shellImageNoBash)
		l := testLauncher(t)
		stdout, stderr, err := runShell(t, l, driven.DevcontainerShellOptions{ContainerID: id, User: "root", Workspace: "/tmp"},
			"echo \"bash=${BASH_VERSION:-none}\"\necho \"arg0=$0\"\nid -un\npwd\n")
		if err != nil {
			t.Fatalf("Shell: %v\n%s", err, stderr)
		}
		if want := "bash=none\narg0=sh\nroot\n/tmp\n"; stdout != want {
			t.Errorf("standard output = %q, want %q", stdout, want)
		}
	})
	t.Run("image without any shell", func(t *testing.T) {
		id := newShellContainer(t, shellImageNoBash)
		mustDocker(t, "exec", "-u", "root", id, "rm", "-f", "/bin/sh")
		stdout, _, err := runShell(t, testLauncher(t), driven.DevcontainerShellOptions{ContainerID: id, User: "root", Workspace: "/tmp"}, "touch "+marker+"\n")
		if !errors.Is(err, driven.ErrDevcontainerRuntime) {
			t.Fatalf("error = %v, want ErrDevcontainerRuntime", err)
		}
		for _, want := range []string{`bash: OCI runtime exec failed`, `sh: OCI runtime exec failed`, "executable file not found"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("error %q lacks %q", err, want)
			}
		}
		if stdout != "" {
			t.Errorf("standard output = %q, want nothing", stdout)
		}
		noSession(t, id)
	})
}

// shellQuote quotes s for `sh -c`.
func shellQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'" }

// shimDocker writes a `docker` script that passes every call to the real
// docker. Before the call that opens the session (`exec -i …`) it runs action
// ($REAL is the real docker); a file `down` in the returned directory makes
// every later call fail to reach the engine. The probe has already ended
// when the action runs, so the action is the event "between the probe and the
// session" at a point the test controls.
func shimDocker(t *testing.T, action string) (binary, dir string) {
	t.Helper()
	dockerPath, err := exec.LookPath("docker")
	if err != nil {
		t.Fatalf("docker CLI not on PATH: %v", err)
	}
	dir = t.TempDir()
	script := "#!/bin/sh\nREAL=" + shellQuote(dockerPath) + "\nD=" + shellQuote(dir) + "\n" +
		"case \"$1 $*\" in\n\"exec \"*\" -i \"*) " + action + " ;;\nesac\n" +
		"if [ -e \"$D/down\" ]; then DOCKER_HOST=tcp://127.0.0.1:1; export DOCKER_HOST; fi\n" +
		"exec \"$REAL\" \"$@\"\n"
	binary = filepath.Join(dir, "docker")
	if err := os.WriteFile(binary, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return binary, dir
}

// LH-FA-CLI-006, ADR-0031 Entscheidung 3, B18: the probe succeeds and the
// session does not begin. The events happen between the probe and the session
// (a script in place of `docker` runs them before the session's `docker
// exec`); the real client ends with its own error, and the launcher returns
// ErrDevcontainerRuntime for every case it can tell from the exit status of a
// shell. Guards: a shell that ends with the same status 1, and a container
// that ends during the session, are no error.
func TestShell_Docker_B18_SessionDoesNotBegin(t *testing.T) {
	cases := []struct {
		name   string
		setup  func(t *testing.T, id string)
		action func(id string) string
		want   string
		wantID bool
	}{
		{"container killed", nil, func(id string) string { return `"$REAL" kill ` + id + ` >/dev/null` }, "is not running", true},
		{"container removed", nil, func(id string) string { return `"$REAL" rm -f ` + id + ` >/dev/null` }, "No such container", true},
		{"engine unreachable", nil, func(string) string { return `touch "$D/down"` }, "Cannot connect to the Docker daemon", false},
		{"user removed", func(t *testing.T, id string) {
			mustDocker(t, "exec", "-u", "root", id, "useradd", "-m", "ubluser")
		}, func(id string) string { return `"$REAL" exec -u root ` + id + ` userdel ubluser >/dev/null` }, "unable to find user ubluser", false},
		{"workspace removed", nil, func(id string) string {
			return `"$REAL" exec -u root ` + id + ` rm -rf ` + shellWorkspaceDir + ` >/dev/null`
		}, "chdir to cwd", false},
		{"shell removed", nil, func(id string) string {
			return `"$REAL" exec -u root ` + id + ` rm -f /usr/bin/bash /bin/bash >/dev/null`
		}, `"bash": executable file not found`, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			id := newVscodeContainer(t)
			user := "vscode"
			if tc.setup != nil {
				tc.setup(t, id)
				user = "ubluser"
			}
			binary, _ := shimDocker(t, tc.action(id))
			l := docker.NewLauncherForTest(binary, os.LookupEnv, 5*time.Second)
			var stdout, stderr bytes.Buffer
			ctx, cancel := context.WithTimeout(context.Background(), launcherTestBudget)
			defer cancel()
			err := l.Shell(ctx, driven.DevcontainerShellOptions{
				ContainerID: id, User: user, Workspace: shellWorkspaceDir,
				Stdin: strings.NewReader("echo session-ran\n"), Stdout: &stdout, Stderr: &stderr,
			})
			t.Logf("measured: error = %v\n  client stderr = %q\n  client stdout = %q", err, stderr.String(), stdout.String())
			if !errors.Is(err, driven.ErrDevcontainerRuntime) || errors.Is(err, driven.ErrDockerUnavailable) {
				t.Fatalf("error = %v, want ErrDevcontainerRuntime", err)
			}
			if !strings.Contains(err.Error(), "before the session began") || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("error %q lacks %q", err, tc.want)
			}
			if tc.wantID && !strings.Contains(err.Error(), id) {
				t.Errorf("error %q does not name the container", err)
			}
			if strings.Contains(stdout.String(), "session-ran") {
				t.Errorf("the session ran: %q", stdout.String())
			}
		})
	}

	t.Run("guard: the shell ends with status 1 in a running container", func(t *testing.T) {
		id := newVscodeContainer(t)
		stdout, stderr, err := runShell(t, testLauncher(t), vscodeShell(id, nil), "echo session-ran\nexit 1\n")
		if err != nil || stdout != "session-ran\n" || stderr != "" {
			t.Errorf("error = %v, stdout = %q, stderr = %q; want the session without an error", err, stdout, stderr)
		}
	})
	t.Run("the container ends during the session", func(t *testing.T) {
		id := newVscodeContainer(t)
		in, out, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = in.Close(); _ = out.Close() })
		var stdout, stderr syncBuffer
		ctx, cancel := context.WithTimeout(context.Background(), launcherTestBudget)
		defer cancel()
		done := make(chan error, 1)
		go func() {
			done <- testLauncher(t).Shell(ctx, vscodeShell(id, func(o *driven.DevcontainerShellOptions) {
				o.Stdin, o.Stdout, o.Stderr = in, &stdout, &stderr
			}))
		}()
		if _, err := out.WriteString("echo started\nsleep 60\n"); err != nil {
			t.Fatal(err)
		}
		deadline := time.Now().Add(30 * time.Second)
		for !strings.Contains(stdout.String(), "started") {
			if time.Now().After(deadline) {
				t.Fatal("the session did not start")
			}
			time.Sleep(50 * time.Millisecond)
		}
		mustDocker(t, "kill", id)
		select {
		case err := <-done:
			if err != nil {
				t.Errorf("error = %v, want none: the container ended during the session (stderr %q)", err, stderr.String())
			}
		case <-time.After(60 * time.Second):
			t.Fatal("Shell did not return after the container ended")
		}
	})
}

// syncBuffer is a buffer the writer goroutine of os/exec and the test read
// and write concurrently.
type syncBuffer struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (s *syncBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.Write(p)
}

func (s *syncBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.String()
}

// scriptSession is the test binary as helper under `script`: the helper's
// standard streams are a pseudo-terminal; the test types into the standard
// input of `script` and reads its output.
type scriptSession struct {
	cmd    *exec.Cmd
	in     io.WriteCloser
	out    *syncBuffer
	result string
	// exited is closed when `script` has ended; waitErr is set before.
	exited  chan struct{}
	waitErr error
}

func startUnderScript(t *testing.T, mode, id string) *scriptSession {
	t.Helper()
	requireDocker(t)
	if _, err := exec.LookPath("script"); err != nil {
		t.Fatalf("`script` not on PATH but the terminal test needs it: %v", err)
	}
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	s := &scriptSession{out: &syncBuffer{}, result: filepath.Join(t.TempDir(), "result"), exited: make(chan struct{})}
	s.cmd = exec.Command("script", "-q", "-e", "-c", shellQuote(self), "/dev/null")
	s.cmd.Env = append(os.Environ(),
		shellHelperEnv+"="+mode,
		"UBL_SHELL_HELPER_ID="+id,
		"UBL_SHELL_HELPER_USER=vscode",
		"UBL_SHELL_HELPER_WORKSPACE="+shellWorkspaceDir,
		"UBL_SHELL_HELPER_RESULT="+s.result,
	)
	s.cmd.Stdout, s.cmd.Stderr = s.out, s.out
	if s.in, err = s.cmd.StdinPipe(); err != nil {
		t.Fatal(err)
	}
	if err := s.cmd.Start(); err != nil {
		t.Fatal(err)
	}
	go func() { s.waitErr = s.cmd.Wait(); close(s.exited) }()
	t.Cleanup(func() {
		_ = s.cmd.Process.Kill()
		select {
		case <-s.exited:
		case <-time.After(10 * time.Second):
		}
	})
	return s
}

func (s *scriptSession) send(text string) { _, _ = io.WriteString(s.in, text) }

func (s *scriptSession) waitFor(text string, within time.Duration) bool {
	deadline := time.Now().Add(within)
	for time.Now().Before(deadline) {
		if strings.Contains(s.out.String(), text) {
			return true
		}
		time.Sleep(50 * time.Millisecond)
	}
	return strings.Contains(s.out.String(), text)
}

// sync types a command until its output shows that the shell reads input.
// The output carries tag+"2"; the typed command, echoed by a terminal,
// carries `$((1+1))` instead.
func (s *scriptSession) sync(t *testing.T, tag string) {
	t.Helper()
	start := time.Now()
	deadline := start.Add(90 * time.Second)
	for attempt := 1; time.Now().Before(deadline); attempt++ {
		s.send("echo " + tag + "$((1+1))\n")
		if s.waitFor(tag+"2", 3*time.Second) {
			t.Logf("the shell answered after %v (%d attempt(s))", time.Since(start).Round(time.Millisecond), attempt)
			return
		}
	}
	t.Fatalf("the shell never answered %s2; output so far:\n%q", tag, s.out.String())
}

// finish waits for the end of the helper and returns its result file.
func (s *scriptSession) finish(t *testing.T) (exit int, result string) {
	t.Helper()
	select {
	case <-s.exited:
		var exitErr *exec.ExitError
		if errors.As(s.waitErr, &exitErr) {
			exit = exitErr.ExitCode()
		} else if s.waitErr != nil {
			t.Fatalf("script: %v", s.waitErr)
		}
	case <-time.After(60 * time.Second):
		t.Fatalf("the helper did not end; output:\n%q", s.out.String())
	}
	raw, _ := os.ReadFile(s.result)
	return exit, string(raw)
}

// clientPIDs returns the PIDs of the `docker exec` clients of the session in
// the container (host-side processes of the test container).
func clientPIDs(id string) []int {
	var pids []int
	entries, _ := filepath.Glob("/proc/[0-9]*/cmdline")
	for _, e := range entries {
		raw, err := os.ReadFile(e)
		if err != nil {
			continue
		}
		args := strings.Split(string(raw), "\x00")
		if len(args) > 3 && filepath.Base(args[0]) == "docker" && args[1] == "exec" && args[2] == "-i" && contains(args, id) {
			if pid, err := strconv.Atoi(strings.Split(strings.TrimPrefix(e, "/proc/"), "/")[0]); err == nil {
				pids = append(pids, pid)
			}
		}
	}
	return pids
}

// LH-FA-DEV-010 ("interaktive Shell"), ADR-0031 Entscheidung 4: the terminal
// path. The helper process calls Shell with Terminal and runs under `script`,
// so its standard input is a pseudo-terminal; in the session the standard
// input and output are terminals, and Ctrl-C reaches the shell in the
// container. Gegenprobe: the same helper under the same `script` and the same
// commands, with Terminal off (no `-t`), finds no terminal in the session.
func TestShell_Docker_Terminal_UnderScript(t *testing.T) {
	const check = "if test -t 0 && test -t 1; then echo TTY$((20+1))=yes; else echo TTY$((20+1))=no; fi\n"
	t.Run("with a terminal", func(t *testing.T) {
		id := newVscodeContainer(t)
		// Each startup file sets a variable; the session shows which file the
		// shell read.
		mustDocker(t, "exec", "-u", "vscode", id, "sh", "-c",
			"echo 'export UBL_BASHRC_READ=yes' >> /home/vscode/.bashrc && echo 'export UBL_PROFILE_READ=yes' >> /home/vscode/.profile")
		s := startUnderScript(t, "terminal", id)
		s.sync(t, "SYNC")
		s.send(check)
		if !s.waitFor("TTY21=yes", 20*time.Second) {
			t.Fatalf("the session found no terminal on standard input and output:\n%q", s.out.String())
		}
		// The kind of shell: not a login shell, started as `bash`, interactive,
		// reads ~/.bashrc and not ~/.profile. The typed commands carry
		// `$((…))`, `$0` and `${…}`, so only the answers contain the values.
		s.send("shopt -q login_shell && echo LOGIN$((1+1))=yes || echo LOGIN$((1+1))=no\n")
		if !s.waitFor("LOGIN2=", 20*time.Second) {
			t.Fatalf("the shell gave no answer about the login shell:\n%q", s.out.String())
		}
		if !strings.Contains(s.out.String(), "LOGIN2=no") {
			t.Errorf("the shell is a login shell or answered otherwise:\n%q", s.out.String())
		}
		s.send(`echo "ARG0=[$0] INTERACTIVE=${-//[^i]/} BASHRC=${UBL_BASHRC_READ:-unset} PROFILE=${UBL_PROFILE_READ:-unset} KIND$((1+1))"` + "\n")
		if !s.waitFor("KIND2", 20*time.Second) {
			t.Fatalf("the shell gave no answer about its kind:\n%q", s.out.String())
		}
		if want := "ARG0=[bash] INTERACTIVE=i BASHRC=yes PROFILE=unset KIND2"; !strings.Contains(s.out.String(), want) {
			t.Errorf("the answer about the kind of shell lacks %q:\n%q", want, s.out.String())
		}
		s.send("sleep 300\n")
		time.Sleep(time.Second)
		s.send("\x03")
		s.send("echo CTRLC$((1+1))\n")
		if !s.waitFor("CTRLC2", 20*time.Second) {
			t.Fatalf("Ctrl-C did not reach the shell in the container (sleep still holds the terminal):\n%q", s.out.String())
		}
		s.send("exit\n")
		exit, result := s.finish(t)
		if exit != 0 || result != "ok" {
			t.Errorf("helper exit = %d, result = %q; want 0 and ok (the end of the shell is no error)", exit, result)
		}
	})
	t.Run("guard: the same run without a terminal finds none", func(t *testing.T) {
		id := newVscodeContainer(t)
		s := startUnderScript(t, "noterminal", id)
		s.sync(t, "SYNC")
		s.send(check)
		if !s.waitFor("TTY21=no", 20*time.Second) {
			t.Fatalf("the session without `-t` did not answer no:\n%q", s.out.String())
		}
		if strings.Contains(s.out.String(), "TTY21=yes") {
			t.Errorf("the session without `-t` found a terminal:\n%q", s.out.String())
		}
		s.send("exit\n")
		exit, result := s.finish(t)
		if exit != 0 || result != "ok" {
			t.Errorf("helper exit = %d, result = %q; want 0 and ok", exit, result)
		}
	})
}

// ADR-0031 Entscheidung 4, Prüfauftrag: does the shell in the container end
// with the end of the client? Observed with engine and client 29.9.2 on Linux
// amd64: it ends only for an idle shell without a terminal (its standard
// input reaches the end); a shell that runs a command, and every shell with a
// terminal, outlives the client. The tests pin that observation; a test that
// turns red means the engine or the client behaves differently, and what
// manual and specification say about it changes with it.
func TestShell_Docker_ClientEnd(t *testing.T) {
	// startNoTerminal opens a session without a terminal whose standard input
	// stays open, types the input and waits for the shell to answer; the
	// returned function cancels the context and returns the error of Shell.
	startNoTerminal := func(t *testing.T, id string, input string) (cancel func() error) {
		t.Helper()
		in, out, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = in.Close(); _ = out.Close() })
		ctx, stop := context.WithCancel(context.Background())
		t.Cleanup(stop)
		var stdout syncBuffer
		done := make(chan error, 1)
		go func() {
			done <- docker.NewLauncher().Shell(ctx, vscodeShell(id, func(o *driven.DevcontainerShellOptions) {
				o.Stdin, o.Stdout = in, &stdout
			}))
		}()
		if _, err := out.WriteString(input); err != nil {
			t.Fatal(err)
		}
		deadline := time.Now().Add(30 * time.Second)
		for !strings.Contains(stdout.String(), "started") {
			if time.Now().After(deadline) {
				t.Fatal("the session did not start")
			}
			time.Sleep(50 * time.Millisecond)
		}
		return func() error {
			stop()
			select {
			case err := <-done:
				return err
			case <-time.After(60 * time.Second):
				return errors.New("Shell did not return after the cancel")
			}
		}
	}

	t.Run("without a terminal, idle shell: ends", func(t *testing.T) {
		id := newVscodeContainer(t)
		cancel := startNoTerminal(t, id, "echo started\n")
		if procs := containerProcs(t, id, "bash"); len(procs) != 1 {
			t.Fatalf("bash processes before the cancel = %v, want the one of the session", procs)
		}
		if err := cancel(); !errors.Is(err, context.Canceled) {
			t.Errorf("Shell error = %v, want the context error", err)
		}
		if procs := waitProcs(t, id, "bash", 15*time.Second, func(n int) bool { return n == 0 }); len(procs) != 0 {
			t.Errorf("the idle shell outlived the client: %v", procs)
		}
	})
	t.Run("without a terminal, busy shell: outlives the client", func(t *testing.T) {
		id := newVscodeContainer(t)
		cancel := startNoTerminal(t, id, "echo started\nsleep 300\n")
		requireProcs(t, id, "bash", 1, "before the cancel")
		requireProcs(t, id, "sleep", 1, "before the cancel (the command must run, or the case is the idle one)")
		if err := cancel(); !errors.Is(err, context.Canceled) {
			t.Errorf("Shell error = %v, want the context error", err)
		}
		time.Sleep(3 * time.Second)
		if procs := containerProcs(t, id, "bash"); len(procs) != 1 {
			t.Errorf("bash processes after the client ended = %v; the observation is that the busy shell stays", procs)
		}
		if procs := containerProcs(t, id, "sleep"); len(procs) != 1 {
			t.Errorf("sleep processes after the client ended = %v; the observation is that the command stays", procs)
		}
	})
	for _, busy := range []bool{false, true} {
		name := "with a terminal, idle shell: outlives the client"
		if busy {
			name = "with a terminal, busy shell: outlives the client"
		}
		t.Run(name, func(t *testing.T) {
			id := newVscodeContainer(t)
			s := startUnderScript(t, "terminal", id)
			s.sync(t, "SYNC")
			wantSleep := 0
			if busy {
				s.send("sleep 300\n")
				wantSleep = 1
			}
			requireProcs(t, id, "bash", 1, "before the kill")
			requireProcs(t, id, "sleep", wantSleep, "before the kill")
			pids := clientPIDs(id)
			if len(pids) != 1 {
				t.Fatalf("docker exec clients of the session = %v, want one", pids)
			}
			if err := syscall.Kill(pids[0], syscall.SIGKILL); err != nil {
				t.Fatal(err)
			}
			s.finish(t)
			time.Sleep(3 * time.Second)
			if procs := containerProcs(t, id, "bash"); len(procs) != 1 {
				t.Errorf("bash processes after the client was killed = %v; the observation is that the shell stays", procs)
			}
			if procs := containerProcs(t, id, "sleep"); len(procs) != wantSleep {
				t.Errorf("sleep processes after the client was killed = %v, want %d; the observation is that the command stays with its shell", procs, wantSleep)
			}
		})
	}
}

// requireProcs fails the test unless exactly want processes called comm run in
// the container within 15 s; when names what the check stands for.
func requireProcs(t *testing.T, id, comm string, want int, when string) {
	t.Helper()
	procs := waitProcs(t, id, comm, 15*time.Second, func(n int) bool { return n == want })
	if len(procs) != want {
		t.Fatalf("%s processes %s = %v, want %d", comm, when, procs, want)
	}
}
