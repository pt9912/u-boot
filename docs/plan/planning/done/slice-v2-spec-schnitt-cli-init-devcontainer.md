# Slice V2: Lastenheft-Schnitt 1: CLI-Grundverhalten, Projektinitialisierung, Devcontainer

> **Status:** **Done** (2026-10-02, **Delivery-Hash: `8119daa`**).

**Welle:** `welle-spec-technik-stratum` (siehe [`roadmap.md`](../in-progress/roadmap.md)).
**Bezug:** Lastenheft §4.1 CLI-Grundverhalten, §4.2 Projektinitialisierung, §4.3 Devcontainer-Unterstützung, §6.1 Kommandozeilenschnittstelle (Scope).
**Berührte Spec-Stellen:** `spec/spezifikation.md` §1–§4, §6; `spec/lastenheft.md` §4.1–§4.3, §6.1.
**Autor:** pt9912. **Datum:** 2026-10-02.

## Auslöser

Größte Technik-Dichte im Lastenheft: JSON-Schema der `--dry-run`/`--diff`-Ausgabe, Diagnose-Codes, Ablaufbeschreibungen der
Init-Zustandsmaschine, Devcontainer-Defaults. Dieser Slice ist der **Pilot** des Schnitts: Er legt Form und Granularität
der Verfeinerungen fest, an der die übrigen Schnitt-Slices sich orientieren.

## Ziel und Abgrenzung

**Ziel:** Die Kapitel §4.1–§4.3 und §6.1 enthalten nur noch Vertrag (Was, Akzeptanzkriterien, Out-of-Scope); jede technische Festlegung steht wörtlich in der Spezifikation mit Kennung.

**Schnittregeln:** siehe [`slice-v2-spec-technik-stratum-grundlage`](../done/slice-v2-spec-technik-stratum-grundlage.md) §Schnittregeln (bleibt im Lastenheft / zieht in die Spezifikation, Decken-Regel, Verbleib-Tabelle, Zeilenverweise).

**Ausdrücklich NICHT in diesem Slice:**

- **Kapitel §4.4 und folgende** — eigene Schnitt-Slices; ein Pilot, der alles anfasst, ließe sich nicht mehr gegenlesen.
- **Neuformulierung von Anforderungen** — der Schnitt entscheidet nichts neu; Unklares bleibt im Lastenheft (siehe Offen).


## Definition of Done

- [x] §4.1–§4.3, §6.1 geschnitten; Verbleib-Tabelle vollständig (jede Quelle genau ein Ziel); `LH-*`-Überschriften und Anker unverändert.
- [x] Spezifikation um die übernommenen Stellen ergänzt (Verfeinerungen `.a`, `SPEC-<NNN>`), ohne ADR- und Slice-Verweis; Zeilenverweise auf die berührten Stellen umgestellt.
- [x] `make gates` grün (inklusive `make docs-check`), Exit-Code separat geprüft.
- [x] Review durch eine andere Rolle als die Umsetzung (`harness/review.md`), Report unter `docs/reviews/`.

## Tranchen

| T | Inhalt |
| - | ------ |
| T1 | **Inventur.** Je Anforderung in §4.1–§4.3 und §6.1 klassifizieren (bleibt / zieht um), Verbleib-Tabelle im Slice anlegen. |
| T2 | **Übernahme CLI.** §4.1 und §6.1: Schemata, Codes, Beispielinstanzen in die Spezifikation (§2, §4). |
| T3 | **Übernahme Init und Devcontainer.** §4.2 und §4.3: Zustandsfolgen, Defaults, Dateiformate. |
| T4 | **Gegenlesen.** Verbleib-Tabelle gegen Diff prüfen, `make gates`, Review. |

## Risiken

- Anforderungstext und Technik sind in einem Absatz verwoben (zum Beispiel Akzeptanzkriterium mit eingebettetem Schema): Schnitt an der Satzgrenze, Rest bleibt im Lastenheft; Zweifel unter „Offen“ statt raten.
- Link-Ziele auf verschobene Beispiele (Code-Kommentare, Tests) brechen: T4 sucht per `grep` nach `lastenheft.md:`-Zeilenverweisen.

