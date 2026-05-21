// Package domain holds u-boot's pure domain types and invariants.
// It must not depend on application, port, adapter, or any I/O library
// (LH-FA-ARCH-002, LH-FA-ARCH-003).
package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ProjectName is the validated name of a u-boot project. The validation
// rules come from LH-FA-INIT-006:
//
//   - only lowercase letters, digits, and `-`
//   - must start with a lowercase letter
//   - must end with a lowercase letter or digit (1-character names are
//     allowed and trivially satisfy this rule)
//   - length 1..63 characters
//
// The zero value is invalid; use [NewProjectName] to construct.
type ProjectName string

const projectNameMaxLen = 63

var projectNamePattern = regexp.MustCompile(`^[a-z]([a-z0-9-]{0,61}[a-z0-9])?$`)

// ErrInvalidProjectName signals that a string does not satisfy
// LH-FA-INIT-006. The wrapped error message explains which rule failed
// so the caller can surface it to the user.
var ErrInvalidProjectName = errors.New("invalid project name")

// NewProjectName validates the given string against LH-FA-INIT-006 and
// returns a [ProjectName] on success or [ErrInvalidProjectName]-wrapped
// error on failure. The error is sentinel-comparable via [errors.Is].
func NewProjectName(raw string) (ProjectName, error) {
	if raw == "" {
		return "", fmt.Errorf("%w: empty", ErrInvalidProjectName)
	}
	if len(raw) > projectNameMaxLen {
		return "", fmt.Errorf("%w: length %d exceeds maximum %d", ErrInvalidProjectName, len(raw), projectNameMaxLen)
	}
	if !projectNamePattern.MatchString(raw) {
		return "", fmt.Errorf("%w: %q does not match %s (LH-FA-INIT-006)", ErrInvalidProjectName, raw, projectNamePattern)
	}
	return ProjectName(raw), nil
}

// String returns the name as a plain string. It satisfies fmt.Stringer.
func (p ProjectName) String() string { return string(p) }

// NormalizeProjectName applies the deterministic normalization rules
// from LH-FA-INIT-002 to derive a candidate project name from an
// arbitrary string (typically the working-directory basename):
//
//  1. lowercase
//  2. map every character outside [a-z0-9-] to `-`
//  3. collapse consecutive `-` to a single `-`
//  4. trim leading/trailing `-` and spaces
//  5. clamp length to 63
//  6. re-trim leading/trailing `-` after the clamp
//
// The returned string is a *candidate*; the caller still has to pass it
// through [NewProjectName] for the final validation, because some
// inputs (e.g. all-digits, all-`-`) cannot be normalized into a valid
// name.
func NormalizeProjectName(raw string) string {
	lowered := strings.ToLower(raw)

	mapped := make([]rune, 0, len(lowered))
	for _, r := range lowered {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			mapped = append(mapped, r)
			continue
		}
		mapped = append(mapped, '-')
	}

	collapsed := collapseDashes(string(mapped))
	trimmed := strings.Trim(collapsed, "- ")

	if len(trimmed) > projectNameMaxLen {
		trimmed = trimmed[:projectNameMaxLen]
	}
	return strings.Trim(trimmed, "-")
}

// collapseDashes is a tiny helper kept local to the file because it is
// not part of the domain contract — it only supports the
// [NormalizeProjectName] normalization.
func collapseDashes(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	prevDash := false
	for _, r := range s {
		if r == '-' {
			if prevDash {
				continue
			}
			prevDash = true
		} else {
			prevDash = false
		}
		b.WriteRune(r)
	}
	return b.String()
}
