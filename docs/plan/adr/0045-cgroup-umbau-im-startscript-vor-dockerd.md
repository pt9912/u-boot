# ADR 0045: Umbau der cgroup im Startscript `docker-init.sh` vor `dockerd` — Startort `postStartCommand`, Fehlschlag best-effort ohne Statusänderung

**Status:** Accepted

**Datum:** 2026-11-04

**Autor:** pt9912 (Rolle Architect)

**Bezug:** [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Container-Runtime im Sandbox-Devcontainer, Degradationstabelle), [`LH-FA-DEV-001`](../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) (erzeugte Dateien)

**Schärft:** [`SPEC-056`](../../../spec/spezifikation.md#spec-056--cgroup-umbau-im-startscript-docker-initsh-bedingung-schritte-und-fehlschlag) (neu: Bedingung, Schritte, Ort und Fehlschlag des Umbaus), [`SPEC-045`](../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) (der Umbau schreibt keine Statusdatei und ändert `daemon` nicht) und [`SPEC-053`](../../../spec/spezifikation.md#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs) (Wirkung der Grenzen auf den verschachtelten Docker)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR); Form nach [`MR-008`](../../../harness/conventions/MR-008-adr-form.md) (neue ADRs MADR).

---

## Kontext

In einem Sandbox-Devcontainer mit `nestedRuntime: docker` scheitert jeder verschachtelte Start mit `--memory` oder `--cpus` mit Exit-Code `125` (`unable to apply cgroup configuration: cannot enter cgroupv2 "/sys/fs/cgroup/docker" with domain controllers -- it is in threaded mode`).

**Gemessen** (an einem Container, dem Entwicklungs-Devcontainer dieses Repos, 2026-11-04; nicht an anderen Hosts oder Engines):

- Die Wurzel-cgroup des Containers war `domain threaded`: Prozesse in der Wurzel und die threaded Controller `cpuset cpu pids` in `cgroup.subtree_control`, `memory` nicht.
- Die Folge *(1) `-cpuset -cpu -pids` aus `subtree_control` der Wurzel, (2) Unter-cgroup `init` anlegen und alle Prozesse der Wurzel hineinschreiben, (3) alle Controller aus `cgroup.controllers` in `subtree_control` aktivieren, (4) `dockerd` starten* setzte die Wurzel auf `domain`; `docker run --memory=64m --cpus=0.5` zeigte danach `memory.max = 67108864` und `cpu.max = 50000 100000`. Ohne Schritt 1 schlugen 2 und 3 fehl.
- Der Zustand lebt in der cgroup des Containers und endet mit ihr; er muss bei jedem Start neu hergestellt werden. Das ist das Muster des `dind`-Entrypoints.

**Gelesen (2026-11-04, nur lesend):** Das Hilfs-Image von `u-boot devcontainer up` (`u-boot-devcontainer-cli:0.89.0-f1a94cce795e`, Dev Container CLI 0.89.0, gepinnt in `internal/adapter/driven/docker/devcontainercli/package.json`) enthält in `/opt/devcontainer-cli/node_modules/@devcontainers/cli/dist/spec-node/devContainersSpecCLI.js` im Pfad für einen einzelnen Container die Argumente `["--entrypoint","/bin/sh"]` und `["-c","echo Container started …"]`, und nur der Zweig `if(t.overrideCommand===!1){ … Config.Entrypoint … Config.Cmd … }` hängt Entrypoint und Cmd des Images an. Im Compose-Pfad der CLI ist `overrideCommand` wahr gleichbedeutend mit „Entrypoint und Cmd des Images entfallen“. Gemessen am laufenden Container: PID 1 ist `/bin/sh -c "echo Container started …"`.

**Nicht belegt:** ob die Erweiterung des Editors (eigene Kopie der CLI) in jeder Version denselben Pfad nimmt; der Wert von `overrideCommand` ist in der erzeugten `devcontainer.json` nicht gesetzt (`templates/devcontainer/devcontainer.json.tmpl`, gelesen: kein Vorkommen).

**Stand im Code (gelesen):** `internal/hexagon/application/templates/devcontainer/docker-init.sh.tmpl` startet `dockerd` per `setsid nohup` erst, wenn `ready` (`docker info`) fehlschlägt und kein lebender `dockerd` mit PID-Datei vorhanden ist; `Dockerfile.tmpl` kopiert das Skript im Block `{{- if .NestedDocker}}` nach `/usr/local/bin/u-boot-docker-init`; `devcontainerTemplateData` (`devcontainer_sandbox.go`) setzt `postStartCommand` auf `sudo sh /usr/local/bin/u-boot-docker-init` (mit Egress davor verkettet). Die Statusdatei `daemon` kennt `ok`, `fallback`, `failed` ([`SPEC-045`](../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)); `fallback` bedeutet dort „Container ohne Daemon“, `failed` Exit-Code `11` unter `fail`. [`SPEC-053`](../../../spec/spezifikation.md#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs) sagt die Wirkung der Grenzen auf den verschachtelten Docker ausdrücklich nicht zu; das Lastenheft nennt cgroups nirgends (Suche in `spec/lastenheft.md`, kein Treffer).

## Entscheidung

Wir wählen **den Umbau im bestehenden Skript `docker-init.sh`, vor dem Start von `dockerd`, ausgelöst durch `postStartCommand` bei jedem Start, mit Fehlschlag als best-effort**:

1. **Startort (a): `docker-init.sh`, nicht `ENTRYPOINT`.** Das Skript läuft bei jedem Start (`postStartCommand`) und damit nach jeder Neuanlage und jedem Neustart des Containers; der Zustand lebt nur bis zum Ende der cgroup und wird so jedes Mal neu hergestellt. Der `ENTRYPOINT` (b) trägt nicht: die CLI ersetzt ihn mit `--entrypoint /bin/sh`, solange `overrideCommand` nicht `false` ist (Fundstelle oben); ihn zu erzwingen hieße `overrideCommand: false` zu setzen, die Start-Hülle der Erweiterung („Container started“, Warten) aufzugeben und eine zweite Start-Kette neben `postStartCommand` zu führen — Mehraufwand ohne Gewinn, der Umbau muss ohnehin nur vor `dockerd` laufen, das dasselbe Skript startet. (c) „beides nicht“ lässt den Fehler des Befunds bestehen.
2. **Ort im Skript.** Der Umbau steht im Zweig, der `dockerd` startet (nach der Bereitschaftsprüfung, nach dem Entfernen der verwaisten PID- und Socket-Dateien, unmittelbar vor `setsid nohup dockerd`). Der Zweig wird nur betreten, wenn kein laufender `dockerd` erkannt wird (Bereitschaftsprüfung und PID-Datei); ein bereiter oder startender Daemon löst keinen Umbau aus. Ein `dockerd`, den das Skript nicht erkennt (etwa ohne PID-Datei), ist die Ausnahme: der Umbau liefe dann unter ihm.
3. **Bedingung (Idempotenz).** Nur bei cgroup v2 (`/sys/fs/cgroup/cgroup.controllers` vorhanden) und nur, wenn `memory` in `cgroup.controllers`, aber noch nicht in `cgroup.subtree_control` der Wurzel steht. Steht `memory` schon dort, ist die Wurzel bereits umgebaut oder vom Host so geliefert; der Umbau entfällt. Alle Schritte sind POSIX-`sh` (das Skript läuft als `sh`, Debian `dash`).
4. **Schritte.** Wie im Befund: `-<Controller>` für jeden Eintrag von `subtree_control` der Wurzel (im Befund `-cpuset -cpu -pids`; der Schritt entfällt bei leerem Inhalt); Unter-cgroup `init` anlegen und die Prozesse der Wurzel einzeln hineinschreiben (ein Fehler bei einem einzelnen Prozess, beendet oder nicht verschiebbar, wird ignoriert und zeigt sich erst als Fehlschlag des dritten Schritts); alle Controller aus `cgroup.controllers` aktivieren. Die Einzelheiten stehen in der Spezifikation.
5. **Fehlschlag: best-effort, ohne Statusänderung, ohne Abbruch.** Scheitert ein Schritt, schreibt das Skript eine Zeile auf stderr (`u-boot sandbox: cgroup preparation failed …`), bricht den Umbau ab und startet `dockerd` wie bisher. Exit-Code und Statusdatei `daemon` bleiben die des Daemon-Starts; `onUnavailable` wirkt nicht auf den Umbau. Begründung: (i) Die Zeile `daemon` beschreibt die Verfügbarkeit des Daemons; `fallback` hieße für `doctor` „Container ohne Daemon“ (Kennzeichen `DegradationFallback`) und wäre bei laufendem Daemon eine falsche Aussage, `failed` mit Exit-Code `11` unter `fail` machte einen funktionierenden Daemon zum Fehler. (ii) Die Statusdatei trägt keinen Freitext und genau drei Werte, streng gelesen ([`SPEC-045`](../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)); ein vierter Wert oder ein Grund wäre eine Änderung dieses Formats und der Degradationstabelle, also des Vertrags. (iii) Die Wirkung der Grenzen auf verschachtelte Container ist im Lastenheft nicht zugesagt und in [`SPEC-053`](../../../spec/spezifikation.md#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs) ausdrücklich nicht zugesagt; ein Ausbleiben ist keine Degradationsstufe der Container-Runtime. Daraus folgt: **kein Lastenheft-CR nötig**. Würde die Wirkung auf den verschachtelten Docker später zugesagt, ist das ein eigener Vertragsvorgang.
6. **Zusage-Umfang.** Zugesagt wird nur, was die Belege der Spezifikation tragen: der Umbau läuft vor `dockerd`, ist bedingt und ändert Exit-Code und Statusdatei nicht. Dass die Grenzen in jedem Container ankommen, wird nicht zugesagt (gemessen: ein Container).

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun | keine Änderung | jede Grenze (`--memory`, `--cpus`) bricht den verschachtelten Start in diesem Container; zwei Docker-Tests bleiben rot |
| **B — Umbau in `docker-init.sh` vor `dockerd` (gewählt)** | läuft bei jedem Start; ein vorhandener Pfad (`postStartCommand`, Skript im Image); keine neue Start-Kette; Test an der erzeugten Datei möglich | der Umbau läuft erst nach dem Start-Hook der Erweiterung, Prozesse der Erweiterung in der Wurzel werden mitverschoben (Risiko unten) |
| C — `ENTRYPOINT` im Image (dind-Muster) | der Umbau läuft vor allen anderen Prozessen | die CLI ersetzt den `ENTRYPOINT` (`--entrypoint /bin/sh`), solange `overrideCommand` nicht `false` ist; sonst eigener Start-Mechanismus neben `postStartCommand` |
| D — Fehlschlag über `onUnavailable` und Statusdatei `daemon` | Fehlschlag sichtbar in `doctor` | falsche Aussage bei laufendem Daemon oder neuer Statuswert; Vertragsänderung (Lastenheft-CR) |
| E — `--cgroupns=host` oder Host-Einstellung | umgeht den Umbau | hebt die Isolation auf; ausgeschlossen |

## Konsequenzen

- Positiv: ein Pfad, jeder Start, kein neuer Mechanismus; der Umbau ist an der erzeugten Datei testbar; Fehlschlag gefährdet den Daemon-Start nicht.
- Negativ (akzeptiert): (a) Der Fehlschlag ist nur auf stderr sichtbar (Protokoll des Start-Hooks), nicht in `doctor`. (b) Der Umbau verschiebt alle Prozesse der Wurzel in `init`, auch die der Erweiterung, die vor `postStartCommand` gestartet sind. (c) **Ungemessen:** nach dem Umbau ist die Wurzel eine Domain-cgroup mit gesetztem `subtree_control`; nach Kenntnis des Kernels nimmt sie keine neuen Prozesse mehr auf (`EBUSY`) — `docker exec` und neue Terminals der Erweiterung könnten betroffen sein. Gemessen ist das nicht. **Akzeptiertes Negativ** mit Auffangnetz: der Docker-Test der Spezifikation führt nach dem Umbau ein `docker exec` aus; schlägt es fehl, greift der Re-Evaluierungs-Trigger.
- Folgepflicht: [`SPEC-056`](../../../spec/spezifikation.md#spec-056--cgroup-umbau-im-startscript-docker-initsh-bedingung-schritte-und-fehlschlag) mit Beleg-Tabelle (angelegt); [`SPEC-053`](../../../spec/spezifikation.md#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs) verweist auf sie.

### Architektur-Befund (hexagonale Passung)

Keine Schichtänderung: die Änderung betrifft nur ein Template im Application-Service (`internal/hexagon/application/templates/devcontainer/`); kein neuer Port, kein Adapter, keine Importregel.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Tests (`internal/hexagon/application`) | Umbau steht vor `dockerd`, Bedingung, keine Statusänderung, nur bei `nestedRuntime: docker`; Beleg-Tabelle in [`SPEC-056`](../../../spec/spezifikation.md#spec-056--cgroup-umbau-im-startscript-docker-initsh-bedingung-schritte-und-fehlschlag) | `make test` |
| Docker-Test (Build-Tag `docker`) | `memory.max` und `cpu.max` im verschachtelten Container nach dem Umbau, `docker exec` danach | `make test-docker` |
| `docs-check` | Kennungen und Anker | `make docs-check` |

## Re-Evaluierungs-Trigger

- **Messung nach einer Neuanlage des Devcontainers** (`make dev-up RECREATE=1`, nur der Projektinhaber löst sie aus): Zeigt sie, dass (i) der Zustand nach Neuanlage nicht hergestellt ist (die Wurzel bleibt `domain threaded`, weil etwas vor `postStartCommand` Prozesse einfügt), (ii) `docker exec` oder ein neues Terminal nach dem Umbau mit `EBUSY` scheitert, oder (iii) die Erweiterung durch das Verschieben ihrer Prozesse ausfällt, gilt diese Entscheidung nicht: Folge-ADR mit `Supersedes`, nicht Änderung dieser.
- Die CLI oder die Erweiterung nimmt den `ENTRYPOINT` des Images mit (`overrideCommand: false` als Standard): Startort (b) neu bewerten.
- Eine Zusage der Wirkung auf verschachtelte Container im Lastenheft: Fehlschlag-Verhalten (Entscheidung 5) neu bewerten.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-11-04 | Accepted | Entscheidung des Startorts und des Fehlschlag-Verhaltens |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-0045` (Baseline-Regelwerk `modul-04-adrs.md`).
