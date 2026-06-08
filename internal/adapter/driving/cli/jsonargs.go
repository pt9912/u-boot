package cli

import (
	"github.com/spf13/cobra"
)

// jsonArgsValidator wraps a positional-args base validator with the
// LH-NFA-USE-004 §1841 envelope hook: when --json is active and the
// base validator rejects (wrong arg count), the error is emitted as a
// spec-valid envelope on stdout BEFORE Cobra returns its usage error.
// Without it a machine consumer gets a bare stderr message and no
// JSON on stdout — a §1841 violation.
//
// Cluster-Konsolidierung (slice-v1-cli-json-envelope-consolidation
// T1, SD-A (a)): this is the single shared form that config's
// configArgsValidator and remove's validateRemoveArgs delegate to,
// and that add/init/generate adopt — removing the R15-Cross-Slice-1
// pattern drift.
//
// previewFlags selects the §1842 schema policy (Schutzplanke 1):
//   - true for modifying forms (add/init/generate/remove/`config
//     set`): reads the actual --dry-run/--diff state so a wrong-arg
//     `--dry-run --json add` still emits the Voll-Schema.
//   - false for read-only forms (`config show`/`config get`) which
//     REGISTER --dry-run/--diff as reject-flags but must stay on the
//     Minimal-Envelope for arg errors. (Replaces config's earlier
//     `if subcommand == "set"` guard, now made explicit.)
//
// base is a per-command closure (not necessarily a raw cobra
// validator) so a command can surface a custom missing-arg sentinel
// — e.g. remove returns ErrServiceNameMissing for len(args)==0
// (Schutzplanke 2).
func jsonArgsValidator(
	a *App,
	command, subcommand string,
	base cobra.PositionalArgs,
	mapErr func(error) diagnosticItem,
	previewFlags bool,
) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		err := base(cmd, args)
		if err == nil {
			return nil
		}
		if a.json {
			dryRun, diffFlag := false, false
			if previewFlags {
				dryRun, _ = cmd.Flags().GetBool("dry-run")
				diffFlag, _ = cmd.Flags().GetBool("diff")
			}
			_ = writeErrorEnvelopeSub(cmd.OutOrStdout(), err, nil, dryRun, diffFlag, command, subcommand, mapErr, nil)
		}
		return err
	}
}
