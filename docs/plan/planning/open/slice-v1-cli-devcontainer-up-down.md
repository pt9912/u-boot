# Slice v1: Befehl `u-boot devcontainer up` und `down` (Application, CLI, Ausgabe, Handbuch, Docker-Tests S1, S3, S4, S8, S9, Einzelmessung Colima)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-devcontainer-start-ohne-editor`](../welle-devcontainer-start-ohne-editor.md). Dritter von vier Folge-Slices von [`slice-v1-devcontainer-start-ohne-cli`](../in-progress/slice-v1-devcontainer-start-ohne-cli.md); nach [`slice-v1-devcontainer-start-fundament`](../open/slice-v1-devcontainer-start-fundament.md) und vor [`slice-v2-spec-devcontainer-start`](../open/slice-v2-spec-devcontainer-start.md). Release-Zuordnung offen (Vorschlag des Planners: V2, nächstes Minor nach diesem Slice; entscheidet der Projektinhaber).

**Bezug:** die neue Anforderung des Lastenhefts aus [`slice-cr-devcontainer-start-ohne-editor`](../open/slice-cr-devcontainer-start-ohne-editor.md) (Devcontainer ohne Editor starten; Kennung steht nach dem CR im Lastenheft), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (`--quiet`, Ausweisung), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (10, 11, 12), [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) und [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) (`--json`, Enum `command`), [`LH-FA-DEV-004`](../../../../spec/lastenheft.md#lh-fa-dev-004--benutzerrechte) (UID), [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Ausweisung der Fallbacks), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zugangsdaten), [`LH-NFA-PORT-004`](../../../../spec/lastenheft.md#lh-nfa-port-004--ausgeliefertes-container-image); [`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) (muss `Accepted` sein; Entscheidungen 3, 6, 8, 9).

**Berührte Spec-Stellen:** [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) (Ausnahme 1 nennt die Ausweisungen des Sandbox-Profils; die Socket-Lockerung des Befehls kommt hinzu), [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) (Ausweisung der Fallbacks nach dem Start), [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) (Enum `command`), `architecture.md` §2.3 und §7. Der Slice **liest** sie und ändert keine; die Nachführung ist der Spezifikations-Slice. Der Verweis zeigt aufwärts.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Architect). **Datum:** 2026-10-20.

---

## 1. Ziel und Abgrenzung

**Ziel:** `u-boot devcontainer up [--recreate]` und `u-boot devcontainer down` laufen auf dem Fundament des vorigen Slice: Application-Service und Driving-Port, Cobra-Befehle mit Ausgabe, `--quiet`, `--json`, Exit-Codes 2/10/11/12/14 und der Ausweisung der Socket-Lockerung und der Fallbacks nach [`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidungen 3, 4 (Ausweisung) und 6. Dazu die Docker-Tests mit Wächter für S1, S3, S4, S8, S9, die **Einzelmessung** des ausgelieferten Befehls auf Colima durch den Projektinhaber, und Handbuch, `CHANGELOG.md` und JSON-Kontrakt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Port, Adapter, eingebetteter Kontext, Pins** — [`slice-v1-devcontainer-start-fundament`](../open/slice-v1-devcontainer-start-fundament.md); dieser Slice ruft sie nur auf.
- **`spec/`** (Lastenheft, Spezifikation, Sicht) — der Lastenheft-CR liegt davor, die Spezifikation danach ([`slice-v2-spec-devcontainer-start`](../open/slice-v2-spec-devcontainer-start.md)); Bedingungen im Handbuch brauchen einen Test aus diesem Slice oder entfallen (`AGENTS.md` §3.10).
- **`u-boot up` startet den Devcontainer mit** — `up` betrifft den Compose-Stack; ein gemeinsamer Start wäre eine eigene Entscheidung.
- **Weitere CLI-Optionen** (`--build-no-cache`, Verzeichnis-Option, `--remove` bei `down`) — nicht Teil der Entscheidung ([`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 3); sie lassen sich ohne Folge-ADR ergänzen, wenn ein Nutzer sie braucht.
- **Docker Desktop, Podman, entfernte Engine, Windows, rootless Docker, SELinux** — nicht zugesagt; keine Probe dafür.
- **Das Hilfs-Image im Runtime-Image** — der Befehl läuft dort nicht (Exit 11 über die Vorab-Prüfung); das Image bleibt unverändert.

## 2. Definition of Done

