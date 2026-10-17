# Review-Report: slice-v1-release-cut-v0.10.0 (smoke fix) — 2026-10-17 <!-- d-check:ignore (Slice-Kennung im Titel) -->

**Review-Art:** Code — Diff gegen Plan, Entscheidungen und Hard Rules (Maintainability); keine DoD-Prüfung (Verifier).

**Gegenstand:** Commit `6e397e0` (`git show HEAD`) im Slice [`slice-v1-release-cut-v0.10.0`](../plan/planning/done/slice-v1-release-cut-v0.10.0.md): `.github/workflows/publish.yml`, `scripts/package-asset-name.sh`, `cmd/uboot/package_asset_name_test.go`, `docs/maintainer/releasing.md`, Slice §1/§3/§6. Anlass: `publish`-Lauf 39049650475 (Tag `v0.10.0-rc.1`), beide Paket-Smoke-Jobs mit `curl: (22) 404`.

**Skill:** `.harness/skills/reviewer.md` @ `6e397e0` · <!-- d-check:ignore (Skill-Pfad) -->
**Modell:** Sonnet 5.5 · **Datum:** 2026-10-17

**Eingangs-Kontext:**

- Slice [`slice-v1-release-cut-v0.10.0`](../plan/planning/done/slice-v1-release-cut-v0.10.0.md) (§1 Befund, §3 Plan-Tabelle, §6 Risiken)
- [`ADR-0017`](../plan/adr/0017-linux-pakete-nfpm.md) (Linux-Pakete per nfpm)
- `docs/maintainer/releasing.md`, `scripts/build-packages.sh`, Workflow von `v0.9.0` (`git show v0.9.0:.github/workflows/publish.yml`)
- `AGENTS.md` §3.6, §3.7, §3.9, §3.10

Sensors dieses Laufs: `make test` (Basis grün, Exit 0; vier Mutationen, siehe unten), `make docs-check` (nach dem Report). Nicht gelaufen: `make gates`/`make lint` (Diff berührt keine Lint-Regel; der Verifier fährt `make gates`), `make test-docker`, `make verify-depguard` (kein Docker-Pfad, keine Importregel berührt). Der Workflow selbst läuft nur im echten Tag-Lauf; `gh run view 39049650475` bestätigt nur den Ausgangsbefund (`publish` success, beide Smoke-Jobs failure, `tap` skipped).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Auflösung ist im Workflow nicht fail-closed: `echo "deb_asset=$(scripts/package-asset-name.sh deb <<<"$NAMES")" >> "$GITHUB_OUTPUT"` liefert den Exit-Status von `echo`, nicht des Skripts; `set -euo pipefail` greift nicht. Lokal nachgestellt (Skript mit Exit 1 bei 0 Treffern): der Schritt läuft mit Exit 0 weiter und schreibt `deb_asset=` leer. | Maintainability (Plan §3: „genau ein Treffer, sonst Exit 1“) | `.github/workflows/publish.yml` · `echo "deb_asset=$(scripts/package-asset-name.sh deb <<<"$NAMES")"` (ebenso `rpm_asset`) | nein — nur im Tag-Lauf; Shell-Verhalten lokal reproduziert | Exit-Status in Command-Substitution verschluckt |
| F-2 | INFO | `SHA256SUMS` wird aus `packages/` mit den Dateinamen von nfpm erzeugt (`~` bei Vorab-Tags), die Assets heißen mit `.`. Bei Vorab-Tags passt keine der vier Paket-Zeilen zu einem Asset-Namen: `sha256sum -c` meldet „No such file“, mit `--ignore-missing` bleiben die Pakete ungeprüft. Bei stabilen Tags sind beide Namen gleich; die Formel liest nur die vier Brew-Binaries. Der Slice nennt den Befund als offen (§6), `releasing.md` §1 ebenfalls. | Maintainability | `.github/workflows/publish.yml` · `(cd packages && sha256sum ./* \| sed 's# \./# #')` | nein — Befund aus Skript und Workflow gelesen; nfpm nicht erneut gefahren | Prüfsummen-Dateiname weicht vom Asset-Namen ab |
| F-3 | INFO | Der Test deckt die Auswahl (Vorab- und stabile Liste, arm64/aarch64 als Distraktor) und Fail-closed des Skripts ab, aber weder den Workflow-Schritt (F-1) noch den Fall „kein rpm“ (nur „kein deb“, „zwei deb“, „unbekanntes Format“). Der rpm-Zählzweig ist durch `Selects` mitgeprüft (Mutation M3), der Null-Treffer-Zweig nur über deb. | Maintainability | `cmd/uboot/package_asset_name_test.go` · `TestPackageAssetName_FailsClosed` | ja — `make test` | Test ohne Workflow-Verdrahtung |

