package application_test

import (
	"errors"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
)

func TestBackupPath_File_FirstSuffixWhenFree(t *testing.T) {
	fs := newFakeFS()
	src := "/proj/u-boot.yaml"
	if err := fs.WriteFile(src, []byte("payload"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	dst, err := application.BackupPath(fs, src)
	if err != nil {
		t.Fatalf("BackupPath: %v", err)
	}
	if want := src + ".bak"; dst != want {
		t.Errorf("dst = %q, want %q", dst, want)
	}
}

func TestBackupPath_File_NumericSuffixOnCollision(t *testing.T) {
	fs := newFakeFS()
	src := "/proj/u-boot.yaml"
	for _, p := range []string{src, src + ".bak"} {
		if err := fs.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatalf("setup %s: %v", p, err)
		}
	}

	dst, err := application.BackupPath(fs, src)
	if err != nil {
		t.Fatalf("BackupPath: %v", err)
	}
	if want := src + ".bak.1"; dst != want {
		t.Errorf("dst = %q, want %q", dst, want)
	}
}

func TestBackupPath_File_PicksSmallestFreeSuffix(t *testing.T) {
	// Why: spec §607 requires *smallest* free numeric suffix, not the
	// next-after-highest. With .bak and .bak.2 occupied but .bak.1
	// free, BackupPath must pick .bak.1.
	fs := newFakeFS()
	src := "/proj/u-boot.yaml"
	for _, p := range []string{src, src + ".bak", src + ".bak.2"} {
		if err := fs.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatalf("setup %s: %v", p, err)
		}
	}

	dst, err := application.BackupPath(fs, src)
	if err != nil {
		t.Fatalf("BackupPath: %v", err)
	}
	if want := src + ".bak.1"; dst != want {
		t.Errorf("dst = %q, want %q", dst, want)
	}
}

func TestBackupPath_File_ContentPreserved(t *testing.T) {
	fs := newFakeFS()
	src := "/proj/compose.yaml"
	body := []byte("services:\n  app:\n    image: foo\n")
	if err := fs.WriteFile(src, body, 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	dst, err := application.BackupPath(fs, src)
	if err != nil {
		t.Fatalf("BackupPath: %v", err)
	}
	got, err := fs.ReadFile(dst)
	if err != nil {
		t.Fatalf("ReadFile(%s): %v", dst, err)
	}
	if string(got) != string(body) {
		t.Errorf("backup body = %q, want %q", got, body)
	}
}

func TestBackupPath_Directory_RecursiveCopy(t *testing.T) {
	fs := newFakeFS()
	src := "/proj/docker"
	fs.markDirExists(src)
	if err := fs.WriteFile(filepath.Join(src, "Dockerfile"), []byte("FROM scratch\n"), 0o644); err != nil {
		t.Fatalf("setup Dockerfile: %v", err)
	}
	nested := filepath.Join(src, "scripts")
	if err := fs.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("setup nested mkdir: %v", err)
	}
	if err := fs.WriteFile(filepath.Join(nested, "entry.sh"), []byte("#!/bin/sh\n"), 0o644); err != nil {
		t.Fatalf("setup entry.sh: %v", err)
	}

	dst, err := application.BackupPath(fs, src)
	if err != nil {
		t.Fatalf("BackupPath: %v", err)
	}
	if want := src + ".bak"; dst != want {
		t.Errorf("dst = %q, want %q", dst, want)
	}

	for srcPath, wantContent := range map[string]string{
		filepath.Join(dst, "Dockerfile"):           "FROM scratch\n",
		filepath.Join(dst, "scripts", "entry.sh"): "#!/bin/sh\n",
	} {
		got, err := fs.ReadFile(srcPath)
		if err != nil {
			t.Errorf("ReadFile(%s): %v", srcPath, err)
			continue
		}
		if string(got) != wantContent {
			t.Errorf("ReadFile(%s) = %q, want %q", srcPath, got, wantContent)
		}
	}
}

func TestBackupPath_MissingSourceReturnsErr(t *testing.T) {
	fs := newFakeFS()
	_, err := application.BackupPath(fs, "/proj/does-not-exist")
	if err == nil {
		t.Fatalf("BackupPath(missing): expected error, got nil")
	}
	if !errors.Is(err, application.ErrBackupSourceMissing) {
		t.Errorf("BackupPath(missing): error %v does not wrap ErrBackupSourceMissing", err)
	}
}

