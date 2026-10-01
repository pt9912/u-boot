//go:build docker

// LH-AK-003 Keycloak-Acceptance-Flow Pin (slice-v1-keycloak T3).
//
// Spec §LH-AK-003: the sequence
//
//   u-boot init
//   u-boot add keycloak
//   u-boot up
//
// must succeed end-to-end. Acceptance criteria:
//
//   - Keycloak-Service exists in `compose.yaml`;
//   - `.env.example` lists `KEYCLOAK_ADMIN`, `KEYCLOAK_ADMIN_PASSWORD`
//     (with CHANGEME-Placeholder values per Spec §2351);
//   - the container reaches healthcheck status `healthy` within the
//     UpService timeout;
//   - the admin endpoint at `http://localhost:8080/` responds with
//     HTTP 200 or 302 (Spec §2352 toleriert beide — Keycloak
//     redirected ältere Versionen, neuere antworten direkt).
//
// Boot-Zeit-Carveout (slice-v1-keycloak.md §T3): Keycloak JVM-Boot
// + Realm-Init dauert 30–90 s warm, mit Cold-Image-Pull in CI
// deutlich länger. Postgres-Test verwendet 90 s — Keycloak braucht
// 4 Minuten Healthcheck-Timeout + 6 Minuten Gesamt-Context.
//
// Test prerequisites: docker CLI + compose plugin, test process
// shares network namespace mit docker daemon. `make test-docker`
// wirelt das via `test-docker-tools`-Stage + `--network=host`.

package e2e_test

import (
	"context"
	"testing"
	"time"
)

func TestE2E_LHAK003_KeycloakAcceptanceFlow(t *testing.T) {
	res := runAcceptanceFlow(t, acceptanceFlow{
		projectName: "t-uboot-e2e-kc",
		serviceName: "keycloak",
		envKeys:     []string{"KEYCLOAK_ADMIN", "KEYCLOAK_ADMIN_PASSWORD"},
		upTimeout:   4 * time.Minute,
		ctxTimeout:  6 * time.Minute,
	})

	// LH-AK-003 (3): Stabilized + keycloak healthcheck `healthy`.
	stabilizationCheck(t, res, "keycloak")

	// LH-AK-003 (4): admin endpoint reachable mit HTTP 200 oder 302.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	probeHTTPEndpoint(ctx, t, "http://localhost:8080/", 200, 302)
}
