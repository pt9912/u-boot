# Verification-Report: slice-v1-forge-cli-quellen-messung — 2026-10-23

**Rolle:** Verifier (Modul 11, [`harness/verification.md`](../../harness/verification.md)), frischer Kontext, kein Self-Review, kein Reparieren. Frage: Wurde der Slice richtig gebaut, gemessen an Plan, DoD und Entscheidungslage?

**Gegenstand:** [`slice-v1-forge-cli-quellen-messung`](../plan/planning/done/slice-v1-forge-cli-quellen-messung.md) (Welle [`welle-forge-cli-gh-glab`](../plan/planning/done/welle-forge-cli-gh-glab.md)), Range `origin/main..fe4d9548` (elf Commits, nichts gepusht, kein Tag), [`ADR-0035`](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) (`Proposed`), Messprotokoll und Wortlaut-Entwurf unter [`docs/messungen/slice-v1-forge-cli-quellen-messung/`](../messungen/slice-v1-forge-cli-quellen-messung/README.md), Skripte `tools/messung/forge-cli/`. Verträge: [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer).

**Eingang gelesen:** Slice (alle Sektionen), ADR, Messprotokoll, Wortlaut-Entwurf, beide Review-Reports und `git show fe4d9548`, Welle-Datei §4 bis §5, `Dockerfile.tmpl` und `devcontainer_sandbox.go` (für Aussagen der ADR über den Code), Lastenheft [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte).

## Ergebnis

**Nicht schließbar.** Der Messungs-Anteil ist belastbar (Zahlen nachgerechnet, Token-Suche leer, Probe grün, Gates grün). Es fehlt, was kein Rollenlauf liefern kann: die Entscheidung des Projektinhabers zu den Punkten a bis g der ADR (DoD 3 und Closure-Trigger verlangen `Accepted`). Zwei Dinge muss der Planner zusätzlich entscheiden: den Haken von DoD 2 und den Umgang mit dem offenen Befund D7 vor `Accepted`. Kein Code-, Gate- oder Sicherheitsmangel.

## DoD Punkt für Punkt

| DoD-Punkt | Haken im Slice | Urteil | Beleg / Grund |
| --- | --- | --- | --- |
| 1 Messung je Werkzeug | [x] | **trägt**, mit benannter Einschränkung | Quellen, Lieferkette, Metadaten, Bau, Hosts im Protokoll A bis D; „Form, die `u-boot` ausliefert“ gibt es für einen Block nicht, das Protokoll führt Messvariante gegen Auslieferung offen (Abschnitt *Was Messvariante ist und was Auslieferung*); ein Block ohne Sandbox-Profil (ohne `USER root`) ist nicht gebaut, das Zielbild der ADR (Entscheidung 3) ist nicht gemessen und in „Nicht gemessen“ als Bau eines vom Generator erzeugten Dockerfiles geführt |
| 2 Umgang mit Tokens | [x] | **trägt nicht vollständig**, siehe V-2 | DoD fragt, ob ein Aufruf mit gesetzter Variable ohne `auth login` **gelingt**, und ob Git-Operationen einen Credential-Helper brauchen; gemessen ist: Variable wird verwendet, Anfrage geht raus, Antwort `401` (Platzhalter), Dateien, Verhalten ohne Variable; Credential-Helper nur gegen öffentliche Repositories |
| 3 Folge-ADR und Wortlaut | [ ] | **offen, richtig offen** | ADR steht auf `Proposed`; Index-Zeile vorhanden; Wortlaut-Entwurf liegt vor; Freigabe durch den Projektinhaber fehlt, siehe V-1 |
| 4 Gates | [x] | **trägt** | `make gates` Exit 0 (selbst gefahren, unten); `make docs-check` 0 Befunde (in `gates`); `make ablage-pruefung` 26 Fälle, 0 Fehlschläge; `make doc-commits RANGE="origin/main..HEAD"`: 11 Commits, 0 Befunde |
| 5 Review | [ ] | **erfüllt, Haken nachzuziehen** | Erst-Report und Nachlauf-Report liegen unter `docs/reviews/`; der Haken fehlt noch |
| 6 Doku-Update | [x] | **trägt** | ADR-Index (Zeile 0035), `docs/messungen/README.md`, Roadmap, Welle-Datei, `open/README.md` im Diff; `spec/`, `internal/`, `cmd/`, `docs/user/`, READMEs, CHANGELOG ohne Änderung (`git diff --stat` leer) |
| 7 bis 10 Closure-Notiz, Register, Risiken, Paarungen | [ ] | **offen, erwartet** | §7 „Wird bei der Closure gefüllt“; Register-Fortschreibung, Ausgänge der fünf Risiken aus §6 und die drei Paarungen sind erst mit der Closure zu schreiben |

