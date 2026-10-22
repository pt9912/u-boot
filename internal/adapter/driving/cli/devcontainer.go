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
var ErrDevcontainerSubcommandRequired = errors.New("u-boot devcontainer requires a subcommand (try `u-boot devcontainer up`, `u-boot devcontainer down` or `u-boot devcontainer shell`)")

// ErrDevcontainerPreviewNotApplicable is returned when `--dry-run` or
// `--diff` is given to `devcontainer up`, `devcontainer down` or
// `devcontainer shell`: none of the commands writes a project file, so there
// is nothing to preview. Maps to exit 2.
var ErrDevcontainerPreviewNotApplicable = errors.New("--dry-run/--diff is not applicable to `devcontainer up`, `devcontainer down` and `devcontainer shell` (they write no project file)")

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
// `devcontainer down` and `devcontainer shell`. DryRun and Diff are registered
// only so Cobra parses them and the RunE can reject them with an envelope.
type devcontainerFlags struct {
	Recreate bool
	DryRun   bool
	Diff     bool
	Quiet    bool
	JSON     bool
}

// devcontainerUpData is the typed `data` carrier of `devcontainer up
// --json`; no field is omitempty. `alreadyRunning` is true when the container
// ran before the call and the command ran without `--recreate`.
type devcontainerUpData struct {
	ContainerID           string `json:"containerId"`
	RemoteUser            string `json:"remoteUser"`
	RemoteWorkspaceFolder string `json:"remoteWorkspaceFolder"`
	AlreadyRunning        bool   `json:"alreadyRunning"`
}

// devcontainerDownData is the typed `data` carrier of `devcontainer down
// --json`; no field is omitempty. `containerId` is the empty string when
// there is no container for the project.
type devcontainerDownData struct {
	ContainerID string `json:"containerId"`
	Stopped     bool   `json:"stopped"`
}

// newDevcontainerCommand builds the `u-boot devcontainer` group
// (LH-FA-DEV-010). Children: `up`, `down`, `shell`.
func newDevcontainerCommand(a *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "devcontainer",
		Short: "Start, stop and enter the devcontainer of the project without an editor",
		Long: `Start, stop and enter the devcontainer of the project without a
devcontainer-capable editor, without Node.js and without the Dev Container CLI
on the host; Docker (client and local engine) is enough.

Subcommands:

  u-boot devcontainer up [--recreate]   start the devcontainer
  u-boot devcontainer down              stop the devcontainer
  u-boot devcontainer shell             start the devcontainer and open a shell in it`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return ErrDevcontainerSubcommandRequired
		},
	}
	cmd.AddCommand(newDevcontainerUpCommand(a), newDevcontainerDownCommand(a), newDevcontainerShellCommand(a))
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
that are set on the host reach the helper container, by name. The command
does not take a value into a file or an image; with containerEnv the resolved
value stays in the metadata of the started container, remoteEnv is the form
meant for credentials.

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

Exit codes (LH-FA-CLI-006):
  0   success
  2   CLI validation (unknown flag, arguments, --dry-run/--diff)
  10  no .devcontainer/devcontainer.json in the current directory; a directory
      path with a comma, double quote or line break; a ${localEnv:NAME}
      reference whose name is not a plain variable name
  11  no docker binary on the PATH (for example in the released container
      image), the daemon does not answer, the engine is not local, or a start
      script of the container aborted for an environment reason (also in an
      earlier start of the container, if u-boot.yaml configures the
      capability)
  12  any other failure of the build or the start
  14  filesystem read failure

Examples:
  u-boot devcontainer up
  u-boot devcontainer up --recreate --json`,
		Args: jsonArgsValidator(a, "devcontainer", "up", cobra.NoArgs, mapDevcontainerUpErrorToDiagnostic, false),
		RunE: func(cmd *cobra.Command, _ []string) error {
			flags.Quiet = a.quiet()
			flags.JSON = a.json
			return runDevcontainerUp(cmd.Context(), cmd.OutOrStdout(), cmd.ErrOrStderr(), *flags, a.devcontainerUseCase, a.getwd)
		},
	}
	cmd.Flags().BoolVar(&flags.Recreate, "recreate", false,
		"remove an existing container and create it anew; named volumes stay")
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
up" starts the same container again.

With no container for the project, or with one that is already stopped, the
command ends with exit 0.

Exit codes (LH-FA-CLI-006):
  0   success, also when there is nothing to stop
  2   CLI validation (unknown flag, arguments, --dry-run/--diff)
  11  no docker binary on the PATH, the daemon does not answer, or the engine
      is not local
  12  docker stop failed

Examples:
  u-boot devcontainer down
  u-boot devcontainer down --json`,
		Args: jsonArgsValidator(a, "devcontainer", "down", cobra.NoArgs, mapDevcontainerDownErrorToDiagnostic, false),
		RunE: func(cmd *cobra.Command, _ []string) error {
			flags.Quiet = a.quiet()
			flags.JSON = a.json
			return runDevcontainerDown(cmd.Context(), cmd.OutOrStdout(), *flags, a.devcontainerUseCase, a.getwd)
		},
	}
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
host are set in the shell; the others are absent.

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

