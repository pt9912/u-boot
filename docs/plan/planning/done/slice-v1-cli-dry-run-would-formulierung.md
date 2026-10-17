# Slice v1: Dry-Run — die Zusammenfassung sagt „Would …“ in allen Subkommandos

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — keine Closure-Bedingung, die von der DoD dieses Slice verschieden ist, und kein Release-Bezug.
Herkunft: Review-Finding F-3 im [Report](../../../reviews/2026-10-16-slice-v1-cli-quiet-vorschau.md) des Slice
[`slice-v1-cli-quiet-vorschau`](../done/slice-v1-cli-quiet-vorschau.md) aus der Welle [`welle-cli-verbosity`](../done/welle-cli-verbosity.md).

**Bezug:** [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) (Dry-Run zeigt, was geschehen würde),
[`LH-FA-CLI-008`](../../../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) (Diff als Vorschau).

**Berührte Spec-Stellen:** `lastenheft.md` §3 ([`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run)) — wird in diesem Slice nicht geändert; bei Lösung C (siehe §1)
entscheidet ein CR, kein Code. `spezifikation.md` und `architecture.md` — keine Berührung.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Zusammenfassung eines Dry-Runs sagt in allen Subkommandos mit Vorschau, was geschehen **würde**, oder der Slice stellt fest, dass der Vertrag die Zeitform nicht festlegt.

**Stand heute (im Code gefunden):** Im Dry-Run melden `init` und `add` mit „Would …“ (`Would initialize`/`Would create:`, `Would add service`/`Would repair service`/`Would change:`).
`remove`, `generate` und `config set` melden dieselbe Zeile wie nach dem Schreiben, in der Vergangenheitsform (`Removed service …`, `Generated …`, `config: a → b.`), obwohl nichts geschrieben wurde.
Der Befund besteht **auch ohne `--quiet`** und hat keinen Bezug zu `--quiet`; er fiel im Review des Slice
[`slice-v1-cli-quiet-vorschau`](../done/slice-v1-cli-quiet-vorschau.md) auf (F-3), der unter `--quiet` ausgibt, was ohne `--quiet` ausgegeben wird.

**Wortlaut [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run):** „Der Dry-Run muss anzeigen, welche Dateien erzeugt, geändert oder gelöscht **würden**, ohne Änderungen am Dateisystem vorzunehmen.“
Die Zusage nennt den Konjunktiv für den Inhalt der Anzeige, aber kein Ausgabeformat und keine Zeilen-Formulierung; die Zeitform der Zusammenfassungszeile legt sie nicht ausdrücklich fest.

**Offene Frage an den Projektinhaber** (Lastenheft-Auslegung oder reiner Code-Befund):

- Widerspricht die Vergangenheitsform („Removed service …“ im Dry-Run) dem „würden“ von [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run)? Die Zusammenfassung behauptet ein Ergebnis, das nicht eingetreten ist;
  eine Leserin, die nur die Zusammenfassung sieht, kann einen Dry-Run für einen echten Lauf halten. Das spricht für einen Verstoß gegen den Sinn der Zusage.
  Gegen einen Verstoß spricht: Die Zusage verlangt den Inhalt (welche Dateien), nicht die Zeitform; die Anzeige der Dateien ist vorhanden.
- Das ist eine Auslegung des Lastenhefts und gehört dem Projektinhaber; der Planer entscheidet sie nicht.

**Lösungen:**

- **A — Code vereinheitlicht:** `remove`, `generate` und `config set` (einzeln und mit mehreren Paaren) sagen im Dry-Run „Would …“; Tests pinnen Text je Subkommando, Handbuch und CHANGELOG folgen.
  Gewählt, wenn die Auslegung „Vergangenheitsform widerspricht [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run)“ gilt.
- **B — kein Wechsel:** [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) legt die Zeitform nicht fest; der Slice schließt ohne Code (CR-frei) mit einer Begründung in §7 und `Gegenstand:` — entfallen.
- **C — Vertrag schärfen:** Der Projektinhaber will die Zeitform im Vertrag festhalten; dann entsteht ein CR nach dem Muster [`slice-cr-cli-quiet-vorschau`](../done/slice-cr-cli-quiet-vorschau.md), der Code folgt (Lösung A).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Verhalten unter `--quiet`** — [`slice-v1-cli-quiet-vorschau`](../done/slice-v1-cli-quiet-vorschau.md) setzt es um; der Befund besteht davon unabhängig, und dieser Slice ändert keine Quiet-Zweige.
- **Die Rangfolge der Verbosity-Optionen** — [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../done/slice-v1-cli-verbosity-letzte-option-gewinnt.md).
- **Das JSON-Format des Dry-Runs** (`--dry-run --json`) — ein maschinenlesbarer Vertrag ([`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run), Pflichtfelder), der von der Zeitform der Text-Zeile unberührt bleibt.
- **Die Texte von `init` und `add`** — sie sagen schon „Would …“; Änderung nur, falls die Vereinheitlichung auf eine gemeinsame Form es verlangt (dann im Plan §3 ergänzen).

