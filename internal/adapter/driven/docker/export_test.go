package docker

import (
	"context"
	"io"
	"time"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// Test-only bridges for the package-private helpers in engine.go and
// devcontainer_launcher.go.
// Convention: keep production internals unexported, expose narrow
// surfaces here so the `_test` test package can exercise them
// without polluting the public API.

// ParseComposePsOutputForTest exposes [parseComposePsOutput] for
// unit testing the NDJSON/array parser without going through the
// full Engine.ComposePs path (which needs a real docker binary).
func ParseComposePsOutputForTest(raw []byte) ([]driven.ComposeService, error) {
	return parseComposePsOutput(raw)
}

// ProgressSinkOrDiscardForTest exposes [progressSinkOrDiscard] so
// the nil-tolerance convention can be unit-tested.
func ProgressSinkOrDiscardForTest(sink io.Writer) io.Writer {
	return progressSinkOrDiscard(sink)
}

// WrapComposeRunErrorForTest exposes [wrapComposeRunError] for the
// slice-v1-logs Review-Followup F3 unit-test of SIGINT-Pass-Through
// Schicht 1. The helper was extracted from ComposeLogs so the
// "ctx.Err() unverdeckt"-contract can be pinned without a real
// subprocess.
func WrapComposeRunErrorForTest(ctx context.Context, runErr error, kind string) error {
	return wrapComposeRunError(ctx, runErr, kind)
}

// Process is the exported alias of the inspector's command seam, so the
// `_test` package can supply a fake runner.
type Process = process

// NewInspectorWithRunnerForTest builds an [Inspector] whose commands run
// through start instead of os/exec, with an injectable environment.
func NewInspectorWithRunnerForTest(
	start func(ctx context.Context, args ...string) (Process, error),
	getenv func(string) string,
) *Inspector {
	return &Inspector{binary: "docker", getenv: getenv, start: start}
}

// NewLauncherForTest builds a [Launcher] that runs the given docker binary,
// reads the environment through lookupEnv and waits abortGrace for a
// cancelled helper container before it kills `docker run`.
func NewLauncherForTest(binary string, lookupEnv func(string) (string, bool), abortGrace time.Duration) *Launcher {
	l := WithLauncherBinary(binary)
	l.lookupEnv = lookupEnv
	l.abortGrace = abortGrace
	return l
}

// HelperContainerPrefixForTest is the name prefix of every helper
// container.
const HelperContainerPrefixForTest = helperContainerPrefix

// RunArgsForTest exposes [runArgs].
func RunArgsForTest(name, image string, envNames []string, opts driven.DevcontainerUpOptions) []string {
	return runArgs(name, image, envNames, opts)
}

// ParseUpResultForTest exposes [parseUpResult].
func ParseUpResultForTest(stdout string) (driven.DevcontainerUpResult, error) {
	return parseUpResult(stdout)
}

// HelperContextFilesForTest exposes [helperContextFiles].
func HelperContextFilesForTest() (map[string][]byte, error) { return helperContextFiles() }

// HelperContextNamesForTest exposes [helperContextNames].
func HelperContextNamesForTest() []string { return helperContextNames() }

// HelperContextTarForTest exposes [helperContextTar].
func HelperContextTarForTest() ([]byte, error) { return helperContextTar() }

// HelperImageTagForTest exposes [helperImageTag].
func HelperImageTagForTest() (string, error) { return helperImageTag() }

// HelperCLIVersionForTest exposes [helperCLIVersion] for the embedded
// package.json.
func HelperCLIVersionForTest() (string, error) {
	files, err := helperContextFiles()
	if err != nil {
		return "", err
	}
	return helperCLIVersion(files["package.json"])
}
