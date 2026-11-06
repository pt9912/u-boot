# Verification-Report: slice-v2-spec-devcontainer-remove — 2026-11-06

**Rolle:** Verifier (frischer Kontext; Baseline-Regelwerk `modul-11`)
**Gegenstand:** [`slice-v2-spec-devcontainer-remove`](../plan/planning/done/slice-v2-spec-devcontainer-remove.md) — Commit-Range `570c86ce..98a21a37` (HEAD, drei Commits: `638996ca` Start, `f64294a0` Spezifikation, `98a21a37` Fix-Runde nach dem Review; acht Dateien, 340 Einfügungen / 36 Löschungen; nur Markdown, kein Code). Abgleichstand **nach der letzten Fix-Runde = `98a21a37`** (= Repo-Tip, Arbeitsbaum vor diesem Report sauber).

**Eingangs-Kontext:**

- Slice-Plan (Kopf, §1 bis §6, §8), DoD (eine inhaltliche Zeile, vier Pflichten der Closure), Plan-Tabelle §3 (elf Zeilen)
- [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes) (ganz), die Nachzüge in [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes), unter `--quiet`, im Feld `command`, in der Verweisliste und in der Historie; [`spec/architecture.md`](../../spec/architecture.md); [`docs/maintainer/cli-json-contract.md`](../maintainer/cli-json-contract.md) (§6.11, Matrix)
- Lastenheft `0.37.3` gegen [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes): [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Spiegelstriche zu `remove`, Fehlerklassen, Zusage-Menge, Nicht-Ziele), [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)
- Messprotokoll [`slice-v1-devcontainer-remove-messung`](../messungen/slice-v1-devcontainer-remove-messung/README.md) (Kurzfassung, M1, M2a bis M2f, M4; gegen die vier Stellen der Fix-Runde gelesen)
- Folge-Slice [`slice-v1-devcontainer-remove`](../plan/planning/open/slice-v1-devcontainer-remove.md) (Diff, DoD, Plan-Tabelle)
- Review-Report (Eingang, nicht Evidence): [`2026-11-06-slice-v2-spec-devcontainer-remove`](2026-11-06-slice-v2-spec-devcontainer-remove.md) (0 HIGH / 3 MEDIUM / 3 LOW / 4 INFO, merge-blockierend)
- [`harness/verification.md`](../../harness/verification.md), [`harness/roles.md`](../../harness/roles.md), [`AGENTS.md`](../../AGENTS.md) (§3.1, §3.3 bis §3.7, §3.9, §3.10, §5)

**Eingang des Implementers:** Der Auftrag trug keine Sensor-Belege des Implementers. Alle Sensoren dieses Reports sind eigene Läufe des Verifiers auf `98a21a37`.

---

## Verification Evidence

