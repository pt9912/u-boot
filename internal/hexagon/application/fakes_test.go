package application_test

import (
	"context"
	"errors"
	iofs "io/fs"
	"path/filepath"
	"sort"
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
	mu      sync.Mutex
	files   map[string][]byte
	dirs    map[string]bool
	writes  []string // ordered: every successful WriteFile path
	mkdirs  []string // ordered: every MkdirAll path
	failOn  string   // when non-empty, WriteFile returns failErr for that path
	failErr error
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
	_, fileOK := f.files[path]
	_, dirOK := f.dirs[path]
	return fileOK || dirOK, nil
}

func (f *fakeFS) ReadFile(path string) ([]byte, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
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

func (f *fakeFS) ReadDir(_ string) ([]iofs.DirEntry, error) {
	// Not needed for M3-T2 tests; return empty to satisfy the port.
	return nil, nil
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

// mkdirPaths returns the recorded MkdirAll paths sorted.
func (f *fakeFS) mkdirPaths() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]string, len(f.mkdirs))
	copy(out, f.mkdirs)
	sort.Strings(out)
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
