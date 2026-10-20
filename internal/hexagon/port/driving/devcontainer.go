package driving

import (
	"context"
	"errors"
	"io"
)

// DevcontainerUpRequest is the input for [DevcontainerUseCase.Up]
// (`u-boot devcontainer up`, LH-FA-DEV-010).
type DevcontainerUpRequest struct {
	// BaseDir is the absolute project folder; the CLI adapter passes the
	// working directory. Mandatory.
	BaseDir string

	// Recreate mirrors `--recreate`: an existing container is removed and
	// created anew; named volumes stay.
	Recreate bool

	// ProgressSink receives the build output of the helper image and the
	// output of the devcontainer CLI; nil discards it.
	ProgressSink io.Writer

	// OnWarning receives the disclosure of the socket relaxation
	// (LH-FA-DEV-010) synchronously, before the helper container is
	// started, so the adapter can show it ahead of a long start. It is not
	// called when the request fails before that point (no devcontainer
	// configuration, a folder or variable name the helper cannot carry). nil
	// discards it.
	OnWarning func(WarningEntry)
}

// DevcontainerUpResponse is the output of [DevcontainerUseCase.Up] on
// success.
type DevcontainerUpResponse struct {
	// ContainerID is the full ID of the started container.
	ContainerID string
	// RemoteUser is the user the devcontainer CLI reports for the container.
	RemoteUser string
	// RemoteWorkspaceFolder is the workspace path inside the container.
	RemoteWorkspaceFolder string
	// Warnings holds the degradation fallbacks the container reports after
	// the start (LH-FA-DEV-007); it does not hold the disclosure that
	// [DevcontainerUpRequest.OnWarning] receives.
	Warnings []WarningEntry
}

// DevcontainerDownRequest is the input for [DevcontainerUseCase.Down]
// (`u-boot devcontainer down`).
type DevcontainerDownRequest struct {
	// BaseDir is the absolute project folder. Mandatory.
	BaseDir string
}

// DevcontainerDownResponse is the output of [DevcontainerUseCase.Down] on
// success.
type DevcontainerDownResponse struct {
	// ContainerID is the full ID of the container found for the project;
	// empty when there is none.
	ContainerID string
	// Stopped is true when this call stopped a running container. It is
	// false when there is no container and when the container was already
	// stopped; both end with exit code 0.
	Stopped bool
}

// ErrDevcontainerConfigMissing signals that the project folder has no
// `.devcontainer/devcontainer.json`. Maps to exit code 10 (LH-FA-DEV-010,
// LH-FA-CLI-006).
var ErrDevcontainerConfigMissing = errors.New("devcontainer: .devcontainer/devcontainer.json not found")

// ErrDevcontainerProjectFolderInvalid signals a project folder whose path
// holds a character a `--mount` value cannot carry (comma, double quote,
// line break) or that is not absolute. Maps to exit code 10.
var ErrDevcontainerProjectFolderInvalid = errors.New("devcontainer: project folder path is not usable")

// ErrDevcontainerEnvNameInvalid signals a `${localEnv:<NAME>}` reference in
// `devcontainer.json` whose name is not a plain variable name. Maps to exit
// code 10.
var ErrDevcontainerEnvNameInvalid = errors.New("devcontainer: localEnv reference with an invalid variable name")

// ErrDevcontainerFileSystem wraps a read failure of the devcontainer
// configuration. Maps to exit code 14 (LH-NFA-REL-003).
var ErrDevcontainerFileSystem = errors.New("devcontainer: filesystem read failed")

// ErrDevcontainerStartScriptFailed signals that the devcontainer CLI did
// not succeed and a start script of the container reported `result=failed`
// in its status file (ADR-0029 decision 1). Maps to exit code 11, like the
// same cause in `u-boot doctor`. The error does not wrap
// `driven.ErrDevcontainerRuntime` (exit code 12), so the two classes do not
// overlap.
var ErrDevcontainerStartScriptFailed = errors.New("devcontainer: a container start script reported a failure")

// DevcontainerUseCase is the driving port of `u-boot devcontainer up` and
// `u-boot devcontainer down` (LH-FA-DEV-010).
//
// Contract:
//
//   - Up needs `.devcontainer/devcontainer.json` below BaseDir
//     ([ErrDevcontainerConfigMissing]) and a folder path the helper container
//     can mount ([ErrDevcontainerProjectFolderInvalid]); it names no
//     variable but those the configuration references through
//     `${localEnv:<NAME>}` ([ErrDevcontainerEnvNameInvalid]).
//   - Up hands every failure of the engine through with its sentinel
//     identity (`driven.ErrDockerUnavailable`,
//     `driven.ErrContainerEngineNotLocal`, `driven.ErrDevcontainerRuntime`);
//     a failed start whose container reports `result=failed` is
//     [ErrDevcontainerStartScriptFailed] instead of the runtime sentinel.
//   - Down stops a running container and succeeds when there is none or it
//     is stopped. It never removes a container or a volume.
type DevcontainerUseCase interface {
	Up(ctx context.Context, req DevcontainerUpRequest) (DevcontainerUpResponse, error)
	Down(ctx context.Context, req DevcontainerDownRequest) (DevcontainerDownResponse, error)
}
