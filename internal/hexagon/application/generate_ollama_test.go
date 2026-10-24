package application_test

import (
	"context"
	"errors"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// Golden cases of the catalogue key `ollama` (LH-FA-DEV-003): a Dockerfile
// block without an OCI render key, a model volume, two egress hosts, one
// supply-chain notice, nothing that starts or exposes the server.

const (
	ollamaSHAAMD64 = "c164ab8b1b28e2a247c29719b8e404c423c5b83ae64c3ab4e59d50c7ff7fc5e0"
	ollamaSHAARM64 = "0a279d9428ad17b06fdb7dddf3b059b2c148433585fa3b30676a102fc459bb9f"

	ollamaFeature  = "  features:\n    ollama:\n      enabled: true\n"
	sandboxProfile = "devcontainer:\n  enabled: true\n  profile: sandbox\n"
	egressOn       = "  sandbox:\n    egress:\n      enabled: true\n"

	workspaceMountLine = `"source=t-uboot-gen-workspace-${devcontainerId},target=/workspaces/t-uboot-gen,type=volume"`
	containersMount    = `"source=t-uboot-gen-containers-${devcontainerId},target=/home/vscode/.local/share/containers,type=volume"`
	dockerMount        = `"source=t-uboot-gen-docker-${devcontainerId},target=/var/lib/docker,type=volume"`
	ollamaMount        = `"source=t-uboot-gen-ollama-${devcontainerId},target=/home/vscode/.ollama,type=volume"`
)

// ollamaNotices returns the warnings that name the ollama feature.
func ollamaNotices(ws []driving.WarningEntry) []driving.WarningEntry {
	var out []driving.WarningEntry
	for _, w := range ws {
		if strings.Contains(w.Message, "feature ollama ") {
			out = append(out, w)
		}
	}
	return out
}

// assertOllamaDockerfile pins the Ollama block of the generated Dockerfile
// (ADR-0024 decisions 1, 3, 6, 8): version and both checksums stand in one
// place, the check runs before the extraction, the GPU directories are left
// out, `zstd` is removed again, the model store belongs to `vscode`, and no
// instruction starts the server or sets a listen address (decision 5).
func assertOllamaDockerfile(t *testing.T, df string) {
	t.Helper()
	for _, want := range []string{
		"ARG OLLAMA_VERSION=0.40.6\n",
		"ARG OLLAMA_SHA256_AMD64=" + ollamaSHAAMD64 + "\n",
		"ARG OLLAMA_SHA256_ARM64=" + ollamaSHAARM64 + "\n",
		"https://github.com/ollama/ollama/releases/download/v${OLLAMA_VERSION}/ollama-linux-${arch}.tar.zst",
		`amd64) sha256="${OLLAMA_SHA256_AMD64}" ;; arm64) sha256="${OLLAMA_SHA256_ARM64}" ;; *) echo "ollama has no archive for $arch" >&2; exit 1 ;; esac`,
		"--no-same-owner",
		"--exclude='lib/ollama/cuda_*'",
		"--exclude='lib/ollama/vulkan'",
	} {
		if !strings.Contains(df, want) {
			t.Errorf("Dockerfile lacks %q:\n%s", want, df)
		}
	}
	for _, once := range []string{ollamaSHAAMD64, ollamaSHAARM64, "0.40.6"} {
		if n := strings.Count(df, once); n != 1 {
			t.Errorf("Dockerfile names the pin %q %d times, want exactly once (one place)", once, n)
		}
	}
	// Order: install zstd, check the checksum, extract, remove zstd, prepare the store.
	last := -1
	for _, step := range []string{
		"install -y --no-install-recommends zstd",
		"sha256sum -c -",
		"tar --zstd -xf",
		"apt-get purge -y --auto-remove zstd",
		"mkdir -p /home/vscode/.ollama",
		"chown vscode:vscode /home/vscode/.ollama",
	} {
		i := strings.Index(df, step)
		if i < 0 || i < last {
			t.Errorf("Dockerfile step %q missing or out of order (index %d, previous %d):\n%s", step, i, last, df)
		}
		last = i
	}
	for _, line := range strings.Split(df, "\n") {
		if strings.HasPrefix(line, "#") {
			continue
		}
		for _, bad := range []string{"OLLAMA_HOST", "ollama serve", "ENV ", "EXPOSE", "ENTRYPOINT", "CMD "} {
			if strings.Contains(line, bad) {
				t.Errorf("Dockerfile line %q must not contain %q (the image starts no server and sets no listen address)", line, bad)
			}
		}
	}
	if strings.Index(df, "USER root") > strings.Index(df, "ARG OLLAMA_VERSION") {
		t.Errorf("the Ollama block must stand behind USER root:\n%s", df)
	}
	if strings.Index(df, "ARG OLLAMA_VERSION") > strings.LastIndex(df, "USER vscode") {
		t.Errorf("Dockerfile must end on USER vscode:\n%s", df)
	}
}

// blockSandboxDockerfile renders the Dockerfile of a sandbox project with the
// Dockerfile-block feature `key` enabled; the catalogue test looks for the
// pin there.
func blockSandboxDockerfile(t *testing.T, key string) string {
	t.Helper()
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, sandboxProfile+"  features:\n    "+key+":\n      enabled: true\n")
	if _, err := generateSandbox(svc, false); err != nil {
		t.Fatalf("generate: %v", err)
	}
	df, _ := fs.ReadFile(dockerfilePath())
	return string(df)
}

