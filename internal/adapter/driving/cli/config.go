package cli

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// configShowData is the typed `data` carrier for `u-boot config`
// (bare → subcommand "show"), slice-v1-cli-json-dry-run-config
// T0-(c). Body is the full u-boot.yaml as a Go string (JSON-safe;
// UTF-8-escape on CR/Tab/non-printables). No omitempty — an empty
// u-boot.yaml is a legitimate `""` body, not absence.
type configShowData struct {
	Body string `json:"body"`
}

// configGetData is the typed `data` carrier for `config get`
// (T0-(c)). Both fields without omitempty: a bare-scalar value of
// `""` is legitimate (e.g. an unset-but-present optional).
type configGetData struct {
	Path  string `json:"path"`
	Value string `json:"value"`
}

// configSetData is the typed `data` carrier for `config set`
// (T0-(c)/(d)). NoOp without omitempty (a legitimate success=false
// signal); OldValue/NewValue without omitempty (empty-string `""`
// = legitimate initial-unset). AppendedSources is set ONLY for the
// `devcontainer.featureSources.allow` path (T0-(c) Hybrid): it
// echoes the raw `--allow-external-feature-sources` input so the
// consumer sees what the flag tried to append (omitempty elsewhere).
type configSetData struct {
	Path            string   `json:"path"`
	OldValue        string   `json:"oldValue"`
	NewValue        string   `json:"newValue"`
	NoOp            bool     `json:"noOp"`
	AppendedSources []string `json:"appendedSources,omitempty"`
}

// configShowFlags / configGetFlags bundle the read-only forms'
// per-invocation flag state. JSON/Quiet read through from the root
// (Cluster-T0-(a) doctor-Pattern: `--quiet --json` ≡ `--json`).
// DryRun/Diff are registered on these read-only commands ONLY so
// Cobra parses them cleanly and the RunE can reject them
// Envelope-konform with [ErrDryRunNotApplicable] (T0-(g) Option
// (i.a)) instead of letting Cobra emit a raw `unknown flag`.
type configShowFlags struct {
	JSON   bool
	Quiet  bool
	DryRun bool
	Diff   bool
}

type configGetFlags struct {
	JSON   bool
	Quiet  bool
	DryRun bool
	Diff   bool
}

// configSetFlags bundles the per-invocation flag state of
// `u-boot config set`. The LH-FA-DEV-003 allowlist seed flag is
// only meaningful when the positional path is
// `devcontainer.featureSources.allow` (Spec §717); the use case
// re-checks before applying. DryRun/Diff drive the LH-FA-CLI-007/008
// preview-mode (T5); JSON/Quiet read through from the root.
type configSetFlags struct {
	AllowExternalFeatureSources []string
	JSON                        bool
	Quiet                       bool
	DryRun                      bool
	Diff                        bool
}

// newConfigCommand builds the `u-boot config` Cobra subcommand
// tree (LH-FA-CONF-001 / LH-FA-CONF-005). Three forms:
//
//   - `u-boot config`              — print full u-boot.yaml (Show).
//   - `u-boot config get <path>`   — print a single scalar.
//   - `u-boot config set <path> <value>` — schema-validated write.
//
// The parent command runs Show when no subcommand is given; `get`
// and `set` are children so help / error messages stay in the
// canonical `u-boot config <verb>` namespace.
func newConfigCommand(a *App) *cobra.Command {
	flags := &configShowFlags{}
	cmd := &cobra.Command{
		Use:   "config [get|set]",
		Short: "Show or modify the u-boot project configuration",
		Long: `Inspect or change values in u-boot.yaml.

Three forms:

  u-boot config                       # print the full configuration
  u-boot config get <path>            # print a single value
  u-boot config set <path> <value>    # write a single value

Writable paths (LH-FA-CONF-001 / §D1 of slice-m8-config.md):

  project.name                 # the project identifier (LH-FA-INIT-006 rules)
  devcontainer.enabled         # bool

Readable paths additionally include:

  services.<svc>.enabled       # bool; write goes through u-boot add/remove

set runs a two-stage schema validation before touching the file
(slice §D3): the patched bytes round-trip through the
ubootYAMLConfig struct AND each path's domain validator. On
validation failure the file stays byte-identical. set is
idempotent. --dry-run/--diff preview without writing (only valid
for set; rejected on the read-only forms). --json emits the
LH-FA-CLI-007/008 envelope.

Exit codes (LH-FA-CLI-006):
  0   success
  2   CLI validation (unknown subcommand, missing args, --dry-run/
      --diff on a read-only form)
  10  fachlich (no u-boot.yaml; unknown path; invalid value;
      write-rejected path; post-patch sanity; value not set on Get)
  14  filesystem error

Examples:
  u-boot config
  u-boot config get project.name
  u-boot config set project.name my-service
  u-boot config set devcontainer.enabled true --dry-run --json`,
		Args: configArgsValidator(a, "show", cobra.NoArgs),
		RunE: func(cmd *cobra.Command, _ []string) error {
			flags.JSON = a.json
			flags.Quiet = a.quiet
			return runConfigShow(cmd.Context(), cmd.OutOrStdout(), *flags, a.configUseCase, a.getwd)
		},
	}
	registerConfigPreviewRejectFlags(cmd, "config", &flags.DryRun, &flags.Diff)

	cmd.AddCommand(newConfigGetCommand(a))
	cmd.AddCommand(newConfigSetCommand(a))
	return cmd
}

