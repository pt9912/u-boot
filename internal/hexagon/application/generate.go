package application

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	iofs "io/fs"
	"path/filepath"
	"regexp"

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

// The remaining handler stub uses `_` as the receiver because T1
// does not touch s.fs / s.yaml / s.logger yet. T5 renames `_` back
// to `s` when it wires the real devcontainer handler.

// changelogHeaderRE matches a Keep-a-Changelog second-level header
// line (`## [<name>]`), capturing the inner identifier so callers
// can branch on whether it is `Unreleased` or a release version
// (`[0.1.0]`, etc.). Pinned at line-start so the regex itself
// rejects inline mentions; fenced-code-block headers are filtered
// by [isOffsetInsideFencedBlock] in the call sites — a quoted
// `## [1.2.3]` example before any real release section would
// otherwise trigger the RepairedManual splice path *inside* the
// fence and corrupt the user's Markdown (review-followup S1).
var changelogHeaderRE = regexp.MustCompile(`(?m)^##\s*\[([^\]]+)\]`)

// fenceMarkerRE matches a Markdown fenced-code-block opener/closer
// (3+ backticks at line start). Tilde-fenced blocks and indented
// code blocks are not recognised — neither is common in Keep-a-
// Changelog projects, and the heuristic stays cheap.
var fenceMarkerRE = regexp.MustCompile("(?m)^`{3,}")

// normaliseLF returns content with CRLF (\r\n) sequences replaced
// by LF (\n). Used by the [bytes.Equal] freshness comparison in
// [generateManagedFile] and [generateChangelog] so a file edited
// on Windows (CRLF) does not falsely register as user-edited
// against a template that ships with LF (review-followup S2).
//
// The returned slice is only used for the equality comparison —
// the actual splice into the existing file uses the un-normalised
// `existing` bytes, so the user's original line endings outside
// the spliced block stay intact.
func normaliseLF(content []byte) []byte {
	if !bytes.Contains(content, []byte{'\r', '\n'}) {
		return content
	}
	return bytes.ReplaceAll(content, []byte{'\r', '\n'}, []byte{'\n'})
}

// classifyExistingBlock wraps [managedblock.Find] for the three M7
// handlers (generateManagedFile, generateChangelog, planDevcontainerFile)
// so the BlockNotFound / BlockMalformed branches all surface the
// same [driving.ErrGenerateManualConflict] with a deterministic
// repair-hint message (review-followup N4). Other find errors get
// a generic wrap; the byte range is unset on any error.
func classifyExistingBlock(content []byte, marker managedblock.Marker, relPath string) (start, end int, err error) {
	start, end, findErr := managedblock.Find(content, marker)
	switch {
	case errors.Is(findErr, managedblock.ErrBlockNotFound):
		return 0, 0, fmt.Errorf(
			"%w: %q exists without an `init` managed block; rename the file or insert the format-appropriate BEGIN/END markers from LH-SA-FILE-002 manually",
			driving.ErrGenerateManualConflict, relPath)
	case errors.Is(findErr, managedblock.ErrBlockMalformed):
		return 0, 0, fmt.Errorf(
			"%w: %q has a malformed `init` managed block (%v); rename the file or repair the BEGIN/END markers manually",
			driving.ErrGenerateManualConflict, relPath, findErr)
	case findErr != nil:
		return 0, 0, fmt.Errorf("find init block in %q: %w", relPath, findErr)
	}
	return start, end, nil
}

// changelogUnreleasedStub is the bare scaffold inserted by the
// RepairedManual path when the user removed the `## [Unreleased]`
// section while cutting a release. Empty Added/Changed/Fixed
// subsections so the user can pick up immediately without manual
// formatting.
const changelogUnreleasedStub = "## [Unreleased]\n\n### Added\n\n### Changed\n\n### Fixed\n\n"

