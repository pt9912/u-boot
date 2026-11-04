package application_test

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// Golden cases and render behaviour of the project environment values
// `devcontainer.containerEnv` (lastenheft 0.32.0, Absatz
// Projekt-Umgebungswerte): the generate run carries the map into the
// `containerEnv` object of the generated `.devcontainer/devcontainer.json`
// beside the cache values of the two Go keys — without duplicate entries,
// idempotent (regeneration rewrites the entries, a removed entry
// disappears) — with the collision abweisung of the cache names and the
// disclosure per rendered value. The generated Dockerfile carries none of
// the values.

// containerEnvFixtureBlock is the `containerEnv:` block of the fixture:
// three entries in the sort order of the map keys.
const containerEnvFixtureBlock = "  containerEnv:\n" +
	"    BUILD_MODE: local\n" +
	"    GOFLAGS: \"-mod=readonly -buildvcs=false\"\n" +
	"    GOTOOLCHAIN: local\n"

// containerEnvValues is the value map of the fixture block.
func containerEnvValues() map[string]string {
	return map[string]string{
		"BUILD_MODE":  "local",
		"GOFLAGS":     "-mod=readonly -buildvcs=false",
		"GOTOOLCHAIN": "local",
	}
}

// containerEnvJSONGolden is the devcontainer.json of the default profile
// for a project with project environment values and no Go key: only the
// map entries under `containerEnv`, sorted, no cache values.
const containerEnvJSONGolden = "// BEGIN U-BOOT MANAGED BLOCK: init\n" +
	"{\n" +
	"  \"name\": \"t-uboot-gen\",\n" +
	"  \"build\": {\n" +
	"    \"dockerfile\": \"./Dockerfile\",\n" +
	"    \"context\": \".\"\n" +
	"  },\n" +
	"  \"runArgs\": [\"--env-file=${localWorkspaceFolder}/.devcontainer/devcontainer.env\"],\n" +
	"  \"containerEnv\": {\n" +
	"    \"BUILD_MODE\": \"local\",\n" +
	"    \"GOFLAGS\": \"-mod=readonly -buildvcs=false\",\n" +
	"    \"GOTOOLCHAIN\": \"local\"\n" +
	"  },\n" +
	"  \"remoteUser\": \"vscode\"\n" +
	"}\n" +
	"// END U-BOOT MANAGED BLOCK: init\n"

// mixedCacheJSONGolden is the devcontainer.json of the default profile for
// a project with `go-toolchain` and one project value: the cache values stand
// once, the project value follows behind them, sorted, without duplicate
// entries.
const mixedCacheJSONGolden = "// BEGIN U-BOOT MANAGED BLOCK: init\n" +
	"{\n" +
	"  \"name\": \"t-uboot-gen\",\n" +
	"  \"build\": {\n" +
	"    \"dockerfile\": \"./Dockerfile\",\n" +
	"    \"context\": \".\"\n" +
	"  },\n" +
	"  \"runArgs\": [\"--env-file=${localWorkspaceFolder}/.devcontainer/devcontainer.env\"],\n" +
	"  \"containerEnv\": {\n" +
	"    \"GOCACHE\": \"${containerWorkspaceFolder}/.cache/go-build\",\n" +
	"    \"GOMODCACHE\": \"${containerWorkspaceFolder}/.cache/go-mod\",\n" +
	"    \"BUILD_MODE\": \"local\"\n" +
	"  },\n" +
	"  \"remoteUser\": \"vscode\"\n" +
	"}\n" +
	"// END U-BOOT MANAGED BLOCK: init\n"

// containerEnvFixture is a u-boot.yaml fixture with the given containerEnv
// block behind the devcontainer head and the feature entries of
// toolchainFixture.
func containerEnvFixture(envBlock string, featureEntries ...string) string {
	return toolchainFixture(featureEntries...) + envBlock
}

