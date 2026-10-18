package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// LH-FA-DEV-004 / -006 / -007: closed value sets and UID range.
func TestSandboxSettings_Validation(t *testing.T) {
	t.Parallel()
	t.Run("profile", func(t *testing.T) {
		t.Parallel()
		for _, ok := range []string{"default", "sandbox"} {
			if _, err := domain.NewDevcontainerProfile(ok); err != nil {
				t.Errorf("profile %q: %v", ok, err)
			}
		}
		for _, bad := range []string{"", "Sandbox", "strict"} {
			if _, err := domain.NewDevcontainerProfile(bad); !errors.Is(err, domain.ErrInvalidSandboxSetting) {
				t.Errorf("profile %q: err = %v, want ErrInvalidSandboxSetting", bad, err)
			}
		}
	})
	t.Run("nestedRuntime", func(t *testing.T) {
		t.Parallel()
		for _, ok := range []string{"none", "podman", "docker"} {
			if _, err := domain.NewNestedRuntime(ok); err != nil {
				t.Errorf("nestedRuntime %q: %v", ok, err)
			}
		}
		for _, bad := range []string{"", "Docker", "PODMAN", "containerd"} {
			if _, err := domain.NewNestedRuntime(bad); !errors.Is(err, domain.ErrInvalidSandboxSetting) {
				t.Errorf("nestedRuntime %q: err = %v", bad, err)
			}
		}
	})
	t.Run("onUnavailable", func(t *testing.T) {
		t.Parallel()
		for _, ok := range []string{"warn", "fail"} {
			if _, err := domain.NewOnUnavailable(ok); err != nil {
				t.Errorf("onUnavailable %q: %v", ok, err)
			}
		}
		for _, bad := range []string{"", "ignore"} {
			if _, err := domain.NewOnUnavailable(bad); !errors.Is(err, domain.ErrInvalidSandboxSetting) {
				t.Errorf("onUnavailable %q: err = %v", bad, err)
			}
		}
	})
}

func TestContainerUID(t *testing.T) {
	t.Parallel()
	for _, ok := range []string{"1", "501", "1000", " 65535 "} {
		if _, err := domain.ParseContainerUID(ok); err != nil {
			t.Errorf("uid %q: %v", ok, err)
		}
	}
	for _, bad := range []string{"0", "-1", "65536", "abc", "", "1.5"} {
		if _, err := domain.ParseContainerUID(bad); !errors.Is(err, domain.ErrInvalidSandboxSetting) {
			t.Errorf("uid %q: err = %v, want ErrInvalidSandboxSetting", bad, err)
		}
	}
	if domain.DefaultContainerUID != 1000 {
		t.Errorf("DefaultContainerUID = %d, want 1000", domain.DefaultContainerUID)
	}
}

// LH-FA-DEV-008: an entry covers its subdomains, so a wildcard entry is
// rejected (exit code 10 via ErrInvalidSandboxSetting) with the name to
// enter instead; other malformed entries keep the generic text.
func TestNewEgressHost(t *testing.T) {
	t.Parallel()
	for _, ok := range []string{"example.org", " api.example.org ", "a-b.c1.example.org"} {
		if _, err := domain.NewEgressHost(ok); err != nil {
			t.Errorf("host %q: %v", ok, err)
		}
	}
	cases := []struct {
		raw     string
		want    []string
		wantNot []string
	}{
		{"*.example.org", []string{"wildcard", "already included", `"example.org"`}, []string{"not a lowercase"}},
		{"*example.org", []string{"wildcard", "already included", "without `*`"}, []string{"not a lowercase"}},
		{"*.*.example.org", []string{"wildcard", "without `*`"}, []string{"not a lowercase"}},
		{"api.*.example.org", []string{"wildcard", "without `*`"}, nil},
		{"Example.org", []string{"not a lowercase host name"}, []string{"already included"}},
		{"https://example.org", []string{"not a lowercase host name"}, []string{"already included"}},
		{"localhost", []string{"not a lowercase host name"}, []string{"already included"}},
	}
	for _, tc := range cases {
		_, err := domain.NewEgressHost(tc.raw)
		if !errors.Is(err, domain.ErrInvalidSandboxSetting) {
			t.Errorf("host %q: err = %v, want ErrInvalidSandboxSetting", tc.raw, err)
			continue
		}
		for _, w := range tc.want {
			if !strings.Contains(err.Error(), w) {
				t.Errorf("host %q: %q lacks %q", tc.raw, err.Error(), w)
			}
		}
		for _, w := range tc.wantNot {
			if strings.Contains(err.Error(), w) {
				t.Errorf("host %q: %q contains %q", tc.raw, err.Error(), w)
			}
		}
	}
}

// The length limit of 253 applies to the name behind `*.` as well: a
// 253-character name is offered back, a 254-character one is not. The
// error quotes the entry as typed, not trimmed.
func TestNewEgressHost_WildcardLengthAndRawQuote(t *testing.T) {
	t.Parallel()
	label := func(n int) string { return strings.Repeat("a", n) }
	name253 := strings.Join([]string{label(63), label(63), label(63), label(61)}, ".")
	name254 := strings.Join([]string{label(63), label(63), label(63), label(62)}, ".")
	if len(name253) != 253 || len(name254) != 254 {
		t.Fatalf("fixture lengths %d, %d", len(name253), len(name254))
	}
	_, err := domain.NewEgressHost("*." + name253)
	if err == nil || !strings.Contains(err.Error(), "enter "+`"`+name253+`"`+" instead") {
		t.Errorf("253-character name behind `*.`: err = %v, want the name offered", err)
	}
	_, err = domain.NewEgressHost("*." + name254)
	if err == nil || !strings.Contains(err.Error(), "wildcard") || strings.Contains(err.Error(), " instead") {
		t.Errorf("254-character name behind `*.`: err = %v, want the generic wildcard text", err)
	}
	_, err = domain.NewEgressHost("  *.example.org ")
	if err == nil || !strings.Contains(err.Error(), `entry "  *.example.org " is a wildcard`) ||
		!strings.Contains(err.Error(), `enter "example.org" instead`) {
		t.Errorf("untrimmed wildcard entry: err = %v, want raw quote and trimmed name", err)
	}
}
