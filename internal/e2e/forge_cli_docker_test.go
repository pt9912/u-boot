//go:build docker

// End-to-end check of the catalogue features `gh` and `glab` (LH-FA-DEV-003,
// LH-FA-DEV-008, LH-FA-DEV-009): the generated Dockerfile builds in both
// profiles with the real Docker engine, the pinned archives are checked
// against their SHA-256, `glab` reaches `gitlab.com` behind the egress
// restriction only with the host in the list, and the token variables of the
// host reach the container without a value in any file, image or volume. The
// build downloads from github.com and gitlab.com; a network failure turns the
// tests red without a change of the code.

package e2e_test

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const (
	forgeCanaryGh   = "ghp_forgeCANARYcanary0123456789abcdef"
	forgeCanaryGlab = "glpat-forgeCANARYcanary0123456789"
)

// forgeImageScript runs as the image user without volume and without egress:
// the pinned versions, the owner of the binaries, the user, and no store of a
// token below the home directory. The check of the store runs first: both tools
// create their configuration directory on their first call.
const forgeImageScript = `test ! -e /home/vscode/.config/gh && test ! -e /home/vscode/.config/glab-cli && echo NO_TOKEN_STORE
gh --version | head -1
glab --version | head -1
echo "OWNER:$(stat -c %U /usr/local/bin/gh /usr/local/bin/glab | tr '\n' ' ')"
echo "USER:$(id -un)"`

// LH-FA-DEV-003 (features `gh` and `glab`), B1 and B3: the Dockerfile
// generated for the default profile and for the sandbox profile builds, the
// image holds the pinned versions, the binaries belong to root, the image user
// is vscode, and the image stores no token. The default profile is built here
// for the first time with a `USER root` that only the block causes.
func TestE2E_ForgeCLI_BuildsInBothProfiles(t *testing.T) {
	requireDocker(t)
	for name, sandbox := range map[string]bool{"default profile": false, "sandbox profile": true} {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 25*time.Minute)
			defer cancel()
			dir := dcGenerateProject(ctx, t, sandbox, map[string]string{
				"devcontainer.features.gh.enabled":   "true",
				"devcontainer.features.glab.enabled": "true",
			})
			image := buildProjectImage(ctx, t, dir, "uboot-e2e-forge")
			out, code := dockerOutput(ctx, t, "run", "--rm", image, "bash", "-c", forgeImageScript)
			if code != 0 {
				t.Fatalf("image checks exit %d:\n%s", code, out)
			}
			for _, want := range []string{"gh version 2.103.0 ", "glab 1.124.0 ", "OWNER:root root \n", "USER:vscode\n", "NO_TOKEN_STORE"} {
				if !strings.Contains(out, want) {
					t.Errorf("image checks lack %q:\n%s", want, out)
				}
			}
		})
	}
}

// ADR-0035 decision 1, negative probe per feature: the build stops before the
// extraction when the pinned checksum does not match; `--build-arg` replaces
// both architecture values, so the probe holds on `amd64` and `arm64` hosts.
// The failed build leaves no container with the zero pin.
func TestE2E_ForgeCLI_PinMismatchBreaksBuild(t *testing.T) {
	requireDocker(t)
	for _, upper := range []string{"GH", "GLAB"} {
		t.Run(upper, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 25*time.Minute)
			defer cancel()
			dir := dcGenerateProject(ctx, t, false, map[string]string{
				"devcontainer.features." + strings.ToLower(upper) + ".enabled": "true",
			})
			zeros := strings.Repeat("0", 64)
			tag := "uboot-e2e-forge-pin-negative-" + strings.ToLower(upper) + ":probe"
			pin := upper + "_SHA256_AMD64=" + zeros
			before := containerIDs(ctx, t)
			// The classic builder keeps the container of the failed step
			// unless `--force-rm` is set.
			out, code := dockerOutput(ctx, t, "build", "-q", "--force-rm", "-t", tag,
				"--build-arg", pin, "--build-arg", upper+"_SHA256_ARM64="+zeros,
				filepath.Join(dir, ".devcontainer"))
			t.Cleanup(func() {
				for _, id := range newContainersWithCommand(context.Background(), t, before, pin) {
					if rmOut, rmCode := dockerOutput(context.Background(), t, "rm", "-f", id); rmCode != 0 {
						t.Errorf("the failed build left container %.12s with the zero pin; removing it failed: %s", id, rmOut)
						continue
					}
					t.Errorf("the failed build left container %.12s with the zero pin; cleanup removed it", id)
				}
			})
			if code == 0 {
				t.Cleanup(func() { _, _ = dockerOutput(context.Background(), t, "rmi", "-f", tag) })
				t.Fatalf("build with a wrong checksum succeeded, want a failure before the extraction:\n%s", out)
			}
			if !strings.Contains(out, "FAILED") {
				t.Errorf("build output lacks the checksum failure:\n%s", out)
			}
		})
	}
}

