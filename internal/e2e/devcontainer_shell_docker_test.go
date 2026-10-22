//go:build docker

// End-to-end check of `u-boot devcontainer shell` (LH-FA-DEV-010, ADR-0031,
// ADR-0034): the built binary runs in a generated project against the engine of
// the runner, with a script on its standard input (no terminal) or under
// `script` (a pseudo-terminal). The tests observe the session from outside: the
// lines the shell writes, the containers of the engine and the exit code of the
// process. A missing docker binary or a missing `script` fails the test, it
// does not skip it.

package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
)

// dcHelperPrefix is the name prefix of the helper container of the start; the
// binary carries it without a test prefix (the constant `helperContainerPrefix`
// of the launcher). The tests see the helper during the start; a prefix that
// differs from the product's fails that guard.
const dcHelperPrefix = "u-boot-devcontainer-up-"

// lockedBuffer is a buffer that the writer goroutine of os/exec and the test
// read and write concurrently.
type lockedBuffer struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (l *lockedBuffer) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.b.Write(p)
}

func (l *lockedBuffer) String() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.b.String()
}

// dcQuote quotes s as one word for `sh -c`.
func dcQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'" }

// dcHelperNames lists the containers (running or not) whose name holds the
// prefix of the helper container.
func dcHelperNames() map[string]bool {
	out, err := dcDocker("ps", "-a", "--filter", "name="+dcHelperPrefix, "--format", "{{.Names}}")
	names := map[string]bool{}
	if err != nil {
		return names
	}
	for _, n := range strings.Fields(out) {
		names[n] = true
	}
	return names
}

// dcNewHelpers lists the helper containers that exist now and did not exist in
// baseline, sorted.
func dcNewHelpers(baseline map[string]bool) []string {
	var fresh []string
	for n := range dcHelperNames() {
		if !baseline[n] {
			fresh = append(fresh, n)
		}
	}
	sort.Strings(fresh)
	return fresh
}

// dcWatchHelpers polls the helper containers every 100 ms until the returned
// function is called, which returns the sorted names of the helper containers
// that appeared in the meantime.
func dcWatchHelpers(baseline map[string]bool) (stop func() []string) {
	seen := map[string]bool{}
	var mu sync.Mutex
	quit, done := make(chan struct{}), make(chan struct{})
	go func() {
		defer close(done)
		for {
			for _, n := range dcNewHelpers(baseline) {
				mu.Lock()
				seen[n] = true
				mu.Unlock()
			}
			select {
			case <-quit:
				return
			case <-time.After(100 * time.Millisecond):
			}
		}
	}()
	return func() []string {
		close(quit)
		<-done
		mu.Lock()
		defer mu.Unlock()
		var names []string
		for n := range seen {
			names = append(names, n)
		}
		sort.Strings(names)
		return names
	}
}

// dcShellProc is one run of `u-boot devcontainer shell` whose standard input
// the test writes to. out and errOut are two buffers, so the command writes its
// two streams to distinct writers.
type dcShellProc struct {
	cmd     *exec.Cmd
	in      io.WriteCloser
	out     *lockedBuffer
	errOut  *lockedBuffer
	done    chan struct{}
	waitErr error
}

// startDCShell starts the binary (or, for a pseudo-terminal, `script` around
// it) in dir with the environment of the test. With underScript both streams
// of the command are the pseudo-terminal, which `script` writes to out.
func startDCShell(ctx context.Context, t *testing.T, bin, dir string, underScript bool) *dcShellProc {
	t.Helper()
	var cmd *exec.Cmd
	if underScript {
		if _, err := exec.LookPath("script"); err != nil {
			t.Fatalf("`script` not on PATH but the terminal test needs it: %v", err)
		}
		cmd = exec.CommandContext(ctx, "script", "-q", "-e", "-c", dcQuote(bin)+" devcontainer shell", "/dev/null")
	} else {
		cmd = exec.CommandContext(ctx, bin, "devcontainer", "shell")
	}
	cmd.Dir = dir
	cmd.Env = os.Environ()
	p := &dcShellProc{cmd: cmd, out: &lockedBuffer{}, errOut: &lockedBuffer{}, done: make(chan struct{})}
	cmd.Stdout, cmd.Stderr = p.out, p.errOut
	if underScript {
		cmd.Stderr = p.out
	}
	in, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	p.in = in
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	go func() { p.waitErr = cmd.Wait(); close(p.done) }()
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		select {
		case <-p.done:
		case <-time.After(10 * time.Second):
		}
	})
	return p
}

