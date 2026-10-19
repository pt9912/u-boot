# Verification-Report: slice-v2-spec-zielbild — 2026-10-19

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“, gegen Plan, DoD und Entscheidungen). Kein Review, keine Validierung, keine Reparatur.

**Gegenstand:** Slice `slice-v2-spec-zielbild` (`docs/plan/planning/in-progress/`), Commit-Range `6af82fd..HEAD` (7 Commits, HEAD `5c19443`, 5 Dateien: `CHANGELOG.md`, `docs/user/benutzerhandbuch.md`, `spec/spezifikation.md`, die Slice-Datei, der Review-Report). Inhalts-Commits: `4080858` (Grenze zu `nestedRuntime` in Handbuch und CHANGELOG), `bb9863a` (Spezifikation auf Zielbild), `9f748b6` (Historienzeile verlinkt), `963415a` (Entscheidungen und Abgleich im Slice), `969f2a5` und `a659f1f` (Review-Report), `5c19443` (Fix-Runde). Kein Produktcode, kein Test, kein Lastenheft, keine Architektur-Sicht, keine ADR im Diff. <!-- d-check:ignore (Kennung als Zitat) -->

**Eingang:** Review `docs/reviews/2026-10-19-slice-v2-spec-zielbild.md` (kein HIGH/MEDIUM; F-1, F-2 LOW; F-3 bis F-6 INFO), Slice §1 (Fundstellen-Tabelle), §2 (DoD), §3 (Plan und Abschnitt *Umsetzung*), §6 (Risiken); Anweisungen des Projektinhabers vom 2026-10-19 („specs sind Zielbild“, „ins Handbuch kommt der Ist-Zustand“). <!-- d-check:ignore (Kennung als Zitat) -->

**Gesamturteil: GO mit Auflagen.** Die Spezifikation und die Architektur-Sicht tragen nach eigener Suche keinen Messstand, keinen Ist-Zustand ohne Zusage und kein Chronik-Wort in Prosa mehr; jede entfernte Klasse-A-Aussage steht im Handbuch im Ist-Zustand (Stelle für Stelle gegengelesen); die Grenze zu `nestedRuntime` ist vor dem Streichen übertragen und stimmt mit dem Code überein; Lastenheft, ADRs, Code, Tests und Reports sind unverändert; alle aufgetragenen Sensoren sind grün. Eine DoD-Verletzung liegt nicht vor. Zwei Restpunkte der Spezifikation (V-1 LOW, V-2 INFO) sind kein Hinderungsgrund. Die Auflagen sind Closure-Arbeit des Planners (A-1 bis A-6) und zwei Punkte für den Projektinhaber (Abschnitt 9).

## 1. Sensoren (selbst gefahren)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | pass, **Exit 0** | `EXIT=0` in der Log-Datei; `coverage-gate: OK — coverage 91.40% meets threshold 90%`; HEAD `5c19443`; Arbeitsbaum vorher sauber |
| `make doc-trace` | pass, Exit 0 | „83 Anforderung(en), 0 Waise(n).“ |
| `make docs-check` | pass, Exit 0 | 478 Dateien, 0 Befunde (vor Anlage dieses Reports) |
| `make doc-planning` | pass, Exit 0 | 478 Dateien, 0 Befunde |
| `make doc-immutable RANGE="6af82fd..HEAD"` | pass, Exit 0 | 0 Befunde (Accepted-ADRs unverändert) |
| `git diff --stat 6af82fd..HEAD -- internal cmd spec/lastenheft.md spec/architecture.md docs/plan/adr docs/maintainer README.md README.de.md Dockerfile Makefile` | leer | kein Produktcode, kein Lastenheft, keine Architektur-Sicht, keine ADR, keine Build-Verträge berührt |
| `make test-docker` | nicht gelaufen | kein Code und kein Test im Diff; die Pins wurden gelesen (Abschnitt 4), nicht gefahren. Die E2E-Pins laufen nur lokal mit Docker-Engine und Netz (Register `BEO-E2E/ci-ueberspringt-nested-podman-build`, Slice §6 Risiko 6); die Spezifikation behauptet nichts über CI |
| `make verify-depguard` | nicht gelaufen | keine Importregel im Diff |

## 2. Prüfpunkt 1: Spezifikation und Architektur-Sicht ohne Messstand

