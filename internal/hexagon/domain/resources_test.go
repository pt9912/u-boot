package domain_test

import (
	"errors"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// LH-FA-CONF-002: the character set of devcontainer.resources.memory.
func TestNewMemoryLimit(t *testing.T) {
	t.Parallel()
	for _, ok := range []string{"3g", "512m", "64k", "1024b", "1048576", "01g"} {
		if got, err := domain.NewMemoryLimit(ok); err != nil || string(got) != ok {
			t.Errorf("memory %q: got %q, err %v", ok, got, err)
		}
	}
	bad := []string{"", "0", "0g", "00m", "-1g", "1.5g", "3G", "3gb", "3 g", " 3g", "3g ", "4t", "g",
		`3g"`, `3\g`, "3g;x", "3g\n", "--memory=3g", "1g1"}
	for _, v := range bad {
		if _, err := domain.NewMemoryLimit(v); !errors.Is(err, domain.ErrInvalidResourceLimit) {
			t.Errorf("memory %q: err = %v, want ErrInvalidResourceLimit", v, err)
		}
	}
}

// LH-FA-CONF-002: the character set of devcontainer.resources.cpus.
func TestNewCPULimit(t *testing.T) {
	t.Parallel()
	for _, ok := range []string{"2", "1.5", "0.5", "0.01", "10", "02"} {
		if got, err := domain.NewCPULimit(ok); err != nil || string(got) != ok {
			t.Errorf("cpus %q: got %q, err %v", ok, got, err)
		}
	}
	bad := []string{"", "0", "0.0", "-1", ".5", "2.", "1,5", "1e1", "2x", " 2", "2 ", "1.2.3",
		`2"`, `2\`, "2;x", "2\n", "--cpus=2"}
	for _, v := range bad {
		if _, err := domain.NewCPULimit(v); !errors.Is(err, domain.ErrInvalidResourceLimit) {
			t.Errorf("cpus %q: err = %v, want ErrInvalidResourceLimit", v, err)
		}
	}
}
