package docker_test

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driven/docker"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

const testContainerID = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

// fakeProcess is a started command with canned stdout, stderr and exit
// status. aborted records whether the adapter ended it early.
type fakeProcess struct {
	stdout  io.Reader
	stderr  string
	waitErr error
	aborted bool
}

func (p *fakeProcess) Read(b []byte) (int, error) { return p.stdout.Read(b) }
func (p *fakeProcess) Wait() (string, error)      { return p.stderr, p.waitErr }
func (p *fakeProcess) Abort()                     { p.aborted = true }

// fakeRunner answers each docker command by its first argument and records
// every argv, so a test can assert which engine commands ran.
type fakeRunner struct {
	mu        sync.Mutex
	calls     [][]string
	byCommand map[string]func(args []string) (docker.Process, error)
}

func (r *fakeRunner) starter() func(context.Context, ...string) (docker.Process, error) {
	return func(_ context.Context, args ...string) (docker.Process, error) {
		r.mu.Lock()
		r.calls = append(r.calls, append([]string(nil), args...))
		r.mu.Unlock()
		handler, ok := r.byCommand[args[0]]
		if !ok {
			return nil, errors.New("fake runner: no handler for " + args[0])
		}
		return handler(args)
	}
}

func textProcess(out string) func([]string) (docker.Process, error) {
	return func([]string) (docker.Process, error) { return &fakeProcess{stdout: strings.NewReader(out)}, nil }
}

func failingProcess(stderr string) func([]string) (docker.Process, error) {
	return func([]string) (docker.Process, error) {
		return &fakeProcess{stdout: strings.NewReader(""), stderr: stderr, waitErr: errors.New("exit status 1")}, nil
	}
}

func newTestInspector(r *fakeRunner, dockerHost string) *docker.Inspector {
	return docker.NewInspectorWithRunnerForTest(r.starter(), func(key string) string {
		if key == "DOCKER_HOST" {
			return dockerHost
		}
		return ""
	})
}

// LH-FA-DEV-007: the container is found by the devcontainer CLI label.
func TestInspector_FindContainer(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name        string
		ps          func([]string) (docker.Process, error)
		wantFound   bool
		wantRunning bool
		wantErr     bool
	}{
		{"running", textProcess(testContainerID + " running\n"), true, true, false},
		{"stopped", textProcess(testContainerID + " exited\n"), true, false, false},
		{"none", textProcess(""), false, false, false},
		{"blank output", textProcess("\n"), false, false, false},
		{"several, first wins", textProcess(testContainerID + " running\n" + strings.Repeat("a", 64) + " exited\n"), true, true, false},
		{"id not hex", textProcess("zzzzzzzzzzzz running\n"), false, false, true},
		{"id too short", textProcess("abc running\n"), false, false, true},
		{"missing state", textProcess(testContainerID + "\n"), false, false, true},
		{"engine failure", failingProcess("Cannot connect to the Docker daemon"), false, false, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := &fakeRunner{byCommand: map[string]func([]string) (docker.Process, error){
				"context": textProcess("unix:///var/run/docker.sock\n"),
				"ps":      tc.ps,
			}}
			got, found, err := newTestInspector(r, "").FindContainer(context.Background(), "/home/dev/proj")
			if (err != nil) != tc.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tc.wantErr)
			}
			if found != tc.wantFound || got.Running != tc.wantRunning {
				t.Errorf("found=%v running=%v, want found=%v running=%v", found, got.Running, tc.wantFound, tc.wantRunning)
			}
			if tc.wantFound && got.ID != testContainerID {
				t.Errorf("ID = %q", got.ID)
			}
		})
	}
}

