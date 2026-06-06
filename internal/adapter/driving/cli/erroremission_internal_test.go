package cli

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// TestMapAddErrorToDiagnostic_AllCases covers every switch-case in
// add.go's mapAddErrorToDiagnostic to lift the coverage from the
// 30% spot (only the default-fallback case was exercised through
// add_json_test indirectly). Each sentinel-class gets its own
// expected LH-code; the default path catches anything unknown.
//
// T0-(f) Switch-Order-Pflicht: ErrAddFileSystem MUST be checked
// first — multi-`%w` wraps that happen to include both
// ErrAddFileSystem AND a fachlich sentinel route to LH-NFA-REL-003
// (FS-class), not to LH-FA-ADD-005 (fachlich-class). The
// MultiWrap-case at the end of the table pins this.
func TestMapAddErrorToDiagnostic_AllCases(t *testing.T) {
	cases := []struct {
		name    string
		err     error
		wantLH  string
		wantLvl string
	}{
		{"ErrAddFileSystem", driving.ErrAddFileSystem, "LH-NFA-REL-003", "error"},
		{"ErrProjectNotInitialized", driving.ErrProjectNotInitialized, "LH-FA-ADD-001", "error"},
		{"ErrServiceUnsupported", driving.ErrServiceUnsupported, "LH-FA-ADD-002", "error"},
		{"ErrServiceInconsistent", driving.ErrServiceInconsistent, "LH-FA-ADD-005", "error"},
		{"ErrDependenciesRequired", driving.ErrDependenciesRequired, "LH-FA-ADD-006", "error"},
		{"ErrInvalidServiceName", domain.ErrInvalidServiceName, "LH-FA-INIT-006", "error"},
		{"ErrFileExists", driving.ErrFileExists, "LH-FA-INIT-004", "error"},
		{"ErrProjectExists", driving.ErrProjectExists, "LH-FA-INIT-004", "error"},
		{"ErrBackupSuffixExhausted", driving.ErrBackupSuffixExhausted, "LH-NFA-REL-003", "error"},
		{"ErrBackupSourceMissing", driving.ErrBackupSourceMissing, "LH-NFA-REL-003", "error"},
		{"unknown sentinel → default LH-FA-CLI-006", errors.New("unknown"), "LH-FA-CLI-006", "error"},
		// T0-(f) Switch-Order pin: multi-`%w` mit ErrAddFileSystem +
		// ErrProjectNotInitialized muss als FS-class klassifiziert
		// werden (FS-first checks). Ohne FS-first würde der wrap als
		// fachlich (LH-FA-ADD-001) durchgehen — Exit-10 statt Exit-14.
		{
			"multi-%w: ErrAddFileSystem + ErrProjectNotInitialized → FS-class wins",
			fmt.Errorf("write x: %w: %w", driving.ErrAddFileSystem, driving.ErrProjectNotInitialized),
			"LH-NFA-REL-003", "error",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			diag := mapAddErrorToDiagnostic(tc.err)
			if diag.Code != tc.wantLH {
				t.Errorf("code: want %q, got %q", tc.wantLH, diag.Code)
			}
			if diag.Level != tc.wantLvl {
				t.Errorf("level: want %q, got %q", tc.wantLvl, diag.Level)
			}
			if diag.Message == "" {
				t.Errorf("message must not be empty")
			}
		})
	}
}

