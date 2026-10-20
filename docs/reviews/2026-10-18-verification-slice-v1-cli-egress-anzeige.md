# Verification-Report: slice-v1-cli-egress-anzeige — 2026-10-18 <!-- d-check:ignore (Slice-Kennung im Titel, Zitat) -->

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“, gegen Plan, DoD und Entscheidungen). Kein Review, keine Validierung, keine Reparatur.

**Gegenstand:** Slice `slice-v1-cli-egress-anzeige` (`docs/plan/planning/in-progress/`), Befehl `u-boot egress list`; Inhalts-Commit `94581b9`, Review-Report `95f9b03`, Fix-Runde `09496b0` (HEAD); Range `9953000..HEAD` (27 Dateien, 2048 Einfügungen, 75 Löschungen). Arbeitsbaum zu Beginn und am Ende sauber. Keine Änderung an Lastenheft, Spezifikation, Architektur-Sicht oder ADR im Diff. <!-- d-check:ignore (Kennungen als Zitat) -->

**Eingang:** Review `docs/reviews/2026-10-18-slice-v1-cli-egress-anzeige.md` (F-1 bis F-3 MEDIUM, F-4 bis F-7 LOW, F-8 bis F-14 INFO), Fix-Runde `09496b0`, Umsetzungsstand und Plan-Tabelle im Slice (§1, §3), Sensor-Behauptung und Mutationsliste des Implementers (die Mutationsliste lag dem Verifier als Auftragstext vor, nicht als Datei). <!-- d-check:ignore (Kennungen als Zitat) -->

**Gesamturteil: GO mit Auflagen.** Alle Sensoren, die der Verifier selbst gefahren hat, sind grün; jede Teilzusage des Lastenheft-Spiegelstrichs ist an Code und Test belegt; die zwei nachgefahrenen Mutationen werden von den Tests gefangen; die Plan-Tabelle stimmt mit dem Diff überein; die Findings F-1 bis F-14 sind behoben oder begründet. Die Entscheidungen 1 bis 5 des Implementers werden bestätigt (Abschnitt 5). Auflagen (Abschnitt 9): Closure-Notiz, DoD-Haken, Zustandsmarker im Slice, Folge-Slice für das Enum `command`, Register-Einträge, Spec-Nachzug. Keine DoD-Verletzung, kein HIGH, kein MEDIUM.

## 1. Sensoren (selbst gefahren)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | pass, **Exit 0** (Exit-Code des `make`-Prozesses selbst erfasst, Ausgabe in eine Datei im Scratchpad umgeleitet) | Lauf nach HEAD `09496b0`; im Log gelesen: `lint` „0 issues.“, `coverage`-Stage mit allen Paketen `ok`; Rest über den Exit-Code |
| `make doc-trace` | pass, **Exit 0** | Zeile `LH-FA-DEV-008`: ADRs `ADR-0012` bis `ADR-0025`, Slices u. a. `slice-v1-cli-egress-anzeige`, Status `ok`; Zeile `LH-FA-DEV-006` trägt `slice-v1-cli-egress-anzeige` ebenfalls | <!-- d-check:ignore (Kennungen als Zitat) -->
| `make verify-depguard` | pass, **Exit 0** | „all 9 depguard rules verified“, jede Regel `PASS` |
| `make docs-check` | pass, **Exit 0** | „d-check: 433 Datei(en) geprüft, 0 Befund(e)“ (vor Anlage dieses Reports; mit diesem Report erneut gefahren, siehe Abschnitt 10) |
| `make test` mit Mutation 1 und 2 | **rot, Exit 2** je Mutation (Abschnitt 4) | Logs im Scratchpad; Namen der roten Tests per Docker-Container aus dem Log gezogen |
| `make test-docker` | nicht gelaufen | kein Docker-Pfad berührt: der Diff enthält keine Datei unter `internal/adapter/driven/`, keine unter `internal/e2e/`, kein `Dockerfile`, kein `Makefile`; `cmd/uboot/main.go` bekommt nur eine Zeile, die den neuen Service mit den bestehenden `fsAdapter`, `yamlAdapter`, `logAdapter` baut; `EgressListService` hält nur `FileSystem`, `YAMLCodec`, `Logger` (kein Docker-Port) |
| `make doc-immutable` | nicht gelaufen | keine ADR im Diff |
| `TestE2E_SandboxDevcontainer_EgressSubdomains` (Beleg des Subdomain-Satzes) | nicht gelaufen | braucht `NET_ADMIN` und Netz (Docker-Pfad, Vorgänger-Slice); der Slice ändert ihn nicht |

## 2. Lastenheft 0.14.0: jede Teilzusage am Code und Test

Quelle: `spec/lastenheft.md` Zeile 472 (Spiegelstrich in `LH-FA-DEV-008`), Zeile 196 (`LH-FA-CLI-005`), Zeile 1121 (`LH-SA-CLI-002`). <!-- d-check:ignore (Kennungen als Zitat) -->

