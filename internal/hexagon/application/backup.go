package application

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// backupSuffixCap caps the .bak.N suffix search at 1000. A user who
// has accumulated this many stale backups under one path must clean
// them up manually; the cap exists so a bug in the chooser cannot
// loop forever.
const backupSuffixCap = 1000

// ErrBackupSourceMissing is returned by [BackupPath] when src does
// not exist. It signals an invariant violation in the caller (the
// service must check existence before invoking BackupPath); the
// sentinel keeps the failure recognizable in tests.
var ErrBackupSourceMissing = errors.New("backup source does not exist")

// ErrBackupSuffixExhausted is returned when <src>.bak through
// <src>.bak.{backupSuffixCap} are all occupied so no fresh backup
// path can be chosen. The CLI adapter will map this to a technical
// filesystem error (exit code 14) when M3-T4b wires it through.
var ErrBackupSuffixExhausted = errors.New("backup suffix exhausted")

// BackupPath copies src to a sibling backup path and returns the
// chosen backup path. Suffix selection follows LH-FA-INIT-005:
// <src>.bak first, then <src>.bak.1, .bak.2, ... — the smallest free
// numeric suffix is picked, so existing backups are never overwritten.
// Works for files and directory trees; on partial directory-copy
// failure, BackupPath rolls back by removing the partial destination
// (POSIX atomicity for recursive trees is not guaranteed — see
// LH-FA-INIT-005 §608).
//
// Mode policy: files are written with 0o644 and directories with
// 0o755 (the canonical u-boot modes). Backups are read by u-boot
// itself, so preserving arbitrary original modes would add port
// surface (Stat returning FileInfo) for no observable benefit.
func BackupPath(fs driven.FileSystem, src string) (string, error) {
	exists, err := fs.Exists(src)
	if err != nil {
		return "", fmt.Errorf("check %s: %w", src, err)
	}
	if !exists {
		return "", fmt.Errorf("%w: %s", ErrBackupSourceMissing, src)
	}

	dst, err := chooseBackupPath(fs, src)
	if err != nil {
		return "", err
	}

	isDir, err := fs.IsDir(src)
	if err != nil {
		return "", fmt.Errorf("stat %s: %w", src, err)
	}
	if isDir {
		if err := copyTree(fs, src, dst); err != nil {
			// Rollback the partial tree; ignore the rollback error
			// because the original copy error is the user-relevant one.
			_ = fs.RemoveAll(dst)
			return "", err
		}
		return dst, nil
	}

	if err := copyFile(fs, src, dst); err != nil {
		return "", err
	}
	return dst, nil
}

// chooseBackupPath returns the smallest-suffix backup destination
// that does not yet exist: <src>.bak, then .bak.1, .bak.2, ...
// Iteration is bounded by backupSuffixCap so a runaway loop cannot
// happen even if Exists keeps returning true.
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
		ErrBackupSuffixExhausted, src, src, backupSuffixCap)
}

// copyFile reads src and writes it to dst with the canonical u-boot
// file mode 0o644.
func copyFile(fs driven.FileSystem, src, dst string) error {
	data, err := fs.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read %s: %w", src, err)
	}
	if err := fs.WriteFile(dst, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", dst, err)
	}
	return nil
}

// copyTree recursively copies the directory at src into dst. The
// destination is created if missing. Children are walked in the
// order ReadDir returns them; the caller (BackupPath) handles
// rollback on error.
func copyTree(fs driven.FileSystem, src, dst string) error {
	if err := fs.MkdirAll(dst, 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", dst, err)
	}
	entries, err := fs.ReadDir(src)
	if err != nil {
		return fmt.Errorf("read dir %s: %w", src, err)
	}
	for _, entry := range entries {
		srcChild := filepath.Join(src, entry.Name())
		dstChild := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			if err := copyTree(fs, srcChild, dstChild); err != nil {
				return err
			}
			continue
		}
		if err := copyFile(fs, srcChild, dstChild); err != nil {
			return err
		}
	}
	return nil
}
