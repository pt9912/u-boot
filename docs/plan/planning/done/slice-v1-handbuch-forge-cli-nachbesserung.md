# Slice V1: Handbuch — Standardprofil-Weg zum Token, Einschränkung der `shell`-Sätze, Anmeldung gebündelt und Abhilfe bei ungültiger `version` (Auflagen des Validators zu `gh` und `glab`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-forge-cli-gh-glab`](welle-forge-cli-gh-glab.md) (Release `v0.15.0`). Fünfter von sechs Slices, **vor** dem Release-Cut (`slice-v1-release-cut-v0.15.0`, Datei erst beim Trigger): Er setzt die Auflagen A-2 und A-3 und die Lücken M-1, M-2 (nur den Handbuch-Satz) und L-1 des Validator-Berichts [`2026-10-23-welle-forge-cli-gh-glab-validierung`](../../../reviews/2026-10-23-welle-forge-cli-gh-glab-validierung.md) um; der Cut nimmt den Stand des Handbuchs mit. Die Welle schließt erst, wenn dieser Slice in `done/` liegt (Welle-Datei §3 Punkt 1, §4).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Katalog-Features `gh`, `glab`), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Token als Umgebungsvariable des Hosts), [`ADR-0035`](../../adr/0035-forge-clis-gh-und-glab-im-devcontainer.md). Der Slice ändert **keine Zusage**: Er schränkt Sätze des Handbuchs auf das Profil ein, in dem sie gemessen sind, und benennt, was ungeprüft ist.

