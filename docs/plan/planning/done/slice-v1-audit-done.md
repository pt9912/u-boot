# Slice V1: Audit-Done — drei V1-IDs verifizieren und schließen

## Auslöser

Drei V1-Spec-IDs sind vermutlich bereits durch den bestehenden
Code- und Doku-Stand erfüllt, aber nicht formal als ✅
markiert. v0.3.0-Milestone-Kickoff: bevor Code-Slices
(`addons-deps`, `keycloak`, `otel`) starten, diese drei Loose-
Ends als Aufwärm-Tranche verifizieren — Doku-Only, keine
Code-Änderung.

Spec-IDs:

- [`LH-FA-BUILD-006`](../../../../spec/lastenheft.md#lh-fa-build-006--aggregator-targets) Aggregator-Targets
- [`LH-NFA-MAINT-004`](../../../../spec/lastenheft.md#lh-nfa-maint-004--dokumentierte-schnittstellen) Dokumentierte Schnittstellen
- [`LH-NFA-PORT-003`](../../../../spec/lastenheft.md#lh-nfa-port-003--containerfreundlichkeit) Containerfreundlichkeit

## Aufhebungsbedingung

Pro ID:
- Spec-Text gegen aktuellen Repo-Stand verglichen.
- Konkrete Evidence dokumentiert (Datei + Zeilenbereich oder
  Slice/ADR-Ref).
- Anhang in der `roadmap.md`-MVP-Bilanz oder gleichwertiger
  V1-Bilanz-Tabelle gesetzt; v0.3.0-Milestone-Tabelle markiert
  [`slice-v1-audit-done`](slice-v1-audit-done.md) als ✅.

## Akzeptanzkriterien

- ✅ [`LH-FA-BUILD-006`](../../../../spec/lastenheft.md#lh-fa-build-006--aggregator-targets) — Makefile-Aggregator-Targets `gates`,
  `ci`, `fullbuild` sind vorhanden und decken die spec-
  geforderten Sub-Targets ab.
- ✅ [`LH-NFA-MAINT-004`](../../../../spec/lastenheft.md#lh-nfa-maint-004--dokumentierte-schnittstellen) — Add-on- und Template-Schnittstellen
  sind via ADRs + Port-Doc-Comments + Slice-Dokus dokumentiert.
- ✅ [`LH-NFA-PORT-003`](../../../../spec/lastenheft.md#lh-nfa-port-003--containerfreundlichkeit) — u-boot läuft in Container/Devcontainer;
  GHCR-Distroless-Image + container-aware `doctor` decken den
  Container-Pfad; Binary-Distribution decken den Host-Pfad; `init
  --devcontainer` erzeugt Devcontainer-Files.

## Audit-Evidence

### [`LH-FA-BUILD-006`](../../../../spec/lastenheft.md#lh-fa-build-006--aggregator-targets) Aggregator-Targets

Spec-Anforderung (V1): `gates`, `ci`, `fullbuild`-Targets im
Makefile; bei Subtarget-Failure Non-Zero-Exit mit klarer Fehler-
ursache.

Stand im `Makefile`:

```text
gates:      lint test coverage-gate docs-check   ## Inner-loop mandatory gates.
ci:         gates govulncheck image-scan         ## Gates plus govulncheck plus image-scan (mirrors ci.yml).
fullbuild:  ci build                             ## CI plus runtime image (full closure).
```

Sub-Target-Abdeckung:

- `gates` enthält `lint` + `test` + `coverage-gate` (plus
  `docs-check`, der über die Spec hinausgeht — additive
  Markdown-Link-Validierung, kein Spec-Bruch).
- `ci` enthält `gates` + `govulncheck` + `image-scan`. Der
  `govulncheck`-Target-Comment trägt bereits den
  [`LH-FA-BUILD-006`](../../../../spec/lastenheft.md#lh-fa-build-006--aggregator-targets)-Anker.
- `fullbuild` enthält `ci` + `build` — Closure-Lauf.

Non-Zero-Exit: standardmäßige Make-Semantik bricht bei jedem
fehlschlagenden Subtarget ab. Beim `make ci` ist die Reihenfolge
Lint → Test → Coverage → docs-check → govulncheck → image-scan;
jeder Subtarget-Failure bricht die Pipeline mit klarer Stage-
Ausgabe ab. SBOM bleibt optional (Spec-Wortlaut: „bleibt
optional").

✅ erfüllt.

### [`LH-NFA-MAINT-004`](../../../../spec/lastenheft.md#lh-nfa-maint-004--dokumentierte-schnittstellen) Dokumentierte Schnittstellen

Spec-Anforderung (V1): interne Schnittstellen für Add-ons und
Templates dokumentieren.

**Add-on-Schnittstelle:**

- [ADR-0008](../../adr/0008-plugin-system-statisch.md) — Add-on-
  System bleibt statisch (keine Plugins); vier Re-Eval-Trigger
  in §Folgepunkte.
- `internal/hexagon/port/driving/addservice.go` — `AddServiceUseCase`
  + `AddServiceRequest`/`AddServiceResponse` + Sentinels
  (`ErrServiceUnsupported`, `ErrServiceInconsistent`,
  `ErrProjectNotInitialized`) mit Doc-Comments.
- `internal/hexagon/port/driving/removeservice.go` — `RemoveServiceUseCase`
  mit dem [`LH-FA-ADD-005`](../../../../spec/lastenheft.md#lh-fa-add-005--mehrfaches-hinzufügen-verhindern)-State-Machine-Vertrag dokumentiert.
- [`done/slice-m5-add-postgres.md`](slice-m5-add-postgres.md) —
  detaillierte Add-on-Mechanik-Doku.
- [`done/slice-v1-add-remove.md`](slice-v1-add-remove.md) —
  [`LH-FA-ADD-007`](../../../../spec/lastenheft.md#lh-fa-add-007--service-entfernen)-State-Machine-Spiegelung.

Add-on-Katalog ist heute statisch (postgres only); neue Add-ons
werden direkt im Code (`isSupportedService`/`supportedServices`
in `addservice.go`) registriert. Per [ADR-0008](../../adr/0008-plugin-system-statisch.md) ist das die
verbindliche Add-on-Schnittstelle für Maintainer.

**Template-Schnittstelle:**

- [ADR-0009](../../adr/0009-template-format-yaml-files.md) —
  Template-Format YAML-Metadaten + `text/template`-Files;
  drei Implementierungs-Slices in §Folgepunkte.
- `internal/hexagon/port/driven/template_catalog.go` —
  `TemplateCatalog`-Port für Listing.
- `internal/hexagon/port/driven/template_files.go` — `TemplateFiles`-
  Port für File-Tree-Zugriff (Render-Pfad).
- `internal/hexagon/domain/template_metadata.go` —
  `TemplateMetadata`-Domain-Struct mit `template.yaml`-Schema-
  Doku-Comment.
- [`done/slice-v1-template-list.md`](slice-v1-template-list.md) +
  [`done/slice-v1-template-init.md`](slice-v1-template-init.md) —
  Template-Listing + Render-Pfad ausführlich beschrieben.

Beide Schnittstellen haben ADR-Plan-Anker (verbindliche
Entscheidung) + Port-Doc-Comments (Interface-Vertrag im Code) +
Slice-Dokus (How-To für Maintainer). ✅ erfüllt.

### [`LH-NFA-PORT-003`](../../../../spec/lastenheft.md#lh-nfa-port-003--containerfreundlichkeit) Containerfreundlichkeit

Spec-Anforderung (V1): u-boot selbst muss in Container oder
Devcontainer ausführbar sein.

**Container-Distribution:**

- Distroless-Image: `Dockerfile` baut `gcr.io/distroless/static-
  debian12:nonroot` (UID 65532).
- GHCR-Tags: `ghcr.io/pt9912/u-boot:0.1.0`, `:0.2.0`, `:latest`
  über `publish.yml` auf jedem `v*`-Tag-Push.
- [ADR-0007](../../adr/0007-distributionswege-ghcr.md) §Entscheidung
  setzt GHCR als primären Distributionsweg.

**Container-Aware `doctor`:**

- [`done/slice-v0.1.1-doctor-container-awareness.md`](slice-v0.1.1-doctor-container-awareness.md)
  ergänzt `driven.RuntimeEnvironment`-Port mit
  `/.dockerenv`/`/run/.containerenv`-Probes.
- Die vier Host-Prerequisite-Checks (`git.installed`,
  `docker.installed`, `docker.reachable`,
  `docker.compose.installed`) werden im Container-Modus mit
  `SeverityInfo` skipped statt als False-Positive-Errors zu
  feuern. Exit 0 bei sonst gesundem Projekt (vorher 11).

**Host-Pfad ergänzt:**

- [`done/slice-v2-binary-distribution.md`](slice-v2-binary-distribution.md)
  liefert sechs Plattformen (Linux/macOS/Windows × amd64/arm64)
  als GitHub-Release-Asset ab v0.2.0 — der host-native Pfad für
  `doctor` und andere Subkommandos.

**Devcontainer-Files:**

- `u-boot init --devcontainer` schreibt
  `.devcontainer/devcontainer.json` und `.devcontainer/Dockerfile`
  + setzt `devcontainer.enabled: true` in `u-boot.yaml`.
  Akzeptanztest-gepinnt durch [`LH-AK-005`](../../../../spec/lastenheft.md#lh-ak-005--devcontainer-flow)-Test (`bfe6416`).

✅ erfüllt. Container-Pfad (GHCR + container-aware `doctor`),
Host-Pfad (Binary) und Devcontainer-Generierung sind alle drei
abgedeckt.

## Tranchen

| T | Commit | Inhalt |
| - | ------ | ------ |
| T1 | dieser Commit | Audit-Evidence pro Spec-ID dokumentiert (siehe Sektion oben); Slice-Plan direkt in `done/` (Doku-only, kein open/-Zwischenstand); `roadmap.md` v0.3.0-Milestone-Tabelle markiert [`slice-v1-audit-done`](slice-v1-audit-done.md) als ✅; neuer V1-Audit-Sub-Block am Ende des `roadmap.md`-§MVP-Bilanz-Bereichs als Anker für künftige V1-Audits; `CHANGELOG.md ## [Unreleased]` Notes-Eintrag mit den drei verifizierten Spec-IDs. `make docs-check` grün. |

## Out of Scope

- **Code-Änderung:** dieser Slice ist Doku-Only. Wenn der Audit
  Lücken aufdeckt, werden sie in eigenen Slices behoben.
- **MVP-Spec-IDs:** nur die drei genannten V1-IDs. Andere V1-Folgen
  (Generators, --json/--dry-run, Logs) sind nicht Teil dieses
  Audits.

## Bezug

- Spec: [`LH-FA-BUILD-006`](../../../../spec/lastenheft.md#lh-fa-build-006--aggregator-targets), [`LH-NFA-MAINT-004`](../../../../spec/lastenheft.md#lh-nfa-maint-004--dokumentierte-schnittstellen), [`LH-NFA-PORT-003`](../../../../spec/lastenheft.md#lh-nfa-port-003--containerfreundlichkeit)
  (alle V1).
- v0.3.0-Milestone: erster Slice der „Add-on Catalogue Expansion"-
  Reihenfolge per [roadmap.md §v0.3.0](../in-progress/roadmap.md);
  ursprünglich als Aufwärm-Tranche geplant, kommt jetzt nach
  [`slice-v1-add-remove`](slice-v1-add-remove.md) (umgekehrte Reihenfolge gegenüber Plan
  ohne semantischen Verlust — die drei Audit-IDs hängen nicht an
  add-remove und vice versa).
- ADR-Anker für die Add-on-Schnittstelle:
  [ADR-0008](../../adr/0008-plugin-system-statisch.md).
- ADR-Anker für die Template-Schnittstelle:
  [ADR-0009](../../adr/0009-template-format-yaml-files.md).
- ADR-Anker für die Distribution:
  [ADR-0007](../../adr/0007-distributionswege-ghcr.md).
- Phase: V1 (laufend; v0.3.0-Milestone-Slice).
