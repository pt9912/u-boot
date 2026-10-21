# Verification-Report: [`slice-v1-release-cut-v0.13.0`](../plan/planning/done/slice-v1-release-cut-v0.13.0.md) — 2026-10-21

**Rolle:** Verifier (frischer Kontext, kein Self-Review, repariert nichts). Frage: Wurde der Cut richtig gebaut, gemessen an Plan, DoD, Entscheidungen des Projektinhabers und Sensors? Form nach [`harness/verification.md`](../../harness/verification.md). Kein Inhalt geändert außer dieser Datei, nichts committet, nichts gepusht, kein Tag, kein `gh workflow run`; `gh` nur lesend.

**Gegenstand:** Stand **vor dem Push (M1)**. `HEAD` = `5ec45eee`, `origin/main` = `3e7799b6`, acht Commits darüber (`c26cfd43` Plan, `53a7cdd6` reiner Move, `6084a110` Link-Nachzug, `736d894a` Start, `43760c4d` L-10, `9f90e9da` Cut, `62ccd9cd` §3 Umsetzung, `5ec45eee` Review-Fix-Runde samt Review-Report); `git diff --stat 3e7799b6..HEAD`: 16 Dateien, 531 Einfügungen, 41 Löschungen. Kein Tag `v0.13*` lokal (`git tag -l`) und auf `origin` (`git ls-remote --tags origin`: nur `v0.12.0`). Arbeitsbaum vor dieser Datei sauber.

**Eingang (als Eingang gelesen, nicht als Beweis):** Slice §1 bis §8 samt *Umsetzung*, *Review-Fix-Runde* und den sieben Entscheidungen des Projektinhabers vom 2026-10-21 (§7); Review-Report [`2026-10-21-slice-v1-release-cut-v0.13.0`](2026-10-21-slice-v1-release-cut-v0.13.0.md) (0 HIGH, 0 MEDIUM, 2 LOW, 6 INFO); [`releasing.md`](../maintainer/releasing.md) §1 bis §8; Vorbild [`slice-v1-release-cut-v0.12.0`](../plan/planning/done/slice-v1-release-cut-v0.12.0.md) und [Verification v0.12.0](2026-10-19-verification-slice-v1-release-cut-v0.12.0.md); [Validator-Bericht der Welle](2026-10-21-validation-welle-devcontainer-start-ohne-editor.md).

**Gesamturteil: GO für M1 (Push von `main`)** — nach ausdrücklicher Freigabe des Projektinhabers im Gespräch und mit den Auflagen 1 und 2 unten. Keine DoD-Verletzung im vor dem Push erfüllbaren Teil; keine Aussage im Release-Text ohne Beleg gefunden; sieben INFO-Befunde (V-1 bis V-7), keiner blockiert.

## Verification Evidence

