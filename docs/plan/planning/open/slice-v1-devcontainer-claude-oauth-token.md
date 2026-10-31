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

**Verantwortlich:** —

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

- [ ] **Messung** ([`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)): mit einem Platzhalter-Wert kommt die Variable über `u-boot devcontainer shell` im Container an; ungesetzt auf dem Host kommt sie ungesetzt an, leer gesetzt als leere; `claude -p` mit dem Platzhalter endet mit einer Antwort des Dienstes (kein gültiges Token). Ein Lauf mit einem **gültigen** Token ist nur zugesagt, wenn der Projektinhaber ihn gefahren hat und das Protokoll ihn nennt. Mindestens zwei Läufe je Zelle, datiert, mit Umgebung. *Beleg:* Messprotokoll unter `docs/messungen/<kennung>/` mit Rohdaten, Skript unter `tools/messung/<name>/` ([`harness/conventions.md`](../../../../harness/conventions.md#ablage-von-messprotokollen-und-skripten)); keine Zugangsdaten in den Rohdaten.
- [ ] **Der gewählte Weg** (A, B oder C): bei A oder B reicht `u-boot generate devcontainer` die Variable als Referenz durch, `u-boot doctor` prüft sie nach der festgelegten Regel, Tests je Bedingung und je Mutation; bei C steht der Ablauf mit Datum und Messstand im Handbuch.
- [ ] **Spezifikation, Handbuch, `CHANGELOG.md` nach dem Code** mit Beleg-Tabelle (Bedingung → Test oder Streichung); Handbuch und READMEs ohne Kennungen ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6); der Satz „wird weder eingerichtet noch geprüft“ in [`SPEC-044`](../../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) entfällt oder bleibt, je nach Weg.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register ([`observations/`](../observations/)) fortgeschrieben, oder in §7 „keine Beobachtung angefallen“.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

Gelesen am 2026-10-31; Zeilennummern verschieben sich, der Implementer liest sie am Ist-Stand nach ([`AGENTS.md`](../../../../AGENTS.md) §3.10). Welche Zeilen der Weg berührt, entscheidet der Projektinhaber.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/messung/<name>/`, `docs/messungen/<kennung>/` | neu | Skript und Protokoll der Messung (Platzhalter; ein gültiges Token nur durch den Projektinhaber) |
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

1. **`doctor`-Regel bei zwei Anmeldewegen** (API-Schlüssel und OAuth-Token sind Alternativen): verlangt `u-boot doctor` beide Referenzen, wäre jede Warnung falsch, sobald nur ein Weg genutzt wird. — **Ausgang:** bei Closure.
2. **Ein langlebiges OAuth-Token in der Umgebung eines privilegierten Containers** (das Sandbox-Profil mit verschachteltem Docker läuft mit `--privileged`). — **Ausgang:** bei Closure; die Dokumentation nennt Widerruf und Reichweite des Tokens.
3. **Wechselwirkung mit dem Ollama-Ablauf:** Er ist nur bei ungesetztem `ANTHROPIC_API_KEY` gemessen; ein gesetztes `CLAUDE_CODE_OAUTH_TOKEN` ist nicht gemessen. — **Ausgang:** bei Closure.
4. **Zusage weiter als die Messung** ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)): Ein gültiges Token kann der Agent nicht messen. — **Ausgang:** bei Closure.
5. **Golden Cases** (`harness/replay.md`): jeder `remoteEnv` mit `claude-code` ändert sich bei Weg A. — **Ausgang:** bei Closure.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt: Was hat funktioniert · Was ging anders als geplant · Steering-Loop-Eintrag · Beobachtungs-Register · Folge-Slices · Risiken aus §6 · Drei Paarungen.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** und die **vier Pflichtkriterien**.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden je nach Weg `internal/hexagon/application` (APP), `tools/` (WZG, Messskript) sowie `spec/`, `harness/`, `docs/plan/` und `docs/messungen/` (DOK). Alle sind in der Modus-Deklaration von [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) als Greenfield auditiert.

**Vorgelagert — offene Beobachtungen sichten:** Das Register ([`observations/`](../observations/)) ist nach den Titeln der Verzeichnisse durchgesehen, nicht inhaltlich gelesen. Titeltreffer: [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) und [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) betreffen die Zusage in Spezifikation und Handbuch (Risiko 4); die Zähler und den Inhalt liest der Planner beim Start.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
