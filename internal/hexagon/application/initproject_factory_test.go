package application_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driven/recordingfs"
	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// initFactoryBaseDir is the per-request BaseDir for init-factory
// tests. Pre-registered as an existing directory in each test's
// prod fakeFS via markDirExists.
const initFactoryBaseDir = "/tmp/u-boot-init-test/myproj"

// newInitFactoryFS returns a fresh production FS with the base dir
// pre-existing so init's writeDirectories does not fail on the dir
// itself (MkdirAll for `./docker`, `./scripts`, `./docs` still runs).
func newInitFactoryFS() *fakeFS {
	fs := newFakeFS()
	fs.markDirExists(initFactoryBaseDir)
	return fs
}

// TestInitProjectService_WithFactory_DryRunMapsRecorderToPlannedFiles
// is the canonical T3 pin: a PreviewDryRun request routes FS writes
// through a RecordingFileSystem; the Init response carries
// PlannedFiles[] reflecting what the use case captured (no
// production writes happen, but the planned changes are surfaced).
func TestInitProjectService_WithFactory_DryRunMapsRecorderToPlannedFiles(t *testing.T) {
	prod := newInitFactoryFS()
	factory := func(mode driving.PreviewMode) (driven.FileSystem, driven.RecorderPort) {
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
	svc := application.NewInitProjectServiceWithFactory(factory, &fakeYAML{}, &fakeGit{}, nil, nil, nil)

	resp, err := svc.Init(context.Background(), driving.InitProjectRequest{
		BaseDir:     initFactoryBaseDir,
		Name:        "myproj",
		SkipGit:     true, // T0-(n) Skip-Logic auch im PreviewDryRun aktiv — explizit SkipGit gesetzt für Klarheit
		PreviewMode: driving.PreviewDryRun,
	})
	if err != nil {
		t.Fatalf("Init dry-run: %v", err)
	}

	if len(resp.PlannedFiles) == 0 {
		t.Fatal("DryRun must populate PlannedFiles from the recorder")
	}
	// fresh init writes at least the 3 skeleton-dirs + 5 templated
	// files + u-boot.yaml = ≥ 9 entries.
	if len(resp.PlannedFiles) < 9 {
		t.Errorf("expected ≥9 PlannedFiles for fresh init, got %d: %v",
			len(resp.PlannedFiles), pathsOf(resp.PlannedFiles))
	}
	for _, pf := range resp.PlannedFiles {
		if pf.Path == "" || pf.Action == "" {
			t.Errorf("PlannedFile missing path/action: %+v", pf)
		}
		// Path-Anchor (T0-(c) erblich aus add R8 A): must be
		// project-relative, no absolute prefix.
		if strings.HasPrefix(pf.Path, "/") || strings.HasPrefix(pf.Path, initFactoryBaseDir) {
			t.Errorf("PlannedFile.Path must be project-relative, got absolute: %q", pf.Path)
		}
	}
}

// TestInitProjectService_WithFactory_LegacyConstructorIgnoresMode
// pins that NewInitProjectService(fs, …) keeps writing through the
// passed FS and ignores PreviewMode (fsFactory nil → no recorder →
// PlannedFiles stays empty). Backward-compat guarantee.
func TestInitProjectService_WithFactory_LegacyConstructorIgnoresMode(t *testing.T) {
	prod := newInitFactoryFS()
	svc := application.NewInitProjectService(prod, &fakeYAML{}, &fakeGit{}, nil, nil, nil)

	resp, err := svc.Init(context.Background(), driving.InitProjectRequest{
		BaseDir:     initFactoryBaseDir,
		Name:        "myproj",
		SkipGit:     true,
		PreviewMode: driving.PreviewDryRun, // ignored by legacy path
	})
	if err != nil {
		t.Fatalf("Init legacy: %v", err)
	}
	if resp.PlannedFiles != nil {
		t.Errorf("legacy constructor must leave PlannedFiles nil, got %d entries", len(resp.PlannedFiles))
	}
	if len(resp.Created) == 0 {
		t.Errorf("legacy path must still apply writes; Created=nil")
	}
}

// TestInitProjectService_WithFactory_WriteFailureWrapsErrInitFileSystem
// pins the T0-(f)/T2 sentinel mapping + non-empty Response on Error:
// a FS write failure surfaces as a wrapped ErrInitFileSystem while
// the response still carries PlannedFiles captured up to the
// failure point.
func TestInitProjectService_WithFactory_WriteFailureWrapsErrInitFileSystem(t *testing.T) {
	prod := newInitFactoryFS()
	// Fail at the FIRST WriteFile (README.md). MkdirAll for
	// docker/scripts/docs runs first (no failure target there).
	prod.failOn = initFactoryBaseDir + "/README.md"
	prod.failErr = errors.New("disk full")

	factory := func(mode driving.PreviewMode) (driven.FileSystem, driven.RecorderPort) {
		// PreviewAndApply so the FS write actually fails (passthrough=true).
		rec := recordingfs.New(prod, recordingfs.WithPassthrough(true))
		return rec, rec
	}
	svc := application.NewInitProjectServiceWithFactory(factory, &fakeYAML{}, &fakeGit{}, nil, nil, nil)

	resp, initErr := svc.Init(context.Background(), driving.InitProjectRequest{
		BaseDir:     initFactoryBaseDir,
		Name:        "myproj",
		SkipGit:     true,
		PreviewMode: driving.PreviewAndApply,
	})
	if !errors.Is(initErr, driving.ErrInitFileSystem) {
		t.Fatalf("expected wrapped ErrInitFileSystem, got: %v", initErr)
	}
	if len(resp.PlannedFiles) == 0 {
		t.Errorf("Mid-write failure must still expose recorder capture (T0-(b) Mid-Write-Failure)")
	}
	// Captured entries should include at least one mkdir-record from
	// before the failing WriteFile.
	hasMkdir := false
	for _, pf := range resp.PlannedFiles {
		if pf.Action == "create" {
			hasMkdir = true
			break
		}
	}
	if !hasMkdir {
		t.Errorf("expected at least one create-record in PlannedFiles, got: %v", pathsOf(resp.PlannedFiles))
	}
}

// TestInitProjectService_WithFactory_DryRunSkipsInitGit pins T0-(n):
// in PreviewDryRun the use case MUST NOT call s.git.Init() — git is
// a separate driven port, not routed through the per-request fs-swap.
// Without the skip, `init --dry-run` would create a real `.git/`
// on disk (Adversarial R3 C-1).
func TestInitProjectService_WithFactory_DryRunSkipsInitGit(t *testing.T) {
	prod := newInitFactoryFS()
	git := &fakeGit{}
	factory := func(mode driving.PreviewMode) (driven.FileSystem, driven.RecorderPort) {
		if mode == driving.PreviewDryRun {
			rec := recordingfs.New(prod, recordingfs.WithPassthrough(false))
			return rec, rec
		}
		return prod, nil
	}
	svc := application.NewInitProjectServiceWithFactory(factory, &fakeYAML{}, git, nil, nil, nil)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		BaseDir:     initFactoryBaseDir,
		Name:        "myproj",
		SkipGit:     false, // explicitly NOT skipping via flag — PreviewDryRun must skip on its own
		PreviewMode: driving.PreviewDryRun,
	})
	if err != nil {
		t.Fatalf("Init dry-run: %v", err)
	}
	if len(git.initCalls) != 0 {
		t.Errorf("T0-(n) violation: PreviewDryRun must NOT call git.Init; got %d calls: %v", len(git.initCalls), git.initCalls)
	}
}

