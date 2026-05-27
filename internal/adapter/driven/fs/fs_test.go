package fs_test

import (
	"errors"
	iofs "io/fs"
	"os"
	"path/filepath"
	"runtime"
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

func TestFS_Lstat(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "f.txt")
	if err := os.WriteFile(file, []byte("hi"), 0o600); err != nil {
		t.Fatalf("setup: %v", err)
	}
	adapter := fs.New()

	info, err := adapter.Lstat(dir)
	if err != nil {
		t.Fatalf("Lstat(dir): %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("Lstat(dir).IsDir() = false, want true")
	}

	info, err = adapter.Lstat(file)
	if err != nil {
		t.Fatalf("Lstat(file): %v", err)
	}
	if info.IsDir() {
		t.Fatalf("Lstat(file).IsDir() = true, want false")
	}
	if info.Size() != 2 {
		t.Errorf("Lstat(file).Size() = %d, want 2", info.Size())
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("Lstat(file).Mode().Perm() = %o, want 0o600", info.Mode().Perm())
	}

	_, err = adapter.Lstat(filepath.Join(dir, "missing"))
	if !errors.Is(err, iofs.ErrNotExist) {
		t.Fatalf("Lstat(missing): want ErrNotExist, got %v", err)
	}
}

func TestFS_Lstat_DoesNotFollowSymlink(t *testing.T) {
	// Why: pins the no-follow semantics that the LH-FA-INIT-005
	// backup strategy relies on for symlink detection. A naive
	// os.Stat-based impl would silently report the link's target.
	dir := t.TempDir()
	target := filepath.Join(dir, "target.txt")
	if err := os.WriteFile(target, []byte("hi"), 0o644); err != nil {
		t.Fatalf("setup target: %v", err)
	}
	link := filepath.Join(dir, "link.txt")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlink unsupported on this platform: %v", err)
	}

	info, err := fs.New().Lstat(link)
	if err != nil {
		t.Fatalf("Lstat(link): %v", err)
	}
	if info.Mode()&iofs.ModeSymlink == 0 {
		t.Errorf("Lstat(link).Mode() = %v, want ModeSymlink bit set", info.Mode())
	}
}

func TestFS_WriteFileExclusive_FailsOnExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "f.txt")
	if err := os.WriteFile(path, []byte("first"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	err := fs.New().WriteFileExclusive(path, []byte("second"), 0o644)
	if !errors.Is(err, iofs.ErrExist) {
		t.Fatalf("WriteFileExclusive(existing): want ErrExist, got %v", err)
	}
	// Original content untouched.
	got, _ := os.ReadFile(path)
	if string(got) != "first" {
		t.Errorf("file content = %q, want %q (exclusive write must not clobber)", got, "first")
	}
}

func TestFS_WriteFileExclusive_CreatesNew(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "deep", "f.txt")

	if err := fs.New().WriteFileExclusive(path, []byte("payload"), 0o600); err != nil {
		t.Fatalf("WriteFileExclusive: %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if string(got) != "payload" {
		t.Errorf("content = %q, want %q", got, "payload")
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("mode = %o, want 0o600", info.Mode().Perm())
	}
}

func TestFS_Mkdir_FailsOnExisting(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "sub")
	if err := os.Mkdir(target, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	err := fs.New().Mkdir(target, 0o755)
	if !errors.Is(err, iofs.ErrExist) {
		t.Fatalf("Mkdir(existing): want ErrExist, got %v", err)
	}
}

func TestFS_Mkdir_CreatesNew(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "sub")

	if err := fs.New().Mkdir(target, 0o755); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}
	info, err := os.Stat(target)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if !info.IsDir() {
		t.Errorf("not a directory")
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

func TestFS_Copy(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.bin")
	body := []byte("u-boot streaming copy test\n")
	if err := os.WriteFile(src, body, 0o600); err != nil {
		t.Fatalf("setup: %v", err)
	}
	dst := filepath.Join(dir, "nested", "dst.bin")
	adapter := fs.New()

	if err := adapter.Copy(src, dst, 0o644); err != nil {
		t.Fatalf("Copy: %v", err)
	}
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("ReadFile dst: %v", err)
	}
	if string(got) != string(body) {
		t.Errorf("dst contents = %q, want %q", got, body)
	}
	info, err := os.Stat(dst)
	if err != nil {
		t.Fatalf("Stat dst: %v", err)
	}
	if info.Mode().Perm() != 0o644 {
		t.Errorf("dst mode = %v, want 0o644", info.Mode().Perm())
	}
}

func TestFS_CopyExclusive_FailsOnExisting(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	dst := filepath.Join(dir, "dst")
	if err := os.WriteFile(src, []byte("x"), 0o644); err != nil {
		t.Fatalf("setup src: %v", err)
	}
	if err := os.WriteFile(dst, []byte("existing"), 0o644); err != nil {
		t.Fatalf("setup dst: %v", err)
	}
	adapter := fs.New()

	err := adapter.CopyExclusive(src, dst, 0o644)
	if !errors.Is(err, iofs.ErrExist) {
		t.Errorf("CopyExclusive on existing: err = %v, want wrapped ErrExist", err)
	}
}

func TestFS_Copy_LargeFile_BoundedMemoryFootprint(t *testing.T) {
	// Why: streaming-copy must not scale in-process memory with file
	// size. Sets up a sparse 1-GiB source via os.Truncate (the kernel
	// records the hole; no 1 GiB of bytes actually written) and
	// measures heap-alloc-delta around the Copy call via
	// runtime.MemStats. A streaming implementation should stay well
	// under 100 MiB of new allocation (io.Copy's internal buffer is
	// 32 KiB per iteration, all reusable); a ReadFile+WriteFile
	// re-regression would allocate ~1 GiB and trip the assert.
	//
	// The threshold (100 MiB) is generous to avoid flake from GC
	// timing — the real ReadFile+WriteFile path would be off by an
	// order of magnitude, so a one-order-of-magnitude check is
	// reliable.
	dir := t.TempDir()
	src := filepath.Join(dir, "huge.bin")
	f, err := os.Create(src)
	if err != nil {
		t.Fatalf("create huge: %v", err)
	}
	if err := f.Truncate(1 << 30); err != nil { // 1 GiB
		_ = f.Close()
		t.Fatalf("truncate huge: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close huge: %v", err)
	}
	dst := filepath.Join(dir, "huge.bak")
	adapter := fs.New()

	runtime.GC()
	var before runtime.MemStats
	runtime.ReadMemStats(&before)

	if err := adapter.Copy(src, dst, 0o644); err != nil {
		t.Fatalf("Copy huge: %v", err)
	}

	var after runtime.MemStats
	runtime.ReadMemStats(&after)
	const memBound = 100 << 20 // 100 MiB
	if delta := after.TotalAlloc - before.TotalAlloc; delta > memBound {
		t.Errorf("Copy(1 GiB) allocated %d bytes (cap %d) — streaming regression?", delta, memBound)
	}

	info, err := os.Stat(dst)
	if err != nil {
		t.Fatalf("stat dst: %v", err)
	}
	if info.Size() != 1<<30 {
		t.Errorf("dst size = %d, want %d", info.Size(), 1<<30)
	}
}

// The static FS↔driven.FileSystem contract check lives in fs.go (see
// `var _ driven.FileSystem = (*FS)(nil)`), not here.
