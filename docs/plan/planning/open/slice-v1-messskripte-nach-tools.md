# Slice v1: Messskripte und Hilfs-Dockerfiles aus dem Planning-Baum nach `tools/messung/`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Der Slice ist ein eigener kleiner Aufräum-Slice **nach** dem Release-Cut `v0.12.0` (Entscheidung des Projektinhabers vom 2026-10-19: „ein eigener kleiner Aufräum-Slice … nach dem Release-Cut“); er hält den Release nicht auf und hat keine Closure-Bedingung, die von seiner DoD verschieden wäre.

**Bezug:** kein Lastenheft-Eintrag und keine ADR als Scope: der Slice ändert weder Produktverhalten noch einen Vertrag. Grundlage ist die Ablage-Regel des Projektinhabers (2026-10-19): „Skripte haben im Planning-Baum nichts verloren – es gibt doch `./tools`.“ Berührte Entscheidungen: [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) und [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md) (beide `Accepted`, bleiben unverändert, siehe §6 und die offene Entscheidung (b)); Lehre der Beobachtung [`BEO-DOK/messprotokoll-ausserhalb-des-repos`](../observations/BEO-DOK/messprotokoll-ausserhalb-des-repos/observation.md) (Protokolle gehören ins Repo, nicht ins Home-Verzeichnis; hier geht es um das *Wo im Repo*).

**Berührte Spec-Stellen:** —

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Shell-Skripte und Hilfs-Dockerfiles der drei bestehenden Messprotokolle liegen unter `tools/messung/<name>/` (Muster: `tools/harness/*.sh`, `tools/messung/ollama-amd64/`), die Rohdaten und `README.md` bleiben im Evidence-Verzeichnis, und jeder Verweis auf einen verschobenen Pfad löst auf.

### Gefundener Stand (gelesen am 2026-10-19)

- **Was im Planning-Baum liegt** (`find docs tools -type f` auf `*.sh`, `Dockerfile*`, `*.awk`, `*.py`, `*.mk`, danach alle Dateien unter `docs/` außer `*.md`, `*.txt`, `*.log`): 25 Skripte (24 Shell, 1 Perl) und 5 Hilfs-Dockerfiles in drei Protokoll-Verzeichnissen. `*.awk`, `*.py` und `*.mk` kommen unter `docs/` nicht vor; `*.json` unter `docs/` sind nur die erzeugten `devcontainer.json` der amd64-Rohdaten.

| Protokoll-Verzeichnis (unter `docs/plan/planning/observations/BEO-DOK/`) | Skripte | Dockerfiles |
|---|---|---|
| `messprotokoll-ausserhalb-des-repos/evidence/messung-slice-v1-devcontainer-ollama-messung/` | 13: `01-basis-und-optionen-c-d.sh` bis `12-negativprobe-pin.sh` (12 nummerierte) und `inner-05.sh` | 3: `Dockerfile.b`, `Dockerfile.cli`, `Dockerfile.probe` |
| `messprotokoll-ausserhalb-des-repos/evidence/messung-slice-v1-sandbox-egress-neuladen-entwurf/` | 9: `01-projekt-erzeugen.sh` bis `05-m4-neustart.sh` (5 nummerierte), `inner-m1.sh`, `inner-m2.sh`, `reload-atomar.sh` und die Perl-Sonde `probe.pl` | 2: `Dockerfile.cli`, `Dockerfile.messung` |
| `spec-bedingung-ohne-code-entsprechung/evidence/messung-slice-v1-egress-wildcard-pruefung/` | 3: `run.sh`, `inner.sh`, `inner-cd.sh` | — |

  Der Ort der Wildcard-Messung ist `spec-bedingung-ohne-code-entsprechung/`, nicht `messprotokoll-ausserhalb-des-repos/`.
