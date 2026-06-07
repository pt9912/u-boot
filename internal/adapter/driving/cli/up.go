package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"

	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// upFlags bundles the per-invocation flag state of `u-boot up`.
// TimeoutSec is the raw seconds value from `--timeout`; runUp
// converts it to `time.Duration` via `time.Duration(sec) * time.Second`
// — never as raw `time.Duration(sec)` (which would be nanoseconds).
// JSON is read through from the App's persistent root flag so
// `u-boot --json up` and `u-boot up --json` behave identically
// (slice-v1-cli-json-dry-run-up-down T5).
type upFlags struct {
	TimeoutSec int
	Quiet      bool
	JSON       bool
}

// serviceStatus is the per-service wire shape carried in
// [upStatusData.Services]. Fields match the LH-FA-UP-003 four-column
// status table (Name / ContainerStatus / Port / Healthcheck);
// `Port` and `Healthcheck` use omitempty because services without an
// exposed port or healthcheck render as "-" in the human-mode table
// — Key-Absence is the JSON-equivalent (T0-(g) R3-MED-3 plain-String
// + omitempty, NICHT Pointer-Wrap weil keine Three-State-Sit).
type serviceStatus struct {
	Name        string `json:"name"`
	State       string `json:"state"`
	Port        string `json:"port,omitempty"`
	Healthcheck string `json:"healthcheck,omitempty"`
}

// upStatusData is the typed `data` carrier for the `--json` envelope
// of `u-boot up` (slice-v1-cli-json-dry-run-up-down T0-(g)/(j)). Two
// notable fields:
//
//   - Services: NO omitempty — empty slice MUST serialize as `[]`,
//     not `null`. CLI-Layer initializes nil-Slice to
//     `[]serviceStatus{}` before passing to newDataEnvelope so the
//     Empty-Array-Pin (T0-(j) + R5-LOW-3) holds for fire-and-forget
//     (`--timeout=0` returns `Services: nil` from the application).
//   - TimeoutFireAndForget: `*bool` with omitempty — present only in
//     fire-and-forget mode (`--timeout=0`). Pointer-wrap mirrors
//     remove's `VolumesPurged *bool` Key-Absence-Disambiguation
//     (T0-(j) R4 marker discipline).
type upStatusData struct {
	Services             []serviceStatus `json:"services"`
	TimeoutFireAndForget *bool           `json:"timeoutFireAndForget,omitempty"`
}

// newUpCommand builds the `u-boot up` Cobra subcommand
// (LH-FA-UP-001..003).
//
// Local flag:
//
//	--timeout <sec>   maximum wall-clock duration to wait for every
//	                  declared service to stabilize (default 60).
//	                  `0` short-circuits to fire-and-forget (no
//	                  polling, no port/healthcheck probes — see
//	                  LH-FA-UP-001 §970). Negative values are rejected
//	                  with [ErrInvalidTimeout] (exit code 2).
//
// The persistent flags --quiet (suppresses the status table) /
// --verbose / --debug (LH-FA-CLI-005) are read from the App after
// Cobra parses them. The Compose progress stream (LH-NFA-PERF-002)
// always reaches stderr regardless of --quiet; in --json mode it is
// silenced via UpRequest.SilenceProgress (slice-v1-cli-json-dry-run-
// up-down T0-(c) form (d)).
func newUpCommand(a *App) *cobra.Command {
	flags := &upFlags{}

	cmd := &cobra.Command{
		Use:   "up",
		Short: "Start the Compose environment and wait for every service to stabilize",
		Long: `Bring the Compose environment defined in compose.yaml up via
docker compose up -d, then poll docker compose ps every 500ms until
every declared service reaches healthy (when a healthcheck is defined)
or running (when no healthcheck). LH-FA-UP-001 §966-§969 stabilization.

Stabilization semantics per LH-FA-UP-001:
  - --timeout <sec>      maximum wait (default 60). Negative ⇒ exit 2.
  - --timeout=0          fire-and-forget. No polling, no probes;
                         status table omitted, info diagnostic shown
                         (LH-FA-UP-001 §970).

LH-FA-CLI-006 exit codes:
  - 10  no u-boot.yaml or compose.yaml in the current directory
  - 11  Docker daemon unreachable / compose plugin missing
  - 12  Compose runtime failure or stabilization timeout
  - 14  filesystem read failure (u-boot.yaml / compose.yaml)

LH-NFA-PERF-002 progress: pull/create/start/healthcheck phases stream
to stderr live (unaffected by --quiet; silenced in --json mode).`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			flags.Quiet = a.quiet
			flags.JSON = a.json
			return runUp(cmd.Context(), cmd.OutOrStdout(), cmd.ErrOrStderr(), *flags, a.upUseCase, a.getwd)
		},
	}
	cmd.Flags().IntVar(&flags.TimeoutSec, "timeout", 60,
		"maximum seconds to wait for stabilization; 0 = fire-and-forget (LH-FA-UP-001 §963/§970)")
	return cmd
}

