package application

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// DownService implements [driving.DownUseCase]. It orchestrates
// `u-boot down`: validate the project layout (mirrors [UpService]),
// route the destructive `--volumes` confirmation per the M6 slice
// §T5 truth table, then hand off to the [driven.DockerEngine] for
// the actual `compose down`.
//
// Truth table (LH-FA-CLI-005A §254 / slice plan §T5):
//
//	RemoveVolumes | AssumeYes | NonInteractive | behaviour
//	false         | *         | *              | proceed, no confirmer call
//	true          | true      | *              | proceed, no confirmer call
//	true          | false     | true           | fail-fast with ErrConfirmationRequired,
//	              |           |                | no confirmer call, no engine call
//	true          | false     | false          | call confirmer; (true, nil) → proceed,
//	              |           |                | (false, nil) → ErrConfirmationRequired,
//	              |           |                | error → wrap and return
type DownService struct {
	fs        driven.FileSystem
	engine    driven.DockerEngine
	confirmer driven.Confirmer
	logger    driven.Logger
}

// Static check: DownService satisfies the driving port.
var _ driving.DownUseCase = (*DownService)(nil)

// NewDownService constructs the service with the driven adapters the
// M6 down-flow needs. A nil logger is routed to the package-level
// noopLogger so tests and dry-runs do not need a stub.
func NewDownService(fs driven.FileSystem, engine driven.DockerEngine, confirmer driven.Confirmer, logger driven.Logger) *DownService {
	if logger == nil {
		logger = noopLogger{}
	}
	return &DownService{fs: fs, engine: engine, confirmer: confirmer, logger: logger}
}

// Down implements [driving.DownUseCase.Down].
func (s *DownService) Down(ctx context.Context, req driving.DownRequest) (driving.DownResponse, error) {
	if req.BaseDir == "" {
		return driving.DownResponse{}, errors.New("down service: BaseDir is empty")
	}
	if err := s.checkProjectInitialized(req.BaseDir); err != nil {
		return driving.DownResponse{}, err
	}
	if err := s.checkComposeFilePresent(req.BaseDir); err != nil {
		return driving.DownResponse{}, err
	}

	if err := s.runConfirmationGate(ctx, req); err != nil {
		return driving.DownResponse{}, err
	}

	if err := s.engine.ComposeDown(ctx, req.BaseDir, driven.ComposeDownOptions{
		RemoveVolumes: req.RemoveVolumes,
		ProgressSink:  req.ProgressSink,
	}); err != nil {
		return driving.DownResponse{}, fmt.Errorf("down service: ComposeDown on %q: %w", req.BaseDir, err)
	}
	return driving.DownResponse{RemovedVolumes: req.RemoveVolumes}, nil
}

// checkProjectInitialized verifies that `<BaseDir>/u-boot.yaml`
// exists. Mirrors [UpService.checkProjectInitialized] so the two
// use cases produce identical sentinel-mapping behaviour at the CLI.
func (s *DownService) checkProjectInitialized(baseDir string) error {
	path := filepath.Join(baseDir, "u-boot.yaml")
	exists, err := s.fs.Exists(path)
	if err != nil {
		return fmt.Errorf("down service: Exists(%q): %w", path, err)
	}
	if !exists {
		return fmt.Errorf("down service: %q absent: %w", path, driving.ErrProjectNotInitialized)
	}
	return nil
}

// checkComposeFilePresent verifies that `<BaseDir>/compose.yaml`
// exists. `down` does not need to parse the file (no per-service
// classification like Up); it only confirms presence so the engine
// invocation has a target.
func (s *DownService) checkComposeFilePresent(baseDir string) error {
	path := filepath.Join(baseDir, "compose.yaml")
	exists, err := s.fs.Exists(path)
	if err != nil {
		return fmt.Errorf("down service: Exists(%q): %w", path, err)
	}
	if !exists {
		return fmt.Errorf("down service: %q absent: %w", path, driving.ErrComposeFileMissing)
	}
	return nil
}

// runConfirmationGate implements the §T5 truth table. Returns nil
// when the request should proceed to ComposeDown; returns wrapped
// [driving.ErrConfirmationRequired] (CLI code 10) when the
// destructive op is refused or skipped per LH-FA-CLI-005A §254.
func (s *DownService) runConfirmationGate(ctx context.Context, req driving.DownRequest) error {
	if !req.RemoveVolumes {
		return nil // Row 1: non-destructive; no confirmation needed.
	}
	if req.AssumeYes {
		return nil // Row 2: explicit --yes auto-approves.
	}
	if req.NonInteractive {
		// Row 3: LH-FA-CLI-005A §254 — non-interactive without
		// --yes ⇒ fail-fast with ErrConfirmationRequired (code 10).
		// No confirmer call, no engine call.
		return fmt.Errorf("down service: --volumes refused in --no-interactive without --yes: %w", driving.ErrConfirmationRequired)
	}
	// Row 4: interactive confirmation.
	confirmed, err := s.confirmer.ConfirmRemoveVolumes(ctx, req.BaseDir)
	if err != nil {
		return fmt.Errorf("down service: confirmer error: %w", err)
	}
	if !confirmed {
		return fmt.Errorf("down service: --volumes declined by user: %w", driving.ErrConfirmationRequired)
	}
	return nil
}
