package cli

import (
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// TestPreviewModeFromFlags pins the slice-v1-cli-json-dry-run-add
// T0-(b) Wahrheitstabelle: four flag combinations map to three modes.
// Anti-Drift: a future refactor that adds a fourth mode or flips the
// (yes, yes) cell to PreviewAndApply would fail this test before it
// reaches the CLI integration tests.
func TestPreviewModeFromFlags(t *testing.T) {
	cases := []struct {
		name     string
		dryRun   bool
		diffFlag bool
		want     driving.AddPreviewMode
	}{
		{name: "no flags → PreviewNone (normal write)", dryRun: false, diffFlag: false, want: driving.PreviewNone},
		{name: "--dry-run → PreviewDryRun (no write)", dryRun: true, diffFlag: false, want: driving.PreviewDryRun},
		{name: "--diff → PreviewAndApply (write + capture)", dryRun: false, diffFlag: true, want: driving.PreviewAndApply},
		{name: "--dry-run --diff → PreviewDryRun (preview wins)", dryRun: true, diffFlag: true, want: driving.PreviewDryRun},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := previewModeFromFlags(tc.dryRun, tc.diffFlag)
			if got != tc.want {
				t.Errorf("previewModeFromFlags(%v, %v) = %v, want %v", tc.dryRun, tc.diffFlag, got, tc.want)
			}
		})
	}
}
