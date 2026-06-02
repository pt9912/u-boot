package docker

import (
	"context"
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

// WrapComposeRunErrorForTest exposes [wrapComposeRunError] for the
// slice-v1-logs Review-Followup F3 unit-test of SIGINT-Pass-Through
// Schicht 1. The helper was extracted from ComposeLogs so the
// "ctx.Err() unverdeckt"-contract can be pinned without a real
// subprocess.
func WrapComposeRunErrorForTest(ctx context.Context, runErr error, kind string) error {
	return wrapComposeRunError(ctx, runErr, kind)
}
