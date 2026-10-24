package application_test

import (
	"bytes"
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

const (
	dockerCLIKey = "ghcr.io/devcontainers/features/docker-outside-of-docker:1"

	sandboxFeaturesHead  = "devcontainer:\n  enabled: true\n  profile: sandbox\n  features:\n"
	sandboxDockerCLIYAML = sandboxFeaturesHead + "    docker-cli:\n      enabled: true\n"
	defaultDockerCLIYAML = "devcontainer:\n  enabled: true\n  features:\n    docker-cli:\n      enabled: true\n"
)

// requireWayOutIsDisablingTheFeature pins the way out that the message of the
// rejection (in `generate devcontainer` and in `doctor`) names, in this order:
// disabling the feature; the additional `nestedRuntime: docker`, named as a
// conditional step ("if you need Docker inside the container, also set …") for
// a Docker daemon of its own in the container; generating the devcontainer
// again; recreating an existing container. `nestedRuntime: docker` is never an
// alternative to disabling, because it alone does not lift the rejection
// (TestGenerateDevcontainer_SandboxNestedDocker_DoesNotLiftTheRejection).
func requireWayOutIsDisablingTheFeature(t *testing.T, message string) {
	t.Helper()
	if !strings.Contains(message, "docker-cli") {
		t.Errorf("message %q does not name the feature", message)
	}
	steps := []string{
		"disable the feature (`u-boot config set devcontainer.features.docker-cli.enabled false`)",
		"if you need Docker inside the container, also set `u-boot config set devcontainer.sandbox.nestedRuntime docker` " +
			"for a Docker daemon of its own in the container",
		"then generate the devcontainer again",
		"recreate an existing container (`u-boot devcontainer up --recreate`)",
	}
	at := make([]int, len(steps))
	for i, step := range steps {
		at[i] = strings.Index(message, step)
		if at[i] < 0 {
			t.Errorf("message %q lacks %q", message, step)
			return
		}
		if i > 0 && at[i] < at[i-1] {
			t.Errorf("message names %q before %q: %q", step, steps[i-1], message)
			return
		}
	}
	if strings.Contains(message[at[0]:at[1]], " or ") {
		t.Errorf("message offers nestedRuntime as an alternative to disabling the feature: %q", message)
	}
}

func generateDevcontainerDryRun(svc *application.GenerateService, sandbox bool) (driving.GenerateResponse, error) {
	return svc.Generate(context.Background(), driving.GenerateRequest{
		BaseDir: generateTestBaseDir, Artifact: domain.ArtifactDevcontainer, Sandbox: sandbox,
		PreviewMode: driving.PreviewDryRun,
	})
}

// LH-FA-DEV-006: a built-in feature that mounts the container-engine socket
// of the host (`docker-cli`) is a business error (exit 10) in the sandbox
// profile, whether the profile comes from u-boot.yaml or from `--sandbox`
// alone, with or without a `version` override and under `--dry-run`. The
// call writes nothing (no file, no directory, no `u-boot.yaml` change, so
// `--sandbox` does not persist the profile) and names the way out.
func TestGenerateDevcontainer_SandboxRejectsHostSocketFeature(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		yaml    string
		sandbox bool
		dryRun  bool
	}{
		{"profile in u-boot.yaml", sandboxDockerCLIYAML, false, false},
		{"profile from --sandbox only", defaultDockerCLIYAML, true, false},
		{"with a version override", sandboxDockerCLIYAML + "      version: \"2\"\n", false, false},
		{"profile in u-boot.yaml, dry run", sandboxDockerCLIYAML, false, true},
		{"profile from --sandbox only, dry run", defaultDockerCLIYAML, true, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, tc.yaml)
			seedGitOrigin(t, fs, "git@github.com:pt9912/demo.git")
			yamlPath := filepath.Join(generateTestBaseDir, "u-boot.yaml")
			yamlBefore, _ := fs.ReadFile(yamlPath)
			writesBefore, mkdirsBefore := fs.writtenPaths(), fs.mkdirPaths()

			var err error
			if tc.dryRun {
				_, err = generateDevcontainerDryRun(svc, tc.sandbox)
			} else {
				_, err = generateSandbox(svc, tc.sandbox)
			}
			if !errors.Is(err, domain.ErrInvalidSandboxSetting) {
				t.Fatalf("err = %v, want domain.ErrInvalidSandboxSetting", err)
			}
			requireWayOutIsDisablingTheFeature(t, err.Error())
			if got := fs.writtenPaths(); !reflect.DeepEqual(got, writesBefore) {
				t.Errorf("the rejected call wrote %v (before: %v)", got, writesBefore)
			}
			if got := fs.mkdirPaths(); !reflect.DeepEqual(got, mkdirsBefore) {
				t.Errorf("the rejected call created directories %v (before: %v)", got, mkdirsBefore)
			}
			if yamlAfter, _ := fs.ReadFile(yamlPath); !bytes.Equal(yamlAfter, yamlBefore) {
				t.Errorf("u-boot.yaml changed (the profile of --sandbox must not be persisted):\n%s", yamlAfter)
			}
			if ok, _ := fs.Exists(devcontainerJSONPath()); ok {
				t.Error("devcontainer.json exists after the rejected call")
			}
		})
	}
}