// generateChangelog implements the M7-T4 state machine for
// `CHANGELOG.md` (LH-FA-GEN-002 / LH-AK-007 / LH-FA-GEN-005). The
// template's `init` managed block carries the initial scaffold
// **and** the `## [Unreleased]` section, which the user typically
// edits (adds entries, moves to a release). A blind block-replace
// would destroy those edits and violate LH-AK-007 ("vorhandene
// Inhalte werden nicht zerstört"). The handler is therefore the
// only M7 generator that does not call [generateManagedFile]; its
// state machine has a different shape:
//
//	absent                          → render full template, write file → Created
//	present-no-block                → ErrGenerateManualConflict (Code 10)
//	present-malformed               → ErrGenerateManualConflict (Code 10)
//	present, block == rendered      → NoOp (idempotent re-run on a fresh file)
//	present, block ≠ rendered       → user-edited; do NOT touch the block.
//	  └─ has `## [Unreleased]` anywhere → NoOp
//	  └─ no `## [Unreleased]`,
//	     has a release section outside the block → insert an
//	     Unreleased stub before the first release section
//	     → RepairedManual
//	  └─ no Unreleased and no release section either → NoOp (the
//	     user has a non-Keep-a-Changelog layout; leave it alone)
//
// Known fragility (slice plan §"Bekannte Fragilität der Hash-
// Heuristik"): any future change to changelog.md.tmpl flips every
// existing project's changelog into the user-edited branch because
// the heuristic is `bytes.Equal(existing, rendered)` — there is no
// template-version marker. M7 freezes the M3 templates; a future
// `--migrate` flag or versioned marker (`init v2`) is V1.
//
// Same fragility on project rename (review-followup S3): when
// `u-boot.yaml.project.name` changes, the block's rendered
// `**<name>**` line no longer matches the existing one, so the
// handler routes to the user-edited branch and leaves the old
// name in place. T2 env-example and T3 readme rewrite the name
// transparently (they call back into the splice path on diff);
// changelog stays stuck on purpose to avoid clobbering user
// entries. The expected workaround is a hand-edit of the existing
// block. CRLF line-endings used to share the same flip — that
// path is now defended in [normaliseLF] (review-followup S2).
func (s *GenerateService) generateChangelog(_ context.Context, req driving.GenerateRequest) (driving.GenerateResponse, error) {
	cfg, err := s.readProjectConfig(req.BaseDir)
	if err != nil {
		return driving.GenerateResponse{}, err
	}

	const relPath = "CHANGELOG.md"
	targetPath := filepath.Join(req.BaseDir, relPath)
	marker := managedblock.Marker{Style: managedblock.StyleHTMLComment, Name: managedblock.InitName}

	rendered, err := renderTemplate("changelog.md.tmpl", templateData{Name: cfg.Project.Name})
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
		s.logger.Info("generate: created",
			"artifact", req.Artifact.String(), "path", relPath, "project", cfg.Project.Name)
		return driving.GenerateResponse{
			Artifact: req.Artifact,
			Action:   driving.GenerateActionCreated,
			Changed:  []string{relPath},
		}, nil
	}

	existing, err := s.fs.ReadFile(targetPath)
	if err != nil {
		return driving.GenerateResponse{}, fmt.Errorf("%w: read %q: %v",
			driving.ErrGenerateFileSystem, targetPath, err)
	}

	renderedBlock, err := renderManagedBlockOnly(rendered, marker)
	if err != nil {
		return driving.GenerateResponse{}, fmt.Errorf(
			"extract init block from rendered changelog.md.tmpl: %w", err)
	}

	start, end, err := classifyExistingBlock(existing, marker, relPath)
	if err != nil {
		return driving.GenerateResponse{}, err
	}

	// State: present, block matches rendered → NoOp (fresh).
	// LF-normalise both sides so a CRLF-edited file does not flip
	// into the user-edited branch (review-followup S2).
	if bytes.Equal(normaliseLF(existing[start:end]), normaliseLF(renderedBlock)) {
		s.logger.Debug("generate: no-op (fresh block)",
			"artifact", req.Artifact.String(), "path", relPath)
		return driving.GenerateResponse{
			Artifact: req.Artifact,
			Action:   driving.GenerateActionNoOp,
			Changed:  nil,
		}, nil
	}

	// State: user-edited block. Never re-render — only consider an
	// Unreleased-stub repair outside the block.
	repaired, doRepair := changelogRepairUnreleased(existing, end)
	if !doRepair {
		s.logger.Debug("generate: no-op (user-edited)",
			"artifact", req.Artifact.String(), "path", relPath)
		return driving.GenerateResponse{
			Artifact: req.Artifact,
			Action:   driving.GenerateActionNoOp,
			Changed:  nil,
		}, nil
	}
	if err := s.fs.WriteFile(targetPath, repaired, defaultFileMode); err != nil {
		return driving.GenerateResponse{}, fmt.Errorf("%w: write %q: %v",
			driving.ErrGenerateFileSystem, targetPath, err)
	}
	s.logger.Info("generate: repaired Unreleased header",
		"artifact", req.Artifact.String(), "path", relPath, "project", cfg.Project.Name)
	return driving.GenerateResponse{
		Artifact: req.Artifact,
		Action:   driving.GenerateActionRepairedManual,
		Changed:  []string{relPath},
	}, nil
}

