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
		{"quiet-lowers-to-warn", []string{"--quiet", "doctor"}, slog.LevelWarn},
		{"verbose-raises-to-debug", []string{"--verbose", "doctor"}, slog.LevelDebug},
		{"debug-raises-to-debug", []string{"--debug", "doctor"}, slog.LevelDebug},
		{"debug-takes-precedence-over-quiet", []string{"--debug", "--quiet", "doctor"}, slog.LevelDebug},
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
