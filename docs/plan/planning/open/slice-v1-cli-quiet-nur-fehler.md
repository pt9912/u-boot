# Slice v1: `--quiet` zeigt nur Fehler — Code für alle Subkommandos und vier Ausnahmen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-cli-verbosity`](../welle-cli-verbosity.md) — dritter Slice; folgt auf
[`slice-cr-cli-quiet-fehler-und-warnungen`](../done/slice-cr-cli-quiet-fehler-und-warnungen.md) und
[`slice-cr-cli-quiet-weitere-ausnahmen`](../in-progress/slice-cr-cli-quiet-weitere-ausnahmen.md) (die vier Ausnahmen stehen erst nach
Lastenheft 0.9.0 im Vertrag) und geht
[`slice-v1-cli-verbosity-letzte-option-gewinnt`](../open/slice-v1-cli-verbosity-letzte-option-gewinnt.md) voraus (siehe §4: beide ändern
`root.go` und `verbosity_test.go`).

**Bezug:** [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Scope),
[`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil),
[`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer),
[`LH-FA-ADD-007`](../../../../spec/lastenheft.md#lh-fa-add-007--service-entfernen),
[`LH-NFA-PERF-002`](../../../../spec/lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker),
[`LH-FA-DIAG-003`](../../../../spec/lastenheft.md#lh-fa-diag-003--fehlerklassifikation).
Herkunft: Entscheidungen des Projektinhabers (2026-10-16), festgehalten in den CRs
[`slice-cr-cli-quiet-fehler-und-warnungen`](../done/slice-cr-cli-quiet-fehler-und-warnungen.md) und
[`slice-cr-cli-quiet-weitere-ausnahmen`](../in-progress/slice-cr-cli-quiet-weitere-ausnahmen.md); Befunde des Reviews
[`2026-10-16-slice-v1-cli-quiet-nur-fehler`](../../../reviews/2026-10-16-slice-v1-cli-quiet-nur-fehler.md) (F-1 bis F-6).

**Berührte Spec-Stellen:** `lastenheft.md` §3 ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)) — wird
nicht geändert; der Code folgt dem Lastenheft (Fassung aus dem CR, Lastenheft 0.9.0).
`spezifikation.md` und `architecture.md` beschreiben das Verhalten erst nach
allen Slices der Welle, im letzten
([`slice-v1-cli-verbosity-letzte-option-gewinnt`](../open/slice-v1-cli-verbosity-letzte-option-gewinnt.md)).

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

