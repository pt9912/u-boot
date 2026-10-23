package application_test

import (
	"errors"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// Golden cases of the catalogue keys `gh` and `glab` (LH-FA-DEV-003,
// LH-FA-DEV-009): a Dockerfile block with a pinned release archive in both
// profiles, no OCI render key, no volume, one supply-chain notice per active
// feature, the hosts in the egress list and, in the sandbox profile, the
// reference of the token variable in `remoteEnv`.

const (
	ghSHAAMD64   = "fdd768eb70f6409f80fcfa9f3a1c0adf0fc98d827ab3ce0bed2ac22762c5d786"
	ghSHAARM64   = "7f322df34e624d14f5ec0c4acb121fd6dfa14392fdbbd1066aa8e0a05ac3721e"
	glabSHAAMD64 = "63df684842b0697e418649dd64374d37ac1da0fbcdcb422d6a5557222f8bb40f"
	glabSHAARM64 = "0716395a0efeb991fc844b33fe55310645afcbd9b9746e50c725e9d1cd1b8d16"

	defaultProfile = "devcontainer:\n  enabled: true\n"
	ghEntry        = "    gh:\n      enabled: true\n"
	glabEntry      = "    glab:\n      enabled: true\n"
)

// forgeSpec describes the Dockerfile block of one forge-CLI feature.
type forgeSpec struct {
	key, upper, version, shaAMD64, shaARM64, url, member string
}

func forgeSpecs() map[string]forgeSpec {
	return map[string]forgeSpec{
		"gh": {"gh", "GH", "2.103.0", ghSHAAMD64, ghSHAARM64,
			"https://github.com/cli/cli/releases/download/v${GH_VERSION}/gh_${GH_VERSION}_linux_${arch}.tar.gz",
			`"gh_${GH_VERSION}_linux_${arch}/bin/gh"`},
		"glab": {"glab", "GLAB", "1.124.0", glabSHAAMD64, glabSHAARM64,
			"https://gitlab.com/gitlab-org/cli/-/releases/v${GLAB_VERSION}/downloads/glab_${GLAB_VERSION}_linux_${arch}.tar.gz",
			"bin/glab"},
	}
}

// forgeProfiles are the three ways a Dockerfile is generated: the default
// profile, the sandbox profile and the `--sandbox` flag.
func forgeProfiles() map[string]struct {
	head    string
	sandbox bool
} {
	return map[string]struct {
		head    string
		sandbox bool
	}{
		"default profile": {defaultProfile, false},
		"sandbox profile": {sandboxProfile, false},
		"sandbox flag":    {defaultProfile, true},
	}
}

func featuresYAML(entries ...string) string {
	if len(entries) == 0 {
		return ""
	}
	return "  features:\n" + strings.Join(entries, "")
}

// assertForgeBlock pins the Dockerfile block of one forge-CLI feature (ADR-0035
// decisions 1 and 3): version and both checksums stand in one place, the check
// runs between the download and the extraction, only the binary is extracted
// without the owner of the archive, and the block stands behind `USER root`.
func assertForgeBlock(t *testing.T, df string, s forgeSpec) {
	t.Helper()
	for _, want := range []string{
		"ARG " + s.upper + "_VERSION=" + s.version + "\n",
		"ARG " + s.upper + "_SHA256_AMD64=" + s.shaAMD64 + "\n",
		"ARG " + s.upper + "_SHA256_ARM64=" + s.shaARM64 + "\n",
		s.url,
		`amd64) sha256="${` + s.upper + `_SHA256_AMD64}" ;; arm64) sha256="${` + s.upper + `_SHA256_ARM64}" ;; *) echo "` + s.key + ` has no archive for $arch" >&2; exit 1 ;; esac`,
	} {
		if !strings.Contains(df, want) {
			t.Errorf("Dockerfile lacks %q:\n%s", want, df)
		}
	}
	for _, once := range []string{s.shaAMD64, s.shaARM64, s.version} {
		if n := strings.Count(df, once); n != 1 {
			t.Errorf("Dockerfile names the pin %q %d times, want exactly once (one place)", once, n)
		}
	}
	start := strings.Index(df, "ARG "+s.upper+"_VERSION=")
	end := strings.Index(df[start+1:], `rm -rf "$work"`)
	if start < 0 || end < 0 {
		t.Fatalf("block of %s not found:\n%s", s.key, df)
	}
	block := df[start : start+1+end]
	last := -1
	for _, step := range []string{"curl -fsSLO", "sha256sum -c -", "tar -xzf"} {
		i := strings.Index(block, step)
		if i < 0 || i < last {
			t.Errorf("block step %q missing or out of order (index %d, previous %d):\n%s", step, i, last, block)
		}
		last = i
	}
	var tarLine string
	for _, line := range strings.Split(block, "\n") {
		if strings.Contains(line, "tar -xzf") {
			tarLine = strings.TrimSuffix(strings.TrimSpace(line), " \\")
		}
	}
	if !strings.Contains(tarLine, "--no-same-owner") || !strings.HasSuffix(tarLine, " "+s.member) {
		t.Errorf("tar line %q must keep --no-same-owner and extract only %s", tarLine, s.member)
	}
	if i := strings.Index(df, "USER root"); i < 0 || i > start {
		t.Errorf("the %s block must stand behind USER root:\n%s", s.key, df)
	}
	if start > strings.LastIndex(df, "USER vscode") {
		t.Errorf("Dockerfile must end on USER vscode:\n%s", df)
	}
}

// assertNoForgeSecrets pins ADR-0035 decisions 1 and 5 for a Dockerfile that
// carries only forge-CLI blocks: no signature or attestation step is claimed,
// no sign-in, no stored token, no start instruction. Comment lines are
// skipped.
func assertNoForgeSecrets(t *testing.T, df string) {
	t.Helper()
	for _, line := range strings.Split(df, "\n") {
		if strings.HasPrefix(line, "#") {
			continue
		}
		for _, bad := range []string{"auth login", "attestation", "cosign", "gpg", "--verify", "GH_TOKEN", "GITLAB_TOKEN", "ENV ", "EXPOSE", "ENTRYPOINT", "CMD "} {
			if strings.Contains(line, bad) {
				t.Errorf("Dockerfile line %q must not contain %q", line, bad)
			}
		}
	}
}

// forgeNotices returns the supply-chain notices of the forge-CLI feature key.
func forgeNotices(ws []driving.WarningEntry, key string) []driving.WarningEntry {
	var out []driving.WarningEntry
	for _, w := range ws {
		if strings.HasPrefix(w.Message, "devcontainer feature "+key+" installs") {
			out = append(out, w)
		}
	}
	return out
}

// B1, B3, B7 (fresh state): in the default profile, the sandbox profile and
// with `--sandbox` the enabled feature gets its pinned block behind a single
// `USER root`; the other feature leaves no trace; the second run is a NoOp.
func TestGenerateDevcontainer_ForgeCLI_BlocksInBothProfiles(t *testing.T) {
	t.Parallel()
	for pname, p := range forgeProfiles() {
		for fname, keys := range map[string][]string{"gh": {"gh"}, "glab": {"glab"}, "both": {"gh", "glab"}} {
			t.Run(pname+"/"+fname, func(t *testing.T) {
				t.Parallel()
				var entries []string
				for _, k := range keys {
					entries = append(entries, "    "+k+":\n      enabled: true\n")
				}
				svc, fs := newGenerateService(t)
				seedUBootYAMLWithFeatures(t, fs, p.head+featuresYAML(entries...))
				if _, err := generateSandbox(svc, p.sandbox); err != nil {
					t.Fatalf("generate: %v", err)
				}
				raw, _ := fs.ReadFile(dockerfilePath())
				df := string(raw)
				specs := forgeSpecs()
				for _, k := range keys {
					assertForgeBlock(t, df, specs[k])
				}
				for k, s := range specs {
					if !contains(keys, k) && strings.Contains(df, s.upper+"_VERSION") {
						t.Errorf("Dockerfile carries the block of %s although it is not enabled:\n%s", k, df)
					}
				}
				if n := strings.Count(df, "USER root"); n != 1 {
					t.Errorf("USER root appears %d times, want 1:\n%s", n, df)
				}
				assertNoForgeSecrets(t, df)
				if r2, err := generateSandbox(svc, p.sandbox); err != nil || r2.Action != driving.GenerateActionNoOp {
					t.Errorf("second run: action=%v err=%v, want NoOp", r2.Action, err)
				}
			})
		}
	}
}

func contains(list []string, s string) bool {
	for _, e := range list {
		if e == s {
			return true
		}
	}
	return false
}

// B1: in the default profile `USER root` stands only with a reason: the
// feature or a UID override; with both it stands once, ahead of the UID step
// and the block.
func TestGenerateDevcontainer_ForgeCLI_UserRootOnlyWithReason(t *testing.T) {
	t.Parallel()
	const uid = "  user:\n    uid: 501\n"
	for name, tc := range map[string]struct {
		yaml     string
		wantRoot int
		wantGh   bool
		wantUID  bool
	}{
		"no feature, no UID":      {defaultProfile, 0, false, false},
		"gh":                      {defaultProfile + featuresYAML(ghEntry), 1, true, false},
		"UID only":                {defaultProfile + uid, 1, false, true},
		"gh with UID override":    {defaultProfile + uid + featuresYAML(ghEntry), 1, true, true},
		"glab with UID override":  {defaultProfile + uid + featuresYAML(glabEntry), 1, false, true},
		"gh switched off, no UID": {defaultProfile + featuresYAML("    gh:\n      enabled: false\n"), 0, false, false},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, tc.yaml)
			if _, err := generateSandbox(svc, false); err != nil {
				t.Fatalf("generate: %v", err)
			}
			raw, _ := fs.ReadFile(dockerfilePath())
			df := string(raw)
			if n := strings.Count(df, "USER root"); n != tc.wantRoot {
				t.Errorf("USER root appears %d times, want %d:\n%s", n, tc.wantRoot, df)
			}
			if got := strings.Contains(df, "ARG GH_VERSION="); got != tc.wantGh {
				t.Errorf("gh block present = %v, want %v:\n%s", got, tc.wantGh, df)
			}
			if got := strings.Contains(df, "ARG USER_UID=501"); got != tc.wantUID {
				t.Errorf("UID step present = %v, want %v:\n%s", got, tc.wantUID, df)
			}
			if tc.wantGh {
				assertForgeBlock(t, df, forgeSpecs()["gh"])
			}
		})
	}
}

