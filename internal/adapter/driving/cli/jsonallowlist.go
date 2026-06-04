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
//   - 2 Migrate: "u-boot doctor", "u-boot template list"
//     (Flag-Schnitt ohne Envelope-Migration — Carveout).
//   - 11 Reject (heute): init, add, remove, up, down, logs, generate,
//     config (bare), config get, config set, template (bare).
//
// Cluster-T_close: alle 13 in Allowlist ODER Mechanik komplett raus.
func jsonAllowlist() map[string]bool {
	return map[string]bool{
		"u-boot doctor":        true, // slice-v1-cli-json-dry-run-doctor (this slice)
		"u-boot template list": true, // slice-v1-cli-json-dry-run-doctor (Flag-Schnitt, Carveout)
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
// set, template list) all share the parent's slice suffix.
func jsonSliceSuffix(cmdPath string) string {
	const root = "u-boot "
	rest := strings.TrimPrefix(cmdPath, root)
	first := strings.SplitN(rest, " ", 2)[0]
	switch first {
	case "up", "down":
		return "up-down"
	}
	return first
}

// applyJSONRejectGate runs at PersistentPreRunE time: if --json is
// set and the running cmd's path is not in the allowlist, return
// the reject error. Otherwise no-op. Help and version commands are
// always allowed through (the user is asking Cobra, not the
// subcommand, for output).
func applyJSONRejectGate(cmd *cobra.Command, jsonFlag bool) error {
	if !jsonFlag {
		return nil
	}
	if cmd == nil {
		return nil
	}
	// Cobra invokes PersistentPreRunE on the leaf command. The
	// builtin "help" command (and Cobra's internal __complete) live
	// under the root but are user-facing escape hatches — letting
	// them through avoids breaking `u-boot doctor --help --json`-
	// style invocations.
	if cmd.Name() == "help" || cmd.Name() == "__complete" {
		return nil
	}
	path := cmd.CommandPath()
	if jsonAllowlist()[path] {
		return nil
	}
	return jsonRejectError(path)
}
