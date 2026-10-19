# Verification-Report: slice-v1-egress-dns-upstream-schliessen — 2026-10-19

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“, gegen Plan, DoD und Entscheidungen). Kein Review, keine Validierung, keine Reparatur.

**Gegenstand:** Slice [`slice-v1-egress-dns-upstream-schliessen`](../plan/planning/done/slice-v1-egress-dns-upstream-schliessen.md) (`in-progress/`), Commit-Range `d2f8294..HEAD` (HEAD `5b44920`). Phase 1 (Messung): `eab3c2a`, `445fa99`, `1910089`. [`ADR-0026`](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md): `ba2e788`, `4524ffe`. Phase 2 (Umsetzung): `f691fab`, `6d85215`, `c1e6c07`, `62c5f26`. Review-Report `3ebae0b`. Fix-Runde: `52c7d1d`, `8d914a3`, `5b44920`. Bezug: [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`ADR-0012`](../plan/adr/0012-devcontainer-egress-firewall.md), [`ADR-0025`](../plan/adr/0025-egress-liste-im-laufenden-container.md).

**Eingang:** Implementer-Zusage (DoD, Plan-Tabelle, Sensor- und Mutationstabellen, Fix-Runde), Review [`2026-10-19-slice-v1-egress-dns-upstream-schliessen`](2026-10-19-slice-v1-egress-dns-upstream-schliessen.md) (0 HIGH, 2 MEDIUM, 4 LOW, 5 INFO).

**Gesamturteil: GO mit Auflagen.** Die Zusage trägt, und ich habe sie selbst nachgefahren: alle fünf Sensoren und `make gates` grün; das volle `make test-docker` grün (e2e-Paket 629,4 s); die gezielten Docker-Tests PASS ohne SKIP, darunter der bisher rote Test `TestE2E_SandboxDevcontainer_NestedPodman_BuildEgress`; die Mutationen am Template färben den Docker-Test beziehungsweise den Unit-Test rot; ein Wiederholungslauf des Startscripts unter Last endete in 260 von 260 Läufen mit Exit 0. Eine DoD-Verletzung (Zusage ohne Beleg) liegt nicht vor. Die Auflagen sind Closure-Arbeit des Planners, ein Folge-Slice-Anker für die Spezifikation fehlt noch (Auflage 1), und drei kleine Restfunde (Abschnitt 7) gehören in die Fix-Runde oder die Closure.

## 1. Sensoren (selbst gefahren, nur `make` beziehungsweise das Image des `make test-docker`-Ziels)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | pass, **Exit 0** | eigener Lauf auf `5b44920`, sauberer Arbeitsbaum |
| `make doc-trace` | pass | „83 Anforderung(en), 0 Waise(n).“ |
| `make docs-check` | pass | 525 Dateien, 0 Befunde (vor Anlage dieses Reports) |
| `make doc-planning` | pass | 525 Dateien, 0 Befunde |
| `make ablage-pruefung` | pass | „26 Faelle, 0 Fehlschlaege“ |
| `make doc-immutable RANGE="d2f8294..HEAD"` | pass (13 Commits, 0 Befunde) | zählt für Inline-Status-ADRs nicht (bekannte Lücke, [Beobachtung](../plan/planning/observations/BEO-DOK/doc-immutable-greift-nicht-bei-inline-status/observation.md)); deshalb Diff gelesen, Abschnitt 2 |
| `make test-docker` (voll) | pass, **Exit 0**, Dauer 641 s | `ok internal/e2e 629.362s`, `ok internal/hexagon/application 17.183s`, alle übrigen Pakete `ok`; das Ziel läuft ohne `-v`, ein SKIP wäre darin nicht sichtbar, deshalb Abschnitt 3 gezielt mit `-v` |
| gezielte Docker-Tests mit `-v` | pass | Abschnitt 3 |
| Mutationen in Wegwerf-Kopie | wie erwartet | Abschnitt 4 |
| `make verify-depguard` | nicht gelaufen | keine Importregel berührt (Diff: Template, Tests, ein Meldungstext in `application`, ein Kommentar in `domain`, Doku); `make lint` in `make gates` grün |

## 2. Unveränderte Vertragsebenen (selbst gelesen)

