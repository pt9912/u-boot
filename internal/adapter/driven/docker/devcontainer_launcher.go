package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

const (
	// helperImageRepo is the repository part of the helper image tag; the
	// tag part derives from the CLI version and the context hash.
	helperImageRepo = "u-boot-devcontainer-cli"
	// helperContainerPrefix starts the name of every helper container, so
	// a leftover is recognisable by name.
	helperContainerPrefix = "u-boot-devcontainer-up-"
	// dockerSocketPath is the engine socket on the machine of the daemon;
	// it is mounted at the same path into the helper container.
	dockerSocketPath = "/var/run/docker.sock"
	// cliPackageName is the dependency in package.json whose exact version
	// is the CLI pin.
	cliPackageName = "@devcontainers/cli"

	// defaultAbortGrace is the time `docker run` gets after SIGINT before it is
	// killed.
	defaultAbortGrace = 10 * time.Second
	// helperRemoveTimeout bounds `docker rm -f` of a helper container after
	// a failed or cancelled run.
	helperRemoveTimeout = 30 * time.Second
	// upStdoutLimit bounds the stdout kept from the helper container (its
	// tail: the CLI prints the result line last).
	upStdoutLimit = 1 << 20
	// upStderrTail is the stderr tail kept for error messages.
	upStderrTail = 4096
)

// helperContextFS is the build context of the helper image: Dockerfile,
// package.json and package-lock.json.
//
//go:embed devcontainercli/Dockerfile devcontainercli/package.json devcontainercli/package-lock.json
var helperContextFS embed.FS

// helperContextNames lists the files of the build context in tar order.
func helperContextNames() []string {
	return []string{"Dockerfile", "package.json", "package-lock.json"}
}

// Launcher is the production [driven.DevcontainerLauncher] adapter. It runs
// `docker build -`, `docker run`, `docker stop`, `docker rm -f` and the
// environment checks; it never uses Compose. Construct with [NewLauncher]
// or [WithLauncherBinary].
type Launcher struct {
	binary     string
	lookupEnv  func(string) (string, bool)
	probe      *Probe
	abortGrace time.Duration
	// namePrefix starts the name of every helper container this Launcher
	// starts; the default is helperContainerPrefix.
	namePrefix string
	// imageRepo is the repository part of the helper image tag; the default
	// is helperImageRepo. A test that removes the image uses its own
	// repository, so it neither removes nor races with the image another
	// package of the same test run starts containers from.
	imageRepo string
	// buildArgs are inserted into `docker build` before the tag; empty by
	// default.
	buildArgs []string
	output    func(ctx context.Context, args ...string) (string, error)
}

// Static check: Launcher satisfies the DevcontainerLauncher port.
var _ driven.DevcontainerLauncher = (*Launcher)(nil)

// NewLauncher returns a Launcher that shells out to the `docker` binary on
// `$PATH`.
func NewLauncher() *Launcher { return WithLauncherBinary("docker") }

// WithLauncherBinary overrides the docker binary path; intended for tests.
func WithLauncherBinary(path string) *Launcher {
	return &Launcher{
		binary:     path,
		lookupEnv:  os.LookupEnv,
		probe:      WithBinary(path),
		abortGrace: defaultAbortGrace,
		namePrefix: helperContainerPrefix,
		imageRepo:  helperImageRepo,
		output:     WithInspectorBinary(path).output,
	}
}

func (l *Launcher) getenv(key string) string {
	value, _ := l.lookupEnv(key)
	return value
}

// preflight is the two-plus-one environment check: the Docker binary is on
// the PATH, the daemon answers, the engine is local. The Compose plugin is
// not needed, so [Engine.preflight] does not apply.
func (l *Launcher) preflight(ctx context.Context) error {
	if _, err := exec.LookPath(l.binary); err != nil {
		return fmt.Errorf("docker binary not on PATH (%s): %w", err.Error(), driven.ErrDockerUnavailable)
	}
	if err := l.probe.Info(ctx); err != nil {
		return fmt.Errorf("docker daemon unreachable (%s): %w", err.Error(), driven.ErrDockerUnavailable)
	}
	if err := requireLocalEngine(ctx, l.getenv, l.output); err != nil {
		if errors.Is(err, driven.ErrContainerEngineNotLocal) {
			return err
		}
		return fmt.Errorf("docker engine not determinable (%s): %w", err.Error(), driven.ErrDockerUnavailable)
	}
	return nil
}