func newConfigGetCommand(a *App) *cobra.Command {
	flags := &configGetFlags{}
	cmd := &cobra.Command{
		Use:   "get <path>",
		Short: "Print a single configuration value",
		Args:  configArgsValidator(a, "get", cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.JSON = a.json
			flags.Quiet = a.quiet
			return runConfigGet(cmd.Context(), cmd.OutOrStdout(), args, *flags, a.configUseCase, a.getwd)
		},
	}
	registerConfigPreviewRejectFlags(cmd, "config get", &flags.DryRun, &flags.Diff)
	return cmd
}

func newConfigSetCommand(a *App) *cobra.Command {
	flags := &configSetFlags{}
	cmd := &cobra.Command{
		Use:   "set <path> <value>",
		Short: "Set a single configuration value (schema-validated)",
		Args:  configArgsValidator(a, "set", cobra.ExactArgs(2)),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.JSON = a.json
			flags.Quiet = a.quiet
			return runConfigSet(cmd.Context(), cmd.OutOrStdout(), args, *flags, a.configUseCase, a.getwd)
		},
	}
	cmd.Flags().StringSliceVar(&flags.AllowExternalFeatureSources, "allow-external-feature-sources", nil,
		"additional URLs to append to devcontainer.featureSources.allow (LH-FA-DEV-003; only valid when <path> is devcontainer.featureSources.allow; cumulative with the positional value).")
	cmd.Flags().BoolVar(&flags.DryRun, "dry-run", false,
		"preview the planned change without writing the file (LH-FA-CLI-007)")
	cmd.Flags().BoolVar(&flags.Diff, "diff", false,
		"render a unified diff of the planned change (LH-FA-CLI-008)")
	return cmd
}

// registerConfigPreviewRejectFlags registers --dry-run/--diff on a
// read-only config command (bare/get) so Cobra parses them cleanly;
// the RunE rejects them via [ErrDryRunNotApplicable] (T0-(g) Option
// (i.a)). The flags bind to the form's flag-struct fields so the
// RunE actually sees the user's input; the help text states the
// rejection so the synthetic flags are not mistaken for working
// previews.
func registerConfigPreviewRejectFlags(cmd *cobra.Command, form string, dryRun, diff *bool) {
	cmd.Flags().BoolVar(dryRun, "dry-run", false,
		fmt.Sprintf("not applicable to `%s` (read-only); only valid for `config set` — rejected with exit 2", form))
	cmd.Flags().BoolVar(diff, "diff", false,
		fmt.Sprintf("not applicable to `%s` (read-only); only valid for `config set` — rejected with exit 2", form))
}

