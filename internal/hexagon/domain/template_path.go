package domain

import (
	"errors"
	"fmt"
	"path"
	"strings"
)

// ErrInvalidTemplatePath signals that a raw path string violates
// the ADR-0009 §Entscheidung Pfad-Eskalation-Vertrag. The CLI
// adapter wraps this via [driving.ErrInvalidTemplatePath] to map
// to LH-FA-CLI-006 exit code 10 (user must fix the template).
var ErrInvalidTemplatePath = errors.New("invalid template path")

// TemplatePath is a typed, validated relative path inside an
// external template's file tree (ADR-0009 §Entscheidung promises
// this validator as the slice-v1-template-init artifact). The
// constructor [NewTemplatePath] is the only producer; the zero
// value is invalid.
//
// Rejected by design (security-relevant):
//
//   - empty string,
//   - absolute paths (`/foo`, `\foo`) — would write outside the
//     project base dir,
//   - any `..` segment in the raw input — even if `path.Clean`
//     would normalize it away (`foo/../bar` → `bar`), the presence
//     of `..` signals an escape attempt or author confusion,
//   - Windows drive letters (`C:foo`).
//
// Accepted: relative paths with `.` / `//` collapsing (canonical
// form stored after `path.Clean`).
type TemplatePath struct {
	raw string
}

// NewTemplatePath parses raw and returns the validated path.
// See [TemplatePath] for the rejection rules.
//
// Review-followup F2 + F4 strengthen the rejection list:
//   - NUL bytes anywhere in raw — OS-level filesystems reject them
//     deterministically; surfacing the rejection at the domain
//     boundary gives the template author a clearer error message
//     than EINVAL from the eventual fs.WriteFile call.
//   - Backslash anywhere in raw — POSIX filesystems accept `\` as
//     a literal filename character but Windows / future filesystem
//     adapters interpret it as a path separator. `docker\..\..\etc\
//     passwd` would slip past the `/`-splitting `..`-detector and
//     escape the project tree on a Windows build. Reject up front
//     to keep the validator cross-platform-safe.
func NewTemplatePath(raw string) (TemplatePath, error) {
	if raw == "" {
		return TemplatePath{}, fmt.Errorf("%w: empty path", ErrInvalidTemplatePath)
	}
	if strings.ContainsRune(raw, 0) {
		return TemplatePath{}, fmt.Errorf("%w: contains NUL byte", ErrInvalidTemplatePath)
	}
	if strings.ContainsRune(raw, '\\') {
		return TemplatePath{}, fmt.Errorf("%w: %q contains backslash; use forward slashes only", ErrInvalidTemplatePath, raw)
	}
	if strings.HasPrefix(raw, "/") {
		return TemplatePath{}, fmt.Errorf("%w: %q is absolute", ErrInvalidTemplatePath, raw)
	}
	// Windows drive letter (`C:foo`, `D:\bar`). Two-char lookahead.
	// The `\\` prefix is no longer reachable (rejected above by the
	// backslash check) but a colon at position 1 still triggers.
	if len(raw) >= 2 && raw[1] == ':' {
		return TemplatePath{}, fmt.Errorf("%w: %q has a drive letter", ErrInvalidTemplatePath, raw)
	}
	// Check raw segments before normalisation — `path.Clean` would
	// resolve `foo/../bar` to `bar`, masking the escape attempt.
	for _, seg := range strings.Split(raw, "/") {
		if seg == ".." {
			return TemplatePath{}, fmt.Errorf("%w: %q contains a `..` segment", ErrInvalidTemplatePath, raw)
		}
	}
	// `path.Clean` collapses `.`, `//`, and trailing `/`. After the
	// `..`-segment guard above the cleaned result is safe to store.
	return TemplatePath{raw: path.Clean(raw)}, nil
}

// String returns the canonical (cleaned) form of the path. Round-
// trip: `NewTemplatePath(p.String())` returns an equal value.
func (p TemplatePath) String() string { return p.raw }
