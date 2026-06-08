package application_test

import (
	"context"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driven/recordingfs"
	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// configFactoryForProd mirrors the cmd/uboot/main.go newPreviewFSFactory
// closure: PreviewNone → production FS (no recorder), PreviewDryRun →
// recording FS with passthrough off (capture, no disk write),
// PreviewAndApply → recording FS with passthrough on.
func configFactoryForProd(prod driven.FileSystem) func(driving.PreviewMode) (driven.FileSystem, driven.RecorderPort) {
	return func(mode driving.PreviewMode) (driven.FileSystem, driven.RecorderPort) {
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
}

// TestConfigService_WithFactory_DryRunDoesNotTouchProduction pins the
// T4 contract: a PreviewDryRun `config set` routes the WriteFile
// through the recording FS (passthrough off), so the production
// u-boot.yaml is NOT modified — a follow-up Get still returns the old
// value — while the recorder DID capture the intended mutation.
func TestConfigService_WithFactory_DryRunDoesNotTouchProduction(t *testing.T) {
	t.Parallel()
	prod := newFakeFS()
	prod.markDirExists(configTestBaseDir)
	seedConfigUbootYAML(t, prod) // project.name: t-uboot-config

	factory := configFactoryForProd(prod)
	svc := application.NewConfigServiceWithFactory(factory, &fakeYAML{}, nil)

	resp, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir:     configTestBaseDir,
		Path:        mustConfigPath(t, "project.name"),
		Value:       "renamed-in-dryrun",
		PreviewMode: driving.PreviewDryRun,
	})
	if err != nil {
		t.Fatalf("Set: %v", err)
	}
	if resp.OldValue != "t-uboot-config" || resp.NewValue != "renamed-in-dryrun" {
		t.Errorf("response = {Old:%q New:%q}, want {t-uboot-config renamed-in-dryrun}",
			resp.OldValue, resp.NewValue)
	}

	// R-T4-1: the captured mutation must surface on the response so
	// the CLI --diff renderer has its byte source. Exactly one entry
	// (u-boot.yaml), with the patched NewContent populated.
	if len(resp.PlannedFiles) != 1 {
		t.Fatalf("PreviewDryRun must surface 1 PlannedFile (u-boot.yaml); got %d: %+v",
			len(resp.PlannedFiles), resp.PlannedFiles)
	}
	pf := resp.PlannedFiles[0]
	if pf.Path != "u-boot.yaml" {
		t.Errorf("PlannedFile.Path = %q, want project-relative u-boot.yaml", pf.Path)
	}
	if len(pf.NewContent) == 0 {
		t.Error("PlannedFile.NewContent empty — --diff renderer would have no patched bytes")
	}

	// Production must be untouched: a fresh Get reads via s.fs (prod)
	// and still sees the original value.
	got, err := svc.Get(context.Background(), driving.ConfigGetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "project.name"),
	})
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Value != "t-uboot-config" {
		t.Errorf("dry-run leaked to production: Get = %q, want unchanged t-uboot-config", got.Value)
	}
}

// TestConfigService_WithFactory_PreviewNoneWritesProduction pins the
// default branch: PreviewNone routes to the production FS, so the
// write persists and a follow-up Get returns the new value.
func TestConfigService_WithFactory_PreviewNoneWritesProduction(t *testing.T) {
	t.Parallel()
	prod := newFakeFS()
	prod.markDirExists(configTestBaseDir)
	seedConfigUbootYAML(t, prod)

	factory := configFactoryForProd(prod)
	svc := application.NewConfigServiceWithFactory(factory, &fakeYAML{}, nil)

	resp, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir:     configTestBaseDir,
		Path:        mustConfigPath(t, "project.name"),
		Value:       "renamed-for-real",
		PreviewMode: driving.PreviewNone,
	})
	if err != nil {
		t.Fatalf("Set: %v", err)
	}
	// PreviewNone uses production FS (no recorder) → no PlannedFiles.
	if resp.PlannedFiles != nil {
		t.Errorf("PreviewNone must not surface PlannedFiles, got %+v", resp.PlannedFiles)
	}

	got, err := svc.Get(context.Background(), driving.ConfigGetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "project.name"),
	})
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Value != "renamed-for-real" {
		t.Errorf("PreviewNone must persist: Get = %q, want renamed-for-real", got.Value)
	}
}

// TestConfigService_LegacyConstructorIgnoresPreviewMode pins the
// backward-compat guarantee: NewConfigService (nil factory) ignores
// PreviewMode entirely — even PreviewDryRun writes to the passed FS.
func TestConfigService_LegacyConstructorIgnoresPreviewMode(t *testing.T) {
	t.Parallel()
	prod := newFakeFS()
	prod.markDirExists(configTestBaseDir)
	seedConfigUbootYAML(t, prod)

	svc := application.NewConfigService(prod, &fakeYAML{}, nil)

	resp, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir:     configTestBaseDir,
		Path:        mustConfigPath(t, "project.name"),
		Value:       "legacy-write",
		PreviewMode: driving.PreviewDryRun, // ignored by the nil-factory path
	})
	if err != nil {
		t.Fatalf("Set: %v", err)
	}
	// nil factory → nil recorder → no PlannedFiles (backward-compat).
	if resp.PlannedFiles != nil {
		t.Errorf("legacy path must not surface PlannedFiles, got %+v", resp.PlannedFiles)
	}

	got, err := svc.Get(context.Background(), driving.ConfigGetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "project.name"),
	})
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Value != "legacy-write" {
		t.Errorf("legacy path must ignore PreviewMode and write: Get = %q, want legacy-write", got.Value)
	}
}
