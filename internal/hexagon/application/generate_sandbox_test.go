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
	// The default bind mount is switched off; the workspace is a
	// per-instance volume in `mounts` (${devcontainerId} is resolved
	// there, not in workspaceMount).
	if got := m["workspaceMount"]; got != "" {
		t.Errorf("workspaceMount = %q, want empty (no host bind mount)", got)
	}
	wsMounts, _ := m["mounts"].([]any)
	if len(wsMounts) != 1 || wsMounts[0] != "source=t-uboot-gen-workspace-${devcontainerId},target=/workspaces/t-uboot-gen,type=volume" {
		t.Errorf("mounts = %v, want the per-instance workspace volume", m["mounts"])
	}
	if m["workspaceFolder"] != "/workspaces/t-uboot-gen" {
		t.Errorf("workspaceFolder = %v", m["workspaceFolder"])
	}
	if m["postCreateCommand"] != "[ -d .git ] || git clone -- git@github.com:pt9912/demo.git ." {
		t.Errorf("postCreateCommand = %v", m["postCreateCommand"])
	}
	for _, forbidden := range []string{"runArgs", "privileged", "capAdd", "securityOpt"} {
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
	// LH-FA-DEV-009: token only as a runtime reference, never a value.
	env, _ := m["remoteEnv"].(map[string]any)
	if env["GIT_TOKEN"] != "${localEnv:GIT_TOKEN}" {
		t.Errorf("remoteEnv = %v, want GIT_TOKEN pass-through", m["remoteEnv"])
	}
	if !strings.Contains(string(df), "credential.helper") {
		t.Errorf("Dockerfile lacks the runtime credential helper:\n%s", df)
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

const sandboxPodmanYAML = "devcontainer:\n  enabled: true\n  profile: sandbox\n  sandbox:\n    nestedRuntime: podman\n"

func sandboxInitPath() string {
	return filepath.Join(generateTestBaseDir, ".devcontainer", "sandbox-init.sh")
}

// LH-FA-DEV-007: nestedRuntime podman adds runArgs with exactly the
// measured relaxations, the storage volume, the startup script (COPY
// + postCreateCommand) and reports every relaxation individually.
func TestGenerateDevcontainer_SandboxPodman(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, sandboxPodmanYAML)
	seedGitOrigin(t, fs, "https://github.com/pt9912/demo.git")

	resp, err := generateSandbox(svc, false)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	m := devcontainerJSONMap(t, fs)
	wantArgs := []any{"--cap-add=SYS_ADMIN", "--security-opt=seccomp=unconfined",
		"--security-opt=apparmor=unconfined", "--security-opt=systempaths=unconfined", "--device=/dev/fuse"}
	gotArgs, _ := m["runArgs"].([]any)
	if len(gotArgs) != len(wantArgs) {
		t.Fatalf("runArgs = %v, want %v", gotArgs, wantArgs)
	}
	for i := range wantArgs {
		if gotArgs[i] != wantArgs[i] {
			t.Errorf("runArgs[%d] = %v, want %v", i, gotArgs[i], wantArgs[i])
		}
	}
	for _, forbidden := range []string{"privileged", "capAdd", "securityOpt"} {
		if _, ok := m[forbidden]; ok {
			t.Errorf("must not contain %q (no --privileged)", forbidden)
		}
	}
	mounts, _ := m["mounts"].([]any)
	if len(mounts) != 2 ||
		mounts[0] != "source=t-uboot-gen-workspace-${devcontainerId},target=/workspaces/t-uboot-gen,type=volume" ||
		mounts[1] != "source=t-uboot-gen-containers-${devcontainerId},target=/home/vscode/.local/share/containers,type=volume" {
		t.Errorf("mounts = %v", mounts)
	}
	for _, mnt := range mounts {
		if strings.Contains(mnt.(string), "docker.sock") || strings.Contains(mnt.(string), "type=bind") {
			t.Errorf("no socket/bind mounts allowed: %v", mnt)
		}
	}
	if m["postCreateCommand"] != "sh /usr/local/bin/u-boot-sandbox-init && { [ -d .git ] || git clone -- https://github.com/pt9912/demo.git .; }" {
		t.Errorf("postCreateCommand = %v", m["postCreateCommand"])
	}
	df, _ := fs.ReadFile(dockerfilePath())
	for _, want := range []string{"podman uidmap fuse-overlayfs passt", "COPY sandbox-init.sh /usr/local/bin/u-boot-sandbox-init", "ln -sf /usr/bin/podman /usr/local/bin/docker"} {
		if !strings.Contains(string(df), want) {
			t.Errorf("Dockerfile lacks %q:\n%s", want, df)
		}
	}
	if strings.Index(string(df), "USER root") > strings.LastIndex(string(df), "USER vscode") {
		t.Errorf("Dockerfile must end on USER vscode")
	}
	script, err := fs.ReadFile(sandboxInitPath())
	if err != nil {
		t.Fatalf("sandbox-init.sh not generated: %v", err)
	}
	if !strings.Contains(string(script), `on_unavailable="warn"`) || !strings.Contains(string(script), "exit 11") {
		t.Errorf("script lacks policy/exit 11:\n%s", script)
	}
	relaxations := 0
	for _, w := range resp.Warnings {
		if w.Code == "LH-FA-DEV-007" && strings.Contains(w.Message, "security relaxation") {
			relaxations++
		}
	}
	if relaxations != 5 {
		t.Errorf("relaxation warnings = %d, want 5: %+v", relaxations, resp.Warnings)
	}

	// Idempotent.
	if resp2, err := generateSandbox(svc, false); err != nil || resp2.Action != driving.GenerateActionNoOp {
		t.Errorf("second run: action=%v err=%v, want NoOp", resp2.Action, err)
	}
}

// onUnavailable: fail is baked into the script; without a remote the
// postCreateCommand is only the init script.
func TestGenerateDevcontainer_SandboxPodman_FailPolicyNoRemote(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, sandboxPodmanYAML+"    onUnavailable: fail\n")

	if _, err := generateSandbox(svc, false); err != nil {
		t.Fatalf("generate: %v", err)
	}
	script, _ := fs.ReadFile(sandboxInitPath())
	if !strings.Contains(string(script), `on_unavailable="fail"`) {
		t.Errorf("fail policy not baked into script:\n%s", script)
	}
	if got := devcontainerJSONMap(t, fs)["postCreateCommand"]; got != "sh /usr/local/bin/u-boot-sandbox-init" {
		t.Errorf("postCreateCommand = %v", got)
	}
}