Scope:
- Slice: [`slice-v2-spec-devcontainer-remove`](../plan/planning/done/slice-v2-spec-devcontainer-remove.md)
- IDs: [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`ADR-0015`](../plan/adr/0015-sandbox-volumes-pro-instanz.md); Spezifikation [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes) (neu) und [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (Nachzug)
- Artefakte: `spec/spezifikation.md`, `spec/architecture.md`, `docs/maintainer/cli-json-contract.md`, der Slice, der Folge-Slice, `docs/plan/planning/in-progress/roadmap.md`, `docs/plan/planning/open/README.md`, der Review-Report

### DoD-Abgleich (Plan §2)

| # | DoD-Punkt | Zustand | Evidence |
|---|---|---|---|
| 1 | **Spezifikation** (`[x]`): Aufruf, Suche, Volume-Auswahl, Abfrage, Ausgabe (Text und `--json`), Exit-Codes und Grenzen in `spec/spezifikation.md`; jede Bedingung, die den Code beschreibt, bekommt im Code-Slice einen pinnenden Test oder entfällt; die Kopien der Befehlsliste nennen `remove`; keine Aussage widerspricht dem Lastenheft | **erbracht, Haken zu Recht gesetzt** | Teilpunkte unten |
| 1a | Aufruf, Suche, Volumes, Abfrage, Ausgabe, Exit-Codes, Grenzen | erbracht | [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes): Absätze *Aufruf*, *Ablauf* (neun Schritte), *Auswahl des Containers*, *Menge der Volumes*, *Abfrage und Freigabe* (Tabelle mit sechs Fällen), *Fehler beim Entfernen* (Tabelle der Zustände), *Ausgabe* (Text, `--quiet`, `--json`), *Pfade und Exit-Codes* (15 Zeilen), *Reichweite*. Die Entscheidungen des Slice-Ziels (§1) sind je getragen: Volume eines anderen Containers (Absatz *Menge der Volumes*, zweiter Spiegelstrich), anonyme Volumes (`docker rm -v`), Reihenfolge Abfrage, Stoppen, Entfernen (Ablauf, Schritte 6 bis 8), Exit-Codes `10`/`11`/`12`, Messgrenze (*Reichweite*). |
| 1b | Kopien der Befehlsliste nennen `remove` | erbracht | Am Ist-Stand gelesen: `--quiet`-Liste (Zeile 619) und Absatz (Zeile 631); [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) Einleitung (864), *Aufruf* (866), *Instanzen* (870), *Mit `--instance`* (876), *Container neben dem Container ohne Instanz* (877), Zeiger unter der Tabelle der Exit-Codes (948), *Reichweite* (950: „alle vier Unterbefehle“); Wert `devcontainer` des Feldes `command` (1479); Verweisliste zu [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (1727); `architecture.md` (130, 133, 304 bis 308, 385, 660); `cli-json-contract.md` (§6.11 Überschrift 1244, Absatz 1290, Matrix 1379). |
| 1c | keine Aussage widerspricht dem Lastenheft („räumt keine auf“) | erbracht | Suche `räumt keine` in `spec/`, `docs/maintainer/`: Zeile 870 lautet „räumt keine von sich aus auf“ wie das Lastenheft (Nicht-Ziele), die Historie zitiert die Änderung. Wortvergleich der Lastenheft-Sätze gegen [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes) siehe Lastenheft-Prüfung. |
| 1d | jede Bedingung bekommt im Code-Slice einen Test oder entfällt | erbracht als **Übergabe**, nicht als Test | Tabelle *Belege* (zwölf Zeilen) in [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes); der Träger steht beim Empfänger: DoD-Zeile *Belege der Spezifikation* in [`slice-v1-devcontainer-remove`](../plan/planning/open/slice-v1-devcontainer-remove.md) (je Zeile ein Test der genannten Ebene oder Streichung; Entfernen der Zielbild-Marker), Bezug und Plan-Zeilen für `spec/spezifikation.md`, `spec/architecture.md`, `docs/maintainer/cli-json-contract.md`. Die Tests selbst gehören nach dem Slice-Ziel (§1, „Code“) nicht in diesen Slice. |

**Pflichten der Closure** (zählen nicht mit, Haken stehen aus):

| Pflicht | Zustand | Evidence |
|---|---|---|
| `make gates` grün | **erbracht (eigener Lauf), Haken steht aus** | siehe Sensors |
| Review durchgeführt, Report | **erbracht, Haken steht aus** | [`2026-11-06-slice-v2-spec-devcontainer-remove`](2026-11-06-slice-v2-spec-devcontainer-remove.md); die Befunde sind unten gegen den Ist-Stand geprüft. Ob der Lauf einen frischen Kontext hatte, belegt kein Artefakt (V-4). |
| Verification-Report | **mit diesem Report erbracht, Haken steht aus** | dieses Dokument |
| Closure-Notiz mit Steering-Loop-Lerneintrag; Ausgang je Risiko | **offen** | §7 trägt den Platzhalter „(Wird bei der Closure gefüllt.)“; §6 trägt ein Risiko mit Ausgang „übergeben“ (bestätigt, siehe Risiken). |

### Prüfung der Review-Behebung gegen den Ist-Stand (Fix-Runde `98a21a37`)

| Finding | Behebung trifft? | Beleg |
|---|---|---|
| F-1 (MEDIUM, `std` und `anon` als Beleg der leeren Menge) | **ja** | Zeile 1013: „ohne Katalog-Feature **und ohne `VOLUME` im Dockerfile des Projekts** … (Messung M1, `std`; das Projekt `anon` mit einem von Hand ergänzten `VOLUME` hält ein anonymes)“. Gegen das Messprotokoll gelesen: Tabelle M1, `std` 0 Volumes, `anon` 1 anonymes Volume; Kurzfassung „Das Standardprofil ohne Katalog-Feature hält keins“. Die Beleg-Zeile zur Menge (1089) trägt „leere Menge im Standardprofil ohne `VOLUME` im Dockerfile“. |
| F-2 (MEDIUM, vier Aussagen ohne Beleg) | **ja** | Drei neue Zeilen der Tabelle *Belege*: kein Lesen einer Projektdatei, kein Hilfs-Container, kein Socket, Exit-Code `14` (Zeile 1095, mit Testebene Dienst/Fake-Port, Docker-Adapter-Argumente `run`/`--mount`, CLI-Adapter mit ungültiger `u-boot.yaml`); *Nach dem Entfernen* (1096, Docker-Test mit Lesestand); Inhalt der Meldungen zu `12` und `10` (1097). Jede der vier Aussagen des Findings (a) bis (d) hat eine Zeile. |
| F-3 (MEDIUM, Übergabe ohne Träger im Folge-Slice) | **ja** | Diff des Folge-Slice: Bezug nennt [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes); berührte Spec-Stellen; Nicht-Ziel „Vertrag und Spezifikation“ erlaubt die benannten Änderungen mit Rückführung; neue DoD-Zeile *Belege der Spezifikation*; drei neue Plan-Zeilen. Die Zählung „zwölf Zeilen“ im Folge-Slice stimmt mit der Tabelle (Zeilen 1086 bis 1097 = 12). Slice-§6 sagt jetzt „übergeben“ mit Verweis auf diese DoD-Zeile. |
| F-4 (LOW, `--yes`/`--no-interactive` ohne `--volumes`) | **ja** | Zeile 980: „Ohne `--volumes` hat jede der beiden Optionen für sich keine Wirkung …; zusammen enden sie mit Exit-Code `2`.“ Kein Widerspruch mehr zum Satz davor; die Tabelle der Freigabe (1026) nimmt Aufrufe mit beiden Optionen aus. |
| F-5 (LOW, `docker rm -v` ohne `-f` nicht gemessen) | **ja** | Schritt 8 (991), *Menge der Volumes* (1007: „Messung M2e für `docker rm -f -v`; ob `docker rm -v` ohne `-f` an einem gestoppten Container sie entfernt, belegt der Docker-Test“) und *Reichweite* (1080: als nicht gemessen geführt). Messprotokoll M2e: nur `docker rm -f` und `-f -v`; stimmt. |
| F-6 (LOW, Zielbild-Marker im JSON-Vertrag) | **ja** | `cli-json-contract.md`: Satz vor dem Beispiel (1281), Klammer im Absatz (1290), Matrix-Zeile (1379) tragen „Zielbild“; „der Befehl hat im Code noch keine Entsprechung“ — am Ist-Stand gelesen: unter `internal/adapter/driving/cli/devcontainer*.go` kein Treffer auf einen Befehl `remove` für `devcontainer`. |
| F-7 bis F-10 (INFO) | bewusst offen / zur Kenntnis | F-7 (Linktext Messprotokoll) und F-9 (Fehler-Umschlag ohne `data` nach Teilzustand) bleiben unverändert und sind kein Mangel; F-8 (Image-Satz mit Ersatzbeleg) liegt beim Verifier des Code-Slice (V-3); F-10 (Zeiger in `architecture.md`) ist behoben: Zeile 133 nennt [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes) neben [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes). |

### Lastenheft-Prüfung (Rang 1 gegen Rang 2)

| Prüfpunkt | Ergebnis | Beleg |
|---|---|---|
| Suche (Instanz; ohne Instanz neben Instanzen; nur Instanzen → `10` mit Hinweis; kein Container → `0`) | **deckt** | Lastenheft Zeile 517 Sätze 7 bis 9; [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes) *Auswahl des Containers*, Zeilen der Tabelle *Pfade* |
| Abfrage vor jedem Eingriff, Nennung der Volumes vor der Abfrage und danach | **deckt** | Lastenheft Zeile 517 Satz 4 und 6; *Ablauf* Schritte 5 bis 7, Absatz *Abfrage und Freigabe*, Absatz *Ausgabe* |
| Ohne Freigabe geschieht nichts (laufender Container läuft weiter) | **deckt** | Lastenheft Satz 5; Spec-Tabelle der Freigabe und Absatz danach |
| Ausnahme „nur wenn es einen Container zu entfernen gibt“ (Zeile 214, 517, 522) | **deckt** | Freigabe-Tabelle Zeile 2; Pfade-Tabelle Zeile 4; Auswahl, Spiegelstrich *Kein Container* |
| Fehlerklassen `10`, `11`, `12` | **deckt** | Lastenheft Zeile 522; Pfade-Tabelle (ungültiger Name `10`, Engine `11`, Stoppen/Entfernen `12`); Zusatz der Spec (Volume besteht nach Schritt 9 → `12`) liegt im Lastenheft-Wortlaut „Stoppen oder Entfernen … seiner Volumes scheitert“ |
| Spannung Exit-Code `2` (Lastenheft Zeile 212) und `10` (Zeile 214) | **von der Spec benannt, nicht gelöst** | Absatz nach der Freigabe-Tabelle; das Lastenheft trägt die Entscheidung in der Ausnahme für `remove --volumes`; offene Frage des CR-Vorschlags an Architect und Projektinhaber |
| `--quiet` (Abfrage bleibt, Erfolgsmeldungen entfallen) | **deckt** | [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Aufzählung nennt `devcontainer remove`); Absatz *Ausgabe*, `--quiet` |
| Zusage-Menge und Image-Satz | **deckt, mit Ersatzbeleg** | Lastenheft Zeile 523; *Reichweite*: Linux-Host `amd64`, lokale Engine, „nur, soweit ein Test mit Build-Tag `docker` ihn belegt“; Image-Satz auf den Aufruf ohne `docker` im `PATH` abgebildet. `Dockerfile` Stage `runtime` (`gcr.io/distroless/static-debian12:nonroot`, Zeile 173) enthält kein `docker`; gelesen. Ob der Ersatzbeleg die Bedingung des Folge-Slice erfüllt, entscheidet dessen Verifier (V-3). |
| Nicht-Ziele (keine Liste, Instanzen nicht von sich aus, Volumes nicht mehr vorhandener Container) | **deckt** | Lastenheft Zeile 524; *Fehler beim Entfernen* (letzter Absatz), Zeile 870 |
| Keine Erweiterung des Lastenhefts durch die Spezifikation | **bestätigt** | Alles über das Lastenheft Hinausgehende ist als Festlegung der Spezifikation ausgewiesen (Abgrenzungen der Slice-§1: anonyme Volumes, gehaltenes Volume, Reihenfolge); `git diff 570c86ce..HEAD -- spec/lastenheft.md` leer |

### Plan-vs-Code (Plan §3 gegen `git diff --stat 570c86ce..98a21a37`, Stand nach der letzten Fix-Runde)

| Diff-Datei | Plan-Zeile (§3) | Abgleich |
|---|---|---|
| `spec/spezifikation.md` (+154/−) | Zeilen 1 und 2, 5, 6, 7 der Tabelle | deckt (Einzelstellen gelesen, siehe DoD 1b) |
| `spec/architecture.md` (+/−26) | Zeilen „`spec/architecture.md`“ und „(weitere Stellen)“ | deckt; die Port-Zeilen (`DevcontainerUseCase (Up/Down/Shell)`, `DevcontainerLauncher (Up/Stop/Shell)` Zeile 249, `DevcontainerInspector`) sind laut Plan unverändert und es steht tatsächlich so |
| `docs/maintainer/cli-json-contract.md` (+24) | Zeilen „`cli-json-contract.md`“ und „(weitere Stelle)“ | deckt |
| [`slice-v1-devcontainer-remove`](../plan/planning/open/slice-v1-devcontainer-remove.md) (+10) | Zeile „Träger der Übergabe“ | deckt |
| [`slice-v2-spec-devcontainer-remove`](../plan/planning/done/slice-v2-spec-devcontainer-remove.md) (+15) | Zeile „diese Tabelle, Häkchen …“ | deckt |
| `docs/plan/planning/in-progress/roadmap.md` (−2: Zeile „Nichts in Arbeit.“ entfernt, Commit `638996ca`) | **„gelesen, unverändert“** | **Abweichung V-1** |
| `docs/plan/planning/open/README.md` (Zeile 25 geändert, Commit `638996ca`) | **„gelesen, unverändert“** | **Abweichung V-1** |
| [`2026-11-06-slice-v2-spec-devcontainer-remove.md`](2026-11-06-slice-v2-spec-devcontainer-remove.md) (neu), dieser Report (neu) | keine Zeile | Reports der Reviewer- und Verifier-Rolle, kein Plan-Eintrag (wie im Verification-Report zu [`slice-cr-devcontainer-remove`](../plan/planning/done/slice-cr-devcontainer-remove.md) geführt); kein Befund |
| `docs/plan/adr/`, `docs/plan/carveouts/`, `harness/`, `Makefile`, `repo.mk`, `.golangci.yml`, `tools/`, `.github/`, `internal/`, `cmd/`, `spec/lastenheft.md`, `README*`, `docs/user/`, `CHANGELOG.md` | — | `git diff --name-only 570c86ce..HEAD` über diese Pfade: leer; deckt die Aussage des Slice, der Code und die Nutzer-Doku liegen im Folge-Slice |

### Entscheidungs-Konformität

| Quelle | Ergebnis | Beleg |
|---|---|---|
| [`ADR-0015`](../plan/adr/0015-sandbox-volumes-pro-instanz.md) (Volumes je Instanz) | **konform** | [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes) bestimmt die Menge aus den Einhängungen des Containers, nie aus einem Namensmuster; die ADR ist im Diff nicht berührt. |
| Accepted-ADR nicht umgeschrieben (`AGENTS.md` §3.5) | **bestätigt** | Kein File unter `docs/plan/adr/` im Diff; `make doc-immutable RANGE="570c86ce..HEAD"` grün (3 Commits) |
| Gate nicht gelockert (`AGENTS.md` §3.6) | **bestätigt** | Kein Eintrag in `Makefile`, `repo.mk`, `tools/`, `.golangci.yml`, `.github/`, `docs/plan/carveouts/` |
| Architektur sprach- und meilensteinfrei (`AGENTS.md` §3.4) | **bestätigt** | Hinzugefügte Zeilen von `spec/architecture.md`: keine Treffer auf `ADR-`, `slice`, `welle`, `Meilenstein`; hinzugefügte Zeilen von `spec/spezifikation.md`: nach unten nur der Link auf `docs/messungen/` (Review F-7, Präzedenz [`SPEC-057`](../../spec/spezifikation.md#spec-057--homebrew-formeln-im-devcontainer-katalog-schlüssel-brew-pfade-freigabe-der-taps-konflikt-mit-homebrew-pin-meldungen-hosts-und-grenzen)); keine ADR, kein Slice, keine Welle |
| Move-Disziplin (`AGENTS.md` §3.3) | **bestätigt** | `git diff --name-status 570c86ce..HEAD`: nur `A` und `M`, keine `R`; der Slice bleibt in `in-progress/` |
| Kommentar-Selbstcheck (`AGENTS.md` §3.9), Suppression-Verbot (§3.2) | **nicht anwendbar / bestätigt** | Kein Code-, Skript- oder Konfigurationskommentar im Diff; kein `//nolint` |
| Nutzer-Doku ohne Kennungen (`AGENTS.md` §5 Regel 6) | **nicht berührt** | Keine Datei unter `README*`, `docs/user/` im Diff |
| Sprachvertrag ([`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) Anwendung; Texte der CLI Englisch) | **bestätigt** | Abfrage und Ausgabezeilen in [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes) sind Englisch, die Prosa Deutsch |
| §3.10, Aussagen über Bestand | **bestätigt** | Gelesen: `internal/hexagon/application/downservice.go` (`Down` ruft `runConfirmationGate` in Zeile 63 vor `projectVolumeSnapshot` und `ComposeDown`: Freigabe vor dem Engine-Zugriff bei `down --volumes`), `internal/adapter/driven/confirm/confirm.go` (`y`/`yes` ohne Beachtung der Schreibung, Ende der Eingabe als `N`, anderer Lesefehler als Fehler), `internal/hexagon/port/driving/removeservice.go` (`ErrConfirmerUnavailable`), `Dockerfile` Zeile 173. Die Zielbild-Marker sind wahr: kein `remove`-Befehl im CLI-Adapter. |

Sensors:

| Sensor | Ergebnis | Evidence |
|---|---|---|
| `make gates` | **pass (eigener Lauf)** | Im Devcontainer-Container (`/.dockerenv` vorhanden), Exit 0 auf `98a21a37` vor Anlage dieses Reports: `tool-pruefung` OK (go 1.27.3, golangci-lint 2.14.0, make >= 4.3), `gate-reihenfolge` OK (2 Fälle, 0 Fehlschläge), `baseline-verify` v6.18.1 OK (54 Dateien), `docs-check` 0 Befunde (1176 Dateien, `pinsync` 3 Pins), `lint` „0 issues“, `test` grün (alle Pakete; `(cached)`, die Range enthält keinen Code), `coverage-gate` 92,30 % (Schwelle 90 %), `ablage-pruefung` OK, `ablage-pruefung-test` 26 Fälle / 0 Fehlschläge. |
| `make docs-check` | **pass** | im Lauf von `make gates`: 1176 Dateien, 0 Befunde; mit diesem Report erneut gefahren: siehe Abschnitt „Docs-check mit diesem Report“ |
| `make doc-immutable RANGE="570c86ce..HEAD"` | **pass** | `history-range-guard`: Range aufgelöst, 3 Commits, OK; d-check (`vcs`) 1176 Dateien, 0 Befunde. Keine Accepted-ADR im Diff. |
| Wortvergleich Lastenheft gegen [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes) | **pass (gelesen)** | siehe Lastenheft-Prüfung (Sätze 1 bis 10 der Zeile 517, Zeilen 214, 522 bis 524) |
| Messprotokoll gegen die Fundstellen `Messung M…` der vier Fix-Stellen (M1, M2a bis M2f, M4) | **pass (gelesen)** | siehe Prüfung der Review-Behebung, F-1 und F-5 |
| Zählung der Beleg-Zeilen (`zwölf`) gegen die Tabelle | **pass** | Zeilen 1086 bis 1097 von `spec/spezifikation.md` |

Traceability:

| ID | Beleg |
|---|---|
| [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) | [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes) (Verfeinerung), Verweisliste Zeile 1727, Review-Report, dieser Report. Code und Tests: keine; die Tabelle *Belege* ist die Arbeitsliste des Folge-Slice (dort DoD-Zeile *Belege der Spezifikation*). |
| [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) | Tabelle *Pfade und Exit-Codes* in [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes); Beleg-Zeile „Exit-Codes jeder Zeile … in beiden Klassifikatoren“ (Dual-Classifier-Regel); `architecture.md` Zeile 660 (Fehlerklasse `12`). |
| [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | Tabelle der Freigabe und Absatz *Ausgabe* (`--quiet`); Beleg-Zeilen 2 und 6 der Tabelle *Belege*. |
| [`ADR-0015`](../plan/adr/0015-sandbox-volumes-pro-instanz.md) | unberührt (`make doc-immutable`); Vereinbarkeit gelesen. |

Carveouts:
- Neu: none
- Gelöst: none
- Unverändert: [`CO-001`](../plan/carveouts/CO-001-crlf-logs.md) (der Diff berührt `docs/plan/carveouts/` nicht)

### Beobachtungs-Register

Der Verifier schreibt nicht ins Register; die Closure-Notiz entscheidet (der Plan §8 trägt noch kein Ergebnis der Sichtung). Befund an diesem Slice:

| Beobachtung | Befund | Für die Closure |
|---|---|---|
| [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../plan/planning/observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (verkörpert in `AGENTS.md` §3.10) | **Review F-2 ist ein weiterer Rückfall der Klasse** (vier Aussagen der Spezifikation ohne Beleg-Zeile, vor dem Review nicht gefangen); behoben durch drei Beleg-Zeilen. Dass die Spezifikation selbst eine vollständige Beleg-Tabelle führt, ist die Form, die der Lese-Schritt zu `welle-devcontainer-start-ohne-editor` als Vorschlag (i) vorlegt. | Evidence-Datei zu diesem Slice anlegen (Klasse: Bedingungen im Spec-Text ohne Zeile der eigenen Beleg-Tabelle), Zählung nachziehen. |
| [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md) | **Review F-1 und F-5 sind Belege der Klasse** (Messung auf andere Variante übertragen: `anon` als Standardprofil ohne Volume; `docker rm -f -v` auf `docker rm -v` ohne `-f`); beide behoben, die Spezifikation führt die nicht gemessene Form jetzt in *Reichweite*. | Evidence-Datei oder Zählung nachziehen; die Klasse tritt hier in einer Spezifikation auf, nicht im Lastenheft. |
| [`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md) | **V-1 ist ein Beleg der Klasse** (Plan-Tabelle führt zwei Dateien als „gelesen, unverändert“, die der Diff derselben Slice-Arbeit berührt). | Nach der Berichtigung (V-1) Evidence-Datei, falls die Closure die Klasse zählt. |

### Risiken (Plan §6) und was die Closure braucht

| Risiko | Zustand | Für die Closure |
|---|---|---|
| Die Spezifikation beschreibt Verhalten, das der Code noch nicht hat (Zielbild) — Ausgang „übergeben; DoD-Zeile *Belege der Spezifikation* verlangt je Zeile der Tabelle einen Test und das Entfernen der Zielbild-Marker“ | **bestätigt**: Träger beim Empfänger vorhanden (F-3 behoben), Marker stehen in [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes) (Einleitung), im JSON-Vertrag (drei Stellen) und im Folge-Slice (DoD, Plan-Zeile). | Schließen mit „übergeben, Träger: [`slice-v1-devcontainer-remove`](../plan/planning/open/slice-v1-devcontainer-remove.md) DoD *Belege der Spezifikation*“. Das Zielbild bleibt in der Spezifikation stehen, bis der Folge-Slice schließt; die Marker laufen mit dem Folge-Slice ab (V-2). |

Nicht ausgeführt:
- `make test-docker` — kein Docker-Pfad berührt (nur Markdown); die Belege des Zielbilds liefert der Folge-Slice.
- `make verify-depguard` — keine Importregel berührt.
- `make ci`, `make fullbuild` — kein Security-/Release-Pfad.
- Replay-/Golden-Evidence nach [`harness/replay.md`](../../harness/replay.md) — kein Generator berührt.
- Prüfung, ob Review und Verifikation in frischem Kontext liefen — kein Artefakt dazu (V-4).
- Wiederholung der Messungen M1 bis M5 — nicht Aufgabe des Verifiers dieses Slice; die Fundstellen sind gegen das Protokoll gelesen, nicht gegen die Rohdaten.

Commit / Artefakt:
- `98a21a37` (Tip); Slice-Commits `638996ca` (Start), `f64294a0` (Spezifikation), `98a21a37` (Fix-Runde mit Review-Report)

## Abweichungen und Hinweise

| ID | Gewicht | Befund | Fundstelle |
|---|---|---|---|
| V-1 | LOW | Plan-Tabelle gegen den kumulativen Diff (`AGENTS.md` §3.10, Abgleich nach der letzten Fix-Runde): Die Zeile zu `docs/plan/planning/open/README.md` und `docs/plan/planning/in-progress/roadmap.md` führt beide als „gelesen, unverändert“, der Start-Commit `638996ca` dieses Slice ändert beide (Roadmap: Zeile „Nichts in Arbeit.“ entfernt; `open/README.md`: Zeile 25 trägt „liegt in `in-progress/`, in Umsetzung“). Der Inhalt der Zeile stimmt am Ist-Stand (Zeile 25, Zeile 45 der Roadmap), nur die Art ist falsch. Der Review prüfte den Teil-Diff `f64294a0` und fand sie deshalb nicht. Dieselbe Klasse wie V-1 im Verification-Report zu [`slice-cr-devcontainer-remove`](../plan/planning/done/slice-cr-devcontainer-remove.md). Ein Zusatz: die Stand-Zelle in `open/README.md` Zeile 25 führt „Backlog“ und „liegt in `in-progress/`“ nebeneinander; die Closure setzt sie ohnehin um. | [`slice-v2-spec-devcontainer-remove`](../plan/planning/done/slice-v2-spec-devcontainer-remove.md) §3, Zeile „`open/README.md`, `roadmap.md`“ |
| V-2 | INFO | Zielbild-Marker mit Ablaufpunkt: Einleitung von [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes) („beschreibt das Zielbild“), drei Stellen im JSON-Vertrag und die Historie-Zeile vom 2026-11-06 („als Zielbild festgelegt“). Sie laufen mit dem Folge-Slice ab (dessen DoD und Plan-Zeilen nennen Einleitung, Beispiel, Absatz und Matrix; die Historie-Zeile ist eine Chronik-Zeile und bleibt). Die Zeile 25 in `open/README.md` und die Welle zeigen auf den Folge-Slice. Kein Mangel; der Marker ist wahr, solange `remove` im Code fehlt. | `spec/spezifikation.md` Zeile 978 und 1968; `docs/maintainer/cli-json-contract.md` Zeilen 1281, 1290, 1379 |
| V-3 | INFO | Image-Satz des Lastenhefts (Zeile 523) mit Ersatzbeleg „Aufruf ohne `docker` im `PATH`“: Die Spec weist das ausdrücklich aus und führt die Beleg-Zeile; ob ein Test mit `PATH` ohne `docker` den Satz „läuft dort nicht“ pinnt oder die Nicht-Zusage im Handbuch genügt, entscheidet der Verifier des Folge-Slice (dessen DoD-Zeile (1)). Review F-8, hier nicht entschieden. | `spec/spezifikation.md` Zeile 1080 und 1094; [`slice-v1-devcontainer-remove`](../plan/planning/open/slice-v1-devcontainer-remove.md) §2 |
| V-4 | INFO | Unabhängigkeit des Reviews nicht belegbar: der Report nennt dasselbe Modell (`claude-sonnet-5-5`) und den Skill `.harness/skills/reviewer.md`; „kein Self-Review“ (`harness/roles.md`) hängt am frischen Kontext, den der Report nicht ausweist. Kein Befund gegen den Slice; die Closure-Notiz nennt die Grenze. | Kopfzeile des Review-Reports |
| V-5 | INFO | Offene Fragen an Architect und Projektinhaber, nicht an diesen Slice: (1) Spannung Exit-Code `2` gegen `10` bei `--no-interactive` (Lastenheft Zeilen 212 und 214), von der Spec benannt, nicht gelöst; (2) Fehler-Umschlag ohne `data` nach einem Teilzustand (Container entfernt, Volumes teilweise): der Zustand steht nur im Meldungstext (Review F-9). Die Closure führt beide mit Ablage, nicht als Ausgang eines Risikos. | `spec/spezifikation.md` Absatz nach der Freigabe-Tabelle, Absatz *Ausgabe* (`--json`) |
| V-6 | INFO | Zeilenangaben der Plan-Tabelle (`Zeilen 299 bis 314`, `Zeile 619`, `Zeile 656`, `Zeile 866` …) gelten für den Stand `7f5110cc` vor der Änderung; die Tabelle sagt es nur in der ersten Zeile. Nach dem Diff sind sie verschoben (Beispiel: der Absatz zum CLI-Adapter in `architecture.md` liegt jetzt bei 304 bis 308). Kein Mangel der Aussage, ein Leser der Closure-Notiz sollte sie als Stand vor der Änderung lesen. | [`slice-v2-spec-devcontainer-remove`](../plan/planning/done/slice-v2-spec-devcontainer-remove.md) §3 |

## Urteil

**Verifiziert mit Auflagen.** Die einzige inhaltliche DoD-Zeile ist erfüllt und am Ist-Stand belegt: [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes) legt Aufruf, Suche, Volume-Menge, Abfrage, Ausgabe, Exit-Codes und Grenzen fest, jede Kopie der Befehlsliste nennt `remove`, keine Aussage widerspricht dem Lastenheft `0.37.3`, und die Pflicht zum Beleg je Aussage hat beim Empfänger einen Träger. Alle drei MEDIUM- und drei LOW-Findings des Reviews sind im Ist-Stand behoben (gegen Messprotokoll und Text gelesen). `make gates` im Devcontainer-Container (eigener Lauf, Exit 0), `make docs-check` (0 Befunde) und `make doc-immutable` sind grün; Gate nicht gelockert, keine ADR berührt, Carveouts unverändert, kein Code im Diff.

Keine Abweichung ist ein Blocker für die Spezifikation; V-1 ist Planungs-Hygiene (LOW). Eine Rückführung nach `next/` (§4) ist nicht nötig: der CR-Wortlaut trägt die Spezifikation.

Vor dem Verschieben nach `done/` fehlen (Reihenfolge nach `AGENTS.md` §3.3: Inhalt zuerst, dann der reine `git mv`):

1. V-1 beheben: Plan-Zeile zu `open/README.md` und `roadmap.md` auf „update (Start-Commit)“ setzen und den Abgleich der Plan-Tabelle kumulativ gegen `git diff --stat` wiederholen.
2. Haken der drei Pflichten der Closure setzen (`make gates`, Review, Verification-Report) und §7 mit Steering-Loop-Lerneintrag schreiben; die Klassen nennt der Abschnitt Beobachtungs-Register (Aussage ohne Beleg-Zeile; Messung auf andere Form übertragen; Plan-Tabelle gegen kumulativen Diff). §6 mit „übergeben, Träger: Folge-Slice“ schließen; V-5 mit Ablage führen.
3. Stand-Zellen der Closure: `open/README.md` Zeile 25 („Backlog“ entfällt) und die Zeile der Roadmap/Welle auf den Ist-Stand; `make gates` nach der letzten Änderung wiederholen, danach der reine `git mv`.

## Docs-check mit diesem Report

`make docs-check` mit diesem Report: 1177 Dateien, 0 Befunde (`pinsync` 3 Pins). Der erste Lauf mit dem Report fand 2 unverlinkte Kennungen im Report selbst (Pfade in der Tabelle Plan-vs-Code); berichtigt. Der Lauf vor Anlage: 1176 Dateien, 0 Befunde.

## Handoff

```text
Role: Verifier
Input context: Slice-DoD, Plan-Tabelle, Diff 570c86ce..98a21a37, Lastenheft 0.37.3, SPEC-058, Review-Report, eigene Sensor-Läufe
Changed artefacts: docs/reviews/verification-2026-11-06-slice-v2-spec-devcontainer-remove.md (neu)
Evidence: make gates (Exit 0, Devcontainer-Container), make docs-check (0 Befunde), make doc-immutable (OK), Lesen von Lastenheft, Messprotokoll, Code-Fundstellen
Open risks: V-1 (LOW); V-2 bis V-6 (INFO); Spannung Exit-Code 2 gegen 10 und Fehler-Umschlag ohne data bei Architect und Projektinhaber; Review-Kontext nicht belegbar
Next role: Planner (Plan-Tabelle, Haken, Closure-Notiz, Stand-Zellen, Move), danach Implementer des Folge-Slice slice-v1-devcontainer-remove
```
