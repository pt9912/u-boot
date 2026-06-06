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

	// PreviewMode encodes the --dry-run × --diff flag combination per
	// slice-v1-cli-json-dry-run-remove T0-(b) (inherited 1:1 from
	// init T0-(b)/add T0-(b) truth table — kein remove-Prefix-Alias,
	// init-Slice-T0-(c) Alias-Lebensdauer-Pflicht zwingt direktes
	// [PreviewMode]). Default zero value [PreviewNone] preserves the
	// existing production-write behaviour; non-zero values route every
	// FS access of this Remove() invocation through the recordingfs
	// adapter (Composition-Root removeFSFactory closure, T4).
	PreviewMode PreviewMode

	// SilenceConfirmer disables the [driven.Confirmer]-prompt during
	// the Remove() call (slice-v1-cli-json-dry-run-remove T0-(j), NEW
	// pattern — NOT inherited from init's SilenceProgress because
	// init swaps ProgressPort, not Confirmer). The CLI adapter sets
	// this to true when --json is set so the JSON envelope on stdout
	// isn't corrupted by interactive prompts.
	//
	// Semantics-Klarstellung: das ist KEIN Silencing (keine UX-
	// Information-Verlust-Symmetrie zu noopProgress), sondern eine
	// bewusste Behaviour-Change im JSON-Mode — der --purge-Gate
	// trifft den noopConfirmer (returnt false, nil) → der Use-Case
	// generiert ErrConfirmationRequired ohne User-Prompt. User muss
	// explizit --yes setzen um im JSON-Mode zu purgen.
	SilenceConfirmer bool
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

	// PlannedFiles is the FS-plan emitted when [RemoveServiceRequest.
	// PreviewMode] is non-zero (slice-v1-cli-json-dry-run-remove T2 /
	// inherited from init T2 / add T0-(i) / generate T2). One entry
	// per mutated path captured by the recorder, in the order the use
	// case attempted them. Empty for PreviewNone (no recorder wired)
	// and for true no-ops. Includes `delete`-Action captures for
	// RemoveAll on extraFiles (T0-(p) — remove is the first end-to-
	// end-visible delete-action producer).
	//
	// Mid-Write-Failure-Semantik (R4 Recorder-Realität,
	// recordingfs.go:139 zeichnet vor Delegieren auf): bei einem
	// underlying.WriteFile/RemoveAll-Failure enthält PlannedFiles
	// trotzdem die i Captures inkl. des fehlgeschlagenen — die T0-(i)
	// Mid-Write-Failure-AK und der T6 Pin testen exakt diese Form.
	PlannedFiles []PlannedFile

	// Changes mirrors PlannedFiles' paths with their line-count
	// summaries (LH-FA-CLI-007 §365-371). Filled only in preview
	// modes; nil for PreviewNone. Count semantics follow add T0-(g)
	// (1:1 inherited): create = CountLines(NewContent); modify = sum
	// of `+`-lines via diff.CountAdditions; delete = 0 (T0-(p)).
	// Today populated by the CLI-adapter's mapPlannedFilesToWire
	// helper, not directly by the application service.
	Changes []ChangeEntry

	// Warnings carries soft-warning Diagnostics that the CLI adapter
	// renders into the JSON envelope's `diagnostics[]` array with
	// `level: "warn"` (slice-v1-cli-json-dry-run-remove T0-(g) +
	// R7-MED-F2 + R8-MED-F2 + R9-MED-F2). Use-Case is Source-of-Truth
	// for WARN because it knows the Catalog (volumeOptional lookup,
	// T0-(g) R3-HIGH-F1 / R5-MED-F5); CLI maps via
	// `mapWarningsToDiagnostics(resp.Warnings) []diagnosticItem`.
	//
	// Today emitted only by the `--purge && !VolumesPurged && catalog
	// has named volume` path with Code=LH-FA-ADD-007. Empty/nil in
	// all other paths.
	Warnings []WarningEntry
}

