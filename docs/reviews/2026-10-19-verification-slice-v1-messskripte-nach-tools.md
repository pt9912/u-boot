# Verification-Report: slice-v1-messskripte-nach-tools — 2026-10-19

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“, gegen Plan, DoD und Entscheidungen). Kein Review, keine Validierung, keine Reparatur.

**Gegenstand:** Slice [`slice-v1-messskripte-nach-tools`](../plan/planning/done/slice-v1-messskripte-nach-tools.md) (`in-progress/`), Commit-Range `e952311..HEAD` (HEAD `77dde24`, lokal). Implementierung: `4e3a4d4` (Move Werkzeug), `b541202` (Move Protokolle), `98f76a4` (Pfad-Berichtigung), `79c749f` (ADR-Berichtigung), `ae39ec2` (Nachzug), `42bb8b2` (Entscheidung (c)), `6bf964b` (Hinweis a1), `7f0bdc4` (Gate a3), `2f21a06` (Plan-Tabelle). Review-Report `ccb4e31`. Fix-Runde: `afce363` (F-2), `c9dc3b3` (F-3), `51173b8` (F-4), `3c094c3` (F-5), `7d27441` (F-6), `3e36489` (F-1, Beobachtung), `589464d` (F-1, F-7, F-8, Plan), `77dde24` (Plan-Tabelle nach der Fix-Runde). `git diff --stat -M e952311..HEAD`: 86 Dateien, davon 52 reine Renames, 3 Renames mit Inhaltsänderung (zwei Protokoll-READMEs, `06-option-a.sh`), 8 neue Dateien, 23 geänderte.

**Eingang:** Implementer-Zusage (DoD, Plan-Tabelle, Fix-Runde), Review [`2026-10-19-slice-v1-messskripte-nach-tools`](2026-10-19-slice-v1-messskripte-nach-tools.md) (kein HIGH; F-1, F-2 MEDIUM; F-3 bis F-7 LOW; F-8 bis F-12 INFO), Vorgaben des Projektinhabers vom 2026-10-19: (1) im Planning-Baum nur Dokumente, unter `observations/` nur Beobachtungen; (2) Messprotokolle unter `docs/messungen/<slice>/` nur `.md`, `.txt`, `.log`; (3) Skripte, Dockerfiles, Sonde unter `tools/messung/<name>/`; (4) entschieden (a) a1+a3, (b) b3, (c) Vorschlag übernommen, nicht widersprochen.

**Gesamturteil: GO mit Auflagen.** Der Umzug ist rein (55 Renames `R100` in den Move-Commits, Rohdaten blob-gleich), die ADR-Änderung besteht aus genau der Messungs-Zeile und je einer Zeile in *Geschichte*, das neue Gate ist eine echte Positivliste und färbt sich an meinen eigenen Probedateien rot, `make gates` und die Doku-Sensoren sind grün, Lastenheft, Spezifikation, Rohdaten, Workflows, `.d-check.yml` und `internal/**` sind unverändert. Eine DoD-Verletzung (Zusage ohne Beleg) liegt nicht vor. Die Auflagen sind Closure-Arbeit des Planners, eine Wortlaut-Angleichung der DoD (Abschnitt 3, Punkt 7) und zwei Entscheidungen des Projektinhabers (Gate-Lücke `doc-immutable`, Archiv-Hinweis). **Nicht gelaufen:** kein Messskript an seinem neuen Ort (Abschnitt 2).

