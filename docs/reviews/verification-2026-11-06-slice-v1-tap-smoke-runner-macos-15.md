# Verification-Report (erster Teil): slice-v1-tap-smoke-runner-macos-15 — 2026-11-06

**Rolle:** Verifier (frischer Kontext; Baseline-Regelwerk `modul-11`)

**Gegenstand:** [`slice-v1-tap-smoke-runner-macos-15`](../plan/planning/done/slice-v1-tap-smoke-runner-macos-15.md), Commit-Spanne `34de6de2^..3e1d81a9` (Plananlage bis lokaler Tip). Die Quelländerung ist `737dae18` (eine Datei, 4 Einfügungen, 1 Löschung); `3e1d81a9` legt den Review-Report ab und berichtigt den Plan (F-1 bis F-3). **Abgleichstand nach der letzten Fix-Runde = `3e1d81a9`**; der Arbeitsbaum war vor und nach allen Läufen sauber (`git status --short` leer).

**Umfang dieses Teils.** Nach Plan §2 (Zeile „Verification-Report“) prüft dieser Teil vor dem Push in den Tap: die Quelle, die Gates und die Plan-Tabelle. Die Liefer-Punkte **Tap** und **Smoke** stehen aus; der Hauptagent führt sie nach diesem Bericht aus, ein zweiter, nur lesender Teil prüft sie nach. Dieser Bericht behauptet für beide Punkte nichts.

**Eingangs-Kontext:**

- Slice-Plan §1 (Planentscheidung `macos-15` statt `macos-latest`, Nicht-Ziele), §2 (DoD), §3 (Plan-Tabelle), §4 bis §8
- [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung) (Paketierung: Homebrew, eigener Tap), [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) (Entscheidung 5: Smoke im Tap; nennt kein Runner-Label)
- Review-Report [`2026-11-06-slice-v1-tap-smoke-runner-macos-15`](2026-11-06-slice-v1-tap-smoke-runner-macos-15.md) (0 HIGH / 0 MEDIUM / 3 LOW / 2 INFO; Eingang, nicht Evidence)
- [`harness/verification.md`](../../harness/verification.md), [`harness/roles.md`](../../harness/roles.md), [`AGENTS.md`](../../AGENTS.md) (§3.1, §3.2, §3.6, §3.7, §3.9, §3.10, §5)
- Vorbild der Form: [`verification-2026-11-06-slice-v1-release-cut-v0.24.0`](verification-2026-11-06-slice-v1-release-cut-v0.24.0.md)

---

## Verification Evidence

Scope:
- Slice: [`slice-v1-tap-smoke-runner-macos-15`](../plan/planning/done/slice-v1-tap-smoke-runner-macos-15.md)
- IDs: [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung), [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md)
- Artefakte: `packaging/homebrew-tap/.github/workflows/smoke.yml` (Quelle); Tap-Repo `pt9912/homebrew-u-boot`, `.github/workflows/smoke.yml` (extern, noch unverändert)

### DoD-Abgleich (Plan §2)

Die Häkchen im Slice stehen alle offen; sie gehören dem Planner. Der Zustand je Punkt:

