# Verification-Report: slice-v1-release-cut-v0.14.0 — 2026-10-23

**Rolle:** Verifier (Modul 11, [`harness/verification.md`](../../harness/verification.md)), frischer Kontext, kein Self-Review, kein Reparieren. Frage: Wurde der Release-Cut richtig gebaut, gemessen an Plan, DoD und ausgeführten Sensoren, **bis zur Push-Reife (M1)**? Die Punkte nach M1 (Tag, `publish`, GHCR, Tap-Job, CI auf dem Push-Tip) sind nicht Gegenstand dieses Berichts.

**Gegenstand:** [`slice-v1-release-cut-v0.14.0`](../plan/planning/in-progress/slice-v1-release-cut-v0.14.0.md), Range `f9680634..0721a30d` (neun Commits vor `origin/main`, nichts gepusht, kein Tag `v0.14*`), Vertrag [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten); Eingang: Plan §2/§3, Validierungsbericht [`2026-10-22-welle-devcontainer-shell-und-zweites-up-validierung`](2026-10-22-welle-devcontainer-shell-und-zweites-up-validierung.md), Review-Report [`2026-10-23-slice-v1-release-cut-v0.14.0`](2026-10-23-slice-v1-release-cut-v0.14.0.md) (Eingang, nicht Beleg), [`releasing.md`](../maintainer/releasing.md).

## Ergebnis

**GO für M1 (Push).** Kein Blocker. Zwei Hinweise für M2 (Abschnitt Befunde je Stufe), keine Änderung am Cut nötig. Alle Sensoren, deren Ausgabe ich nicht aus dem Bericht des Implementers übernehmen wollte, habe ich selbst gefahren (Host `linux/amd64`, Docker Engine 29.9.2).

**Push-Tip:** Dieser Report ist ein weiterer Commit auf `0721a30d`. Der Push-Tip, M2 und der Tag beziehen sich auf den Commit, der diesen Report trägt (Hash in der Übergabe an den Hauptagenten), nicht auf `0721a30d`.

## Sensoren — selbst gefahren

| Sensor | Ergebnis | Beleg |
| --- | --- | --- |
| `make gates` (auf `0721a30d`) | **Exit 0** | `lint` 0 issues; `coverage-gate` 92,00 % gegen 90 %; `docs-check` 722 Dateien, 0 Befunde; `baseline-verify` v6.18.1 OK, 54 Dateien; `ablage-pruefung` 26 Fälle, 0 Fehlschläge |
| `make ci` | **Exit 0** | `govulncheck` „No vulnerabilities found“; `image-scan` (Trivy 0.70.0): `u-boot:latest (debian 12.15)` 0, `usr/local/bin/u-boot` 0, **HIGH/CRITICAL gesamt 0** |
| `make fullbuild` | **Exit 0** | `ci` + Runtime-Image; danach `docker run --rm u-boot:latest --version` = `u-boot version 0.14.0-dev`, OCI-Label `0.14.0-dev` |
| `make build VERSION=0.14.0` | **Exit 0** | `--version` = `u-boot version 0.14.0`, OCI-Label `org.opencontainers.image.version` = `0.14.0` (lokales Image, danach durch `fullbuild` auf den Standardbau zurückgesetzt) |
| `make doc-trace` | Exit 0 | 84 Anforderungen, 0 Waisen |
| `make doc-immutable RANGE="v0.13.1..HEAD"` | Exit 0 | 127 Commits, 0 Befunde (keine Accepted-ADR umgeschrieben) |
| `make doc-commits RANGE="origin/main..HEAD"` | Exit 0 | 9 Commits, 0 Befunde |
| `make doc-planning` | Exit 0 | 722 Dateien, 0 Befunde |
| `make docs-check` (Tip, vor diesem Report) | Exit 0 | 722 Dateien, 0 Befunde; Arbeitsbaum sauber (`git status --short` leer) |
| `make verify-depguard` | nicht gelaufen | Grund: keine Importregel und kein Go-Import berührt; der einzige Code-Diff ist `var version` in `cmd/uboot/main.go` (Diff gelesen) |
| `make test-docker` (voll) | nicht wiederholt | Grund: Laufzeit über dem Paket-Zeitlimit von 30 Minuten (Befund des Implementers, Exit 2 wie in den Cuts `v0.13.x`); stattdessen die neun neuen Tests und `EgressAptUpdate` gezielt (unten) |

## Tap-Trockenlauf M-T — selbst wiederholt

