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

// TestRecordingFS_WriteFile_SynthesisesParentMkdirAll pins the
// slice-v1-cli-json-dry-run-add T0-(b) requirement that the recorder
// mirrors the production-FS implicit MkdirAll on the parent dir.
// Without this, --dry-run would silently drop the directory-creation
// hint that a sub-path write (e.g. `otel/collector-config.yaml`)
// implies in production.
func TestRecordingFS_WriteFile_SynthesisesParentMkdirAll(t *testing.T) {
	prod := newFakeFS(nil)
	rec := recordingfs.New(prod)
	if err := rec.WriteFile("otel/collector.yaml", []byte("body\n"), 0o644); err != nil {
		t.Fatalf("WriteFile sub-path: %v", err)
	}
	captured := rec.Captured()
	if len(captured) != 2 {
		t.Fatalf("expected 2 records (mkdir + write), got %d: %+v", len(captured), captured)
	}
	if captured[0].Path != "otel" || captured[0].Action != "create" {
		t.Errorf("first record: want {otel, create}, got {%s, %s}", captured[0].Path, captured[0].Action)
	}
	if captured[1].Path != "otel/collector.yaml" || captured[1].Action != "create" {
		t.Errorf("second record: want {otel/collector.yaml, create}, got {%s, %s}", captured[1].Path, captured[1].Action)
	}
}

// TestRecordingFS_WriteFile_NoMkdirForExistingDir pins the idempotent
// no-op branch: when the parent dir already exists, no synthetic
// MkdirAll record is emitted — matches production MkdirAll's idempotent
// behaviour and keeps the plannedFiles[] view free of noise.
func TestRecordingFS_WriteFile_NoMkdirForExistingDir(t *testing.T) {
	prod := newFakeFS(map[string]string{"otel/.keep": ""})
	rec := recordingfs.New(prod)
	if err := rec.WriteFile("otel/collector.yaml", []byte("body\n"), 0o644); err != nil {
		t.Fatalf("WriteFile sub-path: %v", err)
	}
	captured := rec.Captured()
	if len(captured) != 1 {
		t.Fatalf("expected 1 record (no mkdir for existing dir), got %d: %+v", len(captured), captured)
	}
	if captured[0].Path != "otel/collector.yaml" {
		t.Errorf("only record: want otel/collector.yaml, got %s", captured[0].Path)
	}
}

// TestRecordingFS_WriteFile_NoMkdirForFlatPath pins the flat-path
// case: writing `compose.yaml` (filepath.Dir = ".") must NOT emit a
// synthetic mkdir — the dir-anchor is the CWD which always exists.
func TestRecordingFS_WriteFile_NoMkdirForFlatPath(t *testing.T) {
	prod := newFakeFS(nil)
	rec := recordingfs.New(prod)
	if err := rec.WriteFile("compose.yaml", []byte("body\n"), 0o644); err != nil {
		t.Fatalf("WriteFile flat: %v", err)
	}
	captured := rec.Captured()
	if len(captured) != 1 {
		t.Fatalf("expected 1 record (no mkdir for flat path), got %d: %+v", len(captured), captured)
	}
}

