package application_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
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

// T2's stub-pin test is removed in T4 — all three handlers are
// real now. errStubConfigHandler + ErrStubConfigHandlerForTest
// are removed along with this test (slice-m8-config.md §T4 DoD).

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

// --- T4: Set ----------------------------------------------------------

func TestConfigSet_ProjectName_HappyPath(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAML(t, fs)

	resp, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "project.name"),
		Value:   "new-name",
	})
	if err != nil {
		t.Fatalf("Set: %v", err)
	}
	if got, want := resp.OldValue, "t-uboot-config"; got != want {
		t.Errorf("OldValue = %q, want %q", got, want)
	}
	if got, want := resp.NewValue, "new-name"; got != want {
		t.Errorf("NewValue = %q, want %q", got, want)
	}
	// File now carries the new name.
	body, err := fs.ReadFile(configTestBaseDir + "/u-boot.yaml")
	if err != nil {
		t.Fatalf("read post-set: %v", err)
	}
	if !bytes.Contains(body, []byte("name: new-name")) {
		t.Errorf("post-set body missing new name:\n%s", body)
	}
}

func TestConfigSet_ProjectName_InvalidName_NoWrite(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAML(t, fs)

	writesBefore := len(fs.writtenPaths())
	_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "project.name"),
		Value:   "Invalid Name", // capital + space — fails LH-FA-INIT-006 regex
	})
	if !errors.Is(err, driving.ErrConfigValueInvalid) {
		t.Fatalf("err = %v, want wrap of ErrConfigValueInvalid", err)
	}
	// Transactional contract: no WriteFile on validation failure.
	if delta := len(fs.writtenPaths()) - writesBefore; delta != 0 {
		t.Errorf("invalid-name set produced %d WriteFile call(s), want 0; writes = %v",
			delta, fs.writtenPaths())
	}
}

// TestConfigSet_ServicesEnabled_Rejected_WithHint pins the M1
// review-fix: services.<svc>.enabled is Get-only because the
// LH-FA-ADD-005 state machine owns the lifecycle. Set must
// reject with ErrConfigValueInvalid + a hint pointing at
// `u-boot add <svc>`.
func TestConfigSet_ServicesEnabled_Rejected_WithHint(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAML(t, fs)

	writesBefore := len(fs.writtenPaths())
	_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "services.postgres.enabled"),
		Value:   "true",
	})
	if !errors.Is(err, driving.ErrConfigValueInvalid) {
		t.Fatalf("err = %v, want wrap of ErrConfigValueInvalid", err)
	}
	if !strings.Contains(err.Error(), "u-boot add postgres") {
		t.Errorf("error message %q lacks `u-boot add postgres` hint", err.Error())
	}
	if delta := len(fs.writtenPaths()) - writesBefore; delta != 0 {
		t.Errorf("rejected set produced %d WriteFile call(s), want 0", delta)
	}
}

func TestConfigSet_DevcontainerEnabled_HappyPath(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAML(t, fs)

	resp, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "devcontainer.enabled"),
		Value:   "true",
	})
	if err != nil {
		t.Fatalf("Set: %v", err)
	}
	if resp.OldValue != "" {
		t.Errorf("OldValue = %q, want empty (initial unset)", resp.OldValue)
	}
	if resp.NewValue != "true" {
		t.Errorf("NewValue = %q, want %q", resp.NewValue, "true")
	}
	// Get returns the value now.
	getResp, err := svc.Get(context.Background(), driving.ConfigGetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "devcontainer.enabled"),
	})
	if err != nil {
		t.Fatalf("post-set Get: %v", err)
	}
	if getResp.Value != "true" {
		t.Errorf("Get-after-Set Value = %q, want %q", getResp.Value, "true")
	}
}

func TestConfigSet_DevcontainerEnabled_InvalidBool_NoWrite(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAML(t, fs)

	writesBefore := len(fs.writtenPaths())
	_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "devcontainer.enabled"),
		Value:   "vielleicht",
	})
	if !errors.Is(err, driving.ErrConfigValueInvalid) {
		t.Fatalf("err = %v, want wrap of ErrConfigValueInvalid", err)
	}
	if delta := len(fs.writtenPaths()) - writesBefore; delta != 0 {
		t.Errorf("invalid-bool set produced %d WriteFile call(s), want 0", delta)
	}
}

// TestConfigSet_NoOp_SameValue_SkipsWrite pins the §T4 NoOp
// short-circuit: setting a path to its current value produces a
// response with OldValue == NewValue and no WriteFile call.
func TestConfigSet_NoOp_SameValue_SkipsWrite(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAML(t, fs)

	writesBefore := len(fs.writtenPaths())
	resp, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "project.name"),
		Value:   "t-uboot-config", // same as seed
	})
	if err != nil {
		t.Fatalf("Set: %v", err)
	}
	if resp.OldValue != resp.NewValue {
		t.Errorf("OldValue %q != NewValue %q (expected idempotent NoOp)", resp.OldValue, resp.NewValue)
	}
	if delta := len(fs.writtenPaths()) - writesBefore; delta != 0 {
		t.Errorf("NoOp set produced %d WriteFile call(s), want 0", delta)
	}
}

