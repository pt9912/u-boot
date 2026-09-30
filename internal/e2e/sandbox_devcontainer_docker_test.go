//go:build docker

// End-to-end check of the Lastenheft 0.3.x sandbox devcontainer
// (LH-FA-DEV-004/-006/-007/-009): generate the files with the
// production adapters, build the image with the real Docker engine,
// start it with the generated runArgs and run the generated startup
// script plus a nested `podman run`. Colima/macOS and Podman as the
// host engine are not covered (ADR 0014, Nachholmessung).

package e2e_test

import (
	"context"
	"encoding/json"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	confirmadapter "github.com/pt9912/u-boot/internal/adapter/driven/confirm"
	fsadapter "github.com/pt9912/u-boot/internal/adapter/driven/fs"
	gitadapter "github.com/pt9912/u-boot/internal/adapter/driven/git"
	loggeradapter "github.com/pt9912/u-boot/internal/adapter/driven/logger"
	progressadapter "github.com/pt9912/u-boot/internal/adapter/driven/progress"
	yamladapter "github.com/pt9912/u-boot/internal/adapter/driven/yaml"
	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

const sandboxCanaryToken = "ghp_canaryCANARYcanary0123456789abcdef"

func dockerOutput(ctx context.Context, t *testing.T, args ...string) (string, int) {
	t.Helper()
	out, err := exec.CommandContext(ctx, "docker", args...).CombinedOutput()
	code := 0
	if err != nil {
		ee, ok := err.(*exec.ExitError)
		if !ok {
			t.Fatalf("docker %v: %v", args, err)
		}
		code = ee.ExitCode()
	}
	return string(out), code
}

// generateSandboxProject produces a sandbox project with nested
// Podman in a temp dir and returns the dir.
func generateSandboxProject(ctx context.Context, t *testing.T, onUnavailable string) string {
	t.Helper()
	dir := t.TempDir()
	fsys := fsadapter.New()
	yaml := yamladapter.New()
	logger := loggeradapter.New(os.Stderr, loggeradapter.FormatText, nil)
	initSvc := application.NewInitProjectService(fsys, yaml, gitadapter.New(),
		progressadapter.NewText(os.Stderr), confirmadapter.New(strings.NewReader(""), os.Stderr), logger)
	if _, err := initSvc.Init(ctx, driving.InitProjectRequest{
		BaseDir: dir, Name: "sbxe2e", SkipGit: true, Devcontainer: true, Sandbox: true, SilenceProgress: true,
	}); err != nil {
		t.Fatalf("init --devcontainer --sandbox: %v", err)
	}
	cfgSvc := application.NewConfigService(fsys, yaml, logger)
	for path, value := range map[string]string{
		"devcontainer.sandbox.nestedRuntime": "podman",
		"devcontainer.sandbox.onUnavailable": onUnavailable,
	} {
		p, err := domain.NewConfigPath(path)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := cfgSvc.Set(ctx, driving.ConfigSetRequest{BaseDir: dir, Path: p, Value: value}); err != nil {
			t.Fatalf("config set %s: %v", path, err)
		}
	}
	genSvc := application.NewGenerateService(fsys, yaml, logger)
	if _, err := genSvc.Generate(ctx, driving.GenerateRequest{BaseDir: dir, Artifact: domain.ArtifactDevcontainer}); err != nil {
		t.Fatalf("generate devcontainer: %v", err)
	}
	return dir
}

func generatedMounts(t *testing.T, dir string) []string {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(dir, ".devcontainer", "devcontainer.json"))
	if err != nil {
		t.Fatal(err)
	}
	var lines []string
	for _, l := range strings.Split(string(body), "\n") {
		if !strings.HasPrefix(strings.TrimSpace(l), "//") {
			lines = append(lines, l)
		}
	}
	var dc struct {
		WorkspaceMount string   `json:"workspaceMount"`
		Mounts         []string `json:"mounts"`
	}
	if err := json.Unmarshal([]byte(strings.Join(lines, "\n")), &dc); err != nil {
		t.Fatalf("devcontainer.json: %v", err)
	}
	if dc.WorkspaceMount != "" {
		t.Fatalf("workspaceMount = %q, want empty (no host bind mount)", dc.WorkspaceMount)
	}
	return dc.Mounts
}

