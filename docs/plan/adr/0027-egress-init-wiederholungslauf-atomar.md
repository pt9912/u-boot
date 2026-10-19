# ADR 0027: Wiederholungslauf von `egress-init.sh` — eigene Tabellen in einer einzigen `nft -f`-Transaktion, Resolver-Wechsel ohne Not, Zusage nur im Messstand

**Status:** Accepted

**Datum:** 2026-10-19

**Autor:** pt9912

**Bezug:** [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [ADR-0012](0012-devcontainer-egress-firewall.md), [ADR-0022](0022-fallback-ausweisung-statusdateien-im-container.md), [ADR-0025](0025-egress-liste-im-laufenden-container.md), [ADR-0026](0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md)

**Schärft:** [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Mechanismus: nur das Verhalten des Startscripts bei einem zweiten Lauf; Allowlist, Guardrail-Satz und Degradation bleiben) sowie — als Stelle, die nach `Accepted` zu prüfen ist — [`SPEC-045`](../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung). Gelesen am 2026-10-19: `spec/spezifikation.md` trägt keinen Satz zum zweiten Lauf des Startscripts (Suche nach „zweiten Lauf“, „idempot“, „erneut“, „wiederhol“; die Treffer betreffen Generator-Idempotenz und Namenskürzung, nicht `egress-init.sh`); [`SPEC-019`](../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) führt die Allowlist, nicht den Wiederholungslauf. Eine Nachführung ist damit nicht nötig; ob die Spezifikation eine Zusage samt Grenze aufnimmt, entscheidet der Spezifikations-Slice nach dem Code, diese ADR ändert sie nicht. **Ergänzt**, ersetzt nichts: [ADR-0025](0025-egress-liste-im-laufenden-container.md) (dort die Konsequenz „das Startscript ist bei einem zweiten Lauf nicht atomar“, Messung M2) und [ADR-0012](0012-devcontainer-egress-firewall.md) bleiben unverändert `Accepted`; ebenso [ADR-0026](0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[ADR-0025](0025-egress-liste-im-laufenden-container.md) maß als M2 (`arm64`, Colima), dass der erneute Lauf des Startscripts die Tabellen löscht und danach neu anlegt: Fenster ohne Filter-Regelwerk 24 bis 45 ms, Ausfall des lokalen `dnsmasq` 28 bis 55 ms; eine **atomare Variante** (Messskript, kein Produktcode: Wechsel in einer nft-Transaktion, `dnsmasq` erst nach dem Ende des alten) zeigte in fünf Läufen keinen Zustand `offen` und einen Ausfall von höchstens 24 ms. Die Messfrage A7 auf `amd64` ([Messprotokoll](../../messungen/slice-v1-devcontainer-ollama-amd64-messung/README.md), Befund B-1) bestätigte M2 für die heutige Fassung: Verbindungsversuch zum Gateway nicht aufgehalten in 27 von 30 Läufen (Spanne bis 48 ms), im zweiten Lauf 28 von 30 (bis 50 ms); die Namensauflösung des Containers fiel in 30 und 27 Läufen kurz aus ([`rohdaten-08.txt`](../../messungen/slice-v1-devcontainer-ollama-amd64-messung/rohdaten-08.txt), [`rohdaten-lauf2-08.txt`](../../messungen/slice-v1-devcontainer-ollama-amd64-messung/rohdaten-lauf2-08.txt)). Die Zählung ist eine Untergrenze (Sonde im Abstand von rund 7 ms, Verbindungsversuch zum Gateway auf Port 9, IPv4).

Der Projektinhaber entschied am 2026-10-19: Die Änderung bekommt diese kurze Folge-ADR, die [ADR-0025](0025-egress-liste-im-laufenden-container.md) ergänzt; **eine neue Messung entfällt** (abgebrochen). Belege sind M2 (`arm64`, atomare Variante) und B-1 (`amd64`, heutige Fassung). Daraus folgt, was diese ADR zusagen darf und was nicht (Entscheidung 3).

**Stand im Code (gelesen 2026-10-19).** `internal/hexagon/application/templates/devcontainer/egress-init.sh.tmpl`: Abschnitt „Re-run safe“ — `kill` des Resolvers aus der Pid-Datei ohne Warten, `rm -f` der Pid-Datei, `nft delete table ip u_boot_egress` und `nft delete table ip6 u_boot_egress6`; danach die `dnsmasq`-Konfiguration schreiben, dann `nft add table`/`set`/`chain` und jede `rule` einzeln (IPv4: Loopback, `established,related`, `@allowed`, je Upstream zwei Regeln mit `meta skuid "$dns_uid"`; IPv6: Tabelle, Kette `policy drop`, Loopback), dann `dnsmasq --conf-file` mit Prüfung des Rückgabewerts und der Prozess-UID (Exit `11`, Statusdatei `failed`), danach `echo "nameserver 127.0.0.1" > /etc/resolv.conf` und `write_status ok`. Die `nft add`-Aufrufe werden nicht auf ihr Ergebnis geprüft. Exit-Codes: `0` (Fallback ohne Capability, Erfolg) und `11` (fehlende Capability bei `fail`; Benutzer fehlt, `dnsmasq` startet nicht, `dnsmasq` unter fremder UID).

## Entscheidung

1. **Der Wiederholungslauf baut die eigenen Tabellen in einer einzigen `nft -f`-Transaktion neu auf.** Beide Tabellen (`ip u_boot_egress` mit Set, Kette und allen Regeln; `ip6 u_boot_egress6` mit Kette und Loopback-Regel) werden in **einem** `nft -f`-Aufruf ersetzt, ohne vorheriges eigenständiges `nft delete table` als eigener Schritt: Zwischen altem und neuem Regelwerk gibt es keinen Zustand ohne Tabelle. Die Syntax innerhalb der Transaktion (`delete` und Neuanlage im selben Eingabestrom oder ein gleichwertiges Idiom) legt der Plan der Umsetzung fest; sie muss in der Transaktion auf dem ersten Lauf (keine Tabelle vorhanden) und auf dem Wiederholungslauf funktionieren. Die Regeln, auch die `meta skuid`-Bedingung am Upstream-DNS ([ADR-0026](0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md)), bleiben inhaltlich gleich. **Der Resolver-Wechsel** (`dnsmasq`): Die neue Konfiguration wird vor dem Wechsel geschrieben; der alte Resolver bekommt `kill`, das Script wartet begrenzt auf sein Ende (der Port ist bis dahin belegt, `bind-interfaces`), startet dann den neuen und prüft Rückgabewert und UID wie bisher; `/etc/resolv.conf` wird **zuletzt** geschrieben, nach dem geprüften Start. „Ohne Not“ heißt: kein Ausfall, der sich vermeiden ließe; ein kurzer Ausfall des lokalen Resolvers zwischen Ende des alten und Start des neuen bleibt, weil beide denselben Port binden. **Unverändert:** die Exit-Codes `0` und `11` mit ihren Bedingungen, die Statusdatei (`ok`/`fallback`/`failed`, Schreibweise über Temporärdatei und `mv`), die `skuid`-Regel und die Prüfung der Prozess-UID aus [ADR-0026](0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md), die Wahl des Benutzers `nobody`. Der Inhalt des Sets `allowed` wird beim Wiederholungslauf wie bisher verworfen (M2: Set nach dem Lauf leer); die ADR sagt nicht, dass er erhalten bleibt.

2. **Abgrenzung.** (a) Kein Neu-Laden vom Host: [ADR-0025](0025-egress-liste-im-laufenden-container.md) Entscheidung 1 gilt unverändert; der Weg „Liste ändern“ bleibt `u-boot generate devcontainer` und Neuaufbau des Containers. Die Änderung betrifft nur, was geschieht, wenn jemand im Container (mit `sudo`) das Startscript erneut ausführt, und der `postStartCommand` beim Container-Start. (b) R1 bis R4 und S bleiben offen und werden nicht entschieden ([ADR-0025](0025-egress-liste-im-laufenden-container.md) Entscheidung 3); die atomare Variante ist dort als Zelle von R1 geführt und wird hier nur für das Startscript übernommen, ohne R1 zu wählen. (c) Kein neuer Port, kein Programmaufruf des Hosts im Container, keine Änderung des Lastenhefts.

3. **Die Zusage und ihre Grenze.** Die Restriktion bleibt eine Leitplanke, keine Sicherheitsgrenze ([`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)); der Wiederholungslauf braucht `sudo`, mit dem sich die Restriktion ohnehin aufheben lässt ([ADR-0025](0025-egress-liste-im-laufenden-container.md) Entscheidung 2). Zugesagt wird **die Struktur**: Das Script ersetzt die eigenen Tabellen in einer Transaktion und schreibt `resolv.conf` zuletzt; ein Golden Case pinnt sie. **Nicht zugesagt wird ein Fenster von null und kein Ausfall der Namensauflösung.** Belegt ist:
   - `arm64` (Colima), atomare Variante als Messskript: in fünf Läufen kein Zustand `offen`, Ausfall des lokalen `dnsmasq` höchstens 24 ms ([ADR-0025](0025-egress-liste-im-laufenden-container.md) M2).
   - `amd64`: nur die **heutige** Fassung gemessen (B-1), also die Größe des Fensters ohne die Änderung; für die Variante ist auf `amd64` **nichts gemessen**.
   - **Ungemessen und ausgewiesen:** die Wirkung der Variante auf `amd64`, auf Docker Desktop und in der CI; die IPv6-Tabelle (in M2 und B-1 nicht gemessen; die Sonden decken IPv4); das Verhalten unter Last des Hosts; der Wiederholungslauf mit `nestedRuntime`; ob die Transaktion im Produkt-Script dieselben Zeiten zeigt wie das Messskript (die Variante war ein Skript, kein Produktcode). Die Zählungen sind Untergrenzen; eine Null in einer späteren Messung hieße „kein Treffer“, nicht „kein Fenster“.
   - Eine Bedingung im Handbuch oder in der Spezifikation zum zweiten Lauf hat einen Test oder entfällt (`AGENTS.md` §3.10; Beobachtung `BEO-DOK/zusage-weiter-als-messung`): Das Handbuch trägt für die neue Fassung nur, was gemessen oder gepinnt ist, und weist die Zählungen von B-1 als Stand der alten Fassung aus.

4. **Eine fehlgeschlagene Transaktion**: Ein `nft -f`, das nicht durchläuft, ändert nichts am Regelwerk (die Transaktion ist atomar, nach der `nft`-Dokumentation; im Produkt-Script ungemessen). Das Script bricht dann mit Meldung, Statusdatei `failed` und Exit-Code `11` ab, unabhängig von `onUnavailable`, wie die Abbruchfälle aus [ADR-0026](0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) Entscheidung 2 (Umgebungsproblem, vom Generator nicht zu verantworten). Es führt keinen neuen Exit-Code ein; heute ist der Fall ungeprüft (Statusdatei `ok` trotz fehlender Regeln möglich). Beim Wiederholungslauf bleibt dabei das alte Regelwerk bestehen.

**Freigabe-Punkte des Projektinhabers vor `Accepted`** (Empfehlungen, vom Projektinhaber alle angenommen, siehe Geschichte):

- a) Beide Tabellen in einer `nft -f`-Transaktion ersetzen, ohne eigenes `delete`-Fenster; `dnsmasq`-Wechsel mit begrenztem Warten und `resolv.conf` zuletzt; Exit-Codes, Statusdatei und `skuid`-Regel unverändert (Entscheidung 1). Empfehlung: ja.
- b) Abgrenzung: kein Neu-Laden vom Host, R1 bis R4 und S offen (Entscheidung 2). Empfehlung: ja.
- c) Zugesagt wird die Struktur, nicht ein Fenster von null; `amd64`-Wirkung, IPv6, Last, `nestedRuntime` und Produkt-Zeiten ungemessen und ausgewiesen (Entscheidung 3). Empfehlung: ja.
- d) Fehlgeschlagene Transaktion endet mit Exit-Code `11`, `failed`, auch bei `warn` (Entscheidung 4). Das ist ein neuer Abbruchzweig innerhalb des bestehenden Codes und der Klasse aus [ADR-0026](0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md). Empfehlung: ja; wenn nein, entfällt Entscheidung 4 und der Fall bleibt wie heute ungeprüft.

## Verglichene Alternativen

Jede Zeile nennt, was die Belege tragen; „nicht gebaut“ heißt, dass die Zelle aus M2 und B-1 abgeleitet ist, nicht aus einer Umsetzung im Template.

| Option | Pro | Contra |
|---|---|---|
| Nichts tun (Ist: `delete`, dann Neuaufbau mit einzelnen `nft add`) | keine Änderung | Fenster ohne Restriktion und Ausfall der Auflösung auf `arm64` (M2) und `amd64` (B-1: 27 und 28 von 30 Läufen); das Handbuch müsste die Grenze weiter ausweisen |
| **Eine `nft -f`-Transaktion, Resolver-Wechsel geordnet** | auf `arm64` kein Zustand `offen`, Ausfall ≤ 24 ms (M2, Messskript); eine Änderung im Template, keine neue Fähigkeit, keine neue Messung nötig | `amd64`, IPv6 und Produkt-Zeiten ungemessen (Entscheidung 3); ein kurzer Resolver-Ausfall bleibt (gleicher Port); Set `allowed` wird weiter verworfen |
| Neu-Laden vom Host (R1/R2) | Auslöser auf dem Host | [ADR-0025](0025-egress-liste-im-laufenden-container.md) Entscheidung 1 und 3: ausgeschlossen bzw. offen; braucht Port, Spezifikations- und Lastenheft-Änderung; löst das Fenster des Startscripts nicht von selbst |
| Zweiter `dnsmasq` auf anderem Port, Umschalten per `resolv.conf` | kein Resolver-Ausfall denkbar | nicht gemessen; ändert Konfiguration, Prüfung der UID und Port-Belegung; mehr Zustand als der Befund verlangt |
| Set `allowed` beim Wechsel erhalten | bestehende Auflösungen blieben erreichbar | nicht gemessen; widerspricht dem Befund M2 (Set leer) nicht, ändert aber den Zustand nach dem Lauf; eigener Vorgang, wenn Nutzer es verlangen |

## Konsequenzen

- Positiv: Der Befund B-1 hat eine begründete Gegenmaßnahme mit einem Beleg auf `arm64`; das Script behauptet seine Wiederholbarkeit mit einer Struktur, die ein Golden Case pinnt; ein nicht durchlaufender Neuaufbau wird laut statt still (Entscheidung 4, falls freigegeben).
- Negativ: Die Wirkung auf dem Zielpfad `amd64` ist nicht belegt; ein kurzer Resolver-Ausfall bleibt; die Sonde misst eine Untergrenze, ein Docker-Test kann das Fenster nur mit Wiederholungen und nur als Treffer abweisen (Risiko im Slice).
- Folgepflicht (eigene Schritte der Umsetzung, nicht dieser ADR): Template `egress-init.sh.tmpl`, Golden Case des erzeugten Scripts (Fresh-State und zweiter Lauf, `harness/replay.md`), Docker-Test, Handbuch (Abschnitt *Internet-Beschränkung (Egress)*: Zählungen und „Untergrenze“ als Stand der alten Fassung, Messstand der neuen) und `CHANGELOG.md`.
- **Nicht verifiziert:** die Wirkung der Variante auf `amd64`, Docker Desktop und CI; die IPv6-Tabelle; Last des Hosts; `nestedRuntime`; die Dauer der Transaktion im Produkt-Script; die Atomarität eines fehlgeschlagenen `nft -f` im Produkt-Script. **Akzeptiertes Negativ:** Die Messung entfällt auf Entscheidung des Projektinhabers (2026-10-19); die Lücken bleiben damit ausgewiesen, eine Folgepflicht „später messen“ entsteht nicht, weil die Wirkung ein Guardrail-Fenster für einen gleichzeitigen Prozess während eines `sudo`-Laufs betrifft ([ADR-0025](0025-egress-liste-im-laufenden-container.md) Entscheidung 2).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Golden Case / Unit | Das erzeugte `egress-init.sh` ersetzt beide Tabellen in einer `nft -f`-Eingabe, trägt kein eigenständiges `nft delete table` vor dem Neuaufbau, wartet nach dem `kill` begrenzt auf das Ende des Resolvers und schreibt `resolv.conf` nach dessen geprüftem Start; Exit-Codes und Statusdatei wie zuvor; Konfigurationen ohne Egress rendern byte-identisch | `make test` |
| Integrationstest (`//go:build docker`) | Das Startscript läuft wiederholt; Exit `0`, Statusdatei `ok`, genau ein `dnsmasq` mit UID der Regel, beide Tabellen vorhanden, ein erlaubter Name löst danach auf. Ein Fenster wird nur als Treffer abgewiesen (Untergrenze), eine Nullzählung trägt keine Zusage | `make test-docker` |

## Re-Evaluierungs-Trigger

Eine Messung auf `amd64`, Docker Desktop oder in der CI, die für die neue Fassung ein Fenster oder einen längeren Ausfall der Auflösung zeigt; eine Messung der IPv6-Tabelle mit Befund; ein Nutzerbericht, dass ein Wiederholungslauf legitime Verbindungen stört; die Entscheidung für Neu-Laden vom Host (R1/R2), `sudo` zu entziehen (S) oder die Restriktion außerhalb des Containers zu führen (R4) ([ADR-0025](0025-egress-liste-im-laufenden-container.md)); ein Verhalten von `nft` oder `dnsmasq`, das den geordneten Wechsel ändert.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-19 | Entscheidung des Projektinhabers im Gespräch: kurze Folge-ADR, die [ADR-0025](0025-egress-liste-im-laufenden-container.md) ergänzt; keine neue Messung (abgebrochen), Belege M2 und B-1; Entwurf, `Proposed`, `Accepted` setzt der Projektinhaber | Vereinbarung mit dem Projektinhaber |
| 2026-10-19 | `Accepted` durch den Projektinhaber („ok“, im Gespräch, auf die Vorlage der Freigabe-Punkte a bis d mit dem Hinweis, dass Punkt d wegfällt, wenn er nicht gewollt ist): (a) bis (d) wie empfohlen | Vereinbarung mit dem Projektinhaber |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes` auf [ADR-0027](0027-egress-init-wiederholungslauf-atomar.md).
