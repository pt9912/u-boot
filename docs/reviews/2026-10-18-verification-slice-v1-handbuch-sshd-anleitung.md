# Verification-Report: slice-v1-handbuch-sshd-anleitung — 2026-10-18 <!-- d-check:ignore (Slice-Kennung im Titel, Zitat) -->

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“, gegen Plan und DoD). Kein Review, keine Validierung, keine Reparatur.

**Gegenstand:** Slice `slice-v1-handbuch-sshd-anleitung` (`docs/plan/planning/in-progress/`), Commit-Range `0bfc46f..HEAD` (HEAD `541ad35`); Entscheidung des Projektinhabers 2026-10-18 „Das Handbuch können wir dazu verbessern, sonst nichts“. <!-- d-check:ignore (Slice-Kennung als Zitat) -->

**Eingang:** Review `docs/reviews/2026-10-18-slice-v1-handbuch-sshd-anleitung.md` (F-1, F-2 MEDIUM; F-3, F-4 LOW; F-5, F-6 INFO) und die Zusage des Implementers, alle seien behoben. Geprüft am Text des Handbuchs, nicht an der Prosa; Sensoren selbst gefahren (Host: Colima, Docker, kein Host-go). <!-- d-check:ignore (Slice-Kennung als Zitat) -->

**Gesamturteil:** Die sachliche DoD (Zeilen 1, 2, 4) ist erfüllt, alle sechs Review-Findings sind am Text behoben. Offen sind nur Closure-Zeilen, die dieser Rolle vorgelagert sind (Closure-Notiz, Register, Risiko-Ausgänge, Paarungen) und der Verifier-Report selbst als Review-Folge. Eine kleine Restdiskrepanz (D-1, INFO) blockiert die Closure nicht.

## 1. Sensoren (selbst gefahren)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | pass (Exit 0) | lint, test, `coverage-gate`, `docs-check`, `baseline-verify` grün |
| `make doc-trace` | pass | „83 Anforderung(en), 0 Waise(n)“ |
| `make docs-check` | nach Ablage dieser Datei erneut gefahren, 0 Befunde (siehe Übergabe) | — |
| `make test-docker`, `make verify-depguard` | nicht gefahren | Diff enthält ausschließlich Markdown; kein Code-, Docker- oder Importpfad berührt |
| Praxistest der Anleitung | pass, mit Näherung (Abschnitt 4) | Exit 0, `uid=1000(vscode)` |
| `grep` `LH-`, `ADR-`, `SPEC-`, `docs/plan` über `docs/user/` | 0 Treffer | — |
| unsichtbare Zeichen (U+00AD, U+200B/C/D, U+2060, U+FEFF) in `docs/user/` | 0 Treffer | — |

## 2. DoD-Abgleich (jede Zeile)

| # | DoD-Zeile | Stand | Evidence |
|---|---|---|---|
| 1 | Anleitung im Handbuch (aktivieren, Zugang, verbinden, Sicherheitshinweis; jede Aussage belegt, „nicht gemessen“ oder gestrichen) | erfüllt | Abschnitt „Login-Dienst `sshd` (optional)“ (`benutzerhandbuch.md` ab Zeile 1211): Schritt 1 (Feature aktivieren, neu aufbauen), Schritt 2 (Schlüssel bevorzugt, Passwort Alternative), Schritt 3 (Verbinden), „Sicherheitshinweis“ mit Anker `sicherheitshinweis-sshd`; Vorspann trennt „Gemessen“ (Feature: Port 2222, `sshd -T`-Werte, gesperrte Konten; Näherungsaufbau: `ssh-keygen`, 700/600, `ssh -i … -p 2222`) von „Nicht gemessen“ (Login mit dem Feature, Editor-/Engine-Weiterleitung, `docker exec`). Aussagen zu Code gegen Quelle gelesen: `forwardPorts` nur bei nicht leerer Liste (`devcontainer.json.tmpl:43`), `collectDevcontainerForwardPorts` (`generate.go:934`), Warntext (`devcontainer_sandbox.go:321`), `--quiet`-Zeile (Handbuch Zeile 310). |
| 2 | `make gates` grün; keine `LH-*`/`ADR-*`/`SPEC-*`/Slice-Kennung, keine Chronik im Handbuch | erfüllt | Gates grün; `grep` 0 Treffer; Abschnitt nennt Stand („gemessen mit …“), keine Chronik |
| 3 | Review durchgeführt, Report liegt vor | erfüllt | `docs/reviews/2026-10-18-slice-v1-handbuch-sshd-anleitung.md`; Findings-Behebung in `541ad35` | <!-- d-check:ignore (Dateipfad des Reviews) -->
| 4 | Doku-Update `docs/user/benutzerhandbuch.md` | erfüllt | Diff: 127 Zeilen im Handbuch; einziger berührter öffentlicher Vertrag |
| 5 | Closure-Notiz mit Steering-Loop-Lerneintrag | offen | §7 steht auf „Wird bei der Closure gefüllt.“ |
| 6 | Beobachtungs-Register fortgeschrieben oder „keine Beobachtung“ in §7 | offen | Closure |
| 7 | Jedes Risiko aus §6 mit Ausgang | offen | alle drei Ausgänge stehen auf „offen“ (Vorschläge in Abschnitt 5) |
| 8 | Drei Paarungen (Anker, Folge-Slice, Register) | offen | Closure |

