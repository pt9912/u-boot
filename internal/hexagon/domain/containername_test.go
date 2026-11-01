package domain_test

import (
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// LH-FA-DEV-010: the container name of the started devcontainer joins the
// prefix, the project part and, with an instance, the instance part.
func TestComposeDevcontainerContainerName(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		project  string
		instance string
		want     string
	}{
		{"without instance", "demo", "", "u-boot-demo"},
		{"with instance", "demo", "agent-1", "u-boot-demo-agent-1"},
		{"project of digits", "9bench", "a", "u-boot-9bench-a"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := domain.ComposeDevcontainerContainerName(tc.project, tc.instance)
			if err != nil || got != tc.want {
				t.Errorf("Compose(%q, %q) = %q, %v; want %q, nil", tc.project, tc.instance, got, err, tc.want)
			}
		})
	}
	for _, instance := range []string{"", "agent-1"} {
		t.Run("empty project is not composable with instance "+instance, func(t *testing.T) {
			t.Parallel()
			if _, err := domain.ComposeDevcontainerContainerName("", instance); err == nil {
				t.Errorf("Compose(\"\", %q) composed", instance)
			}
		})
	}
}
