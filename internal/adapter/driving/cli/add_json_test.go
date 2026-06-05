package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/adapter/driving/cli/jsontestutil"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// addUseCaseStub returns a fixed response/err and ignores the request.
// Wraps the test-package fakeAddServiceUseCase pattern with a
// shape tailored for the JSON-pin tests below — every test fills in
// the resp.PlannedFiles / err directly.
type addUseCaseStub struct {
	resp driving.AddServiceResponse
	err  error
}

func (s *addUseCaseStub) Add(_ context.Context, _ driving.AddServiceRequest) (driving.AddServiceResponse, error) {
	return s.resp, s.err
}

// newAppWithAddStub wires the add use-case stub plus a deterministic
// getwd stub so the tests do not depend on the runner's CWD (review
// #13). cli_test.go's newApp helpers don't isolate getwd by default;
// this wrapper closes that gap for every test in this file and
// add_acceptance_test.go.
func newAppWithAddStub(stub driving.AddServiceUseCase) *cli.App {
	return newAppWithAdd(stub, cli.WithGetwd(func() (string, error) { return "/tmp/u-boot-add-test/demo", nil }))
}

func newAddSvcName(t *testing.T, raw string) domain.ServiceName {
	t.Helper()
	name, err := domain.NewServiceName(raw)
	if err != nil {
		t.Fatalf("NewServiceName(%q): %v", raw, err)
	}
	return name
}

// TestAddJSON_BareUsesMinimalEnvelope pins T0-(k): `u-boot add ...
// --json` ohne --dry-run/--diff trägt nur den Minimalkontrakt
// (Spec §1841) — keine plannedFiles/changes/dryRun/diff Felder.
func TestAddJSON_BareUsesMinimalEnvelope(t *testing.T) {
	stub := &addUseCaseStub{
		resp: driving.AddServiceResponse{
			ServiceName: newAddSvcName(t, "postgres"),
			Changed:     []string{"u-boot.yaml", "compose.yaml", ".env.example"},
		},
	}
	app := newAppWithAddStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "add", "postgres"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("add"),
		jsontestutil.WithExitCode(0),
	)
}

// TestAddJSON_DryRunUsesFullEnvelope pins the --dry-run --json path:
// voll-schema with dryRun=true, diff=false, plannedFiles from the
// recorder mapping.
func TestAddJSON_DryRunUsesFullEnvelope(t *testing.T) {
	stub := &addUseCaseStub{
		resp: driving.AddServiceResponse{
			ServiceName: newAddSvcName(t, "postgres"),
			PlannedFiles: []driving.PlannedFile{
				{Path: "u-boot.yaml", Action: "modify", OldContent: []byte("services:\n"), NewContent: []byte("services:\n  postgres: {}\n")},
				{Path: "compose.yaml", Action: "create", NewContent: []byte("services:\n  postgres:\n    image: postgres:16\n")},
				{Path: ".env.example", Action: "create", NewContent: []byte("POSTGRES_DB=app\n")},
			},
		},
	}
	app := newAppWithAddStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"add", "postgres", "--dry-run", "--json"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertFullEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("add"),
		jsontestutil.WithExitCode(0),
	)
	// Pin dryRun=true, diff=false, plannedFiles[].action shape.
	var env map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}
	if got, _ := env["dryRun"].(bool); !got {
		t.Errorf("dryRun: want true, envelope=%s", stdout.String())
	}
	if got, _ := env["diff"].(bool); got {
		t.Errorf("diff: want false, envelope=%s", stdout.String())
	}
	pfs, _ := env["plannedFiles"].([]any)
	if len(pfs) != 3 {
		t.Errorf("plannedFiles: want 3, got %d (envelope=%s)", len(pfs), stdout.String())
	}
}

