package docker_test

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/u-boot/internal/adapter/driven/docker"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

const (
	launcherTestFolder = "/tmp/ubl-unit/project"
	launcherTestImage  = "u-boot-devcontainer-cli:test"
	launcherSuccessOut = `{"outcome":"success","containerId":"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef","remoteUser":"vscode","remoteWorkspaceFolder":"/workspaces/project"}`
)

// fakeDocker is a shell script that stands in for the docker binary. It
// logs every argv (one argument per line, calls separated by `--`) and
// answers by the flag files in its directory.
type fakeDocker struct {
	dir string
	bin string
}

const fakeDockerScript = `#!/bin/sh
D='%DIR%'
{ for a in "$@"; do printf '%s\n' "$a"; done; echo '--'; } >> "$D/argv.log"
case "$1" in
version)
  [ -e "$D/daemon-down" ] && exit 1
  [ -e "$D/daemon-down-late" ] && grep -q -x -e '-i' "$D/argv.log" && exit 1
  echo 29.0.0; exit 0 ;;
context) [ -e "$D/context-fail" ] && exit 1; cat "$D/context-host"; exit 0 ;;
image) [ -e "$D/image-missing" ] && exit 1; exit 0 ;;
build) cat > "$D/context.tar"; [ -e "$D/build-fail" ] && { echo 'build-boom' >&2; exit 1; }; exit 0 ;;
run)
  [ -e "$D/run-hang" ] && exec sleep 30
  [ -e "$D/run-stderr" ] && cat "$D/run-stderr" >&2
  [ -e "$D/run-stderr-b" ] && { cat "$D/run-stderr-a" >&2; sleep 0.3; cat "$D/run-stderr-b" >&2; }
  [ -e "$D/run-fail" ] && { echo 'run-boom' >&2; exit 3; }
  cat "$D/run-stdout"; exit 0 ;;
rm)
  [ -e "$D/rm-fail" ] && { echo 'rm-boom' >&2; exit 1; }
  [ -e "$D/rm-nosuch" ] && { echo 'Error response from daemon: No such container: x' >&2; exit 1; }
  exit 0 ;;
stop) [ -e "$D/stop-fail" ] && { echo 'stop-boom' >&2; exit 1; }; exit 0 ;;
exec)
  case " $* " in
  *" -i "*)
    [ -e "$D/session-record-fds" ] && { l0=$(readlink /proc/$$/fd/0); l1=$(readlink /proc/$$/fd/1); echo "$l0" > "$D/fd0.link"; echo "$l1" > "$D/fd1.link"; }
    [ -e "$D/session-copy-stdin" ] && cat > "$D/stdin.copy"
    [ -e "$D/session-stdout" ] && cat "$D/session-stdout"
    [ -e "$D/session-stderr" ] && cat "$D/session-stderr" >&2
    [ -e "$D/session-hang" ] && exec sleep 30
    [ -e "$D/session-delay" ] && { sleep 1; echo session-done; }
    [ -e "$D/session-exit" ] && exit "$(cat "$D/session-exit")"
    exit 0 ;;
  esac
  n=$#; i=0
  for a in "$@"; do i=$((i+1)); [ "$i" -eq $((n-2)) ] && shell=$a; done
  [ -e "$D/user-missing" ] && { echo 'Error response from daemon: unable to find user nosuchuser: no matching entries in passwd file' >&2; exit 1; }
  [ -e "$D/workspace-missing" ] && { echo 'OCI runtime exec failed: exec failed: unable to start container process: chdir to cwd ("/nonexistent") set in config.json failed: no such file or directory'; exit 127; }
  [ -e "$D/$shell-missing" ] && { echo "OCI runtime exec failed: exec failed: unable to start container process: exec: \"$shell\": executable file not found in \$PATH"; exit 127; }
  exit 0 ;;
esac
exit 0
`

