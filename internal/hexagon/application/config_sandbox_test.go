package application_test

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

const fixtureSandboxBase = `schemaVersion: 1
project:
  name: t-uboot-config
devcontainer:
  enabled: true
`

// LH-FA-DEV-004/-006/-007: the four Lastenheft 0.3.0 keys are
// settable, gettable and round-trip through the stage pipeline.
func TestConfigSetGet_SandboxKeys(t *testing.T) {
	t.Parallel()
	cases := []struct{ path, value string }{
		{"devcontainer.user.uid", "501"},
		{"devcontainer.profile", "sandbox"},
		{"devcontainer.sandbox.nestedRuntime", "podman"},
		{"devcontainer.sandbox.onUnavailable", "fail"},
		{"devcontainer.sandbox.repository", "git@github.com:other/fork.git"},
		{"devcontainer.sandbox.egress.enabled", "true"},
		{"devcontainer.sandbox.egress.allow", "api.anthropic.com,registry.npmjs.org"},
	}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()
			svc, fs := newConfigService(t)
			seedConfigUbootYAMLWithDevcontainer(t, fs, fixtureSandboxBase)
			p := mustConfigPath(t, tc.path)

			// Unset → ErrConfigValueNotSet (exit 10 class).
			if _, err := svc.Get(context.Background(), driving.ConfigGetRequest{BaseDir: configTestBaseDir, Path: p}); !errors.Is(err, driving.ErrConfigValueNotSet) {
				t.Fatalf("Get unset: err = %v, want ErrConfigValueNotSet", err)
			}
			resp, err := svc.Set(context.Background(), driving.ConfigSetRequest{BaseDir: configTestBaseDir, Path: p, Value: tc.value})
			if err != nil {
				t.Fatalf("Set: %v", err)
			}
			if resp.NewValue != tc.value {
				t.Errorf("NewValue = %q, want %q", resp.NewValue, tc.value)
			}
			got, err := svc.Get(context.Background(), driving.ConfigGetRequest{BaseDir: configTestBaseDir, Path: p})
			if err != nil {
				t.Fatalf("Get: %v", err)
			}
			if got.Value != tc.value {
				t.Errorf("Get = %q, want %q", got.Value, tc.value)
			}
			// Second identical Set is a NoOp.
			resp2, err := svc.Set(context.Background(), driving.ConfigSetRequest{BaseDir: configTestBaseDir, Path: p, Value: tc.value})
			if err != nil || resp2.OldValue != resp2.NewValue {
				t.Errorf("second Set: err=%v old=%q new=%q, want NoOp", err, resp2.OldValue, resp2.NewValue)
			}
		})
	}
}

// Invalid values → ErrConfigValueInvalid (exit 10); the file stays
// untouched.
func TestConfigSet_SandboxKeys_InvalidValue(t *testing.T) {
	t.Parallel()
	cases := []struct{ path, value string }{
		{"devcontainer.user.uid", "0"},
		{"devcontainer.user.uid", "65536"},
		{"devcontainer.user.uid", "abc"},
		{"devcontainer.profile", "strict"},
		{"devcontainer.sandbox.nestedRuntime", "docker"},
		{"devcontainer.sandbox.onUnavailable", "ignore"},
		{"devcontainer.sandbox.repository", "https://token@github.com/o/r.git"},
		{"devcontainer.sandbox.repository", "https://github.com/o/r.git;rm -rf /"},
		{"devcontainer.sandbox.egress.enabled", "maybe"},
		{"devcontainer.sandbox.egress.allow", "https://api.anthropic.com"},
		{"devcontainer.sandbox.egress.allow", "*.github.com"},
		{"devcontainer.sandbox.egress.allow", "Example.COM"},
		{"devcontainer.sandbox.egress.allow", "localhost"},
	}
	for _, tc := range cases {
		t.Run(tc.path+"="+tc.value, func(t *testing.T) {
			t.Parallel()
			svc, fs := newConfigService(t)
			seedConfigUbootYAMLWithDevcontainer(t, fs, fixtureSandboxBase)
			_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
				BaseDir: configTestBaseDir, Path: mustConfigPath(t, tc.path), Value: tc.value})
			if !errors.Is(err, driving.ErrConfigValueInvalid) {
				t.Fatalf("err = %v, want ErrConfigValueInvalid", err)
			}
			body, _ := fs.ReadFile(filepath.Join(configTestBaseDir, "u-boot.yaml"))
			if string(body) != fixtureSandboxBase {
				t.Errorf("u-boot.yaml changed on rejected Set:\n%s", body)
			}
		})
	}
}