// TestMapInitErrorToDiagnostic_AllCases mirrors the Add-test for
// init's mapInitErrorToDiagnostic so each switch-case is covered.
// Pins T0-(f) Switch-Order: ErrInitFileSystem FIRST + Backup-FS-
// sentinels also routed to LH-NFA-REL-003.
func TestMapInitErrorToDiagnostic_AllCases(t *testing.T) {
	cases := []struct {
		name   string
		err    error
		wantLH string
	}{
		{"ErrInitFileSystem", driving.ErrInitFileSystem, "LH-NFA-REL-003"},
		{"ErrBackupSuffixExhausted", driving.ErrBackupSuffixExhausted, "LH-NFA-REL-003"},
		{"ErrBackupSourceMissing", driving.ErrBackupSourceMissing, "LH-NFA-REL-003"},
		{"ErrTemplateConflictsWithFlag", driving.ErrTemplateConflictsWithFlag, "LH-FA-CLI-006"},
		{"ErrConfirmationRequired", driving.ErrConfirmationRequired, "LH-FA-INIT-005"},
		{"ErrForceRequiresBackup", driving.ErrForceRequiresBackup, "LH-FA-INIT-005"},
		{"ErrBackupUnsupportedKind", driving.ErrBackupUnsupportedKind, "LH-FA-INIT-005"},
		{"ErrProjectExists", driving.ErrProjectExists, "LH-FA-INIT-004"},
		{"ErrFileExists", driving.ErrFileExists, "LH-FA-INIT-004"},
		{"ErrInvalidProjectName", domain.ErrInvalidProjectName, "LH-FA-INIT-006"},
		// Review-Round-9 #1: ErrInvalidFeatureSource (LH-FA-DEV-003)
		// stammt aus `init --allow-external-feature-sources` ohne
		// `--devcontainer`. Vorher fiel der Wert zur default-CLI-006-
		// Klasse während der Exit-Code via isConfigValidationError
		// bereits Code-10 lieferte — Envelope-Code/Exit-Klassen
		// drifteten auseinander.
		{"ErrInvalidFeatureSource", domain.ErrInvalidFeatureSource, "LH-FA-DEV-003"},
		{"unknown → default LH-FA-CLI-006", errors.New("unknown"), "LH-FA-CLI-006"},
		// Multi-`%w` Switch-Order-Pin (analog Add-Test): FS-first
		// klassifiziert auch wenn der wrap ein fachlich Sentinel
		// enthält.
		{
			"multi-%w: ErrInitFileSystem + ErrProjectExists → FS-class wins",
			fmt.Errorf("init: write x: %w: %w", driving.ErrInitFileSystem, driving.ErrProjectExists),
			"LH-NFA-REL-003",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			diag := mapInitErrorToDiagnostic(tc.err)
			if diag.Code != tc.wantLH {
				t.Errorf("code: want %q, got %q", tc.wantLH, diag.Code)
			}
		})
	}
}

// TestWriteDiff_EdgeCases covers the three non-trivial writeDiff
// paths that the acceptance-tests only exercise indirectly:
// (a) binary content → "(binary content — diff suppressed)" hint;
// (b) content-identical modify → "(no changes)" hint;
// (c) multi-file output → blank-line separator between files.
func TestWriteDiff_EdgeCases(t *testing.T) {
	t.Run("binary content → suppressed hint", func(t *testing.T) {
		var buf bytes.Buffer
		planned := []driving.PlannedFile{
			{Path: "blob.bin", Action: "create", NewContent: []byte{0xff, 0xfe, 0xfd}},
		}
		if err := writeDiff(&buf, planned); err != nil {
			t.Fatalf("writeDiff: %v", err)
		}
		out := buf.String()
		if !strings.Contains(out, "--- blob.bin (create)") {
			t.Errorf("missing file header, got: %q", out)
		}
		if !strings.Contains(out, "binary content") {
			t.Errorf("missing binary-suppressed hint, got: %q", out)
		}
	})

	t.Run("content-identical modify → (no changes) hint", func(t *testing.T) {
		var buf bytes.Buffer
		body := []byte("same\n")
		planned := []driving.PlannedFile{
			{Path: "file.txt", Action: "modify", OldContent: body, NewContent: body},
		}
		if err := writeDiff(&buf, planned); err != nil {
			t.Fatalf("writeDiff: %v", err)
		}
		if !strings.Contains(buf.String(), "(no changes)") {
			t.Errorf("missing (no changes) hint, got: %q", buf.String())
		}
	})

	t.Run("multi-file → blank line separator", func(t *testing.T) {
		var buf bytes.Buffer
		planned := []driving.PlannedFile{
			{Path: "a.txt", Action: "create", NewContent: []byte("a\n")},
			{Path: "b.txt", Action: "create", NewContent: []byte("b\n")},
		}
		if err := writeDiff(&buf, planned); err != nil {
			t.Fatalf("writeDiff: %v", err)
		}
		out := buf.String()
		// The separator-line between files is a bare \n; find both
		// file headers and ensure there's at least one blank line
		// (\n\n) between them.
		idxA := strings.Index(out, "--- a.txt")
		idxB := strings.Index(out, "--- b.txt")
		if idxA == -1 || idxB == -1 || idxA > idxB {
			t.Fatalf("file headers missing or out-of-order, got: %q", out)
		}
		between := out[idxA:idxB]
		if !strings.Contains(between, "\n\n") {
			t.Errorf("expected blank-line separator between a and b, got: %q", between)
		}
	})
}

