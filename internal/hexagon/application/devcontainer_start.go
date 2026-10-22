package application

import (
	"context"
	"errors"
	"fmt"
	"io"
	iofs "io/fs"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// devcontainerConfigRel is the devcontainer configuration `devcontainer up`
// starts, relative to the project folder.
const devcontainerConfigRel = ".devcontainer/devcontainer.json"

// socketRelaxationCode is the requirement the disclosure of the socket
// relaxation names (LH-FA-DEV-010).
const socketRelaxationCode = "LH-FA-DEV-010"

// socketRelaxationText is the disclosure text. duration names how long the
// helper container holds the socket: "call" for `up`, "start" for `shell`,
// whose session binds no socket. Each of its parts has a test: the helper
// holds the engine socket (`TestRunArgs`), the generated configuration mounts
// it if a feature that does so is enabled, such as `docker-cli`
// (`TestDevcontainerUp_Docker_StandardProfileDockerCLIMountsTheSocket`; the
// project without that feature, `TestDevcontainerUp_Docker_S1_S9_StandardProfile`,
// is the counter-sample), the sandbox profile rejects that feature
// (`TestGenerateDevcontainer_SandboxRejectsHostSocketFeature`), a hand-written
// configuration can mount it (`TestLauncher_Docker_S6_GuardHandWrittenSocketMountIsDetected`),
// and an `initializeCommand` runs with the socket
// (`TestLauncher_Docker_S6_InitializeCommandRunsWithSocket`). The text does
// not say "only": a feature from an external source can mount the socket
// without u-boot seeing it.
func socketRelaxationText(duration string) string {
	return "this command mounts the host container-engine socket into a short-lived helper container " +
		"(root-equivalent access to the engine for the duration of the " + duration + "); " +
		"the devcontainer configuration that u-boot generates mounts it if a feature that does so is enabled, " +
		"such as the built-in feature docker-cli (the sandbox profile rejects that feature); " +
		"a hand-written " + devcontainerConfigRel + " can mount it, " +
		"and an initializeCommand declared there runs with this access"
}

// localEnvRef matches `${localEnv:NAME}` and `${localEnv:NAME:default}` and
// captures NAME.
var localEnvRef = regexp.MustCompile(`\$\{localEnv:([^}:]*)[:}]`)

// localEnvNames returns the sorted, distinct variable names the text
// references through `${localEnv:<NAME>}`. It reads the raw text: a
// reference inside a comment counts as well, and no other source of
// references (features, other configuration files) is read.
func localEnvNames(text string) []string {
	seen := map[string]bool{}
	var names []string
	for _, m := range localEnvRef.FindAllStringSubmatch(text, -1) {
		if !seen[m[1]] {
			seen[m[1]] = true
			names = append(names, m[1])
		}
	}
	sort.Strings(names)
	return names
}

// sessionEnvEntry matches an entry `"<NAME>": "${localEnv:<NAME>}"` and
// captures the name of the key and the name of the reference; the two must be
// equal.
var sessionEnvEntry = regexp.MustCompile(`"([A-Za-z_][A-Za-z0-9_]*)"\s*:\s*"\$\{localEnv:([A-Za-z_][A-Za-z0-9_]*)\}"`)

// sessionEnvNames returns the sorted, distinct variable names of the entries
// `"<NAME>": "${localEnv:<NAME>}"` in the text of devcontainer.json: the form
// that the generated configuration writes for `remoteEnv`. It reads the text
// without its comments (an entry that is commented out passes nothing), finds
// the form wherever it stands, not only below `remoteEnv`, and does not ask
// whether a name is set on the host (the launcher leaves an unset name out).
// A reference of another form (a different name, a default, a path in `mounts`)
// is no entry; no other source of references is read.
func sessionEnvNames(text string) []string {
	seen := map[string]bool{}
	var names []string
	for _, m := range sessionEnvEntry.FindAllStringSubmatch(withoutComments(text), -1) {
		if m[1] == m[2] && !seen[m[1]] {
			seen[m[1]] = true
			names = append(names, m[1])
		}
	}
	sort.Strings(names)
	return names
}

// withoutComments returns text without the comments of JSON with comments:
// `//` up to the end of the line and `/* … */`. Characters inside a string are
// kept, so a `//` in a URL starts no comment. A line comment leaves its line
// break and a block comment one space, so the tokens next to it stay apart; an
// unterminated block comment runs to the end of the text.
func withoutComments(text string) string {
	var b strings.Builder
	for i := 0; i < len(text); {
		switch {
		case text[i] == '"':
			end := stringEnd(text, i)
			b.WriteString(text[i:end])
			i = end
		case strings.HasPrefix(text[i:], "//"):
			if nl := strings.IndexByte(text[i:], '\n'); nl < 0 {
				i = len(text)
			} else {
				i += nl
			}
		case strings.HasPrefix(text[i:], "/*"):
			if end := strings.Index(text[i+2:], "*/"); end < 0 {
				i = len(text)
			} else {
				i += end + 4
			}
			b.WriteByte(' ')
		default:
			b.WriteByte(text[i])
			i++
		}
	}
	return b.String()
}

// stringEnd returns the index after the string that begins with the quote at
// text[start]; an unterminated string ends at the end of the text. A backslash
// escapes the next character.
func stringEnd(text string, start int) int {
	for i := start + 1; i < len(text); i++ {
		switch text[i] {
		case '\\':
			i++
		case '"':
			return i + 1
		}
	}
	return len(text)
}

// DevcontainerStartService implements [driving.DevcontainerUseCase]
// (`u-boot devcontainer up`, `down` and `shell`, LH-FA-DEV-010). It reads the
// project files through the [driven.FileSystem] port, starts, stops and enters
// through the [driven.DevcontainerLauncher] and finds the container and its
// status files through the [driven.DevcontainerInspector].
type DevcontainerStartService struct {
	fs        driven.FileSystem
	yaml      driven.YAMLCodec
	launcher  driven.DevcontainerLauncher
	inspector driven.DevcontainerInspector
	logger    driven.Logger
}

// Static check: DevcontainerStartService satisfies the driving port.
var _ driving.DevcontainerUseCase = (*DevcontainerStartService)(nil)

// NewDevcontainerStartService constructs the service. Launcher and inspector
// are mandatory ports without a nil-tolerant default: a missing one makes
// `Up`, `Down` and `Shell` fail instead of reporting "nothing to do". A nil logger
// is routed to [noopLogger].
func NewDevcontainerStartService(
	fs driven.FileSystem,
	yaml driven.YAMLCodec,
	launcher driven.DevcontainerLauncher,
	inspector driven.DevcontainerInspector,
	logger driven.Logger,
) *DevcontainerStartService {
	if logger == nil {
		logger = noopLogger{}
	}
	return &DevcontainerStartService{fs: fs, yaml: yaml, launcher: launcher, inspector: inspector, logger: logger}
}

func (s *DevcontainerStartService) wired(command, baseDir string) error {
	switch {
	case baseDir == "":
		return fmt.Errorf("devcontainer %s: BaseDir is required", command)
	case s.launcher == nil:
		return fmt.Errorf("devcontainer %s: launcher not wired", command)
	case s.inspector == nil:
		return fmt.Errorf("devcontainer %s: inspector not wired", command)
	}
	return nil
}

// Up implements [driving.DevcontainerUseCase.Up].
func (s *DevcontainerStartService) Up(ctx context.Context, req driving.DevcontainerUpRequest) (driving.DevcontainerUpResponse, error) {
	resp, _, err := s.start(ctx, startRequest{
		command: "up", baseDir: req.BaseDir, recreate: req.Recreate, reportRunning: true,
		duration: "call", progress: req.ProgressSink, onWarning: req.OnWarning,
	})
	return resp, err
}

// Shell implements [driving.DevcontainerUseCase.Shell]: the start of Up
// without Recreate, then one session in the started container. The fallback
// warnings of the start go to OnWarning before the session begins, because the
// call returns only when the session has ended. The variables of the session
// come from the entries of the configuration that the start has just read.
func (s *DevcontainerStartService) Shell(ctx context.Context, req driving.DevcontainerShellRequest) error {
	resp, config, err := s.start(ctx, startRequest{
		command: "shell", baseDir: req.BaseDir, duration: "start", progress: req.ProgressSink, onWarning: req.OnWarning,
	})
	if err != nil {
		return err
	}
	if req.OnWarning != nil {
		for _, w := range resp.Warnings {
			req.OnWarning(w)
		}
	}
	names := sessionEnvNames(config)
	s.logger.Debug("devcontainer shell: opening the session", "envNames", len(names), "terminal", req.Terminal)
	return s.launcher.Shell(ctx, driven.DevcontainerShellOptions{
		ContainerID: resp.ContainerID,
		User:        resp.RemoteUser,
		Workspace:   resp.RemoteWorkspaceFolder,
		EnvNames:    names,
		Stdin:       req.Stdin,
		Stdout:      req.Stdout,
		Stderr:      req.Stderr,
		Terminal:    req.Terminal,
	})
}

// startRequest is what `up` and `shell` hand to [DevcontainerStartService.start].
type startRequest struct {
	// command names the command in error texts: "up" or "shell".
	command  string
	baseDir  string
	recreate bool
	// reportRunning asks the inspector, before the launcher, whether the
	// container already runs (ADR-0030 decision 2); only `up` reports it.
	reportRunning bool
	// duration completes the disclosure of the socket, see [socketRelaxationText].
	duration  string
	progress  io.Writer
	onWarning func(driving.WarningEntry)
}

// start is the start phase that `up` and `shell` share (ADR-0031 decision 1):
// validation, the disclosure of the socket, the launcher, and the judgement by
// the state of the container ([DevcontainerStartService.startOutcome]). It
// returns the response and the text of devcontainer.json it started from.
func (s *DevcontainerStartService) start(ctx context.Context, r startRequest) (driving.DevcontainerUpResponse, string, error) {
	if err := s.wired(r.command, r.baseDir); err != nil {
		return driving.DevcontainerUpResponse{}, "", err
	}
	if driven.ValidateDevcontainerProjectFolder(r.baseDir) != nil {
		return driving.DevcontainerUpResponse{}, "", fmt.Errorf("devcontainer %s: %w: the path must be absolute and hold no comma, double quote or line break",
			r.command, driving.ErrDevcontainerProjectFolderInvalid)
	}
	config, err := s.readConfig(r.command, r.baseDir)
	if err != nil {
		return driving.DevcontainerUpResponse{}, "", err
	}
	names, err := referencedEnvNames(r.command, config)
	if err != nil {
		return driving.DevcontainerUpResponse{}, "", err
	}
	if r.onWarning != nil {
		r.onWarning(driving.WarningEntry{
			Code: socketRelaxationCode, Level: "warn", Message: socketRelaxationText(r.duration), SecurityRelaxation: true,
		})
	}
	s.logger.Debug("devcontainer "+r.command, "recreate", r.recreate, "envNames", len(names))
	runningBefore := ""
	if r.reportRunning {
		runningBefore = s.runningContainerID(ctx, r.baseDir, r.recreate)
	}
	res, err := s.launcher.Up(ctx, driven.DevcontainerUpOptions{
		ProjectFolder: r.baseDir,
		EnvNames:      names,
		Recreate:      r.recreate,
		Progress:      r.progress,
	})
	resp, err := s.startOutcome(ctx, r.command, r.baseDir, runningBefore, res, err)
	return resp, config, err
}

// runningContainerID returns the ID of the project's container if it runs
// before the launcher is called (ADR-0030 decision 2), else the empty string:
// with `--recreate` (the container is removed, so the query is skipped), with
// no container, a stopped one, and when the engine cannot be asked. A failed
// query is no failure of the command; the launcher reports an engine problem
// itself.
func (s *DevcontainerStartService) runningContainerID(ctx context.Context, baseDir string, recreate bool) string {
	if recreate {
		return ""
	}
	container, found, err := s.inspector.FindContainer(ctx, baseDir)
	if err != nil || !found || !container.Running {
		return ""
	}
	return container.ID
}

// startOutcome judges the result of the launcher by the state of the container
// (ADR-0030 decision 1). Only a start that succeeded or failed with the runtime
// class reads the status files, through one [fallbackReader]: an engine
// problem or an abort keeps its class and asks nothing. A `result=failed` of a
// configured capability ends the call with
// [driving.ErrDevcontainerStartScriptFailed], whether the CLI failed (the
// message keeps its output) or succeeded (the file stems from an earlier start
// of the container). Otherwise the response carries the fallback warnings and
// whether the container ran before the call. command names the command in the
// error texts.
func (s *DevcontainerStartService) startOutcome(
	ctx context.Context, command, baseDir, runningBefore string, res driven.DevcontainerUpResult, launchErr error,
) (driving.DevcontainerUpResponse, error) {
	if launchErr != nil && !errors.Is(launchErr, driven.ErrDevcontainerRuntime) {
		return driving.DevcontainerUpResponse{}, launchErr
	}
	probes := s.fallbackProbes(baseDir)
	reader := &fallbackReader{inspector: s.inspector, baseDir: baseDir}
	failedProbe, failedReading, failed := firstFailedScript(ctx, reader, probes)
	switch {
	case launchErr != nil && failed:
		return driving.DevcontainerUpResponse{}, fmt.Errorf("devcontainer %s: %w: %v", command, driving.ErrDevcontainerStartScriptFailed, launchErr)
	case launchErr != nil:
		return driving.DevcontainerUpResponse{}, launchErr
	case failed:
		return driving.DevcontainerUpResponse{}, fmt.Errorf("devcontainer %s: %w: %s",
			command, driving.ErrDevcontainerStartScriptFailed, failedAfterStartText(failedProbe, failedReading))
	}
	return driving.DevcontainerUpResponse{
		ContainerID:           res.ContainerID,
		RemoteUser:            res.RemoteUser,
		RemoteWorkspaceFolder: res.RemoteWorkspaceFolder,
		Warnings:              fallbackWarnings(ctx, reader, probes, startFallbackWording()),
		AlreadyRunning:        runningBefore != "" && runningBefore == res.ContainerID,
	}, nil
}

// readConfig returns the text of `.devcontainer/devcontainer.json`. A missing
// file is [driving.ErrDevcontainerConfigMissing], a read failure
// [driving.ErrDevcontainerFileSystem]. command names the command in the error
// texts.
func (s *DevcontainerStartService) readConfig(command, baseDir string) (string, error) {
	path := filepath.Join(baseDir, devcontainerConfigRel)
	exists, err := s.fs.Exists(path)
	if err != nil {
		return "", fmt.Errorf("devcontainer %s: probe %q: %w: %w", command, path, driving.ErrDevcontainerFileSystem, err)
	}
	if !exists {
		return "", fmt.Errorf("devcontainer %s: %w (looked for %s)", command, driving.ErrDevcontainerConfigMissing, devcontainerConfigRel)
	}
	body, err := s.fs.ReadFile(path)
	if err != nil {
		if errors.Is(err, iofs.ErrNotExist) {
			return "", fmt.Errorf("devcontainer %s: %w (looked for %s)", command, driving.ErrDevcontainerConfigMissing, devcontainerConfigRel)
		}
		return "", fmt.Errorf("devcontainer %s: read %q: %w: %w", command, path, driving.ErrDevcontainerFileSystem, err)
	}
	return string(body), nil
}

// referencedEnvNames returns the variable names the text of devcontainer.json
// references through `${localEnv:<NAME>}`, for the helper container of the
// start. A name the launcher could not pass is
// [driving.ErrDevcontainerEnvNameInvalid]. command names the command in the
// error text.
func referencedEnvNames(command, text string) ([]string, error) {
	names := localEnvNames(text)
	for _, name := range names {
		if driven.ValidateDevcontainerEnvName(name) != nil {
			return nil, fmt.Errorf("devcontainer %s: %w: %q in %s", command, driving.ErrDevcontainerEnvNameInvalid, name, devcontainerConfigRel)
		}
	}
	return names, nil
}

// fallbackProbes returns the status files of the capabilities the
// project's `u-boot.yaml` asks for (LH-FA-DEV-007). Without the file, with
// an unreadable or unparsable one, and without the sandbox profile it
// returns none: the fallback disclosure and the exit code 11 of a start
// script abort need the configuration, and their absence is no error here.
func (s *DevcontainerStartService) fallbackProbes(baseDir string) []fallbackProbe {
	path := filepath.Join(baseDir, "u-boot.yaml")
	if exists, err := s.fs.Exists(path); err != nil || !exists {
		return nil
	}
	body, err := s.fs.ReadFile(path)
	if err != nil {
		return nil
	}
	var cfg ubootYAMLConfig
	if err := s.yaml.Unmarshal(body, &cfg); err != nil {
		return nil
	}
	return configuredFallbackProbes(cfg.Devcontainer, profileIsSandbox(cfg.Devcontainer))
}

// firstFailedScript returns the first of the probes, in their order, whose
// status file reports `result=failed`, with the reading (ADR-0029 decision 1,
// ADR-0030 decision 1). The file describes the last run of its script in the
// container, not necessarily this call.
func firstFailedScript(ctx context.Context, r *fallbackReader, probes []fallbackProbe) (fallbackProbe, statusReading, bool) {
	for _, p := range probes {
		if rd := r.read(ctx, p); rd.outcome == statusReported && rd.status.Result() == domain.SandboxStatusFailed {
			return p, rd, true
		}
	}
	return fallbackProbe{}, statusReading{}, false
}

// failedAfterStartText is the message of a `result=failed` that the devcontainer
// CLI's successful run leaves standing (ADR-0030 decision 1): the sentence of
// the capability (the one `doctor` shows), the time of the report, that the
// container runs and this call did not run the start script, and the way
// out.
func failedAfterStartText(p fallbackProbe, rd statusReading) string {
	return p.failed + " (reported at " + rd.status.At().Format(time.RFC3339) + "); " +
		"the container is running and this call did not run the start script; " +
		"`u-boot devcontainer up --recreate` creates the container anew"
}

// Down implements [driving.DevcontainerUseCase.Down].
func (s *DevcontainerStartService) Down(ctx context.Context, req driving.DevcontainerDownRequest) (driving.DevcontainerDownResponse, error) {
	if err := s.wired("down", req.BaseDir); err != nil {
		return driving.DevcontainerDownResponse{}, err
	}
	container, found, err := s.inspector.FindContainer(ctx, req.BaseDir)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return driving.DevcontainerDownResponse{}, ctxErr
		}
		if errors.Is(err, driven.ErrContainerEngineNotLocal) {
			return driving.DevcontainerDownResponse{}, fmt.Errorf("devcontainer down: %w", err)
		}
		return driving.DevcontainerDownResponse{}, fmt.Errorf("devcontainer down: the container engine could not be queried: %w (%s)",
			driven.ErrDockerUnavailable, err)
	}
	if !found {
		return driving.DevcontainerDownResponse{}, nil
	}
	if !container.Running {
		return driving.DevcontainerDownResponse{ContainerID: container.ID}, nil
	}
	if err := s.launcher.Stop(ctx, container.ID); err != nil {
		return driving.DevcontainerDownResponse{}, err
	}
	return driving.DevcontainerDownResponse{ContainerID: container.ID, Stopped: true}, nil
}
