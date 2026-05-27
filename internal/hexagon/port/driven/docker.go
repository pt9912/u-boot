package driven

import "context"

// DockerProbe abstracts the read-only docker CLI operations the M4
// doctor uses to diagnose the local environment:
//
//   - `Version`: the docker client version (proves the binary is on
//     PATH; works even when the daemon is unreachable).
//   - `Info`: a roundtrip to the daemon; non-nil error means the
//     daemon is not reachable (engine down, socket permissions, etc.).
//   - `ComposeVersion`: the docker compose plugin version (proves the
//     plugin is installed and produces a parseable semver).
//
// Methods take a [context.Context] because the underlying adapter
// shells out to `docker`, which can block on network or hang on a
// stale socket.
//
// Layer rules (LH-FA-ARCH-002, LH-FA-ARCH-003): driven port; the
// production implementation lives in `internal/adapter/driven/docker`
// and is the only place `os/exec docker ...` runs.
//
// Separate from a future `DockerEngine` port (M6, `Up`/`Down`/...)
// — `DockerProbe` is read-only diagnostics, `DockerEngine` mutates
// state; splitting keeps each port narrow.
type DockerProbe interface {
	// Version returns the docker client version string (e.g.
	// `"24.0.7"`). A non-nil error signals the binary is missing or
	// failed to run.
	Version(ctx context.Context) (string, error)

	// Info performs a single roundtrip to the docker daemon (e.g.
	// `docker info` or `docker version --format '{{.Server.Version}}'`).
	// A nil error means the daemon answered; non-nil signals
	// unreachability — the doctor maps this to `docker.reachable:
	// error` with a "start docker" hint.
	Info(ctx context.Context) error

	// ComposeVersion returns the bare semver of the docker compose
	// plugin (e.g. `"2.20.0"`, without a leading `v`). A non-nil
	// error signals the plugin is missing.
	ComposeVersion(ctx context.Context) (string, error)
}