// Up implements [driven.DevcontainerLauncher.Up].
func (l *Launcher) Up(ctx context.Context, opts driven.DevcontainerUpOptions) (driven.DevcontainerUpResult, error) {
	if err := driven.ValidateDevcontainerProjectFolder(opts.ProjectFolder); err != nil {
		return driven.DevcontainerUpResult{}, rejectedInput(err)
	}
	envNames, err := l.passedEnvNames(opts.EnvNames)
	if err != nil {
		return driven.DevcontainerUpResult{}, err
	}
	if err := l.preflight(ctx); err != nil {
		return driven.DevcontainerUpResult{}, err
	}
	image, err := l.EnsureImage(ctx, opts.Progress)
	if err != nil {
		return driven.DevcontainerUpResult{}, err
	}
	return l.runUp(ctx, image, l.namePrefix+randomSuffix(), envNames, opts)
}

// Stop implements [driven.DevcontainerLauncher.Stop].
func (l *Launcher) Stop(ctx context.Context, containerID string) error {
	if !containerIDPattern.MatchString(containerID) {
		return fmt.Errorf("docker stop: %q is not a container ID: %w", containerID, driven.ErrDevcontainerRuntime)
	}
	if err := l.preflight(ctx); err != nil {
		return err
	}
	out, err := exec.CommandContext(ctx, l.binary, "stop", containerID).CombinedOutput()
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return fmt.Errorf("docker stop failed (%s: %s): %w", err.Error(), tailOf(out), driven.ErrDevcontainerRuntime)
	}
	return nil
}

// EnsureImage returns the tag of the helper image and builds the image from
// the embedded context when the tag is missing. The build output goes to
// progress (nil discards it). `--force-rm` removes the intermediate
// container of a failed step, so a failed build leaves no container.
func (l *Launcher) EnsureImage(ctx context.Context, progress io.Writer) (string, error) {
	tag, err := helperImageTag(l.imageRepo)
	if err != nil {
		return "", fmt.Errorf("helper image context: %s: %w", err.Error(), driven.ErrDevcontainerRuntime)
	}
	if exec.CommandContext(ctx, l.binary, "image", "inspect", "--format", "{{.Id}}", tag).Run() == nil {
		return tag, nil
	}
	if ctxErr := ctx.Err(); ctxErr != nil {
		return "", ctxErr
	}
	buildContext, err := helperContextTar()
	if err != nil {
		return "", fmt.Errorf("helper image context: %s: %w", err.Error(), driven.ErrDevcontainerRuntime)
	}
	stderr := &tailBuffer{limit: upStderrTail}
	args := append([]string{"build", "--force-rm"}, l.buildArgs...)
	cmd := exec.CommandContext(ctx, l.binary, append(args, "-t", tag, "-")...)
	cmd.Stdin = bytes.NewReader(buildContext)
	cmd.Stdout = progressSinkOrDiscard(progress)
	cmd.Stderr = io.MultiWriter(stderr, progressSinkOrDiscard(progress))
	if err := cmd.Run(); err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return "", ctxErr
		}
		return "", fmt.Errorf("docker build of the helper image failed (%s: %s): %w", err.Error(), stderr.String(), driven.ErrDevcontainerRuntime)
	}
	return tag, nil
}