// changelogRepairUnreleased considers whether to insert a fresh
// `## [Unreleased]` stub before the first release section in the
// user-edited branch of [generateChangelog]. Returns the repaired
// bytes and `true` when a stub was inserted; returns `(nil, false)`
// when no repair is appropriate (Unreleased already present, or no
// release section after the init-block END marker to anchor at).
// blockEnd is the byte offset just past the END-marker line so the
// release-section scan ignores anything inside the managed block.
func changelogRepairUnreleased(existing []byte, blockEnd int) ([]byte, bool) {
	if hasChangelogUnreleased(existing) {
		return nil, false
	}
	tailOffset := firstReleaseSectionOffset(existing[blockEnd:])
	if tailOffset < 0 {
		return nil, false
	}
	insertAt := blockEnd + tailOffset
	repaired := make([]byte, 0, len(existing)+len(changelogUnreleasedStub))
	repaired = append(repaired, existing[:insertAt]...)
	repaired = append(repaired, []byte(changelogUnreleasedStub)...)
	repaired = append(repaired, existing[insertAt:]...)
	return repaired, true
}

// hasChangelogUnreleased reports whether content contains a
// Keep-a-Changelog `## [Unreleased]` header. Used by
// [generateChangelog] to decide between the NoOp and RepairedManual
// user-edited branches. Headers inside a backtick-fenced code block
// are skipped so a documentation example does not falsely register.
func hasChangelogUnreleased(content []byte) bool {
	for _, m := range changelogHeaderRE.FindAllSubmatchIndex(content, -1) {
		if isOffsetInsideFencedBlock(content, m[0]) {
			continue
		}
		name := content[m[2]:m[3]]
		if bytes.Equal(name, []byte("Unreleased")) {
			return true
		}
	}
	return false
}

// firstReleaseSectionOffset returns the byte offset of the first
// `## [<release>]` header in content where `<release>` is not
// `Unreleased`. Returns -1 when no such header is present. Used by
// [generateChangelog] to position the Unreleased stub before the
// existing release history. Headers inside a backtick-fenced code
// block are skipped — splicing the Unreleased stub into a fence
// would corrupt the user's Markdown (review-followup S1).
func firstReleaseSectionOffset(content []byte) int {
	for _, m := range changelogHeaderRE.FindAllSubmatchIndex(content, -1) {
		if isOffsetInsideFencedBlock(content, m[0]) {
			continue
		}
		// m[0]/m[1] cover the full header line offset; m[2]/m[3]
		// cover the captured group (the inner identifier).
		name := content[m[2]:m[3]]
		if !bytes.Equal(name, []byte("Unreleased")) {
			return m[0]
		}
	}
	return -1
}

// isOffsetInsideFencedBlock reports whether the byte at `offset`
// falls inside an open Markdown fenced code block. Counts the
// backtick-fence markers in content[:offset]; an odd count means
// the offset is inside an unclosed fence. Cheap O(N) regex scan
// over the prefix; for the typical sub-kilobyte CHANGELOG.md this
// is negligible.
//
// Limitations (review-followup S1):
//   - Tilde-fenced blocks (`~~~`) are not recognised.
//   - Indented code blocks (4-space indent) are not recognised.
//   - A closing fence with fewer backticks than the opener is
//     counted as a separate fence, which can mis-pair on
//     intentionally-mismatched fences. Both are rare in practice
//     and a future migration to a CommonMark parser would close
//     these gaps.
func isOffsetInsideFencedBlock(content []byte, offset int) bool {
	return len(fenceMarkerRE.FindAllIndex(content[:offset], -1))%2 == 1
}