| Teilzusage | Beleg im Code | Beleg im Test | Ergebnis |
|---|---|---|---|
| Wirksame Allowlist und Herkunft je Host (Basis, Container-Runtime, Feature, Clone-Quelle, Nutzerliste) | `egressHostSources` (`internal/hexagon/application/devcontainer_egress.go`), Namen `driving.EgressSource*` (`internal/hexagon/port/driving/egress.go`) | `TestEgressList_OriginPerHost` (je Host alle Quellen, feste Reihenfolge, `ghcr.io` mit drei, `registry.npmjs.org` und `github.com` mit mehreren), `TestEgressList_FeatureOriginsAreSortedByName` (100 Läufe) | erfüllt; sechster Name `docker-apt` siehe V-3 |
| Liest nur `u-boot.yaml` und die Git-Konfiguration | `EgressListService` hält nur `FileSystem`, `YAMLCodec`, `Logger`; `readConfig`, `readGitOriginURL` | `TestEgressList_OnlyProbesAndReadsTheTwoProjectFiles` (Aufzeichner `recordingFS` über **alle** `FileSystem`-Methoden; erlaubt nur `Exists` und `ReadFile` auf die zwei Pfade, beide Aufrufe müssen vorkommen) | erfüllt (schließt F-9 des Reviews) |
| Kein Netz, kein Docker, kein Container | kein Docker-/Netz-Port im Konstruktor; depguard `application-no-net` PASS | `TestRun_EgressList_MatchesTheGeneratedStartScript` läuft mit den echten Adaptern ohne Docker | erfüllt |
| Ändert nichts | kein Schreibaufruf im Use Case | `snapshotTree` vor und nach dem Lauf (`cmd/uboot/main_test.go`, beide Nutzerpfad-Tests, auch der Fehlerfall); `recordingFS` | erfüllt |
| Dieselbe Liste wie `generate devcontainer` | `egressHosts` (Startscript, aus `devcontainerTemplateData`) und `EgressListService.List` rufen beide `egressHostSources` | `TestEgressList_EqualsTheStartScriptList` (elf Konfigurationen gegen die `hosts=`-Zeile der Generator-Ausgabe), `TestRun_EgressList_MatchesTheGeneratedStartScript` (Nutzerpfad), `TestGenerateDevcontainer_SandboxEgress_HostsLineGolden` (Literal) | erfüllt |
| Aktivierungsstatus (Sandbox-Profil und `egress.enabled`); nicht aktiv: Liste, die bei Aktivierung gälte, und sagt das | `EgressListResponse.Active`, `egressInactiveReason`, `renderEgressListText` | `TestEgressList_ActivationStateAndSameList` (fünf Zustände, gleiche Liste), `TestExecute_EgressList_NotActiveStatesTheReason`, `TestExecute_EgressList_JSON` (Zustand ist Datum, kein Diagnostic) | erfüllt |
| Subdomains-Satz | Zeile in der Ausgabe (`renderEgressListText`), Hilfetext | `TestExecute_EgressList_Notes` pinnt den Wortlaut; das **Verhalten** trägt `ADR-0012` Entscheidung 2 und `TestE2E_SandboxDevcontainer_EgressSubdomains` (nicht gefahren) | Wortlaut belegt, Verhalten durch Vorgänger-Beleg (V-5) | <!-- d-check:ignore (Kennungen als Zitat) -->
| `LH-FA-CLI-005`: Ausgabe des Abfragebefehls bleibt unter `--quiet` | `runEgressList` schreibt dieselben Bytes | `TestExecute_EgressList_QuietKeepsTheCompleteOutput` (Text und `--json`, aktiv und nicht aktiv), `TestExecute_Quiet_QueryCommandsKeepTheirOutput` (Fall `egress list`) | erfüllt | <!-- d-check:ignore (Kennungen als Zitat) -->
| `LH-FA-CLI-006`: Exit-Codes | `isFilesystemError` (14), `isUsageError` (2), zentrale Fachfehler-Klassifikation (10) | `TestExecute_EgressList_Errors` (10, 10, 10, 14, 1 mit Diagnose-Code je Klasse), `TestExecute_EgressList_RejectsPreviewFlagsAndArguments` (2), `TestExecute_Egress_Bare` (2), `TestRun_EgressList_UnsafeOriginExits10WithoutActiveRestriction` (10 mit den echten Adaptern) | erfüllt | <!-- d-check:ignore (Kennungen als Zitat) -->
| `LH-SA-CLI-002`: Befehlsübersicht (Zeile 1121, abgeleitet, nicht freigegeben) | — | — | im Lastenheft; Freigabe des Projektinhabers weiter ausstehend (Auflage 6) | <!-- d-check:ignore (Kennungen als Zitat) -->

