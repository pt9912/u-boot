// Package driven holds the driven-port interfaces of u-boot — the
// abstractions through which the application layer reaches the outside
// world (filesystem, YAML codec, git, clock, …). Concrete
// implementations live in `internal/adapter/driven/`.
//
// Layer rules (LH-FA-ARCH-002, LH-FA-ARCH-003): driven ports may only
// depend on `internal/hexagon/domain` and the Go standard library.
package driven

import "io/fs"

// FileSystem abstracts the file-system operations the application
// layer needs. It is intentionally small — every method maps to one
// LH-FA-INIT-005 / LH-FA-INIT-003 use case so fakes stay trivial to
// write in tests.
type FileSystem interface {
	// Exists reports whether the given path exists. It returns
	// `(false, nil)` for a non-existent path and `(_, err)` only on
	// non-categorical errors (e.g. permission denied on a parent dir).
	Exists(path string) (bool, error)

	// ReadFile reads the file at path. It mirrors os.ReadFile.
	ReadFile(path string) ([]byte, error)

	// WriteFile writes data to the file at path, creating parent
	// directories with mode 0o755 as needed and the file itself with
	// the given mode. It is non-atomic at this layer; atomic semantics
	// are an adapter concern.
	WriteFile(path string, data []byte, mode fs.FileMode) error

	// MkdirAll creates the directory and all parents with the given
	// mode, mirroring os.MkdirAll. Idempotent.
	MkdirAll(path string, mode fs.FileMode) error

	// Rename moves src to dst, mirroring os.Rename. Used by the backup
	// strategy in LH-FA-INIT-005.
	Rename(src, dst string) error

	// ReadDir lists the directory entries at path.
	ReadDir(path string) ([]fs.DirEntry, error)

	// IsDir reports whether path exists and is a directory. Returns
	// `(false, nil)` for a non-existent path so callers can use it as a
	// "kind probe" without a separate Exists call. Real I/O errors
	// (permission denied on a parent, etc.) propagate. Added for the
	// LH-FA-INIT-005 backup strategy, which copies file-vs-directory
	// trees differently.
	IsDir(path string) (bool, error)

	// RemoveAll deletes path and any children, mirroring os.RemoveAll.
	// Used by the LH-FA-INIT-005 backup strategy as the rollback action
	// when a partial tree copy fails partway through.
	RemoveAll(path string) error
}
