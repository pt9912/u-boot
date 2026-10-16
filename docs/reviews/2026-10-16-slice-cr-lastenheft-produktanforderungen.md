# Review-Report: slice-cr-lastenheft-produktanforderungen — 2026-10-16

**Review-Art:** Plan | Design — Diff gegen Plan, Decken-Regel und Hard Rules (Vertrags-Stratum; kein Code).

**Gegenstand:** Diff `938df12..HEAD`, Schwerpunkt Commits `167ee32` (Lastenheft 0.6.0), `8f3bc0a` (Spezifikation), `baf87c7` (lebende Verweise), `a5ae0b0` (Slice-Stand).

**Skill:** `.harness/skills/reviewer.md` @ HEAD `a5ae0b0` · <!-- d-check:ignore (Skill-Pfad als Zitat) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-16

**Eingangs-Kontext:**

- `slice-cr-lastenheft-produktanforderungen` (inkl. Anhang Zuordnung) <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
- Lastenheft vor und nach dem CR (`git show 938df12:spec/lastenheft.md`), Spezifikation, Architektur-Sicht
- [`ADR-0002`](../plan/adr/0002-hexagonale-architektur.md), [`ADR-0013`](../plan/adr/0013-dokumentationsreferenzmodell.md) (Immutabilität, Links)
- [`LH-LESE-003`](../../spec/lastenheft.md#lh-lese-003--dokumentenordnung), [`LH-NFA-PORT-004`](../../spec/lastenheft.md#lh-nfa-port-004--ausgeliefertes-container-image), [`LH-LESE-004`](../../spec/lastenheft.md#lh-lese-004--prüfbarkeit)
- `AGENTS.md` (Hard Rules)

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Die Verfeinerung zur CI-Pipeline nennt zweimal denselben Eltern-Verweis: „in SPEC-008 (MVP) und SPEC-008 (V1) definiert“ und „unter SPEC-008 (V1)“. Die Umnummerierung hat zwei verschiedene Eltern (Build-Targets MVP/V1) auf eine Kennung zusammengeführt; die MVP/V1-Unterscheidung hängt jetzt an einer Kennung ohne Unterscheidungsmerkmal. | Maintainability | `spec/spezifikation.md` · Abschnitt „Komposition der CI-Pipeline“, Zeile „Drei Jobs, parallel“ | nein — Review-only | Automatische Kennungs-Ersetzung erzeugt Doppel-Verweis | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| F-2 | LOW | Die Architektur-Sicht trägt `**Bezug:** SPEC-012..SPEC-034` als Bereich; er umfasst alle Kennungen dazwischen, gemeint sind drei (012, 013, 034). Die Zeile „Letzte Änderung: 2026-10-02“ ist trotz geänderter Datei nicht nachgezogen. | Maintainability | `spec/architecture.md` · Zeile 5 und 7 | nein — Review-only | Bereichs-Notation über unzusammenhängende Kennungen | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| F-3 | LOW | Die Einleitung zu §16 Historie begründet „das Lastenheft verweist nie abwärts auf Planung“ weiter mit einem Link auf die entfallene Referenzmodell-Anforderung (Grabstein). Eine geltende Regel wird mit einer entfallenen belegt; die Regel selbst steht nur noch im Baseline-Regelwerk. | [`LH-LESE-003`](../../spec/lastenheft.md#lh-lese-003--dokumentenordnung) | `spec/lastenheft.md` · §16, Absatz „Die Spalte Verweis“ | nein — Review-only | Lebender Verweis auf Grabstein als Geltungsgrund |
| F-4 | LOW | Grabstein-Texte der Bau-Anforderungen sagen „steht im technischen Dokument“, die der Architektur-Anforderungen „in der Architektur-Sicht und im technischen Dokument“. Das sind Namensnennungen ohne Link, aber Abwärtszeiger in die Ränge 2 und 3, während derselbe Absatz festhält, das Lastenheft verweise nicht auf das technische Dokument. Die Grabsteine nennen zudem keine Kennung; wer über einen alten Anker kommt, findet den neuen Ort nur über den Slice-Anhang. | [`LH-LESE-003`](../../spec/lastenheft.md#lh-lese-003--dokumentenordnung), `v6.18.1` · `regelwerk/grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP) | `spec/lastenheft.md` · Grabsteine unter §4.11 und §4.13 | nein — Review-only | Abwärtszeiger in Grabstein-Prosa |
| F-5 | LOW | Zwei ersetzte Verweise tragen die Regel nicht: `docs/archive/README.md` verweist für die Archiv-Konvention (Verschieben nach `docs/archive/`, Hinweiszeile) auf `modul-05-planning-harness.md`; dort steht nur die Wellen-Archivierung (Stubs, ZIP), keine `docs/archive/`-Konvention (`grep docs/archive` im Regelwerk ohne Treffer). `roadmap.md` und `planning/README.md` verweisen für die `slice-`/`tranche-`-Namenskonvention auf §Lifecycle als State Machine; die Abschnittsüberschrift existiert, die Namenskonvention steht dort nicht. | Maintainability | `docs/archive/README.md` · „Baseline-Regelwerk `modul-05-planning-harness.md`“; `docs/plan/planning/in-progress/roadmap.md` · „einzige zulässige Ausnahme“ | nein — Review-only | Ersatz-Verweis trägt die ersetzte Aussage nicht |
| F-6 | LOW | Lebende Stellen zitieren weiter entfallene Anforderungen als Geltungsgrund: `Dockerfile`, `repo.mk`, `.golangci.yml`, `.dockerignore` (Kommentare) und die `depguard`-`desc`-Zeilen in `spec/architecture.md` §4. Für Kommentare ist das im Plan bewusst zugelassen; die Architektur-Sicht ist aber ein lebendes Dokument, dessen `Bezug`-Zeile in demselben Diff auf die neuen Kennungen umgestellt wurde. Der Grabstein führt keinen Zeiger auf den neuen Ort. | AGENTS.md §3.7 (Rang-Zeiger), Slice-DoD „Lebende Verweise nachgezogen“ | `spec/architecture.md` · Zeilen 357–426; `Dockerfile` · Zeilen 6–22 | ja — `grep` nach den 22 Kennungen außerhalb `done/`, Reviews, Accepted ADRs | Zeiger auf Grabstein ohne Weiterleitung |
| F-7 | INFO | Die Traceability-Matrix zählt die 22 Grabsteine weiter als Anforderungen (83, 0 Waisen) und belegt sie über den CR-Slice selbst, der sie nennt; „0 Waisen“ trägt damit für diese 22 keine Aussage. `LH-NFA-PORT-004` steht nicht in der Matrix (Familienfilter nur `LH-FA-*`/`LH-QA-*`). | Maintainability | `.d-check.yml` · Abschnitt `trace`; `make doc-complete` | ja — `make doc-complete` | Matrix zählt Grabsteine als belegt | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| F-8 | INFO | Der Diff enthält über den CR hinaus Planänderungen an `slice-harness-wellen-betrieb` (Entscheidung „Archivierung nein“, Umbau von Ziel und DoD). Beleg ist die Notiz „Entscheidung (2026-10-16)“ im Slice selbst; der CR-Slice verweist darauf nur als Nicht-Ziel. Zuständigkeit des Verifiers, den Rollenwechsel zu prüfen. | Maintainability | `docs/plan/planning/open/slice-harness-wellen-betrieb.md` | nein — Review-only | Scope-Mischung im Commit-Bereich | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Vertragsintegrität: alle 22 Grabstein-Überschriften | geprüft, ohne Befund: der Vergleich der Überschriftenlisten `938df12` gegen HEAD (`##`–`####`) zeigt als einzige Abweichungen die zwei neuen Überschriften (Prüfbarkeit, Ausgeliefertes Container-Image); alle Anker bleiben gültig |
| Grabstein-Text | geprüft, ohne Befund: kein Link auf ADR, Slice, Carveout, Roadmap oder `SPEC-*`; Ausnahme siehe F-4 (Namensnennung) |
| Produktzusagen in den 22 Anforderungen | geprüft, ohne Befund: Runtime-Image-Eigenschaften leben in der neuen Anforderung; Pin-Politik, Mindest-Testumfang, CI, Lint, Layout, Schichten sind Prozess- bzw. Entwurfsvorgaben ohne Nutzerzusage |
| Deckungsgleichheit neue Image-Anforderung mit entfallener Runtime-Stage-Anforderung | geprüft, ohne Befund: minimal, shell-los (keine Shell), Non-root, OCI-Labels, keine Build-Toolchain, `ENTRYPOINT` auf das Binary; Priorität MVP und Einordnung in §5.4 (Portabilität, neben dem Container-Betrieb) stimmig; Bezug auf die Paketierung (GHCR-Image als primärer Weg) trägt |
| Lastenheft-Konsistenz: Kopf 0.6.0, §16-Zeile, §12 MVP-Umfang, §9/§10/§11/Glossar | geprüft, ohne Befund: Format der §16-Zeile wie die Vorgängerzeilen, Verweis „Vereinbarung mit dem Projektinhaber“; kein verbleibender Verweis auf Prozess-Inhalt außerhalb der Grabsteine; Wortlaut der beiden neuen Lesehinweise geprüft |
| Spezifikation: 13 Verfeinerungen | geprüft, ohne Befund: Zuordnung Anhang gegen `spec/spezifikation.md` stimmig (11 neue `SPEC-025`..`SPEC-035`, 2 Suffix-Verfeinerungen der neuen Image-Anforderung); alte Anker bleiben als HTML-Anker erhalten, der Slug entspricht der früheren Überschrift; Kennungs-Regel im Kopf deckt Suffix-Form (nur die Image-Anforderung) und `SPEC-<NNN>` ab; Eltern-Verweise (Makefile-Anforderungen auf Makefile-Verfeinerung, Repo-Layout auf Go-Layout, Schichten auf Layout-Verfeinerung) inhaltlich passend; keine zirkulären Selbstverweise; Ausnahme F-1 | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| Lebende Verweise (Stichprobe) | geprüft, ohne Befund: `AGENTS.md` §3.8, beide READMEs, `docs/maintainer/quality.md`, `branch-protection.md`, Planning-READMEs, `carveouts.md`, Skills, Sensor-Dateien, MR-Einträge; Baseline-Stellen `modul-05-planning-harness.md` §Lifecycle als State Machine, `modul-07-carveouts.md` und `grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP) existieren in `v6.18.1`; keine Historie-Erzählung auf spätere `SPEC-*` umgebogen (die Zeile zu `slice-m3-depguard-aktivierung-verifizieren` im Audit-Trail nennt die heutige Fundstelle der Regel, nicht die damalige); Ausnahmen F-5, F-6 | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| Accepted ADRs | geprüft, ohne Befund: `git diff 938df12..HEAD -- docs/plan/adr/` ändert nur `README.md`; die ADR-Links auf Grabstein-Anker bleiben gültig |
| Gates | geprüft, ohne Befund: `make docs-check` (234 Dateien, 0 Befunde), `make doc-complete` (83 Anforderungen, 0 Waisen), `make doc-planning` und `make doc-targets` (0 Befunde), jeweils per Docker; Hinweis zur Matrix siehe F-7 |
| Code, `make gates`, `make test-docker`, `make verify-depguard` | nicht geprüft: Diff ändert keinen Code und keine Importregeln; `make gates` ist Nachweis des Implementers bzw. Aufgabe des Verifiers |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 6 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Automatische Kennungs-Ersetzung erzeugt Doppel-Verweis · Bereichs-Notation über unzusammenhängende Kennungen · Lebender Verweis auf Grabstein als Geltungsgrund · Abwärtszeiger in Grabstein-Prosa · Ersatz-Verweis trägt die ersetzte Aussage nicht · Zeiger auf Grabstein ohne Weiterleitung · Matrix zählt Grabsteine als belegt · Scope-Mischung im Commit-Bereich

## Verdikt

**Merge-blockierend:** nein — kein HIGH/MEDIUM; die sechs LOW-Findings sind einmalige Inkonsistenzen ohne Vertragsbruch, die Vertragsintegrität (Anker, Produktzusage, Decken-Regel in den Grabsteinen als Link-Verbot, Accepted ADRs) hält.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen zusätzlich in die Slice-Closure §7. Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation — DoD-Konformität prüft der Verifier.