## 3. Plan-vs-Diff (Plan-Tabelle §3 gegen `git diff 9953000..HEAD --stat`, Stand nach der Fix-Runde)

27 Dateien im Diff. Jede Datei hat eine Plan-Zeile, jede Plan-Zeile mit Datei-Eintrag einen Diff:

| Diff-Datei(en) | Plan-Zeile | Befund |
|---|---|---|
| `devcontainer_egress.go`, `egress_list.go`, `port/driving/egress.go` | 1, 3, 4 | ok |
| `devcontainer_egress_test.go`, `egress_list_test.go`, `generate_sandbox_test.go` | 10, 11 | ok |
| `cli/egress.go`, `cli/egress_test.go`, `cli/quiet_modifying_test.go` | 6, 12 | ok (Plan nennt die Fundstelle `quiet_modifying_test.go` richtig) |
| `cli/root.go`, `cli/cli.go`, `cli/template.go`, `cli/rootjson_test.go` | 7, 8, 9 | ok (zwei Zeilen als „Fix-Runde, nicht im Plan“ ausgewiesen) |
| `cmd/uboot/main.go`, `cmd/uboot/main_test.go` | 10 | ok |
| `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md`, `harness/replay.md` | 13 | ok |
| `docs/maintainer/cli-json-contract.md` | 14 | ok |
| `done/slice-v2-spec-sandbox-egress-ollama.md` | 15 | ok | <!-- d-check:ignore (Kennungen als Zitat) -->
| `in-progress/roadmap.md` | 16 | ok |
| `docs/reviews/2026-10-18-slice-v1-cli-egress-anzeige.md` | 17 | ok | <!-- d-check:ignore (Kennungen als Zitat) -->
| Slice-Datei | 18 | ok |
| `open/README.md`, `done/slice-v1-devcontainer-start-ohne-cli.md` | — | **gehören nicht zu diesem Slice** (stammen aus `6743254`); der Plan führt sie zu Recht nicht, der Review weist sie aus; die Range `9953000..HEAD` enthält sie dennoch (V-6) | <!-- d-check:ignore (Kennungen als Zitat) -->
| `jsontestutil/coderegistry.go`, `devcontainer_sandbox.go` | 2, 12 („unverändert, bedingt entfallen“) | bestätigt: nicht im Diff; die gelesenen Fundstellen stimmen (`devcontainerTemplateData` ruft weiter `egressHosts`) |

Mit diesem Report wächst die Range auf 28 Dateien; die Closure ändert den Slice, die Roadmap und das Register. Der Abgleich ist nach der Closure ein drittes Mal zu wiederholen (`AGENTS.md` §3.10, V-1).

## 4. Mutationsstichprobe (selbst gefahren, einzige Produktcode-Änderung, danach per `git checkout` zurückgenommen)

| # | Mutation | Ergebnis | Rote Tests |
|---|---|---|---|
| 1 | `egress_list.go`, `List`: `profileIsSandbox(cfg.Devcontainer)` statt `true` im Aufruf von `resolveSandboxClone` | `make test` **Exit 2** | `TestRun_EgressList_UnsafeOriginExits10WithoutActiveRestriction`, `TestEgressList_UnsafeCloneUrlIsRefusedWithoutActiveRestriction/origin,_no_sandbox_profile`, `TestEgressList_CloneHostIsResolvedWithoutActiveRestriction` (drei Unterfälle) |
| 2 | `devcontainer_egress.go`, `enabledFeatureNames`: `sort.Strings(names)` entfernt | `make test` **Exit 2** | `TestEgressList_FeatureOriginsAreSortedByName` und `TestRun_EgressList_MatchesTheGeneratedStartScript` |

Beide Mutationen wurden gefangen; der Befund F-2 (Clone-Quelle ohne aktives Profil ohne Test) und die „nur wahrscheinlichkeitsweise“-Lücke der Sortierung (F-10) sind damit geschlossen. Arbeitsbaum danach: `git status` leer. Nicht nachgefahren sind die übrigen Mutationen des Implementer-Berichts (Duplikat-Filter, Reihenfolge, `Active()`, Exit-14-Eintrag, Verdrahtung in `main.go`); sie sind aus den Tests gelesen und passen zur Struktur der Tests, nicht gefahren.

## 5. Entscheidungen 1 bis 5 des Implementers (ohne Architect-Artefakt): bestätigt

Prüfmaßstab: `spec/architecture.md` §2.2 bis §2.4 und §7, `ADR-0012`, `ADR-0013`, Hexagon-Regeln (`make verify-depguard`, 9 Regeln PASS). Die Bestätigung ist die eines Verifiers (Konformität gegen Text und Gates), kein Architect-Artefakt; sie ersetzt nicht den Rollenwechsel, schließt aber die Lücke, die der Slice und der Review (F-4) benennen. <!-- d-check:ignore (Kennungen als Zitat) -->