// configArgsValidator is the shared custom [cobra.PositionalArgs]
// validator for all three config forms (T0-(l) consolidation —
// validateConfigShowArgs/Get/Set folded into one closure). It runs
// the form-specific base validator (NoArgs for bare, ExactArgs(N)
// for get/set) and, on failure with --json active, emits the
// Envelope-konformen reject on stdout BEFORE returning the error to
// Cobra (Spec §1841/§1842 — `cobra.ExactArgs` alone would fire its
// raw stderr error before RunE and the consumer would get no
// envelope). Pattern-Erbe remove's validateRemoveArgs.
//
// For bare `config`, cobra.NoArgs yields the `unknown command "<x>"`
// error when an unrecognised non-subcommand token is passed
// (R2-HIGH-3: `u-boot config foo` dispatches to the parent with
// args=["foo"]). All three map to Exit 2 via [isUsageError]
// (`unknown command` / `accepts ` prefix).
func configArgsValidator(a *App, subcommand string, base cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if err := base(cmd, args); err != nil {
			if a.json {
				dryRun, _ := cmd.Flags().GetBool("dry-run")
				diffFlag, _ := cmd.Flags().GetBool("diff")
				_ = writeErrorEnvelopeSub(cmd.OutOrStdout(), err, nil, dryRun, diffFlag, "config", subcommand, mapConfigErrorToDiagnostic, nil)
			}
			return err
		}
		return nil
	}
}

// ErrDryRunNotApplicable is returned when `--dry-run` (or `--diff`)
// is given to a read-only config form — `u-boot config` (Show) or
// `u-boot config get` (slice-v1-cli-json-dry-run-config T0-(g)
// Option (i.a)). Only the modifying form `config set` carries the
// preview flags (Cluster-Plan Z. 91-100). Pattern-Erbe logs'
// [ErrFollowJSONNotSupported]: the flag is registered on the command
// so Cobra parses it cleanly, and the RunE rejects it Envelope-
// konform (LH-NFA-USE-004 §1841). Maps to Exit 2 via [isUsageError].
var ErrDryRunNotApplicable = errors.New("--dry-run/--diff is only valid for `config set`")

// runConfigShow streams the full u-boot.yaml body (bare `config`,
// subcommand "show"). Read-only → rejects --dry-run/--diff.
func runConfigShow(ctx context.Context, out io.Writer, flags configShowFlags, uc driving.ConfigUseCase, getwd func() (string, error)) error {
	mapErr := mapConfigErrorToDiagnostic
	if flags.DryRun || flags.Diff {
		return reportErrorSub(out, ErrDryRunNotApplicable, nil, false, false, flags.JSON, "config", "show", mapErr, nil)
	}
	cwd, err := getwd()
	if err != nil {
		return reportErrorSub(out, fmt.Errorf("determine working directory: %w", err), nil, false, false, flags.JSON, "config", "show", mapErr, nil)
	}
	resp, err := uc.Show(ctx, driving.ConfigShowRequest{BaseDir: cwd})
	if err != nil {
		return reportErrorSub(out, sanitizeBaseDir(err, cwd), nil, false, false, flags.JSON, "config", "show", mapErr, nil)
	}
	if flags.JSON {
		data := configShowData{Body: string(resp.Body)}
		return writeEnvelope(out, newDataEnvelope("config", "show", data, nil, 0))
	}
	if _, err := out.Write(resp.Body); err != nil {
		return fmt.Errorf("write config body to stdout: %w", err)
	}
	return nil
}

// runConfigGet prints a single scalar (subcommand "get"). Read-only
// → rejects --dry-run/--diff.
func runConfigGet(
	ctx context.Context,
	out io.Writer,
	args []string,
	flags configGetFlags,
	uc driving.ConfigUseCase,
	getwd func() (string, error),
) error {
	mapErr := mapConfigErrorToDiagnostic
	if flags.DryRun || flags.Diff {
		return reportErrorSub(out, ErrDryRunNotApplicable, nil, false, false, flags.JSON, "config", "get", mapErr, nil)
	}
	path, err := domain.NewConfigPath(args[0])
	if err != nil {
		return reportErrorSub(out, fmt.Errorf("%w: %v", driving.ErrConfigPathUnknown, err), nil, false, false, flags.JSON, "config", "get", mapErr, nil)
	}
	cwd, err := getwd()
	if err != nil {
		return reportErrorSub(out, fmt.Errorf("determine working directory: %w", err), nil, false, false, flags.JSON, "config", "get", mapErr, nil)
	}
	resp, err := uc.Get(ctx, driving.ConfigGetRequest{BaseDir: cwd, Path: path})
	if err != nil {
		return reportErrorSub(out, sanitizeBaseDir(err, cwd), nil, false, false, flags.JSON, "config", "get", mapErr, nil)
	}
	if flags.JSON {
		data := configGetData{Path: path.String(), Value: resp.Value}
		return writeEnvelope(out, newDataEnvelope("config", "get", data, nil, 0))
	}
	fmt.Fprintln(out, resp.Value)
	return nil
}