// TestAddJSON_DiffWithoutDryRunRendersHunks pins the --diff --json
// preview-and-apply path: voll-schema with diff=true, dryRun=false,
// and plannedFiles[].hunks populated for non-binary files.
func TestAddJSON_DiffWithoutDryRunRendersHunks(t *testing.T) {
	stub := &addUseCaseStub{
		resp: driving.AddServiceResponse{
			ServiceName: newAddSvcName(t, "postgres"),
			PlannedFiles: []driving.PlannedFile{
				{Path: "compose.yaml", Action: "create", NewContent: []byte("services:\n  postgres:\n    image: postgres:16\n")},
			},
		},
	}
	app := newAppWithAddStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"add", "postgres", "--diff", "--json"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertFullEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("add"),
		jsontestutil.WithExitCode(0),
	)
	var env map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}
	if got, _ := env["diff"].(bool); !got {
		t.Errorf("diff: want true, envelope=%s", stdout.String())
	}
	if got, _ := env["dryRun"].(bool); got {
		t.Errorf("dryRun: want false, envelope=%s", stdout.String())
	}
	pfs, _ := env["plannedFiles"].([]any)
	if len(pfs) != 1 {
		t.Fatalf("plannedFiles: want 1, got %d (envelope=%s)", len(pfs), stdout.String())
	}
	first, _ := pfs[0].(map[string]any)
	hunks, ok := first["hunks"].([]any)
	if !ok || len(hunks) == 0 {
		t.Errorf("plannedFiles[0].hunks: want non-empty, got %v (envelope=%s)", first["hunks"], stdout.String())
	}
}

// TestAddJSON_ValidationErrorShipsMinimalEnvelope pins the pre-write
// validation-error path (T0-(b) Scenario 3): an invalid service name
// fails before the use case runs; PlannedFiles stays empty and the
// envelope ships the minimal shape with exitCode 10 / LH-FA-INIT-006.
func TestAddJSON_ValidationErrorShipsMinimalEnvelope(t *testing.T) {
	stub := &addUseCaseStub{} // unused — domain.NewServiceName fails first
	app := newAppWithAddStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "add", "INVALID NAME WITH SPACES"}, &stdout, &stderr)
	// Review #2: validation errors propagate so cli.ExitCode maps to
	// 10 (LH-FA-INIT-006). The envelope on stdout carries the same
	// LH-Code; shell exit and envelope MUST agree.
	if err == nil {
		t.Fatal("expected error to propagate (review #2)")
	}
	if got := cli.ExitCode(err); got != 10 {
		t.Errorf("cli.ExitCode: want 10 (LH-FA-INIT-006), got %d (err=%v)", got, err)
	}
	if stdout.Len() == 0 {
		t.Fatalf("expected JSON envelope on stdout, got empty (err=%v)", err)
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("add"),
		jsontestutil.WithExpectedCodes("LH-FA-INIT-006"),
		jsontestutil.WithExitCode(10),
	)
}

// TestAddJSON_DryRunDiffCombo pins the canonical Plan §Aufhebungsbedingung
// `--dry-run --diff --json` 3-flag combo (review #12): voll-schema with
// BOTH dryRun=true AND diff=true, plannedFiles[].hunks populated, no
// production write. Anti-Drift: a refactor that flipped the (true,true)
// cell of previewModeFromFlags to PreviewAndApply would slip past the
// other JSON tests.
func TestAddJSON_DryRunDiffCombo(t *testing.T) {
	stub := &addUseCaseStub{
		resp: driving.AddServiceResponse{
			ServiceName: newAddSvcName(t, "postgres"),
			PriorState:  domain.ServiceStateUnregistered,
			State:       domain.ServiceStateActive,
			PlannedFiles: []driving.PlannedFile{
				{Path: "compose.yaml", Action: "modify", OldContent: []byte("services:\n  redis: {}\n"), NewContent: []byte("services:\n  redis: {}\n  postgres:\n    image: postgres:16\n")},
			},
		},
	}
	app := newAppWithAddStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"add", "postgres", "--dry-run", "--diff", "--json"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	raw := bytes.TrimSpace(stdout.Bytes())
	jsontestutil.AssertFullEnvelope(t, raw,
		jsontestutil.WithCommand("add"),
		jsontestutil.WithExitCode(0),
	)
	var env map[string]any
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got, _ := env["dryRun"].(bool); !got {
		t.Errorf("dryRun: want true (3-flag combo), envelope=%s", raw)
	}
	if got, _ := env["diff"].(bool); !got {
		t.Errorf("diff: want true (3-flag combo), envelope=%s", raw)
	}
	pfs, _ := env["plannedFiles"].([]any)
	if len(pfs) != 1 {
		t.Fatalf("plannedFiles: want 1, got %d", len(pfs))
	}
	pf, _ := pfs[0].(map[string]any)
	hunks, _ := pf["hunks"].([]any)
	if len(hunks) == 0 {
		t.Errorf("--dry-run --diff --json: hunks must be present, got %v", pf["hunks"])
	}
}

