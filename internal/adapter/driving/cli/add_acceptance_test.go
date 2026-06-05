package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli/jsontestutil"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// Acceptance fixtures: spec-konforme content snapshots from the Plan
// §Aufhebungsbedingung. Variante A = frisch-init project (compose.yaml
// does not exist, 12-line postgres block created); Variante B =
// existing setup (compose.yaml already has redis, 6 lines appended).
//
// composeFreshAdd is the new compose.yaml body for the create-from-
// scratch case — 12 lines exactly (Spec §430 `count: 12`).
const composeFreshAdd = `services:
  postgres:
    image: postgres:16
    restart: unless-stopped
    environment:
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    volumes:
      - postgres-data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
`

// composeExistingOld is the pre-add compose body (redis only).
const composeExistingOld = `services:
  redis:
    image: redis:7
    restart: unless-stopped
`

// composeExistingNew appends a six-line postgres block (Spec §477
// `count: 6` — the added lines, not context). The block intentionally
// totals SIX additive lines (matching the Spec example exactly), not
// seven — including a fourth attribute would push count to 7 and
// drift against the §477 canonical example.
const composeExistingNew = `services:
  redis:
    image: redis:7
    restart: unless-stopped
  postgres:
    image: postgres:16
    environment:
      POSTGRES_DB: ${POSTGRES_DB}
    ports:
      - "5432:5432"
`

// pgServiceName is a shorthand for the validated postgres ServiceName.
func pgServiceName(t *testing.T) domain.ServiceName {
	t.Helper()
	n, err := domain.NewServiceName("postgres")
	if err != nil {
		t.Fatalf("NewServiceName(postgres): %v", err)
	}
	return n
}

// unmarshalEnv parses the JSON output for further structural pins.
func unmarshalEnv(t *testing.T, raw []byte) map[string]any {
	t.Helper()
	var env map[string]any
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatalf("unmarshal envelope: %v\nraw=%s", err, raw)
	}
	return env
}

// TestAddAcceptance_VarianteA_FreshInit_PinsCreateCount12 is the
// Variante-A pin from Plan §Aufhebungsbedingung: the create-from-
// scratch postgres block has count=12 (Spec §430). Anti-Drift:
// changing the CountLines formula or breaking the create-action
// path would surface here before downstream consumers notice.
func TestAddAcceptance_VarianteA_FreshInit_PinsCreateCount12(t *testing.T) {
	stub := &addUseCaseStub{
		resp: driving.AddServiceResponse{
			ServiceName: pgServiceName(t),
			PriorState:  domain.ServiceStateUnregistered,
			State:       domain.ServiceStateActive,
			PlannedFiles: []driving.PlannedFile{
				{Path: "compose.yaml", Action: "create", NewContent: []byte(composeFreshAdd)},
			},
		},
	}
	app := newAppWithAddStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"add", "postgres", "--dry-run", "--json"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	raw := bytes.TrimSpace(stdout.Bytes())
	jsontestutil.AssertFullEnvelope(t, raw,
		jsontestutil.WithCommand("add"),
		jsontestutil.WithExitCode(0),
	)

	env := unmarshalEnv(t, raw)
	changes, _ := env["changes"].([]any)
	if len(changes) != 1 {
		t.Fatalf("changes: want 1, got %d (envelope=%s)", len(changes), raw)
	}
	first, _ := changes[0].(map[string]any)
	if got, _ := first["count"].(float64); int(got) != 12 {
		t.Errorf("Variante A compose.yaml count: want 12 (Spec §430), got %v", first["count"])
	}
	pfs, _ := env["plannedFiles"].([]any)
	pf, _ := pfs[0].(map[string]any)
	if act, _ := pf["action"].(string); act != "create" {
		t.Errorf("Variante A action: want create, got %q", act)
	}
}

