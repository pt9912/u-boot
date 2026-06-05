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

// initUseCaseStub returns a fixed response/err and records the last
// request. Stripped-down compared to the test-package fakeInitUseCase
// so each acceptance test fills only what it needs.
type initUseCaseStub struct {
	resp    driving.InitProjectResponse
	err     error
	lastReq driving.InitProjectRequest
}

func (s *initUseCaseStub) Init(_ context.Context, req driving.InitProjectRequest) (driving.InitProjectResponse, error) {
	s.lastReq = req
	return s.resp, s.err
}

// newAppWithInitStub wires the init use-case stub plus a deterministic
// getwd-stub so the tests do not depend on the runner's CWD
// (T0-(l)/(k) Path-Anchor-Konsistenz).
func newAppWithInitStub(stub driving.InitProjectUseCase) *cli.App {
	return newApp(stub, cli.WithGetwd(func() (string, error) { return "/tmp/u-boot-init-test", nil }))
}

// mustProjectName is defined in cli_test.go (same package).

// =====================================================================
// JSON-Pfade (T0-(k) drei Output-Formen)
// =====================================================================

// TestInitJSON_BareUsesMinimalEnvelope pins T0-(j): `u-boot init --json`
// ohne --dry-run/--diff trägt nur den Spec-§1841-Minimalkontrakt
// (analog add T0-(k)).
func TestInitJSON_BareUsesMinimalEnvelope(t *testing.T) {
	stub := &initUseCaseStub{
		resp: driving.InitProjectResponse{Project: domain.NewProject(mustProjectName(t, "myproj"))},
	}
	app := newAppWithInitStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "init", "myproj"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("init"),
		jsontestutil.WithExitCode(0),
	)
}

// TestInitJSON_DryRunUsesFullEnvelope pins the --dry-run --json
// Voll-Schema-Pfad: dryRun=true, diff=false, plannedFiles from
// recorder mapping.
func TestInitJSON_DryRunUsesFullEnvelope(t *testing.T) {
	stub := &initUseCaseStub{
		resp: driving.InitProjectResponse{
			Project: domain.NewProject(mustProjectName(t, "myproj")),
			PlannedFiles: []driving.PlannedFile{
				{Path: "docker", Action: "create"},
				{Path: "scripts", Action: "create"},
				{Path: "docs", Action: "create"},
				{Path: "README.md", Action: "create", NewContent: []byte("# myproj\n")},
				{Path: "u-boot.yaml", Action: "create", NewContent: []byte("project:\n  name: myproj\n")},
			},
		},
	}
	app := newAppWithInitStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"init", "myproj", "--dry-run", "--json"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertFullEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("init"),
		jsontestutil.WithExitCode(0),
	)
	env := unmarshalInitEnv(t, stdout.Bytes())
	if got, _ := env["dryRun"].(bool); !got {
		t.Errorf("dryRun: want true, envelope=%s", stdout.String())
	}
	if got, _ := env["diff"].(bool); got {
		t.Errorf("diff: want false, envelope=%s", stdout.String())
	}
	pfs, _ := env["plannedFiles"].([]any)
	if len(pfs) != 5 {
		t.Errorf("plannedFiles: want 5, got %d (envelope=%s)", len(pfs), stdout.String())
	}
}

// TestInitJSON_DiffWithoutDryRunRendersHunks pins the --diff --json
// preview-and-apply-Pfad: diff=true, dryRun=false, plannedFiles[].
// hunks populiert für non-binary files.
func TestInitJSON_DiffWithoutDryRunRendersHunks(t *testing.T) {
	stub := &initUseCaseStub{
		resp: driving.InitProjectResponse{
			Project: domain.NewProject(mustProjectName(t, "myproj")),
			PlannedFiles: []driving.PlannedFile{
				{Path: "u-boot.yaml", Action: "create", NewContent: []byte("project:\n  name: myproj\n  enabled: true\n")},
			},
		},
	}
	app := newAppWithInitStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"init", "myproj", "--diff", "--json"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertFullEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("init"),
		jsontestutil.WithExitCode(0),
	)
	env := unmarshalInitEnv(t, stdout.Bytes())
	if got, _ := env["diff"].(bool); !got {
		t.Errorf("diff: want true, envelope=%s", stdout.String())
	}
	if got, _ := env["dryRun"].(bool); got {
		t.Errorf("dryRun: want false, envelope=%s", stdout.String())
	}
	pfs, _ := env["plannedFiles"].([]any)
	if len(pfs) == 0 {
		t.Fatalf("plannedFiles: empty, envelope=%s", stdout.String())
	}
	first, _ := pfs[0].(map[string]any)
	hunks, ok := first["hunks"].([]any)
	if !ok || len(hunks) == 0 {
		t.Errorf("plannedFiles[0].hunks: want non-empty, got %v", first["hunks"])
	}
}

