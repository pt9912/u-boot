# Slice V1: `config set` prüft die `version` eines Dockerfile-Block-Features wie `generate` (Befund M-2 des Validators zu `gh` und `glab`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Folge-Punkt von [`welle-forge-cli-gh-glab`](../done/welle-forge-cli-gh-glab.md) (Befund M-2 des Validator-Berichts [`2026-10-23-welle-forge-cli-gh-glab-validierung`](../../../reviews/2026-10-23-welle-forge-cli-gh-glab-validierung.md)); **nicht vor dem Release-Cut `v0.15.0`**, die Release-Zuordnung trifft der Projektinhaber. Es gibt keine Closure-Bedingung, die von der DoD dieses Slice verschieden wäre. Den Handbuch-Satz mit der Abhilfe liefert vor dem Cut [`slice-v1-handbuch-forge-cli-nachbesserung`](../done/slice-v1-handbuch-forge-cli-nachbesserung.md).

**Bezug:** [`LH-FA-CONF-005`](../../../../spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern) („Beim Setzen muss die geänderte Konfiguration auf Schema-Konformität geprüft werden“), [`LH-FA-CONF-002`](../../../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) (ungültige Werte: fachlicher Validierungsfehler, Exit-Code `10`), [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Katalog), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (Exit-Codes); [`ADR-0035`](../../adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) (Entscheidung 2: `version` ohne `source` ist ein Validierungsfehler) und [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) (dieselbe Regel für `ollama`).

**Berührte Spec-Stellen:** [`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config) (Wertetabelle, Zeile der `version`; Stufe 6 der Schreibpipeline), [`SPEC-051`](../../../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) (Absatz „`version` und `source`“). Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie.

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

**Ziel:** `u-boot config set devcontainer.features.<feature>.version <wert>` weist einen Wert ab, den `generate devcontainer` ablehnt (`version` ohne `source` bei einem Dockerfile-Block-Feature: `ollama`, `gh`, `glab`), mit Exit-Code `10` und ohne die Datei zu ändern — statt Erfolg zu melden und das Projekt zu sperren.

### Befund (Code gelesen am 2026-10-23, Validator-Bericht M-2)

- **Wirkung.** `u-boot config set devcontainer.features.gh.version 2.0.0` meldet Erfolg und schreibt die Datei (Validator, ein Lauf). Danach lehnen `generate devcontainer` (Exit-Code `10`), `doctor` (`uboot.yaml.valid`, Fehler) und **jedes weitere** `config set` die Datei ab; ein `config unset` gibt es nicht; die Meldung von `generate` beginnt mit „managed block missing or malformed:“, was nicht stimmt. Abhilfe ist nur das Bearbeiten von `u-boot.yaml` von Hand.
- **Ursache im Code.** Die Schreibpipeline von `config set` hat zwei Prüfstellen: Beim Lesen ruft `readUbootYAMLBody` (`internal/hexagon/application/config.go`, Zeilen 652 bis 681) `validateDevcontainer` und damit `validateDevcontainerFeatures` (`devcontainer_features.go`, Zeilen 613 bis 647; die Regel in den Zeilen 641 bis 644: `entry.Source == "" && entry.Version != "" && featureCatalogue()[k].dockerfileBlock` → `ErrInvalidFeatureSource`). Nach dem Patch prüft `revalidateConfigDomain` (`config.go`, Zeilen 488 bis 539) den geschriebenen Wert, und für `.version` prüft `revalidateFeatureEntry` (Zeilen 576 bis 581) nur `entry.Version == ""`. Die Regel von `generate` (`readProjectConfig`, `generate.go`, Zeilen 279 bis 293, ebenfalls `validateDevcontainer`) läuft beim Schreiben nicht: Erst das **nächste** Lesen lehnt die Datei ab. Dieselbe Lücke hat `SetMany` (`config_many.go`, Zeile 184: `revalidateConfigDomain` je Paar) und gilt für `ollama` ebenso wie für `gh` und `glab`, weil `dockerfileBlock` alle drei trägt.
- **Vertrag.** Die Zusage besteht schon: [`LH-FA-CONF-005`](../../../../spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern) („Beim Setzen muss die geänderte Konfiguration auf Schema-Konformität geprüft werden“), [`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config) Stufe 5 („Das Ergebnis muss sich wieder als Projektkonfiguration einlesen lassen (sonst Schemafehler, Exit-Code `10`)“) und das Handbuch (Abschnitt „Werte ändern“: „Jeder Aufruf wird gegen das Schema geprüft, **bevor** geschrieben wird. Ungültige Werte führen zu Exit-Code 10, ohne die Datei anzufassen“). Der Code hält Stufe 5 für diesen Fall nicht: Er liest das Ergebnis nur durch `yaml.Unmarshal` zurück, nicht durch `validateDevcontainer`. Das ist ein **Fehler gegen einen bestehenden Vertrag**, kein neuer.

### Lastenheft: keine Änderung, keine Freigabe nötig (Planner, vorgelegt)

