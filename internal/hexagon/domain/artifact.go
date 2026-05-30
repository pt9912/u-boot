package domain

import (
	"errors"
	"fmt"
)

// Artifact identifies the kind of artefact `u-boot generate` produces
// (LH-FA-GEN-001). The four MVP catalogue values map 1:1 to the
// positional argument of `u-boot generate <artifact>`:
//
//   - changelog    → CHANGELOG.md
//   - readme       → README.md
//   - env-example  → .env.example
//   - devcontainer → .devcontainer/devcontainer.json + Dockerfile
//
// The zero value is invalid; use [NewArtifact] to construct from the
// CLI argument string. Unknown values must surface as exit-code 2
// (CLI validation) per LH-FA-CLI-006 — the driving-port-level wrap
// (`driving.ErrArtifactUnknown`) handles the exit-code mapping; the
// domain layer raises only [ErrInvalidArtifact].
type Artifact int

const (
	// ArtifactChangelog corresponds to `u-boot generate changelog`
	// (LH-FA-GEN-002 / LH-AK-007).
	ArtifactChangelog Artifact = iota

	// ArtifactReadme corresponds to `u-boot generate readme`
	// (LH-FA-GEN-003).
	ArtifactReadme

	// ArtifactEnvExample corresponds to `u-boot generate env-example`
	// (LH-FA-GEN-004).
	ArtifactEnvExample

	// ArtifactDevcontainer corresponds to `u-boot generate devcontainer`
	// (LH-FA-DEV-001 / LH-FA-DEV-004 / LH-FA-DEV-005).
	ArtifactDevcontainer
)

// String returns the canonical CLI-argument form ("changelog",
// "readme", "env-example", "devcontainer"). Stable: the M7 spec
// catalogue and the CLI's allowed-values message both pin these
// exact strings.
func (a Artifact) String() string {
	switch a {
	case ArtifactChangelog:
		return "changelog"
	case ArtifactReadme:
		return "readme"
	case ArtifactEnvExample:
		return "env-example"
	case ArtifactDevcontainer:
		return "devcontainer"
	default:
		return "unknown"
	}
}

// ErrInvalidArtifact signals that a string does not match any
// [Artifact] in the MVP catalogue. The wrapped message lists the
// catalogue so the CLI surfaces it verbatim — `u-boot generate
// <unknown>` then maps via `driving.ErrArtifactUnknown` to exit code
// 2 (LH-FA-CLI-006).
var ErrInvalidArtifact = errors.New("invalid artifact")

// artifactCatalogue is the source of truth for [NewArtifact] and the
// error message. Function (not package-level var) to keep the
// gochecknoglobals lint clean — same shape as supportedServices in
// the application layer.
func artifactCatalogue() []string {
	return []string{"changelog", "readme", "env-example", "devcontainer"}
}

// NewArtifact parses raw into an [Artifact]. Returns
// [ErrInvalidArtifact]-wrapped error on unknown input; the error
// message includes the catalogue list so the CLI surfaces it
// directly.
func NewArtifact(raw string) (Artifact, error) {
	switch raw {
	case "changelog":
		return ArtifactChangelog, nil
	case "readme":
		return ArtifactReadme, nil
	case "env-example":
		return ArtifactEnvExample, nil
	case "devcontainer":
		return ArtifactDevcontainer, nil
	}
	return 0, fmt.Errorf("%w: %q is not one of %v",
		ErrInvalidArtifact, raw, artifactCatalogue())
}
