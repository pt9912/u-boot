package domain

import (
	"errors"
	"fmt"
	"strings"
)

// ErrInvalidAddOnDependency signals that an [AddOnDependency]
// declaration is malformed: missing Requires, WhenPath, or
// EqualsValue. The application-layer catalogue source-of-truth
// (`dependenciesFor` in addservice.go) is the only producer; tests
// branch on this via `errors.Is(err, domain.ErrInvalidAddOnDependency)`.
var ErrInvalidAddOnDependency = errors.New("invalid add-on dependency declaration")

// AddOnDependency models a conditional, YAML-path-triggered
// dependency between two add-ons per `LH-FA-ADD-006` (V1). The
// declaration is evaluated by the application-layer resolver
// (`resolveAddDependencies` — slice-v1-addons-deps T2) against
// the current `u-boot.yaml` config: when the scalar at WhenPath
// equals EqualsValue, the service named by Requires is treated
// as required for the add-on the declaration is attached to.
//
// Example (spec): Keycloak with
// `services.keycloak.persistence: external-postgres` requires
// postgres:
//
//	AddOnDependency{
//	    Requires:    ServiceName("postgres"),
//	    WhenPath:    "services.keycloak.persistence",
//	    EqualsValue: "external-postgres",
//	}
//
// Today's catalogue (`postgres` only) declares no dependencies;
// this type ships ahead of its first user (slice-v1-keycloak) so
// the dependency-resolution mechanism can be developed and
// reviewed in isolation.
type AddOnDependency struct {
	// Requires is the service this add-on conditionally depends on
	// when the WhenPath/EqualsValue condition is met. Must be a
	// validated, non-empty [ServiceName].
	Requires ServiceName

	// WhenPath is the dotted YAML path to a scalar field in
	// `u-boot.yaml` that triggers the dependency, e.g.
	// `services.keycloak.persistence`. The application-layer
	// resolver navigates this path against the current config.
	WhenPath string

	// EqualsValue is the literal scalar value at WhenPath that
	// triggers the dependency. Compared as a string after YAML-
	// canonical coercion (bool, int, etc. all become their
	// canonical string form). An empty trigger value would never
	// fire because a config-absent path coerces to `""`, so
	// EqualsValue must be non-empty.
	EqualsValue string
}

// Validate enforces the LH-FA-ADD-006 minimum: Requires must be a
// non-empty [ServiceName], WhenPath and EqualsValue both non-
// whitespace strings. Returns an [ErrInvalidAddOnDependency]-wrapped
// error that lists every issue, so a catalogue-load surfaces
// multiple authoring problems in one pass (mirror of
// [TemplateMetadata.Validate]).
func (d AddOnDependency) Validate() error {
	var probs []string
	if strings.TrimSpace(d.Requires.String()) == "" {
		probs = append(probs, "Requires is required")
	}
	if strings.TrimSpace(d.WhenPath) == "" {
		probs = append(probs, "WhenPath is required")
	}
	if strings.TrimSpace(d.EqualsValue) == "" {
		probs = append(probs, "EqualsValue is required")
	}
	if len(probs) > 0 {
		return fmt.Errorf("%w: %s", ErrInvalidAddOnDependency, strings.Join(probs, "; "))
	}
	return nil
}
