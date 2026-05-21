package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// newInitCommand builds the `u-boot init` Cobra subcommand.
//
// Flags supported in M3-T3:
//
//	[name]      positional, optional — explicit project name
//	            (LH-FA-INIT-002); when omitted, derived from the
//	            working directory's basename.
//	--no-git    skip the `git init` step (LH-FA-INIT-007).
//
// Flags planned for M3-T4 / `slice-m4-soft-existing-detection`:
//
//	--backup            LH-FA-INIT-005
//	--force             LH-FA-INIT-005
//	--no-interactive    LH-FA-CLI-005A
//	--yes               LH-FA-CLI-005A
//	--assume-existing   LH-FA-INIT-004 / LH-FA-CLI-005A
func newInitCommand(a *App) *cobra.Command {
	var skipGit bool

	cmd := &cobra.Command{
		Use:   "init [name]",
		Short: "Initialize a new u-boot project in the current directory",
		Long: `Create the mandatory project structure (LH-FA-INIT-003) plus a
u-boot.yaml (LH-FA-CONF-002) and optionally a git repository
(LH-FA-INIT-007). The current working directory is used as BaseDir;
the project name is derived from its basename (LH-FA-INIT-002) unless
a [name] argument is given.

Examples:
  u-boot init               # name from current directory
  u-boot init my-service    # explicit name
  u-boot init --no-git      # do not initialize a git repository`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(cmd.Context(), cmd.OutOrStdout(), args, skipGit, a.initUseCase, a.getwd)
		},
	}

	cmd.Flags().BoolVar(&skipGit, "no-git", false, "skip `git init` (LH-FA-INIT-007)")
	return cmd
}

// runInit is split from the Cobra closure so it can be unit-tested
// with a fake `getwd` and a fake use-case (LH-FA-ARCH-003,
// LH-FA-INIT-002). Context is taken as the first parameter
// explicitly (instead of via cmd.Context()) so contextcheck can see
// the propagation and so the function is straightforward to test
// without a Cobra command.
func runInit(
	ctx context.Context,
	out io.Writer,
	args []string,
	skipGit bool,
	uc driving.InitProjectUseCase,
	getwd func() (string, error),
) error {
	cwd, err := getwd()
	if err != nil {
		return fmt.Errorf("determine working directory: %w", err)
	}

	req := driving.InitProjectRequest{
		BaseDir: cwd,
		SkipGit: skipGit,
	}
	if len(args) == 1 {
		req.Name = args[0]
	}

	resp, err := uc.Init(ctx, req)
	if err != nil {
		return err
	}

	printInitSummary(out, resp)
	return nil
}

// printInitSummary writes a deterministic, human-friendly summary
// of what was created. Order follows resp.Created (which the
// application service guarantees).
func printInitSummary(out io.Writer, resp driving.InitProjectResponse) {
	fmt.Fprintf(out, "Initialized u-boot project %q.\n\nCreated:\n", resp.Project.Name)
	for _, entry := range resp.Created {
		fmt.Fprintln(out, "  - "+entry)
	}
}
