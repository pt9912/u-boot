# ADR 0002: Hexagonale Architektur mit driving/driven-Split

## Status

Accepted

## Datum

2026-05-21

## Kontext

u-boot ist ein CLI-Tool, das eine wachsende Zahl externer Systeme
orchestriert: Docker Engine (`docker`/`docker compose`), Dateisystem
(Projektstruktur, Templates, Backups), YAML-/JSON-Codecs, Git und
perspektivisch externe Devcontainer-Feature-Quellen sowie
Template-Registries. (Ein Plugin-System wurde mit
[ADR-0008](0008-plugin-system-statisch.md) am 2026-05-31
ausgeschlossen — Add-on-System bleibt statisch.)

Ohne klare Schichten droht der Standard-Drift einer CLI:

- Cobra-Commands greifen direkt auf Docker-SDK, `os.WriteFile` und
  YAML-Encoder zu.
- Geschäftslogik (z. B. *„darf ein Service idempotent reaktiviert
  werden?"*, [`LH-FA-ADD-005`](../../../spec/lastenheft.md#lh-fa-add-005-mehrfaches-hinzufügen-verhindern)) verteilt sich auf Commands und Helper.
- Tests werden Integrationstests gegen die echte Docker-Engine oder
  enden als Mock-Wüsten.
- Wechsel der Docker-Bindung (z. B. von `os/exec` auf das offizielle
  Docker-SDK) trifft die gesamte Codebase.

Drei untersuchte Referenzprojekte (`k-deskflight` in Go, `m-trace` in
TypeScript, `grid-gym` in Python) lösen das durchgängig mit dem
hexagonalen Pattern (Ports & Adapters, Cockburn 2005). `m-trace` und
`grid-gym` nutzen zusätzlich einen `driving`/`driven`-Split unter `port/`
und `adapter/`, was die Verantwortungsrichtung an Verzeichnis-Ebene
sichtbar macht. `k-deskflight` (Go) hält die Adapter flach
(`adapter/k8s/`, `adapter/check/`).

Lastenheft-Bezug:

- [`LH-FA-ARCH-001`](../../../spec/lastenheft.md#lh-fa-arch-001-hexagonales-pattern) – Hexagonales Pattern (Pflicht).
- [`LH-FA-ARCH-002`](../../../spec/lastenheft.md#lh-fa-arch-002-schichten-und-verzeichnislayout) – Schichten und Verzeichnislayout.
- [`LH-FA-ARCH-003`](../../../spec/lastenheft.md#lh-fa-arch-003-import-regeln-und-enforcement) – Import-Regeln und Enforcement via `depguard`.
- [`LH-NFA-MAINT-001`](../../../spec/lastenheft.md#lh-nfa-maint-001-modulare-architektur) – modulare Architektur.
- [`LH-NFA-MAINT-003`](../../../spec/lastenheft.md#lh-nfa-maint-003-testbarkeit) – fachliche Funktionen automatisiert testbar.

## Entscheidung

u-boot folgt dem **hexagonalen Architektur-Pattern mit
`driving`/`driven`-Split** in Ports **und** Adaptern.

Konkrete Setzungen:

- Verzeichnislayout unter `internal/`:

  ```
  internal/
  ├── hexagon/
  │   ├── domain/
  │   ├── application/
  │   └── port/
  │       ├── driving/
  │       └── driven/
  └── adapter/
      ├── driving/
      └── driven/
  ```

- Wiring-Schicht ausschließlich in `cmd/uboot/` ([`LH-FA-BUILD-009`](../../../spec/lastenheft.md#lh-fa-build-009-repository-layout)).
- Import-Regeln aus [`LH-FA-ARCH-003`](../../../spec/lastenheft.md#lh-fa-arch-003-import-regeln-und-enforcement) werden im `lint`-Stage per
  `golangci-lint depguard` PR-blockierend durchgesetzt.
- `//nolint:depguard` ist verboten; Carveouts werden zentral in
  `.golangci.yml` mit `Why:`-Kommentar dokumentiert.
- Detail-Spezifikation (Beispiel-Inhalte je Schicht, depguard-Schema,
  Test-Pattern, Anti-Patterns) liegt in
  [`spec/architecture.md`](../../../spec/architecture.md).

Im MVP-Bootstrap (dieser Commit) sind die Verzeichnisse mit `README.md`
angelegt, aber leer; `depguard` ist in `.golangci.yml` aktiviert mit
leerem `rules`-Map. Die konkreten Regelblöcke werden mit dem ersten
fachlichen Inkrement scharf geschaltet (M3 `u-boot init`).

## Konsequenzen

Positiv:

- **Domäne und Application sind I/O-frei** und ohne Docker-Engine
  testbar (Fake-`DockerEngine`-Implementierung im `_test.go`-Paket).
- **Adapter-Wechsel** (z. B. `os/exec docker` → offizielles
  Docker-SDK) trifft nur `adapter/driven/docker/`, nicht `application`
  oder `domain`.
- **Verantwortungsrichtung sichtbar** auf Datei-Ebene durch
  `driving`/`driven`-Split — Newcomer erkennen sofort, wer reinruft
  (CLI) und wen das Application nach außen ruft (Docker, FS, YAML).
- **CI-blockierte Architektur-Regeln** verhindern Drift ab Tag 1, ohne
  manuelle Review-Disziplin zu erzwingen.
- ~~**Plugin-System** ([`LH-OPEN-003`](../../../spec/lastenheft.md#lh-open-003-plugin-system-entschieden)) lässt sich später als zusätzlicher
  `driven`-Port (`PluginRegistry`) sauber integrieren.~~
  **Überholt:** mit [ADR-0008](0008-plugin-system-statisch.md) am
  2026-05-31 entschieden, dass das Add-on-System statisch bleibt
  und kein `PluginRegistry`-Driven-Port eingeführt wird; die
  Re-Evaluation-Trigger sind in [ADR-0008](0008-plugin-system-statisch.md) §Folgepunkte
  dokumentiert. [ADR-0002](0002-hexagonale-architektur.md) wird durch [ADR-0008](0008-plugin-system-statisch.md) in diesem Punkt
  überschrieben ([`LH-FA-PROJDOCS-002`](../../../spec/lastenheft.md#lh-fa-projdocs-002-adr-format)).

Negativ / Trade-offs:

- **Mehr Boilerplate** als ein flacher Layout-Stil: jeder
  Use-Case braucht ein Driving-Port-Interface plus eine
  Application-Implementierung. Für ein CLI-Tool mit ~10 Subkommandos
  ist das überschaubar, aber spürbar.
- **Einarbeitung:** Beitragende ohne Hex-Erfahrung brauchen einmal das
  Layout-Diagramm aus `spec/architecture.md` und die Import-Tabelle
  aus [`LH-FA-ARCH-003`](../../../spec/lastenheft.md#lh-fa-arch-003-import-regeln-und-enforcement). `depguard`-Fehlermeldungen sind dabei
  selbsterklärend.
- **Slight overhead** im `lint`-Stage durch `depguard`-Auswertung;
  vernachlässigbar gegenüber `golangci-lint`-Gesamtlaufzeit.

Alternative-Optionen (verworfen):

- **Flacher Layout-Stil wie k-deskflight** (`adapter/<technologie>/`
  ohne `driving`/`driven`-Split): knapper, aber Verantwortungsrichtung
  ist nur per Konvention erkennbar. Für u-boot mit mehreren möglichen
  Driving-Pfaden (CLI heute, perspektivisch HTTP/Daemon) ist der
  explizite Split sauberer.
- **Clean Architecture / Onion Architecture:** semantisch nah am Hex,
  aber mit anderer Vokabular-Konvention. Hex ist in den drei
  Referenzprojekten bereits etabliert und vereinheitlicht die Begriffe
  org-weit.
- **Schicht-loser Layered-CLI-Stil:** würde gegen [`LH-NFA-MAINT-001`](../../../spec/lastenheft.md#lh-nfa-maint-001-modulare-architektur)
  (modulare Architektur) und [`LH-NFA-MAINT-003`](../../../spec/lastenheft.md#lh-nfa-maint-003-testbarkeit) (Testbarkeit)
  arbeiten.