// TestInitJSON_DryRunDiffCombo pins the 3-Flag-Combo
// `--dry-run --diff --json` (T0-(b) (yes, yes)-Cell: --dry-run wins;
// envelope MUSS dryRun=true UND diff=true tragen).
func TestInitJSON_DryRunDiffCombo(t *testing.T) {
	stub := &initUseCaseStub{
		resp: driving.InitProjectResponse{
			Project: domain.NewProject(mustProjectName(t, "myproj")),
			PlannedFiles: []driving.PlannedFile{
				{Path: "u-boot.yaml", Action: "create", NewContent: []byte("project:\n  name: myproj\n")},
			},
		},
	}
	app := newAppWithInitStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"init", "myproj", "--dry-run", "--diff", "--json"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	env := unmarshalInitEnv(t, stdout.Bytes())
	if got, _ := env["dryRun"].(bool); !got {
		t.Errorf("3-flag combo: dryRun must be true")
	}
	if got, _ := env["diff"].(bool); !got {
		t.Errorf("3-flag combo: diff must be true")
	}
	pfs, _ := env["plannedFiles"].([]any)
	if len(pfs) == 0 {
		t.Errorf("3-flag combo: plannedFiles must be populated")
	}
}

// =====================================================================
// Error-Pfade (T0-(b) Failure-Scenarios + T0-(i) Template-Reject)
// =====================================================================

// TestInitJSON_SoftExistingShipsMinimalErrorEnvelope pins T0-(b)
// Soft-Existing-Detection-Failure: ErrProjectExists → exitCode 10 /
// LH-FA-INIT-004 / minimal-envelope (Planning-Phase, kein Recorder).
func TestInitJSON_SoftExistingShipsMinimalErrorEnvelope(t *testing.T) {
	stub := &initUseCaseStub{err: driving.ErrProjectExists}
	app := newAppWithInitStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "init", "myproj"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrProjectExists) {
		t.Fatalf("expected ErrProjectExists, got: %v", err)
	}
	if got := cli.ExitCode(err); got != 10 {
		t.Errorf("cli.ExitCode: want 10, got %d", got)
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("init"),
		jsontestutil.WithExpectedCodes("LH-FA-INIT-004"),
		jsontestutil.WithExitCode(10),
	)
}

// TestInitJSON_PlanningPhaseForceFailure pins T0-(q): `init --force`
// ohne --backup auf existierender .gitignore failed in planFile
// BEVOR Recorder anyone captures. PlannedFiles bleibt leer; envelope
// trägt Voll-Schema (--dry-run gesetzt) mit error-Diagnostic
// LH-FA-INIT-005 / exitCode 10.
func TestInitJSON_PlanningPhaseForceFailure(t *testing.T) {
	stub := &initUseCaseStub{err: driving.ErrForceRequiresBackup}
	app := newAppWithInitStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "init", "myproj", "--force", "--dry-run"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrForceRequiresBackup) {
		t.Fatalf("expected ErrForceRequiresBackup, got: %v", err)
	}
	if got := cli.ExitCode(err); got != 10 {
		t.Errorf("cli.ExitCode: want 10 (Usage-Klasse), got %d", got)
	}
	jsontestutil.AssertFullEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("init"),
		jsontestutil.WithExpectedCodes("LH-FA-INIT-005"),
		jsontestutil.WithExitCode(10),
	)
	env := unmarshalInitEnv(t, stdout.Bytes())
	pfs, _ := env["plannedFiles"].([]any)
	if len(pfs) != 0 {
		t.Errorf("planning-phase failure: plannedFiles must be empty (kein Recorder-Capture), got %d", len(pfs))
	}
}