## Closure

- **Geliefert (`8119daa`):** Lastenheft §4.1–§4.3 und §6.1 geschnitten. In die Spezifikation gezogen: JSON-Schema der Vorschau-Ausgabe, zwei Beispielinstanzen, Normalisierung des Projektnamens (Verfeinerung), Projektnamen-Muster, Übergabe der Benutzer-ID an den Image-Build (Verfeinerung). Anforderungstext, Akzeptanzkriterien und Produktverträge blieben im Lastenheft; alle `LH-*`-Überschriften und Anker sind unverändert.
- **Schnittregel in der Praxis (Pilot-Erkenntnis):** Als Technik gilt Mechanismus, Format, Algorithmus und Wert; beobachtbares Verhalten bleibt Vertrag. Deshalb blieben Entscheidungslogik (Bestätigungsmodi), Degradationstabelle und Sicherheitsverhalten im Lastenheft. Die Kürzung ist damit kleiner als die Vorlage-Idealform; die Folge-Slices wenden dieselbe Regel an.
- **Zeilenverweise:** 482 Verweise der Form `Spec §NNN` in 99 Dateien (Go-Kommentare, READMEs, Wartungs-Doku) wurden auf Anforderungs-Kennungen umgestellt. Die Nummern stammten teils aus älteren Lastenheft-Ständen; aufgelöst wurde über `git blame` (Stand des Kommentar-Commits). Erzeugte Vorlagen (`.tmpl`) mit eingebettetem `§611` blieben unangetastet (Produktausgabe).
- **Gegenlesen:** Alle 147 aus dem Lastenheft entfernten Zeilen sind in der Spezifikation wiederzufinden (Skript-Abgleich, 0 fehlend).
- **Sensoren:** `make gates` grün (lint, test, coverage-gate, docs-check). Nicht ausgeführt: `make ci`, `make test-docker`.
- **Review:** unabhängiges Review über die gesamte Welle am Abschluss.
- **Lerneintrag:** Zeilenverweise sind eine Drift-Quelle ohne Sensor; mit Kennungen ist der Schnitt der Folge-Slices für den Code folgenlos.

### Verbleib-Tabelle

| Quelle (Lastenheft) | Ziel (Spezifikation) | Ort |
|---|---|---|
| §4.1 [LH-FA-CLI-007](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) (Schema, 320–418) | [SPEC-001](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) | spezifikation.md §2 |
| §4.1 [LH-FA-CLI-007](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) (Beispielinstanz, 422–441) | [SPEC-002](../../../../spec/spezifikation.md#spec-002--beispielinstanz-einer-vorschau-ausgabe-add---dry-run---json) | spezifikation.md §2 |
| §4.1 [LH-FA-CLI-008](../../../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) (Beispielinstanz, 470–487) | [SPEC-003](../../../../spec/spezifikation.md#spec-003--beispielinstanz-einer-diff-ausgabe-add---diff---json-ohne---dry-run) | spezifikation.md §2 |
| §4.2 [LH-FA-INIT-002](../../../../spec/lastenheft.md#lh-fa-init-002--projektname) (Normalisierungsschritte 1–7, 518–524) | [LH-FA-INIT-002.a](../../../../spec/spezifikation.md#lh-fa-init-002a--normalisierung-des-abgeleiteten-projektnamens) | spezifikation.md §1 |
| §4.2 [LH-FA-INIT-006](../../../../spec/lastenheft.md#lh-fa-init-006--projektnamen-validierung) (regulärer Ausdruck, 640) | [SPEC-004](../../../../spec/spezifikation.md#spec-004--projektnamen-muster) | spezifikation.md §2 |
| §4.3 [LH-FA-DEV-004](../../../../spec/lastenheft.md#lh-fa-dev-004--benutzerrechte) (Build-Argument, 734) | [LH-FA-DEV-004.a](../../../../spec/spezifikation.md#lh-fa-dev-004a--übergabe-der-benutzer-id-an-den-image-build) | spezifikation.md §1 |