// TestApplyJSONRejectGate_Branches covers the early-return
// branches that the existing acceptance-tests do not exercise:
// jsonFlag=false (pass-through), cmd==nil (defensive), help
// subcommand, and __complete (Cobra-internal shell-completion
// escape hatch). The Allowlist-hit + reject paths are exercised
// in jsonallowlist_test.go via the full Execute() flow.
func TestApplyJSONRejectGate_Branches(t *testing.T) {
	t.Run("jsonFlag=false → no-op", func(t *testing.T) {
		if err := applyJSONRejectGate(nil, false); err != nil {
			t.Errorf("want nil, got %v", err)
		}
	})
	t.Run("cmd==nil → defensive no-op", func(t *testing.T) {
		if err := applyJSONRejectGate(nil, true); err != nil {
			t.Errorf("want nil, got %v", err)
		}
	})
	t.Run("cmd.Name()==help → escape hatch", func(t *testing.T) {
		c := &cobra.Command{Use: "help"}
		if err := applyJSONRejectGate(c, true); err != nil {
			t.Errorf("help cmd must pass through, got %v", err)
		}
	})
	t.Run("cmd.Name()==__complete → escape hatch", func(t *testing.T) {
		c := &cobra.Command{Use: "__complete"}
		if err := applyJSONRejectGate(c, true); err != nil {
			t.Errorf("__complete cmd must pass through, got %v", err)
		}
	})
}

// TestHelpRequested_NoFlag covers the nil-flag branch of
// helpRequested — a Cobra command that has not had a --help
// persistent flag registered. The applyJSONRejectGate caller relies
// on this defensive fallback so a malformed command tree does not
// crash the gate.
func TestHelpRequested_NoFlag(t *testing.T) {
	c := &cobra.Command{Use: "x"}
	if helpRequested(c) {
		t.Errorf("cmd without --help flag must return false")
	}
}

// TestJSONSliceSuffix_DefaultCases covers the jsonSliceSuffix
// branches not exercised through the public TreeWalk-paths: the
// no-root-prefix (defensive), the empty-first-segment after the
// prefix (orphan "u-boot " with trailing space), and the recognised
// "up"/"down" → "up-down" collapse.
func TestJSONSliceSuffix_DefaultCases(t *testing.T) {
	cases := []struct {
		path string
		want string
	}{
		{"weird/no-root", "unknown"},
		{"u-boot ", "unknown"},
		{"u-boot up", "up-down"},
		{"u-boot down", "up-down"},
		{"u-boot logs", "logs"},
	}
	for _, tc := range cases {
		got := jsonSliceSuffix(tc.path)
		if got != tc.want {
			t.Errorf("path %q: want %q, got %q", tc.path, tc.want, got)
		}
	}
}

// TestStatusFromDiagnostics_AllBranches pins Spec §447 / §1837 —
// error wins over warn, warn wins over ok, ok = empty/info-only.
func TestStatusFromDiagnostics_AllBranches(t *testing.T) {
	cases := []struct {
		name  string
		diags []diagnosticItem
		want  string
	}{
		{"empty → ok", []diagnosticItem{}, "ok"},
		{"warn-only → warn", []diagnosticItem{{Level: "warn"}}, "warn"},
		{"error-only → error", []diagnosticItem{{Level: "error"}}, "error"},
		{"warn+error → error", []diagnosticItem{{Level: "warn"}, {Level: "error"}}, "error"},
		{"info-only → ok", []diagnosticItem{{Level: "info"}}, "ok"},
		{"error first → short-circuit error", []diagnosticItem{{Level: "error"}, {Level: "warn"}}, "error"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := statusFromDiagnostics(tc.diags)
			if got != tc.want {
				t.Errorf("want %q, got %q", tc.want, got)
			}
		})
	}
}

