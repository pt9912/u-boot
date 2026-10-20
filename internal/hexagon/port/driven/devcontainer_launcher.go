package driven

import (
	"context"
	"errors"
	"io"
)

// ErrDevcontainerRuntime signals that starting or stopping a devcontainer
// failed after the environment checks passed: the build of the helper
// image, the start of the helper container, a CLI outcome other than
// `success`, or `docker stop` ([DevcontainerLauncher]). The CLI maps it to
// exit code 12 (LH-FA-CLI-006); a missing Docker binary, an unreachable
// daemon and a remote engine stay [ErrDockerUnavailable] and
// [ErrContainerEngineNotLocal] (exit code 11).
var ErrDevcontainerRuntime = errors.New("devcontainer runtime failed")

// DevcontainerUpOptions parameterises [DevcontainerLauncher.Up].
type DevcontainerUpOptions struct {
	// ProjectFolder is the absolute project folder. The launcher mounts it
	// into the helper container at the same path and hands that path to
	// the devcontainer CLI, so the label `devcontainer.local_folder` of the
	// started container carries the path of the host.
	ProjectFolder string
	// EnvNames lists the host environment variables handed to the helper
	// container by name; the value never appears on a command line.
	// Names that are not set on the host are skipped.
	EnvNames []string
	// Recreate removes an existing container before the start
	// (`--remove-existing-container`); volumes stay.
	Recreate bool
	// ReadOnlyProject mounts the project folder read-only into the helper
	// container. A lifecycle command that writes into the project then
	// fails.
	ReadOnlyProject bool
	// Progress receives the build output of the helper image and the
	// stderr of the devcontainer CLI; nil discards it.
	Progress io.Writer
}

// DevcontainerUpResult is the result line of a successful
// `devcontainer up`.
type DevcontainerUpResult struct {
	// ContainerID is the ID of the started devcontainer (lower-case hex).
	ContainerID string
	// RemoteUser is the user the CLI reports for the container.
	RemoteUser string
	// RemoteWorkspaceFolder is the workspace path inside the container.
	RemoteWorkspaceFolder string
}

// DevcontainerLauncher starts and stops the devcontainer of a project
// without an editor and without Node on the host: `Up` runs the Dev
// Container CLI in a short-lived helper container that holds the engine
// socket; `Stop` stops the container with `docker stop`.
//
// Both methods check the environment first (Docker binary on the `PATH`,
// daemon answers, engine is local) and fail with [ErrDockerUnavailable] or
// [ErrContainerEngineNotLocal]; every later failure is
// [ErrDevcontainerRuntime]. A cancelled context ends `Up` with the context
// error after the helper container has been removed.
//
// Layer rules (LH-FA-ARCH-002, LH-FA-ARCH-003): driven port; the
// production implementation lives in `internal/adapter/driven/docker`.
type DevcontainerLauncher interface {
	// Up starts the devcontainer of opts.ProjectFolder and returns the
	// result line of the CLI.
	Up(ctx context.Context, opts DevcontainerUpOptions) (DevcontainerUpResult, error)

	// Stop stops the container with the given ID (full ID or prefix of at
	// least 12 hex characters). It neither removes the container nor
	// starts a helper container.
	Stop(ctx context.Context, containerID string) error
}
