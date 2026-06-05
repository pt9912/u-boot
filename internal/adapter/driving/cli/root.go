package cli

import (
	"log/slog"

	"github.com/spf13/cobra"
)

// buildRootCommand assembles the root `u-boot` command with every
// subcommand registered. Kept in its own file (separate from
// per-subcommand files) so adding a new subcommand only needs a
// new file plus a single `AddCommand` line here.
func buildRootCommand(a *App) *cobra.Command {
	root := &cobra.Command{
		Use:   "u-boot",
		Short: "Developer environment bootloader for Docker-based projects",
		Long: `u-boot bootstraps reproducible development environments:
project structure, Docker Compose stack, devcontainer configuration,
service add-ons, and the usual recurring artefacts (README, CHANGELOG,
.env.example).

For command details, run u-boot <command> --help.
Quickstart and user docs: https://github.com/pt9912/u-boot#readme`,
		Version: a.version,
		// Disable Cobra's auto-suggest on unknown commands so the
		// error message is plain (`unknown command "frobnicate"`)
		// and the LH-FA-CLI-006 exit-code mapping in ExitCode is
		// stable across Cobra upgrades.
		DisableSuggestions: true,
		SilenceUsage:       true,
		SilenceErrors:      true,
	}

	// LH-FA-CLI-005A persistent flags. They apply to every
	// subcommand that takes confirmation decisions (init today; add,
	// remove, config set, down --volumes in M5+). Living on the root
	// command also means they appear in `u-boot --help` once instead
	// of being duplicated per subcommand.
	root.PersistentFlags().BoolVar(&a.yes, "yes", false,
		"answer confirmation prompts automatically; exclusive with --no-interactive")
	root.PersistentFlags().BoolVar(&a.noInteractive, "no-interactive", false,
		"do not prompt; fail when a required confirmation cannot be resolved; exclusive with --yes")

	// LH-FA-CLI-005 verbosity flags. Persistent so subcommands read
	// a single source of truth. --quiet is load-bearing for the
	// doctor subcommand (filters SeverityOK items from the rendered
	// report) and the up/down subcommands (suppress status table /
	// success message). --quiet / --verbose / --debug also raise
	// or lower the logger level via the PersistentPreRunE below.
	root.PersistentFlags().BoolVar(&a.quiet, "quiet", false,
		"reduce normal output; keep warnings and errors where applicable")
	root.PersistentFlags().BoolVar(&a.verbose, "verbose", false,
		"show additional detail")
	root.PersistentFlags().BoolVar(&a.debug, "debug", false,
		"show internal diagnostic output")

	// LH-NFA-USE-004 machine-readable output (slice-v1-cli-json-dry-
	// run-doctor T3). Persistent flag, inherited by every subcommand.
	// Migrated subcommands read a.json from their RunE; non-migrated
	// subcommand forms are rejected by the PersistentPreRunE below
	// via ErrJSONNotImplemented (exit code 2). See
	// docs/user/cli-json-output.md §6 for the migration roadmap.
	root.PersistentFlags().BoolVar(&a.json, "json", false,
		"emit machine-readable JSON output where supported")

	// LH-FA-CLI-005 verbosity → slog level wiring. Runs after Cobra
	// has parsed flags, before the subcommand's RunE. The LevelVar
	// instance is shared with the logger adapter (wired in
	// cmd/uboot/main.go), so the Set call here changes the level of
	// every Logger.Debug/Info/... call that follows.
	//
	// Precedence: --debug > --verbose > --quiet > default(Info).
	// --debug and --verbose both map to LevelDebug today; a future
	// slice can introduce a Verbose-only level if a service-specific
	// pegel between Info and Debug becomes useful.
	root.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		// LH-NFA-USE-004 reject gate first: a non-migrated subcommand
		// form receiving --json must fail with exit code 2 BEFORE
		// any side effects (logger level switching, RunE invocation).
		// See slice-v1-cli-json-dry-run-doctor §T0-(g).
		if err := applyJSONRejectGate(cmd, a.json); err != nil {
			return err
		}

		if a.logLevel == nil {
			return nil
		}
		switch {
		case a.debug:
			a.logLevel.Set(slog.LevelDebug)
		case a.verbose:
			a.logLevel.Set(slog.LevelDebug)
		case a.quiet:
			a.logLevel.Set(slog.LevelWarn)
		default:
			a.logLevel.Set(slog.LevelInfo)
		}
		return nil
	}

	root.AddCommand(newInitCommand(a))
	root.AddCommand(newDoctorCommand(a))
	root.AddCommand(newAddCommand(a))
	root.AddCommand(newRemoveCommand(a))
	root.AddCommand(newUpCommand(a))
	root.AddCommand(newDownCommand(a))
	root.AddCommand(newLogsCommand(a))
	root.AddCommand(newGenerateCommand(a))
	root.AddCommand(newConfigCommand(a))
	root.AddCommand(newTemplateCommand(a))
	return root
}
