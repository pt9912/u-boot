# Slice V1: Release-Pipeline (Image-Publish + Trivy + Branch-Protection)

## Auslöser

ADR-0004 schließt drei bewusst aus dem M2c-CI aus
(`LH-FA-PROJDOCS-005`):

1. Image-Publish-Workflow (`.github/workflows/publish.yml`) — kommt mit
   dem Release-Slice, gekoppelt an `LH-OPEN-002` (Paketierung).
2. Trivy-Image-Scan — optionaler dritter Job, der das
   `runtime`-Image scannt und CRITICAL/HIGH-Findings blockiert.
3. **Branch-Protection** im GitHub-UI — Required-Status-Checks für die
   tatsächlichen GitHub-Check-Namen `gates (lint + test +
   coverage-gate)` und `security-gates (govulncheck)` sind manuell zu
   aktivieren, sonst sind beide Jobs zwar grün, aber nicht
   PR-blockierend (`LH-QA-003`).
   Bei M3-Anker-Triage 2026-05-27 in diesen Slice gebündelt, weil die
   gleiche Sitzung (erster Release / erster externer PR) auch
   Image-Publish + Trivy aufsetzt; Standalone wäre Disziplin-Overhead.

Solange `u-boot` keine ersten Releases hat und kein externer
Contributor PRs öffnet, sind alle drei Lücken akzeptabel. Der Slice hat
zwei Auslösepfade:

- **Erster externer PR:** Branch-Protection muss vor dem Merge
  dokumentiert und im GitHub-UI aktiviert sein. Image-Publish und Trivy
  dürfen offen bleiben, solange kein Release vorbereitet wird.
- **Erster Release / erster GHCR-Tag:** Image-Publish, Trivy-Scan und
  Branch-Protection müssen vollständig umgesetzt sein, bevor der Release
  gemacht wird.

## Aufhebungsbedingung

Erster Release (`v0.1.0` oder ähnlich) wird vorbereitet. Dieser Slice
liefert den Workflow + Trivy-Scan, bevor der Release gemacht wird.

Vorheriger Teilabschluss ist erlaubt, wenn zuerst ein externer PR
ansteht: Dann wird nur der Branch-Protection-Teil umgesetzt und die
Release-Teile bleiben als offene Restarbeit in diesem Slice.

## Akzeptanzkriterien

- `.github/workflows/publish.yml`:
  - Trigger: `push` von Tags `v*`.
  - Früher Validierungsstep prüft Publish-SemVer-Tags
    (`vMAJOR.MINOR.PATCH`, ggf. mit SemVer-Prerelease, aber ohne
    Build-Metadaten, weil `+` kein gültiges Docker/GHCR-Tag-Zeichen ist)
    und bricht bei anderen `v*`-Tags vor Login/Build/Push ab.
  - Job baut das Runtime-Image über `make build`, pushed nach
    `ghcr.io/pt9912/u-boot:<version>`; `:latest` wird nur für stabile
    `vMAJOR.MINOR.PATCH`-Tags gesetzt, nicht für Prereleases.
  - `permissions: contents: read, packages: write` (Per-Job minimal).
  - SHA-pinned `docker/login-action`, `docker/build-push-action` o. ä.
  - OCI-Labels aus `LH-FA-BUILD-002` sind im gepushten Image gesetzt.
- `.github/workflows/ci.yml` bekommt einen optionalen dritten Job
  `image-scan` (oder eigener Workflow `image-scan.yml`), der nach `make
  build` `trivy image --severity HIGH,CRITICAL --exit-code 1`
  ausführt.
- Branch-Protection nimmt alle aktivierten PR-Gates auf:
  - Mindestmenge beim externen-PR-Pfad: `gates (lint + test +
    coverage-gate)` und `security-gates (govulncheck)`; falls der
    Workflow später auf kürzere Job-Namen umgestellt wird, muss die
    Checkliste die dann tatsächlichen GitHub-Check-Namen verwenden.
  - Sobald `image-scan` existiert: zusätzlich `image-scan`.
