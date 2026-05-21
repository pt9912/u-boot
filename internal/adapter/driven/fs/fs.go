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

// WriteFile creates parent directories with mode 0o755 and writes the
// file with the given mode. The write itself is non-atomic.
func (FS) WriteFile(path string, data []byte, mode iofs.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, mode)
}

// MkdirAll mirrors os.MkdirAll.
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