- **Nicht Gegenstand, obwohl Skript oder Dockerfile** (die Sonde `probe.pl` ist dagegen Gegenstand; das `README.md` der Neuladen-Messung nennt sie in einem Punkt, Zeile 52, und beschreibt ihre Zustandswechsel; wer sie aufruft, ist nicht gelesen und wird beim Lauf geprüft): `…/evidence/messung-slice-v1-devcontainer-ollama-amd64-messung/erzeugt/*/Dockerfile` und `erzeugt/oll/egress-init.sh` sind **vom Generator erzeugte Rohdaten** der amd64-Messung (Evidence, kein Werkzeug); das Werkzeug dieser Messung liegt bereits unter `tools/messung/ollama-amd64/` (`01-messstand.sh`, `02-projekte.sh`, `03-bau.sh`, `lib.sh`). Der Slice [`slice-v1-devcontainer-ollama-amd64-messung`](../in-progress/slice-v1-devcontainer-ollama-amd64-messung.md) arbeitet dort; dieser Slice fasst weder ihn noch seine Dateien an.
- **Wie die Skripte ihren Ort benutzen** (gelesen, Zeilen der Skripte): `02-build-b.sh` und `12-negativprobe-pin.sh` bauen mit `-f "$here/Dockerfile.b"` und `$here` als Kontext; `04-build-probe.sh` mit `-f "$here/Dockerfile.probe"`; `05-laufzeit.sh` hängt `$here` als `/m` ein und ruft `bash /m/inner-05.sh`; `02-m1-rechte.sh` liest `$here/inner-m1.sh`, `03-m2-neuladen.sh` hängt `$here` als `/m` ein (`/m/inner-m2.sh`); `run.sh` hängt `$here` als `/m` ein (`/m/inner.sh`, `/m/inner-cd.sh`). Skript und Hilfsdatei müssen daher **im selben Verzeichnis** bleiben: Skripte und Dockerfiles wandern zusammen.
- **Wo ein Skript in sein eigenes Verzeichnis schreibt:** `06-option-a.sh` Zeile 37 kopiert das Build-Log nach `"$(dirname "$0")/rohdaten-06-build-a.log"`; nach dem Umzug läge die Rohdatei im Werkzeug-Verzeichnis statt im Evidence-Verzeichnis. Die übrigen Skripte schreiben Rohdaten über die Umleitung des Aufrufers (`… > rohdaten-NN.txt 2>&1`, Kopfzeile `Aufruf:`); der Aufruf steht in `README.md` des Protokolls („Jedes Skript wird im eigenen Verzeichnis aufgerufen“). Weitere `$0`-abgeleitete Schreibziele sind beim Implementer-Lauf gegen alle 25 Skripte zu prüfen; hier gelesen sind nur die oben genannten.
- **Verweis auf die Repo-Wurzel:** `01-projekt-erzeugen.sh`, `03-projekt-erzeugen.sh`, `05-m4-neustart.sh` und `run.sh` nehmen die Repo-Wurzel als Argument (`${1:?repo-wurzel}`); `run.sh` liest daraus `internal/hexagon/application/templates/devcontainer/egress-init.sh.tmpl`. Eine Pfadannahme über die Lage im Repo steht dort nicht.
- **Muster in `tools/`:** `tools/messung/ollama-amd64/lib.sh` leitet `ROOT` drei Ebenen über dem Skript ab und nimmt `OUT` (Evidence-Verzeichnis) und `WORK` (Messordner außerhalb des Repos) aus der Umgebung; `tools/harness/*.sh` sind Werkzeuge des Harness (Aufruf über `make`). Die Sub-Area `tools/`, `scripts/` ist WZG, Greenfield ([`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area)).
- **Eingehende Verweise — Links** (`](…)` auf ein verschobenes Skript, gesucht über alle 30 Dateinamen in allen `*.md`): fünf, in zwei Dateien: [`slice-v1-egress-wildcard-pruefung`](../done/slice-v1-egress-wildcard-pruefung.md) (zwei Links auf `run.sh`, Zeilen 108 und 146) und die Evidence `spec-bedingung-ohne-code-entsprechung/evidence/slice-v1-egress-wildcard-pruefung.md` (drei Links auf `run.sh`, `inner.sh`, `inner-cd.sh`, Zeile 4). Diese brechen beim Umzug, falls der `links`-Check von d-check das Ziel prüft; ob er Nicht-`md`-Ziele auflöst, ist hier nicht gelesen (`make docs-check` am 2026-10-19: 480 Dateien, 0 Befunde, die fünf Links lösen heute auf) und wird im ersten Implementer-Lauf durch einen Probe-Umzug belegt.
- **Eingehende Verweise — Prosa** (Dateiname ohne Link, gesucht über die 29 Namen der Shell-Skripte und Dockerfiles; `probe.pl` kommt außerdem in denselben Dateien vor, die die Liste nennt (Slice [`slice-v1-sandbox-egress-neuladen-entwurf`](../done/slice-v1-sandbox-egress-neuladen-entwurf.md) und zwei Reviews), und in `inner-m2.sh`, dem Nachbarn im Protokoll): 18 Dateien außerhalb der Protokoll-Verzeichnisse: `docs/plan/adr/0024-ollama-im-sandbox-devcontainer.md` (Zeile 28, nennt `Dockerfile.probe` ohne Pfad), die Slices [`slice-v1-devcontainer-ollama-messung`](../done/slice-v1-devcontainer-ollama-messung.md), [`slice-v1-egress-wildcard-pruefung`](../done/slice-v1-egress-wildcard-pruefung.md), [`slice-v1-sandbox-egress-neuladen-entwurf`](../done/slice-v1-sandbox-egress-neuladen-entwurf.md), [`slice-v2-spec-sandbox-egress-ollama`](../done/slice-v2-spec-sandbox-egress-ollama.md), [`slice-v1-devcontainer-start-ohne-cli`](slice-v1-devcontainer-start-ohne-cli.md) und der laufende amd64-Slice, die Beobachtungs-Dateien unter `BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/`, `BEO-DOK/spec-bedingung-ohne-code-entsprechung/evidence/` und `BEO-WZG/host-werkzeug-ausserhalb-docker/evidence/` sowie sechs Dateien unter `docs/reviews/`. Ein Dateiname ohne Pfad bricht nicht; **falsch werden** dort die Ortsangaben („liegen im Repo, neben der Evidence“, „im Verzeichnis dieser Datei“).
- **ADRs:** [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) (Zeile 28) und [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md) (Zeile 23) verlinken ausschließlich das `README.md` des jeweiligen Protokolls (`…/evidence/messung-slice-v1-devcontainer-ollama-messung/README.md` beziehungsweise `…/messung-slice-v1-sandbox-egress-neuladen-entwurf/README.md`), das **an seinem Ort bleibt**; kein Link der ADRs zeigt auf ein Skript oder Dockerfile. [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) und [`ADR-0022`](../../adr/0022-fallback-ausweisung-statusdateien-im-container.md) trugen in der Wortsuche nur den Namen `egress-init.sh` des Generators (kein Messskript). Der Bruch bei den ADRs ist daher eine **Ortsaussage in Prosa** ([`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md): „Befehle, Dockerfiles und Rohdaten: Messprotokoll (13 Skripte, 3 Dockerfiles, 15 Rohdaten-Dateien)“ — das Messprotokoll-README listet sie nach dem Umzug nicht mehr als Nachbarn), kein toter Link.
- **Werkzeug für den Nachzug:** `make slice-mv` (`tools/harness/slice-mv.sh`) bewegt nur Slices und zieht deren Verweise nach (Kopfzeilen des Skripts); für Skripte gibt es kein Werkzeug. `tools/harness/selbstpruefung.sh` trägt in seinem Kopf „DIESE DATEI IST KONVERGENT: jeder Lauf des Werkzeugs schreibt sie kanonisch neu“; ein Edit dort ginge beim nächsten `ai-harness-init`-Lauf verloren.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Messläufe wiederholen** — die Rohdaten sind Belege eines Messstands (Host, Docker-Version, Datum); ein Lauf nach dem Umzug wäre eine neue Messung, kein Umzug.
- **Rohdaten ändern** (`rohdaten-*.txt`, Logs, `erzeugt/`) — sie bleiben bytegleich im Evidence-Verzeichnis; ein Diff darauf wäre ein Befund.
- **`README.md` der Protokolle inhaltlich umschreiben** — es bekommt den neuen Skriptpfad und die angepassten Aufrufe, Messstand und Ergebnis bleiben.
- **Ein Accepted-ADR ändern** ([`AGENTS.md`](../../../../AGENTS.md) §3.5) — wie der Bruch in [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) und [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md) aufgefangen wird, ist die offene Entscheidung (b) des Projektinhabers.
- **Ein neues Gate, ein neues `make`-Target oder ein Sensor ohne Entscheidung** ([`AGENTS.md`](../../../../AGENTS.md) §3.6) — Entscheidung (a).
- **Das Werkzeug und die Rohdaten der amd64-Messung** — sie liegen schon richtig (`tools/messung/ollama-amd64/`, Evidence-Verzeichnis); sie sind Gegenstand des laufenden Slice [`slice-v1-devcontainer-ollama-amd64-messung`](../in-progress/slice-v1-devcontainer-ollama-amd64-messung.md).
- **Zeitdokumente in `docs/reviews/` und `done/`-Slices inhaltlich korrigieren** über Querverweise und Archiv-Hinweise hinaus — `docs/plan/planning/README.md` erlaubt in `done/` nur korrigierende Änderungen (Tippfehler, Querverweise, Archiv-Hinweise); ein Review-Report ist ein Zeitdokument des Prüfstands.

