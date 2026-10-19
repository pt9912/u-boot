# Review-Report: slice-v1-egress-dns-upstream-schliessen — 2026-10-19

**Review-Art:** Code- und Design-Review (Sicherheits-Template `egress-init.sh.tmpl`, Unit- und Docker-Tests, Doctor-Meldung, Handbuch, READMEs, CHANGELOG, Messprotokoll, [`ADR-0026`](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md)) gegen Plan, Entscheidungen und Hard Rules; Template, Tests, Doku und Slice selbst gelesen, keine Einschätzung des Implementers ungeprüft übernommen.

**Gegenstand:** `git diff -M d2f8294..HEAD` (Phase 1 Messung: `eab3c2a`, `445fa99`, `1910089`; [`ADR-0026`](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md): `ba2e788`, `4524ffe`; Phase 2 Umsetzung: `f691fab`, `6d85215`, `c1e6c07`, `62c5f26`): 54 Dateien, 2572 Einfügungen, 46 Löschungen; davon 11 Dateien nach `4524ffe` (Phase 2).

**Skill:** `.harness/skills/reviewer.md` · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-19

**Eingangs-Kontext:**

- Slice [`slice-v1-egress-dns-upstream-schliessen`](../plan/planning/in-progress/slice-v1-egress-dns-upstream-schliessen.md) (§1, §2 DoD, §3 Plan, „Umsetzung (Implementer)“ Phase 1 und 2, §4, §6).
- [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes); [`ADR-0026`](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) (`Accepted`, Entscheidungen 1 bis 6), [`ADR-0012`](../plan/adr/0012-devcontainer-egress-firewall.md), [`ADR-0025`](../plan/adr/0025-egress-liste-im-laufenden-container.md); [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung); `AGENTS.md` §3.1, §3.2, §3.6, §3.7, §3.9, §3.10, §5 Regel 6; `harness/review.md`; `spec/architecture.md` Abschnitt Tests.
- Selbst gelesen: das ganze Template `egress-init.sh.tmpl` nach dem Diff und der Diff dagegen; `postStart`-Verkettung (`devcontainer_sandbox.go`); `egressProbe`, `SandboxStatusFailed`; `TestGenerateDevcontainer_SandboxEgress_UpstreamDNSOnlyForTheResolverUser`, `TestGenerateDevcontainer_SandboxStartScriptsWriteStatusFiles`; `TestE2E_SandboxDevcontainer_EgressUpstreamDNS` samt Sonden-Skripten, `generatedRunArgs`, `TestE2E_SandboxDevcontainer_StatusFiles_Egress`; Dockerfile-Template (`USER vscode`); `10-aufraeumen.sh`, `variante.sh`; Handbuch-Abschnitte zu Egress, Statusdatei-Tabelle, Fehlersuche Exit 11; Spezifikation (Degradationstabelle, [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)) nur gelesen.
- Sensor-Läufe des Reviewers (nur `make`): `make lint` (Exit 0), `make test` (Exit 0), `make docs-check` (524 Dateien, 0 Befunde), `make doc-trace` (83 Anforderungen, 0 Waisen), `make doc-planning` (0 Befunde), `make ablage-pruefung` (26 Fälle, 0 Fehlschläge). `make test-docker` nicht gefahren (Verifier-Aufgabe; der Slice-Bericht nennt einen vollen Lauf und zwei gezielte Läufe, nicht wiederholt). `make verify-depguard` nicht gefahren: keine Importregel berührt.

---

## Findings

Keine HIGH. Zwei MEDIUM, vier LOW, fünf INFO. Der Konflikt-Pfad ist nicht ausgelöst.

