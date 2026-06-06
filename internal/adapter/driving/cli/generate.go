package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// newGenerateCommand builds the `u-boot generate <artifact>` Cobra
// subcommand (LH-FA-GEN-001). The positional argument is parsed via
// [domain.NewArtifact]; on failure the error is wrapped in
// [driving.ErrArtifactUnknown] so [ExitCode] maps it to exit code
// **2** (CLI validation), unlike `u-boot add <unknown>` which maps
// to code 10. The distinction is spec-mandated (§LH-FA-GEN-001):
// "Bei unbekanntem Artefakt muss der Befehl mit Exit Code 2
// abbrechen und die erlaubten Werte explizit zurückgeben."
// generateFlags bundles the per-invocation flag state of
// `u-boot generate`. The LH-FA-DEV-003 allowlist seed flag is only
// meaningful for `generate devcontainer`; the use case re-checks
// the artefact kind before applying.
type generateFlags struct {
	AllowExternalFeatureSources []string

	// DryRun / Diff / JSON (slice-v1-cli-json-dry-run-generate T5):
	// LH-FA-CLI-007/008/§1841 flags. DryRun/Diff route Generate()
	// through the RecordingFileSystem via the per-request fsFactory
	// (T4); together with JSON they form the three voll-schema/
	// minimal output paths analog to add/init. JSON is read-through
	// from the root persistent --json (LH-NFA-USE-004).
	DryRun bool
	Diff   bool
	JSON   bool
}

func newGenerateCommand(a *App) *cobra.Command {
	flags := &generateFlags{}
	cmd := &cobra.Command{
		Use:   "generate <artifact>",
		Short: "Generate or update a u-boot-managed artefact (changelog/readme/env-example/devcontainer)",
		Long: `Generate or update one of the four u-boot-managed project artefacts:

  changelog     CHANGELOG.md (LH-FA-GEN-002 / LH-AK-007)
  readme        README.md (LH-FA-GEN-003)
  env-example   .env.example (LH-FA-GEN-004)
  devcontainer  .devcontainer/devcontainer.json + Dockerfile (LH-FA-DEV-001)

The command runs in an initialised project (u-boot.yaml present) and
is idempotent (LH-FA-GEN-005): a second invocation against an
artefact whose managed block already matches the rendered template
returns a no-op without touching the file. Existing manual content
outside the managed block is preserved byte-identically.

For ` + "`changelog`" + `, the handler is conservative when the init block
has been user-edited (entries added under ` + "`### Added`" + `): it does
not re-render the block; it only inserts a missing ` + "`## [Unreleased]`" + `
header before the first release section if needed (RepairedManual
action).

Exit codes (LH-FA-CLI-006):
  0   success
  2   unknown artefact name; missing positional argument
  10  no u-boot.yaml (project not initialised);
      file present without an ` + "`init`" + ` managed block; malformed block
  14  filesystem error (permission, disk full, race)

Examples:
  u-boot generate changelog        # create or refresh CHANGELOG.md
  u-boot generate readme           # create or refresh README.md
  u-boot generate env-example      # create or refresh .env.example
  u-boot generate devcontainer     # both .devcontainer/ files`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Read-through persistent --json (LH-NFA-USE-004) from
			// the App; Cobra has parsed it by RunE-time.
			flags.JSON = a.json
			return runGenerate(cmd.Context(), cmd.OutOrStdout(), args, *flags, a.generateUseCase, a.getwd)
		},
	}
	cmd.Flags().StringSliceVar(&flags.AllowExternalFeatureSources, "allow-external-feature-sources", nil,
		"append the given URLs to devcontainer.featureSources.allow before generating (LH-FA-DEV-003; only valid for `generate devcontainer`; comma-separated, repeatable). `--yes` does not substitute (LH-NFA-SEC-004).")
	cmd.Flags().BoolVar(&flags.DryRun, "dry-run", false,
		"preview the planned changes without writing files (LH-FA-CLI-007)")
	cmd.Flags().BoolVar(&flags.Diff, "diff", false,
		"render a unified diff of the planned changes (LH-FA-CLI-008)")
	return cmd
}

// generateEnvelopeData is the typed `data` carrier for the JSON
// envelope of `u-boot generate` (slice-v1-cli-json-dry-run-generate
// T0-(f)/(q)/(p)). Artifact is always set; Action stays empty on
// the error path so omitempty drops it (T0-(q) Festzurrung —
// Use-Case-Response is Zero on the error path, no Action exists).
type generateEnvelopeData struct {
	Artifact string `json:"artifact"`
	Action   string `json:"action,omitempty"`
}

