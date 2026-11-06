package application

import (
	"context"
	"encoding/json"
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
// (`u-boot devcontainer up`, `down`, `shell` and `remove`, LH-FA-DEV-010). It
// reads the project files through the [driven.FileSystem] port, starts, stops,
// enters and removes through the [driven.DevcontainerLauncher], finds the
// container, its status files and its volumes through the
// [driven.DevcontainerInspector], asks the confirmation of `remove --volumes`
// through the [driven.Confirmer] and, when wired, resolves the git token and
// the git identity of the host through the [driven.HostGitSource] port
// (LH-FA-DEV-009, LH-FA-DEV-010).
type DevcontainerStartService struct {
	fs         driven.FileSystem
	yaml       driven.YAMLCodec
	launcher   driven.DevcontainerLauncher
	inspector  driven.DevcontainerInspector
	hostSource driven.HostGitSource
	confirmer  driven.Confirmer
	logger     driven.Logger
}

// Static check: DevcontainerStartService satisfies the driving port.
var _ driving.DevcontainerUseCase = (*DevcontainerStartService)(nil)

// NewDevcontainerStartService constructs the service. Launcher and inspector
// are mandatory ports without a nil-tolerant default: a missing one makes
// `Up`, `Down`, `Shell` and `Remove` fail instead of reporting "nothing to do". A nil logger
// is routed to [noopLogger]. A nil host source switches the host git
// resolution off: the caller of the construction decides; the production
// wiring sets it.
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
	return &DevcontainerStartService{fs: fs, yaml: yaml, launcher: launcher, inspector: inspector, confirmer: noopConfirmer{}, logger: logger}
}

// WithConfirmer wires the confirmation of `remove --volumes`
// ([driven.Confirmer]); a nil value leaves the default, which declines.
func (s *DevcontainerStartService) WithConfirmer(confirmer driven.Confirmer) *DevcontainerStartService {
	if confirmer != nil {
		s.confirmer = confirmer
	}
	return s
}

// WithHostGitSource wires the host git resolution
// ([driven.HostGitSource]); a nil value switches it off.
func (s *DevcontainerStartService) WithHostGitSource(hostSource driven.HostGitSource) *DevcontainerStartService {
	s.hostSource = hostSource
	return s
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
		command: "up", baseDir: req.BaseDir, instance: req.Instance, recreate: req.Recreate, reportRunning: true,
		duration: "call", progress: req.ProgressSink, onWarning: req.OnWarning,
	})
	return resp, err
}

// Shell implements [driving.DevcontainerUseCase.Shell]: the start of Up
// without Recreate, then one session in the started container. The fallback
// warnings of the start go to OnWarning before the session begins, because the
// call returns only when the session has ended; the response goes to OnStart
// after them, so the caller prints the container of the session beside its ID
// before the session owns the standard output. The variables of the session
// come from the entries of the configuration that the start has just read.
func (s *DevcontainerStartService) Shell(ctx context.Context, req driving.DevcontainerShellRequest) error {
	resp, resolution, err := s.start(ctx, startRequest{
		command: "shell", baseDir: req.BaseDir, instance: req.Instance, duration: "start",
		progress: req.ProgressSink, onWarning: req.OnWarning, onStart: req.OnStart,
	})
	if err != nil {
		return err
	}
	if req.OnWarning != nil {
		for _, w := range resp.Warnings {
			req.OnWarning(w)
		}
	}
	if req.OnStart != nil {
		req.OnStart(resp)
	}
	names := sessionEnvNames(resolution.config)
	s.logger.Debug("devcontainer shell: opening the session", "envNames", len(names), "terminal", req.Terminal)
	return s.launcher.Shell(ctx, driven.DevcontainerShellOptions{
		ContainerID: resp.ContainerID,
		User:        resp.RemoteUser,
		Workspace:   resp.RemoteWorkspaceFolder,
		EnvNames:    names,
		// The session receives the resolved values of the start; a token
		// resolved there reaches the session (LH-FA-DEV-009, LH-FA-DEV-010).
		EnvValues: resolution.envValues,
		Stdin:     req.Stdin,
		Stdout:    req.Stdout,
		Stderr:    req.Stderr,
		Terminal:  req.Terminal,
	})
}

