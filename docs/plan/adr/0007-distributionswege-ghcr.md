# ADR 0007: Distributionswege — GHCR primär, andere Wege vertagt/verworfen

## Status

Accepted

## Datum

2026-05-31

## Kontext

`LH-OPEN-002` Paketierung (siehe `spec/lastenheft.md` §14) ist seit
M0 offen. Mit dem Abschluss des MVP (Stand `e0d6c87`) und dem
bevorstehenden ersten Release-Schnitt (`v0.1.0` über
[`slice-v1-release-pipeline`](../planning/done/slice-v1-release-pipeline.md))
muss die Frage konkret beantwortet werden, bevor
`.github/workflows/publish.yml` geschrieben wird.

Sechs Optionen stehen laut `spec/lastenheft.md` §14 zur Wahl:

1. einzelnes Binary
2. Container Image
3. npm package
4. pip package
5. Homebrew
6. Debian/RPM-Paket

Rahmenbedingungen aus dem bisherigen Projektstand:

- u-boot **orchestriert** Docker (`LH-FA-UP-001`, M5/M6); Docker ist
  faktisch Laufzeit-Voraussetzung, nicht nur Bau-Voraussetzung.
- Build ist Docker-only (`LH-FA-BUILD-007`), Runtime-Image existiert
  bereits als Mehr-Stufen-Dockerfile mit OCI-Labels (`LH-FA-BUILD-002`).
- `docker/login-action` + `docker/build-push-action` für GHCR sind
  in der Go-CLI-Domäne (`gh`, `helm`, `kubectl`) Standard, SHA-pinbar
  und benötigen kein zusätzliches Secret außer dem automatischen
  `GITHUB_TOKEN` mit `packages: write`.
- Solo-Projekt-Status (siehe `docs/user/branch-protection.md`); jeder
  zusätzliche Distributionsweg ist Wartungs-Overhead pro Release.

Vorlagen:

- `gh` (cli/cli): GHCR + Homebrew + Debian/RPM + Binary über
  goreleaser — Multi-Channel, aber Team-Projekt mit dedizierter
  Release-Engineer-Rolle.
- `kubectl`: GHCR (via krel/cloud-build) + Binary; kein npm/pip,
  kein Homebrew als Erstanlauf.
- `k-deskflight` (Referenzprojekt): GHCR-only, Binary über Tag-Asset
  bei Bedarf nachgereicht.

## Entscheidung

**Container Image über GHCR** (`ghcr.io/pt9912/u-boot`) ist der
primäre und für `v0.1.0` einzige Distributionsweg. Die
[`slice-v1-release-pipeline`](../planning/done/slice-v1-release-pipeline.md)-Tranche
T2 liefert die Workflow-Mechanik, T3 den Trivy-Scan.

Über die sechs Optionen aus §14:

| Option | Entscheidung | Begründung |
| ------ | ------------ | ---------- |
| Container Image (GHCR) | **Gewählt** | Geringster Reibungspunkt, passt zu Docker-only-Build und Docker-Laufzeit-Voraussetzung; OCI-Labels schon da; SHA-pinbare Standard-Actions. |
| Einzelnes Binary | **Gewählt** (zusätzlich zu GHCR, ab v0.1.1) | Trigger ist am 2026-05-31 materialisiert: der v0.1.1-doctor-container-Befund (siehe §Folgepunkte) zeigte, dass `doctor` als Host-Diagnostik-Subkommando eine host-native Form braucht — das distroless-Image bringt keine `docker`/`git`-Binaries für die LH-FA-DIAG-002-Probes mit. Geliefert via [`slice-v2-binary-distribution`](../planning/done/slice-v2-binary-distribution.md): sechs Plattformen (Linux/macOS/Windows × amd64/arm64), `make build-binaries` cross-kompiliert im pinned `golang:$(GO_VERSION)`-Container, `publish.yml` hängt sie als GitHub-Release-Asset an jeden `v*`-Tag. |
| npm package | **Verworfen** | Sprach-Ökosystem-Mismatch (Go-Binary ist kein Node-Modul). npm/yarn-Distribution hätte nur Mehrwert bei Frontend-/JS-Team-Targets; u-boot ist Tooling für beliebige Stacks. |
| pip package | **Verworfen** | Analog npm: u-boot ist Go, kein Python-Modul. `pipx`-Wrapper für Go-Binaries hat keinen erkennbaren Mehrwert über GHCR. |
| Homebrew | **Vertagt** mit eigenem Slice-Trigger | Natürlicher Folgeweg nach Binary; ohne Binary keine Homebrew-Formula. Trigger: macOS-Nutzer-Nachfrage. Trigger-Slice: `slice-v2-homebrew-formula.md` (zu erstellen). |
| Debian/RPM | **Vertagt** mit eigenem Slice-Trigger | Hoher Tooling-Overhead (`debhelper`, `rpmbuild`, Repository-Hosting); nur sinnvoll bei konkreter Distro-Anfrage. Trigger-Slice: `slice-v2-distro-pakete.md` (zu erstellen). |

Konkrete Setzungen für den GHCR-Weg:

- **Registry:** `ghcr.io` (kein Docker Hub — kein zusätzliches Konto,
  kein zusätzliches Secret).