// The filter value is one argv entry (no shell), and the label is the
// devcontainer CLI's.
func TestInspector_FindContainer_Argv(t *testing.T) {
	t.Parallel()
	r := &fakeRunner{byCommand: map[string]func([]string) (docker.Process, error){
		"context": textProcess("unix:///var/run/docker.sock\n"),
		"ps":      textProcess(""),
	}}
	folder := "/home/dev/my proj; rm -rf /"
	if _, _, err := newTestInspector(r, "").FindContainer(context.Background(), folder); err != nil {
		t.Fatal(err)
	}
	var ps []string
	for _, c := range r.calls {
		if c[0] == "ps" {
			ps = c
		}
	}
	if ps == nil {
		t.Fatalf("no docker ps call in %v", r.calls)
	}
	want := "label=devcontainer.local_folder=" + folder
	found := false
	for _, a := range ps {
		if a == want {
			found = true
		}
	}
	if !found {
		t.Errorf("ps argv %q lacks the single filter argument %q", ps, want)
	}
}

func TestInspector_FindContainer_RejectsBadFolder(t *testing.T) {
	t.Parallel()
	for _, folder := range []string{"", "relative/path", "/a\nb", "/a\x00b"} {
		r := &fakeRunner{byCommand: map[string]func([]string) (docker.Process, error){}}
		if _, _, err := newTestInspector(r, "").FindContainer(context.Background(), folder); err == nil {
			t.Errorf("folder %q: expected an error", folder)
		}
		if len(r.calls) != 0 {
			t.Errorf("folder %q: engine called: %v", folder, r.calls)
		}
	}
}

// A remote engine cannot be tied to this project folder: neither DOCKER_HOST
// nor the Docker context may name anything but a local unix socket.
func TestInspector_FindContainer_RemoteEngine(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name       string
		dockerHost string
		ctxHost    string
		wantRemote bool
	}{
		{"env tcp", "tcp://10.0.0.5:2375", "", true},
		{"env ssh", "ssh://dev@build-host", "", true},
		{"env unix", "unix:///run/user/1000/docker.sock", "", false},
		{"context ssh", "", "ssh://dev@build-host\n", true},
		{"context tcp", "", "tcp://10.0.0.5:2375\n", true},
		{"context unix", "", "unix:///var/run/docker.sock\n", false},
		{"env wins over a remote context", "unix:///var/run/docker.sock", "ssh://dev@build-host\n", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := &fakeRunner{byCommand: map[string]func([]string) (docker.Process, error){
				"context": textProcess(tc.ctxHost),
				"ps":      textProcess(testContainerID + " running\n"),
			}}
			_, found, err := newTestInspector(r, tc.dockerHost).FindContainer(context.Background(), "/home/dev/proj")
			if tc.wantRemote {
				if !errors.Is(err, driven.ErrContainerEngineNotLocal) || found {
					t.Errorf("err = %v found = %v, want ErrContainerEngineNotLocal", err, found)
				}
				for _, c := range r.calls {
					if c[0] == "ps" {
						t.Errorf("docker ps ran against a remote engine: %v", c)
					}
				}
				return
			}
			if err != nil || !found {
				t.Errorf("err = %v found = %v, want a local hit", err, found)
			}
		})
	}
}

func TestInspector_FindContainer_ContextInspectFails(t *testing.T) {
	t.Parallel()
	r := &fakeRunner{byCommand: map[string]func([]string) (docker.Process, error){
		"context": failingProcess("context store broken"),
		"ps":      textProcess(testContainerID + " running\n"),
	}}
	_, found, err := newTestInspector(r, "").FindContainer(context.Background(), "/home/dev/proj")
	if err == nil || found || errors.Is(err, driven.ErrContainerEngineNotLocal) {
		t.Errorf("err = %v found = %v, want a plain error and no hit", err, found)
	}
}

func TestInspector_FindContainer_StartFails(t *testing.T) {
	t.Parallel()
	r := &fakeRunner{byCommand: map[string]func([]string) (docker.Process, error){
		"context": func([]string) (docker.Process, error) { return nil, errors.New("docker: executable file not found") },
	}}
	if _, found, err := newTestInspector(r, "").FindContainer(context.Background(), "/home/dev/proj"); err == nil || found {
		t.Errorf("err = %v found = %v", err, found)
	}
}

