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
//
// Pass-by-value is correct today (six bools = six bytes); revisit
// the call signature if this grows past roughly ten bool/string
// fields or starts holding slices/maps — at that point a pointer
// receiver or a builder pattern becomes the cheaper option.
type initFlags struct {
	SkipGit        bool
	Force          bool
	Backup         bool
	AssumeExisting bool
	Devcontainer   bool
	Yes            bool
	NoInteractive  bool

	// Template is the external project template name selected via
	// `--template <name>` (LH-FA-TPL-001 / slice-v1-template-init T4).
	// Empty keeps the M3 default-init render path; non-empty
	// dispatches to the TemplateInitUseCase via the
	// InitProjectService delegation introduced in T4.
	Template string

	// AllowExternalFeatureSources is the LH-FA-DEV-003 allowlist seed
	// from `--allow-external-feature-sources <quelle>[,<quelle>...]`
	// (Spec §714). Multi-flag occurrences cumulate, comma-separated
	// values split per Cobra StringSlice. Only meaningful when
	// `--devcontainer` is also set. Slice-v1-devcontainer-features T4.
	AllowExternalFeatureSources []string
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
//	--devcontainer      also write the LH-FA-DEV-001 devcontainer
//	                    files (`.devcontainer/devcontainer.json`
//	                    + `Dockerfile`) and set
//	                    `devcontainer.enabled: true` in u-boot.yaml
//	                    (LH-AK-005). Same --force/--backup discipline
//	                    as M3-templated files: existing devcontainer
//	                    files with an `init` block (e.g. from a prior
//	                    `u-boot generate devcontainer`) re-splice;
//	                    without the marker the call aborts with
//	                    ErrFileExists unless --force --backup is set.
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

Soft-existing-detection (LH-FA-INIT-004): when BaseDir lacks the
hard markers (u-boot.yaml / compose.yaml / .env.example) but already
contains ≥3 LH-FA-INIT-003 structure elements (README.md, CHANGELOG.md,
docs/, scripts/, docker/, .devcontainer/devcontainer.json), the
init refuses to proceed without confirmation:

  - --assume-existing  asserts existence non-interactively (exit 10
                       unless --backup / --force).
  - --no-interactive   skips the detection entirely (deterministic
                       fresh init; per-file collisions may still fail).
  - interactive (default) prompts "[y/N]" and aborts on y.

The mode-flag mutual-exclusion check (--yes + --no-interactive →
exit 2) is unchanged.

Examples:
  u-boot init                            # name from current directory
  u-boot init my-service                 # explicit name
  u-boot init --no-git                   # skip git init
  u-boot init --force                    # refresh u-boot blocks only
  u-boot init --backup                   # full overwrite with .bak[*]
  u-boot init --force --backup           # block edit + safety backup
  u-boot init --no-interactive --force   # CI-safe re-init
  u-boot init --assume-existing --backup # re-init a partial layout
  u-boot init --devcontainer             # LH-AK-005 devcontainer flow
  u-boot init --devcontainer --force --backup  # re-splice over generate output`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Read-through persistent flags from the App; Cobra has
			// already parsed them by the time RunE fires.
			flags.Yes = a.yes
			flags.NoInteractive = a.noInteractive
			return runInit(cmd.Context(), cmd.OutOrStdout(), cmd.ErrOrStderr(), args, *flags, a.initUseCase, a.getwd)
		},
	}

	cmd.Flags().BoolVar(&flags.SkipGit, "no-git", false, "skip `git init` (LH-FA-INIT-007)")
	cmd.Flags().BoolVar(&flags.Force, "force", false,
		"replace U-BOOT MANAGED BLOCK content in existing files (LH-FA-INIT-005)")
	cmd.Flags().BoolVar(&flags.Backup, "backup", false,
		"back up existing files to <name>.bak[.N] before overwriting (LH-FA-INIT-005)")
	cmd.Flags().BoolVar(&flags.AssumeExisting, "assume-existing", false,
		"assert existing project in non-interactive runs; aborts unless --backup/--force (LH-FA-INIT-004, LH-FA-CLI-005A §238)")
	cmd.Flags().BoolVar(&flags.Devcontainer, "devcontainer", false,
		"also generate `.devcontainer/devcontainer.json` + `Dockerfile` and set devcontainer.enabled=true in u-boot.yaml (LH-AK-005)")
	cmd.Flags().StringVar(&flags.Template, "template", "",
		"render the project from an external template instead of the default flow (`u-boot template list` for the catalog; LH-FA-TPL-001 / slice-v1-template-init T4 — fresh-init only, mutex with --devcontainer/--force/--backup)")
	cmd.Flags().StringSliceVar(&flags.AllowExternalFeatureSources, "allow-external-feature-sources", nil,
		"seed devcontainer.featureSources.allow with the given URLs (LH-FA-DEV-003; comma-separated, repeatable). Requires --devcontainer; `--yes` does not substitute (LH-NFA-SEC-004).")
	return cmd
}

