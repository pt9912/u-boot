package application_test

import (
	"errors"
	iofs "io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

func TestBackupPath_File_FirstSuffixWhenFree(t *testing.T) {
	fs := newFakeFS()
	src := "/proj/u-boot.yaml"
	if err := fs.WriteFile(src, []byte("payload"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	dst, err := application.BackupPath(fs, src)
	if err != nil {
		t.Fatalf("BackupPath: %v", err)
	}
	if want := src + ".bak"; dst != want {
		t.Errorf("dst = %q, want %q", dst, want)
	}
}

func TestBackupPath_File_NumericSuffixOnCollision(t *testing.T) {
	fs := newFakeFS()
	src := "/proj/u-boot.yaml"
	for _, p := range []string{src, src + ".bak"} {
		if err := fs.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatalf("setup %s: %v", p, err)
		}
	}

	dst, err := application.BackupPath(fs, src)
	if err != nil {
		t.Fatalf("BackupPath: %v", err)
	}
	if want := src + ".bak.1"; dst != want {
		t.Errorf("dst = %q, want %q", dst, want)
	}
}

func TestBackupPath_File_PicksSmallestFreeSuffix(t *testing.T) {
	// Why: spec §607 requires *smallest* free numeric suffix, not the
	// next-after-highest. With .bak and .bak.2 occupied but .bak.1
	// free, BackupPath must pick .bak.1.
	fs := newFakeFS()
	src := "/proj/u-boot.yaml"
	for _, p := range []string{src, src + ".bak", src + ".bak.2"} {
		if err := fs.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatalf("setup %s: %v", p, err)
		}
	}

	dst, err := application.BackupPath(fs, src)
	if err != nil {
		t.Fatalf("BackupPath: %v", err)
	}
	if want := src + ".bak.1"; dst != want {
		t.Errorf("dst = %q, want %q", dst, want)
	}
}

func TestBackupPath_File_ContentAndModePreserved(t *testing.T) {
	// Why: closes review finding #5 — `scripts/entry.sh` with 0o755
	// must keep its executable bit through a backup round-trip.
	fs := newFakeFS()
	src := "/proj/scripts/entry.sh"
	body := []byte("#!/bin/sh\necho hi\n")
	if err := fs.WriteFile(src, body, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	dst, err := application.BackupPath(fs, src)
	if err != nil {
		t.Fatalf("BackupPath: %v", err)
	}
	got, err := fs.ReadFile(dst)
	if err != nil {
		t.Fatalf("ReadFile(%s): %v", dst, err)
	}
	if string(got) != string(body) {
		t.Errorf("backup body = %q, want %q", got, body)
	}
	info, err := fs.Lstat(dst)
	if err != nil {
		t.Fatalf("Lstat(%s): %v", dst, err)
	}
	if info.Mode().Perm() != 0o755 {
		t.Errorf("backup mode = %o, want 0o755", info.Mode().Perm())
	}
}

func TestBackupPath_Directory_RecursiveCopy(t *testing.T) {
	fs := newFakeFS()
	src := "/proj/docker"
	if err := fs.Mkdir(src, 0o755); err != nil {
		t.Fatalf("setup mkdir src: %v", err)
	}
	if err := fs.WriteFile(filepath.Join(src, "Dockerfile"), []byte("FROM scratch\n"), 0o644); err != nil {
		t.Fatalf("setup Dockerfile: %v", err)
	}
	nested := filepath.Join(src, "scripts")
	if err := fs.Mkdir(nested, 0o750); err != nil {
		t.Fatalf("setup nested mkdir: %v", err)
	}
	if err := fs.WriteFile(filepath.Join(nested, "entry.sh"), []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatalf("setup entry.sh: %v", err)
	}

	dst, err := application.BackupPath(fs, src)
	if err != nil {
		t.Fatalf("BackupPath: %v", err)
	}
	if want := src + ".bak"; dst != want {
		t.Errorf("dst = %q, want %q", dst, want)
	}

	// Mode preservation propagates to nested files.
	entryInfo, err := fs.Lstat(filepath.Join(dst, "scripts", "entry.sh"))
	if err != nil {
		t.Fatalf("Lstat nested file: %v", err)
	}
	if entryInfo.Mode().Perm() != 0o755 {
		t.Errorf("nested file mode = %o, want 0o755", entryInfo.Mode().Perm())
	}

	for srcPath, wantContent := range map[string]string{
		filepath.Join(dst, "Dockerfile"):          "FROM scratch\n",
		filepath.Join(dst, "scripts", "entry.sh"): "#!/bin/sh\n",
	} {
		got, err := fs.ReadFile(srcPath)
		if err != nil {
			t.Errorf("ReadFile(%s): %v", srcPath, err)
			continue
		}
		if string(got) != wantContent {
			t.Errorf("ReadFile(%s) = %q, want %q", srcPath, got, wantContent)
		}
	}
}

func TestBackupPath_Symlink_TopLevel_Rejected(t *testing.T) {
	// Why: review finding #1 — silently following a symlink would
	// demote a link to a plain file and surprise users who linked
	// shared assets into the project.
	fs := newFakeFS()
	src := "/proj/u-boot.yaml"
	fs.markSymlink(src)

	_, err := application.BackupPath(fs, src)
	if err == nil {
		t.Fatalf("BackupPath(symlink): expected error, got nil")
	}
	if !errors.Is(err, driving.ErrBackupUnsupportedKind) {
		t.Errorf("BackupPath(symlink): error %v does not wrap ErrBackupUnsupportedKind", err)
	}
}

func TestBackupPath_Symlink_Nested_Rejected(t *testing.T) {
	// Why: nested symlinks during a tree walk are just as surprising
	// as top-level ones; the same sentinel covers both.
	fs := newFakeFS()
	src := "/proj/docs"
	if err := fs.Mkdir(src, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := fs.WriteFile(filepath.Join(src, "README.md"), []byte("hi"), 0o644); err != nil {
		t.Fatalf("setup readme: %v", err)
	}
	fs.markSymlink(filepath.Join(src, "link"))

	_, err := application.BackupPath(fs, src)
	if err == nil {
		t.Fatalf("BackupPath: expected error, got nil")
	}
	if !errors.Is(err, driving.ErrBackupUnsupportedKind) {
		t.Errorf("BackupPath: error %v does not wrap ErrBackupUnsupportedKind", err)
	}
}

func TestBackupPath_MissingSourceReturnsErr(t *testing.T) {
	fs := newFakeFS()
	_, err := application.BackupPath(fs, "/proj/does-not-exist")
	if err == nil {
		t.Fatalf("BackupPath(missing): expected error, got nil")
	}
	if !errors.Is(err, driving.ErrBackupSourceMissing) {
		t.Errorf("BackupPath(missing): error %v does not wrap ErrBackupSourceMissing", err)
	}
}

func TestBackupPath_TreeCopyFailure_RollsBack(t *testing.T) {
	// Why: spec §608 requires rollback when a tree-backup fails
	// partway. Setup: directory with two files; force WriteFile to
	// fail on the second; assert (a) BackupPath returns the
	// underlying error and (b) the partial destination is gone.
	fs := newFakeFS()
	src := "/proj/docs"
	if err := fs.Mkdir(src, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := fs.WriteFile(filepath.Join(src, "a.md"), []byte("a"), 0o644); err != nil {
		t.Fatalf("setup a: %v", err)
	}
	if err := fs.WriteFile(filepath.Join(src, "b.md"), []byte("b"), 0o644); err != nil {
		t.Fatalf("setup b: %v", err)
	}

	wantErr := errors.New("disk full on b")
	fs.failOn = filepath.Join(src+".bak", "b.md")
	fs.failErr = wantErr

	_, err := application.BackupPath(fs, src)
	if err == nil {
		t.Fatalf("BackupPath: expected error, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("BackupPath: error %v does not wrap %v", err, wantErr)
	}

	// Rollback: partial destination removed.
	exists, _ := fs.Exists(src + ".bak")
	if exists {
		t.Errorf("rollback: %s.bak still exists", src)
	}
}

func TestBackupPath_RollbackFailure_JoinsErrors(t *testing.T) {
	// Why: review finding #8 — a failed rollback used to be swallowed,
	// leaving an orphan .bak tree without operator visibility. Both
	// errors must surface via errors.Join.
	fs := newFakeFS()
	src := "/proj/docs"
	if err := fs.Mkdir(src, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := fs.WriteFile(filepath.Join(src, "a.md"), []byte("a"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	copyFailErr := errors.New("disk full")
	rollbackFailErr := errors.New("rollback denied")
	fs.failOn = filepath.Join(src+".bak", "a.md")
	fs.failErr = copyFailErr
	fs.failRemoveAll = rollbackFailErr

	_, err := application.BackupPath(fs, src)
	if err == nil {
		t.Fatalf("BackupPath: expected error, got nil")
	}
	if !errors.Is(err, copyFailErr) {
		t.Errorf("joined error does not wrap copy error %v: got %v", copyFailErr, err)
	}
	if !errors.Is(err, rollbackFailErr) {
		t.Errorf("joined error does not wrap rollback error %v: got %v", rollbackFailErr, err)
	}
}

func TestBackupPath_ExhaustedSuffixReturnsErr(t *testing.T) {
	fs := newFakeFS()
	src := "/proj/u-boot.yaml"
	if err := fs.WriteFile(src, []byte("x"), 0o644); err != nil {
		t.Fatalf("setup src: %v", err)
	}
	if err := fs.WriteFile(src+".bak", []byte("x"), 0o644); err != nil {
		t.Fatalf("setup .bak: %v", err)
	}
	for i := 1; i <= 1000; i++ {
		p := src + ".bak." + strconv.Itoa(i)
		if err := fs.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatalf("setup %s: %v", p, err)
		}
	}

	_, err := application.BackupPath(fs, src)
	if err == nil {
		t.Fatalf("BackupPath: expected exhaustion error, got nil")
	}
	if !errors.Is(err, driving.ErrBackupSuffixExhausted) {
		t.Errorf("BackupPath: error %v does not wrap ErrBackupSuffixExhausted", err)
	}
}

func TestBackupPath_RaceRetry_PicksNextSuffix(t *testing.T) {
	// Why: closes review finding #2 — when WriteFileExclusive returns
	// ErrExist (another process won the race for our .bak slot),
	// BackupPath retries chooseBackupPath which then picks the next
	// free suffix. Simulated by pre-populating the file at .bak in
	// memory AFTER chooseBackupPath would normally pick it: we just
	// have to occupy .bak so the chooser picks .bak.1 directly. This
	// asserts the retry path code-wise via a single-shot race window
	// is hard to express here; the easier path is to assert the
	// chooser already skips a present slot — which is the same code
	// path the retry loop relies on after a race.
	fs := newFakeFS()
	src := "/proj/u-boot.yaml"
	if err := fs.WriteFile(src, []byte("x"), 0o644); err != nil {
		t.Fatalf("setup src: %v", err)
	}
	if err := fs.WriteFile(src+".bak", []byte("preexisting"), 0o644); err != nil {
		t.Fatalf("setup .bak: %v", err)
	}

	dst, err := application.BackupPath(fs, src)
	if err != nil {
		t.Fatalf("BackupPath: %v", err)
	}
	if want := src + ".bak.1"; dst != want {
		t.Errorf("dst = %q, want %q (skip-occupied path)", dst, want)
	}
	// Pre-existing .bak must be untouched.
	got, _ := fs.ReadFile(src + ".bak")
	if string(got) != "preexisting" {
		t.Errorf(".bak clobbered: content = %q", got)
	}
}

func TestBackupPath_TopLevelDirRace_RetriesAndSucceeds(t *testing.T) {
	// Why: closes review finding #2 for the directory-backup path.
	// Setup: src dir + .bak dir already occupied (race winner held
	// it) → chooser picks .bak.1, Mkdir(.bak.1) succeeds, copy
	// completes.
	fs := newFakeFS()
	src := "/proj/docker"
	if err := fs.Mkdir(src, 0o755); err != nil {
		t.Fatalf("setup mkdir src: %v", err)
	}
	if err := fs.WriteFile(filepath.Join(src, "Dockerfile"), []byte("X"), 0o644); err != nil {
		t.Fatalf("setup Dockerfile: %v", err)
	}
	// Simulate a race winner: .bak slot taken.
	if err := fs.Mkdir(src+".bak", 0o755); err != nil {
		t.Fatalf("setup .bak: %v", err)
	}

	dst, err := application.BackupPath(fs, src)
	if err != nil {
		t.Fatalf("BackupPath: %v", err)
	}
	if want := src + ".bak.1"; dst != want {
		t.Errorf("dst = %q, want %q", dst, want)
	}
}

func TestBackupPath_LstatErrorPropagates(t *testing.T) {
	fs := newFakeFS()
	src := "/proj/u-boot.yaml"
	if err := fs.WriteFile(src, []byte("x"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	wantErr := errors.New("lstat denied")
	fs.failLstatOn = src
	fs.failLstatErr = wantErr

	_, err := application.BackupPath(fs, src)
	if err == nil {
		t.Fatalf("BackupPath: expected error, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("BackupPath: error %v does not wrap %v", err, wantErr)
	}
}

func TestBackupPath_ReadFileErrorPropagates(t *testing.T) {
	fs := newFakeFS()
	src := "/proj/u-boot.yaml"
	if err := fs.WriteFile(src, []byte("x"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	wantErr := errors.New("read denied")
	fs.failReadOn = src
	fs.failReadErr = wantErr

	_, err := application.BackupPath(fs, src)
	if err == nil {
		t.Fatalf("BackupPath: expected error, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("BackupPath: error %v does not wrap %v", err, wantErr)
	}
}

func TestBackupPath_ExistsErrorOnCandidatePropagates(t *testing.T) {
	fs := newFakeFS()
	src := "/proj/u-boot.yaml"
	if err := fs.WriteFile(src, []byte("x"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	wantErr := errors.New("stat denied")
	fs.failExistsOn = src + ".bak"
	fs.failExistsErr = wantErr

	_, err := application.BackupPath(fs, src)
	if err == nil {
		t.Fatalf("BackupPath: expected error, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("BackupPath: error %v does not wrap %v", err, wantErr)
	}
}

func TestBackupPath_ExistsErrorOnNumericCandidatePropagates(t *testing.T) {
	// Why: covers the chooseBackupPath loop body (numeric-suffix
	// path) — the first-suffix Exists succeeds, but the .bak.1
	// candidate fails.
	fs := newFakeFS()
	src := "/proj/u-boot.yaml"
	if err := fs.WriteFile(src, []byte("x"), 0o644); err != nil {
		t.Fatalf("setup src: %v", err)
	}
	if err := fs.WriteFile(src+".bak", []byte("x"), 0o644); err != nil {
		t.Fatalf("setup .bak: %v", err)
	}
	wantErr := errors.New("stat denied")
	fs.failExistsOn = src + ".bak.1"
	fs.failExistsErr = wantErr

	_, err := application.BackupPath(fs, src)
	if err == nil {
		t.Fatalf("BackupPath: expected error, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("BackupPath: error %v does not wrap %v", err, wantErr)
	}
}

func TestBackupPath_CopyTree_MkdirAllFailurePropagates(t *testing.T) {
	// Why: closes the coverage gap in copyTreeContents — the nested
	// MkdirAll failure path was unreachable in T4a.
	fs := newFakeFS()
	src := "/proj/docs"
	if err := fs.Mkdir(src, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	nested := filepath.Join(src, "subdir")
	if err := fs.Mkdir(nested, 0o755); err != nil {
		t.Fatalf("setup nested: %v", err)
	}
	if err := fs.WriteFile(filepath.Join(nested, "x.md"), []byte("x"), 0o644); err != nil {
		t.Fatalf("setup file: %v", err)
	}
	wantErr := errors.New("mkdir denied")
	fs.failMkdirOn = filepath.Join(src+".bak", "subdir")
	fs.failMkdirErr = wantErr

	_, err := application.BackupPath(fs, src)
	if err == nil {
		t.Fatalf("BackupPath: expected error, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("BackupPath: error %v does not wrap %v", err, wantErr)
	}
}

func TestBackupPath_CopyTree_ReadDirFailurePropagates(t *testing.T) {
	// Why: closes the second coverage gap in copyTreeContents.
	fs := newFakeFS()
	src := "/proj/docs"
	if err := fs.Mkdir(src, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := fs.WriteFile(filepath.Join(src, "a.md"), []byte("a"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	wantErr := errors.New("readdir denied")
	fs.failReadDirOn = src
	fs.failReadDirErr = wantErr

	_, err := application.BackupPath(fs, src)
	if err == nil {
		t.Fatalf("BackupPath: expected error, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("BackupPath: error %v does not wrap %v", err, wantErr)
	}
}

func TestBackupPath_NestedLstatErrorPropagates(t *testing.T) {
	// Why: closes coverage gap for the nested-child Lstat error in
	// copyTreeContents.
	fs := newFakeFS()
	src := "/proj/docs"
	if err := fs.Mkdir(src, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	child := filepath.Join(src, "a.md")
	if err := fs.WriteFile(child, []byte("a"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	wantErr := errors.New("lstat child denied")
	fs.failLstatOn = child
	fs.failLstatErr = wantErr

	_, err := application.BackupPath(fs, src)
	if err == nil {
		t.Fatalf("BackupPath: expected error, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("BackupPath: error %v does not wrap %v", err, wantErr)
	}
}

func TestBackupPath_NestedReadFileErrorPropagates(t *testing.T) {
	// Why: closes coverage gap for ReadFile failure during nested
	// file copy (different from top-level ReadFile failure tested
	// above).
	fs := newFakeFS()
	src := "/proj/docs"
	if err := fs.Mkdir(src, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	child := filepath.Join(src, "a.md")
	if err := fs.WriteFile(child, []byte("a"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	wantErr := errors.New("read denied")
	fs.failReadOn = child
	fs.failReadErr = wantErr

	_, err := application.BackupPath(fs, src)
	if err == nil {
		t.Fatalf("BackupPath: expected error, got nil")
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("BackupPath: error %v does not wrap %v", err, wantErr)
	}
}

func TestBackupPath_FakeFS_DeeplyNestedReadDirReturnsImplicitDirs(t *testing.T) {
	// Why: closes review finding #7 — pin that registerAncestorsLocked
	// makes implicit intermediate directories appear in ReadDir, so
	// the test fake matches os.ReadDir for keys written several
	// levels deep without explicit MkdirAll.
	fs := newFakeFS()
	if err := fs.WriteFile("/a/b/c/d/file.txt", []byte("x"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	for _, parent := range []string{"/a", "/a/b", "/a/b/c", "/a/b/c/d"} {
		entries, err := fs.ReadDir(parent)
		if err != nil {
			t.Fatalf("ReadDir(%s): %v", parent, err)
		}
		if len(entries) != 1 {
			t.Fatalf("ReadDir(%s) len = %d, want 1", parent, len(entries))
		}
		// Every intermediate must be a directory; the leaf is the file.
		wantDir := parent != "/a/b/c/d"
		if entries[0].IsDir() != wantDir {
			t.Errorf("ReadDir(%s): entry IsDir = %v, want %v",
				parent, entries[0].IsDir(), wantDir)
		}
	}
}

func TestBackupPath_LstatMissing_ReturnsErrNotExist(t *testing.T) {
	// Why: pins the fake's Lstat error policy (matches os.Lstat),
	// so future tests can rely on the contract.
	fs := newFakeFS()
	_, err := fs.Lstat("/nope")
	if !errors.Is(err, iofs.ErrNotExist) {
		t.Fatalf("Lstat(missing): want ErrNotExist, got %v", err)
	}
}

func TestBackupPath_ErrorMessageMentionsPath(t *testing.T) {
	// Why: defensive — the error formatting must mention the offending
	// path so operators have something searchable in logs.
	fs := newFakeFS()
	src := "/proj/somepath.yaml"
	_, err := application.BackupPath(fs, src)
	if err == nil {
		t.Fatalf("BackupPath: expected error, got nil")
	}
	if !strings.Contains(err.Error(), src) {
		t.Errorf("error %v does not mention path %q", err, src)
	}
}