// newFakeDocker legt ein ausführbares docker-Skript an. Die Tests dieser Datei
// laufen nicht parallel: Ein parallel forkender Test hält die frisch
// geschriebene Datei offen, und der exec scheitert mit "text file busy".
func newFakeDocker(t *testing.T) *fakeDocker {
	t.Helper()
	dir := t.TempDir()
	f := &fakeDocker{dir: dir, bin: filepath.Join(dir, "docker")}
	script := strings.ReplaceAll(fakeDockerScript, "%DIR%", dir)
	if err := os.WriteFile(f.bin, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	f.write(t, "context-host", "unix:///var/run/docker.sock\n")
	f.write(t, "run-stdout", launcherSuccessOut+"\n")
	return f
}

func (f *fakeDocker) write(t *testing.T, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(f.dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func (f *fakeDocker) flag(t *testing.T, names ...string) {
	t.Helper()
	for _, name := range names {
		f.write(t, name, "")
	}
}

// calls returns the recorded argv lists.
func (f *fakeDocker) calls(t *testing.T) [][]string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(f.dir, "argv.log"))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	var calls [][]string
	var current []string
	for _, line := range strings.Split(strings.TrimSuffix(string(raw), "\n"), "\n") {
		if line == "--" {
			calls = append(calls, current)
			current = nil
			continue
		}
		current = append(current, line)
	}
	return calls
}

// call returns the first recorded argv whose first argument is verb.
func (f *fakeDocker) call(t *testing.T, verb string) []string {
	t.Helper()
	for _, c := range f.calls(t) {
		if len(c) > 0 && c[0] == verb {
			return c
		}
	}
	return nil
}

func (f *fakeDocker) launcher(env map[string]string) *docker.Launcher {
	return docker.NewLauncherForTest(f.bin, func(key string) (string, bool) {
		v, ok := env[key]
		return v, ok
	}, 2*time.Second)
}

func upOptions() driven.DevcontainerUpOptions {
	return driven.DevcontainerUpOptions{ProjectFolder: launcherTestFolder}
}

// LH-FA-DEV-006, ADR-0028 Entscheidung 4: the argument list of the helper
// container.
func TestRunArgs(t *testing.T) {
	sock := "type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock"
	folder := "type=bind,source=" + launcherTestFolder + ",target=" + launcherTestFolder
	cases := []struct {
		name string
		env  []string
		opts driven.DevcontainerUpOptions
		want []string
	}{
		{
			name: "defaults",
			opts: upOptions(),
			want: []string{"run", "--rm", "--name", "helper", "--mount", sock, "--mount", folder,
				launcherTestImage, "devcontainer", "up", "--workspace-folder", launcherTestFolder},
		},
		{
			name: "env by name, recreate, read-only",
			env:  []string{"GIT_TOKEN", "OTHER"},
			opts: driven.DevcontainerUpOptions{ProjectFolder: launcherTestFolder, Recreate: true, ReadOnlyProject: true},
			want: []string{"run", "--rm", "--name", "helper", "--mount", sock, "--mount", folder + ",readonly",
				"-e", "GIT_TOKEN", "-e", "OTHER",
				launcherTestImage, "devcontainer", "up", "--workspace-folder", launcherTestFolder, "--remove-existing-container"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := docker.RunArgsForTest("helper", launcherTestImage, tc.env, tc.opts)
			if strings.Join(got, "\x00") != strings.Join(tc.want, "\x00") {
				t.Errorf("args =\n%q\nwant\n%q", got, tc.want)
			}
			for _, arg := range got {
				if arg == "-v" || arg == "--volume" || arg == "--privileged" || arg == "--network" || strings.HasPrefix(arg, "--cap-add") {
					t.Errorf("argument %q is not part of the helper container contract", arg)
				}
			}
		})
	}
}

// LH-FA-DEV-009: the value of a passed variable is on no command line; a
// name that is not set on the host is not passed.
func TestLauncher_Up_PassesEnvByNameOnly(t *testing.T) {
	f := newFakeDocker(t)
	const secret = "placeholder-token-value"
	opts := upOptions()
	opts.EnvNames = []string{"GIT_TOKEN", "NOT_SET"}
	res, err := f.launcher(map[string]string{"GIT_TOKEN": secret, "UNLISTED": "placeholder-unlisted"}).Up(context.Background(), opts)
	if err != nil {
		t.Fatalf("Up: %v", err)
	}
	if res.RemoteUser != "vscode" || res.RemoteWorkspaceFolder != "/workspaces/project" || len(res.ContainerID) != 64 {
		t.Errorf("result = %+v", res)
	}
	run := f.call(t, "run")
	joined := strings.Join(f.flatCalls(t), "\n")
	if strings.Contains(joined, secret) || strings.Contains(joined, "placeholder-unlisted") {
		t.Errorf("a value reached a command line:\n%s", joined)
	}
	var passed []string
	for i, arg := range run {
		if arg == "-e" {
			passed = append(passed, run[i+1])
		}
	}
	if strings.Join(passed, ",") != "GIT_TOKEN" {
		t.Errorf("passed variables = %v, want [GIT_TOKEN]", passed)
	}
}

func (f *fakeDocker) flatCalls(t *testing.T) []string {
	t.Helper()
	var all []string
	for _, c := range f.calls(t) {
		all = append(all, strings.Join(c, " "))
	}
	return all
}

// ADR-0028 Entscheidung 1: the image is built from the tar context only
// when the tag is missing.
func TestLauncher_Up_BuildsHelperImageOnlyWhenTagIsMissing(t *testing.T) {
	present := newFakeDocker(t)
	if _, err := present.launcher(nil).Up(context.Background(), upOptions()); err != nil {
		t.Fatalf("Up: %v", err)
	}
	if present.call(t, "build") != nil {
		t.Errorf("docker build ran although the image tag exists")
	}

	missing := newFakeDocker(t)
	missing.flag(t, "image-missing")
	if _, err := missing.launcher(nil).Up(context.Background(), upOptions()); err != nil {
		t.Fatalf("Up: %v", err)
	}
	build := missing.call(t, "build")
	tag, err := docker.HelperImageTagForTest()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(build, " ") != "build --force-rm -t "+tag+" -" {
		t.Errorf("build argv = %q, want build --force-rm -t %s -", build, tag)
	}
	run := missing.call(t, "run")
	if !contains(run, tag) {
		t.Errorf("run argv %q does not use the image %s", run, tag)
	}
	raw, err := os.ReadFile(filepath.Join(missing.dir, "context.tar"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := docker.HelperContextTarForTest()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(raw, want) {
		t.Errorf("tar stream on stdin of docker build differs from the embedded context")
	}
}

func contains(list []string, s string) bool {
	for _, e := range list {
		if e == s {
			return true
		}
	}
	return false
}

// ADR-0028 Entscheidung 7: the environment check classifies 11 vs 12 and
// runs before any build or run.
func TestLauncher_Preflight(t *testing.T) {
	cases := []struct {
		name  string
		setup func(t *testing.T, f *fakeDocker)
		env   map[string]string
		want  error
	}{
		{"daemon down", func(t *testing.T, f *fakeDocker) { f.flag(t, "daemon-down") }, nil, driven.ErrDockerUnavailable},
		{"context cannot be read", func(t *testing.T, f *fakeDocker) { f.flag(t, "context-fail") }, nil, driven.ErrDockerUnavailable},
		{"remote DOCKER_HOST", func(*testing.T, *fakeDocker) {}, map[string]string{"DOCKER_HOST": "tcp://10.0.0.1:2375"}, driven.ErrContainerEngineNotLocal},
		{"remote context", func(t *testing.T, f *fakeDocker) { f.write(t, "context-host", "ssh://user@host\n") }, nil, driven.ErrContainerEngineNotLocal},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newFakeDocker(t)
			tc.setup(t, f)
			l := f.launcher(tc.env)
			_, upErr := l.Up(context.Background(), upOptions())
			stopErr := l.Stop(context.Background(), strings.Repeat("a", 12))
			shellErr := l.Shell(context.Background(), shellOptions(nil))
			for name, err := range map[string]error{"Up": upErr, "Stop": stopErr, "Shell": shellErr} {
				if !errors.Is(err, tc.want) {
					t.Errorf("%s error = %v, want %v", name, err, tc.want)
				}
				if errors.Is(err, driven.ErrDevcontainerRuntime) {
					t.Errorf("%s error %v is classified as runtime failure", name, err)
				}
			}
			for _, verb := range []string{"build", "run", "stop", "exec"} {
				if f.call(t, verb) != nil {
					t.Errorf("docker %s ran although the environment check failed", verb)
				}
			}
		})
	}
}

func TestLauncher_Preflight_MissingBinary(t *testing.T) {
	l := docker.NewLauncherForTest("/does/not/exist/docker-binary", func(string) (string, bool) { return "", false }, time.Second)
	if _, err := l.Up(context.Background(), upOptions()); !errors.Is(err, driven.ErrDockerUnavailable) {
		t.Errorf("Up error = %v, want ErrDockerUnavailable", err)
	}
	if err := l.Stop(context.Background(), strings.Repeat("a", 12)); !errors.Is(err, driven.ErrDockerUnavailable) {
		t.Errorf("Stop error = %v, want ErrDockerUnavailable", err)
	}
	if err := l.Shell(context.Background(), shellOptions(nil)); !errors.Is(err, driven.ErrDockerUnavailable) {
		t.Errorf("Shell error = %v, want ErrDockerUnavailable", err)
	}
}

// ADR-0028 Entscheidung 3: every failure after the environment check is
// ErrDevcontainerRuntime.
func TestLauncher_Up_RuntimeFailures(t *testing.T) {
	cases := []struct {
		name     string
		setup    func(t *testing.T, f *fakeDocker)
		wantText string
		wantRun  bool
		wantRm   bool
	}{
		{"build fails", func(t *testing.T, f *fakeDocker) { f.flag(t, "image-missing", "build-fail") }, "build-boom", false, false},
		{"run exits non-zero", func(t *testing.T, f *fakeDocker) { f.flag(t, "run-fail") }, "run-boom", true, true},
		{"outcome is not success", func(t *testing.T, f *fakeDocker) {
			f.write(t, "run-stdout", `{"outcome":"error","message":"hook-boom"}`+"\n")
		}, "hook-boom", true, false},
		{"no result line", func(t *testing.T, f *fakeDocker) { f.write(t, "run-stdout", "") }, "no result line", true, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newFakeDocker(t)
			tc.setup(t, f)
			_, err := f.launcher(nil).Up(context.Background(), upOptions())
			if !errors.Is(err, driven.ErrDevcontainerRuntime) {
				t.Fatalf("error = %v, want ErrDevcontainerRuntime", err)
			}
			if errors.Is(err, driven.ErrDockerUnavailable) {
				t.Errorf("error %v is also classified as ErrDockerUnavailable", err)
			}
			if !strings.Contains(err.Error(), tc.wantText) {
				t.Errorf("error %q does not carry %q", err, tc.wantText)
			}
			if (f.call(t, "run") != nil) != tc.wantRun {
				t.Errorf("docker run ran = %v, want %v", f.call(t, "run") != nil, tc.wantRun)
			}
			if (f.call(t, "rm") != nil) != tc.wantRm {
				t.Errorf("docker rm -f ran = %v, want %v", f.call(t, "rm") != nil, tc.wantRm)
			}
		})
	}
}