## 1. Sensoren (selbst gefahren, nur `make`)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | pass, **Exit 0** | Kette `docs-check`, `baseline-verify`, `ablage-pruefung`, `ablage-pruefung-test`, `record-gates` (`make -n gates`); `coverage-gate: OK — coverage 91.40% meets threshold 90%` |
| `make ablage-pruefung` | pass | `ablage-pruefung: OK`; Wächter-Test „26 Faelle, 0 Fehlschlaege“ |
| `make docs-check` | pass | 509 Dateien, 0 Befunde (vor Anlage dieses Reports) |
| `make doc-trace` | pass | „83 Anforderung(en), 0 Waise(n).“ |
| `make doc-planning` | pass | 509 Dateien, 0 Befunde |
| `find docs/plan/planning/observations -type f`, abzüglich `observations/README.md`, `observation.md`, `state.md`, `evidence/<id>.md` | **leer** | nur Beobachtungsform |
| `find docs/plan/planning -type f ! -name '*.md'` | genau die drei `.gitkeep` | `open/`, `next/`, `done/` |
| `find docs -type f` auf `*.sh`, `Dockerfile*`, `*.pl` | **leer** | `git ls-files docs` ohne `md`/`txt`/`log`: nur `.gitkeep` (ADR, Carveouts, `open/`, `next/`, `done/`, `reviews/`) |
| Probe am Gate (lokal, nicht committet) | siehe Abschnitt 4 | rot an `docs/messungen/x/probe.sh` und `docs/messungen/RUN.SH`; Arbeitsbaum danach sauber |
| `make doc-immutable STAGED=1` mit Probe-Änderung an zwei Überschriften | Abschnittsform rot, Inline-Form grün | Abschnitt 5 (F-1 unabhängig bestätigt); Probe zurückgenommen |
| Suche im Arbeitsbaum ohne `.git` nach den sechs vom Projektinhaber benannten Mustern fremder Docker-Objekte und Registry-Pfade | **0 Treffer** | die Historie in `.git` ist nicht Gegenstand |
| `bash -n` auf alle `tools/messung/**/*.sh` | pass | kein Fehler |
| Aufruf ohne Argumente an `06-option-a.sh`, `01-projekt-erzeugen.sh`, `run.sh` | erwartete Fehlmeldung | `arbeitsverzeichnis` beziehungsweise `repo-wurzel` (`${1:?…}`) |
| `ai-harness-init .` in einem Klon im Scratchpad | Exit 0 | `ablage-pruefung.sh` und `-test.sh` bleiben stehen (MR-012-Aussage bestätigt); das Werkzeug schreibt `harness/mk/traeger.mk` neu (Review F-12 bestätigt) |
| `git diff --stat -M e952311..HEAD -- spec .github internal .d-check.yml Makefile Dockerfile .golangci.yml tools/harness/selbstpruefung.sh` | **leer** | kein Lastenheft, keine Spezifikation, kein Workflow, kein Code, `selbstpruefung.sh` unverändert |
| `make test-docker`, `make verify-depguard`, `make doc-immutable RANGE="e952311..HEAD"` als Beleg | nicht als Beleg gezählt | kein Docker-Pfad und keine Importregel im Diff; `doc-immutable` meldet 0 Treffer, weil das Gate die zwei ADRs nicht schützt (Abschnitt 5), der ADR-Diff ist gelesen |
| Messläufe der Skripte an ihrem neuen Ort | **nicht gelaufen** | brauchen den Messhost (Docker-Socket, Netz, Anmeldung) und verändern dessen Docker-Bestand; kein Auftrag zum Messlauf, der Slice schließt ihn aus (§1). Die Aussage „sie arbeiten vom neuen Ort aus“ beruht auf dem Lesen der Ortsannahmen (Abschnitt 2), `bash -n` und dem Aufruf ohne Argumente |

## 2. Moves rein, Pfade, Ortsannahmen

