# Slice V1: Docker-Integrationstests (Build-Tag-Pfad)

## Auslöser

`spec/architecture.md` §5 beschreibt eine Build-Tag-Konvention für
Adapter-Integrationstests gegen die echte Docker-Engine:

```
//go:build docker
```

mit `go test -tags docker ./...`. Aktuell gibt es weder einen
`adapter/driven/docker/`-Adapter noch einen entsprechenden CI-Stage
oder Make-Target — die Konvention ist nur dokumentiert
(`LH-FA-PROJDOCS-005`).

## Aufhebungsbedingung

Sobald `internal/adapter/driven/docker/` existiert und über
`port/driven.DockerEngine` aufgerufen wird (`LH-FA-UP-001..004`,
`LH-SA-DOCKER-001`/`-002`), wird:

1. mindestens ein Test mit `//go:build docker` angelegt, der gegen
   eine echte Docker-Engine läuft;
2. ein neuer Make-Target `make test-docker` erzeugt, das die
   getaggten Tests ausführt;
3. ein neuer Dockerfile-Stage oder CI-Job das Docker-Socket mountet
   und `make test-docker` aufruft (ergänzt `make ci`, nicht
   `make gates`).

## Akzeptanzkriterien

- `internal/adapter/driven/docker/docker_integration_test.go` mit
  `//go:build docker` läuft lokal via `go test -tags docker
  ./internal/adapter/driven/docker/...` gegen die Host-Docker-Engine.
- `make test-docker` führt den getaggten Pfad in einer Docker-in-
  Docker-Variante aus.
- `.github/workflows/ci.yml` bekommt einen neuen Job
  `integration-docker` (optional `continue-on-error` bis stabilisiert)
  oder ein eigener Workflow `.github/workflows/integration.yml`.
- `docs/user/quality.md` §2 Tests wird um den Docker-Pfad ergänzt.
- Zeile in `carveouts.md` entweder entfernen oder mit Verweis auf den Aufhebungs-Commit als gelöst markieren.

## Out of Scope

- Andere Build-Tags (`//go:build keycloak`, `//go:build otel`) —
  separate Slices pro Adapter.
- Kubernetes-Smoke — u-boot orchestriert Compose-Stacks, nicht
  Kubernetes; ein Cluster-Smoke-Pfad ist nicht im Roadmap-Bereich.

## Bezug

- Auslösende Spec: `spec/architecture.md` §5 Build-Tag-Konvention.
- Inventar-Eintrag: [`carveouts.md`](../in-progress/carveouts.md) →
  Docker-Integrationstests fehlen.
- Hängt von: M3+ Adapter-Slice (`internal/adapter/driven/docker/`).
