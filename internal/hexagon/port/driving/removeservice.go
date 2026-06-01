package driving

import (
	"context"
	"errors"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// RemoveServiceRequest is the input for [RemoveServiceUseCase.Remove].
// It is the application-layer expression of `u-boot remove <service>`
// per LH-FA-ADD-007 (V1) — the inverse of [AddServiceRequest], with
// the additional `--purge`-destructive opt-in for volume removal.
//
// V1-shape — kept symmetric to [AddServiceRequest] so the M5 state-
// machine code paths can be mirrored in the application service.
// Add-on-specific cleanup hooks (Keycloak's realm-export, OTel's
// collector-config) are out of scope for this slice; they land
// when the respective add-on slices need them.
type RemoveServiceRequest struct {
	// BaseDir is the absolute path of the initialized u-boot project
	// the service is removed from. Mandatory; the CLI adapter
	// defaults it to the current working directory.
	BaseDir string

	// ServiceName is the validated identifier of the service to
	// remove. The application service rejects names that are not in
	// its built-in catalogue with [ErrServiceUnsupported] — same
	// catalogue as [AddServiceRequest], mirrored on purpose.
	ServiceName domain.ServiceName

	// Purge enables the destructive volume-removal path
	// (LH-FA-ADD-007 §"Volumes nur auf explizite Anforderung"). When
	// false (default), the service's named volumes stay on disk
	// after the compose- and env-block removal — data survives the
	// remove. When true, the LH-FA-CLI-005A §254 confirmation gate
	// fires (mediated by [Yes] / [NoInteractive] below) before the
	// destructive step. [VolumesPurged] in the response reflects
	// whether the purge actually ran.
	Purge bool

	// Yes is the persistent root flag value (LH-FA-CLI-005A §237);
	// when true together with [Purge], the confirmation prompt is
	// skipped and the volume removal proceeds. CLI parses the
	// `--yes` PersistentFlag and the request constructor copies it.
	Yes bool

	// NoInteractive is the persistent root flag value; when true
	// together with [Purge] and [Yes]=false, the use case returns
	// [ErrConfirmationRequired] before any side effect. Mirrors
	// the `down --volumes` gate from M6.
	NoInteractive bool
}

// RemoveServiceResponse is the output of [RemoveServiceUseCase.Remove].
// The CLI adapter renders it as a short summary, using PriorState +
// State + Changed to choose the right phrasing.
type RemoveServiceResponse struct {
	// ServiceName echoes the name that was processed.
	ServiceName domain.ServiceName

	// PriorState is the [domain.ServiceState] observed before the
	// remove ran. Drives the CLI message:
	//
	//   - PriorState=Active → State=Deactivated: "Removed X."
	//   - PriorState=Deactivated → State=Deactivated, Changed=nil:
	//     "X is already disabled; no changes."
	//   - PriorState=EnabledUnset → State=Deactivated:
	//     "Normalised X (enabled key was missing)."
	PriorState domain.ServiceState

	// State is the resulting [domain.ServiceState] after the remove.
	// On a successful Active-or-EnabledUnset transition this is
	// [domain.ServiceStateDeactivated]; on the already-Deactivated
	// idempotent path it is unchanged.
	State domain.ServiceState

	// Changed lists the project-relative paths the use case
	// mutated. Empty signals a true no-op (already-disabled).
	// Non-empty entries today: `compose.yaml` (managed-block
	// removed), `.env.example` (managed-block removed),
	// `u-boot.yaml` (enabled flipped to false).
	Changed []string

	// VolumesPurged is true when [RemoveServiceRequest.Purge] was
	// set AND the confirmation gate passed AND the volume removal
	// succeeded. False in every other case, including the gate-
	// refused and no-volume-known paths.
	VolumesPurged bool
}

// All Remove sentinels below live in the `driving` package so the
// CLI adapter can branch via [errors.Is] without importing
// `application` (LH-FA-ARCH-003 depguard rule). All four map to
// LH-FA-CLI-006 exit code 10 (validation) via the existing
// `isValidationError` classifier — except [ErrConfirmationRequired]
// which is already wired for the M6 `down --volumes` flow.
//
// Sentinels reused from the add / M6 flows:
//
//   - [ErrServiceUnsupported]     → unknown service name
//   - [ErrServiceInconsistent]    → managed-block orphan
//   - [ErrProjectNotInitialized]  → no u-boot.yaml
//   - [ErrConfirmationRequired]   → `--purge` non-interactive without `--yes`

// ErrServiceUnregistered signals that the requested service has
// never been added to the project — there is no
// `services.<name>` entry in `u-boot.yaml` and no managed compose-
// block. Idempotent semantics live one state up: an already-
// disabled service produces a no-op success response, not this
// error. Maps to LH-FA-CLI-006 exit code 10.
//
// Distinct from [ErrServiceUnsupported]: that one means "u-boot
// has no catalogue entry for this name"; this one means "the
// catalogue knows about it but the project does not have it".
var ErrServiceUnregistered = errors.New("service not registered")

// RemoveServiceUseCase is the driving-port for `u-boot remove
// <service>` (LH-FA-ADD-007).
//
// Contract:
//
//   - On success the response carries PriorState and State.
//     Changed is empty for the idempotent already-disabled path
//     and non-empty for state-transitioning calls (Active or
//     EnabledUnset → Deactivated). VolumesPurged reflects the
//     destructive purge step.
//   - On failure the response is the zero value and the error
//     wraps one of the documented sentinels (or
//     [domain.ErrInvalidServiceName] for a syntactically invalid
//     name).
//
// Idempotence guarantee: calling [Remove] twice with the same
// request is safe. The second call returns
// PriorState=Deactivated, State=Deactivated, Changed=nil,
// error=nil.
type RemoveServiceUseCase interface {
	Remove(ctx context.Context, req RemoveServiceRequest) (RemoveServiceResponse, error)
}