Eigene Mustersuche über die **ganze** Datei `spec/spezifikation.md` (nicht nur die Plan-Stellen) und `spec/architecture.md`: `gemessen`, `messstand`, `messung`, `colima`, `arm64`, `amd64`, `ein lauf`, `nicht verifiziert`, `nicht gemessen`, `nicht geprüft`, `stand 20`, `ungemessen`, `stabilität`, `wie bisher`, `früher`, `nicht mehr`, `weiterhin`, `jetzt`, `bisher`, `derzeit`, `noch nicht`, `nicht ausgeschlossen`, `docker 29`, `docker desktop`, `aktuell`; zusätzlich Versionsnummern (`0.89`, `2.91`, `0.40`, `29.5`, `0.10.3`), `nicht zugesagt`, `analog`, Datumsangaben. Jeden Treffer im Absatz gelesen.

| Treffer | Einordnung |
|---|---|
| `LH-FA-DEV-007.b`: `amd64` und `arm64` zu Podman 6.1.4; `SPEC-047`: `OLLAMA_SHA256_AMD64`/`_ARM64` | Architekturangabe des gepinnten Archivs, Pin-Zusage; kein Messstand | <!-- d-check:ignore (Kennung als Zitat) -->
| `SPEC-047` Einleitung: „Der Messstand der Zusagen steht im Lastenheft“ | Verweis aufwärts, Klasse D, vom Plan ausgenommen; kein Messstand in der Spezifikation (INFO V-2) | <!-- d-check:ignore (Kennung als Zitat) -->
| `SPEC-047` *Volume* „nicht zugesagt“, *Kein Start, nur Loopback* „nicht zugesagt“, *Clone-Quelle* „nicht zugesagt“ | Grenzen der Zusage (Klasse D) | <!-- d-check:ignore (Kennung als Zitat) -->
| §7 Historie, fünf Zeilen vom 2026-10-19 | Chronik der Änderung; vier Zeilen im Wortlaut unverändert (Diff), die neue Zeile nennt keinen Messstand, „nicht mehr“ steht dort als Chronik der Änderung |
| Zeile 216 „weiterhin registriert“, 291 „aktuellen Spezifikation“, 428/435 „noch nicht stabilisiert“, 653 „noch nie gestartet“, „früheres Script“ | Zustands- bzw. Verhaltensaussagen ohne Chronik- oder Messbezug; nicht im Slice, nicht beanstandet |
| `SPEC-016`: „analog Referenzprojekt“, `SPEC-006`/`SPEC-016` ohne `Stand 2026-10-14` | Pin-Stand entfernt (Diff); der Wert steht im `Dockerfile` (`ARG GO_VERSION=1.27.2`, gelesen) | <!-- d-check:ignore (Kennung als Zitat) -->
| `wie bisher` | **kein Treffer** in `spec/spezifikation.md` und `spec/lastenheft.md` (Review F-2 behoben, Diff bestätigt) |
| `spec/architecture.md` | Treffer nur `Coverage-Messung` (Prozess, kein Messstand), „aktuelle Liste steht im Code“ und die Schlusszeile „Geplante Erweiterungen … keine“; Datei im Diff unverändert |

**Abgrenzung Messstand (raus) gegen Zusage-Grenze (bleibt):** eingehalten. Geblieben sind nur Grenzen, die das Lastenheft selbst trägt (GPU, Dienst des Hosts, Modellverwaltung, Überstehen des Neubaus, Wert von `OLLAMA_HOST`, abgeleiteter Host eines lokalen Pfads). Entfernt sind Host, Version, Anzahl der Läufe und „nicht gemessen“-Sätze. Ergebnis: **kein Messstand, kein Ist-Zustand ohne Zusage, kein Chronik-Wort in Prosa**.

## 3. Prüfpunkt 2: Nichts geht für Nutzer verloren

Jede entfernte Klasse-A-Aussage gegen das Handbuch (`docs/user/benutzerhandbuch.md`) gelesen, mit Zeile im Stand `5c19443`:

