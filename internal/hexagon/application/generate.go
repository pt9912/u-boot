package application

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	iofs "io/fs"
	"path/filepath"

	"github.com/pt9912/u-boot/internal/hexagon/application/managedblock"
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

// readProjectConfig reads and parses `<baseDir>/u-boot.yaml`. Used by
// per-artefact handlers that need the project name (`{{.Name}}`) for
// template rendering. Maps TOCTOU file-vanish and parse errors to
// [driving.ErrProjectNotInitialized] to mirror the M5
// `detectServiceState`-classifier behaviour — a missing/malformed
// config is a fachliche precondition failure, not a technical FS
// fault.
//
// The handlers call this *after* `checkProjectInitialized` has
// already returned nil at dispatch entry; the duplicate work
// (Exists + ReadFile) is intentional and mirrors how `AddServiceService`
// re-reads u-boot.yaml inside `detectServiceState`. Folding both into
// a single `readProjectConfig` at the dispatcher would force every
// handler that doesn't need the parsed config (none today, but T3/T4
// will at least need the Name field too) to receive the cfg as a
// parameter — preferable to leave the read at the handler boundary so
// future handlers stay self-sufficient.
func (s *GenerateService) readProjectConfig(baseDir string) (ubootYAMLConfig, error) {
	yamlPath := filepath.Join(baseDir, "u-boot.yaml")
	body, err := s.fs.ReadFile(yamlPath)
	if err != nil {
		if errors.Is(err, iofs.ErrNotExist) {
			return ubootYAMLConfig{}, fmt.Errorf("%w: %s vanished between Exists and ReadFile",
				driving.ErrProjectNotInitialized, yamlPath)
		}
		return ubootYAMLConfig{}, fmt.Errorf("%w: read u-boot.yaml: %v",
			driving.ErrGenerateFileSystem, err)
	}
	var cfg ubootYAMLConfig
	if err := s.yaml.Unmarshal(body, &cfg); err != nil {
		return ubootYAMLConfig{}, fmt.Errorf("%w: parse u-boot.yaml: %v",
			driving.ErrProjectNotInitialized, err)
	}
	return cfg, nil
}

// The remaining handler stubs use `_` as the receiver because T1
// does not touch s.fs / s.yaml / s.logger yet. Each of T3..T5
// renames `_` back to `s` when it wires the real handler, which the
// revive unused-receiver rule then accepts.

func (*GenerateService) generateChangelog(_ context.Context, _ driving.GenerateRequest) (driving.GenerateResponse, error) {
	return driving.GenerateResponse{}, fmt.Errorf("generate changelog: %w", errStubHandler)
}

func (*GenerateService) generateReadme(_ context.Context, _ driving.GenerateRequest) (driving.GenerateResponse, error) {
	return driving.GenerateResponse{}, fmt.Errorf("generate readme: %w", errStubHandler)
}

