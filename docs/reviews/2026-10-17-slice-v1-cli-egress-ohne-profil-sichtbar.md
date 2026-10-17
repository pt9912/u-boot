# Review-Report: slice-v1-cli-egress-ohne-profil-sichtbar — 2026-10-17 <!-- d-check:ignore (Slice-Kennung im Titel) -->

**Review-Art:** Code (gegen Plan, Spezifikation und Hard Rules)

**Gegenstand:** Slice [`slice-v1-cli-egress-ohne-profil-sichtbar`](../plan/planning/in-progress/slice-v1-cli-egress-ohne-profil-sichtbar.md), Commit-Range `3063cb1..d08146d` (Welle [`welle-devcontainer-sshd-claude-code`](../plan/planning/welle-devcontainer-sshd-claude-code.md), sechster Slice).

**Skill:** `.harness/skills/reviewer.md` @ d08146d · <!-- d-check:ignore (Pfad des Skills, Zitat des Laufs) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-17

**Eingangs-Kontext:**

- Slice-Plan (oben), insbesondere §2 DoD-Rahmen, §3 Plan-Tabelle
- [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Ausnahme 1, im Lastenheft ungeändert; trägt „dazu gehört die Meldung, dass die Egress-Restriktion ohne Sandbox-Profil nicht greift“), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)
- [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten)
- `AGENTS.md` §3.4, §3.7, §3.9, §3.10; `harness/review.md`; `harness/replay.md`

**Vorgehen (Sensors):** `make docs-check` (305 Dateien, 0 Befunde); `make test` als Mutationslauf, sechs Mutationen, danach jeweils `git checkout -- <Datei>`, Arbeitsbaum am Ende sauber (`git status` leer bis auf diesen Report). `make gates`, `make test-docker`, `make verify-depguard` nicht gelaufen: keine Importregel und kein Docker-Pfad berührt; `make gates` ist Sache des Implementers/Verifiers.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | INFO | Die Kennzeichnung wird im APP-Test am echten `GenerateService`, die Sichtbarkeit unter `--quiet` im CLI-Test mit einem Fake-Use-Case gebunden; es gibt keinen Test, der beides in einem Lauf durchreicht. Die Kopplung ist das Feld `SecurityRelaxation`, beide Hälften sind einzeln mutationsrot (M3, M4, M5), die Bindung genügt. | [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) Ausnahme 1 | `internal/adapter/driving/cli/quiet_test.go` · `TestExecute_Generate_Quiet_KeepsEgressNoEffectDropsNestedRuntimeNoEffect` | nein (Review-only) | Wächter-Hälften ohne gemeinsamen End-to-End-Lauf |
| F-2 | INFO | Zwei vom Implementer gemeldete Lücken bewertet: (a) die `--json`-Seite ist im neuen CLI-Test gedeckt (`--quiet --json` gleich `--json`, beide Meldungen vorhanden, kein „securityrelaxation“ im Text); die Meldung des Implementers trifft nur die Stellung, nicht die Abdeckung. (b) Der `docker`-Fall von `nestedRuntime` ist im CLI-Test nicht gespiegelt, im APP-Test (`TestGenerateDevcontainer_SandboxDocker_NotGenerated`) aber geprüft und durch M1 rot; das CLI-Verhalten hängt nur am Kennzeichen, nicht am Runtime-Wert. | [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) | `internal/hexagon/application/generate_sandbox_test.go` · `nestedRuntime no-effect warning is flagged` | ja — `make test` | Spiegelfall nur auf einer Schicht |
| F-3 | INFO | Mutation M6 (Tag `json:"-"` an `SecurityRelaxation` entfernt) bleibt grün. Das Kennzeichen erreicht die Wire-Form nicht über das Struct-Tag, sondern weil `mapWarningsToDiagnostics` nur Code, Level und Message abbildet; der Leak-Schutz im CLI-Test prüft das Ergebnis und ist damit nicht tautologisch, die Mutation ist äquivalent. | Maintainability | `internal/hexagon/port/driving/removeservice.go` · `SecurityRelaxation  bool   \`json:"-"\`` | ja — `make test` | äquivalente Mutation (Tag redundant zur Abbildung) |
| F-4 | INFO | Das Handbuch führt die `--quiet`-Zeile zweimal (Zeilen 299 und 1417) mit je einer Zelle von über 2000 Zeichen und verschachtelter Klammer; beide Zellen sind wortgleich nachgezogen und sachlich wahr, aber nur schwer lesbar und doppelt zu pflegen. Der Roadmap-Ruhe-Marker „Nichts in Arbeit.“ ist korrekt entfernt, solange der Slice in `in-progress/` liegt; die Wiederherstellung gehört in die Closure. | `AGENTS.md` §3.7 (Rang-Zeiger statt Wiederholung) | `docs/user/benutzerhandbuch.md` · `in der Befehlsausgabe erzeugt nur \`generate devcontainer\` diese Warnungen` | nein (Review-only) | Handbuch-Zelle doppelt und überlang |

## Prüfung der Schwerpunkte

