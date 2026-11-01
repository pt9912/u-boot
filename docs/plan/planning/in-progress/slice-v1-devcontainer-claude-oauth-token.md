# Slice V1: `claude-code` im Sandbox-Devcontainer — das OAuth-Token des Abos (`CLAUDE_CODE_OAUTH_TOKEN`) als Referenz durchreichen, Wege A bis C zur Entscheidung des Projektinhabers

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** Release `v0.17.0` (geplant, Roadmap *Nächste Wellen*; die Welle-Datei entsteht bei der Eröffnung). Die Closure-Bedingung der Welle (alle Slices in `done/`, `make gates` und Replay grün, Release-Cut mit Tag) ist von der DoD dieses Slice verschieden. Anlass: Der Projektinhaber meldet sich bei Claude Code mit einem OAuth-Token an (2026-10-31); `u-boot` reicht dieses Token nicht in den Devcontainer durch.

**Bezug:** [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zugangsdaten eines Agenten nur zur Laufzeit, nur als Referenz), [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Katalog-Feature `claude-code`).

**Berührte Spec-Stellen:** [`SPEC-044`](../../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) (Absatz *Zugangsdaten des Agenten und der Forge-Werkzeuge*, Zeile 608; Prüfung `devcontainer.sandbox.credentials`, Punkt 3, Zeile 624), je nach Weg [`SPEC-051`](../../../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) (Absatz *Zugangsdaten*, Zeile 769; nicht gelesen — prüfen).

**Verantwortlich:** pt9912

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-31.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Wer Claude Code mit einem OAuth-Token des Abos nutzt, bekommt es im Sandbox-Devcontainer mit aktivem `claude-code` zur Laufzeit als Umgebungsvariable des Hosts, ohne dass ein Wert in einer erzeugten Datei steht — oder `u-boot` und das Handbuch sagen, was gemessen und was nicht zugesagt ist.

**Befund** (gelesen am 2026-10-31, Fundstellen je Punkt; [`AGENTS.md`](../../../../AGENTS.md) §3.10):

- `agentCredentialEnv` (`internal/hexagon/application/devcontainer_sandbox.go`, Zeilen 446 bis 458): liefert für ein aktives `claude-code` nur `ANTHROPIC_API_KEY`, daneben `GH_TOKEN` und `GITLAB_TOKEN`. Der Aufruf `data.AgentEnv = agentCredentialEnv(dc)` steht in Zeile 101; die Vorlage `templates/devcontainer/devcontainer.json.tmpl` rendert die Liste in `remoteEnv` (Zeilen 33 bis 35, `range .AgentEnv`).
- `checkSandboxCredentials` (`doctor_sandbox.go`, ab Zeile 134) liest dieselbe Liste (Zeile 147) und meldet eine fehlende Referenz; `credentialPatterns` (Zeile 54) erkennt `CLAUDE_CODE_OAUTH_TOKEN=<Wert>` schon als Klartext in einer Projektdatei.
- [`SPEC-044`](../../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) (Zeile 608) sagt ausdrücklich: andere Variablen reicht der Generator nicht durch; ein Anmeldeweg über `CLAUDE_CODE_OAUTH_TOKEN` „wird weder eingerichtet noch geprüft“.
- Handbuch (`docs/user/benutzerhandbuch.md`, Zeilen 1914 bis 1930): `u-boot` reicht nur `ANTHROPIC_API_KEY` durch; `CLAUDE_CODE_OAUTH_TOKEN` gehöre nicht zur Zusage und sei nicht gemessen.
- Beobachtet am 2026-10-31 (Linux `amd64`, neu angelegter Container, `u-boot devcontainer shell`): Im Container ist keine der Variablen `ANTHROPIC_*`, `CLAUDE_*`, `GH_*`, `GIT_*` gesetzt, weil sie auf dem Host nicht gesetzt waren; `devcontainer.json` verweist nur auf `GIT_TOKEN`, `ANTHROPIC_API_KEY` und `GH_TOKEN`.

