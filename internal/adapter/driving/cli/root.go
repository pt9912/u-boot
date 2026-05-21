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

	root.AddCommand(newInitCommand(a))
	return root
}