// generateEnvExample implements the M7-T2 state machine for
// `.env.example` (LH-FA-GEN-004 / LH-FA-GEN-005). The slice plan's
// four-state table:
//
//	absent              → render full template, write file → Created
//	present-with-block  → splice re-rendered block → UpdatedBlock or NoOp
//	present-no-block    → ErrGenerateManualConflict (Code 10)
//	present-malformed   → ErrGenerateManualConflict (Code 10, different detail)
//
// Idempotency contract (LH-FA-GEN-005): a second invocation against
// an artefact that already matches the rendered block must return
// [driving.GenerateActionNoOp] with `Changed = nil` and zero
// WriteFile calls — the T2 NoOp-pin test asserts both. Content
// outside the BEGIN/END region (service add-on blocks, free-form
// user variables) survives any UpdatedBlock splice byte-identically;
// that is the [managedblock.Replace] contract, asserted by the
// T2 add-on-preservation test.
func (s *GenerateService) generateEnvExample(_ context.Context, req driving.GenerateRequest) (driving.GenerateResponse, error) {
	cfg, err := s.readProjectConfig(req.BaseDir)
	if err != nil {
		return driving.GenerateResponse{}, err
	}

	const relPath = ".env.example"
	targetPath := filepath.Join(req.BaseDir, relPath)
	marker := managedblock.Marker{Style: managedblock.StyleHash, Name: managedblock.InitName}

	rendered, err := renderTemplate("env.example.tmpl", templateData{Name: cfg.Project.Name})
	if err != nil {
		return driving.GenerateResponse{}, err
	}

	exists, err := s.fs.Exists(targetPath)
	if err != nil {
		return driving.GenerateResponse{}, fmt.Errorf("%w: Exists(%q): %v",
			driving.ErrGenerateFileSystem, targetPath, err)
	}

	// State: absent — write the whole rendered template.
	if !exists {
		if err := s.fs.WriteFile(targetPath, rendered, defaultFileMode); err != nil {
			return driving.GenerateResponse{}, fmt.Errorf("%w: write %q: %v",
				driving.ErrGenerateFileSystem, targetPath, err)
		}
		s.logger.Info("generate env-example: created",
			"path", relPath, "project", cfg.Project.Name)
		return driving.GenerateResponse{
			Artifact: req.Artifact,
			Action:   driving.GenerateActionCreated,
			Changed:  []string{relPath},
		}, nil
	}

	existing, err := s.fs.ReadFile(targetPath)
	if err != nil {
		// TOCTOU: file vanished between Exists and ReadFile. Surface
		// as a filesystem error rather than re-classify as absent —
		// the user-visible message is more useful, and an immediate
		// retry will hit the absent branch cleanly.
		return driving.GenerateResponse{}, fmt.Errorf("%w: read %q: %v",
			driving.ErrGenerateFileSystem, targetPath, err)
	}

	// Extract the BEGIN..END region from the freshly-rendered template.
	// Template-side ErrBlockNotFound/Malformed would mean the embedded
	// template has rotted — a programmer error, not a user-side issue,
	// so it surfaces as a plain error without an M7 sentinel.
	renderedBlock, err := renderManagedBlockOnly(rendered, marker)
	if err != nil {
		return driving.GenerateResponse{}, fmt.Errorf(
			"extract init block from rendered env.example.tmpl: %w", err)
	}

	// Classify the existing file's block.
	start, end, findErr := managedblock.Find(existing, marker)
	switch {
	case errors.Is(findErr, managedblock.ErrBlockNotFound):
		return driving.GenerateResponse{}, fmt.Errorf(
			"%w: %q exists without an `init` managed block; rename the file or insert a `# BEGIN U-BOOT MANAGED BLOCK: init` marker manually",
			driving.ErrGenerateManualConflict, relPath)
	case errors.Is(findErr, managedblock.ErrBlockMalformed):
		return driving.GenerateResponse{}, fmt.Errorf(
			"%w: %q has a malformed `init` managed block (%v); rename the file or repair the BEGIN/END markers manually",
			driving.ErrGenerateManualConflict, relPath, findErr)
	case findErr != nil:
		return driving.GenerateResponse{}, fmt.Errorf("find init block in %q: %w", relPath, findErr)
	}

	// State: present-with-block. NoOp if the existing block bytes are
	// already equal to the rendered block bytes (idempotency).
	if bytes.Equal(existing[start:end], renderedBlock) {
		s.logger.Debug("generate env-example: no-op", "path", relPath)
		return driving.GenerateResponse{
			Artifact: req.Artifact,
			Action:   driving.GenerateActionNoOp,
			Changed:  nil,
		}, nil
	}

	// State: present-with-block, block stale — splice the new bytes.
	updated, err := managedblock.Replace(existing, marker, renderedBlock)
	if err != nil {
		return driving.GenerateResponse{}, fmt.Errorf("replace init block in %q: %w", relPath, err)
	}
	if err := s.fs.WriteFile(targetPath, updated, defaultFileMode); err != nil {
		return driving.GenerateResponse{}, fmt.Errorf("%w: write %q: %v",
			driving.ErrGenerateFileSystem, targetPath, err)
	}
	s.logger.Info("generate env-example: updated block",
		"path", relPath, "project", cfg.Project.Name)
	return driving.GenerateResponse{
		Artifact: req.Artifact,
		Action:   driving.GenerateActionUpdatedBlock,
		Changed:  []string{relPath},
	}, nil
}

func (*GenerateService) generateDevcontainer(_ context.Context, _ driving.GenerateRequest) (driving.GenerateResponse, error) {
	return driving.GenerateResponse{}, fmt.Errorf("generate devcontainer: %w", errStubHandler)
}
