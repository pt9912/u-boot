# Verification Evidence: slice-v2-homebrew-formula — 2026-10-16

**Rolle:** Verifier (frischer Kontext, kein Self-Review). Frage: Wurde der Slice richtig gebaut, gemessen an Plan, DoD und Sensors? Form nach `harness/verification.md`. <!-- d-check:ignore (Pfad als Zitat) -->

**Gegenstand:** Slice [`slice-v2-homebrew-formula`](../plan/planning/done/slice-v2-homebrew-formula.md) (Delivery-Hash `0cab536`, Folgecommits bis `a4fd3a4`), Welle `welle-macos-verifikation` §3 (Closure-Trigger). <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) --> Der Slice ist in älterer Form geschrieben (Akzeptanzkriterien statt §2-DoD, Tranchen T1..T4); Maßstab sind seine vier Akzeptanzkriterien, [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md), [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung) und die Closure-Trigger der Welle.

## Verification Evidence

Scope:
- Slice: [`slice-v2-homebrew-formula`](../plan/planning/done/slice-v2-homebrew-formula.md)
- IDs: [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung), [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md), [`ADR-0007`](../plan/adr/0007-distributionswege-ghcr.md) (nur als Bezug)
- Artefakte: `scripts/homebrew-formula.rb.tmpl`, `scripts/homebrew-formula-fill.sh`, `scripts/tap-nachzug.sh`, `.github/workflows/publish.yml` (Jobs `publish`/`tap`), `packaging/homebrew-tap/**`, `cmd/uboot/homebrew_scripts_test.go`, README EN/DE, Tap `pt9912/homebrew-u-boot`

