package cli

import (
	"io"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// Test-only bridges for unexported helpers in statusview.go. The
// helpers stay package-private in production; tests in
// `cli_test` package access them through these wrappers. Pattern
// borrowed from `internal/hexagon/application/export_test.go`.

// RenderUpStatusForTest exposes [renderUpStatus] (M6-T6).
func RenderUpStatusForTest(out io.Writer, services []domain.ServiceStatus) error {
	return renderUpStatus(out, services)
}

// RenderUpDiagnosticsForTest exposes [renderUpDiagnostics] (M6-T6).
func RenderUpDiagnosticsForTest(out io.Writer, diagnostics []domain.Diagnostic, quiet bool) {
	renderUpDiagnostics(out, diagnostics, quiet)
}

// RenderDownSuccessForTest exposes [renderDownSuccess] (M6-T6).
func RenderDownSuccessForTest(out io.Writer, removedVolumes, quiet bool) {
	renderDownSuccess(out, removedVolumes, quiet)
}
