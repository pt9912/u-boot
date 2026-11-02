# Slice Harness: Carveout-Ablage auf den Baseline-Ort — `docs/plan/carveouts/` trägt die CO-Form, das Master-Inventar löst sich auf

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle (Harness-Wartung; es gibt keine Closure-Bedingung, die
von der DoD dieses Slice verschieden ist).

**Anlass (Projektinhaber, 2026-11-02).** Anweisung: „Die carveouts in
`docs/plan/planning/in-progress` gehören in den Ordner `docs/plan/carveouts`“
— und auf Nachfrage die Werkzeug-Wahl: **volle Baseline-Form** (CO-Dateien,
permanente Einträge über den Trichter). Das kehrt die Adaption
[`MR-002`](../../../../harness/conventions/MR-002-carveout-inventar.md) (Baseline-Regelwerk
`modul-07-carveouts.md` §Ziel-Form: Carveout) um. Wegen des WIP-Limits 1
(belegt durch [`slice-cr-devcontainer-java-toolchain`](../in-progress/slice-cr-devcontainer-java-toolchain.md))
läuft die Umsetzung in einem eigenen Slice, sobald kein anderer Slice in
`in-progress/` ist; dieser Plan wartet in `open/`.

**Berührte Spec-Stellen:** `—` — es wird keine Spec-Stelle berührt. Die
dauerhaften Abweichungen sind bereits getragen: [SPEC-034](../../../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln)
(depguard-Regeln gelten production-only, `*_test.go`-Dateien ausgenommen),
[SPEC-031](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) (`make` als Carveout zu
[`LH-NFA-PORT-002`](../../../../spec/lastenheft.md#lh-nfa-port-002--keine-unnötigen-systemabhängigkeiten)),
[SPEC-035](../../../../spec/spezifikation.md#spec-035--pflichten-des-lint-profils) (Carveouts zentral mit
`Why:`-Kommentar). Der neue ADR (Zählnummer aus [`docs/plan/adr/README.md`](../../adr/README.md)
bei der Anlage) trägt die zwei dauerhaften Lint-Ausnahmen, die bislang keinen
ADR-Namen tragen.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner; Anweisung und Werkzeug-Wahl des
Projektinhabers 2026-11-02). **Datum:** 2026-11-02.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Die Carveout-Ablage trägt die Baseline-Ziel-Form
(`modul-07-carveouts.md` §Ziel-Form: Carveout): [`docs/plan/carveouts/`](../../carveouts/) hält
einen README-Index und je aktiven Carveout eine CO-Datei aus der Vorlage; der
Master-Inventar `docs/plan/planning/in-progress/carveouts.md` geht auf. Für
jede Zeile des Inventars trägt ein Anker den Inhalt weiter: der eine
temporäre Carveout (CRLF-Edge-Case in `splitLogLines`) als CO-Datei mit
Folge-Slice-Plan, die dauerhaften Einträge an den Ankern, die sie bereits
tragen ([ADR-0003](../../adr/0003-solid-nahes-lint-profil.md) trägt vier der
acht; [SPEC-031](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) den
`make`-Carveout; [SPEC-034](../../../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln) den
depguard-Test-Ausschluss) — und ein neuer ADR trägt die zwei, die bislang
nur in `.golangci.yml`-Kommentaren stehen (`contextcheck` für
`internal/adapter/driving/cli/`, `interfacebloat` für
`internal/hexagon/port/driven/filesystem.go`). [`MR-002`](../../../../harness/conventions/MR-002-carveout-inventar.md)
wird aufgelöst (reiner `git mv` nach `conventions/done/`, Muster
[`MR-011`](../../../../harness/conventions/done/MR-011-werkzeug-fragmente-repo-gefuehrt.md)).

**Zuordnung der acht dauerhaften Inventar-Zeilen** (geprüft gegen die Anker,
gelesen 2026-11-02 — Fundstellen: [ADR-0003](../../adr/0003-solid-nahes-lint-profil.md) Zeilen
62 bis 72; [SPEC-034](../../../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln); [SPEC-031](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows); `.golangci.yml`
`issues.exclude-rules`):

| Inventar-Zeile | Neuer Anker | Berührung im Slice |
|---|---|---|
| `errcheck` `fmt.Fprint*` | [ADR-0003](../../adr/0003-solid-nahes-lint-profil.md) §Entscheidung (getragen) | Why-Kommentar nennt den ADR |
| `testpackage` + `gochecknoglobals` für `cmd/uboot/` | [ADR-0003](../../adr/0003-solid-nahes-lint-profil.md) §Entscheidung (getragen) | Why-Kommentar nennt den ADR |
| Test-Carveouts `_test.go` (acht Linter) | [ADR-0003](../../adr/0003-solid-nahes-lint-profil.md) §Entscheidung (getragen) | Warum-Kommentare nennen den ADR |
| depguard `!**/*_test.go` als erstes files-Pattern | [SPEC-034](../../../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln) (getragen) | — |
| GNU `make` am Host | [SPEC-031](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) (getragen) | — |
| `contextcheck` für `internal/adapter/driving/cli/` | **neuer ADR** (bislang nur `Why:`-Kommentar) | Anlage + Why-Kommentar nennt den ADR |
| `interfacebloat` für `driven/filesystem.go` | **neuer ADR** (bislang nur `Why:`-Kommentar) | Anlage + Why-Kommentar nennt den ADR |
| Master-Inventar selbst (Einrichten, Audit-Trail) | `docs/plan/carveouts/README.md` (Index) + Closure dieses Slice | — |

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Keine Neubewertung der Begründungen.** Die Gründe der acht dauerhaften
  Zeilen stehen in [ADR-0003](../../adr/0003-solid-nahes-lint-profil.md), [SPEC-031](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows), [SPEC-034](../../../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln) und den
  `Why:`-Kommentaren; der Slice transponiert Ablage und Anker und urteilt
  nicht neu (das wäre ein anderer Vorgang — Werkzeug-Wahl des Trichters ist
  mit der Projektinhaber-Wahl 2026-11-02 abgeschlossen).
- **Kein Fix des CRLF-Edge-Cases.** Er ist das Objekt des Folge-Slice (siehe
  [`slice-v1-logs-crlf.md`](slice-v1-logs-crlf.md)); dieses Slice legt nur dessen Plan
  und die CO-Datei an.
- **Kein Auflösungs-Entzug der Lint-Ausnahmen.** Die `.golangci.yml`-Werte
  (Schwellen, Ausschlüsse, Pattern) bleiben wortwörtlich; der Diff berührt
  nur Kommentare. Eine geänderte Ausnahme wäre eine Gate-Änderung (AGENTS.md
  §3.6) und ein eigener Lauf.
- **Keine Änderung an den drei Spec-Straten.** Lastenheft, Spezifikation und
  Architektur-Sicht tragen nichts zu übergeben (dazu Decken-Regel, AGENTS.md
  §3.4); die Referenz-Richtung bleibt: Slice/ADR zeigen aufwärts, die Spec
  nennt diesen Slice nie.
- **Kein Umbau des Beobachtungs-Registers.** Keine Beobachtung zählt neu; das
  Register liest nur (Closure), es wird nicht geschärft.

**Keine Mindestzahl** — vier Ausschlüsse sind die echten Grenzen.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — ≤ 3 Liefer-Punkte.

- [ ] [`docs/plan/carveouts/`](../../carveouts/) trägt die Baseline-Ziel-Form:
      README-Index (aus `README.template.md`) und eine CO-Datei für den
      CRLF-Carveout (aus `carveout.template.md`); das `.gitkeep` ist
      entfallen; jede Zeile des früheren Inventars trägt einen Anker nach der
      Zuordnungstabelle in §1 — geprüft je Zeile mit Fundstelle in der
      Closure.
- [ ] [`MR-002`](../../../../harness/conventions/MR-002-carveout-inventar.md) ist aufgelöst: reiner
      `git mv` nach `harness/conventions/done/` (AGENTS.md §3.3), die
      Tabellen in [`conventions.md`](../../../../harness/conventions.md) zeigen die Resolution;
      [`.d-check.yml`](../../../../.d-check.yml) (matrix `carveout`, ids-Target `CO`), die
      Gate-Index-/Pointer-Zeilen ([`harness/README.md`](../../../../harness/README.md) §Guides,
      [`harness/verification.md`](../../../../harness/verification.md), `.harness/skills/closure-note-reviewer.md`,
      [`docs/plan/planning/README.md`](../README.md), [`in-progress/README.md`](../in-progress/README.md),
      [`roadmap.md`](../in-progress/roadmap.md)) zeigen auf die neuen Orte; `make docs-check` 0
      Befunde.
- [ ] Der neue ADR (MADR, Folgeentscheidung „Ergänzt, ersetzt nichts“ zu
      [ADR-0003](../../adr/0003-solid-nahes-lint-profil.md)) trägt `contextcheck` und `interfacebloat`
      (Fundstellen der bislang kommentar-getragenen Ausnahmen gelesen:
      `.golangci.yml` `issues.exclude-rules`); die `Why:`-Kommentare und
      [`docs/maintainer/quality.md`](../../../../docs/maintainer/quality.md) §1.3 nennen die tragenen ADRs;
      `make lint` grün (Konfig-Werte unverändert).
- [ ] `make gates` grün. *(nicht gezählt — Gate-Lauf)*
- [ ] Review, Verifikation, Closure-Notiz, Beobachtungs-Register, Risiko-Ausgänge,
      Drei-Paarungen. *(nicht gezählt — Closure-Pflichten)*

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — Pfad-Kandidaten für §8.

| Datei / Komponente | Art | Begründung / Fundstelle |
|---|---|---|
| docs/plan/carveouts/README.md (Link bei der Anlage, die Datei entsteht in diesem Slice) | neu | Index aus `README.template.md` (Baseline `templates/docs/plan/carveouts/`); Nachfolger des Master-Inventars |
| `docs/plan/carveouts/CO-<NNN>-crlf-logs.md` | neu | aus `carveout.template.md`; trägt den CRLF-Edge-Case (Fundstellen gelesen 2026-11-02: `splitLogLines` in `internal/adapter/driving/cli/logs.go`, Absatz in `docs/maintainer/cli-json-contract.md` §6.8); die Kennung ist bei der Anlage gesetzt, dann Berichtigung der Verweise in dieser Datei |
| `docs/plan/carveouts/.gitkeep` | gelöscht | das Verzeichnis trägt Dateien; kein Bestand im `ablage-prüfung`-Positivbestand (mutiert werden nur die drei `.gitkeep` unter `docs/plan/planning/`) |
| [`slice-v1-logs-crlf.md`](slice-v1-logs-crlf.md) | neu | Folge-Slice-Plan des CRLF-Carveouts, aus `slice.template.md` |
| `docs/plan/adr/0042-<slug>.md` | neu | aus `NNNN-titel.template.md` (MADR, [MR-008](../../../../harness/conventions/MR-008-adr-form.md)); trägt `contextcheck` (Fundstelle: `Why:`-Kommentar zu `path: internal/adapter/driving/cli/`) und `interfacebloat` (Fundstelle: `Why:`-Kommentar zu `path: internal/hexagon/port/driven/filesystem\.go$`) |
| [`docs/plan/adr/README.md`](../../adr/README.md) | update | Index-Zeile für den neuen ADR (AGENTS.md §5.3) |
| [`.golangci.yml`](../../../../.golangci.yml) | update | nur Warum-Kommentare: die Kommentar-Blöcke der `exclude-rules` nennen den tragenden ADR (`errcheck`-Settings, `_test.go`-Zeilen, `cmd/uboot/`, `contextcheck`, `interfacebloat`); kein Wert oder Regeln geändert |
| [`docs/maintainer/quality.md`](../../../../docs/maintainer/quality.md) | update | §1.3 um die drei der Tabelle fehlenden Ausnahmen erweitert (gelesen: Tabelle trägt vier Zeilen, `_test\.go$`/`cmd/uboot/`/`errcheck`); Referenzmodell: maintainer-Doku trägt ADR-Verweise |
| [`docs/plan/planning/README.md`](../README.md) | update | die Carveout-Ausnahme und die zwei Inventar-Verweise berichtigt (Zeilen 51 und 57 bis 58, 93 bis 94) |
| [`docs/plan/planning/in-progress/README.md`](../in-progress/README.md) | update | Bestand-Zeile `carveouts.md` entfällt; die Dauerdokumente-Anmerkung |
| [`docs/plan/planning/in-progress/roadmap.md`](../in-progress/roadmap.md) | update | Pointer-Zeilen (Zeilen 137 bis 139) auf docs/plan/carveouts/README.md (Link bei der Anlage, die Datei entsteht in diesem Slice) |
| [`.d-check.yml`](../../../../.d-check.yml) | update | `matrix`-Eintrag `carveout` → `docs/plan/carveouts/CO-*.md`; `ids`-Teil: Target der CO-Regel → `docs/plan/carveouts/` (aktuell `spec/lastenheft.md`, keine CO-ID in Verwendung — geprüft per Suche, 0 Verwendungen) |
| [`harness/README.md`](../../../../harness/README.md) | update | Guides-Zeile (Carveout-Verweis) und die Master-Inventar-Zeile in der Rote-Gates-Liste |
| [`harness/verification.md`](../../../../harness/verification.md) | update | Zeile 93 bis 94: Parallel-Eintrag der neuen temporären Carveouts zeigt auf [`docs/plan/carveouts/`](../../carveouts/) |
| `.harness/skills/closure-note-reviewer.md` | update | die zwei Inventar-Verweise (Zeilen 31 bis 33 und 75 bis 76) auf den Index |
| [`harness/conventions.md`](../../../../harness/conventions.md) | update | Aktive-Tabelle: `MR-002`-Zeile in die Aufgelöste-Tabelle (aufgelöst durch diesen Slice); Bereichsangabe der Zeile berichtigt |
| `harness/conventions/MR-002-carveout-inventar.md` | git mv nach `done/` | reiner Move-Commit, Inhalt unverändert (Muster [`MR-011`](../../../../harness/conventions/done/MR-011-werkzeug-fragmente-repo-gefuehrt.md)) |
| Verweis-Nachzug | update | alle Dateien, die `planning/in-progress/carveouts.md` hrefen oder nennen (132 Treffer im Baum, geprüft am 2026-11-02: 69 in `../in-progress/carveouts.md`-Form aus `done/`/`next/`/`open/`, 63 in `docs/plan/planning/in-progress/carveouts.md`-Form aus Reviews, Maintainer-Doku, Archiv und Konventionen) — mechanische Berichtigung; Abgleich gegen `git diff --stat` nach jeder Fix-Runde (AGENTS.md §3.10) |

**Ansatz als Liste** (die Referenz-Berichtigung ist gruppen-gleichrangig):

- Die Verweis-Berichtigung folgt zwei Regeln: (1) jede Zeile, deren Pfadteil
  `docs/plan/planning/in-progress/carveouts.md` trägt, zeigt künftig auf
  `docs/plan/carveouts/README.md` am selben Wurzel-Anker; (2) die
  Geschwister-Form `../in-progress/carveouts.md` aus `done/`, `next/` und
  `open/` wird zu `../../carveouts/README.md`. Textnachweise (Backticks,
  Prosa) folgten derselben Ersetzung.
- Die historischen Done-Reporte halten ihre Zitate (Berichte sind
  korrigierend änderbar, see `.d-check.yml` Kommentar): nur der href
  wandert, kein Inhalt.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Der Projektinhaber gibt den Auftrag;
kein anderer Slice in `in-progress/` (WIP-Limit 1 — Stand 2026-11-02:
belegt durch [`slice-cr-devcontainer-java-toolchain`](../in-progress/slice-cr-devcontainer-java-toolchain.md),
Projektinhaber-Entscheidung 2026-11-02: dieser Plan wartet, siehe Anlass).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): —, kein Fall
  vorweg zu sehen; der Umfang ist die Ablage-Transposition.
