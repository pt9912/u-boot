# Slice v1: Degradations-Fallbacks (Daemon, Egress, `vfs`) host-seitig ermitteln und ausweisen (`LH-FA-DEV-007`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — kein Release-Bezug festgelegt, keine Closure-Bedingung, die von der DoD dieses Slice verschieden wäre.
Eine künftige Release-Welle kann den Slice aufnehmen; dann ändert sich dieses Feld, nicht der Plan.
Herkunft: [`welle-cli-verbosity`](../done/welle-cli-verbosity.md), Slice
[`slice-v1-cli-quiet-nur-fehler`](../done/slice-v1-cli-quiet-nur-fehler.md) (zweiter Review F-4, Verifier-Auflage 3).

**Bezug:** [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Absatz Degradation: jeder Fallback wird in der Befehlsausgabe und in `u-boot doctor` ausgewiesen),
[`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Ausnahme 2),
[`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Wegfall der Egress-Restriktion als Fallback),
[`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Grenze: kein Host-Socket, keine Host-Dateien im Container),
[`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md), [`ADR-0013`](../../adr/0013-dokumentationsreferenzmodell.md),
[`ADR-0018`](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md), [`ADR-0019`](../../adr/0019-docker-daemon-aus-docker-apt-repo.md),
[`ADR-0022`](../../adr/0022-fallback-ausweisung-statusdateien-im-container.md) (`Accepted`, Entscheidung dieses Slice).
Das Lastenheft (0.13.0, Stand nach dem vorgelagerten CR) ändert dieser Slice **nicht**: Es ist Zielbild, der Code wird so gebaut, dass er die Zusage erfüllt
(Entscheidung des Projektinhabers 2026-10-17, „A: Code zum Ziel“). Die eine Stelle, an der die Messung das Ziel als unerreichbar zeigt (`/dev/fuse`), ändert der vorgelagerte
CR [`slice-cr-sandbox-fuse-startfehler`](../done/slice-cr-sandbox-fuse-startfehler.md) (Lastenheft 0.13.0), nicht dieser Slice.

**Berührte Spec-Stellen:**
`spezifikation.md` ([`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten),
[`LH-FA-DEV-008.a`](../../../../spec/spezifikation.md#lh-fa-dev-008a--verhalten-bei-nicht-gewährbarer-capability),
[`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) Ausnahme 2, §7 Historie);
`architecture.md` §2 (Driven-Ports, Adapter `docker/`; Abschnitt Fehlermodelle nur, falls der Port eine eigene Fehlerklasse bekommt). `lastenheft.md` bleibt unverändert.
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-17.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die drei Degradations-Fallbacks von [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Container ohne Docker-Daemon, Wegfall der Egress-Restriktion, `vfs`-Storage) sind
host-seitig ermittelbar und werden in `u-boot doctor` und, wo ein Befehl sie feststellt, in der Befehlsausgabe ausgewiesen — auch unter `--quiet`
([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) Ausnahme 2).

### Entscheidungen des Projektinhabers (2026-10-17, im Gespräch)

- **F1/F3–F7 und [`ADR-0022`](../../adr/0022-fallback-ausweisung-statusdateien-im-container.md):** `Accepted`. Statusdateien im Container (`/var/lib/u-boot/status/{vfs,egress,daemon}`), vom Host mit `docker cp` gelesen,
  neuer Driven-Port `DevcontainerInspector`, drei neue `doctor`-Items (`devcontainer.sandbox.fallback.vfs`, `….daemon`, `….egress`).
- **F2 „Befehlsausgabe“ — „doctor, plus generate“:** (i) `doctor` weist die Fallbacks aus; (ii) zusätzlich fragen `generate devcontainer` und `init` **nach dem Schreiben** denselben Port ab und warnen bei einem
  bestehenden Container mit `result=fallback` (Warnung mit `DegradationFallback`; best effort, still ohne Engine, ohne Container, ohne Datei oder bei ungültigem Dokument; unter `--quiet` sichtbar als Ausnahme 2;
  nicht bei `--dry-run`, weil nichts geschrieben wird — Festlegung der ADR, Entscheidung 10). `up` und ein neuer Befehl bekommen die Abfrage **nicht**.
- **F8 `/dev/fuse` — „A: Lastenheft-CR“:** Ein fehlendes `/dev/fuse` auf dem Container-Host wird wie `/dev/net/tun` als Startfehler geführt (Exit-Code `127` der Engine, kein stiller `vfs`-Fallback unter Docker);
  die Host-Prüfung in `checkSandboxRuntime` sagt dann nicht mehr „fällt auf `vfs` zurück“. Das legt der CR [`slice-cr-sandbox-fuse-startfehler`](../done/slice-cr-sandbox-fuse-startfehler.md) fest.

- **Zeitbudget von `make test-docker` (Entscheidung des Projektinhabers, 2026-10-18, zum Review-Befund F-1): „-timeout 30m in repo.mk“.** Das `go test`-Rezept des Ziels `test-docker` bekommt `-timeout 30m`; das e2e-Paket braucht gemessen sieben bis elf Minuten und lief im Review-Lauf in das Standardlimit von zehn Minuten. Nur dieses Target, keine Test- und keine Gate-Lockerung; das Limit des CI-Jobs (`integration.yml`, `timeout-minutes: 30`) passt dazu.

**Reihenfolge und Abhängigkeit: CR → dieser Slice.** Der CR setzt den Wortlaut der Host-Prüfung in `checkSandboxRuntime` und die Zeilen in
[`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten) (Tabelle) und
[`LH-FA-DEV-007.b`](../../../../spec/spezifikation.md#lh-fa-dev-007b--ausstattung-bei-nestedruntime-podman); dieser Slice zieht darüber hinaus die Aussage („Jeder Fallback wird …“, Prüfsatz),
[`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) Ausnahme 2 und die Item-Beschreibung nach und ändert die vom CR gesetzten Zeilen nicht erneut.
Der `vfs`-Fallback bleibt nach dem CR ein Fallback mit enger Ursache (Container ohne `--device=/dev/fuse` gestartet, etwa nach geänderten `runArgs`); sein `doctor`-Item `devcontainer.sandbox.fallback.vfs` und die Statusdatei `vfs` bleiben im Umfang.
Der Slice startet erst, wenn der CR in `done/` liegt (§4).

**Übernimmt:** [`slice-cr-sandbox-fallback-ausweisung`](../done/slice-cr-sandbox-fallback-ausweisung.md) (Gegenstand entfallen als CR; hier der Code-Zweig, den
[`slice-v1-sandbox-fallback-ausweisung`](../done/slice-v1-sandbox-fallback-ausweisung.md) als Lösung A beschrieb).

### Gefundener Stand im Code (Ist-Stand, gelesen am 2026-10-17)