// B5 (ADR-0035 decision 6): exactly one notice per active feature, in every
// profile, with the code of LH-NFA-SEC-004. The text of `gh` names the
// attestation that `u-boot` does not verify and the SHA-256 check; the text of
// `glab` names no attestation. Neither carries a version, and neither carries
// SecurityRelaxation or DegradationFallback, so `--quiet` drops them (the
// filter is the CLI's, see cmd/uboot TestRun_Generate_ForgeCLINotices).
func TestGenerateDevcontainer_ForgeCLI_Notices(t *testing.T) {
	t.Parallel()
	for pname, p := range forgeProfiles() {
		t.Run(pname, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, p.head+featuresYAML(ghEntry, glabEntry))
			resp, err := generateSandbox(svc, p.sandbox)
			if err != nil {
				t.Fatalf("generate: %v", err)
			}
			for key, tc := range map[string]struct {
				must    []string
				mustNot []string
			}{
				"gh":   {[]string{"cli/cli", "GH_VERSION", "SHA-256", "no signature was found", "not exhaustive", "attestation", "not verified by u-boot"}, []string{"2.103.0", "GLAB"}},
				"glab": {[]string{"gitlab-org/cli", "GLAB_VERSION", "SHA-256", "no signature was found", "not exhaustive"}, []string{"attestation", "1.124.0", "GH_"}},
			} {
				notices := forgeNotices(resp.Warnings, key)
				if len(notices) != 1 {
					t.Fatalf("%s notices = %d, want 1: %+v", key, len(notices), resp.Warnings)
				}
				w := notices[0]
				if w.Code != "LH-NFA-SEC-004" || w.Level != "warn" || w.Subject != ".devcontainer/devcontainer.json" {
					t.Errorf("%s notice = %+v, want a warn entry with code LH-NFA-SEC-004 on devcontainer.json", key, w)
				}
				if w.SecurityRelaxation || w.DegradationFallback {
					t.Errorf("%s notice must be unflagged (dropped under --quiet): %+v", key, w)
				}
				for _, part := range tc.must {
					if !strings.Contains(w.Message, part) {
						t.Errorf("%s notice lacks %q: %q", key, part, w.Message)
					}
				}
				for _, part := range tc.mustNot {
					if strings.Contains(w.Message, part) {
						t.Errorf("%s notice must not contain %q: %q", key, part, w.Message)
					}
				}
			}
			// One feature, one notice: the other feature adds none.
			svc2, fs2 := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs2, p.head+featuresYAML(glabEntry))
			resp2, err := generateSandbox(svc2, p.sandbox)
			if err != nil {
				t.Fatalf("generate: %v", err)
			}
			if n := len(forgeNotices(resp2.Warnings, "gh")); n != 0 {
				t.Errorf("gh notices without gh = %d, want 0", n)
			}
		})
	}
}