// TestInitJSON_MidWriteFailureShipsFullEnvelope pins T0-(b) Mid-Write-
// Failure: ErrInitFileSystem mit non-empty PlannedFiles (Recorder
// capturte Calls bis Failure-Stelle). Voll-Schema-Envelope mit
// LH-NFA-REL-003 / exitCode 14 / file-Pointer auf letzte Capture.
func TestInitJSON_MidWriteFailureShipsFullEnvelope(t *testing.T) {
	stub := &initUseCaseStub{
		resp: driving.InitProjectResponse{
			Project: domain.NewProject(mustProjectName(t, "myproj")),
			PlannedFiles: []driving.PlannedFile{
				{Path: "docker", Action: "create"},
				{Path: "README.md", Action: "create", NewContent: []byte("# myproj\n")},
			},
		},
		err: driving.ErrInitFileSystem,
	}
	app := newAppWithInitStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "init", "myproj", "--diff"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrInitFileSystem) {
		t.Fatalf("expected ErrInitFileSystem, got: %v", err)
	}
	if got := cli.ExitCode(err); got != 14 {
		t.Errorf("cli.ExitCode: want 14 (LH-NFA-REL-003), got %d", got)
	}
	jsontestutil.AssertFullEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("init"),
		jsontestutil.WithExpectedCodes("LH-NFA-REL-003"),
		jsontestutil.WithExitCode(14),
	)
	env := unmarshalInitEnv(t, stdout.Bytes())
	pfs, _ := env["plannedFiles"].([]any)
	if len(pfs) != 2 {
		t.Errorf("mid-write failure: plannedFiles must carry capture-up-to-failure (2), got %d", len(pfs))
	}
	// File-Annotation: diagnostics[0].file points to the last
	// captured path (T0-(f) lastPlannedPath erblich).
	diags, _ := env["diagnostics"].([]any)
	if len(diags) > 0 {
		diag, _ := diags[0].(map[string]any)
		if file, _ := diag["file"].(string); file != "README.md" {
			t.Errorf("diagnostics[0].file: want \"README.md\" (last captured path), got %q", file)
		}
	}
}

// TestInitJSON_TemplateRejectShipsErrorEnvelope pins T0-(i)
// Out-of-Scope-Carveout: `--template + --dry-run|--diff` rejects am
// CLI-Level mit ErrTemplateConflictsWithFlag / LH-FA-CLI-006 /
// exitCode 2. CLI ruft uc.Init NICHT auf.
//
// Envelope-Form (add R6 #4 erblich): wantsFullSchema ist true,
// sobald --dry-run ODER --diff gesetzt ist (User-Intent
// reflektieren); plannedFiles/changes bleiben empty arrays weil
// kein Recorder-Capture stattfand. Plan-Pin-Form aus der Aufhebungs-
// bedingung (Minimal-Envelope) war hier inkonsistent zur add-R6-#4-
// Regel — der Test pinnt die konsistente Voll-Schema-Form.
func TestInitJSON_TemplateRejectShipsErrorEnvelope(t *testing.T) {
	stub := &initUseCaseStub{}
	app := newAppWithInitStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "init", "myproj", "--template", "basic", "--dry-run"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrTemplateConflictsWithFlag) {
		t.Fatalf("expected ErrTemplateConflictsWithFlag, got: %v", err)
	}
	if got := cli.ExitCode(err); got != 2 {
		t.Errorf("cli.ExitCode: want 2 (Usage-Klasse), got %d", got)
	}
	jsontestutil.AssertFullEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("init"),
		jsontestutil.WithExpectedCodes("LH-FA-CLI-006"),
		jsontestutil.WithExitCode(2),
	)
	env := unmarshalInitEnv(t, stdout.Bytes())
	pfs, _ := env["plannedFiles"].([]any)
	if len(pfs) != 0 {
		t.Errorf("template-reject: plannedFiles must be empty (CLI-Level vor uc.Init), got %d", len(pfs))
	}
	if stub.lastReq.BaseDir != "" {
		t.Errorf("uc.Init must NOT be called when --template + --dry-run rejects at CLI level; lastReq=%+v", stub.lastReq)
	}
}

// TestInitJSON_ConflictingModeFlagsShipsMinimalEnvelope pins that
// --yes + --no-interactive + --json produces a JSON envelope (not raw
// error on stderr) — add review #5 erblich für init.
func TestInitJSON_ConflictingModeFlagsShipsMinimalEnvelope(t *testing.T) {
	stub := &initUseCaseStub{}
	app := newAppWithInitStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "--yes", "--no-interactive", "init", "myproj"}, &stdout, &stderr)
	if !errors.Is(err, cli.ErrConflictingModeFlags) {
		t.Fatalf("expected ErrConflictingModeFlags, got: %v", err)
	}
	if got := cli.ExitCode(err); got != 2 {
		t.Errorf("cli.ExitCode: want 2, got %d", got)
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("init"),
		jsontestutil.WithExitCode(2),
	)
}