| # | Entscheidung | Urteil | Begründung |
|---|---|---|---|
| 1 | Top-Level `egress list`; Verdrahtung über Option `cli.WithEgressListUseCase` | **bestätigt** | Der Name `u-boot egress list` steht im Lastenheft 0.14.0 (Zeile 472, 1121), er ist keine freie Wahl. Die Option folgt dem Muster der bestehenden Optionen (`WithGetwd`, `WithLogLevel`); `architecture.md` §2.3 verlangt für den Port Request/Response mit Domänen- und Primitivwerten, nicht die Form der Verdrahtung. Ohne Verdrahtung: Fehler, kein Panic (`TestExecute_EgressList_NotWired`). |
| 2 | `egressHostSources` als einzige Berechnung; `egressHosts` leitet ab | **bestätigt** | Beide Aufrufer teilen eine Funktion; `egressHosts` bildet nur die Host-Namen ab; die Gleichheit trägt `TestEgressList_EqualsTheStartScriptList`. Byte-Gleichheit des Generators: Golden-Literal und bestehende Egress-Tests grün. Typ `driving.EgressHostEntry` im Driving-Port ist nach §2.3 zulässig (Strings), `cli` darf `application` nicht importieren (`adapter-no-application` PASS). |
| 3 | Herkunftsnamen `base`, `runtime`, `docker-apt`, `feature:<name>`, `clone`, `user` | **bestätigt mit Hinweis** | Eine feste, kleine Menge, je Host jede Quelle einmal, feste Reihenfolge, Konstanten im Port (Teil des `--json`-Vertrags). Hinweis V-3: das Lastenheft nennt fünf Klassen; `docker-apt` ist eine Verfeinerung der Klasse Container-Runtime und muss in der Spezifikation als solche stehen. |
| 4 | Ohne aktive Restriktion: Liste, die bei Aktivierung gälte, mit Zustandszeile; `--json` führt `active`, `sandboxProfile`, `egressEnabled`, `cloneSource` als Daten | **bestätigt** | Wortlaut des Lastenhefts („zeigt er die Liste, die bei Aktivierung gälte, und nennt das“); Zustand ist Datum, keine Diagnose (`TestExecute_EgressList_JSON`). |
| 5 | Fachfehler (Exit 10, Diagnose `LH-FA-DEV-006`) bei unzulässiger `origin`-URL auch ohne aktive Restriktion | **bestätigt (Variante „ohne Clone-Host“ bleibt offene Entscheidung des Projektinhabers)** | Gleichlauf mit `generate devcontainer --sandbox` (`resolveSandboxClone`), gepinnt durch `TestEgressList_UnsafeCloneUrlIsRefusedWithoutActiveRestriction`, `TestRun_EgressList_UnsafeOriginExits10WithoutActiveRestriction` und Mutation 1; Dual-Classifier-Regel (`architecture.md` §7.2) eingehalten (zentral: `isFilesystemError`, `isUsageError`, Fachfehler; je Subkommando `mapEgressErrorToDiagnostic`). Das Lastenheft schweigt zum Fehlerfall; die Spezifikation muss ihn tragen (Auflage 3). | <!-- d-check:ignore (Kennungen als Zitat) -->

ADR-0013 (Referenzmodell): im Diff keine abwärts gerichtete Kennung in Lastenheft, Spezifikation oder Architektur-Sicht (diese sind nicht im Diff); `docs/user/`, `README.md`, `README.de.md` tragen keine Kennungen (`AGENTS.md` §5 Zeile 6; die Diffs der drei Dateien sind gelesen: keine `LH-`, `ADR-`, `SPEC-`, Slice-Kennung). `ADR-0012`: die Basis-Hosts im Code stimmen mit Entscheidung 4 überein; Subdomains der Namen: Entscheidung 2. <!-- d-check:ignore (Kennungen als Zitat) -->

## 6. Handbuch, README, CHANGELOG, Messaussage gegen Code und Test (`AGENTS.md` §3.10)