DoD-Abgleich (Akzeptanzkriterien des Slice):
- [x] Formel im Tap, SHA256-Pin pro Plattform — Evidence: Tap `Formula/u-boot.rb` trägt `version "0.9.0"` (per `gh api` gelesen); Fill-Skript fail-closed, Go-Test `homebrew_scripts_test.go` (läuft in `make test`, grün); Tap-Commit `3059299` (v0.7.0) belegt den Nachzug im echten Lauf.
- [x] `publish.yml` aktualisiert die Formel pro Tag-Push — Evidence: Job `tap` + `scripts/tap-nachzug.sh`; Secret besteht `tap-check`; Nachzüge v0.7.0..v0.9.0 (Slice §Stand, Welle §2). Nachzug-Skript und Job selbst ohne Test (Review F-10, INFO; durch echte Läufe belegt).
- [x] README-Install-Block EN + DE — Evidence: `README.md:93` und `README.de.md:96` („macOS with Homebrew", nur macOS; F-1 behoben), `brew trust`-Hinweis vorhanden. Die Reihenfolge „zweite-empfohlene Variante nach Binary" wurde im Block nicht einzeln nachgeprüft; kein Review-Befund dazu.
- [~] Tap-Smoke mit `brew install --build-from-source` und den drei Pre-Checks (init/add/doctor) — teilweise: Der Smoke (`packaging/homebrew-tap/.github/workflows/smoke.yml`, identisch im Tap, Commit `a8498fd` ändert nur `smoke.yml`) führt tap, `brew install` (Binary-Pfad), `--version`, `init`, `doctor` (Exit 0 oder 11 zulässig) und `brew audit --strict --online … || true`. Es fehlen `--build-from-source` und `add`; der Audit ist nicht gatend. Bewertung siehe Offene Punkte.

Closure-Trigger der Welle (§3):
- [ ] Slice liegt in `done/` — offen: Slice liegt in `in-progress/`; Closure-Schritte (Verschiebung, §7-Notiz der Welle „Offen.") stehen aus. Folgt auf dieses Verdikt.
- [x] Smoke zur aktuellen Formel auf Apple Silicon und Intel grün — Evidence: `gh run view 38900477533 --repo pt9912/homebrew-u-boot`: `install (macos-14)` und `install (macos-15-intel)` grün, Formel 0.9.0; Log zeigt `Formula u-boot (0.9.0)` und `u-boot version 0.9.0` (macos-14). Die Annotationen im Lauf (AWS-Tap-Hinweis, `vagrant.rb`-Importfehler) stammen von vorinstallierten Runner-Taps, nicht von der Formel.
- [~] Erprobung auf echtem Mac — Evidence: nur Prosa des Projektinhabers (Slice §Stand nach `v0.9.0`: tap/trust/install, `--version` = 0.9.0, `init`, `doctor` Exit 0, 0 Fehler, 1 Warnung). Vom Verifier nicht reproduzierbar (Docker-only, kein Host-brew); Beleg ist die Zusage des Inhabers, der Smoke auf zwei echten macOS-Runnern stützt sie unabhängig (Apple Silicon `--version`/`init`/`doctor`). Akzeptiert, als Zusage gekennzeichnet.
- [x] `make gates` grün — siehe Sensors.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | selbst gefahren, Exit 0; `coverage-gate: OK — coverage 90.90% meets threshold 90%` (lint, test, docs-check, baseline-verify Teil von `gates`) |
| `make docs-check` | pass | Teil von `make gates`; Review-Lauf meldete 0 Befunde, nach `a4fd3a4` erneut grün |
| `make test-docker` | not run | kein Docker-Laufzeit-Pfad berührt (Diff: Skripte, Workflow, Doku, Tap-Vorlagen) |
| `make verify-depguard` | not run | keine Importregel berührt |
| Tap-Smoke (extern) | pass | Lauf `38900477533`, beide Jobs grün (per `gh run view` und Log selbst gelesen) |

Traceability:
| ID | Beleg |
| --- | --- |
| [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung) | Tap-Formel 0.9.0, Smoke-Lauf `38900477533`, `README.md:93`, `homebrew_scripts_test.go`; Carveout-Zeile noch offen (siehe Carveouts) |
| [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) | Entscheidungen 2 (Formel als Release-Asset, `u-boot.rb`), Tap-Nachzug per Job `tap`, Tap `pt9912/homebrew-u-boot` (belegt); 6 (Linux-Brew nicht beworben) — konform seit `a4fd3a4` (F-1 behoben, README/Handbuch nur macOS) |

Plan-vs-Code-Diff:
- Plan: vier Plattformen, Test-Block `--version`, fail-closed Fill — Code: vorhanden, vom Go-Test gedeckt.
- Plan T1 „`brew audit --strict`" — Code/Beleg: nie gatend gelaufen (Smoke: `|| true`); die Slice-Notiz „Offen: `brew audit --strict` auf einem Mac" wurde nicht ausdrücklich geschlossen.
- Plan T4 (CHANGELOG, Carveout reduziert, [`ADR-0007`](../plan/adr/0007-distributionswege-ghcr.md)-Anpassung) — CHANGELOG `[Unreleased]` gefüllt (F-6 behoben); [`ADR-0007`](../plan/adr/0007-distributionswege-ghcr.md) bewusst nicht umgeschrieben, ersetzt durch [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) (Accepted-immutable, Slice-Lieferstand dokumentiert); Carveout-Zeile **noch nicht** aufgelöst.
- Slice-Kopf/Stand (F-3) auf aktuellen Zustand gebracht („geliefert und belegt", Smoke-Lauf `38900477533` genannt).

Carveouts:
- Neu: none
- Gelöst: none (noch)
- Unverändert: [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung)-Zeile in `in-progress/carveouts.md` (temporär, Schließbedingung „grüner macOS-Smoke zur Formel 0.7.0 und Prüfung auf einem Mac" ist erfüllt, Zeile steht noch). Sie muss in der Closure aufgelöst werden; ihr Wortlaut nennt noch 0.7.0 statt der aktuellen Formel.

Nicht ausgeführt:
- `make test-docker`, `make verify-depguard` — kein berührter Pfad (siehe oben).
- Eigener Mac-/`brew`-Lauf, `brew audit` — Docker-only-Umgebung, kein Host-brew (AGENTS.md §3.1).

Commit / Artefakt:
- Delivery `0cab536`; Review-Fixes `a4fd3a4`; Tap-Commit `a8498fd` (Smoke-Vorlage), Formel 0.9.0; Smoke-Lauf `38900477533`; Review-Report [`2026-10-16-welle-macos-verifikation.md`](2026-10-16-welle-macos-verifikation.md) (0 HIGH, 1 MEDIUM, 7 LOW, 3 INFO).

## Review-Findings: Stand und Bewertung

| Finding | Stand | Bewertung für die Closure |
|---|---|---|
| F-1 (MEDIUM), F-2, F-3, F-6, F-7 | in `a4fd3a4` behoben; F-1 und F-3 am Ist-Stand gegengeprüft (README nur macOS, Slice-Kopf aktuell, CHANGELOG gefüllt) | erledigt |
| F-4 Smoke schwächer als Kriterium | offen | **nicht blockierend, aber zu entscheiden**: ein Binär-Install ist der Pfad, den Nutzer gehen (`--build-from-source` würde hier nur Go-Bau-Dauer testen, die Formel hat kein Source-Build). Das Akzeptanzkriterium bleibt damit wörtlich unerfüllt. Auflage: Slice §Lieferstand benennt die Abweichung (kein `--build-from-source`, kein `add`, Audit nicht gatend) ausdrücklich als Setzung, oder ein Folge-Slice nimmt `add` und gatendes `brew audit` auf. Der Smoke deckt die Welle-Trigger (install, `--version`, `init`, `doctor`) ab. |
| F-5 deutsche Formel-`desc` | offen | nicht blockierend: Distributions-Artefakt, kein CLI-Output (Sprachvertrag [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache) trifft CLI-Ausgabe und generierte Dateien); Folge-Arbeit (Tmpl-Zeile + Tap-Nachzug), mit nächstem Release erledigbar. |
| F-8 „sechs Binaries"-Kommentar | offen | nicht blockierend: Kommentar-/ADR-Text-Drift ohne Verhaltenswirkung ([`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) ist Accepted, nicht umzuschreiben; Kommentar in `publish.yml` bei nächster Berührung). |

## Verdikt

**Bestanden, mit Auflagen für die Closure** (kein Befund, der das Schließen sperrt):

1. Akzeptanzkriterium 4 ist nur teilweise erfüllt (F-4); die Abweichung wird im Slice (Closure-Notiz) ausdrücklich als Setzung oder Folge-Slice festgehalten, sonst bleibt „✅" ein Wortlaut-Widerspruch.
2. Carveout-Zeile [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung) in `carveouts.md` auflösen und ins „Gelöst"-Inventar führen (Schließbedingung erfüllt); Slice-Closure-Notiz §7 und Welle-§7 schreiben, dann reiner `git mv` nach `done/` (AGENTS.md §3.3, Inhalt zuerst).
3. Mac-Erprobung bleibt Zusage des Projektinhabers (nicht reproduzierbar durch den Verifier); in der Closure-Notiz als solche ausweisen.
4. `brew audit --strict` (T1) ist nie gatend gelaufen; in der Closure-Notiz als nicht erbracht oder als Folge-Arbeit führen.

Folge-Arbeit ohne Closure-Bezug: F-5 (englische `desc`), F-8 (Kommentar), F-4 (gatender Audit, `add`-Schritt).