// TestAddAcceptance_VarianteB_Existing_PinsModifyAndAddedLines is
// the Variante-B pin from Plan §Aufhebungsbedingung: an existing
// compose.yaml with another service gets 6 lines appended. The
// formal T0-(g) form (CountFromHunks = sum(hunk.NewLines)) reports
// inserts + context, so the floor invariant is six `+` lines in
// the hunk content (Spec §477 `count: 6` refers to the additions,
// the formal sum may be higher).
func TestAddAcceptance_VarianteB_Existing_PinsModifyAndAddedLines(t *testing.T) {
	stub := &addUseCaseStub{
		resp: driving.AddServiceResponse{
			ServiceName: pgServiceName(t),
			PriorState:  domain.ServiceStateUnregistered,
			State:       domain.ServiceStateActive,
			PlannedFiles: []driving.PlannedFile{
				{
					Path:       "compose.yaml",
					Action:     "modify",
					OldContent: []byte(composeExistingOld),
					NewContent: []byte(composeExistingNew),
				},
			},
		},
	}
	app := newAppWithAddStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"add", "postgres", "--diff", "--json"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	raw := bytes.TrimSpace(stdout.Bytes())
	jsontestutil.AssertFullEnvelope(t, raw,
		jsontestutil.WithCommand("add"),
		jsontestutil.WithExitCode(0),
	)

	env := unmarshalEnv(t, raw)
	pfs, _ := env["plannedFiles"].([]any)
	if len(pfs) != 1 {
		t.Fatalf("plannedFiles: want 1, got %d (envelope=%s)", len(pfs), raw)
	}
	pf, _ := pfs[0].(map[string]any)
	if act, _ := pf["action"].(string); act != "modify" {
		t.Errorf("Variante B action: want modify, got %q", act)
	}
	hunks, _ := pf["hunks"].([]any)
	if len(hunks) == 0 {
		t.Fatalf("Variante B hunks: want non-empty, got %v", pf["hunks"])
	}
	hunk0, _ := hunks[0].(map[string]any)
	content, _ := hunk0["content"].(string)
	additions := strings.Count(content, "\n+") + boolToInt(strings.HasPrefix(content, "+"))
	if additions != 6 {
		t.Errorf("Variante B: want exactly 6 '+' lines (Spec §477), got %d in content=%q", additions, content)
	}
	// changes[].count = CountAdditions (review-round-7 B): true
	// additive lines, NOT additions + context. Spec §477 example
	// pins this to 6 exactly.
	changes, _ := env["changes"].([]any)
	first, _ := changes[0].(map[string]any)
	if got, _ := first["count"].(float64); int(got) != 6 {
		t.Errorf("Variante B count: want exactly 6 (Spec §477), got %v", first["count"])
	}
}

// TestAddAcceptance_DiffJSON_HunkStructurePin pins concrete hunk
// shape fields against the Variante-B fixture so a regression in
// LCS coordinate-arithmetic surfaces as a coordinate mismatch
// rather than a vague test failure. OldLines/NewLines are
// formally constrained; OldStart/NewStart must be 1-based when
// the matching Lines field is positive (T0-(l)).
func TestAddAcceptance_DiffJSON_HunkStructurePin(t *testing.T) {
	stub := &addUseCaseStub{
		resp: driving.AddServiceResponse{
			ServiceName: pgServiceName(t),
			PriorState:  domain.ServiceStateUnregistered,
			State:       domain.ServiceStateActive,
			PlannedFiles: []driving.PlannedFile{
				{
					Path:       "compose.yaml",
					Action:     "modify",
					OldContent: []byte(composeExistingOld),
					NewContent: []byte(composeExistingNew),
				},
			},
		},
	}
	app := newAppWithAddStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"add", "postgres", "--diff", "--json"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	env := unmarshalEnv(t, bytes.TrimSpace(stdout.Bytes()))
	pfs, _ := env["plannedFiles"].([]any)
	pf, _ := pfs[0].(map[string]any)
	hunks, _ := pf["hunks"].([]any)
	hunk0, _ := hunks[0].(map[string]any)

	// OldStart must be ≥ 1 when OldLines > 0 — Variante B has the
	// redis block as 4 context lines before the addition. With the
	// default 3-line context, the hunk starts at op-index 1 (after
	// the first Equal "services:" line gets clipped off), so
	// oldStart=newStart=2 ("  redis:"). The invariant pinned here is
	// 1-based-coordinates, not a specific value.
	oldStart, _ := hunk0["oldStart"].(float64)
	oldLines, _ := hunk0["oldLines"].(float64)
	newStart, _ := hunk0["newStart"].(float64)
	newLines, _ := hunk0["newLines"].(float64)
	if oldLines <= 0 {
		t.Errorf("Variante B: oldLines must be > 0 (context lines present), got %v", oldLines)
	}
	if oldStart < 1 || newStart < 1 {
		t.Errorf("Variante B: hunk 1-based coordinates must be ≥ 1 (T0-(l)), got oldStart=%v newStart=%v", oldStart, newStart)
	}
	if newLines <= oldLines {
		t.Errorf("Variante B: newLines (%v) must be > oldLines (%v) for an additive modify", newLines, oldLines)
	}
}

