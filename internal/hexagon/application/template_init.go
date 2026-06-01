package application

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	iofs "io/fs"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// renderedFileMode is the mode for files written by the template
// render loop. Matches [defaultFileMode] used by
// [InitProjectService]; sharing it would require moving the
// constant to a shared file, which is more churn than the value
// itself.
const renderedFileMode iofs.FileMode = 0o644

// renderedDirMode is the mode for parent directories created on
// the fly by the render loop (e.g. `docker/` for a template that
// writes `docker/Dockerfile`).
const renderedDirMode iofs.FileMode = 0o755

// templateMetadataFilename is the per-template metadata file the
// render loop skips. Already validated by the listing slice's
// catalog adapter; copying it into the user's project would leak
// internal schema into rendered output.
const templateMetadataFilename = "template.yaml"

// templateSuffix marks files that get rendered through Go
// `text/template`. Non-suffixed files are copied byte-identically
// (ADR-0009 §Entscheidung).
const templateSuffix = ".tmpl"

// TemplateInitService implements [driving.TemplateInitUseCase] —
// the render stage of `u-boot init --template <name>`
// (slice-v1-template-init). Stateless and concurrent-safe; every
// call walks the template's file tree afresh.
type TemplateInitService struct {
	files  driven.TemplateFiles
	fs     driven.FileSystem
	logger driven.Logger
}

// Static check.
var _ driving.TemplateInitUseCase = (*TemplateInitService)(nil)

// NewTemplateInitService constructs the service. files and fs are
// mandatory; logger accepts nil and is routed to [noopLogger]
// (matching the codebase convention).
func NewTemplateInitService(files driven.TemplateFiles, fs driven.FileSystem, logger driven.Logger) *TemplateInitService {
	if logger == nil {
		logger = noopLogger{}
	}
	return &TemplateInitService{files: files, fs: fs, logger: logger}
}

// Init opens the template's file tree, renders every `.tmpl`
// against the request's projection, copies every non-`.tmpl`
// 1:1, and writes both into BaseDir.
//
// Errors are classified for the CLI exit-code mapping:
//
//   - Unknown template name → [driving.ErrTemplateNotFound] (10)
//   - Path-escape attempt → [driving.ErrInvalidTemplatePath] (10)
//   - Render / IO failure → [driving.ErrTemplateRender] (14)
//
// Two-phase render (review-followup F1): the walk-loop first
// renders every file into memory; only after every render
// succeeds does the second loop write to disk. Bad templates
// (parse errors, path-escape) thus no longer leave half-populated
// project directories — the typical author-error class
// short-circuits before any side effect. Genuine IO failures
// during the write phase (disk full, permissions) still produce
// partial state, but at that point the user has to clean up
// anyway. WalkDir-layer errors (review-followup F5) get the
// ErrTemplateRender wrap so the CLI maps them to exit code 14
// instead of falling through to 1.
func (s *TemplateInitService) Init(ctx context.Context, req driving.TemplateInitRequest) (driving.TemplateInitResponse, error) {
	if req.BaseDir == "" {
		return driving.TemplateInitResponse{}, errors.New("BaseDir is required")
	}
	sub, err := s.files.Open(ctx, req.TemplateName)
	if err != nil {
		if errors.Is(err, driven.ErrTemplateNotFound) {
			return driving.TemplateInitResponse{}, fmt.Errorf("%w: %w", driving.ErrTemplateNotFound, err)
		}
		return driving.TemplateInitResponse{}, fmt.Errorf("%w: %w", driving.ErrTemplateRender, err)
	}

	data := templateData{Name: req.ProjectName.String()}

	// Phase 1: walk + render to memory. No disk writes here.
	planned, err := s.planRender(sub, data)
	if err != nil {
		return driving.TemplateInitResponse{}, err
	}

	// Phase 2: apply. Every render succeeded; now persist.
	created := make([]string, 0, len(planned))
	for _, f := range planned {
		outAbs := filepath.Join(req.BaseDir, f.path)
		if err := s.fs.MkdirAll(filepath.Dir(outAbs), renderedDirMode); err != nil {
			return driving.TemplateInitResponse{}, fmt.Errorf("%w: mkdir %s: %w", driving.ErrTemplateRender, filepath.Dir(outAbs), err)
		}
		if err := s.fs.WriteFile(outAbs, f.content, renderedFileMode); err != nil {
			return driving.TemplateInitResponse{}, fmt.Errorf("%w: write %s: %w", driving.ErrTemplateRender, outAbs, err)
		}
		created = append(created, f.path)
	}

	sort.Strings(created)
	s.logger.Debug("template init: rendered", "name", req.TemplateName, "files", len(created))
	return driving.TemplateInitResponse{Created: created}, nil
}

// renderedFile is one entry in the in-memory render plan: the
// project-relative output path (canonicalised by [domain.NewTemplatePath])
// and the bytes to write at it. The Init pipeline produces a slice
// of these during phase 1 and consumes it during phase 2.
type renderedFile struct {
	path    string
	content []byte
}

// planRender walks the template tree and produces an in-memory
// render plan without touching the destination filesystem. Every
// per-file failure short-circuits the walk; on success returns
// the full list ready for phase-2 writes.
func (s *TemplateInitService) planRender(sub iofs.FS, data templateData) ([]renderedFile, error) {
	var planned []renderedFile
	walkErr := iofs.WalkDir(sub, ".", func(p string, d iofs.DirEntry, walkErr error) error {
		if walkErr != nil {
			// review-followup F5: classify walk-layer errors so the
			// CLI exit-code mapping fires correctly.
			return fmt.Errorf("%w: walk %s: %w", driving.ErrTemplateRender, p, walkErr)
		}
		if p == "." || d.IsDir() {
			return nil
		}
		if p == templateMetadataFilename {
			return nil
		}
		f, err := s.renderOne(p, sub, data)
		if err != nil {
			return err
		}
		planned = append(planned, f)
		return nil
	})
	if walkErr != nil {
		return nil, walkErr
	}
	return planned, nil
}

// renderOne handles a single non-directory, non-metadata entry from
// the template's file tree: validate path, read source, render or
// copy. Returns the rendered file (path + bytes) for the phase-2
// write loop; no disk side effect.
func (*TemplateInitService) renderOne(p string, sub iofs.FS, data templateData) (renderedFile, error) {
	outRel := strings.TrimSuffix(p, templateSuffix)
	tp, err := domain.NewTemplatePath(outRel)
	if err != nil {
		return renderedFile{}, fmt.Errorf("%w: %w", driving.ErrInvalidTemplatePath, err)
	}

	src, err := iofs.ReadFile(sub, p)
	if err != nil {
		return renderedFile{}, fmt.Errorf("%w: read %s: %w", driving.ErrTemplateRender, p, err)
	}

	var content []byte
	if strings.HasSuffix(p, templateSuffix) {
		content, err = renderTemplateBytes(p, src, data)
		if err != nil {
			return renderedFile{}, fmt.Errorf("%w: %s: %w", driving.ErrTemplateRender, p, err)
		}
	} else {
		content = src
	}

	return renderedFile{path: tp.String(), content: content}, nil
}

// renderTemplateBytes runs a single `text/template` execution
// against arbitrary source bytes. The parallel to the M3-T2
// [renderTemplate] helper in templates.go is intentional — both
// use the same engine and same [templateData] projection so the
// rendered shape stays consistent between the internal init flow
// and the external template-init flow.
func renderTemplateBytes(name string, src []byte, data templateData) ([]byte, error) {
	tmpl, err := template.New(name).Parse(string(src))
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("execute: %w", err)
	}
	return buf.Bytes(), nil
}