// LH-FA-DEV-006 / LH-FA-DEV-007: `nestedRuntime: docker` gives the container a
// Docker daemon of its own and does not lift the rejection of `docker-cli`:
// with both set, the call ends with the business error. With the feature
// switched off the same configuration generates: without an egress
// restriction `--privileged` is the only `runArgs` entry (with one,
// `--cap-add=NET_ADMIN` joins it), and there is no socket mount and no key of
// the feature.
func TestGenerateDevcontainer_SandboxNestedDocker_DoesNotLiftTheRejection(t *testing.T) {
	t.Parallel()
	const nestedDocker = "devcontainer:\n  enabled: true\n  profile: sandbox\n  sandbox:\n    nestedRuntime: docker\n  features:\n    docker-cli:\n      enabled: "
	t.Run("docker-cli active", func(t *testing.T) {
		t.Parallel()
		svc, fs := newGenerateService(t)
		seedUBootYAMLWithFeatures(t, fs, nestedDocker+"true\n")
		seedGitOrigin(t, fs, "git@github.com:pt9912/demo.git")
		_, err := generateSandbox(svc, false)
		if !errors.Is(err, domain.ErrInvalidSandboxSetting) {
			t.Fatalf("err = %v, want domain.ErrInvalidSandboxSetting", err)
		}
		if ok, _ := fs.Exists(devcontainerJSONPath()); ok {
			t.Error("devcontainer.json exists after the rejected call")
		}
	})
	t.Run("docker-cli switched off", func(t *testing.T) {
		t.Parallel()
		svc, fs := newGenerateService(t)
		seedUBootYAMLWithFeatures(t, fs, nestedDocker+"false\n")
		seedGitOrigin(t, fs, "git@github.com:pt9912/demo.git")
		if _, err := generateSandbox(svc, false); err != nil {
			t.Fatalf("generate: %v", err)
		}
		body, _ := fs.ReadFile(devcontainerJSONPath())
		m := devcontainerJSONMap(t, fs)
		if args, _ := m["runArgs"].([]any); len(args) != 1 || args[0] != "--privileged" {
			t.Errorf("runArgs = %v, want [--privileged]", m["runArgs"])
		}
		if bytes.Contains(body, []byte("docker.sock")) || bytes.Contains(body, []byte("docker-outside-of-docker")) {
			t.Errorf("devcontainer.json carries the socket or the feature:\n%s", body)
		}
	})
}

// LH-FA-DEV-006: the rejection is a rule of the sandbox profile. The
// standard profile keeps `docker-cli` selectable: exit 0, the key is
// rendered, and no warning names the socket or an option.
func TestGenerateDevcontainer_StandardProfileDockerCLI_Renders(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, defaultDockerCLIYAML)
	resp, err := generateSandbox(svc, false)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	body, _ := fs.ReadFile(devcontainerJSONPath())
	if _, ok := parseDevcontainerFeatures(t, body)[dockerCLIKey]; !ok {
		t.Errorf("the key %s is not rendered:\n%s", dockerCLIKey, body)
	}
	if len(resp.Warnings) != 0 {
		t.Errorf("warnings = %+v, want none for docker-cli in the standard profile", resp.Warnings)
	}
}