// TestNewMinimalEnvelope_NilDiagsAreReplaced covers the nil-guard
// branch of newMinimalEnvelope — analog to newFullEnvelope's.
func TestNewMinimalEnvelope_NilDiagsAreReplaced(t *testing.T) {
	env := newMinimalEnvelope("u-boot doctor", "", nil, 0)
	if env.Diagnostics == nil {
		t.Errorf("nil diagnostics must be replaced with []")
	}
	// Minimal envelope must NOT carry the full-fields.
	if env.PlannedFiles != nil || env.Changes != nil || env.DryRun != nil || env.Diff != nil {
		t.Errorf("minimal envelope must omit full-fields, got planned=%v changes=%v dryRun=%v diff=%v",
			env.PlannedFiles, env.Changes, env.DryRun, env.Diff)
	}
}

// TestNewFullEnvelope_NilSlicesAreReplaced covers the three nil-
// guard branches in newFullEnvelope: nil diags / nil planned /
// nil changes must all be replaced with non-nil empty slices so the
// rendered JSON ships `"diagnostics": []`, `"plannedFiles": []`,
// `"changes": []` (Spec-Required-Set, never `null`).
func TestNewFullEnvelope_NilSlicesAreReplaced(t *testing.T) {
	env := newFullEnvelope("u-boot init", "", false, false, nil, nil, nil, nil, 0)
	if env.PlannedFiles == nil {
		t.Errorf("nil planned must be replaced (got nil pointer)")
	} else if *env.PlannedFiles == nil {
		t.Errorf("planned slice must be non-nil empty, got nil slice")
	}
	if env.Changes == nil {
		t.Errorf("nil changes must be replaced (got nil pointer)")
	} else if *env.Changes == nil {
		t.Errorf("changes slice must be non-nil empty, got nil slice")
	}
	if env.Diagnostics == nil {
		t.Errorf("nil diagnostics must be replaced with []")
	}
	if env.DryRun == nil || *env.DryRun {
		t.Errorf("DryRun must point to false")
	}
	if env.Diff == nil || *env.Diff {
		t.Errorf("Diff must point to false")
	}
}

// TestAddSummaryHeader_AllBranches covers the four addSummaryHeader
// state-combinations: already-active, repair-existing,
// register/reactivate, plus their dry-run variants. This pins the
// human-mode lead-in text for every PriorState×Changed combination.
func TestAddSummaryHeader_AllBranches(t *testing.T) {
	svc, _ := domainServiceName(t, "postgres")
	cases := []struct {
		name   string
		resp   driving.AddServiceResponse
		dryRun bool
		want   string
	}{
		{
			"already active, no changes",
			driving.AddServiceResponse{ServiceName: svc, PriorState: domain.ServiceStateActive, State: domain.ServiceStateActive},
			false,
			`is already active`,
		},
		{
			"repair (same state, has changes)",
			driving.AddServiceResponse{ServiceName: svc, PriorState: domain.ServiceStateActive, State: domain.ServiceStateActive, Changed: []string{"compose.yaml"}},
			false,
			"Repaired service",
		},
		{
			"repair dry-run",
			driving.AddServiceResponse{ServiceName: svc, PriorState: domain.ServiceStateActive, State: domain.ServiceStateActive, Changed: []string{"compose.yaml"}},
			true,
			"Would repair service",
		},
		{
			"register (state change)",
			driving.AddServiceResponse{ServiceName: svc, PriorState: domain.ServiceStateUnregistered, State: domain.ServiceStateActive},
			false,
			"Added service",
		},
		{
			"register dry-run",
			driving.AddServiceResponse{ServiceName: svc, PriorState: domain.ServiceStateUnregistered, State: domain.ServiceStateActive},
			true,
			"Would add service",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := addSummaryHeader(tc.resp, tc.dryRun)
			if !strings.Contains(got, tc.want) {
				t.Errorf("header missing %q in: %q", tc.want, got)
			}
		})
	}
}

