package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// ErrDevcontainerSubcommandRequired is returned by bare `u-boot
// devcontainer`, with and without `--json`: the group has no data of its own.
// Like bare `egress --json` it ends with exit 2 (usage class via
// [isUsageError]) and without an envelope.
var ErrDevcontainerSubcommandRequired = errors.New("u-boot devcontainer requires a subcommand (try `u-boot devcontainer up`, `u-boot devcontainer down`, `u-boot devcontainer shell` or `u-boot devcontainer remove`)")

// ErrDevcontainerPreviewNotApplicable is returned when `--dry-run` or
// `--diff` is given to `devcontainer up`, `devcontainer down`,
// `devcontainer shell` or `devcontainer remove`: none of the commands writes a
// project file, so there is nothing to preview. Maps to exit 2.
var ErrDevcontainerPreviewNotApplicable = errors.New("--dry-run/--diff is not applicable to `devcontainer up`, `devcontainer down`, `devcontainer shell` and `devcontainer remove` (they write no project file)")

// ErrDevcontainerShellJSONNotApplicable is returned when `--json` is given to
// `devcontainer shell`: the command opens an interactive session, and the
// standard output belongs to it. Maps to exit 2; with `--json` the rejection
// itself is written as an envelope (LH-NFA-USE-004).
var ErrDevcontainerShellJSONNotApplicable = errors.New("--json is not applicable to `devcontainer shell` (the command opens an interactive session)")

// isDevcontainerUsageError reports the usage-class sentinels of the
// devcontainer commands (exit 2); [isUsageError] takes them as one
// condition.
func isDevcontainerUsageError(err error) bool {
	return errors.Is(err, ErrDevcontainerSubcommandRequired) ||
		errors.Is(err, ErrDevcontainerPreviewNotApplicable) ||
		errors.Is(err, ErrDevcontainerShellJSONNotApplicable)
}

// devcontainerFlags bundles the flag state of `devcontainer up`,
// `devcontainer down`, `devcontainer shell` and `devcontainer remove`. DryRun
// and Diff are registered only so Cobra parses them and the RunE can reject
// them with an envelope. Volumes is the local flag of `remove`; Yes and
// NoInteractive are read from the persistent root flags by `remove` alone.
type devcontainerFlags struct {
	Recreate      bool
	Volumes       bool
	Yes           bool
	NoInteractive bool
	// Instance is the value of `--instance`; InstanceGiven is true when the
	// flag was on the command line, so an empty value is told from no flag.
	Instance      string
	InstanceGiven bool
	DryRun        bool
	Diff          bool
	Quiet         bool
	JSON          bool
}

// instanceFlagUsage is the help text of `--instance`.
const instanceFlagUsage = "act on the container of this instance of the project (a-z, 0-9 and '-', 1 to 32 characters, " +
	"beginning and ending with a letter or digit); without it, on the container without an instance"

// registerDevcontainerInstanceFlag adds `--instance` to a subcommand.
func registerDevcontainerInstanceFlag(cmd *cobra.Command, flags *devcontainerFlags) {
	cmd.Flags().StringVar(&flags.Instance, "instance", "", instanceFlagUsage)
}

// emptyInstanceError is the error of a `--instance` given with the empty
// string. The requests carry "" for no instance, so the adapter rejects the
// empty value itself; nil for every other flag state. command names the
// subcommand in the error text.
func emptyInstanceError(flags devcontainerFlags, command string) error {
	if flags.InstanceGiven && flags.Instance == "" {
		return fmt.Errorf("devcontainer %s: %w: --instance needs a name", command, driving.ErrDevcontainerInstanceInvalid)
	}
	return nil
}

// devcontainerUpData is the typed `data` carrier of `devcontainer up
// --json`; no field is omitempty. `alreadyRunning` is true when the container
// ran before the call and the command ran without `--recreate`.
type devcontainerUpData struct {
	ContainerID           string `json:"containerId"`
	RemoteUser            string `json:"remoteUser"`
	RemoteWorkspaceFolder string `json:"remoteWorkspaceFolder"`
	AlreadyRunning        bool   `json:"alreadyRunning"`
	// The name of the source of the git token of the resolution (LH-FA-
	// DEV-009); empty with no token or no resolution.
	GitTokenSource string `json:"gitTokenSource"`
}

// devcontainerDownData is the typed `data` carrier of `devcontainer down
// --json`; no field is omitempty. `containerId` is the empty string when
// there is no container for the project.
type devcontainerDownData struct {
	ContainerID string `json:"containerId"`
	Stopped     bool   `json:"stopped"`
}

// devcontainerRemoveData is the typed `data` carrier of `devcontainer remove
// --json`; no field is omitempty. `containerId` is the empty string when there
// is no container for the project, and `removedVolumeNames` is an array, empty
// when no volume was removed.
type devcontainerRemoveData struct {
	ContainerID        string   `json:"containerId"`
	Stopped            bool     `json:"stopped"`
	Removed            bool     `json:"removed"`
	RemovedVolumeNames []string `json:"removedVolumeNames"`
}