// ADR-0028 Entscheidung 4, Lebensdauer: a cancelled run ends with the
// context error and removes the helper container by name.
func TestLauncher_Up_CancelRemovesHelperContainer(t *testing.T) {
	f := newFakeDocker(t)
	err := cancelHangingUp(t, f)
	if !errors.Is(err, context.Canceled) || errors.Is(err, driven.ErrDevcontainerRuntime) {
		t.Errorf("error = %v, want the context error", err)
	}
	run, rm := f.call(t, "run"), f.call(t, "rm")
	name := run[3]
	if !strings.HasPrefix(name, docker.HelperContainerPrefixForTest) {
		t.Fatalf("helper container name %q lacks the prefix", name)
	}
	if strings.Join(rm, " ") != "rm -f "+name {
		t.Errorf("rm argv = %q, want rm -f %s", rm, name)
	}
}

// cancelHangingUp starts an Up whose `docker run` hangs, cancels it once
// `docker run` was seen and returns the error of Up.
func cancelHangingUp(t *testing.T, f *fakeDocker) error {
	t.Helper()
	f.flag(t, "run-hang")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() {
		_, err := f.launcher(nil).Up(ctx, upOptions())
		done <- err
	}()
	deadline := time.Now().Add(20 * time.Second)
	for f.call(t, "run") == nil {
		if time.Now().After(deadline) {
			t.Fatal("docker run did not start")
		}
		time.Sleep(20 * time.Millisecond)
	}
	cancel()
	select {
	case err := <-done:
		return err
	case <-time.After(30 * time.Second):
		t.Fatal("Up did not return after cancel")
		return nil
	}
}

