# Review-Report (Nachlauf): slice-v1-forge-cli-quellen-messung, Fix-Runde 995c5dbe — 2026-10-23

**Rolle:** Reviewer (Modul 10, [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md), [`harness/review.md`](../../harness/review.md)), frischer Kontext. Prüfung gegen Plan, Entscheidungen und Hard Rules, nicht gegen die DoD (Verifier).

**Gegenstand:** Commit `995c5dbe` (7 Dateien) gegen den Erst-Report [`2026-10-23-slice-v1-forge-cli-quellen-messung.md`](2026-10-23-slice-v1-forge-cli-quellen-messung.md): Slice [`slice-v1-forge-cli-quellen-messung`](../plan/planning/done/slice-v1-forge-cli-quellen-messung.md), [`ADR-0035`](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) (`Proposed`), Messprotokoll und Wortlaut-Entwurf unter [`docs/messungen/slice-v1-forge-cli-quellen-messung/`](../messungen/slice-v1-forge-cli-quellen-messung/README.md), Skripte unter `tools/messung/forge-cli/`.

**Eingang gelesen:** `git show 995c5dbe` vollständig, Erst-Report (F-1 bis F-9), Messprotokoll (geänderte Abschnitte und D-Zeilen), die ADR (geänderte Zeilen und Restsuche), Wortlaut-Entwurf, Plan-Tabelle des Slice, `AGENTS.md` §3.7, §3.9, §3.10.

## Ergebnis

Kein HIGH, kein MEDIUM. Die sieben bearbeiteten Findings (F-1 bis F-5, F-7, F-9) sind behoben; zwei neue LOW und zwei INFO. F-6 (Prozessabweichung) und F-8 (DoD-Haken) lagen nicht in der Runde und bleiben unverändert INFO/LOW wie im Erst-Report.

## Nachprüfung der Findings

| Finding | Stand | Beleg |
|---|---|---|
| F-1 (MEDIUM, `redact`) | behoben | `tools/messung/forge-cli/redact-probe.sh` selbst gelaufen (Docker, über `lib.sh`): `redact-probe: 20 Formen geschwaerzt, neutrale Zeile unveraendert`, Exit 0. Gegenprobe rot: `lib.sh` aus `995c5dbe~` in ein Scratch-Verzeichnis, dieselbe Probe: Exit 1, genau 14 Zeilen `DURCHGELASSEN` (`authorization: token`, drei Token-Header, `GITLAB_TOKEN=`, `GH_TOKEN=`, `export GITHUB_TOKEN=`, drei JSON-Felder, drei URL-Formen) — deckt sich mit „14 von 20“ im Protokoll. Neues `redact` über alle `rohdaten-*` (10 402 Zeilen): `cmp` je Datei, keine Änderung (Aussage „ändert keine Zeile“ stimmt). |
| F-2 (Läufe-Spalte) | behoben | Spalte A1 bis B6 überall `1`, Definition nennt „deshalb überall auf 1“; „früherer“ zu „frühere“ berichtigt. |
| F-3 (Einschränkungen) | behoben | ADR-Tabelle (A3 nur zwei Mechanismen, A12 Suchmuster nicht abschließend, C4 bis C6 als Handvariante, B1 baut nicht), Entscheidung 6 und Konsequenz „Negativ“ tragen den Vorbehalt; Kommentar in `block-glabtgz.dockerfile` ebenso. Die ADR bleibt `Proposed`, §3.5 nicht berührt. |
| F-4 (Messwerte im Wortlaut) | behoben | S6 im Absatz ohne Docker-/Werkzeugversionen; Messstand in der Tabellenzeile. |
| F-5 (Zeiger ins Leere) | behoben | `lib.sh` zeigt auf „Abschnitt Lesehilfen“; der Abschnitt existiert (README Zeile 181). Suche nach „Wiederholung“ im Baum: kein weiterer Zeiger (§3.9 b). |
| F-7 (`glab`-Konfiguration) | behoben, als D7 geführt | D7 gegen `rohdaten-07-glabsk-discover.txt` gelesen: Schlüsselliste (Zeile 229) enthält `duo_cli_*` (sechs), `orbit_cli_*` (sechs) und `notify_skill_updates`; `glab config get duo_cli_auto_download`/`duo_cli_auto_run` geben in der Zeile 166 `duo_cli_auto_download=duo_cli_auto_run=…` ohne Wert aus; `orbit_cli_*` nicht abgefragt — alles wie in D7 formuliert. |
| F-9 (D3-Spanne) | behoben; **Implementer-Wert gilt** | Rohdaten `gesamtdauer`, Modus `discover`, `glab` 1.124.0: `glabtgz` 19,6 · `glabsk` 19,5 · `glabzk` 18,2 · `glabpk` 18,8, also **18,2 bis 19,6 s**; hinter der Restriktion 127,3 (`glabpk`) · 127,8 · 129,7 · 130,5, also **127,3 bis 130,5 s**; `glabdeb` 8,5 gegen 15,3 s. Mein Wert „19,8“ im Erst-Report war falsch (Maximum 19,6; 19,8 steht in keiner Rohdatei); die Berichtigung durch den Implementer ist richtig. |

