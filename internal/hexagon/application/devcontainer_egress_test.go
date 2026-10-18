package application_test

import (
	"context"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

func newEgressListService(fs *fakeFS) *application.EgressListService {
	return application.NewEgressListService(fs, &fakeYAML{}, nil)
}

func listEgress(t *testing.T, fs *fakeFS) driving.EgressListResponse {
	t.Helper()
	resp, err := newEgressListService(fs).List(context.Background(), driving.EgressListRequest{BaseDir: generateTestBaseDir})
	if err != nil {
		t.Fatalf("egress list: %v", err)
	}
	return resp
}

func egressHostNames(resp driving.EgressListResponse) []string {
	names := make([]string, 0, len(resp.Hosts))
	for _, h := range resp.Hosts {
		names = append(names, h.Host)
	}
	return names
}

func egressOrigins(resp driving.EgressListResponse) map[string][]string {
	out := map[string][]string{}
	for _, h := range resp.Hosts {
		out[h.Host] = h.Sources
	}
	return out
}

// egressConfig builds the `devcontainer:` body of a sandbox project with the
// egress restriction: features and sandbox keys are the YAML fragments
// (already indented below `devcontainer:` and `sandbox:`), allow the user list.
func egressConfig(features, sandboxKeys string, allow ...string) string {
	body := "devcontainer:\n  enabled: true\n  profile: sandbox\n"
	if features != "" {
		body += "  features:\n" + features
	}
	body += "  sandbox:\n" + sandboxKeys + "    egress:\n      enabled: true\n"
	if len(allow) > 0 {
		body += "      allow:\n"
		for _, h := range allow {
			body += "        - " + h + "\n"
		}
	}
	return body
}

const (
	featureNode       = "    node:\n      enabled: true\n"
	featureClaudeCode = "    claude-code:\n      enabled: true\n"
	featureHomebrew   = "    homebrew:\n      enabled: true\n"
	featureGoOff      = "    go:\n      enabled: false\n"
)

// LH-FA-DEV-008: `egress list` shows the list that `generate devcontainer`
// writes into the start script — the host names of the listing equal the
// `hosts=` line of `egress-init.sh` for the same project state. The cases
// cover the sources of a host (runtime, feature, clone, user) and the clone
// source states (origin, sandbox.repository, no remote, a `.git` file).
func TestEgressList_EqualsTheStartScriptList(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		yaml    string
		origin  string // "" = no .git/config
		gitFile bool   // a `.git` file (worktree/submodule) instead of `.git/config`
		clone   bool   // a clone host must contribute
	}{
		"default":                  {yaml: egressConfig("", ""), clone: false},
		"podman":                   {yaml: egressConfig("", "    nestedRuntime: podman\n"), clone: false},
		"docker":                   {yaml: egressConfig("", "    nestedRuntime: docker\n"), clone: false},
		"features":                 {yaml: egressConfig(featureNode+featureClaudeCode+featureHomebrew+featureGoOff, ""), clone: false},
		"features and podman":      {yaml: egressConfig(featureHomebrew, "    nestedRuntime: podman\n"), clone: false},
		"clone from origin":        {yaml: egressConfig("", ""), origin: "git@gitlab.example.org:o/r.git", clone: true},
		"clone from repository":    {yaml: egressConfig("", "    repository: https://git.example.net/o/r.git\n"), origin: "git@gitlab.example.org:o/r.git", clone: true},
		"user list":                {yaml: egressConfig("", "", "api.example.org", "ghcr.io", "api.example.org"), clone: false},
		"no remote":                {yaml: egressConfig(featureNode, "", "api.example.org"), clone: false},
		"git file instead of conf": {yaml: egressConfig("", ""), gitFile: true, clone: false},
		"everything": {yaml: egressConfig(featureNode+featureHomebrew, "    nestedRuntime: docker\n", "api.example.org", "ghcr.io"),
			origin: "git@gitlab.example.org:o/r.git", clone: true},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, tc.yaml)
			if tc.origin != "" {
				seedGitOrigin(t, fs, tc.origin)
			}
			if tc.gitFile {
				if err := fs.WriteFile(filepath.Join(generateTestBaseDir, ".git"), []byte("gitdir: ../elsewhere\n"), 0o644); err != nil {
					t.Fatalf("seed .git file: %v", err)
				}
			}
			if _, err := generateSandbox(svc, false); err != nil {
				t.Fatalf("generate: %v", err)
			}
			script := egressHostList(t, fs)
			resp := listEgress(t, fs)
			if got := egressHostNames(resp); !reflect.DeepEqual(got, script) {
				t.Errorf("listing differs from the start script:\n list  = %v\n script= %v", got, script)
			}
			if !resp.Active() {
				t.Errorf("Active() = false for a sandbox project with egress.enabled")
			}
			if resp.CloneSource != tc.clone {
				t.Errorf("CloneSource = %v, want %v", resp.CloneSource, tc.clone)
			}
			if hasClone := originsContain(resp, "clone"); hasClone != tc.clone {
				t.Errorf("a host with origin clone: %v, want %v", hasClone, tc.clone)
			}
		})
	}
}