Nach [`releasing.md`](../maintainer/releasing.md) §4: Image `u-boot:test`, `--network none`, Arbeitsbaum als `:ro`-Mount, Skript im Scratchpad (nicht im Repo), Fixture nach `sumsFull` in `cmd/uboot/homebrew_scripts_test.go`, Stellvertreter für `curl` und `git` im `PATH`.

| Schritt | Exit | Beobachtung |
| --- | --- | --- |
| `homebrew-formula-fill.sh v0.14.0` | 0 | vier `url`-Zeilen mit `releases/download/v0.14.0/` |
| `tap-formula-check.sh v0.14.0` | **0** | |
| `tap-formula-check.sh v0.13.1` (Gegenprobe) | **1** | „found 4 url lines, 0 of v0.13.1“ |
| `tap-nachzug.sh`, `TAG=v0.14.0`, Dummy-Token | 0 | Protokoll: `curl …/v0.14.0/u-boot.rb`, `git clone`, `git diff --quiet`, `git add`, `git commit -m u-boot 0.14.0`, `git push origin HEAD`; „updated to v0.14.0“ |
| `tap-nachzug.sh`, `TAG=v0.13.1` (Gegenprobe) | 1 | Protokoll nur der `curl`-Aufruf, **kein** `clone`, kein `push` |

**Grenze (wie in M-T benannt):** nicht bewiesen sind `is_stable == 'true'`, der Workflow-Schritt, der echte `curl`, der echte `git push` und das Token; das bleibt `tap-check` (M4) und dem ersten stabilen Lauf vorbehalten.

## Docker-Tests der Welle — gezielt gefahren

Befehl wie beauftragt, Image `u-boot:test-docker-tools` (`repo.mk`, Ziel `test-docker`), `-run` auf die Namen am Baum:

| Test | Datei | Ergebnis |
| --- | --- | --- |
| `TestShell_Docker_NoTerminal_UserWorkspaceAndEnvironment` | `internal/adapter/driven/docker/devcontainer_shell_docker_test.go` | PASS (5,05 s) |
| `TestShell_Docker_Probe` | dieselbe | PASS (4,95 s) |
| `TestShell_Docker_B18_SessionDoesNotBegin` | dieselbe | PASS (8,74 s) |
| `TestShell_Docker_Terminal_UnderScript` | dieselbe | PASS (3,44 s) |
| `TestShell_Docker_ClientEnd` | dieselbe | PASS (13,92 s) |
| `TestDevcontainerShell_Docker_Session` | `internal/e2e/devcontainer_shell_docker_test.go` | PASS (14,95 s) |
| `TestDevcontainerShell_Docker_FailedStatusOpensNoSession` | dieselbe | PASS (10,90 s) |
| `TestDevcontainerUp_Docker_S8_SecondUpAfterStartScriptAbort` | `internal/e2e/devcontainer_start_docker_test.go` | PASS (17,07 s) |
| `TestDevcontainerUp_Docker_AlreadyRunningLifecycleCounters` | dieselbe | PASS (9,62 s) |

Exit 0, **0 `--- SKIP`**, beide Pakete `ok`. **Zählung:** am Baum sind es **neun** Tests (fünf `TestShell_Docker_*` im Adapter, vier im Paket `internal/e2e`), nicht zehn; der Plan (M2) führt ebenfalls neun. Vorher/nachher: Container (`pgr-gp`, `pgr-dev2`, `funny_dhawan`, `awesome_shockley`) und Images unverändert, Volumes 645 vorher und nachher gleich; die vier Container wurden nicht angefasst.

## Befunde je Stufe