- `git diff --stat d2f8294..HEAD -- docs/plan/adr spec`: nur [`ADR-0026`](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) (neu, 124 Zeilen) und der ADR-Index (eine Zeile). Der Diff der übrigen ADR-Dateien ist leer: [`ADR-0012`](../plan/adr/0012-devcontainer-egress-firewall.md), [`ADR-0025`](../plan/adr/0025-egress-liste-im-laufenden-container.md) und alle weiteren `Accepted`-ADRs sind unverändert. Lastenheft und `spec/` sind unverändert (Diff leer).
- [`ADR-0026`](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) wurde `Proposed` angelegt (`ba2e788`) und mit `4524ffe` auf `Accepted` gesetzt; die Geschichte-Zeile trägt „Alle sieben (a bis g) wie empfohlen“ mit Punkten a bis g. Der Text „Wortlaut der Punkte 1 bis 6 steht zur Freigabe“ und der Block „Freigabe-Punkte … vor `Accepted`“ stehen in der angenommenen ADR (Review INFO 1): als Zustandsmarker veraltet, aber wegen der Immutabilität nicht mehr zu ändern; **begründet offen**.

## 3. Docker-Tests selbst gefahren

Aufruf wie das Ziel `test-docker` in `repo.mk` (Bau-Schritt `--target test-docker-tools`, Exit 0; dann im Image mit Host-Netz, Docker-Socket und `/tmp`): `go test -tags docker -v -count=1 -run '^TestE2E_SandboxDevcontainer_(EgressUpstreamDNS|NestedPodman_BuildEgress|StatusFiles_Egress|Egress|EgressSubdomains)$' ./internal/e2e/`.

| Test | Ergebnis | Dauer |
|---|---|---|
| `TestE2E_SandboxDevcontainer_Egress` | **PASS** | 13,73 s |
| `TestE2E_SandboxDevcontainer_EgressSubdomains` | **PASS** | 10,55 s |
| `TestE2E_SandboxDevcontainer_EgressUpstreamDNS` | **PASS** | 22,98 s |
| ↳ `direct_path_closed` | **PASS** | 12,98 s |
| ↳ `abort_branches` | **PASS** | 0,79 s |
| `TestE2E_SandboxDevcontainer_NestedPodman_BuildEgress` (bisher rot) | **PASS** | 58,21 s |
| `TestE2E_SandboxDevcontainer_StatusFiles_Egress` | **PASS** | 14,08 s |

Gesamt `ok … 119.554s`, Gesamtdauer des Aufrufs 126 s. **Kein SKIP, kein FAIL.** Der Test `…_NestedPodman_BuildEgress` ist laut `git diff` unverändert (nur Einfügungen im Diff der Datei) und ohne Anpassung grün; die Messung des Implementers (62 s) und meine (58 s) liegen zusammen. Das volle `make test-docker` lief danach noch einmal grün (Abschnitt 1). Nach den Läufen blieb kein Container, kein Image und kein Volume mit den Präfixen `uboot-e2e-`, `u-boot-verify-dns-` zurück; ich habe nur eigene Objekte (Präfix `u-boot-verify-dns-`) entfernt.

## 4. Unabhängige Wirksamkeit (Mutationen in Wegwerf-Kopien)

Wegwerf-Kopien des Repos (`git archive HEAD`) im Scratchpad, nicht im Arbeitsbaum (`git status` danach sauber); je eine Änderung am Template, per `diff` gegen HEAD bestätigt; gezielter Lauf von `TestE2E_SandboxDevcontainer_EgressUpstreamDNS` im Image des `test-docker`-Ziels.

| Mutation | Erwartung | Gesehen |
|---|---|---|
| `meta skuid` aus beiden Upstream-Regeln entfernt | Docker-Test FAIL | **FAIL**: alle vier `DIRECT … example.com=answer`, `output lacks "DIRECT udp example.com=none"`, `rules with meta skuid = 0, want 2 per upstream`; der Unit-Test `…_UpstreamDNSOnlyForTheResolverUser` ebenfalls rot („upstream DNS rule without a user condition“) |
| `meta skuid` nur aus der `tcp`-Regel entfernt | Docker-Test FAIL (Fix-Runde, Review LOW) | **FAIL**: `DIRECT tcp example.com=answer`, `DIRECT udp … =none`; `abort_branches` bleibt PASS: die positive TCP-Kontrolle (`LIMIT tcp … =answer`) und die TCP-Sonde tragen jetzt |
| UID-Prüfung nach dem Start abgeschaltet (`if false`) | `abort_branches` FAIL | **FAIL**: `ABORT uid-differs exit=0 result=ok`, `resolver-in-resolv.conf=1` (der Zustand der Gegenprobe aus der Messung) |
| Zeile `user=` entfernt | Docker-Test bleibt PASS (Verhalten gleich, `nobody` ist der Standard), Unit-Test FAIL | **bestätigt**: Docker-Test PASS (22,6 s), Unit-Test FAIL (`echo "user=$dns_user"` kommt 0-mal vor). Die `user=`-Zeile trägt nur der Unit-Test; vom Slice und der ADR so ausgewiesen |