// forgeEgressScript runs as `vscode` in the container behind the egress
// restriction: the list of the start script, the version of `gh`, a call of
// `glab` to gitlab.com with the host in the list, and the same call after the
// host was removed from the list.
const forgeEgressScript = `sudo sh /usr/local/bin/u-boot-egress-init
hosts="$(grep '^hosts=' /usr/local/bin/u-boot-egress-init | cut -d'"' -f2)"
echo "HOSTS:$hosts"
gh --version | head -1
if glab api projects/gitlab-org%2Fcli --silent > /tmp/allowed.log 2>&1; then echo GLAB_ALLOWED_OK; else echo GLAB_ALLOWED_FAILED; head -5 /tmp/allowed.log; fi
rest="$(echo "$hosts" | tr ' ' '\n' | grep -vx 'gitlab.com' | tr '\n' ' ')"
sudo sed -i "s|^hosts=.*|hosts=\"$rest\"|" /usr/local/bin/u-boot-egress-init
sudo sh /usr/local/bin/u-boot-egress-init
glab api projects/gitlab-org%2Fcli --silent > /tmp/blocked.log 2>&1
echo "GLAB_BLOCKED_RC:$?"
grep -F 'lookup gitlab.com' /tmp/blocked.log > /dev/null && echo LOOKUP_FAILED_MESSAGE`