// startResolution is what the start keeps for the caller: the text of
// devcontainer.json it started from and the resolved values of the host
// git sources; the session of `shell` receives the values (LH-FA-DEV-009,
// LH-FA-DEV-010).
type startResolution struct {
	config string
	// envValues carries the resolved token (the key is the name in the
	// configuration); nil-free: an empty map resolves nothing.
	envValues map[string]string
}

// startRequest is what `up` and `shell` hand to [DevcontainerStartService.start].
type startRequest struct {
	// command names the command in error texts: "up" or "shell".
	command string
	baseDir string
	// instance is the instance name; empty selects the container without one.
	instance string
	recreate bool
	// reportRunning asks the inspector, before the launcher, whether the
	// container already runs (ADR-0030 decision 2); only `up` reports it.
	reportRunning bool
	// duration completes the disclosure of the socket, see [socketRelaxationText].
	duration  string
	progress  io.Writer
	onWarning func(driving.WarningEntry)
	// onStart receives the response of the start phase before the session of
	// `shell` begins; nil in `up`.
	onStart func(driving.DevcontainerUpResponse)
	// name is the composed container name of the request
	// (`u-boot-<project>[-<instance>]`); nameHanded reports whether the
	// launcher receives it with a per-call configuration. An invalid part has
	// already ended the call when the name was composed.
	name       string
	nameHanded bool
	// overlay is the per-call configuration the launcher hands to the CLI
	// when nameHanded holds; nil otherwise.
	overlay []byte
}

// start is the start phase that `up` and `shell` share (ADR-0031 decision 1):
// validation, the disclosure of the socket, the launcher, and the judgement by
// the state of the container ([DevcontainerStartService.startOutcome]). It
// returns the response and the [startResolution] it started from.
func (s *DevcontainerStartService) start(ctx context.Context, r startRequest) (driving.DevcontainerUpResponse, startResolution, error) {
	if err := s.wired(r.command, r.baseDir); err != nil {
		return driving.DevcontainerUpResponse{}, startResolution{}, err
	}
	if driven.ValidateDevcontainerProjectFolder(r.baseDir) != nil {
		return driving.DevcontainerUpResponse{}, startResolution{}, fmt.Errorf("devcontainer %s: %w: the path must be absolute and hold no comma, double quote or line break",
			r.command, driving.ErrDevcontainerProjectFolderInvalid)
	}
	if err := validInstance(r.command, r.instance); err != nil {
		return driving.DevcontainerUpResponse{}, startResolution{}, err
	}
	name, nameHanded, err := s.containerName(r)
	if err != nil {
		return driving.DevcontainerUpResponse{}, startResolution{}, err
	}
	r.name = name
	r.nameHanded = nameHanded
	config, err := s.readConfig(r.command, r.baseDir)
	if err != nil {
		return driving.DevcontainerUpResponse{}, startResolution{}, err
	}
	names, err := referencedEnvNames(r.command, config)
	if err != nil {
		return driving.DevcontainerUpResponse{}, startResolution{}, err
	}
	r.applyNameOverlay(config)
	envValues, gitIdentity, hostGitWarnings, tokenSource := s.resolveHostGit(ctx, r, names)
	if r.onWarning != nil {
		r.onWarning(driving.WarningEntry{
			Code: socketRelaxationCode, Level: "warn", Message: socketRelaxationText(r.duration), SecurityRelaxation: true,
		})
	}
	runningBefore, err := s.containerBeforeStart(ctx, r)
	if err != nil {
		return driving.DevcontainerUpResponse{}, startResolution{}, err
	}
	s.logger.Debug("devcontainer "+r.command, "recreate", r.recreate, "instance", r.instance, "envNames", len(names), "nameHanded", r.nameHanded)
	options := r.launcherOptions(names, envValues, gitIdentity)
	res, err := s.launcher.Up(ctx, options)
	resp, err := s.startOutcome(ctx, r, runningBefore, res, err)
	if err != nil {
		return resp, startResolution{config: config}, err
	}
	resp.Warnings = append(resp.Warnings, hostGitWarnings...)
	resp.GitTokenSource = tokenSource
	if gitIdentity.Name != "" || gitIdentity.Email != "" {
		if !res.GitIdentityApplied {
			resp.Warnings = append(resp.Warnings, warningNoAbort("LH-FA-DEV-010", gitIdentityStepFailedWarning))
		}
	}
	return resp, startResolution{config: config, envValues: envValues}, err
}