// =====================================================================
// Pflicht-Pins T0-(o)/(k)
// =====================================================================

// TestInitJSON_StdoutCleanlinessSinglJSONObject pins T0-(o) JSON-
// stdout-Cleanliness: stdout enthält EXAKT ein JSON-Object. ProgressPort
// wird durch SilenceProgress=true (CLI setzt es bei flags.JSON) in der
// Application-Schicht stillgelegt; bei einem Stub-UC hier ist
// ProgressPort gar nicht aktiv, aber der Pin sichert die Empty-Stream-
// Cleanliness und dokumentiert das Anti-Drift-Contract.
func TestInitJSON_StdoutCleanlinessSinglJSONObject(t *testing.T) {
	stub := &initUseCaseStub{
		resp: driving.InitProjectResponse{
			Project: domain.NewProject(mustProjectName(t, "myproj")),
			PlannedFiles: []driving.PlannedFile{
				{Path: "u-boot.yaml", Action: "create", NewContent: []byte("a\n")},
			},
		},
	}
	app := newAppWithInitStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"init", "myproj", "--dry-run", "--json"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	dec := json.NewDecoder(bytes.NewReader(stdout.Bytes()))
	var first map[string]any
	if err := dec.Decode(&first); err != nil {
		t.Fatalf("first JSON-decode failed: %v\nstdout=%s", err, stdout.String())
	}
	// Second decode must return io.EOF — kein zweites Object,
	// kein nachgehängter Text.
	var second any
	if err := dec.Decode(&second); err == nil {
		t.Errorf("stdout contains MORE than one JSON object — T0-(o) violation\nfirst=%v\nsecond=%v\nfull stdout=%s", first, second, stdout.String())
	}

	// Stub-UC empfängt SilenceProgress=true vom CLI-RunE (T0-(o)
	// SilenceProgress = flags.JSON).
	if !stub.lastReq.SilenceProgress {
		t.Errorf("CLI must set req.SilenceProgress=true in JSON-mode (T0-(o)), got false")
	}
}

// TestInitJSON_PreviewModeFromFlags pins that CLI-RunE correctly
// maps the --dry-run/--diff flag combinations to req.PreviewMode
// (Wahrheitstabelle T0-(b) — über previewModeFromFlags-Helper aus
// T1-B).
func TestInitJSON_PreviewModeFromFlags(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want driving.PreviewMode
	}{
		{"none", []string{"init", "myproj"}, driving.PreviewNone},
		{"dry-run", []string{"init", "myproj", "--dry-run"}, driving.PreviewDryRun},
		{"diff", []string{"init", "myproj", "--diff"}, driving.PreviewAndApply},
		{"dry-run + diff", []string{"init", "myproj", "--dry-run", "--diff"}, driving.PreviewDryRun},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stub := &initUseCaseStub{
				resp: driving.InitProjectResponse{Project: domain.NewProject(mustProjectName(t, "myproj"))},
			}
			app := newAppWithInitStub(stub)
			var stdout, stderr bytes.Buffer
			if err := app.Execute(context.Background(), tc.args, &stdout, &stderr); err != nil {
				t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
			}
			if stub.lastReq.PreviewMode != tc.want {
				t.Errorf("PreviewMode: want %v, got %v", tc.want, stub.lastReq.PreviewMode)
			}
		})
	}
}

// =====================================================================
// Human-Mode Tests
// =====================================================================

