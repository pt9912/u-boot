// Package fs is the real-filesystem implementation of the
// `port/driven.FileSystem` interface (LH-FA-ARCH-002).
//
// Layer rule: adapters may import the domain and their driven-port
// interface, plus external libraries; they may not import application
// or other adapter packages (LH-FA-ARCH-003, depguard-enforced).
package fs

import (
	"errors"
	"io"
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

// Copy streams src to dst (non-exclusive: truncate-overwrites
// existing dst). Parent directories are created with mode 0o755 like
// [WriteFile]. Memory footprint is bounded by io.Copy's internal
// buffer (~32 KiB) regardless of file size.
func (FS) Copy(src, dst string, mode iofs.FileMode) error {
	return streamCopy(src, dst, mode, os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
}

// CopyExclusive streams src to dst with O_CREATE|O_EXCL — fails
// fast with a wrapped fs.ErrExist if dst already exists. Companion
// to [WriteFileExclusive].
func (FS) CopyExclusive(src, dst string, mode iofs.FileMode) error {
	return streamCopy(src, dst, mode, os.O_CREATE|os.O_EXCL|os.O_WRONLY)
}

// streamCopy is the shared implementation of [Copy] and
// [CopyExclusive]. The only difference is the open-flag for the
// destination file; the rest (mkdir parents, open source, io.Copy,
// close-with-error-join) is identical.
//
// Close-on-error: when io.Copy fails midway, the destination must
// still be closed (otherwise the file handle leaks). errors.Join
// surfaces both the copy error and any close error so the operator
// can see both.
func streamCopy(src, dst string, mode iofs.FileMode, flag int) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }() // read-side close error is non-actionable
	out, err := os.OpenFile(dst, flag, mode)
	if err != nil {
		return err
	}
	if _, copyErr := io.Copy(out, in); copyErr != nil {
		if closeErr := out.Close(); closeErr != nil {
			return errors.Join(copyErr, closeErr)
		}
		return copyErr
	}
	return out.Close()
}