| Aussage | Beleg | Ergebnis |
|---|---|---|
| Befehlsübersicht: Zeile `u-boot egress list`, Verweis auf Abschnitt 4.11 | Überschrift `### 4.11 Einen Sandbox-Devcontainer …` an Zeile 921 | stimmt |
| „Ohne Docker funktionieren … `egress list`“ (Handbuch zweimal) | Konstruktor ohne Docker-Port; Nutzerpfad-Test | belegt |
| Herkunftsnamen, Mehrfach-Herkunft („`ghcr.io` bei `podman` oder `docker` zusammen mit `homebrew`“) | `TestEgressList_OriginPerHost` | belegt |
| „Ein Test vergleicht seine Liste mit der Zeile `hosts=` des erzeugten `egress-init.sh`“ | `TestEgressList_EqualsTheStartScriptList`, `TestRun_EgressList_MatchesTheGeneratedStartScript` | belegt |
| Bedingung „auch ohne Sandbox-Profil und ohne aktive Beschränkung weist `egress list` eine unzulässige `origin`-URL ab (Exit 10)“ | Tests und Mutation 1 | belegt (F-6 behoben) |
| „`devcontainer.sandbox.repository` mit Zugangsdaten weist schon die Prüfung der Konfiguration ab (Exit 10)“ | Fall `repository, no sandbox profile` in `TestEgressList_UnsafeCloneUrlIsRefusedWithoutActiveRestriction` | belegt |
| `--quiet` ändert die Ausgabe nicht; `--dry-run`/`--diff` Exit 2; bare `egress --json` Exit 2 ohne JSON | `QuietKeepsTheCompleteOutput`, `RejectsPreviewFlagsAndArguments`, `TestExecute_Egress_Bare` | belegt |
| JSON-Felder (`active`, `sandboxProfile`, `egressEnabled`, `cloneSource`, `hosts`) | `TestExecute_EgressList_JSON` | belegt |
| Beispiel der Ausgabe als „Auszug (drei von mehreren Hosts, ohne `Notes:`)“ gekennzeichnet | Handbuch-Text; Format entspricht `renderEgressListText` (Kopf, Leerzeile, Tabelle) | stimmt (F-6 behoben) |
| **Messaussage** „Neuaufbau statt Neustart“ (Handbuch, Hilfe-Ausgabe, CHANGELOG) | Quelle: Messung M4 im Vorgänger-Slice (Messstand: Docker 29.5.2, Colima 0.10.3, arm64, Dev-Container-CLI 0.89.0, **ein Lauf**; nicht gemessen: Docker Desktop, Linux ohne VM, `amd64`, VS Code). Inhalt stimmt: Stoppen und Starten ließ die alte Liste, `devcontainer up --remove-existing-container` setzte die neue durch. Grenzen in der Aussage: Handbuch nennt CLI-Version, Docker, Colima, arm64, „nicht gemessen: Docker Desktop, Linux ohne VM, Editor-Aktion“; die Ausgabe nennt CLI-Version, Colima, arm64, „editor actions not measured“; das CHANGELOG nennt CLI-Version, Colima, arm64, Editor-Aktionen | als Messaussage gekennzeichnet und inhaltlich korrekt; Grenzen unvollständig, siehe V-2 | <!-- d-check:ignore (Kennungen als Zitat) -->
| `CHANGELOG.md` `[Unreleased]` `Added`: Quellen, Reihenfolge, `--quiet`, `--json`-Felder, Exit 2 und 10, Clone-Auflösung | Tests wie oben | belegt |
| `harness/replay.md`: Golden-Zeile nennt Fälle und Tests | genannte Tests existieren (`TestEgressList_*`, `TestGenerateDevcontainer_SandboxEgress_HostsLineGolden`, `TestRun_EgressList_MatchesTheGeneratedStartScript`, `egress_test.go`) | belegt |

## 7. Review-Findings F-1 bis F-14

Der Review-Report trägt keine Spalte „Behandlung“; die Zuordnung steht in der Commit-Message `09496b0` und im Slice (§1, §3, §6). Verifier-Prüfung am Baum:

| Finding | Stand | Verifier-Beleg |
|---|---|---|
| F-1 MEDIUM Alleinstellung „einzige Reject-Form“ | **behoben** | Diff gelesen: `root.go` (zwei Stellen), `cli.go` (Feld `json`), `rootjson_test.go` (Doc-Kommentar und Tabelle `bareRejects`, Baumlauf prüft beide), `template.go`, `cli-json-contract.md` §6.1 nennt zwei Formen; `git grep` nach „only `--json` reject“, „einzige verbliebene“, „only RunE-borne“, „the only gate-less“ in `internal`, `docs/maintainer`, `docs/user`: kein Treffer |
| F-2 MEDIUM Clone-Quelle im Zustand „nicht aktiv“ ohne Test | **behoben, durch Mutation 1 belegt** | `TestEgressList_CloneHostIsResolvedWithoutActiveRestriction`, `TestEgressList_UnsafeCloneUrlIsRefusedWithoutActiveRestriction`, `TestRun_EgressList_UnsafeOriginExits10WithoutActiveRestriction`: rot unter Mutation 1 |
| F-3 MEDIUM Vermerk zum Enum `command` | **behoben bis auf den Folge-Slice** | `cli-json-contract.md` §6.10 nennt `LH-NFA-USE-004`, `LH-FA-CLI-007`, `SPEC-001`; die Spezifikationsseite ist als Übergabe im Slice `slice-v2-spec-sandbox-egress-ollama` eingetragen (Diff gelesen); der Lastenheft-Anteil bleibt offen (Auflage 2) | <!-- d-check:ignore (Kennungen als Zitat) -->
| F-4 LOW Entscheidungen ohne Architect | **ausgewiesen und bestätigt** | Slice §1 „Abweichung von der Rollenfolge“; Bestätigung in Abschnitt 5 |
| F-5 LOW Messaussage ohne Messgrenze | **behoben, mit Rest (V-2)** | Ausgabe nennt Colima, arm64, „editor actions not measured“; pinnt `TestExecute_EgressList_Notes` |
| F-6 LOW Handbuch-Bedingung, Auszug | **behoben** | Abschnitt 6 |
| F-7 LOW Kommentar `WithEgressListUseCase` | **behoben** | Diff `cli.go`: der Kommentar nennt den Zustand („Without it `egress list` ends with a "not wired" error“); kein Vergleich mit einer verworfenen Form |
| F-8 INFO Entscheidung 5 vertretbar | **bestätigt** | Abschnitt 5 |
| F-9 INFO Test beobachtet nicht alle Port-Methoden | **behoben** | `recordingFS` überschreibt alle zwölf Methoden; der Test erlaubt genau vier Aufrufe |
| F-10 INFO Mutationsfestigkeit | **geschlossen für die zwei offenen Zweige** | Abschnitt 4 |
| F-11 INFO `cloneSourceHost` bei lokalen Pfaden | **übergeben, nicht gefahren** | Slice §6 (Beobachtung) und Übergabe an `slice-v2-spec-sandbox-egress-ollama`; Register-Eintrag steht aus (Auflage 4); der Verifier hat das Verhalten gelesen (Code von `cloneSourceHost`: ohne `://` und ohne `@` Schnitt an `:`), **nicht gefahren** | <!-- d-check:ignore (Kennungen als Zitat) -->
| F-12 INFO weitere Aufzählungen (`SPEC-043`, `SPEC-001`, `architecture.md`) | **übergeben** | Slice `slice-v2-spec-sandbox-egress-ollama` trägt sie; Spec und Architektur im Diff nicht angefasst, wie es die Source Precedence verlangt | <!-- d-check:ignore (Kennungen als Zitat) -->
| F-13 INFO Plan-Tabelle nach der Fix-Runde | **erfüllt für den Stand `09496b0`** | Abschnitt 3; nach der Closure zu wiederholen |
| F-14 INFO Regelverstöße des Implementers | **ausgewiesen, Register offen** | Abschnitt 10 und Auflage 4 |

## 8. Findings des Verifiers

| ID | Klasse | Befund |
|---|---|---|
| V-1 | LOW | Zustandsmarker im Slice: §1 führt die Entscheidungen 1 bis 5 als „Vorschlag, den der Verifier bestätigt oder verwirft“, die DoD-Haken (§2) und §7 sind offen, die Plan-Tabelle trägt „Stand: nach der Fix-Runde“. Mit diesem Report sind die Marker abgelaufen: Die Entscheidungen sind bestätigt (Abschnitt 5), die Tabelle ist nach der Closure erneut gegen `git diff --stat` abzugleichen (`AGENTS.md` §3.10). |
| V-2 | LOW | Messgrenzen unvollständig: Die Messung M4 ist **ein Lauf**; nicht gemessen sind Docker Desktop, Linux ohne VM, `amd64`, VS Code. Die Laufzeit-Ausgabe nennt davon nur „Colima, arm64, editor actions not measured“ und trägt die CLI-Version 0.89.0 als Literal ohne Bindung an die Messung; das Handbuch nennt Docker Desktop und Linux ohne VM, nicht „ein Lauf“ und nicht `amd64`. Die Aussage ist korrekt übernommen und als Messung gekennzeichnet; kein Blocker. |
| V-3 | LOW | Sechster Herkunftsname `docker-apt`: Das Lastenheft führt fünf Klassen (Basis, Container-Runtime, Feature, Clone-Quelle, Nutzerliste). `docker-apt` (Host `download.docker.com`, nur bei `nestedRuntime: docker`) ist eine Verfeinerung der Klasse Container-Runtime. Die Namen sind Teil des `--json`-Vertrags; die Spezifikation muss die sechs Namen und die Zuordnung nennen (`SPEC-019` und der neue Eintrag der Egress-Anzeige). | <!-- d-check:ignore (Kennungen als Zitat) -->
| V-4 | INFO | Der Spiegelstrich verlangt „liest nur `u-boot.yaml` und die Git-Konfiguration“: Der Beleg ist doppelt (Konstruktor des Service ohne weitere Ports, `recordingFS`-Test über alle `FileSystem`-Methoden). Ein späteres zusätzliches Lesen anderer Dateien bricht den Test. |
| V-5 | INFO | Der Subdomains-Satz ist in Ausgabe, Hilfe, Handbuch und CHANGELOG nur im Wortlaut gepinnt (`TestExecute_EgressList_Notes`); das Verhalten stützt sich auf `ADR-0012` Entscheidung 2 und den Docker-Test des Vorgänger-Slices, der hier nicht lief (Abschnitt 1). Kein Befund gegen den Code, eine nicht erneuerte Evidenz. | <!-- d-check:ignore (Kennungen als Zitat) -->
| V-6 | INFO | Die Range `9953000..HEAD` enthält zwei Dateien aus `6743254`, die nicht zu diesem Slice gehören; die Plan-Tabelle sollte sie als „nicht Teil dieses Slice“ nennen, damit der Abgleich gegen die Range ohne Rückfrage aufgeht. | <!-- d-check:ignore (Kennungen als Zitat) -->
| V-7 | INFO | Review-Fix-Runde `09496b0` wurde nicht erneut reviewt; der Verifier hat den Diff der Fix-Runde (Tests, Kommentare, Handbuch, CHANGELOG, Vertrags-Doku) gelesen und die zwei Mutationen gefahren, ohne Befund. Wer die Rollentrennung eng liest, kann einen Review der Fix-Runde verlangen. |

