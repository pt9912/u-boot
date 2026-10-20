package cli

import (
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
)

const (
	repoURL    = "https://github.com/pt9912/u-boot"
	manualPath = "docs/user/benutzerhandbuch.md"
)

// releaseVersionPattern captures the documented release forms `X.Y.Z` and
// `X.Y.Z-rc.N` (numbers without leading zeros). It is a subset of what the
// publish workflow accepts as a tag; any other version gets the `main` link
// only.
const releaseVersionPattern = `^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-rc\.(0|[1-9][0-9]*))?$`

// rootLongText renders the long help text of the root command. The manual
// link of the version (tag `v<version>`) appears only when version is a
// release form (`X.Y.Z` or `X.Y.Z-rc.N`); every other version (`-dev`,
// test or self-built, empty) gets the link to `main` only, because a tag
// of that form is not promised.
func rootLongText(version string) string {
	text := `u-boot bootstraps reproducible development environments:
project structure, Docker Compose stack, devcontainer configuration,
service add-ons, and the usual recurring artefacts (README, CHANGELOG,
.env.example).

For command details, run u-boot <command> --help.
Quickstart and user docs: ` + repoURL + `#readme
`
	if regexp.MustCompile(releaseVersionPattern).MatchString(version) {
		text += fmt.Sprintf("User manual (German), this version: %s/blob/v%s/%s\n", repoURL, version, manualPath)
	}
	return text + fmt.Sprintf("User manual (German), latest: %s/blob/main/%s", repoURL, manualPath)
}

// buildRootCommand assembles the root `u-boot` command with every
// subcommand registered. Kept in its own file (separate from
// per-subcommand files) so adding a new subcommand only needs a
// new file plus a single `AddCommand` line here.
func buildRootCommand(a *App) *cobra.Command {
	root := &cobra.Command{
		Use:     "u-boot",
		Short:   "Developer environment bootloader for Docker-based projects",
		Long:    rootLongText(a.version),
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
	// a single source of truth. The three flags share one value
	// ([verbosityFlag]): the option given last on the command line
	// decides, a repeated or conflicting one is not rejected.
	// --quiet shows errors only (LH-FA-CLI-005): every modifying
	// subcommand drops its success and summary output, doctor renders
	// only SeverityError items, degradation-fallback items and
	// security-relaxation items, and the
	// warning lines of init/generate/up keep only security relaxations
	// and degradation fallbacks, see [visibleWarnings]. The query
	// commands (config get/list, egress list, logs, template list) keep
	// their output. The selected verbosity also sets the logger level in the
	// PersistentPreRunE below.
	a.verbosity = verbosityDefault
	addVerbosityFlag(root, a, "quiet", verbosityQuiet,
		"show errors only; security relaxations, degradation fallbacks of the sandbox profile, `up` progress and the `--force` path summary stay reported; the last of --quiet, --verbose and --debug wins")
	addVerbosityFlag(root, a, "verbose", verbosityVerbose,
		"show additional detail; the last of --quiet, --verbose and --debug wins")
	addVerbosityFlag(root, a, "debug", verbosityDebug,
		"show internal diagnostic output; the last of --quiet, --verbose and --debug wins")

	// LH-NFA-USE-004 machine-readable output. Persistent flag,
	// inherited by every subcommand. All spec-enum subcommand forms
	// read a.json from their RunE and emit the envelope (the
	// migration completed with the Cluster-T_close; the transitional
	// allowlist reject gate does not exist). Two `--json` rejects remain: bare
	// `u-boot template` and bare `u-boot egress` (RunE-borne, see
	// template.go and egress.go). See
	// docs/maintainer/cli-json-contract.md §6.
	root.PersistentFlags().BoolVar(&a.json, "json", false,
		"emit machine-readable JSON output where supported")

	// LH-FA-CLI-005 verbosity → slog level wiring. Runs after Cobra
	// has parsed flags, before the subcommand's RunE. The LevelVar
	// instance is shared with the logger adapter (wired in
	// cmd/uboot/main.go), so the Set call here changes the level of
	// every Logger.Debug/Info/... call that follows.
	//
	// The level follows the verbosity option given last, see
	// [verbosity.slogLevel]; --debug and --verbose both map to
	// LevelDebug.
	root.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		// LH-FA-CLI-005 verbosity → slog level. The LH-NFA-USE-004
		// `--json` reject gate (allowlist + applyJSONRejectGate) was
		// REMOVED at Cluster-T_close: all spec-enum subcommand forms
		// carry `--json` now, so there is nothing left to reject. The
		// remaining `--json` rejects are bare `u-boot template` and
		// bare `u-boot egress`, handled RunE-borne via
		// [ErrTemplateSubcommandRequired] and [ErrEgressSubcommandRequired]
		// (slice-v1-cli-json-dry-run-template T0-(a)). See
		// docs/maintainer/cli-json-contract.md §6.1.
		if a.logLevel == nil {
			return nil
		}
		a.logLevel.Set(a.verbosity.slogLevel())
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
	root.AddCommand(newEgressCommand(a))
	root.AddCommand(newDevcontainerCommand(a))
	return root
}

// addVerbosityFlag registers one verbosity option as a value-less
// persistent boolean flag of root that selects mode when given.
func addVerbosityFlag(root *cobra.Command, a *App, name string, mode verbosity, usage string) {
	f := root.PersistentFlags().VarPF(&verbosityFlag{app: a, mode: mode}, name, "", usage)
	f.NoOptDefVal = "true"
}
