// Package netprobe is the os/net-backed implementation of the
// `port/driven.NetProbe` interface (M6-T3). It performs short-lived
// TCP-reachability probes used by the [UpService] polling loop to
// confirm declared service ports answer on `localhost` after Compose
// reports a service as `healthy` / `running`.
//
// Layer rules (LH-FA-ARCH-003): driven adapter — imports
// `hexagon/port/driven` + stdlib only; the application layer is
// blocked from importing the `net` stdlib directly by the depguard
// rule `application-no-net` (`.golangci.yml`), so all TCP-probe
// traffic in the project funnels through here.
//
// Directory naming: the slice plan originally specified
// `internal/adapter/driven/net/`, but that would have placed a
// production package named `net` next to the stdlib `net`,
// requiring every caller that also imports stdlib `net` to alias
// one of them. The directory was renamed to `netprobe/` so
// `package netprobe` matches the path's last segment without
// clashing. The slice plan (`docs/plan/planning/in-progress/
// slice-m6-up-down.md`) carries the rename.
package netprobe

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// Probe is the production NetProbe adapter. Construct with [New].
// The struct is intentionally empty: each [DialTCP] call opens and
// closes its own socket; there is no per-Probe state to carry
// across invocations.
type Probe struct{}

// Static check: Probe satisfies the NetProbe port.
var _ driven.NetProbe = (*Probe)(nil)

// New returns a Probe ready to use.
func New() *Probe { return &Probe{} }

// DialTCP implements [driven.NetProbe]. It uses
// `(*net.Dialer).DialContext` with `Dialer.Timeout` set per the
// caller's `timeout` argument — the linter (`noctx`) rejects
// `net.DialTimeout`, and `DialContext` also gives free
// `ctx.Err()`-precedence: when `ctx` is cancelled mid-dial the
// returned error wraps `context.Canceled` /
// `context.DeadlineExceeded` directly, satisfying the port's
// contract without manual pre/post `ctx.Err()` checks.
//
// The connection is closed immediately on success — this is a
// reachability check, not a session. Errors from the close are
// swallowed (a successful dial already proves reachability; a
// failing close at this point would not change the outcome).
func (Probe) DialTCP(ctx context.Context, host string, port int, timeout time.Duration) error {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	_ = conn.Close()
	return nil
}
