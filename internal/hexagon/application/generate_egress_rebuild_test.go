package application_test

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// LH-FA-DEV-008: the notice that an existing container keeps the egress
// list of its image appears when `generate devcontainer` replaces the
// managed block of an existing egress-init.sh, and only then.

const (
	egressRebuildSubject = ".devcontainer/egress-init.sh"
	egressRebuildAllowA  = "devcontainer:\n  enabled: true\n  profile: sandbox\n  sandbox:\n    egress:\n      enabled: true\n      allow:\n        - a.example.com\n"
	egressRebuildAllowB  = "devcontainer:\n  enabled: true\n  profile: sandbox\n  sandbox:\n    egress:\n      enabled: true\n      allow:\n        - a.example.com\n        - b.example.com\n"
)

// egressRebuildNotices returns the warnings that carry the rebuild notice.
func egressRebuildNotices(ws []driving.WarningEntry) []driving.WarningEntry {
	var out []driving.WarningEntry
	for _, w := range ws {
		if w.Subject == egressRebuildSubject {
			out = append(out, w)
		}
	}
	return out
}

func generateDevcontainerMode(svc *application.GenerateService, mode driving.PreviewMode) (driving.GenerateResponse, error) {
	return svc.Generate(context.Background(), driving.GenerateRequest{
		BaseDir: generateTestBaseDir, Artifact: domain.ArtifactDevcontainer, PreviewMode: mode})
}

func TestGenerateDevcontainer_EgressRebuildNotice(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)

	// Fresh state: the script is written for the first time, no notice.
	seedUBootYAMLWithFeatures(t, fs, egressRebuildAllowA)
	resp, err := generateSandbox(svc, false)
	if err != nil {
		t.Fatalf("fresh generate: %v", err)
	}
	if n := egressRebuildNotices(resp.Warnings); len(n) != 0 {
		t.Errorf("fresh state gave the rebuild notice: %+v", n)
	}

	// Same configuration again: no-op, no notice.
	resp, err = generateSandbox(svc, false)
	if err != nil {
		t.Fatalf("second generate: %v", err)
	}
	if resp.Action != driving.GenerateActionNoOp || len(egressRebuildNotices(resp.Warnings)) != 0 {
		t.Errorf("unchanged run: action=%v notices=%+v, want NoOp and none", resp.Action, egressRebuildNotices(resp.Warnings))
	}

	// A host added to the list: the script changes, the notice appears.
	seedUBootYAMLWithFeatures(t, fs, egressRebuildAllowB)
	resp, err = generateSandbox(svc, false)
	if err != nil {
		t.Fatalf("changed generate: %v", err)
	}
	n := egressRebuildNotices(resp.Warnings)
	if len(n) != 1 {
		t.Fatalf("rebuild notices = %d, want 1: %+v", len(n), resp.Warnings)
	}
	w := n[0]
	if w.Code != "LH-FA-DEV-008" || w.Level != "warn" || w.SecurityRelaxation || w.DegradationFallback {
		t.Errorf("notice = %+v, want code LH-FA-DEV-008, level warn, neither flag set (--quiet drops it)", w)
	}
	for _, want := range []string{"keeps the egress list of the image it was built from", "rebuild the container", "devcontainer up --remove-existing-container"} {
		if !strings.Contains(w.Message, want) {
			t.Errorf("notice %q lacks %q", w.Message, want)
		}
	}
}

// A change that leaves egress-init.sh alone gives no notice: the project
// file `devcontainer.json` changes (a feature is added), the script does not.
func TestGenerateDevcontainer_EgressRebuildNotice_NotForOtherFiles(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, egressRebuildAllowA)
	if _, err := generateSandbox(svc, false); err != nil {
		t.Fatalf("fresh generate: %v", err)
	}
	seedUBootYAMLWithFeatures(t, fs, egressRebuildAllowA+"  features:\n    sshd:\n      enabled: true\n")
	resp, err := generateSandbox(svc, false)
	if err != nil {
		t.Fatalf("changed generate: %v", err)
	}
	if resp.Action != driving.GenerateActionUpdatedBlock {
		t.Fatalf("action = %v, want UpdatedBlock (devcontainer.json changed)", resp.Action)
	}
	if n := egressRebuildNotices(resp.Warnings); len(n) != 0 {
		t.Errorf("a change without egress-init.sh gave the rebuild notice: %+v", n)
	}
}