// newDevcontainerCommand builds the `u-boot devcontainer` group
// (LH-FA-DEV-010). Children: `up`, `down`, `shell`, `remove`.
func newDevcontainerCommand(a *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "devcontainer",
		Short: "Start, stop, enter and remove the devcontainer of the project without an editor",
		Long: `Start, stop, enter and remove the devcontainer of the project without a
devcontainer-capable editor, without Node.js and without the Dev Container CLI
on the host; Docker (client and local engine) is enough.

Subcommands:

  u-boot devcontainer up [--recreate]   start the devcontainer
  u-boot devcontainer down              stop the devcontainer
  u-boot devcontainer shell             start the devcontainer and open a shell in it
  u-boot devcontainer remove [--volumes]  remove the container, with --volumes its volumes

Each subcommand takes --instance NAME to act on one of several devcontainers of
the same project folder; see the help of the subcommands.`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return ErrDevcontainerSubcommandRequired
		},
	}
	cmd.AddCommand(newDevcontainerUpCommand(a), newDevcontainerDownCommand(a), newDevcontainerShellCommand(a), newDevcontainerRemoveCommand(a))
	return cmd
}

func newDevcontainerUpCommand(a *App) *cobra.Command {
	flags := &devcontainerFlags{}
	cmd := &cobra.Command{
		Use:   "up",
		Short: "Start the devcontainer of the project (Dev Container CLI in a helper container)",
		Long: `Start the devcontainer that .devcontainer/devcontainer.json in the current
directory describes. The command runs the Dev Container CLI in a short-lived
helper container; the host needs Docker (client and local engine), nothing else.

The helper image is built locally from pinned sources (the Dev Container CLI
from the npm registry with a lock file, the base images node and docker by
digest). The first call needs network access and builds the image; later calls
reuse it.

The helper container holds the engine socket of the host for the duration of
the call (root-equivalent access to the engine). The command says so before it
starts the helper, also under --quiet. The devcontainer configuration that
u-boot generates mounts the socket if a feature that does so is enabled, such
as the built-in feature docker-cli (the sandbox profile rejects that feature);
a hand-written .devcontainer/devcontainer.json can mount it, and its
initializeCommand runs with this access: read the file before the first start
if somebody else can write to the project folder.

Only the variables that devcontainer.json references as ${localEnv:NAME} and
that are set on the host reach the helper container, by name. With the sandbox
profile and the switch devcontainer.sandbox.hostGitToken switched on, the
command also resolves the git token on the host when GIT_TOKEN is not set
there (the gh session of the host answers; the source is named "gh auth
token"), names the source it used, never the value, and points out without
aborting when no source is found. The command does not take a value into a
file or an image;
with containerEnv the resolved value stays in the metadata of the started
container, remoteEnv is the form meant for credentials.

--recreate removes an existing container and creates it anew; named volumes
stay. It is the editor-neutral form of "devcontainer up
--remove-existing-container".

After the start the command reads the status files that the container's start
scripts wrote (through the capabilities configured in u-boot.yaml; without
u-boot.yaml there is nothing to read). It reports a degradation fallback of the
sandbox profile as a warning. A start script that reported a failure ends the
command with exit 11, also when the container was already running and this
call did not run the script; "u-boot devcontainer up --recreate" creates the
container anew. If the container was already running and the command runs
without --recreate, the result line says "already running" and the JSON data
carry alreadyRunning: true; with --recreate the line says "started" and
alreadyRunning is false.

--instance NAME starts a further devcontainer of the same project folder, next
to the one without an instance; each instance has its own container and its own
volumes. NAME has 1 to 32 characters of a-z, 0-9 and '-' and begins and ends
with a letter or digit. Without --instance the command acts on the container
without an instance only; if the project has only containers of instances, it
ends with exit 10 and starts nothing.

The started container carries the name u-boot-<project> (with --instance NAME,
u-boot-<project>-<NAME>); the success line shows it beside the container ID.
The devcontainer.json of the project is not changed by the name.

Exit codes:
  0   success
  2   CLI validation (unknown flag, arguments, --dry-run/--diff)
  10  no .devcontainer/devcontainer.json in the current directory; a directory
      path with a comma, double quote or line break; a ${localEnv:NAME}
      reference whose name is not a plain variable name; an instance name that
      is not admissible; a container-name part the engine would refuse for
      --name, or a container of another project folder that already holds the
      name; no --instance although the project has only
      containers of instances
  11  no docker binary on the PATH (for example in the released container
      image), the daemon does not answer, the engine is not local, or a start
      script of the container aborted for an environment reason (also in an
      earlier start of the container, if u-boot.yaml configures the
      capability)
  12  any other failure of the build or the start
  14  filesystem read failure

Examples:
  u-boot devcontainer up
  u-boot devcontainer up --instance agent-1
  u-boot devcontainer up --recreate --json`,
		Args: jsonArgsValidator(a, "devcontainer", "up", cobra.NoArgs, mapDevcontainerUpErrorToDiagnostic, false),
		RunE: func(cmd *cobra.Command, _ []string) error {
			flags.Quiet = a.quiet()
			flags.JSON = a.json
			flags.InstanceGiven = cmd.Flags().Changed("instance")
			return runDevcontainerUp(cmd.Context(), cmd.OutOrStdout(), cmd.ErrOrStderr(), *flags, a.devcontainerUseCase, a.getwd)
		},
	}
	cmd.Flags().BoolVar(&flags.Recreate, "recreate", false,
		"remove an existing container and create it anew; named volumes stay")
	registerDevcontainerInstanceFlag(cmd, flags)
	registerDevcontainerPreviewFlags(cmd, flags)
	return cmd
}