## 2. Definition of Done

- [ ] **Umzug:** Alle 25 Skripte (einschließlich `probe.pl`) und 5 Dockerfiles aus §1 liegen per reinem `git mv` (30 Dateien, eigener Commit, [`AGENTS.md`](../../../../AGENTS.md) §3.3) unter `tools/messung/<name>/` (Verzeichnisnamen: Entscheidung (c)); im Planning-Baum bleibt kein Skript und kein Hilfs-Dockerfile außer den vom Generator erzeugten Rohdaten (`erzeugt/`); `find docs -name '*.sh' -o -name 'Dockerfile*'` zeigt nur diese.
- [ ] **Pfade und Verweise nachgezogen** (Inhaltscommits nach dem Move): die Pfade in den Skripten arbeiten vom neuen Ort aus (`06-option-a.sh` schreibt die Rohdatei ins Evidence-Verzeichnis über `OUT` oder ein Argument, der Muster-Ort ist `tools/messung/ollama-amd64/lib.sh`); `README.md` der drei Protokolle nennt den Skriptpfad und die angepassten Aufrufe; die fünf Links aus §1 und die Ortsaussagen in Prosa lösen auf beziehungsweise stimmen (soweit die Zeitdokument-Grenze aus §1 es zulässt); die Rohdaten sind bytegleich (`git diff --stat` gegen den Move-Commit zeigt keine Zeile unter `rohdaten-*` und `erzeugt/`); ein Skript, das von seinem neuen Ort aus startet, wird mit `bash -n` und einem Aufruf ohne Argumente (Fehlmeldung `${1:?…}`) geprüft — **kein** Messlauf.
- [ ] **ADR-Verweise nach Entscheidung (b)** und **Regel nach Entscheidung (a)** umgesetzt, wie der Projektinhaber sie entscheidet (§6); fällt (a) auf „Sensor“, ist das eine Gate-Änderung mit eigenem Plan-Anker und Eintrag in [`harness/README.md`](../../../../harness/README.md) §Sensors.
- [ ] `make gates` grün, `make doc-planning` und `make doc-trace` grün; Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — kein Self-Review; Doku-Update soweit ein öffentlicher Vertrag berührt ist (keiner erwartet: die Skripte sind Maintainer-Werkzeug, [`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist; jedes Risiko aus §6 trägt einen Ausgang; die drei Paarungen sind getragen.

## 3. Plan (vor Code)

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; der Implementer ergänzt sie in seinem ersten Lauf und gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab ([`AGENTS.md`](../../../../AGENTS.md) §3.10). Jede Zeile ist eine Erwartung, keine gelesene Aussage über unveränderte Dateien. Reihenfolge der Commits: (1) reiner `git mv` aller 30 Dateien, (2) Pfad-Berichtigung in den Skripten, (3) Nachzug in READMEs, Links und Prosa, (4) Entscheidungen (a) und (b), (5) Closure. Move und Inhalt nie im selben Commit.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| 13 Skripte und 3 Dockerfiles aus `…/messprotokoll-ausserhalb-des-repos/evidence/messung-slice-v1-devcontainer-ollama-messung/` | move (`git mv`) nach `tools/messung/<name>/` | Gegenstand; Skript und Hilfsdatei bleiben zusammen (§1 *Wie die Skripte ihren Ort benutzen*) |
| 9 Skripte (einschließlich `probe.pl`) und 2 Dockerfiles aus `…/messung-slice-v1-sandbox-egress-neuladen-entwurf/` | move nach `tools/messung/<name>/` | dito |
| 3 Skripte aus `…/spec-bedingung-ohne-code-entsprechung/evidence/messung-slice-v1-egress-wildcard-pruefung/` | move nach `tools/messung/<name>/` | dito |
| die verschobenen Skripte (insbesondere `06-option-a.sh`, Zeile 37; die übrigen `$0`-Schreibziele beim Lauf prüfen) | update | Pfad-Berichtigung nach dem Move, getrennter Commit; Skripte nehmen den Evidence-Ort als Argument oder Variable (`OUT`, Muster `lib.sh` der amd64-Messung) |
| `README.md` der drei Protokolle (`…/messung-slice-v1-devcontainer-ollama-messung/README.md`, `…/messung-slice-v1-sandbox-egress-neuladen-entwurf/README.md`, falls vorhanden das der Wildcard-Messung, sonst `spec-bedingung-ohne-code-entsprechung/evidence/slice-v1-egress-wildcard-pruefung.md`) | update | neuer Skriptpfad, angepasste Aufrufe (`docker build -f Dockerfile.cli`, Ausgabe-Umleitung) |
| `docs/plan/planning/done/slice-v1-egress-wildcard-pruefung.md` (Zeilen 108, 146) | update (nur Querverweis) | zwei Links auf `run.sh`; `done/` erlaubt Querverweis-Korrekturen |
| `docs/plan/planning/observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/evidence/slice-v1-egress-wildcard-pruefung.md` (Zeile 4) | update | drei Links auf `run.sh`, `inner.sh`, `inner-cd.sh` |
| übrige 16 Dateien mit Prosa-Nennung (§1: [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md), Slices in `done/` und `open/`, Beobachtungs-Dateien, `docs/reviews/`) | unverändert oder update, je nach Entscheidung (b) und (c) | Dateiname ohne Pfad bricht nicht; Ortsaussagen in Zeitdokumenten werden nicht umgeschrieben (§1); der Implementer liest die Stelle, bevor er „unverändert“ einträgt (§3.10) |
| `docs/plan/adr/0024-ollama-im-sandbox-devcontainer.md`, `…/0025-egress-liste-im-laufenden-container.md` | unverändert (Accepted) | [`AGENTS.md`](../../../../AGENTS.md) §3.5; Auffangen des Bruchs nach Entscheidung (b) im `README.md` des Protokolls |
| `harness/conventions.md` (Modus-Deklaration, Zeile `tools/`, `scripts/`) und/oder `docs/plan/planning/README.md`, oder ein Prüf-Skript unter `tools/harness/` mit Anbindung in `repo.mk` | update oder neu, nach Entscheidung (a) | Regel gegen neue Skripte im Planning-Baum |
| `docs/reviews/…-slice-v1-messskripte-nach-tools.md` (Review) und `…-verification-…` | neu | Review und Verifikation |
| diese Slice-Datei, `open/README.md`, Register, Roadmap bei Start und Closure | update | Lifecycle |

## 4. Trigger

**Start** (`open` → `in-progress`, ggf. über `next`): `Verantwortlich:` gesetzt; der Release-Cut `v0.12.0` ist geschnitten; der Vorgänger [`slice-v1-devcontainer-ollama-amd64-messung`](../in-progress/slice-v1-devcontainer-ollama-amd64-messung.md) liegt in `done/`; `in-progress/` ist leer (WIP-Limit 1); der Projektinhaber hat die Entscheidungen (a), (b) und (c) aus §6 getroffen, mindestens (b) und (c) vor dem Move.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Nachzug der Pfade in den Skripten braucht mehr als eine Variable je Skript (zum Beispiel Skripte, die ihre Umgebung aus dem Ort ableiten); dann trennt der Architect den Umzug (reiner Move, Nachzug der Links) von der Skript-Überarbeitung.
- `in-progress` → `open` (blockiert): Die Entscheidung (b) oder (c) fehlt; oder ein Slice mit Zugriff auf dieselben Evidence-Verzeichnisse (WIP-Limit) liegt in `in-progress/`.

## 5. Closure-Trigger

DoD vollständig, Review-Report und Verification-Evidence liegen vor, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7); die Rohdaten sind bytegleich belegt.

## 6. Risiken und offene Punkte

**Offene Entscheidungen des Projektinhabers** (der Plan entscheidet sie nicht):

- **(a) Regel gegen neue Skripte im Planning-Baum.** Optionen:
  - *a1 — Hinweis ohne Sensor:* ein Satz in [`harness/conventions.md`](../../../../harness/conventions.md) (Zeile `tools/`, `scripts/` der Modus-Deklaration) und in `docs/plan/planning/README.md`: Skripte und Hilfs-Dockerfiles liegen unter `tools/`, im Planning-Baum nur Rohdaten und `README.md`. Aufwand: klein, zwei Textstellen; kein Gate. Durchsetzung nur durch Review.
  - *a2 — Prüfung in der Doku-Konfiguration:* falls `d-check` ein Verbot für ein Dateimuster kennt (`.d-check.yml`). Ob es das kann, ist hier nicht gelesen; der Architect liest die d-check-Fähigkeiten vor dieser Option. Aufwand: klein bei Erfolg, Gate-Änderung ([`AGENTS.md`](../../../../AGENTS.md) §3.6).
  - *a3 — eigenes Prüf-Skript:* ein Skript unter `tools/harness/` (zum Beispiel `git ls-files docs | grep -E '\.sh$|Dockerfile'` mit Ausnahme von `erzeugt/`), angebunden an ein Target in `repo.mk` und in `harness/README.md` §Sensors, Aufnahme in `make gates` oder als on-demand-Werkzeug. Aufwand: mittel (Skript, Selbsttest gegen eine Probe-Datei, Gate-Index, Ausnahmenliste für `erzeugt/`); eine neue Gate-Zeile ist eine Gate-Änderung mit eigenem Plan-Anker. Ein Edit in `tools/harness/selbstpruefung.sh` scheidet aus (konvergente Datei, §1).
  - *a4 — keine Regel:* der Umzug allein; Risiko: das Muster kehrt bei der nächsten Messung zurück (Beobachtung [`BEO-DOK/messprotokoll-ausserhalb-des-repos`](../observations/BEO-DOK/messprotokoll-ausserhalb-des-repos/observation.md) zählt es dann).
- **(b) Umgang mit den Verweisen in den Accepted-ADRs** ([`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) Zeile 28, [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md) Zeile 23). Der Link zeigt auf das `README.md` des Protokolls und bleibt gültig (§1); falsch wird nur die Nachbar-Aussage („Dockerfiles … im Messprotokoll“). Optionen:
  - *b1 — ADRs unverändert, Auffangen im Protokoll:* das `README.md` des Protokolls (Ziel des ADR-Links) trägt oben einen Hinweis „Skripte und Dockerfiles liegen seit … unter `tools/messung/<name>/`“. Aufwand: eine Zeile je Protokoll; §3.5 gewahrt, der Leser des ADR kommt über den Link zum Hinweis.
  - *b2 — Folge-ADR oder Fortschreibung als neue Entscheidung:* ein neuer ADR mit `Supersedes`/Verweis. Aufwand: hoch für einen reinen Ablageort; der Plan hält ihn für unverhältnismäßig, weil keine Entscheidung der ADRs berührt ist, nur der Fundort eines Belegs.
  - *b3 — Querverweis-Korrektur in den ADRs:* ein Satz im ADR-Text. Das ist eine Änderung eines Accepted-ADR und berührt §3.5 (und `make doc-immutable`); sie geschieht nur mit ausdrücklicher Freigabe des Projektinhabers.
  Empfehlung des Planners: b1.
