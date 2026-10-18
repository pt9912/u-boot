package application_test

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

func egressListErr(fs *fakeFS) error {
	_, err := newEgressListService(fs).List(context.Background(), driving.EgressListRequest{BaseDir: generateTestBaseDir})
	return err
}

// LH-FA-DEV-008, LH-FA-CLI-006: without a project, with an unreadable or an
// invalid u-boot.yaml, and with an `origin` URL that `generate devcontainer
// --sandbox` refuses, the listing fails with the sentinel of the neighbours;
// the CLI maps them to exit 10 (exit 14 for the read failure).
func TestEgressList_Errors(t *testing.T) {
	t.Parallel()
	t.Run("no u-boot.yaml", func(t *testing.T) {
		t.Parallel()
		_, fs := newGenerateService(t)
		if err := egressListErr(fs); !errors.Is(err, driving.ErrProjectNotInitialized) {
			t.Fatalf("err = %v, want ErrProjectNotInitialized", err)
		}
	})
	t.Run("unparsable u-boot.yaml", func(t *testing.T) {
		t.Parallel()
		_, fs := newGenerateService(t)
		if err := fs.WriteFile(filepath.Join(generateTestBaseDir, "u-boot.yaml"), []byte("project: [unclosed\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := egressListErr(fs); !errors.Is(err, driving.ErrConfigSchemaInvalid) {
			t.Fatalf("err = %v, want ErrConfigSchemaInvalid", err)
		}
	})
	t.Run("invalid egress entry", func(t *testing.T) {
		t.Parallel()
		_, fs := newGenerateService(t)
		seedUBootYAMLWithFeatures(t, fs, egressConfig("", "", "*.example.org"))
		if err := egressListErr(fs); !errors.Is(err, driving.ErrConfigSchemaInvalid) {
			t.Fatalf("err = %v, want ErrConfigSchemaInvalid", err)
		}
	})
	t.Run("read failure", func(t *testing.T) {
		t.Parallel()
		_, fs := newGenerateService(t)
		seedUBootYAMLWithFeatures(t, fs, egressConfig("", ""))
		fs.failReadOn = filepath.Join(generateTestBaseDir, "u-boot.yaml")
		fs.failReadErr = errors.New("disk on fire")
		err := egressListErr(fs)
		if !errors.Is(err, driving.ErrEgressFileSystem) || errors.Is(err, driving.ErrProjectNotInitialized) {
			t.Fatalf("err = %v, want ErrEgressFileSystem only", err)
		}
	})
	t.Run("probe failure", func(t *testing.T) {
		t.Parallel()
		_, fs := newGenerateService(t)
		fs.failExistsOn = filepath.Join(generateTestBaseDir, "u-boot.yaml")
		fs.failExistsErr = errors.New("stat failed")
		if err := egressListErr(fs); !errors.Is(err, driving.ErrEgressFileSystem) {
			t.Fatalf("err = %v, want ErrEgressFileSystem", err)
		}
	})
	t.Run("empty base dir", func(t *testing.T) {
		t.Parallel()
		_, fs := newGenerateService(t)
		_, err := newEgressListService(fs).List(context.Background(), driving.EgressListRequest{})
		if err == nil || errors.Is(err, driving.ErrProjectNotInitialized) {
			t.Fatalf("err = %v, want a non-sentinel error", err)
		}
	})
	t.Run("origin with credentials is refused like generate", func(t *testing.T) {
		t.Parallel()
		svc, fs := newGenerateService(t)
		seedUBootYAMLWithFeatures(t, fs, egressConfig("", ""))
		seedGitOrigin(t, fs, "https://user:secret@example.org/o/r.git")
		if err := egressListErr(fs); !errors.Is(err, domain.ErrInvalidSandboxSetting) {
			t.Fatalf("list err = %v, want ErrInvalidSandboxSetting", err)
		}
		if _, err := generateSandbox(svc, false); !errors.Is(err, domain.ErrInvalidSandboxSetting) {
			t.Fatalf("generate err = %v, want the same sentinel (the two commands agree)", err)
		}
	})
}

// LH-FA-DEV-008: no remote and no `sandbox.repository` is not an error: the
// list has no clone host and says so (`CloneSource` false), where `generate`
// warns and writes no clone step.
func TestEgressList_NoCloneSourceIsNotAnError(t *testing.T) {
	t.Parallel()
	_, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, egressConfig("", ""))
	resp := listEgress(t, fs)
	if resp.CloneSource || originsContain(resp, "clone") {
		t.Errorf("no remote: CloneSource=%v, clone origin=%v, want neither", resp.CloneSource, originsContain(resp, "clone"))
	}
}

