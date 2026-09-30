package application_test

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

func seedGitOrigin(t *testing.T, fs *fakeFS, url string) {
	t.Helper()
	body := "[core]\n\trepositoryformatversion = 0\n[remote \"origin\"]\n\turl = " + url + "\n\tfetch = +refs/heads/*:refs/remotes/origin/*\n"
	if err := fs.WriteFile(filepath.Join(generateTestBaseDir, ".git", "config"), []byte(body), 0o644); err != nil {
		t.Fatalf("seed .git/config: %v", err)
	}
}

func generateSandbox(svc *application.GenerateService, sandbox bool) (driving.GenerateResponse, error) {
	return svc.Generate(context.Background(), driving.GenerateRequest{
		BaseDir: generateTestBaseDir, Artifact: domain.ArtifactDevcontainer, Sandbox: sandbox})
}

func devcontainerJSONMap(t *testing.T, fs *fakeFS) map[string]any {
	t.Helper()
	body, err := fs.ReadFile(devcontainerJSONPath())
	if err != nil {
		t.Fatalf("read devcontainer.json: %v", err)
	}
	stripped := application.StripJSONCForTest(body)
	var m map[string]any
	if err := json.Unmarshal(stripped, &m); err != nil {
		t.Fatalf("devcontainer.json invalid: %v\n%s", err, stripped)
	}
	return m
}

// LH-FA-DEV-006: --sandbox renders a named-volume workspace with a
// clone step, no bind mount / socket mount / privileged, and persists
// the profile.
func TestGenerateDevcontainer_Sandbox_WithRemote(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLPostgres(t, fs)
	seedGitOrigin(t, fs, "git@github.com:pt9912/demo.git")

	resp, err := generateSandbox(svc, true)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if len(resp.Warnings) != 0 {
		t.Errorf("unexpected warnings: %+v", resp.Warnings)
	}
	m := devcontainerJSONMap(t, fs)
	if got := m["workspaceMount"]; got != "source=t-uboot-gen-workspace,target=/workspaces/t-uboot-gen,type=volume" {
		t.Errorf("workspaceMount = %v", got)
	}
	if m["workspaceFolder"] != "/workspaces/t-uboot-gen" {
		t.Errorf("workspaceFolder = %v", m["workspaceFolder"])
	}
	if m["postCreateCommand"] != "[ -d .git ] || git clone -- git@github.com:pt9912/demo.git ." {
		t.Errorf("postCreateCommand = %v", m["postCreateCommand"])
	}
	for _, forbidden := range []string{"runArgs", "mounts", "privileged", "capAdd", "securityOpt"} {
		if _, ok := m[forbidden]; ok {
			t.Errorf("sandbox devcontainer.json must not contain %q", forbidden)
		}
	}
	if m["remoteUser"] != "vscode" {
		t.Errorf("remoteUser = %v, want vscode (LH-FA-DEV-004)", m["remoteUser"])
	}
	df, _ := fs.ReadFile(dockerfilePath())
	if !strings.Contains(string(df), `mkdir -p "/workspaces/t-uboot-gen"`) {
		t.Errorf("Dockerfile lacks workspace dir:\n%s", df)
	}
	yamlBody, _ := fs.ReadFile(filepath.Join(generateTestBaseDir, "u-boot.yaml"))
	if !strings.Contains(string(yamlBody), "profile: sandbox") {
		t.Errorf("profile not persisted:\n%s", yamlBody)
	}
	found := false
	for _, c := range resp.Changed {
		found = found || c == "u-boot.yaml"
	}
	if !found {
		t.Errorf("Changed lacks u-boot.yaml: %v", resp.Changed)
	}

	// Idempotent: second run (profile now persisted) is a NoOp.
	resp2, err := generateSandbox(svc, true)
	if err != nil {
		t.Fatalf("second generate: %v", err)
	}
	if resp2.Action != driving.GenerateActionNoOp {
		t.Errorf("second run action = %v, want NoOp", resp2.Action)
	}
}

// Profile from config (no flag) renders the sandbox shape and does
// not rewrite u-boot.yaml.
func TestGenerateDevcontainer_Sandbox_FromConfigProfile(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, "devcontainer:\n  enabled: true\n  profile: sandbox\n")
	seedGitOrigin(t, fs, "https://github.com/pt9912/demo.git")

	resp, err := generateSandbox(svc, false)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if _, ok := devcontainerJSONMap(t, fs)["workspaceMount"]; !ok {
		t.Errorf("profile: sandbox from config did not render workspaceMount")
	}
	for _, c := range resp.Changed {
		if c == "u-boot.yaml" {
			t.Errorf("u-boot.yaml must not be rewritten when the profile is already persisted")
		}
	}
}

// No remote: no clone step, warning, exit 0 (LH-FA-DEV-006, 0.3.1).
func TestGenerateDevcontainer_Sandbox_NoRemote_WarnsWithoutCloneStep(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLPostgres(t, fs)

	resp, err := generateSandbox(svc, true)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if len(resp.Warnings) != 1 || resp.Warnings[0].Code != "LH-FA-DEV-006" || resp.Warnings[0].Level != "warn" {
		t.Fatalf("warnings = %+v, want one LH-FA-DEV-006 warn", resp.Warnings)
	}
	m := devcontainerJSONMap(t, fs)
	if _, ok := m["postCreateCommand"]; ok {
		t.Errorf("postCreateCommand must be absent without a remote")
	}
	if _, ok := m["workspaceMount"]; !ok {
		t.Errorf("workspaceMount must still be present")
	}
}

