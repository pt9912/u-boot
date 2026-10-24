package application_test

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// Golden cases of the catalogue key `ollama-host` (LH-FA-DEV-003,
// LH-FA-DEV-006, LH-FA-CLI-005): one `runArgs` line in both profiles and
// nothing else for the connection, one notice per run, no host in an egress
// list, no variable, no relay.

const (
	hostLine        = "--add-host=host.docker.internal:host-gateway"
	hostLineJSON    = `  "runArgs": ["` + hostLine + `"],` + "\n"
	ollamaHostOn    = "  features:\n    ollama-host:\n      enabled: true\n"
	hostNoticeStart = "devcontainer feature ollama-host adds the container option"
)

// seedHostYAML writes a u-boot.yaml without services (so no `forwardPorts`
// appears) and with the given `devcontainer:` sub-tree.
func seedHostYAML(t *testing.T, fs *fakeFS, devcontainerBody string) {
	t.Helper()
	body := "schemaVersion: 1\nproject:\n  name: t-uboot-gen\n" + devcontainerBody
	if err := fs.WriteFile(filepath.Join(generateTestBaseDir, "u-boot.yaml"), []byte(body), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
}

// hostNotices returns the notices of the connection (not the warning about a
// feature from an external source, which also names the feature).
func hostNotices(ws []driving.WarningEntry) []driving.WarningEntry {
	var out []driving.WarningEntry
	for _, w := range ws {
		if strings.Contains(w.Message, hostNoticeStart) {
			out = append(out, w)
		}
	}
	return out
}

// withoutHostNotices drops the notices of the connection from ws.
func withoutHostNotices(ws []driving.WarningEntry) []driving.WarningEntry {
	var out []driving.WarningEntry
	for _, w := range ws {
		if !strings.Contains(w.Message, hostNoticeStart) {
			out = append(out, w)
		}
	}
	return out
}

// runArgsOf returns the `runArgs` list of the generated devcontainer.json;
// nil when the key is absent.
func runArgsOf(t *testing.T, fs *fakeFS) []string {
	t.Helper()
	raw, ok := devcontainerJSONMap(t, fs)["runArgs"].([]any)
	if !ok {
		return nil
	}
	out := make([]string, len(raw))
	for i, v := range raw {
		out[i], _ = v.(string)
	}
	return out
}

// generatedFile returns the content of a file of the project, "" when absent.
func generatedFile(fs *fakeFS, rel string) string {
	raw, err := fs.ReadFile(filepath.Join(generateTestBaseDir, rel))
	if err != nil {
		return ""
	}
	return string(raw)
}

// projectTexts returns every file of the project directory by relative path.
func projectTexts(fs *fakeFS) map[string]string {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	out := map[string]string{}
	for p, b := range fs.files {
		if rel, err := filepath.Rel(generateTestBaseDir, p); err == nil && !strings.HasPrefix(rel, "..") {
			out[rel] = string(b)
		}
	}
	return out
}

// hostRun is one generate run of a fixture.
type hostRun struct {
	resp driving.GenerateResponse
	fs   *fakeFS
	svc  *application.GenerateService
}

func runHost(t *testing.T, yaml string, sandboxFlag bool) hostRun {
	t.Helper()
	svc, fs := newGenerateService(t)
	seedHostYAML(t, fs, yaml)
	resp, err := generateSandbox(svc, sandboxFlag)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	return hostRun{resp: resp, fs: fs, svc: svc}
}

// H1 (Z13, default profile): the line is the only entry of `runArgs`, between
// `build` and the following keys, valid JSONC; the file is the run without
// the feature plus that one line, character for character.
func TestGenerateDevcontainer_OllamaHost_DefaultProfileRunArgs(t *testing.T) {
	t.Parallel()
	with := runHost(t, defaultProfile+ollamaHostOn, false)
	without := runHost(t, defaultProfile, false)
	if got := runArgsOf(t, with.fs); !reflect.DeepEqual(got, []string{hostLine}) {
		t.Fatalf("runArgs = %v, want exactly [%s]", got, hostLine)
	}
	raw := generatedFile(with.fs, ".devcontainer/devcontainer.json")
	if !strings.Contains(raw, "  },\n"+hostLineJSON+"  \"remoteUser\"") {
		t.Errorf("the line must stand between build and remoteUser:\n%s", raw)
	}
	if got := strings.Replace(raw, hostLineJSON, "", 1); got != generatedFile(without.fs, ".devcontainer/devcontainer.json") {
		t.Errorf("devcontainer.json differs from the run without the feature by more than the line:\n%s", raw)
	}
	for _, rel := range []string{".devcontainer/Dockerfile"} {
		if generatedFile(with.fs, rel) != generatedFile(without.fs, rel) {
			t.Errorf("%s differs from the run without the feature", rel)
		}
	}
}

// H1: the line stands before `forwardPorts` and `features` when those keys
// exist (a service with a port, an OCI feature).
func TestGenerateDevcontainer_OllamaHost_DefaultProfileKeyOrder(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, defaultProfile+"  features:\n    node:\n      enabled: true\n    ollama-host:\n      enabled: true\n")
	if err := fs.WriteFile(filepath.Join(generateTestBaseDir, "compose.yaml"), []byte(composeYAMLPostgres), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := generateSandbox(svc, false); err != nil {
		t.Fatalf("generate: %v", err)
	}
	raw := generatedFile(fs, ".devcontainer/devcontainer.json")
	last := -1
	for _, key := range []string{`"build"`, `"runArgs"`, `"forwardPorts"`, `"features"`, `"remoteUser"`} {
		i := strings.Index(raw, key)
		if i < 0 || i < last {
			t.Fatalf("key %s missing or out of order (index %d, previous %d):\n%s", key, i, last, raw)
		}
		last = i
	}
	if got := runArgsOf(t, fs); !reflect.DeepEqual(got, []string{hostLine}) {
		t.Errorf("runArgs = %v, want exactly [%s]", got, hostLine)
	}
}

// H2 (Z13, sandbox profile and `--sandbox`): the line stands behind the
// options of the egress restriction and is the only entry of `runArgs`
// without it; `mounts`, `postStartCommand` and `remoteEnv` are those of the
// run without the feature.
func TestGenerateDevcontainer_OllamaHost_SandboxRunArgs(t *testing.T) {
	t.Parallel()
	const profileDefaultHead = "devcontainer:\n  enabled: true\n"
	cases := map[string]struct {
		head    string
		sandbox bool
		want    []string
	}{
		"sandbox profile":                {sandboxProfile, false, []string{hostLine}},
		"sandbox flag":                   {profileDefaultHead, true, []string{hostLine}},
		"sandbox profile with egress":    {sandboxProfile + egressOn, false, []string{"--cap-add=NET_ADMIN", hostLine}},
		"sandbox flag with egress":       {profileDefaultHead + egressOn, true, []string{"--cap-add=NET_ADMIN", hostLine}},
		"sandbox profile, podman":        {sandboxProfile + "  sandbox:\n    nestedRuntime: podman\n", false, append(podmanOptions(), hostLine)},
		"sandbox profile, nested docker": {sandboxProfile + "  sandbox:\n    nestedRuntime: docker\n", false, []string{"--privileged", hostLine}},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			with := runHost(t, tc.head+ollamaHostOn, tc.sandbox)
			without := runHost(t, tc.head, tc.sandbox)
			if got := runArgsOf(t, with.fs); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("runArgs = %v, want %v", got, tc.want)
			}
			a, b := devcontainerJSONMap(t, with.fs), devcontainerJSONMap(t, without.fs)
			delete(a, "runArgs")
			delete(b, "runArgs")
			if !reflect.DeepEqual(a, b) {
				t.Errorf("devcontainer.json differs from the run without the feature beyond runArgs:\n%v\n---\n%v", a, b)
			}
			if generatedFile(with.fs, ".devcontainer/Dockerfile") != generatedFile(without.fs, ".devcontainer/Dockerfile") {
				t.Error("Dockerfile differs from the run without the feature")
			}
		})
	}
}