// ADR-0028 Entscheidung 4, Lebensdauer: a failing `docker rm -f` is
// reported, except for a container that is already gone.
func TestLauncher_Up_ReportsFailedHelperRemoval(t *testing.T) {
	t.Run("run fails, rm fails", func(t *testing.T) {
		f := newFakeDocker(t)
		f.flag(t, "run-fail", "rm-fail")
		_, err := f.launcher(nil).Up(context.Background(), upOptions())
		if !errors.Is(err, driven.ErrDevcontainerRuntime) || !strings.Contains(err.Error(), "may still run") || !strings.Contains(err.Error(), "rm-boom") {
			t.Errorf("error = %v, want ErrDevcontainerRuntime that reports the failed rm", err)
		}
	})
	t.Run("cancelled, rm fails", func(t *testing.T) {
		f := newFakeDocker(t)
		f.flag(t, "rm-fail")
		err := cancelHangingUp(t, f)
		if !errors.Is(err, context.Canceled) || errors.Is(err, driven.ErrDevcontainerRuntime) ||
			!strings.Contains(err.Error(), "may still run") || !strings.Contains(err.Error(), "rm-boom") {
			t.Errorf("error = %v, want the context error that reports the failed rm", err)
		}
	})
	t.Run("run fails, container already gone", func(t *testing.T) {
		f := newFakeDocker(t)
		f.flag(t, "run-fail", "rm-nosuch")
		_, err := f.launcher(nil).Up(context.Background(), upOptions())
		if !errors.Is(err, driven.ErrDevcontainerRuntime) || strings.Contains(err.Error(), "may still run") {
			t.Errorf("error = %v, want ErrDevcontainerRuntime without a rm report", err)
		}
	})
}

