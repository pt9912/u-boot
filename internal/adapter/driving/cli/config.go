package cli

import (
	"context"
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

func newConfigSetCommand(a *App) *cobra.Command {
	return &cobra.Command{
		Use:   "set <path> <value>",
		Short: "Set a single configuration value (schema-validated)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigSet(cmd.Context(), cmd.OutOrStdout(), args, a.configUseCase, a.getwd)
		},
	}
}

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
	resp, err := uc.Set(ctx, driving.ConfigSetRequest{
		BaseDir: cwd, Path: path, Value: args[1],
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