// generateReadme is the thin T3 wrapper over generateManagedFile for
// the `README.md` artefact (LH-FA-GEN-003): StyleHTMLComment marker
// and the M3 readme.md.tmpl template. User content after the init
// block (Markdown sections the user adds for project-specific
// documentation) is preserved byte-identically per the
// managedblock.Replace contract.
func (s *GenerateService) generateReadme(_ context.Context, req driving.GenerateRequest) (driving.GenerateResponse, error) {
	return s.generateManagedFile(req, "README.md", "readme.md.tmpl", managedblock.StyleHTMLComment)
}

// generateManagedFile is the M7-T2/T3 shared state machine for any
// single-file artefact whose template ships with one `init` managed
// block. Used by `generateEnvExample` (StyleHash, .env.example) and
// `generateReadme` (StyleHTMLComment, README.md). T4 changelog and
// T5 devcontainer have additional concerns (Unreleased-Repair pfad
// and atomic two-file plan respectively) and stay separate.
//
// State table (LH-FA-GEN-001/004/005):
//
//	absent              → render full template, write file → Created
//	present-with-block  → splice re-rendered block → UpdatedBlock or NoOp
//	present-no-block    → ErrGenerateManualConflict (Code 10)
//	present-malformed   → ErrGenerateManualConflict (Code 10, different detail)
//
// Idempotency contract (LH-FA-GEN-005): a second invocation against
// an artefact that already matches the rendered block must return
// [driving.GenerateActionNoOp] with `Changed = nil` and zero
// WriteFile calls — the per-tranche NoOp-pin tests assert both.
// Content outside the BEGIN/END region survives any UpdatedBlock
// splice byte-identically; that is the [managedblock.Replace]
// contract, asserted by the per-tranche content-preservation tests
// (T2: add-on block, T3: user Markdown after init block).
//
// The marker name is fixed to [managedblock.InitName] — all four M7
// artefacts share the `init` block name so that future
// `init --devcontainer` (LH-AK-005, MVP-Closure) can reactivate the
// same block without a parallel marker (see slice-m7-generate.md
// §Architektur-Punkt "Block-Name in allen generierten Dateien").
func (s *GenerateService) generateManagedFile(
	req driving.GenerateRequest,
	relPath, templateName string,
	style managedblock.Style,
) (driving.GenerateResponse, error) {
	cfg, err := s.readProjectConfig(req.BaseDir)
	if err != nil {
		return driving.GenerateResponse{}, err
	}

	targetPath := filepath.Join(req.BaseDir, relPath)
	marker := managedblock.Marker{Style: style, Name: managedblock.InitName}

	rendered, err := renderTemplate(templateName, templateData{Name: cfg.Project.Name})
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
		s.logger.Info("generate: created",
			"artifact", req.Artifact.String(), "path", relPath, "project", cfg.Project.Name)
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
			"extract init block from rendered %s: %w", templateName, err)
	}

	// Classify the existing file's block.
	start, end, err := classifyExistingBlock(existing, marker, relPath)
	if err != nil {
		return driving.GenerateResponse{}, err
	}

	// State: present-with-block. NoOp if the existing block bytes are
	// already equal to the rendered block bytes (idempotency).
	// LF-normalised so a CRLF-edited file does not flip into
	// UpdatedBlock and silently rewrite line endings to LF
	// (review-followup S2).
	if bytes.Equal(normaliseLF(existing[start:end]), normaliseLF(renderedBlock)) {
		s.logger.Debug("generate: no-op",
			"artifact", req.Artifact.String(), "path", relPath)
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
	s.logger.Info("generate: updated block",
		"artifact", req.Artifact.String(), "path", relPath, "project", cfg.Project.Name)
	return driving.GenerateResponse{
		Artifact: req.Artifact,
		Action:   driving.GenerateActionUpdatedBlock,
		Changed:  []string{relPath},
	}, nil
}