Scope:
- Slice: [`slice-v1-release-cut-v0.13.0`](../plan/planning/done/slice-v1-release-cut-v0.13.0.md)
- IDs: [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-FA-DIAG-002`](../../spec/lastenheft.md#lh-fa-diag-002--lokale-voraussetzungen-prüfen), [ADR-0028](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md), [ADR-0029](../plan/adr/0029-devcontainer-start-exit-11-und-socket-haltung.md), [ADR-0026](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md), [ADR-0027](../plan/adr/0027-egress-init-wiederholungslauf-atomar.md), [ADR-0004](../plan/adr/0004-ci-system.md), [ADR-0007](../plan/adr/0007-distributionswege-ghcr.md), [ADR-0016](../plan/adr/0016-homebrew-distribution-per-tap.md), [ADR-0017](../plan/adr/0017-linux-pakete-nfpm.md)
- Artefakte: `CHANGELOG.md`, `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md`, `cmd/uboot/main.go`, `repo.mk`, `Dockerfile`, `.github/workflows/ci.yml`, `docs/maintainer/quality.md`, `docs/maintainer/releasing.md`, Roadmap, `open/README.md`, Colima-Folge-Slice, Vermerk im Welle-Kopf, Review-Report

### DoD-Abgleich (§2 des Slice; Stand vor M1)

- [x] **Release-Text und Handbuch** — erfüllt im vor dem Tag erfüllbaren Teil (Einzelnachweise unten). `## [0.13.0] - 2026-10-21` mit Lead nach Auftrag 3 und 5, Vergleichslinks, Versionsstrings `0.13.0-dev` an den drei Stellen, beide READMEs (Status, Tabellenzeile, `VERSION=0.13.0`, L-7), Handbuch (Kopf, drei Versionsbeispiele, Gültigkeitsbereich, L-1 in 4.9 und 4.11, L-2 in 4.9 und §7), keine Kennungen in den hinzugefügten Zeilen von READMEs und Handbuch. **Zwei Bedingungen bleiben am Tag hängen:** das Datum (M0 am Tag-Tag) und der CI-Satz (M2 auf dem Push-Tip).
- [ ] **L-10** — Form erfüllt, Beleg offen (M2). `ci.yml` `with.version: 'v0.70.0'` (Zeile 122), Kommentar in `ci.yml`, `repo.mk` Zeile 60 und `quality.md` Zeile 253 nennen die Eingabe `version`; die gepinnte Action kennt genau diese Eingabe (selbst gelesen, Einzelnachweise). Der Beleg „image-scan grün, keine Annotation *Unexpected input(s)*, Log `version: v0.70.0`“ entsteht erst auf dem Push-Tip; V-4 schärft, welcher Teil davon trennt.
- [ ] **Tag und Kontrolle** — vor dem Tag erfüllbarer Teil: `make gates`, `make ci`, `make fullbuild` **selbst gefahren, Exit 0** (Sensors); `make test-docker` ist gelaufen und steht mit Zahlen in §3 (Implementer: Exit 2 am Paket-Zeitlimit, 2833 PASS, 0 SKIP, 0 FAIL) — erfüllt im Wortlaut der DoD, die kein Grün verlangt; Nachweis für die Release-Aussage ist M2 (Review F-3). M-T festgehalten und **selbst nachgefahren, gleiches Ergebnis**. **Offen:** `tap-check` (M4), `ci` und `integration-docker` auf dem Push-Tip (M2), Tag (M5), die fünf Kontrollzeilen von M6.
- [x] **`make gates` grün nach dem letzten Commit vor dem Handoff** — selbst auf `5ec45eee` (Exit 0, Stempel `.harness/state/gates-passed.head` = `5ec45eee…`). „Nach der Closure“ offen.
- [x] **Review durchgeführt, Report liegt vor** — [`2026-10-21-slice-v1-release-cut-v0.13.0`](2026-10-21-slice-v1-release-cut-v0.13.0.md), im Repo seit `5ec45eee`; Findings einzeln unten. Verification-Report getrennt: diese Datei, **vor M1 zu committen** (Auflage 1).
- [~] **Doku-Update** — CHANGELOG, READMEs, Handbuch, `releasing.md` §3 („zuletzt“ → [`slice-v1-release-cut-v0.12.0`](../plan/planning/done/slice-v1-release-cut-v0.12.0.md)) und die abgelaufenen Zustandsmarker in den lebenden Dokumenten (Roadmap-Zeile der Welle, `open/README.md` Zeilen 42 und 43, Colima-Slice §4, Vermerk im Welle-Kopf) tragen `v0.13.0`. **Offen bis nach dem Tag:** Roadmap-Zeile `v0.13.0` mit Tag-Hash (M8). Die Stellen in `done/` bleiben nach Review F-8.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag — offen (Closure nach M6).
- [ ] Beobachtungs-Register — offen (Closure; `state.md` von [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../plan/planning/observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md) steht in der Plan-Tabelle).
- [ ] Jedes Risiko aus §6 mit Ausgang — offen (Closure); Eingaben dafür in V-2 und V-3.
- [ ] Drei Paarungen — offen (Closure).

**M0 und Vorab-Prüfung nach [`releasing.md`](../maintainer/releasing.md) §4 und §5:**

- M0: `date +%F` = `2026-10-21` = `CHANGELOG.md` `## [0.13.0] - 2026-10-21` = Tabellenzeile `README.md` = Tabellenzeile `README.de.md` = Handbuch `Stand: 2026-10-21` — vier Stellen gleich. Gilt nur für einen Tag am 2026-10-21 (Auflage 2).
- §4 Make-Ziele: `gates`, `ci`, `fullbuild` selbst grün; `test-docker` lokal nicht erneut gefahren (Auftrag; Nachweis M2).
- §4 Trockenlauf des Tap-Pfads: selbst grün, Grenze wie benannt.
- §4 „`main` ist gepusht und der letzte CI-Lauf auf `main` grün“: **offen** (M1, M2). Der letzte Lauf auf `origin/main` (`3e7799b6`) ist grün (V-3).

### Einzelnachweise (selbst gelesen, nicht aus dem Slice übernommen)

- **CHANGELOG `[0.13.0]`** (`git diff 3e7799b6..HEAD -- CHANGELOG.md`): Der alte Block `[Unreleased]` steht vollständig unter `[0.13.0]`; geändert sind im *Added*-Eintrag zu `devcontainer up` genau drei Stellen (Satz zu `containerEnv`/`remoteEnv` für beide Formen, Exit-Code `14` in der Aufzählung, ein neuer Satz zu L-1 und L-2), dazu der Lead und die Vergleichslinks (`[Unreleased]: …compare/v0.13.0...HEAD`, neu `[0.13.0]: …compare/v0.12.0...v0.13.0`); leerer Anker `## [Unreleased]`. „Thirteenth release.“ folgt der Zählung der Minor-Releases (Vorgänger „Twelfth release.“; Patch `v0.8.1` zählt dort nicht mit).
- **Die Sätze der Fix-Runde F-1**, Wort für Wort gegen Handbuch 4.9 *Erneutes `up` nach einem Fehler oder Abbruch* (Zeilen 874 bis 883) und den Validator-Bericht (Tabelle *Fehlerfälle*, Zeilen 49 bis 52):
  - Strg-C nur „during a `postCreateCommand`“ — Validator: `demo-int`, `postCreateCommand: sleep 120`, SIGINT; „Ein folgendes `up` endet in 1,1 s mit Exit 0, ohne das abgebrochene `postCreateCommand` zu wiederholen“, „der Devcontainer läuft“. Lead und *Added* binden „the failed `postStartCommand`“ und „the aborted `postCreateCommand`“ je an ihren Fall; die Verallgemeinerung aus F-1 ist weg.
  - „so this exit code `0` does not mean that the lifecycle commands completed“ — Handbuch Zeilen 880 bis 882 („Ein Exit-Code 0 des zweiten `up` heißt in diesen Fällen nicht, dass die Lebenszyklus-Befehle durchgelaufen sind“). Inhaltsgleich, auf diese Fälle begrenzt.
  - „named volumes stay“ — Handbuch Zeilen 846 bis 847 und 882 bis 883; gepinnt durch `TestDevcontainerUp_Docker_S1_S4_SandboxProfileVolumes` (`internal/e2e/devcontainer_start_docker_test.go`, Kopfkommentar: „a marker in the workspace volume survives `down`, a following `up` and `up --recreate`“).
  - „Linux `amd64` with Docker 29.9.1“ — Umgebung des Validators (Bericht Zeile 15: native Docker Engine 29.9.1); sie trägt die Beobachtungen zu Exit 12 und Strg-C. Der Satz zu `docker exec -e` nennt keine Docker-Version, er sagt „one run each“; dahinter steht die Gegenprobe des Implementers auf Docker 29.9.2 (Slice §3), das Handbuch nennt sie mit 29.9.2. Kein Widerspruch.
  - „other failures, such as a failed build, and exit code `11` from a start script not measured“ — gleich Handbuch Zeilen 879 bis 880 und Beleg-Tabelle („nur Simulation `demo-egfail`“).
- **Weitere Zahlen und Bedingungen des Lead**, je an der Fundstelle gelesen: 1716 MiB ≈ 1,7 GiB und 597 HIGH/CRITICAL-Funde im Betriebssystem-Anteil des Node-Images (Handbuch Zeilen 762 bis 768; Validator Schritt 3: „≈ 1716 MiB“); „about ten seconds“ nach Strg-C (Handbuch 868 bis 870: 10,1 s CI, 10,5 s lokal; Validator 10,35 s); Maskierung ab vier Byte (796 bis 798); Benutzer-ID (885 bis 886); Exit-Codes `10`/`11`/`12`/`14`/`2` (Tabelle 853 bis 859; Klasse `14` gepinnt in `internal/adapter/driving/cli/devcontainer_test.go` Fall `file system`, Erwartung 14); Verhaltensänderungen aus den *Changed*-Einträgen (Zeilen 48 bis 54: Exit `11` auch unter `warn` in vier Fällen, `podman run --dns`, 36 s statt 20 ms, Messstand „Docker 29.9.1 in two runs“, Transaktion „pinned by a test“); „Lastenheft 0.18.0“ = `spec/lastenheft.md` Zeile 9; CI-Satz „four tests of nested Podman were skipped“ = vier Aufrufer von `skipIfNestedUsernsRestricted` (`sandbox_devcontainer_docker_test.go` Zeilen 691, 737, 778; `sandbox_status_docker_test.go` Zeile 173), Skip-Grund `kernel.apparmor_restrict_unprivileged_userns=1` (Zeile 653), und die zwei letzten CI-Läufe (V-3).
- **Handbuch 4.9, Exit-Code-Tabelle Zeile `11` (Fix F-5)** gegen den Code:
  - `egress-init.sh.tmpl`: vier Abbruchzweige mit `write_status failed` und `exit 11` **ohne** Abfrage von `on_unavailable` — Benutzer `nobody` fehlt (Zeilen 41 bis 45), `nft -f` scheitert (97 bis 101), `dnsmasq` startet nicht (120 bis 124), Resolver unter anderer UID (127 bis 131); nur fehlendes `NET_ADMIN` hängt an `fail` (25 bis 33). Das deckt „das Startscript der Egress-Beschränkung auch unter `warn`, wenn der Benutzer `nobody` im Image fehlt, die nftables-Regeln sich nicht anwenden lassen oder der lokale DNS-Server nicht startet oder nicht unter diesem Benutzer läuft“ vollständig und ohne Überschuss.
  - `sandbox-init.sh.tmpl` und `docker-init.sh.tmpl`: `unavailable()` schreibt `failed` und endet mit `11` nur bei `fail`, sonst `fallback` (Zeilen 17 bis 24 bzw. 21 bis 28). Der Zweig „nested user namespaces are blocked“ in `sandbox-init.sh.tmpl` (Zeilen 27 bis 30) endet mit `11`, schreibt aber keine Statusdatei; er fällt damit nicht unter „meldet `failed` in seiner Statusdatei“ — die Tabellenzeile behauptet für ihn nichts, richtig.
  - `startScriptFailed` (`internal/hexagon/application/devcontainer_start.go` Zeilen 195 bis 205) prüft nur `result=failed`; Abbildung auf `11` nur zusammen mit `ErrDevcontainerRuntime` (Zeile 131).
  - Tests: `TestDevcontainerUp_RuntimeFailureClassification` Fall `egress failed` (`devcontainer_start_test.go` Zeile 265); `TestE2E_SandboxDevcontainer_EgressUpstreamDNS` Unterfall `abort branches` (drei Zweige, Projekt mit `warn`) und `TestE2E_SandboxDevcontainer_EgressRerun` Unterfall `failed transaction` (vierter Zweig, ebenfalls `warn`). Wortlaut gleich §7 (Zeilen 2351 bis 2353) und *Fallbacks sehen* (Zeile 1290). Kennungsfrei.
- **Handbuch L-1 (4.9 *Im Container arbeiten*, 4.11):** `GIT_TOKEN` steht im Template nur unter `remoteEnv` (`devcontainer.json.tmpl` Zeilen 33 bis 34); `docker exec` wendet `remoteEnv` nicht an. Die Docker-Semantik von `-e` habe ich selbst nachgefahren (Wegwerf-Container aus `u-boot:test`, `--network none`, Docker 29.9.2, danach entfernt, Filter leer): ohne `-e` Exit 1, `-e PROBE_VAR` mit Wert in der Umgebung gesetzt (Exit 0), `-e PROBE_VAR` ohne Wert nicht gesetzt (Exit 1). Den `u-boot`-Lauf der Gegenprobe habe ich nicht wiederholt (Nicht ausgeführt).
- **Handbuch L-2 (4.9, §7 Schritt 3):** 4.9 gleich dem Validator; §7 Schritt 3 kurz gefasst, dazu V-1.
- **READMEs:** „Podman ≥ 4.0“ = [`LH-FA-DIAG-002`](../../spec/lastenheft.md#lh-fa-diag-002--lokale-voraussetzungen-prüfen) (Zeile 683) = Handbuch Zeile 66 und Tabelle Zeile 2531; Ausnahme „nicht zugesagt für `devcontainer up`/`down`“ = Handbuch Zeile 69 (Entscheidung 3). *Optional* nennt Editor, Dev-Container-CLI oder Docker allein. *Documentation* und *Dokumentation* verlinken den Unterabschnitt in 4.9, die deutsche Liste zusätzlich 4.11 (die englische hatte ihn schon); Anker von `docs-check` geprüft. Status-Block ohne Datum und ohne „released“; EN und DE inhaltsgleich gelesen; die Kurzform sagt nichts, was der Lead nicht sagt. Suche nach `LH-`, `ADR-`, `SPEC-`, `ARC-`, `slice-`, `welle-`, `BEO-`, `spec/`, `docs/plan`, `roadmap`, `MR-` in den hinzugefügten Zeilen von READMEs und Handbuch: 0 Treffer. Ein grüner Docker-Lauf wird nur als CI-Aussage behauptet (Lead, beide Tabellenzeilen), nie „alle Docker-Tests“.
- **Versionsstrings und Handbuch-Kopf:** `git grep 0.13.0-dev` = sieben Stellen (`main.go` 1, `repo.mk` 3, `Dockerfile` 3); `0.12.0-dev` außerhalb von Historie und Slice-Text nur in der Messumgebung von [ADR-0026](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md) (Historie). Muster nach [`releasing.md`](../maintainer/releasing.md) §3 Nr. 2 und gegen den Tag `v0.12.0` (`main.go:43`, `repo.mk:55`, `Dockerfile:39` je `0.12.0-dev`). `publish.yml` übergibt `VERSION` aus dem Tag (Zeilen 79 bis 99) und pinnt Label und `--version`. Handbuch: `Handbuch-Version: 2.0` unverändert, `Software-Version: v0.13.0`, `Stand: 2026-10-21`, Beispiele Zeilen 82, 108, 135, Gültigkeitsbereich Zeile 2656.
- **L-10, ADR unverändert:** `gh api repos/aquasecurity/trivy-action/git/ref/tags/v0.36.0` → Tag-Objekt `a9c7b0f…` → Commit `ed142fd0673e97e23eac54620cfb913e5ce36c25` (= Pin in `ci.yml` Zeile 99). `action.yaml` an diesem Commit: Eingabe `version` („Trivy version to use“, Default `'v0.70.0'`, Zeilen 98 bis 101), weitergereicht als `version: ${{ inputs.version }}` an `aquasecurity/setup-trivy@3fb12ec…` (Zeile 131); `trivy-version` kommt in der Datei nicht vor. `git grep trivy-version` außerhalb von `done/`, `docs/reviews/`, Baseline: nur [ADR-0004](../plan/adr/0004-ci-system.md) Zeilen 93 und 96 und der Slice-Text. Kommentare in `ci.yml` (Indikativ, Zusage und Grenze), `repo.mk` und `quality.md` gleich benannt. `make doc-immutable RANGE="v0.12.0..HEAD"`: 0 Befunde. Entscheidung 6 (ADR unverändert, Abweichung benannt) ist eingehalten; keine Gate-Lockerung: der gescannte Stand bleibt `v0.70.0`.
- **Plan-Tabelle gegen `git diff --stat c26cfd43..HEAD`** ([`AGENTS.md`](../../AGENTS.md) §3.10, Stand nach der Fix-Runde): 16 Dateien, jede mit Zeile; umgekehrt jede Zeile mit Diff oder als „kein Diff“/„Closure“ geführt. `git diff --stat c26cfd43..HEAD -- spec docs/plan/adr internal docs/maintainer/cli-json-contract.md docs/messungen .github/workflows/publish.yml .github/workflows/integration.yml .github/workflows/tap-check.yml scripts packaging tools Makefile d-check.mk harness AGENTS.md .golangci.yml`: leer. Pipeline seit `v0.12.0` (`git diff --stat v0.12.0..HEAD -- .github scripts packaging tools/harness Makefile Dockerfile`): `ci.yml` (L-10), `Dockerfile` (Versionszeilen), zwei Skripte `tools/harness/ablage-pruefung*.sh` (Gate, nicht im Release-Pfad). Produktcode seit dem zuletzt in der CI gefahrenen Stand (`git diff --stat 0b8f268e..HEAD -- internal cmd Dockerfile repo.mk Makefile .github`): nur die vier Dateien der Versions-, Kommentar- und L-10-Zeilen. Abweichungen der Form: V-6.
- **§3.3, §3.7, §3.9:** `53a7cdd6` ist ein reiner Rename (0 Zeilen); der Vermerk im Welle-Kopf ist an den Absatz angehängt, kein Ursprungstext gelöscht (`git diff -U0 … -- docs/plan/planning/done`: eine Zeile, gleicher Anfang). Kommentare in `Dockerfile`, `repo.mk`, `main.go`, `ci.yml` im Indikativ, ohne Chronik; die ersetzten Werte stehen an keiner weiteren Stelle. `make doc-commits RANGE="3e7799b6..HEAD"`: acht Commits, 0 Befunde.
- **Zustandsmarker:** Roadmap Zeile 20 „(Tag-Push steht aus)“ ist der geltende Zustand bis M5; Zeile 99 trägt „Release `v0.13.0` … bestätigt am 2026-10-21“; `open/README.md` Zeilen 42 und 43 „(vergeben: `v0.13.0`)“; Colima-Slice §4 „vergeben als `v0.13.0`“, „eingetreten“ folgt nach dem Tag. „Die Nummer vergibt der Release-Cut“ steht in lebenden Dokumenten nur noch als Ursprungstext neben dem Vermerk im Welle-Kopf. `in-progress/` enthält nur die Master-Dokumente und diesen Slice (WIP 1).

### Sensors

Alle selbst gefahren auf `5ec45eee`, Arbeitsbaum sauber (vor dieser Datei), Docker Engine 29.9.2, Host `x86_64`, Ubuntu 24.04.5:

| Sensor | Ergebnis | Evidence | Laufzeit |
| --- | --- | --- | --- |
| `make gates` | pass | Exit 0; `lint` `0 issues`; `coverage-gate: OK — coverage 91.70% meets threshold 90%`; `d-check: 593 Datei(en) geprüft, 0 Befund(e)`; `baseline-verify: v6.18.1 OK — 54 Dateien`; `ablage-pruefung: OK`, Wächter-Test 26 Fälle, 0 Fehlschläge; Stempel `gates-passed.head` = `5ec45eee…` | 47 s |
| `make ci` | pass | Exit 0; Schwachstellen-Scan des Moduls: `No vulnerabilities found.`, „Your code is affected by 0 vulnerabilities“ (eine Schwachstelle in einem benötigten Modul, vom Code nicht aufgerufen); `image-scan`: `u-boot:latest (debian 12.15)` 0, `usr/local/bin/u-boot` 0 (HIGH/CRITICAL); Trivy meldet `0.76.0` verfügbar, gepinnt `0.70.0`; `[ci] gates + govulncheck + image-scan green` | 94 s |
| `make fullbuild` | pass | Exit 0; `UBOOT_VERSION=0.13.0-dev` im Build-Aufruf; `[fullbuild] ci + runtime image green` | 90 s |
| `make build VERSION=0.13.0` | pass | Exit 0; `docker run --rm u-boot:latest --version` = `u-boot version 0.13.0`; OCI-Label `org.opencontainers.image.version` = `0.13.0` (lokal, nicht gepusht); danach `make build`: `0.13.0-dev` | 7 s |
| `make doc-trace` | pass | `84 Anforderung(en), 0 Waise(n).` | — |
| `make doc-planning` | pass | `593 Datei(en) geprüft, 0 Befund(e)` | — |
| `make doc-immutable RANGE="v0.12.0..HEAD"` | pass | `219 Commit(s) — OK`; 0 Befunde (Accepted-ADRs unverändert, darunter [ADR-0004](../plan/adr/0004-ci-system.md)) | — |
| `make doc-commits RANGE="3e7799b6..HEAD"` | pass | `8 Commit(s) — OK`; 0 Befunde | — |
| M-T (Trockenlauf Tap-Pfad), selbst: `docker run -i --rm --network none -v /Development/u-boot:/src:ro u-boot:test bash < <Skript im Scratchpad>`, Fixture nach `sumsFull` | pass | (1) `homebrew-formula-fill.sh v0.13.0` Exit 0, vier Zeilen `releases/download/v0.13.0/`; (2a) `tap-formula-check.sh v0.13.0` Exit 0; (2b) mit `v0.12.0` Exit 1 („found 4 url lines, 0 of v0.12.0“); (3) `tap-nachzug.sh` mit `TAP_TOKEN=dummy-token`, `TAG=v0.13.0` und Stellvertretern: Exit 0, Protokoll `curl …/download/v0.13.0/u-boot.rb`, `git clone --depth 1` (Token maskiert), `git diff --quiet`, `git … add`, `git … commit -m u-boot 0.13.0`, `git push origin HEAD`, Meldung `updated to v0.13.0.`; (3b) `TAG=v0.12.0`: Exit 1, nur `curl`, **kein** `clone`. Rot gesehen in (2b) und (3b) | ca. 5 s |
| Lesend: `gh run list`, `gh api …/actions/jobs/<id>/logs`, `gh api …/check-runs/<id>/annotations`, `gh release view v0.12.0`, `git ls-remote --tags origin`, `git merge-base` | Befund | V-2 bis V-5 | — |
| `make docs-check`, `make gates` mit dieser Datei | pass | siehe Abschluss | — |

Traceability:

| ID | Beleg |
| --- | --- |
| [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [ADR-0028](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md), [ADR-0029](../plan/adr/0029-devcontainer-start-exit-11-und-socket-haltung.md) | Lead, *Added*, READMEs, Handbuch 4.9 gegen Validator-Bericht, Code (`startScriptFailed`, Templates) und Tests (`TestDevcontainerUp_Docker_S1_S4_SandboxProfileVolumes`, `…_S1_S9_StandardProfile`, `TestDevcontainerUp_RuntimeFailureClassification`); `internal/` ohne Diff, `make gates` grün |
| [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) | Handbuch 4.9/4.11 zu `GIT_TOKEN` in der `docker exec`-Shell gegen `devcontainer.json.tmpl` Zeilen 33 bis 34 und die selbst nachgefahrene `-e`-Semantik |
| [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [ADR-0026](../plan/adr/0026-upstream-dns-im-sandbox-devcontainer-nur-fuer-den-dnsmasq.md), [ADR-0027](../plan/adr/0027-egress-init-wiederholungslauf-atomar.md) | Lead (Verhaltensänderungen) gegen *Changed*-Einträge; Handbuch-Zeile `11` gegen `egress-init.sh.tmpl` und die e2e-Tests `…_EgressUpstreamDNS` und `…_EgressRerun` |
| [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) | CI-Satz (vier Skips) gegen `skipIfNestedUsernsRestricted` und die CI-Läufe 39540751579 und 39545000986 |
| [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) | Exit-Codes `10`/`11`/`12`/`14`/`2` im Lead gegen Handbuch-Tabelle und Klassen-Pin `cli/devcontainer_test.go` |
| [`LH-FA-DIAG-002`](../../spec/lastenheft.md#lh-fa-diag-002--lokale-voraussetzungen-prüfen) | „Podman ≥ 4.0“ beider READMEs gegen Lastenheft Zeile 683 und Handbuch Zeile 66 |
| [ADR-0004](../plan/adr/0004-ci-system.md) | L-10 gegen `action.yaml` der gepinnten Action; `doc-immutable` 0 Befunde; `make ci` grün |
| [ADR-0007](../plan/adr/0007-distributionswege-ghcr.md), [ADR-0017](../plan/adr/0017-linux-pakete-nfpm.md) | `make build VERSION=0.13.0` (Label und `--version`); `publish.yml`, `packaging/`, `scripts/` ohne Diff seit `v0.12.0`; stabiler Lauf `v0.12.0` (39224848402) in allen vier Jobs `success`, zwölf Assets |
| [ADR-0016](../plan/adr/0016-homebrew-distribution-per-tap.md) | M-T selbst nachgefahren; Job `tap` im Lauf 39224848402 `success` mit `tap-nachzug: pt9912/homebrew-u-boot updated to v0.12.0.` |

Carveouts:
- Neu: none
- Gelöst: none
- Unverändert: alle (`carveouts.md` ohne Diff; kein Gate, keine Schwelle gelockert; L-10 schärft den Pin, lockert nichts)

Nicht ausgeführt:
- `make test-docker` lokal — Auftrag (bekanntes Zeitlimit des Pakets `internal/e2e` auf diesem Host, Slice §3); Nachweis ist `integration-docker` auf dem Push-Tip (M2). Zwei CI-Läufe auf produktgleichem Code liegen vor (V-3).
- `make verify-depguard` — keine Importregel und kein Go-Import berührt (Produktcode im Diff: nur `var version`).
- `u-boot devcontainer up` mit Gegenprobe `docker exec -e` und die Validator-Messungen zu L-1/L-2 — nicht wiederholt; nur die Docker-Semantik von `-e` nachgefahren. Rohdaten liegen gewollt nicht im Repo.
- `git push`, `git tag`, `gh workflow run tap-check`, `gh release` — Auftragsgrenze; M1 bis M8 unten.
- Erreichbarkeit der Vergleichslinks im Netz — nur Form geprüft (`v0.13.0` existiert erst nach M5).
- Colima, Docker Desktop, `arm64` — nicht Gegenstand, nicht zugesagt.

Commit / Artefakt:
- Geprüft auf `5ec45eee`; lokales Image `u-boot:latest` am Ende wieder `0.13.0-dev`; kein Release-Artefakt.

## Review-Findings einzeln

| Finding | Schwere | Stand | Beleg des Verifiers |
| --- | --- | --- | --- |
| F-1 Abbruch im CHANGELOG verallgemeinert | LOW | **behoben** | Lead Zeile 22 und *Added* Zeile 38 nennen Strg-C nur „during a `postCreateCommand`“, je Fall den nicht erneut gelaufenen Befehl und die nicht gemessenen Fälle; gleich Handbuch 4.9 und Validator (oben). |
| F-2 Marker `open/README.md` Zeile 42 | LOW | **behoben** | „(vergeben: `v0.13.0`; …)“ wie Zeile 43; Zuordnung selbst bestätigt: `git merge-base --is-ancestor 08c0d10a v0.12.0` endet mit Exit 1 (kein Vorfahre des Tags). |
| F-3 lokales `test-docker` rot, CI genügt | INFO | **übernommen in M2** | Prüfliste M2 liest Job-Konklusion **und** Zählung; Vorab-Datenpunkt V-3. |
| F-4 Kopplungs-Kommentar Zeitlimit (`repo.mk`) | INFO | **übergeben (Planner, Closure)** | Bestand außerhalb des Diffs; die CI-Jobdauer schwankt nach oben (V-3). |
| F-5 Exit-Code-Tabelle Zeile `11` | INFO | **behoben** | gegen Templates, `startScriptFailed` und Tests gelesen (oben). |
| F-6 [ADR-0004](../plan/adr/0004-ci-system.md) nennt `trivy-version` | INFO | **entschieden (Projektinhaber, Entscheidung 6)** | `doc-immutable` 0 Befunde; Abweichung in §6 benannt. |
| F-7 erneutes `up` nach Exit 11 | INFO | **offene Entscheidung (a)** | Handbuch sagt „nicht gemessen“ (4.9 und §7). |
| F-8 Marker in `done/` | INFO | **übergeben (Planner, Closure)** | ein Hop über den Vermerk im Welle-Kopf, geprüft. |

## Befunde und Diskrepanzen des Verifiers

| ID | Klasse | Befund | Wirkung |
|---|---|---|---|
| V-1 | INFO | Handbuch §7 Schritt 3 fasst L-2 kürzer als 4.9: „führte in den beobachteten Fällen die Lebenszyklus-Befehle nicht erneut aus (nach Exit-Code 12 und nach einem Abbruch, Abschnitt 4.9)“. Gemessen ist je ein Befehl (das gescheiterte `postStartCommand`, das abgebrochene `postCreateCommand`) und nur Strg-C **während eines `postCreateCommand`**. Der Satz ist durch „in den beobachteten Fällen“ und den Verweis begrenzt, also keine Zusage über das Gemessene hinaus; die Beleg-Zeile zum Strg-C-Fall nennt als Ort aber nur 4.9 und CHANGELOG, nicht §7. | Keine Blockade. Wer präzisieren will („nach Strg-C während eines `postCreateCommand`“), tut es **vor dem Tag**: `--help` verlinkt das Handbuch am Tag der Version, eine spätere Änderung auf `main` erreicht diesen Link nicht. Entscheidung Implementer/Projektinhaber. |
| V-2 | INFO | Slice §6 (*Runner-Wechsel*) führt das Betriebssystem des letzten stabilen `publish`-Laufs als „nicht gelesen (`gh run view --job 117968525028 --log` liefert keine Zeilen)“. Über `gh api repos/pt9912/u-boot/actions/jobs/117968525028/logs` ist es lesbar: Ubuntu 24.04.5, Runner Image `ubuntu-24.04` Version 20261012.334.1. `gh run view --job … --log` lieferte auch für den aktuellen Job 119041060670 keine Zeile; der API-Weg liefert das volle Log (9630 Zeilen). | Eingabe für den Ausgang des Risikos (Closure). Für M2 und M6 den API-Weg nutzen. |
| V-3 | INFO | Neuerer Vorab-Datenpunkt als der im Slice zitierte: Auf `origin/main` = `3e7799b6` (Elter des Cuts, produktgleich) liefen `ci` 39545000945 und `integration-docker` 39545000986, beide `success`. `integration-docker`: 2835 `--- PASS`, 4 `--- SKIP` (`…_NestedPodman_Build`, `…_NestedPodman_BuildEgress`, `…_NestedPodman`, `…_StatusFiles_VFS`), 0 `--- FAIL`, alle Pakete `ok`; Ubuntu 24.04.5 (Image `ubuntu-24.04` 20261018.339.1); Job 13 min 32 s (`internal/e2e` 747,7 s), gegen 11 min 28 s (631,998 s) auf `0b8f268e` (Lauf 39540751579, Zählung von mir bestätigt: 2835/4/0). `image-scan` auf `3e7799b6` trägt erwartungsgemäß noch zweimal `Unexpected input(s) 'trivy-version'` und die Notice zur Migration von `ubuntu-latest` auf Ubuntu 26. | Der CI-Satz des Release ist durch zwei Läufe auf produktgleichem Code getragen; bindend bleibt M2. Die Jobdauer schwankt um zwei Minuten bei 30 Minuten Limit — Eingabe zu F-4 und Restrisiko (b) des Reviews. |
| V-4 | INFO | Der L-10-Beleg „das Log nennt `version: v0.70.0`“ trennt nicht: Das Log von `image-scan` auf `3e7799b6` zeigt im `with:`-Echo des Schritts `version: v0.70.0` (Default der Action) **neben** `trivy-version: v0.70.0`, obwohl die Eingabe dort nicht angenommen wurde. Trennend sind (1) keine Annotation `Unexpected input(s)` (auf `3e7799b6`: zwei, am Schrittbeginn und im Post-Schritt) und (2) keine Zeile `trivy-version:` im `with:`-Echo. Der Slice benennt die Grenze (*Befund L-10*), die DoD-Zeile nennt das Log trotzdem als Beleg. | M2-Prüfung unten so geschärft; bei der Closure den Beleg in dieser Form festschreiben. |
| V-5 | INFO | Für M6 („keine Warnung zum fehlenden Secret“): Das Log des Jobs `tap` im Lauf `v0.12.0` enthält den Text `HOMEBREW_TAP_GITHUB_TOKEN is not configured` im Echo des Schritt-Skripts (Zeile 162), obwohl der Nachzug lief (Zeile 176: `tap-nachzug: pt9912/homebrew-u-boot updated to v0.12.0.`). Eine Textsuche im Log meldet den Fall also falsch-positiv. | M6 prüft die Annotationen des Jobs und die Erfolgszeile, nicht eine Textsuche. |
| V-6 | INFO | Plan-Tabelle, Form: (a) Die Zeile der Reports führt den Verification-Report nach dem Muster `<Datum>-verification-<Slice-Kennung>.md`; diese Datei folgt auf Auftrag dem Muster `<Datum>-<Slice-Kennung>-verify.md` (Form der Reports seit 2026-10-20). (b) Die Zeile `.github/workflows/publish.yml, …` begründet „kein Diff“ mit „`git diff --stat v0.12.0..HEAD -- .github scripts packaging` ist leer“ und „`git grep trivy-version` trifft … `ci.yml:117`, `repo.mk:60`, `quality.md:253`“ — Stand `3e7799b6` laut Kopf von §3; seit `43760c4d` trifft das Diff `ci.yml`, und `trivy-version` steht nur noch in der ADR. Der Abschnitt *Umsetzung* nennt den Ist-Stand richtig. | Abgleich nach der letzten Änderung ([`AGENTS.md`](../../AGENTS.md) §3.10) bei der Closure: Dateiname nachziehen, die beiden Begründungen auf den Ist-Stand setzen. |
| V-7 | INFO | §7 *Vorgemerkt*: „Tag-Datum … gegen die fünf Datumsstellen aus M0“; M0 nennt vier (CHANGELOG, Tabellenzeile `README.md`, Tabellenzeile `README.de.md`, Handbuch-`Stand`). | Planner nennt die fünfte oder schreibt vier. |

Keine DoD-Verletzung vor dem Push gefunden. Keine Behauptung ohne Beleg in den geprüften Nutzer-Aussagen.

## Auflagen und Prüfliste

### Vor M1 (Push von `main`)

1. **Diese Datei committen** (nur sie; Slice-Kennung und eine `LH-*`- oder `ADR-*`-Kennung in der Nachricht), damit der geprüfte Stand im Push liegt; danach `make docs-check` und `make gates` auf dem neuen `HEAD` (Stempel). `git log --oneline origin/main..main` zeigt dann neun Commits; `git status` leer.
2. **M0 am Push-Tag:** `date +%F` = `2026-10-21`. Weicht der Tag des späteren Tags ab, ist es ein weiterer Cut-Commit vor dem Tag (vier Datumsstellen), dann Gates, Push und M2 erneut.
3. Optional, nur vor dem Tag: V-1 entscheiden.
4. Push nur durch den Hauptagenten nach ausdrücklicher Freigabe des Projektinhabers im Gespräch. Die Ausgabe des Push lesen (beim Cut `v0.12.0` meldete GitHub „Bypassed rule violations for refs/heads/main“, Register).

### M2 — CI auf dem Push-Tip (vor jedem Tag)

- `gh run list --limit 4 --json databaseId,workflowName,headSha,status,conclusion`: `ci` und `integration-docker` mit `headSha` = vollständiger Hash von `git rev-parse HEAD` nach M1, beide `completed`/`success`. Keine weiteren Commits zwischen M2 und M5 ohne neuen Lauf.
- **`ci`, Job `image-scan (trivy HIGH+CRITICAL)` (L-10):** `gh api repos/pt9912/u-boot/check-runs/<job-id>/annotations` — **kein** Eintrag `Unexpected input(s)` (zulässig ist die Notice zu `ubuntu-latest`); `gh api repos/pt9912/u-boot/actions/jobs/<job-id>/logs`: im `with:`-Echo des Schritts „Trivy scan (image)“ steht `version: v0.70.0` und **keine** Zeile `trivy-version:`; Schritt `setup-trivy` mit `version: v0.70.0`, Cache-Schlüssel `trivy-binary-v0.70.0-Linux-X64`; Tabelle 0 HIGH/CRITICAL (V-4). Ein roter `image-scan` (neuer Advisory) ist Rückführung nach §4 des Slice.
- **`integration-docker`:** Job-Konklusion `success` (nicht nur die Zählung, F-3); im Log über die API: `--- FAIL` 0, kein `panic: test timed out`, alle Pakete `ok`; `--- SKIP` **genau 4** und genau `TestE2E_SandboxDevcontainer_NestedPodman`, `…_NestedPodman_Build`, `…_NestedPodman_BuildEgress`, `…_StatusFiles_VFS` mit Grund `kernel.apparmor_restrict_unprivileged_userns=1`; Referenz `--- PASS` 2835. Gruppe „Operating System“ lesen (Referenz Ubuntu 24.04.5; Image `ubuntu-24.04` 20261018.339.1). Jobdauer gegen 30 Minuten (Referenz 11 min 28 s und 13 min 32 s; über 25 Minuten als Befund für F-4 festhalten).
- **Abweichung:** andere SKIP-Zahl oder andere Namen (etwa nach dem Wechsel auf Ubuntu 26) → Release-Aussage vor dem Tag nachziehen: CHANGELOG-Lead Zeile 27, Status-Block beider READMEs, Tabellenzeile beider READMEs; dann Gates, Push, M2 erneut. Ein `--- FAIL` oder Zeitlimit → Rückführung nach §4.

### M3 — entfällt

Entscheidung 2 des Projektinhabers (direkt stabil). Der stabile Lauf ist damit die erste Probe auf dem aktuellen Runner-Image; das in der Closure benennen.

### M4 — `tap-check` unmittelbar vor M5

- `gh workflow run tap-check` (Maintainer), dann `gh run list --workflow tap-check --limit 1`: `success`, angelegt nach M2. Letzter Lauf: 39224782459 vom 2026-10-19 (`success`) — zu alt, der PAT kann abgelaufen sein ([`releasing.md`](../maintainer/releasing.md) §4).
- M0 erneut (`date +%F`).

### M5 — Tag (nur Hauptagent, nur nach ausdrücklicher Freigabe im Gespräch)

- `git switch main && git pull --ff-only`; `git rev-parse HEAD` = `headSha` aus M2; `git tag v0.13.0` (leichtgewichtig wie `v0.12.0`, `git cat-file -t v0.12.0` = `commit`); `git push origin v0.13.0`. Danach wird der Tag nicht verschoben.

### M6 — Kontrolle nach [`releasing.md`](../maintainer/releasing.md) §6

- `publish`-Lauf: Jobs `publish (ghcr.io image push)`, `tap (homebrew formula)`, `package smoke (rpm, fedora)`, `package smoke (deb, ubuntu)` je `success`; im Log `Resolved version=0.13.0, is_stable=true`; Gruppe „Operating System“ der Jobs auf `ubuntu-latest` (Zeilen 33 und 316 von `publish.yml`) lesen — Ausgang des Risikos *Runner-Wechsel* (Referenz `v0.12.0`: Ubuntu 24.04.5, V-2). Lauf-ID in §7.
- `gh release view v0.13.0 --json assets,isPrerelease,isDraft`: zwölf Assets (Muster `v0.12.0`: `SHA256SUMS`, `u-boot-0.13.0-1.aarch64.rpm`, `u-boot-0.13.0-1.x86_64.rpm`, sechs Binaries `u-boot-{darwin,linux,windows}-{amd64,arm64}` mit `.exe` für Windows, `u-boot.rb`, `u-boot_0.13.0_amd64.deb`, `u-boot_0.13.0_arm64.deb`), `isPrerelease` und `isDraft` `false`. Release-Notes: Auto-Text mit Verweis auf `CHANGELOG.md` §[0.13.0].
- `docker run --rm ghcr.io/pt9912/u-boot:0.13.0 --version` = `u-boot version 0.13.0`; Digest von `:0.13.0` = `:latest`.
- Job `tap` **wirklich gelaufen**: Log-Zeile `tap-nachzug: pt9912/homebrew-u-boot updated to v0.13.0.` (oder `tap already at v0.13.0; nothing to do.`) und in den Annotationen des Jobs keine Warnung `HOMEBREW_TAP_GITHUB_TOKEN is not configured` — nicht per Textsuche im Log (V-5).

### M7 — Mac (Projektinhaber)

- `brew tap pt9912/u-boot`, `brew trust pt9912/u-boot`, `brew install u-boot` oder `brew upgrade u-boot`, `u-boot --version` = `0.13.0`; Smoke-Workflow im Tap-Repo nach dem Formel-Commit grün.

### M8 — Nacharbeit und Closure (Planner, Implementer)

- Roadmap: Zeile `v0.13.0` mit Datum und Tag-Hash nach Zeile 99, Ruhe-Marker zurück, `**Letzte Änderung:**`; `git log -1 --format=%cs v0.13.0` gegen die vier Datumsstellen (V-7).
- Colima-Folge-Slice: Vermerk „eingetreten: Tag `v0.13.0`“.
- `state.md` von [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../plan/planning/observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md) mit SKIP-Zahl und Betriebssystem des Push-Tips und des `publish`-Laufs.
- §6: Ausgang je Risiko (V-2, V-3 als Eingaben); F-4, F-8, V-1, V-6, V-7 erledigen oder mit Grund belassen; Plan-Tabelle gegen `git diff --stat c26cfd43..HEAD` nach der letzten Änderung (V-6); Freigabe von Push und Tag mit Zeitpunkt und Ausführendem.
- Erst Inhalt (Closure-Notiz, DoD-Haken), dann reiner `git mv` nach `done/`, dann Link-Nachzug — im selben Push.

## Minimal Handoff

```text
Role: Verifier
Input context: Slice §1–§8, Review-Report vom 2026-10-21, releasing.md, git diff 3e7799b6..HEAD, Code, Templates und Tests zu den neuen Nutzer-Sätzen, CI-Läufe 39540751579, 39545000986, 39224848402 (lesend)
Changed artefacts: dieser Verification-Report (neu, nicht committet)
Evidence: make gates, ci, fullbuild, build VERSION=0.13.0, doc-trace, doc-planning, doc-immutable, doc-commits grün; M-T nachgefahren; docker exec -e nachgefahren
Open risks: M2 (SKIP-Zahl nach Runner-Wechsel, Jobdauer), M4 (PAT), Tag-Datum, erster stabiler Lauf ohne Vorab-Tag
Next role: Hauptagent (M1 nach Freigabe), dann M2 bis M8; Planner (Closure)
```

## Abschluss

Sensor-Lauf nach dem Anlegen dieser Datei (Arbeitsbaum = `HEAD` `5ec45eee` plus diese Datei, ungestagt): `make docs-check` Exit 0, `d-check: 594 Datei(en) geprüft, 0 Befund(e)`; `make gates` Exit 0. Ein erster `docs-check`-Lauf fand zwei unverlinkte Kennungen in dieser Datei (`id-unlinked`); sie sind verlinkt bzw. umformuliert.
