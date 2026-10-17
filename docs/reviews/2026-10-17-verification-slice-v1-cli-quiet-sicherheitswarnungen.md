# Verification Evidence: slice-v1-cli-quiet-sicherheitswarnungen — 2026-10-17 <!-- d-check:ignore (Slice-Kennung im Titel) -->

**Rolle:** Verifier (frischer Kontext, kein Self-Review). Frage: Wurde der Slice richtig gebaut, gemessen an Plan, DoD und Sensors? Form nach `harness/verification.md`. <!-- d-check:ignore (Pfad als Zitat) -->

**Gegenstand:** Slice [`slice-v1-cli-quiet-sicherheitswarnungen`](../plan/planning/done/slice-v1-cli-quiet-sicherheitswarnungen.md) (Welle [`welle-devcontainer-sshd-claude-code`](../plan/planning/welle-devcontainer-sshd-claude-code.md), zweiter Slice), Commits `763cb6f` (Code, Tests), `d316d47` (Handbuch, CHANGELOG, Slice-§3) und `b4e90b2` (Review-Fixes F-1/F-2, Review-Report). Maßstab: §1 (Ziel, Einstufungstabelle, Nicht-Ziele), §2 (DoD), §3 (Plan), §6 (Risiken), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Lastenheft 0.11.0, Ausnahme 1), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer). Review-Report: [`2026-10-17-slice-v1-cli-quiet-sicherheitswarnungen`](2026-10-17-slice-v1-cli-quiet-sicherheitswarnungen.md) (0 HIGH, 0 MEDIUM, 2 LOW, 2 INFO; lief gegen `d316d47`, die Fixes in `b4e90b2`).

## Verification Evidence

Scope:
- Slice: [`slice-v1-cli-quiet-sicherheitswarnungen`](../plan/planning/done/slice-v1-cli-quiet-sicherheitswarnungen.md)
- IDs: [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)
- Artefakte: `internal/hexagon/domain/diagnostic.go`, `internal/hexagon/application/doctor_sandbox.go`, `internal/adapter/driving/cli/doctor.go`, `root.go`, die Tests, `docs/user/benutzerhandbuch.md`, `CHANGELOG.md` <!-- d-check:ignore (Pfade als Zitat) -->

DoD-Abgleich (§2 des Slice; die Häkchen in der Datei sind noch leer, das ist Closure-Arbeit):
- [x] **Kennzeichen und Filter** — `SecurityRelaxation` am `domain.Diagnostic` (`diagnostic.go`, ohne JSON-Tag, Wire-Form `mapDiagnosticsToJSON` nutzt nur `Level`/`Code`/`Message`); vier Warn-Zweige setzen es (siehe (a)); `writeDoctorReport` zeigt unter `--quiet` `SeverityError || DegradationFallback || SecurityRelaxation`, kein Kopf, keine Summary; Hilfetext „errors, degradation fallbacks and security warnings only“; Kommentare nennen die geltende Zusage (`AGENTS.md` §3.9 (a) bis (c) am Diff gelesen: kein „früher/no longer“, `b4e90b2` zieht `root.go` und `strictWarnFailure` nach). <!-- d-check:ignore (Pfade und Bezeichner als Zitat) -->
- [x] **Tests und Mutationen** — Happy, Boundary, Negative und Pins sind vorhanden (siehe (c)); die Mutationen fallen rot (eigene Läufe unten).
- [~] **Doku** — Handbuch (beide `--quiet`-Zeilen, `doctor`-Abschnitt) und CHANGELOG stimmen mit dem Code (siehe (e), (f)); **Abweichung:** `docs/maintainer/cli-json-contract.md` §2.2 sagt weiter „`--quiet` zeigt im Plaintext-Output nur Fehler-Einträge und die Degradations-Fallback-Einträge“ und steht weder im Diff noch in §3 (V-1). <!-- d-check:ignore (Pfad als Zitat) -->
- [x] `make gates` grün — selbst gefahren, Exit 0.
- [x] Review durchgeführt, Report liegt vor; Bewertung der Findings unten.
- [x] Doku-Update bei öffentlichem Vertrag — siehe `[~] Doku`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag — offen (§7 „Wird bei der Closure gefüllt“); Finding-Klassen aus dem Review („Kommentar nennt geänderten Verhaltensumfang nicht mit“, „Lastenheft-Klasse ohne Warn-Eintrag im Code“, „Spec-Nachzug offen“) und V-1 gehören hinein.
- [ ] Beobachtungs-Register fortgeschrieben oder „keine Beobachtung“ in §7 — offen, Closure. Die Klasse „Kommentar beschreibt Abwesendes“ trat im Review wieder auf (F-1/F-2, Zähler für [BEO-CLI/kommentar-beschreibt-abwesendes](../plan/planning/observations/BEO-CLI/kommentar-beschreibt-abwesendes/observation.md) prüfen; die Regel ist schon verkörpert).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang — Ausgang steht an allen drei, Form bei Risiko 2 unscharf (V-2).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) — Closure.

