// Package fs is the real-filesystem implementation of the
// `port/driven.FileSystem` interface (LH-FA-ARCH-002).
//
// Layer rule: adapters may import the domain and their driven-port
// interface, plus external libraries; they may not import application
// or other adapter packages (LH-FA-ARCH-003, depguard-enforced).
package fs

import (
	"errors"
	iofs "io/fs"
	"os"
	"path/filepath"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// FS is the production filesystem adapter. It delegates to the
// standard library; the implementation lives behind the port interface
// so application-layer tests can substitute a fake without touching
// disk.
type FS struct{}

// Static check: FS satisfies the FileSystem port. The line lives in
// the adapter (not in a `_test.go` file) so a mismatch breaks the
// package build, not only the test build.
var _ driven.FileSystem = (*FS)(nil)

// New returns a ready-to-use FS adapter.
func New() *FS { return &FS{} }

// Exists reports whether path exists. It distinguishes
// "does not exist" (returns `(false, nil)`) from a real I/O error
// (returns `(false, err)`).
func (FS) Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, iofs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// ReadFile mirrors os.ReadFile.
func (FS) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile creates parent directories with the canonical
// project-directory mode 0o755 (LH-FA-INIT-003 — directories are
// shared with collaborators and CI runners, neither benefits from a
// more restrictive default) and writes the file itself with the
// caller-supplied mode. Truncate-overwrites an existing file.
func (FS) WriteFile(path string, data []byte, mode iofs.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, mode)
}

// WriteFileExclusive uses O_CREATE|O_EXCL|O_WRONLY so the write fails
// with a wrapped os.ErrExist (which is fs.ErrExist) when path already
// exists. Parent directories are created with mode 0o755 like
// [WriteFile]. The os.OpenFile + Write + Close path is the
// canonical Go way to express atomic-create-then-write.
func (FS) WriteFileExclusive(path string, data []byte, mode iofs.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	if _, writeErr := f.Write(data); writeErr != nil {
		closeErr := f.Close()
		if closeErr != nil {
			return errors.Join(writeErr, closeErr)
		}
		return writeErr
	}
	return f.Close()
}

// Mkdir mirrors os.Mkdir — single directory, no parents, fails with
// fs.ErrExist when path is taken. Use [MkdirAll] when idempotent
// semantics are wanted.
func (FS) Mkdir(path string, mode iofs.FileMode) error {
	return os.Mkdir(path, mode)
}

// MkdirAll mirrors os.MkdirAll. Idempotent.
func (FS) MkdirAll(path string, mode iofs.FileMode) error {
	return os.MkdirAll(path, mode)
}

// Rename mirrors os.Rename.
func (FS) Rename(src, dst string) error {
	return os.Rename(src, dst)
}

// ReadDir mirrors os.ReadDir.
func (FS) ReadDir(path string) ([]iofs.DirEntry, error) {
	return os.ReadDir(path)
}

// Lstat mirrors os.Lstat — does not follow symlinks. The backup
// strategy relies on this so symlinks are detectable via
// `info.Mode()&fs.ModeSymlink != 0`.
func (FS) Lstat(path string) (iofs.FileInfo, error) {
	return os.Lstat(path)
}

// RemoveAll mirrors os.RemoveAll.
func (FS) RemoveAll(path string) error {
	return os.RemoveAll(path)
}
