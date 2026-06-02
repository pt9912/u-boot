package application_test

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// TestInit_AllowExternalFeatureSources_Seeds pins the LH-FA-DEV-003
// Spec §714 init wiring: passing `--allow-external-feature-sources
// URL[,URL]` together with `--devcontainer` seeds the freshly-
// written u-boot.yaml's `devcontainer.featureSources.allow` list.
func TestInit_AllowExternalFeatureSources_Seeds(t *testing.T) {
	svc, fs, _, _ := newService(t)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:                        "demo",
		BaseDir:                     testBaseDir,
		SkipGit:                     true,
		Devcontainer:                true,
		AllowExternalFeatureSources: []string{"https://example.test/a", "https://example.test/b"},
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	body, err := fs.ReadFile(filepath.Join(testBaseDir, "u-boot.yaml"))
	if err != nil {
		t.Fatalf("read u-boot.yaml: %v", err)
	}
	for _, want := range []string{"https://example.test/a", "https://example.test/b"} {
		if !strings.Contains(string(body), want) {
			t.Errorf("u-boot.yaml missing %q\nbody:\n%s", want, body)
		}
	}
	if !strings.Contains(string(body), "featureSources:") {
		t.Errorf("u-boot.yaml missing featureSources block\nbody:\n%s", body)
	}
}

// TestInit_AllowExternalFeatureSources_RequiresDevcontainer pins the
// Spec §714 constraint: the flag is only valid together with
// `--devcontainer`. Without it the use case rejects before any FS
// side effect with the LH-FA-DEV-003 sentinel.
func TestInit_AllowExternalFeatureSources_RequiresDevcontainer(t *testing.T) {
	svc, _, _, _ := newService(t)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:                        "demo",
		BaseDir:                     testBaseDir,
		SkipGit:                     true,
		Devcontainer:                false,
		AllowExternalFeatureSources: []string{"https://example.test/a"},
	})
	if err == nil {
		t.Fatalf("Init: expected error, got nil")
	}
	if !errors.Is(err, application.ErrInvalidFeatureSource) {
		t.Errorf("err = %v, want wrap of ErrInvalidFeatureSource", err)
	}
}

// TestInit_AllowExternalFeatureSources_InvalidURL pins that a bad
// URL on the init flag rejects with the LH-FA-DEV-003 sentinel
// (validateFeatureSource catches malformed entries before marshal).
func TestInit_AllowExternalFeatureSources_InvalidURL(t *testing.T) {
	svc, _, _, _ := newService(t)

	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name:                        "demo",
		BaseDir:                     testBaseDir,
		SkipGit:                     true,
		Devcontainer:                true,
		AllowExternalFeatureSources: []string{"not-a-url"},
	})
	if err == nil {
		t.Fatalf("Init: expected error, got nil")
	}
	if !errors.Is(err, application.ErrInvalidFeatureSource) {
		t.Errorf("err = %v, want wrap of ErrInvalidFeatureSource", err)
	}
}
