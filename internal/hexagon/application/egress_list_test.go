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
