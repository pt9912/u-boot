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
[`MR-002`](../../../../harness/conventions/done/MR-002-carveout-inventar.md) (Baseline-Regelwerk
`modul-07-carveouts.md` §Ziel-Form: Carveout) um. Wegen des WIP-Limits 1
(belegt durch [`slice-cr-devcontainer-java-toolchain`](../done/slice-cr-devcontainer-java-toolchain.md))
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

**Verantwortlich:** pt9912 (Implementer; Umsetzung nach dem freigegebenen Plan des Projektinhabers, Start beim freien WIP-Limit)

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
`internal/hexagon/port/driven/filesystem.go`). [`MR-002`](../../../../harness/conventions/done/MR-002-carveout-inventar.md)
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
| Test-Carveouts `_test.go` (acht Linter) | [ADR-0003](../../adr/0003-solid-nahes-lint-profil.md) §Entscheidung (getragen) | die `Why:`-Kommentare der Ausschluss-Blöcke nennen den ADR — **berichtigung nach F-1(c)**: der achte („revive`/`unused-parameter`) bekam seinen Block-Kommentar erst in der Fix-Runde; der Nachbar-Block (`unused-receiver`, kein Inventar-Eintrag) trägt seinen Namen ohne ADR-Ausweisung (Bestand) |
| depguard `!**/*_test.go` als erstes files-Pattern | [SPEC-034](../../../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln) (getragen) | **berichtigung nach F-1(b):** der depguard-Kopf-Kommentar nennt [ADR-0042](../../adr/0042-dauerhafte-lint-ausnahmen.md) entschieden — Entscheidung 3 (bis zur Fix-Runde stand hier „—“) |
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
  [`slice-v1-logs-crlf.md`](../open/slice-v1-logs-crlf.md)); dieses Slice legt nur dessen Plan
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

- [x] [`docs/plan/carveouts/`](../../carveouts/) trägt die Baseline-Ziel-Form:
      README-Index (aus `README.template.md`) und eine CO-Datei für den
      CRLF-Carveout (aus `carveout.template.md`); das `.gitkeep` ist
      entfallen; jede Zeile des früheren Inventars trägt einen Anker nach der
      Zuordnungstabelle in §1 — geprüft je Zeile mit Fundstelle in der
      Closure.
- [x] [`MR-002`](../../../../harness/conventions/done/MR-002-carveout-inventar.md) ist aufgelöst: reiner
      `git mv` nach `harness/conventions/done/` (AGENTS.md §3.3), die
      Tabellen in [`conventions.md`](../../../../harness/conventions.md) zeigen die Resolution;
      [`.d-check.yml`](../../../../.d-check.yml) (matrix `carveout`, ids-Target `CO`), die
      Gate-Index-/Pointer-Zeilen ([`harness/README.md`](../../../../harness/README.md) §Guides,
      [`harness/verification.md`](../../../../harness/verification.md), `.harness/skills/closure-note-reviewer.md`,
      [`docs/plan/planning/README.md`](../README.md), die Sub-README unter `in-progress/` (im Slice gelöscht — nicht Regelwerks-konform),
      [`roadmap.md`](../in-progress/roadmap.md)) zeigen auf die neuen Orte; `make docs-check` 0
      Befunde.
- [x] Der neue ADR (MADR, Folgeentscheidung „Ergänzt, ersetzt nichts“ zu
      [ADR-0003](../../adr/0003-solid-nahes-lint-profil.md)) trägt `contextcheck` und `interfacebloat`
      (Fundstellen der bislang kommentar-getragenen Ausnahmen gelesen:
      `.golangci.yml` `issues.exclude-rules`); die `Why:`-Kommentare und
      [`docs/maintainer/quality.md`](../../../../docs/maintainer/quality.md) §1.3 nennen die tragenen ADRs;
      `make lint` grün (Konfig-Werte unverändert).
- [x] `make gates` grün. *(nicht gezählt — Gate-Lauf)*
- [x] Review, Verifikation, Closure-Notiz, Beobachtungs-Register, Risiko-Ausgänge,
      Drei-Paarungen. *(nicht gezählt — Closure-Pflichten)*

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — Pfad-Kandidaten für §8.

| Datei / Komponente | Art | Begründung / Fundstelle |
|---|---|---|
| docs/plan/carveouts/README.md (Link bei der Anlage, die Datei entsteht in diesem Slice) | neu | Index aus `README.template.md` (Baseline `templates/docs/plan/carveouts/`); Nachfolger des Master-Inventars |
| `docs/plan/carveouts/CO-<NNN>-crlf-logs.md` | neu | aus `carveout.template.md`; trägt den CRLF-Edge-Case (Fundstellen gelesen 2026-11-02: `splitLogLines` in `internal/adapter/driving/cli/logs.go`, Absatz in `docs/maintainer/cli-json-contract.md` §6.8); die Kennung ist bei der Anlage gesetzt — CO-001, Verweise in dieser Datei berichtigt |
| `docs/plan/planning/in-progress/carveouts.md` | gelöscht | die Inventar-Ablage geht auf — jede Zeile trägt einen Anker weiter (§1), Löschung im Umsetzungs-Commit |
| `docs/plan/carveouts/.gitkeep` | gelöscht | das Verzeichnis trägt Dateien; kein Bestand im `ablage-prüfung`-Positivbestand (mutiert werden nur die drei `.gitkeep` unter `docs/plan/planning/`) |
| [`slice-v1-logs-crlf.md`](../open/slice-v1-logs-crlf.md) | neu | Folge-Slice-Plan des CRLF-Carveouts, aus `slice.template.md` |
| `docs/plan/adr/0042-<slug>.md` | neu | aus `NNNN-titel.template.md` (MADR, [MR-008](../../../../harness/conventions/MR-008-adr-form.md)); trägt `contextcheck` (Fundstelle: `Why:`-Kommentar zu `path: internal/adapter/driving/cli/`) und `interfacebloat` (Fundstelle: `Why:`-Kommentar zu `path: internal/hexagon/port/driven/filesystem\.go$`) |
| [`docs/plan/adr/README.md`](../../adr/README.md) | update | Index-Zeile für den neuen ADR (AGENTS.md §5.3) |
| [`.golangci.yml`](../../../../.golangci.yml) | update | nur Warum-Kommentare: die Kommentar-Blöcke der `exclude-rules` nennen den tragenden ADR (`errcheck`-Settings, `_test.go`-Zeilen, `cmd/uboot/`, `contextcheck`, `interfacebloat`); kein Wert oder Regeln geändert |
| [`docs/maintainer/quality.md`](../../../../docs/maintainer/quality.md) | update | §1.3 um die drei der Tabelle fehlenden Ausnahmen erweitert (gelesen: Tabelle trägt vier Zeilen, `_test\.go$`/`cmd/uboot/`/`errcheck`); Referenzmodell: maintainer-Doku trägt ADR-Verweise |
| [`docs/plan/planning/README.md`](../README.md) | update | die Carveout-Ausnahme und die zwei Inventar-Verweise berichtigt (Zeilen 51 und 57 bis 58, 93 bis 94) |
| `docs/plan/planning/in-progress/README.md` | gelöscht | nicht Regelwerks-konform — das Baseline-Regelwerk kennt keinen Sub-README; Projektinhaber-Anweisung 2026-11-02, Bestand war veraltet 
| [`docs/plan/planning/in-progress/roadmap.md`](../in-progress/roadmap.md) | update | Pointer-Zeilen (Zeilen 137 bis 139) auf docs/plan/carveouts/README.md (Link bei der Anlage, die Datei entsteht in diesem Slice) |
| [`docs/plan/planning/open/README.md`](../open/README.md) | update | Bestand-Zeilen im Verweis-Nachzug; Bestand-Zeile für den Folge-Slice-Plan [`slice-v1-logs-crlf`](../open/slice-v1-logs-crlf.md) (Review F-6) und die Trichter-Formulierung; Link-Lauf des CR-Slices |
| Berichte unter `docs/reviews/` (Review des Laufs, Verify des Laufs) | neu | Review-Runde und Verification-Runde (DoD-Zeilen des Slice); **Nachtrag nach Verifier V-1** — die Zeilen fehlten in der Erst-Tabelle |
| diese Datei | update | Plan-Berichtigungen der Runden (F-1, F-2, V-1; §3.10-Abgleich nach jeder Runde) |
| [`.d-check.yml`](../../../../.d-check.yml) | update | `matrix`-Eintrag `carveout` → `docs/plan/carveouts/CO-*.md`; `ids`-Teil: Target der CO-Regel → `docs/plan/carveouts/` (aktuell `spec/lastenheft.md`, keine CO-ID in Verwendung — geprüft per Suche, 0 Verwendungen) |
| [`harness/README.md`](../../../../harness/README.md) | update | Guides-Zeile (Carveout-Verweis) und die Master-Inventar-Zeile in der Rote-Gates-Liste |
| [`harness/verification.md`](../../../../harness/verification.md) | update | Zeile 93 bis 94: Parallel-Eintrag der neuen temporären Carveouts zeigt auf [`docs/plan/carveouts/`](../../carveouts/) |
| `.harness/skills/closure-note-reviewer.md` | update | die zwei Inventar-Verweise (Zeilen 31 bis 33 und 75 bis 76) auf den Index |
| [`harness/conventions.md`](../../../../harness/conventions.md) | update | Aktive-Tabelle: `MR-002`-Zeile in die Aufgelöste-Tabelle (aufgelöst durch diesen Slice); Bereichsangabe der Zeile berichtigt |
| `harness/conventions/MR-002-carveout-inventar.md` | git mv nach `done/` | reiner Move-Commit (Muster [`MR-011`](../../../../harness/conventions/done/MR-011-werkzeug-fragmente-repo-gefuehrt.md)); **berichtigung nach F-1(a)/F-2:** die Formulierung „Inhalt unverändert“ trug der Stand nicht — der Record bekam danach zwei Link-/Pfad-Berichtigungs-Commits (f021780f, e04e126c, dann die F-2-Rücksetzung der Prosa); der Herkunfts-Regel zuliebe: der `git mv` zieht die Pfad-Berichtigung nach sich, sie berührt Referenzen, nicht die Adaption-Sätze |
| Verweis-Nachzug | update | alle Dateien, die `planning/in-progress/carveouts.md` hrefen oder nennen (132 Treffer im Baum, geprüft am 2026-11-02: 69 in `../../carveouts/README.md`-Form aus `done/`/`next/`/`open/`, 63 in `docs/plan/planning/in-progress/carveouts.md`-Form aus Reviews, Maintainer-Doku, Archiv und Konventionen) — mechanische Berichtigung; Abgleich gegen `git diff --stat` nach jeder Fix-Runde (AGENTS.md §3.10) |

**Ansatz als Liste** (die Referenz-Berichtigung ist gruppen-gleichrangig):

- Die Verweis-Berichtigung folgt zwei Regeln: (1) jede Zeile, deren Pfadteil
  `docs/plan/planning/in-progress/carveouts.md` trägt, zeigt künftig auf
  `docs/plan/carveouts/README.md` am selben Wurzel-Anker; (2) die
  Geschwister-Form `../in-progress/carveouts.md` aus `done/`, `next/` und
  `open/` wird zu `../../carveouts/README.md`. Textnachweise in **lebenden**
  Dokumenten (Planung, Gate-Index, Skills) folgten derselben Ersetzung;
  **berichtigung nach F-2:** die link-freien Backtick-Nennungen der alten
  Ablage im aufgelösten Konvention-Record ([`MR-002`](../../../../harness/conventions/done/MR-002-carveout-inventar.md))
  bleiben unverändert stehen — die Pfad-Berichtigung deckt Referenzen, nicht
  die Inhaltsaussage eines Records.
- Die historischen Done-Reporte halten ihre Zitate (Berichte sind
  korrigierend änderbar, see `.d-check.yml` Kommentar): nur der href
  wandert, kein Inhalt.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Der Projektinhaber gibt den Auftrag;
kein anderer Slice in `in-progress/` (WIP-Limit 1 — Stand 2026-11-02:
belegt durch [`slice-cr-devcontainer-java-toolchain`](../done/slice-cr-devcontainer-java-toolchain.md),
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
  `roadmap.md`/die Sub-README unter `in-progress/` (im Slice gelöscht — nicht Regelwerks-konform) am Rand — gelesen: der CR-Slice
  fügt Wellen- und Bestand-Zeilen hinzu, der hier nur eine Carveout-Zeile
  korrigiert). — **Ausgang:** bei Closure; der Start-Trigger bindet den
  Slice an die freie WIP-Lage.

## 7. Closure-Notiz

Der Slice ist umgesetzt, Review- und Verifier-Reports liegen vor und ihre Befunde sind berichtigt; geschlossen wird er mit dem Closure-Block am Datei-Ende (das DoD-Häckchen der Pflichten und der Verifikations-Lauf nach dem Move und dem Verweis-Nachzug).

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
### Closure (Planner, 2026-11-02, nach Verifier-Fix-Runde `ff528d08`)

**Evidence:** Review-Report [`2026-11-02-slice-harness-carveouts-baseline-form`](../../../reviews/2026-11-02-slice-harness-carveouts-baseline-form.md) (0 HIGH / 2 MEDIUM / 6 LOW / 1 INFO, Verdikt *nicht merge-blockierend*; Negativbefunde: `.golangci.yml` wortwörtlich ohne Werte-Diff, ADR-0003 unberührt, Ziel-Form von CO-001 und Index gegen die Baseline-Vorlagen konform, MR-002-Move reines `R100`, Trichter-Spur lesbar, `.d-check.yml`-Bindung korrekt) und Verification-Report [`2026-11-02-slice-harness-carveouts-baseline-form-verify`](../../../reviews/2026-11-02-slice-harness-carveouts-baseline-form-verify.md) (DoD-Zeilen 1 bis 3 belegt erfüllt — je-Zeile-Anker der acht dauerhaften Inventar-Zeilen voll deckend geprüft, `make lint` Exit 0; die Reste V-1 und V-2 (LOW, reine Planning-Texte) sind berichtigt `ff528d08`). Sensors der Runden: `make docs-check` 0 Befunde (1045/1046/1047 Dateien je Stand), `make baseline-verify` v6.18.1 OK, `make ablage-pruefung` 26 Fälle 0 Fehlschläge, `make lint` 0 issues.

**Zuordnung der dauerhaften Inventar-Zeilen zu ihren Ankern** (Beleg der DoD-Zeile 1, Fundstellen je Zeile stehen in §1; Verifier geprüft): `errcheck` `fmt.Fprint*`, `testpackage`/`gochecknoglobals`, `_test.go`-Cluster → [ADR-0003](../../adr/0003-solid-nahes-lint-profil.md); depguard `!**/*_test.go`, `contextcheck`, `interfacebloat` → [ADR-0042](../../adr/0042-dauerhafte-lint-ausnahmen.md) (E3, E1, E2) bzw. [SPEC-034](../../../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln); GNU `make` → [SPEC-031](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows). Der historische Audit-Trail des Inventory (die Tabelle der Auflösungs-Slices) lebt in der Git-Historie der Datei `docs/plan/planning/in-progress/carveouts.md` und in deren Einträgen der done-Slices; eine Kopie entsteht nicht (Referenzmodell, [ADR-0013](../../adr/0013-dokumentationsreferenzmodell.md)).

**Risiken aus §6, Ausgang je Punkt:** (1) Verweis-Nachzug unvollständig — **entfallen** (zwei Nachzugs-Runden; `make docs-check` 0 am Schluss; der Verifier-Fund der Reste war Planning-Text, kein Link). (2) Inventar-Zeile ohne Anker — **entfallen** (je-Zeile-Prüfung voll deckend, Verifier). (3) CO-ID-Linkpflichten — **entfallen** (die Pflichten sind getragen: CO-001 gelinkt im Index, ADR-0042 im Index; `make docs-check` 0). (4) Konflikt mit dem laufenden CR-Slice — **entfallen** (der CR-Slice schloss mit `c1de2582` vor dem Start dieses Laufs; der Start-Trigger hielt die Reihenfolge, WIP-Limit frei).

**Steering-Loop-Eintrag** (gezählt, nicht verkörpert; kein Eintrag trägt das Feld `liegt in`): *Lehre:* der §3.10-Abgleich lief nicht vollständig nach den Runden — die Berichte-Dateien fehlten als §3-Zeilen und sechs „Inventar“-Vokabeln des Nachzugs standen (Verifier V-1, V-2); Register [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md). Die Lehre der letzten Closures gilt: der Abgleich läuft nach jeder Berichtigungs-Runde erneut — die Klasse trat hier auf, obwohl die Lehre getragen wurde, weil der Abgleich den Umfang der Runden (Berichte, Vokabel-Nachzug) nicht reichhaltig nahm. Die Sensor-Fragen der früheren Closures bleiben beim Projektinhaber. *Neuer Sensor:* keiner. *Benannte Spec-Lücke:* keine.

**Beobachtungs-Register:** *Evidence neu:* [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (`evidence/slice-harness-carveouts-baseline-form.md`: F-1, F-2, V-1, V-2). *Nicht gezählt:* [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md) — kein Fund (kein Vertrags-Diff im Slice).

**Carveouts:** keiner — die Ablage-Form selbst ist der Gegenstand; CO-001 bleibt als aktiver temporärer Carveout mit Folge-Slice (kein Gate gelockert, keine Schwelle gesenkt; Review-Negativbefund).

**Weiter offen:** der Folge-Slice [`slice-v1-logs-crlf`](../../../planning/open/slice-v1-logs-crlf.md) (CO-001-Konsequenz, Start am Trigger); der Umsetzungs-Slice [`slice-v1-katalog-feature-jdk-gepinnt`](../../../planning/open/slice-v1-katalog-feature-jdk-gepinnt.md) mit seinen Rest-Aufträgen und die Ableitung des Release-Cut (F-6 des CR-Reviews: kein Cut, der den java-Schlüssel als implementiert meldet, solange der Katalog-Code ihn nicht trägt); die open/-README-Sub-Datei und das `next/`-README analog zur in-progress/-Sub-Datei zu prüfen (gleiche Regeln, kein Aufhebungsauftrag).

**Drei Paarungen:** Anker — kein `liegt in`-Feld, nichts verkörpert; Folge-Slice — [`slice-v1-logs-crlf`](../../../planning/open/slice-v1-logs-crlf.md) liegt als Datei in `open/` (`ls` bei der Closure); Register — jede genannte Kennung ist ein Verzeichnis mit der Evidence-Datei dieser Closure (`ls`). Der Welle-Bezug ist der der nächsten Welle-Closure (Modul 6).

**Verifikation** (Docker-Modus): `make gates` läuft nach Closure-Inhalt, Move, Register- und Verweis-Nachzug auf dem Abschluss-Commit; der `record-gates`-Stempel trägt den Nachweis; die Commit-Meldung des Nachzugs trägt das Ergebnis. Der Lauf gilt für den Stand vor dieser Zeile; die Zeile ändert nur Prosa.