## 3. Review-Findings am Text

| ID | Befund | Stand | Beleg im Text |
|---|---|---|---|
| F-1 | Loopback-Bindung | behoben | Schritt 3: „nur an der Loopback-Adresse des Hosts, nach dem Muster `127.0.0.1:2222:2222`, nicht als `2222:2222`“ mit Begründung; Sicherheitshinweis wiederholt die Regel (`127.0.0.1`). Der Praxistest bindet so und `docker port` zeigt `2222/tcp -> 127.0.0.1:2222`. Die Bindung selbst bleibt als nicht gemessen geführt (Plan-Tabelle). |
| F-2 | Fremdwerkzeug-Aussagen ohne Beleg | behoben | Aussagen zu `ssh-keygen`, 700/600, `ssh -i` tragen den Beleg „Praxistest des Reviewers 2026-10-18“ (Plan-Tabelle) und sind im Text als Näherungsaufbau gekennzeichnet; `docker exec`, Editor-/Engine-Weiterleitung stehen als „nicht gemessen“. Der Verifier hat den Praxistest unabhängig wiederholt (Abschnitt 4). Die frühere Annahme-Aussage zu den Rechten („Annahme über die Prüfung“) ist aus dem Text verschwunden. |
| F-3 | Schlüssel bevorzugt | behoben | Schritt 2: „Bevorzugt ist der Schlüssel; das Passwort ist die Alternative“; Passwort mit „starken Passwort“ und Hinweis auf Root-Login |
| F-4 | „neu aufbauen“ nur als belegt | behoben | Schritt 1 stützt sich auf „Installationsskript läuft beim Image-Build“ (Handbuch, Abschnitt Claude Code); ob ein Neustart genügt, steht als „Nicht gemessen“; Abgrenzung zum „neu starten“ der Egress-Beschränkung. Der Satz „Bauen Sie … neu auf“ bleibt eine Folgerung aus belegter Prämisse und ist nicht als gemessen ausgegeben. |
| F-5 | Messaussage `command -v ssh-keygen` | behoben (anders) | Die Klammer „gemessen mit `command -v`“ steht nicht mehr im Text; Rest-Aussage „`ssh-keygen` liegt im Basis-Image“ (Schritt 2, siehe D-1). Die Messung des Reviewers (`/usr/bin/ssh-keygen` im Image) ist im Review dokumentiert; der Verifier hat sie nicht wiederholt (der Praxistest erzeugte den Schlüssel auf dem Host). |
| F-6 | „unverändert“-Zeilen ohne gelesene Aussage | behoben | Plan-Tabelle: `open/README.md` mit gelesener Aussage (Zeiger `../in-progress/…`, Gegenstand zutreffend); gegen den Baum geprüft: Zeile 15 der Datei trägt genau diesen Zeiger. `CHANGELOG.md`, README-Dateien mit gelesener Aussage. |

## 4. Praxistest (Näherung benannt)

Aufbau: Wegwerf-Container aus `mcr.microsoft.com/devcontainers/base:debian`, Port veröffentlicht mit `docker run -p 127.0.0.1:2222:2222`, im Container das Debian-Paket `openssh-server` auf Port 2222 (Root- und Passwort-Login `yes`, Konten `root`/`vscode` mit `passwd -l` gesperrt). Dann die Schritte der Anleitung: `ssh-keygen -t ed25519 -f …` auf dem Host, `~/.ssh` 700 und `authorized_keys` 600 als `vscode` im Container (`docker exec`), `ssh -i … -p 2222 vscode@localhost` mit `BatchMode`.

Ergebnis: Listener `0.0.0.0:2222` und `[::]:2222` im Container, Login erfolgreich (`uid=1000(vscode)`, Exit 0), `docker port` zeigt `127.0.0.1:2222`. Container, Schlüssel und Skript danach gelöscht, kein Container `vsshd` übrig.

**Näherung (nicht belegt durch diesen Test):** Debian-Paket statt Feature-Installationsskript; keine Dev-Container-CLI, keine Editor-Weiterleitung; Passwort-Login nicht ausgeführt; Schritt 1 (Neubau-Frage) und der Weg des Ports über einen Editor sind nicht abgedeckt. `docker exec` ist hier erstmals benutzt worden und lief, im Handbuch bleibt es „nicht gemessen“ (konservativ, nicht falsch). Die Näherung deckt sich mit der im Text genannten.

## 5. Plan-vs-Diff (§3.10)