// runUp starts the helper container and reads the result line. A helper
// container that outlives a failed or cancelled run is removed by name.
// stdout is kept unmasked and only parsed: the result line is data for the
// caller, and the error texts that quote it are masked. stderr goes through
// the masking writer before it reaches the tail buffer and the progress.
func (l *Launcher) runUp(ctx context.Context, image, name string, envNames []string, opts driven.DevcontainerUpOptions) (driven.DevcontainerUpResult, error) {
	stdout := &tailBuffer{limit: upStdoutLimit}
	stderr := &tailBuffer{limit: upStderrTail}
	values := l.passedValues(envNames)
	maskedStderr := newMaskingWriter(io.MultiWriter(stderr, progressSinkOrDiscard(opts.Progress)), values)
	cmd := exec.CommandContext(ctx, l.binary, runArgs(name, image, envNames, opts)...)
	cmd.Cancel = func() error { return cmd.Process.Signal(os.Interrupt) }
	cmd.WaitDelay = l.abortGrace
	cmd.Stdout = stdout
	cmd.Stderr = maskedStderr
	runErr := cmd.Run()
	_ = maskedStderr.Flush()
	if runErr != nil {
		detail := ""
		if rmErr := l.removeHelper(ctx, name); rmErr != nil {
			detail = "; the helper container may still run: " + rmErr.Error()
		}
		if ctxErr := ctx.Err(); ctxErr != nil {
			if detail != "" {
				return driven.DevcontainerUpResult{}, fmt.Errorf("%w%s", ctxErr, detail)
			}
			return driven.DevcontainerUpResult{}, ctxErr
		}
		return driven.DevcontainerUpResult{}, fmt.Errorf("devcontainer up failed (%s: %s)%s: %w", runErr.Error(), stderr.String(), detail, driven.ErrDevcontainerRuntime)
	}
	result, err := parseUpResult(stdout.String())
	if err != nil {
		return driven.DevcontainerUpResult{}, fmt.Errorf("devcontainer up: %s (%s): %w", maskString(err.Error(), values), stderr.String(), driven.ErrDevcontainerRuntime)
	}
	return result, nil
}

// removeHelper removes the helper container by name, also when ctx is
// cancelled; an absent container is not an error, any other failure is
// returned.
func (l *Launcher) removeHelper(ctx context.Context, name string) error {
	rmCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), helperRemoveTimeout)
	defer cancel()
	out, err := exec.CommandContext(rmCtx, l.binary, "rm", "-f", name).CombinedOutput()
	if err != nil && !strings.Contains(strings.ToLower(string(out)), "no such container") {
		return fmt.Errorf("docker rm -f %s failed (%s: %s)", name, err.Error(), tailOf(out))
	}
	return nil
}

// passedValues returns the host values of the passed variables; the
// launcher masks them in everything it forwards or reports.
func (l *Launcher) passedValues(names []string) []string {
	values := make([]string, 0, len(names))
	for _, name := range names {
		values = append(values, l.getenv(name))
	}
	return values
}

// rejectedInput classifies a request the launcher refuses before any docker
// call.
func rejectedInput(err error) error {
	return fmt.Errorf("%s: %w", err.Error(), driven.ErrDevcontainerRuntime)
}

// passedEnvNames validates the requested names and keeps those set on the
// host.
func (l *Launcher) passedEnvNames(names []string) ([]string, error) {
	var passed []string
	for _, name := range names {
		if err := driven.ValidateDevcontainerEnvName(name); err != nil {
			return nil, rejectedInput(err)
		}
		if _, ok := l.lookupEnv(name); ok {
			passed = append(passed, name)
		}
	}
	return passed, nil
}

// runArgs is the argument list of the helper container: removed at the end
// of the run, the engine socket and the project folder (same path as on
// the host) as `--mount`, variables by name only, no privileges.
func runArgs(name, image string, envNames []string, opts driven.DevcontainerUpOptions) []string {
	folder := opts.ProjectFolder
	folderMount := "type=bind,source=" + folder + ",target=" + folder
	if opts.ReadOnlyProject {
		folderMount += ",readonly"
	}
	args := []string{
		"run", "--rm", "--name", name,
		"--mount", "type=bind,source=" + dockerSocketPath + ",target=" + dockerSocketPath,
		"--mount", folderMount,
	}
	for _, envName := range envNames {
		args = append(args, "-e", envName)
	}
	args = append(args, image, "devcontainer", "up", "--workspace-folder", folder)
	if opts.Recreate {
		args = append(args, "--remove-existing-container")
	}
	return args
}

// upResultLine is the result line the devcontainer CLI prints as the last
// non-empty line of its stdout.
type upResultLine struct {
	Outcome               string `json:"outcome"`
	ContainerID           string `json:"containerId"`
	RemoteUser            string `json:"remoteUser"`
	RemoteWorkspaceFolder string `json:"remoteWorkspaceFolder"`
	Message               string `json:"message"`
}