// B6 (ADR-0035 decision 4): the hosts of the features in the effective
// allowlist. The two GitHub hosts of `gh` are part of the base list already;
// `glab` adds `gitlab.com`; nothing else is added for either.
func TestGenerateDevcontainer_ForgeCLI_EgressHosts(t *testing.T) {
	t.Parallel()
	base := []string{"api.github.com", "codeload.github.com", "deb.debian.org", "github.com",
		"objects.githubusercontent.com", "raw.githubusercontent.com", "security.debian.org"}
	withHosts := func(extra ...string) []string {
		out := append(append([]string{}, base...), extra...)
		sort.Strings(out)
		return out
	}
	cases := map[string]struct {
		features  string
		wantHosts []string
	}{
		"no feature":    {"", withHosts()},
		"gh only":       {featuresYAML(ghEntry), withHosts()},
		"glab only":     {featuresYAML(glabEntry), withHosts("gitlab.com")},
		"both":          {featuresYAML(ghEntry, glabEntry), withHosts("gitlab.com")},
		"glab disabled": {featuresYAML("    glab:\n      enabled: false\n"), withHosts()},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, sandboxProfile+tc.features+egressOn)
			if _, err := generateSandbox(svc, false); err != nil {
				t.Fatalf("generate: %v", err)
			}
			hosts := egressHostList(t, fs)
			if !sort.StringsAreSorted(hosts) || !reflect.DeepEqual(hosts, tc.wantHosts) {
				t.Errorf("egress hosts =\n%v\nwant (sorted, no duplicates)\n%v", hosts, tc.wantHosts)
			}
			if r2, err := generateSandbox(svc, false); err != nil || r2.Action != driving.GenerateActionNoOp {
				t.Errorf("second run: action=%v err=%v, want NoOp", r2.Action, err)
			}
		})
	}
}

