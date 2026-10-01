package application_test

import (
	"bytes"
	"context"
	"errors"
	iofs "io/fs"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// truncatingFailFS reproduces the real failure topology of os.WriteFile:
// the target is truncated/partially written, THEN the call fails
// (disk full, signal, …). Pin B needs it: a no-op-failing fake would let
// an implementation without a real restore path pass.
type truncatingFailFS struct {
	*fakeFS
	failPath string
	failed   bool // fail only the first write (the restore must succeed)
}

func (f *truncatingFailFS) WriteFile(path string, data []byte, mode iofs.FileMode) error {
	if path == f.failPath && !f.failed {
		f.failed = true
		_ = f.fakeFS.WriteFile(path, data[:len(data)/2], mode) // partial mutation
		return errors.New("simulated disk full after truncate")
	}
	return f.fakeFS.WriteFile(path, data, mode)
}

func dirExists(t *testing.T, fs *fakeFS, path string) bool {
	t.Helper()
	ok, err := fs.Exists(path)
	if err != nil {
		t.Fatal(err)
	}
	return ok
}

const rollbackSeedYAML = "schemaVersion: 1\nproject:\n  name: t-uboot-gen\nservices:\n  postgres:\n    enabled: true\n"

func rollbackGen(svc *application.GenerateService) (driving.GenerateResponse, error) {
	return svc.Generate(context.Background(), driving.GenerateRequest{
		BaseDir: generateTestBaseDir, Artifact: domain.ArtifactDevcontainer,
		AllowExternalFeatureSources: []string{"https://example.test/features/x"}})
}

// Pin A: second file fails → first file removed, the freshly created
// .devcontainer/ directory removed, u-boot.yaml byte-identical.
func TestGenerateDevcontainer_Rollback_SecondFileFails(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLPostgres(t, fs)
	yamlBefore, _ := fs.ReadFile(filepath.Join(generateTestBaseDir, "u-boot.yaml"))
	fs.failOn, fs.failErr = dockerfilePath(), errors.New("simulated write failure")

	_, err := rollbackGen(svc)
	if !errors.Is(err, driving.ErrGenerateFileSystem) {
		t.Fatalf("err = %v, want ErrGenerateFileSystem", err)
	}
	if strings.Contains(err.Error(), "rollback incomplete") {
		t.Errorf("clean rollback must not report an incomplete state: %v", err)
	}
	if dirExists(t, fs, devcontainerJSONPath()) {
		t.Errorf("devcontainer.json must be rolled back")
	}
	if dirExists(t, fs, filepath.Join(generateTestBaseDir, ".devcontainer")) {
		t.Errorf("freshly created .devcontainer/ must be removed")
	}
	if after, _ := fs.ReadFile(filepath.Join(generateTestBaseDir, "u-boot.yaml")); !bytes.Equal(after, yamlBefore) {
		t.Errorf("u-boot.yaml changed:\n%s", after)
	}
}

// Pin B: both devcontainer files succeed, the u-boot.yaml write fails
// after truncating it → everything is restored, the yaml byte-for-byte.
func TestGenerateDevcontainer_Rollback_YAMLWriteFailsAfterTruncate(t *testing.T) {
	t.Parallel()
	fs := newFakeFS()
	fs.markDirExists(generateTestBaseDir)
	yamlPath := filepath.Join(generateTestBaseDir, "u-boot.yaml")
	if err := fs.WriteFile(yamlPath, []byte(rollbackSeedYAML), 0o644); err != nil {
		t.Fatal(err)
	}
	svc := application.NewGenerateService(&truncatingFailFS{fakeFS: fs, failPath: yamlPath}, &fakeYAML{}, nil)

	_, err := rollbackGen(svc)
	if !errors.Is(err, driving.ErrGenerateFileSystem) {
		t.Fatalf("err = %v, want ErrGenerateFileSystem", err)
	}
	if dirExists(t, fs, devcontainerJSONPath()) || dirExists(t, fs, dockerfilePath()) ||
		dirExists(t, fs, filepath.Join(generateTestBaseDir, ".devcontainer")) {
		t.Errorf("devcontainer files / directory must be rolled back")
	}
	after, _ := fs.ReadFile(yamlPath)
	if string(after) != rollbackSeedYAML {
		t.Errorf("u-boot.yaml not restored after the truncating failure:\n%q", after)
	}
}

// Existing files are restored to their previous content and mode, not
// deleted.
func TestGenerateDevcontainer_Rollback_RestoresExistingFiles(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLPostgres(t, fs)
	if _, err := generateDevcontainer(t, svc); err != nil {
		t.Fatalf("first generate: %v", err)
	}
	jsonBefore, _ := fs.ReadFile(devcontainerJSONPath())
	dfBefore, _ := fs.ReadFile(dockerfilePath())
	// Make both files stale so the next generate rewrites them, then fail on the Dockerfile.
	stale := bytes.ReplaceAll(jsonBefore, []byte(`"name": "t-uboot-gen"`), []byte(`"name": "old"`))
	if err := fs.WriteFile(devcontainerJSONPath(), stale, 0o640); err != nil {
		t.Fatal(err)
	}
	staleDF := bytes.ReplaceAll(dfBefore, []byte("t-uboot-gen"), []byte("old"))
	if err := fs.WriteFile(dockerfilePath(), staleDF, 0o644); err != nil {
		t.Fatal(err)
	}
	fs.failOn, fs.failErr = dockerfilePath(), errors.New("simulated write failure")

	if _, err := generateDevcontainer(t, svc); err == nil {
		t.Fatalf("expected the write failure")
	}
	if got, _ := fs.ReadFile(devcontainerJSONPath()); !bytes.Equal(got, stale) {
		t.Errorf("devcontainer.json not restored to its pre-call content")
	}
	info, _ := fs.Lstat(devcontainerJSONPath())
	if info.Mode().Perm() != 0o640 {
		t.Errorf("mode = %v, want 0640 restored", info.Mode().Perm())
	}
	if got, _ := fs.ReadFile(dockerfilePath()); !bytes.Equal(got, staleDF) {
		t.Errorf("Dockerfile changed")
	}
}

// A failing rollback is reported, the original error chain stays.
func TestGenerateDevcontainer_Rollback_FailureIsReported(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLPostgres(t, fs)
	fs.failOn, fs.failErr = dockerfilePath(), errors.New("simulated write failure")
	fs.failRemoveAll = errors.New("cannot remove")

	_, err := rollbackGen(svc)
	if !errors.Is(err, driving.ErrGenerateFileSystem) || !strings.Contains(err.Error(), "rollback incomplete") {
		t.Errorf("err = %v, want ErrGenerateFileSystem + rollback incomplete", err)
	}
}

// Dry-run never writes, so there is nothing to roll back and the
// recorder capture stays untouched.
func TestGenerateDevcontainer_Rollback_NotInDryRun(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLPostgres(t, fs)
	resp, err := svc.Generate(context.Background(), driving.GenerateRequest{
		BaseDir: generateTestBaseDir, Artifact: domain.ArtifactDevcontainer, PreviewMode: driving.PreviewDryRun})
	if err != nil {
		t.Fatalf("dry-run: %v", err)
	}
	_ = resp
}