| Slice-Zeile | Entfernte Aussage | Handbuch (Ist-Zustand) |
|---|---|---|
| 1 | Messstand Blob-Host | Zeile 1589 (Punkt *Egress*, „Nicht gemessen: ob der Name des zweiten Hosts …“), Zeile 1243 (Egress-Tabelle) |
| 4, 5 | `dnsmasq` 2.91, Docker 29.5.2/Colima, `REFUSED`, gleiche IPs, `www.example.org` nicht gemessen, Adresse „nicht geprüft und nicht ausgeschlossen“ | Zeilen 1309 bis 1315, *Messstand* 1324 |
| 6 | offener Upstream-Resolver, gemessen in einem Lauf | Absatz *Grenze der Namensauflösung* (Zeile 1220) mit Messstand |
| 7 | `CNAME`-Ziele | Punkt *Namen hinter einem CNAME* (Zeilen 1317 bis 1323) |
| 9 | „das setzt das Devcontainer-Werkzeug“ | Zeile 1134 (Dev-Container-Werkzeug muss den Container mit dem Projektordner erzeugt haben) |
| 10 | Prozess ohne `sudo` ersetzt die Statusdatei, Lauf auf Colima/arm64 | Zeile 1160 |
| 11 | gestoppter Container behält die Liste, Dev-Container-CLI 0.89.0, nicht gemessen: Docker Desktop, Linux ohne VM, `amd64`, Editor | Zeilen 1332 bis 1343 |
| 16 | `amd64` nicht verifiziert | Zeile 1595 (*Grenzen der Messung*) |
| 17 | Lauf von `ollama` mit `nestedRuntime` nicht gemessen | **neu** in Handbuch 1597 bis 1599 und CHANGELOG `[Unreleased]` (Commit `4080858` **vor** dem Streichen in `bb9863a`) |
| 19 | Kurznamen/Images gemessen, `vfs` mit Volume gemessen | Zeilen 1050, 1052, 1061, 1066 (*Images bauen*, Absatz **Gemessen**) |
| 20 | `Stand 2026-10-14: 1.27.2` | `Dockerfile` Zeile 29 `ARG GO_VERSION=1.27.2` |

Handbuch-Diff (ein Aufzählungspunkt) enthält keine `LH-*`-, `SPEC-*`-, `ADR-*`-Kennung und keinen Verweis auf Planung oder Spec (`AGENTS.md` §5 Regel 6); eine Suche über das ganze Handbuch ohne Treffer. CHANGELOG-Diff: ein Halbsatz im Eintrag `ollama`, Englisch wie der Rest, gleiche Grenze wie das Handbuch.

**Die Grenze zu `nestedRuntime` gegen den Code:** `devcontainer_sandbox.go` (Zeilen 98 bis 117) hängt das Workspace-Volume, dann (`nestedRuntime`) `-containers` bzw. `-docker`, dann `-ollama` an, in dieser Reihenfolge. `generate_ollama_test.go` (Zeilen 353 und 355, Fälle „podman and ollama“, „docker and ollama“) pinnt die Reihenfolge byte-genau. Handbuch („die Reihenfolge der Volumes wird erzeugt, der Lauf ist nicht gemessen“), CHANGELOG („the order of the volumes is generated“) und `SPEC-047` *Volume* („nach dem Workspace-Volume und dem Volume der Container-Runtime, in dieser Reihenfolge“) sagen dasselbe wie der Code. `ADR-0024` verwirft einen eigenen Schlüssel und widerspricht nicht. Die Spezifikation führt die Kombination nicht mehr als Nicht-Zusage; sie ist im Lastenheft durch „für die gemessenen Hosts und Versionen“ gedeckt (INFO V-3). <!-- d-check:ignore (Kennung als Zitat) -->

Ergebnis: **kein Verlust**. Die DoD-Zeile „Nichts geht verloren“ ist erfüllt.

## 4. Prüfpunkt 3: Pins der verbliebenen Zusagen

Existenz per `grep` im Baum, Rümpfe der Stichprobe gelesen:

| Zusage in der Spezifikation | Test (Fundstelle) | Pinnt wirklich? |
|---|---|---|
| `example.org`, `www.example.org`, `s3.us-east-1.amazonaws.com` erreichbar, `example.com` und bei `www.example.org`-Eintrag `example.org` per Name nicht | `TestE2E_SandboxDevcontainer_EgressSubdomains` (`sandbox_devcontainer_docker_test.go:198`) | ja (Review gelesen: Schlüssel `A …`, `B example.org`, Status `000`); nur lokal sichtbar |
| `*`-Eintrag Exit-Code `10` mit Hinweis | `TestNewEgressHost` (`domain/sandbox_test.go:75`), `TestConfigSet_EgressAllowWildcard_Exit10WithHint` (`config_egress_wildcard_test.go:22`) | Namen existieren; Inhalt nicht erneut gelesen (Slice-Klasse B, im Review gelesen) |
| Verzeichnis `/var/lib/u-boot/status` gehört `vscode` | `TestGenerateDevcontainer_SandboxStartScriptsWriteStatusFiles` (`generate_sandbox_test.go:1076`): `RUN mkdir -p /var/lib/u-boot/status && chown vscode:vscode …` | ja für den Besitz; „ein Prozess ohne `sudo` kann ersetzen“ folgt aus Verzeichnisbesitz (kein Test des Ersetzens, Zusage in der Form des Besitz-Pins) |
| Container über Label `devcontainer.local_folder` | `TestE2E_SandboxDevcontainer_StatusFiles_InspectorAgainstEngine` (`sandbox_status_docker_test.go:348`) | ja für die Abfrage; die Aussage zum Werkzeug ist gestrichen |
| Startscript mit der Liste wird beim Bau ins Image kopiert | `TestGenerateDevcontainer_SandboxEgress` (`generate_sandbox_test.go:576`): `COPY egress-init.sh /usr/local/bin/u-boot-egress-init` | ja |
| Ausgabe von `egress list` nennt den Behalt der Liste | `TestExecute_EgressList_Notes` (`egress_test.go:135`) pinnt „a container keeps the list of the image it was built from“ | ja für den Behalt; für „Neuaufbau braucht“ siehe V-1 |
| `registries.conf` `unqualified-search-registries = ['docker.io']`, `vfs`-Storage | `TestGenerateDevcontainer_SandboxPodman` (`:301`), `TestE2E_SandboxDevcontainer_NestedPodman` (`:455`, `REGISTRIES_DOCKER_IO`, `VFS_DRIVER`) | ja |
| genau zwei Hosts für `ollama`, `ollama.com` nicht | `TestGenerateDevcontainer_Ollama_EgressHostsAndTogetherWithClaudeCode` (`generate_ollama_test.go:383`); `TestE2E_SandboxDevcontainer_Ollama` (`:928`: `PULL_OK`, `RUN_OK`, `FOREIGN_BLOCKED`, `LISTEN:127.0.0.1:11434`, `LIBRARY_CPU`, `MODEL_SURVIVES_CONTAINER_REMOVAL`) | ja; trägt auch die Sätze zu Loopback, CPU und Modell nach Entfernen des Containers (`SPEC-047`) | <!-- d-check:ignore (Kennung als Zitat) -->
| Build (Podman) kalt, Cache, `--no-cache`, mit/ohne Egress | `TestE2E_SandboxDevcontainer_NestedPodman_Build` (`:380`), `..._BuildEgress` (`:422`) | Namen existieren |
| Reihenfolge der Volumes | `generate_ollama_test.go:353-355` | ja |

**Zeile 11 (Behalt der Liste, Neuaufbau):** Die Verhaltensbedingung „behält die Liste des Images“ / „wirkt nach einem Neuaufbau“ ist in der Fix-Runde `5c19443` aus `SPEC-046` gestrichen (Diff gelesen); die Spezifikation trägt dort jetzt den Mechanismus (Pin: `COPY`) und den Satz zur Ausgabe (Pin: Notes-Test). Restaussagen in `spec/spezifikation.md` mit dem Muster Neubau/Neuaufbau/behält gesucht: `SPEC-045` *Nicht ermittelbar* („braucht … einen Neubau“) und *Abbildung* („erst nach einem Neubau wirken“, Warnungstext) sind Vorbestand, Hinweise des Befehls an den Nutzer, nicht Gegenstand des Slice (Wortlaut des Warnungstexts als Ausgabe, nicht als Verhaltenszusage); ich werte sie nicht als Restaussage dieser Klasse. Eine Restaussage ohne Pin in `SPEC-046` bleibt nur in V-1. <!-- d-check:ignore (Kennung als Zitat) -->

## 5. Prüfpunkt 4: Ein Vertrag, vier Stellen

Lastenheft 0.17.0 (im Slice unverändert, Diff leer), Spezifikation, Handbuch, CHANGELOG:

