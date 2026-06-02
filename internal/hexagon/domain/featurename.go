package domain

import (
	"errors"
	"fmt"
	"regexp"
)

// FeatureName is the validated identifier of a u-boot devcontainer
// feature catalogue entry (e.g. `node`, `java`, `go`). Validation
// rules are identical to [ServiceName]:
//
//   - only lowercase letters, digits, and `-`
//   - must start with a lowercase letter
//   - must end with a lowercase letter or digit
//   - length 1..32 characters.
//
// The zero value is invalid; use [NewFeatureName] to construct.
//
// FeatureName is name-only (no version slot). Version pins live in
// the catalogue (`featureCatalogueEntry.defaultVersion`) and in the
// per-feature YAML override (`ubootYAMLDevcontainerFeature.Version`).
// See slice-v1-devcontainer-features §T0-Outcomes (c) for the
// rationale.
type FeatureName string

const featureNameMaxLen = 32

var featureNamePattern = regexp.MustCompile(`^[a-z]([a-z0-9-]{0,30}[a-z0-9])?$`)

// ErrInvalidFeatureName signals that a string does not satisfy the
// FeatureName rules. The wrapped error message explains which rule
// failed so the CLI can surface it. Sentinel-comparable via
// [errors.Is].
var ErrInvalidFeatureName = errors.New("invalid feature name")

// ErrInvalidFeatureSource signals that a string in the
// `devcontainer.featureSources.allow` list (or a
// `devcontainer.features.<name>.source` override) does not satisfy
// the LH-FA-DEV-003 source-format rules (non-empty, scheme in
// http/https/oci, non-empty host). Lives in the domain layer so the
// CLI adapter can include it in the exit-code-10 mapping
// ([cli.isValidationError]) without depending on the application
// layer (depguard `adapter-no-application`). The format-checking
// helper itself (`application.validateFeatureSource`) lives in the
// application layer because the parser intentionally avoids
// `net/url` per the `application-no-net` depguard rule
// (LH-FA-ARCH-003).
var ErrInvalidFeatureSource = errors.New("invalid feature source")

// NewFeatureName validates the given string and returns a
// [FeatureName] on success or [ErrInvalidFeatureName]-wrapped error
// on failure.
func NewFeatureName(raw string) (FeatureName, error) {
	if raw == "" {
		return "", fmt.Errorf("%w: empty", ErrInvalidFeatureName)
	}
	if len(raw) > featureNameMaxLen {
		return "", fmt.Errorf("%w: length %d exceeds maximum %d",
			ErrInvalidFeatureName, len(raw), featureNameMaxLen)
	}
	if !featureNamePattern.MatchString(raw) {
		return "", fmt.Errorf("%w: %q does not match %s",
			ErrInvalidFeatureName, raw, featureNamePattern)
	}
	return FeatureName(raw), nil
}

// String returns the name as a plain string; satisfies fmt.Stringer.
func (n FeatureName) String() string { return string(n) }
