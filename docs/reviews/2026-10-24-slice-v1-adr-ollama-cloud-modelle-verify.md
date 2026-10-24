# Verification-Report: slice-v1-adr-ollama-cloud-modelle — 2026-10-24

**Rolle:** Verifier (Modul 8) · **Frage:** Wurde der Slice richtig gebaut, gemessen an Plan, DoD, ADR und ausgeführten Sensors? (nicht: Review, nicht: Validierung) · **Modell:** claude-sonnet-5-5 · frischer Kontext.

**Gegenstand:** Slice [`slice-v1-adr-ollama-cloud-modelle`](../plan/planning/open/slice-v1-adr-ollama-cloud-modelle.md) (liegt in `open/`, DoD-Häkchen offen: setzt der Planner bei der Closure), Stand HEAD `d9dd171e` (9 Commits `origin/main..HEAD`). Geprüft gegen [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) (Accepted), den ADR-Index, die Reports [`…-modelle.md`](2026-10-24-slice-v1-adr-ollama-cloud-modelle.md) (F-1 bis F-16) und [`…-nachlauf.md`](2026-10-24-slice-v1-adr-ollama-cloud-modelle-nachlauf.md) (N-1 bis N-6), [`harness/verification.md`](../../harness/verification.md).

**Nicht Gegenstand (Auftrag):** Nachzug der Welle-Datei/Roadmap (Planner-Closure), Gesprächsverlauf mit dem Projektinhaber.

---

## Verification Evidence

Scope:
- Slice: [`slice-v1-adr-ollama-cloud-modelle`](../plan/planning/open/slice-v1-adr-ollama-cloud-modelle.md)
- IDs: [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md) (unverändert), [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md)
- Artefakte (Diff `origin/main..HEAD`, 8 Dateien): die ADR-Datei (neu, 170 Zeilen), der ADR-Index (+1), der Slice-Plan (neu), das Verzeichnis-README von `open/` (+1), die Roadmap (Offen-Zeile), die Welle-Datei (Zeile §4, Satz, Frage 2), zwei Review-Reports. Kein Go-Code, kein Docker-Pfad, keine Importregel.

### Sensors (selbst gefahren)

| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | Exit 0 (zwei Läufe; Log im Scratchpad); `d-check: 842 Datei(en) geprüft, 0 Befund(e)`, Coverage-Gate, `ablage-pruefung-test: 26 Faelle, 0 Fehlschlaege`; `git status` danach sauber |
| `make docs-check` | pass | 842 Dateien, 0 Befunde; mit diesem Report erneut gefahren: 843 Dateien, 0 Befunde |
| `make doc-commits RANGE="origin/main..HEAD"` | pass | 9 Commits aufgelöst, 0 Befunde |
| `make doc-immutable RANGE="origin/main..HEAD"` | pass | 0 Befunde. Hinweis: [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) ist in der Range neu; die Aussage „Accepted-ADR unverändert“ trägt hier nur für [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md) |
| `git diff origin/main..HEAD` über die Datei von [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md) | pass | 0 Zeilen (unverändert) |
| `git diff d9dd171e..HEAD` über die ADR-Datei | pass | leer (Accepted-Commit = HEAD zum Prüfzeitpunkt; DoD-Zusatz „inline Status“ erfüllt) |
| `make test-docker`, `make verify-depguard`, `make ci` | not run | nur Markdown berührt (siehe Diff-Liste); `verification.md` §Sensor-Auswahl: Mindest-Sensor `make docs-check`, normaler Closure-Sensor ebenso; `make gates` zusätzlich gefahren |

## DoD-Abgleich (Slice §2, Punkt für Punkt)

