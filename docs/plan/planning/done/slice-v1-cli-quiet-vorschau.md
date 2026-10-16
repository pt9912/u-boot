# Slice v1: Verbosity — die Vorschau von `--dry-run` und `--diff` bleibt unter `--quiet`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-cli-verbosity`](welle-cli-verbosity.md) — fünfter Slice der Welle; folgt auf
[`slice-cr-cli-quiet-vorschau`](../done/slice-cr-cli-quiet-vorschau.md) (der Vorschau-Satz steht erst mit Lastenheft 0.10.0 im Vertrag) und auf
[`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md) (berührt dieselben Quiet-Zweige); ihm folgt
[`slice-v1-cli-verbosity-letzte-option-gewinnt`](../done/slice-v1-cli-verbosity-letzte-option-gewinnt.md).

**Bezug:** [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Vorschau bleibt unter `--quiet`, ab 0.10.0),
[`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) (Dry-Run zeigt die Dateien) und
[`LH-FA-CLI-008`](../../../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) (Diff als Vorschau).
Herkunft: Review-Finding F-2 im [Report](../../../reviews/2026-10-16-slice-v1-cli-quiet-nur-fehler-2.md).

**Berührte Spec-Stellen:** `lastenheft.md` §3 ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)) — wird
nicht geändert; der Code folgt dem Lastenheft. `spezifikation.md` und `architecture.md` beschreiben das Verhalten erst nach allen Slices der Welle.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

**Ziel:** `--dry-run --quiet` und `--diff --quiet` geben die Vorschau aus (Zusammenfassung des Plans und Diff), Fehler bleiben, wie
[`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) ab 0.10.0 es zusagt; ein reines `--dry-run --quiet` schweigt nicht mehr.

**Stand heute (im Code gefunden):** Alle fünf Subkommandos mit `--dry-run`/`--diff` (`init`, `generate`, `add`, `remove`, `config set`
einzeln und mit mehreren Paaren) schreiben den Diff vor dem Quiet-Zweig, also bleibt der Diff unter `--quiet --diff` schon heute. Unter `--quiet` entfällt dagegen
die Zusammenfassung des Plans (`Would initialize …`/`Would create:`, `Would add service …`/`Would change:`, bei `remove`, `generate` und `config set` die
Zusammenfassungszeile), sodass `--dry-run --quiet` ohne Ausgabe und mit Exit-Code 0 endet. Kein Test kombiniert `--quiet` mit `--dry-run` oder `--diff`.
`config get`, `config list`, `config show` und `config` lehnen die Optionen mit `ErrDryRunNotApplicable` ab; `up`, `down`, `doctor`, `logs` und `template` kennen sie nicht.

**Regel (Lastenheft 0.10.0 entscheidet, bestätigt vom Projektinhaber am 2026-10-16):**

| Aufruf unter `--quiet` | Vorschau-Zusammenfassung | Diff | Erfolgsmeldung nach dem Schreiben |
|---|---|---|---|
| `--dry-run` | bleibt | — | — (es wird nicht geschrieben) |
| `--dry-run --diff` | bleibt | bleibt | — |
| `--diff` (ohne `--dry-run`) | — | bleibt | entfällt (keine Vorschau) |
| weder noch | — | — | entfällt (wie bisher) |

Die Sicherheits-/Fallback-Ausnahmen und Fehler bleiben unabhängig davon; sonstige Hinweise und Warnungen entfallen weiter. `--json` ist unberührt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine „Would …"-Formulierung für `remove`, `generate` und `config set` im Dry-Run** — heute melden diese Subkommandos auch ohne `--quiet` in der Vergangenheitsform
  (`Removed service …`, `Generated …`, `config: a → b.`); nur `init` und `add` unterscheiden. Das ist ein Befund ohne Bezug zu `--quiet`; dieser Slice gibt unter `--quiet` aus, was ohne `--quiet`
  ausgegeben wird. Offene Frage an den Projektinhaber, getragen vom Folge-Slice [`slice-v1-cli-dry-run-would-formulierung`](../open/slice-v1-cli-dry-run-would-formulierung.md) (siehe §6).
- **Der Compose-Strom von `down`, der CHANGELOG-Widerspruch (Review F-1, F-3)** — anderer Befund, Arbeit am Slice
  [`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md).
