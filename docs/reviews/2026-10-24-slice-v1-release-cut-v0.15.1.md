# Review-Report: slice-v1-release-cut-v0.15.1

* Rolle: Reviewer (`.harness/skills/reviewer.md`), frischer Kontext
* Gegenstand: `git diff 6335bf18..HEAD` (HEAD `40195a96`), Slice
  [`slice-v1-release-cut-v0.15.1`](../plan/planning/in-progress/slice-v1-release-cut-v0.15.1.md)
  (Patch `v0.15.1`, Einstufung *Fixed*, kein Advisory); Muster: Cuts `v0.15.0` und `v0.13.1`
* Bezug: [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-CONF-005`](../../spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern)
* Belege der Inhalte: Review- und Verify-Reports beider Fix-Slices vom 2026-10-24
* Nicht Gegenstand: DoD-Prüfung (Verifier)

## Urteil

Kein HIGH, kein MEDIUM. Ein LOW (Messort im Lead breiter als die Probe), drei INFO.

## Sensoren (selbst gefahren)

| Sensor | Ergebnis |
|---|---|
| `make gates` | Exit-Code 0 (`EXIT=0` im Log; `ablage-pruefung` 26 Fälle, 0 Fehlschläge als letzter Block) |
| `make test-docker` | nicht gelaufen: Auftrag (Speicher); der Diff berührt keinen Produktcode (eine Zeile `var version`) |
| `make verify-depguard` | nicht gelaufen: keine Importregel berührt |
| `make ci`, `make fullbuild` | nicht selbst gefahren; die Darstellung des Implementers (Trivy-DB 404, Wiederholung grün) ist plausibel, aber von mir nicht reproduziert |

## Findings

### LOW-1

* `kategorie`: LOW
* `quelle`: `AGENTS.md` §3.10 (Aussage nennt nur, was gemessen wurde), `BEO-DOK/zusage-weiter-als-messung`
* `pfad`: `CHANGELOG.md` · „**Measured:** both fixes at the built binary on one Linux host“; gleichlautend im Status-Block beider READMEs („both are measured at the built binary“)
* `klasse`: messort-breiter-als-probe
* `befund`: Der Lead nennt `ollama`, `gh` und `glab` für den Fix in `config set` und sagt, beide Fixes seien am gebauten Binary gemessen. Die Binary-Proben des Verify-Reports `config set` laufen mit `gh`, `glab` und `java`; `ollama` ist auf Service-Ebene durch Tests gepinnt, am Binary nicht gesehen.
* `risiko`: Leser nimmt `ollama` als am Binary gemessen; geringes Risiko, weil Regel und Test den Fall tragen (eine Regel in `featureVersionNeedsSource`, je Fall in beiden Testfunktionen).
* `verifizierbar`: Review-only

### INFO-1

* `befund`: Die DoD des Slice verlangt vor dem Tag „`make test-docker` ist gelaufen“; §3 Umsetzung Punkt 8 ersetzt das durch den CI-Lauf in M2 (Entscheidung des Auftrags). Der Haken ist offen und wird nicht gesetzt; Wortlaut der DoD und Ist-Stand laufen auseinander, bis Closure oder Planner sie angleichen. Zuständig: Verifier/Planner.

### INFO-2

* `befund`: Der Roadmap-Zeiger („Push und Tag stehen aus“) und die Marker „in Arbeit“ sind Zustandsmarker mit Ablaufpunkt; die Closure-Zeile der DoD nennt ihre Entfernung. Kein Befund jetzt; bei der Closure prüfen (§3.10).

### INFO-3

* `befund`: Datum 2026-10-24 gilt für Cut, Tabellenzeilen und Handbuch-`Stand`. Wird an einem anderen Tag getaggt, greift M0 (neuer Cut-Commit, M2 erneut). Zuständig: Hauptagent vor M1/M5.

## Negativbefunde (geprüft, ohne Befund)

* **CHANGELOG `[0.15.1]`**: Lead trägt nur belegte Sätze. Gegen Code gelesen: `featureVersionNeedsSource` (`devcontainer_features.go`) fragt `dockerfileBlock` (gh/glab/ollama), beide Pfade (`validateDevcontainerFeatures`, `revalidateFeatureEntry` in `config.go`) nutzen es; Freigabeprüfung in `generate.go` (`unlistedFeatureSources`, `ErrConfigValueInvalid` → Exit 10, vor Schreiben). Vorher-Verhalten, zwei Auswege, `--dry-run`/`--diff`, gesperrte Bestandsdatei, Ungemessenes decken sich mit den Verify-Reports. Keine Sicherheits-/Advisory-Aussage, kein CI-Satz, „one Linux host“ ohne `amd64` (Reports nennen keine Architektur), „not measured“ statt „nicht unterstützt“. `[Unreleased]` leer; Vergleichslinks `[Unreleased]: v0.15.1...HEAD` und `[0.15.1]: v0.15.0...v0.15.1` stimmen; Fixed-Einträge wörtlich unverändert, Text unter `[0.15.0]` unberührt.
* **Versionsstellen**: eigener Grep über den Baum nach `0.15.0` (ohne `.harness`, `planning`, `reviews`, `messungen`): Treffer nur CHANGELOG-Geschichte, die zwei `v0.15.0`-Tabellenzeilen, `releasing.md:57` (Zeiger „zuletzt“, Closure) und Spec-Versionen (`lastenheft.md:1591`, `spezifikation.md:1435`). `0.15.1`: `main.go:43`, `repo.mk:51/54/55`, `Dockerfile:35/37/39`, READMEs, Handbuch `:4`, `:5`, `:82`, `:108`, `:135`, `:2998`. `0.15.0-dev` nirgends mehr. Vollständig gegen `releasing.md`-Liste und Plan.
* **Kennungsfreiheit (§5 Regel 6)**: Diff von README, README.de und Handbuch ohne `LH-`/`ADR-`/`SPEC-`/Slice-Verweise; keine CI-Behauptung; Zeilen `v0.15.0` unangetastet; bestehender CI-Satz im Status-Block unverändert.
* **Datum**: `date +%F` = 2026-10-24 = CHANGELOG, beide Tabellenzeilen, Handbuch-`Stand`.
* **M-T-Protokoll**: plausibel (2b/3b rote Gegenprobe, Grenzen benannt: `is_stable`, echter `curl`/`push`, Token).
* **Sensoren-Darstellung**: `make ci`-Erstlauf (Trivy-DB 404, Mirror), Diagnose-Aufruf ausdrücklich „kein Gate, kein Ersatz“, Wiederholung grün, `test-docker` nicht gelaufen und M2 ausstehend ehrlich ausgewiesen.
* **Plan-Tabelle gegen `git diff --stat`**: zehn Dateien, jede mit Plan-Eintrag; „kein Diff“-Zeilen (`spec/`, `releasing.md`, `.github/`, `scripts/`, `internal/`) stehen nicht im Diff.
* **Kommentare in `repo.mk`/`Dockerfile`/`main.go`**: reine Zahlenänderung, Indikativ, keine Chronik (§3.7/§3.9 a–c), kein Verweis auf umbenannte Symbole.
* **Commits/Lifecycle**: Moves `f246e155`, `f053b1ec` reine Renames (0 Zeilen), Verweis-Nachzüge getrennt; Zwischenstände nicht Spitze, 8 Commits lokal vor `origin/main` (`d98f5755`), kein Push.
* **Tags**: `git tag -l 'v0.15*'` = nur `v0.15.0`, zeigt weiter auf `a895ffd2`; kein `v0.15.1`.
* **Harte Regeln**: kein `//nolint`, kein Gate gelockert, keine Spec-/ADR-Änderung, `architecture.md` unberührt.

## Nicht geprüft

* `make test-docker`, `make verify-depguard` (Gründe oben); `make ci`/`fullbuild` nicht reproduziert; Tap-Trockenlauf nicht nachgefahren (Verifier-Aufgabe laut Plan); CI auf dem Push-Tip (M2, steht aus).
