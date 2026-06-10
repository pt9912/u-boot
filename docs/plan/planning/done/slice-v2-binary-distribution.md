# Slice V2: Binary-Distribution

## Auslöser

**[ADR-0007](../../adr/0007-distributionswege-ghcr.md) §Folgepunkte 1 hat am 2026-05-31 seinen ersten konkreten
Trigger gefunden:** der Doctor-Container-Awareness-Befund (siehe
[`slice-v0.1.1-doctor-container-awareness`](slice-v0.1.1-doctor-container-awareness.md))
zeigt, dass GHCR-Container-only für ein Subkommando-Set, das
Host-Diagnostik macht (`doctor`), praktisch nicht reicht. Bis
v0.1.1 wird `doctor` per Container-Skip entschärft; mittelfristig
braucht u-boot eine **lokale Binary**, mit der `doctor`,
`init` und die anderen Subkommandos ohne Container-Wrapping
laufen können.

[ADR-0007](../../adr/0007-distributionswege-ghcr.md) §Entscheidung hatte „Einzelnes Binary" ursprünglich als
**vertagt mit Trigger-Slice** markiert. Dieses Slice ist die
Trigger-Auflösung und hat die Entscheidung mit T4 auf
**„Gewählt (zusätzlich zu GHCR)"** gehoben.

## Aufhebungsbedingung

