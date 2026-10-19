# Slice v1: Egress — Pfad zum Upstream-Resolver für Prozesse ohne `sudo` schließen oder messen (`LH-FA-DEV-008`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Der Slice ist ein eigener Slice **nach** [`welle-sandbox-egress-ollama`](../done/welle-sandbox-egress-ollama.md) (Entscheidung des Projektinhabers, 2026-10-18) und hält deren Closure und das Release nicht auf; er hat keine Closure-Bedingung, die von seiner DoD verschieden wäre. Reihenfolge: nach der Closure jener Welle.

**Bezug:** [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Egress-Restriktion: „Allowlist für ausgehenden Netzwerkverkehr“; Zielbild, Code folgt), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil: Schadensbegrenzung); [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) (Mechanismus `dnsmasq` + nftables-Set, `Accepted`), [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md) (Messbefund M1 a, `Accepted`).
Anlass: Befund M1 (a) aus [`slice-v1-sandbox-egress-neuladen-entwurf`](../done/slice-v1-sandbox-egress-neuladen-entwurf.md); Entscheidung des Projektinhabers (2026-10-18): Code zum Ziel, eigener Folge-Slice, das Lastenheft wird nicht abgeschwächt.

**Berührte Spec-Stellen:** [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) *Grenze* (nur der Teil zur Schreibbarkeit der Statusdatei; dessen Nachführung gehört dem Spezifikations-Slice der Welle, nicht diesem Slice), [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) (Allowlist; der Architect liest, ob ein DNS-Satz dort steht). Der Slice **liest** sie und ändert keine; ist nach der Messung eine Spec-Änderung nötig, benennt der Architect die Stelle, und ein eigener Spezifikations-Slice folgt dem Code.

**Verantwortlich:** pt9912

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Prozess im Sandbox-Container ohne `sudo` bekommt über den Upstream-Resolver keine Adresse mehr für einen Namen außerhalb der Allowlist — oder die Messung zeigt, dass und warum das mit dem Mechanismus von [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) nicht geht, und das Ergebnis ist ausgewiesen. **Messung zuerst**, Umsetzung nur nach dem Messbefund.

### Gefundener Stand (gelesen am 2026-10-18)

- **Die Regeln lassen DNS zu jedem Upstream für alle Prozesse zu.** `internal/hexagon/application/templates/devcontainer/egress-init.sh.tmpl`: `upstreams` kommt aus den IPv4-`nameserver`-Zeilen der gesicherten `/etc/resolv.conf.u-boot-orig` (Zeile 36); die Schleife in den Zeilen 63 bis 66 legt je Upstream `ip daddr "$ns" udp dport 53 accept` und `… tcp dport 53 accept` in die Kette `output` — ohne Bedingung an Benutzer oder Prozess. Der lokale `dnsmasq` fragt nur die Upstreams für erlaubte Namen (`server=/$h/$ns`, Zeile 52; `no-resolv`, Zeile 47); der Container zeigt über `/etc/resolv.conf` auf `127.0.0.1` (Zeile 72). Ein Prozess, der den Upstream direkt anspricht, umgeht den lokalen Server.
- **Gemessen** (Messprotokoll [`rohdaten-02-m1.txt`](../../../messungen/slice-v1-sandbox-egress-neuladen-entwurf/rohdaten-02-m1.txt), Abschnitt „DNS direkt an den Upstream-Resolver 192.168.5.1 (ohne sudo)“): `example.org` löst mit zwei A-Einträgen auf; die Verbindung zur Adresse läuft in den Timeout (`Connection timed out after 5006 milliseconds`). Je ein Lauf, Messstand Docker 29.5.2 auf Colima 0.10.3 (arm64). Ein Datenkanal über DNS ist nicht gemessen.
- **Beleg für verschachtelte Container** (Messung des Cut-Slice [`slice-v1-release-cut-v0.12.0`](../done/slice-v1-release-cut-v0.12.0.md), Abschnitt *Umsetzung (Implementer)*, *Befund zum roten Test*, 2026-10-19, Linux `amd64`, native Docker Engine 29.9.1, zwei Läufe): Der Docker-Test `TestE2E_SandboxDevcontainer_NestedPodman_BuildEgress` (`internal/e2e`) ist rot: In `STEP 3/3: RUN ! getent hosts example.com` löst der Name auf (nur IPv6), weil der verschachtelte Podman-Container als Rückfall `8.8.8.8` fragt, ein Upstream des Docker-Hosts, den die Regeln aus Zeilen 63 bis 66 für jeden Prozess zulassen; die TCP-Verbindung bleibt gesperrt; mit `--dns 9.9.9.9` ist der Test grün. Dieselbe Ursache wie oben, gelesen und gemessen, nicht für Colima oder die CI.
- **Was heute dazu gesagt wird:** [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) Alternative C2 („nicht erlaubte Namen sind nicht auflösbar“) gilt für den lokalen `dnsmasq`; das Handbuch (Abschnitt *Internet-Beschränkung (Egress)*, Absatz *Grenze der Namensauflösung*) weist die direkte Abfrage aus. [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) bleibt unverändert: Der Guardrail-Satz beschreibt die Aufhebung durch Capability, der offene DNS-Pfad ist eine Lücke des Codes gegenüber dem Zielbild, keine Abschwächung des Lastenhefts.

### Messfragen (Arbeitsplan; Docker auf Colima, Container wie im Messprotokoll der Vorgänger-Messung; der Architect entscheidet den Weg)

