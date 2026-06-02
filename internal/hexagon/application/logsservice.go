package application

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// LogsService implements [driving.LogsUseCase] for LH-FA-UP-005
// (`u-boot logs`). Orchestrates the Project-State-Check (u-boot.yaml
// + compose.yaml present, analog to [UpService]) and forwards the
// streaming compose-logs call to the driven adapter.
//
// Slice-v1-logs §T0-Outcomes pinned the surface:
//
//   - T0-(a): Service-Filter is Compose-Default — `Service == ""`
//     means "all services in compose.yaml", no
//     `activeServiceNames(cfg)`-Filter.
//   - T0-(b): Service-Name-Validation is regex-only via the CLI
//     adapter (`domain.NewServiceName`). The application layer
//     trusts the value (no double-validation here). Unknown
//     services at runtime → `driven.ErrComposeRuntime` (Exit 12).
//     **Non-CLI callers** (future RPC, direct tests, other
//     adapters) MUST validate `req.Service` via
//     `domain.NewServiceName` themselves before calling
//     `LogsService.Logs`, OR accept that an invalid name is
//     forwarded verbatim to Compose and surfaces as Exit-12
//     instead of Exit-10. Review-Followup F7 made the validation-
//     distribution explicit.
//   - T0-(c): Empty Tail normalises to Compose's `"all"` constant
//     before reaching the adapter. Numeric values pass through
//     verbatim.
//   - SIGINT-Pass-Through: `context.Canceled` /
//     `context.DeadlineExceeded` from the adapter surface as
//     `(LogsResponse{}, nil)` so the CLI exits 0.
type LogsService struct {
	fs     driven.FileSystem
	engine driven.DockerEngine
	logger driven.Logger
}

// NewLogsService wires the driven ports the use case needs. The
// logger is optional (slice-v1-logs follows the M6 pattern of
// best-effort structured logging); fs and engine are mandatory.
func NewLogsService(fs driven.FileSystem, engine driven.DockerEngine, logger driven.Logger) *LogsService {
	if logger == nil {
		logger = noopLogger{}
	}
	return &LogsService{fs: fs, engine: engine, logger: logger}
}

// Logs implements [driving.LogsUseCase.Logs]. The flow:
//
//  1. BaseDir present + u-boot.yaml exists (ErrProjectNotInitialized
//     on miss) + compose.yaml exists (ErrComposeFileMissing on miss).
//     Shared semantics with [UpService] / [DownService].
//  2. Normalise Tail: empty → `"all"` (T0-(c)).
//  3. Build [driven.ComposeLogsOptions] and call
//     [driven.DockerEngine.ComposeLogs]. Single-service requests
//     become a one-element Services slice; empty Service leaves the
//     slice nil (Compose-Default).
//  4. SIGINT-Pass-Through: if the adapter returns a wrapped
//     `context.Canceled` / `context.DeadlineExceeded`, the use case
//     short-circuits to `(LogsResponse{}, nil)` so the CLI exit-code
//     is 0 (LH-FA-UP-005 SIGINT-Vertrag).
func (s *LogsService) Logs(ctx context.Context, req driving.LogsRequest) (driving.LogsResponse, error) {
	if req.BaseDir == "" {
		return driving.LogsResponse{}, errors.New("logs service: BaseDir is empty")
	}
	if err := s.checkProjectInitialized(req.BaseDir); err != nil {
		return driving.LogsResponse{}, err
	}
	if err := s.checkComposeFile(req.BaseDir); err != nil {
		return driving.LogsResponse{}, err
	}

	opts := driven.ComposeLogsOptions{
		Follow: req.Follow,
		Tail:   normaliseTail(req.Tail),
		Sink:   req.OutputSink,
	}
	if req.Service != "" {
		opts.Services = []string{req.Service}
	}

	s.logger.Debug("logs: invoking compose logs",
		"baseDir", req.BaseDir,
		"service", req.Service,
		"follow", req.Follow,
		"tail", opts.Tail)
	err := s.engine.ComposeLogs(ctx, req.BaseDir, opts)
	if err != nil {
		// SIGINT-Pass-Through (slice-v1-logs §SIGINT-Vertrag
		// Schicht 2/3): a cancelled context is the user's
		// expected exit path for `--follow`, not an error.
		// `context.DeadlineExceeded` is included for symmetry
		// with future automated callers (test timeouts, RPC).
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			s.logger.Debug("logs: context cancelled, treating as success",
				"baseDir", req.BaseDir)
			return driving.LogsResponse{}, nil
		}
		return driving.LogsResponse{}, fmt.Errorf("logs service: ComposeLogs on %q: %w",
			req.BaseDir, err)
	}
	return driving.LogsResponse{}, nil
}

// checkProjectInitialized verifies that `<BaseDir>/u-boot.yaml`
// exists (analog [UpService.checkProjectInitialized]). Permission/IO
// errors are wrapped without the project-not-initialized sentinel —
// they are technical, not fachlich.
func (s *LogsService) checkProjectInitialized(baseDir string) error {
	path := filepath.Join(baseDir, "u-boot.yaml")
	exists, err := s.fs.Exists(path)
	if err != nil {
		return fmt.Errorf("logs service: Exists(%q): %w", path, err)
	}
	if !exists {
		return fmt.Errorf("logs service: %q absent: %w", path, driving.ErrProjectNotInitialized)
	}
	return nil
}

// checkComposeFile verifies that `<BaseDir>/compose.yaml` exists
// (analog [UpService.readComposeFile] — but `logs` does not parse
// the file contents; Compose itself drives the streaming, so a bare
// existence check is enough).
func (s *LogsService) checkComposeFile(baseDir string) error {
	path := filepath.Join(baseDir, "compose.yaml")
	exists, err := s.fs.Exists(path)
	if err != nil {
		return fmt.Errorf("logs service: Exists(%q): %w", path, err)
	}
	if !exists {
		return fmt.Errorf("logs service: %q absent: %w", path, driving.ErrComposeFileMissing)
	}
	return nil
}

// normaliseTail enforces the T0-(c) contract at the application
// boundary: empty incoming Tail → Compose's `"all"` constant
// (verbatim string the adapter forwards as `--tail all`). Numeric
// strings pass through verbatim; the CLI layer is responsible for
// rejecting negative / non-numeric inputs before reaching the
// service, and for parsing the format-validated service name via
// `domain.NewServiceName` (T0-(b)).
//
// Review-Followup F5: free function rather than `*LogsService`-
// method on purpose — pure function over its input, no service-
// state needed. Keeping it free makes the T0-(c) contract
// trivially unit-testable and signals "this is the only place
// that produces the `\"all\"` constant" (Plan-Followup N2
// invariant).
func normaliseTail(in string) string {
	if in == "" {
		return "all"
	}
	return in
}