// generateEnvExample is the thin T2 wrapper over generateManagedFile
// for the `.env.example` artefact (LH-FA-GEN-004): StyleHash marker
// and the M3 env.example.tmpl template.
func (s *GenerateService) generateEnvExample(_ context.Context, req driving.GenerateRequest) (driving.GenerateResponse, error) {
	return s.generateManagedFile(req, ".env.example", "env.example.tmpl", managedblock.StyleHash)
}

// devcontainerFileAction classifies what generateDevcontainer must do
// for one of the two devcontainer files during its two-phase plan-
// and-execute. Computed during the plan phase; consumed by the
// execute phase. NoOp files are skipped during execute.
type devcontainerFileAction int

const (
	devcontainerActionNoOp devcontainerFileAction = iota
	devcontainerActionWrite        // file absent → write the full rendered template
	devcontainerActionReplaceBlock // file present, block stale → splice
)

// devcontainerFilePlan is the per-file plan computed in phase 1 of
// generateDevcontainer. Phase 2 reads only the fields it needs for
// the chosen action and never re-classifies.
type devcontainerFilePlan struct {
	relPath       string
	targetPath    string
	rendered      []byte
	renderedBlock []byte
	marker        managedblock.Marker
	action        devcontainerFileAction
	existing      []byte // populated only when action == devcontainerActionReplaceBlock
}

// generateDevcontainer implements the M7-T5 two-file state machine
// for `.devcontainer/devcontainer.json` + `.devcontainer/Dockerfile`
// (LH-FA-DEV-001 / LH-FA-DEV-004 / LH-FA-DEV-005 / LH-FA-GEN-005).
// `forwardPorts` is derived from the active services' compose-side
// container ports via the shared `activeServiceNames` +
// `collectActiveServicePorts` helpers — the same source of truth as
// the doctor `devcontainer.forwardPorts.consistency` check. Pinned
// by the T5 anti-drift test in generate_test.go.
//
// Atomic plan-and-execute: phase 1 classifies both files (absent /
// present-with-block-fresh / present-with-block-stale / present-no-
// block / present-malformed) without writing anything. Phase 2
// executes only after phase 1 has confirmed every file is either
// writable or a clean NoOp. If any file is present-no-block or
// malformed, the call returns [driving.ErrGenerateManualConflict]
// with **no** WriteFile invocations — half-written state would be
// re-classified as a conflict on the next call.
func (s *GenerateService) generateDevcontainer(_ context.Context, req driving.GenerateRequest) (driving.GenerateResponse, error) {
	// LH-FA-DEV-003 / Spec §715 — validate the
	// `--allow-external-feature-sources` flag entries early (so a
	// bad URL fails the generate before any FS side effect), but
	// **defer the u-boot.yaml mutation** until after the
	// devcontainer plan-and-execute succeeds. Review-Followup R2:
	// the prior order wrote u-boot.yaml first and produced a
	// mutation leak (allowlist extended + comments lost) when the
	// later plan-and-execute aborted with ErrGenerateManualConflict.
	// T3's renderer does not enforce the allowlist, so the rest of
	// the flow reading an unmodified u-boot.yaml is safe.
	if err := validateAllowExternalFeatureSourcesEntries(req.AllowExternalFeatureSources); err != nil {
		return driving.GenerateResponse{}, err
	}

	cfg, err := s.readProjectConfig(req.BaseDir)
	if err != nil {
		return driving.GenerateResponse{}, err
	}

	ports, err := s.collectDevcontainerForwardPorts(req.BaseDir, cfg)
	if err != nil {
		return driving.GenerateResponse{}, err
	}

	features := collectDevcontainerFeatures(cfg)

	data := templateData{
		Name:         cfg.Project.Name,
		ForwardPorts: ports,
		Features:     features,
	}
	plans, err := s.planDevcontainerFiles(req.BaseDir, data)
	if err != nil {
		return driving.GenerateResponse{}, err
	}

	changed, hasWrite, hasReplace, err := s.executeDevcontainerPlans(plans)
	if err != nil {
		return driving.GenerateResponse{}, err
	}

	// Allowlist write LAST — only after every other write
	// succeeded. Any failure above this point leaves u-boot.yaml
	// byte-identical (no comment loss, no half-mutated state).
	if err := s.applyAllowExternalFeatureSources(req.BaseDir, req.AllowExternalFeatureSources); err != nil {
		return driving.GenerateResponse{}, err
	}

	action := devcontainerAggregateAction(hasWrite, hasReplace)
	if action == driving.GenerateActionNoOp {
		s.logger.Debug("generate devcontainer: no-op",
			"project", cfg.Project.Name,
			"forwardPorts", ports, "features", len(features))
		return driving.GenerateResponse{Artifact: req.Artifact, Action: action}, nil
	}
	s.logger.Info("generate devcontainer: "+action.String(),
		"project", cfg.Project.Name,
		"forwardPorts", ports, "features", len(features),
		"changed", changed)
	return driving.GenerateResponse{
		Artifact: req.Artifact,
		Action:   action,
		Changed:  changed,
	}, nil
}

