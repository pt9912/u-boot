# ADR 0007: Distributionswege — GHCR primär, andere Wege vertagt/verworfen

## Status

Accepted

## Datum

2026-05-31

## Kontext

[`LH-OPEN-002`](../../../spec/lastenheft.md#lh-open-002-paketierung) Paketierung (siehe `spec/lastenheft.md` §14) ist seit
M0 offen. Mit dem Abschluss des MVP (Stand `e0d6c87`) und dem
bevorstehenden ersten Release-Schnitt (`v0.1.0`) muss die Frage konkret
beantwortet werden, bevor `.github/workflows/publish.yml` geschrieben
wird.

Sechs Optionen stehen laut `spec/lastenheft.md` §14 zur Wahl:

1. einzelnes Binary
2. Container Image
3. npm package
4. pip package
5. Homebrew
6. Debian/RPM-Paket

Rahmenbedingungen aus dem bisherigen Projektstand:

- u-boot **orchestriert** Docker ([`LH-FA-UP-001`](../../../spec/lastenheft.md#lh-fa-up-001-umgebung-starten), M5/M6); Docker ist
  faktisch Laufzeit-Voraussetzung, nicht nur Bau-Voraussetzung.
- Build ist Docker-only ([`LH-FA-BUILD-007`](../../../spec/lastenheft.md#lh-fa-build-007-docker-only-workflow)), Runtime-Image existiert
  bereits als Mehr-Stufen-Dockerfile mit OCI-Labels ([`LH-FA-BUILD-002`](../../../spec/lastenheft.md#lh-fa-build-002-runtime-stage-pflichten)).
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
Workflow-Mechanik muss das unten definierte Tag-Schema und die
Security-Gates aus dem Lastenheft erfüllen.

Über die sechs Optionen aus §14:

| Option | Entscheidung | Begründung |
| ------ | ------------ | ---------- |
| Container Image (GHCR) | **Gewählt** | Geringster Reibungspunkt, passt zu Docker-only-Build und Docker-Laufzeit-Voraussetzung; OCI-Labels schon da; SHA-pinbare Standard-Actions. |
| Einzelnes Binary | **Gewählt** (zusätzlich zu GHCR, ab v0.1.1) | Trigger ist am 2026-05-31 materialisiert: der v0.1.1-doctor-container-Befund zeigte, dass `doctor` als Host-Diagnostik-Subkommando eine host-native Form braucht — das distroless-Image bringt keine `docker`/`git`-Binaries für die [`LH-FA-DIAG-002`](../../../spec/lastenheft.md#lh-fa-diag-002-lokale-voraussetzungen-prüfen)-Probes mit. Distribution umfasst sechs Plattformen (Linux/macOS/Windows × amd64/arm64). |
| npm package | **Verworfen** | Sprach-Ökosystem-Mismatch (Go-Binary ist kein Node-Modul). npm/yarn-Distribution hätte nur Mehrwert bei Frontend-/JS-Team-Targets; u-boot ist Tooling für beliebige Stacks. |
| pip package | **Verworfen** | Analog npm: u-boot ist Go, kein Python-Modul. `pipx`-Wrapper für Go-Binaries hat keinen erkennbaren Mehrwert über GHCR. |
| Homebrew | **Vertagt** mit Trigger | Natürlicher Folgeweg nach Binary; ohne Binary keine Homebrew-Formula. Trigger: macOS-Nutzer-Nachfrage. |
| Debian/RPM | **Vertagt** mit Trigger | Hoher Tooling-Overhead (`debhelper`, `rpmbuild`, Repository-Hosting); nur sinnvoll bei konkreter Distro-Anfrage. |

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
- **OCI-Labels:** Aus [`LH-FA-BUILD-002`](../../../spec/lastenheft.md#lh-fa-build-002-runtime-stage-pflichten) im Dockerfile gesetzt; der
  Workflow lädt nichts nach.

## Konsequenzen

Positiv:

- **Eine** Distributionsmechanik für `v0.1.0`, statt mehrerer
  parallel zu wartender Pipelines.
- **GHCR ist OCI-konform** — Wechsel zu Docker Hub / Quay / privatem
  Mirror ist später ohne Schema-Bruch möglich.
- **[`LH-OPEN-002`](../../../spec/lastenheft.md#lh-open-002-paketierung) ist für den GHCR-Anteil entschieden;** die anderen
  Wege haben jeweils konkrete Trigger, sind also nicht mehr „offen ohne
  Plan", sondern „vertagt mit Trigger".
- **Versions-Pin verifiziert vor dem Push.** Die Tag-`VERSION` wird
  konsistent durch alle drei Layer geführt: `-X main.version`
  im Go-Binary (`cmd/uboot/main.go:36`-Pattern), Build-Arg
  `UBOOT_VERSION` im Dockerfile, und der OCI-Label
  `org.opencontainers.image.version`. `publish.yml` pinnt vor dem
  GHCR-Push (1) Label gegen Tag-VERSION und (2) Live-`--version`-
  Smoke gegen Tag-VERSION; ein vergessener Build-Arg (Regression
  zu `0.1.0-dev`) bricht den Workflow vor dem Push.

Negativ / Trade-offs:

- **Docker-Pflicht für Endnutzer.** Wer u-boot nutzt, muss bereits
  Docker installiert haben — was wegen [`LH-FA-DIAG-002`](../../../spec/lastenheft.md#lh-fa-diag-002-lokale-voraussetzungen-prüfen) ohnehin
  geprüft wird. Für reines Lokal-CLI-Tooling ohne Docker-Use-Case
  gibt es bis zur Binary-Vertagung keinen Distributionsweg.
- **`:latest`-Floating-Tag.** Konventionsgemäß bequem, aber bei
  Breaking Changes mit Vorsicht zu nutzen. Wird in
  `docs/user/branch-protection.md` / Release-Doku adressiert, sobald
  ein Breaking Change ansteht.
- **[`LH-OPEN-002`](../../../spec/lastenheft.md#lh-open-002-paketierung) bleibt formal offen,** weil Binary / Homebrew /
  Distro-Pakete weiter aussehen — aber jeweils mit Trigger, statt als
  Plan-Loch.

Alternativen (verworfen):

- **Multi-Channel ab v0.1.0 (GHCR + Binary + Homebrew):** wäre der
  Stand vergleichbarer Tools (`gh`), kostet aber pro Release deutlich
  Wartungszeit ohne erkennbare Nutzer-Nachfrage als Solo-Projekt.
  Wird beim ersten konkreten Trigger nachgeholt.
- **Docker Hub statt GHCR:** zusätzliches Konto, zusätzliches Secret,
  Rate-Limits für Anonymous-Pulls; kein Vorteil gegenüber GHCR
  solange das Projekt auf GitHub liegt.

## Folgepunkte

- Wenn Homebrew oder Debian/RPM konkret nachgefragt werden, wird die
  jeweilige Distributionsentscheidung neu bewertet.
- Wenn ein weiterer Distributionsweg gewählt wird, muss [`LH-OPEN-002`](../../../spec/lastenheft.md#lh-open-002-paketierung)
  im Lastenheft auf den neuen Produktstand gebracht werden.
