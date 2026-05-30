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
	// is the second-call shape).
	Action GenerateAction

	// Changed lists project-relative paths the handler mutated.
	// Sorted deterministically so the CLI summary order is stable.
	// Empty for NoOp; for Created / UpdatedBlock / RepairedManual it
	// names the artefact files actually written.
	Changed []string
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

// ErrGenerateManualConflict signals that the target artefact exists
// but has no `U-BOOT MANAGED BLOCK: init` marker (or the marker is
// malformed: BEGIN without END / duplicate BEGIN). The handler
// refuses to overwrite user content without an anchor; the CLI maps
// the sentinel to LH-FA-CLI-006 exit code 10 with a repair hint
// pointing at manual file rename / marker insertion (M7 has no
// `--replace` flag, see slice plan §Out of Scope).
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
