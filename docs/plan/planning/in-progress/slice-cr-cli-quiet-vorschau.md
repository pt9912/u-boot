# Slice CR: Die Vorschau von `--dry-run` und `--diff` bleibt unter `--quiet` (`LH-FA-CLI-005`)

> **CR-Vehikel.** u-boot führt kein separates CR-Artefakt — dieser Slice
> **ist** der Change Request am Vertrags-Stratum (Precedent:
> [`slice-cr-cli-quiet-weitere-ausnahmen`](../done/slice-cr-cli-quiet-weitere-ausnahmen.md)).
> Anlass: Der zweite Review des Code-Slices
> [`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md)
> ([Report](../../../reviews/2026-10-16-slice-v1-cli-quiet-nur-fehler-2.md), F-2) fand eine Spannung zwischen
> [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) und
> [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run). Der Projektinhaber hat am 2026-10-16
> entschieden: Die Vorschau-Ausgabe von `--dry-run` und `--diff` ist das Ergebnis des Aufrufs und bleibt unter `--quiet`
> (wie die Abfragebefehle); das Lastenheft wird präzisiert, der Code folgt.

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-cli-verbosity`](../welle-cli-verbosity.md) — vierter Slice der Welle; ihm gehen
[`slice-cr-cli-quiet-fehler-und-warnungen`](../done/slice-cr-cli-quiet-fehler-und-warnungen.md),
[`slice-cr-cli-quiet-weitere-ausnahmen`](../done/slice-cr-cli-quiet-weitere-ausnahmen.md) und
[`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md) voraus, ihm folgen
[`slice-v1-cli-quiet-vorschau`](../open/slice-v1-cli-quiet-vorschau.md) und
[`slice-v1-cli-verbosity-letzte-option-gewinnt`](../open/slice-v1-cli-verbosity-letzte-option-gewinnt.md).

**Bezug:** [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Änderung),
[`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) (Dry-Run zeigt die Dateien) und
[`LH-FA-CLI-008`](../../../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) (Diff-Ausgabe als Vorschau).
Verwandt: [`slice-v1-cli-quiet-vorschau`](../open/slice-v1-cli-quiet-vorschau.md) (Code: setzt die Präzisierung um).

**Berührte Spec-Stellen:** `lastenheft.md` §3 ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)), §16
(Historie), Versionsfeld. `spezifikation.md` und `architecture.md` berührt der
CR nicht; sie beschreiben das Verhalten erst nach allen Slices der Welle.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

**Ziel:** Lastenheft-Änderung **0.9.0 → 0.10.0**: [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) benennt die
Ausgabe der Vorschau-Optionen `--dry-run` und `--diff` als Ergebnis des Aufrufs, das unter `--quiet` bleibt.

**Warum Präzisierung statt Auslegung:** [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) verlangt „Der Dry-Run muss anzeigen,
welche Dateien erzeugt, geändert oder gelöscht würden“ und, mit `--diff`, „die Ausgabe besteht ausschließlich aus dem geplanten Änderungsplan und der
Diff-Darstellung“; [`LH-FA-CLI-008`](../../../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) verlangt, dass `--diff` den geplanten Endzustand als Vorschau zeigt.
[`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (0.9.0) lässt unter `--quiet` die Zusammenfassungen der ändernden Subkommandos entfallen und
nennt den Plan weder als Ausnahme noch als Ergebnis eines Abfragebefehls. Ein `--dry-run --quiet` ohne Ausgabe und mit Exit-Code 0 erfüllt
die Zusage von [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) nicht; ohne Benennung müsste der Code den Vorrang still entscheiden.

**Warum keine vierte Ausnahme, sondern ein Ergebnis-Satz:** Die vier Ausnahmen schützen Ausweisungen *neben* dem Ergebnis (Lockerung, Fallback,
Fortschritt, Überschreib-Pfade). Die Vorschau *ist* das Ergebnis des Aufrufs, genau wie die Ausgabe von `config get`; sie gehört in den
Satz zu den Abfragebefehlen und lässt die Zählung „vier Ausnahmen“ unverändert.

**Warum 0.10.0 und nicht 0.9.1:** Die Historie (§16) führt `0.3.x`-Patches für Klarstellungen, die einen Fall nachtragen, ohne eine Zusage zu
ändern, und Minor-Wechsel für jede Änderung, die das Verhalten eines Befehls berührt (`0.5.0`, `0.7.0`, `0.8.0`, `0.9.0`, alle an
[`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) oder an Zusagen mit Code-Folge). Hier ändert sich das beobachtbare Verhalten (`--dry-run --quiet` gibt aus statt zu schweigen),
und ein Umsetzungs-Slice folgt; das ist eine Minor-Änderung. Der Review nannte 0.9.1 als offenes INFO-Thema; die Wahl 0.10.0 folgt dem Muster der
vorherigen CRs zu [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) und ist vom Projektinhaber zu bestätigen.

### Gefundener Stand im Code

