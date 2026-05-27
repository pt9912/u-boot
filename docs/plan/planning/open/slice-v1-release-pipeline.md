# Slice V1: Release-Pipeline (Image-Publish + Trivy + Branch-Protection)

## Auslöser

ADR-0004 schließt drei bewusst aus dem M2c-CI aus
(`LH-FA-PROJDOCS-005`):

1. Image-Publish-Workflow (`.github/workflows/publish.yml`) — kommt mit
   dem Release-Slice, gekoppelt an `LH-OPEN-002` (Paketierung).
2. Trivy-Image-Scan — optionaler dritter Job, der das
   `runtime`-Image scannt und CRITICAL/HIGH-Findings blockiert.
3. **Branch-Protection** im GitHub-UI — Required-Status-Checks für
   `gates` und `security-gates` sind manuell zu aktivieren, sonst sind
   beide Jobs zwar grün, aber nicht PR-blockierend (`LH-QA-003`).
   Bei M3-Anker-Triage 2026-05-27 in diesen Slice gebündelt, weil die
   gleiche Sitzung (erster Release / erster externer PR) auch
   Image-Publish + Trivy aufsetzt; Standalone wäre Disziplin-Overhead.

Solange `u-boot` keine ersten Releases hat und kein externer
Contributor PRs öffnet, sind alle drei Lücken akzeptabel. Sobald der
erste GHCR-Tag fällt oder der erste externe PR ansteht, müssen sie
existieren.

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
- `docs/user/branch-protection.md` beschreibt Schritt-für-Schritt die
  einmalige UI-Aktivierung:
  - Settings → Branches → Add rule für `main`.
  - Required status checks: `gates` und `security-gates`.
  - Require PR before merging (Solo-Projekt: 0 Approvals, dokumentiert).
  - Block force-pushes auf `main`, block branch deletion.
  - Optional: linear history erzwingen.
- Optional `docs/user/branch-protection-ruleset.json` als
  GitHub-Repository-Ruleset-Export (importierbar via UI/API).
- README (de/en) Section „Setup" verweist auf die Branch-Protection-
  Checkliste.
- Alle drei Zeilen in `carveouts.md` (Image-Publish/Trivy, Branch-
  Protection, `LH-OPEN-002`) als gelöst markiert oder entfernt.

## Out of Scope

- DCO-Bot-Aktivierung (separater ADR-0004-Folgepunkt; lebt im
  GitHub-Marketplace, kein Repo-Artefakt).
- CODEOWNERS-Datei (eigener Slice, wenn Teilautoren dazukommen).

## Bezug

- Auslösende ADR: `0004-ci-system.md` Folgepunkte (3 davon).
- Auslösende Spec: `LH-OPEN-002` Paketierung, `LH-QA-003` PR-Blocking.
- Inventar-Eintrag: [`carveouts.md`](../in-progress/carveouts.md) →
  Image-Publish + Trivy + Branch-Protection.
- Hängt von: erstem Release-Wunsch oder erstem externen PR-Workflow.
- Absorbiert (2026-05-27): vormalig eigenständiges
  `slice-m3-branch-protection-checkliste.md`.