| # | Frage | Warum |
|---|---|---|
| D1 | Läuft der `dnsmasq` unter einem Benutzer, den eine nftables-Bedingung (`meta skuid`) von `vscode` unterscheidet? Welche Benutzer-IDs tragen `dnsmasq` und der Nutzer im gemessenen Image? | Voraussetzung für eine Regel „DNS zum Upstream nur für den `dnsmasq`“; der Benutzername des `dnsmasq` ist hier nicht gelesen, nur gemessen werden kann er |
| D2 | Mit einer solchen Regel (Upstream-DNS nur für den `dnsmasq`, andere Prozesse verworfen): löst ein Prozess ohne `sudo` einen Namen außerhalb der Liste direkt noch auf? Lösen erlaubte Namen weiter auf (Wirkung auf `api.github.com` und die Default-Allowlist)? | Wirksamkeit und Nebenwirkung |
| D3 | Wege um die Regel: TCP statt UDP, ein anderer Upstream (IPv6 ist gesperrt, Zeilen 67 bis 69), DNS über eine erlaubte Adresse, eine andere Portnummer; was bleibt offen? | Grenze der Lösung ehrlich benennen |
| D4 | Bleibt das Startscript beim zweiten Lauf und beim Neuaufbau wie bisher lauffähig (Messvariante aus dem Vorgänger-Messprotokoll)? | keine Regression der Start-Eigenschaften |
| D5 | Offen für den Architect: ein Datenkanal über die erlaubten Namen (Anfragen an erlaubte Namen werden an den Upstream weitergeleitet) — vorhanden, ja oder nein, mit welchem Durchsatz? | nur wenn D2 den direkten Pfad schließt; sonst Grenze im Handbuch |

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Entziehen von `sudo`** (Option S aus [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md)) — ändert Benutzerrechte und alle Startscripte; mit `sudo` bleibt die Restriktion ein Guardrail ([`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)), das Ziel hier ist nur der Weg *ohne* `sudo`.
- **Neu-Laden der Liste im laufenden Container** — entschieden gegen (R0, [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md)).
- **Eine Restriktion außerhalb des Containers** (Proxy, Netz-/DNS-Ebene) — eigener Vorgang ([`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) Alternative D).
- **Die Nachführung von [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) *Grenze* (Schreibbarkeit der Statusdatei)** — Spezifikations-Slice der Welle ([`slice-v2-spec-sandbox-egress-ollama`](../done/slice-v2-spec-sandbox-egress-ollama.md)).
- **Änderung des Lastenhefts** — nach dem Stand oben nicht erwartet (Zielbild, Code folgt); zeigt D3, dass der Pfad grundsätzlich nicht schließbar ist, geht der Befund an den Projektinhaber, und das Handbuch bleibt bei der ausgewiesenen Grenze.
- **Docker Desktop, Linux-Host ohne VM, `amd64`** — nicht gemessen, nicht zugesagt; jede Aussage nennt den Messstand.

## 2. Definition of Done

- [x] **Messung D1 bis D4** (D5 falls D2 den Pfad schließt) gemessen, mit Messstand, Skripten und Rohdaten im Repo (Ablage wie das Messprotokoll des Vorgängers); je Frage ein Ergebnis oder die Angabe „nicht gemessen“ mit Grund. Beleg: `eab3c2a` und `445fa99` (Skripte, Rohdaten, [Protokoll](../../../messungen/slice-v1-egress-dns-upstream-schliessen/README.md)), Ergebnis je Frage in §3 *Umsetzung (Implementer)*; Messstand ein nativer Linux-x86_64-Host (nicht Colima).
- [x] **Entscheidungs-ADR und Umsetzung:** (Entscheidungs-ADR: [`ADR-0026`](../../adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md), `Accepted` (2026-10-19); Umsetzung: `f691fab`, `6d85215`, `c1e6c07`, Belege in §3 *Umsetzung (Implementer)* Phase 2.) Eine Folge-ADR zu [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) (`Accepted` bleibt unverändert) trägt den Messbefund und den gewählten Weg; bei einem Weg, der schließt: Änderung am Startscript-Template, Golden Case und Docker-Test (`//go:build docker`), der den direkten Upstream-Pfad ohne `sudo` abweist und die erlaubten Namen weiter auflöst; Handbuch (Absatz *Grenze der Namensauflösung*) und `CHANGELOG.md` auf den Ist-Stand. Bei einem Weg, der nicht schließt: die Grenze steht ausgewiesen im Handbuch, und der Projektinhaber hat entschieden.
- [ ] **Der Test `TestE2E_SandboxDevcontainer_NestedPodman_BuildEgress` ist abgenommen:** Er wird mit dieser Änderung grün oder wird mit Begründung angepasst (Ausgang in §7); bis dahin bleibt `make test-docker` auf einem Host, dessen Upstream-Resolver ein Rückfall-Nameserver von Podman ist, rot. (Beleg des Implementers: grün und unverändert gefahren, §3 Phase 2; Abnahme durch den Verifier steht aus, der Ausgang in §7 folgt bei der Closure.)
- [ ] **Die Spec-Stelle**, die nach dem Befund nachzuführen ist (falls eine), ist benannt; die Nachführung selbst entsteht als eigener Spezifikations-Slice nach dem Code.
- [x] `make gates` grün; `make test-docker` grün, wenn Code geändert wurde (Belege: §3 Phase 2).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: Handbuch, ADR-Index, `CHANGELOG.md`, soweit der öffentliche Vertrag berührt ist (Handbuch, READMEs und `CHANGELOG.md`: `c1e6c07`; der ADR-Index trägt [`ADR-0026`](../../adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) seit dem ADR-Commit).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo mit Wellen von der nächsten Welle-Closure geprüft.