// collectDevcontainerForwardPorts derives the container-side ports
// of every active service. Reuses the doctor helpers
// activeServiceNames + collectActiveServicePorts so the generator
// and the `devcontainer.forwardPorts.consistency` doctor check
// share a single source of truth (anti-drift pin in the T5 tests).
// Treats a missing compose.yaml as "no ports" rather than an error
// so a project whose compose.yaml has not been generated yet still
// gets a syntactically-valid devcontainer.json (LH-FA-DEV-005
// allows omitting forwardPorts).
//
// Parse-error classification (slice-v1-yaml-parse-error-sentinel):
// when compose.yaml fails to parse, the doctor helper now propagates
// a [driven.ErrYAMLParse]-wrapped error; this caller routes it to
// [driving.ErrGenerateManualConflict] → LH-FA-CLI-006 exit code 10
// (fachlich, user must fix the YAML). All other helper failures
// stay on [driving.ErrGenerateFileSystem] → exit code 14
// (technical).
func (s *GenerateService) collectDevcontainerForwardPorts(baseDir string, cfg ubootYAMLConfig) ([]int, error) {
	services := activeServiceNames(cfg)
	if len(services) == 0 {
		return nil, nil
	}
	composeExists, err := s.fs.Exists(filepath.Join(baseDir, "compose.yaml"))
	if err != nil {
		return nil, fmt.Errorf("%w: Exists(compose.yaml): %v",
			driving.ErrGenerateFileSystem, err)
	}
	if !composeExists {
		return nil, nil
	}
	ports, err := collectActiveServicePorts(s.fs, s.yaml, baseDir, services)
	if err != nil {
		if errors.Is(err, driven.ErrYAMLParse) {
			return nil, fmt.Errorf(
				"%w: compose.yaml is unparseable (%v); repair the YAML manually",
				driving.ErrGenerateManualConflict, err)
		}
		return nil, fmt.Errorf("%w: collectActiveServicePorts: %v",
			driving.ErrGenerateFileSystem, err)
	}
	return ports, nil
}

// planDevcontainerFiles classifies every devcontainer file and
// returns the per-file plan. Fails with ErrGenerateManualConflict on
// the first file whose existing content lacks an init block or has
// a malformed one. No file writes happen here.
func (s *GenerateService) planDevcontainerFiles(baseDir string, data templateData) ([]devcontainerFilePlan, error) {
	specs := []struct {
		relPath  string
		template string
		style    managedblock.Style
	}{
		{".devcontainer/devcontainer.json", "devcontainer/devcontainer.json.tmpl", managedblock.StyleDoubleSlash},
		{".devcontainer/Dockerfile", "devcontainer/Dockerfile.tmpl", managedblock.StyleHash},
	}
	plans := make([]devcontainerFilePlan, 0, len(specs))
	for _, spec := range specs {
		plan, err := s.planDevcontainerFile(baseDir, data, spec.relPath, spec.template, spec.style)
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}
	return plans, nil
}

