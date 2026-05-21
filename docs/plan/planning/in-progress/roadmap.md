# u-boot Roadmap

Übergreifendes Master-Dokument zum Stand aller Slices und Tranchen
(`LH-FA-PROJDOCS-003`). Wird laufend gepflegt und liegt deshalb dauerhaft
in `in-progress/`.

| Phase | Status | Beschreibung | Artefakt |
| ----- | ------ | ------------ | -------- |
| M0 Spec | Done | Lastenheft v0.1.0 (Sektionen 1–14, inkl. 4.11 Build-/CI-Infrastruktur, 4.12 Doku-Struktur) | [`spec/lastenheft.md`](../../../../spec/lastenheft.md) |
| M0 ADRs | Done | ADR-0001 Implementierungssprache Go | [`docs/plan/adr/0001-implementierungssprache-go.md`](../../adr/0001-implementierungssprache-go.md) |
| M1 Repo-Skeleton | Done | Multi-Stage Dockerfile, Makefile, .dockerignore, Repo-Layout (`LH-FA-BUILD-001..009`), Doku-Struktur (`LH-FA-PROJDOCS-001..003`), `u-boot --help` / `--version`-Stub | Commit `7da05c7` |
| M2 Architektur | Done | Hexagonale Architektur (`LH-FA-ARCH-001..003`), `spec/architecture.md`, ADR-0002, `internal/{hexagon,adapter}/`-Skeleton, depguard mit aktiven Schicht-Regeln (match nichts, bis erste Pakete in `./internal/...`) | Commit `9d191a5` + Review-Fixes |
| M2b SOLID-Lint | Done | SOLID-nahes Lint-Profil (`LH-QA-004` auf MVP gehoben), 5 Default-Linter + 24 SOLID-nahe Linter (inkl. `depguard`), `docs/user/quality.md`, ADR-0003 | Commit `365e532` + Review-Fixes |
| M2c CI | Done | GitHub-Actions-CI (`LH-QA-003` auf konkret gehoben), `.github/workflows/ci.yml` mit Jobs `gates` + `security-gates` (beide PR-blockierend), SHA-pinned Actions, Docker-only, ADR-0004 | Commit `9a74e35` |
| M2d Carveouts | Done | Carveout-Disziplin (`LH-FA-PROJDOCS-005` MVP-Pflicht), Master-Inventar [`carveouts.md`](carveouts.md), 7 neue Slice-Pläne in [`open/`](../open/) für offene Carveouts; permanente Carveouts dokumentiert | dieser Commit |
| M3 `u-boot init` | In progress | Projektstruktur erzeugen (`LH-FA-INIT-001..007`), `u-boot.yaml` schreiben, Git-Init; Coverage-Carveout bereits aufgelöst, depguard-Verifikation folgt mit T5. Detail: [`slice-m3-init-flow.md`](slice-m3-init-flow.md). **Stand:** T1 ✅ `132d1a1` + `f5c784a`; T2 ✅ `aaf4d8d` + Review-Fixes; T3/T4/T5 offen |
| M4 `u-boot doctor` | Open | Lokale Voraussetzungen prüfen (`LH-FA-DIAG-001..004`), JSON-Output (`LH-NFA-USE-004`) | offen |
| M5 `u-boot add postgres` | Open | PostgreSQL-Add-on (`LH-FA-ADD-001..005`), Compose-Block, `.env.example`-Block, Healthcheck | offen |
| M6 `u-boot up` / `down` | Open | Compose-Wrapper (`LH-FA-UP-001..004`), Healthcheck-Polling, `--timeout`, `--volumes` | offen |
| M7 `u-boot generate` | Open | `generate changelog`/`readme`/`env-example`/`devcontainer` (`LH-FA-GEN-001..005`) | offen |
| M8 `u-boot config` | Open | `config get`/`set`/Anzeigen (`LH-FA-CONF-001..005`), Schema-Validierung | offen |
| MVP-Closure | Open | Devcontainer-Mindestumfang (`LH-FA-DEV-001..005`), MVP-Acceptance-Flows (`LH-AK-001..002`, `LH-AK-005..007`) | offen |
| V1 Keycloak / OTel | Open | `LH-FA-ADD-003`, `LH-FA-ADD-004`, `LH-AK-003`, `LH-AK-004` | offen |
| V1 Templates | Open | `LH-FA-TPL-001..004` | offen |
| V1 Logs / Dry-Run / Diff | Open | `LH-FA-UP-005`, `LH-FA-CLI-007/008` | offen |
| Later Migration / Custom Templates | Open | `LH-FA-CONF-006`, `LH-FA-TPL-003`, `LH-DA-004` | offen |

## Nächste Schritte

1. **M2** abschließen: Sektion 4.13 (Architektur) im Lastenheft + ADR-0002 (Hex-Pattern) + `spec/architecture.md` (Import-Graph, Layer-Regeln) + `internal/`-Skeleton (`hexagon/{domain,application,port/{driving,driven}}/` + `adapter/{driving,driven}/`).
2. **M3** starten: erster fachlicher Slice, `slice-m3-init-flow.md` in `next/` anlegen, dann nach `in-progress/`.

## Lifecycle-Hinweis

Diese Datei ist die einzige zulässige Ausnahme von der `slice-`/`tranche-`-Konvention für Dateinamen in `docs/plan/planning/` (siehe `LH-FA-PROJDOCS-003` und [`../README.md`](../README.md)).