| # | DoD-Punkt | Zustand | Evidence |
|---|---|---|---|
| 1 | **Quelle:** Matrix `runner: [macos-15, macos-15-intel]`; Diff = diese Zeile und ein Kommentar, sonst nichts | **erbracht, eigene Lesung** | `git show 737dae18`: ein Hunk, `-        runner: [macos-14, macos-15-intel]`, `+` drei Kommentarzeilen und `+        runner: [macos-15, macos-15-intel]`. Endstand gelesen: Kommentar Zeilen 24 bis 26, Matrix Zeile 27, gleiche Einrückung innerhalb von `matrix:`. Trigger (Zeilen 10 bis 14), `permissions` (16/17), alle Schritte und das gatende `brew audit --strict --online` (Zeilen 61 bis 69, `exit 1` bei Befund) sind nicht im Diff. Blob der Quelle am Tip: `db62a07fb710fb2752857b76bcb09422806a2800` |
| 2 | **Tap:** Datei im Tap byte-gleich mit der Quelle (`cmp`, Exit 0), Tap-Commit in §7 | **OFFEN — erwartet offen** | Push steht nach diesem Bericht aus (Plan §3, Ansatz). Ausgangslage gelesen, siehe *Ausgangslage des Tap* |
| 3 | **Smoke:** Lauf `smoke` im Tap, `install (macos-15)` und `install (macos-15-intel)` `success`; arm64-Beleg aus Job-Metadaten und „Set up job“; Formel-Version genannt | **OFFEN — erwartet offen** | Lauf steht nach dem Push aus; Nachlauf im zweiten Teil. Welche Messstellen der Nachlauf liest, steht in *Übergabe* |
| 4 | `make gates` grün (Devcontainer-Container) | **erbracht, eigener Lauf** | siehe *Sensors* |
| 5 | Review durchgeführt, Report liegt vor | **erbracht** | [Review-Report](2026-11-06-slice-v1-tap-smoke-runner-macos-15.md), Gegenstand `737dae18`, in `3e1d81a9` abgelegt; nicht der Implementer (Rollenwechsel) |
| 6 | Verification-Report vor dem Push (Quelle, Gates, Plan-Tabelle) | **erbracht mit diesem Bericht**; Nachlauf für Tap und Smoke offen | dieser Bericht |
| 7 | Doku-Update: keine Datei berührt; „unverändert“-Zeilen belegt, bei der Closure gegen `git diff --stat` abgeglichen | **erbracht für den Stand `3e1d81a9`**; Abgleich bei der Closure zu wiederholen | siehe *Plan-vs-Code*; kein öffentlicher Vertrag im Diff |
| 8 | Closure-Notiz mit Steering-Loop-Lerneintrag | **offen — Closure** | §7 trägt den Stand „noch nicht geschlossen“ (richtig vor Tap und Smoke) |
| 9 | Beobachtungs-Register fortgeschrieben | **offen — Closure** | `git diff --stat 34de6de2^..HEAD` berührt `docs/plan/planning/observations/` nicht; Plan §3 und §5 sehen das bei der Closure vor |
| 10 | Jedes Risiko aus §6 trägt einen Ausgang | **offen — Closure** | §6 nennt je Risiko, woran der Ausgang hängt |
| 11 | Die drei Paarungen (Anker · Folge-Slice · Register) | **offen — Closure** | — |

### Prüfung der Review-Übergaben

| Finding | Auflösung | Beleg |
|---|---|---|
| F-1 (LOW, Zeilenverweise im Plan nach dem Diff veraltet) | **behoben in `3e1d81a9`** | Plan nennt jetzt „im Stand vor der Änderung in Zeile 24“ und „im Endstand … Zeile 27“ und „Zeilen 61 bis 69, Schritt ‚Formel pruefen (gatend)‘, Endstand“. Gegen `smoke.yml` gelesen: Zeile 27 Matrix, Zeile 61 `- name: Formel pruefen (gatend)`, Zeile 69 `fi`; die Zeile 24 des Vorstands ist im Hunk von `737dae18` die gelöschte Zeile. „Zeilen 10 bis 13“ (Push-Trigger) und „10 bis 14“ (mit `workflow_dispatch`) stimmen am Endstand |
| F-2 (LOW, ASCII-Aussage ungenau) | **behoben in `3e1d81a9`** | Der Plan sagt jetzt: genau ein Nicht-ASCII-Zeichen, der Gedankenstrich in Zeile 6. Eigene Suche (`grep -nP '[^\x00-\x7F]'` über `smoke.yml`): einziger Treffer Zeile 6 |
| F-3 (LOW, §7 „liegt in `open/`“) | **behoben in `3e1d81a9`** | §7 sagt „liegt in `in-progress/`“, nennt `737dae18` und den Review-Report. Die Datei liegt in `docs/plan/planning/in-progress/` (`ls`) |
| F-4 (INFO, Fremdzustand bis zur Messung) | **an den Nachlauf übergeben** | V-2 unten; der Lauf mit dem neuen Label bleibt der Beleg |
| F-5 (INFO, Kopplung Quelle und Tap ohne Sensor) | **an den Projektinhaber übergeben** (Plan §6) | Der Verifier hält den Zustand zum Zeitpunkt des Vergleichs fest; ein Verfahren ist nicht Teil dieses Slice |

Die Plan-Berichtigungen in `3e1d81a9` haben kein Nachlauf-Review; ich habe den Diff `737dae18..3e1d81a9` des Plans gegen die Datei gelesen (V-5).