## 3. Plan (vor Code)

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; der Implementer ergänzt sie in seinem ersten Lauf und gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab (`AGENTS.md` §3.10). Jede Zeile ist eine Erwartung, keine gelesene Aussage über unveränderte Dateien.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Messprotokoll unter `docs/messungen/slice-v1-egress-dns-upstream-schliessen/` (Rohdaten, README mit Messstand und Voraussetzungen) und Skripte unter `tools/messung/egress-dns-upstream/` (24 Dateien: `lib.sh`, `00-` bis `10-` mit `09b-gegenprobe.sh`, `inner-*.sh`, `variante.sh`, `vergleich.sh`, `dnsq.pl`, `Dockerfile.messung`) | neu (`eab3c2a`, `445fa99`) | D1 bis D5; im Repo statt außerhalb ([`BEO-DOK/messprotokoll-ausserhalb-des-repos`](../observations/BEO-DOK/messprotokoll-ausserhalb-des-repos/observation.md)) |
| `docs/messungen/README.md` | update (`eab3c2a`) | Tabelle *Bestand*: eine Zeile für das neue Protokoll (gelesen: die Tabelle führt je Protokoll Verzeichnis und Skriptpfad) |
| `docs/plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md` und `docs/plan/adr/README.md` | neu / update | Folge-ADR zu [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) mit Messbefund und Weg: [`ADR-0026`](../../adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md), `Accepted` (2026-10-19; Weg A, Freigabe-Punkte a bis g) |
| `internal/hexagon/application/templates/devcontainer/egress-init.sh.tmpl` | update (`f691fab`) | `dns_user`/`dns_uid`, `user=` in der Konfiguration, `meta skuid` in den zwei Upstream-Regeln, zwei Abbruchzweige (Exit 11, `failed`) |
| `internal/hexagon/application/generate_sandbox_test.go` (Test des erzeugten Scripts; gelesen: es gibt keine Golden-Dateien zu `egress-init.sh`, nur Teilstring-Tests) | update (`f691fab`) | neuer Test `TestGenerateDevcontainer_SandboxEgress_UpstreamDNSOnlyForTheResolverUser`; `TestGenerateDevcontainer_SandboxStartScriptsWriteStatusFiles` zählt für `egress` vier `failed`-Zweige statt einem (Fix-Runde: Rückgabewert des `dnsmasq`-Starts geprüft) |
| `harness/replay.md` | update (`6d85215`) | zwei Zeilen der Replay-Tabelle nennen die neuen Zusagen und ihre Evidence-Tests |
| `internal/e2e/sandbox_devcontainer_docker_test.go` (Build-Tag `docker`) | update (`f691fab`) | neuer Test `TestE2E_SandboxDevcontainer_EgressUpstreamDNS`; `TestE2E_SandboxDevcontainer_NestedPodman_BuildEgress` (roter Test, Abnahme in der DoD) bleibt unverändert und lief grün |
| `internal/hexagon/application/devcontainer_fallback.go` (`egressProbe`), `internal/hexagon/domain/sandbox_status.go` (Kommentar zu `SandboxStatusFailed`) | update (`f691fab`, Fix-Runde) | die Doctor-Meldung für `failed` nennt die drei Abbruchgründe, der Hinweis beide Wege (Capability oder Meldung des Startscripts im Container-Log) |
| `internal/hexagon/application/devcontainer_fallback_test.go` | update (Fix-Runde) | neuer Test `TestDoctor_SandboxFallback_EgressFailedNamesEveryCause` pinnt Text und Hinweis |
| `docs/messungen/slice-v1-egress-dns-upstream-schliessen/README.md` | update (Fix-Runde) | Abschnitt *Nicht gemessen*: Zustandsmarker zu Template und `make test-docker` auf den Ist-Stand |
| `docs/reviews/2026-10-19-slice-v1-egress-dns-upstream-schliessen.md` | neu (Reviewer, `3ebae0b`) | Review-Report |
| `docs/reviews/2026-10-19-verification-slice-v1-egress-dns-upstream-schliessen.md` | neu (Verifier, `012abab`) | Verification-Report |
| `docs/plan/planning/open/slice-v1-egress-init-ohne-luecke.md` | update (Verifier-Auflage V4) | Zeilenbezüge auf das Template durch Symbolbezüge ersetzt, sonst unverändert |
| `docs/user/benutzerhandbuch.md`, `CHANGELOG.md`, `README.md`, `README.de.md` | update (`c1e6c07`) | Absatz *Namensauflösung und ihre Grenze* (vormals *Grenze der Namensauflösung*), verschachtelte Container, Statusdatei-Tabelle, zwei Verweise; README-Satz zum verschachtelten `podman build`; CHANGELOG `Changed` |
| diese Slice-Datei, `open/README.md`, Roadmap bei Start und Closure | update; Moves | Lifecycle (Phase 2: nur diese Datei; `open/README.md` und Roadmap unverändert, nicht gelesen auf Berührung) |
| `docs/reviews/` | neu (Reviewer, Verifier) | Review und Verifikation |

### Umsetzung (Implementer, 2026-10-19; Phase 1: nur die Messung; nichts gepusht, kein Tag)

**Auftrag.** Entscheidung des Projektinhabers (2026-10-19): Messung zuerst, Umsetzung und Folge-ADR nach dem Befund. Diese Phase **setzt nichts um**: kein Edit am Template `egress-init.sh.tmpl`, keine ADR, kein Golden Case, kein Docker-Test, kein Handbuch- und kein CHANGELOG-Eintrag zu einer Lösung, kein Eingriff in Lastenheft und Spezifikation. Die Messvariante ist ein Messskript (`tools/messung/egress-dns-upstream/variante.sh`), das im Wegwerf-Container eine `sed`-Kopie des erzeugten Startscripts ausführt. Die Entscheidung über Weg, Folge-ADR zu [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) und Umsetzung liegt beim Projektinhaber.

Commits: `eab3c2a` (Skripte, Rohdaten, Protokoll, Bestandszeile), `445fa99` (Gegenprobe `09b`: Skript, Rohdaten, Aussagen von Lauf 2); der Stand dieses Abschnitts folgt als eigener Commit. Protokoll und Rohdaten: [`docs/messungen/slice-v1-egress-dns-upstream-schliessen/README.md`](../../../messungen/slice-v1-egress-dns-upstream-schliessen/README.md).

**Messstand — ein anderer als im Slice vorgesehen.** Der Slice nennt Docker auf Colima mit `NET_ADMIN`. Gemessen wurde auf einem **nativen Linux-x86_64-Host ohne VM** (Linux 6.8.0-139, Docker 29.9.1, 2026-10-19, `u-boot 0.12.0-dev` am Commit `d2f8294`); die Rückführung „Messumgebung nicht verfügbar“ (§4) greift nicht, weil Container mit `NET_ADMIN` laufen. **Zwei vollständige Läufe**, 180 `AUSSAGE`-Zeilen, beide gleich (`vergleich.sh`, Exit 0); Durchsatz und Dauern sind Messwerte und stehen außerhalb des Vergleichs. Alle Aussagen gelten nur für diesen Messstand; Colima/arm64, Docker Desktop, die CI und der GitHub-Runner sind nicht gemessen.

**Ergebnis je Frage** (Messwert, nicht Auslegung; Rohdaten je Zeile im Protokoll):