- `git diff -M --numstat 4e3a4d4^ b541202`: keine Zeile mit Änderungen; 55 Renames, alle `R100` (30 Werkzeug-, 25 Protokoll-Dateien), keine Moduswechsel. Die Rohdaten (23) sind damit blob-gleich; bis `HEAD` ändert sich keine Rohdatei (nur die zwei READMEs und `06-option-a.sh` tragen `R084`/`R086`/`R089`, alle drei sind Inhalt der Nachzug-Commits).
- Pfad-Berichtigung `98f76a4` und F-4-Fix `51173b8` zusammen gelesen: `06-option-a.sh` bekommt `out=${OUT:?…}` **vor** `rm -rf "$w"` (Abbruch vor jedem Lauf, wenn `OUT` fehlt) und kopiert mit `cp "$w/build-a.log" "$out/rohdaten-06-build-a.log"` ohne `2>/dev/null || true` (ein falsches `OUT` wird gemeldet). Der Aufruf im Kopf und das Protokoll-README nennen `OUT` als Pflicht. Minimal und korrekt; F-4 behoben.
- Ortsannahmen: `grep` nach `dirname "$0"`, `BASH_SOURCE`, `$here` über `tools/messung`: Alle Skripte, die ihren Ort ableiten, suchen Nachbarn (`Dockerfile.*`, `inner*.sh`, `probe.pl`, `reload-atomar.sh`) im eigenen Verzeichnis; die vier Skripte mit Repo-Wurzel nehmen sie als Argument; kein Schreibziel liegt mehr im Werkzeug-Verzeichnis. Die Bind-Mounts `-v "$here":/m:ro` hängen jetzt nur Skripte ein.
- Offen (benannt, nicht verdeckt): `run.sh` führt im Kopfkommentar (Zeile 3) weiterhin den Fall `d`, den das `case` nicht kennt; das Protokoll-README ist seit `c9dc3b3` richtig. Das Skript ist bewusst nur umgezogen; die Kommentarzeile ist ein LOW-Rest (Auflage 7).

## 3. DoD Punkt für Punkt

Die Kästchen in der Slice-Datei sind alle noch offen; das Setzen ist Closure-Arbeit des Planners.

| # | DoD-Punkt | Befund | Evidence |
|---|---|---|---|
| 1 | Umzug Werkzeug: 25 Skripte (mit `probe.pl`) und 5 Dockerfiles per reinem `git mv`, eigener Commit, unter `tools/messung/<name>/` | **erfüllt** | `4e3a4d4`: 30 × `R100`; drei Verzeichnisse `devcontainer-ollama`, `sandbox-egress-neuladen-entwurf`, `egress-wildcard-pruefung` neben `ollama-amd64` |
| 2 | Umzug Protokoll: 25 Dateien per reinem `git mv` unter `docs/messungen/<kennung>/`; Register nur Beobachtungsform; nur drei `.gitkeep` | **erfüllt** | `b541202`: 25 × `R100`; die zwei `find`-Zählungen leer beziehungsweise drei `.gitkeep` (Abschnitt 1). Die 25 amd64-Rohdaten sind durch den Schwester-Slice aus dem Register entfallen |
| 3 | Pfade und Verweise nachgezogen; Rohdaten bytegleich; `bash -n` und Aufruf ohne Argumente, **kein** Messlauf | **erfüllt** | Abschnitt 2; `docs-check` und `doc-planning` grün (die Links auf Skripte und Rohdaten lösen auf); `git diff --numstat` auf Rohdaten: 0/0. Ob die Skripte nach dem Umzug messen, ist wie vorgesehen nicht belegt |
| 4 | `docs/messungen/` angelegt: README, Sub-Area-Zeile | **erfüllt** | `docs/messungen/README.md` (Zweck, Ablage, Bestand mit vier Verzeichnissen); `harness/conventions.md`: Zeile DOK führt `docs/messungen/`, Zeile WZG führt `tools/messung/`; `docs-check` ohne Befund |
| 5 | ADR-Verweise nach (b3), eigener Commit, Geschichte-Zeile im Wortlaut, `.d-check.yml` unverändert, `doc-immutable` 0 Treffer | **erfüllt** (Abschnitt 5) | `git diff e952311^ HEAD -- docs/plan/adr/` gelesen: je ADR genau eine geänderte Messungs-Zeile (Link-Ziel `../../messungen/…/README.md` und Ortsangabe) und eine neue Zeile in *Geschichte* mit „(a) a1+a3, (b) b3“ und Datum 2026-10-19; sonst nichts. Die Erwartung „0 Treffer“ ist im Plan jetzt als Gate-Lücke ausgewiesen (F-7 behoben) |
| 6 | Regel (a1): Hinweis einmal kanonisch, zwei Verweissätze | **erfüllt** | Abschnitt „Ablage von Messprotokollen und Skripten“ steht genau einmal in `harness/conventions.md` (Wortsuche im Baum außerhalb `docs/reviews/`: eine Fundstelle); `docs/plan/planning/README.md` (§Ablage) und `observations/README.md` verweisen mit je einem Satz und ohne Wiederholung; Anker löst auf |
| 7 | Regel (a3): Skript mit Positivliste, Target, Aufnahme in `make gates`, Sensors-Zeile, Wächter-Test, `selbstpruefung.sh` unverändert, MR-012 | **erfüllt, Wortlaut weicht ab** | Abschnitt 4. Die DoD sagt „Test im Docker-Betrieb des Repos“; gebaut ist ein Host-Test mit `bash` und `git` (Muster `history-range-guard.sh`). §3 der Slice-Datei und der Review führen den Host-Betrieb, der DoD-Wortlaut ist nicht nachgezogen (Auflage 1) |
| 8 | `make gates`, `make doc-planning`, `make doc-trace` grün; Review durchgeführt, Report liegt vor; Doku-Update | **erfüllt** | Abschnitt 1; Review `ccb4e31`. Öffentlicher Vertrag: nur die Gate-Kette in `README.md`, `README.de.md`, `docs/maintainer/releasing.md`, `docs/maintainer/quality.md` (Zeile 72), und sie stimmt mit `make -n gates` überein |
| 9 | Closure-Notiz mit Steering-Loop-Eintrag, Register, Risiko-Ausgänge, Paarungen | **offen: Planner** | §7 der Slice-Datei trägt „Wird bei der Closure gefüllt“; Abschnitt 6 |