### Ausgangslage des Tap (gelesen, nur lesend)

| Aussage | Beleg (`gh api`, 2026-11-06) |
|---|---|
| Datei im Tap = Stand vor der Änderung | `gh api repos/pt9912/homebrew-u-boot/contents/.github/workflows/smoke.yml --jq .content \| base64 -d` gegen `git show 026fefb5:packaging/homebrew-tap/.github/workflows/smoke.yml`: `cmp` Exit 0 |
| Datei im Tap ≠ Quelle am Tip (Erwartung vor dem Push) | `cmp` gegen `packaging/homebrew-tap/.github/workflows/smoke.yml`: `differ: byte 614, line 24`, Exit 1 |
| Tap-Kopf | `b9f23331a5` „u-boot 0.24.0“ (`commits/main`) |
| `main` im Tap nicht geschützt | `branches/main` → `protected` `false` |
| Rechte des Tokens am Tap-Repo | `permissions`: `push` und `admin` wahr. Ob das Token Dateien unter `.github/workflows/` pushen darf (eigenes Recht, Plan §6), ist **nicht gemessen**; die Messstelle ist der Push |

### Fremdzustand zum Anlass (Auftragspunkt: Behauptung ohne Bestätigung)

Der Plan stützt seinen Anlass auf `gh api`-Lesungen. Ich habe sie nachgelesen:

| Plan-Aussage | Befund |
|---|---|
| Lauf `42191356009`: Event `push`, Kopf `b9f23331a5` „u-boot 0.24.0“, `run_attempt` 2, `cancelled` | bestätigt |
| Job `127811577331` `install (macos-15-intel)` `success`, Runner vergeben | bestätigt (`runner_name` `GitHub Actions 1000036007`) |
| Job `127811574012` `install (macos-14)` `cancelled`, Label `macos-14`, `runner_name` leer | bestätigt |
| Annotation „The configuration 'macos-14-arm64-us-default' is not supported“ | bestätigt (`check-runs/127811574012/annotations`) |
| Läufe `41636802191` (2026-11-03) und `41844822254` (2026-11-04) `success` | bestätigt |
| README von `actions/runner-images`: `macos-15` = „macOS 15 Arm64“ (arm64), `macos-15-intel` = „macOS 15“ (x64), `macos-latest` = „macOS 26 Arm64“ | bestätigt für die Labels `macos-15` (Zeile „macOS 15 Arm64“, arm64) und `macos-15-intel` (Zeile „macOS 15“, x64); die `-latest`-Passage (Wanderung über 1 bis 2 Monate) steht im Abschnitt zum `-latest`-Label |
| Ankündigung vom 2026-10-01 zur Außerbetriebnahme von `macos-14` | **nicht nachgelesen** (kein Abruf der Blog-Seite in diesem Lauf); die Annotation und die README-Kennzeichnung der Zeile „macOS 14“ als „deprecated“ stützen sie |

Dass die Außerbetriebnahme die Ursache ist, bleibt bis zum grünen Lauf mit dem neuen Label eine Annahme (Plan §6, erstes Risiko).

### Entscheidungs-Konformität (Auftragspunkt: ADR und Spec)

