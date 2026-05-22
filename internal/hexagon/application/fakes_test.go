package application_test

import (
	"context"
	"errors"
	iofs "io/fs"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// fakeFS is an in-memory FileSystem implementation for application-
// layer tests. It records every WriteFile and MkdirAll call so tests
// can assert on creation order; Exists/Read/Rename/ReadDir behave
// consistently with WriteFile/MkdirAll history.
//
// The fake is intentionally tiny — application tests need
// deterministic in-memory behaviour, not a full ioutil emulator.
type fakeFS struct {
	mu            sync.Mutex
	files         map[string][]byte
	dirs          map[string]bool
	writes        []string // ordered: every successful WriteFile path
	mkdirs        []string // ordered: every MkdirAll path
	failOn        string   // when non-empty, WriteFile returns failErr for that path
	failErr       error
	failReadOn    string // when non-empty, ReadFile returns failReadErr for that path
	failReadErr   error
	failExistsOn  string // when non-empty, Exists returns failExistsErr for that path
	failExistsErr error
	failIsDirOn   string // when non-empty, IsDir returns failIsDirErr for that path
	failIsDirErr  error
}

func newFakeFS() *fakeFS {
	return &fakeFS{
		files: make(map[string][]byte),
		dirs:  make(map[string]bool),
	}
}

func (f *fakeFS) Exists(path string) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failExistsOn != "" && path == f.failExistsOn {
		return false, f.failExistsErr
	}
	_, fileOK := f.files[path]
	_, dirOK := f.dirs[path]
	return fileOK || dirOK, nil
}

func (f *fakeFS) ReadFile(path string) ([]byte, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failReadOn != "" && path == f.failReadOn {
		return nil, f.failReadErr
	}
	data, ok := f.files[path]
	if !ok {
		return nil, iofs.ErrNotExist
	}
	out := make([]byte, len(data))
	copy(out, data)
	return out, nil
}

func (f *fakeFS) WriteFile(path string, data []byte, _ iofs.FileMode) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failOn != "" && path == f.failOn {
		return f.failErr
	}
	stored := make([]byte, len(data))
	copy(stored, data)
	f.files[path] = stored
	f.dirs[filepath.Dir(path)] = true
	f.writes = append(f.writes, path)
	return nil
}

func (f *fakeFS) MkdirAll(path string, _ iofs.FileMode) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.dirs[path] = true
	f.mkdirs = append(f.mkdirs, path)
	return nil
}

func (f *fakeFS) Rename(src, dst string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	data, ok := f.files[src]
	if !ok {
		return iofs.ErrNotExist
	}
	f.files[dst] = data
	delete(f.files, src)
	return nil
}

// ReadDir returns the direct children of path. The fake reconstructs
// them from the recorded files/dirs maps so the BackupPath tests can
// walk a directory tree without touching disk. A child is classified
// as a directory when (a) it is itself a recorded dir, or (b) it is
// not a recorded file but appears as a parent of some deeper key
// (implicit intermediate directory).
func (f *fakeFS) ReadDir(path string) ([]iofs.DirEntry, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if !f.dirs[path] {
		return nil, iofs.ErrNotExist
	}
	prefix := path + string(filepath.Separator)
	seen := make(map[string]bool)
	var entries []iofs.DirEntry
	addChild := func(key string) {
		if !strings.HasPrefix(key, prefix) {
			return
		}
		name := strings.SplitN(strings.TrimPrefix(key, prefix), string(filepath.Separator), 2)[0]
		if name == "" || seen[name] {
			return
		}
		seen[name] = true
		direct := filepath.Join(path, name)
		switch {
		case f.dirs[direct]:
			entries = append(entries, fakeDirEntry{name: name, isDir: true})
		case fileExistsLocked(f.files, direct):
			entries = append(entries, fakeDirEntry{name: name, isDir: false})
		default:
			// Implicit intermediate directory — no explicit entry, but
			// deeper keys live under it.
			entries = append(entries, fakeDirEntry{name: name, isDir: true})
		}
	}
	for k := range f.files {
		addChild(k)
	}
	for k := range f.dirs {
		addChild(k)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	return entries, nil
}

// fileExistsLocked is a tiny helper to keep the ReadDir switch
// readable; the lookup is mutex-protected by the caller.
func fileExistsLocked(files map[string][]byte, path string) bool {
	_, ok := files[path]
	return ok
}

// IsDir reports whether path is a recorded directory. Missing paths
// return `(false, nil)` to match the real adapter's policy. Tests
// inject a forced failure via failIsDirOn / failIsDirErr.
func (f *fakeFS) IsDir(path string) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failIsDirOn != "" && path == f.failIsDirOn {
		return false, f.failIsDirErr
	}
	return f.dirs[path], nil
}