## Selbst gefahrene Sensoren

| Sensor | Ergebnis | Beleg |
| --- | --- | --- |
| `make gates` | Exit 0 | endet mit `ablage-pruefung-test: 26 Faelle, 0 Fehlschlaege`; Arbeitsbaum vorher und nachher sauber |
| `make doc-commits RANGE="origin/main..HEAD"` | Exit 0 | 11 Commits, 737 Dateien, 0 Befunde |
| `tools/messung/forge-cli/redact-probe.sh` | Exit 0 | „20 Formen geschwaerzt, neutrale Zeile unveraendert“ |
| Token-Suche über `docs/messungen/`, `tools/messung/`, `docs/plan/`, `docs/reviews/` | kein Wert | Präfixe `ghp_ gho_ ghu_ ghs_ ghr_ github_pat_ glpat- glptt- gldt- glrt- sk-ant-` mit mindestens 8 Zeichen: Treffer nur in den Beispielzeilen der Probe (`PROBEWERT`) und in einem fremden Slice (Beschreibung eines Erkennungsausdrucks); `Authorization`/`Bearer`/`PRIVATE-TOKEN`/`JOB-TOKEN`/`X-Amz`/`Signature=`/`sig=`/URL-Zugangsdaten in Rohdaten: kein Treffer, in den Skripten nur `Authorization: Bearer $tok` (Variable) und die Regeln von `redact`; `NAME=wert` mit Token-Namen in Rohdaten: nur Variablenverweise; `$HOME`: kein Treffer |
| `make test-docker`, `make verify-depguard`, `make ci` | nicht gelaufen | kein Go-Code, keine Importregel, kein Image im Diff (`internal/`, `cmd/`, `spec/` ohne Änderung) |

## Zahlen selbst nachgerechnet (Protokoll gegen Rohdaten und Netz)

- **Bau C1 bis C10:** Image-Größen in `rohdaten-05.txt` und `-06.txt` gelesen; Differenzen nachgerechnet: `ghk` +60 037 684, `ghnk` +60 030 973, `ghdeb` +48 809 142, `glabdeb` +40 798 568, `ghtgz` +56 100 573, `glabtgz` +67 439 177, alle gleich dem Protokoll; Dauern 59,2 / 71,0 / 70,3 / 59,9 / 69,6 / 36,2 / 16,0 / 17,2 s, Exit-Codes der sechs Bauten über `source:` (alle 1, 0,8 bis 0,9 s), Zeilen `…tar.gz: OK` und `FAILED` der Negativprobe, Versionen 2.103.0 und 1.124.0 stimmen. Paketdifferenzen 1 / 0 / 11 / 0 / 10 / 1 stimmen.
- **Lieferkette A10 bis A14:** `gh` v2.103.0 `immutable: true`, 22 Assets, Digest und `checksums.txt` des `amd64`-Archivs gleich (`fdd768eb…`), 0 Signatur-Assets, Attestierungs-Schnittstelle HTTP 200 mit zwei Attestierungen; `VALIDSIG` mit `7F38BBB5…`; Debian 2.46.0-3 und 1.53.0-1+b3. A13: acht `gh`-Releases von 2026-07-31 bis 2026-10-21 (82 Tage, nachgerechnet) in `rohdaten-03.txt`; die acht `glab`-Releases von 2026-09-08 bis 2026-10-20 **live gegen die GitLab-API** wiederholt, gleiche Daten.
- **Pins der ADR (Entscheidung 1), live am 2026-10-23:** `checksums.txt` von `gh` v2.103.0 (GitHub) und `glab` v1.124.0 (Paketregister): `amd64` und `arm64` beider Werkzeuge stimmen mit den vier Werten der ADR überein; `gh` v2.103.0 `immutable` weiterhin wahr.
- **Hosts D1 bis D3:** `gesamtdauer` der vier `glab`-1.124.0-Varianten: `discover` 19,6 / 19,5 / 18,2 / 18,8 s (Spanne 18,2 bis 19,6), `restricted` 127,8 / 130,5 / 129,7 / 127,3 s (Spanne 127,3 bis 130,5); `glabdeb` 15,3 gegen 8,5 s; zehn Zeilen `server misbehaving` in `glabsk-restricted`; DNS-Namen je Variante: nur `api.github.com`, `github.com` beziehungsweise `gitlab.com`.
- **Token E1 bis E8:** Dateigrößen (`device-id` 36, `aliases.yml` 32, `config.lock` 0, `config.yml` 7 534 bis 7 559 Byte, `.gitconfig` 181), `platzhalter-treffer=0`, Exit-Codes (`gh api` ohne Variable 4, mit Platzhalter 1 mit `HTTP 401`, `glab api` ohne Variable 0, mit Platzhalter 1), `git_protocol=https` bei `gh`, `gh auth status` mit Platzhalter Exit 1 (2.103.0) gegen Exit 0 (2.46.0): stimmen.
- **Nicht nachgerechnet:** A1/A2 (Feature-Index, Digests der Nicht-`gh`-Features), A6/A8 über die Größen hinaus, C3 (Link-Details), E5, E6, E7, E9 und die Läufe `ghx`/`glabp`/`glabz` in D4 im Einzelnen. Eine Abweichung ist mir in dem Gelesenen nicht aufgefallen; vollständig abgeglichen ist es nicht.

