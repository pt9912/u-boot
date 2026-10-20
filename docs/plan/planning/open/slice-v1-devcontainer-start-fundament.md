# Slice v1: Devcontainer-Start — Fundament (Port, Adapter, eingebetteter Kontext des Hilfs-Images mit Pins, Docker-Tests zu Pfad, Socket und Lieferkette)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-devcontainer-start-ohne-editor`](../welle-devcontainer-start-ohne-editor.md). Zweiter von vier Folge-Slices von [`slice-v1-devcontainer-start-ohne-cli`](../in-progress/slice-v1-devcontainer-start-ohne-cli.md); nach [`slice-cr-devcontainer-start-ohne-editor`](../open/slice-cr-devcontainer-start-ohne-editor.md) und vor [`slice-v1-cli-devcontainer-up-down`](../open/slice-v1-cli-devcontainer-up-down.md). Release-Zuordnung offen (Vorschlag des Planners: V2, nächstes Minor nach dem Befehl-Slice; entscheidet der Projektinhaber). Der Slice liefert noch keinen Befehl und keine nutzersichtbare Änderung.

**Bezug:** [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (gepinnter externer Code), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zugangsdaten nur als Umgebung), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Socket nicht im erzeugten Ergebnis), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (Sentinels für 11 und 12), [`LH-NFA-PORT-002`](../../../../spec/lastenheft.md#lh-nfa-port-002--keine-unnötigen-systemabhängigkeiten); [`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) (muss `Accepted` sein; Entscheidungen 1, 2, 4, 5, 7), Muster [`ADR-0023`](../../adr/0023-podman-aus-podman-static-im-sandbox-devcontainer.md) und [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md).

**Berührte Spec-Stellen:** [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) (Label `devcontainer.local_folder`: der Spiegelpfad-Test S2 trägt, dass die Stelle für Container dieses Wegs wahr bleibt), [`SPEC-042`](../../../../spec/spezifikation.md#spec-042--aufruf-der-compose-engine-und-preflight) (gilt für Compose-Aufrufe; der Launcher hat eine eigene Vorab-Prüfung), `architecture.md` §2.4 und §2.6 (Port und Adapter). Der Slice **liest** sie und ändert keine; die Nachführung ist der Spezifikations-Slice. Der Verweis zeigt aufwärts.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Architect). **Datum:** 2026-10-20.

---

## 1. Ziel und Abgrenzung

**Ziel:** Das Hexagon kann einen Devcontainer über ein Hilfs-Image starten und stoppen: Port `DevcontainerLauncher`, Adapter im Paket `internal/adapter/driven/docker/`, eingebetteter Kontext (`Dockerfile`, `package.json`, `package-lock.json`) mit den Pins nach [`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 2, und Docker-Tests mit Wächter, die Pfad-Abbildung (S2), Socket-Haltung (S6), Lieferkette (S5) und die Linux-Hälfte von S7 für den Adapter tragen. Ohne Befehl: Der Adapter wird in Tests aufgerufen, nicht im Wiring (`cmd/uboot/main.go` bleibt unverändert).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Befehl, Application-Service, Ausgabe, Handbuch, `CHANGELOG.md`** — [`slice-v1-cli-devcontainer-up-down`](../open/slice-v1-cli-devcontainer-up-down.md); der Fundament-Slice ist ohne nutzersichtbare Änderung einzeln lieferbar und prüfbar.
- **Tests S1, S3, S4, S8, S9** — sie brauchen ein von `u-boot` erzeugtes Projekt und den Befehl; sie stehen im Befehl-Slice.
- **Spezifikation und Sicht** — [`slice-v2-spec-devcontainer-start`](../open/slice-v2-spec-devcontainer-start.md), nach dem Code.
- **Dauer-Scan des Hilfs-Images und Aufräumen alter Hilfs-Images** — akzeptierte Negative der Entscheidung ([`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Konsequenzen): lokal gebautes, nicht ausgeliefertes Image, Ereignis-Trigger statt Dauer-Scan.
- **Docker Desktop, Podman, rootless Docker, SELinux, `arm64` unter Linux** — nicht zugesagt; keine Probe dafür.

## 2. Definition of Done

Drei Liefer-Punkte; Gate-Läufe und Closure-Pflichten darunter zählen nicht mit.

- [ ] **Port, Adapter, Sentinel und eingebetteter Kontext:** `DevcontainerLauncher` (`Up`, `Stop`) und `ErrDevcontainerRuntime` in `port/driven`; der Adapter prüft vorab (Docker-Binary, Daemon, lokale Engine; zwei plus eine Prüfung, nicht `Engine.preflight`), baut das Hilfs-Image aus dem Tar-Kontext (nur wenn der Tag fehlt), startet es mit der Argumentliste aus [`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 4 (`--rm`, `--mount`, `-e <NAME>` ohne Wert, Spiegelpfad) und räumt bei Abbruch auf. Pins als Digest der Basis-Images, exakte CLI-Version und Lockdatei; Ziel `--ignore-scripts` trägt oder der Slice geht zurück an den Architect.
- [ ] **Docker-Tests mit Wächter** (`//go:build docker`, in `make test-docker`), je Zusage der Tabelle in [`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 9: S2 Pfad, S6 Socket, Umgebung und Lebensdauer, S5 Lieferkette, S7 Linux; dazu Unit-Tests ohne Docker (Argumentliste, Vorab-Prüfung, Ergebniszeile, Pins in den eingebetteten Dateien). Die Wächter stehen in §3.
- [ ] **Einmal-Messung** (Größe, Bauzeit kalt und warm, Dauer des ersten `up`, Trivy HIGH/CRITICAL) als Messprotokoll unter `docs/messungen/slice-v1-devcontainer-start-fundament/`, Skripte unter `tools/messung/devcontainer-start-fundament/`, nur `make` und Docker am Host, Rohdaten ohne Zugangsdaten; die Zahlen fließen als Einzelmessung in das Handbuch des Befehl-Slice, keine Zahl wird zugesagt.
- [ ] `make gates` grün, `make test-docker` grün (Exit-Code von `make` selbst); `make verify-depguard` nur, wenn eine Importregel berührt wird.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: ADR-Index unverändert; Roadmap-Zeile „In Arbeit“ bei der Closure; Handbuch und `CHANGELOG.md` unverändert (kein nutzersichtbarer Vertrag).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben: dritter Beleg der Pin-Alterung ([`BEO-APP/podman-static-pin-alterung`](../observations/BEO-APP/podman-static-pin-alterung/observation.md)), oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; der Implementer ergänzt sie in seinem Lauf und gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab (`AGENTS.md` §3.10). Eine Zeile „prüfen“ nennt eine Datei, die der Architect nicht gelesen hat.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/port/driven/devcontainer_launcher.go` | neu | Port `DevcontainerLauncher` (`Up(ctx, Optionen) (Ergebnis, error)`, `Stop(ctx, containerID) error`), Sentinel `ErrDevcontainerRuntime`. Vorbild `devcontainer_inspector.go` (gelesen: Sentinel als `var … = errors.New`, `DevcontainerContainer` als Ergebnistyp, Doku-Kommentar nennt Layer-Regeln und Implementierungsort). Wiederverwendet: `ErrDockerUnavailable` (`docker_engine.go`), `ErrContainerEngineNotLocal` (`devcontainer_inspector.go`) |
| `internal/adapter/driven/docker/devcontainer_launcher.go` | neu | Adapter: Vorab-Prüfung, Image-Bau per `docker build -` mit Tar-Strom, `docker run`-Argumente, SIGINT und `docker rm -f` bei Abbruch, Ergebniszeile (letzte nicht leere Zeile der Standardausgabe: `outcome`, `containerId` gegen `[0-9a-f]{12,64}`, `remoteUser`, `remoteWorkspaceFolder`), `Stop` per `docker stop`. Statischer Port-Pin `var _ driven.DevcontainerLauncher = (*Launcher)(nil)` wie `Inspector` und `Engine` (gelesen) |
| `internal/adapter/driven/docker/devcontainercli/Dockerfile`, `package.json`, `package-lock.json` | neu | eingebetteter Kontext (`go:embed`); `FROM node:22-bookworm@sha256:…`, `COPY --from=docker:29-cli@sha256:…` für `docker` und `docker-buildx` (Messweg `Dockerfile.cli`, gelesen), `npm ci`, exakte Version von `@devcontainers/cli` (Belegbasis 0.89.0); die Digests und die Lockdatei setzt der Slice |
| `internal/adapter/driven/docker/inspector.go` | update (wenn nötig) | gelesen: `requireLocalEngine` ist Methode von `Inspector` (liest `DOCKER_HOST`, sonst `docker context inspect`, nur `unix://` gilt als lokal); der Launcher braucht dieselbe Prüfung ohne Kopie, etwa als frei stehende Funktion mit unverändertem Verhalten und den vorhandenen Tests (`inspector_test.go`, nicht gelesen) |
| `internal/adapter/driven/docker/devcontainer_launcher_test.go` | neu | Unit: Argumentliste (`--mount` statt `-v`, `--rm`, Name, `-e <NAME>` ohne Wert, kein `--privileged`), Vorab-Prüfung mit gefälschtem Binary, Ergebniszeile (Erfolg, `outcome` ≠ `success`, Müll), Pins in den eingebetteten Dateien (Digest bei beiden `FROM`, `npm ci`, exakte Version) |
| `internal/adapter/driven/docker/devcontainer_launcher_docker_test.go` | neu | Docker-Tests (Build-Tag `docker`), Muster `engine_docker_test.go` (Datei vorhanden, nicht gelesen): die Wächter unten |
| `internal/adapter/driven/docker/docker_test.go`, `export_test.go` | prüfen | vorhandene Test-Hilfen (nicht gelesen); wiederverwenden statt neu schreiben |
| `tools/messung/devcontainer-start-fundament/`, `docs/messungen/slice-v1-devcontainer-start-fundament/` | neu | Einmal-Messung; Ablage nach [`harness/conventions.md`](../../../../harness/conventions.md#ablage-von-messprotokollen-und-skripten) (Skripte unter `tools/messung/`, Protokoll und Rohdaten unter `docs/messungen/`); Trivy-Aufruf wie `make image-scan` (Aufrufform in `repo.mk`, nicht gelesen) |
| `.golangci.yml` | **unverändert** | Abschnitt `depguard` gelesen (Regeln `domain-isoliert`, `application-no-adapter`, `application-no-net`, `port-no-application` …): ein Port unter `port/driven` und ein Adapter unter `adapter/driven` verletzen keine; löst `make lint` doch etwas aus, steht der Anlass hier, kein `nolint` |
| `internal/hexagon/port/driven/README.md` | **unverändert** | gelesen: nennt Schicht, Zweck und Verweise (Architektur, [`SPEC-042`](../../../../spec/spezifikation.md#spec-042--aufruf-der-compose-engine-und-preflight) für die Compose-Engine, [`SPEC-013`](../../../../spec/spezifikation.md#spec-013--import-regel-tabelle-der-schichten)); keine Port-Liste, die den Launcher nachführen müsste |
| `cmd/uboot/main.go`, `internal/adapter/driving/cli/` | **unverändert** | kein Wiring, kein Befehl in diesem Slice (`main.go` gelesen: der Adapter wird erst im Befehl-Slice konstruiert) |
| `Dockerfile` (Stufe `test-docker-tools`), `repo.mk` (`test-docker`), `.github/workflows/integration.yml` | **unverändert**, Risiko Laufzeit (§6) | gelesen: `test-docker-tools` bringt `docker-ce-cli` und das Compose-Plugin, kein buildx-Plugin; `test-docker` hängt Socket und `/tmp` unter demselben Pfad ein und läuft mit `-timeout 30m`; `integration.yml` hat `timeout-minutes: 30`. Läuft der Bau des Kontexts mit diesem Client nicht, ist das ein Befund an den Architect, keine stille Änderung der Fristen oder des Test-Images |
| `docs/user/benutzerhandbuch.md`, `CHANGELOG.md`, `spec/` | **unverändert** | kein nutzersichtbarer Vertrag; Handbuch und `CHANGELOG.md` im Befehl-Slice, `spec/` im Spezifikations-Slice |

**Docker-Tests mit Wächter** (jede Zusage hat eine Probe, die rot wird, wenn sie falsch wäre; Projekt in `t.TempDir()` unter dem geteilten `/tmp`, ein handgeschriebenes `devcontainer.json` ohne Features auf `mcr.microsoft.com/devcontainers/base:debian`):

- **S2 Pfad.** Nach `Up` trägt der Container das Label `devcontainer.local_folder` mit dem Pfad des Projekts, `Inspector.FindContainer` findet ihn. *Wächter:* derselbe Hilfs-Container mit dem Projekt unter `/w/…` (Messweg) liefert ein Label ≠ Pfad des Hosts, das der Test als Abweichung meldet.
- **S6 Socket im gestarteten Container.** `docker inspect` des Dev-Containers zeigt kein `docker.sock` in `Mounts`. *Wächter:* ein `devcontainer.json`, das den Socket selbst einhängt, wird vom selben Prüfer erkannt.
- **S6 Lebensdauer.** Nach `Up` und nach einem Kontextabbruch mitten im Lauf bleibt kein Container mit dem Namenspräfix des Hilfs-Containers. *Wächter:* ein Lauf ohne `--rm` (Test-Variante der Argumentliste) hinterlässt ihn und wird gefunden.
- **S6 Umgebung.** Eine gesetzte, in `remoteEnv` mit `${localEnv:<NAME>}` referenzierte Variable erreicht den Dev-Container, eine gesetzte, nicht referenzierte nicht; der Wert steht in keiner Argumentliste (Unit). *Wächter:* mit der Variante „alle Variablen durchreichen“ erreicht die nicht referenzierte Variable den Container, der Test wird rot. Zusätzlich ein Test auf die Allaussage „nie in Image, Volume oder erzeugter Datei“: Mit einem Platzhalter-Wert als Variable steht dieser Wert weder in `docker inspect` des Hilfs-Images und des gestarteten Dev-Containers (Labels, `Config.Env` des Images), noch in `docker inspect` des laufenden Hilfs-Containers (dort steht er in `Config.Env`, wenn `-e NAME=Wert` gesetzt wird; die Form `-e NAME` ohne Wert gibt der Test vor), noch in einer Datei unter dem Projektordner; *Wächter:* die Variante `-e NAME=Wert` macht den Test rot. Trägt der Test die Aussage nicht, entfällt die Allaussage in Wortlaut und Handbuch.
- **S6 `readonly`.** `Up` mit schreibgeschützt eingehängtem Projektordner; trägt es, bleibt es, sonst schreibend (Ausgang in §7, [`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 4).
- **S5 Lieferkette.** Bau aus dem eingebetteten Kontext nach `docker rmi` des Hilfs-Images (kein lokaler Cache, [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md)); `devcontainer --version` im Image entspricht dem Pin; die CLI startet mit `--ignore-scripts`. *Wächter:* eine verfälschte `integrity` im Lockfile bricht den Bau ab, das unveränderte Lockfile baut.
- **S7 Linux.** Vorab-Prüfung: fehlendes Docker-Binary → `ErrDockerUnavailable`, `DOCKER_HOST=tcp://…` → `ErrContainerEngineNotLocal` (Unit, ohne Engine); der Lauf des Adapters gegen die Engine des Runners ist die Zusage-Grundlage für Linux. Ein Test, der in der CI überspringt, trägt nichts ([`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md)).
- **Aufräumen.** Jeder Test entfernt, was er anlegt (Container, Volumes, Images ohne den Cache des Hilfs-Images), auch bei Fehlschlag ([`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md)).

## 4. Trigger

**Start** (`open` → `in-progress`, ggf. über `next`): `Verantwortlich:` gesetzt; [`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) ist `Accepted`; [`slice-cr-devcontainer-start-ohne-editor`](../open/slice-cr-devcontainer-start-ohne-editor.md) ist geschlossen; kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Bau des Kontexts braucht Änderungen am Test-Image oder an den Fristen von `make test-docker`/`integration.yml` (Gate-Änderung, `AGENTS.md` §3.6), oder die Docker-Tests brauchen mehr als einen Tag Laufzeit-Anpassung.
- `in-progress` → `open` (blockiert): ein Docker-Test widerlegt eine Zeile der Entscheidung — das Label trägt nicht den Pfad des Hosts (S2), der Dev-Container trägt den Socket (S6), `--ignore-scripts` trägt nicht (S5) —; der Slice geht mit dem Befund an den Architect (Ergänzung per Folge-ADR), nicht in eine stille Änderung der Zusage.

## 5. Closure-Trigger

DoD vollständig, Review-Report und Verification-Evidence liegen vor, jede Zeile der Zusage-Tabelle ([`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 9) für den Adapter hat einen grünen, nicht übersprungenen Test samt Wächter oder steht als „nicht zugesagt“ in §7, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Laufzeit von `make test-docker`:** das e2e-Paket braucht sieben bis elf Minuten (Kommentar in `repo.mk`, gelesen), das Limit ist 30 Minuten (`-timeout 30m`, `timeout-minutes: 30`); der Bau des Hilfs-Images und ein `devcontainer up` kommen dazu. — **Ausgang:** weiter offen; der Slice misst die zusätzliche Dauer und berichtet sie; reicht das Limit nicht, ist die Anhebung ein Plan-Eintrag, kein stilles Ändern.
- **Digest und Plattform:** Der Digest eines Index-Manifests gilt für `amd64` und `arm64`; ein Basis-Image ohne `arm64`-Variante bräche den Bau auf Colima. — **Ausgang:** weiter offen; der Bau auf `arm64` (Einzelmessung) und `amd64` (CI) prüft es.
- **Ein harter Abbruch von `u-boot` (SIGKILL) lässt den Hilfs-Container bis zum Ende der CLI laufen.** — **Ausgang:** weiter offen als akzeptiertes Negativ: `--rm` entfernt ihn am Ende des Laufs; die Lebensdauer ist die Dauer des `devcontainer up`, nicht unbegrenzt; kein Wächter für SIGKILL.
- **Pin-Alterung der drei Quellen** (CLI, Node-Image, Docker-CLI-Image). — **Ausgang:** weiter offen → `BEO-APP/podman-static-pin-alterung` (dritter Beleg bei der Closure); Trigger sind Ereignisse nach [`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 2.
- **`npm ci` im Hilfs-Image zieht Pakete mit Installationsskripten.** — **Ausgang:** weiter offen; entweder trägt `--ignore-scripts` (Test), oder der Slice geht zurück an den Architect.
- **Das Messprotokoll trägt Zugangsdaten.** — **Ausgang:** weiter offen; Rohdaten nur mit Platzhaltern, nie ein echtes Token ([`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md)).

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `internal/hexagon/port` (PORT), `internal/adapter/driven` (ADP, einschließlich der Docker-Tests im Adapter-Paket) und `tools/` (WZG), laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) alle Greenfield; die Schwelle aus zwei von drei Achsen ist dort erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-APP/podman-static-pin-alterung`](../observations/BEO-APP/podman-static-pin-alterung/observation.md) (drei gepinnte Quellen jetzt), [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md), [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md), [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md), [`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](../observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md) (die Tests laufen am eingebetteten Kontext, nicht am Messweg), [`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md), [`BEO-WZG/host-werkzeug-ausserhalb-docker`](../observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md) (Messskripte nur mit `make` und Docker). Evidenz-/Diskrepanz-Risiko: mittel (Pfad, Plattform, Laufzeit).

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
