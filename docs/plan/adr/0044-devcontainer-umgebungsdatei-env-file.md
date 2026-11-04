# ADR 0044: Umgebungsdatei des Devcontainers — `--env-file` in `runArgs`, Anlage der Datei durch `generate devcontainer`, benutzergeführt und nie überschrieben

**Status:** Accepted

**Datum:** 2026-11-04

**Autor:** pt9912 (Rolle Architect — Entscheidungen des Projektinhabers 2026-11-04 des Planungsvorgangs vom selben Tag: Zweck `--env-file` ohne Umweg über `${localEnv:<NAME>}`, Vorlage mit Platzhaltern, Anlage der echten Datei durch `generate devcontainer`; der Lastenheft-CR [`cr-vorschlag-devcontainer-env-datei`](../cr/cr-vorschlag-devcontainer-env-datei.md) ist freigegeben und übernommen, Lastenheft 0.35.0)

**Bezug:** [`LH-FA-DEV-012`](../../../spec/lastenheft.md#lh-fa-dev-012--umgebungsdatei-des-devcontainers), [`LH-FA-DEV-009`](../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-DEV-001`](../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), [`LH-FA-INIT-005`](../../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz); Messbasis [`docs/messungen/slice-v1-devcontainer-env-datei/README.md`](../../messungen/slice-v1-devcontainer-env-datei/README.md) (E1–E6 samt Grenzen)

**Schärft:** [`SPEC-055`](../../../spec/spezifikation.md#spec-055--umgebungsdatei-des-devcontainers-pfad-inhalt-reihenfolge-in-runargs-anlage-und-prüfung) (neu: Pfad, Inhalt, Reihenfolge in `runArgs`, Anlage, Prüfung `devcontainer.envfile`), [`SPEC-053`](../../../spec/spezifikation.md#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs) (Reihenfolge der `runArgs`) und [`SPEC-052`](../../../spec/spezifikation.md#spec-052--feature-ollama-host-im-devcontainer-eintragstyp-container-option-meldung-und-grenzen) (Zahl der Einträge von `runArgs`)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR); Form nach [`MR-008`](../../../harness/conventions/MR-008-adr-form.md) (neue ADRs MADR).

---

## Kontext

Ein Projekt will Umgebungsvariablen seines Devcontainers in einer Datei führen, die `docker run --env-file` liest. Bisher kennt der erzeugte `devcontainer.json` nur Verweise `${localEnv:<NAME>}` in `remoteEnv` (Sandbox-Profil) und Werte in `containerEnv`; Zugangsdaten stehen nach [`LH-FA-DEV-009`](../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) in keiner von `u-boot` erzeugten Datei. Der Lastenheft-CR (0.35.0) hat die benutzergeführte Datei `.devcontainer/devcontainer.env` davon ausgenommen und [`LH-FA-DEV-012`](../../../spec/lastenheft.md#lh-fa-dev-012--umgebungsdatei-des-devcontainers) eingeführt. Diese ADR trägt die Entscheidungen dahinter.

**Gemessen** (Messbasis, 2026-11-04, Docker Engine 29.9.3 `linux/amd64`, verschachtelt, ein Lauf je Fall, nur `docker run --env-file`):

- E1: eine fehlende Datei bricht den Start mit Exit-Code `125` ab (`--env-file: open …: no such file or directory`).
- E2/E3: eine Datei nur aus Kommentarzeilen startet mit Exit-Code `0`; `NAME=WERT` setzt den Wert samt Leerzeichen, Anführungszeichen bleiben Teil des Werts, `NAME=` setzt einen leeren Wert.
- E4: ein Präfix `export` bricht den Start mit Exit-Code `125` ab.
- E5: ein relativer Pfad wird gegen das Arbeitsverzeichnis des Clients aufgelöst.
- E6: ein mit `-e` gesetzter Wert gewinnt gegen denselben Namen in der Datei.

**Nicht gemessen:** der Start über `u-boot devcontainer up` und den Editor, andere Engines und Hosts, die Reihenfolge gegen `remoteEnv` und `containerEnv`, `CRLF`, mehrzeilige Werte.

**Stand im Code (gelesen am 2026-11-04, mit Fundstelle).**

- `internal/hexagon/application/devcontainer_sandbox.go`, `devcontainerTemplateData`: der Standardzweig setzt `RunArgs = resourceRunArgs + hostGatewayRunArgs` (leer ohne beide Schlüssel), der Sandbox-Zweig hängt der Reihe nach Podman-Optionen, `--privileged`, `--cap-add=NET_ADMIN`, `resources`, `ollama-host` an. Das Template `templates/devcontainer/devcontainer.json.tmpl` rendert `runArgs` nur bei nicht leerer Liste (zwei Stellen, Sandbox und Standard).
- `internal/hexagon/application/generate.go`, `planDevcontainerFiles` und `planDevcontainerFile`: jede Datei trägt einen verwalteten Block `init`; `planDevcontainerFile` liest eine vorhandene Datei und klassifiziert sie (`classifyExistingBlock`); `executeDevcontainerPlans` kennt die Aktionen `NoOp`, `Write`, `ReplaceBlock`, snapshottet vor jeder Mutation im `generateJournal` (`captureFile` liest eine vorhandene Datei) und nimmt bei einem Fehler alles zurück.
- `internal/hexagon/application/initproject.go`, `devcontainerFileTemplates` und `planFile`: `init --devcontainer` plant die zwei Dateien über einen **zweiten** Pfad mit eigener Entscheidungsregel (vorhandene Datei ohne `--force`/`--backup` bricht ab, mit `--backup` wird gesichert und überschrieben); `devcontainerTemplateData` wird dort mit `dc == nil` aufgerufen (Zeile 571).
- `internal/hexagon/application/templates/gitignore.tmpl` ist über `fileTemplates` eine ganz verwaltete Datei von `init`.
- `internal/hexagon/application/doctor_sandbox.go`, `credentialScanFiles`: die Klartext-Prüfung liest sieben feste Dateien; `.devcontainer/devcontainer.env` ist nicht darunter.
- `internal/hexagon/application/devcontainer_start.go`, `containerNameOverlay`: hängt `--name=<Name>` an `runArgs` der Konfigurationskopie des Hilfs-Containers; [`SPEC-049`](../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (Spiegelpfad) legt den Projektordner im Hilfs-Container unter denselben absoluten Pfad wie auf dem Host.

## Entscheidung

Wir wählen **die Umgebungsdatei als `--env-file` in `runArgs`, von `generate devcontainer` (und `init --devcontainer`) ohne Werte angelegt, benutzergeführt, nie überschrieben**:

1. **`--env-file` in `runArgs`, nicht `${localEnv:NAME}`-Fallback.** Die erzeugte Konfiguration trägt `--env-file=${localWorkspaceFolder}/.devcontainer/devcontainer.env` als **ersten** Eintrag von `runArgs`, in beiden Profilen und immer. `u-boot` liest die Datei nicht und baut keinen zweiten Mechanismus mit eigener Rangfolge; die Engine setzt die Variablen.
2. **Absoluter Pfad über `${localWorkspaceFolder}`.** Ein relativer Pfad löst gegen das Arbeitsverzeichnis des Clients auf (E5); der Editor und die CLI starten aus wechselnden Verzeichnissen. Die Variable steht wörtlich in der Datei; `u-boot` löst sie nicht auf.
3. **Erste Stelle der `runArgs`.** Die Position ist unabhängig von allen anderen Optionen (Podman-, Docker-, Egress-, `resources`-, `ollama-host`-Zeilen); die bestehenden Sätze „`ollama-host` ist der letzte Eintrag“ und „`resources` vor `ollama-host`“ bleiben wahr. Die Engine wertet die Reihenfolge der Optionen untereinander nicht aus, der Vorrang gegen `-e` ist gemessen (E6) und hängt nicht an ihr.
4. **`generate devcontainer` legt die Datei an, wenn sie fehlt.** Eine fehlende Datei lässt `docker run` mit Exit-Code `125` abbrechen (E1); die Konfiguration darf nie auf eine fehlende Datei zeigen. Der Inhalt sind Kommentarzeilen und Leerzeilen (E2: startet), keine Variablen. `init --devcontainer` legt sie ebenso an.
5. **Die Datei ist benutzergeführt: kein Managed Block, nie gelesen, nie überschrieben, nie gesichert.** Eine vorhandene Datei bleibt byte-gleich, auch bei `init --force --backup`. Das trägt die Ausnahme der Zugangsdaten-Zusage aus dem Lastenheft: `u-boot` gibt keinen Wert aus, weil es keinen liest.
6. **`.devcontainer/devcontainer.env.example` ist eine verwaltete Datei** (Block `init`, `#`-Kommentare wie `Dockerfile`), mit Platzhaltern `CHANGEME_…`, nie reale Werte. Sie wird getrackt.
7. **`.gitignore` nur über `u-boot init`.** `gitignore.tmpl` trägt die Zeile `.devcontainer/devcontainer.env`; `generate devcontainer` schreibt die `.gitignore` nie. Bestehende Projekte erhalten die Zeile über `init --backup` oder von Hand.
8. **`doctor` prüft nur die Existenz.** Die Prüfung `devcontainer.envfile` meldet `warn`, wenn die Konfiguration die Zeile trägt und die Datei fehlt; sie liest die Datei nicht, fragt Git nicht und prüft die Ignorierung nicht.
9. **Abgrenzung der Zusage.** Zugesagt ist die erzeugte Konfiguration und die Anlage der Dateien, soweit Tests es belegen. Der Start über `u-boot devcontainer up` ist erst zugesagt, wenn der Docker-Test ihn belegt (Lesart: der Spiegelpfad aus [`SPEC-049`](../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) legt die Datei im Hilfs-Container unter denselben Pfad; gemessen ist das nicht).

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun (Benutzer pflegt `runArgs` und die Datei selbst) | kein Vertragsdiff | die verwaltete `devcontainer.json` überschreibt Zeilen im Block bei jedem Lauf; der Benutzer müsste hinter den Block schreiben; jedes Projekt baut sich die Datei neu |
| B — `${localEnv:NAME}`-Fallback: `u-boot` liest die Datei und füllt `remoteEnv`/`containerEnv` | kein `runArgs`, kein Pfadproblem | `u-boot` liest Werte (Zugangsdaten-Zusage fiele für den Leser); zweiter Mechanismus mit eigener Rangfolge gegen `${localEnv:…}`; widerspricht der Entscheidung des Projektinhabers |
| C — `--env-file` in `runArgs`, Datei nicht angelegt (Benutzer legt an) | `u-boot` fasst die Datei nie an | eine fehlende Datei lässt jeden Start mit Exit 125 scheitern (E1), auch frisch erzeugte Projekte |
| **D — `--env-file`, Anlage ohne Werte, benutzergeführt (gewählt)** | Start gelingt nach `generate`; `u-boot` liest nie; wenige Änderungen im Generator | jedes erzeugte `devcontainer.json` ändert sich (neue Zeile `runArgs`); Datei kann gelöscht werden (dann Exit 125, `doctor` warnt) |
| E — `.gitignore` auch von `generate devcontainer` ergänzen | bestehende Projekte geschützt | `generate devcontainer` fasst eine vollständig verwaltete Datei von `init` an, die der Benutzer pflegt; eigener Block-Mechanismus nötig; im CR abgelehnt (Gegenposition benannt) |
| F — `doctor` prüft zusätzlich „getrackt oder nicht ignoriert“ | schützt vor einem Commit von Zugangsdaten | verlangt eine Git-Abfrage (Port-Tragfähigkeit ungeprüft); im CR auf „fehlt“ begrenzt |

## Konsequenzen

- Positiv: ein Mechanismus, den die Engine trägt; keine Lese-Pflicht für `u-boot`; Anlage ohne Werte hält die Zugangsdaten-Zusage für alles Erzeugte; die Position „erste Stelle“ lässt alle vorhandenen Reihenfolge-Sätze gelten.
- Negativ (akzeptiert): (a) die Ausgabe **jedes** Projekts ändert sich — das Standardprofil erzeugt bisher ohne `resources`/`ollama-host` kein `runArgs`; das Template-Datum „byte-identisch“ im Doc-Kommentar von `devcontainerTemplateData` entfällt für das UID-Default und ist zu berichtigen. (b) Bestehende Projekte ohne `init --backup` ignorieren die Datei nicht; die Gefahr, Zugangsdaten zu committen, liegt dort beim Benutzer (Entscheidung des Projektinhabers, Gegenposition im CR benannt, Maßnahme: Nutzer-Doku). (c) Der Dateimodus ist `defaultFileMode` (`0o644`, `initproject.go`); `0o600` wäre sicherer, aber ob der Benutzer im Hilfs-Container die Datei lesen darf, ist nicht gemessen. (d) Zwischen Planung und Schreiben kann eine Datei entstehen und überschrieben werden; dasselbe Fenster besteht für alle bisher geschriebenen Dateien, ein Einzelnutzer-CLI, kein eigener Schutz. (e) Löscht der Benutzer die Datei, scheitert der Start mit Exit 125, bis `generate devcontainer` sie neu anlegt; `doctor` weist darauf hin.
- Folgepflicht: [`SPEC-055`](../../../spec/spezifikation.md#spec-055--umgebungsdatei-des-devcontainers-pfad-inhalt-reihenfolge-in-runargs-anlage-und-prüfung) mit Beleg-Tabelle (angelegt); [`SPEC-053`](../../../spec/spezifikation.md#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs), [`SPEC-052`](../../../spec/spezifikation.md#spec-052--feature-ollama-host-im-devcontainer-eintragstyp-container-option-meldung-und-grenzen) und die Sätze mit „kein `runArgs`“ ([`SPEC-047`](../../../spec/spezifikation.md#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen), [`SPEC-044`](../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten)) berichtigt; die Prüfung `devcontainer.envfile` braucht die Zeile in der Code-Registry (`DefaultAllowedCodes`, Tabelle in der Spezifikation) und ist unter `--quiet` ungekennzeichnet; Doc-Kommentar von `devcontainerTemplateData` nachziehen; Nutzer-Doku, CHANGELOG, Maintainer-Doku (Slice DoD, Liefer-Punkt Doku). Kein neuer Folge-Slice.

### Architektur-Befund (hexagonale Passung)

Die Anlage der Datei gehört in den **Application-Service** (`GenerateService`, Hexagon), nicht in einen Adapter: sie braucht nur `Exists`, `MkdirAll` und `WriteFile` des vorhandenen Ports `driven.FileSystem`; kein neuer Port, keine neue Port-Methode, kein Domain-Typ (der Pfad und die Zeile sind Konstanten der Application, `${localWorkspaceFolder}` ist ein Wert der Dev-Container-Konfiguration, kein Domain-Begriff). Das `doctor`-Check gehört in `DoctorService` (Muster `doctor_sandbox.go`), liest `devcontainer.json` und prüft `Exists`; die Registry-Zeile in `coderegistry.go` ist Test-Infrastruktur des CLI-Adapters, kein Architekturbruch. Auffälligkeiten:

1. **Die Datei passt nicht in `planDevcontainerFile`.** Diese Funktion rendert einen Block, liest die vorhandene Datei und klassifiziert sie; für die Umgebungsdatei ist beides verboten (kein Block, kein Lesen). Der Weg mit dem geringsten Zusatzaufwand: eine zweite Planfunktion `planDevcontainerEnvFile` liefert denselben `devcontainerFilePlan` — Aktion `Write` mit gerendertem Inhalt, wenn `Exists` falsch ist, sonst `NoOp`. `executeDevcontainerPlans`, `writeDevcontainerNewFile`, das Journal und die Aggregation bleiben unverändert; die Atomarität (Phase 1 ohne Schreiben, Phase 2 mit Rollback) gilt automatisch. Dass eine vorhandene Datei nie gelesen wird, folgt daraus, dass `captureFile` (liest die Datei) nur bei Aktion ungleich `NoOp` läuft; ein Test pinnt es mit einem lesezählenden Fake.
2. **`--dry-run`** braucht keinen Code: der Dry-Run-Lauf ersetzt das Dateisystem (`generate.go`, `PreviewMode != PreviewDryRun`), `Changed` führt die Dateien auf, die angelegt würden.
3. **Zweiter Erzeugungspfad.** `init --devcontainer` plant mit `planFile`; dessen Regel „vorhandene Datei ohne `--force` bricht ab, mit `--backup` wird gesichert und überschrieben“ ist für die Umgebungsdatei falsch. Der Init-Pfad muss die Datei vor oder neben `planFile` behandeln (Anlage nur wenn fehlend; kein Abbruch, kein `.bak`). Die Zeile in `runArgs` entsteht dort ohne weiteres Zutun, weil `devcontainerTemplateData` gemeinsam ist. Die Rendering-Hilfe für Inhalt und Pfad liegt zweckmäßig in einem gemeinsamen Helfer der Application, damit die zwei Pfade nicht auseinanderlaufen (Anti-Drift-Test wie bei `forwardPorts`).
4. **`devcontainerAggregateAction`.** Eine neu angelegte Datei zählt als `hasWrite`; ist sie die einzige Änderung (Datei gelöscht, Rest aktuell), ist die Aktion `created`. Das ist gewollt und gepinnt.
5. **`doctor`-Anwendbarkeit.** Projekte, die vor dieser Änderung erzeugt wurden, tragen die Zeile nicht; die Prüfung wird dort übersprungen (`ok` mit Meldung), sonst warnte `doctor` für jedes Altprojekt. Die Klartext-Prüfung `devcontainer.sandbox.credentials` bleibt unverändert: sie liest die Datei nicht, weil sie nicht in `credentialScanFiles` steht (gepinnt, damit die Aufnahme dort nicht still den Inhalt liest).
6. **`containerNameOverlay`** (Hilfs-Container) hängt `--name=` an `runArgs` und lässt `--env-file` stehen; ob `${localWorkspaceFolder}` im Hilfs-Container auflöst, ist nicht gemessen — dafür der Docker-Test des Slice.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Tests (`internal/hexagon/application`) | Reihenfolge der `runArgs`, vorhandene Datei nie gelesen/geschrieben, Beispieldatei verwaltet, `doctor`-Befund; Beleg-Tabelle in [`SPEC-055`](../../../spec/spezifikation.md#spec-055--umgebungsdatei-des-devcontainers-pfad-inhalt-reihenfolge-in-runargs-anlage-und-prüfung) | `make test` |
| Registry-Test (Code gegen Spezifikation) | `devcontainer.envfile` in `DefaultAllowedCodes` und Tabelle | `make test` |
| Docker-Test (Build-Tag `docker`) | Variable aus der Datei im gestarteten Container, auch über `u-boot devcontainer up` | `make test-docker` |
| `docs-check` | Kennungen und Anker | `make docs-check` |

## Re-Evaluierungs-Trigger

- Der Docker-Test zeigt, dass `${localWorkspaceFolder}` im Hilfs-Container von `u-boot devcontainer up` nicht zum Host-Pfad auflöst: dann ist Entscheidung 2 für diesen Pfad zu überprüfen (Folge-ADR oder Abgrenzung in der Zusage).
- Ein Meldungsbedarf „getrackt oder nicht ignoriert“ (Alternative F) oder ein Schreiben der `.gitignore` durch `generate devcontainer` (Alternative E): Folge-ADR mit `Supersedes`, nicht Änderung dieser.
- Eine Engine, die `--env-file` anders liest als die gemessene (Quoting, `export`, CRLF): Messung wiederholen.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-11-04 | Accepted | Planungsvorgang zur Umgebungsdatei des Devcontainers |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-0044` (Baseline-Regelwerk `modul-04-adrs.md`).
