//go:build docker

// End-to-end check of the catalogue key `brew` (LH-FA-DEV-013, LH-FA-DEV-008):
// the generated Dockerfile builds against the real engine in both profiles,
// Homebrew stands at the pinned commit and the formulas of the list are
// installed as the container user (a formula of Homebrew itself and a formula
// of a tap), `brew` is on the `PATH` of the sessions the product starts, a
// pin that does not exist breaks the build, the build switch keeps the
// revision when the automatic update would move it, and a formula of Homebrew
// installs behind the egress restriction with the hosts of the generated list.
// The measurement behind the promise: docs/messungen/slice-v1-devcontainer-
// homebrew-pakete. The build downloads from github.com, ghcr.io and
// formulae.brew.sh; a network failure turns the tests red without a change of
// the code. The test runs on the architecture of the engine; `arm64` is not
// measured.

package e2e_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	confirmadapter "github.com/pt9912/u-boot/internal/adapter/driven/confirm"
	fsadapter "github.com/pt9912/u-boot/internal/adapter/driven/fs"
	gitadapter "github.com/pt9912/u-boot/internal/adapter/driven/git"
	loggeradapter "github.com/pt9912/u-boot/internal/adapter/driven/logger"
	progressadapter "github.com/pt9912/u-boot/internal/adapter/driven/progress"
	yamladapter "github.com/pt9912/u-boot/internal/adapter/driven/yaml"
	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

const (
	// The commit of Homebrew/brew that the template pins (release 7.0.12) and
	// the commit of the release before it (7.0.10), read with `git ls-remote
	// --tags https://github.com/Homebrew/brew`. Both age with the releases; a
	// red test names them in its failure.
	brewPin    = "bf464f57275cae0be106b43d190c98564ac59b0e"
	brewPinOld = "8f10429b56be28a25032b87463fffc95e556b5bc"
	// A formula of a tap that exists: the tap of this project (its formula
	// installs the binary of the release). The tap name is the approval entry.
	brewTapFormula = "pt9912/u-boot/u-boot"
	brewTap        = "pt9912/u-boot"
)

// kv is one `config set` of the project set-up, applied in the given order (a
// map would shuffle the order of the tap list and the formulas).
type kv struct{ path, value string }

// dcGenerateOrderedProject mirrors dcGenerateProject with the settings in
// their order: `init --devcontainer`, the `config set` calls, then `generate
// devcontainer`, all with the production adapters.
func dcGenerateOrderedProject(ctx context.Context, t *testing.T, sandbox bool, settings []kv) string {
	t.Helper()
	dir := dcProjectDir(t)
	fsys := fsadapter.New()
	yaml := yamladapter.New()
	logger := loggeradapter.New(os.Stderr, loggeradapter.FormatText, nil)
	initSvc := application.NewInitProjectService(fsys, yaml, gitadapter.New(),
		progressadapter.NewText(os.Stderr), confirmadapter.New(strings.NewReader(""), os.Stderr), logger)
	if _, err := initSvc.Init(ctx, driving.InitProjectRequest{
		BaseDir: dir, Name: filepath.Base(dir), SkipGit: true, Devcontainer: true, Sandbox: sandbox, SilenceProgress: true,
	}); err != nil {
		t.Fatalf("init --devcontainer (sandbox=%v): %v", sandbox, err)
	}
	cfgSvc := application.NewConfigService(fsys, yaml, logger)
	for _, s := range settings {
		p, err := domain.NewConfigPath(s.path)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := cfgSvc.Set(ctx, driving.ConfigSetRequest{BaseDir: dir, Path: p, Value: s.value}); err != nil {
			t.Fatalf("config set %s: %v", s.path, err)
		}
	}
	genSvc := application.NewGenerateService(fsys, yaml, logger)
	if _, err := genSvc.Generate(ctx, driving.GenerateRequest{BaseDir: dir, Artifact: domain.ArtifactDevcontainer}); err != nil {
		t.Fatalf("generate devcontainer: %v", err)
	}
	if err := os.Chmod(dir, 0o777); err != nil {
		t.Fatal(err)
	}
	return dir
}

func brewSettings(extra ...kv) []kv {
	return append([]kv{
		{"devcontainer.brewTaps.allow", brewTap},
		{"devcontainer.features.brew.enabled", "true"},
		{"devcontainer.features.brew.packages", "ripgrep," + brewTapFormula},
	}, extra...)
}

// brewImageScript runs as the image user: who runs, where `brew` stands, the
// revision of the Homebrew checkout, the owner of the prefix, the installed
// formulas and the programs of the list.
const brewImageScript = `echo "USER:$(id -un)"
echo "BREW:$(command -v brew)"
echo "REVISION:$(git -C /home/linuxbrew/.linuxbrew/Homebrew rev-parse HEAD)"
echo "OWNER:$(stat -c %U /home/linuxbrew/.linuxbrew)"
echo "FORMULAS:$(brew list --formula | tr '\n' ' ')"
echo "RG:$(rg --version | head -1)"
echo "TAP:$(u-boot --version | head -1)"
echo "NOENV:$(env | grep -c '^HOMEBREW')"`

