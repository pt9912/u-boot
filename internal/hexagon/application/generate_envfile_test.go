package application_test

import (
	"bytes"
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// Cases of the environment file of the devcontainer (LH-FA-DEV-012,
// ADR-0044): the line in `runArgs`, the creation of the file, the managed
// example file and the guarantee that an existing file is never touched.

const (
	envFileRel        = ".devcontainer/devcontainer.env"
	envFileExampleRel = ".devcontainer/devcontainer.env.example"
	envDryRunYAML     = "schemaVersion: 1\nproject:\n  name: t-uboot-gen\ndevcontainer:\n  enabled: true\n"
)

func envFilePath() string { return filepath.Join(generateTestBaseDir, envFileRel) }

func devEnvExampleFilePath() string { return filepath.Join(generateTestBaseDir, envFileExampleRel) }

// Happy: the standard profile carries the line as its only entry; with
// `resources` and `ollama-host` it stands first, followed by the resource
// options and the line of the feature.
func TestGenerateDevcontainer_EnvFile_RunArgsDefaultProfile(t *testing.T) {
	t.Parallel()
	only := runHost(t, defaultProfile, false)
	if got := fullRunArgsOf(t, only.fs); !reflect.DeepEqual(got, []string{envFileOption}) {
		t.Errorf("runArgs = %v, want exactly [%s]", got, envFileOption)
	}
	if raw := generatedFile(only.fs, ".devcontainer/devcontainer.json"); !strings.Contains(raw, envOnlyJSON) {
		t.Errorf("devcontainer.json lacks the literal runArgs line %q:\n%s", envOnlyJSON, raw)
	}
	all := runHost(t, defaultProfile+resourcesBoth+ollamaHostOn, false)
	want := []string{envFileOption, "--memory=3g", "--cpus=1.5", hostLine}
	if got := fullRunArgsOf(t, all.fs); !reflect.DeepEqual(got, want) {
		t.Errorf("runArgs = %v, want %v", got, want)
	}
}

// Boundary: in the sandbox profile the line stands before the options of
// `nestedRuntime`, the egress restriction, the resources and `ollama-host`.
func TestGenerateDevcontainer_EnvFile_RunArgsSandboxProfileOrder(t *testing.T) {
	t.Parallel()
	const egress = "    egress:\n      enabled: true\n"
	tail := resourcesBoth + ollamaHostOn
	cases := map[string]struct {
		yaml string
		want []string
	}{
		"podman": {
			sandboxProfile + "  sandbox:\n    nestedRuntime: podman\n" + egress + tail,
			append(append([]string{envFileOption}, podmanOptions()...), "--cap-add=NET_ADMIN", "--memory=3g", "--cpus=1.5", hostLine),
		},
		"docker": {
			sandboxProfile + "  sandbox:\n    nestedRuntime: docker\n" + egress + tail,
			[]string{envFileOption, "--privileged", "--cap-add=NET_ADMIN", "--memory=3g", "--cpus=1.5", hostLine},
		},
		"no runtime": {sandboxProfile, []string{envFileOption}},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := fullRunArgsOf(t, runHost(t, tc.yaml, false).fs); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("runArgs = %v, want %v", got, tc.want)
			}
		})
	}
}

// Happy: a second run changes nothing.
func TestGenerateDevcontainer_EnvFile_SecondRunNoOp(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)
	if _, err := generateDevcontainer(t, svc); err != nil {
		t.Fatalf("first run: %v", err)
	}
	before := projectTexts(fs)
	resp, err := generateDevcontainer(t, svc)
	if err != nil {
		t.Fatalf("second run: %v", err)
	}
	if resp.Action != driving.GenerateActionNoOp || len(resp.Changed) != 0 {
		t.Errorf("second run: action=%v changed=%v, want NoOp and nothing", resp.Action, resp.Changed)
	}
	if !reflect.DeepEqual(projectTexts(fs), before) {
		t.Error("second run changed a file")
	}
}

// Boundary: a deleted file is the only change of the next run; the action is
// `created`.
func TestGenerateDevcontainer_EnvFile_DeletedIsRecreated(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)
	if _, err := generateDevcontainer(t, svc); err != nil {
		t.Fatalf("first run: %v", err)
	}
	if err := fs.RemoveAll(envFilePath()); err != nil {
		t.Fatal(err)
	}
	resp, err := generateDevcontainer(t, svc)
	if err != nil {
		t.Fatalf("second run: %v", err)
	}
	if resp.Action != driving.GenerateActionCreated || !reflect.DeepEqual(resp.Changed, []string{envFileRel}) {
		t.Errorf("action=%v changed=%v, want Created and [%s]", resp.Action, resp.Changed, envFileRel)
	}
	if !dirExists(t, fs, envFilePath()) {
		t.Error("the file was not written again")
	}
}

