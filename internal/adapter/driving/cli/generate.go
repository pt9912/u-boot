package cli

import (
	"context"
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
func newGenerateCommand(a *App) *cobra.Command {
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

For `+"`changelog`"+`, the handler is conservative when the init block
has been user-edited (entries added under `+"`### Added`"+`): it does
not re-render the block; it only inserts a missing `+"`## [Unreleased]`"+`
header before the first release section if needed (RepairedManual
action).

Exit codes (LH-FA-CLI-006):
  0   success
  2   unknown artefact name; missing positional argument
  10  no u-boot.yaml (project not initialised);
      file present without an `+"`init`"+` managed block; malformed block
  14  filesystem error (permission, disk full, race)

Examples:
  u-boot generate changelog        # create or refresh CHANGELOG.md
  u-boot generate readme           # create or refresh README.md
  u-boot generate env-example      # create or refresh .env.example
  u-boot generate devcontainer     # both .devcontainer/ files`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenerate(cmd.Context(), cmd.OutOrStdout(), args, a.generateUseCase, a.getwd)
		},
	}
	return cmd
}

// runGenerate is split from the Cobra closure for direct unit-testing
// (no Cobra command construction needed). ctx is the first parameter
// explicitly so contextcheck can see the propagation.
//
// Parses the positional argument via [domain.NewArtifact]; on
// validation failure wraps the domain error in
// [driving.ErrArtifactUnknown] so the CLI exit-code mapping sends
// the right code (2, not 10). Delegates the project-state and per-
// artefact state-machine concerns to the use case.
func runGenerate(
	ctx context.Context,
	out io.Writer,
	args []string,
	uc driving.GenerateUseCase,
	getwd func() (string, error),
) error {
	artifact, err := domain.NewArtifact(args[0])
	if err != nil {
		return fmt.Errorf("%w: %v", driving.ErrArtifactUnknown, err)
	}

	cwd, err := getwd()
	if err != nil {
		return fmt.Errorf("determine working directory: %w", err)
	}

	resp, err := uc.Generate(ctx, driving.GenerateRequest{
		BaseDir:  cwd,
		Artifact: artifact,
	})
	if err != nil {
		return err
	}
	printGenerateSummary(out, resp)
	return nil
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