## 5. Template-Review unabhängig (`egress-init.sh.tmpl`, ganz gelesen, 100 Zeilen)

- **Regel und Variable:** `dns_user=nobody` einmal gesetzt (Zeile 39), `dns_uid` per `id -u` (Zeile 40, numerisch, wie gemessen), `user=$dns_user` (Zeile 64) und `meta skuid "$dns_uid"` (Zeilen 79 und 80) lesen denselben Wert; je Upstream und Protokoll eine Regel; `lo`, `established,related` und `@allowed` stehen unverändert davor; die IPv6-Tabelle hat keine Upstream-Regel. Nur POSIX-Konstrukte, kein `set -e`, kein Bashismus (Aufruf `sudo sh …`).
- **Reihenfolge der Abbrüche:** Capability-Prüfung, dann Benutzerprüfung (vor `resolv.conf`-Sicherung, `kill`, `nft delete`, Anlegen der Tabellen), dann Tabellen, dann `dnsmasq`-Start mit Rückgabewert-Prüfung, dann UID-Prüfung über `/proc/<pid>/status`, erst danach `echo "nameserver 127.0.0.1" > /etc/resolv.conf` und `write_status ok`. Die Pid-Datei wird vor dem Start gelöscht (Zeile 55): ein alter Eintrag besteht die Prüfung nicht; eine leere UID führt zum Abbruch, nie zu `ok`.
- **Fail-open / fail-closed (gelesen und durch `abort_branches` gepinnt):**
  - *Benutzer fehlt:* Abbruch vor den Tabellen, Exit 11, `failed`. Im **Erstlauf** bleibt der Container **ohne Restriktion** (`STATE user-missing tables=0`): fail-open mit Signal (Statusdatei `failed`, Meldung auf `stderr`, `u-boot doctor` zeigt `failed`). Im Wiederholungslauf bleiben die Tabellen des vorigen Laufs stehen.
  - *`dnsmasq` startet nicht / UID weicht ab:* Tabellen mit `policy drop` stehen, `/etc/resolv.conf` unverändert, `failed`, Exit 11: fail-closed (kein DNS, keine offene Verbindung). Die Vorgabe der [`ADR-0026`](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) Entscheidung 2 („kein Zustand, in dem `resolv.conf` auf einen nicht antwortenden `dnsmasq` zeigt, ohne dass die Statusdatei `failed` trägt“) ist erfüllt: der Docker-Test pinnt `resolver-in-resolv.conf=0`.
  - *Kette:* `generate_sandbox_test.go` pinnt `postStartCommand` als `… u-boot-egress-init && … u-boot-docker-init`; endet das Egress-Script mit 11, läuft das Docker-Script nicht (auch unter `warn`). Das Handbuch weist das aus (Abschnitt 6) und nennt es „gelesen, kein Test mit Abbruch“; die Verkettung selbst ist durch den Unit-Test gepinnt (Verhalten der Shell-Kette ist nicht ausgeführt). Ob die Dev-Container-Werkzeuge den Container bei fehlgeschlagenem `postStartCommand` weiterlaufen lassen, ist nicht gemessen und nicht behauptet.
- **Race der UID-Prüfung:** Die Prüfung liest `/proc/<pid>/status` unmittelbar nach der Rückkehr von `dnsmasq`; sie setzt voraus, dass Pid-Datei und Benutzerwechsel beim Ende des Elternprozesses vollzogen sind (Verhalten von `dnsmasq` 2.91, nicht im Quelltext gelesen). Ein Fehlalarm wäre fail-loud (Exit 11), nie ein stilles `ok`. Gemessen: Abschnitt 5a.
- **Doctor-Meldung:** `egressProbe` nennt im Text drei Gründe (Capability unter `fail`; fehlender Benutzer `nobody`; Resolver startet nicht oder läuft unter anderem Benutzer) und im Hinweis beide Wege (Capability, oder Meldung des Startscripts im Container-Log); `TestDoctor_SandboxFallback_EgressFailedNamesEveryCause` pinnt beides (in `make gates` grün). Review MEDIUM 1 ist damit behoben.

### 5a. Wiederholungslauf unter Last (selbst gemessen)

Wegwerf-Container (Präfix `u-boot-verify-dns-`) aus dem Image, das `u-boot generate devcontainer` erzeugt (Sandbox-Profil, Egress aktiv, `onUnavailable: warn`, ohne `nestedRuntime`; Basis wie im Slice), Skripte im Scratchpad, nicht im Repo. Der Erstlauf endete mit Exit 0; danach erneute Läufe `sudo sh /usr/local/bin/u-boot-egress-init` als `vscode`:

