# Slice v1: Verbosity — die zuletzt angegebene Option gewinnt

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-cli-verbosity`](welle-cli-verbosity.md) — letzter Slice der Welle; folgt auf
[`slice-cr-cli-quiet-fehler-und-warnungen`](../done/slice-cr-cli-quiet-fehler-und-warnungen.md),
[`slice-cr-cli-quiet-weitere-ausnahmen`](../done/slice-cr-cli-quiet-weitere-ausnahmen.md),
[`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md),
[`slice-cr-cli-quiet-vorschau`](../done/slice-cr-cli-quiet-vorschau.md) und
[`slice-v1-cli-quiet-vorschau`](../done/slice-v1-cli-quiet-vorschau.md), weil die Pegelabbildung von
`--quiet` (Fehler-Pegel) aus dem Slice [`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md) kommt.

**Bezug:** [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Scope).
Verwandt: [`slice-cr-cli-quiet-fehler-und-warnungen`](../done/slice-cr-cli-quiet-fehler-und-warnungen.md)
(klärt die Zusage von `--quiet`) und [`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md)
(setzt sie im Code um; dieser Slice ändert nur die Wahl zwischen mehreren Optionen und ist logisch davon unabhängig).
Herkunft: Inventur in
[`slice-harness-internal-readmes-festlegungen-in-spec`](../done/slice-harness-internal-readmes-festlegungen-in-spec.md) §9.

**Berührte Spec-Stellen:** `lastenheft.md` §3 ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)) — wird
nicht geändert; der Code folgt dem Lastenheft. `spezifikation.md` und
`architecture.md` beschreiben das Verhalten erst **nach** diesem Slice, dem Slice
[`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md) **und**
den CRs [`slice-cr-cli-quiet-fehler-und-warnungen`](../done/slice-cr-cli-quiet-fehler-und-warnungen.md),
[`slice-cr-cli-quiet-weitere-ausnahmen`](../done/slice-cr-cli-quiet-weitere-ausnahmen.md) und
[`slice-cr-cli-quiet-vorschau`](../done/slice-cr-cli-quiet-vorschau.md) sowie dem Slice
[`slice-v1-cli-quiet-vorschau`](../done/slice-v1-cli-quiet-vorschau.md);
vorher nicht in die Spec schreiben.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

**Ziel:** Bei mehreren Verbosity-Optionen (`--quiet`, `--verbose`, `--debug`)
bestimmt die zuletzt auf der Kommandozeile angegebene den Pegel, ohne
Abweisung — wie es [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) zusagt. Der Code
(`internal/adapter/driving/cli/root.go`) wendet heute eine feste Rangfolge
`--debug` > `--verbose` > `--quiet` an.

**Warum Code statt Lastenheft:** Die Regel „letzte gewinnt" erlaubt, einen
Alias zu überstimmen (`alias ub='u-boot --quiet'`, dann `ub --verbose`); eine
feste Rangfolge macht das unmöglich. Der Projektinhaber hat entschieden, den
Code anzupassen.

**Ansatz:** Cobra-Bool-Flags kennen die Reihenfolge nicht. Die drei Optionen
teilen sich einen pflag-Wert, der den zuletzt gesetzten Pegel festhält
(Optionen bleiben einzeln als `--quiet`/`--verbose`/`--debug` ohne Wert
benutzbar).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Änderung der Pegelabbildung von `--quiet`** (Fehler-Pegel, Ausnahme für
  Sicherheitslockerungen) — das ist der Slice [`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md)
  (Zusage im CR [`slice-cr-cli-quiet-fehler-und-warnungen`](../done/slice-cr-cli-quiet-fehler-und-warnungen.md)); dieser Slice ändert nur,
  *welche* Option gilt, nicht, *was* sie bewirkt.
- **Neue Stufe zwischen Info und Debug** — `--verbose` und `--debug` bleiben
  beide `Debug`; der Kommentar in `root.go` nennt das als vorläufig. Eine neue
  Stufe wäre eine neue Zusage und gehört in einen eigenen Slice, wenn ein
  Dienst sie braucht.