| # | DoD-Punkt | Urteil | Evidence |
|---|---|---|---|
| 1 | ADR (E1 bis E5) in MADR-Form, `Schärft:`, `Supersedes (teilweise)` auf [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md) mit Kopf-Abgleich; je Entscheidung Gegenoption, Messbeleg-Zeile (Rohdatei, Läufe, Pfad), Ungemessenes; Tabelle *Freigabe-Punkte* mit *Herkunft*; Fitness Function ohne Netz zu `ollama.com`; Re-Evaluierungs-Trigger; Index-Eintrag; Status `Accepted` durch den Projektinhaber | erfüllt | `Schärft:` und `Supersedes (teilweise)` im Kopf (Spanne auf die Vorgabe im Kontext von [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md) begrenzt, Nummerierungsregel „ohne Zusatz = diese ADR“); Entscheidungen 1 bis 11 mit Gegenoption (E1 → 2, 3, 6, 7, 9; E2 → 1, 4; E3 → 5, 6; E4 → 10, 11; E5 → Folgepflicht, 11); Messbeleg mit Rohdatei/Läufe/Pfad in 2 bis 7 und Tabelle Z1 bis Z12; Freigabe-Punkte a bis p mit Herkunft; Fitness Function (Golden/Unit/`docker`-Test, „kein Aufruf an `ollama.com`, kein Konto“); Trigger-Abschnitt; Index-Zeile 64 in `docs/plan/adr/README.md` (Status Accepted); Status `Accepted` im Kopf, Geschichte-Zeile |
| 2 | Zusage-Menge als eigene Tabelle (*zugesagt (bedingungsgleich)* / *nicht gemessen, nicht zugesagt*), Beleg je Zeile, Ausgang A/B benannt inkl. Antwort auf „trägt ein Lauf mit Konto …?“ | erfüllt | Tabelle Z1 bis Z12 mit Beleg- und Pfad-Spalte, Z13 als Produktversprechen mit noch zu lieferndem Test (ausgewiesen); Absatz *Nicht gemessen, nicht zugesagt*; Entscheidung 10 benennt Ausgang A, beantwortet die Frage und nennt die Rückfallebene B mit den zu streichenden Zeilen |
| 3 | Folgepflicht und Welle-Zeilen (Welle §4/§5, Roadmap) nachgezogen; `Entfällt`-Bedingung des Code-Slice benannt | teilweise (Planner-Anteil nicht Gegenstand) | Folgepflicht der ADR (Lastenheft inkl. [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), Spec-Stellen inkl. der Klammer in [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten), Handbuch, Code, Reihenfolge CR vor Code) vorhanden; „Code-Slice entfällt nicht“ in *Konsequenzen* benannt. Welle §4, Satz und §5 Frage 2 sowie Roadmap-Zeile auf „drei Dateien“ nachgezogen (Diff). **Offen, Planner-Closure:** die Welle-Zeile des Code-Slice führt weiter „entfällt, wenn die ADR keinen Host aufnimmt und keine Meldung ändert“ und die Zeile des Lastenheft-CR nennt [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) und [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) nicht; beides ist mit Accepted überholt, ausdrücklich nicht Prüfgegenstand |
| 4 | `make gates`, `make docs-check`, `doc-immutable`, `git diff` über die ADR leer, `doc-commits` | erfüllt | Sensors-Tabelle oben |
| 5 | Review durchgeführt, Report unter `docs/reviews/` | erfüllt | zwei Reports (Design und Nachlauf), kein Self-Review: Skill `reviewer.md`, anderer Kontext |
| 6 | Doku-Update: Welle, Roadmap, `open/README.md`, ADR-Index; kein öffentlicher Vertrag berührt | erfüllt (Welle/Roadmap-Stand s. 3) | Diff-Liste; `spec/`, `docs/user/`, `README*`, `CHANGELOG.md` nicht im Diff |
| 7 | Closure-Notiz mit Steering-Loop-Lerneintrag | offen | Slice §7 trägt den Platzhalter; Aufgabe der Planner-Closure, nicht Verifier |
| 8 | Beobachtungs-Register fortgeschrieben | offen | `observations/` nicht im Diff; Closure-Aufgabe |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | offen | alle neun Risiken in §6 stehen auf „offen bis zur Closure“; Closure-Aufgabe |
| 10 | Drei Paarungen (Anker · Folge-Slice · Register) | offen | Closure-Aufgabe |

Die Punkte 7 bis 10 sind Closure-Pflichten des Planners (Slice §4/§5) und vor `done/` zu erfüllen; sie sind hier kein Mangel des Baus. Die Bewegung nach `done/` ist ohne sie nicht zulässig.

## Folgen der Reviews: F-1 bis F-16 und N-1 bis N-6 am Text

