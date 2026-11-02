package application_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// Golden cases of the catalogue keys `go-toolchain` and `golangci-lint`
// (LH-FA-DEV-003, lastenheft 0.30.0, SPEC-054): a Dockerfile block with a
// pinned release archive in both profiles, the cache values of
// `containerEnv` in the generated devcontainer.json, no OCI render key, one
// supply-chain notice per active key, the `version:` override with the read
// checksums of the official source and the exit-10 rejections.

// The checksum values of the tests are the recorded values of the two
// sources (go.dev/dl listing for 1.27.3 and 1.26.10, the checksums file of
// the golangci-lint releases v2.14.0 and v2.13.2); the Docker tests of the
// package (`internal/e2e`) read the live sources, these fixtures pin the
// parsing offline.
const (
	goDefaultVersion   = "1.27.3"
	goSHAAMD64Default  = "65efbc80fdc1fa2cd2631cad5b481777a68e8725f34e89f50fc6ac714e379b40"
	goSHAARM64Default  = "a06e99309dbcbdca1a8ccd2282c3ee146e55930886e26b4f1193269587165f49"
	golangciDefaultVer = "v2.14.0"
	golangciSHAAMD64   = "ab90aeb7b066f92a33415b638a50fe5344bbb75a0d32ad30cc248d88f81032ab"
	golangciSHAARM64   = "ee7ec5f3453d15ddf106fae5a4d6c71737712348a979d1fe9cd52ec7ea299bae"

	goOverrideVersion = "1.26.10"
	goOverrideAMD64   = "ccfe7ae822d2c2f41f4fe7c956199967a4dac22d0ee0d2b2bfdb440aab3b20c5"
	goOverrideARM64   = "3f4227ce7b206a513e794f1393adb473a07b46e94d18bd15d1c684bbb7319797"
	lintOverrideTag   = "v2.13.2"
	lintOverrideAMD64 = "2277d43b98ec0054280f2ac26b53268bae97682444678a59a657dd565da021d6"
	lintOverrideARM64 = "a2a4e0065aa41be71f7c5ac90f271b61751331e5d04314e62afe4027855f0893"

	// The java-toolchain key (lastenheft 0.31.0): recorded values of the
	// Temurin checksums files of jdk-21.0.13+9 (the default pin) and
	// jdk-21.0.9+10 (the override fixtures); the Docker tests read the live
	// sources, these constants pin the parsing offline.
	jdkDefaultVer  = "21.0.13+9"
	jdkDefaultX64  = "d7f5d5cf865a2902b70ff549b97586a7d3f234e016f041ab447ea740c7d66d1c"
	jdkDefaultA64  = "e3522e707fc7dbe2a72c30746b1e7068301a0213ecd213d4491f35a14eb55578"
	jdkOverrideVer = "21.0.9+10"
	jdkOverrideX64 = "810d3773df7e0d6c4394e4e244b264c8b30e0b05a0acf542d065fd78a6b65c2f"
	jdkOverrideA64 = "edf0da4debe7cf475dbe320d174d6eed81479eb363f41e38a2efb740428c603a"
)

// toolchainSpec describes the Dockerfile block of one pinned-toolchain key.
type toolchainSpec struct {
	key, arg, version, shaAMD64, shaARM64, downloadURL string
	// caseLine is the architecture case line of the block; empty means the
	// plain `sha256="${ARG_SHA256_*}"` template (the two Go keys; the JDK
	// block resolves the Temurin slug in the same line).
	caseLine string
}

func toolchainSpecs() map[string]toolchainSpec {
	return map[string]toolchainSpec{
		"go-toolchain": {"go-toolchain", "GO", goDefaultVersion, goSHAAMD64Default, goSHAARM64Default,
			"https://go.dev/dl/go${GO_VERSION}.linux-${arch}.tar.gz", ""},
		"golangci-lint": {"golangci-lint", "GOLANGCI_LINT", golangciDefaultVer, golangciSHAAMD64, golangciSHAARM64,
			"https://github.com/golangci/golangci-lint/releases/download/${GOLANGCI_LINT_VERSION}/golangci-lint-${GOLANGCI_LINT_VERSION#v}-linux-${arch}.tar.gz", ""},
		"java-toolchain": {"java-toolchain", "JDK", jdkDefaultVer, jdkDefaultX64, jdkDefaultA64,
			"https://github.com/adoptium/temurin21-binaries/releases/download/jdk-$(printf '%s' \"$JDK_VERSION\" | sed 's/+/%2B/')/${asset}.tar.gz",
			`case "$arch" in amd64) slug="x64" sha256="${JDK_SHA256_AMD64}" ;; arm64) slug="aarch64" sha256="${JDK_SHA256_ARM64}" ;; *) echo "the JDK has no archive for $arch" >&2; exit 1 ;; esac`},
	}
}

// toolName is the archive naming of the block error text: the upstream tool
// name, not the catalogue key (the go-toolchain block installs the `go`
// archive; the other keys carry their key as the name).
func toolName(key string) string {
	if key == "go-toolchain" {
		return "go"
	}
	if key == "java-toolchain" {
		return "the JDK"
	}
	return key
}

// toolchainJSONGolden is the devcontainer.json of the default profile for a
// project with `go-toolchain` enabled: no key under `features` (the block
// carries the installation), the two cache values under `containerEnv`.
const toolchainJSONGolden = "// BEGIN U-BOOT MANAGED BLOCK: init\n" +
	"{\n" +
	"  \"name\": \"t-uboot-gen\",\n" +
	"  \"build\": {\n" +
	"    \"dockerfile\": \"./Dockerfile\",\n" +
	"    \"context\": \".\"\n" +
	"  },\n" +
	"  \"containerEnv\": {\n" +
	"    \"GOCACHE\": \"${containerWorkspaceFolder}/.cache/go-build\",\n" +
	"    \"GOMODCACHE\": \"${containerWorkspaceFolder}/.cache/go-mod\"\n" +
	"  },\n" +
	"  \"remoteUser\": \"vscode\"\n" +
	"}\n" +
	"// END U-BOOT MANAGED BLOCK: init\n"

// javaJSONGolden is the devcontainer.json of the default profile for a
// project with `java-toolchain` enabled: no key under `features`, and NO
// `containerEnv` values — the JDK key carries JAVA_HOME and PATH as ENV
// lines of its Dockerfile block (SPEC-054, Cache-Werte-Absatz).
const javaJSONGolden = "// BEGIN U-BOOT MANAGED BLOCK: init\n" +
	"{\n" +
	"  \"name\": \"t-uboot-gen\",\n" +
	"  \"build\": {\n" +
	"    \"dockerfile\": \"./Dockerfile\",\n" +
	"    \"context\": \".\"\n" +
	"  },\n" +
	"  \"remoteUser\": \"vscode\"\n" +
	"}\n" +
	"// END U-BOOT MANAGED BLOCK: init\n"

