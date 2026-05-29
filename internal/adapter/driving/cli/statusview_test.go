package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// statusview tests use the export_test.go bridge to reach the
// unexported renderUpStatus / renderUpDiagnostics /
// renderDownSuccess helpers — keeps the helpers package-private
// in production while satisfying the testpackage convention.

func TestRenderUpStatus_EmptyServices_NoOutput(t *testing.T) {
	t.Parallel()
	// Fire-and-forget mode (--timeout=0) returns Services=nil; the
	// CLI must NOT print an empty header row. Pin so a future
	// renderer refactor doesn't regress LH-NFA-USE-004 stability.
	var buf bytes.Buffer
	if err := cli.RenderUpStatusForTest(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %q", buf.String())
	}
	if strings.Contains(buf.String(), "SERVICE") {
		t.Errorf("output contains header row in fire-and-forget mode: %q", buf.String())
	}
}

func TestRenderUpStatus_SingleServiceNoPortNoHealth_RendersDash(t *testing.T) {
	t.Parallel()
	services := []domain.ServiceStatus{
		{Name: "worker", ContainerStatus: domain.StateRunning},
	}
	var buf bytes.Buffer
	if err := cli.RenderUpStatusForTest(&buf, services); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "SERVICE") || !strings.Contains(out, "CONTAINER") ||
		!strings.Contains(out, "PORT") || !strings.Contains(out, "HEALTH") {
		t.Errorf("output missing header columns: %q", out)
	}
	if !strings.Contains(out, "worker") {
		t.Errorf("output missing service row: %q", out)
	}
	if !strings.Contains(out, "running") {
		t.Errorf("output missing container state: %q", out)
	}
	// Both port and health must be "-" for a bare running service.
	// Count "-" occurrences: should be ≥2.
	if strings.Count(out, " -") < 2 {
		t.Errorf("output missing '-' placeholders for empty port/health: %q", out)
	}
}

func TestRenderUpStatus_FullService_RendersAllFields(t *testing.T) {
	t.Parallel()
	services := []domain.ServiceStatus{
		{Name: "postgres", ContainerStatus: domain.StateRunning, Port: "5432:5432", Healthcheck: "healthy"},
	}
	var buf bytes.Buffer
	if err := cli.RenderUpStatusForTest(&buf, services); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"postgres", "running", "5432:5432", "healthy"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q: %s", want, out)
		}
	}
}

func TestRenderUpStatus_MultipleServices_OrderPreserved(t *testing.T) {
	t.Parallel()
	// The application service is responsible for sorting; this
	// function MUST preserve incoming order. Verify by feeding an
	// already-sorted slice and checking the lines come out in that
	// order.
	services := []domain.ServiceStatus{
		{Name: "alpha", ContainerStatus: domain.StateRunning, Port: "8080:80", Healthcheck: "healthy"},
		{Name: "postgres", ContainerStatus: domain.StateRunning, Port: "5432:5432", Healthcheck: "healthy"},
	}
	var buf bytes.Buffer
	if err := cli.RenderUpStatusForTest(&buf, services); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	alphaIdx := strings.Index(out, "alpha")
	postgresIdx := strings.Index(out, "postgres")
	if alphaIdx < 0 || postgresIdx < 0 {
		t.Fatalf("rows missing: %q", out)
	}
	if alphaIdx > postgresIdx {
		t.Errorf("rows out of order: alpha after postgres in %q", out)
	}
}

func TestRenderUpDiagnostics_QuietSuppressesAll(t *testing.T) {
	t.Parallel()
	diagnostics := []domain.Diagnostic{
		{ID: "up.port.x.0", Severity: domain.SeverityWarn, Message: "non-TCP port"},
		{ID: "up.fire-and-forget", Severity: domain.SeverityInfo, Message: "fire-and-forget"},
	}
	var buf bytes.Buffer
	cli.RenderUpDiagnosticsForTest(&buf, diagnostics, true)
	if buf.Len() != 0 {
		t.Errorf("--quiet should suppress diagnostics, got: %q", buf.String())
	}
}

func TestRenderUpDiagnostics_InfoAndWarnShown_OKFiltered(t *testing.T) {
	t.Parallel()
	diagnostics := []domain.Diagnostic{
		{ID: "filtered.ok", Severity: domain.SeverityOK, Message: "should-not-show"},
		{ID: "up.fire-and-forget", Severity: domain.SeverityInfo, Message: "fire-and-forget here", Hint: "run u-boot doctor"},
		{ID: "up.port.x.0", Severity: domain.SeverityWarn, Message: "udp warn here"},
	}
	var buf bytes.Buffer
	cli.RenderUpDiagnosticsForTest(&buf, diagnostics, false)
	out := buf.String()
	if strings.Contains(out, "should-not-show") {
		t.Errorf("OK diagnostic leaked into output: %q", out)
	}
	if !strings.Contains(out, "fire-and-forget here") {
		t.Errorf("Info diagnostic missing: %q", out)
	}
	if !strings.Contains(out, "udp warn here") {
		t.Errorf("Warn diagnostic missing: %q", out)
	}
	if !strings.Contains(out, "hint: run u-boot doctor") {
		t.Errorf("hint missing for diagnostic with Hint set: %q", out)
	}
}

func TestRenderDownSuccess_NoVolumes(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	cli.RenderDownSuccessForTest(&buf, false, false)
	if got := strings.TrimSpace(buf.String()); got != "environment stopped" {
		t.Errorf("got %q, want \"environment stopped\"", got)
	}
}

func TestRenderDownSuccess_WithVolumes(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	cli.RenderDownSuccessForTest(&buf, true, false)
	if got := strings.TrimSpace(buf.String()); got != "environment stopped, volumes removed" {
		t.Errorf("got %q, want \"environment stopped, volumes removed\"", got)
	}
}

func TestRenderDownSuccess_QuietProducesEmptyOutput(t *testing.T) {
	t.Parallel()
	// Asymmetric --quiet contract from M6 slice §T6: down --quiet
	// suppresses the one-line success message entirely; CI scripts
	// rely on empty stdout for success.
	for _, removed := range []bool{false, true} {
		var buf bytes.Buffer
		cli.RenderDownSuccessForTest(&buf, removed, true)
		if buf.Len() != 0 {
			t.Errorf("--quiet should suppress success (removed=%v), got: %q", removed, buf.String())
		}
	}
}
