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

	// PreviewMode selects the FS-mutation regime per slice-v1-cli-json-
	// dry-run-add T0-(e) Option 4. The CLI maps --dry-run/--diff flag
	// combinations to one of three modes; the Composition-Root's
	// fsFactory routes FS calls through the recording adapter
	// accordingly:
	//
	//   PreviewNone       -- direct production FS (no capture)
	//   PreviewDryRun     -- capture only; production FS untouched
	//   PreviewAndApply   -- capture AND write through (LH-FA-CLI-008
	//                        Preview-and-Apply mode)
	//
	// Default zero value PreviewNone preserves backward compatibility
	// with the existing non-JSON code path; callers that don't set it
	// see today's production-write behaviour.
	PreviewMode AddPreviewMode

	// ServiceName is the validated identifier of the service to add
	// (`postgres` in MVP; the application service rejects names that
	// are not in its built-in catalogue with
	// [ErrServiceUnsupported]).
	ServiceName domain.ServiceName

	// WithDeps, Yes, NoInteractive drive the LH-FA-ADD-006 four-mode
	// dispatch when the requested add-on declares dependencies that
	// are not yet registered in u-boot.yaml:
	//
	//   --with-deps (WithDeps=true): auto-install missing deps
	//   without prompting. Propagates recursively so transitive deps
	//   inherit the flag.
	//
	//   --yes (Yes=true): pre-confirm any interactive prompt. Same
	//   effect as --with-deps for missing deps; also auto-confirms
	//   future destructive prompts.
	//
	//   --no-interactive (NoInteractive=true) without --yes /
	//   --with-deps: refuse to prompt and return
	//   [ErrDependenciesRequired] (LH-FA-CLI-006 exit code 10).
	//
	//   default (all zero): prompt via [driven.Confirmer.
	//   ConfirmAddDependency]; user "no" or EOF also returns
	//   [ErrDependenciesRequired].
	//
	// Add() carries these flags through the recursive sub-call so a
	// single `--with-deps` at the top level installs the whole chain.
	WithDeps      bool
	Yes           bool
	NoInteractive bool
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

	// PlannedFiles is the FS-plan emitted when [AddServiceRequest.
	// PreviewMode] is non-zero (slice-v1-cli-json-dry-run-add T0-(i)).
	// One entry per mutated path captured by the recorder, in the
	// order the use case attempted them. Empty for PreviewNone (no
	// recorder wired) and for true no-ops. Carries NewContent and
	// OldContent for the CLI-adapter diff renderer; these two fields
	// stay out of the JSON wire-format via `json:"-"` (Spec §326 has
	// no place for raw bytes).
	PlannedFiles []PlannedFile

	// Changes mirrors PlannedFiles' paths with their line-count
	// summaries (LH-FA-CLI-007 §365-371). Filled only in preview
	// modes; nil for PreviewNone. Count semantics follow T0-(g):
	// newLines/totalLines.
	Changes []ChangeEntry
}

// PreviewMode encodes the four flag combinations for
// `u-boot <modifying-subcommand>` (--dry-run × --diff) per
// slice-v1-cli-json-dry-run-add T0-(b) truth table:
//
//	flags                    | PreviewMode      | production write?
//	-------------------------+------------------+------------------
//	(neither)                | PreviewNone      | yes
//	--dry-run                | PreviewDryRun    | no
//	--diff                   | PreviewAndApply  | yes
//	--dry-run --diff         | PreviewDryRun    | no
//
// The CLI adapter computes the mode from the parsed flags and writes
// it into [AddServiceRequest.PreviewMode] (or [InitProjectRequest.
// PreviewMode] in the init slice) before calling the use case.
// Composition-Root reads the mode in its fsFactory closure to pick
// between production FS and the RecordingFileSystem variants.
//
// Originally named AddPreviewMode; slice-v1-cli-json-dry-run-init
// T0-(c) renamed to PreviewMode because the type is consumed by
// every modifying subcommand (add, init, generate, remove,
// config set). AddPreviewMode remains as a type-alias for backward
// compatibility — see below.
type PreviewMode int

// AddPreviewMode is a backward-compat type-alias for [PreviewMode]
// (slice-v1-cli-json-dry-run-init T0-(c) Carveout). The `=` syntax
// makes them the IDENTICAL type, so existing call-sites that say
// `driving.AddPreviewMode` (and the matching function-types) stay
// assignable to the renamed canonical form without source edits.
// Carveout removal owner: slice-v1-cli-cleanup-add-preview-mode-
// alias (T8 of init-slice creates the open/-stub).
type AddPreviewMode = PreviewMode

