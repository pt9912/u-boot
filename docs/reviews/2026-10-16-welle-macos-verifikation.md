# Review-Report: welle-macos-verifikation — 2026-10-16

**Review-Art:** Code | Plan — Diff gegen Plan, [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) und [`ADR-0018`](../plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md), Lastenheft 0.7.0 und Hard Rules (Maintainability; nicht gegen die DoD).

**Gegenstand:** Slice [`slice-v2-homebrew-formula`](../plan/planning/done/slice-v2-homebrew-formula.md) samt Artefakten (Skeleton, Fill-Skript, Nachzug-Skript, `publish.yml` Jobs `publish`/`tap`, `packaging/homebrew-tap/**`, `cmd/uboot/homebrew_scripts_test.go`, README EN/DE, [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md)) sowie Commits `109ce70` (Smoke-Runner) und `c873814` (Lastenheft 0.7.0). <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

**Skill:** `.harness/skills/reviewer.md` @ HEAD `c873814` · <!-- d-check:ignore (Skill-Pfad als Zitat) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-16

**Eingangs-Kontext:**

- [`slice-v2-homebrew-formula`](../plan/planning/done/slice-v2-homebrew-formula.md) und `welle-macos-verifikation` (Welle-Datei) <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
- [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) (Homebrew-Tap), [`ADR-0018`](../plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md) (Nested Docker), [`ADR-0019`](../plan/adr/0019-docker-daemon-aus-docker-apt-repo.md) (nur gelesen)
- [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache)
- `AGENTS.md` (Hard Rules), `harness/conventions.md`
- Belege geprüft: Tap-Smoke-Lauf `38900477533` (per `gh run view` und Log; beide Jobs grün, Formel 0.9.0, `--version`, `init`, `doctor`); `smoke.yml` im Tap ist byte-identisch mit der Vorlage `packaging/homebrew-tap/.github/workflows/smoke.yml`; Tap-Commit `a8498fd` ändert nur `smoke.yml`.
- Sensors: `make docs-check` (0 Befunde), `make doc-immutable RANGE="c873814~1..c873814"` (OK). Nicht gelaufen: `make gates`, `make test` (Diff ohne Go-Code-Änderung; der Go-Test `homebrew_scripts_test.go` wurde nur gelesen).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | README, Benutzerhandbuch und Plattformliste bewerben Homebrew für „macOS / Linux" (vier Plattformen). [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) §Entscheidung 6 nimmt Linux-Brew als beworbenen Weg aus, der Slice führt es als Out-of-Scope („bietet keinen Mehrwert"); der Smoke läuft nur auf macOS, die Linux-Zweige der Formel sind nicht belegt. | [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md), Maintainability | README.md:93 · „macOS / Linux with Homebrew (stable releases, four platforms)"; README.de.md:96; docs/user/benutzerhandbuch.md:109,1745 | nein — Review-only (`make docs-check` prüft das nicht) | Nutzerdoku bewirbt Weg, den ADR ausnimmt |
| F-2 | LOW | Benutzerhandbuch nennt in der `--privileged`-Warnung weiter „die VM von Docker Desktop oder Colima" als macOS-Fall von `nestedRuntime: docker`, während §Grenzen und [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (0.7.0) nur Linux und Colima zusagen. | [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) | docs/user/benutzerhandbuch.md:1005 | nein — Review-only | Restzusage Docker Desktop |
| F-3 | LOW | Der Slice-Kopf („repo-seitig geliefert, Tap-Einrichtung offen", 2026-10-01), §Offene Einrichtung (Schritte 2 und 3 offen) und §Stand nach v0.9.0 („Offen: Vorlage im Tap einspielen … Smoke grün sehen") widersprechen dem vorliegenden Zustand: Secret gesetzt, Vorlage im Tap (`a8498fd`), Smoke-Lauf `38900477533` grün. Der Lauf steht nicht im Slice. | AGENTS.md §3.7 (Zustandsfelder) | [`slice-v2-homebrew-formula.md`:3](../plan/planning/done/slice-v2-homebrew-formula.md) · „Tap-Einrichtung offen" | nein — Review-only | Slice-Stand veraltet |
| F-4 | LOW | Das Akzeptanzkriterium verlangt im Smoke `brew install --build-from-source` und die drei `LH-AK`-Pre-Checks (init/add/doctor). `smoke.yml` führt `brew install` (Binary), `init` und `doctor`, kein `add`; `brew audit --strict --online … \|\| true` kann nie fehlschlagen, T1 führt `brew audit` weiter als offen. | [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung), Maintainability | packaging/homebrew-tap/.github/workflows/smoke.yml · „brew audit --strict --online pt9912/u-boot/u-boot \|\| true" | nein — Review-only (Tap-Lauf) | Smoke schwächer als Akzeptanzkriterium |
| F-5 | LOW | `desc` der Formel ist deutsch und ohne Umlaute („Bootstrapper fuer reproduzierbare Docker-/Devcontainer-Stacks"); `brew info`/`brew search` zeigen sie Endnutzern, die übrige Nutzer-Ausgabe ist Englisch. | [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache) | scripts/homebrew-formula.rb.tmpl · `desc "Bootstrapper fuer …"` | nein — Review-only | Distributions-Artefakt mit deutschem Nutzertext |
| F-6 | LOW | Die Einengung der Zusage (Docker Desktop entfällt) ist nutzersichtbar, `CHANGELOG.md` `[Unreleased]` ist leer. | AGENTS.md §5 Regel 5 | CHANGELOG.md · `## [Unreleased]` (leer) | nein — Review-only | Öffentlicher Vertrag ohne CHANGELOG |
| F-7 | LOW | Die DoD des herausgelösten Desktop-Slice enthält keinen Schritt, der die Zusage im Lastenheft wiederherstellt; die Verschiebungsnotiz verlangt dafür einen CR. | [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) | [`slice-v1-sandbox-docker-nachholmessung-desktop.md`](../plan/planning/open/slice-v1-sandbox-docker-nachholmessung-desktop.md) · §2 DoD, §1 Ziel | nein — Review-only | Folge-Slice trägt Rückkante zur Spec nicht |
| F-8 | LOW | Der Kommentar in `publish.yml` und [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) sprechen von einer `SHA256SUMS` der „sechs Binaries"; die Datei enthält seit [`ADR-0017`](../plan/adr/0017-linux-pakete-nfpm.md) auch die Pakete. Der Job `tap` bestimmt „stabil" selbst über `!contains(ref_name,'-')`, `publish` über den Step-Output `is_stable`. | Maintainability | .github/workflows/publish.yml:212,290 | nein — Review-only | Doppelte Stabil-Bestimmung / Kommentar-Drift |
| F-9 | INFO | `doctor`-Hint „Install Docker Desktop or Docker Engine" betrifft den Host-Docker von u-boot selbst, nicht `nestedRuntime: docker`; kein Widerspruch zu 0.7.0. | [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) | internal/hexagon/application/doctor.go:321 | — | — |
| F-10 | INFO | `scripts/tap-nachzug.sh` und der Job `tap` haben keinen Test (nur das Fill-Skript ist per Go-Test gedeckt); der Nachzug ist durch den echten Lauf (Tap-Commit `3059299`, v0.7.0) belegt. | [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) | scripts/tap-nachzug.sh | Review-only | — |
| F-11 | INFO | Die manuelle Erprobung auf Apple Silicon (u-boot 0.9.0, Colima) ist nur als Prosa im Slice belegt; nicht prüfbar. Zuständig: Verifier. | verification.md | [`slice-v2-homebrew-formula`](../plan/planning/done/slice-v2-homebrew-formula.md).md §Stand nach v0.9.0 | — | — |

## Konsistenzprüfung 0.7.0 (Auftrag 2)

- `spec/lastenheft.md`: Version 0.7.0, §16-Zeile 0.7.0 (Datum 2026-10-16, Beschreibung, Quelle „Vereinbarung mit dem Projektinhaber") folgt dem Schema der Vorzeilen und ist chronologisch am Ende geführt; Zusage in [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) Zeile 452 und Historie deckungsgleich; kein Abwärts-Verweis im Lastenheft.
- [`ADR-0018`](../plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md): Der Diff fügt genau eine Zeile in die Geschichte ein, der Rumpf bleibt unverändert (Accepted bleibt immutable; `history-range-guard` OK). Die Rumpf-Aussagen zu Docker Desktop („Nachholmessung") sind Chronik und bleiben korrekt stehen; die Geschichte-Zeile verweist aufwärts auf das Lastenheft (zulässig).
- Welle-Datei, Desktop-Slice (`open/`, „ohne Welle") und `open/README.md` sind untereinander konsistent; kein Lifecycle-Doppelort.
- Restzusagen für Docker Desktop in README, README.de, CHANGELOG (nur `/.dockerenv`-Erkennung), `spec/spezifikation.md`, `spec/architecture.md`: keine. Rest: F-2.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `scripts/homebrew-formula-fill.sh` (fail-closed, Platzhalter, Exit-Codes) | geprüft, ohne Befund |
| `cmd/uboot/homebrew_scripts_test.go` (Fill, fail-closed, leerer Tag) | geprüft, ohne Befund |
| `scripts/tap-nachzug.sh` (Version-Check, Idempotenz, Token nur aus Env) | geprüft, ohne Befund (F-10 nur Hinweis) |
| `publish.yml` Job `tap` (Secret fehlt → Warnung, `persist-credentials: false`, Rechte `contents: read`, Tag/Token nur über Step-Env) | geprüft, ohne Befund (F-8 nur Wartbarkeit) |
| `109ce70` (`macos-13` → `macos-15-intel`, Lauf `38900477533` grün auf beiden) | geprüft, ohne Befund |
| Hexagon-/depguard-/Dual-Classifier-Linsen | nicht berührt, kein Go-Produktivcode im Diff |
| Sprachvertrag CLI-Ausgabe | nicht berührt (F-5 betrifft nur die Formel-`desc`) |
| Gate-Lockerung / Carveout-Buchführung | geprüft, ohne Befund (`brew audit \|\| true` im Tap-Smoke ist kein Repo-Gate, siehe F-4) |
| Referenzmodell Rang 1–4 (kein Abwärtsverweis) | geprüft, ohne Befund |
| Baseline-Bump | nicht berührt |
| `make gates`, `make test`, `make test-docker`, `make verify-depguard` | nicht ausgeführt: kein Go-/Docker-Pfad und keine Importregel im Diff; `make docs-check` und `make doc-immutable` liefen |
| Linux-Brew-Pfad der Formel, `brew audit` | nicht geprüft: kein Linux-/Audit-Lauf verfügbar (Docker-only, kein Host-brew) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 7 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Nutzerdoku bewirbt Weg, den ADR ausnimmt · Restzusage Docker Desktop · Slice-Stand veraltet · Smoke schwächer als Akzeptanzkriterium

## Verdikt

**Merge-blockierend:** ja für F-1 (MEDIUM; Quelle [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) und Nutzerdoku widersprechen sich), sonst nein. Kein HIGH, kein Rollen-Widerspruch.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7. Dieser Report ersetzt keine Verifikation (DoD, Belege F-11: Verifier).
