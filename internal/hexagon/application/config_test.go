package application_test

import (
	"bytes"
	"context"
	"errors"
	"path/filepath"
	"strings"
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

// TestConfig_StubHandler_OnlySetReturnsErrStubConfigHandler pins
// the remaining unimplemented handler (T4 fills it). T3 reduces
// the pin count from 3 (T2) to 1; T4 removes the test entirely.
func TestConfig_StubHandler_OnlySetReturnsErrStubConfigHandler(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAML(t, fs)

	_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "project.name"),
		Value:   "foo",
	})
	if err == nil {
		t.Fatal("expected stub-handler error, got nil")
	}
	if !errors.Is(err, application.ErrStubConfigHandlerForTest) {
		t.Errorf("err = %v, want wrap of errStubConfigHandler", err)
	}
}

// --- T3: Get + Show -------------------------------------------------

func TestConfigGet_ProjectName_ReturnsValue(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAML(t, fs)

	resp, err := svc.Get(context.Background(), driving.ConfigGetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "project.name"),
	})
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got, want := resp.Value, "t-uboot-config"; got != want {
		t.Errorf("Value = %q, want %q", got, want)
	}
	if resp.Path.String() != "project.name" {
		t.Errorf("Path = %v, want project.name", resp.Path)
	}
}

func TestConfigGet_DevcontainerEnabled_FromExistingBlock(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	body := "schemaVersion: 1\nproject:\n  name: t-uboot-config\ndevcontainer:\n  enabled: true\n"
	if err := fs.WriteFile(configTestBaseDir+"/u-boot.yaml", []byte(body), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	resp, err := svc.Get(context.Background(), driving.ConfigGetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "devcontainer.enabled"),
	})
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got, want := resp.Value, "true"; got != want {
		t.Errorf("Value = %q, want %q", got, want)
	}
}

func TestConfigGet_DevcontainerEnabled_MissingBlock_ReturnsNotSet(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAML(t, fs) // no devcontainer block

	_, err := svc.Get(context.Background(), driving.ConfigGetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "devcontainer.enabled"),
	})
	if !errors.Is(err, driving.ErrConfigValueNotSet) {
		t.Fatalf("err = %v, want wrap of ErrConfigValueNotSet", err)
	}
	// User-visible hint must point at the canonical write path.
	if !strings.Contains(err.Error(), "init --devcontainer") &&
		!strings.Contains(err.Error(), "config set devcontainer.enabled") {
		t.Errorf("error message %q lacks a write-path hint", err.Error())
	}
}

func TestConfigGet_ServiceEnabled_FromExistingMap(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	body := "schemaVersion: 1\nproject:\n  name: t-uboot-config\nservices:\n  postgres:\n    enabled: false\n"
	if err := fs.WriteFile(configTestBaseDir+"/u-boot.yaml", []byte(body), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	resp, err := svc.Get(context.Background(), driving.ConfigGetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "services.postgres.enabled"),
	})
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got, want := resp.Value, "false"; got != want {
		t.Errorf("Value = %q, want %q", got, want)
	}
}

func TestConfigGet_ServiceEnabled_MissingService_ReturnsNotSet(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAML(t, fs) // no services block

	_, err := svc.Get(context.Background(), driving.ConfigGetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "services.postgres.enabled"),
	})
	if !errors.Is(err, driving.ErrConfigValueNotSet) {
		t.Fatalf("err = %v, want wrap of ErrConfigValueNotSet", err)
	}
	if !strings.Contains(err.Error(), "u-boot add postgres") {
		t.Errorf("error message %q lacks `u-boot add postgres` hint", err.Error())
	}
}

func TestConfigGet_ProjectName_MissingName_ReturnsSchemaInvalid(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	// schemaVersion present but project.name missing — corrupt
	// config (LH-FA-CONF-002 §1308 requires the name).
	if err := fs.WriteFile(configTestBaseDir+"/u-boot.yaml",
		[]byte("schemaVersion: 1\nproject: {}\n"), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}

	_, err := svc.Get(context.Background(), driving.ConfigGetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "project.name"),
	})
	if !errors.Is(err, driving.ErrConfigSchemaInvalid) {
		t.Errorf("err = %v, want wrap of ErrConfigSchemaInvalid", err)
	}
}

func TestConfigShow_ReturnsByteIdenticalBody(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	// Seed with comments + non-default formatting to pin that
	// Show does NOT re-marshal.
	body := []byte("# u-boot project config\nschemaVersion: 1\nproject:\n  name: t-uboot-config  # display name\n")
	if err := fs.WriteFile(configTestBaseDir+"/u-boot.yaml", body, 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}

	resp, err := svc.Show(context.Background(), driving.ConfigShowRequest{
		BaseDir: configTestBaseDir,
	})
	if err != nil {
		t.Fatalf("Show: %v", err)
	}
	if !bytes.Equal(resp.Body, body) {
		t.Errorf("Show.Body byte-mismatch:\n got:%q\nwant:%q", resp.Body, body)
	}
}

func TestConfigGet_CorruptYAML_ReturnsSchemaInvalid(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	// Syntactically corrupt YAML: leading `:`.
	if err := fs.WriteFile(configTestBaseDir+"/u-boot.yaml",
		[]byte(":\n  bad yaml\n"), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}

	_, err := svc.Get(context.Background(), driving.ConfigGetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "project.name"),
	})
	if !errors.Is(err, driving.ErrConfigSchemaInvalid) {
		t.Errorf("err = %v, want wrap of ErrConfigSchemaInvalid", err)
	}
}
