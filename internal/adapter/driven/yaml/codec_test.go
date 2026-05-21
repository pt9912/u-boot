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

func TestCodec_RoundTrip(t *testing.T) {
	in := sample{Name: "demo", Version: 1}

	data, err := yaml.New().Marshal(in)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	got := strings.TrimSpace(string(data))
	want := "name: demo\nversion: 1"
	if got != want {
		t.Fatalf("Marshal output:\n got: %q\nwant: %q", got, want)
	}

	var back sample
	if err := yaml.New().Unmarshal(data, &back); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if back != in {
		t.Fatalf("Round-trip: got %+v, want %+v", back, in)
	}
}

func TestCodec_Unmarshal_InvalidYAMLReturnsError(t *testing.T) {
	var dst sample
	err := yaml.New().Unmarshal([]byte(":\n  not yaml"), &dst)
	if err == nil {
		t.Fatalf("Unmarshal(invalid): expected error, got nil")
	}
}
