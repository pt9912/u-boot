package jsontestutil_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli/jsontestutil"
)

// recordingT captures Errorf calls so the helper-tests can assert
// negative cases without actually failing the test.
type recordingT struct {
	testing.TB
	errors []string
}

func (r *recordingT) Errorf(format string, args ...any) {
	r.errors = append(r.errors, fmt.Sprintf(format, args...))
}

func (r *recordingT) Helper() {}

func newRecorder(t *testing.T) *recordingT {
	t.Helper()
	return &recordingT{TB: t}
}

func TestAssertMinimalEnvelope_AllOK(t *testing.T) {
	r := newRecorder(t)
	raw := []byte(`{"status":"ok","command":"doctor","diagnostics":[],"exitCode":0}`)
	jsontestutil.AssertMinimalEnvelope(r, raw,
		jsontestutil.WithCommand("doctor"),
		jsontestutil.WithExitCode(0),
	)
	if len(r.errors) != 0 {
		t.Errorf("expected no errors, got %d:\n  %s", len(r.errors), strings.Join(r.errors, "\n  "))
	}
}

func TestAssertMinimalEnvelope_RejectsLevelOK(t *testing.T) {
	r := newRecorder(t)
	raw := []byte(`{
		"status":"ok","command":"doctor",
		"diagnostics":[{"level":"ok","code":"docker.installed","message":"x"}],
		"exitCode":0
	}`)
	jsontestutil.AssertMinimalEnvelope(r, raw)
	if !containsSubstring(r.errors, "must be warn or error") {
		t.Errorf("expected level-ok reject, errors=%v", r.errors)
	}
}

func TestAssertMinimalEnvelope_RejectsLevelInfo(t *testing.T) {
	r := newRecorder(t)
	raw := []byte(`{
		"status":"ok","command":"doctor",
		"diagnostics":[{"level":"info","code":"docker.installed","message":"x"}],
		"exitCode":0
	}`)
	jsontestutil.AssertMinimalEnvelope(r, raw)
	if !containsSubstring(r.errors, "must be warn or error") {
		t.Errorf("expected level-info reject, errors=%v", r.errors)
	}
}

func TestAssertMinimalEnvelope_RejectsUndocumentedCode(t *testing.T) {
	r := newRecorder(t)
	raw := []byte(`{
		"status":"warn","command":"doctor",
		"diagnostics":[{"level":"warn","code":"made.up.code","message":"x"}],
		"exitCode":0
	}`)
	jsontestutil.AssertMinimalEnvelope(r, raw)
	if !containsSubstring(r.errors, "not in DefaultAllowedCodes") {
		t.Errorf("expected undocumented-code reject, errors=%v", r.errors)
	}
}

func TestAssertMinimalEnvelope_AcceptsLHCode(t *testing.T) {
	r := newRecorder(t)
	raw := []byte(`{
		"status":"warn","command":"add",
		"diagnostics":[{"level":"warn","code":"LH-FA-CLI-007","message":"x"}],
		"exitCode":0
	}`)
	jsontestutil.AssertMinimalEnvelope(r, raw)
	if len(r.errors) != 0 {
		t.Errorf("LH-codes must pass (Spec §445), got errors: %v", r.errors)
	}
}

func TestAssertMinimalEnvelope_RejectsFullSchemaFields(t *testing.T) {
	r := newRecorder(t)
	raw := []byte(`{
		"status":"ok","command":"doctor",
		"dryRun":false,"diff":false,"plannedFiles":[],"changes":[],
		"diagnostics":[],"exitCode":0
	}`)
	jsontestutil.AssertMinimalEnvelope(r, raw)
	for _, k := range []string{"dryRun", "diff", "plannedFiles", "changes"} {
		if !containsSubstring(r.errors, k) {
			t.Errorf("expected reject for full-schema field %q in minimal envelope, errors=%v", k, r.errors)
		}
	}
}

func TestAssertMinimalEnvelope_RejectsStatusDecoupling(t *testing.T) {
	r := newRecorder(t)
	raw := []byte(`{
		"status":"ok","command":"doctor",
		"diagnostics":[{"level":"error","code":"docker.installed","message":"x"}],
		"exitCode":11
	}`)
	jsontestutil.AssertMinimalEnvelope(r, raw)
	if !containsSubstring(r.errors, "decoupled from highest") {
		t.Errorf("expected status-coupling reject, errors=%v", r.errors)
	}
}

