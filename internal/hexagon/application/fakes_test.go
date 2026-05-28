package application_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	iofs "io/fs"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"

	yamladapter "github.com/pt9912/u-boot/internal/adapter/driven/yaml"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// fakeFS is an in-memory FileSystem implementation for application-
// layer tests. It records every WriteFile and MkdirAll call so tests
// can assert on creation order; Exists/Read/Rename/ReadDir/Lstat
// behave consistently with WriteFile/MkdirAll history. WriteFile
// and MkdirAll register *every* path ancestor in `dirs` so the
// fake's ReadDir matches the real adapter's behaviour for deeply
// nested keys.
//
// The fake is intentionally small — application tests need
// deterministic in-memory behaviour, not a full ioutil emulator.
type fakeFS struct {
	mu            sync.Mutex
	files         map[string][]byte
	fileModes     map[string]iofs.FileMode
	dirs          map[string]bool
	dirModes      map[string]iofs.FileMode
	symlinks      map[string]bool // path is a symlink for Lstat purposes
	writes        []string        // ordered: every successful WriteFile path
	mkdirs        []string         // ordered: every MkdirAll path
	failOn        string           // when non-empty, WriteFile / WriteFileExclusive returns failErr for that path
	failErr       error
	failReadOn    string // when non-empty, ReadFile returns failReadErr for that path
	failReadErr   error
	failExistsOn  string // when non-empty, Exists returns failExistsErr for that path
	failExistsErr error
	failLstatOn   string // when non-empty, Lstat returns failLstatErr for that path
	failLstatErr  error
	failMkdirOn   string // when non-empty, Mkdir / MkdirAll returns failMkdirErr for that path
	failMkdirErr  error
	failReadDirOn string // when non-empty, ReadDir returns failReadDirErr for that path
	failReadDirErr error
	failRemoveAll  error          // when non-nil, RemoveAll returns this error
	readFileCalls  map[string]int // per-path ReadFile call counter (tests assert no double-reads)
}

func newFakeFS() *fakeFS {
	return &fakeFS{
		files:         make(map[string][]byte),
		fileModes:     make(map[string]iofs.FileMode),
		dirs:          make(map[string]bool),
		dirModes:      make(map[string]iofs.FileMode),
		symlinks:      make(map[string]bool),
		readFileCalls: make(map[string]int),
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
	f.readFileCalls[path]++
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

// readFileCallCount returns how many times ReadFile was invoked for
// path. Tests use it to pin the "no double-read" invariant from the
// T4b-review (plan caches the body, execute uses the cache).
func (f *fakeFS) readFileCallCount(path string) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.readFileCalls[path]
}

func (f *fakeFS) WriteFile(path string, data []byte, mode iofs.FileMode) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failOn != "" && path == f.failOn {
		return f.failErr
	}
	f.writeFileLocked(path, data, mode)
	return nil
}

func (f *fakeFS) WriteFileExclusive(path string, data []byte, mode iofs.FileMode) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failOn != "" && path == f.failOn {
		return f.failErr
	}
	if _, ok := f.files[path]; ok {
		return iofs.ErrExist
	}
	if f.dirs[path] {
		return iofs.ErrExist
	}
	f.writeFileLocked(path, data, mode)
	return nil
}

// writeFileLocked stores the file plus its mode and registers every
// ancestor directory so ReadDir reflects the implicit tree (matches
// os.WriteFile + os.MkdirAll behaviour). Caller must hold f.mu.
func (f *fakeFS) writeFileLocked(path string, data []byte, mode iofs.FileMode) {
	stored := make([]byte, len(data))
	copy(stored, data)
	f.files[path] = stored
	f.fileModes[path] = mode
	f.registerAncestorsLocked(filepath.Dir(path))
	f.writes = append(f.writes, path)
}

func (f *fakeFS) Mkdir(path string, mode iofs.FileMode) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failMkdirOn != "" && path == f.failMkdirOn {
		return f.failMkdirErr
	}
	if _, ok := f.files[path]; ok {
		return iofs.ErrExist
	}
	if f.dirs[path] {
		return iofs.ErrExist
	}
	f.dirs[path] = true
	f.dirModes[path] = mode
	f.mkdirs = append(f.mkdirs, path)
	return nil
}

func (f *fakeFS) MkdirAll(path string, mode iofs.FileMode) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failMkdirOn != "" && path == f.failMkdirOn {
		return f.failMkdirErr
	}
	if !f.dirs[path] {
		f.dirs[path] = true
		f.dirModes[path] = mode
	}
	f.registerAncestorsLocked(filepath.Dir(path))
	f.mkdirs = append(f.mkdirs, path)
	return nil
}