// ADR-0024 decisions 1, 2, 3, 5, 8 (fresh state): the sandbox profile
// installs Ollama from the pinned archive, adds the model volume next to
// the workspace volume, renders no key under `features`, reports the
// supply-chain limit once and adds nothing that starts the server or sets a
// listen address; the second run is a NoOp.
func TestGenerateDevcontainer_Ollama_Sandbox(t *testing.T) {
	t.Parallel()
	for name, tc := range map[string]struct {
		yaml    string
		sandbox bool
	}{
		"profile":      {sandboxProfile + ollamaFeature, false},
		"sandbox flag": {"devcontainer:\n  enabled: true\n" + ollamaFeature, true},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, tc.yaml)
			resp, err := generateSandbox(svc, tc.sandbox)
			if err != nil {
				t.Fatalf("generate: %v", err)
			}
			df, _ := fs.ReadFile(dockerfilePath())
			assertOllamaDockerfile(t, string(df))

			raw, _ := fs.ReadFile(devcontainerJSONPath())
			body := string(raw)
			if feats := parseDevcontainerFeatures(t, raw); len(feats) != 0 {
				t.Errorf("ollama renders no key under features, got %v", feats)
			}
			want := "  \"mounts\": [\n    " + workspaceMountLine + ",\n    " + ollamaMount + "\n  ],\n"
			if !strings.Contains(body, want) {
				t.Errorf("devcontainer.json lacks the mounts block\n%s\nin:\n%s", want, body)
			}
			// ADR-0024 decisions 4 and 5: no start command, no listen address, no port.
			for _, bad := range []string{"OLLAMA_HOST", "forwardPorts", "postStartCommand", "postCreateCommand", "ollama serve", "11434"} {
				if strings.Contains(body, bad) {
					t.Errorf("devcontainer.json must not contain %q:\n%s", bad, body)
				}
			}
			m := devcontainerJSONMap(t, fs)
			if _, ok := m["runArgs"]; ok {
				t.Errorf("ollama adds no runArgs: %v", m["runArgs"])
			}

			notices := ollamaNotices(resp.Warnings)
			if len(notices) != 1 {
				t.Fatalf("ollama notices = %d, want 1: %+v", len(notices), resp.Warnings)
			}
			w := notices[0]
			if w.Code != "LH-NFA-SEC-004" || w.Level != "warn" || w.Subject != ".devcontainer/devcontainer.json" {
				t.Errorf("notice = %+v, want a warn entry with code LH-NFA-SEC-004 on devcontainer.json", w)
			}
			for _, part := range []string{"github.com/ollama/ollama", "OLLAMA_VERSION", "SHA-256", "publishes no signature"} {
				if !strings.Contains(w.Message, part) {
					t.Errorf("notice lacks %q: %q", part, w.Message)
				}
			}
			if strings.Contains(w.Message, "third-party publisher") || strings.Contains(w.Message, "install script") {
				t.Errorf("notice is no OCI third-party notice: %q", w.Message)
			}

			if r2, err := generateSandbox(svc, tc.sandbox); err != nil || r2.Action != driving.GenerateActionNoOp {
				t.Errorf("second run: action=%v err=%v, want NoOp", r2.Action, err)
			}
		})
	}
}

