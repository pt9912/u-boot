package cli_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
)

// Verbosity-wiring pin tests (slice-followup-verbosity-wiring).
// Live in their own file to keep cli_test.go focused on the
// existing Execute-/ExitCode-pins; the verbosity surface is
// orthogonal and easier to revisit when its own file holds it.

func TestExecute_VerbosityFlag_AdjustsLoggerLevel(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		args []string
		want slog.Level
	}{
		{"no-flag-defaults-to-info", []string{"doctor"}, slog.LevelInfo},
		{"quiet-sets-error", []string{"--quiet", "doctor"}, slog.LevelError},
		{"verbose-raises-to-debug", []string{"--verbose", "doctor"}, slog.LevelDebug},
		{"debug-raises-to-debug", []string{"--debug", "doctor"}, slog.LevelDebug},
		// LH-FA-CLI-005: the option given last decides; the former
		// fixed order (--debug before --quiet) is gone.
		{"quiet-after-debug-wins", []string{"--debug", "--quiet", "doctor"}, slog.LevelError},
		{"debug-after-quiet-wins", []string{"--quiet", "--debug", "doctor"}, slog.LevelDebug},
		{"verbose-after-quiet-wins", []string{"--quiet", "--verbose", "doctor"}, slog.LevelDebug},
		{"quiet-after-verbose-wins", []string{"--verbose", "--quiet", "doctor"}, slog.LevelError},
		{"repeated-quiet-verbose-quiet", []string{"--quiet", "--verbose", "--quiet", "doctor"}, slog.LevelError},
		{"repeated-verbose-quiet-verbose", []string{"--verbose", "--quiet", "--verbose", "doctor"}, slog.LevelDebug},
		{"repeated-same-option", []string{"--quiet", "--quiet", "doctor"}, slog.LevelError},
		{"verbose-and-debug-stay-debug", []string{"--verbose", "--debug", "doctor"}, slog.LevelDebug},
		{"option-after-subcommand-counts", []string{"--quiet", "doctor", "--verbose"}, slog.LevelDebug},
		{"explicit-false-clears-own-option", []string{"--quiet", "--quiet=false", "doctor"}, slog.LevelInfo},
		{"explicit-false-of-other-option-keeps-current", []string{"--quiet", "--verbose=false", "doctor"}, slog.LevelError},
		{"explicit-true-selects", []string{"--verbose=true", "doctor"}, slog.LevelDebug},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			levelVar := new(slog.LevelVar)
			levelVar.Set(slog.LevelInfo)
			getwd := func() (string, error) { return "/tmp/proj", nil }
			var stdout, stderr bytes.Buffer
			if err := newApp(&fakeInitUseCase{}, cli.WithLogLevel(levelVar), cli.WithGetwd(getwd)).Execute(context.Background(), tc.args, &stdout, &stderr); err != nil {
				t.Fatalf("Execute %v: %v", tc.args, err)
			}
			if got := levelVar.Level(); got != tc.want {
				t.Errorf("LogLevel after %v = %v, want %v", tc.args, got, tc.want)
			}
		})
	}
}

func TestExecute_WithoutLogLevel_NoPanic(t *testing.T) {
	t.Parallel()
	// Nil-tolerance pin: WithLogLevel is optional. Apps constructed
	// without it must not panic when the persistent verbosity flags
	// fire PersistentPreRunE.
	getwd := func() (string, error) { return "/tmp/proj", nil }
	var stdout, stderr bytes.Buffer
	if err := newApp(&fakeInitUseCase{}, cli.WithGetwd(getwd)).Execute(context.Background(), []string{"--verbose", "doctor"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute --verbose doctor without WithLogLevel: %v", err)
	}
}
