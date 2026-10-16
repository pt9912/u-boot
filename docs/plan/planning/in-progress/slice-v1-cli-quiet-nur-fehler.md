# Slice v1: `--quiet` zeigt nur Fehler — Code und Sicherheitslockerungs-Ausnahme

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-cli-verbosity`](../welle-cli-verbosity.md) — folgt auf
[`slice-cr-cli-quiet-fehler-und-warnungen`](../done/slice-cr-cli-quiet-fehler-und-warnungen.md) (die Ausnahme steht erst nach
Lastenheft 0.8.0 im Vertrag) und geht
[`slice-v1-cli-verbosity-letzte-option-gewinnt`](../open/slice-v1-cli-verbosity-letzte-option-gewinnt.md) voraus (siehe §4: beide ändern
`root.go` und `verbosity_test.go`).

**Bezug:** [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Scope),
[`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil),
[`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer).
Herkunft: Entscheidung des Projektinhabers (2026-10-16), festgehalten im CR
[`slice-cr-cli-quiet-fehler-und-warnungen`](../done/slice-cr-cli-quiet-fehler-und-warnungen.md).

**Berührte Spec-Stellen:** `lastenheft.md` §3 ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)) — wird
nicht geändert; der Code folgt dem Lastenheft (Fassung aus dem CR, Lastenheft 0.8.0).
`spezifikation.md` und `architecture.md` beschreiben das Verhalten erst nach
allen drei Slices der Welle, im letzten
([`slice-v1-cli-verbosity-letzte-option-gewinnt`](../open/slice-v1-cli-verbosity-letzte-option-gewinnt.md)).

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

**Ziel:** `--quiet` zeigt nur Fehler: der Logger-Pegel ist Fehler (heute Warn),
Status- und Erfolgsmeldungen sowie Warnungen entfallen, und die Ausweisung von
Sandbox-Sicherheitslockerungen bleibt bei `init`, `generate` und `up` sichtbar.

**Gefundener Stand im Code (Ausgangslage):**

- `cli/root.go` setzt bei `--quiet` `slog.LevelWarn`. Produktionscode ruft
  `Logger.Warn` nirgends auf; der Pegelwechsel auf Fehler ändert die sichtbare
  Ausgabe deshalb nicht, bringt den Code aber auf den Wortlaut des Lastenhefts.
- Sandbox-Lockerungs-Hinweise (`sandboxRelaxationWarnings`, Codes
  [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)/[`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)) laufen als `WarningEntry` durch
  die Use-Case-Antwort und werden in `cli/generate.go` und `cli/init.go` ohne
  `--quiet`-Prüfung ausgegeben; sie bleiben dort sichtbar.