## Prüfung nach Schwerpunkten

1. **§3.10 (Aussagen ohne Messung).** Jede neue Zahl in D3, D7 und der Probe-Beschreibung ist an Rohdaten oder am Probe-Lauf nachvollzogen (siehe Tabelle). Die Token-Suche ist wiederholt (Präfixe inklusive `glptt-`, `gldt-`, `glrt-`, `sk-ant-` mit mindestens 8 Zeichen; `Authorization`, `Bearer`, `PRIVATE-TOKEN`, `JOB-TOKEN`, `X-Amz`, `Signature=`, `sig=`, URL-Zugangsdaten, JSON-Felder mit Token-Namen, Namen des Home-Verzeichnisses, `NAME=wert` mit Token-Namen) über `docs/messungen/`, `tools/messung/`, `docs/plan/`: in den Rohdateien kein Wert (Treffer sind `"GIT_TOKEN": "${localEnv:GIT_TOKEN}"`, also Variablenverweise, und Prüfsummen in `rohdaten-03.txt`); in `docs/plan/` nur die Beschreibung des Erkennungsausdrucks in einem fremden Slice.
2. **ADR gegen Protokoll.** Die Einschränkungen in der ADR widersprechen dem Protokoll nicht mehr (A3, A12, C4 bis C6, B1). Restsuche „Signatur“ in ADR, Slice, Skripten: eine unbedingte Fassung bleibt im Wortlaut-Entwurf (L-1).
3. **Kommentare (§3.7/§3.9).** Der neue Kommentar in `lib.sh` und der in `redact-probe.sh` sind Indikativ und tragen Zusage und Grenze („Zusage: die Formen der Liste; nicht zugesagt: …“); kein „früher/ohne X wäre“, kein Konjunktiv über Verworfenes. Der Kommentar in `block-glabtgz.dockerfile` nennt die Grenze des Suchmusters. Keiner verweist auf ein geändertes Symbol (Suche nach altem Zeiger leer).
4. **Plan-Tabelle (§3.10).** Die 7 Dateien des Commits sind gedeckt: `tools/messung/forge-cli/` (Zeile nennt jetzt `redact-probe.sh`; `lib.sh`, `block-glabtgz.dockerfile` unter „Messskripte … vier `block-*.dockerfile`“), der Messordner (`README.md`, `wortlaut-entwurf.md`), die ADR, „diese Datei“. Kein Eintrag ohne Diff in der Runde, keine Zeile ohne Plan-Eintrag.
5. **Sensors.** `make gates` grün (Exit 0, `ablage-pruefung-test: 26 Faelle, 0 Fehlschlaege`). `make test-docker` und `make verify-depguard` nicht gelaufen: kein Go-Code, keine Importregel berührt.

## Findings

### L-1 — LOW — Wortlaut-Entwurf führt „Die Archive sind nicht signiert“ unbedingt, die ADR nur „nicht gefunden“

