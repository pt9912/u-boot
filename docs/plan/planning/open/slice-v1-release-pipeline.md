# Slice V1: Release-Pipeline (Image-Publish + Trivy)

## Auslöser

ADR-0004 schließt zwei bewusst aus dem M2c-CI aus
(`LH-FA-PROJDOCS-005`):

1. Image-Publish-Workflow (`.github/workflows/publish.yml`) — kommt mit
   dem Release-Slice, gekoppelt an `LH-OPEN-002` (Paketierung).
2. Trivy-Image-Scan — optionaler dritter Job, der das
   `runtime`-Image scannt und CRITICAL/HIGH-Findings blockiert.

Solange `u-boot` keine ersten Releases hat, sind beide Lücken
akzeptabel. Sobald der erste GHCR-Tag fällt, müssen beide existieren.

## Aufhebungsbedingung

Erster Release (`v0.1.0` oder ähnlich) wird vorbereitet. Dieser Slice
liefert den Workflow + Trivy-Scan, bevor der Release gemacht wird.

## Akzeptanzkriterien

- `.github/workflows/publish.yml`:
  - Trigger: `push` von Tags `v*` (semver-konforme Tags).
  - Job baut das Runtime-Image über `make build`, pushed nach
    `ghcr.io/pt9912/u-boot:<version>` und `:latest`.
  - `permissions: contents: read, packages: write` (Per-Job minimal).
  - SHA-pinned `docker/login-action`, `docker/build-push-action` o. ä.
  - OCI-Labels aus `LH-FA-BUILD-002` sind im gepushten Image gesetzt.
- `.github/workflows/ci.yml` bekommt einen optionalen dritten Job
  `image-scan` (oder eigener Workflow `image-scan.yml`), der nach `make
  build` `trivy image --severity HIGH,CRITICAL --exit-code 1`
  ausführt. PR-blockierend, sobald aktiviert.
- `docs/user/quality.md` §6 wird um die neuen Workflows erweitert.
- `LH-OPEN-002` (Paketierung) wird mit der konkreten Distributions-
  Entscheidung geschlossen (mindestens GHCR-Image, ggf. Binary-
  Release).
- Eintrag in `carveouts.md` von „temporär" auf aufgehoben verschoben.

## Bezug

- Auslösende ADR: `0004-ci-system.md` Folgepunkte.
- Auslösende Spec: `LH-OPEN-002` Paketierung (heute offen).
- Inventar-Eintrag: [`carveouts.md`](../in-progress/carveouts.md) →
  Image-Publish + Trivy fehlen.
- Hängt von: erstem Release-Wunsch (vermutlich nach MVP-Closure
  `LH-MVP-001`).
