package application_test

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
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

// LH-FA-DEV-006: `init --devcontainer --sandbox` renders the sandbox
// shape and persists the profile; without a remote it warns.
func TestInit_Sandbox_RendersProfileAndWarnsWithoutRemote(t *testing.T) {
	svc, fs, _, _ := newService(t)

	resp, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name: "demo", BaseDir: testBaseDir, SkipGit: true, Devcontainer: true, Sandbox: true,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	if len(resp.Warnings) != 1 || resp.Warnings[0].Code != "LH-FA-DEV-006" {
		t.Errorf("warnings = %+v, want one LH-FA-DEV-006", resp.Warnings)
	}
	dc, err := fs.ReadFile(filepath.Join(testBaseDir, ".devcontainer", "devcontainer.json"))
	if err != nil {
		t.Fatalf("read devcontainer.json: %v", err)
	}
	if !strings.Contains(string(dc), `"workspaceMount": "source=demo-workspace,target=/workspaces/demo,type=volume"`) {
		t.Errorf("devcontainer.json lacks sandbox workspaceMount:\n%s", dc)
	}
	if strings.Contains(string(dc), "postCreateCommand") {
		t.Errorf("no remote: postCreateCommand must be absent:\n%s", dc)
	}
	y, _ := fs.ReadFile(filepath.Join(testBaseDir, "u-boot.yaml"))
	if !strings.Contains(string(y), "profile: sandbox") {
		t.Errorf("u-boot.yaml lacks profile: sandbox:\n%s", y)
	}
}

// Default init (no --sandbox) keeps the pre-0.3.0 devcontainer.json.
func TestInit_Devcontainer_NoSandbox_Unchanged(t *testing.T) {
	svc, fs, _, _ := newService(t)
	resp, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name: "demo", BaseDir: testBaseDir, SkipGit: true, Devcontainer: true,
	})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	if len(resp.Warnings) != 0 {
		t.Errorf("unexpected warnings: %+v", resp.Warnings)
	}
	dc, _ := fs.ReadFile(filepath.Join(testBaseDir, ".devcontainer", "devcontainer.json"))
	for _, s := range []string{"workspaceMount", "postCreateCommand", "USER_UID"} {
		if strings.Contains(string(dc), s) {
			t.Errorf("default devcontainer.json must not contain %q", s)
		}
	}
}

// `--sandbox` without `--devcontainer` → domain sentinel (exit 10)
// before any write.
func TestInit_Sandbox_RequiresDevcontainer(t *testing.T) {
	svc, fs, _, _ := newService(t)
	_, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name: "demo", BaseDir: testBaseDir, SkipGit: true, Sandbox: true,
	})
	if !errors.Is(err, domain.ErrInvalidSandboxSetting) {
		t.Fatalf("err = %v, want ErrInvalidSandboxSetting", err)
	}
	if exists, _ := fs.Exists(filepath.Join(testBaseDir, "u-boot.yaml")); exists {
		t.Errorf("u-boot.yaml must not be written")
	}
}
