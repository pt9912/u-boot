# Verification-Report: [`slice-v1-devcontainer-claude-oauth-token`](../plan/planning/done/slice-v1-devcontainer-claude-oauth-token.md) — 2026-11-01

**Rolle:** Verifier (frischer Kontext, kein Self-Review, repariert nichts). Frage: Wurde der Slice richtig gebaut, gemessen an Plan, DoD, Entscheidung Weg A, [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten), [`SPEC-051`](../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) und den ausgeführten Sensoren? Form nach [`harness/verification.md`](../../harness/verification.md). Außer dieser Datei ist nichts geändert, nichts committet, nichts gepusht.

**Gegenstand:** Branch `worktree-agent-aa69746b53d530dc4`, Commits `fbf310dc`, `c08a9216`, `c2e68d7b`, `4abb1906` gegenüber der Basis `b7afa8ae` (20 Dateien, `git diff --stat b7afa8ae..HEAD`). Eingang (gelesen, nicht als Beweis genommen): der [Review-Report samt Nachlauf](2026-11-01-slice-v1-devcontainer-claude-oauth-token.md), das Messprotokoll [`README.md`](../messungen/slice-v1-devcontainer-claude-oauth-token/README.md) mit Rohdaten.

**Verdikt:** Der Code-, Spec- und Handbuch-Teil ist **erfüllt** (DoD-Punkte 2 und 3, `make gates` grün). Die DoD-Zeile **Messung ist zu Recht offen**: `claude -p` lief einmal statt zweimal, ein Lauf mit gültigem Token fehlt. Ob der Slice damit schließen kann, ist eine **Planner-Entscheidung** (Abschnitt 4); ich entscheide sie nicht. 0 HIGH, 0 MEDIUM, 1 LOW, 3 INFO. Die Closure-Aufgaben (Review-Häkchen, Closure-Notiz, Register, Risiko-Ausgänge, Paarungen) stehen offen; nichts im Stand verhindert sie.

## Verification Evidence

Scope:
- Slice: [`slice-v1-devcontainer-claude-oauth-token`](../plan/planning/done/slice-v1-devcontainer-claude-oauth-token.md) (liegt in `open/`)
- IDs: [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten), [`SPEC-051`](../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten), [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md)
- Artefakte: `internal/hexagon/application/devcontainer_sandbox.go` (`agentCredentialEnv`), `doctor_sandbox.go`, sechs Testdateien der Schicht, `internal/e2e/claude_oauth_token_docker_test.go`, `.devcontainer/devcontainer.json`, `spec/spezifikation.md`, Handbuch, Wartungsdoku, `CHANGELOG.md`, `harness/replay.md`, Messprotokoll

### 1. DoD-Abgleich (gegen gelesene Fundstellen)

