package driving

import (
	"context"
	"errors"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// AddServiceRequest is the input for [AddServiceUseCase.Add]. It is
// the application-layer expression of `u-boot add <service>` per
// LH-FA-ADD-001 / LH-FA-ADD-002; the CLI adapter translates the
// positional service-name argument into [domain.ServiceName].
//
// MVP-shape — kept minimal to mirror M5's `postgres`-only scope:
// add-on-specific options (Keycloak's `--persistence`, OTel's
// `--exporter`, ...) are out of scope until LH-FA-ADD-003/-004
// (V1) land. `--with-deps` (LH-FA-ADD-006) is also V1.
type AddServiceRequest struct {
	// BaseDir is the absolute path of the initialized u-boot project
	// the service is added to. Mandatory; the CLI adapter defaults it
	// to the current working directory (mirroring `u-boot init`).
	BaseDir string

	// ServiceName is the validated identifier of the service to add
	// (`postgres` in MVP; the application service rejects names that
	// are not in its built-in catalogue with
	// [ErrServiceUnsupported]).
	ServiceName domain.ServiceName
}

// AddServiceResponse is the output of [AddServiceUseCase.Add]. The
// CLI adapter renders it as a short summary. Consumers should use
// Changed, not PriorState alone, to detect a no-op: an already-active
// service may still repair missing service artefacts.
type AddServiceResponse struct {
	// ServiceName echoes the name that was processed — useful for
	// callers that batch invocations.
	ServiceName domain.ServiceName

	// PriorState is the [domain.ServiceState] observed before the
	// add ran. Together with [State] it lets the CLI render a
	// meaningful transition message:
	//
	//   - PriorState=Unregistered → State=Active: "Added X."
	//   - PriorState=Deactivated  → State=Active: "Reactivated X."
	//   - PriorState=Active → State=Active, Changed=nil:
	//     "X already active (no changes)."
	//   - PriorState=Active → State=Active, Changed!=nil:
	//     "Repaired X artefacts."
	//
	// Inconsistent-state aborts never produce a response; they
	// return [ErrServiceInconsistent] instead.
	PriorState domain.ServiceState

	// State is the resulting [domain.ServiceState] after the add. On a
	// successful call this is always [domain.ServiceStateActive].
	State domain.ServiceState

	// Changed lists the project-relative paths the use case mutated
	// (`compose.yaml`, `.env.example`, `u-boot.yaml`). Empty means a
	// true no-op: the service was already active and all service
	// artefacts were present. PriorState may still be
	// [domain.ServiceStateActive] with non-empty Changed when Add
	// repairs missing PostgreSQL artefacts such as the volume or env
	// managed block.
	Changed []string
}

// All Add sentinels below live in the `driving` package (not in
// `application`) so the CLI adapter can branch on them via
// [errors.Is] without importing `application` — the LH-FA-ARCH-003
// depguard rule forbids that cross-layer import. The CLI maps each
// to LH-FA-CLI-006 exit code 10 (validation).

// ErrServiceUnsupported signals that the requested service name is
// valid syntactically (passes [domain.NewServiceName]) but is not
// in the built-in catalogue the application service knows how to
// add. MVP catalogue: only `postgres`.
var ErrServiceUnsupported = errors.New("service not supported")

// ErrServiceInconsistent signals an LH-FA-ADD-005-§895 condition:
// a managed `BEGIN/END U-BOOT MANAGED BLOCK: service.<name>` block
// is present in `compose.yaml` but the matching `services.<name>`
// entry is missing from `u-boot.yaml` — the YAML anchor has been
// removed but the orphan compose-block survived (typically a
// partial cleanup). The add use-case refuses to silently re-anchor
// because doing so could be the wrong recovery for an
// intentionally-different state. The CLI surfaces it with a repair
// hint pointing at manual cleanup.
var ErrServiceInconsistent = errors.New("service state inconsistent")

// ErrProjectNotInitialized signals that BaseDir contains no
// `u-boot.yaml` (or one that cannot be parsed into the expected
// schema). LH-FA-ADD-001 requires an initialized project; the use
// case refuses to invent a config. The CLI surfaces it with a
// "run u-boot init" hint.
var ErrProjectNotInitialized = errors.New("project not initialized")

// AddServiceUseCase is the driving-port for `u-boot add <service>`.
// The CLI adapter holds a reference and calls [Add] from the Cobra
// command handler.
//
// Contract:
//
//   - On success State is [domain.ServiceStateActive]. Changed is
//     empty only for a true no-op. It is non-empty for state
//     transitions (Unregistered/Deactivated/inconsistent-block →
//     Active) and for Active → Active artefact repairs.
//   - On failure the response is the zero value and the error wraps
//     one of the sentinels above (or [domain.ErrInvalidServiceName]
//     for a syntactically invalid name).
//
// Idempotence guarantee: calling [Add] twice with the same request
// is safe. The second call returns PriorState=Active, State=Active,
// Changed=nil, error=nil.
type AddServiceUseCase interface {
	Add(ctx context.Context, req AddServiceRequest) (AddServiceResponse, error)
}
