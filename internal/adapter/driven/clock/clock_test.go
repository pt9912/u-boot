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
