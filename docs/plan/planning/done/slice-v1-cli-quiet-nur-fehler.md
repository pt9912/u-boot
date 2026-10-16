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
[`slice-cr-cli-quiet-weitere-ausnahmen`](../done/slice-cr-cli-quiet-weitere-ausnahmen.md) (die vier Ausnahmen stehen erst nach
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
[`slice-cr-cli-quiet-weitere-ausnahmen`](../done/slice-cr-cli-quiet-weitere-ausnahmen.md); Befunde des Reviews
[`2026-10-16-slice-v1-cli-quiet-nur-fehler`](../../../reviews/2026-10-16-slice-v1-cli-quiet-nur-fehler.md) (erster Review, F-1 bis F-7) und
[`2026-10-16-slice-v1-cli-quiet-nur-fehler-2`](../../../reviews/2026-10-16-slice-v1-cli-quiet-nur-fehler-2.md) (zweiter Review, F-1 bis F-6).

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
  (host-seitig hat nur `internal/hexagon/application/doctor_sandbox.go` ein Item, und zwar für `/dev/fuse` → `vfs`; der Wegfall der
  Egress-Restriktion und der Container ohne Daemon entstehen erst im Startscript im Container und haben kein `doctor`-Item, die
  Ausweisung dieser beiden ist daher nur dort möglich, soweit das Startscript sie meldet); `doctor`
  zeigt unter `--quiet` nur Fehler-Items, die Fallback-Items der Stufe Warn entfielen also. Nicht jedes Warn-Item von `doctor` ist ein
  Fallback (z. B. die Credentials-Prüfung).
- Der Compose-Fortschritt von `up` läuft auf stderr und ist von `--quiet` unberührt.

**Ansatz:** Der Pegel kommt aus `root.go`; die Ausgabe-Filter bleiben an `a.quiet` gekoppelt. Zusammenfassungen und
Erfolgsmeldungen der Subkommandos werden unter `--quiet` nicht geschrieben (jede Subkommando-Flag-Struct bekommt `Quiet`, wie
in `config.go`). Ein Hinweis, der eine Ausnahme trägt, wird am Eintrag gekennzeichnet (`SecurityRelaxation`, `DegradationFallback`);
Fallback-Items von `doctor` werden analog gekennzeichnet, damit `doctor --quiet` sie neben den Fehler-Items zeigt. Der
Compose-Strom von `down` läuft unter `--quiet` ohne `--json` in einen Puffer: Der `ProgressSink` ist dann ein Puffer, der nur
bei Fehlschlag auf stderr geschrieben und im Erfolgsfall verworfen wird; so bleiben die Fehlerdetails des Compose-Aufrufs
erhalten, die der Adapter nicht in die Fehlermeldung übernimmt. Der JSON-Modus bleibt unverändert (`--quiet --json` ist ein No-op).

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
- **Die Vorschau von `--dry-run` und `--diff` unter `--quiet`** — das Lastenheft nennt sie weder als Ausnahme noch als Abfrage; der
  Vertrag wird im CR [`slice-cr-cli-quiet-vorschau`](../done/slice-cr-cli-quiet-vorschau.md) entschieden, der Code folgt in
  [`slice-v1-cli-quiet-vorschau`](../open/slice-v1-cli-quiet-vorschau.md). Bis dahin ist `--dry-run --quiet` leer bei Exit 0.
- **Die host-seitige Ausweisung der Fallbacks Daemon und Egress** — die Lücke besteht unabhängig von `--quiet` und ist ein Vertrags-
  oder Code-Vorgang von [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer); sie liegt in
  [`slice-v1-sandbox-fallback-ausweisung`](../open/slice-v1-sandbox-fallback-ausweisung.md).
- **Andere Hinweis-Klassen als die vier Ausnahmen** (Recreate-Warnung
  [`LH-FA-UP-003`](../../../../spec/lastenheft.md#lh-fa-up-003--startstatus-anzeigen), Drittanbieter-Features
  [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), „no git remote“, Purge-Warnungen von `remove`)
  — sie entfallen unter `--quiet`, soweit sie keine Fehler sind.

## 2. Definition of Done

- [x] [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) erfüllt (Lastenheft 0.9.0): `--quiet` setzt den Fehler-Pegel; `init`, `generate`, `add`, `remove`,
      `config set`, `up` und `down` unterdrücken Erfolgs- und Zusammenfassungsmeldungen, `doctor --quiet` zeigt nur
      Fehler-Einträge und Fallback-Items (Ausnahme 2); `remove` hat die Option `--quiet`; `config get`, `config list`,
      `logs` und `template list` behalten ihre Ausgabe (Prüfpunkt §1, gepinnt). Tests Happy/Boundary/Negative je Subkommando
      referenzieren die IDs.
- [x] Die Purge-Warnungen und die Abhängigkeits-Warnung von `remove`
      ([`LH-FA-ADD-007`](../../../../spec/lastenheft.md#lh-fa-add-007--service-entfernen)) entfallen unter `--quiet`, soweit sie keine Fehler sind: Der
      Vertrag nennt sie nicht als Ausnahme, sie fallen unter „alle anderen Hinweise und Warnungen entfallen“. Ein Test pinnt
      die Setzung ausdrücklich.
- [x] Die vier Ausnahmen bleiben unter `--quiet` sichtbar und sind je einzeln gepinnt: Sicherheitslockerungs-Hinweise bei
      `init`, `generate` und `up`; Degradations-Fallbacks
      ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer): Container ohne Daemon, `vfs`-Storage, Wegfall der Egress-Restriktion) in der Befehlsausgabe
      und in `doctor`, soweit der Host sie ermittelt (Ist-Stand: nur der `vfs`-Fall hat ein `doctor`-Item; die beiden anderen
      entstehen im Startscript im Container, die Lücke ist hier benannt und kein neuer Nutzer-Vertrag); Compose-Fortschritt von `up` auf stderr ([`LH-NFA-PERF-002`](../../../../spec/lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker));
      die Pfad-Zusammenfassung von `--force` vor dem Überschreiben bei `init` (`generate` hat keinen `--force`-Pfad und keinen
      `AffectedFiles`-Produzenten)
      ([`LH-FA-CLI-005A`](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz)).
- [x] Die Review-Befunde sind behoben: F-1 (der Kommentar am Doctor-Test in `cli_test.go` nennt die geltende Zusage, nicht die frühere Fassung),
      F-3 (`doctor --quiet --strict` mit nur Warn-Items endet mit Exit 11 bei leerer stdout; die Fehlermeldung auf stderr nennt für
      Skript-Nutzer das auslösende Item bzw. die Ursache, ein Test pinnt Exit-Code und stderr —
      [`LH-FA-DIAG-003`](../../../../spec/lastenheft.md#lh-fa-diag-003--fehlerklassifikation)), F-4 (der Kommentar in `root.go` verweist auf `visibleWarnings`, nicht auf das nicht
      vorhandene `[writeWarning]`).
- [x] Doku-Update: `docs/user/benutzerhandbuch.md` (Optionstabellen, Zeilen um
      299 und 1280; `remove --quiet`) nennt `--quiet` als „nur Fehler“ mit den vier Ausnahmen;
      CHANGELOG-Eintrag.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo mit Wellen-Betrieb von der Welle-Closure geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapter/driving/cli/root.go` | update | `PersistentPreRunE`: `case a.quiet` setzt `slog.LevelError`; Flag-Hilfetext und Kommentare passen (F-4: Verweis auf `visibleWarnings`) |
| `internal/adapter/driving/cli/cli.go` | update | Feldkommentar „--quiet lowers it to Warn“ auf Fehler-Pegel |
| `internal/adapter/driving/cli/init.go`, `generate.go` | update | Zusammenfassung unter `--quiet` nicht schreiben (`printInitSummary`, Generate-Zusammenfassung); Lockerungs-Hinweise und Fallback-Ausweisung bleiben über `visibleWarnings` |
| `internal/adapter/driving/cli/add.go` | update | `printAddSummary` unter `--quiet` nicht schreiben; Flag-Struct bekommt `Quiet` |
| `internal/adapter/driving/cli/remove.go` | update | neue Option `--quiet` (Flag-Struct `Quiet`, `a.quiet` durchreichen); Erfolgsmeldung, „Removed volumes:“ und Purge-Warnungen ([`LH-FA-ADD-007`](../../../../spec/lastenheft.md#lh-fa-add-007--service-entfernen)) unter `--quiet` unterdrücken; Fehlerpfade bleiben |
| `internal/adapter/driving/cli/up.go` | update / prüfen | Tabelle und Diagnosen unterdrückt, Lockerungs-Hinweise sichtbar; Compose-Fortschritt auf stderr bleibt (Kommentar auf die Ausnahme im Lastenheft umstellen) |
| `internal/adapter/driving/cli/down.go` | update | Unter `--quiet` ohne `--json` ist der `ProgressSink` ein Puffer: Compose-Strom nur bei Fehlschlag auf stderr, im Erfolgsfall verworfen; Erfolgsmeldung und Volume-Zeile entfallen |
| `internal/adapter/driving/cli/statusview.go` | prüfen | unterdrückt die Tabelle bereits; Tests pinnen es |
| `internal/adapter/driving/cli/config.go`, `template.go`, `logs.go` | prüfen | Prüfpunkt §1: Ergebnis-Ausgaben bleiben, `config set`-Zusammenfassung zuordnen und pinnen |
| `internal/adapter/driving/cli/doctor.go` | update | `writeDoctorReport`: unter `--quiet` Fehler-Items und Fallback-Items; Kopf-/Summary-Zeile entfallen (Pin); F-3: Erklärung auf stderr für `--quiet --strict` mit nur Warn-Items |
| `internal/hexagon/port/driving/removeservice.go`, `internal/hexagon/application/doctor_sandbox.go`, `internal/hexagon/application/devcontainer_sandbox.go`, `internal/hexagon/domain/diagnostic.go` | update | Kennzeichen „Fallback“ (`DegradationFallback`) an `WarningEntry` und `domain.Diagnostic`; produktiv setzt es nur `doctor_sandbox.go` (`vfs`); `devcontainer_sandbox.go` erzeugt keine Fallbacks, sondern kennzeichnet die Lockerungen; nicht im JSON; [`LH-FA-ARCH-003`](../../../../spec/lastenheft.md#lh-fa-arch-003--import-regeln-und-enforcement)-konform, kein I/O |
| `internal/adapter/driving/cli/verbosity_test.go` | update | `quiet-lowers-to-warn` → Fehler-Pegel; `debug-takes-precedence-over-quiet` bleibt unverändert (kehrt sich erst in [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../open/slice-v1-cli-verbosity-letzte-option-gewinnt.md) um) |
| `internal/adapter/driving/cli/cli_test.go` | update | F-1: Kommentar am Doctor-Test auf die geltende Zusage; F-3-Pin `doctor --quiet --strict` |
| `internal/adapter/driving/cli/sandbox_test.go`, `updown_cluster_test.go`, `statusview_test.go`, `down_acceptance_test.go`, `up_acceptance_test.go`, `doctor_json_test.go` sowie Tests zu `init`/`generate`/`add`/`remove` | update / prüfen | neue Fälle Happy/Boundary/Negative je Subkommando und je Ausnahme; `--quiet --json` unverändert — nach [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) |
| `docs/user/benutzerhandbuch.md` | update | Optionstabellen (Zeilen um 299 und 1280): „nur Fehler; vier Ausnahmen“, `remove --quiet` |
| `CHANGELOG.md` | update | sichtbare Verhaltensänderung (Erfolgsmeldungen von `init`/`generate`/`add`/`remove` entfallen unter `--quiet`, `doctor --quiet` zeigt keine übrigen Warn-Items mehr) |

## 4. Trigger

**Start** (`next` → `in-progress`): Lastenheft 0.9.0 ist angenommen (der CR
[`slice-cr-cli-quiet-weitere-ausnahmen`](../done/slice-cr-cli-quiet-weitere-ausnahmen.md) liegt in `done/`) und der Projektinhaber
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
  [`slice-cr-cli-quiet-weitere-ausnahmen`](../done/slice-cr-cli-quiet-weitere-ausnahmen.md) noch nicht in `done/`) oder
  der Projektinhaber ändert eine der vier Ausnahmen. Dieser Slice liegt bis zur Schließung des CR in `open/`.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit
Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **Umfang wächst:** Sieben Subkommandos, vier Ausnahmen und drei Review-Befunde in einem Slice. — **Ausgang:** entfallen:
  der Slice wurde ungeteilt geschlossen; zwei Reviews und eine Verification trugen den Umfang, die Zerlegung war nicht nötig.
- **Abgrenzung der Ausnahmen:** Die Kennzeichnung kann eine Lockerung oder einen Fallback übersehen oder einen Hinweis zu
  Unrecht ausnehmen. — **Ausgang:** entfallen: der Test pinnt alle Lockerungs-Hinweise aus `sandboxRelaxationWarnings` (Podman,
  Docker, Egress), die drei Fallbacks und die ausgenommenen Hinweise einzeln.
- **`doctor --quiet` verliert übrige Warn-Items:** Skripte, die die Warn-Zeilen lasen, sehen sie nicht mehr; der Exit-Code
  trägt das Ergebnis. — **Ausgang:** entfallen: die Folge ist gewollt (Entscheidung des Projektinhabers) und steht im
  CHANGELOG und in den Optionstabellen des Benutzerhandbuchs (DoD); F-3 sichert die Erklärbarkeit bei `--strict`.
- **Vorschau unter `--quiet`:** `--dry-run --quiet` gibt nichts aus und endet mit Exit 0; das Lastenheft nennt die Vorschau weder
  als Ausnahme noch als Abfrage (zweiter Review, F-2). — **Ausgang:** eingetreten:
  [`slice-cr-cli-quiet-vorschau`](../done/slice-cr-cli-quiet-vorschau.md) und
  [`slice-v1-cli-quiet-vorschau`](../open/slice-v1-cli-quiet-vorschau.md).
- **Degradations-Fallbacks Daemon und Egress ohne host-seitige Ausweisung:** produktiv weist nur der `vfs`-Fall im `doctor` aus; für
  die beiden anderen gibt es keine host-seitige Quelle, und für `init`/`generate`/`up` erzeugt keine produktive Quelle einen
  Fallback-Eintrag. Die Lücke besteht unabhängig von `--quiet`
  ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)). — **Ausgang:** eingetreten:
  [`slice-v1-sandbox-fallback-ausweisung`](../open/slice-v1-sandbox-fallback-ausweisung.md).
- **`down`-Puffer ist eine Setzung ohne ausdrückliche Lastenheft-Zeile:** Ausnahme 3 nennt nur `up`; `down` steht unter den ändernden
  Subkommandos mit „nur Fehler“. — **Ausgang:** entfallen: Entscheidung des Projektinhabers vom 2026-10-16; der Test
  `TestExecute_Down_Quiet_ComposeStream` pinnt Erfolgs- und Fehlerfall, der Fehlerfall ist durch „nur Fehler“ gedeckt.
- **Reihenfolge zum Folge-Slice:** beide ändern `root.go` und `verbosity_test.go`. — **Ausgang:** entfallen: logisch
  unabhängig (§4); ein Merge-Konflikt wird beim zweiten Slice aufgelöst.

## 7. Closure-Notiz

**Belege:** Review-Reports
[`2026-10-16-slice-v1-cli-quiet-nur-fehler`](../../../reviews/2026-10-16-slice-v1-cli-quiet-nur-fehler.md) (Erstlauf: 1 HIGH, 2 MEDIUM, 1 LOW, 3 INFO;
F-1, F-3 und F-4 behoben, F-2 durch die Reichweite „alle Subkommandos“ behoben, F-5 und F-6 durch Lastenheft 0.9.0 aufgelöst, F-7 bleibt als
Grenze) und
[`2026-10-16-slice-v1-cli-quiet-nur-fehler-2`](../../../reviews/2026-10-16-slice-v1-cli-quiet-nur-fehler-2.md) (zweiter Lauf: 0 HIGH, 3 MEDIUM, 3 INFO;
F-1 `down`-Strom und F-3 CHANGELOG behoben, F-2 Vorschau in die Folge-Slices ausgelagert, F-4 in den Fallback-Slice, F-5 und F-6 als Grenze
bzw. gepinnt). Verification Evidence
[`2026-10-16-verification-slice-v1-cli-quiet-nur-fehler`](../../../reviews/2026-10-16-verification-slice-v1-cli-quiet-nur-fehler.md): bestanden, mit
sechs Auflagen für die Closure (hier eingelöst). `make gates` grün, Coverage 91,00 % (Schwelle 90 %); `make doc-trace`: 83 Anforderungen, 0 Waisen.

- **Was hat funktioniert:** Die Pins prüfen Verhalten, keine Tautologien: Die Mutations-Gegenproben M1 bis M12 fielen rot, ebenso die drei des Verifiers (`down`-Strom im Fehlerfall nicht ausgeben, Fallback-Filter im `doctor`
  entfernt, Puffer nie aktiv). Die Fehlerpins (`…_ErrorStays` je ändernder Subkommando, `…_PurgeRefusedStaysAnError`) halten den Gegenpol
  „nicht zu viel unterdrücken“; der Wire-Ausschluss des Fallback-Flags ist für `WarningEntry` und das Doctor-Item gepinnt
  (`TestDoctorJSON_FallbackFlagNotInWireForm`).
- **Was ging anders als geplant:** Die Entscheidung zu den Fallbacks aus dem Erstlauf wurde durch Lastenheft 0.9.0 überholt: Aus „Spannung,
  Architect entscheidet“ wurde Ausnahme 2 im Vertrag, der Code folgte. Der zweite Review fand drei Dinge, die der Erstlauf nicht sah: den
  `down`-Strom (Compose-Zeilen unter `--quiet`, Fehlerdetails nur dort), die Vorschau (`--dry-run --quiet` leer bei Exit 0) und den
  CHANGELOG-Widerspruch zur Regel „letzte Option gewinnt“. Die Zählung der Ausnahmen („eine“, „drei“, „vier“) war an mehreren Stellen
  nachzuziehen. Der Slice-Plan nannte `doctor`-Items für Daemon und Egress, die es nicht gibt; der zweite Review fand es.
- **Steering-Loop-Eintrag:** Ein Slice-Text beschreibt, was da ist: Der Plan nannte Quellen im Code (`doctor`-Items in `doctor_sandbox.go`), die
  es nicht gab, und der Review fand es erst nach dem Bau. Plan-Aussagen über Code werden vor dem Schreiben im Code geprüft (Suche nach dem
  Symbol oder der Item-ID). Zusammen mit der Finding-Klasse „Plan nennt Quellen, die es nicht gibt“ der beiden Reviews gezählt, nicht verkörpert.
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung angelegt; das Register ist leer.
- **Folge-Slices:** [`slice-cr-cli-quiet-vorschau`](../done/slice-cr-cli-quiet-vorschau.md) (CR: Vorschau unter `--quiet`, Lastenheft 0.10.0),
  [`slice-v1-cli-quiet-vorschau`](../open/slice-v1-cli-quiet-vorschau.md) (Code dazu),
  [`slice-v1-sandbox-fallback-ausweisung`](../open/slice-v1-sandbox-fallback-ausweisung.md) (Lücke der Fallback-Ausweisung, Entscheidung Code oder CR beim
  Projektinhaber) und
  [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../open/slice-v1-cli-verbosity-letzte-option-gewinnt.md) (Rangfolge bei Mehrfachangabe, Spec-Eintrag).
- **Risiken aus §6:** alle tragen ihren Ausgang: eingetreten (Vorschau, Fallback-Ausweisung) oder entfallen (Umfang, Abgrenzung der Ausnahmen,
  `doctor`-Warn-Items, `down`-Puffer, Reihenfolge zum Folge-Slice).
- **Release-Hinweis:** kein Release-Cut vor `slice-v1-cli-verbosity-letzte-option-gewinnt` und den beiden Vorschau-Slices; der
  CHANGELOG sagt im Präsens „follows in a later release“.
- **Grenzen:** (1) Der `--force`-Pin (`TestExecute_Init_Quiet_DoesNotSilenceForcePathSummary`) prüft nur `SilenceProgress` im Request, nicht die
  tatsächliche Ausgabe; `generate` hat keinen `--force`-Pfad und keinen `AffectedFiles`-Produzenten. (2) Die Pins der Ausnahme 2 für `init`,
  `generate` und `up` laufen über Fake-Use-Cases, weil keine produktive Quelle einen Fallback-Eintrag erzeugt; ebenso der
  Lockerungs-Zweig in `up` (Erstreview F-7). Produktiv ist nur `vfs` im `doctor` erreichbar. (3) `config show` steht in keiner der
  Lastenheft-Listen und verhält sich als Abfrage.

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