### (a) Einstufung gegen den Code (`doctor_sandbox.go`, alle `SeverityWarn`-Zweige gelesen)

| Zweig (Zeile) | Kennzeichen im Code | Lastenheft-Wortlaut | Befund |
|---|---|---|---|
| `credentials`: Klartext-Schlüssel (129) | `SecurityRelaxation` | Quelle der Zugangsdaten | konform |
| `credentials`: https ohne Token-Quelle (140) | `SecurityRelaxation` | Quelle der Zugangsdaten | konform |
| `egress`: ohne Sandbox-Profil „has no effect“ (194) | `SecurityRelaxation` | „dazu gehört die Meldung, dass die Egress-Restriktion ohne Sandbox-Profil nicht greift“ | konform |
| `egress`: mit `nestedRuntime: docker` (200) | `SecurityRelaxation` | Egress-Restriktion | konform |
| `runtime`: `nestedRuntime` ohne Sandbox-Profil (80) | keines | „Konfigurationsmeldung … gehört nicht dazu“ | konform |
| `runtime`: `/dev/fuse` fehlt, Politik `warn` (105 bis 115) | `DegradationFallback` (nur bei `warn`) | Ausnahme 2 | konform, unverändert |
| `runtime`: `docker` (84) | Stufe OK, kein Eintrag | leere Menge (F-3) | konform |
| Fehler-Zweige (`fail`, AppArmor-Sysctl) | `SeverityError` | zeigt `--quiet` ohnehin | konform |

Keine weiteren `SeverityWarn` in `doctor_sandbox.go` (Suche nach `SeverityWarn`: sechs Fundstellen, alle oben). Die OK-Zweige tragen das Kennzeichen nicht (`TestDoctor_SandboxOKItemsCarryNoQuietFlag`); die Feature-Prüfungen `allowlist` und `drift` tragen es nicht (`TestDoctor_FeatureWarnItemsCarryNoQuietFlag`).

### (b) Wire-Form und Exit-Codes

