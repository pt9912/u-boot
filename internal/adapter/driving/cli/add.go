package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// addFlags bundles the per-invocation flag state of
// `u-boot add <service>` (--yes / --no-interactive read through
// from the root command, plus the add-specific --with-deps from
// LH-FA-ADD-006). Kept as a struct so [runAdd] stays consistent
// with [runInit] and so new per-service flags (`--persistence`,
// `--exporter`, … for V1 keycloak / otel) plug in without a
// signature churn.
type addFlags struct {
	Yes           bool
	NoInteractive bool
	WithDeps      bool
}

// newAddCommand builds the `u-boot add <service>` Cobra subcommand.
//
// Positional argument layout: `add <service>` instead of fixed
// sub-subcommands (`add postgres`, `add keycloak`, …) so a typo or
// an unsupported service name reaches [domain.NewServiceName] and the
// in-app catalogue check. Either rejection then maps to exit-code 10
// via [ExitCode]; a fixed sub-subcommand structure would surface
// `unknown command` as exit-code 2, hiding the fachliche distinction
// between "invalid name" and "name not in catalogue".
//
// The persistent --yes / --no-interactive flags are read through
// from the root command; today they are no-op for the postgres MVP
// (no per-service interactive prompts), but the mutual-exclusion
// check still fires here so future add-ons inherit the same usage
// contract.
func newAddCommand(a *App) *cobra.Command {
	flags := &addFlags{}

	cmd := &cobra.Command{
		Use:   "add <service>",
		Short: "Add a service add-on (PostgreSQL, …) to the u-boot project",
		Long: `Add an integrated service add-on to the current u-boot project. The
positional <service> argument names the catalogue entry; today MVP
supports only postgres (LH-FA-ADD-002). Keycloak (LH-FA-ADD-003) and
OpenTelemetry (LH-FA-ADD-004) join in V1.

The command runs in an initialized project (u-boot.yaml present per
LH-FA-ADD-001) and is idempotent: calling it twice on the same active
service is a no-op. The state-machine (LH-FA-ADD-005) handles
re-activation, block rebuild after a partial cleanup, and abort with
a repair hint when the project state is inconsistent (orphan compose
block or user-managed entry without u-boot marker).

Output: a short summary lists the relative file paths the use case
mutated (u-boot.yaml, compose.yaml, .env.example). On a true no-op
(active service with all artefacts present) nothing is printed
beyond the leading "already active" line.

Examples:
  u-boot add postgres                 # first add: register + write
  u-boot add postgres                 # idempotent re-run: no-op
  u-boot add redis                    # exit 10 — not in catalogue
  u-boot add keycloak --with-deps     # auto-install missing deps`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.Yes = a.yes
			flags.NoInteractive = a.noInteractive
			return runAdd(cmd.Context(), cmd.OutOrStdout(), args, *flags, a.addServiceUseCase, a.getwd)
		},
	}

	cmd.Flags().BoolVar(&flags.WithDeps, "with-deps", false,
		"auto-install missing add-on dependencies (LH-FA-ADD-006) without prompting")

	return cmd
}

// runAdd is split from the Cobra closure for direct unit-testing
// (no Cobra command construction needed). ctx is taken as the first
// parameter explicitly so contextcheck can see the propagation.
//
// The function performs the LH-FA-CLI-005A mutual-exclusion check on
// --yes / --no-interactive (same shape as runInit), validates the
// positional service name via [domain.NewServiceName] (LH-FA-ADD-001
// invalid-name path), then delegates to the AddServiceUseCase. The
// use case owns the LH-FA-ADD-005 dispatch — runAdd only translates
// the response into a human-readable summary.
func runAdd(
	ctx context.Context,
	out io.Writer,
	args []string,
	flags addFlags,
	uc driving.AddServiceUseCase,
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

	resp, err := uc.Add(ctx, driving.AddServiceRequest{
		BaseDir:       cwd,
		ServiceName:   svcName,
		WithDeps:      flags.WithDeps,
		Yes:           flags.Yes,
		NoInteractive: flags.NoInteractive,
	})
	if err != nil {
		return err
	}

	printAddSummary(out, resp)
	return nil
}

// printAddSummary writes a short, deterministic summary of the add
// outcome. Three shapes:
//
//   - PriorState=Active and Changed=nil: prints
//     "Service <name> is already active; no changes."
//   - PriorState=Active and Changed!=nil: prints
//     "Repaired service <name> artefacts." plus the changed paths.
//   - Other PriorState transitions: prints
//     "Added service <name>." plus the changed paths.
//
// Order of Changed follows the AddServiceUseCase contract
// (u-boot.yaml → compose.yaml → .env.example), so the user sees a
// stable rollback hint without re-sorting.
func printAddSummary(out io.Writer, resp driving.AddServiceResponse) {
	name := resp.ServiceName.String()
	switch {
	case resp.PriorState == resp.State && len(resp.Changed) == 0:
		fmt.Fprintf(out, "Service %q is already active; no changes.\n", name)
		return
	case resp.PriorState == resp.State && len(resp.Changed) > 0:
		fmt.Fprintf(out, "Repaired service %q artefacts.\n\nChanged:\n", name)
	default:
		fmt.Fprintf(out, "Added service %q.\n\nChanged:\n", name)
	}
	for _, p := range resp.Changed {
		fmt.Fprintln(out, "  - "+p)
	}
}