| Fall | Läufe | Exit 0 | Exit ≠ 0 | Zustand danach |
|---|---|---|---|---|
| je ein `docker exec`, ohne Last | 30 | 30 | 0 | ein `dnsmasq`, `result=ok`, 7 `accept`-Regeln |
| je ein `docker exec`, 40 Busy-Loops (davon 9 über 50 % CPU) im Container | 30 | 30 | 0 | ein `dnsmasq`, `result=ok`, 7 `accept`-Regeln |
| 100 Läufe hintereinander im Container (ohne `docker exec` dazwischen), CPU-Quote 1, ohne Last | 100 | 100 | 0 | ein `dnsmasq`, `result=ok` |
| dasselbe, 20 Busy-Loops, CPU-Quote 1 | 100 | 100 | 0 | ein `dnsmasq`, `result=ok` |

**Der Fall „alter `dnsmasq` hält den Port noch, neuer Start scheitert, Exit 11“ trat in 260 Läufen nicht ein.** Das Risiko aus dem Slice („Wiederholungslauf unter Last nicht gemessen“) ist damit für diesen Messstand (nativer Linux-x86_64-Host, Docker 29.9.1, ein Host, ein Image, die Last war ein Busy-Loop-Modell und nicht der Last eines echten Builds) **nicht eingetreten**; es bleibt eine Zusage für diesen Messstand, kein Beweis. Nicht gemessen: Colima, Docker Desktop, `arm64`, CI, Last durch parallelen Image-Bau, andere `dnsmasq`-Versionen.

## 6. Doku, Vertrag an vier Stellen

| Stelle | Befund |
|---|---|
| Code (Template, Doctor-Text) | Abschnitt 5 |
| Handbuch (`docs/user/benutzerhandbuch.md`) | Kennungsfrei (Suche im Diff nach `LH-`, `ADR-`, `SPEC-`, Slice-/Wellen-Kennungen und `spec/`: keine Treffer, `AGENTS.md` §5 Regel 6). Ist-Zustand; Messstand „nativer Linux-x86_64-Host, Docker 29.9.1, zwei Läufe“, Datum steht im Messprotokoll und im Slice, nicht im Handbuch (kein Befund). Messwerte 36 s statt 20 ms, 43 statt 7 s, `--dns`, Standard-Namensauflösung und `--network=host` stehen als „gemessen; kein Test pinnt sie“; die 3 650 bis 5 640 Anfragen je Sekunde mit „ohne Test, der es pinnt“. Grenzen: `sudo`, Adressen erlaubter Namen, Unternamen; ungemessen und nicht zugesagt: Colima, Docker Desktop, `arm64`, CI, `nestedRuntime: docker`, benutzerdefiniertes Netz, Podman-Netzwerkmodi, IPv6-Upstream. Neue Absätze zu Abbruch, den zwei Zuständen und der Daemon-Kette sind da; Fehlersuche (Exit 11), `onUnavailable`-Tabelle und Statusdatei-Tabelle nennen den neuen Grund (Review LOW 3 behoben). Querverweise auf „Namensauflösung und ihre Grenze“ nachgezogen; „Grenze der Namensauflösung“ steht außerhalb von Planung, Messprotokoll und Reviews nirgends mehr (Suche im Baum). |
| `README.md`, `README.de.md` | je ein Satz zum verschachtelten `podman build` auf den gemessenen Stand mit Wartezeit und Messstand umgestellt; kennungsfrei |
| `CHANGELOG.md` `[Unreleased]` | Englisch, `Changed`, Kennung verlinkt; nennt Regel, 36 s, 43 gegen 7 s, `--dns`, Exit 11 auch unter `warn`, Offenes, Messstand, Neuaufbau des Containers |
| `harness/replay.md` | zwei Zeilen ergänzt; Zählung der `failed`-Zweige veraltet, Abschnitt 7 (V1) |
| Doctor-Text | stimmt mit Template und Handbuch überein (drei Gründe) |

Die vier Stellen (Code, Handbuch, CHANGELOG, Doctor-Text) sagen dasselbe: Regel nur für `nobody`, Abbruch mit Exit 11 auch unter `warn`, Grenze, Messstand. Abweichung: Abschnitt 7 (V2, Zustand nach Abbruch im Wiederholungslauf).

## 7. Restfunde dieser Verifikation (keine HIGH)