func newDevcontainerDownCommand(a *App) *cobra.Command {
	flags := &devcontainerFlags{}
	cmd := &cobra.Command{
		Use:   "down",
		Short: "Stop the devcontainer of the project (it stays; not the Compose \"down\")",
		Long: `Stop the container that belongs to the project in the current directory.

Unlike "u-boot down", which tears the Compose environment down, this command
only stops the container: it removes neither the container nor a volume, and it
needs no helper container and no socket mount. A later "u-boot devcontainer
up" starts the same container again. "u-boot devcontainer remove" removes the
container.

With no container for the project, or with one that is already stopped, the
command ends with exit 0.

--instance NAME stops the container of that instance (1 to 32 characters of
a-z, 0-9 and '-', beginning and ending with a letter or digit); with no
container of that instance the command ends with exit 0. Without --instance the
command acts on the container without an instance only; if the project has
only containers of instances, it ends with exit 10 and stops nothing.

Exit codes:
  0   success, also when there is nothing to stop
  2   CLI validation (unknown flag, arguments, --dry-run/--diff)
  10  an instance name that is not admissible; no --instance although the
      project has only containers of instances
  11  no docker binary on the PATH, the daemon does not answer, or the engine
      is not local
  12  docker stop failed

Examples:
  u-boot devcontainer down
  u-boot devcontainer down --instance agent-1
  u-boot devcontainer down --json`,
		Args: jsonArgsValidator(a, "devcontainer", "down", cobra.NoArgs, mapDevcontainerDownErrorToDiagnostic, false),
		RunE: func(cmd *cobra.Command, _ []string) error {
			flags.Quiet = a.quiet()
			flags.JSON = a.json
			flags.InstanceGiven = cmd.Flags().Changed("instance")
			return runDevcontainerDown(cmd.Context(), cmd.OutOrStdout(), *flags, a.devcontainerUseCase, a.getwd)
		},
	}
	registerDevcontainerInstanceFlag(cmd, flags)
	registerDevcontainerPreviewFlags(cmd, flags)
	return cmd
}

func newDevcontainerShellCommand(a *App) *cobra.Command {
	flags := &devcontainerFlags{}
	cmd := &cobra.Command{
		Use:   "shell",
		Short: "Start the devcontainer of the project and open a shell in it",
		Long: `Start the devcontainer that .devcontainer/devcontainer.json in the current
directory describes, as "u-boot devcontainer up" does without --recreate, and
open an interactive shell in the container. The host needs Docker (client and
local engine), nothing else.

The shell runs as the remote user in the workspace folder that the start
reports. Of the host variables that the generated devcontainer.json hands to
the container as "NAME": "${localEnv:NAME}" entries, those that are set on the
host are set in the shell; the others are absent. Before the session begins,
the command names the container of the session beside its ID on the standard
error stream.

The helper container of the start holds the engine socket of the host for the
duration of the start (root-equivalent access to the engine). The command says
so before it starts the helper, also under --quiet. The shell runs without a
helper container and without a socket that this command mounts.

Everything u-boot itself prints (that notice, the progress of the start,
warnings and errors) goes to standard error; standard output belongs to the
shell. --quiet holds the progress back. With a terminal on standard input the
shell is attached to it; otherwise it reads its commands from standard input.

A start script of the container that reported a failure ends the command with
exit 11 and opens no shell, also when the container was already running; this
needs the u-boot.yaml of the project, without it no status file is read.
"u-boot devcontainer up --recreate" creates the container anew. Once the shell
has begun, the command ends with exit 0 whatever the exit status of the shell
was; to run one command in the container and get its exit code, use "docker
exec".

--instance NAME opens the shell in the container of that instance, and starts
it as "u-boot devcontainer up --instance NAME" does (1 to 32 characters of a-z,
0-9 and '-', beginning and ending with a letter or digit). Without --instance
the command acts on the container without an instance only; if the project has
only containers of instances, it ends with exit 10 and starts nothing.

--json, --dry-run, --diff and arguments are rejected with exit 2. --yes and
--no-interactive have no effect: the command asks for no confirmation.

Exit codes:
  0   the shell has begun and ended
  2   CLI validation (unknown flag, arguments, --json, --dry-run/--diff)
  10  no .devcontainer/devcontainer.json in the current directory; a directory
      path with a comma, double quote or line break; a ${localEnv:NAME}
      reference whose name is not a plain variable name; an instance name that
      is not admissible; a container-name part the engine would refuse for
      --name, or a container of another project folder that already holds the
      name; no --instance although the project has only
      containers of instances
  11  no docker binary on the PATH (for example in the released container
      image), the daemon does not answer, the engine is not local, or a start
      script of the container aborted for an environment reason (also in an
      earlier start of the container, if u-boot.yaml configures the
      capability)
  12  any other failure of the build or the start; the check before the shell
      opens fails (user, workspace folder or shell missing in the container)
  14  filesystem read failure

Examples:
  u-boot devcontainer shell
  u-boot devcontainer shell --instance agent-1
  u-boot --quiet devcontainer shell`,
		Args: jsonArgsValidator(a, "devcontainer", "shell", cobra.NoArgs, mapDevcontainerShellErrorToDiagnostic, false),
		RunE: func(cmd *cobra.Command, _ []string) error {
			flags.Quiet = a.quiet()
			flags.JSON = a.json
			flags.InstanceGiven = cmd.Flags().Changed("instance")
			return runDevcontainerShell(cmd.Context(),
				devcontainerStreams{in: cmd.InOrStdin(), stdout: cmd.OutOrStdout(), stderr: cmd.ErrOrStderr()},
				*flags, a.devcontainerUseCase, a.getwd, a.terminalProbe)
		},
	}
	registerDevcontainerInstanceFlag(cmd, flags)
	registerDevcontainerPreviewFlags(cmd, flags)
	return cmd
}

