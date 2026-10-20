package application

import (
	"context"
	"errors"
	"fmt"
	iofs "io/fs"
	"path/filepath"
	"regexp"
	"sort"

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

// socketRelaxationMessage is the disclosure text. Each of its parts has a
// test: the helper holds the engine socket (`TestRunArgs`), the generated
// configuration does not mount it (`TestDevcontainerUp_Docker_S1_S9_StandardProfile`, `TestDevcontainerUp_Docker_S1_S4_SandboxProfileVolumes`) and a
// hand-written one can (`TestLauncher_Docker_S6_GuardHandWrittenSocketMountIsDetected`),
// and an `initializeCommand` runs with the socket
// (`TestLauncher_Docker_S6_InitializeCommandRunsWithSocket`).
const socketRelaxationMessage = "this command mounts the host container-engine socket into a short-lived helper container " +
	"(root-equivalent access to the engine for the duration of the call); " +
	"the devcontainer configuration that u-boot generates does not mount it, " +
	"but a hand-written " + devcontainerConfigRel + " can mount it, " +
	"and an initializeCommand declared there runs with this access"

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

// DevcontainerStartService implements [driving.DevcontainerUseCase]
// (`u-boot devcontainer up` and `down`, LH-FA-DEV-010). It reads the
// project files through the [driven.FileSystem] port, starts and stops
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
// `Up` and `Down` fail instead of reporting "nothing to do". A nil logger
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
	if err := s.wired("up", req.BaseDir); err != nil {
		return driving.DevcontainerUpResponse{}, err
	}
	if driven.ValidateDevcontainerProjectFolder(req.BaseDir) != nil {
		return driving.DevcontainerUpResponse{}, fmt.Errorf("devcontainer up: %w: the path must be absolute and hold no comma, double quote or line break",
			driving.ErrDevcontainerProjectFolderInvalid)
	}
	names, err := s.referencedEnvNames(req.BaseDir)
	if err != nil {
		return driving.DevcontainerUpResponse{}, err
	}
	if req.OnWarning != nil {
		req.OnWarning(driving.WarningEntry{
			Code: socketRelaxationCode, Level: "warn", Message: socketRelaxationMessage, SecurityRelaxation: true,
		})
	}
	s.logger.Debug("devcontainer up", "recreate", req.Recreate, "envNames", len(names))
	res, err := s.launcher.Up(ctx, driven.DevcontainerUpOptions{
		ProjectFolder: req.BaseDir,
		EnvNames:      names,
		Recreate:      req.Recreate,
		Progress:      req.ProgressSink,
	})
	probes := s.fallbackProbes(req.BaseDir)
	reader := &fallbackReader{inspector: s.inspector, baseDir: req.BaseDir}
	if err != nil {
		if errors.Is(err, driven.ErrDevcontainerRuntime) && startScriptFailed(ctx, reader, probes) {
			return driving.DevcontainerUpResponse{}, fmt.Errorf("devcontainer up: %w: %v", driving.ErrDevcontainerStartScriptFailed, err)
		}
		return driving.DevcontainerUpResponse{}, err
	}
	return driving.DevcontainerUpResponse{
		ContainerID:           res.ContainerID,
		RemoteUser:            res.RemoteUser,
		RemoteWorkspaceFolder: res.RemoteWorkspaceFolder,
		Warnings:              fallbackWarnings(ctx, reader, probes, startFallbackWording()),
	}, nil
}

// referencedEnvNames reads `.devcontainer/devcontainer.json` and returns the
// variable names it references through `${localEnv:<NAME>}`. A missing file
// is [driving.ErrDevcontainerConfigMissing], a read failure
// [driving.ErrDevcontainerFileSystem], a name the launcher could not pass
// [driving.ErrDevcontainerEnvNameInvalid].
func (s *DevcontainerStartService) referencedEnvNames(baseDir string) ([]string, error) {
	path := filepath.Join(baseDir, devcontainerConfigRel)
	exists, err := s.fs.Exists(path)
	if err != nil {
		return nil, fmt.Errorf("devcontainer up: probe %q: %w: %w", path, driving.ErrDevcontainerFileSystem, err)
	}
	if !exists {
		return nil, fmt.Errorf("devcontainer up: %w (looked for %s)", driving.ErrDevcontainerConfigMissing, devcontainerConfigRel)
	}
	body, err := s.fs.ReadFile(path)
	if err != nil {
		if errors.Is(err, iofs.ErrNotExist) {
			return nil, fmt.Errorf("devcontainer up: %w (looked for %s)", driving.ErrDevcontainerConfigMissing, devcontainerConfigRel)
		}
		return nil, fmt.Errorf("devcontainer up: read %q: %w: %w", path, driving.ErrDevcontainerFileSystem, err)
	}
	names := localEnvNames(string(body))
	for _, name := range names {
		if driven.ValidateDevcontainerEnvName(name) != nil {
			return nil, fmt.Errorf("devcontainer up: %w: %q in %s", driving.ErrDevcontainerEnvNameInvalid, name, devcontainerConfigRel)
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

// startScriptFailed reports whether one of the status files reports
// `result=failed` (ADR-0029 decision 1). The file describes the last run of
// its script, not necessarily this start.
func startScriptFailed(ctx context.Context, r *fallbackReader, probes []fallbackProbe) bool {
	for _, p := range probes {
		if rd := r.read(ctx, p); rd.outcome == statusReported && rd.status.Result() == domain.SandboxStatusFailed {
			return true
		}
	}
	return false
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
