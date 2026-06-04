package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driven/recordingfs"
	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// TestAddService_WithFactory_DryRunMapsRecorderToPlannedFiles is the
// canonical T1-C pin: a PreviewDryRun request routes FS writes through
// a RecordingFileSystem; the Add response carries PlannedFiles[]
// reflecting what the use case captured (no production writes
// happen, but the planned changes are surfaced).
func TestAddService_WithFactory_DryRunMapsRecorderToPlannedFiles(t *testing.T) {
	prod := newFakeFS()
	prod.markDirExists(addTestBaseDir)
	seedUBootYAML(t, prod, "schemaVersion: 1\nproject:\n  name: test\n")

	factory := func(mode driving.AddPreviewMode) (driven.FileSystem, driven.RecorderPort) {
		switch mode {
		case driving.PreviewDryRun:
			rec := recordingfs.New(prod, recordingfs.WithPassthrough(false))
			return rec, rec
		case driving.PreviewAndApply:
			rec := recordingfs.New(prod, recordingfs.WithPassthrough(true))
			return rec, rec
		default:
			return prod, nil
		}
	}
	svc := application.NewAddServiceServiceWithFactory(factory, &fakeYAML{}, nil, nil)

	resp, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: mustServiceName(t, "postgres"),
		PreviewMode: driving.PreviewDryRun,
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	if len(resp.PlannedFiles) == 0 {
		t.Fatal("DryRun must populate PlannedFiles from the recorder")
	}
	// add postgres on a fresh setup writes at least u-boot.yaml +
	// compose.yaml + .env.example (slice T0-(b) Mutations-Matrix).
	if len(resp.PlannedFiles) < 3 {
		t.Errorf("expected ≥3 PlannedFiles for fresh add postgres, got %d: %v",
			len(resp.PlannedFiles), pathsOf(resp.PlannedFiles))
	}
	// Plan files must carry their snapshots so the CLI-adapter diff
	// renderer (T2) has something to work with.
	for _, pf := range resp.PlannedFiles {
		if pf.Path == "" || pf.Action == "" {
			t.Errorf("PlannedFile missing path/action: %+v", pf)
		}
	}
}

// TestAddService_WithFactory_LegacyConstructorIgnoresMode pins that the
// today's NewAddServiceService(fs, ...) path keeps writing through the
// passed FS and ignores PreviewMode (fsFactory is nil → no recorder →
// PlannedFiles stays empty). Backward-compat guarantee.
func TestAddService_WithFactory_LegacyConstructorIgnoresMode(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs, "schemaVersion: 1\nproject:\n  name: test\n")

	resp, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: mustServiceName(t, "postgres"),
		PreviewMode: driving.PreviewDryRun, // ignored by legacy path
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if resp.PlannedFiles != nil {
		t.Errorf("legacy constructor must leave PlannedFiles nil, got %d entries",
			len(resp.PlannedFiles))
	}
	if len(resp.Changed) == 0 {
		t.Errorf("legacy path must still apply writes; Changed=nil")
	}
}

// TestAddService_WithFactory_WriteFailureWrapsErrAddFileSystem pins the
// T0-(j) sentinel-mapping plus non-empty Response on Error: a FS write
// failure surfaces as a wrapped ErrAddFileSystem while the response
// still carries PlannedFiles captured up to the failure point.
func TestAddService_WithFactory_WriteFailureWrapsErrAddFileSystem(t *testing.T) {
	prod := newFakeFS()
	prod.markDirExists(addTestBaseDir)
	seedUBootYAML(t, prod, "schemaVersion: 1\nproject:\n  name: test\n")
	// T0-(b) Mid-Write-Failure: the 2nd WriteFile fails. The first
	// production write (u-boot.yaml) succeeds; the second
	// (compose.yaml) returns failErr.
	prod.failOn = addTestBaseDir + "/compose.yaml"
	prod.failErr = errors.New("disk full")

	factory := func(mode driving.AddPreviewMode) (driven.FileSystem, driven.RecorderPort) {
		// PreviewAndApply: capture AND delegate, so the production
		// write actually fails and surfaces the error.
		if mode == driving.PreviewAndApply {
			rec := recordingfs.New(prod, recordingfs.WithPassthrough(true))
			return rec, rec
		}
		return prod, nil
	}
	svc := application.NewAddServiceServiceWithFactory(factory, &fakeYAML{}, nil, nil)

	resp, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: mustServiceName(t, "postgres"),
		PreviewMode: driving.PreviewAndApply,
	})
	if !errors.Is(err, driving.ErrAddFileSystem) {
		t.Fatalf("want wrapped ErrAddFileSystem, got %v", err)
	}
	if len(resp.PlannedFiles) == 0 {
		t.Errorf("non-empty Response on Error: PlannedFiles must show calls up to failure, got none")
	}
}

func pathsOf(pfs []driving.PlannedFile) []string {
	out := make([]string, len(pfs))
	for i, pf := range pfs {
		out[i] = pf.Path + "(" + pf.Action + ")"
	}
	return out
}

