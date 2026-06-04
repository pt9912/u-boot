package recordingfs_test

import (
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/pt9912/u-boot/internal/adapter/driven/recordingfs"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// fakeFS is a recording-aware test double that counts mutation calls
// (to verify the recorder's passthrough switch routes correctly) and
// serves a small pre-state via testing/fstest.MapFS for reads.
type fakeFS struct {
	files          fstest.MapFS
	writeCalls     int
	mkdirCalls     int
	renameCalls    int
	removeCalls    int
	copyCalls      int
	failWritePath  string
	failWriteError error
}

func (f *fakeFS) Exists(path string) (bool, error) {
	_, err := f.files.Open(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func (f *fakeFS) ReadFile(path string) ([]byte, error) { return f.files.ReadFile(path) }

func (f *fakeFS) ReadDir(path string) ([]fs.DirEntry, error) { return f.files.ReadDir(path) }

func (f *fakeFS) Lstat(path string) (fs.FileInfo, error) {
	info, err := f.files.Stat(path)
	return info, err
}

func (f *fakeFS) WriteFile(path string, data []byte, _ fs.FileMode) error {
	if path == f.failWritePath {
		return f.failWriteError
	}
	f.writeCalls++
	f.files[path] = &fstest.MapFile{Data: data}
	return nil
}

func (f *fakeFS) WriteFileExclusive(path string, data []byte, _ fs.FileMode) error {
	f.writeCalls++
	f.files[path] = &fstest.MapFile{Data: data}
	return nil
}

func (f *fakeFS) Mkdir(path string, _ fs.FileMode) error {
	f.mkdirCalls++
	f.files[path] = &fstest.MapFile{Mode: fs.ModeDir}
	return nil
}

func (f *fakeFS) MkdirAll(path string, _ fs.FileMode) error {
	f.mkdirCalls++
	f.files[path] = &fstest.MapFile{Mode: fs.ModeDir}
	return nil
}

func (f *fakeFS) Rename(_, _ string) error { f.renameCalls++; return nil }

func (f *fakeFS) RemoveAll(_ string) error { f.removeCalls++; return nil }

func (f *fakeFS) Copy(_, _ string, _ fs.FileMode) error { f.copyCalls++; return nil }

func (f *fakeFS) CopyExclusive(_, _ string, _ fs.FileMode) error { f.copyCalls++; return nil }

func newFakeFS(seed map[string]string) *fakeFS {
	files := fstest.MapFS{}
	for path, body := range seed {
		files[path] = &fstest.MapFile{Data: []byte(body)}
	}
	return &fakeFS{files: files}
}

// TestRecordingFS_DryRun_WriteFile_DoesNotDelegate is the canonical
// negative-pin from slice-v1-cli-json-dry-run-add T0-(b): under
// Passthrough=false the underlying FS sees zero mutation calls.
func TestRecordingFS_DryRun_WriteFile_DoesNotDelegate(t *testing.T) {
	prod := newFakeFS(nil)
	rec := recordingfs.New(prod)
	if err := rec.WriteFile("compose.yaml", []byte("a\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if prod.writeCalls != 0 {
		t.Errorf("dry-run passthrough=false: production FS must see 0 writes, got %d", prod.writeCalls)
	}
	captured := rec.Captured()
	if len(captured) != 1 {
		t.Fatalf("expected 1 captured record, got %d", len(captured))
	}
	if captured[0].Path != "compose.yaml" || captured[0].Action != "create" {
		t.Errorf("captured record %+v: want path=compose.yaml action=create", captured[0])
	}
	if string(captured[0].NewContent) != "a\n" {
		t.Errorf("NewContent: want %q, got %q", "a\n", captured[0].NewContent)
	}
}

func TestRecordingFS_Passthrough_WriteFile_Delegates(t *testing.T) {
	prod := newFakeFS(nil)
	rec := recordingfs.New(prod, recordingfs.WithPassthrough(true))
	if err := rec.WriteFile("compose.yaml", []byte("body"), 0o644); err != nil {
		t.Fatal(err)
	}
	if prod.writeCalls != 1 {
		t.Errorf("passthrough=true: want 1 production write, got %d", prod.writeCalls)
	}
	if len(rec.Captured()) != 1 {
		t.Errorf("passthrough=true: capture must still happen, got %d records", len(rec.Captured()))
	}
}

// TestRecordingFS_WriteFile_PreStateClassification pins the
// create-vs-modify classification: if the file existed before the
// write, Action="modify" with the old bytes; otherwise "create".
func TestRecordingFS_WriteFile_PreStateClassification(t *testing.T) {
	t.Run("missing file → create", func(t *testing.T) {
		rec := recordingfs.New(newFakeFS(nil))
		_ = rec.WriteFile("new.yaml", []byte("new"), 0o644)
		got := rec.Captured()[0]
		if got.Action != "create" || got.OldContent != nil {
			t.Errorf("want action=create OldContent=nil, got %+v", got)
		}
	})
	t.Run("existing file → modify", func(t *testing.T) {
		prod := newFakeFS(map[string]string{"x.yaml": "old"})
		rec := recordingfs.New(prod)
		_ = rec.WriteFile("x.yaml", []byte("new"), 0o644)
		got := rec.Captured()[0]
		if got.Action != "modify" || string(got.OldContent) != "old" {
			t.Errorf("want action=modify OldContent=old, got action=%s OldContent=%q",
				got.Action, got.OldContent)
		}
	})
}

func TestRecordingFS_DryRun_MutationFailureNeverHappens(t *testing.T) {
	prod := newFakeFS(nil)
	prod.failWritePath = "compose.yaml"
	prod.failWriteError = errors.New("disk full")
	rec := recordingfs.New(prod) // passthrough=false
	if err := rec.WriteFile("compose.yaml", []byte("a"), 0o644); err != nil {
		t.Errorf("dry-run must never see production errors, got %v", err)
	}
}

func TestRecordingFS_Passthrough_MutationFailureSurfaced(t *testing.T) {
	prod := newFakeFS(nil)
	prod.failWritePath = "compose.yaml"
	prod.failWriteError = errors.New("disk full")
	rec := recordingfs.New(prod, recordingfs.WithPassthrough(true))
	err := rec.WriteFile("compose.yaml", []byte("a"), 0o644)
	if err == nil || err.Error() != "disk full" {
		t.Errorf("passthrough=true must surface production errors, got %v", err)
	}
	// Capture must still record the attempted call — slice T0-(b)
	// Mid-Write-Failure scenario relies on this.
	if len(rec.Captured()) != 1 {
		t.Errorf("capture must include the failed attempt, got %d records", len(rec.Captured()))
	}
}

// TestRecordingFS_ReadsAlwaysDelegate sweeps the 4 read methods to
// confirm they don't end up in the capture log.
func TestRecordingFS_ReadsAlwaysDelegate(t *testing.T) {
	prod := newFakeFS(map[string]string{"x.yaml": "body"})
	rec := recordingfs.New(prod)
	_, _ = rec.Exists("x.yaml")
	_, _ = rec.ReadFile("x.yaml")
	_, _ = rec.ReadDir(".")
	_, _ = rec.Lstat("x.yaml")
	if len(rec.Captured()) != 0 {
		t.Errorf("reads must not be captured, got %d records", len(rec.Captured()))
	}
}

// TestRecordingFS_AllEightMutationMethodsCaptured is the cluster-T0-(b)
// drift guard: even mutations add doesn't use today (RemoveAll/Rename/
// Mkdir/Copy/...) must be recorded so future use cases don't slip past
// the dry-run filter.
func TestRecordingFS_AllEightMutationMethodsCaptured(t *testing.T) {
	prod := newFakeFS(nil)
	rec := recordingfs.New(prod)

	_ = rec.WriteFile("w.txt", nil, 0o644)
	_ = rec.WriteFileExclusive("we.txt", nil, 0o644)
	_ = rec.Mkdir("d", 0o755)
	_ = rec.MkdirAll("d/d2", 0o755)
	_ = rec.Rename("a", "b")
	_ = rec.RemoveAll("r")
	_ = rec.Copy("src", "dst", 0o644)
	_ = rec.CopyExclusive("src", "dst2", 0o644)

	captured := rec.Captured()
	// Rename produces TWO records (delete src + create dst); everything
	// else produces one. Total: 9.
	if len(captured) != 9 {
		t.Errorf("expected 9 captured records (Rename produces 2), got %d", len(captured))
	}
	if prod.writeCalls+prod.mkdirCalls+prod.renameCalls+prod.removeCalls+prod.copyCalls != 0 {
		t.Errorf("dry-run mode must leave production FS untouched")
	}
}

// TestRecordingFS_NewContentDefensivelyCopied pins that the recorder
// stores a defensive copy of the WriteFile body so later caller-side
// buffer reuse cannot corrupt past records.
func TestRecordingFS_NewContentDefensivelyCopied(t *testing.T) {
	rec := recordingfs.New(newFakeFS(nil))
	body := []byte("original")
	_ = rec.WriteFile("x.txt", body, 0o644)
	body[0] = 'X' // mutate caller buffer
	captured := rec.Captured()
	if string(captured[0].NewContent) != "original" {
		t.Errorf("NewContent must be defensively copied; got %q", captured[0].NewContent)
	}
}

// TestRecordingFS_CapturedReturnsDefensiveCopy pins the same for the
// reverse direction — caller-side mutation of the returned slice must
// not affect future Captured() results.
func TestRecordingFS_CapturedReturnsDefensiveCopy(t *testing.T) {
	rec := recordingfs.New(newFakeFS(nil))
	_ = rec.WriteFile("x.txt", []byte("v"), 0o644)
	first := rec.Captured()
	first[0].Path = "tampered"
	second := rec.Captured()
	if second[0].Path != "x.txt" {
		t.Errorf("Captured must return defensive copy; got %q", second[0].Path)
	}
}

// Compile-time check (mirrors the production fs adapter): drift would
// surface as a build error here, not only a test-time error.
var _ driven.FileSystem = (*recordingfs.RecordingFileSystem)(nil)
var _ driven.RecorderPort = (*recordingfs.RecordingFileSystem)(nil)
