# Verification Evidence: slice-v2-homebrew-smoke-haertung — 2026-10-17

**Rolle:** Verifier (frischer Kontext, kein Self-Review). Frage: Wurde der Slice richtig gebaut, gemessen an Plan, DoD und Sensors? Form nach `harness/verification.md`. <!-- d-check:ignore (Pfad als Zitat) -->

**Gegenstand:** Slice [`slice-v2-homebrew-smoke-haertung`](../plan/planning/done/slice-v2-homebrew-smoke-haertung.md), Commits `8b5bd3b`, `ddf2813`, `82677f1`, `88a7259`, `8877fa3` (ohne Welle). Maßstab: §1 (Ziel, Nicht-Ziele), §2 (DoD), §3 (Plan), §6 (Risiko), [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md), [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung), [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache). Review-Report: [`2026-10-17-slice-v2-homebrew-smoke-haertung.md`](2026-10-17-slice-v2-homebrew-smoke-haertung.md).

## Verification Evidence

Scope:
- Slice: [`slice-v2-homebrew-smoke-haertung`](../plan/planning/done/slice-v2-homebrew-smoke-haertung.md)
- IDs: [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung), [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache), [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md)
- Artefakte: `packaging/homebrew-tap/.github/workflows/smoke.yml`, `scripts/homebrew-formula.rb.tmpl`, `cmd/uboot/homebrew_scripts_test.go`, `.github/workflows/publish.yml`, `packaging/nfpm.yaml`, `docs/maintainer/releasing.md`

DoD-Abgleich (§2 des Slice):
- [ ] Tap-Smoke (Build-from-source bzw. gleichwertig, `u-boot add`, gatendes `brew audit --strict --online`, Beleg grüner Smoke-Lauf im Tap) — **NICHT erfüllt.** Die Vorlage im Repo ist umgesetzt (siehe Diff), aber der Beleg fehlt: `gh api` (nur lesend) zeigt im Tap `pt9912/homebrew-u-boot` HEAD `a8498fd` und eine `smoke.yml`, die der alten Fassung entspricht (kein `--build-from-source`, kein `brew test`, kein `add postgres`, Audit mit `|| true`; `diff` gegen die Vorlage zeigt genau die Slice-Änderungen). Es gab keinen Lauf der neuen Fassung, weder auf `macos-14` noch auf `macos-15-intel`, und `brew audit --strict --online` lief nie gatend. Akzeptanzkriterium 4 ist damit nicht belegt.
- [x] F-5: Formel-`desc` Englisch/ASCII, Test angepasst — Evidence: Diff `ddf2813`; eigene Mutation (siehe Sensors) färbt `make test` rot.
- [x] F-8: Kommentar und Stabilitäts-Bestimmung — Evidence: Diff `82677f1`/`8877fa3` (`outputs.is_stable`, `needs.publish.outputs.is_stable == 'true'`, Kommentar „six binaries and the four Linux packages“ am passenden Step). Verdrahtung gelesen, nicht gelaufen (kein Workflow-Sensor, siehe F-1).
- [x] `make gates` grün — selbst gefahren, Exit 0.
- [x] Review durchgeführt, Report liegt vor (0 HIGH, 2 MEDIUM, 1 LOW, 3 INFO); F-2 und F-3 in `8877fa3` behoben.
- [x] Doku-Update — `docs/maintainer/releasing.md` §6 gegen Smoke geprüft, stimmt (Install mit `--build-from-source`, `brew test`, `--version`, `init`, `add postgres`, Audit gatet).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag — offen (§7 steht auf „Offen.“); gehört in die Closure.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang — offen: das Audit-Risiko hat noch keinen Ausgang, weil der erste gatende Audit nie lief.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | selbst gefahren, Exit 0; `coverage-gate: OK — coverage 91.10% meets threshold 90%`; `d-check: 281 Datei(en) geprüft, 0 Befund(e)`; `baseline-verify: v6.18.1 OK` |
| `make doc-trace` | pass | selbst gefahren: `83 Anforderung(en), 0 Waise(n)` |
| Eigene Mutation `make test` | pass (Test rot) | `desc` auf Nicht-ASCII-Zeichen gesetzt: `FAIL: TestHomebrewFormulaDesc_EnglishASCII ... holds the non-ASCII rune`, Exit 2; danach `git checkout -- scripts/homebrew-formula.rb.tmpl`, `git status` sauber |
| `make docs-check` | pass | Teil von `make gates`; nach Anlage dieser Evidence erneut gefahren (siehe Übergabe) |
| `make test-docker` | not run | kein Docker-Laufzeitpfad berührt (Diff: Workflows, Skript-Template, Test, nfpm, Doku) |
| `make verify-depguard` | not run | keine Importregel berührt |
| Tap-Smoke (extern) | **nicht gelaufen** | neue `smoke.yml` nicht ausgerollt; Host-`brew` und Tap-CI sind für den Verifier gesperrt, Push ins Tap braucht die Freigabe des Projektinhabers |
| `brew audit --strict --online` | **nicht gelaufen** | nur im Tap-CI möglich |