// assertContainerEnvJSON pins the `containerEnv` map of a generated
// devcontainer.json: exactly the given entries, no others.
func assertContainerEnvJSON(t *testing.T, body []byte, want map[string]string) {
	t.Helper()
	m := strippedJSONMap(t, body)
	env, ok := m["containerEnv"].(map[string]any)
	if !ok {
		t.Fatalf("devcontainer.json carries no containerEnv:\n%s", string(body))
	}
	got := map[string]string{}
	for k, v := range env {
		s, isStr := v.(string)
		if !isStr {
			t.Fatalf("containerEnv entry %q is no string: %v\n%s", k, v, string(body))
		}
		got[k] = s
	}
	if len(got) != len(want) {
		t.Errorf("containerEnv = %v, want %v", got, want)
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("containerEnv[%q] = %q, want %q", k, got[k], v)
		}
	}
}

// assertCacheValuesOnce pins the cache entries of the two Go keys in the
// raw body: each name stands exactly once (no duplicate entry in the
// merged object).
func assertCacheValuesOnce(t *testing.T, body []byte) {
	t.Helper()
	for _, name := range []string{"GOCACHE", "GOMODCACHE"} {
		if n := strings.Count(string(body), name); n != 1 {
			t.Errorf("%s appears %d times, want once:\n%s", name, n, string(body))
		}
	}
}

// Fresh state per profile (lastenheft 0.32.0, Beleg-Tabelle row of the
// golden cases): the map lands sorted under `containerEnv`, the
// Dockerfile carries none of the values, and the second run is a NoOp.
func TestGenerateDevcontainer_ContainerEnv_InBothProfiles(t *testing.T) {
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
		t.Run(pname, func(t *testing.T) {
			t.Parallel()
			svc, fs := newToolchainService(t, newFakeFetcher())
			seedToolchainYAML(t, fs, p.head+containerEnvFixtureBlock)
			if _, err := generateSandbox(svc, p.sandbox); err != nil {
				t.Fatalf("generate: %v", err)
			}
			body, _ := fs.ReadFile(devcontainerJSONPath())
			assertContainerEnvJSON(t, body, containerEnvValues())
			// Sorted entry order in the raw bytes: the byte order of the
			// three lines follows the map-key order.
			positions := []int{
				strings.Index(string(body), `"BUILD_MODE":`),
				strings.Index(string(body), `"GOFLAGS":`),
				strings.Index(string(body), `"GOTOOLCHAIN":`),
			}
			for i := 1; i < len(positions); i++ {
				if positions[i-1] < 0 || positions[i-1] > positions[i] {
					t.Fatalf("containerEnv entries are not sorted in the file:\n%s", string(body))
				}
			}
			df, err := fs.ReadFile(dockerfilePath())
			if err != nil {
				t.Fatalf("read Dockerfile: %v", err)
			}
			for _, name := range []string{"BUILD_MODE", "GOFLAGS", "GOTOOLCHAIN"} {
				if strings.Contains(string(df), name+"=") {
					t.Errorf("the Dockerfile carries %s although the value belongs to containerEnv alone:\n%s", name, string(df))
				}
			}
			if r2, err := generateSandbox(svc, p.sandbox); err != nil || r2.Action != driving.GenerateActionNoOp {
				t.Errorf("second run: action=%v err=%v, want NoOp", r2.Action, err)
			}
		})
	}
	t.Run("golden default profile", func(t *testing.T) {
		t.Parallel()
		svc, fs := newToolchainService(t, newFakeFetcher())
		seedToolchainYAML(t, fs, containerEnvFixture(containerEnvFixtureBlock))
		if _, err := generateSandbox(svc, false); err != nil {
			t.Fatalf("generate: %v", err)
		}
		raw, _ := fs.ReadFile(devcontainerJSONPath())
		if got := string(raw); got != containerEnvJSONGolden {
			t.Errorf("devcontainer.json is not the golden:\n--- got ---\n%s\n--- want ---\n%s", got, containerEnvJSONGolden)
		}
	})
}