// registerAncestorsLocked walks up from `start` and marks every
// directory ancestor as existing. Stops when filepath.Dir is a
// fixed point ("/" on POSIX, "." on relative paths). Caller must
// hold f.mu. Ancestor dirModes default to 0o755 so Lstat returns a
// sensible mode for implicit intermediate directories.
func (f *fakeFS) registerAncestorsLocked(start string) {
	p := start
	for p != "" && p != "." {
		if !f.dirs[p] {
			f.dirs[p] = true
			f.dirModes[p] = 0o755
		}
		parent := filepath.Dir(p)
		if parent == p {
			break
		}
		p = parent
	}
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

// ReadDir returns the direct children of path. Tests inject failures
// via failReadDirOn / failReadDirErr.
func (f *fakeFS) ReadDir(path string) ([]iofs.DirEntry, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failReadDirOn != "" && path == f.failReadDirOn {
		return nil, f.failReadDirErr
	}
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
		// Symlinks beat dir classification — Lstat-time semantics
		// match os.ReadDir's: an entry is a symlink if the link
		// itself exists, regardless of what it points at.
		if f.symlinks[direct] {
			entries = append(entries, fakeDirEntry{name: name, isDir: false})
			return
		}
		if f.dirs[direct] {
			entries = append(entries, fakeDirEntry{name: name, isDir: true})
			return
		}
		entries = append(entries, fakeDirEntry{name: name, isDir: false})
	}
	for k := range f.files {
		addChild(k)
	}
	for k := range f.dirs {
		addChild(k)
	}
	for k := range f.symlinks {
		addChild(k)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	return entries, nil
}

// Lstat returns FileInfo for path without following symlinks.
// Behaviour matches os.Lstat: symlinks report ModeSymlink, regular
// files report the recorded mode, directories report ModeDir | mode.
// Missing paths return iofs.ErrNotExist.
func (f *fakeFS) Lstat(path string) (iofs.FileInfo, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failLstatOn != "" && path == f.failLstatOn {
		return nil, f.failLstatErr
	}
	if f.symlinks[path] {
		return &fakeFileInfo{name: filepath.Base(path), mode: iofs.ModeSymlink}, nil
	}
	if data, ok := f.files[path]; ok {
		return &fakeFileInfo{
			name: filepath.Base(path),
			size: int64(len(data)),
			mode: f.fileModes[path],
		}, nil
	}
	if f.dirs[path] {
		return &fakeFileInfo{
			name: filepath.Base(path),
			mode: iofs.ModeDir | f.dirModes[path],
		}, nil
	}
	return nil, iofs.ErrNotExist
}

// RemoveAll deletes path and any recorded children. Idempotent —
// removing a missing path is a no-op, matching os.RemoveAll. Tests
// inject a forced failure via failRemoveAll.
func (f *fakeFS) RemoveAll(path string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failRemoveAll != nil {
		return f.failRemoveAll
	}
	delete(f.files, path)
	delete(f.fileModes, path)
	delete(f.dirs, path)
	delete(f.dirModes, path)
	delete(f.symlinks, path)
	prefix := path + string(filepath.Separator)
	for k := range f.files {
		if strings.HasPrefix(k, prefix) {
			delete(f.files, k)
			delete(f.fileModes, k)
		}
	}
	for k := range f.dirs {
		if strings.HasPrefix(k, prefix) {
			delete(f.dirs, k)
			delete(f.dirModes, k)
		}
	}
	for k := range f.symlinks {
		if strings.HasPrefix(k, prefix) {
			delete(f.symlinks, k)
		}
	}
	return nil
}

// Copy streams src to dst (truncate-overwrites). In the in-memory
// fake, "streaming" is a no-op slot: we just clone the byte slice.
// Honors failReadOn (for src) and failOn (for dst) so the existing
// error-injection knobs work transparently across the
// ReadFile+WriteFile → Copy refactor.
func (f *fakeFS) Copy(src, dst string, mode iofs.FileMode) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failReadOn != "" && src == f.failReadOn {
		return f.failReadErr
	}
	if f.failOn != "" && dst == f.failOn {
		return f.failErr
	}
	data, ok := f.files[src]
	if !ok {
		return fmt.Errorf("open %s: %w", src, iofs.ErrNotExist)
	}
	body := append([]byte(nil), data...)
	f.writeFileLocked(dst, body, mode)
	return nil
}

