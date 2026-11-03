package application_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// Validation and write path of the project environment key
// `devcontainer.containerEnv` (lastenheft 0.32.0, Absatz
// Projekt-Umgebungswerte): the entries are written per name with
// `u-boot config set devcontainer.containerEnv.<name>`, the name is
// a valid dev-container environment name, the value is a string
// (an empty value included), and every fachlicher Fehler ends the
// call with the exit-10 class before anything is written. The base
// path itself has no whole-map write — the constructor rejects it
// (TestNewConfigPath_ContainerEnv in internal/hexagon/domain).

func configUbootYAMLPath() string {
	return configTestBaseDir + "/u-boot.yaml"
}

// setContainerEnvEntryAt seeds the config and writes one entry; the
// returned body is u-boot.yaml after the call.
func setContainerEnvEntryAt(t *testing.T, seed, path, value string) (string, error) {
	t.Helper()
	svc, fs := newConfigService(t)
	if err := fs.WriteFile(configUbootYAMLPath(), []byte(seed), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
	_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir, Path: mustConfigPath(t, path), Value: value,
	})
	body, _ := fs.ReadFile(configUbootYAMLPath())
	return string(body), err
}

// The happy path: the entry lands under `devcontainer.containerEnv`
// — also when the `devcontainer:` block does not exist yet — an
// empty value is a set value, and a second call with the same value
// is a NoOp; the comments of the surrounding document stay.
func TestConfigSet_ContainerEnv_HappyPathAndNoOp(t *testing.T) {
	t.Parallel()
	const seed = "schemaVersion: 1\n# project comment\nproject:\n  name: t-uboot-config\n"
	body, err := setContainerEnvEntryAt(t, seed, "devcontainer.containerEnv.BUILD_MODE", "local")
	if err != nil {
		t.Fatalf("set: %v", err)
	}
	for _, want := range []string{"# project comment", "devcontainer:", "  containerEnv:", "    BUILD_MODE: local"} {
		if !strings.Contains(body, want) {
			t.Errorf("u-boot.yaml lacks %q:\n%s", want, body)
		}
	}

	// An empty value is a string value: the entry is written with it.
	emptyBody, err := setContainerEnvEntryAt(t, body, "devcontainer.containerEnv.FOO", "")
	if err != nil {
		t.Fatalf("set empty value: %v", err)
	}
	if !strings.Contains(emptyBody, `FOO: ""`) {
		t.Errorf("the empty value did not land as an entry:\n%s", emptyBody)
	}

	svc, fs := newConfigService(t)
	if err := fs.WriteFile(configUbootYAMLPath(), []byte(emptyBody), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
	get, err := svc.Get(context.Background(), driving.ConfigGetRequest{
		BaseDir: configTestBaseDir, Path: mustConfigPath(t, "devcontainer.containerEnv.FOO"),
	})
	if err != nil || get.Value != "" {
		t.Errorf("get of the empty entry: value = %q err = %v, want the empty value", get.Value, err)
	}

	// NoOp: the same value over an existing entry changes nothing.
	noOpBody, err := setContainerEnvEntryAt(t, emptyBody, "devcontainer.containerEnv.BUILD_MODE", "local")
	if err != nil {
		t.Fatalf("no-op set: %v", err)
	}
	if noOpBody != emptyBody {
		t.Errorf("the equal-value call changed the file:\n--- before ---\n%s\n--- after ---\n%s", emptyBody, noOpBody)
	}
}

// The SetMany route applies the item with the same stages — the
// comment-preserving patch included.
func TestConfigSetMany_ContainerEnv_Entry(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	const seed = "schemaVersion: 1\nproject:\n  name: t-uboot-config\ndevcontainer:\n  enabled: true\n"
	if err := fs.WriteFile(configUbootYAMLPath(), []byte(seed), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
	resp, err := svc.SetMany(context.Background(), driving.ConfigSetManyRequest{
		BaseDir: configTestBaseDir,
		Items: []driving.ConfigSetItem{
			{Path: mustConfigPath(t, "devcontainer.containerEnv.BUILD_MODE"), Value: "local"},
			{Path: mustConfigPath(t, "devcontainer.enabled"), Value: "true"},
		},
	})
	if err != nil {
		t.Fatalf("set many: %v", err)
	}
	if len(resp.Entries) != 2 {
		t.Fatalf("entries = %d, want 2:\n%v", len(resp.Entries), resp.Entries)
	}
	body, _ := fs.ReadFile(configUbootYAMLPath())
	for _, want := range []string{"  enabled: true", "  containerEnv:", "    BUILD_MODE: local"} {
		if !strings.Contains(string(body), want) {
			t.Errorf("u-boot.yaml lacks %q:\n%s", want, body)
		}
	}
}

// The name is a dev-container environment name: every other name is
// a fachlicher Fehler of the exit-10 class, and the file stays
// byte-identical.
func TestConfigSet_ContainerEnv_InvalidNameExits10(t *testing.T) {
	t.Parallel()
	const seed = "schemaVersion: 1\nproject:\n  name: t-uboot-config\ndevcontainer:\n  enabled: true\n  containerEnv:\n    GOOD: ok\n"
	for _, bad := range []string{"9START", "HAS SPACE", "WITH.DOT", "dash-name", "ÜBER", "WITH$", "WITH\n"} {
		path := "devcontainer.containerEnv." + bad
		body, err := setContainerEnvEntryAt(t, seed, path, "x")
		if !errors.Is(err, driving.ErrConfigValueInvalid) {
			t.Errorf("set %q: err = %v, want a wrap of ErrConfigValueInvalid (exit 10)", path, err)
			continue
		}
		if body != seed {
			t.Errorf("the rejected call wrote to the file:\n--- before ---\n%s\n--- after ---\n%s", seed, body)
		}
		if !strings.Contains(err.Error(), path) && !strings.Contains(err.Error(), "not a plain name") {
			t.Errorf("err = %v, want the message to name the path or the rule", err)
		}
	}
}

// The extraction: an existing entry reads verbatim (an empty value
// included), an absent entry is NotSet with the repair hint, and
// the list carries one path per name of the map.
func TestConfigGet_ContainerEnv_Extraction(t *testing.T) {
	t.Parallel()
	svc, fs := newConfigService(t)
	const seed = "schemaVersion: 1\nproject:\n  name: t-uboot-config\ndevcontainer:\n  enabled: true\n  containerEnv:\n    BUILD_MODE: local\n    PADDING: \"\"\n"
	if err := fs.WriteFile(configUbootYAMLPath(), []byte(seed), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
	get, err := svc.Get(context.Background(), driving.ConfigGetRequest{
		BaseDir: configTestBaseDir, Path: mustConfigPath(t, "devcontainer.containerEnv.BUILD_MODE"),
	})
	if err != nil || get.Value != "local" {
		t.Errorf("get BUILD_MODE: value = %q err = %v, want local", get.Value, err)
	}
	if _, err := svc.Get(context.Background(), driving.ConfigGetRequest{
		BaseDir: configTestBaseDir, Path: mustConfigPath(t, "devcontainer.containerEnv.ABSENT"),
	}); !errors.Is(err, driving.ErrConfigValueNotSet) {
		t.Errorf("get ABSENT: err = %v, want ErrConfigValueNotSet", err)
	}
	list, err := svc.List(context.Background(), driving.ConfigListRequest{BaseDir: configTestBaseDir})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	var found []string
	for _, e := range list.Entries {
		if strings.HasPrefix(e.Path.String(), "devcontainer.containerEnv.") {
			found = append(found, e.Path.String())
		}
	}
	// The list keeps the shared filter of the paths with a value: the
	// empty-valued entry stays out of the list — the get reads it
	// verbatim, the list names the valued entry.
	if len(found) != 1 || found[0] != "devcontainer.containerEnv.BUILD_MODE" {
		t.Errorf("list carries %v, want the valued containerEnv entry", found)
	}
}

// The load path accepts a hand-edited map: the schema read of the
// generate run parses the entries — the map is a plain
// name-to-string map beside the known subtrees.
func TestContainerEnv_LoadPathParsesTheMap(t *testing.T) {
	t.Parallel()
	const seed = "schemaVersion: 1\nproject:\n  name: t-uboot-config\ndevcontainer:\n  enabled: true\n  containerEnv:\n    BUILD_MODE: local\n    GOFLAGS: \"-mod=readonly\"\n"
	if err := application.ValidateDevcontainerFeaturesForTest(t, []byte(seed)); err != nil {
		t.Errorf("the schema read rejected the containerEnv map: %v", err)
	}
}