// TestAddJSON_FilesystemErrorShipsFullEnvelope pins the mid-write
// failure path (T0-(b) Scenario 2): ErrAddFileSystem produces a
// non-empty Response with PlannedFiles captured up to the failure
// point; the envelope ships voll-schema with diagnostic
// LH-NFA-REL-003 / exitCode 14 (T0-(j)).
func TestAddJSON_FilesystemErrorShipsFullEnvelope(t *testing.T) {
	stub := &addUseCaseStub{
		resp: driving.AddServiceResponse{
			ServiceName: newAddSvcName(t, "postgres"),
			PlannedFiles: []driving.PlannedFile{
				{Path: "u-boot.yaml", Action: "modify", OldContent: []byte("a\n"), NewContent: []byte("a\nb\n")},
				{Path: "compose.yaml", Action: "create", NewContent: []byte("svc: {}\n")},
			},
		},
		err: driving.ErrAddFileSystem,
	}
	app := newAppWithAddStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "add", "postgres", "--diff"}, &stdout, &stderr)
	// Review #2: the error MUST propagate so cli.ExitCode picks up
	// 14 (envelope body and shell exit must agree).
	if !errors.Is(err, driving.ErrAddFileSystem) {
		t.Fatalf("expected ErrAddFileSystem to propagate so ExitCode maps to 14; got err=%v", err)
	}
	if got := cli.ExitCode(err); got != 14 {
		t.Errorf("cli.ExitCode: want 14 (LH-NFA-REL-003), got %d", got)
	}
	if stdout.Len() == 0 {
		t.Fatalf("expected JSON envelope on stdout, got empty (err=%v)", err)
	}
	jsontestutil.AssertFullEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("add"),
		jsontestutil.WithExpectedCodes("LH-NFA-REL-003"),
		jsontestutil.WithExitCode(14),
	)
	var env map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}
	pfs, _ := env["plannedFiles"].([]any)
	if len(pfs) != 2 {
		t.Errorf("plannedFiles: want 2 (capture up to failure), got %d", len(pfs))
	}
}

// TestAddHumanDryRun_NoChangeLine pins that --dry-run without --json
// renders "Would add" / "Would change" instead of "Added" / "Changed"
// so the human-mode user sees the operation was a preview.
func TestAddHumanDryRun_NoChangeLine(t *testing.T) {
	stub := &addUseCaseStub{
		resp: driving.AddServiceResponse{
			ServiceName: newAddSvcName(t, "postgres"),
			PriorState:  domain.ServiceStateUnregistered,
			State:       domain.ServiceStateActive,
			Changed:     []string{"u-boot.yaml", "compose.yaml"},
		},
	}
	app := newAppWithAddStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"add", "postgres", "--dry-run"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Would add") {
		t.Errorf("--dry-run output should announce preview, got: %q", stdout.String())
	}
	if strings.Contains(stdout.String(), "Added service") {
		t.Errorf("--dry-run must not report past-tense add, got: %q", stdout.String())
	}
}

// TestAddHumanDiff_RendersHunksOnStdout pins the human-mode --diff
// path: the unified-diff body lands on stdout before the summary.
func TestAddHumanDiff_RendersHunksOnStdout(t *testing.T) {
	stub := &addUseCaseStub{
		resp: driving.AddServiceResponse{
			ServiceName: newAddSvcName(t, "postgres"),
			Changed:     []string{"compose.yaml"},
			PlannedFiles: []driving.PlannedFile{
				{Path: "compose.yaml", Action: "create", NewContent: []byte("services:\n  postgres: {}\n")},
			},
		},
	}
	app := newAppWithAddStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"add", "postgres", "--diff"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "@@") {
		t.Errorf("--diff output should contain @@-headers, got: %q", out)
	}
	if !strings.Contains(out, "+services:") {
		t.Errorf("--diff output should contain '+services:' line, got: %q", out)
	}
}