// runUp is split from the Cobra closure for testability with a fake
// use-case + fake getwd. Returns the CLI-level error; the wiring
// layer maps it to an exit code via [ExitCode].
//
// Two output paths (slice-v1-cli-json-dry-run-up-down T5):
//
//   - Human (no --json) → [renderUpStatus] table + [renderUpDiagnostics]
//     section, both suppressed by --quiet (legacy contract).
//   - --json → Minimal+Data envelope via [newDataEnvelope] with
//     [upStatusData] carrier. --quiet does NOT suppress the envelope
//     (T0-(b) Cluster-T0-(a) doctor-Pattern: --quiet --json is
//     semantically identical to --json).
//
// Error path: errors flow through [reportError] with [sanitizeBaseDir]
// applied to keep absolute filesystem paths out of diagnostic.message
// (R2-MED-5 path-leak defense). `data` is nil on the error path
// (R2-MED-4: interface-nil, NOT zero-value-struct, otherwise
// `services: null` would break the empty-array-pin).
func runUp(ctx context.Context, stdout, stderr io.Writer, flags upFlags, useCase driving.UpUseCase, getwd func() (string, error)) error {
	mapErr := mapUpErrorToDiagnostic

	if flags.TimeoutSec < 0 {
		return reportError(stdout, ErrInvalidTimeout, nil, false, false, flags.JSON, "up", mapErr, nil)
	}
	cwd, err := getwd()
	if err != nil {
		return reportError(stdout, fmt.Errorf("determine working directory: %w", err), nil, false, false, flags.JSON, "up", mapErr, nil)
	}
	// Explicit `* time.Second` conversion — slice plan T1
	// timeout-semantics pin. `time.Duration(sec)` alone would be
	// nanoseconds.
	timeout := time.Duration(flags.TimeoutSec) * time.Second
	resp, err := useCase.Up(ctx, driving.UpRequest{
		BaseDir:         cwd,
		Timeout:         timeout,
		ProgressSink:    stderr,
		SilenceProgress: flags.JSON,
	})
	if err != nil {
		return reportError(stdout, sanitizeBaseDir(err, cwd), nil, false, false, flags.JSON, "up", mapErr, nil)
	}
	if flags.JSON {
		return writeUpJSON(stdout, resp, flags.TimeoutSec == 0)
	}
	// LH-FA-CLI-005 + M6 slice §T6 binding contract: `up --quiet`
	// must suppress BOTH the status table AND the diagnostic
	// section so CI scripts can rely on empty stdout for success.
	// The Compose progress stream on stderr is NOT affected
	// (LH-NFA-PERF-002 requires phase visibility regardless).
	if flags.Quiet {
		return nil
	}
	if err := renderUpStatus(stdout, resp.Result.Services); err != nil {
		return fmt.Errorf("render status: %w", err)
	}
	renderUpDiagnostics(stdout, resp.Result.Diagnostics, flags.Quiet)
	return nil
}

