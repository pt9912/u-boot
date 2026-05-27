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
	// are an adapter concern. Truncate-overwrites an existing file —
	// see [WriteFileExclusive] for the race-safe variant.
	WriteFile(path string, data []byte, mode fs.FileMode) error

	// WriteFileExclusive writes data to path with O_CREATE|O_EXCL
	// semantics — it succeeds only if path did not yet exist. Returns
	// a wrapped fs.ErrExist when the slot is taken. The LH-FA-INIT-005
	// backup strategy uses this to close the TOCTOU window between
	// suffix selection (`<src>.bak.N` chosen via [Exists]) and the
	// actual write, so concurrent runs cannot clobber each other's
	// fresh backups (spec §607: "ohne vorhandene Backups zu
	// überschreiben").
	WriteFileExclusive(path string, data []byte, mode fs.FileMode) error

	// Mkdir creates a single directory at path with O_EXCL semantics,
	// mirroring os.Mkdir. Returns a wrapped fs.ErrExist when the path
	// is taken. Companion to [WriteFileExclusive]; the LH-FA-INIT-005
	// backup strategy uses it to reserve the top-level <src>.bak
	// directory atomically.
	Mkdir(path string, mode fs.FileMode) error

	// MkdirAll creates the directory and all parents with the given
	// mode, mirroring os.MkdirAll. Idempotent — use [Mkdir] when you
	// need exclusive create semantics.
	MkdirAll(path string, mode fs.FileMode) error

	// Rename moves src to dst, mirroring os.Rename. Used by the backup
	// strategy in LH-FA-INIT-005.
	Rename(src, dst string) error

	// ReadDir lists the directory entries at path.
	ReadDir(path string) ([]fs.DirEntry, error)

	// Lstat returns file info for path without following symlinks
	// (mirrors os.Lstat). Callers consume Mode() for the kind probe
	// (regular file / directory / symlink) and for mode preservation;
	// fs.ErrNotExist is returned wrapped for a missing path so callers
	// can branch on errors.Is. The LH-FA-INIT-005 backup strategy
	// uses Lstat (not Stat) so that symlinks are detectable and can
	// be rejected with ErrBackupUnsupportedKind rather than silently
	// followed.
	Lstat(path string) (fs.FileInfo, error)

	// RemoveAll deletes path and any children, mirroring os.RemoveAll.
	// Used by the LH-FA-INIT-005 backup strategy as the rollback action
	// when a partial tree copy fails partway through.
	RemoveAll(path string) error

	// Copy streams src to dst, creating parent directories with mode
	// 0o755 as needed and writing dst with the given mode. Truncate-
	// overwrites an existing dst. Used by the LH-FA-INIT-005 backup
	// strategy for nested file copies inside an already-reserved top-
	// level backup dir. The streaming form replaces the prior
	// ReadFile+WriteFile pair so the memory footprint stays bounded by
	// io.Copy's internal buffer (typically 32 KiB) regardless of file
	// size; the 256-MiB safety cap of the M3-T4a-MVP is gone.
	Copy(src, dst string, mode fs.FileMode) error

	// CopyExclusive streams src to dst with O_CREATE|O_EXCL semantics —
	// it succeeds only if dst did not yet exist. Returns a wrapped
	// fs.ErrExist when the slot is taken. The LH-FA-INIT-005 backup
	// strategy uses this for the top-level <src>.bak file slot to
	// close the TOCTOU window between suffix selection and the actual
	// write.
	CopyExclusive(src, dst string, mode fs.FileMode) error
}