// tarOf builds a one-entry tar archive the way `docker cp` streams it.
func tarOf(t *testing.T, hdr tar.Header, body string) string {
	t.Helper()
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	if hdr.Typeflag == tar.TypeReg {
		hdr.Size = int64(len(body))
	}
	hdr.Mode = 0o644
	if err := tw.WriteHeader(&hdr); err != nil {
		t.Fatal(err)
	}
	if hdr.Typeflag == tar.TypeReg {
		if _, err := tw.Write([]byte(body)); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

func cpRunner(p *fakeProcess) *fakeRunner {
	return &fakeRunner{byCommand: map[string]func([]string) (docker.Process, error){
		"cp": func([]string) (docker.Process, error) { return p, nil },
	}}
}

const statusPath = "/var/lib/u-boot/status/egress"

// LH-FA-DEV-007 / trust boundary: only a regular file within the limit is
// read; every other entry type and an oversized file are refused before the
// body is read.
func TestInspector_ReadFile_EntryTypes(t *testing.T) {
	t.Parallel()
	const maxBytes = 512
	cases := []struct {
		name     string
		archive  func(t *testing.T) string
		want     string
		wantErr  error
		wantBody bool
	}{
		{"regular file", func(t *testing.T) string {
			return tarOf(t, tar.Header{Name: "egress", Typeflag: tar.TypeReg}, "u-boot-status 1\nresult=ok\n")
		}, "u-boot-status 1\nresult=ok\n", nil, true},
		{"empty regular file", func(t *testing.T) string {
			return tarOf(t, tar.Header{Name: "egress", Typeflag: tar.TypeReg}, "")
		}, "", nil, true},
		{"file at the limit", func(t *testing.T) string {
			return tarOf(t, tar.Header{Name: "egress", Typeflag: tar.TypeReg}, strings.Repeat("x", maxBytes))
		}, strings.Repeat("x", maxBytes), nil, true},
		{"file over the limit", func(t *testing.T) string {
			return tarOf(t, tar.Header{Name: "egress", Typeflag: tar.TypeReg}, strings.Repeat("x", maxBytes+1))
		}, "", driven.ErrContainerFileInvalid, false},
		{"symlink", func(t *testing.T) string {
			return tarOf(t, tar.Header{Name: "egress", Typeflag: tar.TypeSymlink, Linkname: "/etc/shadow"}, "")
		}, "", driven.ErrContainerFileInvalid, false},
		{"fifo", func(t *testing.T) string {
			return tarOf(t, tar.Header{Name: "egress", Typeflag: tar.TypeFifo}, "")
		}, "", driven.ErrContainerFileInvalid, false},
		{"directory", func(t *testing.T) string {
			return tarOf(t, tar.Header{Name: "egress/", Typeflag: tar.TypeDir}, "")
		}, "", driven.ErrContainerFileInvalid, false},
		{"hard link", func(t *testing.T) string {
			return tarOf(t, tar.Header{Name: "egress", Typeflag: tar.TypeLink, Linkname: "other"}, "")
		}, "", driven.ErrContainerFileInvalid, false},
		{"character device", func(t *testing.T) string {
			return tarOf(t, tar.Header{Name: "egress", Typeflag: tar.TypeChar}, "")
		}, "", driven.ErrContainerFileInvalid, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			p := &fakeProcess{stdout: strings.NewReader(tc.archive(t))}
			got, err := newTestInspector(cpRunner(p), "").ReadFile(context.Background(), testContainerID, statusPath, maxBytes)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) || got != nil {
					t.Errorf("got %q, %v; want %v and no content", got, err, tc.wantErr)
				}
			} else if err != nil || string(got) != tc.want {
				t.Errorf("got %q, %v; want %q", got, err, tc.want)
			}
			if !p.aborted {
				t.Error("the stream stays open after the read; the adapter must end it")
			}
		})
	}
}

func TestInspector_ReadFile_Missing(t *testing.T) {
	t.Parallel()
	p := &fakeProcess{
		stdout:  strings.NewReader(""),
		stderr:  "Error response from daemon: Could not find the file " + statusPath + " in container x\n",
		waitErr: errors.New("exit status 1"),
	}
	_, err := newTestInspector(cpRunner(p), "").ReadFile(context.Background(), testContainerID, statusPath, 512)
	if !errors.Is(err, driven.ErrContainerFileNotFound) {
		t.Errorf("err = %v, want ErrContainerFileNotFound", err)
	}
}

