# Verifikations-Bericht: slice-v1-release-cut-v0.10.0 (smoke fix, Review F-1 und F-3) — 2026-10-17 <!-- d-check:ignore (Slice-Kennung im Titel) -->

**Rolle:** Verifier — „Bauen wir es richtig?“ gegen Plan und DoD; keine Reparatur.

**Gegenstand:** Commit `5e0897f` (`git log --oneline 6e397e0..HEAD`) im Slice [`slice-v1-release-cut-v0.10.0`](../plan/planning/in-progress/slice-v1-release-cut-v0.10.0.md); Review: [`2026-10-17-slice-v1-release-cut-v0.10.0-smoke-fix`](2026-10-17-slice-v1-release-cut-v0.10.0-smoke-fix.md) (F-1 MEDIUM, F-3 INFO).

## Ergebnis

**GO** für Push und `v0.10.0-rc.2` aus Sicht der Verifikation; die Freigabe liegt beim Projektinhaber. F-1 und F-3 sind am Code behoben; F-2 ist offen und im Slice vermerkt.

## Sensor-Belege (selbst gefahren)

| Sensor | Ergebnis |
|---|---|
| `make gates` | Exit 0, Coverage 91,20 % (Schwelle 90 %) |
| `make doc-trace` | Exit 0, 83 Anforderungen, 0 Waisen |
| `make ci` | Exit 0 (`gates` + `govulncheck` + `image-scan`) |
| `make docs-check` | 0 Befunde (Teil von `gates`; nach Anlage dieses Berichts erneut gelaufen) |
| Mutation 1 (`make test`) | Leerprüfung entfernt und `\|\| true` an beiden Namensaufrufen: `TestPackageAssetOutputs` rot |
| Mutation 2 (`make test`) | `deb_asset` vor der rpm-Auflösung geschrieben: `TestPackageAssetOutputs` rot (doppelte Ausgabe, Ausgabe bei Fehler) |
| `make test-docker`, `make verify-depguard` | nicht gelaufen: kein Docker-Pfad, keine Importregel berührt |

Beide Mutationen wurden per `git checkout -- <Datei>` zurückgenommen; Arbeitsbaum vor und nach den Läufen sauber.

## (1) Schritt-Form nachgestellt

Lokal die Schritt-Form `set -euo pipefail` plus `scripts/package-asset-outputs.sh <<<"$NAMES" >> "$GITHUB_OUTPUT"` unter `bash -e` gefahren:

| Fall | Exit | `$GITHUB_OUTPUT` |
|---|---|---|
| 0 Treffer, leere Liste | 1 | leer |
| 2 deb, 2 rpm | 1 | leer |
| nur deb, nur rpm | 1 | leer |
| rc.2-Liste mit arm64/aarch64-Distraktoren | 0 | `deb_asset=u-boot_0.10.0.rc.2_amd64.deb`, `rpm_asset=u-boot-0.10.0.rc.2-1.x86_64.rpm` |
| stabile Liste mit `u-boot.rb` | 0 | `u-boot_0.10.0_amd64.deb`, `u-boot-0.10.0-1.x86_64.rpm` |

Das Skript schreibt erst nach der Prüfung beider Namen, ein Fehler hinterlässt nichts. Auch ein Fehler von `gh release view` in `NAMES="$(…)"` bricht den Schritt ab (nachgestellt, Exit 1).

## (2) Workflow ganz gelesen

Unverändert korrekt: Job-Outputs `deb_asset`/`rpm_asset` (`steps.assets.outputs`, `id: assets` passt), `is_stable`, `needs: publish` in beiden Smoke-Jobs, `env: ASSET` aus `needs.publish.outputs.*`, Install-Schritt, `u-boot --version | grep -F "$VERSION"`, `init`-Test, Tap-Job, Formel-Upload. Der Schritt liegt nach dem Paket-Upload. Die Differenz gegen `v0.9.0` besteht nur aus der dokumentierten Differenz (Outputs, Asset-Auflösung, ASSET statt tag-gebildeter Namen, `tap`-Bedingung über `is_stable`, SHA256SUMS über Pakete, Kommentarumstellung).

## (3) Tests

`TestPackageAssetName_NoRpmFails` und `TestPackageAssetOutputs` sind nicht tautologisch: Sie führen die echten Skripte in der Schritt-Form aus und prüfen Exit und Dateiinhalt; die zwei Mutationen oben fallen auf.

## (4) Plan-Konformität (§3.10)

Die Tabelle §3 gegen `git diff --name-only 630f427..HEAD` (20 Dateien): alle Dateien stehen drin, mit Ausnahme des Slice-Plans selbst und der Berichte `2026-10-17-slice-v1-release-cut-v0.10.0.md` und `2026-10-17-verification-slice-v1-release-cut-v0.10.0.md` (der Smoke-Fix-Report ist eingetragen). Kleine Lücke gegen §3.10, nicht blockierend. `releasing.md` (§1 `~`/`.`, Paket-Smoke-Zeile, §8-404-Zeile) stimmt mit Skripten und Workflow überein. Slice §6 nennt F-2 (`SHA256SUMS` führt `~`, Assets `.`) als offen und die Grenze „der YAML-Schritt läuft nur im Workflow“. <!-- d-check:ignore (Dateinamen der Berichte) -->

## Erst im echten rc.2-Lauf prüfbar

- `gh release view "v$VERSION" --json assets -q '.assets[].name'` liefert die Namen im erwarteten Format (Release existiert, `GH_TOKEN`-Rechte).
- Die Job-Outputs kommen bei den Smoke-Jobs an (`needs.publish.outputs.*`).
- `curl`-Download, `apt-get install` und `dnf install` des Vorab-Pakets, `u-boot --version | grep -F "0.10.0-rc.2"`, `init`-Test.
- Asset-Sichtbarkeit direkt nach dem Upload (Zeitverhalten).
- Die ganze Kette `publish` (GHCR, Binaries, Pakete, `SHA256SUMS`) mit `is_stable=false`, `tap` übersprungen.

## Offen, nicht blockierend

- F-2: `SHA256SUMS` führt bei Vorab-Tags `~`-Dateinamen, die Assets heißen mit `.`; bei stabilen Tags gleich, ohne Wirkung auf die Formel. Im Slice §6 und `releasing.md` §1 vermerkt.
