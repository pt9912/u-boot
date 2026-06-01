//go:build docker

// LH-AK-002 PostgreSQL-Acceptance-Flow Pin (M6-docker-int Sub-T3).
//
// Spec §LH-AK-002: the sequence
//
//   u-boot init
//   u-boot add postgres
//   u-boot up
//
// must succeed end-to-end. Acceptance criteria:
//
//   - PostgreSQL-Service exists in `compose.yaml`;
//   - `.env.example` lists `POSTGRES_USER`, `POSTGRES_PASSWORD`,
//     `POSTGRES_DB`;
//   - the container reaches healthcheck status `healthy` within 90 s;
//   - port 5432 is reachable on `localhost`.
//
// Refactored in slice-v1-keycloak T3 to share its scaffolding
// (init + add + up + Compose-Down cleanup + pre-up env-block
// assertions) with the new Keycloak Acceptance test via
// [runAcceptanceFlow] in acceptance_helpers.go.
//
// Test prerequisites (slice §Strukturelle Bedingungen): docker CLI
// + compose plugin available; test process and docker daemon must
// share a network namespace. `make test-docker` satisfies both via
// the `test-docker-tools` Dockerfile stage and `--network=host`.

package e2e_test

import (
	"context"
	"testing"
	"time"
)

func TestE2E_LHAK002_PostgresAcceptanceFlow(t *testing.T) {
	res := runAcceptanceFlow(t, acceptanceFlow{
		projectName: "t-uboot-e2e-acc",
		serviceName: "postgres",
		envKeys:     []string{"POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB"},
		upTimeout:   90 * time.Second,
		ctxTimeout:  5 * time.Minute,
	})

	// LH-AK-002 (3): Stabilized + postgres healthcheck `healthy`.
	stabilizationCheck(t, res, "postgres")

	// LH-AK-002 (4): port 5432 reachable on localhost. With a
	// healthcheck UpService treats port-probe failure as warn-only,
	// so stabilizationCheck above does not by itself prove the port
	// is reachable — dial directly to lock that part of the contract.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dialTCP(ctx, t, "localhost:5432")
}
