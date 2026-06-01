package application

import (
	"strconv"
	"strings"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// resolveAddDependencies returns the [domain.ServiceName] values
// that the caller would need to register before activating an
// add-on with the given dependency declarations against the
// currently-loaded `u-boot.yaml` config (slice-v1-addons-deps T2).
//
// A dependency contributes to the missing list when BOTH:
//
//  1. `resolveScalarPath(cfg, dep.WhenPath) == dep.EqualsValue`
//     — the conditional trigger is met (the spec example
//     `services.keycloak.persistence: external-postgres`).
//  2. `cfg.Services[dep.Requires.String()]` is absent — the
//     required service was never registered. A registered-but-
//     disabled entry counts as PRESENT here; re-enabling is a
//     separate `add` invocation, not a dependency resolution
//     concern.
//
// Same-Requires duplicates are deduplicated; the returned slice
// preserves the order of first encounter so CLI error messages
// stay deterministic.
//
// Pure function — no I/O. The application service (Add() at the
// call site) provides the cfg by loading and parsing
// `u-boot.yaml`; the resolver only consumes the projection.
func resolveAddDependencies(cfg ubootYAMLConfig, deps []domain.AddOnDependency) []domain.ServiceName {
	if len(deps) == 0 {
		return nil
	}
	var missing []domain.ServiceName
	seen := make(map[string]bool)
	for _, d := range deps {
		if resolveScalarPath(cfg, d.WhenPath) != d.EqualsValue {
			continue
		}
		required := d.Requires.String()
		if _, registered := cfg.Services[required]; registered {
			continue
		}
		if seen[required] {
			continue
		}
		seen[required] = true
		missing = append(missing, d.Requires)
	}
	return missing
}

// resolveScalarPath returns the canonical string representation of
// the scalar at path inside cfg, or `""` if the path is unknown or
// the field is unset. Bool fields render as `"true"`/`"false"`;
// string fields render verbatim.
//
// Supported paths today (v0.3.0 catalogue scope):
//
//   - `project.name`
//   - `devcontainer.enabled`
//   - `services.<svcname>.enabled`
//
// Returns `""` for any other path. slice-v1-keycloak extends this
// switch with `services.<svcname>.persistence` when the Keycloak
// add-on lands; slice-v1-otel may add OTel-specific paths.
//
// Empty-string return doubles as "path not found" AND "field is
// unset" — both produce no-trigger semantics in
// [resolveAddDependencies] because `EqualsValue` must be non-
// whitespace per `domain.AddOnDependency.Validate`.
func resolveScalarPath(cfg ubootYAMLConfig, path string) string {
	parts := strings.Split(path, ".")
	switch {
	case len(parts) == 2 && parts[0] == "project" && parts[1] == "name":
		return cfg.Project.Name

	case len(parts) == 2 && parts[0] == "devcontainer" && parts[1] == "enabled":
		if cfg.Devcontainer == nil || cfg.Devcontainer.Enabled == nil {
			return ""
		}
		return strconv.FormatBool(*cfg.Devcontainer.Enabled)

	case len(parts) == 3 && parts[0] == "services" && parts[2] == "enabled":
		svc, ok := cfg.Services[parts[1]]
		if !ok || svc.Enabled == nil {
			return ""
		}
		return strconv.FormatBool(*svc.Enabled)
	}
	return ""
}

// missingServiceNamesAsStrings projects a slice of ServiceName to
// its raw string representation for embedding into error messages.
// Lives next to the resolver because both surface the same
// formatting concern for the CLI display path.
func missingServiceNamesAsStrings(missing []domain.ServiceName) []string {
	out := make([]string, len(missing))
	for i, s := range missing {
		out[i] = s.String()
	}
	return out
}
