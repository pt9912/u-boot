# Slice V1: Handbuch — Standardprofil-Weg zum Token, Einschränkung der `shell`-Sätze, Anmeldung gebündelt und Abhilfe bei ungültiger `version` (Auflagen des Validators zu `gh` und `glab`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-forge-cli-gh-glab`](../welle-forge-cli-gh-glab.md) (Release `v0.15.0`). Fünfter von sechs Slices, **vor** dem Release-Cut (`slice-v1-release-cut-v0.15.0`, Datei erst beim Trigger): Er setzt die Auflagen A-2 und A-3 und die Lücken M-1, M-2 (nur den Handbuch-Satz) und L-1 des Validator-Berichts [`2026-10-23-welle-forge-cli-gh-glab-validierung`](../../../reviews/2026-10-23-welle-forge-cli-gh-glab-validierung.md) um; der Cut nimmt den Stand des Handbuchs mit. Die Welle schließt erst, wenn dieser Slice in `done/` liegt (Welle-Datei §3 Punkt 1, §4).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Katalog-Features `gh`, `glab`), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Token als Umgebungsvariable des Hosts), [`ADR-0035`](../../adr/0035-forge-clis-gh-und-glab-im-devcontainer.md). Der Slice ändert **keine Zusage**: Er schränkt Sätze des Handbuchs auf das Profil ein, in dem sie gemessen sind, und benennt, was ungeprüft ist.

