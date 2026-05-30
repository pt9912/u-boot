package application_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

const configTestBaseDir = "/tmp/u-boot-config-test/demo"

// newConfigService constructs a service against the in-memory FS
// plus the yaml.v3-backed fake codec. The BaseDir is pre-
// registered so the Exists() check on u-boot.yaml fails on
// file-absence, not on directory-absence.
func newConfigService(t *testing.T) (*application.ConfigService, *fakeFS) {
	t.Helper()
	fs := newFakeFS()
	fs.markDirExists(configTestBaseDir)
	y := &fakeYAML{}
	svc := application.NewConfigService(fs, y, nil)
	return svc, fs
}

func seedConfigUbootYAML(t *testing.T, fs *fakeFS) {
	t.Helper()
	body := "schemaVersion: 1\nproject:\n  name: t-uboot-config\n"
	if err := fs.WriteFile(filepath.Join(configTestBaseDir, "u-boot.yaml"),
		[]byte(body), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
}

func mustConfigPath(t *testing.T, raw string) domain.ConfigPath {
	t.Helper()
	p, err := domain.NewConfigPath(raw)
	if err != nil {
		t.Fatalf("NewConfigPath(%q): %v", raw, err)
	}
	return p
}

func TestConfig_BaseDirEmpty_NonSentinelError(t *testing.T) {
	t.Parallel()
	svc, _ := newConfigService(t)
	cases := []struct {
		name string
		fn   func() error
	}{
		{"Get", func() error {
			_, err := svc.Get(context.Background(), driving.ConfigGetRequest{Path: mustConfigPath(t, "project.name")})
			return err
		}},
		{"Set", func() error {
			_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
				Path: mustConfigPath(t, "project.name"), Value: "foo",
			})
			return err
		}},
		{"Show", func() error {
			_, err := svc.Show(context.Background(), driving.ConfigShowRequest{})
			return err
		}},
	}
	for _, tc := range cases {
		err := tc.fn()
		if err == nil {
			t.Errorf("%s: expected non-nil error for empty BaseDir", tc.name)
			continue
		}
		if errors.Is(err, driving.ErrProjectNotInitialized) {
			t.Errorf("%s: empty BaseDir leaked ErrProjectNotInitialized (it should be a non-sentinel error): %v", tc.name, err)
		}
	}
}

func TestConfig_NoUbootYAML_ReturnsErrProjectNotInitialized(t *testing.T) {
	t.Parallel()
	svc, _ := newConfigService(t) // no seedConfigUbootYAML — fresh BaseDir
	cases := []struct {
		name string
		fn   func() error
	}{
		{"Get", func() error {
			_, err := svc.Get(context.Background(), driving.ConfigGetRequest{
				BaseDir: configTestBaseDir, Path: mustConfigPath(t, "project.name"),
			})
			return err
		}},
		{"Set", func() error {
			_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
				BaseDir: configTestBaseDir, Path: mustConfigPath(t, "project.name"), Value: "foo",
			})
			return err
		}},
		{"Show", func() error {
			_, err := svc.Show(context.Background(), driving.ConfigShowRequest{
				BaseDir: configTestBaseDir,
			})
			return err
		}},
	}
	for _, tc := range cases {
		err := tc.fn()
		if !errors.Is(err, driving.ErrProjectNotInitialized) {
			t.Errorf("%s: err = %v, want wrap of ErrProjectNotInitialized", tc.name, err)
		}
	}
}

// TestConfig_StubHandlers pins that every M8-T2 handler is
// reachable past the project-state gate but unimplemented. T3
// reduces this to one (Set); T4 removes the test entirely once
// the last stub is gone.
func TestConfig_StubHandlers_ThreeReturnErrStubConfigHandler(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAML(t, fs)

	cases := []struct {
		name string
		fn   func() error
	}{
		{"Get", func() error {
			_, err := svc.Get(context.Background(), driving.ConfigGetRequest{
				BaseDir: configTestBaseDir, Path: mustConfigPath(t, "project.name"),
			})
			return err
		}},
		{"Set", func() error {
			_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
				BaseDir: configTestBaseDir, Path: mustConfigPath(t, "project.name"), Value: "foo",
			})
			return err
		}},
		{"Show", func() error {
			_, err := svc.Show(context.Background(), driving.ConfigShowRequest{
				BaseDir: configTestBaseDir,
			})
			return err
		}},
	}
	for _, tc := range cases {
		err := tc.fn()
		if err == nil {
			t.Errorf("%s: expected stub-handler error, got nil", tc.name)
			continue
		}
		if !errors.Is(err, application.ErrStubConfigHandlerForTest) {
			t.Errorf("%s: err = %v, want wrap of errStubConfigHandler", tc.name, err)
		}
	}
}
