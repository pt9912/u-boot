# Verification-Report: slice-v1-handbuch-forge-cli-nachbesserung — 2026-10-23

**Rolle:** Verifier ([`harness/verification.md`](../../harness/verification.md)), frischer Kontext, kein Self-Review, kein Reparieren. Frage: Wurde der Slice richtig gebaut, gemessen an Plan, DoD und Aussage-Beleg-Tabelle?

**Gegenstand:** [`slice-v1-handbuch-forge-cli-nachbesserung`](../plan/planning/done/slice-v1-handbuch-forge-cli-nachbesserung.md) (Welle [`welle-forge-cli-gh-glab`](../plan/planning/done/welle-forge-cli-gh-glab.md)), Diff `f010a1f0..36d24f2a`, Vertrag [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer). Review: [`2026-10-23-slice-v1-handbuch-forge-cli-nachbesserung.md`](2026-10-23-slice-v1-handbuch-forge-cli-nachbesserung.md) (0 HIGH, 0 MEDIUM, 4 LOW, 3 INFO), Fix-Runde `36d24f2a`.

## Ergebnis

**Schließbar** nach den Closure-Pflichten des Planners (Haken Review, Closure-Notiz, Register, Risiko-Ausgänge, Paarungen, Move nach `done/`). Kein Blocker, kein MEDIUM. Ein neues LOW (V-1, Plan-Aussage der Welle-Datei veraltet) und zwei INFO. Alle Sensoren, deren Ausgabe ich nicht übernehmen wollte, habe ich selbst gefahren (Linux `amd64`, lokale Docker Engine 29.9.2).

## Sensoren — selbst gefahren

| Sensor | Ergebnis | Beleg |
| --- | --- | --- |
| `make gates` (auf `36d24f2a`) | **Exit 0** (sichtbar als `EXIT=0` am Logende) | `ablage-pruefung` 26 Fälle, 0 Fehlschläge |
| `make doc-trace` | Exit 0 | 84 Anforderungen, 0 Waisen |
| `make doc-planning` | Exit 0 | 775 Dateien, 0 Befunde |
| `make doc-immutable RANGE="4655bcff..HEAD"` | Exit 0 | 59 Commits, 0 Befunde |
| `make test-docker`, `make verify-depguard` | nicht gelaufen | Grund: kein Docker-Pfad, keine Importregel berührt; der Diff ändert nur Testdateien im Paket `application`, Handbuch, Plan und einen Spec-Satz |

Arbeitsbaum vor und nach allen Läufen und Mutationen: `git status --short` leer.

## Mutationen (Rücknahme per `git apply -R`, Baum danach leer)

Die Genehmigungsprüfung hat die Mutationen nicht abgelehnt. Ausgangslauf ohne Mutation: `make test` Exit 0.

| Nr. | Mutation | Ergebnis (`make test`, Exit 2) |
| --- | --- | --- |
| A | `validateDevcontainerFeatures` (`devcontainer_features.go`, Zeile 641): Bedingung um `&& false` erweitert (Prüfung aus) | rot: `TestConfigSet_ForgeVersionWithoutSource_RejectsEverySet` **alle acht Teilfälle**, darunter die neuen `gh/…version`, `glab/…version`, `gh/devcontainer.featureSources.allow`, `glab/devcontainer.featureSources.allow`; `TestDoctor_UbootYaml_ErrorOnForgeVersionWithoutSource` (gh, glab); `TestGenerateDevcontainer_ForgeCLI_VersionOverrideWithoutSourceRejected` (4 Fälle); dazu `…Ollama_VersionOverrideWithoutSourceRejected` und die CLI-Tests `TestRun_Generate_ForgeCLINotices` (`version_without_source_gh`, `_glab`) und `TestRun_Generate_OllamaNotices` |
| B | dieselbe Bedingung ohne `entry.Version != ""` | rot: `TestGenerateDevcontainer_ForgeCLI_RemovingVersionEntryLetsGenerateRun` (T3), `TestConfigSet_ForgeVersionWithoutSource_SetItselfSucceeds` (T4), dazu viele Bestandstests des Forge-Pfads |
| C | `revalidateFeatureEntry` (`config.go`, Zeile 577): `.version` ohne `source` bei `gh`/`glab` abgelehnt | rot: nur `TestConfigSet_ForgeVersionWithoutSource_SetItselfSucceeds` (T4) |

