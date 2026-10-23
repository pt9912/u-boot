# Verification-Report: slice-v1-release-cut-v0.15.0 — 2026-10-23

**Rolle:** Verifier ([`harness/verification.md`](../../harness/verification.md)), frischer Kontext, kein Self-Review, kein Reparieren. Frage: Wurde der Release-Cut richtig gebaut, gemessen an Plan, DoD und ausgeführten Sensoren, **bis zur Push-Reife (M1)**? Die Punkte nach M1 (CI auf dem Push-Tip, `tap-check`, Tag, `publish`, GHCR, Tap-Job, Mac) sind nicht Gegenstand dieses Berichts und stehen aus.

**Gegenstand:** [`slice-v1-release-cut-v0.15.0`](../plan/planning/done/slice-v1-release-cut-v0.15.0.md), Range `07353780..6f79fc73` (acht Commits vor `origin/main`, nichts gepusht, kein Tag `v0.15*`), Bezug [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) und [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer). Eingang: Plan §2/§3, Review-Report [`2026-10-23-slice-v1-release-cut-v0.15.0`](2026-10-23-slice-v1-release-cut-v0.15.0.md) (Eingang, nicht Beleg), [`releasing.md`](../maintainer/releasing.md), Muster [`2026-10-23-slice-v1-release-cut-v0.14.0-verify`](2026-10-23-slice-v1-release-cut-v0.14.0-verify.md).

## Ergebnis

**Schließbar bis M1: GO für den Push (M1). Der Slice als Ganzes ist noch nicht schließbar** — die Closure-DoD hängt an M2, M4 bis M7 und M8, die nach dem Push liegen. Kein Blocker, kein Befund, der Produktcode oder eine Änderung am Cut verlangt. Zwei Hinweise für M2 und die Closure (unten). Alle Sensoren, deren Ausgabe ich nicht aus dem Bericht des Implementers übernehmen wollte, habe ich selbst gefahren (Host `linux/amd64`, Docker Engine 29.9.2).

**Push-Tip:** Dieser Report ist ein weiterer Commit auf `6f79fc73`. Push-Tip, M2 und Tag beziehen sich auf den Commit, der diesen Report trägt, nicht auf `6f79fc73`.

## Sensoren — selbst gefahren

| Sensor | Ergebnis | Beleg |
| --- | --- | --- |
| `make gates` (auf `6f79fc73`) | **Exit 0** (`GATES_EXIT=0`) | `lint` 0 issues; `docs-check` 781 Dateien, 0 Befunde; `baseline-verify` v6.18.1 OK, 54 Dateien; `ablage-pruefung` 26 Fälle, 0 Fehlschläge; `coverage-gate` lief im Gate (Exit 0 des Ganzen) |
| `make ci` | **Exit 0** | `govulncheck` ohne Befund; `image-scan`: `u-boot:latest (debian 12.15)` 0, `usr/local/bin/u-boot` 0 |
| `make build VERSION=0.15.0` | **Exit 0** | `docker run --rm u-boot:latest --version` = `u-boot version 0.15.0`, OCI-Label `org.opencontainers.image.version` = `0.15.0` (danach durch `make ci` wieder Standardbau `0.15.0-dev`) |
| `make doc-trace` | Exit 0 | 84 Anforderungen, 0 Waisen |
| `make doc-planning` | Exit 0 | 781 Dateien, 0 Befunde |
| `make doc-commits RANGE=07353780..HEAD` | Exit 0 | 8 Commits, 0 Befunde |
| `make doc-immutable RANGE=v0.14.0..HEAD` | Exit 0 | 93 Commits, 0 Befunde (keine Accepted-ADR umgeschrieben) |
| `make verify-depguard` | nicht gelaufen | Grund: keine Importregel, kein Go-Import berührt; einziger Code-Diff ist `var version` in `cmd/uboot/main.go` (Diff gelesen) |
| `make fullbuild` | nicht wiederholt | Grund: `ci` (gefahren, Exit 0) plus Runtime-Image-Bau (`make build`, gefahren, Exit 0) decken seinen Inhalt; Implementer-Beleg Exit 0 |
| `make test-docker` | **nicht gefahren** | Grund: Laufzeit über dem Paket-Zeitlimit von 30 Minuten (Befund des Implementers: Exit 2 am Zeitlimit, 1 `--- FAIL` `NestedDocker_Restart` durch Netzstörung, Einzellauf PASS, die vier Tests der Welle PASS); der lokale Lauf ersetzt ohnehin die CI nicht. Das Ergebnis des Implementers wurde von mir **nicht reproduziert**, nur auf Darstellung geprüft: es ist ehrlich ausgewiesen (Exit 2, FAIL und Zeitlimit benannt, kein „grün“ behauptet) |

## Tap-Trockenlauf M-T — selbst wiederholt

