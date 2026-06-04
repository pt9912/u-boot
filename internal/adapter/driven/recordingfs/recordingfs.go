// Package recordingfs is the recording-wrapper implementation of the
// `port/driven.FileSystem` interface for u-boot's preview/dry-run
// path (slice-v1-cli-json-dry-run-add T1-B; Cluster T0-(b) Variante 2).
//
// The wrapper satisfies two driven ports at once:
//
//   - [driven.FileSystem]: every method either delegates to the
//     underlying production FS (reads always; mutations only when
//     Passthrough=true) or, for mutations under Passthrough=false,
//     records the call without touching disk.
//   - [driven.RecorderPort]: a single [RecordingFileSystem.Captured]
//     method returns the captured mutation log in call order.
//
// The CLI adapter never imports this package directly — depguard
// `adapter-driving-no-driven` forbids it. The Composition-Root
// (`cmd/uboot/main.go`) constructs an instance per
// [driving.AddPreviewMode] and hands it to the application service
// as both ports. See the slice's T0-(i) Outcome for the layer rules.
//
// Layer rule: adapters may import their driven-port interface plus
// external libs; they may not import application or other adapter
// packages (LH-FA-ARCH-003, depguard-enforced).
package recordingfs

import (
	"errors"
	"io/fs"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// Action constants mirror the LH-FA-CLI-007 §354 enum.
const (
	actionCreate = "create"
	actionModify = "modify"
	actionDelete = "delete"
)

// RecordingFileSystem wraps a production [driven.FileSystem] and
// records every mutation method-call. Reads always delegate through;
// mutations follow the Passthrough switch:
//
//   - Passthrough=false → record only; the underlying FS is untouched.
//     This is the LH-FA-CLI-007 dry-run path (slice T0-(b) Variante 2).
//   - Passthrough=true  → record AND delegate. This is the
//     LH-FA-CLI-008 preview-and-apply path (Spec §465-470): the
//     CLI shows the user what is about to be written and writes it.
//
// Pre-write capture (slice T0-(b) Mid-Failure semantics): every
// mutation method fetches the pre-state via the underlying FS's
// ReadFile/Lstat BEFORE applying the action, so [FileMutationRecord.
// OldContent] and the action classification (create vs. modify vs.
// delete) reflect reality at the moment of the call. Read errors on
// the pre-state are swallowed (treated as "did not exist" → "create")
// — the recorder is best-effort by design; the actual mutation's
// own error path (Passthrough=true) is the authoritative signal of
// success/failure.
type RecordingFileSystem struct {
	underlying  driven.FileSystem
	passthrough bool
	records     []driven.FileMutationRecord
}

// Static checks: RecordingFileSystem satisfies both ports.
var _ driven.FileSystem = (*RecordingFileSystem)(nil)
var _ driven.RecorderPort = (*RecordingFileSystem)(nil)

// Option mutates a [RecordingFileSystem] during [New]. Functional-
// options pattern matches the rest of the codebase (cli.New, etc.).
type Option func(*RecordingFileSystem)

// WithPassthrough toggles the passthrough switch. Default is false
// (record-only, dry-run mode).
func WithPassthrough(on bool) Option {
	return func(r *RecordingFileSystem) { r.passthrough = on }
}

// New wraps underlying in a recorder. underlying must be non-nil; the
// Composition-Root is responsible for that contract.
func New(underlying driven.FileSystem, opts ...Option) *RecordingFileSystem {
	r := &RecordingFileSystem{underlying: underlying}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Captured implements [driven.RecorderPort]. Returns a defensive copy
// so callers may mutate it freely.
func (r *RecordingFileSystem) Captured() []driven.FileMutationRecord {
	if len(r.records) == 0 {
		return nil
	}
	out := make([]driven.FileMutationRecord, len(r.records))
	copy(out, r.records)
	return out
}

// ----- Read methods: always delegate ------------------------------

// Exists delegates to the underlying FS.
func (r *RecordingFileSystem) Exists(path string) (bool, error) {
	return r.underlying.Exists(path)
}

// ReadFile delegates to the underlying FS.
func (r *RecordingFileSystem) ReadFile(path string) ([]byte, error) {
	return r.underlying.ReadFile(path)
}

// ReadDir delegates to the underlying FS.
func (r *RecordingFileSystem) ReadDir(path string) ([]fs.DirEntry, error) {
	return r.underlying.ReadDir(path)
}

// Lstat delegates to the underlying FS.
func (r *RecordingFileSystem) Lstat(path string) (fs.FileInfo, error) {
	return r.underlying.Lstat(path)
}

// ----- Mutation methods: record + optionally delegate -------------

// WriteFile records the call; delegates only when passthrough is on.
func (r *RecordingFileSystem) WriteFile(path string, data []byte, mode fs.FileMode) error {
	r.recordWrite(path, data)
	if r.passthrough {
		return r.underlying.WriteFile(path, data, mode)
	}
	return nil
}

// WriteFileExclusive records the call; delegates only when passthrough
// is on.
func (r *RecordingFileSystem) WriteFileExclusive(path string, data []byte, mode fs.FileMode) error {
	r.recordWrite(path, data)
	if r.passthrough {
		return r.underlying.WriteFileExclusive(path, data, mode)
	}
	return nil
}

// Mkdir records the call; delegates only when passthrough is on. The
// recorder treats Mkdir as a synthetic "create" with empty NewContent —
// the underlying file is a directory marker, not a regular file.
func (r *RecordingFileSystem) Mkdir(path string, mode fs.FileMode) error {
	r.recordDir(path)
	if r.passthrough {
		return r.underlying.Mkdir(path, mode)
	}
	return nil
}

// MkdirAll records the call; delegates only when passthrough is on.
// Mirrors [Mkdir]'s "create" classification.
func (r *RecordingFileSystem) MkdirAll(path string, mode fs.FileMode) error {
	r.recordDir(path)
	if r.passthrough {
		return r.underlying.MkdirAll(path, mode)
	}
	return nil
}

// Rename records the call as a "delete" on src and a "create" or
// "modify" on dst (depending on dst's pre-state); delegates only when
// passthrough is on.
func (r *RecordingFileSystem) Rename(src, dst string) error {
	srcContent := r.snapshot(src)
	r.records = append(r.records, driven.FileMutationRecord{
		Path: src, Action: actionDelete, OldContent: srcContent,
	})
	r.recordCopyOrMove(dst, srcContent)
	if r.passthrough {
		return r.underlying.Rename(src, dst)
	}
	return nil
}

// RemoveAll records the call; delegates only when passthrough is on.
func (r *RecordingFileSystem) RemoveAll(path string) error {
	r.records = append(r.records, driven.FileMutationRecord{
		Path: path, Action: actionDelete, OldContent: r.snapshot(path),
	})
	if r.passthrough {
		return r.underlying.RemoveAll(path)
	}
	return nil
}

// Copy records the call (reading src for NewContent before applying);
// delegates only when passthrough is on.
func (r *RecordingFileSystem) Copy(src, dst string, mode fs.FileMode) error {
	srcContent := r.snapshot(src)
	r.recordCopyOrMove(dst, srcContent)
	if r.passthrough {
		return r.underlying.Copy(src, dst, mode)
	}
	return nil
}

// CopyExclusive records the call (reading src for NewContent before
// applying); delegates only when passthrough is on.
func (r *RecordingFileSystem) CopyExclusive(src, dst string, mode fs.FileMode) error {
	srcContent := r.snapshot(src)
	r.recordCopyOrMove(dst, srcContent)
	if r.passthrough {
		return r.underlying.CopyExclusive(src, dst, mode)
	}
	return nil
}

// ----- Internal capture helpers -----------------------------------

// recordWrite is the shared body for WriteFile/WriteFileExclusive.
// Fetches the pre-state to classify create vs. modify.
func (r *RecordingFileSystem) recordWrite(path string, data []byte) {
	old := r.snapshot(path)
	action := actionCreate
	if old != nil {
		action = actionModify
	}
	r.records = append(r.records, driven.FileMutationRecord{
		Path:       path,
		Action:     action,
		NewContent: append([]byte(nil), data...),
		OldContent: old,
	})
}

// recordDir is the shared body for Mkdir/MkdirAll. Directories carry
// empty NewContent — they are markers, not regular files. Action is
// always "create" (MkdirAll on an existing dir is a no-op and
// recording it as such has no UX value; the LCS diff would render
// nothing).
func (r *RecordingFileSystem) recordDir(path string) {
	r.records = append(r.records, driven.FileMutationRecord{
		Path:   path,
		Action: actionCreate,
	})
}

// recordCopyOrMove classifies the dst side of a Copy/Rename: if dst
// existed before the call, action is "modify"; otherwise "create".
// NewContent is the src snapshot fetched by the caller.
func (r *RecordingFileSystem) recordCopyOrMove(dst string, newContent []byte) {
	old := r.snapshot(dst)
	action := actionCreate
	if old != nil {
		action = actionModify
	}
	r.records = append(r.records, driven.FileMutationRecord{
		Path:       dst,
		Action:     action,
		NewContent: newContent,
		OldContent: old,
	})
}

// snapshot reads path's current content via the underlying FS or
// returns nil when the file does not exist / cannot be read. The
// recorder is best-effort: read errors during pre-state capture are
// treated as "did not exist" so the action falls back to "create".
// The actual mutation's own error path (passthrough=true) is the
// authoritative success/failure signal — recorder reads are a UX
// aid, not a correctness gate.
func (r *RecordingFileSystem) snapshot(path string) []byte {
	data, err := r.underlying.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return nil
	}
	return data
}
