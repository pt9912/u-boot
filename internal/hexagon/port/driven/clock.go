package driven

import "time"

// Clock abstracts time for the application layer. Tests inject a fake
// clock so deterministic time-dependent assertions (e.g. CHANGELOG
// date stamps) are possible without freezing real time.
type Clock interface {
	// Now returns the current instant.
	Now() time.Time
}