## 2. Definition of Done

- [x] Die Auslegung von [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) ist vom Projektinhaber entschieden und in §7 festgehalten; bei Lösung A erfüllt: je Subkommando
      (`remove`, `generate`, `config set` einzeln und mit mehreren Paaren) belegen Tests Happy (Dry-Run sagt „Would …“), Boundary (kein Dry-Run: Vergangenheitsform bleibt) und Negative (Fehler und `--json` unverändert).
- [x] Doku-Update, falls Lösung A: `docs/user/benutzerhandbuch.md` und `CHANGELOG.md` nennen die neue Formulierung; bei B oder C entfällt es bzw. folgt dem CR.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne Wellen-Betrieb hier geprüft.

## 3. Plan (vor Code)

Gilt für Lösung A; bei B entfällt der Code, bei C kommt der CR davor.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapter/driving/cli/remove.go` <!-- d-check:ignore (Pfad-Zitat) --> | update | Zusammenfassung im Dry-Run als „Would remove service …“ (Zeile „Changed:“ analog „Would change:“) |
| `internal/adapter/driving/cli/generate.go` <!-- d-check:ignore (Pfad-Zitat) --> | update | Zusammenfassung im Dry-Run als „Would generate …“ |
| `internal/adapter/driving/cli/config.go` <!-- d-check:ignore (Pfad-Zitat) --> | update | `config set` (einzeln und mit mehreren Paaren) im Dry-Run als „Would set …“ statt `config: a → b.` |
| `internal/adapter/driving/cli/quiet_preview_test.go` <!-- d-check:ignore (Pfad-Zitat) --> und Tests je Subkommando <!-- d-check:ignore (Pfad-Zitat) --> | update | die Matrix pinnt heute die Vergangenheitsform als `summary`; Texte nachziehen, Happy/Boundary/Negative nach [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) und [`LH-FA-CLI-008`](../../../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) |
| `internal/adapter/driving/cli/dry_run_would_test.go` <!-- d-check:ignore (Pfad-Zitat) --> | add | Happy/Boundary/Negative je Subkommando (`printGenerateSummary` je `GenerateAction`, `printRemoveSummary`, `printConfigSetSummary`; Dry-Run mit und ohne `--diff`/`--quiet`, Schreiben mit und ohne `--diff`, `--dry-run --json`, Fehler; Drift-Wächter über alle `GenerateAction`-Werte, die der Port benennt, und `config set` mit gemischten Paaren: unverändert → „already X; no changes“, geändert → „Would set …“) |
| `internal/adapter/driving/cli/remove_acceptance_test.go` <!-- d-check:ignore (Pfad-Zitat) --> | update | `TestRemove_PurgeDryRun_HumanMode_NoVolumeWarningPollution` pinnt die Dry-Run-Headline; der Kommentar „Headline-Wording bleibt heute …“ beschreibt den neuen Zustand |
| `docs/user/benutzerhandbuch.md`, `CHANGELOG.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | sichtbare Textänderung im Dry-Run (Handbuch: Abschnitt „Vorschau statt Überraschung“; CHANGELOG: `[Unreleased]` · Changed) |
| `docs/plan/planning/in-progress/roadmap.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | der Ruhe-Marker in „Offene Wellen“ entfällt, solange der Slice in `in-progress/` liegt (d-check `planning`); der Planner setzt ihn nach der Closure zurück |
| `docs/reviews/2026-10-17-slice-v1-cli-dry-run-would-formulierung.md`, `docs/reviews/2026-10-17-verification-slice-v1-cli-dry-run-would-formulierung.md` <!-- d-check:ignore (Pfad-Zitat) --> | add | Review-Report und Verification Evidence (DoD-Artefakte) |
| `docs/plan/planning/observations/BEO-WZG/host-werkzeug-ausserhalb-docker/`, `docs/plan/planning/observations/BEO-CLI/dry-run-log-vergangenheitsform/` <!-- d-check:ignore (Pfad-Zitat) --> | add | Beobachtungs-Register: zwei neue Beobachtungen mit je einer Evidence-Datei (Closure §7) |
| `docs/plan/planning/open/README.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | die Bestandszeile dieses Slice entfällt (er liegt nicht mehr in `open/`) |