// validInstance rejects an instance name that
// [driven.ValidateDevcontainerInstance] does not admit; the empty name (no
// instance) passes. command names the command in the error text.
func validInstance(command, instance string) error {
	if instance == "" {
		return nil
	}
	if err := driven.ValidateDevcontainerInstance(instance); err != nil {
		return fmt.Errorf("devcontainer %s: %w: %v", command, driving.ErrDevcontainerInstanceInvalid, err)
	}
	return nil
}

// onlyInstancesError is the error of a call without an instance for a project
// that has only containers of instances; it names the way out.
func onlyInstancesError(command string) error {
	return fmt.Errorf("devcontainer %s: %w; pass --instance <name> to act on one of them",
		command, driving.ErrDevcontainerOnlyInstances)
}

// handedName is the container name the launcher receives for the request; the
// empty string means the started container carries no name.
func handedName(r startRequest) string {
	if r.nameHanded {
		return r.name
	}
	return ""
}

// applyNameOverlay merges the name into the configuration text of the call
// ([containerNameOverlay]); a text that carries no merge hands no name, and
// the start keeps its behavior without one.
func (r *startRequest) applyNameOverlay(config string) {
	if !r.nameHanded {
		return
	}
	if r.overlay = containerNameOverlay(config, r.name); r.overlay == nil {
		r.nameHanded = false
	}
}

// launcherOptions is the launcher request of the start; a name that is not
// handed travels as the empty name without a per-call configuration.
func (r startRequest) launcherOptions(
	names []string, envValues map[string]string, gitIdentity driven.DevcontainerGitIdentity,
) driven.DevcontainerUpOptions {
	name := ""
	overrideConfig := []byte(nil)
	if r.nameHanded {
		name = r.name
		overrideConfig = r.overlay
	}
	return driven.DevcontainerUpOptions{
		ProjectFolder:  r.baseDir,
		Instance:       r.instance,
		Name:           name,
		OverrideConfig: overrideConfig,
		EnvNames:       names,
		EnvValues:      envValues,
		GitIdentity:    gitIdentity,
		Recreate:       r.recreate,
		Progress:       r.progress,
	}
}

// containerName composes the container name of the start
// (`u-boot-<project>[-<instance>]`, LH-FA-DEV-010) and hands it, when
// accepted. The project part is the project name of u-boot.yaml, or, without
// the file, the normalized basename of the folder; a part the engine would
// refuse ends the call as a validation error (exit code 10), before the
// engine is asked and before the disclosure of the socket. A name the engine
// accepts is still not handed when the configuration text cannot carry it
// ([containerNameOverlay]); nameHanded reports that case: the start keeps its
// behavior and the container carries no name.
func (s *DevcontainerStartService) containerName(r startRequest) (string, bool, error) {
	project := s.projectName(r.baseDir)
	name, err := domain.ComposeDevcontainerContainerName(project, r.instance)
	if err != nil {
		return "", false, fmt.Errorf("devcontainer %s: %w: %v", r.command, driving.ErrDevcontainerContainerNameInvalid, err)
	}
	if err := driven.ValidateDockerContainerName(name); err != nil {
		return "", false, fmt.Errorf("devcontainer %s: %w: %v", r.command, driving.ErrDevcontainerContainerNameInvalid, err)
	}
	return name, true, nil
}

// projectName is the project part of the container name: the project name of
// u-boot.yaml, or, without the file, an unparsable one or without a project
// name, the normalized basename of the folder (LH-FA-INIT-002). A missing or
// unparsable file is no failure of `up`; the name keeps its sources.
func (s *DevcontainerStartService) projectName(baseDir string) string {
	path := filepath.Join(baseDir, domain.ConfigFileName)
	if exists, err := s.fs.Exists(path); err == nil && exists {
		if body, rerr := s.fs.ReadFile(path); rerr == nil {
			var cfg ubootYAMLConfig
			if serr := s.yaml.Unmarshal(body, &cfg); serr == nil {
				return cfg.Project.Name
			}
		}
	}
	return domain.NormalizeProjectName(filepath.Base(baseDir))
}

