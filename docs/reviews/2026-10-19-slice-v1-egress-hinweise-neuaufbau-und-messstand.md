# Review-Report: slice-v1-egress-hinweise-neuaufbau-und-messstand — 2026-10-19

**Review-Art:** Code-/Doku-Review (Hinweis zum Neuaufbau bei `generate devcontainer`, gekürzte Subdomain-Texte, Spezifikation, Handbuch, CHANGELOG) gegen Plan, Entscheidungen und Hard Rules; Code, Aufrufer und Tests selbst gelesen, nicht aus dem Bericht des Implementers übernommen.

**Gegenstand:** `git diff -M 1e445d3..HEAD` (Commits `c772352`, `31a55b6`, `ee092eb`): 14 Dateien (13 plus die Slice-Datei), 284 Einfügungen, 25 Löschungen; eine neue Datei (`internal/hexagon/application/generate_egress_rebuild_test.go`), 13 geänderte.

**Skill:** `.harness/skills/reviewer.md` · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-19

**Eingangs-Kontext:**

- Slice [`slice-v1-egress-hinweise-neuaufbau-und-messstand`](../plan/planning/done/slice-v1-egress-hinweise-neuaufbau-und-messstand.md) (§1, §2 DoD, §3 Plan, §4, §6); Entscheidungen des Projektinhabers vom 2026-10-19: (a) Subdomain-Texte auf das Gemessene kürzen, ohne Vorbehalt und Messstand; (b) der Neuaufbau-Hinweis entfällt unter `--quiet`.
- [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes); `AGENTS.md` §3.4, §3.7, §3.9, §3.10, §5 Regel 6; `harness/review.md`, `harness/replay.md`.
- Selbst gelesen: `egressRebuildWarnings` und `generateDevcontainer` (`internal/hexagon/application/generate.go`), `Generate` (Wahl des Dateisystems je `PreviewMode`), `existingContainerFallbackWarnings` als Schwester-Guard, `visibleWarnings` (`internal/adapter/driving/cli/warnings.go`), `printGenerateSummary` und `writeGenerateJSON` (`internal/adapter/driving/cli/generate.go`), alle Aufrufer von `planDevcontainerFiles` und `generateDevcontainer` (nur `Generate`; `init`, `add` und `config set` rufen den Pfad nicht), `egress.go`, `sandbox.go` (`NewEgressHost`), die Tests dieses Diffs, `TestE2E_SandboxDevcontainer_EgressSubdomains` samt Probeskript, die Evidence-Datei und das Protokoll der Wildcard-Messung, [`SPEC-019`](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) (*Namen und Subdomains*), [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), [`SPEC-046`](../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung), Handbuch (Abschnitt Egress) und CHANGELOG `[Unreleased]`.
- Sensor-Läufe des Reviewers (nur `make`): `make lint` (Exit 0), `make test` (Exit 0), `make docs-check` (510 Dateien, 0 Befunde), `make doc-trace` (83 Anforderungen, 0 Waisen), `make doc-planning` (0 Befunde), `make ablage-pruefung` (Exit 0, 26 Fälle, 0 Fehlschläge). `make gates` nicht erneut gefahren.

---

## Findings

Keine HIGH-Findings; der Konflikt-Pfad ist nicht ausgelöst. Zwei LOW, vier INFO.

