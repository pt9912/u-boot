package docker

import (
	"io"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// Test-only bridges for the package-private helpers in engine.go.
// Convention: keep production internals unexported, expose narrow
// surfaces here so the `_test` test package can exercise them
// without polluting the public API.

// ParseComposePsOutputForTest exposes [parseComposePsOutput] for
// unit testing the NDJSON/array parser without going through the
// full Engine.ComposePs path (which needs a real docker binary).
func ParseComposePsOutputForTest(raw []byte) ([]driven.ComposeService, error) {
	return parseComposePsOutput(raw)
}

// ProgressSinkOrDiscardForTest exposes [progressSinkOrDiscard] so
// the nil-tolerance convention can be unit-tested.
func ProgressSinkOrDiscardForTest(sink io.Writer) io.Writer {
	return progressSinkOrDiscard(sink)
}