// `--dry-run` writes nothing, so it gives no notice, and it leaves the
// script untouched.
func TestGenerateDevcontainer_EgressRebuildNotice_DryRunGivesNone(t *testing.T) {
	t.Parallel()
	prod := newFakeFS()
	prod.markDirExists(generateTestBaseDir)
	yamlPath := filepath.Join(generateTestBaseDir, "u-boot.yaml")
	seed := func(body string) {
		t.Helper()
		if err := prod.WriteFile(yamlPath, []byte("schemaVersion: 1\nproject:\n  name: t-uboot-gen\n"+body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	svc := application.NewGenerateServiceWithFactory(configFactoryForProd(prod), &fakeYAML{}, nil)
	seed(egressRebuildAllowA)
	if _, err := generateDevcontainerMode(svc, driving.PreviewNone); err != nil {
		t.Fatalf("fresh generate: %v", err)
	}
	before, _ := prod.ReadFile(egressInitPath())
	seed(egressRebuildAllowB)

	resp, err := generateDevcontainerMode(svc, driving.PreviewDryRun)
	if err != nil {
		t.Fatalf("dry-run: %v", err)
	}
	if n := egressRebuildNotices(resp.Warnings); len(n) != 0 {
		t.Errorf("--dry-run gave the rebuild notice: %+v", n)
	}
	if after, _ := prod.ReadFile(egressInitPath()); string(after) != string(before) {
		t.Errorf("--dry-run changed egress-init.sh")
	}

	resp, err = generateDevcontainerMode(svc, driving.PreviewNone)
	if err != nil {
		t.Fatalf("real run: %v", err)
	}
	if n := egressRebuildNotices(resp.Warnings); len(n) != 1 {
		t.Errorf("the real run after the dry-run gave %d notices, want 1", len(n))
	}
}

// Activating and switching off the restriction give no notice: the first
// activation creates egress-init.sh and replaces the managed blocks of
// devcontainer.json and Dockerfile; switching off replaces those two blocks
// and leaves egress-init.sh in place (it is no longer planned).
func TestGenerateDevcontainer_EgressRebuildNotice_NotForActivationAndSwitchOff(t *testing.T) {
	t.Parallel()
	const off = "devcontainer:\n  enabled: true\n  profile: sandbox\n  sandbox:\n    egress:\n      enabled: false\n"
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, off)
	if _, err := generateSandbox(svc, false); err != nil {
		t.Fatalf("fresh generate: %v", err)
	}
	if ok, _ := fs.Exists(egressInitPath()); ok {
		t.Fatalf("egress-init.sh exists while the restriction is off")
	}
	jsonBefore, _ := fs.ReadFile(devcontainerJSONPath())
	dockerBefore, _ := fs.ReadFile(dockerfilePath())

	seedUBootYAMLWithFeatures(t, fs, egressRebuildAllowA)
	resp, err := generateSandbox(svc, false)
	if err != nil {
		t.Fatalf("activation: %v", err)
	}
	if n := egressRebuildNotices(resp.Warnings); len(n) != 0 {
		t.Errorf("first activation gave the rebuild notice: %+v", n)
	}
	if ok, _ := fs.Exists(egressInitPath()); !ok {
		t.Errorf("activation did not create egress-init.sh")
	}
	jsonOn, _ := fs.ReadFile(devcontainerJSONPath())
	dockerOn, _ := fs.ReadFile(dockerfilePath())
	if string(jsonOn) == string(jsonBefore) || string(dockerOn) == string(dockerBefore) {
		t.Errorf("activation left devcontainer.json or Dockerfile unchanged")
	}

	seedUBootYAMLWithFeatures(t, fs, off)
	resp, err = generateSandbox(svc, false)
	if err != nil {
		t.Fatalf("switch off: %v", err)
	}
	if n := egressRebuildNotices(resp.Warnings); len(n) != 0 {
		t.Errorf("switching off gave the rebuild notice: %+v", n)
	}
	if ok, _ := fs.Exists(egressInitPath()); !ok {
		t.Errorf("switching off removed egress-init.sh")
	}
	jsonOff, _ := fs.ReadFile(devcontainerJSONPath())
	dockerOff, _ := fs.ReadFile(dockerfilePath())
	if string(jsonOff) == string(jsonOn) || string(dockerOff) == string(dockerOn) {
		t.Errorf("switching off left devcontainer.json or Dockerfile unchanged")
	}
}