Kein Finding wurde auf Zuruf herabgestuft; es gab keinen Implementer-Widerspruch.

## Prüfung im Einzelnen

**(1) Workflow Schritt für Schritt** — geprüft, bis auf F-1 ohne Befund:

- Reihenfolge im Job `publish`: Upload der Binaries → `Build and upload Linux packages` (`gh release upload "v$VERSION" packages/* --clobber`) → `Resolve the package asset names` → `Upload SHA256SUMS and Homebrew formula`. Der Auflösungsschritt liegt nach dem Paket-Upload; das Release existiert (der Upload-Schritt davor legt es an oder nutzt das vorhandene), ist weder Draft noch versteckt (`--prerelease` bei Vorab-Tags, kein `--draft`), `gh release view` ist per `GH_TOKEN: ${{ github.token }}` und Job-`permissions: contents: write` lesbar. Der Schritt hat `actions/checkout` im selben Job (das Skript existiert im Arbeitsbaum, Modus `100755`).
- `id: assets` und `steps.assets.outputs.deb_asset`/`rpm_asset` stimmen mit den Namen der Job-`outputs:` überein; `needs: publish` steht in beiden Smoke-Jobs; die Smoke-Jobs lesen `needs.publish.outputs.*_asset` über `env:`, nicht in `run:`. Eine Injektionsfläche über `${{ }}` in Skriptkörpern entsteht nicht; der Asset-Name stammt aus dem eigenen Release, das Skript lässt nur eine Zeile (anchored, genau ein Treffer) durch, ein Zeilenumbruch im `$GITHUB_OUTPUT` ist damit ausgeschlossen.
- Skript: `set -euo pipefail`; `grep -E` mit `\.` für die Punkte, `~`/`.`-Schreibweise egal (beide matchen `.*`); `grep … || true` und `grep -c . || true` sind bei null Treffern korrekt (Zählwert `0`, Exit 1 des Skripts, nachgestellt). `^u-boot_.*_amd64\.deb$` und `^u-boot-.*-1\.x86_64\.rpm$` schließen arm64/aarch64 aus (Mutationen M1, M4); `u-boot-linux-amd64` und `u-boot.rb` matchen keines der Muster.
- rpm-Smoke im Fedora-Container: `gh` wird dort nicht benutzt. Der Name wird im Job `publish` aufgelöst und nur als String-Output übergeben; Smoke lädt per `curl` die öffentliche Asset-URL. Bestätigt.

**(2) Unverändert (nichts gelockert)** — geprüft, ohne Befund: Install-Schritt (`apt-get install -y`/`dnf install -y` auf die geladene Datei), `u-boot --version | grep -F "$VERSION"` (das Binary trägt die Tag-Version, vom Schritt „matches tag“ im Job belegt), `u-boot init demo --no-git && test -f u-boot.yaml`, `needs`-Kette, `is_stable`-Logik und `tap`-Bedingung, Formel-Upload. Gegenüber dem `v0.9.0`-Workflow ändert dieser Commit in den Smoke-Jobs nur Env/Asset-Name; für stabile Tags sind Datei- und Asset-Name gleich, die gewählten Namen sind dieselben wie die frühere Tag-Ableitung (`u-boot_<v>_amd64.deb`, `u-boot-<v>-1.x86_64.rpm`). Kein Gate, keine Schwelle berührt (`AGENTS.md` §3.6).

**(3) Mutationen** (jeweils `make test`, danach `git checkout -- scripts/package-asset-name.sh`; Basis Exit 0):

| Nr. | Mutation | Ergebnis |
|---|---|---|
| M1 | deb-Muster `amd64` → `arm64` | rot (`Selects`: arm64-deb gewählt) |
| M2 | Zählprüfung `-ne 1` → `-lt 1` | rot (`FailsClosed`: „two deb assets must fail“) |
| M3 | rpm-Muster auf `^u-boot-.*.*\.rpm$` aufgeweicht | rot (`Selects`: „found 2“) |
| M4 | rpm `x86_64` → `aarch64` | rot (`Selects`) |

Arbeitsbaum danach sauber (`git status --short` leer). Der Test ist nicht tautologisch: die erwarteten Namen stehen literal, die Listen enthalten Distraktoren.

**(4) SHA256SUMS** — siehe F-2. Für stabile Tags gilt der Befund nicht; das Benutzerhandbuch verweist für `v0.10.0` auf `SHA256SUMS` (stabiler Tag, nicht betroffen).