## 4. Neues Gate `ablage-pruefung`

- **Positivliste (gelesen, `tools/harness/ablage-pruefung.sh`):** Eingabe ist `git ls-files -s -z -- docs` (Index, nicht Arbeitsbaum). Regeln: nur Modus `100644`; Planning-Baum nur `*.md` und die drei `.gitkeep`; unter `observations/` nur `README.md` und `BEO-<KUERZEL>/<slug>/{observation.md,state.md,evidence/<id>.md}`; `docs/messungen/` nur nicht-versteckte `.md`, `.txt`, `.log` in Kleinschreibung; im übrigen `docs/` kein Name auf `.sh .bash .pl .py .awk .mk .ps1` oder `Dockerfile*` (Groß-/Kleinschreibung gleich) und kein Shebang. Exit 0 (alles drin), 1 (Verstoß, Pfad und Grund auf stderr), 2 (kein git-Arbeitsbaum). Review F-2 (Zusage größer als Prüfung) ist mit `afce363` behoben: eine echte Positivliste mit 26 Test-Fällen (Skript ohne Endung mit Ausführungsbit und mit Shebang, `RUN.SH`, `dockerfile`, `.ps1`, `.json`, `.png`, `.TXT`, versteckte Datei, Datei ohne Endung, Symlink, `.md` mit Ausführungsbit, `.gitkeep` in einem vierten Verzeichnis, Grundzustand grün, ungetrackte Datei zählt nicht).
- **Eigene Probe (Exit-Codes):** `docs/messungen/x/probe.sh` und `docs/messungen/RUN.SH` angelegt. Ungetrackt: `make ablage-pruefung` Exit 0 (Grenze, siehe unten). Mit `git add -N`: `ablage-pruefung: docs/messungen/RUN.SH -- unter docs/messungen/ nur .md, .txt, .log (…)` und dasselbe für `probe.sh`, „2 Datei(en) ausserhalb der Positivliste“, `make gates` Exit 2. Danach Index und Dateien zurückgenommen, `git status` leer, `make ablage-pruefung` wieder grün.
- **Grenze (Hinweis, kein Befund):** Das Gate prüft den **Index**. Eine Probedatei im Arbeitsbaum, die nie `git add` sah, färbt `make gates` nicht; der Skript-Kopf und der Fall „ungetrackte Datei zählt nicht“ im Test weisen das aus, die Sensors-Zeile in `harness/README.md` und die Regel in `harness/conventions.md` nennen es nicht. Der Stop-Hook liest den Stempel von `record-gates`, also fängt ihn das Gate erst, wenn die Datei gestaged oder committet ist (Auflage 6).
- **Host-Werkzeug:** `bash` und `git`, kein Docker. Das ist mit [`AGENTS.md`](../../AGENTS.md) §3.1 vereinbar, soweit §3.1 den Go-Toolchain-Vertrag meint (Build, Lint, Tests, Coverage, Security laufen in Docker); das Muster ist im Repo etabliert (`baseline-verify.sh`, `history-range-guard.sh` stehen in derselben Kette), und `repo.mk` benennt es im Kommentar am Target. Die Beobachtung [`BEO-WZG/host-werkzeug-ausserhalb-docker`](../plan/planning/observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md) bleibt davon unberührt.
- **Einbindung:** `repo.mk`: `.PHONY`, Target `ablage-pruefung` (ruft Skript, danach Test), `GATE_CHECKS += ablage-pruefung`; Kommentarzeile zur Kette nachgezogen. `harness/README.md` §Sensors: neue Zeile mit Bindung (Ablage-Regel und Plan-Anker), Zeile `make gates` um das Gate ergänzt; kein `gate-undocumented`. Kette in `README.md`, `README.de.md`, `releasing.md`, `quality.md` (Zeile 72) gleich `make -n gates`.
- **Scharfschaltung:** zuletzt in der Commit-Reihenfolge (`7f0bdc4`), kein Carveout, keine Bootstrap-Stufe; `selbstpruefung.sh` unverändert.
- **MR-012-Ergänzung (Festlegung 3 und Geltungsbereich, Index in `harness/conventions.md`):** stimmig, deckungsgleich; die Aussage „ein zweiter Lauf entfernt eine fremde `*.sh` dort nicht“ habe ich im Klon unabhängig bestätigt (Abschnitt 1).
- **Drift, nur gemeldet:** Der Kopf von `.github/workflows/ci.yml` (Zeile 7) führt `gates` weiterhin als „lint + test + coverage-gate“, `docs/maintainer/quality.md` die Jobnamen `gates (lint + test + coverage-gate)` (Zeilen 300, 317). Vorbestehend (Review F-9), Workflow und Jobnamen sind bewusst unverändert; die Befehlszeile in `quality.md` Zeile 72 ist nachgezogen.

