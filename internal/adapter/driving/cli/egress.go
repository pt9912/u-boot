package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// ErrEgressSubcommandRequired is returned by bare `u-boot egress --json`:
// `egress` is a help parent without data of its own, so it cannot emit an
// envelope for a leaf. Like bare `template --json` it rejects with exit 2
// (usage class via [isUsageError]) and without an envelope.
var ErrEgressSubcommandRequired = errors.New("u-boot egress requires a subcommand (try `u-boot egress list`)")

// ErrEgressPreviewNotApplicable is returned when `--dry-run` or `--diff`
// is given to `egress list`, a read-only query. Maps to exit 2.
var ErrEgressPreviewNotApplicable = errors.New("--dry-run/--diff is not applicable to `egress list` (read-only query)")

// isEgressUsageError reports the two usage-class sentinels of the egress
// commands (exit 2); [isUsageError] takes them as one condition.
func isEgressUsageError(err error) bool {
	return errors.Is(err, ErrEgressSubcommandRequired) || errors.Is(err, ErrEgressPreviewNotApplicable)
}

// egressListFlags bundles the flag state of `u-boot egress list`.
// DryRun/Diff are registered only so Cobra parses them and the RunE can
// reject them with an envelope, as `config list` does.
type egressListFlags struct {
	JSON   bool
	DryRun bool
	Diff   bool
}

// egressHostData is one host of the `data.hosts` array.
type egressHostData struct {
	Host    string   `json:"host"`
	Sources []string `json:"sources"`
}

// egressListData is the typed `data` carrier of `egress list --json`. No
// field is omitempty: `hosts` is an array and the three state fields are
// always present.
type egressListData struct {
	Active         bool             `json:"active"`
	SandboxProfile bool             `json:"sandboxProfile"`
	EgressEnabled  bool             `json:"egressEnabled"`
	CloneSource    bool             `json:"cloneSource"`
	Hosts          []egressHostData `json:"hosts"`
}

// newEgressCommand builds the `u-boot egress` group (LH-FA-DEV-008).
// Children: `list`.
func newEgressCommand(a *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "egress",
		Short: "Inspect the egress allowlist of the sandbox devcontainer",
		Long: `Inspect the egress allowlist of the sandbox devcontainer.

Subcommands:

  u-boot egress list [--json]   show the allowlist with the origin of each host`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if a.json {
				return ErrEgressSubcommandRequired
			}
			return cmd.Help()
		},
	}
	cmd.AddCommand(newEgressListCommand(a))
	return cmd
}

func newEgressListCommand(a *App) *cobra.Command {
	flags := &egressListFlags{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Show the egress allowlist with the origin of each host",
		Long: `Show the egress allowlist that "u-boot generate devcontainer" writes into the
start script of the sandbox devcontainer, with the origin of each host: base,
runtime, docker-apt, feature:<name>, clone or user. A host that several origins
contribute lists each of them.

The command reads u-boot.yaml and the git configuration (.git/config) of the
project only. It needs no network, no Docker and no running container, and
it changes nothing.

It states whether the restriction is active for the configuration: the sandbox
profile (devcontainer.profile) and devcontainer.sandbox.egress.enabled must both
be set. If it is not active, the list is the one that would apply on activation.
A name in the list also allows its subdomains.

The list is computed from the project files.
It is not read from a container: a container keeps the list of the image it was built from.

The output stays under --quiet (it is the result of the command). --json emits
the envelope with data.{active,sandboxProfile,egressEnabled,cloneSource,hosts}.

Exit codes (LH-FA-CLI-006):
  0   success
  2   CLI validation (unknown flag, arguments, --dry-run/--diff)
  10  fachlich (no u-boot.yaml; invalid configuration; a git remote URL that
      cannot be written into the devcontainer)
  14  filesystem error

Examples:
  u-boot egress list
  u-boot egress list --json`,
		Args: jsonArgsValidator(a, "egress", "list", cobra.NoArgs, mapEgressErrorToDiagnostic, false),
		RunE: func(cmd *cobra.Command, _ []string) error {
			flags.JSON = a.json
			return runEgressList(cmd.Context(), cmd.OutOrStdout(), *flags, a.egressListUseCase, a.getwd)
		},
	}
	cmd.Flags().BoolVar(&flags.DryRun, "dry-run", false,
		"not applicable to egress list (read-only query); rejected with exit 2")
	cmd.Flags().BoolVar(&flags.Diff, "diff", false,
		"not applicable to egress list (read-only query); rejected with exit 2")
	return cmd
}