// nestedRuntime none (default) and podman-without-sandbox generate no
// nested runtime; the latter warns.
func TestGenerateDevcontainer_SandboxPodman_NotGenerated(t *testing.T) {
	t.Parallel()
	for name, tc := range map[string]struct {
		yaml     string
		wantWarn bool
	}{
		"none":            {"devcontainer:\n  enabled: true\n  profile: sandbox\n", false},
		"podman, default": {"devcontainer:\n  enabled: true\n  sandbox:\n    nestedRuntime: podman\n", true},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, tc.yaml)
			resp, err := generateSandbox(svc, false)
			if err != nil {
				t.Fatalf("generate: %v", err)
			}
			if exists, _ := fs.Exists(sandboxInitPath()); exists {
				t.Errorf("sandbox-init.sh must not be generated")
			}
			m := devcontainerJSONMap(t, fs)
			if _, ok := m["runArgs"]; ok {
				t.Errorf("runArgs must be absent")
			}
			mnts, hasMounts := m["mounts"].([]any)
			if wantMounts := name == "none"; hasMounts != wantMounts || (hasMounts && len(mnts) != 1) {
				t.Errorf("mounts = %v, want only the workspace volume for the sandbox profile and none otherwise", m["mounts"])
			}
			df, _ := fs.ReadFile(dockerfilePath())
			if strings.Contains(string(df), "podman") {
				t.Errorf("Dockerfile must not install podman:\n%s", df)
			}
			hasWarn := false
			for _, w := range resp.Warnings {
				hasWarn = hasWarn || (w.Code == "LH-FA-DEV-007" && strings.Contains(w.Message, "no effect"))
			}
			if hasWarn != tc.wantWarn {
				t.Errorf("no-effect warning = %v, want %v (%+v)", hasWarn, tc.wantWarn, resp.Warnings)
			}
		})
	}
}