// B9 (ADR-0035 decision 5): in the sandbox profile `remoteEnv` carries the
// reference of the token variable of each active feature next to `GIT_TOKEN`,
// in a fixed order, and only a reference. The Dockerfile names no token
// variable. A switched-off feature, an absent feature and a `source:` override
// add nothing.
func TestGenerateDevcontainer_ForgeCLI_RemoteEnvReferences(t *testing.T) {
	t.Parallel()
	const external = "ghcr.io/example/features/gh"
	ref := func(name string) string { return `    "` + name + `": "${localEnv:` + name + `}"` }
	env := func(names ...string) string {
		lines := []string{ref("GIT_TOKEN")}
		for _, n := range names {
			lines = append(lines, ref(n))
		}
		return "  \"remoteEnv\": {\n" + strings.Join(lines, ",\n") + "\n  },\n"
	}
	const claude = "    claude-code:\n      enabled: true\n    node:\n      enabled: true\n"
	cases := map[string]struct {
		yaml string
		want string
	}{
		"gh":             {sandboxProfile + featuresYAML(ghEntry), env("GH_TOKEN")},
		"glab":           {sandboxProfile + featuresYAML(glabEntry), env("GITLAB_TOKEN")},
		"both":           {sandboxProfile + featuresYAML(glabEntry, ghEntry), env("GH_TOKEN", "GITLAB_TOKEN")},
		"with claude":    {sandboxProfile + featuresYAML(glabEntry, ghEntry, claude), env("ANTHROPIC_API_KEY", "GH_TOKEN", "GITLAB_TOKEN")},
		"disabled":       {sandboxProfile + featuresYAML("    gh:\n      enabled: false\n"), env()},
		"no feature":     {sandboxProfile, env()},
		"source for gh":  {sandboxProfile + "  featureSources:\n    allow:\n      - https://" + external + "\n" + featuresYAML("    gh:\n      enabled: true\n      source: https://"+external+"\n"), env()},
		"glab beside gh": {sandboxProfile + featuresYAML(ghEntry, "    glab:\n      enabled: false\n"), env("GH_TOKEN")},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, tc.yaml)
			if _, err := generateSandbox(svc, false); err != nil {
				t.Fatalf("generate: %v", err)
			}
			raw, _ := fs.ReadFile(devcontainerJSONPath())
			if !strings.Contains(string(raw), tc.want) {
				t.Errorf("devcontainer.json lacks\n%s\nin:\n%s", tc.want, raw)
			}
			m := devcontainerJSONMap(t, fs)
			remoteEnv, _ := m["remoteEnv"].(map[string]any)
			for k, v := range remoteEnv {
				s, _ := v.(string)
				if want := "${localEnv:" + k + "}"; s != want {
					t.Errorf("remoteEnv[%s] = %q, want the reference %q", k, s, want)
				}
			}
			df, _ := fs.ReadFile(dockerfilePath())
			for _, bad := range []string{"GH_TOKEN", "GITLAB_TOKEN"} {
				if strings.Contains(string(df), bad) {
					t.Errorf("Dockerfile must not name %s:\n%s", bad, df)
				}
			}
		})
	}
}