func TestBackupPath_TreeCopyFailure_RollsBack(t *testing.T) {
	// Why: spec §608 requires a rollback when a tree-backup fails
	// partway. Setup: directory with two files; force WriteFile to
	// fail on the second one; assert (a) BackupPath returns the
	// underlying error and (b) the partial destination is gone after
	// rollback.
	fs := newFakeFS()
	src := "/proj/docs"
	fs.markDirExists(src)
	if err := fs.WriteFile(filepath.Join(src, "a.md"), []byte("a"), 0o644); err != nil {
		t.Fatalf("setup a: %v", err)
	}
	if err := fs.WriteFile(filepath.Join(src, "b.md"), []byte("b"), 0o644); err != nil {
		t.Fatalf("setup b: %v", err)
	}

	wantErr := errors.New("disk full on b")
	fs.failOn = filepath.Join(src+".bak", "b.md")
	fs.failErr = wantErr

	_, err := application.BackupPath(fs, src)
	if err == nil {
		t.Fatalf("BackupPath: expected error, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("BackupPath: error %v does not wrap %v", err, wantErr)
	}

	// Rollback: partial destination removed.
	exists, _ := fs.Exists(src + ".bak")
	if exists {
		t.Errorf("rollback: %s.bak still exists", src)
	}
	if exists, _ := fs.Exists(filepath.Join(src+".bak", "a.md")); exists {
		t.Errorf("rollback: leaked a.md in partial backup")
	}
}

func TestBackupPath_ExhaustedSuffixReturnsErr(t *testing.T) {
	fs := newFakeFS()
	src := "/proj/u-boot.yaml"
	// Source + .bak + .bak.1 ... .bak.1000 all occupied.
	if err := fs.WriteFile(src, []byte("x"), 0o644); err != nil {
		t.Fatalf("setup src: %v", err)
	}
	if err := fs.WriteFile(src+".bak", []byte("x"), 0o644); err != nil {
		t.Fatalf("setup .bak: %v", err)
	}
	for i := 1; i <= 1000; i++ {
		p := src + ".bak." + strconv.Itoa(i)
		if err := fs.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatalf("setup %s: %v", p, err)
		}
	}

	_, err := application.BackupPath(fs, src)
	if err == nil {
		t.Fatalf("BackupPath: expected exhaustion error, got nil")
	}
	if !errors.Is(err, application.ErrBackupSuffixExhausted) {
		t.Errorf("BackupPath: error %v does not wrap ErrBackupSuffixExhausted", err)
	}
}

func TestBackupPath_ReadFileErrorPropagates(t *testing.T) {
	fs := newFakeFS()
	src := "/proj/u-boot.yaml"
	if err := fs.WriteFile(src, []byte("x"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	wantErr := errors.New("read denied")
	fs.failReadOn = src
	fs.failReadErr = wantErr

	_, err := application.BackupPath(fs, src)
	if err == nil {
		t.Fatalf("BackupPath: expected error, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("BackupPath: error %v does not wrap %v", err, wantErr)
	}
}

func TestBackupPath_ExistsErrorOnSrcPropagates(t *testing.T) {
	fs := newFakeFS()
	src := "/proj/u-boot.yaml"
	wantErr := errors.New("stat denied")
	fs.failExistsOn = src
	fs.failExistsErr = wantErr

	_, err := application.BackupPath(fs, src)
	if err == nil {
		t.Fatalf("BackupPath: expected error, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("BackupPath: error %v does not wrap %v", err, wantErr)
	}
}

func TestBackupPath_ExistsErrorOnCandidatePropagates(t *testing.T) {
	fs := newFakeFS()
	src := "/proj/u-boot.yaml"
	if err := fs.WriteFile(src, []byte("x"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	wantErr := errors.New("stat denied")
	fs.failExistsOn = src + ".bak"
	fs.failExistsErr = wantErr

	_, err := application.BackupPath(fs, src)
	if err == nil {
		t.Fatalf("BackupPath: expected error, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("BackupPath: error %v does not wrap %v", err, wantErr)
	}
}

func TestBackupPath_ExistsErrorOnNumericCandidatePropagates(t *testing.T) {
	// Why: covers the chooseBackupPath loop body (numeric-suffix
	// path) — the first-suffix Exists succeeds, but the .bak.1
	// candidate fails.
	fs := newFakeFS()
	src := "/proj/u-boot.yaml"
	if err := fs.WriteFile(src, []byte("x"), 0o644); err != nil {
		t.Fatalf("setup src: %v", err)
	}
	if err := fs.WriteFile(src+".bak", []byte("x"), 0o644); err != nil {
		t.Fatalf("setup .bak: %v", err)
	}
	wantErr := errors.New("stat denied")
	fs.failExistsOn = src + ".bak.1"
	fs.failExistsErr = wantErr

	_, err := application.BackupPath(fs, src)
	if err == nil {
		t.Fatalf("BackupPath: expected error, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("BackupPath: error %v does not wrap %v", err, wantErr)
	}
}

func TestBackupPath_IsDirErrorPropagates(t *testing.T) {
	fs := newFakeFS()
	src := "/proj/u-boot.yaml"
	if err := fs.WriteFile(src, []byte("x"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	wantErr := errors.New("isdir denied")
	fs.failIsDirOn = src
	fs.failIsDirErr = wantErr

	_, err := application.BackupPath(fs, src)
	if err == nil {
		t.Fatalf("BackupPath: expected error, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("BackupPath: error %v does not wrap %v", err, wantErr)
	}
}
