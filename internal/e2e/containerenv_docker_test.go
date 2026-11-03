//go:build docker

// End-to-end check of the project environment values `devcontainer.containerEnv`
// (lastenheft 0.32.0, LH-FA-DEV-001): the values reach the started container
// as the environment of the container processes — resolved from the
// `containerEnv` object of the generated devcontainer.json — and they are
// not part of the image: the probe of a bare `docker run` of the built image
// sees none of them, and the generated Dockerfile carries none of the names.

package e2e_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The containerEnv values of the measurement: BUILD_MODE mirrors the
// project value the own devcontainer of u-boot carries; the marker name
// stays outside every other env source of the image.
const (
	e2eContainerEnvMarker      = "UBOOT_E2E_CONTAINERENV"
	e2eContainerEnvMarkerValue = "container-process-carries-me"
)

// LH-FA-DEV-001 (lastenheft 0.32.0): the project environment values stand
// in the `containerEnv` level alone — the generated Dockerfile carries none
// of them, a bare `docker run` of the built image has them in no
// environment, and the started container carries them in the environment of
// its processes (the `containerEnv` resolution of the Dev Container CLI).
func TestE2E_ContainerEnv_ValuesReachTheContainerProcess(t *testing.T) {
	requireDocker(t)
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := dcGenerateProject(ctx, t, false, map[string]string{
		"devcontainer.containerEnv.BUILD_MODE":               "local",
		"devcontainer.containerEnv." + e2eContainerEnvMarker: e2eContainerEnvMarkerValue,
	})
	df, err := os.ReadFile(filepath.Join(dir, ".devcontainer", "Dockerfile"))
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"BUILD_MODE", e2eContainerEnvMarker} {
		if strings.Contains(string(df), name+"=") || strings.Contains(string(df), "ENV "+name) {
			t.Errorf("the Dockerfile carries %s although the value belongs to the containerEnv level alone:\n%s", name, string(df))
		}
	}
	jsonBody, err := os.ReadFile(filepath.Join(dir, ".devcontainer", "devcontainer.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`"BUILD_MODE": "local"`,
		`"` + e2eContainerEnvMarker + `": "` + e2eContainerEnvMarkerValue + `"`,
	} {
		if !strings.Contains(string(jsonBody), want) {
			t.Errorf("devcontainer.json lacks %q:\n%s", want, jsonBody)
		}
	}

	// The negative probe of the level: the image the build produces carries
	// none of the values in its environment — they are per container
	// process, resolved at start.
	image := buildProjectImage(ctx, t, dir, "uboot-e2e-containerenv")
	out, code := dockerOutput(ctx, t, "run", "--rm", image, "sh", "-c",
		`echo "IMAGE:BUILD_MODE=${BUILD_MODE:-__unset__}"; echo "IMAGE:MARKER=${`+e2eContainerEnvMarker+`:-__unset__}"`)
	if code != 0 {
		t.Fatalf("image probe exit %d:\n%s", code, out)
	}
	if !strings.Contains(out, "IMAGE:BUILD_MODE=__unset__") || !strings.Contains(out, "IMAGE:MARKER=__unset__") {
		t.Errorf("the image carries a containerEnv value in its environment:\n%s", out)
	}

	// The started container: the environment of its processes carries both
	// values (the containerEnv resolution).
	first, upCode := dcUpJSON(ctx, t, dir)
	if upCode != 0 {
		t.Fatalf("up: exit %d\nlast diagnostics: %+v", upCode, first.Diagnostics)
	}
	if id := first.Data.ContainerID; id == "" {
		t.Fatal("up left no container id in the envelope")
	} else {
		out, code = dockerOutput(ctx, t, "exec", id, "sh", "-c",
			`echo "PROC:BUILD_MODE=$BUILD_MODE"; echo "PROC:MARKER=$`+e2eContainerEnvMarker+`"`)
		if code != 0 {
			t.Fatalf("container process probe exit %d:\n%s", code, out)
		}
		for _, want := range []string{
			"PROC:BUILD_MODE=local",
			"PROC:MARKER=" + e2eContainerEnvMarkerValue,
		} {
			if !strings.Contains(out, want) {
				t.Errorf("the container process environment lacks %q:\n%s", want, out)
			}
		}
	}
}
