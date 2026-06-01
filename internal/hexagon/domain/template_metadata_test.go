package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

func TestTemplateMetadata_Validate(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		meta      domain.TemplateMetadata
		wantErr   bool
		wantProbs []string
	}{
		{
			name: "valid minimum",
			meta: domain.TemplateMetadata{
				Name:        "basic",
				Description: "minimal skeleton",
				Version:     "0.1.0",
			},
		},
		{
			name: "valid hyphenated name",
			meta: domain.TemplateMetadata{
				Name:        "micronaut-sveltekit",
				Description: "combined stack",
				Version:     "1.0.0",
			},
		},
		{
			name:      "missing name",
			meta:      domain.TemplateMetadata{Description: "x", Version: "1"},
			wantErr:   true,
			wantProbs: []string{"name is required"},
		},
		{
			name:      "whitespace-only name",
			meta:      domain.TemplateMetadata{Name: "   ", Description: "x", Version: "1"},
			wantErr:   true,
			wantProbs: []string{"name is required"},
		},
		{
			name:      "uppercase name rejected",
			meta:      domain.TemplateMetadata{Name: "Basic", Description: "x", Version: "1"},
			wantErr:   true,
			wantProbs: []string{`name "Basic" must be kebab-case`},
		},
		{
			name:      "trailing dash rejected",
			meta:      domain.TemplateMetadata{Name: "basic-", Description: "x", Version: "1"},
			wantErr:   true,
			wantProbs: []string{`name "basic-" must be kebab-case`},
		},
		{
			name:      "consecutive dashes rejected (review-followup N2)",
			meta:      domain.TemplateMetadata{Name: "my--bad", Description: "x", Version: "1"},
			wantErr:   true,
			wantProbs: []string{`name "my--bad" must be kebab-case`},
		},
		{
			name:      "leading dash rejected",
			meta:      domain.TemplateMetadata{Name: "-foo", Description: "x", Version: "1"},
			wantErr:   true,
			wantProbs: []string{`name "-foo" must be kebab-case`},
		},
		{
			name:      "missing description",
			meta:      domain.TemplateMetadata{Name: "basic", Version: "1"},
			wantErr:   true,
			wantProbs: []string{"description is required"},
		},
		{
			name:      "missing version",
			meta:      domain.TemplateMetadata{Name: "basic", Description: "x"},
			wantErr:   true,
			wantProbs: []string{"version is required"},
		},
		{
			name:    "all fields missing — every problem reported",
			meta:    domain.TemplateMetadata{},
			wantErr: true,
			wantProbs: []string{
				"name is required",
				"description is required",
				"version is required",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.meta.Validate()
			if !tc.wantErr {
				if err != nil {
					t.Fatalf("Validate() = %v, want nil", err)
				}
				return
			}
			if err == nil {
				t.Fatal("Validate() = nil, want error")
			}
			if !errors.Is(err, domain.ErrInvalidTemplate) {
				t.Errorf("err = %v, want wrap of domain.ErrInvalidTemplate", err)
			}
			msg := err.Error()
			for _, want := range tc.wantProbs {
				if !strings.Contains(msg, want) {
					t.Errorf("err.Error() = %q, missing fragment %q", msg, want)
				}
			}
		})
	}
}

func TestTemplateMetadata_ZeroValueIsInvalid(t *testing.T) {
	t.Parallel()
	var meta domain.TemplateMetadata
	if err := meta.Validate(); err == nil {
		t.Error("zero TemplateMetadata.Validate() = nil, want error")
	}
}
