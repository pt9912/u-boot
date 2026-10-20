# Slice v1: Devcontainer ohne Editor und ohne Node/CLI auf dem Host starten — Entscheidung und Entwurf (`LH-FA-DEV-001`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-devcontainer-start-ohne-editor`](../welle-devcontainer-start-ohne-editor.md) (Eltern-Slice der Welle). Der Slice ist ein eigener Slice **nach** [`welle-sandbox-egress-ollama`](../done/welle-sandbox-egress-ollama.md) (Entscheidung des Projektinhabers, 2026-10-18; Reihenfolge wie [`slice-v1-egress-dns-upstream-schliessen`](../done/slice-v1-egress-dns-upstream-schliessen.md)) und hält deren Closure und das Release nicht auf; er hat keine Closure-Bedingung, die von seiner DoD verschieden wäre.

**Bezug:** [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) (Devcontainer erzeugen: Mindestdateien), [`LH-FA-DEV-002`](../../../../spec/lastenheft.md#lh-fa-dev-002--vs-code-kompatibilität) (VS-Code-Kompatibilität), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil: „keinen Container-Runtime-Socket des Hosts einbinden“; Lockerungen in der Befehlsausgabe ausweisen), [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) (Befehlsübersicht), [`LH-NFA-PORT-002`](../../../../spec/lastenheft.md#lh-nfa-port-002--keine-unnötigen-systemabhängigkeiten) (möglichst wenige externe Systemabhängigkeiten), [`LH-NFA-PORT-004`](../../../../spec/lastenheft.md#lh-nfa-port-004--ausgeliefertes-container-image) (Runtime-Image minimal, keine Shell, keine Build-Toolchain), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (keine verdeckte Ausführung fremden Codes), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zugangsdaten nie in Image, Volume oder erzeugter Datei), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (Exit-Codes); [`ADR-0007`](../../adr/0007-distributionswege-ghcr.md) und [`ADR-0004`](../../adr/0004-ci-system.md) (Distribution und CI, `Accepted`), [`ADR-0019`](../../adr/0019-docker-daemon-aus-docker-apt-repo.md), [`ADR-0023`](../../adr/0023-podman-aus-podman-static-im-sandbox-devcontainer.md), [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) (Muster für Pin und Lieferkette, `Accepted`), [`ADR-0022`](../../adr/0022-fallback-ausweisung-statusdateien-im-container.md) (Host liest den Devcontainer-Container über Label und `docker cp`, `Accepted`), [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md) (Neustart und Neuanlage mit der Dev-Container-CLI, `Accepted`).
Anlass: Entscheidung des Projektinhabers (2026-10-18): Zum Starten des Devcontainers braucht man heute einen Devcontainer-fähigen Editor oder die `devcontainer`-CLI (npm, also Node); der Wunsch: „können wir das nicht besser machen. kann man das in das u-boot image einbauen, damit man mit weniger tools auskommt“.

**Berührte Spec-Stellen:** [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) (der Host findet den Container über das Label `devcontainer.local_folder`; der Architect liest, wie sich ein vom neuen Weg gestarteter Container dazu verhält), [`SPEC-042`](../../../../spec/spezifikation.md#spec-042--aufruf-der-compose-engine-und-preflight) (Aufruf der Compose-Engine und Preflight; Titel gelesen, Inhalt nicht — der Architect liest, ob der Preflight für einen Port zum Container-Start gilt), `architecture.md` §2 (Ports `DockerProbe`, `DockerEngine`, `DevcontainerInspector`) und §3 (`ARC-008`, Container-Runtime als externe Abhängigkeit). Der Slice **liest** sie und ändert keine; ist nach dem ADR eine Spec-Änderung nötig, benennt der Architect die Stelle, und ein eigener Spezifikations-Slice folgt dem Code. Der Verweis zeigt aufwärts: Die Spec nennt diesen Slice nie.

**Verantwortlich:** pt9912

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-20.

---

## 1. Ziel und Abgrenzung

**Ziel:** Entscheiden und entwerfen, wie `u-boot` den Devcontainer eines Projekts ohne Editor und ohne Host-Installation von Node oder der Dev-Container-CLI startet, sodass auf dem Host nur Docker nötig bleibt. Der Weg ist vorgegeben (Hilfs-Image mit der Dev-Container-CLI, lokal aus eingebettetem Dockerfile gebaut, Docker-Socket nur für den Befehl, siehe *Entscheidungen des Projektinhabers*); der Slice liefert das **ADR** (`Accepted`), den **Wortlaut-Vorschlag** für das Lastenheft und die **Folge-Slices als Dateien**. **Kein Produktcode, keine neue Messung:** Belegbasis ist das vorhandene Protokoll; was es nicht trägt, steht im ADR als „nicht gemessen, nicht zugesagt“ und wird im Umsetzungs-Slice mit Docker-Tests geprüft.

### Gefundener Stand (gelesen am 2026-10-20)

- **`u-boot` erzeugt nur Dateien, es startet den Devcontainer nicht.** Handbuch `docs/user/benutzerhandbuch.md` Abschnitt 4.9, *Voraussetzung*: „Zum Starten des Devcontainers brauchen Sie einen Devcontainer-fähigen Editor (zum Beispiel VS Code) oder die `devcontainer`-CLI.“ Abschnitt *Ergebnis*: `.devcontainer/devcontainer.json` und `.devcontainer/Dockerfile` entstehen. [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) nennt als Mindestdateien nur diese beiden. Die Befehle des Binary sind `init`, `doctor`, `add`, `remove`, `up`, `down`, `logs`, `generate`, `config`, `template`, `egress` (`internal/adapter/driving/cli/root.go`, die `AddCommand`-Zeilen in `buildRootCommand`); die Befehlsübersicht [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) führt `u-boot egress list` (V2) als jüngste Zeile und keinen Befehl, der einen Devcontainer startet. `u-boot up` betrifft den Compose-Stack: `internal/adapter/driving/cli/up.go`, Langtext von `newUpCommand` („Bring the Compose environment defined in compose.yaml up via docker compose up -d“); [`LH-FA-UP-001`](../../../../spec/lastenheft.md#lh-fa-up-001--umgebung-starten) beschreibt dasselbe.
- **Docker-Aufrufe des Hosts sitzen in einem Adapter-Paket.** `internal/hexagon/port/driven/docker.go` (`DockerProbe`), `docker_engine.go` (`DockerEngine`: `docker compose …`, Sentinels `ErrDockerUnavailable` und `ErrComposeRuntime`; Abbildung auf Exit-Code 11 und 12 in `internal/adapter/driving/cli/composesentinel.go`) und `devcontainer_inspector.go` (`DevcontainerInspector`: `FindContainer` über das Label `devcontainer.local_folder=<absoluter Ordner>`, `ReadFile` über `docker cp`, Sentinel `ErrContainerEngineNotLocal`); Implementierung in `internal/adapter/driven/docker/` (`engine.go`, `inspector.go`; `engine.go` ist für diesen Plan nicht gelesen). Der Inspector weist eine entfernte Engine ab (`requireLocalEngine` in `inspector.go`: `DOCKER_HOST` oder `docker context inspect`, nur `unix://` gilt als lokal). Einen Port, der einen Container **startet**, gibt es nicht.
- **Was das erzeugte `devcontainer.json` von dem Werkzeug verlangt** (`internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl`): `build.dockerfile` mit Kontext `.` und optional `USER_UID` als Build-Argument; im Sandbox-Profil `workspaceMount: ""`, `workspaceFolder`, `runArgs`, `mounts` mit benannten Volumes, `postCreateCommand`, `postStartCommand`, `remoteEnv` mit `${localEnv:GIT_TOKEN}`; `remoteUser: vscode`. Die Volume-Namen tragen `${devcontainerId}` (`internal/hexagon/application/devcontainer_sandbox.go`, `volumeMount`). `updateRemoteUserUID` und `initializeCommand` kommen weder in der Vorlage noch in `internal/` vor (Suche über beide). Im Standardprofil bindet das Werkzeug den Projektordner ein (Handbuch 4.9, *Hinweise*) — der **Host-Pfad** spielt dort eine Rolle.
- **Das Runtime-Image trägt weder Docker noch Shell.** `Dockerfile`, Stufe `runtime`: `FROM gcr.io/distroless/static-debian12:nonroot`, nur `COPY … /usr/local/bin/u-boot`, `USER 65532:65532`. [`LH-NFA-PORT-004`](../../../../spec/lastenheft.md#lh-nfa-port-004--ausgeliefertes-container-image) schreibt „minimal … keine Shell und keine Build-Toolchain, … Non-root-Benutzer“ vor. Ein Befehl, der Docker aufruft, läuft deshalb im Image so wenig wie `up`, `down` und `logs` (Handbuch Abschnitt *Container-Image* nennt für das Image die Einschränkungen). Die Stufe `test-docker-tools` desselben `Dockerfile` bringt `docker-ce-cli` und das Compose-Plugin aus dem Docker-APT-Repo, `repo.mk` (`test-docker`) hängt `/var/run/docker.sock` in einen Container: Muster im Repo vorhanden.
- **Distribution.** [`ADR-0007`](../../adr/0007-distributionswege-ghcr.md): GHCR-Image `ghcr.io/pt9912/u-boot` plus Binary; `npm package` ist dort **verworfen** („Go-Binary ist kein Node-Modul“) — das betrifft die Distribution von `u-boot`, nicht den Einsatz der Dev-Container-CLI als Werkzeug. Ein Binary aus Homebrew oder `.deb` trägt den Quelltext nicht; ein Dockerfile, das `u-boot` bauen soll, muss deshalb im Binary eingebettet sein (`go:embed`, wie die Vorlagen unter `internal/hexagon/application/templates/`).

### Belegbasis und Grenzen

**Belegbasis ist das Messprotokoll** [`docs/messungen/slice-v1-sandbox-egress-neuladen-entwurf`](../../../messungen/slice-v1-sandbox-egress-neuladen-entwurf/README.md) (Skripte unter `tools/messung/sandbox-egress-neuladen-entwurf/`, Messschritt 4 = `05-m4-neustart.sh`, Rohdaten `rohdaten-05-m4.txt`). Es gibt **keine neue Messreihe**. Messstand: Dev-Container-CLI 0.89.0 (`@devcontainers/cli 0.89.0`, Node.js v22.23.3) in einem Wegwerf-Container (`Dockerfile.cli`: `FROM docker:29-cli`, `FROM node:22-bookworm`), Colima 0.10.3, Docker 29.5.2, arm64, **ein Lauf je Fall**; Aufruf `docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v "$w":/w <cli-image> devcontainer up --workspace-folder <ordner>` (Funktion `dc` in `05-m4-neustart.sh`). Das Skript reicht nur die letzten drei Zeilen der zusammengeführten Ausgabe weiter (`tail -n 3`).

| Aussage | Stand | Wo geprüft |
|---|---|---|
| `up` legt den Container an; Ausgabezeile mit `outcome`, `containerId`, `remoteUser`, `remoteWorkspaceFolder` (`rohdaten-05-m4.txt`); Sandbox-Profil mit Egress | **gemessen** (Messstand oben) | ADR (Belegbasis) |
| Stop, dann `up`: derselbe Container (gleiche `containerId`), altes Image; `up --remove-existing-container`: neuer Container, neues Image (`rohdaten-05-m4.txt`; Ergebnis in [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md)) | **gemessen** (Messstand oben) | ADR (Belegbasis) |
| `postStartCommand` läuft beim Start | **gemessen** für das Egress-Script des Sandbox-Profils; `postCreateCommand` und ein Katalog-Feature: nicht gelesen | ADR nennt die Grenze |
| Standardprofil (Bind-Mount des Projektordners), `updateRemoteUserUID`, `USER_UID` auf einem Linux-Host mit anderer UID (S1, S3) | **nicht gemessen, nicht zugesagt** | Umsetzungs-Slice, Docker-Test mit Wächter |
| Pfad-Abbildung: Host-Pfad im Hilfs-Container, Wert von `devcontainer.local_folder`, `FindContainer` findet den Container (S2). Im Protokoll lag das Projekt im Hilfs-Container unter `/w/proj`, das Label trägt diesen Pfad, nicht den des Hosts | **nicht gemessen, nicht zugesagt** | Umsetzungs-Slice, Docker-Test mit Wächter |
| Socket-Haltung: der **gestartete** Dev-Container trägt den Socket nicht; der Socket bleibt auf die Lebensdauer des Hilfs-Containers (`--rm`) begrenzt; Zugangsdaten nur als Umgebung des Befehls (S6) | **nicht gemessen, nicht zugesagt** | Umsetzungs-Slice, Docker-Test mit Wächter |
| Lieferkette: Pin von `@devcontainers/cli`, Digest-Pins der Basis-Images, Integritätsprüfung, Image-Größe, Bauzeit, Image-Scan (S5) | **nicht gemessen** (nur Version 0.89.0 und die Basis-Images des Messstands sind belegt) | ADR legt die Pin-Form fest (Muster [`ADR-0023`](../../adr/0023-podman-aus-podman-static-im-sandbox-devcontainer.md), [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md)); Zahlen und Pins im Umsetzungs-Slice |
| Andere Engines als Colima: Linux-Host, Docker Desktop, Podman, entfernte Engine, Windows, `amd64` (S7) | **nicht gemessen, nicht zugesagt** | Umsetzungs-Slice prüft, was seine Docker-Tests tragen; die Zusage-Menge folgt dem Ergebnis |
| Fehlerbild der CLI (fehlendes Docker, falscher Pfad, Build-Fehler, `postStartCommand`-Fehler) und Abbildung auf Exit-Code 11/12 (S8) | **nicht gemessen**; die Abbildung ist Entwurf im ADR | Umsetzungs-Slice, Docker-Test |
| Editor-Vergleich: Labels, `${devcontainerId}`, Volume-Namen wie bei VS Code (S9) | **nicht gemessen, nicht zugesagt** | Umsetzungs-Slice, soweit ohne Editor prüfbar; sonst ausgewiesene Grenze |
| Kalt-/Warm-Dauer des Starts | **nicht gemessen** | — |

Eine **Zusage** im ADR und im Lastenheft-Wortlaut deckt nur, was Belegbasis oder ein Test des Umsetzungs-Slice trägt; alles andere heißt „nicht gemessen, nicht zugesagt“. Ein Umsetzungs-Slice, dessen Docker-Test eine Zeile der Tabelle widerlegt, geht zurück an den Architect (ADR ergänzen), nicht in eine stille Änderung der Zusage.

### Entscheidungen des Projektinhabers (beantwortet am 2026-10-20)

1. **Reichweite:** Linux und macOS mit Colima; Windows ist nicht zugesagt.
2. **Beschaffung:** Das Hilfs-Image baut `u-boot` lokal aus einem **eingebetteten** Dockerfile (`go:embed`, Pins) — Variante **a2**; ein GHCR-Artefakt (a1) entfällt.
3. **Socket:** Der Docker-Socket darf in einen kurzlebigen Hilfs-Container (`--rm`) und wird in der Befehlsausgabe als Sicherheitslockerung ausgewiesen.
4. **Befehlsform:** `u-boot devcontainer up` / `u-boot devcontainer down` (Vorschlag; das ADR konkretisiert Flags, Exit-Codes 11 und 12, `--quiet` und `--json`). Priorität und Release-Zuordnung schlägt der Planner vor (unten), entscheidet der Projektinhaber.
5. **Variante a′** (vorhandene `devcontainer`-CLI bevorzugen) entfällt zunächst: immer das Hilfs-Image.

### Optionen (Stand nach den Entscheidungen)

| Opt. | Weg | Stand |
|---|---|---|
| **a2** | Hilfs-Image mit der Dev-Container-CLI, lokal aus eingebettetem Dockerfile gebaut, `docker run --rm` mit Docker-Socket und eingehängtem Arbeitsordner | **gewählt** (Entscheidung 2 und 3); das ADR begründet Pin, Socket-Haltung und Grenzen |
| a1 | Hilfs-Image als GHCR-Artefakt (neuer Release-Schritt, `docs/maintainer/releasing.md` §1) | nicht gewählt (Entscheidung 2); das ADR nennt den Grund |
| a′ | vorhandene `devcontainer`-CLI auf dem `PATH` bevorzugen | zunächst nicht (Entscheidung 5); zwei Pfade mit gleicher Zusage wären doppelt zu prüfen |
| b | CLI und Docker im veröffentlichten Runtime-Image | Vergleich im ADR: widerspricht dem Wortlaut von [`LH-NFA-PORT-004`](../../../../spec/lastenheft.md#lh-nfa-port-004--ausgeliefertes-container-image) (minimal, keine Shell) und hilft nur Nutzern des Images |
| c | eigene Go-Implementierung des Starts ohne CLI | Vergleich im ADR: Features, Lifecycle-Hooks, `${devcontainerId}` und `updateRemoteUserUID` nachzubauen ist eine Größenordnung mehr Aufwand und driftet von der Referenz-CLI; gehört bei Wahl zurück zur Zerlegung |
| d | Status quo, Handbuch nennt die Voraussetzung | Vergleich im ADR: der Wunsch bliebe unerfüllt |

Eine Option **e** (Socket-Proxy mit eingeschränkter API) ist nicht im Plan; der Architect nimmt sie als Vergleich auf, wenn das ADR die Socket-Reichweite für die Zusage zu weit findet.

### Vorschlag des Planners: Priorität und Release-Zuordnung (Entscheidung beim Projektinhaber)

- **Priorität V2.** Der Befehl ist Komfort: Der Start mit Editor oder `devcontainer`-CLI funktioniert, `u-boot` erzeugt die Dateien bereits ([`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) MVP). Vergleichbar steht `u-boot egress list` mit V2 in [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle). Die neue Zusage lockert den Socket (Entscheidung 3), das spricht gegen V1.
- **Release-Zuordnung:** offen (Entscheidung des Projektinhabers); Vorschlag: im nächsten Minor-Release nach dem Code-Slice, getragen von [`welle-devcontainer-start-ohne-editor`](../welle-devcontainer-start-ohne-editor.md); weder die Welle `welle-ollama-cloud-modelle` noch `welle-forge-cli-gh-glab` (Roadmap, *Nächste Wellen*) hängt daran, und der Befehl hält keine von beiden auf. Die Nummer vergibt der Release-Cut.

### Größenprüfung (Planner)

Der Slice liefert ein ADR, einen Wortlaut-Vorschlag und Folge-Slice-Dateien mit Plan-Tabellen; er bleibt klein genug, solange das ADR ohne neue Messung auskommt und der Architect höchstens **fünf** Folge-Slices schneidet (Mindestzuschnitt: Lastenheft-CR · Code · Spezifikation; Code darf in Fundament — Port, Adapter, eingebettetes Dockerfile, Pin — und Befehl zerfallen). Wächst eines davon, gilt die Rückführung in §4; der Slice erweitert sich nicht still.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Produktcode, Befehl, Port, eingebettetes Dockerfile, Test, Workflow-Änderung** — dieser Slice entscheidet und entwirft; die Zusage ist neu (ein Befehl, der Container startet), und das Lastenheft ist das Zielbild: Wortlaut erst nach Freigabe des Projektinhabers. Umsetzung: Folge-Slices nach dem ADR (§7, *Folge-Slices*).
- **Neue Messreihen und Messskripte** — Entscheidung des Projektinhabers (2026-10-20); was das Protokoll nicht trägt, prüft der Umsetzungs-Slice als Docker-Test mit Wächter.
- **Änderung von `spec/lastenheft.md`, `spec/spezifikation.md`, `spec/architecture.md`** — der Slice legt einen Wortlaut-Vorschlag vor; geschrieben wird er in einem eigenen CR-Slice nach der Freigabe ([`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md) ist das Muster dieser Trennung).
- **Handbuch und `CHANGELOG.md`** — Abschnitt 4.9 beschreibt den Ist-Stand richtig (der Befehl existiert nicht); die Änderung gehört zur Lieferung des Befehls.
- **Option c (Nachbau ohne CLI) als Umsetzung** — nur als Vergleich im ADR.
- **VS Code, JetBrains oder andere Editoren** — nicht der Gegenstand.
- **`u-boot up` startet den Devcontainer mit** — `up` betrifft den Compose-Stack ([`LH-FA-UP-001`](../../../../spec/lastenheft.md#lh-fa-up-001--umgebung-starten)); ein gemeinsamer Start wäre eine Änderung dieser Zusage und braucht eine eigene Entscheidung.
- **Windows-Host, Docker Desktop, Podman als Host-Engine, entfernte Engine, `amd64`** — nicht gemessen, nicht zugesagt (Tabelle oben).
- **Einen Agenten im Devcontainer starten** — Nicht-Ziel des Sandbox-Profils ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil): „`u-boot` startet keinen Agenten“).

## 2. Definition of Done

Die Liefer-Punkte sind drei (ADR · Lastenheft-Wortlaut-Vorschlag · Folge-Slices); die Gate-Läufe und Closure-Pflichten darunter zählen nicht mit.

- [ ] **Entscheidungs-ADR (Rolle Architect):** ein neues ADR (nächste freie Nummer, MADR-Form; `Accepted` durch den Projektinhaber, erfolgt; die Review-Befunde F-2 und F-5 behebt die Folge-ADR [`ADR-0029`](../../adr/0029-devcontainer-start-exit-11-und-socket-haltung.md), `Accepted`) hält fest: Wahl a2 mit Vergleich der Optionen a1, a′, b, c, d; Belegbasis und die Liste „nicht gemessen, nicht zugesagt“ aus der Tabelle oben; Lieferkette (Pin-Form für `@devcontainers/cli` und die Basis-Images, Integritätsprüfung, Pflege) nach dem Muster von [`ADR-0023`](../../adr/0023-podman-aus-podman-static-im-sandbox-devcontainer.md) und [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md); Socket-Haltung (kurzlebig, `--rm`, nicht im gestarteten Container) und ihre Ausweisung in der Befehlsausgabe; Befehlsform `u-boot devcontainer up` / `down` (Flags, Exit-Codes 11 und 12, `--quiet`, `--json`); Port im Hexagon (I/O im Driven-Adapter); Zusage-Menge, die dem Wortlaut „Linux und macOS mit Colima“ nur so weit folgt, wie Belegbasis und Tests des Umsetzungs-Slice sie tragen. `docs/plan/adr/README.md` führt es im Index.
- [ ] **Lastenheft-Wortlaut-Vorschlag** in der Closure-Notiz (§7), als Alt/Neu je Stelle: die neue Anforderung (Kennung nach Schema des Lastenhefts durch den Architect), die Zeile(n) in [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) samt Priorität, der Bezug zu [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) und [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (der Hilfs-Container ist nicht das erzeugte Ergebnis, die Lockerung wird ausgewiesen), ggf. [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) zu `--quiet`. Das Lastenheft trägt keine Messstände (nur Zusage-Mengen). Geschrieben wird er erst im CR-Slice nach Freigabe.
- [ ] **Folge-Slices als Dateien in `open/`** (Kennungen vergibt der Architect nach dem ADR; mindestens Lastenheft-CR · Code · Spezifikation, höchstens fünf, Handbuch und `CHANGELOG.md` im Code-Slice), jeder mit Plan-Tabelle (§3.10 gelesen, nicht vermutet), Bezug, DoD und Trigger; `docs/plan/planning/open/README.md` führt sie im Bestand. Der Code-Slice nennt die Docker-Tests mit Wächter für S1, S2, S3, S4, S5, S6, S7, S8 und S9 aus der Tabelle oben.
- [ ] **Priorität und Release-Zuordnung** sind vom Projektinhaber entschieden und im ADR oder in §7 vermerkt (Vorschlag oben).
- [ ] `make gates` grün (Exit-Code von `make` selbst), `make docs-check` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: ADR-Index, `open/README.md`, Roadmap-Zeile „In Arbeit“ bei der Closure; Handbuch und `CHANGELOG.md` bleiben unverändert (öffentlicher Vertrag nicht berührt).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der Slice-Closure geprüft; die Welle prüft sie bei ihrer Closure für die Welle als Ganzes.

## 3. Plan (vor Code)

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; der Implementer und der Architect ergänzen sie in ihrem Lauf und gleichen sie nach der letzten Fix-Runde gegen `git diff --stat` ab (`AGENTS.md` §3.10). Eine Zeile „unverändert“ trägt ihre gelesene Fundstelle.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/adr/NNNN-….md` | neu | Entscheidungs-ADR (Rolle Architect), `Accepted` |
| `docs/plan/adr/README.md` | update | ADR-Index |
| Folge-Slice-Dateien in `docs/plan/planning/open/` (Kennungen vergibt der Architect) | neu | Lastenheft-CR, Code, Spezifikation (höchstens fünf), je mit Plan-Tabelle |
| `docs/plan/planning/open/README.md` | update | Bestand der Folge-Slices |
| diese Slice-Datei (§7 Wortlaut-Vorschlag, Closure-Notiz) | update | Closure; Moves per `git mv` |
| `docs/plan/planning/in-progress/roadmap.md` | update | Zeile „In Arbeit“ bei der Closure (heute: Zeile 20 nennt diesen Slice) |
| `docs/reviews/` | neu (Reviewer, Verifier) | Review und Verifikation |
| `docs/user/benutzerhandbuch.md`, `CHANGELOG.md` | **unverändert** | Abschnitt 4.9 *Voraussetzung* (Zeilen 692 bis 696) gelesen: beschreibt den Ist-Stand ohne Befehl; die Änderung gehört in den Code-Slice |
| `Dockerfile`, `repo.mk`, `internal/`, `cmd/`, `.github/workflows/`, `tools/`, `docs/messungen/`, `spec/` | **unverändert** | kein Code, keine Messung, kein Lastenheft-Wortlaut in diesem Slice; `Dockerfile` Stufe `runtime` gelesen (Zeilen 175 bis 192: distroless, `COPY`, `USER 65532:65532`), die übrigen Pfade ändert der Plan nicht |

## 4. Trigger

**Start** (`open` → `in-progress`, ggf. über `next`): `Verantwortlich:` gesetzt; [`welle-sandbox-egress-ollama`](../done/welle-sandbox-egress-ollama.md) ist geschlossen (Reihenfolge-Entscheidung des Projektinhabers); kein anderer Slice in `in-progress/` (WIP-Limit 1). Der Start hängt nicht am Release. Rollenfolge: Planner → Architect (ADR, Wortlaut-Vorschlag, Folge-Slices) → Reviewer → Verifier → Planner (Closure).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): das ADR braucht für die Wahl eine neue Messung (der Slice hat dann seine Voraussetzung gewechselt), oder der Architect schneidet mehr als fünf Folge-Slices, oder wählt Option c.
- `in-progress` → `open` (blockiert): die Freigabe des Projektinhabers zu Priorität oder Release-Zuordnung steht aus und blockiert das ADR; die Evidence nennt den Grund (kein stilles Rot).

## 5. Closure-Trigger

DoD vollständig, Review-Report und Verification-Evidence liegen vor, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7); die Entscheidung des Projektinhabers zu Priorität und Release-Zuordnung ist im Gespräch belegt und in §7 vermerkt; kein Lastenheft-Wortlaut ohne Freigabe geschrieben.

## 6. Risiken und offene Punkte

- **Der Docker-Socket im Hilfs-Container ist Root auf dem Host-Docker.** — **Ausgang:** weiter offen; Haltung (kurzlebig, ausgewiesen) ist entschieden (Entscheidung 3), die Prüfung „der gestartete Container trägt den Socket nicht“ steht im Docker-Test des Umsetzungs-Slice (Tabelle, S6).
- **Die Pfad-Abbildung stimmt auf Hosts außerhalb des Messstands nicht** (Windows, Docker Desktop, Colima mit Projekt außerhalb des Homes, entfernte Engine); ein falscher Pfad kann still einen falschen Bind-Mount erzeugen. — **Ausgang:** weiter offen; Docker-Test mit Wächter im Umsetzungs-Slice (S2, S7), Zusage nur für Getestetes.
- **Neue Lieferkette** (npm-Paket mit Abhängigkeiten, Node-Basisimage, Docker-CLI-Quelle) und Pin-Pflege (neue Advisories im Image-Scan). — **Ausgang:** weiter offen; Pin-Form legt das ADR fest, Pins und Zahlen der Umsetzungs-Slice (S5).
- **Ein per `u-boot` gestarteter Container weicht vom Editor-Start ab** (Labels, `${devcontainerId}`, Volumes); zwei getrennte Umgebungen wären die Folge. — **Ausgang:** weiter offen; ohne VS Code im Test bleibt es eine ausgewiesene Grenze (S9).
- **Zusage-Menge gegen Belegbasis:** Entscheidung 1 nennt Linux und macOS mit Colima als Reichweite, das Protokoll belegt nur Colima/arm64. — **Ausgang:** weiter offen; das ADR formuliert die Zusage-Menge nach dem, was der Umsetzungs-Slice testet (Linux in den Docker-Tests der CI, macOS/Colima ohne CI-Lauf), und weist die Differenz aus.
- **Eine neue Zusage ohne Freigabe:** ein Befehl, der Container startet, ist im Lastenheft nicht vorgesehen. — **Ausgang:** weiter offen; Wortlaut-Vorschlag in §7, Freigabe des Projektinhabers vor dem CR-Slice ([`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md): Vertragsänderungen fielen im Review auf).
- **Der Slice wächst in Richtung Umsetzung oder Messung** (Option c, Code „nur zum Probieren“, Messreihen „zur Sicherheit“). — **Ausgang:** weiter offen; Rückführung `in-progress` → `next` (§4), Code und Docker-Tests gehören in einen Folge-Slice.
- **Das ADR sagt mehr zu als das Protokoll trägt** (Ist-Stand-Aussagen über andere Engines, Pfade, Dauer). — **Ausgang:** weiter offen; der Reviewer prüft jede Aussage des ADR gegen die Belegbasis-Tabelle (`AGENTS.md` §3.10).

### Offene Fragen an den Projektinhaber

1. **Reichweite „nur Docker“** — beantwortet 2026-10-20: Linux und macOS mit Colima, Windows nicht zugesagt.
2. **Beschaffung des Hilfs-Images** — beantwortet 2026-10-20: lokal aus eingebettetem Dockerfile (a2).
3. **Socket** — beantwortet 2026-10-20: ja, kurzlebiger Hilfs-Container (`--rm`), als Sicherheitslockerung in der Befehlsausgabe ausgewiesen.
4. **Befehlsform** — beantwortet 2026-10-20: `u-boot devcontainer up` / `down`, Details im ADR. **Priorität und Release-Zuordnung** — offen: Vorschlag V2, nächster Minor-Release nach dem Code-Slice (§1); Entscheidung beim Projektinhaber.
5. **Vorhandene CLI bevorzugen (a′)** — beantwortet 2026-10-20: nein, immer das Hilfs-Image.

## 7. Closure-Notiz

*(Die Closure-Bullets — was funktioniert hat, Steering-Loop-Eintrag, Register, Risiken, Paarungen — füllt der Planner bei der Closure. Die drei Abschnitte darunter sind das Ergebnis der Rolle Architect und die Liefer-Punkte der DoD.)*

### Entscheidungs-ADR

[`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md), `Accepted` (Freigabe-Punkte a bis i durch den Projektinhaber). Gewählt ist a2 mit Vergleich gegen a1, a′, b, c, d und e. **Verlagerung von S1, S3, S4 und S8 in die Docker-Tests des Umsetzungs-Slice: bestätigt.** S4 ist am Messweg teilweise belegt (Stop/`up`, `--remove-existing-container`), die Tests prüfen ihn am Befehl und ergänzen `down` und das Überleben eines benannten Volumes; die DoD-Zeile „Folge-Slices“ nannte S4 nicht und nennt ihn jetzt. Die Verlagerung von S2, S5, S6, S7 und S9 bleibt wie in §1.

### Lastenheft-Wortlaut-Vorschlag (Alt/Neu je Stelle; geschrieben wird er erst im CR-Slice nach Freigabe)

Gelesen am 2026-10-20 in `spec/lastenheft.md`: [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) (Mindestdateien, der Befehl erzeugt keine Datei), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (letzte der Familie `LH-FA-DEV`, die neue Anforderung bekommt die Nummer 010), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle), Versionsfeld `0.17.0`, §16 Historie (letzte Zeile `0.17.0`). Das Lastenheft trägt keine Messstände und keine Prozesszustände; die Zusage-Menge folgt dem Wortlaut von [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Version 0.7.0: „unter Docker auf Linux und unter Colima (Linux-VM)“). Die Priorität `V2` und die Versionsnummer (nächste freie Minor-Version zum Zeitpunkt des CR; `0.18.0` ist durch [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md) vorgesehen) sind offen.

**Stelle 1 — neue Anforderung** in §4.3, nach [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer). Alt: — (neu). Neu:

> ### LH-FA-DEV-010 – Devcontainer ohne Editor starten <!-- d-check:ignore (Kennung des Vorschlags, im Lastenheft erst nach dem CR) -->
>
> *Priorität: V2.* Das Produkt soll den erzeugten Devcontainer eines Projekts starten und stoppen können, ohne dass auf dem Host ein Devcontainer-fähiger Editor, Node.js oder die Dev-Container-CLI installiert ist; auf dem Host genügen Docker (Client und lokale Engine).
>
> - `u-boot devcontainer up` startet den Devcontainer des Projekts im Basisverzeichnis mit der Dev-Container-CLI, die `u-boot` in einem kurzlebigen Hilfs-Container ausführt. Das Hilfs-Image baut `u-boot` lokal aus gepinnten Quellen; der erste Aufruf braucht Netzzugriff. `--recreate` entfernt einen vorhandenen Container und legt ihn neu an; benannte Volumes bleiben erhalten. Nach dem Start weist der Befehl die Degradations-Fallbacks des Containers aus ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)).
> - `u-boot devcontainer down` stoppt den Container; er entfernt weder den Container noch seine Volumes.
> - Der Hilfs-Container bindet für die Dauer des Aufrufs den Container-Runtime-Socket des Hosts ein. Der Befehl weist diese Sicherheitslockerung in der Befehlsausgabe aus; der gestartete Devcontainer erhält den Socket durch diesen Befehl nicht. Das erzeugte Ergebnis nach [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) bleibt davon unberührt.
> - Zugangsdaten gibt der Befehl nur als Umgebungsvariablen weiter, auf die die erzeugte Konfiguration mit `${localEnv:<NAME>}` verweist ([`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)); sie stehen nie in einer Datei oder einem Image.
> - Fehlerklassen ([`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)): Exit-Code `10`, wenn das Basisverzeichnis keine Devcontainer-Konfiguration enthält; `11`, wenn Docker nicht erreichbar oder die Engine nicht lokal ist oder wenn ein Startscript des Containers wegen eines Umgebungsproblems abbricht (wie bei `doctor`); `12` bei Fehlern beim Bau oder Start, die kein Umgebungsproblem sind.
> - Der Befehl startet unter Docker auf einem Linux-Host mit lokaler Engine (`amd64`); alle übrigen Hosts und Engines, darunter Colima, Docker Desktop, rootless Docker, Podman als Host-Engine, eine entfernte Engine, Linux-`arm64` und Windows, sind nicht zugesagt. Die Fallback-Ausweisung nach dem Start setzt eine `u-boot.yaml` im Basisverzeichnis voraus. Der Exit-Code `11` für den Startscript-Abbruch setzt die Statusdatei des Containers voraus, die der Befehl über die `u-boot.yaml` im Basisverzeichnis findet: ohne diese Datei endet der Befehl bei einem Startscript-Abbruch mit `12`. Im ausgelieferten Container-Image ([`LH-NFA-PORT-004`](../../../../spec/lastenheft.md#lh-nfa-port-004--ausgeliefertes-container-image)) läuft er nicht.
> - Nicht-Ziele: `u-boot` startet keinen Agenten ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)), und der Befehl startet nicht den Compose-Stack ([`LH-FA-UP-001`](../../../../spec/lastenheft.md#lh-fa-up-001--umgebung-starten)).

Grenze des Wortlauts: Er nennt weder Standard- noch Sandbox-Profil einzeln und sagt „Linux-Host, `amd64`, lokale Engine“ als Menge, nicht als Messstand; Colima steht bewusst nicht darin (Entscheidung des Projektinhabers, 2026-10-20: erst zusagen, wenn das Messprotokoll der Einzelmessung vorliegt). Das weicht von [`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 9 und Freigabe-Punkt h ab (dort ist Colima als Einzelmessung zugesagt); die ADR bleibt unverändert (`AGENTS.md` §3.5), und für die Zusage gilt dieser Vermerk, bis der Folge-CR Colima aufnimmt. Die Stellen in Welle und Befehl-Slice sind entsprechend berichtigt. Liegt es vor und trägt, nimmt ein Folge-CR Colima in die Zusage auf; bis dahin sagt das Handbuch „nicht zugesagt“. Was von der Linux-Zusage trägt, entscheiden die Docker-Tests des Befehl-Slice; widerlegt eine Zeile die Zusage, geht der Befehl-Slice zurück an den Architect, und ein Folge-CR schneidet die Zusage zurück (Muster der Version 0.7.0 von [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)). „Weist aus“ ist hier gewollt: Die Ausweisung bleibt unter `--quiet` (Stelle 3).

**Stelle 2 — [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle), Befehlsübersicht** (zwei Zeilen nach `u-boot egress list`; Priorität nach der Entscheidung des Projektinhabers).

Alt: die Tabelle endet bei `u-boot egress list` (V2) und `u-boot template list` (V1). Neu:

> | `u-boot devcontainer up`     | Devcontainer ohne Editor starten   | V2        |
> | `u-boot devcontainer down`   | Devcontainer stoppen               | V2        |

**Stelle 3 — [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), `--quiet`** (zwei Satzteile).

Alt (a): „Unter `--quiet` geben die Subkommandos, die etwas ändern (`init`, `generate`, `add`, `remove`, `config set`, `up`, `down`), nur Fehler aus; …“ Neu: „… (`init`, `generate`, `add`, `remove`, `config set`, `up`, `down`, `devcontainer up`, `devcontainer down`) …“.

Alt (b), Ausnahme 1: „die Ausweisung von Sicherheitslockerungen (zusätzliche Capabilities, `--privileged`, Seccomp-/AppArmor-Anpassungen, zusätzliche Devices) und eines im Container aktivierten Login-Dienstes (Passwort-Login) in der Befehlsausgabe ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), …)“. Neu: „… zusätzliche Devices, der Einbindung des Container-Runtime-Sockets in den Hilfs-Container von `u-boot devcontainer up`) und eines im Container aktivierten Login-Dienstes … ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), …, `LH-FA-DEV-010`)“ <!-- d-check:ignore (Kennung des Vorschlags) -->. Die Ausnahmen 2 bis 4 bleiben; der Fortschritt von `devcontainer up` ist keine Ausnahme und entfällt unter `--quiet` bis auf Fehler.

