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
[`slice-v2-homebrew-formula`](../done/slice-v2-homebrew-formula.md).

**Berührte Spec-Stellen:** —

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Homebrew-Smoke im Tap erfüllt das Akzeptanzkriterium 4 von
[`slice-v2-homebrew-formula`](../done/slice-v2-homebrew-formula.md) voll:
`brew install --build-from-source` (bzw. gleichwertig), `u-boot add` im Smoke und
ein gatendes `brew audit --strict` ohne `|| true`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Eine Änderung von [`ADR-0016`](../../adr/0016-homebrew-distribution-per-tap.md):
  die ADR ist Accepted und immutable; F-8 betrifft nur den Kommentar im Workflow,
  nicht die Entscheidung.
- Neue Plattformen oder Distributionswege im Tap: der Slice härtet den
  bestehenden Smoke, er erweitert den Umfang der Zusage
  [`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung) nicht.
- Änderungen an `done/slice-v2-homebrew-formula.md`, Welle, Roadmap oder
  Carveout-Inventar: das führt der Verantwortliche selbst nach.

## 2. Definition of Done

- [x] Tap-Smoke: Build-from-source (bzw. gleichwertig) und `u-boot add` laufen im
  Smoke; `brew audit --strict` gatet ohne `|| true` (erst `--online` entscheiden);
  Akzeptanzkriterium 4 erfüllt, Beleg ist ein grüner Smoke-Lauf im Tap (Lauf 2, §7).
- [x] F-5: Formel-`desc` Englisch und ASCII in `scripts/homebrew-formula.rb.tmpl`,
  Go-Test `cmd/uboot/homebrew_scripts_test.go` angepasst.
- [x] F-8: Kommentar „sechs Binaries“ in `.github/workflows/publish.yml` und die
  Stabilitäts-Bestimmung (`is_stable` / `!contains(github.ref_name, '-')`)
  angeglichen — nur Kommentar bzw. Bedingung im Workflow.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [x] Doku-Update, falls öffentlicher Vertrag berührt (Maintainer-Doku zum Tap).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne Wellen-Betrieb nach dem `git mv` geprüft.

## 3. Plan (vor Code)

Alle Läufe Docker-only (`AGENTS.md` §3.1); kein Host-Python und kein Host-Paketmanager
als Nachweis — der Homebrew-Lauf findet im Smoke des Tap-Repos auf dem Runner statt.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `packaging/homebrew-tap/.github/workflows/smoke.yml` | update | Tap-Vorlage: Build-from-source, `brew test`, `u-boot add`, `brew audit --strict` gatend |
| Tap-Repo `pt9912/homebrew-u-boot` (`.github/workflows/smoke.yml`) | update | ausgerollte Fassung der Vorlage (Tap-Commit `76a7faa`) |
| Tap-Repo `pt9912/homebrew-u-boot` (`Formula/u-boot.rb`) | update | die ausgelieferte Formel 0.9.0 trägt keine `version`-Zeile mehr (Tap-Commit `5129ebc`); Homebrew liest die Version aus der URL |
| `scripts/homebrew-formula.rb.tmpl` | update | F-5: `desc` Englisch, ASCII; keine `version`-Zeile |
| `scripts/homebrew-formula-fill.sh` | update | der Platzhalter für `version` entfällt |
| `cmd/uboot/homebrew_scripts_test.go` | update | F-5: Test der `desc`; Test prüft das Fehlen der `version`-Zeile |
| `.github/workflows/publish.yml` | update | F-8: Kommentar „sechs Binaries“ und Stabilitäts-Bestimmung angleichen |
| `packaging/nfpm.yaml` | update | Review F-2: Paketbeschreibung Englisch ([`LH-LESE-002`](../../../../spec/lastenheft.md#lh-lese-002--sprache)) |
| `docs/maintainer/releasing.md` | update | Doku: Maintainer-Beschreibung des gehärteten Tap-Smoke |

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
  **Ausgang:** eingetreten und behoben: ein Befund (redundante `version`-Zeile),
  Lauf 1 rot, behoben in `8d81b47` und Tap-Commit `5129ebc`, Lauf 2 grün.

## 7. Closure-Notiz

**Belege:** Review-Report
[`2026-10-17-slice-v2-homebrew-smoke-haertung`](../../../reviews/2026-10-17-slice-v2-homebrew-smoke-haertung.md) (0 HIGH, 2 MEDIUM, 1 LOW, 3 INFO;
F-2 und F-3 behoben in `8877fa3`; F-1, F-5 und F-6 als Grenze bzw. Setzung unten; F-4 durch die Läufe unten eingelöst). Verification Evidence
[`2026-10-17-verification-slice-v2-homebrew-smoke-haertung`](../../../reviews/2026-10-17-verification-slice-v2-homebrew-smoke-haertung.md): bestanden mit Auflagen, hier erfüllt.
`make gates` grün; `make doc-trace`: 83 Anforderungen, 0 Waisen. Tap-Läufe nach der Freigabe des Projektinhabers vom 2026-10-17 (Tap-Commit `76a7faa`, neue `smoke.yml`):

- **Lauf 1** (`39008767100`, `workflow_dispatch`): rot auf `macos-14` und `macos-15-intel`, nur im Schritt „Formel pruefen (gatend)“: `brew audit --strict --online` meldete
  „Stable: `version 0.9.0` is redundant with version scanned from URL“. Installation, `brew test u-boot`, `init`, `add postgres` (Managed-Block-Marker) und `doctor` waren grün.
- **Lauf 2** (`39009128343`, `push` auf `Formula/u-boot.rb`): grün auf `macos-14` und `macos-15-intel`, einschließlich des gatenden Audits. Das Gate wurde nie auf `|| true` zurückgenommen.

Akzeptanzkriterium 4 aus [`slice-v2-homebrew-formula`](../done/slice-v2-homebrew-formula.md) ist damit belegt
([`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung), [`ADR-0016`](../../adr/0016-homebrew-distribution-per-tap.md)).

