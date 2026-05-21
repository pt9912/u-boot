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

// projectStructureDirs returns the directories from LH-FA-INIT-003
// that every initialized project gets. The order is the deterministic
// order in which they appear in
// [driving.InitProjectResponse.Created]. Implemented as a function
// (not a package var) to avoid the gochecknoglobals false-positive
// on immutable list constants.
func projectStructureDirs() []string {
	return []string{"docker", "scripts", "docs"}
}

// projectMarkerFiles returns the steering files whose presence in
// req.BaseDir marks the directory as "already an initialized
// u-boot project" (LH-FA-INIT-004 — at least one of these is enough).
func projectMarkerFiles() []string {
	return []string{"u-boot.yaml", "compose.yaml", ".env.example"}
}

// ubootYAMLConfig is the YAML-marshalable shape of u-boot.yaml as
// required by LH-FA-CONF-002 (schemaVersion + project.name). The
// struct lives in the application layer because the YAML schema is
// part of the application contract; the YAMLCodec port stays
// schema-agnostic.
type ubootYAMLConfig struct {
	SchemaVersion int `yaml:"schemaVersion"`
	Project       struct {
		Name string `yaml:"name"`
	} `yaml:"project"`
}

// InitProjectService implements [driving.InitProjectUseCase]. It
// orchestrates the driven ports (FileSystem, YAMLCodec, Git) to
// realize the LH-FA-INIT-001..007 flow.
type InitProjectService struct {
	fs   driven.FileSystem
	yaml driven.YAMLCodec
	git  driven.Git
}

// Static check: InitProjectService satisfies the driving port.
var _ driving.InitProjectUseCase = (*InitProjectService)(nil)

// NewInitProjectService constructs the service with the driven
// adapters injected by the wiring layer (cmd/uboot).
func NewInitProjectService(fs driven.FileSystem, yaml driven.YAMLCodec, git driven.Git) *InitProjectService {
	return &InitProjectService{fs: fs, yaml: yaml, git: git}
}

// Init runs the init flow per LH-FA-INIT-001..007 / LH-FA-CONF-001..003.
// M3-T2 covers the happy path plus the default overwrite-rejection
// branch (LH-FA-INIT-004); --backup / --force handling lands in
// M3-T4.
func (s *InitProjectService) Init(ctx context.Context, req driving.InitProjectRequest) (driving.InitProjectResponse, error) {
	if req.BaseDir == "" {
		return driving.InitProjectResponse{}, errors.New("BaseDir is required")
	}

	name, err := s.resolveProjectName(req)
	if err != nil {
		return driving.InitProjectResponse{}, err
	}

	if err := s.rejectIfExisting(req.BaseDir); err != nil {
		return driving.InitProjectResponse{}, err
	}

	project := domain.NewProject(name)
	dirs := projectStructureDirs()
	templates := fileTemplates()
	created := make([]string, 0, len(dirs)+len(templates)+1)

	dirEntries, err := s.writeDirectories(req.BaseDir)
	if err != nil {
		return driving.InitProjectResponse{}, err
	}
	created = append(created, dirEntries...)

	fileEntries, err := s.writeTemplatedFiles(req.BaseDir, project)
	if err != nil {
		return driving.InitProjectResponse{}, err
	}
	created = append(created, fileEntries...)

	yamlEntry, err := s.writeUBootYAML(req.BaseDir, project)
	if err != nil {
		return driving.InitProjectResponse{}, err
	}
	created = append(created, yamlEntry)

	if !req.SkipGit {
		if err := s.initGit(ctx, req.BaseDir); err != nil {
			return driving.InitProjectResponse{}, err
		}
	}

	return driving.InitProjectResponse{Project: project, Created: created}, nil
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

// rejectIfExisting returns wrapped [driving.ErrProjectExists] when
// any marker file from LH-FA-INIT-004 is present.
func (s *InitProjectService) rejectIfExisting(baseDir string) error {
	for _, marker := range projectMarkerFiles() {
		path := filepath.Join(baseDir, marker)
		exists, err := s.fs.Exists(path)
		if err != nil {
			return fmt.Errorf("check %s: %w", marker, err)
		}
		if exists {
			return fmt.Errorf("%w: %s present", driving.ErrProjectExists, marker)
		}
	}
	return nil
}

// writeDirectories creates the LH-FA-INIT-003 mandatory subdirs.
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

// writeTemplatedFiles renders and writes the embedded templates from
// templates.go (README, CHANGELOG, compose.yaml, .env.example,
// .gitignore).
func (s *InitProjectService) writeTemplatedFiles(baseDir string, project domain.Project) ([]string, error) {
	data := templateData{Name: project.Name.String()}
	templates := fileTemplates()
	created := make([]string, 0, len(templates))
	for _, ft := range templates {
		body, err := renderTemplate(ft.TemplateName, data)
		if err != nil {
			return nil, fmt.Errorf("render %s: %w", ft.Path, err)
		}
		path := filepath.Join(baseDir, ft.Path)
		if err := s.fs.WriteFile(path, body, 0o644); err != nil {
			return nil, fmt.Errorf("write %s: %w", ft.Path, err)
		}
		created = append(created, ft.Path)
	}
	return created, nil
}

// writeUBootYAML marshals and writes u-boot.yaml per LH-FA-CONF-002.
func (s *InitProjectService) writeUBootYAML(baseDir string, project domain.Project) (string, error) {
	cfg := ubootYAMLConfig{SchemaVersion: project.SchemaVersion}
	cfg.Project.Name = project.Name.String()

	body, err := s.yaml.Marshal(cfg)
	if err != nil {
		return "", fmt.Errorf("marshal u-boot.yaml: %w", err)
	}
	path := filepath.Join(baseDir, "u-boot.yaml")
	if err := s.fs.WriteFile(path, body, 0o644); err != nil {
		return "", fmt.Errorf("write u-boot.yaml: %w", err)
	}
	return "u-boot.yaml", nil
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