Traceability:
| ID | Beleg |
| --- | --- |
| [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung) | Vorlage `smoke.yml`, Formel-`desc`, Maintainer-Doku; **Smoke-Lauf der gehärteten Fassung steht aus** |
| [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache) | `TestHomebrewFormulaDesc_EnglishASCII` (mutationsgeprüft), englische nfpm-Beschreibung |
| [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) | im Diff unverändert (Accepted-immutable); Tap-Nachzug-Mechanik unverändert |

Plan-vs-Code-Diff (§3 gegen `9ec0f12..HEAD`, nur Produktdateien):
| Plan-Zeile | Diff | Befund |
|---|---|---|
| `packaging/homebrew-tap/.github/workflows/smoke.yml` update | geändert | konform |
| Tap-Repo `smoke.yml` update | **nicht geändert** (Tap = alte Fassung, HEAD `a8498fd`) | **offen**, Closure-Voraussetzung |
| `scripts/homebrew-formula.rb.tmpl` update | geändert (`desc`) | konform |
| `cmd/uboot/homebrew_scripts_test.go` update | geändert (neuer Test) | konform; Plan nennt zusätzlich „Akzeptanzkriterium 4“ als Test-Zweck, der Test deckt nur die `desc` (Kriterium 4 ist Workflow-Sache, nur im Tap-Lauf prüfbar) |
| `.github/workflows/publish.yml` update | geändert | konform |
| nicht im Plan: `packaging/nfpm.yaml`, `docs/maintainer/releasing.md` | geändert (`8877fa3`, `88a7259`) | Abweichung vom Plan, vertretbar: nfpm folgt aus Review F-2 ([`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache), §3.9-Suche nach altem Wortlaut), Doku aus DoD-Punkt „Doku-Update“. Beide sind in §3 nachzutragen, damit Plan und Diff stimmen. |
| Nicht-Ziele §1 (ADR, `done/`-Slice, Welle, Roadmap, Carveout, neue Plattformen) | nicht berührt; `8877fa3` legt vier fremde Slice-Dateien in `open/` und eine Welle-Datei an (Commit-Bündelung, fremder Umfang, siehe Auflage 6) | Slice-Scope selbst eingehalten |

Externes Tap-Repo: unverändert gegenüber dem Stand vor dem Slice (lesend per `gh api`; HEAD `a8498fd`, `smoke.yml` = alte Fassung).

Carveouts:
- Neu: none
- Gelöst: none
- Unverändert: [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung)-Zeile in `in-progress/carveouts.md` verweist auf diesen Slice als offene Härtung.

## Bewertung der Review-Findings

| Finding | Bewertung |
|---|---|
| F-1 (`publish.yml` ohne Sensor) | Grenze genügt für den Slice, aber nur mit Benennung. Verdrahtung korrekt gelesen (`outputs`, Step-Id `tag`, `needs`); Verhalten für gültige Tags gleich der alten Bedingung. Das Restrisiko (stilles Überspringen von `tap` bei Umbenennung) ist von keinem Sensor und von keinem Pre-Release-Tag abgedeckt; belegt erst der nächste stabile Tag. Auflage: in §7 als Grenze benennen (kein Workflow-Linter im Repo) und beim nächsten stabilen Release prüfen, dass `tap` läuft. Nicht blockierend. |
| F-4 (Beleg Smoke-Lauf steht aus) | **Blockierend für die Closure.** Siehe Verdikt. |
| F-5 (`--build-from-source` bei Binary-Formel) | Akzeptanzkriterium 4 („bzw. gleichwertig“) ist **dem Wortlaut nach erfüllt** (der Schritt steht wörtlich im Smoke), **inhaltlich** trägt `brew test u-boot` den Gegenwert: er führt den `test do`-Block der installierten Formel aus. Es gibt keinen Quellbau, weil die Formel keine Quelle hat; der Kommentar im Workflow und das Commit sagen das ehrlich. Bedingung: Die Closure-Notiz hält diese Setzung ausdrücklich fest. Erfüllung gilt erst mit grünem Lauf. |
| F-6 (Denylist-Test `desc`) | Test ist nicht tautologisch (eigene Mutation rot). Die Denylist (`fuer`, ` und `, `reproduzierbar`) plus Pflicht-Substring `reproducible` plus ASCII-Prüfung ist eine Grenze: ein deutsches Wort ohne Umlaut außerhalb der Liste bleibt unentdeckt. Für eine einzeilige `desc` akzeptabel; nicht blockierend. |

## Was bis zur Closure fehlt

1. Tap-Rollout: `smoke.yml` der Vorlage ins Tap `pt9912/homebrew-u-boot` bringen (per `scripts/tap-nachzug.sh`; Freigabe des Projektinhabers nötig, Push in externes Repo).
2. Grüner Lauf der neuen Fassung auf `macos-14` und `macos-15-intel` (Lauf-ID und Jobs in §7 nennen); alle Schritte, besonders `brew install --build-from-source`, `brew test`, `add postgres` (Managed-Block-`grep`).
3. Ergebnis des gatenden `brew audit --strict --online pt9912/u-boot/u-boot` (Befund-Ausgabe lesen, nicht nur Statusfarbe).
4. Closure-Notiz §7 mit Steering-Loop-Lerneintrag, Setzung zu `--build-from-source` (F-5) und Grenze F-1; §6-Ausgang für das Audit-Risiko.
5. §3 um `packaging/nfpm.yaml` und `docs/maintainer/releasing.md` ergänzen.

**Roter Lauf (echte Audit-Befunde):** Der Audit ist neu gatend, und ein roter Lauf ist ein möglicher, kein unwahrscheinlicher Ausgang. Behandlung: (a) Befund und Lauf-ID in §6 als Ausgang des Risikos festhalten („eingetreten: <Befund>“). (b) Ist die Behebung klein (Formel-Template, ein Nachzug), im Slice beheben und den Lauf wiederholen. (c) Ist sie eigener Umfang, greift die vorab benannte Rückführung §4 `in-progress` → `open` (blockiert) mit Verweis auf den Befund und einem Folge-Slice für die Formel. Das Gate darf dabei nicht wieder auf `|| true` zurückgenommen werden (AGENTS.md §3.6; eine Lockerung bräuchte Plan-/ADR-Anker). Ein befristetes Aussetzen nur einzelner Audit-Regeln wäre ein Carveout mit Trigger und Folge-Slice. Auch bei einem Befund, der erst durch den Rollout sichtbar wird, bleibt DoD-Punkt 1 offen, bis der Lauf grün ist.

## Verdikt

**Bestanden mit Auflagen für die Closure — der Slice ist ohne den Tap-Lauf nicht vollständig und noch nicht schließbar.** Die Repo-seitige Lieferung ist konform, die Sensors (`make gates`, `make doc-trace`, Mutation) sind grün. Das erste DoD-Kriterium (grüner Smoke-Lauf im Tap, Akzeptanzkriterium 4) und die Closure-Punkte aus „Was bis zur Closure fehlt“ sind offen. Der Slice bleibt `in-progress`.

Auflagen:
1. Tap-Rollout und grüner Lauf auf beiden Runnern; Audit-Ergebnis belegen.
2. Closure-Notiz mit Setzungen (F-5, F-1) und Steering-Loop-Eintrag; §6-Ausgang eintragen; §3 nachtragen.
3. Bei rotem Lauf: Behandlung wie oben, kein Zurücknehmen des Gates.
4. Commit-Hygiene: `8877fa3` bündelt fremde Planungsdateien (vier Slices in `open/`, eine Welle) mit dem Review-Fix. Ohne Wirkung auf die DoD, aber Traceability-Hinweis: Umfang gehört einem eigenen Commit.

Nicht ausgeführt: `make test-docker`, `make verify-depguard` (kein berührter Pfad); Tap-CI-Lauf und `brew audit` (Docker-only, Host-`brew` gesperrt, externer Push nicht freigegeben).
