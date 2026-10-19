# Messprotokoll: Pfad zum Upstream-Resolver für Prozesse ohne `sudo`

Rohdaten und Ablauf der Messung (Fragen D1 bis D5) des Slice
[`slice-v1-egress-dns-upstream-schliessen`](../../plan/planning/in-progress/slice-v1-egress-dns-upstream-schliessen.md)
(Auswertung, Einordnung und Empfehlung: dort, Abschnitt *Umsetzung (Implementer)*). Skripte, Hilfs-Dockerfile und Sonden liegen unter
`tools/messung/egress-dns-upstream/` (Ablage: [`docs/messungen/README.md`](../README.md)). Gemessen am 2026-10-19.

Gegenstand ist der Pfad, den die Regeln des erzeugten Startscripts `egress-init.sh` offen lassen: Die Zeilen zum Upstream-Resolver erlauben
`ip daddr <Upstream> udp/tcp dport 53` für **jeden** Prozess. Die Messung prüft eine **Messvariante** (nft-Regel nur für den `dnsmasq`, `meta skuid`)
und misst, was sie schließt und was offen bleibt. Die Variante ist ein Messskript (`variante.sh`), kein Produktcode; Template, Golden Case,
Handbuch und ADR sind unverändert.

## Messstand

**Anderer Messstand als im Slice und in den Vorgänger-Messungen:** Der Slice nennt Docker auf Colima; dieses Protokoll misst auf einem
**nativen Linux-x86_64-Host ohne VM**. Aussagen gelten nur dafür.

| | |
|---|---|
| Host | Linux 6.8.0-139-generic x86_64 (Ubuntu 24.04.5), 20 CPUs, cgroup v2, `kernel.apparmor_restrict_unprivileged_userns=0` |
| Docker | Client und Server 29.9.1 (linux/amd64), Speichertreiber overlay2, kein Colima |
| Datum, Läufe | 2026-10-19; **zwei vollständige Läufe** (`00-alles.sh`), 180 `AUSSAGE`-Zeilen, beide gleich (`vergleich.sh`, Exit 0); die Rohdaten dieses Verzeichnisses sind Lauf 1, `rohdaten-lauf2-aussagen.txt` trägt die Aussagen von Lauf 2 |
| u-boot | `u-boot version 0.12.0-dev`, `make build` (Runtime-Image), Stand des Commits `d2f8294`; die Messprojekte erzeugt `u-boot init --devcontainer --sandbox`, `egress.enabled: true`, `onUnavailable: warn`, ohne und mit `nestedRuntime: podman` |
| Image | Basis `mcr.microsoft.com/devcontainers/base:debian` (Debian 13 trixie), `dnsmasq` 2.91, `nftables` 1.1.3, `sudo` 1.9.16p2, Podman 6.1.4 (nur Image `pod`); `vscode` UID 1000, `nobody` UID 65534, `dnsmasq` UID 999 |
| Upstream-Liste | `1.1.1.1` und `8.8.8.8` (aus der `resolv.conf` des Containers, die Docker vom Host ableitet) |
| Container | wie die erzeugten `runArgs` (`rohdaten-02.txt`): `--cap-add=NET_ADMIN`, im Fall `pod` zusätzlich die Lockerungen für verschachteltes Podman; Prozess 1 `sleep infinity` unter `--init`; Standardnetz (Brücke), IPv4 ohne IPv6-Adresse |
| Netzzugriff | öffentliche Hosts der Default-Allowlist, `example.com`, `1.1.1.1`, `8.8.8.8` (DNS, wenige Pakete; ICMP an `1.1.1.1` je Sonde), Docker Hub für den Pull des Podman-Falls; das Gegenüber für `D3` und `D5` ist ein Container der Messung |

## Voraussetzungen am Host (nicht im Repo)

- Docker mit Zugriff auf `debian:stable` und `mcr.microsoft.com/devcontainers/base:debian` und das Netz, `make`, `bash`, `git`, `grep`, `diff`, `sort`,
  `wc`, `cmp`, `tr`, `cut`, `sed`, `id`; keine Go-Toolchain, kein Python. Aus dem Repo nutzt `02-projekte.sh` `make build IMAGE=…`.
- **Messordner außerhalb des Repos** (`WORK`): hier entstehen die Messprojekte `p-sbx` und `p-pod` und die Bau-Logs; `OUT` ist das Ausgabeverzeichnis
  der Rohdaten (Pflicht; für dieses Protokoll das Verzeichnis dieser Datei).
- Alle Docker-Objekte der Messung tragen das Präfix `u-boot-messung-dns` (Container, Images, das Netz `u-boot-messung-dns-net`); aufgeräumt wird nur nach
  diesem Präfix (`10-aufraeumen.sh`).

## Reihenfolge und Aufruf

