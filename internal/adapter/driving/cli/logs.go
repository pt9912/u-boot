package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// logsFlags bundles the per-invocation flag state of `u-boot logs`.
// Slice-v1-logs §T0-Outcomes pinned the surface: only `--follow`
// and `--tail` are exposed (T0-(d) Spec-treu — no
// `--no-log-prefix`/`--timestamps`). Slice-v1-cli-json-dry-run-logs
// T2 adds JSON/Quiet read-through from the App's persistent root
// flags so `u-boot --json logs` and `u-boot logs --json` behave
// identically (T0-(j)(ii) Cluster-Pattern analog up/down).
type logsFlags struct {
	// Follow mirrors `docker compose logs --follow`. Blocks until
	// SIGINT cancels `cmd.Context()`; the application service
	// short-circuits the resulting context.Canceled to a nil error
	// so the CLI exits 0.
	Follow bool

	// Tail is the raw string from `--tail <n>`. Empty means "flag
	// not set" — the application service normalises that to
	// Compose's `"all"`. Negative or non-numeric inputs (except
	// the internal `"all"`, which is not a user-supplied value)
	// are rejected by [validateLogsTailFlag] before the use case
	// runs.
	Tail string

	// NoLogPrefix / Timestamps / Since / Until are the Compose
	// pass-through knobs of slice-v1-logs-format-flags and
	// slice-v1-logs-time-range-filter (Since/Until validated by
	// [validateLogsTimeFlag]).
	NoLogPrefix bool
	Timestamps  bool
	Since       string
	Until       string

	// JSON read-through from the App's persistent root flag
	// (slice-v1-cli-json-dry-run-logs T0-(j)(ii)). When true,
	// runLogs routes via the Single-Envelope JSON path (T0-(a)
	// Option (A) festgezurrt): `--follow --json` is rejected with
	// [ErrFollowJSONNotSupported] / Exit 2 before the use case
	// runs; bounded-output (`--tail=N --json` or `--json` without
	// `--follow`) buffers the Compose-Stream and emits a
	// Minimal+Data envelope with `data.lines []string`.
	JSON bool

	// Quiet read-through from the App's persistent root flag.
	// Today logs streams Compose-Output to OutputSink directly and
	// `--quiet` does not silence it (Status-quo from
	// slice-v1-logs). In JSON mode `--quiet` is a no-op
	// (Cluster-T0-(a) doctor-Pattern: `--quiet --json` is
	// semantically identical to `--json`).
	Quiet bool
}

// ErrInvalidLogsTail is returned by `u-boot logs` when `--tail` is
// neither empty nor a non-negative integer string. Slice-v1-logs
// §T0-(c) + §AK: numeric ≥ 0 accepted, everything else rejected at
// the CLI Stage-1 with Exit-Code 2 — the value never reaches the
// application service. Lives in the cli package because the
// LH-FA-CLI-006 mapping to exit code 2 is a CLI concern.
var ErrInvalidLogsTail = errors.New("--tail must be a non-negative integer")

// ErrInvalidLogsTime is returned by `u-boot logs` when `--since` /
// `--until` is neither a positive relative duration (`90s`, `30m`,
// `1h`) nor an absolute timestamp (RFC 3339, or `YYYY-MM-DD`,
// `YYYY-MM-DDTHH:MM[:SS]`). Maps to LH-FA-CLI-006 exit code 2.
var ErrInvalidLogsTime = errors.New("--since/--until must be a duration like 30m or a timestamp like 2026-06-07T12:00:00Z")

// ErrFollowJSONNotSupported is returned by `u-boot logs --follow
// --json` (slice-v1-cli-json-dry-run-logs T0-(a) Option (A)
// festgezurrt): the unbounded streaming use case cannot be
// reconciled with the LH-NFA-USE-004 §1841 Single-Envelope-pro-
// Aufruf-Vertrag. Konsumenten die strukturiert streamen wollen
// können `--tail=N --json` für Bounded-Snapshots nutzen ODER
// einen Folge-Slice mit NDJSON-Carveout abwarten. Maps to Exit-
// Code 2 (LH-FA-CLI-006-Klasse) via [isUsageError].
var ErrFollowJSONNotSupported = errors.New("--follow is not supported with --json (use --tail=N --json for bounded snapshots)")