Kein HIGH, kein MEDIUM, keine DoD-Verletzung, kein Gate gelockert (`.golangci.yml`, `Makefile`, `Dockerfile`, `.d-check.yml` nicht im Diff), kein `//nolint`, keine Accepted-ADR im Diff.

## 9. DoD-Abgleich (Slice §2) und Auflagen für die Closure

| DoD-Zeile | Stand |
|---|---|
| Befehl und Berechnung (Tests belegen die Fälle) | **erfüllt**, Haken noch nicht gesetzt: die geforderten Fälle (Standard, `podman`, `docker`, Features mit zwei Quellen, Clone-Quelle, Nutzerliste, Aktivierungszustände, ohne Remote, `--quiet`) stehen in Abschnitt 2 |
| Vertrag nachgezogen (Handbuch, README, CHANGELOG, Golden-Zeile, `make doc-trace`) | **erfüllt** (Abschnitte 1 und 6) |
| Hilfe-Ausgabe (Zweck, Beispiel, Test pinnt Kurzbeschreibung) | **erfüllt** (`TestExecute_Egress_Help`) |
| `make gates` grün | **erfüllt** (Abschnitt 1), Haken noch nicht gesetzt |
| Review, Report liegt vor | **erfüllt** (Report vom 2026-10-18, Fix-Runde `09496b0`) |
| Doku-Update, Closure-Notiz, Steering-Loop-Eintrag, Register, Risiken §6, drei Paarungen | **offen, zu Recht**: §7 trägt den Platzhalter; Risiken in §6 tragen „offen“ oder „weiter offen“ |

**Auflagen:**