## 5. ADR-Diff und Gate-Lücke `doc-immutable` (unabhängig bestätigt)

**ADR-Diff** (`git diff e952311^ HEAD -- docs/plan/adr/`, gelesen, nicht über das Gate): in [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md) und [`ADR-0025`](../plan/adr/0025-egress-liste-im-laufenden-container.md) je eine geänderte Zeile (Messprotokoll-Link und Ortsangabe der Skripte) und eine neue Zeile in der Tabelle *Geschichte* (2026-10-19, Freigabe im Wortlaut „(a) a1+a3, (b) b3“). Keine Entscheidung, kein Status, keine Kennung eines Slice ändert sich; die Link-Ziele lösen auf. Review F-5 (Paraphrase statt Wortlaut) ist mit `3c094c3` behoben. Der Diff deckt sich mit der Freigabe und ist nicht weiter.

**Lücke** (`BEO-DOK/doc-immutable-greift-nicht-bei-inline-status`, [Beobachtung](../plan/planning/observations/BEO-DOK/doc-immutable-greift-nicht-bei-inline-status/observation.md)):
- `.d-check.yml` (`vcs`): `immutable-when: '^Accepted$'`, `exclude-sections: [Status, Geschichte]`; der Kommentar darüber setzt die Form „eigene Zeile unter `## Status`“ voraus und nennt `status-line` als nicht gewählte Option.
- Zählung der ADR-Dateien: 12 mit Abschnitt `## Status`, 13 ohne (alle mit `**Status:** Accepted`: ADR 0012 und 0014 bis 0025); eine Datei (ADR 0011) ist `Proposed`. Das ergibt 13 von 24 Accepted-ADRs ohne Schutz, wie behauptet.
- Eigene Probe mit `make doc-immutable STAGED=1`: Überschrift `## Kontext` in ADR 0001 (Abschnittsform) geändert, gestaged: `core-drift-vcs`, Exit 1. Überschrift `## Entscheidung` in ADR 0025 (Inline-Form) geändert, gestaged: in derselben Ausgabe kein Befund für ADR 0025 (nur der Befund für ADR 0001, „1 Befund(e)“). Probe zurückgenommen, Arbeitsbaum sauber.

