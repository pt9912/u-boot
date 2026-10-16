# Review-Report: slice-v1-cli-quiet-nur-fehler — 2026-10-16

**Review-Art:** Code — Diff gegen Plan, [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Lastenheft 0.8.0), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) und Hard Rules (Maintainability; nicht gegen die DoD).

**Gegenstand:** Commit-Range `aca134c..HEAD` (Implementer-Commits `c277c47`, `72006d3`, `7a244b9`) für den Slice [`slice-v1-cli-quiet-nur-fehler`](../plan/planning/done/slice-v1-cli-quiet-nur-fehler.md) (Welle `welle-cli-verbosity`). <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

**Skill:** `.harness/skills/reviewer.md` @ HEAD `7a244b9` · <!-- d-check:ignore (Skill-Pfad als Zitat) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-16

**Eingangs-Kontext:**

- Slice [`slice-v1-cli-quiet-nur-fehler`](../plan/planning/done/slice-v1-cli-quiet-nur-fehler.md) (§1 Abgrenzung, §3 Plan) <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
- [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DIAG-003`](../../spec/lastenheft.md#lh-fa-diag-003--fehlerklassifikation), [`LH-NFA-PERF-002`](../../spec/lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker)
- `AGENTS.md` (Hard Rules §3.2, §3.7), `harness/conventions.md`, `harness/review.md`
- Sensors: `make lint` (Exit 0); Mutation `doctor.go` (`SeverityError` → `!= SeverityOK`) per `make test`: drei Doctor-Pins rot, Arbeitsbaum per `git checkout` zurückgesetzt (`git status` sauber). Nicht gelaufen: `make test-docker` (kein Docker-Pfad berührt), `make verify-depguard` (keine Importregel geändert; `warnings.go` importiert nur den Port `driving`, wie die Nachbardateien), `make gates` (Teilsensoren `lint`, `test` einzeln gelaufen).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | Der Kommentar am umgeschriebenen Doctor-Test beschreibt abwesenden Text: „Changed from "hides OK items only": the warn item and the summary line used to stay." Er erzählt die frühere Fassung statt der geltenden Zusage. Szenario: Der nächste Änderer des Tests liest die Chronik als Teil der Zusage und hält „Warn bleibt" für eine noch geltende Variante; kein Gate widerlegt es. | AGENTS.md §3.7 (Reviewer-Skill HIGH „Kommentar trägt keine der Kommentar-Klassen") | internal/adapter/driving/cli/cli_test.go · „the summary line used to stay.)" | nein — Review-only | Kommentar beschreibt abwesenden Text |
| F-2 | MEDIUM | `--quiet` ist Root-Persistent-Flag mit Hilfetext „show errors only", Handbuch und CHANGELOG sagen „nur Fehler". `init` und `generate` schreiben ihre Zusammenfassung weiter, `add` und `remove` ihre Erfolgsmeldung („Removed service …", „Changed:"), `remove` zusätzlich die Purge-Warnungen auf stderr ([`LH-FA-ADD-007`](../../spec/lastenheft.md#lh-fa-add-007--service-entfernen)); `remove` liest `a.quiet` nicht. Szenario: Ein Skript mit `u-boot init --quiet`, `add --quiet` oder `remove --quiet` erwartet nach [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) („`--quiet` – nur Fehler", „Alle anderen Hinweise entfallen") leere Ausgabe und bekommt Status- und Warntext. Der Slice schließt `init`/`generate`-Zusammenfassung und `remove` per §1 aus; ein Slice-Plan (Rang 5) setzt das Lastenheft (Rang 1) nicht außer Kraft. Handbuch und Hilfetext nennen diese Ausnahmen nicht (Doku-Drift). | [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | internal/adapter/driving/cli/root.go · „show errors only; security relaxations of the sandbox profile stay reported"; internal/adapter/driving/cli/remove.go · „WARNING: %s"; internal/adapter/driving/cli/init.go · „Initialized" | ja — `make test` (Pin `remove --quiet` / `init --quiet` ohne Statuszeile würde heute rot) | `--quiet` nicht an allen Statuspfaden durchgesetzt |
| F-3 | MEDIUM | `doctor --quiet --strict` mit nur Warn-Items endet mit Exit 11 ([`LH-FA-DIAG-003`](../../spec/lastenheft.md#lh-fa-diag-003--fehlerklassifikation), `--strict`), die Ausgabe ist leer; auf stderr steht nur „u-boot: doctor report contains failures". Es gibt keinen Test für die Kombination. Szenario: Ein CI-Skript ruft `doctor --quiet --strict` auf, der Schritt scheitert mit 11, und weder stdout noch stderr nennen das auslösende Warn-Item. Der Exit-Code selbst ist unverändert (der Diff berührt weder `doctorExitCode` noch `ErrDoctorFailures`), die Erklärbarkeit nicht gepinnt. | [`LH-FA-DIAG-003`](../../spec/lastenheft.md#lh-fa-diag-003--fehlerklassifikation), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) | internal/adapter/driving/cli/doctor.go · „if item.Severity == domain.SeverityError {" | ja — `make test` (fehlender Pin `--quiet --strict`) | Exit 11 ohne sichtbare Ursache, Kombination ungepinnt |
| F-4 | LOW | Der Kommentar in `root.go` verweist auf `[writeWarning]`; das Symbol existiert nirgends (der Helfer heißt `visibleWarnings` in `warnings.go`). | AGENTS.md §3.7 (Rang-Zeiger) | internal/adapter/driving/cli/root.go · „remain, see [writeWarning])" | nein — Review-only | Rang-Zeiger auf nicht existierendes Symbol |
| F-5 | INFO | Ausnahme-Umfang: `sandboxRelaxationWarnings` kennzeichnet alle Lockerungen, die `devcontainer_sandbox.go` in `RunArgs` setzt: Podman fünf (`--cap-add=SYS_ADMIN`, `seccomp`/`apparmor`/`systempaths=unconfined`, `--device=/dev/fuse`), Docker `--privileged`, Egress `--cap-add=NET_ADMIN`. Es gibt keine gesetzte, aber ungewarnte und keine gewarnte, aber ungekennzeichnete Lockerung (Gegenprobe: Compose-/Dockerfile-Templates tragen keine). Die Degradations-Fallbacks (Daemon, `vfs`, Wegfall Egress) entstehen beim Containerstart und in `doctor`, nicht als Render-Warnung; ihr Ausschluss aus der Ausnahme folgt der aufgezählten Klasse des Lastenhefts (Capabilities, `--privileged`, Seccomp/AppArmor, Devices). [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) („jeder Fallback wird in der Befehlsausgabe und in `u-boot doctor` ausgewiesen") und CLI-005 („alle anderen Hinweise entfallen") stehen in Spannung; der Wegfall der Egress-Restriktion ist ein Sicherheitsverlust, aber keine Lockerung im Sinn der Aufzählung. Entscheidung gehört dem Architect/Projektinhaber, kein Code-Befund. | [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) | internal/hexagon/application/devcontainer_sandbox.go · „are not produced here: they arise at container start" | Review-only | Spec-Spannung Fallback-Ausweisung vs. `--quiet` |
| F-6 | INFO | Der Compose-Fortschritt auf stderr bleibt unter `up --quiet` bestehen (Kommentar in `up.go`: „requires phase visibility regardless"). Das ist ein Widerspruch zweier Lastenheft-Sätze („nur Fehler" gegen „Fortschritt einzelner Services muss sichtbar sein"), kein Code-Fehler; das Lastenheft löst ihn nicht auf. Zuständig: Architect (Spec-Klärung, ggf. CR). | [`LH-NFA-PERF-002`](../../spec/lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | internal/adapter/driving/cli/up.go · „requires phase visibility regardless" | Review-only | Lastenheft-Widerspruch Fortschritt vs. quiet |
| F-7 | INFO | Der Ausnahmezweig in `up` ist heute ohne Quelle: `UpService` erzeugt nur die Recreate-Warnung, keine Lockerung; der Zweig ist nur durch Fake-Use-Cases gepinnt (Slice §1 nennt das ausdrücklich). Er wirkt erst, wenn `up` künftig gekennzeichnete Warnungen erzeugt. Zuständig: Verifier (DoD-Formulierung „bleiben bei `up` sichtbar"). | [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | internal/adapter/driving/cli/up.go · „for _, w := range visibleWarnings(resp.Warnings, flags.Quiet) {" | Review-only | Ausnahmepfad ohne produktive Quelle |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Dual-Classifier / Exit-Codes (Reviewer-Skill HIGH #1) | geprüft, ohne Befund: kein neuer Sentinel, `ExitCode`, `doctorExitCode` und `ErrDoctorFailures` unverändert; die Pegel- und Filteränderung berührt nur die Textausgabe. F-3 betrifft Erklärbarkeit, nicht den Code. |
| JSON-Vertrag ([`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe)) | geprüft, ohne Befund: in `init` und `generate` läuft der JSON-Zweig vor dem Filter, `up` filtert nur den Human-Pfad, `doctor --json` ist unverändert; `SecurityRelaxation` trägt `json:"-"`. `TestExecute_QuietJSON_KeepsAllWarningsAsDiagnostics` zählt vier `warn`-Einträge je init/generate/up und sucht den Feldnamen im Wire-Format; `doctor_json_test.go` pinnt `--json --quiet` gegen `--json`. |
| Tests/Wächter | geprüft, ohne Befund: Die Doctor-Mutation wurde von drei Pins gefangen (nicht tautologisch). Geänderte Bestandspins (`verbosity_test.go` Warn→Error, `cli_test.go` Doctor-Quiet) sind Soll-Änderungen nach Lastenheft 0.8.0, keine Lockerung; `debug-takes-precedence-over-quiet` bleibt unverändert. Die Kennzeichnung hängt am Eintrag, nicht am Code (Test `SameCodeBothFlags`). |
| Suppression, Architekturgrenzen, Sprache | geprüft, ohne Befund: kein `//nolint` im Diff; `warnings.go` und das Port-Feld bringen keine neuen Import-Richtungen; neue CLI-Texte (Hilfetext, Warnzeilen) Englisch ([`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache)). |
| Handbuch/CHANGELOG gegen Lastenheft 0.8.0 | geprüft, ohne Befund über F-2 hinaus: Ausnahme und `doctor --quiet` korrekt beschrieben; die Annahme „nur Fehler" ohne Einschränkung ist F-2. |
| Referenzmodell, Lifecycle, Carveouts | nicht berührt: der Diff ändert keine Spec-, ADR- oder Planning-Dateien. |
| `make gates`, `make test-docker`, `make verify-depguard` | nicht vollständig ausgeführt: siehe Eingangs-Kontext (`make lint` und die Mutation über `make test` liefen; kein Docker-Pfad, keine Importregel geändert). |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 2 |
| LOW | 1 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Kommentar beschreibt abwesenden Text · `--quiet` nicht an allen Statuspfaden durchgesetzt · Exit 11 ohne sichtbare Ursache, Kombination ungepinnt · Rang-Zeiger auf nicht existierendes Symbol · Spec-Spannung Fallback-Ausweisung vs. quiet · Lastenheft-Widerspruch Fortschritt vs. quiet · Ausnahmepfad ohne produktive Quelle

## Verdikt

**Merge-blockierend:** ja — F-1 (HIGH) und F-2/F-3 (MEDIUM) sind vor Merge zu klären. F-2 ist zusätzlich eine Abgrenzungsfrage an den Architect: Hält der Slice-Ausschluss von `init`/`generate`/`add`/`remove` gegen Lastenheft 0.8.0, oder braucht das Lastenheft eine Präzisierung (Ranghöhere Quelle gewinnt bis dahin).

**Übergabe:** Findings gehen an den Implementer; F-2, F-5 und F-6 gehen zusätzlich als Spec-Klärungsfragen an den Architect (Konflikt-Pfad bei HIGH-Widerspruch als Sequenz mit Übergabe-Artefakt); die Finding-Klassen gehen in die Slice-Closure §7. Dieser Report ersetzt keine Verifikation.