1. **Closure-Notiz (§7) füllen:** Befehlsort, Typ-Ort, Herkunftsnamen mit Fundstellen, Verweis auf diesen Report als Bestätigung der Entscheidungen 1 bis 5 (Zustandsmarker in §1 auf den Ist-Stand setzen, V-1), Steering-Loop-Eintrag, Folge-Slices, Ausgang je Risiko aus §6 (erwartet: Risiken 1 und 2 entfallen, 3 entschieden und belegt, 4 und 5 weiter offen mit Übergabe, 6 weiter offen als Zusage); DoD-Haken setzen; Plan-Tabelle nach der letzten Änderung gegen `git diff --stat` abgleichen und die zwei Fremddateien aus `6743254` nennen (V-6). <!-- d-check:ignore (Kennungen als Zitat) -->
2. **Folge-Slice für das Enum `command`:** Der Planer legt einen CR-Slice an, der `egress` in das Enum von `LH-FA-CLI-007` und `SPEC-001` aufnimmt (`LH-NFA-USE-004` bindet es für alle `--json`-Ausgaben); die Lastenheft-Änderung braucht die Freigabe des Projektinhabers; bis dahin liegt `command: "egress"` außerhalb des Vertrags, und kein Gate prüft das Enum. <!-- d-check:ignore (Kennungen als Zitat) -->
3. **Spezifikations-Slice** `slice-v2-spec-sandbox-egress-ollama`: nimmt die Übergaben (Fehlerfall der Clone-Quelle, Enum `command` nach dem Lastenheft, `SPEC-043`, `architecture.md`, sechs Herkunftsnamen samt `docker-apt` als Verfeinerung der Container-Runtime, Beobachtung `cloneSourceHost`) auf; die Variante „Liste ohne Clone-Host statt Fehler“ bleibt Entscheidung des Projektinhabers. <!-- d-check:ignore (Kennungen als Zitat) -->
4. **Beobachtungs-Register (Planer):** (a) `cloneSourceHost`: eine lokale Pfad-URL (`/srv/git/repo.git`) ergibt den „Host“ `/srv/git/repo.git` in der `hosts=`-Zeile und in `egress list` (vorbestehend, gelesen, nicht gefahren); im Register unter `docs/plan/planning/observations/` ist dazu **kein Eintrag** vorhanden (Verifier: Dateiliste und Inhaltssuche ohne Treffer); (b) `BEO-WZG/host-werkzeug-ausserhalb-docker`: dieser Slice ist der **fünfte Vorgang** (Stand des Registers: vier Belege, Schwelle erreicht, Ausgang nicht entschieden) und der **erste mit schreibendem Aufruf** (zwei Dateien per Shell-Heredoc `cat >>` am Host, laut Review und Auftrag; Inhalt im Diff gelesen und gate-grün); Entscheidung des Projektinhabers über den Ausgang steht aus. Der Verifier-Lauf selbst ist ein weiterer Beleg (Abschnitt 10, nur lesend).
5. **Entscheidung des Projektinhabers** bleibt offen: Variante „Liste ohne Clone-Host statt Fehler“ bei unzulässiger `origin`-URL.
6. **Freigabe `LH-SA-CLI-002`** (Zeile der Befehlsübersicht mit Priorität V2, abgeleitet) und die Klausel „Bestätigung ausstehend“ in der Historie-Zeile 0.14.0 des Lastenhefts: vor dem Release vom Projektinhaber zu bestätigen (Auflage aus dem Verification-Report des CR-Slice, hier unverändert offen). <!-- d-check:ignore (Kennungen als Zitat) -->
7. Roadmap auf „Nichts in Arbeit.“ zurücksetzen und `git mv` nach `done/` als reiner Move nach der Closure-Notiz (`AGENTS.md` §3.3), über das Werkzeug `make slice-mv`. <!-- d-check:ignore (Make-Target, keine Slice-Kennung) -->
8. Optional: Review der Fix-Runde `09496b0` (V-7).

## 10. Regelverstöße dieses Verifier-Laufs (ehrlich ausgewiesen)

Regel des Projektinhabers: „alles mit make und docker machen — keine tools ausserhalb des dockerfiles“. Gelesen wurde überwiegend mit dem Read-Tool, geprüft mit `make gates`, `make doc-trace`, `make verify-depguard`, `make docs-check`, `make test`; Logs wurden in eine Datei im Scratchpad umgeleitet; rote Testnamen zog ein `docker run … golang:1.27.2` aus dem Log. Die Mutationen schrieb das Write-Tool (das Edit-Werkzeug stand nicht zur Verfügung, daher wurde je Datei vollständig geschrieben und mit `git checkout --` zurückgenommen).

Verstöße, alle **nur lesend** und ohne Wirkung auf den Baum: host-seitig `cat`, `head`, `wc`, `cut` und `grep` als Pipeline-Glieder hinter `git` (Lesen des Exit-Status, Kürzen und Zählen von `git grep`-/`git show`-Ausgaben), ein `grep -o` hinter `git show`, `seq` und eine `read -t`-Warteschleife in der Shell, um auf das Ende des Hintergrundlaufs zu warten (kein Dateizugriff). Kein host-seitiges `go`, `sed`, `awk`, `cp`, `perl`, `python`, `jq`, `curl`, `gh`. Kein schreibender Aufruf am Host außerhalb des Write-Tools. Der Lauf ist damit ein weiterer lesender Beleg für `BEO-WZG/host-werkzeug-ausserhalb-docker`.

`make docs-check` fand beim ersten Lauf über diesen Report 17 Befunde (Kennungen ohne Link in Zitat-Zeilen); der Report trägt jetzt je betroffene Zeile den Ausnahme-Vermerk, wie der Report des Vorgänger-Slice; Ergebnis des zweiten Laufs: siehe Übergabe.

## 11. Nicht verifiziert

Die Freigabe des Wortlauts und der Zeile `LH-SA-CLI-002` durch den Projektinhaber (Gespräch, nicht prüfbar); `make test-docker` und `TestE2E_SandboxDevcontainer_EgressSubdomains` (kein Docker-Pfad berührt; Evidenz des Subdomain-Satzes aus dem Vorgänger); das Verhalten von `cloneSourceHost` bei lokalen Pfaden (gelesen, nicht gefahren); die übrigen Mutationen des Implementer-Berichts (gelesen, nicht gefahren); Laufzeitverhalten mit echter Dev-Container-CLI, `devcontainer up` und Editor (Messaussage stammt aus dem Vorgänger-Slice, Grenzen in V-2); `make doc-immutable` (keine ADR im Diff); die Wirkung der Hinweise auf einen Leser (Validator-Frage). <!-- d-check:ignore (Kennungen als Zitat) -->