- **Eine Änderung des Lastenhefts** — liegt im CR [`slice-cr-cli-quiet-vorschau`](../done/slice-cr-cli-quiet-vorschau.md).
- **Die Wahl zwischen mehreren Verbosity-Optionen** — Slice [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../done/slice-v1-cli-verbosity-letzte-option-gewinnt.md).
- **Eine Spec-Beschreibung in `spezifikation.md`** — erst nach allen Slices der Welle, sonst beschreibt sie einen Zwischenstand.

## 2. Definition of Done

- [x] [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (0.10.0) erfüllt, [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) und
      [`LH-FA-CLI-008`](../../../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) wieder tragend: je Subkommando mit Vorschau (§3) belegen Tests Happy
      (`--dry-run --quiet`, `--dry-run --diff --quiet`, `--diff --quiet`), Boundary (Erfolgsmeldung nach dem Schreiben bei `--diff --quiet` entfällt; No-op-Plan) und Negative (Fehler bleibt sichtbar; `--json --quiet` unverändert).
- [x] Doku-Update: `docs/user/benutzerhandbuch.md` (Optionstabellen) und `CHANGELOG.md` nennen die Vorschau als Ergebnis unter `--quiet`.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne Wellen-Betrieb hier geprüft.

## 3. Plan (vor Code)

Subkommandos mit `--dry-run`/`--diff` (per Suche in `internal/adapter/driving/cli/` gefunden, ohne `*_test.go`), und was unter `--quiet` heute entfällt:

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapter/driving/cli/init.go` <!-- d-check:ignore (Pfad-Zitat) --> | update | `runInit`: `if flags.Quiet` lässt `printInitSummary` (`Would initialize`/`Would create:`, Backups) entfallen; bei Vorschau (`DryRun`) ausgeben, Warnungen wie bisher gefiltert. Gilt nicht bei `--template` (dort sind die Optionen ausgeschlossen) |
| `internal/adapter/driving/cli/add.go` <!-- d-check:ignore (Pfad-Zitat) --> | update | `runAdd`: `if flags.Quiet { return nil }` lässt `printAddSummary` (`Would add service`/`Would repair`/`Would change:`) entfallen; bei `DryRun` ausgeben |
| `internal/adapter/driving/cli/remove.go` <!-- d-check:ignore (Pfad-Zitat) --> | update | `runRemove`: `if !flags.Quiet` um `printRemoveSummary` lässt Zusammenfassung und Volumes-Zeile entfallen; bei `DryRun` ausgeben; Purge-Warnungen bleiben unter `--quiet` entfallen (`PreviewDryRun` erzeugt sie ohnehin nicht) |
| `internal/adapter/driving/cli/generate.go` <!-- d-check:ignore (Pfad-Zitat) --> | update | `runGenerate`: `if !flags.Quiet` um `printGenerateSummary` lässt die Zeile entfallen; bei `DryRun` ausgeben |
| `internal/adapter/driving/cli/config.go` <!-- d-check:ignore (Pfad-Zitat) --> | update | `runConfigSet` (`if !flags.Quiet` um `printConfigSetSummary`) und `runConfigSetMany` (`if flags.Quiet { return nil }`): bei `DryRun` ausgeben |
| `internal/adapter/driving/cli/previewmode.go` <!-- d-check:ignore (Pfad-Zitat) --> | prüfen | Wahrheitstabelle bleibt; ein kleiner Helfer (etwa „Vorschau-Zusammenfassung unterdrücken: `quiet && !dryRun`“) vermeidet fünf Kopien der Bedingung |
| `internal/adapter/driving/cli/quiet_preview_test.go` <!-- d-check:ignore (Pfad-Zitat) --> | neu | Happy/Boundary/Negative je Subkommando nach [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-FA-CLI-008`](../../../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe); die Bestandstests ohne Vorschau (`quiet_modifying_test.go`) bleiben unverändert (Kontrolle: Erfolgsmeldung entfällt weiter) |
| `docs/user/benutzerhandbuch.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | die beiden `--quiet`-Tabellenzeilen (Zeilen um 299 und 1280) und der Absatz zu `--dry-run` nennen die Vorschau als Ergebnis |
| `CHANGELOG.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | sichtbare Verhaltensänderung: `--dry-run --quiet` schweigt nicht mehr |