// A hand-edited invalid value is rejected on load (schema invalid,
// exit 10) instead of being silently accepted.
func TestConfigGet_SandboxKeys_InvalidOnLoad(t *testing.T) {
	t.Parallel()
	bodies := map[string]string{
		"uid zero":  "devcontainer:\n  enabled: true\n  user:\n    uid: 0\n",
		"profile":   "devcontainer:\n  enabled: true\n  profile: paranoid\n",
		"nested":    "devcontainer:\n  enabled: true\n  sandbox:\n    nestedRuntime: lxc\n",
		"onUnavail": "devcontainer:\n  enabled: true\n  sandbox:\n    onUnavailable: maybe\n",
	}
	for name, dc := range bodies {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newConfigService(t)
			seedConfigUbootYAMLWithDevcontainer(t, fs, "schemaVersion: 1\nproject:\n  name: t-uboot-config\n"+dc)
			_, err := svc.Get(context.Background(), driving.ConfigGetRequest{
				BaseDir: configTestBaseDir, Path: mustConfigPath(t, "devcontainer.enabled")})
			if !errors.Is(err, driving.ErrConfigSchemaInvalid) {
				t.Fatalf("err = %v, want ErrConfigSchemaInvalid", err)
			}
			if !strings.Contains(err.Error(), "devcontainer") {
				t.Errorf("error should name the devcontainer subtree: %v", err)
			}
		})
	}
}

// egress.allow appends and de-duplicates (list path); the load
// validator rejects hand-edited bad hosts.
func TestConfigSet_EgressAllow_AppendDedupe(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAMLWithDevcontainer(t, fs, fixtureSandboxBase)
	p := mustConfigPath(t, "devcontainer.sandbox.egress.allow")
	set := func(v string) driving.ConfigSetResponse {
		t.Helper()
		resp, err := svc.Set(context.Background(), driving.ConfigSetRequest{BaseDir: configTestBaseDir, Path: p, Value: v})
		if err != nil {
			t.Fatalf("Set(%q): %v", v, err)
		}
		return resp
	}
	if r := set("a.example.com,b.example.com"); r.NewValue != "a.example.com,b.example.com" {
		t.Errorf("first NewValue = %q", r.NewValue)
	}
	if r := set("b.example.com,c.example.com"); r.NewValue != "a.example.com,b.example.com,c.example.com" {
		t.Errorf("merged NewValue = %q", r.NewValue)
	}
	if r := set("a.example.com"); r.OldValue != r.NewValue {
		t.Errorf("re-setting an existing host must be a NoOp: %+v", r)
	}
	got, err := svc.Get(context.Background(), driving.ConfigGetRequest{BaseDir: configTestBaseDir, Path: p})
	if err != nil || got.Value != "a.example.com,b.example.com,c.example.com" {
		t.Errorf("Get = %q, %v", got.Value, err)
	}
}

func TestConfigGet_EgressAllow_InvalidOnLoad(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAMLWithDevcontainer(t, fs, fixtureSandboxBase+"  sandbox:\n    egress:\n      allow:\n        - \"https://bad.example\"\n")
	_, err := svc.Get(context.Background(), driving.ConfigGetRequest{BaseDir: configTestBaseDir, Path: mustConfigPath(t, "devcontainer.enabled")})
	if !errors.Is(err, driving.ErrConfigSchemaInvalid) {
		t.Fatalf("err = %v, want ErrConfigSchemaInvalid", err)
	}
}
