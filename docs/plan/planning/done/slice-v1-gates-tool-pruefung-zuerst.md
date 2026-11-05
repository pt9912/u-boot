# Slice V1: `tool-pruefung` läuft in `make gates` als erster Check — eine falsche Werkzeugkette fällt vor Lint, Test und Coverage auf

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle, Release-Zuordnung offen (kleine Änderung an der Gate-Kette;
keine Closure-Bedingung, die von der DoD dieses Slice verschieden ist).

**Bezug:** [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) (Scope: das Gate `tool-pruefung` gehört zur Pin-Inventur), [`ADR-0043`](../../adr/0043-nur-devcontainer-bau-pfad.md) (Bau-Pfad; die Werkzeuge des laufenden Kontexts tragen).

**Berührte Spec-Stellen:** — (das Lastenheft nennt `tool-pruefung` als Gate der Pin-Inventur, eine Reihenfolge der Checks sagt es nicht zu; gelesen: `spec/lastenheft.md` Pin-Inventur-Absatz bei [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette)).

**Verantwortlich:** —

**Autor:** Claude (Planner-Entwurf auf Wunsch des Projektinhabers). **Datum:** 2026-11-05.

---

## 1. Ziel und Abgrenzung

**Ziel:** In einem seriellen `make gates` läuft `tool-pruefung` vor `baseline-verify`, `docs-check`, `lint`, `test` und `coverage-gate`.

Ausgangslage (gelesen): `gates` hängt an `record-gates`, und `record-gates` an `$(GATE_CHECKS)` (`Makefile`). Die Liste entsteht in dieser Folge: die Fragmente unter `harness/mk/` (`baseline-verify`, `docs-check`), dann `repo.mk` (`lint test coverage-gate`), dann `tools/mk/*.mk` (`ablage-pruefung`, `tool-pruefung`), eingebunden in `repo.mk` hinter der Zeile `GATE_CHECKS += lint test coverage-gate`. `tool-pruefung` stand damit am Ende; eine abweichende Werkzeugkette fiel erst nach den langen Gates auf.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Reihenfolge unter `make -j`** — `Makefile` hält `-j` bewusst offen (`.NOTPARALLEL` ist nicht gewählt); eine Zusage für den parallelen Lauf bräuchte Order-only-Voraussetzungen an `lint`, `test` und die übrigen Checks. Weder `scripts/bau.sh` noch ein Workflow unter `.github/workflows/` ruft `make` mit `-j` auf (gelesen per Suche nach `-j`, ohne Treffer). Kommt `-j` in die Rezepte, ist das ein Folge-Slice.
- **Änderung des Gate-Inhalts** (`tools/harness/tool-pruefung.sh`) — die Prüfung bleibt wie sie ist; verschoben wird nur der Zeitpunkt.
- **Reihenfolge der übrigen Checks untereinander** — kein Befund, der sie begründet.

## 2. Definition of Done

- [x] **Reihenfolge:** `make -n gates` zeigt `tool-pruefung` als ersten Check; ein Test oder Sensor pinnt das, damit die Zusage aus dem Kommentar in `tools/mk/tool-pruefung.mk` nicht still bricht (Form wählt der Implementer, Kandidat: ein Shell-Test unter `tools/harness/` neben `ablage-pruefung-test.sh`).

Pflichten der Closure (zählen nicht mit):

- [x] `make gates` grün (im Devcontainer-Container, `AGENTS.md` §3.1).
- [x] Review durchgeführt, Report unter `docs/reviews/` (`.harness/skills/reviewer.md`), kein Self-Review.
- [x] Doku-Zeiger nachgesehen: `README.md`, `README.de.md`, `harness/README.md` nennen `tool-pruefung` als Teil der Menge der Gates; eine Reihenfolge behaupten sie nach Lesung nicht.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag; jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/mk/tool-pruefung.mk` | update | `GATE_CHECKS := tool-pruefung $(GATE_CHECKS)` statt `+=` (Commit `6c4edaa9`); Kommentar nennt die Zusage (nur serieller Lauf); das Rezept von `tool-pruefung` ruft Sensor und Wächter-Test auf |
| `tools/harness/gate-reihenfolge.sh` | neu | Sensor: liest die Prerequisite-Liste von `record-gates` aus `make -pn gates` und verlangt `tool-pruefung` an erster Stelle ([`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette)) |
| `tools/harness/gate-reihenfolge-test.sh` | neu | Wächter-Test: Grundzustand grün, Mutation (`+=`) rot; Form nach `ablage-pruefung-test.sh` |
| `docs/reviews/2026-11-05-slice-v1-gates-tool-pruefung-zuerst.md` | neu | Review-Report (Rolle Reviewer, 0 HIGH, 1 MEDIUM) |
| `docs/reviews/verification-2026-11-05-slice-v1-gates-tool-pruefung-zuerst.md` | neu | Verification-Report (Rolle Verifier) |
| `docs/plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/` | update | `state.md` und `evidence/slice-v1-gates-tool-pruefung-zuerst.md` (Closure, §7) |
| `docs/plan/planning/in-progress/roadmap.md` | update | Sektion „Offene Wellen“ nennt den Slice statt des Ruhe-Markers (Drift-Befund `planning-drift` von `docs-check`); bei der Closure kehrt der Ruhe-Marker zurück |
| `README.md`, `README.de.md`, `harness/README.md` | unverändert | gelesen: Zeilen mit `make gates` führen die Gates als Menge mit `+` und nennen keine Reihenfolge (`README.md:281`, `README.de.md:287`, `harness/README.md:72`) |
| `docs/plan/planning/open/README.md` | update | Zeile des Slice (Zustand `in-progress/`, kein Chronik-Satz) |
| `docs/plan/planning/done/slice-v1-devcontainer-homebrew-pakete.md` | neu | anderer Slice, im selben Anlege-Commit `9182d4bd`; nicht Gegenstand dieses Slice |
| `docs/plan/planning/done/slice-v1-gates-tool-pruefung-zuerst.md` | neu / update | dieser Plan selbst (angelegt in `open/`, per `make slice-mv` hierher) |
| `harness/README.md`, `harness/conventions.md`, `harness/conventions/MR-012-ai-harness-init-werkzeug-des-repos.md` | update | Review F-5: Zeile `make tool-pruefung` nennt den Sensor; die zwei neuen Skripte stehen als repo-geführte Ausnahme in `tools/harness/` |
| `repo.mk`, `Makefile` | unverändert | gelesen: `repo.mk:254` hängt `lint test coverage-gate` an, `repo.mk:256` bindet `tools/mk/*.mk` ein; die Reihenfolge ergibt sich allein aus der Position der Anhängestelle |