// Two instances of the same project (different ${devcontainerId},
// as the Dev Containers tooling resolves it per folder) must not
// share the workspace or the Podman storage volumes.
func TestE2E_SandboxDevcontainer_TwoInstancesIsolated(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker CLI not on PATH: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	dir := generateSandboxProject(ctx, t, "warn")
	image := "uboot-e2e-sandbox-multi:" + strings.ReplaceAll(filepath.Base(dir), "/", "-")
	if out, code := dockerOutput(ctx, t, "build", "-q", "-t", image, filepath.Join(dir, ".devcontainer")); code != 0 {
		t.Fatalf("docker build failed:\n%s", out)
	}
	t.Cleanup(func() { _, _ = dockerOutput(context.Background(), t, "rmi", "-f", image) })

	mounts := generatedMounts(t, dir)
	if len(mounts) != 2 {
		t.Fatalf("mounts = %v, want workspace + containers volume", mounts)
	}
	runArgs := generatedRunArgs(t, dir)
	start := func(id string) string {
		name := "uboot-e2e-inst-" + id
		args := []string{"run", "-d", "--name", name}
		for _, m := range mounts {
			if strings.Contains(m, "type=bind") {
				t.Fatalf("bind mount in sandbox mounts: %s", m)
			}
			args = append(args, "--mount", strings.ReplaceAll(m, "${devcontainerId}", id))
		}
		args = append(append(args, runArgs...), image, "sleep", "600")
		if out, code := dockerOutput(ctx, t, args...); code != 0 {
			t.Fatalf("start instance %s: %s", id, out)
		}
		t.Cleanup(func() {
			_, _ = dockerOutput(context.Background(), t, "rm", "-f", name)
			for _, m := range mounts {
				src := strings.TrimPrefix(strings.Split(strings.ReplaceAll(m, "${devcontainerId}", id), ",")[0], "source=")
				_, _ = dockerOutput(context.Background(), t, "volume", "rm", "-f", src)
			}
		})
		return name
	}
	a, b := start("idaaaa"), start("idbbbb")
	ws := "/workspaces/sbxe2e"
	if out, code := dockerOutput(ctx, t, "exec", a, "sh", "-c", "echo instance-a > "+ws+"/marker"); code != 0 {
		t.Fatalf("write in instance a: %s", out)
	}
	if out, code := dockerOutput(ctx, t, "exec", b, "sh", "-c", "test ! -e "+ws+"/marker && echo instance-b > "+ws+"/marker"); code != 0 {
		t.Errorf("instance b sees instance a's workspace: %s", out)
	}
	if out, _ := dockerOutput(ctx, t, "exec", a, "cat", ws+"/marker"); strings.TrimSpace(out) != "instance-a" {
		t.Errorf("instance a workspace content = %q", out)
	}
}

func generatedRunArgs(t *testing.T, dir string) []string {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(dir, ".devcontainer", "devcontainer.json"))
	if err != nil {
		t.Fatal(err)
	}
	var lines []string
	for _, l := range strings.Split(string(body), "\n") {
		if !strings.HasPrefix(strings.TrimSpace(l), "//") {
			lines = append(lines, l)
		}
	}
	var dc struct {
		RunArgs []string `json:"runArgs"`
	}
	if err := json.Unmarshal([]byte(strings.Join(lines, "\n")), &dc); err != nil {
		t.Fatalf("devcontainer.json: %v", err)
	}
	if len(dc.RunArgs) == 0 {
		t.Fatalf("generated devcontainer.json has no runArgs")
	}
	return dc.RunArgs
}