Die neuen Teilfälle `…version` und `…featureSources.allow` (Fix F-3) werden bei Mutation A rot: Der Allquantor „jedes weitere `config set`“ ist für vier Pfade je Feature gepinnt. `SetMany` gibt es laut Plan in `ConfigService` nicht (Aussage des Implementers, `config.go`; ich habe die Aussage nicht gesondert gesucht, der Handbuchsatz nennt nur `config set`).

## Messung A4 — wiederholt

Alpine, `GH_TOKEN=PLACEHOLDER-x`, `GITLAB_TOKEN=PLACEHOLDER-y` im Host, Container mit `sleep`: `docker exec` ohne `-e`: beide ungesetzt; mit `-e GH_TOKEN -e GITLAB_TOKEN`: beide mit dem Wert des Hosts; auf dem Host ungesetzt und `-e GH_TOKEN`: im Container ungesetzt. Reproduziert, Aussagen des Handbuchs (Platzhalter, ein Lauf, Container ohne Forge-CLIs) stimmen mit dem Lauf. `-u <Benutzer>` war nicht Teil meines Laufs und nicht des Implementer-Laufs; der Text führt es aus der Messung des Validators.

## DoD-Abgleich

| DoD-Punkt | Befund |
| --- | --- |
| Liefer-Punkt 1 (Standardprofil, A-2, A-3-Satz) | erfüllt. Spiegelstrich „Standardprofil, keine Token-Referenz“ nennt den Weg `docker exec -e` mit Platzhalter-Vorbehalt, Aufruf mit gültigem Token „nicht gemessen und nicht zugesagt“, Login im Container „nicht geprüft“ samt Verbleib des Tokens; beide `shell`-Sätze („Token-Variablen“, „Zugangsdaten“) auf das Sandbox-Profil eingeschränkt. |
| Liefer-Punkt 2 (Anmeldung gebündelt) | erfüllt. Zwei Absätze „Anmeldung in kurzer Form“ und „Ohne Variable …“ am Kopf; der `up`-Satz ist gestrichen (kein Pin, F-2); jeder Satz wiederholt A2 bis A5 oder einen Spiegelstrich (`401` und Exit-Code 4 stehen in den Spiegelstrichen „Ohne Variable“ und „Mit Variable“). Der Vorbehalt „Gemessen mit einem Platzhalter ist nur, dass die Variable dann im Container ankommt“ steht an der Aussage (F-1). |
| Liefer-Punkt 3 (Abhilfe bei `version`) | erfüllt. Absatz vor den Spiegelstrichen; Bedingungen durch T1 (config set), T2 (doctor), `…VersionOverrideWithoutSourceRejected` (generate), T3 (Abhilfe wirkt), T4 (derzeit Erfolg bei `config set …version`, Ablaufpunkt „Zurzeit“) gepinnt, Mutationen A bis C belegen die Pins. |
| `make gates` grün, `docs-check`, Kennungsfreiheit | erfüllt. Im Handbuch-Diff Suche nach `LH-`, `ADR-`, `SPEC-`, `ARC-`, `slice-`, `welle-`: 0 Treffer (§5 Regel 6 gehalten); `uboot.yaml.valid` ist ein Prüfungsname. |
| Doku-Update | erfüllt (Handbuch; Roadmap-Zeile, `open/README.md`, Welle-Datei, Folge-Slice-Zeiger). |
| Review durchgeführt | Report liegt vor, Haken in §2 noch offen: Aufgabe des Planners nach diesem Bericht. |
| Closure-Notiz, Register, Risiko-Ausgänge, drei Paarungen | **offen**: §7 trägt „(Wird bei der Closure gefüllt …)“, §6 die Ausgänge „offen bis zur Closure“. Auftrag des Planners vor dem Move, keine Lücke der Umsetzung. |

## Aussage → Beleg, A1 bis A11

| # | Befund |
| --- | --- |
| A1, A2, A3 | Bestandstests existieren (`TestGenerateDevcontainer_ForgeCLI_DefaultProfileJSONUnchanged`, `TestSessionEnvNames`, `…RemoteEnvReferences`); die Mutations-Läufe führen sie mit aus (grün). Die Spec-Einschränkung („Im Sandbox-Profil gibt `shell` …“) ist durch `TestSessionEnvNames` und `…RemoteEnvReferences` gepinnt, die Standardprofil-Seite durch A1/A2. |
| A4 | M, selbst reproduziert (siehe oben). Grenzen im Plan ausgewiesen (F-6): `-u` nicht gelaufen, Image ohne `gh`/`glab`, kein Protokoll unter `docs/messungen/`. |
| A5 | G: nur „nicht geprüft“, kein Volume-Halbsatz, keine „geht verloren“-Annahme. Im Handbuch gelesen. |
| A6 | Jeder Satz der Einleitung nennt seine Quelle in der Plan-Tabelle „Stand“; gegen die Spiegelstriche gelesen, kein Satz neu belegt. |
| A7 | Test vorhanden, in Mutation A rot. |
| A8 | T2 vorhanden, in Mutation A rot (gh, glab). |
| A9 | T1 vorhanden, in Mutation A rot (alle acht Teilfälle inkl. `version`, `featureSources.allow`). |
| A10 | T3 vorhanden, in Mutation B rot (beide Schlüssel, beide Profile). |
| A11 | T4 vorhanden, in Mutationen B und C rot; Handbuchsatz sagt „Zurzeit“ (Ablaufpunkt). |