| Gegenstand | Lastenheft | Spezifikation | Handbuch | CHANGELOG | Urteil |
|---|---|---|---|---|---|
| Egress-Anzeige | offline, aus Projektdateien, ohne Container | `SPEC-046`: Liste der Projektdateien, nicht der Zustand eines Containers | Befehlsbeschreibung | Eintrag `egress list` (unverändert) | kohärent | <!-- d-check:ignore (Kennung als Zitat) -->
| Subdomain-Semantik | Name erlaubt Subdomains (Lastenheft nennt keine Einzelfälle) | `example.org` ⇒ `www…`; `www.example.org` ⇒ `example.org` per Name nicht (gepinnt) | zusätzlich Messstand, Adresse bei geteilter IP nicht geprüft | — | kohärent; die Spezifikation behauptet nichts, was das Handbuch einschränkt (F-5 des Reviews) |
| Neuaufbau | schweigt | nur Mechanismus und Ausgabe | Zeilen 1332 bis 1343 mit Messstand | Eintrag `egress list` (Behalt) | kohärent; Spezifikation enger als Handbuch, kein Widerspruch |
| Upstream-Resolver | schweigt (Suche nach `resolver`, `upstream`, `dns`, `cname` im Absatz zu `LH-FA-DEV-008` ohne Treffer) | schweigt | beschreibt (Zeile 1220, mit Messstand) | schweigt (Suche nach `dns`/`resolver`/`upstream` ohne Treffer; Vorbestand) | **bewertet: tragfähig.** Eine Spezifikation darf nichts zusagen, was der Code nicht hat (der Resolver ist offen, Folge-Slice `slice-v1-egress-dns-upstream-schliessen`); Schweigen ist keine Zusage der Bindung. Schließt der Folge-Slice, ändern sich Handbuch und CHANGELOG, nicht die Spezifikation <!-- d-check:ignore (Kennung als Zitat) --> |
| `CNAME` | schweigt | schweigt | beschreibt | schweigt | wie Resolver |
| Statusdatei-Grenze | „Prozesse mit der nötigen Capability können sie aufheben“ (Grenze der Restriktion) | `SPEC-045` *Grenze*: Verzeichnis gehört `vscode`, Datei `root`, Prozess ohne `sudo` kann sie ersetzen; „Hinweis, kein Beleg“ | dasselbe mit Messstand (Zeile 1160) | — | kohärent; **Folgepflicht aus `ADR-0025`** („die Spezifikation führt den Absatz *Grenze* von `SPEC-045` nach“) bleibt erfüllt: der Absatz trägt die Zusage, nur der Messstand entfiel; die ADR lässt offen, wo der Resolver steht <!-- d-check:ignore (Kennung als Zitat) --> |
| `ollama` | CPU, Volume, nicht zugesagt: GPU, Host-Dienst, Modellverwaltung; „für die gemessenen Hosts und Versionen“ | `SPEC-047` Grenzen, Reihenfolge der Volumes, Verweis aufs Lastenheft | *Grenzen der Messung* inkl. `nestedRuntime` | Eintrag `ollama` inkl. `nestedRuntime` | kohärent | <!-- d-check:ignore (Kennung als Zitat) -->

Widerspruch zwischen den vier Stellen: **keiner**.

## 6. Prüfpunkt 5: Spec-Straten und Unveränderliches

- `AGENTS.md` §3.4: Der hinzugefügte Text in `spec/spezifikation.md` (Diff gelesen) enthält keinen Verweis auf ADR, Slice, Carveout, Roadmap, Welle oder `BEO-*`; die Historienzeile verlinkt seitwärts auf `LH-FA-DEV-007.b`. `spec/architecture.md` unverändert. <!-- d-check:ignore (Kennung als Zitat) -->
- `git diff --stat 6af82fd..HEAD` zeigt genau fünf Dateien; Lastenheft, ADRs, Code, Tests, Makefile, Dockerfile, Workflows, bestehende Reports sind nicht darunter (Abschnitt 1). `make doc-immutable` Exit 0.
- Kommentar-/Zustandsregeln (`AGENTS.md` §3.7, §3.9): Die hinzugefügte Prosa der Spezifikation ist im Indikativ; die Suche nach Chronik-Wörtern (Abschnitt 2) hat nur die Historie getroffen. Kein `//nolint`, keine Gate-Lockerung (§3.6), kein `git mv` im Diff (§3.3).

## 7. DoD-Abgleich (Slice §2) und Plan-vs-Diff (§3.10)

