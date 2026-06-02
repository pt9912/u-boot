# Roadmap History v0.1-v0.3

> Archiviert am 2026-06-02; die lebendige Roadmap liegt in
> [`docs/plan/planning/in-progress/roadmap.md`](../plan/planning/in-progress/roadmap.md).
> Diese Datei entlastet die Roadmap von historischen Release-Details,
> ohne die Audit-Spur zu verlieren.

## Releases

| Version | Datum | Tag-Commit | Highlights | Detail-Slice |
| --- | --- | --- | --- | --- |
| v0.1.0 | 2026-05-31 | `49ec464` | MVP-Core M1..M8, Release-Pipeline, GHCR-Distribution | [`slice-v1-release-cut-v0.1.0`](../plan/planning/done/slice-v1-release-cut-v0.1.0.md) |
| v0.2.0 | 2026-06-01 | `595acdf` | Container-aware `doctor`, sechs Plattform-Binaries, Template-Katalog | [`slice-v1-release-cut-v0.2.0`](../plan/planning/done/slice-v1-release-cut-v0.2.0.md) |
| v0.3.0 | 2026-06-01 | `54bc384` | Add-on Catalogue Expansion: `remove`, `--with-deps`, Keycloak, OpenTelemetry | [`slice-v1-release-cut-v0.3.0`](../plan/planning/done/slice-v1-release-cut-v0.3.0.md) |

## v0.1.0 MVP-Cluster

| Cluster | Geliefert | Detailquellen |
| --- | --- | --- |
| Spec und ADR-Basis | Lastenheft, Go-Entscheidung, hexagonale Architektur, CI- und Quality-Entscheidungen | [`spec/lastenheft.md`](../../spec/lastenheft.md), [`ADR-0001`](../plan/adr/0001-implementierungssprache-go.md), [`ADR-0002`](../plan/adr/0002-hexagonale-architektur.md), [`ADR-0003`](../plan/adr/0003-solid-nahes-lint-profil.md), [`ADR-0004`](../plan/adr/0004-ci-system.md) |
| Build und Gates | Docker-only Makefile, lint/test/coverage, docs-check, security-gates, image-scan | [`slice-m1-repo-skeleton`](../plan/planning/done/slice-m1-repo-skeleton.md), [`slice-m2b-solid-lint-profil`](../plan/planning/done/slice-m2b-solid-lint-profil.md), [`slice-m2c-ci-pipeline`](../plan/planning/done/slice-m2c-ci-pipeline.md), [`slice-v1-markdown-link-validator`](../plan/planning/done/slice-v1-markdown-link-validator.md), [`slice-v1-release-pipeline`](../plan/planning/done/slice-v1-release-pipeline.md) |
| Core CLI | `init`, `doctor`, `add postgres`, `up`, `down`, `generate`, `config` | [`slice-m3-init-flow`](../plan/planning/done/slice-m3-init-flow.md), [`slice-m4-doctor`](../plan/planning/done/slice-m4-doctor.md), [`slice-m5-add-postgres`](../plan/planning/done/slice-m5-add-postgres.md), [`slice-m6-up-down`](../plan/planning/done/slice-m6-up-down.md), [`slice-m7-generate`](../plan/planning/done/slice-m7-generate.md), [`slice-m8-config`](../plan/planning/done/slice-m8-config.md) |
| Closure und Carveouts | MVP-Acceptance, Carveout-Disziplin, depguard-/gomodguard-/coverage-Auflösung | [`slice-mvp-closure`](../plan/planning/done/slice-mvp-closure.md), [`slice-m2d-carveout-disziplin`](../plan/planning/done/slice-m2d-carveout-disziplin.md), [`carveouts.md`](../plan/planning/in-progress/carveouts.md) |

## v0.2.0 Cluster

| Cluster | Geliefert | Detailquellen |
| --- | --- | --- |
| Container-aware Runtime | `doctor` erkennt Container-Kontext und vermeidet Host-False-Positives | [`slice-v0.1.1-doctor-container-awareness`](../plan/planning/done/slice-v0.1.1-doctor-container-awareness.md) |
| Binary Distribution | sechs Plattform-Binaries Linux/macOS/Windows × amd64/arm64 als Release-Assets | [`slice-v2-binary-distribution`](../plan/planning/done/slice-v2-binary-distribution.md), [`ADR-0007`](../plan/adr/0007-distributionswege-ghcr.md) |
| Template-Katalog | `template list`, `init --template basic`, statisches Plugin-/Template-Modell | [`slice-v1-template-list`](../plan/planning/done/slice-v1-template-list.md), [`slice-v1-template-init`](../plan/planning/done/slice-v1-template-init.md), [`ADR-0008`](../plan/adr/0008-plugin-system-statisch.md), [`ADR-0009`](../plan/adr/0009-template-format-yaml-files.md) |

## v0.3.0 Cluster

| Cluster | Geliefert | Detailquellen |
| --- | --- | --- |
| Add-on Lifecycle | `remove <service>`, dependency resolution, `--with-deps` | [`slice-v1-add-remove`](../plan/planning/done/slice-v1-add-remove.md), [`slice-v1-addons-deps`](../plan/planning/done/slice-v1-addons-deps.md) |
| Add-on Catalogue | Keycloak und OpenTelemetry inklusive Tests, Templates und Doctor-/E2E-Pfade | [`slice-v1-keycloak`](../plan/planning/done/slice-v1-keycloak.md), [`slice-v1-otel`](../plan/planning/done/slice-v1-otel.md) |
| V1 Audit | Makefile-Aggregatoren, dokumentierte Schnittstellen, Containerfreundlichkeit | [`slice-v1-audit-done`](../plan/planning/done/slice-v1-audit-done.md) |
| Devcontainer Features | Feature-Katalog, allowlist, drift-doctor, README-Entlastung | [`slice-v1-devcontainer-features`](../plan/planning/done/slice-v1-devcontainer-features.md), [`slice-followup-devcontainer-features-drift-doctor`](../plan/planning/done/slice-followup-devcontainer-features-drift-doctor.md) |

## Offene Restwege aus den Releases

| Restweg | Status | Aktuelle Quelle |
| --- | --- | --- |
| Homebrew | Plan-Stub `open/`, Trigger: erste macOS-Nutzeranfrage | [`slice-v2-homebrew-formula`](../plan/planning/open/slice-v2-homebrew-formula.md) |
| Debian/RPM | Plan-Stub `open/`, Trigger: konkrete Distro-Anfrage | [`slice-v2-distro-pakete`](../plan/planning/open/slice-v2-distro-pakete.md) |
| Keycloak-CI-Flake | Plan-Stub `open/`, Trigger: belastbarer CI-/Quay-Befund | [`slice-v1-keycloak-ci-flake`](../plan/planning/open/slice-v1-keycloak-ci-flake.md) |
| Branch Protection | Nutzeraktion im GitHub-UI, nicht versionierbar | [`docs/user/branch-protection.md`](../user/branch-protection.md) |

## Hinweis zur Detailtiefe

Die alten Roadmap-Tabellen enthielten viele Commit-IDs und
Tranche-Details. Die kanonische Detailquelle dafür sind jetzt die
jeweiligen `done/`-Slices und die Release-Cut-Slices. Diese Archivseite
nennt bewusst nur Cluster und Anker, damit historische Information
auffindbar bleibt, ohne die lebendige Roadmap zu überladen.