// Happy: the created file names the format and the example file and holds
// only blank and comment lines: no assignment, no managed block.
func TestGenerateDevcontainer_EnvFile_CreatedWithCommentsOnly(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)
	if _, err := generateDevcontainer(t, svc); err != nil {
		t.Fatalf("generate: %v", err)
	}
	body := generatedFile(fs, envFileRel)
	if body == "" {
		t.Fatal("the environment file was not created")
	}
	for _, line := range strings.Split(body, "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			t.Errorf("line %q is neither blank nor a comment", line)
		}
	}
	if strings.Contains(body, "U-BOOT MANAGED BLOCK") {
		t.Errorf("the file carries a managed block marker:\n%s", body)
	}
	for _, want := range []string{"NAME=VALUE", "export", "devcontainer.env.example"} {
		if !strings.Contains(body, want) {
			t.Errorf("the file does not mention %q:\n%s", want, body)
		}
	}
}

// Negative: an existing file, with content and empty, stays byte-equal and
// is neither read, written nor backed up, in a first and a second run.
func TestGenerateDevcontainer_EnvFile_ExistingNeverTouched(t *testing.T) {
	t.Parallel()
	for name, content := range map[string]string{"with content": "TOKEN=abc\n# mine\n", "empty": ""} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedGenerateUbootYAML(t, fs)
			if err := fs.WriteFile(envFilePath(), []byte(content), 0o600); err != nil {
				t.Fatal(err)
			}
			writesBefore := len(fs.writes)
			for run := 1; run <= 2; run++ {
				resp, err := generateDevcontainer(t, svc)
				if err != nil {
					t.Fatalf("run %d: %v", run, err)
				}
				for _, c := range resp.Changed {
					if c == envFileRel {
						t.Errorf("run %d reports the existing file as changed: %v", run, resp.Changed)
					}
				}
			}
			if n := fs.readFileCallCount(envFilePath()); n != 0 {
				t.Errorf("the existing file was read %d times", n)
			}
			for _, p := range fs.writes[writesBefore:] {
				if p == envFilePath() || strings.HasPrefix(p, envFilePath()+".bak") {
					t.Errorf("the existing file or a copy of it was written: %s", p)
				}
			}
			if got := generatedFile(fs, envFileRel); got != content {
				t.Errorf("content = %q, want %q", got, content)
			}
			if mode := fs.fileModes[envFilePath()]; mode != 0o600 {
				t.Errorf("mode = %v, want 0600", mode)
			}
		})
	}
}

// Negative: when a later step fails, the newly created file is removed;
// a file of the user in the same directory stays.
func TestGenerateDevcontainer_EnvFile_RollbackRemovesCreatedFile(t *testing.T) {
	t.Parallel()
	fs := newFakeFS()
	fs.markDirExists(generateTestBaseDir)
	yamlPath := filepath.Join(generateTestBaseDir, "u-boot.yaml")
	if err := fs.WriteFile(yamlPath, []byte(rollbackSeedYAML), 0o644); err != nil {
		t.Fatal(err)
	}
	keep := filepath.Join(generateTestBaseDir, ".devcontainer", "keep.txt")
	if err := fs.WriteFile(keep, []byte("mine"), 0o644); err != nil {
		t.Fatal(err)
	}
	svc := application.NewGenerateService(&truncatingFailFS{fakeFS: fs, failPath: yamlPath}, &fakeYAML{}, nil)

	_, err := rollbackGen(svc)
	if !errors.Is(err, driving.ErrGenerateFileSystem) {
		t.Fatalf("err = %v, want ErrGenerateFileSystem", err)
	}
	for _, p := range []string{envFilePath(), devEnvExampleFilePath(), devcontainerJSONPath()} {
		if dirExists(t, fs, p) {
			t.Errorf("%s survived the rollback", p)
		}
	}
	if got, _ := fs.ReadFile(keep); string(got) != "mine" {
		t.Errorf("the file of the user changed: %q", got)
	}
}