## Findings

### V-1 — Blocker für die Schließung — ADR-0035 `Proposed` gegen DoD „Accepted“

Slice-DoD 3 und Closure-Trigger (§5) verlangen `Accepted`. Die Entscheidung liegt ausdrücklich beim Projektinhaber: Quelle, die vom Muster abweicht (Dockerfile-Block, Drittanbieter, externe Quelle) entscheidet er (DoD 3 Schlusssatz), die Punkte a bis g stehen in der ADR als „Empfehlung des Implementers (keine Entscheidung)“. Im Repo gibt es **keine** Spur einer Entscheidung: Welle-Datei §5 Frage 4 sagt „liegen dem Projektinhaber vor“, die Fragen 1 bis 3 sind ebenfalls unbeantwortet.

Was für die Schließung fehlt (in dieser Reihenfolge):

1. Entscheidung des Projektinhabers zu a (Quelle `gh`), b (Quelle `glab`), c (Token-Variablen), d (beide Profile, Template-Änderung), e (Telemetrie), f (Nachmessung mit Token oder Zusage nur für Gemessenes), g (ein oder zwei Code-Slices). Der Wortlaut-Entwurf hängt an a bis d; er braucht dieselbe Freigabe (Teil von DoD 3).
2. Entscheidung des Planners über D7 und über den Haken von DoD 2 (V-2, V-3), **bevor** die ADR `Accepted` wird; danach ist sie nach [`AGENTS.md`](../../AGENTS.md) §3.5 nicht mehr umzuschreiben.
3. Danach (Arbeit des Implementers oder Architekten, nicht dieses Reports): ADR-Text von „Vorschlag“ und „Empfehlung des Implementers“ auf die getroffene Entscheidung umschreiben (Kopf `Status`, Abschnitt *Entscheidung*, Tabelle *Freigabe-Punkte*, *Geschichte*), Index-Zeile in `docs/plan/adr/README.md` (Status `Accepted`, „Vorschlag“ streichen), Wortlaut-Entwurf an die Wahl angleichen (Voraussetzungs-Absatz), Welle-Datei §5 Frage 4 und Slice §1/§2 nachziehen.
4. Die Closure-Pflichten aus §2 (Häkchen 5, Closure-Notiz mit Steering-Loop-Eintrag, Register, Risiko-Ausgänge, Paarungen) und der Verifier-Lauf über den dann geänderten Stand.

Der Verifier setzt `Accepted` nicht und schlägt keine Wahl vor; die Empfehlungen der ADR sind in sich stimmig (geprüft unter *Zahlen*), das ersetzt die Entscheidung nicht.

### V-2 — Entscheidung für den Planner — Haken von DoD 2 („Umgang mit Tokens gemessen“)

Die DoD-Frage hat fünf Teile; Stand der Messung:

| Teilfrage | Stand | Beleg |
| --- | --- | --- |
| welche Variable trägt die Anmeldung | beantwortet | E1 |
| gelingt ein Aufruf mit gesetzter Variable ohne `auth login` | **nur „wird versucht, `401` mit Platzhalter“**; das Gelingen ist nicht gemessen | E2, F |
| ob und wo das Werkzeug eine Datei schreibt | beantwortet (ohne Anmeldung) | E4 |
| Ende ohne Variable | beantwortet | E3 |
| Credential-Helper für Git-Operationen | **nur öffentliche Repositories**; privates Repository und GitLab nicht gemessen | E5, F |

Der DoD-Text erlaubt „ein Aufruf ohne Token, das Protokoll nennt welches“, und das Protokoll nennt es (Platzhalter, kein Token angelegt). Formal ist das gedeckt; inhaltlich ist die Frage „gelingt“ nicht beantwortet, sondern ausdrücklich ausgeklammert. Ehrlich ist es im Protokoll, in der ADR und im Wortlaut (S6 nimmt Aufrufe mit gültigem Token aus der Zusage), also keine Überzusage.

**Urteil:** Der Haken trägt nur, wenn die DoD-Formulierung angepasst wird. Zwei saubere Wege, die Wahl liegt beim Planner und hängt an Punkt f des Projektinhabers:

- (a) Zusage nur für das Gemessene (Empfehlung der ADR): DoD 2 im Plan auf „Umgang mit Tokens **ohne gültiges Token** gemessen; Erfolgspfad und Credential-Helper gegen private Repositories ausdrücklich nicht gemessen (Protokoll F)“ ändern, mit Freigabe des Projektinhabers; der Haken bleibt dann.
- (b) Nachmessung mit kurzlebigem Token vor der Schließung: Haken bleibt offen, die Nachmessung ist ein eigener Auftrag, und das Protokoll braucht dafür das Lesen der Rohdaten (siehe `redact`-Grenze unten).

Ohne eine der beiden Entscheidungen bleibt der Haken ein „Haken über dem Wortlaut der DoD“; er ist bis dahin als nicht belegt zu führen.

### V-3 — MEDIUM — D7 (`glab` lädt womöglich weitere Binaries) ist vor `Accepted` zu entscheiden; Wortlaut-Entwurf trägt es nicht