| ID | Urteil | Beleg im Text (ADR @ HEAD) |
|---|---|---|
| F-1 | behoben | Suche in der ADR: kein Link/Zitat auf Welle, Roadmap, Slice-Datei; die einzigen `slice-`-Treffer sind Pfade nach `docs/messungen/` (Zeile 33); `docs-check` 0 |
| F-2, F-3 | behoben (F-3 mit Rest N-3, s. u.) | Regel in *Messung*; Z1 bis Z4, Z7 bis Z11 „Handvariante“; *Nicht gemessen* führt das erzeugte Sandbox-Profil ohne Egress |
| F-4 | behoben | Kontext (`relaxations` speist nur Ausweisung), O1 „neuer Eintragstyp mit neuem Rendering-Pfad“, *Negativ* |
| F-5 | behoben | Entscheidung 8 trennt Sandbox (`SecurityRelaxation`, sichtbar) und Standardprofil (entfällt); Fitness Function angeglichen |
| F-6 | behoben | Spalte *Herkunft* an das Muster von [`ADR-0035`](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) angeglichen; durch den Accepted-Commit auf „Projektinhaber, 2026-10-24 (Freigabe der Empfehlung)“ gesetzt |
| F-7 | behoben | „In dieser Lage trägt die erzeugte Zeile allein nicht“ mit Messlage |
| F-8 | behoben | Kopf: Spanne auf die Vorgabe im Kontext, Wortlaut-Fundort benannt, Nummerierungsregel |
| F-9 bis F-13 | behoben | `rohdaten-05-std-1-111134.txt` zitiert; „Messlauf war bereinigt“; „28 Aufrufe mit Exit 0 und Ausgabe, davon 14 `claude -p` mit Kurzantwort“; B-Streichliste vollständig; „abgeleitet“ gekennzeichnet |
| F-14 | behoben | Z13 als Produktversprechen getrennt |
| F-15 | behoben | Doc-Kommentar von `egressFeatureHosts` und [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) in der Folgepflicht (Rest N-2) |
| F-16 | behoben | Reihenfolge CR vor Code in der Folgepflicht; [`ADR-0012`](../plan/adr/0012-devcontainer-egress-firewall.md) in Entscheidung 3 im Text |
| N-1 | behoben | [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) in `Bezug`, `Schärft`, Entscheidung 8, Folgepflicht (Ausnahme 1 ergänzen), Fitness Function („nach dem CR“), Geschichte |
| N-2 | behoben | Folgepflicht zu [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) nennt die Klammer „(`runArgs` stammt allein aus `nestedRuntime` und Egress-Restriktion)“ als zu berichtigen |
| N-3 | behoben | Z6: „Auslieferung bis auf `devcontainer up` der CLI (Volume der CLI, Start von Hand)“, analog Z5 |
| N-4 | behoben | Suche nach `Planner`/`Welle`/`Roadmap` in der ADR: kein Treffer; die Planner-Pflichten stehen nicht mehr in der ADR |
| N-5 | behoben | Entscheidung 10: Zusage-Menge = Z1 bis Z12 (Messaussagen) und Z13 (Produktversprechen); Folgepflicht „Z1 bis Z13“ |
| N-6 | behoben | Marker „wird vor `Accepted` gefragt“ entfernt; Herkunft von i auf den Ist-Stand gesetzt; Geschichte-Zeile `Accepted` |

Die Reviewer-Findings sind keine Evidence (`verification.md`); die Urteile oben stützen sich auf den eigenen Textabgleich der ADR @ HEAD, nicht auf die Report-Aussagen.

## Entscheidungen des Projektinhabers (2026-10-24) gegen die Herkunftsspalte

| Entscheidung | ADR-Stelle | Herkunft in der Tabelle | Urteil |
|---|---|---|---|
| H2 (Host-Server), begrenzt auf das Gemessene | Punkt a, Entscheidung 4 | „Projektinhaber, 2026-10-24 (Begrenzung … Messung)“ | konform |
| Egress: `ollama.com` nur durch den Nutzer | Punkt b, Entscheidung 3 | Projektinhaber | konform |
| Geräteschlüssel: Handbuch-Hinweis, nicht werten, kein Code | Punkt c, Entscheidung 5 | Projektinhaber | konform |
| Kein Nachmessen | Punkt d, Entscheidungen 10, 11 | Projektinhaber | konform |
| O1 = K2, Katalogeintrag `ollama-host` | Punkt n, Tabelle O1 | Projektinhaber | konform |
| O2 = Warnung | Punkt o, Tabelle O2, Entscheidung 8 | Projektinhaber | konform |
| O3 bestätigt | Punkt p, Tabelle O3 | Projektinhaber | konform |
| j bis m | Punkte j, k, l, m | Projektinhaber | konform |
| e bis i (Freigabe der Empfehlung), Ausgang A | Punkte e bis i, Entscheidung 10 | „Projektinhaber, 2026-10-24 (Freigabe der Empfehlung)“ | konform |
| `Accepted` | Kopf, Geschichte, Index | Status `Accepted`, Datum 2026-10-24 | konform |

Grenze: Die Herkunftsangaben sind gegen den Auftrag des Aufrufers abgeglichen; die Entscheidungen selbst stehen als Gespräch außerhalb des Repos und sind nicht nachprüfbar.

## Plan-vs-Diff (Slice §3 gegen `git diff --stat origin/main..HEAD`)