**Schwere:** MEDIUM für das Repo, ohne Wirkung auf dieses Slice. Die Zusage aus [`AGENTS.md`](../../AGENTS.md) §3.5 (Accepted-ADRs sind immutable) ist für mehr als die Hälfte der Accepted-ADRs nicht gate-gestützt, und es sind **genau die zwölf jüngsten** (0014 bis 0025) plus 0012. Das Gate bleibt grün, die Lücke ist ohne Gegenprobe unsichtbar, und sie wächst mit jedem neuen ADR in der Inline-Form. Für diesen Slice ist der Schaden ausgeschlossen, weil der Diff durch Lesen belegt ist.

**Empfehlung an den Projektinhaber:** den **Sensor** beide Statusformen lesen lassen und nicht die ADRs angleichen. Gründe: (1) Die Gate-Änderung ist eine Verschärfung, keine Lockerung; sie braucht trotzdem einen eigenen Plan-Anker ([`AGENTS.md`](../../AGENTS.md) §3.6). (2) Die Angleichung wäre eine Änderung an 13 Accepted-ADRs ([`AGENTS.md`](../../AGENTS.md) §3.5) mit je einer Freigabe, ändert also genau die Dateien, die das Gate schützen soll, und lässt künftige Inline-ADRs ungeschützt. (3) Die Konfiguration kennt laut Kommentar eine Option `status-line`; ob sie neben der Abschnittsform wirkt oder ob der Sensor ein eigenes Skript braucht, ist **nicht gelesen** und gehört in die Planung. Nach der Änderung ist `doc-immutable` rückwirkend über die Historie der 13 ADRs zu fahren, um zu sehen, ob in der Zwischenzeit Drift entstand. Bis zur Entscheidung gilt: jede Änderung an einem der 13 Accepted-ADRs wird durch Lesen des Diffs belegt, wie hier. Keine Änderung an `.d-check.yml` in diesem Lauf.

## 6. Review-Findings einzeln

| ID | Klasse | Stand | Beleg |
|---|---|---|---|
| F-1 | MEDIUM | **übergeben** (Gate-Lücke im Repo, nicht im Diff) | Beobachtung angelegt (`3e36489`); Abschnitt 5; Entscheidung beim Projektinhaber, Empfehlung dort |
| F-2 | MEDIUM | **behoben** | `afce363`: Positivliste, 26 Fälle; eigene Probe `RUN.SH` und `probe.sh` rot (Abschnitt 4) |
| F-3 | LOW | **behoben** | `c9dc3b3`: README nennt `a a2 b c d2 d3 cd-c cd-d`, deckungsgleich mit der Standardliste und dem `case`; Rest im Kommentar von `run.sh` (Auflage 7) |
| F-4 | LOW | **behoben** | `51173b8`: `OUT` Pflicht, Fehler beim Kopieren sichtbar; README angepasst |
| F-5 | LOW | **behoben** | `3c094c3`: Wortlaut „(a) a1+a3, (b) b3“ in beiden Geschichte-Zeilen |
| F-6 | LOW | **behoben** | `7d27441`: kein „heute“ mehr in `state.md` und der Evidence-Datei unter der Beobachtung (`grep` leer) |
| F-7 | LOW | **behoben** | `589464d`: DoD und §6 sagen „0 Treffer, Ursache F-1“, Risiko mit eigenem Ausgang „getragen“; gelesen in der Slice-Datei |
| F-8 | INFO | **übergeben, begründet offen** | Archiv-Hinweis in `done/`-Slices ist als offene Frage mit Empfehlung (kein Hinweis; `docs/messungen/README.md` als Auffangpunkt) im Slice §6 geführt; Entscheidung beim Projektinhaber, bis dahin keine Änderung |
| F-9 | INFO | **offen, übergeben** | `ci.yml`-Kopf und Jobnamen in `quality.md` unverändert (Abschnitt 4); Entscheidung Planner/Projektinhaber |
| F-10 | INFO | **begründet offen** | Zeile DOK trägt `docs/messungen/` mit der alten Begründung; Greenfield vertretbar, Begründungsspalte nennt Messprotokolle nicht |
| F-11 | INFO | **begründet offen** | Reihenfolge `98f76a4` vor `79c749f` weicht von „ADR unmittelbar nach dem Move“ ab; beide Commits liegen im selben Bereich, kein Push-Tip dazwischen |
| F-12 | INFO | **übergeben** | `ai-harness-init` schreibt `harness/mk/traeger.mk` (Tag), bestätigt im Klon; berührt diesen Slice nicht |