// podmanOptions is the golden list of the options of nested Podman, in the
// order of `runArgs`.
func podmanOptions() []string {
	return []string{
		"--cap-add=SYS_ADMIN",
		"--security-opt=seccomp=unconfined",
		"--security-opt=apparmor=unconfined",
		"--security-opt=systempaths=unconfined",
		"--device=/dev/fuse",
		"--device=/dev/net/tun",
	}
}

// H3 (byte identity): an entry that is absent, `enabled: false` or without
// `enabled` produces the files of the run without the entry, in the default
// profile and in the sandbox profiles with and without the egress
// restriction.
func TestGenerateDevcontainer_OllamaHost_InactiveIsByteIdentical(t *testing.T) {
	t.Parallel()
	heads := map[string]string{
		"default":                 defaultProfile,
		"sandbox":                 sandboxProfile,
		"sandbox with egress":     sandboxProfile + egressOn,
		"default with UID":        defaultProfile + "  user:\n    uid: 501\n",
		"sandbox, nested podman":  sandboxProfile + "  sandbox:\n    nestedRuntime: podman\n",
		"sandbox, ollama and gh":  sandboxProfile + "  features:\n    ollama:\n      enabled: true\n    gh:\n      enabled: true\n",
		"default, gh, glab":       defaultProfile + "  features:\n    gh:\n      enabled: true\n    glab:\n      enabled: true\n",
		"default, node (feature)": defaultProfile + "  features:\n    node:\n      enabled: true\n",
	}
	entries := map[string]string{
		"absent":           "",
		"enabled false":    "    ollama-host:\n      enabled: false\n",
		"enabled missing":  "    ollama-host:\n      version: \"\"\n",
		"enabled null key": "    ollama-host: {}\n",
	}
	for hname, head := range heads {
		base := runHost(t, head, false)
		want := projectTexts(base.fs)
		for ename, entry := range entries {
			t.Run(hname+"/"+ename, func(t *testing.T) {
				t.Parallel()
				yaml := head
				if entry != "" {
					if strings.Contains(head, "  features:\n") {
						yaml = head + entry
					} else {
						yaml = head + "  features:\n" + entry
					}
				}
				got := runHost(t, yaml, false)
				files := projectTexts(got.fs)
				delete(files, "u-boot.yaml")
				wantFiles := map[string]string{}
				for k, v := range want {
					if k != "u-boot.yaml" {
						wantFiles[k] = v
					}
				}
				if !reflect.DeepEqual(files, wantFiles) {
					t.Errorf("generated files differ from the run without the entry:\n%v\n---\n%v", files, wantFiles)
				}
				if n := len(hostNotices(got.resp.Warnings)); n != 0 {
					t.Errorf("notices = %d, want 0 for an inactive entry", n)
				}
			})
		}
	}
}

