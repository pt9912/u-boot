# Review-Report: [`slice-v1-handbuch-forge-cli-nachbesserung`](../plan/planning/done/slice-v1-handbuch-forge-cli-nachbesserung.md) — 2026-10-23

**Review-Art:** Code (Handbuch-Text und Tests gegen Plan und Hard Rules)

**Gegenstand:** Diff `f010a1f0..dd9d8b79` (Slice [`slice-v1-handbuch-forge-cli-nachbesserung`](../plan/planning/done/slice-v1-handbuch-forge-cli-nachbesserung.md), Welle `welle-forge-cli-gh-glab`; Anlass: Befunde M-1, M-2 (Handbuch-Teil) und L-1 des Validator-Berichts `2026-10-23-welle-forge-cli-gh-glab-validierung`)

**Skill:** `.harness/skills/reviewer.md` @ `dd9d8b79`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-23

**Eingangs-Kontext:**

- Slice-Plan (§1 bis §3, Tabelle „Aussage → Beleg“, Stand der Umsetzung)
- [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-CONF-005`](../../spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern); [`ADR-0035`](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md); [`SPEC-051`](../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten), [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)
- `AGENTS.md` §3.7, §3.9, §3.10, §5 Zeile 6
- Validator-Bericht `2026-10-23-welle-forge-cli-gh-glab-validierung`

**Sensors dieses Laufs:**

- `make gates`: Exit-Code 0 (selbst gefahren).
- Mutationen A, B, C selbst nachgefahren (Patch per `git apply`, `make test`, Rücknahme per `git apply -R`, `git status --short` danach leer):
  - A (`&& false` an der Prüfung in `validateDevcontainerFeatures`): rot `TestConfigSet_ForgeVersionWithoutSource_RejectsEverySet` (alle vier Teilfälle), `TestDoctor_UbootYaml_ErrorOnForgeVersionWithoutSource` (gh, glab), `…VersionOverrideWithoutSourceRejected`.
  - B (Bedingung ohne `entry.Version != ""`): rot `TestGenerateDevcontainer_ForgeCLI_RemovingVersionEntryLetsGenerateRun` (T3) und `TestConfigSet_ForgeVersionWithoutSource_SetItselfSucceeds` (T4), dazu Bestandstests.
  - C (`revalidateFeatureEntry` lehnt `.version` ohne `source` ab): rot nur T4. Die Mutationen sind plausibel und tragen die Tests.
- Messung A4 selbst wiederholt (`alpine`, Engine 29.9.2, Platzhalter): ohne `-e` ungesetzt; mit `-e GH_TOKEN -e GITLAB_TOKEN` beide mit dem Wert des Hosts; Host ungesetzt und `-e GH_TOKEN`: im Container ungesetzt. Reproduziert, wie im Text gesagt.
- `make test-docker`, `make verify-depguard`: nicht gelaufen. Grund: Der Diff ändert weder einen Docker-Pfad noch Importregeln (nur Handbuch, Plan, Testdateien im Paket `application`).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Die Einleitung schreibt zum Standardprofil-Weg „(gemessen mit einem Platzhalter)“. Gemessen ist nur, dass die Variable im Container ankommt, in einem Container ohne `gh`/`glab`; eine Anfrage der Werkzeuge auf diesem Weg ist nicht gemessen. Der Vorbehalt steht erst im nächsten Absatz („ein Aufruf mit einem gültigen Token ist nicht gemessen“) und im Spiegelstrich, nicht an der Klammer. | `AGENTS.md` §3.10 (Aussagebreite); Plan §6 Risiko „Wiederholung der Überbehauptung in der Einleitung“ | `docs/user/benutzerhandbuch.md` · „<Container> …` mit (gemessen mit einem Platzhalter).“ | Review-only | Messvorbehalt an der Aussage statt im Nachbarsatz |
| F-2 | LOW | Die Einleitung sagt, `u-boot devcontainer up` und `shell` gäben die Variable im Sandbox-Profil „beide“ weiter. A6 führt den Satz auf A3 zurück; A3 pinnt die Referenz in der `devcontainer.json` und `shell` (`TestSessionEnvNames`), nicht das Verhalten von `up`. Ein Test oder eine Messung für `up` ist im Plan nicht genannt. | `AGENTS.md` §3.10 (Bedingung im Handbuch ohne Pin) | `docs/user/benutzerhandbuch.md` · „`u-boot devcontainer up` oder `u-boot devcontainer shell` aufrufen; beide geben sie weiter.“ | Review-only | Handbuch-Bedingung ohne Pin |
| F-3 | LOW | „jedes weitere `u-boot config set`“ lehnt ab: T1 pinnt zwei Pfade je Feature (`project.name`, `…enabled`). Nicht gepinnt sind der Pfad `…version` und `devcontainer.featureSources.allow` auf der schon ungültigen Datei sowie `SetMany`. Der Allquantor ist breiter als der Beleg; der Plan (A9) führt ihn als eine Zeile. | `AGENTS.md` §3.10 (Bedingung im Handbuch ohne Pin) | `docs/user/benutzerhandbuch.md` · „und jedes weitere `u-boot config set` die Datei mit Exit-Code `10` ab“ | `make test` (Test ergänzen) | Allquantor breiter als Test |
| F-4 | LOW | Zwei Plan-Zustandsangaben stimmen nicht mehr: (a) Die Zeile des Slice in `docs/plan/planning/open/README.md` endet auf „nur Plan“, obwohl der Slice Handbuch und vier Tests liefert; (b) die Plan-Tabelle §3 führt die Datei des Folge-Slice zur Prüfung der Feature-Version beim `config set` nicht, obwohl der Diff sie berührt (3 Verweise, 6 Zeilen). Alle anderen Diff-Dateien stehen in der Tabelle; die als „unverändert“ geführten Dateien (`CHANGELOG.md`, READMEs, Spec) sind im Diff nicht vorhanden. | `AGENTS.md` §3.10 (Abgleich Tabelle gegen `git diff --stat`, Zustandsmarker beim Ablauf) | `docs/plan/planning/open/README.md` · „vier Tests ohne Produktcode; keine neue Zusage“ | `git diff --stat` | Plan-Tabelle/Zustandsmarker nicht nachgezogen |
| F-5 | INFO | Der Spezifikationsabschnitt zu `gh` und `glab` im Devcontainer sagt im Absatz „Zugangsdaten“: „`u-boot devcontainer shell` gibt beide Variablen an die Sitzung weiter, soweit sie auf dem Host gesetzt sind“, im selben Absatz wie die Aussage zum Standardprofil, aber ohne Profil-Einschränkung im Satz. Der Plan liest das als richtig und lässt die Spezifikation unverändert; dieselbe Lesart führte im Handbuch zu M-1. Entscheidung beim Architect/Projektinhaber; kein Fehler dieses Diffs. | Maintainability | `spec/spezifikation.md` · „gibt beide Variablen an die Sitzung weiter“ | Review-only | Spec-Satz ohne Profilbezug |
| F-6 | INFO | Die Messung A4 ist nur im Plan (Stand der Umsetzung) beschrieben, ohne Messprotokoll unter `docs/messungen/`; der Befehl steht nicht im Repo. Ich habe sie mit eigenen Befehlen reproduziert (siehe oben). Der Handbuchtext nennt Umgebung, Engine, „ein Lauf“ und „ohne die Forge-CLIs“ ausreichend; `-u <Benutzer>` war nicht Teil des Laufs (stammt aus der Messung des Validators). | Ablage-Regel (`harness/conventions.md`) | Slice-Plan §3 · „Gemessen vom Implementer mit Platzhaltern“ | Review-only | Messung ohne Protokoll |
| F-7 | INFO | Abweichung „Heredoc beim Anhängen der Tests“: Der Inhalt der 117 hinzugefügten Zeilen ist unauffällig (gofmt-konform, `make lint` grün, Tabs, Dateiende mit genau einem Zeilenumbruch, keine CRLF, keine Leerzeichen am Zeilenende; ein typografisches „…“ in einem Test-Kommentar). Bewertung: harmlos, Prozessabweichung ohne Wirkung im Ergebnis. | Maintainability | `internal/hexagon/application/config_features_test.go` · „Characterisation of the state before“ | `make lint` | — |

## Prüfpunkte laut Auftrag

- **Einleitung gegen Spiegelstriche:** Aus „gemessen mit Platzhalter“ wird kein „funktioniert“: Der Absatz nennt „ein Aufruf mit einem gültigen Token ist nicht gemessen“ und „`gh auth login` … nicht geprüft“; Befund nur als F-1 (Platzierung des Vorbehalts). `401`, Exit-Code 4, `glab` liest öffentliches Projekt: deckungsgleich mit den Spiegelstrichen „Ohne Variable“ und „Mit Variable, ohne Anmeldung durch `u-boot`“.
- **docker-exec-Weg und A4:** plausibel und reproduziert; die Aussagebreite „Container ohne die Forge-CLIs“, „Host ungesetzt → im Container fehlt“ stimmt mit dem Lauf überein; „ebenso `-e GITLAB_TOKEN`“ ist gemessen (zusammen mit `GH_TOKEN`).
- **`shell`-Sätze:** beide auf das Sandbox-Profil eingeschränkt; die Standardprofil-Aussage deckt sich mit `TestSessionEnvNames`, `…DefaultProfileJSONUnchanged` und der Spezifikation.
- **Abhilfe-Absatz (Exit-Code 10 bei `generate` und `config set`):** stimmt. T1 pinnt `ErrConfigSchemaInvalid` für `config set` auf einer schon ungültigen Datei; der Exit-Code-Pin im CLI-Adapter wurde selbst gelesen: `internal/adapter/driving/cli/cli_test.go`, Fall „ErrConfigSchemaInvalid (validation)“ → 10 (Zeile 988), und die Zuordnung in `cli.go` (Zeile 464, Kohorte 10). Kein Widerspruch zum Satz „zurzeit meldet `config set …version` Erfolg“: T4 beschreibt eine gültige Datei, T1 eine bereits ungültige. `doctor`-Fehler in `uboot.yaml.valid`: T2 (gh, glab). Die Abhilfe wirkt: T3 (beide Schlüssel, beide Profile). Der Satz „Zurzeit“ hat seinen Ablaufpunkt in der Plan-Tabelle des Folge-Slice.
- **Kennungsfreiheit:** Der Handbuch-Diff trägt keine `LH-*`-, `ADR-*`-, `SPEC-*`- oder Slice-Kennung; `uboot.yaml.valid` ist ein Prüfungsname, keine Kennung.
- **Kommentare (§3.7/§3.9):** Die Test-Kommentare tragen Zusage, Kopplung und Abgrenzung; keine Muster (a) bis (c). Der Kommentar „The assertion turns to a rejection with the follow-up …“ nennt den Ablaufpunkt ohne Kennung und ist akzeptabel.
- **Roadmap-Zeile:** Zustand und Anker, keine Chronik.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `docs/user/benutzerhandbuch.md` (Abhilfe-Absatz, Einleitung, Spiegelstriche, Zugangsdaten-Satz) | geprüft; Befunde F-1 bis F-3, sonst ohne Befund |
| `internal/hexagon/application/*_test.go` (T1 bis T4) | geprüft, ohne Befund (Mutationen A, B, C nachgefahren) |
| `internal/adapter/driving/cli` (Exit-Code-Pin, Dual-Classifier) | geprüft, ohne Befund; kein Sentinel neu, kein Produktcode im Diff |
| Roadmap, Welle-Datei, Folge-Slice (Zeiger) | geprüft; Zeiger korrekt nachgezogen, Befund F-4 |
| Hard Rules (Docker-only, Suppression, Lifecycle `git mv`, Referenzmodell) | geprüft, ohne Befund (reiner Move-Commit `088c1684` vor der Inhaltsänderung) |
| `make test-docker`, `make verify-depguard` | nicht geprüft: kein Docker-Pfad, keine Importregel berührt |
| DoD-Abgleich, Closure-Notiz, Register | nicht geprüft: Aufgabe des Verifiers |

## Summary

0 HIGH, 0 MEDIUM, 4 LOW (F-1 bis F-4), 3 INFO (F-5 bis F-7). `make gates` Exit-Code 0. Der Diff trägt seinen Plan; die Handbuch-Aussagen sind bis auf die Randstellen in F-1 bis F-3 durch Test oder Messung gedeckt.