MEDIUM `internal/hexagon/application/devcontainer_fallback.go` · "hint:     \"Start the container with the generated runArgs (`--cap-add=NET_ADMIN`), or restrict egress at network or DNS level.\"," - Reparaturhinweis zu `failed` passt nur zu einem der zwei Gründe
Klasse: Meldung und Hinweis driften auseinander
Quelle: [`ADR-0026`](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) Entscheidung 2 (neuer Abbruchgrund `failed`), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion); Skill MEDIUM „Doku … driftet gegen das implementierte Verhalten“
Befund: Der Text zu `failed` nennt jetzt zwei Gründe (Capability fehlt unter `fail`; lokaler Resolver läuft nicht unter dem Benutzer der Regeln, jede Richtlinie), der `hint` desselben Eintrags gilt unverändert nur dem ersten. `hint` ist laut Struktur „the repair hint of the warning and error items“, also auch für den Fehler-Eintrag von `u-boot doctor`. Die Statusdatei trägt keinen Freitext, der Host kann die Ursache nicht unterscheiden.
Risiko: Ein Nutzer, dessen Image `nobody` nicht trägt, bekommt von `u-boot doctor` (Exit 11) den Hinweis, den Container mit `--cap-add=NET_ADMIN` zu starten, obwohl die Capability gewährt war; der Hinweis führt in die falsche Richtung, kein Test pinnt Text oder Hinweis des Egress-`failed`-Eintrags (`doctor_fallback_test.go` und `devcontainer_fallback_test.go` prüfen den Egress-Eintrag auf Schwere, nicht auf Wortlaut). Der Slice führt „kein Test pinnt den Wortlaut“ als gelesen, zieht daraus aber nur die Änderung des Satzes, nicht des Hinweises.
Verifizierbar: Review-only (Lesen von `egressProbe`); nach einer Änderung `make test`

MEDIUM `docs/user/benutzerhandbuch.md` · "Das Startscript bricht mit Exit-Code 11 ab und schreibt `failed` in die Statusdatei," - Folge des Abbruchs für die Kette der Startscripte und für den Zustand des Containers nicht ausgewiesen
Klasse: Handbuch benennt Folge einer neuen Bedingung nicht
Quelle: `AGENTS.md` §3.10 (Bedingung/Folge in Handbuch gegen Code), [`ADR-0026`](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) Entscheidung 2 („Die Kette `egress && …` im `postStartCommand` bricht dann ab“), [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) *Stand*
Befund: `devcontainer_sandbox.go` setzt `egress-init` mit ` && ` vor `docker-init`; endet es mit 11, läuft `docker-init.sh` nicht (die Spezifikation führt das als Verhalten bei `fail` und die Datei `daemon` vom vorigen Start bleibt stehen). Neu gilt das auch unter `warn`, wenn `nobody` fehlt oder die UID abweicht. Das Handbuch nennt den Abbruch und die Statusdatei, nicht, dass mit `nestedRuntime: docker` kein verschachtelter Daemon startet, und nicht, in welchem Zustand der Container bleibt: bei fehlendem Benutzer ohne jede Restriktion (Abbruch vor dem Anlegen der Tabellen, im Erstlauf offen), bei UID-Abweichung nach dem Anlegen mit Tabellen `drop` und unverändertem `resolv.conf` (kein DNS, fail-closed). Die Übergabe des Slice („Offen für Folge-Vorgänge“) nennt [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) und die Degradationszeile, nicht diese Kette und nicht den Zustand.
Risiko: Ein Nutzer mit `warn` und `nestedRuntime: docker` liest „Warnung statt Abbruch“ in der Tabelle zu `onUnavailable` (Handbuch, Abschnitt zu den Fähigkeiten) und findet den Container ohne Daemon und mit offenem Netz (erster Fall) oder ganz ohne DNS (zweiter Fall), ohne dass das Handbuch es sagt; der Spezifikations-Slice bekommt die Folge nicht als Eingang.
Verifizierbar: Review-only (Lesen von `postStart`-Verkettung und Handbuch); die Wirkung der Kette ist nicht durch einen Test gepinnt

LOW `internal/hexagon/application/generate_sandbox_test.go` · "// `failed` is written right before the exit code 11 of the `fail` policy." - Kommentar beschreibt die frühere Zusage
Klasse: Kommentar nennt abgelaufene Zusage
Quelle: `AGENTS.md` §3.9 (b), §3.7
Befund: Die Zählung dahinter wurde für `egress` von 1 auf 3 geändert; zwei der drei Zweige schreiben `failed` und enden mit 11 unabhängig von der Richtlinie. Der Kommentar vor der Prüfung behauptet weiter nur die Richtlinie `fail`; das Doc des Tests („`failed` before `exit 11`“) trifft zu.
Risiko: Der nächste Änderer liest „exit 11 der Richtlinie fail“ und hält die beiden unbedingten Abbruchzweige für einen Zählfehler; kein Gate widerlegt den Kommentar.
Verifizierbar: Review-only

