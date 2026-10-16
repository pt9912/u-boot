# Verification Evidence: slice-cr-cli-quiet-fehler-und-warnungen — 2026-10-16

**Rolle:** Verifier (frischer Kontext, kein Self-Review). Frage: Wurde der CR richtig gebaut, gemessen an Plan, DoD und Sensors? Form nach `harness/verification.md`. <!-- d-check:ignore (Pfad als Zitat) -->

**Gegenstand:** CR-Slice [`slice-cr-cli-quiet-fehler-und-warnungen`](../plan/planning/done/slice-cr-cli-quiet-fehler-und-warnungen.md) (Delivery `b5fa4b4`, Review-Fix `5a20103`), Welle [`welle-cli-verbosity`](../plan/planning/welle-cli-verbosity.md) §3 bis §5 (CR zuerst). Maßstab: Slice §1 (CR-Kern, Nicht-Ziele), §2 (DoD), §3, §6. Review-Report: [`2026-10-16-slice-cr-cli-quiet-fehler-und-warnungen.md`](2026-10-16-slice-cr-cli-quiet-fehler-und-warnungen.md).

## Verification Evidence

Scope:
- Slice: [`slice-cr-cli-quiet-fehler-und-warnungen`](../plan/planning/done/slice-cr-cli-quiet-fehler-und-warnungen.md)
- IDs: [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Änderung), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)
- Artefakte: `spec/lastenheft.md` (Versionsfeld, `LH-FA-CLI-005`, §16), `CHANGELOG.md` `[Unreleased]`; kein Go-Code im Diff. <!-- d-check:ignore (Pfade als Zitat) -->

DoD-Abgleich (Slice §2):
- [x] Lastenheft 0.8.0 — Evidence: Versionsfeld `0.8.0`; Ausnahme-Absatz in `LH-FA-CLI-005` vor dem Mehrfachangabe-Satz (`spec/lastenheft.md` Z. 189); §16-Zeile 0.8.0 mit Verweis „Vereinbarung mit dem Projektinhaber"; die vier Stufenzeilen und der Mehrfachangabe-Satz unverändert (`git show` beider Commits). <!-- d-check:ignore (Pfad als Zitat) -->
- [x] Umsetzungsbedarf einem Slice zugeordnet — Evidence: [`slice-v1-cli-quiet-nur-fehler`](../plan/planning/in-progress/slice-v1-cli-quiet-nur-fehler.md) liegt in `open/` und trägt Abnahmekriterium 1 (Tests je Subkommando; Lockerungs-Hinweis bei `init`/`generate`/`up`) in seiner DoD und §3 (Testzeile), Kriterium 2 (`docs/user/benutzerhandbuch.md`, Zeilen um 299 und 1280) in DoD und §3, Kriterium 3 (CHANGELOG) in DoD und §3. <!-- d-check:ignore (Pfad als Zitat) -->
- [x] Doku-Update CHANGELOG — Evidence: `[Unreleased]` `Changed`-Eintrag mit Link auf `LH-FA-CLI-005` (`b5fa4b4`). Präsens-Hinweis siehe Review-Findings (F-5). <!-- d-check:ignore (Kennung als Zitat) -->
- [x] `make gates` grün — siehe Sensors (selbst gefahren).
- [x] Review-Report liegt vor — Evidence: `docs/reviews/2026-10-16-slice-cr-cli-quiet-fehler-und-warnungen.md` (0 HIGH, 1 MEDIUM, 4 INFO), Reviewer-Kontext getrennt vom Autor. <!-- d-check:ignore (Pfad als Zitat) -->
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag — offen: Slice §7 trägt „Wird bei der Closure gefüllt." Folgt nach diesem Verdikt (Auflage 1).
- [ ] Beobachtungs-Register — offen: §8 notiert „Register leer" für die Sichtung; die Fortschreibung bzw. der „keine Beobachtung"-Vermerk in §7 steht mit der Closure-Notiz aus (Auflage 1).
- [~] Jedes Risiko aus §6 trägt einen Ausgang — beide Risiken tragen „entfallen" mit Begründung; die Begründung des ersten verschiebt die Code-Grenze in den Umsetzungs-Slice, der sie trägt (siehe F-2).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) — Folge-Slice getragen (`slice-v1-cli-quiet-nur-fehler` in `open/`, in der Welle §4 gelistet); Anker und Register entstehen mit der Closure-Notiz §7 (Auflage 1). <!-- d-check:ignore (Kennung als Zitat) -->