// LH-FA-DEV-006: only an active catalogue entry is rejected. A switched-off
// or unset entry renders nothing, and an entry with a `source:` override is
// an external feature: it is not `docker-cli`, the sandbox profile warns
// about it and renders it.
func TestGenerateDevcontainer_SandboxHostSocketFeature_OnlyActiveEntryIsRejected(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name         string
		yaml         string
		wantExternal int
	}{
		{"switched off", sandboxFeaturesHead + "    docker-cli:\n      enabled: false\n", 0},
		{"enabled unset", sandboxFeaturesHead + "    docker-cli: {}\n", 0},
		{"another entry", sandboxFeaturesHead + "    node:\n      enabled: true\n", 0},
		{
			"source override",
			"devcontainer:\n  enabled: true\n  profile: sandbox\n  featureSources:\n    allow:\n      - oci://registry.example.test/features/docker-outside-of-docker\n  features:\n    docker-cli:\n      enabled: true\n      source: oci://registry.example.test/features/docker-outside-of-docker\n",
			1,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, tc.yaml)
			seedGitOrigin(t, fs, "git@github.com:pt9912/demo.git")
			resp, err := generateSandbox(svc, false)
			if err != nil {
				t.Fatalf("generate: %v", err)
			}
			body, _ := fs.ReadFile(devcontainerJSONPath())
			if _, ok := parseDevcontainerFeatures(t, body)[dockerCLIKey]; ok {
				t.Errorf("the key %s is rendered:\n%s", dockerCLIKey, body)
			}
			external := 0
			for _, w := range resp.Warnings {
				if strings.Contains(w.Message, "comes from the external source") {
					external++
				}
			}
			if external != tc.wantExternal {
				t.Errorf("warnings about an external feature = %d, want %d: %+v", external, tc.wantExternal, resp.Warnings)
			}
		})
	}
}

// LH-FA-DEV-006 (Golden): the sandbox profile with `go` renders the key
// `…/go:1`, and `runArgs` and `mounts` equal those of the sandbox profile
// without any feature: the options of the feature are disclosed by two
// warnings and not written into the generated file.
func TestGenerateDevcontainer_SandboxGo_RendersKeyWithoutOptions(t *testing.T) {
	t.Parallel()
	render := func(t *testing.T, yaml string) (map[string]any, []byte, driving.GenerateResponse) {
		t.Helper()
		svc, fs := newGenerateService(t)
		seedUBootYAMLWithFeatures(t, fs, yaml)
		seedGitOrigin(t, fs, "git@github.com:pt9912/demo.git")
		resp, err := generateSandbox(svc, false)
		if err != nil {
			t.Fatalf("generate: %v", err)
		}
		body, _ := fs.ReadFile(devcontainerJSONPath())
		return devcontainerJSONMap(t, fs), body, resp
	}
	withGo, goBody, goResp := render(t, sandboxFeaturesHead+"    go:\n      enabled: true\n")
	without, _, _ := render(t, "devcontainer:\n  enabled: true\n  profile: sandbox\n")

	if _, ok := parseDevcontainerFeatures(t, goBody)["ghcr.io/devcontainers/features/go:1"]; !ok {
		t.Errorf("the key …/go:1 is not rendered:\n%s", goBody)
	}
	for _, key := range []string{"runArgs", "mounts"} {
		if !reflect.DeepEqual(withGo[key], without[key]) {
			t.Errorf("%s with go = %v, without any feature = %v", key, withGo[key], without[key])
		}
	}
	for _, option := range []string{"SYS_PTRACE", "seccomp"} {
		if bytes.Contains(goBody, []byte(option)) {
			t.Errorf("devcontainer.json carries %q:\n%s", option, goBody)
		}
	}
	flagged := 0
	for _, w := range goResp.Warnings {
		if w.SecurityRelaxation && strings.Contains(w.Message, "feature go") {
			flagged++
		}
	}
	if flagged != 2 {
		t.Errorf("relaxation warnings for go = %d, want 2: %+v", flagged, goResp.Warnings)
	}
}
