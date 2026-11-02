# Slice V1: CRLF-Behandlung in `u-boot logs --json` — `splitLogLines` normalisiert Wagenrückläufe vor dem Split; die CO-Datei des temporären Carveouts löst sich auf

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle, Release-Zuordnung offen (kleiner Produktfix; es gibt
keine Closure-Bedingung, die von der DoD verschieden ist).

**Anlass.** Der CRLF-Edge-Case aus [`slice-v1-cli-json-dry-run-logs`](../done/slice-v1-cli-json-dry-run-logs.md)
(Pre-T8-Bestätigungsrunde LOW-2, 2026-06-07): `splitLogLines`
(`internal/adapter/driving/cli/logs.go`, Zeilen 309 bis 320, gelesen am
2026-11-02) splittet nur auf `\n`; bei einem hypothetischen
CRLF-terminierten Compose-Output trügen alle `lines[i]` ein trailing `\r`.
Der Fall ist dokumentiert als Vertrags-Grenze
([`docs/maintainer/cli-json-contract.md`](../../../../docs/maintainer/cli-json-contract.md) §6.8,
Absatz „CRLF-Edge-Case“) und wird auf Projektinhaber-Entscheidung 2026-11-02
als temporärer Carveout mit diesem Plan als Folge-Slice geführt (Anlage
im [`slice-harness-carveouts-baseline-form.md`](../next/slice-harness-carveouts-baseline-form.md),
Baseline-Ziel-Form `docs/plan/carveouts/`). Hier entsteht seine Auflösung,
wenn der Trigger eintritt.