| Quelle | Ergebnis | Beleg |
|---|---|---|
| [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) Entscheidung 5 (Smoke im Tap, Datei `smoke.yml` als Vorlage im Hauptrepo) | **konform** | Zeile 29 gelesen: nennt Datei und Prüfungen, kein Label; die Änderung bleibt in der Vorlage und (nach dem Push) im Tap. Keine Datei unter `docs/plan/adr/` im Diff |
| [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung) (Homebrew, eigener Tap) | **konform** | Distributionsweg unverändert; Binary, Formel und Pakete im Diff nicht berührt |
| [`AGENTS.md`](../../AGENTS.md) §3.6 (keine Gate-Lockerung) | **bestätigt** | `brew audit --strict --online` mit `exit 1`, kein `continue-on-error`, kein `\|\| true` (Zeilen 61 bis 69 gelesen); `Makefile`, `.golangci.yml`, `.d-check.yml`, `.github/workflows/` nicht im Diff |
| [`AGENTS.md`](../../AGENTS.md) §3.2 (Suppression) | **bestätigt** | kein `//nolint`; `git diff --stat` außerhalb `docs/` und `packaging/` ist leer |
| [`AGENTS.md`](../../AGENTS.md) §3.5 (Accepted-ADRs unverändert) | **bestätigt** | keine ADR im Diff |
| [`AGENTS.md`](../../AGENTS.md) §3.7, §3.9 (Kommentar) | **bestätigt** | die drei Kommentarzeilen sind Indikativ und tragen Abgrenzung (kein wanderndes `macos-latest`) und Kopplung (Vorlage und Tap); keine Chronik, kein „früher“, kein „ohne X wäre“. Der alte Wert `macos-14` steht im Baum nur noch in `CHANGELOG.md` Zeile 454 (Abschnitt `[0.10.0]`, Stand dieses Release), in Abschlussnotizen unter `docs/plan/planning/done/`, in Reports und im Anlass des Slice |
| [`AGENTS.md`](../../AGENTS.md) §5 (öffentliche Verträge; Regel 6 Nutzer-Doku) | **bestätigt** | kein Runner-Label in `README.md`, `README.de.md`, `docs/user/`, `docs/maintainer/`; kein Nutzer-Doku-Diff |
| [`AGENTS.md`](../../AGENTS.md) §3.1 (Docker-only) | **bestätigt** | Nachweis ist `make gates` im Devcontainer-Container (`/.dockerenv` vorhanden); kein Host-`go test` |
| [`AGENTS.md`](../../AGENTS.md) §3.3, Lifecycle | **bestätigt** | der Slice liegt nur in `in-progress/` (`ls`); die Moves `a47d113b` und `d87425da` sind reine `git mv`-Commits laut Commit-Nachricht; Verweis-Nachzüge getrennt (`4b767878`, `ac5acc2c`) |

### Plan-vs-Code (§3 gegen `git diff --stat 34de6de2^..HEAD`, Stand nach der letzten Fix-Runde)

Diff: 5 Dateien, 385 Einfügungen, 2 Löschungen.

| Diff-Datei | Plan-Zeile (§3) | Abgleich |
|---|---|---|
| `packaging/homebrew-tap/.github/workflows/smoke.yml` (+4/−1) | update (Matrix-Zeile 27, Kommentar 24 bis 26) | deckt; Zeilennummern des Endstands gelesen |
| [Slice-Plan](../plan/planning/done/slice-v1-tap-smoke-runner-macos-15.md) (+300) | neu | deckt |
| `docs/plan/planning/open/README.md` (+1) | update (Bestandszeile) | deckt; Zeile 14 gelesen: verweist auf `../in-progress/`, Text „liegt in `in-progress/`, in Umsetzung“ |
| `docs/plan/planning/in-progress/roadmap.md` (−1) | update (Ruhe-Marker entfällt) | deckt; `git show 026fefb5`: gelöscht ist die Zeile „Nichts in Arbeit.“; sie steht nicht mehr in `roadmap.md` (Suche) |
| [Review-Report](2026-11-06-slice-v1-tap-smoke-runner-macos-15.md) (+80) | neu (Review-Report) | deckt |
| dieser Report (neu, nach `3e1d81a9`) | `docs/reviews/<Datum>-verification-slice-…` neu | **Name weicht ab** (V-1) |
| Tap-Repo `smoke.yml` | update (extern) | im Diff nicht enthalten, wie geplant; Push steht aus |
| `docs/plan/planning/observations/BEO-<KUERZEL>/<slug>/…` | neu oder update, bei der Closure | nicht im Diff, wie geplant |
| als „keine Änderung / unverändert“ geführt: Tests, `.github/workflows/*.yml`, `scripts/tap-nachzug.sh`, `scripts/tap-check.sh`, `packaging/homebrew-tap/README.md`, `docs/maintainer/releasing.md`, `CHANGELOG.md`, [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) | unverändert | im Diff nicht enthalten, deckt. Gelesene Stellen nachgeprüft: `packaging/homebrew-tap/README.md` Zeilen 19 bis 21 (nennt `.github/workflows/smoke.yml`, kein Label); `docs/maintainer/releasing.md` Zeilen 50/51 und 162 bis 164; `CHANGELOG.md` Zeile 12 (`[Unreleased]`), 14 (`[0.24.0]`), 423 (`[0.10.0]`), 454 (nennt `macos-14` und `macos-15-intel`); `scripts/tap-nachzug.sh` Zeilen 26/27 (`mkdir -p "$work/tap/Formula"`, `cp … Formula/u-boot.rb`), keine Nennung von `smoke` oder `workflows` in `tap-nachzug.sh` und `tap-check.sh`; `cmd/uboot/homebrew_scripts_test.go` Zeilen 21/35 (`runFill`), 120/129 (`runTapCheck`), 70 (Kommentar), 73 (liest `scripts/homebrew-formula.rb.tmpl`); `internal/adapter/driven/docker/inspector_test.go` Zeile 47 (`fake runner: …`). Suche nach `smoke.yml`, `macos-1`, `runner:` über `*.go`, `*.sh`, `*.mk`, `*.py`, `*.yml`, `*.yaml`, `*.rb`, `*.tmpl` außerhalb der Datei selbst: nur dieser eine Treffer. `.github/workflows/` nennt kein `macos` |

