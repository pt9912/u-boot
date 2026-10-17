# Review-Report: slice-v1-cli-quiet-sicherheitswarnungen — 2026-10-17 <!-- d-check:ignore (Slice-Kennung im Titel) -->

**Review-Art:** Code — Code-Review gegen Plan + Konventionen.

**Gegenstand:** Commit-Range `97daf20..HEAD` (`763cb6f` Code und Tests, `d316d47` Handbuch, CHANGELOG, Slice-§3) des Slice `slice-v1-cli-quiet-sicherheitswarnungen` (Welle `welle-devcontainer-sshd-claude-code`). <!-- d-check:ignore (Slice- und Wellen-Kennung als Zitat) -->

**Skill:** `.harness/skills/reviewer.md` @ HEAD `d316d47` <!-- d-check:ignore (Pfad als Zitat) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-17

**Eingangs-Kontext:**

- Slice-Plan `slice-v1-cli-quiet-sicherheitswarnungen` (in-progress) <!-- d-check:ignore (Slice-Kennung als Zitat) -->
- [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Lastenheft 0.11.0, Ausnahme 1), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)
- [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) und [`SPEC-017`](../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry) (Nachzug im Spec-Slice der Welle)
- `AGENTS.md` (Hard Rules, §3.7, §3.9), `harness/review.md`

