package application

import (
	"context"
	"errors"
	"fmt"
	"io"
	iofs "io/fs"
	"path/filepath"
	"sort"

	"github.com/pt9912/u-boot/internal/hexagon/application/managedblock"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// defaultFileMode is the canonical mode for freshly-created u-boot
// project files (LH-FA-INIT-003) — used when the source has no mode
// to preserve (i.e. the file does not yet exist).
const defaultFileMode iofs.FileMode = 0o644

// projectStructureDirs returns the directories from LH-FA-INIT-003
// that every initialized project gets. The order is the deterministic
// order in which they appear in
// [driving.InitProjectResponse.Created]. Implemented as a function
// (not a package var) to avoid the gochecknoglobals false-positive
// on immutable list constants.
func projectStructureDirs() []string {
	return []string{"docker", "scripts", "docs"}
}

// ubootYAMLProject is the `project:` sub-tree of u-boot.yaml
// (LH-FA-CONF-002). It is its own type — not an anonymous nested
// struct — so future additions (description, template reference,
// etc.) stay readable and can be tested in isolation.
type ubootYAMLProject struct {
	Name string `yaml:"name"`
}

// ubootYAMLConfig is the YAML-marshalable shape of u-boot.yaml as
// required by LH-FA-CONF-002 (schemaVersion + project + later
// services + devcontainer + template). The struct lives in the
// application layer because the YAML schema is part of the
// application contract; the YAMLCodec port stays schema-agnostic.
//
// Future M3+/M5+ slices will add:
//
//   - Services    ubootYAMLServices    `yaml:"services,omitempty"`
//   - Devcontainer ubootYAMLDevcontainer `yaml:"devcontainer,omitempty"`
//   - Template    string               `yaml:"template,omitempty"`
type ubootYAMLConfig struct {
	SchemaVersion int              `yaml:"schemaVersion"`
	Project       ubootYAMLProject `yaml:"project"`
}

// InitProjectService implements [driving.InitProjectUseCase]. It
// orchestrates the driven ports (FileSystem, YAMLCodec, Git) to
// realize the LH-FA-INIT-001..007 flow.
type InitProjectService struct {
	fs       driven.FileSystem
	yaml     driven.YAMLCodec
	git      driven.Git
	progress io.Writer
}

// Static check: InitProjectService satisfies the driving port.
var _ driving.InitProjectUseCase = (*InitProjectService)(nil)

// NewInitProjectService constructs the service with the driven
// adapters and a progress writer injected by the wiring layer
// (cmd/uboot). progress receives the LH-FA-INIT-005 §609 / LH-FA-
// CLI-005A §262 "affected paths" summary before any write happens
// on re-init. Pass [io.Discard] for tests that do not assert on
// summary output; nil is treated as Discard for the same reason.
func NewInitProjectService(fs driven.FileSystem, yaml driven.YAMLCodec, git driven.Git, progress io.Writer) *InitProjectService {
	if progress == nil {
		progress = io.Discard
	}
	return &InitProjectService{fs: fs, yaml: yaml, git: git, progress: progress}
}

// fileAction classifies what the service should do with a single
// templated file at execute time. The plan phase computes this for
// every file before any write happens, so a summary can be emitted
// and an abort can fire before partial side effects.
type fileAction int

const (
	// actionWrite means the file does not exist yet — render and
	// write fresh.
	actionWrite fileAction = iota
	// actionReplaceBlock means the file exists with a
	// `U-BOOT MANAGED BLOCK: init` marker; splice in the new block
	// (LH-FA-INIT-005 §613–§614).
	actionReplaceBlock
	// actionOverwriteFull means the file exists and gets fully
	// rewritten. Always paired with backup=true in the plan
	// (LH-FA-INIT-005 §617/§619 require backup before any full
	// overwrite of an existing file).
	actionOverwriteFull
)

// filePlan is the planned action for a single templated file plus
// whether a backup is taken before the action runs. Backup is
// independent of Action — it can be true for both ReplaceBlock and
// OverwriteFull. Body and Mode capture the existing file's content
// and mode at plan time; they are re-used by the execute phase to
// (a) avoid a second Lstat+ReadFile round-trip (and the extra
// TOCTOU window that would open with it), and (b) preserve the
// original file mode across the write — fixing the T4b-review
// mode-regression by mirroring T4a-review's backup-mode policy.
// For actionWrite (file did not exist at plan time) Body is nil
// and Mode is the zero value; executors fall back to
// [defaultFileMode].
type filePlan struct {
	Template fileTemplate
	Action   fileAction
	Backup   bool
	Body     []byte
	Mode     iofs.FileMode
}

// Init runs the init flow per LH-FA-INIT-001..007 / LH-FA-CONF-001..003,
// extended in M3-T4b with the LH-FA-INIT-005 re-init paths
// (--force / --backup, managed-block-only edits).
//
// TOCTOU note: between the plan phase (Lstat + ReadFile) and the
// execute phase (WriteFile / BackupPath), a concurrent process
// could change the file system. For a CLI one-shot the race is
// benign — the worst case is that the execute step trips its own
// error. BackupPath itself is TOCTOU-safe via WriteFileExclusive +
// Mkdir (T4a-review).
func (s *InitProjectService) Init(ctx context.Context, req driving.InitProjectRequest) (driving.InitProjectResponse, error) {
	if req.BaseDir == "" {
		return driving.InitProjectResponse{}, errors.New("BaseDir is required")
	}

	baseExists, err := s.fs.Exists(req.BaseDir)
	if err != nil {
		return driving.InitProjectResponse{}, fmt.Errorf("check BaseDir: %w", err)
	}
	if !baseExists {
		return driving.InitProjectResponse{}, fmt.Errorf("%w: %s", driving.ErrBaseDirMissing, req.BaseDir)
	}

	name, err := s.resolveProjectName(req)
	if err != nil {
		return driving.InitProjectResponse{}, err
	}
	project := domain.NewProject(name)

	// Plan: decide per-file action before writing anything.
	plans, err := s.planTemplatedFiles(req)
	if err != nil {
		return driving.InitProjectResponse{}, err
	}
	yamlPlan, err := s.planUBootYAML(req)
	if err != nil {
		return driving.InitProjectResponse{}, err
	}

	// Summary: emit before any side effect (LH-FA-INIT-005 §609).
	s.emitSummary(req.BaseDir, plans, yamlPlan)

	// Execute.
	created := make([]string, 0)
	backups := make([]driving.BackupAction, 0)

	dirEntries, err := s.writeDirectories(req.BaseDir)
	if err != nil {
		return driving.InitProjectResponse{}, err
	}
	created = append(created, dirEntries...)

	fileEntries, fileBackups, err := s.executeTemplatedFiles(req.BaseDir, project, plans)
	if err != nil {
		return driving.InitProjectResponse{}, err
	}
	created = append(created, fileEntries...)
	backups = append(backups, fileBackups...)

	yamlEntry, yamlBackup, err := s.executeUBootYAML(req.BaseDir, project, yamlPlan)
	if err != nil {
		return driving.InitProjectResponse{}, err
	}
	created = append(created, yamlEntry)
	if yamlBackup != nil {
		backups = append(backups, *yamlBackup)
	}

	if !req.SkipGit {
		if err := s.initGit(ctx, req.BaseDir); err != nil {
			return driving.InitProjectResponse{}, err
		}
	}

	return driving.InitProjectResponse{Project: project, Created: created, Backups: backups}, nil
}

// resolveProjectName derives and validates the project name per
// LH-FA-INIT-002 / LH-FA-INIT-006.
func (s *InitProjectService) resolveProjectName(req driving.InitProjectRequest) (domain.ProjectName, error) {
	raw := req.Name
	if raw == "" {
		raw = domain.NormalizeProjectName(filepath.Base(req.BaseDir))
	}
	name, err := domain.NewProjectName(raw)
	if err != nil {
		return "", err
	}
	return name, nil
}

// planTemplatedFiles computes the per-file plan for every templated
// file (README, CHANGELOG, compose.yaml, .env.example, .gitignore).
// Returns the first abort-error encountered, so no side effect runs.
func (s *InitProjectService) planTemplatedFiles(req driving.InitProjectRequest) ([]filePlan, error) {
	templates := fileTemplates()
	plans := make([]filePlan, 0, len(templates))
	for _, ft := range templates {
		fp, err := s.planFile(req.BaseDir, ft, req.Force, req.Backup)
		if err != nil {
			return nil, err
		}
		plans = append(plans, fp)
	}
	return plans, nil
}

// planUBootYAML computes the plan for u-boot.yaml. The file is
// treated as fully managed (no managed-block marker support — per
// LH-SA-FILE-002 §615 strict-YAML / steering file), so a re-init
// without --backup always aborts.
func (s *InitProjectService) planUBootYAML(req driving.InitProjectRequest) (filePlan, error) {
	ft := fileTemplate{Path: "u-boot.yaml", TemplateName: "", Managed: false}
	return s.planFile(req.BaseDir, ft, req.Force, req.Backup)
}

// planFile applies the LH-FA-INIT-005 decision tree to a single file:
//
//   - file does not exist                        → actionWrite
//   - is symlink                                 → ErrBackupUnsupportedKind
//   - exists + (--force AND has managed block)   → actionReplaceBlock
//   - exists + --backup                          → actionOverwriteFull (with backup)
//   - exists + --force (no block, no --backup)   → ErrForceRequiresBackup
//   - exists + (no --force, no --backup)         → ErrProjectExists / ErrFileExists
//
// The backup flag in the returned plan is true whenever --backup is
// set AND the action mutates the file, so a managed-block-only edit
// with --backup still gets a safety copy. For an existing file
// planFile captures Mode + (for Managed templates) Body in the
// returned plan so the execute phase can preserve mode and avoid a
// second read.
//
// Symlinks are rejected with [driving.ErrBackupUnsupportedKind] — the
// same sentinel used by [BackupPath] in T4a, for the same reason:
// silently following a `.env.example -> /etc/passwd` symlink would
// have the re-init read and overwrite the link target instead of
// the link itself.
//
// Collision errors split the spec-§604 marker files (u-boot.yaml,
// compose.yaml, .env.example) from the rest: marker collisions
// produce [driving.ErrProjectExists] (the directory really is an
// existing u-boot project), non-marker collisions produce
// [driving.ErrFileExists] (a stray README.md does not prove a
// u-boot project).
func (s *InitProjectService) planFile(baseDir string, ft fileTemplate, force, backup bool) (filePlan, error) {
	fp := filePlan{Template: ft}
	fullPath := filepath.Join(baseDir, ft.Path)

	info, err := s.fs.Lstat(fullPath)
	if err != nil {
		if errors.Is(err, iofs.ErrNotExist) {
			fp.Action = actionWrite
			return fp, nil
		}
		return fp, fmt.Errorf("lstat %s: %w", ft.Path, err)
	}
	if info.Mode()&iofs.ModeSymlink != 0 {
		return fp, fmt.Errorf("%w: %s is a symlink", driving.ErrBackupUnsupportedKind, ft.Path)
	}
	fp.Mode = info.Mode().Perm()

	hasBlock, body, err := s.fileHasManagedBlock(fullPath, ft)
	if err != nil {
		return fp, err
	}
	fp.Body = body

	switch {
	case force && hasBlock:
		fp.Action = actionReplaceBlock
		fp.Backup = backup
		return fp, nil
	case backup:
		fp.Action = actionOverwriteFull
		fp.Backup = true
		return fp, nil
	case force:
		return fp, fmt.Errorf("%w: %s exists without a managed block; add --backup to overwrite",
			driving.ErrForceRequiresBackup, ft.Path)
	default:
		return fp, collisionError(ft.Path)
	}
}

// collisionError picks the right sentinel for a re-init collision.
// LH-FA-INIT-004 marker files (u-boot.yaml, compose.yaml,
// .env.example) signal a genuine existing u-boot project →
// [driving.ErrProjectExists]; anything else (a stray README.md in a
// non-u-boot directory) signals only a name collision →
// [driving.ErrFileExists]. Both map to exit code 10; splitting them
// keeps the CLI message faithful to what the tool actually
// observed.
func collisionError(path string) error {
	for _, marker := range []string{"u-boot.yaml", "compose.yaml", ".env.example"} {
		if path == marker {
			return fmt.Errorf("%w: %s present", driving.ErrProjectExists, path)
		}
	}
	return fmt.Errorf("%w: %s present; pass --backup or --force to proceed", driving.ErrFileExists, path)
}

// fileHasManagedBlock reports whether the existing file at fullPath
// contains the `U-BOOT MANAGED BLOCK: init` marker for the
// template's declared style, and returns the file body so the
// caller can cache it in [filePlan]. Returns (false, nil, nil) for
// non-Managed templates (e.g. .gitignore) — these never have
// inline block markers and the execute path will only need the
// body if a full overwrite happens (whole-file backup via
// [BackupPath] reads disk directly; no in-process body needed).
func (s *InitProjectService) fileHasManagedBlock(fullPath string, ft fileTemplate) (bool, []byte, error) {
	if !ft.Managed {
		return false, nil, nil
	}
	content, err := s.fs.ReadFile(fullPath)
	if err != nil {
		return false, nil, fmt.Errorf("read %s: %w", ft.Path, err)
	}
	marker := managedblock.Marker{Style: ft.Style, Name: managedblock.InitName}
	return managedblock.Has(content, marker), content, nil
}

// emitSummary writes the LH-FA-INIT-005 §609 / LH-FA-CLI-005A §262
// affected-paths summary to the configured progress writer. Only
// fires when at least one plan would *replace a block* or *fully
// overwrite* a file — purely additive runs (fresh init, all
// actionWrite) stay quiet. If a future action ever mutates a file
// without falling into ReplaceBlock/OverwriteFull, extend the
// classifier in [shouldSummarize] below.
func (s *InitProjectService) emitSummary(baseDir string, plans []filePlan, yamlPlan filePlan) {
	type affected struct {
		Path   string
		Action string
		Backup bool
	}
	var rows []affected
	collect := func(p filePlan) {
		switch p.Action {
		case actionReplaceBlock:
			rows = append(rows, affected{Path: p.Template.Path, Action: "replace managed block", Backup: p.Backup})
		case actionOverwriteFull:
			rows = append(rows, affected{Path: p.Template.Path, Action: "full overwrite", Backup: p.Backup})
		}
	}
	for _, p := range plans {
		collect(p)
	}
	collect(yamlPlan)
	if len(rows) == 0 {
		return
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].Path < rows[j].Path })
	fmt.Fprintf(s.progress, "Affected files in %s:\n", baseDir)
	for _, r := range rows {
		marker := ""
		if r.Backup {
			marker = " (with backup)"
		}
		fmt.Fprintf(s.progress, "  - %s — %s%s\n", r.Path, r.Action, marker)
	}
}