## §3.4, §3.10, Plan-Tabelle

- **§3.4 (Spec-Änderung):** Diff der Spezifikation ist ein Satz in [`SPEC-051`](../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) „Zugangsdaten“ („Im Sandbox-Profil gibt …“). Suche im Spec-Diff nach `slice`, `welle`, `ADR-0`: 0 Treffer; `architecture.md` und Lastenheft ohne Diff.
- **§3.10:** Die Bedingung in der Spec („Im Sandbox-Profil …“) hat Pins (siehe A1 bis A3). Zustandsmarker „Stand der Umsetzung“ trägt Ist-Stand inklusive Fix-Runde (T1-Pfade, A4-Grenzen).
- **Plan-Tabelle gegen `git diff --stat f010a1f0..HEAD`:** elf Dateien, jede hat eine Zeile der Tabelle (Handbuch, drei Testdateien, Spec, Folge-Slice-Zeiger, `open/README.md`, Welle-Datei, Roadmap, Review, diese Slice-Datei); „unverändert“ geführt: `CHANGELOG.md`, `README.md`, `README.de.md` (Diff leer, selbst geprüft), `spec/lastenheft.md` (Diff leer). Mit diesem Report kommt `docs/reviews/…-verify.md` unter die Plan-Zeile `docs/reviews/`.

## Findings

| ID | Klasse | Befund | Fundstelle |
| --- | --- | --- | --- |
| V-1 | LOW | Die Zeile des Slice in der Welle-Datei (Slice-Tabelle) endet auf „Spezifikation und Lastenheft unverändert“. Seit der Fix-Runde (F-5) ändert der Slice einen Satz in [`SPEC-051`](../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten); die Aussage ist falsch (§3.10: Zustandsaussage über den Diff). Die Plan-Tabelle des Slice und sein Kopf nennen die Spec-Änderung richtig. | `docs/plan/planning/welle-forge-cli-gh-glab.md` Zeile 70 („Spezifikation und Lastenheft unverändert“) |
| V-2 | INFO | Messung A4 ohne Protokoll unter `docs/messungen/` und ohne `-u <Benutzer>`; im Text ausgewiesen, ändert den Haken nicht. | Plan §3 Zeile A4 |
| V-3 | INFO | `SetMany` als weiterer Schreibpfad: die Aussage „gibt es nicht“ stammt vom Implementer und ist hier nicht selbst gegengelesen. Der Handbuchtext nennt nur `config set`, daher ohne Folge. | Plan §3 Zeile A9 |

## Ungemessen oder nicht gepinnt — trägt es den Haken?

| Gegenstand | Stand im Text | Trägt den Haken? |
| --- | --- | --- |
| Aufruf von `gh`/`glab` mit gültigem Token (auch auf dem Weg `docker exec -e`) | „nicht gemessen und nicht zugesagt“; es gibt kein Token | ja — als Nicht-Zusage geschrieben, keine Abhilfe ohne Probe behauptet |
| Antwort von `gh`/`glab` im Standardprofil via `docker exec -e` (Platzhalter, Image mit CLI) | Text sagt nur, dass die Variable im Container ankommt; die Antwort `401` steht im Spiegelstrich „Mit Variable“ aus dem Sandbox-Lauf des Vorgänger-Slice | ja — Aussagebreite deckt sich mit der Messung; die Kombination Weg + Werkzeug bleibt ungemessen und ist so benannt |
| `gh auth login`, `glab auth login`, Verbleib eines Tokens nach Neuanlegen | „nicht geprüft“ | ja |
| `up` reicht Variablen im Sandbox-Profil durch | Aussage gestrichen | ja |

## Urteil

**Schließbar.** V-1 ist eine Berichtigung in einer Planungsdatei (eine Zeile) und gehört in die Closure des Planners; V-2 und V-3 sind Hinweise. Offen bleiben ausschließlich die Closure-Pflichten des Planners.