Die Änderung reicht über die drei Dateien hinaus nicht: `init` und `add` bleiben unverändert (ihre Formen sind schon „Would …“). Fundstellen: `printRemoveSummary` (`remove.go`), `printGenerateSummary` (`generate.go`), `printConfigSetSummary` (`config.go`, einzeln und über `runConfigSetMany`). Spec und Spezifikation geben den Text nicht vor.

Abweichung vom Vorgehen (2026-10-17): Der Implementer hat einmal `python3` auf dem Host ausgeführt (Verstoß gegen `AGENTS.md` §3.1). Der Aufruf änderte nur Zeilen in `remove_acceptance_test.go` (Kommentar-Berichtigung und eine zusätzliche Headline-Assertion); der Nachweis dieser Änderung läuft allein über `make test` und `make gates` im Docker. Seitdem kein Host-Werkzeug mehr.

Grenze: Der defensive `default`-Zweig von `printGenerateSummary` kennt `dryRun` nicht. Ein Drift-Wächter (`TestDryRun_GenerateEveryActionHasADryRunForm`) wird rot, wenn der Port eine neue `GenerateAction` benennt, die keine „Would …“-Form hat.

## 4. Trigger

**Start** (`next` → `in-progress`): Der Projektinhaber hat die offene Frage aus §1 entschieden (A, B oder C) und priorisiert den Slice.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Vereinheitlichung verlangt auch eine Änderung der Texte von `init` und `add` oder des JSON-Formats.
- `in-progress` → `open` (blockiert): Lösung C — der CR zur Zeitform ist noch nicht in `done/`.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit
Steering-Loop-Lerneintrag geschrieben (§7). Bei Lösung B: Begründung in §7, keine Liefer-Punkte.

## 6. Risiken und offene Punkte

- **Auslegung von [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) bleibt offen:** ohne Entscheidung des Projektinhabers ist unklar, ob Code oder Vertrag zu ändern ist. — **Ausgang:** entfallen: entschieden, Lösung A (Projektinhaber, 2026-10-17).
- **Textänderung ist für Skripte sichtbar:** wer die Dry-Run-Zeile maschinell las, sieht einen anderen Text. — **Ausgang:** entfallen: [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) weist `--json` als maschinenlesbaren Weg aus; der CHANGELOG nennt die Änderung.
- **Bestehende Pins schreiben die Vergangenheitsform fest:** die Matrix in `quiet_preview_test.go` <!-- d-check:ignore (Pfad-Zitat) --> und Bestandstests müssen mitgezogen werden. — **Ausgang:** entfallen: Der Plan (§3) nennt sie; der Implementer sucht die Texte per Grep im Test-Baum, bevor er sie ändert.

## 7. Closure-Notiz