// RemoveAll deletes path and any recorded children. Idempotent —
// removing a missing path is a no-op, matching os.RemoveAll.
func (f *fakeFS) RemoveAll(path string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.files, path)
	delete(f.dirs, path)
	prefix := path + string(filepath.Separator)
	for k := range f.files {
		if strings.HasPrefix(k, prefix) {
			delete(f.files, k)
		}
	}
	for k := range f.dirs {
		if strings.HasPrefix(k, prefix) {
			delete(f.dirs, k)
		}
	}
	return nil
}

// fakeDirEntry is the minimal iofs.DirEntry the BackupPath walker
// needs: Name() and IsDir() only. Type()/Info() return the zero values
// because the algorithm never consults them.
type fakeDirEntry struct {
	name  string
	isDir bool
}

func (e fakeDirEntry) Name() string               { return e.name }
func (e fakeDirEntry) IsDir() bool                { return e.isDir }
func (e fakeDirEntry) Type() iofs.FileMode        { return 0 }
func (e fakeDirEntry) Info() (iofs.FileInfo, error) { return nil, errors.New("fakeDirEntry.Info: not implemented") }

// markDirExists pre-registers a directory so Exists returns true.
// Used by test setup to satisfy the BaseDir-existence check
// without going through a real MkdirAll call (which the test
// otherwise wants to count).
func (f *fakeFS) markDirExists(path string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.dirs[path] = true
}

// writtenPaths returns the recorded WriteFile paths in deterministic
// order.
func (f *fakeFS) writtenPaths() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]string, len(f.writes))
	copy(out, f.writes)
	return out
}

// mkdirPaths returns the recorded MkdirAll paths in the order the
// service called them. Tests assert on the real call order so a
// reorder in writeDirectories cannot pass silently.
func (f *fakeFS) mkdirPaths() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]string, len(f.mkdirs))
	copy(out, f.mkdirs)
	return out
}

// fakeYAML uses gopkg.in/yaml.v3 directly. The application layer
// imports the YAMLCodec port, the test imports the library — the
// `*_test.go` Carveout (LH-FA-ARCH-003) covers that.
type fakeYAML struct {
	failMarshal bool
}

func (f *fakeYAML) Marshal(v any) ([]byte, error) {
	if f.failMarshal {
		return nil, errors.New("fakeYAML.Marshal: forced failure")
	}
	return yaml.Marshal(v)
}

func (f *fakeYAML) Unmarshal(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}

// fakeGit records IsRepository / Init calls and lets each test
// configure the IsRepository return values and Init error.
type fakeGit struct {
	isRepoCalls []string
	initCalls   []string
	isRepo      bool
	isRepoErr   error
	initErr     error
}

func (f *fakeGit) IsRepository(_ context.Context, dir string) (bool, error) {
	f.isRepoCalls = append(f.isRepoCalls, dir)
	return f.isRepo, f.isRepoErr
}

func (f *fakeGit) Init(_ context.Context, dir string) error {
	f.initCalls = append(f.initCalls, dir)
	return f.initErr
}
