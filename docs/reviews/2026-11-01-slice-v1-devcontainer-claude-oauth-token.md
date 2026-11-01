# Review-Report: slice-v1-devcontainer-claude-oauth-token — 2026-11-01

**Review-Art:** Code (Diff gegen Plan, Entscheidung Weg A, Spec, ADR und Hard Rules)

**Gegenstand:** [`slice-v1-devcontainer-claude-oauth-token`](../plan/planning/done/slice-v1-devcontainer-claude-oauth-token.md), Commits `fbf310dc` und `c08a9216` gegenüber ihrer Basis `b7afa8ae` (18 Dateien, `git diff --stat fbf310dc~1 c08a9216`), Branch `worktree-agent-aa69746b53d530dc4`.

**Skill:** `.harness/skills/reviewer.md` (Stand des Arbeitsbaums) · Kopie der Vorlage per `cp`, dann mit Write gefüllt.
**Modell:** Sonnet 5.5 · **Datum:** 2026-11-01. `make gates` im Arbeitsbaum gefahren (Exit 0, Coverage 92,20 % bei Schwelle 90 %, `baseline-verify` und `ablage-pruefung` grün). `make test-docker` nicht gefahren (Anweisung des Projektinhabers: 26 Minuten, geteilte Engine); der Docker-Test dieses Slice ist deshalb nur über das Messprotokoll des Implementers belegt, nicht von mir wiederholt.

**Eingangs-Kontext:**

