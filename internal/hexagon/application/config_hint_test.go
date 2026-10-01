package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

func hintOf(t *testing.T, err error) driving.ConfigHint {
	t.Helper()
	var he *driving.ConfigHintError
	if !errors.As(err, &he) {
		t.Fatalf("err = %v, want a *ConfigHintError in the chain", err)
	}
	return he.Hint
}

// slice-v1-config-structured-hint: a write-rejected service path
// carries `u-boot add <svc>` as a structured hint; the sentinel chain
// stays intact.
func TestConfigSet_WriteRejected_StructuredHint(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	seedConfigUbootYAML(t, fs)
	_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir, Path: mustConfigPath(t, "services.postgres.enabled"), Value: "true"})
	if !errors.Is(err, driving.ErrConfigWriteRejected) {
		t.Fatalf("err = %v, want ErrConfigWriteRejected", err)
	}
	h := hintOf(t, err)
	if h.Command != "u-boot add postgres" || h.Action != "add" || h.Argument != "postgres" {
		t.Errorf("hint = %+v", h)
	}
}

// ErrConfigValueNotSet hints: services → add; devcontainer.enabled
// and every other path → config set.
func TestConfigGet_NotSet_StructuredHint(t *testing.T) {
	t.Parallel()
	cases := []struct {
		path string
		want driving.ConfigHint
	}{
		{"services.postgres.enabled", driving.ConfigHint{Command: "u-boot add postgres", Action: "add", Argument: "postgres"}},
		{"devcontainer.enabled", driving.ConfigHint{Command: "u-boot config set devcontainer.enabled <true|false>", Action: "config-set", Argument: "devcontainer.enabled"}},
		{"devcontainer.profile", driving.ConfigHint{Command: "u-boot config set devcontainer.profile <value>", Action: "config-set", Argument: "devcontainer.profile"}},
		{"devcontainer.features.node.version", driving.ConfigHint{Command: "u-boot config set devcontainer.features.node.version <value>", Action: "config-set", Argument: "devcontainer.features.node.version"}},
	}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()
			svc, fs := newConfigService(t)
			seedConfigUbootYAML(t, fs)
			_, err := svc.Get(context.Background(), driving.ConfigGetRequest{BaseDir: configTestBaseDir, Path: mustConfigPath(t, tc.path)})
			if !errors.Is(err, driving.ErrConfigValueNotSet) {
				t.Fatalf("err = %v, want ErrConfigValueNotSet", err)
			}
			if got := hintOf(t, err); got != tc.want {
				t.Errorf("hint = %+v, want %+v", got, tc.want)
			}
		})
	}
}