LOW `spec/spezifikation.md` · "Der Hinweis erscheint nicht beim ersten Schreiben der Datei" - Lücke beim ersten Aktivieren nicht ausgewiesen
Klasse: Zusage-Grenze nicht benannt
Quelle: [`AGENTS.md`](../../AGENTS.md) §3.10 (Bedingung in Spec und Handbuch nennt ihre Grenze), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)
Befund: Der Hinweis entsteht nur, wenn der verwaltete Block eines **bestehenden** `egress-init.sh` ersetzt wird (`p.action == devcontainerActionReplaceBlock`). Wer die Egress-Restriktion in einem Projekt mit bestehendem Container erstmals aktiviert, schreibt die Datei zum ersten Mal und liest keinen Hinweis, obwohl auch dort ein Neuaufbau nötig ist (die neuen `runArgs` und das Dockerfile wirken erst im neuen Container); das Abschalten (die Datei wird nicht mehr geplant) gibt ebenfalls keinen. [`SPEC-046`](../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung), Handbuch und CHANGELOG nennen die Fälle „nicht beim ersten Schreiben, nicht ohne Änderung, nicht unter `--dry-run`“ und „der Befehl prüft nicht, ob ein Container existiert“, sagen aber nicht, dass das erste Aktivieren in einem Projekt mit Container ohne Hinweis bleibt.
Risiko: Ein Nutzer mit laufendem Container aktiviert die Restriktion, sieht keinen Hinweis und liest das als „wirksam“; der Container läuft ohne Restriktion weiter (`doctor` weist den Zustand beim Start aus, [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), aber nur, wer ihn ruft). Der Plan (DoD, §6) beschränkt den Hinweis auf die Änderung von `egress-init.sh`; die Lücke ist Scope-konform, nur nicht ehrlich benannt.
Verifizierbar: Review-only (Lesen von `egressRebuildWarnings` und der drei Dokumentstellen)

LOW `internal/hexagon/domain/sandbox.go` · "an allowed name already covers its subdomains" - Fehlermeldung sagt mehr als die gepinnte Tiefe
Klasse: Zusage weiter als Messung
Quelle: [`AGENTS.md`](../../AGENTS.md) §3.10, Beobachtung [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md)
Befund: Die Langhilfe und `Notes:` begrenzen die Aussage auf „one and two levels below it“ (gepinnt: `www.example.org` unter `example.org`, `s3.us-east-1.amazonaws.com` unter `amazonaws.com`, zwei Ebenen). Die Meldung zu einem `*`-Eintrag („an allowed name already covers its subdomains“, zweimal in `NewEgressHost`) sagt „its subdomains“ ohne Tiefe; drei und mehr Ebenen sind weder gemessen noch gepinnt (die Messung führt dort nur einen nicht existierenden Namen, `NXDOMAIN`).
Risiko: Der nächste Änderer oder Nutzer liest die Meldung als Zusage für jede Tiefe, während dieselbe Ausgabe zwei Ebenen nennt; kein Gate widerlegt die Lesart.
Verifizierbar: Review-only (Vergleich der Meldung mit `TestE2E_SandboxDevcontainer_EgressSubdomains` und dem Protokoll)

INFO `internal/e2e/sandbox_devcontainer_docker_test.go` · "makes s3.us-east-1.amazonaws.com (two levels below it) reachable" - Pin der Subdomain-Tiefe liegt hinter dem Build-Tag
Klasse: Pin außerhalb der Standard-Sensoren
Quelle: Slice §2 DoD (Tests pinnen den Wortlaut); Zuständigkeit Verifier
Befund: Die Behauptungen „one and two levels“ und „does not open its parent name … by name“ sind durch diesen Docker-Test belegt (Status `000` für `example.org` bei Eintrag `www.example.org`; zwei Ebenen unter `amazonaws.com`); ich habe ihn gelesen, nicht gefahren. `make test` pinnt nur den Wortlaut der Ausgabe (`TestExecute_EgressList_Notes`, `TestConfigSet_EgressAllowWildcard_Exit10WithHint`, `TestNewEgressHost`). Der Text sagt nicht mehr als der Test: „by name“ ist die Grenze des Status `000`, und die geteilte Adresse bleibt ungesagt; ausgenommen die Tiefe der Meldung (siehe oben).
Risiko: Keines für den Diff; `make test-docker` ist die Bestätigung und läuft beim Verifier oder im CI.
Verifizierbar: `make test-docker`