// Credentials or unsafe characters in origin → domain error (exit 10),
// nothing written (LH-FA-DEV-009).
func TestGenerateDevcontainer_Sandbox_RejectsUnsafeOrigin(t *testing.T) {
	t.Parallel()
	for _, url := range []string{
		"https://ghp_secret@github.com/o/r.git",
		"https://user:pass@github.com/o/r.git",
		"ssh://git:pw@github.com/o/r.git",
		"https://github.com/o/r.git;rm -rf /",
		"https://github.com/o/$(id).git",
		"-oProxyCommand=evil",
	} {
		t.Run(url, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLPostgres(t, fs)
			seedGitOrigin(t, fs, url)
			_, err := generateSandbox(svc, true)
			if !errors.Is(err, domain.ErrInvalidSandboxSetting) {
				t.Fatalf("err = %v, want ErrInvalidSandboxSetting", err)
			}
			if exists, _ := fs.Exists(devcontainerJSONPath()); exists {
				t.Errorf("devcontainer.json must not be written on rejected origin")
			}
		})
	}
	// Allowed forms: ssh:// with user only, scp-like, https w/o userinfo.
	for _, url := range []string{"ssh://git@github.com/o/r.git", "git@github.com:o/r.git", "https://github.com/o/r.git"} {
		svc, fs := newGenerateService(t)
		seedUBootYAMLPostgres(t, fs)
		seedGitOrigin(t, fs, url)
		if _, err := generateSandbox(svc, true); err != nil {
			t.Errorf("origin %q rejected: %v", url, err)
		}
	}
}

func TestGenerateDevcontainer_Sandbox_OnlyForDevcontainer(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLPostgres(t, fs)
	_, err := svc.Generate(context.Background(), driving.GenerateRequest{
		BaseDir: generateTestBaseDir, Artifact: domain.ArtifactChangelog, Sandbox: true})
	if !errors.Is(err, domain.ErrInvalidSandboxSetting) {
		t.Fatalf("err = %v, want ErrInvalidSandboxSetting", err)
	}
}

// LH-FA-DEV-004: UID 501 → build arg + usermod; default (unset or
// 1000) leaves the output free of both.
func TestGenerateDevcontainer_UID(t *testing.T) {
	t.Parallel()
	t.Run("custom", func(t *testing.T) {
		t.Parallel()
		svc, fs := newGenerateService(t)
		seedUBootYAMLWithFeatures(t, fs, "devcontainer:\n  enabled: true\n  user:\n    uid: 501\n")
		if _, err := generateSandbox(svc, false); err != nil {
			t.Fatalf("generate: %v", err)
		}
		build, _ := devcontainerJSONMap(t, fs)["build"].(map[string]any)
		args, _ := build["args"].(map[string]any)
		if args["USER_UID"] != "501" {
			t.Errorf("build.args = %v, want USER_UID=501", build["args"])
		}
		df, _ := fs.ReadFile(dockerfilePath())
		for _, want := range []string{"ARG USER_UID=501", `usermod --uid "${USER_UID}" vscode`, "USER vscode"} {
			if !strings.Contains(string(df), want) {
				t.Errorf("Dockerfile lacks %q:\n%s", want, df)
			}
		}
		if strings.Index(string(df), "USER root") > strings.LastIndex(string(df), "USER vscode") {
			t.Errorf("Dockerfile must end on USER vscode (LH-FA-DEV-004):\n%s", df)
		}
	})
	for name, yamlDC := range map[string]string{
		"unset":   "devcontainer:\n  enabled: true\n",
		"default": "devcontainer:\n  enabled: true\n  user:\n    uid: 1000\n",
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, yamlDC)
			if _, err := generateSandbox(svc, false); err != nil {
				t.Fatalf("generate: %v", err)
			}
			build, _ := devcontainerJSONMap(t, fs)["build"].(map[string]any)
			if _, ok := build["args"]; ok {
				t.Errorf("no build args expected for default UID")
			}
			df, _ := fs.ReadFile(dockerfilePath())
			if strings.Contains(string(df), "usermod") || strings.Contains(string(df), "USER root") {
				t.Errorf("default UID must keep the pre-0.3.0 Dockerfile:\n%s", df)
			}
		})
	}
}

// Switching an existing default devcontainer to the sandbox profile
// replaces only the managed block: content outside stays intact
// (replay rule: stale block replacement + user content preservation).
func TestGenerateDevcontainer_Sandbox_SwitchKeepsUserContent(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLPostgres(t, fs)
	seedGitOrigin(t, fs, "https://github.com/pt9912/demo.git")
	if _, err := generateSandbox(svc, false); err != nil {
		t.Fatalf("default generate: %v", err)
	}
	df, _ := fs.ReadFile(dockerfilePath())
	custom := string(df) + "\nRUN echo user-content\n"
	if err := fs.WriteFile(dockerfilePath(), []byte(custom), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	resp, err := generateSandbox(svc, true)
	if err != nil {
		t.Fatalf("sandbox generate: %v", err)
	}
	if resp.Action != driving.GenerateActionUpdatedBlock {
		t.Errorf("action = %v, want UpdatedBlock", resp.Action)
	}
	after, _ := fs.ReadFile(dockerfilePath())
	if !strings.HasSuffix(string(after), "\nRUN echo user-content\n") {
		t.Errorf("user content after the managed block was lost:\n%s", after)
	}
	if !strings.Contains(string(after), "workspaces/t-uboot-gen") {
		t.Errorf("block not switched to sandbox:\n%s", after)
	}
}