// ADR-0024 decision 3 and 5, Verification-Report pins V-3 and V-5: the
// supply-chain notice and the "no effect" notice carry neither
// SecurityRelaxation nor DegradationFallback, so `--quiet` drops both
// (the filter is the CLI's, see cmd/uboot TestRun_Generate_OllamaNotices).
// Without the sandbox profile the Dockerfile has no block, the JSON no
// volume, and the notice says that nothing was generated.
func TestGenerateDevcontainer_Ollama_NoticesAreNotRelaxations(t *testing.T) {
	t.Parallel()
	for name, tc := range map[string]struct {
		yaml     string
		wantCode string
		wantText string
	}{
		"sandbox":         {sandboxProfile + ollamaFeature, "LH-NFA-SEC-004", "publishes no signature"},
		"default profile": {"devcontainer:\n  enabled: true\n" + ollamaFeature, "LH-FA-DEV-003", "has no effect without devcontainer.profile=sandbox"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, tc.yaml)
			resp, err := generateSandbox(svc, false)
			if err != nil {
				t.Fatalf("generate: %v", err)
			}
			notices := ollamaNotices(resp.Warnings)
			if len(notices) != 1 {
				t.Fatalf("ollama notices = %d, want 1: %+v", len(notices), resp.Warnings)
			}
			w := notices[0]
			if w.Code != tc.wantCode || !strings.Contains(w.Message, tc.wantText) {
				t.Errorf("notice = %+v, want code %s and %q", w, tc.wantCode, tc.wantText)
			}
			if w.SecurityRelaxation || w.DegradationFallback {
				t.Errorf("notice must be unflagged (dropped under --quiet): %+v", w)
			}
		})
	}
}

// Verification-Report pin V-5 (condition "effective in the sandbox profile"):
// without the sandbox profile `ollama` generates no installation and no
// volume; the Dockerfile and devcontainer.json stay byte-identical to the
// run without the feature, and the second run is a NoOp. The variant with a
// UID override covers the Dockerfile branch that carries `USER root` without
// the sandbox profile.
func TestGenerateDevcontainer_Ollama_DefaultProfileGeneratesNothing(t *testing.T) {
	t.Parallel()
	for name, head := range map[string]string{
		"default profile":          "devcontainer:\n  enabled: true\n",
		"default profile, UID 501": "devcontainer:\n  enabled: true\n  user:\n    uid: 501\n",
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			with, fsWith := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fsWith, head+ollamaFeature)
			if _, err := generateSandbox(with, false); err != nil {
				t.Fatalf("generate: %v", err)
			}
			without, fsWithout := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fsWithout, head)
			if _, err := generateSandbox(without, false); err != nil {
				t.Fatalf("generate: %v", err)
			}
			dfWith, _ := fsWith.ReadFile(dockerfilePath())
			dfWithout, _ := fsWithout.ReadFile(dockerfilePath())
			if string(dfWith) != string(dfWithout) {
				t.Errorf("Dockerfile differs by the ollama feature without the sandbox profile:\n%s\n---\n%s", dfWith, dfWithout)
			}
			if strings.Contains(strings.ToLower(string(dfWith)), "ollama") {
				t.Errorf("Dockerfile mentions ollama without the sandbox profile:\n%s", dfWith)
			}
			jsonWith, _ := fsWith.ReadFile(devcontainerJSONPath())
			jsonWithout, _ := fsWithout.ReadFile(devcontainerJSONPath())
			if string(jsonWith) != string(jsonWithout) {
				t.Errorf("devcontainer.json differs by the ollama feature without the sandbox profile:\n%s\n---\n%s", jsonWith, jsonWithout)
			}
			if r2, err := generateSandbox(with, false); err != nil || r2.Action != driving.GenerateActionNoOp {
				t.Errorf("second run: action=%v err=%v, want NoOp", r2.Action, err)
			}
		})
	}
}

