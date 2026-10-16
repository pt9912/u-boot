# Slice Harness: Festlegungen aus den `internal/`-READMEs in die Spezifikation

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-harness-wartung`](../welle-harness-wartung.md) — folgt auf
[`slice-harness-dcheck-bump-reviewer-skill`](../done/slice-harness-dcheck-bump-reviewer-skill.md), damit sein Review mit dem aufgefrischten Skill läuft.

**Bezug:** [`LH-LESE-003`](../../../../spec/lastenheft.md#lh-lese-003--dokumentenordnung)
(Dokumentenordnung: technische Festlegungen stehen im technischen Dokument),
[`SPEC-012`](../../../../spec/spezifikation.md#spec-012--schichten-und-verzeichnislayout-unter-internal),
Sub-Area `internal/**/README.md` ([`conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration); Vorgänger:
[`slice-harness-internal-readme-kennungs-retrofit`](../done/slice-harness-internal-readme-kennungs-retrofit.md).
Auslöser: Frage des Projektinhabers (2026-10-16), ob die Inhalte der Paket-READMEs nicht in
die Spezifikation gehören. Keine Anforderungsänderung.

**Berührte Spec-Stellen:** `spec/spezifikation.md` (neue `SPEC-<NNN>` für Festlegungen,
die bisher nur im Code und in den READMEs stehen); `spec/architecture.md` (Schichtbeschreibung,
nur falls dort eine Lücke bleibt).

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die sieben `internal/**/README.md` (zusammen 490 Zeilen, rund 110 `LH-`-Verweise)
tragen nur noch, was in ein Paket-README gehört: Zweck der Schicht und Verweis auf die
Spezifikation. Technische Festlegungen (Zahlenwerte, Schwellen, Wahrheitstabellen,
Zustandsautomaten, Whitelists) stehen als `SPEC-<NNN>` in der Spezifikation; die
Zuordnung Service → Anforderung entfällt, sie folgt aus den Kennungen (`make doc-trace`).

**Befund (2026-10-16), Stichprobe `internal/hexagon/application/README.md`:**
`pollInterval=500ms`, `dialTimeout=300ms`, Restart-Loop-Zähler mit Schwelle 3, Dead-Allowlist,
§T5-Wahrheitstabelle des `--volumes`-Pfads und die 3-Element-Whitelist von `config` stehen nur
im Code und im README, nicht in der Spezifikation (0 Treffer in `spezifikation.md` und
`architecture.md`).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Änderung des Codes oder der Werte: der Slice dokumentiert den Ist-Stand; weicht ein Wert
  vom README ab, gilt der Code, und die Abweichung wird im Slice vermerkt.