// seedToolchainYAML writes a u-boot.yaml fixture under the generate base
// directory (the seeder of the add tests writes under its own base dir).
func seedToolchainYAML(t *testing.T, fs *fakeFS, body string) {
	t.Helper()
	if err := fs.WriteFile(filepath.Join(generateTestBaseDir, "u-boot.yaml"),
		[]byte(body), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
}

// toolchainFixture is one u-boot.yaml fixture of the test: the given
// `features:` entries behind a minimal devcontainer of the default profile.
func toolchainFixture(featureEntries ...string) string {
	body := "schemaVersion: 1\nproject:\n  name: t-uboot-gen\ndevcontainer:\n  enabled: true\n"
	if len(featureEntries) > 0 {
		body += "  features:\n" + strings.Join(featureEntries, "")
	}
	return body
}

// toolchainEntries builds the `features:` lines of the named keys, sorted so
// the fixture is deterministic.
func toolchainEntries(keys ...string) []string {
	enabled := make([]string, 0, len(keys))
	for _, k := range keys {
		enabled = append(enabled, "    "+k+":\n      enabled: true\n")
	}
	return enabled
}

const goToolchainEntry = "    go-toolchain:\n      enabled: true\n"

const golangciEntry = "    golangci-lint:\n      enabled: true\n"

const javaEntry = "    java-toolchain:\n      enabled: true\n"

const javaFeatureEntry = "    java:\n      enabled: true\n"

const goFeatureEntry = "    go:\n      enabled: true\n"

// goListingFixture is the shape of the go.dev/dl listing for the read of the
// named versions; two checksum values per version (the records of the two
// real answers) and no further files — the production listing is larger and
// carries no rows of concern.
func goListingFixture(version string) string {
	amd64, arm64 := goOverrideAMD64, goOverrideARM64
	if version == goDefaultVersion {
		amd64, arm64 = goSHAAMD64Default, goSHAARM64Default
	}
	vars := []struct{ arch, sha string }{{"amd64", amd64}, {"arm64", arm64}}
	var s strings.Builder
	s.WriteString("{\n\"version\": \"go" + version + "\",\n\"files\": [\n")
	for i, a := range vars {
		if i > 0 {
			s.WriteString(",\n")
		}
		s.WriteString(`{"filename": "go` + version + `.linux-` + a.arch + `.tar.gz", "os": "linux", "arch": "` + a.arch + `", "sha256": "` + a.sha + `", "kind": "archive"}`)
	}
	s.WriteString("\n]\n}")
	return s.String()
}

// toolchainListing builds a complete listing fixture of the named versions
// with the comma-joined release entries.
func toolchainListing(versions ...string) []byte {
	var s strings.Builder
	s.WriteString("[\n")
	for i, v := range versions {
		if i > 0 {
			s.WriteString(",\n")
		}
		s.WriteString(goListingFixture(v))
	}
	s.WriteString("\n]\n")
	return []byte(s.String())
}

// checksumsFileFixture is the shape of a release checksums file: the pairs
// of the two archives beside rows the read must ignore.
func checksumsFileFixture(tag, amd64, arm64 string) []byte {
	base := strings.TrimPrefix(tag, "v")
	var b strings.Builder
	b.WriteString(strings.Repeat("0", 64) + "  golangci-lint-" + base + "-linux-amd64.deb\n")
	b.WriteString(strings.Repeat("1", 64) + "  golangci-lint-" + base + "-linux-amd64.rpm\n")
	b.WriteString(amd64 + "  golangci-lint-" + base + "-linux-amd64.tar.gz\n")
	b.WriteString(strings.Repeat("2", 64) + "  golangci-lint-" + base + "-linux-amd64.tar.gz.sbom.json\n")
	b.WriteString(arm64 + "  golangci-lint-" + base + "-linux-arm64.tar.gz\n")
	return []byte(b.String())
}

// assertToolchainBlock pins the Dockerfile block of one pinned-toolchain key
// (SPEC-054): version and both checksums stand in one place, the check runs
// between the download and the extraction, and the block stands behind the
// one `USER root` of its file.
func assertToolchainBlock(t *testing.T, df string, s toolchainSpec) {
	t.Helper()
	caseLine := s.caseLine
	if caseLine == "" {
		caseLine = `case "$arch" in amd64) sha256="${` + s.arg + `_SHA256_AMD64}" ;; arm64) sha256="${` + s.arg + `_SHA256_ARM64}" ;; *) echo "` + toolName(s.key) + ` has no archive for $arch" >&2; exit 1 ;; esac`
	}
	for _, want := range []string{
		"ARG " + s.arg + "_VERSION=" + s.version + "\n",
		"ARG " + s.arg + "_SHA256_AMD64=" + s.shaAMD64 + "\n",
		"ARG " + s.arg + "_SHA256_ARM64=" + s.shaARM64 + "\n",
		s.downloadURL,
		caseLine,
	} {
		if !strings.Contains(df, want) {
			t.Errorf("Dockerfile lacks %q:\n%s", want, df)
		}
	}
	for _, once := range []string{s.shaAMD64, s.shaARM64, s.version} {
		if n := strings.Count(df, once); n != 1 {
			t.Errorf("Dockerfile carries the pin %q %d times, want exactly once (one place)", once, n)
		}
	}
	start := strings.Index(df, "ARG "+s.arg+"_VERSION=")
	end := strings.Index(df[start+1:], `rm -rf "$work"`)
	if start < 0 || end < 0 {
		t.Fatalf("block of %s not found:\n%s", s.key, df)
	}
	block := df[start : start+1+end]
	last := -1
	for _, step := range []string{"curl -fsSLO", "sha256sum -c -", "tar -xzf"} {
		i := strings.Index(block, step)
		if i < 0 || i < last {
			t.Errorf("block step %q of %s missing or out of order (index %d, previous %d)", step, s.key, i, last)
		}
		last = i
	}
}

// assertToolchainJSON pins the cache values of the generated
// devcontainer.json (LH-FA-DEV-001, SPEC-054): both values stand under
// `containerEnv`, resolved by the container workspace variable.
func assertToolchainJSON(t *testing.T, body []byte) {
	t.Helper()
	env := toolchainJSONMap(t, body)
	want := map[string]any{
		"GOCACHE":    "${containerWorkspaceFolder}/.cache/go-build",
		"GOMODCACHE": "${containerWorkspaceFolder}/.cache/go-mod",
	}
	if !reflect.DeepEqual(env, want) {
		t.Errorf("containerEnv = %v, want %v", env, want)
	}
}

// assertNoContainerEnv pins the Cache-Werte-Abgrenzung of the java-toolchain
// key (SPEC-054): the generated devcontainer.json carries no containerEnv
// value — JAVA_HOME and PATH are ENV lines of the Docker block.
func assertNoContainerEnv(t *testing.T, body []byte) {
	t.Helper()
	m := strippedJSONMap(t, body)
	if env, ok := m["containerEnv"]; ok {
		t.Errorf("devcontainer.json carries containerEnv although no Go key is active: %v\n%s", env, string(body))
	}
}

// toolchainJSONMap decodes the containerEnv map of a generated
// devcontainer.json.
func toolchainJSONMap(t *testing.T, body []byte) map[string]any {
	t.Helper()
	m := strippedJSONMap(t, body)
	env, ok := m["containerEnv"].(map[string]any)
	if !ok {
		t.Fatalf("devcontainer.json carries no containerEnv:\n%s", string(body))
	}
	return env
}

func strippedJSONMap(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(application.StripJSONCForTest(body), &m); err != nil {
		t.Fatalf("devcontainer.json invalid: %v\n%s", err, string(body))
	}
	return m
}

// newToolchainService builds the generate service with the fake fetcher; a
// nil fetcher keeps the noop default (its errors are the unreadable-source
// rejections).
func newToolchainService(t *testing.T, fetch *fakeFetcher) (*application.GenerateService, *fakeFS) {
	t.Helper()
	svc, fs := newGenerateService(t)
	if fetch != nil {
		svc.WithChecksumFetcher(fetch)
	}
	return svc, fs
}

// fresh state (SPEC-054, S1/S2): in the default profile, the sandbox profile
// and with the `--sandbox` flag the enabled key gets its pinned block and
// the cache values; the other key leaves no trace, no key stands under
// `features`, and the second run is a NoOp.
func TestGenerateDevcontainer_Toolchain_BlocksInBothProfiles(t *testing.T) {
	t.Parallel()
	type profile struct {
		head    string
		sandbox bool
	}
	profiles := map[string]profile{
		"default": {toolchainFixture(), false},
		"sandbox": {"schemaVersion: 1\nproject:\n  name: t-uboot-gen\ndevcontainer:\n  enabled: true\n  profile: sandbox\n", false},
		"flag":    {toolchainFixture(), true},
	}
	for pname, p := range profiles {
		for fname, keys := range map[string][]string{
			"go-toolchain":  {"go-toolchain"},
			"golangci-lint": {"golangci-lint"},
			"both":          {"go-toolchain", "golangci-lint"},
		} {
			t.Run(pname+"/"+fname, func(t *testing.T) {
				t.Parallel()
				fetch := newFakeFetcher()
				fetch.setBody(application.GoChecksumsSourceURLForTest(), toolchainListing(goDefaultVersion))
				fetch.setBody(application.GolangciChecksumsURLForTest(golangciDefaultVer),
					checksumsFileFixture(golangciDefaultVer, golangciSHAAMD64, golangciSHAARM64))
				svc, fs := newToolchainService(t, fetch)
				seedToolchainYAML(t, fs, p.head+"  features:\n"+strings.Join(toolchainEntries(keys...), ""))
				if _, err := generateSandbox(svc, p.sandbox); err != nil {
					t.Fatalf("generate: %v", err)
				}
				raw, _ := fs.ReadFile(dockerfilePath())
				df := string(raw)
				specs := toolchainSpecs()
				for _, k := range keys {
					assertToolchainBlock(t, df, specs[k])
				}
				for k, s := range specs {
					if !contains(keys, k) && strings.Contains(df, s.arg+"_VERSION") {
						t.Errorf("Dockerfile carries the block of %s although it is not enabled:\n%s", k, df)
					}
				}
				body, _ := fs.ReadFile(devcontainerJSONPath())
				assertToolchainJSON(t, body)
				if _, ok := strippedJSONMap(t, body)["features"]; ok {
					t.Errorf("devcontainer.json renders no key under features for the pinned toolchain:\n%s", body)
				}
				if n := strings.Count(df, "USER root"); n != 1 {
					t.Errorf("USER root appears %d times, want 1:\n%s", n, df)
				}
				if r2, err := generateSandbox(svc, p.sandbox); err != nil || r2.Action != driving.GenerateActionNoOp {
					t.Errorf("second run: action=%v err=%v, want NoOp", r2.Action, err)
				}
			})
		}
	}
}

// The default profile with `go-toolchain` alone is the byte golden; without
// the keys neither file carries a toolchain trace.
func TestGenerateDevcontainer_Toolchain_JSONGolden(t *testing.T) {
	t.Parallel()
	t.Run("golden with go-toolchain", func(t *testing.T) {
		t.Parallel()
		fetch := newFakeFetcher()
		svc, fs := newToolchainService(t, fetch)
		seedToolchainYAML(t, fs, toolchainFixture(goToolchainEntry))
		if _, err := generateSandbox(svc, false); err != nil {
			t.Fatalf("generate: %v", err)
		}
		raw, _ := fs.ReadFile(devcontainerJSONPath())
		if got := string(raw); got != toolchainJSONGolden {
			t.Errorf("devcontainer.json is not the golden:\n--- got ---\n%s\n--- want ---\n%s", got, toolchainJSONGolden)
		}
	})
	t.Run("without the keys nothing changes", func(t *testing.T) {
		t.Parallel()
		svc, fs := newToolchainService(t, newFakeFetcher())
		seedToolchainYAML(t, fs, toolchainFixture())
		if _, err := generateSandbox(svc, false); err != nil {
			t.Fatalf("generate: %v", err)
		}
		without, _ := fs.ReadFile(devcontainerJSONPath())
		if strings.Contains(string(without), "containerEnv") {
			t.Errorf("devcontainer.json without the keys must stay the old shape:\n%s", without)
		}
		df, _ := fs.ReadFile(dockerfilePath())
		for _, arg := range []string{"ARG GO_VERSION", "ARG GOLANGCI_LINT_VERSION"} {
			if strings.Contains(string(df), arg) {
				t.Errorf("Dockerfile without the keys must carry no %s:\n%s", arg, df)
			}
		}
	})
	t.Run("golden with java-toolchain carries no containerEnv values", func(t *testing.T) {
		t.Parallel()
		svc, fs := newToolchainService(t, newFakeFetcher())
		seedToolchainYAML(t, fs, toolchainFixture(javaEntry))
		if _, err := generateSandbox(svc, false); err != nil {
			t.Fatalf("generate: %v", err)
		}
		raw, _ := fs.ReadFile(devcontainerJSONPath())
		if got := string(raw); got != javaJSONGolden {
			t.Errorf("devcontainer.json is not the golden:\n--- got ---\n%s\n--- want ---\n%s", got, javaJSONGolden)
		}
	})
	t.Run("both keys carry the values once", func(t *testing.T) {
		t.Parallel()
		fetch := newFakeFetcher()
		fetch.setBody(application.GoChecksumsSourceURLForTest(), toolchainListing(goDefaultVersion))
		fetch.setBody(application.GolangciChecksumsURLForTest(golangciDefaultVer),
			checksumsFileFixture(golangciDefaultVer, golangciSHAAMD64, golangciSHAARM64))
		svc, fs := newToolchainService(t, fetch)
		seedToolchainYAML(t, fs, toolchainFixture(goToolchainEntry, golangciEntry))
		if _, err := generateSandbox(svc, true); err != nil {
			t.Fatalf("generate: %v", err)
		}
		body, _ := fs.ReadFile(devcontainerJSONPath())
		for _, name := range []string{"GOCACHE", "GOMODCACHE"} {
			if n := strings.Count(string(body), name); n != 1 {
				t.Errorf("%s appears %d times, want once with both keys:\n%s", name, n, body)
			}
		}
		assertToolchainJSON(t, body)
	})
}

// newToolchainOverrideFetcher wires the fixture read for the older versions
// (the go listing with the two released versions, the golangci checksums of
// the override tag, the two Temurin checksums files of the override release).
func newToolchainOverrideFetcher() *fakeFetcher {
	fetch := newFakeFetcher()
	fetch.setBody(application.GoChecksumsSourceURLForTest(), toolchainListing(goDefaultVersion, goOverrideVersion))
	fetch.setBody(application.GolangciChecksumsURLForTest(lintOverrideTag),
		checksumsFileFixture(lintOverrideTag, lintOverrideAMD64, lintOverrideARM64))
	fetch.setBody(application.JavaChecksumsURLForTest(jdkOverrideVer, "x64"),
		jdkChecksumsFileFixture(jdkOverrideX64, "OpenJDK21U-jdk_x64_linux_hotspot_21.0.9_10.tar.gz"))
	fetch.setBody(application.JavaChecksumsURLForTest(jdkOverrideVer, "aarch64"),
		jdkChecksumsFileFixture(jdkOverrideA64, "OpenJDK21U-jdk_aarch64_linux_hotspot_21.0.9_10.tar.gz"))
	return fetch
}

// jdkChecksumsFileFixture builds one Temurin checksums file: the single
// `<sha256>  <archive>` line form (one file per architecture).
func jdkChecksumsFileFixture(checksum, archive string) []byte {
	return []byte(checksum + "  " + archive + "\n")
}

// The override (SPEC-054, S4/S5): the generated Dockerfile carries the read
// checksums and the override version, the containerEnv block is unchanged,
// the double run re-reads and stays a NoOp, and the three static pin places
// of the repo carry the default pin still.
func TestGenerateDevcontainer_Toolchain_OverrideRendersReadChecksums(t *testing.T) {
	t.Parallel()
	svc, fs := newToolchainService(t, newToolchainOverrideFetcher())
	seedToolchainYAML(t, fs, toolchainFixture(
		"    go-toolchain:\n      enabled: true\n      version: "+goOverrideVersion+"\n",
		"    golangci-lint:\n      enabled: true\n      version: "+lintOverrideTag+"\n",
		"    java-toolchain:\n      enabled: true\n      version: "+jdkOverrideVer+"\n"))
	if _, err := generateSandbox(svc, false); err != nil {
		t.Fatalf("generate: %v", err)
	}
	raw, _ := fs.ReadFile(dockerfilePath())
	df := string(raw)
	assertToolchainBlock(t, df, toolchainSpec{
		key: "go-toolchain", arg: "GO", version: goOverrideVersion, shaAMD64: goOverrideAMD64, shaARM64: goOverrideARM64,
		downloadURL: "https://go.dev/dl/go${GO_VERSION}.linux-${arch}.tar.gz"})
	assertToolchainBlock(t, df, toolchainSpec{
		key: "golangci-lint", arg: "GOLANGCI_LINT", version: lintOverrideTag, shaAMD64: lintOverrideAMD64, shaARM64: lintOverrideARM64,
		downloadURL: "https://github.com/golangci/golangci-lint/releases/download/${GOLANGCI_LINT_VERSION}/golangci-lint-${GOLANGCI_LINT_VERSION#v}-linux-${arch}.tar.gz"})
	assertToolchainBlock(t, df, toolchainSpec{
		key: "java-toolchain", arg: "JDK", version: jdkOverrideVer, shaAMD64: jdkOverrideX64, shaARM64: jdkOverrideA64,
		downloadURL: "https://github.com/adoptium/temurin21-binaries/releases/download/jdk-$(printf '%s' \"$JDK_VERSION\" | sed 's/+/%2B/')/${asset}.tar.gz",
		caseLine:    `case "$arch" in amd64) slug="x64" sha256="${JDK_SHA256_AMD64}" ;; arm64) slug="aarch64" sha256="${JDK_SHA256_ARM64}" ;; *) echo "the JDK has no archive for $arch" >&2; exit 1 ;; esac`})
	assertToolchainJSON(t, mustRead(t, fs, devcontainerJSONPath()))
	// The second run re-reads the source and stays a NoOp: the same read
	// yields the same bytes.
	if r2, err := generateSandbox(svc, false); err != nil || r2.Action != driving.GenerateActionNoOp {
		t.Errorf("second run with override: action=%v err=%v, want NoOp", r2.Action, err)
	}
}

// The not-covered rejection (SPEC-054, S5): a valid listing that does not
// name the chosen version ends the run with the exit-10 class and nothing
// written.
func TestGenerateDevcontainer_Toolchain_VersionNotCoveredExit10(t *testing.T) {
	t.Parallel()
	fetch := newFakeFetcher()
	fetch.setBody(application.GoChecksumsSourceURLForTest(), toolchainListing(goDefaultVersion))
	svc, fs := newToolchainService(t, fetch)
	seedToolchainYAML(t, fs, toolchainFixture("    go-toolchain:\n      enabled: true\n      version: 9.9.9\n"))
	_, err := generateSandbox(svc, false)
	if !errors.Is(err, driving.ErrConfigValueInvalid) {
		t.Fatalf("generate: err = %v, want a wrap of ErrConfigValueInvalid (exit 10)", err)
	}
	if !strings.Contains(err.Error(), "does not cover") || !strings.Contains(err.Error(), "9.9.9") {
		t.Errorf("err = %v, want the not-covered wording and the version", err)
	}
	assertNothingWritten(t, fs)

	// The java-toolchain key reads one file per architecture; a file without
	// the archive row is the same not-covered rejection.
	fetch2 := newFakeFetcher()
	fetch2.setBody(application.JavaChecksumsURLForTest("21.0.99+99", "x64"),
		jdkChecksumsFileFixture(jdkDefaultX64, "OpenJDK21U-jdk_x64_linux_hotspot_21.0.13_9.tar.gz"))
	fetch2.setBody(application.JavaChecksumsURLForTest("21.0.99+99", "aarch64"),
		jdkChecksumsFileFixture(jdkDefaultA64, "OpenJDK21U-jdk_aarch64_linux_hotspot_21.0.13_9.tar.gz"))
	svc2, fs2 := newToolchainService(t, fetch2)
	seedToolchainYAML(t, fs2, toolchainFixture("    java-toolchain:\n      enabled: true\n      version: 21.0.99+99\n"))
	_, err = generateSandbox(svc2, false)
	if !errors.Is(err, driving.ErrConfigValueInvalid) {
		t.Fatalf("generate (java-toolchain not covered): err = %v, want a wrap of ErrConfigValueInvalid (exit 10)", err)
	}
	if !strings.Contains(err.Error(), "does not cover") || !strings.Contains(err.Error(), "OpenJDK21U-jdk_x64") {
		t.Errorf("err = %v, want the not-covered wording and the archive name", err)
	}
	assertNothingWritten(t, fs2)
}

// The unreadable-source rejection (SPEC-054, S5): a failing read of the
// configured source ends the run with the exit-10 class and nothing written;
// the message names the source and the URL. Without a wired reader an
// override of either key is the same rejection (the noop default).
func TestGenerateDevcontainer_Toolchain_UnreadableChecksumSourceExit10(t *testing.T) {
	t.Parallel()
	fetch := newFakeFetcher()
	fetch.setErr(application.GoChecksumsSourceURLForTest(), errors.New("connection refused"))
	svc, fs := newToolchainService(t, fetch)
	seedToolchainYAML(t, fs, toolchainFixture("    go-toolchain:\n      enabled: true\n      version: 1.26.10\n"))
	_, err := generateSandbox(svc, false)
	if !errors.Is(err, driving.ErrConfigValueInvalid) {
		t.Fatalf("generate: err = %v, want a wrap of ErrConfigValueInvalid (exit 10)", err)
	}
	if !strings.Contains(err.Error(), "is not readable") {
		t.Errorf("err = %v, want the unreadable wording", err)
	}
	assertNothingWritten(t, fs)

	// Without a wired reader an override of either key is the same
	// rejection (the noop default of the service).
	svc2, fs2 := newGenerateService(t)
	seedToolchainYAML(t, fs2, toolchainFixture("    golangci-lint:\n      enabled: true\n      version: v2.15.0\n"))
	if _, err := generateSandbox(svc2, false); !errors.Is(err, driving.ErrConfigValueInvalid) {
		t.Errorf("generate without a wired reader: err = %v, want the unreadable-source class (exit 10)", err)
	}

	// The java-toolchain key reads the checksums files of one release; a
	// failing read of the configured source is the unreadable-source
	// rejection with the source and the URL in the message.
	fetch3 := newFakeFetcher()
	fetch3.setErr(application.JavaChecksumsURLForTest(jdkOverrideVer, "x64"), errors.New("connection refused"))
	svc3, fs3 := newToolchainService(t, fetch3)
	seedToolchainYAML(t, fs3, toolchainFixture("    java-toolchain:\n      enabled: true\n      version: "+jdkOverrideVer+"\n"))
	_, err = generateSandbox(svc3, false)
	if !errors.Is(err, driving.ErrConfigValueInvalid) {
		t.Fatalf("generate (java-toolchain unreadable): err = %v, want a wrap of ErrConfigValueInvalid (exit 10)", err)
	}
	if !strings.Contains(err.Error(), "is not readable") || !strings.Contains(err.Error(), "Temurin") {
		t.Errorf("err = %v, want the unreadable wording and the source name", err)
	}
	assertNothingWritten(t, fs3)
}

// A `version:` value outside the release-version character set is unusable
// input (SPEC-054, S5): the run ends with the exit-10 class, nothing written,
// and the fetcher is not asked.
func TestGenerateDevcontainer_Toolchain_InvalidVersionValueExit10(t *testing.T) {
	t.Parallel()
	for _, bad := range []string{"1 2.3", "1.27;rm-rf", "$(danger)"} {
		bad := bad
		t.Run(bad, func(t *testing.T) {
			t.Parallel()
			fetch := newFakeFetcher()
			svc, fs := newToolchainService(t, fetch)
			seedToolchainYAML(t, fs, toolchainFixture("    go-toolchain:\n      enabled: true\n      version: \""+bad+"\"\n"))
			_, err := generateSandbox(svc, false)
			if !errors.Is(err, driving.ErrConfigValueInvalid) {
				t.Fatalf("generate with version %q: err = %v, want a wrap of ErrConfigValueInvalid (exit 10)", bad, err)
			}
			assertNothingWritten(t, fs)
			if n := len(fetch.calls); n != 0 {
				t.Errorf("the run asked %d source read(s) for an unusable version, want none", n)
			}
		})
	}
}

// The S8 conflict (SPEC-054): an enabled external `go` feature next to the
// active `go-toolchain` key ends the run with the exit-10 class and nothing
// written; the external `go` alone renders as before, and `go` next to
// `golangci-lint` is not a conflict.
func TestGenerateDevcontainer_Toolchain_GoFeatureConflictExit10(t *testing.T) {
	t.Parallel()
	t.Run("go plus go-toolchain is rejected", func(t *testing.T) {
		t.Parallel()
		svc, fs := newGenerateService(t)
		seedToolchainYAML(t, fs, toolchainFixture(goFeatureEntry, goToolchainEntry))
		_, err := generateSandbox(svc, false)
		if !errors.Is(err, driving.ErrConfigValueInvalid) {
			t.Fatalf("generate: err = %v, want a wrap of ErrConfigValueInvalid (exit 10)", err)
		}
		if !strings.Contains(err.Error(), "go-toolchain") {
			t.Errorf("err = %v, want the message to name both keys", err)
		}
		assertNothingWritten(t, fs)
	})
	t.Run("go alone renders", func(t *testing.T) {
		t.Parallel()
		svc, fs := newToolchainService(t, newFakeFetcher())
		seedToolchainYAML(t, fs, toolchainFixture(goFeatureEntry))
		if _, err := generateSandbox(svc, false); err != nil {
			t.Fatalf("generate: %v", err)
		}
		raw, _ := fs.ReadFile(devcontainerJSONPath())
		if _, ok := strippedJSONMap(t, raw)["containerEnv"]; ok {
			t.Errorf("devcontainer.json carries containerEnv although no pinned-toolchain key is active:\n%s", raw)
		}
	})
	t.Run("go plus golangci-lint renders", func(t *testing.T) {
		t.Parallel()
		fetch := newFakeFetcher()
		fetch.setBody(application.GolangciChecksumsURLForTest(golangciDefaultVer),
			checksumsFileFixture(golangciDefaultVer, golangciSHAAMD64, golangciSHAARM64))
		svc, fs := newToolchainService(t, fetch)
		seedToolchainYAML(t, fs, toolchainFixture(goFeatureEntry, golangciEntry))
		if _, err := generateSandbox(svc, false); err != nil {
			t.Fatalf("generate: %v", err)
		}
		raw, _ := fs.ReadFile(dockerfilePath())
		if !strings.Contains(string(raw), "ARG GOLANGCI_LINT_VERSION=") {
			t.Errorf("the golangci-lint block is missing next to the external go feature:\n%s", raw)
		}
	})
}

// The java conflict (SPEC-054, Muster Go-Fall): an enabled external `java`
// feature next to the active `java-toolchain` key ends the run with the
// exit-10 class and nothing written; the external `java` alone renders as
// before, and `java` next to the two Go keys is not a conflict.
func TestGenerateDevcontainer_Toolchain_JavaFeatureConflictExit10(t *testing.T) {
	t.Parallel()
	t.Run("java plus java-toolchain is rejected", func(t *testing.T) {
		t.Parallel()
		svc, fs := newGenerateService(t)
		seedToolchainYAML(t, fs, toolchainFixture(javaFeatureEntry, javaEntry))
		_, err := generateSandbox(svc, false)
		if !errors.Is(err, driving.ErrConfigValueInvalid) {
			t.Fatalf("generate: err = %v, want a wrap of ErrConfigValueInvalid (exit 10)", err)
		}
		if !strings.Contains(err.Error(), "java-toolchain") || !strings.Contains(err.Error(), "two JDKs") {
			t.Errorf("err = %v, want the message to name both keys and the two JDKs", err)
		}
		assertNothingWritten(t, fs)
	})
	t.Run("java alone renders", func(t *testing.T) {
		t.Parallel()
		svc, fs := newToolchainService(t, newFakeFetcher())
		seedToolchainYAML(t, fs, toolchainFixture(javaFeatureEntry))
		if _, err := generateSandbox(svc, false); err != nil {
			t.Fatalf("generate: %v", err)
		}
		raw, _ := fs.ReadFile(devcontainerJSONPath())
		if _, ok := strippedJSONMap(t, raw)["containerEnv"]; ok {
			t.Errorf("devcontainer.json carries containerEnv although no pinned-toolchain key is active:\n%s", raw)
		}
	})
	t.Run("no sdkman rejection case (F7)", func(t *testing.T) {
		// The F7 pin (cr-vorschlag-java-toolchain, Stelle 5): the external
		// `java` feature is SDKMAN-based ("Java + SDKMAN" in the catalogue;
		// the Feature-Metadaten 1.8.3 say "Installs Java, SDKMAN!", gelesen
		// 2026-11-02), no separate sdkman rejection case exists in the
		// catalogue or the toolchain machinery, and this pin catches a new
		// one.
		t.Parallel()
		name, errName := domain.NewFeatureName("java")
		if errName != nil {
			t.Fatalf("java is not a FeatureName: %v", errName)
		}
		entry, ok := application.FeatureForTest(name)
		if !ok {
			t.Fatal("the catalogue does not carry the external java feature")
		}
		if !strings.Contains(entry.ShortDesc, "SDKMAN") {
			t.Errorf("the catalogue entry of the external java feature does not name SDKMAN: %q", entry.ShortDesc)
		}
	})
	t.Run("java plus go-toolchain renders", func(t *testing.T) {
		t.Parallel()
		svc, fs := newToolchainService(t, newFakeFetcher())
		seedToolchainYAML(t, fs, toolchainFixture(javaFeatureEntry, goToolchainEntry))
		if _, err := generateSandbox(svc, false); err != nil {
			t.Fatalf("generate: %v", err)
		}
		raw, _ := fs.ReadFile(dockerfilePath())
		if !strings.Contains(string(raw), "ARG GO_VERSION=") {
			t.Errorf("the go-toolchain block is missing next to the external java feature:\n%s", raw)
		}
	})
}

// The notices (SPEC-054, S7/S9): one `warn` entry per active key with the
// LH-NFA-SEC-004 code, the default-pin or override wording, the checksum
// source at an override and the unsigned limit; neither
// [driving.WarningEntry.SecurityRelaxation] nor DegradationFallback (the CLI
// adapter drops them under `--quiet`, pinned at the command level).
func TestGenerateDevcontainer_Toolchain_Notices(t *testing.T) {
	t.Parallel()
	t.Run("one notice per key", func(t *testing.T) {
		t.Parallel()
		fetch := newFakeFetcher()
		fetch.setBody(application.GoChecksumsSourceURLForTest(), toolchainListing(goDefaultVersion))
		fetch.setBody(application.GolangciChecksumsURLForTest(golangciDefaultVer),
			checksumsFileFixture(golangciDefaultVer, golangciSHAAMD64, golangciSHAARM64))
		svc, fs := newToolchainService(t, fetch)
		seedToolchainYAML(t, fs, toolchainFixture(goToolchainEntry, golangciEntry, javaEntry))
		resp, err := generateSandbox(svc, false)
		if err != nil {
			t.Fatalf("generate: %v", err)
		}
		notices := toolchainNotices(resp.Warnings)
		if len(notices) != 3 {
			t.Fatalf("%d notices (%v), want one per key", len(notices), notices)
		}
		wantKeys := []string{"feature go-toolchain ", "feature golangci-lint ", "feature java-toolchain "}
		for _, key := range wantKeys {
			found := false
			for _, n := range notices {
				if strings.Contains(n.Message, key) {
					found = true
				}
			}
			if !found {
				t.Errorf("no notice carries %q:\n%v", key, notices)
			}
		}
		for _, n := range notices {
			if n.Code != "LH-NFA-SEC-004" || n.Level != "warn" {
				t.Errorf("notice = %+v, want code LH-NFA-SEC-004 at level warn", n)
			}
			if n.SecurityRelaxation || n.DegradationFallback {
				t.Errorf("notice = %+v, want no SecurityRelaxation and no DegradationFallback", n)
			}
		}
	})
	t.Run("override names the source", func(t *testing.T) {
		t.Parallel()
		svc, fs := newToolchainService(t, newToolchainOverrideFetcher())
		seedToolchainYAML(t, fs, toolchainFixture(
			"    go-toolchain:\n      enabled: true\n      version: "+goOverrideVersion+"\n"))
		resp, err := generateSandbox(svc, false)
		if err != nil {
			t.Fatalf("generate: %v", err)
		}
		notices := toolchainNotices(resp.Warnings)
		if len(notices) != 1 {
			t.Fatalf("%d notices (%v), want one", len(notices), notices)
		}
		for _, want := range []string{"override of the default pin", "go.dev/dl", "not signed"} {
			if !strings.Contains(notices[0].Message, want) {
				t.Errorf("override notice lacks %q:\n%s", want, notices[0].Message)
			}
		}
	})
	t.Run("default pin is named", func(t *testing.T) {
		t.Parallel()
		svc, fs := newToolchainService(t, newFakeFetcher())
		seedToolchainYAML(t, fs, toolchainFixture(golangciEntry))
		resp, err := generateSandbox(svc, false)
		if err != nil {
			t.Fatalf("generate: %v", err)
		}
		notices := toolchainNotices(resp.Warnings)
		if len(notices) != 1 || !strings.Contains(notices[0].Message, "default pin of the catalogue") {
			t.Fatalf("notices = %v, want the default-pin wording", notices)
		}
	})
}

// toolchainNotices returns the warn entries that name one of the three keys
// (the supply-chain notices of the pinned toolchain; SPEC-054).
func toolchainNotices(ws []driving.WarningEntry) []driving.WarningEntry {
	var out []driving.WarningEntry
	for _, w := range ws {
		if strings.Contains(w.Message, "feature go-toolchain ") || strings.Contains(w.Message, "feature golangci-lint ") || strings.Contains(w.Message, "feature java-toolchain ") {
			out = append(out, w)
		}
	}
	return out
}

// mustRead reads one file out of the fake FS or fails the test.
func mustRead(t *testing.T, fs *fakeFS, path string) []byte {
	t.Helper()
	body, err := fs.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return body
}

// assertNothingWritten pins the rejections of SPEC-054: the exit-10 class
// ends the run before the plan and write, so both devcontainer files are
// absent (the fixtures of the tests seed u-boot.yaml only).
func assertNothingWritten(t *testing.T, fs *fakeFS) {
	t.Helper()
	for _, path := range []string{devcontainerJSONPath(), dockerfilePath()} {
		if _, err := fs.ReadFile(path); err == nil {
			t.Errorf("the rejected run must write none of the two files; %s exists", path)
		}
	}
}

// mustReadRepoFile reads a file of the repo from the package dir of the test
// binary (cwd is `internal/hexagon/application`; three levels up is the repo
// root).
func mustReadRepoFile(t *testing.T, rel string) string {
	t.Helper()
	body, err := os.ReadFile(filepath.Join("..", "..", "..", rel))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	return string(body)
}

// newSet is the distinct-value set of a place→value map.
func newSet(v map[string]string) map[string]bool {
	out := map[string]bool{}
	for _, val := range v {
		out[val] = true
	}
	return out
}

// pinArgVersionOneFile extracts the version of one `ARG <name>=<version>` or
// `<name> ?=` line; the regexp is anchored at the name so the Go pin does not
// read the GH/GLAB/PODMAN/OLLAMA pins. A template-escape line (the override
// branch of the catalogue template) does not count: only the literal pin
// line is the place.
func pinArgVersionOneFile(t *testing.T, text, name, context string) string {
	t.Helper()
	re := regexp.MustCompile(`(?m)^ARG ` + name + `=(\S+)$`)
	for _, m := range re.FindAllStringSubmatch(text, -1) {
		if !strings.Contains(m[1], "{{") {
			return m[1]
		}
	}
	reQ := regexp.MustCompile(`(?m)^` + name + `\s*\??=\s*(\S+)$`)
	if m := reQ.FindStringSubmatch(text); m != nil {
		return m[1]
	}
	t.Fatalf("the %s file does not carry a %s pin line", context, name)
	return ""
}

// LH-FA-DEV-011 (SPEC-054): the pin values of the four places of the
// coupling stand equal — the root Dockerfile (ARG), repo.mk (?=), the build
// arguments of the catalogue template and the generated .devcontainer/
// Dockerfile of the own project — and a `version:` override of a project
// changes none of the four places: the interaction subtest reads the three
// static places before and after an override generate and pins the
// generated Dockerfile as the only carrier of the override.
func TestToolchainPinPlacesEqual(t *testing.T) {
	t.Parallel()
	dockerfile := mustReadRepoFile(t, "Dockerfile")
	repoMk := mustReadRepoFile(t, "repo.mk")
	template := mustReadRepoFile(t, filepath.Join("internal", "hexagon", "application", "templates", "devcontainer", "Dockerfile.tmpl"))
	generated := mustReadRepoFile(t, filepath.Join(".devcontainer", "Dockerfile"))

	places := map[string]string{
		"root Dockerfile":         dockerfile,
		"repo.mk":                 repoMk,
		"catalogue template":      template,
		"generated .devcontainer": generated,
	}
	values := map[string]map[string]string{
		"GO_VERSION":            {},
		"GOLANGCI_LINT_VERSION": {},
	}
	for name, text := range places {
		for pin := range values {
			values[pin][name] = pinArgVersionOneFile(t, text, pin, name)
		}
	}
	for pin, v := range values {
		if len(newSet(v)) != 1 {
			t.Errorf("the %s pins of the four places differ: %v", pin, v)
		}
		want := map[string]string{"GO_VERSION": goDefaultVersion, "GOLANGCI_LINT_VERSION": golangciDefaultVer}[pin]
		if v["root Dockerfile"] != want {
			t.Errorf("the %s pin = %q, want %q", pin, v["root Dockerfile"], want)
		}
	}

	t.Run("an override changes none of the four places", func(t *testing.T) {
		t.Parallel()
		before := []struct {
			name, rel, content string
		}{
			{"root Dockerfile", "Dockerfile", dockerfile},
			{"repo.mk", "repo.mk", repoMk},
			{"catalogue template", filepath.Join("internal", "hexagon", "application", "templates", "devcontainer", "Dockerfile.tmpl"), template},
			{"own u-boot.yaml", "u-boot.yaml", mustReadRepoFile(t, "u-boot.yaml")},
		}
		svc, fs := newToolchainService(t, newToolchainOverrideFetcher())
		seedToolchainYAML(t, fs, toolchainFixture(
			"    go-toolchain:\n      enabled: true\n      version: "+goOverrideVersion+"\n"))
		if _, err := generateSandbox(svc, true); err != nil {
			t.Fatalf("generate: %v", err)
		}
		df, _ := fs.ReadFile(dockerfilePath())
		for _, pin := range []string{"GO", "GOLANGCI_LINT"} {
			if strings.Contains(string(df), pin+"_VERSION="+map[string]string{"GO": goDefaultVersion, "GOLANGCI_LINT": golangciDefaultVer}[pin]) {
				t.Errorf("the generated Dockerfile carries the default %s pin although the project overrode it", pin)
			}
		}
		assertToolchainBlock(t, string(df), toolchainSpec{
			key: "go-toolchain", arg: "GO", version: goOverrideVersion, shaAMD64: goOverrideAMD64, shaARM64: goOverrideARM64,
			downloadURL: "https://go.dev/dl/go${GO_VERSION}.linux-${arch}.tar.gz"})
		for _, f := range before {
			after, err := os.ReadFile(filepath.Join("..", "..", "..", f.rel))
			if err != nil {
				t.Fatalf("re-read %s: %v", f.name, err)
			}
			if string(after) != f.content {
				t.Errorf("the %s changed by the generate run of the test project", f.name)
			}
		}
	})
}

// The egress guard (SPEC-054 trägt die Hosts der beiden Schlüssel nicht): the
// allowlist script of the sandbox profile is byte-identical with the keys
// enabled and without them — `u-boot` adds no host for the toolchains; the
// archive download belongs to the image build, which the restriction does
// not gate.
func TestGenerateDevcontainer_Toolchain_EgressAddsNoHost(t *testing.T) {
	t.Parallel()
	egressScript := func(features string) string {
		svc, fs := newToolchainService(t, newFakeFetcher())
		seedToolchainYAML(t, fs, "schemaVersion: 1\nproject:\n  name: t-uboot-gen\ndevcontainer:\n  enabled: true\n  profile: sandbox\n"+features+"  sandbox:\n    egress:\n      enabled: true\n")
		if _, err := generateSandbox(svc, false); err != nil {
			t.Fatalf("generate: %v", err)
		}
		body, err := fs.ReadFile(egressInitPath())
		if err != nil {
			t.Fatalf("read egress script: %v", err)
		}
		return string(body)
	}
	with := egressScript("  features:\n" + goToolchainEntry + golangciEntry)
	without := egressScript("  features:\n    git:\n      enabled: true\n")
	if with != without {
		t.Errorf("the toolchain keys changed the egress allowlist:\nwith:\n%s\nwithout:\n%s", with, without)
	}
}

// SPEC-054 (Eintragstyp): `u-boot doctor` reports no finding of
// `devcontainer.features.allowlist` and `devcontainer.features.drift` for
// the three pinned-toolchain keys — the catalogue keys render no key under
// `features` and are no orphans.
func TestDoctor_ToolchainKeys_NoFindings(t *testing.T) {
	t.Parallel()
	const yaml = "schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n  features:\n    go-toolchain:\n      enabled: true\n    golangci-lint:\n      enabled: true\n    java-toolchain:\n      enabled: true\n"
	for name, jsonBody := range map[string]string{
		"features empty": `{"features": {}}`,
		"no features":    `{"name": "demo"}`,
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			svc, fs, _, _, _ := newDoctorService(t)
			driftSeed(t, fs, yaml, jsonBody)
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

// SPEC-054 (S5/S6), the type boundary of the `version:` override: an entry
// of the type without its official checksum source rejects the override, and
// so does every entry outside the type; an entry of the type with its source
// accepts the override without a `source:` (the load path and the write path
// of `config set` share [validateDevcontainerFeatures]). The synthetic
// entries run through the bridge; the two real keys through the load path.
func TestFeatureCatalogue_ToolchainChecksumSources(t *testing.T) {
	t.Parallel()
	accepted := application.ToolchainVersionOverrideAcceptedForTest
	if !accepted(true, "the download listing of go.dev/dl", "1.26.10", "") {
		t.Errorf("an entry of the type with its checksum source must accept the override")
	}
	if !accepted(true, "the checksums file of the same release of the Temurin binaries", jdkOverrideVer, "") {
		t.Errorf("an entry of the type with the Temurin checksum source must accept the override")
	}
	for _, tc := range []struct {
		name            string
		dockerfileBlock bool
		checksumSource  string
		version         string
		sourceOverride  string
		want            bool
	}{
		{"type without source rejects", true, "no real key of the catalogue carries none; the case is synthetic", "9.9.9", "", false},
		{"outside the type rejects", true, "", "1.2.3", "", false},
		{"with a source override the block rule does not apply", true, "", "9.9.9", "https://ghcr.io/x/y", true},
		{"an OCI entry accepts without a checksum source", false, "", "1.2.3", "", true},
	} {
		if got := accepted(tc.dockerfileBlock, tc.checksumSource, tc.version, tc.sourceOverride); got != tc.want {
			t.Errorf("case %q: accepted = %v, want %v", tc.name, got, tc.want)
		}
	}
	// The real catalogue keys carry their source (the S6 rejection stays
	// synthetic); the load path accepts the overrides of the two keys and
	// keeps the rejection for a block entry without source.
	const withOverride = "schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n  features:\n    go-toolchain:\n      enabled: true\n      version: 1.26.10\n    golangci-lint:\n      enabled: true\n      version: v2.13.2\n    java-toolchain:\n      enabled: true\n      version: 21.0.9+10\n"
	if err := application.ValidateDevcontainerFeaturesForTest(t, []byte(withOverride)); err != nil {
		t.Errorf("the load path rejected the overrides of the two keys: %v", err)
	}
	const blockWithoutSource = "schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n  features:\n    gh:\n      enabled: true\n      version: \"9.9.9\"\n"
	if err := application.ValidateDevcontainerFeaturesForTest(t, []byte(blockWithoutSource)); err == nil {
		t.Errorf("the load path accepted a version override of a block entry without source (gh)")
	}
}

// The parsers and the URL forms of the two checksum sources (SPEC-054):
// the read matches the file names, the not-covered rejection and the
// malformed listing are the error cases, and the URL of the read stands in
// the forms of the two real sources (the go.dev/dl listing, the release
// checksums file of golangci/golangci-lint).
func TestToolchainChecksumParsing(t *testing.T) {
	t.Parallel()
	t.Run("listing", func(t *testing.T) {
		t.Parallel()
		pin, err := application.GoChecksumsForTest(toolchainListing(goDefaultVersion, goOverrideVersion), goOverrideVersion)
		if err != nil {
			t.Fatalf("read: %v", err)
		}
		want := application.ToolchainPinForTest{Version: goOverrideVersion, SHA256AMD64: goOverrideAMD64, SHA256ARM64: goOverrideARM64}
		if pin != want {
			t.Errorf("pin = %+v, want %+v", pin, want)
		}
		// The spelling with the `go` prefix reads the same source.
		pin, err = application.GoChecksumsForTest(toolchainListing(goDefaultVersion, goOverrideVersion), "go"+goOverrideVersion)
		if err != nil || pin.Version != goOverrideVersion {
			t.Errorf("read with the go prefix: pin = %+v err = %v, want the bare form", pin, err)
		}
		_, err = application.GoChecksumsForTest([]byte("not json"), goDefaultVersion)
		if err == nil || !strings.Contains(err.Error(), "is not the version listing") {
			t.Errorf("a malformed listing must carry the unreadable-source wording, got %v", err)
		}
		_, err = application.GoChecksumsForTest(toolchainListing(goDefaultVersion), goOverrideVersion)
		if err == nil || !strings.Contains(err.Error(), "does not cover") {
			t.Errorf("a listing without the version must be the not-covered rejection, got %v", err)
		}
	})
	t.Run("checksums file", func(t *testing.T) {
		t.Parallel()
		pin, err := application.GolangciChecksumsForTest(checksumsFileFixture(lintOverrideTag, lintOverrideAMD64, lintOverrideARM64), lintOverrideTag)
		if err != nil {
			t.Fatalf("read: %v", err)
		}
		want := application.ToolchainPinForTest{Version: lintOverrideTag, SHA256AMD64: lintOverrideAMD64, SHA256ARM64: lintOverrideARM64}
		if pin != want {
			t.Errorf("pin = %+v, want %+v", pin, want)
		}
		// The tag-less spelling parses to the same file names.
		pin, err = application.GolangciChecksumsForTest(checksumsFileFixture(lintOverrideTag, lintOverrideAMD64, lintOverrideARM64), strings.TrimPrefix(lintOverrideTag, "v"))
		if err != nil || pin.SHA256AMD64 != lintOverrideAMD64 {
			t.Errorf("read without the v prefix: pin = %+v err = %v, want the archive checksum", pin, err)
		}
		_, err = application.GolangciChecksumsForTest(checksumsFileFixture(lintOverrideTag, lintOverrideAMD64, lintOverrideARM64), lintOverrideTag+"x")
		if err == nil || !strings.Contains(err.Error(), "does not cover") {
			t.Errorf("a checksums file without the archive must be the not-covered rejection, got %v", err)
		}
	})
	t.Run("url forms", func(t *testing.T) {
		t.Parallel()
		if got, want := application.GolangciChecksumsURLForTest(lintOverrideTag),
			"https://github.com/golangci/golangci-lint/releases/download/v2.13.2/golangci-lint-2.13.2-checksums.txt"; got != want {
			t.Errorf("golangci checksums URL = %q, want %q", got, want)
		}
		if got, want := application.GoChecksumsSourceURLForTest(), "https://go.dev/dl/?mode=json&include=all"; got != want {
			t.Errorf("go listing URL = %q, want %q", got, want)
		}
		// The Temurin files: one per architecture, the tag carries the
		// `+` URL-encoded (jdk-21.0.9%2B10), and the version slug of the
		// asset name carries `_` for `+`.
		if got, want := application.JavaChecksumsURLForTest(jdkOverrideVer, "x64"),
			"https://github.com/adoptium/temurin21-binaries/releases/download/jdk-21.0.9%2B10/OpenJDK21U-jdk_x64_linux_hotspot_21.0.9_10.tar.gz.sha256.txt"; got != want {
			t.Errorf("java checksums URL (x64) = %q, want %q", got, want)
		}
		if got, want := application.JavaChecksumsURLForTest(jdkOverrideVer, "aarch64"),
			"https://github.com/adoptium/temurin21-binaries/releases/download/jdk-21.0.9%2B10/OpenJDK21U-jdk_aarch64_linux_hotspot_21.0.9_10.tar.gz.sha256.txt"; got != want {
			t.Errorf("java checksums URL (aarch64) = %q, want %q", got, want)
		}
		// The jdk- prefix spelling reads the same files.
		if got, want := application.JavaChecksumsURLForTest("jdk-21.0.9+10", "x64"),
			application.JavaChecksumsURLForTest(jdkOverrideVer, "x64"); got != want {
			t.Errorf("java checksums URL with the jdk- prefix = %q, want %q", got, want)
		}
	})
	t.Run("temurin checksums files", func(t *testing.T) {
		t.Parallel()
		fetch := newFakeFetcher()
		fetch.setBody(application.JavaChecksumsURLForTest(jdkOverrideVer, "x64"),
			jdkChecksumsFileFixture(jdkOverrideX64, "OpenJDK21U-jdk_x64_linux_hotspot_21.0.9_10.tar.gz"))
		fetch.setBody(application.JavaChecksumsURLForTest(jdkOverrideVer, "aarch64"),
			jdkChecksumsFileFixture(jdkOverrideA64, "OpenJDK21U-jdk_aarch64_linux_hotspot_21.0.9_10.tar.gz"))
		pin, err := application.JavaChecksumsForTest(context.Background(), fetch, jdkOverrideVer)
		if err != nil {
			t.Fatalf("read: %v", err)
		}
		want := application.ToolchainPinForTest{Version: jdkOverrideVer, SHA256AMD64: jdkOverrideX64, SHA256ARM64: jdkOverrideA64}
		if pin != want {
			t.Errorf("pin = %+v, want %+v", pin, want)
		}
		_, err = application.JavaChecksumsForTest(context.Background(), newFakeFetcher(), jdkOverrideVer)
		if err == nil || !strings.Contains(err.Error(), "is not readable") {
			t.Errorf("an unreadable checksums file must be the unreadable-source rejection, got %v", err)
		}
	})
}

// SPEC-054 (the normalisation of S4): the two spellings of a version read
// the same source — `go1.26.10` overrides `go-toolchain` and the tag-less
// spelling overrides `golangci-lint`; the generated Dockerfile carries the
// normalised forms (the ARG convention of the key).
func TestGenerateDevcontainer_Toolchain_OverrideSpellings(t *testing.T) {
	t.Parallel()
	fetch := newFakeFetcher()
	fetch.setBody(application.GoChecksumsSourceURLForTest(), toolchainListing(goDefaultVersion, goOverrideVersion))
	fetch.setBody(application.GolangciChecksumsURLForTest(lintOverrideTag),
		checksumsFileFixture(lintOverrideTag, lintOverrideAMD64, lintOverrideARM64))
	svc, fs := newToolchainService(t, fetch)
	seedToolchainYAML(t, fs, toolchainFixture(
		"    go-toolchain:\n      enabled: true\n      version: go"+goOverrideVersion+"\n",
		"    golangci-lint:\n      enabled: true\n      version: "+strings.TrimPrefix(lintOverrideTag, "v")+"\n"))
	if _, err := generateSandbox(svc, false); err != nil {
		t.Fatalf("generate: %v", err)
	}
	raw, _ := fs.ReadFile(dockerfilePath())
	df := string(raw)
	assertToolchainBlock(t, df, toolchainSpec{
		key: "go-toolchain", arg: "GO", version: goOverrideVersion, shaAMD64: goOverrideAMD64, shaARM64: goOverrideARM64,
		downloadURL: "https://go.dev/dl/go${GO_VERSION}.linux-${arch}.tar.gz"})
	assertToolchainBlock(t, df, toolchainSpec{
		key: "golangci-lint", arg: "GOLANGCI_LINT", version: lintOverrideTag, shaAMD64: lintOverrideAMD64, shaARM64: lintOverrideARM64,
		downloadURL: "https://github.com/golangci/golangci-lint/releases/download/${GOLANGCI_LINT_VERSION}/golangci-lint-${GOLANGCI_LINT_VERSION#v}-linux-${arch}.tar.gz"})
	assertToolchainJSON(t, mustRead(t, fs, devcontainerJSONPath()))
	if !strings.Contains(df, "ARG GOLANGCI_LINT_VERSION="+lintOverrideTag+"\n") {
		t.Errorf("the golangci-lint ARG carries not the v-form:\n%s", df)
	}
	t.Run("java spellings", func(t *testing.T) {
		// The Temurin tag form "jdk-21.0.9+10" and the bare form read the
		// same checksums files and normalise to the bare ARG form.
		t.Parallel()
		fetch := newToolchainOverrideFetcher()
		svc, fs := newToolchainService(t, fetch)
		seedToolchainYAML(t, fs, toolchainFixture(
			"    java-toolchain:\n      enabled: true\n      version: jdk-"+jdkOverrideVer+"\n"))
		if _, err := generateSandbox(svc, false); err != nil {
			t.Fatalf("generate: %v", err)
		}
		raw, _ := fs.ReadFile(dockerfilePath())
		assertToolchainBlock(t, string(raw), toolchainSpec{
			key: "java-toolchain", arg: "JDK", version: jdkOverrideVer, shaAMD64: jdkOverrideX64, shaARM64: jdkOverrideA64,
			downloadURL: "https://github.com/adoptium/temurin21-binaries/releases/download/jdk-$(printf '%s' \"$JDK_VERSION\" | sed 's/+/%2B/')/${asset}.tar.gz",
			caseLine:    `case "$arch" in amd64) slug="x64" sha256="${JDK_SHA256_AMD64}" ;; arm64) slug="aarch64" sha256="${JDK_SHA256_ARM64}" ;; *) echo "the JDK has no archive for $arch" >&2; exit 1 ;; esac`})
		assertNoContainerEnv(t, mustRead(t, fs, devcontainerJSONPath()))
	})
}
