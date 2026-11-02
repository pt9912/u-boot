//go:build docker

// End-to-end check of the catalogue keys of the pinned toolchain
// (`go-toolchain`, `golangci-lint`; LH-FA-DEV-003, lastenheft 0.30.0,
// SPEC-054): the generated Dockerfile builds against the real engine, the
// tools report the pinned versions, a falsified checksum breaks the build
// before the extraction, a `version:` override builds against the checksums
// read from the official source of the key, and the caches of
// `containerEnv` live under the workspace volume and survive the recreate.
// The build downloads from go.dev and github.com; a network failure turns
// the tests red without a change of the code.

package e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	confirmadapter "github.com/pt9912/u-boot/internal/adapter/driven/confirm"
	fetchadapter "github.com/pt9912/u-boot/internal/adapter/driven/fetch"
	fsadapter "github.com/pt9912/u-boot/internal/adapter/driven/fs"
	gitadapter "github.com/pt9912/u-boot/internal/adapter/driven/git"
	loggeradapter "github.com/pt9912/u-boot/internal/adapter/driven/logger"
	progressadapter "github.com/pt9912/u-boot/internal/adapter/driven/progress"
	yamladapter "github.com/pt9912/u-boot/internal/adapter/driven/yaml"
	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// The override versions of the Docker test: one release behind the pin of
// the catalogue template, so the read source is a different answer than the
// pin (SPEC-054). The values age with the releases; a red test here names
// the versions in its failure.
const (
	e2eGoOverrideVersion    = "1.26.10"
	e2eLintOverrideVersion  = "v2.13.2"
	e2eGoDefaultVersion     = "1.27.3"
	e2eGolangciDefaultVer   = "v2.14.0"
	goCanaryMarker          = "go-build-marker-toolchain"
	modCanaryMarker         = "go-mod-marker-toolchain"
	e2eChecksumsReadTimeout = 60 * time.Second
)

// toolchainImageScript runs as the image user: the pinned versions of both
// tools and the existence of the binaries under their install paths.
const toolchainImageScript = `echo "GO:$(go version)"
echo "LINT:$(golangci-lint version 2>&1 | head -1)"
echo "GOMOD:$(go env GOMODCACHE 2>/dev/null || echo unset)"
echo "BINS:$(test -x /usr/local/go/bin/go && test -x /usr/local/bin/golangci-lint && echo both)"`