// LH-FA-DEV-013 (B1, B3): the Dockerfile of the default profile and of the
// sandbox profile (with the egress restriction, which the build does not run
// behind) builds with a formula of Homebrew itself and a formula of an
// approved tap; Homebrew stands at the pinned commit, `rg` and the program of
// the tap formula run as the container user, the image carries no `HOMEBREW`
// variable (the build switches are no `ENV`), and `brew` is on the `PATH` of a
// non-login shell and of the session `docker exec` starts. Measured and
// reported, not promised: the login shell of `bash` and `sh` resets the `PATH`
// (Debian `/etc/profile`).
func TestE2E_Brew_BuildsInBothProfiles(t *testing.T) {
	requireDocker(t)
	for name, tc := range map[string]struct {
		sandbox bool
		extra   []kv
	}{
		"default profile": {false, nil},
		"sandbox profile": {true, []kv{{"devcontainer.sandbox.egress.enabled", "true"}}},
	} {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
			defer cancel()
			dir := dcGenerateOrderedProject(ctx, t, tc.sandbox, brewSettings(tc.extra...))
			image := buildProjectImage(ctx, t, dir, "uboot-e2e-brew")
			out, code := dockerOutput(ctx, t, "run", "--rm", image, "bash", "-c", brewImageScript)
			if code != 0 {
				t.Fatalf("image checks exit %d:\n%s", code, out)
			}
			for _, want := range []string{
				"USER:vscode",
				"BREW:/home/linuxbrew/.linuxbrew/bin/brew",
				"REVISION:" + brewPin,
				"OWNER:vscode",
				"RG:ripgrep ",
				"TAP:u-boot version ",
				"NOENV:0",
			} {
				if !strings.Contains(out, want) {
					t.Errorf("image checks lack %q:\n%s", want, out)
				}
			}
			if f := lineValue(out, "FORMULAS:"); !strings.Contains(" "+f, " ripgrep ") || !strings.Contains(" "+f, " u-boot ") {
				t.Errorf("formulas = %q, want ripgrep and u-boot among them", f)
			}
			for _, shell := range [][]string{{"bash", "-c"}, {"sh", "-c"}} {
				args := append(append([]string{"run", "--rm", image}, shell...), "command -v brew")
				got, c := dockerOutput(ctx, t, args...)
				if c != 0 || strings.TrimSpace(got) != "/home/linuxbrew/.linuxbrew/bin/brew" {
					t.Errorf("%v: brew on the PATH = %q (exit %d), want the prefix of Homebrew", shell, got, c)
				}
			}
		})
	}
}