// writeDirectories creates the LH-FA-INIT-003 mandatory subdirs.
// MkdirAll is idempotent, so re-init on an existing project just
// re-creates the dirs (no-op on disk).
func (s *InitProjectService) writeDirectories(baseDir string) ([]string, error) {
	dirs := projectStructureDirs()
	created := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		path := filepath.Join(baseDir, dir)
		if err := s.fs.MkdirAll(path, 0o755); err != nil {
			return nil, fmt.Errorf("mkdir %s: %w", dir, err)
		}
		created = append(created, dir+"/")
	}
	return created, nil
}

// executeTemplatedFiles runs the plan for every templated file:
// render the template, then dispatch to the action-specific helper
// ([writeNewFile], [replaceManagedBlock], [backupAndOverwrite]).
func (s *InitProjectService) executeTemplatedFiles(baseDir string, project domain.Project, plans []filePlan) ([]string, []driving.BackupAction, error) {
	data := templateData{Name: project.Name.String()}
	created := make([]string, 0, len(plans))
	backups := make([]driving.BackupAction, 0)
	for _, p := range plans {
		body, err := renderTemplate(p.Template.TemplateName, data)
		if err != nil {
			return nil, nil, fmt.Errorf("render %s: %w", p.Template.Path, err)
		}
		entry, backup, err := s.executeFile(baseDir, p, body)
		if err != nil {
			return nil, nil, err
		}
		created = append(created, entry)
		if backup != nil {
			backups = append(backups, *backup)
		}
	}
	return created, backups, nil
}