**Wege** (Entscheidung des Projektinhabers vor dem Start; der Slice wählt nicht):

- **A, `CLAUDE_CODE_OAUTH_TOKEN` zusätzlich:** `agentCredentialEnv` liefert für `claude-code` beide Variablen; `remoteEnv` trägt beide Referenzen. Offen: Ob `u-boot doctor` beide verlangt oder mindestens eine (zwei gleichwertige Anmeldewege); die heutige Regel nennt bei mehreren fehlenden Referenzen alle.
- **B, eine allgemeine Liste durchgereichter Variablen in `u-boot.yaml`** (Arbeitsname `passEnv`): neuer Schlüssel in Schema, `config set`, Rendern und Drift-Prüfung. Deckt auch andere Fälle, ist aber vermutlich ein Lastenheft-CR ([`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) nennt Umgebungsvariablen des Agenten nicht einzeln; auf Optionen und Listen nicht durchgesehen — prüfen).
- **C, nichts im Produkt:** Das Handbuch beschreibt den Weg `read -rs CLAUDE_CODE_OAUTH_TOKEN && export CLAUDE_CODE_OAUTH_TOKEN` in der geöffneten Shell, mit Datum und Messstand. Das Problem bleibt für jeden, der das Token beim Start mitgeben will.

### Entscheidung (Weg A, Projektinhaber 2026-11-01; Teilfragen: Implementer)

Weg A ist vom Projektinhaber entschieden. Die Teilfragen entscheidet der Implementer mit der kleinsten Lösung, die [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`SPEC-044`](../../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) und [ADR-0036](../../adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) trägt; kein Konflikt mit Lastenheft oder ADR (Lastenheft nennt „die Umgebungsvariable des Agenten“ ohne Zahl, [ADR-0036](../../adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) setzt nur die Bedingung „ohne gesetzten `ANTHROPIC_API_KEY` im Aufruf“ und bleibt unberührt).

- **`doctor` verlangt beide Referenzen.** `checkSandboxCredentials` prüft, ob `devcontainer.json` die Referenz enthält (`missingAgentSources`), nicht die Umgebung des Hosts. Die Regel „mindestens eine“ hätte eine neue Verzweigung gebraucht und eine ältere Datei (nur `ANTHROPIC_API_KEY`) stillschweigend durchgelassen, obwohl das Token dann nie im Container ankommt. Mit „beide“ ist die heutige Regel unverändert (Meldung nennt alle fehlenden); die Warnung bei einer älteren Datei nennt nur die fehlende Variable, der Hinweis führt zu `u-boot generate devcontainer`. Eine Host-Umgebung mit nur einem Weg erzeugt keine Warnung, weil die Prüfung keine Host-Variable liest. Pin: `TestDoctor_SandboxCredentials_Agent` (Fälle „only the API key reference“, „only the OAuth token reference“).
- **Ollama-Ablauf unverändert.** [ADR-0036](../../adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) bleibt wörtlich gültig (Bedingung nur für `ANTHROPIC_API_KEY`). Das Verhalten mit gesetztem `CLAUDE_CODE_OAUTH_TOKEN` im Ollama-Ablauf ist nicht gemessen; das Handbuch rät, die Variable dort nicht zu setzen, und sagt dazu nichts Weiteres (Risiko 3).
- **Keine ADR, kein Lastenheft-CR.** Der Weg erweitert eine Liste in einer vorhandenen Funktion; Zusage und Rendern sind in [`SPEC-044`](../../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) geregelt.
- **Reihenfolge der Referenzen:** `ANTHROPIC_API_KEY`, `CLAUDE_CODE_OAUTH_TOKEN`, `GH_TOKEN`, `GITLAB_TOKEN` (feste Renderfolge, in [`SPEC-044`](../../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) und [`SPEC-051`](../../../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) nachgezogen).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **`ANTHROPIC_AUTH_TOKEN`.** Das Handbuch setzt ihn im Aufruf selbst (`ANTHROPIC_AUTH_TOKEN=<Platzhalter> claude …` im Ollama-Ablauf); eine Referenz in `remoteEnv` würde diesen Ablauf stören. Kein Bedarf des Projektinhabers.
- **Anmeldung per `/login` im Container.** Nicht gemessen, die Hosts der Anmeldung stehen nicht in der Default-Allowlist, und das Home von `vscode` liegt nicht im Workspace-Volume. Ein anderer Vorgang.
- **Der Ollama-Ablauf mit Claude Code** ([ADR-0036](../../adr/0036-cloud-modelle-von-ollama-im-devcontainer.md)): Er gilt nur bei ungesetztem `ANTHROPIC_API_KEY`; wie sich ein gesetztes `CLAUDE_CODE_OAUTH_TOKEN` dort verhält, ist nicht gemessen und bleibt außerhalb der Zusage (Risiko 3).
- **`gh`- und `glab`-Token, `GIT_TOKEN`.** Unverändert; sie sind in [`SPEC-051`](../../../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) und [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) geregelt.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **Messung** ([`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)), angepasst auf das Gemessene (Entscheidung des Projektinhabers 2026-11-01, Anker §7 *DoD-Anpassung*): mit einem Platzhalter-Wert kommt die Variable über `u-boot devcontainer shell` im Container an; ungesetzt auf dem Host kommt sie ungesetzt an, leer gesetzt als leere, je Zelle zweimal im Docker-Test; `claude -p` mit dem Platzhalter lief **einmal** und endete mit der Antwort des Dienstes `401 OAuth access token is invalid`. Ein Lauf mit einem **gültigen** Token und ein zweiter `claude -p`-Lauf sind nicht gefahren und nicht zugesagt (Risiko 4, weiter offen). *Beleg:* Messprotokoll [`docs/messungen/slice-v1-devcontainer-claude-oauth-token/README.md`](../../../messungen/slice-v1-devcontainer-claude-oauth-token/README.md) mit Rohdaten (Abschnitt Grenzen); ein Skript unter `tools/messung/<name>/` entfällt, die Messung ist der Docker-Test `internal/e2e/claude_oauth_token_docker_test.go` ([`harness/conventions.md`](../../../../harness/conventions.md#ablage-von-messprotokollen-und-skripten)); keine Zugangsdaten in den Rohdaten.
- [x] **Der gewählte Weg** (A, B oder C): bei A oder B reicht `u-boot generate devcontainer` die Variable als Referenz durch, `u-boot doctor` prüft sie nach der festgelegten Regel, Tests je Bedingung und je Mutation; bei C steht der Ablauf mit Datum und Messstand im Handbuch.
- [x] **Spezifikation, Handbuch, `CHANGELOG.md` nach dem Code** mit Beleg-Tabelle (Bedingung → Test oder Streichung); Handbuch und READMEs ohne Kennungen ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6); der Satz „wird weder eingerichtet noch geprüft“ in [`SPEC-044`](../../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) entfällt oder bleibt, je nach Weg.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register ([`observations/`](../observations/)) fortgeschrieben, oder in §7 „keine Beobachtung angefallen“.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

Gelesen am 2026-10-31; Zeilennummern verschieben sich, der Implementer liest sie am Ist-Stand nach ([`AGENTS.md`](../../../../AGENTS.md) §3.10). Welche Zeilen der Weg berührt, entscheidet der Projektinhaber.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/messungen/<kennung>/` (`tools/messung/<name>/` entfällt) | neu | Protokoll der Messung (Platzhalter; ein gültiges Token nur durch den Projektinhaber); die Messung ist der Docker-Test `internal/e2e/claude_oauth_token_docker_test.go`, kein Skript |
| `internal/hexagon/application/devcontainer_sandbox.go` | update (Weg A) | `agentCredentialEnv` (Zeilen 446 bis 458) liefert für `claude-code` zusätzlich `CLAUDE_CODE_OAUTH_TOKEN`; der Doc-Kommentar nennt die feste Reihenfolge (Zeilen 438 bis 445) |
| `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl` | unverändert (Weg A) | Zeilen 33 bis 35 rendern `.AgentEnv` per `range`; eine weitere Variable braucht keine Vorlagenänderung |
| `internal/hexagon/application/doctor_sandbox.go` | update (Weg A, Regel offen) | `checkSandboxCredentials` (ab Zeile 134) nutzt `agentCredentialEnv` (Zeile 147) und `missingAgentSources`; der Hinweistext (Zeile 144) nennt die Variablen; `credentialPatterns` (Zeile 54) deckt das Token als Klartext schon ab |
| Quelle von `u-boot devcontainer shell` | prüfen | Laut Handbuch gibt `shell` Variablen weiter, auf die `devcontainer.json` verweist; im Code nicht gelesen |
| `internal/hexagon/application/generate_features_test.go`, `doctor_sandbox_test.go`, `devcontainer_shell_test.go`, `generate_forge_cli_test.go`, `generate_ollama_test.go`, `generate_ollama_host_test.go` | update (Weg A) | Sie enthalten `ANTHROPIC_API_KEY` (Wortsuche, nicht auf die Aussage gelesen); jeder `remoteEnv`-Erwartungswert mit `claude-code` ändert sich — prüfen |
| `internal/hexagon/application/initproject.go`, `config.go` und das Schema | update (Weg B) | neuer Schlüssel und `config set`-Pfad; nicht gelesen — prüfen |
| `spec/spezifikation.md` | update (nach dem Code) | [`SPEC-044`](../../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) Zeilen 608 und 624; [`SPEC-051`](../../../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) Zeile 769 (prüfen) |
| `spec/lastenheft.md` | **unverändert**, prüfen | [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zeilen 489 bis 495) nennt keine Variable des Agenten einzeln; bei Weg B ein CR-Slice vor dem Code (§4) |
| `docs/user/benutzerhandbuch.md` | update | Zeilen 1264, 1914 bis 1930, 2410, 2433 bis 2440 und die Variablentabelle bei Zeile 2651 |
| `README.md`, `README.de.md`, `CHANGELOG.md`, `harness/replay.md` | update (prüfen) | je nach Weg; `harness/replay.md` Zeile 57 führt einen Golden Case mit `claude-code` (nur die Treffer-Zeile gelesen) |
| `docs/plan/adr/` | neu, wenn ein Weg das verlangt | eine ADR hält die Wahl und die Zusage-Menge (Vorbild: [ADR-0036](../../adr/0036-cloud-modelle-von-ollama-im-devcontainer.md)) |
| diese Datei, `docs/plan/planning/open/README.md` | update | Zeile dieses Slice; Moves per `make slice-mv` |

**Abgleich mit `git diff --stat` (Weg A, nach der Umsetzung):**

| Datei | Art | Beleg |
|---|---|---|
| `internal/hexagon/application/devcontainer_sandbox.go` | update | `agentCredentialEnv` liefert `ANTHROPIC_API_KEY`, `CLAUDE_CODE_OAUTH_TOKEN`; Doc-Kommentar nennt die Folge |
| `internal/hexagon/application/doctor_sandbox.go` | update | nur der Hinweistext in `checkSandboxCredentials` (Beispiel-Variablen); `missingAgentSources` und die Regel unverändert gelesen |
| `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl` | unverändert | `range .AgentEnv` (Zeilen 33 bis 35) rendert jede Variable der Liste |
| `internal/hexagon/application/devcontainer_shell_test.go` | update (nur Kommentar und Fallname) | `sessionConfig` ist eine Handkonfiguration (Kommentar berichtigt, Review F-3), `TestSessionEnvNames` liest die Namen aus dem Text der Datei (Zeilen 20 bis 60 gelesen); der e2e-Test misst die Weitergabe mit dem echten Namen |
| `generate_features_test.go`, `generate_forge_cli_test.go`, `generate_ollama_host_test.go`, `generate_ollama_test.go`, `doctor_sandbox_test.go` | update | Erwartungswerte von `remoteEnv` mit `claude-code`; Doctor-Fälle „only the API key reference“, „only the OAuth token reference“ |
| `internal/e2e/claude_oauth_token_docker_test.go` | neu | Messung (Variable gesetzt, ungesetzt, leer; `claude -p` mit Platzhalter) |
| `docs/messungen/slice-v1-devcontainer-claude-oauth-token/`, `docs/messungen/README.md` | neu / update | Protokoll und Rohdaten; ein Skript unter `tools/messung/` entfällt, die Messung ist der Test |
| `.devcontainer/devcontainer.json` | update | die vom Generator erzeugte Datei dieses Repos trägt die neue Referenz |
| `spec/spezifikation.md` | update | [`SPEC-044`](../../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) (Absatz Zugangsdaten, Prüfung Punkt 3), [`SPEC-051`](../../../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) (Reihenfolge) |
| `spec/lastenheft.md` | unverändert | [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) nennt die Variable des Agenten ohne Namen, gelesen |
| `docs/user/benutzerhandbuch.md` | update | Abschnitt Claude Code, Zugangsdaten, Schritt der Sandbox-Anleitung, Variablentabelle |
| `docs/maintainer/zusagen-und-messstaende.md` | update | Messstand und Grenze des Tokens |
| `README.md`, `README.de.md` | unverändert | Zeile 153 bzw. 156 nennt `claude-code` ohne Zugangsdaten, gelesen |
| `CHANGELOG.md`, `harness/replay.md` | update | Eintrag `[Unreleased]`; Golden-Case-Zeile `claude-code` |
| `docs/plan/adr/` | unverändert | keine ADR nötig (Entscheidung oben) |
| diese Datei | update | `Verantwortlich`, Entscheidung, Abgleich; `roadmap.md` und `open/README.md` macht der Auftraggeber |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; der Projektinhaber hat den Weg (A, B oder C) entschieden; kein anderer Slice in `in-progress/`.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Weg B berührt das Lastenheft — dann ein Lastenheft-CR-Slice vor dem Code, mit Wortlaut-Freigabe des Projektinhabers.
- `in-progress` → `open` (blockiert): Der Projektinhaber entscheidet den Weg nicht, oder die Messung mit dem Platzhalter ist im Container nicht wiederholbar.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; `make gates` ist grün; die Zusage in Spezifikation und Handbuch nennt nur das Gemessene (Platzhalter oder, wenn gefahren, gültiges Token); Review-Report und Verification-Report liegen vor; die Closure-Notiz mit Steering-Loop-Lerneintrag ist geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

1. **`doctor`-Regel bei zwei Anmeldewegen** (API-Schlüssel und OAuth-Token sind Alternativen): verlangt `u-boot doctor` beide Referenzen, wäre jede Warnung falsch, sobald nur ein Weg genutzt wird. — **Ausgang:** entfallen. `checkSandboxCredentials` liest die Referenzen in `devcontainer.json`, nicht die Host-Umgebung; ein Nutzer mit einem Weg bekommt keine Warnung, solange die Datei neu erzeugt ist (Pin: `TestDoctor_SandboxCredentials_Agent`, Fälle „only the API key reference“, „only the OAuth token reference“). Folge, die bleibt: Eine ältere `devcontainer.json` ohne die Token-Referenz löst nach dem Upgrade eine `doctor`-Warnung aus, auch bei Nutzern des API-Schlüssels; sie ist in Slice, Spezifikation, Handbuch und `CHANGELOG.md` benannt, der Hinweis führt zu `u-boot generate devcontainer`. Akzeptiert (Review F-5).
2. **Ein langlebiges OAuth-Token in der Umgebung eines privilegierten Containers** (das Sandbox-Profil mit verschachteltem Docker läuft mit `--privileged`). — **Ausgang:** entfallen (die Dokumentation nennt Widerruf und Reichweite des Tokens), Handbuch Abschnitt Claude Code, Punkt Zugangsdaten (`3e3ed1c9`, Verifier F-1): Widerruf, Reichweite (jeder Prozess der Sitzung) und `--privileged` bei verschachteltem Docker stehen dort. Offen bleibt, dass die Vererbung an Prozesse der Sitzung kein Test pinnt (§7, Lerneintrag).
3. **Wechselwirkung mit dem Ollama-Ablauf:** Er ist nur bei ungesetztem `ANTHROPIC_API_KEY` gemessen; ein gesetztes `CLAUDE_CODE_OAUTH_TOKEN` ist nicht gemessen. — **Ausgang:** weiter offen (Stand 2026-11-01): [ADR-0036](../../adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) bleibt wörtlich gültig (Bedingung nur für `ANTHROPIC_API_KEY`), das Handbuch rät, die Variable im Ollama-Ablauf nicht zu setzen; das Verhalten ist nicht gemessen, Anker `docs/messungen/slice-v1-devcontainer-claude-oauth-token/README.md` Abschnitt Grenzen.
4. **Zusage weiter als die Messung** ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)): Ein gültiges Token kann der Agent nicht messen. — **Ausgang:** weiter offen (Stand 2026-11-01): `claude -p` mit dem Platzhalter lief einmal statt zweimal, ein Lauf mit gültigem Token fehlt; eine zweite Messung läuft nicht lokal (Entscheidung des Projektinhabers). Handbuch und Wartungsdoku nennen nur das Gemessene; die DoD-Zeile *Messung* ist auf diesen Stand angepasst (§7). Anker: `docs/messungen/slice-v1-devcontainer-claude-oauth-token/README.md` Abschnitt Grenzen und [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md).
5. **Golden Cases** (`harness/replay.md`): jeder `remoteEnv` mit `claude-code` ändert sich bei Weg A. — **Ausgang:** entfallen: die Zeile `claude-code` in `harness/replay.md` führt beide Referenzen, gedeckt durch `generate_features_test.go` (Review, Negativbefund `harness/replay.md`).

## 7. Closure-Notiz

Reports: [Review](../../../reviews/2026-11-01-slice-v1-devcontainer-claude-oauth-token.md) (erster Lauf 0 HIGH, 1 MEDIUM F-1, 3 LOW, 2 INFO; Nachlauf nach der Fix-Runde `c2e68d7b` 0 HIGH, 0 MEDIUM, 2 LOW N-1 und N-2, beide behoben), [Verification](../../../reviews/2026-11-01-slice-v1-devcontainer-claude-oauth-token-verify.md) (0 HIGH, 0 MEDIUM, 1 LOW F-1, 3 INFO; F-1 behoben in `3e3ed1c9`). Messprotokoll: [`docs/messungen/slice-v1-devcontainer-claude-oauth-token/`](../../../messungen/slice-v1-devcontainer-claude-oauth-token/README.md). `make gates` Exit 0 im Lauf des Verifiers (Worktree) und im Lauf der Closure auf `main`.

- **Was hat funktioniert:** Weg A ist wie entschieden umgesetzt (Verification §2: alle Entscheidungen konform): `agentCredentialEnv` liefert für `claude-code` `ANTHROPIC_API_KEY` und `CLAUDE_CODE_OAUTH_TOKEN`, `doctor` verlangt beide Referenzen in der Datei, die Vorlage blieb unverändert. Der Plan-Abgleich in §3 führte alle Dateien des Diffs (20 laut Verification §2), und die „unverändert“-Zeilen trugen Fundstellen; der Review fand die Lücke trotzdem in der Handbuch-Aussage, nicht im Abgleich. Der Docker-Test `TestE2E_ClaudeOAuthToken_ReachesTheSession` lief im Lauf des Implementers grün (Zellen gesetzt, ungesetzt, leer je zweimal; `claude -p` einmal mit Platzhalter, Antwort `401 OAuth access token is invalid`).
- **Was ging anders als geplant:**
  1. **DoD-Anpassung (Entscheidung des Projektinhabers 2026-11-01, „DoD anpassen“).** Die Zeile *Messung* verlangte mindestens zwei Läufe je Zelle, auch für `claude -p`, und ein Skript unter `tools/messung/`. Gemessen sind die Variablen-Zellen zweimal, `claude -p` einmal; ein zweiter Lauf ist lokal nicht gefahren (Entscheidung des Projektinhabers: `make test-docker` nicht lokal), ein gültiges Token fehlt. Die Zeile steht jetzt auf dem Gemessenen, das Skript entfällt zugunsten des Tests. Das ist eine Plan-Änderung, keine Gate-Lockerung: Kein Gate, keine Coverage-Schwelle und keine Architekturregel wurde angefasst (`AGENTS.md` §3.6); der Anker ist diese Zeile, das Drift-Log der Roadmap trägt nur Umplanungen von Wellen und bekommt keinen Eintrag. Was nicht gemessen ist, bleibt als Risiko 4 offen und im Handbuch als Grenze, nicht als Zusage.
  2. **Docker-Nachweis des Gesamtlaufs nicht erbracht.** Im Lauf des Implementers waren drei fremde Tests rot (`TestDevcontainerShell_Docker_Session`, `TestE2E_ForgeCLI_PinMismatchBreaksBuild/GLAB`, `TestE2E_LHAK003_KeycloakAcceptanceFlow`); keiner berührt `agentCredentialEnv`. Die Ursache (Parallelität auf der geteilten Engine) ist vermutet, nicht bewiesen. Ein lokaler Wiederholungslauf ist auf Entscheidung des Projektinhabers nicht gefahren; der CI-Lauf auf `main` entscheidet, Anker [`slice-v1-test-docker-in-ci`](../open/slice-v1-test-docker-in-ci.md). Der Slice führt den Docker-Nachweis des Gesamtlaufs nicht als erfüllt; belegt ist der Docker-Test dieses Slice im Lauf des Implementers, weder vom Reviewer noch vom Verifier wiederholt.
  3. **Folge für Bestandsprojekte:** Eine ältere `devcontainer.json` ohne die Token-Referenz löst nach dem Upgrade eine `doctor`-Warnung aus (Risiko 1); in Slice, Spezifikation, Handbuch und `CHANGELOG.md` benannt.
  4. Zwei Plan-Zeilen hatten nach der Fix-Runde nicht zum Abgleich gepasst (Review N-2: Skript entfällt) und eine Kommentar-Aussage verwies auf die Generator-Ausgabe, die der Diff geändert hat (Review F-3); beide berichtigt.
- **Steering-Loop-Eintrag:** Eine Lehre, gezählt und nicht verkörpert (kein `liegt in`): Eine Aussage im Handbuch über die **Reichweite** eines Wertes („in der Umgebung jedes Prozesses im Container“, Review F-1; „nicht in der Container-Konfiguration“, Nachlauf N-1; „Jeder Prozess in dieser Sitzung“, Verification F-1 nach der Berichtigung) ging in drei Fassungen über das hinaus, was ein Test pinnt. Gemessen ist: die Sitzung von `u-boot devcontainer shell` sieht die Variable (Shell-Ebene, `VALUE:`-Zeilen); `Config.Env` ist für `GIT_TOKEN`, den Agenten-Schlüssel und die Forge-Token mit Kanarienwert gepinnt, nicht für `CLAUDE_CODE_OAUTH_TOKEN`. Die Klasse ist die in `AGENTS.md` §3.10 verkörperte (Bedingung in Spec oder Handbuch ohne pinnenden Test), erneuter Rückfall nach der Verkörperung, hier in der Form „Aussage über Reichweite“; die Prüfung gegen die bestehende Beobachtung [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) ergibt: gleiche Klasse, Beleg ergänzt, keine neue Regel. **Offener Rest:** Der Handbuchsatz „Jeder Prozess in dieser Sitzung sieht die Variable, auch jeder Befehl, den ein Agent dort ausführt“ ist durch keinen Test an einem Kindprozess gepinnt (der Test liest die Variable in der Sitzungs-Shell); ob er bleibt oder durch einen Test gepinnt wird, ist Sache eines Folge-Eintrags im Register, nicht dieses Slice. Dazu eine **benannte Lücke ohne Spec-Änderung:** Der Wortlaut der DoD-Zeile *Messung* verlangte mehr, als der Plan zu messen vorsah (zwei Läufe `claude -p` bei einem Test von 26 Minuten auf geteilter Engine); die Klasse „Zusage/DoD weiter als die Messung“ ist die Beobachtung [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md).
- **Beobachtungs-Register (`../observations/`):** keine neue Beobachtung angelegt; zwei Evidence-Dateien `evidence/slice-v1-devcontainer-claude-oauth-token.md` ergänzt, `state.md` je nachgeführt, kein Zähler gesetzt und kein Ausgang zugewiesen:
  - [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md): Review F-1, Nachlauf N-1, Verification F-1 (Handbuch-Reichweite ohne Test); Rückfall nach der Verkörperung.
  - [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md): DoD-Zeile *Messung* weiter als die Messung (Review F-4 und F-6, Verification Abschnitt 4); durch Anpassung der Zeile aufgelöst, die Grenzen des Messprotokolls bleiben.
- **Folge-Slices:** [`slice-v1-test-docker-in-ci`](../open/slice-v1-test-docker-in-ci.md) (Entwurf in `open/`, Anker des offenen Docker-Nachweises). Kein weiterer Slice.
- **Risiken aus §6:** fünf Risiken mit Ausgang: drei entfallen (1 `doctor`-Regel, mit Folge für Bestandsdateien akzeptiert; 2 Reichweite und Widerruf; 5 Golden Cases), **zwei weiter offen** (3 Ollama-Ablauf mit gesetztem Token, 4 Zusage weiter als die Messung: zweiter `claude -p`-Lauf und gültiges Token), Anker `docs/messungen/slice-v1-devcontainer-claude-oauth-token/README.md` Abschnitt Grenzen.
- **Drei Paarungen:** Anker: kein `liegt in` in diesem Slice, nichts zu prüfen. Folge-Slice: [`slice-v1-test-docker-in-ci`](../open/slice-v1-test-docker-in-ci.md) liegt in `open/`. Register: die zwei Evidence-Dateien liegen unter `../observations/` (Ergebnis nach dem `git mv` geprüft, siehe Häkchen in §2).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** und die **vier Pflichtkriterien**.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden je nach Weg `internal/hexagon/application` (APP), `tools/` (WZG, Messskript) sowie `spec/`, `harness/`, `docs/plan/` und `docs/messungen/` (DOK). Alle sind in der Modus-Deklaration von [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) als Greenfield auditiert.

**Vorgelagert — offene Beobachtungen sichten:** Das Register ([`observations/`](../observations/)) ist nach den Titeln der Verzeichnisse durchgesehen, nicht inhaltlich gelesen. Titeltreffer: [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) und [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) betreffen die Zusage in Spezifikation und Handbuch (Risiko 4); die Zähler und den Inhalt liest der Planner beim Start.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