| Fallback | Wo er entsteht | Was heute host-seitig ausgewiesen wird | Fundstelle |
|---|---|---|---|
| `vfs`-Storage | Im Container: `[ ! -c /dev/fuse ]`, schreibt `storage.conf` mit `driver = "vfs"`; Zeile auf Standardfehlerausgabe, bei `fail` Exit `11` (`postCreateCommand`, als `sh` mit dem Container-Benutzer) | `doctor`-Item `devcontainer.sandbox.runtime` mit `DegradationFallback`, **nur** auf einem Linux-Host (Prüfung `/proc/sys/kernel/osrelease`) und nur bei Policy `warn`; auf macOS/Colima meldet das Item OK und verweist auf das Startscript | `internal/hexagon/application/doctor_sandbox.go`, `checkSandboxRuntime`; `templates/devcontainer/sandbox-init.sh.tmpl`, Funktion `unavailable`; Einhängung `composePostCreate` in `devcontainer_sandbox.go` |
| Container ohne Docker-Daemon | `docker-init.sh`: `dockerd` startet nicht oder ist nach 60 s nicht bereit; Zeile auf Standardfehlerausgabe, Exit `0` bei `warn`, `11` bei `fail` (`postStartCommand`, per `sudo sh`) | nichts; `checkSandboxRuntime` meldet für `nestedRuntime: docker` OK mit dem Hinweis, das Startscript prüfe die Bereitschaft | `templates/devcontainer/docker-init.sh.tmpl`, Funktion `unavailable`; `doctor_sandbox.go`, `checkSandboxRuntime` (Zweig `nestedDocker`) |
| Wegfall der Egress-Restriktion | `egress-init.sh`: `nft list tables` scheitert ohne `NET_ADMIN`; Zeile auf Standardfehlerausgabe, Exit `0` bei `warn`, `11` bei `fail` (`postStartCommand`, per `sudo sh`, vor dem Daemon-Start) | nichts; `checkSandboxEgress` prüft nur die Konfiguration und nennt, die Capability werde beim Containerstart geprüft | `templates/devcontainer/egress-init.sh.tmpl`; `doctor_sandbox.go`, `checkSandboxEgress`; Reihenfolge in `devcontainerTemplateData` (`devcontainer_sandbox.go`) |

Weitere Befunde, an der Stelle gelesen:

- **u-boot kennt den Devcontainer-Container nicht.** `devcontainer.json` ist im Build-Modus (`build.dockerfile`) mit `runArgs` erzeugt (`templates/devcontainer/devcontainer.json.tmpl`); gestartet wird er vom
  Devcontainer-Werkzeug (Editor oder CLI), nicht von u-boot. `DockerEngine` kennt nur Compose-Projekte (`internal/hexagon/port/driven/docker_engine.go`), `DockerProbe` nur `Version`, `Info`, `ComposeVersion`
  (`docker.go`). Weder Containername noch -ID sind u-boot bekannt; der Workspace-Volume-Name trägt `${devcontainerId}`, das das Werkzeug bildet.
- Für `init`, `generate` und `up` erzeugt keine produktive Quelle einen Eintrag mit `DegradationFallback` (Kommentar an `sandboxRelaxationWarnings`, `devcontainer_sandbox.go`). `generate` kennt nur die Konfiguration,
  `up` bedient nur Compose-Projekte — die Fallbacks entstehen erst beim Start des Devcontainers, den keiner dieser Befehle auslöst. Als „Befehl, der sie feststellt“ bleibt damit nach heutigem Stand nur `doctor` (siehe Frage F2).
- Die Kennzeichnung samt Filter unter `--quiet` ist vorhanden (`internal/adapter/driving/cli/warnings.go`, `visibleWarnings`; `internal/adapter/driving/cli/doctor.go`, `writeDoctorReport`) und über Fake-Einträge gepinnt
  (`internal/adapter/driving/cli/quiet_test.go`, `quietFallbacks`); ein neues `doctor`-Item mit `DegradationFallback: true` ist unter `--quiet` ohne CLI-Änderung sichtbar.