- **Konfigurations- und Umgebungsquellen für den Pegel** — im Code gibt es
  keine (geprüft: `grep` nach `Getenv`/`LookupEnv` in `internal/` und `cmd/`
  ohne Treffer; der Pegel kommt nur aus den drei Flags). Es gibt also nichts
  zu priorisieren; eine solche Quelle wäre ein neuer Vertrag.
- **Die zusätzlichen `--quiet`-Wirkungen** (`doctor` zeigt nur Fehler-Items, `up`
  und `down` unterdrücken Tabelle/Erfolgsmeldung, die vier Ausnahmen aus
  [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) bleiben sichtbar; `--quiet --json` ist No-op) — sie kommen aus
  [`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md), lesen `a.quiet` und bleiben unverändert; Folge: `a.quiet` muss bei
  „letzte gewinnt" weiter genau dann wahr sein, wenn `--quiet` die gültige
  Option ist (siehe §6).

## 2. Definition of Done

- [x] [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) erfüllt: die zuletzt angegebene Verbosity-Option gewinnt;
      Tests Happy/Boundary/Negative referenzieren die ID (`--quiet --verbose`,
      `--verbose --quiet`, `--debug --quiet`, `--quiet --debug`, Wiederholung
      wie `--quiet --verbose --quiet`, Einzel-Optionen, keine Option).
- [x] Doku-Update: `docs/user/benutzerhandbuch.md` (Optionstabellen, Zeilen
      um 299 und 1280) nennt „letzte Option gewinnt"; CHANGELOG-Eintrag.
- [x] Spezifikation: ein `SPEC-<NNN>` in `spezifikation.md` beschreibt Pegelabbildung und Rangfolge, deckungsgleich mit
      Lastenheft 0.10.0 (nach den CRs) und Code, an [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) verankert; Zahlenwerte aus dem Code.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne Wellen-Betrieb hier geprüft.

## 3. Plan (vor Code)

Stand der Rangfolge heute (im Code und in der Doku gefunden):

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapter/driving/cli/root.go` | update | Flag-Definition (`BoolVar` für `quiet`/`verbose`/`debug`) und `PersistentPreRunE`-Switch mit fester Rangfolge; Kommentar „Precedence: --debug > --verbose > --quiet" ist zu ersetzen |
| `internal/adapter/driving/cli/cli.go` | update | Felder `quiet`/`verbose`/`debug` am `App` (Kommentar „--quiet lowers it to Warn"); Lesestellen `a.quiet` in `config.go`, `up.go`, `down.go`, `doctor.go` müssen mit dem neuen Wert konsistent bleiben |
| `internal/adapter/driving/cli/verbosity_test.go` | update | `TestExecute_VerbosityFlag_AdjustsLoggerLevel` pinnt heute `debug-takes-precedence-over-quiet` (`--debug --quiet` → Debug) — dieser Fall kehrt sich um; die Zeile `quiet-lowers-to-warn` ist bereits durch den Slice [`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md) auf Fehler-Pegel gesetzt; neue Fälle Happy/Boundary/Negative nach [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) |
| `internal/adapter/driving/cli/doctor_json_test.go`, `up_acceptance_test.go`, `logs_acceptance_test.go` | prüfen | verwenden `--quiet --json`/`--json --quiet`; unberührt, sofern `a.quiet` für Einzel-`--quiet` gleich bleibt |
| `internal/adapter/driving/cli/verbosity.go` <!-- d-check:ignore (Pfad-Zitat) --> | neu | der gemeinsame pflag-Wert (`verbosityFlag`) und der abgeleitete Pegel; in der Plan-Tabelle zunächst nicht genannt (Ansatz in §1) |
| `internal/adapter/driving/cli/verbosity_last_wins_test.go` <!-- d-check:ignore (Pfad-Zitat) --> | neu | `LastWins_*`-Tests: Wirkung auf Zusammenfassung, Warnungsfilter, `doctor`-Bericht und `down`-Strom; Zustand je `Execute`; ungültiger Wert |
| `cmd/uboot/main.go` <!-- d-check:ignore (Pfad-Zitat) --> | update | Wiring-Kommentar nannte für `--quiet` den Pegel Warn (Review F-2); in der Plan-Tabelle nicht genannt |
| `spec/spezifikation.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) beschreibt Pegelabbildung, Rangfolge und Wirkung von `--quiet` (DoD §2) |
| `docs/user/benutzerhandbuch.md` | update | Optionstabellen (zwei Stellen) nennen die Optionen, aber nicht die Rangfolge; Satz „letzte gewinnt" ergänzen |
| `CHANGELOG.md` | update | sichtbare Verhaltensänderung (`--debug --quiet` wird zu `--quiet`, also Fehler-Pegel) |

Die Lesestellen (elf Dateien, zwölf Stellen) lesen jetzt `a.quiet()` statt des Feldes `a.quiet`: mechanisch auf `a.quiet()` umgestellt, weil das Feld
mit dem gemeinsamen Zustand entfällt; kein Verhalten wurde dabei geändert.

Weitere Fundstellen (nicht zu ändern): `docs/maintainer/cli-json-contract.md`
(`--quiet --json` No-op) und `internal/adapter/driven/logger/logger.go`
(LevelVar-Beschreibung).

## 4. Trigger

**Start** (`next` → `in-progress`): Priorisierung durch den Projektinhaber.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die gemeinsame
  Flag-Wert-Lösung zwingt zu einem Umbau der `a.quiet`-Lesestellen über mehrere
  Subkommandos hinaus — dann Umbau als eigener Slice vorziehen.
- `in-progress` → `open` (blockiert): der Slice [`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md)
  ist nicht abgeschlossen oder die Pegelabbildung von `--quiet` ändert sich
  durch eine neue Entscheidung des Projektinhabers, sodass die Umsetzung hier
  berührt ist.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit
Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **`a.quiet` driftet vom Pegel:** Nach „letzte gewinnt" kann `--quiet --verbose`
  den Pegel auf Debug setzen, während `a.quiet` (Filter in `doctor`/`up`/`down`)
  noch wahr ist. — **Ausgang:** entfallen: ein gemeinsamer Zustand, `a.quiet()`
  daraus abgeleitet; die `LastWins_*`-Tests belegen es (Zusammenfassung, Warnungen,
  `doctor`, `down`-Strom).
- **Reihenfolge-Abhängigkeit:** Die Spezifikation darf das Verhalten erst nach
  diesem Slice, den Slices [`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md) und
  [`slice-v1-cli-quiet-vorschau`](../done/slice-v1-cli-quiet-vorschau.md) und den CRs
  [`slice-cr-cli-quiet-fehler-und-warnungen`](../done/slice-cr-cli-quiet-fehler-und-warnungen.md),
  [`slice-cr-cli-quiet-weitere-ausnahmen`](../done/slice-cr-cli-quiet-weitere-ausnahmen.md) und
  [`slice-cr-cli-quiet-vorschau`](../done/slice-cr-cli-quiet-vorschau.md) beschreiben. — **Ausgang:** entfallen: der
  Spezifikations-Eintrag liegt in diesem Slice, dem letzten der Welle
  [`welle-cli-verbosity`](welle-cli-verbosity.md).

## 7. Closure-Notiz

**Belege:** Review-Report
[`2026-10-16-slice-v1-cli-verbosity-letzte-option-gewinnt`](../../../reviews/2026-10-16-slice-v1-cli-verbosity-letzte-option-gewinnt.md) (1 HIGH, 0 MEDIUM, 1 LOW, 2 INFO; F-1 und F-2 behoben in `4260372`,
F-3 und F-4 als Grenzen). Verification Evidence
[`2026-10-16-verification-slice-v1-cli-verbosity-letzte-option-gewinnt`](../../../reviews/2026-10-16-verification-slice-v1-cli-verbosity-letzte-option-gewinnt.md): bestanden, mit Auflagen
für die Closure (hier eingelöst: S1 bis S3; W1 bis W4 gehen an die Welle-Closure). `make gates` grün, Coverage 91,10 % (Schwelle 90 %); `make doc-trace`: 83 Anforderungen, 0 Waisen.
Mutationen: drei des Implementers, drei des Reviewers und zwei des Verifiers fielen rot. Spezifikation: [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet),
an [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) verankert.

- **Was hat funktioniert:** Ein gemeinsamer Zustand statt vier Felder: `a.quiet()` wird aus ihm abgeleitet, das Risiko eines Auseinanderlaufens von Pegel und Ausgabe-Filter entfällt strukturell.
  Die Spezifikation entstand nach dem Code und wurde gegen den Code gelesen; [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet)
  stimmt Aussage für Aussage mit Code und Lastenheft 0.10.0 überein (Verification Prüfung 1). Die Lehre aus der ersten Welle hat gegriffen: Der Spezifikations-Eintrag kam erst nach den fünf anderen Slices.
- **Was ging anders als geplant:** Zwei Kommentare beschrieben den abgelösten Zustand: im Test „the former fixed order … is gone" (Review F-1, HIGH) und in `cmd/uboot/main.go` „--quiet lowers to Warn" (F-2);
  beide behoben in `4260372`. Das Beobachten von `--flag=false` und `--quiet=maybe` ist ein Verhalten, das der Plan nicht nannte. §3 nannte `verbosity.go`, `verbosity_last_wins_test.go`, `cmd/uboot/main.go`
  und `spec/spezifikation.md` nicht; berichtigt. <!-- d-check:ignore (Pfad-Zitate) -->
- **Grenzen:** (1) F-3: `--flag=false` (`--quiet --verbose=false` bleibt `--quiet`) ist nur im Code und Test festgelegt, nicht in [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) oder im Handbuch; der ungültige Wert (`--quiet=maybe`) ist nur als
  `err != nil` gepinnt, nicht als Exit-Code. (2) F-4: `--quiet doctor --verbose` (Option nach dem Subkommando) ist nur im Logger-Pegel-Test gepinnt, nicht in den Lesestellen-Tests; ein gemeinsamer Zustand macht
  eine Abweichung unwahrscheinlich.
- **Steering-Loop-Eintrag:** In dieser Welle gab es vier Kommentar-Befunde nach `AGENTS.md` §3.7 (Kommentar beschreibt nicht, was da ist): F-1 von
  [`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md), F-1 von [`slice-v1-cli-quiet-vorschau`](../done/slice-v1-cli-quiet-vorschau.md) sowie F-1 und F-2 dieses Slice. Die 3×-Schwelle ist erreicht.
  Die Anlage der Beobachtung `BEO-CLI/…` und die Verkörperung entscheidet die Welle-Closure (Schritt 3), nicht dieser Slice; hier wird nichts verkörpert, daher kein `liegt in`.
- **Beobachtungs-Register (`../observations/`):** nicht fortgeschrieben (das Register enthält nur die README). Belege für die spätere Beobachtung: dieselben vier Befunde (Review-Reports der drei Slices,
  Verification Evidence W2); die Anlage folgt bei der Welle-Closure.
- **Folge-Slices:** keine aus diesem Slice. [`slice-v1-sandbox-fallback-ausweisung`](../open/slice-v1-sandbox-fallback-ausweisung.md) und
  [`slice-v1-cli-dry-run-would-formulierung`](../done/slice-v1-cli-dry-run-would-formulierung.md) bleiben unberührt in `open/`.
- **Release-Hinweis:** Der CHANGELOG ist vollständig (nichts steht mehr als „folgt"); der Release-Cut ist nicht mehr blockiert.
- **Risiken aus §6:** alle tragen ihren Ausgang: entfallen (`a.quiet` driftet: ein Zustand, `a.quiet()` abgeleitet; Reihenfolge-Abhängigkeit: Spec nach den fünf anderen Slices).
- **Drei Paarungen:** Anker (kein `liegt in`, nichts verkörpert) · Folge-Slice (keine neuen; zwei Bestands-Slices unberührt in `open/`) · Register (nicht fortgeschrieben, Anlage bei der Welle-Closure) — benannt;
  Bestandseintrag in `open/README.md` bereinigt. Die Paarungen werden erst nach dem `git mv` abschließend getragen, daher bleibt das DoD-Häkchen dafür offen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist `internal/adapter/driving/cli`
(Kürzel CLI, Modus Greenfield laut
[`harness/conventions.md`](../../../../harness/conventions.md#mr-006), Abschnitt Modus-Deklaration) — erfüllt die Schwelle
≥ 2 von 3 Achsen; keine Ausdifferenzierung nötig. `docs/user/` und
`CHANGELOG.md` sind öffentliche Verträge, keine Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register
([`../observations/`](../observations/README.md)) durchgegangen — keine Treffer
(Register leer).

Alle berührten Sub-Areas GF.