**Entscheidung des Projektinhabers (2026-10-17):** Lösung A, „A: Code vereinheitlichen“. [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) bleibt unverändert (kein CR); der Code folgt dem Sinn der Zusage, dass der Dry-Run nichts als geschehen behauptet. Die Lösungen B und C entfallen.

**Belege:** Review-Report [`2026-10-17-slice-v1-cli-dry-run-would-formulierung`](../../../reviews/2026-10-17-slice-v1-cli-dry-run-would-formulierung.md) (0 HIGH, 0 MEDIUM, 1 LOW, 4 INFO; Mutationsläufe M1 bis M6, alle rot). Verification Evidence [`2026-10-17-verification-slice-v1-cli-dry-run-would-formulierung`](../../../reviews/2026-10-17-verification-slice-v1-cli-dry-run-would-formulierung.md): Liefer-Punkte erfüllt, keine Diskrepanz zwischen Plan, Code, Handbuch und CHANGELOG; Auflagen 1 bis 4 in dieser Closure erfüllt. `make gates` grün (Coverage 91,20 %), `make doc-trace`: 83 Anforderungen, 0 Waisen. Commits `73f1f4b`, `1536c8e`, `27723e2`.

**Fundstellen je Aussage (`AGENTS.md` §3.10):**

| Aussage | Fundstelle |
|---|---|
| `remove` sagt im Dry-Run „Would remove service …“ und „Would change:“; `--diff` allein bleibt Vergangenheitsform | `printRemoveSummary` (Kriterium `previewMode == driving.PreviewDryRun`) in `internal/adapter/driving/cli/remove.go` |
| `generate` sagt im Dry-Run „Would generate/update/repair …“ | `printGenerateSummary` (Parameter `dryRun`) in `internal/adapter/driving/cli/generate.go` |
| `config set` sagt im Dry-Run „Would set …“, einzeln und je Paar | `printConfigSetSummary` und `runConfigSetMany` (reicht `flags.DryRun` durch) in `internal/adapter/driving/cli/config.go` |
| Happy, Boundary, Negative je Subkommando | `TestDryRun_SummarySaysWould`, `TestDryRun_ActualWriteKeepsPastForm`, `TestDryRun_JSONAndErrorsUnchanged` in `internal/adapter/driving/cli/dry_run_would_test.go` |
| Je `GenerateAction`, Drift-Wächter für den `default`-Zweig, Mischfall `config set` | `TestDryRun_GenerateActions`, `TestDryRun_GenerateEveryActionHasADryRunForm`, `TestDryRun_ConfigSetManyMixedPairs` in `dry_run_would_test.go` |
| Die Dry-Run-Headline von `remove --purge` ist gepinnt | `TestRemove_PurgeDryRun_HumanMode_NoVolumeWarningPollution` in `internal/adapter/driving/cli/remove_acceptance_test.go` |
| Handbuch und CHANGELOG nennen die Formulierung, den Quiet-Zusatz und den Hinweis auf `--json` | Abschnitt „Vorschau statt Überraschung“ in `docs/user/benutzerhandbuch.md`; `[Unreleased]` · Changed in `CHANGELOG.md` |