func newDevcontainerRemoveCommand(a *App) *cobra.Command {
	flags := &devcontainerFlags{}
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove the container of the devcontainer of the project; with --volumes its volumes (data loss)",
		Long: `Remove the container that belongs to the project in the current directory. A
running container is stopped first. The command needs no helper container, no
socket mount and no project file; the host needs Docker (client and local
engine), nothing else. It is not "u-boot remove", which disables an add-on
service of the Compose environment, and not "u-boot devcontainer down", which
stops the container and removes nothing.

Without --volumes the command removes the container only and asks nothing. The
volumes of the container stay, and a later "u-boot devcontainer up" mounts the
named ones again with their content; an anonymous volume (a VOLUME of the
Dockerfile of the project) stays behind unused.

--volumes also removes the volumes that the container mounts, named and
anonymous (in the sandbox profile the workspace volume, with the nested runtime
the volume of its storage, with the feature ollama the volume of the models).
THE DATA IN THESE VOLUMES IS PERMANENTLY DELETED. The command reads the volumes
from the mounts of the container, never from a name pattern, and names them
before it asks. It asks for confirmation on the standard error stream
("[y/N]"; only y and yes agree, an empty line or the end of the input declines):
  - --yes                     confirm without asking
  - --no-interactive          exit 10, nothing is stopped or removed
  - --json                    never asks, exit 10 unless --yes is given
The question comes only when the project has a container to remove; with no
container the command ends with exit 0 and asks nothing. A volume that another
container holds cannot be removed: the command removes the others, names the
one that remains, and ends with exit 12. A volume of a container that no longer
exists is not found by the command; remove such a volume with "docker volume
rm".

--instance NAME removes the container of that instance (1 to 32 characters of
a-z, 0-9 and '-', beginning and ending with a letter or digit) and leaves the
other containers of the project and their volumes alone; with no container of
that instance the command ends with exit 0. Without --instance the command acts
on the container without an instance only; if the project has only containers
of instances, it ends with exit 10 and removes nothing, also with --volumes
--yes.

--quiet holds back the result lines; the question and every error stay.
--dry-run, --diff and arguments are rejected with exit 2, as is --yes together
with --no-interactive.

Exit codes:
  0   success, also when there is nothing to remove
  2   CLI validation (unknown flag, arguments, --dry-run/--diff, --yes together
      with --no-interactive)
  10  an instance name that is not admissible; no --instance although the
      project has only containers of instances; --volumes was not confirmed
      (declined, answer not readable, --no-interactive or --json without --yes)
  11  no docker binary on the PATH (for example in the released container
      image), the daemon does not answer, the engine is not local, or the
      engine could not be queried
  12  docker stop, docker rm or docker volume rm failed, or a volume remains

Examples:
  u-boot devcontainer remove
  u-boot devcontainer remove --volumes
  u-boot devcontainer remove --instance agent-1 --volumes --yes
  u-boot devcontainer remove --json`,
		Args: jsonArgsValidator(a, "devcontainer", "remove", cobra.NoArgs, mapDevcontainerRemoveErrorToDiagnostic, false),
		RunE: func(cmd *cobra.Command, _ []string) error {
			flags.Quiet = a.quiet()
			flags.JSON = a.json
			flags.Yes = a.yes
			flags.NoInteractive = a.noInteractive
			flags.InstanceGiven = cmd.Flags().Changed("instance")
			return runDevcontainerRemove(cmd.Context(), cmd.OutOrStdout(), *flags, a.devcontainerUseCase, a.getwd)
		},
	}
	cmd.Flags().BoolVar(&flags.Volumes, "volumes", false,
		"also remove the volumes that the container mounts (named and anonymous; the data is permanently deleted)")
	registerDevcontainerInstanceFlag(cmd, flags)
	registerDevcontainerPreviewFlags(cmd, flags)
	return cmd
}

