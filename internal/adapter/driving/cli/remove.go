package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// removeFlags bundles the per-invocation flag state of
// `u-boot remove <service>`. Purge is the local destructive opt-in
// (LH-FA-ADD-007 §"Volumes nur auf explizite Anforderung"); Yes /
// NoInteractive read through from the root command's persistent
// PersistentFlags so `u-boot --yes remove postgres --purge` and
// `u-boot remove postgres --purge --yes` behave identically.
type removeFlags struct {
	Purge         bool
	Yes           bool
	NoInteractive bool
}

// newRemoveCommand builds the `u-boot remove <service>` Cobra
// subcommand (LH-FA-ADD-007).
//
// Symmetric to `add <service>`: positional service name validated
// via [domain.NewServiceName] before reaching the use case; catalog
// rejection and state-machine mismatches both surface as exit code
// 10 via [ExitCode]. The local `--purge` flag opts into volume
// removal (LH-FA-CLI-005A §254-style confirmation gate); the
// application service handles the gate consistency with `down
// --volumes`.
func newRemoveCommand(a *App) *cobra.Command {
	flags := &removeFlags{}

	cmd := &cobra.Command{
		Use:   "remove <service>",
		Short: "Remove a service add-on from the u-boot project",
		Long: `Remove an integrated service add-on from the current u-boot project
(LH-FA-ADD-007 — mirror of "u-boot add"). The positional <service>
argument names the catalogue entry; today MVP supports only
postgres. Keycloak / OpenTelemetry will follow in their own V1
slices.

The command runs in an initialized project (u-boot.yaml present per
LH-FA-ADD-001) and is idempotent: removing an already-disabled
service is a no-op with a clear message. The state-machine handles
the inverse of add — strip the service.<name> managed block from
compose.yaml and .env.example, then set services.<name>.enabled to
false in u-boot.yaml.

--purge is the explicit destructive opt-in for volume removal
(LH-FA-ADD-007). The same LH-FA-CLI-005A §254 confirmation gate as
"u-boot down --volumes" applies: --no-interactive without --yes
exits 10 (ErrConfirmationRequired); interactive mode prompts with a
safe default-No. In v0.3.0 the gate's "approved" outcome does NOT
auto-remove volumes yet — the CLI surfaces the deferred work so the
user can clean up manually with "docker volume rm <name>".

Exit codes (LH-FA-CLI-006):
  0   success (state transition OR idempotent no-op)
  2   CLI / flag errors (unknown subcommand, missing positional,
      conflicting mode flags)
  10  fachlich: ErrServiceUnsupported (not in catalogue),
      ErrServiceUnregistered (never added), ErrServiceInconsistent
      (orphan block or missing entry), ErrProjectNotInitialized,
      ErrConfirmationRequired (purge gate refused)

Examples:
  u-boot remove postgres                    # state-transitioning remove
  u-boot remove postgres                    # idempotent re-run: no-op
  u-boot remove postgres --purge --yes      # opt into volume cleanup
  u-boot remove redis                       # exit 10 — not in catalogue`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.Yes = a.yes
			flags.NoInteractive = a.noInteractive
			return runRemove(cmd.Context(), cmd.OutOrStdout(), args, *flags, a.removeServiceUseCase, a.getwd)
		},
	}

	cmd.Flags().BoolVar(&flags.Purge, "purge", false,
		"also request volume removal for the service (LH-FA-ADD-007). Destructive: triggers the LH-FA-CLI-005A §254 confirmation gate (refuses in --no-interactive without --yes). v0.3.0 does NOT auto-remove volumes after approval — the summary points at `docker volume rm` for manual cleanup.")
	return cmd
}

// runRemove is split from the Cobra closure for direct unit-testing
// (no Cobra construction needed). Mirrors [runAdd]'s shape; the
// mutual-exclusion check on --yes / --no-interactive lives here for
// the same reason (CLI-level usage error, not a use case concern).
func runRemove(
	ctx context.Context,
	out io.Writer,
	args []string,
	flags removeFlags,
	uc driving.RemoveServiceUseCase,
	getwd func() (string, error),
) error {
	if flags.Yes && flags.NoInteractive {
		return ErrConflictingModeFlags
	}

	svcName, err := domain.NewServiceName(args[0])
	if err != nil {
		return err
	}

	cwd, err := getwd()
	if err != nil {
		return fmt.Errorf("determine working directory: %w", err)
	}

	resp, err := uc.Remove(ctx, driving.RemoveServiceRequest{
		BaseDir:       cwd,
		ServiceName:   svcName,
		Purge:         flags.Purge,
		Yes:           flags.Yes,
		NoInteractive: flags.NoInteractive,
	})
	if err != nil {
		return err
	}

	printRemoveSummary(out, resp, flags.Purge)
	return nil
}

// printRemoveSummary writes a short, deterministic summary of the
// remove outcome. Three shapes:
//
//   - Idempotent no-op (PriorState == State == Deactivated,
//     Changed=nil):
//     "Service <name> is already disabled; no changes."
//   - State transition (PriorState=Active|EnabledUnset,
//     State=Deactivated):
//     "Removed service <name>." + list of changed paths.
//
// When `--purge` was requested AND the gate let us through (a
// successful response, not [driving.ErrConfirmationRequired]), the
// summary appends the v0.3.0 manual-cleanup hint so the user knows
// the volume cleanup is still on them. T3-Decision: the application
// service handles the gate but does not auto-remove volumes yet.
func printRemoveSummary(out io.Writer, resp driving.RemoveServiceResponse, purge bool) {
	name := resp.ServiceName.String()

	if len(resp.Changed) == 0 {
		fmt.Fprintf(out, "Service %q is already disabled; no changes.\n", name)
	} else {
		fmt.Fprintf(out, "Removed service %q.\n\nChanged:\n", name)
		for _, p := range resp.Changed {
			fmt.Fprintln(out, "  - "+p)
		}
	}

	if purge && !resp.VolumesPurged {
		fmt.Fprintf(out,
			"\nNOTE: --purge was requested; volume removal is not auto-handled in v0.3.0.\n"+
				"      Remove the %s volumes manually with `docker volume rm <name>` once you have\n"+
				"      confirmed the data is no longer needed (`docker volume ls` to list candidates).\n",
			name)
	}
}