Gelesen am 2026-10-23: [`LH-FA-CONF-005`](../../../../spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern) und [`LH-FA-CONF-002`](../../../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) tragen die Zusage schon; [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) nennt weder `version` noch `config set` zu `version`. Der Slice **schreibt keine Zeile in `spec/lastenheft.md`**. Eine Entscheidung des Projektinhabers steht für diesen Slice nicht an; die Wahl der Sentinel und der Diagnose-Zuordnung ist die des Architect (Empfehlung in §3), und ob der Wortlaut der Fehlermeldung von `generate` („managed block missing or malformed“ für einen Schemafehler) Gegenstand wird, ist ein eigener Vorgang (§6). Hielte der Architect beim Lesen eine Lastenheft-Zeile für nötig (etwa weil `config set … version` eine neue, vom Lastenheft nicht getragene Ablehnung bekäme), ist das eine **Entscheidung des Projektinhabers** und stoppt den Slice (Rückführung in §4).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein `config unset`** (die Datei lässt sich dann ohne Handarbeit berichtigen): neuer Befehl, neuer Vertrag (Befehlstabelle, [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle)); braucht eine Entscheidung des Projektinhabers. Der Fix verhindert die Sperre; den Weg heraus für Dateien, die schon gesperrt sind, bleibt die Handarbeit (Handbuch-Satz aus dem Schwester-Slice).
- **Der Wortlaut „managed block missing or malformed“ bei einem Schemafehler** (`ErrGenerateManualConflict`, `internal/hexagon/port/driving/generate.go`, Zeile 199; wird in `generate.go` Zeilen 291 und 989 für Schemafehler und in `cli.go` Zeile 402 für Exit-Code `10` genutzt): geteilter Sentinel mit mehreren Verwendungen; eine Wortlautänderung berührt alle Aufrufer und gehört in einen eigenen Vorgang (§6).
- **Meldungen und Hinweise von `generate` und `doctor`, die die Abhilfe nennen** (`checkUbootYaml` trägt den Hint „Fix the offending entry (URL scheme, host, or feature name)“, `doctor.go`, Zeilen 617 bis 625, ohne den Fall `version`): Wortlaut; Kandidat für den Architect, wenn die Prüfung ohnehin berührt wird, sonst Notiz in §6.
- **Die Prüfungen anderer Pfade** (`devcontainer.profile`, `sandbox.*`, `source`): prüfen beim Schreiben schon (`validateDevcontainerSandbox`, Allow-Liste); der Slice berührt nur die Zeile der `version`.
- **L-2 und L-3 des Validators** (`glab` ohne Hinweis auf die Egress-Liste; Länge und Wiederholung der Warnungen): keine Änderung, nur Notiz in §6.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **Liefer-Punkt 1, Code** (`5c9a312c`)**:** [`LH-FA-CONF-005`](../../../../spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern) und [`LH-FA-CONF-002`](../../../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) erfüllt: `config set` und `config set` mit mehreren Paaren lehnen `devcontainer.features.<feature>.version` für `ollama`, `gh` und `glab` ohne `source` mit Exit-Code `10` ab, nach demselben Maßstab wie `generate devcontainer`; die Datei bleibt byte-identisch (auch unter `--dry-run` und `--diff` kommt der Fehler). Mit gesetzter `source` (und freigegebener Quelle) gelingt das Setzen, in einem Aufruf mit mehreren Paaren in der Reihenfolge `source`, dann `version`. Ein Test färbt die Ablehnung rot, wenn die Prüfung fehlt (Mutation im Review nachweisbar); der Test T4 des Schwester-Slice ([`slice-v1-handbuch-forge-cli-nachbesserung`](../done/slice-v1-handbuch-forge-cli-nachbesserung.md)), falls er schon in `done/` liegt, ist auf die Ablehnung gedreht.
- [x] **Liefer-Punkt 2, Doku** (`00e3ca3d`)**:** [`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config) (Wertetabelle, Zeile der `version`, und Stufe 6 der Schreibpipeline) und [`SPEC-051`](../../../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) (Absatz „`version` und `source`“: gilt auch beim Schreiben mit `config set`) sagen das Verhalten; das Handbuch nimmt den Teilsatz „`config set` … meldet derzeit Erfolg“ (A11 des Schwester-Slice) zurück und nennt in der Pfad-Tabelle (Zeile `devcontainer.features.<name>.version`) den Fall; jede Bedingung dieser Texte hat einen pinnenden Test (Tabelle „Aussage → Test“ in §3) oder ist gestrichen. Lastenheft unverändert.
- [x] `make gates` grün (inkl. `docs-check`, `doc-trace`); `make test-docker` nicht nötig (kein Docker-Pfad berührt), das Fehlen dieses Laufs steht in der Closure-Notiz.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: `CHANGELOG.md` (`[Unreleased]`, Eintrag *Fixed*), Handbuch, Spezifikation (siehe §3).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht. Jede Aussage über Dateien, Funktionen und Verhalten wurde am 2026-10-23 an der genannten Stelle gelesen ([`AGENTS.md`](../../../../AGENTS.md) §3.10); der Architect und der Implementer erweitern die Tabelle um jede weitere Aussage samt Fundstelle und gleichen sie nach der letzten Fix-Runde gegen `git diff --stat` ab.

| Datei / Komponente | Änderungs-Art | Begründung (Fundstelle) |
|---|---|---|
| `internal/hexagon/application/config.go` | update | `revalidateFeatureEntry`, Fall `domain.ConfigDevcontainerFeatureVersion` (Zeilen 576 bis 581): neben `entry.Version == ""` prüft der Fall dieselbe Regel wie `validateDevcontainerFeatures` (`devcontainer_features.go`, Zeilen 641 bis 644: Dockerfile-Block-Eintrag mit `Version` und ohne `Source`). Die Prüfung steht an **einer** Stelle für `Set` (`config.go`, Zeile 266) und `SetMany` (`config_many.go`, Zeile 184, über `revalidateConfigDomain`); eine Doppelung der Regel ist nicht erlaubt (der Architect wählt: kleine Prüffunktion neben `validateDevcontainerFeatures`, die beide aufrufen, oder Aufruf von `validateDevcontainerFeatures` auf dem gepatchten Teilbaum). Doc-Kommentar von `revalidateFeatureEntry` (Zeilen 541 bis 551, nennt bisher nur `.source`) nachziehen |
| `internal/hexagon/application/config_features_test.go` | update | Tests: (1) `Set` `devcontainer.features.gh.version` ohne `source` → Fehler, der `driving.ErrConfigValueInvalid` wickelt (Empfehlung, siehe unten), Datei byte-identisch; dieselben Fälle für `glab` und `ollama`; (2) `Set` mit `source` und freigegebener Quelle gelingt; (3) ein Feature ohne Dockerfile-Block (`java.version`, Zeile 227 pinnt den leeren Wert) bleibt setzbar. Bestehender Test als Vorlage: `TestConfigSet_FeatureVersion_EmptyRejected` (Zeile 227) |
| `internal/hexagon/application/config_many_test.go` | update | `SetMany`: `gh.version` allein → Fehler, Datei unverändert (atomar, `config_many.go` Zeile 100 bis 105); `gh.source` dann `gh.version` im selben Aufruf gelingt (Reihenfolge der Paare, Handbuch „Werte ändern“); umgekehrte Reihenfolge schlägt fehl |
| `internal/hexagon/application/devcontainer_features.go` | update | Umsetzung (Implementer): `featureVersionNeedsSource` ist die eine Aussage der Regel; `validateDevcontainerFeatures` (Zeile 641) und `revalidateFeatureEntry` fragen sie |
| `docs/plan/planning/in-progress/roadmap.md`, `docs/plan/planning/open/README.md` | update | Zustandsmarker beim Start: „Nichts in Arbeit.“ (Roadmap, Sektion „Offene Wellen“; `docs-check` verlangt die Benennung) und „nur Plan“ → „in Arbeit“ in der Slice-Zeile |
| `internal/adapter/driving/cli/config_test.go` (oder `config_acceptance_test.go`) | unverändert (Umsetzung) | gelesen: ohne neue Sentinel bleiben beide Klassifikatoren und ihre Pins unverändert: `ErrConfigValueInvalid` → 10 (`cli_test.go` Zeile 987), → [`LH-FA-CONF-001`](../../../../spec/lastenheft.md#lh-fa-conf-001--projektkonfiguration) (`config_internal_test.go` Zeile 31). Ursprünglicher Plan-Eintrag: Exit-Code `10` und Diagnose-Code am Adapter: `ExitCode(err)` (`cli.go`, `isConfigValidationError`, Zeilen 453 bis 465, listet `ErrConfigValueInvalid`, `ErrConfigSchemaInvalid`, `ErrConfigPostPatchSanityFailed`) und `mapConfigErrorToDiagnostic` (`config.go` des Adapters, Zeilen 520 bis 546); bestehende Pins: `cli_test.go` Zeilen 986 bis 988 |
| `internal/hexagon/application/generate_forge_cli_test.go`, Docker-Tests | unverändert | gelesen: `TestGenerateDevcontainer_ForgeCLI_VersionOverrideWithoutSourceRejected` (Zeile 443) pinnt `generate`; der Slice ändert `generate` nicht |
| `spec/spezifikation.md` | update | [`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config): Zeile der Wertetabelle `devcontainer.features.<feature>.version` (Zeile 518: bisher „nicht-leerer String“) um die Bedingung für Dockerfile-Block-Features erweitern; Stufe 6 der Schreibpipeline (Zeile 537: bisher nur `.source` genannt) um `.version`; [`SPEC-051`](../../../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) Absatz „`version` und `source`“ (Zeile 750: bisher nur `generate devcontainer`) um `config set`; jede Bedingung mit pinnendem Test (Tabelle unten) |
| `docs/user/benutzerhandbuch.md` | update | Zeile der Pfad-Tabelle `devcontainer.features.<name>.version` (Zeile 2370: „zum Beispiel `21`“); der Halbsatz zur Ablehnung im Abschnitt „Forge-CLIs (optional)“ (Zeilen 2052 und 2053) bleibt wahr und nennt `config set` mit; der Teilsatz A11 aus dem Schwester-Slice entfällt (Fundstelle nach dessen Closure) |
| `CHANGELOG.md` | update | `[Unreleased]` (Zeilen 12 bis 21 tragen *Added* und *Changed*; *Fixed* kommt in älteren Blöcken vor, Zählung 3): Eintrag *Fixed* — `u-boot config set` weist eine `version` ohne `source` für `ollama`, `gh` und `glab` mit Exit-Code `10` ab, statt das Projekt zu sperren |
| `spec/lastenheft.md` | unverändert | gelesen: [`LH-FA-CONF-002`](../../../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) (Zeilen 824 bis 839: „ungültige Werte führen zu einem fachlichen Validierungsfehler (Exit-Code `10`)“) und [`LH-FA-CONF-005`](../../../../spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern) (Zeilen 849 bis 853: „Beim Setzen muss die geänderte Konfiguration auf Schema-Konformität geprüft werden“) tragen die Zusage; [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) nennt `version` bei Dockerfile-Block-Features nicht |
| `spec/architecture.md` | unverändert | gelesen: §7.2 „Zwei Regeln, die die Klassifikation zusammenhalten“ (Zeilen 640 bis 658); ohne neue Sentinel ändert sich nichts. Wählt der Architect eine neue Sentinel, wandert sie in beide Klassifikatoren (siehe Ansatz) |
| `docs/reviews/2026-10-24-slice-v1-config-set-validiert-feature-version.md` | neu | Review-Report (der Dateiname trägt das Datum des Laufs, 2026-10-24) |
| `docs/reviews/2026-10-24-slice-v1-config-set-validiert-feature-version-verify.md` | neu | Verify-Report (Verifier, schließbar) |
| `docs/plan/planning/done/slice-v1-handbuch-forge-cli-nachbesserung.md`, `docs/plan/planning/done/slice-v1-release-cut-v0.15.0.md`, `docs/plan/planning/done/welle-forge-cli-gh-glab.md` | update (nur Verweise) | Verweis-Nachzug von `make slice-mv` bei den zwei Lifecycle-Wechseln dieses Slice (`open/` → `in-progress/` → `done/`): nur Pfade in Link-Zeilen, kein Inhalt (gelesen: `git show 54aba114`, 28 geänderte Zeilen in den drei Dateien, jede trägt die Kennung dieses Slice) |
| `docs/plan/planning/observations/BEO-CLI/kennung-in-cli-ausgabe/evidence/welle-forge-cli-gh-glab.md`, `docs/plan/planning/observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/evidence/welle-forge-cli-gh-glab.md` | update (nur Verweise) | Evidence-Dateien der Welle-Closure, die diesen Slice verlinken: derselbe Verweis-Nachzug, je eine Zeile |
| `docs/plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/` (`evidence/slice-v1-config-set-validiert-feature-version.md`, `state.md`), `BEO-WZG/host-werkzeug-ausserhalb-docker/` (`evidence/slice-v1-config-set-validiert-feature-version.md`, `state.md`) | neu bzw. update | Closure: je ein Beleg dieses Slice (LOW-1 des Reviews; Selbstauskunft zu Heredoc-Commits), Stand-Absatz ohne Zähler und ohne Ausgang |
| `docs/plan/planning/observations/BEO-APP/pruefung-nur-beim-lesen/` (`observation.md`, `state.md`, `evidence/slice-v1-config-set-validiert-feature-version.md`) | neu | Closure: neue Beobachtung, erster Vorgang |
| `docs/plan/planning/in-progress/roadmap.md`, `docs/plan/planning/open/README.md` | update (Closure) | Ruhe-Marker „Nichts in Arbeit.“ und Zustand der Slice-Zeile bei der Closure; die Datei des Slice selbst wandert per `make slice-mv` nach `done/` |