// runConfigSet parses the path, dispatches to the use case under the
// flag-derived PreviewMode, and renders the outcome (subcommand
// "set"). Modifying form → carries --dry-run/--diff. Error path
// flows through [reportErrorSub] with [sanitizeBaseDir] (Path-Leak-
// Defense, T0-(p)).
func runConfigSet(
	ctx context.Context,
	out io.Writer,
	args []string,
	flags configSetFlags,
	uc driving.ConfigUseCase,
	getwd func() (string, error),
) error {
	mapErr := mapConfigErrorToDiagnostic

	path, err := domain.NewConfigPath(args[0])
	if err != nil {
		return reportErrorSub(out, fmt.Errorf("%w: %v", driving.ErrConfigPathUnknown, err), nil, flags.DryRun, flags.Diff, flags.JSON, "config", "set", mapErr, nil)
	}

	// Spec §714-717: --allow-external-feature-sources is only valid
	// on devcontainer.featureSources.allow (Pre-UC-Validation, T0-(i)).
	if len(flags.AllowExternalFeatureSources) > 0 &&
		path.Kind != domain.ConfigDevcontainerFeatureSourcesAllow {
		err := fmt.Errorf(
			"%w: --allow-external-feature-sources is only valid for `config set devcontainer.featureSources.allow` (Spec §714-717); got path %s",
			driving.ErrConfigPathUnknown, path)
		return reportErrorSub(out, err, nil, flags.DryRun, flags.Diff, flags.JSON, "config", "set", mapErr, nil)
	}

	cwd, err := getwd()
	if err != nil {
		return reportErrorSub(out, fmt.Errorf("determine working directory: %w", err), nil, flags.DryRun, flags.Diff, flags.JSON, "config", "set", mapErr, nil)
	}

	mode := previewModeFromFlags(flags.DryRun, flags.Diff)
	resp, setErr := uc.Set(ctx, driving.ConfigSetRequest{
		BaseDir:                     cwd,
		Path:                        path,
		Value:                       args[1],
		AllowExternalFeatureSources: flags.AllowExternalFeatureSources,
		PreviewMode:                 mode,
		SilenceLogger:               flags.JSON,
	})
	if setErr != nil {
		return reportErrorSub(out, sanitizeBaseDir(setErr, cwd), resp.PlannedFiles, flags.DryRun, flags.Diff, flags.JSON, "config", "set", mapErr, nil)
	}

	if flags.JSON {
		return writeConfigSetJSON(out, resp, path, flags)
	}

	if flags.Diff {
		if err := writeDiff(out, resp.PlannedFiles); err != nil {
			return err
		}
	}
	printConfigSetSummary(out, resp)
	return nil
}

// writeConfigSetJSON renders the success-path envelope for
// `config set`. Two shapes (Cluster-Pattern remove writeRemoveJSON):
//
//   - plain --json (no preview flag) → [newDataEnvelope] with the
//     configSetData carrier.
//   - --dry-run / --diff → [newFullEnvelope] with plannedFiles/
//     changes (from the recorder, R-T4-1) + the same data carrier;
//     hunks in the --diff path. NoOp → empty plannedFiles (T0-(d)).
//
// resp.Warnings (Orphan-Feature-WARN, T0-(n)) map into diagnostics[]
// via [mapWarningsToDiagnostics] — warn-only keeps exit 0.
func writeConfigSetJSON(out io.Writer, resp driving.ConfigSetResponse, path domain.ConfigPath, flags configSetFlags) error {
	data := configSetData{
		Path:     path.String(),
		OldValue: resp.OldValue,
		NewValue: resp.NewValue,
		NoOp:     resp.OldValue == resp.NewValue,
	}
	if path.Kind == domain.ConfigDevcontainerFeatureSourcesAllow && len(flags.AllowExternalFeatureSources) > 0 {
		data.AppendedSources = flags.AllowExternalFeatureSources
	}
	warnDiags := mapWarningsToDiagnostics(resp.Warnings)
	if !flags.DryRun && !flags.Diff {
		return writeEnvelope(out, newDataEnvelope("config", "set", data, warnDiags, 0))
	}
	pfs, chs := mapPlannedFilesToWire(resp.PlannedFiles, flags.Diff)
	return writeEnvelope(out, newFullEnvelope("config", "set", flags.DryRun, flags.Diff, pfs, chs, data, warnDiags, 0))
}