// logsStatusData is the typed `data` carrier for the `--json`
// envelope of `u-boot logs` (slice-v1-cli-json-dry-run-logs T5,
// T0-(a) Option (A)). Compose-Output wird im CLI-Layer in einen
// bytes.Buffer gepuffert; nach UC-Return wird der Buffer per
// Newline gesplittet und in `data.lines` serialisiert.
//
// **Lines: NO omitempty** — leerer Stream MUSS als `[]`
// serialisieren, nicht `null`. CLI-Layer initialisiert nil-Slice
// mit `[]string{}` bevor er das Envelope baut. Pattern-Erbe
// up-down T0-(j) Empty-Array-Pin.
type logsStatusData struct {
	Lines []string `json:"lines"`
}

// newLogsCommand builds the `u-boot logs [service]` Cobra
// subcommand (LH-FA-UP-005). Slice-v1-logs §AK contract:
//
//   - Positional argument: optional single service name; validated
//     via [domain.NewServiceName] (regex-only — Compose checks
//     runtime existence and surfaces a Compose-Runtime-Error if the
//     service is unknown). Format failures map to Exit-Code 10 via
//     [isServiceValidationError].
//   - `--follow` (default false). Blocks on `cmd.Context()`;
//     SIGINT cancels and the application service returns nil so
//     the CLI exits 0.
//   - `--tail <n>` (default empty → Compose-Default "all" after
//     normalisation in the use case). Accepts non-negative integer
//     strings only; otherwise [ErrInvalidLogsTail] / Exit-Code 2.
//
// Output: Compose-Default (with service-prefix, without
// timestamps). The CLI writes Compose stdout/stderr through the
// application service's OutputSink to `cmd.OutOrStdout()` in
// Human-Mode. In `--json` mode (slice-v1-cli-json-dry-run-logs T5,
// T0-(a) Option (A)) the Compose-Stream is buffered and emitted as
// a Minimal+Data envelope after UC-Return; `--follow --json` is
// rejected with [ErrFollowJSONNotSupported] / Exit 2 before the
// use case runs.
func newLogsCommand(a *App) *cobra.Command {
	flags := &logsFlags{}
	cmd := &cobra.Command{
		Use:   "logs [service...]",
		Short: "Stream Compose logs of every service or one selected service",
		Long: `Stream Docker Compose logs for the project's services.

Without a positional argument, all services declared in compose.yaml
are streamed (Compose-Default — no u-boot.yaml filter). With a
positional arguments, only those services stream (one or several, e.g.
"logs api db"). Unknown
services at runtime map to LH-FA-CLI-006 exit code 12 via the
Compose runtime error path.

Flags:
  --follow         stream until Ctrl-C (LH-FA-UP-005); SIGINT exits 0.
                   NOT supported with --json (use --tail=N --json
                   for bounded snapshots).
  --tail <n>       show only the last n lines per service. Default
                   shows all lines (Compose-Default). Negative or
                   non-numeric inputs ⇒ exit 2.
  --no-log-prefix  suppress the "service  |" prefix.
  --timestamps     prefix every line with its timestamp.
  --since <t>      only logs since t: a duration (30m, 1h) or a
                   timestamp (2026-06-07T12:00:00Z). Invalid ⇒ exit 2.
  --until <t>      only logs before t (same formats).

LH-FA-CLI-006 exit codes:
  - 0   success (incl. --follow terminated by SIGINT)
  - 2   --tail with invalid value, or malformed CLI usage
        (--follow --json combo also returns exit 2)
  - 10  no u-boot.yaml / compose.yaml; or invalid service name
        (regex-only, per T0-(b))
  - 11  Docker daemon unreachable / compose plugin missing
  - 12  Compose runtime failure (unknown service at runtime, etc.)
  - 14  filesystem read failure (u-boot.yaml / compose.yaml)`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// JSON/Quiet read-through from the App's persistent
			// root flags (slice-v1-cli-json-dry-run-logs T2/T5).
			flags.JSON = a.json
			flags.Quiet = a.quiet
			return runLogs(cmd.Context(), cmd.OutOrStdout(), cmd.ErrOrStderr(), args, *flags, a.logsUseCase, a.getwd)
		},
	}
	cmd.Flags().BoolVar(&flags.Follow, "follow", false,
		"stream logs continuously until Ctrl-C (LH-FA-UP-005 §1038)")
	cmd.Flags().StringVar(&flags.Tail, "tail", "",
		"show only the last n lines per service (non-negative integer; default = all)")
	cmd.Flags().BoolVar(&flags.NoLogPrefix, "no-log-prefix", false,
		"suppress the service prefix of each line (docker compose logs --no-log-prefix)")
	cmd.Flags().BoolVar(&flags.Timestamps, "timestamps", false,
		"show a timestamp on every line (docker compose logs --timestamps)")
	cmd.Flags().StringVar(&flags.Since, "since", "",
		"only logs since a duration (30m, 1h) or timestamp (2026-06-07T12:00:00Z)")
	cmd.Flags().StringVar(&flags.Until, "until", "",
		"only logs before a duration (30m, 1h) or timestamp (2026-06-07T12:00:00Z)")
	return cmd
}

