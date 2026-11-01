# Messprotokoll: wo entsteht die Zeile `unable to find user <Benutzer>: …`, und wie sieht sie aus, wenn die Sitzung sie selbst schreibt?

Messung zum Slice [`slice-v1-shell-erkennung-fehlender-benutzer-eng`](../../plan/planning/done/slice-v1-shell-erkennung-fehlender-benutzer-eng.md) ([`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), DoD-Punkt 1). Gemessen ist nur das Verhalten des `docker`-Clients und der Engine bei `docker exec -u <Benutzer> …`; der Launcher von `u-boot` lief nicht.

## Messstand

| Eigenschaft | Wert |
|---|---|
| Datum | 2026-11-01 |
| Engine | **29.9.2** (API 1.56), `linux/amd64`, lokale Engine |
| Client | 29.9.2 (API 1.56), Host-Binary; derselbe Client läuft in den Fällen der Gruppe E auch im Sitzungs-Container (Binary des Hosts, nur lesend eingehängt) |
| Host | Linux 7.0.0-42-generic `x86_64` |
| Images | Sitzungs-Container und gestoppter Container: `mcr.microsoft.com/devcontainers/base:debian` (`sleep 600`; Benutzer `root` und `vscode` vorhanden, `nouser` fehlt); anderer Container: `alpine:3` (`sleep 600`; `vscode` und `nouser` fehlen) |
| Läufe | **zwei Skriptläufe**, in jedem läuft jede Zelle **zweimal**: vier Beobachtungen je Zelle |
| Durchführung | Hauptagent (Rolle Implementer), nur Docker am Host |

## Skript und Aufruf

[`tools/messung/shell-benutzer-fehlt/messen.sh`](../../../tools/messung/shell-benutzer-fehlt/messen.sh), aus der Repo-Wurzel:

```bash
bash tools/messung/shell-benutzer-fehlt/messen.sh > docs/messungen/slice-v1-shell-erkennung-fehlender-benutzer-eng/rohdaten-lauf1.txt 2>&1
bash tools/messung/shell-benutzer-fehlt/messen.sh > docs/messungen/slice-v1-shell-erkennung-fehlender-benutzer-eng/rohdaten-lauf2.txt 2>&1
```

Die beiden Rohdateien stammen aus dem Skript in der abgelegten Fassung (Lauf 1 um 2026-11-01T06:45:55Z, Lauf 2 um 06:46:05Z). Zeitlimit je Docker-Aufruf `timeout -k 5 30`; aufgeräumt wird nach Label `ubmess-shell=1`, nach beiden Läufen blieb kein Container. Lockerung: Der Sitzungs-Container bekommt den Docker-Socket des Hosts und das `docker`-Binary des Hosts (nur lesend), damit die Sitzung wirklich `docker exec` gegen den anderen Container ausführt. Ein früherer Versuch mit einer Zelle (Terminal-Sitzung ohne Eingabe) lief in das Zeitlimit; die Zelle wurde mit Eingabe wiederholt, die Rohdaten stammen aus den Läufen der abgelegten Fassung.

Rohdaten: [`rohdaten-lauf1.txt`](rohdaten-lauf1.txt), [`rohdaten-lauf2.txt`](rohdaten-lauf2.txt). Jede Zelle nennt Befehl, Exit-Status, Byte-Zahl je Strom und beide Ströme getrennt.

## Beobachtung

Die Zellen sind in beiden Läufen und in beiden Wiederholungen gleich ausgefallen.

| Gruppe | Zelle | Aufruf | Exit-Status | Standardausgabe | Fehlerausgabe |
|---|---|---|---|---|---|
| A, Benutzer `nouser` fehlt, Container läuft, kein Terminal | A1 bis A4 | `docker exec [-i] -u nouser -w / <c> bash [-c true]`, Standardeingabe `/dev/null` oder Pipe | **1** | leer | `Error response from daemon: unable to find user nouser: no matching entries in passwd file` (91 Byte) |
| B, mit Terminal (`script`) | B1, B2 | `-i -t` / `-t` | 1 | dieselbe Zeile, mit `\r\n`, im Terminal-Strom | (im Terminal-Strom) |
| B, mit Terminal, Fehlerausgabe des Clients getrennt | B3 | `-i -t`, `2>&3` | 1 | leer | dieselbe Zeile (91 Byte) |
| C, Container gestoppt | C1 bis C3 | Benutzer fehlend oder `root` | 1 | leer | `Error response from daemon: container <ID> is not running`; die Meldung zum gestoppten Container kommt vor der zum Benutzer, auch bei fehlendem Benutzer |
| D, Kontrolle: Benutzer vorhanden | D1, D2 | `-u root` | 0, bei `exit 3` der Shell 3 | leer | leer |
| E, die Sitzung erzeugt die Zeile selbst (echtes `docker exec -u nouser <anderer>` im Container, Sitzung als `root`) | E1, E2 | ohne Terminal, innerer Aufruf ohne und mit `-i` | 1 (Status des inneren Aufrufs, letzter Befehl der Shell) | leer | dieselbe Zeile mit Präfix `Error response from daemon: ` (91 Byte) |
| E, mit Terminal | E3, E3b | Sitzung unter `script` | 1 | die Zeile steht im Terminal-Strom der Sitzung | **leer** (E3b: der äußere Client schrieb nichts auf seine Fehlerausgabe) |
| E, Kontrolle: Benutzer des inneren Aufrufs vorhanden | E4 | `docker exec -u root <anderer> …` | 0 | leer | leer |
| E, **derselbe Name** in Sitzung und innerem Aufruf (`vscode`: im Sitzungs-Container vorhanden, im anderen fehlend) | E5 | Sitzung als `vscode`, innerer Aufruf `-u vscode`, ohne Terminal | 1 | leer | `Error response from daemon: unable to find user vscode: no matching entries in passwd file` |
| dasselbe, Shell endet danach mit `exit 0` | E6 | wie E5 | 0 | leer | dieselbe Zeile (91 Byte) |
| dasselbe, innerer Aufruf mit `2>&1` | E7 | wie E5 | 1 | die Zeile mit Präfix (91 Byte) | leer |
| F, nachgestellt (von Hand geschrieben, **kein Erzeuger**) | F1 | Sitzung `echo "unable to find user nouser: …"; exit 126` | 126 | die Zeile **ohne** Präfix (63 Byte) | leer |
| G, Probe wie vor der Sitzung, beide Ströme zusammen | G1, G2 | Benutzer vorhanden / fehlend | 0 / 1 | G2: die Zeile mit Präfix | — |

Die Kontrollzeilen vor den Zellen (im Kopf jeder Rohdatei) bestätigen: `nouser` fehlt im Sitzungs-Container, `vscode` ist dort vorhanden und fehlt im anderen Container, `vscode` erreicht den Docker-Socket.

## Auswertung

**1. Wo die Zeile entsteht (gemessen, Engine 29.9.2).** Die Engine erzeugt die Zeile; der Client reicht sie als `Error response from daemon: unable to find user <Benutzer>: no matching entries in passwd file` auf der **Fehlerausgabe** durch, mit **Exit-Status 1**, bei laufendem Container, mit und ohne `-i`, mit und ohne Standardeingabe, und mit Terminal (B3: die Fehlerausgabe des Clients bleibt getrennt von der Standardausgabe). Es gibt in keiner Zelle eine Form auf der Standardausgabe, und nirgends Status 126. Bei gestopptem Container steht stattdessen die Meldung zum Container (C).

**2. Wie die Zeile aussieht, wenn die Sitzung sie selbst schreibt (gemessen, Engine 29.9.2).** Ein `docker exec -u <fehlender Benutzer> <anderer Container>` im Container, mit dem Client 29.9.2 gegen dieselbe Engine, erzeugt **dieselbe Zeile in derselben Form**: Präfix `Error response from daemon: `, auf der Fehlerausgabe des inneren Clients, der in einer Sitzung ohne Terminal die Fehlerausgabe des äußeren Clients ist (E1, E2, E5). Der Status des äußeren Befehls ist der Status des letzten Befehls der Shell: 1 (E1, E2, E5) oder 0 nach `exit 0` (E6). Mit Terminal landet die Zeile im Terminal-Strom (E3, E3b).

**3. Kein Unterschied in der Form zwischen „der Client meldet den fehlenden Benutzer der Sitzung“ und „die Sitzung meldet einen fehlenden Benutzer“.** Der Fall A (Client, Benutzer der Sitzung fehlt) und der Fall E5 (Sitzung, innerer Aufruf mit demselben Namen) liefern auf der Engine 29.9.2 dieselbe Zeile auf demselben Strom mit demselben Status 1. Die Unterscheidung bleibt allein der Probe vor der Sitzung (D1: der Benutzer ist vorhanden) und einer Probe danach (G1 / G2) überlassen. Eine erneute Probe unterscheidet auch hier: G1 endet mit 0, G2 mit 1 und der Zeile. Ein Benutzer, der erst nach dem Beginn der Sitzung verschwindet, ist in dieser Messung nicht erzeugt worden.

**4. Die Form „Zeile ohne Präfix auf der Standardausgabe mit Status 126“** (die Form, die `userMissingOnStdout` erkennt) **ist auf der Engine 29.9.2 nicht aufgetreten.** Sie steht in der Spezifikation für die Engine 28.0.4 und stammt aus einer früheren Messung auf der Runner-Engine 28.0.4 ([`zusagen-und-messstaende.md`](../../maintainer/zusagen-und-messstaende.md), Abschnitt 3.15), für die es kein Protokoll unter `docs/messungen/` gibt. F1 ist keine Beobachtung eines Erzeugers: Sie zeigt nur, dass ein von Hand geschriebener Text dieser Form durch den Client unverändert, mit Status 126, auf der Standardausgabe ankommt.

## Urteil zur Kernfrage

*Erzeugt der Client die Zeile bei fehlendem Benutzer auf der Standardausgabe, wie `userMissingOnStdout` annimmt, oder auf der Fehlerausgabe?*

**Auf der Engine 29.9.2 auf der Fehlerausgabe, mit Präfix und Status 1; die Annahme der Standardausgabe mit Status 126 trifft für diese Engine nicht zu.** Für die Engine 28.0.4, auf die sich `userMissingOnStdout` bezieht, liefert diese Messung keinen Befund: Sie stand nicht zur Verfügung. Daraus folgt:

- Die Frage, ob der **dritte Aufruf** (`userMissingNow`) für die Form von 28.0.4 begründet ist, lässt sich mit dieser Messung **nicht entscheiden**: Ob ein `docker exec -u <Benutzer> <anderer Container>` in der Sitzung auf 28.0.4 die Zeile auf die Standardausgabe schreibt (wie der Client selbst dort), ist nicht gemessen. Auf 29.9.2 kann der dritte Aufruf nie greifen, weil `userMissingOnStdout` den Status 126 und die Standardausgabe verlangt; auf der Engine 29.9.2 ist der Pfad inaktiv.
- **Befund auf der Engine, die gemessen ist (29.9.2):** Der Fall „die Sitzung schreibt die Zeile selbst“ läuft über die **Fehlerausgabe-Erkennung** (`clientMessageOnStderr`), und der dritte Aufruf schützt sie nicht. Nach dem Lesen von `internal/adapter/driven/docker/devcontainer_shell.go`: Eine Sitzung ohne Terminal, deren letzter Befehl `docker exec -u <Benutzer der Sitzung> <anderer Container>` mit fehlendem Benutzer ist (E5: Status 1, Standardausgabe leer, Zeile mit Präfix und dem Namen des Benutzers der Sitzung), erfüllt dort alle Bedingungen — Präfix `Error response from daemon: `, Name in der Zeile, Standardausgabe leer, Status ungleich 0 — und gilt als Fehler des Clients; der Befehl endete mit `12` statt `0`. Das ist **kein gefahrener Lauf des Launchers**, sondern eine Lesung des Codes gegen die gemessene Form E5. Der Pfad ist in der Spezifikation für 29.9.2 so beschrieben (Punkt „Grenzen der Erkennung“, erster Aufzählungspunkt), ohne die Ausnahme „die Sitzung schreibt die Zeile selbst“.
- Der Befund gilt nur für den Namen des Benutzers der Sitzung. Eine Zeile über einen anderen Benutzer (E1 mit `nouser`, Sitzung als `root`) erfüllt „nennt die Kennung oder den Benutzer der Sitzung“ nicht und gilt als Sitzung.

Ob Code und die Folge-Entscheidung ([`ADR-0039`](../../plan/adr/0039-devcontainer-shell-dritter-aufruf-bei-fehlendem-benutzer.md)) daraufhin angepasst werden, entscheidet der Auftraggeber; diese Messung ändert beides nicht.

## Nicht gemessen

- **Engine 28.0.4** (keine Engine dieser Version verfügbar): die Form „Standardausgabe, ohne Präfix, Status 126“ und damit der Erzeuger der Zeile in der Sitzung auf dieser Engine. Alle Aussagen oben gelten für 29.9.2.
- Andere Clients als 29.9.2, andere Hosts als Linux `amd64`, nicht lokale Engines (Colima).
- Der Launcher von `u-boot` selbst (`make test-docker` lief nicht); die Folgerung zu `clientMessageOnStderr` ist eine Lesung des Codes.
- Ein Benutzer, der nach dem Beginn der Sitzung im Container verschwindet (`userdel` in der Sitzung).
- Eine Shell, die kein `bash` ist (der Container der Messung führt `bash`; die Probe fällt sonst auf `sh`).
- Mehr als vier Beobachtungen je Zelle; die Zellen streuten nicht.
