# Verification-Report: slice-v1-tap-nachzug-verfahren-doku — 2026-11-06

**Rolle:** Verifier (frischer Kontext; Baseline-Regelwerk `modul-11`)

**Gegenstand:** [`slice-v1-tap-nachzug-verfahren-doku`](../plan/planning/in-progress/slice-v1-tap-nachzug-verfahren-doku.md), Commit-Spanne `8074e4cf^..8cc97257` (Plananlage bis Tip). Die Umsetzung ist `11fc3b01` (zwei Dateien), der Review-Report liegt in `8558da8b`, der Review-Fix in `8cc97257` (eine Datei, `releasing.md`). **Abgleichstand nach der letzten Fix-Runde = `8cc97257`**; der Arbeitsbaum war vor und nach allen Läufen sauber (`git status --short` leer).

**Umfang.** Der Slice schreibt nichts in den Tap; es gibt keinen Nachlauf-Teil. Dieser Bericht prüft die Dokumentation, die Gates, die Plan-Tabelle, den Tap-Zustand (nur `gh api` GET) und führt die lesenden Befehle von `releasing.md` §9 selbst aus. Der Schreibschritt (Schritt 3), der Dispatch (Schritt 5, erste Zeile) und der Rückweg (Schritt 6) sind von mir nicht ausgeführt; das ist Auftragsgrenze, kein Befund.

**Eingangs-Kontext:**

- Slice-Plan §1 bis §8 (DoD, Plan-Tabelle, Entscheidungen, Risiken)
- [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung) (Paketierung: Homebrew, eigener Tap), [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) (Entscheidung 5: Smoke im Tap)
- Review-Report [`2026-11-06-slice-v1-tap-nachzug-verfahren-doku`](2026-11-06-slice-v1-tap-nachzug-verfahren-doku.md) (0 HIGH / 1 MEDIUM / 2 LOW / 4 INFO; Eingang, nicht Evidence)
- Entscheidungen des Hauptagenten zu Plan §4 (Angabe im Auftrag, kein Artefakt): Ort `releasing.md` §9; Sichtschritt in §4 aufgenommen; Drift-Gefahr nur als Zeiger auf Variante A; Push-Aufruf aus dem Sitzungsprotokoll
- [`harness/verification.md`](../../harness/verification.md), [`harness/roles.md`](../../harness/roles.md), [`AGENTS.md`](../../AGENTS.md) (§3.1, §3.7, §3.9, §3.10, §5)

---

## Verification Evidence

Scope:
- Slice: [`slice-v1-tap-nachzug-verfahren-doku`](../plan/planning/in-progress/slice-v1-tap-nachzug-verfahren-doku.md)
- IDs: [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung), [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md)
- Artefakte: `docs/maintainer/releasing.md` (§2, §4, §8, neuer §9), `docs/maintainer/README.md` (Zeile 10); Tap-Repo `pt9912/homebrew-u-boot` (nur gelesen)

### DoD-Abgleich (Plan §2)

Die Häkchen im Slice stehen alle offen; sie gehören dem Planner. Der Zustand je Punkt:

| # | DoD-Punkt | Zustand | Evidence |
|---|---|---|---|
| 1 | **Verfahren:** §9 mit (a) Geltung, (b) Rechte-Hinweis, (c) sechs Schritte, (d) Grenze | **erbracht bis auf eine wörtliche Abweichung (V-1)** | (a) bis (d) im Diff `11fc3b01` und `8cc97257` gelesen; Titel „Dateien des Tap-Repos nachziehen (`packaging/homebrew-tap/`)“ stimmt mit der DoD. (b) Zitat gegen die GitHub-Seite gelesen, siehe *Fremdzustand*. (c) Schritt 3 ohne `branch` (offengelegt, V-4). (d) im Indikativ, ohne „nur/nie/immer“ als Zusage. **Abweichung:** der Abschnitt trägt die Commit-Kennung `662c5fb5`, die DoD verbietet Commit-IDs „als Beleg“ (V-1) |
| 1a | „Jeder lesende Befehl des Abschnitts ist vom Implementer gegen den Tap ausgeführt (Auszug in §7)“ | **eigene Ausführung durch den Verifier; §7-Auszug steht aus (Closure)** | alle lesenden Befehle der Schritte 1, 3 (Vorbedingung), 4, 5 und die lokale `B64`-Zeile selbst ausgeführt, siehe *Lesende Befehle* |
| 1b | Schreibbefehl als ungeprobt ausgewiesen; Wortlaut vom Hauptagenten | **teilweise (V-2)** | „Gelaufen und nicht gelaufen“ weist Schritt 3 für `README.md`, den Rückweg und den Dispatch als nicht gelaufen aus; der Wortlaut des gelaufenen Aufrufs liegt als Artefakt nicht vor, nur sein Ergebnis ist lesbar |
| 2 | **Auffindbarkeit:** vier Stellen, je höchstens ein Satz, Zeiger auf §9 | **erbracht** | §2 (Klammer im Punkt „Tap-Repo“, `releasing.md` Zeile 51), §4 (ein Satz „Sichtschritt vor einem stabilen Tag, kein Gate: …“), §8 (eine Tabellenzeile hinter „Formel fehlt im Tap“), Index `docs/maintainer/README.md` Zeile 10 („Nachziehen der Tap-Dateien“; der Index nennt den Gegenstand, verlinkt keinen Anker). Nummerierung §1 bis §8 unverändert. Der Satz in §4 entspricht der bejahten Entscheidung 2 (V-3) |
| 3 | **Nichts im Tap und in `packaging/` geändert** | **erbracht, eigene Lesung** | `git diff --stat 8074e4cf^ HEAD -- packaging scripts .github` leer. Tap: `cmp` Exit `0` und sha gleich für beide Dateien, Tap-Kopf unverändert `662c5fb5…`, siehe *Tap-Zustand* |
| 4 | `make gates` grün (Devcontainer-Container) | **erbracht, eigener Lauf** | siehe *Sensors* |
| 5 | Review durchgeführt, Report liegt vor | **erbracht** | [Review-Report](2026-11-06-slice-v1-tap-nachzug-verfahren-doku.md), Gegenstand `11fc3b01`, abgelegt in `8558da8b`; Rollenwechsel (nicht der Implementer). Der Review-Fix `8cc97257` hat kein Nachlauf-Review (V-1 entstand dort) |
| 6 | Verification-Report liegt vor; Verifier führt die lesenden Befehle selbst aus und gleicht die Plan-Tabelle ab | **erbracht mit diesem Bericht** | dieser Bericht |
| 7 | Doku-Update: keine Datei unter `docs/user/`, kein README, kein CHANGELOG; „unverändert“-Zeilen belegt | **erbracht für den Stand `8cc97257`; Abgleich bei der Closure zu wiederholen** | siehe *Plan-vs-Code* |
| 8 | Closure-Notiz mit Steering-Loop-Lerneintrag | **offen — Closure** | §7 trägt „Noch nicht geschlossen“, richtig vor der Closure |
| 9 | Beobachtungs-Register fortgeschrieben (oder „keine Beobachtung“ in §7) | **offen — Closure** | `docs/plan/planning/observations/` ist im Diff der Spanne nicht berührt; Plan §3 sieht das für die Closure vor |
| 10 | Jedes Risiko aus §6 trägt einen Ausgang | **offen — Closure**, Empfehlung je Risiko unten | siehe *Risiken aus Plan §6* |
| 11 | Die drei Paarungen | **offen — Closure** | — |

### Befunde