| # | Schwere | Fund | Beleg | Vorschlag |
|---|---|---|---|---|
| V1 | LOW | Die Zählung der `failed`-Zweige ist nach der Fix-Runde an drei Stellen veraltet. Das Template hat **vier** (Capability unter `fail`, Benutzer fehlt, `dnsmasq` startet nicht, UID weicht ab); der Unit-Test zählt `failed: 4` (richtig). `harness/replay.md` nennt „drei `failed`-Zweige“ ohne „startet nicht“ und in der ersten Zeile nur „fehlender Benutzer und abweichende Prozess-UID“; der Slice-Abschnitt Phase 2 trägt „`failed` 1 → 3“; der Test-Kommentar in `generate_sandbox_test.go` sagt „the two resolver branches“, es sind drei. | `egress-init.sh.tmpl` Z. 28, 43, 88, 95; `generate_sandbox_test.go` (Zählung 4); Diff `harness/replay.md` | Implementer: Replay-Tabelle, Test-Kommentar; Planner: Zählung in der Phase-2-Notiz bei der Closure auf 4 |
| V2 | LOW | Das Handbuch sagt für „Resolver startet nicht oder läuft unter anderem Benutzer“, `/etc/resolv.conf` zeige „nicht auf den lokalen DNS-Server“. Das gilt nur im **Erstlauf**. Im Wiederholungslauf bleibt der Eintrag `127.0.0.1` des vorigen Laufs stehen (das Script schreibt `resolv.conf` erst am Ende), ohne dass ein Resolver antwortet; die Folge (keine Namensauflösung, keine offene Verbindung) und `failed` sind gleich. Der Docker-Test pinnt nur den Erstlauf. | Template Z. 98 (einzige Schreibstelle); `abort_branches` | Implementer: Handbuchsatz auf „im ersten Lauf“ einschränken oder beide Fälle nennen |
| V3 | INFO | Die Rest-Aussagen zu „Plan-Tabelle … nicht gelesen auf Berührung“ (Zeile zu `open/README.md` und Roadmap) genügen `AGENTS.md` §3.10 nicht: Roadmap Z. 20 trägt „Umsetzung erst nach Messbefund und Entscheidung“ (überholt: die Umsetzung liegt vor), `open/README.md` Z. 23 trägt „wenn die Messung einen Weg trägt … Messung zuerst“. | gelesen | Planner bei der Closure (Auflage 4) |
| V4 | INFO | Der Slice [`slice-v1-egress-init-ohne-luecke`](../plan/planning/done/slice-v1-egress-init-ohne-luecke.md) nennt Zeilen des Templates (40 bis 43, 57 ff., 40 bis 72), die nach diesem Diff um bis zu 26 verschoben sind (Template 100 statt 74 Zeilen, `kill` jetzt Z. 54). | gelesen | Planner: Zeilenbezüge neu lesen, bevor der Slice startet (`AGENTS.md` §3.10) |
| V5 | INFO | `user=` ist nur durch den Unit-Test gepinnt (Abschnitt 4). | Mutation | keine Handlung; vom Slice und der ADR so ausgewiesen |

## 8. Review-Findings einzeln

| Finding | Stand | Beleg |
|---|---|---|
| MEDIUM 1: Doctor-Hinweis zu `failed` nur für einen Grund | **behoben** | `52c7d1d`, Text und Hinweis nennen drei Gründe und beide Wege; `TestDoctor_SandboxFallback_EgressFailedNamesEveryCause`; `make gates` grün |
| MEDIUM 2: Folge des Abbruchs (Kette, Zustand) nicht ausgewiesen | **behoben** (im Handbuch), mit V2 als Einschränkung | `8d914a3`: Kette `egress` vor Daemon-Start, zwei Zustände, ein Docker-Test pinnt sie; der Spezifikations-Übergabe fehlte die Zeile „Kette `egress && docker`“ (Slice trägt sie jetzt) |
| LOW 1: Test-Kommentar nennt abgelaufene Zusage | **behoben**, mit V1 | der Kommentar nennt jetzt „every exit code 11“ und „the two resolver branches“ (richtig wären drei) |
| LOW 2: „nicht gestartet“ von „anderer Benutzer“ nicht getrennt | **behoben** | eigener Zweig mit eigener Meldung (`did not start`), Docker-Test pinnt Text und `STATE not-started`; Mutation „Rückgabewert-Prüfung aus“ rot laut Slice |
| LOW 3: Fehlersuche im Handbuch | **behoben** | Diff: Fehlersuche, `onUnavailable`-Tabelle, Statusdatei-Tabelle |
| LOW 4: TCP-Sonde ohne positive Kontrolle | **behoben** | `LIMIT tcp sudo-as-resolver-user=answer` ist Pflichtzeile; tcp-only-Mutation färbt den Docker-Test rot (Abschnitt 4) |
| LOW 5: Zustandsmarker im Messprotokoll | **behoben** | Abschnitt *Nicht gemessen* auf den Ist-Stand gesetzt (gelesen) |
| INFO 1: „steht zur Freigabe“ in `Accepted`-ADR | **begründet offen** | Immutabilität, Abschnitt 2 |
| INFO 2: UID-Prüfung hängt am `dnsmasq`-Verhalten | **begründet offen**, jetzt gemessen | Abschnitt 5a, 260 Läufe ohne Fehlalarm |
| INFO 3: Messwerte ohne pinnenden Test | **behoben** | Handbuch: „gemessen; kein Test pinnt sie“ |
| INFO 4: Test hängt am echten Netz | **begründet offen** | wie im Review; laut, nicht still; offen im Register der E2E-Beobachtungen |
| INFO 5: Übergabe Verifier / Spezifikation | **übergeben** | Auflagen 1 bis 3 |