// runLogs is split from the Cobra closure for testability. Parses
// the positional service argument and the --tail flag, then
// delegates to the LogsUseCase.
//
// Slice-v1-logs §AK pinned the validation order, extended by
// slice-v1-cli-json-dry-run-logs T5 with JSON-Mode-Reject and
// envelope-error-emission via [reportError]:
//
//  1. --follow --json combo Reject (Exit-Code 2 via
//     [ErrFollowJSONNotSupported]) — slice-v1-cli-json-dry-run-logs
//     T0-(a) Option (A) festgezurrt.
//  2. --tail Stage-1 validation (Exit-Code 2 on invalid).
//  3. Service-name format validation via [domain.NewServiceName]
//     (Exit-Code 10 via isServiceValidationError on regex failure).
//     Skipped when no positional argument is given.
//  4. Working-directory probe → BaseDir.
//  5. Use-Case dispatch. In Human-Mode the Compose-Stream goes
//     straight to stdout; in JSON-Mode it is buffered for the
//     Minimal+Data envelope after UC-Return.
//
// Error path: errors flow through [reportError] with
// [sanitizeBaseDir] applied to keep absolute filesystem paths out
// of diagnostic.message (cluster-weite Path-Leak-Defense via
// `cli/sanitize.go` — Helper aus up-down T5). `data` is nil on the
// error path (interface-nil, NOT zero-value-struct, otherwise
// `lines: null` would break the empty-array-pin).
func runLogs(
	ctx context.Context,
	stdout, errOut io.Writer,
	args []string,
	flags logsFlags,
	uc driving.LogsUseCase,
	getwd func() (string, error),
) error {
	mapErr := mapLogsErrorToDiagnostic
	_ = errOut // human-mode renders straight to stdout; reserved for future stderr-routed diagnostics

	// Pre-UC-Validation 1: --follow + --json reject (T0-(a)
	// Option (A) festgezurrt).
	if flags.Follow && flags.JSON {
		return reportError(stdout, ErrFollowJSONNotSupported, nil, false, false, flags.JSON, "logs", mapErr, nil)
	}

	// Pre-UC-Validation 2: --tail format.
	if err := validateLogsTailFlag(flags.Tail); err != nil {
		return reportError(stdout, err, nil, false, false, flags.JSON, "logs", mapErr, nil)
	}

	// Pre-UC-Validation 3: positional service name format.
	for _, f := range []struct{ name, val string }{{"--since", flags.Since}, {"--until", flags.Until}} {
		if err := validateLogsTimeFlag(f.name, f.val); err != nil {
			return reportError(stdout, err, nil, false, false, flags.JSON, "logs", mapErr, nil)
		}
	}

	services, err := parseLogsServices(args)
	if err != nil {
		return reportError(stdout, err, nil, false, false, flags.JSON, "logs", mapErr, nil)
	}

	cwd, err := getwd()
	if err != nil {
		return reportError(stdout, fmt.Errorf("determine working directory: %w", err), nil, false, false, flags.JSON, "logs", mapErr, nil)
	}

	// Output sink selection: JSON-Mode buffers, Human-Mode streams
	// straight to stdout. Buffer-Größe ist proportional zu
	// --tail × durchschnittliche Zeilenlänge; bei Default
	// --tail=all der gesamte Compose-Log-Inhalt. Heute kein
	// expliziter Cap — Konsumenten nutzen typisch --tail=100..1000.
	var sink io.Writer
	var buf *bytes.Buffer
	if flags.JSON {
		buf = &bytes.Buffer{}
		sink = buf
	} else {
		sink = stdout
	}

	_, err = uc.Logs(ctx, driving.LogsRequest{
		BaseDir:     cwd,
		Services:    services,
		Follow:      flags.Follow,
		Tail:        flags.Tail,
		NoLogPrefix: flags.NoLogPrefix,
		Timestamps:  flags.Timestamps,
		Since:       flags.Since,
		Until:       flags.Until,
		OutputSink:  sink,
	})
	if err != nil {
		return reportError(stdout, sanitizeBaseDir(err, cwd), nil, false, false, flags.JSON, "logs", mapErr, nil)
	}

	if flags.JSON {
		return writeLogsJSON(stdout, buf.Bytes())
	}
	return nil
}