// ADR-0028 Entscheidung 3: the result line is the last non-empty line, also
// when stdout is longer than the buffer the adapter keeps.
func TestLauncher_Up_ResultLineSurvivesLargeStdout(t *testing.T) {
	f := newFakeDocker(t)
	noise := strings.Repeat("hook output line that is not the result\n", (2<<20)/40)
	f.write(t, "run-stdout", noise+launcherSuccessOut+"\n")
	res, err := f.launcher(nil).Up(context.Background(), upOptions())
	if err != nil || len(res.ContainerID) != 64 {
		t.Fatalf("Up with %d bytes of stdout: result %+v, error %v", len(noise), res, err)
	}
}

// LH-FA-DEV-009: the value of a passed variable is replaced by a placeholder
// in the progress the adapter forwards and in its error texts, also when the
// value is split over two writes.
func TestLauncher_Up_MasksPassedValues(t *testing.T) {
	const secret = "placeholder-token-value"
	runLine := "docker run -e X=" + secret + " --rm image\n"
	cases := []struct {
		name            string
		setup           func(t *testing.T, f *fakeDocker)
		fails           bool
		placeholderIn   string // "progress" or "error"
		placeholderWant bool
	}{
		{"stderr, run succeeds", func(t *testing.T, f *fakeDocker) { f.write(t, "run-stderr", runLine) }, false, "progress", true},
		{"stderr, run fails", func(t *testing.T, f *fakeDocker) { f.write(t, "run-stderr", runLine); f.flag(t, "run-fail") }, true, "error", true},
		{"stderr without final newline", func(t *testing.T, f *fakeDocker) {
			f.write(t, "run-stderr", "docker run -e X="+secret)
			f.flag(t, "run-fail")
		}, true, "error", true},
		{"value split over two writes", func(t *testing.T, f *fakeDocker) {
			f.write(t, "run-stderr-a", "docker run -e X=placeholder-tok")
			f.write(t, "run-stderr-b", "en-value --rm image\n")
			f.flag(t, "run-fail")
		}, true, "error", true},
		{"outcome message", func(t *testing.T, f *fakeDocker) {
			f.write(t, "run-stdout", `{"outcome":"error","message":"hook printed `+secret+`"}`+"\n")
		}, true, "error", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newFakeDocker(t)
			tc.setup(t, f)
			opts := upOptions()
			opts.EnvNames = []string{"GIT_TOKEN"}
			var progress bytes.Buffer
			opts.Progress = &progress
			_, err := f.launcher(map[string]string{"GIT_TOKEN": secret}).Up(context.Background(), opts)
			if (err != nil) != tc.fails {
				t.Fatalf("Up error = %v, want failure %v", err, tc.fails)
			}
			texts := map[string]string{"progress": progress.String()}
			if err != nil {
				texts["error"] = err.Error()
			}
			for where, text := range texts {
				if strings.Contains(text, secret) {
					t.Errorf("the value is in the %s: %q", where, text)
				}
			}
			if tc.placeholderWant && !strings.Contains(texts[tc.placeholderIn], "[masked]") {
				t.Errorf("%s %q shows no placeholder: the probe cannot tell masking from absence", tc.placeholderIn, texts[tc.placeholderIn])
			}
		})
	}
}

