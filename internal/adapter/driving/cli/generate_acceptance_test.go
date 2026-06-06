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

// =====================================================================
// generate-Acceptance-Tests (slice-v1-cli-json-dry-run-generate T6)
//
// Block 1: JSON-Pfade (4 Modi × representative artifact)
// Block 2: Error-Scenarios (ManualConflict pro Artefakt + URL-Reject +
//          ArtifactUnknown + ProjectNotInitialized + FS-Failure)
// Block 3: Action-Discriminator + NoOp Special-Pins
// Block 4: Human-Mode (printGenerateSummary + writeDiff)
//
// Alle Tests nutzen fakeGenerateUseCase aus cli_test.go plus
// WithGetwd-Stub für deterministische cwd-Auflösung.
// =====================================================================

// newAppWithGenerateStub wires the generate fake plus the deterministic
// getwd stub (T0-(k) Path-Anchor-Konsistenz analog init).
func newAppWithGenerateStub(stub driving.GenerateUseCase) *cli.App {
	return newAppWithGenerate(stub, cli.WithGetwd(func() (string, error) {
		return "/tmp/u-boot-generate-test", nil
	}))
}

// unmarshalGenerateEnv decodes the JSON envelope into a map for
// field-by-field inspection. Fails the test on parse error.
func unmarshalGenerateEnv(t *testing.T, raw []byte) map[string]any {
	t.Helper()
	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal envelope: %v\nraw=%s", err, raw)
	}
	return got
}

// =====================================================================
// Block 1: JSON-Pfade (T0-(f)/(m) Envelope-Shape)
// =====================================================================

// TestGenerateJSON_BareUsesDataEnvelope pins T0-(p): `u-boot --json
// generate changelog` ohne --dry-run/--diff produziert Minimal-Envelope
// mit data.artifact und data.action; KEIN plannedFiles/changes.
func TestGenerateJSON_BareUsesDataEnvelope(t *testing.T) {
	stub := &fakeGenerateUseCase{
		resp: driving.GenerateResponse{
			Artifact: domain.ArtifactChangelog,
			Action:   driving.GenerateActionCreated,
			Changed:  []string{"CHANGELOG.md"},
		},
	}
	app := newAppWithGenerateStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "generate", "changelog"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("generate"),
		jsontestutil.WithExitCode(0),
	)
	env := unmarshalGenerateEnv(t, stdout.Bytes())
	data, ok := env["data"].(map[string]any)
	if !ok {
		t.Fatalf("data missing or wrong type: %v", env["data"])
	}
	if got := data["artifact"]; got != "changelog" {
		t.Errorf("data.artifact: want changelog, got %v", got)
	}
	if got := data["action"]; got != "created" {
		t.Errorf("data.action: want created, got %v", got)
	}
}

// TestGenerateJSON_DryRunUsesFullEnvelope pins --dry-run --json:
// Voll-Schema mit plannedFiles/changes plus data.action.
func TestGenerateJSON_DryRunUsesFullEnvelope(t *testing.T) {
	stub := &fakeGenerateUseCase{
		resp: driving.GenerateResponse{
			Artifact: domain.ArtifactReadme,
			Action:   driving.GenerateActionCreated,
			Changed:  []string{"README.md"},
			PlannedFiles: []driving.PlannedFile{
				{Path: "README.md", Action: "create", NewContent: []byte("# project\n")},
			},
		},
	}
	app := newAppWithGenerateStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"generate", "readme", "--dry-run", "--json"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertFullEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("generate"),
		jsontestutil.WithExitCode(0),
	)
	env := unmarshalGenerateEnv(t, stdout.Bytes())
	if got, _ := env["dryRun"].(bool); !got {
		t.Errorf("dryRun: want true, envelope=%s", stdout.String())
	}
	if got, _ := env["diff"].(bool); got {
		t.Errorf("diff: want false, envelope=%s", stdout.String())
	}
	data := env["data"].(map[string]any)
	if got := data["artifact"]; got != "readme" {
		t.Errorf("data.artifact: want readme, got %v", got)
	}
	// PreviewMode propagation pin: stub captured PreviewDryRun.
	if stub.lastReq.PreviewMode != driving.PreviewDryRun {
		t.Errorf("req.PreviewMode: want PreviewDryRun, got %v", stub.lastReq.PreviewMode)
	}
}