LOW `internal/hexagon/application/templates/devcontainer/egress-init.sh.tmpl` · "echo \"u-boot sandbox: egress restriction failed: the local resolver is not running as '$dns_user' (uid $dns_uid)\" >&2" - Die Prüfung trennt „anderer Benutzer“ nicht von „nicht gestartet“
Klasse: Abbruchgrund breiter als Meldung
Quelle: [`ADR-0026`](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) Entscheidung 2 (b); `AGENTS.md` §3.10 (Bedingung im Handbuch: „wenn … der lokale DNS-Server nicht unter diesem Benutzer läuft“)
Befund: Der Zweig greift auch, wenn `dnsmasq` nicht startet (Konfigurationsfehler, Port belegt, Programm fehlt): Die Pid-Datei fehlt, `running_uid` ist leer, das Script meldet „not running as 'nobody'“. Gelesen: Das Script prüft den Rückgabewert von `dnsmasq --conf-file` nicht. Beim erneuten Lauf entsteht der Fall, wenn der alte Prozess nach `kill` (SIGTERM, ohne Warten) den Port noch hält, während der neue startet. Kein Test pinnt den Fall „startet nicht“; die D4-Messung (50 Wiederholungen) zeigt ihn nicht.
Risiko: Meldung und Doctor-Text benennen eine Ursache (anderer Benutzer), die der Zweig nicht prüft; wer ein Rennen beim Wiederholungslauf erlebt, sucht in der Benutzerfrage. Vorher endete der Fall still mit `ok` (kein neues Fenster, aber jetzt sichtbar als Exit 11).
Verifizierbar: Review-only; ein Pin bräuchte einen Docker-Test mit fehlschlagendem Start

LOW `docs/user/benutzerhandbuch.md` · "oder ein im Container fehlendes `/dev/fuse` unter `onUnavailable: fail`." - Fehlersuche zu Exit-Code 11 im Startscript kennt den neuen Grund nicht
Klasse: Handbuch-Nachführung unvollständig
Quelle: `AGENTS.md` §5 Regel 5 (öffentliche Verträge nachziehen); Skill MEDIUM-Muster „Doku driftet“ nur bei Wiederholung, hier LOW
Befund: Der Abschnitt „Fehler: Sandbox-Container startet nicht (Exit-Code 11 im Startscript)“ zählt Ursachen auf (blockierte User-Namespaces, `/dev/fuse` unter `fail`, Daemon ohne `--privileged`) und nennt den Abbruch der Egress-Beschränkung bei fehlendem Benutzer oder abweichender UID nicht; ebenso nennt die Tabelle zu `onUnavailable` (Zustand/`warn`/`fail`) keine Zeile dafür. Die Statusdatei-Tabelle trägt die UID-Abweichung, aber nicht den fehlenden Benutzer.
Risiko: Ein Nutzer mit Exit 11 aus dem Egress-Script findet in der Fehlersuche keinen Eintrag zu seinem Fall.
Verifizierbar: Review-only

LOW `internal/e2e/sandbox_devcontainer_docker_test.go` · "DIRECT udp example.com=none\", \"DIRECT tcp example.com=none\"," - Die TCP-Sonde hat keine positive Kontrolle
Klasse: Negativsonde ohne Gegenprobe
Quelle: `AGENTS.md` §3.10 (Handbuch-Bedingung „UDP und TCP“ braucht pinnenden Test); `harness/review.md` Linse Tests
Befund: Die Sonde `dnsq` wird nur für UDP positiv gesehen (`LIMIT sudo-as-resolver-user=answer`); `tcp … =none` kann ebenso aus einem Fehler der TCP-Rahmung der Perl-Sonde folgen. Die Mutation „`meta skuid` nur aus der `tcp`-Regel entfernt“ färbt laut Slice-Bericht nur den Unit-Test rot, nicht den Docker-Test.
Risiko: Ein Fehler in der TCP-Sonde lässt den Docker-Test grün; die Aussage „Docker-Test belegt UDP und TCP“ stützt sich für TCP auf die Teilstring-Prüfung des Unit-Tests. Gering: Der Unit-Test pinnt beide Regeln einzeln (Zählung, „Regel ohne Benutzerbedingung“).
Verifizierbar: `make test-docker` (Mutation, Verifier), sonst Review-only

