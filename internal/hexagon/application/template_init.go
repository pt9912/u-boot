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
//nolint:gocyclo,cyclop // The walk-and-classify body is naturally branchy; splitting
// the per-entry handler into a helper hides the order constraint
// (skip-then-validate-then-render-then-write) without reducing
// effective cognitive load.
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
	var created []string

	walkErr := iofs.WalkDir(sub, ".", func(p string, d iofs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if p == "." || d.IsDir() {
			return nil
		}
		if p == templateMetadataFilename {
			return nil
		}
		written, err := s.renderOne(p, sub, req.BaseDir, data)
		if err != nil {
			return err
		}
		created = append(created, written)
		return nil
	})
	if walkErr != nil {
		return driving.TemplateInitResponse{}, walkErr
	}

	sort.Strings(created)
	s.logger.Debug("template init: rendered", "name", req.TemplateName, "files", len(created))
	return driving.TemplateInitResponse{Created: created}, nil
}

// renderOne handles a single non-directory, non-metadata entry from
// the template's file tree: validate path, read source, render or
// copy, ensure parent dir, write. Returns the project-relative
// output path that the caller appends to the Created list.
func (s *TemplateInitService) renderOne(p string, sub iofs.FS, baseDir string, data templateData) (string, error) {
	outRel := strings.TrimSuffix(p, templateSuffix)
	tp, err := domain.NewTemplatePath(outRel)
	if err != nil {
		return "", fmt.Errorf("%w: %w", driving.ErrInvalidTemplatePath, err)
	}

	src, err := iofs.ReadFile(sub, p)
	if err != nil {
		return "", fmt.Errorf("%w: read %s: %w", driving.ErrTemplateRender, p, err)
	}

	var content []byte
	if strings.HasSuffix(p, templateSuffix) {
		content, err = renderTemplateBytes(p, src, data)
		if err != nil {
			return "", fmt.Errorf("%w: %s: %w", driving.ErrTemplateRender, p, err)
		}
	} else {
		content = src
	}

	outAbs := filepath.Join(baseDir, tp.String())
	if err := s.fs.MkdirAll(filepath.Dir(outAbs), renderedDirMode); err != nil {
		return "", fmt.Errorf("%w: mkdir %s: %w", driving.ErrTemplateRender, filepath.Dir(outAbs), err)
	}
	if err := s.fs.WriteFile(outAbs, content, renderedFileMode); err != nil {
		return "", fmt.Errorf("%w: write %s: %w", driving.ErrTemplateRender, outAbs, err)
	}
	return tp.String(), nil
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