## 7. Plan-Tabelle gegen `git diff --stat -M e952311..HEAD` (nach dem letzten Commit, [`AGENTS.md`](../../AGENTS.md) §3.10)

Ich habe die Dateiliste (`--name-status -M`) Zeile für Zeile gegen §3 abgeglichen. Jede der 86 Dateien hat einen Plan-Eintrag: 30 + 25 Moves; `06-option-a.sh`; die zwei Protokoll-READMEs (`R086`, `R089`); vier neue Dateien in `docs/messungen/` und `tools/harness/`; die drei neuen Dateien der Beobachtung `doc-immutable-greift-nicht-bei-inline-status`; die zwei ADRs; `harness/README.md`, `harness/conventions.md`, MR-012, `repo.mk`; die drei READMEs im Planning-Baum und `open/README.md`; `README.md`, `README.de.md`, `releasing.md`, `quality.md`; vier Register-Dateien; drei `done/`-Slices; zwei `open/`-Slices; die Slice-Datei; der Review-Report. Plan-Zeilen ohne Diff: Verifikationsreport (dieser Report) und Closure-Dateien. Es gibt weder eine Datei im Diff ohne Plan-Eintrag noch einen Plan-Eintrag ohne Diff außer den genannten. Die Zeile „unverändert“ für `selbstpruefung.sh`, `.d-check.yml`, Workflows, Lastenheft und Spezifikation ist durch den leeren `git diff --stat` belegt (Abschnitt 1).

## 8. Auflagen

An den Planner (Closure und Plan-Pflege, kein Eingriff durch den Verifier):

1. **DoD-Wortlaut Punkt 7 (a3)** von „Test im Docker-Betrieb des Repos“ auf „Test mit `bash` und `git` (Muster `history-range-guard.sh`)“ setzen; §3 und §6 (*Host-Werkzeug*) tragen den Host-Betrieb bereits, die DoD nicht.
2. **Closure-Notiz §7** mit Steering-Loop-Lerneintrag; die Verification-Evidence weist aus: Messläufe an neuem Ort **nicht gelaufen** (Abschnitt 1), `doc-immutable` 0 Treffer mit Ursache F-1, `make test-docker` und `make verify-depguard` nicht gelaufen mit Grund.
3. **Register fortschreiben** (`state.md` der betroffenen Beobachtungen; die Beobachtung zur Gate-Lücke trägt Auftreten und Bezug) und die **drei Paarungen** tragen.
4. **Risiko-Ausgänge §6:** elf Risiken tragen „weiter offen“ und brauchen bei der Closure einen Ausgang, den dieser Report belegen kann: Skript bricht nach dem Move (Abschnitt 2), Link bricht ohne Befund von `docs-check` (alle Links lösen auf, `docs-check` grün), Rohdaten ändern sich (R100, 0/0), Zeitdokumente behalten falsche Orte (Entscheidung F-8), Altlast macht das Gate rot (`make gates` grün), Positivliste zu eng (im Skript lesbar, 26 Fälle), `ai-harness-init`-Konvergenz (Klon-Beleg), Host-Werkzeug (Muster etabliert, Abschnitt 4), Sub-Area von `docs/messungen/` (Zeile DOK), amd64-Rohdaten (entfallen); Wildcard-Messung ohne README (Umsetzung im Slice vermerkt, Ausgang fehlt).
5. **Setzen der DoD-Kästchen** und **Move nach `done/`** in der Reihenfolge von [`AGENTS.md`](../../AGENTS.md) §3.3 (erst Inhalt mit Commit-Hash, dann reiner `git mv`); Roadmap und `open/README.md` mitziehen.
6. **Hinweis zur Index-Grenze des Gates** (eine Zeile in der Sensors-Zeile oder der Regel in `harness/conventions.md`: das Gate prüft den Index; eine ungestagte Datei färbt `make gates` nicht), falls der Projektinhaber es will; ohne Eingriff ist die Grenze im Skript-Kopf und im Test dokumentiert.
7. **Optional, LOW:** die Kommentarzeile 3 in `tools/messung/egress-wildcard-pruefung/run.sh` (Fall `d`) berichtigen; ein Inhaltscommit nach dem reinen Move, kein neuer Befund.