- `init`, `generate`, `add`, `remove` und `config set` schreiben den Diff vor dem Quiet-Zweig; unter `--quiet --diff` bleibt der Diff also schon heute.
- Die Zusammenfassung der Vorschau (`Would initialize …`, `Would add service …`, `Would change:` und die Entsprechungen bei `remove`)
  entfällt unter `--quiet`; ein reines `--dry-run --quiet` gibt nichts aus und endet mit Exit-Code 0 (Review F-2).
- Kein Test kombiniert `--quiet` mit `--dry-run` oder `--diff`.

### Prüfung gegen die älteren Zusagen (Steering-Loop-Lehre aus dem CR `slice-cr-cli-quiet-weitere-ausnahmen`)

Vor der Änderung geprüft per Suche im Lastenheft (`--dry-run`, `--diff`, `Vorschau` und die Kennungen der Dry-Run- und Diff-Anforderung):

| Stelle | Berührung | Ergebnis |
|---|---|---|
| [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) (Plan „muss anzeigen“; „ausschließlich Änderungsplan und Diff-Darstellung“) | Kollision mit `--quiet` | **aufgelöst durch den neuen Satz** |
| [`LH-FA-CLI-008`](../../../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) (`--diff` zeigt den Endzustand als Vorschau; mit `--dry-run` „dieselbe Vorschau bei vollständigem Schreibschutz“) | Kollision mit `--quiet` | **aufgelöst durch den neuen Satz** |
| [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) / [`LH-FA-CLI-008`](../../../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) mit `--json` | keine: `--quiet --json` ändert die JSON-Ausgabe nicht ([`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) wird nicht berührt) | geprüft, unverändert |
| Exit-Codes für Vorschau-Workflows (CLI-008, „selben Exit-Codes“; [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)) | keine | geprüft, unverändert |
| [`LH-FA-CLI-005A`](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) / [`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz) (`--force`-Zusammenfassung, Ausnahme 4) | keine: bleibt Ausnahme 4, unabhängig von der Vorschau | geprüft |
| weitere Fundstellen von `--dry-run`/`--diff` (Abschnitt zur maschinenlesbaren Ausgabe, Schema-Absätze) | betreffen nur `--json` | geprüft, keine Kollision |

Die Suche fand keine weitere Zusage, die mit „`--quiet` zeigt nur Fehler“ kollidiert. Zwei Fälle bleiben außerhalb der Präzisierung
und sind im Umsetzungs-Slice zu pinnen: die Erfolgsmeldung nach dem Schreiben bei `--diff` *ohne* `--dry-run` (sie ist keine Vorschau und entfällt) und
Warnungen neben der Vorschau (gelten nach den bestehenden Regeln; nur die vier Ausnahmen bleiben).

### CR-Kern (vorgeschlagene Fassung des Absatzes in [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging))

Der Absatz nach der Aufzählung der vier Ausnahmen lautet in 0.10.0 (Änderung kursiv markiert, im Lastenheft ohne Markierung):

> Alle anderen Hinweise und Warnungen entfallen unter `--quiet`. Die Ausgabe der Abfragebefehle, die nichts verändern (`config get`, `config list`, `logs`, `template list`), ist das Ergebnis des Befehls und bleibt. *Ebenso ist die Ausgabe der Vorschau-Optionen `--dry-run` und `--diff` ([`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-FA-CLI-008`](../../../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe)) — der geplante Änderungsplan und die Diff-Darstellung — das Ergebnis des Aufrufs und bleibt. Die Erfolgs- und Zusammenfassungsmeldung, die ein tatsächliches Schreiben abschließt, ist keine Vorschau und entfällt auch bei `--diff` ohne `--dry-run`.*

**Unverändert, mit Begründung:**

- Die vier Ausnahmen und ihre Zählung: Die Vorschau ist Ergebnis, keine Ausweisung neben dem Ergebnis (siehe oben).
- Die Stufen und der Satz zu mehreren Optionen: unberührt.
- [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) und [`LH-FA-CLI-008`](../../../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe): bleiben im Wortlaut; die Kollision löst der Satz in der Verbosity-Anforderung, nicht eine Änderung dort.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Umsetzung im Code** — Adapter ändern sich in [`slice-v1-cli-quiet-vorschau`](../open/slice-v1-cli-quiet-vorschau.md); dieser CR ändert nur das Vertrags-Stratum.
- **Der Compose-Strom von `down` (Review F-1)** — ein anderer Befund an einer anderen Zusage; er gehört in die Arbeit am Slice
  [`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md), nicht in die Vorschau.
- **Der CHANGELOG-Widerspruch (Review F-3)** — Doku des Code-Slices, kein Vertragsthema.
- **Eine Spec-Beschreibung in `spezifikation.md`** — erst nach allen Slices der Welle, sonst beschreibt sie einen Zwischenstand.

**Abnahmekriterien (für die Umsetzung, nicht für den Spec-Change)** — sie liegen in
[`slice-v1-cli-quiet-vorschau`](../open/slice-v1-cli-quiet-vorschau.md):

1. Je Subkommando mit `--dry-run`/`--diff` belegt ein Test, dass `--dry-run --quiet` und `--diff --quiet` die Vorschau ausgeben und Fehler bleiben.
2. Ein Test pinnt, dass die Erfolgsmeldung nach einem tatsächlichen Schreiben (`--diff` ohne `--dry-run`) unter `--quiet` entfällt.
3. `docs/user/benutzerhandbuch.md` und `CHANGELOG.md` nennen die Vorschau als Ergebnis unter `--quiet`.

**Historie-Zeile (§16, bei Spec-Change zu setzen):**

```text
| 0.10.0 | 2026-10-16 | LH-FA-CLI-005: Die Ausgabe der Vorschau-Optionen
`--dry-run` und `--diff` (Änderungsplan und Diff-Darstellung, nach
LH-FA-CLI-007 und LH-FA-CLI-008) ist das Ergebnis des Aufrufs und bleibt
unter `--quiet`; die Erfolgsmeldung eines tatsächlichen Schreibens entfällt
weiter; die vier Ausnahmen, die Stufen und die Mehrfachangabe-Regel bleiben
unverändert. | Vereinbarung mit dem Projektinhaber |
```

## 2. Definition of Done

- [ ] **Lastenheft 0.10.0**: [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) in der Fassung oben, Versionsfeld,
      Historie-Zeile (Verweis „Vereinbarung mit dem Projektinhaber").
- [ ] **Umsetzungsbedarf** einem Slice zugeordnet: Abnahmekriterien 1–3 liegen in
      [`slice-v1-cli-quiet-vorschau`](../open/slice-v1-cli-quiet-vorschau.md).
- [ ] Doku-Update: CHANGELOG (öffentlicher Vertrag berührt).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), Versionsfeld 0.10.0, §16-Zeile |
| `CHANGELOG.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | öffentlicher Vertrag berührt |

Reihenfolge: Lastenheft ändern → Umsetzung im Slice
[`slice-v1-cli-quiet-vorschau`](../open/slice-v1-cli-quiet-vorschau.md) (der Vorschau-Satz steht erst nach 0.10.0 im Vertrag).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): sofort priorisierbar; die Entscheidung des Projektinhabers vom 2026-10-16 liegt vor.
Er steht in der Welle hinter [`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md) (Lastenheft-Seite unabhängig davon, der Code-Slice dort bleibt auf 0.9.0).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Fassung verlangt, weitere Ausgabeklassen mitzuregeln
  (etwa die Konfliktmeldungen bei `--diff` ohne `--dry-run`) — dann je Klasse ein Umsetzungs-Slice.
- `in-progress` → `open` (blockiert): der Projektinhaber ändert die Entscheidung zur Vorschau unter `--quiet`.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit
Steering-Loop-Lerneintrag geschrieben (§7). Lastenheft 0.10.0 ist angenommen;
der Umsetzungs-Slice [`slice-v1-cli-quiet-vorschau`](../open/slice-v1-cli-quiet-vorschau.md) liegt in `open/` oder weiter.

## 6. Risiken und offene Punkte

- **Reichweite „Vorschau“ bei `--diff` ohne `--dry-run`:** Dieser Modus schreibt tatsächlich; Diff ist Vorschau, die anschließende Erfolgsmeldung nicht.
  Der Satz trennt beides ausdrücklich, aber der Projektinhaber hat nur „Vorschau-Ausgabe bleibt“ entschieden. — **Ausgang:** offen bis zur Bestätigung durch den
  Projektinhaber (siehe Übergabe); bei Gegenentscheidung entfällt der Schlusssatz des CR-Kerns.
- **Weitere Kollisionen mit älteren Zusagen:** Die Suche (siehe §1) fand keine über CLI-007 und CLI-008 hinaus. — **Ausgang:** entfallen: Suche nach
  `--dry-run`, `--diff`, `Vorschau`, den Kennungen der Dry-Run- und Diff-Anforderung und `--quiet` ausgewertet; der Verifier wiederholt sie vor dem Schließen.
- **Zwischenstand der Spec:** Bis alle Slices der Welle geschlossen sind, darf die Spezifikation das Verhalten nicht
  beschreiben. — **Ausgang:** entfallen: die Welle [`welle-cli-verbosity`](../welle-cli-verbosity.md) schließt erst mit dem
  Spezifikations-Eintrag im letzten Slice.

## 7. Closure-Notiz

Offen.

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** [`slice-v1-cli-quiet-vorschau`](../open/slice-v1-cli-quiet-vorschau.md),
  [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../open/slice-v1-cli-verbosity-letzte-option-gewinnt.md).
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist `spec/`, `harness/`,
`docs/plan/` (eine Sub-Area laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006), Abschnitt Modus-Deklaration,
Modus Greenfield). `CHANGELOG.md` ist ein öffentlicher Vertrag, keine Sub-Area; der Code
(`internal/adapter/driving/cli`) wird von diesem CR nicht geändert.

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register
([`../observations/`](../observations/README.md)) durchgegangen — keine Treffer
(Register leer).

Alle berührten Sub-Areas GF.
