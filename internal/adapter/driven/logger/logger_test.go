package logger_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driven/logger"
)

func TestNew_TextFormat_EmitsKeyEqValue(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	l := logger.New(&buf, logger.FormatText, slog.LevelInfo)

	l.Info("project ready", "name", "demo", "dirs", 3)

	got := buf.String()
	for _, want := range []string{`msg="project ready"`, `name=demo`, `dirs=3`} {
		if !strings.Contains(got, want) {
			t.Errorf("text output missing %q: %q", want, got)
		}
	}
}

func TestNew_JSONFormat_EmitsValidJSON(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	l := logger.New(&buf, logger.FormatJSON, slog.LevelInfo)

	l.Warn("disk near full", "free_mb", 128)

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, buf.String())
	}
	if entry["msg"] != "disk near full" {
		t.Errorf("msg = %v, want %q", entry["msg"], "disk near full")
	}
	// JSON unmarshal lands numeric values as float64.
	if entry["free_mb"] != float64(128) {
		t.Errorf("free_mb = %v, want 128", entry["free_mb"])
	}
	if entry["level"] != "WARN" {
		t.Errorf("level = %v, want WARN", entry["level"])
	}
}

func TestNew_LevelFilter_DropsBelowMinimum(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	l := logger.New(&buf, logger.FormatText, slog.LevelWarn)

	l.Debug("dropped")
	l.Info("also dropped")
	l.Warn("kept")

	got := buf.String()
	if strings.Contains(got, "dropped") || strings.Contains(got, "also dropped") {
		t.Errorf("level filter let through sub-warn entries: %q", got)
	}
	if !strings.Contains(got, "kept") {
		t.Errorf("level filter dropped a warn entry: %q", got)
	}
}

func TestNew_AllLevels_RoutedCorrectly(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	l := logger.New(&buf, logger.FormatJSON, slog.LevelDebug)

	l.Debug("d")
	l.Info("i")
	l.Warn("w")
	l.Error("e")

	// One JSON object per line.
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 4 {
		t.Fatalf("got %d lines, want 4: %q", len(lines), buf.String())
	}
	wantLevels := []string{"DEBUG", "INFO", "WARN", "ERROR"}
	for i, line := range lines {
		var entry map[string]any
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			t.Fatalf("line %d invalid JSON: %v\n%s", i, err, line)
		}
		if entry["level"] != wantLevels[i] {
			t.Errorf("line %d level = %v, want %v", i, entry["level"], wantLevels[i])
		}
	}
}

func TestNew_UnknownFormat_FallsBackToText(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	// Any Format value other than FormatJSON should yield text.
	l := logger.New(&buf, logger.Format(99), slog.LevelInfo)

	l.Info("fallback test", "k", "v")

	if !strings.Contains(buf.String(), "k=v") {
		t.Errorf("expected text fallback to emit `k=v`, got: %q", buf.String())
	}
}
