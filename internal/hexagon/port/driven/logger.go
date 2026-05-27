package driven

// Logger is the application's structured-logging side-channel. The
// interface mirrors `log/slog`'s variadic key-value form so the
// slog-backed adapter is a one-liner per method; consumers in the
// application layer use the port directly and never import slog.
//
// The four levels match LH-QA-004's expected default profile (debug,
// info, warn, error). Adapters MAY emit additional structure (source
// location, timestamps) but the port surface stays narrow so callers
// do not depend on an adapter-specific log envelope.
//
// Calls are best-effort and must not return errors — a logging
// failure is never reason to abort the use case (mirrors the
// [ProgressPort] contract). args are alternating key/value pairs
// per the slog convention:
//
//	log.Info("project initialized", "name", projectName, "dirs", 3)
//
// Odd-numbered or malformed key-value sequences are the adapter's
// problem to render gracefully; the port does not validate.
//
// Layer rules (LH-FA-ARCH-003): consumers in `internal/hexagon/
// application` reference this interface only; the slog-backed
// implementation lives in `internal/adapter/driven/logger` and is
// injected by the wiring layer (`cmd/uboot`). Resolves the
// `forbidigo.msg` carveout (slice-m4-logging-port).
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}