func registerDevcontainerPreviewFlags(cmd *cobra.Command, flags *devcontainerFlags) {
	cmd.Flags().BoolVar(&flags.DryRun, "dry-run", false,
		"not applicable (the command writes no project file); rejected with exit 2")
	cmd.Flags().BoolVar(&flags.Diff, "diff", false,
		"not applicable (the command writes no project file); rejected with exit 2")
}

// reportDevcontainerError is the error gate of the subcommands. In JSON mode
// it writes one envelope whose `diagnostics[]` holds the warnings already
// raised (the socket disclosure) followed by the error, and returns err for
// [ExitCode]; in human mode it returns err.
func reportDevcontainerError(
	out io.Writer, err error, flags devcontainerFlags, subcommand string,
	mapErr func(error) diagnosticItem, warnings []driving.WarningEntry,
) error {
	if !flags.JSON {
		return err
	}
	diags := append(mapWarningsToDiagnostics(warnings), mapErr(err))
	if envErr := writeEnvelope(out, newMinimalEnvelope("devcontainer", subcommand, diags, ExitCode(err))); envErr != nil {
		return envErr
	}
	return err
}

// runDevcontainerUp delegates to the use case and renders the result.
//
// The socket disclosure reaches the user through the request's OnWarning
// before the helper container starts: human mode writes it to stdout at once
// (it is a security relaxation, so it stays under `--quiet`), JSON mode
// collects it for the envelope. Progress goes to stderr; under `--quiet` it is
// held back and written only when the call fails, under `--json` it is
// dropped.
func runDevcontainerUp(
	ctx context.Context, stdout, stderr io.Writer, flags devcontainerFlags,
	uc driving.DevcontainerUseCase, getwd func() (string, error),
) error {
	var raised []driving.WarningEntry
	fail := func(err error) error {
		return reportDevcontainerError(stdout, err, flags, "up", mapDevcontainerUpErrorToDiagnostic, raised)
	}
	if flags.DryRun || flags.Diff {
		return fail(ErrDevcontainerPreviewNotApplicable)
	}
	if uc == nil {
		return fail(errors.New("devcontainer up: use case not wired"))
	}
	if err := emptyInstanceError(flags, "up"); err != nil {
		return fail(err)
	}
	cwd, err := getwd()
	if err != nil {
		return fail(fmt.Errorf("determine working directory: %w", err))
	}
	var sink io.Writer
	var held *bytes.Buffer
	switch {
	case flags.JSON:
		// no progress sink: the envelope is the only output
	case flags.Quiet:
		held = &bytes.Buffer{}
		sink = held
	default:
		sink = stderr
	}
	resp, err := uc.Up(ctx, driving.DevcontainerUpRequest{
		BaseDir:      cwd,
		Recreate:     flags.Recreate,
		Instance:     flags.Instance,
		ProgressSink: sink,
		OnWarning: func(w driving.WarningEntry) {
			raised = append(raised, w)
			if !flags.JSON {
				writeWarnings(stdout, []driving.WarningEntry{w})
			}
		},
	})
	if err != nil {
		if held != nil {
			_, _ = stderr.Write(held.Bytes())
		}
		return fail(sanitizeBaseDir(err, cwd))
	}
	if flags.JSON {
		data := devcontainerUpData{
			ContainerID: resp.ContainerID, RemoteUser: resp.RemoteUser,
			RemoteWorkspaceFolder: resp.RemoteWorkspaceFolder, AlreadyRunning: resp.AlreadyRunning,
			GitTokenSource: resp.GitTokenSource,
		}
		diags := mapWarningsToDiagnostics(append(raised, resp.Warnings...))
		return writeEnvelope(stdout, newDataEnvelope("devcontainer", "up", data, diags, 0))
	}
	if !flags.Quiet {
		fmt.Fprintf(stdout, "%s, remote user %s, workspace %s\n",
			devcontainerStartLine(resp), resp.RemoteUser, resp.RemoteWorkspaceFolder)
	}
	writeWarnings(stdout, visibleWarnings(resp.Warnings, flags.Quiet))
	return nil
}