// LH-FA-DEV-009, ADR-0028 Entscheidung 3: the result line is read from the
// unmasked stdout. A passed value that occurs in it (as the outcome, as a
// part of the container ID, as a field value) neither disturbs the detection
// of success nor changes a field of the result; the fields are data for the
// caller, not output.
func TestLauncher_Up_ResultLineMayCarryThePassedValue(t *testing.T) {
	for _, value := range []string{"success", "outcome", "0123", "abcdef0123456789", "vscode", "/workspaces/project", "containerId"} {
		t.Run(value, func(t *testing.T) {
			f := newFakeDocker(t)
			opts := upOptions()
			opts.EnvNames = []string{"GIT_TOKEN"}
			res, err := f.launcher(map[string]string{"GIT_TOKEN": value}).Up(context.Background(), opts)
			if err != nil {
				t.Fatalf("Up with the value %q in the result line: %v", value, err)
			}
			want := driven.DevcontainerUpResult{
				ContainerID:           strings.Repeat("0123456789abcdef", 4),
				RemoteUser:            "vscode",
				RemoteWorkspaceFolder: "/workspaces/project",
			}
			if res != want {
				t.Errorf("result = %+v, want %+v", res, want)
			}
		})
	}
}

// A passed value in the message of a failed outcome, or equal to the outcome,
// is masked in the error text; the failure itself stays ErrDevcontainerRuntime.
func TestLauncher_Up_MasksValueInFailedOutcome(t *testing.T) {
	f := newFakeDocker(t)
	f.write(t, "run-stdout", `{"outcome":"error","message":"hook printed error-token-value"}`+"\n")
	opts := upOptions()
	opts.EnvNames = []string{"GIT_TOKEN", "OTHER_TOKEN"}
	env := map[string]string{"GIT_TOKEN": "error-token-value", "OTHER_TOKEN": "error"}
	_, err := f.launcher(env).Up(context.Background(), opts)
	if !errors.Is(err, driven.ErrDevcontainerRuntime) {
		t.Fatalf("error = %v, want ErrDevcontainerRuntime", err)
	}
	for _, value := range env {
		if strings.Contains(err.Error(), value) {
			t.Errorf("the value %q is in the error text %q", value, err)
		}
	}
	if !strings.Contains(err.Error(), "[masked]") {
		t.Errorf("error %q shows no placeholder", err)
	}
}

// The held-back end of stderr is released when the call ends. On success the
// closing line of the CLI (no final newline), which a value of 23 bytes holds
// back, reaches the progress; on failure the fake binary's last line
// `run-boom` does, in the progress and in the error text.
func TestLauncher_Up_ReleasesHeldStderrTail(t *testing.T) {
	const secret = "placeholder-token-value"
	const closing = "closing line of the CLI"
	for _, fails := range []bool{false, true} {
		name, stderr, want := "run succeeds", "docker run -e X="+secret+" --rm\n"+closing, "docker run -e X=[masked] --rm\n"+closing
		if fails {
			name, stderr = "run fails", stderr+"\n"
			want += "\nrun-boom"
		}
		t.Run(name, func(t *testing.T) {
			f := newFakeDocker(t)
			f.write(t, "run-stderr", stderr)
			if fails {
				f.flag(t, "run-fail")
			}
			opts := upOptions()
			opts.EnvNames = []string{"GIT_TOKEN"}
			var progress bytes.Buffer
			opts.Progress = &progress
			_, err := f.launcher(map[string]string{"GIT_TOKEN": secret}).Up(context.Background(), opts)
			if (err != nil) != fails {
				t.Fatalf("Up error = %v, want failure %v", err, fails)
			}
			if got := progress.String(); got != want+"\n" && got != want {
				t.Errorf("progress = %q, want %q", got, want)
			}
			if fails && !strings.Contains(err.Error(), want) {
				t.Errorf("error %q lacks the end of stderr %q", err, want)
			}
		})
	}
}

// LH-FA-DEV-009: a value is masked only when it is passed, set and at least
// four bytes long; other text stays as it is.
func TestLauncher_Up_MasksOnlyPassedLongValues(t *testing.T) {
	const text = "short=abc unlisted=placeholder-unlisted unset=ok\n"
	f := newFakeDocker(t)
	f.write(t, "run-stderr", text)
	opts := upOptions()
	opts.EnvNames = []string{"SHORT", "NOT_SET"}
	var progress bytes.Buffer
	opts.Progress = &progress
	env := map[string]string{"SHORT": "abc", "UNLISTED": "placeholder-unlisted"}
	if _, err := f.launcher(env).Up(context.Background(), opts); err != nil {
		t.Fatal(err)
	}
	if got := progress.String(); got != text {
		t.Errorf("progress = %q, want the text unchanged", got)
	}
}

