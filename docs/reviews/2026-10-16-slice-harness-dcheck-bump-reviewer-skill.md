# Review-Report: slice-harness-dcheck-bump-reviewer-skill — 2026-10-16

**Review-Art:** Code (Diff gegen Plan + Konventionen; Doku-/Skill-Änderung)

**Gegenstand:** Commit `cf2f3e7` (`.harness/skills/reviewer.md`, `harness/review.md`, `docs/plan/planning/in-progress/roadmap.md`)

**Skill:** `.harness/skills/reviewer.md` @ cf2f3e7 (Stand nach dem Commit, bewusst)
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-16

**Eingangs-Kontext:**

- [`slice-harness-dcheck-bump-reviewer-skill`](../plan/planning/in-progress/slice-harness-dcheck-bump-reviewer-skill.md) §1–§3, §6 (Maßstab: Plan, nicht DoD)
- `MR-009` (Reviewer-Skill repo-gepflegt)
- Reviewer-Vorlage `v6.18.1` · `templates/.harness/skills/reviewer.template.md` und `templates/docs/reviews/review-report.template.md`
- `AGENTS.md` (Hard Rules §3.2, §3.7)
- Keine `LH-*`-IDs berührt (Skill-Pflege, keine Anforderungsänderung)

Sensors: nicht gelaufen. Reiner Doku-/Skill-Diff ohne Docker-Pfad und ohne Importregeln, daher kein `make test-docker`/`make verify-depguard`. `make gates` ist Gegenstand der DoD und damit des Verifiers.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die neue Regel „kein HIGH/MEDIUM ohne Failure-Szenario“ und die neuen HIGH-Anker für Kommentar-/Zustandsfeld-Hygiene stehen unaufgelöst nebeneinander: Ein Kommentar, der die verworfene Alternative beschreibt, hat per se kein konkretes Versagen, der Anker zwingt aber HIGH (blockierend). Risiko: Zwei Sessions stufen denselben Kommentar als HIGH bzw. LOW ein — genau die Drift, die der Skill verhindern soll. | `AGENTS.md` §3.7; Maintainability | `harness/review.md` · "Kein HIGH- und kein MEDIUM-Finding ohne Failure-Szenario" (Gegenstück: "Ein Kommentar beschreibt die verworfene Alternative"); ebenso `.harness/skills/reviewer.md` · "Kein HIGH- oder MEDIUM-Finding ohne Failure-Szenario" | nein — Review-only | Kategorien-Regel kollidiert mit Anker ohne Versagensbild |
| F-2 | LOW | Skill und `harness/review.md` sind nach dem Commit nicht deckungsgleich: Inline-`//nolint` als HIGH-Anker und die Drift-Log-Klausel des Zustandsfeld-Ankers stehen nur im Skill, `review.md` führt beides nicht. Risiko: Der Skill verlangt selbst, beide dürften nicht auseinanderlaufen. Die Lücke ist bei weiteren Skill-HIGHs (Dual-Classifier, Referenzmodell u. a.) älter; der Commit vergrößert sie um zwei Stellen und beteuert im Plan „Vorlage übernehmen nur, was review.md nicht widerspricht“ (Widerspruch geprüft: keiner, aber Ergänzung nicht gespiegelt). | `.harness/skills/reviewer.md` §Pflege; `MR-009` | `.harness/skills/reviewer.md` · "Inline-`//nolint` statt zentraler Ausnahme in" und "protokolliert Schließungen und erreichte Meilensteine"; `harness/review.md` · "traegt Chronik statt Zustand und" | nein — Review-only | Skill/review.md-Klassifikation divergiert |
| F-3 | LOW | Das Report-Gerüst `review-report.template.md` hat keine Spalte `Risiko`, das Skill-Schema (`risiko`) und `harness/review.md` (`Risiko:`) verlangen das Feld. Das Gerüst spiegelt `klasse` und `pfad` korrekt, aber nicht das repo-eigene Pflichtfeld. Risiko: Der Report-Autor streicht `risiko` oder baut eine eigene Spalte; der Skill gilt laut Gerüst als verbindlich, das Gerüst ist nicht deckungsgleich. | `MR-009` (Skill repo-gepflegt; Vorlage kein Pflichtstand) | `.harness/skills/reviewer.md` · "`risiko`: warum das relevant ist" | nein — Review-only | Report-Gerüst spiegelt Repo-Schema unvollständig |
| F-4 | LOW | Vier zusätzlich übernommene HIGH-Anker (Kommentar, Zustandsfeld, Norm-nur-im-Template-Kommentar, Inline-`//nolint`) stehen nicht in der DoD-Aufzählung der Vorlagen-Änderungen (`pfad`+Kurzzitat, LOW-Anker, Kein Stil-Polizist, Failure-Szenario). Auch `harness/review.md` wurde um zwei davon erweitert, Plan §3 nennt dort nur „Output-Schema `pfad` und Kategorien-Regel“. Inhaltlich sind sie in `AGENTS.md` §3.2/§3.7 verankert (kein Widerspruch zum Repo); formal ist es ein Scope-Zuwachs. Der §6-Risiko-Ausgang ist offen. | Slice §2/§3; Maintainability | `docs/plan/planning/in-progress/slice-harness-dcheck-bump-reviewer-skill.md` · "trägt die Vorlagen-Änderungen:" | nein — Review-only | Scope über DoD-Aufzählung hinaus | <!-- d-check:ignore (Pfad als Zitat) -->
| F-5 | INFO | Zwei neue Zeilen in `harness/review.md` brechen die Ohne-Umlaute-Konvention der Datei: `§` in der Kommentar-Anker-Zeile und `·` in der Schema-Zeile. Die Datei trug schon vorher ein `—` (Negativbefund-Block), die Konvention war also schon locker. Kein Konventions-Anker (kein Linter), daher kein LOW. | Maintainability | `harness/review.md` · "<CATEGORY> <path> · " und "(`AGENTS.md`](../AGENTS.md) §3.7" | ja — kein Gate; Review-only | Zeichensatz-Konvention review.md |
| F-6 | INFO | Der Commit streicht in `roadmap.md` die Zeile „Nichts in Arbeit.“; die Datei steht nicht in Plan §3. Die Änderung ist inhaltlich korrekt (Welle `welle-harness-wartung` ist gelistet), gehört aber zur Wellen-Eröffnung, nicht zum Skill-Abgleich. | Slice §3 | `docs/plan/planning/in-progress/roadmap.md` · "- [`welle-harness-wartung`](../plan/planning/welle-harness-wartung.md)" | nein — Review-only | Datei außerhalb Plan §3 |
| F-7 | INFO | Zuständigkeit Verifier: Slice §2 (DoD-Häkchen, `make gates`, Closure-Notiz §7 „Offen.“, §6-Ausgang „offen bis zur Closure“) ist nicht Gegenstand dieses Reviews. | Rollenvertrag | `docs/plan/planning/in-progress/slice-harness-dcheck-bump-reviewer-skill.md` · "Offen." | nein | Zuständigkeit Verifier | <!-- d-check:ignore (Pfad als Zitat) -->

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Template-Hinweis-Block und Platzhalter (`<Repo-Name>`, `<ADR-NNNN…>`, `<agent-review-Make-Target>`) im Skill | geprüft, ohne Befund: nicht übernommen, Kopf des Skills unverändert repo-spezifisch |
| Widerspruch Vorlagen-Aussage zu `harness/review.md` (LOW-Anker, Kein Stil-Polizist, `pfad`-Kurzzitat, Failure-Szenario) | geprüft, ohne Befund: in beiden Dateien inhaltlich gleichlautend |
| `klasse` / `Klasse:` in Skill und `review.md` | geprüft, ohne Befund: gleiche Bedeutung, beide verweisen auf den Steering-Loop-Zähler (Skill §Pflege, `review.md` „siehe unten“) |
| Skill-Schema gegen Report-Gerüst (`pfad` Kurzzitat, `klasse`, Negativbefund, Summary) | geprüft; Abweichung nur F-3 (`Risiko`) |
| Anker-Verweise (`AGENTS.md` §3.2, §3.7, Baseline §Template-Schichtung, §Was ein Kommentar trägt) | geprüft, ohne Befund: Abschnitte existieren in `grundlagen-harness-dateien.md` |
| Inline-`//nolint`-Anker gegen `AGENTS.md` §3.2 | geprüft, ohne Befund: deckt sich mit der Hard Rule (zentrale Ausnahme statt Inline); Doppelung mit bestehendem `//nolint:depguard` nur redundant |
| Kommentar-Hygiene (`AGENTS.md` §3.7) in den neuen Zeilen | geprüft, ohne Befund |
| `closure-note-reviewer.md`, `MR-009`-Doku | nicht geprüft — nicht im Diff; Plan §1 nennt Mitziehen nur bei Pfad-Anker-Bedarf, Doku-Update §2 ist Verifier-Frage |
| `make gates`, `make test-docker`, `make verify-depguard` | nicht gelaufen — Doku-only Diff, kein Docker-Pfad, keine Importregel; `make gates` gehört zur DoD (Verifier) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 3 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Kategorien-Regel kollidiert mit Anker ohne Versagensbild · Skill/review.md-Klassifikation divergiert · Report-Gerüst spiegelt Repo-Schema unvollständig · Scope über DoD-Aufzählung hinaus · Zeichensatz-Konvention review.md · Datei außerhalb Plan §3 · Zuständigkeit Verifier

## Verdikt

**Merge-blockierend:** nein — kein HIGH. Das MEDIUM (F-1) ist vor Closure zu klären (Auflösung der Spannung zwischen Failure-Szenario-Regel und den Kommentar-/Zustandsfeld-Ankern in Skill und `review.md`, Entscheidung beim Implementer).

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7. Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation.
