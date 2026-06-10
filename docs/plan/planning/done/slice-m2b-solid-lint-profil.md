# Slice M2b: SOLID-nahes Lint-Profil

> **Status:** Done
> **DoD:** Commit `365e532`
> **Retro-Plan:** Retroaktiv geschrieben 2026-05-27 (siehe [`slice-m3-retroaktive-slice-plaene`](slice-m3-retroaktive-slice-plaene.md))

## Auslöser

[`LH-QA-004`](../../../../spec/lastenheft.md#lh-qa-004-linting-solid-nahes-lint-profil) war bis M2 generisch („Linting (V1)") und das u-boot-Repo
hatte nur die 5 Default-Linter aktiv. Die Schwester-Projekte `m-trace`
und `k-deskflight` hatten bereits ein etabliertes SOLID-nahes Profil
(5 Default + 24 SOLID-nahe Linter + `depguard`), das sich in Reviews
bewährt hatte. u-boot sollte die gleiche Profil-Stärke bekommen,
**bevor** der erste produktive Code in `internal/` landet —
nachgelagertes Härten erzeugt sonst große Refactoring-Schübe.

## Lieferumfang

- **Spec-Verschärfung**: [`LH-QA-004`](../../../../spec/lastenheft.md#lh-qa-004-linting-solid-nahes-lint-profil) von „V1" auf MVP-Pflicht gehoben,
  Verweis auf `docs/user/quality.md` §1.2 / §1.3 als SSOT für Linter-
  Liste + Carveouts; `//nolint`-Pragmas verboten; Detail in [ADR-0003](../../adr/0003-solid-nahes-lint-profil.md).
- **Neue Doku** `docs/user/quality.md`: Statische Analyse §1, SOLID-
  Linter-Tabelle §1.2 (24 Linter), Carveout-Tabelle §1.3, Tests §2,
  Coverage §3, Security §4, Architektur-Enforcement §5. TypeScript-Slot
  §1.1 reserviert (u-boot bleibt Go-only).
- **[ADR-0003](../../adr/0003-solid-nahes-lint-profil.md)** (`docs/plan/adr/0003-solid-nahes-lint-profil.md`):
  Vorlage `m-trace`/`k-deskflight`, Schwellen 1:1 übernommen, Trade-offs
  (Lint-Stage-Geschwindigkeit + Einarbeitung) gegen verworfene
  Alternativen (nur Defaults, engere Schwellen, volles revive-Custom).
- **`.golangci.yml`** von 5+depguard auf 5+24+depguard erweitert:
  - Schwellen analog `m-trace` (`cyclop=15`, `funlen=100/60`,
    `gocognit=20`, `gocyclo=15`, `interfacebloat=10`, `maintidx=20`,
    `nestif=5`, `dupl=150`).
  - `gomodguard_v2` statt v1 (v1 in golangci-lint v2.12.0 deprecated).
  - `forbidigo` verbietet `fmt.Print*` (Logging gehört in `log/slog`).
  - `ireturn`-Allowlist enthält die Hex-Ports.
  - Permanente Carveouts mit `Why:`-Kommentar (test-files, cmd/uboot
    wiring layer).
- READMEs (de/en) verweisen auf `docs/user/quality.md`.

## Akzeptanz

- `make gates` grün mit dem erweiterten Lint-Profil (im Bootstrap-
  Modus, da `./internal/...` leer).
- 24 SOLID-nahe Linter aktiv und in `quality.md` §1.2 dokumentiert.
- [`LH-QA-004`](../../../../spec/lastenheft.md#lh-qa-004-linting-solid-nahes-lint-profil) MVP-Pflicht erfüllt; [`LH-MVP-001`](../../../../spec/lastenheft.md#lh-mvp-001-muss-im-mvp-enthalten-sein) ergänzt.
- [ADR-0003](../../adr/0003-solid-nahes-lint-profil.md) verlinkt aus `quality.md`.

## Bezug

- Auslösende Spec: [`LH-QA-004`](../../../../spec/lastenheft.md#lh-qa-004-linting-solid-nahes-lint-profil).
- ADR: `0003-solid-nahes-lint-profil.md`.
- Vorgänger: [`slice-m2-hexagonale-architektur`](slice-m2-hexagonale-architektur.md).
- Nachfolger: M2c (CI bekommt die Gates).
- Nachfolge-Slices: [`slice-m3-depguard-aktivierung-verifizieren`](slice-m3-depguard-aktivierung-verifizieren.md)
  (verifiziert die 8 depguard-Blöcke real), [`slice-m3-gomodguard-rules`](slice-m3-gomodguard-rules.md)
  (füllt den ursprünglich leeren gomodguard-Block).
