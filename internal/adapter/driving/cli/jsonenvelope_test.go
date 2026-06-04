package cli_test

import (
	"encoding/json"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
)

// TestMinimalEnvelope_AllOK pins the canonical All-OK shape from
// Lastenheft §1846-1852: only `status`, `command`, `diagnostics`,
// `exitCode` — no voll-schema fields.
func TestMinimalEnvelope_AllOK(t *testing.T) {
	raw, err := cli.MinimalEnvelopeForTest("doctor", "", nil, 0)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal failed: %v\nraw=%s", err, raw)
	}

	wantKeys := map[string]any{
		"status":      "ok",
		"command":     "doctor",
		"diagnostics": []any{},
		"exitCode":    float64(0),
	}
	for k, want := range wantKeys {
		if got[k] == nil && want != nil {
			t.Errorf("key %q missing in JSON: %s", k, raw)
			continue
		}
		// Slices need a deeper compare; for the rest a string repr is fine.
		switch want := want.(type) {
		case []any:
			arr, ok := got[k].([]any)
			if !ok || len(arr) != len(want) {
				t.Errorf("key %q: want empty array, got %#v", k, got[k])
			}
		default:
			if got[k] != want {
				t.Errorf("key %q: want %v, got %v", k, want, got[k])
			}
		}
	}

	forbidden := []string{"dryRun", "diff", "plannedFiles", "changes", "data"}
	for _, k := range forbidden {
		if _, present := got[k]; present {
			t.Errorf("minimal envelope must not contain %q (Spec §1841)", k)
		}
	}
}

// TestMinimalEnvelope_StatusCoupling pins Spec §447/§1837: status
// follows the highest diagnostics-level present.
func TestMinimalEnvelope_StatusCoupling(t *testing.T) {
	cases := []struct {
		name       string
		diags      []cli.DiagnosticItemForTest
		wantStatus string
	}{
		{"empty → ok", nil, "ok"},
		{"warn only → warn", []cli.DiagnosticItemForTest{
			{Level: "warn", Code: "uboot.yaml.valid", Message: "x"},
		}, "warn"},
		{"warn + error → error", []cli.DiagnosticItemForTest{
			{Level: "warn", Code: "uboot.yaml.valid", Message: "x"},
			{Level: "error", Code: "docker.installed", Message: "y"},
		}, "error"},
		{"error only → error", []cli.DiagnosticItemForTest{
			{Level: "error", Code: "docker.installed", Message: "y"},
		}, "error"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			raw, err := cli.MinimalEnvelopeForTest("doctor", "", tc.diags, 0)
			if err != nil {
				t.Fatalf("marshal failed: %v", err)
			}
			var got map[string]any
			if err := json.Unmarshal(raw, &got); err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}
			if got["status"] != tc.wantStatus {
				t.Errorf("status: want %q, got %v", tc.wantStatus, got["status"])
			}
		})
	}
}

// TestFullEnvelope_DryRunFalseDiffFalse_Serialised is the M1
// anti-drift pin from doctor.md T0-(d): newFullEnvelope with
// dryRun=false/diff=false MUST serialise as `"dryRun":false,
// "diff":false` (not omitted). If anyone refactors *bool → bool,
// this test breaks immediately.
func TestFullEnvelope_DryRunFalseDiffFalse_Serialised(t *testing.T) {
	raw, err := cli.FullEnvelopeForTest("add", "", false, false, nil, nil, nil, 0)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal failed: %v\nraw=%s", err, raw)
	}
	for _, key := range []string{"dryRun", "diff"} {
		v, present := got[key]
		if !present {
			t.Errorf("%q missing — Spec §326 required-set violated", key)
			continue
		}
		if v != false {
			t.Errorf("%q: want false, got %v", key, v)
		}
	}
}

// TestFullEnvelope_RequiredSet pins Spec §326: required keys are
// status, command, dryRun, diff, plannedFiles, changes,
// diagnostics, exitCode (all eight).
func TestFullEnvelope_RequiredSet(t *testing.T) {
	raw, err := cli.FullEnvelopeForTest(
		"add", "", true, false,
		[]cli.PlannedFileForTest{{Path: "compose.yaml", Action: "create"}},
		[]cli.ChangeEntryForTest{{Path: "compose.yaml", Count: 12}},
		nil, 0,
	)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	required := []string{"status", "command", "dryRun", "diff", "plannedFiles", "changes", "diagnostics", "exitCode"}
	for _, k := range required {
		if _, present := got[k]; !present {
			t.Errorf("required field %q missing in full envelope (Spec §326)", k)
		}
	}
}

// TestFullEnvelope_EmptyArraysSerialiseAsBracketBracket pins that
// PlannedFiles/Changes serialise as `[]` (not `null` or omitted)
// in the full envelope when empty. The full constructor normalises
// nil → empty slice so the modifying-pin always shows `[]`.
func TestFullEnvelope_EmptyArraysSerialiseAsBracketBracket(t *testing.T) {
	raw, err := cli.FullEnvelopeForTest("add", "", true, false, nil, nil, nil, 0)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	for _, k := range []string{"plannedFiles", "changes", "diagnostics"} {
		v, present := got[k]
		if !present {
			t.Errorf("%q missing", k)
			continue
		}
		arr, ok := v.([]any)
		if !ok {
			t.Errorf("%q: want []any, got %T", k, v)
			continue
		}
		if len(arr) != 0 {
			t.Errorf("%q: want empty array, got %#v", k, arr)
		}
	}
}

// TestMinimalEnvelope_SubcommandOmittedWhenEmpty pins that
// Subcommand is omitempty (Spec §1827: optional for non-grouped
// commands, mandatory for template/config).
func TestMinimalEnvelope_SubcommandOmittedWhenEmpty(t *testing.T) {
	raw, err := cli.MinimalEnvelopeForTest("doctor", "", nil, 0)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if _, present := got["subcommand"]; present {
		t.Errorf("subcommand should be omitted for ungrouped commands")
	}
}

// TestMinimalEnvelope_SubcommandPresentWhenSet pins that the
// template/config gruppierten Befehle carry subcommand verbatim.
func TestMinimalEnvelope_SubcommandPresentWhenSet(t *testing.T) {
	raw, err := cli.MinimalEnvelopeForTest("template", "list", nil, 0)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if got["subcommand"] != "list" {
		t.Errorf("subcommand: want \"list\", got %v", got["subcommand"])
	}
}
