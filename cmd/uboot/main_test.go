package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestRun_VersionFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run(context.Background(), []string{"--version"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("--version exit code: got %d, want 0", code)
	}
	if got := strings.TrimSpace(stdout.String()); got == "" {
		t.Fatalf("--version stdout: empty; expected non-empty version string")
	}
}

func TestRun_HelpFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run(context.Background(), []string{"--help"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("--help exit code: got %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "u-boot") {
		t.Fatalf("--help stdout: missing 'u-boot'; got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "init") {
		t.Fatalf("--help stdout: missing subcommand 'init'; got %q", stdout.String())
	}
}

func TestRun_UnknownCommandExits2(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run(context.Background(), []string{"frobnicate"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("unknown command exit code: got %d, want 2 (LH-FA-CLI-006)", code)
	}
	if !strings.Contains(stderr.String(), "u-boot:") {
		t.Fatalf("unknown command stderr: missing 'u-boot:' prefix; got %q", stderr.String())
	}
}

func TestRun_UnknownFlagExits2(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run(context.Background(), []string{"init", "--no-such-flag"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("unknown flag exit code: got %d, want 2 (LH-FA-CLI-006)", code)
	}
}

func TestVersionConstantNonEmpty(t *testing.T) {
	if version == "" {
		t.Fatal("package-level version must not be empty")
	}
}

func TestRun_InitHappyPath_WiresAllAdapters(t *testing.T) {
	// Integration-style test against the real wiring: real fs, yaml,
	// git adapters all exercised once. Without --no-git this would
	// also exercise the git adapter (host git binary required);
	// staying with --no-git keeps the test self-contained on lean
	// CI runners.
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	dir := t.TempDir()
	t.Chdir(dir) // Go 1.24+; PWD is the resolveProjectName fallback source.

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"init", "demo", "--no-git"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("init demo --no-git: exit %d, stderr=%q", code, stderr.String())
	}

	// Spot-check the artefacts produced (LH-FA-INIT-003 mandatory
	// set). Full per-file content is asserted in the application
	// layer's unit tests.
	for _, rel := range []string{
		"u-boot.yaml", "compose.yaml", "README.md", "CHANGELOG.md",
		".env.example", ".gitignore",
		"docker", "scripts", "docs",
	} {
		if _, err := os.Stat(filepath.Join(dir, rel)); err != nil {
			t.Errorf("init did not create %s: %v", rel, err)
		}
	}

	if !strings.Contains(stdout.String(), "demo") {
		t.Errorf("stdout does not mention project name; got %q", stdout.String())
	}
}

// LH-FA-DEV-008 at the user path with the real adapters: `generate devcontainer`
// writes the egress script, `egress list` shows the same host list (with the
// origin of each host), and the listing changes no file.
func TestRun_EgressList_MatchesTheGeneratedStartScript(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	dir := t.TempDir()
	t.Chdir(dir)
	writeFile := func(rel, body string) {
		t.Helper()
		path := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	writeFile("u-boot.yaml", "schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n  profile: sandbox\n"+
		"  features:\n    node:\n      enabled: true\n    claude-code:\n      enabled: true\n"+
		"  sandbox:\n    egress:\n      enabled: true\n      allow:\n        - api.example.org\n")
	writeFile(".git/config", "[remote \"origin\"]\n\turl = git@gitlab.example.org:o/r.git\n")

	var out, errOut bytes.Buffer
	if code := run(context.Background(), []string{"generate", "devcontainer"}, &out, &errOut); code != 0 {
		t.Fatalf("generate devcontainer: exit %d, stderr=%q", code, errOut.String())
	}
	script, err := os.ReadFile(filepath.Join(dir, ".devcontainer", "egress-init.sh"))
	if err != nil {
		t.Fatalf("egress-init.sh: %v", err)
	}
	var scriptHosts []string
	for _, line := range strings.Split(string(script), "\n") {
		if rest, ok := strings.CutPrefix(line, `hosts="`); ok {
			scriptHosts = strings.Fields(strings.TrimSuffix(rest, `"`))
		}
	}
	if len(scriptHosts) == 0 {
		t.Fatalf("no hosts= line in the script:\n%s", script)
	}
	before := snapshotTree(t, dir)

	out.Reset()
	errOut.Reset()
	if code := run(context.Background(), []string{"egress", "list", "--json"}, &out, &errOut); code != 0 {
		t.Fatalf("egress list: exit %d, stderr=%q", code, errOut.String())
	}
	var env struct {
		Data struct {
			Active bool `json:"active"`
			Hosts  []struct {
				Host    string   `json:"host"`
				Sources []string `json:"sources"`
			} `json:"hosts"`
		} `json:"data"`
	}
	if err := json.Unmarshal(out.Bytes(), &env); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	var listed []string
	origins := map[string]string{}
	for _, h := range env.Data.Hosts {
		listed = append(listed, h.Host)
		origins[h.Host] = strings.Join(h.Sources, ",")
	}
	if !env.Data.Active || !reflect.DeepEqual(listed, scriptHosts) {
		t.Errorf("active=%v; listing differs from the script:\n list   %v\n script %v", env.Data.Active, listed, scriptHosts)
	}
	for host, want := range map[string]string{
		"gitlab.example.org": "clone", "api.example.org": "user", "registry.npmjs.org": "feature:claude-code,feature:node", "github.com": "base",
	} {
		if origins[host] != want {
			t.Errorf("origin of %s = %q, want %q", host, origins[host], want)
		}
	}
	if after := snapshotTree(t, dir); !reflect.DeepEqual(before, after) {
		t.Errorf("egress list changed the project:\n before %v\n after  %v", before, after)
	}
}

// LH-FA-DEV-008, LH-FA-DEV-006, LH-FA-CLI-006 at the user path: an `origin` URL
// with credentials ends `egress list` with exit 10 and the diagnostic
// LH-FA-DEV-006 also for a project without the sandbox profile and without an
// active restriction; the listing changes no file.
func TestRun_EgressList_UnsafeOriginExits10WithoutActiveRestriction(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	dir := t.TempDir()
	t.Chdir(dir)
	for rel, body := range map[string]string{
		"u-boot.yaml": "schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n",
		".git/config": "[remote \"origin\"]\n\turl = https://user:secret@example.org/o/r.git\n",
	} {
		path := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	before := snapshotTree(t, dir)

	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"egress", "list", "--json"}, &out, &errOut)
	if code != 10 {
		t.Fatalf("egress list: exit %d, want 10; stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	var env struct {
		Status      string `json:"status"`
		ExitCode    int    `json:"exitCode"`
		Diagnostics []struct {
			Code string `json:"code"`
		} `json:"diagnostics"`
	}
	if err := json.Unmarshal(out.Bytes(), &env); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if env.Status != "error" || env.ExitCode != 10 || len(env.Diagnostics) != 1 || env.Diagnostics[0].Code != "LH-FA-DEV-006" {
		t.Errorf("envelope = %+v, want status error, exitCode 10 and one diagnostic LH-FA-DEV-006", env)
	}
	if after := snapshotTree(t, dir); !reflect.DeepEqual(before, after) {
		t.Errorf("the failed egress list changed the project:\n before %v\n after  %v", before, after)
	}
}

// snapshotTree maps every file below dir to its content.
func snapshotTree(t *testing.T, dir string) map[string]string {
	t.Helper()
	files := map[string]string{}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		body, err := os.ReadFile(path)
		files[path] = string(body)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}