// An existing user-owned sandbox-init.sh without managed block is a
// manual conflict and nothing is written (no partial write).
func TestGenerateDevcontainer_SandboxPodman_ScriptConflict(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, sandboxPodmanYAML)
	if err := fs.WriteFile(sandboxInitPath(), []byte("#!/bin/sh\necho mine\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	_, err := generateSandbox(svc, false)
	if !errors.Is(err, driving.ErrGenerateManualConflict) {
		t.Fatalf("err = %v, want ErrGenerateManualConflict", err)
	}
	if exists, _ := fs.Exists(devcontainerJSONPath()); exists {
		t.Errorf("devcontainer.json must not be written on conflict")
	}
}

// LH-FA-DEV-006 (0.3.2): devcontainer.sandbox.repository replaces
// origin as the clone source; it works without any origin.
func TestGenerateDevcontainer_Sandbox_RepositoryOverridesOrigin(t *testing.T) {
	t.Parallel()
	yaml := "devcontainer:\n  enabled: true\n  profile: sandbox\n  sandbox:\n    repository: git@github.com:other/fork.git\n"
	t.Run("with origin", func(t *testing.T) {
		t.Parallel()
		svc, fs := newGenerateService(t)
		seedUBootYAMLWithFeatures(t, fs, yaml)
		seedGitOrigin(t, fs, "https://github.com/pt9912/demo.git")
		resp, err := generateSandbox(svc, false)
		if err != nil {
			t.Fatalf("generate: %v", err)
		}
		if got := devcontainerJSONMap(t, fs)["postCreateCommand"]; got != "[ -d .git ] || git clone -- git@github.com:other/fork.git ." {
			t.Errorf("postCreateCommand = %v, want clone of the configured repository", got)
		}
		if len(resp.Warnings) != 0 {
			t.Errorf("unexpected warnings: %+v", resp.Warnings)
		}
	})
	t.Run("without origin", func(t *testing.T) {
		t.Parallel()
		svc, fs := newGenerateService(t)
		seedUBootYAMLWithFeatures(t, fs, yaml)
		resp, err := generateSandbox(svc, false)
		if err != nil || len(resp.Warnings) != 0 {
			t.Fatalf("err=%v warnings=%+v, want a clone step without warning", err, resp.Warnings)
		}
		if _, ok := devcontainerJSONMap(t, fs)["postCreateCommand"]; !ok {
			t.Errorf("clone step missing")
		}
	})
	t.Run("credentials rejected on load", func(t *testing.T) {
		t.Parallel()
		svc, fs := newGenerateService(t)
		seedUBootYAMLWithFeatures(t, fs, "devcontainer:\n  enabled: true\n  profile: sandbox\n  sandbox:\n    repository: https://tok@github.com/o/r.git\n")
		_, err := generateSandbox(svc, false)
		if err == nil || !errors.Is(err, driving.ErrGenerateManualConflict) || !strings.Contains(err.Error(), "credentials") {
			t.Fatalf("err = %v, want schema error naming credentials", err)
		}
	})
}

const sandboxEgressYAML = "devcontainer:\n  enabled: true\n  profile: sandbox\n  sandbox:\n    egress:\n      enabled: true\n      allow:\n        - api.anthropic.com\n"

func egressInitPath() string {
	return filepath.Join(generateTestBaseDir, ".devcontainer", "egress-init.sh")
}

// LH-FA-DEV-008 / ADR-0012: egress restriction adds NET_ADMIN, the
// script with default + user hosts (sorted, clone host included), the
// postStart hook, the nftables/dnsmasq install; nothing without the flag.
func TestGenerateDevcontainer_SandboxEgress(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, sandboxEgressYAML+"  features:\n    node:\n      enabled: true\n")
	seedGitOrigin(t, fs, "git@gitlab.example.org:o/r.git")

	resp, err := generateSandbox(svc, false)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	m := devcontainerJSONMap(t, fs)
	args, _ := m["runArgs"].([]any)
	if len(args) != 1 || args[0] != "--cap-add=NET_ADMIN" {
		t.Errorf("runArgs = %v, want only --cap-add=NET_ADMIN", args)
	}
	if m["postStartCommand"] != "sudo sh /usr/local/bin/u-boot-egress-init" {
		t.Errorf("postStartCommand = %v", m["postStartCommand"])
	}
	if m["postCreateCommand"] != "[ -d .git ] || git clone -- git@gitlab.example.org:o/r.git ." {
		t.Errorf("postCreateCommand = %v (clone must still run before the restriction)", m["postCreateCommand"])
	}
	script, err := fs.ReadFile(egressInitPath())
	if err != nil {
		t.Fatalf("egress-init.sh not generated: %v", err)
	}
	for _, want := range []string{"api.anthropic.com", "github.com", "gitlab.example.org", "registry.npmjs.org", "deb.debian.org", `on_unavailable="warn"`} {
		if !strings.Contains(string(script), want) {
			t.Errorf("script lacks host/policy %q:\n%s", want, script)
		}
	}
	for _, unwanted := range []string{"registry-1.docker.io", "proxy.golang.org"} {
		if strings.Contains(string(script), unwanted) {
			t.Errorf("script must not contain %q (podman/go not enabled)", unwanted)
		}
	}
	df, _ := fs.ReadFile(dockerfilePath())
	if !strings.Contains(string(df), "dnsmasq-base nftables") || !strings.Contains(string(df), "COPY egress-init.sh /usr/local/bin/u-boot-egress-init") {
		t.Errorf("Dockerfile lacks the nftables/dnsmasq/COPY steps:\n%s", df)
	}
	found := false
	for _, w := range resp.Warnings {
		found = found || (w.Code == "LH-FA-DEV-008" && strings.Contains(w.Message, "NET_ADMIN"))
	}
	if !found {
		t.Errorf("NET_ADMIN relaxation not reported: %+v", resp.Warnings)
	}
	if r2, err := generateSandbox(svc, false); err != nil || r2.Action != driving.GenerateActionNoOp {
		t.Errorf("second run: %v %v, want NoOp", r2.Action, err)
	}
}

