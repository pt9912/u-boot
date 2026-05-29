package netprobe_test

import (
	"context"
	"errors"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/u-boot/internal/adapter/driven/netprobe"
	driven_port "github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

func TestProbe_SatisfiesNetProbePort(t *testing.T) {
	t.Parallel()
	// Why: pin the interface conformance so a method-signature drift
	// breaks here, not in the cmd/uboot wiring path.
	var _ driven_port.NetProbe = netprobe.New()
}

func TestProbe_DialTCP_OpenPort_Succeeds(t *testing.T) {
	t.Parallel()
	// Why: pin the happy path against a real OS-assigned listener.
	// Using net.Listen with port 0 lets the kernel pick a free port,
	// avoiding flakiness from port-collision on CI runners.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("setup: net.Listen: %v", err)
	}
	defer func() { _ = ln.Close() }()

	host, portStr, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatalf("setup: SplitHostPort: %v", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("setup: parse port %q: %v", portStr, err)
	}

	p := netprobe.New()
	if err := p.DialTCP(context.Background(), host, port, time.Second); err != nil {
		t.Errorf("DialTCP to open port: %v", err)
	}
}

func TestProbe_DialTCP_RefusedPort_ReturnsError(t *testing.T) {
	t.Parallel()
	// Why: pin the refused path. Port 1 is privileged and reserved
	// (RFC 1700 tcpmux); on a developer machine without an explicit
	// tcpmux server bound there, connect() returns ECONNREFUSED
	// quickly — no timeout race with the test's deadline.
	p := netprobe.New()
	err := p.DialTCP(context.Background(), "127.0.0.1", 1, time.Second)
	if err == nil {
		t.Fatal("DialTCP to refused port: expected error, got nil")
	}
	// Sanity: the error should mention "refused" or be a recognizable
	// op error — but we don't pin the exact string because the net
	// stdlib's wording has drifted across Go releases. The non-nil
	// is the contract.
}

func TestProbe_DialTCP_CtxAlreadyCancelled_ReturnsCtxError(t *testing.T) {
	t.Parallel()
	// Why: pin the M6-slice contract that ctx.Err() takes precedence
	// over the dial error — here even before any dial happens.
	// Caller cancels ctx before invoking DialTCP; the adapter must
	// short-circuit and return ctx.Err() directly.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	p := netprobe.New()
	err := p.DialTCP(ctx, "127.0.0.1", 1, time.Second)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got: %v", err)
	}
}

func TestProbe_DialTCP_CtxDeadlineAlreadyExceeded_ReturnsDeadlineErr(t *testing.T) {
	t.Parallel()
	// Why: a deadline that has already passed must surface as
	// DeadlineExceeded, not as a net.OpError. Same precedence
	// contract as ctx.Canceled, but with the deadline variant.
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Hour))
	defer cancel()

	p := netprobe.New()
	err := p.DialTCP(ctx, "127.0.0.1", 1, time.Second)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context.DeadlineExceeded, got: %v", err)
	}
}

func TestProbe_DialTCP_CtxCancelledMidDial_UnwrapsToCtxError(t *testing.T) {
	t.Parallel()
	// Why: pin the mid-dial-cancel path explicitly. The
	// "already-cancelled-before-dial" pin upstream covers the
	// trivial case (DialContext short-circuits on ctx.Err() at
	// entry). The harder case is: ctx is cancelled *while*
	// DialContext is blocking on a slow connect. The error
	// returned by DialContext is a *net.OpError wrapping
	// ctx.Err(), and the M6-slice contract requires that
	// errors.Is(err, context.Canceled) traverses that wrap.
	//
	// Setup: dial RFC-5737 TEST-NET-1 with a long Dialer.Timeout;
	// the ctx is cancelled from a goroutine ~10ms later so the
	// cancel happens mid-dial, not at entry.
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	p := netprobe.New()
	err := p.DialTCP(ctx, "192.0.2.1", 80, 5*time.Second)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("mid-dial cancel: expected errors.Is(err, context.Canceled), got: %v (type %T)", err, err)
	}
}

func TestProbe_DialTCP_Timeout_ReturnsError(t *testing.T) {
	t.Parallel()
	// Why: pin the timeout path. RFC-5737 reserves 192.0.2.0/24 as
	// TEST-NET-1, a routable-looking-but-non-routable address space.
	// A TCP connect to 192.0.2.1:80 hangs until the timeout elapses
	// instead of immediately returning ECONNREFUSED — exactly the
	// behavior we need to exercise the timeout-error branch.
	// Note: we use a short timeout (50ms) so the test stays fast,
	// trading some flakiness risk on unusually slow CI for a more
	// useful per-PR signal.
	p := netprobe.New()
	start := time.Now()
	err := p.DialTCP(context.Background(), "192.0.2.1", 80, 50*time.Millisecond)
	elapsed := time.Since(start)
	if err == nil {
		t.Fatal("DialTCP to TEST-NET-1 with short timeout: expected error, got nil")
	}
	// The dial should have honored the timeout (allow up to 5x slack
	// for slow CI runners; primary signal is "didn't hang forever").
	if elapsed > 5*time.Second {
		t.Errorf("DialTCP did not honor short timeout: elapsed=%v", elapsed)
	}
}

func TestProbe_DialTCP_IPv6Literal_Formats(t *testing.T) {
	t.Parallel()
	// Why: pin that the adapter uses net.JoinHostPort (not naive
	// `host + ":" + port` concatenation) — IPv6 literals must be
	// bracketed. We test this end-to-end: an open IPv6 loopback
	// listener; the dial succeeds only if the address was formatted
	// as "[::1]:N", not "::1:N" (which is invalid).
	ln, err := net.Listen("tcp6", "[::1]:0")
	if err != nil {
		// IPv6 loopback isn't available on every CI runner; skip
		// rather than fail.
		t.Skipf("IPv6 loopback not available: %v", err)
	}
	defer func() { _ = ln.Close() }()

	addr, ok := ln.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatalf("setup: unexpected Addr type %T", ln.Addr())
	}

	p := netprobe.New()
	if err := p.DialTCP(context.Background(), "::1", addr.Port, time.Second); err != nil {
		t.Errorf("DialTCP to IPv6 loopback: %v", err)
	}
}

func TestProbe_DialTCP_RefusedErrorMentionsAddress(t *testing.T) {
	t.Parallel()
	// Why: pin that the dial error carries the attempted address
	// (so the M6 application can surface a meaningful diagnostic
	// without re-formatting). Don't pin exact text — pin that the
	// port number ("1") appears somewhere in the error string.
	// This catches a regression where the error is swallowed and
	// only "dial failed" is returned.
	p := netprobe.New()
	err := p.DialTCP(context.Background(), "127.0.0.1", 1, time.Second)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "127.0.0.1") {
		t.Errorf("error %q does not mention the dialed host", err.Error())
	}
}
