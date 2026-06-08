package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

func TestTemplateList_HumanReadable(t *testing.T) {
	uc := &fakeTemplateListUseCase{
		resp: driving.TemplateListResponse{
			Templates: []domain.TemplateMetadata{
				{Name: "alpha", Description: "First template.", Version: "1.0.0"},
				{Name: "basic", Description: "Minimal skeleton.", Version: "0.1.0"},
			},
		},
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithTemplateList(uc).Execute(context.Background(), []string{"template", "list"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Execute: %v (stderr=%q)", err, stderr.String())
	}
	if !uc.called {
		t.Error("template list use case was not called")
	}

	out := stdout.String()
	// Header line uppercased so it stands out from data rows.
	if !strings.Contains(out, "NAME") || !strings.Contains(out, "DESCRIPTION") || !strings.Contains(out, "VERSION") {
		t.Errorf("stdout missing column headers; got:\n%s", out)
	}
	for _, want := range []string{"alpha", "First template.", "1.0.0", "basic", "Minimal skeleton.", "0.1.0"} {
		if !strings.Contains(out, want) {
			t.Errorf("stdout missing %q; got:\n%s", want, out)
		}
	}
	// Tabwriter padding: between two columns there must be at least
	// two spaces (the minwidth=0, tabwidth=0, padding=2 setup).
	// Cheap regression pin: the alpha row must not run name and
	// description together.
	if strings.Contains(out, "alphaFirst") {
		t.Errorf("name and description appear glued together; got:\n%s", out)
	}
}

func TestTemplateList_EmptyCatalogRendersMessage(t *testing.T) {
	uc := &fakeTemplateListUseCase{resp: driving.TemplateListResponse{Templates: nil}}
	var stdout, stderr bytes.Buffer
	if err := newAppWithTemplateList(uc).Execute(context.Background(), []string{"template", "list"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute: %v (stderr=%q)", err, stderr.String())
	}
	if got := strings.TrimSpace(stdout.String()); got != "No templates available." {
		t.Errorf("stdout = %q, want %q", got, "No templates available.")
	}
}

func TestTemplateList_JSON(t *testing.T) {
	uc := &fakeTemplateListUseCase{
		resp: driving.TemplateListResponse{
			Templates: []domain.TemplateMetadata{
				{
					Name:            "basic",
					Description:     "Minimal skeleton.",
					Version:         "0.1.0",
					SupportedAddOns: []string{"postgres"},
					GeneratedFiles:  []string{"u-boot.yaml", "compose.yaml"},
					RequiredTools:   nil, // pin nil-→-empty-array normalisation
					Variables: []domain.TemplateVariable{
						{Name: "groupId", Description: "Maven group", Default: "com.example", Required: true},
					},
				},
			},
		},
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithTemplateList(uc).Execute(context.Background(), []string{"template", "list", "--json"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Execute: %v (stderr=%q)", err, stderr.String())
	}

	// slice-v1-cli-json-dry-run-template T2: output migrated from a
	// raw array to the Minimalkontrakt-Envelope; the projection now
	// rides in `data`.
	var env struct {
		Command    string           `json:"command"`
		Subcommand string           `json:"subcommand"`
		Data       []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &env); err != nil {
		t.Fatalf("json.Unmarshal: %v; output was:\n%s", err, stdout.String())
	}
	if env.Command != "template" || env.Subcommand != "list" {
		t.Errorf("envelope command/subcommand = %q/%q, want template/list", env.Command, env.Subcommand)
	}
	got := env.Data
	if len(got) != 1 {
		t.Fatalf("len(data) = %d, want 1", len(got))
	}
	want := map[string]any{
		"name":            "basic",
		"description":     "Minimal skeleton.",
		"version":         "0.1.0",
		"supportedAddOns": []any{"postgres"},
		"generatedFiles":  []any{"u-boot.yaml", "compose.yaml"},
		"requiredTools":   []any{}, // nil normalised to empty array
		"variables": []any{
			map[string]any{
				"name":        "groupId",
				"description": "Maven group",
				"default":     "com.example",
				"required":    true,
			},
		},
	}
	for k, v := range want {
		if !deepEqualJSON(got[0][k], v) {
			t.Errorf("got[%q] = %#v, want %#v", k, got[0][k], v)
		}
	}
}

func TestTemplateList_EmptyCatalog_JSONIsEmptyArray(t *testing.T) {
	uc := &fakeTemplateListUseCase{resp: driving.TemplateListResponse{Templates: nil}}
	var stdout, stderr bytes.Buffer
	if err := newAppWithTemplateList(uc).Execute(context.Background(), []string{"template", "list", "--json"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute: %v (stderr=%q)", err, stderr.String())
	}
	// T2: empty catalog → envelope with `data: []` (not null).
	if !strings.Contains(stdout.String(), "\"data\":[]") {
		t.Errorf("empty catalog must emit data:[] (not null); got %s", stdout.String())
	}
	var env struct {
		Data []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &env); err != nil {
		t.Fatalf("json.Unmarshal: %v; output was:\n%s", err, stdout.String())
	}
	if len(env.Data) != 0 {
		t.Errorf("len(data) = %d, want 0 (empty catalog → []); got = %v", len(env.Data), env.Data)
	}
}

func TestTemplateList_UseCaseError_PropagatesAsCode14(t *testing.T) {
	uc := &fakeTemplateListUseCase{
		err: errors.New("wrapped: " + driving.ErrTemplateCatalog.Error()),
	}
	// Wrap manually with %w so the sentinel chain is preserved.
	uc.err = wrapTemplateCatalog(uc.err)

	var stdout, stderr bytes.Buffer
	err := newAppWithTemplateList(uc).Execute(context.Background(), []string{"template", "list"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("Execute: want error, got nil")
	}
	if code := cli.ExitCode(err); code != 14 {
		t.Errorf("ExitCode = %d, want 14 (LH-FA-CLI-006 technical/persistence)", code)
	}
}

func TestTemplateList_HelpListsListSubcommand(t *testing.T) {
	// `u-boot template --help` should mention the `list` subcommand
	// so users can discover it without grep'ing the docs.
	var stdout, stderr bytes.Buffer
	uc := &fakeTemplateListUseCase{}
	if err := newAppWithTemplateList(uc).Execute(context.Background(), []string{"template", "--help"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute: %v (stderr=%q)", err, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "list") {
		t.Errorf("template --help missing `list` mention; got:\n%s", out)
	}
}

// wrapTemplateCatalog mimics the application service's wrap so the
// CLI test does not need to import the application package directly.
func wrapTemplateCatalog(cause error) error {
	return wrapErr(driving.ErrTemplateCatalog, cause)
}

func wrapErr(sentinel, cause error) error {
	if cause == nil {
		return sentinel
	}
	return joinErrs(sentinel, cause)
}

// joinErrs avoids fmt.Errorf in tests so we can keep the cli_test
// import set tight. errors.Join (Go 1.20) wraps both arguments
// so errors.Is hits either side.
func joinErrs(a, b error) error {
	return errors.Join(a, b)
}

func deepEqualJSON(a, b any) bool {
	ab, err := json.Marshal(a)
	if err != nil {
		return false
	}
	bb, err := json.Marshal(b)
	if err != nil {
		return false
	}
	return string(ab) == string(bb)
}