Befund zur Tabelle: keine Diff-Datei ohne Plan-Eintrag, kein Plan-Eintrag ohne Diff außer den als später oder extern geführten Zeilen (Tap, Register, dieser Report). Eine Abweichung: der Dateiname dieses Reports (V-1).

### Sensors

| Sensor | Ergebnis | Evidence |
|---|---|---|
| `make gates` (Devcontainer-Container, `/.dockerenv` vorhanden; Tip `3e1d81a9`, Arbeitsbaum sauber) | **pass** | eigener Lauf, Exit `0` (Log im Scratchpad, 1097 Zeilen). `tool-pruefung: … (go 1.27.3, golangci-lint 2.14.0, make >= 4.3) … — OK`; `gate-reihenfolge: tool-pruefung is the first check of record-gates — OK`; `gate-reihenfolge-test: 2 Faelle, 0 Fehlschlaege`; `baseline-verify: v6.18.1 OK — 54 Dateien`; `d-check: pinsync: 3 Pin(s) geprüft, 11 Quelle(n) gelesen`; `d-check: 1202 Datei(en) geprüft, 0 Befund(e)`; `0 issues.` (lint); `coverage-gate: OK — coverage 92.40% meets threshold 90%`; `ablage-pruefung: OK`; `ablage-pruefung-test: 26 Faelle, 0 Fehlschlaege`. Die Testpakete meldeten `(cached)`: der Diff enthält keine Go-Datei, der Cache ist für diesen Stand gültig |
| `make docs-check` nach Anlage dieses Reports | siehe *Sensor-Lauf* | — |
| Quelle gelesen (Matrix, Kommentar, übrige Datei, Nicht-ASCII-Suche) | **pass** | siehe DoD-Punkt 1, F-2 |
| `make test-docker` | not run | der Diff enthält keinen Go-Code und keine Datei, die ein Test liest (Suche oben); kein Docker-Pfad berührt |
| `make verify-depguard` | not run | keine Importregel berührt |
| `make ci`, `make fullbuild` | not run | kein Produkt-, Image- oder Abhängigkeitsdiff; `govulncheck` und `image-scan` haben keinen Gegenstand in dieser Änderung |
| YAML-Parser über `smoke.yml` | not run | im Container steht kein YAML-Parser außerhalb von `make`-Targets zur Verfügung (Host-Toolchain-Verbot, [`AGENTS.md`](../../AGENTS.md) §3.1); die Gültigkeit stützt sich auf die Lesung (Kommentarzeilen mit der Einrückung von `runner:` innerhalb von `matrix:`) und den Review. Das Parsen leistet erst der Lauf im Tap |
| Vergleich Quelle gegen Tap nach dem Push | not run | Gegenstand von DoD „Tap“, nach dem Push |
| Lauf `smoke` im Tap | not run | Gegenstand von DoD „Smoke“, nach dem Push |

### Traceability

| ID | Beleg |
|---|---|
| [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung) | Der Distributionsweg (Homebrew, eigener Tap) bleibt; belegt wird die Änderung der CI-Vorlage durch den Diff `737dae18` und, nach dem Push, durch den Lauf `smoke` im Tap. Es gibt keinen Test im Hauptrepo, der `smoke.yml` liest (Plan §1, Suche oben); das ist die Planentscheidung, kein Befund |
| [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) | Entscheidung 5 (Smoke im Tap): Datei und Prüfungen unverändert, nur das Label der Matrix; belegt durch Diff und Lesung; kein Label in der ADR |

