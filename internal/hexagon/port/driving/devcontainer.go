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

	// Instance mirrors `--instance`: the name of an instance of the project
	// folder. The empty string selects the container without an instance.
	// The CLI adapter checks a flag that was given empty itself
	// ([ErrDevcontainerInstanceInvalid]), so "" here always means none.
	Instance string

	// ProgressSink receives the build output of the helper image and the
	// output of the devcontainer CLI; nil discards it.
	ProgressSink io.Writer

	// OnWarning receives the disclosure of the socket relaxation
	// (LH-FA-DEV-010) synchronously, before the helper container is
	// started, so the adapter can show it ahead of a long start. It is not
	// called when the request fails before that point (no devcontainer
	// configuration, a folder or variable name the helper cannot carry, an
	// instance name that is not admissible). A call without an instance for a
	// project with only containers of instances ends after it, without a
	// helper container. nil
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
	// the start (LH-FA-DEV-007) and the findings of the host git resolution
	// (LH-FA-DEV-009, LH-FA-DEV-010: the missing token source, the missing
	// identity, a failed identity step); it does not hold the disclosure
	// that [DevcontainerUpRequest.OnWarning] receives.
	Warnings []WarningEntry
	// GitTokenSource names the source of the git token this resolution
	// handed to the container (LH-FA-DEV-009): "GIT_TOKEN" for the
	// environment variable of the host, "gh auth token" for the gh session.
	// It is empty when no token flows — no reference in the configuration,
	// the switch off, or no source found (the response carries a warning
	// about the latter). The caller prints the name of the source, never a
	// value (LH-FA-DEV-009).
	GitTokenSource string
	// AlreadyRunning is true when the container of the project ran before the
	// call and the devcontainer CLI reports the same ID (ADR-0030 decision 2).
	// It is false with [DevcontainerUpRequest.Recreate], for a container that
	// was stopped or absent before the call, for another ID, and when the
	// engine could not be asked before the call.
	AlreadyRunning bool
	// ContainerName is the name the launcher handed to the CLI for the
	// started container (`u-boot-<project>` or, with an instance,
	// `u-boot-<project>-<instance>`); it is empty when no name was handed.
	// The caller prints it beside the container ID; the --json output of
	// `up` carries no field for it.
	ContainerName string
}

// DevcontainerShellRequest is the input for [DevcontainerUseCase.Shell]
// (`u-boot devcontainer shell`, LH-FA-DEV-010).
type DevcontainerShellRequest struct {
	// BaseDir is the absolute project folder; the CLI adapter passes the
	// working directory. Mandatory.
	BaseDir string

	// Instance mirrors `--instance`, see [DevcontainerUpRequest.Instance].
	Instance string

	// ProgressSink receives the build output of the helper image and the
	// output of the devcontainer CLI during the start; nil discards it.
	ProgressSink io.Writer

	// OnWarning receives, synchronously, first the disclosure of the socket
	// relaxation before the helper container is started (it is not called
	// when the request fails before that point), then each degradation
	// fallback the container reports after the start, before the session
	// begins. nil discards them.
	OnWarning func(WarningEntry)

	// OnStart receives the response of the start phase once the container is
	// up and before the session begins, after the fallback warnings — the
	// caller prints the container of the session beside its ID there (the
	// output reaches the standard error stream; the session then owns the
	// standard output). nil discards it.
	OnStart func(DevcontainerUpResponse)

	// Stdin, Stdout and Stderr are the streams of the session. The use case
	// writes nothing of its own to them. Stdout and Stderr must not be one
	// shared writer: the session writes to both from separate goroutines. A
	// nil stream is the null device.
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	// Terminal allocates a terminal for the session. The adapter sets it when
	// its standard input is a terminal; Stdin is then the file of that
	// terminal. Stdout is passed as given and is not checked: it may be no
	// terminal (a pipe), a combination that no test covers.
	Terminal bool
}

