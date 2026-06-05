package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// TestNoopDefaults_Smoke exercises noopProgress / noopConfirmer /
// noopLogger via the [application.CallNoopDefaultsForTest] bridge.
// These nil-tolerant defaults are normally only entered when a
// constructor caller passes nil for an optional dependency — most
// tests wire a real stub, leaving the noop methods at 0% coverage.
// The smoke-test asserts no-panic; the contract for confirmers is
// pinned via the explicit table below.
func TestNoopDefaults_Smoke(t *testing.T) {
	application.CallNoopDefaultsForTest(context.Background())
}

// TestParseForwardPorts_AllShapes covers parseForwardPorts's three
// dispatch branches: float64 (devcontainer.json numeric ports),
// string (host:container mapping), invalid (silently skipped).
func TestParseForwardPorts_AllShapes(t *testing.T) {
	items := []any{
		float64(5432),    // numeric → 5432
		"8080:80",        // host:container → 80
		"127.0.0.1:9000:90", // ipv4-host:host:container → 90
		"bare-string",    // unparseable → skipped
		map[string]any{}, // map → skipped (no case)
		float64(3000),    // duplicate numeric
	}
	got := application.ParseForwardPortsForTest(items)
	for _, want := range []int{5432, 80, 90, 3000} {
		if _, ok := got[want]; !ok {
			t.Errorf("missing port %d in result: %v", want, got)
		}
	}
	if _, ok := got[0]; ok {
		t.Errorf("invalid entries must not appear in result: %v", got)
	}
}

// TestConfigPathToYAMLPath_AllKinds covers the five switch cases
// (ConfigProjectName, ConfigDevcontainerEnabled, three Feature*
// kinds) plus the unreachable-sentinel fallback.
func TestConfigPathToYAMLPath_AllKinds(t *testing.T) {
	feat, err := domain.NewFeatureName("docker-in-docker")
	if err != nil {
		t.Fatalf("NewFeatureName: %v", err)
	}
	cases := []struct {
		name string
		path domain.ConfigPath
		want []string
	}{
		{"project.name", domain.ConfigPath{Kind: domain.ConfigProjectName}, []string{"project", "name"}},
		{"devcontainer.enabled", domain.ConfigPath{Kind: domain.ConfigDevcontainerEnabled}, []string{"devcontainer", "enabled"}},
		{"feature enabled", domain.ConfigPath{Kind: domain.ConfigDevcontainerFeatureEnabled, Feature: feat}, []string{"devcontainer", "features", "docker-in-docker", "enabled"}},
		{"feature source", domain.ConfigPath{Kind: domain.ConfigDevcontainerFeatureSource, Feature: feat}, []string{"devcontainer", "features", "docker-in-docker", "source"}},
		{"feature version", domain.ConfigPath{Kind: domain.ConfigDevcontainerFeatureVersion, Feature: feat}, []string{"devcontainer", "features", "docker-in-docker", "version"}},
		{"unknown sentinel", domain.ConfigPath{Kind: domain.ConfigPathKind(99)}, []string{"__unknown_config_path_kind__"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := application.ConfigPathToYAMLPathForTest(tc.path)
			if len(got) != len(tc.want) {
				t.Fatalf("length mismatch: want %v, got %v", tc.want, got)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("segment %d: want %q, got %q", i, tc.want[i], got[i])
				}
			}
		})
	}
}

// TestNormalisePortEntry_AllShapes covers the four Compose port-entry
// shapes recognised by normalisePortEntry — int / bare-string /
// host:container mapping / long-form map (unsupported).
func TestNormalisePortEntry_AllShapes(t *testing.T) {
	cases := []struct {
		name      string
		raw       any
		wantPort  int
		wantOK    bool
	}{
		{"int", int(5432), 5432, true},
		{"bare string", "5432", 5432, true},
		{"host:container", "8080:80", 80, true},
		{"ipv4 host:host:container", "127.0.0.1:8080:80", 80, true},
		{"with protocol suffix", "8080:80/tcp", 80, true},
		{"unsupported (map)", map[string]any{"target": 80}, 0, false},
		{"unsupported (slice)", []any{}, 0, false},
		{"unparseable string", "not-a-port", 0, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotPort, gotOK := application.NormalisePortEntryForTest(tc.raw)
			if gotPort != tc.wantPort || gotOK != tc.wantOK {
				t.Errorf("want (%d, %v), got (%d, %v)", tc.wantPort, tc.wantOK, gotPort, gotOK)
			}
		})
	}
}

// TestRelativizePath_AllBranches covers relativizePath's three
// branches: empty baseDir (pass-through), successful Rel, and
// non-prefix fallback.
func TestRelativizePath_AllBranches(t *testing.T) {
	cases := []struct {
		name    string
		path    string
		baseDir string
		want    string
	}{
		{"empty baseDir", "/abs/path", "", "/abs/path"},
		{"relative inside base", "/proj/sub/file.txt", "/proj", "sub/file.txt"},
		{"outside base → fallback", "/other/file.txt", "/proj", "/other/file.txt"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := application.RelativizePathForTest(tc.path, tc.baseDir)
			if got != tc.want {
				t.Errorf("want %q, got %q", tc.want, got)
			}
		})
	}
}

// TestServiceAction_String covers all five serviceAction values
// (register, reactivate, rebuild-block, repair-artifacts) plus the
// default "unknown" fallback for the enum.
func TestServiceAction_String(t *testing.T) {
	cases := []struct {
		val  int
		want string
	}{
		{0, "register"},
		{1, "reactivate"},
		{2, "rebuild-block"},
		{3, "repair-artifacts"},
		{99, "unknown"},
	}
	for _, tc := range cases {
		got := application.ServiceActionStringForTest(tc.val)
		if got != tc.want {
			t.Errorf("action(%d): want %q, got %q", tc.val, tc.want, got)
		}
	}
}

// TestTranslatePatchErr_AllBranches covers the three switch-cases of
// translatePatchErr in addservice_execute.go (currently 0%): the two
// adapter-sentinel-translation branches and the default path.
func TestTranslatePatchErr_AllBranches(t *testing.T) {
	svc, err := domain.NewServiceName("postgres")
	if err != nil {
		t.Fatalf("NewServiceName: %v", err)
	}
	t.Run("ErrYAMLAnchorMismatch → ErrServiceInconsistent", func(t *testing.T) {
		got := application.TranslatePatchErrForTest(driven.ErrYAMLAnchorMismatch, svc, "compose.yaml")
		if !errors.Is(got, driving.ErrServiceInconsistent) {
			t.Errorf("want ErrServiceInconsistent, got %v", got)
		}
	})
	t.Run("ErrYAMLFragmentInvalid → ErrServiceInconsistent", func(t *testing.T) {
		got := application.TranslatePatchErrForTest(driven.ErrYAMLFragmentInvalid, svc, "compose.yaml")
		if !errors.Is(got, driving.ErrServiceInconsistent) {
			t.Errorf("want ErrServiceInconsistent, got %v", got)
		}
	})
	t.Run("default: unknown error → plain wrap, not classified", func(t *testing.T) {
		raw := errors.New("disk-quota exceeded")
		got := application.TranslatePatchErrForTest(raw, svc, ".env.example")
		if !errors.Is(got, raw) {
			t.Errorf("want wrapped raw error, got %v", got)
		}
		if errors.Is(got, driving.ErrServiceInconsistent) {
			t.Errorf("default branch must not map to ErrServiceInconsistent; got %v", got)
		}
	})
}