// writeUpJSON renders the success-path JSON envelope. Always
// Minimal+Data form (up is read-only — no `--dry-run`/`--diff`,
// Cluster-Slice Z. 464-467). `data.services[]` is populated from
// the use-case response; `data.timeoutFireAndForget` is set only
// in fire-and-forget mode (`--timeout=0`).
//
// Empty-Array-Pin (T0-(j) + R5-LOW-3): nil-Slice MUST serialize as
// `[]`, not `null`. The use-case returns `Result.Services: nil` for
// fire-and-forget; this function initializes to `[]serviceStatus{}`
// in that case so the wire form is consistent.
//
// resp.Warnings (T2 Cluster-Vorlauf) map via [mapWarningsToDiagnostics]
// to `diagnostics[]` with `level: "warn"`. Recreate-Warnings
// detection is a follow-up slice (T0-(k) carveout) — for now
// resp.Warnings is empty on the happy path.
func writeUpJSON(out io.Writer, resp driving.UpResponse, fireAndForget bool) error {
	services := make([]serviceStatus, 0, len(resp.Result.Services))
	for _, s := range resp.Result.Services {
		services = append(services, serviceStatus{
			Name:        s.Name,
			State:       s.ContainerStatus.String(),
			Port:        s.Port,
			Healthcheck: s.Healthcheck,
		})
	}
	data := upStatusData{Services: services}
	if fireAndForget {
		t := true
		data.TimeoutFireAndForget = &t
	}
	warnDiags := mapWarningsToDiagnostics(resp.Warnings)
	env := newDataEnvelope("up", "", data, warnDiags, 0)
	return writeEnvelope(out, env)
}

// mapUpErrorToDiagnostic maps an up-path error to a [diagnosticItem]
// with the spec-konforme LH-Kennung per T0-(e) Switch-Order-Pflicht.
//
// Switch-Order verbindlich (slice plan T0-(e), R3-HIGH-1 Reihenfolge-
// Pin): FS-Sentinel FIRST (driving.ErrUpFileSystem) damit ein
// Multi-`%w`-Wrap mit FS+Docker auf LH-NFA-REL-003/Exit 14 fällt,
// NICHT auf Docker/Exit 11 (Defense-only-Pin gegen synthetisch
// konstruierte Multi-Wraps — heute existiert kein realer Code-Pfad
// der beide Klassen chained, weil readComposeFile vor ComposeUp
// failed). Danach Docker/Compose-Runtime via shared
// [mapComposeRuntimeSentinel] helper, dann up-spezifische
// Sentinels, dann cross-cutting fachlich, dann CLI-form, dann
// Default.
//
// Cross-Slice-Klassen-Pin (R4-MED-2):
// `driving.ErrProjectNotInitialized` mappt hier auf `LH-FA-INIT-001`
// (Pattern-Erbe generate als Environment-Operation), NICHT auf
// `LH-FA-ADD-001` (add/remove als Service-Operations). Konsumenten
// disambiguieren über `(code, exitCode)`-Tupel (T8 §6.7 Doku-Pin).
// nolint:dupl // Per-Subcommand-Mapper-Pattern: mapUp und mapDown
// teilen Rows 6+7 (ComposeFileMissing + ProjectNotInitialized) und
// Default — strukturelle Ähnlichkeit ist bewusst (T0-(e) Tabellen-
// Form), Konsolidierung in einen geteilten Helper würde die per-
// Subcommand-Switch-Order-Klarheit auflösen.
func mapUpErrorToDiagnostic(err error) diagnosticItem {
	switch {
	// Row 1: FS-class first (Multi-`%w`-defense).
	case errors.Is(err, driving.ErrUpFileSystem):
		return diagnosticItem{Level: "error", Code: "LH-NFA-REL-003", Message: err.Error()}
	}
	// Rows 2-3: shared Docker/Compose-runtime via helper.
	if code, matched := mapComposeRuntimeSentinel(err); matched {
		return diagnosticItem{Level: "error", Code: code, Message: err.Error()}
	}
	switch {
	// Row 4: up-only runtime class.
	case errors.Is(err, driving.ErrStabilizationTimeout):
		return diagnosticItem{Level: "error", Code: "LH-FA-UP-001", Message: err.Error()}
	// Row 6: shared fachliche Validierung (also in mapDown).
	case errors.Is(err, driving.ErrComposeFileMissing):
		return diagnosticItem{Level: "error", Code: "LH-FA-UP-001", Message: err.Error()}
	// Row 7: cross-cutting project-init.
	case errors.Is(err, driving.ErrProjectNotInitialized):
		return diagnosticItem{Level: "error", Code: "LH-FA-INIT-001", Message: err.Error()}
	// Row 8: up-only CLI-form.
	case errors.Is(err, ErrInvalidTimeout):
		return diagnosticItem{Level: "error", Code: "LH-FA-CLI-006", Message: err.Error()}
	default:
		return diagnosticItem{Level: "error", Code: "LH-FA-CLI-006", Message: err.Error()}
	}
}