// TestInitJSON_PositionalNameWithSlashRejectsAsLHInit006 pinnt
// T0-(k) Path-Anchor (Review-Round-9 #4): wenn der User
// `u-boot init myproj/` (Trailing-Slash) oder `u-boot init ./myproj`
// eingibt, MUSS das Verhalten deterministisch sein — der Name wird
// vom CLI-RunE roh an req.Name gereicht, die Application validiert
// via domain.NewProjectName (Regex `^[a-z]([a-z0-9-]{0,61}[a-z0-9])?$`),
// und der Mapper klassifiziert als LH-FA-INIT-006 (Exit-Code 10).
// Ohne diesen Pin könnte eine zukünftige Normalisierungs-Phase
// silently 'myproj/' → 'myproj' machen und den User-Intent (Trailing-
// Slash falsch dahin interpretieren, dass eine Sub-Directory gemeint
// war) verschleiern.
func TestInitJSON_PositionalNameWithSlashRejectsAsLHInit006(t *testing.T) {
	// Stub gibt den raw Name an domain.NewProjectName weiter und
	// returnt den resulting error — emuliert das echte
	// resolveProjectName-Verhalten ohne den ganzen Application-
	// Initialisierungs-Path zu staging.
	cases := []struct {
		name string
		arg  string
	}{
		{"trailing slash", "myproj/"},
		{"leading dot-slash", "./myproj"},
		{"trailing dot", "myproj."},
		{"absolute path", "/abs/myproj"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stub := &validatingInitUseCaseStub{}
			app := newAppWithInitStub(stub)

			var stdout, stderr bytes.Buffer
			execErr := app.Execute(context.Background(),
				[]string{"--json", "init", tc.arg}, &stdout, &stderr)

			if !errors.Is(execErr, domain.ErrInvalidProjectName) {
				t.Fatalf("execute %q: want ErrInvalidProjectName, got %v", tc.arg, execErr)
			}
			if got := cli.ExitCode(execErr); got != 10 {
				t.Errorf("exit code: want 10 (LH-FA-INIT-006), got %d", got)
			}
			jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
				jsontestutil.WithCommand("init"),
				jsontestutil.WithExpectedCodes("LH-FA-INIT-006"),
				jsontestutil.WithExitCode(10),
			)
		})
	}
}

// validatingInitUseCaseStub emuliert das resolveProjectName-Verhalten
// der Application-Schicht: der raw req.Name wird durch
// domain.NewProjectName validiert und ein Reject wird als Sentinel
// returned. So bleibt der T0-(k)-Path-Anchor-Test CLI-acceptance-
// orientiert (Execute durch RunE), ohne die kompletten Init-
// Application-Wiring zu staging.
type validatingInitUseCaseStub struct{}

func (validatingInitUseCaseStub) Init(_ context.Context, req driving.InitProjectRequest) (driving.InitProjectResponse, error) {
	if _, err := domain.NewProjectName(req.Name); err != nil {
		return driving.InitProjectResponse{}, err
	}
	return driving.InitProjectResponse{
		Project: domain.NewProject(domain.ProjectName(req.Name)),
	}, nil
}

// TestInitHumanDryRun_ShowsWouldInitiate pins the human-mode `--dry-run`
// output: "Would initialize" / "Would create:" statt "Initialized"
// — analog add's "Would add"-Prefix (T5 printInitSummary dryRun-aware).
func TestInitHumanDryRun_ShowsWouldInitiate(t *testing.T) {
	stub := &initUseCaseStub{
		resp: driving.InitProjectResponse{
			Project: domain.NewProject(mustProjectName(t, "myproj")),
			Created: []string{"docker/", "scripts/", "README.md", "u-boot.yaml"},
		},
	}
	app := newAppWithInitStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"init", "myproj", "--dry-run"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "Would initialize") {
		t.Errorf("--dry-run must show 'Would initialize', got: %q", out)
	}
	if !strings.Contains(out, "Would create:") {
		t.Errorf("--dry-run must show 'Would create:', got: %q", out)
	}
	if strings.Contains(out, "Initialized u-boot project") {
		t.Errorf("--dry-run must NOT use past-tense 'Initialized', got: %q", out)
	}
}

// TestInitHumanDiff_RendersHunksOnStdout pins the human-mode --diff
// output: writeDiff emits @@-headers + + lines on stdout.
func TestInitHumanDiff_RendersHunksOnStdout(t *testing.T) {
	stub := &initUseCaseStub{
		resp: driving.InitProjectResponse{
			Project: domain.NewProject(mustProjectName(t, "myproj")),
			PlannedFiles: []driving.PlannedFile{
				{Path: "u-boot.yaml", Action: "create", NewContent: []byte("project:\n  name: myproj\n")},
			},
		},
	}
	app := newAppWithInitStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"init", "myproj", "--diff"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "--- u-boot.yaml (create)") {
		t.Errorf("--diff must show file header, got: %q", out)
	}
	if !strings.Contains(out, "+project:") {
		t.Errorf("--diff must show +project: line, got: %q", out)
	}
}

// =====================================================================
// Internal helpers
// =====================================================================

// unmarshalInitEnv parses the JSON output for further structural
// pins. (Duplicate-free vs. unmarshalEnv in add_acceptance_test.go
// would collide; the init-specific name keeps the per-test-file
// scope clear.)
func unmarshalInitEnv(t *testing.T, raw []byte) map[string]any {
	t.Helper()
	var env map[string]any
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatalf("unmarshal envelope: %v\nraw=%s", err, raw)
	}
	return env
}