// runDevcontainerDown delegates to the use case and renders the result. The
// outcome line is the result of the command and disappears under `--quiet`
// like the success line of every modifying command.
func runDevcontainerDown(
	ctx context.Context, stdout io.Writer, flags devcontainerFlags,
	uc driving.DevcontainerUseCase, getwd func() (string, error),
) error {
	fail := func(err error) error {
		return reportDevcontainerError(stdout, err, flags, "down", mapDevcontainerDownErrorToDiagnostic, nil)
	}
	if flags.DryRun || flags.Diff {
		return fail(ErrDevcontainerPreviewNotApplicable)
	}
	if uc == nil {
		return fail(errors.New("devcontainer down: use case not wired"))
	}
	if err := emptyInstanceError(flags, "down"); err != nil {
		return fail(err)
	}
	cwd, err := getwd()
	if err != nil {
		return fail(fmt.Errorf("determine working directory: %w", err))
	}
	resp, err := uc.Down(ctx, driving.DevcontainerDownRequest{BaseDir: cwd, Instance: flags.Instance})
	if err != nil {
		return fail(sanitizeBaseDir(err, cwd))
	}
	if flags.JSON {
		data := devcontainerDownData{ContainerID: resp.ContainerID, Stopped: resp.Stopped}
		return writeEnvelope(stdout, newDataEnvelope("devcontainer", "down", data, nil, 0))
	}
	if flags.Quiet {
		return nil
	}
	switch {
	case resp.Stopped:
		fmt.Fprintf(stdout, "Devcontainer stopped: container %s (not removed)\n", shortContainerID(resp.ContainerID))
	case resp.ContainerID != "":
		fmt.Fprintf(stdout, "Devcontainer container %s is not running; nothing to stop\n", shortContainerID(resp.ContainerID))
	default:
		fmt.Fprintln(stdout, "No devcontainer container found for this project; nothing to stop")
	}
	return nil
}

// runDevcontainerRemove delegates to the use case and renders the result. The
// result lines are the result of the command and disappear under `--quiet`
// like the success line of every modifying command; the question of
// `--volumes` is the confirmer's and stays. `--dry-run`, `--diff` and `--yes`
// together with `--no-interactive` end with exit 2 before anything else, an
// instance flag given empty with exit 10.
func runDevcontainerRemove(
	ctx context.Context, stdout io.Writer, flags devcontainerFlags,
	uc driving.DevcontainerUseCase, getwd func() (string, error),
) error {
	fail := func(err error) error {
		return reportDevcontainerError(stdout, err, flags, "remove", mapDevcontainerRemoveErrorToDiagnostic, nil)
	}
	switch {
	case flags.DryRun || flags.Diff:
		return fail(ErrDevcontainerPreviewNotApplicable)
	case flags.Yes && flags.NoInteractive:
		return fail(ErrConflictingModeFlags)
	case uc == nil:
		return fail(errors.New("devcontainer remove: use case not wired"))
	}
	if err := emptyInstanceError(flags, "remove"); err != nil {
		return fail(err)
	}
	cwd, err := getwd()
	if err != nil {
		return fail(fmt.Errorf("determine working directory: %w", err))
	}
	resp, err := uc.Remove(ctx, driving.DevcontainerRemoveRequest{
		BaseDir:          cwd,
		Instance:         flags.Instance,
		Volumes:          flags.Volumes,
		AssumeYes:        flags.Yes,
		NonInteractive:   flags.NoInteractive,
		SilenceConfirmer: flags.JSON,
	})
	if err != nil {
		return fail(sanitizeBaseDir(err, cwd))
	}
	if flags.JSON {
		names := resp.RemovedVolumeNames
		if names == nil {
			names = []string{}
		}
		data := devcontainerRemoveData{ContainerID: resp.ContainerID, Stopped: resp.Stopped, Removed: resp.Removed, RemovedVolumeNames: names}
		return writeEnvelope(stdout, newDataEnvelope("devcontainer", "remove", data, nil, 0))
	}
	if !flags.Quiet {
		writeDevcontainerRemoveResult(stdout, resp, flags.Volumes)
	}
	return nil
}

// writeDevcontainerRemoveResult renders the result lines of `remove`: the
// container, with `--volumes` the volumes removed, and the hint of a call with
// no container to remove.
func writeDevcontainerRemoveResult(out io.Writer, resp driving.DevcontainerRemoveResponse, volumes bool) {
	if !resp.Removed {
		fmt.Fprintln(out, "No devcontainer container found for this project; nothing to remove")
		if volumes {
			fmt.Fprintln(out, "Volumes of a container that no longer exists are not removed by u-boot.")
		}
		return
	}
	line := "Devcontainer removed: container " + shortContainerID(resp.ContainerID)
	if resp.Stopped {
		line += " (stopped first)"
	}
	fmt.Fprintln(out, line)
	if !volumes {
		return
	}
	if len(resp.RemovedVolumeNames) == 0 {
		fmt.Fprintln(out, "No volume removed (the container held no volume)")
		return
	}
	fmt.Fprintf(out, "Volumes removed: %d\n", len(resp.RemovedVolumeNames))
	for _, name := range resp.RemovedVolumeNames {
		fmt.Fprintf(out, "  - %s\n", name)
	}
}