// DevcontainerDownRequest is the input for [DevcontainerUseCase.Down]
// (`u-boot devcontainer down`).
type DevcontainerDownRequest struct {
	// BaseDir is the absolute project folder. Mandatory.
	BaseDir string

	// Instance mirrors `--instance`, see [DevcontainerUpRequest.Instance].
	Instance string
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

// DevcontainerRemoveRequest is the input for [DevcontainerUseCase.Remove]
// (`u-boot devcontainer remove`).
type DevcontainerRemoveRequest struct {
	// BaseDir is the absolute project folder. Mandatory.
	BaseDir string

	// Instance mirrors `--instance`, see [DevcontainerUpRequest.Instance].
	Instance string

	// Volumes mirrors `--volumes`: the volumes the container mounts are removed
	// together with it. It is the only destructive part of the command, so
	// the confirmation gate below applies to it alone.
	Volumes bool

	// AssumeYes mirrors `--yes`: with Volumes the confirmation is given.
	AssumeYes bool

	// NonInteractive mirrors `--no-interactive`: with Volumes and without
	// AssumeYes the call ends with [ErrConfirmationRequired] and asks nothing.
	NonInteractive bool

	// SilenceConfirmer is set for `--json`: with Volumes and without AssumeYes
	// the question is not asked and the answer counts as declined.
	SilenceConfirmer bool
}

// DevcontainerRemoveResponse is the output of [DevcontainerUseCase.Remove] on
// success.
type DevcontainerRemoveResponse struct {
	// ContainerID is the full ID of the container found for the project; empty
	// when there is none.
	ContainerID string
	// Stopped is true when this call stopped the container before it removed it.
	Stopped bool
	// Removed is true when this call removed the container.
	Removed bool
	// RemovedVolumeNames lists the removed volumes of the container, sorted. It
	// is never nil and empty without Volumes, without a container and for a
	// container that mounts no volume.
	RemovedVolumeNames []string
}

// ErrDevcontainerConfigMissing signals that the project folder has no
// `.devcontainer/devcontainer.json`. Maps to exit code 10 (LH-FA-DEV-010,
// LH-FA-CLI-006).
var ErrDevcontainerConfigMissing = errors.New("devcontainer: .devcontainer/devcontainer.json not found")

// ErrDevcontainerProjectFolderInvalid signals a project folder whose path
// holds a character a `--mount` value cannot carry (comma, double quote,
// line break) or that is not absolute. Maps to exit code 10.
var ErrDevcontainerProjectFolderInvalid = errors.New("devcontainer: project folder path is not usable")

// ErrDevcontainerInstanceInvalid signals an instance name that
// `driven.ValidateDevcontainerInstance` rejects, or a `--instance` flag given
// empty. Maps to exit code 10; it is returned before the engine is asked.
var ErrDevcontainerInstanceInvalid = errors.New("devcontainer: instance name is not usable")

// ErrDevcontainerOnlyInstances signals a call without an instance for a
// project whose containers are all of instances. Maps to exit code 10; it is
// returned before the devcontainer CLI or `docker stop` runs, and its text
// points at `--instance`.
var ErrDevcontainerOnlyInstances = errors.New("devcontainer: the project has only containers of instances")

// ErrDevcontainerEnvNameInvalid signals a `${localEnv:<NAME>}` reference in
// `devcontainer.json` whose name is not a plain variable name. Maps to exit
// code 10.
var ErrDevcontainerEnvNameInvalid = errors.New("devcontainer: localEnv reference with an invalid variable name")

// ErrDevcontainerContainerNameInvalid signals a composed container name
// (`u-boot-<project>[-<instance>]`) that `driven.ValidateDockerContainerName`
// rejects — a part the engine would refuse for `--name`. Maps to exit code 10;
// it is returned before the engine is asked.
var ErrDevcontainerContainerNameInvalid = errors.New("devcontainer: container name is not usable")

// ErrDevcontainerContainerNameConflict signals that a container of another
// project folder already holds the container name of the start (a running and
// a stopped holder refuse the name alike; the engine's error words the
// conflict). Maps to exit code 10; it is returned before the helper container
// starts, and the holder is neither removed nor renamed.
var ErrDevcontainerContainerNameConflict = errors.New("devcontainer: container name is held by a container of another project folder")

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
// `u-boot devcontainer down`, `u-boot devcontainer shell` and
// `u-boot devcontainer remove` (LH-FA-DEV-010).
//
// Contract:
//
//   - Up needs `.devcontainer/devcontainer.json` below BaseDir
//     ([ErrDevcontainerConfigMissing]) and a folder path the helper container
//     can mount ([ErrDevcontainerProjectFolderInvalid]); it names no
//     variable but those the configuration references through
//     `${localEnv:<NAME>}` ([ErrDevcontainerEnvNameInvalid]).
//   - Up, Down and Shell act on the container of the instance the request
//     names; without an instance only on the container without one. An
//     instance name that is not admissible ends the call with
//     [ErrDevcontainerInstanceInvalid], a call without an instance for a
//     project that has only containers of instances with
//     [ErrDevcontainerOnlyInstances]; both are returned before the launcher
//     runs.
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
//   - Remove acts on the container of the instance the request names, like
//     Down, and succeeds without asking when there is none. It reads no
//     project file and starts no helper container. With Volumes it reads the
//     volumes the container mounts, then asks for the confirmation
//     ([ErrConfirmationRequired] or [ErrConfirmerUnavailable] without it,
//     before anything is stopped or removed), then stops the container if it
//     runs, removes it and removes each of its volumes. A step that fails ends
//     the call with the sentinel identity of the engine
//     (`driven.ErrDevcontainerRuntime`), and the error names what was removed
//     and what remains.
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
	Remove(ctx context.Context, req DevcontainerRemoveRequest) (DevcontainerRemoveResponse, error)
}