// containerNameOverlay merges the name into the configuration text as a
// `runArgs` value (`--name=<name>`, measured for the pinned CLI 0.89.0):
// an existing list receives the value, a missing property gains the list.
// A text that u-boot cannot read merges nothing, nil, and the start starts
// without a name — the configuration of the project stays the source of truth,
// the overlay replaces it for this call only, and comments and formatting of
// the file are not part of the copy.
func containerNameOverlay(config, name string) []byte {
	var shape map[string]any
	if err := json.Unmarshal([]byte(withoutComments(config)), &shape); err != nil {
		return nil
	}
	args := []any{}
	if raw, ok := shape["runArgs"]; ok {
		list, ok := raw.([]any)
		if !ok {
			return nil
		}
		args = list
	}
	args = append(args, "--name="+name)
	shape["runArgs"] = args
	out, err := json.Marshal(shape)
	if err != nil {
		return nil
	}
	return out
}

// containerBeforeStart asks the inspector for the container of the request
// before the launcher is called, for two purposes. Without an instance it
// ends the call with [driving.ErrDevcontainerOnlyInstances] when the project
// has only containers of instances, because the devcontainer CLI takes the
// newest instance container for a call without id labels. For `up` it returns
// the ID of the container if it runs (ADR-0030 decision 2), else the empty
// string: with `--recreate` (the container is removed, so none is reported as
// running), with no container, a stopped one, and when the engine cannot be
// asked. A failed query is no failure of the command; the launcher reports an
// engine problem itself.
func (s *DevcontainerStartService) containerBeforeStart(ctx context.Context, r startRequest) (string, error) {
	wantRunning := r.reportRunning && !r.recreate
	if r.instance != "" && !wantRunning {
		return "", nil
	}
	container, found, err := s.inspector.FindContainer(ctx, r.baseDir, r.instance)
	if errors.Is(err, driven.ErrOnlyInstanceContainers) {
		return "", onlyInstancesError(r.command)
	}
	if !wantRunning || err != nil || !found || !container.Running {
		return "", nil
	}
	return container.ID, nil
}