Drei Liefer-Punkte; Gate-Läufe und Closure-Pflichten darunter zählen nicht mit.

- [ ] **Befehl:** `devcontainer up` (`--recreate`) und `devcontainer down` mit Application-Service, Driving-Port, Cobra-Befehlen, Ausgabe und Exit-Codes nach [`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 3 und 6 (bare `u-boot devcontainer` und `--dry-run`/`--diff`: Exit 2; kein `.devcontainer/devcontainer.json` oder ein Pfad mit unzulässigem Zeichen: 10; 11 über `ErrDockerUnavailable`/`ErrContainerEngineNotLocal`; 12 über `ErrDevcontainerRuntime`, in **beiden** Klassifikatoren eingetragen, je mit Pin-Test); Namen der `${localEnv:<NAME>}`-Referenzen aus dem Text von `devcontainer.json` (reine Funktion); Ausweisung der Socket-Lockerung vor dem Start des Hilfs-Containers, sichtbar unter `--quiet`, im JSON als `diagnostics[]`-Eintrag; Fallback-Ausweisung nach `up` über die vorhandene Mechanik; `--json` mit genau einem Umschlag (`command: devcontainer`, `subcommand: up|down`, `data`).
- [ ] **Docker-Tests mit Wächter und Einzelmessung:** S1 (Start eines erzeugten Projekts, Standard- und Sandbox-Profil), S3 (UID und `updateRemoteUserUID`), S4 (`down` stoppt und entfernt nicht, Stop und `up` gleicher Container, `--recreate` neuer Container, benanntes Volume mit Marker überlebt), S8 (Fehlerbild und Abbildung 10/11/12), S9 (Labels des Wegs, Grenze zum Editor); die Wächter stehen in §3. Dazu die **Einzelmessung des ausgelieferten Befehls** auf macOS mit Colima durch den Projektinhaber als Messprotokoll unter `docs/messungen/slice-v1-cli-devcontainer-up-down/` (Skript unter `tools/messung/`, nur `make` und Docker, Rohdaten ohne Zugangsdaten).
- [ ] **Doku:** Handbuch (Abschnitt 4.9 *Voraussetzung*, ein Unterabschnitt zum Befehl mit Erststart, Socket-Ausweisung, Grenzen und Einzelmessung, Abschnitt *Container-Image*), `CHANGELOG.md` (`[Unreleased]`), `docs/maintainer/cli-json-contract.md` (neuer Abschnitt, Matrix-Zeile); alles ohne Anforderungs-, ADR- oder Spec-Kennungen im Handbuch (`AGENTS.md` §5 Regel 6).
- [ ] `make gates` grün, `make test-docker` grün (Exit-Code von `make` selbst).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: Roadmap-Zeile „In Arbeit“ bei der Closure; ADR-Index unverändert.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; widerlegt ein Test eine Zeile der Zusage, steht der Befund in §7 und der Slice geht zurück an den Architect.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; der Implementer ergänzt sie in seinem Lauf und gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab (`AGENTS.md` §3.10). Eine Zeile „prüfen“ nennt eine Datei, die der Architect nicht gelesen hat.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/port/driving/devcontainer.go` | neu | Driving-Port `DevcontainerUseCase` (`Up`, `Down`) mit Request/Response und den Validierungs-Sentinels (kein `devcontainer.json`, unzulässiger Pfad). Vorbild `down.go` (gelesen: `DownRequest.BaseDir`, Sentinel-Konvention) |
| `internal/hexagon/application/devcontainer_start.go` | neu | Service: liest `.devcontainer/devcontainer.json` über `driven.FileSystem`, prüft den Pfad, extrahiert die `${localEnv:<NAME>}`-Namen (reine Funktion, `application-no-net` gelesen: kein `net`), ruft `DevcontainerLauncher`, weist danach die Fallbacks über `fallbackWarnings` aus (`devcontainer_fallback.go`, gelesen: `configuredFallbackProbes(dc, sandbox)`, `fallbackWarnings(ctx, r, probes)`; die Konfiguration kommt aus `u-boot.yaml`, wie `generate` sie lädt — Lader nicht gelesen, prüfen; ohne `u-boot.yaml` entfällt die Ausweisung, kein Fehler) |
| `internal/hexagon/application/noop.go` | **unverändert** | gelesen: nil-tolerante Defaults gibt es für Fortschritt, Bestätigung, Logger und Inspector; der Launcher ist wie `DockerEngine` in `NewDownService` ein Pflicht-Port im Konstruktor (`cmd/uboot/main.go`, gelesen) |
| `internal/adapter/driving/cli/devcontainer.go` | neu | Cobra-Eltern `devcontainer` (ohne Unterbefehl: Exit 2, auch mit `--json`) mit `up` (`--recreate`) und `down`; `--dry-run`/`--diff` mit Exit 2 abgewiesen; Umschlag über `newDataEnvelope`. Vorbild `egress.go` (gelesen: `ErrEgressSubcommandRequired`, `ErrEgressPreviewNotApplicable`, `isEgressUsageError`) und `down.go` (gelesen: `downFlags`, `reportError`-Pfad, Fortschritt zurückgehalten unter `--quiet`) |
| `internal/adapter/driving/cli/root.go` | update | eine `AddCommand`-Zeile in `buildRootCommand` (gelesen: Zeilen 128 bis 138); `rootLongText`/Hilfe nennt den Befehl (prüfen) |
| `internal/adapter/driving/cli/cli.go` | update | `New`-Option für den Use-Case (Muster `WithEgressListUseCase`), `ExitCode`: `driven.ErrDevcontainerRuntime` → 12 (gelesen: Zeilen 283 ff., Sentinels vor den Application-Sentinels); die neuen Usage-Sentinels in `isUsageError` (Kommentar zu `ErrEgressSubcommandRequired` gelesen) |
| `internal/adapter/driving/cli/devcontainersentinel.go` oder Erweiterung von `composesentinel.go` | neu / update | Mapper für `diagnostics[]` des Befehls; zwei Klassifikatoren (Mapper **und** `ExitCode`), Pin-Test je Sentinel. `composesentinel.go` gelesen: Compose-Cluster für `up` und `down`; `ErrComposeRuntime` bleibt Compose vorbehalten |
| `internal/adapter/driving/cli/devcontainer_test.go` und weitere | neu | Akzeptanz ohne Docker mit gefälschtem Launcher: Exit-Codes je Sentinel, `--quiet` (Ausweisung bleibt, Erfolgsmeldung und Fortschritt entfallen), `--quiet --json` gleich `--json`, bare `devcontainer`, `--dry-run`, Umschlag von `up` und `down`; Muster `down_acceptance_test.go`, `quiet_modifying_test.go` (Dateien vorhanden, nicht gelesen) |
| `cmd/uboot/main.go` | update | Wiring: `docker.NewLauncher()` (Name im Fundament-Slice), Service, `cli.WithDevcontainerUseCase(…)` (gelesen: `run` konstruiert Adapter und Services, `cli.New(…)` mit Optionen) |
| `internal/e2e/devcontainer_start_docker_test.go` | neu | Docker-Tests S1, S3, S4, S8, S9 mit den Wächtern unten; Hilfen aus `acceptance_helpers.go` (nicht gelesen) und `sandbox_status_docker_test.go` (gelesen: `requireDocker`, `buildProjectImage`, `startLabeled`; dort bildet der Test die CLI nach — hier läuft die echte) |
| `docs/user/benutzerhandbuch.md` | update | Abschnitt 4.9 *Voraussetzung* (gelesen: Zeilen 692 bis 696 nennen Editor oder `devcontainer`-CLI) um den Befehl; neuer Unterabschnitt (Erststart braucht Netz und baut, Socket-Ausweisung, `down` stoppt nur, `--recreate`, Grenzen, Einzelmessung Colima, Zahlen aus dem Messprotokoll des Fundament-Slice); Abschnitt *Container-Image* (gelesen, Zeile 146 ff.: nennt, was dort nicht läuft) um `devcontainer up` und `down` |
| `CHANGELOG.md` | update | `[Unreleased]` / `Added` (gelesen: Abschnitt vorhanden, englisch, Kennungen erlaubt); Zusage-Menge und Einzelmessung benannt |
| `docs/maintainer/cli-json-contract.md` | update | neuer Abschnitt analog §6.10 (gelesen: `egress list --json`), Zeile in §7 (Mutations-Matrix, alle Spalten „—“: der Befehl schreibt keine Datei) |
| `README.md`, `README.de.md` | prüfen | `README.md` Zeile 15 (gelesen) nennt „four disclosures“ unter `--quiet`; die Socket-Ausweisung gehört zur ersten, die Zahl bleibt vier; ob die Feature-Aufzählung den Befehl nennen soll, entscheidet der Implementer (nicht gelesen: `README.de.md`) |
| `internal/adapter/driving/cli/jsontestutil/` | **unverändert** | gelesen: `jsontestutil.go:324` nimmt jeden Code mit Präfix `LH-` an; der Code der Ausweisung ist die Kennung der Anforderung, kein neuer tool-interner Code |
| `Dockerfile` (Stufe `runtime`), `.github/workflows/`, `spec/` | **unverändert** | Runtime-Image unverändert (Entscheidung 8); `spec/` im CR- und im Spezifikations-Slice |
| `tools/messung/devcontainer-start-befehl/`, `docs/messungen/slice-v1-cli-devcontainer-up-down/` | neu | Einzelmessung auf Colima (Projektinhaber); Ablage nach [`harness/conventions.md`](../../../../harness/conventions.md#ablage-von-messprotokollen-und-skripten) |

**Docker-Tests mit Wächter** (die Wächter des Fundament-Slice gelten fort; hier die Zusagen am Befehl; der Test ruft den Use-Case mit den Produktions-Adaptern auf):

- **S1 Start eines erzeugten Projekts.** Je `init --devcontainer` (Standardprofil) und `--sandbox`: Exit 0, `containerId` läuft (`docker inspect`), `remoteUser` ist der erzeugte Benutzer, Sandbox: Statusdatei `result=ok` und keine Fallback-Ausweisung; mindestens ein Katalog-Feature, soweit der Lauf es trägt (sonst steht „Features nicht zugesagt“ im Handbuch). *Wächter:* ein `postStartCommand`, das mit Exit 1 endet, ergibt Exit 12.
- **S3 UID.** Mit `devcontainer.user.uid: 4321`: im Container hat der Benutzer die UID 4321, eine im Standardprofil vom Container aus angelegte Datei gehört im Host-Ordner der UID 4321; das Ergebnis von `updateRemoteUserUID` aus dem Hilfs-Container (der als `root` läuft) wird beobachtet und im Handbuch nur so beschrieben. *Wächter:* mit der Standard-UID 1000 ist der Eigentümer ein anderer.
- **S4 Neustart-Semantik.** `up` → `down` → `up`: gleiche `containerId`; `up --recreate`: neue `containerId`; `down` entfernt den Container nicht (`docker inspect` kennt ihn, `Running` falsch); `down` ohne Container und ein zweites `down`: Exit 0; im Sandbox-Profil überlebt ein Marker im benannten Volume `--recreate`. *Wächter:* ein Lauf, der das Volume mit `docker volume rm` entfernt, lässt den Marker-Test rot werden.
- **S8 Fehlerbild.** Docker-Binary nicht auf dem `PATH` → 11; `DOCKER_HOST=tcp://…` → 11; kein `devcontainer.json` → 10; ein Dockerfile mit `RUN false` → 12 mit dem Ende der CLI-Ausgabe in der Meldung; `postStartCommand: exit 3` → 12; unter `--json` stimmt `exitCode` mit dem Prozess-Exit überein; welche Ausgabe der CLI auf stdout und welche auf stderr liegt, wird beobachtet und ins Handbuch nur so aufgenommen. *Wächter:* die Tabelle „Fall → erwarteter Exit“ steht im Test, ein abweichender Exit bricht ihn.
- **S9 Editor-Vergleich.** `docker inspect` des gestarteten Containers: `devcontainer.local_folder` = Basisverzeichnis, `devcontainer.config_file` = Pfad der Konfiguration (beobachtet, nicht angenommen), kein `u-boot`-eigenes Label ([`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)); ob VS Code dieselben Werte setzt, bleibt ausgewiesene Grenze (nicht gemessen). *Wächter:* ein Container mit anderem Label-Wert wird von `FindContainer` nicht gefunden.
- **Einzelmessung Colima (Projektinhaber).** Derselbe Ablauf wie S1 und S4 mit dem ausgelieferten Binary auf macOS/Colima; das Protokoll nennt Colima- und Docker-Version, Architektur und Projektpfad (im geteilten Home und, als Probe der Grenze, außerhalb); das Ergebnis trägt die Zusage-Zeile „Colima“ oder nimmt sie zurück.

## 4. Trigger

**Start** (`open` → `in-progress`, ggf. über `next`): `Verantwortlich:` gesetzt; [`slice-v1-devcontainer-start-fundament`](../open/slice-v1-devcontainer-start-fundament.md) ist geschlossen; kein anderer Slice in `in-progress/` (WIP-Limit 1); der Projektinhaber steht für die Einzelmessung auf Colima bereit (die Einzelmessung ist Teil der Verification, nicht des Starts).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Docker-Tests S1, S3, S4, S8, S9 brauchen mehr als einen Slice an Test-Infrastruktur, oder die Einzelmessung verlangt eine neue Messreihe.
- `in-progress` → `open` (blockiert): ein Test widerlegt eine Zeile der Zusage (Standardprofil mit UID ≠ 1000 trägt nicht, Features nicht, `--recreate` entfernt das Volume, Colima trägt nicht); der Slice geht mit dem Befund an den Architect, und ein Folge-CR schneidet die Zusage zurück.

## 5. Closure-Trigger

DoD vollständig, Review-Report und Verification-Evidence liegen vor (die Einzelmessung auf Colima ist Teil der Evidence), jede Zeile der Zusage-Tabelle ([`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 9) hat einen grünen, nicht übersprungenen Test samt Wächter oder steht im Handbuch als „nicht zugesagt“, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Die Zusage reicht weiter als die Tests** (Colima ohne CI, Standardprofil mit UID ≠ 1000, Features). — **Ausgang:** weiter offen; Tests und Einzelmessung entscheiden, ein widerlegender Befund geht an den Architect ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)).
- **Der Erfolgsumschlag trägt immer `status: warn`**, weil die Socket-Ausweisung eine Warnung ist (wie jede Warnung: Exit-Code `0`). Ein Skript, das auf `status: ok` prüft, schlägt an. — **Ausgang:** weiter offen; Handbuch und JSON-Kontrakt sagen es; ein Wechsel auf einen Hinweis-Eintrag wäre eine Lastenheft-Änderung ([`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) kennt nur `warn` und `error`).
- **`down` heißt „stoppen“, `u-boot down` heißt „abbauen“.** — **Ausgang:** weiter offen; die Hilfe sagt es; die Alternative `stop` ist Freigabe-Punkt b des ADR.
- **Die Fallback-Ausweisung nach `up` liest `u-boot.yaml`**; ein Projekt ohne die Datei zeigt keine. — **Ausgang:** weiter offen; der Test S1 fixiert das Verhalten, das Handbuch nennt die Grenze.
- **Hinweistexte nennen `devcontainer up --remove-existing-container`** (`egress.go`, `generate.go`, gelesen), der Befehl heißt jetzt auch `u-boot devcontainer up --recreate`. — **Ausgang:** entfallen als Änderung dieses Slice: die gemessene, editor-neutrale Form bleibt ([`SPEC-046`](../../../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung) und Tests pinnen sie); das Handbuch nennt die Entsprechung.
- **Die Einzelmessung hängt am Projektinhaber** (macOS/Colima). — **Ausgang:** weiter offen; ohne sie bleibt die Zeile „Colima“ unbelegt und der CR nimmt sie zurück.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `internal/hexagon/application` (APP), `internal/hexagon/port` (PORT, driving), `internal/adapter/driving/cli` (CLI), `internal/e2e` (E2E), `cmd/uboot` (Wiring, keine Sub-Area), `docs/` (DOK) und `tools/` (WZG), laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) alle Greenfield. Die Rückstufung der CLI-Sub-Area droht, wenn eine tragende Adapter-Konvention nur im Code entsteht: Die zwei Klassifikatoren für `ErrDevcontainerRuntime` und die Usage-Sentinels gehören deshalb in die Sicht (Spezifikations-Slice), nicht nur in den Code.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (die „prüfen“-Zeilen oben), [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (jede Bedingung im Handbuch hat einen Test), [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md), [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md), [`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md) (`GIT_TOKEN` in der Einzelmessung nur als Platzhalter), [`BEO-WZG/host-werkzeug-ausserhalb-docker`](../observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md), [`BEO-CLI/kommentar-beschreibt-abwesendes`](../observations/BEO-CLI/kommentar-beschreibt-abwesendes/observation.md) (Kommentare beschreiben, was da ist, `AGENTS.md` §3.7). Evidenz-/Diskrepanz-Risiko: mittel (Standardprofil, Colima, Features).

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
