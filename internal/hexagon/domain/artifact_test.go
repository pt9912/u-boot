package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

func TestNewArtifact_KnownValues(t *testing.T) {
	t.Parallel()
	cases := map[string]domain.Artifact{
		"changelog":    domain.ArtifactChangelog,
		"readme":       domain.ArtifactReadme,
		"env-example":  domain.ArtifactEnvExample,
		"devcontainer": domain.ArtifactDevcontainer,
	}
	for raw, want := range cases {
		got, err := domain.NewArtifact(raw)
		if err != nil {
			t.Errorf("NewArtifact(%q): unexpected error %v", raw, err)
			continue
		}
		if got != want {
			t.Errorf("NewArtifact(%q) = %v, want %v", raw, got, want)
		}
	}
}

func TestNewArtifact_Unknown_ReturnsErrInvalidArtifact(t *testing.T) {
	t.Parallel()
	for _, raw := range []string{"", "Changelog", "unknown", "dockerfile", "compose"} {
		_, err := domain.NewArtifact(raw)
		if err == nil {
			t.Errorf("NewArtifact(%q): expected error, got nil", raw)
			continue
		}
		if !errors.Is(err, domain.ErrInvalidArtifact) {
			t.Errorf("NewArtifact(%q): error %v does not wrap ErrInvalidArtifact", raw, err)
		}
	}
}

func TestNewArtifact_UnknownError_MentionsCatalogue(t *testing.T) {
	t.Parallel()
	_, err := domain.NewArtifact("dockerfile")
	if err == nil {
		t.Fatal("expected error")
	}
	msg := err.Error()
	for _, want := range []string{"changelog", "readme", "env-example", "devcontainer"} {
		if !strings.Contains(msg, want) {
			t.Errorf("error message %q does not contain catalogue entry %q", msg, want)
		}
	}
}

func TestArtifact_String_RoundTrip(t *testing.T) {
	t.Parallel()
	for _, a := range []domain.Artifact{
		domain.ArtifactChangelog,
		domain.ArtifactReadme,
		domain.ArtifactEnvExample,
		domain.ArtifactDevcontainer,
	} {
		s := a.String()
		got, err := domain.NewArtifact(s)
		if err != nil {
			t.Errorf("NewArtifact(%v.String()=%q): unexpected error %v", a, s, err)
			continue
		}
		if got != a {
			t.Errorf("round-trip: NewArtifact(%v.String()=%q) = %v, want %v", a, s, got, a)
		}
	}
}

func TestArtifact_String_Unknown(t *testing.T) {
	t.Parallel()
	// Out-of-range value (defensive default branch). The CLI never
	// constructs Artifact values directly; this only protects against
	// a misuse like Artifact(99) sneaking into a log line.
	if got := domain.Artifact(99).String(); got != "unknown" {
		t.Errorf("Artifact(99).String() = %q, want %q", got, "unknown")
	}
}