| ID | Kategorie | Befund | Beleg | Übergabe |
|---|---|---|---|---|
| V-1 | MEDIUM (DoD, Verifier-only-Klasse) | **Eine DoD-Zusage ist wörtlich nicht erfüllt.** Die DoD „Verfahren“ sagt: der Abschnitt „enthält keine Run-, Commit- oder Job-IDs als Beleg“. Der Absatz „Gelaufen und nicht gelaufen“ trägt seit `8cc97257` die Commit-Kennung `662c5fb5` und die Prozessaussage „Reviewer und Verifier haben ihn gelesen“. Beides entstand im Review-Fix zu F-1 (der Review verlangte einen Bezugspunkt); der Review-Report prüfte „keine Commit-IDs“ am Stand `11fc3b01` (sein Negativbefund), der Fix hat kein Nachlauf-Review. Die Prozessaussage hat in der Dauer-Doku keinen auflösbaren Anker (kein Report-Verweis) und war zur Zeit des Fix-Commits unbelegt; sie ist seit diesem Bericht wahr (Commit gelesen, siehe *Tap-Zustand*). | `releasing.md` §9, Absatz „Gelaufen und nicht gelaufen“; Plan §2 DoD „Verfahren“, letzter Absatz; `git show 8cc97257` | Planner/Implementer: entweder die Kennung und die Prozessaussage aus §9 nehmen (Bezugspunkt bleibt Datum und „Hauptagent“; Kennung und Lesebefunde stehen in Plan §7) oder die DoD-Formulierung berichtigen. Beides ist Entscheidung des Planners, nicht des Verifiers |
| V-2 | LOW | **Herkunft des Push-Aufrufs ist nicht prüfbar, der abgedruckte Block ist nicht der gelaufene.** Der Wortlaut „aus dem Sitzungsprotokoll des Hauptagenten“ liegt in keinem Artefakt des Repos; der Absatz sagt selbst, die `B64`-Zeile des Blocks sei nicht gelaufen. Ich habe das Ergebnis gelesen (Tap-Commit auf `main`, Elternteil `b9f23331`, eine Datei, +4 −1, Patch gleich dem Quell-Diff `737dae18`), nicht den Aufruf. Das ist mit der Aussage des Abschnitts vereinbar, belegt die Befehlsform aber nicht. | `gh api repos/pt9912/homebrew-u-boot/commits/662c5fb5…`; Plan §6 Risiko 2 | Closure: Risiko 2 nicht als „entfallen“ führen, sondern „weiter offen → Grenze (d)“ (der Wortlaut entspricht nur in Teilen dem berichteten Aufruf) |
| V-3 | INFO | **Entscheidungen 1 bis 4 liegen nur als Angabe des Hauptagenten vor.** Plan §4 „Start“ nennt die Freigabe des Projektinhabers; ein Artefakt dazu fehlt (Review F-7). Konformität gelesen: Ort §9 (ja); Sichtschritt in §4 (ja, vorhanden); Drift-Gefahr nur als Zeiger (ja: keine neue Beobachtung im Register, §9 Grenze und Plan §6 zeigen auf Variante A); Push-Aufruf aus dem Sitzungsprotokoll (ja, so im Absatz ausgewiesen). | Plan §4; `git diff --stat` ohne `observations/` | Closure: in §7 vermerken, dass die Entscheidungen als Angabe des Hauptagenten vorlagen |
| V-4 | INFO | Plan DoD (c) Schritt 3 nennt `branch`; §9 lässt ihn weg (offengelegt). Plan DoD (c) Schritt 5 verlangt, den neuen Lauf an `created_at` **und** `head_sha` zu erkennen; §9 erkennt ihn an `created_at` und sagt, `head_sha` sei der Commit aus Schritt 3 und bei zwei Dispatch-Läufen auf demselben Commit gleich. Ich habe das gegen die Lauf-Liste gelesen: beide Dispatch-Läufe (`42212837961`, `42196563230`) tragen `662c5fb5…`; §9 ist an der Stelle genauer als der Plan. | `gh api …/workflows/smoke.yml/runs`; `releasing.md` §9 Schritt 5 | keine Änderung nötig; Plan-Wortlaut kann die Closure berichtigen |
| V-5 | INFO | Review F-5 (Annotationen als Rot-Signal) ist nicht aufgenommen und war nicht pflichtig. Ich habe es nachgelesen: der grüne Lauf trägt im Job `install (macos-15-intel)` eine Annotation der Stufe `failure` („Failed to import: …/hashicorp/homebrew-tap/Formula/vagrant.rb“) und zwei Warnungen, im Job `install (macos-15)` einen Hinweis („capacity constraints“). §9 nennt nur „not supported“ als Erkennungszeichen eines Jobs ohne Runner; es steht dort nicht, dass Annotationen allein kein Rot belegen. | `gh api …/check-runs/<job-id>/annotations` für beide Jobs des Laufs `42212837961` | optional, Folgevorgang |
| V-6 | INFO (außerhalb des Slice) | Die Roadmap (`in-progress/roadmap.md` Zeile 43) führt [`slice-v1-bau-nur-devcontainer`](../plan/planning/done/slice-v1-bau-nur-devcontainer.md) als „liegt in `in-progress/`“, die Datei liegt in `done/`. Vorbestand, nicht im Diff; die Closure stellt den Ruhe-Marker an dieser Stelle wieder her und sollte die Zeile beim Lesen sehen. | `ls docs/plan/planning/done/` (Datei vorhanden) | Hinweis an die Closure |

### Lesende Befehle von §9 (selbst ausgeführt, nur GET bzw. lokal)