### Carveouts

- Neu: none
- Gelöst: none
- Unverändert: none berührt (kein Eintrag unter `docs/plan/carveouts/` im Diff)

### Nicht ausgeführt

- **Tap-Push, `cmp` Quelle gegen Tap, Lauf `smoke`** — stehen nach diesem Bericht aus (Plan §3, Ansatz); keine Erfolgsmeldung für den Tap oder für macOS wird behauptet.
- **`make test-docker`, `make verify-depguard`, `make ci`, `make fullbuild`** — Begründung in der Sensor-Tabelle.
- **Ankündigung von GitHub zur Außerbetriebnahme** — nicht nachgelesen (V-2).
- **YAML-Parser über `smoke.yml`** — Begründung in der Sensor-Tabelle.

### Commit / Artefakt

- Quellcommit `737dae18`; Review- und Plan-Berichtigung `3e1d81a9` (lokaler Tip, `main`).
- Blob der Quelle `packaging/homebrew-tap/.github/workflows/smoke.yml` am Tip: `db62a07fb710fb2752857b76bcb09422806a2800`.
- Tap-Kopf vor dem Push: `b9f23331a5`. Kein Tap-Commit, kein Lauf, kein Tag, kein Release-Asset.

---

## Findings des Verifiers

| ID | Kategorie | Befund | Auflösung |
|---|---|---|---|
| V-1 | LOW | **Der Dateiname des Verification-Reports weicht von der Plan-Tabelle ab.** Plan §3 führt die Form `<Datum>-verification-slice-…` („Namensform wie die Berichte zu [`slice-v2-homebrew-smoke-haertung`](../plan/planning/done/slice-v2-homebrew-smoke-haertung.md)“, Beispiel `2026-10-17-verification-slice-…`). Dieser Report heißt `verification-2026-11-06-` plus Slice-Kennung, die Form der neueren Berichte (`ls docs/reviews`: `verification-2026-11-05-…`, `verification-2026-11-06-…`) und des Auftrags. Die Plan-Aussage über den Dateinamen stimmt damit nicht mit der Datei überein (`AGENTS.md` §3.10). | Planner berichtigt die Zeile in §3 auf den Ist-Namen (Tabelle gegen `git diff --stat` nach der letzten Fix-Runde). Kein Eingriff am Report. |
| V-2 | INFO | **Fremdzustand: bestätigt, soweit lesbar; Ursache nicht bewiesen.** Lauf, Jobs, Annotation, Vorläufe und die Label-Zeilen der README (`macos-15` = arm64, `macos-15-intel` = x64) stimmen mit dem Plan. Nicht nachgelesen: die Ankündigung zur Außerbetriebnahme. Dass `macos-15` einen Runner bekommt und der Lauf grün ist, kann erst der Lauf zeigen (Plan §6, F-4 des Reviews). | Nachlauf im zweiten Teil: Job-Metadaten (`labels`, `runner_name`), Image-Abschnitt von „Set up job“ per `gh api`, Formel-Version. |
| V-3 | INFO | **Schreibrecht für Workflow-Dateien ist nicht gemessen.** `push` und `admin` am Repo sind wahr, `main` ist ungeschützt; das Recht für `.github/workflows/` hängt bei GitHub an einem eigenen Scope des Tokens (Plan §6, Risiko „Schreibzugriff auf den Tap“). Die Messstelle ist der Push. Verweigert GitHub ihn, gilt die Rückführung nach `open/` (Plan §4). | Hauptagent: Ergebnis des Push (Annahme oder Ablehnung) berichten. |
| V-4 | INFO | **Der Zustandsmarker in §7 läuft mit diesem Bericht ab.** §7 sagt „Verification, Tap-Push und Smoke-Lauf stehen aus“; nach diesem Bericht steht die Verification nicht mehr aus. `open/README.md` Zeile 14 sagt „in Umsetzung“ und bleibt richtig. | Planner setzt §7 bei der Closure auf den Ist-Stand (`AGENTS.md` §3.10); die Aussage zum Tap-Lauf trennt dort, was der Planner las und was der Hauptagent berichtet. |
| V-5 | INFO | **Die Plan-Berichtigungen aus `3e1d81a9` haben kein Nachlauf-Review.** Der Review-Report liest `737dae18`; die Berichtigungen sind Text im Plan. Ich habe sie gegen die Datei gelesen (Zeilen 24 bis 27, 61 bis 69, Nicht-ASCII-Suche): keine Abweichung. Eine Verifier-Lesung, kein Review. | Entscheidung über ein Nachlauf-Review liegt beim Planner; aus dem Diff ergibt sich keine Abweichung. |

