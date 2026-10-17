# Verification Evidence: slice-v2-spec-devcontainer-sshd-claude-code — 2026-10-17 <!-- d-check:ignore (Slice-Kennung im Titel) -->

**Rolle:** Verifier (frischer Kontext, kein Self-Review). Frage: Wurde der Slice richtig gebaut, gemessen an Plan, DoD und Sensors? Form nach `harness/verification.md`. <!-- d-check:ignore (Pfad als Zitat) -->

**Gegenstand:** Slice [`slice-v2-spec-devcontainer-sshd-claude-code`](../plan/planning/in-progress/slice-v2-spec-devcontainer-sshd-claude-code.md) (Welle [`welle-devcontainer-sshd-claude-code`](../plan/planning/welle-devcontainer-sshd-claude-code.md), letzter Slice), Commits `f266bd4..e53f38a` (`e32859a`, `38df18d`, `ef43e46`, `e53f38a`): `spec/spezifikation.md`, Roadmap-Ruhe-Marker, Plan-Tabelle. Review-Report: [`2026-10-17-slice-v2-spec-devcontainer-sshd-claude-code.md`](2026-10-17-slice-v2-spec-devcontainer-sshd-claude-code.md).

## Verification Evidence

Scope:
- Slice: [`slice-v2-spec-devcontainer-sshd-claude-code`](../plan/planning/in-progress/slice-v2-spec-devcontainer-sshd-claude-code.md)
- IDs: [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`ADR-0020`](../plan/adr/0020-drittanbieter-features-im-katalog.md)
- Spec-Stellen: [`SPEC-017`](../../spec/spezifikation.md#spec-017--diagnose-codes-der-doctor-prüfungen-code-registry), [`SPEC-019`](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), [`SPEC-023`](../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen), [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten), §7 Historie
- Diff-Umfang: `spec/spezifikation.md`, `docs/plan/planning/in-progress/roadmap.md`, der Slice (kein Code, kein Test, kein Handbuch)

Vorgehen: Jede Spec-Aussage wurde unabhängig am Code nachgelesen (`devcontainer_features.go` `featureCatalogue`, `thirdPartyFeatureWarnings`; `devcontainer_egress.go` `egressFeatureHosts`, `egressHosts`; `devcontainer_sandbox.go` `sandboxRelaxationWarnings`, `catalogueFeatureActive`, `agentCredentialEnv`, `agentNodeWarnings`, `devcontainerTemplateData`; `doctor_sandbox.go` `credentialPatterns`, `credentialScanFiles`, `checkSandboxCredentials`, `checkSandboxEgress`, `checkSandboxRuntime`, `hasTokenSource`, `missingAgentSources`; `generate.go` `generateDevcontainer` (einziger Aufrufer der drei Warn-Funktionen); `templates/devcontainer/devcontainer.json.tmpl`; `cli/warnings.go` `visibleWarnings`; `cli/doctor.go` `writeDoctorReport`; `jsontestutil/coderegistry.go`). Alle Hosts, Warncodes, Kennzeichen, Dateien und Muster wurden einzeln verglichen, nicht stichprobenhaft.

DoD-Abgleich (§2 des Slice):
- [x] **Spezifikation** — Evidence (je Teilaussage gegen den Code):
  - `SPEC-023`: `sshd` → `sshd`, Standardversion `1`; `claude-code` → `ghcr.io/anthropics/devcontainer-features/claude-code`, Version `1`, Herausgeber `anthropics`; `homebrew` Version `2`, `meaningful-ooo`; alle Schlüssel ohne `featureSources.allow` aktivierbar; Abhängigkeit `node` (`featureCatalogue`, `agentNodeWarnings`). Stimmt. <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->
  - `SPEC-019`: Zeile `claude-code` = genau `registry.npmjs.org`, `api.anthropic.com`; `sshd` ohne Host; `registry.npmjs.org` gehört auch zu `node` (`egressFeatureHosts`). Stimmt. <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->
  - `SPEC-017`: Beschreibung nennt Agenten-Klartext, Token-Quelle bei https-Clone, Referenz `${localEnv:<NAME>}` (`checkSandboxCredentials`); Schlüssel in `coderegistry.go` unverändert, kein neuer Code. Stimmt. <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->
  - `SPEC-044` Warnungstabelle: `sshd` nur mit Sandbox-Profil, Code `LH-FA-DEV-006`, gekennzeichnet; `homebrew`/`claude-code` ohne eigene `source`, `LH-NFA-SEC-004`, ungekennzeichnet, beide Profile; `node`-Hinweis `LH-FA-DEV-003`, ungekennzeichnet, beide Profile; Level `warn`, Subjekt `.devcontainer/devcontainer.json`; Erzeuger nur `generate devcontainer` (`generate.go:781-783`; `initproject.go` übergibt `nil`). Stimmt. <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->
  - `SPEC-044` Sandbox-Regel und Zugangsdaten: Template schreibt nur `features`-Schlüssel; `runArgs` allein aus `nestedRuntime`/Egress; `remoteEnv` nur im Sandbox-Zweig, `ANTHROPIC_API_KEY` als Referenz bei aktivem `claude-code` (`agentCredentialEnv`). Stimmt. <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->
  - `SPEC-044` Prüfung: Dateiliste (7), alle sieben Muster mit Mindestlängen, Reihenfolge Klartext → https-Clone → Agenten-Referenz, „höchstens ein Befund“, Platzhalter-Grenze, `readonly`-Erkennung ohne `ro`, nicht lesbare Datei = ohne Referenz, Kennzeichen an allen drei Warn-Zweigen. Stimmt. <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->
  - `SPEC-043`: `doctor`-Zeile und Ausnahme 1 (Filter `item.SecurityRelaxation` in `writeDoctorReport`, `visibleWarnings`); `nestedRuntime: docker` in `doctor` Stufe OK (`checkSandboxRuntime`); Egress „has no effect“ gekennzeichnet in `doctor`, ungekennzeichnet in `generate devcontainer`; `nestedRuntime` „has no effect“ ungekennzeichnet in beiden; `features.allowlist`/`features.drift` ungekennzeichnet (kein Flag in `doctor.go`). Stimmt. <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->
  - §7: Zeile vom 2026-10-17 vorhanden.
  - Einschränkung: Der Slice-Text selbst trägt noch veraltete Zustandsaussagen (siehe Diskrepanz D-1); das betrifft den Plan, nicht die Spec.
- [x] **Gegen den Code gelesen** — Evidence: Spec-Commits `e32859a`, `ef43e46` liegen nach den Moves der vier anderen Slices nach `done/` (alle vier Vorgänger stehen in `done/`, siehe Slice-Kopf Zeile 12). Fundstellen führt diese Datei je Aussage; die Closure-Notiz muss sie übernehmen (siehe Auflage A-2).
- [x] **Traceability** — Evidence: `make doc-trace` selbst gefahren, `83 Anforderung(en), 0 Waise(n)`; `make docs-check` (Teil von `make gates`) `302 Datei(en) geprüft, 0 Befund(e)`; jede neue `SPEC-044`-Referenz löst auf (Überschrift `SPEC-044 — Features sshd und claude-code im Devcontainer: Warnungen und Zugangsdaten`). <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->
- [x] `make gates` grün — selbst gefahren, Exit 0 (Belege unten).
- [x] Review durchgeführt, Report liegt vor — Evidence: Report `2026-10-17-slice-v2-spec-devcontainer-sshd-claude-code.md` (0 HIGH, 1 MEDIUM, 4 LOW, 3 INFO). Er ist im Arbeitsbaum noch **nicht committet** (Auflage A-5). <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->
- [x] Doku-Update, falls öffentlicher Vertrag berührt — Evidence: kein öffentlicher Vertrag geändert; die Spec ist das Ziel. Ein Rest-Hinweis zum Handbuch steht als Hinweis H-1.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag — offen: §7 steht auf „Wird bei der Closure gefüllt.“ Gehört in die Closure.
- [ ] Beobachtungs-Register fortgeschrieben oder in §7 notiert — offen, gehört in die Closure.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang — offen: die drei Ausgänge stehen im Slice (entfallen, weiter offen, entfallen), der Punkt „weiter offen“ trägt kein Ziel; F-1 ist als neues Risiko aufzunehmen (Auflage A-3).
- [ ] Drei Paarungen (Anker · Folge-Slice · Register) — offen: von der Welle-Closure zu tragen; Folge-Slice für F-1 ist noch nicht angelegt (Auflage A-3).

Review-Findings am Text (nicht an der Prosa des Berichts) geprüft:
| Finding | Ergebnis | Beleg im Text |
| --- | --- | --- |
| F-1 (Egress „has no effect“ unter `--quiet` in `generate devcontainer`) | bewusst offen (Entscheidung des Projektinhabers 2026-10-17: Code-Folge-Slice) | Spec benennt den Ist-Zustand klar: `SPEC-043` „Ungekennzeichnet sind die Meldungen „has no effect“ der Warnungsausgabe von `generate devcontainer` (für `nestedRuntime` wie für die Egress-Restriktion) … Gekennzeichnet, also sichtbar, ist dagegen der `doctor`-Eintrag“; Code bestätigt (`warn("LH-FA-DEV-008", …)` ohne Flag in `sandboxRelaxationWarnings`, `SecurityRelaxation: true` in `checkSandboxEgress`). | <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->
| F-2 (aktiv über aufgelöste Quelle) | behoben | `SPEC-044` „verglichen wird die aufgelöste Quelle über alle aktivierten Einträge, nicht der Schlüssel: Auch ein Eintrag mit anderem Schlüssel und `source` gleich der Katalogquelle ist aktiv“; entspricht `catalogueFeatureActive` (`source == want`). | <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->
| F-3 (`ro` zählt nicht) | behoben | `SPEC-044` Prüfschritt 2 „… `/run/secrets` und `readonly` zusammen (Secret-Mount; die Kurzform `ro` der Mount-Option zählt nicht)“; entspricht `hasTokenSource`. | <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->
| F-4 (Erzeuger-Aufzählung) | behoben | `SPEC-043` Ausnahme 1 „in der Warnungsausgabe von `generate devcontainer` (nur dieser Befehl erzeugt sie)“; entspricht `generate.go:781-783`. | <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->
| F-5 (Plan-Tabelle) | teilweise behoben | Plan-Tabelle §3 gegen `git diff --stat f266bd4..HEAD` abgeglichen (drei Dateien, jede mit Eintrag; Spec-Zeile nennt Absatz und drei Unterpunkte; Roadmap-Zeile da). Die veralteten Marker „SPEC-044 ist noch nicht vergeben“ stehen aber weiter in Kopf Zeile 16, §1 Zeile 27 und DoD 1 Zeile 46 (D-1). | <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | selbst gefahren auf HEAD `e53f38a`, Exit 0 (`gates-exit=0`); `coverage-gate: OK — coverage 91.20% meets threshold 90%`; `d-check: 302 Datei(en) geprüft, 0 Befund(e)`; `baseline-verify: v6.18.1 OK — 54 Dateien`; lint und test grün |
| `make doc-trace` | pass | selbst gefahren: `83 Anforderung(en), 0 Waise(n)` |
| `make docs-check` | pass | Teil von `make gates` (0 Befunde); nach Anlage dieser Datei erneut gefahren (siehe Übergabe) |

Traceability:
| ID | Beleg |
| --- | --- |
| [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) | `SPEC-023`, `SPEC-044` (`node`-Hinweis); Code `featureCatalogue`, `agentNodeWarnings` | <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->
| [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) | `SPEC-044` (Login-Dienst); Code `sandboxRelaxationWarnings` | <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->
| [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) | `SPEC-019`, `SPEC-043`; Code `egressFeatureHosts`, `checkSandboxEgress` | <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->
| [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) | `SPEC-017`, `SPEC-044`; Code `checkSandboxCredentials`, Template `remoteEnv` | <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->
| [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | `SPEC-043`; Code `visibleWarnings`, `writeDoctorReport` | <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->
| [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) | `SPEC-044`; Code `thirdPartyFeatureWarnings` | <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->
| [`ADR-0020`](../plan/adr/0020-drittanbieter-features-im-katalog.md) | `SPEC-023` (Drittanbieter-Einstufung) | <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->

Hinweis zur Bindung: Kein Gate prüft die Spec-Tabellen gegen den Code (nur die `doctor`-Code-Registry hat einen Drift-Test). Der Beleg dieses Slice ist die hier dokumentierte Lesung, keine maschinelle Prüfung.

Carveouts:
- Neu: none
- Gelöst: none
- Unverändert: none (F-1 ist kein Carveout, sondern ein Folge-Slice; siehe Auflage A-3)

Nicht ausgeführt:
- `make test-docker` — kein Docker-Pfad berührt (reine Spec-/Plan-Änderung).
- `make verify-depguard` — keine Importregel berührt.
- `make ci` / `make fullbuild` — kein Security- oder Release-Pfad berührt.
- `make doc-immutable` — kein Accepted-ADR berührt (der Diff enthält keine ADR).

Commit / Artefakt:
- Spec-Änderung `e32859a`, `ef43e46`; Plan `38df18d`, `e53f38a`; geprüfter Stand HEAD `e53f38a`.

## Diskrepanzen

**D-1 (Plan, veraltete Zustandsmarker).** Der Slice trägt den Zustand „`SPEC-044` ist noch nicht vergeben“ an drei Stellen, obwohl der Eintrag vergeben ist und die Plan-Tabelle ihn schon verlinkt (Verstoß gegen `AGENTS.md` §3.7 Zustandsfelder und §3.10). Wortlaut-Vorschlag, je Stelle, dabei die `d-check:ignore`-Kommentare entfernen: <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->
- Kopf Zeile 16: „… neuer Eintrag `SPEC-044` und §7 (Historie).“ <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->
- §1 Zeile 27: „… und der Eintrag `SPEC-044` für die Sandbox-Regel beider Features, …“ <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->
- DoD 1 Zeile 46: „der Eintrag `SPEC-044` beschreibt Sandbox-Regel, Zugangsdaten-Übergabe und Warnungen; …“ <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->
Danach sind die drei `d-check:ignore (SPEC-044 ist noch nicht vergeben)` entbehrlich; `make docs-check` danach erneut fahren. <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->

**D-2 (Plan, Aussage zu den Findings).** DoD 1 verlangt, dass `SPEC-043` Egress „has no effect“ als sichtbar festhält; das ist erfüllt, nur bezogen auf `doctor`. In `generate devcontainer` entfällt dieselbe Meldung unter `--quiet` (F-1). Der DoD-Satz „Egress ‚has no effect‘ sichtbar bleibt, `nestedRuntime` ‚has no effect‘ dagegen entfällt“ ist ohne den Zusatz „(`doctor`)“ mehrdeutig; Vorschlag: „… dass in `doctor` Egress ‚has no effect‘ sichtbar bleibt, `nestedRuntime` ‚has no effect‘ dagegen entfällt, und dass die Warnungsausgabe von `generate devcontainer` beide Meldungen ungekennzeichnet führt (Ist-Zustand)“. Die Plan-Aussage in §1 Zeile 29 („Egress ‚has no effect‘ bleibt unter `--quiet` sichtbar“) trägt dieselbe Mehrdeutigkeit. <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->

**Hinweis H-1 (außerhalb des Slice).** `docs/user/benutzerhandbuch.md` (Zeilen 299 und 1417) nennt als Erzeuger der übrigen Lockerungen „je nach Befehl `init`, `generate`, `up`“; Spec und Code kennen nur `generate devcontainer`. Der Handbuch-Satz liegt nicht im Diff und gehört in die Handbuch-Pflege (oder den Code-Folge-Slice zu F-1, der die Aufzählung ohnehin berührt).

**Hinweis H-2 (Spec, kein Befund).** Die Spec nennt den Ist-Zustand klar; sie weist die Egress-Asymmetrie nicht als vorläufig aus. Das ist richtig: Die Spec darf keine Slices nennen (`AGENTS.md` §3.4, [`ADR-0013`](../plan/adr/0013-dokumentationsreferenzmodell.md)). Das Nachziehen von `SPEC-043`/`SPEC-044` gehört in den Code-Folge-Slice. <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->

## Vorschlag: Ausweis von F-1 in der Closure

**§6 (Risiken):** neuer Punkt mit Ausgang „weiter offen“:
„**Egress ‚has no effect‘ entfällt unter `--quiet` in `generate devcontainer`** (Review F-1): Der Nutzer, der `egress.enabled: true` ohne Sandbox-Profil mit `--quiet` erzeugt, erhält keine Meldung, `doctor` weist sie aus. Die Spezifikation beschreibt den Ist-Zustand. — **Ausgang:** weiter offen: Entscheidung des Projektinhabers 2026-10-17, im Code-Folge-Slice `slice-v2-quiet-egress-ohne-sandbox-kennzeichen` (Arbeitstitel; Kennung beim Anlegen zu vergeben) zu schließen: `SecurityRelaxation` an der Meldung in `sandboxRelaxationWarnings`; derselbe Slice zieht `SPEC-043` (Einstufungsliste, Ausnahme 1) und `SPEC-044` nach.“ <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->

**§7 (Closure-Notiz):** in „Risiko-Ausgänge“ dieselbe Zeile (eingetreten: nein; weiter offen: ja, mit Folge-Slice-Kennung); in „Steering-Loop-Einträge“ die Finding-Klassen des Reviews (Spec legt Code-Asymmetrie fest, Lastenheft mehrdeutig · Spec-Aussage über Code ungenau · Plan-Tabelle gegen Diff nicht abgeglichen · Lastenheft-Wortlaut weiter oder enger als Verhalten) mit Zählung; ohne `liegt in` (gezählt, nicht verkörpert), solange keine Schwelle erreicht ist. Der Folge-Slice liegt vor dem Move nach `done/` in `open/` oder `next/` (Folge-Slice-Paarung), und die Welle-Datei nennt ihn unter ihren Ausgängen.

## Übergabe

Verdikt: **DoD-Zeilen 1 bis 6 und 7 erfüllt, Zeilen 8 bis 11 (Closure-Notiz, Register, Risiko-Ausgänge, Paarungen) offen** — gehören in die Closure, nicht in die Verifikation. Die Spec stimmt in allen geprüften Aussagen mit dem Code überein; F-2 bis F-4 sind am Text behoben, F-5 nur teilweise (D-1). F-1 bleibt bewusst offen.

Auflagen für die Closure:
- A-1: D-1 beheben (drei Marker, `d-check:ignore` entfernen), D-2 klarstellen.
- A-2: Closure-Notiz §7 mit den Fundstellen aus dieser Datei (Datei und Symbol je Aussage) füllen.
- A-3: Folge-Slice zu F-1 anlegen (Code, `SecurityRelaxation` in `generate devcontainer`; zieht `SPEC-043`/`SPEC-044` nach), in §6 und §7 ausweisen (Vorschlag oben). <!-- d-check:ignore (Spec- und Slice-Kennungen als Zitat) -->
- A-4: Beobachtungs-Register fortschreiben oder in §7 notieren, dass keine Beobachtung anfiel; Klasse „Plan-Aussage über Code“ trat hier erneut auf (D-1, Review F-5), Zählung im Register prüfen (`BEO-DOK/plan-aussage-ueber-code`).
- A-5: Review-Report und diese Datei sind nicht committet (`git status`: untracked); mit der Closure committen, `make docs-check` danach erneut fahren.
- A-6: Erst Inhalt (DoD-Häkchen, Closure-Notiz), dann der reine `git mv` nach `done/` (`AGENTS.md` §3.3); die Welle ist danach schließbar.