- [`slice-v1-devcontainer-claude-oauth-token`](../plan/planning/done/slice-v1-devcontainer-claude-oauth-token.md) (Entscheidung Weg A, Plan §3, DoD, Risiken §6)
- [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) (Ollama-Ablauf nur bei ungesetztem `ANTHROPIC_API_KEY`)
- [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)
- [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten), [`SPEC-051`](../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten)
- `AGENTS.md` §3.7, §3.9, §3.10, §5 Regel 6; `harness/review.md`; `docs/user/benutzerhandbuch.md`

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Das Handbuch sagt, das Token stehe „in der Umgebung jedes Prozesses im Container“. Dasselbe Handbuch sagt an anderer Stelle, in einer Shell per `docker exec` gelte `remoteEnv` nicht, und [`SPEC-051`](../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) hält fest, der Wert stehe nicht in der Umgebung des Containers (`Config.Env`). Gemessen ist nur die Sitzung von `u-boot devcontainer shell`; „jedes Prozesses“ ist weder gemessen noch durch einen Test gepinnt. Szenario: Ein Nutzer öffnet eine Shell mit `docker exec` und erwartet das Token dort, oder er hält den Container für gesichert, weil er die Reichweite anders liest. | `AGENTS.md` §3.10 (Aussage über Verhalten ohne Test oder Messung), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) | `docs/user/benutzerhandbuch.md` · „Das Token steht in der Umgebung“ | nein — Review-only (`make docs-check` sieht es nicht) | Handbuch-Aussage über Verhalten weiter als die Messung |
| F-2 | LOW | Zwei Handbuch-Sätze und ein Satz der Spezifikation beschreiben die weitergegebenen Variablen noch im Singular oder mit „beide“: „der Schlüssel eines Agenten“ (Abschnitt *devcontainer shell*, Punkt Umgebung), „dem Schlüssel des Agenten“ (Abschnitt *Zugangsdaten*, Punkt zu `shell`) und „gibt `u-boot devcontainer shell` beide Variablen an die Sitzung weiter“ in [`SPEC-051`](../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten). Mit `claude-code`, `gh` und `glab` sind es bis zu vier Referenzen; der Diff änderte die Zahl, die Sätze blieben. | `AGENTS.md` §3.9 (b) | `docs/user/benutzerhandbuch.md` · „dem Schlüssel des Agenten und den Token-Variablen von `gh`“; `spec/spezifikation.md` · „gibt `u-boot devcontainer shell` beide Variablen“ | nein — Review-only | Zahlenaussage nach Listenerweiterung nicht nachgezogen |
| F-3 | LOW | Der Kommentar an `sessionConfig` und der Testfall „the shape the generator writes“ behaupten, die Konfiguration entspreche dem, was der Generator für das Sandbox-Profil schreibt. Mit `claude-code` schreibt er jetzt zusätzlich `CLAUDE_CODE_OAUTH_TOKEN`; die Handkonfiguration trägt diese Referenz nicht. Der Plan führt die Datei als „unverändert“ mit dieser Begründung (Handkonfiguration), den Kommentar hat der Diff nicht berichtigt. | `AGENTS.md` §3.9 (b), §3.7 | `internal/hexagon/application/devcontainer_shell_test.go` · „sessionConfig is a configuration in the shape the generator writes“ | nein — Review-only | Kommentar verweist auf Generator-Ausgabe, die der Diff geändert hat |
| F-4 | LOW | Die Messung weicht von der DoD-Zeile ab: Sie verlangt „Skript unter `tools/messung/<name>/`“, der Implementer setzt den Docker-Test an dessen Stelle (Plan-Abgleich: „Skript entfällt“). Die DoD-Zeile selbst ist nicht angepasst. | `harness/conventions.md` §Ablage von Messprotokollen und Skripten | Slice-Datei, Abschnitt DoD, · „Skript unter `tools/messung/<name>/`“ | nein — Review-only | Abweichung vom Plan im Abgleich, nicht in der DoD |
| F-5 | INFO | Die Entscheidung „`doctor` verlangt beide Referenzen“ ist mit [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) (Punkt 3 nachgezogen) und den Tests konsistent; [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) nennt „die Umgebungsvariable des Agenten“ ohne Zahl und `warn` ohne Abbruch. Folge für Bestandsprojekte: Eine mit der früheren Fassung erzeugte `devcontainer.json` löst nach dem Upgrade eine Warnung (mit `SecurityRelaxation`) für `CLAUDE_CODE_OAUTH_TOKEN` aus, auch bei Nutzern, die nur den API-Schlüssel verwenden. Das ist in Slice, Spezifikation, Handbuch und CHANGELOG benannt (Hinweis: neu erzeugen); kein Abbruch, kein Exit-Code. Der Ollama-Ablauf nach [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) bleibt unberührt (Bedingung nur für `ANTHROPIC_API_KEY`); das Verhalten mit gesetztem Token dort ist ausdrücklich nicht gemessen und nicht zugesagt. | [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) | `internal/hexagon/application/doctor_sandbox.go` · „missingAgentSources“ | ja — `make test` (Fälle „only the API key reference“, „only the OAuth token reference“) | — (Abwägung, keine Änderungspflicht) |
| F-6 | INFO | Die Messung ist unvollständig und die DoD sagt das ehrlich: Zeile „Messung“ bleibt `[ ]`, im Protokoll steht `claude -p` mit einem Lauf (gefordert: mindestens zwei je Zelle) und kein Lauf mit gültigem Token. Folge: Der Slice kann erst schließen, wenn ein zweiter `claude -p`-Lauf mit Platzhalter gefahren ist (Implementer) oder der Planner die DoD-Zeile ändert; ein gültiges Token ist nur zugesagt, wenn der Projektinhaber es fährt. Risiko 4 (Zusage weiter als Messung) trägt bis dahin keinen Ausgang; Handbuch und Wartungsdoku nennen die Grenzen („nicht gefahren“, „nicht gemessen“) als Befund, nicht als Zusage. Außerdem schlugen im Lauf von `make test-docker` drei fremde Tests fehl (im Protokoll benannt); der Docker-Test dieses Slice ist grün, aber nicht in einem sauberen Paketlauf belegt. Zuständig: Verifier und Planner. | Slice-DoD, Risiko 4 | Messprotokoll des Slice, README, · „Grenzen“ (Abschnitt mit der Zeile zu `claude -p`) | Review-only | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Reihenfolge `ANTHROPIC_API_KEY`, `CLAUDE_CODE_OAUTH_TOKEN`, `GH_TOKEN`, `GITLAB_TOKEN` | geprüft, ohne Befund: Doc-Kommentar und Liste in `agentCredentialEnv`, [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten), [`SPEC-051`](../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten), `harness/replay.md` Zeile `claude-code`, und die Tests (`generate_features_test.go`, `generate_forge_cli_test.go`, `generate_ollama_host_test.go` mit Reihenfolge-Literal, `generate_ollama_test.go`, `doctor_sandbox_test.go`) tragen dieselbe Folge; Suche im Baum nach dem alten Stand fand außer F-3 keine Stelle |
| `.devcontainer/devcontainer.json` (von Hand geändert) | geprüft, ohne Befund: Zeile und Position entsprechen der Vorlage `devcontainer.json.tmpl` (`range .AgentEnv`, Komma-Folge) und der Reihenfolge der Liste; der Generator selbst wurde nicht ausgeführt (kein Host-Go), die Gleichheit ist aus Vorlage und Liste gelesen, nicht gefahren |
| `internal/hexagon/application/devcontainer_sandbox.go`, `doctor_sandbox.go` (Kommentare) | geprüft, ohne Befund: Indikativ, keine Chronik, kein Konjunktiv über Verworfenes (§3.7, §3.9 a, c) |
| `internal/e2e/claude_oauth_token_docker_test.go` | geprüft, ohne Befund: Build-Tag `docker`, Platzhalter erfunden, Prüfung, dass der Platzhalter nicht in `devcontainer.json` und nicht in `stderr` steht; Kommentare tragen Zusage und Grenze |
| Sicherheit: Token in Dateien, Protokoll, Rohdaten | geprüft, ohne Befund: `rohdaten-01.txt` enthält nur Testzeilen und die Antwort „401 OAuth access token is invalid“, keinen Wert; README nennt nur das Präfix des Platzhalters; `devcontainer.json` trägt nur die Referenz |
| Handbuch: Kennungen, Zusage-/Mess-Sprache (`AGENTS.md` §5 Regel 6) | geprüft, ohne Befund außer F-1 und F-2: keine `LH-`/`SPEC-`/`ADR-`-Kennung im Diff des Handbuchs; „Verhalten im Ollama-Ablauf nicht gemessen“ steht nicht im Handbuch als Zusage, sondern als Hinweis „setzen Sie dort auch `CLAUDE_CODE_OAUTH_TOKEN` nicht“ und in der Wartungsdoku als Grenze |
| Spec-Bedingungen mit Test (`AGENTS.md` §3.10) | geprüft, ohne Befund: Rendern beider Referenzen und deren Zahl (`generate_features_test.go`, `generate_ollama_test.go`), Reihenfolge (`generate_forge_cli_test.go`, `generate_ollama_host_test.go`), Doctor „beide verlangt, nur die fehlende genannt“ (`doctor_sandbox_test.go`), Weitergabe ungesetzt/leer/gesetzt (nur im Docker-Test; nicht im Gate `make gates`, nicht von mir gefahren) |
| Plan-Tabelle gegen `git diff --stat` (§3.10) | geprüft, ohne Befund: alle 18 Dateien des Diffs stehen im Abgleich (Quelle, Doctor, fünf Testdateien, Docker-Test, Handbuch, Spec, CHANGELOG, `replay.md`, Wartungsdoku, Protokoll mit `docs/messungen/README.md`, `.devcontainer/devcontainer.json`, Slice); „unverändert“-Zeilen (Vorlage, `devcontainer_shell_test.go`, README, Lastenheft) tragen Fundstelle; Ausnahme siehe F-3 |
| `harness/replay.md` | geprüft, ohne Befund: Änderung nur an der `claude-code`-Zeile (zwei Referenzen statt einer), gedeckt durch `generate_features_test.go`; die Zeilen für `gh`/`glab` und `ollama-host` nennen keine `claude-code`-Referenz und bleiben stimmig. Kein Gate gelockert |
| Referenzmodell, Sprachvertrag, Architekturgrenzen, Lifecycle | geprüft, ohne Befund: Spec verweist nicht abwärts auf Slice; Meldungstext (`doctor`) englisch; keine Import- oder Port-Änderung; Slice bleibt in `open/` |
| `make gates` | gefahren im Arbeitsbaum, Exit 0 |
| `make test-docker`, `make verify-depguard` | nicht gefahren: erstes laut Anweisung des Projektinhabers untersagt (Dauer, geteilte Engine), zweites nicht nötig (keine Importregel berührt) |
| Slice-Closure (§7, Risiken §6) | nicht geprüft: Slice nicht geschlossen, §7 leer und Risiken ohne Ausgang sind vor der Closure zulässig |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 3 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Handbuch-Aussage über Verhalten weiter als die Messung · Zahlenaussage nach Listenerweiterung nicht nachgezogen · Kommentar verweist auf Generator-Ausgabe, die der Diff geändert hat · Abweichung vom Plan im Abgleich, nicht in der DoD

