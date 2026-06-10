# Slice V1: Release-Cut `v0.5.0`

> **Status:** ready to execute — T1–T3 sind reine Doku-/Versionsstring-
> Mechanik (ein Commit), T4 ist die Nutzer-Aktion (Tag-Push). Vorab in
> `open/` angelegt, damit der Cut morgen ohne Nachdenken läuft.
>
> **Scope-Entscheidung (vor T1 bestätigen):** v0.5.0 ist ein **schlankes
> Ein-Feature-Minor** — einziger auslieferbarer Inhalt seit `v0.4.0` ist
> `u-boot init --template ./pfad` (lokale User-Templates, [`LH-FA-TPL-003`](../../../../spec/lastenheft.md#lh-fa-tpl-003-eigene-templates),
> [`done/slice-later-local-templates.md`](../done/slice-later-local-templates.md)).
> Alles andere seit v0.4.0 ist `docs(plan)` ([ADR-0011](../../adr/0011-agent-harness-scaffolding.md)/[ADR-0012](../../adr/0012-devcontainer-egress-firewall.md) Proposed,
> Roadmap-APs, `examples.md`) — release-neutral. Falls vor dem Cut mehr
> gebündelt werden soll (ADR ratifizieren+bauen, Backlog-Slice), zuerst
> das tun; sonst v0.5.0 = local-templates.

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
> harmonised exit codes ([`LH-FA-TPL-003`](../../../../spec/lastenheft.md#lh-fa-tpl-003-eigene-templates), [ADR-0009](../../adr/0009-template-format-yaml-files.md)). Details below.

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
  ([`LH-FA-TPL-003`](../../../../spec/lastenheft.md#lh-fa-tpl-003-eigene-templates)), [ADR-0009](../../adr/0009-template-format-yaml-files.md).
- Roadmap: [`roadmap.md`](../in-progress/roadmap.md) §Aktueller Snapshot.
- `publish.yml` (Tag-Push-Trigger), `Makefile`/`Dockerfile`/`main.go`
  (Versionsstrings).