- `in-progress` → `open` (blockiert — Carveout?): Falls die ADR-Grund-
  lage (`contextcheck`/`interfacebloat`) im Lauf inhaltlich umstreit wird,
  kehrt der Slice zurück; die ADR-Grundlage ist dann beim Projektinhaber
  zur Entscheidung vorzulegen.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln.

DoD vollständig (drei Liefer-Punkte plus die Gate- und Closure-Pflichten);
Review und Verifikation; Closure-Notiz mit Zuordnungstabelle (jede
Inventar-Zeile → neuer Anker) als Beleg; `git mv` nach `done/`.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Verweis-Nachzug unvollständig** (132 Treffer, mehrere Pfad-Formen): eine
  verbleibende Referenz ins Leere bricht `make docs-check`. — **Ausgang:**
  bei Closure (`make docs-check` 0 Befunde; Abgleich gegen
  `git diff --stat`).
- **Zeile des Inventars ohne neuen Anker:** der Trichter lässt eine der acht
  dauerhaften Zeilen ohne getragenen Anker. — **Ausgang:** bei Closure
  (Zuordnungstabelle in der §1-Fundstellen-Prüfung je Zeile; keine Zeile
  ohne Anker).
- **CO-ID-Linkpflichten des d-check** — die neuen CO-/ADR-Kennungen
  erzeugen Linkpflichten an Stellen, die sie noch nicht verlinkten. —
  **Ausgang:** bei Closure (`make docs-check` iterativ; die Pflichten
  stehen in `.d-check.yml` §ids).