| DoD-Zeile | Haken | Stand | Beleg |
|---|---|---|---|
| Spezifikation | gesetzt | **erfüllt** | Abschnitt 2; Klassen A entschieden und umgesetzt, C-Zeilen 9, 11, 19 tragen Entscheidung (gestrichen; Zeile 11 mit Mechanismus-Pin), B und D im Wortlaut (die geplanten Teilstreichungen in Zeilen 3 und 12 ausgenommen), Historie trägt eine neue Zeile ohne Messstand, keine Abwärtsverweise |
| Nichts geht verloren, ein Vertrag an vier Stellen | gesetzt | **erfüllt** | Abschnitte 3 und 5 |
| Review und Verifikation | offen | **erfüllt mit diesem Report**; Haken beim Planner | Review-Report liegt vor; dieser Report liest die Kopplung gegen Code und Tests |
| `make gates`, `doc-trace`, `doc-planning`, `docs-check` | offen | **erfüllt** (Abschnitt 1); Läufe und Exit-Codes in §7 nachzutragen | Haken und Eintrag beim Planner |
| Closure-Notiz, Register, Risiken §6, Paarungen | offen | offen, zu Recht | §7 trägt den Platzhalter; `observations/` nicht im Diff (A-1 bis A-4) |

**Plan-Tabelle §3 gegen `git diff --stat 6af82fd..HEAD`** (Stand `5c19443`): `spec/spezifikation.md` (19 Zeilen), `docs/user/benutzerhandbuch.md` (4), `CHANGELOG.md` (2), Slice-Datei, Review-Report: je mit Plan-Zeile (Report unter `docs/reviews/`). Keine Datei im Diff ohne Plan-Eintrag. Plan-Zeilen ohne Diff: Verification-Report (dieser Report), Register, `open/README.md` und Roadmap (Closure; Abgleich-Absatz des Slice benennt sie). Der Abgleich des Slice nach der Fix-Runde stimmt mit dem Diff überein. Er ist nach der Closure (Move, Nachzüge, Register, dieser Report) zu wiederholen (A-5). Die „unverändert“-Zeilen (`spec/lastenheft.md`, `docs/plan/adr/`, `spec/architecture.md`, `internal/**`) sind im Diff leer.

## 8. Review-Findings F-1 bis F-6 einzeln

| Finding | Stand | Verifier-Beleg |
|---|---|---|
| F-1 LOW Zeile 11 nur zum Mechanismus gepinnt | **behoben** (Rest V-1) | `5c19443` streicht „behält … Liste des Images, aus dem er gebaut wurde“ und „wirkt also nach einem Neuaufbau“ aus `SPEC-046`; verbleibend Mechanismus (`COPY`-Pin) und Ausgabesatz (Notes-Pin). Die Teilpinnung der zweiten Satzhälfte ist V-1 | <!-- d-check:ignore (Kennung als Zitat) -->
| F-2 LOW „wie bisher“ in `LH-FA-DEV-007.b` | **behoben** | Diff: „kommt nur mit der Egress-Restriktion ins Image“; `nftables` ist in `Dockerfile.tmpl` hinter `{{- if .Egress}}` (Slice-Zeile 11a gelesen); Suche nach `wie bisher` im Baum der beiden Dateien: 0 Treffer | <!-- d-check:ignore (Kennung als Zitat) -->
| F-3 INFO Kopplung mit dem Code-Slice zu den Hinweisen | **übergeben, getragen** | `SPEC-046` nennt keinen Wortlaut des Messstands; Notes-Text trägt weiter „measured on Colima, arm64“ im Code (`egress.go:218` bis `220`), der offene Code-Slice ändert ihn; wer von beiden zuerst schließt, zieht nach (Slice §4/§6, Register `BEO-WZG/kopplung-nicht-nachgezogen`) | <!-- d-check:ignore (Kennung als Zitat) -->
| F-4 INFO Spezifikation schweigt zu Resolver/`CNAME` | **bestätigt, kein Mangel** | Abschnitt 5 |
| F-5 INFO Subdomain-Zusage ohne Gegengrenze | **bestätigt, kein Mangel** | Zusage ist auf „bei der Verbindung per Name“ eingeschränkt und gepinnt (Status `000`); das Handbuch trägt die Grenze; kein Widerspruch (Abschnitt 5) |
| F-6 INFO Absatz zur Beleglage in `LH-FA-DEV-007.b` | **begründet offen (keine Änderungspflicht)** | „Exit-Code `11`“ ist Zusage, „etwa Ubuntu 24.04 mit Standardeinstellung“ nennt den Auslöser als Beispiel, keinen Messstand; Klasse D, nicht Gegenstand des Slice | <!-- d-check:ignore (Kennung als Zitat) -->