## 9. DoD, Punkt für Punkt

| DoD-Punkt (Slice §2) | Befund |
|---|---|
| Messung D1 bis D5 mit Messstand, Skripten, Rohdaten (`[x]`) | **erfüllt, bestätigt.** Verzeichnis das Messprotokoll-Verzeichnis unter `docs/messungen/` hat nur `.md` (1) und `.txt` (15), `tools/messung/egress-dns-upstream/` 24 Dateien; Suche im Messprotokoll, in den Skripten, im ADR, im Slice, im Handbuch, in den READMEs, im CHANGELOG und in `harness/` nach den Namen fremder Docker-Objekte: keine Treffer; `vergleich.sh` Exit 0 („gleich: 180 Aussagen“), ändert nichts im Bestand; `10-aufraeumen.sh` arbeitet nur auf dem Präfix `u-boot-messung-dns-` (nicht ausgeführt: es entstand kein Messbestand). `make ablage-pruefung` grün. |
| Entscheidungs-ADR und Umsetzung (`[x]`) | **erfüllt, bestätigt.** [`ADR-0026`](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) `Accepted` (2026-10-19); Entscheidungen 1 bis 6 und Fitness Function gegen die Umsetzung: Abschnitt 10. |
| Test `…_NestedPodman_BuildEgress` abgenommen (`[ ]`) | **von mir bestätigt**: unverändert, im gezielten Lauf PASS in 58,21 s ohne SKIP, im vollen `make test-docker` grün. Das Häkchen und der Ausgang in §7 setzt der Planner. Rest-Hinweis: Der Bau dauert mit der Regel etwa 43 statt 7 s (gemessen, im Handbuch ausgewiesen); der Test hat 25 Minuten Frist. |
| Spec-Stelle benannt (`[ ]`) | **benannt, nicht verankert.** Der Slice nennt [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) (`failed` nur bei `fail`, `spec/spezifikation.md` Z. 642), die Capability-Zeile der Degradationstabelle (Exit 11 nur bei `fail`, Z. 201) und die Zeile „Kette `egress && docker`“ im Abschnitt *Stand* (Z. 655); ich habe alle drei gelesen, sie sind vom Code überholt (Abbruch auch unter `warn`). **Ein Folge-Slice oder Eintrag, der sie trägt, existiert nicht** (Suche in `open/`, `next/`, `in-progress/`, Roadmap): Auflage 1. |
| `make gates` grün; `make test-docker` grün (`[x]`) | **erfüllt, bestätigt** (Abschnitt 1) |
| Review durchgeführt (`[ ]`) | Report liegt vor (`3ebae0b`); Häkchen setzt der Planner |
| Doku-Update (`[x]`) | **erfüllt**, mit V1 und V2 (Abschnitt 6) |
| Closure-Notiz, Register, Risiken §6, drei Paarungen (`[ ]`) | **offen**, Planner (Auflagen 2 bis 5) |

## 10. ADR-0026: Entscheidungen 1 bis 6 und Fitness Function gegen die Umsetzung