func originsContain(resp driving.EgressListResponse, source string) bool {
	for _, h := range resp.Hosts {
		for _, s := range h.Sources {
			if s == source {
				return true
			}
		}
	}
	return false
}

// LH-FA-DEV-008: every host carries each origin once, in the fixed order
// base, runtime, docker-apt, feature (by key), clone, user — also for a host
// with several origins (`ghcr.io`, `registry.npmjs.org`, `github.com`).
func TestEgressList_OriginPerHost(t *testing.T) {
	t.Parallel()
	t.Run("docker, node and homebrew, clone and user list", func(t *testing.T) {
		t.Parallel()
		_, fs := newGenerateService(t)
		seedUBootYAMLWithFeatures(t, fs, egressConfig(featureNode+featureHomebrew, "    nestedRuntime: docker\n", "api.example.org", "ghcr.io"))
		seedGitOrigin(t, fs, "git@gitlab.example.org:o/r.git")
		want := map[string][]string{
			"api.example.org":                      {"user"},
			"api.github.com":                       {"base"},
			"auth.docker.io":                       {"runtime"},
			"codeload.github.com":                  {"base"},
			"deb.debian.org":                       {"base"},
			"download.docker.com":                  {"docker-apt"},
			"formulae.brew.sh":                     {"feature:homebrew"},
			"ghcr.io":                              {"runtime", "feature:homebrew", "user"},
			"github.com":                           {"base"},
			"gitlab.example.org":                   {"clone"},
			"objects.githubusercontent.com":        {"base"},
			"pkg-containers.githubusercontent.com": {"feature:homebrew"},
			"production.cloudflare.docker.com":     {"runtime"},
			"production.cloudfront.docker.com":     {"runtime"},
			"raw.githubusercontent.com":            {"base"},
			"registry-1.docker.io":                 {"runtime"},
			"registry.npmjs.org":                   {"feature:node"},
			"security.debian.org":                  {"base"},
		}
		if got := egressOrigins(listEgress(t, fs)); !reflect.DeepEqual(got, want) {
			t.Errorf("origins per host =\n%v\nwant\n%v", got, want)
		}
	})
	t.Run("one host from two features, clone host equal to a base host, duplicate user entry", func(t *testing.T) {
		t.Parallel()
		_, fs := newGenerateService(t)
		seedUBootYAMLWithFeatures(t, fs, egressConfig(featureNode+featureClaudeCode, "", "a.example.org", "a.example.org", "github.com"))
		seedGitOrigin(t, fs, "https://github.com/pt9912/demo.git")
		got := egressOrigins(listEgress(t, fs))
		for host, want := range map[string][]string{
			"registry.npmjs.org": {"feature:claude-code", "feature:node"},
			"api.anthropic.com":  {"feature:claude-code"},
			"github.com":         {"base", "clone", "user"},
			"a.example.org":      {"user"},
		} {
			if !reflect.DeepEqual(got[host], want) {
				t.Errorf("origins of %s = %v, want %v", host, got[host], want)
			}
		}
	})
	t.Run("a disabled feature contributes nothing", func(t *testing.T) {
		t.Parallel()
		_, fs := newGenerateService(t)
		seedUBootYAMLWithFeatures(t, fs, egressConfig(featureGoOff, ""))
		if _, ok := egressOrigins(listEgress(t, fs))["proxy.golang.org"]; ok {
			t.Errorf("a disabled go feature must not add proxy.golang.org")
		}
	})
}

