package domain

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrInvalidSandboxSetting signals that a value for one of the
// LH-FA-DEV-004 / LH-FA-DEV-006 / LH-FA-DEV-007 devcontainer
// settings (`devcontainer.user.uid`, `devcontainer.profile`,
// `devcontainer.sandbox.nestedRuntime`,
// `devcontainer.sandbox.onUnavailable`) is outside its closed value
// set. Lives in the domain layer so the CLI adapter can include it in
// the exit-code-10 mapping without importing the application layer.
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
)

// NewNestedRuntime validates raw against the closed set.
func NewNestedRuntime(raw string) (NestedRuntime, error) {
	switch r := NestedRuntime(raw); r {
	case NestedRuntimeNone, NestedRuntimePodman:
		return r, nil
	}
	return "", fmt.Errorf("%w: devcontainer.sandbox.nestedRuntime %q; allowed: none, podman",
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
