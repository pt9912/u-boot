package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/adapter/driving/cli/jsontestutil"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// TestTemplateListJSON_EnvelopeWithData pins the slice-v1-cli-json-
// dry-run-template T2 migration: `template list --json` emits the
// Minimalkontrakt-Envelope (command="template", subcommand="list")
// with the `[]templateJSON` projection in `data`.
func TestTemplateListJSON_EnvelopeWithData(t *testing.T) {
	uc := &fakeTemplateListUseCase{resp: driving.TemplateListResponse{
		Templates: []domain.TemplateMetadata{
			{Name: "basic", Description: "Minimal starter", Version: "1.0.0"},
		},
	}}
	app := newAppWithTemplateList(uc)

	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--json", "template", "list"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("template"),
		jsontestutil.WithSubcommand("list"),
		jsontestutil.WithExitCode(0),
	)
	var env map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	data, ok := env["data"].([]any)
	if !ok || len(data) != 1 {
		t.Fatalf("data must be a 1-element array, got %v", env["data"])
	}
	first, _ := data[0].(map[string]any)
	if first["name"] != "basic" {
		t.Errorf("data[0].name = %v, want basic", first["name"])
	}
}

// TestTemplateListJSON_EmptyCatalogDataArray pins the empty-catalog
// path (R2-LOW-2): zero templates → `data: []` (not null).
func TestTemplateListJSON_EmptyCatalogDataArray(t *testing.T) {
	app := newAppWithTemplateList(&fakeTemplateListUseCase{})

	var stdout bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--json", "template", "list"}, &stdout, &bytes.Buffer{}); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(stdout.String(), "\"data\":[]") {
		t.Errorf("empty catalog must emit data:[] (not null); got %s", stdout.String())
	}
}

// TestTemplateListJSON_ErrorEnvelope pins the T0-(f) error path:
// a catalog failure (ErrTemplateCatalog) surfaces as an error
// envelope with command/subcommand set and exit 14 (technical-
// persistence class).
func TestTemplateListJSON_ErrorEnvelope(t *testing.T) {
	uc := &fakeTemplateListUseCase{err: driving.ErrTemplateCatalog}
	app := newAppWithTemplateList(uc)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "template", "list"}, &stdout, &stderr)
	if cli.ExitCode(err) != 14 {
		t.Fatalf("exit = %d, want 14 (ErrTemplateCatalog → FS class)", cli.ExitCode(err))
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("template"),
		jsontestutil.WithSubcommand("list"),
		jsontestutil.WithExitCode(14),
		jsontestutil.WithExpectedCodes("LH-NFA-REL-003"),
	)
}

// TestTemplateListText_StillWorks pins the non-JSON path stays
// intact: the human tabular form is unchanged by the T2 migration.
func TestTemplateListText_StillWorks(t *testing.T) {
	uc := &fakeTemplateListUseCase{resp: driving.TemplateListResponse{
		Templates: []domain.TemplateMetadata{{Name: "basic", Description: "Minimal", Version: "1.0.0"}},
	}}
	app := newAppWithTemplateList(uc)

	var stdout bytes.Buffer
	if err := app.Execute(context.Background(), []string{"template", "list"}, &stdout, &bytes.Buffer{}); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(stdout.String(), "NAME") || !strings.Contains(stdout.String(), "basic") {
		t.Errorf("text form should render header + row; got %s", stdout.String())
	}
}