// writeLogsJSON emits the success-path Minimal+Data envelope.
// Compose-Output wird per Newline gesplittet und in
// `data.lines []string` serialisiert. Empty-Array-Pin: nil/empty
// Buffer wird zu `[]string{}` initialisiert damit `lines: []`
// (nicht `null`) im JSON erscheint (Cluster-Pattern aus up-down
// T0-(j) R5-LOW-3).
//
// Trailing-Newline-Handling: ein leerer letzter Token von
// strings.Split nach `\n`-Split (Buffer endet typisch mit `\n`)
// wird gestrippt — sonst hätte der letzte Line-Eintrag immer
// `""` als Wert.
func writeLogsJSON(out io.Writer, raw []byte) error {
	lines := splitLogLines(raw)
	if lines == nil {
		lines = []string{}
	}
	data := logsStatusData{Lines: lines}
	env := newDataEnvelope("logs", "", data, nil, 0)
	return writeEnvelope(out, env)
}

// splitLogLines splits the Compose-output buffer on `\n` and
// strips a trailing empty token (the buffer usually ends with a
// newline, which would otherwise produce a phantom empty line).
func splitLogLines(raw []byte) []string {
	if len(raw) == 0 {
		return nil
	}
	s := string(raw)
	s = strings.TrimSuffix(s, "\n")
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

// mapLogsErrorToDiagnostic maps a logs-path error to a
// [diagnosticItem] with the spec-konforme LH-Kennung per T0-(e)
// Switch-Order-Pflicht.
//
// Switch-Order verbindlich (slice plan T0-(e), R3-HIGH-1
// Reihenfolge-Pin): FS-Sentinel FIRST (driving.ErrLogsFileSystem)
// damit ein Multi-`%w`-Wrap mit FS+Docker auf LH-NFA-REL-003/Exit
// 14 fällt. Danach Docker/Compose-Runtime via shared
// [mapComposeRuntimeSentinel] helper, dann fachliche Validierung
// (Compose-File / Project-Init / Service-Name), dann CLI-Form
// (ErrFollowJSONNotSupported / ErrInvalidLogsTail), dann Default.
//
// Cross-Slice-Klassen-Pin (Pattern-Erbe up-down R4-MED-2):
// `driving.ErrProjectNotInitialized` mappt hier auf
// `LH-FA-INIT-001` (Pattern-Erbe generate als Environment-
// Operation, identisch zu up/down).
//
// nolint:dupl // Per-Subcommand-Mapper-Pattern: alle Cluster-
// Mapper teilen Rows 4-5 (ComposeFileMissing + ProjectNotInitialized)
// und Default — strukturelle Ähnlichkeit ist bewusst (T0-(e)
// Tabellen-Form), Konsolidierung würde die per-Subcommand-Switch-
// Order-Klarheit auflösen. (Bestätigungsrunde MED-1: Row-Nummern an
// logs-Mapper-Tabelle angepasst — up/down hat eigene Rows 4-5 wie
// StabilizationTimeout/ConfirmerRefuse, logs nicht.)
func mapLogsErrorToDiagnostic(err error) diagnosticItem {
	switch {
	// Row 1: FS-class first (Multi-`%w`-defense).
	case errors.Is(err, driving.ErrLogsFileSystem):
		return diagnosticItem{Level: "error", Code: "LH-NFA-REL-003", Message: err.Error()}
	}
	// Rows 2-3: shared Docker/Compose-runtime via helper.
	if code, matched := mapComposeRuntimeSentinel(err); matched {
		return diagnosticItem{Level: "error", Code: code, Message: err.Error()}
	}
	switch {
	// Row 4: shared fachliche Validierung (ComposeFileMissing).
	case errors.Is(err, driving.ErrComposeFileMissing):
		return diagnosticItem{Level: "error", Code: "LH-FA-UP-001", Message: err.Error()}
	// Row 5: cross-cutting project-init (LH-FA-INIT-001, identisch
	// zu up/down/generate).
	case errors.Is(err, driving.ErrProjectNotInitialized):
		return diagnosticItem{Level: "error", Code: "LH-FA-INIT-001", Message: err.Error()}
	// Row 6: domain-level service-name validation.
	case errors.Is(err, domain.ErrInvalidServiceName):
		return diagnosticItem{Level: "error", Code: "LH-FA-INIT-006", Message: err.Error()}
	// Row 7: logs-only CLI-form follow+json reject.
	case errors.Is(err, ErrFollowJSONNotSupported):
		return diagnosticItem{Level: "error", Code: "LH-FA-CLI-006", Message: err.Error()}
	// Row 8: logs-only CLI-form tail validation.
	case errors.Is(err, ErrInvalidLogsTail), errors.Is(err, ErrInvalidLogsTime):
		return diagnosticItem{Level: "error", Code: "LH-FA-CLI-006", Message: err.Error()}
	default:
		return diagnosticItem{Level: "error", Code: "LH-FA-CLI-006", Message: err.Error()}
	}
}

// validateLogsTailFlag enforces T0-(c) at the CLI Stage-1: empty
// (flag not set) passes through; otherwise the value must parse as
// a non-negative integer. The internal `"all"` constant is NOT a
// valid user-supplied value — only the application service
// produces it via normaliseTail.
//
// Review-Followup F1: Compose-CLI users tend to type `--tail all`
// out of muscle memory; the special-case explains that the
// implicit default already streams all lines, so the user can drop
// the flag entirely.
//
// Review-Followup F8: validation rejects signs and whitespace
// deterministically. Slice-v1-logs T0-(c) deliberately sets no upper
// bound; Compose receives very large decimal strings and decides
// whether it can handle them.
func validateLogsTailFlag(raw string) error {
	if raw == "" {
		return nil
	}
	if raw == "all" {
		return fmt.Errorf(
			"%w: `--tail \"all\"` is the implicit default; omit the flag to stream all lines",
			ErrInvalidLogsTail)
	}
	if !isDecimalDigits(raw) {
		return fmt.Errorf("%w: got %q", ErrInvalidLogsTail, raw)
	}
	return nil
}

func isDecimalDigits(raw string) bool {
	for _, ch := range raw {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

// parseLogsServices validates every positional service name (regex
// only, T0-(b)) and drops duplicates, keeping the first-seen order.
func parseLogsServices(args []string) ([]string, error) {
	var out []string
	seen := map[string]bool{}
	for _, a := range args {
		svc, err := domain.NewServiceName(a)
		if err != nil {
			return nil, err
		}
		if !seen[svc.String()] {
			seen[svc.String()] = true
			out = append(out, svc.String())
		}
	}
	return out, nil
}

// validateLogsTimeFlag checks a `--since` / `--until` value: empty =
// not set; otherwise a positive duration or an absolute timestamp in
// one of the accepted layouts (the value is forwarded verbatim to
// Compose, which accepts exactly these forms).
func validateLogsTimeFlag(name, raw string) error {
	if raw == "" {
		return nil
	}
	if d, err := time.ParseDuration(raw); err == nil {
		if d <= 0 {
			return fmt.Errorf("%w: %s %q must be a positive duration", ErrInvalidLogsTime, name, raw)
		}
		return nil
	}
	for _, layout := range []string{time.RFC3339, time.RFC3339Nano, "2006-01-02T15:04:05", "2006-01-02T15:04", "2006-01-02"} {
		if _, err := time.Parse(layout, raw); err == nil {
			return nil
		}
	}
	return fmt.Errorf("%w: got %s %q", ErrInvalidLogsTime, name, raw)
}