- `up` gibt keine Sandbox-Hinweise aus (der Up-Service erzeugt nur die
  Recreate-Warnung [`LH-FA-UP-003`](../../../../spec/lastenheft.md#lh-fa-up-003--startstatus-anzeigen)); `cli/up.go` springt bei `flags.Quiet` vor der
  Warnungs-Schleife zurück. Die Ausnahme läuft dort heute leer; der Slice sichert
  sie, indem `up` unter `--quiet` Lockerungs-Hinweise ausgibt und nur die übrigen
  Warnungen unterdrückt.
- Nicht jede Warnung aus der Sandbox-Erzeugung ist eine Lockerung: der
  „no git remote"-Hinweis (Code [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)) und die „has no effect"-Hinweise
  sind es nicht. Der Code unterscheidet sie heute nicht.
- `doctor --quiet` blendet nur OK-Items aus; Warn-Items bleiben sichtbar.
  `down --quiet` unterdrückt Erfolgsmeldung und Volume-Liste (bereits).

**Ansatz:** Der Pegel kommt aus `root.go`; die Ausgabe-Filter bleiben an
`a.quiet` gekoppelt. Ein Lockerungs-Hinweis wird am `WarningEntry` gekennzeichnet
(ein Feld, das nicht ins JSON geht), und ein gemeinsamer Ausgabe-Helfer für die
Warnungs-Schleifen in `init`, `generate` und `up` gibt unter `--quiet` nur
gekennzeichnete Hinweise aus. Der JSON-Modus bleibt unverändert (`--quiet --json`
ist ein No-op).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Rangfolge bei Mehrfachangabe** (`--quiet --verbose` usw.) — das ist
  [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../open/slice-v1-cli-verbosity-letzte-option-gewinnt.md); dieser Slice ändert, *was* `--quiet`
  bewirkt, nicht, *welche* Option gilt.
- **Eine neue Pegelstufe** (etwa zwischen Info und Debug) — eine neue Zusage;
  das Lastenheft nennt vier Stufen.
- **Die Sandbox-Hinweise auf den Logger umstellen** — ein anderer Vorgang:
  Die Hinweise sind Teil der Befehlsausgabe, die
  [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) verlangt; über den Logger liefen sie nach Pegel und nicht nach
  Pflicht.
- **Die Zusammenfassungen von `init` und `generate` unter `--quiet`** — die
  Entscheidung nennt Statustabelle von `up`, Erfolgsmeldung von `down` und
  OK-Items von `doctor`; ein `--quiet` für `init`/`generate` wäre ein neuer
  Vertrag und braucht eine eigene Entscheidung.
- **Andere Warn-Klassen als Sicherheitslockerungen** (Recreate-Warnung
  [`LH-FA-UP-003`](../../../../spec/lastenheft.md#lh-fa-up-003--startstatus-anzeigen), Drittanbieter-Features
  [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), „no git remote") — die Ausnahme ist genau eine; sie
  fallen unter `--quiet` weg.

## 2. Definition of Done

- [ ] [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) und [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) erfüllt: `--quiet` setzt den
      Fehler-Pegel; `up`, `down` und `doctor` unterdrücken Status-/Erfolgsmeldungen,
      `doctor --quiet` zeigt nur Fehler-Items; Sandbox-Lockerungs-Hinweise bleiben
      bei `init`, `generate` und `up` sichtbar. Tests Happy/Boundary/Negative je
      Subkommando referenzieren die IDs.
- [ ] Die Degradations-Fallbacks aus [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)
      (Container ohne Daemon, `vfs`-Storage, Wegfall der Egress-Restriktion) sind ausdrücklich der Ausnahme zugeordnet oder
      ausgeschlossen, und ein Test pinnt die Zuordnung (das Lastenheft nennt sie weder noch schließt es sie aus).
- [ ] Doku-Update: `docs/user/benutzerhandbuch.md` (Optionstabellen, Zeilen um
      299 und 1280) nennt `--quiet` als „nur Fehler" mit der Ausnahme;
      CHANGELOG-Eintrag.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo mit Wellen-Betrieb von der Welle-Closure geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapter/driving/cli/root.go` | update | `PersistentPreRunE`: `case a.quiet` setzt `slog.LevelError`; Flag-Hilfetext („keep warnings and errors where applicable") und Kommentare passen |
| `internal/adapter/driving/cli/cli.go` | update | Feldkommentar „--quiet lowers it to Warn" auf Fehler-Pegel |
| `internal/hexagon/port/driving/removeservice.go` | update | `WarningEntry` bekommt ein Kennzeichen für Sicherheitslockerungen (nicht im JSON); [`LH-FA-ARCH-003`](../../../../spec/lastenheft.md#lh-fa-arch-003--import-regeln-und-enforcement)-konform, kein I/O |
| `internal/hexagon/application/devcontainer_sandbox.go` | update | `sandboxRelaxationWarnings` kennzeichnet nur die „security relaxation"-Hinweise; „has no effect"-Hinweise und der „no remote"-Hinweis bleiben ungekennzeichnet |
| `internal/adapter/driving/cli/generate.go`, `init.go` | update | Warnungs-Schleifen über den gemeinsamen Helfer; unter `--quiet` nur gekennzeichnete Hinweise (heute ohne `--quiet`-Prüfung). Die Subkommandos lesen `a.quiet` bisher nicht; die Flag-Structs bekommen `Quiet` wie in `config.go` |
| `internal/adapter/driving/cli/up.go` | update | Rücksprung bei `flags.Quiet` unterdrückt Tabelle und Diagnosen, die Warnungs-Schleife läuft über den Helfer (Lockerungs-Hinweise sichtbar) |
| `internal/adapter/driving/cli/remove.go` | prüfen | schreibt Warnungen auf `errOut` (Purge-Warnungen); keine Sandbox-Lockerung, bleibt unberührt, sofern der Implementer keine gekennzeichneten Hinweise findet |
| `internal/adapter/driving/cli/doctor.go` | update | `writeDoctorReport`: unter `--quiet` nur Fehler-Items (heute nur OK ausgeblendet); Kopf- und Summary-Zeile gelten als Statusmeldung und entfallen — die Festlegung wird im Test gepinnt |
| `internal/adapter/driving/cli/down.go`, `statusview.go` | prüfen | unterdrücken Erfolgsmeldung und Volume-Liste bereits; Tests pinnen es |
| `internal/adapter/driving/cli/verbosity_test.go` | update | `TestExecute_VerbosityFlag_AdjustsLoggerLevel` pinnt `quiet-lowers-to-warn` → Fehler-Pegel; der Fall `debug-takes-precedence-over-quiet` bleibt in diesem Slice unverändert (er kehrt sich erst in [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../open/slice-v1-cli-verbosity-letzte-option-gewinnt.md) um) |
| `internal/adapter/driving/cli/sandbox_test.go`, `updown_cluster_test.go`, `statusview_test.go`, `down_acceptance_test.go`, `up_acceptance_test.go`, `doctor_json_test.go` | update / prüfen | vorhandene Pins (Warnzeilen der Sandbox- und Recreate-Hinweise ohne `--quiet`, `--quiet`-Unterdrückung von Tabelle/Erfolgsmeldung, `--quiet --json` No-op) bleiben grün; neue Fälle Happy/Boundary/Negative: `init`/`generate`/`up` mit `--quiet` zeigen den Lockerungs-Hinweis, unterdrücken „no remote" und Recreate-Warnung; `doctor --quiet` zeigt Fehler, keine Warn-/OK-Items; `down --quiet` leer; `--quiet --json` unverändert — nach [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) |
| `docs/user/benutzerhandbuch.md` | update | Optionstabellen (Zeilen um 299 und 1280): „nur Fehler; Sicherheitslockerungen im Sandbox-Profil bleiben ausgewiesen" |
| `CHANGELOG.md` | update | sichtbare Verhaltensänderung (`doctor --quiet` zeigt keine Warn-Items mehr; Warnungen unter `--quiet` entfallen außer Sicherheitslockerungen) |

## 4. Trigger

**Start** (`next` → `in-progress`): Lastenheft 0.8.0 ist angenommen (der CR
[`slice-cr-cli-quiet-fehler-und-warnungen`](../done/slice-cr-cli-quiet-fehler-und-warnungen.md) liegt in `done/`) und der Projektinhaber
priorisiert den Slice.

**Reihenfolge zum Slice [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../open/slice-v1-cli-verbosity-letzte-option-gewinnt.md):**
logisch sind beide unabhängig — dieser Slice ändert die Abbildung `--quiet` →
Pegel und die Ausgabe-Filter, jener die Auswahl unter mehreren Optionen. Sie
berühren aber dieselbe Stelle (`PersistentPreRunE`-Switch in `root.go`) und
dieselbe Testfunktion in `verbosity_test.go`. Die Reihenfolge ist deshalb eine
Konflikt-Vermeidung, keine Voraussetzung; läuft jener zuerst, ändert sich nur
die Zeile `case`-Pegel, die hier gesetzt wird.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Kennzeichnung am
  `WarningEntry` zieht weitere Use-Cases (`add`, `remove`) nach sich — dann die
  Kennzeichnung als eigener Slice vorziehen.
- `in-progress` → `open` (blockiert): Lastenheft 0.8.0 ist nicht angenommen oder
  der Projektinhaber ändert die Ausnahme.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit
Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **Abgrenzung „Sicherheitslockerung":** Die Kennzeichnung am `WarningEntry` kann
  eine Lockerung übersehen oder einen Hinweis zu Unrecht ausnehmen. — **Ausgang:**
  entfallen: der Test pinnt alle „security relaxation"-Hinweise aus
  `sandboxRelaxationWarnings` (Podman, Docker, Egress) und die ausgenommenen
  Hinweise einzeln.
- **`doctor --quiet` verliert Warn-Items:** Skripte, die die Warn-Zeilen lasen,
  sehen sie nicht mehr; der Exit-Code trägt das Ergebnis. — **Ausgang:**
  entfallen: die Folge ist gewollt (Entscheidung des Projektinhabers) und
  steht im CHANGELOG und in den Optionstabellen des Benutzerhandbuchs (DoD).
- **Reihenfolge zum Folge-Slice:** beide ändern `root.go` und
  `verbosity_test.go`. — **Ausgang:** entfallen: logisch unabhängig (§4); ein
  Merge-Konflikt wird beim zweiten Slice aufgelöst.

## 7. Closure-Notiz

Wird bei der Closure gefüllt.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `internal/adapter/driving/cli`
(Kürzel CLI), `internal/hexagon/port` (PORT) und `internal/hexagon/application`
(APP) — je Modus Greenfield laut
[`harness/conventions.md`](../../../../harness/conventions.md#mr-006) (Abschnitt Modus-Deklaration), jede erfüllt die Schwelle
≥ 2 von 3 Achsen; keine Ausdifferenzierung nötig. `docs/user/` und
`CHANGELOG.md` sind öffentliche Verträge, keine Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register
([`../observations/`](../observations/README.md)) durchgegangen — keine Treffer
(Register leer).

Alle berührten Sub-Areas GF.