- **Was hat funktioniert:** Die Mutationsläufe des Reviews (M1 bis M6) trafen jeden Zweig (Kriterium `--diff` allein, Mehrpaar-Durchreichung, `Repaired`-Zweig); der Praxislauf des Verifiers am Runtime-Image bestätigte alle Formen am Binary. Die Fix-Runde (`27723e2`) schloss F-1 bis F-4 vor der Verification.
- **Was ging anders als geplant:** Der Implementer führte einmal `python3` auf dem Host aus (Verstoß gegen `AGENTS.md` §3.1, Review F-5). Der Aufruf änderte nur Kommentar-Berichtigung und eine Headline-Assertion in `remove_acceptance_test.go`; der Nachweis läuft über `make test` und `make gates` im Docker; im Plan §3 vermerkt. §6 (erstes Risiko) und §7 standen nach der Lieferung noch auf „offen“ und wurden erst in der Closure nachgezogen (Verifier-Auflage 2).
- **Finding-Klassen des Reviews:** Doku-Aussage mehrdeutig zum Verbosity-Zusammenspiel (F-1, CHANGELOG-Satz zu `--quiet`; behoben) · Defensiver Zweig nicht dry-run-aware (F-2; Grenze in §3 benannt, Drift-Wächter) · Mehrpaar-Mischfall nicht gepinnt (F-3; behoben) · Teilstring-Negativ-Assert zu breit (F-4; behoben) · Host-Werkzeug außerhalb Docker (F-5, einmalig).
- **Steering-Loop-Eintrag:** **Gezählt, nicht verkörpert.** Kein `liegt in`, nichts verkörpert, keine Entscheidung über Verkörperung. (1) Neu im Register mit je einem Beleg: `BEO-WZG/host-werkzeug-ausserhalb-docker` (F-5; eine Sperre hat den Host-Aufruf nicht verhindert) und `BEO-CLI/dry-run-log-vergangenheitsform` (Hinweis des Verifiers). (2) Die Klassen F-1 bis F-4 sind einmalig und ohne bestehenden Pfad; sie stehen hier benannt, nicht im Register: ein Zähler 1 ohne Wiederauftreten sagt nichts, tritt eine Klasse in einem zweiten Vorgang auf, legt dessen Closure den Pfad an. (3) `BEO-DOK/plan-aussage-ueber-code` und `BEO-CLI/kommentar-beschreibt-abwesendes`: die Klassen traten in diesem Vorgang nicht ein (Plan-Tabelle gegen `git diff --stat` ohne Abweichung, Kommentar-Selbstcheck ohne Befund); Zähler unverändert. **Lehre:** Ein Risiko-Ausgang „weiter offen“, der auf eine Entscheidung wartet, wird mit der Entscheidung berichtigt, nicht erst in der Closure; ein einzelner Vorgang ändert den Stand keiner Beobachtung.
- **Beobachtungs-Register (`../observations/`):** zwei neue Verzeichnisse, je eine Evidence-Datei: [`BEO-WZG/host-werkzeug-ausserhalb-docker`](../observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md) (eine bestehende Beobachtung passt nicht: `BEO-WZG/kopplung-nicht-nachgezogen` meint ein Skript, das einer geänderten Form nicht folgt) und [`BEO-CLI/dry-run-log-vergangenheitsform`](../observations/BEO-CLI/dry-run-log-vergangenheitsform/observation.md). Beide stehen `offen`, unter der Schwelle von drei.
- **Folge-Slices:** keine. Die INFO-Logzeilen im Dry-Run bleiben im Register; ein Folge-Slice entsteht nur auf Entscheidung des Projektinhabers (ob Logzeilen im Dry-Run in der Vergangenheitsform stehen dürfen, ist eine Auslegung von [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) und [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)).
- **Risiken aus §6:** alle drei tragen einen Ausgang: entfallen (Auslegung entschieden), entfallen (Skript-Sichtbarkeit), entfallen (bestehende Pins).
- **Roadmap:** Der Ruhe-Marker „Nichts in Arbeit.“ steht nach dem Move wieder unter „Offene Wellen“ (`in-progress/` trägt keinen Slice; d-check `planning`).
- **Drei Paarungen:** Anker (kein `liegt in`, nichts verkörpert) · Folge-Slice (keiner) · Register (zwei Pfade, je mit nicht leerem `evidence/`; Belege liegen als `evidence/slice-v1-cli-dry-run-would-formulierung.md`) — nach dem `git mv` nach `done/` geprüft.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist `internal/adapter/driving/cli` (Sub-Area CLI laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006), Modus Greenfield, graduiert)
sowie `docs/plan/`, `CHANGELOG.md` und `docs/user/` als öffentliche Verträge (keine Sub-Area).

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register
([`../observations/`](../observations/README.md)) durchgegangen — keine Treffer
(Register leer).

Alle berührten Sub-Areas GF.