Unter `--quiet` entfällt heute: bei `init`, `add`, `remove`, `generate`, `config set` und `config set` mit mehreren Paaren jeweils die **Zusammenfassung des Plans bei `--dry-run`**
(und `--dry-run --diff`); der Diff selbst bleibt bereits. Reihenfolge der Ausgabe bleibt: Diff, dann Zusammenfassung.

**Reihenfolge:** erst der CR [`slice-cr-cli-quiet-vorschau`](../done/slice-cr-cli-quiet-vorschau.md) (Lastenheft 0.10.0), dann dieser Slice; er berührt dieselben Quiet-Zweige wie
[`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md) und beginnt erst nach dessen Closure, damit keine Merge-Konflikte entstehen.

## 4. Trigger

**Start** (`next` → `in-progress`): Lastenheft 0.10.0 liegt vor (CR [`slice-cr-cli-quiet-vorschau`](../done/slice-cr-cli-quiet-vorschau.md) in `done/`) und
[`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md) ist geschlossen.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die „Would"-Formulierung für `remove`, `generate` und `config set` wird verlangt — dann eigener Slice.
- `in-progress` → `open` (blockiert): der Projektinhaber ändert die Entscheidung zur Vorschau unter `--quiet`, oder 0.10.0 ist nicht angenommen.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit
Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **Auslegung „Vorschau“ bei `--diff` ohne `--dry-run`:** Die Regel in §1 lässt die Erfolgsmeldung nach dem Schreiben entfallen und den Diff bleiben. — **Ausgang:**
  entfallen: [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) 0.10.0 entscheidet wörtlich (die Erfolgsmeldung entfällt auch bei `--diff` ohne `--dry-run`),
  bestätigt vom Projektinhaber; Code und Pin folgen.
- **Vergangenheitsform im Dry-Run (`remove`, `generate`, `config set`):** Unter `--quiet --dry-run` erscheint dieselbe Vergangenheitsform wie heute ohne `--quiet`,
  obwohl nichts geschrieben wurde. — **Ausgang:** eingetreten:
  [`slice-v1-cli-dry-run-would-formulierung`](../open/slice-v1-cli-dry-run-would-formulierung.md).