- Die Spezifikation formuliert die Zusage enger als das Lastenheft ([`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) Ausnahme 2 „soweit sie vom Host aus ermittelbar sind“,
  [`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten) „(soweit vom Host aus ermittelbar)“,
  [`LH-FA-DEV-008.a`](../../../../spec/spezifikation.md#lh-fa-dev-008a--verhalten-bei-nicht-gewährbarer-capability) „die Capability wird beim Containerstart geprüft“). Mit der Entscheidung „Lastenheft ist Zielbild“ ist die Spezifikation
  dem Lastenheft anzugleichen, nicht umgekehrt (Nachzug in diesem Slice).
- Die Startscripte laufen unter dem Devcontainer-Werkzeug (`postCreateCommand`/`postStartCommand`); die Skripte selbst schreiben nur auf die Standardfehlerausgabe ihres Prozesses.

### Was der Host weiß — und was erst zur Laufzeit

| Zustand | Ohne laufenden (oder je gelaufenen) Container | Erst zur Laufzeit / aus dem Start |
|---|---|---|
| `vfs`: `/dev/fuse` fehlt | **Nur auf einem Linux-Host**, dessen Engine auf dem Host selbst läuft: `/dev/fuse` lesbar (heutige Prüfung, bleibt). Auf macOS/Colima und anderen VM-Engines: nicht ermittelbar. | Das Startscript prüft `[ ! -c /dev/fuse ]` im Container. |
| Daemon nicht bereit | Allenfalls eine **Vorhersage** (`--privileged` wird nicht gewährt, etwa bei einer rootless Engine; Annahme, über `docker info` zu messen) — keine Feststellung. | `docker-init.sh`: Daemon-Start und Bereitschaft nach bis zu 60 s. |
| Egress-Capability fehlt | Nichts verlässlich (so schon [`LH-FA-DEV-008.a`](../../../../spec/spezifikation.md#lh-fa-dev-008a--verhalten-bei-nicht-gewährbarer-capability): vom Host aus nicht bestimmbar); `doctor` prüft die Konfiguration. | `egress-init.sh`: `nft list tables`. |

Folge: Für Daemon- und Egress-Fallback (und für `vfs` außerhalb eines Linux-Hosts) braucht die Zusage einen **Kanal vom Start im Container zum Host** oder ein Lesen des Container-Zustands durch `doctor`.
Genau das ist die Entwurfsfrage dieses Slice; Entscheidung beim Architect und beim Projektinhaber (siehe unten), nicht beim Planner.

### Entwurfsfrage: Wie kommt der Fallback zum Host? (Alternativen, keine Entscheidung)

Randbedingung aus [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil): kein Host-Socket im Container, keine Host-Dateien im Container, Host-Arbeitsverzeichnis nicht eingebunden. Ein Kanal darf deshalb **nur den Host lesen lassen**
(Host → Container-Zustand über die Engine), nicht den Container zum Host schreiben lassen. Die Annahmen, die noch **gemessen** werden müssen, sind als *Annahme* markiert; das Messprotokoll gehört in die Folge-ADR (Muster:
[`ADR-0019`](../../adr/0019-docker-daemon-aus-docker-apt-repo.md) gegenüber [`ADR-0018`](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md): Messung auf Colima, Entscheidung, verglichene Alternativen).

| Alt. | Mechanismus | Vorteile | Nachteile / Risiken | Sicherheit | Aufwand | Messbarkeit Colima / Linux |
|---|---|---|---|---|---|---|
| **(a)** | `doctor` liest die Startausgabe über die Engine (`docker logs`, `docker inspect`, Exit-Code des `postStartCommand`) | Keine Änderung der Startscripte; kein neuer Dateipfad im Container | *Annahme:* `docker logs` zeigt nur die Ausgabe des Hauptprozesses (PID 1), nicht die von `postCreateCommand`/`postStartCommand` — diese laufen als `exec` des Devcontainer-Werkzeugs, und deren Ausgabe liegt im Log des Werkzeugs; ein Exit-Code eines `exec` steht nicht am Container. Würde (a) tragen, müssten die Skripte nach `/proc/1/fd/2` schreiben (das `postCreateCommand` läuft ohne `sudo`, PID 1 gehört vermutlich root). Log-Rotation und Neustart-Anhängen machen den Stand unscharf. Braucht Container-Identifikation und einen laufenden Container | Lesen nur über die Engine; Log-Inhalt ist vom Container geschrieben (nicht vertrauenswürdig) | mittel (Port + Adapter + Identifikation) | Messbar auf beiden Hosts per `docker logs` nach einem Lauf der Skripte (die E2E-Tests in `internal/e2e/sandbox_devcontainer_docker_test.go` starten die Skripte per `docker run`/`exec`, nicht per Devcontainer-Werkzeug — ein Unterschied, der in der ADR zu messen ist) |
| **(b)** | Die Startscripte schreiben ein **Statusdokument im Container** (etwa `/run/u-boot/…` oder unter dem Benutzer-Home); `doctor` liest es über die Engine (`docker cp <container>:<pfad> -` auch bei gestopptem Container, oder `docker exec … cat` bei laufendem). Variante (b2): Statusdatei auf dem Workspace-Volume, Lesen per `docker run --rm -v <volume>:… cat` | Ein stabiler, versionierbarer Vertrag (Schlüssel, Stand, Zeitstempel); überlebt ein Stoppen; unabhängig vom Log des Werkzeugs; `docker cp` ist lesend. Host-seitig gleich auf Colima (CLI über Socket) und Linux | Braucht Container-Identifikation (Frage F3) und einen neuen Driven-Port mit `docker ps`/`inspect`/`cp`; Stand kann veraltet sein (`postCreateCommand` läuft nur einmal je Image-Bau, `postStartCommand` je Start) → Zeitstempel und Container-ID im Dokument. (b2) braucht ein Image auf dem Host und startet einen Hilfs-Container; kein Host-Bind-Mount, aber ein neuer Lauf auf dem Host | Das Statusdokument ist vom (womöglich kompromittierten) Agenten schreibbar: **nicht vertrauenswürdig**, nur Hinweis, keine Sicherheitsaussage (das Profil ist Schadensbegrenzung, keine Grenze); strenges Parsen (geschlossene Schlüsselmenge, Größenlimit, keine Steuerzeichen in der `doctor`-Ausgabe). Kein neuer Weg Container → Host | mittel bis hoch (Skripte, Statusformat, Port, Adapter, Tests) | Messbar auf beiden: `make test-docker` startet die Skripte schon; Erweiterung um Lesen des Dokuments. Linux-CI braucht zum Auslösen von `vfs` das Weglassen von `--device=/dev/fuse` (so schon `TestE2E_SandboxDevcontainer_FailPolicyWithoutFuse`) |
| **(c)** | Das Startscript setzt ein Container-Label oder löst ein Docker-Event aus; `doctor` liest `docker inspect` | Wäre der schmalste Host-Zugriff (nur `inspect`) | *Annahme (nach Docker-Semantik, vom Architect zu bestätigen):* Labels sind nach `create` unveränderlich, ein Prozess im Container kann kein Docker-Event erzeugen, ohne den Daemon-Socket zu erreichen — den das Sandbox-Profil gerade nicht einbindet. **Vermutlich nicht umsetzbar.** Variante (c′): ein `HEALTHCHECK`/`--health-cmd`, der das Statusdokument ausgibt; `docker inspect` zeigt `State.Health` samt Ausgabe des letzten Checks | (c′): Ausgabe vom Container geschrieben (nicht vertrauenswürdig); Healthcheck ist keine Sicherheitslockerung | (c′): mittel; ändert `runArgs`/Dockerfile und damit die erzeugte Ausgabe bestehender Projekte nach dem nächsten `generate` (wie beim sechsten `runArgs`-Eintrag in [`ADR-0021`](../../adr/0021-podman-build-im-sandbox-devcontainer.md)); der Container erscheint unter Umständen als `unhealthy` im Editor und in `docker ps` | (c′): `docker inspect` ist auf beiden Hosts gleich; Ausgabegröße des Health-Logs begrenzt (zu messen) |
| **(d)** | Mischform: (b) als Hauptkanal, die heutige host-seitige `/dev/fuse`-Prüfung bleibt (Linux-Host, ohne laufenden Container), optional (c′) als Fallback, wenn die Container-Identifikation scheitert | Deckt jeden Zustand ab, den der Host ohne Container kennt, und jeden, den der Start erzeugt; „nicht ermittelbar“ wird ehrlich als Info ausgewiesen statt verschwiegen | Zwei Quellen, eine Vorrangregel (laufendes Dokument vor Host-Vorhersage); größerer Umfang → Slice-Zerlegung möglich (siehe §4) | wie (b) | hoch | wie (b) plus die bestehende Linux-Prüfung |

**Einordnung des Planers (keine Entscheidung):**

- (a) hängt an einer ungemessenen Annahme über den Ort der Skript-Ausgabe; ohne Messung ist sie die unsicherste Option, und der Exit-Code des `postStartCommand` ist nach meinem Wissen vom Container aus nicht lesbar.
- (c) in der Reinform scheitert vermutlich an Docker-Semantik (zu bestätigen); (c′) ist denkbar, ändert aber sichtbar das erzeugte Ergebnis und den Container-Zustand im Editor.
- (b) mit `docker cp`/`exec` hält die Grenzen von [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (kein Pfad Container → Host, keine neue `runArgs`), liefert einen versionierbaren Vertrag und trägt auch einen gestoppten Container;
  sein Preis ist die Container-Identifikation und die Vertrauensfrage zum Dokumentinhalt. Innerhalb des Devcontainers kann `doctor` das Dokument direkt über `FileSystem` lesen (`RuntimeEnvironment.InContainer()`, heute überspringt `doctor` dort die Host-Prüfungen) — ohne neuen Port.
- (d) ist (b) plus das, was der Host ohnehin weiß; wenn (b) gewählt wird, ist (d) der natürliche Schnitt.
- Auf macOS/Colima ist der Docker-CLI-Zugriff host-seitig derselbe wie auf Linux (Socket), anders als die `/dev/fuse`-Prüfung, die dort nichts sieht: Ein engine-vermittelter Kanal löst genau die Lücke der Hosts, die der Slice braucht.

**Fragen an den Architect und den Projektinhaber** (die Folge-ADR beantwortet sie, der Implementer baut erst danach):

- **F1 Mechanismus:** (a), (b), (c′), (d) oder anders? Messung der Annahmen (Ort der Skript-Ausgabe, Lesbarkeit per `docker cp`, Health-Log) auf Colima und Linux gehört zur ADR.
- **F2 „Befehlsausgabe“:** Das Lastenheft sagt „in der Befehlsausgabe und in `u-boot doctor`“. Nach Ist-Stand stellt nur `doctor` einen Fallback fest. Genügt die Ausgabe von `doctor` als „Befehlsausgabe“,
  oder soll ein weiterer Befehl (etwa `generate devcontainer` bei laufendem Container, `up`) den Stand ausweisen? Ein neuer Befehl oder eine neue Flag-Fläche wäre eine Lastenheft-Änderung und ist hier ausgeschlossen.
  **Entschieden (2026-10-17): „doctor, plus generate“** — siehe oben.
- **F3 Container-Identifikation:** Wie findet `doctor` den Container? Label des Devcontainer-Werkzeugs (Annahme: `devcontainer.local_folder`, zu messen; gilt je nach Werkzeug — Editor, CLI, Codespaces — möglicherweise nicht gleich), ein von u-boot gesetztes Label in `runArgs`
  (`--label u-boot.project=<name>`, ändert die erzeugte Ausgabe), oder der Containername? Ohne Treffer: „nicht ermittelbar“ als Info, nie als Fallback.
- **F4 Vertrauensgrenze:** Das Dokument stammt aus einem Container, in dem ein Agent läuft. Format, Parser und Ausgabe-Bereinigung; und die ausdrückliche Aussage, dass „kein Fallback gemeldet“ nicht „kein Fallback“ beweist.
- **F5 Port-Schnitt:** `DockerProbe` erweitern (heute read-only und schmal, [`ADR-0002`](../../adr/0002-hexagonale-architektur.md) / Sicht-Spec) oder ein eigener Driven-Port (etwa für Container-Zustand)? Welche Engine-Kommandos sind erlaubt (`ps`, `inspect`, `cp`; nicht `exec`)?
- **F6 Stand-Semantik:** Wie behandelt `doctor` ein Dokument eines früheren Starts (Zeitstempel, Container-ID, gestoppter Container)? Gilt „nie gestartet“ als OK-Info?
- **F7 Strenge:** Unter `onUnavailable: fail` endet das Startscript mit `11`; `doctor` zeigt dann einen Fehler (kein `DegradationFallback`), unter `warn` eine Warnung mit `DegradationFallback`. Und `doctor --strict` endet bei der Warnung mit `11` (zu bestätigen an der Stelle) — gewollt?
- **F8 `--device=/dev/fuse`-Befund:** `podmanRelaxations` setzt `--device=/dev/fuse` (`devcontainer_sandbox.go`); fehlt das Device auf dem Container-Host, lehnt die Engine den Start womöglich ab (wie bei `/dev/net/tun`, [`ADR-0021`](../../adr/0021-podman-build-im-sandbox-devcontainer.md)). Dann ist der `vfs`-Fallback unter Docker auf Hosts ohne `/dev/fuse` nicht erreichbar,
  und die Zusage hat dort nichts auszuweisen. Ist das zu messen und, falls ja, Gegenstand dieses Slice oder eines eigenen?
  **Gemessen (M13) und entschieden (2026-10-17): „A: Lastenheft-CR“** — eigener Vorgang, vorgelagert (siehe oben).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Änderung des Lastenhefts** — Entscheidung des Projektinhabers: das Lastenheft ist Zielbild. Der Slice liefert Code, Tests, Spezifikations-Nachzug, Handbuch, CHANGELOG. Einzige Ausnahme `/dev/fuse`: der vorgelagerte CR [`slice-cr-sandbox-fuse-startfehler`](../done/slice-cr-sandbox-fuse-startfehler.md) (Lastenheft 0.13.0), nicht dieser Slice.
- **Ein neuer Befehl oder eine neue Option** (etwa `u-boot devcontainer status`) — wäre eine neue Zusage am Vertrags-Stratum (siehe F2); eigener Vorgang mit Change Request.
- **Ein Kanal Container → Host mit Schreibrichtung** (Socket, Bind-Mount, Rückmeldung an den Host-Prozess) — verletzt [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil); nur das Lesen durch den Host über die Engine ist im Rahmen.
- **Änderung an `--quiet` oder an der Kennzeichnung `DegradationFallback`** — vorhanden und gepinnt; neue Items nutzen sie.
- **Neue Fallback-Klassen oder Sicherheitslockerungen** — [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) zählt drei Fallbacks; das fehlende `/dev/net/tun` bleibt ausdrücklich keine Degradationsstufe.
- **Podman als Host-Engine** — für `nestedRuntime` nicht zugesagt ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)); der Kanal wird auf Docker (Linux, Colima) gemessen.

## 2. Definition of Done

- [ ] **Ermittlung und Ausweisung (Code):** Daemon-, Egress- und `vfs`-Fallback werden gemäß der Folge-ADR host-seitig ermittelt und in `u-boot doctor` (je Fallback ein Item mit `DegradationFallback`, unter `--quiet` sichtbar) und in der Warnungsausgabe von `generate devcontainer` und `init` (nach dem Schreiben, bei einem bestehenden Container mit `result=fallback`; best effort, still sonst, nicht bei `--dry-run`, unter `--quiet` sichtbar) ausgewiesen; Test je Fallback (Happy, Boundary „nicht ermittelbar“, Negative: fehlerhaftes Dokument) und je Befehl `generate devcontainer`/`init` (Warnung bei `fallback`, still bei `ok`/`failed`/ungültig/Engine-Fehler, keine Abfrage bei `--dry-run`); Replay-/Golden-Evidence für die geänderten Startscripte nach [`harness/replay.md`](../../../../harness/replay.md)
      ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) Ausnahme 2, [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)).
- [ ] **Spezifikation und Sicht nachgezogen:** [`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten), [`LH-FA-DEV-008.a`](../../../../spec/spezifikation.md#lh-fa-dev-008a--verhalten-bei-nicht-gewährbarer-capability) und
      [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) Ausnahme 2 tragen die Aussage des Lastenhefts ohne die Einschränkung „soweit vom Host aus ermittelbar“ (statt dessen das gewählte Verfahren und die Grenze „nicht ermittelbar“); Statusformat als Spezifikations-Stelle;
      Port und Adapter in `spec/architecture.md` (ohne ADR-, Slice- oder Wellen-Verweise, [`AGENTS.md`](../../../../AGENTS.md) §3.4); Historie-Zeile in `spezifikation.md` §7.
- [ ] Doku-Update: `docs/user/benutzerhandbuch.md` und CHANGELOG nennen den geltenden Stand der Ausweisung (öffentlicher Vertrag).
- [ ] `make gates` grün; `make test-docker` auf dem Entwicklungs-Host (Colima) grün oder mit benanntem Grund nicht gelaufen; `make doc-trace` ohne neue Waise.
- [x] Folge-ADR (Architect) ist `Accepted` und im ADR-Index eingetragen, bevor der Implementer beginnt.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne Wellen-Betrieb hier geprüft.

## 3. Plan (vor Code)

Rollenreihenfolge: **Planner** (dieser Plan) → **Architect** (Folge-ADR zu F1–F8, Messung auf Colima/Linux; Modul 8: Entscheidungen prüfen und vorschlagen) und **Projektinhaber** (F1, F2, F7, F8) → **Implementer** → **Reviewer** → **Verifier**.
Die Tabelle nennt jede Datei, die der Diff voraussichtlich berührt (§3.10 von [`AGENTS.md`](../../../../AGENTS.md)). Der Implementer hat sie nach der letzten Änderung gegen `git diff --stat` abgeglichen; Zeilen „unverändert“ und die zwei Zeilen „in der Planung nicht aufgeführt“ halten das Ergebnis des Abgleichs fest.

| Datei / Komponente | Änderungs-Art | Begründung | abhängig von |
|---|---|---|---|
| `docs/plan/adr/` (neue ADR, nächste freie Nummer) und `docs/plan/adr/README.md` (Index) | update (erledigt: `Accepted` 2026-10-18, Index nachgezogen) | Folge-ADR: Mechanismus, Messung, verglichene Alternativen; Schärft [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer); Muster [`ADR-0019`](../../adr/0019-docker-daemon-aus-docker-apt-repo.md) | Architect, Projektinhaber |
| `internal/hexagon/application/templates/devcontainer/sandbox-init.sh.tmpl`, `docker-init.sh.tmpl`, `egress-init.sh.tmpl` | update | je Skript die Funktion `write_status` (Temp-Datei, `mv -f`, Fehler ändern den Exit-Code nicht) und ein Aufruf in jedem Zweig (`ok`, `fallback`, `failed` vor `exit 11`); Exit-Codes und Stderr-Zeilen bleiben | — |
| `internal/hexagon/application/templates/devcontainer/Dockerfile.tmpl` | update | `RUN mkdir -p /var/lib/u-boot/status && chown vscode:vscode …`, nur bei `NestedPodman`, `NestedDocker` oder `Egress` | — |
| `internal/hexagon/application/devcontainer_sandbox.go` | update (nur Kommentar) | Kommentar an `sandboxRelaxationWarnings` nennt die Ausweisung der Fallbacks aus den Statusdateien; kein Label, keine `runArgs`, `composePostCreate` und die Startkommandos bleiben | — |
| `internal/hexagon/domain/sandbox_status.go`, `sandbox_status_test.go` | neu | Wertobjekt `SandboxStatus` mit `ParseSandboxStatus` (geschlossene Schlüssel- und Wertemenge, 512 Byte, strikter UTC-Zeitstempel) und `ErrInvalidSandboxStatus`; I/O-frei | — |
| `internal/hexagon/port/driven/devcontainer_inspector.go` | neu | Driven-Port `DevcontainerInspector` (`FindContainer`, `ReadFile`) mit `ErrContainerFileNotFound`, `ErrContainerFileInvalid`, `ErrContainerEngineNotLocal`; `port/driven/README.md` unverändert | — |
| `internal/adapter/driven/docker/inspector.go`, `inspector_test.go`, `export_test.go` | neu / update | `Inspector` neben `Probe` und `Engine` (Port-Pin `var _ driven.DevcontainerInspector`): `docker ps`, `docker cp … -` mit Tar-Kopf-Prüfung, `docker context inspect` für die Engine-Prüfung; einziger Ort, an dem dieser `os/exec docker …` läuft; Fake-Runner-Tests | — |
| `internal/hexagon/application/devcontainer_fallback.go` | neu | in der Planung nicht aufgeführt: die drei `fallbackProbe` je Fähigkeit, `fallbackReader` (Host über den Port, im Devcontainer über `FileSystem`), Abbildung auf `domain.Diagnostic` und auf Warnungen | — |
| `internal/hexagon/application/doctor.go` | update | Feld `inspector`, `WithDevcontainerInspector` (statt neuem Konstruktor-Parameter, damit die Aufrufstellen unverändert bleiben), `Check` hängt `checkSandboxFallbacks` an, `DoctorCheckIDs` um die drei Kennungen | — |
| `internal/hexagon/application/doctor_sandbox.go` | update | `checkSandboxFallbacks`; die Texte von `checkSandboxRuntime` (Zweig `nestedDocker`) und `checkSandboxEgress` verweisen auf die neuen Items; Kommentare an den Kennungen berichtigt | — |
| `internal/hexagon/application/generate.go` | update | Feld `inspector`, `WithDevcontainerInspector`, `existingContainerFallbackWarnings` nach `writeDevcontainerPhase`: nicht bei `PreviewDryRun`, nicht ohne konfigurierte Fähigkeit | — |
| `internal/hexagon/application/initproject.go` | update | Feld `inspector`, Option `WithDevcontainerInspector`, `existingContainerFallbackWarnings` nach `executeUBootYAML` (liest die geschriebene `u-boot.yaml`) | — |
| `internal/hexagon/application/noop.go` | update | `noopDevcontainerInspector` (findet keinen Container) | — |
| `cmd/uboot/main.go` | update | `docker.NewInspector()` an `doctorSvc`, `generateSvc` und die Option des Init-Service | — |
| `internal/hexagon/application/devcontainer_fallback_test.go` | neu | Fake-Inspector; `doctor` je Fähigkeit und Zustand, nur konfigurierte Fähigkeiten, im Devcontainer über `FileSystem`; `generate devcontainer` und `init`: Warnung bei `fallback`, still sonst, kein Aufruf bei `--dry-run` und ohne Fähigkeit | — |
| `internal/hexagon/application/doctor_sandbox_test.go` | update | zwei Pins nachgezogen (Egress-Text, Kommentar zum Kennzeichen); `fakes_test.go` und `doctor_test.go` bleiben | — |
| `internal/hexagon/application/generate_sandbox_test.go` | update | Replay-/Golden-Assertions der Startscripte und des Dockerfile ([`harness/replay.md`](../../../../harness/replay.md)) | — |
| `harness/replay.md` | update | zwei Zeilen: Startscripte mit Statusdateien, Abfrage in `generate`/`init` | — |
| `internal/adapter/driving/cli/doctor_fallback_test.go` | neu | `doctor` je Zustand: `--quiet`, `--strict`, `--json`, Volltext; `doctor.go`, `warnings.go` und `quiet_test.go` bleiben | — |
| `internal/adapter/driving/cli/jsontestutil/coderegistry.go` | update | in der Planung nicht aufgeführt: die drei Kennungen in `DefaultAllowedCodes` (Drift-Test gegen die Registry der Spezifikation) | — |
| `internal/e2e/sandbox_status_docker_test.go` (Build-Tag `docker`) | neu | echte Engine: Container mit Label `devcontainer.local_folder`, Startscripte per `docker exec`, Lesen mit Inspector und `doctor`, je Fähigkeit `ok`/`fallback`/`failed`; Inspector gegen Symlink, FIFO, Verzeichnis, Übergröße, gestoppten und mehrere Container, entfernte Engine; `sandbox_devcontainer_docker_test.go` bleibt | — |
| `spec/spezifikation.md` | update | [`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten), [`LH-FA-DEV-008.a`](../../../../spec/spezifikation.md#lh-fa-dev-008a--verhalten-bei-nicht-gewährbarer-capability), [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) Ausnahme 2, neu [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) (Statusformat, Ermittlung, Abbildung), drei Zeilen in der Registry von [`SPEC-017`](../../../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry), §7 Historie | — |
| `spec/architecture.md` | update | Port `DevcontainerInspector` (Abschnitt Driven-Ports), Adapter im Paket `docker/`, Zeile `ARC-008`; kein neuer Fehlerklassen-Abschnitt (die Sentinels sind keine Exit-Code-Fälle) | — |
| `docs/user/benutzerhandbuch.md` | update | Abschnitt „Fallbacks sehen“, `--quiet`-Zeile, doctor-Abschnitt, Daemon- und Egress-Absatz, Fehlerbehebung; die `/dev/fuse`-Zeilen setzt der vorgelagerte CR | — |
| `docs/maintainer/cli-json-contract.md` | unverändert | §2.2 gilt unverändert; neue Kennungen erscheinen in `diagnostics[]` nur als `warn`/`error` | — |
| `CHANGELOG.md` | update | Eintrag unter `[Unreleased]` (Added) | — |
| `README.md`, `README.de.md` | unverändert | nennen „vier Ausweisungen“, keine Einzelheiten | — |
| `docs/plan/planning/open/README.md` | unverändert | die Zeile dieses Slice steht in `next/README.md` und zeigt auf `in-progress/` | — |
| `docs/plan/planning/in-progress/roadmap.md` | update | Ruhe-Marker „Nichts in Arbeit.“ entfernt, solange dieser Slice in `in-progress/` liegt (d-check `planning`) | — |
| `repo.mk` | update | in der Planung nicht aufgeführt: `-timeout 30m` im `go test`-Rezept von `test-docker` (Entscheidung des Projektinhabers 2026-10-18, §1); Zeitbudget, keine Prüfung gelockert | — |
| `docs/reviews/2026-10-18-slice-v1-sandbox-fallback-host-ausweisung.md` | neu | Review-Report (Reviewer); die Befunde F-1 bis F-6 sind behoben, F-7 bis F-11 vermerkt | — |
| `docs/plan/planning/in-progress/slice-v1-sandbox-fallback-host-ausweisung.md` | update | diese Tabelle gegen `git diff --stat` abgeglichen, DoD-Zeile „Folge-ADR `Accepted`“ abgehakt, Lastenheft-Stand 0.13.0 (nach dem vorgelagerten CR) | — |
| `docs/plan/planning/in-progress/carveouts.md`, `.golangci.yml` (depguard), `.d-check.yml` | unverändert | kein neuer Carveout, kein neues Paket | — |
| `spec/lastenheft.md` | **unverändert** | Entscheidung 2026-10-17 | — |

### Architect: Entwurf und Messprotokoll (2026-10-17)

Folge-ADR: [`ADR-0022`](../../adr/0022-fallback-ausweisung-statusdateien-im-container.md) (`Accepted` am 2026-10-18; Messprotokoll M1–M14, verglichene Alternativen, Messbefehle). F2 und F8 sind entschieden (Entscheidungen 10 und 11); der Implementer beginnt erst, wenn der vorgelagerte CR in `done/` liegt (§4).
Empfehlung: Alternative (b), Statusdateien im Container (`/var/lib/u-boot/status/{vfs,egress,daemon}`), die die Startscripte bei jedem Lauf schreiben und `doctor` mit `docker cp` liest; Container über das Label `devcontainer.local_folder` der Devcontainer-CLI.

**Gemessen** (Colima 0.10.3, Docker 29.5.2, Devcontainer-CLI 0.89.0, arm64): Die Ausgabe von `postCreate`/`postStart` landet nur auf der Standardfehlerausgabe der CLI, nicht in `docker logs` (Annahme aus §1 bestätigt, (a) fällt); der Exit-Code eines Lifecycle-Skripts ist am Container nicht ablesbar (`State.ExitCode` ist der von PID 1; die CLI selbst endet bei Exit `11` des Skripts mit `1`, der Container bleibt `Up`); Labels sind nach `create` unveränderlich ((c) fällt); die CLI setzt `devcontainer.local_folder` und `devcontainer.config_file` mit absoluten Pfaden, `--id-label` ersetzt sie; `docker cp` liest Dateien laufender und gestoppter Container und weist Symlink, FIFO und Verzeichnis als Typ im Tar-Header aus, die Größe steht im Header; `--health-cmd` kürzt die Ausgabe auf 4096 Zeichen und zeigt den Container als `(unhealthy)` ((c′) fällt).
**Angenommen, nicht gemessen:** Übertragbarkeit der Kommandos auf Linux ohne VM und amd64; dieselben Labels bei VS Code; Verhalten bei mehreren Containern desselben Ordners und bei Symlink-Pfaden; `docker ps` im Container bei laufendem verschachteltem Daemon.

| Frage | Antwort | Stand |
|---|---|---|
| F1 Mechanismus | (b); (a), (a′), (b2), (c), (c′), (e), (f), (g) verworfen, (d) ohne (c′) ist (b) | gemessen (a, c, c′), konstruiert (b2, e, f, g) |
| F2 „Befehlsausgabe“ | **Entschieden (Projektinhaber, 2026-10-17): „doctor, plus generate“** — (i) und (ii); (iii) und (iv) entfallen. Vorher offen: Nach Ist-Stand wissen `generate`, `init` und `up` zur Laufzeit des Containers nichts (der Container wird vom Devcontainer-Werkzeug gestartet). Optionen ohne Änderung der Zusage: (i) `doctor` ist die Befehlsausgabe, die Scripte weisen ohnehin auf der Ausgabe des Start-Werkzeugs aus (gemessen M1, M2); (ii) `generate devcontainer`/`init` rufen nach dem Schreiben denselben Port auf und warnen bei bestehendem Container mit `result=fallback`; (iii) auch `up` — nicht empfohlen; (iv) neuer Befehl — Lastenheft-Änderung, ausgeschlossen. Empfehlung der ADR war (i), (ii) als Nachrüstung | entschieden; [`ADR-0022`](../../adr/0022-fallback-ausweisung-statusdateien-im-container.md) Entscheidung 10 |
| F3 Container-Identifikation | `devcontainer.local_folder` = Basisverzeichnis von `doctor`; kein eigenes Label, keine neue `runArgs`; mehrere Treffer: zuletzt erzeugter; kein Treffer: „nicht ermittelbar“ (Info) | Label gemessen (CLI), VS Code angenommen |
| F4 Vertrauensgrenze | geschlossene Schlüssel-/Wertemenge, 512 Byte, strikter Zeitstempel, nur reguläre Datei, nie Rohausgabe; „nicht gemeldet ≠ nicht aufgetreten“ im Handbuch | Verhalten von `docker cp` gemessen |
| F5 Port-Schnitt | neuer Port `DevcontainerInspector` (`FindContainer`, `ReadFile`), im Paket `adapter/driven/docker` neben `Probe`/`Engine`; nur `ps` und `cp`, nie `exec`; `DockerProbe` unverändert; Wertobjekt in `domain`; kein neues Paket, keine neue `depguard`-Regel | Architektur gelesen (`spec/architecture.md` §2.4, §2.6, §5) |
| F6 Stand-Semantik | `postStart` überschreibt je Start, `postCreate` (`vfs`) gilt je Container; gestoppt: gelesen, mit Vermerk und Zeitstempel; ersetzter oder alter Container ohne Dateien: „nicht ermittelbar“ | Persistenz über Stop/Start gemessen (M6) |
| F7 Strenge | `fallback` → Warnung mit `DegradationFallback`, `failed` (Policy `fail`) → Fehler, `ok` → OK, nicht ermittelbar → Info, ungültig → Warnung ohne Kennzeichen; `--strict` bei Warnung Exit `11` wie bei jeder Warnung (`strictWarnFailure`, `driving/cli/doctor.go`); gewollt, keine neue Semantik | gelesen |
| F8 `/dev/fuse` | **Gemessen:** ein fehlendes Device verhindert den Start (`docker run` Exit `127`, `docker start` Exit `1`: `error gathering device information while adding custom device … no such file or directory`; gleich für `/dev/net/tun`). Der `vfs`-Fallback entsteht bei den erzeugten `runArgs` unter Docker nur nach Hand-Änderung (ohne `--device=/dev/fuse`, gemessen M14); die bestehende Linux-Host-Prüfung sagt „fällt auf vfs zurück“ voraus. **Entschieden (Projektinhaber, 2026-10-17): „A: Lastenheft-CR“**; Optionen waren: (A) Change Request, fehlendes `/dev/fuse` wie `/dev/net/tun` als Startfehler zu führen, oder (B) Zielbild halten, `vfs` wird gebaut und ausgewiesen, die Host-Prüfung bleibt ohne Wirkung unter Docker. Gewählt (A), umgesetzt im vorgelagerten CR; der Slice verliert nichts ([`ADR-0022`](../../adr/0022-fallback-ausweisung-statusdateien-im-container.md) Entscheidung 11) | gemessen auf Colima; Linux angenommen |

Weitere Messbefunde für den Implementer: (1) `doctor` aus dem Runtime-Image (`InContainer()`) überspringt die Docker-Prüfungen und erreicht keine Engine — die Ausweisung vom Host gilt nur für einen Host-Install (Binary, Paket, Homebrew); im Devcontainer selbst liest `doctor` die Dateien direkt über `FileSystem`; (2) `/run` und `/var/log` sind keine tmpfs-Mounts, daher ist der Statusort `/var/lib/u-boot/status/` im Dockerfile anzulegen; (3) der erste CLI-Lauf nach einem `sed -i` auf virtiofs meldete einmalig „config not found“, kein Befund über u-boot.
Plan-Tabelle: Die Zeile für `Dockerfile.tmpl` gilt (Statusverzeichnis), `devcontainer_sandbox.go` bleibt bei (b) unverändert (kein Label, keine `runArgs`); die neue Port-Datei heißt `devcontainer_inspector.go` (Name vorgeschlagen), die Wertobjekt-Datei in `domain` `sandbox_status.go`.

**Messung des Implementers (2026-10-18, Colima 0.10.3, Docker 29.5.2, arm64).** `docker ps -a --no-trunc --filter label=… --format '{{.ID}} {{.State}}'` listet mehrere Container desselben Labels nach Erzeugung absteigend (der zuletzt erzeugte zuerst, auch wenn ein älterer läuft); der Adapter nimmt die erste Zeile, ein Integrationstest mit zwei Containern pinnt das. `docker context inspect --format '{{.Endpoints.docker.Host}}'` nennt den Socket des aktuellen Kontexts (`unix:///Users/…/.colima/default/docker.sock`); ein nicht lokaler Wert (`tcp://`, `ssh://`) oder ein solcher `DOCKER_HOST` gilt als nicht ermittelbar, damit ein Container auf einer entfernten Engine mit demselben Ordnerwert kein Falsch-Positiv erzeugt. Beleg für die Startscripte: `make test-docker` (Tests `TestE2E_SandboxDevcontainer_StatusFiles_*`, Lauf mit `-run StatusFiles -v -count=1` im Image `u-boot:test-docker-tools`, alle vier ohne Überspringen); `make test-docker` insgesamt grün; das e2e-Paket brauchte in ruhigen Läufen 477 und (mit `-timeout 30m`, vollständiger Lauf) 635 Sekunden, im Review-Lauf über 600 (Standardlimit von `go test`: zehn Minuten), daher `-timeout 30m` (§1); die drei Container-Tests laufen parallel. Hinweise für die Closure: Der Befund des Sensors `verify-depguard` im Fall `domain-isoliert` (Import-Zyklus, auch am Basiscommit) liegt im Bestand; die Message von `816f6ec` nennt eine CLI-Änderung (Info-Glyph, Summary-Zähler), die `32371a3` zurücknahm; `init` schreibt keine abfragbare Fähigkeit, ADR-Entscheidung 10 verlangt keine Abfrage der Fähigkeiten der überschriebenen `u-boot.yaml`. Die Devcontainer-CLI ist nicht Teil des Repo-Werkzeugs (kein Node im Image `u-boot:test-docker-tools`); die Tests erzeugen den Container selbst mit dem Label `devcontainer.local_folder` und führen die erzeugten Lifecycle-Kommandos per `docker exec` als Container-Benutzer aus. Für die CLI selbst bleiben die Messläufe M1 bis M14 der ADR der Beleg.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): **Erfüllt für `open` → `next` (2026-10-18):** Die Folge-ADR zu F1–F8 ist geschrieben und vom Projektinhaber angenommen (`Accepted`); F2 („doctor, plus generate“) und F8 („A: Lastenheft-CR“) sind beantwortet. **Zusätzliche Bedingung für `next` → `in-progress`:** der CR [`slice-cr-sandbox-fuse-startfehler`](../done/slice-cr-sandbox-fuse-startfehler.md) liegt in `done/` (Lastenheft 0.13.0, Wortlaut von Host-Prüfung und Spezifikations-Zeilen gesetzt); vorher startet der Implementer nicht.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Wählt die ADR (d) oder verlangt F8 eine eigene Messung, wird in Slice 1 (Statusdokument, Skripte, Port, Adapter, `doctor`-Items) und Slice 2 (weitere Quellen, Healthcheck, Label) geschnitten; die Liefer-Punkte der DoD sind schon jetzt an der Grenze von drei.
- `in-progress` → `open` (blockiert): Die Messung zeigt, dass kein engine-vermittelter Weg den Start-Zustand verlässlich liefert (etwa Skript-Ausgabe nicht lesbar, Container nicht identifizierbar). Dann ist die Zusage am Vertrag zu klären — Entscheidung des Projektinhabers, nicht des Plans; das Ergebnis der Messung geht in die ADR.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit
Steering-Loop-Lerneintrag geschrieben (§7). Das Lastenheft steht danach unverändert auf 0.13.0; die Spezifikation trägt dieselbe Zusage ohne die Einschränkung.

## 6. Risiken und offene Punkte

Jedes Risiko bekommt bei der Closure genau einen Ausgang (eingetreten / entfallen / weiter offen).

- **Messung auf Colima trägt nicht:** Die Annahmen (Ort der Skript-Ausgabe, Lesbarkeit per `docker cp`, Label des Devcontainer-Werkzeugs, Health-Log) sind ungemessen; Colima läuft in einer Linux-VM, Linux-CI ohne VM, Docker Desktop und Podman sind nicht zugesagt. — **Ausgang:** offen bis zur Folge-ADR.
- **Floating Container-Zustand:** Das Dokument stammt von einem früheren Start, der Container ist gestoppt, neu gebaut oder ersetzt; `postCreateCommand` läuft nur einmal je Image-Bau. `doctor` kann einen veralteten Stand als aktuellen ausweisen. — **Ausgang:** offen (F6).
- **Container-Identifikation:** u-boot kennt den Container nicht (siehe §1); je Werkzeug (Editor, CLI, Codespaces) können Labels abweichen, mehrere Container desselben Projekts sind möglich (`${devcontainerId}` je Instanz). — **Ausgang:** offen (F3).
- **`docker`-Zugriff im `doctor`:** `doctor` ruft heute `docker version`/`compose version` über `DockerProbe`; neue Kommandos (`ps`, `inspect`, `cp`) fügen Fehlerfälle hinzu (Daemon nicht erreichbar, Container fehlt, Rechte). `doctor` muss einen Bericht liefern (Diagnose liefert immer einen Report, Sicht-Spec) und darf die bestehende Klassifikation der Exit-Codes `11`/`12` nicht verschieben. — **Ausgang:** offen (F5, F7).
- **Vertrauen in das Dokument:** Es stammt aus einem Container mit Agent; es kann falsch („alles ok“) oder bösartig (Steuerzeichen, Übergröße) sein. — **Ausgang:** offen (F4); Mitigation: geschlossenes Schema, Größenlimit, Bereinigung, „nicht gemeldet ≠ nicht aufgetreten“ in der Doku.
- **`--quiet`-Sichtbarkeit:** Ein neues Item sieht unter `--quiet` nur, wer `DegradationFallback` trägt; „nicht ermittelbar“ ist kein Fallback und bleibt unter `--quiet` stumm — das ist gewollt, muss aber im Handbuch stehen, sonst liest es sich als „kein Fallback“. — **Ausgang:** offen.
- **Docker-only-Tests:** Die Nachweise brauchen eine echte Engine (`make test-docker`, Build-Tag `docker`); `make gates` führt sie nicht aus. Ein Lauf ohne Engine darf nicht als grün berichtet werden ([`harness/README.md`](../../../../harness/README.md) §Sensors). — **Ausgang:** offen.
- **Änderung der erzeugten Ausgabe bestehender Projekte:** Neue `runArgs` (Label, `--health-cmd`) oder geänderte Startscripte erscheinen nach dem nächsten `u-boot generate devcontainer` (managed Blocks, Replay-Cases). — **Ausgang:** offen; im CHANGELOG zu nennen.
- **Spec-Nachzug ohne Verweise nach unten:** Spezifikation und Sicht dürfen weder ADR noch diesen Slice nennen; die Entscheidung steht in der ADR (`Schärft:`), die Spec trägt die Aussage selbst. — **Ausgang:** offen; Prüfung durch `make docs-check` (Referenzmatrix) und Reviewer.
- **`--device=/dev/fuse`-Befund (F8):** gemessen (M13): ein fehlendes `/dev/fuse` verhindert den Start, der `vfs`-Fallback ist unter den erzeugten `runArgs` unerreichbar. — **Ausgang:** eingetreten: [`slice-cr-sandbox-fuse-startfehler`](../done/slice-cr-sandbox-fuse-startfehler.md) (Lastenheft-CR, Entscheidung 2026-10-17).
- **Abfrage in `generate`/`init` (F2 ii):** zusätzliche Engine-Aufrufe (`docker ps`, `docker cp`) in zwei Befehlen, die bisher keine Engine brauchen: Latenz, hängende Engine, Rauschen in `--json`. — **Ausgang:** offen; Zeitlimit je Aufruf im Adapter, alle Fehler still, Test „keine Engine“ und „--dry-run ruft den Port nicht“.
- **Warnung beschreibt einen alten Container:** Die Abfrage nach `generate` liest den Stand des letzten Starts des bestehenden Containers; nach dem Neubau kann er überholt sein. — **Ausgang:** offen; Wortlaut der Warnung nennt „Stand des letzten Starts“ ([`ADR-0022`](../../adr/0022-fallback-ausweisung-statusdateien-im-container.md) Entscheidung 10).
- **Reihenfolge CR → Slice:** der Slice darf nicht vor dem CR starten, sonst ändern beide dieselben Spec-Zeilen. — **Ausgang:** offen; Start-Trigger in §4.
- **Umfang:** drei Liefer-Punkte bei vier Quellen (Skripte, Port/Adapter, `doctor`, Spec/Doku) sind die Obergrenze. — **Ausgang:** offen; Rückführung siehe §4.

## 7. Closure-Notiz

Wird bei der Closure gefüllt.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind laut §3
`internal/hexagon/domain` (DOM), `internal/hexagon/application` (APP), `internal/hexagon/port` (PORT),
`internal/adapter/driven` (ADP), `internal/adapter/driving/cli` (CLI, nur Tests), `internal/e2e` (E2E),
`spec/`, `harness/`, `docs/plan/` (DOK) — alle Greenfield laut
[`harness/conventions.md`](../../../../harness/conventions.md#mr-006) (Abschnitt Modus-Deklaration), jede erfüllt die Schwelle
≥ 2 von 3 Achsen; keine Ausdifferenzierung nötig. `docs/user/` und `CHANGELOG.md` sind öffentliche Verträge, keine Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register
([`../observations/`](../observations/README.md)) durchgegangen — Verzeichnisse `BEO-CLI`, `BEO-DOK`, `BEO-E2E`, `BEO-WZG` liegen vor; keiner betrifft die Ausweisung von Fallbacks.
`BEO-DOK/plan-aussage-ueber-code` (Plan-Aussagen über Code) ist beachtet: Jede Aussage über Code in §1 und §3 steht mit Fundstelle; was ungelesen ist, steht als *Annahme* und trägt eine Messung.

Alle berührten Sub-Areas GF.