- **Was hat funktioniert:** Das gatende Audit zeigte beim ersten echten Lauf genau einen Befund. Die Prüfschritte `brew test u-boot`, `add postgres` mit Managed-Block-Marker
  und die `doctor`-Toleranz liefen grün. Die Warnung des Reviewers zur `version`-Zeile (Review F-4) trat ein.
- **Was ging anders als geplant:** Der Slice sollte nur den Smoke härten; die Formel brauchte zusätzlich einen Fix, weil das Audit zum ersten Mal gatend lief
  (Repo: `8d81b47`, `version`-Zeile aus Vorlage und `scripts/homebrew-formula-fill.sh` entfernt, Test prüft das Fehlen). Die Formel im Tap wurde mit Freigabe des Projektinhabers
  einmalig von Hand korrigiert (Tap-Commit `5129ebc`: nur diese Zeile); der Release-Pfad zieht die korrigierte Vorlage beim nächsten stabilen Tag nach. §3 nannte `packaging/nfpm.yaml`
  und `docs/maintainer/releasing.md` nicht und ist berichtigt.
- **Setzungen:**
  - F-5: `--build-from-source` bei einer Binary-Formel ist dieselbe Installation wie ohne Flag; der inhaltliche Gegenwert ist `brew test u-boot`. Akzeptanzkriterium 4 („bzw. gleichwertig“)
    gilt damit als erfüllt — als Setzung festgehalten.
  - F-1: Für `publish.yml` gibt es keinen Sensor. Die Verdrahtung des Job-Outputs `is_stable` ist gelesen und korrekt; ein Umbenennen von Output oder Step-Id würde `tap` bei stabilen Tags
    lautlos überspringen. Erst der nächste stabile Tag belegt es. Grenze, benannt.
  - F-6: Die Denylist-Prüfung der `desc` kennt deutsche Wörter nur aus einer Liste; für eine einzeilige `desc` akzeptiert.
- **Steering-Loop-Eintrag:** Ein Gate, das zum ersten Mal gatend läuft, liefert im Zielsystem Befunde, die kein Lesen vorhersagt (hier: ein Befund, den Review und Verifier aus der Formel allein
  nicht entscheiden konnten). Es ist vor dem Slice-Abschluss im echten Lauf zu prüfen, nicht auf den nächsten Release zu verschieben. Gezählt, nicht verkörpert: erstes Auftreten dieser Klasse.
  Die Klasse „Plan-Aussage über Code stimmt nicht“ trat in drei Vorgängen der letzten Welle auf; siehe den nächsten Punkt für den vierten.
- **Beobachtungs-Register (`../observations/`):** keine neue Beobachtung angelegt. Die Plan-Aussage in §3 (zwei Dateien fehlten: `packaging/nfpm.yaml`, `docs/maintainer/releasing.md`) ist ein
  weiterer Beleg für „Plan-Aussage über Code stimmt nicht“ — ein vierter Vorgang. Ob dafür eine Beobachtung angelegt wird, entscheidet der Projektinhaber; hier wird nichts angelegt.
- **Folge-Slices:** keine.
- **Release-Hinweis:** Die korrigierte Formel-Vorlage und der Job-Output `is_stable` greifen erst beim nächsten stabilen Tag; bis dahin ist `publish.yml` ungetestet — Grenze.
- **Risiken aus §6:** alle tragen ihren Ausgang (eingetreten und behoben: der Audit-Befund).
- **Drei Paarungen:** Anker (kein Eintrag mit Zielort, nichts verkörpert) · Folge-Slice (keine) · Register (keine neue Beobachtung, notiert) — nach dem `git mv` zu prüfen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `cmd/uboot` (Test, keine
Sub-Area), `tools/`, `scripts/` (WZG) und `packaging/` samt
`.github/workflows/` als ausführbarer Vertrag
([`conventions.md`](../../../../harness/conventions.md) §Modus-Deklaration);
keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** u-boot führt bisher kein
Beobachtungs-Register; keine Treffer.

**Modus:** alle berührten Sub-Areas GF.
