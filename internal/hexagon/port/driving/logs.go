package driving

import (
	"context"
	"io"
)

// LogsRequest is the input for [LogsUseCase.Logs] (LH-FA-UP-005).
// It is the application-layer expression of `u-boot logs [service]
// [--follow] [--tail <n>]`; the CLI adapter translates the
// positional argument and flags into this struct.
//
// Slice-v1-logs §T0-Outcomes pinned the contract:
//
//   - Service-Filter (T0-(a)): empty Service = Compose-Default
//     (all services in compose.yaml); no `activeServiceNames(cfg)`
//     filter. A single service name is forwarded as the only
//     element of `ComposeLogsOptions.Services`.
//   - Service-Name-Validation (T0-(b)): the CLI layer validates
//     the positional argument via [domain.NewServiceName] (regex
//     only). Existence at runtime is delegated to Compose;
//     unknown services map to `driven.ErrComposeRuntime` →
//     Exit-Code 12.
//   - Tail-Default (T0-(c)): empty Tail = Compose-Default. The
//     [LogsService] normalises `Tail == ""` to `"all"` before
//     calling the adapter; numeric inputs ≥ 0 are accepted
//     verbatim. CLI Stage-1-validation rejects negative or
//     non-numeric strings with Exit-Code 2 (before the use case
//     runs).
//   - Format-Flags (T0-(d)): only `--follow` and `--tail` are
//     exposed. No `--no-log-prefix`, no `--timestamps` — those
//     live in a future follow-up slice with a documented trigger.
type LogsRequest struct {
	// BaseDir is the absolute path of the initialized u-boot
	// project. Mandatory; the CLI adapter defaults it to the
	// current working directory (mirroring `u-boot up` / `down`).
	BaseDir string

	// Service is the optional positional argument from
	// `u-boot logs [service]`. Empty (default) means "all
	// Compose services" (T0-(a)). When non-empty, the value is
	// trusted as already-format-validated by the CLI adapter
	// (`domain.NewServiceName`).
	Service string

	// Follow mirrors `docker compose logs --follow`. When true,
	// the call blocks until SIGINT (Ctrl-C) cancels the context.
	// The application service short-circuits the resulting
	// `context.Canceled` to a nil error so the CLI exits 0.
	Follow bool

	// Tail mirrors `docker compose logs --tail <value>`. Accepts
	// either an empty string (Compose-Default = all) or a
	// non-negative decimal integer string (`"0"`, `"100"`, …).
	// The application service normalises empty → `"all"` for the
	// downstream adapter call.
	Tail string

	// OutputSink is the writer the application passes to the
	// `DockerEngine.ComposeLogs` adapter so the compose log
	// stream reaches the user. Analog to
	// [UpRequest.ProgressSink] / [DownRequest.ProgressSink] from
	// M6 — the CLI wires this to `cmd.OutOrStdout()`. `nil` is
	// treated as `io.Discard` by the application service.
	OutputSink io.Writer
}

// LogsResponse is the output of [LogsUseCase.Logs]. Empty by
// design — the actual log lines stream live through
// [LogsRequest.OutputSink]; there is no structured per-line
// return value the CLI needs to render.
type LogsResponse struct{}

// LogsUseCase is the driving-port for `u-boot logs` (LH-FA-UP-005).
// The CLI adapter holds a reference and calls [Logs] from the
// Cobra command handler.
//
// Sentinel contract (slice-v1-logs §Aufhebungsbedingung):
//
//   - SIGINT cancellation → `(LogsResponse{}, nil)`: the
//     application service intercepts `context.Canceled` /
//     `context.DeadlineExceeded` from the adapter and returns
//     success so the CLI exits 0 (tail-konform).
//   - Project not initialised (no `u-boot.yaml`) →
//     [ErrProjectNotInitialized] (Exit 10).
//   - Compose file missing → [ErrComposeFileMissing] (Exit 10).
//   - Docker environment unavailable → wraps
//     `driven.ErrDockerUnavailable` (Exit 11).
//   - Compose runtime failure (unknown service, exit ≠ 0) →
//     wraps `driven.ErrComposeRuntime` (Exit 12).
//   - Invalid service-name format → wraps
//     `domain.ErrInvalidServiceName` (Exit 10).
type LogsUseCase interface {
	Logs(ctx context.Context, req LogsRequest) (LogsResponse, error)
}