// runEgressList delegates to the use case and renders the result. It
// writes the same bytes with and without `--quiet`: the output is the
// result of a query command (LH-FA-CLI-005).
func runEgressList(
	ctx context.Context,
	out io.Writer,
	flags egressListFlags,
	uc driving.EgressListUseCase,
	getwd func() (string, error),
) error {
	fail := func(err error) error {
		return reportErrorSub(out, err, nil, false, false, flags.JSON, "egress", "list", mapEgressErrorToDiagnostic, nil)
	}
	if flags.DryRun || flags.Diff {
		return fail(ErrEgressPreviewNotApplicable)
	}
	if uc == nil {
		return fail(errors.New("egress list: use case not wired"))
	}
	cwd, err := getwd()
	if err != nil {
		return fail(fmt.Errorf("determine working directory: %w", err))
	}
	resp, err := uc.List(ctx, driving.EgressListRequest{BaseDir: cwd})
	if err != nil {
		return fail(sanitizeBaseDir(err, cwd))
	}
	if flags.JSON {
		return writeEnvelope(out, newDataEnvelope("egress", "list", toEgressListData(resp), nil, 0))
	}
	return renderEgressListText(out, resp)
}

func toEgressListData(resp driving.EgressListResponse) egressListData {
	data := egressListData{
		Active:         resp.Active(),
		SandboxProfile: resp.SandboxProfile,
		EgressEnabled:  resp.EgressEnabled,
		CloneSource:    resp.CloneSource,
		Hosts:          make([]egressHostData, 0, len(resp.Hosts)),
	}
	for _, h := range resp.Hosts {
		sources := h.Sources
		if sources == nil {
			sources = []string{}
		}
		data.Hosts = append(data.Hosts, egressHostData{Host: h.Host, Sources: sources})
	}
	return data
}

// egressInactiveReason names why the restriction is not active.
func egressInactiveReason(resp driving.EgressListResponse) string {
	switch {
	case !resp.SandboxProfile && !resp.EgressEnabled:
		return "devcontainer.profile is not `sandbox` and devcontainer.sandbox.egress.enabled is not true"
	case !resp.SandboxProfile:
		return "devcontainer.profile is not `sandbox`"
	default:
		return "devcontainer.sandbox.egress.enabled is not true"
	}
}

// renderEgressListText writes the human form: the state line, the table
// of hosts with their origins and the notes. Every line is part of the
// result; nothing is a hint that `--quiet` could drop.
func renderEgressListText(out io.Writer, resp driving.EgressListResponse) error {
	var b strings.Builder
	if resp.Active() {
		b.WriteString("Egress restriction: active.\n\n")
	} else {
		fmt.Fprintf(&b, "Egress restriction: not active (%s). The list below is the one that would apply once it is active.\n\n",
			egressInactiveReason(resp))
	}
	tw := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "HOST\tORIGIN")
	for _, h := range resp.Hosts {
		fmt.Fprintf(tw, "%s\t%s\n", h.Host, strings.Join(h.Sources, ", "))
	}
	if err := tw.Flush(); err != nil {
		return fmt.Errorf("flush egress-list writer: %w", err)
	}
	b.WriteString("\nNotes:\n")
	b.WriteString("- A name in the list also allows its subdomains.\n")
	if !resp.CloneSource {
		b.WriteString("- No clone host: there is no git remote 'origin' in .git/config and no devcontainer.sandbox.repository.\n")
	}
	b.WriteString("- The list is computed from u-boot.yaml and the git configuration, as `u-boot generate devcontainer` writes it into the start script.\n")
	b.WriteString("  It is not read from a container: a container keeps the list of the image it was built from (measured with the Dev Container CLI 0.89.0:\n")
	b.WriteString("  stopping and starting it again kept the old list, `devcontainer up --remove-existing-container` applied the new one;\n")
	b.WriteString("  measured on Colima, arm64; editor actions not measured).\n")
	if _, err := io.WriteString(out, b.String()); err != nil {
		return fmt.Errorf("write egress list: %w", err)
	}
	return nil
}

// mapEgressErrorToDiagnostic maps an `egress list` error to its
// diagnostic. Only existing LH codes are emitted: FS first (exit 14), then
// the fachlich classes, then the CLI-form rejects.
func mapEgressErrorToDiagnostic(err error) diagnosticItem {
	switch {
	case errors.Is(err, driving.ErrEgressFileSystem):
		return diagnosticItem{Level: "error", Code: "LH-NFA-REL-003", Message: err.Error()}
	case errors.Is(err, driving.ErrConfigSchemaInvalid):
		return diagnosticItem{Level: "error", Code: "LH-FA-CONF-002", Message: err.Error()}
	case errors.Is(err, domain.ErrInvalidSandboxSetting):
		return diagnosticItem{Level: "error", Code: "LH-FA-DEV-006", Message: err.Error()}
	case errors.Is(err, driving.ErrProjectNotInitialized):
		return diagnosticItem{Level: "error", Code: "LH-FA-INIT-001", Message: err.Error()}
	default:
		return diagnosticItem{Level: "error", Code: "LH-FA-CLI-006", Message: err.Error()}
	}
}