Image `u-boot:test`, `--network none`, Arbeitsbaum als `:ro`-Mount, eigenes Skript im Scratchpad (nicht im Repo), Fixture nach `sumsFull` in `cmd/uboot/homebrew_scripts_test.go`, Stellvertreter für `curl` und `git` im `PATH`. Kein Netz, kein Schreiben ins Tap, kein Push.

| Schritt | Exit | Beobachtung |
| --- | --- | --- |
| `homebrew-formula-fill.sh v0.15.0` | 0 | vier `url`-Zeilen mit `v0.15.0` |
| `tap-formula-check.sh v0.15.0` | **0** | |
| `tap-formula-check.sh v0.14.0` (Gegenprobe) | **1** | „found 4 url lines, 0 of v0.14.0“ |
| `tap-nachzug.sh`, `TAG=v0.15.0`, Dummy-Token | 0 | Protokoll: `curl …/v0.15.0/u-boot.rb`, `git clone`, `git diff --quiet`, `git add`, `git commit -m u-boot 0.15.0`, `git push origin HEAD`; „updated to v0.15.0“ |
| `tap-nachzug.sh`, `TAG=v0.14.0` (Gegenprobe) | 1 | Protokoll nur der `curl`-Aufruf, **kein** `clone`, kein `push` |

Deckt sich mit dem Protokoll des Implementers (§3). **Beweist nicht:** `is_stable`, echter `curl` gegen das Release-Asset, echter `git push`, Token. Das ist die Aufgabe von `tap-check` (M4) und des Jobs `tap` (M6).

## DoD Punkt für Punkt (Plan §2)

| DoD-Punkt | Stand | Beleg / Bemerkung |
| --- | --- | --- |
| Release-Text und Handbuch | **erfüllt** | `## [0.15.0] - 2026-10-23` mit Lead (A-1, M-2, Nicht-Zusagen), Vergleichslinks `[Unreleased]…v0.15.0...HEAD`, `[0.15.0]…v0.14.0...v0.15.0`; `0.15.0-dev` an `cmd/uboot/main.go`, `repo.mk` (3 Stellen), `Dockerfile` (3 Stellen); beide READMEs (Status, `VERSION=`, Tabellenzeile) und Handbuch (Kopf, drei Beispiele, Gültigkeitsbereich) auf `v0.15.0`. `grep` nach `0.14.0` außerhalb von CHANGELOG/Planung/Reviews/Messprotokoll/Baseline: nur Tabellenzeilen `v0.14.0`, Lastenheft-Historie (Zeilen 1590 und 1593, anderer Gegenstand) und `releasing.md:57` (Zeiger „zuletzt“, bis zur Closure richtig). Nichts vergessen. README/README.de/Handbuch: im `+`-Diff keine `LH-/ADR-/SPEC-/ARC-/slice-/welle-`-Kennung, kein Verweis auf Planung (Wortsuche rc=1 und Diff gelesen) |
| Tag und Kontrolle | **teils** | Vor Tag erfüllt: `make gates`, `make ci` (selbst gefahren), Trockenlauf M-T (selbst wiederholt), Review-Report liegt vor; `fullbuild` über `ci` + `build` gedeckt. **Stehen aus:** `make test-docker` ohne Zeitlimit-Abbruch (lokal Exit 2, siehe oben), M2 (CI auf dem Push-Tip), M4 `tap-check`, M5 Tag, M6 `publish`/GHCR/Assets/Notes/Job `tap`, M7 Mac |
| Befunde des Validators haben einen Ausgang | **erfüllt (Text)** | A-1 im Lead (Zusage-Menge, „kein authentifizierter Aufruf gemessen“, Nicht-Zusagen), M-2 im Lead und im Eintrag *Added*; A-2/A-3 im Handbuch-Slice (`done/`), L-1 bis L-3, I-1, I-2 nicht entschieden (Plan §1) |
| `make gates` grün | **erfüllt** für den Stand `6f79fc73`; nach der Closure erneut nötig |
| Review durchgeführt | **erfüllt** | Report vorhanden, 1 MEDIUM (Reichweite von `version:`), 1 LOW, INFO; MEDIUM in `6f79fc73` behoben (unten) |
| Doku-Update | **teils** | CHANGELOG, READMEs, Handbuch, Roadmap-Zeiger (Slice), Welle-Datei §4 und `open/README` auf Ist-Stand; **stehen aus (Closure/M8):** `releasing.md` Zeiger „zuletzt“, Roadmap-Zeilen der Welle und `v0.15.0` mit Tag-Hash |
| Closure-Notiz, Register, Risiken, Paarungen | **steht aus** | §7 trägt bewusst nur die Start-Entscheidungen; §6 trägt „Ausgang: bei der Closure“ — erst nach M6 |