// TestAddAcceptance_IdempotentNoOp_EmptyPlanAndChanges pins the
// Spec-§326 voll-schema shape for the idempotent re-add case
// (PriorState=Active, Changed=nil): plannedFiles: [] AND
// changes: [] (both required even when empty), status: ok,
// exitCode: 0. This is the success path of T0-(b) Scenario 1's
// no-op variant — the recorder captured nothing because the
// service is already in the right state.
func TestAddAcceptance_IdempotentNoOp_EmptyPlanAndChanges(t *testing.T) {
	stub := &addUseCaseStub{
		resp: driving.AddServiceResponse{
			ServiceName:  pgServiceName(t),
			PriorState:   domain.ServiceStateActive,
			State:        domain.ServiceStateActive,
			PlannedFiles: nil, // recorder captured no mutations
		},
	}
	app := newAppWithAddStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"add", "postgres", "--dry-run", "--json"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	raw := bytes.TrimSpace(stdout.Bytes())
	jsontestutil.AssertFullEnvelope(t, raw,
		jsontestutil.WithCommand("add"),
		jsontestutil.WithExitCode(0),
	)

	env := unmarshalEnv(t, raw)
	pfs, ok := env["plannedFiles"].([]any)
	if !ok || len(pfs) != 0 {
		t.Errorf("idempotent no-op: plannedFiles must be present as empty array, got %v", env["plannedFiles"])
	}
	chs, ok := env["changes"].([]any)
	if !ok || len(chs) != 0 {
		t.Errorf("idempotent no-op: changes must be present as empty array, got %v", env["changes"])
	}
	if env["status"] != "ok" {
		t.Errorf("idempotent no-op: status: want ok, got %v", env["status"])
	}
}

// TestAddAcceptance_SuccessScenario_AllThreeFilesCaptured is the
// happy-path pin from T0-(b) Scenario 1: a fresh add captures all
// three files (u-boot.yaml + compose.yaml + .env.example). The
// envelope ships voll-schema with three plannedFiles, three
// changes, status: ok, exitCode: 0.
func TestAddAcceptance_SuccessScenario_AllThreeFilesCaptured(t *testing.T) {
	stub := &addUseCaseStub{
		resp: driving.AddServiceResponse{
			ServiceName: pgServiceName(t),
			PriorState:  domain.ServiceStateUnregistered,
			State:       domain.ServiceStateActive,
			PlannedFiles: []driving.PlannedFile{
				{Path: "u-boot.yaml", Action: "modify", OldContent: []byte("services:\n"), NewContent: []byte("services:\n  postgres:\n    enabled: true\n")},
				{Path: "compose.yaml", Action: "create", NewContent: []byte(composeFreshAdd)},
				{Path: ".env.example", Action: "create", NewContent: []byte("POSTGRES_DB=app\nPOSTGRES_USER=app\nPOSTGRES_PASSWORD=changeme\n")},
			},
		},
	}
	app := newAppWithAddStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"add", "postgres", "--dry-run", "--json"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	raw := bytes.TrimSpace(stdout.Bytes())
	jsontestutil.AssertFullEnvelope(t, raw,
		jsontestutil.WithCommand("add"),
		jsontestutil.WithExitCode(0),
	)
	env := unmarshalEnv(t, raw)
	pfs, _ := env["plannedFiles"].([]any)
	if len(pfs) != 3 {
		t.Errorf("Scenario 1 success: want 3 plannedFiles, got %d", len(pfs))
	}
	chs, _ := env["changes"].([]any)
	if len(chs) != 3 {
		t.Errorf("Scenario 1 success: want 3 changes, got %d", len(chs))
	}
	// Pin compose.yaml's count to 12 — the canonical Spec §430 value.
	for _, raw := range chs {
		item, _ := raw.(map[string]any)
		if item["path"] == "compose.yaml" {
			if got, _ := item["count"].(float64); int(got) != 12 {
				t.Errorf("compose.yaml count: want 12 (Spec §430), got %v", item["count"])
			}
		}
	}
}

// TestAddAcceptance_HumanDiff_RendersPostgresBlock pins the
// human-mode --diff output for Variante B: the unified-diff body
// must include the `+postgres:` line, the @@ header, and the
// service-trailer line.
func TestAddAcceptance_HumanDiff_RendersPostgresBlock(t *testing.T) {
	stub := &addUseCaseStub{
		resp: driving.AddServiceResponse{
			ServiceName: pgServiceName(t),
			PriorState:  domain.ServiceStateUnregistered,
			State:       domain.ServiceStateActive,
			Changed:     []string{"compose.yaml"},
			PlannedFiles: []driving.PlannedFile{
				{
					Path:       "compose.yaml",
					Action:     "modify",
					OldContent: []byte(composeExistingOld),
					NewContent: []byte(composeExistingNew),
				},
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
	for _, want := range []string{
		"--- compose.yaml (modify)",
		"@@ -",
		"+  postgres:",
		"+    image: postgres:16",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("human-mode diff missing %q\nfull output:\n%s", want, out)
		}
	}
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
