// Package templateyaml parses and validates a single external
// template's `template.yaml` metadata file per ADR-0009 §Entscheidung.
//
// It is the shared parser for every TemplateCatalog / TemplateFiles
// adapter: the embed.FS-backed `externaltemplates` catalog and the
// filesystem-backed `localtemplates` resolver (slice-later-local-
// templates) both call [Read], so the apiVersion gate, strict-field
// decoding, and domain validation live in exactly one place. Keeping
// the parser in its own neutral package — instead of exporting it from
// `externaltemplates` — avoids adapter-to-adapter coupling
// (slice-later-local-templates T0-(c)).
package templateyaml

import (
	"bytes"
	"fmt"
	iofs "io/fs"
	"path"

	"gopkg.in/yaml.v3"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// SupportedAPIVersion is the only `apiVersion` value the parser
// accepts. ADR-0009 §Entscheidung pins the v1 schema; future schema
// evolution must land in a slice that bumps this constant (and either
// rejects or migrates v1-shaped files). The gate runs at parse time so
// a packaging mistake — or a template author copying a v2 example —
// fails fast with [domain.ErrInvalidTemplate] rather than silently
// rendering under wrong assumptions.
const SupportedAPIVersion = "github.com/pt9912/u-boot/template/v1"

// MetadataFile is the per-template metadata filename. ADR-0009
// §Entscheidung pins it; any template directory without this file is
// treated as an invalid template by its resolver.
const MetadataFile = "template.yaml"

// Read parses the `template.yaml` under dir in fs, enforces the
// apiVersion gate plus strict (KnownFields) decoding, maps the result
// to the domain projection, and validates the LH-FA-TPL-002 metadata
// minimum.
//
// Strict decoding: the decoder runs with KnownFields(true) so a typo
// like `requiredTool:` (singular) or `addOns:` (instead of
// `supportedAddOns:`) fails at parse time instead of silently dropping
// the author's intent.
//
// Errors:
//
//   - read failure → wrapped [iofs.PathError].
//   - malformed YAML / unknown field → wrapped parse error.
//   - apiVersion ≠ [SupportedAPIVersion] → [domain.ErrInvalidTemplate].
//   - metadata fails [domain.TemplateMetadata.Validate] →
//     [domain.ErrInvalidTemplate].
//
// Read stays layer-neutral: it returns [domain.ErrInvalidTemplate]
// (the build-time / catalog class) and plain parse errors. Adapters
// that need a user-facing exit-10 class (the filesystem resolver) wrap
// the result with their own driven sentinel
// (`driven.ErrTemplateInvalid`).
func Read(fs iofs.FS, dir string) (domain.TemplateMetadata, error) {
	metaPath := path.Join(dir, MetadataFile)
	data, err := iofs.ReadFile(fs, metaPath)
	if err != nil {
		return domain.TemplateMetadata{}, fmt.Errorf("read %s: %w", metaPath, err)
	}
	var raw rawTemplateYAML
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(&raw); err != nil {
		return domain.TemplateMetadata{}, fmt.Errorf("parse %s: %w", metaPath, err)
	}
	if raw.APIVersion != SupportedAPIVersion {
		return domain.TemplateMetadata{}, fmt.Errorf("%s: %w: apiVersion %q is not supported (want %q)",
			metaPath, domain.ErrInvalidTemplate, raw.APIVersion, SupportedAPIVersion)
	}
	meta := raw.toDomain()
	if err := meta.Validate(); err != nil {
		return domain.TemplateMetadata{}, fmt.Errorf("%s: %w", metaPath, err)
	}
	return meta, nil
}

// rawTemplateYAML mirrors the on-disk schema per ADR-0009
// §Entscheidung. Private so this package is the only place that knows
// the `apiVersion`-tagged YAML shape; consumers see only the curated
// domain projection.
type rawTemplateYAML struct {
	APIVersion      string           `yaml:"apiVersion"`
	Name            string           `yaml:"name"`
	Description     string           `yaml:"description"`
	Version         string           `yaml:"version"`
	SupportedAddOns []string         `yaml:"supportedAddOns"`
	GeneratedFiles  []string         `yaml:"generatedFiles"`
	RequiredTools   []string         `yaml:"requiredTools"`
	Variables       []rawTemplateVar `yaml:"variables"`
}

// rawTemplateVar is the on-disk variable schema.
type rawTemplateVar struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Default     string `yaml:"default"`
	Required    bool   `yaml:"required"`
}

// toDomain copies the raw YAML projection into the domain value. Flat
// copy by design — any future schema evolution (vX → vX+1) lives in
// this package so the domain shape stays stable.
func (r rawTemplateYAML) toDomain() domain.TemplateMetadata {
	out := domain.TemplateMetadata{
		Name:            r.Name,
		Description:     r.Description,
		Version:         r.Version,
		SupportedAddOns: r.SupportedAddOns,
		GeneratedFiles:  r.GeneratedFiles,
		RequiredTools:   r.RequiredTools,
	}
	if len(r.Variables) > 0 {
		out.Variables = make([]domain.TemplateVariable, len(r.Variables))
		for i, v := range r.Variables {
			out.Variables[i] = domain.TemplateVariable{
				Name:        v.Name,
				Description: v.Description,
				Default:     v.Default,
				Required:    v.Required,
			}
		}
	}
	return out
}