// H4 (Z13 "nothing else", fitness function): none of the five names appears
// in any file a run writes, no `forwardPorts`, no `11434`, in either profile
// and with every feature that touches Ollama or the agent switched on.
func TestGenerateDevcontainer_OllamaHost_WritesNoVariableAndNoRelay(t *testing.T) {
	t.Parallel()
	const features = "  features:\n    ollama:\n      enabled: true\n    claude-code:\n      enabled: true\n    node:\n      enabled: true\n    ollama-host:\n      enabled: true\n"
	cases := map[string]struct {
		yaml    string
		sandbox bool
	}{
		"default profile":        {defaultProfile + features, false},
		"sandbox profile":        {sandboxProfile + features, false},
		"sandbox flag":           {defaultProfile + features, true},
		"sandbox, egress":        {sandboxProfile + features + "  sandbox:\n    egress:\n      enabled: true\n", false},
		"sandbox, podman":        {sandboxProfile + features + "  sandbox:\n    nestedRuntime: podman\n", false},
		"sandbox, nested docker": {sandboxProfile + features + "  sandbox:\n    nestedRuntime: docker\n", false},
	}
	names := []string{"OLLAMA_HOST", "ANTHROPIC_BASE_URL", "ANTHROPIC_AUTH_TOKEN", "OLLAMA_API_KEY", "OLLAMA_NO_CLOUD", "11434"}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			run := runHost(t, tc.yaml, tc.sandbox)
			files := projectTexts(run.fs)
			if len(files) < 3 {
				t.Fatalf("the run wrote %d files, want the project files: %v", len(files), files)
			}
			for rel, body := range files {
				for _, bad := range names {
					if strings.Contains(body, bad) {
						t.Errorf("%s contains %q", rel, bad)
					}
				}
			}
			if strings.Contains(generatedFile(run.fs, ".devcontainer/devcontainer.json"), "forwardPorts") {
				t.Error("devcontainer.json must carry no forwardPorts for the connection")
			}
		})
	}
}