func (p *dcShellProc) send(text string) { _, _ = io.WriteString(p.in, text) }

// waitForOut waits until out holds text; it reports whether it does.
func (p *dcShellProc) waitForOut(text string, within time.Duration) bool {
	deadline := time.Now().Add(within)
	for time.Now().Before(deadline) {
		if strings.Contains(p.out.String(), text) {
			return true
		}
		select {
		case <-p.done:
			return strings.Contains(p.out.String(), text)
		case <-time.After(100 * time.Millisecond):
		}
	}
	return strings.Contains(p.out.String(), text)
}

// finish waits for the end of the process and returns its exit code.
func (p *dcShellProc) finish(t *testing.T) int {
	t.Helper()
	select {
	case <-p.done:
	case <-time.After(2 * time.Minute):
		t.Fatalf("the command did not end; stdout:\n%q\nstderr:\n%q", p.out.String(), p.errOut.String())
	}
	var ee *exec.ExitError
	if errors.As(p.waitErr, &ee) {
		return ee.ExitCode()
	}
	if p.waitErr != nil {
		t.Fatalf("wait: %v", p.waitErr)
	}
	return 0
}

// runBinary runs the built binary in dir with the text on its standard input
// and returns stdout, stderr and the exit code. The two streams are distinct
// buffers.
func runBinary(ctx context.Context, t *testing.T, bin, dir, stdin string, args ...string) (string, string, int) {
	t.Helper()
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Dir = dir
	cmd.Env = os.Environ()
	cmd.Stdin = strings.NewReader(stdin)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	code := 0
	if err := cmd.Run(); err != nil {
		var ee *exec.ExitError
		if !errors.As(err, &ee) {
			t.Fatalf("run %v: %v", args, err)
		}
		code = ee.ExitCode()
	}
	return stdout.String(), stderr.String(), code
}

// runDCShell runs `devcontainer shell` with the script on its standard input.
func runDCShell(ctx context.Context, t *testing.T, bin, dir, script string) (string, string, int) {
	t.Helper()
	return runBinary(ctx, t, bin, dir, script, "devcontainer", "shell")
}

// dcAddRemoteEnv adds two entries of the shape the generator writes for an
// agent (`"<NAME>": "${localEnv:<NAME>}"`) to `remoteEnv` of the generated
// devcontainer.json.
func dcAddRemoteEnv(t *testing.T, dir string, names ...string) {
	t.Helper()
	path := filepath.Join(dir, ".devcontainer", "devcontainer.json")
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	const gitEntry = `"GIT_TOKEN": "${localEnv:GIT_TOKEN}"`
	if !bytes.Contains(body, []byte(gitEntry)) {
		t.Fatalf("devcontainer.json has no remoteEnv entry for GIT_TOKEN:\n%s", body)
	}
	entries := gitEntry
	for _, n := range names {
		entries += ",\n    \"" + n + "\": \"${localEnv:" + n + "}\""
	}
	if err := os.WriteFile(path, bytes.Replace(body, []byte(gitEntry), []byte(entries), 1), 0o644); err != nil {
		t.Fatal(err)
	}
}

// dcUpData is the `data` of the envelope of `devcontainer up --json` as far as
// the result line of the shell tests reads it.
type dcUpData struct {
	Data struct {
		ContainerID           string `json:"containerId"`
		RemoteUser            string `json:"remoteUser"`
		RemoteWorkspaceFolder string `json:"remoteWorkspaceFolder"`
		AlreadyRunning        bool   `json:"alreadyRunning"`
	} `json:"data"`
}