**Berührte Spec-Stellen:** — (das Handbuch ist ein öffentlicher Vertrag, kein Spec-Stratum). Gelesen am 2026-10-23: [`SPEC-051`](../../../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) (Absatz „Wirkung in beiden Profilen“: Standardprofil `devcontainer.json` byte-gleich, ohne `remoteEnv`) und [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (Absätze „Umgebung und Zugangsdaten“ und „Umgebung der Sitzung“: die Namen kommen aus dem Text der `devcontainer.json`) sagen das Standardprofil richtig; **die Spezifikation bekommt nur eine Einschränkung** (der `shell`-Satz in [`SPEC-051`](../../../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) „Zugangsdaten“ nennt das Sandbox-Profil, Fix-Runde, Befund F-5). Der Fehler liegt allein in zwei Handbuch-Sätzen (§3).

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-23.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein Nutzer im Standardprofil, der `gh` oder `glab` aktiviert, findet im Handbuch einen Weg, das Token in den Container zu bringen (oder die ausdrückliche Aussage, was ungeprüft ist), die `shell`-Sätze behaupten nichts, was im Standardprofil nicht gilt, die Antwort auf „wie melde ich mich an?“ steht an einer Stelle statt über vier Spiegelstriche verteilt, und wer eine ungültige `version` in `u-boot.yaml` stehen hat, erfährt die Abhilfe.

### Befund (Validator-Bericht, gelesen am 2026-10-23)

- **M-1 / A-2 (MEDIUM).** Im Standardprofil bekommt der Container keine Token-Referenz; `u-boot devcontainer shell` gibt dort `GH_TOKEN` und `GITLAB_TOKEN` nicht weiter (Validator, Punkt 5: im Container beide Variablen ungesetzt). Zwei Sätze des Handbuchs lesen sich im Standardprofil falsch: im Abschnitt „Forge-CLIs (optional)“ der Spiegelstrich „Token-Variablen“ (Satz „`u-boot devcontainer shell` öffnet die Shell mit den beiden Variablen, soweit sie auf dem Host gesetzt sind …“) und in „Zugangsdaten (Git und Agent)“ der Satz „`u-boot devcontainer shell` öffnet die Shell mit `GIT_TOKEN`, dem Schlüssel des Agenten und den Token-Variablen von `gh` und `glab`, soweit sie auf dem Host gesetzt sind“. Der gemessene Weg `docker exec -e GH_TOKEN -u vscode <Container> …` (Validator: Variable gesetzt, `gh api` erreicht GitHub, Antwort `Bad credentials` mit dem Platzhalter) fehlt im Forge-Abschnitt; der Spiegelstrich „Standardprofil, keine Token-Referenz“ sagt dort „ein bestimmter Weg ist weder zugesagt noch gemessen“.
- **A-3-Satz (zu M-1).** Ob `gh auth login` im Container taugt und was mit einem dabei abgelegten Token beim Neuanlegen des Containers geschieht, steht nirgends. Der Validator hat es **nicht gefahren**; „geht mit dem Container verloren“ ist seine Annahme, keine Messung.
- **L-1 (LOW).** Der Abschnitt hat rund 110 Zeilen Spiegelstriche mit Messvorbehalten; die Antwort auf „wie melde ich mich an?“ ist über „Token-Variablen“, „Ohne Variable“, „Mit Variable, ohne Anmeldung durch `u-boot`“ und „Standardprofil, keine Token-Referenz“ verteilt. Der Abschnitt liegt unter 4.11 (Sandbox), obwohl die Features in beiden Profilen wirken.
- **M-2, nur der Handbuch-Teil (Auflage A-3 des Validators).** Steht in `u-boot.yaml` ein `version` zu `gh` oder `glab` ohne `source`, lehnen `generate devcontainer` (Exit-Code `10`), `doctor` (`uboot.yaml.valid`, Fehler) und jedes weitere `config set` die Datei ab; die Abhilfe (den Eintrag `version` in `u-boot.yaml` von Hand entfernen) nennt keine Meldung und das Handbuch nicht. Die Ursache (`config set …version` prüft nicht wie `generate`) behebt [`slice-v1-config-set-validiert-feature-version`](../in-progress/slice-v1-config-set-validiert-feature-version.md) und nicht dieser Slice.

### Entscheidung zu L-1: Einleitung statt Verschieben (Planner, vorgelegt)

Geprüft am 2026-10-23: Ein Verschieben der Überschrift `#### Forge-CLIs (optional)` aus 4.11 in 4.10 bräche keinen Link, denn alle fünf Verweise sind Überschrift-Anker ohne Nummer (`#forge-clis-optional`, fünf Stellen in `docs/user/benutzerhandbuch.md`: die Katalogzeilen `gh` und `glab` von 4.10, die Egress-Zeilen `gh` und `glab` und der Abschnitt „Zugangsdaten“), und das Inhaltsverzeichnis führt nur die Kapitel der obersten Ebene. Der Slice wählt trotzdem den **Einleitungssatz am Kopf des Abschnitts** (Rezept je Profil), weil der Abschnitt Sandbox-Sachverhalte trägt (Egress, `doctor`-Prüfung, Token-Referenz), die in 4.11 stehen, ein Verschieben von rund 110 Zeilen den Diff und den Review um ein Vielfaches vergrößert und die Auffindbarkeit über die Katalogzeilen in 4.10 schon gegeben ist. **Wünscht der Projektinhaber das Verschieben, ist es ein eigener Handbuch-Slice** (reiner Umbau, keine Zusage).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Ursache von M-2** (`config set …version` meldet Erfolg und sperrt danach das Projekt): Code-Änderung, eigener Slice [`slice-v1-config-set-validiert-feature-version`](../in-progress/slice-v1-config-set-validiert-feature-version.md), ohne Welle und nicht vor dem Cut. Dieser Slice ändert keinen Code.
- **Eine Zusage oder ein Feature** (ein Weg, auf dem `u-boot` das Token im Standardprofil durchreicht, eine `remoteEnv`-Referenz im Standardprofil, ein Befehl zum Anmelden): wäre eine Vertragsänderung ([`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`SPEC-051`](../../../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten): das Standardprofil bleibt ohne `remoteEnv`).
- **Eine Messung mit einem gültigen Token** (`gh auth login` interaktiv, Anmeldung, `gh pr create`): es gibt kein Token für Messungen; was ungeprüft bleibt, steht als „nicht geprüft“ im Text und nicht als Aussage.
- **`CHANGELOG.md` und Lead** (Auflage A-1 des Validators: Zusage-Menge Linux `amd64`, lokale Engine, kein authentifizierter Aufruf gemessen): gehören in den Release-Cut-Slice (Welle-Datei §5).
- **Wortlaute der Meldungen** (L-2: `glab` meldet bei fehlendem Host `lookup gitlab.com … server misbehaving`; L-3: Länge und Wiederholung der Warnungen mit Kennung als Präfix): Code beziehungsweise Werkzeug-Verhalten; notiert in [`slice-v1-config-set-validiert-feature-version`](../in-progress/slice-v1-config-set-validiert-feature-version.md) §6.
- **Das Verschieben des Abschnitts** nach 4.10: siehe Entscheidung zu L-1.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **Liefer-Punkt 1, Standardprofil (M-1, A-2, A-3-Satz):** Im Abschnitt „Forge-CLIs (optional)“ nennt der Spiegelstrich „Standardprofil, keine Token-Referenz“ den Weg `docker exec -e GH_TOKEN` (und `GITLAB_TOKEN`) als mit einem Platzhalter gemessen, ein Aufruf mit gültigem Token als nicht gemessen; der `shell`-Satz im Spiegelstrich „Token-Variablen“ und der `shell`-Satz in „Zugangsdaten (Git und Agent)“ sind auf das Sandbox-Profil eingeschränkt; ein Satz sagt, dass `gh auth login` und `glab auth login` im Container nicht geprüft sind (kein Aufruf mit Token, Verbleib eines abgelegten Tokens über das Neuanlegen des Containers nicht geprüft). Jede Aussage hat einen Test, einen genannten Beleg oder ist gestrichen (Tabelle „Aussage → Beleg“ in §3).
- [x] **Liefer-Punkt 2, Anmeldung gebündelt (L-1):** Der Abschnitt beginnt mit einer Einleitung von höchstens zwei Absätzen, die „wie melde ich mich an?“ je Profil beantwortet (Sandbox: Variable auf dem Host setzen, `shell` — für `up` gibt es keinen Pin, der Handbuchtext sagt dazu nichts; Standardprofil: `docker exec -e`; ohne Variable: das gemessene Verhalten; Login im Container: nicht geprüft) und auf die Spiegelstriche verweist; die Spiegelstriche bleiben als Belege stehen. Kein Satz der Einleitung ist neu belegt: jeder wiederholt einen Satz aus Liefer-Punkt 1 oder einen schon vorhandenen, gemessenen Spiegelstrich (Fundstelle in §3).
- [x] **Liefer-Punkt 3, Abhilfe bei ungültiger `version` (M-2, A-3):** Ein Satz im Abschnitt „Forge-CLIs (optional)“ (an der Stelle, die `version` ohne `source` erwähnt, Absatz vor den Spiegelstrichen) nennt: Steht in `u-boot.yaml` ein `version` zu `gh` oder `glab` ohne `source`, lehnen `generate devcontainer`, `doctor` und jedes `config set` die Datei mit Exit-Code `10` beziehungsweise als Fehler ab; die Abhilfe ist, den Eintrag `version` in `u-boot.yaml` von Hand zu entfernen. Die Bedingungen des Satzes sind durch Tests gepinnt (§3, Zeilen T1 bis T4), oder der Teilsatz entfällt.
- [x] `make gates` grün; `docs-check` ohne Befund; `docs/user/benutzerhandbuch.md` trägt keine `LH-*`-, `ADR-*`-, `SPEC-*`- oder Slice-Kennung und keine Chronik ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 6).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: `docs/user/benutzerhandbuch.md` (öffentlicher Vertrag; der einzige berührte neben dem Planning-Baum, siehe §3).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht. Jede Aussage über Dateien, Funktionen und Verhalten wurde am 2026-10-23 an der genannten Stelle gelesen ([`AGENTS.md`](../../../../AGENTS.md) §3.10); der Implementer erweitert die Tabelle um jede Aussage seiner Texte samt Fundstelle und gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab.

| Datei / Komponente | Änderungs-Art | Begründung (Fundstelle) |
|---|---|---|
| `docs/user/benutzerhandbuch.md` | update | Abschnitt `#### Forge-CLIs (optional)` (Überschrift Zeile 2036): Einleitung nach den Befehlen und dem Absatz „Die Features `gh` und `glab` installieren …“ (Zeilen 2044 bis 2055; dort steht der Halbsatz „`version` ohne `source` weist `u-boot` mit Exit-Code `10` ab“ in den Zeilen 2052 und 2053); Spiegelstrich „Token-Variablen“ (Zeilen 2069 bis 2080, der `shell`-Satz in den Zeilen 2079 und 2080); Spiegelstrich „Standardprofil, keine Token-Referenz“ (Zeilen 2095 bis 2100: „ein bestimmter Weg ist weder zugesagt noch gemessen“); Abschnitt `#### Zugangsdaten (Git und Agent)` (Überschrift Zeile 2149), der `shell`-Satz in den Zeilen 2168 bis 2170 |
| `internal/hexagon/application/config_features_test.go` (oder die Testdatei, die der Implementer für `ConfigService.Set` findet) | update | Tests T1 und T4 (siehe unten); Testart Negative und Charakterisierung. Fundstelle: `readUbootYAMLBody` (`config.go`, Zeilen 652 bis 681) ruft `validateDevcontainer` und wickelt einen Fehler in `driving.ErrConfigSchemaInvalid`; `revalidateFeatureEntry` (`config.go`, Zeilen 552 bis 585) prüft bei `.version` nur `entry.Version == ""` |
| `internal/hexagon/application/doctor_features_test.go` (oder die Testdatei der Prüfung `uboot.yaml.valid`) | update | Test T2. Fundstelle: `checkUbootYaml` (`doctor.go`, Zeilen 573 bis 626) ruft `validateDevcontainer` und gibt `domain.SeverityError` mit `checkIDUbootYaml` zurück |
| `internal/hexagon/application/generate_forge_cli_test.go` | update | Test T3 (Abhilfe wirkt). Fundstelle: `TestGenerateDevcontainer_ForgeCLI_VersionOverrideWithoutSourceRejected` (Zeile 443) pinnt die Ablehnung durch `generate` für `gh` und `glab` in beiden Profilen schon |
| `CHANGELOG.md` | unverändert | gelesen (Zeilen 16 und 20): Der Block `[Unreleased]` nennt „`version:` without `source:` ends with exit code `10`“ (Zeile 16); der Satz gilt für `generate`, nicht für `config set` — die Berichtigung gehört in den Release-Cut-Slice (Auflage A-1, Welle-Datei §5), weil der Cut den Block dem Release zuordnet |
| `README.md`, `README.de.md` | unverändert | gelesen: Suche nach `GH_TOKEN` in beiden Dateien ohne Treffer (2026-10-23); beide führen `gh` und `glab` nur als Katalog-Schlüssel |
| `spec/spezifikation.md` | update (Fix-Runde, nur Einschränkung) | Absatz „Zugangsdaten“ von [`SPEC-051`](../../../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) (Zeile 763): der Satz zu `shell` beginnt mit „Im Sandbox-Profil“; keine neue Zusage, die Zeile steht im selben Absatz wie die Aussage zum Standardprofil (Befund F-5 des Reviews) |
| `spec/lastenheft.md` | unverändert | gelesen (Zeile 413): keine Zeile nennt `shell` zusammen mit `GH_TOKEN` |
| `docs/plan/planning/in-progress/slice-v1-config-set-validiert-feature-version.md` | update (Zeiger) | `make slice-mv` hat die Verweise auf den Folge-Slice nachgezogen (3 Verweise, 6 Zeilen im Diff); kein Inhalt geändert |
| `docs/plan/planning/open/README.md`, `docs/plan/planning/welle-forge-cli-gh-glab.md`, `docs/plan/planning/in-progress/roadmap.md` | update | Zeiger auf diesen Slice (Planner; Welle §4 und §5, Roadmap-Zeile der Welle); `make slice-mv` zieht die Zeiger bei den Moves nach |
| `docs/reviews/2026-10-23-slice-v1-handbuch-forge-cli-nachbesserung.md` | neu | Review-Report des Reviewers |
| `docs/reviews/2026-10-23-slice-v1-handbuch-forge-cli-nachbesserung-verify.md` | neu | Verify-Report des Verifiers (schließbar) |
| `docs/plan/planning/observations/BEO-DOK/{zusage-weiter-als-messung,spec-bedingung-ohne-code-entsprechung,plan-aussage-ueber-code}/` | update (Closure) | je eine Datei `evidence/slice-v1-handbuch-forge-cli-nachbesserung.md` und ein Absatz in `state.md` (§7) |
| `docs/plan/planning/in-progress/slice-v1-config-set-validiert-feature-version.md` | update (Closure) | Notiz in §6 zum fehlenden Pin für `SetMany` (§7) |

**Aussage → Beleg** (jede neue oder geänderte Aussage des Handbuchs; Klasse: **T** = Test im Repo, **M** = gemessen und im Text als „mit einem Platzhalter, ein Lauf, Linux `amd64`“ gekennzeichnet, **G** = gestrichen):

| # | Aussage im Handbuch | Klasse | Beleg |
|---|---|---|---|
| A1 | Im Standardprofil trägt die `devcontainer.json` keine Token-Referenz und kein `remoteEnv` (steht schon im Text, wird wiederholt) | T | `TestGenerateDevcontainer_ForgeCLI_DefaultProfileJSONUnchanged` (`generate_forge_cli_test.go`, Zeile 410: kein `remoteEnv`, kein `GH_TOKEN`, kein `GITLAB_TOKEN`) |
| A2 | `u-boot devcontainer shell` gibt nur Variablen weiter, auf die die `devcontainer.json` mit `"NAME": "${localEnv:NAME}"` verweist; im Standardprofil also weder `GH_TOKEN` noch `GITLAB_TOKEN` | T | `TestSessionEnvNames` (`devcontainer_shell_test.go`, Zeile 32: Fall „none“, „another name in the reference“, „a reference in mounts“) zusammen mit A1; gemessen vom Validator (Standardprofil, beide Variablen im Container ungesetzt, ein Lauf) |
| A3 | Im Sandbox-Profil trägt die `devcontainer.json` die Referenzen auf `GH_TOKEN` und `GITLAB_TOKEN`, `shell` gibt sie weiter, soweit auf dem Host gesetzt | T | `TestGenerateDevcontainer_ForgeCLI_RemoteEnvReferences` (`generate_forge_cli_test.go`, Zeile 351) und `TestSessionEnvNames` (Fall „the shape the generator writes“) |
| A4 | Im Standardprofil bringt `docker exec -e GH_TOKEN -u <Benutzer> <Container> …` die Variable in den Aufruf, wenn sie auf dem Host gesetzt ist; ohne `-e` fehlt sie (mit Platzhalter gemessen; ein Aufruf mit gültigem Token nicht gemessen) | M | Validator-Bericht, Punkt 5 („Ohne `-e`: ungesetzt“; mit `-e`: `gh api` erreicht GitHub, `Bad credentials`); derselbe Engine-Mechanismus steht für `GIT_TOKEN` in Abschnitt 4.9 „Im Container arbeiten“ (Zeilen 875 bis 883, „geprüft in einem Lauf“). Der Implementer wiederholt den Lauf mit einem Platzhalter, wenn die lokale Docker-Engine es zulässt (`GITLAB_TOKEN` nur, wenn er es misst); sonst bleibt `GITLAB_TOKEN` im Text gestrichen und es steht nur `GH_TOKEN` |
| A5 | `gh auth login` und `glab auth login` sind im Container nicht geprüft; wohin das Werkzeug ein Token schreibt und ob es das Neuanlegen des Containers übersteht, ist nicht geprüft | G für die Folge | Gestrichen: „geht mit dem Container verloren“ (Annahme des Validators, nicht gemessen; ein Platzhalter kann `gh auth login --with-token` nicht abschließen, weil das Werkzeug das Token gegen den Server prüft — nicht gelesen, nicht gemessen, daher keine Begründung im Text). Der Satz sagt nur „nicht geprüft“; er sagt **nicht**, dass `u-boot` kein Volume anlegt (`TestGenerateDevcontainer_ForgeCLI_DefaultProfileJSONUnchanged` pinnt `mounts` nur für das Standardprofil, `generate_forge_cli_test.go` Zeilen 410 bis 430; für das Sandbox-Profil fand der Planner keinen Test auf `mounts` im Forge-Test — der Implementer ergänzt den Halbsatz nur mit einem Test oder lässt ihn weg). Die Formulierung folgt dem Muster des Claude-Code-Abschnitts („Wohin ein Login im Container seine Zugangsdaten schreibt, ist nicht gemessen“, Zeilen 1903 bis 1906) |
| A6 | Einleitung (Rezept je Profil): jeder Satz wiederholt A2 bis A5 oder einen vorhandenen Spiegelstrich („Ohne Variable“, Zeilen 2081 bis 2087, gemessen mit einem Platzhalter beziehungsweise ohne Variable) | T/M | die Belege von A2 bis A5; der Implementer nennt je Satz die Zeile des Spiegelstrichs, aus dem er stammt |
| A7 | Steht `version` ohne `source` zu `gh` oder `glab` in `u-boot.yaml`, lehnt `generate devcontainer` die Datei mit Exit-Code `10` ab | T | `TestGenerateDevcontainer_ForgeCLI_VersionOverrideWithoutSourceRejected` (Zeile 443) |
| A8 | …ebenso `u-boot doctor` (Fehler in der Prüfung `uboot.yaml.valid`) | T | **T2** (neu): Konfiguration mit `gh` oder `glab` und `version` ohne `source`, `checkUbootYaml` meldet `SeverityError` mit `checkIDUbootYaml`; der Implementer sucht zuerst, ob ein bestehender Test denselben Fall für die Feature-Schema-Prüfung abdeckt, und ergänzt `gh` dort |
| A9 | …ebenso jedes weitere `u-boot config set` (Exit-Code `10`) | T | **T1** (neu): `ConfigService.Set` auf einer Datei mit `gh.version: "9.9.9"` ohne `source` endet mit einem Fehler, der `driving.ErrConfigSchemaInvalid` wickelt, und die Datei bleibt byte-identisch; dazu der `ExitCode`-Pin des CLI-Adapters für diese Sentinel (`cli_test.go`, Zeile 988 prüft `ErrConfigSchemaInvalid` → 10 schon) |
| A10 | Abhilfe: den Eintrag `version` in `u-boot.yaml` von Hand entfernen — danach läuft `generate devcontainer` wieder | T | **T3** (neu): dieselbe Konfiguration ohne den Eintrag `version`: `generate` endet ohne Fehler und schreibt das Dockerfile mit dem Block (die Wirkungsprobe der Abhilfe vor dem Text, `BEO-DOK/abhilfe-ohne-wirkungsprobe`) |
| A11 | `config set devcontainer.features.gh.version <Wert>` meldet derzeit Erfolg und schreibt die Datei, bevor die Prüfung greift (dieser Teilsatz hat einen Ablaufpunkt: er entfällt mit [`slice-v1-config-set-validiert-feature-version`](../in-progress/slice-v1-config-set-validiert-feature-version.md)) | T | **T4** (neu, Charakterisierung): `Set` auf einer gültigen Datei ohne `version` für `gh` mit `2.0.0` ohne `source` endet ohne Fehler und schreibt; der Folge-Slice dreht diesen Test auf die Ablehnung. Entscheidet der Implementer, den Teilsatz zu streichen, entfällt T4 und der Satz sagt nur den Zustand der Datei (A7 bis A10) |

**Stand der Umsetzung (Implementer, 2026-10-23).** Die Entscheidung zu L-1 (Einleitung statt Verschieben) ist bestätigt; Frage 5 der Welle bleibt offen. Erreichter Beleg je Zeile:

| # | Stand |
|---|---|
| A1, A2, A3, A7 | Bestehende Tests, unverändert gelesen und gelaufen (`TestGenerateDevcontainer_ForgeCLI_DefaultProfileJSONUnchanged`, `TestSessionEnvNames`, `TestGenerateDevcontainer_ForgeCLI_RemoteEnvReferences`, `TestGenerateDevcontainer_ForgeCLI_VersionOverrideWithoutSourceRejected`) |
| A4 | Gemessen vom Implementer mit Platzhaltern in `GH_TOKEN` und `GITLAB_TOKEN`, Linux `amd64`, Docker Engine 29.9.2, ein Lauf, `alpine` ohne die Forge-CLIs: ohne `-e` beide ungesetzt; mit `-e GH_TOKEN -e GITLAB_TOKEN` beide mit dem Wert des Hosts; auf dem Host ungesetzt und `-e GH_TOKEN` fehlt die Variable im Container (keine leere Zeichenkette). **Grenzen der Messung:** `-u <Benutzer>` war nicht Teil des Laufs (der Text führt es aus dem Lauf des Validators), das Image trug weder `gh` noch `glab`, und die Messung hat kein Protokoll unter `docs/messungen/`; der Text sagt deshalb nur, dass die Variable im Container ankommt. Die Antwort von `gh` auf den Platzhalter (`Bad credentials`) hat der Implementer nicht wiederholt (kein Image mit `gh` lokal); der Text verweist dafür nur auf den schon vorhandenen Spiegelstrich „Mit Variable, ohne Anmeldung durch `u-boot`“ (`401`) |
| A5 | Als „nicht geprüft“ geschrieben; der Halbsatz „`u-boot` legt kein Volume an“ steht nicht in der Einleitung oder im Standardprofil-Spiegelstrich (der Spiegelstrich „Beide Profile“ trägt ihn schon: „Die Features legen kein Volume an“, unverändert) |
| A6 | Einleitung, Absatz „Anmeldung in kurzer Form“: Sandbox-Satz = A3 und Spiegelstrich „Token-Variablen“; Standardprofil-Satz = A2 und A4; „Ohne Variable“ = Spiegelstrich „Ohne Variable“; `401` und „gültiges Token nicht gemessen“ = Spiegelstrich „Mit Variable, ohne Anmeldung durch `u-boot`“; Login = A5 |
| A8 | T2: `TestDoctor_UbootYaml_ErrorOnForgeVersionWithoutSource` (`doctor_features_test.go`, `gh` und `glab`); rot gesehen mit Mutation A |
| A9 | T1: `TestConfigSet_ForgeVersionWithoutSource_RejectsEverySet` (`config_features_test.go`; `gh` und `glab`, Pfade `project.name`, `devcontainer.features.<Name>.enabled` und `…<Name>.version`, dazu `devcontainer.featureSources.allow`; pinnt nur `Set`; `SetMany` (`config_many.go`, Zeile 106) ist ein zweiter Schreibpfad und hier **nicht** gepinnt (Berichtigung der Closure, siehe §7); Fehler wickelt `ErrConfigSchemaInvalid`, Datei byte-identisch); rot gesehen mit Mutation A |
| A10 | T3: `TestGenerateDevcontainer_ForgeCLI_RemovingVersionEntryLetsGenerateRun` (`generate_forge_cli_test.go`, beide Schlüssel, beide Profile); rot gesehen mit Mutation B |
| A11 | T4 (Charakterisierung): `TestConfigSet_ForgeVersionWithoutSource_SetItselfSucceeds` (`config_features_test.go`); rot gesehen mit Mutation B und Mutation C |

Mutationen (alle am geprüften Code, vorher committet, Rücknahme per `git apply -R`): **A** `devcontainer_features.go`, Bedingung in `validateDevcontainerFeatures` um `&& false` erweitert (Prüfung aus) — rot: T1, T2, `…VersionOverrideWithoutSourceRejected`; **B** dieselbe Bedingung ohne `entry.Version != ""` (lehnt auch ohne `version` ab) — rot: T3, T4; **C** `config.go`, `revalidateFeatureEntry` lehnt `.version` ohne `source` ab (Folge-Slice) — rot: T4.

Die Tests T1 bis T4 ändern keinen Produktcode; fällt einer rot aus, weil der Code etwas anderes tut als die Aussage, ist die Aussage falsch gelesen, und der Satz wird berichtigt oder gestrichen, nicht der Code geändert.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Der Projektinhaber gibt den Slice frei (Release-Zuordnung `v0.15.0`, vor dem Cut); kein anderer Slice liegt in `in-progress/` (WIP-Limit 1). Der Slice darf vor oder nach [`slice-v1-config-set-validiert-feature-version`](../in-progress/slice-v1-config-set-validiert-feature-version.md) laufen, hängt aber nicht an ihm: Der Satz zu M-2 ist so geschrieben, dass er vor dem Fix wahr ist (A11 mit Ablaufpunkt) und danach nur dort geändert wird.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Implementer will die Überschrift verschieben oder mehr als eine Einleitung umbauen; das ist der Handbuch-Umbau außerhalb dieses Plans (§1).
- `in-progress` → `open` (blockiert): Eine Aussage aus A4 lässt sich weder messen noch streichen, ohne dass der Projektinhaber über den Wortlaut entscheidet.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor (Findings behoben oder mit Ausgang), `make gates` grün, Closure-Notiz mit Lerneintrag geschrieben (§7). Der Release-Cut der Welle startet erst, wenn dieser Slice in `done/` liegt (Welle-Datei §4).

## 6. Risiken und offene Punkte

Jedes Risiko bekommt bei der Closure genau einen Ausgang (eingetreten / entfallen / weiter offen).

- **Aussage ohne Beleg (`BEO-DOK/spec-bedingung-ohne-code-entsprechung`, `BEO-DOK/plan-aussage-ueber-code`):** Der Standardprofil-Weg und der Satz zu `version` sind Handbuch-Bedingungen über Verhalten; der Validator hat sie an einem Binary gemessen, der Plan liest sie am Code. Eine Aussage über das Werkzeug (A4, A5) ist keine über `u-boot`-Code und kann nur mit Messung oder Streichung bestehen. — **Ausgang:** eingetreten, behoben: der `up`-Satz stand ohne Pin (F-2, gestrichen), der Allquantor „jedes weitere `config set`“ trug zwei Pfade (F-3, auf vier erweitert), und der Plan nannte `SetMany` als nicht vorhanden (V-3, Berichtigung in §3 A9 und §7). Die Tabelle „Aussage → Beleg“ in §3 trägt je Aussage Test, Messung oder Streichung; der Verifier las sie gegen die Mutationen A bis C.
- **Wiederholung der Überbehauptung in der Einleitung (L-1):** Eine gebündelte Antwort verleitet dazu, Messvorbehalte zu kürzen und aus „mit Platzhalter gemessen“ ein „funktioniert“ zu machen — derselbe Fehler, der M-1 erzeugte (`BEO-DOK/zusage-weiter-als-messung`). — **Ausgang:** eingetreten, behoben: Der Vorbehalt „mit einem Platzhalter“ stand zunächst im Nachbarsatz statt an der Aussage (F-1, in der Fix-Runde an die Aussage gesetzt); der Verifier las die Einleitung gegen die Spiegelstriche und fand keine Aussage breiter als Test oder Messung.
- **Der Teilsatz A11 veraltet mit dem Folge-Slice:** Er sagt „derzeit“ und hat einen benannten Ablaufpunkt; geht [`slice-v1-config-set-validiert-feature-version`](../in-progress/slice-v1-config-set-validiert-feature-version.md) zuerst nach `done/`, streicht dieser Slice den Teilsatz selbst (kein T4). — **Ausgang:** weiter offen → [`slice-v1-config-set-validiert-feature-version`](../in-progress/slice-v1-config-set-validiert-feature-version.md): Der Folge-Slice liegt in `open/`, der Teilsatz steht mit „Zurzeit“ im Handbuch, T4 pinnt ihn (Mutation C rot); der Folge-Slice dreht T4 auf die Ablehnung und streicht den Teilsatz.
- **Verschieben des Abschnitts (L-1) bleibt ungeklärt:** Der Plan entscheidet für die Einleitung; der Projektinhaber kann das Verschieben wollen. — **Ausgang:** weiter offen: Frage 5 der Welle-Datei §5 steht beim Projektinhaber aus; falls gewünscht, ein eigener Handbuch-Slice (reiner Umbau, keine Zusage).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

Stand der Closure: 2026-10-23. Quellen: Review-Report [`2026-10-23-slice-v1-handbuch-forge-cli-nachbesserung`](../../../reviews/2026-10-23-slice-v1-handbuch-forge-cli-nachbesserung.md) (0 HIGH, 0 MEDIUM, 4 LOW, 3 INFO), Verify-Report [`2026-10-23-slice-v1-handbuch-forge-cli-nachbesserung-verify`](../../../reviews/2026-10-23-slice-v1-handbuch-forge-cli-nachbesserung-verify.md) (schließbar, `cad2d243`; `make gates` Exit 0 selbst gefahren, Mutationen A, B, C, Messung A4 wiederholt).

**Geliefert.** Im Handbuch-Abschnitt „Forge-CLIs (optional)“: der Absatz „Anmeldung in kurzer Form“ am Kopf (Liefer-Punkt 2), der Weg `docker exec -e GH_TOKEN` (und `GITLAB_TOKEN`) im Standardprofil mit Platzhalter-Vorbehalt, die zwei `shell`-Sätze auf das Sandbox-Profil eingeschränkt, „nicht geprüft“ zu `gh auth login` und `glab auth login` (Liefer-Punkt 1), der Absatz zur Abhilfe bei `version` ohne `source` (Liefer-Punkt 3). Vier Tests ohne Produktcode: `TestConfigSet_ForgeVersionWithoutSource_RejectsEverySet` (T1), `TestDoctor_UbootYaml_ErrorOnForgeVersionWithoutSource` (T2), `TestGenerateDevcontainer_ForgeCLI_RemovingVersionEntryLetsGenerateRun` (T3), `TestConfigSet_ForgeVersionWithoutSource_SetItselfSucceeds` (T4). Ein Satz in [`SPEC-051`](../../../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) („Zugangsdaten“) auf das Sandbox-Profil eingeschränkt (Fix-Runde, F-5); das Lastenheft ist unverändert.

**Aussage → Beleg am erreichten Stand.** A1 bis A3 und A7: Bestandstests, im Mutationslauf grün mitgeführt. A4: Messung (Linux `amd64`, Engine 29.9.2, Platzhalter, ein Lauf, Container ohne die Forge-CLIs), vom Verifier reproduziert; `-u <Benutzer>` und die Antwort der Werkzeuge auf diesem Weg ungemessen und im Text so ausgewiesen, kein Messprotokoll unter `docs/messungen/` (F-6, V-2). A5: gestrichen bis auf „nicht geprüft“. A6: Einleitung gegen die Spiegelstriche gelesen (Reviewer, Verifier). A8 (T2) und A9 (T1, vier Pfade je Feature, nur `Set`): rot bei Mutation A. A10 (T3): rot bei Mutation B. A11 (T4): rot bei Mutation B und C; der Satz sagt „Zurzeit“.

**Was funktionierte.** Die Tabelle „Aussage → Beleg“ vor dem Text und die drei Mutationen am geprüften Code: Sie machten aus „der Test pinnt es“ eine rot gesehene Aussage, und der Verifier fand die neuen T1-Pfade bei Mutation A rot. Das Streichen statt Absichern (der `up`-Satz, die „geht verloren“-Annahme des Validators) hielt den Text auf dem Gemessenen. Die Messung A4 wurde vom Implementer, Reviewer und Verifier je selbst gefahren und stimmte jedes Mal.

**Was anders lief.**

- **F-1 bis F-6 des Reviews** (alle in der Fix-Runde `36d24f2a`): F-1 Messvorbehalt an die Aussage statt in den Nachbarsatz; F-2 der Satz „`up` und `shell` geben die Variable weiter“ gestrichen, weil nur `shell` und die Referenz in der `devcontainer.json` gepinnt sind; F-3 T1 um die Pfade `…version` und `devcontainer.featureSources.allow` erweitert (der Allquantor „jedes weitere `config set`“); F-4 `open/README.md` und Plan-Tabelle auf den Ist-Stand; F-5 Spec-Satz auf das Sandbox-Profil eingeschränkt (der Plan hatte die Spezifikation als „unverändert“ geführt, derselbe Satz hatte im Handbuch M-1 erzeugt); F-6 die Grenzen der Messung A4 im Plan ausgewiesen.
- **V-1** (Verifier, LOW): Die Zeile des Slice in der Welle-Datei sagte „Spezifikation und Lastenheft unverändert“; berichtigt in der Closure (Welle-Datei §4, Zeile des Slice).
- **V-3** (Verifier, INFO, vom Planner am Code gelesen): Die Aussage der Plan-Tabelle, `ConfigService` habe keinen weiteren Schreibpfad neben `Set`, war falsch. `SetMany` (`internal/hexagon/application/config_many.go`, Zeile 106; vom CLI-Adapter in `internal/adapter/driving/cli/config.go`, Zeile 666 aufgerufen) liest die Datei über `readUbootYAMLBody` (`config.go`, Zeile 652, ruft `validateDevcontainer`), also über dieselbe Eingangsprüfung wie `Set`, und lehnt eine Datei mit `version` ohne `source` ebenso ab. Der Handbuchsatz „jedes weitere `config set`“ stimmt damit am Code, ist aber für die Mehrpfad-Form **nicht gepinnt**. Die `WriteFile`-Stellen von `ConfigService` (`grep -n WriteFile` über `config.go`, `config_many.go`, `devcontainer_egress.go`): `config.go` Zeilen 276 (`Set`) und 999 (`setFeatureSourcesAllow`, aufgerufen in Zeile 216) und `devcontainer_egress.go` Zeile 227 (`setEgressAllow`, aufgerufen in `config.go` Zeile 220) liegen alle hinter dem `readUbootYAMLBody` von `Set`; `config_many.go` Zeile 145 gehört zu `SetMany`. Ein Schreibpfad von `config set` ohne Eingangsprüfung existiert nicht (andere Dienste, die `u-boot.yaml` schreiben, etwa `add` und `remove`, sind nicht `config set` und nicht Gegenstand des Handbuchsatzes). Die Berichtigung steht in §3 (Zeile A9, Stand); die Lücke des Pins ist als Notiz V-3 in §6 von [`slice-v1-config-set-validiert-feature-version`](../in-progress/slice-v1-config-set-validiert-feature-version.md) geführt (Ausgang weiter offen: ein `SetMany`-Fall beim Drehen von T4, oder der Handbuchsatz nennt nur die Einzelform).
- **Heredoc beim Anhängen der Tests:** Der Implementer hängte die Testfälle per Heredoc an die Testdateien an. Das Ergebnis ist unauffällig (F-7: gofmt-konform, `make lint` grün, keine CRLF, Dateiende korrekt); eine Prozessabweichung ohne Wirkung.
- **Abgelehnter Commit und abgelehnte Mutationsprobe:** Die Genehmigungsprüfung lehnte im Implementer-Lauf einen Commit und die Mutationsprobe ab. Gelöst durch die Arbeitsteilung: Der Projektleiter-Kontext committete die geprüften Dateien, die Mutation fuhr der Verifier (A, B, C, Rücknahme per `git apply -R`, Baum danach leer; der Verifier meldet, dass die Prüfung seine Läufe nicht ablehnte). Die Rot-Belege der Tabelle stammen damit aus zwei Kontexten, dem des Implementers (Plan §3, Stand) und dem des Verifiers (Verify-Report), und stimmen überein.

**Steering-Loop-Eintrag (Lerneintrag).** Gelernt: Eine Handbuchaussage mit „jedes“ oder „kein weiterer“ gilt für die Menge der Eingangswege, nicht für den Weg, den der Test nimmt. Der Implementer pinnte `Set` und schrieb „`SetMany` gibt es nicht“ in den Plan; kein Gate, Review und Verify fingen es als Falschaussage: Der Reviewer sah den Pin-Rand (F-3), der Verifier wies die Negativaussage als ungelesen aus (V-3), erst das Lesen des Codes durch den Planner der Closure löste sie auf. Eine Negativaussage („gibt es nicht“) braucht eine Suche mit Fundstelle (`grep` über das Paket, Treffer oder 0 Treffer), keine Lesart. Die Regel steht schon als `AGENTS.md` §3.10 (Aussage über Code lesen und die Fundstelle nennen); der Slice ist ein Rückfall der Klasse „Negativaussage ohne Fundstelle“ ([`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md), 48 Belege) und der Klasse „Allaussage breiter als der Pin“ ([`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md), 21 Belege). **Keine neue Verkörperung in diesem Slice** (kein Feld `liegt in`): Beide Einträge sind schon verkörpert, eine Schärfung entscheidet der Projektinhaber. Vorschlag zur Entscheidung: `AGENTS.md` §3.10 um den Satz ergänzen „Eine Negativaussage über Code (‚gibt es nicht‘, ‚kein weiterer‘) trägt die Suche, die sie belegt, mit Treffern oder dem Befund ‚0 Treffer‘“; Herkunfts-Anker dann `seit slice-v1-handbuch-forge-cli-nachbesserung`. Benannte Spec-Lücke: keine (die Spezifikation sagte das Standardprofil richtig; die Einschränkung in [`SPEC-051`](../../../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) war eine Präzisierung desselben Absatzes, keine neue Zusage).

**Beobachtungs-Register.** Je ein Beleg dieses Slice, ohne Zähler und ohne Ausgang: [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (F-1), [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (F-2, F-3, F-5, V-3), [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (F-4, V-1, V-3); Dateien `evidence/slice-v1-handbuch-forge-cli-nachbesserung.md`. Nicht aufgetreten: [`BEO-DOK/abhilfe-ohne-wirkungsprobe`](../observations/BEO-DOK/abhilfe-ohne-wirkungsprobe/observation.md) (die Abhilfe hat ihre Wirkungsprobe in T3 vor dem Text), `host-werkzeug-ausserhalb-docker` (kein Host-Werkzeug im Nachweis). **Nicht in diesem Slice eingetragen, sondern mit der Welle-Closure** (Welle-Datei §5, Aufträge an die Welle-Closure; vor deren Schließung wird keine Evidence-Datei `welle-forge-cli-gh-glab.md` angelegt, weil der Vorgang nicht abgeschlossen ist): M-1 für [`BEO-DOK/vertrag-schweigt-zum-nutzerpfad`](../observations/BEO-DOK/vertrag-schweigt-zum-nutzerpfad/observation.md), L-3 für [`BEO-CLI/kennung-in-cli-ausgabe`](../observations/BEO-CLI/kennung-in-cli-ausgabe/observation.md), M-2 für `spec-bedingung-ohne-code-entsprechung`. Lese-Schritt (Schwelle drei): Kein Ausgang wird vergeben; die Einträge über der Schwelle behalten „Entscheidung des Projektinhabers steht aus“.

**Folge-Slices.** [`slice-v1-config-set-validiert-feature-version`](../in-progress/slice-v1-config-set-validiert-feature-version.md) (Ursache M-2, ohne Welle, nicht vor dem Cut; trägt jetzt die Notiz V-3 und dreht T4, streicht den Satz „Zurzeit“). Kein neuer Slice. **Übergabe an `slice-v1-release-cut-v0.15.0`** (Welle-Datei §5): A-1 (`CHANGELOG.md` und Lead nennen die Zusage-Menge Linux `amd64` mit lokaler Engine und dass kein authentifizierter Aufruf gemessen ist); der Satz „`version:` without `source:` ends with exit code `10`“ in `CHANGELOG.md` (Zeile 16 des Blocks `[Unreleased]`) wird auf `generate devcontainer` präzisiert, solange der Folge-Slice nicht in `done/` liegt (der Satz gilt heute nur dort; `config set` lehnt eine schon ungültige Datei ab, schreibt aber `…version` ohne `source` zurzeit selbst). Das Handbuch nimmt der Cut in seinem Stand mit; Frage 5 (Lage des Abschnitts) steht beim Projektinhaber aus.

**Trigger-Audit.** Carveouts: keiner berührt (kein Gate gelockert, kein roter Lauf). Bootstrap-aware Gate: keines berührt. ADR: keine geändert; [`ADR-0035`](../../adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) bleibt Accepted und unberührt. Je Klasse 0 offen.

**Risiken aus §6.** Alle vier tragen ihren Ausgang (§6): zwei eingetreten und behoben, zwei weiter offen mit Anker.

**Drei Paarungen.** (a) Anker: Kein Eintrag trägt das Feld `liegt in`; die Paarung ist leer erfüllt. (b) Folge-Slice: `slice-v1-config-set-validiert-feature-version` liegt als Datei in `open/`; `slice-v1-release-cut-v0.15.0` ist als Datei erst beim Trigger vorgesehen (Welle-Datei §4), kein Versprechen dieses Slice. (c) Register: Die drei genannten Kennungen bestehen als Verzeichnisse mit nicht leerem `evidence/`, jedes mit der neuen Datei dieses Slice.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind `docs/plan/` (DOK, Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006)) und für die Tests T1 bis T4 `internal/hexagon/application` (APP, Greenfield), nur Testdateien. `docs/user/` ist ein öffentlicher Vertrag, keine Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register ([`../observations/`](../observations/README.md)) durchgegangen am 2026-10-23, Zähler = Dateien in `evidence/`:

- [`BEO-DOK/vertrag-schweigt-zum-nutzerpfad`](../observations/BEO-DOK/vertrag-schweigt-zum-nutzerpfad/observation.md) (3×, über der Schwelle, kein Ausgang zugewiesen — Entscheidung des Projektinhabers steht aus): M-1 ist das Muster (der Vertrag sagt „weder zugesagt noch gemessen“, der Nutzer braucht einen Satz). Dieser Slice schließt die Lücke im Text; der Beleg der Welle kommt mit der Welle-Closure, der Slice weist der Beobachtung keinen Ausgang zu.
- [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (20 Belege, verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.10) und [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (47 Belege, verkörpert): gelten für jede Handbuch-Bedingung dieses Slice; angewandt in der Tabelle „Aussage → Beleg“ (§3). Der Halbsatz „`version` ohne `source` weist `u-boot` mit Exit-Code `10` ab“ (Zeilen 2052 und 2053) ist die Klasse, in der der Fund M-2 liegt: wahr für `generate`, ohne Aussage zu `config set`.
- [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (19 Belege) und [`BEO-DOK/abhilfe-ohne-wirkungsprobe`](../observations/BEO-DOK/abhilfe-ohne-wirkungsprobe/observation.md) (1): A4 bis A6 und A10 sind die Fälle (Messvorbehalt nicht kürzen; Wirkungsprobe der Abhilfe vor dem Text).
- Kein Treffer sonst.

Alle berührten Sub-Areas GF.
