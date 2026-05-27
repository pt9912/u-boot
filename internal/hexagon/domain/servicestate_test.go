package domain_test

import (
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

func TestServiceState_String(t *testing.T) {
	t.Parallel()
	cases := []struct {
		state domain.ServiceState
		want  string
	}{
		{domain.ServiceStateUnregistered, "unregistered"},
		{domain.ServiceStateActive, "active"},
		{domain.ServiceStateDeactivated, "deactivated"},
		{domain.ServiceStateEnabledUnset, "enabled-unset"},
		{domain.ServiceStateInconsistentYAML, "inconsistent-yaml"},
		{domain.ServiceStateInconsistentBlock, "inconsistent-block"},
		{domain.ServiceState(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.state.String(); got != tc.want {
			t.Errorf("ServiceState(%d).String() = %q, want %q", tc.state, got, tc.want)
		}
	}
}
