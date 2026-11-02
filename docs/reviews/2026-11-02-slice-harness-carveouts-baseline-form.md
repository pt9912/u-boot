# Review-Report: slice-harness-carveouts-baseline-form — 2026-11-02

**Review-Art:** Code-Review (Doku-/Konfig-Diff, kein Produktcode) — geprüft gegen den
Slice-Plan (§1 Zuordnungstabelle, §2 DoD, §3 Plan-Tabelle, §6 Risiken), die Hard Rules
(`AGENTS.md` §3, §3.10-Abgleich), die Baseline-Ziel-Forms
(`v6.18.1` · `regelwerk/modul-07-carveouts.md` §Ziel-Form: Carveout,
`regelwerk/modul-04-adrs.md` §Ziel-Form: ADR (MADR),
`regelwerk/grundlagen-traceability.md` §Herkunfts-Anker) und
`harness/conventions.md` (Adaptions-Ledger).

**Gegenstand:** die Slice-Commits `2b705b00` (Carveout-Ablage Baseline-Form),
`8a3aa1b0` (MR-002-Move), `f021780f` (Umsetzung: [ADR-0042](../plan/adr/0042-dauerhafte-lint-ausnahmen.md), Why-Kommentare,
Nachzug), `e04e126c` (Nachzug zweite Runde), `106f1114` (Sub-README gelöscht) und
`b290fc05` = `HEAD` (Plan-Tabelle berichtigt); die CR-Slice-Moves im Bereich
(`a1a5ae14`, `5a0fc507`) berühren den Gegenstand nur am Rand.

**Skill:** `v6.18.1` · `.harness/skills/reviewer.md` @ `b290fc05` · <!-- d-check:ignore (Adopter-spezifischer Skill-Pfad, existiert im Ziel-Repo ggf. nicht) -->
**Modell:** claude (glm-5.3-flash) · **Datum:** 2026-11-02

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- [Slice-Plan](../plan/planning/done/slice-harness-carveouts-baseline-form.md) (mit der berichtigten §3-Tabelle, Stand `b290fc05`)
- [ADR-0042](../plan/adr/0042-dauerhafte-lint-ausnahmen.md) (neu, `Accepted`),
  [ADR-0003](../plan/adr/0003-solid-nahes-lint-profil.md) (unverändert, `Accepted`),
  [ADR-0037](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md) (Kontext `make`-Carveout)
