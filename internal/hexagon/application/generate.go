package application

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// GenerateService implements [driving.GenerateUseCase] for the M7
// generators (LH-FA-GEN-001..005). T1 ships only the skeleton: the
// `u-boot.yaml`-existence check, the per-artefact dispatch, and four
// not-yet-implemented handler stubs that surface [errStubHandler].
// T2..T5 replace the stubs one by one; T5 removes the
// [errStubHandler]-pin test entirely.
//
// Driven-Sentinel-Scan (M7-T1 DoD): the pre-T1 scan of
// `internal/hexagon/port/driven/` confirmed that no
// `driven.ErrFileSystem*` sentinel exists today. T2..T5 therefore
// wrap unexpected `FileSystem.ReadFile`/`WriteFile`/`Stat` errors in
// [driving.ErrGenerateFileSystem] rather than relying on
// `errors.Is` against a driven sentinel. If a future slice
// introduces a driven filesystem sentinel, the wrap can collapse to
// a direct `errors.Is` without touching the CLI exit-code mapping.
type GenerateService struct {
	fs     driven.FileSystem
	yaml   driven.YAMLCodec
	logger driven.Logger
}

// Static check: GenerateService satisfies the driving port.
var _ driving.GenerateUseCase = (*GenerateService)(nil)

// NewGenerateService constructs the service with the driven adapters
// injected by the wiring layer. logger accepts nil and is routed to
// the package-local [noopLogger] (same nil-tolerance contract as
// [NewAddServiceService] and [NewUpService]). fs and yaml are
// mandatory.
func NewGenerateService(fs driven.FileSystem, yaml driven.YAMLCodec, logger driven.Logger) *GenerateService {
	if logger == nil {
		logger = noopLogger{}
	}
	return &GenerateService{fs: fs, yaml: yaml, logger: logger}
}

// errStubHandler is the package-internal pin that proves
// `u-boot generate <artifact>` is reachable at runtime but the
// per-artefact handler has not been wired yet. Unexported on purpose:
// the driving-port API does not surface it, so future callers cannot
// accidentally branch on it. The T1 test pins
// `errors.Is(err, errStubHandler)` for all four artefacts; each of
// T2..T5 reduces the pinned count by one, and T5 removes the test
// when the last stub is replaced (see slice-m7-generate.md DoD).
var errStubHandler = errors.New("generate: handler not implemented")

// Generate implements [driving.GenerateUseCase.Generate]. The
// dispatch order mirrors `AddServiceService.Add`:
//
//  1. validate BaseDir is non-empty (non-sentinel error).
//  2. check that `<BaseDir>/u-boot.yaml` exists; otherwise return
//     [driving.ErrProjectNotInitialized] (LH-FA-INIT-001 precondition,
//     reused from M5/M6).
//  3. dispatch on req.Artifact to the per-artefact handler. T1 stubs
//     all four; T2..T5 replace them with real implementations.
//
// ctx is threaded to the handlers so the T2..T5 implementations can
// honour cancellation without changing the call site here.
func (s *GenerateService) Generate(ctx context.Context, req driving.GenerateRequest) (driving.GenerateResponse, error) {
	if req.BaseDir == "" {
		return driving.GenerateResponse{}, errors.New("BaseDir is required")
	}

	if err := s.checkProjectInitialized(req.BaseDir); err != nil {
		return driving.GenerateResponse{}, err
	}

	s.logger.Debug("generate dispatch",
		"baseDir", req.BaseDir,
		"artifact", req.Artifact.String(),
	)

	switch req.Artifact {
	case domain.ArtifactChangelog:
		return s.generateChangelog(ctx, req)
	case domain.ArtifactReadme:
		return s.generateReadme(ctx, req)
	case domain.ArtifactEnvExample:
		return s.generateEnvExample(ctx, req)
	case domain.ArtifactDevcontainer:
		return s.generateDevcontainer(ctx, req)
	}
	// Unreachable in practice — the CLI validates via
	// [domain.NewArtifact] before constructing the request. Defensive
	// branch maps any future enum value (added without updating this
	// switch) to ErrInvalidArtifact rather than silently no-op.
	return driving.GenerateResponse{}, fmt.Errorf("%w: %v", domain.ErrInvalidArtifact, req.Artifact)
}

// checkProjectInitialized mirrors `DownService.checkProjectInitialized`
// and `UpService.checkProjectInitialized` so M7 produces the same
// sentinel-mapping behaviour at the CLI.
func (s *GenerateService) checkProjectInitialized(baseDir string) error {
	path := filepath.Join(baseDir, "u-boot.yaml")
	exists, err := s.fs.Exists(path)
	if err != nil {
		return fmt.Errorf("generate service: Exists(%q): %w", path, err)
	}
	if !exists {
		return fmt.Errorf("generate service: %q absent: %w", path, driving.ErrProjectNotInitialized)
	}
	return nil
}

// The four handler stubs use `_` as the receiver because T1 does not
// touch s.fs / s.yaml / s.logger yet. Each of T2..T5 renames `_` back
// to `s` when it wires the real handler, which the revive
// unused-receiver rule then accepts.

func (*GenerateService) generateChangelog(_ context.Context, _ driving.GenerateRequest) (driving.GenerateResponse, error) {
	return driving.GenerateResponse{}, fmt.Errorf("generate changelog: %w", errStubHandler)
}

func (*GenerateService) generateReadme(_ context.Context, _ driving.GenerateRequest) (driving.GenerateResponse, error) {
	return driving.GenerateResponse{}, fmt.Errorf("generate readme: %w", errStubHandler)
}

func (*GenerateService) generateEnvExample(_ context.Context, _ driving.GenerateRequest) (driving.GenerateResponse, error) {
	return driving.GenerateResponse{}, fmt.Errorf("generate env-example: %w", errStubHandler)
}

func (*GenerateService) generateDevcontainer(_ context.Context, _ driving.GenerateRequest) (driving.GenerateResponse, error) {
	return driving.GenerateResponse{}, fmt.Errorf("generate devcontainer: %w", errStubHandler)
}