func dcUpResultLine(ctx context.Context, t *testing.T, bin, dir string) dcUpData {
	t.Helper()
	stdout, stderr, code := runBinary(ctx, t, bin, dir, "", "devcontainer", "up", "--json")
	var env dcUpData
	if err := json.Unmarshal([]byte(stdout), &env); err != nil || code != 0 {
		t.Fatalf("devcontainer up --json: exit %d, %v\nstdout:\n%s\nstderr:\n%s", code, err, stdout, stderr)
	}
	return env
}

func dcRunningContainer(t *testing.T, dir string) string {
	t.Helper()
	ids := strings.Fields(dcMustDocker(t, "ps", "-q", "--no-trunc", "--filter", "label="+dcFolderLbl+"="+dir))
	if len(ids) != 1 {
		t.Fatalf("running containers of the project = %v, want exactly one", ids)
	}
	return ids[0]
}

func envHolds(doc dcInspectDoc, value string) bool {
	for _, e := range doc.Config.Env {
		if strings.Contains(e, value) {
			return true
		}
	}
	return false
}

// LH-FA-DEV-010, LH-FA-DEV-009, LH-FA-DEV-006, LH-FA-CLI-005, LH-FA-CLI-006,
// ADR-0031 decisions 1 to 3: `u-boot devcontainer shell` in a generated sandbox
// project, started from nothing, with a script on the standard input.
//
//   - The session runs as the user and in the workspace of the result line of
//     the start (read from `up --json` afterwards); `GIT_TOKEN` and the variable
//     of an agent arrive with their values, a variable that is unset on the
//     host is absent. Guard: the same `docker exec` without `-e` does not carry
//     the token.
//   - The standard output holds the lines of the session and nothing of u-boot;
//     the disclosure of the socket (for the duration of the start) and the
//     progress are on the standard error, the session's lines are not.
//   - During the session no helper container exists. Guards: the helper was
//     seen during the start (so the prefix matches the product), and the search
//     finds a container started by hand with that prefix.
//   - The value is not in `Config.Env` of the devcontainer. Guard: the check
//     finds a value passed with `docker run -e`.
//   - A script that ends with `exit 3` ends the command with 0.
//
// A second subtest runs the same command under `script`: the standard input is
// a pseudo-terminal, the session finds a terminal on standard input and output
// and the command ends with 0. Guard: the same session through a pipe finds
// none. This is the pin of the terminal path at the command, through the
// production detection.
func TestDevcontainerShell_Docker_Session(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	suffix := dcSuffix()
	agentName, unsetName := "UBL_AGENT_KEY_"+suffix, "UBL_UNSET_"+suffix
	agentValue := "placeholder-agent-" + suffix
	t.Setenv("GIT_TOKEN", sandboxCanaryToken)
	t.Setenv(agentName, agentValue)
	_ = os.Unsetenv(unsetName)
	dir := dcGenerateProject(ctx, t, true, map[string]string{
		"devcontainer.sandbox.egress.enabled": "true",
		"devcontainer.sandbox.onUnavailable":  "warn",
	})
	dcAddRemoteEnv(t, dir, agentName, unsetName)
	bin := ubootBinary(t)

	var id, wantUser, wantWorkspace string
	t.Run("session without a terminal", func(t *testing.T) {
		baseline := dcHelperNames()
		stopWatch := dcWatchHelpers(baseline)
		p := startDCShell(ctx, t, bin, dir, false)
		p.send("echo begun-" + suffix + "\n")
		if !p.waitForOut("begun-"+suffix, 18*time.Minute) {
			stopWatch()
			t.Fatalf("the session did not begin\nstdout:\n%s\nstderr:\n%s", p.out.String(), p.errOut.String())
		}
		duringStart := stopWatch()
		if len(duringStart) == 0 {
			t.Errorf("guard: no helper container was seen during the start; the prefix %q may no longer be the product's", dcHelperPrefix)
		}
		if during := dcNewHelpers(baseline); len(during) != 0 {
			t.Errorf("helper containers exist during the session: %v", during)
		}
		p.send("id -un\npwd\nprintenv GIT_TOKEN\nprintenv " + agentName + "\nprintenv " + unsetName + " || echo unset-absent\n" +
			"echo end-" + suffix + "\nexit 3\n")
		if code := p.finish(t); code != 0 {
			t.Errorf("exit code = %d, want 0 whatever the exit status of the shell was (the script ended with exit 3)", code)
		}
		stdout, stderr := p.out.String(), p.errOut.String()
		t.Logf("stderr of the command (first 1500 bytes):\n%.1500s", stderr)

		lines := strings.Split(strings.TrimSpace(stdout), "\n")
		if len(lines) != 7 || lines[0] != "begun-"+suffix || lines[3] != sandboxCanaryToken || lines[4] != agentValue ||
			lines[5] != "unset-absent" || lines[6] != "end-"+suffix {
			t.Fatalf("standard output = %q, want the lines of the session only: begin, user, workspace, token, agent value, unset-absent, end", stdout)
		}
		sessionUser, sessionWorkspace := lines[1], lines[2]
		if !strings.Contains(stderr, "Warning (LH-FA-DEV-010): this command mounts the host container-engine socket") ||
			!strings.Contains(stderr, "for the duration of the start)") || strings.Contains(stderr, "duration of the call") {
			t.Errorf("standard error lacks the disclosure of the socket for the duration of the start:\n%.1500s", stderr)
		}
		if strings.Contains(stderr, sandboxCanaryToken) || strings.Contains(stderr, agentValue) || strings.Contains(stderr, "begun-"+suffix) {
			t.Errorf("standard error carries a value or a line of the session:\n%.1500s", stderr)
		}

		// The user and the workspace are those of the result line of the start.
		result := dcUpResultLine(ctx, t, bin, dir)
		if result.Data.RemoteUser != sessionUser || result.Data.RemoteWorkspaceFolder != sessionWorkspace || !result.Data.AlreadyRunning {
			t.Errorf("session user %q workspace %q, result line %+v", sessionUser, sessionWorkspace, result.Data)
		}
		id, wantUser, wantWorkspace = result.Data.ContainerID, result.Data.RemoteUser, result.Data.RemoteWorkspaceFolder
		if id != dcRunningContainer(t, dir) {
			t.Errorf("the result line names container %s, the project runs %s", id, dcRunningContainer(t, dir))
		}
		if wantUser != "vscode" || wantWorkspace != "/workspaces/dcstart" {
			t.Errorf("result line user %q workspace %q, want vscode and /workspaces/dcstart", wantUser, wantWorkspace)
		}
	})

	t.Run("guards of the session", func(t *testing.T) {
		if id == "" {
			t.Fatal("the session subtest left no container")
		}
		// The same exec without `-e` does not carry the token.
		if out, err := dcDocker("exec", "-u", wantUser, "-w", wantWorkspace, id, "printenv", "GIT_TOKEN"); err == nil || out != "" {
			t.Errorf("guard: the token is in the container without `-e`: %q (%v); the session test proves nothing", out, err)
		}
		// The value is not in the metadata of the devcontainer.
		doc := dcInspect(t, id)
		for _, v := range []string{sandboxCanaryToken, agentValue} {
			if envHolds(doc, v) {
				t.Errorf("Config.Env of the devcontainer holds %q", v)
			}
		}
		guard := "ubl-shellguard-" + suffix
		t.Cleanup(func() { _, _ = dcDocker("rm", "-f", guard) })
		dcMustDocker(t, "run", "-d", "--name", guard, "-e", "UBL_GUARD="+sandboxCanaryToken, dcBaseImage, "sleep", "60")
		if !envHolds(dcInspect(t, guard), sandboxCanaryToken) {
			t.Errorf("guard: the check of Config.Env does not find a value passed with `docker run -e`")
		}
		// The search for a helper container finds one started by hand.
		hand := dcHelperPrefix + "guard-" + suffix
		t.Cleanup(func() { _, _ = dcDocker("rm", "-f", hand) })
		before := dcHelperNames()
		dcMustDocker(t, "run", "-d", "--name", hand, dcBaseImage, "sleep", "60")
		if found := dcNewHelpers(before); len(found) != 1 || found[0] != hand {
			t.Errorf("guard: the search for helper containers found %v, want the one started by hand", found)
		}
	})

	t.Run("terminal under script", func(t *testing.T) {
		const check = "if test -t 0 && test -t 1; then echo TTY$((20+1))=yes; else echo TTY$((20+1))=no; fi\n"
		// The container runs, so the start inside the session is short.
		p := startDCShell(ctx, t, bin, dir, true)
		synced := false
		for deadline := time.Now().Add(3 * time.Minute); time.Now().Before(deadline) && !synced; {
			p.send("echo SYNC$((1+1))\n")
			synced = p.waitForOut("SYNC2", 3*time.Second)
		}
		if !synced {
			t.Fatalf("the shell never answered under script; output:\n%q", p.out.String())
		}
		p.send(check)
		if !p.waitForOut("TTY21=yes", 20*time.Second) {
			t.Fatalf("the session found no terminal on standard input and output under `script`:\n%q", p.out.String())
		}
		p.send("exit\n")
		if code := p.finish(t); code != 0 {
			t.Errorf("exit code = %d, want 0; output:\n%q", code, p.out.String())
		}

		// Guard: the same command through a pipe finds no terminal.
		stdout, stderr, code := runDCShell(ctx, t, bin, dir, check)
		if code != 0 || !strings.Contains(stdout, "TTY21=no") || strings.Contains(stdout, "TTY21=yes") {
			t.Errorf("guard: through a pipe exit %d, stdout %q, stderr %.500q; want exit 0 and `TTY21=no`", code, stdout, stderr)
		}
	})
}