// The Mischfall (Beleg-Tabelle, golden cases): the cache values stand
// exactly once and the project values follow behind them — the generated
// object carries one entry per name, no duplicates.
func TestGenerateDevcontainer_ContainerEnv_MixedWithTheCacheValues(t *testing.T) {
	t.Parallel()
	t.Run("golden with go-toolchain", func(t *testing.T) {
		t.Parallel()
		fetch := newFakeFetcher()
		svc, fs := newToolchainService(t, fetch)
		seedToolchainYAML(t, fs, containerEnvFixture("  containerEnv:\n    BUILD_MODE: local\n", goToolchainEntry))
		if _, err := generateSandbox(svc, false); err != nil {
			t.Fatalf("generate: %v", err)
		}
		raw, _ := fs.ReadFile(devcontainerJSONPath())
		if got := string(raw); got != mixedCacheJSONGolden {
			t.Errorf("devcontainer.json is not the golden:\n--- got ---\n%s\n--- want ---\n%s", got, mixedCacheJSONGolden)
		}
	})
	t.Run("golangci-lint only, sandbox flag", func(t *testing.T) {
		t.Parallel()
		svc, fs := newToolchainService(t, newFakeFetcher())
		seedToolchainYAML(t, fs, containerEnvFixture("  containerEnv:\n    GOFLAGS: \"-mod=readonly\"\n", golangciEntry))
		if _, err := generateSandbox(svc, true); err != nil {
			t.Fatalf("generate: %v", err)
		}
		body, _ := fs.ReadFile(devcontainerJSONPath())
		assertContainerEnvJSON(t, body, map[string]string{
			"GOCACHE":    "${containerWorkspaceFolder}/.cache/go-build",
			"GOMODCACHE": "${containerWorkspaceFolder}/.cache/go-mod",
			"GOFLAGS":    "-mod=readonly",
		})
		assertCacheValuesOnce(t, body)
		if r2, err := generateSandbox(svc, true); err != nil || r2.Action != driving.GenerateActionNoOp {
			t.Errorf("second run: action=%v err=%v, want NoOp", r2.Action, err)
		}
	})
	t.Run("java-toolchain takes no cache space", func(t *testing.T) {
		t.Parallel()
		svc, fs := newToolchainService(t, newFakeFetcher())
		seedToolchainYAML(t, fs, containerEnvFixture(containerEnvFixtureBlock, javaEntry))
		if _, err := generateSandbox(svc, false); err != nil {
			t.Fatalf("generate: %v", err)
		}
		body, _ := fs.ReadFile(devcontainerJSONPath())
		assertContainerEnvJSON(t, body, containerEnvValues())
	})
}

// Without the map the render stays byte-identical to the stand before
// this key: the golden of the two cache-only cases keeps its bytes and
// the no-keys project carries no containerEnv object.
func TestGenerateDevcontainer_ContainerEnv_WithoutTheMapByteIdentical(t *testing.T) {
	t.Parallel()
	t.Run("cache values golden unchanged", func(t *testing.T) {
		t.Parallel()
		svc, fs := newToolchainService(t, newFakeFetcher())
		seedToolchainYAML(t, fs, toolchainFixture(goToolchainEntry))
		if _, err := generateSandbox(svc, false); err != nil {
			t.Fatalf("generate: %v", err)
		}
		raw, _ := fs.ReadFile(devcontainerJSONPath())
		if got := string(raw); got != toolchainJSONGolden {
			t.Errorf("devcontainer.json changed against the cache-only golden:\n--- got ---\n%s\n--- want ---\n%s", got, toolchainJSONGolden)
		}
	})
	t.Run("no keys, no map", func(t *testing.T) {
		t.Parallel()
		svc, fs := newToolchainService(t, newFakeFetcher())
		seedToolchainYAML(t, fs, toolchainFixture())
		if _, err := generateSandbox(svc, false); err != nil {
			t.Fatalf("generate: %v", err)
		}
		raw, _ := fs.ReadFile(devcontainerJSONPath())
		if got := string(raw); got != javaJSONGolden {
			t.Errorf("devcontainer.json carries the containerEnv object without map and keys:\n--- got ---\n%s\n--- want ---\n%s", got, javaJSONGolden)
		}
	})
}

