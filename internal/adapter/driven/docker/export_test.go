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

// NewLauncherWithNamePrefixForTest builds a [Launcher] on the `docker`
// binary on `$PATH` whose helper containers are named `<prefix><8 hex>`, so
// a test can find and remove its own containers only.
func NewLauncherWithNamePrefixForTest(prefix string) *Launcher {
	l := NewLauncher()
	l.namePrefix = prefix
	return l
}

// MaskingWriterForTest is the masking writer of the launcher streams.
type MaskingWriterForTest struct{ w *maskingWriter }

// NewMaskingWriterForTest exposes [newMaskingWriter].
func NewMaskingWriterForTest(dst io.Writer, values []string) *MaskingWriterForTest {
	return &MaskingWriterForTest{w: newMaskingWriter(dst, values)}
}

// Write forwards to the masking writer.
func (m *MaskingWriterForTest) Write(p []byte) (int, error) { return m.w.Write(p) }

// Flush forwards to the masking writer.
func (m *MaskingWriterForTest) Flush() error { return m.w.Flush() }

// HelperContainerPrefixForTest is the name prefix of every helper
// container.
const HelperContainerPrefixForTest = helperContainerPrefix

// RunArgsForTest exposes [runArgs].
func RunArgsForTest(name, image string, envNames []string, opts driven.DevcontainerUpOptions) []string {
	return runArgs(name, image, envNames, opts)
}

// ShellArgsForTest exposes [shellArgs].
func ShellArgsForTest(opts driven.DevcontainerShellOptions, envNames []string, shell string) []string {
	return shellArgs(opts, envNames, shell)
}

// ProbeArgsForTest exposes [probeArgs].
func ProbeArgsForTest(opts driven.DevcontainerShellOptions, shell string) []string {
	return probeArgs(opts, shell)
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
func HelperImageTagForTest() (string, error) { return helperImageTag(helperImageRepo) }

// HelperImageTagInRepoForTest exposes [helperImageTag] for another repository.
func HelperImageTagInRepoForTest(repo string) (string, error) { return helperImageTag(repo) }

// NewLauncherWithImageRepoForTest builds a [Launcher] on the `docker` binary
// on `$PATH` whose helper image has the tag `<repo>:<version>-<hash>` and is
// built without the layer cache, so a test can remove and rebuild its image
// from scratch without touching the product tag.
func NewLauncherWithImageRepoForTest(repo string) *Launcher {
	l := NewLauncher()
	l.imageRepo = repo
	l.buildArgs = []string{"--no-cache"}
	return l
}

// HelperCLIVersionForTest exposes [helperCLIVersion] for the embedded
// package.json.
func HelperCLIVersionForTest() (string, error) {
	files, err := helperContextFiles()
	if err != nil {
		return "", err
	}
	return helperCLIVersion(files["package.json"])
}
