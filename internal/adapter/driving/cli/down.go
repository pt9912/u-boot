package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// downFlags bundles the per-invocation flag state of `u-boot down`.
// The local flag Volumes is CLI-only; Yes / NoInteractive / Quiet
// are read through from the App's persistent root flags
// (LH-FA-CLI-005 / LH-FA-CLI-005A). The spec LH-FA-CLI-005A explicitly names
// `u-boot down --volumes` among the commands governed by the
// persistent `--yes` / `--no-interactive` switches, so reusing the
// root values keeps `u-boot --yes down --volumes` working
// identically to `u-boot down --volumes --yes`. JSON read through
// from the root persistent flag (slice-v1-cli-json-dry-run-up-down
// T5).
type downFlags struct {
	Volumes       bool
	Yes           bool
	NoInteractive bool
	Quiet         bool
	JSON          bool
}

// downStatusData is the typed `data` carrier for the `--json`
// envelope of `u-boot down` (slice-v1-cli-json-dry-run-up-down
// T0-(h)). Single field `removedVolumes bool` (NO omitempty, R5-LOW-1
// + R6-MED-1): `false` is the legitimate success value "nothing
// removed" — the consumer must be able to distinguish it from
// key-absence. Pattern matches today's `DownResponse.RemovedVolumes
// bool` port contract (`port/driving/down.go:80`); named-list form
// (`[]string`) is out-of-scope V1 carveout (Folge-Slice
// `slice-v1-down-volumes-named-list`).
type downStatusData struct {
	RemovedVolumes bool `json:"removedVolumes"`

	// RemovedVolumeNames lists the Docker volumes that `down --volumes`
	// removed (slice-v1-down-volumes-named-list). Additive; always an
	// array (`[]` when none or not determinable), never null.
	RemovedVolumeNames []string `json:"removedVolumeNames"`
}

// newDownCommand builds the `u-boot down` Cobra subcommand
// (LH-FA-UP-004).
//
// Local flags:
//
//	--volumes  remove named Compose volumes alongside containers
//	           (LH-FA-UP-004 destructive op; default false).
//	--yes      auto-confirm the destructive --volumes prompt
//	           (LH-FA-CLI-005A / LH-FA-CLI-005A).
//
// The persistent flags --no-interactive (LH-FA-CLI-005A / LH-FA-CLI-005A)
// and --quiet (LH-FA-CLI-005) are read from the App after Cobra
// parses them. `--json` triggers refuse-by-default for the
// destructive `--volumes`-gate (T0-(d) Option (b) Request-time Gate-
// Branch): without `--yes` the use case returns
// [ErrConfirmationRequired]/Exit 10.
//
// Mode-flag mutual exclusion: `--yes` AND `--no-interactive` set
// together returns [ErrConflictingModeFlags] (exit code 2,
// LH-FA-CLI-005A). Independent of whether --volumes is set —
// the LH-FA-CLI-005A rule is a global CLI-validation, not a destructive-path
// concern.
func newDownCommand(a *App) *cobra.Command {
	flags := &downFlags{}

	cmd := &cobra.Command{
		Use:   "down",
		Short: "Stop the Compose environment and optionally remove its named volumes",
		Long: `Tear down the Compose environment via docker compose down.
With --volumes the named Compose volumes are removed as well (data
loss; LH-FA-UP-004).

Destructive confirmation gate (LH-FA-CLI-005A):
  - --yes                     auto-confirm
  - --no-interactive          fail-fast with exit 10 (no confirmation
                              possible without user input)
  - interactive (default)     "[y/N]" prompt; "n" / empty / EOF aborts
  - --json                    refuse-by-default (no prompt on stdin);
                              user must opt in via --yes for --volumes

The mode-flag mutual exclusion (--yes + --no-interactive → exit 2)
is checked before any use-case logic runs.

LH-FA-CLI-006 exit codes:
  - 10  no u-boot.yaml / no compose.yaml / destructive confirmation refused
  - 11  Docker daemon unreachable / compose plugin missing
  - 12  Compose runtime failure
  - 14  filesystem read failure (u-boot.yaml / compose.yaml)

LH-NFA-PERF-002 progress: compose down phases stream to stderr live
(unaffected by --quiet; silenced in --json mode).`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			// LH-FA-CLI-005A: --yes / --no-interactive are
			// persistent root flags that govern down --volumes
			// (along with init, add, remove, config set). Read
			// the parsed values into the per-invocation struct
			// so runDown stays unit-testable without poking at
			// global state.
			flags.Yes = a.yes
			flags.NoInteractive = a.noInteractive
			flags.Quiet = a.quiet
			flags.JSON = a.json
			return runDown(cmd.Context(), cmd.OutOrStdout(), cmd.ErrOrStderr(), *flags, a.downUseCase, a.getwd)
		},
	}
	cmd.Flags().BoolVar(&flags.Volumes, "volumes", false,
		"also remove named Compose volumes (data loss; LH-FA-UP-004)")
	return cmd
}