`git diff --stat 0bfc46f..HEAD`: vier Dateien — `docs/user/benutzerhandbuch.md`, `docs/plan/planning/in-progress/roadmap.md` (−2, „Nichts in Arbeit.“ entfernt), der Slice selbst (+26/−), der Review-Report. Plan-Tabelle §3 führt alle vier (Handbuch update, Roadmap update, Slice update, Review neu); der Review-Report steht als „neu“. Keine Diff-Datei ohne Eintrag, kein Eintrag ohne Diff. „unverändert“-Zeilen (`CHANGELOG.md`, `README.md`, `README.de.md`, `open/README.md`) tragen gelesene Aussagen; sie decken sich mit dem Diff (keine Änderung dort).

## 6. Diskrepanzen

- **D-1 (INFO):** Schritt 2 sagt „Auf dem Host ein Schlüsselpaar erzeugen (`ssh-keygen` liegt im Basis-Image …)“. Der Nachsatz belegt das Werkzeug im Container, nicht auf dem Host, auf dem der Befehl läuft; die Aussage ist wahr, aber für die Handlung ohne Bezug. Kein Gate, keine Falschaussage; bei Gelegenheit streichen oder auf „auf dem Host (`ssh-keygen` gehört zu OpenSSH)“ ändern — letzteres wäre eine neue Fremdwerkzeug-Aussage und braucht Beleg.
- **D-2 (INFO):** Die Egress-Tabelle (Zeile 1179) und die Feature-Tabelle (Zeile 797) sind unverändert und widersprechen dem neuen Abschnitt nicht (Querverweis `#login-dienst-sshd-optional` löst auf).

## 7. Auflagen für die Closure

- **A-1 Closure-Notiz (§7) mit Steering-Loop-Lerneintrag:** Finding-Klassen des Reviews (Anleitung-Rat ohne Sicherheitsgrenze; Fremdwerkzeug-Aussage ohne Quelle; Sicherheitsaussage gleichrangig neben bequemerem Weg; Folgerung als Anweisung, Begriff uneinheitlich; Messaussage nicht vom Autor wiederholt; „unverändert“-Zeile ohne gelesene Aussage); Grenzen benennen: Feature-Login und Weiterleitung nicht gemessen, Loopback-Bindung nicht gemessen; D-1. Die Klassen „Fremdwerkzeug-Aussage“ und „unverändert ohne gelesene Aussage“ am Register zählen (ein Vorgang zählt einmal), keine „3×“ behaupten, wo es weniger sind. DoD-Zeilen 5 bis 8 erst nach Schreiben abhaken.
- **A-2 Risiko-Ausgänge §6** (je genau einer; Vorschlag): „Aussage ohne Beleg“ → entfallen für die gemessenen Aussagen (Praxistest, Verifier wiederholt), weiter offen für die als „nicht gemessen“ gekennzeichneten (Feature-Login, Editor-/Engine-Weiterleitung, Neustart statt Neubau); „Zusage durch die Hintertür“ → entfallen (Text nennt den Stand „u-boot veröffentlicht … nicht“, keine Garantie); „Sicherheitsaussage verwässert“ → entfallen (Hinweis vor Schritt 1 verlinkt, Loopback in Schritt 3 und im Hinweis, Schlüssel bevorzugt).
- **A-3 Beobachtungs-Register:** Vorgang an `BEO-DOK/plan-aussage-ueber-code` (verkörpert, `AGENTS.md` §3.10) fortschreiben, soweit der Rückfall (Fremdwerkzeug-Aussage ohne Beleg, „unverändert“-Zeile) dazugehört; sonst in §7 „keine neue Beobachtung“ vermerken. Zählbasis ehrlich. <!-- d-check:ignore (Beobachtungs-Kennung als Zitat) -->
- **A-4 Roadmap-Ruhe-Marker nach dem Move:** Commit `git mv` nach `done/` als reiner Move (erst Inhalt, dann Move, `AGENTS.md` §3.3). Danach liegt nichts mehr in `in-progress/`: der Marker „Nichts in Arbeit.“ gehört zurück unter „Offene Wellen“ in `roadmap.md`; `open/README.md`, `done/README.md` nachziehen; Pfade in der Slice-Datei vom Ruheort aus (`../done/…`, `../observations/…`). Der Cut `slice-v1-release-cut-v0.11.0` startet erst nach dem Move. <!-- d-check:ignore (Slice-Kennung als Zitat) -->
- **A-5 CHANGELOG:** unverändert lassen; der Release-Cut entscheidet, ob er die Handbuch-Verbesserung nennt (Plan-Tabelle).

## 8. Nicht-Befunde

Keine Verletzung von §3.1 (nur `make` und Docker; Praxistest ausschließlich in Docker, Schlüssel nur im Scratchpad), §3.4 (keine Architektur-Änderung), §5 Zeile 6 (keine Kennungen im Handbuch).
