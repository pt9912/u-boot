// Package logger is the slog-backed implementation of
// [driven.Logger] (LH-QA-004 logging port). Renders structured logs
// to a configurable [io.Writer] in text or JSON form, controlled by
// the build-time format / level pair.
//
// Layer rules (LH-FA-ARCH-003): driven adapter — imports
// `hexagon/port/driven` (and `log/slog`); must not import
// `hexagon/application` or `adapter/driving`.
package logger

import (
	"io"
	"log/slog"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// Format selects the on-wire shape of the log entries.
type Format int

const (
	// FormatText emits the slog text handler's `key=value` lines.
	// Default human-readable form for terminal output.
	FormatText Format = iota
	// FormatJSON emits one JSON object per entry — preferred for
	// `--json` runs (LH-NFA-USE-004) and CI log scrapers.
	FormatJSON
)

// slogLogger wraps a *slog.Logger and exposes the [driven.Logger]
// surface. The wrapper is the bridge between the slog API
// (consumers in `application` would otherwise import `log/slog`)
// and the project's own port.
type slogLogger struct {
	l *slog.Logger
}

// Static check: slogLogger satisfies the port.
var _ driven.Logger = (*slogLogger)(nil)

// New returns a [driven.Logger] backed by `log/slog`. `out` is the
// sink (typically `os.Stderr`); `format` selects text vs JSON;
// `level` is the minimum slog level emitted — entries below it are
// dropped at the handler level (no string formatting cost).
func New(out io.Writer, format Format, level slog.Level) driven.Logger {
	opts := &slog.HandlerOptions{Level: level}
	var handler slog.Handler
	switch format {
	case FormatJSON:
		handler = slog.NewJSONHandler(out, opts)
	default:
		handler = slog.NewTextHandler(out, opts)
	}
	return &slogLogger{l: slog.New(handler)}
}

func (s *slogLogger) Debug(msg string, args ...any) { s.l.Debug(msg, args...) }
func (s *slogLogger) Info(msg string, args ...any)  { s.l.Info(msg, args...) }
func (s *slogLogger) Warn(msg string, args ...any)  { s.l.Warn(msg, args...) }
func (s *slogLogger) Error(msg string, args ...any) { s.l.Error(msg, args...) }
