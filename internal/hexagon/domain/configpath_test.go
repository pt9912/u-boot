package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

func TestNewConfigPath_KnownPaths(t *testing.T) {
	t.Parallel()
	cases := []struct {
		raw          string
		wantKind     domain.ConfigPathKind
		wantWrite    bool
		wantService  string // empty unless kind == ConfigServiceEnabled
	}{
		{"project.name", domain.ConfigProjectName, true, ""},
		{"devcontainer.enabled", domain.ConfigDevcontainerEnabled, true, ""},
		{"services.postgres.enabled", domain.ConfigServiceEnabled, false, "postgres"},
		{"services.keycloak.enabled", domain.ConfigServiceEnabled, false, "keycloak"},
		{"services.otel.enabled", domain.ConfigServiceEnabled, false, "otel"},
	}
	for _, tc := range cases {
		got, err := domain.NewConfigPath(tc.raw)
		if err != nil {
			t.Errorf("NewConfigPath(%q): unexpected error %v", tc.raw, err)
			continue
		}
		if got.Kind != tc.wantKind {
			t.Errorf("NewConfigPath(%q).Kind = %v, want %v", tc.raw, got.Kind, tc.wantKind)
		}
		if got.WriteAllowed != tc.wantWrite {
			t.Errorf("NewConfigPath(%q).WriteAllowed = %v, want %v", tc.raw, got.WriteAllowed, tc.wantWrite)
		}
		if got.Service.String() != tc.wantService {
			t.Errorf("NewConfigPath(%q).Service = %q, want %q", tc.raw, got.Service.String(), tc.wantService)
		}
	}
}

func TestNewConfigPath_UnknownPath_ReturnsErrInvalidConfigPath(t *testing.T) {
	t.Parallel()
	for _, raw := range []string{
		"",
		"unknown.path",
		"project",                       // missing leaf segment
		"project.foo",                   // wrong leaf
		"project.name.extra",            // trailing segment
		"services",                      // top-level only
		"services.postgres",             // missing leaf
		"services.postgres.persistence", // V1 field, not in MVP whitelist
		"devcontainer",                  // top-level only
		"devcontainer.foo",              // wrong leaf
		"schemaVersion",                 // read-only field, not whitelisted
		"foo.bar.baz",
	} {
		_, err := domain.NewConfigPath(raw)
		if err == nil {
			t.Errorf("NewConfigPath(%q): expected error, got nil", raw)
			continue
		}
		if !errors.Is(err, domain.ErrInvalidConfigPath) {
			t.Errorf("NewConfigPath(%q): err %v does not wrap ErrInvalidConfigPath", raw, err)
		}
	}
}

// TestNewConfigPath_InvalidServiceName_WrapsBothSentinels pins
// that a malformed `services.<svc>.enabled` path surfaces both
// [domain.ErrInvalidConfigPath] AND [domain.ErrInvalidServiceName]
// — the application layer can branch on either, and the user-
// visible message names the offending service name explicitly.
func TestNewConfigPath_InvalidServiceName_WrapsBothSentinels(t *testing.T) {
	t.Parallel()
	cases := []string{
		"services..enabled",          // empty service name
		"services.UPPERCASE.enabled", // wrong case (LH-FA-INIT-006 regex)
		"services.with space.enabled",
		"services.-leading-dash.enabled",
	}
	for _, raw := range cases {
		_, err := domain.NewConfigPath(raw)
		if err == nil {
			t.Errorf("NewConfigPath(%q): expected error, got nil", raw)
			continue
		}
		if !errors.Is(err, domain.ErrInvalidConfigPath) {
			t.Errorf("NewConfigPath(%q): err %v does not wrap ErrInvalidConfigPath", raw, err)
		}
		if !errors.Is(err, domain.ErrInvalidServiceName) {
			t.Errorf("NewConfigPath(%q): err %v does not chain ErrInvalidServiceName", raw, err)
		}
	}
}

func TestNewConfigPath_ErrorMessage_NamesAllowedPaths(t *testing.T) {
	t.Parallel()
	// The CLI surfaces the error message verbatim; users need to
	// see the catalogue to know how to fix their typo. Pin that
	// every whitelist entry appears in the unknown-path message.
	_, err := domain.NewConfigPath("totally.unknown.path")
	if err == nil {
		t.Fatal("expected error")
	}
	msg := err.Error()
	for _, want := range []string{
		"project.name",
		"devcontainer.enabled",
		"services.<svc>.enabled",
	} {
		if !strings.Contains(msg, want) {
			t.Errorf("error message %q does not mention whitelist entry %q", msg, want)
		}
	}
}

func TestConfigPath_String_RoundTrip(t *testing.T) {
	t.Parallel()
	for _, raw := range []string{
		"project.name",
		"devcontainer.enabled",
		"services.postgres.enabled",
		"services.keycloak.enabled",
	} {
		p, err := domain.NewConfigPath(raw)
		if err != nil {
			t.Errorf("NewConfigPath(%q): %v", raw, err)
			continue
		}
		if got := p.String(); got != raw {
			t.Errorf("String() round-trip: got %q, want %q", got, raw)
		}
		// Equality round-trip: re-parsing produces an equal value.
		p2, err := domain.NewConfigPath(p.String())
		if err != nil {
			t.Errorf("re-parse %q: %v", p.String(), err)
			continue
		}
		if p != p2 {
			t.Errorf("round-trip inequality: %+v vs %+v", p, p2)
		}
	}
}

func TestConfigPath_String_DefensiveBranch(t *testing.T) {
	t.Parallel()
	// Out-of-range Kind value renders as a deterministic
	// `ConfigPath(kind=N)` form so debug logs surface the actual
	// integer instead of an opaque empty string. Mirrors the
	// review-followup S4 pattern from domain.Artifact.
	p := domain.ConfigPath{Kind: domain.ConfigPathKind(99)}
	if got := p.String(); !strings.Contains(got, "99") {
		t.Errorf("ConfigPath{Kind: 99}.String() = %q, want substring %q", got, "99")
	}
}