Befund im Dokument: ADR Entscheidung 4 und 6 führen D7 als „offene Randbedingung“ und schreiben „der Projektinhaber sieht ihn bei Punkt f“. Punkt f der Freigabe-Tabelle fragt aber nur nach der Nachmessung mit Token; D7 steht in keiner Zeile der Tabelle, also nicht als eigener Freigabe-Punkt. Der Wortlaut-Entwurf nennt D7 nur in der Spalte „Nicht gemessen“ von S4; der Satz S4 (Hosts von `glab` = `gitlab.com`) und die Stelle 3 („[`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte): keine Änderung“, Auslegung wie `ollama`) sagen es unbedingt. Folge: Wer nur den Wortlaut freigibt, sieht D7 nicht.

Eigene Lesung am 2026-10-23 (ein Lauf, `glab` 1.124.0 `amd64` aus dem Archiv mit SHA-256-Prüfung im Basis-Image, nichts angemeldet, kein Download ausgelöst): alle zehn Schlüssel `duo_cli_*`/`orbit_cli_*` (Auto-Run, Auto-Download, Pfad, Version, Prüfsumme) haben nach dem ersten Start **keinen Wert**; `glab --help` führt `duo` und `orbit` als eigene Unterbefehle; der Hilfetext von `glab orbit` sagt, `glab` lade, prüfe und aktualisiere das Orbit-Binary **bei der ersten Benutzung**, die Konfiguration kenne eine Bestätigungsabfrage (`orbit_cli_auto_download`: „Skip the download confirmation prompt“), `glab duo cli` „installs and runs“ das Duo-CLI-Binary. Das spricht dafür, dass D7 ein **ausdrücklicher Unterbefehl mit Bestätigung** ist, kein Hintergrundladen. Nicht gemessen und nicht belegt: von welchem Host das Binary kommt, ob die Abfrage in einer Sitzung ohne Terminal ausbleibt, was `auto_download` per Umgebungsvariable bewirkt. Das ist Hilfetext, keine Messung eines Downloads.

Optionen: (a) D7 als eigener Freigabe-Punkt (h) in die ADR und eine Zeile „Nicht zugesagt: `glab duo` und `glab orbit` und die Hosts ihres Nachladens“ in S6, oder (b) kleine Messung des Download-Hosts im Messungs-Slice (hinter `dnsmasq`, wie D1 bis D3) vor `Accepted`. Beides ist Planner- und Projektinhaber-Sache; die ADR ist bis zur Annahme noch änderbar.

### V-4 — LOW — `arm64`-Pin von `glab` ist in keiner Rohdatei als gelesener Wert belegt

Das Protokoll (zu C) und die ADR sagen, der `arm64`-Pin sei „aus den Prüfsummen-Dateien gelesen“. `0716395a…` steht in den Rohdaten nur im Dockerfile-Block und im Bau-Log von `glabtgz`, nicht als Ausgabe von `inner-03.sh` (die Rohdaten drucken nur den `amd64`-Wert von `glab`; beim `gh` steht der `arm64`-Digest in `rohdaten-03.txt` Zeile 10). Der Wert selbst ist richtig (live gegen das Paketregister nachgelesen, 2026-10-23). Beleg-Lücke, kein Wertfehler; mit der ADR-Annahme ist ein Satz „live nachgelesen“ oder die Aufnahme der Zeile in `inner-03.sh` der saubere Weg.

### V-5 — LOW — Plan-Tabelle (§3) nennt `docs/reviews/` nicht

Der Diff trägt zwei Dateien unter `docs/reviews/` (Review und Nachlauf, jetzt auch dieser Report), die Plan-Tabelle führt kein Verzeichnis dafür ([`AGENTS.md`](../../AGENTS.md) §3.10: jede berührte Datei, auch Dokumentation). Zeile nachtragen.

### V-6 — INFO — Fortbestehende Grenzen und offene Reviewer-Punkte

- **`redact` formgebunden (Reviewer I-1):** Die Grenze steht im Protokoll; für die Nachmessung mit echtem Token (Variante b) gilt: Rohdatei lesen, nicht nur suchen. Eigene Wiederholung der Suche ändert daran nichts.
- **Reviewer I-2 (S6 und Plattformen)** und **F-6 (Heredoc-Commits, Host-Python)** liegen unverändert vor; ich habe die Commit-Texte nicht erneut auf Heredoc-Reste geprüft. Zur Klarheit: Ich habe selbst einmal `python3` auf dem Host für eine Rechenprobe aufgerufen (kein Beleg, die Differenzen oben sind in `bash` wiederholt); Nachweise stammen aus Docker und `make`.
- **Risiken §6, Vorschläge für die Ausgänge** (Entscheidung beim Implementer/Planner bei der Closure): 1 weiter offen, Vorsorge getragen (Läufe-Spalte, Nicht-gemessen-Listen je Satz); 2 eingetreten und benannt (Handvarianten `*k`, `*tgz`, `*deb`; kein Block im Generator); 3 entfallen (kein Token, Suche leer), mit Hinweis auf die `redact`-Grenze; 4 zur Hälfte eingetreten (Befund B1 zum `source:`-Weg, Empfehlung ohne externe Quelle); 5 entfallen (kalte Bauten, Aufräumen `objekte-mit-praefix-nachher=0`).

## Traceability

| ID | Beleg |
| --- | --- |
| [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) | Protokoll A bis C, ADR Entscheidung 1 bis 3; Wortlaut Stelle 1; Katalog im Code unverändert gelesen |
| [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) | A5 bis A12, C10/C11; Zusage ab D7 offen (V-3) |
| [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) | E1 bis E9, B5; Erfolgspfad nicht gemessen (V-2) |
| [`ADR-0035`](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) | `Proposed`; Index-Zeile; `Schärft:` nach oben; Aussagen über Code (Template, `egressFeatureHosts`, Remote-Env) gelesen und zutreffend |

## Carveouts

Neu: keine. Gelöst: keine. Unverändert: alle. Kein Gate gelockert, keine Schwelle berührt.

## Nicht ausgeführt

`make test-docker`, `make verify-depguard`, `make ci`: kein Go-Code, keine Importregel, kein Image im Diff. Keine Messung mit gültigem Token (kein Token bereitgestellt). Kein Download-Lauf zu D7.

## Commit / Artefakt

Stand der Prüfung: `fe4d9548`; dieser Report ist der nächste Commit darauf.
