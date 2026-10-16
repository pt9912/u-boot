# Slice Harness: Baseline-Bump Regelwerk v6.13.0 → v6.18.1

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Phase:** Harness-/Baseline-Wartung (kein Produkt-Meilenstein, keine Welle).

**Bezug:** Baseline-Pin in
[`harness/conventions.md`](../../../../harness/conventions.md) Abschnitt
Baseline; [`MR-004`](../../../../harness/conventions/MR-004-regelwerk-vendored.md)
Bump-Prozedur und Freshness-Audit. Berührt `MR-001`..`MR-010` als
Gegenproben, ändert sie nicht (außer dem Pin-Pfad). Kein `LH`-Neubezug.

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel

Review-Bump des adoptierten Betriebsregelwerks (`pt9912/ai-harness-course`)
von `v6.13.0` auf `v6.18.1` als Einheit nach `MR-004` (Pin, Vendor-Pfad,
`AGENTS.md`, `harness/README.md`) samt Delta-Lektüre und `MR-*`-Gegenprobe.
Referenz für den Soll-Stand war ein frisch mit `ai-harness-init` erzeugtes
Repo auf `v6.18.1`; dessen `SHA256SUMS` ist mit dem hier vendorten identisch.

## 2. Verification-Evidence

- **Sensor:** `tools/harness/fetch-baseline-cache.sh --check-freshness` —
  vor dem Bump **Exit 3** (Pin `v6.13.0`, neuer: `v6.14.0`, `v6.14.1`,
  `v6.15.0`, `v6.16.0`, `v6.17.0`, `v6.18.0`, `v6.18.1`); nach dem Bump
  **Exit 0** (Pin ist der neueste Tag).
- **Sensor:** `tools/harness/fetch-baseline-cache.sh --verify` — **grün**
  (`54` Dateien, vollständig).
- **Sensor:** `make doc-check` — grün.
- **Skript-Fix:** `fetch-baseline-cache.sh` verglich Dateizahlen als
  Strings; `wc -l` liefert unter macOS führende Leerzeichen, der Vendor-Lauf
  brach mit „Kopierschritt hat Dateien übergangen“ ab. Die beiden Zählungen
  sind getrimmt.

## 3. Delta-Lektüre (v6.13.0 → v6.18.1)

Umfang: 36 von 54 Dateien geändert (26 `regelwerk/`, 10 `templates/`), keine
Datei hinzugekommen oder entfallen. Der Großteil der `regelwerk/`-Zeilen sind
Tag-Umschreibungen der Quell-Links.

- **Werkzeug-Fragmente:** `grundlagen-harness-dateien.md` führt
  `harness/mk/<werkzeug>.md` als werkzeug-eigenen Teil des Gate-Index
  (einer je Werkzeug, disjunkt, verlinkt; d-check ≥ v0.82.0/v0.83.0).
- **Spec-Stratum für Werkzeug-Festlegungen:** Was ein Gate prüft und wie es an
  Randformen entscheidet, steht in der Spezifikation (neuer Abschnitt
  „Festlegungen der Harness-Werkzeuge“); die Gate-ADR schärft diese Stelle.
- **Review:** Deckungs-Sensor meldet Leerlauf (`require-promises`); Report-
  Dateiname trägt die volle Slice-Kennung; `pfad` im Finding ist ein
  wörtliches Kurzzitat statt Zeilennummer; Reviewer-Skill: kein Stil-Polizist,
  kein HIGH/MEDIUM ohne Failure-Szenario, LOW nur mit Konventions-Anker.
- **Beobachtungs-Register:** Kennung `BEO-<KUERZEL>/<slug>`, Ablage als
  Verzeichnis `docs/plan/planning/observations/`.
- **Templates:** Pfade auf `.harness/baseline/<tag>/templates/…`, `.d-check.yml`
  mit `reviews`-Schlüsseln (`match: name`, `require-promises`, Archiv-Stubs).

## 4. MR-Gegenproben

| MR | Befund |
|---|---|
| `MR-000`/`MR-001` | **valide.** Drei obligatorische Straten bleiben; der neue Spec-Abschnitt für Werkzeug-Festlegungen ist additiv. |
| `MR-002`/`MR-003` | **valide.** Carveout- und Roadmap-Form unverändert. |
| `MR-004` | **Vollzug** durch diesen Slice (Pin, Vendor-Pfad, `AGENTS.md`, `harness/README.md`). |
| `MR-005` | **valide.** u-boots `.d-check.yml` setzt `reviews` nicht; die neuen Schlüssel erfordern d-check ≥ v0.86.0 (Pin: v0.51.1) und sind erst mit dem Folge-Slice relevant. |
| `MR-006`/`MR-007` | **valide.** |
| `MR-008` | **valide.** ADR-Vorlagen-Delta ist eine Ergänzung im `Schärft`-Feld (Gate-ADRs), keine Formänderung. |
| `MR-009` | **valide, Hinweis.** Die Reviewer-Skill-Vorlage schärft Finding-Anker und Kategorien (§3); Auffrischung von `.harness/skills/reviewer.md` ist offen. |
| `MR-010` | **valide.** `spezifikation.template.md` hat einen weiteren Abschnitt; `architecture.template.md` unverändert. |

BEO-Register: u-boot führt keines; die Kennungsänderung ist ohne Nachlauf.

## 5. Definition of Done

- [x] Vendor-Lauf `v6.18.1` (`54` Dateien) mit `--verify` grün.
- [x] Alten Bestand `.harness/baseline/v6.13.0/` entfernt.
- [x] Pin-Einheit nach `MR-004`: `conventions.md` (Stand, Adoptions-Historie,
  Integritäts-Pin, Quellen, Index-Spalte), `AGENTS.md`, `harness/README.md`.
- [x] Live-Verweise nachgezogen: `MR-001`..`MR-010`,
  `.harness/skills/reviewer.md`, `docs/reviews/README.md`,
  `docs/plan/planning/in-progress/roadmap.md`. Historische Nennungen in
  `done/` und `docs/reviews/` bleiben stehen.
- [x] `MR`-Gegenprobe dokumentiert (§4).
- [x] `make doc-check` grün.
- [x] Review durchgeführt: nicht erforderlich — reiner Vendor-/Pin-Bump ohne
  Produktcode; Gegenproben in §4.

## 6. Nicht-Ziele

- Keine Produktcode-Änderung.
- Keine Übernahme der Werkzeug-Schicht aus dem `ai-harness-init`-Referenzrepo
  (`harness/mk/`, `tools/harness/*.sh`, `.githooks/`, `.claude/`):
  [`slice-harness-werkzeug-schicht-ai-harness-init`](../in-progress/slice-harness-werkzeug-schicht-ai-harness-init.md).

## 7. Closure-Notiz

Vollzogen am 2026-10-16. Der Bump ist als Einheit ohne Adaptions-Konflikt
gelaufen; die Delta-Punkte aus §3 erzeugen keinen Pflichtbestand, außer der
offenen Auffrischung des Reviewer-Skills (§4, `MR-009`).