// X1 (ADR-0035 decisions 3 and 5): in the default profile the features change
// the Dockerfile and nothing else; devcontainer.json is byte-identical to the
// run without them, has no `remoteEnv`, no `mounts` and no `features` key, and
// the second run is a NoOp.
func TestGenerateDevcontainer_ForgeCLI_DefaultProfileJSONUnchanged(t *testing.T) {
	t.Parallel()
	with, fsWith := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fsWith, defaultProfile+featuresYAML(ghEntry, glabEntry))
	if _, err := generateSandbox(with, false); err != nil {
		t.Fatalf("generate: %v", err)
	}
	without, fsWithout := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fsWithout, defaultProfile)
	if _, err := generateSandbox(without, false); err != nil {
		t.Fatalf("generate: %v", err)
	}
	jsonWith, _ := fsWith.ReadFile(devcontainerJSONPath())
	jsonWithout, _ := fsWithout.ReadFile(devcontainerJSONPath())
	if string(jsonWith) != string(jsonWithout) {
		t.Errorf("devcontainer.json differs by the forge-CLI features in the default profile:\n%s\n---\n%s", jsonWith, jsonWithout)
	}
	for _, bad := range []string{"remoteEnv", "mounts", "GH_TOKEN", "GITLAB_TOKEN"} {
		if strings.Contains(string(jsonWith), bad) {
			t.Errorf("devcontainer.json must not contain %q in the default profile:\n%s", bad, jsonWith)
		}
	}
	if feats := parseDevcontainerFeatures(t, jsonWith); len(feats) != 0 {
		t.Errorf("gh and glab render no key under features, got %v", feats)
	}
	if r2, err := generateSandbox(with, false); err != nil || r2.Action != driving.GenerateActionNoOp {
		t.Errorf("second run: action=%v err=%v, want NoOp", r2.Action, err)
	}
}

// X2 (ADR-0035 decision 2): a `version:` override without `source:` ends with
// exit code 10 (ErrGenerateManualConflict) in both profiles and writes
// nothing; the message names the key.
func TestGenerateDevcontainer_ForgeCLI_VersionOverrideWithoutSourceRejected(t *testing.T) {
	t.Parallel()
	for _, key := range []string{"gh", "glab"} {
		for pname, head := range map[string]string{"default profile": defaultProfile, "sandbox profile": sandboxProfile} {
			t.Run(key+"/"+pname, func(t *testing.T) {
				t.Parallel()
				svc, fs := newGenerateService(t)
				seedUBootYAMLWithFeatures(t, fs, head+"  features:\n    "+key+":\n      enabled: true\n      version: \"9.9.9\"\n")
				_, err := generateSandbox(svc, false)
				if !errors.Is(err, driving.ErrGenerateManualConflict) || err == nil {
					t.Fatalf("generate error = %v, want a wrap of ErrGenerateManualConflict (exit code 10)", err)
				}
				if !strings.Contains(err.Error(), "features."+key+".version") {
					t.Errorf("error does not name the key: %v", err)
				}
				if exists, _ := fs.Exists(dockerfilePath()); exists {
					t.Errorf("nothing may be written on a validation error")
				}
			})
		}
	}
}