--json, --dry-run, --diff and arguments are rejected with exit 2. --yes and
--no-interactive have no effect: the command asks for no confirmation.

Exit codes:
  0   the shell has begun and ended
  2   CLI validation (unknown flag, arguments, --json, --dry-run/--diff)
  10  no .devcontainer/devcontainer.json in the current directory; a directory
      path with a comma, double quote or line break; a ${localEnv:NAME}
      reference whose name is not a plain variable name
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
  u-boot --quiet devcontainer shell`,
		Args: jsonArgsValidator(a, "devcontainer", "shell", cobra.NoArgs, mapDevcontainerShellErrorToDiagnostic, false),
		RunE: func(cmd *cobra.Command, _ []string) error {
			flags.Quiet = a.quiet()
			flags.JSON = a.json
			return runDevcontainerShell(cmd.Context(),
				devcontainerStreams{in: cmd.InOrStdin(), stdout: cmd.OutOrStdout(), stderr: cmd.ErrOrStderr()},
				*flags, a.devcontainerUseCase, a.getwd, a.terminalProbe)
		},
	}
	registerDevcontainerPreviewFlags(cmd, flags)
	return cmd
}

func registerDevcontainerPreviewFlags(cmd *cobra.Command, flags *devcontainerFlags) {
	cmd.Flags().BoolVar(&flags.DryRun, "dry-run", false,
		"not applicable (the command writes no project file); rejected with exit 2")
	cmd.Flags().BoolVar(&flags.Diff, "diff", false,
		"not applicable (the command writes no project file); rejected with exit 2")
}

// reportDevcontainerError is the error gate of both subcommands. In JSON mode
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
		}
		diags := mapWarningsToDiagnostics(append(raised, resp.Warnings...))
		return writeEnvelope(stdout, newDataEnvelope("devcontainer", "up", data, diags, 0))
	}
	if !flags.Quiet {
		outcome := "started"
		if resp.AlreadyRunning {
			outcome = "already running"
		}
		fmt.Fprintf(stdout, "Devcontainer %s: container %s, remote user %s, workspace %s\n",
			outcome, shortContainerID(resp.ContainerID), resp.RemoteUser, resp.RemoteWorkspaceFolder)
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
	cwd, err := getwd()
	if err != nil {
		return fail(fmt.Errorf("determine working directory: %w", err))
	}
	resp, err := uc.Down(ctx, driving.DevcontainerDownRequest{BaseDir: cwd})
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
		ProgressSink: sink,
		OnWarning: func(w driving.WarningEntry) {
			writeWarnings(streams.stderr, visibleWarnings([]driving.WarningEntry{w}, flags.Quiet))
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
		fmt.Fprintf(out, "Warning (%s): %s\n", w.Code, w.Message)
	}
}

// shortContainerID is the first 12 characters of a container ID, the form
// `docker ps` shows.
func shortContainerID(id string) string {
	if len(id) > 12 {
		return id[:12]
	}
	return id
}

// mapDevcontainerEngineSentinel classifies the engine-side sentinels both
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
// its diagnostic: the engine sentinels, then the rest. `down` reads no
// project file and returns none of the validation sentinels.
func mapDevcontainerDownErrorToDiagnostic(err error) diagnosticItem {
	if code, matched := mapDevcontainerEngineSentinel(err); matched {
		return diagnosticItem{Level: "error", Code: code, Message: err.Error()}
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
// up`; [isValidationError] takes them as one condition.
func isDevcontainerValidationError(err error) bool {
	return errors.Is(err, driving.ErrDevcontainerConfigMissing) ||
		errors.Is(err, driving.ErrDevcontainerProjectFolderInvalid) ||
		errors.Is(err, driving.ErrDevcontainerEnvNameInvalid)
}