`WORK=<messordner> OUT=<ausgabeverzeichnis> tools/messung/egress-dns-upstream/00-alles.sh` führt alle Schritte aus; einzeln:

| Schritt | Skript (mit Hilfsdateien) | Rohdaten | Frage |
|---|---|---|---|
| 1 | `01-messstand.sh` | `rohdaten-01-messstand.txt`, `rohdaten-11-vorher.txt` | Messstand; Bestand der Docker-Objekte vorher (nur Anzahlen und Präfix-Objekte) |
| 2 | `02-projekte.sh`, `Dockerfile.messung` | `rohdaten-02.txt`, `rohdaten-02-erzeugt-egress-init.txt` | Bau von u-boot, zwei Projekte, drei Images; das erzeugte `egress-init.sh` |
| 3 | `03-d1.sh` | `rohdaten-03-d1.txt` | D1: Benutzer von `dnsmasq` und `vscode`, Unterscheidbarkeit per `meta skuid` |
| 4 | `04-d2.sh`, `inner-d2.sh`, `variante.sh` | `rohdaten-04-d2.txt` | D2: direkter Pfad und erlaubte Namen, Ist gegen Variante |
| 5 | `05-d2-nested.sh`, `inner-nested.sh` | `rohdaten-05-d2-nested.txt` | D2 für verschachteltes Podman (roter Test), Ist gegen Variante |
| 6 | `06-d3.sh`, `inner-d3.sh`, `inner-resp.sh` | `rohdaten-06-d3.txt` | D3: Umgehungswege, Fälle ohne Startscript / Ist / Variante |
| 7 | `07-d3-netz.sh`, `inner-d3-netz.sh` | `rohdaten-07-d3-netz.txt` | D3 Zusatz: benutzerdefiniertes Netz mit IPv6, Docker-eigener Resolver |
| 8 | `08-d5.sh`, `inner-d5.sh`, `dnsq.pl` | `rohdaten-08-d5.txt` | D5: Datenkanal über erlaubte Namen |
| 9 | `09-d4.sh`, `inner-d4.sh` | `rohdaten-09-d4.txt` | D4: zweiter Lauf, Neustart, Neuaufbau, 50 Läufe |
| 9b | `09b-gegenprobe.sh` | `rohdaten-09b-gegenprobe.txt` | Gegenprobe: Variante mit falscher UID sperrt den `dnsmasq`; die Prüfung aus Schritt 9 muss anschlagen |
| 10 | `10-aufraeumen.sh` | `rohdaten-11.txt`, `rohdaten-11-nachher.txt` | Aufräumen nach Präfix |
| — | `vergleich.sh <lauf-1> <lauf-2>` | `rohdaten-lauf2-aussagen.txt` | vergleicht die `AUSSAGE`-Zeilen zweier Läufe |

## Die Messvariante

`variante.sh skuid` führt eine **Kopie** des erzeugten Startscripts aus, in der `sed` genau zwei Zeilen ändert (UDP- und TCP-Regel je Upstream); das Skript
bricht ab, wenn nicht genau zwei Regeln ersetzt wurden:

```
nft add rule ip u_boot_egress output meta skuid 65534 ip daddr "$ns" udp dport 53 accept
nft add rule ip u_boot_egress output meta skuid 65534 ip daddr "$ns" tcp dport 53 accept
```

`variante.sh ist` führt das erzeugte Script unverändert aus (Gegenprobe). Der Zusatz `zaehler` hängt **Messinstrumente** an: eine nicht terminierende
Zählerregel vor den Upstream-Regeln (Versuche anderer Benutzer als `nobody` an Port 53, nicht über `lo`) und eine Regel `counter drop` am Ende, die
Pakete an Port 53 zählt, die bis dahin nicht angenommen wurden. Beide ändern die Entscheidung nicht (die Richtlinie der Kette ist `drop`). Die Fälle
`D4` und `D5` laufen ohne `zaehler`.

## Ergebnisse (Kurzfassung; Einordnung im Slice)

**D1.** `dnsmasq` läuft nach dem Startscript als `nobody` (UID 65534; die Konfiguration trägt keine Zeile `user=`/`group=`), `vscode` hat UID 1000, ein
Benutzer `dnsmasq` (UID 999) existiert, trägt den Prozess aber nicht. Zwei Zähler per `meta skuid` unterscheiden beide: eine Anfrage von `vscode`
direkt an den Upstream zählt unter 1000 (1 Paket), die Anfrage des `dnsmasq` für einen erlaubten Namen unter 65534 (2 Pakete).

**D2 (ohne `nestedRuntime`).** Ist: `example.com` löst direkt am Upstream auf (UDP und TCP, beide Upstreams); mit der Variante läuft jede direkte
Anfrage in die Frist (4 von 4, dazu `api.github.com` direkt). Erlaubte Namen: alle sieben Namen der Default-Allowlist lösen über den lokalen Resolver
in beiden Fällen auf, `getent` und `https://api.github.com` / `https://deb.debian.org` (200) ebenso; `example.com` bleibt über den lokalen Resolver
`REFUSED`.

