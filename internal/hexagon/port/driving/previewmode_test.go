package driving_test

import (
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// Compile-time identity checks for the slice-v1-cli-json-dry-run-
// init T0-(c) Carveout: `type AddPreviewMode = PreviewMode` (with
// `=` — type-alias, NOT a new defined type). The alias makes both
// type-names refer to the same underlying type, so existing call-
// sites that say `driving.AddPreviewMode` keep compiling.
//
// Anti-Drift: if a future maintainer accidentally rewrites the alias
// as `type AddPreviewMode PreviewMode` (defined type, no `=`), the
// var-declarations below stop compiling — the package fails to
// build before any test runs.

// Value-level identity: direct assignment between the two type-
// names compiles only when they are the same type.
var _ driving.AddPreviewMode = driving.PreviewMode(driving.PreviewNone)
var _ driving.PreviewMode = driving.AddPreviewMode(driving.PreviewNone)

// Function-type identity: a function value declared with one
// name must be assignable to the other without conversion.
// addFSFactory in cmd/uboot/main.go uses
// func(driving.PreviewMode); the alias must keep callers that
// declared their closures as func(driving.AddPreviewMode)
// assignable to the same callee without source edits.
var _ func(driving.AddPreviewMode) int = (func(driving.PreviewMode) int)(nil)
var _ func(driving.PreviewMode) int = (func(driving.AddPreviewMode) int)(nil)

// TestPreviewMode_ConstantsStillUsable is the runtime smoke for the
// renamed constants. The constants share the canonical PreviewMode
// type after T1-A's rename; alias keeps the AddPreviewMode references
// in caller code compiling.
func TestPreviewMode_ConstantsStillUsable(t *testing.T) {
	t.Parallel()
	for _, m := range []driving.PreviewMode{driving.PreviewNone, driving.PreviewDryRun, driving.PreviewAndApply} {
		if m < 0 {
			t.Errorf("PreviewMode constant out of range: %d", m)
		}
	}
}