**Stufe Sensoren und Tests (e) — `EgressAptUpdate`.** `TestE2E_SandboxDevcontainer_NestedDocker_EgressAptUpdate` habe ich einmal gefahren: **PASS** (52,85 s, Paket `internal/e2e` `ok`); ein weiterer Lauf war nicht nötig, der Befund ist damit nicht reproduziert (Implementer: FAIL, FAIL, PASS; Netzfehler `SSL routines::unexpected eof` bei `download.docker.com`). **„Transient, ohne Code-Bezug“ ist belegt, über die Dateiliste, nicht über die Wortsuche** (das deckt sich mit dem LOW-Finding F-2 des Reviews): `git diff v0.13.1..HEAD --stat -- internal/e2e Dockerfile` zeigt in `internal/e2e` nur `devcontainer_shell_docker_test.go` (neu) und `devcontainer_start_docker_test.go`; `git diff v0.13.1..HEAD --stat -- internal/hexagon/application/templates internal/hexagon/application/devcontainer_egress.go internal/e2e/sandbox_devcontainer_docker_test.go` ist leer (Test, Dockerfile-Vorlage und Egress-Erzeugung unverändert); das `Dockerfile` des Repos hat nur die drei Versionszeilen. Eine Rate ist nicht belegt (ein Lauf von mir). **Was M2 prüfen muss:** `integration-docker` auf dem Push-Tip ohne `--- FAIL` und ohne `panic: test timed out` (ein roter `EgressAptUpdate` dort wäre erneut ein Netzbefund, nicht automatisch ein Code-Befund, aber die Rückführung nach `open` nach §4 des Plans bliebe); die **Dauer** des Jobs gegen `timeout-minutes: 30` (lokal läuft das Paket `internal/e2e` am Limit aus und erreicht die Tests hinter dem Ollama-Test nicht; **lokal ist daher nicht belegt, dass jeder Test des Pakets besteht**, nur die neun und `EgressAptUpdate`); die Zahl der `--- SKIP` (CI erwartet vier Nested-Podman-Tests, lokal liefen sie); „Operating System“ der Runner-Logs.

**Stufe Versionsstellen (b).** Wortsuche über den ganzen Baum (außer `docs/plan`, `docs/reviews`, `.harness`, CHANGELOG), jede Fundstelle gelesen: `cmd/uboot/main.go:43`, `repo.mk:51/54/55`, `Dockerfile:35/37/39` tragen `0.14.0-dev` (Kommentare nennen den geltenden Wert); `README.md`/`README.de.md`: Status, `VERSION=0.14.0`, Tabellenzeile `v0.14.0`; Handbuch: Kopf `v0.14.0` und `Stand: 2026-10-23`, `:82`, `:108`, `:135`, Gültigkeitsbereich. Übrig bleibt `0.13.x` nur in den Tabellenzeilen `v0.13.0`/`v0.13.1` (Historie), im CHANGELOG-Kopf `[0.13.1]` samt Vergleichslinks, in Handbuch `:955` (Messstand „`u-boot` 0.13.0, Docker 29.9.1“, als alter Messstand richtig) und in [`releasing.md`](../maintainer/releasing.md) Zeile 57 (Zeiger „zuletzt“ auf den Cut `v0.13.1`, Closure-Sache, nach Plan). Treffer in `spec/lastenheft.md`/`spec/spezifikation.md` sind Dokumentversionen, keine Release-Nummer. `.github/workflows/*`, `scripts/`, `packaging/`, `tools/` tragen keine Versionsnummer; die Workflows leiten die Version aus dem Tag ab und prüfen `--version` und OCI-Label gegen ihn (`publish.yml`, Schritt mit `IMAGE_VERSION`/`BINARY_VERSION`) — beides habe ich lokal mit `make build VERSION=0.14.0` nachgezeichnet. Keine übersehene Stelle. **Grenze:** kein Test pinnt die Nummer; eine falsche Versionszeile färbt kein Gate (Plan §3 Abweichung 6, bestätigt).

**Stufe Release-Text und Handbuch (d).** CHANGELOG-Lead (`[0.14.0] - 2026-10-23`, „Fourteenth release.“) gegen Validator-Bericht, Test oder Plan-Beleg-Tabelle gehalten:

- **A-1:** „wiederholtes `up` ohne `--recreate` führt `postCreateCommand`/`postStartCommand` nicht erneut aus, `postAttachCommand` bei jedem `up`“: gepinnt durch `TestDevcontainerUp_Docker_AlreadyRunningLifecycleCounters` (PASS, selbst gefahren). Exit `12` (ein `postStartCommand` endet mit 3) und Strg-C (Rückkehr nach 10,5 s) mit folgendem `up` Exit `0`/`already running`: Validator Punkt (ii), im Text als „one run each, Linux `amd64`, Docker 29.9.2“ gekennzeichnet; der Test `…S8_SecondUpAfterStartScriptAbort` (PASS) stützt den Exit-`12`-Fall. Ein gescheiterter Bau ist im Text ausdrücklich „not measured“.
- **A-2:** Zusage-Menge (Linux `amd64`, lokale Docker-Engine, Messbasis Docker 29.9.2) und sieben Nicht-Zusagen stehen im Lead; Terminal „measured under `script`“: `TestShell_Docker_Terminal_UnderScript` (PASS); ein Terminal eines Menschen ist als nicht gemessen ausgewiesen.
- **A-3:** „nach Behebung der Ursache“ im Handbuch an den zwei geplanten Stellen (Spiegelstrich zu `shell` bei Klasse 11; zweiter Spiegelstrich *Erneutes `up`*) und im Lead; die Aussage „bei fortbestehender Ursache wieder Exit 11“ ist als „ein Lauf, Linux `amd64`, Docker 29.9.2“ gekennzeichnet, Beleg Validator Punkt (ii). Sie hat keinen Test (Bedingung in Handbuch ohne Test, [`AGENTS.md`](../../AGENTS.md) §3.10, letzter Punkt): der Plan führt sie als Empfehlung mit Messbeleg „ein Lauf“ (§6) — **kein Befund**, die Aussage steht als Messung, nicht als Eigenschaft des Codes.
- `shell`-Aussagen (Exit `0` nach Beginn der Shell, `failed` öffnet keine Sitzung/Exit `11`, `--json`/`--dry-run`/Argumente Exit `2`): `TestDevcontainerShell_Docker_Session` und `…FailedStatusOpensNoSession` (PASS) und Validator Punkte (iii), (iv), (vii).
- Nicht im Text, wie Entscheidung (m) es festlegt: Grenze (e), CI-Satz im Lead (Plan §3 Abweichung 1 — die Status-Blöcke der READMEs behalten ihren bisherigen CI-Satz, der Lead sagt nichts zur CI).
- **Kennungsfreiheit (f):** in den hinzugefügten Zeilen von `README.md`, `README.de.md` und `docs/user/` 0 Treffer auf `LH-`, `ADR-`, `SPEC-`, `ARC-`, `slice-`, `welle`, `Roadmap`; im CHANGELOG-Lead nur „Lastenheft 0.21.0“ (erlaubt).
- Hinweis (INFO, kein Befund): Die Aussage der Status-Blöcke „the four tests of nested Podman are skipped there“ bleibt unverändert aus `v0.13.1` und ist für den Push-Tip erst durch M2 bestätigt; der Plan führt das als offenes Risiko.

**Stufe Plan-vs-Diff (AGENTS.md §3.10).** `git diff --stat f9680634..HEAD` zeigt 13 Dateien (vor diesem Report): `CHANGELOG.md`, `Dockerfile`, `README.md`, `README.de.md`, `cmd/uboot/main.go`, `repo.mk`, `docs/user/benutzerhandbuch.md`, `docs/plan/planning/in-progress/roadmap.md`, `docs/plan/planning/open/README.md`, die Welle-Datei, der Slice selbst, der Validierungsbericht (Eingang, vor dem Plan-Commit) und der Review-Report. Jede Datei hat eine Zeile in der Plan-Tabelle (§3); jeder Eintrag mit „update/neu“ hat einen Diff. Die als „kein Diff“ geführten (`releasing.md`, `spec/`, `docs/plan/adr/`, Workflows, `scripts/`, `packaging/`, `Makefile`, `go.mod`, `tools/`, `internal/`, übrige `cmd/uboot/`) haben keinen Diff (diffstat leer; `doc-immutable` Exit 0). Die Roadmap-Zeile nennt „Push und Tag stehen aus“ — ein Zustandsmarker mit Ablaufpunkt (M8), im Plan vorgesehen. Dieser Report ist die in der Plan-Tabelle geführte Datei `…-verify.md`. Abweichung: keine.

**Stufe Push-Reife (g).** Reihenfolge: `ff351704` Start, `8bb06e9f` reiner Rename (`open/` → `in-progress/`, 0 Änderungen), `d0cd12c5` Verweis-Nachzug (drei Dateien, 8 Zeilen), `51a1efe4` Start-Zeiger, `9c6fedf1` Cut, `08babc88` Umsetzung, `0721a30d` Review. Move und Nachzug liegen im selben Push; der Zwischenstand `8bb06e9f` allein ist nach dem Review docs-check-rot und darf nicht Spitze eines Push werden — die Spitze (`0721a30d` bzw. dieser Report) ist grün (`make docs-check` Exit 0). Datum: `date +%F` = 2026-10-23, `git log -1 --format=%cs` = 2026-10-23, an allen vier Stellen 2026-10-23; kein Tag `v0.14*` vorhanden.

## DoD-Abgleich (Punkt für Punkt, Stand vor M1)

