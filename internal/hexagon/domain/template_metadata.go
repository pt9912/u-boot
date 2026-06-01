package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// templateNameRE pins the kebab-case identifier shape from
// ADR-0009 §Entscheidung (`name: micronaut`, `name: sveltekit`,
// `name: micronaut-sveltekit`). The pattern enforces single-dash
// segment separators: an alphanumeric segment, optionally followed
// by `-segment` repetitions. Single-character names are allowed
// (a single `[a-z0-9]+` segment with no separators).
//
// Rejected by design:
//   - empty string (no segment)
//   - leading or trailing dash (`-foo`, `foo-`)
//   - consecutive dashes (`my--bad`) — would violate the
//     kebab-case intent ADR-0009 examples imply.
var templateNameRE = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// ErrInvalidTemplate signals that a TemplateMetadata value failed
// the LH-FA-TPL-002 / ADR-0009 §Entscheidung minimum (Name +
// Description + Version present, Name kebab-case). The
// `external-templates` adapter wraps the per-template parse error
// with this sentinel so the catalog-load-time test can branch on
// `errors.Is(err, domain.ErrInvalidTemplate)` without importing
// the adapter.
var ErrInvalidTemplate = errors.New("invalid template metadata")

// TemplateMetadata describes one external project template per
// ADR-0009 §Entscheidung. Populated read-only by the
// `port/driven.TemplateCatalog` adapters from each `template.yaml`
// file; consumed by the application layer for `u-boot template list`
// (this slice) and `u-boot init --template <name>` (slice-v1-
// template-init, future).
//
// The struct mirrors the `template.yaml` v1 schema literally so the
// adapter's `toDomain` mapping is a flat copy; downstream code does
// not see the raw YAML shape (`apiVersion`-tagged) at all.
type TemplateMetadata struct {
	// Name is the kebab-case identifier the user types in
	// `--template <name>`. Validated by [TemplateMetadata.Validate]
	// against [templateNameRE].
	Name string

	// Description is the one-line human label rendered by
	// `u-boot template list`. Required (LH-FA-TPL-002).
	Description string

	// Version is a free-form template version string (semver-ish
	// by convention but not enforced). Required (LH-FA-TPL-002).
	Version string

	// SupportedAddOns lists the `u-boot add <svc>` services the
	// template can be combined with cleanly. Empty list is
	// allowed and means "no add-on integration claims" — not
	// "all add-ons", to keep the contract honest.
	SupportedAddOns []string

	// GeneratedFiles is the LH-FA-TPL-002 listing of paths the
	// template renders. Empty allowed for the listing slice;
	// `slice-v1-template-init` consults it as the render plan.
	GeneratedFiles []string

	// RequiredTools surfaces external CLI / runtime dependencies
	// the template assumes (e.g. `jdk:>=21`). Free-form strings;
	// `u-boot doctor` consults the list in a future slice.
	RequiredTools []string

	// Variables is the LH-FA-TPL-002 / ADR-0009 §Entscheidung
	// variable schema the user can pass via
	// `--var key=value` in `slice-v1-template-init`. For the
	// listing slice we surface the names + descriptions but do
	// not resolve them.
	Variables []TemplateVariable
}

// TemplateVariable is one entry in [TemplateMetadata.Variables].
// Mirror of the ADR-0009 §Entscheidung variable-schema fields.
type TemplateVariable struct {
	Name        string
	Description string
	Default     string
	Required    bool
}

// Validate enforces the LH-FA-TPL-002 metadata minimum and the
// kebab-case Name format. Returns an [ErrInvalidTemplate]-wrapped
// error that lists every issue, so a catalog-load surfaces multiple
// authoring problems in one pass instead of one-at-a-time.
//
// Empty SupportedAddOns / GeneratedFiles / RequiredTools / Variables
// are intentionally allowed: the listing slice does not enforce
// rendering claims, and the bootstrap `basic` template legitimately
// has an empty variable set.
func (t TemplateMetadata) Validate() error {
	var probs []string
	switch {
	case strings.TrimSpace(t.Name) == "":
		probs = append(probs, "name is required")
	case !templateNameRE.MatchString(t.Name):
		probs = append(probs, fmt.Sprintf("name %q must be kebab-case (lowercase alphanumeric + dashes)", t.Name))
	}
	if strings.TrimSpace(t.Description) == "" {
		probs = append(probs, "description is required")
	}
	if strings.TrimSpace(t.Version) == "" {
		probs = append(probs, "version is required")
	}
	if len(probs) > 0 {
		return fmt.Errorf("%w: %s", ErrInvalidTemplate, strings.Join(probs, "; "))
	}
	return nil
}
