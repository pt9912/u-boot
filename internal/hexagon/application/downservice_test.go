package application_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// downFixture bundles the DownService under test with all fakes,
// pre-populated with u-boot.yaml + compose.yaml so happy-path tests
// don't repeat boilerplate.
type downFixture struct {
	svc       *application.DownService
	fs        *fakeFS
	engine    *fakeDockerEngine
	confirmer *fakeConfirmer
}

func newDownFixture(t *testing.T) *downFixture {
	t.Helper()
	fs := newFakeFS()
	if err := fs.WriteFile("/proj/u-boot.yaml", []byte("schemaVersion: 1\nproject:\n  name: demo\n"), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
	if err := fs.WriteFile("/proj/compose.yaml", []byte("services:\n  postgres:\n    image: postgres:16-alpine\n"), 0o644); err != nil {
		t.Fatalf("seed compose.yaml: %v", err)
	}
	engine := newFakeDockerEngine()
	engine.scriptDown(nil) // default: success
	confirmer := &fakeConfirmer{}
	svc := application.NewDownService(fs, engine, confirmer, nil)
	return &downFixture{svc: svc, fs: fs, engine: engine, confirmer: confirmer}
}

func TestDownService_BaseDirEmpty_ReturnsValidationError(t *testing.T) {
	t.Parallel()
	f := newDownFixture(t)
	_, err := f.svc.Down(context.Background(), driving.DownRequest{BaseDir: ""})
	if err == nil {
		t.Fatal("expected error for empty BaseDir")
	}
	if errors.Is(err, driving.ErrProjectNotInitialized) ||
		errors.Is(err, driving.ErrComposeFileMissing) ||
		errors.Is(err, driving.ErrConfirmationRequired) {
		t.Errorf("validation error should not match a use-case sentinel: %v", err)
	}
}

func TestDownService_MissingUbootYAML_ReturnsErrProjectNotInitialized(t *testing.T) {
	t.Parallel()
	f := newDownFixture(t)
	if err := f.fs.RemoveAll("/proj/u-boot.yaml"); err != nil {
		t.Fatalf("setup: %v", err)
	}
	_, err := f.svc.Down(context.Background(), driving.DownRequest{BaseDir: "/proj"})
	if !errors.Is(err, driving.ErrProjectNotInitialized) {
		t.Errorf("expected ErrProjectNotInitialized, got: %v", err)
	}
	// Pin: failed pre-check must NOT have reached the engine.
	if f.engine.downCallCount != 0 {
		t.Errorf("ComposeDown was called despite pre-check failure")
	}
}

func TestDownService_MissingComposeYAML_ReturnsErrComposeFileMissing(t *testing.T) {
	t.Parallel()
	f := newDownFixture(t)
	if err := f.fs.RemoveAll("/proj/compose.yaml"); err != nil {
		t.Fatalf("setup: %v", err)
	}
	_, err := f.svc.Down(context.Background(), driving.DownRequest{BaseDir: "/proj"})
	if !errors.Is(err, driving.ErrComposeFileMissing) {
		t.Errorf("expected ErrComposeFileMissing, got: %v", err)
	}
	if f.engine.downCallCount != 0 {
		t.Errorf("ComposeDown was called despite missing compose.yaml")
	}
}

// --- §T5 truth-table tests (4 rows × variants) ---

func TestDownService_TruthTable_Row1_NoVolumes_NoConfirmer(t *testing.T) {
	t.Parallel()
	// Row 1: RemoveVolumes=false → proceed straight to engine.
	f := newDownFixture(t)
	resp, err := f.svc.Down(context.Background(), driving.DownRequest{BaseDir: "/proj", RemoveVolumes: false})
	if err != nil {
		t.Fatalf("Down: %v", err)
	}
	if resp.RemovedVolumes {
		t.Errorf("RemovedVolumes = true, want false")
	}
	if len(f.confirmer.removeVolumesCalls) != 0 {
		t.Errorf("Confirmer was called for non-destructive down")
	}
	if f.engine.downCallCount != 1 {
		t.Errorf("ComposeDown called %d times, want 1", f.engine.downCallCount)
	}
	if f.engine.downOptions.RemoveVolumes {
		t.Errorf("Engine called with RemoveVolumes=true, want false")
	}
}

func TestDownService_TruthTable_Row2_VolumesYesAssume_NoConfirmer(t *testing.T) {
	t.Parallel()
	// Row 2: RemoveVolumes=true, AssumeYes=true → no confirmer.
	// NonInteractive does not matter — two sub-cases.
	subCases := []struct {
		name           string
		nonInteractive bool
	}{
		{"NonInteractive=false", false},
		{"NonInteractive=true", true},
	}
	for _, sc := range subCases {
		sc := sc
		t.Run(sc.name, func(t *testing.T) {
			t.Parallel()
			f := newDownFixture(t)
			resp, err := f.svc.Down(context.Background(), driving.DownRequest{
				BaseDir:        "/proj",
				RemoveVolumes:  true,
				AssumeYes:      true,
				NonInteractive: sc.nonInteractive,
			})
			if err != nil {
				t.Fatalf("Down: %v", err)
			}
			if !resp.RemovedVolumes {
				t.Errorf("RemovedVolumes = false, want true")
			}
			if len(f.confirmer.removeVolumesCalls) != 0 {
				t.Errorf("Confirmer was called despite --yes")
			}
			if !f.engine.downOptions.RemoveVolumes {
				t.Errorf("Engine called with RemoveVolumes=false, want true")
			}
		})
	}
}

func TestDownService_TruthTable_Row3_VolumesNonInteractive_FailFast(t *testing.T) {
	t.Parallel()
	// Row 3: RemoveVolumes=true, AssumeYes=false, NonInteractive=true
	// → ErrConfirmationRequired without confirmer or engine call.
	f := newDownFixture(t)
	_, err := f.svc.Down(context.Background(), driving.DownRequest{
		BaseDir:        "/proj",
		RemoveVolumes:  true,
		AssumeYes:      false,
		NonInteractive: true,
	})
	if !errors.Is(err, driving.ErrConfirmationRequired) {
		t.Errorf("expected ErrConfirmationRequired, got: %v", err)
	}
	if len(f.confirmer.removeVolumesCalls) != 0 {
		t.Errorf("Confirmer was called in --no-interactive path; should fail-fast")
	}
	if f.engine.downCallCount != 0 {
		t.Errorf("ComposeDown was called despite confirmation refusal")
	}
}

func TestDownService_TruthTable_Row4_InteractiveConfirmTrue_Proceeds(t *testing.T) {
	t.Parallel()
	// Row 4a: RemoveVolumes=true, AssumeYes=false, NonInteractive=false,
	// confirmer answers (true, nil) → proceed with RemoveVolumes=true.
	f := newDownFixture(t)
	f.confirmer.removeVolumesAnswer = true
	resp, err := f.svc.Down(context.Background(), driving.DownRequest{
		BaseDir:       "/proj",
		RemoveVolumes: true,
	})
	if err != nil {
		t.Fatalf("Down: %v", err)
	}
	if !resp.RemovedVolumes {
		t.Errorf("RemovedVolumes = false, want true")
	}
	if len(f.confirmer.removeVolumesCalls) != 1 {
		t.Errorf("Confirmer was called %d times, want 1", len(f.confirmer.removeVolumesCalls))
	}
	if f.confirmer.removeVolumesCalls[0].BaseDir != "/proj" {
		t.Errorf("Confirmer called with BaseDir=%q, want /proj", f.confirmer.removeVolumesCalls[0].BaseDir)
	}
	if !f.engine.downOptions.RemoveVolumes {
		t.Errorf("Engine called with RemoveVolumes=false, want true")
	}
}

func TestDownService_TruthTable_Row4_InteractiveConfirmFalse_ReturnsErrConfirmation(t *testing.T) {
	t.Parallel()
	// Row 4b: confirmer answers (false, nil) → ErrConfirmationRequired,
	// no engine call.
	f := newDownFixture(t)
	f.confirmer.removeVolumesAnswer = false
	_, err := f.svc.Down(context.Background(), driving.DownRequest{
		BaseDir:       "/proj",
		RemoveVolumes: true,
	})
	if !errors.Is(err, driving.ErrConfirmationRequired) {
		t.Errorf("expected ErrConfirmationRequired, got: %v", err)
	}
	if len(f.confirmer.removeVolumesCalls) != 1 {
		t.Errorf("Confirmer was called %d times, want exactly 1", len(f.confirmer.removeVolumesCalls))
	}
	if f.engine.downCallCount != 0 {
		t.Errorf("ComposeDown was called despite confirmation refusal")
	}
}

var errConfirmerSimulated = errors.New("confirmer io failure")

func TestDownService_TruthTable_Row4_InteractiveConfirmError_PropagatesNonSentinel(t *testing.T) {
	t.Parallel()
	// Row 4c: confirmer returns error → wrapped error (not a use-
	// case sentinel; this is a technical failure).
	f := newDownFixture(t)
	f.confirmer.removeVolumesErr = errConfirmerSimulated
	_, err := f.svc.Down(context.Background(), driving.DownRequest{
		BaseDir:       "/proj",
		RemoveVolumes: true,
	})
	if err == nil {
		t.Fatal("expected error from confirmer failure")
	}
	if errors.Is(err, driving.ErrConfirmationRequired) {
		t.Errorf("confirmer IO error leaked into ErrConfirmationRequired: %v", err)
	}
	if !errors.Is(err, errConfirmerSimulated) {
		t.Errorf("underlying confirmer error not preserved: %v", err)
	}
	if f.engine.downCallCount != 0 {
		t.Errorf("ComposeDown was called despite confirmer error")
	}
}

// --- Engine sentinel pass-through pins ---

func TestDownService_EngineReturnsErrDockerUnavailable_PassesThrough(t *testing.T) {
	t.Parallel()
	f := newDownFixture(t)
	f.engine.scriptDown(driven.ErrDockerUnavailable)
	_, err := f.svc.Down(context.Background(), driving.DownRequest{BaseDir: "/proj"})
	if !errors.Is(err, driven.ErrDockerUnavailable) {
		t.Errorf("expected errors.Is(err, ErrDockerUnavailable), got: %v", err)
	}
}

func TestDownService_EngineReturnsErrComposeRuntime_PassesThrough(t *testing.T) {
	t.Parallel()
	f := newDownFixture(t)
	f.engine.scriptDown(driven.ErrComposeRuntime)
	_, err := f.svc.Down(context.Background(), driving.DownRequest{BaseDir: "/proj"})
	if !errors.Is(err, driven.ErrComposeRuntime) {
		t.Errorf("expected errors.Is(err, ErrComposeRuntime), got: %v", err)
	}
}

func TestDownService_ProgressSinkWiredToEngine(t *testing.T) {
	t.Parallel()
	f := newDownFixture(t)
	sink := io.Discard
	_, err := f.svc.Down(context.Background(), driving.DownRequest{
		BaseDir:      "/proj",
		ProgressSink: sink,
	})
	if err != nil {
		t.Fatalf("Down: %v", err)
	}
	if f.engine.downOptions.ProgressSink != sink {
		t.Errorf("ComposeDown called with ProgressSink = %v, want the caller's sink %v", f.engine.downOptions.ProgressSink, sink)
	}
}

// TestDownService_SilenceConfirmer_True_SwapsToNoopConfirmer pins the
// Request-time Gate-Branch from slice-v1-cli-json-dry-run-up-down
// T0-(d) Option (b) + T3-Implementation: when `req.SilenceConfirmer ==
// true`, `runConfirmationGate` Row 4 MUST substitute a local
// noopConfirmer{} BEFORE calling ConfirmRemoveVolumes — the wired
// `s.confirmer` MUST NOT be called.
//
// T7-HIGH-1 closes the test-coverage gap that the CLI-Acceptance-
// Tests left open: CLI-Acceptance stubs the use case before
// runConfirmationGate executes, so the real branch was never
// exercised. This Application-Layer test pins the branch directly.
//
// Refuse-by-Default contract (R2-MED-2 festgezurrt): noopConfirmer.
// ConfirmRemoveVolumes returns `(false, nil)` → runConfirmationGate
// Z. 138 returns wrapped ErrConfirmationRequired. JSON-Mode-Konsument
// MUST opt in via `--yes` for destructive `--volumes`.
func TestDownService_SilenceConfirmer_True_SwapsToNoopConfirmer(t *testing.T) {
	t.Parallel()
	f := newDownFixture(t)
	// Setup: configure the wired confirmer to return (true, nil) —
	// IF the branch fails to swap, the use case would proceed to
	// ComposeDown. Our pin asserts it does NOT proceed AND the wired
	// confirmer is NOT called.
	f.confirmer.removeVolumesAnswer = true
	_, err := f.svc.Down(context.Background(), driving.DownRequest{
		BaseDir:          "/proj",
		RemoveVolumes:    true,
		AssumeYes:        false, // force Row 4 (not the AssumeYes-shortcut)
		NonInteractive:   false, // force Row 4 (not the NonInteractive-shortcut)
		SilenceConfirmer: true,  // T3 branch trigger
	})
	if !errors.Is(err, driving.ErrConfirmationRequired) {
		t.Fatalf("expected ErrConfirmationRequired (Refuse-by-Default), got: %v", err)
	}
	// Defense-Pin: the wired confirmer MUST NOT have been called.
	// If the branch fails to swap, fakeConfirmer.removeVolumesCalls
	// would be 1 (with removeVolumesAnswer=true), and the use case
	// would proceed to ComposeDown.
	if len(f.confirmer.removeVolumesCalls) != 0 {
		t.Errorf("wired confirmer was called %d times despite SilenceConfirmer=true; the branch did not swap to noopConfirmer", len(f.confirmer.removeVolumesCalls))
	}
	// Defense-Pin: ComposeDown MUST NOT have been called either
	// (the Gate fail-fasts before the engine call).
	if f.engine.downCallCount != 0 {
		t.Errorf("ComposeDown was called %d times despite refused confirmation", f.engine.downCallCount)
	}
}

// TestDownService_SilenceConfirmer_False_UsesWiredConfirmer is the
// contrast-pin to TestDownService_SilenceConfirmer_True: without the
// flag, the wired confirmer IS called and its answer drives the gate.
func TestDownService_SilenceConfirmer_False_UsesWiredConfirmer(t *testing.T) {
	t.Parallel()
	f := newDownFixture(t)
	f.confirmer.removeVolumesAnswer = true // proceed
	_, err := f.svc.Down(context.Background(), driving.DownRequest{
		BaseDir:          "/proj",
		RemoveVolumes:    true,
		AssumeYes:        false,
		NonInteractive:   false,
		SilenceConfirmer: false, // wired confirmer MUST be used
	})
	if err != nil {
		t.Fatalf("Down: %v", err)
	}
	if len(f.confirmer.removeVolumesCalls) != 1 {
		t.Errorf("wired confirmer was called %d times, want exactly 1", len(f.confirmer.removeVolumesCalls))
	}
	// T8-Bestätigungsrunde LOW-1 Symmetrie-Pin: bei wired-Confirmer
	// + answer=true MUSS ComposeDown durchgereicht werden — schließt
	// die Lücke dass ein Branch sowohl Confirmer-Call als auch
	// Engine-Call stilllegen könnte.
	if f.engine.downCallCount != 1 {
		t.Errorf("ComposeDown was called %d times, want exactly 1 (wired-Confirmer + answer=true MUST proceed)", f.engine.downCallCount)
	}
}
