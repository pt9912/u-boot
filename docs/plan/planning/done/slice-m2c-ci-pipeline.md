# Slice M2c: CI-Pipeline

> **Status:** Done
> **DoD:** Commit `9a74e35`
> **Retro-Plan:** Retroaktiv geschrieben 2026-05-27 (siehe [`slice-m3-retroaktive-slice-plaene`](slice-m3-retroaktive-slice-plaene.md))
>
> **Folge-Stand (post-Lieferung):** Der dritte PR-blockierende Job
> `image-scan (trivy HIGH+CRITICAL)` wurde mit
> [`slice-v1-release-pipeline`](slice-v1-release-pipeline.md) T3
> (`8212889`) ergänzt; [`LH-QA-003`](../../../../spec/lastenheft.md#lh-qa-003-ci-fähigkeit-github-actions) ist seither auf drei Jobs gehoben.
> Dieser Slice-Plan beschreibt absichtlich den ursprünglichen
> M2c-Lieferstand (zwei Jobs) und wird nicht rückwirkend
> aktualisiert.

## Auslöser

[`LH-QA-003`](../../../../spec/lastenheft.md#lh-qa-003-ci-fähigkeit-github-actions) war bis M2b generisch („soll in CI testbar sein") ohne
konkreten Workflow. Vor dem ersten produktiven Slice (M3) musste die
PR-Pipeline stehen, damit jeder Feature-Slice von Anfang an gegen
`make gates` + `make govulncheck` läuft und der „dev-fertig vs.
CI-fertig"-Drift erst gar nicht entsteht.

## Lieferumfang

- **Spec-Verschärfung** [`LH-QA-003`](../../../../spec/lastenheft.md#lh-qa-003-ci-fähigkeit-github-actions): GitHub Actions konkret,
  `.github/workflows/ci.yml` mit zwei Jobs `gates` + `security-gates`,
  beide PR-blockierend, Docker-only-Pfad, SHA-pinned Actions,
  Top-Level `permissions: {}`, Timeout 20 min. [`LH-MVP-001`](../../../../spec/lastenheft.md#lh-mvp-001-muss-im-mvp-enthalten-sein) ergänzt.
- **[ADR-0004](../../adr/0004-ci-system.md)** (`docs/plan/adr/0004-ci-system.md`): `k-deskflight` als
  Vorlage statt `m-trace`'s 9 Workflows (zu reich); Trade-off
  GitHub-Actions-Vendor-Bindung kompensiert durch Make-Targets als SSOT;
  Folgepunkte explizit aufgelistet (Image-Publish, Trivy,
  Branch-Protection — alle drei sind in
  [`slice-v1-release-pipeline`](../done/slice-v1-release-pipeline.md)
  gebündelt).
- **Workflow** `.github/workflows/ci.yml`:
  - Zwei parallele Jobs: `gates` (`make gates`) und `security-gates`
    (`make govulncheck`).
  - Trigger: `pull_request` + `push` auf `main`.
  - Runner `ubuntu-latest` mit `DOCKER_BUILDKIT=1`.
  - `actions/checkout` SHA-gepinnt (`de0fac2 # v6.0.2`).
- **Doku** `docs/user/quality.md` §6 CI-Pipeline ergänzt mit Pflichten
  und Hinweis auf Branch-Protection im UI.
- **Roadmap**: M2c als Done markiert.

## Akzeptanz

- `actionlint` (rhysd/actionlint:latest) grün gegen den neuen Workflow.
- `make gates` weiterhin grün.
- [`LH-QA-003`](../../../../spec/lastenheft.md#lh-qa-003-ci-fähigkeit-github-actions) (konkrete Form) und [`LH-MVP-001`](../../../../spec/lastenheft.md#lh-mvp-001-muss-im-mvp-enthalten-sein)-Ergänzung abgehakt.
- [ADR-0004](../../adr/0004-ci-system.md) mit Folgepunkten dokumentiert.

## Bezug

- Auslösende Spec: [`LH-QA-003`](../../../../spec/lastenheft.md#lh-qa-003-ci-fähigkeit-github-actions).
- ADR: `0004-ci-system.md`.
- Vorgänger: [`slice-m2b-solid-lint-profil`](slice-m2b-solid-lint-profil.md).
- Nachfolger: M2d (Carveout-Disziplin sammelt alle bewussten Lücken
  inkl. [ADR-0004](../../adr/0004-ci-system.md)-Folgepunkte ein).
- Aufhebt mit: [`slice-v1-release-pipeline`](../done/slice-v1-release-pipeline.md)
  (Image-Publish + Trivy + Branch-Protection — 3 Folgepunkte gebündelt).
