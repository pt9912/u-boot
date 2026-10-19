# Verification-Report: slice-v1-egress-hinweise-neuaufbau-und-messstand — 2026-10-19

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“, gegen Plan, DoD und Entscheidungen). Kein Review, keine Validierung, keine Reparatur.

**Gegenstand:** Slice [`slice-v1-egress-hinweise-neuaufbau-und-messstand`](../plan/planning/in-progress/slice-v1-egress-hinweise-neuaufbau-und-messstand.md) (`in-progress/`), Commit-Range `1e445d3..HEAD` (HEAD `e78d71f`, lokal). Implementierung: `c772352` (Hinweis und Texte), `31a55b6` (Spezifikation, Handbuch, CHANGELOG, Replay-Tabelle), `ee092eb` (Plan-Tabelle). Review-Report `c5be3c7`. Fix-Runde: `dccfa44` (Wildcard-Fehler nennt die Tiefe, Test für erste Aktivierung und Abschalten), `12b3dca` (Spezifikation, Handbuch, CHANGELOG nennen, wann der Hinweis nicht erscheint), `e78d71f` (Plan-Tabelle). `git diff --stat -M 1e445d3..HEAD`: 15 Dateien, 434 Einfügungen, 25 Löschungen; eine neue Testdatei, ein neuer Review-Report.

**Eingang:** Implementer-Zusage (DoD, Plan-Tabelle, Fix-Runde), Review [`2026-10-19-slice-v1-egress-hinweise-neuaufbau-und-messstand`](2026-10-19-slice-v1-egress-hinweise-neuaufbau-und-messstand.md) (0 HIGH, 0 MEDIUM, 2 LOW, 4 INFO), Entscheidungen des Projektinhabers vom 2026-10-19: (a) Subdomain-Texte auf das Gemessene kürzen, (b) der Neuaufbau-Hinweis bei `generate devcontainer` entfällt unter `--quiet`.

**Gesamturteil: GO mit Auflagen.** Der Hinweis erscheint in meiner End-zu-End-Probe genau im beschriebenen Fall und in keinem Gegenfall; der Docker-Test, an dem die gekürzten Subdomain-Texte hängen, ist von mir gefahren und grün; `make gates` und alle Doku-Sensoren sind grün; Lastenheft, ADRs, Workflows, `.d-check.yml` und Templates sind unverändert. Eine DoD-Verletzung (Zusage ohne Beleg) liegt nicht vor. Auflagen: Closure-Arbeit des Planners, ein LOW (Zusage in Spezifikation, Handbuch und CHANGELOG ohne pinnenden Test) und eine stale Zeilenangabe in der Plan-Tabelle.

## 1. Sensoren (selbst gefahren, nur `make` beziehungsweise das Image des `make test-docker`-Ziels)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | pass, **Exit 0** | `coverage-gate: OK — coverage 91.40% meets threshold 90%`; Kette wie in `harness/README.md` |
| `make doc-trace` | pass | „83 Anforderung(en), 0 Waise(n).“ |
| `make docs-check` | pass | 511 Dateien, 0 Befunde (vor Anlage dieses Reports) |
| `make doc-planning` | pass | 511 Dateien, 0 Befunde |
| `make ablage-pruefung` | pass | Wächter-Test „26 Faelle, 0 Fehlschlaege“ |
| `make doc-immutable RANGE="1e445d3..HEAD"` | 0 Befunde, **nicht als Beleg gezählt** | Gate schützt 13 Inline-Status-ADRs nicht ([`BEO-DOK/doc-immutable-greift-nicht-bei-inline-status`](../plan/planning/observations/BEO-DOK/doc-immutable-greift-nicht-bei-inline-status/observation.md)); Beleg stattdessen: `git diff --stat 1e445d3..HEAD -- docs/plan/adr` leer |
| `TestE2E_SandboxDevcontainer_EgressSubdomains` (Build-Tag `docker`) | **PASS** | Abschnitt 2 |
| `make test-docker` (voll) | nicht gelaufen | zehn Minuten, bekannt rot durch einen anderen, nicht berührten Test (`..._NestedPodman_BuildEgress`); der berührte Docker-Test lief gezielt |
| `make verify-depguard` | nicht gelaufen | keine Importregel berührt (`make lint` in `make gates` grün; kein neuer Import in `application`) |
| `git diff --stat 1e445d3..HEAD -- spec/lastenheft.md docs/plan/adr .github .d-check.yml .harness harness/conventions.md` | **leer** | Lastenheft, ADRs, Workflows, `.d-check.yml`, Templates und Baseline unverändert |