## Verdikt

**Merge-blockierend:** ja — wegen F-1 (MEDIUM, Handbuch-Aussage über Verhalten ohne Beleg und im Widerspruch zu Handbuch und Spec). F-2 bis F-4 sind nicht blockierend. F-5 und F-6 ohne Änderungspflicht am Code; F-6 ist Sache von Verifier und Planner und hält die DoD-Zeile „Messung“ offen.

**Übergabe:** Findings gehen an den Implementer; F-6 zusätzlich an den Verifier. Die Finding-Klassen gehen in die Slice-Closure §7. Dieser Report ersetzt keine Verifikation.

---

## Nachlauf

**Gegenstand:** Fix-Runde `c2e68d7b` (5 Dateien: Slice, Handbuch, `devcontainer_shell_test.go`, Spezifikation, dieser Report) gegen F-1 bis F-4. Geprüft sind nur die Fix-Runde und ihre neuen Vertragstexte. **Sensoren:** `make docs-check` im Arbeitsbaum: 952 Dateien, 0 Befunde. `make gates` nicht gefahren (der Code-Diff der Runde ist ein Kommentar und ein Testfall-Name); `make test-docker` nicht gefahren (Anweisung).

### Stand der Findings

| Finding | Stand | Beleg |
|---|---|---|
| F-1 (MEDIUM) | behoben bis auf N-1 (LOW) | Handbuch, Abschnitt Claude Code, Punkt Zugangsdaten: „Das Token steht in der Umgebung der Sitzung von `u-boot devcontainer shell`“; Sitzung gemessen durch `TestE2E_ClaudeOAuthToken_ReachesTheSession`; „in einer Shell, die Sie mit `docker exec` öffnen, gilt `remoteEnv` nicht“ steht im Handbuch (Abschnitt devcontainer shell, „In dieser Shell gilt `remoteEnv` aus `devcontainer.json` nicht“) und ist für `GIT_TOKEN` durch den Wächter „The same exec without `-e` does not carry the token“ in `internal/e2e/devcontainer_shell_docker_test.go` gepinnt |
| F-2 (LOW) | behoben | beide Handbuchstellen nennen die Variablen des Agenten; [`SPEC-051`](../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten): „jede Variable, auf die `remoteEnv` verweist (bis zu vier: …)“ stimmt mit der Liste in [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) und `agentCredentialEnv` überein; Wortsuche nach „beide Variablen“, „Schlüssel des/eines Agenten“, „jedes Prozesses“ außerhalb von Reports und archivierten Slices: keine Treffer (Archivtexte unter `done/` sind Chronik, kein Nachzug) |
| F-3 (LOW) | behoben | Kommentar von `sessionConfig` sagt „hand-written … form of the `remoteEnv` entries“, Fallname „the entry form the generator writes“; der Wortlaut des Falls steht in `done/`-Slices als Zitat (Chronik, unberührt) |
| F-4 (LOW) | teilweise behoben, N-2 (LOW) | DoD-Zeile „Messung“ benennt: Skript entfällt, Test ist die Messung, `claude -p` einmal, gültiges Token fehlt, Risiko 4 |