// executeFile applies one filePlan to disk and returns the written
// path (relative) plus any backup action taken. Mode preservation:
// for re-init actions (ReplaceBlock / OverwriteFull) the existing
// file's mode is captured in plan.Mode by planFile and used for
// WriteFile, so a `chmod 600 .env.example` survives the round-trip.
// actionWrite (new file) falls back to [defaultFileMode].
func (s *InitProjectService) executeFile(baseDir string, plan filePlan, body []byte) (string, *driving.BackupAction, error) {
	fullPath := filepath.Join(baseDir, plan.Template.Path)
	switch plan.Action {
	case actionWrite:
		if err := s.fs.WriteFile(fullPath, body, defaultFileMode); err != nil {
			return "", nil, fmt.Errorf("write %s: %w", plan.Template.Path, err)
		}
		return plan.Template.Path, nil, nil
	case actionReplaceBlock:
		return s.executeReplaceBlock(baseDir, plan, body)
	case actionOverwriteFull:
		return s.executeOverwriteFull(baseDir, plan, body)
	}
	return "", nil, fmt.Errorf("unknown action %d for %s", plan.Action, plan.Template.Path)
}

// executeReplaceBlock splices the rendered managed block into the
// existing file body (captured by planFile in plan.Body to avoid a
// second read + extra TOCTOU window) and writes the result with the
// preserved mode. Optionally backs up the original first when
// plan.Backup is true.
func (s *InitProjectService) executeReplaceBlock(baseDir string, plan filePlan, body []byte) (string, *driving.BackupAction, error) {
	fullPath := filepath.Join(baseDir, plan.Template.Path)
	var backup *driving.BackupAction
	if plan.Backup {
		ba, err := s.runBackup(baseDir, plan.Template.Path)
		if err != nil {
			return "", nil, err
		}
		backup = ba
	}
	marker := managedblock.Marker{Style: plan.Template.Style, Name: managedblock.InitName}
	updated, err := managedblock.Replace(plan.Body, marker, body)
	if err != nil {
		return "", nil, fmt.Errorf("replace block in %s: %w", plan.Template.Path, err)
	}
	if err := s.fs.WriteFile(fullPath, updated, plan.Mode); err != nil {
		return "", nil, fmt.Errorf("write %s: %w", plan.Template.Path, err)
	}
	return plan.Template.Path, backup, nil
}