// TestGenerateJSON_DiffWithoutDryRunRendersHunks pins --diff --json
// preview-and-apply: diff=true, dryRun=false, plannedFiles[].hunks
// populated.
func TestGenerateJSON_DiffWithoutDryRunRendersHunks(t *testing.T) {
	stub := &fakeGenerateUseCase{
		resp: driving.GenerateResponse{
			Artifact: domain.ArtifactEnvExample,
			Action:   driving.GenerateActionUpdatedBlock,
			Changed:  []string{".env.example"},
			PlannedFiles: []driving.PlannedFile{
				{
					Path:       ".env.example",
					Action:     "modify",
					OldContent: []byte("FOO=1\n"),
					NewContent: []byte("FOO=1\nBAR=2\n"),
				},
			},
		},
	}
	app := newAppWithGenerateStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"generate", "env-example", "--diff", "--json"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	env := unmarshalGenerateEnv(t, stdout.Bytes())
	if got, _ := env["dryRun"].(bool); got {
		t.Errorf("dryRun: want false, envelope=%s", stdout.String())
	}
	if got, _ := env["diff"].(bool); !got {
		t.Errorf("diff: want true, envelope=%s", stdout.String())
	}
	pfs, _ := env["plannedFiles"].([]any)
	if len(pfs) != 1 {
		t.Fatalf("plannedFiles: want 1, got %d", len(pfs))
	}
	first := pfs[0].(map[string]any)
	hunks, _ := first["hunks"].([]any)
	if len(hunks) == 0 {
		t.Errorf("plannedFiles[0].hunks must be non-empty for --diff modify; envelope=%s", stdout.String())
	}
	if stub.lastReq.PreviewMode != driving.PreviewAndApply {
		t.Errorf("req.PreviewMode: want PreviewAndApply, got %v", stub.lastReq.PreviewMode)
	}
}

