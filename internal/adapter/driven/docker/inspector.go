package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// devcontainerFolderLabel is the label the devcontainer CLI sets on the
// container it creates, with the absolute project folder as its value.
const devcontainerFolderLabel = "devcontainer.local_folder"

const (
	// inspectorTimeout bounds every engine call of the inspector: the
	// callers (doctor, generate, init) must not hang on a stale socket.
	inspectorTimeout = 10 * time.Second
	// inspectorOutputLimit bounds the text read from `docker ps` and
	// `docker context inspect`.
	inspectorOutputLimit = 1 << 20
	// inspectorStderrLimit bounds the stderr text kept for classifying a
	// failed command.
	inspectorStderrLimit = 4096
)

var containerIDPattern = regexp.MustCompile(`^[0-9a-f]{12,64}$`)

// Inspector is the production [driven.DevcontainerInspector] adapter. It
// runs only `docker ps`, `docker cp <id>:<path> -` and, to tell a local
// from a remote engine, `docker context inspect`; never `docker exec`
// and never `docker cp -L`. Construct with [NewInspector] or
// [WithInspectorBinary].
type Inspector struct {
	binary string
	getenv func(string) string
	start  startFunc
}

// Static check: Inspector satisfies the DevcontainerInspector port.
var _ driven.DevcontainerInspector = (*Inspector)(nil)

// NewInspector returns an Inspector that shells out to the `docker`
// binary on `$PATH`.
func NewInspector() *Inspector { return WithInspectorBinary("docker") }

// WithInspectorBinary overrides the docker binary path; intended for
// tests.
func WithInspectorBinary(path string) *Inspector {
	return &Inspector{binary: path, getenv: os.Getenv, start: execStart(path)}
}

// process is a started docker command whose stdout the caller reads.
type process interface {
	io.Reader
	// Wait returns the stderr text and the exit error of a command whose
	// stdout was read to EOF.
	Wait() (stderr string, err error)
	// Abort ends the command without regard to its exit status; the
	// caller read all it needs.
	Abort()
}

// startFunc starts the docker command with args.
type startFunc func(ctx context.Context, args ...string) (process, error)

// FindContainer implements [driven.DevcontainerInspector.FindContainer].
func (i *Inspector) FindContainer(ctx context.Context, localFolder string) (driven.DevcontainerContainer, bool, error) {
	if !strings.HasPrefix(localFolder, "/") || strings.ContainsAny(localFolder, "\x00\n\r") {
		return driven.DevcontainerContainer{}, false, fmt.Errorf("docker ps: project folder %q is not an absolute path", localFolder)
	}
	ctx, cancel := context.WithTimeout(ctx, inspectorTimeout)
	defer cancel()
	if err := i.requireLocalEngine(ctx); err != nil {
		return driven.DevcontainerContainer{}, false, err
	}
	out, err := i.output(ctx, "ps", "-a", "--no-trunc",
		"--filter", "label="+devcontainerFolderLabel+"="+localFolder,
		"--format", "{{.ID}} {{.State}}")
	if err != nil {
		return driven.DevcontainerContainer{}, false, fmt.Errorf("docker ps: %w", err)
	}
	first, _, _ := strings.Cut(strings.TrimSpace(out), "\n")
	if first == "" {
		return driven.DevcontainerContainer{}, false, nil
	}
	fields := strings.Fields(first)
	if len(fields) != 2 || !containerIDPattern.MatchString(fields[0]) {
		return driven.DevcontainerContainer{}, false, errors.New("docker ps: unexpected output")
	}
	return driven.DevcontainerContainer{ID: fields[0], Running: fields[1] == "running"}, true, nil
}