## 9. Findings des Verifiers

| ID | Klasse | Befund |
|---|---|---|
| V-1 | LOW | `SPEC-046` sagt: „Die Ausgabe des Befehls nennt, dass ein Container die Liste seines Images behält und eine geänderte Liste einen Neuaufbau braucht.“ Die Ausgabe (`egress.go:218` bis `220`, Hilfetext `:102`) nennt den Behalt und den Wirkungsweg `devcontainer up --remove-existing-container … applied the new one`; das Wort „rebuild“/„Neuaufbau“ kommt im Ausgabetext nicht vor (Suche `rebuild` in `egress*.go`: 0 Treffer). `TestExecute_EgressList_Notes` pinnt den Behalt und den Satz zum `--remove-existing-container`, nicht „braucht einen Neuaufbau“. Die zweite Satzhälfte ist also mittelbar getragen, nicht wörtlich gepinnt (`AGENTS.md` §3.10, Bedingung ohne Test). Kein Failure-Szenario (der Sinn ist durch den genannten Weg erkennbar); keine Änderungspflicht in diesem Slice, aber die Lücke gehört dem Code-Slice `slice-v1-egress-hinweise-neuaufbau-und-messstand`, dessen Texte den Neuaufbau nennen sollen: er pinnt dann die Satzhälfte. Alternativ streicht die Spezifikation „und eine geänderte Liste einen Neuaufbau braucht“ (Entscheidung des Planners). | <!-- d-check:ignore (Kennung als Zitat) -->
| V-2 | INFO | `SPEC-047` Einleitung trägt das Wort „Messstand“ („Der Messstand der Zusagen steht im Lastenheft“). Verweis aufwärts, vom Plan als Klasse D belassen; es steht kein Messstand in der Spezifikation. Eine Umformulierung („Die Grenzen der Zusagen stehen im Lastenheft“) würde die Suche nach `Messstand` leer machen; nicht nötig. | <!-- d-check:ignore (Kennung als Zitat) -->
| V-3 | INFO | Die Kombination `ollama` mit `nestedRuntime` ist in der Spezifikation weder als Zusage noch als Nicht-Zusage geführt; getragen ist sie über „für die gemessenen Hosts und Versionen“ im Lastenheft und die Grenze in Handbuch und CHANGELOG. Das entspricht der Entscheidung (Zeile 17, Klasse A). |
| V-4 | INFO | Die Fix-Runde `5c19443` änderte `spec/spezifikation.md` nach dem Review-Report. Der Endtext ist von mir gelesen (Diff `5c19443`: 4 Zeilen: `SPEC-046` und `LH-FA-DEV-007.b`) und stimmt mit den Entscheidungen F-1 und F-2 überein; ein Nachlauf-Review ist für diese Streichung nicht nötig. Es ist ein weiteres Auftreten von `BEO-DOK/vertragsaenderung-nach-review` (Spezifikations-Text in der Fix-Runde nach dem Review), Zählung beim Planner (A-3). | <!-- d-check:ignore (Kennung als Zitat) -->
| V-5 | INFO (Release-Gate, kein Mangel) | Gates 1 und 2 (Befehlsübersicht `egress list`, Enum `egress`) sind beim Projektinhaber offen, kein Closure-Kriterium dieses Slice. Nicht im Verifier-Zugriff. |

Kein HIGH, kein MEDIUM, keine DoD-Verletzung, keine Abweichung vom Plan-Wortlaut außer den im Slice benannten (Zeilen 6 und 7: Spezifikation schweigt, Beschreibung im Handbuch).

## 10. Auflagen

**Beim Planner vor `done/`:**

