//go:build docker

// LH-NFA-PERF-002 progress-stream pin (M6-docker-int Sub-T1).
//
// Pins the binding contract that the Compose stderr stream
// (`Pulling…` / `Creating…` / `Starting…` / `Healthchecking…`)
// reaches `opts.ProgressSink` **live** — i.e. the first chunk
// arrives at the sink **before** `ComposeUp` returns. A buffered
// or post-hoc flush implementation would reverse that ordering.
//
// Operational definition of "live" is the happens-before relation
// `events[0].recvAt < composeUpReturnedAt`. No absolute wall-clock
// thresholds → robust on both fast (cache-hit) and slow (cold-pull)
// runners.
//
// Fixture: a service whose `image:` points at a non-existent
// registry hostname. The pull fails fast, producing real Compose
// stderr (`unable to get image` / `manifest unknown`), and
// ComposeUp returns non-zero. The failure path is the harder pin:
// a buffering implementation that flushes only on successful
// completion would still pass an all-green-path test.

package docker_test

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pt9912/u-boot/internal/adapter/driven/docker"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// composeProgressStreamFixture pins an image at a registry hostname
// that cannot resolve. Compose emits stderr immediately ("unable to
// resolve…"/"no such host"), then exits non-zero — both halves of
// the LH-NFA-PERF-002 contract are exercised.
const composeProgressStreamFixture = `services:
  doomed:
    image: zzz-nonexistent-uboot-test-host.invalid/nope:notreal
`

func TestEngine_RealDocker_ProgressStream_LivePerHappensBefore(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skipf("docker not on PATH: %v", err)
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "compose.yaml"), []byte(composeProgressStreamFixture), 0o644); err != nil {
		t.Fatalf("write compose.yaml: %v", err)
	}

	type chunkEvent struct {
		chunk  []byte
		recvAt time.Time
	}

	r, w := io.Pipe()
	defer func() { _ = r.Close() }()

	var (
		eventsMu sync.Mutex
		events   []chunkEvent
	)
	done := make(chan struct{})
	go func() {
		defer close(done)
		buf := make([]byte, 4096)
		for {
			n, err := r.Read(buf)
			if n > 0 {
				chunk := make([]byte, n)
				copy(chunk, buf[:n])
				eventsMu.Lock()
				events = append(events, chunkEvent{chunk: chunk, recvAt: time.Now()})
				eventsMu.Unlock()
			}
			if err != nil {
				return
			}
		}
	}()

	engine := docker.NewEngine()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, upErr := engine.ComposeUp(ctx, dir, driven.ComposeUpOptions{
		Detach:       true,
		ProgressSink: w,
	})
	composeUpReturnedAt := time.Now()

	// Close pipe so the reader goroutine drains and exits.
	_ = w.Close()
	<-done

	// The doomed-image fixture must make ComposeUp fail — that's
	// what produces the stderr output we are pinning.
	if upErr == nil {
		t.Fatal("expected ComposeUp to fail on non-existent image; got nil error")
	}

	eventsMu.Lock()
	defer eventsMu.Unlock()

	// (a) Some output was received at all.
	if len(events) == 0 {
		t.Fatal("no chunks received from ProgressSink — Compose stderr was not forwarded")
	}

	// (b) The first chunk arrived BEFORE ComposeUp returned.
	// This is the load-bearing pin: a buffered implementation
	// would emit all events at flush time, which is after the
	// call has already returned.
	if !events[0].recvAt.Before(composeUpReturnedAt) {
		t.Errorf(
			"first chunk arrived at %v, after ComposeUp returned at %v (Δ = %v); "+
				"this means stderr was buffered, not live — LH-NFA-PERF-002 violation",
			events[0].recvAt, composeUpReturnedAt, events[0].recvAt.Sub(composeUpReturnedAt),
		)
	}

	// Soft sanity: the output should mention something pull-related.
	// Logged-only, not asserted, because some Compose versions
	// truncate the stderr on quick failures.
	var combined strings.Builder
	for _, e := range events {
		combined.Write(e.chunk)
	}
	lower := strings.ToLower(combined.String())
	if !strings.Contains(lower, "pull") && !strings.Contains(lower, "manifest") &&
		!strings.Contains(lower, "resolve") && !strings.Contains(lower, "error") {
		t.Logf("compose stderr did not contain expected phase keywords; got: %q", combined.String())
	}
}
