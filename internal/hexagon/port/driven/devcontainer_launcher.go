package driven

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
)

// ErrDevcontainerRuntime signals that starting, stopping or entering a
// devcontainer failed: the build of the helper image, the start of the helper
// container, a CLI outcome other than `success`, `docker stop`, or the probe
// or the start of a session ([DevcontainerLauncher]); also a request the
// launcher rejects before any docker call: the launcher applies
// [ValidateDevcontainerProjectFolder], [ValidateDevcontainerEnvName],
// [ValidateDevcontainerShellUser] and [ValidateDevcontainerWorkspace]
// itself, so a caller that did not run them first gets this error. The CLI
// maps it to exit code 12 (LH-FA-CLI-006); a
// missing Docker binary, an unreachable daemon and a remote engine stay
// [ErrDockerUnavailable] and [ErrContainerEngineNotLocal] (exit code 11).
var ErrDevcontainerRuntime = errors.New("devcontainer runtime failed")

// ValidateDevcontainerProjectFolder rejects a folder that is not an absolute
// path or that holds a character a `--mount` value cannot carry (comma,
// double quote, line break) or a NUL. It is the one definition of the
// admissible character set; the launcher applies it again.
func ValidateDevcontainerProjectFolder(folder string) error {
	if !strings.HasPrefix(folder, "/") || strings.ContainsAny(folder, ",\"\n\r\x00") {
		return fmt.Errorf("project folder %q is not an absolute path a --mount value can carry", folder)
	}
	return nil
}

// ValidateDevcontainerEnvName rejects a name that is not a plain variable
// name (`[A-Za-z_][A-Za-z0-9_]*`); the launcher hands the name to
// `docker run -e` and must not let it carry a value.
func ValidateDevcontainerEnvName(name string) error {
	for i, r := range name {
		letter := r == '_' || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
		if letter || (i > 0 && r >= '0' && r <= '9') {
			continue
		}
		return fmt.Errorf("environment variable name %q is not a plain name", name)
	}
	if name == "" {
		return fmt.Errorf("environment variable name %q is not a plain name", name)
	}
	return nil
}

// ValidateDevcontainerShellUser rejects a user the launcher cannot hand to
// `docker exec -u` as one argument: an empty one, one that begins with a
// dash, and one that holds a line break or a NUL. A name, a numeric ID and
// `name:group` pass; whether the user exists in the container is the
// business of the probe of [DevcontainerLauncher.Shell]. The launcher applies
// it again.
func ValidateDevcontainerShellUser(user string) error {
	if user == "" || strings.HasPrefix(user, "-") || strings.ContainsAny(user, "\n\r\x00") {
		return fmt.Errorf("user %q is not a user `docker exec -u` can take as one argument", user)
	}
	return nil
}

// ValidateDevcontainerWorkspace rejects a working directory that is not an
// absolute path or that holds a line break or a NUL. Whether the directory
// exists in the container is the business of the probe of
// [DevcontainerLauncher.Shell]. The launcher applies it again.
func ValidateDevcontainerWorkspace(workspace string) error {
	if !strings.HasPrefix(workspace, "/") || strings.ContainsAny(workspace, "\n\r\x00") {
		return fmt.Errorf("workspace %q is not an absolute path `docker exec -w` can take as one argument", workspace)
	}
	return nil
}

// DevcontainerShellOptions parameterises [DevcontainerLauncher.Shell].
type DevcontainerShellOptions struct {
	// ContainerID is the ID of the running devcontainer (lower-case hex, the
	// full ID or a prefix of at least 12 characters), as the result of
	// [DevcontainerLauncher.Up] reports it.
	ContainerID string
	// User is the user of the session (`docker exec -u`): the RemoteUser of
	// the result of [DevcontainerLauncher.Up].
	User string
	// Workspace is the working directory of the session (`docker exec -w`):
	// the RemoteWorkspaceFolder of the result of [DevcontainerLauncher.Up].
	Workspace string
	// EnvNames lists the host environment variables handed to the session by
	// name (`docker exec -e NAME`); the value stays out of every command
	// line, the `docker` client takes it from the environment of the
	// process. Names that are not set on the host are skipped, so a variable
	// the host lacks is absent in the session. The launcher masks nothing in
	// the session: what the shell prints reaches the caller as it is.
	EnvNames []string
	// Stdin, Stdout and Stderr are the streams of the session. The launcher
	// writes nothing of its own to Stdout. Without Terminal, Stdout receives
	// the standard output of the shell, and Stderr its standard error and the
	// messages of the `docker` client. With Terminal the shell writes through
	// the terminal, so Stdout receives everything the shell prints and
	// Stderr only the messages of the `docker` client. A nil stream is the
	// null device. With Terminal, Stdin and Stdout are handed to the
	// `docker` client as they are when they are files, so the client finds
	// the terminal on them; Stdin must then be a terminal.
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	// Terminal allocates a terminal for the session (`-t`). The caller sets
	// it when its standard input is a terminal.
	Terminal bool
}

