package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

func TestNewProjectName_ValidAccepts(t *testing.T) {
	cases := []string{
		"a",
		"ab",
		"a1",
		"my-service",
		"a-b-c",
		"abc-123-def",
		strings.Repeat("a", 63),
		"a" + strings.Repeat("b", 62),
	}
	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			got, err := domain.NewProjectName(in)
			if err != nil {
				t.Fatalf("NewProjectName(%q): unexpected error: %v", in, err)
			}
			if got.String() != in {
				t.Fatalf("NewProjectName(%q).String() = %q, want %q", in, got.String(), in)
			}
		})
	}
}

func TestNewProjectName_InvalidRejects(t *testing.T) {
	cases := map[string]string{
		"empty":                  "",
		"starts with digit":      "1abc",
		"starts with dash":       "-abc",
		"ends with dash":         "abc-",
		"contains uppercase":     "Abc",
		"contains underscore":    "a_b",
		"contains space":         "a b",
		"contains slash":         "a/b",
		"too long (64)":          strings.Repeat("a", 64),
		"way too long":           strings.Repeat("a", 200),
		"only digits":            "123",
		"only dashes":            "---",
		"starts with dot":        ".abc",
		"contains unicode lower": "äbc",
	}
	for name, in := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := domain.NewProjectName(in)
			if err == nil {
				t.Fatalf("NewProjectName(%q): expected error, got nil", in)
			}
			if !errors.Is(err, domain.ErrInvalidProjectName) {
				t.Fatalf("NewProjectName(%q): error %v does not wrap ErrInvalidProjectName", in, err)
			}
		})
	}
}

func TestNormalizeProjectName(t *testing.T) {
	cases := map[string]string{
		"my-service":             "my-service",
		"MyService":              "myservice",
		"My_Service":             "my-service",
		"My Service":             "my-service",
		"---my--service---":      "my-service",
		"my.service/v1":          "my-service-v1",
		"  spaces  around  ":     "spaces-around",
		"a":                      "a",
		"---":                    "",
		"":                       "",
		"123":                    "123",
		strings.Repeat("a", 80):  strings.Repeat("a", 63),
		"a-" + strings.Repeat("b", 62) + "-c": "a-" + strings.Repeat("b", 61),
	}
	for in, want := range cases {
		t.Run(in, func(t *testing.T) {
			got := domain.NormalizeProjectName(in)
			if got != want {
				t.Fatalf("NormalizeProjectName(%q) = %q, want %q", in, got, want)
			}
		})
	}
}

func TestNormalizeProjectName_DoubleDashCollapse(t *testing.T) {
	// Why: explicit single-purpose test for LH-FA-INIT-002 rule 3
	// (collapse consecutive dashes to a single dash before trim).
	got := domain.NormalizeProjectName("foo----bar")
	if got != "foo-bar" {
		t.Fatalf("NormalizeProjectName(\"foo----bar\") = %q, want %q", got, "foo-bar")
	}
}

func TestProjectName_StringRoundTrip(t *testing.T) {
	name, err := domain.NewProjectName("hello-world")
	if err != nil {
		t.Fatalf("NewProjectName: %v", err)
	}
	if name.String() != "hello-world" {
		t.Fatalf("ProjectName.String() = %q, want %q", name.String(), "hello-world")
	}
}