- [x] **Release-Text und Handbuch** (DoD 1, im Slice abgehakt) — bestätigt: CHANGELOG, Versionsstrings an drei Stellen, beide READMEs, Handbuch (Kopf, Beispiele, Gültigkeitsbereich, A-3 an zwei Stellen), kennungsfrei; Belege siehe Stufen oben.
- [ ] **Tag und Kontrolle** (DoD 2) — vor M1 erfüllbar und erfüllt: `make gates`, `make ci`, `make fullbuild` grün (selbst gefahren); `make test-docker` ist gelaufen (Implementer: Exit 2 am Paket-Zeitlimit, 3136 PASS, 0 SKIP, 1 FAIL `EgressAptUpdate` transient; von mir nur gezielt wiederholt, siehe oben); M-T festgehalten und von mir wiederholt. **Offen bis nach M1/M5, korrekt unabgehakt:** `tap-check` (M4), `ci`/`integration-docker` auf dem Push-Tip (M2), Tag (M5), Kontrolle (M6) mit allen vier Unterpunkten.
- [ ] **Befunde des Validators haben einen Ausgang** (DoD 3) — A-1, A-2, A-3 im Text mit Beleg (bestätigt); N-1 bis N-6 als Folge mit Träger in §6 des Slice benannt; die Closure-Notiz mit den Ausgängen folgt nach M6. Nicht vor M1 abzuhaken.
- [ ] `make gates` grün (DoD 4) — am Tip grün (selbst gefahren); „nach der Closure“ steht aus.
- [ ] **Review durchgeführt** (DoD 5) — Report liegt vor ([`2026-10-23-slice-v1-release-cut-v0.14.0`](2026-10-23-slice-v1-release-cut-v0.14.0.md), GO, 0 HIGH, 0 MEDIUM); Verification-Report ist dieser. Der Haken gehört dem Planner.
- [ ] Doku-Update, Closure-Notiz, Beobachtungs-Register, Risiken-Ausgänge, Paarungen (DoD 6 bis 10) — Closure-Pflichten nach M6, vor M1 nicht fällig; Roadmap-Zeiger, `open/README.md` und Welle-Datei tragen den Ist-Stand „in `in-progress/`, Push und Tag stehen aus“.

**Traceability.** [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten): die neun Docker-Tests und die Handbuch-/CHANGELOG-Aussagen (oben); Release-Pfad [`ADR-0004`](../plan/adr/0004-ci-system.md) und [`ADR-0007`](../plan/adr/0007-distributionswege-ghcr.md): `make ci`, `make fullbuild`, `make build VERSION=0.14.0`; [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md): M-T.

**Carveouts:** neu keine, gelöst keine, unverändert alle (Inventar nicht berührt).

## Hinweise (keine Befunde)

1. **F-3 des Reviews bleibt wirksam:** Beim Lauf dieses Berichts (kurz nach Mitternacht Ortszeit) ist `date +%F` der 23., `date -u` noch der 22. Der Hauptagent hält `date +%F` und `git log -1 --format=%cs` vor M1 und vor M5 gegen die vier Datumsstellen; ein Tag vor 02:00 CEST zeigt im GitHub-Release das UTC-Datum des 22.
2. **Lokaler `test-docker`-Befund** ersetzt den CI-Lauf nicht (Plan M2): Das Paket `internal/e2e` erreicht lokal die Tests hinter dem Ollama-Test nicht; M2 ist die einzige Stelle, an der das ganze Paket besteht oder nicht.
3. **Zählung:** neun statt zehn neue Docker-Tests (Plan und Baum stimmen überein).

## Nicht geprüft

CI-Läufe auf dem Push-Tip (M2), `tap-check` (M4), Tag, `publish`, GHCR-Digest, Release-Assets, Job `tap` (nach M5); `make verify-depguard` (entfällt, Begründung oben); der volle `make test-docker` (Zeitlimit, Begründung oben); macOS/Colima, `arm64`, Docker Desktop, rootless, Terminal eines Menschen (nicht zugesagt, nicht verfügbar).

## Übergabe

**GO für M1** durch den Hauptagenten, nach ausdrücklicher Freigabe des Projektinhabers im Gespräch. Vor dem Push: `date +%F`, `make docs-check`, Arbeitsbaum sauber. Alle Commits der Range in einem Push (Move `8bb06e9f` und Nachzug `d0cd12c5` zusammen). Danach M2 mit den oben genannten Prüfpunkten; der Tag gehört auf den in M2 geprüften Commit.
