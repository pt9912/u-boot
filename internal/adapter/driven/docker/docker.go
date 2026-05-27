// Package docker is the os/exec-backed implementation of the
// `port/driven.DockerProbe` interface (LH-FA-DIAG-002). It shells
// out to the host `docker` binary for read-only diagnostics: client
// version, daemon reachability, compose plugin version.
//
// Layer rules (LH-FA-ARCH-003): driven adapter — imports
// `hexagon/port/driven` + stdlib; must not import application or
// driving.
package docker

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// Probe is the production DockerProbe adapter. Construct with [New].
type Probe struct {
	// binary lets tests substitute a stub via [WithBinary]; production
	// code uses the default "docker".
	binary string
}

// Static check: Probe satisfies the DockerProbe port.
var _ driven.DockerProbe = (*Probe)(nil)

// New returns a Probe that shells out to the `docker` binary on
// `$PATH`.
func New() *Probe { return &Probe{binary: "docker"} }

// WithBinary overrides the docker binary path; intended for tests
// (integration tests under build-tag `docker` may point this at a
// container-runtime alias).
func WithBinary(path string) *Probe { return &Probe{binary: path} }

// Version returns the docker client version. Uses the templated
// output `--format '{{.Client.Version}}'`, which works even when the
// daemon is unreachable — it queries only the local client.
func (p Probe) Version(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, p.binary, "version", "--format", "{{.Client.Version}}")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker version failed: %w (output: %s)", err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

// Info pings the docker daemon. A nil error means the daemon
// answered; non-nil signals unreachability (engine not running,
// socket permission denied, ...). The actual `docker info` output is
// discarded — the doctor only needs the reachability boolean.
func (p Probe) Info(ctx context.Context) error {
	// Using `version --format {{.Server.Version}}` instead of
	// `info` because:
	//   - it forces the client to talk to the daemon (returning the
	//     server version requires a roundtrip)
	//   - it is faster than `info` (which prints config / system data)
	//   - the exit code semantics are identical (non-zero if daemon
	//     unreachable)
	cmd := exec.CommandContext(ctx, p.binary, "version", "--format", "{{.Server.Version}}")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker daemon unreachable: %w (output: %s)", err, string(out))
	}
	return nil
}

// ComposeVersion returns the bare semver of the docker compose
// plugin. Uses `docker compose version --short`, which prints just
// the version number (e.g. `"2.20.0"`). Some older compose plugin
// versions print a leading `v` (`"v2.20.0"`); the adapter strips it
// so the application layer sees a clean semver every time.
func (p Probe) ComposeVersion(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, p.binary, "compose", "version", "--short")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker compose version failed: %w (output: %s)", err, string(out))
	}
	raw := strings.TrimSpace(string(out))
	return strings.TrimPrefix(raw, "v"), nil
}
