package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// templateFlags bundles the per-invocation flag state of
// `u-boot template list`. Today only --json (LH-FA-TPL-004 machine-
// readable form); future flags (`--filter`, `--sort`) can land
// additively without changing the runner signature.
type templateFlags struct {
	JSON bool
}

// newTemplateCommand builds the `u-boot template` Cobra subcommand
// group (slice-v1-template-list T3). The parent has no Run path —
// `u-boot template` alone prints the help text via Cobra's default
// behaviour. Children today: `list` only; `init --template <name>`
// lives on the existing `init` subcommand and gets implemented
// by slice-v1-template-init.
func newTemplateCommand(a *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Inspect the catalog of project templates",
		Long: `Inspect the catalog of external project templates u-boot ships.

Subcommands:

  u-boot template list [--json]   list available templates (LH-FA-TPL-004)

Coming in later slices:

  u-boot init --template <name>   render a template into a new project
                                  (slice-v1-template-init)
  u-boot init --template ./path   use a local template directory
                                  (slice-later-local-templates)`,
		Args: cobra.NoArgs,
	}
	cmd.AddCommand(newTemplateListCommand(a))
	return cmd
}

// newTemplateListCommand builds the `u-boot template list` Cobra
// leaf (LH-FA-TPL-004). Surfaces every template metadata entry the
// driven catalog enumerates; --json flips the renderer between the
// human tabular form (default) and a structured JSON array.
func newTemplateListCommand(a *App) *cobra.Command {
	flags := &templateFlags{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available project templates (LH-FA-TPL-004)",
		Long: `Print every built-in template with its name, description, and version.

The default render is a tabwriter-aligned table for humans; --json
emits a structured array (one object per template, all LH-FA-TPL-002
fields) suitable for downstream tooling.

Exit codes (LH-FA-CLI-006):
  0   success
  2   pure CLI / flag errors (unknown subcommand, unknown flag)
  14  catalog adapter failure (filesystem IO, malformed
      template.yaml in a fixture-backed adapter — the production
      embed.FS adapter validates at load time and never reaches
      this path in a passing CI build)

Examples:
  u-boot template list           # tabular layout
  u-boot template list --json    # JSON array for tooling`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runTemplateList(cmd.Context(), cmd.OutOrStdout(), *flags, a.templateListUseCase)
		},
	}
	cmd.Flags().BoolVar(&flags.JSON, "json", false,
		"emit a JSON array instead of the human-readable table (LH-FA-TPL-004 machine-readable form)")
	return cmd
}

// runTemplateList delegates to the use case and dispatches to the
// matching renderer. Split from the Cobra closure so unit tests can
// drive it with a fake use case + an io.Writer buffer, no Cobra
// involved.
func runTemplateList(
	ctx context.Context,
	out io.Writer,
	flags templateFlags,
	uc driving.TemplateListUseCase,
) error {
	resp, err := uc.List(ctx, driving.TemplateListRequest{})
	if err != nil {
		return err
	}
	if flags.JSON {
		return renderTemplateListJSON(out, resp.Templates)
	}
	return renderTemplateListText(out, resp.Templates)
}

// renderTemplateListText writes the human-readable tabwriter form:
// header line `NAME  DESCRIPTION  VERSION` followed by one row per
// template, padding to align columns. Empty catalog renders a
// short empty-state message instead of a header-only table.
func renderTemplateListText(out io.Writer, metas []domain.TemplateMetadata) error {
	if len(metas) == 0 {
		if _, err := fmt.Fprintln(out, "No templates available."); err != nil {
			return fmt.Errorf("write empty-state message: %w", err)
		}
		return nil
	}
	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(tw, "NAME\tDESCRIPTION\tVERSION"); err != nil {
		return fmt.Errorf("write template-list header: %w", err)
	}
	for _, m := range metas {
		if _, err := fmt.Fprintf(tw, "%s\t%s\t%s\n", m.Name, m.Description, m.Version); err != nil {
			return fmt.Errorf("write template-list row: %w", err)
		}
	}
	if err := tw.Flush(); err != nil {
		return fmt.Errorf("flush template-list writer: %w", err)
	}
	return nil
}

// renderTemplateListJSON marshals the metadata list as a pretty-
// indented JSON array. The CLI uses a private DTO ([templateJSON])
// instead of tagging the domain type so the domain layer stays
// presentation-agnostic (hexagonal purity) — adapters own the
// rendering shape.
//
// Nil slices are normalised to `[]` so JSON consumers always see
// arrays (never `null`); a missing `requiredTools` in `template.yaml`
// is semantically "no required tools", not "field absent".
func renderTemplateListJSON(out io.Writer, metas []domain.TemplateMetadata) error {
	dtos := make([]templateJSON, 0, len(metas))
	for _, m := range metas {
		dtos = append(dtos, toTemplateJSON(m))
	}
	data, err := json.MarshalIndent(dtos, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal template list: %w", err)
	}
	if _, err := out.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("write json: %w", err)
	}
	return nil
}

// templateJSON is the CLI-local JSON projection of
// [domain.TemplateMetadata]. Lives in the driving adapter so the
// domain layer stays presentation-agnostic (LH-FA-ARCH-002 /
// ADR-0002); the field tags pin the wire shape against future
// rename refactors in the domain struct.
type templateJSON struct {
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	Version         string                 `json:"version"`
	SupportedAddOns []string               `json:"supportedAddOns"`
	GeneratedFiles  []string               `json:"generatedFiles"`
	RequiredTools   []string               `json:"requiredTools"`
	Variables       []templateVariableJSON `json:"variables"`
}

// templateVariableJSON mirrors [domain.TemplateVariable] for JSON
// output. Default uses `omitempty` because a missing-default has a
// real semantic difference from `default: ""` once
// slice-v1-template-init lands the prompt path.
type templateVariableJSON struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Default     string `json:"default,omitempty"`
	Required    bool   `json:"required"`
}

func toTemplateJSON(m domain.TemplateMetadata) templateJSON {
	out := templateJSON{
		Name:            m.Name,
		Description:     m.Description,
		Version:         m.Version,
		SupportedAddOns: m.SupportedAddOns,
		GeneratedFiles:  m.GeneratedFiles,
		RequiredTools:   m.RequiredTools,
	}
	if out.SupportedAddOns == nil {
		out.SupportedAddOns = []string{}
	}
	if out.GeneratedFiles == nil {
		out.GeneratedFiles = []string{}
	}
	if out.RequiredTools == nil {
		out.RequiredTools = []string{}
	}
	out.Variables = make([]templateVariableJSON, 0, len(m.Variables))
	for _, v := range m.Variables {
		out.Variables = append(out.Variables, templateVariableJSON{
			Name:        v.Name,
			Description: v.Description,
			Default:     v.Default,
			Required:    v.Required,
		})
	}
	return out
}
