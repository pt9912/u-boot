# Slice v1: Verbosity — die zuletzt angegebene Option gewinnt

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-cli-verbosity`](../welle-cli-verbosity.md) — letzter Slice der Welle; folgt auf
[`slice-cr-cli-quiet-fehler-und-warnungen`](../in-progress/slice-cr-cli-quiet-fehler-und-warnungen.md) und
[`slice-v1-cli-quiet-nur-fehler`](slice-v1-cli-quiet-nur-fehler.md), weil die Pegelabbildung von
`--quiet` (Fehler-Pegel) aus dem Slice [`slice-v1-cli-quiet-nur-fehler`](slice-v1-cli-quiet-nur-fehler.md) kommt.

**Bezug:** [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Scope).
Verwandt: [`slice-cr-cli-quiet-fehler-und-warnungen`](../in-progress/slice-cr-cli-quiet-fehler-und-warnungen.md)
(klärt die Zusage von `--quiet`) und [`slice-v1-cli-quiet-nur-fehler`](slice-v1-cli-quiet-nur-fehler.md)
(setzt sie im Code um; dieser Slice ändert nur die Wahl zwischen mehreren Optionen und ist logisch davon unabhängig).
Herkunft: Inventur in
[`slice-harness-internal-readmes-festlegungen-in-spec`](../done/slice-harness-internal-readmes-festlegungen-in-spec.md) §9.

**Berührte Spec-Stellen:** `lastenheft.md` §3 ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)) — wird
nicht geändert; der Code folgt dem Lastenheft. `spezifikation.md` und
`architecture.md` beschreiben das Verhalten erst **nach** diesem Slice, dem Slice
[`slice-v1-cli-quiet-nur-fehler`](slice-v1-cli-quiet-nur-fehler.md) **und**
dem CR [`slice-cr-cli-quiet-fehler-und-warnungen`](../in-progress/slice-cr-cli-quiet-fehler-und-warnungen.md);
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
  Sicherheitslockerungen) — das ist der Slice [`slice-v1-cli-quiet-nur-fehler`](slice-v1-cli-quiet-nur-fehler.md)
  (Zusage im CR [`slice-cr-cli-quiet-fehler-und-warnungen`](../in-progress/slice-cr-cli-quiet-fehler-und-warnungen.md)); dieser Slice ändert nur,
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
  und `down` unterdrücken Tabelle/Erfolgsmeldung, Sandbox-Lockerungs-Hinweise
  bleiben sichtbar; `--quiet --json` ist No-op) — sie kommen aus
  [`slice-v1-cli-quiet-nur-fehler`](slice-v1-cli-quiet-nur-fehler.md), lesen `a.quiet` und bleiben unverändert; Folge: `a.quiet` muss bei
  „letzte gewinnt" weiter genau dann wahr sein, wenn `--quiet` die gültige
  Option ist (siehe §6).

## 2. Definition of Done

- [ ] [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) erfüllt: die zuletzt angegebene Verbosity-Option gewinnt;
      Tests Happy/Boundary/Negative referenzieren die ID (`--quiet --verbose`,
      `--verbose --quiet`, `--debug --quiet`, `--quiet --debug`, Wiederholung
      wie `--quiet --verbose --quiet`, Einzel-Optionen, keine Option).
- [ ] Doku-Update: `docs/user/benutzerhandbuch.md` (Optionstabellen, Zeilen
      um 299 und 1280) nennt „letzte Option gewinnt"; CHANGELOG-Eintrag.
- [ ] Spezifikation: ein `SPEC-<NNN>` in `spezifikation.md` beschreibt Pegelabbildung und Rangfolge, deckungsgleich mit
      Lastenheft 0.8.0 (nach dem CR) und Code, an [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) verankert; Zahlenwerte aus dem Code.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne Wellen-Betrieb hier geprüft.

## 3. Plan (vor Code)

Stand der Rangfolge heute (im Code und in der Doku gefunden):

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapter/driving/cli/root.go` | update | Flag-Definition (`BoolVar` für `quiet`/`verbose`/`debug`) und `PersistentPreRunE`-Switch mit fester Rangfolge; Kommentar „Precedence: --debug > --verbose > --quiet" ist zu ersetzen |
| `internal/adapter/driving/cli/cli.go` | update | Felder `quiet`/`verbose`/`debug` am `App` (Kommentar „--quiet lowers it to Warn"); Lesestellen `a.quiet` in `config.go`, `up.go`, `down.go`, `doctor.go` müssen mit dem neuen Wert konsistent bleiben |
| `internal/adapter/driving/cli/verbosity_test.go` | update | `TestExecute_VerbosityFlag_AdjustsLoggerLevel` pinnt heute `debug-takes-precedence-over-quiet` (`--debug --quiet` → Debug) — dieser Fall kehrt sich um; die Zeile `quiet-lowers-to-warn` ist bereits durch den Slice [`slice-v1-cli-quiet-nur-fehler`](slice-v1-cli-quiet-nur-fehler.md) auf Fehler-Pegel gesetzt; neue Fälle Happy/Boundary/Negative nach [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) |
| `internal/adapter/driving/cli/doctor_json_test.go`, `up_acceptance_test.go`, `logs_acceptance_test.go` | prüfen | verwenden `--quiet --json`/`--json --quiet`; unberührt, sofern `a.quiet` für Einzel-`--quiet` gleich bleibt |
| `docs/user/benutzerhandbuch.md` | update | Optionstabellen (zwei Stellen) nennen die Optionen, aber nicht die Rangfolge; Satz „letzte gewinnt" ergänzen |
| `CHANGELOG.md` | update | sichtbare Verhaltensänderung (`--debug --quiet` wird zu `--quiet`, also Fehler-Pegel) |

Weitere Fundstellen (nicht zu ändern): `docs/maintainer/cli-json-contract.md`
(`--quiet --json` No-op) und `internal/adapter/driven/logger/logger.go`
(LevelVar-Beschreibung).

## 4. Trigger

**Start** (`next` → `in-progress`): Priorisierung durch den Projektinhaber.

**Rückführungen — vorab benennen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die gemeinsame
  Flag-Wert-Lösung zwingt zu einem Umbau der `a.quiet`-Lesestellen über mehrere
  Subkommandos hinaus — dann Umbau als eigener Slice vorziehen.
- `in-progress` → `open` (blockiert): der Slice [`slice-v1-cli-quiet-nur-fehler`](slice-v1-cli-quiet-nur-fehler.md)
  ist nicht abgeschlossen oder die Pegelabbildung von `--quiet` ändert sich
  durch eine neue Entscheidung des Projektinhabers, sodass die Umsetzung hier
  berührt ist.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit
Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **`a.quiet` driftet vom Pegel:** Nach „letzte gewinnt" kann `--quiet --verbose`
  den Pegel auf Debug setzen, während `a.quiet` (Filter in `doctor`/`up`/`down`)
  noch wahr ist. — **Ausgang:** offen bis zur Umsetzung; Test `--quiet --verbose`
  muss auch die Ausgabe-Filter prüfen.
- **Reihenfolge-Abhängigkeit:** Die Spezifikation darf das Verhalten erst nach
  diesem Slice, dem Slice [`slice-v1-cli-quiet-nur-fehler`](slice-v1-cli-quiet-nur-fehler.md) und dem CR
  [`slice-cr-cli-quiet-fehler-und-warnungen`](../in-progress/slice-cr-cli-quiet-fehler-und-warnungen.md) beschreiben. — **Ausgang:** entfallen: der
  Spezifikations-Eintrag liegt in diesem Slice, dem letzten der Welle
  [`welle-cli-verbosity`](../welle-cli-verbosity.md).

## 7. Closure-Notiz

Wird bei der Closure gefüllt.

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
