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

// LH-FA-DEV-003, LH-FA-CLI-005 at the user path (Verification-Report pins V-3
// and V-5): `generate devcontainer` reports the supply-chain limit of the
// `ollama` feature with the sandbox profile and the "no effect" statement
// without it; `--quiet` drops both, and a `version` override without
// `source` ends with exit 10.
func TestRun_Generate_OllamaNotices(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	const head = "schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n"
	const feature = "  features:\n    ollama:\n      enabled: true\n"
	for name, tc := range map[string]struct {
		yaml     string
		wantCode int
		wantText string
	}{
		"sandbox":         {head + "  profile: sandbox\n" + feature, 0, "Warning (LH-NFA-SEC-004): devcontainer feature ollama"},
		"default profile": {head + feature, 0, "Warning (LH-FA-DEV-003): devcontainer feature ollama has no effect without devcontainer.profile=sandbox"},
		"version without source": {head + "  profile: sandbox\n  features:\n    ollama:\n      enabled: true\n      version: \"0.40.7\"\n",
			10, ""},
	} {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			t.Chdir(dir)
			if err := os.WriteFile(filepath.Join(dir, "u-boot.yaml"), []byte(tc.yaml), 0o644); err != nil {
				t.Fatal(err)
			}
			if err := os.MkdirAll(filepath.Join(dir, ".git"), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, ".git", "config"), []byte("[remote \"origin\"]\n\turl = git@gitlab.example.org:o/r.git\n"), 0o644); err != nil {
				t.Fatal(err)
			}

			var out, errOut bytes.Buffer
			if code := run(context.Background(), []string{"generate", "devcontainer"}, &out, &errOut); code != tc.wantCode {
				t.Fatalf("generate devcontainer: exit %d, want %d; stdout=%q stderr=%q", code, tc.wantCode, out.String(), errOut.String())
			}
			if tc.wantText != "" && !strings.Contains(out.String(), tc.wantText) {
				t.Errorf("stdout lacks %q:\n%s", tc.wantText, out.String())
			}
			if tc.wantCode != 0 {
				return
			}

			out.Reset()
			errOut.Reset()
			if code := run(context.Background(), []string{"--quiet", "generate", "devcontainer"}, &out, &errOut); code != 0 {
				t.Fatalf("--quiet generate devcontainer: exit %d; stdout=%q stderr=%q", code, out.String(), errOut.String())
			}
			if strings.Contains(strings.ToLower(out.String()+errOut.String()), "ollama") {
				t.Errorf("--quiet must drop the ollama notice:\nstdout=%q\nstderr=%q", out.String(), errOut.String())
			}
		})
	}
}