**Berührte Spec-Stellen:** `—` — die `logs`-Vertrags-Doku
(`docs/maintainer/cli-json-contract.md` §6.8) trägt den Case; dort steht
bereits der Auflösungsweg („Falls je relevant: Folge-Slice ergänzt explizites
`\r\n`-Handling“). Die Zusage bleibt [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe)
(maschinenlesbare Ausgabe, Single-Envelope) und
[`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (Exit-Codes) unverändert.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner; Anlage auf Projektinhaber-Entscheidung
2026-11-02). **Datum:** 2026-11-02.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** `splitLogLines` behandelt CRLF-terminierten Compose-Output: keine
`data.lines[i]` der `--json`-Antwort trägt ein trailing `\r` — Wagenrücklauf
unmittelbar vor dem Zeilen-Trenner wird vor dem Split entfernt, der Rest der
Semantik (Split auf `\n`, Stripping des trailing-`\n`-Tokens) bleibt
unverändert. Damit löst sich der Carveout auf (per reiner `git mv` nach
`docs/plan/carveouts/done/`, Index aktualisiert, Vertrags-Doku setzt den
Ist-Stand).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Kein `--follow --json`-Streaming** (Ablehnung T0-(a), unverändert) —
  der Edge-Case betrifft nur den gepufferten `--json`-Weg; ein Streaming-
  Pfad ist eine eigene Zusage, kein Fix.
- **Keine Änderung an den Zeit-/Format-Flags** (`--timestamps`, `--since`,
  `--until`, `--no-log-prefix`) — ihre Semantik ist unabhängig vom
  Zeilentrenner und berührt den Edge-Case nicht.
- **Keine Änderung im Docker-/Compose-Adapter.** Die Quelle der Zeilen bleibt
  der Compose-Output; das Normalisieren sitzt am einzigen Konsumenten-Pfad
  (`writeLogsJSON` → `splitLogLines`), nicht an der Quelle.
- **Kein neuer Exit-Code und keine Meldung.** Das Verhalten bleibt
  still (Normalisierung, kein Befund); der Edge-Case ist kein
  Validierungsfall.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — ≤ 3 Liefer-Punkte.

- [ ] **CRLF-Behandlung gepinnt.** Unit-Tests gegen `splitLogLines`
      (Fundstelle der Test-Nachbarn gelesen:
      `internal/adapter/driving/cli/logs_acceptance_test.go`, Zeile 381
      Kommentar zum trailing-`\n`-Strip): CRLF-only-Zeilen, gemischte
      `\n`/`\r\n`-Zeilen, trailing `\r\r\n`, unverändertes `\n`-Verhalten;
      je Fall Happy/Edge (gepint gegen den `data.lines`-Envelope).
- [ ] **Der Carveout ist aufgelöst.** Die CO-Datei (aus
      `docs/plan/carveouts/`) wandert per reinen `git mv` nach
      `docs/plan/carveouts/done/` (Muster `modul-07-carveouts.md`
      §Ziel-Form: Carveout); der README-Index trägt die Zeile unter
      „Aufgelöste Carveouts“; der
      `Folge-Slice`-Verweis in der CO-Datei wird nach dem Move berichtigt
      (Ruheort-Regel, `grundlagen-traceability.md`).
- [ ] **Vertrags-Doku setzt den Ist-Stand.**
      [`docs/maintainer/cli-json-contract.md`](../../../../docs/maintainer/cli-json-contract.md) §6.8:
      der Absatz „CRLF-Edge-Case“ nennt den gemessenen/pintent Zustand, nicht
      die Grenze; AGENTS.md §3.10 (Zustandsmarker mit Ablaufpunkt).
- [ ] `make gates` grün. *(nicht gezählt — Gate-Lauf)*
- [ ] Review, Verifikation, Closure-Notiz, Register, Risiko-Ausgänge,
      Drei-Paarungen. *(nicht gezählt — Closure-Pflichten)*

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — Pfad-Kandidaten für §8.

| Datei / Komponente | Art | Begründung / Fundstelle |
|---|---|---|
| `internal/adapter/driving/cli/logs.go` | update | `splitLogLines` (Zeilen 309 bis 320): Normalisierung des `\r` vor dem `\n`-Split, Doc-Kommentar trägt die Zusage |
| `internal/adapter/driving/cli/logs_acceptance_test.go` | update | Unit-Pins je Fall (Happy/Edge/Negativ — getestet gegen den `data.lines`-Envelop) |
| `docs/plan/carveouts/CO-<NNN>-crlf-logs.md` | update → `done/` | Status „Aufgelöst“, Auflösungs-Semantik (`git mv` nach `done/`, Verifikations-Block abgehakt) |
| `docs/plan/carveouts/README.md` | update | Index-Zeile: von „Aktive Carveouts“ zu „Aufgelöste Carveouts“ (neuer Pfad `done/`) |
| [`docs/maintainer/cli-json-contract.md`](../../../../docs/maintainer/cli-json-contract.md) | update | §6.8: Absatz setzt den Ist-Anstand (Zustandsmarker mit Ablaufpunkt entfernt, AGENTS.md §3.10) |
| Diese Datei, Verweise | update | href-Berichtigungen nach dem `git mv` (Ruheort-Regel) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): **Erster bestätigter Befund** über ein
trailing `\r` in der `--json`-Antwort von `u-boot logs` (Konsumer-Meldung,
Issue oder Review-Finding); Projektinhaber gibt den Auftrag; kein anderer
Slice in `in-progress/` (WIP-Limit 1). Ohne Befund bleibt der Plan in
`open/`; die Vertrags-Doku trägt die Grenze (das ist der aktuelle Stand,
nicht ein Mangel).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß): —, kein Fall vorweg zu sehen (der Fix
  ist eine Funktion plus Tests).
- `in-progress` → `open` (blockiert): Falls der Befund sich als nicht
  reproduzierbar erweist (Compose-Output ist `\n`-terminiert am gemessenen
  Pfad), kehrt der Slice zurück; der Carveout-Trigger ist dann mit dem
  Befund neu zu beurteilen.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln.

DoD vollständig (drei Liefer-Punkte plus Gate- und Closure-Pflichten);
Review und Verifikation; Closure-Notiz; der Carveout-Status in
`docs/plan/carveouts/` aufgelöst; `git mv` nach `done/`.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Über-Normalisierung:** `\r` außerhalb des Zeilenende (mid-line) könnte
  gestript werden, wenn der Fix breit fasst. — **Ausgang:** bei Closure
  (der Fix stript nur `\r` unmittelbar vor `\n`; ein Pin dafür).
- **Der Unit-Pin beweist kein echtes Verhalten** (der Edge-Case ist
  hypothetisch — kein bekannter Konsument-Pfad). — **Ausgang:** bei Closure;
  das Start-Trigger-Ereignis (echter Befund) ist der Beweis; der Pin trägt
  die Regression-Absicherung.
- **Vertrags-Doku-Differenz:** §6.8 könnte den Ist-Stand über- oder
  unter-sagen. — **Ausgang:** bei Closure (DoD-Punkt 3 + AGENTS.md §3.10:
  Aussage über Verhalten nennt nur das Gemessene).

## 7. Closure-Notiz

Noch nicht geschlossen; der Slice liegt in `open/` (Entwurf).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

**Vorgelagert — Sub-Area-Wahl prüfen:** `internal/adapter/driving/cli` ist
die deklarierte Sub-Area CLI (GF, graduiert 2026-07-25,
[`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area)) —
keine neue Kandidatin.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen am
2026-11-02 — keine Treffer: die offenen BEO-Einträge betreffen Devcontainer-
und Katalog-Verhalten, nicht die CLI-Zeilen-Verarbeitung; die verkörperten
BEO-CLI-Einträge (`kennung-in-cli-ausgabe`, `kommentar-beschreibt-abwesendes`)
sind andere Klassen. Keine Zählung.

**Modus-Begründungsblock:** entfällt — alle berührten Sub-Areas GF
(Refactor ohne neue Sub-Area-Berührung).