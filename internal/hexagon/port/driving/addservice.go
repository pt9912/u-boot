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
// CLI adapter renders it as a short summary; consumers can also
// branch on [PriorState] for idempotent-detected messages ("was
// already active; no changes").
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
	//   - PriorState=Active       → State=Active: "X already active (no changes)."
	//
	// Inconsistent-state aborts never produce a response; they
	// return [ErrServiceInconsistent] instead.
	PriorState domain.ServiceState

	// State is the resulting [domain.ServiceState] after the add.
	// On a successful call this is always [domain.ServiceStateActive]
	// (no-op when PriorState was already Active; flipped to Active
	// otherwise).
	State domain.ServiceState

	// Changed lists the project-relative paths the use case mutated
	// (`compose.yaml`, `.env.example`, `u-boot.yaml`). Empty when
	// [PriorState] was [domain.ServiceStateActive] — that path runs
	// no writes.
	Changed []string
}

// ErrServiceUnsupported signals that the requested service name is
// valid syntactically (passes [domain.NewServiceName]) but is not
// in the built-in catalogue the application service knows how to
// add. MVP catalogue: only `postgres`. The CLI maps this to
// LH-FA-CLI-006 exit code 10 (validation).
var ErrServiceUnsupported = errors.New("service not supported")

// ErrServiceInconsistent signals an LH-FA-ADD-005-§896 condition:
// a managed `BEGIN/END U-BOOT MANAGED BLOCK: service.<name>` block
// is present in `compose.yaml` but the matching `services.<name>`
// entry is missing from `u-boot.yaml` — the YAML anchor has been
// removed but the orphan compose-block survived (typically a
// partial cleanup). The add use-case refuses to silently re-anchor
// because doing so could be the wrong recovery for an
// intentionally-different state. The CLI maps it to exit code 10
// with a repair hint pointing at manual cleanup.
//
// Sentinel kept here (not in the application package) so the CLI
// adapter can branch on [errors.Is] without an `application` import
// (LH-FA-ARCH-003).
var ErrServiceInconsistent = errors.New("service state inconsistent")

// ErrProjectNotInitialized signals that BaseDir contains no
// `u-boot.yaml` (or one that cannot be parsed into the expected
// schema). LH-FA-ADD-001 requires an initialized project; the use
// case refuses to invent a config. The CLI maps it to exit code 10
// with a "run u-boot init" hint.
var ErrProjectNotInitialized = errors.New("project not initialized")

// AddServiceUseCase is the driving-port for `u-boot add <service>`.
// The CLI adapter holds a reference and calls [Add] from the Cobra
// command handler.
//
// Contract:
//
//   - On success the response always has State =
//     [domain.ServiceStateActive] and a non-empty Changed slice
//     (unless PriorState was already Active).
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
