package application

import (
	"errors"
	"fmt"
	iofs "io/fs"
	"path/filepath"
	"strconv"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// backupSuffixCap caps the .bak.N suffix search (and the
// race-retry budget) at 1000. A user who has accumulated this many
// stale backups under one path must clean them up manually; the cap
// exists so a runaway loop cannot happen even if every Exists check
// keeps returning true (e.g. a hostile sibling process). The same
// budget bounds the TOCTOU race-retry — see [BackupPath].
const backupSuffixCap = 1000

// maxBackupFileSize caps individual files at 256 MiB. The current
// FileSystem port reads files via ReadFile (full content into
// memory), so a multi-GB asset under `docs/` or `docker/` would OOM
// the process. The cap is documented as a temporary MVP carveout in
// `docs/plan/planning/in-progress/carveouts.md` and lifted by
// `slice-v1-backup-streaming-copy.md` once a streaming copy
// primitive lands on the FileSystem port.
const maxBackupFileSize = 256 << 20

// BackupPath copies src to a sibling backup path and returns the
// chosen backup path. Suffix selection follows LH-FA-INIT-005 §607:
// <src>.bak first, then <src>.bak.1, .bak.2, ... — smallest free
// numeric suffix, existing backups are never overwritten.
//
// File-vs-directory dispatch comes from [driven.FileSystem.Lstat]
// (Lstat, not Stat, so symlinks are detectable). Symlinks are
// rejected with [driving.ErrBackupUnsupportedKind] — LH-FA-INIT-005
// §608 does not specify symlink semantics, and silently following
// would surprise users who symlink shared assets into the project.
//
// TOCTOU: two concurrent runs can both pick the same `.bak.N` slot.
// BackupPath defends against this by using the exclusive-create
// primitives ([driven.FileSystem.WriteFileExclusive] for files,
// [driven.FileSystem.Mkdir] for the top-level backup dir) and
// retrying on fs.ErrExist — chooseBackupPath will skip the now-
// occupied slot on the next attempt. The retry budget is shared
// with the suffix cap, so the worst case is the same
// [driving.ErrBackupSuffixExhausted] outcome.
//
// On partial directory-copy failure, BackupPath rolls back by
// removing the partial destination (LH-FA-INIT-005 §608 explicitly
// requires this; POSIX atomicity for recursive trees is not
// guaranteed). A non-nil rollback error is joined to the original
// via errors.Join so the operator can see both — silently swallowing
// it would leave half-written .bak trees lying around.
//
// Mode preservation: file/dir modes come from the source's Lstat
// info, so a `0o755 scripts/entry.sh` is backed up `0o755`. A
// future re-installer that restores the backup keeps the
// executable bit.
//
// Returns wrapped sentinels in [driving] so the CLI adapter can map
// each to its LH-FA-CLI-006 exit code without violating the
// adapter→application depguard rule.
func BackupPath(fs driven.FileSystem, src string) (string, error) {
	info, err := fs.Lstat(src)
	if err != nil {
		if errors.Is(err, iofs.ErrNotExist) {
			return "", fmt.Errorf("%w: %s", driving.ErrBackupSourceMissing, src)
		}
		return "", fmt.Errorf("lstat %s: %w", src, err)
	}
	if info.Mode()&iofs.ModeSymlink != 0 {
		return "", fmt.Errorf("%w: %s is a symlink", driving.ErrBackupUnsupportedKind, src)
	}
	if !info.IsDir() && info.Size() > maxBackupFileSize {
		return "", fmt.Errorf("%w: %s is %d bytes (cap %d)",
			driving.ErrBackupTooLarge, src, info.Size(), maxBackupFileSize)
	}

	for attempt := 0; attempt < backupSuffixCap; attempt++ {
		dst, err := chooseBackupPath(fs, src)
		if err != nil {
			return "", err
		}
		copyErr := createBackup(fs, src, dst, info)
		if copyErr == nil {
			return dst, nil
		}
		if errors.Is(copyErr, iofs.ErrExist) {
			// Race-retry: another process took dst between our
			// chooseBackupPath and our exclusive create. Loop —
			// chooseBackupPath will pick a different slot.
			continue
		}
		// Real failure during dir copy → rollback the partial tree.
		// File backups are atomic (WriteFileExclusive), so no rollback
		// is needed for them.
		if info.IsDir() {
			if rmErr := fs.RemoveAll(dst); rmErr != nil {
				return "", errors.Join(copyErr,
					fmt.Errorf("rollback %s: %w", dst, rmErr))
			}
		}
		return "", copyErr
	}
	return "", fmt.Errorf("%w: %s.bak exhausted after %d race retries",
		driving.ErrBackupSuffixExhausted, src, backupSuffixCap)
}

// chooseBackupPath returns the smallest-suffix backup destination
// that does not yet exist: <src>.bak, then .bak.1, .bak.2, ...
// Iteration is bounded by backupSuffixCap.
func chooseBackupPath(fs driven.FileSystem, src string) (string, error) {
	candidate := src + ".bak"
	exists, err := fs.Exists(candidate)
	if err != nil {
		return "", fmt.Errorf("check %s: %w", candidate, err)
	}
	if !exists {
		return candidate, nil
	}
	for i := 1; i <= backupSuffixCap; i++ {
		candidate = src + ".bak." + strconv.Itoa(i)
		exists, err := fs.Exists(candidate)
		if err != nil {
			return "", fmt.Errorf("check %s: %w", candidate, err)
		}
		if !exists {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("%w: %s.bak through %s.bak.%d are all occupied",
		driving.ErrBackupSuffixExhausted, src, src, backupSuffixCap)
}

// createBackup dispatches to the file or directory backup path,
// preserving src's mode. Returns iofs.ErrExist (wrapped) when a
// concurrent process beat us to the slot.
func createBackup(fs driven.FileSystem, src, dst string, info iofs.FileInfo) error {
	if info.IsDir() {
		return createBackupDir(fs, src, dst, info.Mode().Perm())
	}
	return createBackupFile(fs, src, dst, info.Mode().Perm())
}

// createBackupFile copies src into dst with O_EXCL semantics — fails
// fast with iofs.ErrExist on TOCTOU collisions.
func createBackupFile(fs driven.FileSystem, src, dst string, mode iofs.FileMode) error {
	data, err := fs.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read %s: %w", src, err)
	}
	if err := fs.WriteFileExclusive(dst, data, mode); err != nil {
		return fmt.Errorf("write %s: %w", dst, err)
	}
	return nil
}

// createBackupDir reserves the top-level dst directory atomically
// via [driven.FileSystem.Mkdir] (returns iofs.ErrExist on race),
// then recurses to copy children with [copyTreeContents].
func createBackupDir(fs driven.FileSystem, src, dst string, mode iofs.FileMode) error {
	if err := fs.Mkdir(dst, mode); err != nil {
		return fmt.Errorf("mkdir %s: %w", dst, err)
	}
	return copyTreeContents(fs, src, dst)
}

// copyTreeContents walks src and replicates its contents into dst.
// Per-entry work is delegated to [copyTreeEntry] so the top-level
// loop stays inside the gocognit budget; nested dirs and files use
// MkdirAll / WriteFile (non-exclusive) because they live inside the
// top-level dst directory that [createBackupDir] already reserved
// atomically.
func copyTreeContents(fs driven.FileSystem, src, dst string) error {
	entries, err := fs.ReadDir(src)
	if err != nil {
		return fmt.Errorf("read dir %s: %w", src, err)
	}
	for _, entry := range entries {
		if err := copyTreeEntry(fs, src, dst, entry); err != nil {
			return err
		}
	}
	return nil
}

// copyTreeEntry dispatches a single ReadDir entry to the file or
// directory copy path. Symlinks are rejected with the same sentinel
// used for the top-level case ([driving.ErrBackupUnsupportedKind]).
func copyTreeEntry(fs driven.FileSystem, srcDir, dstDir string, entry iofs.DirEntry) error {
	srcChild := filepath.Join(srcDir, entry.Name())
	dstChild := filepath.Join(dstDir, entry.Name())
	info, err := fs.Lstat(srcChild)
	if err != nil {
		return fmt.Errorf("lstat %s: %w", srcChild, err)
	}
	if info.Mode()&iofs.ModeSymlink != 0 {
		return fmt.Errorf("%w: %s is a symlink", driving.ErrBackupUnsupportedKind, srcChild)
	}
	if entry.IsDir() {
		return copyTreeNestedDir(fs, srcChild, dstChild, info.Mode().Perm())
	}
	return copyTreeNestedFile(fs, srcChild, dstChild, info)
}

// copyTreeNestedDir mkdir's dst with src's mode, then recurses into
// copyTreeContents. MkdirAll (not Mkdir) is correct here: the
// top-level reservation already ruled out concurrent winners for
// the outermost slot, so idempotence is fine.
func copyTreeNestedDir(fs driven.FileSystem, src, dst string, mode iofs.FileMode) error {
	if err := fs.MkdirAll(dst, mode); err != nil {
		return fmt.Errorf("mkdir %s: %w", dst, err)
	}
	return copyTreeContents(fs, src, dst)
}

// copyTreeNestedFile copies a single file from src to dst preserving
// the source mode and enforcing the maxBackupFileSize cap.
func copyTreeNestedFile(fs driven.FileSystem, src, dst string, info iofs.FileInfo) error {
	if info.Size() > maxBackupFileSize {
		return fmt.Errorf("%w: %s is %d bytes (cap %d)",
			driving.ErrBackupTooLarge, src, info.Size(), maxBackupFileSize)
	}
	data, err := fs.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read %s: %w", src, err)
	}
	if err := fs.WriteFile(dst, data, info.Mode().Perm()); err != nil {
		return fmt.Errorf("write %s: %w", dst, err)
	}
	return nil
}