// WarningEntry is a soft-warning Diagnostic emitted by a Use-Case
// (slice-v1-cli-json-dry-run-remove T2 R8-MED-F2 + R9-MED-F2).
// Generic Cluster-Vorlauf-Type for the following Folge-Slices
// (6/9 up/down recreate-warnings, 8/9 config-set value-warnings)
// to inherit without breaking Type-Change — analogous to
// PreviewMode-rename in init T0-(c).
//
// Layer-Heim: `driving`-package because Use-Cases produce
// WarningEntry values; the CLI adapter consumes them via a thin
// mapping helper (`mapWarningsToDiagnostics`).
//
// Field-Schema mirrors the CLI-Wire `diagnosticItem` form
// (`cli/jsonenvelope.go:diagnosticItem`):
//   - Code is the LH-Kennung (`LH-FA-ADD-007` for remove's
//     deferred-volumes WARN; future slices add their own).
//   - Level is "warn" today (Spec §1834 allows warn | error;
//     the field is kept for symmetry with diagnosticItem and
//     future error-level Use-Case-Diagnostics, e.g. doctor's
//     readonly diagnostics-emission pattern).
//   - Message is the user-facing text.
//   - Subject is OPTIONAL (R12-LOW-F4 proactive Cluster-Vorlauf
//     for up/down per-Service-WARN and config-set per-Key-WARN).
//     Remove leaves it empty (`""` → omitempty drops it). The
//     CLI adapter maps Subject to `diagnostics[].file` or a
//     future `diagnostics[].subject`-Field — implementation in
//     T5 (mapWarningsToDiagnostics).
type WarningEntry struct {
	Code    string `json:"code"`
	Level   string `json:"level"`
	Message string `json:"message"`
	Subject string `json:"subject,omitempty"`
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

// ErrRemoveFileSystem signals that the remove use case hit a raw
// filesystem error during write/remove/read (slice-v1-cli-json-
// dry-run-remove T2 / T0-(d) inherited from init's ErrInitFileSystem
// + generate's ErrGenerateFileSystem). T3 wraps the 8 FS-Wrap-
// Stellen in removeservice.go (Z. 235, 241, 272, 282, 286, 321, 325,
// 358) with multi-`%w`-form (Go 1.20+):
//
//	`fmt.Errorf("remove %s: %w: %w", path, ErrRemoveFileSystem,
//	            rawErr)`
//
// Switch-Order in `mapRemoveErrorToDiagnostic` (T0-(e)) MUST check
// ErrRemoveFileSystem FIRST so multi-`%w` chains that include both
// ErrRemoveFileSystem AND a fachlich sentinel route to the FS-class
// (LH-NFA-REL-003 / exit 14), not the fachlich-class. Maps to
// LH-FA-CLI-006 exit code 14 via cli's `isFilesystemError`.
//
// Note: Z. 304 (managedblock-malformed), Z. 307 (scanner default-
// branch) and Z. 330 (yaml.PatchScalar) are NOT FS-wraps — they
// carry [ErrServiceInconsistent] per T0-(d) R4-HIGH-F1
// classification fix.
var ErrRemoveFileSystem = errors.New("remove: filesystem mutation failed")

// ErrConfirmerUnavailable signals that the [driven.Confirmer]
// ConfirmRemoveVolumes call returned an I/O error (stdin EOF,
// pipe break, terminal lost) DURING the --purge confirmation gate
// (slice-v1-cli-json-dry-run-remove T2 R2-HIGH-F1). Distinct from
// [ErrConfirmationRequired] which is the User-Refusal path
// (Confirmer returns `false, nil`).
//
// runPurgeGate wraps the Confirmer's I/O error as:
//
//	`fmt.Errorf("remove: confirmer: %w: %w", ErrConfirmerUnavailable,
//	            rawErr)`
//
// (multi-`%w` analog to FS-Wraps). Switch-Order in
// `mapRemoveErrorToDiagnostic` (T0-(e)) places ErrConfirmerUnavailable
// AFTER ErrRemoveFileSystem and BEFORE the fachlich service sentinels
// — both Infrastruktur-class Sentinels (FS + Confirmer) win over
// fachlich (R3-MED-F3 + R5-MED-F2 classification).
//
// Maps to LH-FA-CLI-005A / exit 10 (Confirmation-Gate-Klasse, Spec
// §254). The User sees the same exit code as ErrConfirmationRequired
// but with a different diagnostic message — both are gate-failures
// from the same Spec-anchor.
var ErrConfirmerUnavailable = errors.New("remove: confirmer unavailable")

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
