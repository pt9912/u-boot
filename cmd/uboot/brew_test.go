package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The Homebrew block at the user path (real adapters, LH-FA-DEV-013): the tap
// is a source of its own that needs the approval of `devcontainer.brewTaps.allow`,
// `--yes` and `--no-interactive` do not replace it, the rejection ends with
// exit code 10 before anything is written (also under `--dry-run`), the notices
// of the block and of each tap appear once and `--quiet` drops them.

const brewProjectHead = "schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n"

func runBrew(t *testing.T, args ...string) (code int, stdout, stderr string) {
	t.Helper()
	var out, errOut bytes.Buffer
	code = run(context.Background(), args, &out, &errOut)
	return code, out.String(), errOut.String()
}

func mustBrew(t *testing.T, args ...string) string {
	t.Helper()
	code, out, errOut := runBrew(t, args...)
	if code != 0 {
		t.Fatalf("%v: exit %d, stdout=%q stderr=%q", args, code, out, errOut)
	}
	return out + errOut
}

// A tap formula without the approval: `generate devcontainer` ends with exit
// code 10 and writes nothing, also with `--yes`, with `--no-interactive`, with
// the flag `--allow-external-feature-sources` naming the tap and under
// `--dry-run --json` (the envelope carries the code of the failing step and the
// command that approves the tap). The approval lets the same project generate.
func TestRun_Generate_BrewTapNeedsApprovalAndYesDoesNotSubstitute(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	dir := writeFeatureProject(t, brewProjectHead+"  features:\n    brew:\n      enabled: true\n      packages:\n        - ripgrep\n        - pt9912/u-boot/u-boot\n")
	before, _ := os.ReadFile(filepath.Join(dir, "u-boot.yaml"))
	for _, args := range [][]string{
		{"generate", "devcontainer"},
		{"--yes", "generate", "devcontainer"},
		{"--no-interactive", "generate", "devcontainer"},
		{"generate", "devcontainer", "--dry-run"},
		{"generate", "devcontainer", "--allow-external-feature-sources", "pt9912/u-boot"},
	} {
		code, out, errOut := runBrew(t, args...)
		if code != 10 {
			t.Fatalf("%v: exit %d, want 10; stdout=%q stderr=%q", args, code, out, errOut)
		}
		if _, err := os.Stat(".devcontainer"); !os.IsNotExist(err) {
			t.Fatalf("%v: the rejected run must write nothing; .devcontainer exists (%v)", args, err)
		}
	}
	code, out, errOut := runBrew(t, "generate", "devcontainer", "--dry-run", "--json")
	if code != 10 || !strings.Contains(out, "LH-FA-DEV-003") ||
		!strings.Contains(out, "u-boot config set devcontainer.brewTaps.allow pt9912/u-boot") {
		t.Fatalf("--dry-run --json: exit %d, stdout=%q stderr=%q, want 10 with the code and the approving command", code, out, errOut)
	}
	if after, _ := os.ReadFile(filepath.Join(dir, "u-boot.yaml")); !bytes.Equal(before, after) {
		t.Fatalf("u-boot.yaml changed by the rejected runs (the flag must not append to the list):\n%s", after)
	}

	mustBrew(t, "config", "set", "devcontainer.brewTaps.allow", "pt9912/u-boot")
	mustBrew(t, "--yes", "generate", "devcontainer")
	df, err := os.ReadFile(filepath.Join(dir, ".devcontainer", "Dockerfile"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(df), "brew install --formula 'ripgrep' 'pt9912/u-boot/u-boot'") {
		t.Errorf("Dockerfile lacks the install line:\n%s", df)
	}
}

// The approval on the write path: with the key enabled, `config set` of the
// packages refuses a formula of an unapproved tap (exit code 10, the file as it
// was); a list written while the key is off is refused at the enabling; the
// tap list first lets both go through.
func TestRun_ConfigSet_BrewTapApprovalOnTheWritePath(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	dir := writeFeatureProject(t, brewProjectHead)
	mustBrew(t, "config", "set", "devcontainer.features.brew.enabled", "true")
	before, _ := os.ReadFile(filepath.Join(dir, "u-boot.yaml"))
	if code, out, errOut := runBrew(t, "--yes", "config", "set", "devcontainer.features.brew.packages", "ripgrep,pt9912/u-boot/u-boot"); code != 10 {
		t.Fatalf("packages with an unapproved tap: exit %d, want 10; stdout=%q stderr=%q", code, out, errOut)
	}
	if after, _ := os.ReadFile(filepath.Join(dir, "u-boot.yaml")); !bytes.Equal(before, after) {
		t.Fatalf("u-boot.yaml changed by the rejected call:\n%s", after)
	}
	mustBrew(t, "config", "set", "devcontainer.features.brew.enabled", "false")
	mustBrew(t, "config", "set", "devcontainer.features.brew.packages", "ripgrep,pt9912/u-boot/u-boot")
	if code, out, errOut := runBrew(t, "config", "set", "devcontainer.features.brew.enabled", "true"); code != 10 {
		t.Fatalf("enabling with an unapproved tap in the list: exit %d, want 10; stdout=%q stderr=%q", code, out, errOut)
	}
	mustBrew(t, "config", "set", "devcontainer.brewTaps.allow", "pt9912/u-boot")
	mustBrew(t, "config", "set", "devcontainer.features.brew.enabled", "true")
	mustBrew(t, "generate", "devcontainer")
}

// The notices: one for the block and one per approved tap, once each, in the
// default and the sandbox profile; `--quiet` drops them and keeps the block.
func TestRun_Generate_BrewNotices(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	for name, profile := range map[string]string{"default profile": "", "sandbox profile": "  profile: sandbox\n"} {
		t.Run(name, func(t *testing.T) {
			dir := writeFeatureProject(t, brewProjectHead+profile)
			mustBrew(t, "config", "set", "devcontainer.brewTaps.allow", "pt9912/u-boot")
			mustBrew(t, "config", "set", "devcontainer.features.brew.packages", "ripgrep,pt9912/u-boot/u-boot")
			mustBrew(t, "config", "set", "devcontainer.features.brew.enabled", "true")
			out := mustBrew(t, "generate", "devcontainer")
			for _, want := range []string{
				"Warning: devcontainer feature brew installs Homebrew from github.com/Homebrew/brew",
				"Warning: devcontainer feature brew installs formulas from the tap pt9912/u-boot",
			} {
				if n := strings.Count(out, want); n != 1 {
					t.Errorf("output contains %q %d times, want 1:\n%s", want, n, out)
				}
			}
			if strings.Contains(out, "bf464f57275cae0be106b43d190c98564ac59b0e") {
				t.Errorf("the notices name the pinned commit, which stands once in the Dockerfile:\n%s", out)
			}
			// A missing Dockerfile forces the write, so the notices are asked again under --quiet.
			if err := os.Remove(filepath.Join(dir, ".devcontainer", "Dockerfile")); err != nil {
				t.Fatal(err)
			}
			if quiet := strings.ToLower(mustBrew(t, "--quiet", "generate", "devcontainer")); strings.Contains(quiet, "feature brew") {
				t.Errorf("--quiet must drop the notices of the block:\n%s", quiet)
			}
			if _, err := os.Stat(filepath.Join(dir, ".devcontainer", "Dockerfile")); err != nil {
				t.Errorf("--quiet must still write the Dockerfile: %v", err)
			}
		})
	}
}

// The conflict with the OCI feature `homebrew` at the user path: `generate`
// ends with exit code 10 and writes nothing, `config set` refuses the enabling
// of either key next to the other, and either key switched off lets the
// project generate.
func TestRun_Brew_ConflictWithHomebrewFeatureExits10(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	writeFeatureProject(t, brewProjectHead+"  features:\n    brew:\n      enabled: true\n    homebrew:\n      enabled: true\n")
	code, out, errOut := runBrew(t, "generate", "devcontainer", "--dry-run", "--json")
	if code != 10 || !strings.Contains(out, "install two Homebrews") {
		t.Fatalf("generate --dry-run --json: exit %d, stdout=%q stderr=%q, want 10 naming the conflict", code, out, errOut)
	}
	if _, err := os.Stat(".devcontainer"); !os.IsNotExist(err) {
		t.Errorf("the rejected run must write nothing; .devcontainer exists (%v)", err)
	}
	mustBrew(t, "config", "set", "devcontainer.features.homebrew.enabled", "false")
	if code, out, errOut := runBrew(t, "config", "set", "devcontainer.features.homebrew.enabled", "true"); code != 10 {
		t.Fatalf("enabling homebrew next to the other key: exit %d, want 10; stdout=%q stderr=%q", code, out, errOut)
	}
	mustBrew(t, "generate", "devcontainer")
}

// A package list at another key and an entry outside the form: `config set`
// ends with exit code 10 and leaves `u-boot.yaml` as it was.
func TestRun_ConfigSet_BrewPackagesFormAndPlacementExit10(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	dir := writeFeatureProject(t, brewProjectHead)
	before, _ := os.ReadFile(filepath.Join(dir, "u-boot.yaml"))
	for _, args := range [][]string{
		{"config", "set", "devcontainer.features.brew.packages", "rip grep"},
		{"config", "set", "devcontainer.features.brew.packages", "ripgrep;id"},
		{"config", "set", "devcontainer.features.node.packages", "ripgrep"},
		{"config", "set", "devcontainer.features.homebrew.packages", "ripgrep"},
		{"config", "set", "devcontainer.brewTaps.allow", "pt9912"},
		{"config", "set", "devcontainer.brewTaps.allow", "a/b/c"},
	} {
		if code, out, errOut := runBrew(t, args...); code != 10 {
			t.Errorf("%v: exit %d, want 10; stdout=%q stderr=%q", args, code, out, errOut)
		}
	}
	after, _ := os.ReadFile(filepath.Join(dir, "u-boot.yaml"))
	if !bytes.Equal(before, after) {
		t.Errorf("u-boot.yaml changed by the rejected calls:\n%s", after)
	}
	mustBrew(t, "config", "set", "devcontainer.features.brew.packages", "ripgrep, jq")
	if got := strings.TrimSpace(mustBrew(t, "config", "get", "devcontainer.features.brew.packages")); !strings.HasSuffix(got, "ripgrep,jq") {
		t.Errorf("config get = %q, want ripgrep,jq", got)
	}
}
