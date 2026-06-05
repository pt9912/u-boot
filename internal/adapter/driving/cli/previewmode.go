package cli

import (
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// previewModeFromFlags maps the --dry-run / --diff flag combination
// to the [driving.PreviewMode] enum per slice-v1-cli-json-dry-run-
// add T0-(b) Wahrheitstabelle:
//
//	--dry-run | --diff | mode              | production write?
//	-----------+--------+-------------------+------------------
//	no        | no    | PreviewNone       | yes (Normal-Mode)
//	yes       | no    | PreviewDryRun     | no  (Plan only)
//	no        | yes   | PreviewAndApply   | yes (Plan + Write)
//	yes       | yes   | PreviewDryRun     | no  (Diff preview)
//
// The --dry-run-wins rule on the (yes, yes) cell matches LH-FA-CLI-
// 007 (dry-run is a hard "no write") combined with LH-FA-CLI-008
// (--diff alone is preview-and-apply). With both flags set the user
// asked for a diff preview WITHOUT writing — that's PreviewDryRun.
//
// Originally lived in add.go as a private add-only helper; extracted
// here in slice-v1-cli-json-dry-run-init T1-B so init/generate/
// remove/config-set can consume the same Wahrheitstabelle without
// copy-paste. File follows the established pattern from
// jsonenvelope.go / statusview.go / jsonallowlist.go (pure cli-
// package helpers, no RunE binding).
func previewModeFromFlags(dryRun, diffFlag bool) driving.PreviewMode {
	if dryRun {
		return driving.PreviewDryRun
	}
	if diffFlag {
		return driving.PreviewAndApply
	}
	return driving.PreviewNone
}
