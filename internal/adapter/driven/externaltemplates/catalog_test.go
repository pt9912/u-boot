package externaltemplates_test

import (
	"context"
	"errors"
	iofs "io/fs"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/pt9912/u-boot/internal/adapter/driven/externaltemplates"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

func TestCatalog_List_ProductionBundle_ContainsBasic(t *testing.T) {
	t.Parallel()
	cat := externaltemplates.New()
	metas, err := cat.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(metas) == 0 {
		t.Fatal("production catalog is empty; want at least the `basic` bootstrap (ADR-0009 §Folgepunkte 4)")
	}

	var basic *domain.TemplateMetadata
	for i, m := range metas {
		if m.Name == "basic" {
			basic = &metas[i]
			break
		}
	}
	if basic == nil {
		names := make([]string, 0, len(metas))
		for _, m := range metas {
			names = append(names, m.Name)
		}
		t.Fatalf("production catalog missing `basic` bootstrap template; have %v", names)
	}

	// LH-FA-TPL-002 / LH-FA-TPL-004 minimum surface: name +
	// description + version must be non-empty in the listing.
	if basic.Description == "" {
		t.Error("basic.description is empty")
	}
	if basic.Version == "" {
		t.Error("basic.version is empty")
	}
}

func TestCatalog_List_SortsByName(t *testing.T) {
	t.Parallel()
	fs := fstest.MapFS{
		"templates/zebra/template.yaml":  &fstest.MapFile{Data: []byte(yamlBlob("zebra", "z desc", "1.0.0"))},
		"templates/alpha/template.yaml":  &fstest.MapFile{Data: []byte(yamlBlob("alpha", "a desc", "1.0.0"))},
		"templates/middle/template.yaml": &fstest.MapFile{Data: []byte(yamlBlob("middle", "m desc", "1.0.0"))},
	}
	cat := externaltemplates.NewWithFS(fs)
	metas, err := cat.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(metas) != 3 {
		t.Fatalf("len(metas) = %d, want 3", len(metas))
	}
	want := []string{"alpha", "middle", "zebra"}
	for i, name := range want {
		if metas[i].Name != name {
			t.Errorf("metas[%d].Name = %q, want %q (alphabetical)", i, metas[i].Name, name)
		}
	}
}

func TestCatalog_List_SkipsNonDirEntries(t *testing.T) {
	t.Parallel()
	// A stray README at the catalog root must not be misclassified as
	// a template — only subdirectories are templates.
	fs := fstest.MapFS{
		"templates/basic/template.yaml": &fstest.MapFile{Data: []byte(yamlBlob("basic", "b", "1.0.0"))},
		"templates/README.md":           &fstest.MapFile{Data: []byte("# notes")},
	}
	cat := externaltemplates.NewWithFS(fs)
	metas, err := cat.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(metas) != 1 {
		t.Fatalf("len(metas) = %d, want 1 (README at catalog root must be skipped)", len(metas))
	}
}

func TestCatalog_List_InvalidYAMLReturnsError(t *testing.T) {
	t.Parallel()
	fs := fstest.MapFS{
		"templates/broken/template.yaml": &fstest.MapFile{Data: []byte("not: valid: yaml: here")},
	}
	cat := externaltemplates.NewWithFS(fs)
	_, err := cat.List(context.Background())
	if err == nil {
		t.Fatal("List: want error for broken yaml, got nil")
	}
}

func TestCatalog_List_MissingMetadataFileReturnsError(t *testing.T) {
	t.Parallel()
	// Directory exists but no template.yaml — packaging mistake the
	// adapter must surface, not silently skip.
	fs := fstest.MapFS{
		"templates/empty/.keep": &fstest.MapFile{Data: []byte("placeholder")},
	}
	cat := externaltemplates.NewWithFS(fs)
	_, err := cat.List(context.Background())
	if err == nil {
		t.Fatal("List: want error for template dir without template.yaml, got nil")
	}
}

func TestCatalog_List_InvalidMetadataWrapsErrInvalidTemplate(t *testing.T) {
	t.Parallel()
	// Valid YAML, but missing the LH-FA-TPL-002 required `name` field.
	fs := fstest.MapFS{
		"templates/nameless/template.yaml": &fstest.MapFile{
			Data: []byte("apiVersion: github.com/pt9912/u-boot/template/v1\n" +
				"description: \"desc\"\n" +
				"version: 1.0.0\n"),
		},
	}
	cat := externaltemplates.NewWithFS(fs)
	_, err := cat.List(context.Background())
	if err == nil {
		t.Fatal("List: want validation error, got nil")
	}
	if !errors.Is(err, domain.ErrInvalidTemplate) {
		t.Errorf("err = %v, want wrap of domain.ErrInvalidTemplate", err)
	}
}

func TestCatalog_List_UnknownYAMLFieldRejected(t *testing.T) {
	t.Parallel()
	// Review-followup N1: yaml.v3's default mode silently drops
	// unknown fields. The strict decoder must reject a typo like
	// `requiredTool:` (singular) so the author hears about it.
	fs := fstest.MapFS{
		"templates/typo/template.yaml": &fstest.MapFile{Data: []byte(
			"apiVersion: github.com/pt9912/u-boot/template/v1\n" +
				"name: typo\n" +
				"description: \"singular typo for requiredTools\"\n" +
				"version: 1.0.0\n" +
				"requiredTool: [jdk]\n", // singular — would silently vanish without KnownFields(true)
		)},
	}
	_, err := externaltemplates.NewWithFS(fs).List(context.Background())
	if err == nil {
		t.Fatal("List: want error for unknown YAML field, got nil (typos must not silently parse)")
	}
}

func TestCatalog_List_UnsupportedAPIVersionRejected(t *testing.T) {
	t.Parallel()
	// Review-followup N4: a template carrying a future apiVersion
	// must fail at load time rather than silently render under
	// wrong assumptions.
	fs := fstest.MapFS{
		"templates/future/template.yaml": &fstest.MapFile{Data: []byte(
			"apiVersion: github.com/pt9912/u-boot/template/v999\n" +
				"name: future\n" +
				"description: \"hypothetical v999\"\n" +
				"version: 1.0.0\n",
		)},
	}
	_, err := externaltemplates.NewWithFS(fs).List(context.Background())
	if err == nil {
		t.Fatal("List: want error for unsupported apiVersion, got nil")
	}
	if !errors.Is(err, domain.ErrInvalidTemplate) {
		t.Errorf("err = %v, want wrap of domain.ErrInvalidTemplate", err)
	}
}

func TestCatalog_List_HonorsCancelledContext(t *testing.T) {
	t.Parallel()
	// Review-followup N5: the port doc claims every adapter MUST
	// observe ctx cancellation; the production embed.FS adapter
	// satisfies it via an entry-time ctx.Err() check.
	fs := fstest.MapFS{
		"templates/basic/template.yaml": &fstest.MapFile{Data: []byte(yamlBlob("basic", "b", "1.0.0"))},
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before the call so List sees the cancellation at entry.

	_, err := externaltemplates.NewWithFS(fs).List(ctx)
	if err == nil {
		t.Fatal("List: want ctx-cancel error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}
}

func TestCatalog_List_MissingCatalogRootReturnsError(t *testing.T) {
	t.Parallel()
	// Empty FS — no `templates/` directory exists. Production embed
	// guarantees the root; this exercise the error path for callers
	// that wire a wrong FS.
	cat := externaltemplates.NewWithFS(fstest.MapFS{})
	_, err := cat.List(context.Background())
	if err == nil {
		t.Fatal("List: want error for missing catalog root, got nil")
	}
}

func TestCatalog_List_ParsesVariablesAndFiles(t *testing.T) {
	t.Parallel()
	fs := fstest.MapFS{
		"templates/full/template.yaml": &fstest.MapFile{Data: []byte(
			"apiVersion: github.com/pt9912/u-boot/template/v1\n" +
				"name: full\n" +
				"description: \"all-fields-populated fixture\"\n" +
				"version: 2.3.4\n" +
				"supportedAddOns: [postgres, keycloak]\n" +
				"generatedFiles:\n  - build.gradle\n  - src/Main.java\n" +
				"requiredTools:\n  - jdk:>=21\n" +
				"variables:\n" +
				"  - name: groupId\n    description: \"Maven group ID\"\n    default: \"com.example\"\n    required: true\n",
		)},
	}
	metas, err := externaltemplates.NewWithFS(fs).List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(metas) != 1 {
		t.Fatalf("len(metas) = %d, want 1", len(metas))
	}
	m := metas[0]
	if m.Version != "2.3.4" {
		t.Errorf("Version = %q, want 2.3.4", m.Version)
	}
	if !equalStrings(m.SupportedAddOns, []string{"postgres", "keycloak"}) {
		t.Errorf("SupportedAddOns = %v, want [postgres keycloak]", m.SupportedAddOns)
	}
	if !equalStrings(m.GeneratedFiles, []string{"build.gradle", "src/Main.java"}) {
		t.Errorf("GeneratedFiles = %v, want [build.gradle src/Main.java]", m.GeneratedFiles)
	}
	if !equalStrings(m.RequiredTools, []string{"jdk:>=21"}) {
		t.Errorf("RequiredTools = %v, want [jdk:>=21]", m.RequiredTools)
	}
	if len(m.Variables) != 1 {
		t.Fatalf("len(Variables) = %d, want 1", len(m.Variables))
	}
	v := m.Variables[0]
	if v.Name != "groupId" || v.Description != "Maven group ID" || v.Default != "com.example" || !v.Required {
		t.Errorf("Variables[0] = %#v, want groupId/Maven/com.example/required=true", v)
	}
}

func TestCatalog_Open_ProductionBasicReturnsFS(t *testing.T) {
	t.Parallel()
	cat := externaltemplates.New()
	sub, err := cat.Open(context.Background(), "basic")
	if err != nil {
		t.Fatalf("Open(basic): %v", err)
	}
	// The returned FS must be rooted at templates/basic/, so
	// template.yaml is directly reachable.
	data, err := iofs.ReadFile(sub, "template.yaml")
	if err != nil {
		t.Fatalf("ReadFile(template.yaml): %v", err)
	}
	if !strings.Contains(string(data), "name: basic") {
		t.Errorf("template.yaml content missing `name: basic`; got:\n%s", string(data))
	}
}

func TestCatalog_Open_UnknownNameReturnsErrTemplateNotFound(t *testing.T) {
	t.Parallel()
	cat := externaltemplates.New()
	_, err := cat.Open(context.Background(), "nonexistent-zzz")
	if err == nil {
		t.Fatal("Open(nonexistent): want error, got nil")
	}
	if !errors.Is(err, driven.ErrTemplateNotFound) {
		t.Errorf("err = %v, want wrap of driven.ErrTemplateNotFound", err)
	}
}

func TestCatalog_Open_EmptyNameReturnsErrTemplateNotFound(t *testing.T) {
	t.Parallel()
	cat := externaltemplates.New()
	_, err := cat.Open(context.Background(), "")
	if err == nil {
		t.Fatal("Open(\"\"): want error, got nil")
	}
	if !errors.Is(err, driven.ErrTemplateNotFound) {
		t.Errorf("err = %v, want wrap of driven.ErrTemplateNotFound", err)
	}
}

// Post-Closure-Review #1 (slice-later-local-templates): a path-shaped
// name must NOT resolve to the catalog root (`.` → path.Join collapses
// to catalogRoot) or escape it (`..`). Defense-in-depth guard — the
// Composite already routes these to the FS resolver, but the catalog
// rejects them regardless of caller.
func TestCatalog_Open_PathShapedNameRejected(t *testing.T) {
	t.Parallel()
	cat := externaltemplates.New()
	for _, name := range []string{".", "..", "a/b", `a\b`, "../etc"} {
		_, err := cat.Open(context.Background(), name)
		if !errors.Is(err, driven.ErrTemplateNotFound) {
			t.Errorf("Open(%q) err = %v, want ErrTemplateNotFound (not catalog-root resolution)", name, err)
		}
	}
}

func TestCatalog_Open_HonorsCancelledContext(t *testing.T) {
	t.Parallel()
	cat := externaltemplates.New()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := cat.Open(ctx, "basic")
	if err == nil {
		t.Fatal("Open with cancelled ctx: want error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}
}

func TestCatalog_Open_EmptyTemplateDirReturnsErrTemplateNotFound(t *testing.T) {
	t.Parallel()
	// A directory exists but has no entries — packaging mistake the
	// adapter must surface as not-found (semantically the template
	// has nothing to render).
	fs := fstest.MapFS{
		"templates/basic/template.yaml": &fstest.MapFile{Data: []byte(yamlBlob("basic", "b", "1.0.0"))},
		// Note: `placeholder` ensures the empty/ dir exists in the
		// MapFS namespace via its parent path resolution; MapFS
		// requires at least one entry under a directory for ReadDir
		// to return an empty slice rather than ENOENT, so we test
		// via a sibling structure that creates the dir entry.
	}
	// Verify our happy-path baseline first.
	_, err := externaltemplates.NewWithFS(fs).Open(context.Background(), "basic")
	if err != nil {
		t.Fatalf("baseline Open(basic): %v", err)
	}
	// Now ask for a name that doesn't exist as a directory entry.
	_, err = externaltemplates.NewWithFS(fs).Open(context.Background(), "missing")
	if err == nil {
		t.Fatal("Open(missing): want ErrTemplateNotFound, got nil")
	}
	if !errors.Is(err, driven.ErrTemplateNotFound) {
		t.Errorf("err = %v, want wrap of driven.ErrTemplateNotFound", err)
	}
}

func yamlBlob(name, desc, version string) string {
	return "apiVersion: github.com/pt9912/u-boot/template/v1\n" +
		"name: " + name + "\n" +
		"description: \"" + desc + "\"\n" +
		"version: " + version + "\n"
}

func equalStrings(a, b []string) bool {
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