// CopyExclusive streams src to dst with O_EXCL — returns wrapped
// iofs.ErrExist if dst is already present. Same error-injection
// hooks as [Copy].
func (f *fakeFS) CopyExclusive(src, dst string, mode iofs.FileMode) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failReadOn != "" && src == f.failReadOn {
		return f.failReadErr
	}
	if f.failOn != "" && dst == f.failOn {
		return f.failErr
	}
	if _, ok := f.files[dst]; ok {
		return fmt.Errorf("create %s: %w", dst, iofs.ErrExist)
	}
	if f.dirs[dst] {
		return fmt.Errorf("create %s: %w", dst, iofs.ErrExist)
	}
	data, ok := f.files[src]
	if !ok {
		return fmt.Errorf("open %s: %w", src, iofs.ErrNotExist)
	}
	body := append([]byte(nil), data...)
	f.writeFileLocked(dst, body, mode)
	return nil
}

// markDirExists pre-registers a directory so Exists returns true.
// Used by test setup to satisfy the BaseDir-existence check
// without going through a real MkdirAll call (which the test
// otherwise wants to count). Mode defaults to 0o755.
func (f *fakeFS) markDirExists(path string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.dirs[path] = true
	if _, ok := f.dirModes[path]; !ok {
		f.dirModes[path] = 0o755
	}
}

// markSymlink registers path as a symlink for Lstat purposes and
// makes the entry appear in ReadDir of its parent — so the backup
// walker can encounter and reject it. Ancestors are registered the
// same way WriteFile does for files (review finding #7).
func (f *fakeFS) markSymlink(path string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.symlinks[path] = true
	f.registerAncestorsLocked(filepath.Dir(path))
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

// fakeDirEntry is the minimal iofs.DirEntry the BackupPath walker
// needs: Name() and IsDir() only. Type()/Info() return the zero
// values because the algorithm never consults them.
type fakeDirEntry struct {
	name  string
	isDir bool
}

func (e fakeDirEntry) Name() string                 { return e.name }
func (e fakeDirEntry) IsDir() bool                  { return e.isDir }
func (e fakeDirEntry) Type() iofs.FileMode          { return 0 }
func (e fakeDirEntry) Info() (iofs.FileInfo, error) { return nil, errors.New("fakeDirEntry.Info: not implemented") }

// fakeFileInfo backs the Lstat return value with just the fields
// the backup service consults: Name / Size / Mode / IsDir.
type fakeFileInfo struct {
	name string
	size int64
	mode iofs.FileMode
}

func (i *fakeFileInfo) Name() string       { return i.name }
func (i *fakeFileInfo) Size() int64        { return i.size }
func (i *fakeFileInfo) Mode() iofs.FileMode { return i.mode }
func (i *fakeFileInfo) ModTime() time.Time { return time.Time{} }
func (i *fakeFileInfo) IsDir() bool        { return i.mode.IsDir() }
func (i *fakeFileInfo) Sys() any           { return nil }

// fakeYAML uses gopkg.in/yaml.v3 directly. The application layer
// imports the YAMLCodec port, the test imports the library — the
// `*_test.go` Carveout (LH-FA-ARCH-003) covers that.
type fakeYAML struct {
	failMarshal bool
	// failPatchOn forces PatchScalar to fail when the patched path's
	// joined-by-"." representation equals this string. Tests assert
	// no FS write happens on a patch failure.
	failPatchOn  string
	failPatchErr error
	// patchCalls records every PatchScalar invocation so tests can
	// assert call ordering / count.
	patchCalls []fakePatchCall
	// failPatchMappingOn / failPatchMappingErr fail
	// PatchMappingEntryYAML for a specific parentKey.entryKey pair.
	failPatchMappingOn  string
	failPatchMappingErr error
	// mappingPatchCalls records every PatchMappingEntryYAML call.
	mappingPatchCalls []fakeMappingPatchCall
}

type fakePatchCall struct {
	Path  []string
	Value any
}

type fakeMappingPatchCall struct {
	ParentKey  string
	EntryKey   string
	ValueYAML  []byte
	MarkerName string
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

// PatchScalar mirrors the production yaml.Node-based adapter just
// closely enough for tests: it parses content, walks path through
// nested mappings (creating missing intermediates), and sets a
// scalar leaf. Errors mirror the [driven.ErrYAMLPathInvalid]
// contract so tests can assert on the sentinel without depending on
// the production adapter.
func (f *fakeYAML) PatchScalar(content []byte, path []string, value any) ([]byte, error) {
	if f.failPatchOn != "" && strings.Join(path, ".") == f.failPatchOn {
		return nil, f.failPatchErr
	}
	f.patchCalls = append(f.patchCalls, fakePatchCall{Path: append([]string(nil), path...), Value: value})
	if len(path) == 0 {
		return nil, fmt.Errorf("%w: empty path", driven.ErrYAMLPathInvalid)
	}
	scalar, err := fakeScalarNode(value)
	if err != nil {
		return nil, err
	}
	var doc yaml.Node
	if len(bytes.TrimSpace(content)) > 0 {
		if err := yaml.Unmarshal(content, &doc); err != nil {
			return nil, fmt.Errorf("fakeYAML.PatchScalar parse: %w", err)
		}
	}
	root := fakeDocumentRoot(&doc)
	if root.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("%w: root is not a mapping", driven.ErrYAMLPathInvalid)
	}
	cur := root
	for i, key := range path {
		isLast := i == len(path)-1
		child, found := fakeMappingChild(cur, key)
		switch {
		case isLast && found:
			*child = *scalar
		case isLast && !found:
			fakeAppendMappingChild(cur, key, scalar)
		case !isLast && found:
			if child.Kind != yaml.MappingNode {
				return nil, fmt.Errorf("%w: intermediate key %q is not a mapping",
					driven.ErrYAMLPathInvalid, key)
			}
			cur = child
		case !isLast && !found:
			next := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
			fakeAppendMappingChild(cur, key, next)
			cur = next
		}
	}
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(&doc); err != nil {
		return nil, fmt.Errorf("fakeYAML.PatchScalar encode: %w", err)
	}
	if err := enc.Close(); err != nil {
		return nil, fmt.Errorf("fakeYAML.PatchScalar close: %w", err)
	}
	return buf.Bytes(), nil
}

func fakeDocumentRoot(doc *yaml.Node) *yaml.Node {
	if doc.Kind == 0 {
		doc.Kind = yaml.DocumentNode
		doc.Content = []*yaml.Node{{Kind: yaml.MappingNode, Tag: "!!map"}}
		return doc.Content[0]
	}
	if doc.Kind == yaml.DocumentNode {
		if len(doc.Content) == 0 {
			doc.Content = []*yaml.Node{{Kind: yaml.MappingNode, Tag: "!!map"}}
		}
		return doc.Content[0]
	}
	return doc
}

func fakeMappingChild(mapping *yaml.Node, key string) (*yaml.Node, bool) {
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		k := mapping.Content[i]
		if k.Kind == yaml.ScalarNode && k.Value == key {
			return mapping.Content[i+1], true
		}
	}
	return nil, false
}