// The masking writer replaces every occurrence of every value, whatever the
// chunk boundaries are: every split of the stream into two writes, and a
// write per byte, give the result of a one-shot replacement.
func TestMaskingWriter_AllSplitsEqualOneShot(t *testing.T) {
	values := []string{"secret-one", "secret-one-longer", "ab", "line1\nline2", ""}
	stream := "x secret-one-longer y secret-one z secret-onesecret-one\nline1\nline2 ab secret-on"
	want := strings.NewReplacer("secret-one-longer", "[masked]", "secret-one", "[masked]", "line1\nline2", "[masked]").Replace(stream)
	for cut := 0; cut <= len(stream); cut++ {
		var out bytes.Buffer
		w := docker.NewMaskingWriterForTest(&out, values)
		for _, part := range []string{stream[:cut], stream[cut:]} {
			if n, err := w.Write([]byte(part)); err != nil || n != len(part) {
				t.Fatalf("Write = %d, %v", n, err)
			}
		}
		if err := w.Flush(); err != nil {
			t.Fatal(err)
		}
		if out.String() != want {
			t.Fatalf("cut %d: output\n%q\nwant\n%q", cut, out.String(), want)
		}
	}
	var out bytes.Buffer
	w := docker.NewMaskingWriterForTest(&out, values)
	for i := 0; i < len(stream); i++ {
		_, _ = w.Write([]byte{stream[i]})
	}
	_ = w.Flush()
	if out.String() != want {
		t.Errorf("byte-wise output\n%q\nwant\n%q", out.String(), want)
	}
}

func TestLauncher_Up_RejectsUnsafeInput(t *testing.T) {
	f := newFakeDocker(t)
	l := f.launcher(map[string]string{"A": "x"})
	for _, folder := range []string{"relative/dir", "", "/tmp/a,b", `/tmp/a"b`, "/tmp/a\nb", "/tmp/a\x00b"} {
		_, err := l.Up(context.Background(), driven.DevcontainerUpOptions{ProjectFolder: folder})
		if !errors.Is(err, driven.ErrDevcontainerRuntime) {
			t.Errorf("folder %q: error = %v, want ErrDevcontainerRuntime", folder, err)
		}
	}
	for _, name := range []string{"A=x", "1A", "", "A B", "A-B"} {
		opts := upOptions()
		opts.EnvNames = []string{name}
		_, err := l.Up(context.Background(), opts)
		if !errors.Is(err, driven.ErrDevcontainerRuntime) {
			t.Errorf("variable name %q: error = %v, want ErrDevcontainerRuntime", name, err)
		}
	}
	if len(f.calls(t)) != 0 {
		t.Errorf("docker ran for rejected input: %v", f.flatCalls(t))
	}
}

func TestLauncher_Stop(t *testing.T) {
	id := strings.Repeat("ab", 32)
	f := newFakeDocker(t)
	if err := f.launcher(nil).Stop(context.Background(), id); err != nil {
		t.Fatalf("Stop: %v", err)
	}
	if got := strings.Join(f.call(t, "stop"), " "); got != "stop "+id {
		t.Errorf("stop argv = %q", got)
	}
	for _, bad := range []string{"", "xyz", "--all", "abc def", strings.Repeat("a", 11)} {
		if err := f.launcher(nil).Stop(context.Background(), bad); !errors.Is(err, driven.ErrDevcontainerRuntime) {
			t.Errorf("Stop(%q) error = %v, want ErrDevcontainerRuntime", bad, err)
		}
	}
	f.flag(t, "stop-fail")
	err := f.launcher(nil).Stop(context.Background(), id)
	if !errors.Is(err, driven.ErrDevcontainerRuntime) || !strings.Contains(err.Error(), "stop-boom") {
		t.Errorf("failing stop: error = %v", err)
	}
}

