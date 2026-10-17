# Slice v1: Egress „has no effect“ bleibt unter `--quiet` in `generate devcontainer` sichtbar

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-devcontainer-sshd-claude-code`](../welle-devcontainer-sshd-claude-code.md) — Folge-Slice nach dem Spezifikations-Slice [`slice-v2-spec-devcontainer-sshd-claude-code`](../done/slice-v2-spec-devcontainer-sshd-claude-code.md) (Review F-1, Verification D-2); die Welle schließt erst nach ihm (Closure-Bedingung, die von der DoD dieses Slice verschieden ist: Lastenheft, Code, Handbuch, CHANGELOG und Spezifikation sagen über `--quiet` dasselbe).

**Bezug:** [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Ausnahme 1: „dazu gehört die Meldung, dass die Egress-Restriktion ohne Sandbox-Profil nicht greift“), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Egress-Restriktion), [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (die Meldung zu `nestedRuntime` ist eine Konfigurationsmeldung). Entscheidung des Projektinhabers 2026-10-17: Code-Folge-Slice, kein Lastenheft-CR. Review: [`2026-10-17-slice-v2-spec-devcontainer-sshd-claude-code`](../../../reviews/2026-10-17-slice-v2-spec-devcontainer-sshd-claude-code.md) (F-1); Verification Evidence: [`2026-10-17-verification-slice-v2-spec-devcontainer-sshd-claude-code`](../../../reviews/2026-10-17-verification-slice-v2-spec-devcontainer-sshd-claude-code.md) (D-2, H-1, H-2).

**Berührte Spec-Stellen:** `spezifikation.md`: [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) (Einstufungsliste, Ausnahme 1), [`SPEC-044`](../../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) (Warnungstabelle: Kennzeichen der Egress-Meldung, falls sie dort geführt wird — beim Schreiben prüfen); `lastenheft.md` und `architecture.md` bleiben unberührt (das Lastenheft trägt die Zusage bereits).

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-17.

---

## 1. Ziel und Abgrenzung

**Ziel:** `u-boot generate devcontainer --quiet` zeigt die Meldung „Egress-Restriktion ohne Sandbox-Profil ohne Wirkung“, und `doctor --quiet` tut es schon: Ein Nutzer, der `egress.enabled: true` mit Profil `default` erzeugt, erfährt unter `--quiet`, dass keine Restriktion entstand.

Ist-Zustand (Lesestand 2026-10-17, Fundstelle `sandboxRelaxationWarnings` in `internal/hexagon/application/devcontainer_sandbox.go`): Die Meldung `devcontainer.sandbox.egress.enabled has no effect without devcontainer.profile=sandbox; …` entsteht über `warn(…)` (Code der Anforderung Egress-Restriktion), also ohne `SecurityRelaxation`, und entfällt unter `--quiet`; der `doctor`-Eintrag `devcontainer.sandbox.egress` für denselben Fall trägt das Kennzeichen (`checkSandboxEgress` in `doctor_sandbox.go`, `SecurityRelaxation: true`). Der einzige Aufrufer von `sandboxRelaxationWarnings` ist `generateDevcontainer` (`generate.go`, Zeile mit `warnings = append(warnings, sandboxRelaxationWarnings(…)…)`).

**Soll:**

- Die Meldung zu `egressEnabled(dc)` ohne Sandbox-Profil trägt `SecurityRelaxation: true`. Die Meldungen mit Sandbox-Profil (`security relaxation for the egress restriction: …`) sind schon gekennzeichnet und bleiben.
- Die Meldungen „has no effect“ zu `nestedRuntime` (`podman`, `docker`) bleiben ungekennzeichnet; sie sind eine Konfigurationsmeldung ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)) und entfallen unter `--quiet`.
- Der Text der Meldungen bleibt (Golden- und Replay-Fälle unverändert); geändert wird allein das Kennzeichen.
- Die Spezifikation und das Handbuch nennen den Zustand nach dem Code: [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) führt die Egress-Meldung von `generate devcontainer` unter den gekennzeichneten Einträgen, die ungekennzeichneten „has no effect“-Meldungen nennen nur noch `nestedRuntime`; das Handbuch nennt als Erzeuger der Warnungen zu Sicherheitslockerungen nur `generate devcontainer` (Verification H-1: `init` ruft `sandboxRelaxationWarnings` nicht auf, `initTemplateData` übergibt `nil`; `up` erzeugt keine Lockerung).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Lastenheft-Klammern** („Passwort-Login“, „Sicherheitslockerungen des Docker-Daemons“ als `doctor`-Eintrag, Login-Dienst nur im Sandbox-Profil; Review F-6 bis F-8) — Zuständigkeit des Projektinhabers im Lastenheft, kein Code-Gegenstand; eine Klarstellung wäre ein eigener CR.
- **Änderung der Meldungstexte oder des Codes [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)** — die Texte sind Vertrag der Golden-Fälle; der Slice ändert das Kennzeichen, nicht den Wortlaut.
- **Die `doctor`-Einträge** — sie tragen das Kennzeichen bereits; `nestedRuntime: docker` hat dort die Stufe OK, und das bleibt (Review F-7).
- **Neue Erzeuger der Warnungen in `init` oder `up`** — `init` ruft `sandboxRelaxationWarnings` nicht auf; ob es sie aufrufen soll, ist eine Produktentscheidung und kein Teil des Kennzeichens.

## 2. Definition of Done

- [x] **Code:** `sandboxRelaxationWarnings` kennzeichnet die Egress-Meldung ohne Sandbox-Profil mit `SecurityRelaxation`; die `nestedRuntime`-Meldungen „has no effect“ bleiben ungekennzeichnet; der Test `TestSandboxRelaxationWarnings_FlagIsOnTheEntryNotTheCode` bildet die neue Einteilung ab (Egress ohne Profil: gekennzeichnet; `podman` und `docker` ohne Profil: ungekennzeichnet), und ein CLI-Test belegt, dass `generate devcontainer --quiet` die Egress-Meldung zeigt und die `nestedRuntime`-Meldung unterdrückt ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)).
- [x] **Spezifikation und Doku:** [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) (Einstufungsliste) und, wo die Egress-Meldung dort steht, [`SPEC-044`](../../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) nennen den Zustand; `docs/user/benutzerhandbuch.md` (die beiden `--quiet`-Zeilen, `generate devcontainer` als einziger Erzeuger der Lockerungs-Warnungen, die Egress-Meldung von `generate devcontainer` als sichtbare Ausnahme) und `CHANGELOG.md` (`[Unreleased]`, Eintrag zum `--quiet`-Vertrag) sagen dasselbe wie Code und Spezifikation.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update, falls ein öffentlicher Vertrag berührt ist (hier: Handbuch und CHANGELOG, siehe oben).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne Wellen-Betrieb hier geprüft, im Repo mit Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Lesestand 2026-10-17; der Implementer liest jede Zeile vor der Änderung an der Stelle im Code (`AGENTS.md` §3.10) und gleicht die Tabelle nach der letzten Fix-Runde gegen `git diff --stat <Basis>..HEAD` ab, nicht nur vor der ersten Review.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/application/devcontainer_sandbox.go` | update | Zweig `egressEnabled(dc)` ohne Sandbox: `relax(…)` statt `warn(…)` für die Meldung `… has no effect …`; Doc-Kommentar von `sandboxRelaxationWarnings` („Only the relaxation warnings carry …“, „a „no effect“ notice per feature that is configured without the sandbox profile“) nennt die Egress-Meldung als gekennzeichnet (`AGENTS.md` §3.7 und §3.9 (b): Verweise auf Geändertes nachziehen, auch im ganzen Baum suchen) |
| `internal/hexagon/application/devcontainer_sandbox_relaxation_test.go` | update | `TestSandboxRelaxationWarnings_FlagIsOnTheEntryNotTheCode`: Fall `egress-without-sandbox` von 0 gekennzeichnet / 1 ungekennzeichnet auf 1 / 0; die Prüfung „gekennzeichneter Eintrag enthält `security relaxation`“ gilt nicht mehr für ihn (gekennzeichnet ist er mit „has no effect“-Text), der Test unterscheidet ihn; der Doc-Kommentar des Tests (Aussage, die „has no effect“-Meldungen mit demselben Code wie die Lockerungen blieben ungekennzeichnet) wird berichtigt; `TestSandboxRelaxationWarnings_SameCodeBothFlags` und `…PolicyChangesNoRenderTimeWarning` bleiben (nachlesen) |
| `internal/adapter/driving/cli/quiet_test.go` | update | neuer Fall für `generate devcontainer --quiet`: gekennzeichnete Egress-Meldung neben ungekennzeichneter `nestedRuntime`-Meldung — die erste bleibt, die zweite entfällt; nach dem Muster von `TestExecute_Quiet_KeepsSshdLoginWarningDropsThirdParty` (Hilfsfunktionen `quietNoEffect`, `quietRun`, `newAppWithGenerate`); nachzulesen, ob `quiet_test.go` oder `verbosity_last_wins_test.go` (Test `TestExecute_Verbosity_LastWins_WarningFilter`, prüft `has no effect` als „nicht-gekennzeichnet“) die Aufteilung ebenfalls trägt |
| `internal/hexagon/application/generate_sandbox_test.go` | update | `TestGenerateDevcontainer_SandboxEgress_NotGenerated` prüft am echten `GenerateService`, dass die Egress-Meldung `SecurityRelaxation` trägt; `…SandboxPodman_NotGenerated` und `…SandboxDocker_NotGenerated` prüfen, dass die `nestedRuntime`-Meldung „has no effect“ ungekennzeichnet bleibt (die Fixture-Fakes des CLI-Tests setzen das Kennzeichen selbst und binden den Code nicht) |
| `docs/plan/planning/in-progress/roadmap.md` | update | der Ruhe-Marker „Nichts in Arbeit.“ in „Offene Wellen“ entfällt, solange der Slice in `in-progress/` liegt (d-check `planning`) |
| `docs/plan/planning/open/README.md` | update | die drei Bestandszeilen zu Slices, die nicht mehr in `open/` liegen (`slice-v1-devcontainer-feature-claude-code`, `slice-v2-spec-devcontainer-sshd-claude-code`, dieser Slice), entfallen |
| `internal/hexagon/application/export_test.go` | lesen | `SandboxRelaxationWarningsForTest` bleibt, wenn sich die Signatur nicht ändert |
| `spec/spezifikation.md` ([`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), gegebenenfalls [`SPEC-044`](../../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten), §7 Historie) | update | Einstufungsliste: Egress-Meldung von `generate devcontainer` unter „gekennzeichnet“, „ungekennzeichnet“ nennt für diesen Befehl nur noch `nestedRuntime`; Historie-Zeile |
| `docs/user/benutzerhandbuch.md` | update | die zwei `--quiet`-Tabellenzeilen (Zeilen 299 und 1417 beim Lesestand): „die übrigen Lockerungen je nach Befehl `init`, `generate`, `up`“ wird zu `generate devcontainer` als einzigem Erzeuger (Verification H-1), und die Egress-Meldung ohne Sandbox-Profil steht in der Befehlsausgabe wie in `doctor` unter den sichtbaren Meldungen; die Zeilen 611 bis 613 (`doctor`) bleiben, soweit sie nur `doctor` beschreiben |
| `CHANGELOG.md` | update | `[Unreleased]`, Abschnitt Changed, Eintrag zum `--quiet`-Vertrag (Lesestand: Zeilen 38 bis 61): die Egress-Meldung von `generate devcontainer` bleibt unter `--quiet` sichtbar, die `nestedRuntime`-Meldung nicht; dort nennt der Eintrag bisher nur `doctor` für diese Meldung |
| `harness/replay.md`, Golden-Fälle | lesen | der Wortlaut der Meldungen ändert sich nicht; zu prüfen, dass kein Golden-Fall das Kennzeichen trägt |

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; der Spezifikations-Slice [`slice-v2-spec-devcontainer-sshd-claude-code`](../done/slice-v2-spec-devcontainer-sshd-claude-code.md) liegt in `done/`; kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Slice überschreitet 3 Liefer-Punkte, etwa weil ein Golden-Fall oder die JSON-Ausgabe (`warnings[]`) das Kennzeichen sichtbar trägt und mitzieht.
- `in-progress` → `open` (blockiert): der Projektinhaber entscheidet im Lastenheft gegen die Lesart „Egress-Meldung gehört in der Befehlsausgabe zu den sichtbaren“ (CR statt Code).

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Verification-Evidence und Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7). Mit diesem Slice ist die Welle [`welle-devcontainer-sshd-claude-code`](../welle-devcontainer-sshd-claude-code.md) schließbar.

## 6. Risiken und offene Punkte

- **Das Kennzeichen erscheint in `--json` oder in einem Golden-Fall** und ändert eine maschinenlesbare Ausgabe. Der Implementer liest `WarningEntry` im JSON-Umschlag und die Golden-Fälle vor der Änderung; trägt eine Ausgabe das Kennzeichen, geht der Slice nach `next` (§4). — **Ausgang:** entfallen: Das Kennzeichen erreicht keine maschinenlesbare Ausgabe; `TestExecute_Generate_Quiet_KeepsEgressNoEffectDropsNestedRuntimeNoEffect` (`quiet_test.go`) belegt `--quiet --json` gleich `--json` ohne Kennzeichen im Text, `mapWarningsToDiagnostics` bildet nur Code, Level und Message ab; die Replay-Fälle prüfen Dateiinhalte, kein Golden-Fall trägt das Kennzeichen.
- **Plan-Aussage über Code stimmt nicht** (Klasse [BEO-DOK/plan-aussage-ueber-code](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md), verkörpert in `AGENTS.md §3.10`; die Plan-Tabelle altert mit jedem Commit nach dem Abgleich): Der Abgleich der Tabelle gegen `git diff --stat` wird nach der letzten Fix-Runde wiederholt. — **Ausgang:** nicht eingetreten: Die Plan-Tabelle wurde nach der letzten Fix-Runde (`d08146d`) gegen `git diff --stat 3063cb1..HEAD` abgeglichen und vom Verifier erneut gelesen; keine Abweichung (jede berührte Datei hat einen Eintrag, jeder Änderungs-Eintrag einen Diff). Register: [BEO-DOK/plan-aussage-ueber-code](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md), Zähler unverändert.

## 7. Closure-Notiz

**Belege:** Review-Report [`2026-10-17-slice-v1-cli-egress-ohne-profil-sichtbar`](../../../reviews/2026-10-17-slice-v1-cli-egress-ohne-profil-sichtbar.md) (0 HIGH, 0 MEDIUM, 0 LOW, 4 INFO; Mutationslauf M1 bis M6, fünf rot, M6 äquivalent). Verification Evidence [`2026-10-17-verification-slice-v1-cli-egress-ohne-profil-sichtbar`](../../../reviews/2026-10-17-verification-slice-v1-cli-egress-ohne-profil-sichtbar.md): DoD-Liefer-Punkte erfüllt, keine Diskrepanz zwischen Plan, Code, Spezifikation, Handbuch, CHANGELOG und Lastenheft; Auflagen 1 bis 6 in der Closure erfüllt bzw. beachtet. `make gates` grün, `make doc-trace`: 83 Anforderungen, 0 Waisen, `make docs-check`: 0 Befunde. Commits `12a1aae`, `fdbb555`, `88c03b9`, `d08146d`.

**Fundstellen je Aussage (`AGENTS.md` §3.10):**

| Aussage | Fundstelle |
|---|---|
| Die Egress-Meldung ohne Sandbox-Profil trägt `SecurityRelaxation`, die `nestedRuntime`-Meldungen „has no effect“ nicht | `sandboxRelaxationWarnings` (`relax` im Zweig `egressEnabled(dc)`, `warn` für `nestedRuntime`) in `internal/hexagon/application/devcontainer_sandbox.go` |
| Unter `--quiet` bleibt die gekennzeichnete Meldung, die ungekennzeichnete entfällt; `--quiet --json` gleich `--json` | `TestExecute_Generate_Quiet_KeepsEgressNoEffectDropsNestedRuntimeNoEffect` in `internal/adapter/driving/cli/quiet_test.go`; Filter `visibleWarnings` in `internal/adapter/driving/cli/warnings.go` |
| Das Kennzeichen erreicht die Wire-Form nicht | `mapWarningsToDiagnostics` (bildet Code, Level, Message ab) im CLI-Adapter; Leak-Prüfung im selben CLI-Test |
| Kennzeichnung am echten Service | `TestGenerateDevcontainer_SandboxEgress_NotGenerated`, `…SandboxPodman_NotGenerated`, `…SandboxDocker_NotGenerated` in `internal/hexagon/application/generate_sandbox_test.go` |
| Einteilung je Fall | `TestSandboxRelaxationWarnings_FlagIsOnTheEntryNotTheCode` in `internal/hexagon/application/devcontainer_sandbox_relaxation_test.go` |
| Spezifikation und Doku sagen dasselbe | [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) (Einstufungsliste, Historie) in `spec/spezifikation.md`; `--quiet`-Zeilen und Erzeuger-Satz in `docs/user/benutzerhandbuch.md`; `[Unreleased]` in `CHANGELOG.md`; `doctor` unverändert (`checkSandboxEgress` in `doctor_sandbox.go`) |

- **Was hat funktioniert:** Die Plan-Tabelle wurde nach der letzten Fix-Runde gegen `git diff --stat` abgeglichen (`d08146d`), die Verification fand keine Abweichung; das Kennzeichen am Eintrag statt am Code (Test `…FlagIsOnTheEntryNotTheCode`) trug die Einteilung ohne Textänderung, sodass Golden- und Replay-Fälle unberührt blieben. Der Wächter ist auf die „no effect“-Meldung beschränkt (`88c03b9`); das `nestedRuntime`-Verhalten ist durch beide Schichten einzeln mutationsrot belegt.
- **Was ging anders als geplant:** Nichts Wesentliches; der Slice blieb bei den geplanten Liefer-Punkten (Code, Spezifikation und Doku). Der Commit `88c03b9` engte den Wächter auf die „no effect“-Meldung ein.
- **Steering-Loop-Eintrag:** **Gezählt, nicht verkörpert; kein neuer Eintrag.** `BEO-DOK/plan-aussage-ueber-code`: die Klasse trat in diesem Vorgang **nicht** ein (Plan-Tabelle gegen Diff ohne Abweichung); Zähler unverändert. `BEO-CLI/kommentar-beschreibt-abwesendes`: Klasse nicht eingetreten (`AGENTS.md` §3.9-Selbstcheck laut Review ohne Befund, auch am geänderten Doc-Kommentar von `sandboxRelaxationWarnings`); Zähler unverändert. **Lehre:** Beide verkörperten Regeln (§3.9, §3.10) hielten in einem Vorgang, in dem sie vorab im Plan benannt waren (§6, §8); ein einzelner Vorgang ohne Auftreten ist kein Beleg für Wirksamkeit und ändert den Stand beider Beobachtungen nicht (`state.md` unverändert; die Entscheidung zu `BEO-DOK/plan-aussage-ueber-code` Sensor oder Begründung steht beim Projektinhaber, spätestens bei der Welle-Closure, und wird hier nicht entschieden). Kein `liegt in`, nichts verkörpert.
- **Beobachtungs-Register (`../observations/`):** keine neue Evidence-Datei, keine neue Beobachtung. Review F-4 (Handbuch-`--quiet`-Zelle in Zeilen 299 und 1417 doppelt und überlang) wird nicht ins Register aufgenommen: Es ist ein einmaliger Review-only-INFO-Befund einer kosmetischen Doku-Eigenschaft, der aus dem Bestand stammt, ohne Wiederauftreten und ohne Klasse, die einen zweiten Beleg zählen ließe; das Register zählt Wiederauftreten einer Klasse, ein Einzelfall wäre ein Zähler 1 ohne Aussage. Er geht als **Hinweis an den Projektinhaber** (kein Folge-Slice): Ob die Zelle gekürzt und die Doppelung durch einen Verweis ersetzt werden soll (`AGENTS.md` §3.7, Rang-Zeiger statt Wiederholung), ist seine Entscheidung; ein Folge-Slice nach `open/` entsteht nur auf sie. Tritt eine zweite überlange oder doppelte Handbuch-Zelle auf, wird die Klasse dann im Register geführt.
- **Folge-Slices:** keine. Die Hinweise F-1 (kein gemeinsamer End-to-End-Lauf von Anwendung und CLI; beide Hälften einzeln mutationsrot), F-2 und F-3 (äquivalente Mutation) erfordern keine Änderung.
- **Risiken aus §6:** beide tragen einen Ausgang: entfallen (Kennzeichen in `--json` oder Golden-Fall) und nicht eingetreten (Plan-Aussage über Code).
- **Roadmap:** Der Ruhe-Marker „Nichts in Arbeit.“ bleibt entfernt, weil die Welle [`welle-devcontainer-sshd-claude-code`](../welle-devcontainer-sshd-claude-code.md) noch unter „Offene Wellen“ steht; er kehrt mit ihrem Schluss zurück (Verifier-Auflage 1, d-check `planning`).
- **Drei Paarungen:** Anker (kein `liegt in`, nichts verkörpert) · Folge-Slice (keiner) · Register (keine neue Evidence-Datei, Zähler unverändert) — nach dem `git mv` nach `done/` geprüft: nichts zu paaren.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `internal/hexagon/application` (APP, Code und Test), `internal/adapter/driving/cli` (CLI, nur Test) und, für Doku und Spezifikation, `docs/`, `spec/` (DOK) — alle Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006) (Abschnitt Modus-Deklaration), jede erfüllt die Schwelle ≥ 2 von 3 Achsen; keine Ausdifferenzierung nötig. Weil CLI berührt wird, gilt `AGENTS.md §3.9` (Kommentar-Selbstcheck) für den Diff, `AGENTS.md §3.10` für jede Plan-Aussage über Code.

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register ([`../observations/`](../observations/README.md)) durchgegangen — Treffer: [BEO-DOK/plan-aussage-ueber-code](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (DOK, verkörpert in `AGENTS.md §3.10`; Plan-Tabelle gegen den Diff nach der letzten Fix-Runde erneut abgleichen) und `BEO-CLI/kommentar-beschreibt-abwesendes` (CLI, verkörpert in `AGENTS.md §3.9`; der Doc-Kommentar von `sandboxRelaxationWarnings` ändert sich).

Alle berührten Sub-Areas GF.
