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

- [ ] **Code:** `sandboxRelaxationWarnings` kennzeichnet die Egress-Meldung ohne Sandbox-Profil mit `SecurityRelaxation`; die `nestedRuntime`-Meldungen „has no effect“ bleiben ungekennzeichnet; der Test `TestSandboxRelaxationWarnings_FlagIsOnTheEntryNotTheCode` bildet die neue Einteilung ab (Egress ohne Profil: gekennzeichnet; `podman` und `docker` ohne Profil: ungekennzeichnet), und ein CLI-Test belegt, dass `generate devcontainer --quiet` die Egress-Meldung zeigt und die `nestedRuntime`-Meldung unterdrückt ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)).
- [ ] **Spezifikation und Doku:** [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) (Einstufungsliste) und, wo die Egress-Meldung dort steht, [`SPEC-044`](../../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) nennen den Zustand; `docs/user/benutzerhandbuch.md` (die beiden `--quiet`-Zeilen, `generate devcontainer` als einziger Erzeuger der Lockerungs-Warnungen, die Egress-Meldung von `generate devcontainer` als sichtbare Ausnahme) und `CHANGELOG.md` (`[Unreleased]`, Eintrag zum `--quiet`-Vertrag) sagen dasselbe wie Code und Spezifikation.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update, falls ein öffentlicher Vertrag berührt ist (hier: Handbuch und CHANGELOG, siehe oben).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne Wellen-Betrieb hier geprüft, im Repo mit Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Lesestand 2026-10-17; der Implementer liest jede Zeile vor der Änderung an der Stelle im Code (`AGENTS.md` §3.10) und gleicht die Tabelle nach der letzten Fix-Runde gegen `git diff --stat <Basis>..HEAD` ab, nicht nur vor der ersten Review.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/application/devcontainer_sandbox.go` | update | Zweig `egressEnabled(dc)` ohne Sandbox: `relax(…)` statt `warn(…)` für die Meldung `… has no effect …`; Doc-Kommentar von `sandboxRelaxationWarnings` („Only the relaxation warnings carry …“, „a „no effect“ notice per feature that is configured without the sandbox profile“) nennt die Egress-Meldung als gekennzeichnet (`AGENTS.md` §3.7 und §3.9 (b): Verweise auf Geändertes nachziehen, auch im ganzen Baum suchen) |
| `internal/hexagon/application/devcontainer_sandbox_relaxation_test.go` | update | `TestSandboxRelaxationWarnings_FlagIsOnTheEntryNotTheCode`: Fall `egress-without-sandbox` von 0 gekennzeichnet / 1 ungekennzeichnet auf 1 / 0; die Prüfung „gekennzeichneter Eintrag enthält `security relaxation`“ gilt nicht mehr für ihn (gekennzeichnet ist er mit „has no effect“-Text), der Test unterscheidet ihn; der Doc-Kommentar des Tests (Aussage, die „has no effect“-Meldungen mit demselben Code wie die Lockerungen blieben ungekennzeichnet) wird berichtigt; `TestSandboxRelaxationWarnings_SameCodeBothFlags` und `…PolicyChangesNoRenderTimeWarning` bleiben (nachlesen) |
| `internal/adapter/driving/cli/quiet_test.go` | update | neuer Fall für `generate devcontainer --quiet`: gekennzeichnete Egress-Meldung neben ungekennzeichneter `nestedRuntime`-Meldung — die erste bleibt, die zweite entfällt; nach dem Muster von `TestExecute_Quiet_KeepsSshdLoginWarningDropsThirdParty` (Hilfsfunktionen `quietNoEffect`, `quietRun`, `newAppWithGenerate`); nachzulesen, ob `quiet_test.go` oder `verbosity_last_wins_test.go` (Test `TestExecute_Verbosity_LastWins_WarningFilter`, prüft `has no effect` als „nicht-gekennzeichnet“) die Aufteilung ebenfalls trägt |
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

- **Das Kennzeichen erscheint in `--json` oder in einem Golden-Fall** und ändert eine maschinenlesbare Ausgabe. — **Ausgang:** weiter offen: der Implementer liest `WarningEntry` im JSON-Umschlag und die Golden-Fälle vor der Änderung; trägt eine Ausgabe das Kennzeichen, geht der Slice nach `next` (§4).
- **Plan-Aussage über Code stimmt nicht** (Klasse [BEO-DOK/plan-aussage-ueber-code](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md), verkörpert in `AGENTS.md §3.10`; die Plan-Tabelle altert mit jedem Commit nach dem Abgleich): Der Abgleich der Tabelle gegen `git diff --stat` wird nach der letzten Fix-Runde wiederholt. — **Ausgang:** weiter offen: → [BEO-DOK/plan-aussage-ueber-code](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) im Register.

## 7. Closure-Notiz

Wird bei der Closure gefüllt.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `internal/hexagon/application` (APP, Code und Test), `internal/adapter/driving/cli` (CLI, nur Test) und, für Doku und Spezifikation, `docs/`, `spec/` (DOK) — alle Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006) (Abschnitt Modus-Deklaration), jede erfüllt die Schwelle ≥ 2 von 3 Achsen; keine Ausdifferenzierung nötig. Weil CLI berührt wird, gilt `AGENTS.md §3.9` (Kommentar-Selbstcheck) für den Diff, `AGENTS.md §3.10` für jede Plan-Aussage über Code.

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register ([`../observations/`](../observations/README.md)) durchgegangen — Treffer: [BEO-DOK/plan-aussage-ueber-code](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (DOK, verkörpert in `AGENTS.md §3.10`; Plan-Tabelle gegen den Diff nach der letzten Fix-Runde erneut abgleichen) und `BEO-CLI/kommentar-beschreibt-abwesendes` (CLI, verkörpert in `AGENTS.md §3.9`; der Doc-Kommentar von `sandboxRelaxationWarnings` ändert sich).

Alle berührten Sub-Areas GF.
