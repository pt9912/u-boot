package domain_test

import (
	"slices"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

func TestConfigFileCandidates_HoldsOnlyConfigFileName(t *testing.T) {
	t.Parallel()
	got := domain.ConfigFileCandidates()
	want := []string{domain.ConfigFileName}
	if !slices.Equal(got, want) {
		t.Errorf("ConfigFileCandidates() = %v, want %v", got, want)
	}
	if domain.ConfigFileName != "u-boot.yaml" {
		t.Errorf("ConfigFileName = %q, want %q", domain.ConfigFileName, "u-boot.yaml")
	}
}

func TestConfigFileCandidates_ReturnsFreshSlicePerCall(t *testing.T) {
	t.Parallel()
	first := domain.ConfigFileCandidates()
	first[0] = "changed.yaml"
	_ = append(first[:0], "a", "b")

	second := domain.ConfigFileCandidates()
	if !slices.Equal(second, []string{domain.ConfigFileName}) {
		t.Errorf("second call = %v after modifying the first result, want [%s]", second, domain.ConfigFileName)
	}
}
