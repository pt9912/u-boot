// Package runtime is the filesystem-probing implementation of the
// `port/driven.RuntimeEnvironment` interface (LH-FA-ARCH-002). It
// detects whether u-boot itself is running inside a container so
// the doctor service can skip host-prerequisite checks
// (`docker.*`, `git.installed`) that would otherwise fire on the
// distroless v0.1.0 GHCR image — see
// `slice-v0.1.1-doctor-container-awareness`.
package runtime

import (
	"os"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// containerMarkers returns the ordered list of filesystem markers
// that indicate a container runtime. The list is small on purpose
// — each entry must have a single mechanical meaning, not a
// heuristic. Order matches the historical convention (docker
// first, podman/cri-o second); detection short-circuits on the
// first hit.
//
// Why not `/proc/1/cgroup`-parsing: cgroup v1/v2 ambiguity makes
// the heuristic fragile across kernels, and the markers below
// are written by every modern container runtime u-boot is
// expected to run under (Docker Desktop / Podman / docker-engine).
//
// Why a function (not a package-level `var`): the SOLID-lint
// profile (LH-QA-004 + ADR-0003) forbids package-level globals;
// the list is build-time-constant so a function is the canonical
// substitute.
func containerMarkers() []string {
	return []string{
		"/.dockerenv",        // Docker Engine, Docker Desktop.
		"/run/.containerenv", // Podman, CRI-O, buildah.
	}
}

// FileEnv is the production RuntimeEnvironment adapter; the zero
// value is usable. It probes the well-known container-runtime
// marker files via [os.Stat]. statFunc is a test seam so the
// adapter's table-driven tests can simulate marker presence
// without touching the host filesystem.
type FileEnv struct {
	// statFunc lets tests inject a fake [os.Stat]. Production
	// callers leave it nil and the zero value falls back to the
	// real os.Stat.
	statFunc func(name string) (os.FileInfo, error)
}

// Static check: FileEnv satisfies the RuntimeEnvironment port.
var _ driven.RuntimeEnvironment = (*FileEnv)(nil)

// New returns a ready-to-use FileEnv backed by [os.Stat].
func New() *FileEnv { return &FileEnv{} }

// NewWithStat is the test seam — wraps a stat function for
// table-driven tests that need to simulate marker presence
// without touching `/`. Production code uses [New].
func NewWithStat(stat func(name string) (os.FileInfo, error)) *FileEnv {
	return &FileEnv{statFunc: stat}
}

// InContainer reports true when at least one container-runtime
// marker file exists. Any [os.Stat] error other than ErrNotExist
// is treated as "no marker" — false positives would silently skip
// checks the user wants. Detection is best-effort; the spec
// explicitly says the host check should be the source of truth
// when in doubt.
func (e *FileEnv) InContainer() bool {
	stat := e.statFunc
	if stat == nil {
		stat = os.Stat
	}
	for _, marker := range containerMarkers() {
		if _, err := stat(marker); err == nil {
			return true
		}
	}
	return false
}