// devcontainerStreams are the streams of `devcontainer shell`: the standard
// input of the session and the two output streams of the command.
type devcontainerStreams struct {
	in     io.Reader
	stdout io.Writer
	stderr io.Writer
}

// runDevcontainerShell delegates to the use case, which starts the container
// and opens the session. Everything u-boot itself writes goes to stderr:
// the socket disclosure and the fallback warnings (as they arrive, before the
// session begins), the progress of the start (held back under `--quiet` and
// written only when the call fails) and the error. Stdout belongs to the
// session, and no success line follows it. The streams reach the session
// unchanged (in production the files of the process, so a terminal on the
// standard input is passed on); stdout and stderr must be two distinct
// writers, because the session writes to both from separate goroutines.
//
// `--dry-run`, `--diff` and `--json` end with exit 2 before anything else,
// with `--json` as an envelope; `--yes` and `--no-interactive` are not read,
// the command asks for no confirmation.
func runDevcontainerShell(
	ctx context.Context, streams devcontainerStreams, flags devcontainerFlags,
	uc driving.DevcontainerUseCase, getwd func() (string, error), terminalProbe func(io.Reader) bool,
) error {
	fail := func(err error) error {
		return reportDevcontainerError(streams.stdout, err, flags, "shell", mapDevcontainerShellErrorToDiagnostic, nil)
	}
	switch {
	case flags.DryRun || flags.Diff:
		return fail(ErrDevcontainerPreviewNotApplicable)
	case flags.JSON:
		return fail(ErrDevcontainerShellJSONNotApplicable)
	case uc == nil:
		return fail(errors.New("devcontainer shell: use case not wired"))
	}
	if err := emptyInstanceError(flags, "shell"); err != nil {
		return fail(err)
	}
	cwd, err := getwd()
	if err != nil {
		return fail(fmt.Errorf("determine working directory: %w", err))
	}
	sink := streams.stderr
	var held *bytes.Buffer
	if flags.Quiet {
		held = &bytes.Buffer{}
		sink = held
	}
	err = uc.Shell(ctx, driving.DevcontainerShellRequest{
		BaseDir:      cwd,
		Instance:     flags.Instance,
		ProgressSink: sink,
		OnWarning: func(w driving.WarningEntry) {
			writeWarnings(streams.stderr, visibleWarnings([]driving.WarningEntry{w}, flags.Quiet))
		},
		// The start line names the container of the session beside its ID and
		// goes to the standard error stream: the session then owns the
		// standard output. Under `--quiet` it is a result of a modifying
		// command — it waits in the held buffer, only a failed call shows it.
		OnStart: func(resp driving.DevcontainerUpResponse) {
			if resp.ContainerName == "" {
				return
			}
			line := devcontainerStartLine(resp) + "\n"
			if held != nil {
				held.WriteString(line)
				return
			}
			fmt.Fprint(streams.stderr, line)
		},
		Stdin:    streams.in,
		Stdout:   streams.stdout,
		Stderr:   streams.stderr,
		Terminal: terminalProbe(streams.in),
	})
	if err != nil {
		if held != nil {
			_, _ = streams.stderr.Write(held.Bytes())
		}
		return fail(sanitizeBaseDir(err, cwd))
	}
	return nil
}

// writeWarnings renders warnings in the human form of the other commands.
func writeWarnings(out io.Writer, ws []driving.WarningEntry) {
	for _, w := range ws {
		fmt.Fprintf(out, "Warning: %s\n", w.Message)
	}
}

// devcontainerStartLine is the head of the start line: the outcome, the
// container ID, and, when the launcher handed the CLI a per-call
// configuration, the name of the container; the caller appends what the
// command shows besides it. The name stands beside the ID when the container
// carries one, and is absent when it does not (LH-FA-DEV-010).
func devcontainerStartLine(resp driving.DevcontainerUpResponse) string {
	outcome := "started"
	if resp.AlreadyRunning {
		outcome = "already running"
	}
	line := fmt.Sprintf("Devcontainer %s: container %s", outcome, shortContainerID(resp.ContainerID))
	if resp.ContainerName != "" {
		line += ", name " + resp.ContainerName
	}
	return line
}

// shortContainerID is the first 12 characters of a container ID, the form
// `docker ps` shows.
func shortContainerID(id string) string {
	if len(id) > 12 {
		return id[:12]
	}
	return id
}

// mapDevcontainerEngineSentinel classifies the engine-side sentinels the
// subcommands can return: the environment class (exit 11), the runtime class
// (exit 12) and the start script abort (exit 11). The code is the
// requirement that defines the classes.
func mapDevcontainerEngineSentinel(err error) (code string, matched bool) {
	switch {
	case errors.Is(err, driven.ErrDockerUnavailable),
		errors.Is(err, driven.ErrContainerEngineNotLocal),
		errors.Is(err, driving.ErrDevcontainerStartScriptFailed),
		errors.Is(err, driven.ErrDevcontainerRuntime):
		return "LH-FA-DEV-010", true
	}
	return "", false
}