// ReadFile implements [driven.DevcontainerInspector.ReadFile]: `docker cp
// <id>:<path> -` streams a tar archive whose first header names the entry
// type and size; only a regular file within the limit is read.
func (i *Inspector) ReadFile(ctx context.Context, containerID, path string, maxBytes int) ([]byte, error) {
	if !containerIDPattern.MatchString(containerID) {
		return nil, fmt.Errorf("docker cp: %q is not a container ID", containerID)
	}
	if !strings.HasPrefix(path, "/") || strings.ContainsAny(path, "\x00\n\r") || maxBytes <= 0 {
		return nil, fmt.Errorf("docker cp: invalid request for %q", path)
	}
	ctx, cancel := context.WithTimeout(ctx, inspectorTimeout)
	defer cancel()
	proc, err := i.start(ctx, "cp", containerID+":"+path, "-")
	if err != nil {
		return nil, fmt.Errorf("docker cp: %w", err)
	}
	defer proc.Abort()
	tr := tar.NewReader(proc)
	hdr, err := tr.Next()
	if errors.Is(err, io.EOF) {
		stderr, waitErr := proc.Wait()
		if waitErr != nil && strings.Contains(stderr, "Could not find the file") {
			return nil, driven.ErrContainerFileNotFound
		}
		if waitErr != nil {
			return nil, fmt.Errorf("docker cp: %w (%s)", waitErr, strings.TrimSpace(stderr))
		}
		return nil, fmt.Errorf("%w: empty archive", driven.ErrContainerFileInvalid)
	}
	if err != nil {
		return nil, fmt.Errorf("docker cp: read archive header: %w", err)
	}
	if hdr.Typeflag != tar.TypeReg {
		return nil, fmt.Errorf("%w: entry type %q is not a regular file", driven.ErrContainerFileInvalid, string(hdr.Typeflag))
	}
	if hdr.Size > int64(maxBytes) {
		return nil, fmt.Errorf("%w: %d bytes exceed the limit of %d", driven.ErrContainerFileInvalid, hdr.Size, maxBytes)
	}
	body, err := io.ReadAll(io.LimitReader(tr, hdr.Size))
	if err != nil {
		return nil, fmt.Errorf("docker cp: read archive body: %w", err)
	}
	if int64(len(body)) != hdr.Size {
		return nil, fmt.Errorf("docker cp: archive ended after %d of %d bytes", len(body), hdr.Size)
	}
	return body, nil
}

// requireLocalEngine fails with [driven.ErrContainerEngineNotLocal] when
// the engine the CLI would talk to is not a local unix socket: a
// container on a remote engine that carries the same `local_folder` value
// belongs to another machine.
func (i *Inspector) requireLocalEngine(ctx context.Context) error {
	host := strings.TrimSpace(i.getenv("DOCKER_HOST"))
	if host == "" {
		out, err := i.output(ctx, "context", "inspect", "--format", "{{.Endpoints.docker.Host}}")
		if err != nil {
			return fmt.Errorf("docker context inspect: %w", err)
		}
		host = strings.TrimSpace(out)
	}
	if host != "" && !strings.HasPrefix(host, "unix://") {
		return driven.ErrContainerEngineNotLocal
	}
	return nil
}

// output runs the command to completion and returns its stdout (at most
// [inspectorOutputLimit] bytes).
func (i *Inspector) output(ctx context.Context, args ...string) (string, error) {
	proc, err := i.start(ctx, args...)
	if err != nil {
		return "", err
	}
	defer proc.Abort()
	out, err := io.ReadAll(io.LimitReader(proc, inspectorOutputLimit))
	if err != nil {
		return "", err
	}
	stderr, err := proc.Wait()
	if err != nil {
		return "", fmt.Errorf("%w (%s)", err, strings.TrimSpace(stderr))
	}
	return string(out), nil
}

// execProcess is the os/exec-backed [process].
type execProcess struct {
	cmd    *exec.Cmd
	stdout io.ReadCloser
	stderr *limitedBuffer
	cancel context.CancelFunc
}

func execStart(binary string) startFunc {
	return func(ctx context.Context, args ...string) (process, error) {
		ctx, cancel := context.WithCancel(ctx)
		cmd := exec.CommandContext(ctx, binary, args...)
		stderr := &limitedBuffer{limit: inspectorStderrLimit}
		cmd.Stderr = stderr
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			cancel()
			return nil, err
		}
		if err := cmd.Start(); err != nil {
			cancel()
			return nil, err
		}
		return &execProcess{cmd: cmd, stdout: stdout, stderr: stderr, cancel: cancel}, nil
	}
}

func (p *execProcess) Read(b []byte) (int, error) { return p.stdout.Read(b) }

func (p *execProcess) Wait() (string, error) {
	err := p.cmd.Wait()
	return p.stderr.String(), err
}

func (p *execProcess) Abort() {
	p.cancel()
	_ = p.cmd.Wait()
}

// limitedBuffer keeps at most limit bytes and discards the rest.
type limitedBuffer struct {
	buf   bytes.Buffer
	limit int
}

func (b *limitedBuffer) Write(p []byte) (int, error) {
	if room := b.limit - b.buf.Len(); room > 0 {
		b.buf.Write(p[:min(room, len(p))])
	}
	return len(p), nil
}

func (b *limitedBuffer) String() string { return b.buf.String() }
