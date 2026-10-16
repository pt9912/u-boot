# Verification Evidence: slice-harness-ai-harness-init-adoption — 2026-10-16

Form nach [`harness/verification.md`](../../harness/verification.md). Verifier in frischem Kontext; geprüft gegen den Stand `792c53e`, Arbeitsbaum sauber.

## Verification Evidence

Scope:
- Slice: [`slice-harness-ai-harness-init-adoption`](../plan/planning/done/slice-harness-ai-harness-init-adoption.md)
- IDs: [`MR-012`](../../harness/conventions.md#mr-012), [`MR-011`](../../harness/conventions.md#mr-011), [`MR-004`](../../harness/conventions/MR-004-regelwerk-vendored.md), [`MR-005`](../../harness/conventions/MR-005-gate-haltung-docs-check.md); keine `LH-*`-/`ADR-*`-Änderung (der Slice nennt keine berührte Spec-Stelle).
- Artefakte: `repo.mk`, `Makefile`, `d-check.mk`, `harness/mk/*`, `harness/sensors/*`, `.d-check.yml`, `harness/README.md`, Adaptionen; Commits `eafe636..792c53e`.

DoD-Abgleich:
- [x] Targets in `repo.mk`; `make gates` fährt `lint`, `test`, `coverage-gate`, `docs-check`, `baseline-verify`, endet mit `record-gates`; `ci`/`fullbuild` verfügbar — Evidence: `make -n -j8 gates` zeigt `baseline-verify.sh`, d-check-Lauf, drei `docker build` (lint/test/coverage) und als letzte Zeile `bash tools/harness/record-gates.sh`; `GATE_CHECKS += lint test coverage-gate` und `ci`/`fullbuild` stehen in `repo.mk`. Der echte Lauf endete mit Exit 0, `coverage-gate: OK — coverage 90.90% meets threshold 90%`.
- [x] Werkzeug gelaufen; `harness/README.md` führt „Targets der Werkzeug-Fragmente“; `.d-check.yml` `targets` mit `makefiles`, `doc-tables`, `authority` als Listen und `authority-disjoint: true`; `make docs-check` grün — Evidence: `.d-check.yml` Z. 123–128; `d-check: 232 Datei(en) geprüft, 0 Befund(e)`.
- [x] Sensor-Dateien `harness/sensors/coverage-gate.md` und `verify-depguard.md`, in den Target-Zellen verlinkt — Evidence: `harness/README.md` §Sensors. Abweichung vom Plan §3 (siehe unten): `test-docker` bekommt keine Sensor-Datei; Plan §3 nennt das selbst so.
- [x] `make gates` grün — Evidence: eigener Lauf, Exit 0; Stempel `.harness/state/gates-passed.head` = `792c53e`.
- [x] Review durchgeführt, Report liegt vor — Evidence: [`Report`](2026-10-16-slice-harness-ai-harness-init-adoption.md); Findings siehe unten.
- [x] Doku-Update: `MR-011` nach `conventions/done/` (Commit `20893ce`, reiner Move), `MR-012` neu, `MR-004`/`MR-005` auf dem Werkzeug-Stand (Gate-Image `v0.86.1`, neun `exempt-targets`) — Evidence: Diff und `conventions.md`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag — offen: Slice §7 trägt „Offen.“; folgt nach dieser Verifikation.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang — offen (Risiken ohne Ausgang, siehe unten); folgt mit der Closure.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | Exit 0, Coverage 90.90 % ≥ 90 % |
| `make docs-check` | pass | 232 Dateien, 0 Befunde |
| `tools/harness/fetch-baseline-cache.sh --verify` | pass | „verify ok (54 Dateien, vollständig)“ |
| `make baseline-verify` | pass | „v6.18.1 OK — 54 Dateien“ |
| `make doc-planning` | pass | 0 Befunde |
| `make doc-targets` | pass | 0 Befunde |
| `make -n -j8 gates` | pass | `record-gates.sh` ist die letzte Zeile, nach allen fünf Checks |

Traceability:
| ID | Beleg |
| --- | --- |
| [`MR-012`](../../harness/conventions.md#mr-012) | `repo.mk`, Index-Zweiteilung (`doc-targets` 0 Befunde, `authority-disjoint`), Ordnungskante (`make -n -j8 gates`) |
| [`MR-005`](../../harness/conventions/MR-005-gate-haltung-docs-check.md) | `d-check.mk` (digest-gepinnt `v0.86.1`), `make docs-check` grün, Zahl der `exempt-targets` = 9 |
| [`MR-004`](../../harness/conventions/MR-004-regelwerk-vendored.md) | `--verify` und `make baseline-verify` grün |

Carveouts:
- Neu: none
- Gelöst: none (`MR-011` ist eine Adaption, kein Carveout)
- Unverändert: alle bestehenden

Nicht ausgeführt:
- `make test-docker` — kein Docker-Pfad berührt; das Rezept ist byte-gleich nach `repo.mk` übernommen.
- `make verify-depguard` — Importregeln und `.golangci.yml` unverändert; Rezept byte-gleich übernommen.
- `make ci` / `make fullbuild` — Closure-Sensor für Security-/Release-Pfade; nicht berührt, nur das Vorhandensein der Targets wurde geprüft.

Commit / Artefakt:
- `792c53e` (HEAD); Vorstufe der Kette `eafe636`, `864cb1d`, `20893ce`, `4a81edb`.

## Plan-vs-Code-Diff (Plan §3)

| Plan-Zeile | Ist | Urteil |
|---|---|---|
| `repo.mk` neu vor dem Lauf | `eafe636` vor `864cb1d` | bestätigt |
| `Makefile`, `d-check.mk`, `harness/mk/*.mk`, `tools/harness/*`, `.claude/*` vom Werkzeug | Commit `864cb1d` | bestätigt |
| `harness/mk/ai-harness-init.md` | vorhanden | bestätigt |
| `.d-check.yml` update | `targets`-Block wie geplant | bestätigt |
| `harness/README.md`: Werkzeug-Zeile; `record-gates`/`slice-mv`/`hooks-install` bleiben mit Bindung `MR-012` | so im Ist | bestätigt | <!-- d-check:ignore (Make-Target, keine Slice-Kennung) -->
| `harness/sensors/*.md` für `coverage-gate`, `verify-depguard`; `test-docker` ohne Datei | genau diese zwei Dateien | bestätigt |
| `MR-011` `git mv`, `MR-004`, `MR-005`, `conventions.md` | `20893ce` reiner Move, Inhalt in `4a81edb` | bestätigt |

Abweichung ohne Auswirkung auf die DoD: Der Review-Befund F-4 (Sensor-Datei `test-docker` fehlt) bezieht sich auf eine Lesart des Plans, die der Plan-Text (`test-docker` ohne Sensor-Datei) nicht stützt. Die Ist-Lage entspricht dem Plan.

## Prüfung (a): Target-Migration

Alle Targets des alten `Makefile` (`git show eafe636:Makefile`) sind erreichbar. Der Vergleich der Zielnamen gegen `make help` ergab genau einen Unterschied: `coverage` (Alias ohne `##`-Text; steht in `repo.mk` als `coverage: coverage-gate` und in `exempt-targets`). `packages`, `build-binaries`, `clean`, `deps`, `compile`, `run`, `build`, `lint`, `test`, `test-docker`, `coverage-gate`, `govulncheck`, `image-scan`, `verify-depguard`, `ci`, `fullbuild` stehen in `repo.mk`. Das Risiko „vergessenes Target“ aus Slice §6 ist damit sachlich widerlegt; sein Ausgang ist im Slice noch nicht eingetragen.

## Prüfung (b): Review-Befunde

| Finding | Befund im Ist |
|---|---|
| F-1 | behoben — `docs/maintainer/releasing.md` Z. 61/84 und `quality.md` Z. 202/233 nennen `repo.mk`; kein `Makefile` mehr mit Bump-Anweisung in `docs/maintainer/` |
| F-2 | behoben — `ci.yml`, `publish.yml` und `repo.mk` nennen `repo.mk`; README-Baum nennt `Makefile` als generierten Aggregator |
| F-3 | behoben — `plan-welle.md`, `close-welle.md`, `planner.md` tragen den Hinweis „Ungenutzt, bis `slice-harness-wellen-betrieb` …“ | <!-- d-check:ignore (Slice-Kennung im Zeitdokument) -->
| F-4 | Plan und Ist stimmen überein (siehe oben); keine Änderung nötig |
| F-5 | behoben — `coverage-gate.md` Exit 2 nennt fehlende Eingabe, leere Eingabe, fehlende `total:`-Zeile; `verify-depguard.md` Exit 2 nennt den unsauberen Arbeitsbaum |
| F-6 | behoben — `repo.mk` sagt „nine“ (Kommentar und Hilfetext); `MR-005` sagt „Neun“; `exempt-targets` hat 9 Einträge |
| F-7 | behoben — `MR-005` Gate-Image-Zeile nennt den Stand `v0.86.1` ohne „Bump 2026-10-16“ und ohne „nicht `doc-check`“ |
| F-8 | bewusst offen — `record-gates`/`slice-mv`/`hooks-install` bleiben als Repo-Zeilen mit Handkopie des Vertragstexts; Bindung `MR-012` benennt das | <!-- d-check:ignore (Make-Target, keine Slice-Kennung) -->

F-9 und F-10 (INFO) sind unverändert und ohne Gate-Wirkung.

## Prüfung (c): offene Punkte

Slice §6 trägt vier Risiken, alle ohne Ausgang: Target-Verlust (Vergleich liegt jetzt vor, siehe (a)), `d-check`-Sprung `v0.51.1` auf `v0.86.1` (`make docs-check` liefert 0 Befunde auf 232 Dateien; der Ausgang ist trotzdem noch einzutragen), Wellen-Dateien gegen `MR-003` (Ausgang schon benannt: `slice-harness-wellen-betrieb`) und Träger-Binary (`make traeger-fetch` pro Klon; der Binary fehlt im Klon, `.harness/state/bin/` ist vorhanden). Das ist die erwartete Lage vor der Closure und keine DoD-Verletzung. <!-- d-check:ignore (Slice-Kennung im Zeitdokument) -->

## Urteil

DoD-Punkte 1 bis 6 bestätigt. Punkte 7 (Closure-Notiz) und 8 (Risiko-Ausgänge) bewusst offen und vor `git mv` nach `done/` zu schließen. Keine Abweichung gefunden, die die Closure blockiert; F-8 bleibt als bewusst offener LOW-Befund.
