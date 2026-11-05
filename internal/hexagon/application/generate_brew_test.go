package application_test

import (
	"bytes"
	"errors"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// Golden cases of the catalogue key `brew` (LH-FA-DEV-013): a Dockerfile block
// that checks out Homebrew at a pinned commit and installs the formulas of the
// list as the container user, in both profiles, with no render key under
// `features`, one notice for the block and one per tap, the approval of the
// taps before anything is written, and the conflict with the OCI feature
// `homebrew`.

const (
	// brewPin is the commit of Homebrew/brew the template pins
	// (HOMEBREW_REVISION, release 7.0.12). The template is the one place of the
	// pin; this constant pins the template.
	brewPin = "bf464f57275cae0be106b43d190c98564ac59b0e"

	brewEntryOn     = "    brew:\n      enabled: true\n"
	brewTapAllow    = "  brewTaps:\n    allow:\n      - pt9912/u-boot\n"
	brewTwoFormulas = "    brew:\n      enabled: true\n      packages:\n        - ripgrep\n        - pt9912/u-boot/u-boot\n"
)

// brewBlockLines are the instruction lines of the block for the formulas
// ripgrep and pt9912/u-boot/u-boot: the pin, the build dependencies with the
// profile file that puts Homebrew on the `PATH` of login shells, the
// checkout at the pin with the abort on another commit, the installation as
// the container user with the update and the analytics switched off, and the
// return to root. Comment lines are not part of it.
func brewBlockLines(formulas string) []string {
	lines := []string{
		"ARG HOMEBREW_REVISION=" + brewPin,
		"RUN apt-get update \\",
		" && apt-get install -y --no-install-recommends build-essential procps curl file git ca-certificates \\",
		" && rm -rf /var/lib/apt/lists/* \\",
		" && install -d -o vscode -g vscode /home/linuxbrew/.linuxbrew \\",
		` && printf '%s\n' 'export PATH="/home/linuxbrew/.linuxbrew/bin:/home/linuxbrew/.linuxbrew/sbin:${PATH}"' > /etc/profile.d/u-boot-brew.sh`,
		`ENV PATH="/home/linuxbrew/.linuxbrew/bin:/home/linuxbrew/.linuxbrew/sbin:${PATH}"`,
		"USER vscode",
		"RUN cd /home/linuxbrew/.linuxbrew \\",
		" && git init -q Homebrew \\",
		" && cd Homebrew \\",
		" && git remote add origin https://github.com/Homebrew/brew \\",
		` && git fetch -q --depth 1 origin "${HOMEBREW_REVISION}" \`,
		" && git checkout -q --detach FETCH_HEAD \\",
		` && test "$(git rev-parse HEAD)" = "${HOMEBREW_REVISION}" \`,
		" && mkdir ../bin \\",
		" && ln -s ../Homebrew/bin/brew ../bin/brew",
	}
	if formulas != "" {
		lines = append(lines, "RUN HOMEBREW_NO_AUTO_UPDATE=1 HOMEBREW_NO_ANALYTICS=1 brew install --formula"+formulas)
	}
	return append(lines, "USER root")
}

// brewBlockOf returns the instruction lines of the Homebrew block of df: from
// the pin to the `USER root` that ends it, without comment lines.
func brewBlockOf(t *testing.T, df string) []string {
	t.Helper()
	lines := strings.Split(df, "\n")
	start := -1
	for i, l := range lines {
		if strings.HasPrefix(l, "ARG HOMEBREW_REVISION=") {
			start = i
			break
		}
	}
	if start < 0 {
		t.Fatalf("the Dockerfile carries no Homebrew block:\n%s", df)
	}
	var out []string
	for _, l := range lines[start:] {
		if strings.HasPrefix(l, "#") {
			continue
		}
		out = append(out, l)
		if l == "USER root" {
			return out
		}
	}
	t.Fatalf("the Homebrew block has no closing USER root:\n%s", df)
	return nil
}

// assertBrewBlock pins the instruction lines of the block (the golden) and what
// must be absent from the Dockerfile: no installer script and no download piped
// into a shell, no explicit tap or trust (the approval is the one of u-boot, the
// trust of Homebrew is the one of the fully qualified formula name), no
// sign-in or token, and no `ENV` that carries the update switch into the
// running container.
func assertBrewBlock(t *testing.T, df, formulas string) {
	t.Helper()
	if got, want := brewBlockOf(t, df), brewBlockLines(formulas); !reflect.DeepEqual(got, want) {
		t.Errorf("Homebrew block =\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	if n := strings.Count(df, "ARG HOMEBREW_REVISION="); n != 1 {
		t.Errorf("the pin stands %d times, want once (one place)", n)
	}
	for _, line := range brewBlockOf(t, df) {
		for _, bad := range []string{
			"| sh", "| bash", "|sh", "|bash", "install.sh", "NONINTERACTIVE", "raw.githubusercontent.com/Homebrew",
			"brew tap", "brew trust", "auth login", "TOKEN", "HOMEBREW_GITHUB_API", "ENV HOMEBREW", "sudo",
		} {
			if strings.Contains(line, bad) {
				t.Errorf("block line %q must not contain %q", line, bad)
			}
		}
	}
	if i := strings.Index(df, "USER root"); i < 0 || i > strings.Index(df, "ARG HOMEBREW_REVISION=") {
		t.Errorf("the Homebrew block must stand behind the first USER root:\n%s", df)
	}
	if !strings.HasSuffix(strings.TrimSpace(strings.SplitN(df, "# END U-BOOT MANAGED BLOCK", 2)[0]), "USER vscode") {
		t.Errorf("Dockerfile must end on USER vscode:\n%s", df)
	}
}

func brewNotices(ws []driving.WarningEntry) []driving.WarningEntry {
	var out []driving.WarningEntry
	for _, w := range ws {
		if strings.HasPrefix(w.Message, "devcontainer feature brew installs") {
			out = append(out, w)
		}
	}
	return out
}

// B1, B3 (fresh state): in the default profile, the sandbox profile and with
// `--sandbox` the active key gets the pinned block behind `USER root`, renders
// no key under `features`, and the second run is a NoOp.
func TestGenerateDevcontainer_Brew_BlockInEveryProfile(t *testing.T) {
	t.Parallel()
	for pname, p := range forgeProfiles() {
		t.Run(pname, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, p.head+brewTapAllow+featuresYAML(brewTwoFormulas))
			if _, err := generateSandbox(svc, p.sandbox); err != nil {
				t.Fatalf("generate: %v", err)
			}
			raw, _ := fs.ReadFile(dockerfilePath())
			assertBrewBlock(t, string(raw), " 'ripgrep' 'pt9912/u-boot/u-boot'")
			js, _ := fs.ReadFile(devcontainerJSONPath())
			if feats := parseDevcontainerFeatures(t, js); len(feats) != 0 {
				t.Errorf("brew renders no key under features, got %v", feats)
			}
			if r2, err := generateSandbox(svc, p.sandbox); err != nil || r2.Action != driving.GenerateActionNoOp {
				t.Errorf("second run: action=%v err=%v, want NoOp", r2.Action, err)
			}
		})
	}
}

// B2 (boundary): `enabled: true` with an empty or missing list installs
// Homebrew alone; the list keeps its order and each formula is quoted.
func TestGenerateDevcontainer_Brew_ListForms(t *testing.T) {
	t.Parallel()
	for name, tc := range map[string]struct {
		entry    string
		formulas string
	}{
		"no list":            {brewEntryOn, ""},
		"empty list":         {"    brew:\n      enabled: true\n      packages: []\n", ""},
		"one formula":        {"    brew:\n      enabled: true\n      packages:\n        - ripgrep\n", " 'ripgrep'"},
		"order is the list":  {"    brew:\n      enabled: true\n      packages:\n        - zoxide\n        - bat\n        - node@22\n", " 'zoxide' 'bat' 'node@22'"},
		"special characters": {"    brew:\n      enabled: true\n      packages:\n        - c++filt\n        - python-tk@3.13\n", " 'c++filt' 'python-tk@3.13'"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, defaultProfile+featuresYAML(tc.entry))
			resp, err := generateSandbox(svc, false)
			if err != nil {
				t.Fatalf("generate: %v", err)
			}
			raw, _ := fs.ReadFile(dockerfilePath())
			assertBrewBlock(t, string(raw), tc.formulas)
			if tc.formulas == "" && strings.Contains(string(raw), "brew install --formula") {
				t.Errorf("no formula, no install step:\n%s", raw)
			}
			n := brewNotices(resp.Warnings)
			if len(n) != 1 {
				t.Fatalf("notices = %d, want 1: %+v", len(n), resp.Warnings)
			}
			if tc.formulas == "" && !strings.Contains(n[0].Message, "no formula is listed") {
				t.Errorf("notice for Homebrew alone must say so: %q", n[0].Message)
			}
		})
	}
}

// B4 (boundary): an entry that is not active leaves no trace in any file of
// the project but `u-boot.yaml`: switched off with a list, without `enabled`,
// and with a `source:` (an external feature without the block).
func TestGenerateDevcontainer_Brew_InactiveEntriesLeaveNoTrace(t *testing.T) {
	t.Parallel()
	const external = "ghcr.io/example/features/brew"
	baseline, baseFS := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, baseFS, defaultProfile)
	if _, err := generateSandbox(baseline, false); err != nil {
		t.Fatalf("baseline generate: %v", err)
	}
	wantDF, _ := baseFS.ReadFile(dockerfilePath())
	wantJSON, _ := baseFS.ReadFile(devcontainerJSONPath())
	for name, yaml := range map[string]string{
		"switched off with a list": defaultProfile + featuresYAML("    brew:\n      enabled: false\n      packages:\n        - ripgrep\n        - nobody/tap/formula\n"),
		"without enabled":          defaultProfile + featuresYAML("    brew:\n      packages:\n        - ripgrep\n"),
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, yaml)
			resp, err := generateSandbox(svc, false)
			if err != nil {
				t.Fatalf("generate: %v", err)
			}
			gotDF, _ := fs.ReadFile(dockerfilePath())
			gotJSON, _ := fs.ReadFile(devcontainerJSONPath())
			if !bytes.Equal(gotDF, wantDF) || !bytes.Equal(gotJSON, wantJSON) {
				t.Errorf("an inactive brew entry changed the generated files:\n%s", gotDF)
			}
			if n := len(brewNotices(resp.Warnings)); n != 0 {
				t.Errorf("an inactive entry gives %d notices, want none", n)
			}
		})
	}
	t.Run("with a source", func(t *testing.T) {
		t.Parallel()
		svc, fs := newGenerateService(t)
		seedUBootYAMLWithFeatures(t, fs, defaultProfile+"  featureSources:\n    allow:\n      - "+external+"\n"+
			featuresYAML("    brew:\n      enabled: true\n      source: "+external+"\n"))
		resp, err := generateSandbox(svc, false)
		if err != nil {
			t.Fatalf("generate: %v", err)
		}
		gotDF, _ := fs.ReadFile(dockerfilePath())
		if !bytes.Equal(gotDF, wantDF) {
			t.Errorf("brew with a source is an external feature and adds no block:\n%s", gotDF)
		}
		js, _ := fs.ReadFile(devcontainerJSONPath())
		if feats := parseDevcontainerFeatures(t, js); feats[external+":1"] == nil {
			t.Errorf("features keys = %v, want the external feature", feats)
		}
		if n := len(brewNotices(resp.Warnings)); n != 0 {
			t.Errorf("an external brew feature gives %d block notices, want none", n)
		}
	})
}

// X1 (golden, order): the block stands behind the blocks of the container
// runtime, ollama, gh, glab and the three toolchains, and ahead of the egress
// restriction; the hosts and mounts do not change by it; the second run is a
// NoOp (BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf: the order is
// pinned here, the combination is not built).
func TestGenerateDevcontainer_Brew_OrderAgainstTheOtherBlocks(t *testing.T) {
	t.Parallel()
	const features = "  features:\n    ollama:\n      enabled: true\n    gh:\n      enabled: true\n    glab:\n      enabled: true\n" +
		"    go-toolchain:\n      enabled: true\n    golangci-lint:\n      enabled: true\n    java-toolchain:\n      enabled: true\n" +
		brewTwoFormulas
	for rt, runtimeMarker := range map[string]string{"podman": "ARG PODMAN_VERSION=", "docker": "# Nested Docker daemon"} {
		t.Run(rt, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, sandboxProfile+brewTapAllow+features+"  sandbox:\n    nestedRuntime: "+rt+"\n    egress:\n      enabled: true\n")
			if _, err := generateSandbox(svc, false); err != nil {
				t.Fatalf("generate: %v", err)
			}
			raw, _ := fs.ReadFile(dockerfilePath())
			df := string(raw)
			assertBrewBlock(t, df, " 'ripgrep' 'pt9912/u-boot/u-boot'")
			last := -1
			for _, marker := range []string{
				"USER root", runtimeMarker, "ARG OLLAMA_VERSION=", "ARG GH_VERSION=", "ARG GLAB_VERSION=",
				"ARG GO_VERSION=", "ARG GOLANGCI_LINT_VERSION=", "ARG JDK_VERSION=", "ARG HOMEBREW_REVISION=",
				"# Egress restriction",
			} {
				i := strings.Index(df, marker)
				if i < 0 || i < last {
					t.Errorf("Dockerfile marker %q missing or out of order (index %d, previous %d):\n%s", marker, i, last, df)
				}
				last = i
			}
			// The block hands the user back to root, so the egress block runs as root,
			// and the file ends on the container user.
			if i := strings.Index(df[strings.Index(df, "ARG HOMEBREW_REVISION="):], "\nUSER root\n"); i < 0 ||
				strings.Index(df, "ARG HOMEBREW_REVISION=")+i > strings.Index(df, "# Egress restriction") {
				t.Errorf("the Homebrew block must return to USER root before the egress block:\n%s", df)
			}
			if i := strings.LastIndex(df, "USER vscode"); i < last {
				t.Errorf("Dockerfile must end on USER vscode:\n%s", df)
			}
			js, _ := fs.ReadFile(devcontainerJSONPath())
			if n := strings.Count(string(js), "type=volume"); n != 3 {
				t.Errorf("volumes = %d, want 3 (workspace, runtime, ollama; brew adds none):\n%s", n, js)
			}
			if r2, err := generateSandbox(svc, false); err != nil || r2.Action != driving.GenerateActionNoOp {
				t.Errorf("second run: action=%v err=%v, want NoOp", r2.Action, err)
			}
		})
	}
}

// X2 (default profile): the key changes the Dockerfile and nothing else;
// devcontainer.json is byte-identical to the run without it and the key
// renders no entry under `features`.
func TestGenerateDevcontainer_Brew_DefaultProfileJSONUnchanged(t *testing.T) {
	t.Parallel()
	with, fsWith := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fsWith, defaultProfile+brewTapAllow+featuresYAML(brewTwoFormulas))
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
		t.Errorf("devcontainer.json differs by the brew key in the default profile:\n%s\n---\n%s", jsonWith, jsonWithout)
	}
}

// X3 (replay): enabling the key later replaces the stale managed block
// (UpdatedBlock) and keeps the line the user added behind it; the next run is a
// NoOp; disabling it removes the block and keeps the user line.
func TestGenerateDevcontainer_Brew_ActivateAndDeactivateLater(t *testing.T) {
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
			seedUBootYAMLWithFeatures(t, fs, head+featuresYAML("    brew:\n      enabled: true\n      packages:\n        - ripgrep\n"))
			resp, err := generateSandbox(svc, false)
			if err != nil {
				t.Fatalf("activate: %v", err)
			}
			if resp.Action != driving.GenerateActionUpdatedBlock {
				t.Errorf("activate: action = %v, want UpdatedBlock", resp.Action)
			}
			after, _ := fs.ReadFile(dockerfilePath())
			assertBrewBlock(t, string(after), " 'ripgrep'")
			if !strings.Contains(string(after), "# user addition\nRUN echo user\n") {
				t.Errorf("the user line behind the block was lost:\n%s", after)
			}
			if r2, err := generateSandbox(svc, false); err != nil || r2.Action != driving.GenerateActionNoOp {
				t.Errorf("second run: action=%v err=%v, want NoOp", r2.Action, err)
			}
			seedUBootYAMLWithFeatures(t, fs, head+featuresYAML("    brew:\n      enabled: false\n      packages:\n        - ripgrep\n"))
			if _, err := generateSandbox(svc, false); err != nil {
				t.Fatalf("deactivate: %v", err)
			}
			gone, _ := fs.ReadFile(dockerfilePath())
			if strings.Contains(string(gone), "HOMEBREW_REVISION") {
				t.Errorf("the block must be gone after disabling the key:\n%s", gone)
			}
			if !strings.Contains(string(gone), "# user addition\nRUN echo user\n") {
				t.Errorf("the user line behind the block was lost:\n%s", gone)
			}
		})
	}
}

// N1 (LH-FA-DEV-013, Ausweisung): one notice for the block and one per used
// tap in every profile, code LH-NFA-SEC-004, level warn, no SecurityRelaxation
// and no DegradationFallback (`--quiet` drops them; the filter is the CLI's).
// The block notice names the source, the place of the pin, the formulas and
// that they are not pinned and not checked; it names the commit nowhere. The
// tap notice names the tap as a source of its own.
func TestGenerateDevcontainer_Brew_Notices(t *testing.T) {
	t.Parallel()
	for pname, p := range forgeProfiles() {
		t.Run(pname, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, p.head+"  brewTaps:\n    allow:\n      - pt9912/u-boot\n      - other/tap\n"+
				featuresYAML("    brew:\n      enabled: true\n      packages:\n        - ripgrep\n        - other/tap/b\n        - pt9912/u-boot/u-boot\n        - other/tap/a\n"))
			resp, err := generateSandbox(svc, p.sandbox)
			if err != nil {
				t.Fatalf("generate: %v", err)
			}
			notices := brewNotices(resp.Warnings)
			if len(notices) != 3 {
				t.Fatalf("notices = %d, want 3 (block, tap other/tap, tap pt9912/u-boot): %+v", len(notices), resp.Warnings)
			}
			for _, w := range notices {
				if w.Code != "LH-NFA-SEC-004" || w.Level != "warn" || w.Subject != ".devcontainer/devcontainer.json" {
					t.Errorf("notice = %+v, want a warn entry with code LH-NFA-SEC-004 on devcontainer.json", w)
				}
				if w.SecurityRelaxation || w.DegradationFallback {
					t.Errorf("notice must be unflagged (dropped under --quiet): %+v", w)
				}
				if strings.Contains(w.Message, brewPin) || strings.Contains(w.Message, "7.0.12") {
					t.Errorf("the notice names the pin, which stands once, in the Dockerfile: %q", w.Message)
				}
			}
			for _, part := range []string{
				"github.com/Homebrew/brew", "HOMEBREW_REVISION in .devcontainer/Dockerfile", "without an installer script",
				"ripgrep, other/tap/b, pt9912/u-boot/u-boot, other/tap/a", "not pinned and not checked",
			} {
				if !strings.Contains(notices[0].Message, part) {
					t.Errorf("block notice lacks %q: %q", part, notices[0].Message)
				}
			}
			for i, tap := range []string{"other/tap", "pt9912/u-boot"} {
				m := notices[i+1].Message
				if !strings.Contains(m, "the tap "+tap+" ") || !strings.Contains(m, "source of its own") || !strings.Contains(m, "not pinned and not checked") {
					t.Errorf("tap notice %d for %s = %q", i, tap, m)
				}
			}
			// A formula without a tap adds no tap notice.
			svc2, fs2 := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs2, p.head+featuresYAML("    brew:\n      enabled: true\n      packages:\n        - ripgrep\n"))
			resp2, err := generateSandbox(svc2, p.sandbox)
			if err != nil {
				t.Fatalf("generate: %v", err)
			}
			if n := len(brewNotices(resp2.Warnings)); n != 1 {
				t.Errorf("notices for a core formula = %d, want 1", n)
			}
		})
	}
}

// T1 (LH-FA-DEV-013, Tap nur nach Freigabe): a formula of a tap that is not in
// `devcontainer.brewTaps.allow` ends with ErrConfigValueInvalid (exit code 10)
// and changes nothing — no file, no directory, `u-boot.yaml` byte-identical —
// with and without `--dry-run`. The message names the tap, the formula and the
// ways out; it is the same text the doctor check gives.
func TestGenerateDevcontainer_Brew_UnapprovedTapIsRejectedBeforeWriting(t *testing.T) {
	t.Parallel()
	const yaml = defaultProfile + "  features:\n    brew:\n      enabled: true\n      packages:\n        - ripgrep\n        - pt9912/u-boot/u-boot\n"
	for name, p := range map[string]string{"default profile": defaultProfile, "sandbox profile": sandboxProfile} {
		for _, mode := range []driving.PreviewMode{driving.PreviewNone, driving.PreviewDryRun} {
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, strings.Replace(yaml, defaultProfile, p, 1))
			before, _ := fs.ReadFile(freigabeYAMLPath())
			writes, mkdirs := len(fs.writes), len(fs.mkdirs)
			_, err := generateWithFlag(svc, mode)
			if !errors.Is(err, driving.ErrConfigValueInvalid) {
				t.Fatalf("%s mode %v: err = %v, want ErrConfigValueInvalid", name, mode, err)
			}
			for _, part := range []string{
				`the tap "pt9912/u-boot"`, `"pt9912/u-boot/u-boot"`, "devcontainer.brewTaps.allow",
				"`u-boot config set devcontainer.brewTaps.allow pt9912/u-boot`", "devcontainer.features.brew.enabled false", "--yes does not substitute",
			} {
				if !strings.Contains(err.Error(), part) {
					t.Errorf("%s mode %v: message lacks %q: %v", name, mode, part, err)
				}
			}
			if len(fs.writes) != writes || len(fs.mkdirs) != mkdirs {
				t.Errorf("%s mode %v: rejected call wrote: writes %v, mkdirs %v", name, mode, fs.writes[writes:], fs.mkdirs[mkdirs:])
			}
			for path := range fs.files {
				if strings.Contains(path, ".devcontainer") {
					t.Errorf("%s mode %v: rejected call left %q", name, mode, path)
				}
			}
			if after, _ := fs.ReadFile(freigabeYAMLPath()); !bytes.Equal(before, after) {
				t.Errorf("%s mode %v: u-boot.yaml changed:\n%s", name, mode, after)
			}
		}
	}
}

// T2 (boundary): the approval is byte-equal and its own list. A near name, a
// case variant, the entry in `featureSources.allow` and the flag
// `--allow-external-feature-sources` of the call free no tap; the exact entry
// frees it, and a formula without a tap needs no approval.
func TestGenerateDevcontainer_Brew_ApprovalIsByteEqualAndItsOwnList(t *testing.T) {
	t.Parallel()
	const entry = "  features:\n    brew:\n      enabled: true\n      packages:\n        - pt9912/u-boot/u-boot\n"
	for name, tc := range map[string]struct {
		extra   string
		flag    []string
		wantErr bool
	}{
		"exact entry frees":                    {"  brewTaps:\n    allow:\n      - pt9912/u-boot\n", nil, false},
		"another tap frees nothing":            {"  brewTaps:\n    allow:\n      - pt9912/other\n", nil, true},
		"case variant frees nothing":           {"  brewTaps:\n    allow:\n      - Pt9912/u-boot\n", nil, true},
		"prefix frees nothing":                 {"  brewTaps:\n    allow:\n      - pt9912/u-boo\n", nil, true},
		"featureSources frees nothing":         {"  featureSources:\n    allow:\n      - pt9912/u-boot\n", nil, true},
		"the flag of the call frees nothing":   {"", []string{"pt9912/u-boot"}, true},
		"empty approval list frees nothing":    {"  brewTaps:\n    allow: []\n", nil, true},
		"list with both entries frees the tap": {"  brewTaps:\n    allow:\n      - pt9912/other\n      - pt9912/u-boot\n", nil, false},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, defaultProfile+tc.extra+entry)
			_, err := generateWithFlag(svc, driving.PreviewNone, tc.flag...)
			if tc.wantErr != errors.Is(err, driving.ErrConfigValueInvalid) {
				t.Fatalf("err = %v, want rejection = %v", err, tc.wantErr)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("generate: %v", err)
			}
			if exists, _ := fs.Exists(dockerfilePath()); exists == tc.wantErr {
				t.Errorf("Dockerfile exists = %v after a call with rejection = %v", exists, tc.wantErr)
			}
		})
	}
	t.Run("a formula without a tap needs no approval", func(t *testing.T) {
		t.Parallel()
		svc, fs := newGenerateService(t)
		seedUBootYAMLWithFeatures(t, fs, defaultProfile+featuresYAML("    brew:\n      enabled: true\n      packages:\n        - ripgrep\n"))
		if _, err := generateWithFlag(svc, driving.PreviewNone); err != nil {
			t.Fatalf("generate: %v", err)
		}
	})
	t.Run("an inactive entry is not asked", func(t *testing.T) {
		t.Parallel()
		svc, fs := newGenerateService(t)
		seedUBootYAMLWithFeatures(t, fs, defaultProfile+featuresYAML("    brew:\n      enabled: false\n      packages:\n        - pt9912/u-boot/u-boot\n"))
		if _, err := generateWithFlag(svc, driving.PreviewNone); err != nil {
			t.Fatalf("generate: %v", err)
		}
	})
	t.Run("several taps are named in sorted order", func(t *testing.T) {
		t.Parallel()
		svc, fs := newGenerateService(t)
		seedUBootYAMLWithFeatures(t, fs, defaultProfile+featuresYAML("    brew:\n      enabled: true\n      packages:\n        - z/z/f\n        - a/a/f\n        - a/a/g\n"))
		_, err := generateWithFlag(svc, driving.PreviewNone)
		if !errors.Is(err, driving.ErrConfigValueInvalid) {
			t.Fatalf("err = %v, want ErrConfigValueInvalid", err)
		}
		if a, z := strings.Index(err.Error(), `"a/a"`), strings.Index(err.Error(), `"z/z"`); a < 0 || z < a {
			t.Errorf("taps are not named in sorted order: %v", err)
		}
		if !strings.Contains(err.Error(), `"a/a/f", "a/a/g"`) {
			t.Errorf("the formulas of one tap are not grouped: %v", err)
		}
	})
}

// C1 (LH-FA-DEV-013, Konflikt): `brew` next to the enabled OCI feature
// `homebrew`, with and without a `source:` at `homebrew`, ends with
// ErrConfigValueInvalid (exit code 10) before anything is written; the message
// names both keys and both ways out. Either key switched off lets the run pass.
func TestGenerateDevcontainer_Brew_ConflictWithTheHomebrewFeature(t *testing.T) {
	t.Parallel()
	const external = "ghcr.io/example/features/homebrew"
	for name, tc := range map[string]struct {
		features string
		extra    string
		wantErr  bool
	}{
		"both active":                    {brewEntryOn + "    homebrew:\n      enabled: true\n", "", true},
		"homebrew with a source":         {brewEntryOn + "    homebrew:\n      enabled: true\n      source: " + external + "\n", "  featureSources:\n    allow:\n      - " + external + "\n", true},
		"homebrew with a version":        {brewEntryOn + "    homebrew:\n      enabled: true\n      version: \"1\"\n", "", true},
		"brew switched off":              {"    brew:\n      enabled: false\n    homebrew:\n      enabled: true\n", "", false},
		"homebrew switched off":          {brewEntryOn + "    homebrew:\n      enabled: false\n", "", false},
		"homebrew without enabled":       {brewEntryOn + "    homebrew: {}\n", "", false},
		"homebrew alone":                 {"    homebrew:\n      enabled: true\n", "", false},
		"brew with a source is external": {"    brew:\n      enabled: true\n      source: ghcr.io/example/features/brew\n    homebrew:\n      enabled: true\n", "  featureSources:\n    allow:\n      - ghcr.io/example/features/brew\n", false},
	} {
		for _, mode := range []driving.PreviewMode{driving.PreviewNone, driving.PreviewDryRun} {
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				svc, fs := newGenerateService(t)
				seedUBootYAMLWithFeatures(t, fs, defaultProfile+tc.extra+featuresYAML(tc.features))
				writes := len(fs.writes)
				_, err := generateWithFlag(svc, mode)
				if !tc.wantErr {
					if err != nil {
						t.Fatalf("generate: %v", err)
					}
					return
				}
				if !errors.Is(err, driving.ErrConfigValueInvalid) {
					t.Fatalf("err = %v, want ErrConfigValueInvalid", err)
				}
				for _, part := range []string{
					"devcontainer.features.brew", "devcontainer.features.homebrew",
					"`u-boot config set devcontainer.features.brew.enabled false`", "`u-boot config set devcontainer.features.homebrew.enabled false`",
				} {
					if !strings.Contains(err.Error(), part) {
						t.Errorf("message lacks %q: %v", part, err)
					}
				}
				if len(fs.writes) != writes {
					t.Errorf("rejected call wrote: %v", fs.writes[writes:])
				}
				if exists, _ := fs.Exists(dockerfilePath()); exists {
					t.Errorf("rejected call left a Dockerfile")
				}
			})
		}
	}
}

// R1 (form, read path): an entry outside the form, a package list at another
// key, a list at `brew` with a `source:`, a bad tap entry and a `version:` end
// with ErrGenerateManualConflict (exit code 10) and write nothing; the message
// names the field.
func TestGenerateDevcontainer_Brew_InvalidEntriesExit10(t *testing.T) {
	t.Parallel()
	pk := func(entry string) string {
		return "    brew:\n      enabled: true\n      packages:\n        - " + entry + "\n"
	}
	cases := map[string]struct {
		yaml     string
		wantPart string
		formErr  bool // the cause is an entry outside the form (the message says so)
	}{
		"space":               {pk("rip grep"), "features.brew.packages[0]", true},
		"shell character":     {pk("'ripgrep;id'"), "features.brew.packages[0]", true},
		"leading dash":        {pk("--HEAD"), "features.brew.packages[0]", true},
		"two segments":        {pk("pt9912/u-boot"), "features.brew.packages[0]", true},
		"four segments":       {pk("a/b/c/d"), "features.brew.packages[0]", true},
		"url":                 {pk("https://github.com/a/b"), "features.brew.packages[0]", true},
		"second entry bad":    {"    brew:\n      enabled: true\n      packages:\n        - ripgrep\n        - bad entry\n", "features.brew.packages[1]", true},
		"switched off, bad":   {"    brew:\n      enabled: false\n      packages:\n        - bad entry\n", "features.brew.packages[0]", true},
		"list at node":        {"    node:\n      enabled: true\n      packages:\n        - ripgrep\n", "devcontainer.features.node.packages", false},
		"list at homebrew":    {"    homebrew:\n      enabled: true\n      packages:\n        - ripgrep\n", "devcontainer.features.homebrew.packages", false},
		"list at brew+source": {"    brew:\n      enabled: true\n      source: ghcr.io/example/features/brew\n      packages:\n        - ripgrep\n", "devcontainer.features.brew.packages", false},
		"version at brew":     {"    brew:\n      enabled: true\n      version: \"5\"\n", "features.brew.version", false},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs := newGenerateService(t)
			seedUBootYAMLWithFeatures(t, fs, defaultProfile+"  featureSources:\n    allow:\n      - ghcr.io/example/features/brew\n"+featuresYAML(tc.yaml))
			_, err := generateSandbox(svc, false)
			if !errors.Is(err, driving.ErrGenerateManualConflict) {
				t.Fatalf("err = %v, want ErrGenerateManualConflict (exit code 10)", err)
			}
			if !strings.Contains(err.Error(), tc.wantPart) {
				t.Errorf("error does not name %q: %v", tc.wantPart, err)
			}
			if got := strings.Contains(err.Error(), "invalid homebrew package"); got != tc.formErr {
				t.Errorf("error names an entry outside the form = %v, want %v: %v", got, tc.formErr, err)
			}
			if exists, _ := fs.Exists(dockerfilePath()); exists {
				t.Errorf("nothing may be written on a validation error")
			}
		})
	}
	t.Run("bad tap in the approval list", func(t *testing.T) {
		t.Parallel()
		svc, fs := newGenerateService(t)
		seedUBootYAMLWithFeatures(t, fs, defaultProfile+"  brewTaps:\n    allow:\n      - pt9912\n")
		_, err := generateSandbox(svc, false)
		if !errors.Is(err, driving.ErrGenerateManualConflict) || !strings.Contains(err.Error(), "brewTaps.allow[0]") {
			t.Fatalf("err = %v, want ErrGenerateManualConflict naming brewTaps.allow[0]", err)
		}
	})
}

// The validation does not read `enabled`: a switched-off key keeps a valid list
// valid, and the way out of every rejection above is the file or `config set`.
// The sorted key order makes the first error deterministic.
func TestGenerateDevcontainer_Brew_FirstErrorIsDeterministic(t *testing.T) {
	t.Parallel()
	for i := 0; i < 5; i++ {
		svc, fs := newGenerateService(t)
		seedUBootYAMLWithFeatures(t, fs, defaultProfile+featuresYAML(
			"    zzz:\n      enabled: true\n      packages:\n        - a\n    node:\n      enabled: true\n      packages:\n        - a\n"))
		_, err := generateSandbox(svc, false)
		if err == nil || !strings.Contains(err.Error(), "features.node.packages") {
			t.Fatalf("iter %d: err = %v, want the first key by sort order (node)", i, err)
		}
	}
}

// H1 (hosts): the hosts of the egress restriction with `brew` differ from the
// run without it by the three hosts that Homebrew needs in the running
// container (docs/messungen/slice-v1-devcontainer-homebrew-pakete, M4): the
// API host, the registry and the blob host of the bottles. The hosts of a
// tap formula are not listed. A switched-off key adds none.
func TestGenerateDevcontainer_Brew_EgressHosts(t *testing.T) {
	t.Parallel()
	hostsOf := func(features string) []string {
		svc, fs := newGenerateService(t)
		seedUBootYAMLWithFeatures(t, fs, sandboxProfile+brewTapAllow+features+egressOn)
		if _, err := generateSandbox(svc, false); err != nil {
			t.Fatalf("generate: %v", err)
		}
		return egressHostList(t, fs)
	}
	without := hostsOf("")
	with := hostsOf(featuresYAML(brewTwoFormulas))
	added := []string{}
	for _, h := range with {
		if !contains(without, h) {
			added = append(added, h)
		}
	}
	sort.Strings(added)
	if want := []string{"formulae.brew.sh", "ghcr.io", "pkg-containers.githubusercontent.com"}; !reflect.DeepEqual(added, want) {
		t.Errorf("hosts added by brew = %v, want %v", added, want)
	}
	if !sort.StringsAreSorted(with) {
		t.Errorf("hosts are not sorted: %v", with)
	}
	if off := hostsOf(featuresYAML("    brew:\n      enabled: false\n      packages:\n        - ripgrep\n")); !reflect.DeepEqual(off, without) {
		t.Errorf("a switched-off brew key changes the hosts: %v vs %v", off, without)
	}
	for _, h := range []string{"release-assets.githubusercontent.com", "github.com"} {
		if contains(added, h) {
			t.Errorf("the row of brew lists %s, which depends on the formula or is in the base list", h)
		}
	}
}

// H2 (`u-boot egress list`): the three hosts carry the origin `feature:brew`;
// with a nested runtime `ghcr.io` keeps the origin `runtime` in front of it.
func TestEgressList_BrewOrigins(t *testing.T) {
	t.Parallel()
	_, fs := newGenerateService(t)
	seedUBootYAMLWithFeatures(t, fs, egressConfig("    brew:\n      enabled: true\n", "    nestedRuntime: podman\n"))
	origins := egressOrigins(listEgress(t, fs))
	want := map[string][]string{
		"formulae.brew.sh":                     {"feature:brew"},
		"ghcr.io":                              {"runtime", "feature:brew"},
		"pkg-containers.githubusercontent.com": {"feature:brew"},
	}
	for host, w := range want {
		if !reflect.DeepEqual(origins[host], w) {
			t.Errorf("origins of %s = %v, want %v", host, origins[host], w)
		}
	}
}