func fakeAppendMappingChild(mapping *yaml.Node, key string, value *yaml.Node) {
	mapping.Content = append(mapping.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key},
		value,
	)
}

func fakeScalarNode(value any) (*yaml.Node, error) {
	switch v := value.(type) {
	case bool:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: strconv.FormatBool(v)}, nil
	case string:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: v}, nil
	case int:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: strconv.FormatInt(int64(v), 10)}, nil
	default:
		return nil, fmt.Errorf("%w: unsupported %T", driven.ErrYAMLPathInvalid, value)
	}
}

// PatchMappingEntryYAML delegates to the production yaml adapter so
// application-layer tests see the exact same byte-preservation,
// anchor-validation and error-sentinel behaviour the real CLI sees.
// Direct adapter imports are allowed inside `_test.go` files via the
// `application-no-adapter` depguard carveout.
//
// The fake adds two failure hooks (failPatchMappingOn /
// failPatchMappingErr) so tests can assert that no FS write happens
// when PatchMappingEntryYAML fails mid-plan.
func (f *fakeYAML) PatchMappingEntryYAML(content []byte, parentKey, entryKey string, valueYAML []byte, markerName string) ([]byte, error) {
	if f.failPatchMappingOn != "" && f.failPatchMappingOn == parentKey+"."+entryKey {
		return nil, f.failPatchMappingErr
	}
	f.mappingPatchCalls = append(f.mappingPatchCalls, fakeMappingPatchCall{
		ParentKey:  parentKey,
		EntryKey:   entryKey,
		ValueYAML:  append([]byte(nil), valueYAML...),
		MarkerName: markerName,
	})
	return yamladapter.New().PatchMappingEntryYAML(content, parentKey, entryKey, valueYAML, markerName)
}