// startOutcome judges the result of the launcher by the state of the container
// (ADR-0030 decision 1). Only a start that succeeded or failed with the runtime
// class reads the status files, through one [fallbackReader]: an engine
// problem or an abort keeps its class and asks nothing. A `result=failed` of a
// configured capability ends the call with
// [driving.ErrDevcontainerStartScriptFailed], whether the CLI failed (the
// message keeps its output) or succeeded (the file stems from an earlier start
// of the container). A conflict of the container name keeps its text and takes
// its exit class ([driving.ErrDevcontainerContainerNameConflict]). Otherwise
// the response carries the fallback warnings and whether the container ran
// before the call. The request names the command in
// the error texts and the instance whose container the status files are read
// from.
func (s *DevcontainerStartService) startOutcome(
	ctx context.Context, r startRequest, runningBefore string, res driven.DevcontainerUpResult, launchErr error,
) (driving.DevcontainerUpResponse, error) {
	command := r.command
	if errors.Is(launchErr, driven.ErrDevcontainerContainerNameConflict) {
		return driving.DevcontainerUpResponse{}, fmt.Errorf("%s: %w", launchErr.Error(), driving.ErrDevcontainerContainerNameConflict)
	}
	if launchErr != nil && !errors.Is(launchErr, driven.ErrDevcontainerRuntime) {
		return driving.DevcontainerUpResponse{}, launchErr
	}
	probes := s.fallbackProbes(r.baseDir)
	reader := &fallbackReader{inspector: s.inspector, baseDir: r.baseDir, instance: r.instance}
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
		ContainerName:         handedName(r),
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

// gitTokenEnvName is the reference of the git token in the generated
// configuration (LH-FA-DEV-009): it is also the source name the resolution
// reports for the environment variable of the host.
const gitTokenEnvName = driven.TokenSourceEnvironment

// Warning texts of the host git resolution (LH-FA-DEV-009, LH-FA-DEV-010):
// one line each, no identifier, no value, indicative over the state.
const (
	tokenSourceHostEnvWarning = "the git token of this container comes from the environment variable GIT_TOKEN of the host"
	tokenSourceGhWarning      = "the git token of this container comes from the gh session of the host (gh auth token)"
	// %s names the error of the resolution; the value never appears in a
	// warning (LH-FA-DEV-009).
	tokenResolutionFailedWarningT = "the git token of the host could not be resolved (%s); the container carries no token"
	tokenMissingSourceWarning     = "no host source carries a git token (GIT_TOKEN is not set on the host and gh does not answer); git operations in the container that need one will fail or ask"
	tokenSwitchWithoutRefWarning  = "devcontainer.sandbox.hostGitToken is switched on, but the devcontainer configuration references no ${localEnv:GIT_TOKEN}; the switch resolves nothing here"
	// %s names the key (user.name or user.email).
	gitIdentityMissingKeyWarningT  = "the git configuration of the project folder carries no %s; the container sets none for it"
	identityHostConfigReadWarningT = "the git configuration of the host could not be read (%s); the container sets no identity"
	gitIdentityStepFailedWarning   = "the git identity of the host could not be set in the container; the container keeps the identity it carries"
	gitIdentityMissingWarning      = "the git configuration of the host carries no user.name and no user.email; git commit in the container asks for an identity"
)

// warningNoAbort builds one warning entry of the host git resolution:
// Level warn, code the requirement the text belongs to.
func warningNoAbort(code, message string) driving.WarningEntry {
	return driving.WarningEntry{Code: code, Level: "warn", Message: message}
}

// resolveHostGit resolves the host values the start hands to the launcher:
// the token of the opt-in sources (reference present, sandbox profile,
// switch on) and the identity of the host git configuration (both profiles;
// D2 of the slice: the launcher sets user.name/user.email in the container
// after the start; LH-FA-DEV-009, LH-FA-DEV-010). Without a wired
// [driven.HostGitSource] port it resolves nothing: no values, no identity,
// no warning. Host git missing and missing values warn without abort; the
// value of the token appears in no warning and in no output.
func (s *DevcontainerStartService) resolveHostGit(
	ctx context.Context, r startRequest, names []string,
) (envValues map[string]string, gitIdentity driven.DevcontainerGitIdentity, warnings []driving.WarningEntry, tokenSource string) {
	envValues = map[string]string{}
	if s.hostSource == nil {
		return envValues, gitIdentity, nil, ""
	}
	sandboxProfile, switchOn := s.hostGitSwitch(r.baseDir)
	resolveToken := sandboxProfile && switchOn
	envValues, tokenSource, warnings = s.resolveTokenSource(ctx, names, resolveToken, envValues)
	gitIdentity, warnings = s.resolveIdentity(ctx, r, warnings)
	return envValues, gitIdentity, warnings, tokenSource
}

// resolveTokenSource applies the token resolution of the opt-in switch to
// the environment values (LH-FA-DEV-009); with no switch, no sandbox profile
// or no reference it changes nothing. The host variable keeps the transport
// of the present: the port answers its source first; the launcher passes it
// as before.
func (s *DevcontainerStartService) resolveTokenSource(
	ctx context.Context, names []string, resolveToken bool, envValues map[string]string,
) (map[string]string, string, []driving.WarningEntry) {
	var warnings []driving.WarningEntry
	if !resolveToken {
		return envValues, "", warnings
	}
	if !containsString(names, gitTokenEnvName) {
		return envValues, "", append(warnings, warningNoAbort("LH-FA-DEV-009", tokenSwitchWithoutRefWarning))
	}
	value, source, terr := s.hostSource.GitToken(ctx)
	switch {
	case terr != nil:
		warnings = append(warnings, warningNoAbort("LH-FA-DEV-009", fmt.Sprintf(tokenResolutionFailedWarningT, terr.Error())))
	case value == "":
		warnings = append(warnings, warningNoAbort("LH-FA-DEV-009", tokenMissingSourceWarning))
	case source == driven.TokenSourceEnvironment:
		// The host passes its variable as before (LH-FA-DEV-009).
		warnings = append(warnings, warningNoAbort("LH-FA-DEV-009", tokenSourceHostEnvWarning))
	default:
		envValues[gitTokenEnvName] = value
		warnings = append(warnings, warningNoAbort("LH-FA-DEV-009", tokenSourceGhWarning))
	}
	return envValues, source, warnings
}

// resolveIdentity reads the effective values of the project folder on the
// host (LH-FA-DEV-010); a missing key or an empty value sets nothing (the
// launcher keeps a value the container already carries); a failed reading
// and missing values warn without abort.
func (s *DevcontainerStartService) resolveIdentity(
	ctx context.Context, r startRequest, warnings []driving.WarningEntry,
) (driven.DevcontainerGitIdentity, []driving.WarningEntry) {
	var gitIdentity driven.DevcontainerGitIdentity
	name, nerr := s.hostSource.ConfigValue(ctx, r.baseDir, "user.name")
	email, eerr := s.hostSource.ConfigValue(ctx, r.baseDir, "user.email")
	switch {
	case nerr != nil || eerr != nil:
		readErr := nerr
		if readErr == nil {
			readErr = eerr
		}
		warnings = append(warnings, warningNoAbort("LH-FA-DEV-010", fmt.Sprintf(identityHostConfigReadWarningT, readErr.Error())))
	case name == "" && email == "":
		warnings = append(warnings, warningNoAbort("LH-FA-DEV-010", gitIdentityMissingWarning))
	default:
		gitIdentity.Name = name
		gitIdentity.Email = email
		if name == "" {
			warnings = append(warnings, warningNoAbort("LH-FA-DEV-010", fmt.Sprintf(gitIdentityMissingKeyWarningT, "user.name")))
		}
		if email == "" {
			warnings = append(warnings, warningNoAbort("LH-FA-DEV-010", fmt.Sprintf(gitIdentityMissingKeyWarningT, "user.email")))
		}
	}
	return gitIdentity, warnings
}

// hostGitSwitch reads the opt-in switch and the profile of the project
// (`devcontainer.sandbox.hostGitToken`, LH-FA-DEV-009): a sandbox profile
// with the switch on resolves; a missing u-boot.yaml, an unparsable one and
// a missing key resolve nothing (no switch, no resolution).
func (s *DevcontainerStartService) hostGitSwitch(baseDir string) (sandboxProfile, switchOn bool) {
	path := filepath.Join(baseDir, domain.ConfigFileName)
	exists, err := s.fs.Exists(path)
	if err != nil || !exists {
		return false, false
	}
	body, err := s.fs.ReadFile(path)
	if err != nil {
		return false, false
	}
	var cfg ubootYAMLConfig
	if err := s.yaml.Unmarshal(body, &cfg); err != nil {
		return false, false
	}
	dc := cfg.Devcontainer
	return profileIsSandbox(dc), dc != nil && dc.Sandbox != nil && dc.Sandbox.HostGitToken != nil && *dc.Sandbox.HostGitToken
}

func containsString(list []string, what string) bool {
	for _, item := range list {
		if item == what {
			return true
		}
	}
	return false
}

// fallbackProbes returns the status files of the capabilities the
// project's `u-boot.yaml` asks for (LH-FA-DEV-007). Without the file, with
// an unreadable or unparsable one, and without the sandbox profile it
// returns none: the fallback disclosure and the exit code 11 of a start
// script abort need the configuration, and their absence is no error here.
func (s *DevcontainerStartService) fallbackProbes(baseDir string) []fallbackProbe {
	path := filepath.Join(baseDir, domain.ConfigFileName)
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
	if err := validInstance("down", req.Instance); err != nil {
		return driving.DevcontainerDownResponse{}, err
	}
	container, found, err := s.findContainer(ctx, "down", req.BaseDir, req.Instance)
	if err != nil {
		return driving.DevcontainerDownResponse{}, err
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

// findContainer asks the inspector for the container of the project folder and
// the instance and gives the failures of the search the identity that `down`
// and `remove` share: the context error, the project that has only containers
// of instances, the remote engine, and any other failure of the engine as
// [driven.ErrDockerUnavailable]. command names the command in the error text.
func (s *DevcontainerStartService) findContainer(ctx context.Context, command, baseDir, instance string) (driven.DevcontainerContainer, bool, error) {
	container, found, err := s.inspector.FindContainer(ctx, baseDir, instance)
	if err == nil {
		return container, found, nil
	}
	if ctxErr := ctx.Err(); ctxErr != nil {
		return driven.DevcontainerContainer{}, false, ctxErr
	}
	if errors.Is(err, driven.ErrOnlyInstanceContainers) {
		return driven.DevcontainerContainer{}, false, onlyInstancesError(command)
	}
	if errors.Is(err, driven.ErrContainerEngineNotLocal) {
		return driven.DevcontainerContainer{}, false, fmt.Errorf("devcontainer %s: %w", command, err)
	}
	return driven.DevcontainerContainer{}, false, fmt.Errorf("devcontainer %s: the container engine could not be queried: %w (%s)",
		command, driven.ErrDockerUnavailable, err)
}