// DevcontainerUpOptions parameterises [DevcontainerLauncher.Up].
type DevcontainerUpOptions struct {
	// ProjectFolder is the absolute project folder. The launcher mounts it
	// into the helper container at the same path and hands that path to
	// the devcontainer CLI, so the label `devcontainer.local_folder` of the
	// started container carries the path of the host.
	ProjectFolder string
	// EnvNames lists the host environment variables handed to the helper
	// container by name; the launcher's own `docker run` carries no value.
	// The devcontainer CLI itself writes the value of a variable that
	// devcontainer.json references through `${localEnv:NAME}` into the
	// command lines of `docker` processes inside the helper container
	// (`docker run -e`, `docker exec -e`) for the duration of the command;
	// with `containerEnv` the resolved value also stays in the metadata
	// (`Config.Env`) of the started devcontainer for as long as the
	// container exists. The launcher masks neither. Names that are not
	// set on the host are skipped.
	EnvNames []string
	// Recreate removes an existing container before the start
	// (`--remove-existing-container`); volumes stay.
	Recreate bool
	// ReadOnlyProject mounts the project folder read-only into the helper
	// container. A lifecycle command that writes into the project then
	// fails.
	ReadOnlyProject bool
	// Progress receives the build output of the helper image and the
	// stderr of the devcontainer CLI; nil discards it. While EnvNames
	// names a set variable with a maskable value, the launcher holds back
	// the last bytes of the stderr stream (up to the length of the longest
	// value minus one) until the next write or the end of the call, so the
	// last line of a step can appear late. The value of a
	// variable in EnvNames is replaced by a placeholder in what the
	// launcher forwards here and in its error texts, where the value
	// appears byte for byte and is at least four bytes long; a value
	// that the CLI prints in another form (escaped, encoded, split by its
	// own line breaks) or a shorter value is not masked.
	Progress io.Writer
}

// DevcontainerUpResult is the result line of a successful
// `devcontainer up`. The launcher reads the line from the unmasked stdout of
// the CLI and returns its fields as data for the caller, not as an output
// stream: a passed value that equals a field (a user name, a part of the
// container ID) leaves the field unchanged, and a caller that prints a field
// decides what appears.
type DevcontainerUpResult struct {
	// ContainerID is the ID of the started devcontainer (lower-case hex).
	ContainerID string
	// RemoteUser is the user the CLI reports for the container.
	RemoteUser string
	// RemoteWorkspaceFolder is the workspace path inside the container.
	RemoteWorkspaceFolder string
}

// DevcontainerLauncher starts, stops and enters the devcontainer of a
// project without an editor and without Node on the host: `Up` runs the Dev
// Container CLI in a short-lived helper container that holds the engine
// socket; `Stop` stops the container with `docker stop`; `Shell` opens a
// session in the running container with `docker exec` of the host.
//
// Each method first rejects an invalid argument as [ErrDevcontainerRuntime]
// without a docker call (`Up`: project folder, then variable names; `Stop`:
// container ID; `Shell`: container ID, user, workspace, then variable
// names), then checks the environment (Docker binary on the `PATH`,
// daemon answers, engine is local) and fails with [ErrDockerUnavailable] or
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

	// Shell opens a session in the running container and returns when the
	// session ends. It starts no helper container and binds no engine
	// socket, mount or privilege into the session. A probe with the user
	// and workspace of opts runs first (`bash -c true`, else `sh -c true`,
	// each non-interactive and without variables); the first shell that
	// ends with status 0 is the shell of the session and starts without an
	// argument. The probe fails as [ErrDevcontainerRuntime] when neither
	// shell runs: user or workspace missing in the container, no shell in
	// the image, container not running.
	//
	// Once the session has begun, its end is no error: the exit status of
	// the shell, also after the container ended during the session, does
	// not reach the caller. The launcher recognises a failure of the
	// `docker` client before the session begins by the message of the
	// client and returns it as [ErrDevcontainerRuntime]: the container is
	// not running or does not exist, the user does not exist after the
	// probe, the engine does not answer, a terminal was requested on
	// standard input that is none; without Terminal also an `OCI runtime
	// exec failed` of the client (with Terminal the client owns Stdout, so
	// the launcher does not see the line). Any other message of the client
	// is part of the session.
	//
	// A cancelled context ends a session without Terminal (the launcher
	// interrupts the client) with the context error. With Terminal the
	// launcher does not observe the context once the session has begun: the
	// client owns the terminal, and Ctrl-C reaches the shell in the
	// container.
	Shell(ctx context.Context, opts DevcontainerShellOptions) error
}