An den Projektinhaber:

8. **Gate-Lücke `doc-immutable`** (Abschnitt 5): Sensor beide Statusformen lesen lassen (Empfehlung) oder 13 ADRs angleichen (Freigabe je ADR). Beides braucht einen eigenen Plan-Anker ([`AGENTS.md`](../../AGENTS.md) §3.6).
9. **Archiv-Hinweis in `done/`-Slices** (F-8): Empfehlung kein Hinweis; die Entscheidung fehlt, und der Slice führt die Frage als offen.
10. **Drift der Gate-Kette** in `ci.yml`-Kopf (Zeile 7): Workflow-Änderung, nicht Teil dieses Slice.

Entscheidung (c) bleibt im Slice, in `open/README.md` und in `42bb8b2` als „Vorschlag übernommen, nicht widersprochen“ geführt, nirgends als Bestätigung; geprüft.

## 9. Verification Evidence (Block nach `harness/verification.md`)

Scope:
- Slice: [`slice-v1-messskripte-nach-tools`](../plan/planning/done/slice-v1-messskripte-nach-tools.md)
- IDs: [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md), [`ADR-0025`](../plan/adr/0025-egress-liste-im-laufenden-container.md) (nur Messungs-Zeile und *Geschichte*); keine `LH-*`
- Artefakte: `tools/messung/**`, `docs/messungen/**`, `tools/harness/ablage-pruefung.sh`, `tools/harness/ablage-pruefung-test.sh`, `repo.mk`, `harness/README.md`, `harness/conventions.md`, MR-012

DoD-Abgleich:
- [x] Punkte 1 bis 8 (Abschnitt 3), Punkt 7 mit Wortlaut-Auflage
- [ ] Punkt 9 (Closure) — offen: Planner, Auflagen 2 bis 5

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | Exit 0, Coverage 91.40 % |
| `make ablage-pruefung` | pass | 26 Fälle, 0 Fehlschläge; eigene Probe rot |
| `make docs-check` | pass | 509 Dateien, 0 Befunde |
| `make doc-trace`, `make doc-planning` | pass | 0 Waisen, 0 Befunde |
| `make doc-immutable` | 0 Treffer | Lücke F-1, Diff gelesen |
| `make test-docker`, `make verify-depguard` | not run | kein Docker-Pfad, keine Importregel im Diff |
| Messläufe der Skripte am neuen Ort | not run | Messhost nötig, Slice schließt sie aus |

Traceability: [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md) und [`ADR-0025`](../plan/adr/0025-egress-liste-im-laufenden-container.md) durch den gelesenen ADR-Diff und die Freigabe-Zeile; die Ablage-Regel durch das Gate und seinen Wächter-Test.

Carveouts: neu none; gelöst none; unverändert none.

Commit / Artefakt: HEAD `77dde24` vor diesem Report.