- [SPEC-031](../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows),
  [SPEC-034](../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln),
  [SPEC-035](../../spec/spezifikation.md#spec-035--pflichten-des-lint-profils),
  [LH-QA-004](../../spec/lastenheft.md#lh-qa-004--linting-solid-nahes-lint-profil)
- [MR-002](../../harness/conventions/done/MR-002-carveout-inventar.md) (aufgelöst),
  [MR-003](../../harness/conventions/MR-003-roadmap-wellen.md),
  [MR-008](../../harness/conventions/MR-008-adr-form.md),
  [MR-011](../../harness/conventions/done/MR-011-werkzeug-fragmente-repo-gefuehrt.md) (Muster)
- `AGENTS.md` (Hard Rules), `harness/conventions.md`, `harness/review.md`,
  [CO-001](../plan/carveouts/CO-001-crlf-logs.md) und
  [der Folge-Slice-Plan](../plan/planning/open/slice-v1-logs-crlf.md) als geprüfte Artefakte

---

## Findings

Jedes Finding folgt dem **§Output-Schema des Reviewer-Skills** — der
verbindlichen Single Source of Truth. Die Spalten unten sind nur
**gespiegelt** (Bequemlichkeit beim Ausfüllen), nicht neu definiert; bei
Abweichung gilt der Skill bzw. dessen Quelle
`v6.18.1` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der Plan-Diff-Abgleich (AGENTS.md §3.10) ist nach der letzten Fix-Runde unvollständig: drei Plan-Aussagen widersprechen dem Stand des Diffs. (a) Die §3-Zeile zu MR-002 sagt „reiner Move-Commit, Inhalt unverändert" — der Inhalt wurde nach dem Move in `f021780f` und `e04e126c` zweimal geändert (Geltungsbereich-Link, Adaption-Satz). (b) Die §1-Zuordnung sagt zur depguard-Zeile `!**/*_test.go` „Berührung im Slice: —" — [ADR-0042](../plan/adr/0042-dauerhafte-lint-ausnahmen.md) trägt sie als Entscheidung 3 und der depguard-`Why:`-Kommentar benennt sie. (c) Die §1-Zuordnung sagt „Warum-Kommentare nennen den ADR" für die acht Test-Carveout-Linter — der Diff benennt [ADR-0003](../plan/adr/0003-solid-nahes-lint-profil.md) nur je Block für sieben (`cyclop`, `gocognit`, `gocyclo`, `nestif`, `funlen`, `noctx`, `unparam`); der achte (`revive`/`unused-parameter`) trägt keinen `Why:`-Kommentar und bekam keinen. | `AGENTS.md` §3.10 · [BEO-DOK/plan-aussage-ueber-code](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md) | [Slice-Plan](../plan/planning/done/slice-harness-carveouts-baseline-form.md) · "reiner Move-Commit, Inhalt unverändert" | nein — Review-only; Lesung an `git show f021780f` bzw. gegen `b290fc05` | Plan-Aussage über Code/Berührung trägt dem Stand nicht (BEO-DOK/plan-aussage-ueber-code) |
| F-2 | MEDIUM | Der aufgelöste MR-002-Record beschreibt nach der Umstellung eine Adaption, die es so nicht gab: Geltungsbereich und Adaption benennen als die inventarisierte Ablage `docs/plan/carveouts/README.md` — die tatsächliche Ablage war `docs/plan/planning/in-progress/carveouts.md` (im Plan-Text link-frei in Backticks), und der Titel desselben Eintrags lautet unverändert „Carveout-Inventar an fester Stelle **statt `docs/plan/carveouts/`**". Die Pfad-Berichtigung nach dem Umzug („der `git mv` zieht die Pfad-Berichtigung nach sich") deckt Referenzen, nicht die Inhaltsaussage; link-freie Backtick-Nennungen der alten Stelle hätten stehen bleiben dürfen. | `v6.18.1` · `regelwerk/grundlagen-traceability.md` §Herkunfts-Anker (Gegenrichtung) · `AGENTS.md` §3.7 | [MR-002](../../harness/conventions/done/MR-002-carveout-inventar.md) · "Carveouts werden **inventarisiert** in der einen Datei" | nein — Review-only; kein Sensor liest die Adaption-Prosa von Konvention-Einträgen | Aufgelöster Konvention-Record inhaltlich umgeschrieben statt nur Pfad-berichtigt |
| F-3 | LOW | `harness/verification.md` sagt jetzt, neue temporäre Carveouts bräuchten „einen Eintrag in `docs/plan/carveouts/README.md` und einen Plan-Anker" — die Pflicht der CO-Form (erst die [CO-Datei](../plan/carveouts/CO-001-crlf-logs.md) aus der Vorlage; die Zeile des Index ist nur Pointer) fehlt im Wortlaut. | `v6.18.1` · `regelwerk/modul-07-carveouts.md` §Ziel-Form: Carveout (Sechs Pflicht-Header-Felder, Dateikonvention `CO-<NNN>`) | [verification.md](../../harness/verification.md) · "Neue temporaere Carveouts brauchen parallel einen Eintrag in" | nein — Review-only; ein Folge-Lauf, der nur die Index-Zeile anlegt, entginge dem Sensor (d-check prüft die CO-Datei-Pflicht nicht) | Pointer-Wortlaut verliert die CO-Datei-Pflicht |
| F-4 | LOW | `.harness/skills/closure-note-reviewer.md` liest für die „Gegenprobe der Carveout-Zeile" weiterhin „das Carveout-Inventar" als eine Datei (Ziel ist jetzt der Index `docs/plan/carveouts/README.md`); in der Baseline-Form ist der Index eine Pointersammlung — die Gegenprobe liest die CO-Dateien je Carveout. | `v6.18.1` · `regelwerk/modul-07-carveouts.md` §Ziel-Form: Carveout | [Skill](../../.harness/skills/closure-note-reviewer.md) · "das Carveout-Inventar" | nein — Review-only | Pointer-Wortlaut verliert die CO-Datei-Pflicht |
| F-5 | LOW | Der Verweis-Nachzug hat den href auf `README.md` gezogen, den Link-Text nicht: mehrere Stellen tragen als Beleg-Label den Text `carveouts.md`, das Ziel heißt README.md — eine Leserin, die dem Label folgt, sucht eine Datei, die es nicht gibt. Stellen: `docs/maintainer/branch-protection.md`, `.harness/skills/closure-note-reviewer.md` (beide Treffer), `docs/archive/roadmap-history-v0.1-v0.3.md`. | Maintainability (Konsumierbarkeit des Nachzugs) | [branch-protection.md](../../docs/maintainer/branch-protection.md) · "Carveout-Inventar: [`carveouts.md`](" | nein — Review-only (d-check prüft Pfad und Anker, nicht Link-Text) | Link-Label bleibt beim alten Dateinamen, Ziel wechselt |
| F-6 | LOW | Der in diesem Slice neu angelegte Folge-Slice-Plan [slice-v1-logs-crlf](../plan/planning/open/slice-v1-logs-crlf.md) hat keine Zeile im „Bestand" von `docs/plan/planning/open/README.md` — alle anderen offenen Slices dort tragen je eine; der Backlog-Bestand ist für dieses Artefakt unvollständig. (Auch die Roadmap nennt ihn nicht — dort trägt das Format-Regel, dass der Verzeichnis-Stand die Zustandsquelle ist; daran setze ich keinen zweiten Befund.) | Maintainability (Bestandstabelle des Backlogs) | [open/README.md](../../docs/plan/planning/open/README.md) · "## Bestand" | nein — Review-only; kein Sensor prüft Bestand-Vollständigkeit | Neu angelegtes open-Artefakt fehlt im Bestand-Index |
| F-7 | LOW | Im depguard-Regelblock von `.golangci.yml` trägt der neue Kommentar die Kennung als „ADR 0042" (ohne Bindestrich) — die Anzeigeform des Index; die zwei anderen neuen Kommentare in derselben Datei tragen die Schema-Form mit Bindestrich. Die Grepspur aller Anchor-Stellen der ADR in der Konfiguration greift so an zwei von drei Stellen. | `harness/conventions.md` [MR-000](../../harness/conventions.md) (ID-Schema `ADR-<NNNN>`) | `.golangci.yml` · "the named decision is ADR 0042" | nein — Review-only; kein Sensor scannt `.yml`-Kommentare auf Kennungs-Form | Kennung in einer Nicht-Markdown-Ablage trägt die Anzeige- statt der Schema-Form |
| F-8 | LOW | Die Guides-Zeile des Gate-Index beschreibt die neue Ablage als „temporäre und permanente Carveouts (Index und CO-Dateien)" — unter dem Trichter sind die dauerhaften Ausnahmen ADR-getragen ([ADR-0042](../plan/adr/0042-dauerhafte-lint-ausnahmen.md)), nicht CO-Dateien; der eigene [Index](../../docs/plan/carveouts/README.md) sagt „Dauerhafte Begründungen leben nicht hier". Jeder Lauf, der der Guides-Zeile folgt, sucht dauerhafte Carveouts am falschen Ort oder hält eine künftige dauerhafte Ausnahme verfrüht für die CO-Pflicht. | `v6.18.1` · `regelwerk/modul-07-carveouts.md` §Werkzeug-Wahl bei Diskrepanz | [harness/README.md](../../harness/README.md) · "temporäre und permanente Carveouts (Index und CO-Dateien)" | nein — Review-only; d-check prüft Existenz und Anker der Zeile, nicht ihre Trichter-Aussage | Guides-Zeile beschreibt die alte Aufteilung und kollidiert mit der Trichter-Spur |
| F-9 | INFO | In der ADR-Index-Zeile von [ADR-0042](../plan/adr/0042-dauerhafte-lint-ausnahmen.md) steht nach der SPEC-Klammer ein drittes schließendes Klammerzeichen („…production-only, …`regeln`)))"); der Link parst, `make docs-check` grün — Beobachtung ohne Änderungspflicht, der nächste Berührungslauf der Zeile kann sie mitziehen. | Maintainability (Anzeige-Genauigkeit des Index) | [ADR-Index](../plan/adr/README.md) Zeile 70 · <!-- d-check:ignore (Zitat der überzähligen Klammerung der Index-Zeile, kein Verweis dieses Reports) --> "…import-regeln)))" | nein — Review-only | Überzähliges Klammerzeichen in der Index-Zeile |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `.golangci.yml` — Gate-Lockerung (AGENTS.md §3.6) | geprüft, ohne Befund — der Slice-Diff der Datei trägt nur Kommentar-Zeilen; `linters`/`settings`/`depguard`/`exclude-rules`-Werte wortwörtlich unverändert (`git show f021780f -- .golangci.yml`: jede geänderte Zeile Kommentar) |
| Accepted-[ADR-0003](../plan/adr/0003-solid-nahes-lint-profil.md) — Unveränderlichkeit | geprüft, ohne Befund — kein Slice-Commit berührt `docs/plan/adr/0003-*.md` |
| `docs/plan/carveouts/` Ziel-Form gegen die Vorlagen | geprüft, ohne Befund — [CO-001](../plan/carveouts/CO-001-crlf-logs.md) trägt die sechs Pflicht-Header-Felder (Status, Datum angelegt, Letzte Prüfung, betroffenes Gate, Geltungsbereich, Folge-Slice), den beobachtbaren Auflösungs-Trigger („Erster bestätigter Befund …") und den vorhandenen Folge-Slice-Plan in `open/`; die Fundstellen (`splitLogLines` in `internal/adapter/driving/cli/logs.go`, Zeilen 309 bis 320; `docs/maintainer/cli-json-contract.md` §6.8) habe ich am Code gelesen und bestätigt; [README.md](../../docs/plan/carveouts/README.md) folgt der Vorlage, trägt „(noch keine)" ehrlich (es gab nie eine CO-Datei) und die Trichter-Konvention statt eine Quelldaten-Kopie |
| Trichter-Disziplin (dauerhafte Einträge direkt ADR-getragen, nicht als CO-Dateien geboren und gestorben) | geprüft, ohne Befund — [ADR-0042](../plan/adr/0042-dauerhafte-lint-ausnahmen.md) trägt das Verhältnis („Ergänzt, ersetzt nichts", Trichter-Zitat), der README-Konventionseintrag nennt die Anker-Klassen ([ADR-0003](../plan/adr/0003-solid-nahes-lint-profil.md), [ADR-0042](../plan/adr/0042-dauerhafte-lint-ausnahmen.md), [SPEC-031](../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows), [SPEC-034](../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln)) und die Werkzeug-Wahl; die `carveouts/`-Ablage trägt nur die temporäre Ausnahme |
| [ADR-0042](../plan/adr/0042-dauerhafte-lint-ausnahmen.md) gegen `v6.18.1` · `regelwerk/modul-04-adrs.md` §Ziel-Form: ADR (MADR) | geprüft, ohne Befund — fünf verglichene Optionen mit Pro/Contra inklusive „nichts tun" (Option E), Re-Evaluierungs-Trigger vorhanden (permanent mit Ereignis-Trigger), ADR-Index-Zeile aktualisiert (SPEC-Klammer in der zweiten Nachzug-Runde nachgezogen), Decken-Regel: kein Link von der ADR nach unten (Slice, Carveout, Roadmap) in der Datei; das `Schärft`-Feld deklariert die Aufwärts-Kopplung ([SPEC-035](../../spec/spezifikation.md#spec-035--pflichten-des-lint-profils), [SPEC-034](../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln)) |
| MR-002-Auflösung — Move-Reihenfolge | geprüft, ohne Befund — reiner Move-Commit (`8a3aa1b0`, `R100`), Berichtigungen als eigene Commits danach (`f021780f`, `e04e126c`); Aufgelöste-Tabelle in `harness/conventions.md` korrekt; der Anker `mr-002` existiert weiter (keine verwaiste `#mr-002`-Referenz im Baum) |
| Lifecycle-/Move-Disziplin der Slice-Commits | geprüft, ohne Befund — reine `git mv`-Commits für MR-002 und die CR-Vorschläge, Berichtigungs-Nachzug je separater Commit (AGENTS.md §3.3-Muster) |
| `.d-check.yml` neue Bindung | geprüft, ohne Befund — `matrix`-Klasse `carveout` auf `docs/plan/carveouts/CO-*.md`, `ids`-CO-Target auf `docs/plan/carveouts/`, keine verbleibende alte Pfad-Angabe in der Konfiguration |
| Verweis-Nachzug (zweite Runde, ~30 weitere Dateien) | geprüft, ohne Befund (Link-Ziele) — Suchlauf über den Baum nach `planning/in-progress/carveouts` ohne lebende Links; die verbleibenden Treffer sind link-freie historische Zitate im Slice-Plan und in [ADR-0042](../plan/adr/0042-dauerhafte-lint-ausnahmen.md) (Plan-Aussatz: historische Zitate bleiben, nur der href wandert) |
| Master-Inventar — Anker je Zeile | geprüft, ohne Befund (Zuordnung) — jede der acht dauerhaften Zeilen des früheren Inventars trägt nach der Plan-§1-Tabelle einen Anker; der DoD-Beleg „je Zeile mit Fundstelle in der Closure" ist Verifier-Aufgabe, nicht Gegenstand dieses Laufs |
| CR-Slice-Moves `a1a5ae14` / `5a0fc507` im Detail | nicht geprüft — andere Slice-Kennungen mit eigenen Berichten am selben Tag; hier nur als Bereichsstöber-Rand gewertet |

## Sensors

| Sensor | Stand |
|---|---|
| `make docs-check` | gelaufen per Implementer-Handoff (0 Befunde); für diesen Report nach dem Schreiben erneut (die Report-Datei trägt neue Links), Ausgang im Handoff-Text |
| `make baseline-verify` | gelaufen per Implementer-Handoff (`v6.18.1` OK); nicht wiederholt — keine Baseline-Bestandes-Assertion in den Findings |
| `make lint`, `make test`, `make gates` | nicht gelaufen — kein Produktcode-Diff; nur YAML-Kommentare und Markdown; Gate-Evidence ist Verifier-Aufgabe |
| `make test-docker`, `make verify-depguard` | nicht gelaufen — kein Docker-Pfad und keine Import-Regel berührt |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 6 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Plan-Aussage über Code/Berührung trägt dem Stand nicht (BEO-DOK/plan-aussage-ueber-code) · Aufgelöster Konvention-Record inhaltlich umgeschrieben statt nur Pfad-berichtigt · Pointer-Wortlaut verliert die CO-Datei-Pflicht · Link-Label bleibt beim alten Dateinamen, Ziel wechselt · Neu angelegtes open-Artefakt fehlt im Bestand-Index · Kennung in einer Nicht-Markdown-Ablage trägt die Anzeige- statt der Schema-Form · Guides-Zeile beschreibt die alte Aufteilung und kollidiert mit der Trichter-Spur · Überzähliges Klammerzeichen in der Index-Zeile

## Verdikt

**Merge-blockierend:** nein — keine HIGH-Findings. Die zwei MEDIUM-Findings sind
Berichtigt-Läufe (F-1: drei Plan-Zeilen nachziehen; F-2: den MR-002-Record
wieder inhaltlich treu zur alten Ablage zurücksetzen und die Berichtigung nur
auf die Links anwenden); sie blockieren die Slice-Closure nicht, lassen sich
aber in deren Fix-Runde vor der Closure-Notiz lösen (F-2 ist Plan-Rückkante —
der Plan-Aussatz „Textnachweise folgten derselben Ersetzung" schrieb die
übergreifende Regel).

**Übergabe:** Findings gehen an den Implementer (Rückkante
Review → Plan bei Plan-Defekt: F-1); die **Finding-Klassen** gehen zusätzlich
in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst
ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und
muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).