# Release-Ablauf — u-boot

| Dokument    | Release-Anleitung für Maintainer |
| ----------- | -------------------------------- |
| Projektname | `u-boot` |
| Bezug       | [`.github/workflows/publish.yml`](../../.github/workflows/publish.yml), [ADR-0007](../plan/adr/0007-distributionswege-ghcr.md), [ADR-0016](../plan/adr/0016-homebrew-distribution-per-tap.md) |
| Zielgruppe  | Maintainer, die einen Release schneiden |

## Zweck

Beschreibt, wie ein Release von `u-boot` entsteht: was vorab zu prüfen ist, was der
Tag-Push auslöst, was danach zu kontrollieren ist und was bei einem Fehler zu tun ist.
Ein Release ist **öffentlich** (GHCR-Image, Binaries, Homebrew-Tap) und lässt sich nicht
sauber zurücknehmen — der Tag-Push ist deshalb bewusst ein von Hand ausgeführter Schritt.

## 1. Was ein Tag auslöst

Der Push eines Tags `vMAJOR.MINOR.PATCH[-PRERELEASE]` startet den Workflow `publish`:

| Schritt | Ergebnis |
| --- | --- |
| SemVer-Prüfung | bricht vor Login/Build/Push ab bei ungültigen Tags (`v1.2`, `vNext`, Build-Metadaten mit `+`) |
| Runtime-Image bauen | `make build VERSION=<x.y.z>` (Docker-only) |
| OCI-Labels und `--version` prüfen | das Image meldet genau die Tag-Version, sonst Abbruch |
| GHCR-Push | `ghcr.io/pt9912/u-boot:<x.y.z>`; `:latest` **nur** für stabile Tags |
| Binaries | sechs Plattformen (Linux/macOS/Windows × amd64/arm64) als Release-Assets |
| `SHA256SUMS` | Prüfsummen der sechs Binaries als weiteres Asset |
| Homebrew-Formel | nur stabile Tags: `u-boot.rb` aus den Digests von `SHA256SUMS` |
| Job `tap` | nur stabile Tags: zieht die Formel nach `pt9912/homebrew-u-boot` nach |

Ein **Vorab-Tag** (`v0.6.0-rc.1`) durchläuft dieselbe Kette, wird aber als Prerelease
markiert, bekommt kein `:latest`, keine Formel und keinen Tap-Commit. Das ist die
Probe für den Release-Pfad (und hinterher löschbar).

## 2. Voraussetzungen (einmalig)

- **Secret `HOMEBREW_TAP_GITHUB_TOKEN`** im Repo `pt9912/u-boot` (Settings → Secrets and
  variables → Actions): ein PAT mit Schreibrecht (Contents: Read & Write) auf
  `pt9912/homebrew-u-boot`. Ohne Secret **überspringt** der Job `tap` den Nachzug mit einer
  Warnung; der Release selbst hängt nicht am Tap.
- **Tap-Repo** `pt9912/homebrew-u-boot` (README und Smoke-Workflow aus
  [`packaging/homebrew-tap/`](../../packaging/homebrew-tap/)).
- Schreibrecht auf Tags und auf `ghcr.io/pt9912/u-boot` (über `GITHUB_TOKEN` des Workflows).

## 3. Vorbereitung (ein Commit auf `main`)

Der Cut folgt dem Muster der Release-Cut-Slices (zuletzt
[`slice-v1-release-cut-v0.5.0`](../plan/planning/done/slice-v1-release-cut-v0.5.0.md)):

1. **CHANGELOG** (`CHANGELOG.md`): `## [Unreleased]` bleibt als leerer Anker; der bisherige
   Inhalt wandert unter `## [x.y.z] - <Tag-Datum>` mit kurzem Lead-Absatz. Compare-Links am
   Fuß umstellen (`[Unreleased]: …/compare/vx.y.z...HEAD`, neue Zeile `[x.y.z]: …`).
2. **Versionsstrings** `<alt>-dev` → `<neu>-dev` an genau drei Stellen:
   `cmd/uboot/main.go` (`var version`), `Makefile` (`VERSION ?=`), `Dockerfile`
   (`ARG UBOOT_VERSION=`), jeweils samt Kommentar-Erwähnungen. Der Tag-Build überschreibt
   den Wert per `UBOOT_VERSION`; der `-dev`-Stand gilt für lokale Builds.
3. **READMEs** (`README.md`, `README.de.md`): Status-Block und eine Zeile in der
   Releases-Tabelle.
