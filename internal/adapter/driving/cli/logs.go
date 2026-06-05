package cli

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// logsFlags bundles the per-invocation flag state of `u-boot logs`.
// Slice-v1-logs §T0-Outcomes pinned the surface: only `--follow`
// and `--tail` are exposed (T0-(d) Spec-treu — no
// `--no-log-prefix`/`--timestamps`).
type logsFlags struct {
	// Follow mirrors `docker compose logs --follow`. Blocks until
	// SIGINT cancels `cmd.Context()`; the application service
	// short-circuits the resulting context.Canceled to a nil error
	// so the CLI exits 0.
	Follow bool

	// Tail is the raw string from `--tail <n>`. Empty means "flag
	// not set" — the application service normalises that to
	// Compose's `"all"`. Negative or non-numeric inputs (except
	// the internal `"all"`, which is not a user-supplied value)
	// are rejected by [validateLogsTailFlag] before the use case
	// runs.
	Tail string
}

// ErrInvalidLogsTail is returned by `u-boot logs` when `--tail` is
// neither empty nor a non-negative integer string. Slice-v1-logs
// §T0-(c) + §AK: numeric ≥ 0 accepted, everything else rejected at
// the CLI Stage-1 with Exit-Code 2 — the value never reaches the
// application service. Lives in the cli package because the
// LH-FA-CLI-006 mapping to exit code 2 is a CLI concern.
var ErrInvalidLogsTail = errors.New("--tail must be a non-negative integer")

// newLogsCommand builds the `u-boot logs [service]` Cobra
// subcommand (LH-FA-UP-005). Slice-v1-logs §AK contract:
//
//   - Positional argument: optional single service name; validated
//     via [domain.NewServiceName] (regex-only — Compose checks
//     runtime existence and surfaces a Compose-Runtime-Error if the
//     service is unknown). Format failures map to Exit-Code 10 via
//     [isServiceValidationError].
//   - `--follow` (default false). Blocks on `cmd.Context()`;
//     SIGINT cancels and the application service returns nil so
//     the CLI exits 0.
//   - `--tail <n>` (default empty → Compose-Default "all" after
//     normalisation in the use case). Accepts non-negative integer
//     strings only; otherwise [ErrInvalidLogsTail] / Exit-Code 2.
//
// Output: Compose-Default (with service-prefix, without
// timestamps). The CLI writes Compose stdout/stderr through the
// application service's OutputSink to `cmd.OutOrStdout()`.
func newLogsCommand(a *App) *cobra.Command {
	flags := &logsFlags{}
	cmd := &cobra.Command{
		Use:   "logs [service]",
		Short: "Stream Compose logs of every service or one selected service",
		Long: `Stream Docker Compose logs for the project's services.

Without a positional argument, all services declared in compose.yaml
are streamed (Compose-Default — no u-boot.yaml filter). With a
positional argument, only that single service streams. Unknown
services at runtime map to LH-FA-CLI-006 exit code 12 via the
Compose runtime error path.

Flags:
  --follow         stream until Ctrl-C (LH-FA-UP-005); SIGINT exits 0.
  --tail <n>       show only the last n lines per service. Default
                   shows all lines (Compose-Default). Negative or
                   non-numeric inputs ⇒ exit 2.

LH-FA-CLI-006 exit codes:
  - 0   success (incl. --follow terminated by SIGINT)
  - 2   --tail with invalid value, or malformed CLI usage
  - 10  no u-boot.yaml / compose.yaml; or invalid service name
        (regex-only, per T0-(b))
  - 11  Docker daemon unreachable / compose plugin missing
  - 12  Compose runtime failure (unknown service at runtime, etc.)`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogs(cmd.Context(), cmd.OutOrStdout(), args, *flags, a.logsUseCase, a.getwd)
		},
	}
	cmd.Flags().BoolVar(&flags.Follow, "follow", false,
		"stream logs continuously until Ctrl-C (LH-FA-UP-005 §1038)")
	cmd.Flags().StringVar(&flags.Tail, "tail", "",
		"show only the last n lines per service (non-negative integer; default = all)")
	return cmd
}

// runLogs is split from the Cobra closure for testability. Parses
// the positional service argument and the --tail flag, then
// delegates to the LogsUseCase.
//
// Slice-v1-logs §AK pinned the validation order:
//
//  1. --tail Stage-1 validation (Exit-Code 2 on invalid).
//  2. Service-name format validation via [domain.NewServiceName]
//     (Exit-Code 10 via isServiceValidationError on regex failure).
//     Skipped when no positional argument is given.
//  3. Working-directory probe → BaseDir.
//  4. Use-Case dispatch.
func runLogs(
	ctx context.Context,
	out io.Writer,
	args []string,
	flags logsFlags,
	uc driving.LogsUseCase,
	getwd func() (string, error),
) error {
	if err := validateLogsTailFlag(flags.Tail); err != nil {
		return err
	}
	var service string
	if len(args) == 1 {
		svc, err := domain.NewServiceName(args[0])
		if err != nil {
			return err
		}
		service = svc.String()
	}
	cwd, err := getwd()
	if err != nil {
		return fmt.Errorf("determine working directory: %w", err)
	}
	_, err = uc.Logs(ctx, driving.LogsRequest{
		BaseDir:    cwd,
		Service:    service,
		Follow:     flags.Follow,
		Tail:       flags.Tail,
		OutputSink: out,
	})
	return err
}

// validateLogsTailFlag enforces T0-(c) at the CLI Stage-1: empty
// (flag not set) passes through; otherwise the value must parse as
// a non-negative integer. The internal `"all"` constant is NOT a
// valid user-supplied value — only the application service
// produces it via normaliseTail.
//
// Review-Followup F1: Compose-CLI users tend to type `--tail all`
// out of muscle memory; the special-case explains that the
// implicit default already streams all lines, so the user can drop
// the flag entirely.
//
// Review-Followup F8: validation rejects signs and whitespace
// deterministically. Slice-v1-logs T0-(c) deliberately sets no upper
// bound; Compose receives very large decimal strings and decides
// whether it can handle them.
func validateLogsTailFlag(raw string) error {
	if raw == "" {
		return nil
	}
	if raw == "all" {
		return fmt.Errorf(
			"%w: `--tail \"all\"` is the implicit default; omit the flag to stream all lines",
			ErrInvalidLogsTail)
	}
	if !isDecimalDigits(raw) {
		return fmt.Errorf("%w: got %q", ErrInvalidLogsTail, raw)
	}
	return nil
}

func isDecimalDigits(raw string) bool {
	for _, ch := range raw {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}