| # | Ergebnis |
|---|---|
| D1 | **Gemessen: ja.** `dnsmasq` läuft nach dem Startscript als `nobody` (UID 65534; die erzeugte Konfiguration trägt keine Zeile `user=`), `vscode` hat UID 1000, ein Benutzer `dnsmasq` (UID 999) existiert, trägt den Prozess aber nicht. `meta skuid` unterscheidet beide: ein Zähler zählt die direkte Anfrage von `vscode` unter 1000 (1 Paket) und die Anfrage des `dnsmasq` unter 65534 (2 Pakete). Der Benutzername des `dnsmasq` ist damit ein **gemessener**, kein gepinnter Wert (siehe Befund 5). |
| D2 | **Gemessen.** *Ohne `nestedRuntime`:* im Ist-Zustand löst `example.com` direkt an beiden Upstreams (`1.1.1.1`, `8.8.8.8`) per UDP und TCP auf; mit der Variante laufen alle vier direkten Anfragen (und `api.github.com` direkt) in die Frist. Alle sieben Namen der Default-Allowlist lösen über den lokalen Resolver in beiden Fällen weiter auf; `getent` und HTTPS zu `api.github.com` und `deb.debian.org` (200) unverändert; `example.com` über den lokalen Resolver `REFUSED`. *Verschachteltes Podman (roter Test):* Der nachgestellte Ablauf des Tests gibt im Ist-Zustand Exit 1 (7 s, `RUN ! getent hosts example.com` gibt zwei IPv6-Adressen aus, wie im Cut-Befund), mit der Variante Exit 0 (43 s), in beiden Läufen. Der verschachtelte Container fragt **als `vscode`** (Zähler `skuid 1000`) über die Nameserver `169.254.1.1` (pasta), `8.8.8.8`, `8.8.4.4`: Ist 1 Paket im Bau (1 an `8.8.8.8`), Variante 8 Pakete (je 4 an `8.8.8.8` und `8.8.4.4`), alle verworfen. Der Weg über `169.254.1.1` (pasta zum lokalen `dnsmasq`) löst erlaubte Namen in beiden Fällen auf; die Abfrage eines unlisteten Namens braucht mit der Variante etwa 36 s (Ist 20 ms). |
| D3 | **Gemessen** (Kontrolle ohne Startscript je Sonde: offen). *Mit der Variante gesperrt:* UDP und TCP Port 53 an beide Upstreams und an einen fremden Resolver, TCP 853 und 443 und UDP 5353 am Upstream, ICMP an Upstream und fremde Adresse, IPv6 (im benutzerdefinierten Netz mit IPv6: UDP/TCP 53, ICMP gesperrt, Kontrolle offen). Im Ist-Zustand ist davon nur UDP/TCP 53 an die Upstreams offen. *Offen in Ist und Variante:* jeder Port und jedes Protokoll (53, 5353, 853, 443, ICMP, Nutzdaten kamen an) an eine Adresse **im Set `allowed`** (`ip daddr @allowed accept` ohne Bedingung an Port oder Benutzer; das Gegenüber wurde per `sudo nft add element` in das Set eingetragen, so wie der `dnsmasq` es für einen erlaubten Namen täte); alles mit `sudo` (`sudo setpriv --reuid=65534 dig` am Upstream offen; `sudo dig` als `root` ist mit der Variante gesperrt; die `sudoers`-Regel ist `(root) NOPASSWD: ALL`, `sudo -u nobody` scheitert). *Nicht messbar im Standardnetz:* IPv6 (keine Adresse außer `::1`), `127.0.0.11` (nicht vorhanden); dafür der Zusatz im benutzerdefinierten Netz (Befund 6). |
| D4 | **Gemessen.** Ist und Variante verhalten sich gleich: erster Lauf, zweiter Lauf im selben Container, Lauf nach `docker restart`, Lauf im neuen Container aus demselben Image — je ein `dnsmasq` (UID 65534, Pid lebt), 7 `accept`-Regeln (in der Variante 4 davon mit `meta skuid`, keine Anhäufung), beide Tabellen, `api.github.com` lokal `NOERROR`, Statusdatei `result=ok`. 50 erneute Läufe je Fall: 0 auffällige Zustände, 0 Meldungen. **Nicht erneut gemessen:** das Fenster ohne Filter beim erneuten Lauf (die Variante ändert den Ablauf nicht, nur zwei Regeln); das Startrennen unter Last. |
| D5 | **Gemessen** (D2 schließt den direkten Pfad): Der `dnsmasq` leitet Anfragen für beliebige Unternamen eines erlaubten Namens an den Upstream weiter (Variante: alle `kanal1.api.github.com`-Anfragen kamen am Upstream an, keine für `example.com`). Obergrenze im Laboraufbau (Gegenüber auf demselben Host als Upstream, sequenzielle Anfragen mit drei Labels zu 63 Hex-Zeichen): 3 650 bis 5 640 Anfragen je Sekunde, 0,7 bis 1,1 Millionen Label-Zeichen je Sekunde (etwa 0,4 bis 0,7 MB/s bei 5 Bit je Zeichen), Ist wie Variante; die direkte Abfrage am Upstream lag im Ist bei 9 300 und 11 600 Anfragen je Sekunde. Ein Empfänger braucht den autoritativen Server des erlaubten Namens; bei der Default-Allowlist ist das nicht der Messende. Echte Namen mit kontrollierbarem Server: nicht gemessen. |

**Befunde.**