// TestGenerateJSON_DryRunDiffCombo pins the three-flag combination:
// --dry-run --diff --json → dryRun=true, diff=true, no write.
func TestGenerateJSON_DryRunDiffCombo(t *testing.T) {
	stub := &fakeGenerateUseCase{
		resp: driving.GenerateResponse{
			Artifact: domain.ArtifactChangelog,
			Action:   driving.GenerateActionCreated,
			Changed:  []string{"CHANGELOG.md"},
			PlannedFiles: []driving.PlannedFile{
				{Path: "CHANGELOG.md", Action: "create", NewContent: []byte("# Changelog\n")},
			},
		},
	}
	app := newAppWithGenerateStub(stub)

	var stdout bytes.Buffer
	err := app.Execute(context.Background(), []string{"generate", "changelog", "--dry-run", "--diff", "--json"}, &stdout, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	env := unmarshalGenerateEnv(t, stdout.Bytes())
	if got, _ := env["dryRun"].(bool); !got {
		t.Errorf("dryRun: want true")
	}
	if got, _ := env["diff"].(bool); !got {
		t.Errorf("diff: want true")
	}
	if stub.lastReq.PreviewMode != driving.PreviewDryRun {
		t.Errorf("PreviewMode: want PreviewDryRun (dry-run wins over diff)")
	}
}

// =====================================================================
// Block 2: Error-Scenarios (T0-(e) per-Artefakt LH-Code-Tabelle)
// =====================================================================

// TestGenerateJSON_ManualConflictCodePerArtifact pins T0-(e): each
// artifact maps ErrGenerateManualConflict to its own LH-Code.
// changelog→GEN-002, readme→GEN-003, env-example→GEN-004,
// devcontainer→DEV-001.
func TestGenerateJSON_ManualConflictCodePerArtifact(t *testing.T) {
	cases := []struct {
		artifact string
		wantCode string
	}{
		{"changelog", "LH-FA-GEN-002"},
		{"readme", "LH-FA-GEN-003"},
		{"env-example", "LH-FA-GEN-004"},
		{"devcontainer", "LH-FA-DEV-001"},
	}
	for _, tc := range cases {
		t.Run(tc.artifact, func(t *testing.T) {
			stub := &fakeGenerateUseCase{err: driving.ErrGenerateManualConflict}
			app := newAppWithGenerateStub(stub)

			var stdout bytes.Buffer
			execErr := app.Execute(context.Background(), []string{"--json", "generate", tc.artifact}, &stdout, &bytes.Buffer{})
			if execErr == nil || !errors.Is(execErr, driving.ErrGenerateManualConflict) {
				t.Fatalf("want ErrGenerateManualConflict, got %v", execErr)
			}
			if cli.ExitCode(execErr) != 10 {
				t.Errorf("ExitCode: want 10, got %d", cli.ExitCode(execErr))
			}
			env := unmarshalGenerateEnv(t, stdout.Bytes())
			diags, _ := env["diagnostics"].([]any)
			if len(diags) != 1 {
				t.Fatalf("diagnostics: want 1, got %d", len(diags))
			}
			first := diags[0].(map[string]any)
			if got := first["code"]; got != tc.wantCode {
				t.Errorf("diagnostics[0].code: want %s, got %v", tc.wantCode, got)
			}
			// Error-Envelope MUSS data.artifact tragen (T0-(q)).
			data, ok := env["data"].(map[string]any)
			if !ok {
				t.Fatalf("data missing: %s", stdout.String())
			}
			if got := data["artifact"]; got != tc.artifact {
				t.Errorf("data.artifact: want %s, got %v", tc.artifact, got)
			}
			// Error-Envelope DARF KEIN data.action tragen (T0-(q)).
			if _, present := data["action"]; present {
				t.Errorf("data.action must be omitted on error path, got %v", data["action"])
			}
		})
	}
}

// TestGenerateJSON_URLRejectLHDEV003 pins the R6-HIGH-1-Finding:
// invalid --allow-external-feature-sources URL maps to LH-FA-DEV-003 /
// Exit 10 (Spec §720). Without the T3 ErrConfigValueInvalid wrap the
// path would have fallen to default LH-FA-CLI-006 / Exit 1.
func TestGenerateJSON_URLRejectLHDEV003(t *testing.T) {
	stub := &fakeGenerateUseCase{
		err: driving.ErrConfigValueInvalid,
	}
	app := newAppWithGenerateStub(stub)

	var stdout bytes.Buffer
	execErr := app.Execute(context.Background(), []string{
		"--json", "generate", "devcontainer",
		"--allow-external-feature-sources", "https://example.com/feature.json",
	}, &stdout, &bytes.Buffer{})
	if execErr == nil || !errors.Is(execErr, driving.ErrConfigValueInvalid) {
		t.Fatalf("want ErrConfigValueInvalid, got %v", execErr)
	}
	if cli.ExitCode(execErr) != 10 {
		t.Errorf("ExitCode: want 10, got %d", cli.ExitCode(execErr))
	}
	env := unmarshalGenerateEnv(t, stdout.Bytes())
	diags, _ := env["diagnostics"].([]any)
	if len(diags) != 1 {
		t.Fatalf("diagnostics: want 1")
	}
	if got := diags[0].(map[string]any)["code"]; got != "LH-FA-DEV-003" {
		t.Errorf("diagnostics[0].code: want LH-FA-DEV-003, got %v", got)
	}
}

// TestGenerateJSON_ArtifactUnknownExit2 pins LH-FA-GEN-001 Spec §1157:
// unknown artifact → Exit 2 (CLI-validation), not the default Exit 1.
// data is nil here because we have no artifact to embed.
func TestGenerateJSON_ArtifactUnknownExit2(t *testing.T) {
	stub := &fakeGenerateUseCase{}
	app := newAppWithGenerateStub(stub)

	var stdout bytes.Buffer
	execErr := app.Execute(context.Background(), []string{"--json", "generate", "bogus"}, &stdout, &bytes.Buffer{})
	if execErr == nil || !errors.Is(execErr, driving.ErrArtifactUnknown) {
		t.Fatalf("want ErrArtifactUnknown, got %v", execErr)
	}
	if cli.ExitCode(execErr) != 2 {
		t.Errorf("ExitCode: want 2, got %d", cli.ExitCode(execErr))
	}
	env := unmarshalGenerateEnv(t, stdout.Bytes())
	diags, _ := env["diagnostics"].([]any)
	if len(diags) != 1 {
		t.Fatalf("diagnostics: want 1")
	}
	if got := diags[0].(map[string]any)["code"]; got != "LH-FA-CLI-006" {
		t.Errorf("diagnostics[0].code: want LH-FA-CLI-006, got %v", got)
	}
	if _, present := env["data"]; present {
		t.Errorf("data must be omitted for unknown artifact (no classification possible)")
	}
}

// TestGenerateJSON_ProjectNotInitializedExit10 pins the reused
// LH-FA-INIT-001 sentinel: u-boot.yaml missing → fachlicher Exit 10.
func TestGenerateJSON_ProjectNotInitializedExit10(t *testing.T) {
	stub := &fakeGenerateUseCase{err: driving.ErrProjectNotInitialized}
	app := newAppWithGenerateStub(stub)

	var stdout bytes.Buffer
	execErr := app.Execute(context.Background(), []string{"--json", "generate", "changelog"}, &stdout, &bytes.Buffer{})
	if execErr == nil || !errors.Is(execErr, driving.ErrProjectNotInitialized) {
		t.Fatalf("want ErrProjectNotInitialized, got %v", execErr)
	}
	if cli.ExitCode(execErr) != 10 {
		t.Errorf("ExitCode: want 10, got %d", cli.ExitCode(execErr))
	}
	env := unmarshalGenerateEnv(t, stdout.Bytes())
	diags, _ := env["diagnostics"].([]any)
	if got := diags[0].(map[string]any)["code"]; got != "LH-FA-INIT-001" {
		t.Errorf("diagnostics[0].code: want LH-FA-INIT-001, got %v", got)
	}
}

// TestGenerateJSON_FSFailureExit14 pins LH-NFA-REL-003 / Exit 14:
// any ErrGenerateFileSystem-wrapped path maps to the technical-
// persistence class.
func TestGenerateJSON_FSFailureExit14(t *testing.T) {
	stub := &fakeGenerateUseCase{err: driving.ErrGenerateFileSystem}
	app := newAppWithGenerateStub(stub)

	var stdout bytes.Buffer
	execErr := app.Execute(context.Background(), []string{"--json", "generate", "changelog"}, &stdout, &bytes.Buffer{})
	if execErr == nil || !errors.Is(execErr, driving.ErrGenerateFileSystem) {
		t.Fatalf("want ErrGenerateFileSystem, got %v", execErr)
	}
	if cli.ExitCode(execErr) != 14 {
		t.Errorf("ExitCode: want 14, got %d", cli.ExitCode(execErr))
	}
	env := unmarshalGenerateEnv(t, stdout.Bytes())
	diags, _ := env["diagnostics"].([]any)
	if got := diags[0].(map[string]any)["code"]; got != "LH-NFA-REL-003" {
		t.Errorf("diagnostics[0].code: want LH-NFA-REL-003, got %v", got)
	}
}

// TestGenerateJSON_AllowExternalOnNonDevcontainerRejects pins the
// pre-use-case mutex check: --allow-external-feature-sources on
// changelog/readme/env-example rejects without touching the UC.
func TestGenerateJSON_AllowExternalOnNonDevcontainerRejects(t *testing.T) {
	stub := &fakeGenerateUseCase{}
	app := newAppWithGenerateStub(stub)

	var stdout bytes.Buffer
	execErr := app.Execute(context.Background(), []string{
		"--json", "generate", "readme",
		"--allow-external-feature-sources", "https://example.com/x",
	}, &stdout, &bytes.Buffer{})
	if execErr == nil || !errors.Is(execErr, driving.ErrArtifactUnknown) {
		t.Fatalf("want ErrArtifactUnknown, got %v", execErr)
	}
	if cli.ExitCode(execErr) != 2 {
		t.Errorf("ExitCode: want 2, got %d", cli.ExitCode(execErr))
	}
	if stub.called {
		t.Errorf("use-case must NOT be invoked when the pre-mutex check fires")
	}
}

// =====================================================================
// Block 3: Action-Discriminator + NoOp Special-Pin (T0-(f))
// =====================================================================

// TestGenerateJSON_NoOpEmptyArrays pins T0-(f) Festzurrung: NoOp →
// plannedFiles: [], changes: [], data.action: "no-op". Konsumenten
// leiten NoOp aus Array-Leerheit + data.action ab.
func TestGenerateJSON_NoOpEmptyArrays(t *testing.T) {
	stub := &fakeGenerateUseCase{
		resp: driving.GenerateResponse{
			Artifact: domain.ArtifactChangelog,
			Action:   driving.GenerateActionNoOp,
			// PlannedFiles and Changed both nil/empty.
		},
	}
	app := newAppWithGenerateStub(stub)

	var stdout bytes.Buffer
	err := app.Execute(context.Background(), []string{"generate", "changelog", "--dry-run", "--json"}, &stdout, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	env := unmarshalGenerateEnv(t, stdout.Bytes())
	pfs, _ := env["plannedFiles"].([]any)
	if len(pfs) != 0 {
		t.Errorf("plannedFiles: want empty for NoOp, got %d entries", len(pfs))
	}
	chs, _ := env["changes"].([]any)
	if len(chs) != 0 {
		t.Errorf("changes: want empty for NoOp, got %d entries", len(chs))
	}
	data := env["data"].(map[string]any)
	if got := data["action"]; got != "no-op" {
		t.Errorf("data.action: want no-op, got %v", got)
	}
}

// TestGenerateJSON_ActionDiscriminator pins T0-(f) R3-Finding:
// UpdatedBlock and RepairedManual produce identical
// plannedFiles[i].action: "modify"; data.action is the ONLY
// discriminator between them.
func TestGenerateJSON_ActionDiscriminator(t *testing.T) {
	cases := []struct {
		name       string
		action     driving.GenerateAction
		wantAction string
	}{
		{"UpdatedBlock", driving.GenerateActionUpdatedBlock, "updated-block"},
		{"RepairedManual", driving.GenerateActionRepairedManual, "repaired-manual"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stub := &fakeGenerateUseCase{
				resp: driving.GenerateResponse{
					Artifact: domain.ArtifactChangelog,
					Action:   tc.action,
					Changed:  []string{"CHANGELOG.md"},
					PlannedFiles: []driving.PlannedFile{
						{
							Path:       "CHANGELOG.md",
							Action:     "modify",
							OldContent: []byte("a\n"),
							NewContent: []byte("a\nb\n"),
						},
					},
				},
			}
			app := newAppWithGenerateStub(stub)

			var stdout bytes.Buffer
			err := app.Execute(context.Background(), []string{"generate", "changelog", "--dry-run", "--json"}, &stdout, &bytes.Buffer{})
			if err != nil {
				t.Fatalf("execute: %v", err)
			}
			env := unmarshalGenerateEnv(t, stdout.Bytes())
			data := env["data"].(map[string]any)
			if got := data["action"]; got != tc.wantAction {
				t.Errorf("data.action: want %s, got %v", tc.wantAction, got)
			}
			pfs := env["plannedFiles"].([]any)
			first := pfs[0].(map[string]any)
			// FS-Semantik identisch — Discriminator lebt NUR in data.action.
			if got := first["action"]; got != "modify" {
				t.Errorf("plannedFiles[0].action: want modify (FS-semantisch identisch), got %v", got)
			}
		})
	}
}

// =====================================================================
// Block 4: Human-Mode
// =====================================================================

// TestGenerateHuman_NoOpSummary pins the printGenerateSummary NoOp
// path: no Changed list, no diff output.
func TestGenerateHuman_NoOpSummary(t *testing.T) {
	stub := &fakeGenerateUseCase{
		resp: driving.GenerateResponse{
			Artifact: domain.ArtifactChangelog,
			Action:   driving.GenerateActionNoOp,
		},
	}
	app := newAppWithGenerateStub(stub)

	var stdout bytes.Buffer
	err := app.Execute(context.Background(), []string{"generate", "changelog"}, &stdout, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(stdout.String(), "already up to date") {
		t.Errorf("stdout must mention 'already up to date', got %q", stdout.String())
	}
}

// TestGenerateHuman_DiffRendersUnifiedDiff pins LH-FA-CLI-008 human-
// mode: `--diff` without --json renders a unified diff string before
// the summary line.
func TestGenerateHuman_DiffRendersUnifiedDiff(t *testing.T) {
	stub := &fakeGenerateUseCase{
		resp: driving.GenerateResponse{
			Artifact: domain.ArtifactReadme,
			Action:   driving.GenerateActionCreated,
			Changed:  []string{"README.md"},
			PlannedFiles: []driving.PlannedFile{
				{Path: "README.md", Action: "create", NewContent: []byte("# project\n")},
			},
		},
	}
	app := newAppWithGenerateStub(stub)

	var stdout bytes.Buffer
	err := app.Execute(context.Background(), []string{"generate", "readme", "--diff"}, &stdout, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	out := stdout.String()
	if !strings.Contains(out, "--- README.md (create)") {
		t.Errorf("stdout must contain unified-diff header for README.md; got %q", out)
	}
	if !strings.Contains(out, "Generated readme") {
		t.Errorf("stdout must end with printGenerateSummary line; got %q", out)
	}
}