// writeFeatureProject writes a project folder with the given u-boot.yaml and a
// git remote, and makes it the working directory.
func writeFeatureProject(t *testing.T, yaml string) string {
	t.Helper()
	dir := t.TempDir()
	t.Chdir(dir)
	if err := os.WriteFile(filepath.Join(dir, "u-boot.yaml"), []byte(yaml), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".git", "config"), []byte("[remote \"origin\"]\n\turl = git@gitlab.example.org:o/r.git\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

// LH-FA-DEV-006, LH-FA-CLI-005 at the user path (real adapters): in the
// sandbox profile `generate devcontainer` discloses the two options of `go`
// and each feature from an external source, and `--quiet` keeps these
// warnings (they carry SecurityRelaxation); the default profile discloses
// none of them.
func TestRun_Generate_SandboxFeatureWarnings(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	const head = "schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n  featureSources:\n    allow:\n      - oci://registry.example.test/acme\n"
	const features = "  features:\n    go:\n      enabled: true\n    acme:\n      enabled: true\n      source: oci://registry.example.test/acme\n"
	want := []string{
		"Warning (LH-FA-DEV-006): security relaxation for the devcontainer feature go: --cap-add=SYS_PTRACE",
		"Warning (LH-FA-DEV-006): security relaxation for the devcontainer feature go: --security-opt=seccomp=unconfined",
		"Warning (LH-FA-DEV-006): devcontainer feature acme comes from the external source oci://registry.example.test/acme",
	}
	for _, tc := range []struct {
		name  string
		yaml  string
		quiet bool
		shown bool
	}{
		{"sandbox profile", head + "  profile: sandbox\n" + features, false, true},
		{"sandbox profile, quiet", head + "  profile: sandbox\n" + features, true, true},
		{"default profile", head + features, false, false},
		{"default profile, quiet", head + features, true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			writeFeatureProject(t, tc.yaml)
			args := []string{"generate", "devcontainer"}
			if tc.quiet {
				args = append([]string{"--quiet"}, args...)
			}
			var out, errOut bytes.Buffer
			if code := run(context.Background(), args, &out, &errOut); code != 0 {
				t.Fatalf("%v: exit %d, want 0; stdout=%q stderr=%q", args, code, out.String(), errOut.String())
			}
			for _, w := range want {
				if got := strings.Contains(out.String(), w); got != tc.shown {
					t.Errorf("stdout contains %q = %v, want %v:\n%s", w, got, tc.shown, out.String())
				}
			}
		})
	}
}

// LH-FA-DEV-006, LH-FA-CLI-006 at the user path (real adapters): the sandbox
// profile with an active `docker-cli` ends `generate devcontainer` with exit
// 10 and the code LH-FA-DEV-006, also under `--dry-run` and `--json`, and
// creates nothing; `config set …docker-cli.enabled false` then succeeds and
// the next `generate devcontainer` ends with 0. The default profile keeps
// `docker-cli`.
func TestRun_Generate_SandboxHostSocketFeatureIsRejected(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	const base = "schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n"
	const feature = "  features:\n    docker-cli:\n      enabled: true\n"
	const remedy = "`u-boot config set devcontainer.features.docker-cli.enabled false`"

	dir := writeFeatureProject(t, base+"  profile: sandbox\n"+feature)
	before, err := os.ReadFile(filepath.Join(dir, "u-boot.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{
		{"generate", "devcontainer"},
		{"generate", "devcontainer", "--dry-run"},
		{"generate", "devcontainer", "--sandbox"},
	} {
		var out, errOut bytes.Buffer
		if code := run(context.Background(), args, &out, &errOut); code != 10 {
			t.Fatalf("%v: exit %d, want 10; stdout=%q stderr=%q", args, code, out.String(), errOut.String())
		}
		if !strings.Contains(errOut.String()+out.String(), remedy) {
			t.Errorf("%v: no output names the way out %s:\nstdout=%q\nstderr=%q", args, remedy, out.String(), errOut.String())
		}
	}
	var out, errOut bytes.Buffer
	if code := run(context.Background(), []string{"--json", "generate", "devcontainer", "--dry-run"}, &out, &errOut); code != 10 {
		t.Fatalf("--json --dry-run: exit %d, want 10; stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	var env devcontainerRunEnvelope
	if err := json.Unmarshal(out.Bytes(), &env); err != nil {
		t.Fatalf("--json envelope: %v\n%s", err, out.String())
	}
	if env.ExitCode != 10 || len(env.Diagnostics) != 1 || env.Diagnostics[0].Code != "LH-FA-DEV-006" || env.Diagnostics[0].Level != "error" {
		t.Errorf("envelope = %+v, want exit code 10 and one error diagnostic LH-FA-DEV-006", env)
	}
	if _, err := os.Stat(filepath.Join(dir, ".devcontainer")); err == nil {
		t.Error("the rejected calls created .devcontainer")
	}
	if after, _ := os.ReadFile(filepath.Join(dir, "u-boot.yaml")); !bytes.Equal(after, before) {
		t.Errorf("u-boot.yaml changed:\n%s", after)
	}

	out.Reset()
	errOut.Reset()
	if code := run(context.Background(), []string{"config", "set", "devcontainer.features.docker-cli.enabled", "false"}, &out, &errOut); code != 0 {
		t.Fatalf("config set …docker-cli.enabled false: exit %d; stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	out.Reset()
	errOut.Reset()
	if code := run(context.Background(), []string{"generate", "devcontainer"}, &out, &errOut); code != 0 {
		t.Fatalf("generate devcontainer after the way out: exit %d; stdout=%q stderr=%q", code, out.String(), errOut.String())
	}

	writeFeatureProject(t, base+feature)
	out.Reset()
	errOut.Reset()
	if code := run(context.Background(), []string{"generate", "devcontainer"}, &out, &errOut); code != 0 {
		t.Fatalf("default profile with docker-cli: exit %d, want 0; stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
}

// LH-FA-DEV-003, LH-NFA-SEC-004, LH-FA-CLI-006 at the user path (real
// adapters): an enabled feature whose source is outside the allowlist ends
// `generate devcontainer` with exit 10 and the code LH-FA-DEV-003, also under
// `--dry-run` and `--json`, and creates nothing. Each way out of the message
// leads to exit 0: allowing the source in the list, the flag of the call (the
// list gains the source only without `--dry-run`) and switching the feature off.
func TestRun_Generate_UnlistedFeatureSourceIsRejected(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	const src = "https://example.test/features/evil"
	const project = "schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n  features:\n    evil:\n      enabled: true\n      source: " + src + "\n"
	runArgs := func(args ...string) (int, string, string) {
		var out, errOut bytes.Buffer
		code := run(context.Background(), args, &out, &errOut)
		return code, out.String(), errOut.String()
	}

	dir := writeFeatureProject(t, project)
	before, err := os.ReadFile(filepath.Join(dir, "u-boot.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{
		{"generate", "devcontainer"},
		{"generate", "devcontainer", "--dry-run"},
	} {
		if code, out, errOut := runArgs(args...); code != 10 {
			t.Fatalf("%v: exit %d, want 10; stdout=%q stderr=%q", args, code, out, errOut)
		}
	}
	code, out, errOut := runArgs("--json", "generate", "devcontainer")
	if code != 10 {
		t.Fatalf("--json: exit %d, want 10; stdout=%q stderr=%q", code, out, errOut)
	}
	var env devcontainerRunEnvelope
	if err := json.Unmarshal([]byte(out), &env); err != nil {
		t.Fatalf("--json envelope: %v\n%s", err, out)
	}
	if env.ExitCode != 10 || len(env.Diagnostics) != 1 || env.Diagnostics[0].Code != "LH-FA-DEV-003" || env.Diagnostics[0].Level != "error" {
		t.Errorf("envelope = %+v, want exit code 10 and one error diagnostic LH-FA-DEV-003", env)
	}
	if _, err := os.Stat(filepath.Join(dir, ".devcontainer")); err == nil {
		t.Error("the rejected calls created .devcontainer")
	}
	if after, _ := os.ReadFile(filepath.Join(dir, "u-boot.yaml")); !bytes.Equal(after, before) {
		t.Errorf("u-boot.yaml changed:\n%s", after)
	}

	// Way out 1: allow the source in the list.
	writeFeatureProject(t, project)
	if code, out, errOut := runArgs("config", "set", "devcontainer.featureSources.allow", src); code != 0 {
		t.Fatalf("config set …featureSources.allow: exit %d; stdout=%q stderr=%q", code, out, errOut)
	}
	if code, out, errOut := runArgs("generate", "devcontainer"); code != 0 {
		t.Fatalf("generate after the allowlist way out: exit %d; stdout=%q stderr=%q", code, out, errOut)
	}

	// Way out 2: the flag of the same call.
	dir = writeFeatureProject(t, project)
	if code, out, errOut := runArgs("generate", "devcontainer", "--dry-run", "--allow-external-feature-sources", src); code != 0 {
		t.Fatalf("--dry-run with the flag: exit %d; stdout=%q stderr=%q", code, out, errOut)
	}
	if _, err := os.Stat(filepath.Join(dir, ".devcontainer")); err == nil {
		t.Error("--dry-run with the flag created .devcontainer")
	}
	if after, _ := os.ReadFile(filepath.Join(dir, "u-boot.yaml")); !bytes.Equal(after, []byte(project)) {
		t.Errorf("--dry-run with the flag changed u-boot.yaml:\n%s", after)
	}
	if code, out, errOut := runArgs("generate", "devcontainer", "--allow-external-feature-sources", src); code != 0 {
		t.Fatalf("generate with the flag: exit %d; stdout=%q stderr=%q", code, out, errOut)
	}
	if after, _ := os.ReadFile(filepath.Join(dir, "u-boot.yaml")); !strings.Contains(string(after), "- "+src) {
		t.Errorf("the flag source is not in the list:\n%s", after)
	}

	// Way out 3: switch the feature off.
	writeFeatureProject(t, project)
	if code, out, errOut := runArgs("config", "set", "devcontainer.features.evil.enabled", "false"); code != 0 {
		t.Fatalf("config set …evil.enabled false: exit %d; stdout=%q stderr=%q", code, out, errOut)
	}
	if code, out, errOut := runArgs("generate", "devcontainer"); code != 0 {
		t.Fatalf("generate after switching the feature off: exit %d; stdout=%q stderr=%q", code, out, errOut)
	}
}

// devcontainerRunEnvelope is the part of the `--json` envelope the
// devcontainer wiring tests read.
type devcontainerRunEnvelope struct {
	Status      string `json:"status"`
	Command     string `json:"command"`
	Subcommand  string `json:"subcommand"`
	ExitCode    int    `json:"exitCode"`
	Diagnostics []struct {
		Level   string `json:"level"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"diagnostics"`
}

// projectWithDevcontainer writes a project folder with a devcontainer
// configuration, makes it the working directory and removes every `docker`
// binary from the PATH, so the wiring meets the same situation as the
// released container image: no Docker client.
func projectWithDevcontainer(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".devcontainer", "devcontainer.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"name": "demo", "image": "debian"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)
	t.Setenv("PATH", t.TempDir())
	t.Setenv("DOCKER_HOST", "")
}

// LH-FA-DEV-010, LH-FA-CLI-006 at the user path: without a docker binary
// `devcontainer up` ends with exit 11 and names the socket relaxation in the
// envelope; nothing is written into the project.
func TestRun_DevcontainerUp_NoDockerBinaryExits11(t *testing.T) {
	projectWithDevcontainer(t)
	before := snapshotTree(t, ".")

	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"devcontainer", "up", "--json"}, &out, &errOut)
	if code != 11 {
		t.Fatalf("exit %d, want 11; stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	var env devcontainerRunEnvelope
	if err := json.Unmarshal(out.Bytes(), &env); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if env.Status != "error" || env.Command != "devcontainer" || env.Subcommand != "up" || env.ExitCode != 11 ||
		len(env.Diagnostics) != 2 || env.Diagnostics[0].Level != "warn" || env.Diagnostics[1].Level != "error" ||
		env.Diagnostics[0].Code != "LH-FA-DEV-010" || env.Diagnostics[1].Code != "LH-FA-DEV-010" ||
		!strings.Contains(env.Diagnostics[1].Message, "docker binary not on PATH") {
		t.Errorf("envelope = %+v", env)
	}
	if after := snapshotTree(t, "."); !reflect.DeepEqual(before, after) {
		t.Errorf("the failed start changed the project:\n before %v\n after  %v", before, after)
	}
}

// LH-FA-DEV-010: `devcontainer down` cannot ask the engine without a docker
// binary; that is an environment failure (exit 11), not "nothing to stop".
func TestRun_DevcontainerDown_NoDockerBinaryExits11(t *testing.T) {
	projectWithDevcontainer(t)

	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"devcontainer", "down", "--json"}, &out, &errOut)
	if code != 11 {
		t.Fatalf("exit %d, want 11; stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	var env devcontainerRunEnvelope
	if err := json.Unmarshal(out.Bytes(), &env); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if env.Status != "error" || env.Subcommand != "down" || env.ExitCode != 11 || len(env.Diagnostics) != 1 || env.Diagnostics[0].Code != "LH-FA-DEV-010" {
		t.Errorf("envelope = %+v", env)
	}
}

// LH-FA-DEV-010, LH-FA-CLI-006, ADR-0031 decision 2 at the user path: without a
// docker binary `devcontainer shell` ends with exit 11 before a session can
// begin, writes the socket disclosure and the error to stderr and nothing to
// stdout, and changes nothing in the project. `--json` is rejected with exit 2
// before any of that.
func TestRun_DevcontainerShell_NoDockerBinaryExits11(t *testing.T) {
	projectWithDevcontainer(t)
	before := snapshotTree(t, ".")

	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"devcontainer", "shell"}, &out, &errOut)
	if code != 11 || out.Len() != 0 {
		t.Fatalf("exit %d, stdout=%q stderr=%q, want exit 11 and no stdout", code, out.String(), errOut.String())
	}
	if !strings.Contains(errOut.String(), "Warning (LH-FA-DEV-010): this command mounts the host container-engine socket") ||
		!strings.Contains(errOut.String(), "for the duration of the start)") ||
		!strings.Contains(errOut.String(), "docker binary not on PATH") {
		t.Errorf("stderr lacks the disclosure of the start or the error:\n%s", errOut.String())
	}
	if after := snapshotTree(t, "."); !reflect.DeepEqual(before, after) {
		t.Errorf("the failed start changed the project:\n before %v\n after  %v", before, after)
	}

	out.Reset()
	errOut.Reset()
	code = run(context.Background(), []string{"devcontainer", "shell", "--json"}, &out, &errOut)
	var env devcontainerRunEnvelope
	if err := json.Unmarshal(out.Bytes(), &env); err != nil {
		t.Fatalf("--json: invalid JSON: %v\n%s", err, out.String())
	}
	if code != 2 || env.Command != "devcontainer" || env.Subcommand != "shell" || env.ExitCode != 2 || len(env.Diagnostics) != 1 ||
		env.Diagnostics[0].Code != "LH-FA-CLI-006" || strings.Contains(errOut.String(), "Warning") {
		t.Errorf("--json: exit %d, envelope %+v, stderr %q, want exit 2, subcommand shell and no disclosure", code, env, errOut.String())
	}
}

// LH-FA-DEV-010, LH-FA-CLI-006: a folder without `.devcontainer/devcontainer.json`
// ends with exit 10 before any docker call (no docker binary is needed to see
// it) and without the socket disclosure: no helper container was started.
func TestRun_DevcontainerUp_NoConfigurationExits10(t *testing.T) {
	projectWithDevcontainer(t)
	if err := os.Remove(filepath.Join(".devcontainer", "devcontainer.json")); err != nil {
		t.Fatal(err)
	}

	var out, errOut bytes.Buffer
	code := run(context.Background(), []string{"devcontainer", "up"}, &out, &errOut)
	if code != 10 || out.Len() != 0 || !strings.Contains(errOut.String(), ".devcontainer/devcontainer.json") {
		t.Errorf("exit %d, stdout=%q stderr=%q, want exit 10, no stdout and the file named", code, out.String(), errOut.String())
	}
}

// LH-FA-DEV-010: the group is wired into the real command tree.
func TestRun_Devcontainer_WiredAndBareExits2(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run(context.Background(), []string{"devcontainer"}, &out, &errOut); code != 2 {
		t.Errorf("bare devcontainer: exit %d, want 2", code)
	}
	out.Reset()
	if code := run(context.Background(), []string{"devcontainer", "--help"}, &out, &errOut); code != 0 ||
		!strings.Contains(out.String(), "up [--recreate]") || !strings.Contains(out.String(), "down ") || !strings.Contains(out.String(), "shell ") {
		t.Errorf("devcontainer --help: exit %d, stdout=%q", code, out.String())
	}
	out.Reset()
	if code := run(context.Background(), []string{"--help"}, &out, &errOut); code != 0 || !strings.Contains(out.String(), "devcontainer") {
		t.Errorf("--help lacks the devcontainer command: exit %d", code)
	}
}

// LH-FA-DEV-003, LH-FA-CLI-005 at the user path with the real adapters:
// `config set devcontainer.features.gh.enabled true` and the same for `glab`
// are accepted; `generate devcontainer` then writes both blocks and reports
// one supply-chain notice per feature in the default and in the sandbox
// profile; `--quiet` drops both notices, the blocks stay; a `version`
// override without `source` ends with exit 10.
func TestRun_Generate_ForgeCLINotices(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	const head = "schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n"
	for name, profile := range map[string]string{"default profile": "", "sandbox profile": "  profile: sandbox\n"} {
		t.Run(name, func(t *testing.T) {
			dir := writeFeatureProject(t, head+profile)
			for _, key := range []string{"gh", "glab"} {
				var out, errOut bytes.Buffer
				if code := run(context.Background(), []string{"config", "set", "devcontainer.features." + key + ".enabled", "true"}, &out, &errOut); code != 0 {
					t.Fatalf("config set %s: exit %d, stdout=%q stderr=%q", key, code, out.String(), errOut.String())
				}
			}
			var out, errOut bytes.Buffer
			if code := run(context.Background(), []string{"generate", "devcontainer"}, &out, &errOut); code != 0 {
				t.Fatalf("generate devcontainer: exit %d; stdout=%q stderr=%q", code, out.String(), errOut.String())
			}
			for _, want := range []string{
				"Warning (LH-NFA-SEC-004): devcontainer feature gh installs the release archive",
				"Warning (LH-NFA-SEC-004): devcontainer feature glab installs the release archive",
			} {
				if n := strings.Count(out.String(), want); n != 1 {
					t.Errorf("stdout contains %q %d times, want 1:\n%s", want, n, out.String())
				}
			}
			df, err := os.ReadFile(filepath.Join(dir, ".devcontainer", "Dockerfile"))
			if err != nil {
				t.Fatal(err)
			}
			for _, want := range []string{"ARG GH_VERSION=", "ARG GLAB_VERSION="} {
				if !strings.Contains(string(df), want) {
					t.Errorf("Dockerfile lacks %q", want)
				}
			}

			out.Reset()
			errOut.Reset()
			if code := run(context.Background(), []string{"--quiet", "generate", "devcontainer"}, &out, &errOut); code != 0 {
				t.Fatalf("--quiet generate devcontainer: exit %d; stdout=%q stderr=%q", code, out.String(), errOut.String())
			}
			if got := strings.ToLower(out.String() + errOut.String()); strings.Contains(got, "feature gh") || strings.Contains(got, "feature glab") {
				t.Errorf("--quiet must drop the notices of gh and glab:\nstdout=%q\nstderr=%q", out.String(), errOut.String())
			}
		})
	}

	for _, key := range []string{"gh", "glab"} {
		t.Run("version without source "+key, func(t *testing.T) {
			writeFeatureProject(t, head+"  features:\n    "+key+":\n      enabled: true\n      version: \"9.9.9\"\n")
			var out, errOut bytes.Buffer
			if code := run(context.Background(), []string{"generate", "devcontainer"}, &out, &errOut); code != 10 {
				t.Fatalf("generate devcontainer: exit %d, want 10; stdout=%q stderr=%q", code, out.String(), errOut.String())
			}
		})
	}
}

// LH-FA-DEV-008 at the user path: `egress list` names the origin of the hosts
// of the forge-CLI features. The two GitHub hosts of `gh` are in the base list
// as well; `gitlab.com` comes from `glab` alone.
func TestRun_EgressList_ForgeCLIHostOrigins(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	writeFeatureProject(t, "schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n  profile: sandbox\n"+
		"  features:\n    gh:\n      enabled: true\n    glab:\n      enabled: true\n"+
		"  sandbox:\n    egress:\n      enabled: true\n")
	var out, errOut bytes.Buffer
	if code := run(context.Background(), []string{"egress", "list", "--json"}, &out, &errOut); code != 0 {
		t.Fatalf("egress list: exit %d, stderr=%q", code, errOut.String())
	}
	var env struct {
		Data struct {
			Hosts []struct {
				Host    string   `json:"host"`
				Sources []string `json:"sources"`
			} `json:"hosts"`
		} `json:"data"`
	}
	if err := json.Unmarshal(out.Bytes(), &env); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	origins := map[string]string{}
	for _, h := range env.Data.Hosts {
		origins[h.Host] = strings.Join(h.Sources, ",")
	}
	for host, want := range map[string]string{
		"github.com": "base,feature:gh", "api.github.com": "base,feature:gh", "gitlab.com": "feature:glab", "codeload.github.com": "base",
	} {
		if origins[host] != want {
			t.Errorf("origin of %s = %q, want %q", host, origins[host], want)
		}
	}
}