| Punkt | Konformität |
|---|---|
| 1 Regel gilt (Upstream-DNS nur für den Benutzer des `dnsmasq`) | **erfüllt**: Template Z. 79 und 80, Unit- und Docker-Test, Mutationen (Abschnitt 4) |
| 2 Benutzer festgelegt, Regel aus demselben Wert, Abweichung bricht ab | **erfüllt**: `nobody`, `user=`, numerische UID (die ADR empfahl die numerische Form), Abbruch vor den Tabellen und nach dem Start, Exit 11 und `failed` unter jeder Richtlinie; zusätzlich der Zweig „`dnsmasq` startet nicht“ (Fix-Runde): eine Ergänzung im Sinn der Entscheidung, keine Abweichung |
| 3 Grenze ausgewiesen, nicht zugesagt | **erfüllt**: Handbuch (Abschnitt 6); der Docker-Test pinnt die `sudo`-Grenze (`LIMIT … =answer`, UDP und TCP) |
| 4 Verhalten im verschachtelten Container akzeptiert, im Handbuch | **erfüllt**: 36 s, 43 gegen 7 s, `--dns`, mit Messstand und „kein Test pinnt“ |
| 5 Lastenheft und Spezifikation unverändert, Frage an den Spezifikations-Slice | **erfüllt** (Diff leer); der Anker für den Spezifikations-Slice fehlt (Auflage 1) |
| 6 Der Test wird grün und bleibt unverändert | **erfüllt** (Abschnitt 3) |
| Fitness Function: Unit-/Golden-Test, Integrationstest, Wächter-Mutationen, Messprotokoll | **erfüllt**: Unit-Test (Abschnitt 4), Docker-Test (Abschnitt 3), Mutationen (i) bis (iii) von mir wiederholt, Messprotokoll im Repo; die Zeile „bestehende Konfigurationen ohne Egress rendern byte-identisch“ ist durch die bestehenden `…_NotGenerated`-Tests gedeckt (grün in `make test`, nicht einzeln neu gefahren) |

## 11. Plan-vs-Diff (`AGENTS.md` §3.10), nach dem letzten Commit

`git diff --stat -M d2f8294..HEAD` (vor diesem Report): **56 Dateien**, 2790 Einfügungen, 50 Löschungen. Gegen die Plan-Tabelle (§3):

| Plan-Zeile | Diff | Befund |
|---|---|---|
| Messprotokoll (16 Dateien) und Skripte (24 Dateien) | im Diff | stimmt (16 + 24) |
| `docs/messungen/README.md` | im Diff (1 Zeile) | stimmt |
| [`ADR-0026`](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) und ADR-Index | im Diff | stimmt |
| `egress-init.sh.tmpl` | im Diff | stimmt |
| `generate_sandbox_test.go`, `sandbox_devcontainer_docker_test.go` | im Diff | stimmt |
| `devcontainer_fallback.go`, `sandbox_status.go`, `devcontainer_fallback_test.go` | im Diff | stimmt |
| `harness/replay.md` | im Diff | stimmt (Inhalt: V1) |
| Review-Report | im Diff | stimmt |
| Handbuch, `CHANGELOG.md`, `README.md`, `README.de.md` | im Diff | stimmt |
| diese Slice-Datei | im Diff | stimmt |
| `README.md` des Messprotokolls (Fix-Runde) | in den 16 Dateien des Messprotokolls | stimmt |
| `internal/e2e/sandbox_status_docker_test.go`, `spec/`, Lastenheft, `docs/plan/adr/` ohne die neue ADR (als unverändert geführt) | nicht im Diff | stimmt |
| `open/README.md`, Roadmap („unverändert, nicht gelesen auf Berührung“) | nicht im Diff | stimmt für den Stand vor der Closure; die Aussage ist nicht gelesen (V3) |