// `enabled: false`, an absent key and a `source:` override produce no
// block, no volume and no notice of the Dockerfile block; the override
// renders the external feature like any catalogue key (ADR-0024 decision 2)
// and, in the sandbox profile, gets the one warning about a feature from an
// external source (ADR-0032 decision 3).
func TestGenerateDevcontainer_Ollama_NotActive(t *testing.T) {
	t.Parallel()
	const external = "ghcr.io/example/features/ollama"
	for name, tc := range map[string]struct {
		features string
		wantKeys []string
	}{
		"disabled":        {"  features:\n    ollama:\n      enabled: false\n", nil},
		"absent":          {"", nil},
		"source override": {"  features:\n    ollama:\n      enabled: true\n      source: https://" + external + "\n", []string{"https://" + external + ":1"}},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			yaml := sandboxProfile + tc.features
			if tc.wantKeys != nil {
				yaml = sandboxProfile + "  featureSources:\n    allow:\n      - https://" + external + "\n" + tc.features
			}
			seedUBootYAMLWithFeatures(t, fs, yaml)
			resp, err := generateSandbox(svc, false)
			if err != nil {
				t.Fatalf("generate: %v", err)
			}
			df, _ := fs.ReadFile(dockerfilePath())
			raw, _ := fs.ReadFile(devcontainerJSONPath())
			if strings.Contains(strings.ToLower(string(df)), "ollama") || strings.Contains(string(raw), "ollama-") {
				t.Errorf("no block and no volume expected:\n%s\n%s", df, raw)
			}
			wantNotices := 0
			if tc.wantKeys != nil {
				wantNotices = 1
			}
			notices := ollamaNotices(resp.Warnings)
			if len(notices) != wantNotices {
				t.Errorf("ollama notices = %d, want %d: %+v", len(notices), wantNotices, notices)
			}
			for _, n := range notices {
				if !strings.Contains(n.Message, "comes from the external source") {
					t.Errorf("only the warning about an external feature expected, no notice of the Dockerfile block: %+v", n)
				}
			}
			var got []string
			for k := range parseDevcontainerFeatures(t, raw) {
				got = append(got, k)
			}
			sort.Strings(got)
			if !reflect.DeepEqual(got, tc.wantKeys) {
				t.Errorf("features keys = %v, want %v", got, tc.wantKeys)
			}
		})
	}
}

// ADR-0024 decision 2: a `version:` override without `source:` is a
// validation error in `generate devcontainer` (exit code 10 through
// ErrGenerateManualConflict, nothing written); the pin hangs on the checksums
// in the Dockerfile.
func TestGenerateDevcontainer_Ollama_VersionOverrideWithoutSourceRejected(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, sandboxProfile+"  features:\n    ollama:\n      enabled: true\n      version: \"0.40.7\"\n")
	_, err := generateSandbox(svc, false)
	if !errors.Is(err, driving.ErrGenerateManualConflict) || err == nil {
		t.Fatalf("generate error = %v, want a wrap of ErrGenerateManualConflict (exit code 10)", err)
	}
	if !strings.Contains(err.Error(), "features.ollama.version") {
		t.Errorf("error does not name the key: %v", err)
	}
	if exists, _ := fs.Exists(dockerfilePath()); exists {
		t.Errorf("nothing may be written on a validation error")
	}
}

// ADR-0024 decision 8, byte-identity: the `mounts` text of the existing
// configurations is unchanged, and the model volume joins them in the order
// workspace, runtime volume, Ollama.
func TestGenerateDevcontainer_Ollama_MountsComposition(t *testing.T) {
	t.Parallel()
	nested := func(rt string) string { return "  sandbox:\n    nestedRuntime: " + rt + "\n" }
	cases := map[string]struct {
		yaml string
		want string
	}{
		"plain sandbox (unchanged)": {sandboxProfile,
			"  \"mounts\": [" + workspaceMountLine + "],\n"},
		"podman (unchanged)": {sandboxProfile + nested("podman"),
			"  \"mounts\": [\n    " + workspaceMountLine + ",\n    " + containersMount + "\n  ],\n"},
		"docker (unchanged)": {sandboxProfile + nested("docker"),
			"  \"mounts\": [\n    " + workspaceMountLine + ",\n    " + dockerMount + "\n  ],\n"},
		"ollama": {sandboxProfile + ollamaFeature,
			"  \"mounts\": [\n    " + workspaceMountLine + ",\n    " + ollamaMount + "\n  ],\n"},
		"podman and ollama": {sandboxProfile + ollamaFeature + nested("podman"),
			"  \"mounts\": [\n    " + workspaceMountLine + ",\n    " + containersMount + ",\n    " + ollamaMount + "\n  ],\n"},
		"docker and ollama": {sandboxProfile + ollamaFeature + nested("docker"),
			"  \"mounts\": [\n    " + workspaceMountLine + ",\n    " + dockerMount + ",\n    " + ollamaMount + "\n  ],\n"},
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
			if n := strings.Count(string(raw), `"mounts"`); n != 1 {
				t.Errorf("mounts key appears %d times, want 1", n)
			}
		})
	}
}