// ADR-0028 Entscheidung 3: the result line is the last non-empty line.
func TestParseUpResult(t *testing.T) {
	good := launcherSuccessOut
	cases := []struct {
		name    string
		stdout  string
		wantErr string
	}{
		{"single line", good + "\n", ""},
		{"noise before, blank lines after", "log line\n" + good + "\n\n", ""},
		{"empty", "", "no result line"},
		{"not JSON", "Container started\n", "not a JSON result"},
		{"error outcome", `{"outcome":"error","message":"boom"}`, `outcome "error": boom`},
		{"missing container id", `{"outcome":"success","containerId":"zz"}`, "no container ID"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := docker.ParseUpResultForTest(tc.stdout)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Errorf("error = %v, want %q", err, tc.wantErr)
				}
				return
			}
			if err != nil || res.RemoteUser != "vscode" || res.RemoteWorkspaceFolder != "/workspaces/project" {
				t.Errorf("result = %+v, err = %v", res, err)
			}
		})
	}
}

var (
	digestFrom      = regexp.MustCompile(`^FROM \S+@sha256:[0-9a-f]{64}( AS \S+)?$`)
	exactVersion    = regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	integrityPrefix = regexp.MustCompile(`^sha512-[A-Za-z0-9+/]+=*$`)
)

// ADR-0028 Entscheidung 2: the pins are in the embedded files.
func TestHelperContext_Pins(t *testing.T) {
	files, err := docker.HelperContextFilesForTest()
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(docker.HelperContextNamesForTest(), ","); got != "Dockerfile,package.json,package-lock.json" {
		t.Errorf("context files = %s", got)
	}

	var froms int
	for _, line := range strings.Split(string(files["Dockerfile"]), "\n") {
		if strings.HasPrefix(line, "FROM ") {
			froms++
			if !digestFrom.MatchString(line) {
				t.Errorf("FROM line %q is not pinned by digest", line)
			}
		}
	}
	if froms != 2 {
		t.Errorf("Dockerfile has %d FROM lines, want 2 (Node and Docker CLI image)", froms)
	}
	dockerfile := string(files["Dockerfile"])
	if !strings.Contains(dockerfile, "RUN npm ci --ignore-scripts\n") || strings.Contains(dockerfile, "npm install") {
		t.Errorf("Dockerfile does not install with `npm ci --ignore-scripts`")
	}

	version, err := docker.HelperCLIVersionForTest()
	if err != nil || !exactVersion.MatchString(version) {
		t.Errorf("CLI pin = %q (err %v), want an exact version", version, err)
	}

	var lock struct {
		Packages map[string]struct {
			Version   string `json:"version"`
			Integrity string `json:"integrity"`
		} `json:"packages"`
	}
	if err := json.Unmarshal(files["package-lock.json"], &lock); err != nil {
		t.Fatal(err)
	}
	for path, pkg := range lock.Packages {
		if path == "" {
			continue
		}
		if !integrityPrefix.MatchString(pkg.Integrity) {
			t.Errorf("lock entry %s carries no integrity hash (%q)", path, pkg.Integrity)
		}
	}
	cli := lock.Packages["node_modules/@devcontainers/cli"]
	if cli.Version != version {
		t.Errorf("lock pins CLI %q, package.json pins %q", cli.Version, version)
	}
}

// ADR-0028 Entscheidung 1: the tag derives from the CLI version and the
// context hash; the tar stream holds exactly the three files.
func TestHelperContext_TagAndTar(t *testing.T) {
	version, err := docker.HelperCLIVersionForTest()
	if err != nil {
		t.Fatal(err)
	}
	tag, err := docker.HelperImageTagForTest()
	if err != nil {
		t.Fatal(err)
	}
	if !regexp.MustCompile(`^u-boot-devcontainer-cli:` + regexp.QuoteMeta(version) + `-[0-9a-f]{12}$`).MatchString(tag) {
		t.Errorf("tag = %q", tag)
	}
	raw, err := docker.HelperContextTarForTest()
	if err != nil {
		t.Fatal(err)
	}
	files, err := docker.HelperContextFilesForTest()
	if err != nil {
		t.Fatal(err)
	}
	tr := tar.NewReader(bytes.NewReader(raw))
	var names []string
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		body, err := io.ReadAll(tr)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(body, files[hdr.Name]) {
			t.Errorf("tar entry %s differs from the embedded file", hdr.Name)
		}
		names = append(names, hdr.Name)
	}
	if strings.Join(names, ",") != strings.Join(docker.HelperContextNamesForTest(), ",") {
		t.Errorf("tar entries = %v", names)
	}
}
