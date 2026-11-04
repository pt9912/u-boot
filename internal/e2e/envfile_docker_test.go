//go:build docker

// End-to-end check of the environment file of the devcontainer
// (LH-FA-DEV-012, ADR-0044): a variable written to
// `.devcontainer/devcontainer.env` is set in the container that the generated
// `runArgs` start, with `docker run` and through `u-boot devcontainer up`.
// They run on the engine the suite is pointed at and say nothing about
// another engine; only the first test skips itself off Linux `amd64`.

package e2e_test

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const (
	envProbeName  = "UBOOT_ENVFILE_PROBE"
	envProbeQuote = "UBOOT_ENVFILE_QUOTED"
)

// envProbeFile is the content of the environment file of both tests: one
// plain value with a space and one value whose quotes belong to it.
const envProbeFile = envProbeName + "=hello world\n" + envProbeQuote + "=\"a b\"\n"

// F1: `docker run` with the generated `runArgs` sets the variables of the
// file; the file of the same project without the lines leaves them unset.
func TestE2E_EnvFile_VariableReachesContainer(t *testing.T) {
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		t.Skipf("measured on a Linux host (amd64) only (GOOS=%s GOARCH=%s)", runtime.GOOS, runtime.GOARCH)
	}
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := dcGenerateProject(ctx, t, false, nil)
	envPath := filepath.Join(dir, ".devcontainer", "devcontainer.env")
	runArgs := generatedRunArgs(t, dir)
	probe := func() (string, int) {
		args := append(append([]string{"run", "--rm"}, runArgs...), dcBaseImage, "sh", "-c",
			"printf '[%s][%s]' \"${"+envProbeName+"-unset}\" \"${"+envProbeQuote+"-unset}\"")
		return dockerOutput(ctx, t, args...)
	}

	// Counter-sample: the file the generator wrote holds no assignment.
	out, code := probe()
	if code != 0 || !strings.HasSuffix(out, "[unset][unset]") {
		t.Fatalf("created file: exit %d, output %q, want the variables unset", code, out)
	}
	if err := os.WriteFile(envPath, []byte(envProbeFile), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code = probe()
	if code != 0 || !strings.HasSuffix(out, "[hello world][\"a b\"]") {
		t.Errorf("exit %d, output %q, want [hello world] and the quotes kept as part of the value", code, out)
	}
}

// F2: the container that `u-boot devcontainer up` starts from the generated
// project holds the variable (SPEC-049 mirror path: the file stands in the
// helper container under the path of the host).
func TestE2E_EnvFile_VariableReachesContainerViaDevcontainerUp(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), dcTestBudget)
	defer cancel()
	dir := dcGenerateProject(ctx, t, false, nil)
	if err := os.WriteFile(filepath.Join(dir, ".devcontainer", "devcontainer.env"), []byte(envProbeFile), 0o644); err != nil {
		t.Fatal(err)
	}
	s := dcMustUp(ctx, t, dcService(), dir, false)
	out, code := dcExec(ctx, t, s.resp.ContainerID, "vscode", "printf '[%s]' \"$"+envProbeName+"\"")
	if code != 0 || out != "[hello world]" {
		t.Errorf("exec: exit %d, output %q, want [hello world]", code, out)
	}
}