- **Konflikt mit dem Code-Slice [`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md):** beide ändern die Quiet-Zweige derselben fünf Dateien. — **Ausgang:** entfallen:
  dieser Slice startet erst nach dessen Closure (siehe §4).
- **Zwischenstand der Spec:** Die Spezifikation darf das Verhalten erst nach allen Slices der Welle beschreiben. — **Ausgang:** entfallen: der Spezifikations-Eintrag liegt im letzten Slice
  [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../done/slice-v1-cli-verbosity-letzte-option-gewinnt.md).

## 7. Closure-Notiz

**Belege:** Review-Report
[`2026-10-16-slice-v1-cli-quiet-vorschau`](../../../reviews/2026-10-16-slice-v1-cli-quiet-vorschau.md) (0 HIGH, 0 MEDIUM, 1 LOW, 3 INFO; F-1 und F-2 behoben in `7ee6c76`, F-3 in den
Folge-Slice, F-4 als Grenze). Verification Evidence
[`2026-10-16-verification-slice-v1-cli-quiet-vorschau`](../../../reviews/2026-10-16-verification-slice-v1-cli-quiet-vorschau.md): bestanden, mit vier Auflagen für die Closure
(hier eingelöst). `make gates` grün, Coverage 91,10 % (Schwelle 90 %); `make doc-trace`: 83 Anforderungen, 0 Waisen. Fünf Mutationen des Reviewers und zwei des Verifiers fielen rot.

- **Was hat funktioniert:** Die Regeltabelle aus dem CR ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) 0.10.0) diente als Testmatrix: je Subkommando mit Vorschau
  jede Zeile × jede Option, die Erwartung ist die Tabelle, nicht der Code. Der Helfer `summaryVisible` hält die Bedingung an einer Stelle statt in fünf Kopien. Der Reviewer prüfte jede
  Schreibstelle systematisch gegen die Tabelle. Es mussten keine bestehenden Pins geändert werden, weil kein Bestandstest `--quiet` mit `--dry-run` oder `--diff` kombinierte.
- **Was ging anders als geplant:** Der Kommentarblock von `summaryVisible` klebte am Doc-Kommentar von `previewModeFromFlags` (Review F-1). Der CHANGELOG-Satz sagte „the same text“ bei `init`,
  wo Warnungen gefiltert werden (F-2). §3 nannte die falsche Testdatei (`quiet_modifying_test.go` statt der neuen `quiet_preview_test.go`), und §1/§6 behaupteten eine offene Bestätigung, die
  Lastenheft 0.10.0 bereits wörtlich gab; beides berichtigt.
- **Steering-Loop-Eintrag:** Ein neuer Helfer mit Doc-Kommentar gehört **nach** die Funktion, deren Doc-Kommentar davor steht: Ein Einfügen zwischen Doc-Kommentar und Funktion hängt beide
  zusammen. Das ist zum dritten Mal in dieser Welle ein Kommentar-Befund aus `AGENTS.md` §3.7 (Kommentar beschreibt nicht, was da ist). Gezählt, nicht verkörpert. Befund für die Welle-Closure:
  die Klasse „Kommentar und Gegenstand getrennt bzw. falsch zugeordnet“ berührt die 3×-Schwelle der Welle; ob sie als Regel verkörpert wird, entscheidet die Welle-Closure anhand der Slice-Belege,
  hier wird nichts verkörpert.
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung angelegt; das Register ist leer, also keine Beobachtung über der Schwelle.
- **Folge-Slices:** [`slice-v1-cli-dry-run-would-formulierung`](../open/slice-v1-cli-dry-run-would-formulierung.md) (neu: Zeitform der Dry-Run-Zusammenfassung, Review F-3) und
  [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../done/slice-v1-cli-verbosity-letzte-option-gewinnt.md) (Rangfolge bei Mehrfachangabe, Spec-Eintrag).
- **Release-Hinweis:** kein Release-Cut vor `slice-v1-cli-verbosity-letzte-option-gewinnt`; der CHANGELOG sagt „follows in a later release“.
- **Risiken aus §6:** alle tragen ihren Ausgang: eingetreten (Vergangenheitsform im Dry-Run) oder entfallen (Auslegung `--diff` ohne `--dry-run`, Konflikt mit dem Code-Slice, Zwischenstand der Spec).
- **Grenzen:** (1) Kein Dry-Run-Pin mit Warnungen für `generate` (Review F-4); das Verhalten ist durch Lesen bestätigt, der Warn-Pin im Dry-Run besteht nur für `init`.
  (2) Der No-op-Pin (`TestExecute_Add_Quiet_DryRunNoOpPlanStays`) existiert nur für `add`; die übrigen Subkommandos teilen den Zweig.
- **Drei Paarungen:** Anker (kein `liegt in`, nichts verkörpert) · Folge-Slice (beide in `open/`) · Register (leer, notiert) — getragen; Bestandseintrag in `open/README.md` bereinigt.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist `internal/adapter/driving/cli` (Sub-Area CLI laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006), Modus Greenfield, graduiert)
sowie `docs/plan/` und `CHANGELOG.md`/`docs/user/` als öffentliche Verträge (keine Sub-Area).

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register
([`../observations/`](../observations/README.md)) durchgegangen — keine Treffer
(Register leer).

Alle berührten Sub-Areas GF.
