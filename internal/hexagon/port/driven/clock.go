package driven

import "time"

// Clock abstracts time for the application layer. Tests inject a fake
// clock so deterministic time-dependent assertions (e.g. CHANGELOG
// date stamps, polling-loop iteration timing) are possible without
// freezing real wall-clock time.
//
// M6-T4 polling-loop usage: [UpService] calls [Clock.Now] at the
// start of the loop and after each iteration to detect timeout, and
// [Clock.Sleep] between iterations to throttle Compose-Ps polling.
// Tests replace both with a manual-advance fake so the loop runs at
// maximum speed without `time.Sleep` introducing flakiness.
type Clock interface {
	// Now returns the current instant.
	Now() time.Time

	// Sleep pauses the calling goroutine for at least `d`. A fake
	// implementation may advance internal time and return
	// immediately (the slice-plan-mandated "no real sleep in
	// tests" contract — see `slice-m6-up-down.md` §T4).
	Sleep(d time.Duration)
}