- **Konflikt mit dem laufenden CR-Slice** (beide berühren
  `roadmap.md`/[`in-progress/README.md`](../in-progress/README.md) am Rand — gelesen: der CR-Slice
  fügt Wellen- und Bestand-Zeilen hinzu, der hier nur eine Carveout-Zeile
  korrigiert). — **Ausgang:** bei Closure; der Start-Trigger bindet den
  Slice an die freie WIP-Lage.

## 7. Closure-Notiz

Noch nicht geschlossen; der Slice liegt in `open/` (Entwurf).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

**Vorgelagert — Sub-Area-Wahl prüfen:** Die berührten Verzeichnisse liegen
in deklarierten Sub-Areas ([DOK](../../../../harness/conventions.md#modus-deklaration-pro-sub-area)
für `spec/`, `harness/`, `docs/plan/`; WZG für `.d-check.yml` und
`.golangci.yml` als Werkzeug-Dokument), keine neu zu differenzierende
Kandidatin.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen am
2026-11-02 — keine Treffer: [`BEO-DOK/abhilfe-ohne-wirkungsprobe`](../../../../docs/plan/planning/observations/README.md)
(offen, 1×) betrifft Abhilfe-Meldungen der Katalog-Features, nicht die
Carveout-Ablage; [`BEO-DOK/plan-aussage-ueber-code`](../../../../docs/plan/planning/observations/README.md)
ist verkörpert ([AGENTS.md](../../../../AGENTS.md) §3.10). Keine Zählung, keine
Evidenz-Datei im Register für diesen Plan.

**Modus-Begründungsblock:** entfällt — alle berührten Sub-Areas GF (Refactor
ohne neue Sub-Area-Berührung; keine BF-/Hybrid-Markierung).