// TestConfigSet_CommentsSurviveRoundTrip pins that the
// yaml.v3-Node-API-backed PatchScalar preserves comments around
// the touched field. Important because Show is byte-identical
// (§D5); Set must not silently strip comments.
func TestConfigSet_CommentsSurviveRoundTrip(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	body := []byte("# u-boot project config\nschemaVersion: 1\nproject:\n  name: original  # display name\n")
	if err := fs.WriteFile(configTestBaseDir+"/u-boot.yaml", body, 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}

	_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "project.name"),
		Value:   "renamed",
	})
	if err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := fs.ReadFile(configTestBaseDir + "/u-boot.yaml")
	if err != nil {
		t.Fatalf("read post-set: %v", err)
	}
	if !bytes.Contains(got, []byte("# u-boot project config")) {
		t.Errorf("head comment lost; got:\n%s", got)
	}
	if !bytes.Contains(got, []byte("# display name")) {
		t.Errorf("line-trailing comment lost; got:\n%s", got)
	}
	if !bytes.Contains(got, []byte("name: renamed")) {
		t.Errorf("new value not written; got:\n%s", got)
	}
}

// TestConfigSet_WriteFileFails_ReturnsFileSystem covers the
// ErrConfigFileSystem path after a successful patch + validate:
// if the final WriteFile fails, the sentinel surfaces as
// LH-FA-CLI-006 code 14 (technical).
func TestConfigSet_WriteFileFails_ReturnsFileSystem(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAML(t, fs)
	// Arm a WriteFile failure for u-boot.yaml.
	fs.mu.Lock()
	fs.failOn = filepath.Join(configTestBaseDir, "u-boot.yaml")
	fs.failErr = errors.New("disk full")
	fs.mu.Unlock()

	_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "project.name"),
		Value:   "new-name",
	})
	if !errors.Is(err, driving.ErrConfigFileSystem) {
		t.Errorf("err = %v, want wrap of ErrConfigFileSystem", err)
	}
}

// TestConfigSet_PatchScalarFails_ReturnsSchemaInvalid covers the
// PatchScalar-failure branch: when the YAML codec rejects the
// patch (e.g. corrupt content the helper failed to detect
// upfront), the Set surfaces as ErrConfigSchemaInvalid + no
// WriteFile.
func TestConfigSet_PatchScalarFails_ReturnsSchemaInvalid(t *testing.T) {
	t.Parallel()
	fs := newFakeFS()
	fs.markDirExists(configTestBaseDir)
	y := &fakeYAML{
		failPatchOn:  "project.name",
		failPatchErr: fmt.Errorf("PatchScalar mock failure: %w", driven.ErrYAMLParse),
	}
	svc := application.NewConfigService(fs, y, nil)
	if err := fs.WriteFile(configTestBaseDir+"/u-boot.yaml",
		[]byte("schemaVersion: 1\nproject:\n  name: t-uboot-config\n"), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}

	writesBefore := len(fs.writtenPaths())
	_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "project.name"),
		Value:   "new-name",
	})
	if !errors.Is(err, driving.ErrConfigSchemaInvalid) {
		t.Errorf("err = %v, want wrap of ErrConfigSchemaInvalid", err)
	}
	if delta := len(fs.writtenPaths()) - writesBefore; delta != 0 {
		t.Errorf("PatchScalar-fail produced %d WriteFile call(s), want 0", delta)
	}
}

// TestConfigSet_LHFACONF005_RoundTripPin reproduces the spec
// example flow verbatim:
//
//	u-boot config get project.name
//	u-boot config set project.name foo
//	u-boot config get project.name
//
// Asserts that the second Get returns the value the Set wrote.
// Spec anchor: LH-FA-CONF-005 §1366-1380.
func TestConfigSet_LHFACONF005_RoundTripPin(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAML(t, fs)

	// Step 1: read current value.
	resp1, err := svc.Get(context.Background(), driving.ConfigGetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "project.name"),
	})
	if err != nil {
		t.Fatalf("Get #1: %v", err)
	}
	if resp1.Value != "t-uboot-config" {
		t.Errorf("Get #1 Value = %q, want %q", resp1.Value, "t-uboot-config")
	}

	// Step 2: set.
	if _, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "project.name"),
		Value:   "foo",
	}); err != nil {
		t.Fatalf("Set: %v", err)
	}

	// Step 3: read again and observe the new value.
	resp3, err := svc.Get(context.Background(), driving.ConfigGetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "project.name"),
	})
	if err != nil {
		t.Fatalf("Get #3: %v", err)
	}
	if resp3.Value != "foo" {
		t.Errorf("Get #3 Value = %q, want %q (LH-FA-CONF-005 roundtrip)", resp3.Value, "foo")
	}
}