- Neue Produktanforderungen: Festlegungen sind Verfeinerungen vorhandener Anforderungen
  (`LH-FA-UP-*`, `LH-FA-CONF-*` usw.), keine neuen Zusagen
  ([`LH-LESE-003`](../../../../spec/lastenheft.md#lh-lese-003--dokumentenordnung)).
- Go-Paketdoku (`doc.go`) für alle Pakete: nur dort, wo ein README ersetzt wird und der
  Implementierungshinweis (z. B. `noop*`-Defaults) am Code bleiben soll.

## 2. Definition of Done

- [x] **Inventur:** je README eine Liste seiner Aussagen, eingeordnet als Festlegung
  (→ Spezifikation), Rückverfolgbarkeit (→ entfällt), Aufbau der Schicht (→ Architektur-Sicht
  oder Verweis) oder Implementierungshinweis (→ `doc.go` oder Kommentar); Abweichungen vom
  Code sind benannt.
- [x] **Spezifikation:** die Festlegungen stehen als `SPEC-<NNN>` (ab der nächsten freien Nummer) in
  `spezifikation.md`, jede an der Anforderung verankert, die sie verfeinert, mit den
  Zahlenwerten aus dem Code.
- [x] **READMEs** auf Zweck und Verweis gekürzt, die handgepflegte Service-Anforderungs-Liste
  entfällt; `make docs-check` grün.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/README.md` und sechs Paket-READMEs | update | auf Zweck und Verweis kürzen |
| `spec/spezifikation.md` | update | neue `SPEC-<NNN>` für die Festlegungen |
| `spec/architecture.md` | prüfen | Schichtbeschreibung nur bei Lücke |
| `internal/**/doc.go` | neu, wo nötig | Implementierungshinweise am Code |
| `harness/conventions.md` | prüfen | Zeile der Sub-Area `internal/**/README.md` |

Ansatz: ein README je Commit, beginnend mit `internal/hexagon/application/` (die meisten
Festlegungen); Zahlenwerte aus dem Code lesen, nicht aus dem README übernehmen.

## 4. Trigger

**Start** (`next` → `in-progress`): Priorisierung durch den Projektinhaber; kein anderer
Slice in `in-progress/`.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): nach Schicht aufteilen (`application`, `adapter`, `port`,
  `domain`), je ein Slice.
- `in-progress` → `open` (blockiert): ein README-Wert widerspricht dem Code, und die Frage,
  welcher gilt, ist nicht entschieden.

## 5. Closure-Trigger

DoD vollständig, Review- und Verification-Evidence liegen vor, Closure-Notiz
geschrieben, `make gates` grün.

## 6. Risiken und offene Punkte

- Ein README nennt einen Wert, der vom Code abweicht; die Spezifikation übernähme die
  falsche Zahl — **Ausgang:** eingetreten und aufgelöst: der Review fand einen falschen Ist-Stand
  (`SPEC-039`, F-1), er ist gegen den Code korrigiert; die Inventur benennt die übrigen <!-- d-check:ignore (Kennungen als Zitat) -->
  Abweichungen README gegen Code.
- Die Spezifikation wächst um viele Einträge und wird unübersichtlich —
  **Ausgang:** entfallen: sieben Einträge (`SPEC-036`, `LH-FA-UP-001.a`, `SPEC-038` bis `SPEC-042`), <!-- d-check:ignore (Kennungen als Zitat) -->
  je Anforderung gebündelt.
- Die Sub-Area `internal/**/README.md` ist in `conventions.md` als Greenfield deklariert;
  ihre Begründung nennt den Zweck „beschreiben den Code-Bestand“ — **Ausgang:** entfallen: die
  Zeile ist nachgezogen, der Modus bleibt Greenfield.
- Die Inventur fand eine Abweichung Lastenheft gegen Code bei der Verbosity
  ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)), die die Spezifikation
  nicht übernehmen darf — **Ausgang:** eingetreten:
  [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../open/slice-v1-cli-verbosity-letzte-option-gewinnt.md)
  (Code) und [`slice-cr-cli-quiet-fehler-und-warnungen`](../open/slice-cr-cli-quiet-fehler-und-warnungen.md) (Lastenheft).

## 7. Closure-Notiz

**Belege:** Review-Report [`2026-10-16-slice-harness-internal-readmes-festlegungen-in-spec`](../../../reviews/2026-10-16-slice-harness-internal-readmes-festlegungen-in-spec.md)
(0 HIGH, 1 MEDIUM, 2 LOW, 5 INFO; alle behoben); Verification Evidence
[`2026-10-16-verification-slice-harness-internal-readmes-festlegungen-in-spec`](../../../reviews/2026-10-16-verification-slice-harness-internal-readmes-festlegungen-in-spec.md)
(bestanden); `make gates` grün, `make doc-trace`: 83 Anforderungen, 0 Waisen.

- **Was hat funktioniert:** Die Inventur vor der Spezifikation fand die Abweichungen, bevor sie in die Spec wanderten;
  die Werte kamen aus dem Code. Die sieben READMEs sind von 490 auf 109 Zeilen gekürzt.
- **Was ging anders als geplant:** Das README und die Spec sagten Verschiedenes zum Changelog-Pfad (F-1): Die erste
  Fassung von `SPEC-039` beschrieb eine Reparatur, die der Code nicht macht; der Review fing das gegen <!-- d-check:ignore (Kennungen als Zitat) -->
  `hasChangelogUnreleased`. Ein Eintrag verfeinerte nur eine Anforderung und trug die falsche ID-Form
  (`SPEC-037` wurde `LH-FA-UP-001.a`); der Kopf der Spec deckt jetzt auch den Mehr-Anforderungs-Fall ab. Dabei blieb die <!-- d-check:ignore (Kennungen als Zitat) -->
  Nummer 037 als Lücke stehen, wie `SPEC-011`; die Baseline verbietet das Nachbelegen. Die Abweichung bei der <!-- d-check:ignore (Kennungen als Zitat) -->
  Verbosity führte zu zwei Folge-Slices, siehe §6.
- **Steering-Loop-Eintrag:** Ein Review gegen den Code fängt, was ein README-getriebener Entwurf falsch übernimmt: Die
  Zahlenwerte einer Spec stehen aus dem Code, nie aus einem README. Gezählt, nicht verkörpert; das Beobachtungs-Register ist
  leer und die Beobachtung trat einmal auf.
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung angelegt (erstes Auftreten).
- **Folge-Slices:** [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../open/slice-v1-cli-verbosity-letzte-option-gewinnt.md),
  [`slice-cr-cli-quiet-fehler-und-warnungen`](../open/slice-cr-cli-quiet-fehler-und-warnungen.md).
- **Risiken aus §6:** jedes trägt einen Ausgang (siehe §6).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `internal/**/README.md` (Sub-Area
Code-Paket-READMEs) und `spec/`
([`conventions.md`](../../../../harness/conventions.md) §Modus-Deklaration); beide erfüllen
die Schwelle, keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** das Register ist leer.

**Modus:** alle berührten Sub-Areas GF.

## 9. Inventur der README-Aussagen (DoD „Inventur")

Einordnung je Aussage: **F** Festlegung (→ Spezifikation), **R** Rückverfolgbarkeit (entfällt),
**A** Aufbau der Schicht (steht in der Architektur-Sicht oder entfällt), **I** Implementierungshinweis
(bleibt am Code). Zahlenwerte stammen aus dem Quelltext (Stand HEAD), nicht aus den READMEs.
[`SPEC-036`](../../../../spec/spezifikation.md#spec-036--stabilisierung-von-u-boot-up-takt-zustandsklassen-und-entscheidungsmatrix)..[`SPEC-042`](../../../../spec/spezifikation.md#spec-042--aufruf-der-compose-engine-und-preflight) sind die in diesem Slice neu angelegten Einträge.

### `internal/hexagon/application/README.md`

| Aussage | Einordnung | Ziel |
|---|---|---|
| Service je Use-Case-Familie, keine externe I/O, Import-Regeln | A | Architektur §2.2, [`SPEC-012`](../../../../spec/spezifikation.md#spec-012--schichten-und-verzeichnislayout-unter-internal)/[`SPEC-013`](../../../../spec/spezifikation.md#spec-013--import-regel-tabelle-der-schichten) |
| nil-tolerante Ports über `noop*`-Defaults | I | steht an `noop.go` und in Architektur §2.2 |
| Service → Anforderung-Liste (`LH-FA-*`-Spannen je Service) | R | entfällt (`make doc-trace`) |
| `up`: `pollInterval=500ms`, `dialTimeout=300ms`, Zustandsklassen, Dead-Allowlist, Restart-Zähler (Schwelle 3), Healthcheck-Dominanz, Warn-Diagnosen | F | [`SPEC-036`](../../../../spec/spezifikation.md#spec-036--stabilisierung-von-u-boot-up-takt-zustandsklassen-und-entscheidungsmatrix) |
| `parseComposePort`: 8 Syntaxfälle, nicht probebare Formen | F | [`LH-FA-UP-001.a`](../../../../spec/spezifikation.md#lh-fa-up-001a--erkennung-probebarer-ports-in-composeyaml) |
| `down --volumes`: §T5-Wahrheitstabelle | F | [`SPEC-038`](../../../../spec/spezifikation.md#spec-038--bestätigungspfad-von-u-boot-down---volumes) |
| `generate`: 4-Zustands-Maschine, Changelog-Reparaturpfad, atomarer Devcontainer-Plan, `forwardPorts`-Quelle | F | [`SPEC-039`](../../../../spec/spezifikation.md#spec-039--zustandsmaschinen-der-generatoren-u-boot-generate) |
| `config`: Pfad-Whitelist, zweistufige Validierung vor dem Schreiben, `services.<svc>.enabled` nur lesbar | F | [`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config) |
| `compose.yaml.valid` stuft „keine Services" als Warn ein | F | [`SPEC-041`](../../../../spec/spezifikation.md#spec-041--schweregrade-der-prüfung-composeyamlvalid) |
| `init --devcontainer`: zwei Templates, gleiche `planFile`-Pipeline | A | [`LH-AK-005`](../../../../spec/lastenheft.md#lh-ak-005--devcontainer-flow), Architektur §2.2 (Pipeline-Detail bleibt am Code) |
| `init --devcontainer` setzt `devcontainer.enabled: true` | F | [`SPEC-005`](../../../../spec/spezifikation.md#spec-005--schema-der-projektkonfiguration-u-bootyaml) ([`LH-AK-005`](../../../../spec/lastenheft.md#lh-ak-005--devcontainer-flow) nennt den Wert nicht; Wert aus `executeUBootYAML`) |
| `TemplateInitService` wird vom `InitProjectService` konsumiert | A | steht in Architektur §2.2 |
| `managedblock/` (3 Marker-Stile), `templates/` | A / F | Architektur §2.2; Stile in [`SPEC-014`](../../../../spec/spezifikation.md#spec-014--markierungsformat-verwalteter-bereiche-je-dateityp) |
| `remove --purge` nur über das Bestätigungs-Gate | R | [`LH-FA-ADD-007`](../../../../spec/lastenheft.md#lh-fa-add-007--service-entfernen), [`SPEC-038`](../../../../spec/spezifikation.md#spec-038--bestätigungspfad-von-u-boot-down---volumes) für den Bestätigungsmechanismus |
| Anti-Drift-Pin `collectActiveServicePorts` | I | Code-Kommentar in `doctor.go`; Quelle in [`SPEC-039`](../../../../spec/spezifikation.md#spec-039--zustandsmaschinen-der-generatoren-u-boot-generate) |

Abweichungen README ↔ Code:

- „11 Checks" — der Code führt 16 Prüf-Kennungen (`DoctorCheckIDs`, abgeglichen mit [`SPEC-017`](../../../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry)).
- „3-Element-Whitelist von `config`" — der Code kennt 14 Pfad-Arten (`ConfigPathKind`), 13 schreibbar, eine nur lesbar.
- `AddServiceService` nennt die Anforderungen für das Hinzufügen und den Mehrfach-Schutz, nicht aber [`LH-FA-ADD-006`](../../../../spec/lastenheft.md#lh-fa-add-006--add-on-abhängigkeiten); der Code deckt die Abhängigkeitsauflösung ab (`addservice_dependencies.go`).
- „4 Zeilen × 2 Sub-Cases bei AssumeYes-und-NonInteractive" — der Code hat vier Zeilen; bei `AssumeYes` ist `NonInteractive` ohne Wirkung, es gibt keine zwei Sub-Cases.
- „atomarer Two-File-Plan" — der Devcontainer-Plan umfasst bis zu fünf Dateien (zusätzlich `sandbox-init.sh`, `docker-init.sh`, `egress-init.sh` je nach Sandbox-Profil).
- `Doctor::collectActiveServicePorts` — die Funktion heißt `collectActiveServicePorts` (Methoden-Schreibweise fremder Sprache).

### `internal/hexagon/domain/README.md`

| Aussage | Einordnung | Ziel |
|---|---|---|
| Reine Datentypen, nur Standardbibliothek | A | Architektur §2.1, [`SPEC-013`](../../../../spec/spezifikation.md#spec-013--import-regel-tabelle-der-schichten) |
| `ProjectName`/`NormalizeProjectName`/`Project`/`SchemaVersionCurrent = 1` | A | Architektur §2.1; Regel in [`SPEC-004`](../../../../spec/spezifikation.md#spec-004--projektnamen-muster) und [`LH-FA-INIT-002.a`](../../../../spec/spezifikation.md#lh-fa-init-002a--normalisierung-des-abgeleiteten-projektnamens) |
| `ServiceState` (6 Zustände) | A | Architektur §2.1; Semantik in [`LH-FA-ADD-005.a`](../../../../spec/spezifikation.md#lh-fa-add-005a--zustandsregeln-für-registrierte-und-aktive-services) |
| `Severity` 4-stufig inkl. `SeverityInfo` | A | Architektur §2.1 — **Lücke**, dort stand dreistufig; in diesem Slice ergänzt |
| `ContainerState`, Dead-Allowlist, `RestartLoopThreshold = 3`, `StabilizationOutcome` | F | [`SPEC-036`](../../../../spec/spezifikation.md#spec-036--stabilisierung-von-u-boot-up-takt-zustandsklassen-und-entscheidungsmatrix) |
| `ServiceStatus` 4 Spalten, `UpResult` | R | [`LH-FA-UP-003`](../../../../spec/lastenheft.md#lh-fa-up-003--startstatus-anzeigen) nennt die Mindestangaben |
| `Artifact` 4 Werte | R | [`LH-FA-GEN-001`](../../../../spec/lastenheft.md#lh-fa-gen-001--generate-befehl) |
| `ConfigPath` mit `WriteAllowed` | F | [`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config) |
| „Geplante Erweiterungen" (`Port`, `ImageRef`, `ComposeFile`, `EnvVar`) | — | entfällt (Plan, kein Ist-Stand) |

Abweichungen: „3-Pfad-Whitelist" (siehe oben, 14 Pfad-Arten); das README nennt `featurename`, `addon_dependency`, `sandbox`, `template_*` nicht, obwohl sie im Paket liegen.

### `internal/hexagon/port/driving/README.md`

| Aussage | Einordnung | Ziel |
|---|---|---|
| Interfaces je Use-Case-Familie mit Request/Response, Sentinels | A | Architektur §2.3 |
| Sentinel-Listen je Use-Case | R / I | Code (`*.go` neben dem Interface), Klassifikation in Architektur §7 |
| „Geplante Erweiterungen" (`Generate`, `Config`, `Logs`) | — | entfällt |

Abweichungen: `GenerateUseCase`, `ConfigUseCase`, `LogsUseCase` sowie die Remove-/Template-Ports liegen im Paket; „Doctor: Keine Sentinels" stimmt, der Exit-Code-Sentinel `ErrDoctorFailures` liegt im CLI-Adapter.

### `internal/hexagon/port/driven/README.md`

| Aussage | Einordnung | Ziel |
|---|---|---|
| Interfaces, die die Application nach außen nutzt | A | Architektur §2.4 |
| Methodenlisten je Port | R / I | Code (Interface-Doku) |
| `NetProbe` Kontext-Vorrang, `ComposeUp` ohne Status-Folgeabfrage, Preflight-Vertrag, Sentinels 11/12 | F | [`SPEC-042`](../../../../spec/spezifikation.md#spec-042--aufruf-der-compose-engine-und-preflight), Architektur §7.1 |
| `Clock.Sleep` ohne Kontext (Produktion delegiert an `time.Sleep`; nur der Polling-Loop von `up` ruft es) | I | Implementierungshinweis am Code (`internal/hexagon/port/driven/clock.go`), keine Vertragsaussage; weder [`SPEC-042`](../../../../spec/spezifikation.md#spec-042--aufruf-der-compose-engine-und-preflight) noch Architektur §7.1 tragen sie |
| `Confirmer` mit eng benannten Methoden | A | Architektur §2.4 |
| Verweis [`SPEC-027`](../../../../spec/spezifikation.md#spec-027--zusammensetzung-des-lint-profils) als Logging-Port | R | **falscher Anker** ([`SPEC-027`](../../../../spec/spezifikation.md#spec-027--zusammensetzung-des-lint-profils) ist das Lint-Profil); entfällt |
| „Geplante Erweiterungen" (`Logs`, `Exec`) | — | entfällt |

Abweichungen: `DockerEngine` hat acht Methoden (`ComposeUp`, `ComposeDown`, `ComposePs`, `ComposeProject`, `ListVolumeNames`, `RemoveVolume`, `ComposeUpPlan`, `ComposeLogs`), das README nennt drei und `Logs` als geplant; `Confirmer` hat drei Methoden (zusätzlich `ConfirmAddDependency`); `YAMLCodec` hat fünf (`Marshal`, `Unmarshal`, `PatchScalar`, `PatchMappingEntryYAML`, `LocateMarkedEntry`), das README nennt zwei; `Git`, `FileSystem`, `Clock`, `NetProbe`, `Logger`, `ProgressPort` stimmen. Ports `TemplateCatalog`, `TemplateFiles`, `RecorderPort`, `RuntimeEnvironment` fehlen im README.

### `internal/adapter/driven/README.md`

| Aussage | Einordnung | Ziel |
|---|---|---|
| Adapter je Port, `var _ driven.X`-Pin | A | Architektur §2.6 |
| Adapter-Technik (stdlib, `yaml.v3`, `os/exec`, `slog`) | A | Architektur §2.6 |
| Git: nur Exit-Code 128 heißt „kein Repository" | A | Architektur §2.6 |
| Docker-Engine: Aufrufformen, Preflight, Klassifikation 11/12 | F | [`SPEC-042`](../../../../spec/spezifikation.md#spec-042--aufruf-der-compose-engine-und-preflight) |
| `netprobe/` statt `net/` (Alias-Konflikt) | I | Verzeichnisname, bleibt am Code |
| Build-Tag-Pfad `docker` | R | `make test-docker`, Architektur §8 |

Abweichungen: Das README führt neun Adapter; im Paket liegen vierzehn Verzeichnisse (fehlen: `externaltemplates`, `localtemplates`, `recordingfs`, `runtime`, `templateyaml`; in Architektur §2.6 stehen sie).

### `internal/adapter/driving/README.md`

| Aussage | Einordnung | Ziel |
|---|---|---|
| Cobra-Command je Subkommando, Wiring in `cmd/uboot/main.go` | A | Architektur §2.5 / §2.7 |
| Command → Anforderung-Liste | R | entfällt |
| Exit-Code-Klassifikation je Sentinel (0/2/10/11/12/14/1) | A | Architektur §7.1 und §7.2; Code (`ExitCode`) |
| Verbosity-Flags steuern den `slog`-Pegel (`--debug`/`--verbose` Debug, `--quiet` Warn, sonst Info; Vorrang debug > verbose > quiet) | F — **Abweichung vom Lastenheft** | nicht übernommen, siehe unten |
| `statusview.go`: asymmetrischer `--quiet`-Vertrag | I | Code-Kommentar |
| „Geplante Erweiterungen" (`template`, `logs`, `--json`) | — | entfällt |

Abweichungen: `template`, `logs` und `--json` sind umgesetzt; die Verbosity-Auflösung des Codes (feste Rangfolge `--debug` vor `--verbose` vor `--quiet`) widerspricht
[`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) („die zuletzt angegebene Option gewinnt"), und `--quiet` hebt nur den Pegel auf Warn statt „nur Fehler". Die Spezifikation
darf dem Lastenheft nicht widersprechen; beide Punkte stehen deshalb **nicht** in der Spezifikation und sind ein offener Befund für den Projektinhaber (Code-Änderung ist nicht Gegenstand dieses Slice).

### `internal/README.md`

| Aussage | Einordnung | Ziel |
|---|---|---|
| Layout-Baum, Import-Regeln, Verweise auf [`SPEC-012`](../../../../spec/spezifikation.md#spec-012--schichten-und-verzeichnislayout-unter-internal)/[`SPEC-013`](../../../../spec/spezifikation.md#spec-013--import-regel-tabelle-der-schichten)/[`SPEC-034`](../../../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln) | A | [`SPEC-012`](../../../../spec/spezifikation.md#spec-012--schichten-und-verzeichnislayout-unter-internal), [`SPEC-013`](../../../../spec/spezifikation.md#spec-013--import-regel-tabelle-der-schichten), [`SPEC-034`](../../../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln) (bleibt als Verweis) |
| Kurz-Inventar je Schicht | A / R | Architektur §2 |
| CLI-Subcommands-Tabelle | R | entfällt (`make doc-trace`) |
| Coverage-Scope `./internal/...` | F | [`SPEC-032`](../../../../spec/spezifikation.md#spec-032--bootstrap-verhalten-der-coverage-stage), [`SPEC-033`](../../../../spec/spezifikation.md#spec-033--go-layout-des-repositories) |
| `e2e/`- und `acceptance_test.go`-Beschreibung | R / I | Architektur §8 |

Abweichungen: „Alle Subkommandos tragen `--json`, `--dry-run` und `--diff`" ist nicht belegt (zum Beispiel `logs --follow` mit `--json` wird abgewiesen, `ErrFollowJSONNotSupported`).