func lineValue(out, prefix string) string {
	for _, l := range strings.Split(out, "\n") {
		if v, ok := strings.CutPrefix(l, prefix); ok {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

// LH-FA-DEV-013 (B2, pin): a pin that does not exist breaks the build at the
// checkout (`--build-arg` replaces the value of the template), and the failed
// build leaves no image.
func TestE2E_Brew_PinThatDoesNotExistBreaksTheBuild(t *testing.T) {
	requireDocker(t)
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := dcGenerateOrderedProject(ctx, t, false, []kv{{"devcontainer.features.brew.enabled", "true"}})
	tag := "uboot-e2e-brew-badpin:" + strings.ReplaceAll(filepath.Base(dir), "/", "-")
	out, code := dockerOutput(ctx, t, "build", "-q", "--build-arg", "HOMEBREW_REVISION=0000000000000000000000000000000000000001",
		"-t", tag, filepath.Join(dir, ".devcontainer"))
	if code == 0 {
		t.Cleanup(func() { _, _ = dockerOutput(context.Background(), t, "rmi", "-f", tag) })
		t.Fatalf("build with a pin that does not exist succeeded, want a failure at the checkout:\n%s", out)
	}
	if !strings.Contains(out, "git fetch") {
		t.Errorf("the failure is not the checkout of the pin:\n%.1500s", out)
	}
	if _, c := dockerOutput(ctx, t, "image", "inspect", tag); c == 0 {
		t.Errorf("the failed build left the image %s", tag)
	}
}

// brewUpdateScript sets the age of FETCH_HEAD to three days (the update
// automatism asks it) and installs a formula; the first argument is the
// environment of the installation, as the block sets it or empty.
const brewUpdateScript = `repo=/home/linuxbrew/.linuxbrew/Homebrew
echo "BEFORE:$(git -C "$repo" rev-parse HEAD)"
touch -d "3 days ago" "$repo/.git/FETCH_HEAD"
env $ENVS brew install --formula jq > /tmp/install.log 2>&1
echo "INSTALL_RC:$?"
echo "AFTER:$(git -C "$repo" rev-parse HEAD)"`

// LH-FA-DEV-013 (B2, measurement M3): the switch of the block keeps the
// revision. Built on the commit of the release before the pin, a stale
// FETCH_HEAD lets `brew install` move the checkout to the newest release when
// the switch is missing (the control), and keeps it when the block's
// environment is set. The test fails if the control does not move: then the
// probe measures nothing.
func TestE2E_Brew_UpdateSwitchKeepsTheRevision(t *testing.T) {
	requireDocker(t)
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := dcGenerateOrderedProject(ctx, t, false, []kv{{"devcontainer.features.brew.enabled", "true"}})
	tag := "uboot-e2e-brew-old:" + strings.ReplaceAll(filepath.Base(dir), "/", "-")
	if out, code := dockerOutput(ctx, t, "build", "-q", "--build-arg", "HOMEBREW_REVISION="+brewPinOld, "-t", tag, filepath.Join(dir, ".devcontainer")); code != 0 {
		t.Fatalf("build on the older commit failed:\n%s", out)
	}
	t.Cleanup(func() { _, _ = dockerOutput(context.Background(), t, "rmi", "-f", tag) })
	run := func(envs string) string {
		out, code := dockerOutput(ctx, t, "run", "--rm", "-e", "ENVS="+envs, tag, "bash", "-c", brewUpdateScript)
		if code != 0 {
			t.Fatalf("update probe (%q) exit %d:\n%s", envs, code, out)
		}
		if lineValue(out, "BEFORE:") != brewPinOld || lineValue(out, "INSTALL_RC:") != "0" {
			t.Fatalf("update probe (%q) did not start on the older commit or the install failed:\n%s", envs, out)
		}
		return lineValue(out, "AFTER:")
	}
	if after := run("HOMEBREW_NO_ANALYTICS=1"); after == brewPinOld {
		t.Fatalf("the control (no switch) kept the older commit; the probe measures nothing:\n%s", after)
	}
	if after := run("HOMEBREW_NO_AUTO_UPDATE=1 HOMEBREW_NO_ANALYTICS=1"); after != brewPinOld {
		t.Errorf("with the switch of the block the revision moved to %s, want %s", after, brewPinOld)
	}
}

// brewEgressScript runs as `vscode` in the container behind the egress
// restriction: the list of the start script, the install of a formula of
// Homebrew itself, and the install of another formula (its metadata and bottle
// are not cached) after two of the three hosts of Homebrew were removed from
// the list (the control).
const brewEgressScript = `sudo sh /usr/local/bin/u-boot-egress-init
hosts="$(grep '^hosts=' /usr/local/bin/u-boot-egress-init | cut -d'"' -f2)"
echo "HOSTS:$hosts"
if brew install --formula jq > /tmp/allowed.log 2>&1; then echo INSTALL_ALLOWED_OK; else echo INSTALL_ALLOWED_FAILED; tail -5 /tmp/allowed.log; fi
rest="$(echo "$hosts" | tr ' ' '\n' | grep -v -x -e formulae.brew.sh -e pkg-containers.githubusercontent.com | tr '\n' ' ')"
sudo sed -i "s|^hosts=.*|hosts=\"$rest\"|" /usr/local/bin/u-boot-egress-init
sudo sh /usr/local/bin/u-boot-egress-init
brew install --formula tree > /tmp/blocked.log 2>&1
echo "INSTALL_BLOCKED_RC:$?"`

// LH-FA-DEV-008, LH-FA-DEV-013: behind the egress restriction the start script
// lists the three hosts of Homebrew, `brew install` of a formula of Homebrew
// itself succeeds with them, and the same call fails once two of the hosts are
// out of the list. The hosts of a tap formula are not part of the promise (they
// follow the formula; measured: the release host of the one formula is not in
// the list).
func TestE2E_Brew_EgressHostsOfHomebrew(t *testing.T) {
	requireDocker(t)
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := dcGenerateOrderedProject(ctx, t, true, []kv{
		{"devcontainer.features.brew.enabled", "true"},
		{"devcontainer.sandbox.egress.enabled", "true"},
		{"devcontainer.sandbox.onUnavailable", "warn"},
	})
	image := buildProjectImage(ctx, t, dir, "uboot-e2e-brew-egress")
	args := append(append([]string{"run", "--rm"}, generatedRunArgs(t, dir)...), image, "bash", "-c", brewEgressScript)
	out, code := dockerOutput(ctx, t, args...)
	if code != 0 {
		t.Fatalf("egress script exit %d:\n%s", code, out)
	}
	hostsLine := lineValue(out, "HOSTS:")
	for _, host := range []string{"formulae.brew.sh", "ghcr.io", "pkg-containers.githubusercontent.com"} {
		if !strings.Contains(" "+hostsLine+" ", " "+host+" ") {
			t.Errorf("the allowlist of the start script lacks %s: %q", host, hostsLine)
		}
	}
	if !strings.Contains(out, "INSTALL_ALLOWED_OK") {
		t.Errorf("brew install did not succeed with the hosts of the generated list:\n%s", out)
	}
	if rc := lineValue(out, "INSTALL_BLOCKED_RC:"); rc == "" || rc == "0" {
		t.Errorf("brew install succeeded without two of the hosts (rc %q), the control measures nothing:\n%s", rc, out)
	}
}
