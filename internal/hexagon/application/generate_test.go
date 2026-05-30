package application_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

const generateTestBaseDir = "/tmp/u-boot-generate-test/demo"

// newGenerateService constructs the service against the shared
// in-memory FS plus the yaml.v3-backed fake codec. The BaseDir is
// pre-registered as an existing directory so the Exists() check on
// u-boot.yaml fails on file-absence, not on directory-absence.
func newGenerateService(t *testing.T) (*application.GenerateService, *fakeFS) {
	t.Helper()
	fs := newFakeFS()
	fs.markDirExists(generateTestBaseDir)
	y := &fakeYAML{}
	svc := application.NewGenerateService(fs, y, nil)
	return svc, fs
}

// seedGenerateUbootYAML writes a minimal u-boot.yaml under the test
// BaseDir so the project-initialized gate passes.
func seedGenerateUbootYAML(t *testing.T, fs *fakeFS) {
	t.Helper()
	body := "schemaVersion: 1\nproject:\n  name: t-uboot-gen\n"
	if err := fs.WriteFile(filepath.Join(generateTestBaseDir, "u-boot.yaml"),
		[]byte(body), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
}

func TestGenerate_BaseDirEmpty_NonSentinelError(t *testing.T) {
	t.Parallel()
	svc, _ := newGenerateService(t)
	_, err := svc.Generate(context.Background(), driving.GenerateRequest{
		BaseDir:  "",
		Artifact: domain.ArtifactChangelog,
	})
	if err == nil {
		t.Fatal("expected non-nil error for empty BaseDir")
	}
	// BaseDir-empty is a programmer-/CLI-bug, not a fachliche
	// condition — must NOT map to ErrProjectNotInitialized.
	if errors.Is(err, driving.ErrProjectNotInitialized) {
		t.Errorf("empty BaseDir leaked ErrProjectNotInitialized: %v", err)
	}
}

func TestGenerate_NoUbootYAML_ReturnsErrProjectNotInitialized(t *testing.T) {
	t.Parallel()
	svc, _ := newGenerateService(t)
	for _, art := range []domain.Artifact{
		domain.ArtifactChangelog,
		domain.ArtifactReadme,
		domain.ArtifactEnvExample,
		domain.ArtifactDevcontainer,
	} {
		_, err := svc.Generate(context.Background(), driving.GenerateRequest{
			BaseDir:  generateTestBaseDir,
			Artifact: art,
		})
		if !errors.Is(err, driving.ErrProjectNotInitialized) {
			t.Errorf("artifact=%s: err = %v, want wrap of ErrProjectNotInitialized", art, err)
		}
	}
}

// TestGenerate_StubHandlers pins that every M7-T1 handler is reachable
// at runtime but unimplemented. The stub-pin reduces by one with each
// of T2..T5 as those tranches replace handlers; T5 removes the test
// when the last stub is gone (slice-m7-generate.md T5 DoD).
//
// Pin tracks remaining stubs: 4 in T1 → 3 in T2 → 2 in T3 → 1 in T4 →
// 0 in T5 (test deleted).
func TestGenerate_StubHandlers_AllFourReturnErrStubHandler(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)

	cases := []domain.Artifact{
		domain.ArtifactChangelog,
		domain.ArtifactReadme,
		domain.ArtifactEnvExample,
		domain.ArtifactDevcontainer,
	}
	for _, art := range cases {
		_, err := svc.Generate(context.Background(), driving.GenerateRequest{
			BaseDir:  generateTestBaseDir,
			Artifact: art,
		})
		if err == nil {
			t.Errorf("artifact=%s: expected stub-handler error, got nil", art)
			continue
		}
		if !errors.Is(err, application.ErrStubHandlerForTest) {
			t.Errorf("artifact=%s: err = %v, want wrap of errStubHandler", art, err)
		}
	}
}

// TestGenerate_UnknownArtifactValue_DefensiveBranch pins the
// out-of-enum-range guard. The CLI validates via [domain.NewArtifact]
// before constructing the request, so this branch is unreachable in
// production; the test protects against a future enum addition that
// forgets to update the dispatch switch.
func TestGenerate_UnknownArtifactValue_DefensiveBranch(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)

	_, err := svc.Generate(context.Background(), driving.GenerateRequest{
		BaseDir:  generateTestBaseDir,
		Artifact: domain.Artifact(99),
	})
	if !errors.Is(err, domain.ErrInvalidArtifact) {
		t.Errorf("err = %v, want wrap of domain.ErrInvalidArtifact", err)
	}
}
