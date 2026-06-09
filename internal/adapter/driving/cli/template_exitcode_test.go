package cli_test

import (
	"fmt"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// TestExitCode_TemplateInitSentinels pins the LH-FA-CLI-006 exit-code
// classes for the template-init sentinels. slice-later-local-templates
// T1 adds ErrTemplateInvalid to the exit-10 class
// (isTemplateInitValidationError); this is the dual-classifier pin
// guarding against a future ExitCode regression that would silently
// drop it to exit 1.
func TestExitCode_TemplateInitSentinels(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		err  error
		want int
	}{
		{"not found", driving.ErrTemplateNotFound, 10},
		{"invalid path", driving.ErrInvalidTemplatePath, 10},
		{"invalid metadata", driving.ErrTemplateInvalid, 10},
		{"render failure", driving.ErrTemplateRender, 14},
		// Through a multi-%w wrap, mirroring TemplateInitService.Init.
		{"wrapped invalid metadata", fmt.Errorf("%w: bad", driving.ErrTemplateInvalid), 10},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := cli.ExitCode(tc.err); got != tc.want {
				t.Errorf("ExitCode(%v) = %d, want %d", tc.err, got, tc.want)
			}
		})
	}
}