// X3 (ADR-0035, reading of the planner): a `source:` override makes `gh` an
// external OCI feature: it renders its key, the Dockerfile carries no block
// and no `USER root` in the default profile, `remoteEnv` no reference; the
// sandbox profile gives the one warning about an external source and none of
// the supply-chain notice of the block.
func TestGenerateDevcontainer_ForgeCLI_SourceOverrideIsExternalFeature(t *testing.T) {
	t.Parallel()
	const external = "ghcr.io/example/features/gh"
	features := "  featureSources:\n    allow:\n      - https://" + external + "\n" +
		featuresYAML("    gh:\n      enabled: true\n      source: https://"+external+"\n")
	for pname, head := range map[string]string{"default profile": defaultProfile, "sandbox profile": sandboxProfile} {
		t.Run(pname, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, head+features)
			resp, err := generateSandbox(svc, false)
			if err != nil {
				t.Fatalf("generate: %v", err)
			}
			raw, _ := fs.ReadFile(dockerfilePath())
			df := string(raw)
			if strings.Contains(df, "GH_VERSION") || strings.Contains(df, "USER root") && head == defaultProfile {
				t.Errorf("no block and, in the default profile, no USER root expected:\n%s", df)
			}
			js, _ := fs.ReadFile(devcontainerJSONPath())
			if strings.Contains(string(js), "GH_TOKEN") {
				t.Errorf("no token reference expected:\n%s", js)
			}
			if keys := parseDevcontainerFeatures(t, js); len(keys) != 1 || keys["https://"+external+":1"] == nil {
				t.Errorf("features keys = %v, want the external feature", keys)
			}
			if n := len(forgeNotices(resp.Warnings, "gh")); n != 0 {
				t.Errorf("supply-chain notice of the block = %d, want 0: %+v", n, resp.Warnings)
			}
			nExternal := 0
			for _, w := range resp.Warnings {
				if strings.Contains(w.Message, "feature gh comes from the external source") {
					nExternal++
				}
			}
			if want := map[string]int{"default profile": 0, "sandbox profile": 1}[pname]; nExternal != want {
				t.Errorf("warnings about the external source = %d, want %d: %+v", nExternal, want, resp.Warnings)
			}
		})
	}
}

// X4 (kombination-ollama-nested-runtime-ohne-lauf): with a nested runtime,
// `ollama`, `gh`, `glab` and the egress restriction the blocks stand in a
// fixed order (runtime, ollama, gh, glab, egress), `gh` and `glab` add no
// volume, the hosts are the sorted union, the second run is a NoOp.
func TestGenerateDevcontainer_ForgeCLI_TogetherWithOllamaAndRuntime(t *testing.T) {
	t.Parallel()
	const features = "  features:\n    ollama:\n      enabled: true\n    gh:\n      enabled: true\n    glab:\n      enabled: true\n"
	for rt, runtimeMarker := range map[string]string{"podman": "ARG PODMAN_VERSION=", "docker": "# Nested Docker daemon"} {
		t.Run(rt, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, sandboxProfile+features+"  sandbox:\n    nestedRuntime: "+rt+"\n    egress:\n      enabled: true\n")
			if _, err := generateSandbox(svc, false); err != nil {
				t.Fatalf("generate: %v", err)
			}
			raw, _ := fs.ReadFile(dockerfilePath())
			df := string(raw)
			assertForgeBlock(t, df, forgeSpecs()["gh"])
			assertForgeBlock(t, df, forgeSpecs()["glab"])
			last := -1
			for _, marker := range []string{"USER root", runtimeMarker, "ARG OLLAMA_VERSION=", "ARG GH_VERSION=", "ARG GLAB_VERSION=", "# Egress restriction", "USER vscode"} {
				i := strings.Index(df, marker)
				if i < 0 || i < last {
					t.Errorf("Dockerfile marker %q missing or out of order (index %d, previous %d):\n%s", marker, i, last, df)
				}
				last = i
			}
			js, _ := fs.ReadFile(devcontainerJSONPath())
			for _, bad := range []string{"-gh-", "-glab-", "/home/vscode/.config"} {
				if strings.Contains(string(js), bad) {
					t.Errorf("gh and glab add no volume, devcontainer.json contains %q:\n%s", bad, js)
				}
			}
			if n := strings.Count(string(js), `type=volume`); n != 3 {
				t.Errorf("volumes = %d, want 3 (workspace, runtime, ollama):\n%s", n, js)
			}
			hosts := egressHostList(t, fs)
			for _, h := range []string{"gitlab.com", "api.github.com", "github.com", "registry.ollama.ai", "ghcr.io"} {
				if !contains(hosts, h) {
					t.Errorf("egress hosts lack %s: %v", h, hosts)
				}
			}
			if !sort.StringsAreSorted(hosts) {
				t.Errorf("egress hosts are not sorted: %v", hosts)
			}
			seen := map[string]bool{}
			for _, h := range hosts {
				if seen[h] {
					t.Errorf("egress hosts contain %s twice: %v", h, hosts)
				}
				seen[h] = true
			}
			if r2, err := generateSandbox(svc, false); err != nil || r2.Action != driving.GenerateActionNoOp {
				t.Errorf("second run: action=%v err=%v, want NoOp", r2.Action, err)
			}
		})
	}
}

