//go:build docker

// LH-AK-004 OpenTelemetry-Acceptance-Flow Pin (slice-v1-otel T3).
//
// Spec §LH-AK-004: the sequence
//
//   u-boot init
//   u-boot add otel
//   u-boot up
//
// must succeed end-to-end. Acceptance criteria:
//
//   - OpenTelemetry Collector is configured (Compose-Service +
//     `otel-collector-config.yaml`);
//   - the container reaches `running` ODER `healthy` within the
//     UpService timeout (Spec §2374 — Healthcheck nicht zwingend
//     für LH-AK-004; das Mindest-Setup ist healthcheck-frei);
//   - OTLP/gRPC ist auf `localhost:4317` erreichbar;
//   - OTLP/HTTP ist auf `localhost:4318` erreichbar.
//
// UpService-Timeout 60 s reicht (Collector-Boot < 5 s, kein
// JVM-Init wie bei Keycloak). Falls die CI das docker.io-Pull
// des Collector-Image als flaky erlebt (analog Quay/Keycloak),
// eskaliert dieser Test auf `//go:build docker && acceptance_
// extended` und der Folge-Slice `slice-v1-keycloak-ci-flake`
// schließt beide gleichzeitig.

package e2e_test

import (
	"context"
	"testing"
	"time"
)

func TestE2E_LHAK004_OtelAcceptanceFlow(t *testing.T) {
	res := runAcceptanceFlow(t, acceptanceFlow{
		projectName: "t-uboot-e2e-otel",
		serviceName: "otel",
		envKeys:     nil, // OTel default-Setup hat keine .env-Keys.
		upTimeout:   60 * time.Second,
		ctxTimeout:  3 * time.Minute,
	})

	// LH-AK-004 §2374 tolerates `running` OR `healthy`. The
	// generic stabilizationCheck-Helper from acceptance_helpers.go
	// asserts BOTH Stabilized + Healthcheck=="healthy" — that is
	// stricter than the spec. For OTel we therefore only assert
	// Stabilized (Compose hat den Service erfolgreich gebootet)
	// und überspringen die Healthcheck-Assertion.
	if !res.upResp.Result.Stabilized {
		t.Errorf("Stabilized = false, want true")
	}

	// LH-AK-004 (3+4): OTLP/gRPC + OTLP/HTTP reachable on localhost.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	dialTCP(ctx, t, "localhost:4317")
	dialTCP(ctx, t, "localhost:4318")
}
