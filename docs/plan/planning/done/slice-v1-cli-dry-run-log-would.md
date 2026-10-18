# Slice v1: Dry-Run — keine Info-Logzeilen in Vergangenheitsform

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — kein Release-Bezug festgelegt; der Slice läuft **nach** dem Cut [`slice-v1-release-cut-v0.11.0`](../done/slice-v1-release-cut-v0.11.0.md) und ist nicht Teil von `v0.11.0` (Entscheidung des Projektinhabers 2026-10-18). Herkunft: Beobachtung [`BEO-CLI/dry-run-log-vergangenheitsform`](../observations/BEO-CLI/dry-run-log-vergangenheitsform/observation.md) (Zähler 1), Slice [`slice-v1-cli-dry-run-would-formulierung`](../done/slice-v1-cli-dry-run-would-formulierung.md) (Verification, Abschnitt Diskrepanzen und Beobachtungen).

**Bezug:** [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) (Dry-Run: nichts wird geschrieben, die Ausgabe behauptet kein Ergebnis),
[`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Verbosity und Logging).

**Berührte Spec-Stellen:** keine geplant. Das Lastenheft ist Zielbild ([`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) bleibt unverändert, wie beim Vorgänger-Slice); zeigt sich beim Bau, dass die Spezifikation ([`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet)) eine Aussage zu Logzeilen im Dry-Run trägt, die der Code ändert, zieht der Slice sie nach (Plan-Eintrag des Implementers).

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** Im Dry-Run schreibt kein Use-Case eine `level=INFO`-Logzeile, die ein nicht eingetretenes Ergebnis in Vergangenheitsform meldet.

### Entscheidung des Projektinhabers (2026-10-18)

„Code folgt Zielbild“: Der Dry-Run behauptet nichts als geschehen, auch nicht in den Logzeilen auf stderr. Die Beobachtung [`BEO-CLI/dry-run-log-vergangenheitsform`](../observations/BEO-CLI/dry-run-log-vergangenheitsform/observation.md) bleibt bis zur Closure `offen` (Zähler 1) und wird mit diesem Slice beantwortet.

### Gefundener Stand im Code (Ist-Stand, gelesen am 2026-10-18)

Die Use-Cases schreiben Info-Zeilen nach dem Schreibzugriff, ohne den Vorschau-Modus zu prüfen; der Schreibzugriff geht im Dry-Run in einen aufzeichnenden Dateisystem-Ersatz (`selectFS(req.PreviewMode)`), die Logzeile steht danach unbedingt:

| Logzeile | Fundstelle (Datei, Zeile am Ist-Stand) |
|---|---|
| `config set: updated` (Einzelpfad) | `internal/hexagon/application/config.go:283` (nach `fs.WriteFile` aus `selectFS`) |
| `config set: updated` (Listenpfad) | `internal/hexagon/application/config.go:1005` |
| `config set: updated` (Egress-Liste) | `internal/hexagon/application/devcontainer_egress.go:193` |
| `config set (many): updated` (mehrere Pfade) | `internal/hexagon/application/config_many.go:148` |
| `generate: created` | `internal/hexagon/application/generate.go:432` und `:647` |
| `generate: repaired Unreleased header` | `internal/hexagon/application/generate.go:487` |
| `generate: updated block` | `internal/hexagon/application/generate.go:706` |
| `generate devcontainer: <Aktion>` (`created`, `updated block`, `repaired`) | `internal/hexagon/application/generate.go:823` |
| `generate devcontainer: allowlist updated` | `internal/hexagon/application/generate.go:1218` |

Die Vorschau-Unterscheidung steht im Request: `req.PreviewMode == driving.PreviewDryRun` (Beispiel `generate.go:846`, der Zweig der Fallback-Abfrage); `--diff` allein (`PreviewAndApply`) schreibt wirklich, dort stimmt die Vergangenheitsform. Der Vorgänger-Slice hat die Zusammenfassung auf stdout umgestellt (`printRemoveSummary`, `printGenerateSummary`, `printConfigSetSummary` im CLI-Adapter); die Logzeilen blieben (Beobachtung). Ein Test pinnt eine dieser Logzeilen: `internal/hexagon/application/config_t3_test.go` (gefunden mit der Suche nach `config set: updated`, der Implementer liest die Stelle und trägt sie in die Plan-Tabelle ein). Die Suche nach `.Info(` in `internal/hexagon/application/` (ohne Tests) findet außer den Zeilen der Tabelle nur `generate_rollback.go:127` („write failed, rolled back to the previous state“; die Rücknahme gibt es nur bei echten Schreibzugriffen, `s.rollbackFS` bleibt im Dry-Run leer: `generate.go:171`) und `doctor.go:212` („doctor: checks complete“, kein Schreibergebnis). Ob eine dieser Zeilen oder ein Aufruf mit einem anderen Methodennamen im Dry-Run ein nicht eingetretenes Ergebnis meldet, liest der Implementer Zeile für Zeile (nicht nach Stichwort entschieden); `add` und `remove` tauchen in der Suche nicht auf.

### Lösungsraum — der Implementer wählt

Der Projektinhaber legt das Zielbild fest, nicht die Lösung. Zwei Lösungen sind benannt; der Implementer wählt eine und begründet sie im Plan-Eintrag:

- **(a)** Im Dry-Run lautet die Zeile „would …“ (gleiche Ebene, gleiche Felder, Wortlaut im Konjunktiv-Vokabular der Zusammenfassung).
- **(b)** Im Dry-Run steigt die Zeile auf `Debug` ab (sie erscheint nur mit `--debug`, wie die No-op-Zeilen).

Eine dritte Lösung ist erlaubt, wenn sie das Zielbild erfüllt (im Dry-Run keine Vergangenheitsform auf `Info`) und der Plan-Eintrag sie begründet. Die Wahl muss für **jede** Fundstelle der Tabelle und jede, die die Suche des Implementers zusätzlich findet, gelten; ein Rest, der in Vergangenheitsform bleibt, ist ein Befund.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Änderung der Zusammenfassung auf stdout, der `--json`-Ausgabe oder der Exit-Codes:** Der Vorgänger-Slice hat sie festgelegt; der Slice ändert nur Logzeilen.
- **Eine Änderung des Logging-Ports oder der Pegelabbildung** ([`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet)): Die Pegel bleiben; geändert wird, auf welchem Pegel oder mit welchem Wortlaut der Dry-Run meldet.
- **Logzeilen außerhalb des Dry-Runs:** Ohne `--dry-run` (auch bei `--diff` allein) sind sie wahr und bleiben.
- **Ein neues Gate für Logwortlaut:** Der Test pinnt das Verhalten (§2); ein Sensor über alle Logzeilen wäre ein eigener Vorgang.

## 2. Definition of Done

- [x] **Dry-Run ohne Info-Zeile in Vergangenheitsform (Code):** Für jede Fundstelle der Tabelle (und jede, die die Suche des Implementers zusätzlich findet) meldet der Dry-Run kein nicht eingetretenes Ergebnis in Vergangenheitsform auf `Info`; Test je Subkommando (`config set` Einzelpfad, Listenpfad, Egress-Liste, `generate` je Artefakt-Aktion, `generate devcontainer`): Happy (Dry-Run, Zeile im gewählten Vokabular oder nicht auf `Info`), Boundary (ohne `--dry-run` und mit `--diff` allein bleibt die Vergangenheitsform), Negative (kein Schreiben, `--json` unverändert) ([`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run)).
- [x] Doku-Update: Handbuch (Abschnitt „Vorschau statt Überraschung“) und CHANGELOG nennen den geltenden Stand der Logausgabe im Dry-Run, falls sich für den Benutzer sichtbar etwas ändert (öffentlicher Vertrag; Nutzer-Doku ohne Kennungen).
- [x] `make gates` grün; `make doc-trace` ohne neue Waise.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben: `BEO-CLI/dry-run-log-vergangenheitsform` bekommt einen Ausgang (`gestrichen` mit Begründung, weil die Ursache wegfiel, oder bleibt `offen`, wenn der Implementer einen Rest findet).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne Wellen-Betrieb hier geprüft.

## 3. Plan (vor Code)

**Gewählte Lösung: (a), „would …“ auf `Info`, über einen gemeinsamen Helfer.** `internal/hexagon/application/dryrun_log.go`, Symbol `logWrite(l, mode, done, would, args...)`: im Dry-Run (`req.PreviewMode == driving.PreviewDryRun`) steht die `would`-Zeile, sonst die bisherige Zeile; Pegel (`Info`), Felder und Logger-Port bleiben (`internal/hexagon/port/driven/logger.go`, unverändert). Begründung gegen (b): (1) Die Zeile bleibt eine Auskunft auf Standardstufe und trägt dieselben Felder (`path`, `old`, `new`, `changed`); bei (b) fehlte sie im Dry-Run ganz und stünde unter `--debug` weiter in Vergangenheitsform, ein Rest, den das Zielbild nicht trägt („Debug“ beantwortet nur den Pegel, nicht den Wortlaut). (2) Der Wortlaut folgt der Zusammenfassung auf stdout (`Would …`), die der Vorgänger-Slice festgelegt hat. (3) `--quiet` (Pegel `Error`, `internal/adapter/driving/cli/verbosity.go`, `slogLevel`) unterdrückt die Zeile wie bisher. `generate devcontainer` bildet die Aktion in `devcontainerWriteMessage` ab (`created` → `would create`, `updated-block` → `would update block`); `devcontainerAggregateAction` in `generate.go` liefert nur `created`, `updated-block` und `no-op`, und `no-op` endet in `generateDevcontainer` vor der Zeile; einen Zweig für weitere Aktionen gibt es deshalb nicht. Nicht berührt, nach Zeilenlektüre: `maybeWarnOrphanFeatureActivation` in `config.go` (`config set: orphan feature activation …`, Hinweis auf den Zustand der Konfiguration, im Dry-Run wahr); `rollbackOnError` in `generate_rollback.go` (Rücknahme nur bei echten Schreibzugriffen, `s.rollbackFS` ist im Dry-Run leer: `generate.go` in `Generate`); `DoctorService.Check` in `doctor.go` (`doctor: checks complete`, kein Schreibergebnis). `add`, `remove`, `init`, `up`, `down`, `logs`, `template` schreiben ihre Info-Zeilen nicht (nur `Debug`-Aufrufe; Suche nach `.Info(`, `.Warn(`, `.Error(`, `.Log(` in `internal/hexagon/application/` ohne Tests).

Grenze: Ein Service ohne Factory (Legacy-Konstruktor `NewConfigService`/`NewGenerateService`) ignoriert `PreviewMode` und schreibt wirklich (`selectFS` in `config.go` und `generate.go`); die Composition-Root verdrahtet die Factory, das Kriterium bleibt dort der Request-Modus.

Jede Aussage über Code steht mit Fundstelle (gelesen am 2026-10-18; die Tabelle nennt Datei und Symbol, keine Zeilennummer als Beleg). Er gleicht die Tabelle nach der letzten Fix-Runde gegen `git diff --stat` ab ([`AGENTS.md`](../../../../AGENTS.md) §3.10); eine Datei, die er als „unverändert“ führt, belegt er durch eine gelesene Aussage.

| Datei / Komponente | Änderungs-Art | Begründung (Fundstelle) |
|---|---|---|
| `internal/hexagon/application/dryrun_log.go` | neu | `logWrite` (Wortlaut je Vorschau-Modus) und `devcontainerWriteMessage` (Aktion → `would create` / `would update block`) |
| `internal/hexagon/application/config.go` | update | `ConfigService.Set` und `setFeatureSourcesAllow`: `logWrite(logger, req.PreviewMode, "config set: updated", "config set: would update", …)` nach dem Schreiben über `selectFS(req.PreviewMode)` |
| `internal/hexagon/application/devcontainer_egress.go`, `internal/hexagon/application/config_many.go` | update | `setEgressAllow`: dieselbe Zeile im Egress-Listenpfad; `SetMany`: `config set (many): would update` |
| `internal/hexagon/application/generate.go` | update | `generateChangelog` und `generateManagedFile` (`would create`), `generateChangelog` (`would repair Unreleased header`), `generateManagedFile` (`would update block`), `generateDevcontainer` (`devcontainerWriteMessage`), `applyAllowExternalFeatureSources` (`would update allowlist`); `applyAllowExternalFeatureSources` bekommt den Vorschau-Modus als Parameter `mode`, `writeDevcontainerPhase` reicht `req.PreviewMode` durch |
| `internal/hexagon/application/generate_rollback.go`, `doctor.go` | unverändert | `rollbackOnError` und `DoctorService.Check` gelesen: Rücknahme nur bei echten Schreibzugriffen, Prüfabschluss ohne Schreibergebnis (Begründung oben) |
| `internal/hexagon/application/config_t3_test.go` | unverändert | `TestConfigSet_SilenceLoggerFalse_EmitsInfo` läuft ohne Vorschau-Modus (Nullwert `PreviewNone`) und erwartet `config set: updated`; die Erwartung gilt weiter |
| `internal/hexagon/application/dryrun_log_test.go` | neu | `TestDryRunLog_WriteLines` (Happy: Dry-Run `would`; Boundary: `PreviewNone` und `PreviewAndApply` Vergangenheitsform; je Schreibpfad; beide Zeilen tragen dieselben Schlüssel und Werte, `infoArgs`), `TestDryRunLog_GenerateDevcontainerUpdatedBlock`, `TestDryRunLog_DryRunWritesNothingAndNoOpStaysQuiet` (Negative: nichts geschrieben, No-op und Fehler ohne Info-Zeile); `fakeLogger` aus `fakes_test.go` |
| `docs/user/benutzerhandbuch.md`, `CHANGELOG.md` | update | Handbuch: Absatz nach der Zusammenfassungs-Regel in `### Vorschau statt Überraschung`; CHANGELOG: `[Unreleased]` · `### Changed` |
| `docs/plan/planning/in-progress/roadmap.md` | update | Ruhe-Marker „Nichts in Arbeit.“ während der Arbeit durch die Arbeitszeile ersetzt, nach dem Move wieder gesetzt; im Diff bleibt die Linkberichtigung der Welle-Zeile `done/` (`make slice-mv`) |
| `docs/plan/planning/done/slice-v1-release-cut-v0.11.0.md` | update (Closure) | Verweise auf diesen Slice nach `done/` nachgezogen (`make slice-mv`), kein inhaltlicher Eingriff |
| `docs/plan/planning/in-progress/slice-v1-cli-help-handbuch-link.md` | update (Closure, fremder Vorgang) | nur der Verweis auf diesen Slice nachgezogen (`make slice-mv`); der übrige Inhalt gehört `044e0c9` |
| `docs/plan/planning/done/slice-v1-cli-dry-run-log-would.md` | update | dieser Plan (gewählte Lösung, Tabelle) |
| `docs/reviews/2026-10-18-slice-v1-cli-dry-run-log-would.md` | neu | Review-Report (Reviewer-Rolle) |
| `docs/plan/planning/observations/BEO-CLI/dry-run-log-vergangenheitsform/state.md` | update (Closure) | Ausgang nach §2 |
| `docs/plan/planning/observations/BEO-CLI/kommentar-beschreibt-abwesendes/evidence/slice-v1-cli-dry-run-log-would.md` | neu (Closure) | Beleg F-2 und F-4: Rückfall nach Verkörperung (`AGENTS.md` §3.9) |
| `docs/reviews/2026-10-18-verification-slice-v1-cli-dry-run-log-would.md` | neu | Verification Evidence (Verifier-Rolle) |
| `docs/plan/planning/open/README.md` | update (Closure) | Zeile dieses Slice: Ort `done/`, „geliefert“ |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): der Cut [`slice-v1-release-cut-v0.11.0`](../done/slice-v1-release-cut-v0.11.0.md) liegt in `done/` (Entscheidung des Projektinhabers 2026-10-18: nach `v0.11.0`); kein anderer Slice liegt in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Findet die Suche des Implementers weit mehr Fundstellen als die Tabelle (etwa Zeilen anderer Methodennamen oder in anderen Schichten), wird nach Use-Case geschnitten.
- `in-progress` → `open` (blockiert): Die Wahl zwischen (a) und (b) braucht eine Entscheidung über Pegel oder Wortlaut, die das Lastenheft ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)) nicht trägt; dann klärt der Projektinhaber die Auslegung, nicht der Plan.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7); das Register trägt für `BEO-CLI/dry-run-log-vergangenheitsform` einen Ausgang.

## 6. Risiken und offene Punkte

Jedes Risiko bekommt bei der Closure genau einen Ausgang (eingetreten / entfallen / weiter offen).

- **Weitere Fundstellen außerhalb der Tabelle:** Die Tabelle nennt, was die Suche nach `.Info(` fand; Zeilen anderer Methodennamen und Zeilen, die ein Port außerhalb der Anwendungsschicht schreibt, sind nicht gelesen. Ein übersehener Rest lässt die Beobachtung bestehen. — **Ausgang:** entfallen (Zeilenlektüre der Logger-Aufrufe der Anwendungsschicht durch Implementer, Review und Verifier; kein Rest, Test je Subkommando).
- **Wortlaut als Vertrag:** Skripte könnten die Logzeilen auswerten (sie stehen auf stderr, `--json` ist der Vertrag für Maschinen, [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe)); eine Wortlautänderung wäre dann eine Verhaltensänderung. — **Ausgang:** eingetreten als sichtbare Änderung, behandelt (CHANGELOG `[Unreleased]` · Changed nennt die Änderung; `--json` ist der Vertrag).
- **`--diff` allein und Dry-Run verwechselt:** Ein Kriterium, das beide gleich behandelt, ändert die wahre Vergangenheitsform von `--diff` allein. — **Ausgang:** entfallen (Boundary-Test in `TestDryRunLog_WriteLines`, Praxislauf des Verifiers).

## 7. Closure-Notiz

**Entscheidung des Projektinhabers (2026-10-18):** „Code folgt Zielbild“. Der Dry-Run behauptet nichts als geschehen, auch nicht in den Logzeilen auf stderr; [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) bleibt unverändert (kein CR).

**Gewählte Lösung:** (a) „would …“ auf `Info` über den Helfer `logWrite`. Begründung gegen (b) Debug: Die Zeile bleibt eine Auskunft auf Standardstufe mit denselben Feldern; bei (b) fehlte sie im Dry-Run und stünde unter `--debug` weiter in Vergangenheitsform; der Wortlaut folgt der stdout-Zusammenfassung des Vorgängers; `--quiet` unterdrückt sie wie bisher.

**Belege:** Review-Report [`2026-10-18-slice-v1-cli-dry-run-log-would`](../../../reviews/2026-10-18-slice-v1-cli-dry-run-log-would.md) (0 HIGH, 0 MEDIUM, 2 LOW, 3 INFO; sieben Mutationen, drei überlebten bis zur Fix-Runde). Verification Evidence [`2026-10-18-verification-slice-v1-cli-dry-run-log-would`](../../../reviews/2026-10-18-verification-slice-v1-cli-dry-run-log-would.md): DoD 1 bis 4 erfüllt, Auflagen 1 bis 6 in dieser Closure erfüllt. `make gates` grün (Coverage 91,10 %), `make doc-trace`: 83 Anforderungen, 0 Waisen. Commits `1f42dfd`, `a79bf59`, `64e853a`, `878e453`, `b35ad74`, `aacee08`, `bf2d592`, `c94eeb6`.

**Fundstellen je Aussage (`AGENTS.md` §3.10):**

| Aussage | Fundstelle |
|---|---|
| Der Dry-Run schreibt `would …`, sonst die Vergangenheitsform; Kriterium genau `mode == driving.PreviewDryRun` (`--diff` allein bleibt Vergangenheitsform) | `logWrite` in `internal/hexagon/application/dryrun_log.go` |
| `generate devcontainer` bildet die Aktion ab (`created` → `would create`, `updated-block` → `would update block`) | `devcontainerWriteMessage` in `dryrun_log.go`, Aufruf in `generateDevcontainer` (`generate.go`) |
| Alle Schreibzeilen laufen über den Helfer | `ConfigService.Set`, `setFeatureSourcesAllow` (`config.go`), `setEgressAllow` (`devcontainer_egress.go`), `SetMany` (`config_many.go`), `generateChangelog`, `generateManagedFile`, `applyAllowExternalFeatureSources` (`generate.go`) |
| Der Vorschau-Modus erreicht die Allowlist-Zeile | Parameter `mode` von `applyAllowExternalFeatureSources`, übergeben in `writeDevcontainerPhase` (`generate.go`) |
| Nicht berührt: Zustands- und Prüfzeilen | `maybeWarnOrphanFeatureActivation` (`config.go`), `rollbackOnError` (`generate_rollback.go`), `DoctorService.Check` (`doctor.go`) |
| Happy, Boundary, Negative; Schlüssel und Werte beider Zeilen gepinnt | `TestDryRunLog_WriteLines`, `TestDryRunLog_GenerateDevcontainerUpdatedBlock`, `TestDryRunLog_DryRunWritesNothingAndNoOpStaysQuiet` in `internal/hexagon/application/dryrun_log_test.go` |
| Ohne Vorschau-Modus gilt die Vergangenheitsform weiter | `TestConfigSet_SilenceLoggerFalse_EmitsInfo` in `config_t3_test.go` (unverändert) |
| Handbuch und CHANGELOG nennen den Konjunktiv, `--quiet` und `--json` als Vertrag | Abschnitt „Vorschau statt Überraschung“ in `docs/user/benutzerhandbuch.md`; `[Unreleased]` · Changed in `CHANGELOG.md` |

**Praxislauf (Verifier, Runtime-Images aus `fa019d0` und `HEAD`):** `config set project.name renamed --dry-run`: vorher `config set: updated`, nachher `config set: would update` mit `path`, `old`, `new`. `generate devcontainer --dry-run`: vorher `generate devcontainer: created`, nachher `generate devcontainer: would create`, Felder gleich. `generate changelog --dry-run`: `generate: created` → `generate: would create`. `--diff` allein und ohne `--dry-run`: unverändert `updated` / `created`; `--dry-run --quiet`: keine Zeile; `u-boot.yaml` nach den Dry-Runs unverändert.

**Plan-Tabelle gegen `git diff --stat fa019d0..HEAD` (§3.10):** Die zwei verschobenen Zeilenhinweise (Allowlist-Zeile, Signatur von `applyAllowExternalFeatureSources`) und alle weiteren stehen als Datei und Symbol; keine Zeilennummer ist Beleg. Der Diff nennt dreizehn Dateien: elf gehören zum Slice (jede mit Plan-Zeile; dazu die Closure-Dateien dieser Notiz); **fremd** sind `044e0c9` und `docs/plan/planning/open/README.md` samt `in-progress/slice-v1-cli-help-handbuch-link.md` (Planer-Vorgang des anderen Slice), nicht Gegenstand. „Unverändert“ per gelesener Aussage: `generate_rollback.go`, `doctor.go`, `config_t3_test.go`, `logger.go`.

- **Was hat funktioniert:** Die Mutationen des Reviews trafen jeden Zweig des Kriteriums; der Praxislauf bestätigte Vorher und Nachher am Binary; die Fix-Runde schloss F-1 bis F-5 vor der Verification.
- **Was ging anders als geplant:** Der erste Test pinnte nur die Meldung, nicht die Felder (F-1); der Plan trug Zeilenhinweise, die nach der Fix-Runde verschoben standen (F-3, Verifier-Diskrepanz).
- **Finding-Klassen des Reviews:** Test pinnt Zusage „Felder unverändert“ nicht (F-1, LOW; drei überlebende Mutationen; behoben) · Kommentar beschreibt Vorzustand (F-2, LOW, `AGENTS.md` §3.7/§3.9 a; behoben) · Zeilenhinweis im Plan ungenau (F-3, INFO; behoben) · Doc-Kommentar kennt neuen Parameter nicht (F-4, INFO, §3.9 b; behoben) · Toter Rückfallzweig ohne Test (F-5, INFO; Zweig gestrichen).
- **Steering-Loop-Eintrag:** **Gezählt, nicht verkörpert.** Kein `liegt in`, keine Entscheidung über Verkörperung. (1) F-2 und F-4 gehören zu `BEO-CLI/kommentar-beschreibt-abwesendes` (verkörpert in `AGENTS.md` §3.9): ein Rückfall **nach** der Verkörperung, als Beleg `evidence/slice-v1-cli-dry-run-log-would.md` eingetragen; ob die Regel nachgeschärft oder ein Sensor nötig wird, entscheidet diese Closure nicht (offen für den Projektinhaber). (2) F-1, F-3 und F-5 sind einmalig; sie stehen hier benannt, nicht im Register (ein Zähler 1 ohne Wiederauftreten sagt nichts; tritt eine Klasse in einem zweiten Vorgang auf, legt dessen Closure den Pfad an). F-3 (Zeilennummer im Plan verrutscht) berührt die Klasse von `BEO-DOK/plan-aussage-ueber-code` (verkörpert, §3.10); sie ist hier benannt und nicht als Beleg eingetragen. **Lehre:** Ein Test, der eine Zusage „unverändert“ trägt, muss die Zusage vergleichen (Schlüssel und Werte), nicht nur die Meldung; eine Zeilennummer im Plan ist kein Beleg, Datei und Symbol sind es.
- **Beobachtungs-Register (`../observations/`):** [`BEO-CLI/dry-run-log-vergangenheitsform`](../observations/BEO-CLI/dry-run-log-vergangenheitsform/observation.md) steht `gestrichen`: Die Ursache fiel weg; der Dry-Run meldet auf `INFO` nur noch den Konjunktiv, `TestDryRunLog_WriteLines` pinnt ihn. Nicht verkörpert (Zähler 1, unter der Schwelle, kein Zielort). Der eigene Slice zählt nicht als Beleg, er beantwortet die Beobachtung; das Verzeichnis bleibt liegen. [`BEO-CLI/kommentar-beschreibt-abwesendes`](../observations/BEO-CLI/kommentar-beschreibt-abwesendes/observation.md): eine weitere Evidence-Datei (Rückfall nach Verkörperung), Stand unverändert.
- **Folge-Slices:** keine.
- **Risiken aus §6:** alle drei tragen einen Ausgang: entfallen (weitere Fundstellen), eingetreten als sichtbare Änderung und behandelt (Wortlaut als Vertrag: CHANGELOG, `--json` ist der Vertrag), entfallen (`--diff` allein und Dry-Run).
- **Roadmap:** Die Arbeitszeile steht nach dem Move wieder auf „Nichts in Arbeit.“ (`in-progress/` trägt keinen Slice; d-check `planning`).
- **Drei Paarungen:** Anker (kein `liegt in`, nichts verkörpert) · Folge-Slice (keiner) · Register (zwei Pfade, je mit nicht leerem `evidence/`) — nach dem `git mv` nach `done/` geprüft.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `internal/hexagon/application` (APP), `spec/`, `docs/plan/` (DOK) — alle Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006), jede erfüllt die Schwelle ≥ 2 von 3 Achsen; keine Ausdifferenzierung nötig. `docs/user/` und `CHANGELOG.md` sind öffentliche Verträge, keine Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register ([`../observations/`](../observations/README.md)) durchgegangen — `BEO-CLI/dry-run-log-vergangenheitsform` ist der Gegenstand (`offen`, Zähler 1); `BEO-CLI/kommentar-beschreibt-abwesendes` (verkörpert, `AGENTS.md` §3.9) und `BEO-DOK/plan-aussage-ueber-code` (verkörpert, `AGENTS.md` §3.10) gelten für Kommentare und Plan-Aussagen dieses Slice.

Alle berührten Sub-Areas GF.