## 2. Docker-Test selbst gefahren

Aufruf wie das Ziel `test-docker` in `repo.mk`: erst der Bau-Schritt des Ziels (`--target test-docker-tools`, Exit 0), dann im Image mit denselben Mounts (Host-Netz, Docker-Socket, `/tmp`) `go test -tags docker -timeout 30m -count=1 -v -run '^TestE2E_SandboxDevcontainer_EgressSubdomains$' ./internal/e2e/`.

Ergebnis: **`--- PASS: TestE2E_SandboxDevcontainer_EgressSubdomains (10.58s)`**, `ok … 10.588s`, Gesamtdauer des Aufrufs (ohne Bau) 17 s. Weder SKIP noch FAIL.

Nachlauf: Container- und Volume-Bestand vor und nach dem Lauf verglichen: 0 neue Container, 0 neue Volumes, 0 Images mit Präfix `uboot-e2e-`. Es war nichts zu räumen, ich habe nichts entfernt.

Text gegen Test (gelesen in `internal/e2e/sandbox_devcontainer_docker_test.go`, `egressSubdomainProbeScript`):

- „one and two levels below it“: gepinnt durch `A www.example.org` (Eintrag `example.org`, eine Ebene) und `A s3.us-east-1.amazonaws.com` (Eintrag `amazonaws.com`, zwei Ebenen); beide müssen einen Status ungleich `000` haben.
- „does not open its parent name (example.org) for a connection by name“: gepinnt durch `B example.org` (Eintrag `www.example.org`), der Status muss genau `000` sein; „by name“ ist die Grenze dieses Status (Namensauflösung oder verworfene Adresse, der Test trennt beides nicht). Die Verbindung über eine geteilte Adresse ist nicht geprüft und der Text sagt sie nicht zu.
- Drei und mehr Ebenen: nicht gepinnt, und weder Langhilfe, `Notes:`, Wildcard-Fehler noch Spezifikation sagen sie (Wildcard-Fehler nach `dccfa44`: „an allowed name already covers its subdomains, one and two levels below it“). Der Text sagt nicht mehr als der Test.
- Der Testkommentar (Zeilen 196 bis 200) nennt die drei gepinnten Fälle und ist mit dem Test deckungsgleich.

## 3. End-zu-End-Probe mit dem gebauten Binary

`make build` (Exit 0), dann das Image `u-boot:latest` mit gemountetem Wegwerf-Projekt im Scratchpad (`/tmp/claude-1000/-Development-u-boot/b5840ffc-5074-42ee-b764-156f5e7a42fb/scratchpad/verify-hinweis/p/`). `init` endete mit Exit 1, weil das Runtime-Image kein `git` enthält (`git rev-parse failed: exec: "git": executable file not found`); die Projektstruktur entstand, `u-boot.yaml` liegt vor. Das ist ein Probe-Artefakt (kein Git-Repo im distroless-Image) und kein Befund dieses Slice. Gesetzt: `devcontainer.profile sandbox`, `devcontainer.sandbox.egress.enabled true`, wechselnde `devcontainer.sandbox.egress.allow`.

| Schritt | Beobachtet (wörtlich, gekürzt) | Soll |
|---|---|---|
| `generate devcontainer` ohne Egress | „Generated devcontainer (…devcontainer.json, …Dockerfile).“ | kein Hinweis |
| erstes Aktivieren (`egress.enabled true`) | „Updated devcontainer managed block (…, .devcontainer/egress-init.sh).“; Warnungen nur zu `LH-FA-DEV-006` und `LH-FA-DEV-008` (Lockerung `--cap-add=NET_ADMIN`) | **kein** Neuaufbau-Hinweis: erfüllt | <!-- d-check:ignore (wörtliche Ausgabe des Binaries) -->

| `allow` auf `example.org, amazonaws.com`, dann `generate devcontainer` | „Updated devcontainer managed block (.devcontainer/egress-init.sh).“ und **„Warning (LH-FA-DEV-008): .devcontainer/egress-init.sh changed; an existing container keeps the egress list of the image it was built from, so rebuild the container for the new list to take effect (`devcontainer up --remove-existing-container`)“**, Exit 0 | Hinweis: erfüllt | <!-- d-check:ignore (wörtliche Ausgabe des Binaries) -->