// Idempotency (Beleg-Tabelle): the regeneration rewrites the entries and
// a removed entry disappears; adding the entry back renders it again.
func TestGenerateDevcontainer_ContainerEnv_RemovedEntryDisappears(t *testing.T) {
	t.Parallel()
	svc, fs := newToolchainService(t, newFakeFetcher())
	const twoEntries = "  containerEnv:\n    BUILD_MODE: local\n    GOTOOLCHAIN: local\n"
	const oneEntry = "  containerEnv:\n    BUILD_MODE: local\n"
	seedToolchainYAML(t, fs, containerEnvFixture(twoEntries))
	if _, err := generateSandbox(svc, false); err != nil {
		t.Fatalf("generate: %v", err)
	}
	body, _ := fs.ReadFile(devcontainerJSONPath())
	assertContainerEnvJSON(t, body, map[string]string{"BUILD_MODE": "local", "GOTOOLCHAIN": "local"})

	// The hand-edited config without GOTOOLCHAIN: the regeneration drops
	// the entry from the managed block.
	seedToolchainYAML(t, fs, containerEnvFixture(oneEntry))
	r2, err := generateSandbox(svc, false)
	if err != nil {
		t.Fatalf("regenerate: %v", err)
	}
	if r2.Action == driving.GenerateActionNoOp {
		t.Errorf("the regeneration kept the removed entry: action NoOp")
	}
	body, _ = fs.ReadFile(devcontainerJSONPath())
	assertContainerEnvJSON(t, body, map[string]string{"BUILD_MODE": "local"})

	// The third run with the same one-entry map is a NoOp again.
	if r3, err := generateSandbox(svc, false); err != nil || r3.Action != driving.GenerateActionNoOp {
		t.Errorf("third run: action=%v err=%v, want NoOp", r3.Action, err)
	}
}