// LH-FA-DEV-008 (hosts of `gh` and `glab`), B6: behind the egress restriction
// the start script lists the three hosts, `glab` reaches gitlab.com with the
// host in the list, and the same call fails on the name lookup once the host
// is out of the list. The call with a valid token is not part of the promise.
func TestE2E_ForgeCLI_EgressHosts(t *testing.T) {
	requireDocker(t)
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Minute)
	defer cancel()
	dir := dcGenerateProject(ctx, t, true, map[string]string{
		"devcontainer.features.gh.enabled":    "true",
		"devcontainer.features.glab.enabled":  "true",
		"devcontainer.sandbox.egress.enabled": "true",
		"devcontainer.sandbox.onUnavailable":  "warn",
	})
	image := buildProjectImage(ctx, t, dir, "uboot-e2e-forge-egress")
	args := append(append([]string{"run", "--rm"}, generatedRunArgs(t, dir)...), image, "bash", "-c", forgeEgressScript)
	out, code := dockerOutput(ctx, t, args...)
	if code != 0 {
		t.Fatalf("egress script exit %d:\n%s", code, out)
	}
	hostsLine := ""
	for _, line := range strings.Split(out, "\n") {
		if rest, ok := strings.CutPrefix(line, "HOSTS:"); ok {
			hostsLine = rest
		}
	}
	for _, host := range []string{"api.github.com", "github.com", "gitlab.com"} {
		if !strings.Contains(" "+hostsLine+" ", " "+host+" ") {
			t.Errorf("the allowlist of the start script lacks %s: %q", host, hostsLine)
		}
	}
	for _, want := range []string{"gh version 2.103.0 ", "GLAB_ALLOWED_OK", "GLAB_BLOCKED_RC:1", "LOOKUP_FAILED_MESSAGE"} {
		if !strings.Contains(out, want) {
			t.Errorf("egress script output lacks %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "GLAB_ALLOWED_FAILED") {
		t.Errorf("glab did not reach gitlab.com with the host in the list:\n%s", out)
	}
}

// streamContains reports for each needle whether it occurs in the stream, also
// across the boundary of two reads.
func streamContains(r io.Reader, needles ...string) (map[string]bool, error) {
	found := map[string]bool{}
	longest := 0
	for _, n := range needles {
		longest = max(longest, len(n))
	}
	buf := make([]byte, 1<<20)
	var tail []byte
	for {
		n, err := r.Read(buf)
		window := append(tail, buf[:n]...)
		for _, needle := range needles {
			if bytes.Contains(window, []byte(needle)) {
				found[needle] = true
			}
		}
		if keep := longest - 1; len(window) > keep {
			tail = append([]byte(nil), window[len(window)-keep:]...)
		} else {
			tail = append([]byte(nil), window...)
		}
		if err == io.EOF {
			return found, nil
		}
		if err != nil {
			return found, err
		}
	}
}

// Guard of the scan of the saved image: a needle that straddles the boundary of
// two reads is found, an absent needle is not.
func TestStreamContains_FindsNeedleAcrossReads(t *testing.T) {
	r := io.MultiReader(strings.NewReader(strings.Repeat("x", 1<<20-5)+"canary-"), strings.NewReader("value"+strings.Repeat("y", 100)))
	got, err := streamContains(r, "canary-value", "absent-value")
	if err != nil {
		t.Fatal(err)
	}
	if !got["canary-value"] || got["absent-value"] {
		t.Errorf("found = %v, want canary-value only", got)
	}
}

// LH-FA-DEV-009 (token variables of `gh` and `glab`), B9: with the sandbox
// profile the references in `remoteEnv` carry the values of the host into the
// session, and the values stand in no generated file, no saved image, no volume
// and not in the environment of the container. The canary values are invented.
// Guards: the scan of the volume finds a value written into it, and the check
// of the environment finds a value passed with `docker run -e`.
func TestE2E_ForgeCLI_TokenReferencesCarryNoValue(t *testing.T) {
	requireDocker(t)
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	t.Setenv("GH_TOKEN", forgeCanaryGh)
	t.Setenv("GITLAB_TOKEN", forgeCanaryGlab)
	t.Setenv("GIT_TOKEN", sandboxCanaryToken)
	dir := dcGenerateProject(ctx, t, true, map[string]string{
		"devcontainer.features.gh.enabled":   "true",
		"devcontainer.features.glab.enabled": "true",
	})
	canaries := []string{forgeCanaryGh, forgeCanaryGlab}

	_ = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			if b, rerr := os.ReadFile(p); rerr == nil {
				for _, c := range canaries {
					if bytes.Contains(b, []byte(c)) {
						t.Errorf("credential value found in generated file %s", p)
					}
				}
			}
		}
		return nil
	})

	bin := ubootBinary(t)
	stdout, stderr, code := runDCShell(ctx, t, bin, dir, "printenv GH_TOKEN\nprintenv GITLAB_TOKEN\nexit\n")
	if code != 0 {
		t.Fatalf("devcontainer shell exit %d\nstdout:\n%s\nstderr:\n%s", code, stdout, stderr)
	}
	if lines := strings.Split(strings.TrimSpace(stdout), "\n"); len(lines) != 2 || lines[0] != forgeCanaryGh || lines[1] != forgeCanaryGlab {
		t.Errorf("the session sees the variables %q, want the two values of the host (the references work)", stdout)
	}
	for _, c := range canaries {
		if strings.Contains(stderr, c) {
			t.Errorf("standard error of the command carries a value:\n%.1500s", stderr)
		}
	}

	id := dcRunningContainer(t, dir)
	doc := dcInspect(t, id)
	for _, c := range canaries {
		if envHolds(doc, c) {
			t.Errorf("Config.Env of the devcontainer holds %q", c)
		}
	}
	guard := "ubl-forgeguard-" + dcSuffix()
	t.Cleanup(func() { _, _ = dcDocker("rm", "-f", guard) })
	dcMustDocker(t, "run", "-d", "--name", guard, "-e", "UBL_GUARD="+forgeCanaryGh, dcBaseImage, "sleep", "60")
	if !envHolds(dcInspect(t, guard), forgeCanaryGh) {
		t.Errorf("guard: the check of Config.Env does not find a value passed with `docker run -e`")
	}

	cmd := exec.CommandContext(ctx, "docker", "save", doc.Config.Image)
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatalf("docker save: %v", err)
	}
	inImage, scanErr := streamContains(pipe, canaries...)
	if err := cmd.Wait(); err != nil || scanErr != nil {
		t.Fatalf("docker save %s: wait %v, scan %v", doc.Config.Image, err, scanErr)
	}
	for _, c := range canaries {
		if inImage[c] {
			t.Errorf("the saved image holds %q", c)
		}
	}

	workspace := "/workspaces/dcstart"
	var volumes []string
	for _, m := range doc.Mounts {
		if m.Type == "volume" {
			volumes = append(volumes, m.Name)
		}
	}
	for _, v := range volumes {
		for _, c := range canaries {
			if dcVolumeHolds(ctx, t, doc.Config.Image, v, c) {
				t.Errorf("volume %s holds %q", v, c)
			}
		}
	}
	volume := doc.volumeAt(workspace)
	if volume == "" {
		t.Fatalf("no named volume at %s: %+v", workspace, doc.Mounts)
	}
	if out, rc := dockerOutput(ctx, t, "exec", "-u", "root", id, "sh", "-c", `printf '%s\n' "$1" > "$2"`, "sh", forgeCanaryGh, workspace+"/leak.txt"); rc != 0 {
		t.Fatalf("write the guard value: %s", out)
	}
	if !dcVolumeHolds(ctx, t, doc.Config.Image, volume, forgeCanaryGh) {
		t.Errorf("guard: the scan does not find a value written into the workspace volume")
	}
}
