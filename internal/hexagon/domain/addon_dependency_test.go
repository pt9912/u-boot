package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

func TestAddOnDependency_Validate(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		dep       domain.AddOnDependency
		wantErr   bool
		wantProbs []string
	}{
		{
			name: "valid keycloak-style declaration",
			dep: domain.AddOnDependency{
				Requires:    mustServiceNameDomainTest(t, "postgres"),
				WhenPath:    "services.keycloak.persistence",
				EqualsValue: "external-postgres",
			},
		},
		{
			name: "missing Requires (zero-value ServiceName)",
			dep: domain.AddOnDependency{
				WhenPath:    "services.x.foo",
				EqualsValue: "bar",
			},
			wantErr:   true,
			wantProbs: []string{"Requires is required"},
		},
		{
			name: "missing WhenPath",
			dep: domain.AddOnDependency{
				Requires:    mustServiceNameDomainTest(t, "postgres"),
				EqualsValue: "bar",
			},
			wantErr:   true,
			wantProbs: []string{"WhenPath is required"},
		},
		{
			name: "missing EqualsValue",
			dep: domain.AddOnDependency{
				Requires: mustServiceNameDomainTest(t, "postgres"),
				WhenPath: "services.x.foo",
			},
			wantErr:   true,
			wantProbs: []string{"EqualsValue is required"},
		},
		{
			name: "whitespace-only WhenPath rejected",
			dep: domain.AddOnDependency{
				Requires:    mustServiceNameDomainTest(t, "postgres"),
				WhenPath:    "   ",
				EqualsValue: "bar",
			},
			wantErr:   true,
			wantProbs: []string{"WhenPath is required"},
		},
		{
			name: "whitespace-only EqualsValue rejected",
			dep: domain.AddOnDependency{
				Requires:    mustServiceNameDomainTest(t, "postgres"),
				WhenPath:    "services.x.foo",
				EqualsValue: "  \t  ",
			},
			wantErr:   true,
			wantProbs: []string{"EqualsValue is required"},
		},
		{
			name:    "all fields missing — every problem reported",
			dep:     domain.AddOnDependency{},
			wantErr: true,
			wantProbs: []string{
				"Requires is required",
				"WhenPath is required",
				"EqualsValue is required",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.dep.Validate()
			if !tc.wantErr {
				if err != nil {
					t.Fatalf("Validate() = %v, want nil", err)
				}
				return
			}
			if err == nil {
				t.Fatal("Validate() = nil, want error")
			}
			if !errors.Is(err, domain.ErrInvalidAddOnDependency) {
				t.Errorf("err = %v, want wrap of domain.ErrInvalidAddOnDependency", err)
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

func TestAddOnDependency_ZeroValueIsInvalid(t *testing.T) {
	t.Parallel()
	var dep domain.AddOnDependency
	if err := dep.Validate(); err == nil {
		t.Error("zero AddOnDependency.Validate() = nil, want error")
	}
}

func mustServiceNameDomainTest(t *testing.T, raw string) domain.ServiceName {
	t.Helper()
	name, err := domain.NewServiceName(raw)
	if err != nil {
		t.Fatalf("NewServiceName(%q): %v", raw, err)
	}
	return name
}