// The collision abweisung (Beleg-Tabelle): a project value with one of the
// cache directory names beside an active Go key ends the run with the
// exit-10 class before anything is written; without the keys, and beside
// the java-toolchain key (it carries no containerEnv value), the same
// value is allowed.
func TestGenerateDevcontainer_ContainerEnv_CollisionExits10(t *testing.T) {
	t.Parallel()
	t.Run("GOCACHE with go-toolchain", func(t *testing.T) {
		t.Parallel()
		svc, fs := newGenerateService(t)
		seedToolchainYAML(t, fs, containerEnvFixture("  containerEnv:\n    GOCACHE: /tmp/go-cache\n", goToolchainEntry))
		_, err := generateSandbox(svc, false)
		if !errors.Is(err, driving.ErrConfigValueInvalid) {
			t.Fatalf("generate: err = %v, want a wrap of ErrConfigValueInvalid (exit 10)", err)
		}
		if !strings.Contains(err.Error(), "devcontainer.containerEnv.GOCACHE") || !strings.Contains(err.Error(), "go-toolchain") {
			t.Errorf("err = %v, want the message to name the entry and the active key", err)
		}
		assertNothingWritten(t, fs)
	})
	t.Run("GOMODCACHE with golangci-lint", func(t *testing.T) {
		t.Parallel()
		svc, fs := newGenerateService(t)
		seedToolchainYAML(t, fs, containerEnvFixture("  containerEnv:\n    GOMODCACHE: /tmp/go-mod\n", golangciEntry))
		_, err := generateSandbox(svc, false)
		if !errors.Is(err, driving.ErrConfigValueInvalid) {
			t.Fatalf("generate: err = %v, want a wrap of ErrConfigValueInvalid (exit 10)", err)
		}
		assertNothingWritten(t, fs)
	})
	t.Run("both cache names with both keys", func(t *testing.T) {
		t.Parallel()
		svc, fs := newToolchainService(t, newFakeFetcher())
		seedToolchainYAML(t, fs, containerEnvFixture(
			"  containerEnv:\n    GOCACHE: /tmp/go-cache\n    GOMODCACHE: /tmp/go-mod\n",
			goToolchainEntry, golangciEntry))
		_, err := generateSandbox(svc, false)
		if !errors.Is(err, driving.ErrConfigValueInvalid) {
			t.Fatalf("generate: err = %v, want a wrap of ErrConfigValueInvalid (exit 10)", err)
		}
		if !strings.Contains(err.Error(), "go-toolchain and golangci-lint") {
			t.Errorf("err = %v, want the message to name both active keys", err)
		}
		assertNothingWritten(t, fs)
	})
	t.Run("sandbox flag rejection stands before the write", func(t *testing.T) {
		t.Parallel()
		svc, fs := newGenerateService(t)
		seedToolchainYAML(t, fs, containerEnvFixture("  containerEnv:\n    GOCACHE: /tmp/go-cache\n", goToolchainEntry))
		_, err := generateSandbox(svc, true)
		if !errors.Is(err, driving.ErrConfigValueInvalid) {
			t.Fatalf("generate --sandbox: err = %v, want a wrap of ErrConfigValueInvalid (exit 10)", err)
		}
		assertNothingWritten(t, fs)
	})
	t.Run("without the keys the value renders", func(t *testing.T) {
		t.Parallel()
		svc, fs := newToolchainService(t, newFakeFetcher())
		seedToolchainYAML(t, fs, containerEnvFixture("  containerEnv:\n    GOCACHE: /tmp/go-cache\n"))
		if _, err := generateSandbox(svc, false); err != nil {
			t.Fatalf("generate: %v", err)
		}
		body, _ := fs.ReadFile(devcontainerJSONPath())
		assertContainerEnvJSON(t, body, map[string]string{"GOCACHE": "/tmp/go-cache"})
	})
	t.Run("beside java-toolchain the cache name renders", func(t *testing.T) {
		t.Parallel()
		svc, fs := newToolchainService(t, newFakeFetcher())
		seedToolchainYAML(t, fs, containerEnvFixture("  containerEnv:\n    GOMODCACHE: /tmp/go-mod\n", javaEntry))
		if _, err := generateSandbox(svc, false); err != nil {
			t.Fatalf("generate: %v", err)
		}
		body, _ := fs.ReadFile(devcontainerJSONPath())
		assertContainerEnvJSON(t, body, map[string]string{"GOMODCACHE": "/tmp/go-mod"})
	})
	t.Run("the java-feature next to a colliding Go key does not lift the rejection", func(t *testing.T) {
		t.Parallel()
		svc, fs := newGenerateService(t)
		seedToolchainYAML(t, fs, containerEnvFixture("  containerEnv:\n    GOCACHE: /tmp/go-cache\n", goToolchainEntry, javaEntry))
		if _, err := generateSandbox(svc, false); !errors.Is(err, driving.ErrConfigValueInvalid) {
			t.Errorf("generate: err = %v, want the collision rejection", err)
		}
		assertNothingWritten(t, fs)
	})
}

// The disclosure (Beleg-Tabelle): one plain entry per rendered project
// value with the LH-FA-DEV-001 code, the value named, without
// SecurityRelaxation and DegradationFallback; without the map there is
// none.
func TestGenerateDevcontainer_ContainerEnv_Notices(t *testing.T) {
	t.Parallel()
	t.Run("one per entry", func(t *testing.T) {
		t.Parallel()
		svc, fs := newToolchainService(t, newFakeFetcher())
		seedToolchainYAML(t, fs, containerEnvFixture(containerEnvFixtureBlock, goToolchainEntry))
		resp, err := generateSandbox(svc, false)
		if err != nil {
			t.Fatalf("generate: %v", err)
		}
		var notices []driving.WarningEntry
		for _, w := range resp.Warnings {
			if strings.Contains(w.Message, "project environment value ") {
				notices = append(notices, w)
			}
		}
		if len(notices) != 3 {
			t.Fatalf("%d notices (%v), want one per entry", len(notices), notices)
		}
		for name, value := range containerEnvValues() {
			found := false
			for _, n := range notices {
				if strings.Contains(n.Message, name+"="+value) {
					found = true
				}
			}
			if !found {
				t.Errorf("no notice names %s=%s:\n%v", name, value, notices)
			}
		}
		for _, n := range notices {
			if n.Code != "LH-FA-DEV-001" || n.Level != "warn" {
				t.Errorf("notice = %+v, want code LH-FA-DEV-001 at level warn", n)
			}
			if n.SecurityRelaxation || n.DegradationFallback {
				t.Errorf("notice = %+v, want no SecurityRelaxation and no DegradationFallback", n)
			}
		}
	})
	t.Run("without the map none (the toolchain notices stay alone)", func(t *testing.T) {
		t.Parallel()
		svc, fs := newToolchainService(t, newFakeFetcher())
		seedToolchainYAML(t, fs, toolchainFixture(goToolchainEntry))
		resp, err := generateSandbox(svc, false)
		if err != nil {
			t.Fatalf("generate: %v", err)
		}
		for _, w := range resp.Warnings {
			if strings.Contains(w.Message, "project environment value ") {
				t.Errorf("a run without the map emitted a project-value notice: %+v", w)
			}
		}
	})
}

