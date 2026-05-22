package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// initFlags bundles the per-invocation flag state of `u-boot init`
// (local flags plus the read-through from the persistent --yes /
// --no-interactive on the root command). Kept as a struct so
// [runInit] has one parameter instead of six bool arguments.
type initFlags struct {
	SkipGit        bool
	Force          bool
	Backup         bool
	AssumeExisting bool
	Yes            bool
	NoInteractive  bool
}

// newInitCommand builds the `u-boot init` Cobra subcommand.
//
// Local flags (LH-FA-INIT-005 / LH-FA-CLI-005A):
//
//	[name]              positional, optional — explicit project name
//	                    (LH-FA-INIT-002); when omitted, derived from
//	                    the working directory's basename.
//	--no-git            skip the `git init` step (LH-FA-INIT-007).
//	--force             managed-block-only re-write of existing files
//	                    (LH-FA-INIT-005 §609 / §613).
//	--backup            backup-then-full-overwrite of existing files
//	                    (LH-FA-INIT-005 §605 / §607).
//	--assume-existing   accept implicit existing-project detection
//	                    in non-interactive runs (LH-FA-CLI-005A §238);
//	                    no-op until the M4 soft-detection slice lands.
//
// The persistent flags --yes / --no-interactive are bound at the
// root command (LH-FA-CLI-005A); we read their parsed values via
// the App struct.
func newInitCommand(a *App) *cobra.Command {
	flags := &initFlags{}

	cmd := &cobra.Command{
		Use:   "init [name]",
		Short: "Initialize a new u-boot project in the current directory",
		Long: `Create the mandatory project structure (LH-FA-INIT-003) plus a
u-boot.yaml (LH-FA-CONF-002) and optionally a git repository
(LH-FA-INIT-007). The current working directory is used as BaseDir;
the project name is derived from its basename (LH-FA-INIT-002) unless
a [name] argument is given.

Re-running init on an existing project requires --force (managed-block
only edit) or --backup (full overwrite with safety copy), per
LH-FA-INIT-005 §611–§619.

Examples:
  u-boot init                            # name from current directory
  u-boot init my-service                 # explicit name
  u-boot init --no-git                   # skip git init
  u-boot init --force                    # refresh u-boot blocks only
  u-boot init --backup                   # full overwrite with .bak[*]
  u-boot init --force --backup           # block edit + safety backup
  u-boot init --no-interactive --force   # CI-safe re-init`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Read-through persistent flags from the App; Cobra has
			// already parsed them by the time RunE fires.
			flags.Yes = a.yes
			flags.NoInteractive = a.noInteractive
			return runInit(cmd.Context(), cmd.OutOrStdout(), args, *flags, a.initUseCase, a.getwd)
		},
	}

	cmd.Flags().BoolVar(&flags.SkipGit, "no-git", false, "skip `git init` (LH-FA-INIT-007)")
	cmd.Flags().BoolVar(&flags.Force, "force", false,
		"replace U-BOOT MANAGED BLOCK content in existing files (LH-FA-INIT-005)")
	cmd.Flags().BoolVar(&flags.Backup, "backup", false,
		"back up existing files to <name>.bak[.N] before overwriting (LH-FA-INIT-005)")
	cmd.Flags().BoolVar(&flags.AssumeExisting, "assume-existing", false,
		"accept implicit existing-project detection in non-interactive runs (LH-FA-CLI-005A; no-op until M4 soft-detection)")
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
	flags initFlags,
	uc driving.InitProjectUseCase,
	getwd func() (string, error),
) error {
	if flags.Yes && flags.NoInteractive {
		return ErrConflictingModeFlags
	}

	cwd, err := getwd()
	if err != nil {
		return fmt.Errorf("determine working directory: %w", err)
	}

	req := driving.InitProjectRequest{
		BaseDir:        cwd,
		SkipGit:        flags.SkipGit,
		Force:          flags.Force,
		Backup:         flags.Backup,
		AssumeExisting: flags.AssumeExisting,
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
// of what was created and what was backed up. Order follows
// resp.Created and resp.Backups (which the application service
// guarantees).
func printInitSummary(out io.Writer, resp driving.InitProjectResponse) {
	fmt.Fprintf(out, "Initialized u-boot project %q.\n\nCreated:\n", resp.Project.Name)
	for _, entry := range resp.Created {
		fmt.Fprintln(out, "  - "+entry)
	}
	if len(resp.Backups) > 0 {
		fmt.Fprintln(out, "\nBackups:")
		for _, b := range resp.Backups {
			fmt.Fprintf(out, "  - %s -> %s\n", b.Original, b.Backup)
		}
	}
}