Jede Diff-Zeile hat einen Plan-Eintrag, jeder Eintrag mit Diff eine Zeile. Mit diesem Report kommt als 57. Datei `docs/reviews/…` hinzu (durch die Sammelzeile *docs/reviews/* gedeckt); der Endstand ist nach der Closure-Bewegung erneut abzugleichen.

## 12. Kommentare, Suppression, Build-Tag

- `//go:build docker` steht im Docker-Test; Suche im Diff nach `nolint`: kein Treffer im Code (nur eine Erwähnung im Review-Report); keine Änderung an `.golangci.yml`, Gates oder Workflows.
- Kommentare im Indikativ mit Kopplung, Zusage oder Abgrenzung; Suche im Diff nach „former“, „früher“, „originally“, „no longer“, „is gone“: keine Treffer in Code und Skripten (`AGENTS.md` §3.7, §3.9 (a)). Der Test-Kommentar zur Zählung ist in der Zahl falsch (V1). Symbole in Kommentaren (`egressProbe`, `SandboxStatusFailed`, Testnamen) stimmen mit dem Baum.
- Fehlermeldungen von Template und Tests sind Englisch; CHANGELOG Englisch; Handbuch und READMEs Deutsch beziehungsweise nach Dateisprache.

## 13. Carveouts

Neu: keine. Gelöst: keine. Unverändert: alle. Kein Eintrag im Carveout-Inventar berührt.

## 14. Auflagen (Planner)

1. **Folge-Slice für die Spezifikation anlegen** (`open/`) und in §7 verlinken: [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) (`failed` auch bei `warn` für den Egress-Abbruch), die Capability-Zeile der Degradationstabelle (Exit 11 nur bei `fail`), die Zeile „Kette `egress && docker`“ im Abschnitt *Stand*, und die Frage aus [`ADR-0026`](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) Entscheidung 5, ob die Zusage samt Grenze eine eigene Spezifikationsstelle braucht. Ohne Anker ist DoD-Punkt 4 („benannt“) nur halb getragen.
2. **DoD-Rest und Closure:** die Häkchen zu „Test abgenommen“ (mit diesem Report als Beleg) und „Review durchgeführt“ setzen; §7 füllen (Steering-Loop-Lerneintrag, `liegt in` falls eine Regel verkörpert wird); Register fortschreiben oder „keine Beobachtung“ notieren (Kandidaten: Verhalten der `postStartCommand`-Kette bei Abbruch ungemessen; die `user=`-Zeile nur im Unit-Test gepinnt); die drei Paarungen tragen.
3. **Risiko-Ausgänge §6 setzen**: (a) „Regel trifft auch den `dnsmasq` oder erlaubte Namen“: **nicht eingetreten** (Messung D2, Docker-Test, voller Lauf); (b) „nur solange ohne `sudo`“: **weiter offen**, ausgewiesen im Handbuch und in der ADR (Guardrail); (c) „Datenkanal über erlaubte Namen“: **weiter offen**, gemessen bis zum Upstream (D5), ausgewiesen; (d) „Messung gilt nur für Colima/arm64“: gegenstandslos, der Messstand ist ein anderer (nativer Linux-x86_64); Colima, Docker Desktop, `arm64`, CI und `nestedRuntime: docker` **weiter offen**; (e) „Startrennen“: Wiederholungslauf unter Last in 260 von 260 Läufen ohne Abbruch (Abschnitt 5a), **nicht eingetreten für diesen Messstand**; das „Risiko, neu“ aus der Fix-Runde (Port noch belegt) auf denselben Beleg setzen.
4. **Zustandsmarker beim Ablauf berichtigen (V3):** Roadmap Z. 20 („Umsetzung erst nach Messbefund und Entscheidung“) und `open/README.md` Z. 23 bei der Closure auf den Ist-Stand setzen oder entfernen und die Plan-Aussage „nicht gelesen auf Berührung“ durch eine gelesene Aussage ersetzen; Plan-Tabelle nach der Closure-Bewegung gegen `git diff --stat` erneut abgleichen (57 Dateien mit diesem Report).
5. **Zählung der `failed`-Zweige (V1)** in der Phase-2-Notiz des Slice auf vier setzen; `harness/replay.md` und den Test-Kommentar lässt der Implementer berichtigen (kleine Fix-Runde vor der Closure) — beides vor dem Move nach `done/`.
6. **Handbuchsatz (V2)** zum Zustand nach dem Abbruch im Wiederholungslauf berichtigen (Implementer, gleiche kleine Fix-Runde).
7. **Folge-Slice [`slice-v1-egress-init-ohne-luecke`](../plan/planning/done/slice-v1-egress-init-ohne-luecke.md) (V4):** Zeilenbezüge des Templates beim Start neu lesen.
8. **Reihenfolge der Bewegung:** erst Inhalt (DoD-Häkchen, Closure-Notiz), dann reiner `git mv` nach `done/` (`AGENTS.md` §3.3); der Done-Slice trägt den Commit-Hash.

## 15. Nicht ausgeführt

- `make verify-depguard`: keine Importregel berührt.
- `make doc-immutable` mit Aussagekraft für die neue ADR: greift bei Inline-Status-ADRs nicht; Diff selbst gelesen (Abschnitt 2).
- `make image-scan`, `make govulncheck`, `make fullbuild`: kein Release-/Security-Pfad berührt.
- Colima, Docker Desktop, `arm64`, CI, `nestedRuntime: docker`, benutzerdefiniertes Docker-Netz, Podman-Netzwerkmodi außer dem Standard: vom Slice und von der ADR als ungemessen ausgewiesen; nicht von mir gemessen.
- Das Verhalten der Dev-Container-Werkzeuge bei fehlgeschlagenem `postStartCommand` (läuft der Container ohne Restriktion weiter?): nicht gemessen.
- `10-aufraeumen.sh` und die Messskripte der Phase 1: nicht neu gefahren (nur `vergleich.sh` auf den eingecheckten Dateien).

## Commit / Artefakt

Verifiziert auf `5b44920`. Dieser Report liegt unter `docs/reviews/` (Präfix `2026-10-19-verification-`).