1. **Die Variante schließt den direkten Pfad im Messstand vollständig** (UDP und TCP 53, beide Upstreams, vier von vier Anfragen) und lässt die Namensauflösung der erlaubten Namen unberührt (sieben von sieben, `getent`, HTTPS). Das Risiko aus §6 „die Regel trifft auch den `dnsmasq` oder erlaubte Namen“ ist für diesen Messstand **nicht eingetreten**.
2. **Der rote Test wird mit der Variante grün — als nachgestellter Ablauf, nicht als Lauf des Go-Tests.** Der Ablauf des Tests (`nestedPodmanEgressBuildScript`, gelesen in `internal/e2e/sandbox_devcontainer_docker_test.go`) ist Zeile für Zeile nachgestellt (plus ein `RUN cat /etc/resolv.conf`); Exit 1 im Ist, Exit 0 mit der Variante, in beiden Läufen. Die Ursache aus dem Cut-Befund bestätigt der Zähler: der verschachtelte Container fragt mit der UID 1000. **Nebenwirkung:** der Bau dauert 43 statt 7 Sekunden, weil der Rückfall auf `8.8.8.8` und `8.8.4.4` nach dem `REFUSED` des ersten Nameservers in die Frist läuft (36 s je unlistetem Namen); der Test hat 25 Minuten Frist.
3. **`podman run --dns 8.8.8.8` löst mit der Variante nichts mehr auf, auch keinen erlaubten Namen** (20 s je Name). Ein verschachtelter Container, der einen Upstream selbst setzt, verliert die Auflösung; mit dem Standard (`pasta`) und `--network=host` bleibt sie erhalten. Das ist eine Verhaltensänderung für Nutzer des Feature `nestedRuntime`.
4. **Was bleibt offen (D3, D5):** (a) Verkehr an Adressen im Set `allowed` auf jedem Port und Protokoll — das ist die Allowlist selbst, kein DNS-Befund und nicht über eine DNS-Regel zu schließen; (b) `sudo` — Guardrail, wie in [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) und im Handbuch ausgewiesen; (c) der Weg über Unternamen erlaubter Namen (D5). Keiner dieser Wege ist eine Umgehung der Regel der Variante; sie bestehen im Ist-Zustand ebenso.
5. **Die Regel hängt am Benutzer des `dnsmasq`, der heute nur ein Standardwert ist.** Die erzeugte Konfiguration trägt kein `user=`; `dnsmasq` fällt auf `nobody` zurück (gemessen, nicht gelesen im Quelltext von `dnsmasq`). Eine Umsetzung müsste den Benutzer in der Konfiguration festlegen und die Regel daraus ableiten; weicht ein Build ab, sperrt die Regel den `dnsmasq` selbst (die Auflösung bricht für alle, fail-closed). Der Benutzer `nobody` trägt sonst keine Prozesse; mit `sudo` ist er erreichbar (Befund 4b).
6. **Nebenbefund, nicht Gegenstand dieses Slice:** In einem **benutzerdefinierten Docker-Netz** ist der Docker-eigene Resolver `127.0.0.11` der einzige „Upstream“ des Containers; mit dem **erzeugten** Script (Ist wie Variante) ist er nicht erreichbar, und auch der lokale `dnsmasq` löst dort keinen erlaubten Namen auf (Kontrolle ohne Script: offen). Die Ursache ist nicht untersucht (Vermutung: die Anfragen des Docker-Resolvers an den Host-Upstream gehen aus dem Namensraum des Containers und treffen die Richtlinie `drop`). Ob ein solches Netz für das Sandbox-Profil vorgesehen ist, ist hier nicht gelesen.
7. **Nicht gemessen, mit Folge für eine Zusage:** `nestedRuntime: docker` — der verschachtelte `dockerd` läuft als `root`, sein Verkehr läuft über die Weiterleitung und nicht über die Kette `output`; ob die Variante dort etwas ändert, ist offen. Eine Handbuch-Bedingung „ein Prozess ohne `sudo` löst keinen Namen außerhalb der Liste auf“ braucht nach `AGENTS.md` §3.10 einen Test, der sie pinnt, und sie nennt den Messstand.

**Rückführungs-Trigger §4 („D3 zeigt mehrere Umgehungswege, deren Schließung jeweils eine eigene Messreihe braucht“).** Nicht eingetreten, und ich habe deshalb nicht gestoppt: Alle direkten Wege zum Upstream (TCP, UDP/TCP 53 an andere Resolver, andere Ports, 853, 443, ICMP, IPv6) sind mit der einen Regel der Variante geschlossen oder waren es schon; die offenen Wege (Befund 4) sind Eigenschaften der Allowlist und von `sudo`, keine Umgehungen der Regel, und brauchen keine eigene Messreihe für ihre Schließung, weil sie keine Schließung durch diese Regel haben. Die Einschätzung ist meine Lesart des Triggers; der Projektinhaber kann sie anders ziehen.

**Einordnung.**

- *Schließt die Messvariante den direkten Pfad?* Ja, im Messstand (Befund 1), für UDP und TCP 53 an die Upstreams der Liste.
- *Was bleibt offen?* Befund 4: Verkehr an Adressen im Set, `sudo`, der Weg über Unternamen erlaubter Namen. Für die Default-Allowlist (`github.com`, `githubusercontent.com`, `debian.org`) ist ein Datenkanal über Unternamen nur ein Sichtbarkeitsproblem beim Betreiber des Resolvers, kein Weg zu einem Empfänger des Messenden; bei vom Nutzer ergänzten Namen mit eigenem autoritativen Server wäre er einer (nicht gemessen).
- *Wirkung auf den roten Test:* mit der Variante grün (nachgestellt), mit 43 s statt 7 s Bau (Befund 2). Der Test selbst lief nicht.
- *Sicherheitsrelevanz:* Der heutige Pfad erlaubt einem Prozess ohne `sudo` die Anfrage **beliebiger Namen** an die Upstreams und damit einen DNS-Datenkanal zu jedem autoritativen Server der Welt (D5 misst die Weiterleitung bis zum Upstream, mit der Größenordnung 10^3 bis 10^4 Anfragen je Sekunde); die Variante begrenzt ihn auf Unternamen der erlaubten Namen. Das ist ein wirklicher Gewinn gegenüber dem Ist, kein Ersatz für eine Sandbox-Grenze: mit `sudo` ist die Restriktion weiter aufhebbar, und Verkehr an erlaubte Adressen bleibt unbeschränkt im Port.

**Empfehlung (Entscheidung beim Projektinhaber, keine Umsetzung in dieser Phase).**

- **Weg A — `nft meta skuid`-Regel im Template, Folge-ADR zu [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md): empfohlen.** Er schließt den gemessenen Pfad, lässt die Auflösung erlaubter Namen unberührt, macht den roten Test ohne Anpassung grün, und D4 zeigt keine Regression. Bedingungen aus den Befunden: Benutzer des `dnsmasq` in der Konfiguration festlegen (Befund 5); Golden Case und ein Docker-Test, der die Zusage pinnt und die gemessene Grenze benennt (direkter Pfad ohne `sudo` abgewiesen, erlaubte Namen lösen auf; Befund 7); die 36 s Wartezeit und das Verhalten von `--dns` im Handbuch ausweisen (Befunde 2 und 3); `nestedRuntime: docker` und das benutzerdefinierte Netz vor einer Zusage prüfen oder als Grenze ausweisen (Befunde 6 und 7).
- **Weg B — Test anpassen (Upstream pinnen), Grenze im Handbuch lassen:** möglich (der Cut hat `--dns 9.9.9.9` als grün gemessen), lässt den gemessenen offenen Pfad aber stehen und widerspricht dem Ziel des Slice (Code zum Ziel, Entscheidung 2026-10-18).
- **Weg C — nicht schließbar:** durch die Messung nicht gestützt, soweit der Messstand reicht; er ist nicht Colima/arm64, und `nestedRuntime: docker` ist offen.

**Zu jeder Zusage das rot gesehene Gegenbeispiel** (Messpunkt: derselbe Container, Fall `ist` gegen Fall `skuid`; keine Produktänderung):

