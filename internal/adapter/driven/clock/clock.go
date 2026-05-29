// Package clock is the real-time implementation of the
// `port/driven.Clock` interface (LH-FA-ARCH-002). Tests substitute a
// fake to make time-dependent assertions deterministic.
package clock

import (
	"time"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// Clock is the production clock adapter; the zero value is usable.
type Clock struct{}

// Static check: Clock satisfies the Clock port.
var _ driven.Clock = (*Clock)(nil)

// New returns a ready-to-use Clock.
func New() *Clock { return &Clock{} }

// Now returns time.Now() in UTC. UTC is the project-wide default so
// generated artefacts (CHANGELOG date stamps, backup timestamps) are
// not affected by the host's local timezone.
func (Clock) Now() time.Time { return time.Now().UTC() }

// Sleep delegates to time.Sleep. A non-positive duration is a no-op
// per time.Sleep's contract, which is the right behaviour for the
// M6 polling loop (a zero or negative interval just retries
// immediately instead of panicking).
func (Clock) Sleep(d time.Duration) { time.Sleep(d) }
