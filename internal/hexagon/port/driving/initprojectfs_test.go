package driving_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// TestErrInitFileSystem_Identity pins the slice-v1-cli-json-dry-
// run-init T2 sentinel identity: the new ErrInitFileSystem must
// exist, carry the expected message, and stay distinct from other
// init sentinels (so multi-`%w` chains route correctly through
// the CLI's mapInitErrorToDiagnostic switch).
func TestErrInitFileSystem_Identity(t *testing.T) {
	t.Parallel()
	if driving.ErrInitFileSystem == nil {
		t.Fatal("ErrInitFileSystem is nil")
	}
	const want = "init: filesystem mutation failed"
	if got := driving.ErrInitFileSystem.Error(); got != want {
		t.Errorf("ErrInitFileSystem.Error() = %q, want %q", got, want)
	}
}

// TestErrInitFileSystem_DistinctFromOtherInitSentinels pins that
// the new FS-class sentinel is not accidentally aliased to one of
// the fachlich init sentinels — without distinctness, the T0-(f)
// switch-order pflicht (ErrInitFileSystem FIRST) could not
// reliably distinguish the FS-class from the user-action class.
func TestErrInitFileSystem_DistinctFromOtherInitSentinels(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		other error
	}{
		{"ErrProjectExists", driving.ErrProjectExists},
		{"ErrFileExists", driving.ErrFileExists},
		{"ErrBackupSourceMissing", driving.ErrBackupSourceMissing},
		{"ErrBackupSuffixExhausted", driving.ErrBackupSuffixExhausted},
		{"ErrBackupUnsupportedKind", driving.ErrBackupUnsupportedKind},
		{"ErrForceRequiresBackup", driving.ErrForceRequiresBackup},
		{"ErrTemplateConflictsWithFlag", driving.ErrTemplateConflictsWithFlag},
		{"ErrBaseDirMissing", driving.ErrBaseDirMissing},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Bare-sentinel identity check via `==` (errors.New
			// creates pointer-distinct values). If a future refactor
			// aliased two sentinels to the same `errors.New(...)`
			// call, the FS-class switch-order would collapse because
			// errors.Is would match both via the same target.
			if driving.ErrInitFileSystem == tc.other {
				t.Errorf("ErrInitFileSystem must NOT be identical to %s — both would be reachable via the same errors.Is() check and the FS-class switch-order would collapse", tc.name)
			}
		})
	}
}

// TestErrInitFileSystem_MultiWrapMatchesBoth pins the Go 1.20+
// multi-`%w` behaviour the T0-(f) Switch-Order-Pflicht relies on:
// `fmt.Errorf("%w: %w", ErrInitFileSystem, sentinel)` must
// errors.Is-match BOTH sentinels (so the switch-order matters).
// This is the inverse-direction proof of the
// distinctness test above.
func TestErrInitFileSystem_MultiWrapMatchesBoth(t *testing.T) {
	t.Parallel()
	wrapped := fmt.Errorf("init: write u-boot.yaml: %w: %w",
		driving.ErrInitFileSystem, driving.ErrProjectExists)

	if !errors.Is(wrapped, driving.ErrInitFileSystem) {
		t.Error("multi-%w wrap should match ErrInitFileSystem")
	}
	if !errors.Is(wrapped, driving.ErrProjectExists) {
		t.Error("multi-%w wrap should also match ErrProjectExists — this is the precise reason switch-order matters in mapInitErrorToDiagnostic (T0-(f))")
	}
}

// TestInitProjectRequest_PreviewModeField pins the T2-added field's
// type and zero-value. Compile-time check via var-decl: if the
// field is renamed or removed, the test file fails to build.
func TestInitProjectRequest_PreviewModeField(t *testing.T) {
	t.Parallel()
	var req driving.InitProjectRequest
	if req.PreviewMode != driving.PreviewNone {
		t.Errorf("PreviewMode zero-value: want PreviewNone, got %v", req.PreviewMode)
	}
	req.PreviewMode = driving.PreviewDryRun
	if req.PreviewMode != driving.PreviewDryRun {
		t.Errorf("PreviewMode assignment: want PreviewDryRun, got %v", req.PreviewMode)
	}
}

// TestInitProjectRequest_SilenceProgressField pins the T0-(o) field
// that the init-RunE sets to true in JSON-mode to silence the
// ProgressPort during emitSummary.
func TestInitProjectRequest_SilenceProgressField(t *testing.T) {
	t.Parallel()
	var req driving.InitProjectRequest
	if req.SilenceProgress {
		t.Error("SilenceProgress zero-value: want false")
	}
	req.SilenceProgress = true
	if !req.SilenceProgress {
		t.Error("SilenceProgress assignment failed")
	}
}

// TestInitProjectResponse_PlannedFilesAndChangesFields pins the T2-
// added carrier fields (mirrors AddServiceResponse).
func TestInitProjectResponse_PlannedFilesAndChangesFields(t *testing.T) {
	t.Parallel()
	resp := driving.InitProjectResponse{
		PlannedFiles: []driving.PlannedFile{{Path: "u-boot.yaml", Action: "create"}},
		Changes:      []driving.ChangeEntry{{Path: "u-boot.yaml", Count: 5}},
	}
	if len(resp.PlannedFiles) != 1 || resp.PlannedFiles[0].Path != "u-boot.yaml" {
		t.Errorf("PlannedFiles field unusable: %+v", resp.PlannedFiles)
	}
	if len(resp.Changes) != 1 || resp.Changes[0].Count != 5 {
		t.Errorf("Changes field unusable: %+v", resp.Changes)
	}
}