// Egress hosts of the feature (ADR-0024 decision 7): exactly the two
// measured names, only while the feature is enabled, with the sorted and
// duplicate-free union next to other features. Together with
// `claude-code` and `node` the render keys, the credential reference and
// the model volume appear once each and the second run is a NoOp (golden
// case "ollama together with claude-code").
func TestGenerateDevcontainer_Ollama_EgressHostsAndTogetherWithClaudeCode(t *testing.T) {
	t.Parallel()
	const (
		manifestHost = "registry.ollama.ai"
		blobHost     = "dd20bb891979d25aebc8bec07b2b3bbc.r2.cloudflarestorage.com"
	)
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
		"ollama only": {ollamaFeature, withHosts(manifestHost, blobHost)},
		"disabled":    {"  features:\n    ollama:\n      enabled: false\n", withHosts()},
		"no feature":  {"", withHosts()},
		"claude-code and node": {"  features:\n    ollama:\n      enabled: true\n    claude-code:\n      enabled: true\n    node:\n      enabled: true\n",
			withHosts(manifestHost, blobHost, "api.anthropic.com", "registry.npmjs.org")},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, sandboxProfile+tc.features+egressOn)
			resp, err := generateSandbox(svc, false)
			if err != nil {
				t.Fatalf("generate: %v", err)
			}
			hosts := egressHostList(t, fs)
			if !sort.StringsAreSorted(hosts) || !reflect.DeepEqual(hosts, tc.wantHosts) {
				t.Errorf("egress hosts =\n%v\nwant (sorted, no duplicates)\n%v", hosts, tc.wantHosts)
			}
			raw, _ := fs.ReadFile(devcontainerJSONPath())
			body := string(raw)
			if strings.Contains(tc.features, "ollama:\n      enabled: true") {
				if n := strings.Count(body, ollamaMount); n != 1 {
					t.Errorf("model volume appears %d times, want 1:\n%s", n, body)
				}
				if n := len(ollamaNotices(resp.Warnings)); n != 1 {
					t.Errorf("ollama notices = %d, want 1: %+v", n, resp.Warnings)
				}
			}
			if strings.Contains(tc.features, "claude-code") {
				for _, key := range []string{claudeCodeKey + ":1", "ghcr.io/devcontainers/features/node:1"} {
					if n := strings.Count(body, `"`+key+`"`); n != 1 {
						t.Errorf("render key %q appears %d times, want 1:\n%s", key, n, body)
					}
				}
				if n := strings.Count(body, `"ANTHROPIC_API_KEY"`); n != 1 {
					t.Errorf("credential reference appears %d times, want 1", n)
				}
			}
			// Egress adds only the egress script to the start chain; ollama starts nothing.
			if !strings.Contains(body, `"postStartCommand": "sudo sh /usr/local/bin/u-boot-egress-init"`) || strings.Contains(body, "ollama serve") {
				t.Errorf("postStartCommand must be the egress script alone:\n%s", body)
			}
			if r2, err := generateSandbox(svc, false); err != nil || r2.Action != driving.GenerateActionNoOp {
				t.Errorf("second run: action=%v err=%v, want NoOp", r2.Action, err)
			}
		})
	}
}

// `ollama` is catalogue-activatable (no orphan finding) and renders no key,
// so the drift check stays silent whether or not devcontainer.json has a
// `features` entry for it (ADR-0024 decision 2).
func TestDoctor_FeaturesDrift_Ollama(t *testing.T) {
	t.Parallel()
	for name, jsonBody := range map[string]string{
		"features empty": `{"features": {}}`,
		"no features":    `{"name": "demo"}`,
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs, _, _, _ := newDoctorService(t)
			driftSeed(t, fs, "schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n  features:\n    ollama:\n      enabled: true\n", jsonBody)
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
		})
	}
}