// LH-FA-DEV-008: the clone host is part of the list that would apply on
// activation, so it is resolved while the restriction is not active — without
// the sandbox profile, with `egress.enabled` off or unset, and from `origin` as
// from `sandbox.repository` (the resolution of `generate devcontainer --sandbox`).
func TestEgressList_CloneHostIsResolvedWithoutActiveRestriction(t *testing.T) {
	t.Parallel()
	const repository = "    repository: https://git.example.net/o/r.git\n"
	const origin = "git@gitlab.example.org:o/r.git"
	type cloneCase struct{ yaml, origin, host string }
	cases := map[string]cloneCase{
		"no sandbox profile, egress on, origin": {
			"devcontainer:\n  enabled: true\n  sandbox:\n    egress:\n      enabled: true\n", origin, "gitlab.example.org"},
		"sandbox profile, egress off, origin": {
			"devcontainer:\n  enabled: true\n  profile: sandbox\n  sandbox:\n    egress:\n      enabled: false\n", origin, "gitlab.example.org"},
		"sandbox profile, egress not set, origin": {
			"devcontainer:\n  enabled: true\n  profile: sandbox\n", origin, "gitlab.example.org"},
		"no sandbox profile, repository wins over origin": {
			"devcontainer:\n  enabled: true\n  sandbox:\n" + repository, origin, "git.example.net"},
		"no devcontainer section, origin": {"", origin, "gitlab.example.org"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, tc.yaml)
			seedGitOrigin(t, fs, tc.origin)
			resp := listEgress(t, fs)
			if resp.Active() {
				t.Fatalf("the case must not be active: profile=%v enabled=%v", resp.SandboxProfile, resp.EgressEnabled)
			}
			if !resp.CloneSource {
				t.Errorf("CloneSource = false, want true")
			}
			if got := egressOrigins(resp)[tc.host]; !reflect.DeepEqual(got, []string{"clone"}) {
				t.Errorf("origins of %s = %v, want [clone]", tc.host, got)
			}
		})
	}
}

// LH-FA-DEV-008, LH-FA-DEV-006: an `origin` URL that `generate devcontainer
// --sandbox` refuses (credentials) fails the listing also while the restriction
// is not active; the sentinel is the one the CLI classifies as exit 10 with the
// diagnostic LH-FA-DEV-006. A `sandbox.repository` with credentials is refused
// by the validation of the configuration (exit 10, LH-FA-CONF-002), whatever
// the profile.
func TestEgressList_UnsafeCloneUrlIsRefusedWithoutActiveRestriction(t *testing.T) {
	t.Parallel()
	const credentials = "https://user:secret@example.org/o/r.git"
	type unsafeCase struct {
		yaml, origin string
		want         error
	}
	cases := map[string]unsafeCase{
		"origin, no sandbox profile": {
			"devcontainer:\n  enabled: true\n", credentials, domain.ErrInvalidSandboxSetting},
		"origin, sandbox profile, egress off": {
			"devcontainer:\n  enabled: true\n  profile: sandbox\n  sandbox:\n    egress:\n      enabled: false\n", credentials, domain.ErrInvalidSandboxSetting},
		"repository, no sandbox profile": {
			"devcontainer:\n  enabled: true\n  sandbox:\n    repository: " + credentials + "\n", "", driving.ErrConfigSchemaInvalid},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, tc.yaml)
			if tc.origin != "" {
				seedGitOrigin(t, fs, tc.origin)
			}
			if err := egressListErr(fs); !errors.Is(err, tc.want) {
				t.Fatalf("err = %v, want %v", err, tc.want)
			}
		})
	}
}

// A nil logger is routed to the no-op logger.
func TestNewEgressListService_NilLogger(t *testing.T) {
	t.Parallel()
	_, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, egressConfig("", ""))
	if _, err := application.NewEgressListService(fs, &fakeYAML{}, nil).List(context.Background(),
		driving.EgressListRequest{BaseDir: generateTestBaseDir}); err != nil {
		t.Fatalf("List with nil logger: %v", err)
	}
}