// Podman + egress: both relaxation sets in runArgs, both init steps in
// the right lifecycle hook, podman registries in the allowlist.
func TestGenerateDevcontainer_SandboxEgressWithPodman(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, sandboxPodmanYAML+"    egress:\n      enabled: true\n")
	seedGitOrigin(t, fs, "https://github.com/o/r.git")
	if _, err := generateSandbox(svc, false); err != nil {
		t.Fatalf("generate: %v", err)
	}
	m := devcontainerJSONMap(t, fs)
	if args, _ := m["runArgs"].([]any); len(args) != 6 || args[5] != "--cap-add=NET_ADMIN" {
		t.Errorf("runArgs = %v, want 5 podman options + NET_ADMIN", m["runArgs"])
	}
	if m["postCreateCommand"] != "sh /usr/local/bin/u-boot-sandbox-init && { [ -d .git ] || git clone -- https://github.com/o/r.git .; }" {
		t.Errorf("postCreateCommand = %v", m["postCreateCommand"])
	}
	script, _ := fs.ReadFile(egressInitPath())
	if !strings.Contains(string(script), "registry-1.docker.io") {
		t.Errorf("podman registries missing from the allowlist:\n%s", script)
	}
}

// Egress off (default) or without the sandbox profile: no script, no
// NET_ADMIN; the latter warns.
func TestGenerateDevcontainer_SandboxEgress_NotGenerated(t *testing.T) {
	t.Parallel()
	for name, tc := range map[string]struct {
		yaml     string
		wantWarn bool
	}{
		"disabled":        {"devcontainer:\n  enabled: true\n  profile: sandbox\n  sandbox:\n    egress:\n      enabled: false\n", false},
		"default profile": {"devcontainer:\n  enabled: true\n  sandbox:\n    egress:\n      enabled: true\n", true},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, tc.yaml)
			resp, err := generateSandbox(svc, false)
			if err != nil {
				t.Fatalf("generate: %v", err)
			}
			if exists, _ := fs.Exists(egressInitPath()); exists {
				t.Errorf("egress-init.sh must not be generated")
			}
			m := devcontainerJSONMap(t, fs)
			if _, ok := m["runArgs"]; ok {
				t.Errorf("runArgs must be absent: %v", m["runArgs"])
			}
			if _, ok := m["postStartCommand"]; ok {
				t.Errorf("postStartCommand must be absent")
			}
			hasWarn := false
			for _, w := range resp.Warnings {
				hasWarn = hasWarn || (w.Code == "LH-FA-DEV-008" && strings.Contains(w.Message, "no effect"))
			}
			if hasWarn != tc.wantWarn {
				t.Errorf("no-effect warning = %v, want %v", hasWarn, tc.wantWarn)
			}
		})
	}
}