// planDevcontainerFile classifies a single devcontainer file.
func (s *GenerateService) planDevcontainerFile(
	baseDir string, data templateData,
	relPath, templateName string, style managedblock.Style,
) (devcontainerFilePlan, error) {
	targetPath := filepath.Join(baseDir, relPath)
	marker := managedblock.Marker{Style: style, Name: managedblock.InitName}
	rendered, err := renderTemplate(templateName, data)
	if err != nil {
		return devcontainerFilePlan{}, err
	}
	renderedBlock, err := renderManagedBlockOnly(rendered, marker)
	if err != nil {
		return devcontainerFilePlan{}, fmt.Errorf(
			"extract init block from rendered %s: %w", templateName, err)
	}
	plan := devcontainerFilePlan{
		relPath:       relPath,
		targetPath:    targetPath,
		rendered:      rendered,
		renderedBlock: renderedBlock,
		marker:        marker,
	}
	exists, err := s.fs.Exists(targetPath)
	if err != nil {
		return devcontainerFilePlan{}, fmt.Errorf("%w: Exists(%q): %v",
			driving.ErrGenerateFileSystem, targetPath, err)
	}
	if !exists {
		plan.action = devcontainerActionWrite
		return plan, nil
	}
	existing, err := s.fs.ReadFile(targetPath)
	if err != nil {
		return devcontainerFilePlan{}, fmt.Errorf("%w: read %q: %v",
			driving.ErrGenerateFileSystem, targetPath, err)
	}
	start, end, err := classifyExistingBlock(existing, marker, relPath)
	if err != nil {
		return devcontainerFilePlan{}, err
	}
	// LF-normalise both sides so CRLF-edited Dockerfiles /
	// devcontainer.json files do not flip into UpdatedBlock
	// (review-followup S2).
	if bytes.Equal(normaliseLF(existing[start:end]), normaliseLF(renderedBlock)) {
		plan.action = devcontainerActionNoOp
		return plan, nil
	}
	plan.action = devcontainerActionReplaceBlock
	plan.existing = existing
	return plan, nil
}

// executeDevcontainerPlans runs phase 2 of the two-phase plan-and-
// execute: writes new files, splices existing blocks, skips NoOps.
// Returns the sorted list of changed relative paths plus boolean
// flags driving the aggregate action decision.
func (s *GenerateService) executeDevcontainerPlans(plans []devcontainerFilePlan) (changed []string, hasWrite, hasReplace bool, err error) {
	for _, plan := range plans {
		switch plan.action {
		case devcontainerActionWrite:
			if err := s.writeDevcontainerNewFile(plan); err != nil {
				return nil, false, false, err
			}
			changed = append(changed, plan.relPath)
			hasWrite = true
		case devcontainerActionReplaceBlock:
			if err := s.writeDevcontainerBlockReplace(plan); err != nil {
				return nil, false, false, err
			}
			changed = append(changed, plan.relPath)
			hasReplace = true
		case devcontainerActionNoOp:
			// nothing to do
		default:
			// Defensive: a future devcontainerFileAction value
			// added without updating this switch surfaces as a
			// loud programmer error instead of silently no-opping
			// (review-followup N1).
			return nil, false, false, fmt.Errorf(
				"internal: unknown devcontainerFileAction %d for %q",
				int(plan.action), plan.relPath)
		}
	}
	return changed, hasWrite, hasReplace, nil
}

func (s *GenerateService) writeDevcontainerNewFile(plan devcontainerFilePlan) error {
	dir := filepath.Dir(plan.targetPath)
	if err := s.fs.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("%w: MkdirAll(%q): %v",
			driving.ErrGenerateFileSystem, dir, err)
	}
	if err := s.fs.WriteFile(plan.targetPath, plan.rendered, defaultFileMode); err != nil {
		return fmt.Errorf("%w: write %q: %v",
			driving.ErrGenerateFileSystem, plan.targetPath, err)
	}
	return nil
}

func (s *GenerateService) writeDevcontainerBlockReplace(plan devcontainerFilePlan) error {
	updated, err := managedblock.Replace(plan.existing, plan.marker, plan.renderedBlock)
	if err != nil {
		return fmt.Errorf("replace init block in %q: %w", plan.relPath, err)
	}
	if err := s.fs.WriteFile(plan.targetPath, updated, defaultFileMode); err != nil {
		return fmt.Errorf("%w: write %q: %v",
			driving.ErrGenerateFileSystem, plan.targetPath, err)
	}
	return nil
}

