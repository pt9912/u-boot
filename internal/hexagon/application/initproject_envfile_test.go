package application_test

import (
	"context"
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// Cases of `init` for the environment file of the devcontainer
// (LH-FA-DEV-012): creation next to `generate devcontainer`, the
// guarantee for an existing file under the re-init flags, and the
// `.gitignore` line.

func initDevcontainerRunArgs(t *testing.T, fs *fakeFS) []string {
	t.Helper()
	body, err := fs.ReadFile(filepath.Join(testBaseDir, ".devcontainer", "devcontainer.json"))
	if err != nil {
		t.Fatalf("read devcontainer.json: %v", err)
	}
	var m struct {
		RunArgs []string `json:"runArgs"`
	}
	if err := json.Unmarshal(application.StripJSONCForTest(body), &m); err != nil {
		t.Fatalf("devcontainer.json invalid: %v", err)
	}
	return m.RunArgs
}

// Anti-drift: `init --devcontainer` and `generate devcontainer` write the
// same environment file, the same example file and the same `runArgs`.
func TestInit_Devcontainer_EnvFileCreatedAndRunArgsMatchGenerate(t *testing.T) {
	t.Parallel()
	svc, fs, _, _ := newService(t)
	resp, err := svc.Init(context.Background(), driving.InitProjectRequest{
		Name: "t-uboot-gen", BaseDir: testBaseDir, SkipGit: true, Devcontainer: true})
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	for _, rel := range []string{envFileRel, envFileExampleRel} {
		if !containsString(resp.Created, rel) {
			t.Errorf("Created = %v, want it to name %s", resp.Created, rel)
		}
	}
	gen := runHost(t, defaultProfile, false)
	if got, want := initDevcontainerRunArgs(t, fs), fullRunArgsOf(t, gen.fs); !reflect.DeepEqual(got, want) || len(got) == 0 || got[0] != envFileOption {
		t.Errorf("init runArgs = %v, generate runArgs = %v", got, want)
	}
	for _, rel := range []string{envFileRel, envFileExampleRel} {
		got, _ := fs.ReadFile(filepath.Join(testBaseDir, rel))
		if want := generatedFile(gen.fs, rel); string(got) != want || want == "" {
			t.Errorf("%s differs between init and generate:\n--- init ---\n%s\n--- generate ---\n%s", rel, got, want)
		}
	}
}

// Negative: with `--force --backup` an existing environment file does not
// abort the call, stays byte-equal, is not read and gets no copy.
func TestInit_Devcontainer_ExistingEnvFileKept(t *testing.T) {
	t.Parallel()
	for name, content := range map[string]string{"with content": "TOKEN=abc\n", "empty": ""} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs, _, _ := newService(t)
			envPath := filepath.Join(testBaseDir, envFileRel)
			if err := fs.WriteFile(envPath, []byte(content), 0o600); err != nil {
				t.Fatal(err)
			}
			resp, err := svc.Init(context.Background(), driving.InitProjectRequest{
				Name: "t-uboot-gen", BaseDir: testBaseDir, SkipGit: true, Devcontainer: true, Force: true, Backup: true})
			if err != nil {
				t.Fatalf("Init: %v", err)
			}
			if containsString(resp.Created, envFileRel) {
				t.Errorf("Created = %v, must not name the existing file", resp.Created)
			}
			for _, b := range resp.Backups {
				if b.Original == envFileRel {
					t.Errorf("backup of the environment file: %+v", b)
				}
			}
			if n := fs.readFileCallCount(envPath); n != 0 {
				t.Errorf("the existing file was read %d times", n)
			}
			if got, _ := fs.ReadFile(envPath); string(got) != content {
				t.Errorf("content = %q, want %q", got, content)
			}
			for p := range fs.files {
				if strings.HasPrefix(p, envPath+".bak") {
					t.Errorf("copy of the environment file: %s", p)
				}
			}
		})
	}
}

// Happy: the `.gitignore` of `init` names the environment file, with and
// without `--devcontainer`, and not the example file.
func TestInit_Gitignore_IgnoresEnvFileNotExample(t *testing.T) {
	t.Parallel()
	for name, devcontainer := range map[string]bool{"with devcontainer": true, "without devcontainer": false} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs, _, _ := newService(t)
			if _, err := svc.Init(context.Background(), driving.InitProjectRequest{
				Name: "t-uboot-gen", BaseDir: testBaseDir, SkipGit: true, Devcontainer: devcontainer}); err != nil {
				t.Fatalf("Init: %v", err)
			}
			body, _ := fs.ReadFile(filepath.Join(testBaseDir, ".gitignore"))
			lines := strings.Split(string(body), "\n")
			if !containsString(lines, envFileRel) {
				t.Errorf(".gitignore lacks the line %s:\n%s", envFileRel, body)
			}
			if containsString(lines, envFileExampleRel) || strings.Contains(string(body), ".example") {
				t.Errorf(".gitignore names the example file:\n%s", body)
			}
		})
	}
}
