# Verification-Report: slice-v1-devcontainer-brew-ripgrep-im-repo — 2026-11-06

**Rolle:** Verifier (frischer Kontext; Baseline-Regelwerk `modul-11`)

**Gegenstand:** [`slice-v1-devcontainer-brew-ripgrep-im-repo`](../plan/planning/in-progress/slice-v1-devcontainer-brew-ripgrep-im-repo.md), Commit-Spanne `7fbad506..ee2db629` (Start-Commit bis lokaler Tip). Der Umsetzungs-Commit ist `9d255cbb` (4 Dateien, 39 Einfügungen, 5 Löschungen); `ee2db629` ändert nur den Review-Report und eine Zeile der Plan-Tabelle (`git diff --stat 9d255cbb ee2db629`: 2 Dateien, 80 Einfügungen). **Abgleichstand nach der letzten Fix-Runde = `ee2db629`** (Arbeitsbaum vor diesem Report sauber).

**Push-Stand (gelesen am 2026-11-06, 07:42 UTC):** `origin/main` steht auf `9d255cbb`; `ee2db629` ist **nicht gepusht**. Der CI-Beleg dieses Reports gilt dem **Push-Tip `9d255cbb`**, nicht dem lokalen HEAD (siehe V-1).

**Eingangs-Kontext:**