// LH-FA-DEV-010, ADR-0030 decision 1, ADR-0031 decision 1: a container whose
// configured start script reported `failed` opens no session. In the project
// whose egress start script aborts (capability withdrawn, `onUnavailable:
// fail`, u-boot.yaml present) the first `shell` ends with 11 and runs none of
// its script, so does the second one against the container that runs, and the
// message of the second says that this call did not run the start script and
// names `--recreate`. Guard in the same container: without u-boot.yaml the
// next `shell` opens the session and runs the script (exit 0), so the 11 rests
// on the configuration that finds the status file and the script is one that
// would run.
func TestDevcontainerShell_Docker_FailedStatusOpensNoSession(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	suffix := dcSuffix()
	marker := "must-not-run-" + suffix
	script := "echo " + marker + "\n"
	dir := dcEgressFailProject(ctx, t, false)
	bin := ubootBinary(t)

	for _, step := range []string{"first call", "second call"} {
		stdout, stderr, code := runDCShell(ctx, t, bin, dir, script)
		t.Logf("%s: exit %d, stderr (first 900 bytes):\n%.900s", step, code, stderr)
		if code != 11 || strings.Contains(stdout, marker) || stdout != "" {
			t.Errorf("%s: exit %d, stdout %q, want exit 11 and no session", step, code, stdout)
		}
		if !strings.Contains(stderr, "devcontainer shell:") || strings.Contains(stderr, "devcontainer up:") {
			t.Errorf("%s: the error does not name the command:\n%.900s", step, stderr)
		}
	}
	_, stderr, _ := runDCShell(ctx, t, bin, dir, script)
	for _, want := range []string{
		"the container is running and this call did not run the start script",
		"u-boot devcontainer up --recreate",
	} {
		if !strings.Contains(stderr, want) {
			t.Errorf("the message of a repeated call lacks %q:\n%.900s", want, stderr)
		}
	}

	if err := os.Remove(filepath.Join(dir, "u-boot.yaml")); err != nil {
		t.Fatal(err)
	}
	stdout, stderr, code := runDCShell(ctx, t, bin, dir, script)
	if code != 0 || strings.TrimSpace(stdout) != marker {
		t.Errorf("guard: without u-boot.yaml exit %d, stdout %q, stderr %.500q; want the session to open and run the script", code, stdout, stderr)
	}
}