// TestPrintAddSummary_AlreadyActive covers the early-return path in
// printAddSummary where the header carries the full message and no
// Changed lines are emitted.
func TestPrintAddSummary_AlreadyActive(t *testing.T) {
	svc, _ := domainServiceName(t, "postgres")
	var buf bytes.Buffer
	resp := driving.AddServiceResponse{ServiceName: svc, PriorState: domain.ServiceStateActive, State: domain.ServiceStateActive}
	if err := printAddSummary(&buf, resp, false); err != nil {
		t.Fatalf("printAddSummary: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "already active") {
		t.Errorf("missing already-active line, got: %q", out)
	}
	if strings.Contains(out, "\n  - ") {
		t.Errorf("must not emit changed-lines for no-op, got: %q", out)
	}
}

// TestPrintAddSummary_WithChanges covers the loop-over-Changed
// path.
func TestPrintAddSummary_WithChanges(t *testing.T) {
	svc, _ := domainServiceName(t, "postgres")
	var buf bytes.Buffer
	resp := driving.AddServiceResponse{
		ServiceName: svc,
		PriorState:  domain.ServiceStateUnregistered,
		State:       domain.ServiceStateActive,
		Changed:     []string{"u-boot.yaml", "compose.yaml", ".env.example"},
	}
	if err := printAddSummary(&buf, resp, false); err != nil {
		t.Fatalf("printAddSummary: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"Added service", "u-boot.yaml", "compose.yaml", ".env.example"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in: %q", want, out)
		}
	}
}

// domainServiceName is a tiny helper so the printAddSummary tests
// don't repeat the domain.NewServiceName boilerplate.
func domainServiceName(t *testing.T, name string) (domain.ServiceName, error) {
	t.Helper()
	svc, err := domain.NewServiceName(name)
	if err != nil {
		t.Fatalf("NewServiceName(%q): %v", name, err)
	}
	return svc, nil
}

// TestPrintInitSummary_AllPaths covers the four printInitSummary
// branches not reached by the acceptance tests: dryRun-prefix +
// the Backups block (replace / overwrite paths populate Backups —
// the JSON-acceptance fixtures only exercise create-mode init).
func TestPrintInitSummary_AllPaths(t *testing.T) {
	t.Run("non-dry, no backups", func(t *testing.T) {
		var buf bytes.Buffer
		resp := driving.InitProjectResponse{
			Project: domain.Project{Name: "demo", SchemaVersion: 1},
			Created: []string{"u-boot.yaml", ".devcontainer/devcontainer.json"},
		}
		if err := printInitSummary(&buf, resp, false); err != nil {
			t.Fatalf("printInitSummary: %v", err)
		}
		out := buf.String()
		for _, want := range []string{`Initialized u-boot project "demo"`, "Created:", "u-boot.yaml"} {
			if !strings.Contains(out, want) {
				t.Errorf("missing %q in: %s", want, out)
			}
		}
		if strings.Contains(out, "Backups:") {
			t.Errorf("no backups should emit no Backups section, got: %s", out)
		}
	})
	t.Run("dry-run, no backups", func(t *testing.T) {
		var buf bytes.Buffer
		resp := driving.InitProjectResponse{
			Project: domain.Project{Name: "demo", SchemaVersion: 1},
			Created: []string{"u-boot.yaml"},
		}
		if err := printInitSummary(&buf, resp, true); err != nil {
			t.Fatalf("printInitSummary: %v", err)
		}
		out := buf.String()
		for _, want := range []string{`Would initialize u-boot project "demo"`, "Would create:"} {
			if !strings.Contains(out, want) {
				t.Errorf("missing %q in: %s", want, out)
			}
		}
	})
	t.Run("with Backups", func(t *testing.T) {
		var buf bytes.Buffer
		resp := driving.InitProjectResponse{
			Project: domain.Project{Name: "demo", SchemaVersion: 1},
			Created: []string{"u-boot.yaml"},
			Backups: []driving.BackupAction{
				{Original: "u-boot.yaml", Backup: "u-boot.yaml.bak.1"},
				{Original: ".env", Backup: ".env.bak.1"},
			},
		}
		if err := printInitSummary(&buf, resp, false); err != nil {
			t.Fatalf("printInitSummary: %v", err)
		}
		out := buf.String()
		for _, want := range []string{
			"Backups:",
			"u-boot.yaml → u-boot.yaml.bak.1",
			".env → .env.bak.1",
		} {
			if !strings.Contains(out, want) {
				t.Errorf("missing %q in: %s", want, out)
			}
		}
	})
}

