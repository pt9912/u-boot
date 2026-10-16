# Review-Report: slice-harness-wellen-betrieb — 2026-10-16

**Review-Art:** Code (Doku-/Harness-Diff) — geprüft gegen Plan, Entscheidungen (`MR-003`, `MR-012`) und Hard Rules.

**Gegenstand:** Diff `67aa428..HEAD`, insbesondere die Commits `9ed60be` und `6be822c`; Plan: [`slice-harness-wellen-betrieb`](../plan/planning/in-progress/slice-harness-wellen-betrieb.md).

**Skill:** `.harness/skills/reviewer.md` @ `6be822c` · <!-- d-check:ignore (Skill-Pfad, Stand des Laufs) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-16

**Eingangs-Kontext:**

- [`slice-harness-wellen-betrieb`](../plan/planning/in-progress/slice-harness-wellen-betrieb.md) (Scope, DoD, Nicht-Ziele); Entscheidung des Projektinhabers: Welle-Dateien und Register ja, keine Archivierung, keine Results-Datei
- `MR-003`, `MR-012` (`harness/conventions/`), `harness/conventions.md` (ID-Schema, Modus-Deklaration)
- `v6.18.1` · `regelwerk/modul-06-roadmap.md` §Roadmap-Struktur: fünf Abschnitte; `v6.18.1` · `templates/docs/plan/planning/roadmap.template.md`
- `AGENTS.md` (Hard Rules 3.1, 3.3, 3.6), `harness/roles.md`

Sensoren gelaufen (Docker): `make docs-check` (237 Dateien, 0 Befunde), `make doc-planning` (0 Befunde), `make doc-targets` (0 Befunde). Zusatznachweis Marker-Hälfte in einem Wegwerf-Klon (HEAD, außerhalb des Repos, anschließend gelöscht), d-check-Modul `planning` allein:

| Zustand des Klons | Ergebnis |
|---|---|
| Ist: Slice in `in-progress/`, kein Marker | 0 Befunde |
| `in-progress/` ohne Slice, kein Marker | 1 Befund `planning-drift` („trägt den Ruhe-Marker nicht — er gehört dorthin“) |
| `in-progress/` ohne Slice, Marker `Nichts in Arbeit.` im Block `## Offene Wellen` | 0 Befunde |
| Slice in `in-progress/`, Marker steht | 1 Befund `planning-drift` („die Sektion muss die Arbeit benennen“) |

