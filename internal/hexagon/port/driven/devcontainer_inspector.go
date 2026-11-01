package driven

import (
	"context"
	"errors"
)

// ErrContainerFileNotFound signals that the requested file does not
// exist in the container ([DevcontainerInspector.ReadFile]).
var ErrContainerFileNotFound = errors.New("container file not found")

// ErrContainerFileInvalid signals that the path in the container is not
// a regular file or is larger than the limit; nothing of its content was
// returned ([DevcontainerInspector.ReadFile]).
var ErrContainerFileInvalid = errors.New("container file invalid")

// ErrContainerEngineNotLocal signals that the container engine the host
// CLI talks to is not a local socket (a remote `DOCKER_HOST` or Docker
// context). A container found there cannot be tied to this project
// folder ([DevcontainerInspector.FindContainer]).
var ErrContainerEngineNotLocal = errors.New("container engine is not local")

// ErrOnlyInstanceContainers signals that the project folder has containers,
// all of them of instances, and none without an instance
// ([DevcontainerInspector.FindContainer] with an empty instance). It is not a
// failure of the engine: "there is no container" is the false bool without an
// error.
var ErrOnlyInstanceContainers = errors.New("the project has only containers of instances")

// DevcontainerContainer identifies the container a devcontainer tool
// created for a project folder.
type DevcontainerContainer struct {
	// ID is the full container ID (lower-case hex).
	ID string
	// Running reports whether the container is running; a stopped
	// container is still readable.
	Running bool
}

// DevcontainerInspector is the read-only view of the devcontainer
// container of a project (LH-FA-DEV-007: the host determines the
// fallbacks of the container start scripts from the status files they
// write). The two methods are the only engine access it needs; it never
// runs a program in the container and never follows a symlink in it.
//
// Layer rules (LH-FA-ARCH-002, LH-FA-ARCH-003): driven port; the
// production implementation lives in `internal/adapter/driven/docker`
// next to [DockerProbe] and [DockerEngine].
type DevcontainerInspector interface {
	// FindContainer returns the container that carries the label
	// `devcontainer.local_folder=<localFolder>` (absolute project
	// folder); with several, the most recently created. With a non-empty
	// instance the container must also carry `u-boot.instance=<instance>`; an
	// instance that [ValidateDevcontainerInstance] rejects is an error. With
	// an empty instance the container must carry no `u-boot.instance` label:
	// when the folder has containers and all of them carry it, the call fails
	// with [ErrOnlyInstanceContainers]. The bool is false when there is none.
	// Any other non-nil error means the engine could not be asked
	// ([ErrContainerEngineNotLocal] for a remote engine); callers treat it as
	// "not determinable", never as a finding.
	FindContainer(ctx context.Context, localFolder, instance string) (DevcontainerContainer, bool, error)

	// ReadFile returns the content of the regular file `path` (absolute)
	// in the container, at most maxBytes bytes. It fails with
	// [ErrContainerFileNotFound] when the file is absent and with
	// [ErrContainerFileInvalid] for any other entry type (symlink, FIFO,
	// directory, device) or a file larger than maxBytes.
	ReadFile(ctx context.Context, containerID, path string, maxBytes int) ([]byte, error)
}