**Stelle 4 — [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run), Wertebereich von `command`** (nur der Eintrag `devcontainer`; die Wertemenge steht nur dort, [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) verweist darauf).

Alt: „… `command` ist einer der Werte `init`, `add`, `remove`, `up`, `down`, `doctor`, `logs`, `generate`, `config`, `egress` oder `template`. …“ Neu: „… `generate`, `config`, `egress`, `devcontainer` oder `template`. …“. Der Befehl setzt `subcommand` (`up`, `down`); das Beispiel „wie `command == "template"` oder `command == "config"`“ bleibt ein Beispiel, wie beim Zusatz `egress`.

**Stelle 5 — [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), Abgrenzung** (ein Zusatz am Spiegelstrich).

Alt: „- keinen Container-Runtime-Socket des Hosts einbinden;“. Neu: „- keinen Container-Runtime-Socket des Hosts einbinden (der kurzlebige Hilfs-Container von `u-boot devcontainer up` gehört nicht zum erzeugten Ergebnis, siehe `LH-FA-DEV-010`);“ <!-- d-check:ignore (Kennung des Vorschlags) -->.

**Stelle 6 — Versionsfeld und Historie (§16).** Alt: `| Version | 0.17.0 |`. Neu: nächste freie Minor-Version. Neue Zeile: „`LH-FA-DEV-010`: neue Anforderung (Devcontainer ohne Editor starten und stoppen, Hilfs-Container mit Socket-Ausweisung, Zusage-Menge Linux-Host `amd64`); `LH-SA-CLI-002`: zwei Zeilen; `LH-FA-CLI-005`: `--quiet` für die neuen Befehle und die Socket-Ausweisung als Teil von Ausnahme 1; `LH-FA-CLI-007`: Wert `devcontainer`; `LH-FA-DEV-006`: Abgrenzung des Hilfs-Containers.“ <!-- d-check:ignore (Kennung des Vorschlags) -->