// The fields of the map are hand-editable: a value with a quote, a
// backslash and a newline renders escaped, so the generated file stays
// valid JSON, the double run is a NoOp and the escaped form is stable.
func TestGenerateDevcontainer_ContainerEnv_EscapesJSONContent(t *testing.T) {
	t.Parallel()
	svc, fs := newToolchainService(t, newFakeFetcher())
	const tricky = "q\"uote\\back\nline"
	seedToolchainYAML(t, fs, containerEnvFixture(fmt.Sprintf("  containerEnv:\n    TRICKY: %q\n", tricky)))
	if _, err := generateSandbox(svc, false); err != nil {
		t.Fatalf("generate: %v", err)
	}
	body, _ := fs.ReadFile(devcontainerJSONPath())
	assertContainerEnvJSON(t, body, map[string]string{"TRICKY": tricky})
	if n := strings.Count(string(body), "TRICKY"); n != 1 {
		t.Errorf("TRICKY appears %d times:\n%s", n, string(body))
	}
	if r2, err := generateSandbox(svc, false); err != nil || r2.Action != driving.GenerateActionNoOp {
		t.Errorf("second run: action=%v err=%v, want NoOp", r2.Action, err)
	}
}

// The own devcontainer of the u-boot project is converted: the three
// values of the hand-written ENV rest behind its managed block
// (BUILD_MODE, GOFLAGS, GOTOOLCHAIN) move into `devcontainer.containerEnv`
// of the project config, their lines render into the `containerEnv`
// object of the generated .devcontainer/devcontainer.json beside the two
// cache values, and the Dockerfile rest carries none of them.
func TestOwnDevcontainer_ProjectValuesUnderContainerEnv(t *testing.T) {
	t.Parallel()
	yaml := mustReadRepoFile(t, "u-boot.yaml")
	for _, want := range []string{
		"containerEnv:",
		"BUILD_MODE: local",
		"GOTOOLCHAIN: local",
	} {
		if !strings.Contains(yaml, want) {
			t.Errorf("u-boot.yaml of the project lacks %q:\n%s", want, yaml)
		}
	}
	if !strings.Contains(yaml, `GOFLAGS: "-mod=readonly -buildvcs=false"`) &&
		!strings.Contains(yaml, "GOFLAGS: -mod=readonly -buildvcs=false") {
		t.Errorf("u-boot.yaml of the project lacks the GOFLAGS value:\n%s", yaml)
	}
	jsonBody := mustReadRepoFile(t, filepath.Join(".devcontainer", "devcontainer.json"))
	assertContainerEnvJSON(t, []byte(jsonBody), map[string]string{
		"GOCACHE":     "${containerWorkspaceFolder}/.cache/go-build",
		"GOMODCACHE":  "${containerWorkspaceFolder}/.cache/go-mod",
		"BUILD_MODE":  "local",
		"GOFLAGS":     "-mod=readonly -buildvcs=false",
		"GOTOOLCHAIN": "local",
	})
	df := mustReadRepoFile(t, filepath.Join(".devcontainer", "Dockerfile"))
	for _, name := range []string{"ENV BUILD_MODE", "BUILD_MODE=", "ENV GOFLAGS", "GOTOOLCHAIN=local"} {
		if strings.Contains(df, name) {
			t.Errorf("the hand-written rest of .devcontainer/Dockerfile still carries %q:\n%s", name, df)
		}
	}
}