- `docs/user/quality.md` §4 und §6 werden um die neuen Workflows
  erweitert; die bisherige Aussage "Trivy/SBOM folgen mit dem
  Release-Slice" wird aktualisiert oder entfernt.
- `LH-OPEN-002` (Paketierung) wird für den GHCR-Anteil konkret
  entschieden. Weitere Distributionswege (Binary-Release, Homebrew,
  Debian/RPM, npm/pip) werden entweder explizit verworfen/vertagt oder
  bekommen eigene Slice-Pläne; `LH-OPEN-002` wird nur geschlossen, wenn
  diese Restwege ebenfalls entschieden sind.
- `docs/user/branch-protection.md` beschreibt Schritt-für-Schritt die
  einmalige UI-Aktivierung:
  - Settings → Branches → Add rule für `main`.
  - Required status checks: `gates (lint + test + coverage-gate)`,
    `security-gates (govulncheck)` und, sobald vorhanden, `image-scan`.
  - Require PR before merging (Solo-Projekt: 0 Approvals, dokumentiert).
  - Block force-pushes auf `main`, block branch deletion.
  - Optional: linear history erzwingen.
- Optional `docs/user/branch-protection-ruleset.json` als
  GitHub-Repository-Ruleset-Export (importierbar via UI/API).
- README (de/en) Section „Setup" verweist auf die Branch-Protection-
  Checkliste.
- Alle drei Zeilen in `carveouts.md` (Image-Publish/Trivy, Branch-
  Protection, `LH-OPEN-002`) passend aktualisieren:
  - Branch-Protection entfernen/als gelöst markieren, sobald der
    externe-PR-Pfad umgesetzt ist.
  - Image-Publish/Trivy entfernen/als gelöst markieren, sobald der
    Release-Pfad umgesetzt ist.
  - `LH-OPEN-002` nur entfernen/als gelöst markieren, wenn alle
    Distributionswege entschieden sind; andernfalls den Carveout auf
    die verbleibenden Wege reduzieren.

## Tranchen-Schnitt

Fünf Tranchen, in Reihenfolge implementierbar. Stand 2026-05-31
ergänzt im Zuge der Release-Vorbereitung (`v0.1.0`).

### T1 — ADR-0007 Distributionswege + `LH-OPEN-002`-Update

- `docs/plan/adr/0007-distributionswege-ghcr.md` neu: GHCR primär;
  Binary / Homebrew / Distro-Pakete vertagt mit Trigger-Slices;
  npm / pip verworfen.
- `spec/lastenheft.md` §14 `LH-OPEN-002`-Abschnitt + Übersichts-
  tabellen-Zeile auf den Stand „GHCR entschieden, Restwege
  vertagt/verworfen" gehoben.

**DoD T1:**
- [x] ADR-0007 angelegt; Mindest-Abschnitte erfüllt (`adr/README.md`).
- [x] `LH-OPEN-002` §14 enthält Entscheidungs-Tabelle + ADR-Verweis.
- [x] `make gates` grün.
- [x] T1 ✅ `0f64938`.

### T2 — `.github/workflows/publish.yml` (GHCR Image-Publish)

- Trigger `push` von Tags `v*`; früher SemVer-Validate-Step lehnt
  alles ab, was nicht `vMAJOR.MINOR.PATCH(-prerelease)?` ist
  (Build-Metadaten mit `+` rejecten).
- Job baut Runtime-Image via `make build`, pushed nach
  `ghcr.io/pt9912/u-boot:<version>`; `:latest` nur für stable Tags.
- SHA-pinned `docker/login-action` + `docker/build-push-action`;
  `permissions: contents: read, packages: write` per-Job.
- OCI-Labels aus `LH-FA-BUILD-002` im gepushten Image verifizieren.

**DoD T2:**
- [x] Workflow-File angelegt; Validate-Step + Push-Step getrennt.
- [x] Probe-Lauf lokal: `make build` + `docker tag` + Label-Inspect
  gegen die 5 LH-FA-BUILD-002 Labels grün.