CR-Kern gegen Lastenheft (Slice §1 wörtlich gegen `LH-FA-CLI-005`): <!-- d-check:ignore (Kennung als Zitat) -->
- Der Ausnahme-Absatz im CR-Kern und im Lastenheft ist nach `5a20103` zeichengleich (beide mit Klammer „zusätzliche Capabilities, `--privileged`, Seccomp-/AppArmor-Anpassungen, zusätzliche Devices", Zitat von `-006`, `-007`, `-008` und dem Satz „Alle anderen Hinweise entfallen unter `--quiet`."); der Rest des Anforderungstexts (Stufen, Mehrfachangabe) ebenfalls. Geprüft per Diff beider Commits und Lesen beider Texte. <!-- d-check:ignore (Zitat) -->
- Nicht nachgezogen (Drift im Plan, nicht im Vertrag): (a) der Ziel-Absatz in §1 („Ziel", Z. 39 bis 42) nennt als Quelle der Ausnahme nur `-006`/`-007`, nicht `-008`; (b) der Historie-Zeilen-Entwurf in §1 trägt den Vorschlagstext („`--quiet` bleibt „nur Fehler"; …") und weicht von der §16-Zeile des Lastenhefts ab (F-4 des Reviews; die §16-Zeile ist die geltende, sie nennt zusätzlich `-008` und „alle anderen Hinweise entfallen"). Beides ohne Wirkung auf den Vertrag; Auflage 2.
- Nicht-Ziele (Code, „letzte Option gewinnt", Spezifikation) — eingehalten: Diff berührt weder `internal/`, noch `spezifikation.md`, noch `architecture.md`, noch `docs/user/`.

Welle-Reihenfolge (`welle-cli-verbosity` §3 bis §5):
- CR zuerst — Evidence: CR liegt in `in-progress/`, beide Code-Slices in `open/`; §5 und `slice-v1-cli-quiet-nur-fehler` §4 binden deren Start an „CR in `done/`". Konsistent. <!-- d-check:ignore (Kennung als Zitat) -->
- Welle-Closure-Trigger (§3) betreffen diesen Slice nicht; der Spezifikations-Eintrag bleibt dem letzten Slice vorbehalten (Spezifikation im Diff unberührt, §6-Risiko „Zwischenstand der Spec" gehalten).
- `make doc-trace`: 83 Anforderungen, 0 Waisen; `LH-FA-CLI-005`-Zeile mit Slice-Belegen vorhanden (Welle-Trigger „keine Waise unter `LH-FA-CLI-005`" erfüllt). <!-- d-check:ignore (Kennung als Zitat) -->

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | selbst gefahren, Exit 0; `coverage-gate: OK — coverage 90.90% meets threshold 90%` (lint, test, docs-check, baseline-verify Teil von `gates`) |
| `make doc-trace` | pass | selbst gefahren, Exit 0; „83 Anforderung(en), 0 Waise(n)." |
| `make docs-check` | pass | Teil von `make gates`; nach Anlegen dieser Datei erneut gefahren (siehe Abschluss unten) |
| `make test-docker` | not run | kein Docker-Pfad berührt (Diff: Lastenheft, CHANGELOG, Plan, Review) |
| `make verify-depguard` | not run | keine Importregel berührt, kein Go-Code im Diff |

Traceability:
| ID | Beleg |
| --- | --- |
| [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | Lastenheft 0.8.0 (Absatz, §16), CHANGELOG-Eintrag; Umsetzungsbeleg (Tests) liegt im Slice `slice-v1-cli-quiet-nur-fehler`, hier noch nicht erbracht (nicht Ziel des CR) | <!-- d-check:ignore (Kennung als Zitat) -->
| [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) | als Quellen der Ausnahme im Absatz und in §16 zitiert; keine Änderung an den Anforderungen selbst |

Plan-vs-Code-Diff:
- Plan §3 `spec/lastenheft.md` (Anforderung, Versionsfeld, §16-Zeile) — Ist: alle drei Spuren gesetzt (`b5fa4b4`, verfeinert in `5a20103`).
- Plan §3 `CHANGELOG.md` — Ist: `[Unreleased]`-Eintrag, englisch.
- Nicht im Plan, im Diff: Wegfall der Zeile „Nichts in Arbeit." in `roadmap.md` (Commit `b5fa4b4`); vom Review ohne Befund geprüft (Zustandsfeld, keine Chronik). <!-- d-check:ignore (Pfad als Zitat) -->
- Plan §3 Reihenfolge „Lastenheft ändern → Umsetzung" — eingehalten (Code unberührt).
- Slice-Datei trägt nach `5a20103` die Fassung des CR-Kerns und des §6-Ausgangs der Review-Korrektur.

Carveouts:
- Neu: none
- Gelöst: none
- Unverändert: none berührt

Nicht ausgeführt:
- `make test-docker`, `make verify-depguard` — kein berührter Pfad (siehe Sensors).

Commit / Artefakt:
- Delivery `b5fa4b4`; Review-Fix und Report `5a20103`; Review-Report siehe oben.

## Review-Findings: Stand und Bewertung

Der Report entstand gegen `b5fa4b4` und führt F-1 als offen; er wurde im Commit `5a20103` mit der Korrektur abgelegt und nicht nachgeführt (Auflage 3).

| Finding | Stand | Bewertung für die Closure |
|---|---|---|
| F-1 (MEDIUM) Ausnahme ohne Abgrenzung | in `5a20103` behoben, am Ist-Stand gegengeprüft: der Absatz nennt die Arten der Lockerung, zitiert `-008` und schließt alle anderen Hinweise aus; §16-Zeile und CR-Kern/§6 nachgezogen | erledigt. Restlücke: Degradations-Fallbacks aus `-007` (Container ohne Daemon) sind weder genannt noch ausgeschlossen; sie fallen nach „alle anderen Hinweise" unter die Unterdrückung, sofern sie keine der genannten Arten sind. Entscheidung des Umsetzungs-Slice, nicht blockierend. |
| F-2 INFO Risiko-Ausgang verschiebt Grenze in den Code | durch die F-1-Korrektur größtenteils gegenstandslos (der Vertrag trägt jetzt die Grenze) | nicht blockierend |
| F-3 INFO 0.8.0 gegen 0.7.1 | offen, ohne Änderungspflicht | nicht blockierend: §16 kennt keine Schnittregel; das Minor passt zur Einengung/Erweiterung einer Zusage (vgl. 0.7.0) |
| F-4 INFO §16-Zeile weicht vom CR-Entwurf ab | offen | nicht blockierend: die geltende Zeile ist die im Lastenheft; der CR-Entwurf (§1, „Historie-Zeile") ist Plan-Stand und sollte beim Schließen mit der Ist-Zeile gleichgezogen oder als Entwurf gekennzeichnet werden (Auflage 2) |
| F-5 INFO CHANGELOG im Präsens | offen | nicht blockierend: der Satz „the implementation follows in a later release" deckt den Zwischenstand. Auflage 4: kein Release zwischen CR und den beiden Code-Slices, sonst widerspricht das Präsens dem Binary (`--quiet` hebt heute auf Warn, Mehrfachangabe nicht „letzte gewinnt"). |

## Verdikt

**Bestanden, mit Auflagen für die Closure** (kein Befund sperrt das Schließen; kein Defekt in Lastenheft oder CHANGELOG):

1. Closure-Notiz §7 mit Steering-Loop-Lerneintrag schreiben (Finding-Klassen des Reviews eintragen, Zähler beachten), Register-Vermerk („keine Beobachtung" oder Fortschreibung) und die Paarungen belegen; erst danach reiner `git mv` nach `done/` (AGENTS.md §3.3, Inhalt zuerst). DoD-Häkchen setzen.
2. CR §1 an den Ist-Stand angleichen: Ziel-Absatz (Quellen `-006`/`-007`/`-008`) und Historie-Zeilen-Entwurf auf die geltende §16-Zeile; geringe Drift, kein Vertragsfehler.
3. Review-Report ist gegen `b5fa4b4` geschrieben und nennt F-1 als offen; in §7 den Stand „F-1 in `5a20103` behoben" festhalten (der Report selbst wird nicht umgeschrieben).
4. Kein Release-Cut, bevor `slice-v1-cli-quiet-nur-fehler` und `slice-v1-cli-verbosity-letzte-option-gewinnt` geliefert sind (CHANGELOG-Präsens, F-5). Die Degradations-Fallbacks aus `-007` im Umsetzungs-Slice ausdrücklich entscheiden und im Test pinnen. <!-- d-check:ignore (Kennung als Zitat) -->
