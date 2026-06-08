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

// newConfigCommand builds the `u-boot config` Cobra subcommand
// tree (LH-FA-CONF-001 / LH-FA-CONF-005). Three forms:
//
//   - `u-boot config`              — print full u-boot.yaml (Show).
//   - `u-boot config get <path>`   — print a single scalar.
//   - `u-boot config set <path> <value>` — schema-validated write.
//
// The parent command runs Show when no subcommand is given (Cobra's
// natural `Args: NoArgs + RunE` pattern); `get` and `set` are
// children so help / error messages stay in the canonical
// `u-boot config <verb>` namespace.
func newConfigCommand(a *App) *cobra.Command {
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
ubootYAMLConfig struct AND each path's domain validator
(domain.NewProjectName, strconv.ParseBool). On validation
failure the file stays byte-identical. set is idempotent —
writing a value identical to the current one is a no-op.

Exit codes (LH-FA-CLI-006):
  0   success
  2   CLI validation (unknown subcommand, missing args)
  10  fachlich (no u-boot.yaml; unknown path; invalid value;
      schema-invalid post-patch; value not set on Get)
  14  filesystem error

Examples:
  u-boot config
  u-boot config get project.name
  u-boot config set project.name my-service
  u-boot config get devcontainer.enabled
  u-boot config set devcontainer.enabled true`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runConfigShow(cmd.Context(), cmd.OutOrStdout(), a.configUseCase, a.getwd)
		},
	}

	cmd.AddCommand(newConfigGetCommand(a))
	cmd.AddCommand(newConfigSetCommand(a))
	return cmd
}

func newConfigGetCommand(a *App) *cobra.Command {
	return &cobra.Command{
		Use:   "get <path>",
		Short: "Print a single configuration value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigGet(cmd.Context(), cmd.OutOrStdout(), args, a.configUseCase, a.getwd)
		},
	}
}

// configSetFlags bundles the per-invocation flag state of
// `u-boot config set`. The LH-FA-DEV-003 allowlist seed flag is
// only meaningful when the positional path is
// `devcontainer.featureSources.allow` (Spec §717); the use case
// re-checks before applying.
type configSetFlags struct {
	AllowExternalFeatureSources []string

	// JSON read-through from the App's persistent root flag
	// (slice-v1-cli-json-dry-run-config T2, Pattern-Erbe logs T2).
	// At T2 the field is only populated by the Cobra closure; T5
	// routes the JSON-Mode Voll-Schema envelope path off it.
	JSON bool

	// Quiet read-through from the App's persistent root flag. In
	// JSON mode `--quiet` is a no-op (Cluster-T0-(a) doctor-Pattern:
	// `--quiet --json` is semantically identical to `--json`). T2
	// adds the field; T5 consumes it.
	Quiet bool
}

func newConfigSetCommand(a *App) *cobra.Command {
	flags := &configSetFlags{}
	cmd := &cobra.Command{
		Use:   "set <path> <value>",
		Short: "Set a single configuration value (schema-validated)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// JSON/Quiet read-through from the App's persistent root
			// flags (slice-v1-cli-json-dry-run-config T2). T5 will
			// consume these via mapConfigErrorToDiagnostic and the
			// Voll-Schema envelope path.
			flags.JSON = a.json
			flags.Quiet = a.quiet
			return runConfigSet(cmd.Context(), cmd.OutOrStdout(), args, *flags, a.configUseCase, a.getwd)
		},
	}
	cmd.Flags().StringSliceVar(&flags.AllowExternalFeatureSources, "allow-external-feature-sources", nil,
		"additional URLs to append to devcontainer.featureSources.allow (LH-FA-DEV-003; only valid when <path> is devcontainer.featureSources.allow; cumulative with the positional value).")
	return cmd
}

// ErrDryRunNotApplicable is returned when `--dry-run` (or `--diff`)
// is given to a read-only config form — `u-boot config` (Show) or
// `u-boot config get` (slice-v1-cli-json-dry-run-config T0-(g)
// Option (i.a)). Only the modifying form `config set` carries the
// preview flags (Cluster-Plan Z. 91-100: "nur modifying tragen
// Dry-Run"). Pattern-Erbe logs' [ErrFollowJSONNotSupported]: the
// flag is registered on the command so Cobra parses it cleanly, and
// the RunE rejects it Envelope-konform rather than letting Cobra
// emit a raw `unknown flag` to stderr (LH-NFA-USE-004 §1841).
//
// Maps to Exit-Code 2 (LH-FA-CLI-006 usage class) via
// [isUsageError]. Defined in T2; T5 registers the synthetic flags
// on the bare/get commands, wires the reject path, and adds the
// [isUsageError] branch.
var ErrDryRunNotApplicable = errors.New("--dry-run/--diff is only valid for `config set`")

// runConfigShow streams the full u-boot.yaml body to stdout
// byte-identically (slice-m8-config.md §D5).
func runConfigShow(ctx context.Context, out io.Writer, uc driving.ConfigUseCase, getwd func() (string, error)) error {
	cwd, err := getwd()
	if err != nil {
		return fmt.Errorf("determine working directory: %w", err)
	}
	resp, err := uc.Show(ctx, driving.ConfigShowRequest{BaseDir: cwd})
	if err != nil {
		return err
	}
	if _, err := out.Write(resp.Body); err != nil {
		return fmt.Errorf("write config body to stdout: %w", err)
	}
	return nil
}

// runConfigGet parses the path through [domain.NewConfigPath]
// (wrapping any failure in [driving.ErrConfigPathUnknown] for the
// LH-FA-CLI-006 code-10 mapping), invokes the use case, and writes
// the bare scalar value to stdout with a trailing newline
// (slice-m8-config.md §D4).
func runConfigGet(
	ctx context.Context,
	out io.Writer,
	args []string,
	uc driving.ConfigUseCase,
	getwd func() (string, error),
) error {
	path, err := domain.NewConfigPath(args[0])
	if err != nil {
		return fmt.Errorf("%w: %v", driving.ErrConfigPathUnknown, err)
	}
	cwd, err := getwd()
	if err != nil {
		return fmt.Errorf("determine working directory: %w", err)
	}
	resp, err := uc.Get(ctx, driving.ConfigGetRequest{BaseDir: cwd, Path: path})
	if err != nil {
		return err
	}
	fmt.Fprintln(out, resp.Value)
	return nil
}

// runConfigSet parses the path, invokes the use case, and prints a
// short summary line. Two shapes:
//
//   - NoOp (OldValue == NewValue): "config: <path> already <value>; no changes."
//   - Changed:                     "config: <path> <old> → <new>."
//
// The arrow form matches the convention used by `u-boot init`'s
// backup summary so users see one consistent transition glyph
// across subcommands.
func runConfigSet(
	ctx context.Context,
	out io.Writer,
	args []string,
	flags configSetFlags,
	uc driving.ConfigUseCase,
	getwd func() (string, error),
) error {
	path, err := domain.NewConfigPath(args[0])
	if err != nil {
		return fmt.Errorf("%w: %v", driving.ErrConfigPathUnknown, err)
	}

	// Spec §714-717: --allow-external-feature-sources is only
	// valid on the three listed paths. For `config set`, the
	// only valid host is devcontainer.featureSources.allow.
	if len(flags.AllowExternalFeatureSources) > 0 &&
		path.Kind != domain.ConfigDevcontainerFeatureSourcesAllow {
		return fmt.Errorf(
			"%w: --allow-external-feature-sources is only valid for `config set devcontainer.featureSources.allow` (Spec §714-717); got path %s",
			driving.ErrConfigPathUnknown, path)
	}

	cwd, err := getwd()
	if err != nil {
		return fmt.Errorf("determine working directory: %w", err)
	}
	resp, err := uc.Set(ctx, driving.ConfigSetRequest{
		BaseDir:                     cwd,
		Path:                        path,
		Value:                       args[1],
		AllowExternalFeatureSources: flags.AllowExternalFeatureSources,
	})
	if err != nil {
		return err
	}
	printConfigSetSummary(out, resp)
	return nil
}

// printConfigSetSummary writes the post-Set summary line.
func printConfigSetSummary(out io.Writer, resp driving.ConfigSetResponse) {
	if resp.OldValue == resp.NewValue {
		fmt.Fprintf(out, "config: %s already %s; no changes.\n", resp.Path, resp.NewValue)
		return
	}
	fmt.Fprintf(out, "config: %s %s → %s.\n", resp.Path, summaryValue(resp.OldValue), summaryValue(resp.NewValue))
}

// summaryValue renders one side of the OldValue → NewValue line.
// An empty value (first-time write of an optional field, where
// OldValue is `""`) renders as the literal `(unset)` so the user
// sees a deterministic marker instead of an empty pair of spaces.
func summaryValue(s string) string {
	if s == "" {
		return "(unset)"
	}
	return s
}