- [x] `make gates` grün.
- [x] T2 ✅ `93b703e`.

### T3 — Trivy-Image-Scan

- Eigener Workflow `.github/workflows/image-scan.yml` ODER dritter
  CI-Job in `ci.yml`. Entscheidung in der Tranche begründen.
- Nach `make build`: `trivy image --severity HIGH,CRITICAL --exit-code 1`.
- SHA-pinned Action.

**DoD T3:**
- [ ] Workflow / Job angelegt; Probe-Lauf grün auf `main`.
- [ ] `make gates` grün.
- [ ] DoD-Line: `T3 ✅ <commit-hash>`.

### T4 — Doku-Sync

- `docs/user/quality.md` §4 (Trivy-Folgepunkt-Satz aktualisiert) und
  §6 (Image-Publish-Folgepunkt-Satz aktualisiert, neue Workflows
  verlinkt).
- `docs/user/branch-protection.md` Required-Status-Checks um
  `image-scan` ergänzen.
- `README.md` + `README.de.md` Setup-Section um Branch-Protection-
  Verweis prüfen/ergänzen.
- Optional: `docs/user/branch-protection-ruleset.json` als
  importierbarer Repository-Ruleset-Export.

**DoD T4:**
- [ ] `quality.md` §4 + §6 sind ohne offene Carveout-Sätze.
- [ ] `branch-protection.md` listet `image-scan`.
- [ ] READMEs verlinken Branch-Protection.
- [ ] `make gates` grün.
- [ ] DoD-Line: `T4 ✅ <commit-hash>`.

### T5 — Slice-Closure

- `docs/plan/planning/in-progress/carveouts.md`:
  - Image-Publish/Trivy-Zeile entfernen.
  - `LH-OPEN-002`-Zeile auf verbleibende Wege (Binary / Homebrew /
    Distro-Pakete) reduzieren, mit Verweis auf ADR-0007.
- Slice-Plan von `open/` nach `done/` verschieben; DoD-Lines auf
  Commit-Hashes auflösen.
- `docs/plan/planning/in-progress/roadmap.md`: V1-Liste der Trigger-
  getriebenen Slices auf den neuen Stand; MVP-Bilanz V1-Phase
  aktualisieren.

**DoD T5:**
- [ ] `carveouts.md` entsprechend bereinigt.
- [ ] Slice-Plan in `done/` mit allen T1..T5 DoD-Lines.
- [ ] Roadmap-Zeile auf Done; MVP-Bilanz aktualisiert.
- [ ] `make gates` grün.
- [ ] DoD-Line: `T5 ✅ <commit-hash>`.

## Out of Scope

- DCO-Bot-Aktivierung (separater ADR-0004-Folgepunkt; lebt im
  GitHub-Marketplace, kein Repo-Artefakt).
- CODEOWNERS-Datei (eigener Slice, wenn Teilautoren dazukommen).

## Bezug

- Auslösende ADR: `0004-ci-system.md` Folgepunkte (3 davon).
- Auslösende Spec: `LH-OPEN-002` Paketierung, `LH-QA-003` PR-Blocking.
- Inventar-Eintrag: [`carveouts.md`](../in-progress/carveouts.md) →
  Image-Publish + Trivy (Branch-Protection ist mit dem
  Teilabschluss-Commit 2026-05-27 aufgehoben, siehe
  [`docs/user/branch-protection.md`](../../../user/branch-protection.md)).
- Hängt von: erstem Release-Wunsch (GHCR/Trivy/`LH-OPEN-002`).
- **Teilabschluss 2026-05-27:** Branch-Protection-Checkliste in
  `docs/user/branch-protection.md` veröffentlicht. Restscope
  (Image-Publish + Trivy + `LH-OPEN-002`) bleibt offen bis zum
  ersten Release-Wunsch.
- Absorbiert (2026-05-27): vormalig eigenständiges
  `slice-m3-branch-protection-checkliste.md`.