| Zusage | Was rot wäre | Gesehen |
|---|---|---|
| Variante schließt den direkten Pfad | die Regel ohne `meta skuid` (Fall `ist`) | ja: `AUSSAGE ist upstream-direkt … =antwort` (4 von 4) gegen `AUSSAGE skuid … =keine`, `rohdaten-04-d2.txt` |
| Erlaubte Namen lösen weiter auf | eine Regel, die auch den `dnsmasq` sperrt (`DNSUID=1234` in `variante.sh`) | ja: lokale Auflösung von `api.github.com` `REFUSED` statt `NOERROR`, `rohdaten-09b-gegenprobe.txt` (beide Läufe) |
| Der Ablauf des roten Tests wird grün | Fall `ist` | ja: Exit 1 gegen Exit 0, beide Läufe, `rohdaten-05-d2-nested.txt` |
| Sonden erkennen einen offenen Pfad | Fall `ohne` (kein Startscript) | ja: je Sonde offen, mit den vermerkten Ausnahmen (Protokoll, Lesehilfen) |
| D4: jeder Lauf ist auffällig-frei | ein Zustand, der vom Soll abweicht | ja, ein Fall: dieselbe Prüfung (`09-d4.sh`) meldet den Zustand mit falscher UID als auffällig, `rohdaten-09b-gegenprobe.txt`; sie ist am Fall `ist` (`skuid-regeln=0`, `direkt=antwort`) gegen `skuid` (`4`, `keine`) zusätzlich unterscheidend; andere Abweichungen (zweiter `dnsmasq`, Regelanhäufung) wurden nicht künstlich erzeugt |

**Sensoren.**

| Sensor | Ergebnis |
|---|---|
| `make ablage-pruefung` | OK, 26 Fälle, 0 Fehlschläge |
| `make docs-check`, `make doc-planning` | je 523 Dateien, 0 Befunde |
| `make doc-trace` | 83 Anforderungen, 0 Waisen |
| `make gates` | Exit 0 (Coverage 91,40 % gegen 90 %; `ablage-pruefung` OK, 26 Fälle, 0 Fehlschläge) |
| `make test-docker`, `make verify-depguard` | **nicht gelaufen**: kein Produktcode, keine Importregel berührt |