**Sensors:** `make test` (nur zur Mutation, drei Läufe, jeweils zurückgesetzt). `make gates` und `make docs-check` über den Report liefen nicht als Teil der Code-Prüfung; `make docs-check` nach dem Schreiben dieses Reports separat. `make test-docker` und `make verify-depguard` nicht gelaufen: kein Docker-Pfad und keine Importregel berührt (nur ein bool-Feld am Domain-Typ, keine neuen Imports).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Der Kommentar am `--quiet`-Flag sagt, `doctor` rendere unter `--quiet` nur Fehler- und Degradations-Fallback-Einträge; seit dem Diff kommen die Sicherheits-Warn-Einträge hinzu. Die Datei steht nicht im Diff und nicht in Slice §3. | `AGENTS.md` §3.9 (b) und §3.7 | `internal/adapter/driving/cli/root.go` · "doctor renders\n\t// only SeverityError items and degradation-fallback items" | nein — Review-only (`grep` nach dem Wortlaut im ganzen Baum) | Kommentar nennt geänderten Verhaltensumfang nicht mit |
| F-2 | LOW | Der Doc-Kommentar von `strictWarnFailure` sagt, `--quiet` verberge die Warn-Einträge auf stdout, daher sei die Meldung auf stderr der einzige Ort für den Grund. Für die gekennzeichneten Einträge stimmt das nicht mehr: sie stehen auf stdout und erklären den Exit 11. | `AGENTS.md` §3.9 (b) und §3.7 | `internal/adapter/driving/cli/doctor.go` · "`--quiet` hides those items from stdout, so the" | nein — Review-only | Kommentar nennt geänderten Verhaltensumfang nicht mit |
| F-3 | INFO | Der Lastenheft-Wortlaut nennt „Sicherheitslockerungen des Docker-Daemons“ als Warn-Einträge. Der Code hat dafür keinen Warn-Eintrag: `nestedRuntime: docker` ist ein OK-Eintrag, die Lockerung erscheint als Warn nur zusammen mit Egress. Der Slice benennt das ausdrücklich (Stufe OK bleibt); die Einstufung ist damit Ist-Stand, kein Fehler. | [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | `internal/hexagon/application/doctor_sandbox.go` · "Nested Docker daemon configured; it needs `--privileged`" | nein — Review-only | Lastenheft-Klasse ohne Warn-Eintrag im Code |
| F-4 | INFO | [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) (Zeile „`doctor` schreibt nur Fehler-Einträge und Einträge eines Degradations-Fallbacks“) und der Lastenheft-Satz „`u-boot doctor` zeigt nur Fehler-Einträge“ beschreiben den Zwischenstand. Der Spec-Slice der Welle trägt den Nachzug (Einstufung der `doctor`-Einträge unter `--quiet`); kein Befund, Hinweis an den Verifier. | [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) | `spec/spezifikation.md` · "`doctor` schreibt nur Fehler-Einträge und Einträge eines Degradations-Fallbacks" | nein — Verifier-Zuständigkeit | Spec-Nachzug offen, im Spec-Slice getragen |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Einstufung gegen den Code: `devcontainer.sandbox.credentials` (beide Warn-Zweige gekennzeichnet, OK-Zweige nicht) | geprüft, ohne Befund |
| `devcontainer.sandbox.egress` (beide Warn-Zweige gekennzeichnet: ohne Sandbox-Profil, mit `nestedRuntime: docker`) | geprüft, ohne Befund |
| `devcontainer.sandbox.runtime`: Warn „no effect“ ungekennzeichnet (Konfigurationsmeldung, wie Lastenheft), Warn `vfs` über `DegradationFallback` sichtbar, AppArmor-Sysctl und Politik `fail` sind Fehler, `docker` ist OK | geprüft, ohne Befund |
| `devcontainer.features.allowlist` (Verstoß ist `SeverityError`, bleibt unter `--quiet` sichtbar; Orphan/`enabled`-Lücke ist Warn und bleibt unterdrückt), `features.drift`, Host- und Compose-Checks | geprüft, ohne Befund (weder Sicherheitsbezug noch Überkennzeichnung) |
| Mechanismus `writeDoctorReport`: sichtbar = Fehler OR Fallback OR Kennzeichen; ohne `--quiet` unverändert (Summary-Zählung 0/2/1 im Test), Kopf und Summary entfallen, Hint bleibt in `writeDoctorItem` | geprüft, ohne Befund |
| JSON-Vertrag: Feld nicht in der Wire-Form; `TestDoctorJSON_SecurityRelaxationFlagNotInWireForm` vergleicht gekennzeichnet gegen ungekennzeichnet über drei Argument-Sätze plus Namensprüfung, nicht tautologisch (würde bei einer Wire-Änderung rot) | geprüft, ohne Befund |
| Exit-Codes und Dual-Classifier (HIGH #1 des Skills): kein neuer Sentinel, `ErrDoctorFailures` und `strictWarnFailure` unverändert; `doctor --quiet --strict` bleibt Exit 11, die Zeilen auf stdout und die IDs im Fehlertext erklären ihn (Test pinnt Exit-Code, Zeilen und IDs) | geprüft, ohne Befund |
| Wächterqualität, Fixture-Tausch in `..._FallbackItemAndErrorKept_OtherWarnsDropped` (Credentials durch `features.drift`): begründet im Testkommentar, die Negativ-Abdeckung für die gekennzeichneten Einträge liegt in den neuen Tests, nichts gelockert | geprüft, ohne Befund |
| Eigene Mutationen (alle rot): (1) Filter auf `Severity == Warn` statt Kennzeichen → sechs `cli`-Tests rot; (2) Kennzeichen am `nestedRuntime`-„no effect“ ergänzt → `TestDoctor_SandboxFallbackFlagIsOnTheItem` rot; (3) Kennzeichen am Egress-„no effect“ entfernt → derselbe Test rot. Alle Dateien per `git checkout -- <Datei>` zurückgesetzt, Baum sauber | geprüft, ohne Befund |
| Doku: Handbuch (`--quiet`-Zeilen, `doctor`-Abschnitt), Hilfetext „errors, degradation fallbacks and security warnings only“, CHANGELOG (nennt als folgend nur noch Agent-Credential-Übergabe samt neuem `doctor`-Check, Login-Dienst-Ausweisung und die Features; stimmt mit Slice und Lastenheft 0.11.0 überein, verschweigt nichts Geliefertes) | geprüft, ohne Befund |
| Plan-Tabelle Slice §3: jede berührte Datei steht darin, Fundstellen sind Datei und Symbol; Korrektur `cli_test.go` statt `quiet_test.go` ist vermerkt | geprüft, ohne Befund |
| Architekturgrenze (neues bool-Feld am Domain-Typ, keine neuen Imports), Suppression (kein `//nolint`), Sprachvertrag (alle neuen CLI-Texte Englisch), §3.9 (a) und (c) in den geänderten Kommentaren | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Kommentar nennt geänderten Verhaltensumfang nicht mit · Lastenheft-Klasse ohne Warn-Eintrag im Code · Spec-Nachzug offen, im Spec-Slice getragen

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM. Die zwei LOW-Findings sind Kommentar-Nachzüge nach `AGENTS.md` §3.9 (b) und kein Vertragsbruch.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7. Die DoD-Prüfung bleibt beim Verifier.
