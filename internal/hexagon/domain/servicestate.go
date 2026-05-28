package domain

// ServiceState classifies the observable state of a service inside a
// u-boot project, per the LH-FA-ADD-005 state machine. The state is
// the outcome of inspecting two artefacts:
//
//   - the `services.<name>` entry in `u-boot.yaml` (present /
//     absent; `enabled: true|false|unset`);
//   - the `BEGIN/END U-BOOT MANAGED BLOCK: service.<name>` marker in
//     `compose.yaml` (present / absent).
//
// The application layer derives the state during `u-boot add` and
// `u-boot doctor`; the resulting action depends on the state.
type ServiceState int

const (
	// ServiceStateUnregistered means no `services.<name>` entry in
	// u-boot.yaml and no managed compose-block. `u-boot add` may
	// create both.
	ServiceStateUnregistered ServiceState = iota

	// ServiceStateActive means `services.<name>.enabled: true` AND
	// the managed compose-block is present. `u-boot add` is a no-op
	// (idempotent); it returns a nil-error response with Changed=nil.
	ServiceStateActive

	// ServiceStateDeactivated means `services.<name>.enabled: false`
	// (registered but explicitly disabled). `u-boot add` flips
	// enabled→true and re-emits the compose-block.
	ServiceStateDeactivated

	// ServiceStateEnabledUnset means `services.<name>` is present
	// but the `enabled:` key is missing — per LH-FA-ADD-005 §893
	// this counts as deactivated for add purposes and is flagged as
	// a Warn in `u-boot doctor`.
	ServiceStateEnabledUnset

	// ServiceStateInconsistentYAML means a managed compose-block is
	// present but `services.<name>` is missing from u-boot.yaml.
	// The block has no YAML anchor — likely a partial cleanup.
	// `u-boot add` must abort with ErrServiceInconsistent and a
	// repair hint.
	ServiceStateInconsistentYAML

	// ServiceStateInconsistentBlock means `services.<name>.enabled:
	// true` in u-boot.yaml but the managed compose-block is missing.
	// Deterministic recovery: `u-boot add` re-emits the compose-
	// block (no abort).
	ServiceStateInconsistentBlock
)

// String returns the canonical lowercase identifier of the state,
// used by the [Logger]-port and the CLI's diagnostic output. The
// identifiers are stable: CI dashboards and log scrapers may pin
// them.
func (s ServiceState) String() string {
	switch s {
	case ServiceStateUnregistered:
		return "unregistered"
	case ServiceStateActive:
		return "active"
	case ServiceStateDeactivated:
		return "deactivated"
	case ServiceStateEnabledUnset:
		return "enabled-unset"
	case ServiceStateInconsistentYAML:
		return "inconsistent-yaml"
	case ServiceStateInconsistentBlock:
		return "inconsistent-block"
	default:
		return "unknown"
	}
}
