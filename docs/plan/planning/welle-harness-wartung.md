# Welle welle-harness-wartung: Harness-Wartung (Reviewer-Skill und Paket-READMEs)

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/`. Der Zustand ist die Verzeichnis-Position — kein Status-Feld.
**Geplante Wellen bekommen noch keine Datei:** Sie stehen in der Roadmap unter
*Nächste Wellen* und nirgends sonst — zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-10-16.

---

## 1. Welle-Ziel

Der Reviewer-Skill entspricht der Vorlage der Baseline `v6.18.1`, und die Festlegungen aus den
`internal/`-READMEs stehen in der Spezifikation statt in den Paket-READMEs. Der zweite Slice wird mit dem
aufgefrischten Skill reviewt: Der Review-Report des Spezifikations-Slice trägt damit das Schema der Vorlage
(Datei als `pfad`, LOW nur mit Konventions-Anker, kein HIGH/MEDIUM ohne Failure-Szenario).

## 2. Trigger (Welle startet)

- Die Welle [`welle-macos-verifikation`](done/welle-macos-verifikation.md) liegt in `done/` (geschlossen 2026-10-16).
- Beide Slices der Welle liegen in `next/` und tragen einen `Verantwortlich:`-Eintrag.

## 3. Closure-Trigger (Welle schließt)

Das *Mehr* gegenüber den beiden Slice-DoDs: Die Reihenfolge trägt. Der Review des zweiten Slice läuft mit dem
Skill, den der erste geliefert hat.

- Beide Slices der Welle liegen in `done/`.
- Der Review-Report des Slice [`slice-harness-internal-readmes-festlegungen-in-spec`](done/slice-harness-internal-readmes-festlegungen-in-spec.md)
  ist nach dem Abschluss von [`slice-harness-dcheck-bump-reviewer-skill`](done/slice-harness-dcheck-bump-reviewer-skill.md)
  entstanden (Commit-Reihenfolge belegt es).
- `make doc-trace` meldet nach der Verlagerung keine Waise unter den Anforderungen, die vorher über die README-Liste Service → Anforderung belegt waren.
- `make gates` ist grün.

## 4. Slices in dieser Welle

| Slice | Titel | Bezug |
|---|---|---|
| [`slice-harness-dcheck-bump-reviewer-skill`](done/slice-harness-dcheck-bump-reviewer-skill.md) | Auffrischung des Reviewer-Skills (Baseline v6.18.1) | [`MR-009`](../../../harness/conventions/MR-009-skills-und-review-ablage.md) |
| [`slice-harness-internal-readmes-festlegungen-in-spec`](done/slice-harness-internal-readmes-festlegungen-in-spec.md) | Festlegungen aus den `internal/`-READMEs in die Spezifikation | [`LH-LESE-003`](../../../spec/lastenheft.md#lh-lese-003--dokumentenordnung) |

## 5. Abhängigkeiten

- Blockiert: keine Welle.
- Wird blockiert von: keiner Welle. Innerhalb der Welle geht der Skill-Slice dem Spezifikations-Slice voraus (Closure-Trigger).

## 6. Out-of-Scope für diese Welle

- Der Bump von d-check auf `v0.86.1`: ihn führte [`slice-harness-ai-harness-init-adoption`](done/slice-harness-ai-harness-init-adoption.md) aus.
- Änderungen an Code oder Werten: der Spezifikations-Slice dokumentiert den Ist-Stand, ein abweichender Wert wird vermerkt, nicht geändert.
- Neue Produktanforderungen: Festlegungen verfeinern vorhandene Anforderungen ([`LH-LESE-003`](../../../spec/lastenheft.md#lh-lese-003--dokumentenordnung)).
- Die Smoke-Härtung des Homebrew-Taps ([`slice-v2-homebrew-smoke-haertung`](open/slice-v2-homebrew-smoke-haertung.md)) und die Docker-Desktop-Messung:
  andere Gegenstände, ohne Bezug zum Harness.

## 7. Closure-Notiz

Offen.
