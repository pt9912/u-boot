package fs_test

import (
	"errors"
	iofs "io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driven/fs"
)

func TestFS_ReadFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "payload.txt")
	if err := os.WriteFile(path, []byte("hello"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	got, err := fs.New().ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != "hello" {
		t.Fatalf("ReadFile = %q, want %q", got, "hello")
	}

	_, err = fs.New().ReadFile(filepath.Join(dir, "missing.txt"))
	if err == nil {
		t.Fatalf("ReadFile(missing): expected error, got nil")
	}
}

func TestFS_Exists(t *testing.T) {
	dir := t.TempDir()
	existing := filepath.Join(dir, "present.txt")
	if err := os.WriteFile(existing, []byte("hi"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	adapter := fs.New()

	got, err := adapter.Exists(existing)
	if err != nil {
		t.Fatalf("Exists(present): %v", err)
	}
	if !got {
		t.Fatalf("Exists(present) = false, want true")
	}

	got, err = adapter.Exists(filepath.Join(dir, "missing.txt"))
	if err != nil {
		t.Fatalf("Exists(missing): unexpected error: %v", err)
	}
	if got {
		t.Fatalf("Exists(missing) = true, want false")
	}
}

func TestFS_WriteFile_CreatesParents(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "deep", "nested", "file.txt")

	if err := fs.New().WriteFile(target, []byte("payload"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if string(got) != "payload" {
		t.Fatalf("WriteFile payload = %q, want %q", got, "payload")
	}
}

func TestFS_MkdirAll_Idempotent(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "a", "b", "c")
	adapter := fs.New()

	if err := adapter.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("MkdirAll first: %v", err)
	}
	if err := adapter.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("MkdirAll second (idempotent): %v", err)
	}
}

func TestFS_Rename(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")
	if err := os.WriteFile(src, []byte("x"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := fs.New().Rename(src, dst); err != nil {
		t.Fatalf("Rename: %v", err)
	}

	if _, err := os.Stat(src); !errors.Is(err, iofs.ErrNotExist) {
		t.Fatalf("Rename: src still exists, err=%v", err)
	}
	if _, err := os.Stat(dst); err != nil {
		t.Fatalf("Rename: dst missing, err=%v", err)
	}
}

func TestFS_Rename_MissingSourceReturnsError(t *testing.T) {
	// Why: the backup strategy in LH-FA-INIT-005 must be able to tell
	// "no file to back up" from "the OS swallowed our error". Pin the
	// error path explicitly.
	dir := t.TempDir()
	src := filepath.Join(dir, "missing.txt")
	dst := filepath.Join(dir, "dst.txt")

	err := fs.New().Rename(src, dst)
	if err == nil {
		t.Fatalf("Rename(missing src): expected error, got nil")
	}
	if !errors.Is(err, iofs.ErrNotExist) {
		t.Fatalf("Rename(missing src): error %v does not wrap fs.ErrNotExist", err)
	}
}

func TestFS_ReadDir(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"a.txt", "b.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), nil, 0o644); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	entries, err := fs.New().ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("ReadDir len = %d, want 2", len(entries))
	}
}

func TestFS_IsDir(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "f.txt")
	if err := os.WriteFile(file, nil, 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	adapter := fs.New()

	got, err := adapter.IsDir(dir)
	if err != nil {
		t.Fatalf("IsDir(dir): %v", err)
	}
	if !got {
		t.Fatalf("IsDir(dir) = false, want true")
	}

	got, err = adapter.IsDir(file)
	if err != nil {
		t.Fatalf("IsDir(file): %v", err)
	}
	if got {
		t.Fatalf("IsDir(file) = true, want false")
	}

	got, err = adapter.IsDir(filepath.Join(dir, "missing"))
	if err != nil {
		t.Fatalf("IsDir(missing): unexpected error: %v", err)
	}
	if got {
		t.Fatalf("IsDir(missing) = true, want false")
	}
}

func TestFS_RemoveAll(t *testing.T) {
	dir := t.TempDir()
	tree := filepath.Join(dir, "a", "b")
	if err := os.MkdirAll(tree, 0o755); err != nil {
		t.Fatalf("setup mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tree, "c.txt"), []byte("x"), 0o644); err != nil {
		t.Fatalf("setup write: %v", err)
	}
	adapter := fs.New()

	if err := adapter.RemoveAll(filepath.Join(dir, "a")); err != nil {
		t.Fatalf("RemoveAll: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "a")); !errors.Is(err, iofs.ErrNotExist) {
		t.Fatalf("RemoveAll: tree still exists, err=%v", err)
	}

	if err := adapter.RemoveAll(filepath.Join(dir, "missing")); err != nil {
		t.Fatalf("RemoveAll(missing): want nil (idempotent), got %v", err)
	}
}

// The static FS↔driven.FileSystem contract check lives in fs.go (see
// `var _ driven.FileSystem = (*FS)(nil)`), not here.
