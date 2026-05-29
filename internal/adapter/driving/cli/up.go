package cli

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"

	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// upFlags bundles the per-invocation flag state of `u-boot up`.
// TimeoutSec is the raw seconds value from `--timeout`; runUp
// converts it to `time.Duration` via `time.Duration(sec) * time.Second`
// — never as raw `time.Duration(sec)` (which would be nanoseconds).
type upFlags struct {
	TimeoutSec int
	Quiet      bool
}

// newUpCommand builds the `u-boot up` Cobra subcommand
// (LH-FA-UP-001..003).
//
// Local flag:
//
//	--timeout <sec>   maximum wall-clock duration to wait for every
//	                  declared service to stabilize (default 60).
//	                  `0` short-circuits to fire-and-forget (no
//	                  polling, no port/healthcheck probes — see
//	                  LH-FA-UP-001 §970). Negative values are rejected
//	                  with [ErrInvalidTimeout] (exit code 2).
//
// The persistent flags --quiet (suppresses the status table) /
// --verbose / --debug (LH-FA-CLI-005) are read from the App after
// Cobra parses them. The Compose progress stream (LH-NFA-PERF-002)
// always reaches stderr regardless of --quiet.
func newUpCommand(a *App) *cobra.Command {
	flags := &upFlags{}

	cmd := &cobra.Command{
		Use:   "up",
		Short: "Start the Compose environment and wait for every service to stabilize",
		Long: `Bring the Compose environment defined in compose.yaml up via
docker compose up -d, then poll docker compose ps every 500ms until
every declared service reaches healthy (when a healthcheck is defined)
or running (when no healthcheck). LH-FA-UP-001 §966-§969 stabilization.

Stabilization semantics per LH-FA-UP-001:
  - --timeout <sec>      maximum wait (default 60). Negative ⇒ exit 2.
  - --timeout=0          fire-and-forget. No polling, no probes;
                         status table omitted, info diagnostic shown
                         (LH-FA-UP-001 §970).

LH-FA-CLI-006 exit codes:
  - 10  no u-boot.yaml or compose.yaml in the current directory
  - 11  Docker daemon unreachable / compose plugin missing
  - 12  Compose runtime failure or stabilization timeout

LH-NFA-PERF-002 progress: pull/create/start/healthcheck phases stream
to stderr live (unaffected by --quiet).`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			flags.Quiet = a.quiet
			return runUp(cmd.Context(), cmd.OutOrStdout(), cmd.ErrOrStderr(), *flags, a.upUseCase, a.getwd)
		},
	}
	cmd.Flags().IntVar(&flags.TimeoutSec, "timeout", 60,
		"maximum seconds to wait for stabilization; 0 = fire-and-forget (LH-FA-UP-001 §963/§970)")
	return cmd
}

// runUp is split from the Cobra closure for testability with a fake
// use-case + fake getwd. Returns the CLI-level error; the wiring
// layer maps it to an exit code via [ExitCode].
func runUp(ctx context.Context, stdout, stderr io.Writer, flags upFlags, useCase driving.UpUseCase, getwd func() (string, error)) error {
	if flags.TimeoutSec < 0 {
		return ErrInvalidTimeout
	}
	cwd, err := getwd()
	if err != nil {
		return fmt.Errorf("determine working directory: %w", err)
	}
	// Explicit `* time.Second` conversion — slice plan T1
	// timeout-semantics pin. `time.Duration(sec)` alone would be
	// nanoseconds.
	timeout := time.Duration(flags.TimeoutSec) * time.Second
	resp, err := useCase.Up(ctx, driving.UpRequest{
		BaseDir:      cwd,
		Timeout:      timeout,
		ProgressSink: stderr,
	})
	if err != nil {
		return err
	}
	// LH-FA-CLI-005 + M6 slice §T6 binding contract: `up --quiet`
	// must suppress BOTH the status table AND the diagnostic
	// section so CI scripts can rely on empty stdout for success.
	// The Compose progress stream on stderr is NOT affected
	// (LH-NFA-PERF-002 requires phase visibility regardless).
	if flags.Quiet {
		return nil
	}
	if err := renderUpStatus(stdout, resp.Result.Services); err != nil {
		return fmt.Errorf("render status: %w", err)
	}
	renderUpDiagnostics(stdout, resp.Result.Diagnostics, flags.Quiet)
	return nil
}