const (
	// PreviewNone selects the direct production FS path. Default zero
	// value; today's non-JSON code path keeps emitting this and stays
	// unchanged.
	PreviewNone PreviewMode = iota

	// PreviewDryRun captures every mutation in the recorder without
	// touching the production FS. Used for `--dry-run` (with or
	// without `--diff`).
	PreviewDryRun

	// PreviewAndApply captures every mutation AND writes it through
	// to the production FS. Used for `--diff` without `--dry-run`
	// (LH-FA-CLI-008 §465-470 Preview-and-Apply).
	PreviewAndApply
)

// PlannedFile is the wire-shape of one FS mutation in the LH-FA-CLI-007
// §326 voll-schema response. The CLI adapter consumes it for both the
// JSON envelope's `plannedFiles[]` and the human/JSON unified diff.
//
// NewContent and OldContent carry the raw file bytes the recorder
// captured for the CLI-adapter diff renderer. They are excluded from
// the JSON wire-form via `json:"-"` — Spec §326 has no field for
// raw bytes and embedding them would be base64-drift. Diff hunks
// rendered from these bytes land in Hunks below.
type PlannedFile struct {
	Path       string `json:"path"`
	Action     string `json:"action"` // "create" | "modify" | "delete"
	NewContent []byte `json:"-"`
	OldContent []byte `json:"-"`
	Hunks      []Hunk `json:"hunks,omitempty"`
}

// ChangeEntry is the wire-shape of one line-count summary entry. Spec
// §365-371 requires count ≥ 0. Semantics (slice-v1-cli-json-dry-run-add
// T0-(g)):
//
//   - action "create" → total lines in the new file
//   - action "modify" → sum of hunk.NewLines across all hunks
//   - action "delete" → 0
type ChangeEntry struct {
	Path  string `json:"path"`
	Count int    `json:"count"`
}

// Hunk is the wire-shape of one diff hunk in
// [PlannedFile.Hunks] (LH-FA-CLI-008 §477-482). Coordinates are
// 1-based (oldStart/newStart ≥ 1 when the respective Lines > 0).
// Content holds the raw hunk body with `+`/`-`/space line prefixes.
type Hunk struct {
	OldStart int    `json:"oldStart"`
	OldLines int    `json:"oldLines"`
	NewStart int    `json:"newStart"`
	NewLines int    `json:"newLines"`
	Content  string `json:"content"`
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

// ErrDependenciesRequired signals that the add request would
// activate an add-on whose declared dependencies
// (`domain.AddOnDependency`) are not satisfied in the current
// u-boot.yaml — at least one required service is not registered.
// Maps to LH-FA-CLI-006 exit code 10.
//
// The fail-fast path returns this sentinel; slice-v1-addons-deps
// T3 adds the four-mode CLI dispatch (`--with-deps` / `--yes` /
// `--no-interactive` / interactive prompt) on top so the error
// is reachable only when the user explicitly opted into fail-fast
// (no flag + non-interactive shell), or for tests bypassing the
// CLI.
var ErrDependenciesRequired = errors.New("service add-on requires missing dependencies")

// ErrProjectNotInitialized signals that BaseDir contains no
// `u-boot.yaml` (or one that cannot be parsed into the expected
// schema). LH-FA-ADD-001 requires an initialized project; the use
// case refuses to invent a config. The CLI surfaces it with a
// "run u-boot init" hint.
var ErrProjectNotInitialized = errors.New("project not initialized")

// ErrAddFileSystem signals an FS-write failure during `u-boot add`
// (LH-NFA-REL-003 "Abbruch bei kritischen Fehlern"). The use case
// wraps the raw FS error with this sentinel so the CLI's
// [ExitCode] mapper can route it to exit code 14 (technical
// persistence/filesystem-class), not the generic catch-all 1.
//
// On this error path the response is **non-empty**: it carries
// PlannedFiles[] (the calls captured up to the failure point) so the
// JSON envelope can show the partial progress. Per slice-v1-cli-json-
// dry-run-add T0-(b)'s Mid-Write-Failure scenario the user sees the
// captured calls plus a diagnostics[] entry pointing at the failed
// path.
var ErrAddFileSystem = errors.New("add: filesystem mutation failed")

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