**D2 (verschachteltes Podman, roter Test).** Der Ablauf des Tests `TestE2E_SandboxDevcontainer_NestedPodman_BuildEgress` (Dockerfile mit
`RUN getent hosts deb.debian.org` und `RUN ! getent hosts example.com`, vorweg `RUN cat /etc/resolv.conf`): Ist Exit 1 (Schritt `RUN ! getent hosts
example.com` gibt zwei IPv6-Adressen aus; 7 s), Variante Exit 0 (43 s, in beiden Läufen). Der verschachtelte Container bekommt `nameserver 169.254.1.1`
(pasta), `8.8.8.8`, `8.8.4.4`; die Anfragen des verschachtelten Containers gehen mit der UID 1000 aus dem Sandbox-Container (Zähler `skuid 1000`):
Ist 1 Paket im Bau, davon 1 an `8.8.8.8`, keines verworfen; Variante 8 Pakete (je 4 an `8.8.8.8` und `8.8.4.4`), alle verworfen. Der Weg
über `169.254.1.1` (pasta zum lokalen `dnsmasq`) löst erlaubte Namen in beiden Fällen auf (`deb.debian.org`, `api.github.com`). In der Variante dauert
die Abfrage eines unlisteten Namens im verschachtelten Container etwa 36 s (`REFUSED` vom ersten Nameserver, danach laufen die verworfenen Rückfälle in die
Frist), im Ist-Zustand 20 ms. `podman run --dns 8.8.8.8` löst im Ist alles auf, in der Variante nichts mehr, auch keinen erlaubten Namen (20 s je Name);
`--dns 9.9.9.9` löst in beiden Fällen nichts auf; `--network=host` verhält sich wie der Sandbox-Container.

**D3.** Mit der Variante bleibt gesperrt (Kontrolle ohne Startscript: jeweils offen, außer wo unten vermerkt): UDP und TCP an Port 53 jedes Upstreams
und eines fremden Resolvers, TCP 853 und 443 und UDP 5353 an den Upstream, ICMP an Upstream und fremde Adresse, IPv6 (UDP/TCP 53, ICMP; im
Standardnetz hat der Container keine IPv6-Adresse außer `::1`, deshalb zusätzlich im benutzerdefinierten Netz mit IPv6: gesperrt, Kontrolle offen). Im Ist
sind davon nur UDP/TCP 53 an die Upstreams offen. **Offen bleibt in beiden Fällen**: jeder Port, jedes Protokoll (UDP/TCP 53, 5353, 853, 443, ICMP) an eine
Adresse **im Set `allowed`** (`ip daddr @allowed accept` hat keine Bedingung an Port oder Benutzer; gemessen mit einem Gegenüber, dessen Adresse die Messung in das
Set einträgt, so wie es der `dnsmasq` für einen erlaubten Namen täte; Nutzdaten kamen am Gegenüber an); und alles, was `sudo` erlaubt (`sudo dig` als `root` ist
mit der Variante gesperrt, `sudo setpriv --reuid=65534 dig` offen; die `sudoers`-Regel des Images ist `(root) NOPASSWD: ALL`, `sudo -u nobody` scheitert
daran). Zusatz im benutzerdefinierten Netz: Docker-eigener Resolver `127.0.0.11` ist mit dem erzeugten Script (Ist und Variante) **nicht** erreichbar
(Kontrolle ohne Script: offen); auch der lokale `dnsmasq` löst dort keinen erlaubten Namen auf (Ursache nicht untersucht; vermutet: die Anfragen des
Docker-Resolvers an den Host-Upstream gehen aus dem Namensraum des Containers und treffen die Richtlinie `drop`).

**D4.** Variante und Ist verhalten sich gleich: erster Lauf, zweiter Lauf im selben Container, Lauf nach `docker restart`, Lauf im neuen Container
aus demselben Image — je ein `dnsmasq` (UID 65534), Pid lebt, 7 `accept`-Regeln (davon 4 mit `meta skuid` in der Variante, keine Regelanhäufung),
beide Tabellen vorhanden, `api.github.com` lokal `NOERROR`, Statusdatei `result=ok`. 50 erneute Läufe hintereinander je Fall: 0 auffällige
Zustände, 0 Meldungen (Ist und Variante, beide Läufe). Nach dem `docker restart` hat Docker `/etc/resolv.conf` neu geschrieben (kein `127.0.0.1`),
das Startscript stellt es wieder her. Gegenprobe (`09b`): Mit `DNSUID=1234` löst der lokale `dnsmasq` einen erlaubten Namen nicht mehr auf (`status: REFUSED` statt `NOERROR`), die Prüfung aus Schritt 9 meldet den Zustand als auffällig (beide Läufe). Nicht erneut gemessen: das Fenster ohne Filter beim erneuten Lauf (die Variante ändert den Ablauf nicht, nur zwei Regeln).

