# Slice V1: Release-Cut `v0.5.0`

> **Status (2026-07-25):** T1–T3 **ausgeführt**; **T4 offen** — Tag-Push ist
> Nutzer-Aktion. Der Slice bleibt bis dahin in `in-progress/`, siehe
> §Abweichung.
>
> **Scope-Korrektur (2026-07-25, vor T1):** Der Plan beschrieb v0.5.0 als
> **Ein-Feature-Minor** (nur local-templates). Das war beim Schreiben richtig
> und ist es nicht mehr: Seither kam der Go-Toolchain-Bump 1.26.4 → 1.26.5
> dazu, der **CVE-2026-39822** (HIGH, `os.Root` Directory-Traversal) im
> ausgelieferten Runtime-Image schließt. Wer heute `:v0.4.0` zieht, bekommt das
> verwundbare Image — **das** ist der eigentliche Anlass des Cuts, nicht das
> Feature. v0.5.0 ist damit ein **Security-Release mit einem Feature**.
> Ausgeliefert werden beide Einträge der `[Unreleased]`-Sektion:
> `### Security` (Toolchain) und `### Added` (lokale Vorlagen,
> [`LH-FA-TPL-003`](../../../../spec/lastenheft.md#lh-fa-tpl-003--eigene-templates),
> [`done/slice-later-local-templates.md`](../done/slice-later-local-templates.md)).
> Alles andere seit v0.4.0 ist Doku-/Harness-Arbeit — release-neutral.

## Auslöser

Seit dem `v0.4.0`-Release (Tag auf `bce886f`, 2026-06-08) ist genau ein
auslieferbares Feature gelandet. Die CHANGELOG-`## [Unreleased]`-Sektion
trägt den zugehörigen `### Added`-Eintrag bereits (local-templates).
Der Cut folgt dem etablierten Muster aus
[`done/slice-v1-release-cut-v0.4.0.md`](../done/slice-v1-release-cut-v0.4.0.md)
(T1–T3 Doku/Versionsstrings als ein Commit, T4 Nutzer-Aktion).

## Aufhebungsbedingung

`git tag v0.5.0 && git push origin v0.5.0` triggert `publish.yml`
(GHCR-Push `ghcr.io/pt9912/u-boot:0.5.0` + `:latest` + Binary-Upload für
sechs Plattformen); `u-boot --version` eines Release-Binaries zeigt
`0.5.0`; roadmap §Snapshot trägt v0.5.0 als `released`.

## Akzeptanzkriterien (T1–T3, exakte Anker)

- **CHANGELOG.md**: `## [Unreleased]` bleibt als leerer Anker; darunter
  neu `## [0.5.0] - <Tag-Datum>` mit kurzem Lead-Absatz (Vorlage unten),
  gefolgt vom bestehenden local-templates-`### Added`-Block (der aus
  `[Unreleased]` herunterwandert). Optional eine `### Documentation`-Zeile
  für `docs/user/examples.md`.
- **Compare-Links** (CHANGELOG-Fuß): `[Unreleased]:` auf
  `…/compare/v0.5.0...HEAD` umstellen; neue Zeile
  `[0.5.0]: https://github.com/pt9912/u-boot/compare/v0.4.0...v0.5.0`.
- **Versionsstrings** `0.4.0-dev` → `0.5.0-dev` an **genau drei** Stellen:
  - `cmd/uboot/main.go:43` — `var version = "0.4.0-dev"`
  - `Makefile:27` — `VERSION ?= 0.4.0-dev` (+ Kommentar-Erwähnungen Z. 23/26)
  - `Dockerfile:39` — `ARG UBOOT_VERSION=0.4.0-dev` (+ Kommentar Z. 35/37)
- **README.md**: Status-Block (Z. 11–12, „v0.4.0 released …" → v0.5.0)
  + neue Releases-Tabellen-Zeile nach Z. 174 (`v0.4.0`-Zeile).
- **README.de.md**: Status-Block (Z. 11) + Releases-Zeile nach Z. 179.
- **roadmap.md §Aktueller Snapshot**: neue v0.5.0-Zeile nach der
  v0.4.0-Zeile (Z. 16); v0.4.0-Backlog-Tabelle ggf. um die jetzt
  ausgelieferte Zeile bereinigen.
- `make gates` grün.

## Lead-Absatz-Vorlage (CHANGELOG `## [0.5.0]`)

> Fifth release. Local filesystem templates: `u-boot init --template
> ./path` resolves a project template from the real filesystem (not just
> the built-in catalogue), with a pure, platform-independent name-vs-path
> classification, a shared `template.yaml` parser, a symlink guard, and
> harmonised exit codes ([`LH-FA-TPL-003`](../../../../spec/lastenheft.md#lh-fa-tpl-003--eigene-templates), [ADR-0009](../../adr/0009-template-format-yaml-files.md)). Details below.

## Releases-Tabellen-Zeile (Vorlage)

README.md:

```
| `v0.5.0` | <Datum> | "Local templates" — `u-boot init --template ./path` resolves a project from a local directory (`LH-FA-TPL-003`), alongside the built-in catalogue. [GitHub release](https://github.com/pt9912/u-boot/releases/tag/v0.5.0). |
```

README.de.md:

```
| `v0.5.0` | <Datum> | „Lokale Templates" — `u-boot init --template ./pfad` rendert ein Projekt aus einem lokalen Verzeichnis (`LH-FA-TPL-003`), neben dem eingebauten Katalog. [GitHub-Release](https://github.com/pt9912/u-boot/releases/tag/v0.5.0). |
```

## Tranchen

| T | Inhalt |
| - | ------ |
| T1–T3 | (ein Commit) **CHANGELOG** `[Unreleased]` → `[0.5.0] - <Datum>` + Lead + Compare-Links. **Versionsstrings** `0.4.0-dev`→`0.5.0-dev` (main.go/Makefile/Dockerfile). **README.{md,de.md}** Status + Releases-Zeile. **roadmap** §Snapshot. `make gates` grün. Slice-Doc `open/` → `done/`. |
| T4 | **Nutzer-Aktion:** (a) `## [0.5.0] - <Datum>` auf das tatsächliche Tag-Datum setzen. (b) sicherstellen, dass `main` auf `origin/main` ist; ersten grünen CI-Lauf abwarten. (c) `git tag v0.5.0 && git push origin v0.5.0` → `publish.yml` (GHCR + Binaries). (d) Post-Push: roadmap §Snapshot v0.5.0 auf `released` + Tag-Commit-Hash + Datum; `## [Unreleased]` bleibt leerer Anker. |

## Abweichung vom Plan: Lifecycle und Roadmap-Anker

Zwei Anker des Plans waren beim Ausführen veraltet — beide dokumentiert statt
still umgangen:

**1. Der Slice bleibt in `in-progress/`, bis T4 erledigt ist.** Die
Tranchen-Tabelle sieht den Move nach `done/` schon am Ende von T1–T3 vor. Das
widerspricht seit dem 2026-07-25 dem `planning`-Sensor: Er verlangt, dass die
Roadmap den Ruhe-Marker genau dann trägt, wenn **kein** Slice in `in-progress/`
liegt. Die Roadmap führt v0.5.0 als aktive Welle, weil der Release bis zum
Tag-Push nicht fertig ist — also gehört auch der Slice dorthin. Der Move nach
`done/` ist Teil von T4.

**2. „roadmap.md §Aktueller Snapshot" gibt es nicht mehr.** Die Roadmap folgt
seit der Regelwerk-Adoption dem Wellen-Template (`MR-003`); der Snapshot-
Abschnitt ist durch §Aktuelle Welle, §Meilensteine und §Abgeschlossene Wellen
ersetzt. Nachgezogen wurden stattdessen: §Aktuelle Welle auf `v0.5.0` (mit dem
Hinweis, dass T4 aussteht), die Vorwelle `welle-gate-ausbau-v0.51` nach
§Abgeschlossene Wellen, und der Abhängigkeitsgraph.

Ergänzend nachgezogen, weil es sonst zwischen Release und Doku driftet:
[`docs/user/benutzerhandbuch.md`](../../../user/benutzerhandbuch.md) trägt jetzt
Software-Version `v0.5.0` und beschreibt `--template ./pfad`; der Abschnitt
„Noch nicht enthalten" ist entfallen. Das Handbuch entstand einen Tag vor dem
Cut und hatte die unveröffentlichte Funktion bewusst nur als Ausblick geführt.

## Out of Scope

- **Mehr Features bündeln**: Scope-Frage (siehe Status-Block) — falls ja,
  vor T1 erledigen; dieser Plan setzt local-templates-only voraus.
- **[ADR-0011](../../adr/0011-agent-harness-scaffolding.md)/[ADR-0012](../../adr/0012-devcontainer-egress-firewall.md) ratifizieren**: bleiben Proposed; keine v0.5.0-Blocker.
- **Spec-`Datum`/`Version`-Header in `spec/lastenheft.md`**: bleibt
  unverändert (Spec-Version ≠ Release-Version, wie in v0.1–v0.4).

## Bezug

- Prozess-Vorbild:
  [`done/slice-v1-release-cut-v0.4.0.md`](../done/slice-v1-release-cut-v0.4.0.md).
- Inhalt: [`done/slice-later-local-templates.md`](../done/slice-later-local-templates.md)
  ([`LH-FA-TPL-003`](../../../../spec/lastenheft.md#lh-fa-tpl-003--eigene-templates)), [ADR-0009](../../adr/0009-template-format-yaml-files.md).
- Roadmap: [`roadmap.md`](../in-progress/roadmap.md) §Aktueller Snapshot.
- `publish.yml` (Tag-Push-Trigger), `Makefile`/`Dockerfile`/`main.go`
  (Versionsstrings).