func TestAssertMinimalEnvelope_RejectsMissingRequired(t *testing.T) {
	r := newRecorder(t)
	raw := []byte(`{"command":"doctor","diagnostics":[],"exitCode":0}`)
	jsontestutil.AssertMinimalEnvelope(r, raw)
	if !containsSubstring(r.errors, "missing required field") {
		t.Errorf("expected missing-required reject, errors=%v", r.errors)
	}
}

func TestAssertMinimalEnvelope_RejectsBadSubcommandForTemplate(t *testing.T) {
	r := newRecorder(t)
	raw := []byte(`{"status":"ok","command":"template","diagnostics":[],"exitCode":0}`)
	jsontestutil.AssertMinimalEnvelope(r, raw)
	if !containsSubstring(r.errors, "subcommand required") {
		t.Errorf("expected subcommand-required reject for command=template, errors=%v", r.errors)
	}
}

func TestAssertMinimalEnvelope_WithExpectedCodesPin(t *testing.T) {
	r := newRecorder(t)
	raw := []byte(`{
		"status":"warn","command":"doctor",
		"diagnostics":[{"level":"warn","code":"docker.installed","message":"x"}],
		"exitCode":0
	}`)
	jsontestutil.AssertMinimalEnvelope(r, raw,
		jsontestutil.WithExpectedCodes("uboot.yaml.valid"),
	)
	if !containsSubstring(r.errors, "expected code") {
		t.Errorf("expected expected-code mismatch, errors=%v", r.errors)
	}
}

func TestAssertFullEnvelope_HappyPath(t *testing.T) {
	r := newRecorder(t)
	raw := []byte(`{
		"status":"ok","command":"add",
		"dryRun":true,"diff":false,
		"plannedFiles":[{"path":"compose.yaml","action":"create"}],
		"changes":[{"path":"compose.yaml","count":12}],
		"diagnostics":[],"exitCode":0
	}`)
	jsontestutil.AssertFullEnvelope(r, raw,
		jsontestutil.WithCommand("add"),
		jsontestutil.WithExitCode(0),
	)
	if len(r.errors) != 0 {
		t.Errorf("expected no errors, got %d: %v", len(r.errors), r.errors)
	}
}

func TestAssertFullEnvelope_RejectsMissingDryRun(t *testing.T) {
	r := newRecorder(t)
	raw := []byte(`{
		"status":"ok","command":"add",
		"diff":false,
		"plannedFiles":[],"changes":[],
		"diagnostics":[],"exitCode":0
	}`)
	jsontestutil.AssertFullEnvelope(r, raw)
	if !containsSubstring(r.errors, `"dryRun"`) {
		t.Errorf("expected missing dryRun reject, errors=%v", r.errors)
	}
}

func TestAssertFullEnvelope_RejectsBadPlannedAction(t *testing.T) {
	r := newRecorder(t)
	raw := []byte(`{
		"status":"ok","command":"add",
		"dryRun":true,"diff":false,
		"plannedFiles":[{"path":"x","action":"rename"}],
		"changes":[],
		"diagnostics":[],"exitCode":0
	}`)
	jsontestutil.AssertFullEnvelope(r, raw)
	if !containsSubstring(r.errors, "action") {
		t.Errorf("expected action-enum reject for 'rename', errors=%v", r.errors)
	}
}

func TestAssertFullEnvelope_RejectsNegativeCount(t *testing.T) {
	r := newRecorder(t)
	raw := []byte(`{
		"status":"ok","command":"add",
		"dryRun":true,"diff":false,
		"plannedFiles":[],
		"changes":[{"path":"x","count":-1}],
		"diagnostics":[],"exitCode":0
	}`)
	jsontestutil.AssertFullEnvelope(r, raw)
	if !containsSubstring(r.errors, "must be ≥ 0") {
		t.Errorf("expected negative-count reject, errors=%v", r.errors)
	}
}

