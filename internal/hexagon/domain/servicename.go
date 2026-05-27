package domain

import (
	"errors"
	"fmt"
	"regexp"
)

// ServiceName is the validated identifier of a u-boot-managed service
// add-on (e.g. `postgres`, `keycloak`, `otel`). Validation follows the
// same lowercase-letter-start / dash-or-digit-tail pattern as
// [ProjectName] but with a tighter length cap appropriate for
// Compose service keys.
//
//   - only lowercase letters, digits, and `-`
//   - must start with a lowercase letter
//   - must end with a lowercase letter or digit (1-character names
//     are allowed and trivially satisfy this rule)
//   - length 1..32 characters (Compose service keys, k8s labels,
//     env-var prefixes all stay readable within 32 chars).
//
// The zero value is invalid; use [NewServiceName] to construct.
type ServiceName string

const serviceNameMaxLen = 32

var serviceNamePattern = regexp.MustCompile(`^[a-z]([a-z0-9-]{0,30}[a-z0-9])?$`)

// ErrInvalidServiceName signals that a string does not satisfy the
// ServiceName rules. The wrapped error message explains which rule
// failed so the CLI can surface it. Sentinel-comparable via
// [errors.Is].
var ErrInvalidServiceName = errors.New("invalid service name")

// NewServiceName validates the given string and returns a
// [ServiceName] on success or [ErrInvalidServiceName]-wrapped error
// on failure.
func NewServiceName(raw string) (ServiceName, error) {
	if raw == "" {
		return "", fmt.Errorf("%w: empty", ErrInvalidServiceName)
	}
	if len(raw) > serviceNameMaxLen {
		return "", fmt.Errorf("%w: length %d exceeds maximum %d",
			ErrInvalidServiceName, len(raw), serviceNameMaxLen)
	}
	if !serviceNamePattern.MatchString(raw) {
		return "", fmt.Errorf("%w: %q does not match %s",
			ErrInvalidServiceName, raw, serviceNamePattern)
	}
	return ServiceName(raw), nil
}

// String returns the name as a plain string; satisfies fmt.Stringer.
func (s ServiceName) String() string { return string(s) }
