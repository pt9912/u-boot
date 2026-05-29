package clock_test

import (
	"testing"
	"time"

	"github.com/pt9912/u-boot/internal/adapter/driven/clock"
)

func TestClock_NowReturnsUTC(t *testing.T) {
	now := clock.New().Now()
	if now.Location() != time.UTC {
		t.Fatalf("Now().Location() = %v, want UTC", now.Location())
	}
}

func TestClock_NowIsRecent(t *testing.T) {
	before := time.Now().UTC()
	got := clock.New().Now()
	after := time.Now().UTC()

	if got.Before(before) || got.After(after) {
		t.Fatalf("Now() = %v not within [%v, %v]", got, before, after)
	}
}

func TestClock_SleepBlocksAtLeastDuration(t *testing.T) {
	t.Parallel()
	// Why: pin that the production Sleep delegates to time.Sleep
	// (i.e. actually blocks). Short duration (5ms) keeps the test
	// fast; generous upper bound (1s) tolerates slow CI runners.
	// A regression to a no-op would surface as elapsed < d/2.
	const want = 5 * time.Millisecond
	start := time.Now()
	clock.New().Sleep(want)
	elapsed := time.Since(start)
	if elapsed < want/2 {
		t.Errorf("Sleep(%v) returned after %v (likely a no-op regression)", want, elapsed)
	}
	if elapsed > time.Second {
		t.Errorf("Sleep(%v) took %v (unexpectedly long)", want, elapsed)
	}
}

func TestClock_SleepZeroIsNoop(t *testing.T) {
	t.Parallel()
	// Why: pin time.Sleep's contract for non-positive durations.
	// The M6 polling loop may compute a negative remaining timeout
	// in edge cases; Sleep(<=0) must not panic.
	start := time.Now()
	clock.New().Sleep(0)
	clock.New().Sleep(-time.Second)
	if time.Since(start) > 100*time.Millisecond {
		t.Errorf("Sleep(0)/Sleep(-1s) blocked unexpectedly long")
	}
}