| Schritt | Befehl | Ergebnis |
|---|---|---|
| 1 | `git hash-object packaging/homebrew-tap/$P` gegen `gh api …/contents/$P --jq .sha` | `smoke.yml`: beide `db62a07f…`; `README.md`: beide `5b333c7e…` (gleich) |
| 3 (Vorbedingung) | `git fetch origin`, `git status --short`, `git diff origin/main -- packaging/homebrew-tap/` | alle drei Ausgaben leer. Hinweis: `origin/main` stand bei `cf4c3cb8`, das lokale `HEAD` bei `8cc97257`; die Vorbedingung vergleicht nur `packaging/homebrew-tap/`, und der Slice ändert dort nichts |
| 3 (lokal, ohne `PUT`) | `B64=$(base64 -w0 < packaging/homebrew-tap/$P)` | 3368 Zeichen, ohne Zeilenumbruch; `printf %s "$B64" \| base64 -d \| cmp - <Quelle>` Exit `0` |
| 4 | `gh api …/contents/$P --jq .content \| base64 -d \| cmp - packaging/homebrew-tap/$P` und `--jq .sha` | `cmp` Exit `0` für beide Dateien; sha gleich dem aus Schritt 1 |
| 5 | Lauf-Liste mit `--jq '… @tsv'` | liefert id, `created_at`, `head_sha`, Event, Status, Ergebnis; der jüngste Dispatch-Lauf ist `42212837961` (`success`), daneben `42196563230` (`workflow_dispatch`, `queued`, ohne Ergebnis, `total_count` der Jobs `0`) |
| 5 | Jobs des Laufs mit `labels`, `runner_name`, `conclusion` | `install (macos-15)` / `macos-15` / `GitHub Actions 1000036052` / `success`; `install (macos-15-intel)` / `macos-15-intel` / `GitHub Actions 1000036051` / `success` |
| 5 | Annotationen und Log mit `--allow-escape-sequences` und `grep` | `Image: macos-15-arm64` und `Cellar/u-boot/0.24.0` und `u-boot version 0.24.0` im Job `macos-15`; `Image: macos-15` und `Cellar/u-boot/0.24.0` und `u-boot version 0.24.0` im Job `macos-15-intel`. Annotationen: siehe V-5 |
| 6 (Form, lokal) | `git show 737dae18^:packaging/homebrew-tap/$P \| base64 -w0` mit Rücklesen | Rücklesen gleich dem Blob; der Rückweg selbst ist nicht ausgeführt |
| Aussage „Der Push startet keinen Lauf“ | Läufe mit `head_sha=662c5fb5…` | zwei Läufe, beide `workflow_dispatch`; kein `push`-Lauf auf diesem Commit. Der letzte `push`-Lauf (`42191356009`) gehört zum Vor-Commit `b9f23331` und ist `cancelled` |

### Tap-Zustand (gelesen, nur lesend)

| Aussage | Beleg (`gh api`, 2026-11-06) |
|---|---|
| Datei im Tap = Quelle, beide Dateien | `cmp` Exit `0`; sha `db62a07f…` und `5b333c7e…` gleich (siehe oben) |
| Tap-Kopf | `commits/main` → `662c5fb5c57f…`, Nachricht „smoke: runner macos-15 statt macos-14 …“ |
| Tap-Commit `662c5fb5` | Elternteil `b9f23331`; Dateien: `.github/workflows/smoke.yml` mit +4 −1; Patch: Matrix `runner: [macos-14, macos-15-intel]` → `[macos-15, macos-15-intel]` plus drei Kommentarzeilen. Das stimmt mit dem Satz in §9 überein („nur `.github/workflows/smoke.yml` änderte (+4, −1)“) |
| Der Slice schrieb nichts in den Tap | Tap-Kopf gleich dem Stand der Plananlage (`662c5fb5…`); kein neuer Lauf nach `42212837961` |

### Fremdzustand zu den Rechten (Plan §6 Risiko 3)

Die Seite „Create or update file contents“ der GitHub-REST-Dokumentation (`docs.github.com/en/rest/repos/contents`, von mir am 2026-11-06 gelesen) sagt: klassische PATs und OAuth-App-Tokens brauchen den Scope `repo`, „The workflow scope is also required in order to modify files in the .github/workflows directory“; fine-grained: mindestens einer der Sätze „Contents (write)“ oder „Contents (write) and Workflows (write)“. Das Zitat in §9 stimmt damit überein. §9 sagt nichts über `HOMEBREW_TAP_GITHUB_TOKEN`. Die Rechte eines bestimmten Tokens sind weiterhin nicht gemessen.