### Neue Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| N-1 | LOW | Der Satz „nicht in der Container-Konfiguration“ stützt sich auf [`SPEC-051`](../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) („nicht in der Umgebung des Containers (`Config.Env`)“). Gepinnt ist `Config.Env` für `GIT_TOKEN`, den Agenten-Wert und die Forge-Token (`devcontainer_shell_docker_test.go`, `forge_cli_docker_test.go`: Kanarienwerte mit Wächter), nicht für `CLAUDE_CODE_OAUTH_TOKEN`; `claude_oauth_token_docker_test.go` prüft nur `devcontainer.json` und `stderr` auf den Platzhalter. Der Mechanismus (`remoteEnv` mit `${localEnv:…}`) ist derselbe, die Variable steht aber in keiner Messung für `Config.Env`. „Container-Konfiguration“ ist zudem weiter als `Config.Env` (die `devcontainer.json` gehört dazu und trägt die Referenz). Szenario: Der nächste Änderer liest „Konfiguration“ als „keine Spur im Container“ und baut darauf. | `AGENTS.md` §3.10 (Aussage über Verhalten; Bedingung in Spec/Handbuch braucht Test) | `docs/user/benutzerhandbuch.md` · „nicht in der Container-Konfiguration“ | `make test-docker` (nicht gefahren); sonst Review-only | Handbuch-Aussage über die Reichweite stützt sich auf Pins anderer Variablen |
| N-2 | LOW | Die Plan-Tabelle in §3 führt weiter die Zeile „`tools/messung/<name>/`, `docs/messungen/<kennung>/` · neu · Skript und Protokoll der Messung“ (Datei [`slice-v1-devcontainer-claude-oauth-token`](../plan/planning/done/slice-v1-devcontainer-claude-oauth-token.md)); die DoD-Zeile und die Abgleichtabelle sagen, das Skript entfällt. Die Runde hat die DoD-Zeile berichtigt, nicht diese Zeile (Abgleich nach der Fix-Runde, `AGENTS.md` §3.10). | `AGENTS.md` §3.10 (Abgleich nach der letzten Fix-Runde) | Slice-Datei · „Skript und Protokoll der Messung (Platzhalter“ | Review-only | Plan-Zeile ohne Diff-Entsprechung |

### Zu den gestellten Fragen

- **Tragbarkeit von „nicht in der Container-Konfiguration“ ohne neue Messung:** tragbar als Aussage über `Config.Env`, nicht in dieser Wortwahl. Entweder den Satz auf „nicht in der Umgebung des Containers (`docker inspect`)“ schärfen und einen `Config.Env`-Wächter mit dem Platzhalter in den Docker-Test dieses Slice aufnehmen, oder auf die Sitzung beschränken (der erste Halbsatz trägt allein, er ist gemessen). Gestrichen werden muss er nicht; Bedingung ist N-1. Nicht blockierend, solange der Slice offen bleibt (`make test-docker` wird ohnehin für den zweiten `claude -p`-Lauf gebraucht).
- **Handbuch-Ton:** keine Kennung, keine Zusage- oder Mess-Sprache im Diff des Handbuchs; die Variablennamen stehen als Namen, nicht als Anforderungs-IDs. Geprüft, ohne Befund.
- **DoD-Zeile „Messung“ und Risiko 4 „weiter offen“:** ehrlich. Die Zeile bleibt `[ ]`, nennt Test statt Skript, `claude -p` einmal, ein gültiges Token fehlt; Risiko 4 trägt Ausgang „weiter offen“ mit Beleg-Zeiger (Messprotokoll, Abschnitt Grenzen). Deckt sich mit dem Messprotokoll (Zelle `claude -p`: 1 Lauf). Der Halbsatz „(Entscheidung des Projektinhabers)“ ist ein Verweis auf die Anweisung, `make test-docker` nicht lokal zu fahren; im Slice selbst nicht auflösbar belegt (INFO, kein Änderungszwang).
- **Reste nach altem Wortlaut (§3.9 b):** außer N-2 keine; `dcAddRemoteEnv` in `devcontainer_shell_docker_test.go` („two entries of the shape the generator writes for an agent“) bleibt zutreffend (Form des Eintrags, nicht Zahl).

### Negativbefunde Nachlauf

| Bereich | Ergebnis |
|---|---|
| `spec/spezifikation.md` (Absatz Zugangsdaten von [`SPEC-051`](../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten)) | geprüft, ohne Befund: Zahl und Namen stimmen mit [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) und Reihenfolge-Tests; kein Abwärtsverweis |
| `internal/hexagon/application/devcontainer_shell_test.go` | geprüft, ohne Befund: Kommentar im Indikativ; Fallname ändert keinen Erwartungswert |
| Handbuch, Abschnitt Zugangsdaten (Punkt zu `shell`, Standardprofil) | geprüft, ohne Befund: Aussage entspricht dem Wächter und `TestSessionEnvNames` |
| Slice §6 Risiko 4, DoD | geprüft: ehrlich, siehe oben; N-2 bleibt |

### Verdikt Nachlauf

0 HIGH, 0 MEDIUM, 2 LOW, F-6 weiter beim Verifier. Das MEDIUM F-1 ist aufgelöst; die Merge-Blockade aus dem ersten Lauf entfällt. N-1 und N-2 sind nicht blockierend, der Implementer entscheidet über Schärfen oder Streichen.