// executeOverwriteFull backs up the existing file (Backup is always
// true for this action) and then writes the freshly rendered body
// over the whole file, preserving the captured mode.
func (s *InitProjectService) executeOverwriteFull(baseDir string, plan filePlan, body []byte) (string, *driving.BackupAction, error) {
	fullPath := filepath.Join(baseDir, plan.Template.Path)
	ba, err := s.runBackup(baseDir, plan.Template.Path)
	if err != nil {
		return "", nil, err
	}
	if err := s.fs.WriteFile(fullPath, body, plan.Mode); err != nil {
		return "", nil, fmt.Errorf("write %s: %w", plan.Template.Path, err)
	}
	return plan.Template.Path, ba, nil
}

// runBackup wraps [BackupPath] and returns the resulting
// [driving.BackupAction] record for the response.
func (s *InitProjectService) runBackup(baseDir, relPath string) (*driving.BackupAction, error) {
	fullPath := filepath.Join(baseDir, relPath)
	backupPath, err := BackupPath(s.fs, fullPath)
	if err != nil {
		return nil, fmt.Errorf("backup %s: %w", relPath, err)
	}
	return &driving.BackupAction{Original: relPath, Backup: backupPath}, nil
}

// executeUBootYAML marshals and writes u-boot.yaml per
// LH-FA-CONF-002 with the same plan dispatch as the templated files.
// u-boot.yaml is whole-file managed (no inline block marker), so the
// only re-init action is OverwriteFull (with backup).
func (s *InitProjectService) executeUBootYAML(baseDir string, project domain.Project, plan filePlan) (string, *driving.BackupAction, error) {
	cfg := ubootYAMLConfig{
		SchemaVersion: project.SchemaVersion,
		Project:       ubootYAMLProject{Name: project.Name.String()},
	}
	body, err := s.yaml.Marshal(cfg)
	if err != nil {
		return "", nil, fmt.Errorf("marshal u-boot.yaml: %w", err)
	}
	return s.executeFile(baseDir, plan, body)
}

// initGit runs the LH-FA-INIT-007 default path: when BaseDir is not
// yet a git repo, run `git init`. Already-initialized repos are
// left alone.
func (s *InitProjectService) initGit(ctx context.Context, baseDir string) error {
	isRepo, err := s.git.IsRepository(ctx, baseDir)
	if err != nil {
		return fmt.Errorf("git is-repository: %w", err)
	}
	if isRepo {
		return nil
	}
	if err := s.git.Init(ctx, baseDir); err != nil {
		return fmt.Errorf("git init: %w", err)
	}
	return nil
}