### Plan-vs-Code (Plan-Tabelle §3 gegen `git diff --stat 8074e4cf^ HEAD`, Stand `8cc97257`)

| Diff-Datei | Plan-Tabelle |
|---|---|
| `docs/maintainer/releasing.md` (+108 −…) | ja, „update“ (§9 und drei Verweise) |
| `docs/maintainer/README.md` (1 Zeile) | ja, „update“ |
| [Plan-Datei des Slice](../plan/planning/in-progress/slice-v1-tap-nachzug-verfahren-doku.md) in `in-progress/` (neu) | ja, „neu“ |
| `docs/plan/planning/open/README.md` | ja, „update“ (Bestandszeile folgt dem Slice) |
| `docs/plan/planning/in-progress/roadmap.md` (−2: Ruhe-Marker) | ja, „unverändert bis zum Start; dann update, Nettodiff der Spanne leer“ — der Zustand ist für `in-progress` richtig; die Closure stellt den Marker wieder her, danach fehlt die Datei im `--stat` |
| [Review-Report](2026-11-06-slice-v1-tap-nachzug-verfahren-doku.md) (neu) | ja |
| dieser Bericht (neu, ungestaged) | ja, „neu“ |

Plan-Einträge ohne Diff: `docs/plan/planning/observations/…` (bedingt, Closure entscheidet); „Tests“, `smoke.yml`, `packaging/homebrew-tap/README.md`, `scripts/tap-nachzug.sh`, `scripts/tap-check.sh`, `publish.yml`, `tap-check.yml`, `harness/README.md`, `CHANGELOG.md`, Benutzerhandbuch, `branch-protection.md`, [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) — „unverändert“, im Diff tatsächlich nicht enthalten (`git diff --stat` über `packaging scripts .github` leer; `CHANGELOG.md`, `docs/user/`, `harness/`, `docs/plan/adr/` nicht in der Liste). Diff-Dateien ohne Plan-Eintrag: keine. Die Plan-Fundstellen zu den Skripten habe ich nachgelesen: `smoke.yml` Zeilen 10 bis 14 (Trigger), `scripts/tap-nachzug.sh` Zeilen 25 bis 27 und 34 (Klon, `mkdir`, `cp`, `add Formula/u-boot.rb`), `scripts/tap-check.sh` Zeile 34 (`git push --dry-run`) — gleich dem Wortlaut in §9. Dass `tap-check.sh` „kein Recht für `.github/workflows/` belegt“, ist eine Folgerung aus dem Lesen, nicht gemessen (dazu müsste ein Push geschrieben werden).

### Risiken aus Plan §6 (Empfehlung für den Ausgang, Entscheidung bei der Closure)

| Risiko | Befund des Verifiers | Empfohlener Ausgang |
|---|---|---|
| 1 Drift zwischen Quelle und Tap | Stand heute gleich (sha, `cmp`); kein Sensor hält das fest; §9 Grenze führt es im Indikativ | weiter offen → Variante A und Grenze (d) |
| 2 Schreibbefehl ungeprobt | Aufruf nicht lesbar, Block nicht der gelaufene (V-2) | weiter offen → Grenze (d), nicht „entfallen“ |
| 3 Rechte-Hinweis über das Gelesene hinaus | Zitat stimmt mit der Seite; kein Satz über ein bestimmtes Token | entfallen |
| 4 Satz über ein Skript ohne Test | Sätze als Lesebefund mit Fundstelle; kein „nur/nie/immer“ als Zusage; Fundstellen gelesen | entfallen (Stand `8cc97257`, ohne Sensor) |
| 5 (Plan: „Der falsche Lauf wird gelesen“) | Schritt 5 erkennt den Lauf an `created_at`; gegen die Lauf-Liste gelesen | entfallen |

### Prüfung der Review-Übergaben

| Finding | Auflösung | Beleg |
|---|---|---|
| F-1 (MEDIUM, Beleg-Zustand ohne Herkunft) | **behoben in `8cc97257`, mit Folge V-1** | Der Absatz nennt Herkunft (Hauptagent, Sitzungsprotokoll), Bezugspunkt (2026-11-06) und die Abweichung von der `B64`-Zeile; zugleich neue Kennung im Text (V-1) |
| F-2 (LOW, `base64 -w0` GNU-Option) | **behoben** | Grenze: „`base64 -w0` … ist eine GNU-Option; ausgeführt wurde sie unter Linux“; die Zeile lief hier unter Linux (3368 Zeichen, Rücklesen gleich) |
| F-3 (LOW, Arbeitsbaum statt Commit) | **behoben** | Voraussetzung vor Schritt 3 (`git fetch`, `git status --short`, `git diff origin/main -- packaging/homebrew-tap/`), alle drei lesend ausgeführt. Restgrenze: der Vergleich deckt nur `packaging/homebrew-tap/`, nicht den Stand von `HEAD` gegen `origin/main` |
| F-4 bis F-6 (INFO) | keine Änderungspflicht | F-4 siehe V-4, F-5 siehe V-5, F-6 (Zeilenangaben): gegen `8cc97257` richtig gelesen |
| F-7 (INFO, Entscheidung ohne Artefakt) | **übernommen** | V-3 |