4. **Roadmap** (`docs/plan/planning/in-progress/roadmap.md`): die Welle als aktiv führen,
   bis der Tag gesetzt ist; danach nach „Abgeschlossene Wellen".
5. **Release-Cut-Slice** nach dem Muster der früheren Cuts anlegen und nach dem Tag nach
   `done/` verschieben (Planning-Sensor: der Ruhe-Marker der Roadmap gilt nur, wenn kein
   Slice in `in-progress/` liegt).

## 4. Vorab-Prüfung

```bash
make gates            # lint + test + coverage-gate + docs-check
make ci               # gates + govulncheck + image-scan (Trivy HIGH/CRITICAL)
make test-docker      # Integrationstests gegen die echte Docker Engine
make fullbuild        # ci + Runtime-Image-Build
```

Zusätzlich: `main` ist auf `origin/main` gepusht und der **letzte CI-Lauf auf `main` ist
grün** (`gh run list --limit 3`). Ein roter `image-scan` blockiert den Release — typischer
Anlass ist ein neuer Go-stdlib-Advisory; dann den Go-Pin heben (`Dockerfile` `ARG
GO_VERSION`, `Makefile` `GO_VERSION`), siehe CHANGELOG-Einträge zu Toolchain-Bumps.

## 5. Tag setzen (Maintainer-Aktion)

```bash
git switch main && git pull --ff-only
git tag vX.Y.Z
git push origin vX.Y.Z      # löst publish aus
```

Optional zuerst ein Vorab-Tag (`vX.Y.Z-rc.1`), um die Kette ohne `:latest`, Formel und
Tap-Commit zu proben; er lässt sich danach löschen (`gh release delete`, `git push --delete
origin <tag>`, GHCR-Tag entfernen).

## 6. Nach dem Push kontrollieren

```bash
gh run watch                                   # publish-Lauf (Job publish, dann tap)
gh release view vX.Y.Z --json assets -q '.assets[].name'
docker run --rm ghcr.io/pt9912/u-boot:X.Y.Z --version
```

Erwartet: sechs Binaries, `SHA256SUMS`, bei stabilen Tags `u-boot.rb`; `--version` gleich
`X.Y.Z`; `:latest` zeigt auf dieselbe Version (nur stabil). Der Job `tap` ist grün (oder
meldet die Warnung zum fehlenden Secret). Danach auf einem Mac:

```bash
brew tap pt9912/u-boot
brew trust pt9912/u-boot      # Homebrew 5 verlangt das für Drittanbieter-Taps
brew install u-boot
u-boot --version              # X.Y.Z
```

Im Tap-Repo läuft nach dem Formel-Commit der Smoke-Workflow (macOS: Install, `--version`,
`init`; `doctor` ohne Docker nur auf „kein Absturz" geprüft).

## 7. Nacharbeit

- Roadmap: Welle nach „Abgeschlossene Wellen", Tag-Commit-Hash und Datum eintragen.
- Release-Cut-Slice nach `done/` (Closure-Notiz mit Hash und Sensoren).
- CHANGELOG-Datum auf das tatsächliche Tag-Datum prüfen.

## 8. Wenn etwas schiefgeht

| Symptom | Ursache / Maßnahme |
| --- | --- |
| `publish` bricht in „Validate SemVer tag" | Tag ungültig → Tag löschen, korrekt neu setzen (noch nichts veröffentlicht) |
| `--version` des Images ≠ Tag | Build-Arg-Drift → Abbruch vor dem Push; Versionsstrings/`make build VERSION` prüfen |
| Lauf bricht nach dem GHCR-Push | `:x.y.z` ist schon veröffentlicht; Workflow nach Behebung erneut ausführen (`--clobber` ersetzt Assets), nicht den Tag verschieben |
| Job `tap` warnt „secret not configured" | Secret setzen (§2), Job erneut ausführen (`gh run rerun --job`) |
| Job `tap` schlägt fehl | PAT-Rechte/Ablauf prüfen; manuell nachziehen: `TAP_TOKEN=… TAG=vX.Y.Z scripts/tap-nachzug.sh` |
| Formel fehlt im Tap | Asset `u-boot.rb` des Releases fehlt (Prerelease?) oder Tap-Job nicht gelaufen |
| Trivy-`image-scan` rot | Go-Pin heben (siehe §4), erst dann taggen |

Ein bereits veröffentlichter Tag wird **nicht verschoben**; ein Fehler wird mit dem nächsten
Patch-Release behoben.