// LH-FA-DEV-008: the listing states whether the restriction is active (sandbox
// profile and `egress.enabled`); when it is not, it shows the list that applies
// on activation — the same hosts in every state.
func TestEgressList_ActivationStateAndSameList(t *testing.T) {
	t.Parallel()
	const features = "  features:\n" + featureNode
	const allow = "      allow:\n        - api.example.org\n"
	cases := map[string]struct {
		yaml                     string
		profile, enabled, active bool
	}{
		"active":                  {"devcontainer:\n  enabled: true\n  profile: sandbox\n" + features + "  sandbox:\n    egress:\n      enabled: true\n" + allow, true, true, true},
		"egress without sandbox":  {"devcontainer:\n  enabled: true\n" + features + "  sandbox:\n    egress:\n      enabled: true\n" + allow, false, true, false},
		"sandbox, egress off":     {"devcontainer:\n  enabled: true\n  profile: sandbox\n" + features + "  sandbox:\n    egress:\n      enabled: false\n" + allow, true, false, false},
		"sandbox, egress not set": {"devcontainer:\n  enabled: true\n  profile: sandbox\n" + features + "  sandbox:\n    egress:\n" + allow, true, false, false},
		"neither":                 {"devcontainer:\n  enabled: true\n" + features + "  sandbox:\n    egress:\n" + allow, false, false, false},
	}
	var reference []string
	for _, name := range []string{"active", "egress without sandbox", "sandbox, egress off", "sandbox, egress not set", "neither"} {
		tc := cases[name]
		_, fs := newGenerateService(t)
		seedUBootYAMLWithFeatures(t, fs, tc.yaml)
		resp := listEgress(t, fs)
		if resp.SandboxProfile != tc.profile || resp.EgressEnabled != tc.enabled || resp.Active() != tc.active {
			t.Errorf("%s: profile=%v enabled=%v active=%v, want %v %v %v", name,
				resp.SandboxProfile, resp.EgressEnabled, resp.Active(), tc.profile, tc.enabled, tc.active)
		}
		names := egressHostNames(resp)
		if reference == nil {
			reference = names
		}
		if !reflect.DeepEqual(names, reference) {
			t.Errorf("%s: hosts differ from the active list:\n got  %v\n want %v", name, names, reference)
		}
	}
	if !sort.StringsAreSorted(reference) || len(reference) == 0 {
		t.Errorf("reference list must be sorted and not empty: %v", reference)
	}
}

// LH-FA-DEV-008: the listing reads and changes nothing: it writes no file
// and reads only u-boot.yaml and `.git/config`.
func TestEgressList_WritesNothingAndReadsOnlyTheProjectFiles(t *testing.T) {
	t.Parallel()
	_, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, egressConfig(featureNode, ""))
	seedGitOrigin(t, fs, "git@gitlab.example.org:o/r.git")
	before := len(fs.writes)
	listEgress(t, fs)
	if len(fs.writes) != before || len(fs.mkdirs) != 0 {
		t.Errorf("egress list wrote: writes %v, mkdirs %v", fs.writes[before:], fs.mkdirs)
	}
	for path, n := range fs.readFileCalls {
		if n > 0 && !strings.HasSuffix(path, "u-boot.yaml") && !strings.HasSuffix(path, filepath.Join(".git", "config")) {
			t.Errorf("egress list read %s", path)
		}
	}
}