## Vergleich mit dem Cut `v0.14.0` (Muster)

Gleiches Muster: Tabellenzeile pro README mit Zusage-Menge und Grenzen, Lead im Stil „Fifteenth release.“ mit fetter Kernaussage, `releasing.md` §3 Schritte 1 bis 4 (CHANGELOG, drei Versionsstellen samt Kommentaren, READMEs, Roadmap) vollständig; Schritt 5 (Slice) läuft; Move (`9f15adf5`) ist reiner Move, Verweis-Nachzug (`138cdc41`) eigener Commit, danach Inhalt. **Zwischenzustand nicht Spitze:** Move und Nachzug liegen mehrere Commits unter dem Inhalts-Tip; mit diesem Report ist die Spitze ein Inhaltscommit. Abweichung vom Muster, begründet: kein CI-Satz im Lead (M2 steht aus, Abweichung 1 des Implementers) — richtig und in `[0.15.0]` nicht verletzt (kein „CI“-Satz im Lead; die unveränderten Status-Block-Sätze der READMEs zur CI sind Bestand, siehe Hinweis 1).

## Lead des CHANGELOG `[0.15.0]` — Aussage für Aussage

| Aussage | Beleg |
| --- | --- |
| `gh`/`glab` als Katalog-Features, beide Profile, Token nur als Variable des Hosts, Sandbox-Profil `${localEnv:…}`, Standardprofil ohne Referenz, Weg `docker exec -e` (ein Lauf, Platzhalter) | Plan Beleg-Tabelle Zeilen 1 und 5; Validator-Bericht Punkt 5; Handbuch |
| Zusage nur Linux `amd64`, lokale Engine, Docker 29.9.2; Nicht-Zusagen (macOS/Colima auch bei Homebrew, Docker Desktop, rootless, Podman als Host-Engine, entfernte Engine, `arm64`, Windows) | Validator-Bericht *Messumgebung* und *Was ohne gültiges Token ungeprüft blieb*; gleiche Liste wie im Lead von `v0.14.0` |
| „Kein authentifizierter Aufruf gemessen“ (Platzhalter: 401 / `Bad credentials`) | Validator-Bericht |
| Kalt-Bau je Profil ohne Layer-Cache, Basis-Image lokal | Slice §3 Abweichung 2 (eigene Messung, ein Lauf je Profil) — im Text mit Grenze („base image already local“) |
| Nicht zugesagt: `glab duo`, `glab orbit`, selbst gehostete Instanzen, Hosts nach Anmeldung, Telemetrie-/Update-Hosts | Bestandstext, Handbuch `:2159`–`:2167` |
| `version:` ohne `source:` endet bei `gh`, `glab`, `ollama` mit Exit `10`; `java` u. a. nehmen es an | **Fix-Runde `6f79fc73`, selbst gelesen:** `validateDevcontainerFeatures` (`devcontainer_features.go`, Bedingung `entry.Source == "" && entry.Version != "" && featureCatalogue()[k].dockerfileBlock`); `dockerfileBlock: true` in genau drei Katalogeinträgen (`ollama`, `gh`, `glab`, Zeilen 281, 285, 289). Tests gelesen: `TestGenerateDevcontainer_ForgeCLI_VersionOverrideWithoutSourceRejected` (gh/glab × Standard-/Sandbox-Profil, `ErrGenerateManualConflict`, nichts geschrieben), `TestGenerateDevcontainer_Ollama_VersionOverrideWithoutSourceRejected`, Gegenprobe `TestGenerateDevcontainer_FeaturesJavaVersionAndSdkman_RenderTwoKeys` (`java` mit `version: "21"` ohne `source` erzeugt `java:21`, ohne Fehler). Alle laufen in `make gates` (Exit 0) |
| `config set …gh.version` (und `glab`) prüft es nicht, meldet Erfolg, `generate`, `doctor` und jedes weitere `config set` lehnen danach ab; Abhilfe: Eintrag von Hand entfernen | gelesen: `TestConfigSet_ForgeVersionWithoutSource_SetItselfSucceeds` (nur `gh`), `…_RejectsEverySet` (gh und glab, vier Sets, Datei byte-identisch), `doctor_features_test.go:368`; Code-Pfad im Slice §3 Abweichung 3. Handbuch `:2060` ff. gleiche Aussage |
| Konfigurationen ohne `gh`/`glab` rendern byte-identisch; Lastenheft 0.22.0 | Golden Cases des Code-Slice; `spec/lastenheft.md` Zeile 9 |
| **Keine CI-Behauptung vor M2** | im Lead und im *Added*-Eintrag nicht vorhanden (gelesen) |

Die Reichweite von `version:` im Lead, im Eintrag *Added* und in der Tabellenzeile beider READMEs ist nach der Fix-Runde nicht mehr breiter als der Code (MEDIUM des Reviews behoben).