// runInit is split from the Cobra closure so it can be unit-tested
// with a fake `getwd` and a fake use-case (LH-FA-ARCH-003,
// LH-FA-INIT-002). Context is taken as the first parameter
// explicitly (instead of via cmd.Context()) so contextcheck can see
// the propagation and so the function is straightforward to test
// without a Cobra command.
//
// Scope of the mode flags after the M4 soft-detection slice
// (LH-FA-CLI-005A §238 / LH-FA-INIT-004 §247):
//   - --yes / --no-interactive — mutual-exclusion check fires here
//     (exit 2). --no-interactive propagates into the request to
//     suppress the soft-detection prompt.
//   - --assume-existing — propagates into the request; the service
//     uses it together with the soft-detection result to decide
//     between fresh init and a [driving.ErrProjectExists] abort.
//
// The `errOut` parameter is kept on the signature for forward
// compatibility (future warnings, JSON-mode hints); it is unused on
// the success path today but the closures in newInitCommand still
// pass cmd.ErrOrStderr() so adding an emit later is a one-line
// change.
func runInit(
	ctx context.Context,
	out io.Writer,
	_ io.Writer,
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
		BaseDir:                     cwd,
		SkipGit:                     flags.SkipGit,
		Force:                       flags.Force,
		Backup:                      flags.Backup,
		AssumeExisting:              flags.AssumeExisting,
		NoInteractive:               flags.NoInteractive,
		Devcontainer:                flags.Devcontainer,
		Template:                    flags.Template,
		AllowExternalFeatureSources: flags.AllowExternalFeatureSources,
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
//
// Intentional information split with the application's progress
// emitter (driven.ProgressPort):
//
//   - PRE-write the application emits "Affected files in <baseDir>"
//     with action labels — that is the *intent* the user sees
//     before any side effect, per LH-FA-INIT-005 §609.
//   - POST-write printInitSummary lists the resolved backup paths
//     (which may have suffix .bak.N when the .bak slot was taken)
//     — that is the *result* the user needs for rollback.
//
// Both layers mention the same files; the duplication is by design.
// The Unicode arrow (→) in the Backups section matches the broader
// project glyph convention (Unicode dashes/arrows over ASCII fall-
// backs) — closes T4c-review finding #6.
func printInitSummary(out io.Writer, resp driving.InitProjectResponse) {
	fmt.Fprintf(out, "Initialized u-boot project %q.\n\nCreated:\n", resp.Project.Name)
	for _, entry := range resp.Created {
		fmt.Fprintln(out, "  - "+entry)
	}
	if len(resp.Backups) > 0 {
		fmt.Fprintln(out, "\nBackups:")
		for _, b := range resp.Backups {
			fmt.Fprintf(out, "  - %s → %s\n", b.Original, b.Backup)
		}
	}
}