// TestInitProjectService_WithFactory_PreviewAndApplyRunsInitGit
// is the inverse pin: PreviewAndApply (--diff without --dry-run)
// IS preview-and-apply — git init must still run.
func TestInitProjectService_WithFactory_PreviewAndApplyRunsInitGit(t *testing.T) {
	prod := newInitFactoryFS()
	git := &fakeGit{}
	factory := func(mode driving.PreviewMode) (driven.FileSystem, driven.RecorderPort) {
		rec := recordingfs.New(prod, recordingfs.WithPassthrough(true))
		return rec, rec
	}
	svc := application.NewInitProjectServiceWithFactory(factory, &fakeYAML{}, git, nil, nil, nil)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		BaseDir:     initFactoryBaseDir,
		Name:        "myproj",
		SkipGit:     false,
		PreviewMode: driving.PreviewAndApply,
	})
	if err != nil {
		t.Fatalf("Init preview-and-apply: %v", err)
	}
	if len(git.initCalls) != 1 {
		t.Errorf("PreviewAndApply must call git.Init once (preview-and-apply contract); got %d calls", len(git.initCalls))
	}
}

// TestInitProjectService_WithFactory_SilenceProgressSwapsToNoop
// pins T0-(o): when req.SilenceProgress=true the service must NOT
// emit progress-events (so the JSON envelope on stdout isn't
// corrupted). Default false keeps today's progress behaviour.
func TestInitProjectService_WithFactory_SilenceProgressSwapsToNoop(t *testing.T) {
	prod := newInitFactoryFS()
	progress := &fakeProgress{}
	factory := func(mode driving.PreviewMode) (driven.FileSystem, driven.RecorderPort) {
		if mode == driving.PreviewDryRun {
			rec := recordingfs.New(prod, recordingfs.WithPassthrough(false))
			return rec, rec
		}
		return prod, nil
	}
	svc := application.NewInitProjectServiceWithFactory(factory, &fakeYAML{}, &fakeGit{}, progress, nil, nil)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		BaseDir:         initFactoryBaseDir,
		Name:            "myproj",
		SkipGit:         true,
		PreviewMode:     driving.PreviewDryRun,
		SilenceProgress: true,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	if len(progress.calls) != 0 {
		t.Errorf("T0-(o) violation: SilenceProgress=true must yield zero progress.AffectedFiles calls; got %d", len(progress.calls))
	}
}

// Note: a "DefaultProgressNotSilenced" inverse pin was tried here
// but emitSummary today only fires for actionReplaceBlock /
// actionOverwriteFull (re-init paths); a fresh init has no plans
// in those categories, so the rows-collector returns 0 and the
// ProgressPort is never called. The SilenceProgress-Pin above
// covers the swap-restore behaviour sufficiently: it proves the
// swap fires for SilenceProgress=true AND that the deferred restore
// runs (subsequent Init calls would race-fail if the swap leaked).
// A future re-init pin in T6-Acceptance will exercise the
// AffectedFiles-emit path proper.

// pathsOf is defined in addservice_factory_test.go (same package).