// X5 (replay): enabling the feature later replaces the stale managed block
// (UpdatedBlock) and keeps the line the user added behind it; the next run is a
// NoOp; disabling it removes the block again and keeps the user line.
func TestGenerateDevcontainer_ForgeCLI_ActivateAndDeactivateLater(t *testing.T) {
	t.Parallel()
	for pname, head := range map[string]string{"default profile": defaultProfile, "sandbox profile": sandboxProfile} {
		t.Run(pname, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, head)
			if _, err := generateSandbox(svc, false); err != nil {
				t.Fatalf("first generate: %v", err)
			}
			body, _ := fs.ReadFile(dockerfilePath())
			if err := fs.WriteFile(dockerfilePath(), append(body, []byte("# user addition\nRUN echo user\n")...), 0o644); err != nil {
				t.Fatal(err)
			}

			seedUBootYAMLWithFeatures(t, fs, head+featuresYAML(ghEntry))
			resp, err := generateSandbox(svc, false)
			if err != nil {
				t.Fatalf("activate: %v", err)
			}
			if resp.Action != driving.GenerateActionUpdatedBlock {
				t.Errorf("activate: action = %v, want UpdatedBlock", resp.Action)
			}
			after, _ := fs.ReadFile(dockerfilePath())
			assertForgeBlock(t, string(after), forgeSpecs()["gh"])
			if !strings.Contains(string(after), "# user addition\nRUN echo user\n") {
				t.Errorf("the user line behind the block was lost:\n%s", after)
			}
			if r2, err := generateSandbox(svc, false); err != nil || r2.Action != driving.GenerateActionNoOp {
				t.Errorf("second run: action=%v err=%v, want NoOp", r2.Action, err)
			}

			seedUBootYAMLWithFeatures(t, fs, head+featuresYAML("    gh:\n      enabled: false\n"))
			if _, err := generateSandbox(svc, false); err != nil {
				t.Fatalf("deactivate: %v", err)
			}
			gone, _ := fs.ReadFile(dockerfilePath())
			if strings.Contains(string(gone), "GH_VERSION") {
				t.Errorf("the block must be gone after disabling the feature:\n%s", gone)
			}
			if !strings.Contains(string(gone), "# user addition\nRUN echo user\n") {
				t.Errorf("the user line behind the block was lost:\n%s", gone)
			}
		})
	}
}

// The way out of the rejection of X2: the same configuration without the
// `version:` entry generates without error and writes the Dockerfile block.
func TestGenerateDevcontainer_ForgeCLI_RemovingVersionEntryLetsGenerateRun(t *testing.T) {
	t.Parallel()
	for _, key := range []string{"gh", "glab"} {
		for pname, head := range map[string]string{"default profile": defaultProfile, "sandbox profile": sandboxProfile} {
			t.Run(key+"/"+pname, func(t *testing.T) {
				t.Parallel()
				svc, fs := newGenerateService(t)
				seedUBootYAMLWithFeatures(t, fs, head+"  features:\n    "+key+":\n      enabled: true\n")
				if _, err := generateSandbox(svc, false); err != nil {
					t.Fatalf("generate without the version entry: %v", err)
				}
				raw, _ := fs.ReadFile(dockerfilePath())
				assertForgeBlock(t, string(raw), forgeSpecs()[key])
			})
		}
	}
}
