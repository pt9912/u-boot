package application_test

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

const manyFixture = `schemaVersion: 1
project:
  name: t-uboot-config
services:
  postgres:
    enabled: true
devcontainer:
  enabled: true
  profile: sandbox
  features:
    node:
      enabled: true
      version: "22"
`

func seedMany(t *testing.T, fs *fakeFS, body string) {
	t.Helper()
	seedConfigUbootYAMLWithDevcontainer(t, fs, body)
}

func readYAML(t *testing.T, fs *fakeFS) string {
	t.Helper()
	b, err := fs.ReadFile(filepath.Join(configTestBaseDir, "u-boot.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func paths(t *testing.T, raws ...string) []domain.ConfigPath {
	t.Helper()
	out := make([]domain.ConfigPath, 0, len(raws))
	for _, r := range raws {
		out = append(out, mustConfigPath(t, r))
	}
	return out
}

// slice-v1-config-multi-path-get: values in request order, all-or-nothing.
func TestConfigGetMany(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedMany(t, fs, manyFixture)
	resp, err := svc.GetMany(context.Background(), driving.ConfigGetManyRequest{BaseDir: configTestBaseDir,
		Paths: paths(t, "devcontainer.profile", "project.name", "services.postgres.enabled")})
	if err != nil {
		t.Fatalf("GetMany: %v", err)
	}
	want := []string{"sandbox", "t-uboot-config", "true"}
	if len(resp.Entries) != 3 {
		t.Fatalf("entries = %+v", resp.Entries)
	}
	for i, w := range want {
		if resp.Entries[i].Value != w {
			t.Errorf("entry %d = %q, want %q", i, resp.Entries[i].Value, w)
		}
	}
	// One unset path aborts the whole call, with its structured hint.
	_, err = svc.GetMany(context.Background(), driving.ConfigGetManyRequest{BaseDir: configTestBaseDir,
		Paths: paths(t, "project.name", "devcontainer.user.uid")})
	if !errors.Is(err, driving.ErrConfigValueNotSet) {
		t.Fatalf("err = %v, want ErrConfigValueNotSet", err)
	}
	if h := hintOf(t, err); h.Argument != "devcontainer.user.uid" {
		t.Errorf("hint = %+v", h)
	}
}

// slice-v1-config-multi-path-set: all pairs land in one write.
func TestConfigSetMany_AppliesAllInOneWrite(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedMany(t, fs, manyFixture)
	before := len(fs.writes)
	resp, err := svc.SetMany(context.Background(), driving.ConfigSetManyRequest{BaseDir: configTestBaseDir,
		Items: []driving.ConfigSetItem{
			{Path: mustConfigPath(t, "devcontainer.user.uid"), Value: "501"},
			{Path: mustConfigPath(t, "devcontainer.sandbox.nestedRuntime"), Value: "podman"},
			{Path: mustConfigPath(t, "devcontainer.sandbox.egress.allow"), Value: "a.example.com"},
		}})
	if err != nil {
		t.Fatalf("SetMany: %v", err)
	}
	if got := len(fs.writes) - before; got != 1 {
		t.Errorf("writes = %d, want exactly 1", got)
	}
	if len(resp.Entries) != 3 || resp.Entries[0].NewValue != "501" || resp.Entries[2].NewValue != "a.example.com" {
		t.Errorf("entries = %+v", resp.Entries)
	}
	y := readYAML(t, fs)
	for _, want := range []string{"uid: 501", "nestedRuntime: podman", "a.example.com"} {
		if !strings.Contains(y, want) {
			t.Errorf("u-boot.yaml lacks %q:\n%s", want, y)
		}
	}
}

// A failing item leaves the file byte-identical and writes nothing.
func TestConfigSetMany_AtomicOnFailure(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedMany(t, fs, manyFixture)
	before, writes := readYAML(t, fs), len(fs.writes)
	_, err := svc.SetMany(context.Background(), driving.ConfigSetManyRequest{BaseDir: configTestBaseDir,
		Items: []driving.ConfigSetItem{
			{Path: mustConfigPath(t, "devcontainer.user.uid"), Value: "501"},
			{Path: mustConfigPath(t, "devcontainer.profile"), Value: "paranoid"},
		}})
	if !errors.Is(err, driving.ErrConfigValueInvalid) {
		t.Fatalf("err = %v, want ErrConfigValueInvalid", err)
	}
	if readYAML(t, fs) != before || len(fs.writes) != writes {
		t.Errorf("u-boot.yaml changed on a failed SetMany")
	}
	// A write-rejected path aborts before any read/write too.
	_, err = svc.SetMany(context.Background(), driving.ConfigSetManyRequest{BaseDir: configTestBaseDir,
		Items: []driving.ConfigSetItem{
			{Path: mustConfigPath(t, "devcontainer.user.uid"), Value: "501"},
			{Path: mustConfigPath(t, "services.postgres.enabled"), Value: "false"},
		}})
	if !errors.Is(err, driving.ErrConfigWriteRejected) || readYAML(t, fs) != before {
		t.Errorf("write-rejected: err=%v, file changed=%v", err, readYAML(t, fs) != before)
	}
}

// Items apply in order on one document: the allowlist entry is
// available to the feature source set afterwards (LH-FA-DEV-003).
func TestConfigSetMany_OrderedDependency(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedMany(t, fs, manyFixture)
	items := []driving.ConfigSetItem{
		{Path: mustConfigPath(t, "devcontainer.features.node.source"), Value: "https://example.test/features/node"},
		{Path: mustConfigPath(t, "devcontainer.featureSources.allow"), Value: "https://example.test/features/node"},
	}
	// Wrong order: the source is not allowlisted yet → rejected, nothing written.
	before := readYAML(t, fs)
	if _, err := svc.SetMany(context.Background(), driving.ConfigSetManyRequest{BaseDir: configTestBaseDir, Items: items}); err == nil {
		t.Fatalf("source before allowlist must be rejected")
	}
	if readYAML(t, fs) != before {
		t.Errorf("file changed on rejected SetMany")
	}
	// Right order works.
	if _, err := svc.SetMany(context.Background(), driving.ConfigSetManyRequest{BaseDir: configTestBaseDir,
		Items: []driving.ConfigSetItem{items[1], items[0]}}); err != nil {
		t.Fatalf("ordered SetMany: %v", err)
	}
}

// All values already set → no write at all.
func TestConfigSetMany_NoOp(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedMany(t, fs, manyFixture)
	writes := len(fs.writes)
	resp, err := svc.SetMany(context.Background(), driving.ConfigSetManyRequest{BaseDir: configTestBaseDir,
		Items: []driving.ConfigSetItem{
			{Path: mustConfigPath(t, "devcontainer.profile"), Value: "sandbox"},
			{Path: mustConfigPath(t, "devcontainer.enabled"), Value: "true"},
		}})
	if err != nil || len(fs.writes) != writes {
		t.Fatalf("err=%v writes=%d→%d, want a pure no-op", err, writes, len(fs.writes))
	}
	for _, e := range resp.Entries {
		if e.OldValue != e.NewValue {
			t.Errorf("entry %v changed on no-op", e.Path)
		}
	}
}

// slice-v1-config-list-subcommand: sorted path/value pairs, only set paths.
func TestConfigList(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedMany(t, fs, manyFixture)
	resp, err := svc.List(context.Background(), driving.ConfigListRequest{BaseDir: configTestBaseDir})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	got := make([]string, 0, len(resp.Entries))
	for _, e := range resp.Entries {
		got = append(got, e.Path.String()+"="+e.Value)
	}
	want := []string{
		"devcontainer.enabled=true",
		"devcontainer.features.node.enabled=true",
		"devcontainer.features.node.version=22",
		"devcontainer.profile=sandbox",
		"project.name=t-uboot-config",
		"services.postgres.enabled=true",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("entries:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}