**Ansatz (Empfehlung des Planners, der Architect entscheidet):**

- **Eine Prüfung, zwei Aufrufer.** Die Regel „Dockerfile-Block-Feature mit `version` und ohne `source`“ steht heute in `validateDevcontainerFeatures` (Zeilen 641 bis 644). `revalidateFeatureEntry` ruft sie für den Eintrag auf, den der Patch geändert hat; kein zweiter Text der Regel.
- **Sentinel: keine neue.** Die Ablehnung wickelt `driving.ErrConfigValueInvalid`: Der Nutzer hat einen Wert angegeben, den er ändern kann (wie der leere Wert in Zeile 227 des Tests und die Allow-Prüfung der `source`, `config.go` Zeilen 570 bis 574). Das Ergebnis ist Exit-Code `10` über `isConfigValidationError`, Diagnose-Code über `mapConfigErrorToDiagnostic` (`ErrConfigValueInvalid` → [`LH-FA-CONF-001`](../../../../spec/lastenheft.md#lh-fa-conf-001--projektkonfiguration)). Die Alternative `ErrConfigSchemaInvalid` (wie `readUbootYAMLBody`) meldet „Datei kaputt“ und trifft den Nutzer, der nur einen Wert gesetzt hat, falsch. **Fordert der Architect eine neue Sentinel**, gilt die Dual-Classifier-Regel (`architecture.md` §7.2): Eintrag in `isConfigValidationError` **und** in `mapConfigErrorToDiagnostic`, dazu der `ExitCode`-Pin (`cli_test.go` Zeilen 986 bis 988) und ein Diagnose-Pin.
- **Meldung.** Der Text nennt den Pfad und was zu tun ist („set `devcontainer.features.<feature>.source` first, or remove the version“); `u-boot config set` kennt `ConfigHintError` (`driving.NewConfigHintError`, `config.go` Zeilen 363 bis 366) für einen Reparaturbefehl in `data.hint` (`--json`). Ob er hier gesetzt wird, entscheidet der Architect; der Hinweis nennt nur Befehle, die wirken (Wirkungsprobe vor dem Text, `BEO-DOK/abhilfe-ohne-wirkungsprobe`).
- **Reihenfolge mehrerer Paare.** `SetMany` wendet die Paare der Reihe nach auf ein Dokument an (`config_many.go`, Zeilen 126 bis 138) und prüft nach jedem Paar; `source` vor `version` ist die Ordnung, die das Handbuch für Allow-Liste und Quelle schon nennt („ein späteres Paar darf auf ein früheres aufbauen“). Eine Prüfung erst am Ende aller Paare wäre eine Änderung dieser Zusage und gehört nicht in diesen Slice.

**Aussage → Test** (Spezifikation und Handbuch; Klasse **T** = Test, **G** = gestrichen):

| # | Aussage | Klasse | Test |
|---|---|---|---|
| B1 | `config set …<feature>.version` ohne `source` bei `ollama`, `gh`, `glab`: Exit-Code `10`, Datei unverändert | T | Test (1) in `config_features_test.go`; `ExitCode`-Pin am Adapter |
| B2 | …auch bei mehreren Paaren in einem Aufruf (atomar) | T | Test in `config_many_test.go` |
| B3 | …auch unter `--dry-run` und `--diff` | T | Test am Adapter oder Service mit `PreviewMode` (Stufe 5 liegt vor der Wahl des Dateisystems, `config.go` Zeilen 266 bis 276); sonst Streichung der Aussage |
| B4 | mit `source` (und freigegebener Quelle) gelingt das Setzen | T | Test (2) und Test in `config_many_test.go` |
| B5 | andere Features (ohne Dockerfile-Block) behalten ihr `version` | T | Test (3) |
| B6 | Stufe 6 der Schreibpipeline nennt `.version` | T | B1 (der Satz folgt dem Test) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Der Release `v0.15.0` ist geschnitten (`slice-v1-release-cut-v0.15.0`, Datei erst beim Trigger der Welle [`welle-forge-cli-gh-glab`](../done/welle-forge-cli-gh-glab.md); beobachtbar: der Slice liegt in `done/`, Tag `v0.15.0` gesetzt) **und** der Projektinhaber gibt den Slice frei (Release-Zuordnung); kein anderer Slice liegt in `in-progress/` (WIP-Limit 1). Der Schwester-Slice [`slice-v1-handbuch-forge-cli-nachbesserung`](../done/slice-v1-handbuch-forge-cli-nachbesserung.md) liegt dann in `done/` (er ist Teil der Welle und läuft vor dem Cut).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Prüfung lässt sich nicht an einer Stelle für `Set` und `SetMany` führen (zwei Prüf-Pfade); dann trennt der Architect die Mehrpaar-Variante als eigenen Slice.
- `in-progress` → `open` (blockiert): Eine Lastenheft-Zeile ist nötig (Entscheidung des Projektinhabers, §1), oder der Projektinhaber will `config unset`/einen anderen Weg aus der gesperrten Datei zuerst entscheiden.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor (Findings behoben oder mit Ausgang), `make gates` grün, Test-Mutation gezeigt (die Ablehnung entfernt → der Test wird rot), Closure-Notiz mit Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

Jedes Risiko bekommt bei der Closure genau einen Ausgang (eingetreten / entfallen / weiter offen).

- **Die Prüfung greift zu früh bei mehreren Paaren:** Prüft `SetMany` nach jedem Paar, scheitert `version` vor `source`, wo ein Nutzer das Gegenteil erwartet. — **Ausgang:** nicht eingetreten: `TestConfigSetMany_DockerfileBlockVersionNeedsSourceFirst` pinnt `source` vor `version` (gelingt) und die umgekehrte Reihenfolge (Exit-Code `10`, Datei unverändert); die Mutation M5 des Verifiers färbt ihn rot, am Binary gemessen.
- **Ein bereits gesperrtes Projekt bleibt gesperrt:** Der Fix verhindert das Entstehen, nicht die Folge für Dateien, die ein früheres `config set` schon beschrieb. — **Ausgang:** weiter offen: Eine von Hand oder früher gesperrte Datei bleibt gesperrt (am Binary gemessen: `generate devcontainer` Exit-Code `10`, `config set` Exit-Code `10`, `doctor` Exit-Code `11`). Der Handbuch-Satz nennt die Abhilfe von Hand; `config unset` wäre ein neuer Befehl und braucht eine Entscheidung des Projektinhabers ([`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle)). Kein Slice; diese Notiz verspricht ihn nicht als Datei.
- **Sentinel/Diagnose-Zuordnung driftet (Dual-Classifier-Regel):** Eine neue Sentinel nur in einem Klassifikator fällt auf Exit-Code `1`. — **Ausgang:** entfallen: Es entstand keine neue Sentinel; die Ablehnung wickelt `ErrConfigValueInvalid`, beide Klassifikatoren und ihre bestehenden Pins sind unverändert (Review und Verify gelesen; Exit-Code `10` am Binary gemessen).
- **Notiz L-2 (Validator, LOW, nicht geändert):** Fehlt `gitlab.com` in der Egress-Liste, meldet `glab` `lookup gitlab.com on 127.0.0.1:53: server misbehaving` (Verhalten des Werkzeugs, `u-boot` kann es nicht ändern; das Handbuch nennt die Meldung im Abschnitt „Forge-CLIs (optional)“, Spiegelstrich „Egress“, ohne die Abhilfe in einem Satz). Ein Satz „bei `lookup … misbehaving` den Host in `devcontainer.sandbox.egress.allow` eintragen“ stünde dem Nutzer bei — Text, kein Code. — **Ausgang:** weiter offen: kein Gegenstand dieses Slice; wird mit der Welle-Closure im Register bewertet (Welle-Datei §5); ein Handbuch-Satz ist ein eigener kleiner Vorgang, wenn der Projektinhaber ihn will.
- **Notiz L-3 (Validator, LOW, nicht geändert):** Die zwei Warnungen von `generate devcontainer` zu `gh` und `glab` (rund 400 Zeichen, Präfix „Warning ([`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte))“) erscheinen bei jedem Lauf, auch ohne Änderung; das ist zugesagt (Handbuch: „je Feature eine Meldung in jedem Profil“), nur laut. Die Kennung in der Ausgabe ist die Klasse von [`BEO-CLI/kennung-in-cli-ausgabe`](../observations/BEO-CLI/kennung-in-cli-ausgabe/observation.md) (3×, über der Schwelle, Ausgang beim Projektinhaber: Produktfrage (j)). — **Ausgang:** weiter offen → `BEO-CLI/kennung-in-cli-ausgabe` (Beleg mit der Welle-Closure, Welle-Datei §5).
- **Notiz V-3 (Verifier des Schwester-Slice, INFO, gelesen am 2026-10-23):** Das Handbuch sagt, bei `version` ohne `source` lehne „jedes weitere `config set`“ die Datei ab. `SetMany` (die Mehrpfad-Form, `config_many.go` Zeile 106, aufgerufen in `cli/config.go` Zeile 666) liest die Datei über dieselbe Eingangsprüfung (`readUbootYAMLBody`, Zeile 122) und lehnt sie ebenso ab; der Test `TestConfigSet_ForgeVersionWithoutSource_RejectsEverySet` pinnt nur `Set`. — **Ausgang:** entfallen: `TestConfigSetMany_RejectsFileWithForgeVersionWithoutSource` pinnt die Mehrpfad-Form auf der schon ungültigen Datei.
- **Notiz zur Meldung von `generate`:** „managed block missing or malformed“ bei einem Schemafehler der `devcontainer:`-Teilbaums (Sentinel `ErrGenerateManualConflict`, `generate.go` Zeilen 291 und 989) trifft die Ursache nicht; ein eigener Wortlaut wäre eine Änderung des geteilten Sentinel-Textes. — **Ausgang:** weiter offen: eigener Vorgang, Entscheidung des Projektinhabers; im Register erst, wenn sie ein zweites Mal auftritt.

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

Stand der Closure: 2026-10-24. Quellen: Review-Report [`2026-10-24-slice-v1-config-set-validiert-feature-version`](../../../reviews/2026-10-24-slice-v1-config-set-validiert-feature-version.md) (0 HIGH, 0 MEDIUM, 1 LOW, 2 INFO), Verify-Report [`2026-10-24-slice-v1-config-set-validiert-feature-version-verify`](../../../reviews/2026-10-24-slice-v1-config-set-validiert-feature-version-verify.md) (schließbar, `3f5013d7`; `make gates` Exit-Code 0 selbst gefahren, `make doc-trace` 84 Anforderungen und 0 Waisen, Mutationen M1 bis M5, Ende-zu-Ende am gebauten Binary).

**Geliefert.** `config set` (Einzel- und Mehrpaar-Form) weist `devcontainer.features.<feature>.version` für `ollama`, `gh` und `glab` ohne `source` mit Exit-Code `10` ab, auch unter `--dry-run` und `--diff`, die Datei bleibt byte-identisch (Code `5c9a312c`). Die Regel steht an **einer** Stelle, `featureVersionNeedsSource` (`devcontainer_features.go`); `validateDevcontainerFeatures` (Lesepfad, unverändert im Text und in der Sentinel `ErrInvalidFeatureSource`) und `revalidateFeatureEntry` (Schreibpfad, `config.go`) fragen sie. Fehlerklasse `ErrConfigValueInvalid`, keine neue Sentinel, beide Klassifikatoren unverändert. Mit `source` gelingt das Setzen, in einem Aufruf mit mehreren Paaren in der Reihenfolge `source`, dann `version`; Features ohne Dockerfile-Block (`java`) behalten ihr `version`. Spezifikation ([`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config) Tabelle und Stufe 6, [`SPEC-051`](../../../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) Absatz „`version` und `source`“), Handbuch (Teilsatz „meldet derzeit Erfolg“ zurückgenommen, Pfad-Tabelle) und `CHANGELOG.md` `[Unreleased]` (*Fixed*) nachgezogen (`00e3ca3d`); das Lastenheft ist unverändert. Der Charakterisierungstest des Schwester-Slice (`…SetItselfSucceeds`) ist auf die Ablehnung gedreht (`…Rejected`).

**Was funktionierte.** Die Tabelle „Aussage → Test“ vor dem Text: Jede Bedingung in Spezifikation und Handbuch hat einen pinnenden Test (B1 bis B6), und die Mutationen von Reviewer (drei) und Verifier (fünf: Prüfung aus, Aufruf in `SetMany` entfernt, `dockerfileBlock`-Bedingung entfernt, `Source == ""` entfernt, Prüfung nicht pro Paar) färbten je mindestens einen Test rot. Die Empfehlung „keine neue Sentinel“ ersparte die Dual-Classifier-Arbeit (Pins tragen). Die Messung am gebauten Binary (Docker, ohne Host-Toolchain) zeigte `--dry-run`, `--diff`, `java`, drei Paare in einem Aufruf, die umgekehrte Reihenfolge und die gesperrte Bestandsdatei (`generate` Exit-Code `10`, `config set` Exit-Code `10`, `doctor` Exit-Code `11`).

**Was anders lief.**

- **LOW-1 des Reviews** (bestätigt vom Verifier): Die Plan-Tabelle nannte die Dateien des Verweis-Nachzugs unter `done/` (drei Dateien), die zwei Evidence-Dateien des Registers und die Reports nicht, und der Name des Review-Reports trug das Datum 2026-10-23 statt 2026-10-24. Berichtigt in der Closure: Die Tabelle führt jetzt jede Datei des Diffs (Abgleich gegen `git diff --stat` nach der letzten Änderung, siehe Paarungen).
- **Erster `make gates`-Lauf rot** (Auskunft des Implementer-Kontexts, nicht in den Reports): ein Fehler am Roadmap-Marker bzw. Ruhe-Marker (Einzelheit nicht in den Reports); behoben vor dem Review, der Stand in den Reports ist grün.
- **Heredoc beim Commit** (Auskunft des Projektleiter-Kontexts: der Verifier schrieb eine Commit-Message per Heredoc statt per Datei): Prozessabweichung ohne Wirkung auf Dateien und Nachweise; als Beleg im Register geführt.
- **Der Plan-Eintrag „`ExitCode`-Pin am Adapter“ (B1) wurde nicht neu angelegt:** Der bestehende Pin `ErrConfigValueInvalid` → Exit-Code `10` und der Diagnose-Pin tragen die Kette (Plan-Zeile `cli/config_test.go` führt beides); der Verifier belegt sie am Binary.

**Offene Punkte der Verifikation, ehrlich:**

- Es gibt **keinen automatischen Test**, der den CLI-Befehl `config set … --dry-run` oder `--diff` durchläuft. Die Zusage ist auf Service-Ebene gepinnt (`PreviewDryRun`, `PreviewAndApply`, M1 färbt rot) und **einmal** am Binary gemessen; das ersetzt keinen Regressionstest (INFO-2 des Reviews, INFO des Verifiers).
- Die `--json`-Ausgabe der Ablehnung (`data.hint`) ist **nicht geprüft**; der Slice sieht keinen Hinweis vor.
- Nicht gelaufen: `make test-docker` (kein Docker-Pfad berührt), `make verify-depguard` (keine Importregel berührt, keine neuen Importe in `internal/hexagon/application`), `make ci` mit `govulncheck` und `image-scan` (nicht Teil der DoD). Mac, Podman und andere Hosts sind nicht gemessen; der Fix liegt in der Anwendungsschicht ohne Plattformpfad.

**Steering-Loop-Eintrag (Lerneintrag).** Gelernt: Eine Regel, die zwei Pfade prüfen (Lesen beim nächsten Aufruf, Schreiben beim Setzen), gilt für beide nur, wenn beide **dieselbe Prüffunktion** rufen. Stufe 5 der Schreibpipeline sagte „wieder einlesen lassen“, der Code las das Ergebnis nur durch `yaml.Unmarshal` zurück; die Regel der Lesestelle (`validateDevcontainerFeatures`) lief beim Schreiben nicht, und erst das nächste Lesen sperrte das Projekt. Der Fix hielt die Regel **einmal** (`featureVersionNeedsSource`) und ließ beide Aufrufer fragen; ein zweiter Text der Regel wäre beim nächsten Feld wieder gedriftet. Das ist der Auftrag der Planung: Die Klasse ist als erster Vorgang im Register ([`BEO-APP/pruefung-nur-beim-lesen`](../observations/BEO-APP/pruefung-nur-beim-lesen/observation.md)); **keine Verkörperung in diesem Slice** (kein Feld `liegt in`), die Schwelle von drei ist nicht erreicht. **Benannte Spec-Lücke:** keine (Stufe 5 der Spezifikation und die Zusage des Lastenhefts sagten das Verhalten schon; der Code hielt sie nicht).

**Beobachtungs-Register.** Je ein Beleg dieses Slice, ohne Zähler und ohne Ausgang (Dateien `evidence/slice-v1-config-set-validiert-feature-version.md`): [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (LOW-1: Plan-Tabelle ohne die Dateien des Verweis-Nachzugs und der Evidence, Dateiname des Reports; Rückfall nach der Verkörperung in `AGENTS.md` §3.10), [`BEO-WZG/host-werkzeug-ausserhalb-docker`](../observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md) (Selbstauskunft zum Heredoc-Commit), und das neue Verzeichnis [`BEO-APP/pruefung-nur-beim-lesen`](../observations/BEO-APP/pruefung-nur-beim-lesen/observation.md) (erster Vorgang). **Nicht aufgetreten:** [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (jede bedingte Aussage in Spezifikation und Handbuch hat einen pinnenden Test; der Befund M-2 selbst ist schon als Beleg der Welle-Closure von [`welle-forge-cli-gh-glab`](../done/welle-forge-cli-gh-glab.md) gezählt und zählt nicht ein zweites Mal). Lese-Schritt (Schwelle drei): Kein Eintrag dieses Slice erreicht sie neu; die Einträge über der Schwelle behalten ihren Stand, kein Ausgang wird vergeben.

**Folge-Slices.** Keine. Offene Fragen ohne Slice, Entscheidung des Projektinhabers: `config unset` (§6), der Wortlaut „managed block missing or malformed“ bei einem Schemafehler von `generate`, ein Adapter-Test des CLI-Befehls mit `--dry-run` (siehe Offene Punkte). **Release-Zuordnung:** Der Fix steht in `CHANGELOG.md` `[Unreleased]`; welchem Release er zugeordnet wird, entscheidet der Projektinhaber, hier nicht festgelegt. `v0.15.0` ist veröffentlicht, der Tag zeigt weiter auf `a895ffd2`.

**Trigger-Audit.** Carveouts: keiner berührt (kein Gate gelockert, kein roter Lauf in der Closure). Bootstrap-aware Gate: keines berührt. ADR: keine geändert; [`ADR-0035`](../../adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) und [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) bleiben Accepted und unberührt. Je Klasse 0 offen.

**Risiken aus §6.** Alle tragen ihren Ausgang: zwei entfallen (Sentinel-Drift, V-3), eins nicht eingetreten (Reihenfolge), vier weiter offen mit Anker (gesperrte Bestandsdatei, L-2, L-3, Meldung von `generate`).

**Paarungen.** (a) Anker: kein Eintrag trägt das Feld `liegt in`, nichts zu prüfen. (b) Folge-Slice: keiner genannt. (c) Register: `BEO-APP/pruefung-nur-beim-lesen`, `BEO-DOK/plan-aussage-ueber-code`, `BEO-WZG/host-werkzeug-ausserhalb-docker` und `BEO-DOK/spec-bedingung-ohne-code-entsprechung` bestehen als Verzeichnisse mit nicht leerem `evidence/`.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind `internal/hexagon/application` (APP, Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006)), `internal/adapter/driving/cli` (CLI, Greenfield; nur Tests, solange keine neue Sentinel entsteht) und `spec/`, `docs/plan/` (DOK, Greenfield). `docs/user/` ist ein öffentlicher Vertrag, keine Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register ([`../observations/`](../observations/README.md)) durchgegangen am 2026-10-23, Zähler = Dateien in `evidence/`:

- [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (20 Belege, verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.10) und [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (47 Belege, verkörpert): gelten für jede Bedingung in Spezifikation und Handbuch dieses Slice (Tabelle „Aussage → Test“, §3). Der Fund selbst ist ein Fall der Klasse: Stufe 5 der Schreibpipeline sagt „wieder einlesen lassen“, der Code hält sie nicht (ein Beleg für die Welle-Closure, nicht für diesen Slice).
- [`BEO-APP/block-schritt-ohne-golden-pin`](../observations/BEO-APP/block-schritt-ohne-golden-pin/observation.md): berührt den Dockerfile-Block, nicht die Config-Pipeline; kein Treffer für diesen Slice.
- Dual-Classifier-Regel ([`architecture.md`](../../../../spec/architecture.md) §7.2): angewandt in §3 (Ansatz, Sentinel).
- Kein Treffer sonst.

Alle berührten Sub-Areas GF.