// LH-FA-DEV-003 (SPEC-054, S2/S3): the Dockerfile of the default profile
// with both keys builds, `go` reports the pinned version, `golangci-lint`
// reports the pinned version, and both binaries sit at their install paths.
func TestE2E_Toolchain_BuildsAndReportsPins(t *testing.T) {
	requireDocker(t)
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := dcGenerateProject(ctx, t, false, map[string]string{
		"devcontainer.features.go-toolchain.enabled":  "true",
		"devcontainer.features.golangci-lint.enabled": "true",
	})
	image := buildProjectImage(ctx, t, dir, "uboot-e2e-toolchain")
	out, code := dockerOutput(ctx, t, "run", "--rm", image, "bash", "-c", toolchainImageScript)
	if code != 0 {
		t.Fatalf("image checks exit %d:\n%s", code, out)
	}
	for _, want := range []string{
		"GO:go version go" + e2eGoDefaultVersion + " linux/",
		// The version output of golangci-lint carries the release number
		// without the tag's `v` prefix ("has version 2.14.0 built with ...").
		"LINT:golangci-lint has version " + strings.TrimPrefix(e2eGolangciDefaultVer, "v"),
		"BINS:both",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("image checks lack %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "GO:go version go") && !strings.Contains(out, e2eGoDefaultVersion) {
		t.Errorf("go does not report the pin:\n%s", out)
	}
}

// The ARG prefix of the block pin per key (the `go-toolchain` key installs
// the `go` archive and carries its pin as GO_*); shared by the probes that
// address the pin via `--build-arg`.
func toolchainArgPin(key string) string {
	if key == "go-toolchain" {
		return "GO"
	}
	return "GOLANGCI_LINT"
}

// ADR-0041 fitness function (SPEC-054), negative probe per key: the build
// stops before the extraction when the pinned checksum does not match;
// `--build-arg` replaces the architecture value, so the probe holds for both
// archives. The failed build leaves no container with the zero pin
// (BEO-E2E/docker-test-hinterlaesst-zwischencontainer). The catalogue key
// and the ARG prefix of a key differ (the `go-toolchain` key installs the
// `go` archive); the probes follow the ARG prefix.
func TestE2E_Toolchain_PinMismatchBreaksBuild(t *testing.T) {
	requireDocker(t)
	for probe, key := range map[string]string{"go": "go-toolchain", "golangci-lint": "golangci-lint"} {
		arg := toolchainArgPin(key)
		t.Run(probe, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
			defer cancel()
			dir := dcGenerateProject(ctx, t, false, map[string]string{
				"devcontainer.features." + key + ".enabled": "true",
			})
			zeros := strings.Repeat("0", 64)
			tag := "uboot-e2e-toolchain-pin-negative-" + probe + ":probe"
			pin := arg + "_SHA256_AMD64=" + zeros
			before := containerIDs(ctx, t)
			// The classic builder keeps the container of the failed step
			// unless `--force-rm` is set.
			out, code := dockerOutput(ctx, t, "build", "-q", "--force-rm", "-t", tag,
				"--build-arg", pin, "--build-arg", arg+"_SHA256_ARM64="+zeros,
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

// dcGenerateToolchainProject mirrors dcGenerateProject with the checksum
// reader wired: the production adapters plus the production fetch adapter of
// the driven layer; the `version:` overrides read the live checksum sources
// (SPEC-054) and the Docker test reads their result.
func dcGenerateToolchainProject(ctx context.Context, t *testing.T, sandbox bool, settings map[string]string) string {
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
	for path, value := range settings {
		p, err := domain.NewConfigPath(path)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := cfgSvc.Set(ctx, driving.ConfigSetRequest{BaseDir: dir, Path: p, Value: value}); err != nil {
			t.Fatalf("config set %s: %v", path, err)
		}
	}
	genSvc := application.NewGenerateService(fsys, yaml, logger).WithChecksumFetcher(fetchadapter.New(e2eChecksumsReadTimeout))
	if _, err := genSvc.Generate(ctx, driving.GenerateRequest{BaseDir: dir, Artifact: domain.ArtifactDevcontainer}); err != nil {
		t.Fatalf("generate devcontainer: %v", err)
	}
	if err := os.Chmod(dir, 0o777); err != nil {
		t.Fatal(err)
	}
	return dir
}

// The override build (SPEC-054): the generation reads the two official
// sources for the older versions and the build checks the archives against
// the read checksums; the tools report the override versions.
func TestE2E_Toolchain_OverrideBuildsAgainstReadChecksums(t *testing.T) {
	requireDocker(t)
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := dcGenerateToolchainProject(ctx, t, false, map[string]string{
		"devcontainer.features.go-toolchain.enabled":  "true",
		"devcontainer.features.golangci-lint.enabled": "true",
		"devcontainer.features.go-toolchain.version":  e2eGoOverrideVersion,
		"devcontainer.features.golangci-lint.version": e2eLintOverrideVersion,
	})
	df, err := os.ReadFile(filepath.Join(dir, ".devcontainer", "Dockerfile"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(df), "ARG GO_VERSION="+e2eGoOverrideVersion) ||
		!strings.Contains(string(df), "ARG GOLANGCI_LINT_VERSION="+e2eLintOverrideVersion) {
		t.Errorf("the generated Dockerfile does not carry the override versions:\n%s", df)
	}
	image := buildProjectImage(ctx, t, dir, "uboot-e2e-toolchain-override")
	out, code := dockerOutput(ctx, t, "run", "--rm", image, "bash", "-c", toolchainImageScript)
	if code != 0 {
		t.Fatalf("image checks exit %d:\n%s", code, out)
	}
	for _, want := range []string{
		"GO:go version go" + e2eGoOverrideVersion + " linux/",
		"LINT:golangci-lint has version " + strings.TrimPrefix(e2eLintOverrideVersion, "v"),
	} {
		if !strings.Contains(out, want) {
			t.Errorf("image checks lack %q:\n%s", want, out)
		}
	}
}

// LH-FA-DEV-001 (lastenheft 0.30.0, SPEC-054): the caches of `containerEnv`
// lie under the workspace volume and survive the recreate of
// `u-boot devcontainer up --recreate` — the sandbox profile, the markers of
// the M3/M3Z reading (docs/messungen/slice-v1-katalog-gepinnte-toolchain-
// entwurf) and the ownership of the volume content (`vscode:vscode`). The
// environment of the recreated container carries the resolved values again
// (M1c: `.Config.Env`).
func TestE2E_Toolchain_CacheSurvivesRecreate(t *testing.T) {
	requireDocker(t)
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := dcGenerateProject(ctx, t, true, map[string]string{
		"devcontainer.features.go-toolchain.enabled":  "true",
		"devcontainer.features.golangci-lint.enabled": "true",
	})
	jsonBody, err := os.ReadFile(filepath.Join(dir, ".devcontainer", "devcontainer.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`"GOCACHE": "${containerWorkspaceFolder}/.cache/go-build"`,
		`"GOMODCACHE": "${containerWorkspaceFolder}/.cache/go-mod"`,
	} {
		if !strings.Contains(string(jsonBody), want) {
			t.Errorf("devcontainer.json lacks %q:\n%s", want, jsonBody)
		}
	}
	first, code := dcUpJSON(ctx, t, dir)
	if code != 0 {
		t.Fatalf("first up: exit %d\nlast diagnostics: %+v", code, first.Diagnostics)
	}
	firstID := first.Data.ContainerID
	if firstID == "" {
		t.Fatal("first up left no container id in the envelope")
	}
	// The container user writes one marker with padding into each cache; the
	// padding makes a lost content visible and not only a lost marker line.
	cacheWriteScript := `set -e
for pair in "` + goCanaryMarker + `:$GOCACHE" "` + modCanaryMarker + `:$GOMODCACHE"; do
  name="${pair%%:*}"; dir="${pair#*:}"
  test -n "$dir" || { echo "ENV:unset"; exit 3; }
  mkdir -p "$dir"
  printf '%s\n' "$name" > "$dir/marker"
  dd if=/dev/zero bs=1k count=64 2>/dev/null >> "$dir/marker"
  test "$(wc -c < "$dir/marker")" -gt 65536
done
echo "ENV:$(printenv GOCACHE)|$(printenv GOMODCACHE)"
stat -c '%U:%G' "$GOCACHE" "$GOMODCACHE"`
	out, code := dockerOutput(ctx, t, "exec", "-u", "vscode", firstID, "bash", "-c", cacheWriteScript)
	if code != 0 {
		t.Fatalf("cache writing in the first container: exit %d\n%s", code, out)
	}
	if !strings.Contains(out, "ENV:/workspaces/"+filepath.Base(dir)+"/.cache/go-build|/workspaces/"+filepath.Base(dir)+"/.cache/go-mod") {
		t.Errorf("the environment of the container resolves the caches elsewhere:\n%s", out)
	}
	second, code := dcUpJSONArgs(ctx, t, dir, "--recreate")
	if code != 0 {
		t.Fatalf("recreate up: exit %d\nlast diagnostics: %+v", code, second.Diagnostics)
	}
	if secondID := second.Data.ContainerID; secondID == "" || secondID == firstID {
		t.Fatalf("the recreate kept the container: first %.12s second %.12s", firstID, secondID)
	}
	cacheReadScript := `set -e
test -n "$GOCACHE" && test -n "$GOMODCACHE"
echo "ENV:$(printenv GOCACHE)|$(printenv GOMODCACHE)"
echo "MARKERS:$(head -1 "$GOCACHE/marker") $(head -1 "$GOMODCACHE/marker")"
echo "SIZES:$(wc -c < "$GOCACHE/marker") $(wc -c < "$GOMODCACHE/marker")"
stat -c 'OWNER:%U:%G' "$GOCACHE" "$GOMODCACHE"`
	out, code = dockerOutput(ctx, t, "exec", "-u", "vscode", second.Data.ContainerID, "bash", "-c", cacheReadScript)
	if code != 0 {
		t.Fatalf("cache reading in the recreated container: exit %d\n%s", code, out)
	}
	if !strings.Contains(out, goCanaryMarker) || !strings.Contains(out, modCanaryMarker) {
		t.Errorf("the recreated container lost a cache marker:\n%s", out)
	}
	for _, size := range strings.Split(strings.SplitN(out[strings.Index(out, "SIZES:")+6:], "\n", 2)[0], " ") {
		n := 0
		fmt.Sscanf(size, "%d", &n)
		if n <= 65536 {
			t.Errorf("the recreated container carries an incomplete cache content (%s <= the 64KiB pad): %s", size, out)
		}
	}
	if n := strings.Count(out, "OWNER:vscode:vscode"); n != 2 {
		t.Errorf("the cache directories carry the vscode ownership %d times, want twice:\n%s", n, out)
	}
}

// dcUpJSONArgs runs `u-boot devcontainer up --json` with further flags in
// dir and returns the decoded envelope and the exit code of the process.
func dcUpJSONArgs(ctx context.Context, t *testing.T, dir string, extra ...string) (dcEnvelope, int) {
	t.Helper()
	args := append([]string{"devcontainer", "up", "--json"}, extra...)
	stdout, stderr, code := dcRun(ctx, t, dir, nil, args...)
	var env dcEnvelope
	if err := json.Unmarshal([]byte(stdout), &env); err != nil {
		t.Fatalf("stdout is not one JSON envelope: %v\nstdout:\n%s\nstderr:\n%s", err, stdout, stderr)
	}
	return env, code
}

// The read of the `golangci-lint` risk item (SPEC-054: golangci-lint without
// go-toolchain and without an external Go feature): the block builds and the
// binary stands and reports its version. This is a reading, no zusage: the
// lint run of the tool without a Go installation is not part of SPEC-054.
func TestE2E_Toolchain_GolangciLintAloneBuilds(t *testing.T) {
	requireDocker(t)
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := dcGenerateProject(ctx, t, false, map[string]string{
		"devcontainer.features.golangci-lint.enabled": "true",
	})
	image := buildProjectImage(ctx, t, dir, "uboot-e2e-toolchain-lint-alone")
	out, code := dockerOutput(ctx, t, "run", "--rm", image, "bash", "-c",
		`golangci-lint version | head -1; echo "GO-BINARY:$(command -v go || echo missing)"`)
	if code != 0 {
		t.Fatalf("image checks exit %d:\n%s", code, out)
	}
	if !strings.Contains(out, "golangci-lint has version "+strings.TrimPrefix(e2eGolangciDefaultVer, "v")) {
		t.Errorf("golangci-lint does not report the pin:\n%s", out)
	}
	t.Logf("golangci-lint alone, the go binary of the image:\n%s", out)
}