**Ziel:** `--quiet` zeigt nur Fehler: der Logger-Pegel ist Fehler, die Erfolgs- und Zusammenfassungsmeldung **aller**
Subkommandos, die etwas ändern, entfällt (`init`, `generate`, `add`, `remove`, `config set`, `up`, `down`; `doctor` zeigt nur
Fehler-Einträge), `remove` bekommt die Option, die Abfragebefehle (`config get`, `config list`, `logs`, `template list`)
behalten ihre Ausgabe, und vier Ausnahmen bleiben sichtbar: Sicherheitslockerungen, Degradations-Fallbacks
([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)), der Compose-Fortschritt von `up`
([`LH-NFA-PERF-002`](../../../../spec/lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker)) und die Pfad-Zusammenfassung von `--force`
([`LH-FA-CLI-005A`](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung),
[`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz)).

**Gefundener Stand im Code (Ausgangslage):**

- Der Slice ist teilweise umgesetzt (Commits `c277c47`, `72006d3`, `7a244b9`): Fehler-Pegel in `cli/root.go`,
  `driving.WarningEntry` mit Kennzeichen für Sicherheitslockerungen, `cli/warnings.go` (`visibleWarnings`), `doctor --quiet`
  zeigt nur Fehler-Items, `up` und `down` unterdrücken Tabelle bzw. Erfolgsmeldung. Offen sind die Befunde des Reviews und
  die Erweiterung der Reichweite.
- `init` und `generate` schreiben ihre Zusammenfassung (`printInitSummary`, Generate-Zusammenfassung) auch unter `--quiet`;
  `add` schreibt `printAddSummary`; `remove` schreibt „Removed service …“, „Changed:“, „Removed volumes:“ auf stdout und die
  Purge-Warnungen ([`LH-FA-ADD-007`](../../../../spec/lastenheft.md#lh-fa-add-007--service-entfernen)) auf stderr und liest `a.quiet` nicht (keine Option `Quiet` in seinen Flags).
- `config set` schreibt eine Zusammenfassungszeile (`printConfigSetSummary`), `config get`/`config list`, `logs` und `template list`
  geben Werte aus, die der Zweck des Befehls sind.
- Die Degradations-Fallbacks entstehen beim Containerstart (generierte Skripte) und als Items in `doctor`
  (`internal/hexagon/application/doctor_sandbox.go`: `/dev/fuse` → `vfs`, Egress-Capability nicht gewährbar, Daemon); `doctor`
  zeigt unter `--quiet` nur Fehler-Items, die Fallback-Items der Stufe Warn entfielen also. Nicht jedes Warn-Item von `doctor` ist ein
  Fallback (z. B. die Credentials-Prüfung).
- Der Compose-Fortschritt von `up` läuft auf stderr und ist von `--quiet` unberührt.

**Ansatz:** Der Pegel kommt aus `root.go`; die Ausgabe-Filter bleiben an `a.quiet` gekoppelt. Zusammenfassungen und
Erfolgsmeldungen der Subkommandos werden unter `--quiet` nicht geschrieben (jede Subkommando-Flag-Struct bekommt `Quiet`, wie
in `config.go`). Ein Hinweis, der eine Ausnahme trägt, wird am Eintrag gekennzeichnet (heute `SecurityRelaxation`); Fallback-Items
von `doctor` werden analog gekennzeichnet, damit `doctor --quiet` sie neben den Fehler-Items zeigt. Der JSON-Modus bleibt
unverändert (`--quiet --json` ist ein No-op).

**Prüfpunkt — Ergebnis-Ausgaben (das Lastenheft entscheidet, der Code setzt um und pinnt):** Die Abfragebefehle, die nichts
verändern, behalten ihre Ausgabe, weil sie das Ergebnis des Befehls ist: `config get`, `config list` (Werte), `logs`
(Log-Strom), `template list` (Tabelle). `config set` ändert etwas und gehört zu den ändernden Subkommandos: Die Zeile
„config: … → …“ (`printConfigSetSummary`) entfällt unter `--quiet`, nur Fehler bleiben. Test und Handbuch pinnen die Zuordnung.

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
- **Die Ausgabe der generierten Container-Start-Skripte** — sie laufen im Container, nicht im CLI-Prozess, und kennen `--quiet`
  nicht; der Slice sichert die Ausweisung in Befehlsausgabe und `doctor`, wie der Lastenheft-Satz sie nennt.
- **Andere Hinweis-Klassen als die vier Ausnahmen** (Recreate-Warnung
  [`LH-FA-UP-003`](../../../../spec/lastenheft.md#lh-fa-up-003--startstatus-anzeigen), Drittanbieter-Features
  [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), „no git remote“, Purge-Warnungen von `remove`)
  — sie entfallen unter `--quiet`, soweit sie keine Fehler sind.

## 2. Definition of Done

- [ ] [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) erfüllt (Lastenheft 0.9.0): `--quiet` setzt den Fehler-Pegel; `init`, `generate`, `add`, `remove`,
      `config set`, `up` und `down` unterdrücken Erfolgs- und Zusammenfassungsmeldungen, `doctor --quiet` zeigt nur
      Fehler-Einträge und Fallback-Items (Ausnahme 2); `remove` hat die Option `--quiet`; `config get`, `config list`,
      `logs` und `template list` behalten ihre Ausgabe (Prüfpunkt §1, gepinnt). Tests Happy/Boundary/Negative je Subkommando
      referenzieren die IDs.
- [ ] Die Purge-Warnungen und die Abhängigkeits-Warnung von `remove`
      ([`LH-FA-ADD-007`](../../../../spec/lastenheft.md#lh-fa-add-007--service-entfernen)) entfallen unter `--quiet`, soweit sie keine Fehler sind: Der
      Vertrag nennt sie nicht als Ausnahme, sie fallen unter „alle anderen Hinweise und Warnungen entfallen“. Ein Test pinnt
      die Setzung ausdrücklich.
- [ ] Die vier Ausnahmen bleiben unter `--quiet` sichtbar und sind je einzeln gepinnt: Sicherheitslockerungs-Hinweise bei
      `init`, `generate` und `up`; Degradations-Fallbacks
      ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer): Container ohne Daemon, `vfs`-Storage, Wegfall der Egress-Restriktion) in der Befehlsausgabe
      und in `doctor`; Compose-Fortschritt von `up` auf stderr ([`LH-NFA-PERF-002`](../../../../spec/lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker));
      die Pfad-Zusammenfassung von `--force` vor dem Überschreiben bei `init` und `generate`
      ([`LH-FA-CLI-005A`](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz)).
- [ ] Die Review-Befunde sind behoben: F-1 (der Kommentar am Doctor-Test in `cli_test.go` nennt die geltende Zusage, nicht die frühere Fassung),
      F-3 (`doctor --quiet --strict` mit nur Warn-Items endet mit Exit 11 bei leerer stdout; die Fehlermeldung auf stderr nennt für
      Skript-Nutzer das auslösende Item bzw. die Ursache, ein Test pinnt Exit-Code und stderr —
      [`LH-FA-DIAG-003`](../../../../spec/lastenheft.md#lh-fa-diag-003--fehlerklassifikation)), F-4 (der Kommentar in `root.go` verweist auf `visibleWarnings`, nicht auf das nicht
      vorhandene `[writeWarning]`).
- [ ] Doku-Update: `docs/user/benutzerhandbuch.md` (Optionstabellen, Zeilen um
      299 und 1280; `remove --quiet`) nennt `--quiet` als „nur Fehler“ mit den vier Ausnahmen;
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
| `internal/adapter/driving/cli/root.go` | update | `PersistentPreRunE`: `case a.quiet` setzt `slog.LevelError`; Flag-Hilfetext und Kommentare passen (F-4: Verweis auf `visibleWarnings`) |
| `internal/adapter/driving/cli/cli.go` | update | Feldkommentar „--quiet lowers it to Warn“ auf Fehler-Pegel |
| `internal/adapter/driving/cli/init.go`, `generate.go` | update | Zusammenfassung unter `--quiet` nicht schreiben (`printInitSummary`, Generate-Zusammenfassung); Lockerungs-Hinweise und Fallback-Ausweisung bleiben über `visibleWarnings` |
| `internal/adapter/driving/cli/add.go` | update | `printAddSummary` unter `--quiet` nicht schreiben; Flag-Struct bekommt `Quiet` |
| `internal/adapter/driving/cli/remove.go` | update | neue Option `--quiet` (Flag-Struct `Quiet`, `a.quiet` durchreichen); Erfolgsmeldung, „Removed volumes:“ und Purge-Warnungen ([`LH-FA-ADD-007`](../../../../spec/lastenheft.md#lh-fa-add-007--service-entfernen)) unter `--quiet` unterdrücken; Fehlerpfade bleiben |
| `internal/adapter/driving/cli/up.go` | update / prüfen | Tabelle und Diagnosen unterdrückt, Lockerungs-Hinweise sichtbar; Compose-Fortschritt auf stderr bleibt (Kommentar auf die Ausnahme im Lastenheft umstellen) |
| `internal/adapter/driving/cli/down.go`, `statusview.go` | prüfen | unterdrücken Erfolgsmeldung und Volume-Liste bereits; Tests pinnen es |
| `internal/adapter/driving/cli/config.go`, `template.go`, `logs.go` | prüfen | Prüfpunkt §1: Ergebnis-Ausgaben bleiben, `config set`-Zusammenfassung zuordnen und pinnen |
| `internal/adapter/driving/cli/doctor.go` | update | `writeDoctorReport`: unter `--quiet` Fehler-Items und Fallback-Items; Kopf-/Summary-Zeile entfallen (Pin); F-3: Erklärung auf stderr für `--quiet --strict` mit nur Warn-Items |
| `internal/hexagon/port/driving/removeservice.go`, `internal/hexagon/application/doctor_sandbox.go` | update | Kennzeichen „Fallback“ an den Items bzw. `WarningEntry`s der drei Degradations-Fälle (nicht im JSON); [`LH-FA-ARCH-003`](../../../../spec/lastenheft.md#lh-fa-arch-003--import-regeln-und-enforcement)-konform, kein I/O |
| `internal/adapter/driving/cli/verbosity_test.go` | update | `quiet-lowers-to-warn` → Fehler-Pegel; `debug-takes-precedence-over-quiet` bleibt unverändert (kehrt sich erst in [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../open/slice-v1-cli-verbosity-letzte-option-gewinnt.md) um) |
| `internal/adapter/driving/cli/cli_test.go` | update | F-1: Kommentar am Doctor-Test auf die geltende Zusage; F-3-Pin `doctor --quiet --strict` |
| `internal/adapter/driving/cli/sandbox_test.go`, `updown_cluster_test.go`, `statusview_test.go`, `down_acceptance_test.go`, `up_acceptance_test.go`, `doctor_json_test.go` sowie Tests zu `init`/`generate`/`add`/`remove` | update / prüfen | neue Fälle Happy/Boundary/Negative je Subkommando und je Ausnahme; `--quiet --json` unverändert — nach [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) |
| `docs/user/benutzerhandbuch.md` | update | Optionstabellen (Zeilen um 299 und 1280): „nur Fehler; vier Ausnahmen“, `remove --quiet` |
| `CHANGELOG.md` | update | sichtbare Verhaltensänderung (Erfolgsmeldungen von `init`/`generate`/`add`/`remove` entfallen unter `--quiet`, `doctor --quiet` zeigt keine übrigen Warn-Items mehr) |

## 4. Trigger

**Start** (`next` → `in-progress`): Lastenheft 0.9.0 ist angenommen (der CR
[`slice-cr-cli-quiet-weitere-ausnahmen`](../in-progress/slice-cr-cli-quiet-weitere-ausnahmen.md) liegt in `done/`) und der Projektinhaber
priorisiert den Slice.

**Reihenfolge zum Slice [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../open/slice-v1-cli-verbosity-letzte-option-gewinnt.md):**
logisch sind beide unabhängig — dieser Slice ändert die Abbildung `--quiet` →
Pegel und die Ausgabe-Filter, jener die Auswahl unter mehreren Optionen. Sie
berühren aber dieselbe Stelle (`PersistentPreRunE`-Switch in `root.go`) und
dieselbe Testfunktion in `verbosity_test.go`. Die Reihenfolge ist deshalb eine
Konflikt-Vermeidung, keine Voraussetzung.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Reichweite (alle Subkommandos) und Ausnahmen-Kennzeichnung lassen
  sich nicht in einem Review-Umfang halten — Zerlegungsschnitte siehe §6.
- `in-progress` → `open` (blockiert): Lastenheft 0.9.0 ist nicht angenommen (der CR
  [`slice-cr-cli-quiet-weitere-ausnahmen`](../in-progress/slice-cr-cli-quiet-weitere-ausnahmen.md) noch nicht in `done/`) oder
  der Projektinhaber ändert eine der vier Ausnahmen. Dieser Slice liegt bis zur Schließung des CR in `open/`.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit
Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **Umfang wächst:** Sieben Subkommandos, vier Ausnahmen und drei Review-Befunde in einem Slice. — **Ausgang:** entfallen:
  Ist der Umfang beim Start nicht tragbar, wird nach Subkommando zerlegt; Zerlegungsschnitte, in dieser Reihenfolge: (1)
  Fehler-Pegel und Review-Befunde F-1/F-3/F-4 samt `doctor`-Fallback-Items; (2) `init`/`generate`/`add` (Zusammenfassungen);
  (3) `remove` mit neuer Option `--quiet` und Purge-Warnungen; (4) `config set`/`template`-Zuordnung (Prüfpunkt §1). Die
  Ausnahmen-Pins gehören jeweils in den Slice des Subkommandos.
- **Abgrenzung der Ausnahmen:** Die Kennzeichnung kann eine Lockerung oder einen Fallback übersehen oder einen Hinweis zu
  Unrecht ausnehmen. — **Ausgang:** entfallen: der Test pinnt alle Lockerungs-Hinweise aus `sandboxRelaxationWarnings` (Podman,
  Docker, Egress), die drei Fallbacks und die ausgenommenen Hinweise einzeln.
- **`doctor --quiet` verliert übrige Warn-Items:** Skripte, die die Warn-Zeilen lasen, sehen sie nicht mehr; der Exit-Code
  trägt das Ergebnis. — **Ausgang:** entfallen: die Folge ist gewollt (Entscheidung des Projektinhabers) und steht im
  CHANGELOG und in den Optionstabellen des Benutzerhandbuchs (DoD); F-3 sichert die Erklärbarkeit bei `--strict`.
- **Reihenfolge zum Folge-Slice:** beide ändern `root.go` und `verbosity_test.go`. — **Ausgang:** entfallen: logisch
  unabhängig (§4); ein Merge-Konflikt wird beim zweiten Slice aufgelöst.

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
