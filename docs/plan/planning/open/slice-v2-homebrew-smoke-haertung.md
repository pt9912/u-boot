# Slice v2: Härtung des Homebrew-Smoke im Tap

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — keine Closure-Bedingung, die von der DoD dieses Slice
verschieden ist (Nachzug zu Review-Findings am Tap-Smoke).

**Bezug:** [`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung)
(Paketierung), [`ADR-0016`](../../adr/0016-homebrew-distribution-per-tap.md)
(Homebrew-Distribution per Tap). Herkunft: Findings F-4, F-5 und F-8 im
[Review welle-macos-verifikation](../../../reviews/2026-10-16-welle-macos-verifikation.md)
und die
[Verification von `slice-v2-homebrew-formula`](../../../reviews/2026-10-16-verification-slice-v2-homebrew-formula.md);
Akzeptanzkriterium 4 in
[`slice-v2-homebrew-formula`](../in-progress/slice-v2-homebrew-formula.md).

**Berührte Spec-Stellen:** —

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Homebrew-Smoke im Tap erfüllt das Akzeptanzkriterium 4 von
[`slice-v2-homebrew-formula`](../in-progress/slice-v2-homebrew-formula.md) voll:
`brew install --build-from-source` (bzw. gleichwertig), `u-boot add` im Smoke und
ein gatendes `brew audit --strict` ohne `|| true`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Eine Änderung von [`ADR-0016`](../../adr/0016-homebrew-distribution-per-tap.md):
  die ADR ist Accepted und immutable; F-8 betrifft nur den Kommentar im Workflow,
  nicht die Entscheidung.
- Neue Plattformen oder Distributionswege im Tap: der Slice härtet den
  bestehenden Smoke, er erweitert den Umfang der Zusage
  [`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung) nicht.
- Änderungen an `in-progress/slice-v2-homebrew-formula.md`, Welle, Roadmap oder
  Carveout-Inventar: das führt der Verantwortliche selbst nach.

## 2. Definition of Done

- [ ] Tap-Smoke: Build-from-source (bzw. gleichwertig) und `u-boot add` laufen im
  Smoke; `brew audit --strict` gatet ohne `|| true` (erst `--online` entscheiden);
  Akzeptanzkriterium 4 erfüllt, Beleg ist ein grüner Smoke-Lauf im Tap.
- [ ] F-5: Formel-`desc` Englisch und ASCII in `scripts/homebrew-formula.rb.tmpl`,
  Go-Test `cmd/uboot/homebrew_scripts_test.go` angepasst.
- [ ] F-8: Kommentar „sechs Binaries“ in `.github/workflows/publish.yml` und die
  Stabilitäts-Bestimmung (`is_stable` / `!contains(github.ref_name, '-')`)
  angeglichen — nur Kommentar bzw. Bedingung im Workflow.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [ ] Doku-Update, falls öffentlicher Vertrag berührt (Maintainer-Doku zum Tap).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

Alle Läufe Docker-only (`AGENTS.md` §3.1); kein Host-Python und kein Host-Paketmanager
als Nachweis — der Homebrew-Lauf findet im Smoke des Tap-Repos auf dem Runner statt.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `packaging/homebrew-tap/.github/workflows/smoke.yml` | update | Tap-Vorlage: Build-from-source, `u-boot add`, `brew audit --strict` gatend |
| Tap-Repo `pt9912/homebrew-u-boot` (`.github/workflows/smoke.yml`) | update | ausgerollte Fassung der Vorlage, per `scripts/tap-nachzug.sh` nachgezogen |
| `scripts/homebrew-formula.rb.tmpl` | update | F-5: `desc` Englisch, ASCII |
| `cmd/uboot/homebrew_scripts_test.go` | update | F-5: Test der `desc`; Akzeptanzkriterium 4 |
| `.github/workflows/publish.yml` | update | F-8: Kommentar „sechs Binaries“ und Stabilitäts-Bestimmung angleichen |

## 4. Trigger

**Start** (`next` → `in-progress`): nächster stabiler Release oder der Zeitpunkt,
an dem der Tap-Smoke ohnehin angefasst wird; kein Slice in `in-progress/`
(WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Tap-Smoke und Formel-/Workflow-Kommentare
  trennen sich in zwei Slices.
- `in-progress` → `open` (blockiert): `brew audit --strict` liefert echte
  Formel-Befunde, deren Behebung eigener Umfang ist.

## 5. Closure-Trigger

DoD vollständig, Review- und Verification-Evidence liegen vor, Closure-Notiz
mit Steering-Loop-Eintrag geschrieben, `make gates` auf dem Stand vor dem
`git mv` nach `done/` grün.

## 6. Risiken und offene Punkte

- `brew audit --strict` kann auf der Formel echte Befunde liefern —
  **Ausgang:** offen bis zur Closure; bei Eintritt „weiter offen“ mit Verweis auf
  den Befund.

## 7. Closure-Notiz

Offen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `cmd/uboot` (Test, keine
Sub-Area), `tools/`, `scripts/` (WZG) und `packaging/` samt
`.github/workflows/` als ausführbarer Vertrag
([`conventions.md`](../../../../harness/conventions.md) §Modus-Deklaration);
keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** u-boot führt bisher kein
Beobachtungs-Register; keine Treffer.

**Modus:** alle berührten Sub-Areas GF.
