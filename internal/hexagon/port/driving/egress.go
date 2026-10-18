package driving

import (
	"context"
	"errors"
)

// Origin names of an egress host (LH-FA-DEV-008). The spelling is part of
// the `--json` form of `u-boot egress list`. A feature origin is
// [EgressSourceFeaturePrefix] followed by the catalogue key.
const (
	// EgressSourceBase is the default allowlist that applies to every
	// restricted container.
	EgressSourceBase = "base"
	// EgressSourceRuntime is the registry set of `nestedRuntime: podman`
	// or `docker`.
	EgressSourceRuntime = "runtime"
	// EgressSourceDockerAPT is the Docker APT repository of
	// `nestedRuntime: docker`.
	EgressSourceDockerAPT = "docker-apt"
	// EgressSourceFeaturePrefix precedes the key of an enabled
	// devcontainer feature (`feature:node`).
	EgressSourceFeaturePrefix = "feature:"
	// EgressSourceClone is the host of the clone source
	// (`devcontainer.sandbox.repository` or the `origin` remote).
	EgressSourceClone = "clone"
	// EgressSourceUser is the user list `devcontainer.sandbox.egress.allow`.
	EgressSourceUser = "user"
)

// EgressHostEntry is one host of the allowlist with every origin that
// contributes it (a host can have several, e.g. `ghcr.io`).
type EgressHostEntry struct {
	Host    string
	Sources []string
}

// EgressListRequest is the input for [EgressListUseCase.List].
type EgressListRequest struct {
	// BaseDir is the absolute path of the u-boot project.
	BaseDir string
}

// EgressListResponse is the output of [EgressListUseCase.List].
//
// Hosts is the list that `u-boot generate devcontainer` bakes into the
// start script (sorted by host, without duplicates) when the restriction
// is active, and the list that would apply on activation otherwise.
type EgressListResponse struct {
	// SandboxProfile reports `devcontainer.profile: sandbox`.
	SandboxProfile bool
	// EgressEnabled reports `devcontainer.sandbox.egress.enabled: true`.
	EgressEnabled bool
	// CloneSource reports whether a clone source (`sandbox.repository` or
	// the `origin` remote of `.git/config`) contributes a host.
	CloneSource bool
	// Hosts is sorted by host; Sources of each entry are in a fixed order.
	Hosts []EgressHostEntry
}

// Active reports whether the restriction applies to the configuration:
// the sandbox profile and `egress.enabled` must both hold.
func (r EgressListResponse) Active() bool {
	return r.SandboxProfile && r.EgressEnabled
}

// ErrEgressFileSystem wraps a read failure of `u-boot.yaml` in the egress
// listing. It maps to LH-FA-CLI-006 exit code 14 (technical failure).
var ErrEgressFileSystem = errors.New("egress: filesystem error")

// EgressListUseCase implements `u-boot egress list` (LH-FA-DEV-008).
//
// Contract:
//
//   - Reads `u-boot.yaml` and the Git configuration of the project only;
//     needs no network, no Docker and no running container; writes nothing.
//   - Hosts is computed by the function that computes the allowlist of the
//     generated start script, so both lists are equal for one configuration.
//   - A missing `u-boot.yaml` wraps [ErrProjectNotInitialized], an invalid
//     one [ErrConfigSchemaInvalid], an unreadable one [ErrEgressFileSystem].
//     An `origin` URL that `generate devcontainer --sandbox` refuses is
//     refused the same way (a domain error, exit 10).
//   - Without a remote and without `sandbox.repository` the list has no
//     clone host and CloneSource is false.
type EgressListUseCase interface {
	List(ctx context.Context, req EgressListRequest) (EgressListResponse, error)
}
