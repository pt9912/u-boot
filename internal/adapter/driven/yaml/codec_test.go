package yaml_test

import (
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driven/yaml"
)

type sample struct {
	Name    string `yaml:"name"`
	Version int    `yaml:"version"`
}

func TestCodec_MarshalUnmarshalRoundTrip(t *testing.T) {
	// Why: the round-trip is the actual semantic contract — Marshal +
	// Unmarshal must preserve the value. A separate test asserts the
	// presence of expected keys without pinning yaml.v3's formatting
	// (see TestCodec_MarshalEmitsExpectedKeys), so a future yaml.v3
	// bump that changes quoting/indent does not break the round-trip
	// test for the wrong reason.
	in := sample{Name: "demo", Version: 1}

	data, err := yaml.New().Marshal(in)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var back sample
	if err := yaml.New().Unmarshal(data, &back); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if back != in {
		t.Fatalf("Round-trip: got %+v, want %+v", back, in)
	}
}

func TestCodec_MarshalEmitsExpectedKeys(t *testing.T) {
	// Why: a smoke check that the YAML output mentions the struct's
	// `yaml:"…"` keys. Tolerant of yaml.v3 formatting (quoting,
	// indent, trailing newline) — only the key presence is pinned.
	in := sample{Name: "demo", Version: 1}

	data, err := yaml.New().Marshal(in)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	got := string(data)
	for _, key := range []string{"name:", "version:"} {
		if !strings.Contains(got, key) {
			t.Errorf("Marshal output missing key %q; got:\n%s", key, got)
		}
	}
}

func TestCodec_Unmarshal_InvalidYAMLReturnsError(t *testing.T) {
	var dst sample
	err := yaml.New().Unmarshal([]byte(":\n  not yaml"), &dst)
	if err == nil {
		t.Fatalf("Unmarshal(invalid): expected error, got nil")
	}
}
