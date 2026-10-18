# ADR 0025: Egress-Liste im laufenden Container — Ist-Stand gemessen, Neubau statt Neu-Laden (R0), weitere Wege offen gehalten

**Status:** Proposed

**Datum:** 2026-10-18

**Autor:** pt9912

**Bezug:** [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [ADR-0012](0012-devcontainer-egress-firewall.md), [ADR-0022](0022-fallback-ausweisung-statusdateien-im-container.md)

**Schärft:** [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (nur der Satz zur Grenze der Restriktion: „Prozesse mit der nötigen Capability können sie aufheben“; der Weg „Liste ändern“ steht nicht im Lastenheft, sondern im Handbuch) sowie — als Stellen, die nach `Accepted` nachzuführen sind — [`SPEC-045`](../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) (Absatz *Grenze*: Schreibbarkeit der Statusdatei). **Ergänzt**, ersetzt nichts: [ADR-0012](0012-devcontainer-egress-firewall.md) und [ADR-0022](0022-fallback-ausweisung-statusdateien-im-container.md) bleiben unverändert `Accepted`.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Der Projektinhaber fragte am 2026-10-18, ob sich die Egress-Liste im laufenden Sandbox-Container vom Host aus neu laden lässt, ohne Neustart und ohne dass ein Agent im Container sich damit selbst Hosts freischalten kann. Er entschied danach (2026-10-18, Vorschläge bestätigt): **Variante R0, Guardrail-Status quo** — kein Neu-Laden-Befehl vom Host und kein Versprechen „der Agent kann sich nicht selbst freischalten“; ein Agent mit `sudo` kann die Regeln heute schon aufheben, [ADR-0012](0012-devcontainer-egress-firewall.md) und das Handbuch nennen das. Die Messung soll darlegen, was der Ist-Stand ist, und ob ein Container-Neustart statt eines Image-Neubaus nach geänderter Liste genügt. Die Wege R1 bis R4 und S bleiben dokumentiert und werden **nicht entschieden**.

**Stand im Code (gelesen).** `generateDevcontainer` hängt bei Egress `.devcontainer/egress-init.sh` an die erzeugten Dateien an (`internal/hexagon/application/generate.go`, Kommentar „egress restriction script baked via COPY“); `Dockerfile.tmpl` (Block `{{- if .Egress}}`) installiert `dnsmasq-base nftables` und kopiert das Script mit `COPY egress-init.sh /usr/local/bin/u-boot-egress-init`; `egress-init.sh.tmpl` trägt die Hosts als `hosts="…"` (Allowlist). `devcontainerTemplateData` (`devcontainer_sandbox.go`) setzt `data.PostStart` auf `sudo sh /usr/local/bin/u-boot-egress-init` (Egress zuerst, `" && "`-Kette), `egressRelaxation()` liefert `--cap-add=NET_ADMIN`. Das Script ist laut seinem Kommentar „re-run safe“: es beendet den vorigen `dnsmasq` mit `kill` (ohne zu warten), löscht beide Tabellen mit `nft delete table`, schreibt die `dnsmasq`-Konfiguration neu, legt die Tabellen neu an, startet `dnsmasq` und setzt `/etc/resolv.conf` auf `127.0.0.1`; `write_status` schreibt `/var/lib/u-boot/status/egress` über eine Temporärdatei und `mv`. `TestGenerateDevcontainer_SandboxEgress` (`generate_sandbox_test.go`) pinnt, dass das Script die Hosts trägt und das Dockerfile das `COPY` enthält. Der Host führt heute nie ein Programm im Container aus: `DevcontainerInspector` (`internal/hexagon/port/driven/devcontainer_inspector.go`) hat nur `FindContainer` und `ReadFile` ([`SPEC-045`](../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)).

**Messung (2026-10-18).** Docker-Server 29.5.2 (linux/arm64) auf Colima 0.10.3 (VM Ubuntu 24.04.4, Kernel 6.8.0-117-generic, 12 CPUs), Docker-Client 29.9.1, Dev-Container-CLI 0.89.0, Basis `mcr.microsoft.com/devcontainers/base@sha256:6ee2b4097c75d33e9a94d941577a3286e5d629434fba384688b06daeeaaeac99` (Debian 13.7, arm64), `dnsmasq` 2.91, `nftables` 1.1.3; Projekt und `egress-init.sh` von `u-boot` (`0.11.0-dev`) erzeugt. Skripte und Rohdaten: [Messprotokoll](../planning/observations/BEO-DOK/messprotokoll-ausserhalb-des-repos/evidence/messung-slice-v1-sandbox-egress-neuladen-entwurf/README.md). **Gemessen** heißt: auf diesem Host ausgeführt und ausgegeben.

| Frage | Ergebnis (Rohdaten) |
|---|---|
| M1 ohne `sudo` (als `vscode`) | `sudo -n true` Exit 0 und `sudo -l`: `(root) NOPASSWD: ALL` — das Base-Image gibt `vscode` passwortloses `sudo`; ohne `sudo`: `CapEff` 0, `nft list tables` „Operation not permitted“, Anhängen an `/usr/local/bin/u-boot-egress-init`, `/etc/u-boot-egress-dnsmasq.conf`, `/etc/resolv.conf`, `/etc/resolv.conf.u-boot-orig` und die Pid-Datei „Permission denied“, `SIGTERM` an den `dnsmasq` (Benutzer `nobody`) „Operation not permitted“ (Prozess lebt); `https://api.github.com` 200, `https://example.com` nicht auflösbar, `https://1.1.1.1` Timeout (`rohdaten-02-m1.txt`) |
| M1 ohne `sudo`, zwei Befunde neben der Regel | **(a)** Die DNS-Abfrage direkt an den Upstream-Resolver (`192.168.5.1`, aus `/etc/resolv.conf.u-boot-orig`) löst `example.org` auf (zwei A-Einträge); die Verbindung zur Adresse läuft in den Timeout. Die Namensauflösung ist damit auch ohne `sudo` ein offener Pfad nach außen; gemessen ist die Auflösung, nicht der Durchsatz eines Datenkanals. **(b)** `/var/lib/u-boot/status/egress` gehört `root` (nicht schreibbar), das Verzeichnis gehört `vscode`: `vscode` löschte die Datei und legte eine neue mit `result=ok` an (`rohdaten-02-m1.txt`) |
| M1 mit `sudo` | `sudo nft delete table ip u_boot_egress` und `… ip6 u_boot_egress6` Exit 0, danach `https://1.1.1.1` 301; `https://example.com` bleibt nicht auflösbar, solange `/etc/resolv.conf` auf `127.0.0.1` zeigt; nach `sudo cp /etc/resolv.conf.u-boot-orig /etc/resolv.conf` 200 (`rohdaten-02-m1.txt`) |
| M2 erneuter Lauf des gebackenen Scripts, andere Liste | Lauf 50 bis 78 ms (fünf Läufe, `dauer_s`); **Fenster ohne Filter-Regelwerk 24, 30, 30, 34 und 45 ms** (Sonde im Abstand von rund 5 µs, Zustand `offen` bis `gesperrt`); in dieser Zeit kein `nft`-Regelwerk; Ausfall des lokalen `dnsmasq` (`Connection refused` bis Antwort) 28, 30, 34, 42 und 55 ms, in vier der fünf Läufe dazwischen ein Zustand „Loopback gesperrt“ (`sendfehler`), bevor die Regel `oifname lo accept` steht (`rohdaten-03-m2.txt`, `rohdaten-03-m2-wiederholung.txt`, Abschnitt A) |
| M2 Variante (Messskript, kein Produktcode): atomarer Wechsel in einer nft-Transaktion, `dnsmasq` erst nach dem Ende des alten | Lauf 16 bis 24 ms; in fünf Läufen **kein** Zustand `offen` beobachtet; Ausfall des lokalen `dnsmasq` höchstens 24 ms (Dauer des längsten Laufs) (Abschnitt E) |
| M2 Startrennen (`kill` ohne Warten) | 50 Läufe des gebackenen Scripts hintereinander mit wechselnder Liste: 0 mit auffälligem Zustand (genau ein lebender `dnsmasq`, antwortet, neue Konfiguration), 0 mit Fehlermeldung; dieselben 50 mit der Variante ebenso (Abschnitte B, F); der alte `dnsmasq` endet 1,9 bis 3,4 ms nach `SIGTERM` (20 Messungen, Abschnitt G), das ganze Script dauert 50 bis 78 ms. Ohne Last gemessen: das Rennen ist nicht aufgetreten, ein Beleg für Rennfreiheit unter Last ist es nicht |
| M2 laufende Verbindung | Ein Download von `deb.debian.org` (300 kB/s) lief über einen Neu-Lauf mit einer Liste **ohne** diesen Host weiter (Bytes wuchsen 9 s lang; die Regel `ct state established,related accept` gilt für bestehende Verbindungen); eine neue Verbindung zum selben Host: `Could not resolve host` (Abschnitt C) |
| M2 Set nach dem Lauf | Das nftables-Set `allowed` ist nach dem Lauf leer (vorher ein Eintrag mit `expires 59m59s`); eine Verbindung zu einer zuvor aufgelösten Adresse **ohne neue Auflösung** scheitert (`http=000`), mit Auflösung `200`, danach wieder ohne (Abschnitt D) |
| M3 Zustellung vom Host | `docker cp` einer Host-Datei nach `/etc` ergab Eigentümer `501 dialout` (UID und GID des Host-Benutzers; `vscode` schreibt nicht und löscht nicht); `docker cp -a`: `vscode vscode`; Tar-Strom über `docker cp -` mit `owner=0`: `root root` (`vscode` schreibt nicht), mit `owner=1000`: `vscode` schreibt. Eine mit `docker cp` in `/var/lib/u-boot/status` (gehört `vscode`) abgelegte Datei ersetzte `vscode` ohne Weiteres (löschen und neu anlegen). `docker exec -u root … sh /usr/local/bin/u-boot-egress-init` Exit 0, die Liste galt; eine Kopie des Scripts mit anderer Liste per `docker cp` über das gebackene Script, danach der `postStartCommand`-Aufruf: `example.org` löste auf. Gestoppter Container: `docker exec` „is not running“, `docker cp` in und aus dem gestoppten Container Exit 0. Nach `docker start` ohne `postStartCommand`: keine Tabellen, `resolv.conf` wie von Docker, `https://example.com` 200; nach dem Aufruf 000 (`rohdaten-04-m3.txt`). Eine entfernte Engine ist nicht gemessen |
| M4 Liste ändern, Container neu starten | Mit der Dev-Container-CLI 0.89.0 und einem von `u-boot` erzeugten Projekt: `u-boot generate devcontainer` mit geänderter Liste schreibt `.devcontainer/egress-init.sh` (Host-Datei trägt `example.org`); der laufende Container blieb unverändert. **Stoppen und `devcontainer up`**: derselbe Container, dasselbe Image, `postStartCommand` lief (Zeitstempel der Statusdatei neu), die **alte** Liste galt (`example.org` nicht erreichbar). `devcontainer up` gegen den laufenden Container: nichts geändert. **`devcontainer up --remove-existing-container`**: neuer Container, **neue Image-Kennung**, die **neue** Liste galt (`example.org` 200). Der Satz im Handbuch „starten den Container neu“ stimmte nicht (`rohdaten-05-m4.txt`) |

**Nicht gemessen:** Docker Desktop, Linux-Host ohne VM, `amd64`, eine entfernte Engine, VS Code (welche Editor-Aktion dem Neuaufbau entspricht), ein Container mit `nestedRuntime`, Last des Hosts als Auslöser des Startrennens, der Durchsatz eines Datenkanals über die Namensauflösung, mehrere Läufe je Fall bei M1, M3 und M4.

## Entscheidung

**R0 — Guardrail-Status quo** (Projektinhaber, 2026-10-18; Bestätigung dieses Wortlauts steht aus, daher `Proposed`).

1. **Kein Neu-Laden im laufenden Container.** `u-boot` bekommt keinen Befehl, der ein Programm im Container ausführt oder eine Liste hineinlegt; der Satz „der Host führt nie ein Programm im Container aus“ ([`SPEC-045`](../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)) bleibt unverändert. Der Weg „Liste ändern“ ist `u-boot generate devcontainer` und **Neuaufbau des Containers** (M4); das Handbuch sagt das, mit Messstand.
2. **Kein Versprechen, dass ein Agent sich nicht selbst freischalten kann.** Gemessen: Mit dem passwortlosen `sudo` des Base-Images genügt ein `nft`-Aufruf plus das Zurücksetzen von `/etc/resolv.conf` (M1). Ohne `sudo` kann ein Prozess die Regeln nicht ändern (M1), die Namensauflösung über den Upstream-Resolver bleibt aber offen (M1 a); das Lastenheft bleibt unverändert: Die Restriktion ist ein Guardrail und keine Sandbox-Grenze ([`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)).
3. **R1 bis R4 und S bleiben offen und werden nicht entschieden.** Die Tabelle unten trägt, was die Messung zu ihnen sagt; jede spätere Wahl braucht eine eigene ADR und für R1 und R2 einen Lastenheft-CR mit freigegebenem Wortlaut. Folge-Slices zu R1 und R2 entstehen nicht.

## Verglichene Alternativen

Jede Zeile nennt, was die Messung für den Weg trägt; „nicht gebaut“ heißt, dass die Zelle aus M1 bis M4 abgeleitet ist, nicht aus einer Umsetzung.

| Option | Pro | Contra |
|---|---|---|
| **R0** — kein Neu-Laden, Neubau, Handbuch berichtigt | keine neue Fähigkeit, kein Port, keine Lastenheft-Änderung; ein ehrlicher Satz zum Neuaufbau (M4); Guardrail-Aussage unverändert | jede Listenänderung braucht einen Neuaufbau des Images (Dauer für diesen Aufbau nicht gemessen) |
| R1 — Host-Befehl (`docker exec -u root` mit zugestellter Liste) | Auslöser auf dem Host; `docker exec -u root` führt das Script aus (M3); mit Tar-Strom `owner=0` liegt die Liste Root-eigen in `/etc` (M3) | neuer Driven-Port und Spezifikationsänderung (erster Programmaufruf im Container); Reload-Script im Image (einmal Neubau); `docker cp` einer Host-Datei ist **kein** verlässlicher Weg zu `root` als Eigentümer (M3: `501 dialout`, bei `-a` `vscode`); das gebackene Script hat ein Fenster ohne Filter von 24 bis 45 ms und einen Ausfall des lokalen `dnsmasq` von 28 bis 55 ms (M2), die atomare Variante zeigte kein offenes Fenster; das Set ist danach leer, bestehende Verbindungen laufen weiter (M2); löst das Versprechen „kein Selbst-Freischalten“ nicht ein, solange `sudo` ohne Passwort gilt (M1) |
| R2 — zwei Schritte (Host legt die Liste per `docker cp -` Root-eigen ab, Nutzer ruft im Container `sudo sh …` auf) | kleinerer Port (nur Schreiben); Liste Root-eigen (M3) | eine Aktion im Container; sonst wie R1; ein Verzeichnis des Benutzers (`/var/lib/u-boot/status`) taugt nicht als Ablage, `vscode` ersetzt die Datei dort (M3) |
| R3 — im Container ohne `sudo` auslösbar | — | verworfen: ein Prozess im Container könnte sich selbst Hosts freischalten (widerspricht dem Wunsch des Projektinhabers) |
| R4 — Restriktion außerhalb des Containers (Proxy, Netz-/DNS-Ebene) | auch mit `sudo` nicht aufhebbar; nicht vom Container abhängig | Host-Konfiguration ist Nicht-Ziel von `u-boot` ([ADR-0012](0012-devcontainer-egress-firewall.md) Alternative D); nicht gemessen |
| S — `sudo` für `vscode` entziehen (Zusatz zu R1/R2) | nähme dem Agenten den Weg über `sudo` (M1: ohne `sudo` unveränderbar) | ändert Benutzerrechte und alle drei Startscripte (`sudo sh …` im `postStartCommand`); nicht gemessen, eigener Vorgang |

## Konsequenzen

- Positiv: Das Handbuch stimmt mit der Messung überein; die Aussage des Lastenhefts zur Grenze der Restriktion wird durch M1 gedeckt; der Entscheidungsraum für ein späteres Neu-Laden trägt gemessene Zellen (Eigentümer beim `docker cp`, Fenster, Set, bestehende Verbindungen) statt Annahmen.
- Negativ: Ein Neuaufbau ist der einzige Weg zu einer anderen Liste; **gemessene Grenzen, die das Lastenheft und das Handbuch heute nicht nennen:** die Namensauflösung über den Upstream-Resolver ist ohne `sudo` offen (M1 a), ein Prozess ohne `sudo` kann die Statusdatei ersetzen (M1 b; der Absatz *Grenze* in [`SPEC-045`](../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) sagt „mit Root-Rechten frei schreibbar“), das Startscript ist bei einem zweiten Lauf nicht atomar (M2). Ob eine davon ins Handbuch oder in die Spezifikation gehört, ist offen (Entscheidung des Projektinhabers bzw. Spezifikations-Slice der Welle).
- Folgepflicht: Handbuch und `CHANGELOG.md` (erledigt im selben Slice); die Spezifikation führt den Absatz *Grenze* von [`SPEC-045`](../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) im Spezifikations-Slice der Welle nach (ein Eintrag zu Varianten, die nicht gebaut wurden, entsteht nicht). Ein Test zu M4 entsteht nicht: Das Verhalten ist das der Dev-Container-CLI und des Docker-Caches, nicht des Codes; die Zusage im Handbuch trägt Messstand und Grenze statt eines Tests, und was `u-boot` selbst beisteuert (Script mit Liste im Image) pinnt `TestGenerateDevcontainer_SandboxEgress`.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Unit-Test | Das erzeugte Script trägt die Hosts, das Dockerfile das `COPY` des Scripts (`TestGenerateDevcontainer_SandboxEgress`) | `make test` |
| Messprotokoll im Repo | M1 bis M4 mit Skripten und Rohdaten wiederholbar (Voraussetzungen im README des Protokolls); kein Gate | — |

## Re-Evaluierungs-Trigger

Der Projektinhaber verlangt ein Neu-Laden im laufenden Container; Nutzerberichte, dass der Neuaufbau nach einer Listenänderung zu aufwendig ist; die Entscheidung, `sudo` für den Container-Benutzer zu entziehen (S) oder die Restriktion außerhalb des Containers zu führen (R4); ein Messstand außerhalb Colima/arm64 (Docker Desktop, Linux-Host, `amd64`) oder eine andere Version der Dev-Container-CLI, der M4 widerlegt.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-18 | Messung M1 bis M4 (Colima 0.10.3, Docker 29.5.2, Dev-Container-CLI 0.89.0, arm64); Entwurf mit R0 als Entscheidung des Projektinhabers (Vorschläge bestätigt), `Proposed` | Vereinbarung mit dem Projektinhaber |