**(1) Warnungsweg.** `SecurityRelaxation` wird in `visibleWarnings` (`internal/adapter/driving/cli/warnings.go`) ausgewertet, aufgerufen aus `generate.go`, `up.go` und `init.go`, sowie in `doctor.go` für Diagnostic-Einträge. Das Kennzeichen entsteht nur in `sandboxRelaxationWarnings`, dessen einziger Aufrufer `generateDevcontainer` ist; `init` (über `resolveSandboxClone`, Devcontainer-Konfiguration `nil`) und `up` (`recreateWarnings`) erzeugen keine Lockerungs-Warnungen, die Handbuch-Aussage „in der Befehlsausgabe erzeugt nur `generate devcontainer`“ ist am Code wahr. `--json` bildet über `mapWarningsToDiagnostics` ab und filtert nicht; das Feld ist kein Teil der Wire-Form (CLI-Test, M6-Analyse). Exit-Code und `--strict`: Warnungen beeinflussen den Exit-Code von `generate` nicht (kein `--strict`-Pfad in `generate.go`), `doctor --strict` ist unberührt. Meldungstext und Code der Egress-Meldung unverändert; die Replay-Fälle in `harness/replay.md` für `generate devcontainer` prüfen Dateiinhalte, nicht die Kennzeichnung, und sind nicht berührt.

**(2) Wächterqualität.** Meine Mutationen (`make test`, jeweils rot):

| Mutation | Rot durch |
|---|---|
| M1 `nestedRuntime=docker` „no effect“ als `relax` | APP: `…FlagIsOnTheEntryNotTheCode/docker-without-sandbox`, `…SameCodeBothFlags`, `…SandboxDocker_NotGenerated` |
| M2 `nestedRuntime=podman` „no effect“ als `relax` | APP: `…FlagIsOnTheEntryNotTheCode/podman-without-sandbox`, `…SandboxPodman_NotGenerated` |
| M3 Code der Egress-Meldung [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) zu [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) | APP: `…SandboxEgress_NotGenerated` (über den Code-Filter, mittelbar) |
| M4 `visibleWarnings` behält unter `--quiet` alles | CLI: `…Generate_Quiet_KeepsOnlyRelaxationWarnings`, `…Up_Quiet_…`, `…Init_Quiet_…` u. a. |
| M5 Egress-Lockerung mit Sandbox als `warn` | APP: `…FlagIsOnTheEntryNotTheCode/egress-sandbox`, `…PolicyChangesNoRenderTimeWarning` |
| M6 `json:"-"` entfernt | grün, äquivalent (F-3) |

Die Tests sind nicht tautologisch: Der APP-Test prüft das Kennzeichen am echten Service mit dem Code-und-Text-Filter, der CLI-Test die Filterung und die Abwesenheit des Kennzeichens in Text und Wire-Form.

**(3) Spec/Doku.** [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) führt die Egress-Meldung unter den gekennzeichneten Einträgen, `nestedRuntime` bleibt ungekennzeichnet; Historiezeile vorhanden. [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) unverändert richtig: Die Warnungstabelle führt nur Login-Dienst, Drittanbieter und `node`, die Egress-Meldung steht dort nicht. Handbuch (299, 1417), CHANGELOG und Spezifikation sagen dasselbe wie das ungeänderte Lastenheft; die Suche in Handbuch, README, README.de, CHANGELOG und `docs/maintainer/` nach „init“ im Zusammenhang mit Sandbox-Warnungen fand keine widersprechende Stelle. Der CHANGELOG-Eintrag behauptet nichts Ungeliefertes (Text und `--json` unverändert stimmen).

**(4) Plan-Tabelle gegen `git diff --stat 3063cb1..HEAD`.** Alle zehn Dateien des Diffs stehen in §3 (Code, drei Testdateien, Spezifikation, Handbuch, CHANGELOG, Roadmap, `open/README.md`, Slice selbst); kein Eintrag ohne Diff außer den als „lesen“ markierten. `open/README.md`: die drei entfernten Zeilen zeigten auf `done/` (zwei) und `in-progress/` (eine), lagen also nicht mehr in `open/`; die zwei verbleibenden Zeilen liegen in `open/`.

**(5) §3.9.** (a) kein „früher/former/no longer“ in neuen Kommentaren; (b) der Doc-Kommentar von `sandboxRelaxationWarnings` und der Test-Doc-Kommentar sind nachgezogen, die Suche nach „Only the relaxation“ und „stay unflagged“ im Baum fand keine veraltete Stelle; (c) kein Helfer zwischen Doc-Kommentar und Funktion eingefügt. `AGENTS.md` §3.4: `spec/architecture.md` unberührt.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/application/devcontainer_sandbox.go` | geprüft, ohne Befund |
| `internal/adapter/driving/cli/warnings.go`, `generate.go`, `up.go`, `init.go`, `doctor.go` | geprüft, ohne Befund |
| Dual-Classifier (Sentinels) | nicht berührt, kein neuer Sentinel |
| `internal/hexagon/application/*_test.go`, `cli/quiet_test.go` | geprüft, ohne Befund (siehe F-1 bis F-3 als INFO) |
| `spec/spezifikation.md` [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten), Historie | geprüft, ohne Befund |
| `CHANGELOG.md`, Handbuch Zeilen 299/1417 | geprüft, ohne Befund über INFO F-4 hinaus |
| `docs/plan/planning/` Roadmap, `open/README.md`, Plan-Tabelle | geprüft, ohne Befund |
| `harness/replay.md` Golden-Fälle | geprüft, nicht berührt |
| `make gates`, `make test-docker`, `make verify-depguard` | nicht gelaufen (Grund oben) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Wächter-Hälften ohne gemeinsamen End-to-End-Lauf · Spiegelfall nur auf einer Schicht · äquivalente Mutation (Tag redundant zur Abbildung) · Handbuch-Zelle doppelt und überlang

## Verdikt

**Merge-blockierend:** nein — keine HIGH-, MEDIUM- oder LOW-Findings.

**Übergabe:** Die INFO-Findings erfordern keine Änderung; F-2 verweist auf den Verifier (DoD-Abgleich, `make gates`). Dieser Report ersetzt keine Verifikation.