Vor dem Handoff gleicht der Implementer diese Tabelle gegen `git diff --stat` ab (`AGENTS.md` §3.10).

## 4. Trigger

**Start** (`next` → `in-progress`): Projektinhaber priorisiert; keine Abhängigkeit.

**Rückführungen:**

- `in-progress` → `next` (zu groß): wenn der Sensor der Reihenfolge ein eigenes Werkzeug braucht, das mehr als ein Skript ist.
- `in-progress` → `open` (blockiert): wenn die Rezepte von `scripts/bau.sh` `-j` einführen, bevor der Slice schließt — dann ist die Order-only-Variante Teil des Gegenstands.

## 5. Closure-Trigger

DoD vollständig, `make gates` grün im Devcontainer, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

1. **Der Sensor der Reihenfolge hängt am Text von `make -n`** und bricht, wenn sich die Rezeptzeilen ändern. — **Ausgang:** entfallen: der Sensor liest die Position in der Prerequisite-Liste aus `make -pn gates` (`tools/harness/gate-reihenfolge.sh`), nicht den Text einer Rezeptzeile; Läufe mit `MAKEFLAGS=-j4` und `env -i` blieben grün (Review-Report).
2. **Ein späterer `-j`-Lauf** hebt die Zusage auf, ohne dass ein Gate rot wird. — **Ausgang:** entfallen: die Rezepte von `scripts/bau.sh` und die Workflows rufen `make` ohne `-j` auf (gelesen), der Kommentar in `tools/mk/tool-pruefung.mk` benennt die Grenze, und führt ein Rezept `-j` ein, greift die Rückführung nach `open/` aus §4.

## 7. Closure-Notiz

- **Was hat funktioniert:** Die Reihenfolge ist eine Position in der Prerequisite-Liste von `record-gates`, und der Sensor liest genau diese Liste aus `make -pn gates`; der Wächter-Test färbt die Mutation (`+=`) rot. Beide laufen als Teil von `make tool-pruefung` im Gate.
- **Was ging anders als geplant:** Der Wechsel nach `in-progress/` färbte `docs-check` rot, weil die Roadmap den Ruhe-Marker „Nichts in Arbeit.“ trug (`planning-drift`). Die Plan-Tabelle war zweimal unvollständig (Review F-2, Verification A-1), und die Zeile in `open/README.md` trug nach dem Wechsel noch Marker des Zustands davor (Review F-1).
- **Steering-Loop-Eintrag:** `BEO-DOK/plan-aussage-ueber-code` ergänzt (weiterer Beleg, Rückfall nach der Verkörperung in `AGENTS.md` §3.10); keine weitere Schärfung entschieden.
- **Beobachtungs-Register (`../observations/`):** `evidence/slice-v1-gates-tool-pruefung-zuerst.md` in `BEO-DOK/plan-aussage-ueber-code/` ergänzt.
- **Folge-Slices:** keine. Zwei Vorbefunde außerhalb des Slice stehen im Verification-Report und sind nicht Gegenstand: `tools/harness/tool-pruefung.sh` fehlt in der Ausnahmeliste von `MR-012`, `docs/maintainer/quality.md` führt `tool-pruefung` in der Gate-Liste nicht.
- **Risiken aus §6:** beide mit Ausgang `entfallen` (siehe §6).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist `tools/`, `scripts/` (WZG) und `docs/plan/` (DOK); beide stehen in `harness/conventions.md` §Modus-Deklaration.

**Vorgelagert — offene Beobachtungen sichten:** vor dem Übergang nach `next/` auszuführen (Register `docs/plan/planning/observations/`); der Entwurf trägt noch kein Ergebnis.

**Modus-Begründungsblock:** Alle berührten Sub-Areas sind Greenfield. Kein BF-/Hybrid-Block nötig.