| Plan-Zeile (§3) | Diff | Urteil |
|---|---|---|
| `docs/plan/adr/0036-<slug>.md` neu | `0036-cloud-modelle-von-ollama-im-devcontainer.md`, 170 Zeilen | erfüllt |
| `docs/plan/adr/README.md` update | +1 Zeile (Index) | erfüllt |
| Welle-Datei update (Zeile §4, Satz, §5 Frage 2) | 3 Stellen geändert | erfüllt |
| `open/README.md` update | +1 Zeile | erfüllt |
| `in-progress/roadmap.md` update | Offen-Zeile „drei Slices … bis zu vier weitere“ | erfüllt |
| Slice-Datei neu | vorhanden (276 Zeilen) | erfüllt |
| `docs/reviews/…-modelle.md`, `…-verify.md` neu | Review-Report vorhanden; dieser Report ist der Verify-Teil | erfüllt mit Abweichung unten |
| `observations/BEO-DOK/…` update (Closure) | nicht im Diff | offen (Closure) |
| Die Datei von [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md), Messprotokolle, `spec/`, Handbuch, READMEs, `CHANGELOG.md`, `internal/`, `cmd/`, `tools/` „unverändert“ | alle nicht im Diff | erfüllt (Datei von [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md): 0 Diff-Zeilen) |

Abweichungen (für den Closure-Abgleich nach `AGENTS.md` §3.10, nichts hier repariert):
1. Der Diff enthält einen zweiten Review-Report (`…-nachlauf.md`), den §3 nicht als eigene Zeile führt (Muster `…-modelle.md`, `…-verify.md`). Die Plan-Tabelle ist um diese Datei zu ergänzen.
2. Der Absatz „Abgleich mit `git diff --stat` (Stand des Plan-Commits)“ nennt vier Dateien und kündigt den Abgleich der Closure an; er ist zu wiederholen (jetzt acht Dateien plus dieser Report).
3. Slice §6 Fragen 1 bis 11 stehen als „gesammelt, nicht gestellt“; sie sind mit Accepted überholt (Zustandsmarker mit Ablaufpunkt, §3.10) und bei der Closure auf den Ist-Stand zu setzen.

## Traceability

| ID | Beleg |
|---|---|
| [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | `Bezug`/`Schärft` der ADR; Lastenheft nicht im Diff (Änderung nur per CR, Folgepflicht benannt) |
| [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md) | `Supersedes (teilweise)` mit Kopf-Abgleich; Datei unverändert (Diff 0 Zeilen) |
| [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) | Index-Eintrag, Accepted-Commit `d9dd171e`, `git diff` ab Commit leer |
| Commit-Regel (ID in Nachricht) | `make doc-commits` 0 Befunde |

## Carveouts

- Neu: none · Gelöst: none · Unverändert: `carveouts.md` nicht im Diff.

## Nicht ausgeführt

- `make test-docker`, `make verify-depguard`, `make govulncheck`, `make image-scan`: Diff trägt nur Markdown; kein Docker-, Importregel- oder Release-Pfad berührt.
- Replay/Golden (`harness/replay.md`): kein Generator-Output geändert.
- Rohdaten-Nachprüfung der Messzeilen über die Stichproben der Reviewer hinaus: Gegenstand der Messreviews; hier nur der Abgleich ADR gegen Reviews und Plan.

## Ergebnis

- **Bau gegen DoD:** DoD-Punkte 1, 2, 4, 5, 6 erfüllt; Punkt 3 für den ADR-Anteil erfüllt, der Welle-Zeilen-Nachzug nach Accepted ist Planner-Closure (ausdrücklich ausgenommen); Punkte 7 bis 10 sind Closure-Pflichten und offen.
- **Findings des Verifiers:** keine HIGH/MEDIUM. Drei Plan-Abweichungen (Plan-Tabelle ohne Nachlauf-Report, Abgleichsabsatz veraltet, §6-Fragenmarker) für die Closure, siehe oben.
- **Behauptung gegen Beleg:** alle Sensor-Aussagen (gates, docs-check, doc-commits, doc-immutable, Diff der Datei von [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md)) selbst gefahren und bestätigt; F-1 bis F-16 und N-1 bis N-6 am ADR-Text als behoben bestätigt.
- **Übergabe:** an den Planner (Closure: Slice §3/§6/§7, Welle-Zeilen mit Entfällt-Bedingung und CR-Bezug, Register, Risiko-Ausgänge, `git mv` nach `done/` mit Commit-Hash in der DoD-Zeile).

**Commit / Artefakt:** ADR-Commit `d9dd171e`; dieser Report in einem eigenen Commit.
