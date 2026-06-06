package driving

import (
	"context"
	"errors"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// GenerateRequest is the input for [GenerateUseCase.Generate]. It is
// the application-layer expression of `u-boot generate <artifact>`
// (LH-FA-GEN-001). The CLI adapter translates the positional argument
// through [domain.NewArtifact] before constructing the request.
type GenerateRequest struct {
	// BaseDir is the absolute path of the initialized u-boot project
	// the artefact is generated into. Mandatory; the CLI adapter
	// defaults it to the current working directory.
	BaseDir string

	// Artifact selects which artefact handler runs (changelog, readme,
	// env-example, devcontainer). Constructed via [domain.NewArtifact]
	// in the CLI adapter; the application service trusts the value.
	Artifact domain.Artifact

	// AllowExternalFeatureSources lists LH-FA-DEV-003 source URLs
	// that should be appended to `devcontainer.featureSources.allow`
	// before the devcontainer artefact is rendered (Spec §715). Only
	// applied when Artifact == ArtifactDevcontainer; ignored for the
	// other artefacts (changelog / readme / env-example) — the CLI
	// adapter restricts the flag to the devcontainer subcommand,
	// but the use case re-checks for safety. Empty leaves the
	// allowlist untouched. Slice-v1-devcontainer-features T4.
	AllowExternalFeatureSources []string

	// PreviewMode encodes the --dry-run × --diff flag combination per
	// slice-v1-cli-json-dry-run-generate T0-(b) (inherited 1:1 from
	// init T0-(b)/add T0-(b) truth table — kein generate-Prefix-Alias,
	// init-Slice-T0-(c) Alias-Lebensdauer-Pflicht zwingt direktes
	// [PreviewMode]). Default zero value [PreviewNone] preserves the
	// existing production-write behaviour; non-zero values route every
	// FS access of this Generate() invocation through the recordingfs
	// adapter (Composition-Root generateFSFactory closure, T4).
	//
	// Generate hat KEIN SilenceProgress-Feld — der GenerateService
	// nutzt keinen [driven.Progress]-Port (nur driven.FileSystem +
	// driven.YAMLCodec + driven.Logger), kein stdout-Pollution-Risk.
	PreviewMode PreviewMode
}

// GenerateAction classifies what `u-boot generate` did with the
// artefact, so the CLI can render a deterministic one-line summary
// without re-deriving the state. The four values cover the M7-T1
// state-machine outcomes documented in the slice plan; per-artefact
// handlers (T2..T5) populate them according to their own state tables.
type GenerateAction int

const (
	// GenerateActionCreated means the artefact did not exist and was
	// newly written.
	GenerateActionCreated GenerateAction = iota

	// GenerateActionUpdatedBlock means the artefact already existed
	// with a `U-BOOT MANAGED BLOCK: init` marker; only that block
	// was rerendered and spliced. Content outside the block is
	// byte-identical to the prior state.
	GenerateActionUpdatedBlock

	// GenerateActionNoOp means the artefact already existed and the
	// rerendered block was byte-identical to the existing one;
	// nothing was written.
	GenerateActionNoOp

	// GenerateActionRepairedManual means the artefact existed in a
	// user-edited shape that the handler refused to re-render but
	// repaired structurally (e.g. T4 changelog: inserts a missing
	// `## [Unreleased]` header without touching the managed block).
	GenerateActionRepairedManual
)

// String returns the canonical lowercase identifier of the action,
// used by the CLI summary print and the [Logger] port. Stable: CI
// log scrapers may pin these strings.
func (g GenerateAction) String() string {
	switch g {
	case GenerateActionCreated:
		return "created"
	case GenerateActionUpdatedBlock:
		return "updated-block"
	case GenerateActionNoOp:
		return "no-op"
	case GenerateActionRepairedManual:
		return "repaired-manual"
	default:
		return "unknown"
	}
}

// GenerateResponse is the output of [GenerateUseCase.Generate]. The
// CLI adapter renders it via `printGenerateSummary` (T6); the Changed
// list lets the user see at a glance which files were touched.
type GenerateResponse struct {
	// Artifact echoes the artefact that was processed.
	Artifact domain.Artifact

	// Action classifies the outcome (LH-FA-GEN-005 idempotency: NoOp
	// is the second-call shape). The CLI adapter renders this enum
	// to the `data.action: "<created|updated-block|no-op|repaired-
	// manual>"`-field in the JSON envelope (slice-v1-cli-json-dry-
	// run-generate T0-(f) — kein Port-Feld dafür, Rendering lebt im
	// Envelope-Layer T5).
	Action GenerateAction

	// Changed lists project-relative paths the handler mutated.
	// Sorted deterministically so the CLI summary order is stable.
	// Empty for NoOp; for Created / UpdatedBlock / RepairedManual it
	// names the artefact files actually written.
	Changed []string

	// PlannedFiles is the FS-plan emitted when [GenerateRequest.
	// PreviewMode] is non-zero (slice-v1-cli-json-dry-run-generate
	// T2 / inherited from init T2 / add T0-(i)). One entry per
	// mutated path captured by the recorder, in the order the use
	// case attempted them. Empty for PreviewNone (no recorder wired)
	// and for true no-ops. The [PlannedFile] type itself trägt die
	// NewContent/OldContent-Felder mit `json:"-"` (siehe addservice.go
	// Definition), damit Raw-Bytes nicht in der CLI-Wire-Form landen
	// (Spec §326). Der Port serialisiert nichts direkt — die CLI baut
	// eigene Wire-Typen in `cli/jsonenvelope.go` (plannedFile/hunk).
	//
	// Mid-Write-Failure-Semantik (R4 Recorder-Realität,
	// recordingfs.go:139 zeichnet vor Delegieren auf): bei einem
	// underlying.WriteFile-Failure des i-ten Files enthält
	// PlannedFiles trotzdem die i Captures inkl. des fehlgeschlagenen
	// — die T0-(i) Devcontainer-Phase-2-Half-Write-AK und der
	// generate T6 Pin testen exakt diese Form.
	PlannedFiles []PlannedFile

	// Changes mirrors PlannedFiles' paths with their line-count
	// summaries (LH-FA-CLI-007 §365-371). Filled only in preview
	// modes; nil for PreviewNone. Count semantics follow add T0-(g)
	// (1:1 inherited): create = CountLines(NewContent); modify = sum
	// of `+`-lines via diff.CountAdditions; delete = 0. Today
	// populated by the CLI-adapter's mapPlannedFilesToWire helper,
	// not directly by the application service (matches add/init
	// pattern; the use case fills PlannedFiles, the CLI computes
	// Changes from them).
	Changes []ChangeEntry
}

// All Generate sentinels below live in the `driving` package (not in
// `application`) so the CLI adapter can branch on them via
// [errors.Is] without importing `application` — the LH-FA-ARCH-003
// depguard rule forbids that cross-layer import.

// ErrArtifactUnknown signals that the CLI adapter received an
// `<artifact>` positional argument that is not in the LH-FA-GEN-001
// catalogue (`changelog` / `readme` / `env-example` / `devcontainer`).
// The CLI wraps the [domain.ErrInvalidArtifact] from
// [domain.NewArtifact] with this sentinel so [ExitCode] can map it
// to LH-FA-CLI-006 exit code 2 (CLI validation) instead of code 10.
// The distinction matters: `u-boot add <unknown-service>` maps to 10
// (fachlich), but `u-boot generate <unknown-artifact>` maps to 2
// because the spec calls it out explicitly (§LH-FA-GEN-001).
var ErrArtifactUnknown = errors.New("unknown artifact")

// ErrGenerateManualConflict signals a fachlich-blocking state in
// which an automated regeneration would clobber user content. Two
// concrete sub-cases share the sentinel because both map to
// LH-FA-CLI-006 exit code 10 with the same "user must intervene
// manually" semantics:
//
//  1. The target artefact exists but has no
//     `U-BOOT MANAGED BLOCK: init` marker (or the marker is malformed:
//     BEGIN without END / duplicate BEGIN). M7 has no `--replace`
//     flag (see slice-m7-generate.md §Out of Scope) so the handler
//     refuses to overwrite user content without an anchor.
//  2. A neighbouring project YAML file (`compose.yaml` under
//     `u-boot generate devcontainer`) is unparseable
//     (yaml.v3-reported parse error wrapped via
//     [driven.ErrYAMLParse]). The handler routes the parse path
//     through this sentinel instead of [ErrGenerateFileSystem] so
//     the user sees the spec-correct code 10 (fachlich, "fix the
//     YAML") instead of code 14 (technical). Introduced by the
//     `slice-v1-yaml-parse-error-sentinel` follow-up — see that
//     slice for the architectural rationale.
//
// The CLI maps the sentinel via `isValidationError` in cli.go.
var ErrGenerateManualConflict = errors.New("generate: managed block missing or malformed")

// ErrGenerateFileSystem wraps an unexpected I/O or permissions error
// from [driven.FileSystem] (Read/Write/Stat). Mapped to LH-FA-CLI-006
// exit code 14 (technical persistence failure) in T6 via
// `isFilesystemError` in cli.go.
//
// Wrap (not direct [errors.Is] against a driven sentinel) because
// the pre-T1 scan of `internal/hexagon/port/driven/` confirmed there
// is no `driven.ErrFileSystem*` sentinel today — only
// `ErrDockerUnavailable`, `ErrComposeRuntime`, and the YAML-codec
// sentinels. If a future slice introduces `driven.ErrFileSystem`,
// this wrap can be replaced with a `errors.Is` against the driven
// sentinel without touching the CLI exit-code mapping.
var ErrGenerateFileSystem = errors.New("generate: filesystem error")

// GenerateUseCase is the driving-port for `u-boot generate <artifact>`.
// The CLI adapter (T6) holds a reference and calls [Generate] from
// the Cobra command handler.
//
// Contract:
//
//   - On success the response carries Artifact = req.Artifact and an
//     Action / Changed pair consistent with the per-artefact state
//     machine documented in slice-m7-generate.md.
//   - On failure the response is the zero value. The error wraps one
//     of the sentinels above, or [ErrProjectNotInitialized] (reused
//     from M5/M6 — `u-boot generate` requires an initialized project),
//     or [domain.ErrInvalidArtifact] when the CLI did not validate
//     before constructing the request.
//
// Idempotence guarantee (LH-FA-GEN-005): calling [Generate] twice
// with the same request is safe; the second call returns
// Action = [GenerateActionNoOp] and Changed = nil when nothing has
// changed since the first call.
type GenerateUseCase interface {
	Generate(ctx context.Context, req GenerateRequest) (GenerateResponse, error)
}