**(5) Doku und Plan-Tabelle** — `releasing.md` §1 (`0.6.0~rc.1` Datei, `0.6.0.rc.1` Asset, „Stabile Tags haben nur eine Schreibweise“) und die Zeile im Abschnitt Fehlerbilder (404) stimmen mit Skript, Workflow und Slice §1 überein. Plan-Tabelle §3 gegen `git show --stat HEAD`: fünf Dateien im Diff, `publish.yml`, `package-asset-name.sh`, Test, `releasing.md` mit Eintrag, der Slice selbst ist Träger der Tabelle; keine Zeile ohne Diff, keine Diff-Datei ohne Eintrag (`AGENTS.md` §3.10). §3.9: (a) kein Kommentar beschreibt Entferntes („nicht aus dem Tag abgeleitet“ in Workflow-, Skript- und Testkommentar ist eine Abgrenzung, kein abwesender Text); (b) keine Kommentare mit Verweisen auf umbenannte Symbole; (c) kein Helfer zwischen Doc-Kommentar und Funktion (`runAssetName` trägt seinen Block vor der Funktion).

**(6) Würde rc.2 voraussichtlich grün durchlaufen?** Ja, wenn die folgenden Stellen im echten Lauf halten; kein Gate und kein lokaler Lauf prüft sie:

1. `gh release view "v$VERSION" --json assets -q '.assets[].name'` im Job `publish` liefert nach dem Paket-Upload bereits die vier Pakete (Upload ist synchron; nicht lokal belegbar).
2. Die Assets heißen tatsächlich `u-boot_0.10.0.rc.1_amd64.deb` und `u-boot-0.10.0.rc.1-1.x86_64.rpm` (Quelle: `gh release view v0.10.0-rc.1`, Slice §1) — der neue rc.2 ersetzt `rc.1` in den Namen.
3. Die Job-Outputs `deb_asset`/`rpm_asset` kommen in den Smoke-Jobs an (Syntax geprüft, Laufzeit nicht).
4. `dnf`/`apt-get` installieren eine Datei mit `.` in der Versionsangabe des Dateinamens (die Paket-Metadaten tragen `~`; der Dateiname spielt für die Installation keine Rolle).
5. Berechtigung `contents: write` im Job `publish` reicht für `gh release view`; nicht im echten Lauf geprüft.

F-1 ändert das Ergebnis des grünen Pfads nicht; er trifft den Fehlerpfad (siehe Verdikt).

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `.github/workflows/publish.yml` Reihenfolge, Berechtigungen, Outputs, `needs` | geprüft, ohne Befund |
| `.github/workflows/publish.yml` Injektionsfläche (`${{ }}` in `run:`) | geprüft, ohne Befund |
| `scripts/package-asset-name.sh` Muster, Zählung, Format-Argument | geprüft, ohne Befund |
| `cmd/uboot/package_asset_name_test.go` Aussagekraft | geprüft, ohne Befund (siehe F-3, INFO) |
| Smoke-Jobs: Install, `--version`, `init` | geprüft, ohne Befund (unverändert gegenüber `v0.9.0`) |
| `tap`-Job, `is_stable`, Formel-Upload | geprüft, ohne Befund (unverändert) |
| `docs/maintainer/releasing.md` | geprüft, ohne Befund |
| Plan-Tabelle §3 gegen Diff, §3.9 (a)/(b)/(c) | geprüft, ohne Befund |
| Echter Tag-Lauf (Outputs, `gh release view` im Job, Berechtigungen) | nicht geprüft: läuft nur nach Tag-Push; Stellen oben unter (6) |
| `make gates`, `make lint`, `make test-docker`, `make verify-depguard` | nicht gelaufen: Diff berührt weder Lint-Regel noch Docker-Pfad noch Importregel; `make gates` ist Verifier-Aufgabe |
| nfpm-Dateinamen mit `~` | nicht erneut gebaut; Aussage aus Slice §1 (lokal nachgebaut) übernommen und gegen `build-packages.sh` (nfpm wandelt die Vorabversion) gelesen |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Exit-Status in Command-Substitution verschluckt · Prüfsummen-Dateiname weicht vom Asset-Namen ab · Test ohne Workflow-Verdrahtung

## Verdikt

**Merge-blockierend:** ja (ein MEDIUM, F-1). Begründung: Der grüne Pfad ist plausibel, aber ein Fehlschlag der Auflösung wird erst im Smoke-Job als `curl`-/`apt`-Fehler sichtbar statt im Schritt selbst; jeder weitere Vorab-Tag ist ein echter Lauf mit Kosten, der Fehlerpfad sollte vorher stimmen.

**Übergabe:** F-1 an den Implementer; F-2 ist im Slice bereits als offen geführt (Entscheidung Projektinhaber). Die Finding-Klassen gehen in die Slice-Closure §7. Dieser Report ersetzt keine Verifikation.