**Unverändert, mit Grund:** [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) und [`LH-FA-DEV-002`](../../../../spec/lastenheft.md#lh-fa-dev-002--vs-code-kompatibilität) (Dateien und Kompatibilität bleiben, der Befehl erzeugt nichts), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (die Klassen 10, 11 und 12 bestehen), [`LH-NFA-PORT-002`](../../../../spec/lastenheft.md#lh-nfa-port-002--keine-unnötigen-systemabhängigkeiten) und [`LH-NFA-PORT-004`](../../../../spec/lastenheft.md#lh-nfa-port-004--ausgeliefertes-container-image) (der Befehl läuft im Image nicht, das Image bleibt), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (die Auslegung steht im ADR, der Wortlaut trägt sie), [`LH-ABG-002`](../../../../spec/lastenheft.md#lh-abg-002--keine-ide-abhängigkeit) (der Befehl setzt sie um).

### Folge-Slices (Dateien in `open/`)

Vier Folge-Slices, Reihenfolge 1 → 2 → 3 → 4; der Gesamtzuschnitt passt in die Grenze von fünf, **keine Rückführung** nach §4 nötig. Priorität und Release-Zuordnung sind offen (Vorschlag des Planners: V2, nächstes Minor nach dem Befehl-Slice); jeder Folge-Slice trägt das als offene Frage.

| Nr. | Slice | Gegenstand |
|---|---|---|
| 1 | [`slice-cr-devcontainer-start-ohne-editor`](../open/slice-cr-devcontainer-start-ohne-editor.md) | Lastenheft-CR nach dem Wortlaut oben, nach Freigabe |
| 2 | [`slice-v1-devcontainer-start-fundament`](../open/slice-v1-devcontainer-start-fundament.md) | Port `DevcontainerLauncher`, Adapter, eingebetteter Kontext mit Pins, Docker-Tests S2, S5, S6 |
| 3 | [`slice-v1-cli-devcontainer-up-down`](../open/slice-v1-cli-devcontainer-up-down.md) | Befehl `devcontainer up`/`down`, Handbuch, `CHANGELOG.md`, Docker-Tests S1, S3, S4, S8, S9, Einzelmessung Colima |
| 4 | [`slice-v2-spec-devcontainer-start`](../open/slice-v2-spec-devcontainer-start.md) | Spezifikation und Sicht nach dem Code (Stellen aus [`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 10) |

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird `docs/plan/` (DOK, Greenfield: ADR, Slices) laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); die Schwelle aus zwei von drei Achsen ist dort erfüllt. Die Sub-Areas der Umsetzung (`internal/hexagon/port`, `internal/adapter/driven`, `internal/adapter/driving/cli`, `internal/e2e`, `tools/`) berührt dieser Slice nicht; sie stehen in den Folge-Slices.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md) (offen, zwei Belege; eine neue Zusage berührt den Vertrag), [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (verkörpert in `AGENTS.md` §3.10; die Aussagen zu Code und Dateien oben sind gelesen und mit Datei und Symbol genannt, was nicht gelesen ist — `engine.go`, der Inhalt von [`SPEC-042`](../../../../spec/spezifikation.md#spec-042--aufruf-der-compose-engine-und-preflight) — steht als nicht gelesen da) und [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (verkörpert; eine Spec- oder Lastenheft-Bedingung zu Engines und Pfaden braucht einen Test oder entfällt — daher die Docker-Tests mit Wächter im Code-Slice). Die Beobachtungen zu Messprotokollen und Host-Werkzeugen ([`BEO-WZG/host-werkzeug-ausserhalb-docker`](../observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md), [`BEO-DOK/messprotokoll-ausserhalb-des-repos`](../observations/BEO-DOK/messprotokoll-ausserhalb-des-repos/observation.md), [`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md)) betreffen diesen Slice nicht mehr, weil er nicht misst; sie bleiben für die Docker-Tests des Code-Slice relevant (Zugangsdaten nur als Umgebung des Befehls). Evidenz-/Diskrepanz-Risiko: mittel (Zusage-Menge gegen Belegbasis, §6).

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
