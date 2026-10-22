package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ErrInvalidSandboxSetting signals a business error of the
// LH-FA-DEV-004 / LH-FA-DEV-006 / LH-FA-DEV-007 devcontainer
// settings: a value of `devcontainer.user.uid`, `devcontainer.profile`,
// `devcontainer.sandbox.nestedRuntime` or
// `devcontainer.sandbox.onUnavailable` outside its closed value set, or a
// request the sandbox profile rejects (`--sandbox` without a devcontainer,
// an `origin` URL that cannot be written into the devcontainer safely, a
// built-in feature that mounts the container-engine socket of the host).
// Lives in the domain layer so the CLI adapter can include it in the
// exit-code-10 mapping without importing the application layer.
var ErrInvalidSandboxSetting = errors.New("invalid sandbox setting")

// DevcontainerProfile is the closed value set of
// `devcontainer.profile` (LH-FA-DEV-006).
type DevcontainerProfile string

const (
	// ProfileDefault is the pre-0.3.0 devcontainer shape.
	ProfileDefault DevcontainerProfile = "default"
	// ProfileSandbox is the opt-in sandbox profile.
	ProfileSandbox DevcontainerProfile = "sandbox"
)

// NewDevcontainerProfile validates raw against the closed set.
func NewDevcontainerProfile(raw string) (DevcontainerProfile, error) {
	switch p := DevcontainerProfile(raw); p {
	case ProfileDefault, ProfileSandbox:
		return p, nil
	}
	return "", fmt.Errorf("%w: devcontainer.profile %q; allowed: default, sandbox",
		ErrInvalidSandboxSetting, raw)
}

// NestedRuntime is the closed value set of
// `devcontainer.sandbox.nestedRuntime` (LH-FA-DEV-007).
type NestedRuntime string

const (
	// NestedRuntimeNone provides no container runtime inside the
	// devcontainer (Default).
	NestedRuntimeNone NestedRuntime = "none"
	// NestedRuntimePodman provides rootless Podman inside it.
	NestedRuntimePodman NestedRuntime = "podman"
	// NestedRuntimeDocker provides an own Docker daemon inside it
	// (needs `--privileged`).
	NestedRuntimeDocker NestedRuntime = "docker"
)

// NewNestedRuntime validates raw against the closed set.
func NewNestedRuntime(raw string) (NestedRuntime, error) {
	switch r := NestedRuntime(raw); r {
	case NestedRuntimeNone, NestedRuntimePodman, NestedRuntimeDocker:
		return r, nil
	}
	return "", fmt.Errorf("%w: devcontainer.sandbox.nestedRuntime %q; allowed: none, podman, docker",
		ErrInvalidSandboxSetting, raw)
}

// OnUnavailable is the closed value set of
// `devcontainer.sandbox.onUnavailable` (LH-FA-DEV-007 degradation
// table).
type OnUnavailable string

const (
	// OnUnavailableWarn degrades along the fallback chain with a
	// visible warning (Default).
	OnUnavailableWarn OnUnavailable = "warn"
	// OnUnavailableFail turns an unavailable capability into an
	// environment problem (exit 11).
	OnUnavailableFail OnUnavailable = "fail"
)

// NewOnUnavailable validates raw against the closed set.
func NewOnUnavailable(raw string) (OnUnavailable, error) {
	switch o := OnUnavailable(raw); o {
	case OnUnavailableWarn, OnUnavailableFail:
		return o, nil
	}
	return "", fmt.Errorf("%w: devcontainer.sandbox.onUnavailable %q; allowed: warn, fail",
		ErrInvalidSandboxSetting, raw)
}

// Container-user UID bounds (LH-FA-DEV-004): 0 (root) is excluded so
// the non-root default cannot be undone through the UID knob.
const (
	// DefaultContainerUID is the UID used when the key is absent.
	DefaultContainerUID = 1000
	minContainerUID     = 1
	maxContainerUID     = 65535
)

// NewContainerUID validates the UID of the devcontainer user
// (`devcontainer.user.uid`, LH-FA-DEV-004): 1..65535.
func NewContainerUID(uid int) (int, error) {
	if uid < minContainerUID || uid > maxContainerUID {
		return 0, fmt.Errorf("%w: devcontainer.user.uid %d; allowed: %d..%d (0 = root is rejected)",
			ErrInvalidSandboxSetting, uid, minContainerUID, maxContainerUID)
	}
	return uid, nil
}

// ParseContainerUID parses and validates a textual UID.
func ParseContainerUID(raw string) (int, error) {
	uid, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return 0, fmt.Errorf("%w: devcontainer.user.uid %q is not an integer",
			ErrInvalidSandboxSetting, raw)
	}
	return NewContainerUID(uid)
}

var egressHostPattern = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)+$`)

// NewEgressHost validates one entry of
// `devcontainer.sandbox.egress.allow` (LH-FA-DEV-008): a lowercase
// DNS host name with at least one dot and no scheme, port, path or
// wildcard. An allowed name covers its subdomains (resolver rule
// `server=/<name>/`, ADR-0012), so a wildcard entry is rejected with the
// name to enter instead.
func NewEgressHost(raw string) (string, error) {
	h := strings.TrimSpace(raw)
	if len(h) <= 253 && egressHostPattern.MatchString(h) {
		return h, nil
	}
	if strings.Contains(h, "*") {
		if name := strings.TrimPrefix(h, "*."); len(name) <= 253 && egressHostPattern.MatchString(name) {
			return "", fmt.Errorf("%w: devcontainer.sandbox.egress.allow entry %q is a wildcard; an allowed name already covers its subdomains, one and two levels below it, enter %q instead",
				ErrInvalidSandboxSetting, raw, name)
		}
		return "", fmt.Errorf("%w: devcontainer.sandbox.egress.allow entry %q is a wildcard; an allowed name already covers its subdomains, one and two levels below it, enter the host name without `*`",
			ErrInvalidSandboxSetting, raw)
	}
	return "", fmt.Errorf("%w: devcontainer.sandbox.egress.allow entry %q is not a lowercase host name (no scheme, port, path or wildcard)",
		ErrInvalidSandboxSetting, raw)
}