- [ ] **Messung** — offen, ehrlich so geführt. Gelesen: der Docker-Test `TestE2E_ClaudeOAuthToken_ReachesTheSession` fährt je Zelle (Platzhalter gesetzt, ungesetzt, leer) eine Schleife `for run := 1; run <= 2` mit je einer Sitzung `devcontainer shell` und prüft die Zeilen `VALUE:<Platzhalter>` / `VALUE-UNSET` / `VALUE:`, dazu `KEY-UNSET` für den API-Schlüssel und dass der Platzhalter weder in `devcontainer.json` noch in `stderr` steht. Die Zusage „2 Läufe je Zelle“ im Protokoll deckt der Quelltext, die Rohdaten (`rohdaten-01.txt`) nicht: sie tragen nur die `claude -p`-Ausgabe und die Zeile `--- PASS … (68.34s)`; die Zellen stehen als Zusicherungen im grünen Test. `claude -p` mit Platzhalter: **ein** Aufruf im Test (keine Schleife), Ausgabe `Failed to authenticate. API Error: 401 OAuth access token is invalid.`, `CLAUDE-EXIT:1`. Die DoD verlangt „mindestens zwei Läufe je Zelle“, also fehlt ein Lauf in dieser Zelle; ein gültiges Token wurde nicht gefahren. Die Läufe einer Zelle gehören zu einem Containerstart (Protokoll benennt das). Datum und Umgebung stehen im Protokoll (2026-11-01, Linux `x86_64`, Docker Engine 29.9.2, Dev-Container-CLI 0.89.0). Keine Zugangsdaten in den Rohdaten (gelesen: nur Testzeilen und die Dienstantwort). **Nicht von mir nachgeprüft:** der Docker-Test selbst (kein `make test-docker`, Entscheidung des Projektinhabers); die Messzahlen sind die des Implementers.
- [x] **Der gewählte Weg (A)** — Evidence: `agentCredentialEnv` liefert bei aktivem `claude-code` `ANTHROPIC_API_KEY, CLAUDE_CODE_OAUTH_TOKEN` (danach `GH_TOKEN`, `GITLAB_TOKEN` unverändert); die Vorlage ist unverändert (`range .AgentEnv`, im Diff nicht enthalten); `.devcontainer/devcontainer.json` trägt die Zeile `"CLAUDE_CODE_OAUTH_TOKEN": "${localEnv:CLAUDE_CODE_OAUTH_TOKEN}"` an der Position der Liste. `doctor_sandbox.go`: im Diff nur der Hinweistext (Beispiel-Variablen); `missingAgentSources` und die Regel „alle fehlenden nennen“ unverändert, damit verlangt die Prüfung beide Referenzen, wie in der Entscheidung. Tests je Bedingung: `TestDoctor_SandboxCredentials_Agent` mit den Fällen „reference present“ (beide), „reference absent“, „devcontainer.json missing“, „other variable is no source“, „only the API key reference“ (Meldung nennt nur `CLAUDE_CODE_OAUTH_TOKEN`), „only the OAuth token reference“ (nennt nur `ANTHROPIC_API_KEY`); Erwartungswerte von `remoteEnv` in `generate_features_test.go`, `generate_forge_cli_test.go`, `generate_ollama_host_test.go`, `generate_ollama_test.go` im Diff nachgezogen. Mutationstests wurden nicht ausgeführt (kein Werkzeug im Gate); die Pins je Fall stehen im Diff.
- [x] **Spezifikation, Handbuch, `CHANGELOG.md`** — Evidence: [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) (Absatz Zugangsdaten: beide Referenzen, feste Reihenfolge; Prüfung Punkt 3: beide verlangt, nur die fehlende genannt) und [`SPEC-051`](../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) (Reihenfolge mit vier Namen) im Diff; der Satz „wird weder eingerichtet noch geprüft“ ist aus dem Absatz entfallen (Diff). Handbuch: `grep` nach `LH-`, `SPEC-`, `ADR-` im Handbuch liefert keinen Treffer; die Aussage zur Reichweite lautet „in der Umgebung der Sitzung von `u-boot devcontainer shell`“, nicht mehr „jedes Prozesses“; Widerruf des Tokens genannt; Ollama-Hinweis („setzen Sie dort auch `CLAUDE_CODE_OAUTH_TOKEN` nicht“) ohne Verhaltenszusage. Wartungsdoku nennt Grenzen („Ein Lauf mit einem gültigen Token ist nicht gefahren“, „1 Lauf“ für `claude -p`). READMEs: nicht im Diff, laut Plan „ohne Zugangsdaten“ (Plan-Aussage, von mir nicht erneut gelesen). `CHANGELOG.md`: Eintrag `[Unreleased]` / `Added` mit [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) und Regenerierungs-Hinweis. Beleg-Tabelle Bedingung → Test: Rendern beider Referenzen und Reihenfolge (die vier Generator-Testdateien), Doctor (s. o.), Weitergabe gesetzt/ungesetzt/leer (nur der Docker-Test, nicht im Gate).
- [x] **`make gates` grün** — selbst gefahren im Worktree (Abschnitt 3).
- [ ] Review durchgeführt — Report samt Nachlauf liegt vor (0 HIGH, 0 MEDIUM; die Nachlauf-Findings N-1 und N-2 sind im Stand behoben, siehe Abschnitt 2). Das Häkchen setzt der Planner.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag — offen (Planner); §7 ist Platzhalter.
- [ ] Beobachtungs-Register fortgeschrieben — offen (Planner).
- [ ] Jedes Risiko aus §6 mit Ausgang — offen: Risiko 1 bis 3 und 5 tragen „bei Closure“, Risiko 4 „weiter offen (Stand 2026-11-01)“.
- [ ] Drei Paarungen — offen (Planner).

