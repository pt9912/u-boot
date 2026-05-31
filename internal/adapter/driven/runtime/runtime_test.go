package runtime_test

import (
	iofs "io/fs"
	"os"
	"testing"
	"time"

	"github.com/pt9912/u-boot/internal/adapter/driven/runtime"
)

// fakeStatTable is a stat() seam that returns a configured map of
// "this path exists" markers. Any path not in the map returns
// os.ErrNotExist (the natural "no marker" answer). Unexpected
// errors (e.g. permission denied) are simulated by mapping the
// path to a non-nil non-ErrNotExist error.
type fakeStatTable struct {
	existing  map[string]bool
	hardError map[string]error
}

func (f fakeStatTable) stat(name string) (os.FileInfo, error) {
	if err, ok := f.hardError[name]; ok {
		return nil, err
	}
	if f.existing[name] {
		return fakeFileInfo{}, nil
	}
	return nil, &iofs.PathError{Op: "stat", Path: name, Err: iofs.ErrNotExist}
}

// fakeFileInfo is the minimal os.FileInfo we need. None of the
// fields are inspected by the adapter — only the err return value
// from Stat matters.
type fakeFileInfo struct{}

func (fakeFileInfo) Name() string       { return "" }
func (fakeFileInfo) Size() int64        { return 0 }
func (fakeFileInfo) Mode() iofs.FileMode { return 0 }
func (fakeFileInfo) ModTime() time.Time  { return time.Time{} }
func (fakeFileInfo) IsDir() bool         { return false }
func (fakeFileInfo) Sys() any            { return nil }

func TestFileEnv_InContainer_Detection(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		existing map[string]bool
		hardErr  map[string]error
		want     bool
	}{
		{
			name:     "no markers — host",
			existing: map[string]bool{},
			want:     false,
		},
		{
			name:     "only docker marker — Docker Engine / Desktop",
			existing: map[string]bool{"/.dockerenv": true},
			want:     true,
		},
		{
			name:     "only podman marker — Podman / CRI-O",
			existing: map[string]bool{"/run/.containerenv": true},
			want:     true,
		},
		{
			name: "both markers present — defensive",
			existing: map[string]bool{
				"/.dockerenv":        true,
				"/run/.containerenv": true,
			},
			want: true,
		},
		{
			name:     "unrelated path 'exists' — must not match",
			existing: map[string]bool{"/etc/hostname": true},
			want:     false,
		},
		{
			name:     "permission denied on docker marker — best-effort no",
			hardErr:  map[string]error{"/.dockerenv": iofs.ErrPermission},
			existing: map[string]bool{},
			want:     false,
		},
		{
			name: "permission denied on docker, podman marker present — yes",
			hardErr: map[string]error{
				"/.dockerenv": iofs.ErrPermission,
			},
			existing: map[string]bool{"/run/.containerenv": true},
			want:     true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fake := fakeStatTable{existing: tc.existing, hardError: tc.hardErr}
			env := runtime.NewWithStat(fake.stat)

			if got := env.InContainer(); got != tc.want {
				t.Fatalf("InContainer() = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestFileEnv_Production_DoesNotPanic exercises the real
// os.Stat-backed adapter (no fake). The result depends on the
// host running the tests, so we only assert that the call
// returns without panicking. CI runs the gates Docker stage
// which DOES have /.dockerenv (gates run inside `docker build`),
// so the test acts as a real-world smoke that the production
// path works end-to-end.
func TestFileEnv_Production_DoesNotPanic(t *testing.T) {
	t.Parallel()

	_ = runtime.New().InContainer()
}
