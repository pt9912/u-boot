package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

func TestNewServiceName_ValidAccepts(t *testing.T) {
	t.Parallel()
	cases := []string{
		"postgres",
		"keycloak",
		"otel",
		"a",                                  // 1-char edge
		"my-service-7",                       // dashes + digits
		"abcdefghij0123456789abcdefghij01",   // 32-char edge
	}
	for _, in := range cases {
		in := in
		t.Run(in, func(t *testing.T) {
			t.Parallel()
			got, err := domain.NewServiceName(in)
			if err != nil {
				t.Fatalf("NewServiceName(%q): unexpected error: %v", in, err)
			}
			if got.String() != in {
				t.Errorf("NewServiceName(%q).String() = %q, want %q", in, got.String(), in)
			}
		})
	}
}

func TestNewServiceName_InvalidRejects(t *testing.T) {
	t.Parallel()
	cases := []string{
		"",                                    // empty
		"Postgres",                            // uppercase
		"1postgres",                           // leading digit
		"-postgres",                           // leading dash
		"postgres-",                           // trailing dash
		"my_service",                          // underscore (not allowed in our spec; Compose tolerates it
		// but we want a single normalization rule for built-in services)
		"my service",                          // space
		strings.Repeat("a", 33),               // length over cap
	}
	for _, in := range cases {
		in := in
		t.Run(in, func(t *testing.T) {
			t.Parallel()
			_, err := domain.NewServiceName(in)
			if err == nil {
				t.Fatalf("NewServiceName(%q): expected error, got nil", in)
			}
			if !errors.Is(err, domain.ErrInvalidServiceName) {
				t.Errorf("NewServiceName(%q): error %v does not wrap ErrInvalidServiceName", in, err)
			}
		})
	}
}
