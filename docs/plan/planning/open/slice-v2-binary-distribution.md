# Slice V2: Binary-Distribution

## Auslöser

**ADR-0007 §Folgepunkte 1 hat am 2026-05-31 seinen ersten konkreten
Trigger gefunden:** der Doctor-Container-Awareness-Befund (siehe
[`slice-v0.1.1-doctor-container-awareness`](../done/slice-v0.1.1-doctor-container-awareness.md))
zeigt, dass GHCR-Container-only für ein Subkommando-Set, das
Host-Diagnostik macht (`doctor`), praktisch nicht reicht. Bis
v0.1.1 wird `doctor` per Container-Skip entschärft; mittelfristig
braucht u-boot eine **lokale Binary**, mit der `doctor`,
`init` und die anderen Subkommandos ohne Container-Wrapping
laufen können.

ADR-0007 §Entscheidung hat „Einzelnes Binary" explizit als
**vertagt mit Trigger-Slice** markiert. Dieses Slice ist die
Trigger-Auflösung.

## Aufhebungsbedingung

u-boot wird zusätzlich zum GHCR-Image als statisch-gelinktes
Binary für Linux/macOS (mind. amd64 + arm64) distribuiert. Der
Distributionsweg ist konkret entschieden — entweder GitHub-
Releases-Asset (am einfachsten, ohne weitere Infrastruktur) oder
ein dedizierter Apt-/Yum-Mirror (falls T2 aus
`slice-v0.1.1-doctor-container-awareness` weitere Trigger zieht).

## Akzeptanzkriterien

- `.github/workflows/publish.yml` (oder Schwester-Workflow)
  bildet pro Tag `v*` zusätzlich zur GHCR-Image-Pipeline auch
  Linux-amd64-, Linux-arm64- und macOS-amd64-/arm64-Binaries
  und hängt sie als Release-Assets an den GitHub-Release.
- README.{md,de.md} Quickstart bekommt einen „Install" -Block
  mit `curl … | sh`-/`wget`-Beispiel pro Plattform vor dem
  `docker run`-Block.
- ADR-0007 §Entscheidung wird aktualisiert: „Einzelnes Binary"
  von „Vertagt" auf „Gewählt (zusätzlich zu GHCR)" gehoben mit
  Verweis auf diesen Slice.
- `LH-OPEN-002`-Restwege-Carveout in
  [`carveouts.md`](../in-progress/carveouts.md) wird auf
  Homebrew + Debian/RPM reduziert (Binary ist dann ausgeliefert,
  ggf. mit `:latest`-Floating-Konvention für unstable).
- `make build-binaries` Make-Target (lokal reproduzierbar) wird
  ergänzt; Inner-/Outer-Loop-Parität nach `LH-FA-BUILD-007`.

## Tranchen (vorgeschlagen)

| T | Inhalt |
| - | ------ |
| T1 | `make build-binaries` Make-Target (cross-compile via `GOOS`/`GOARCH`); Versions-Pin via `UBOOT_VERSION` wie bisher. |
| T2 | `.github/workflows/publish.yml` erweitert: nach GHCR-Push + Verify auch Binaries bauen und `gh release upload` als Release-Assets. |
| T3 | README.{md,de.md} Install-Block + Quickstart auf Binary-First umgestellt. CHANGELOG `## [Unreleased]`-Sektion ergänzt. |
| T4 | ADR-0007-Update + carveouts.md Restwege-Reduktion + Slice-Closure. |

## Out of Scope

- Homebrew-Formula: eigener Slice
  `slice-v2-homebrew-formula.md` (ADR-0007 §Entscheidung
  „Vertagt").
- Debian/RPM-Pakete: eigener Slice
  `slice-v2-distro-pakete.md`.
- Signature-/SBOM-Verifikation für Binaries: optional in T2 oder
  Folge-Slice.

## Bezug

- Auslöser:
  [`slice-v0.1.1-doctor-container-awareness`](../done/slice-v0.1.1-doctor-container-awareness.md)
  hat den ersten konkreten Bedarf gezeigt
  (`doctor` für Host-Diagnostik braucht eine Binary).
- ADR: [ADR-0007](../../adr/0007-distributionswege-ghcr.md)
  §Entscheidung (Binary-Zeile „Vertagt") + §Folgepunkte 1 als
  Re-Evaluation-Trigger.
- Carveout-Inventar:
  [`carveouts.md`](../in-progress/carveouts.md) →
  `LH-OPEN-002`-Restwege; Binary-Anteil wird mit diesem Slice
  geschlossen.
- Phase: V2 (nach v0.1.1).