Die Marker-Hälfte greift in beide Richtungen; `.d-check.yml` (`heading`, `marker`) trifft den Baseline-Wortlaut. Der Marker steht im Repo nur in `.d-check.yml`, nicht literal in Regel- oder Hinweiszeilen der Roadmap.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Der Abschnitt „Offene Wellen“ nennt als einzigen Slice in Arbeit `slice-v2-homebrew-formula` („ein Slice ohne Welle“); in `in-progress/` liegt zugleich `slice-harness-wellen-betrieb`. Der Homebrew-Slice trägt kein `Welle:`-Feld, die Aussage „ohne Welle“ hat keinen Beleg im Slice. | Maintainability; `v6.18.1` · `regelwerk/modul-06-roadmap.md` §Roadmap-Struktur: fünf Abschnitte (Bullet Offene Wellen: derivativ) | `docs/plan/planning/in-progress/roadmap.md` · „In Arbeit ist ein Slice ohne Welle“ | nein — Review-only (der Sensor prüft nur die Marker-Hälfte) | Derivativer Roadmap-Abschnitt gegen Verzeichnisbestand nicht abgeglichen | <!-- d-check:ignore (Kennungen im Befund-Zitat des eingefrorenen Reports) -->
| F-2 | LOW | Das Register-README ist im Diff unverändert, obwohl die DoD „an u-boots Kürzel angepasst“ als erledigt führt; es sagt weiter „ohne Wellen-Betrieb löst die jeweilige Slice-Closure diesen Schritt selbst aus“. Die Schreibweise `KUERZEL` (README, `close-welle`) weicht von `KÜRZEL` (ID-Schema in `harness/conventions.md`) ab. | Maintainability; `MR-003` | `docs/plan/planning/observations/README.md` · „ohne Wellen-Betrieb“; `harness/conventions.md` · `BEO-<KÜRZEL>/<slug>` | nein — Review-only | Zustandsaussage und Schreibweise zwischen Dokumenten nicht nachgezogen |
| F-3 | LOW | Die Roadmap-Notiz „Closure-Form (Abweichung, MR-003)“ sagt, die Welle-Closure lebe „vollständig“ im Release-Cut-Slice; `MR-003` und `close-welle` legen für neue Wellen die Closure-Notiz in §7 der Welle-Datei und nennen den Release-Cut-Slice als Detailquelle. | `MR-003` (Adaption (b)) | `docs/plan/planning/in-progress/roadmap.md` · „Closure-Form (Abweichung“ | nein — Review-only | Zwei Formulierungen derselben Zusage |
| F-4 | LOW | In `close-welle` tragen die Listenpunkte die Nummern 2–6 und 8 (7 fehlt), die Überschriften „Schritt 1–5“ und „fünf Schritte“ zählen anders; die Querverweise („Schritt 3“, „Schritt 4“, „Schritt 5“) stimmen untereinander. | Maintainability | `.claude/commands/close-welle.md` · „8. `make gates` grün nach dem Commit“ | nein — Review-only | Doppelte Nummerierung (Liste und Schrittname) |
| F-5 | INFO | `harness/mk/archivierung.mk` und `harness/mk/ai-harness-init.md` führen `make archive-welle` weiter; „bleibt ungenutzt“ ist damit eine Aufruf-Disziplin, kein Wächter. Der Träger ist gitignoriert; mit lokal vorhandenem Träger bewegt, löscht und committet ein ausdrücklicher Aufruf. Mit `MR-012` Punkt 4 und dessen Auflösungs-Trigger vereinbar. | `MR-012` Punkt 4 | `harness/mk/archivierung.mk` · „archive-welle:“ | nein — Review-only | Ungenutztes Werkzeug-Ziel ohne Sperre |
| F-6 | INFO | Die Zeiger-Hälfte von „Offene Wellen“ (ein Zeiger je flacher Welle-Datei) hat keinen Sensor; `.d-check.yml` benennt die Lücke ausdrücklich. Zulässig nach Baseline, solange benannt. | `v6.18.1` · `regelwerk/modul-06-roadmap.md` §Roadmap-Struktur: fünf Abschnitte | `.d-check.yml` · „hat keinen Sensor“ | nein — Review-only | Ungewächterte Hälfte, benannt |
| F-7 | INFO | In der Tabelle „Abgeschlossene Wellen“ steht eine Leerzeile nach der Zeile `welle-harness-konformitaet-nachlauf`; die folgenden Zeilen bilden beim Rendern eine zweite Tabelle. Nicht Teil des Diffs. | Maintainability | `docs/plan/planning/in-progress/roadmap.md` · Zeile vor `welle-spec-technik-stratum` | nein — Review-only | Bestandsmangel außerhalb Diff-Scope |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Roadmap-Struktur gegen `roadmap.template.md` und Modul 6 (fünf Abschnitte plus Abhängigkeitsgraph, „Offene Wellen“ derivativ, Marker nicht literal in Regel-/Hinweiszeilen, keine toten Verweise auf „Aktuelle Welle“ — Repo-weit kein Treffer außerhalb Baseline/Archiv) | geprüft, ohne Befund (außer F-1, F-3) |
| `.d-check.yml` Modul `planning` (`heading`, `marker`), `make doc-planning`, Wegwerf-Klon-Nachweis in beide Richtungen | geprüft, ohne Befund |
| `MR-003` Aussagen (a) Ort, (b) keine Results-Datei, (c) keine Archivierung; Index-Zeile in `harness/conventions.md`; Bezug zu Lastenheft 0.6.0 (Planungs-Anforderungen entfallen, kein Abwärtsverweis) | geprüft, ohne Befund |
| `docs/plan/planning/README.md` (Slices vs. Wellen, drei Dateinamen-Formate, Master-Dokumente, Register) | geprüft, ohne Befund |
| Kürzel-Spalte in `harness/conventions.md`: alle neun Sub-Areas belegt (DOK, DOM, APP, PORT, CLI, ADP, E2E, WZG, PKG), eindeutig, `cmd/uboot` mit Strich, Erläuterung unter der Überschrift; ID-Schema um `welle-<slug>` und `BEO-<KÜRZEL>/<slug>` ergänzt | geprüft, ohne Befund (Schreibweise siehe F-2) |
| `plan-welle`, `close-welle`, `planner`: Archivierungs-/Results-Schritte entfernt (grep `archiv`, `results`, `archive-welle`, `ANPASSEN`, `ungenutzt`: nur die beabsichtigten Negativ-Aussagen), kein „sechs“ außer im Erfassungs-Satz, Adaptionsblock u-boot-richtig, Move-/Inhalt-Trennung (Hard Rule 3.3) und Docker-only (3.1) gewahrt | geprüft, ohne Befund (außer F-4) |
| Erfassungs-Satz in allen sechs `.claude/agents/*.md` | geprüft, ohne Befund (wortgleich) |
| `MR-012` Punkt 4 und Auflösungs-Trigger konsistent zu `MR-003` | geprüft, ohne Befund |
| Keine Lockerung von Gates, Schwellen oder Importregeln (Hard Rule 3.6); Modul `planning` bleibt aktiv, `modules` unverändert | geprüft, ohne Befund |
| Go-Code, `make test`, `make test-docker`, `make verify-depguard`, `make gates` | nicht geprüft — Diff berührt keinen Code, keine Docker-Pfade, keine Importregeln; `make gates` ist Sache des Implementers/Verifiers |
| DoD-Closure, Closure-Notiz, Risiko-Ausgänge in §6 des Slice | nicht geprüft — Verifier-Aufgabe |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 4 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Derivativer Roadmap-Abschnitt gegen Verzeichnisbestand nicht abgeglichen · Zustandsaussage und Schreibweise zwischen Dokumenten nicht nachgezogen · Zwei Formulierungen derselben Zusage · Doppelte Nummerierung

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM.

**Übergabe:** Findings F-1 bis F-4 (LOW) gehen an den Implementer, F-5 bis F-7 sind Information. Der Report ersetzt keine Verifikation — DoD-Konformität prüft der Verifier.