func TestE2E_SandboxDevcontainer_NestedPodman(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker CLI not on PATH: %v", err)
	}
	t.Setenv("GIT_TOKEN", sandboxCanaryToken)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	dir := generateSandboxProject(ctx, t, "warn")

	// LH-FA-DEV-009 / AK 4: no credential value in any generated file.
	_ = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			if b, rerr := os.ReadFile(p); rerr == nil && strings.Contains(string(b), sandboxCanaryToken) {
				t.Errorf("credential value found in generated file %s", p)
			}
		}
		return nil
	})

	image := "uboot-e2e-sandbox:" + strings.ReplaceAll(filepath.Base(dir), "/", "-")
	if out, code := dockerOutput(ctx, t, "build", "-q", "-t", image, filepath.Join(dir, ".devcontainer")); code != 0 {
		t.Fatalf("docker build of the generated Dockerfile failed:\n%s", out)
	}
	t.Cleanup(func() { _, _ = dockerOutput(context.Background(), t, "rmi", "-f", image) })

	runArgs := generatedRunArgs(t, dir)
	nested := "sh /usr/local/bin/u-boot-sandbox-init && docker run --rm --network=host docker.io/library/alpine:3 echo NESTED_OK"

	// AK 2: docker build/run inside the container without a socket.
	args := append(append([]string{"run", "--rm"}, runArgs...), image, "bash", "-c", nested)
	if out, code := dockerOutput(ctx, t, args...); code != 0 || !strings.Contains(out, "NESTED_OK") {
		t.Fatalf("nested podman run failed (exit %d):\n%s", code, out)
	}

	// Degradation: without the relaxations the start script aborts
	// with exit 11 (blocked user namespaces).
	if out, code := dockerOutput(ctx, t, "run", "--rm", image, "sh", "/usr/local/bin/u-boot-sandbox-init"); code != 11 {
		t.Errorf("without relaxations: exit %d, want 11:\n%s", code, out)
	}

	// Degradation: without /dev/fuse and policy warn → vfs fallback.
	var noFuse []string
	for _, a := range runArgs {
		if a != "--device=/dev/fuse" {
			noFuse = append(noFuse, a)
		}
	}
	vfs := "sh /usr/local/bin/u-boot-sandbox-init && grep -q vfs ~/.config/containers/storage.conf && echo VFS_FALLBACK"
	args = append(append([]string{"run", "--rm"}, noFuse...), image, "bash", "-c", vfs)
	if out, code := dockerOutput(ctx, t, args...); code != 0 || !strings.Contains(out, "VFS_FALLBACK") {
		t.Errorf("warn policy without /dev/fuse: exit %d:\n%s", code, out)
	}
}

func TestE2E_SandboxDevcontainer_FailPolicyWithoutFuse(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker CLI not on PATH: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	dir := generateSandboxProject(ctx, t, "fail")
	image := "uboot-e2e-sandbox-fail:" + strings.ReplaceAll(filepath.Base(dir), "/", "-")
	if out, code := dockerOutput(ctx, t, "build", "-q", "-t", image, filepath.Join(dir, ".devcontainer")); code != 0 {
		t.Fatalf("docker build failed:\n%s", out)
	}
	t.Cleanup(func() { _, _ = dockerOutput(context.Background(), t, "rmi", "-f", image) })

	var noFuse []string
	for _, a := range generatedRunArgs(t, dir) {
		if a != "--device=/dev/fuse" {
			noFuse = append(noFuse, a)
		}
	}
	args := append(append([]string{"run", "--rm"}, noFuse...), image, "sh", "/usr/local/bin/u-boot-sandbox-init")
	if out, code := dockerOutput(ctx, t, args...); code != 11 {
		t.Errorf("fail policy without /dev/fuse: exit %d, want 11:\n%s", code, out)
	}
}
