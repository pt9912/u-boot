package cli

import (
	"context"
	"errors"
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

// ErrTemplateSubcommandRequired is returned by bare `u-boot template
// --json` (Cluster-T_close; slice-v1-cli-json-dry-run-template
// T0-(a)/(f), moved from template-T3). `template` is a help-parent
// without its own data, and §1838/§420 make `subcommand` mandatory
// for command="template" — so bare template cannot emit a spec-valid
// envelope. It rejects with Exit 2 (LH-FA-CLI-006 usage class via
// [isUsageError]), envelope-LOS by design. Human mode (no --json)
// is unaffected and still prints `cmd.Help()`.
//
// This replaces the former allowlist-gate reject
// (ErrJSONNotImplemented) that the Cluster-T_close removed: it is
// RunE-borne so the behaviour survives the gate teardown without
// leaking help text on `--json template`.
var ErrTemplateSubcommandRequired = errors.New("u-boot template requires a subcommand (try `u-boot template list`)")

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
		// Help-Parent: prints help via cmd.Help() in human mode.
		// In --json mode bare `template` is rejected (Cluster-T_close):
		// §1838/§420 make `subcommand` mandatory for
		// command="template", and the help-parent has no data of its
		// own — so it cannot emit a spec-valid `command:"template"`
		// envelope. The reject is RunE-borne (not gate-borne) since
		// the Cluster-T_close removed the allowlist gate; it must
		// fire here so `--json template` does not leak help text.
		// Envelope-LOS by design — the only gate-less reject form of
		// the cluster (slice-v1-cli-json-dry-run-template T0-(a)/(f),
		// moved here from template-T3).
		RunE: func(cmd *cobra.Command, _ []string) error {
			if a.json {
				return ErrTemplateSubcommandRequired
			}
			return cmd.Help()
		},
	}
	cmd.AddCommand(newTemplateListCommand(a))
	return cmd
}

// newTemplateListCommand builds the `u-boot template list` Cobra
// leaf (LH-FA-TPL-004). Surfaces every template metadata entry the
// driven catalog enumerates; --json flips the renderer between the
// human tabular form (default) and a structured JSON array.
func newTemplateListCommand(a *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available project templates (LH-FA-TPL-004)",
		Long: `Print every built-in template with its name, description, and version.

The default render is a tabwriter-aligned table for humans; the
root --json flag (LH-NFA-USE-004) switches to a structured array
(one object per template, all LH-FA-TPL-002 fields) suitable for
downstream tooling.

Exit codes (LH-FA-CLI-006):
  0   success
  2   pure CLI / flag errors (unknown subcommand, unknown flag)
  14  catalog adapter failure (filesystem IO, malformed
      template.yaml in a fixture-backed adapter — the production
      embed.FS adapter validates at load time and never reaches
      this path in a passing CI build)

Examples:
  u-boot template list           # tabular layout
  u-boot template list --json    # JSON array for tooling (root --json flag)`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			// slice-v1-cli-json-dry-run-doctor T3/T4: lokales --json-
			// Flag wurde entfernt; --json wandert auf den Root-
			// PersistentFlag (a.json). Output-Format bleibt heutige
			// templateJSON-Array-Struktur — Envelope-Migration folgt
			// mit slice-v1-cli-json-dry-run-template (Carveouts-
			// Eintrag).
			return runTemplateList(cmd.Context(), cmd.OutOrStdout(), templateFlags{JSON: a.json}, a.templateListUseCase)
		},
	}
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
	// slice-v1-cli-json-dry-run-template T2: read-only error path
	// flows through reportErrorSub so a List failure produces a
	// spec-§1841 error envelope (subcommand "list"), not a raw
	// Cobra stderr error — Cluster-Symmetrie mit logs. dryRun/diff
	// are false (template list is read-only).
	resp, err := uc.List(ctx, driving.TemplateListRequest{})
	if err != nil {
		return reportErrorSub(out, err, nil, false, false, flags.JSON, "template", "list", mapTemplateErrorToDiagnostic, nil)
	}
	if flags.JSON {
		return writeTemplateListJSON(out, resp.Templates)
	}
	return renderTemplateListText(out, resp.Templates)
}

// mapTemplateErrorToDiagnostic maps a template-path error to a
// [diagnosticItem] (slice-v1-cli-json-dry-run-template T2 / T0-(f)).
// Two rows: the catalog-adapter IO failure
// ([driving.ErrTemplateCatalog]) is a technical-persistence error
// (LH-NFA-REL-003 / Exit 14 via isFilesystemError) — unreachable in
// a passing CI build because the production embed.FS adapter
// validates at load time; the default falls back to LH-FA-CLI-006.
// No new code-registry entry: only existing LH-codes are emitted.
func mapTemplateErrorToDiagnostic(err error) diagnosticItem {
	switch {
	case errors.Is(err, driving.ErrTemplateCatalog):
		return diagnosticItem{Level: "error", Code: "LH-NFA-REL-003", Message: err.Error()}
	default:
		return diagnosticItem{Level: "error", Code: "LH-FA-CLI-006", Message: err.Error()}
	}
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

// writeTemplateListJSON emits the LH-NFA-USE-004 Minimalkontrakt-
// Envelope for `template list --json` (slice-v1-cli-json-dry-run-
// template T2). The `[]templateJSON` projection rides in the
// envelope's `data` field via [newDataEnvelope] with
// `command="template"`, `subcommand="list"` (§322 subcommand-
// Pflicht), `diagnostics=[]`, `exitCode=0`.
//
// Breaking change vs. the pre-migration shape (Doctor-Slice
// carveout): the output moves from a top-level pretty-indented JSON
// array to a single-line envelope object whose `data` holds the
// array. Documented in CHANGELOG `### Changed` (T4).
//
// Nil slices are normalised to `[]` inside [toTemplateJSON] so JSON
// consumers always see arrays (never `null`); an empty catalog
// renders `"data": []`.
func writeTemplateListJSON(out io.Writer, metas []domain.TemplateMetadata) error {
	dtos := make([]templateJSON, 0, len(metas))
	for _, m := range metas {
		dtos = append(dtos, toTemplateJSON(m))
	}
	env := newDataEnvelope("template", "list", dtos, nil, 0)
	return writeEnvelope(out, env)
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
// output. Default is emitted unconditionally (no `omitempty`) so
// template authors who declare `default: ""` are not silently
// stripped — preserving the missing-vs-empty distinction the
// future slice-v1-template-init prompt path will need
// (review-followup N3 — earlier version had `omitempty`, which
// destroyed the very distinction the comment justifying it
// claimed to preserve).
type templateVariableJSON struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Default     string `json:"default"`
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