### Sensors

| Sensor | Ergebnis | Evidence |
|---|---|---|
| `make gates` (Devcontainer-Container, ganzer Lauf, `8cc97257`) | **pass** (Exit 0, 18 s) | `tool-pruefung` OK (go 1.27.3, golangci-lint 2.14.0); `baseline-verify` v6.18.1 OK, 54 Dateien; `d-check`: 1211 Dateien, 0 Befunde; `lint`: 0 issues; `coverage-gate`: 92.40 % ≥ 90 %; `ablage-pruefung-test`: 26 Fälle, 0 Fehlschläge; Stempel `.harness/state/gates-passed.head` = `8cc97257` |
| `make docs-check` für diesen Bericht | **pass** (nach Verlinkung von fünf Kennungen; der erste Lauf meldete 5 `id-unlinked`) | `d-check: 1212 Datei(en) geprüft, 0 Befund(e)` |
| `make test-docker` | **not run** | kein Docker-Pfad berührt (Diff nur Markdown) |
| `make verify-depguard` | **not run** | keine Importregel berührt |

### Traceability

| ID | Beleg |
|---|---|
| [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung) | Der Slice ändert kein Produktverhalten; Beleg ist `releasing.md` §9 als Maintainer-Doku und der Tap-Zustand (Tap-Dateien gleich Quelle). Kein Test pinnt die Doku (Variante A, Plan §1) |
| [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) | Entscheidung 5 (Smoke im Tap) unverändert; §9 beschreibt den Weg der Vorlage in den Tap, ändert die Entscheidung nicht; keine Accepted-ADR berührt |

### Carveouts

- Neu: none
- Gelöst: none
- Unverändert: alle (Diff berührt `docs/plan/carveouts/` nicht)

### Nicht ausgeführt

- Schreibschritt 3 für `README.md`, Dispatch (Schritt 5, `gh workflow run`), Rückweg (Schritt 6) — Auftragsgrenze: nichts im Tap und auf GitHub schreiben; die Befehlsform des `PUT` ist damit weiterhin nur durch das Ergebnis des Tap-Commits `662c5fb5` gestützt, nicht durch den Aufruf (V-2).
- Inhalt des Sitzungsprotokolls des Hauptagenten — nicht gelesen, nicht im Repo.
- `make test-docker`, `make verify-depguard` — siehe Sensors.
- Rechte eines bestimmten Tokens für `.github/workflows/` — nicht gemessen; der Abschnitt behauptet sie nicht.

### Commit / Artefakt

- Geprüfter Stand: `8cc97257` (Implementierung `11fc3b01`, Review-Fix `8cc97257`); Tap-Kopf `662c5fb5c57fc61751f5cff127f213be57573211` (unverändert, gelesen)

## Verdikt

**Slice nicht abschlussreif, ein Befund offen:** V-1 (MEDIUM, wörtliche DoD-Abweichung durch den Review-Fix) geht an den Planner/Implementer; ohne seine Auflösung kann der DoD-Punkt „Verfahren“ nicht abgehakt werden. V-2 und V-3 sind Hinweise für die Closure-Notiz (Ausgang von Risiko 2; Herkunft der Entscheidungen), V-4 bis V-6 sind INFO. Gates, Tap-Gleichheit, Auffindbarkeit, Plan-vs-Diff und die lesenden Befehle von §9 sind belegt; nach der Auflösung von V-1 ist ein Nachlauf-Review des Fix-Commits (Zeilen im Absatz „Gelaufen und nicht gelaufen“) ausreichend, ein neuer Gate-Lauf wegen der Dokumentationsänderung `make docs-check` plus `make gates`.

## Nachtrag

`make docs-check` für diesen Bericht: grün, 1212 Dateien, 0 Befunde (Stand nach dem Schreiben).
