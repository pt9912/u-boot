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

- Die Welle [`welle-macos-verifikation`](welle-macos-verifikation.md) liegt in `done/` (geschlossen 2026-10-16).
- Beide Slices der Welle liegen in `next/` und tragen einen `Verantwortlich:`-Eintrag.

## 3. Closure-Trigger (Welle schließt)

Das *Mehr* gegenüber den beiden Slice-DoDs: Die Reihenfolge trägt. Der Review des zweiten Slice läuft mit dem
Skill, den der erste geliefert hat.

- Beide Slices der Welle liegen in `done/`.
- Der Review-Report des Slice [`slice-harness-internal-readmes-festlegungen-in-spec`](slice-harness-internal-readmes-festlegungen-in-spec.md)
  ist nach dem Abschluss von [`slice-harness-dcheck-bump-reviewer-skill`](slice-harness-dcheck-bump-reviewer-skill.md)
  entstanden (Commit-Reihenfolge belegt es).
- `make doc-trace` meldet nach der Verlagerung keine Waise unter den Anforderungen, die vorher über die README-Liste Service → Anforderung belegt waren.
- `make gates` ist grün.

## 4. Slices in dieser Welle

| Slice | Titel | Bezug |
|---|---|---|
| [`slice-harness-dcheck-bump-reviewer-skill`](slice-harness-dcheck-bump-reviewer-skill.md) | Auffrischung des Reviewer-Skills (Baseline v6.18.1) | [`MR-009`](../../../../harness/conventions/MR-009-skills-und-review-ablage.md) |
| [`slice-harness-internal-readmes-festlegungen-in-spec`](slice-harness-internal-readmes-festlegungen-in-spec.md) | Festlegungen aus den `internal/`-READMEs in die Spezifikation | [`LH-LESE-003`](../../../../spec/lastenheft.md#lh-lese-003--dokumentenordnung) |

## 5. Abhängigkeiten

- Blockiert: keine Welle.
- Wird blockiert von: keiner Welle. Innerhalb der Welle geht der Skill-Slice dem Spezifikations-Slice voraus (Closure-Trigger).

## 6. Out-of-Scope für diese Welle

- Der Bump von d-check auf `v0.86.1`: ihn führte [`slice-harness-ai-harness-init-adoption`](slice-harness-ai-harness-init-adoption.md) aus.
- Änderungen an Code oder Werten: der Spezifikations-Slice dokumentiert den Ist-Stand, ein abweichender Wert wird vermerkt, nicht geändert.
- Neue Produktanforderungen: Festlegungen verfeinern vorhandene Anforderungen ([`LH-LESE-003`](../../../../spec/lastenheft.md#lh-lese-003--dokumentenordnung)).
- Die Smoke-Härtung des Homebrew-Taps ([`slice-v2-homebrew-smoke-haertung`](../open/slice-v2-homebrew-smoke-haertung.md)) und die Docker-Desktop-Messung:
  andere Gegenstände, ohne Bezug zum Harness.

## 7. Closure-Notiz

Detailquellen: die Closure-Notizen der beiden Slices
([`slice-harness-dcheck-bump-reviewer-skill`](slice-harness-dcheck-bump-reviewer-skill.md) §7,
[`slice-harness-internal-readmes-festlegungen-in-spec`](slice-harness-internal-readmes-festlegungen-in-spec.md) §7); die Welle ist
kein Release, einen Release-Cut-Slice gibt es nicht.

- **Geliefert:** Der Reviewer-Skill und `harness/review.md` tragen das Schema der Baseline `v6.18.1` (`pfad` als Kurzzitat,
  `klasse`, Failure-Szenario-Pflicht, LOW mit Konventions-Anker). Sieben Festlegungen aus den `internal/`-READMEs stehen in der
  Spezifikation (`SPEC-036`, `LH-FA-UP-001.a`, `SPEC-038` bis `SPEC-042`), die READMEs sind von 490 auf 109 Zeilen gekürzt. <!-- d-check:ignore (Kennungen als Zitat) -->
- **Closure-Trigger (§3):** Beide Slices liegen in `done/`. Der Review-Report des Spezifikations-Slice (`fceea11`) entstand nach dem
  Move des Skill-Slice (`c9b1f4e`), also mit dem aufgefrischten Skill. `make doc-trace`: 83 Anforderungen, 0 Waisen. `make gates` grün.
- **Was hat funktioniert:** Die Reihenfolge trug: Der Review mit dem neuen Skill fand einen echten Spec-Fehler (`SPEC-039`), <!-- d-check:ignore (Kennungen als Zitat) -->
  den die Inventur nicht gefunden hatte.
- **Was ging anders als geplant:** Die Inventur deckte eine Abweichung Lastenheft gegen Code bei der Verbosity auf, die nicht zur Welle
  gehört. Der Skill-Slice nahm drei HIGH-Anker über seine DoD hinaus mit (die DoD ist nachgezogen).
- **Trigger-Audit:** Carveout: keiner berührt, `carveouts.md` unverändert; 0 offen. Bootstrap-aware Gate: `coverage-gate` 90,90 % gegen
  90 %, unverändert; 0 offen. ADR: keine ADR trägt einen eingetretenen Re-Evaluierungs-Trigger; 0 offen.
- **Steering-Loop-Einträge:** Kein Eintrag über der Schwelle (das Beobachtungs-Register ist leer). Gezählt, nicht verkörpert:
  Skill ⊃ `harness/review.md` (die älteren Skill-HIGHs stehen nur im Skill), und eine Spec, die ein README statt den Code abschreibt.
- **Beobachtungs-Register:** keine Beobachtung angelegt.
- **Folge-Slices:** [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../open/slice-v1-cli-verbosity-letzte-option-gewinnt.md) (Code),
  [`slice-cr-cli-quiet-fehler-und-warnungen`](../open/slice-cr-cli-quiet-fehler-und-warnungen.md) (CR am Lastenheft),
  [`slice-v2-homebrew-smoke-haertung`](../open/slice-v2-homebrew-smoke-haertung.md) (aus der Vorwelle, unverändert offen).
- **Verifikation:** `make gates` Exit 0; Verification Evidence je Slice unter `docs/reviews/`
  ([Skill-Slice](../../../reviews/2026-10-16-verification-slice-harness-dcheck-bump-reviewer-skill.md),
  [Spezifikations-Slice](../../../reviews/2026-10-16-verification-slice-harness-internal-readmes-festlegungen-in-spec.md)).