LOW `docs/messungen/slice-v1-egress-dns-upstream-schliessen/README.md` · "eine\nUmsetzung im Template (keine; Golden Case, Docker-Test und Handbuch sind unberührt)" - Zustandsaussage des Protokolls gilt nach Phase 2 nicht mehr <!-- d-check:ignore (Zitat über Zeilenumbruch, Fundstelle ist die Zeile „Umsetzung im Template (keine;“) -->
Klasse: Zustandsmarker mit Ablaufpunkt bleibt stehen
Quelle: `AGENTS.md` §3.10 (Zustandsmarker beim Ablauf entfernen oder auf den Ist-Stand setzen)
Befund: Der Abschnitt „Nicht gemessen“ führt „eine Umsetzung im Template (keine …)“ und „die Wirkung auf den Go-Test selbst (`make test-docker` lief nicht …)“. Beides war zum Messzeitpunkt wahr; der Diff enthält die Umsetzung und den Lauf im selben Zweig. Der Abschnitt ist als Stand der Messung geschrieben, nennt aber keinen Messzeitpunkt für diese zwei Punkte.
Risiko: Ein Leser des Protokolls hält Template und Go-Test für ungemessen/unverändert; geringes Risiko, weil Slice und ADR den Ist-Stand tragen.
Verifizierbar: Review-only

INFO `docs/plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md` · "Wortlaut der Punkte 1 bis 6 steht zur Freigabe" - Zustandsmarker in einem `Accepted`-ADR
Klasse: Zustandsmarker mit Ablaufpunkt bleibt stehen
Quelle: `AGENTS.md` §3.5 (Immutabilität), §3.10
Befund: Die Entscheidung trägt „steht zur Freigabe“ und den Block „Freigabe-Punkte des Projektinhabers vor `Accepted`“; der Statuswechsel auf `Accepted` (`4524ffe`) änderte nur Status, Index und die Geschichte-Zeile. Geprüft: Wortlaut der Freigabe in der Geschichte-Zeile („Alle sieben (a bis g) wie empfohlen“, Punkte a bis g mit Entscheidungsnummern) deckt die Punkte; [`ADR-0012`](../plan/adr/0012-devcontainer-egress-firewall.md) und [`ADR-0025`](../plan/adr/0025-egress-liste-im-laufenden-container.md) sind im Diff unverändert (`git diff --stat` leer), der ADR-Index trägt Zeile 0026.
Risiko: Keines für den Vertrag; eine Berichtigung wäre jetzt eine Änderung eines `Accepted`-ADR. Nur gemeldet, kein Handlungsauftrag.
Verifizierbar: Review-only

INFO `internal/hexagon/application/templates/devcontainer/egress-init.sh.tmpl` · "dnsmasq --conf-file=\"$conf\"" - Zeitpunkt der UID-Prüfung hängt am Verhalten von `dnsmasq`
Klasse: Annahme über externes Programm
Quelle: [`ADR-0026`](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) Entscheidung 2 (b)
Befund: `dnsmasq` daemonisiert (kein `--keep-in-foreground`); die Prüfung liest `/proc/<pid>/status` unmittelbar nach der Rückkehr des Aufrufs und setzt voraus, dass Pid-Datei geschrieben und Benutzerwechsel vollzogen sind, wenn der Elternprozess endet. Diese Reihenfolge ist Verhalten von `dnsmasq` 2.91 (Image `devcontainers/base:debian`), nicht im Script erzwungen und nicht im Quelltext gelesen; belegt ist sie durch den Docker-Test (Pass), D4 (50 Wiederholungen ohne Auffälligkeit) und die UID-Gleichheit im Test. Ein anderer `dnsmasq` mit anderer Reihenfolge ergäbe einen Fehlalarm Exit 11, nie ein stilles `ok` (Prüfung schlägt bei leerer UID an). Busybox/dash: nur POSIX-Konstrukte (`$(…)`, `[ ]`, `awk`), `sh` ist der Aufruf (`sudo sh …`); kein `set -e`, kein Bashismus gelesen.
Risiko: Gering; fail-loud.
Verifizierbar: `make test-docker` (Verifier); Fehlalarm auf langsamen Hosts nicht gemessen (Messbericht: „Startrennen unter Last“ ausdrücklich nicht gemessen)

INFO `docs/user/benutzerhandbuch.md` · "etwa 36 Sekunden bis zur" - Messwerte ohne pinnenden Test sind teils ausgewiesen, teils nicht
Klasse: Ausweisung uneinheitlich
Quelle: `AGENTS.md` §3.10, [`ADR-0026`](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) Entscheidung 3 und 4
Befund: Die Grenze „3 650 bis 5 640 Anfragen je Sekunde“ trägt „ohne Test, der es pinnt“. Andere Zahlen und Aussagen (36 s statt 20 ms, 43 statt 7 s, `podman run --dns` löst nichts, Standard-DNS und `--network=host` bleiben nutzbar) stehen mit Messstand im Absatz, aber ohne diesen Zusatz; sie sind im Messprotokoll (zwei Läufe) belegt, nicht durch einen Test gepinnt, den der Docker-Test `…_NestedPodman_BuildEgress` trägt nur „Bau läuft durch“. Handbuch und ADR benennen den Messstand („nativer Linux-x86_64-Host, Docker 29.9.1, zwei Läufe“) und die ungemessenen Umgebungen (Colima, Docker Desktop, `arm64`, CI, `nestedRuntime: docker`, benutzerdefiniertes Netz, Podman-Netzwerkmodi, IPv6-Upstream) vollständig.
Risiko: Gering; die Zusagen sind als gemessen ausgewiesen, nicht als gepinnt.
Verifizierbar: Review-only

INFO `internal/e2e/sandbox_devcontainer_docker_test.go` · "const egressUpstreamDNSProbeScript" - Netzabhängigkeit und Skip-Bedingung
Klasse: Test hängt am echten Netz
Quelle: `spec/architecture.md` Abschnitt Tests (Build-Tag `docker`), `BEO-E2E/test-gruen-nur-mit-lokalem-cache`
Befund: Gelesen: Build-Tag `docker` steht, Image und Container räumen auf (`rmi -f` in `t.Cleanup`, `docker run --rm`), keine Uhr im Spiel. Der Test braucht Internet (HTTPS zu `api.github.com`, Auflösung über den echten Upstream; Wiederholungen mit `--retry 3`); offline ist er rot, nicht grün (die positive Kontrolle `LIMIT … =answer` und `HTTPS … =200` schlagen an). Skip nur bei `UPSTREAMS=0`; ein Upstream mit Loopback-Adresse (etwa `127.0.0.x` im benutzerdefinierten Netz) würde als offener direkter Pfad gemessen (`oifname lo accept` greift vor der Regel) und den Test rot machen statt zu überspringen. Der Test fährt als `vscode` (Dockerfile `USER vscode`), nicht als `root`. Das Mutations-Protokoll des Slice stimmt mit den Prüfungen überein: `skuid` aus beiden Regeln (`DIRECT … =answer`), `user=daemon` (`INIT_EXIT=11`), UID-Prüfung aus (`ABORT uid-differs exit=0`), Benutzerprüfung aus (`ABORT user-missing exit=0`); die `user=`-Zeile trägt nur der Unit-Test (vom Slice so ausgewiesen).
Risiko: Gering; Flakiness bei Netzstörung ist laut, nicht still.
Verifizierbar: `make test-docker` (Verifier)

INFO `docs/plan/planning/in-progress/slice-v1-egress-dns-upstream-schliessen.md` · "DoD Punkte 2, 3, 5, 7" - Zuständigkeit Verifier und Architect/Spezifikations-Slice <!-- d-check:ignore (Dateipfad des Slice bzw. Messordners als Fundstelle) -->
Klasse: Übergabe an andere Rolle
Quelle: `harness/roles.md`, [`ADR-0026`](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) Entscheidung 5
Befund: (1) Die DoD-Zeile „Test … abgenommen“ steht mit `[x]` und der Klammer „der Ausgang in §7 folgt bei der Closure“; Closure-Notiz, Register, Risiken §6 sind offen: Verifier-Aufgabe, kein Befund. (2) Die Spezifikation ist unverändert (`git diff --stat -- spec` leer, ebenso [`ADR-0012`](../plan/adr/0012-devcontainer-egress-firewall.md)/-0025), [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) und die Degradationszeile werden als Übergabe an den Spezifikations-Slice geführt, nicht geändert: korrekt. Ergänzend für diese Übergabe: Der Satz zu `failed` („bei der Richtlinie `fail`“, [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) Abschnitt Statusdatei), die Zeile „Kette `egress && docker`“ im Abschnitt *Stand* und die Egress-Capability-Zeile der Degradationstabelle (Exit 11 nur bei `fail`) sind jetzt alle drei vom Code überholt (siehe zweites MEDIUM). (3) Das Lastenheft ist im Diff unverändert.
Risiko: Die Übergabe nennt zwei der drei Stellen.
Verifizierbar: Review-only

---

## Negativbefunde

Geprüft, ohne Befund:

- **Template, Regelsyntax und Schleife:** `output meta skuid "$dns_uid" ip daddr "$ns" udp|tcp dport 53 accept` ist gültige nft-Syntax (Reihenfolge `meta skuid` vor `ip daddr`, wie die Messvariante, numerische UID, Variablen in doppelten Anführungszeichen); eine Regel je Upstream und Protokoll, IPv4-Upstreams (`awk`-Filter `$2 !~ /:/` unverändert), IPv6-Tabelle ohne Upstream-Regel; die Regeln `lo`, `established,related`, `@allowed` bleiben unverändert und vor den Upstream-Regeln.
- **Benutzer und Konfiguration:** `dns_user=nobody` einmal gesetzt; `dns_uid` per `id -u` (numerisch, wie gemessen); `user=$dns_user` und `meta skuid "$dns_uid"` lesen dieselbe Variable; Test zählt je Teilstring genau einen Treffer und zwei `dport 53 accept`, jede Zeile mit `dport 53 accept` trägt `meta skuid`.
- **Fehlerpfade und Reihenfolge:** Benutzerprüfung steht nach der Capability-Prüfung und vor Sicherung von `resolv.conf`, Löschen und Anlegen der Tabellen (Abbruch ohne Teilzustand; im Erstlauf bleibt der Container ohne Restriktion, siehe zweites MEDIUM); die Prozess-UID-Prüfung steht nach dem Start und vor dem Schreiben von `resolv.conf` (Test pinnt die Reihenfolge); Pid-Datei wird vor dem Start gelöscht, ein alter Eintrag besteht die Prüfung nicht; leere oder nicht lesbare `/proc`-Zeile führt zum Abbruch, nie zu `ok`. Abbruch nach dem Anlegen der Tabellen: `resolv.conf` unberührt, Tabellen mit `drop` bleiben, `failed` geschrieben, Exit 11 (fail-closed); im Wiederholungslauf bleiben ein früherer `resolv.conf`-Eintrag und die vorherigen Tabellen bis zum Löschen unverändert, das Fenster aus [`ADR-0025`](../plan/adr/0025-egress-liste-im-laufenden-container.md) wird durch den Diff nicht größer (Reihenfolge Kill, Löschen, Anlegen unverändert, ein `rm -f` zusätzlich).
- **Exit-Code- und Klassifikations-Linse:** Exit 11 im Container-Script ist ein Umgebungsproblem wie die übrigen Startscript-Abbrüche ([`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)); `u-boot doctor` wertet `failed` über den bestehenden Weg als Fehler mit Exit 11 aus (`TestE2E_SandboxDevcontainer_StatusFiles_Egress` pinnt `failed` bei `fail`, grün laut Slice-Bericht); der neue Weg geht über dieselbe Datei und denselben `fallbackProbe`. Kein neuer Sentinel: der Diff berührt weder `internal/adapter/driving/cli` noch Mapper, `ExitCode` oder Diagnostic-Abbildungen; `git diff --stat` bestätigt es; keine Dual-Classifier-Arbeit nötig. Der Kommentar zu `SandboxStatusFailed` ist im Indikativ und trifft zu.
- **Architektur und Importregeln:** Änderungen nur im Template (Daten), in Tests, in einem Meldungstext (`application`) und einem Kommentar (`domain`); keine Importänderung, kein Port, kein Adapter.
- **Unit-Tests:** `TestGenerateDevcontainer_SandboxEgress_UpstreamDNSOnlyForTheResolverUser` läuft für `warn` und `fail`; die Änderung der Zählung von `failed` (1 auf 3) und `exit 11` direkt danach (1 auf 3) entspricht den drei Zweigen im Template (Capability, Benutzer, UID); begründet in der Replay-Tabelle (`harness/replay.md`). Es gibt keine Whole-File-Golden-Datei zu `egress-init.sh`, nur Teilstring- und Reihenfolge-Pins.
- **Docker-Test:** Build-Tag `docker`, Cleanup des Images, `--rm` der Container, `context.WithTimeout`, englische Fehlermeldungen; positive Kontrollen (UDP-Sonde `answer` als `nobody` per `sudo setpriv`, `getent`, HTTPS 200, UID-Gleichheit, Regelzahl `2 × Upstreams`, `STATUS result=ok`); Abbruch-Fälle über `sed`-Kopien des ausgelieferten Scripts (`/usr/local/bin/u-boot-egress-init`), nicht über die Messvariante; Meldungstexte und `exit=11 result=failed` unter `warn` geprüft. Der rote Test `…_NestedPodman_BuildEgress` ist unverändert und wurde laut Bericht gegen das ausgelieferte Template gefahren (PASS in 62 s, nicht wiederholt).
- **Messvariante gegen Umsetzung:** Unterschiede (`user=nobody` gesetzt, UID aus `id -u`, Prüfung nach dem Start, Benutzerprüfung) sind im ADR als „ungemessen“ geführt (`user=nobody` in der Konfiguration, Name statt UID) und durch den Docker-Test mit dem ausgelieferten Script abgedeckt.
- **Messprotokoll:** `docs/messungen/slice-v1-egress-dns-upstream-schliessen/` enthält nur `.md` (1) und `.txt` (15), keine Binärdateien; `tools/messung/egress-dns-upstream/` 24 Dateien wie im Plan; `make ablage-pruefung` grün; Suche im Messverzeichnis, in den Skripten, im ADR und im Slice nach den Namen fremder Objekte ohne Treffer; `10-aufraeumen.sh` entfernt nur Objekte mit dem Präfix `u-boot-messung-dns` (`praefix_entfernen`), `rohdaten-11-nachher.txt` weist fremde Objekte nur als Anzahl aus; `docs/messungen/README.md` trägt die Bestandszeile. <!-- d-check:ignore (Dateipfad des Slice bzw. Messordners als Fundstelle) -->
- **[`ADR-0026`](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) und Spec-Straten:** Status, Index-Zeile, Geschichte mit Wortlaut der Freigabe; [`ADR-0012`](../plan/adr/0012-devcontainer-egress-firewall.md) und [`ADR-0025`](../plan/adr/0025-egress-liste-im-laufenden-container.md) unverändert; Lastenheft und Spezifikation unverändert; die ADR zeigt im `Schärft`-Feld aufwärts, Rang 1 bis 3 verweisen nicht abwärts (Diff berührt `spec/` nicht). Fitness Function: Unit-Test, Docker-Test mit der benannten Grenze und Wächter-Mutationen sind vorhanden und im Slice-Bericht je mit Beobachtung geführt.
- **Nutzer-Doku und CHANGELOG:** Handbuch, `README.md`, `README.de.md` kennungsfrei (`AGENTS.md` §5 Regel 6; Suche im Diff nach `LH-`, `ADR-`, `SPEC-` in diesen Dateien: keine); `CHANGELOG.md` englisch, Kennung verlinkt; Kapitelumbenennung „Namensauflösung und ihre Grenze“ mit allen Querverweisen nachgezogen (Suche im Baum nach „Grenze der Namensauflösung“ außerhalb von Planung, Messprotokoll und Reviews: keine Treffer); die README-Aussage zum `podman build` auf `amd64` (vorher: scheitert) ist auf den gemessenen Stand mit Messstand und Wartezeit umgestellt; der Messstand (nativer Linux-x86_64, Docker 29.9.1, zwei Läufe, 2026-10-19) und die ungemessenen Umgebungen sind genannt.
- **Kommentare (`AGENTS.md` §3.7, §3.9):** Template-Kommentare im Indikativ (Kopplung „`user=` … and `meta skuid` … read the same value“); ein Kommentar „a resolver under another uid would be blocked“ beschreibt die Folge, die die Prüfung abfängt (Grenzfall, kein Finding); `harness/replay.md`-Zeilen stimmen mit Tests überein; Symbole (`egressProbe`, `SandboxStatusFailed`, Testnamen) im Baum gesucht, Namen stimmen. Helfer stehen nicht zwischen Doc-Kommentar und Funktion. Kein `//nolint`, keine Änderung an `.golangci.yml`, Gates oder Workflows.
- **Plan-Tabelle gegen `git diff --stat -M`:** `d2f8294..HEAD`: 54 Dateien (24 Skripte, 16 Protokoll-Dateien, `docs/messungen/README.md`, ADR, ADR-Index, Slice, Handbuch, CHANGELOG, zwei READMEs, `harness/replay.md`, Docker-Test, Unit-Test, `devcontainer_fallback.go`, `sandbox_status.go`, Template) und `4524ffe..HEAD`: 11 Dateien; jede Datei steht in der Tabelle, jeder Eintrag hat einen Diff; als „unverändert“ geführte Dateien (`sandbox_status_docker_test.go`, Spezifikation, Lastenheft) fehlen im Diff.
- **Docker-only und Gates:** nur `make`-Sensoren; keine Gate-Lockerung; Coverage-Schwelle nicht berührt.

Nicht geprüft:

- `make test-docker` und die Wiederholung der Wächter-Mutationen: Verifier-Aufgabe; gelesen, nicht gefahren.
- Verhalten auf Colima, Docker Desktop, `arm64`, in der CI, mit `nestedRuntime: docker`, benutzerdefiniertem Netz und IPv6-Upstream: vom Slice und vom ADR als ungemessen ausgewiesen.
- Quelltext von `dnsmasq` (Reihenfolge Pid-Datei, Benutzerwechsel, Rückkehr des Elternprozesses): nicht gelesen, siehe INFO.
- `make gates`, `make verify-depguard`: Importregeln unberührt; die Einzelsensoren liefen grün.
- DoD-Closure, Closure-Notiz, Register, Ausgang der Risiken aus §6: Verifier-Aufgabe.

## Zusammenfassung

Keine HIGH; zwei MEDIUM, vier LOW, fünf INFO. Das Template setzt [`ADR-0026`](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) getreu um: Regel nach Benutzer mit einem Wert für `user=` und `meta skuid`, Abbruch vor den Tabellen bei fehlendem Benutzer, Prüfung der Prozess-UID nach dem Start und vor `resolv.conf`, fail-closed nach dem Anlegen der Tabellen; die Tests belegen Regelsyntax, Reihenfolge und Wirkung am ausgelieferten Script. Offen für den Implementer: Reparaturhinweis der Doctor-Meldung zu `failed` (erstes MEDIUM), im Handbuch und in der Übergabe an den Spezifikations-Slice die Folge des Abbruchs unter `warn` für die Startscript-Kette und den Containerzustand (zweites MEDIUM), der veraltete Test-Kommentar, die Fehlersuche im Handbuch, die TCP-Gegenprobe und der Zustandsmarker im Messprotokoll.
