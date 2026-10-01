package application

import (
	"errors"
	"fmt"
	iofs "io/fs"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// generateJournal records the pre-state of every path `generate
// devcontainer` is about to mutate, so a failure in the multi-file write
// phase can restore the disk to its state before the call
// (slice-v2-generate-devcontainer-rollback-aware-write, Option 1:
// in-memory snapshot + rollback-on-failure). A nil journal is a no-op
// (dry-run: nothing real was written).
//
// Rollback is best effort and LIFO: files that existed are rewritten
// with their original bytes and mode, files created by the call are
// removed, a `.devcontainer/` directory created by the call is removed.
// A failing restore is reported, not hidden.
type generateJournal struct {
	fs    driven.FileSystem
	files []fileSnapshot
	dirs  []string
	seen  map[string]bool
}

type fileSnapshot struct {
	path    string
	existed bool
	content []byte
	mode    iofs.FileMode
}

// newGenerateJournal returns a journal that restores through fs (the
// production filesystem, never the preview recorder), or nil when fs is
// nil.
func newGenerateJournal(fs driven.FileSystem) *generateJournal {
	if fs == nil {
		return nil
	}
	return &generateJournal{fs: fs, seen: map[string]bool{}}
}

// captureDir records that dir does not exist yet (so rollback removes
// it); an existing directory is left alone.
func (j *generateJournal) captureDir(dir string) error {
	if j == nil || j.seen["dir:"+dir] {
		return nil
	}
	j.seen["dir:"+dir] = true
	exists, err := j.fs.Exists(dir)
	if err != nil {
		return fmt.Errorf("rollback snapshot: Exists(%q): %w", dir, err)
	}
	if !exists {
		j.dirs = append(j.dirs, dir)
	}
	return nil
}

// captureFile snapshots path (content and mode when it exists) before
// its first mutation; repeated calls are no-ops.
func (j *generateJournal) captureFile(path string) error {
	if j == nil || j.seen[path] {
		return nil
	}
	j.seen[path] = true
	exists, err := j.fs.Exists(path)
	if err != nil {
		return fmt.Errorf("rollback snapshot: Exists(%q): %w", path, err)
	}
	snap := fileSnapshot{path: path, existed: exists}
	if exists {
		if snap.content, err = j.fs.ReadFile(path); err != nil {
			return fmt.Errorf("rollback snapshot: ReadFile(%q): %w", path, err)
		}
		info, err := j.fs.Lstat(path)
		if err != nil {
			return fmt.Errorf("rollback snapshot: Lstat(%q): %w", path, err)
		}
		snap.mode = info.Mode().Perm()
	}
	j.files = append(j.files, snap)
	return nil
}

// rollback restores the recorded pre-state (LIFO) and returns every
// restore failure joined, or nil when the disk is back to its pre-state.
func (j *generateJournal) rollback() error {
	if j == nil {
		return nil
	}
	var errs []error
	for i := len(j.files) - 1; i >= 0; i-- {
		f := j.files[i]
		var err error
		if f.existed {
			err = j.fs.WriteFile(f.path, f.content, f.mode)
		} else {
			err = j.fs.RemoveAll(f.path)
		}
		if err != nil {
			errs = append(errs, fmt.Errorf("restore %q: %w", f.path, err))
		}
	}
	for i := len(j.dirs) - 1; i >= 0; i-- {
		if err := j.fs.RemoveAll(j.dirs[i]); err != nil {
			errs = append(errs, fmt.Errorf("remove %q: %w", j.dirs[i], err))
		}
	}
	return errors.Join(errs...)
}

// rollbackOnError runs the rollback after a failed write phase and
// annotates err: a clean rollback keeps err as is (the disk is back to
// its pre-state); a failing rollback is appended so the user knows the
// state is inconsistent. The original error chain stays intact.
func (s *GenerateService) rollbackOnError(j *generateJournal, err error) error {
	if j == nil || err == nil {
		return err
	}
	if rbErr := j.rollback(); rbErr != nil {
		return fmt.Errorf("%w; rollback incomplete, manual cleanup of .devcontainer/ and u-boot.yaml may be needed: %v", err, rbErr)
	}
	s.logger.Info("generate devcontainer: write failed, rolled back to the previous state")
	return err
}
