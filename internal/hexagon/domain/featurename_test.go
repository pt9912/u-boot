package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

func TestNewFeatureName_ValidAccepts(t *testing.T) {
	t.Parallel()
	cases := []string{
		"git",
		"node",
		"java",
		"go",
		"cpp",
		"docker-cli",
		"kubectl-helm",
		"postgres-client",
		"a",                                // 1-char edge
		"feature7",                         // trailing digit
		"abcdefghij0123456789abcdefghij01", // 32-char edge
	}
	for _, in := range cases {
		in := in
		t.Run(in, func(t *testing.T) {
			t.Parallel()
			got, err := domain.NewFeatureName(in)
			if err != nil {
				t.Fatalf("NewFeatureName(%q): unexpected error: %v", in, err)
			}
			if got.String() != in {
				t.Errorf("NewFeatureName(%q).String() = %q, want %q", in, got.String(), in)
			}
		})
	}
}

func TestNewFeatureName_InvalidRejects(t *testing.T) {
	t.Parallel()
	cases := []string{
		"",                      // empty
		"Node",                  // uppercase
		"1node",                 // leading digit
		"-node",                 // leading dash
		"node-",                 // trailing dash
		"my_feature",            // underscore
		"my feature",            // space
		"foo/bar",               // slash (no path-syntax in names)
		strings.Repeat("a", 33), // length over cap
	}
	for _, in := range cases {
		in := in
		t.Run(in, func(t *testing.T) {
			t.Parallel()
			_, err := domain.NewFeatureName(in)
			if err == nil {
				t.Fatalf("NewFeatureName(%q): expected error, got nil", in)
			}
			if !errors.Is(err, domain.ErrInvalidFeatureName) {
				t.Errorf("NewFeatureName(%q): error %v does not wrap ErrInvalidFeatureName", in, err)
			}
		})
	}
}
