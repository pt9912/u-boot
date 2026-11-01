package application_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// Golden cases of `devcontainer.resources.*` (LH-FA-CONF-002,
// LH-FA-DEV-001): `--memory` and `--cpus` in `runArgs` of both profiles,
// byte-identical output without the keys, no security relaxation.

const (
	resourcesBoth   = "  resources:\n    memory: 3g\n    cpus: \"1.5\"\n"
	resourcesMemory = "  resources:\n    memory: 512m\n"
	resourcesCPUs   = "  resources:\n    cpus: \"2\"\n"
)

// R1: the default profile renders the options in the order memory, cpus;
// one key alone renders one option.
func TestGenerateDevcontainer_Resources_DefaultProfile(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		body string
		want []string
	}{
		"both":   {defaultProfile + resourcesBoth, []string{"--memory=3g", "--cpus=1.5"}},
		"memory": {defaultProfile + resourcesMemory, []string{"--memory=512m"}},
		"cpus":   {defaultProfile + resourcesCPUs, []string{"--cpus=2"}},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := runArgsOf(t, runHost(t, tc.body, false).fs); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("runArgs = %v, want %v", got, tc.want)
			}
		})
	}
}

// R2: the sandbox profile (also through `--sandbox`) renders the options
// behind the options of nestedRuntime and egress and before `ollama-host`.
func TestGenerateDevcontainer_Resources_SandboxProfile(t *testing.T) {
	t.Parallel()
	const defaultHead = "devcontainer:\n  enabled: true\n"
	lim := []string{"--memory=3g", "--cpus=1.5"}
	cases := map[string]struct {
		body    string
		sandbox bool
		want    []string
	}{
		"sandbox profile":      {sandboxProfile + resourcesBoth, false, lim},
		"sandbox flag":         {defaultHead + resourcesBoth, true, lim},
		"nested docker":        {sandboxProfile + "  sandbox:\n    nestedRuntime: docker\n" + resourcesBoth, false, append([]string{"--privileged"}, lim...)},
		"egress":               {sandboxProfile + egressOn + resourcesBoth, false, append([]string{"--cap-add=NET_ADMIN"}, lim...)},
		"podman":               {sandboxProfile + "  sandbox:\n    nestedRuntime: podman\n" + resourcesBoth, false, append(podmanOptions(), lim...)},
		"with ollama-host":     {sandboxProfile + ollamaHostOn + resourcesBoth, false, append(append([]string{}, lim...), hostLine)},
		"default with ollama":  {defaultProfile + ollamaHostOn + resourcesBoth, false, append(append([]string{}, lim...), hostLine)},
		"nested docker+egress": {sandboxProfile + "  sandbox:\n    nestedRuntime: docker\n    egress:\n      enabled: true\n" + resourcesMemory, false, []string{"--privileged", "--cap-add=NET_ADMIN", "--memory=512m"}},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := runArgsOf(t, runHost(t, tc.body, tc.sandbox).fs); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("runArgs = %v, want %v", got, tc.want)
			}
		})
	}
}

// R3: without the keys neither option appears and the file does not change
// (an empty `resources:` block counts as absent).
func TestGenerateDevcontainer_Resources_AbsentIsByteIdentical(t *testing.T) {
	t.Parallel()
	heads := []string{
		defaultProfile,
		defaultProfile + ollamaHostOn,
		sandboxProfile,
		sandboxProfile + "  sandbox:\n    nestedRuntime: docker\n",
		sandboxProfile + egressOn + ollamaHostOn,
	}
	for _, head := range heads {
		base := generatedFile(runHost(t, head, false).fs, ".devcontainer/devcontainer.json")
		empty := generatedFile(runHost(t, head+"  resources: {}\n", false).fs, ".devcontainer/devcontainer.json")
		if base != empty {
			t.Errorf("empty resources block changes the file for head %q", head)
		}
		blank := generatedFile(runHost(t, head+"  resources:\n    memory: \"\"\n    cpus: \"\"\n", false).fs, ".devcontainer/devcontainer.json")
		if base != blank {
			t.Errorf("empty-string resources change the file for head %q", head)
		}
		if strings.Contains(base, "--memory") || strings.Contains(base, "--cpus") {
			t.Errorf("limit option without the keys for head %q:\n%s", head, base)
		}
	}
}

// R4: the limits are not a security relaxation (LH-FA-DEV-006): no warning
// names them and the warnings equal those of the run without the keys.
func TestGenerateDevcontainer_Resources_NoRelaxationWarning(t *testing.T) {
	t.Parallel()
	for _, head := range []string{sandboxProfile + "  sandbox:\n    nestedRuntime: docker\n", defaultProfile} {
		with := runHost(t, head+resourcesBoth, false)
		without := runHost(t, head, false)
		if !reflect.DeepEqual(with.resp.Warnings, without.resp.Warnings) {
			t.Errorf("warnings differ: with %v, without %v", with.resp.Warnings, without.resp.Warnings)
		}
		for _, w := range with.resp.Warnings {
			if strings.Contains(w.Message, "--memory") || strings.Contains(w.Message, "--cpus") {
				t.Errorf("warning names a limit: %s", w.Message)
			}
		}
	}
}

// R5: a hand-edited bad value stops `generate` (exit-10 class) and names the key.
func TestGenerateDevcontainer_Resources_InvalidValue(t *testing.T) {
	t.Parallel()
	bad := []string{
		"  resources:\n    memory: 3G\n",
		"  resources:\n    memory: \"0\"\n",
		"  resources:\n    memory: \"3g\\\"\"\n",
		"  resources:\n    cpus: \"1,5\"\n",
		"  resources:\n    cpus: \"-1\"\n",
		"  resources:\n    cpus: \"2\\\\\"\n",
	}
	for _, res := range bad {
		svc, fs := newGenerateService(t)
		seedHostYAML(t, fs, defaultProfile+res)
		_, err := generateSandbox(svc, false)
		if !errors.Is(err, driving.ErrGenerateManualConflict) || !strings.Contains(err.Error(), "devcontainer.resources.") {
			t.Errorf("%q: err = %v, want ErrGenerateManualConflict naming devcontainer.resources.*", res, err)
		}
		if got := generatedFile(fs, ".devcontainer/devcontainer.json"); got != "" {
			t.Errorf("%q: devcontainer.json was written:\n%s", res, got)
		}
	}
}