// TestRecordingFS_ImplicitMkdirDeduplicates pins the slice-v1-cli-
// json-dry-run-init T0-(m) Dedup-Pflicht: an explicit MkdirAll
// followed by a WriteFile to a sub-path must NOT emit a duplicate
// synthetic Mkdir-record. In dry-run mode (passthrough=false) the
// underlying FS never sees the Mkdir, so without the knownDirs-
// dedup the implicit Mkdir would fire again on every WriteFile.
//
// Concrete scenario: `init --devcontainer --dry-run` writes
// `.devcontainer/devcontainer.json` after MkdirAll('.devcontainer').
// Without the fix the envelope would list `.devcontainer/` TWICE.
func TestRecordingFS_ImplicitMkdirDeduplicates(t *testing.T) {
	prod := newFakeFS(nil) // empty — .devcontainer does not pre-exist
	rec := recordingfs.New(prod)

	if err := rec.MkdirAll(".devcontainer", 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := rec.WriteFile(".devcontainer/devcontainer.json", []byte("{}"), 0o644); err != nil {
		t.Fatalf("WriteFile sub-path: %v", err)
	}
	if err := rec.WriteFile(".devcontainer/Dockerfile", []byte("FROM x"), 0o644); err != nil {
		t.Fatalf("WriteFile second sub-path: %v", err)
	}

	captured := rec.Captured()
	// Expected: 3 records exactly — Mkdir('.devcontainer'),
	// Write('.devcontainer/devcontainer.json'),
	// Write('.devcontainer/Dockerfile'). NO duplicate Mkdir-records.
	if len(captured) != 3 {
		t.Fatalf("expected 3 records (1 mkdir + 2 writes), got %d: %+v", len(captured), pathsAndActions(captured))
	}
	mkdirCount := 0
	for _, r := range captured {
		if r.Path == ".devcontainer" && r.Action == "create" {
			mkdirCount++
		}
	}
	if mkdirCount != 1 {
		t.Errorf("expected exactly 1 .devcontainer/ create-record (dedup), got %d", mkdirCount)
	}
}

// TestRecordingFS_ImplicitMkdir_FromMultipleSubWrites pins that two
// WriteFiles into the same not-yet-existing dir emit ONE synthetic
// Mkdir, not two — the dedup-map also fills from the first implicit
// Mkdir-record.
func TestRecordingFS_ImplicitMkdir_FromMultipleSubWrites(t *testing.T) {
	prod := newFakeFS(nil)
	rec := recordingfs.New(prod)

	_ = rec.WriteFile("otel/collector.yaml", []byte("a"), 0o644)
	_ = rec.WriteFile("otel/config.toml", []byte("b"), 0o644)

	captured := rec.Captured()
	if len(captured) != 3 {
		t.Fatalf("expected 3 records (1 synthetic mkdir + 2 writes), got %d: %+v", len(captured), pathsAndActions(captured))
	}
	if captured[0].Path != "otel" || captured[0].Action != "create" {
		t.Errorf("first record should be the synthetic mkdir for otel, got %+v", captured[0])
	}
}

// TestRecordingFS_RemoveAllClearsKnownDirs pins that RemoveAll on a
// previously-recorded dir clears it from knownDirs, so a subsequent
// WriteFile to the same sub-tree emits a FRESH synthetic Mkdir.
// Mirrors the actual filesystem state — RemoveAll undoes the dir.
func TestRecordingFS_RemoveAllClearsKnownDirs(t *testing.T) {
	prod := newFakeFS(nil)
	rec := recordingfs.New(prod)

	_ = rec.MkdirAll(".cache", 0o755)
	_ = rec.RemoveAll(".cache")
	_ = rec.WriteFile(".cache/x", []byte("y"), 0o644)

	captured := rec.Captured()
	// Expected: Mkdir('.cache'), Delete('.cache'), Mkdir('.cache') (re-create), Write('.cache/x') = 4 records.
	if len(captured) != 4 {
		t.Fatalf("expected 4 records (mkdir + delete + re-mkdir + write), got %d: %+v", len(captured), pathsAndActions(captured))
	}
	mkdirCount := 0
	for _, r := range captured {
		if r.Path == ".cache" && r.Action == "create" {
			mkdirCount++
		}
	}
	if mkdirCount != 2 {
		t.Errorf("expected 2 .cache/ create-records (one explicit + one synthetic after RemoveAll), got %d", mkdirCount)
	}
}

func pathsAndActions(records []driven.FileMutationRecord) []string {
	out := make([]string, len(records))
	for i, r := range records {
		out[i] = r.Path + "(" + r.Action + ")"
	}
	return out
}

// Compile-time check (mirrors the production fs adapter): drift would
// surface as a build error here, not only a test-time error.
var _ driven.FileSystem = (*recordingfs.RecordingFileSystem)(nil)
var _ driven.RecorderPort = (*recordingfs.RecordingFileSystem)(nil)