- **A-1 Closure-Notiz §7** füllen: Läufe und Exit-Codes (Abschnitt 1), Stand der Gates 1 und 2 (offen, Auftrag des Release-Cut), Steering-Loop-Lerneintrag, Folge-Slice `slice-v1-devcontainer-ollama-amd64-messung` startet erst nach `done/` (Reihenfolge der Folge). <!-- d-check:ignore (Kennung als Zitat) -->
- **A-2 Entscheidung zu V-1:** die zweite Satzhälfte in `SPEC-046` („und eine geänderte Liste einen Neuaufbau braucht“) entweder streichen (Spezifikations-Änderung, dann Fix-Commit mit Abgleich nach §3.10) oder als Übergabe an `slice-v1-egress-hinweise-neuaufbau-und-messstand` in §6 festhalten, der den Ausgabetext und seinen Pin um den Neuaufbau ergänzt. Eine Ablehnung der Änderung ist zulässig, wenn die Übergabe dokumentiert ist. <!-- d-check:ignore (Kennung als Zitat) -->
- **A-3 Register:** `BEO-DOK/zusage-weiter-als-messung` (kein neues Auftreten festgestellt; der Verifier fand in der Spezifikation keinen Satz, der weiter reicht als sein Pin außer V-1), `BEO-WZG/kopplung-nicht-nachgezogen` (Kopplung getragen: Handbuch, CHANGELOG und Spezifikation stimmen überein), `BEO-DOK/vertragsaenderung-nach-review` (V-4), `BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf` (Grenze jetzt in Handbuch und CHANGELOG übertragen) fortschreiben oder in §7 notieren, dass keine Beobachtung angefallen ist.
- **A-4 Risiken §6** auf Ist-Stände setzen: Risiko 1 (Grenze verloren): behoben, Zeile 17 vor dem Streichen übertragen, vom Verifier gegengelesen; Risiko 2 (Resolver): Entscheidung Schweigen getragen, Prüfung Abschnitt 5; Risiko 3 (Folgepflicht `ADR-0025`): erfüllt, vom Verifier geprüft; Risiko 4 (Kopplung zu zwei Code-Slices): weiter offen; Risiko 5 (Messstand bleibt in der Spezifikation): Suche am Ende der Umsetzung wiederholt, vom Verifier unabhängig bestätigt (Abschnitt 2); Risiko 6 (Pin nur lokal sichtbar): weiter offen, der Lauf der E2E-Pins ist in diesem Report **nicht** gefahren (Abschnitt 1, Grund dort). <!-- d-check:ignore (Kennung als Zitat) -->
- **A-5 Closure-Mechanik:** Roadmap-Arbeitszeile auf „Nichts in Arbeit.“ zurücksetzen; reiner `git mv` nach `done/` nach der Closure-Notiz (`AGENTS.md` §3.3, `make slice-mv`); danach `make docs-check` (Verweise aus `open/README.md`, dem Folge-Slice `slice-v1-devcontainer-ollama-amd64-messung` und der Roadmap nachgezogen); Plan-Tabelle nach der letzten Änderung erneut gegen `git diff --stat` abgleichen (§3.10) — dieser Report ist eine zusätzliche Datei unter `docs/reviews/` (Plan-Zeile vorhanden). <!-- d-check:ignore (Kennung als Zitat) -->
- **A-6 Restliche Haken** setzen (Review und Verifikation, Sensoren, Closure, Register, Risiken); die Paarungen (Anker, Folge-Slice, Register) von der nächsten Welle-Closure prüfen lassen, wie der Slice es festhält.

**Für den Projektinhaber (keine Mängel):**

- **E-1 (Release-Gate):** Gates 1 und 2 der Welle `welle-sandbox-egress-ollama` (Befehlsübersicht `egress list`, Enum `egress`) bestätigen oder ändern; kein Gegenstand dieses Slice.
- **E-2 (optional):** V-2, falls die Spezifikation auch das Wort „Messstand“ nicht mehr tragen soll; Umformulierung wäre eine Spezifikations-Änderung ohne Bedeutungsänderung.

## 11. Nicht verifiziert

`make test-docker` und damit der tatsächliche Lauf der E2E-Pins (`TestE2E_SandboxDevcontainer_EgressSubdomains`, `..._Ollama`, `..._NestedPodman*`, `..._StatusFiles_InspectorAgainstEngine`): nur gelesen, nicht gefahren; `make verify-depguard` (keine Importregel im Diff); Wahrheit der Messstände im Handbuch gegen die Messprotokolle (`rohdaten-*`): nicht Verifier-Kontext, unverändert gegenüber dem Stand vor dem Slice; Inhalt von `TestNewEgressHost` und `TestConfigSet_EgressAllowWildcard_Exit10WithHint` nur auf Existenz geprüft (der Review hat die Rümpfe der übrigen Stichprobe gelesen); die Gates 1 und 2 (Entscheidung des Projektinhabers).
