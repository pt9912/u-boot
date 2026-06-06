package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// jsonAllowlist returns the set of Cobra `cmd.CommandPath()` strings
// for which `--json` is implemented. Drift-Anker für den
// slice-v1-cli-json-dry-run-Cluster (T0-(g) Option 2): pro Folge-
// Slice-Merge wandern die migrierten Forms in diese Map; Cluster-
// T_close entfernt die Allowlist-Mechanik komplett (PersistentPreRunE
// raus, Funktion raus).
//
// gochecknoglobals-konform via Funktions-Wrapper (Repo-Konvention,
// vgl. jsontestutil.DefaultAllowedCodes).
//
// Subcommand-Form-Inventar (13 Forms, slice-doctor §T0-(g)):
//   - 4 Migrate: "u-boot doctor", "u-boot template list", "u-boot add",
//     "u-boot init" (template list: Flag-Schnitt ohne Envelope-
//     Migration — Carveout; add: voll-schema via slice-v1-cli-json-
//     dry-run-add T4; init: voll-schema via slice-v1-cli-json-dry-
//     run-init T5).
//   - 8 Reject (heute): remove, up, down, logs, config (bare),
//     config get, config set, template (bare).
//
// Cluster-T_close: alle 13 in Allowlist ODER Mechanik komplett raus.
func jsonAllowlist() map[string]bool {
	return map[string]bool{
		"u-boot doctor":        true, // slice-v1-cli-json-dry-run-doctor
		"u-boot template list": true, // slice-v1-cli-json-dry-run-doctor (Flag-Schnitt, Carveout)
		"u-boot add":           true, // slice-v1-cli-json-dry-run-add T4
		"u-boot init":          true, // slice-v1-cli-json-dry-run-init T5
		"u-boot generate":      true, // slice-v1-cli-json-dry-run-generate T5
	}
}

// jsonRejectError wraps [ErrJSONNotImplemented] with a concrete
// LH-FA-CLI-006-class message including the rejected CommandPath
// and the follow-up-slice reference. Format-Vorgabe aus slice-doctor
// §Aufhebungsbedingung:
//
//	JSON-Ausgabe für '<cmd.CommandPath()>' ist noch nicht implementiert
//	(siehe slice-v1-cli-json-dry-run-<sub>).
func jsonRejectError(cmdPath string) error {
	return fmt.Errorf(
		"%w: JSON-Ausgabe für '%s' ist noch nicht implementiert (siehe slice-v1-cli-json-dry-run-%s)",
		ErrJSONNotImplemented, cmdPath, jsonSliceSuffix(cmdPath),
	)
}

// jsonSliceSuffix maps a CommandPath to its follow-up-slice suffix
// per Cluster-T0-(e)-Reihenfolge. Compound subcommands (config get/
// set, template list) all share the parent's slice suffix. Defaults
// to "unknown" on empty / unrecognised input so the rendered error
// message never carries a malformed slice reference.
func jsonSliceSuffix(cmdPath string) string {
	const root = "u-boot "
	if !strings.HasPrefix(cmdPath, root) {
		return "unknown"
	}
	rest := strings.TrimPrefix(cmdPath, root)
	first := strings.SplitN(rest, " ", 2)[0]
	if first == "" {
		return "unknown"
	}
	switch first {
	case "up", "down":
		return "up-down"
	}
	return first
}

// applyJSONRejectGate runs at PersistentPreRunE time: if --json is
// set and the running cmd's path is not in the allowlist, return
// the reject error. Otherwise no-op.
//
// Three escape hatches let read-only / introspection paths through
// unmodified:
//
//  1. cmd.Name() == "help" — the builtin Cobra help subcommand.
//  2. cmd.Name() == "__complete" — Cobra-internal shell-completion
//     dispatcher (Bash/Zsh/Fish). Undocumented but stable in
//     cobra v1.10.2 (see go.mod). A Cobra major upgrade must
//     re-verify this internal command name; the
//     TestRootJSON_AcceptsHelpFlag pin in jsonallowlist_test.go
//     catches the regression.
//  3. The --help flag is set on the running command. Cobra parses
//     --help into the same persistent flag inherited from the root;
//     the help path is read-only by definition, so --json on a
//     non-migrated subcommand combined with --help should print
//     help instead of rejecting (Review M6-Findings adressiert).
func applyJSONRejectGate(cmd *cobra.Command, jsonFlag bool) error {
	if !jsonFlag {
		return nil
	}
	if cmd == nil {
		return nil
	}
	if cmd.Name() == "help" || cmd.Name() == "__complete" {
		return nil
	}
	if helpRequested(cmd) {
		return nil
	}
	path := cmd.CommandPath()
	if jsonAllowlist()[path] {
		return nil
	}
	return jsonRejectError(path)
}

// helpRequested checks whether Cobra parsed --help on the running
// command. Cobra registers --help as a persistent flag on every
// subcommand; the value is accessible via cmd.Flag("help").
func helpRequested(cmd *cobra.Command) bool {
	flag := cmd.Flag("help")
	if flag == nil {
		return false
	}
	return flag.Value.String() == "true"
}