- **(c) Verzeichnisnamen unter `tools/messung/`.** Vorschlag: der Protokollname ohne das Präfix `messung-slice-v1-`, also `devcontainer-ollama`, `sandbox-egress-neuladen-entwurf`, `egress-wildcard-pruefung` (neben `ollama-amd64`); Alternative: derselbe Name wie das Evidence-Verzeichnis. Die Entscheidung legt die Pfade fest, die die READMEs und Links nennen.

**Risiken** (Ausgang nach Closure):

- **Ein Skript bricht nach dem Move, weil es den Ort des Evidence-Verzeichnisses aus seinem eigenen Ort ableitet** (`06-option-a.sh`, Zeile 37, gelesen; weitere nicht gelesen). — **Ausgang:** weiter offen bis zum ersten Implementer-Lauf; die Prüfung ist `bash -n` und ein Aufruf ohne Argumente, kein Messlauf; ein unklarer Fall wird dokumentiert statt gedeutet.
- **Ein Link auf ein verschobenes Skript bricht, und `docs-check` fängt es nicht** (die fünf Links aus §1; ob `d-check` Nicht-`md`-Ziele prüft, ist ungelesen). — **Ausgang:** weiter offen; der Probe-Umzug im ersten Lauf belegt es; fängt `docs-check` es nicht, geht der Befund als Beobachtung ins Register (Sensor-Grenze), und die fünf Links werden von Hand geprüft.
- **Die Rohdaten ändern sich unbemerkt** (etwa durch ein Skript, das beim Probe-Aufruf in das Evidence-Verzeichnis schreibt). — **Ausgang:** weiter offen; der Abgleich `git diff --stat` auf `rohdaten-*` und `erzeugt/` nach jedem Commit ist Teil der DoD.
- **Zeitdokumente behalten falsche Ortsangaben** (Reviews, `done/`-Slices, Beobachtungs-Evidence: „liegen im Repo, neben der Evidence“). — **Ausgang:** weiter offen; die Grenze ist in §1 benannt, der Hinweis im `README.md` des Protokolls (b1) ist der Auffangpunkt; ob `done/`-Slices einen Archiv-Hinweis bekommen, entscheidet der Projektinhaber mit (b).
- **Das Muster kehrt zurück** (neue Messung legt Skripte in den Planning-Baum, solange (a) auf a1 oder a4 steht). — **Ausgang:** weiter offen → [`BEO-DOK/messprotokoll-ausserhalb-des-repos`](../observations/BEO-DOK/messprotokoll-ausserhalb-des-repos/observation.md) beziehungsweise nach Entscheidung (a) der Regel-Ort.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `tools/`, `scripts/` (WZG, Greenfield: die Skripte materialisieren zuerst geschriebene Konventionen) und `spec/`, `harness/`, `docs/plan/` sowie `docs/reviews/` (DOK, Greenfield) laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); beide erfüllen die Schwelle aus zwei von drei Achsen laut Modus-Deklaration (eigenes Verzeichnis, eigene Konvention: `tools/harness/*.sh` und die Docker-only-Regel; eigene Inventur-Linie: die Planning-Lifecycle-Prüfung beziehungsweise das Werkzeug-Verzeichnis). Eine zu grobe Sub-Area ist nicht erkennbar; `tools/messung/` ist ein Unterverzeichnis von WZG, keine eigene Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-DOK/messprotokoll-ausserhalb-des-repos`](../observations/BEO-DOK/messprotokoll-ausserhalb-des-repos/observation.md) (offen, 1×; Gegenstand dieses Slice ist der Ort *im* Repo, nicht der Verstoß *außerhalb*; ein Beleg für die Beobachtung ist der Umzug nicht, er stützt aber die Frage, ob die Lehre „Protokoll im Repo“ um „Werkzeug unter `tools/`“ zu schärfen ist), [`BEO-WZG/kopplung-nicht-nachgezogen`](../observations/BEO-WZG/kopplung-nicht-nachgezogen/observation.md) (offen; der Nachzug der Verweise nach dem Move ist die Kopplung), [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.10; „unverändert“ in der Plan-Tabelle braucht eine gelesene Stelle) und [`BEO-WZG/host-werkzeug-ausserhalb-docker`](../observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md) (die Skripte brauchen auf dem Host Docker und die Shell; der Slice führt sie nicht aus). Evidenz-/Diskrepanz-Risiko: niedrig bis mittel (Ortsaussagen in Zeitdokumenten).

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