u-boot wird zusätzlich zum GHCR-Image als statisch-gelinktes
Binary für Linux, macOS und Windows (jeweils amd64 + arm64,
sechs Plattformen gesamt) distribuiert. Distributionsweg:
GitHub-Releases-Asset pro `v*`-Tag (ohne dedizierten Apt-/Yum-
Mirror — Homebrew und Debian/RPM bleiben eigene Trigger-Slices
aus [ADR-0007](../../adr/0007-distributionswege-ghcr.md) §Entscheidung). Windows ist mit dabei, weil das
u-boot-Binary nur `os.Exec` und Filesystem-IO macht (kein Linux-
Syscall); die [`LH-NFA-PORT-002`](../../../../spec/lastenheft.md#lh-nfa-port-002--keine-unnötigen-systemabhängigkeiten)-Permanent-Carveout-Zeile betrifft
den Build auf Windows (`make` fehlt), nicht das Ausführen einer
fertigen `.exe`.

## Akzeptanzkriterien

- ✅ `.github/workflows/publish.yml` bildet pro Tag `v*` zusätzlich
  zur GHCR-Image-Pipeline Binaries für sechs Plattformen
  (Linux/macOS/Windows × amd64/arm64) und hängt sie als
  Release-Assets an den GitHub-Release.
- ✅ README.{md,de.md} Quickstart bekommt einen „Install"-Block
  mit `curl`/PowerShell-`Invoke-WebRequest`-Beispielen vor dem
  `docker run`-Block (Binary-first; GHCR demoted auf „alternative
  — container/CI workflows").
- ✅ [ADR-0007](../../adr/0007-distributionswege-ghcr.md) §Entscheidung aktualisiert: „Einzelnes Binary" von
  „Vertagt" auf „Gewählt (zusätzlich zu GHCR)" gehoben mit
  Verweis auf diesen Slice.
- ✅ [`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung)-Restwege-Carveout in
  [`carveouts.md`](../in-progress/carveouts.md) auf Homebrew +
  Debian/RPM reduziert (Binary ausgeliefert).
- ✅ `make build-binaries` Make-Target lokal reproduzierbar;
  Inner-/Outer-Loop-Parität nach [`LH-FA-BUILD-007`](../../../../spec/lastenheft.md#lh-fa-build-007--docker-only-workflow).

## Tranchen

| T | Commit | Inhalt |
| - | ------ | ------ |
| T1 | `dc9a336` + `f3f1731` | `make build-binaries` Make-Target (Cross-Compile via `GOOS`/`GOARCH` im pinned `golang:$(GO_VERSION)`-Container, CGO disabled, `-ldflags "-s -w -X main.version=$(VERSION)"`, Output `bin/u-boot-<os>-<arch>[.exe]`). T1 (`dc9a336`) initial mit vier Plattformen (Linux/macOS × amd64/arm64); T1-follow (`f3f1731`) ergänzt Windows amd64+arm64 (`.exe`-Suffix). Insgesamt sechs Plattformen. |
| T2 | `5e5166b` | `.github/workflows/publish.yml` erweitert: nach GHCR-Push + OCI-Label-Verify + Live-`--version`-Smoke ruft der Workflow `make build-binaries` mit der Tag-`VERSION` im pinned golang-Container auf und uploadet alle sechs Plattform-Binaries via `gh release upload "$TAG" bin/u-boot-* --clobber` als Release-Assets. Build + Upload laufen AFTER GHCR-Push, damit ein gescheiterter Versions-Pin den Binary-Upload nicht in eine inkonsistente Release verleitet. Release wird falls nötig vorab erzeugt (Idempotenz). |
| T3 | `866f6fd` | README.{md,de.md} Quickstart binary-first umgestellt: neuer „Install pre-built binary (recommended)"-Block vor dem GHCR-Block. Linux/macOS via `curl -sSL` + `uname`-Auto-Detection, Windows via PowerShell `Invoke-WebRequest`. GHCR-Sektion demoted auf „alternative — container/CI workflows". Version-Pin-Beispiel `releases/download/v0.1.1/…` plus `latest/download/`-Caveat (v0.1.0 hatte noch keine Binary-Assets, `latest` greift erst ab v0.1.1). `CHANGELOG.md ## [Unreleased]` ergänzt um Cross-Platform-Binary-Distribution + Binary-First-Quickstart-Umstellung. README-v0.1.1-Status-Block bumped auf „T1 + T2 + T3 shipped". |
| T4 | dieser Commit | Slice-Plan nach `done/`; [ADR-0007](../../adr/0007-distributionswege-ghcr.md) §Entscheidung Binary-Zeile von „Vertagt" auf „Gewählt (zusätzlich zu GHCR)" gehoben mit Verweis auf diesen Slice und §Folgepunkte um den materialisierten Trigger ergänzt; `carveouts.md` [`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung)-Restwege auf Homebrew + Debian/RPM reduziert (Binary ausgeliefert); `roadmap.md` Carveout-Auflösungs-Slices-Tabelle Status `Open` → `Done` mit T3+T4-Update, §Nächste Schritte Punkt 5 Binary-Slice-Verweis von `open/` auf `done/` umgehängt. |

## Out of Scope

- Homebrew-Formula: eigener Slice
  [`slice-v2-homebrew-formula.md`](../open/slice-v2-homebrew-formula.md) ([ADR-0007](../../adr/0007-distributionswege-ghcr.md) §Entscheidung
  „Vertagt").
- Debian/RPM-Pakete: eigener Slice
  [`slice-v2-distro-pakete.md`](../open/slice-v2-distro-pakete.md).
- Signature-/SBOM-Verifikation für Binaries: Folge-Slice bei
  konkretem Bedarf (z. B. Reproducible-Builds-Anfrage oder
  Supply-Chain-Anforderung von einem Konsumenten).

## Bezug

- Auslöser:
  [`slice-v0.1.1-doctor-container-awareness`](slice-v0.1.1-doctor-container-awareness.md)
  hat den ersten konkreten Bedarf gezeigt
  (`doctor` für Host-Diagnostik braucht eine Binary).
- ADR: [ADR-0007](../../adr/0007-distributionswege-ghcr.md)
  §Entscheidung (Binary-Zeile auf „Gewählt" gehoben) +
  §Folgepunkte (Trigger materialisiert).
- Carveout-Inventar:
  [`carveouts.md`](../in-progress/carveouts.md) →
  [`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung)-Restwege auf Homebrew + Debian/RPM reduziert.
- Phase: V2 (nach v0.1.1).