## Weitere Prüfungen

- **Keine Tag-/Push-Aktion:** `git tag -l 'v0.15*'` leer; `origin/main` = `07353780`; `main` acht Commits voraus; Arbeitsbaum vor diesem Report sauber.
- **Datum (M0):** `date +%F` = `2026-10-23` = CHANGELOG `[0.15.0]` = beide Tabellenzeilen = Handbuch-`Stand`. Weicht es vor M1 oder M5 ab, ist neu zu schneiden (`releasing.md` §5).
- **Plan-Tabelle gegen `git diff --stat`:** zwölf Dateien (inklusive Review-Report), jede mit Plan-Eintrag; Workflows, `scripts/`, `packaging/`, `Makefile`, `go.mod`, `internal/`, `spec/`, ADRs ohne Diff.
- **Zustandsmarker:** „Push und Tag stehen aus“ in `roadmap.md` und `open/README.md` sind heute wahr und müssen von der Closure auf den Ist-Stand gesetzt werden (LOW des Reviews, bewusst offen).

## Hinweise (keine Blocker)

1. **M2 muss den unveränderten Status-Block-Satz prüfen.** Beide READMEs tragen unter `v0.15.0` weiterhin „The Docker integration tests run in CI; the four tests of nested Podman are skipped there.“ (Bestandssatz). Er gilt für diesen Release erst, wenn `integration-docker` auf dem Push-Tip grün ist, genau vier `--- SKIP` zeigt und `TestE2E_ForgeCLI_BuildsInBothProfiles`, `…_PinMismatchBreaksBuild`, `…_EgressHosts`, `…_TokenReferencesCarryNoValue` sichtbar `--- PASS` sind (Job-Log über REST-API, Check-Run-Annotationen, nicht per Textsuche; Dauer gegen `timeout-minutes: 30`). Weicht der Lauf ab: Aussage vor dem Tag nachziehen und neu schneiden. Der lokale Docker-Lauf lief in das Paket-Zeitlimit; ob die CI (anderer Runner) darin bleibt, ist **ungemessen**.
2. **`glab` bei `config set` ist nur gelesen, nicht im Einzeltest gepinnt:** `…_SetItselfSucceeds` prüft `gh`; `glab` läuft über dieselbe, schlüsselunabhängige Stelle (`revalidateFeatureEntry`) und ist durch `…_RejectsEverySet` (beide Schlüssel) nur für das Danach gepinnt. Die Aussage „(and `glab`)“ ist damit am Code gelesen, aber ohne eigenen Test für das Schreiben. Gering; kein Handlungsbedarf vor dem Cut.

## Ungemessen — ehrlich ausgewiesen

- **M2:** CI (`ci`, `integration-docker`) auf dem Push-Tip — nicht möglich vor M1.
- **M4:** `tap-check` (echter Klon und `git push --dry-run` mit dem Secret) — Sache des Hauptagenten unmittelbar vor M5.
- **M6:** `publish`-Lauf, GHCR-Digest `:0.15.0` = `:latest`, zwölf Release-Assets, Job `tap` — nach dem Tag.
- **M7:** Homebrew auf einem Mac (`brew install u-boot`, `--version` = `0.15.0`) — Projektinhaber; ersetzt keine Zusage für `gh`/`glab` auf macOS (nicht zugesagt).
- **`make test-docker` voll:** nicht reproduziert (Zeitlimit).
- **Authentifizierter Aufruf von `gh`/`glab`:** nie gemessen (so im Lead gesagt).
- **Echte Tap-Schritte** (`is_stable`, echter `curl`/`push`, Token): vom Trockenlauf nicht abgedeckt.

## Was dem Push und der Schließung vorausgehen muss

1. **Vor M1:** Freigabe des Projektinhabers im Gespräch (je Vorgang); `date +%F` erneut (Soll `2026-10-23`); `make docs-check` lokal auf dem Tip mit diesem Report (Arbeitsbaum sauber); alle Commits in **einem** Push (Move `9f15adf5` und Nachzug `138cdc41` nicht als Spitze).
2. **M2** auf dem Push-Tip (Hinweis 1), erst dann M4 `tap-check` unmittelbar vor M5.
3. **M5** Tag nur auf dem in M2 geprüften Commit, nur nach gesonderter Freigabe; `--tags` nicht; Tag danach nicht verschieben.
4. **M6/M7/M8:** `publish`, GHCR, Assets, Job `tap`, Mac, Roadmap mit Tag-Hash, `releasing.md`-Zeiger, Closure-Notiz mit Lerneintrag, `git mv` nach `done/`, `make gates`; danach `close-welle`.

Urteil: **bis M1 schließbar (GO Push); der Slice selbst nicht schließbar vor M2 bis M8.**
