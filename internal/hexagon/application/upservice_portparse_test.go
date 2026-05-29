package application_test

import (
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
)

// portparse tests use the export_test.go bridge
// `ParseComposePortForTest` so they live in the canonical
// `_test`-package layout while exercising the unexported
// `parseComposePort` helper directly. The eight syntax cases below
// are the primary specification of the helper (M6 slice plan
// §T4-portparse-Tabelle).

func TestParseComposePort_NakedInteger(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		raw  any
		want application.PortProbeTargetForTest
	}{
		{"int-5432", 5432, application.PortProbeTargetForTest{Host: "localhost", Port: 5432}},
		{"int64-5432", int64(5432), application.PortProbeTargetForTest{Host: "localhost", Port: 5432}},
		{"float64-5432", float64(5432), application.PortProbeTargetForTest{Host: "localhost", Port: 5432}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, probable := application.ParseComposePortForTest(tc.raw)
			if !probable {
				t.Fatalf("expected probable=true for %v, got false", tc.raw)
			}
			if got != tc.want {
				t.Errorf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestParseComposePort_StringShortForm(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		raw  string
		want application.PortProbeTargetForTest
	}{
		{"single-port", "5432", application.PortProbeTargetForTest{Host: "localhost", Port: 5432}},
		{"host-container", "5432:5432", application.PortProbeTargetForTest{Host: "localhost", Port: 5432}},
		{"different-host-container", "9091:80", application.PortProbeTargetForTest{Host: "localhost", Port: 9091}},
		{"with-host-ip", "127.0.0.1:5432:5432", application.PortProbeTargetForTest{Host: "127.0.0.1", Port: 5432}},
		{"tcp-explicit", "5432:5432/tcp", application.PortProbeTargetForTest{Host: "localhost", Port: 5432}},
		{"tcp-explicit-upper", "5432:5432/TCP", application.PortProbeTargetForTest{Host: "localhost", Port: 5432}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, probable := application.ParseComposePortForTest(tc.raw)
			if !probable {
				t.Fatalf("expected probable=true for %q, got false", tc.raw)
			}
			if got != tc.want {
				t.Errorf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestParseComposePort_NonProbableStringForms(t *testing.T) {
	t.Parallel()
	// All forms in this table must return probable=false. The
	// caller emits a Severity-warn diagnostic and skips probing —
	// LH-FA-UP-001 §969 mandates the graceful warn path.
	cases := []struct {
		name string
		raw  string
	}{
		{"udp-suffix", "5432:5432/udp"},
		{"sctp-suffix", "5432:5432/sctp"},
		{"udp-uppercase", "5432:5432/UDP"},
		{"range-host-container", "5000-5010:5000-5010"},
		{"range-single-side", "5000-5010:5000"},
		{"empty", ""},
		{"non-numeric-port", "abc"},
		{"non-numeric-host-port", "abc:5432"},
		{"too-many-colons", "127.0.0.1:5432:5432:extra"},
		{"port-out-of-range-high", "70000"},
		{"port-zero", "0"},
		{"port-negative", "-1"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, probable := application.ParseComposePortForTest(tc.raw)
			if probable {
				t.Errorf("%q: expected probable=false (warn-diagnostic path), got true", tc.raw)
			}
		})
	}
}

func TestParseComposePort_LongSyntax(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		raw  map[string]any
		want application.PortProbeTargetForTest
	}{
		{
			"minimal",
			map[string]any{"target": 5432, "published": 5432},
			application.PortProbeTargetForTest{Host: "localhost", Port: 5432},
		},
		{
			"explicit-tcp",
			map[string]any{"target": 5432, "published": 5432, "protocol": "tcp"},
			application.PortProbeTargetForTest{Host: "localhost", Port: 5432},
		},
		{
			"with-host-ip",
			map[string]any{"target": 5432, "published": 5432, "host_ip": "127.0.0.1"},
			application.PortProbeTargetForTest{Host: "127.0.0.1", Port: 5432},
		},
		{
			"published-as-string",
			map[string]any{"target": 5432, "published": "5432"},
			application.PortProbeTargetForTest{Host: "localhost", Port: 5432},
		},
		{
			"published-as-int64",
			map[string]any{"target": int64(5432), "published": int64(5432)},
			application.PortProbeTargetForTest{Host: "localhost", Port: 5432},
		},
		{
			"published-as-float64",
			map[string]any{"target": float64(5432), "published": float64(5432)},
			application.PortProbeTargetForTest{Host: "localhost", Port: 5432},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, probable := application.ParseComposePortForTest(tc.raw)
			if !probable {
				t.Fatalf("expected probable=true, got false (raw=%v)", tc.raw)
			}
			if got != tc.want {
				t.Errorf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestParseComposePort_LongSyntaxNonProbable(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		raw  map[string]any
	}{
		{"udp", map[string]any{"published": 5432, "protocol": "udp"}},
		{"sctp", map[string]any{"published": 5432, "protocol": "sctp"}},
		{"protocol-non-string", map[string]any{"published": 5432, "protocol": 42}},
		{"missing-published", map[string]any{"target": 5432}},
		{"published-range-string", map[string]any{"published": "5000-5010"}},
		{"published-non-numeric-string", map[string]any{"published": "abc"}},
		{"published-fractional-float", map[string]any{"published": 5432.5}},
		{"published-non-numeric-type", map[string]any{"published": true}},
		{"published-zero", map[string]any{"published": 0}},
		{"published-too-high", map[string]any{"published": 70000}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, probable := application.ParseComposePortForTest(tc.raw)
			if probable {
				t.Errorf("%v: expected probable=false, got true", tc.raw)
			}
		})
	}
}

func TestParseComposePort_HostNormalization(t *testing.T) {
	t.Parallel()
	// Bind-to-all sentinels ("0.0.0.0", "::", "") are normalized
	// to "localhost" so the NetProbe doesn't dial a non-routable
	// address.
	cases := []struct {
		name string
		raw  map[string]any
		want string
	}{
		{"empty-host-ip", map[string]any{"published": 5432, "host_ip": ""}, "localhost"},
		{"ipv4-bind-all", map[string]any{"published": 5432, "host_ip": "0.0.0.0"}, "localhost"},
		{"ipv6-bind-all", map[string]any{"published": 5432, "host_ip": "::"}, "localhost"},
		{"explicit-loopback", map[string]any{"published": 5432, "host_ip": "127.0.0.1"}, "127.0.0.1"},
		{"custom-bind", map[string]any{"published": 5432, "host_ip": "192.168.1.10"}, "192.168.1.10"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, probable := application.ParseComposePortForTest(tc.raw)
			if !probable {
				t.Fatalf("expected probable=true")
			}
			if got.Host != tc.want {
				t.Errorf("Host = %q, want %q", got.Host, tc.want)
			}
		})
	}
}

func TestParseComposePort_UnknownTypes(t *testing.T) {
	t.Parallel()
	// Any input type outside int / int64 / float64 / string /
	// map[string]any must fall to the graceful warn path. The
	// list below covers all reasonable Compose-YAML decoding
	// surprises plus a few adversarial inputs.
	cases := []struct {
		name string
		raw  any
	}{
		{"nil", nil},
		{"bool", true},
		{"int-slice", []int{1, 2, 3}},
		{"interface-slice", []any{1, 2, 3}},
		{"empty-map", map[string]any{}},
		{"map-with-non-string-key", map[int]any{1: 2}},
		{"struct", struct{}{}},
		{"fractional-float", 5432.5},
		{"int-out-of-range-high", 70000},
		{"int-zero", 0},
		{"int-negative", -1},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, probable := application.ParseComposePortForTest(tc.raw)
			if probable {
				t.Errorf("%v: expected probable=false (graceful warn), got true", tc.raw)
			}
		})
	}
}

func TestParseComposePort_IPv6HostIPLimitation(t *testing.T) {
	t.Parallel()
	// Why: pin the documented MVP limitation. Compose accepts
	// "[::1]:5432:5432" but our naive split-by-":" parser yields
	// more than three segments. The contract is: this falls to
	// the warn path, not a hard fail. A future slice can add
	// bracket-aware parsing.
	cases := []string{
		"[::1]:5432:5432",
		"[2001:db8::1]:5432:5432",
	}
	for _, raw := range cases {
		raw := raw
		t.Run(raw, func(t *testing.T) {
			t.Parallel()
			_, probable := application.ParseComposePortForTest(raw)
			if probable {
				t.Errorf("%q: expected probable=false (IPv6 MVP limitation), got true", raw)
			}
		})
	}
}