// Boundary: `--dry-run` writes nothing; the list names the file only while
// it is missing.
func TestGenerateDevcontainer_EnvFile_DryRun(t *testing.T) {
	t.Parallel()
	prod := newFakeFS()
	prod.markDirExists(generateTestBaseDir)
	if err := prod.WriteFile(filepath.Join(generateTestBaseDir, "u-boot.yaml"), []byte(envDryRunYAML), 0o644); err != nil {
		t.Fatal(err)
	}
	svc := application.NewGenerateServiceWithFactory(configFactoryForProd(prod), &fakeYAML{}, nil)

	resp, err := generateDevcontainerMode(svc, driving.PreviewDryRun)
	if err != nil {
		t.Fatalf("dry-run: %v", err)
	}
	if !containsString(resp.Changed, envFileRel) {
		t.Errorf("dry-run list = %v, want it to name %s", resp.Changed, envFileRel)
	}
	if dirExists(t, prod, envFilePath()) || dirExists(t, prod, devEnvExampleFilePath()) {
		t.Error("dry-run wrote a file")
	}

	if err := prod.WriteFile(envFilePath(), []byte("A=1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	resp, err = generateDevcontainerMode(svc, driving.PreviewDryRun)
	if err != nil {
		t.Fatalf("dry-run with file: %v", err)
	}
	if containsString(resp.Changed, envFileRel) {
		t.Errorf("dry-run list = %v, must not name the existing file", resp.Changed)
	}
}

// Negative: the file causes no warning.
func TestGenerateDevcontainer_EnvFile_NoWarning(t *testing.T) {
	t.Parallel()
	for name, yaml := range map[string]string{"default": defaultProfile, "sandbox": sandboxProfile} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			for _, w := range runHost(t, yaml, false).resp.Warnings {
				if strings.Contains(w.Message, "devcontainer.env") || strings.Contains(w.Message, "env-file") {
					t.Errorf("warning because of the environment file: %+v", w)
				}
			}
		})
	}
}

// Happy + Boundary: the example file is a managed block with placeholder
// values only; a stale block is replaced and lines outside stay.
func TestGenerateDevcontainer_EnvFile_ExampleManagedBlock(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)
	if _, err := generateDevcontainer(t, svc); err != nil {
		t.Fatalf("generate: %v", err)
	}
	body := generatedFile(fs, envFileExampleRel)
	if !strings.Contains(body, "# BEGIN U-BOOT MANAGED BLOCK: init") || !strings.Contains(body, "# END U-BOOT MANAGED BLOCK: init") {
		t.Fatalf("the example file lacks the block markers:\n%s", body)
	}
	assignments := 0
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		assignments++
		if _, value, ok := strings.Cut(line, "="); !ok || !strings.HasPrefix(value, "CHANGEME_") {
			t.Errorf("line %q is not an assignment with a CHANGEME_ value", line)
		}
	}
	if assignments == 0 {
		t.Error("the example file holds no assignment")
	}

	stale := strings.Replace(body, "CHANGEME_API_TOKEN", "STALE", 1) + "MY_OWN=1\n"
	if err := fs.WriteFile(devEnvExampleFilePath(), []byte(stale), 0o644); err != nil {
		t.Fatal(err)
	}
	resp, err := generateDevcontainer(t, svc)
	if err != nil {
		t.Fatalf("regenerate: %v", err)
	}
	if resp.Action != driving.GenerateActionUpdatedBlock {
		t.Errorf("action = %v, want UpdatedBlock", resp.Action)
	}
	got := generatedFile(fs, envFileExampleRel)
	if strings.Contains(got, "STALE") || !strings.Contains(got, "CHANGEME_API_TOKEN") || !strings.HasSuffix(got, "MY_OWN=1\n") {
		t.Errorf("block not replaced or user line lost:\n%s", got)
	}
}

// Negative: `generate devcontainer` neither creates nor changes the
// `.gitignore`.
func TestGenerateDevcontainer_EnvFile_DoesNotWriteGitignore(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)
	if _, err := generateDevcontainer(t, svc); err != nil {
		t.Fatalf("generate: %v", err)
	}
	if dirExists(t, fs, filepath.Join(generateTestBaseDir, ".gitignore")) {
		t.Error(".gitignore was created")
	}
	mine := filepath.Join(generateTestBaseDir, ".gitignore")
	if err := fs.WriteFile(mine, []byte("node_modules/\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := fs.RemoveAll(envFilePath()); err != nil {
		t.Fatal(err)
	}
	if _, err := generateDevcontainer(t, svc); err != nil {
		t.Fatalf("regenerate: %v", err)
	}
	if got, _ := fs.ReadFile(mine); !bytes.Equal(got, []byte("node_modules/\n")) {
		t.Errorf(".gitignore changed: %q", got)
	}
}
