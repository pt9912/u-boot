# Slice Harness: Festlegungen aus den `internal/`-READMEs in die Spezifikation

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-harness-wartung`](../welle-harness-wartung.md) — folgt auf
[`slice-harness-dcheck-bump-reviewer-skill`](../done/slice-harness-dcheck-bump-reviewer-skill.md), damit sein Review mit dem aufgefrischten Skill läuft.

**Bezug:** [`LH-LESE-003`](../../../../spec/lastenheft.md#lh-lese-003--dokumentenordnung)
(Dokumentenordnung: technische Festlegungen stehen im technischen Dokument),
[`SPEC-012`](../../../../spec/spezifikation.md#spec-012--schichten-und-verzeichnislayout-unter-internal),
Sub-Area `internal/**/README.md` ([`conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration); Vorgänger:
[`slice-harness-internal-readme-kennungs-retrofit`](../done/slice-harness-internal-readme-kennungs-retrofit.md).
Auslöser: Frage des Projektinhabers (2026-10-16), ob die Inhalte der Paket-READMEs nicht in
die Spezifikation gehören. Keine Anforderungsänderung.

**Berührte Spec-Stellen:** `spec/spezifikation.md` (neue `SPEC-<NNN>` für Festlegungen,
die bisher nur im Code und in den READMEs stehen); `spec/architecture.md` (Schichtbeschreibung,
nur falls dort eine Lücke bleibt).

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die sieben `internal/**/README.md` (zusammen 490 Zeilen, rund 110 `LH-`-Verweise)
tragen nur noch, was in ein Paket-README gehört: Zweck der Schicht und Verweis auf die
Spezifikation. Technische Festlegungen (Zahlenwerte, Schwellen, Wahrheitstabellen,
Zustandsautomaten, Whitelists) stehen als `SPEC-<NNN>` in der Spezifikation; die
Zuordnung Service → Anforderung entfällt, sie folgt aus den Kennungen (`make doc-trace`).

**Befund (2026-10-16), Stichprobe `internal/hexagon/application/README.md`:**
`pollInterval=500ms`, `dialTimeout=300ms`, Restart-Loop-Zähler mit Schwelle 3, Dead-Allowlist,
§T5-Wahrheitstabelle des `--volumes`-Pfads und die 3-Element-Whitelist von `config` stehen nur
im Code und im README, nicht in der Spezifikation (0 Treffer in `spezifikation.md` und
`architecture.md`).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Änderung des Codes oder der Werte: der Slice dokumentiert den Ist-Stand; weicht ein Wert
  vom README ab, gilt der Code, und die Abweichung wird im Slice vermerkt.
- Neue Produktanforderungen: Festlegungen sind Verfeinerungen vorhandener Anforderungen
  (`LH-FA-UP-*`, `LH-FA-CONF-*` usw.), keine neuen Zusagen
  ([`LH-LESE-003`](../../../../spec/lastenheft.md#lh-lese-003--dokumentenordnung)).
- Go-Paketdoku (`doc.go`) für alle Pakete: nur dort, wo ein README ersetzt wird und der
  Implementierungshinweis (z. B. `noop*`-Defaults) am Code bleiben soll.

## 2. Definition of Done

- [ ] **Inventur:** je README eine Liste seiner Aussagen, eingeordnet als Festlegung
  (→ Spezifikation), Rückverfolgbarkeit (→ entfällt), Aufbau der Schicht (→ Architektur-Sicht
  oder Verweis) oder Implementierungshinweis (→ `doc.go` oder Kommentar); Abweichungen vom
  Code sind benannt.
- [ ] **Spezifikation:** die Festlegungen stehen als `SPEC-<NNN>` (ab der nächsten freien Nummer) in
  `spezifikation.md`, jede an der Anforderung verankert, die sie verfeinert, mit den
  Zahlenwerten aus dem Code.
- [ ] **READMEs** auf Zweck und Verweis gekürzt, die handgepflegte Service-Anforderungs-Liste
  entfällt; `make docs-check` grün.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/README.md` und sechs Paket-READMEs | update | auf Zweck und Verweis kürzen |
| `spec/spezifikation.md` | update | neue `SPEC-<NNN>` für die Festlegungen |
| `spec/architecture.md` | prüfen | Schichtbeschreibung nur bei Lücke |
| `internal/**/doc.go` | neu, wo nötig | Implementierungshinweise am Code |
| `harness/conventions.md` | prüfen | Zeile der Sub-Area `internal/**/README.md` |

Ansatz: ein README je Commit, beginnend mit `internal/hexagon/application/` (die meisten
Festlegungen); Zahlenwerte aus dem Code lesen, nicht aus dem README übernehmen.

## 4. Trigger

**Start** (`next` → `in-progress`): Priorisierung durch den Projektinhaber; kein anderer
Slice in `in-progress/`.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): nach Schicht aufteilen (`application`, `adapter`, `port`,
  `domain`), je ein Slice.
- `in-progress` → `open` (blockiert): ein README-Wert widerspricht dem Code, und die Frage,
  welcher gilt, ist nicht entschieden.

## 5. Closure-Trigger

DoD vollständig, Review- und Verification-Evidence liegen vor, Closure-Notiz
geschrieben, `make gates` grün.

## 6. Risiken und offene Punkte

- Ein README nennt einen Wert, der vom Code abweicht; die Spezifikation übernähme die
  falsche Zahl — **Ausgang:** offen; Werte kommen aus dem Code.
- Die Spezifikation wächst um viele Einträge und wird unübersichtlich —
  **Ausgang:** offen; Einträge je Anforderung bündeln, nicht je Konstante.
- Die Sub-Area `internal/**/README.md` ist in `conventions.md` als Greenfield deklariert;
  ihre Begründung nennt den Zweck „beschreiben den Code-Bestand“ — **Ausgang:** offen; die
  Zeile wird bei der Closure nachgezogen.

## 7. Closure-Notiz

Offen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `internal/**/README.md` (Sub-Area
Code-Paket-READMEs) und `spec/`
([`conventions.md`](../../../../harness/conventions.md) §Modus-Deklaration); beide erfüllen
die Schwelle, keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** das Register ist leer.

**Modus:** alle berührten Sub-Areas GF.