// runGenerate is split from the Cobra closure for direct unit-testing
// (no Cobra command construction needed). ctx is the first parameter
// explicitly so contextcheck can see the propagation.
//
// Three output paths (slice-v1-cli-json-dry-run-generate T5):
//
//   - Human (no --json): printGenerateSummary plus optional
//     writeDiff when --diff is set.
//   - --json (Minimalkontrakt / minimal+data path): newDataEnvelope
//     with command="generate", data={artifact, action}, no
//     plannedFiles/changes.
//   - --dry-run --json / --diff --json: newFullEnvelope with
//     plannedFiles/changes plus data={artifact, action}; optional
//     hunks when --diff.
//
// Error-Pfad (T0-(q) R6/R7-Festzurrung): reportError mit
// `data={"artifact":"<…>"}` (kein action — Zero-Response). Der
// per-artifact LH-Code für ErrGenerateManualConflict wird via
// mapGenerateErrorToDiagnostic(err, artifact)-Closure aufgelöst.
func runGenerate(
	ctx context.Context,
	out io.Writer,
	args []string,
	flags generateFlags,
	uc driving.GenerateUseCase,
	getwd func() (string, error),
) error {
	artifact, err := domain.NewArtifact(args[0])
	if err != nil {
		// Unknown artifact — data.artifact would be misleading
		// (the user passed something we cannot classify). Pass nil
		// data; the error message itself carries the reject reason.
		wrapped := fmt.Errorf("%w: %v", driving.ErrArtifactUnknown, err)
		// Artifact is unknown — pass zero value (the mapper routes
		// ErrArtifactUnknown to LH-FA-CLI-006 without consulting the
		// artifact field, so the value choice doesn't matter here).
		var zeroArtifact domain.Artifact
		mapErr := func(e error) diagnosticItem {
			return mapGenerateErrorToDiagnostic(e, zeroArtifact)
		}
		return reportError(out, wrapped, nil, flags.DryRun, flags.Diff, flags.JSON, "generate", mapErr, nil)
	}

	data := generateEnvelopeData{Artifact: artifact.String()}
	mapErr := func(e error) diagnosticItem {
		return mapGenerateErrorToDiagnostic(e, artifact)
	}

	// Spec §714-717: --allow-external-feature-sources is only
	// valid for `generate devcontainer`. Reject early on other
	// artefacts so the user gets a clear "wrong command" message
	// rather than a silent no-op.
	if len(flags.AllowExternalFeatureSources) > 0 && artifact != domain.ArtifactDevcontainer {
		wrapped := fmt.Errorf(
			"%w: --allow-external-feature-sources is only valid for `generate devcontainer` (Spec §714-717); got `generate %s`",
			driving.ErrArtifactUnknown, artifact)
		return reportError(out, wrapped, nil, flags.DryRun, flags.Diff, flags.JSON, "generate", mapErr, data)
	}

	cwd, err := getwd()
	if err != nil {
		return reportError(out, fmt.Errorf("determine working directory: %w", err), nil, flags.DryRun, flags.Diff, flags.JSON, "generate", mapErr, data)
	}

	mode := previewModeFromFlags(flags.DryRun, flags.Diff)
	req := driving.GenerateRequest{
		BaseDir:                     cwd,
		Artifact:                    artifact,
		AllowExternalFeatureSources: flags.AllowExternalFeatureSources,
		PreviewMode:                 mode,
	}

	resp, genErr := uc.Generate(ctx, req)
	if genErr != nil {
		// Error-Envelope trägt data.artifact aber kein data.action
		// (Zero-Response auf Error-Pfad — T0-(q)).
		return reportError(out, genErr, resp.PlannedFiles, flags.DryRun, flags.Diff, flags.JSON, "generate", mapErr, data)
	}

	if flags.JSON {
		return writeGenerateJSON(out, resp, flags.DryRun, flags.Diff, artifact)
	}

	if flags.Diff {
		if err := writeDiff(out, resp.PlannedFiles); err != nil {
			return err
		}
	}
	printGenerateSummary(out, resp)
	return nil
}

// writeGenerateJSON renders the success-path JSON envelope. Two
// shapes (T0-(m)/(f) Festzurrungen):
//
//   - dryRun=false && diff=false → newDataEnvelope (Minimalkontrakt
//     plus `data: {"artifact": "<…>", "action": "<…>"}`).
//   - dryRun=true || diff=true   → newFullEnvelope mit
//     plannedFiles/changes plus dem gleichen data-Träger; optional
//     hunks im --diff-Pfad.
//
// Im NoOp-Fall sind `plannedFiles[]` UND `changes[]` leer (T0-(f)
// Festzurrung); Konsumenten leiten NoOp aus `data.action: "no-op"`
// UND der Array-Leerheit ab.
func writeGenerateJSON(out io.Writer, resp driving.GenerateResponse, dryRun, diffFlag bool, artifact domain.Artifact) error {
	data := generateEnvelopeData{
		Artifact: artifact.String(),
		Action:   resp.Action.String(),
	}
	if !dryRun && !diffFlag {
		env := newDataEnvelope("generate", "", data, nil, 0)
		return writeEnvelope(out, env)
	}
	pfs, chs := mapPlannedFilesToWire(resp.PlannedFiles, diffFlag)
	env := newFullEnvelope("generate", "", dryRun, diffFlag, pfs, chs, data, nil, 0)
	return writeEnvelope(out, env)
}