// TestComputeChangeCountAndHunks_AllActions covers the four
// branches in computeChangeCountAndHunks: delete (short-circuit),
// binary non-delete (CountBytesDiff), create (CountLines + hunks),
// modify (CountAdditions + hunks). The default-fallback branch is
// unreachable per spec (action ∈ {create, modify, delete}) but we
// also exercise an unknown action so the switch's default is
// observed by coverage.
func TestComputeChangeCountAndHunks_AllActions(t *testing.T) {
	t.Run("delete → (0, nil) regardless of content", func(t *testing.T) {
		pf := driving.PlannedFile{Path: "x", Action: "delete", OldContent: []byte("anything"), NewContent: nil}
		count, hunks := computeChangeCountAndHunks(pf)
		if count != 0 || hunks != nil {
			t.Errorf("delete: want (0, nil), got (%d, %v)", count, hunks)
		}
	})
	t.Run("binary non-delete → CountBytesDiff, nil hunks", func(t *testing.T) {
		pf := driving.PlannedFile{Path: "blob.bin", Action: "create", NewContent: []byte{0xff, 0xfe, 0xfd, 0xfc}}
		count, hunks := computeChangeCountAndHunks(pf)
		if count <= 0 {
			t.Errorf("binary create: count > 0 expected, got %d", count)
		}
		if hunks != nil {
			t.Errorf("binary: hunks must be nil, got %v", hunks)
		}
	})
	t.Run("create text → CountLines + hunks", func(t *testing.T) {
		pf := driving.PlannedFile{Path: "a.txt", Action: "create", NewContent: []byte("one\ntwo\nthree\n")}
		count, hunks := computeChangeCountAndHunks(pf)
		if count != 3 {
			t.Errorf("create 3-line: count want 3, got %d", count)
		}
		if len(hunks) == 0 {
			t.Errorf("create: expected at least one hunk")
		}
	})
	t.Run("modify text → CountAdditions + hunks", func(t *testing.T) {
		pf := driving.PlannedFile{
			Path:       "f.txt",
			Action:     "modify",
			OldContent: []byte("a\nb\n"),
			NewContent: []byte("a\nb\nc\nd\n"),
		}
		count, hunks := computeChangeCountAndHunks(pf)
		if count != 2 {
			t.Errorf("modify 2-add: count want 2, got %d", count)
		}
		if len(hunks) == 0 {
			t.Errorf("modify: expected at least one hunk")
		}
	})
	t.Run("unknown action → default fallback (CountLines parity)", func(t *testing.T) {
		pf := driving.PlannedFile{Path: "u.txt", Action: "weird", NewContent: []byte("x\ny\n")}
		count, _ := computeChangeCountAndHunks(pf)
		if count != 2 {
			t.Errorf("unknown action: want 2 (CountLines), got %d", count)
		}
	})
}

// TestToCLIHunks_NilAndPopulated covers the two branches of
// toCLIHunks: nil/empty input → nil out (no allocation), populated
// → faithful field-copy.
func TestToCLIHunks_NilAndPopulated(t *testing.T) {
	if got := toCLIHunks(nil); got != nil {
		t.Errorf("nil → nil expected, got %v", got)
	}
	if got := toCLIHunks([]driving.Hunk{}); got != nil {
		t.Errorf("empty → nil expected, got %v", got)
	}
	src := []driving.Hunk{
		{OldStart: 1, OldLines: 0, NewStart: 1, NewLines: 3, Content: "+a\n+b\n+c\n"},
	}
	out := toCLIHunks(src)
	if len(out) != 1 {
		t.Fatalf("want 1 hunk, got %d", len(out))
	}
	if out[0].Content != src[0].Content || out[0].NewLines != 3 {
		t.Errorf("hunk fields lost: %+v vs %+v", out[0], src[0])
	}
}

// TestLastPlannedPath_EdgeCases covers the two branches that
// add_acceptance_test does not exercise directly.
func TestLastPlannedPath_EdgeCases(t *testing.T) {
	t.Run("empty → empty string", func(t *testing.T) {
		if got := lastPlannedPath(nil); got != "" {
			t.Errorf("want empty string for nil planned, got %q", got)
		}
		if got := lastPlannedPath([]driving.PlannedFile{}); got != "" {
			t.Errorf("want empty string for empty slice, got %q", got)
		}
	})
	t.Run("returns tail path", func(t *testing.T) {
		planned := []driving.PlannedFile{
			{Path: "first"},
			{Path: "middle"},
			{Path: "last"},
		}
		if got := lastPlannedPath(planned); got != "last" {
			t.Errorf("want \"last\", got %q", got)
		}
	})
}
