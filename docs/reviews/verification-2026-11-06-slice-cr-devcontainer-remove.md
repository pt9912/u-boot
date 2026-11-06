# Verification-Report: slice-cr-devcontainer-remove — 2026-11-06

**Rolle:** Verifier (frischer Kontext; Baseline-Regelwerk `modul-11`)
**Gegenstand:** [`slice-cr-devcontainer-remove`](../plan/planning/done/slice-cr-devcontainer-remove.md) — Commit-Range `0e8693f8^..45f337a7` (HEAD, 11 Commits, 13 Dateien, 1021 Einfügungen / 45 Löschungen; nur Markdown, kein Code). Abgleichstand **nach der letzten Fix-Runde = `45f337a7`** (= Repo-Tip, Arbeitsbaum vor diesem Report sauber). Lastenheft-Stände in der Range: `0.37.0` (`602d8fa4`), `0.37.1` (`a7914c99`), `0.37.2` (`7f568768`), `0.37.3` (`18f07fcc`); seit `18f07fcc` ist `spec/lastenheft.md` unverändert (`git diff 18f07fcc HEAD -- spec/` leer).

**Eingangs-Kontext:**

- Slice-Plan §1 (Ziel, Entscheidungen 1 bis 5, Auslegungen F1 bis F7, Folgeänderungen F-1 bis F-5, G-1, G-2, H-1), §2 (DoD, 7 Zeilen und 4 Pflichten der Closure), §3 (Plan-Tabelle und Abgleiche), §4, §6 (Risiken), §7 (Closure leer), §8
- [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Änderung), [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) (Aufzählung Zeile 210, Absatz Zeile 214), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Zeile 189), [`LH-SA-CLI-002`](../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) (Zeile 1190), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (nur Anwendung), [`LH-FA-ADD-007`](../../spec/lastenheft.md#lh-fa-add-007--service-entfernen) (Abgrenzung), [`ADR-0015`](../plan/adr/0015-sandbox-volumes-pro-instanz.md)
- CR-Vorschlag [`cr-vorschlag-devcontainer-remove`](../plan/cr/cr-vorschlag-devcontainer-remove.md) (Stellen (a) bis (q)); Welle [`welle-devcontainer-remove`](../plan/planning/welle-devcontainer-remove.md); Folge-Slices [`slice-v2-spec-devcontainer-remove`](../plan/planning/done/slice-v2-spec-devcontainer-remove.md) und [`slice-v1-devcontainer-remove`](../plan/planning/open/slice-v1-devcontainer-remove.md) (nicht umgesetzt, nur Plan-Zeilen im Diff)
- Review-Reports (Eingang, nicht Evidence): [`2026-11-05-slice-cr-devcontainer-remove`](2026-11-05-slice-cr-devcontainer-remove.md) (0 HIGH / 3 MEDIUM / 3 LOW / 3 INFO, merge-blockierend), [`2026-11-05-slice-cr-devcontainer-remove-nachlauf`](2026-11-05-slice-cr-devcontainer-remove-nachlauf.md) (0/2/1/2, merge-blockierend), [`2026-11-06-slice-cr-devcontainer-remove-nachlauf-2`](2026-11-06-slice-cr-devcontainer-remove-nachlauf-2.md) (0/1/1/2, merge-blockierend), [`2026-11-06-slice-cr-devcontainer-remove-nachlauf-3`](2026-11-06-slice-cr-devcontainer-remove-nachlauf-3.md) (0/0/2/2, nicht merge-blockierend)
- Messprotokoll [`slice-v1-devcontainer-remove-messung`](../messungen/slice-v1-devcontainer-remove-messung/README.md) (nur als Grenze des Wortlauts, nicht neu geprüft)
- [`harness/verification.md`](../../harness/verification.md), [`harness/roles.md`](../../harness/roles.md), [`AGENTS.md`](../../AGENTS.md) (§3.1, §3.3, §3.4, §3.5, §3.6, §3.7, §3.10)

**Eingang des Implementers:** Der Auftrag trug keine Sensor-Belege des Implementers. Alle Sensoren dieses Reports sind eigene Läufe des Verifiers auf `45f337a7`.

---

## Verification Evidence

Scope:
- Slice: [`slice-cr-devcontainer-remove`](../plan/planning/done/slice-cr-devcontainer-remove.md)
- IDs: [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-SA-CLI-002`](../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-FA-ADD-007`](../../spec/lastenheft.md#lh-fa-add-007--service-entfernen), [`ADR-0015`](../plan/adr/0015-sandbox-volumes-pro-instanz.md)
- Artefakte: `spec/lastenheft.md` (Versionsfeld Z. 9, Z. 189, 210, 214, 515, 517, 522 bis 524, 1190, Historie Z. 1662 bis 1665), CR-Vorschlag, Slice, Welle, `open/README.md`, `roadmap.md`, beide Folge-Slices, vier Review-Reports

### DoD-Abgleich (Plan §2)

| # | DoD-Punkt | Zustand | Evidence |
|---|---|---|---|
| 1 | **CR-Vorschlag und Freigabe** (`[x]`) | **erbracht, Haken zu Recht gesetzt** | Datei vorhanden, Alt/Neu je Stelle (a) bis (q); Status „freigegeben wie vorgelegt“; F1 bis F7 mit Ausgang (Abschnitt „Bei der Freigabe entschieden“). **Grenze:** die Freigaben (2026-11-05, 2026-11-06) sind im Repo nur als Aussage der Planungsdateien belegt; der Verifier sieht den Sitzungsdialog nicht. |
| 2 | **Übernahme `0.37.0`** (`[x]`) | **erbracht** | Alle Stellen (a) bis (g) stehen im Lastenheft: Einleitungssatz Z. 512, `down`-Satz Z. 515, Spiegelstrich `remove` Z. 517, Fehlerklassen Z. 522, Nicht-Ziele Z. 524, `--quiet`-Liste Z. 189, Aufzählung Z. 210, Absatz Z. 214, Befehlsübersicht Z. 1190, Version, Historie `0.37.0`. Wortvergleich: siehe Lastenheft-Prüfung. |
| 3 | **Folgeänderung `0.37.1`** samt Nachlauf-Review (`[x]`) | **erbracht** | Stellen (h) bis (k) im Lastenheft (Z. 517, 522, 523, Historie Z. 1663); Report [`2026-11-05-slice-cr-devcontainer-remove-nachlauf`](2026-11-05-slice-cr-devcontainer-remove-nachlauf.md) existiert, Commit `db8fa8af`. |
| 4 | **Folgeänderung `0.37.2`** (`[x]`) | **erbracht** | Stellen (l) bis (o): neuer Satz 10 in Z. 517, Klausel in Z. 522, Image-Satz in Z. 523, Historie Z. 1664 (Commit `7f568768`). |
| 5 | **Nachlauf-Review zu `0.37.2`** (`[x]`) | **erbracht** | Report [`2026-11-06-slice-cr-devcontainer-remove-nachlauf-2`](2026-11-06-slice-cr-devcontainer-remove-nachlauf-2.md), Commit `18f07fcc`; sein MEDIUM ist in `0.37.3` behandelt. |
| 6 | **Folgeänderung `0.37.3`** (`[x]`) | **erbracht** | Stellen (p) und (q): Halbsatz in Z. 214, Version Z. 9, Historie Z. 1665 (Commit `18f07fcc`). |
| 7 | **Nachlauf-Review zu `0.37.3`** (`[x]`) | **erbracht** | Report [`2026-11-06-slice-cr-devcontainer-remove-nachlauf-3`](2026-11-06-slice-cr-devcontainer-remove-nachlauf-3.md) (0 HIGH, 0 MEDIUM, 2 LOW, 2 INFO), Commit `d393209c`. Die Zählung im Plan stimmt mit dem Report. Der Wortlaut ändert sich danach nicht (Lastenheft seit `18f07fcc` unverändert). |

**Pflichten der Closure** (zählen nicht mit, Haken stehen aus):

| Pflicht | Zustand | Evidence |
|---|---|---|
| `make gates` grün | **erbracht (eigener Lauf), Haken steht aus** | siehe Sensors |
| Review durchgeführt, Report | **erbracht, Haken steht aus** | vier Reports; die Befunde sind unten gegen den Ist-Stand geprüft. Ob jeder Lauf einen frischen Kontext hatte, belegt kein Artefakt (alle vier nennen dasselbe Modell, siehe V-7). |
| Verification-Report | **mit diesem Report erbracht, Haken steht aus** | dieses Dokument |
| Closure-Notiz mit Steering-Loop-Lerneintrag; Ausgang je Risiko | **offen** | §7 trägt den Platzhalter; die Risiken stehen auf „offen“ (siehe Risiken); die Register-Evidence zur Closure fehlt (siehe Beobachtungs-Register) |

### Prüfung der Review-Behebung gegen den Ist-Stand

| Finding | Behebung trifft? | Beleg |
|---|---|---|
| Review F-1 (MEDIUM, Fehlerklassen von `remove`) | **ja** | Z. 522 ordnet zu: abgelehnte Abfrage `10`, Stoppen oder Entfernen scheitert `12`; die Klauseln tragen Befehls-Subjekte (`up`, `shell`, `remove`). Gegen den Code gelesen: `DevcontainerStartService.Down` (`internal/hexagon/application/devcontainer_start.go`, ab Z. 781) ruft `wired`, `validInstance`, `FindContainer`, `Stop` und liest keine Devcontainer-Konfiguration; die Aussage „`down` braucht keine `devcontainer.json`“ trifft (siehe V-3). |
| Review F-2 (MEDIUM, Kopien der Befehlsliste ohne Träger) | **ja** | Träger stehen in den Plan-Tabellen: [`slice-v2-spec-devcontainer-remove`](../plan/planning/done/slice-v2-spec-devcontainer-remove.md) §3 (`spec/spezifikation.md`, `spec/architecture.md`, `docs/maintainer/cli-json-contract.md`), [`slice-v1-devcontainer-remove`](../plan/planning/open/slice-v1-devcontainer-remove.md) §3 (README-Tabellen, `zusagen-und-messstaende.md`, Benutzerhandbuch, CHANGELOG). Stichprobe der gelesenen Zeilen am Ist-Stand: `spezifikation.md` Z. 619, 631, 866, 870, 875, 876, 938, 948, 1354; `architecture.md` Z. 303, 381, 651, 656; `cli-json-contract.md` Z. 1244, 1356 bis 1358; `README.md` Z. 214, 215, 222; `zusagen-und-messstaende.md` Z. 58 bis 60: jede nennt, was der Plan sagt. |
| Review F-3 (MEDIUM, abgelehnte Abfrage, `--instance` ohne Container) | **ja** | Z. 517 Satz 4 („bevor er etwas stoppt oder entfernt“), Satz 5 („läuft weiter … nichts gestoppt oder entfernt“), Satz 7 („gibt es keinen Container dieser Instanz, … Exit-Code `0`, auch wenn es andere Container des Projekts gibt“). |
| Review F-4 (LOW, Zusage-Menge) | **ja** | Z. 523 „Der Befehl startet, und `u-boot devcontainer remove` entfernt, unter Docker auf einem Linux-Host mit lokaler Engine (`amd64`)“. |
| Review F-5 (LOW, Abgrenzung) | **ja** | Z. 517 letzter Satz nennt Dienst des Compose-Stacks gegenüber Container des Devcontainers und `--purge`/`--volumes`; gelesen gegen [`LH-FA-ADD-007`](../../spec/lastenheft.md#lh-fa-add-007--service-entfernen) (Z. 653 bis 670: „Volumes nur auf explizite Anforderung (`--purge`) löschen“). |
| Review F-6 (LOW, Zustandsmarker) | **ja** | Suche nach „solange die Freigabe“, „Stand vor dem Commit der Übernahme“, „Arbeitsbaum“, „nicht committet“ in den sieben Plan-Dateien: 0 Treffer. |
| Review F-7 (INFO, Spannung Exit-Code `2` und `10`) | **offen, bewusst** | Z. 212 und Z. 214 stehen unverändert; Plan §6 führt die Frage an Architect und Projektinhaber (siehe Risiken). |
| Review F-8 (INFO, Grenzen von „Volumes dieses Containers“) | **übergeben** | [`slice-v1-devcontainer-remove`](../plan/planning/open/slice-v1-devcontainer-remove.md) §2, DoD-Zeile 3. |
| Review F-9 (INFO, DoD) | **ja** | Haken zu Recht gesetzt (siehe DoD-Abgleich). |
| Nachlauf 1 MEDIUM 1 (Zustandsmarker) | **ja** | wie F-6; `git show 7f568768:spec/lastenheft.md` trägt Version `0.37.2`, die Anweisungen nennen Commits statt „HEAD“. |
| Nachlauf 1 MEDIUM 2 (Freigabe vor Containersuche) | **ja, in zwei Schritten** | `0.37.2` (Z. 517 Satz 10, Z. 522) und `0.37.3` (Z. 214). |
| Nachlauf 1 LOW (Image-Satz) | **ja im Wortlaut** | Z. 523 letzter Satz nennt `remove`; der pinnende Test steht aus (siehe V-2). |
| Nachlauf 1 INFO 2 (Terminal für Fall (e)) | **übergeben** | Slice-v1 DoD-Zeile 4, Punkt (2); Fundstelle `devcontainer_shell_test.go` Z. 309 bis 317 am Ist-Stand gelesen: stimmt („covered by the Docker test under `script`“, ob dieser Test existiert, ist nicht gelesen, und der Plan sagt das). |
| Nachlauf 2 MEDIUM (Ausnahme `0.37.2` gegen Z. 214) | **ja** | Z. 214 trägt „bei `u-boot devcontainer remove --volumes` gilt das nur, wenn es einen Container zu entfernen gibt (…)“; Z. 517 Satz 10 und Z. 522 sagen dasselbe. |
| Nachlauf 2 LOW (`open/README.md` trägt Chronik) | **ja** | Zeile 24 nennt Zustand und Zeiger auf den CR-Vorschlag. |
| Nachlauf 3 LOW 1 („nicht versioniert“) | **ja** | Suche nach „nicht versioniert“ in Slice, CR-Vorschlag, Welle, `open/`: 0 Treffer; Abgleichsabsatz in §3 führt den Stand nach `45f337a7`. |
| Nachlauf 3 LOW 2 (Überschrift der DoD-Zeile in Slice-v1) | **ja** | Überschrift „Bedingungen des Lastenhefts `0.37.2` und `0.37.3` … (1) und (2) gehören zu `0.37.2`, (3) zu `0.37.3`“. |

### Lastenheft-Prüfung

| Prüfpunkt | Ergebnis | Beleg |
|---|---|---|
| Wortlaut = CR-Vorschlag, Endzustand der Kette | **deckungsgleich** | Eigener Wortvergleich (einmaliges Skript, nicht im Repo; Link-Form und Fettdruck normiert): Stellen (a), (b), (d2), (e), (f1) bis (f3), (g), (k), (l), (m), (n), (o), (p), (q) stehen als Fließtext im Lastenheft. Die überholten Fassungen (c), (d1), (h), (i), (j) sind ersetzt: (i) mit der Klausel aus (m) und (j) mit dem Satz aus (n) ergibt wörtlich Z. 522 und Z. 523. Kein Lastenheft-Text ohne Entsprechung im CR-Vorschlag (`git diff 602d8fa4^ HEAD -- spec/lastenheft.md`: sechs Hunks, jeder einer oder mehreren Stellen zugeordnet). |
| Version, Datum, Priorität | **stimmen** | Z. 9 `0.37.3`; Historie `0.37.0` bis `0.37.3` mit 2026-11-05 bzw. 2026-11-06; Priorität `V2` in Z. 512 und Z. 1190 (F6). |
| Anforderungsebene (Auftrag: nur Anforderungen) | **gewahrt** | Suche in den hinzugefügten Zeilen nach Labels, Pfaden, Docker-Befehlen, Schlüsselnamen: nur Bedienoberfläche (`--volumes`, `--yes`, `--instance`). |
| Verweis nach unten (`AGENTS.md` §3.4) | **keiner** | Suche in den hinzugefügten Zeilen nach `ADR-`, `slice-`, `welle`, `roadmap`, `docs/`, `SPEC-`, `ARC-`, `.go`, `.md)`: 0 Treffer. |
| Binnenkonsistenz Satz 4 und 5 gegen Satz 10 (Z. 517) | **konsistent, mit Lesehilfe** | Satz 4 und 5 sagen Exit-Code `10` unbedingt, Satz 10 und Z. 214 schränken auf „einen Container zu entfernen“ ein; Z. 522 trägt dieselbe Bedingung. Die Einschränkung steht hinter den Sätzen, die sie begrenzt; der Plan führt die Lesart (CR-Vorschlag Stelle (l)). |
| Fall „nur Instanz-Container, ohne `--instance`“ gegen die Ausnahme | **kein Widerspruch** | Exit-Code `10` kommt dort aus der fehlenden Auswahl (Satz 8), nicht aus der fehlenden Freigabe; Z. 214 und Satz 10 betreffen die Freigabe. |
| `down`, `up --recreate` unverändert | **ja** | Z. 514 und 515 tragen nur den Zusatzsatz „Den Container entfernt `u-boot devcontainer remove`“; die Zusage „stoppt, entfernt weder Container noch Volumes“ steht wörtlich (siehe V-3 zur Fehlerklassen-Klausel). |
| `u-boot remove <service>`, `down --volumes`, Compose | **unverändert** | Der Diff berührt [`LH-FA-ADD-007`](../../spec/lastenheft.md#lh-fa-add-007--service-entfernen) (Z. 653 bis 670) und [`LH-FA-UP-004`](../../spec/lastenheft.md#lh-fa-up-004--umgebung-stoppen) (Z. 710) nicht; [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) wird angewandt, nicht geändert. |
| Weitere Lastenheft-Stellen mit `remove`-relevanter Aussage | **keine übersehen** | Suche nach „entfernt weder“, „räumt“, „bleibt erhalten“, `devcontainer down`, `Workspace-Volume`: nur die bekannten Stellen. |

### Plan-vs-Code (§3 gegen `git diff --stat 0e8693f8^..HEAD`, kumulativ, Stand nach der letzten Fix-Runde)

| Diff-Datei | Plan-Zeile (§3) | Abgleich |
|---|---|---|
| `spec/lastenheft.md` (+15/−9, 6 Hunks) | Z. 97 | deckt; die genannten Zeilen (9, 189, 210, 214, 517, 522, 523, 524, 1190, 1662 bis 1665; gelesen: 210, 212, 668) stimmen mit `45f337a7`. Der Plan nennt die Nicht-Ziele auf „524“ nur in der Übernahme-Notiz (vor der Übernahme 523); kein Befund |
| `docs/plan/cr/cr-vorschlag-devcontainer-remove.md` (neu) | Z. 96 | deckt |
| [`docs/plan/planning/done/slice-cr-devcontainer-remove.md`](../plan/planning/done/slice-cr-devcontainer-remove.md) | Z. 102 | deckt |
| `docs/plan/planning/open/README.md` | Z. 98 | deckt |
| [`docs/plan/planning/open/slice-v1-devcontainer-remove.md`](../plan/planning/open/slice-v1-devcontainer-remove.md) | Z. 100 | deckt |
| [`docs/plan/planning/done/slice-v2-spec-devcontainer-remove.md`](../plan/planning/done/slice-v2-spec-devcontainer-remove.md) | Z. 99 | deckt |
| `docs/plan/planning/welle-devcontainer-remove.md` (**neu** in `a5115ad0`) | Z. 101 („update (Fix-Runde)“) | Zeile vorhanden, Art ungenau (V-1) |
| `docs/plan/planning/in-progress/roadmap.md` (4 Commits; Zeile 45: „Nichts in Arbeit.“ → Verweis auf die Welle) | „Gelesen, nicht berührt“ (Z. 105) | **Abweichung V-1** |
| [`docs/plan/planning/done/slice-v1-devcontainer-remove-messung.md`](../plan/planning/done/slice-v1-devcontainer-remove-messung.md) (`a5115ad0`, Feld **Welle**) | **keine Zeile** | **Abweichung V-1** |
| [`2026-11-05-slice-cr-devcontainer-remove.md`](2026-11-05-slice-cr-devcontainer-remove.md), [`…-nachlauf.md`](2026-11-05-slice-cr-devcontainer-remove-nachlauf.md), [`2026-11-06-…-nachlauf-2.md`](2026-11-06-slice-cr-devcontainer-remove-nachlauf-2.md), [`…-nachlauf-3.md`](2026-11-06-slice-cr-devcontainer-remove-nachlauf-3.md) (vier neue Reports) | keine Zeile | Reports der Reviewer-Rolle, kein Plan-Eintrag (`AGENTS.md` §3.10, Reports ändert keine Fix-Runde); kein Befund |
| [`docs/reviews/verification-2026-11-06-slice-cr-devcontainer-remove.md`](verification-2026-11-06-slice-cr-devcontainer-remove.md) (dieser Report, neu) | keine Zeile | wie die Reports |
| `docs/plan/planning/observations/BEO-DOK/*/evidence/` | Z. 103 („neu bei der Closure“) | `git diff --stat`: ohne Eintrag; steht für die Closure an |
| `docs/plan/adr/`, `docs/plan/carveouts/`, `harness/`, `Makefile`, `repo.mk`, `.golangci.yml`, `tools/`, `.github/`, `internal/`, `cmd/` | — | `git diff --name-only 0e8693f8^..HEAD` über diese Pfade: leer |
| `spec/spezifikation.md`, `spec/architecture.md`, `README*`, `docs/user/`, `CHANGELOG.md`, `docs/maintainer/` | „berührt dieser Slice nicht“ (Kopf, Z. 13) | `git diff --stat`: ohne Eintrag, deckt. Die Fundstellen der dortigen Kopien stehen in den Plan-Tabellen der Folge-Slices (Stichprobe oben bestätigt) |

Die Mechanik-Commits `3e9bcdff` und `b066bde9` (`make slice-mv`, Verweise nachgezogen) liegen vor der Range; sie berühren [`2026-11-05-slice-v1-devcontainer-remove-messung.md`](2026-11-05-slice-v1-devcontainer-remove-messung.md) und [`verification-2026-11-05-slice-v1-devcontainer-remove-messung.md`](verification-2026-11-05-slice-v1-devcontainer-remove-messung.md) (Link-Pfade). Kein Befund. <!-- d-check:ignore (Make-Target, keine Slice-Kennung) -->

### Entscheidungs-Konformität

| Quelle | Ergebnis | Beleg |
|---|---|---|
| [`ADR-0015`](../plan/adr/0015-sandbox-volumes-pro-instanz.md) (Volumes je Instanz) | **konform** | Satz 7 „lässt die übrigen Container des Projekts und ihre Volumes unberührt“ und das Nicht-Ziel zu Volumes nicht mehr vorhandener Container widersprechen der ADR nicht; die ADR ist im Diff nicht berührt. |
| Entscheidungen 1 bis 5, F1 bis F7, F-1 bis F-5, G-1, G-2, H-1 | **im Wortlaut getragen** | Je Entscheidung die Satznummer im CR-Vorschlag („Lesart“); gegen das Lastenheft gelesen: Name (Z. 517 Satz 1), Standard (2, 3), laufender Container (1), Ollama-Volume (3), `--instance` (7, 8), F1/F-3a (4, 5), F2 (9), F-3b (7), G-1 (10), G-2 (Z. 523), H-1 (Z. 214). |
| Accepted-ADR nicht umgeschrieben (`AGENTS.md` §3.5) | **bestätigt** | Kein File unter `docs/plan/adr/` im Diff; `make doc-immutable` grün. |
| Gate nicht gelockert (`AGENTS.md` §3.6) | **bestätigt** | Kein Eintrag in `Makefile`, `repo.mk`, `tools/`, `.golangci.yml`, `.github/`, `docs/plan/carveouts/`. |
| Move-Disziplin (`AGENTS.md` §3.3) | **bestätigt** | `git diff --name-status 0e8693f8^..HEAD`: nur `A` und `M`, keine `R`; der Slice bleibt in `in-progress/`. |
| Kommentar-Selbstcheck (`AGENTS.md` §3.9) | **nicht anwendbar** | Kein Code-, Skript- oder Konfigurationskommentar im Diff. |
| Sprachregel (Nutzer-Doku ohne Kennungen, `AGENTS.md` §5 Regel 6) | **nicht berührt** | Keine Datei unter `README*`, `docs/user/` im Diff. |

Sensors:

| Sensor | Ergebnis | Evidence |
|---|---|---|
| `make gates` | **pass (eigener Lauf)** | Im Devcontainer-Container (`/.dockerenv`, Docker 29.9.3), Exit 0 auf `45f337a7` vor Anlage dieses Reports: `tool-pruefung` OK (go 1.27.3, golangci-lint 2.14.0, make >= 4.3), `gate-reihenfolge` OK, `baseline-verify` v6.18.1 OK (54 Dateien), `docs-check` 0 Befunde (1170 Dateien, `pinsync` 3 Pins), `lint` „0 issues“, `test` grün (alle Pakete; `(cached)`, die Range enthält keinen Code), `coverage-gate` 92,30 % (Schwelle 90 %), `ablage-pruefung` OK, `ablage-pruefung-test` 26 Fälle / 0 Fehlschläge. Stempel `.harness/state/gates-passed.head` = `45f337a7…` = HEAD (Ablage ignoriert, kein Repo-Beleg). |
| `make docs-check` | **pass** | im Lauf von `make gates` 1170 Dateien, 0 Befunde; mit diesem Report erneut gefahren: 1171 Dateien, 0 Befunde |
| `make doc-immutable RANGE="0e8693f8^..HEAD"` | **pass** | `history-range-guard`: Range aufgelöst, 11 Commits, OK; d-check (`vcs`) 1170 Dateien, 0 Befunde. Keine Accepted-ADR im Diff. |
| Wortvergleich CR-Vorschlag gegen Lastenheft | **pass** | siehe Lastenheft-Prüfung |
| Code-Lesen `DevcontainerStartService.Down`, `DownService.Down` | **pass** | belegt die zwei Aussagen des CR-Vorschlags über Code (Gate vor erstem Engine-Zugriff bei `down --volumes`; `down` ohne Konfigurationszugriff) |

Traceability:

| ID | Beleg |
|---|---|
| [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) | Lastenheft Z. 512, 515, 517, 522 bis 524, Historie Z. 1662 bis 1665; CR-Vorschlag (Stellen (a) bis (q)); vier Review-Reports; dieser Report. Code und Tests: keine, sie liegen im Umsetzungs-Slice (dort DoD-Zeilen 1 bis 4; Lücken: V-2). |
| [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) | Z. 210 und Z. 214 (Stellen (f2), (f3), (p)); Nachlauf-Report 3; Pin-Test im Umsetzungs-Slice (DoD-Zeile 4, Punkt (3)). |
| [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-SA-CLI-002`](../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) | Z. 189 und Z. 1190 (Stellen (f1), (e)); Wortvergleich; `make docs-check` (Anker). |
| [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) | Anwendung (`0`, `10`, `11`, `12`); Abschnitt im Diff unberührt. |
| [`LH-FA-ADD-007`](../../spec/lastenheft.md#lh-fa-add-007--service-entfernen) | gelesen (Z. 653 bis 670) als Quelle des Abgrenzungssatzes; unverändert. |
| [`ADR-0015`](../plan/adr/0015-sandbox-volumes-pro-instanz.md) | unberührt (`make doc-immutable`); Vereinbarkeit gelesen. |

Carveouts:
- Neu: none
- Gelöst: none
- Unverändert: [`CO-001`](../plan/carveouts/CO-001-crlf-logs.md) (der Diff berührt `docs/plan/carveouts/` nicht)

### Beobachtungs-Register

Zum Stand `45f337a7` trägt keine der neun im Plan §8 genannten Beobachtungen eine Evidence-Datei zu diesem Slice (`grep -rl devcontainer-remove` über `observations/*/*/evidence`: nur drei Dateien zur Messung). Die Zählung der Dateien je Beobachtung, die der Plan nennt, stimmt am Ist-Stand (39, 18, 43, 5, 1, 5, 10, 3, 2). Der Verifier schreibt nicht ins Register; die Closure-Notiz entscheidet. Befund je Beobachtung:

| Beobachtung | Befund an diesem Slice | Für die Closure |
|---|---|---|
| [`BEO-DOK/vertragsaenderung-nach-review`](../plan/planning/observations/BEO-DOK/vertragsaenderung-nach-review/observation.md) (18 Dateien, offen) | **neuer Beleg (19.).** Der Wortlaut änderte sich dreimal nach einem Review (`0.37.1`, `0.37.2`, `0.37.3`), jedes Mal mit eigenem Nachlauf-Review vor dem Verifier; die ersten beiden fanden MEDIUM, der dritte keines. Wirkung belegt: das MEDIUM des zweiten Nachlaufs (Rang-1-interner Widerspruch) wäre ohne den Nachlauf-Review erst beim Verifier oder beim Umsetzungs-Slice aufgefallen. Die Zeile „offen, dreizehn Belege“ in `state.md` trägt die Zahl 13, die Evidence-Ablage 18 Dateien (Register-Drift, nicht dieser Slice). | Evidence-Datei zu [`slice-cr-devcontainer-remove`](../plan/planning/done/slice-cr-devcontainer-remove.md) anlegen (Klasse: Wortlaut in drei Fix-Runden nach dem Review geändert, je Runde Nachlauf-Review, Verifier las nach); Zählzeile im `state.md` mit der Dateizahl abgleichen. |
| [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (39 Dateien, offen) | **kein neuer Beleg im Wortlaut.** Die Zusage-Menge (Z. 523: Linux-Host `amd64`, lokale Engine) trägt der Docker-Test des Umsetzungs-Slice, nicht die Messung (verschachtelte Engine); der Wortlaut benennt keine Messgrenze. Der CR-Vorschlag führt die Tabelle „Prüfung gegen die Messgrenzen“ je Folgeänderung. | Lese-Schritt ohne Datei; offen, bis der Docker-Test die Zusage auf dem Linux-Host trägt (Plan §6 Risiko 5). |
| [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../plan/planning/observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (43 Dateien, verkörpert in `AGENTS.md` §3.10) | **Befund V-2, noch kein Beleg.** Der Plan sagt in §8, die drei genannten Bedingungen bekämen je einen pinnenden Test; die Testzuordnung im Umsetzungs-Slice deckt eine davon nicht explizit und lässt vier weitere Bedingungen des Wortlauts ohne benannten Fall. | V-2 vor dem Start des Umsetzungs-Slice beheben; sonst Evidence-Datei (Klasse „Bedingung des Wortlauts ohne Testzuordnung“). |
| [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../plan/planning/observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md) (5 Dateien, offen, Schwelle erreicht) | **kein Beleg an diesem Slice.** Im Lastenheft keine vergessene Stelle (eigene Suche, siehe Lastenheft-Prüfung); die Kopien außerhalb tragen Träger (Review F-2, bestätigt). | Lese-Schritt. |
| [`BEO-DOK/vertrag-schweigt-zum-nutzerpfad`](../plan/planning/observations/BEO-DOK/vertrag-schweigt-zum-nutzerpfad/observation.md) (5 Dateien, offen, Schwelle erreicht) | **neuer Beleg (6.).** Drei Nutzerpfade fand erst ein Review (abgelehnte Abfrage bei laufendem Container; `--instance` ohne Container bei vorhandenen anderen; `remove --volumes` ohne Container und ohne `--yes`), alle drei entschied eine Fix-Runde. | Evidence-Datei anlegen (Klasse: Nutzerpfad nach dem Befehl nicht entschieden, drei Fälle). |
| [`BEO-DOK/freigabe-vorlage-verdichtet-den-entwurf`](../plan/planning/observations/BEO-DOK/freigabe-vorlage-verdichtet-den-entwurf/observation.md), [`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](../plan/planning/observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md), [`BEO-APP/container-identitaet-nur-ueber-label`](../plan/planning/observations/BEO-APP/container-identitaet-nur-ueber-label/observation.md), [`BEO-DOK/dod-wortlaut-vom-implementer-geaendert`](../plan/planning/observations/BEO-DOK/dod-wortlaut-vom-implementer-geaendert/observation.md) | **kein belegter Befund.** Die Vorlage zeigt den Wortlaut je Stelle im Volltext; die Messung fuhr `down`/`up`, der Wortlaut trägt für `remove` eine Anforderung (ausgewiesen); die Suche des Containers bleibt der Spezifikation; die DoD-Zeilen 3 bis 7 sind Planner-Zeilen und im Plan §8 ausgewiesen. | Lese-Schritt; der Implementer entscheidet, ob ein Eintrag zählt. |

### Risiken (Plan §6) und was die Closure braucht

| Risiko | Zustand | Für die Closure |
|---|---|---|
| 1 — Wortlaut sagt mehr zu, als die Messung trägt | Ausgang „offen; der Reviewer prüft jeden Satz“ — **überholt**: Review F-4/F-8, die Tabellen „Prüfung gegen die Messgrenzen“ je Folgeänderung und diese Verifikation haben geprüft; Wortlaut nennt keine Mess-, Plattform-, Label- oder Namensgrenze außer der Zusage-Menge. | Schließen mit „geprüft; die Zusage-Menge trägt der Docker-Test des Umsetzungs-Slice“. |
| 2 — Verwechslung mit `u-boot remove <dienst>` | Ausgang „Wortlaut behandelt“ — **bestätigt** (Z. 517 letzter Satz). | Schließen; das Handbuch grenzt im Umsetzungs-Slice ab. |
| 3 — übersehene Stelle der Befehlsübersicht | Ausgang „offen bis zur Closure“ — **geprüft**: im Lastenheft keine vergessene Stelle (Review und eigene Suche), außerhalb Träger vorhanden. | Schließen. |
| 4 — Nicht-Ziel-Satz und Grenze „keine Volumes eines nicht mehr vorhandenen Containers“ | „offen“ — Wortlaut Z. 524 gelesen; die Grenze steht, ein pinnender Test ist nicht zugeordnet (V-2). | V-2 beheben, dann schließen. |
| 5 — Kopien der Befehlsliste | Ausgang „Träger in den Plan-Tabellen eingetragen“ — **bestätigt** (Stichprobe). | Schließen mit Verweis auf die Folge-Slices; bis dahin widerspricht `spec/spezifikation.md` Z. 870 („räumt keine auf“) dem Lastenheft (Rang 1 gilt), siehe V-4. |
| 6 — Zusage und Messgrenze (verschachtelte Engine) | Ausgang „offen bis der Docker-Test … trägt“ | bleibt offen, geht mit dem Risiko an den Umsetzungs-Slice; in der Closure ausdrücklich so benennen. |
| 7 — Änderung nach dem Review | **eingetreten und ausgeführt**: drei Folgeänderungen, drei Nachlauf-Reviews; der Plan nennt jede mit Report. | Schließen mit Verweis auf die Evidence-Datei des Registers. |
| 8 — Spannung Exit-Code `2` (Z. 212) und `10` (Z. 214) | „Frage an den Architect und den Projektinhaber, nicht in diesem Slice gelöst“ | Eine offene Frage ist kein Ausgang im Sinn der Closure-Pflicht („jedes Risiko trägt einen Ausgang“); die Closure nennt „an Architect und Projektinhaber übergeben, Ablage: …“ oder die Entscheidung. Wird sie entschieden, ändert sich Z. 212 oder 214 (neue Wortlaut-Änderung mit Freigabe und Nachlauf-Review). |
| 9 bis 12 — Grenzen von „Volumes dieses Containers“, Reihenfolge Suche und Freigabe, Image-Satz, Terminal für Fall (e) | Ausgänge verweisen auf die Folge-Slices | Als „übergeben“ mit Fundstelle schließen; V-2 beachten. |

Nicht ausgeführt:
- `make test-docker` — kein Docker-Pfad berührt (nur Markdown).
- `make verify-depguard` — keine Importregel berührt.
- `make ci`, `make fullbuild` — kein Security-/Release-Pfad.
- Code-, Replay-, Golden-Evidence nach [`harness/replay.md`](../../harness/replay.md) — kein Generator berührt.
- Prüfung der Freigaben des Projektinhabers (2026-11-05, 2026-11-06) an der Quelle — im Repo nicht sichtbar (siehe DoD-Zeile 1).
- Prüfung, ob die vier Review-Läufe in frischem Kontext liefen — kein Artefakt dazu (V-7).

Commit / Artefakt:
- `45f337a7` (Tip); Lastenheft-Commits `602d8fa4` (`0.37.0`), `a7914c99` (`0.37.1`), `7f568768` (`0.37.2`), `18f07fcc` (`0.37.3`); Reports `7f5110cc`, `db8fa8af`, `18f07fcc`, `d393209c`

## Abweichungen und Hinweise

| ID | Gewicht | Befund | Fundstelle |
|---|---|---|---|
| V-1 | LOW | Plan-Tabelle gegen den kumulativen Diff (`AGENTS.md` §3.10, Abgleich nach der letzten Fix-Runde): (a) `roadmap.md` steht unter „Gelesen, nicht berührt“, der Diff berührt sie in vier Commits (`0e8693f8`, `7e6e6c9f`, `602d8fa4`, `a5115ad0`; Zeile 45 „Nichts in Arbeit.“ → Verweis auf die Welle); (b) [`slice-v1-devcontainer-remove-messung`](../plan/planning/done/slice-v1-devcontainer-remove-messung.md) in `done/` hat keine Zeile (Commit `a5115ad0` setzt das Feld **Welle**); (c) die Zeile zur Welle sagt „update (Fix-Runde)“, die Datei entstand in `a5115ad0` (neu). Die Abgleiche in §3 prüfen je Runde den Teil-Diff, keinen kumulativen; jeder dieser Teil-Abgleiche stimmt. | [`slice-cr-devcontainer-remove`](../plan/planning/done/slice-cr-devcontainer-remove.md) §3, Z. 101 und Z. 105 |
| V-2 | LOW | Testzuordnung der Bedingungen des Wortlauts im Umsetzungs-Slice (`AGENTS.md` §3.10, letzter Punkt). Der Plan §8 und der CR-Vorschlag sagen, „nur auf den Container ohne Instanz“, „ohne `--yes` Exit-Code `10`“ und „ohne Container Exit-Code `0`“ bekämen je einen pinnenden Test. Die DoD-Zeilen 1 und 2 von [`slice-v1-devcontainer-remove`](../plan/planning/open/slice-v1-devcontainer-remove.md) tragen `0` und `10` (Fälle (c), (d), (e), (f)), für „ohne `--instance` nur der Container ohne Instanz, bei nur Instanz-Containern Exit-Code `10` mit Hinweis“ steht kein Fall (die Zeile nennt „Instanz“ ohne Fall; Fall (b) betrifft `--volumes --instance`). Ohne benannten Test bleiben außerdem: das Volume der Ollama-Modelle (Satz 3, nur „Katalog-Features außer `ollama`“ in DoD-Zeile 3), die Nennung der Volumes vor der Abfrage und der entfernten Volumes (Satz 4 und 6), der ungültige Instanzname bei `remove` (Exit-Code `10`, Z. 522), Exit-Code `11` bei nicht erreichbarer oder nicht lokaler Engine für `remove` (Z. 522) und das Nicht-Ziel „Volumes eines nicht mehr vorhandenen Containers entfernt `u-boot` nicht“ (Z. 524). Der Wortlaut ändert sich dadurch nicht; die Zeilen gehören vor den Start des Umsetzungs-Slice in dessen DoD, oder der Plan nennt, welche Bedingung er als Nicht-Zusage führt. | [`slice-v1-devcontainer-remove`](../plan/planning/open/slice-v1-devcontainer-remove.md) §2; [`slice-cr-devcontainer-remove`](../plan/planning/done/slice-cr-devcontainer-remove.md) §8 (Eintrag `spec-bedingung-ohne-code-entsprechung`); `spec/lastenheft.md` Z. 517, 522, 524 |
| V-3 | INFO | Der Plan §1 („Nicht in diesem Slice“) und die Historie sagen, `down` behalte seine Zusage. Die Fassung `0.37.1` ordnet die Klausel „Basisverzeichnis ohne Devcontainer-Konfiguration → Exit-Code `10`“ sowie die Namensklauseln des Containers nur noch `up` und `shell` zu (Z. 522); bis `0.37.0` standen sie ohne Befehl und lasen sich auch für `down`. Gelesen: `DevcontainerStartService.Down` liest keine Konfiguration und prüft keinen Namen (nur `wired`, `validInstance`, `FindContainer`, `Stop`), und [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) führt die Voraussetzungen unter „`up` und `shell`“. Das Lastenheft folgt damit dem Code und der Spezifikation (Rang-1-Text näher an Rang 2 und Code); der CR-Vorschlag weist die Lesart aus (Stelle (i), Punkte (1) bis (4)) und trägt die Frage an den Projektinhaber, die Historie `0.37.1` nennt den Schritt („nennen … den Befehl, für den sie gelten“). Kein Mangel; „behält seine Zusage“ gilt für das Verhalten, nicht für den Wortlaut der Fehlerklasse. | `spec/lastenheft.md` Z. 522, Historie Z. 1663; CR-Vorschlag Stelle (i); `internal/hexagon/application/devcontainer_start.go` `Down` |
| V-4 | INFO | Bis [`slice-v2-spec-devcontainer-remove`](../plan/planning/done/slice-v2-spec-devcontainer-remove.md) widerspricht die Spezifikation (Rang 2) dem Lastenheft (Rang 1): `spec/spezifikation.md` Z. 870 „führt keine Liste der Instanzen und räumt keine auf“ gegen „räumt keine von sich aus auf“, Z. 619, 631, 866, 875, 876, 938, 948 und 1354 nennen `remove` nicht. Bekannt und mit Träger geführt (Plan §6 Risiko 4 und 5, Folge-Slice §3); bis dahin gilt das Lastenheft. | `spec/spezifikation.md` Z. 870 |
| V-5 | INFO | Zustandsmarker in der Roadmap, der nicht zu diesem Slice gehört und neben der von ihm geänderten Zeile steht: `roadmap.md` Z. 43 „… liegt in `in-progress/`“ für [`slice-v1-bau-nur-devcontainer`](../plan/planning/done/slice-v1-bau-nur-devcontainer.md); die Datei liegt in `done/` (`ls`). Der Verweis ist ein leerer Linktext `[]`. Nicht Gegenstand des Diffs; an den Planner. | `docs/plan/planning/in-progress/roadmap.md` Z. 43 |
| V-6 | INFO | Zwei Marker, die bei der Closure auf den Ist-Stand gehen: `open/README.md` Z. 24 und die Welle (§1, §5) nennen „Nachlauf-Review zu `0.37.3` liegt vor, Closure offen“; sie werden mit der Closure auf „geschlossen“ gesetzt (Plan §3 kündigt das für `open/README.md` an). Das Zählwerk „vier Slices“ der Welle ändert sich mit dem Move nach `done/` nicht. | `docs/plan/planning/open/README.md` Z. 24; [`welle-devcontainer-remove`](../plan/planning/welle-devcontainer-remove.md) |
| V-7 | INFO | Unabhängigkeit der Reviews ist nicht belegbar: alle vier Reports nennen dasselbe Modell (`claude-sonnet-5-5`) und den Skill `.harness/skills/reviewer.md`; „kein Self-Review“ (`harness/roles.md`) hängt am frischen Kontext, den die Reports nicht ausweisen. Gleiches gilt für die Freigaben des Projektinhabers (nur als Aussage der Planungsdateien im Repo). Kein Befund gegen den Slice; die Closure-Notiz nennt beide Grenzen. | Reports `2026-11-05-…`, `2026-11-06-…` (Kopfzeilen) |

## Urteil

**Verifiziert mit Auflagen.** Alle sieben inhaltlichen DoD-Zeilen sind erfüllt und am Ist-Stand belegt: [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) trägt den Befehl `u-boot devcontainer remove` in der Fassung `0.37.3`, [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) und [`LH-SA-CLI-002`](../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) nennen ihn; der Lastenheft-Text ist zur Kette der Stellen (a) bis (q) deckungsgleich, enthält nur Anforderungen und verweist nicht nach unten; alle Befunde der vier Reviews (sechs MEDIUM, sieben LOW, neun INFO) sind im Ist-Stand behoben oder mit Ziel übergeben; der Lastenheft-Wortlaut ist seit `18f07fcc` unverändert und hat zu jeder Änderung einen Nachlauf-Review. `make gates` im Devcontainer-Container (eigener Lauf, Exit 0), `make docs-check` (0 Befunde) und `make doc-immutable` sind grün; Gate nicht gelockert, keine ADR berührt, Carveouts unverändert.

Keine Abweichung ist ein Blocker für den Wortlaut; V-1 und V-2 sind Planungs-Hygiene (kein HIGH, kein MEDIUM). Eine Rückführung nach `open/` (§4) ist nicht nötig: die Verifikation verlangt keine Änderung des Wortlauts.

Vor dem Verschieben nach `done/` fehlen (Reihenfolge nach `AGENTS.md` §3.3: Inhalt zuerst, dann der reine `git mv`):

1. V-1 beheben (Plan-Tabelle um `roadmap.md` und die `done/`-Datei ergänzen, Zeile der Welle auf „neu“ berichtigen, „Gelesen, nicht berührt“ für `roadmap.md` entfernen) und den Abgleich kumulativ gegen `git diff --stat` wiederholen, mit diesem Report als neuer Datei unter `docs/reviews/`.
2. V-2: im Umsetzungs-Slice die genannten Bedingungen je einem Fall zuordnen oder als Nicht-Zusage ausweisen; die Formulierung in Plan §8 („je einen pinnenden Test“) auf den Ist-Stand setzen.
3. §6: die Risiken 1 bis 5, 7 und 8 bis 12 mit Ausgang schließen (siehe Tabelle); die Frage zu Exit-Code `2` und `10` mit Ablage versehen.
4. Haken der Closure-Pflichten setzen (`make gates`, Review, Verification-Report), §7 mit Steering-Loop-Lerneintrag schreiben; die Klassen für den Lerneintrag nennt der Abschnitt Beobachtungs-Register.
5. Register-Evidence anlegen: `vertragsaenderung-nach-review` (19. Beleg), `vertrag-schweigt-zum-nutzerpfad` (6. Beleg), gegebenenfalls `spec-bedingung-ohne-code-entsprechung` (wenn V-2 offen bleibt).
6. `open/README.md` Z. 24 und die Welle auf den Ist-Stand setzen (V-6), `make gates` nach der letzten Änderung wiederholen, danach der reine `git mv` (`make slice-mv`). <!-- d-check:ignore (Make-Target, keine Slice-Kennung) -->

## Handoff

```text
Role: Verifier
Input context: Slice-DoD, Plan-Tabelle, Diff 0e8693f8^..45f337a7, Lastenheft 0.37.3, vier Review-Reports, eigene Sensor-Läufe
Changed artefacts: docs/reviews/verification-2026-11-06-slice-cr-devcontainer-remove.md (neu)
Evidence: make gates (Exit 0, Devcontainer-Container), make docs-check (0 Befunde), make doc-immutable (OK), Wortvergleich CR-Vorschlag gegen Lastenheft
Open risks: V-1, V-2 (LOW); V-3 bis V-7 (INFO); Frage Exit-Code 2 gegen 10 bei Architect und Projektinhaber; Freigaben und Reviewer-Kontext nicht belegbar
Next role: Planner (Plan-Tabelle, §6, §7, Register-Evidence, Move), danach Validator der Welle nach dem Umsetzungs-Slice
```