func TestInspector_ReadFile_EngineErrors(t *testing.T) {
	t.Parallel()
	cases := map[string]*fakeProcess{
		"no such container": {stdout: strings.NewReader(""), stderr: "Error response from daemon: No such container: x", waitErr: errors.New("exit status 1")},
		"daemon down":       {stdout: strings.NewReader(""), stderr: "Cannot connect to the Docker daemon", waitErr: errors.New("exit status 1")},
		"truncated header":  {stdout: strings.NewReader(strings.Repeat("x", 100))},
		"empty and success": {stdout: strings.NewReader("")},
	}
	for name, p := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, err := newTestInspector(cpRunner(p), "").ReadFile(context.Background(), testContainerID, statusPath, 512)
			if err == nil || errors.Is(err, driven.ErrContainerFileNotFound) {
				t.Errorf("err = %v, want an error that is not ErrContainerFileNotFound", err)
			}
		})
	}
}

func TestInspector_ReadFile_TruncatedBody(t *testing.T) {
	t.Parallel()
	full := tarOf(t, tar.Header{Name: "egress", Typeflag: tar.TypeReg}, strings.Repeat("x", 400))
	p := &fakeProcess{stdout: strings.NewReader(full[:512+100])}
	if _, err := newTestInspector(cpRunner(p), "").ReadFile(context.Background(), testContainerID, statusPath, 512); err == nil {
		t.Error("expected an error for an archive cut inside the body")
	}
}

func TestInspector_ReadFile_RejectsBadRequest(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name, id, path string
		max            int
	}{
		{"id with option", "--all", statusPath, 512},
		{"id too short", "abc", statusPath, 512},
		{"id upper-case", strings.ToUpper(testContainerID), statusPath, 512},
		{"relative path", testContainerID, "var/lib/x", 512},
		{"newline in path", testContainerID, "/var/lib/x\ny", 512},
		{"zero limit", testContainerID, statusPath, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := cpRunner(&fakeProcess{stdout: strings.NewReader("")})
			if _, err := newTestInspector(r, "").ReadFile(context.Background(), tc.id, tc.path, tc.max); err == nil {
				t.Error("expected an error")
			}
			if len(r.calls) != 0 {
				t.Errorf("engine called for a rejected request: %v", r.calls)
			}
		})
	}
}

// The adapter's whole engine surface is `ps`, `cp` and `context inspect`: never
// `exec`, and `cp` never follows a link (`-L`).
func TestInspector_EngineSurface(t *testing.T) {
	t.Parallel()
	archive := tarOf(t, tar.Header{Name: "egress", Typeflag: tar.TypeReg}, "u-boot-status 1\nresult=ok\n")
	r := &fakeRunner{byCommand: map[string]func([]string) (docker.Process, error){
		"context": textProcess("unix:///var/run/docker.sock\n"),
		"ps":      textProcess(testContainerID + " running\n"),
		"cp": func([]string) (docker.Process, error) {
			return &fakeProcess{stdout: strings.NewReader(archive)}, nil
		},
	}}
	insp := newTestInspector(r, "")
	if _, _, err := insp.FindContainer(context.Background(), "/home/dev/proj"); err != nil {
		t.Fatal(err)
	}
	if _, err := insp.ReadFile(context.Background(), testContainerID, statusPath, 512); err != nil {
		t.Fatal(err)
	}
	for _, c := range r.calls {
		switch c[0] {
		case "ps", "context":
		case "cp":
			for _, a := range c[1:] {
				if a == "-L" || a == "--follow-link" {
					t.Errorf("docker cp follows links: %v", c)
				}
			}
			if c[len(c)-1] != "-" {
				t.Errorf("docker cp does not stream to stdout: %v", c)
			}
		default:
			t.Errorf("unexpected docker command %q: %v", c[0], c)
		}
	}
}