| dritter Lauf unverändert | „devcontainer already up to date; no changes.“ | kein Hinweis (NoOp): erfüllt |
| `allow` geändert, `--quiet generate devcontainer` | nur die Lockerungs-Warnung (`--cap-add=NET_ADMIN`, gekennzeichnet); der Neuaufbau-Hinweis fehlt, Exit 0; Skript war geändert | Entfall unter `--quiet`: erfüllt |
| `allow` geändert, `--dry-run` | „Would update devcontainer managed block (.devcontainer/egress-init.sh).“; kein Neuaufbau-Hinweis | kein Hinweis: erfüllt |
| `allow` geändert, `--json` | Umschlag `"status":"warn"`, `"exitCode":0`, `diagnostics[]` mit drei Einträgen, der dritte `code: LH-FA-DEV-008`, `file: .devcontainer/egress-init.sh`, derselbe Wortlaut | Status `warn`, Exit 0, ein Eintrag: erfüllt | <!-- d-check:ignore (wörtliche Ausgabe des Binaries) -->

| `egress.enabled false`, `generate devcontainer` | „Updated devcontainer managed block (.devcontainer/devcontainer.json, .devcontainer/Dockerfile).“; kein Hinweis; `.devcontainer/egress-init.sh` bleibt stehen (`ls`) | kein Hinweis beim Abschalten, Skript bleibt: erfüllt |
| `egress list` | Notes: „A name in the list also allows its subdomains, one and two levels below it; an entry for a subdomain (www.example.org) does not open its parent name (example.org) for a connection by name.“; „… a container keeps the list of the image it was built from, so a changed list takes effect after a rebuild of the container (measured with the Dev Container CLI 0.89.0: …“ | Texte gekürzt, Neuaufbau mit dem Wort: erfüllt |

Keine Doppelmeldung: in jedem Lauf genau eine Zeile mit „egress-init.sh changed“. Kein neuer Exit-Code, kein Sentinel: der Code [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) ist als Warnungs-Code bereits in Gebrauch (Lockerung, Fallback), `application` liest für den Hinweis nur den Plan der Dateien (`egressRebuildWarnings`, reine Funktion; `internal/adapter/driving/cli/generate.go` und `internal/hexagon/port/driving/` nicht im Diff). Nicht probiert: `--diff` (Spezifikation, Handbuch und CHANGELOG sagen dazu nichts zu; im Review als INFO geführt).

## 4. DoD Punkt für Punkt

Die Kästchen in der Slice-Datei sind alle noch offen; das Setzen ist Closure-Arbeit des Planners.

| # | DoD-Punkt | Befund | Evidence |
|---|---|---|---|
| 1 | Hinweis zum Neuaufbau bei Änderung von `egress-init.sh`; Test pinnt Fall und Gegenfälle (Änderung ohne `egress-init.sh`, NoOp, `--dry-run`); Behandlung unter `--quiet` vor dem Code entschieden und getestet | **erfüllt** | `egressRebuildWarnings` (`internal/hexagon/application/generate.go`): Auslöser nur `p.relPath == egressInitRel && p.action == devcontainerActionReplaceBlock`, Guard `PreviewDryRun`. Tests in `generate_egress_rebuild_test.go`: `TestGenerateDevcontainer_EgressRebuildNotice` (erstes Schreiben, NoOp, Listenänderung mit genau einem Hinweis), `…_NotForOtherFiles`, `…_DryRunGivesNone`, `…_NotForActivationAndSwitchOff`; `quiet_test.go` `TestExecute_Generate_EgressRebuildNotice_QuietAndJSON` (`--quiet` streicht, `--json` mit und ohne `--quiet` trägt Status `warn`, Exit 0, genau einen Eintrag). Entscheidung (b) in [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) geführt, [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) („vier Ausnahmen“) unberührt. Probe Abschnitt 3 bestätigt jeden Fall |
| 2 | `Notes:` nennen den Neuaufbau mit dem Wort; `TestExecute_EgressList_Notes` pinnt den Wortlaut; [`SPEC-046`](../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung) sagt die Hälfte „eine geänderte Liste braucht einen Neuaufbau“ wieder zu | **erfüllt** | `egress_test.go`: Fragment „a container keeps the list of the image it was built from, so a changed list takes effect after a rebuild of the container“ ist Pflicht; [`SPEC-046`](../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung) (*Abgrenzung zum laufenden Container*) sagt den Neuaufbau-Satz, und der Satz hat damit den pinnenden Test (Auflage nicht nötig) |
| 3 | Subdomain-Texte: die vier Stellen sagen nur das Gemessene; Tests pinnen den neuen Wortlaut; Spezifikation, Handbuch, CHANGELOG sagen dasselbe | **erfüllt** | Langhilfe und `Notes:` (`egress.go`), Wildcard-Fehler (`sandbox.go`, `NewEgressHost`, zwei Stellen), Testkommentar (`sandbox_devcontainer_docker_test.go`); Pins: `TestExecute_EgressList_Notes`, `TestConfigSet_EgressAllowWildcard_Exit10WithHint` (vier Fälle), `TestNewEgressHost` (mit Gegenfall); Gemessenes durch den von mir gefahrenen Docker-Test belegt (Abschnitt 2); Handbuch (Abschnitt *Subdomains*, unverändert, gleiche Aussage) |
| 4 | `make gates` grün; Review durchgeführt, Report liegt vor; Doku-Update | **erfüllt** | Abschnitt 1; Review `c5be3c7`; Handbuch, CHANGELOG `[Unreleased]`, Spezifikation, `harness/replay.md` nachgezogen |
| 5 | Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | **offen: Planner** | §7 der Slice-Datei trägt „Wird bei der Closure gefüllt“; §6 trägt dreimal „weiter offen“; Abschnitt 8 |

## 5. Ein Vertrag, vier Stellen

| Aussage | Code | Spezifikation | Handbuch | CHANGELOG `[Unreleased]` |
|---|---|---|---|---|
| Hinweis nur bei ersetztem Block eines bestehenden `egress-init.sh` | `egressRebuildWarnings` | [`SPEC-046`](../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung) *Hinweis bei `u-boot generate devcontainer`* | Egress-Abschnitt (neuer Absatz) | *Added* |
| nicht bei erstem Schreiben, erster Aktivierung, Abschalten, ohne Änderung, `--dry-run` | Tests (Punkt 1) | [`SPEC-046`](../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung) | Handbuch | *Added* |
| `--quiet` entfällt, `--json` Status `warn`, Exit 0 | Kennzeichen-Test, CLI-Test | [`SPEC-046`](../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung), [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) (unter den entfallenden) | Handbuch | *Added* |
| Befehl prüft nicht, ob ein Container existiert | keine Inspektion im Pfad (gelesen) | [`SPEC-046`](../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung) | Handbuch | *Added* |
| „eine und zwei Ebenen“, Elternname „bei der Verbindung per Name“ | `egress.go`, `sandbox.go` | [`SPEC-046`](../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung) *Ausgabe* ([`SPEC-019`](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) *Namen und Subdomains* unverändert, gleiche Aussage) | *Subdomains* (unverändert) | *Changed* |

Widerspruchsfrei. Die Spezifikation trägt keinen Messstand (Zielbild), nur Zusagen mit Test; die neuen Sätze verweisen nur aufwärts oder auf Geschwister-`SPEC-*`. Das Handbuch trägt keine Kennungen und keine Verweise auf Planungsdokumente ([`AGENTS.md`](../../AGENTS.md) §5 Regel 6), die CHANGELOG-Einträge sind Englisch und tragen verlinkte Kennungen. Texte im Code, die CLI-Ausgabe und der CHANGELOG sind Englisch ([`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache)). Alte Wortlaute („already included“, „never opens its parent name“) kommen im Baum außerhalb von Plan, Register, Reviews und früheren CHANGELOG-Einträgen nicht mehr vor (Suche gelaufen). Rest-Klammer mit Messstand in den `Notes:` (Dev Container CLI 0.89.0, Colima, arm64, Editor-Aktionen nicht gemessen): sie gehört zur Aussage „Container behält die Liste, `--remove-existing-container` wendet die neue an“, ist durch `TestExecute_EgressList_Notes` im Wortlaut gepinnt, mit der Messung im Handbuch deckungsgleich und trägt keine Subdomain-Aussage; **vertretbar**. Das Gebot „Meldungen und Handbuch = Ist-Zustand“ verlangt hier keine Streichung.

## 6. Review-Findings einzeln

| ID | Klasse | Stand | Beleg |
|---|---|---|---|
| F-1 | LOW (Grenze beim ersten Aktivieren nicht benannt) | **behoben** | `12b3dca`: [`SPEC-046`](../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung), Handbuch und CHANGELOG nennen „nicht beim ersten Schreiben (auch nicht bei der ersten Aktivierung)“ und „nicht beim Abschalten“; `dccfa44` pinnt beide Fälle mit `…_NotForActivationAndSwitchOff`; meine Probe bestätigt beide (Abschnitt 3) |
| F-2 | LOW (Wildcard-Meldung ohne Tiefe) | **behoben** | `dccfa44`: „an allowed name already covers its subdomains, one and two levels below it“ an beiden Stellen in `NewEgressHost`, gepinnt in `TestNewEgressHost` und `TestConfigSet_EgressAllowWildcard_Exit10WithHint`; Tiefe durch den Docker-Test gedeckt |
| F-3 | INFO (Pin der Subdomain-Tiefe hinter dem Build-Tag) | **behoben durch diesen Report** | Test von mir gefahren: PASS, 10.58 s (Abschnitt 2) |
| F-4 | INFO (Hinweis folgt der Änderung des Skripts, nicht der Liste) | **begründet offen** | Spezifikation und Handbuch treffen den Auslöser („ersetzt den Block“); die Meldung sagt „for the new list“; ein Neuaufbau ist auch bei geändertem Skript ohne Listenänderung richtig; Risiko gering |
| F-5 | INFO (Pin in zwei Hälften) | **begründet offen** | Kennzeichen im Anwendungstest, Filter im CLI-Test mit Fake-Dienst; meine Probe fuhr den echten Dienst durch den CLI-Pfad (`--quiet`, `--json`) und bestätigt das Zusammenspiel; `--diff` weder zugesagt noch gepinnt |
| F-6 | INFO (Messstand-Klammer im Rest der Notes) | **begründet offen, vertretbar** | Abschnitt 5 |

## 7. Plan-Tabelle gegen `git diff --stat -M 1e445d3..HEAD` (nach dem letzten Commit, [`AGENTS.md`](../../AGENTS.md) §3.10)

Dateiliste Zeile für Zeile gegen §3 abgeglichen: 15 Dateien, jede hat eine Zeile oder fällt unter die Sammelzeile (Spezifikation, Handbuch, CHANGELOG, `harness/replay.md`, Slice-Datei, Review-Report). Plan-Zeilen ohne Diff: nur dieser Report und die Closure-Dateien (Register, Roadmap, `open/README.md`, nicht im Diff, weil Start und Closure außerhalb der Range liegen). Die als „unverändert“ gelesenen Dateien (`internal/adapter/driving/cli/generate.go`, `internal/hexagon/port/driving/removeservice.go`, `generate_acceptance_test.go`) fehlen im Diff. Abweichung (INFO): die Zeile zum E2E-Test nennt „Kommentar, Zeile 197“, der Kommentar beginnt jetzt in Zeile 196 (nach `AGENTS.md` §3.10 ein Zeilenverweis ohne Ablauf; Datei und Symbol stimmen).

## 8. Auflagen

An den Planner (Closure und Plan-Pflege, kein Eingriff durch den Verifier):

1. **Zusage ohne pinnenden Test (LOW, [`AGENTS.md`](../../AGENTS.md) §3.10, letzter Punkt):** „Ob ein Container existiert, prüft der Befehl nicht“ steht in [`SPEC-046`](../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung), Handbuch und CHANGELOG und ist wahr (gelesen: `egressRebuildWarnings` erhält nur Plan und Modus), aber kein Test pinnt sie: `generate_egress_rebuild_test.go` enthält keinen Inspector-Fall. Entweder ein Test (der Hinweis erscheint unabhängig vom Befund eines Inspectors, einschließlich „kein Container“) oder die Zusage streichen. Kleiner Nachtrag im Slice oder als eigener Fix, nicht Schließungs-Blocker.
2. **Plan-Tabelle:** „Zeile 197“ (§3, Zeile zum E2E-Test; §1 beschreibt den gefundenen Stand und darf bleiben) auf „Kommentar zu `TestE2E_SandboxDevcontainer_EgressSubdomains`“ setzen.
3. **Closure-Notiz §7** mit Steering-Loop-Lerneintrag; die Verification-Evidence weist aus: voller `make test-docker` und `make verify-depguard` nicht gelaufen mit Grund (Abschnitt 1), `make doc-immutable` 0 Befunde ohne Beleg-Wert (Lücke der Inline-Status-ADRs), `--diff` nicht probiert.
4. **Register fortschreiben:** [`BEO-CLI/subdomain-text-ohne-vorbehalt`](../plan/planning/observations/BEO-CLI/subdomain-text-ohne-vorbehalt/observation.md) (Ausgang: Texte auf das Gemessene gekürzt, kein Vorbehalt; Entscheidung (a)); [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md) hat im Review einen weiteren Beleg (Tiefe der Wildcard-Meldung, F-2): Auftreten fortschreiben; [`BEO-WZG/kopplung-nicht-nachgezogen`](../plan/planning/observations/BEO-WZG/kopplung-nicht-nachgezogen/observation.md): dieser Slice hat Code, Spezifikation, Handbuch, CHANGELOG zusammen nachgezogen (kein neuer Beleg, aber der Nachweis, dass die Zusammenführung hielt); die drei Paarungen tragen.
5. **Risiko-Ausgänge §6 (drei Risiken, alle „weiter offen“):** (i) Hinweis auch ohne existierenden Container: getragen, Wortlaut „an existing container“ und Spezifikationssatz „Ob ein Container existiert, prüft der Befehl nicht“; (ii) Antwort trägt den Grund der Änderung nicht: **entfallen**, der Hinweis entsteht in der Anwendungsschicht aus dem Plan, keine Port-Änderung, die Rückführung `in-progress` → `next` trat nicht ein; (iii) Vorbehalt im Wortlaut als Unsicherheit: entfallen durch Entscheidung (a), die Texte tragen keinen Vorbehalt.
6. **Setzen der DoD-Kästchen** und **Move nach `done/`** in der Reihenfolge von [`AGENTS.md`](../../AGENTS.md) §3.3 (erst Inhalt mit Commit-Hash, dann reiner `git mv`); Roadmap und `open/README.md` mitziehen. Der Slice gehört zum Release-Schnitt nach der Welle; CHANGELOG `[Unreleased]` trägt die Einträge.

An den Projektinhaber: keine neue Entscheidung. Die Lücke von `doc-immutable` bei Inline-Status-ADRs bleibt wie in der Beobachtung geführt.

## 9. Verification Evidence (Block nach `harness/verification.md`)

Scope:
- Slice: [`slice-v1-egress-hinweise-neuaufbau-und-messstand`](../plan/planning/in-progress/slice-v1-egress-hinweise-neuaufbau-und-messstand.md)
- IDs: [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`ADR-0025`](../plan/adr/0025-egress-liste-im-laufenden-container.md) (gelesen, unverändert)
- Artefakte: `internal/hexagon/application/generate.go`, `internal/adapter/driving/cli/egress.go`, `internal/hexagon/domain/sandbox.go`, zugehörige Tests, `spec/spezifikation.md` ([`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), [`SPEC-046`](../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung)), `docs/user/benutzerhandbuch.md`, `CHANGELOG.md`, `harness/replay.md`

DoD-Abgleich:
- [x] Punkte 1 bis 4 (Abschnitt 4)
- [ ] Punkt 5 (Closure) — offen: Planner, Auflagen 3 bis 6

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | Exit 0, Coverage 91.40 % |
| `make docs-check`, `make doc-planning`, `make doc-trace`, `make ablage-pruefung` | pass | 0 Befunde, 0 Waisen, 26 Fälle |
| `TestE2E_SandboxDevcontainer_EgressSubdomains` | pass | 10.58 s, kein Rückstand |
| End-zu-End-Probe mit dem gebauten Binary | pass | Abschnitt 3 |
| `make doc-immutable` | 0 Befunde, kein Beleg-Wert | Lücke Inline-Status-ADRs; ADR-Diff leer |
| `make test-docker` (voll), `make verify-depguard` | not run | Begründung Abschnitt 1 |

Traceability: [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) durch `generate_egress_rebuild_test.go`, `egress_test.go`, `sandbox_test.go` und den Docker-Test; [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) durch `TestExecute_Generate_EgressRebuildNotice_QuietAndJSON` und den Kennzeichen-Test; [`ADR-0025`](../plan/adr/0025-egress-liste-im-laufenden-container.md) (Neuaufbau statt Neu-Laden) unverändert, der Hinweis nennt den Neuaufbau und ändert ihn nicht.

Carveouts: neu none; gelöst none; unverändert none.

Commit / Artefakt: HEAD `e78d71f` vor diesem Report.