## Verdikt

**Erster Teil bestanden — die Quelle ist zum Push in den Tap bereit; der Slice bleibt in `in-progress/`, weil Tap und Smoke ausstehen.** Erbracht und belegt: DoD „Quelle“ (Diff besteht aus der Matrix-Zeile und einem Kommentar; Audit-Gate, Trigger, Schritte unberührt), `make gates` im Devcontainer-Container grün am Tip `3e1d81a9`, der Review mit allen drei LOW-Berichtigungen aufgelöst, die Plan-Tabelle gegen `git diff --stat` abgeglichen (eine Abweichung: der Dateiname dieses Reports, V-1), keine Gate-Lockerung, keine ADR oder Spec berührt.

**Offen, und damit den Slice blockierend:**

1. **Tap** — Push von `smoke.yml` (Blob `db62a07fb710fb2752857b76bcb09422806a2800`) in `pt9912/homebrew-u-boot`, danach `cmp` Exit 0 (V-3: Schreibrecht für Workflow-Dateien).
2. **Smoke** — `install (macos-15)` und `install (macos-15-intel)` `success`, arm64-Beleg aus Job-Metadaten und „Set up job“, Formel-Version (V-2).
3. **Closure** — Notiz mit Lerneintrag, Register, Risiko-Ausgänge, Paarungen; Plan-Tabelle berichtigen (V-1), §7 auf den Ist-Stand (V-4).

## Übergabe

**An den Hauptagenten (vor dem Push):** Die Quelle darf sich bis zum Push nicht mehr ändern; jede weitere Änderung an `smoke.yml` oder am Plan macht den Abgleich in diesem Bericht (Blob `db62a07f…`) ungültig und verlangt ihn neu. Push in den Tap als Kopie der Quelle, dann `gh api repos/pt9912/homebrew-u-boot/contents/.github/workflows/smoke.yml --jq .content | base64 -d | cmp - packaging/homebrew-tap/.github/workflows/smoke.yml`, Tap-Commit und Ausgang des Push berichten; Lauf mit `gh workflow run smoke --repo pt9912/homebrew-u-boot`, Run-ID berichten.

**An den zweiten, lesenden Teil (nach dem Lauf):** (1) `cmp` Quelle gegen Tap wiederholen; (2) Lauf per `gh api repos/pt9912/homebrew-u-boot/actions/runs/<id>` und `…/jobs`: Conclusion des Laufs, beider Jobs, `labels` `["macos-15"]` und `["macos-15-intel"]`, `runner_name` nicht leer; (3) „Set up job“ im Log per `gh api …/actions/jobs/<id>/logs` (Image-Abschnitt), nicht allein per `gh run view --log` ([`BEO-DOK/ci-log-falsch-gelesen`](../plan/planning/observations/BEO-DOK/ci-log-falsch-gelesen/observation.md)); die Aussage „arm64“ nur, wenn das Log sie trägt; (4) Annotationen der Jobs; (5) Formel-Version aus dem Schritt „Version“ (`u-boot --version`) und kein „macOS grün“ darüber hinaus ([`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md)); (6) bei einem Befund unter `macos-15` (Audit, Homebrew-Version) bleibt das Gate, der Befund wird ein eigener Vorgang (Plan §1, §6).

**An den Planner:** V-1 (Name der Verification-Zeile in §3), V-4 (§7 bei der Closure), V-5 (Nachlauf-Review der Plan-Berichtigungen: ja oder nein). Keinen Push, keinen Lauf und keine Änderung an Plan oder Quelle habe ich ausgelöst.

## Sensor-Lauf

`make docs-check` nach Anlage dieses Reports: siehe Ergebniszeile am Ende.

`make docs-check` nach Anlage dieses Reports: `1203 Datei(en) geprüft, 0 Befund(e)` (die ersten Läufe fanden unverlinkte Slice-Kennungen in der Tabelle *Plan-vs-Code* und in V-1; berichtigt).