INFO `internal/hexagon/application/generate.go` · "changed; an existing container keeps the egress list of the image it was built from" - Hinweis folgt der Änderung des Skripts, nicht der Liste
Klasse: Auslöser weiter als Wortlaut
Quelle: Slice §1 (Auslöser „Änderung der Liste“)
Befund: Der Auslöser ist der ersetzte Block von `egress-init.sh`. Ändert sich das Skript, ohne dass sich die Liste ändert (etwa weil eine neue Version von `u-boot` den Block der Vorlage anders rendert), erscheint der Hinweis ebenfalls; die Meldung sagt „changed … rebuild … for the new list to take effect“. Spezifikation („ersetzt den verwalteten Block“) und Handbuch („ändert … ein bestehendes `egress-init.sh`“) treffen den Auslöser genau; nur die Meldung und der Titel des Slice sagen „Liste“. Ein Neuaufbau ist in diesem Fall ebenfalls richtig, weil der Container das Skript des Images behält.
Risiko: Gering; eine Lesart „meine Liste hat sich geändert“ kann falsch sein.
Verifizierbar: Review-only

INFO `internal/adapter/driving/cli/quiet_test.go` · "TestExecute_Generate_EgressRebuildNotice_QuietAndJSON" - Pin-Zuschnitt
Klasse: Pin in zwei Hälften
Quelle: Review-Linse *Tests*
Befund: Die Wirkung unter `--quiet` ist in zwei Tests gepinnt, die zusammen die Zusage tragen: Der Anwendungstest prüft, dass der Eintrag weder `SecurityRelaxation` noch `DegradationFallback` trägt; der CLI-Test prüft mit einem Fake-Use-Case und handgebautem Eintrag, dass `--quiet` ihn streicht und `--json` (auch mit `--quiet`) Status `warn`, Exit-Code 0 und genau einen `diagnostics`-Eintrag trägt. Kein Test fährt den echten Dienst durch den CLI-Pfad. Nicht gepinnt, aber in Spezifikation, Handbuch und CHANGELOG auch nicht zugesagt: `--diff` (`PreviewAndApply` schreibt durch, der Hinweis erscheint nach Code und Plan), `--dry-run --json`. Ein Golden Case entsteht nicht; `harness/replay.md` führt die beiden Tests als Evidence, bestehende Golden Cases sind unverändert (kein Golden im Diff).
Risiko: Eine Änderung der Kennzeichen im Dienst und eine im Filter fallen nur zusammen auf; beide sind einzeln gepinnt.
Verifizierbar: `make test`

INFO `internal/adapter/driving/cli/egress.go` · "(measured with the Dev Container CLI 0.89.0:" - Messstand-Klammer im Rest der Notes
Klasse: Messstand in Meldung
Quelle: Zielbild-Regel (Spezifikation = Zielbild, Meldungen und Handbuch = Ist-Zustand); Entscheidung (a) betrifft die Subdomain-Texte
Befund: Die Klammer (Dev Container CLI 0.89.0, Colima, arm64, Aktionen des Editors nicht gemessen) gehört zur Aussage „Container behält die Liste“ und steht jetzt hinter dem Satz „so a changed list takes effect after a rebuild of the container“. Sie belegt den zweiten Teil (`--remove-existing-container` wendet die neue Liste an) und bleibt wahr; sie ist vertretbar, weil die Subdomain-Texte keinen Messstand mehr tragen und die Messung für den Container-Teil im Handbuch dieselbe ist. [`SPEC-046`](../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung) trägt keinen Messstand (Zielbild), die Zusage „Ausgabe nennt den Neuaufbau“ ist durch `TestExecute_EgressList_Notes` gepinnt.
Risiko: Keines.
Verifizierbar: Review-only

---

## Negativbefunde

Geprüft, ohne Befund:

- **Auslöser des Hinweises** (`egressRebuildWarnings`): erstes Schreiben (Aktion `WriteNew`) gibt keinen, `NoOp` keinen, Änderung nur von `devcontainer.json` keinen (gepinnt: `…_NotForOtherFiles`), `--dry-run` keinen (Guard `PreviewDryRun`, gepinnt samt unverändertem Skript und anschließendem echten Lauf mit genau einem Hinweis); `--diff` (`PreviewAndApply`) schreibt durch und gibt ihn (nicht gepinnt, nicht zugesagt, siehe INFO). Die Pläne entstehen vor dem Schreiben aus dem Zustand der Datei; der Guard liest den Modus des Aufrufs, wie der Schwester-Guard `existingContainerFallbackWarnings`.
- **Aufrufer und Doppelmeldung:** `generateDevcontainer` ist nur über `Generate` erreichbar; `init`, `add`, `config set` und `generate` ohne Artefakt `devcontainer` rufen den Pfad nicht, ein zweiter Hinweis aus einem anderen Pfad entsteht nicht. Der Eintrag trägt als einziger das Subjekt `.devcontainer/egress-init.sh`; die Meldungen zur Lockerung und zum Fallback tragen denselben Code [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), aber andere Subjekte und Wortlaute. Der Test filtert nach Subjekt und zählt genau einen.
- **Code und Fehlermodell (Dual-Classifier-Regel):** Der Code [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) ist als `WarningEntry`-Code bereits in Gebrauch (Lockerung, Fallback); es entsteht weder ein Sentinel noch ein neuer Diagnose-Code, kein Exit-Code ändert sich (Warnung: Exit 0, `status: warn`), die zentrale Exit-Code-Klassifikation und die Abbildung des Subkommandos `generate` bleiben unverändert und werden von keinem Pfad dieses Diffs berührt. Der Umschlag-Status `warn` bei Exit 0 ist gepinnt.
- **`--quiet`, `--json`:** Der Eintrag trägt keines der zwei Kennzeichen und entfällt damit nach `visibleWarnings`; mit `--json` bleibt er in `diagnostics[]` (auch mit `--quiet`); Entscheidung (b) und [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) („vier Ausnahmen“, Lastenheft unverändert) sind gewahrt, weil der Hinweis keine fünfte Ausnahme bildet; [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) führt ihn unter den entfallenden.
- **Wirksamkeit der Tests (Mutationen gelesen, nicht gefahren):** Bedingung `ReplaceBlock` auf „jede Aktion“ → fällt im Fresh-State-Fall und im NoOp-Fall; auf „jede Datei“ → fällt in `…_NotForOtherFiles`; Guard `PreviewDryRun` entfernt → fällt in `…_DryRunGivesNone` (der Plan enthält dort `ReplaceBlock`); `SecurityRelaxation` oder `DegradationFallback` gesetzt → fällt im Kennzeichen-Test; Wort „rebuild“ oder der Aufruf `--remove-existing-container` aus der Meldung entfernt → fällt im Wortlaut-Test; Subjekt oder Code geändert → fällt (`Subject` ist Filter, `Code` ist geprüft). Der Wortlaut der Subdomain-Texte ist an drei Stellen gepinnt (Notes, Wildcard-Fehler im Adapter, Wildcard-Fehler in der Domäne, je mit Gegenfall „not a lowercase“ beziehungsweise `wantNot`).
- **Sprache:** Meldung, Langhilfe, `Notes:`, Fehlertexte und CHANGELOG-Einträge sind Englisch; Spezifikation und Plan Deutsch.
- **Spec-Straten (`AGENTS.md` §3.4):** Die neuen Sätze in [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) und [`SPEC-046`](../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung) verweisen nur aufwärts (Lastenheft) oder gleichrangig ([`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), [`SPEC-046`](../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung)); kein Slice, kein ADR, kein Carveout in `spec/`; `spec/architecture.md` unverändert. Das Lastenheft ist unverändert.
- **Ein Vertrag an vier Stellen:** Hinweis und Texte stehen gleichlautend in Code, [`SPEC-046`](../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung)/[`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), Handbuch und CHANGELOG (Auslöser, Gegenfälle, `--quiet`, `--json`, „prüft nicht, ob ein Container existiert“, „eine und zwei Ebenen“, „bei der Verbindung per Name“); `harness/replay.md` führt die Zeile mit den Tests. Das Handbuch trägt keine Kennungen und keine Verweise auf Planungsdokumente (`AGENTS.md` §5 Regel 6); die CHANGELOG-Einträge tragen verlinkte Kennungen.
- **Spezifikation ohne Messstand:** [`SPEC-046`](../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung) sagt nur, was ein Test pinnt (Ausgabe nennt den Neuaufbau: `TestExecute_EgressList_Notes`; Hinweis: Anwendungs- und CLI-Test); die Aussage „ob ein Container existiert, prüft der Befehl nicht“ ist Code-Aussage (keine Inspektion im Pfad) und im Handbuch wiederholt.
- **Architektur:** Kein neuer Import, kein I/O in `application` (reine Funktion über den Plan), keine Port-Änderung (`internal/hexagon/port/driving/` und `internal/adapter/driving/cli/generate.go` im Diff nicht enthalten); `make lint` (depguard eingeschlossen) grün; `make verify-depguard` nicht nötig (Importregeln unberührt).
- **Kommentare (`AGENTS.md` §3.7, §3.9):** Der Doc-Kommentar von `egressRebuildWarnings` und der Kommentar der Konstante beschreiben Zusage und Grenze im Indikativ, ohne Chronik; der Helfer steht nach `generateDevcontainer` und vor `existingContainerFallbackWarnings` mit eigenem Doc-Kommentar (kein Block-Bruch); die neue Konstante `egressInitRel` ersetzt das Literal in `planDevcontainerFiles` (kein weiteres Vorkommen im Quelltext); die alten Texte („already included“, „subdomains of an allowed name“, „never opens its parent name“) kommen im Baum außerhalb von Plan, Register, Reviews und alten CHANGELOG-Einträgen nicht mehr vor. Der Testkommentar in `sandbox_devcontainer_docker_test.go` nennt die drei gepinnten Fälle.
- **Gate-Lockerung, Suppression:** kein `//nolint`, keine Änderung an `.golangci.yml`, Workflows oder Coverage.
- **Plan-Tabelle gegen `git diff --stat -M 1e445d3..HEAD` (`AGENTS.md` §3.10):** alle 13 Dateien haben eine Zeile oder fallen unter die Sammelzeile (Spezifikation, Handbuch, CHANGELOG, `harness/replay.md`); die als unverändert gelesenen Dateien (`generate.go` des Adapters, `generate_acceptance_test.go`, `internal/hexagon/port/driving/`) fehlen im Diff. Die Zeile zum E2E-Kommentar nennt noch „Zeile 197“; der Kommentar steht nach der Änderung wenige Zeilen davor, die Datei bleibt richtig benannt.

Nicht geprüft:

- `make test-docker` (Docker-Pfad nur durch einen Kommentar berührt; der Test, der die Subdomain-Aussagen pinnt, ist gelesen, nicht gefahren; siehe INFO).
- `make verify-depguard` (Importregeln unberührt).
- `make gates` (nicht erneut nötig; die Einzelsensoren `lint`, `test`, `docs-check` liefen grün).
- Wirkung des Hinweises in einem echten Editor-Lauf (Neuaufbau durch Aktion des Editors): nicht gemessen, nicht zugesagt.
- DoD-Closure, Closure-Notiz, Register-Fortschreibung (`BEO-CLI/subdomain-text-ohne-vorbehalt` bleibt offen, bis der Verifier schließt), Ausgang der drei Risiken in §6 des Slice: Verifier-Aufgabe, `[ ]` der DoD sind offen.

## Zusammenfassung

Keine HIGH, keine MEDIUM; zwei LOW, vier INFO. Der Hinweis entsteht genau im beschriebenen Fall, ist durch Anwendungs- und CLI-Tests gegen die naheliegenden Mutationen gepinnt, erfindet weder Code noch Sentinel und ändert keinen Exit-Code. Die gekürzten Texte sagen, bis auf die Tiefe der Wildcard-Meldung, nicht mehr als der Docker-Test und das Protokoll tragen. Offen für den Implementer sind die Grenze der Zusage beim ersten Aktivieren (Spezifikation, Handbuch, CHANGELOG) und die ungebundene „subdomains“-Aussage der Wildcard-Meldung.