- **Repository:** `ghcr.io/pt9912/u-boot` (folgt Repo-Owner-Pfad).
- **Sichtbarkeit:** `public` (Default für Open-Source-Tooling).
- **Tag-Schema:**
  - `vMAJOR.MINOR.PATCH` → `:MAJOR.MINOR.PATCH` plus `:latest`.
  - `vMAJOR.MINOR.PATCH-<prerelease>` (z. B. `-rc.1`, `-alpha.2`) →
    nur `:MAJOR.MINOR.PATCH-<prerelease>`, kein `:latest`.
  - Andere `v*`-Tags (z. B. Build-Metadaten mit `+`) werden vom
    Workflow vor Login/Build/Push abgewiesen — `+` ist kein gültiges
    OCI-Tag-Zeichen.
- **Auth:** `GITHUB_TOKEN` mit `packages: write`. Kein PAT, kein
  Org-Secret.
- **OCI-Labels:** Aus `LH-FA-BUILD-002` im Dockerfile gesetzt; der
  Workflow lädt nichts nach.

## Konsequenzen

Positiv:

- **Eine** Distributionsmechanik für `v0.1.0`, statt mehrerer
  parallel zu wartender Pipelines.
- **GHCR ist OCI-konform** — Wechsel zu Docker Hub / Quay / privatem
  Mirror ist später ohne Schema-Bruch möglich.
- **`LH-OPEN-002` ist für den GHCR-Anteil entschieden;** die anderen
  Wege haben jeweils ein konkretes Trigger-Slice, sind also nicht
  mehr „offen ohne Plan", sondern „vertagt mit Trigger".
- **Versions-Pin verifiziert vor dem Push.** Mit
  `slice-v1-release-cut-v0.1.0` T1 (`056e4c6`) wird die Tag-`VERSION`
  konsistent durch alle drei Layer geführt: `-X main.version`
  im Go-Binary (`cmd/uboot/main.go:36`-Pattern), Build-Arg
  `UBOOT_VERSION` im Dockerfile, und der OCI-Label
  `org.opencontainers.image.version`. `publish.yml` pinnt vor dem
  GHCR-Push (1) Label gegen Tag-VERSION und (2) Live-`--version`-
  Smoke gegen Tag-VERSION; ein vergessener Build-Arg (Regression
  zu `0.1.0-dev`) bricht den Workflow vor dem Push.

Negativ / Trade-offs:

- **Docker-Pflicht für Endnutzer.** Wer u-boot nutzt, muss bereits
  Docker installiert haben — was wegen `LH-FA-DIAG-002` ohnehin
  geprüft wird. Für reines Lokal-CLI-Tooling ohne Docker-Use-Case
  gibt es bis zur Binary-Vertagung keinen Distributionsweg.
- **`:latest`-Floating-Tag.** Konventionsgemäß bequem, aber bei
  Breaking Changes mit Vorsicht zu nutzen. Wird in
  `docs/user/branch-protection.md` / Release-Doku adressiert, sobald
  ein Breaking Change ansteht.
- **`LH-OPEN-002` bleibt formal offen,** weil Binary / Homebrew /
  Distro-Pakete weiter aussehen — aber jeweils mit Slice-Plan und
  Trigger, statt als Plan-Loch. Carveouts-Eintrag mit T5 des
  Release-Pipeline-Slice (`bc487fc`) entsprechend reduziert.

Alternativen (verworfen):

- **Multi-Channel ab v0.1.0 (GHCR + Binary + Homebrew):** wäre der
  Stand vergleichbarer Tools (`gh`), kostet aber pro Release deutlich
  Wartungszeit ohne erkennbare Nutzer-Nachfrage als Solo-Projekt.
  Wird beim ersten konkreten Trigger nachgeholt.
- **Docker Hub statt GHCR:** zusätzliches Konto, zusätzliches Secret,
  Rate-Limits für Anonymous-Pulls; kein Vorteil gegenüber GHCR
  solange das Projekt auf GitHub liegt.

## Folgepunkte

- [`slice-v1-release-pipeline`](../planning/done/slice-v1-release-pipeline.md)
  T2 (`93b703e`) hat das hier festgelegte Tag-Schema und die
  GHCR-Settings im `publish.yml`-Workflow umgesetzt.
- `spec/lastenheft.md` §14 `LH-OPEN-002`-Abschnitt mit T1 dieses
  Slices (`0f64938`) auf den Stand „GHCR entschieden, Restwege
  vertagt" aktualisiert (analog `LH-OPEN-001`-Pattern nach ADR-0001).
- `docs/plan/planning/in-progress/carveouts.md` Zeile für
  `LH-OPEN-002` mit T5 (`bc487fc`) auf die verbleibenden Wege
  (Binary / Homebrew / Distro-Pakete) reduziert.
- Wenn ein vertagter Distributionsweg ausgelöst wird, wird der
  jeweilige `slice-v2-*`-Plan in `open/` neu angelegt und mit dem
  konkreten Trigger versehen.
- [`slice-v2-binary-distribution`](../planning/done/slice-v2-binary-distribution.md)
  hat den Trigger am 2026-05-31 materialisiert (via v0.1.1-
  doctor-container-Befund) und am 2026-06-01 vollständig
  geliefert (T1 `dc9a336` + `f3f1731` Cross-Compile-Makefile,
  T2 `5e5166b` `publish.yml`-Asset-Upload, T3 `866f6fd` README-
  Install-Block + CHANGELOG `## [Unreleased]`, T4 Slice-Closure).
  Binary-Zeile in §Entscheidung-Tabelle entsprechend von
  „Vertagt" auf „Gewählt (zusätzlich zu GHCR)" gehoben; weitere
  Tranche-Hashes in der DoD-Tabelle des done-Slice.
