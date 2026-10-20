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
version) [ -e "$D/daemon-down" ] && exit 1; echo 29.0.0; exit 0 ;;
context) [ -e "$D/context-fail" ] && exit 1; cat "$D/context-host"; exit 0 ;;
image) [ -e "$D/image-missing" ] && exit 1; exit 0 ;;
build) cat > "$D/context.tar"; [ -e "$D/build-fail" ] && { echo 'build-boom' >&2; exit 1; }; exit 0 ;;
run)
  [ -e "$D/run-hang" ] && exec sleep 30
  [ -e "$D/run-fail" ] && { echo 'run-boom' >&2; exit 3; }
  cat "$D/run-stdout"; exit 0 ;;
stop) [ -e "$D/stop-fail" ] && { echo 'stop-boom' >&2; exit 1; }; exit 0 ;;
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
			for name, err := range map[string]error{"Up": upErr, "Stop": stopErr} {
				if !errors.Is(err, tc.want) {
					t.Errorf("%s error = %v, want %v", name, err, tc.want)
				}
				if errors.Is(err, driven.ErrDevcontainerRuntime) {
					t.Errorf("%s error %v is classified as runtime failure", name, err)
				}
			}
			for _, verb := range []string{"build", "run", "stop"} {
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
		if !errors.Is(err, context.Canceled) || errors.Is(err, driven.ErrDevcontainerRuntime) {
			t.Errorf("error = %v, want the context error", err)
		}
	case <-time.After(30 * time.Second):
		t.Fatal("Up did not return after cancel")
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