// mapDevcontainerUpErrorToDiagnostic maps a `devcontainer up` error to its
// diagnostic. Switch order: filesystem class first (exit 14), then the engine
// sentinels, then the validation sentinels (exit 10), then the rest.
func mapDevcontainerUpErrorToDiagnostic(err error) diagnosticItem {
	if errors.Is(err, driving.ErrDevcontainerFileSystem) {
		return diagnosticItem{Level: "error", Code: "LH-NFA-REL-003", Message: err.Error()}
	}
	if code, matched := mapDevcontainerEngineSentinel(err); matched {
		return diagnosticItem{Level: "error", Code: code, Message: err.Error()}
	}
	if isDevcontainerValidationError(err) {
		return diagnosticItem{Level: "error", Code: "LH-FA-DEV-010", Message: err.Error()}
	}
	return diagnosticItem{Level: "error", Code: "LH-FA-CLI-006", Message: err.Error()}
}

// mapDevcontainerDownErrorToDiagnostic maps a `devcontainer down` error to
// its diagnostic: the engine sentinels, then the two instance sentinels (exit
// 10), then the rest. `down` reads no project file, so the sentinels about the
// configuration, the folder and the variable names do not occur.
func mapDevcontainerDownErrorToDiagnostic(err error) diagnosticItem {
	if code, matched := mapDevcontainerEngineSentinel(err); matched {
		return diagnosticItem{Level: "error", Code: code, Message: err.Error()}
	}
	if isDevcontainerInstanceError(err) {
		return diagnosticItem{Level: "error", Code: "LH-FA-DEV-010", Message: err.Error()}
	}
	return diagnosticItem{Level: "error", Code: "LH-FA-CLI-006", Message: err.Error()}
}

// mapDevcontainerRemoveErrorToDiagnostic maps a `devcontainer remove` error to
// its diagnostic: the engine sentinels (exit 11 and 12), then the instance
// sentinels and the two sentinels of the missing confirmation (exit 10, all
// with the code of LH-FA-DEV-010), the exclusion of `--yes` and
// `--no-interactive` (exit 2, LH-FA-CLI-005A), then the rest. `remove` reads
// no project file, so neither the filesystem class (exit 14) nor the sentinels
// about the configuration, the folder and the variable names occur.
func mapDevcontainerRemoveErrorToDiagnostic(err error) diagnosticItem {
	if code, matched := mapDevcontainerEngineSentinel(err); matched {
		return diagnosticItem{Level: "error", Code: code, Message: err.Error()}
	}
	if isDevcontainerInstanceError(err) ||
		errors.Is(err, driving.ErrConfirmationRequired) || errors.Is(err, driving.ErrConfirmerUnavailable) {
		return diagnosticItem{Level: "error", Code: "LH-FA-DEV-010", Message: err.Error()}
	}
	if errors.Is(err, ErrConflictingModeFlags) {
		return diagnosticItem{Level: "error", Code: "LH-FA-CLI-005A", Message: err.Error()}
	}
	return diagnosticItem{Level: "error", Code: "LH-FA-CLI-006", Message: err.Error()}
}

// mapDevcontainerShellErrorToDiagnostic maps a `devcontainer shell` error to
// its diagnostic. `--json` is rejected before the use case runs, so the only
// errors that reach an envelope are the usage errors of the command (exit 2,
// code LH-FA-CLI-006); an engine sentinel of the use case never appears in
// one, and the mapping does not pretend it could.
func mapDevcontainerShellErrorToDiagnostic(err error) diagnosticItem {
	return diagnosticItem{Level: "error", Code: "LH-FA-CLI-006", Message: err.Error()}
}

// isDevcontainerValidationError reports the exit-10 sentinels of `devcontainer
// up`, `down`, `shell` and `remove`; [isValidationError] takes them as one condition.
func isDevcontainerValidationError(err error) bool {
	return errors.Is(err, driving.ErrDevcontainerConfigMissing) ||
		errors.Is(err, driving.ErrDevcontainerProjectFolderInvalid) ||
		errors.Is(err, driving.ErrDevcontainerEnvNameInvalid) ||
		errors.Is(err, driving.ErrDevcontainerContainerNameInvalid) ||
		errors.Is(err, driving.ErrDevcontainerContainerNameConflict) ||
		isDevcontainerInstanceError(err)
}

// isDevcontainerInstanceError reports the two exit-10 sentinels about the
// instance, which `up`, `down`, `shell` and `remove` share.
func isDevcontainerInstanceError(err error) bool {
	return errors.Is(err, driving.ErrDevcontainerInstanceInvalid) ||
		errors.Is(err, driving.ErrDevcontainerOnlyInstances)
}