- `doctor --json`, `--json --quiet`, `--json --strict`: `TestDoctorJSON_SecurityRelaxationFlagNotInWireForm` vergleicht den Umschlag eines gekennzeichneten mit dem eines ungekennzeichneten Eintrags byteweise über drei Argument-Sätze und prüft, dass der Name nicht erscheint; im Code ist `mapDiagnosticsToJSON` unverändert (`doctor.go` 187 bis 201: nur `Level`, `Code`, `Message`). Befund: konform.
- `doctor --quiet --strict` mit nur diesen Warn-Einträgen: `TestExecute_Doctor_QuietStrict_OnlySecurityRelaxation_Exit11` pinnt `ErrDoctorFailures`, `ExitCode == 11`, die Zeilen auf stdout und die Check-Kennungen im Fehlertext. `doctorExitCode` und der Sentinel sind im Diff unverändert (kein neuer Sentinel). Befund: konform ([`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) unverändert).

### (c) Tests

Neu: `TestExecute_Doctor_Quiet_KeepsEachSecurityRelaxationItem` (vier Zweige), `…_UnflaggedWarnsStayHidden` (Negative: `nestedRuntime`, Allowlist, Drift, Host-Warnung, OK), `…_QuietBoundaryAndFullReport` (gleiche Check-Kennung gekennzeichnet vs. nicht; ohne `--quiet` Summary `0 error, 2 warn, 1 ok`), `…_QuietStrict_…_Exit11`, die JSON-Wire-Form, `TestDoctor_FeatureWarnItemsCarryNoQuietFlag`, `TestDoctor_SandboxOKItemsCarryNoQuietFlag`; `TestDoctor_SandboxFallbackFlagIsOnTheItem` trägt je Zweig beide Kennzeichen. Abdeckung der DoD-Zeile „Tests und Mutationen“ vollständig.

### Eigene Mutationen (`make test`, danach `git checkout -- <Datei>`, Baum sauber)

| # | Mutation | Ergebnis |
|---|---|---|
| M-1 | `doctor.go`: Filter ohne `\|\| item.SecurityRelaxation` | rot: `…KeepsEachSecurityRelaxationItem` (4 Unterfälle), `…QuietBoundaryAndFullReport`, `…QuietStrict_…Exit11` |
| M-2 | `doctor_sandbox.go`: Kennzeichen am Egress-Zweig „has no effect“ entfernt | rot: `TestDoctor_SandboxFallbackFlagIsOnTheItem/egress:_no_effect_without_sandbox` |
| M-3 | `doctor_sandbox.go`: Kennzeichen am `nestedRuntime`-Zweig „has no effect“ ergänzt | rot: `TestDoctor_SandboxFallbackFlagIsOnTheItem/runtime:_no_effect_without_sandbox,_same_ID` |

### (d) Review-Findings

| Finding | Bewertung | Begründung |
|---|---|---|
| F-1 LOW (Kommentar `root.go`) | **behoben, nicht blockierend** | `b4e90b2`: „only SeverityError items, degradation-fallback items and security-relaxation items“; Datei steht jetzt in §3. Der Satz ist ohne „früher“-Formulierung. |
| F-2 LOW (Kommentar `strictWarnFailure`) | **behoben, nicht blockierend** | `b4e90b2`: „hides every warn item that carries no degradation-fallback or security-relaxation flag from stdout … reason for those items“; stimmt mit dem Code (die Funktion nennt weiter alle Warn-Einträge, auch die gekennzeichneten; die Aussage „der einzige Ort“ ist auf die ungekennzeichneten eingegrenzt). |
| F-3 INFO (Docker-Daemon-Lockerung als `doctor`-Warn-Eintrag ist leer) | **nicht blockierend** | Menge leer, im Slice §1 und §6 als „entfallen“ mit Beleg ausgewiesen, im Verifier-Bericht des CR (d) bestätigt; `--privileged` weist die Befehlsausgabe aus; die einzige Daemon-Warnung ist der Egress-Eintrag mit `docker` (gekennzeichnet). Der Spec-Slice trägt es (siehe (g)). |
| F-4 INFO (`SPEC-043`/`SPEC-017`) | **nicht blockierend** | `SPEC-043` ist im Spec-Slice vollständig getragen. `SPEC-017` braucht **keinen** Nachzug: die Registry ändert sich nicht (kein neuer Code, keine Bedeutungsänderung); der Spec-Slice nennt `SPEC-017` nur bedingt (Nicht-Ziel: neuer Code gehört dann dorthin). Die Erwähnung in F-4 ist keine Auflage. | <!-- d-check:ignore (Kennungen als Zitat) -->

### (e) CHANGELOG `[Unreleased]`

Die Absätze nennen: `doctor --quiet` listet Fehler, Degradations-Fallbacks und die Sicherheits-Warnungen (vier Fälle, passend zu (a)); `nestedRuntime` „no effect“ bleibt verborgen; Exit-Codes, `--strict`, `--json` unverändert (Test belegt). Der Schlusssatz „The rest of the contract follows in a later release“ nennt nur, was der Slice nicht liefert: Übergabe der Agenten-Zugangsdaten zur Laufzeit samt neuem `doctor`-Check, Ausweisung des Login-Dienstes, `sshd`- und Claude-Code-Features. Stimmt mit §1 Nicht-Ziele und der Welle überein.

Release-Cut-Szenario (Cut nach diesem Slice, ohne die Feature-Slices): Der Eintrag beschreibt Lastenheft 0.11.0 im Indikativ („grows“, „`doctor` checks the source“) und nimmt den noch nicht gelieferten Teil mit dem Schlusssatz aus; der Leser erfährt, was nicht da ist. Restpunkt (V-4, INFO): die Aufzählung „relaxations of the Docker daemon“ im ersten Absatz nennt eine Klasse ohne `doctor`-Eintrag; der Folgesatz mit den vier konkreten Fällen stellt es richtig. Der Eintrag ist nicht falsch, aber der Cut sollte die Aufzählung glätten, wenn die Features nicht mitkommen. Der letzte landende Code-Slice der Welle entfernt den Schlusssatz (A-6 aus der CR-Verifikation bleibt Closure-Pflicht dort).

### (f) Handbuch

- Zeile um 299 und Zeile um 1289 (Optionstabellen) sind textgleich und nennen: `doctor --quiet` zeigt Fehler plus die Einträge der Ausnahmen 1 und 2; in Ausnahme 1 die drei Fälle (Klartext-Schlüssel, fehlende Token-Quelle bei `https`, Egress-Restriktion, auch „greift ohne Sandbox-Profil nicht“) — gegen (a) gelesen: konform; „andere Warnungen, etwa zur Konfiguration, entfallen“ passt zu `nestedRuntime`.
- `doctor`-Abschnitt (um 606): die vier Fälle, „Kopf- und Summary-Zeile entfallen“, `--strict` mit `11`, `--json` unverändert: konform zu Code und Tests.
- `README.md`, `README.de.md`: nennen `doctor --quiet` nicht (Suche), nichts zu tun.

### (g) Folge-Slices und Welle

- [`slice-v2-spec-devcontainer-sshd-claude-code`](../plan/planning/open/slice-v2-spec-devcontainer-sshd-claude-code.md) trägt: `SPEC-043` Ausnahme 1 (Login-Dienst, `doctor`-Sicherheits-Einträge), die `doctor`-Zeile unter `--quiet`, `nestedRuntime: docker` = OK, Allowlist-/Drift-Warnungen entfallen, Egress „has no effect“ sichtbar, `nestedRuntime` „has no effect“ entfällt, mit Fundstellen im Code (`writeDoctorReport`, Kennzeichen am `Diagnostic`, `sandboxRelaxationWarnings`). `SPEC-017` bedingt (nur bei neuem Code). Die Nachzüge aus F-3/F-4 sind getragen. <!-- d-check:ignore (Identifier als Zitat) -->
- Welle §3 (Closure-Trigger „Sicherheitswarnungen unter `--quiet`“) und §4/§5 stimmen mit dem Verzeichnisstand überein: Reihenfolge CR (`done/`) → dieser Slice (`in-progress/`) → `sshd` → `claude-code` → Spec (alle `open/`); Begründung der Reihenfolge (zwingend vor `claude-code`, aus Wahl vor `sshd`) deckt sich mit Slice §4. Der Spec-Slice nennt in §4 Start nur die Bedingung „vier Slices in `done/`“, konsistent.
- `open/README.md`: Reihenfolge und Verweise stimmen. Zwei Zeilen sind beschreibungsmäßig überholt, ohne dass ein Sensor es meldet: die CR-Zeile endet mit „entfällt, wenn der Projektinhaber keine neue Zusage verlangt“ (der CR ist mit Lastenheft 0.11.0 geliefert), und die Tabelle führt zwei Slices, die nicht mehr in `open/` liegen (gängige Praxis bis zur Closure, V-3 INFO).

## Sensors

| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | selbst gefahren, Exit 0; `0 issues.` (golangci-lint), `coverage-gate: OK — coverage 91.10% meets threshold 90%`, `d-check: 293 Datei(en) geprüft, 0 Befund(e)`, `baseline-verify: v6.18.1 OK` |
| `make doc-trace` | pass | selbst gefahren: `83 Anforderung(en), 0 Waise(n)` |
| `make test` (3 Mutationen) | rot wie erwartet | M-1 bis M-3, jeweils per `git checkout -- <Datei>` zurückgesetzt |
| `make docs-check` | pass | Teil von `make gates`; nach Anlage dieser Evidence erneut gefahren (siehe Übergabe) |
| `make test-docker` | not run | kein Docker-Laufzeitpfad berührt (Diagnose-Kennzeichen und Textbericht; keine Compose-/Engine-Aufrufe) |
| `make verify-depguard` | not run | keine Importregel berührt (ein bool-Feld am Domain-Typ, keine neuen Imports) |

Traceability:
| ID | Beleg |
| --- | --- |
| [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | Ausnahme 1 (`doctor`-Einträge): `writeDoctorReport`, Tests Happy/Boundary/Negative, M-1 bis M-3 |
| [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) | `…QuietStrict_…Exit11` pinnt Exit 11 |
| [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) | `TestDoctorJSON_SecurityRelaxationFlagNotInWireForm` |
| [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) | Einstufung (a), Tests je Zweig |

Plan-vs-Code-Diff (§3 gegen `git diff --stat 97daf20..HEAD`, 11 Dateien):
| Datei im Diff | Plan-Zeile §3 | Befund |
|---|---|---|
| `internal/hexagon/domain/diagnostic.go` | ja | konform |
| `internal/hexagon/application/doctor_sandbox.go` | ja | konform (vier Zweige gekennzeichnet; `nestedRuntime` und `vfs` unverändert) |
| `internal/adapter/driving/cli/doctor.go` | ja | konform (Filter, Hilfetext, Doc-Kommentare, `strictWarnFailure`) |
| `internal/adapter/driving/cli/root.go` | ja (nachgetragen in `d316d47`) | konform |
| `doctor_sandbox_test.go`, `cli_test.go`, `doctor_json_test.go` | ja (Korrektur `cli_test.go` statt `quiet_test.go` vermerkt; `quiet_test.go` im Diff unverändert) | konform |
| `docs/user/benutzerhandbuch.md`, `CHANGELOG.md` | ja | konform |
| Slice-Datei selbst, Review-Report | nicht in §3 | Prozess-Artefakte (Plan-Nachzug, Review-Ausgang), wie bei den Vorgängern üblich; kein Mangel |
| nicht im Diff, aber betroffen: `docs/maintainer/cli-json-contract.md` §2.2 | nein | **Abweichung** (V-1) |

Nicht-Ziele §1: keine Stufenänderung, keine neuen Check-Kennungen (`coderegistry.go` und `SPEC-017` unberührt), Wortlaut der Meldungen unverändert, `thirdPartyFeatureWarnings` ohne Kennzeichen, Spezifikation unberührt (`spec/` im Diff unverändert): eingehalten. <!-- d-check:ignore (Kennungen als Zitat) -->

Carveouts:
- Neu: none
- Gelöst: none
- Unverändert: none berührt

Nicht ausgeführt:
- `make test-docker`, `make verify-depguard` — nicht berührt (siehe Sensors).

Commit / Artefakt:
- `763cb6f`, `d316d47`, `b4e90b2` (HEAD zum Zeitpunkt der Verifikation `b4e90b2`)

## Befunde des Verifiers

| ID | Stufe | Befund | Auflage |
|---|---|---|---|
| V-1 | MEDIUM (Closure) | `docs/maintainer/cli-json-contract.md` §2.2 beschreibt den Plaintext-Umfang von `--quiet` als „nur Fehler-Einträge und die Degradations-Fallback-Einträge“; seit dem Slice kommen die Sicherheits-Warn-Einträge hinzu. Maintainer-Doku (Rang 7) mit Vertragsaussage über `doctor --quiet`, weder im Diff noch in Slice §3; Klasse „Aussage über den Code im Plan/Doku stimmt nicht“. | A-1 |
| V-2 | LOW | Slice §6 Risiko 2 („Neue Warn-Einträge späterer Slices vergessen das Kennzeichen“): Ausgang „eingetreten: `slice-v1-devcontainer-feature-claude-code` trägt es in der DoD“ ist keine kanonische Form; „eingetreten“ meint ein eingetretenes Risiko, hier ist es eine Vorkehrung für einen noch nicht eingetretenen Fall. | A-2 | <!-- d-check:ignore (Kennungen als Zitat) -->
| V-3 | INFO | `open/README.md`: CR-Zeile „entfällt, wenn der Projektinhaber keine neue Zusage verlangt“ überholt; `open/`-Tabelle führt zwei Slices außerhalb von `open/`. | A-3 |
| V-4 | INFO | CHANGELOG erster Absatz zählt „relaxations of the Docker daemon“ unter den `doctor`-Warn-Einträgen (leere Menge); der Folgesatz korrigiert es. | A-3 |

## Verdikt

**Erfüllt mit Auflagen.** Die Einstufung stimmt Zweig für Zweig mit dem Lastenheft-Wortlaut und dem Code überein, die Wire-Form und Exit 11 sind unverändert und gepinnt, drei eigene Mutationen fallen rot, `make gates` und `make doc-trace` sind grün, F-1/F-2 sind behoben, F-3/F-4 sind nicht blockierend und der Spec-Slice trägt die Nachzüge. Der Slice darf **nicht** nach `done/`, bevor A-1 erledigt ist.

Auflagen für die Closure:
- **A-1 (V-1):** `docs/maintainer/cli-json-contract.md` §2.2 auf den geltenden Umfang nachziehen (Fehler, Degradations-Fallbacks, Sicherheits-Warn-Einträge; JSON-Aussage bleibt) und die Datei in §3 des Slice eintragen. <!-- d-check:ignore (Pfad als Zitat) -->
- **A-2 (V-2):** Ausgang von Risiko 2 in kanonische Form bringen (etwa „weiter offen: → BEO-DOK/…“ im Register oder „entfallen: …“ mit Begründung); Planner entscheidet.
- **A-3 (V-3, V-4):** optional/Closure: CR-Zeile in `open/README.md` aktualisieren; Aufzählung im CHANGELOG glätten, falls ein Release-Cut ohne die Feature-Slices erfolgt.
- **A-4:** Closure-Notiz §7 mit Steering-Loop-Lerneintrag (Klassen aus Review und V-1), Beobachtungs-Register (Zähler für `BEO-CLI/kommentar-beschreibt-abwesendes`, Klasse trat mit F-1/F-2 erneut auf) und Ausgängen aller Risiken. <!-- d-check:ignore (Identifier als Zitat) -->
