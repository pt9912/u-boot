package cli

import "github.com/spf13/cobra"

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
service add-ons (PostgreSQL, Keycloak, OpenTelemetry, …), and the
usual recurring artefacts (README, CHANGELOG, .env.example).

See spec/lastenheft.md for the full functional specification.`,
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
		"answer the default to every confirmation (LH-FA-CLI-005A); exclusive with --no-interactive")
	root.PersistentFlags().BoolVar(&a.noInteractive, "no-interactive", false,
		"abort with exit-code 2 on any required confirmation (LH-FA-CLI-005A); exclusive with --yes")

	root.AddCommand(newInitCommand(a))
	return root
}
