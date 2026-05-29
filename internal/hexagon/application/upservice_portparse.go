package application

import (
	"strconv"
	"strings"
)

// portProbeTarget is one TCP endpoint the [UpService] polling loop
// can probe via [driven.NetProbe.DialTCP]. Compose service ports go
// through [parseComposePort] before becoming a probe target.
type portProbeTarget struct {
	// Host is the textual host address — typically "localhost",
	// "127.0.0.1", or an explicit host_ip from the compose long
	// port-syntax. Compose's bind-to-all sentinels ("0.0.0.0",
	// "::", "") are normalized to "localhost".
	Host string
	// Port is the host-side TCP port number in [1, 65535].
	Port int
}

// parseComposePort interprets one element of a Compose `ports:` array
// and returns the TCP probe target plus a `probable` flag that tells
// the caller whether the element can be probed at all.
//
// The function is intentionally **lossy** and **non-failing**:
// Compose accepts a wide variety of port-syntax forms (some of which
// — range syntax, UDP, IPv6 host_ip with bracket-quoting — are not
// individually probable on `localhost`), and forcing a hard error
// on those forms would block `u-boot up` for a service that is
// otherwise healthy. Instead, the caller emits a Severity-warn
// diagnostic (ID prefix `up.port.<service>.<index>`) and proceeds
// without blocking stabilization on the port (LH-FA-UP-001 §969).
//
// Eight syntax cases — the M6 slice plan's robust-port-parsing
// table — are accepted:
//
//  1. naked integer `5432` → probe `localhost:5432`
//  2. string short form `"5432:5432"` → probe `localhost:5432`
//     (Compose convention: `HOST:CONTAINER`, so left is host)
//  3. string short form with host bind `"127.0.0.1:5432:5432"` →
//     probe `127.0.0.1:5432`
//  4. string short form with protocol `"5432:5432/udp"` →
//     probable=false (non-TCP)
//  5. string short form with range `"5000-5010:5000-5010"` →
//     probable=false (not single-port)
//  6. long-syntax mapping `{target:5432, published:5432, protocol:tcp}`
//     → probe `localhost:5432`; `protocol:udp` → probable=false
//  7. long-syntax with `host_ip:127.0.0.1` → probe `127.0.0.1:5432`
//  8. unknown / malformed (slice, bool, empty string, non-numeric
//     host port, IPv6 with bracket-quoting, missing `published`)
//     → probable=false (graceful warn, not a fail)
//
// IPv6 host_ip is a known MVP limitation: Compose accepts
// `"[::1]:5432:5432"` but the naive split-by-":" parser here yields
// more than three segments and falls to the unknown-form branch.
// A future slice can add bracket-aware parsing; for M6 the
// warn-diagnostic path is the safe fallback.
func parseComposePort(raw any) (target portProbeTarget, probable bool) {
	switch v := raw.(type) {
	case int:
		return makeProbe("localhost", v)
	case int64:
		return makeProbe("localhost", int(v))
	case float64:
		// YAML decodes bare numbers as float64 in some codecs;
		// reject fractional values (5432.5 is not a port).
		if v != float64(int(v)) {
			return portProbeTarget{}, false
		}
		return makeProbe("localhost", int(v))
	case string:
		return parseComposePortString(v)
	case map[string]any:
		return parseComposePortMapping(v)
	default:
		return portProbeTarget{}, false
	}
}

// parseComposePortString handles the three short-form variants:
// `"PORT"`, `"HOST_PORT:CONTAINER_PORT"`,
// `"HOST_IP:HOST_PORT:CONTAINER_PORT"`, plus protocol suffixes and
// range syntax.
func parseComposePortString(s string) (portProbeTarget, bool) {
	if s == "" {
		return portProbeTarget{}, false
	}
	// Protocol suffix on the rightmost segment ("5432:5432/udp").
	// Only "tcp" (case-insensitive) is probable; anything else
	// (udp, sctp) returns probable=false.
	if idx := strings.LastIndex(s, "/"); idx != -1 {
		proto := strings.ToLower(s[idx+1:])
		if proto != "tcp" {
			return portProbeTarget{}, false
		}
		s = s[:idx]
	}
	parts := strings.Split(s, ":")
	var hostIP, hostPort string
	switch len(parts) {
	case 1, 2:
		// "5432" or "5432:5432" — left segment is the host port.
		hostPort = parts[0]
		hostIP = "localhost"
	case 3:
		// "127.0.0.1:5432:5432" — host_ip then host_port then
		// container_port.
		hostIP = parts[0]
		hostPort = parts[1]
	default:
		return portProbeTarget{}, false
	}
	if strings.Contains(hostPort, "-") {
		// Range syntax ("5000-5010"); we cannot pin a single
		// probe target, so emit a warn instead.
		return portProbeTarget{}, false
	}
	port, err := strconv.Atoi(hostPort)
	if err != nil {
		return portProbeTarget{}, false
	}
	return makeProbe(hostIP, port)
}

// parseComposePortMapping handles the long-syntax variant: a YAML
// mapping with `target`, `published`, `protocol`, and `host_ip`
// keys. `published` is the host-side port (required for a probe);
// `target` is the container-side port (used by Compose for the
// internal mapping, not by us); `protocol` defaults to tcp;
// `host_ip` defaults to localhost.
func parseComposePortMapping(m map[string]any) (portProbeTarget, bool) {
	if !isTCPProtocol(m["protocol"]) {
		return portProbeTarget{}, false
	}
	port, ok := extractPublishedPort(m["published"])
	if !ok {
		return portProbeTarget{}, false
	}
	host := "localhost"
	if hostIP, ok := m["host_ip"].(string); ok && hostIP != "" {
		host = hostIP
	}
	return makeProbe(host, port)
}

// isTCPProtocol returns true for the absence of a `protocol` key
// (Compose defaults to TCP), or for an explicit "tcp" (case-
// insensitive). Returns false for any other value (udp, sctp, or a
// non-string).
func isTCPProtocol(proto any) bool {
	if proto == nil {
		return true
	}
	s, ok := proto.(string)
	if !ok {
		return false
	}
	return strings.EqualFold(s, "tcp")
}

// extractPublishedPort accepts the same numeric forms as the top-
// level [parseComposePort] plus a string form (Compose docs allow
// `published: "5432"`). Range strings and non-numeric strings
// return ok=false.
func extractPublishedPort(raw any) (int, bool) {
	switch v := raw.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case float64:
		if v != float64(int(v)) {
			return 0, false
		}
		return int(v), true
	case string:
		if strings.Contains(v, "-") {
			return 0, false
		}
		p, err := strconv.Atoi(v)
		if err != nil {
			return 0, false
		}
		return p, true
	default:
		return 0, false
	}
}

// makeProbe builds the final [portProbeTarget] and applies the
// "bind-to-all → localhost" normalization. Returns probable=false
// for non-positive port numbers (Compose accepts 0 as "any port",
// but the polling loop has no host port to dial at that point).
func makeProbe(host string, port int) (portProbeTarget, bool) {
	if port <= 0 || port > 65535 {
		return portProbeTarget{}, false
	}
	switch host {
	case "", "0.0.0.0", "::":
		host = "localhost"
	}
	return portProbeTarget{Host: host, Port: port}, true
}