### 2. Entscheidungs-Konformität (Weg A) und Plan gegen Code

| Entscheidung / Zusage | Befund |
|---|---|
| `doctor` verlangt beide Referenzen | konform; Code unverändert in der Regel, Test pinnt beide Einzelfälle; [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) Punkt 3 nennt es ausdrücklich. Folge für Bestandsprojekte (Warnung nach Upgrade, Hinweis: neu erzeugen) in Slice, Spec, Handbuch, CHANGELOG benannt. |
| Ollama-Ablauf unverändert ([`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md)) | konform; Bedingung bleibt nur für `ANTHROPIC_API_KEY`; Verhalten mit gesetztem Token dort ausdrücklich nicht gemessen. |
| Keine ADR, kein Lastenheft-CR | konform; `spec/lastenheft.md` und `docs/plan/adr/` nicht im Diff. |
| Reihenfolge `ANTHROPIC_API_KEY`, `CLAUDE_CODE_OAUTH_TOKEN`, `GH_TOKEN`, `GITLAB_TOKEN` | konform in Code-Kommentar, beiden Spec-Absätzen, `harness/replay.md` (Zeile `claude-code`), `.devcontainer/devcontainer.json`. |
| Ausschlüsse (`ANTHROPIC_AUTH_TOKEN`, `/login`, `gh`/`glab`/`GIT_TOKEN`) | eingehalten; im Diff keine Änderung daran. |
| Harte Regeln | kein `//nolint` im Diff, keine Gate-Lockerung, Architektur-Sicht und Lastenheft ohne Abwärtsverweis; Kommentare im Indikativ. |

`git diff --stat b7afa8ae..HEAD` gegen Plan §3 (Abgleichtabelle): Alle 20 Diff-Dateien haben eine Plan-Zeile (Quelle, Doctor, Tests, Docker-Test, Protokoll mit `docs/messungen/README.md`, `.devcontainer/devcontainer.json`, Spec, Handbuch, Wartungsdoku, `CHANGELOG.md`, `harness/replay.md`, Slice, Review-Report); „unverändert“-Zeilen (Vorlage, README, Lastenheft, `docs/plan/adr/`) stehen nicht im Diff. Die Nachlauf-Findings N-1 und N-2 sind im Stand behoben: die Zeile „`tools/messung/<name>/` entfällt“ steht in der Plan-Tabelle, und die Wendung „nicht in der Container-Konfiguration“ steht nicht mehr im Handbuch. `devcontainer_shell_test.go`: im Diff nur Kommentar und Fallname, Erwartungswert unverändert (der Plan führt es als „update (nur Kommentar und Fallname)“, stimmt).

### 3. Sensors

| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` (Worktree, Docker-Modus am Host) | pass | Exit 0; `baseline-verify: v6.18.1 OK — 54 Dateien`; `coverage-gate: OK — coverage 92.20% meets threshold 90%` (im Log steht die Zeile zweimal, einmal mit 92.10 %, beide über der Schwelle); `d-check: 952 Datei(en) geprüft, 0 Befund(e)`; `ablage-pruefung: OK` |
| `make docs-check` (mit diesem Report) | siehe Nachtrag am Ende | |
| `make test-docker` | **not run** | Entscheidung des Projektinhabers: nicht lokal. Der Docker-Test des Slice ist nur über das Messprotokoll des Implementers belegt, nicht von mir nachgeprüft. |
| `make verify-depguard` | not run | keine Importregel berührt (keine Port- oder Importänderung, `.golangci.yml` nicht im Diff) |

### 4. Zur Frage: kann der Slice mit der offenen Zeile „Messung“ schließen?

Ich entscheide das nicht; ich lege die Lage und die Folgen der Optionen dar.

Lage: Die Zusage in Handbuch und Wartungsdoku nennt nur das Gemessene und benennt die Lücke (`claude -p` einmal, kein gültiges Token). Insofern ist die Closure-Bedingung „die Zusage nennt nur das Gemessene“ (Slice §5) erfüllt. Verletzt ist allein der Buchstabe der DoD-Zeile („mindestens zwei Läufe je Zelle“, hier: die Zelle `claude -p` hat einen). Eine Wiederholung braucht `make test-docker` (26 Minuten, geteilte Engine) und ist laut Anweisung des Projektinhabers nicht lokal.

Optionen und was jede verlangt (Entscheidung beim Planner / Projektinhaber):

1. **Zweiter `claude -p`-Lauf**, an anderer Stelle gefahren (etwa CI oder Projektinhaber). Dann bleibt die DoD wörtlich und das Häkchen wird gesetzt; Risiko 4 bekommt einen Ausgang nach dem Lauf (der Teil „gültiges Token“ bleibt offen oder wird vom Projektinhaber gefahren).
2. **DoD-Zeile anpassen**, etwa „`claude -p` mit dem Platzhalter: ein Lauf; Variablen-Zellen je zwei“. Das ist eine Änderung des Plans, kein Gate nach [`AGENTS.md`](../../AGENTS.md) §3.6; sie braucht trotzdem einen Anker (Planner-Entscheidung in §7 oder Drift-Zeile der Roadmap), und die Zeile muss ehrlich benennen, was nicht gemessen ist. Der Wortlaut „Ein Lauf mit einem **gültigen** Token ist nur zugesagt, wenn der Projektinhaber ihn gefahren hat“ trägt diese Variante schon.
3. **Risiko 4 mit Ausgang „weiter offen“** (Stand im Slice bereits so geführt) und Beobachtung im Register ([`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md)) fortschreiben. Der Slice-Text lässt „weiter offen“ als Ausgang zu (DoD-Zeile *Jedes Risiko … Ausgang (eingetreten / entfallen / weiter offen)*). Ob §6 („kein Slice geht nach `done/`, während eines ohne Ausgang dasteht“) damit erfüllt ist, hängt daran, ob „weiter offen“ als Ausgang zählt; der Slice führt es so, die Baseline-Stelle habe ich nicht gesondert gelesen.