- Slice-Plan §1 (Ziel, Entscheidungen E1 bis E3, Nicht-Ziele), §2 (DoD, Closure-Pflichten), §3 (Plan-Tabelle), §4 und §5 (Trigger), §6 (Risiken R1 bis R9, E1 bis E5, F1 und F2), §7 (leer, Platzhalter), §8
- [`LH-FA-DEV-013`](../../spec/lastenheft.md#lh-fa-dev-013--homebrew-formeln-im-devcontainer), [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette)
- [`ADR-0046`](../plan/adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md), [`ADR-0043`](../plan/adr/0043-nur-devcontainer-bau-pfad.md); [`SPEC-057`](../../spec/spezifikation.md#spec-057--homebrew-formeln-im-devcontainer-katalog-schlüssel-brew-pfade-freigabe-der-taps-konflikt-mit-homebrew-pin-meldungen-hosts-und-grenzen) (Dockerfile-Block, Pin, Gemessen, Nicht zugesagt)
- Review-Report [`2026-11-06-slice-v1-devcontainer-brew-ripgrep-im-repo`](2026-11-06-slice-v1-devcontainer-brew-ripgrep-im-repo.md) (0 HIGH / 0 MEDIUM / 1 LOW / 3 INFO; Eingang, nicht Evidence); seine Übergaben F-2 bis F-4 an den Verifier sind unten aufgelöst
- Vorbild der Form: [`verification-2026-11-05-slice-v1-devcontainer-homebrew-pakete`](verification-2026-11-05-slice-v1-devcontainer-homebrew-pakete.md); [`harness/verification.md`](../../harness/verification.md), [`harness/roles.md`](../../harness/roles.md), [`AGENTS.md`](../../AGENTS.md) (§3.1, §3.2, §3.6, §3.7, §3.9, §3.10, §5)

---

## Verification Evidence

Scope:
- Slice: [`slice-v1-devcontainer-brew-ripgrep-im-repo`](../plan/planning/in-progress/slice-v1-devcontainer-brew-ripgrep-im-repo.md)
- IDs: [`LH-FA-DEV-013`](../../spec/lastenheft.md#lh-fa-dev-013--homebrew-formeln-im-devcontainer), [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), [`ADR-0046`](../plan/adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md), [`ADR-0043`](../plan/adr/0043-nur-devcontainer-bau-pfad.md)
- Artefakte: `u-boot.yaml`, `.devcontainer/Dockerfile`, `.devcontainer/devcontainer.json`, `docs/maintainer/devcontainer-arbeiten.md`, `docs/plan/planning/open/README.md`, `docs/plan/planning/in-progress/roadmap.md`

### DoD-Abgleich (Plan §2)

Die Häkchen im Slice stehen alle offen; sie gehören dem Planner. Der Zustand je Punkt:

| # | DoD-Punkt | Zustand | Evidence |
|---|---|---|---|
| 1a | `u-boot.yaml`: `brew` mit `enabled: true`, `packages: [ripgrep]`, kein `homebrew` | **erbracht** | `git diff 7fbad506 9d255cbb -- u-boot.yaml`: genau der Austausch (+3/−2: `homebrew:`/`version: 2.0.5` gegen `brew:`/`packages:`/`- ripgrep`); `enabled: true` bleibt, kein `enabled: false`-Rest, übrige Schlüssel unberührt |
| 1b | `.devcontainer/Dockerfile` und `devcontainer.json` sind die Erzeugung aus `u-boot.yaml`, ohne Drift; zweiter Lauf ändert nichts; `doctor` `ok` | **erbracht, eigener Lauf** | Binary aus dem Baum (`BAU_COMPILE_OUT=… bash scripts/bau.sh compile`), Klon von `ee2db629` mit `origin` auf die Remote-URL gesetzt (sonst fehlt `postCreateCommand`, wie der Plan ausweist): `generate devcontainer --yes --no-interactive` → `git status --short` leer; zweiter Lauf → weiter leer; `doctor`: `devcontainer.features.allowlist` und `devcontainer.features.drift` je `✓`. Zahlen des Diffs (`7fbad506..9d255cbb`): Dockerfile +32/−0, `devcontainer.json` +1/−2, `u-boot.yaml` +3/−2 — wie in der Plan-Vorschau |
| 1c | Lock-Datei im Slice-Diff unberührt (E1) | **erbracht** | `git diff --quiet 7fbad506 HEAD -- .devcontainer/devcontainer-lock.json` endet mit `0`; letzte Änderung der Datei `12fe9778` (`git log`) |
| 1d | Drift-Aussage (Lock-Eintrag `homebrew` verwaist zwischen Commit und Lock-Schritt) | **erbracht als benannter Zwischenstand** | die Lock-Datei trägt weiter `homebrew:2.0.5`, `devcontainer.json` führt ihn nicht mehr; Suche nach `devcontainer-lock`/`lockfile` in `*.go`, `*.sh`, `*.mk`, `Makefile`, `*.yml`, `*.yaml`: keine Fundstelle — kein Gate liest die Datei, wie der Plan sagt. Das Verhalten der CLI bei einem verwaisten Eintrag ist ungemessen (Lock-Schritt) |
| 2 | Maintainer-Doku nennt `ripgrep` und nur Gemessenes | **erbracht, mit Vorbehalt Run-ID** | `docs/maintainer/devcontainer-arbeiten.md`: „installiert das Dockerfile beim Bau des Images `ripgrep` als Homebrew-Formel (Schlüssel `brew` in `u-boot.yaml`); die Formel ist nicht gepinnt“. Beleg des Verhaltens: CI-Log, `brew install --formula 'ripgrep'` ohne Fehler (Liefer-Punkt 3). Keine Aussage über Bauzeit oder Image-Größe (gelesen). **Vorbehalt:** der Plan verlangt die Run-ID „in der Closure“; §7 trägt sie noch nicht (Closure-Pflicht, nicht dieses Reports) |
| 3a | CI auf dem Push-Tip: `gates` grün, Schritt *Build the devcontainer image* grün, `RUN`-Schritte im Log, Run-/Job-ID, Bauzeit und Dauer (E3) | **erbracht für `9d255cbb`; Tip-Vorbehalt V-1** | Run [`42180354239`](https://github.com/pt9912/u-boot/actions/runs/42180354239) (`ci`, `success`), Job `127768279203` (`gates (lint + test + coverage-gate)`, `success`). Log über `gh api --allow-escape-sequences repos/pt9912/u-boot/actions/jobs/127768279203/logs` gelesen (nicht `gh run view --log`, [`BEO-DOK/ci-log-falsch-gelesen`](../plan/planning/observations/BEO-DOK/ci-log-falsch-gelesen/observation.md)). Block im Log: Schritt `[ 8/13]` (`apt-get install build-essential procps curl file git ca-certificates`, `/etc/profile.d/u-boot-brew.sh`), `[ 9/13]` (`git init`/`git fetch --depth 1 … bf464f57…`/`git checkout`/`test "$(git rev-parse HEAD)" = …`/`ln -s`; `DONE 1.7s`), `[10/13]` (`brew install --formula 'ripgrep'`; `ripgrep (15.2.0)` plus `zlib-ng-compat 2.3.3_1`, `bzip2 1.0.8`, `pcre2 10.49` als Bottles aus `ghcr.io`; `DONE 7.8s`). Kein `CACHED`-Schritt im Log (kalter Bau). Zeiten siehe Tabelle *Bauzeit und Dauer* |
| 3b | `integration-docker` auf demselben Push-Tip grün (E4) | **erbracht für `9d255cbb`** | Run [`42180353823`](https://github.com/pt9912/u-boot/actions/runs/42180353823) (`integration-docker`, `success`), Job `127768277532`, 07:21:41 bis 07:42:00 = 20 min 19 s (Vorläufer 23 min 16 s und 21 min 50 s, `timeout-minutes: 30`). Der Workflow baut `.devcontainer/` nicht; der Lauf belegt die Abwesenheit einer Nebenwirkung, nicht die Änderung |
| 3c | `rg` im `PATH` eines neu gebauten Containers (Bedingung für `done/`, E2) | **OFFEN — nicht erbracht** | Der Projektinhaber hat den Nachweis (Image- oder Container-Ebene) nicht geliefert; in diesem Container nicht erbringbar (Homebrew kommt hier noch aus dem OCI-Feature, kein `rg`). Das CI-Log belegt nur, dass `brew install` lief und `ripgrep 15.2.0` gegossen wurde, nicht `command -v rg` im fertigen Image und nicht die Login-/Nicht-Login-Shell. **Der Slice bleibt nach E2 in `in-progress/`.** |
| 3d | Lock-Schritt nach dem Neubau (E1) | **OFFEN — hängt an 3c** | Neubau über Editor/CLI außerhalb dieses Containers steht aus; Diff der Lock-Datei, Commit oder Befund fehlen |

Closure-Pflichten (Plan §2, zählen nicht mit):

| Pflicht | Zustand | Evidence |
|---|---|---|
| `make gates` grün, im Devcontainer-Container, nach der letzten Fix-Runde | **erbracht, eigener Lauf** | `make gates` am Stand `ee2db629` im Devcontainer-Container, Exit `0` (Sensor-Tabelle). Zusätzlich die CI auf `9d255cbb` (`make gates` im privilegierten `docker run` auf dem aus `.devcontainer/` gebauten Image, `success`) |
| Review durchgeführt, Report liegt vor, kein Self-Review | **erbracht** | [Review-Report](2026-11-06-slice-v1-devcontainer-brew-ripgrep-im-repo.md) liegt vor, Gegenstand `9d255cbb`; `ee2db629` (Doku-Nachzug F-1) hat kein Nachlauf-Review gelesen — vom Verifier gelesen: eine Zeile Plan-Tabelle, Inhalt stimmt (siehe Plan-vs-Code). Haken gehört dem Planner |
| Plan-Abgleich (`AGENTS.md` §3.10) | **erbracht, mit LOW V-2** | Tabelle unten |
| Closure-Notiz mit Steering-Loop-Lerneintrag | **offen** | §7 trägt den Platzhalter „Wird bei der Closure geschrieben“ |
| Beobachtungs-Register fortgeschrieben | **offen** | Kandidaten in §7 und §8 genannt; kein Evidence-Eintrag angelegt |
| Jedes Risiko trägt einen Ausgang; F1/F2 an Architect und Projektinhaber übergeben | **offen** | §6 trägt die Ausgänge als Vorab-Zuordnung (E1 bis E5), §7 nennt keine; Werte für R1, R2, R5 liefert dieser Report (unten) |
| Drei Paarungen (Anker · Folge-Slice · Register) | **offen** | Closure |

### Bauzeit und Dauer (E3), gelesen aus `gh api …/runs/<id>/jobs`

Eine Messung je Seite; keine Aussage über einen Trend.

| Größe | vorher (Run 42177321325, Tip `073362fb`) | nachher (Run 42180354239, Tip `9d255cbb`) |
|---|---|---|
| Schritt *Build the devcontainer image* | 45 s (06:53:44 bis 06:54:29) | **57 s** (07:21:44 bis 07:22:41) |
| Schritt *Run make gates in the devcontainer container* | 70 s | 47 s |
| Job `gates (lint + test + coverage-gate)` | 2 min 0 s | 1 min 49 s |
| Beitrag des Blocks laut Log (BuildKit-Schritte 8 bis 10) | — | 2,6 s + 1,7 s + 7,8 s ≈ 12 s |

Die +12 s des Bau-Schritts decken sich mit den ≈12 s der drei Block-Schritte. Dass der Job trotz längerem Bau kürzer endet, liegt am kürzeren `make gates`-Schritt (47 s gegen 70 s) und ist eine Streuung zwischen zwei Läufen, kein Effekt des Blocks. Gegen `timeout-minutes: 20` bleibt der Abstand groß. **Image-Größe im Runner: nicht gemessen.** Das Log trägt eine Zeile `Warning: Landlock ABI 10 or later is required to deny all network access; found ABI 8` aus dem `brew install` (Sandbox von Homebrew gegen den Kernel des Runners); der Schritt endet trotzdem mit `DONE`.

### Prüfung der Review-Übergaben

| Finding | Auflösung | Beleg |
|---|---|---|
| F-1 (LOW, `roadmap.md` fehlt in der Plan-Tabelle) | **behoben** | `ee2db629` fügt die Zeile `docs/plan/planning/in-progress/roadmap.md` mit Start-Commit `7fbad506` ein (`git diff 9d255cbb ee2db629` gelesen); die Roadmap-Änderung selbst (zwei Zeilen gelöscht, Ruhe-Marker) steht in `7fbad506` |
| F-2 (INFO, Doku-Satz vor dem Beleg) | **aufgelöst** | CI-Bau grün, Log zeigt `brew install --formula 'ripgrep'` ohne Fehler (3a); Run-ID gehört in §7 |
| F-3 (INFO, erster Bau der Kombination) | **aufgelöst für den Bau, offen für die Laufzeit** | Das Image mit `nestedRuntime`, `ollama`, `gh`, `go-toolchain`, `golangci-lint` und `brew` baute in der CI ohne Fehler, und `make gates` lief darin grün (darunter `docs-check` mit verschachteltem Docker). Beleg für das Register [`BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf`](../plan/planning/observations/BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf/observation.md). Nicht belegt: `arm64`, Container-Ebene mit den OCI-Features |
| F-4 (INFO, `PATH`-Vorrang) | **teilweise aufgelöst, Rest beim Projektinhaber** | In der CI-Image-Umgebung endet `tool-pruefung` grün („go 1.27.3, golangci-lint 2.14.0, make >= 4.3 … bash, curl, gh, git, docker client — OK“): die Gate-Werkzeuge werden nicht überdeckt. Nicht gelesen, welche Namen `bzip2`, `pcre2`, `zlib-ng-compat` und `ripgrep` in `/home/linuxbrew/.linuxbrew/bin` verlinken — das ist der `ls`-Befehl des Projektinhabers (3c) |

### Entscheidungs-Konformität (Auftragspunkt: ADR und Spec)

| Quelle | Ergebnis | Beleg |
|---|---|---|
| [`ADR-0046`](../plan/adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md) (Schlüssel `brew`, Pin als Commit, kein Installer-Skript) | **konform** | Block: Checkout von `Homebrew/brew` per `git fetch --depth 1` auf `HOMEBREW_REVISION`, Abbruch bei abweichendem `git rev-parse HEAD`, Installation unter `USER vscode`, `HOMEBREW_NO_AUTO_UPDATE=1` nur am `RUN`, kein `ENV HOMEBREW_*`, kein `curl … \| sh`. Pin `bf464f57275cae0be106b43d190c98564ac59b0e` gegen `git ls-remote … refs/tags/7.0.12` gelesen: gleich |
| [`SPEC-057`](../../spec/spezifikation.md#spec-057--homebrew-formeln-im-devcontainer-katalog-schlüssel-brew-pfade-freigabe-der-taps-konflikt-mit-homebrew-pin-meldungen-hosts-und-grenzen) (Schlüssel und Pfade, Konflikt, Freigabe der Taps) | **konform** | `brew` erzeugt keinen Schlüssel unter `features`; `homebrew` entfällt, damit kein Exit-Code-`10`-Konflikt (die Erzeugung lief); `devcontainer.brewTaps.allow` ungesetzt, eine Formel ohne Tap braucht keine Freigabe; die Warnung der Erzeugung nennt Quelle, Pin-Ort und die Ungepinntheit der Formel |
| Satz „gebaut ist sie nicht“ ([`SPEC-057`](../../spec/spezifikation.md#spec-057--homebrew-formeln-im-devcontainer-katalog-schlüssel-brew-pfade-freigabe-der-taps-konflikt-mit-homebrew-pin-meldungen-hosts-und-grenzen), Handbuch; E5) | **unverändert, konform** | `spec/spezifikation.md` und `docs/user/` nicht im Diff |
| [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) (Pin-Inventur) | **nicht berührt, bestätigt** | `ARG GO_VERSION=` und `ARG GOLANGCI_LINT_VERSION=` im Dockerfile unverändert; `make gates` einschließlich `tool-pruefung` und `docs-check` (`pinsync`) grün |
| Gate nicht gelockert ([`AGENTS.md`](../../AGENTS.md) §3.2, §3.6) | **bestätigt** | `git diff --stat 7fbad506..ee2db629` ohne `Makefile`, `repo.mk`, `.golangci.yml`, `.d-check.yml`, `.github/`, `scripts/`, `docs/plan/carveouts/`; Coverage 92,40 % gegen Schwelle 90 %; kein `//nolint` im Diff |
| Accepted-ADRs unverändert ([`AGENTS.md`](../../AGENTS.md) §3.5) | **bestätigt** | keine Datei unter `docs/plan/adr/` im Diff |
| Docker-only ([`AGENTS.md`](../../AGENTS.md) §3.1) | **bestätigt** | alle Nachweise dieses Reports sind `make`-Targets, das Skript `scripts/bau.sh` (Bau-Pfad, [`ADR-0043`](../plan/adr/0043-nur-devcontainer-bau-pfad.md)) oder CI-Läufe; kein Host-`go test` als Nachweis |

### Plan-vs-Code (§3 gegen `git diff --stat 80beb49b~1..ee2db629`, Stand nach der letzten Fix-Runde)

| Diff-Datei | Plan-Zeile (§3) | Abgleich |
|---|---|---|
| `u-boot.yaml` (+3/−2) | `u-boot.yaml` update | deckt |
| `.devcontainer/Dockerfile` (+32) | update (erzeugt), +32 hinter Zeile 105 | deckt, Zahl und Ort wie die Vorschau |
| `.devcontainer/devcontainer.json` (+1/−2) | update (erzeugt) | deckt |
| `docs/maintainer/devcontainer-arbeiten.md` (+3/−1) | update, Einleitung | deckt |
| `docs/plan/planning/open/README.md` (+1/−1) | update | deckt |
| `docs/plan/planning/in-progress/roadmap.md` (−2) | update (Start-Commit `7fbad506`), seit `ee2db629` | deckt |
| Slice-Datei (Move, +1 Zeile) | „diese Datei“, neu | deckt |
| [Review-Report](2026-11-06-slice-v1-devcontainer-brew-ripgrep-im-repo.md) (+79, neu) | — | **keine Zeile** (V-2) |
| als „unverändert“ geführt: `devcontainer-lock.json`, `.d-check.yml`, `tool-pruefung.sh`, `repo.mk`, `Makefile`, `bau.sh`, `ci.yml`, `integration.yml`, zwei Test-Dateien, `CHANGELOG.md`, `README*.md`, Handbuch, `spec/spezifikation.md` | je „unverändert (gelesen)“ | **im Diff nicht enthalten, deckt** (`git diff --stat` nennt keine davon) |

### Sensors

| Sensor | Ergebnis | Evidence |
|---|---|---|
| `make gates` (Devcontainer-Container, Stand `ee2db629`) | **pass** | Exit `0`; `tool-pruefung … OK`, `gate-reihenfolge … OK`, `baseline-verify: v6.18.1 OK — 54 Dateien`, `d-check: 1191 Datei(en) geprüft, 0 Befund(e)`, `0 issues.` (lint), `coverage-gate: OK — coverage 92.40% meets threshold 90%`, `ablage-pruefung-test: 26 Faelle, 0 Fehlschlaege` |
| `make docs-check` nach Anlage dieses Reports | **pass** | `1192 Datei(en) geprüft, 0 Befund(e)` (Abschnitt *Sensor-Lauf*) |
| CI `ci` auf `9d255cbb` | **pass** | Run 42180354239; Jobs `gates` (127768279203), `security-gates (govulncheck)` (127768279096), `image-scan (trivy HIGH+CRITICAL)` (127768278879) je `success` |
| CI `integration-docker` auf `9d255cbb` | **pass** | Run 42180353823, Job 127768277532, 20 min 19 s |
| Erzeugung aus `u-boot.yaml` ohne Drift, zweiter Lauf stabil, `doctor` | **pass** | eigener Lauf, siehe DoD 1b |
| Pin gegen Release-Tag | **pass** | `git ls-remote https://github.com/Homebrew/brew refs/tags/7.0.12` |
| `make test-docker` lokal | not run | Der Workflow `integration-docker` ist derselbe Sensor auf dem Push-Tip und lief grün; der Diff enthält keinen Go-Code |
| `make verify-depguard` | not run | keine Importregel berührt |
| `docker build --no-cache` / `rg` im Container (E2) | not run | Sache des Projektinhabers (3c); Bau im Runner ist durch die CI belegt, `command -v rg` nicht |

### Traceability

| ID | Beleg |
|---|---|
| [`LH-FA-DEV-013`](../../spec/lastenheft.md#lh-fa-dev-013--homebrew-formeln-im-devcontainer) | `u-boot.yaml` `devcontainer.features.brew.packages: [ripgrep]`; Erzeugung ohne Drift; CI-Bau `brew install --formula 'ripgrep'` (Run 42180354239); die Anforderung trägt ihre Tests im Produkt (Tabelle *Belege* der Spec), der Slice fügt keinen hinzu |
| [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) | das eigene `.devcontainer/` ist die Erzeugung (`generate devcontainer` ohne Diff, `doctor` `drift ok`) |
| [`ADR-0046`](../plan/adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md) | Dockerfile-Block mit Pin-Prüfung, Ausweisung der Ungepinntheit; Review und dieser Report |
| [`ADR-0043`](../plan/adr/0043-nur-devcontainer-bau-pfad.md) | die CI baut `.devcontainer/Dockerfile` und fährt `make gates` darin (Log 127768279203) |

### Carveouts

- Neu: none
- Gelöst: none
- Unverändert: none berührt (kein Eintrag unter `docs/plan/carveouts/` im Diff)

### Nicht ausgeführt

- **`rg` im `PATH` eines neu gebauten Containers (Bedingung E2)** — der Projektinhaber hat den Nachweis nicht erbracht; in diesem Container nicht erbringbar.
- **Lock-Schritt** (Diff von `.devcontainer/devcontainer-lock.json` nach dem Neubau) — setzt den Neubau außerhalb dieses Containers voraus.
- **`make test-docker` lokal** — durch `integration-docker` auf dem Push-Tip abgedeckt.
- **`make verify-depguard`** — keine Importregel berührt.
- **Image-Größe** — nicht gemessen (kein Wert im Runner-Log, kein Lauf des Projektinhabers).

### Commit / Artefakt

- Umsetzung `9d255cbb`; Doku-Nachzug und Review-Report `ee2db629` (lokal, nicht gepusht); CI-Läufe 42180354239 (`ci`) und 42180353823 (`integration-docker`) auf `9d255cbb`.

---

## Findings des Verifiers

| ID | Kategorie | Befund | Auflösung |
|---|---|---|---|
| V-1 | MEDIUM | **Der geprüfte CI-Tip ist nicht der lokale Tip.** `origin/main` = `9d255cbb`; HEAD `ee2db629` und dieser Report sind ungepusht. Liefer-Punkt 3 verlangt die grünen Läufe „auf dem Push-Tip“ mit Run-IDs in §7. Dass `9d255cbb` grün ist, ist belegt; dass es der Tip bleibt, nicht. Die Spanne `9d255cbb..ee2db629` ist reine Doku, aber `docs-check` ist Teil des PR-blockierenden `gates`. | Nach dem Push der Closure-Commits die Läufe von `ci` und `integration-docker` auf dem neuen Tip erneut lesen (Run-IDs in §7); der Lock-Schritt erzeugt ohnehin einen neuen Tip (DoD 3d). Lokal ist `make gates` am Stand `ee2db629` grün. |
| V-2 | LOW | Die Plan-Tabelle in §3 führt die Review-Reports nicht, die der Diff mitbringt (`docs/reviews/2026-11-06-…`, dazu dieser Report). Die Vorgänger-Slices führen sie als Zeile „Nicht im Diff des Implementers“ (z. B. [`slice-v1-devcontainer-homebrew-pakete`](../plan/planning/done/slice-v1-devcontainer-homebrew-pakete.md) §3). | Planner/Implementer: Zeilen nach `AGENTS.md` §3.10 ergänzen, bevor der Slice nach `done/` geht. |
| V-3 | INFO | Messwerte für §7 liegen vor (Tabelle *Bauzeit und Dauer*): Bau-Schritt 57 s gegen 45 s, Block ≈ 12 s, Job 1 min 49 s gegen 2 min 0 s, `integration-docker` 20 min 19 s; Bottles statt Quelltext-Bau (`Pouring … bottle`), also für R2 gelesen: Homebrew zieht für `amd64`/Debian Bottles. Image-Größe ungemessen. | Closure-Notiz übernimmt die Werte; R1 und R2 „hingenommen“ (E3) mit diesen Zahlen. |
| V-4 | INFO | Die Zeile `Warning: Landlock ABI 10 or later …` im Log des Bau-Schritts ist eine Warnung von Homebrew gegen den Kernel des Runners, kein Fehler. Eine Wiederholung auf anderen Runnern ist ungelesen. | Nur nennen; bei Wiederholung Beobachtung im Register. |

## Verdikt

**DoD nicht vollständig — der Slice bleibt in `in-progress/`.** Erbracht und belegt: Liefer-Punkt 1 (Konfiguration, Erzeugung ohne Drift, Lock-Datei unberührt), Liefer-Punkt 2 (Doku, Satz durch den grünen CI-Bau gedeckt), der CI-Teil von Liefer-Punkt 3 (`gates` und `integration-docker` grün auf `9d255cbb`, Block im Log, Bauzeit und Dauer gelesen) und die Gate-Pflicht (`make gates` grün am Stand `ee2db629`).

**Offen, und damit die Closure blockierend:**

1. **`rg`-Nachweis des Projektinhabers (E2)** — Bedingung für `done/`; liegt nicht vor. Das CI-Log ersetzt ihn nicht.
2. **Lock-Schritt (E1)** — hängt am Neubau; Diff, Commit oder Befund fehlen.
3. **Push-Tip (V-1)** — nach dem Push der Commits die Läufe auf dem neuen Tip lesen.
4. **Closure-Notiz, Register, Risiko-Ausgänge, Paarungen** — §7 ist ein Platzhalter; die Werte dafür liegen in diesem Report.

**Übergabe an den Planner:** Slice nicht schließen. Reihenfolge laut Plan §3: Projektinhaber erbringt `rg`-Nachweis (Befehle stehen in DoD 3c) und Neubau → Lock-Schritt → Push → Läufe auf dem neuen Tip lesen → Closure (§7, Plan-Abgleich einschließlich V-2) → Hakenvergabe.

## Sensor-Lauf

`make docs-check` nach Anlage dieses Reports: `1192 Datei(en) geprüft, 0 Befund(e)` (der erste Lauf fand zwei unverlinkte Kennungen in den Tabellen *Entscheidungs-Konformität* und *Plan-vs-Code*; berichtigt).