// mapConfigErrorToDiagnostic maps a config-path error to a
// [diagnosticItem] with the spec-konforme LH-Kennung per T0-(f)
// Switch-Order-Pflicht.
//
// Switch-Order verbindlich: FS-Sentinel FIRST
// (driving.ErrConfigFileSystem → LH-NFA-REL-003) so a Multi-`%w`-Wrap
// carrying FS + a validation sentinel falls on the technical-
// persistence class (Exit 14). Then schema/post-patch-sanity, then
// the path/write/value/not-set fachlich sentinels, then
// ErrProjectNotInitialized (Environment-Operation Pattern-Erbe
// up/down/generate/logs → LH-FA-INIT-001, NOT LH-FA-ADD-001), then
// the CLI-form reject (ErrDryRunNotApplicable → Exit 2), then Default.
//
// LH-FA-CONF-005 is multi-use (Path-Unknown / Write-Rejected /
// Value-Not-Set, T0-(m)/R3-MED-2): consumers disambiguate by the
// sentinel-class message prefix, not by `code` alone (all Exit 10).
func mapConfigErrorToDiagnostic(err error) diagnosticItem {
	switch {
	case errors.Is(err, driving.ErrConfigFileSystem):
		return diagnosticItem{Level: "error", Code: "LH-NFA-REL-003", Message: err.Error()}
	case errors.Is(err, driving.ErrConfigSchemaInvalid):
		return diagnosticItem{Level: "error", Code: "LH-FA-CONF-002", Message: err.Error()}
	case errors.Is(err, driving.ErrConfigPostPatchSanityFailed):
		return diagnosticItem{Level: "error", Code: "LH-FA-CONF-002", Message: err.Error()}
	case errors.Is(err, driving.ErrConfigPathUnknown):
		return diagnosticItem{Level: "error", Code: "LH-FA-CONF-005", Message: err.Error()}
	case errors.Is(err, driving.ErrConfigWriteRejected):
		return diagnosticItem{Level: "error", Code: "LH-FA-CONF-005", Message: err.Error()}
	case errors.Is(err, driving.ErrConfigValueInvalid):
		return diagnosticItem{Level: "error", Code: "LH-FA-CONF-001", Message: err.Error()}
	case errors.Is(err, driving.ErrConfigValueNotSet):
		return diagnosticItem{Level: "error", Code: "LH-FA-CONF-005", Message: err.Error()}
	case errors.Is(err, driving.ErrProjectNotInitialized):
		return diagnosticItem{Level: "error", Code: "LH-FA-INIT-001", Message: err.Error()}
	case errors.Is(err, ErrDryRunNotApplicable):
		return diagnosticItem{Level: "error", Code: "LH-FA-CLI-006", Message: err.Error()}
	default:
		return diagnosticItem{Level: "error", Code: "LH-FA-CLI-006", Message: err.Error()}
	}
}

// printConfigSetSummary writes the post-Set human-mode summary line.
func printConfigSetSummary(out io.Writer, resp driving.ConfigSetResponse) {
	if resp.OldValue == resp.NewValue {
		fmt.Fprintf(out, "config: %s already %s; no changes.\n", resp.Path, resp.NewValue)
		return
	}
	fmt.Fprintf(out, "config: %s %s → %s.\n", resp.Path, summaryValue(resp.OldValue), summaryValue(resp.NewValue))
}

// summaryValue renders one side of the OldValue → NewValue line.
// An empty value (first-time write of an optional field) renders as
// the literal `(unset)` so the user sees a deterministic marker.
func summaryValue(s string) string {
	if s == "" {
		return "(unset)"
	}
	return s
}