- `kategorie`: LOW · `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.10 (eine Aussage über Verhalten nennt nur, was gemessen wurde) · `klasse`: zusage-weiter-als-messung
- `pfad`: `wortlaut-entwurf.md` im Messordner · „Die Archive sind nicht signiert; für `gh` veröffentlicht das Projekt eine GitHub-Attestierung“ (Satz S3)
- `befund`: Die Fix-Runde setzt in der ADR (Entscheidung 6, Konsequenzen) und im Block-Kommentar „nach dem Suchmuster der Messung keine Signatur gefunden“; der Wortlaut-Entwurf, den der CR-Autor übernimmt, behält die Fassung ohne Vorbehalt. Die Begründungszeile S3 in der Tabelle trägt den Vorbehalt („ob `glab` außerhalb des Suchmusters signiert“), der Satz selbst nicht.
- `risiko`: Der CR-Autor übernimmt den Satz; das Lastenheft sagt dann über `glab` mehr, als A12 trägt, und die ADR und der Entwurf weichen in derselben Aussage voneinander ab.
- `verifizierbar`: Review-only.

### L-2 — LOW — D7 steht im Protokoll, nicht in der Entscheidungsgrundlage der ADR

- `kategorie`: LOW · `quelle`: [`ADR-0035`](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) Entscheidung 4 (Hosts „gemessen D1 bis D3“) · `klasse`: befund-nicht-in-entscheidungsgrundlage
- `pfad`: `docs/plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md` · „Hosts nach gültiger Anmeldung, Telemetrie und Update-Prüfung sind **nicht** gemessen“
- `befund`: D7 hält offen, ob `glab` weitere Binaries von weiteren Hosts lädt und die Host-Liste samt Zusage zu [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) zu knapp wäre. Die Suche nach `duo`, `orbit`, `D7` in der ADR ist leer; der Wortlaut-Entwurf kennt nur „Telemetrie- und Update-Hosts“ allgemein.
- `risiko`: Der Projektinhaber gibt die ADR frei und liest die Host-Liste von `glab` als vollständig, ohne den offenen Befund zu sehen.
- `verifizierbar`: Review-only.

### I-1 — INFO — `redact` bleibt formgebunden; die Grenze ist benannt

Außerhalb der gepinnten Formen bleiben Werte stehen, z. B. `password: wert`, `GITLAB_TOKEN = wert` (Leerzeichen um das Gleichheitszeichen), `Cookie: session=wert`, `invalid token wert` (eigener Lauf). Das Protokoll sagt es („erkennt Formen, nicht Werte“, Rohdatei vor dem Commit lesen); keine Änderungspflicht, aber der Nachmessung mit echtem Token (Punkt f der ADR) sollte die Grenze in die Anweisung eingehen. Zuständigkeit: Implementer der Folgemessung.

### I-2 — INFO — S6 im Wortlaut: Plattformgrenzen nur noch in der Tabellenzeile

Nach der Änderung nennt der Satz S6 `arm64`, macOS, Windows, Podman nicht mehr; die Tabellenzeile sagt, sie seien „im Satz nicht zugesagt, aber auch nicht ausdrücklich ausgeschlossen“, während S1 bis S5 plattformunabhängig formuliert sind. Das ist als Zielbild-Text vertretbar (die ADR nennt `arm64` ungebaut), aber die Tabellenaussage „im Satz nicht zugesagt“ ist eine Lesart, kein Wortlaut. Entscheidung beim CR-Autor und Projektinhaber (Frage 4 der Welle).

## Negativbefund (geprüft, ohne Befund)

- `tools/messung/forge-cli/redact-probe.sh`, `lib.sh` (`redact`, Kommentare) · `block-glabtgz.dockerfile` · Rohdaten auf Zugangsdaten · Läufe-Spalte A1 bis B6 · D3 und D7 gegen Rohdaten · Plan-Tabelle gegen die 7 Dateien · Kommentare auf §3.7/§3.9 · die ADR (bleibt `Proposed`, keine Immutabilitätsverletzung, `Schärft:`-Richtung unverändert).

## Nicht geprüft

- `make test-docker`, `make verify-depguard`: kein Go-Code und keine Importregel im Diff.
- Rohdaten-Abgleich aller übrigen Zahlen des Protokolls (A bis C, E): nicht Gegenstand der Runde; Stichproben des Erst-Reports gelten.
- F-6 (Heredoc-Commits, Host-`python3`) und F-8 (DoD-Haken „Umgang mit Tokens“): nicht Gegenstand der Fix-Runde, Stand unverändert.
