package driven

import (
	"context"
	"time"
)

// NetProbe abstracts TCP-reachability probing used by the
// [UpService] polling loop (M6-T4) to verify declared service
// ports are reachable on `localhost` after Compose stabilized the
// container.
//
// Keeping this behind a driven port instead of using the `net`
// stdlib directly from the application layer enforces the
// LH-FA-ARCH-003 depguard rule `application-no-net` and lets
// tests inject a deterministic fake (T4 adds `fakeNetProbe` to
// `internal/hexagon/application/fakes_test.go` with a
// `{host:port → error}` lookup map).
//
// Layer rules: driven port; the production implementation lives
// in `internal/adapter/driven/netprobe/probe.go` and is the only
// place `net.DialTimeout` runs in the project.
type NetProbe interface {
	// DialTCP attempts a TCP connection to `host:port` within the
	// given `timeout`. Returns nil on a successful TCP handshake
	// (the connection is closed immediately — this is a
	// reachability check, not a session); a non-nil error means
	// the target is unreachable for one of these reasons:
	// timeout exceeded, connection refused, host unresolved, or
	// context-cancellation.
	//
	// Contract:
	//
	//   - `host` is a textual hostname or IP literal (e.g.
	//     "localhost", "127.0.0.1", "[::1]"). NOT a "host:port"
	//     string — the adapter formats `host` and `port` via
	//     `net.JoinHostPort` so IPv6 literals get bracketed
	//     correctly.
	//   - `port` is an integer in [0, 65535]; the adapter does not
	//     range-check (the application's port-parse helper in T4
	//     handles that earlier in the pipeline).
	//   - If `ctx.Err()` is non-nil at any point (before, during,
	//     or just after the dial), the context error is returned
	//     in preference to the underlying net error so a cancelled
	//     call surfaces its cancel reason rather than a generic
	//     timeout error.
	//   - The adapter retains no state across calls; each
	//     invocation opens a new connection.
	DialTCP(ctx context.Context, host string, port int, timeout time.Duration) error
}