Unabhängig davon verlangen DoD und §6 Ausgänge für Risiko 1 (Entscheidung „beide“, belegt durch Spec und Tests), 2 (Handbuch nennt Widerruf; die „Reichweite des Tokens“ und das privilegierte Sandbox-Profil nennt das Handbuch in den gelesenen Zeilen nicht ausdrücklich, siehe F-1), 3 (nicht gemessen, Handbuch rät ab) und 5 (Golden Cases nachgezogen, `harness/replay.md`).

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Slice §6 Risiko 2 verlangt, die Dokumentation nenne „Widerruf und Reichweite des Tokens“. Im Diff des Handbuchs steht der Widerruf („widerrufen Sie es bei Anthropic“), nicht die Reichweite des Tokens und nicht, dass das Sandbox-Profil mit verschachteltem Docker privilegiert läuft. | Slice §6 Risiko 2 | `docs/user/benutzerhandbuch.md` · „Nutzen Sie nur ein Token, das Sie widerrufen können“ | nein — Review-only; der Ausgang von Risiko 2 hängt daran | Risiko-Ausgang ohne Dokumentations-Entsprechung |
| F-2 | INFO | Die Rohdaten tragen nur die `claude -p`-Zeilen; die Zellen „gesetzt/ungesetzt/leer“ mit je zwei Läufen sind nur durch den grünen Test und die Schleifen im Quelltext belegt, nicht durch ausgegebene Rohzeilen. Zudem schlugen im selben Paketlauf drei fremde Tests fehl (im Protokoll benannt); ein sauberer Lauf nur dieses Tests liegt nicht vor. | `harness/verification.md` (Evidence) | Messprotokoll, README · „Die drei Variablen-Zellen stehen im Test als Zusicherungen“ | nein — `make test-docker` (nicht gefahren) | Messbeleg ohne Rohzeilen je Zelle |
| F-3 | INFO | Der Slice liegt in `open/`, nicht in `in-progress/`; im Worktree steht dort [`slice-v1-devcontainer-instanz-label`](../plan/planning/done/slice-v1-devcontainer-instanz-label.md). Der Lifecycle-Zustand ist Sache des Planners (Start-Trigger in Slice §4); ich habe nichts bewegt. | Slice §4 (Start-Trigger) | [Slice-Datei in `open/`](../plan/planning/done/slice-v1-devcontainer-claude-oauth-token.md) | nein — `make docs-check` meldet es nicht | — |
| F-4 | INFO | Mutationstests (DoD „Tests je Bedingung und je Mutation“) sind nicht ausgeführt; belegt sind Fälle je Bedingung im Diff, kein Mutationslauf. | Slice DoD Punkt 2 | `internal/hexagon/application/doctor_sandbox_test.go` · „only the OAuth token reference“ | nein | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Kennungen in Nutzer-Doku (`AGENTS.md` §5 Regel 6) | geprüft, ohne Befund: keine `LH-`/`SPEC-`/`ADR-`-Kennung im Handbuch |
| Klartext-Token in Dateien (Rohdaten, Protokoll, `devcontainer.json`, Test) | geprüft, ohne Befund: nur erfundener Platzhalter im Test, nicht in Rohdaten; `devcontainer.json` trägt nur die Referenz |
| Reste alter Fassung (Singular „der Schlüssel des Agenten“, „beide Variablen“) | geprüft, ohne Befund; Nachlauf bestätigt |
| Lastenheft, ADR, Architektur-Sicht, Import-/Port-Grenzen | geprüft, ohne Befund: nicht im Diff |
| Carveouts | unverändert; kein Eintrag neu oder gelöst |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Risiko-Ausgang ohne Dokumentations-Entsprechung · Messbeleg ohne Rohzeilen je Zelle

## Verdikt

**Bereit für die Closure:** bedingt. Code, Spec, Handbuch und `make gates` tragen den Slice. Offen ist (a) die Entscheidung des Planners zur DoD-Zeile „Messung“ (Abschnitt 4: zweiter Lauf, DoD-Anpassung oder Risiko 4 „weiter offen“), (b) F-1 (Reichweite des Tokens und privilegiertes Profil im Handbuch, oder Risiko 2 mit benanntem Rest), (c) die Closure-Aufgaben. **Merge-blockierend:** nein für den Code; die DoD ist erst nach (a) vollständig.

**Übergabe:** an den Planner (DoD-Entscheidung, Risiko-Ausgänge, Closure-Notiz, Register); F-1 optional an den Implementer. Nicht ausgeführt: `make test-docker` (Anweisung des Projektinhabers), `make verify-depguard` (nicht berührt).

**Nachtrag:** `make docs-check` mit diesem Report: siehe Bericht des Aufrufs (Ergebnis 0 Befunde nach der Berichtigung der Kennungs-Links).
