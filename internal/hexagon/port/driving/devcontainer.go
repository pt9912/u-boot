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
	// AlreadyRunning is true when the container of the project ran before the
	// call and the devcontainer CLI reports the same ID (ADR-0030 decision 2).
	// It is false with [DevcontainerUpRequest.Recreate], for a container that
	// was stopped or absent before the call, for another ID, and when the
	// engine could not be asked before the call.
	AlreadyRunning bool
}

// DevcontainerShellRequest is the input for [DevcontainerUseCase.Shell]
// (`u-boot devcontainer shell`, LH-FA-DEV-010).
type DevcontainerShellRequest struct {
	// BaseDir is the absolute project folder; the CLI adapter passes the
	// working directory. Mandatory.
	BaseDir string

	// ProgressSink receives the build output of the helper image and the
	// output of the devcontainer CLI during the start; nil discards it.
	ProgressSink io.Writer

	// OnWarning receives, synchronously, first the disclosure of the socket
	// relaxation before the helper container is started (it is not called
	// when the request fails before that point), then each degradation
	// fallback the container reports after the start, before the session
	// begins. nil discards them.
	OnWarning func(WarningEntry)

	// Stdin, Stdout and Stderr are the streams of the session. The use case
	// writes nothing of its own to them. Stdout and Stderr must not be one
	// shared writer: the session writes to both from separate goroutines. A
	// nil stream is the null device.
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	// Terminal allocates a terminal for the session. The adapter sets it when
	// its standard input is a terminal; Stdin and Stdout are then the files
	// of that terminal.
	Terminal bool
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

// ErrDevcontainerStartScriptFailed signals that a start script of the
// container reported `result=failed` in its status file: in the run of the
// devcontainer CLI that failed, or in an earlier start of the same container
// while the CLI succeeded (ADR-0029 decision 1, ADR-0030 decision 1). Maps to
// exit code 11, like the same cause in `u-boot doctor`. The error does not
// wrap `driven.ErrDevcontainerRuntime` (exit code 12), so the two classes do
// not overlap.
var ErrDevcontainerStartScriptFailed = errors.New("devcontainer: a container start script reported a failure")

// DevcontainerUseCase is the driving port of `u-boot devcontainer up`,
// `u-boot devcontainer down` and `u-boot devcontainer shell` (LH-FA-DEV-010).
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
//     `driven.ErrContainerEngineNotLocal`, `driven.ErrDevcontainerRuntime`).
//   - After the devcontainer CLI ran, Up reads the status files of the
//     capabilities `u-boot.yaml` configures. A file with `result=failed` ends
//     the call with [ErrDevcontainerStartScriptFailed]: after a failed start
//     instead of the runtime sentinel, after a successful start instead of the
//     response, also when the file stems from an earlier start of the same
//     container.
//   - Up reports in [DevcontainerUpResponse.AlreadyRunning] whether the
//     container ran before the call; a failed query is no failure of Up.
//   - Down stops a running container and succeeds when there is none or it
//     is stopped. It never removes a container or a volume.
//   - Shell runs the start of Up without Recreate: the same validation errors,
//     the same sentinel identities, and the same judgement by the state of the
//     container, so a status file with `result=failed` ends the call with
//     [ErrDevcontainerStartScriptFailed] and opens no session. Its disclosure
//     of the socket names the duration of the start. After the start it opens
//     one interactive session in the container as the remote user in the
//     workspace folder the start reported, with the host variables that
//     `"<NAME>": "${localEnv:<NAME>}"` entries of the configuration name and
//     that are set. It returns when the session ends: nil whatever the exit
//     status of the shell was, an error when the session does not begin
//     (`driven.ErrDevcontainerRuntime` when the probe of the session fails).
type DevcontainerUseCase interface {
	Up(ctx context.Context, req DevcontainerUpRequest) (DevcontainerUpResponse, error)
	Down(ctx context.Context, req DevcontainerDownRequest) (DevcontainerDownResponse, error)
	Shell(ctx context.Context, req DevcontainerShellRequest) error
}
