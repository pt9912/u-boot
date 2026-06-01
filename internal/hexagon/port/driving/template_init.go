package driving

import (
	"context"
	"errors"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// TemplateInitRequest is the input for [TemplateInitUseCase.Init].
//
// The full `u-boot init` flow (LH-FA-INIT-001..007) is implemented
// by [InitProjectUseCase]; this use case is the per-template render
// stage that the init pipeline delegates to when `--template <name>`
// is set. The caller (the CLI adapter or the InitProjectService
// integration in T4) is responsible for project-state checks
// (BaseDir exists, project not already initialised), git init, and
// post-render bookkeeping.
type TemplateInitRequest struct {
	// BaseDir is the absolute target directory. Mandatory; an empty
	// string returns a plain `errors.New` (no sentinel — wiring
	// invariant, not a user error).
	BaseDir string

	// ProjectName is the validated project name; surfaces as
	// `.Name` in every `*.tmpl` rendering.
	ProjectName domain.ProjectName

	// TemplateName is the catalog identifier (e.g. `basic`). The
	// driven [driven.TemplateFiles] adapter resolves it; an unknown
	// name surfaces as [ErrTemplateNotFound].
	TemplateName string
}

// TemplateInitResponse is the output of [TemplateInitUseCase.Init].
// Created is the list of project-relative paths that were written,
// alphabetically sorted — the CLI adapter prints the summary in
// that order.
type TemplateInitResponse struct {
	Created []string
}

// ErrTemplateNotFound signals that the requested `--template <name>`
// does not exist in the catalog. The use case wraps the driven-port
// [driven.ErrTemplateNotFound] with this sentinel so the CLI maps
// to LH-FA-CLI-006 exit code 10 (user-actionable: pick a real
// template via `u-boot template list`).
var ErrTemplateNotFound = errors.New("template: not found")

// ErrInvalidTemplatePath signals that a file inside the template's
// tree has a malformed or escaping path
// (see [domain.NewTemplatePath]'s reject list). Exit code 10 —
// the user must edit the template, not retry.
var ErrInvalidTemplatePath = errors.New("template: invalid path")

// ErrTemplateRender signals a render-time failure: malformed
// `text/template` syntax, evaluation panic, or an IO failure during
// the per-file copy/write. Exit code 14 (technical) — the user
// generally cannot fix it without changing the template.
var ErrTemplateRender = errors.New("template: render failed")

// TemplateInitUseCase implements the render stage of
// `u-boot init --template <name>` (LH-FA-TPL-001 /
// slice-v1-template-init).
//
// Contract:
//
//   - Init walks the template's file tree and writes one file per
//     `.tmpl` (rendered via Go `text/template`) or non-`.tmpl`
//     (copied 1:1) entry, into BaseDir + relative path.
//     `template.yaml` is skipped (metadata, already validated by
//     the listing slice).
//   - Each per-file relative path is validated through
//     [domain.NewTemplatePath]; any escape (`..`, absolute, drive
//     letter) fails the call with [ErrInvalidTemplatePath].
//   - The response Created list is alphabetically sorted.
type TemplateInitUseCase interface {
	Init(ctx context.Context, req TemplateInitRequest) (TemplateInitResponse, error)
}