// parseUpResult reads the last non-empty line of stdout as the result of
// `devcontainer up`.
func parseUpResult(stdout string) (driven.DevcontainerUpResult, error) {
	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	last := strings.TrimSpace(lines[len(lines)-1])
	if last == "" {
		return driven.DevcontainerUpResult{}, errors.New("no result line on stdout")
	}
	var line upResultLine
	if err := json.Unmarshal([]byte(last), &line); err != nil {
		return driven.DevcontainerUpResult{}, errors.New("last stdout line is not a JSON result")
	}
	if line.Outcome != "success" {
		return driven.DevcontainerUpResult{}, fmt.Errorf("outcome %q: %s", line.Outcome, line.Message)
	}
	if !containerIDPattern.MatchString(line.ContainerID) {
		return driven.DevcontainerUpResult{}, errors.New("result line carries no container ID")
	}
	return driven.DevcontainerUpResult{
		ContainerID:           line.ContainerID,
		RemoteUser:            line.RemoteUser,
		RemoteWorkspaceFolder: line.RemoteWorkspaceFolder,
	}, nil
}

// helperContextFiles returns the embedded build context by file name.
func helperContextFiles() (map[string][]byte, error) {
	files := make(map[string][]byte)
	for _, name := range helperContextNames() {
		data, err := helperContextFS.ReadFile("devcontainercli/" + name)
		if err != nil {
			return nil, err
		}
		files[name] = data
	}
	return files, nil
}

// helperCLIVersion is the exact CLI version pinned in package.json.
func helperCLIVersion(packageJSON []byte) (string, error) {
	var pkg struct {
		Dependencies map[string]string `json:"dependencies"`
	}
	if err := json.Unmarshal(packageJSON, &pkg); err != nil {
		return "", err
	}
	version := pkg.Dependencies[cliPackageName]
	if version == "" {
		return "", fmt.Errorf("package.json pins no %s", cliPackageName)
	}
	return version, nil
}

// helperImageTag is `<repo>:<CLI version>-<first 12 hex of the context
// hash>`: a changed pin changes the tag, so the image is built anew.
func helperImageTag(repo string) (string, error) {
	files, err := helperContextFiles()
	if err != nil {
		return "", err
	}
	version, err := helperCLIVersion(files["package.json"])
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	for _, name := range helperContextNames() {
		fmt.Fprintf(hash, "%s\x00%d\x00", name, len(files[name]))
		hash.Write(files[name])
	}
	return fmt.Sprintf("%s:%s-%s", repo, version, hex.EncodeToString(hash.Sum(nil))[:12]), nil
}

// helperContextTar is the build context as a tar stream with fixed
// metadata, so equal files give equal bytes.
func helperContextTar() ([]byte, error) {
	files, err := helperContextFiles()
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	for _, name := range helperContextNames() {
		hdr := &tar.Header{Name: name, Mode: 0o644, Size: int64(len(files[name])), ModTime: time.Unix(0, 0)}
		if err := tw.WriteHeader(hdr); err != nil {
			return nil, err
		}
		if _, err := tw.Write(files[name]); err != nil {
			return nil, err
		}
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// randomSuffix returns eight random hex characters for a container name.
func randomSuffix() string {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%08x", time.Now().UnixNano()&0xffffffff)
	}
	return hex.EncodeToString(b)
}

// tailOf returns the last upStderrTail bytes of out, trimmed.
func tailOf(out []byte) string {
	if len(out) > upStderrTail {
		out = out[len(out)-upStderrTail:]
	}
	return strings.TrimSpace(string(out))
}

// tailBuffer keeps the last limit bytes and discards the earlier ones; the
// backing slice is compacted once it holds twice the limit.
type tailBuffer struct {
	data  []byte
	limit int
}

func (b *tailBuffer) Write(p []byte) (int, error) {
	b.data = append(b.data, p...)
	if len(b.data) > 2*b.limit {
		n := copy(b.data, b.data[len(b.data)-b.limit:])
		b.data = b.data[:n]
	}
	return len(p), nil
}

func (b *tailBuffer) String() string {
	data := b.data
	if len(data) > b.limit {
		data = data[len(data)-b.limit:]
	}
	return strings.TrimSpace(string(data))
}