**D5.** Der `dnsmasq` leitet Anfragen für beliebige Unternamen eines erlaubten Namens an den Upstream weiter (mit der Variante: alle `kanal1.api.github.com`-Anfragen
kamen am Upstream an, keine für `example.com`). Obergrenze im Laboraufbau (Gegenüber auf demselben Host als Upstream, sequenzielle Anfragen mit drei Labels zu je
63 Hex-Zeichen): 3 650 bis 5 640 Anfragen je Sekunde, 0,7 bis 1,1 Millionen Label-Zeichen je Sekunde (bei 5 Bit je Zeichen etwa 0,4 bis 0,7 MB/s), Ist wie
Variante; die direkte Abfrage am Upstream (nur im Ist offen) lag bei 9 300 und 11 600 Anfragen je Sekunde (je ein Lauf). Ein Empfänger braucht den autoritativen Server
des erlaubten Namens; für die Default-Allowlist (`github.com`, `githubusercontent.com`, `debian.org`) ist das nicht der Messende; die Anfragen sind für den
Betreiber des Upstream-Resolvers sichtbar.

## Lesehilfen zu den Rohdaten

- **`AUSSAGE`-Zeilen** sind die stabilen Ergebnisse (`vergleich.sh`); **`MESSWERT`** (Durchsatz) und Dauern sind zeitabhängig und gehen nicht in den Vergleich ein.
  Die zwei Läufe des Durchsatzes (Anfragen in 5 s, Ist/Variante): Lauf 1 18 253 / 20 426, Lauf 2 19 395 / 28 183.
- **Fall `ohne`** (D3) ist die Kontrolle ohne Startscript: Sie zeigt, dass eine Sonde einen offenen Pfad als „offen“ erkennt. Wo auch die Kontrolle „gesperrt“ zeigt,
  sagt die Sonde nichts: `udp5353-upstream` (am echten `1.1.1.1` lauscht nichts auf 5353), `ipv6-extern-udp53` im Standardnetz (keine IPv6-Route),
  `udp53-127.0.0.11` im Standardnetz (kein eingebauter Resolver), `ipv6-tcp853-R6` (der Listener antwortet nicht auf IPv6). Für IPv6 gilt deshalb `rohdaten-07-d3-netz.txt`.
- **Antwort-Container** (`inner-resp.sh`): zwei `dnsmasq` (Port 53 und 5353, beantworten jeden Namen mit `203.0.113.7`) und TCP-Listener 853 und 443; er ist
  Gegenüber für D3 (Adresse im Set per `sudo nft add element`, Messvorbereitung) und Upstream des Sandbox-Containers in D5 (`docker run --dns`).
- **Zähler** (`nft list chain`): `skuid 65534` zählt die Anfragen des `dnsmasq`, `skuid != 65534` die der anderen Benutzer; `counter … drop` am Ende zählt, was
  die Richtlinie verwirft. In `rohdaten-05-d2-nested.txt` zählen die Regeln mit `skuid 1000` auch `ip daddr 8.8.8.8` und `8.8.4.4`.
- **Pakete in der Zählung** sind Pakete, keine Anfragen: ein erneuter Versuch (Wiederholung, TCP-SYN) zählt erneut.
- **Bestand** (`rohdaten-11-*.txt`): Die Liste vorher von Lauf 1 enthält vier Images mit dem Präfix aus den Vorversuchen des Messenden (nach dem Aufräumen 0);
  die Anzahlen fremder Objekte ändern sich durch andere Arbeit auf dem Host (Container 3 → 4 in Lauf 2 nachher, 3 danach), Namen stehen nicht in den Rohdaten.
- Zeitstempel sind UTC.

## Nicht gemessen

Colima, Docker Desktop, Linux hinter einer VM, `arm64`, die CI und der GitHub-Runner; `nestedRuntime: docker` (der verschachtelte `dockerd` läuft als `root`; sein
Verkehr läuft über die Weiterleitung, nicht über die Kette `output`); Podman mit anderem Netzwerkmodus als `pasta` (Standard) und `host`; eine IPv6-fähige Upstream-Liste;
echte Namen mit kontrollierbarem autoritativen Server (D5 misst nur bis zum Upstream); Last des Hosts; das Fenster ohne Filter des erneuten Laufs (D4); eine
Umsetzung im Template (keine; Golden Case, Docker-Test und Handbuch sind unberührt); die Wirkung auf den Go-Test selbst (`make test-docker` lief nicht, der
Ablauf des Tests ist im Wegwerf-Container nachgestellt).