// mapGenerateErrorToDiagnostic maps a generate-path error to a
// diagnosticItem with the spec-konforme LH-Kennung per T0-(e)
// Switch-Order-Pflicht plus T0-(e) per-Artefakt LH-Code-Tabelle.
//
// Order matters (Multi-`%w`-wraps from T3): FS-class FIRST so an
// inner fachlich-Sentinel doesn't downgrade Exit 14 to Exit 10.
// LH-FA-DEV-003 (`ErrConfigValueInvalid`-URL-Reject, Spec §720)
// comes before ManualConflict because the typed wrap is
// devcontainer-specific and unambiguous; ManualConflict
// resolves the per-artefact code from the `artifact`-Parameter
// (changelog→GEN-002, readme→GEN-003, env-example→GEN-004,
// devcontainer→DEV-001).
//
// Default fallback maps to LH-FA-CLI-006 / Exit 1 via cli.ExitCode
// (NICHT automatisch Exit 2 — isUsageError matched nur unsere
// CLI-Validation-Sentinels, nicht den generischen Default-Pfad).
func mapGenerateErrorToDiagnostic(err error, artifact domain.Artifact) diagnosticItem {
	switch {
	case errors.Is(err, driving.ErrGenerateFileSystem):
		return diagnosticItem{Level: "error", Code: "LH-NFA-REL-003", Message: err.Error()}
	case errors.Is(err, driving.ErrConfigValueInvalid):
		return diagnosticItem{Level: "error", Code: "LH-FA-DEV-003", Message: err.Error()}
	case errors.Is(err, driving.ErrGenerateManualConflict):
		return diagnosticItem{Level: "error", Code: manualConflictCodeFor(artifact), Message: err.Error()}
	case errors.Is(err, driving.ErrArtifactUnknown):
		return diagnosticItem{Level: "error", Code: "LH-FA-CLI-006", Message: err.Error()}
	case errors.Is(err, driving.ErrProjectNotInitialized):
		return diagnosticItem{Level: "error", Code: "LH-FA-INIT-001", Message: err.Error()}
	default:
		return diagnosticItem{Level: "error", Code: "LH-FA-CLI-006", Message: err.Error()}
	}
}

// manualConflictCodeFor resolves the per-artefact LH-Code for
// ErrGenerateManualConflict (T0-(e) Diagnostic-Code-Tabelle).
// Unknown / zero-value artefact falls back to changelog's code
// (most common path) so a misclassified branch doesn't drop to
// LH-FA-CLI-006 silently.
func manualConflictCodeFor(artifact domain.Artifact) string {
	switch artifact {
	case domain.ArtifactReadme:
		return "LH-FA-GEN-003"
	case domain.ArtifactEnvExample:
		return "LH-FA-GEN-004"
	case domain.ArtifactDevcontainer:
		return "LH-FA-DEV-001"
	}
	return "LH-FA-GEN-002"
}

// printGenerateSummary writes a short, deterministic summary of the
// generate outcome. Four shapes, mirroring the four
// [driving.GenerateAction] values:
//
//	Created         → "Generated <artifact> (<paths>)."
//	UpdatedBlock    → "Updated <artifact> managed block (<paths>)."
//	NoOp            → "<artifact> already up to date; no changes."
//	RepairedManual  → "Repaired <artifact> structure (<paths>)."
//
// The Changed slice is rendered as a comma-separated parenthetical
// so the caller sees at a glance which files were written; for NoOp
// the line is intentionally bare (no Changed list) because nothing
// was touched.
func printGenerateSummary(out io.Writer, resp driving.GenerateResponse) {
	name := resp.Artifact.String()
	switch resp.Action {
	case driving.GenerateActionNoOp:
		fmt.Fprintf(out, "%s already up to date; no changes.\n", name)
	case driving.GenerateActionCreated:
		fmt.Fprintf(out, "Generated %s (%s).\n", name, strings.Join(resp.Changed, ", "))
	case driving.GenerateActionUpdatedBlock:
		fmt.Fprintf(out, "Updated %s managed block (%s).\n", name, strings.Join(resp.Changed, ", "))
	case driving.GenerateActionRepairedManual:
		fmt.Fprintf(out, "Repaired %s structure (%s).\n", name, strings.Join(resp.Changed, ", "))
	default:
		// Defensive: a future GenerateAction value renders both its
		// String() form and the Changed list rather than silently
		// truncating to "Generated <name>" (review-followup N5).
		fmt.Fprintf(out, "%s action %s; changed: %s\n",
			name, resp.Action, strings.Join(resp.Changed, ", "))
	}
}