// H5 (decision 3, fitness function): with `ollama-host` the egress list is
// the list without it; `ollama.com`, `host.docker.internal` and an address
// are in no list or script; the row `ollama` adds exactly its two names.
func TestGenerateDevcontainer_OllamaHost_AddsNoEgressHost(t *testing.T) {
	t.Parallel()
	const (
		manifestHost = "registry.ollama.ai"
		blobHost     = "dd20bb891979d25aebc8bec07b2b3bbc.r2.cloudflarestorage.com"
	)
	ip := regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$`)
	const others = "    ollama:\n      enabled: true\n    claude-code:\n      enabled: true\n    node:\n      enabled: true\n"
	head := sandboxProfile + egressOn + "  features:\n"
	with := runHost(t, head+others+"    ollama-host:\n      enabled: true\n", false)
	without := runHost(t, head+others, false)
	onlyHost := runHost(t, sandboxProfile+egressOn+ollamaHostOn, false)
	bare := runHost(t, sandboxProfile+egressOn, false)
	onlyOllama := runHost(t, sandboxProfile+egressOn+ollamaFeature, false)

	hosts := egressHostList(t, with.fs)
	if !sort.StringsAreSorted(hosts) || !reflect.DeepEqual(hosts, egressHostList(t, without.fs)) {
		t.Errorf("hosts with the feature =\n%v\nwant the sorted list without it\n%v", hosts, egressHostList(t, without.fs))
	}
	if !reflect.DeepEqual(egressHostList(t, onlyHost.fs), egressHostList(t, bare.fs)) {
		t.Errorf("ollama-host alone changes the list:\n%v\n%v", egressHostList(t, onlyHost.fs), egressHostList(t, bare.fs))
	}
	seen := map[string]bool{}
	for _, h := range hosts {
		if seen[h] {
			t.Errorf("host %s twice in %v", h, hosts)
		}
		seen[h] = true
		if h == "ollama.com" || strings.HasSuffix(h, ".ollama.com") || h == "host.docker.internal" || ip.MatchString(h) {
			t.Errorf("host %q must not be in the list", h)
		}
	}
	script := generatedFile(with.fs, ".devcontainer/egress-init.sh")
	for _, bad := range []string{"host.docker.internal", "ollama.com"} {
		if strings.Contains(script, bad) {
			t.Errorf("egress-init.sh contains %q", bad)
		}
	}
	// The row `ollama` adds exactly its two names to the base list.
	var added []string
	have := map[string]bool{}
	for _, h := range egressHostList(t, bare.fs) {
		have[h] = true
	}
	for _, h := range egressHostList(t, onlyOllama.fs) {
		if !have[h] {
			added = append(added, h)
		}
	}
	if want := []string{blobHost, manifestHost}; !reflect.DeepEqual(added, want) {
		t.Errorf("the row ollama adds %v, want exactly %v", added, want)
	}
}

// H6 (decisions 5 and 6): the credential variables of the sandbox profile are
// those of the other features; `ollama-host` and `ollama` add no variable.
func TestGenerateDevcontainer_OllamaHost_AddsNoCredentialVariable(t *testing.T) {
	t.Parallel()
	const all = "  features:\n    claude-code:\n      enabled: true\n    node:\n      enabled: true\n    gh:\n      enabled: true\n    glab:\n      enabled: true\n    ollama:\n      enabled: true\n"
	with := runHost(t, sandboxProfile+all+"    ollama-host:\n      enabled: true\n", false)
	remoteEnv, _ := devcontainerJSONMap(t, with.fs)["remoteEnv"].(map[string]any)
	var keys []string
	for k := range remoteEnv {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if want := []string{"ANTHROPIC_API_KEY", "GH_TOKEN", "GITLAB_TOKEN", "GIT_TOKEN"}; !reflect.DeepEqual(keys, want) {
		t.Errorf("remoteEnv keys = %v, want exactly %v", keys, want)
	}
	raw := generatedFile(with.fs, ".devcontainer/devcontainer.json")
	if !strings.Contains(raw, "\"GIT_TOKEN\": \"${localEnv:GIT_TOKEN}\",\n    \"ANTHROPIC_API_KEY\": \"${localEnv:ANTHROPIC_API_KEY}\",\n    \"GH_TOKEN\": \"${localEnv:GH_TOKEN}\",\n    \"GITLAB_TOKEN\": \"${localEnv:GITLAB_TOKEN}\"\n") {
		t.Errorf("remoteEnv is not the fixed list of references:\n%s", raw)
	}
}

// H6: `doctor` has no new finding for an active `ollama-host`: the line in
// `runArgs` triggers no credential pattern, and the checks of the features
// stay as without the entry.
func TestDoctor_OllamaHost_NoNewFinding(t *testing.T) {
	t.Parallel()
	const head = sandboxYAMLHead + "  profile: sandbox\n"
	const jsonWith = `{"remoteEnv": {"GIT_TOKEN": "${localEnv:GIT_TOKEN}"}, "runArgs": ["` + hostLine + `"]}`
	const jsonWithout = `{"remoteEnv": {"GIT_TOKEN": "${localEnv:GIT_TOKEN}"}}`
	type finding struct {
		ID       string
		Severity domain.Severity
	}
	collect := func(yaml, json string) ([]finding, map[string]string) {
		items := runSandboxDoctor(t, func(fs *fakeFS) {
			writeDoctorFile(t, fs, "u-boot.yaml", yaml)
			writeDoctorFile(t, fs, ".devcontainer/devcontainer.json", json)
		})
		var out []finding
		msgs := map[string]string{}
		for _, d := range items {
			out = append(out, finding{d.ID, d.Severity})
			msgs[d.ID] = d.Message
			if strings.Contains(d.Message, "ollama-host") || strings.Contains(d.Message, "host.docker.internal") {
				t.Errorf("finding %s names the connection: %q", d.ID, d.Message)
			}
		}
		return out, msgs
	}
	with, withMsgs := collect(head+ollamaHostOn, jsonWith)
	without, withoutMsgs := collect(head, jsonWithout)
	if !reflect.DeepEqual(with, without) {
		t.Errorf("doctor findings (id, severity) differ for an active ollama-host:\n%v\n---\n%v", with, without)
	}
	// The checks that read the generated file or the sandbox settings say
	// the same with and without the entry; the two feature checks
	// legitimately report that a feature is configured.
	for _, id := range []string{"devcontainer.sandbox.credentials", "devcontainer.sandbox.runtime", "devcontainer.sandbox.egress", "uboot.yaml.valid"} {
		if withMsgs[id] != withoutMsgs[id] {
			t.Errorf("%s: %q, want %q as without the entry", id, withMsgs[id], withoutMsgs[id])
		}
	}
}

// H7 (decision 8, sandbox profile): exactly one notice per run, flagged as a
// security relaxation, code LH-FA-DEV-006, on devcontainer.json.
func TestGenerateDevcontainer_OllamaHost_SandboxNoticeIsRelaxation(t *testing.T) {
	t.Parallel()
	for name, tc := range map[string]struct {
		yaml    string
		sandbox bool
	}{
		"sandbox profile": {sandboxProfile + ollamaHostOn, false},
		"sandbox flag":    {defaultProfile + ollamaHostOn, true},
		"with egress":     {sandboxProfile + egressOn + ollamaHostOn, false},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			notices := hostNotices(runHost(t, tc.yaml, tc.sandbox).resp.Warnings)
			if len(notices) != 1 {
				t.Fatalf("notices = %d, want 1: %+v", len(notices), notices)
			}
			w := notices[0]
			if w.Code != "LH-FA-DEV-006" || w.Level != "warn" || w.Subject != ".devcontainer/devcontainer.json" {
				t.Errorf("notice = %+v, want a warn entry with code LH-FA-DEV-006 on devcontainer.json", w)
			}
			if !w.SecurityRelaxation || w.DegradationFallback {
				t.Errorf("the notice must carry SecurityRelaxation and no fallback flag: %+v", w)
			}
		})
	}
}

// H8 (decision 8, default profile): one ordinary warning, dropped under
// `--quiet`, code LH-FA-DEV-003.
func TestGenerateDevcontainer_OllamaHost_DefaultNoticeIsOrdinary(t *testing.T) {
	t.Parallel()
	for name, yaml := range map[string]string{
		"default profile":        defaultProfile + ollamaHostOn,
		"egress flag is no-op":   defaultProfile + egressOn + ollamaHostOn,
		"default profile + UID":  defaultProfile + "  user:\n    uid: 501\n" + ollamaHostOn,
		"default profile + node": defaultProfile + "  features:\n    node:\n      enabled: true\n    ollama-host:\n      enabled: true\n",
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			notices := hostNotices(runHost(t, yaml, false).resp.Warnings)
			if len(notices) != 1 {
				t.Fatalf("notices = %d, want 1: %+v", len(notices), notices)
			}
			w := notices[0]
			if w.Code != "LH-FA-DEV-003" || w.Level != "warn" || w.Subject != ".devcontainer/devcontainer.json" {
				t.Errorf("notice = %+v, want a warn entry with code LH-FA-DEV-003 on devcontainer.json", w)
			}
			if w.SecurityRelaxation || w.DegradationFallback {
				t.Errorf("the notice must be unflagged (dropped under --quiet): %+v", w)
			}
		})
	}
}

// H9 (O2, decision 8): profile times egress restriction. Always one notice,
// exit 0, the same action as without the entry; the sentence that the
// connection does not work behind the restriction only with the sandbox
// profile and an active restriction; no address in the allowlist.
func TestGenerateDevcontainer_OllamaHost_EgressRestrictionIsAWarning(t *testing.T) {
	t.Parallel()
	const behind = "this connection does not work behind it"
	cases := map[string]struct {
		head       string
		wantBehind bool
	}{
		"sandbox, restriction":    {sandboxProfile + egressOn, true},
		"sandbox, no restriction": {sandboxProfile, false},
		"default, flag set":       {defaultProfile + egressOn, false},
		"default":                 {defaultProfile, false},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			with := runHost(t, tc.head+ollamaHostOn, false)
			without := runHost(t, tc.head, false)
			notices := hostNotices(with.resp.Warnings)
			if len(notices) != 1 {
				t.Fatalf("notices = %d, want exactly 1: %+v", len(notices), notices)
			}
			if got := strings.Contains(notices[0].Message, behind); got != tc.wantBehind {
				t.Errorf("the sentence %q present = %v, want %v:\n%s", behind, got, tc.wantBehind, notices[0].Message)
			}
			if tc.wantBehind {
				if !notices[0].SecurityRelaxation || !strings.Contains(notices[0].Message, "adds no address to the allowlist") {
					t.Errorf("the one notice must stay a relaxation and name the allowlist: %+v", notices[0])
				}
			}
			if with.resp.Action != without.resp.Action {
				t.Errorf("action = %v, want %v as without the entry", with.resp.Action, without.resp.Action)
			}
			if !reflect.DeepEqual(withoutHostNotices(with.resp.Warnings), without.resp.Warnings) {
				t.Errorf("the other warnings change with the entry:\n%+v\n---\n%+v", withoutHostNotices(with.resp.Warnings), without.resp.Warnings)
			}
			if generatedFile(with.fs, ".devcontainer/egress-init.sh") != generatedFile(without.fs, ".devcontainer/egress-init.sh") {
				t.Error("egress-init.sh differs from the run without the entry")
			}
		})
	}
}

// H10 (O1, new entry type): the entry renders no key under `features`, alone
// or next to OCI features; the warnings of other predicates do not change
// with it.
func TestCollectDevcontainerFeatures_OllamaHost_RendersNoKey(t *testing.T) {
	t.Parallel()
	const head = "schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n  features:\n"
	if got := application.CollectDevcontainerFeaturesForTest(t, []byte(head+"    ollama-host:\n      enabled: true\n")); len(got) != 0 {
		t.Errorf("features = %v, want none (a key \":\" would be rendered)", got)
	}
	const oci = "    claude-code:\n      enabled: true\n    node:\n      enabled: true\n"
	with := application.CollectDevcontainerFeaturesForTest(t, []byte(head+oci+"    ollama-host:\n      enabled: true\n"))
	without := application.CollectDevcontainerFeaturesForTest(t, []byte(head+oci))
	if !reflect.DeepEqual(with, without) || len(with) != 2 {
		t.Errorf("features with the entry = %v, want %v (two OCI keys)", with, without)
	}
}

// H10: no JSON key, no orphan, no drift: `doctor` is silent for the entry with
// an empty `features` in the devcontainer.json.
func TestDoctor_FeaturesDrift_OllamaHost(t *testing.T) {
	t.Parallel()
	for name, jsonBody := range map[string]string{
		"features empty": `{"features": {}}`,
		"no features":    `{"name": "demo"}`,
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs, _, _, _ := newDoctorService(t)
			driftSeed(t, fs, "schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n"+ollamaHostOn, jsonBody)
			resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
			if err != nil {
				t.Fatalf("Check: %v", err)
			}
			for _, id := range []string{"devcontainer.features.drift", "devcontainer.features.allowlist"} {
				d := findDiagnostic(t, resp.Report.Items, id)
				if d.Severity != domain.SeverityOK {
					t.Errorf("%s: Severity = %v, want OK; message = %q", id, d.Severity, d.Message)
				}
			}
			for _, d := range resp.Report.Items {
				if strings.Contains(strings.ToLower(d.Message), "orphan") && strings.Contains(d.Message, "ollama-host") {
					t.Errorf("%s reports the entry as an orphan: %q", d.ID, d.Message)
				}
			}
		})
	}
}

// H10: the predicates of the other entries do not react to `ollama-host`: the
// warnings of a run with a login service, a feature with options and a
// third-party feature are those of the run without the entry (apart from the
// one notice).
func TestGenerateDevcontainer_OllamaHost_OtherPredicatesUnchanged(t *testing.T) {
	t.Parallel()
	const others = "  features:\n    sshd:\n      enabled: true\n    go:\n      enabled: true\n    homebrew:\n      enabled: true\n    claude-code:\n      enabled: true\n"
	for pname, head := range map[string]string{"default profile": defaultProfile, "sandbox profile": sandboxProfile} {
		t.Run(pname, func(t *testing.T) {
			t.Parallel()
			with := runHost(t, head+others+"    ollama-host:\n      enabled: true\n", false)
			without := runHost(t, head+others, false)
			if !reflect.DeepEqual(withoutHostNotices(with.resp.Warnings), without.resp.Warnings) {
				t.Errorf("warnings change with the entry:\n%+v\n---\n%+v", withoutHostNotices(with.resp.Warnings), without.resp.Warnings)
			}
			if generatedFile(with.fs, ".devcontainer/Dockerfile") != generatedFile(without.fs, ".devcontainer/Dockerfile") {
				t.Error("Dockerfile differs from the run without the entry")
			}
		})
	}
}

// H11 (O1, reading): `version:` without `source:` ends with exit code 10
// (ErrGenerateManualConflict), writes nothing, and the message names the key
// and says that the entry is no OCI feature (not "installation in the
// Dockerfile"); the doctor reports the same case as an error of the file.
func TestGenerateDevcontainer_OllamaHost_VersionWithoutSourceRejected(t *testing.T) {
	t.Parallel()
	const entry = "  features:\n    ollama-host:\n      enabled: true\n      version: \"1\"\n"
	for pname, head := range map[string]string{"default profile": defaultProfile, "sandbox profile": sandboxProfile} {
		t.Run(pname, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedHostYAML(t, fs, head+entry)
			writes := len(fs.writtenPaths())
			_, err := generateSandbox(svc, false)
			if !errors.Is(err, driving.ErrGenerateManualConflict) || err == nil {
				t.Fatalf("generate error = %v, want a wrap of ErrGenerateManualConflict (exit code 10)", err)
			}
			for _, part := range []string{"features.ollama-host.version", "no OCI feature"} {
				if !strings.Contains(err.Error(), part) {
					t.Errorf("error lacks %q: %v", part, err)
				}
			}
			if strings.Contains(err.Error(), "Dockerfile") {
				t.Errorf("error must not speak of an installation in the Dockerfile: %v", err)
			}
			if exists, _ := fs.Exists(dockerfilePath()); exists || len(fs.writtenPaths()) != writes {
				t.Errorf("nothing may be written on a validation error: %v", fs.writtenPaths())
			}
		})
	}
	t.Run("doctor", func(t *testing.T) {
		t.Parallel()
		svc, fs, _, _, _ := newDoctorService(t)
		seedDoctorUbootYAMLFeatures(t, fs, "schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n"+entry)
		resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
		if err != nil {
			t.Fatalf("Check: %v", err)
		}
		d := findDiagnostic(t, resp.Report.Items, "uboot.yaml.valid")
		if d.Severity != domain.SeverityError || !strings.Contains(d.Message, "features.ollama-host.version") {
			t.Errorf("uboot.yaml.valid = %v %q, want an error naming the key", d.Severity, d.Message)
		}
	})
}

// H11 (write path): `config set …ollama-host.version` is rejected like the
// read path (ErrConfigValueInvalid, exit code 10), the file stays
// byte-identical in every preview mode, the text names the way out and says
// that the entry is no OCI feature.
func TestConfigSet_OllamaHostVersionWithoutSource_Rejected(t *testing.T) {
	t.Parallel()
	for _, mode := range []struct {
		name string
		mode driving.PreviewMode
	}{{"apply", driving.PreviewNone}, {"dry-run", driving.PreviewDryRun}, {"diff", driving.PreviewAndApply}} {
		t.Run(mode.name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newConfigService(t)
			seedConfigUbootYAMLWithDevcontainer(t, fs, "schemaVersion: 1\nproject:\n  name: t-uboot-config\n"+
				"devcontainer:\n  enabled: true\n"+ollamaHostOn)
			before, writes := readYAML(t, fs), len(fs.writes)
			_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
				BaseDir:     configTestBaseDir,
				Path:        mustConfigPath(t, "devcontainer.features.ollama-host.version"),
				Value:       "2",
				PreviewMode: mode.mode,
			})
			if !errors.Is(err, driving.ErrConfigValueInvalid) || err == nil {
				t.Fatalf("Set error = %v, want a wrap of ErrConfigValueInvalid", err)
			}
			for _, part := range []string{"devcontainer.features.ollama-host.source", "no OCI feature"} {
				if !strings.Contains(err.Error(), part) {
					t.Errorf("Set error lacks %q: %v", part, err)
				}
			}
			if strings.Contains(err.Error(), "Dockerfile") {
				t.Errorf("Set error must not speak of an installation in the Dockerfile: %v", err)
			}
			if readYAML(t, fs) != before || len(fs.writes) != writes {
				t.Errorf("u-boot.yaml changed by a rejected Set:\n%s", readYAML(t, fs))
			}
		})
	}
}

// H11 (`source:`): a `source:` override makes the entry an external OCI
// feature without the line: the key is rendered, `runArgs` carries no line,
// the sandbox profile gives the one warning about an external source and no
// notice of the connection.
func TestGenerateDevcontainer_OllamaHost_SourceOverrideIsExternalFeature(t *testing.T) {
	t.Parallel()
	const external = "ghcr.io/example/features/ollama-host"
	features := "  featureSources:\n    allow:\n      - https://" + external + "\n" +
		"  features:\n    ollama-host:\n      enabled: true\n      source: https://" + external + "\n"
	for pname, head := range map[string]string{"default profile": defaultProfile, "sandbox profile": sandboxProfile} {
		t.Run(pname, func(t *testing.T) {
			t.Parallel()
			run := runHost(t, head+features, false)
			raw := generatedFile(run.fs, ".devcontainer/devcontainer.json")
			if strings.Contains(raw, "host.docker.internal") || len(runArgsOf(t, run.fs)) != 0 {
				t.Errorf("a source override must not add the line:\n%s", raw)
			}
			if !strings.Contains(raw, `"https://`+external+`:1": {}`) {
				t.Errorf("the external feature is not rendered:\n%s", raw)
			}
			if n := len(hostNotices(run.resp.Warnings)); n != 0 {
				t.Errorf("notices of the connection = %d, want 0", n)
			}
			fromExternal := 0
			for _, w := range run.resp.Warnings {
				if strings.Contains(w.Message, "feature ollama-host comes from the external source") {
					fromExternal++
				}
			}
			if want := map[string]int{"default profile": 0, "sandbox profile": 1}[pname]; fromExternal != want {
				t.Errorf("warnings about the external source = %d, want %d: %+v", fromExternal, want, run.resp.Warnings)
			}
		})
	}
}

// H12 (Z13 "nothing without the wish of the user"): `enabled: false`, a
// missing `enabled` and a missing entry give no line and no notice, in both
// profiles.
func TestGenerateDevcontainer_OllamaHost_NotActive(t *testing.T) {
	t.Parallel()
	for pname, head := range map[string]string{"default profile": defaultProfile, "sandbox profile": sandboxProfile} {
		for ename, features := range map[string]string{
			"disabled":        "  features:\n    ollama-host:\n      enabled: false\n",
			"enabled missing": "  features:\n    ollama-host: {}\n",
			"absent":          "",
		} {
			t.Run(pname+"/"+ename, func(t *testing.T) {
				t.Parallel()
				run := runHost(t, head+features, false)
				if got := runArgsOf(t, run.fs); len(got) != 0 {
					t.Errorf("runArgs = %v, want none", got)
				}
				if strings.Contains(generatedFile(run.fs, ".devcontainer/devcontainer.json"), "host.docker.internal") {
					t.Error("devcontainer.json names host.docker.internal for an inactive entry")
				}
				for _, w := range run.resp.Warnings {
					if strings.Contains(w.Message, "ollama-host") {
						t.Errorf("warning for an inactive entry: %+v", w)
					}
				}
			})
		}
	}
}

// H13 (replay): the second run is a NoOp and repeats the notice; enabling the
// entry later replaces the stale block and keeps the line of the user behind
// it; disabling it removes the option again.
func TestGenerateDevcontainer_OllamaHost_ReplayAndLaterActivation(t *testing.T) {
	t.Parallel()
	for pname, head := range map[string]string{"default profile": defaultProfile, "sandbox profile": sandboxProfile} {
		t.Run(pname, func(t *testing.T) {
			t.Parallel()
			const userLine = "// user addition\n"
			svc, fs := newGenerateService(t)
			seedHostYAML(t, fs, head)
			if _, err := generateSandbox(svc, false); err != nil {
				t.Fatalf("first generate: %v", err)
			}
			body := generatedFile(fs, ".devcontainer/devcontainer.json")
			if err := fs.WriteFile(devcontainerJSONPath(), []byte(body+userLine), 0o644); err != nil {
				t.Fatal(err)
			}

			seedHostYAML(t, fs, head+ollamaHostOn)
			resp, err := generateSandbox(svc, false)
			if err != nil {
				t.Fatalf("activate: %v", err)
			}
			if resp.Action != driving.GenerateActionUpdatedBlock {
				t.Errorf("activate: action = %v, want UpdatedBlock", resp.Action)
			}
			after := generatedFile(fs, ".devcontainer/devcontainer.json")
			if !strings.Contains(after, hostLine) || !strings.HasSuffix(after, userLine) {
				t.Errorf("the line must be in the block and the user line behind it:\n%s", after)
			}
			r2, err := generateSandbox(svc, false)
			if err != nil || r2.Action != driving.GenerateActionNoOp {
				t.Errorf("second run: action=%v err=%v, want NoOp", r2.Action, err)
			}
			if n := len(hostNotices(r2.Warnings)); n != 1 {
				t.Errorf("the NoOp run gives %d notices, want 1 (it arises at every run)", n)
			}

			seedHostYAML(t, fs, head+"  features:\n    ollama-host:\n      enabled: false\n")
			if _, err := generateSandbox(svc, false); err != nil {
				t.Fatalf("deactivate: %v", err)
			}
			gone := generatedFile(fs, ".devcontainer/devcontainer.json")
			if strings.Contains(gone, "host.docker.internal") || !strings.HasSuffix(gone, userLine) {
				t.Errorf("the option must be gone and the user line kept:\n%s", gone)
			}
		})
	}
}

// H14 (interplay): `ollama`, `gh`, `glab` and `ollama-host` with each nested
// runtime and the egress restriction: `runArgs` character for character, the
// mounts and the Dockerfile as without the entry, the hosts the sorted union,
// the second run a NoOp.
func TestGenerateDevcontainer_OllamaHost_TogetherWithOtherFeaturesAndNestedRuntime(t *testing.T) {
	t.Parallel()
	const others = "  features:\n    ollama:\n      enabled: true\n    gh:\n      enabled: true\n    glab:\n      enabled: true\n"
	const hostEntry = "    ollama-host:\n      enabled: true\n"
	cases := map[string]struct {
		sandboxSettings string
		wantRunArgs     string
	}{
		"podman": {"  sandbox:\n    nestedRuntime: podman\n    egress:\n      enabled: true\n",
			`  "runArgs": ["--cap-add=SYS_ADMIN", "--security-opt=seccomp=unconfined", "--security-opt=apparmor=unconfined", "--security-opt=systempaths=unconfined", "--device=/dev/fuse", "--device=/dev/net/tun", "--cap-add=NET_ADMIN", "` + hostLine + `"],` + "\n"},
		"docker": {"  sandbox:\n    nestedRuntime: docker\n    egress:\n      enabled: true\n",
			`  "runArgs": ["--privileged", "--cap-add=NET_ADMIN", "` + hostLine + `"],` + "\n"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			head := sandboxProfile + tc.sandboxSettings
			with := runHost(t, head+others+hostEntry, false)
			without := runHost(t, head+others, false)
			if raw := generatedFile(with.fs, ".devcontainer/devcontainer.json"); !strings.Contains(raw, tc.wantRunArgs) {
				t.Errorf("devcontainer.json lacks\n%s\nin:\n%s", tc.wantRunArgs, raw)
			}
			a, b := devcontainerJSONMap(t, with.fs), devcontainerJSONMap(t, without.fs)
			if !reflect.DeepEqual(a["mounts"], b["mounts"]) {
				t.Errorf("mounts differ: %v vs %v", a["mounts"], b["mounts"])
			}
			if generatedFile(with.fs, ".devcontainer/Dockerfile") != generatedFile(without.fs, ".devcontainer/Dockerfile") {
				t.Error("Dockerfile differs from the run without the entry")
			}
			hosts := egressHostList(t, with.fs)
			if !sort.StringsAreSorted(hosts) || !reflect.DeepEqual(hosts, egressHostList(t, without.fs)) {
				t.Errorf("hosts = %v, want the sorted list of the run without the entry %v", hosts, egressHostList(t, without.fs))
			}
			if r2, err := generateSandbox(with.svc, false); err != nil || r2.Action != driving.GenerateActionNoOp {
				t.Errorf("second run: action=%v err=%v, want NoOp", r2.Action, err)
			}
		})
	}
}

// H15 (decision 8): the text is neutral and limited to what is measured. It
// names the host name, the condition of the reach, the relay of the user and
// that u-boot starts none; it names neither the cloud service nor models, no
// identifier, no version and no date, and states no unconditional reach.
func TestGenerateDevcontainer_OllamaHost_NoticeText(t *testing.T) {
	t.Parallel()
	digits := regexp.MustCompile(`[0-9]`)
	for name, tc := range map[string]struct {
		yaml       string
		wantBehind bool
	}{
		"default profile":          {defaultProfile + ollamaHostOn, false},
		"sandbox profile":          {sandboxProfile + ollamaHostOn, false},
		"sandbox profile + egress": {sandboxProfile + egressOn + ollamaHostOn, true},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			notices := hostNotices(runHost(t, tc.yaml, false).resp.Warnings)
			if len(notices) != 1 {
				t.Fatalf("notices = %d, want 1", len(notices))
			}
			msg := notices[0].Message
			for _, must := range []string{
				"host.docker.internal",
				"only through a relay or another binding that you set up",
				"loopback address",
				"u-boot starts no relay",
				"depends on the address it listens on",
			} {
				if !strings.Contains(msg, must) {
					t.Errorf("text lacks %q:\n%s", must, msg)
				}
			}
			low := strings.ToLower(msg)
			for _, bad := range []string{"ollama.com", "cloud", "lh-", "adr-", "model", "quota", "sign in"} {
				if strings.Contains(low, bad) {
					t.Errorf("text must not contain %q:\n%s", bad, msg)
				}
			}
			if digits.MatchString(strings.NewReplacer("ollama-host", "", "host-gateway", "").Replace(low)) {
				t.Errorf("text must carry no number (version, date, address):\n%s", msg)
			}
			if strings.Contains(low, "is reachable.") || strings.Contains(low, "is always reachable") {
				t.Errorf("text states an unconditional reach:\n%s", msg)
			}
			if got := strings.Contains(msg, "does not work behind it"); got != tc.wantBehind {
				t.Errorf("egress sentence present = %v, want %v", got, tc.wantBehind)
			}
			if got := strings.Contains(msg, "limits damage and is not an isolation boundary"); got != strings.Contains(name, "sandbox") {
				t.Errorf("profile clause present = %v for %s", got, name)
			}
		})
	}
}