// LocateMarkedEntry also delegates so the LocateResult flags the
// fake produces are byte-identical to production. No failure hook —
// LocateMarkedEntry is read-only and is exercised on the happy path
// by the application-layer T4c tests.
func (f *fakeYAML) LocateMarkedEntry(content []byte, parentKey, entryKey, markerName string) (driven.LocateResult, error) {
	return yamladapter.New().LocateMarkedEntry(content, parentKey, entryKey, markerName)
}

// fakeProgress records every AffectedFiles call so tests can
// assert on the structured events the application emits via the
// driven.ProgressPort. A nil pointer is acceptable for tests that
// do not care about progress (the service constructor swaps in a
// no-op); fakeProgress is for tests that DO want to inspect.
type fakeProgress struct {
	calls []fakeProgressCall
}

type fakeProgressCall struct {
	BaseDir string
	Rows    []driven.AffectedFile
}

func (p *fakeProgress) AffectedFiles(baseDir string, rows []driven.AffectedFile) {
	// Defensive copy so the test sees a stable snapshot.
	rowsCopy := make([]driven.AffectedFile, len(rows))
	copy(rowsCopy, rows)
	p.calls = append(p.calls, fakeProgressCall{BaseDir: baseDir, Rows: rowsCopy})
}

// fakeGit records IsRepository / Init / Version calls and lets each
// test configure the return values + errors.
type fakeGit struct {
	isRepoCalls []string
	initCalls   []string
	versionCalls int
	isRepo      bool
	isRepoErr   error
	initErr     error
	version     string
	versionErr  error
}

func (f *fakeGit) IsRepository(_ context.Context, dir string) (bool, error) {
	f.isRepoCalls = append(f.isRepoCalls, dir)
	return f.isRepo, f.isRepoErr
}

func (f *fakeGit) Init(_ context.Context, dir string) error {
	f.initCalls = append(f.initCalls, dir)
	return f.initErr
}

func (f *fakeGit) Version(_ context.Context) (string, error) {
	f.versionCalls++
	return f.version, f.versionErr
}

// fakeDockerProbe is the in-memory DockerProbe used by Doctor-service
// tests. Each method has a separate value+error pair so tests can
// dial in arbitrary combinations (e.g. binary present but daemon
// unreachable: Version OK, Info errors).
type fakeDockerProbe struct {
	versionCalls        int
	infoCalls           int
	composeVersionCalls int
	version             string
	versionErr          error
	infoErr             error
	composeVersion      string
	composeVersionErr   error
}

func (p *fakeDockerProbe) Version(_ context.Context) (string, error) {
	p.versionCalls++
	return p.version, p.versionErr
}

func (p *fakeDockerProbe) Info(_ context.Context) error {
	p.infoCalls++
	return p.infoErr
}

func (p *fakeDockerProbe) ComposeVersion(_ context.Context) (string, error) {
	p.composeVersionCalls++
	return p.composeVersion, p.composeVersionErr
}

// fakeConfirmer records every ConfirmTreatAsExisting call and returns
// the configured answer/err. The default zero value answers false
// (the "no" default of the soft-detection prompt) so a freshly
// constructed instance models a deterministic "user declined".
type fakeConfirmer struct {
	calls  []fakeConfirmerCall
	answer bool
	err    error
}

type fakeConfirmerCall struct {
	BaseDir    string
	Indicators []string
}

func (c *fakeConfirmer) ConfirmTreatAsExisting(_ context.Context, baseDir string, indicators []string) (bool, error) {
	c.calls = append(c.calls, fakeConfirmerCall{BaseDir: baseDir, Indicators: append([]string{}, indicators...)})
	return c.answer, c.err
}

// fakeLogger records every Debug/Info/Warn/Error call so tests can
// assert on the LH-QA-004 driven.Logger port usage. args are
// captured verbatim (slog's alternating key/value convention).
type fakeLogger struct {
	entries []fakeLogEntry
}

type fakeLogEntry struct {
	Level string
	Msg   string
	Args  []any
}

func (l *fakeLogger) Debug(msg string, args ...any) { l.record("DEBUG", msg, args) }
func (l *fakeLogger) Info(msg string, args ...any)  { l.record("INFO", msg, args) }
func (l *fakeLogger) Warn(msg string, args ...any)  { l.record("WARN", msg, args) }
func (l *fakeLogger) Error(msg string, args ...any) { l.record("ERROR", msg, args) }

func (l *fakeLogger) record(level, msg string, args []any) {
	l.entries = append(l.entries, fakeLogEntry{Level: level, Msg: msg, Args: append([]any{}, args...)})
}