// devcontainerAggregateAction maps the per-file action flags to a
// single response action. UpdatedBlock dominates Created when both
// are present (the slice plan §"partial-clean" row: "Aggregat-Action:
// UpdatedBlock wenn mindestens eine geupdated wurde, sonst Created").
func devcontainerAggregateAction(hasWrite, hasReplace bool) driving.GenerateAction {
	switch {
	case hasReplace:
		return driving.GenerateActionUpdatedBlock
	case hasWrite:
		return driving.GenerateActionCreated
	default:
		return driving.GenerateActionNoOp
	}
}

// validateAllowExternalFeatureSourcesEntries is the early-reject
// validation pass that runs BEFORE any FS read or write in
// generateDevcontainer (Review-Followup R2). Pure function over the
// flag slice so it can fail fast without leaving u-boot.yaml in a
// half-mutated state. The full merge + write happens later via
// [GenerateService.applyAllowExternalFeatureSources] after the
// devcontainer plan-and-execute has succeeded.
func validateAllowExternalFeatureSourcesEntries(sources []string) error {
	if len(sources) == 0 {
		return nil
	}
	if _, err := normaliseFeatureSources(sources); err != nil {
		return fmt.Errorf("generate devcontainer: --allow-external-feature-sources: %w", err)
	}
	return nil
}

// applyAllowExternalFeatureSources implements the Spec §715 wiring
// of `--allow-external-feature-sources` for `generate devcontainer`:
// after the devcontainer plan-and-execute succeeded, append the
// flag URLs to `devcontainer.featureSources.allow` and marshal-
// rewrite u-boot.yaml. No-op when the flag slice is empty.
//
// Failure modes mirror the [ConfigService.setFeatureSourcesAllow]
// list-path: invalid URL → [driving.ErrConfigValueInvalid] (Code 10);
// FS error → [driving.ErrGenerateFileSystem] (Code 14).
func (s *GenerateService) applyAllowExternalFeatureSources(baseDir string, sources []string) error {
	if len(sources) == 0 {
		return nil
	}
	yamlPath := filepath.Join(baseDir, "u-boot.yaml")
	body, err := s.fs.ReadFile(yamlPath)
	if err != nil {
		return fmt.Errorf("%w: read %q: %v",
			driving.ErrGenerateFileSystem, yamlPath, err)
	}
	var cfg ubootYAMLConfig
	if err := s.yaml.Unmarshal(body, &cfg); err != nil {
		return fmt.Errorf("%w: parse %q: %v",
			driving.ErrProjectNotInitialized, yamlPath, err)
	}
	var existing []string
	if cfg.Devcontainer != nil && cfg.Devcontainer.FeatureSources != nil {
		existing = cfg.Devcontainer.FeatureSources.Allow
	}
	merged, err := normaliseFeatureSources(append(append([]string{}, existing...), sources...))
	if err != nil {
		return fmt.Errorf("generate devcontainer: --allow-external-feature-sources: %w", err)
	}
	// NoOp short-circuit — every flag URL already in the list.
	if equalAllowLists(existing, merged) {
		return nil
	}
	if cfg.Devcontainer == nil {
		cfg.Devcontainer = &ubootYAMLDevcontainer{}
	}
	if cfg.Devcontainer.FeatureSources == nil {
		cfg.Devcontainer.FeatureSources = &ubootYAMLFeatureSources{}
	}
	cfg.Devcontainer.FeatureSources.Allow = merged
	rewritten, err := s.yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("%w: marshal u-boot.yaml: %v",
			driving.ErrGenerateFileSystem, err)
	}
	if err := s.fs.WriteFile(yamlPath, rewritten, defaultFileMode); err != nil {
		return fmt.Errorf("%w: write %q: %v",
			driving.ErrGenerateFileSystem, yamlPath, err)
	}
	s.logger.Info("generate devcontainer: allowlist updated",
		"added", len(merged)-len(existing), "total", len(merged))
	return nil
}

// equalAllowLists reports whether a and b contain the same entries
// in the same order (byte-equal slice comparison). Used by the
// generate-devcontainer allowlist-append branch to skip the rewrite
// when the flag URLs were already present.
func equalAllowLists(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