// TestAssertFullEnvelope_AcceptsValidHunks pins the positive
// hunk-shape case from slice-v1-cli-json-dry-run-add T0-(l): three
// valid hunks (one pure addition with OldStart=0/OldLines=0, one
// middle-modify, one pure deletion with NewStart=0/NewLines=0) all
// pass without errors.
func TestAssertFullEnvelope_AcceptsValidHunks(t *testing.T) {
	r := newRecorder(t)
	raw := []byte(`{
		"status":"ok","command":"add",
		"dryRun":true,"diff":true,
		"plannedFiles":[{
			"path":"compose.yaml","action":"create",
			"hunks":[
				{"oldStart":0,"oldLines":0,"newStart":1,"newLines":3,"content":"+a\n+b\n+c\n"},
				{"oldStart":10,"oldLines":1,"newStart":13,"newLines":1,"content":"-x\n+X\n"},
				{"oldStart":20,"oldLines":2,"newStart":0,"newLines":0,"content":"-y\n-z\n"}
			]
		}],
		"changes":[{"path":"compose.yaml","count":4}],
		"diagnostics":[],"exitCode":0
	}`)
	jsontestutil.AssertFullEnvelope(r, raw,
		jsontestutil.WithCommand("add"),
	)
	if len(r.errors) != 0 {
		t.Errorf("expected no errors for valid hunks, got %d: %v", len(r.errors), r.errors)
	}
}

// TestAssertFullEnvelope_RejectsHunkFieldNameDrift is the negative
// Pin from T0-(l): a renamed field (`offset` instead of `oldStart`)
// must fail the hunk-shape check. This Drift-Anker catches future
// refactors that accidentally rename a hunk field.
func TestAssertFullEnvelope_RejectsHunkFieldNameDrift(t *testing.T) {
	r := newRecorder(t)
	raw := []byte(`{
		"status":"ok","command":"add",
		"dryRun":true,"diff":true,
		"plannedFiles":[{
			"path":"compose.yaml","action":"create",
			"hunks":[{"offset":1,"oldLines":0,"newStart":1,"newLines":3,"content":"+a\n+b\n+c\n"}]
		}],
		"changes":[{"path":"compose.yaml","count":3}],
		"diagnostics":[],"exitCode":0
	}`)
	jsontestutil.AssertFullEnvelope(r, raw)
	if !containsSubstring(r.errors, `"oldStart"`) {
		t.Errorf("expected oldStart-field-drift reject, errors=%v", r.errors)
	}
}

// TestAssertFullEnvelope_RejectsHunkStartZeroWithLinesPositive pins
// the 1-based-coordinate invariant: when *Lines > 0 the
// corresponding *Start MUST be ≥ 1 (T0-(l) Hunk-Schema).
func TestAssertFullEnvelope_RejectsHunkStartZeroWithLinesPositive(t *testing.T) {
	r := newRecorder(t)
	raw := []byte(`{
		"status":"ok","command":"add",
		"dryRun":true,"diff":true,
		"plannedFiles":[{
			"path":"compose.yaml","action":"modify",
			"hunks":[{"oldStart":0,"oldLines":2,"newStart":1,"newLines":2,"content":"-a\n-b\n+A\n+B\n"}]
		}],
		"changes":[{"path":"compose.yaml","count":2}],
		"diagnostics":[],"exitCode":0
	}`)
	jsontestutil.AssertFullEnvelope(r, raw)
	if !containsSubstring(r.errors, "≥ 1") {
		t.Errorf("expected start-must-be-≥1 reject when lines>0, errors=%v", r.errors)
	}
}

// TestAssertFullEnvelope_HunkAbsenceIsOK pins that the hunks field
// is optional on plannedFile entries — Spec §326 lists it as part of
// the --diff --json subset only, and omitempty omission must not
// trigger checkHunks at all.
func TestAssertFullEnvelope_HunkAbsenceIsOK(t *testing.T) {
	r := newRecorder(t)
	raw := []byte(`{
		"status":"ok","command":"add",
		"dryRun":true,"diff":false,
		"plannedFiles":[{"path":"compose.yaml","action":"create"}],
		"changes":[{"path":"compose.yaml","count":12}],
		"diagnostics":[],"exitCode":0
	}`)
	jsontestutil.AssertFullEnvelope(r, raw,
		jsontestutil.WithCommand("add"),
	)
	if len(r.errors) != 0 {
		t.Errorf("hunks-absence is allowed; got errors: %v", r.errors)
	}
}

func TestDefaultAllowedCodes_NotEmpty(t *testing.T) {
	codes := jsontestutil.DefaultAllowedCodes()
	if len(codes) == 0 {
		t.Errorf("DefaultAllowedCodes must not be empty (Spec §1835)")
	}
}

func containsSubstring(errors []string, want string) bool {
	for _, e := range errors {
		if strings.Contains(e, want) {
			return true
		}
	}
	return false
}