// runDown is split from the Cobra closure for testability.
//
// Two output paths (slice-v1-cli-json-dry-run-up-down T5):
//
//   - Human (no --json) → [renderDownSuccess] (suppressed by --quiet).
//   - --json → Minimal+Data envelope via [newDataEnvelope] with
//     [downStatusData] carrier (`removedVolumes` bool, no omitempty).
//
// Error path: errors flow through [reportError] with [sanitizeBaseDir]
// applied to keep absolute filesystem paths out of diagnostic.message
// (R2-MED-5 path-leak defense). `data` is nil on the error path
// (R2-MED-4 Call-Site-Pin — interface-nil, NOT zero-value-struct).
func runDown(ctx context.Context, stdout, stderr io.Writer, flags downFlags, useCase driving.DownUseCase, getwd func() (string, error)) error {
	mapErr := mapDownErrorToDiagnostic

	// LH-FA-CLI-005A mode-flag exclusion. Note: this checks
	// the LOCAL down --yes flag against the PERSISTENT root
	// --no-interactive — different fields but same exclusivity
	// rule.
	if flags.Yes && flags.NoInteractive {
		return reportError(stdout, ErrConflictingModeFlags, nil, false, false, flags.JSON, "down", mapErr, nil)
	}
	cwd, err := getwd()
	if err != nil {
		return reportError(stdout, fmt.Errorf("determine working directory: %w", err), nil, false, false, flags.JSON, "down", mapErr, nil)
	}
	resp, err := useCase.Down(ctx, driving.DownRequest{
		BaseDir:          cwd,
		RemoveVolumes:    flags.Volumes,
		AssumeYes:        flags.Yes,
		NonInteractive:   flags.NoInteractive,
		ProgressSink:     stderr,
		SilenceConfirmer: flags.JSON,
	})
	if err != nil {
		return reportError(stdout, sanitizeBaseDir(err, cwd), nil, false, false, flags.JSON, "down", mapErr, nil)
	}
	if flags.JSON {
		return writeDownJSON(stdout, resp)
	}
	renderDownSuccess(stdout, resp.RemovedVolumes, flags.Quiet)
	if !flags.Quiet && len(resp.RemovedVolumeNames) > 0 {
		fmt.Fprintf(stdout, "Removed volumes: %s\n", strings.Join(resp.RemovedVolumeNames, ", "))
	}
	return nil
}

// writeDownJSON renders the success-path JSON envelope. Always
// Minimal+Data form (down is read-only on the local FS — Cluster-
// Slice Z. 464-467 — even though it mutates the Docker daemon
// state). `data.removedVolumes` mirrors the use-case response.
func writeDownJSON(out io.Writer, resp driving.DownResponse) error {
	names := resp.RemovedVolumeNames
	if names == nil {
		names = []string{}
	}
	data := downStatusData{RemovedVolumes: resp.RemovedVolumes, RemovedVolumeNames: names}
	env := newDataEnvelope("down", "", data, nil, 0)
	return writeEnvelope(out, env)
}

// mapDownErrorToDiagnostic maps a down-path error to a
// [diagnosticItem] with the spec-konforme LH-Kennung per T0-(e)
// Switch-Order-Pflicht.
//
// Switch-Order verbindlich (slice plan T0-(e), R3-HIGH-1 Reihenfolge-
// Pin): FS-Sentinel FIRST (driving.ErrDownFileSystem) damit ein
// Multi-`%w`-Wrap mit FS+Docker auf LH-NFA-REL-003/Exit 14 fällt
// (analog mapUp). Danach Docker/Compose-Runtime via shared
// [mapComposeRuntimeSentinel] helper, dann down-spezifische
// Sentinels, dann cross-cutting fachlich, dann CLI-form, dann
// Default.
//
// Cross-Slice-Klassen-Pin (R4-MED-2):
// `driving.ErrProjectNotInitialized` mappt hier auf `LH-FA-INIT-001`
// (Pattern-Erbe generate als Environment-Operation), NICHT auf
// `LH-FA-ADD-001`. Identisch zu mapUp — beide Subcommands sind
// Environment-Operations.
// nolint:dupl // Per-Subcommand-Mapper-Pattern: mapUp und mapDown
// teilen Rows 6+7 (ComposeFileMissing + ProjectNotInitialized) und
// Default — strukturelle Ähnlichkeit ist bewusst (T0-(e) Tabellen-
// Form), Konsolidierung in einen geteilten Helper würde die per-
// Subcommand-Switch-Order-Klarheit auflösen.
func mapDownErrorToDiagnostic(err error) diagnosticItem {
	switch {
	// Row 1: FS-class first (Multi-`%w`-defense).
	case errors.Is(err, driving.ErrDownFileSystem):
		return diagnosticItem{Level: "error", Code: "LH-NFA-REL-003", Message: err.Error()}
	}
	// Rows 2-3: shared Docker/Compose-runtime via helper.
	if code, matched := mapComposeRuntimeSentinel(err); matched {
		return diagnosticItem{Level: "error", Code: code, Message: err.Error()}
	}
	switch {
	// Row 5: down-only Confirmer-Refuse.
	case errors.Is(err, driving.ErrConfirmationRequired):
		return diagnosticItem{Level: "error", Code: "LH-FA-INIT-005", Message: err.Error()}
	// Row 6: shared fachliche Validierung.
	case errors.Is(err, driving.ErrComposeFileMissing):
		return diagnosticItem{Level: "error", Code: "LH-FA-UP-001", Message: err.Error()}
	// Row 7: cross-cutting project-init.
	case errors.Is(err, driving.ErrProjectNotInitialized):
		return diagnosticItem{Level: "error", Code: "LH-FA-INIT-001", Message: err.Error()}
	// Row 9: down-only CLI-form mutex (LH-FA-CLI-005A).
	case errors.Is(err, ErrConflictingModeFlags):
		return diagnosticItem{Level: "error", Code: "LH-FA-CLI-005A", Message: err.Error()}
	default:
		return diagnosticItem{Level: "error", Code: "LH-FA-CLI-006", Message: err.Error()}
	}
}