**Plan-Tabelle gegen `git diff --stat -M d2f8294..HEAD`** (Stand nach dem Stand-Commit dieses Abschnitts): 42 Dateien: 24 unter `tools/messung/egress-dns-upstream/`, 16 unter `docs/messungen/slice-v1-egress-dns-upstream-schliessen/` (README und 15 Rohdaten), `docs/messungen/README.md` und diese Slice-Datei; jede ist in der Tabelle geführt (erste Zeile, Zeile zu `docs/messungen/README.md`, Zeile *diese Slice-Datei*). Die Kopie von Lauf 2 liegt als `rohdaten-lauf2-aussagen.txt` im Protokoll, die übrigen Rohdaten von Lauf 2 (Scratchpad) sind nicht im Repo. Gelesen als **unverändert**: `egress-init.sh.tmpl` (Zeilen 36 und 63 bis 66 gelesen; `variante.sh` ersetzt genau diese beiden Regeln, `sed` bricht sonst ab), `internal/e2e/sandbox_devcontainer_docker_test.go` (`nestedPodmanEgressBuildScript` und `TestE2E_SandboxDevcontainer_NestedPodman_BuildEgress` gelesen, Ablauf nachgestellt), `docs/user/benutzerhandbuch.md` (Absatz *Grenze der Namensauflösung* gelesen: Stand ist „Colima, je ein Lauf“ und die Cut-Messung zum verschachtelten Podman; nicht berührt), `docs/plan/adr/` ([`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) Alternative C2 und [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md) M1 a gelesen, nicht berührt).

### Umsetzung (Implementer, 2026-10-19; Phase 2: Umsetzung nach [`ADR-0026`](../../adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md); nichts gepusht, kein Tag)

Commits: `f691fab` (Template, Unit- und Docker-Tests, Doctor-Meldung), `6d85215` (Replay-Tabelle), `c1e6c07` (Handbuch, READMEs, CHANGELOG); der Stand dieses Abschnitts folgt als eigener Commit.

**Umsetzung.** `egress-init.sh.tmpl`: `dns_user=nobody` einmal gesetzt; `dns_uid=$(id -u "$dns_user")` (numerisch, wie gemessen; der Name in `meta skuid` ist nicht gemessen und nicht verwendet); die Konfiguration schreibt `user=$dns_user`; beide Upstream-Regeln tragen `meta skuid "$dns_uid"`. Fehlt der Benutzer (Prüfung vor dem Löschen der Tabellen), oder läuft der Prozess aus der Pid-Datei nach dem Start unter einer anderen UID (Prüfung vor dem Schreiben von `/etc/resolv.conf`), meldet das Script auf `stderr`, schreibt `failed` und endet mit Exit 11, unabhängig von `onUnavailable`. Die Pid-Datei wird vor dem Start gelöscht, damit ein nicht gestarteter `dnsmasq` die Prüfung nicht über einen alten Eintrag besteht. Bei Abbruch nach dem Anlegen der Tabellen bleibt `/etc/resolv.conf` unberührt; die Tabellen (Richtlinie `drop`) bleiben, die Statusdatei trägt `failed`.

**Wortlaut der Meldungen** (`stderr`): `u-boot sandbox: egress restriction not set up: the user 'nobody' for the local resolver does not exist` und `u-boot sandbox: egress restriction failed: the local resolver is not running as 'nobody' (uid 65534)`.

**Abweichung vom Auftrag/Plan:** keine vom ADR. Zusätzlich geändert, weil sonst falsch: die Doctor-Meldung für `failed` des Egress-Eintrags (`egressProbe`, `devcontainer_fallback.go`) und der Kommentar zu `SandboxStatusFailed` nannten nur `onUnavailable: fail`; beide nennen jetzt auch den neuen Grund. Lastenheft und Spezifikation unverändert.

**Sensoren.**

| Sensor | Ergebnis |
|---|---|
| `make test` | grün |
| `make gates` | Exit 0 (Coverage 91,40 % gegen 90 %) |
| `make test-docker` (volles Paket, Stand vor der Doctor-Meldung) | `ok internal/e2e 611,4 s`, alle Pakete grün |
| gezielt `TestE2E_SandboxDevcontainer_NestedPodman_BuildEgress` mit `-v` (Aufruf wie `make test-docker`, mit `-run`) | PASS in 62 s, kein SKIP; vor der Umsetzung rot (Cut-Befund, Messprotokoll) |
| gezielt `TestE2E_SandboxDevcontainer_EgressUpstreamDNS` | PASS in 25 s |
| `make verify-depguard` | nicht gelaufen: keine Importregel berührt |

**Wächter-Mutationen** (je eine Mutation am Template, rot gesehen, danach zurückgesetzt):

| Wächter | Mutation | Gesehen |
|---|---|---|
| Docker-Test „direkt nicht auflösbar“ | `meta skuid` aus beiden Regeln entfernt | rot: `DIRECT udp/tcp example.com=answer` (4 von 4) |
| Unit-Test (Regel je Protokoll) | `meta skuid` nur aus der `tcp`-Regel entfernt | rot: Zählung und „Regel ohne Benutzerbedingung“ |
| Docker-Test (UID des `dnsmasq` gleich Regel) | `user=daemon` statt `user=$dns_user` | rot: `INIT_EXIT=11` |
| Unit-Test (`user=`-Zeile) | `user=`-Zeile gelöscht | rot: `echo "user=$dns_user"` kommt 0-mal vor (das Verhalten bleibt gleich, weil `nobody` der Standard des `dnsmasq` ist; den Wächter trägt der Unit-Test) |
| Docker-Test (UID-Abweichung bricht ab) | UID-Prüfung nach dem Start abgeschaltet | rot: `ABORT uid-differs exit=0 result=ok` (der Zustand der Gegenprobe der Messung) |
| Docker-Test (fehlender Benutzer bricht ab) | Prüfung vor den Tabellen abgeschaltet | rot: `ABORT user-missing exit=0 result=ok` |

**Replay / Golden** (`harness/replay.md`): *Affected generators:* `generate devcontainer` mit `devcontainer.sandbox.egress.enabled`. *Golden cases added:* `TestGenerateDevcontainer_SandboxEgress_UpstreamDNSOnlyForTheResolverUser` (Teilstring- und Reihenfolge-Pins des erzeugten `egress-init.sh`, beide Richtlinien). *Golden cases updated:* `TestGenerateDevcontainer_SandboxStartScriptsWriteStatusFiles` (Fall `egress`: `failed` 1 → 4, `exit 11` direkt danach 1 → 4). *Golden cases replayed:* `TestGenerateDevcontainer_SandboxEgress*`, `…_StartScriptsWriteStatusFiles`, `…_NotGenerated`, `…_EgressRebuildNotice*` (grün, `make test`). Es gibt keine Whole-File-Golden-Datei zu `egress-init.sh`. *Intentional output changes:* der verwaltete Block von `egress-init.sh` (skuid-Bedingung, `user=`-Zeile, zwei Abbruchzweige, `rm -f` der Pid-Datei) wegen [`ADR-0026`](../../adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md); Konfigurationen ohne Egress erzeugen kein Script und bleiben unverändert.

**Offen für Folge-Vorgänge (Meldung, kein Eingriff):** [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) und die Egress-Zeile der Degradation (`spec/spezifikation.md`: Capability-Prüfung, Exit 11 nur bei `onUnavailable: fail`) kennen den Abbruch bei `warn` nicht; die Entscheidung liegt beim Spezifikations-Slice nach [`ADR-0026`](../../adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) Entscheidung 5.

**Plan-Tabelle gegen `git diff --stat -M 4524ffe..HEAD`** (Stand nach dem Stand-Commit dieses Abschnitts): 11 Dateien; jede ist in der Tabelle geführt (`egress-init.sh.tmpl`, `generate_sandbox_test.go`, `sandbox_devcontainer_docker_test.go`, `devcontainer_fallback.go`, `sandbox_status.go`, `harness/replay.md`, `benutzerhandbuch.md`, `CHANGELOG.md`, `README.md`, `README.de.md`, diese Datei). Gelesen als **unverändert**: `internal/e2e/sandbox_status_docker_test.go` (`TestE2E_SandboxDevcontainer_StatusFiles_Egress` gelesen: prüft `ok`, `fallback`, `failed` über die Capability, nicht berührt, im vollen `make test-docker` grün), `spec/spezifikation.md` und `spec/lastenheft.md` (gelesen: kein Satz zur DNS-Regel; Zeile zur Capability-Prüfung siehe oben), `docs/plan/adr/` (nicht berührt).

### Fix-Runde zum Review (Implementer, 2026-10-19)

Review: [`2026-10-19-slice-v1-egress-dns-upstream-schliessen.md`](../../../reviews/2026-10-19-slice-v1-egress-dns-upstream-schliessen.md) (0 HIGH, 2 MEDIUM, 4 LOW, 5 INFO). Behoben: **MEDIUM 1** `egressProbe`: Text nennt drei Gründe (Capability unter `fail`; Benutzer fehlt; Resolver startet nicht oder läuft unter anderem Benutzer), der Hinweis nennt beide Wege (Capability, oder die Meldung des Startscripts im Container-Log); `TestDoctor_SandboxFallback_EgressFailedNamesEveryCause` pinnt beides. **MEDIUM 2** Handbuch: Kette `egress` vor Daemon-Start (gelesen in `devcontainer_sandbox.go`, ohne Test mit Abbruch), Zustand nach dem Abbruch (Docker-Test `STATE`-Zeilen: Benutzer fehlt: keine Tabellen; Resolver startet nicht oder UID weicht ab: Tabellen mit `drop`, `resolv.conf` unverändert). **LOW** Test-Kommentar zur Zählung berichtigt (im Baum gesucht, kein weiterer Treffer); Rückgabewert des `dnsmasq`-Starts geprüft, eigene Meldung `the local resolver did not start`, gleicher Abbruchweg (Exit 11, `failed`) — keine Abweichung von [`ADR-0026`](../../adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) Entscheidung 2, die UID-Prüfung bleibt; Fehlersuche, `onUnavailable`-Tabelle und Statusdatei-Tabelle im Handbuch nachgezogen; positive Kontrolle für TCP als `nobody`; Messprotokoll *Nicht gemessen* auf den Ist-Stand. **INFO** [`ADR-0026`](../../adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) trägt „Wortlaut … steht zur Freigabe“ in einem `Accepted`-ADR (Immutabilität, nicht geändert, nur hier geführt); die Handbuch-Messwerte (36 s, 43 gegen 7 s, `--dns`, `--network=host`) stehen als „gemessen; kein Test pinnt sie“; die DoD-Zeile zur Abnahme des Tests ist bis zum Verifier offen.

**Sensoren der Fix-Runde:** `make test`, `make lint`, `make gates`, `make docs-check`, `make doc-trace`, `make doc-planning`, `make ablage-pruefung` grün; gezielt mit `-v` PASS: `TestE2E_SandboxDevcontainer_EgressUpstreamDNS` (Unterfälle `direct_path_closed`, `abort_branches`), `…_NestedPodman_BuildEgress`, `…_StatusFiles_Egress`, `…_Egress`, `…_EgressSubdomains`.

**Wächter-Mutationen der Fix-Runde (rot gesehen, zurückgesetzt):** `meta skuid` nur aus der `tcp`-Regel entfernt: jetzt auch der Docker-Test rot (`DIRECT tcp example.com=answer`); Prüfung des `dnsmasq`-Rückgabewerts abgeschaltet: Docker-Test rot (`ERR not-started … is not running as 'nobody'` statt `did not start`); `echo "nameserver 127.0.0.1"` vor die UID-Prüfung gesetzt: Docker-Test rot (`STATE uid-differs … resolver-in-resolv.conf=1`); Hinweis der Doctor-Meldung gekürzt (`missing user 'nobody'` entfernt): Unit-Test rot.

**Plan-Tabelle gegen `git diff --stat -M d2f8294..HEAD`** (nach dem Stand-Commit der Fix-Runde): 56 Dateien: 24 Skripte, 16 Dateien des Messprotokolls (inkl. README), `docs/messungen/README.md`, ADR und ADR-Index, Review-Report, diese Datei, Handbuch, `CHANGELOG.md`, zwei READMEs, `harness/replay.md` und sechs Dateien unter `internal/` (Template, `generate_sandbox_test.go`, `devcontainer_fallback.go`, `devcontainer_fallback_test.go`, `sandbox_status.go`, Docker-Test); jede ist in der Tabelle geführt.

**Fix-Runde zum Verification-Report:** V1 (`failed`-Zweige: vier in `harness/replay.md`, Test-Kommentar, diese Datei; Unit-Test zählt 4, richtig), V2 (Handbuch: `resolv.conf` nach dem Abbruch getrennt für Erst- und Wiederholungslauf, Wiederholungslauf als gelesen), V4 (Symbolbezüge in `slice-v1-egress-init-ohne-luecke`). Plan-Tabelle gegen `git diff --stat -M d2f8294..HEAD`: 58 Dateien, die zwei neuen sind eingetragen.

**Risiko, neu:** Mit der Prüfung des Rückgabewerts endet ein Wiederholungslauf, bei dem der alte `dnsmasq` (SIGTERM ohne Warten) den Port noch hält, mit Exit 11 statt still ohne Resolver; die Messung D4 (50 Wiederholungen, keine Auffälligkeit) zeigt den Fall nicht, ein Lauf unter Last ist nicht gemessen.

**Übergabe an den Spezifikations-Slice, ergänzt:** neben [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) (`failed` bei `fail`) und der Capability-Zeile der Degradationstabelle ist auch die Zeile *Kette `egress && docker`* im Abschnitt *Stand* dort (`spec/spezifikation.md` Z. 655, gelesen) vom Code überholt (Abbruch unter `warn`).

## 4. Trigger

**Start** (`open` → `in-progress`, ggf. über `next`): `Verantwortlich:` gesetzt; [`welle-sandbox-egress-ollama`](../done/welle-sandbox-egress-ollama.md) ist geschlossen (Reihenfolge-Entscheidung des Projektinhabers); kein anderer Slice in `in-progress/` (WIP-Limit 1). Der Start hängt nicht am Release.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): D3 zeigt mehrere Umgehungswege, deren Schließung jeweils eine eigene Messreihe braucht.
- `in-progress` → `open` (blockiert): die Messumgebung (Docker auf Colima mit `NET_ADMIN`) steht nicht zur Verfügung; die Evidence nennt den Grund (kein stilles Rot).

## 5. Closure-Trigger

DoD vollständig, Review-Report und Verification-Evidence liegen vor, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7); die Entscheidung des Projektinhabers zum gewählten Weg (schließen oder Grenze ausweisen) ist im Gespräch belegt und in §7 vermerkt.

## 6. Risiken und offene Punkte

- **Die Regel gegen direkte DNS-Abfragen trifft auch den `dnsmasq` oder erlaubte Namen** (Auflösung bricht). — **Ausgang:** weiter offen bis zur Messung D2.
- **Eine Regel nach Benutzer trägt nur, solange der Nutzer ohne `sudo` bleibt;** mit `sudo` ist jede Regel aufhebbar. — **Ausgang:** weiter offen; Grenze im Handbuch und in der Folge-ADR zu benennen (Guardrail, [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)).
- **Ein Datenkanal über erlaubte Namen bleibt möglich** (nicht gemessen, D5). — **Ausgang:** weiter offen bis zur Messung.
- **Die Messung gilt nur für Colima/arm64.** — **Ausgang:** weiter offen; Grenze in Handbuch und Protokoll benennen.
- **Der Slice ändert das Startscript, dessen Startrennen aus [`slice-v1-egress-wildcard-pruefung`](../done/slice-v1-egress-wildcard-pruefung.md) nicht ausgeschlossen ist.** — **Ausgang:** weiter offen; D4 misst die Startrennen-Läufe wie der Vorgänger.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** voraussichtlich berührt werden `internal/hexagon/application` (APP, Greenfield; das Template liegt unter `templates/devcontainer/`), `internal/e2e` (E2E, Greenfield; ein Docker-Test) und `docs/` (DOK, Greenfield) laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); jede erfüllt die Schwelle aus zwei von drei Achsen laut Modus-Deklaration.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md) (offen; ein neuer Docker-Test mit Image-Bau ist ein weiterer Kandidat), [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md) (offen; ein Skip ist im CI-Aufruf ohne `-v` nicht von einem Pass zu trennen), [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (verkörpert; die Aussagen zum Template oben sind an den genannten Zeilen gelesen, die Benutzer-Frage des `dnsmasq` ist als ungelesen markiert) und [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (verkörpert; eine Handbuch-Bedingung zum DNS-Pfad braucht einen Test oder entfällt). Evidenz-/Diskrepanz-Risiko vor der Messung: mittel.

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
