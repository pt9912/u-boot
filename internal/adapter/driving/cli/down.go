package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// downFlags bundles the per-invocation flag state of `u-boot down`.
// The local flags Volumes and Yes are CLI-only; NoInteractive and
// Quiet are read through from the App's persistent root flags
// (LH-FA-CLI-005 / LH-FA-CLI-005A).
type downFlags struct {
	Volumes        bool
	Yes            bool
	NoInteractive  bool
	Quiet          bool
}

// newDownCommand builds the `u-boot down` Cobra subcommand
// (LH-FA-UP-004).
//
// Local flags:
//
//	--volumes  remove named Compose volumes alongside containers
//	           (LH-FA-UP-004 §1015 destructive op; default false).
//	--yes      auto-confirm the destructive --volumes prompt
//	           (LH-FA-CLI-005A §234 / §246).
//
// The persistent flags --no-interactive (LH-FA-CLI-005A §235 / §245)
// and --quiet (LH-FA-CLI-005) are read from the App after Cobra
// parses them.
//
// Mode-flag mutual exclusion: `--yes` AND `--no-interactive` set
// together returns [ErrConflictingModeFlags] (exit code 2,
// LH-FA-CLI-005A §235). Independent of whether --volumes is set —
// the §235 rule is a global CLI-validation, not a destructive-path
// concern.
func newDownCommand(a *App) *cobra.Command {
	flags := &downFlags{}

	cmd := &cobra.Command{
		Use:   "down",
		Short: "Stop the Compose environment and optionally remove its named volumes",
		Long: `Tear down the Compose environment via docker compose down.
With --volumes the named Compose volumes are removed as well (data
loss; LH-FA-UP-004 §1015).

Destructive confirmation gate (LH-FA-CLI-005A §254):
  - --yes                     auto-confirm
  - --no-interactive          fail-fast with exit 10 (no confirmation
                              possible without user input)
  - interactive (default)     "[y/N]" prompt; "n" / empty / EOF aborts

The mode-flag mutual exclusion (--yes + --no-interactive → exit 2)
is checked before any use-case logic runs.

LH-FA-CLI-006 exit codes:
  - 10  no u-boot.yaml / no compose.yaml / destructive confirmation refused
  - 11  Docker daemon unreachable / compose plugin missing
  - 12  Compose runtime failure

LH-NFA-PERF-002 progress: compose down phases stream to stderr live
(unaffected by --quiet).`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			flags.NoInteractive = a.noInteractive
			flags.Quiet = a.quiet
			// --yes is a local flag on this subcommand, NOT
			// read-through from the root --yes flag. The root
			// --yes/--no-interactive are general LH-FA-CLI-005A
			// confirmation hooks; the down --yes is specific to
			// the destructive volume-removal prompt.
			return runDown(cmd.Context(), cmd.OutOrStdout(), cmd.ErrOrStderr(), *flags, a.downUseCase, a.getwd)
		},
	}
	cmd.Flags().BoolVar(&flags.Volumes, "volumes", false,
		"also remove named Compose volumes (data loss; LH-FA-UP-004 §1015)")
	cmd.Flags().BoolVar(&flags.Yes, "yes", false,
		"auto-confirm the destructive --volumes prompt (LH-FA-CLI-005A)")
	return cmd
}

// runDown is split from the Cobra closure for testability.
func runDown(ctx context.Context, stdout, stderr io.Writer, flags downFlags, useCase driving.DownUseCase, getwd func() (string, error)) error {
	// LH-FA-CLI-005A §235 mode-flag exclusion. Note: this checks
	// the LOCAL down --yes flag against the PERSISTENT root
	// --no-interactive — different fields but same exclusivity
	// rule.
	if flags.Yes && flags.NoInteractive {
		return ErrConflictingModeFlags
	}
	cwd, err := getwd()
	if err != nil {
		return fmt.Errorf("determine working directory: %w", err)
	}
	resp, err := useCase.Down(ctx, driving.DownRequest{
		BaseDir:        cwd,
		RemoveVolumes:  flags.Volumes,
		AssumeYes:      flags.Yes,
		NonInteractive: flags.NoInteractive,
		ProgressSink:   stderr,
	})
	if err != nil {
		return err
	}
	renderDownSuccess(stdout, resp.RemovedVolumes, flags.Quiet)
	return nil
}