**Berührte Spec-Stellen:** — (das Handbuch ist ein öffentlicher Vertrag, kein Spec-Stratum). Gelesen am 2026-10-23: [`SPEC-051`](../../../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) (Absatz „Wirkung in beiden Profilen“: Standardprofil `devcontainer.json` byte-gleich, ohne `remoteEnv`) und [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (Absätze „Umgebung und Zugangsdaten“ und „Umgebung der Sitzung“: die Namen kommen aus dem Text der `devcontainer.json`) sagen das Standardprofil richtig; **die Spezifikation bleibt unverändert**. Der Fehler liegt allein in zwei Handbuch-Sätzen (§3).

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
- **M-2, nur der Handbuch-Teil (Auflage A-3 des Validators).** Steht in `u-boot.yaml` ein `version` zu `gh` oder `glab` ohne `source`, lehnen `generate devcontainer` (Exit-Code `10`), `doctor` (`uboot.yaml.valid`, Fehler) und jedes weitere `config set` die Datei ab; die Abhilfe (den Eintrag `version` in `u-boot.yaml` von Hand entfernen) nennt keine Meldung und das Handbuch nicht. Die Ursache (`config set …version` prüft nicht wie `generate`) behebt [`slice-v1-config-set-validiert-feature-version`](../open/slice-v1-config-set-validiert-feature-version.md) und nicht dieser Slice.

### Entscheidung zu L-1: Einleitung statt Verschieben (Planner, vorgelegt)

Geprüft am 2026-10-23: Ein Verschieben der Überschrift `#### Forge-CLIs (optional)` aus 4.11 in 4.10 bräche keinen Link, denn alle fünf Verweise sind Überschrift-Anker ohne Nummer (`#forge-clis-optional`, fünf Stellen in `docs/user/benutzerhandbuch.md`: die Katalogzeilen `gh` und `glab` von 4.10, die Egress-Zeilen `gh` und `glab` und der Abschnitt „Zugangsdaten“), und das Inhaltsverzeichnis führt nur die Kapitel der obersten Ebene. Der Slice wählt trotzdem den **Einleitungssatz am Kopf des Abschnitts** (Rezept je Profil), weil der Abschnitt Sandbox-Sachverhalte trägt (Egress, `doctor`-Prüfung, Token-Referenz), die in 4.11 stehen, ein Verschieben von rund 110 Zeilen den Diff und den Review um ein Vielfaches vergrößert und die Auffindbarkeit über die Katalogzeilen in 4.10 schon gegeben ist. **Wünscht der Projektinhaber das Verschieben, ist es ein eigener Handbuch-Slice** (reiner Umbau, keine Zusage).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Ursache von M-2** (`config set …version` meldet Erfolg und sperrt danach das Projekt): Code-Änderung, eigener Slice [`slice-v1-config-set-validiert-feature-version`](../open/slice-v1-config-set-validiert-feature-version.md), ohne Welle und nicht vor dem Cut. Dieser Slice ändert keinen Code.
- **Eine Zusage oder ein Feature** (ein Weg, auf dem `u-boot` das Token im Standardprofil durchreicht, eine `remoteEnv`-Referenz im Standardprofil, ein Befehl zum Anmelden): wäre eine Vertragsänderung ([`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`SPEC-051`](../../../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten): das Standardprofil bleibt ohne `remoteEnv`).
- **Eine Messung mit einem gültigen Token** (`gh auth login` interaktiv, Anmeldung, `gh pr create`): es gibt kein Token für Messungen; was ungeprüft bleibt, steht als „nicht geprüft“ im Text und nicht als Aussage.
- **`CHANGELOG.md` und Lead** (Auflage A-1 des Validators: Zusage-Menge Linux `amd64`, lokale Engine, kein authentifizierter Aufruf gemessen): gehören in den Release-Cut-Slice (Welle-Datei §5).
- **Wortlaute der Meldungen** (L-2: `glab` meldet bei fehlendem Host `lookup gitlab.com … server misbehaving`; L-3: Länge und Wiederholung der Warnungen mit Kennung als Präfix): Code beziehungsweise Werkzeug-Verhalten; notiert in [`slice-v1-config-set-validiert-feature-version`](../open/slice-v1-config-set-validiert-feature-version.md) §6.
- **Das Verschieben des Abschnitts** nach 4.10: siehe Entscheidung zu L-1.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1, Standardprofil (M-1, A-2, A-3-Satz):** Im Abschnitt „Forge-CLIs (optional)“ nennt der Spiegelstrich „Standardprofil, keine Token-Referenz“ den Weg `docker exec -e GH_TOKEN` (und `GITLAB_TOKEN`) als mit einem Platzhalter gemessen, ein Aufruf mit gültigem Token als nicht gemessen; der `shell`-Satz im Spiegelstrich „Token-Variablen“ und der `shell`-Satz in „Zugangsdaten (Git und Agent)“ sind auf das Sandbox-Profil eingeschränkt; ein Satz sagt, dass `gh auth login` und `glab auth login` im Container nicht geprüft sind (kein Aufruf mit Token, Verbleib eines abgelegten Tokens über das Neuanlegen des Containers nicht geprüft). Jede Aussage hat einen Test, einen genannten Beleg oder ist gestrichen (Tabelle „Aussage → Beleg“ in §3).
- [ ] **Liefer-Punkt 2, Anmeldung gebündelt (L-1):** Der Abschnitt beginnt mit einer Einleitung von höchstens zwei Absätzen, die „wie melde ich mich an?“ je Profil beantwortet (Sandbox: Variable auf dem Host setzen, `up` oder `shell`; Standardprofil: `docker exec -e`; ohne Variable: das gemessene Verhalten; Login im Container: nicht geprüft) und auf die Spiegelstriche verweist; die Spiegelstriche bleiben als Belege stehen. Kein Satz der Einleitung ist neu belegt: jeder wiederholt einen Satz aus Liefer-Punkt 1 oder einen schon vorhandenen, gemessenen Spiegelstrich (Fundstelle in §3).
- [ ] **Liefer-Punkt 3, Abhilfe bei ungültiger `version` (M-2, A-3):** Ein Satz im Abschnitt „Forge-CLIs (optional)“ (an der Stelle, die `version` ohne `source` erwähnt, Absatz vor den Spiegelstrichen) nennt: Steht in `u-boot.yaml` ein `version` zu `gh` oder `glab` ohne `source`, lehnen `generate devcontainer`, `doctor` und jedes `config set` die Datei mit Exit-Code `10` beziehungsweise als Fehler ab; die Abhilfe ist, den Eintrag `version` in `u-boot.yaml` von Hand zu entfernen. Die Bedingungen des Satzes sind durch Tests gepinnt (§3, Zeilen T1 bis T4), oder der Teilsatz entfällt.
- [ ] `make gates` grün; `docs-check` ohne Befund; `docs/user/benutzerhandbuch.md` trägt keine `LH-*`-, `ADR-*`-, `SPEC-*`- oder Slice-Kennung und keine Chronik ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 6).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `docs/user/benutzerhandbuch.md` (öffentlicher Vertrag; der einzige berührte neben dem Planning-Baum, siehe §3).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

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
| `spec/spezifikation.md`, `spec/lastenheft.md` | unverändert | gelesen (Zeilen 738 bis 770 und 814, 858 der Spezifikation; Zeile 413 des Lastenhefts): Standardprofil ohne `remoteEnv` und Namen aus dem Text der `devcontainer.json` stehen dort richtig; keine Zeile des Lastenhefts nennt `shell` zusammen mit `GH_TOKEN` |
| `docs/plan/planning/open/README.md`, `docs/plan/planning/welle-forge-cli-gh-glab.md`, `docs/plan/planning/in-progress/roadmap.md` | update | Zeiger auf diesen Slice (Planner; Welle §4 und §5, Roadmap-Zeile der Welle); `make slice-mv` zieht die Zeiger bei den Moves nach |
| `docs/reviews/2026-10-23-slice-v1-handbuch-forge-cli-nachbesserung.md` | neu | Review-Report des Reviewers |

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
| A11 | `config set devcontainer.features.gh.version <Wert>` meldet derzeit Erfolg und schreibt die Datei, bevor die Prüfung greift (dieser Teilsatz hat einen Ablaufpunkt: er entfällt mit [`slice-v1-config-set-validiert-feature-version`](../open/slice-v1-config-set-validiert-feature-version.md)) | T | **T4** (neu, Charakterisierung): `Set` auf einer gültigen Datei ohne `version` für `gh` mit `2.0.0` ohne `source` endet ohne Fehler und schreibt; der Folge-Slice dreht diesen Test auf die Ablehnung. Entscheidet der Implementer, den Teilsatz zu streichen, entfällt T4 und der Satz sagt nur den Zustand der Datei (A7 bis A10) |

Die Tests T1 bis T4 ändern keinen Produktcode; fällt einer rot aus, weil der Code etwas anderes tut als die Aussage, ist die Aussage falsch gelesen, und der Satz wird berichtigt oder gestrichen, nicht der Code geändert.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Der Projektinhaber gibt den Slice frei (Release-Zuordnung `v0.15.0`, vor dem Cut); kein anderer Slice liegt in `in-progress/` (WIP-Limit 1). Der Slice darf vor oder nach [`slice-v1-config-set-validiert-feature-version`](../open/slice-v1-config-set-validiert-feature-version.md) laufen, hängt aber nicht an ihm: Der Satz zu M-2 ist so geschrieben, dass er vor dem Fix wahr ist (A11 mit Ablaufpunkt) und danach nur dort geändert wird.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Implementer will die Überschrift verschieben oder mehr als eine Einleitung umbauen; das ist der Handbuch-Umbau außerhalb dieses Plans (§1).
- `in-progress` → `open` (blockiert): Eine Aussage aus A4 lässt sich weder messen noch streichen, ohne dass der Projektinhaber über den Wortlaut entscheidet.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor (Findings behoben oder mit Ausgang), `make gates` grün, Closure-Notiz mit Lerneintrag geschrieben (§7). Der Release-Cut der Welle startet erst, wenn dieser Slice in `done/` liegt (Welle-Datei §4).

## 6. Risiken und offene Punkte

Jedes Risiko bekommt bei der Closure genau einen Ausgang (eingetreten / entfallen / weiter offen).

- **Aussage ohne Beleg (`BEO-DOK/spec-bedingung-ohne-code-entsprechung`, `BEO-DOK/plan-aussage-ueber-code`):** Der Standardprofil-Weg und der Satz zu `version` sind Handbuch-Bedingungen über Verhalten; der Validator hat sie an einem Binary gemessen, der Plan liest sie am Code. Eine Aussage über das Werkzeug (A4, A5) ist keine über `u-boot`-Code und kann nur mit Messung oder Streichung bestehen. — **Ausgang:** offen bis zur Closure (Tabelle „Aussage → Beleg“ in §3 trägt je Aussage Test, Messung oder Streichung).
- **Wiederholung der Überbehauptung in der Einleitung (L-1):** Eine gebündelte Antwort verleitet dazu, Messvorbehalte zu kürzen und aus „mit Platzhalter gemessen“ ein „funktioniert“ zu machen — derselbe Fehler, der M-1 erzeugte (`BEO-DOK/zusage-weiter-als-messung`). — **Ausgang:** offen bis zur Closure (A6: jeder Satz der Einleitung nennt seinen Spiegelstrich; der Reviewer liest die Einleitung gegen die Spiegelstriche).
- **Der Teilsatz A11 veraltet mit dem Folge-Slice:** Er sagt „derzeit“ und hat einen benannten Ablaufpunkt; geht [`slice-v1-config-set-validiert-feature-version`](../open/slice-v1-config-set-validiert-feature-version.md) zuerst nach `done/`, streicht dieser Slice den Teilsatz selbst (kein T4). — **Ausgang:** offen bis zur Closure.
- **Verschieben des Abschnitts (L-1) bleibt ungeklärt:** Der Plan entscheidet für die Einleitung; der Projektinhaber kann das Verschieben wollen. — **Ausgang:** weiter offen, falls gewünscht: eigener Handbuch-Slice (Frage in der Welle-Datei §5).

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

*(Wird bei der Closure gefüllt: Was hat funktioniert · Was ging anders als geplant · Steering-Loop-Eintrag (Lerneintrag) · Beobachtungs-Register · Folge-Slices · Risiken aus §6 mit je einem Ausgang · Drei Paarungen. Aufträge an die Closure: die Tabelle „Aussage → Beleg“ aus §3 mit dem erreichten Stand; der Befund M-1 ist ein Beleg für `BEO-DOK/vertrag-schweigt-zum-nutzerpfad` und wird mit der Welle-Closure als Beleg `welle-forge-cli-gh-glab` eingetragen, nicht mit diesem Slice, siehe Welle-Datei §5.)*

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
