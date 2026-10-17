# Review-Report: slice-v2-homebrew-smoke-haertung — 2026-10-17

**Review-Art:** Code (gegen Plan + Konventionen)

**Gegenstand:** Commit-Range `9ec0f12..HEAD` (`8b5bd3b`, `ddf2813`, `82677f1`, `88a7259`) zum Slice `slice-v2-homebrew-smoke-haertung` <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

**Skill:** `.harness/skills/reviewer.md` @ `88a7259` · <!-- d-check:ignore (Pfad als Zitat) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-17

**Eingangs-Kontext:**

- Slice `slice-v2-homebrew-smoke-haertung` (in-progress; Scope F-4, F-5, F-8, Akzeptanzkriterium 4 aus `slice-v2-homebrew-formula`) <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
- [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) (Accepted, im Diff unverändert)
- [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung), [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache)
- `AGENTS.md` (Hard Rules, insbesondere §3.7 und §3.9)

Sensor-Lage: `make test` wurde nur für Mutationen gefahren (alle zurückgesetzt, `git status` sauber). Kein Workflow-Linter im Repo; `make gates` und `make docs-check` liegen beim Implementer bzw. Verifier.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der Job `tap` hängt jetzt an einem Cross-Job-Output (`publish` → `steps.tag.outputs.is_stable`). Die Verdrahtung ist nach Lesen des ganzen Workflows korrekt (Step-Id `tag`, `$GITHUB_OUTPUT`, `needs: publish`, Job-`outputs:`), aber kein Sensor prüft sie. Szenario: wird Output oder Step-Id später umbenannt, ist der Output leer, `tap` wird für stabile Tags lautlos übersprungen und die Formel erreicht das Tap nicht. Ein Pre-Release-Tag deckt das nicht auf, weil `tap` dort in beiden Zweigen übersprungen wird. | [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung), [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md), Maintainability | `.github/workflows/publish.yml` · `if: ${{ needs.publish.outputs.is_stable == 'true' }}` | nein — kein `actionlint` o. ä. im Repo; erst der nächste stabile Tag belegt es | Workflow-Kopplung ohne Sensor, stilles Überspringen |
| F-2 | MEDIUM | Das Debian-/RPM-Paket trägt weiter die deutsche Beschreibung, die F-5 an der Formel behebt. Der neue Test deckt nur `scripts/homebrew-formula.rb.tmpl`. Szenario: `apt show u-boot` und `dnf info u-boot` zeigen deutschen Text, obwohl [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache) Englisch für erzeugte Dateien verlangt. Liegt außerhalb des Diffs; die Suche nach dem alten Wortlaut (§3.9 b) trifft diese Stelle. | [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache), `AGENTS.md` §3.9 | `packaging/nfpm.yaml` · `Bootstrapper fuer reproduzierbare Docker-/Devcontainer-Stacks.` | ja — Test analog `TestHomebrewFormulaDesc_EnglishASCII` wäre rot; sonst Review-only | Sprachvertrag: Distributions-Metadaten deutsch |
| F-3 | LOW | Der im Diff umgeschriebene Kommentar zu `SHA256SUMS` und Formel steht über dem Step „Build and upload Linux packages“, beschreibt aber den zwei Steps später folgenden Step „Upload SHA256SUMS and Homebrew formula“; dazwischen steht der `distro-pakete`-Kommentar. Ein Änderer des Packages-Steps liest eine Aussage über `SHA256SUMS`, die dort nicht gilt. | `AGENTS.md` §3.7 (Kopplung) | `.github/workflows/publish.yml` · `one SHA256SUMS over the six binaries` | Review-only | Kommentar-Block vom beschriebenen Step getrennt |
| F-4 | INFO | Die DoD nennt als Beleg einen grünen Smoke-Lauf im Tap. Die geänderte Vorlage `smoke.yml` ist in diesem Diff nicht ausgerollt (Plan-Zeile „Tap-Repo … ausgerollte Fassung“); `brew audit --strict --online` lief bisher nie gatend, ein echter Formel-Befund ist möglich (Slice §6). Aus der Formel allein kenne ich keinen konkreten Regelverstoß (`desc` beginnt groß, ohne Artikel, ohne Punkt, kurz; `license`, `homepage`, `version`, `sha256`, `test do` vorhanden, Klasse `UBoot`). Ob `version` neben der URL als redundant gilt, ist hier nicht entscheidbar. | [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung) | `packaging/homebrew-tap/.github/workflows/smoke.yml` · `brew audit --strict --online pt9912/u-boot/u-boot` | nein — Zuständigkeit Verifier (Tap-Lauf) | Beleg steht aus (Tap-Lauf) |
| F-5 | INFO | `brew install --build-from-source` bei einer Formel mit vorgebauten Binaries ist dieselbe Installation wie ohne Flag; der Kommentar im Workflow sagt das ehrlich („laeuft auch ohne Flag“). Ob das Akzeptanzkriterium 4 („bzw. gleichwertig“) damit erfüllt ist, entscheidet der Verifier, nicht der Reviewer. | [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung) | `packaging/homebrew-tap/.github/workflows/smoke.yml` · `brew install --build-from-source u-boot` | nein — Verifier | Kriterium-Auslegung |
| F-6 | INFO | `TestHomebrewFormulaDesc_EnglishASCII` ist nicht tautologisch: Mutation auf deutsche `desc` und auf `desc` mit Gedankenstrich (U+2013) färbten `make test` rot, mit der erwarteten Meldung. Die deutschen Fragmente sind eine Denylist (`fuer`, ` und `, `reproduzierbar`); ein deutsches Wort außerhalb davon bleibt unentdeckt, solange `reproducible` im Text steht (aus dem Code gelesen, nicht gefahren). | [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache) | `cmd/uboot/homebrew_scripts_test.go` · `[]string{"fuer", " und ", "reproduzierbar"}` | ja — `make test` | Test-Schärfe, Denylist |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `publish.yml` Syntax: `outputs:` am Job `publish`, `steps.tag.outputs.is_stable`, `needs: publish` am Job `tap` | geprüft, ohne Befund |
| `publish.yml` Validate-Step: Regex nimmt `vX.Y.Z` und `vX.Y.Z-pre` an, `is_stable` ist `false` genau bei `-` in der Version; `v0.9.0` und `v1.0.0` ergeben `true`, `v1.0.0-rc.1` ergibt `false`; Verhalten gleich wie `!contains(github.ref_name,'-')` (Build-Metadaten `+` sind abgewiesen) | geprüft, ohne Befund |
| `publish.yml` Kommentar „SHA256SUMS … six binaries and the four Linux packages … four brew digests“ gegen Workflow und `scripts/homebrew-formula-fill.sh` | geprüft, ohne Befund (Inhalt stimmt: eine Datei, vier Digests gelesen) |
| `smoke.yml` Struktur, Matrix (`macos-14`, `macos-15-intel`), `permissions: contents: read`, `HOMEBREW_GITHUB_API_TOKEN`, Audit ohne `\|\| true` | geprüft, ohne Befund |
| `smoke.yml` `grep -F "BEGIN U-BOOT MANAGED BLOCK: service.postgres"`: Marker stimmt mit `internal/adapter/driven/yaml/mapping_patch.go` (`# BEGIN U-BOOT MANAGED BLOCK: `) und Präfix `service.` in `addservice.go` überein; `add` braucht kein Docker | geprüft, ohne Befund |
| `smoke.yml` Kommentare (Indikativ, keine Chronik, Kennungen als Herkunft) | geprüft, ohne Befund |
| `scripts/homebrew-formula.rb.tmpl` `desc` Englisch und ASCII; kein alter Wortlaut im Baum außer F-2 | geprüft, ohne Befund |
| `docs/maintainer/releasing.md` §6 gegen Workflow und Smoke | geprüft, ohne Befund |
| [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) unverändert, keine Abwärts-Referenz aus Rang 1–3, keine Gate-Lockerung (Audit wird strenger) | geprüft, ohne Befund |

Nicht geprüft:

- `brew audit --strict` gegen die echte Formel und Smoke-Läufe im Tap — Docker-only, Host-`brew` gesperrt; Zuständigkeit Verifier (F-4).
- `make gates`, `make docs-check` über den Diff — Sensor des Implementers; hier nur `make docs-check` für diesen Report.
- `make test-docker`, `make verify-depguard` — kein Docker-Pfad und keine Importregel berührt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 1 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Workflow-Kopplung ohne Sensor, stilles Überspringen · Sprachvertrag: Distributions-Metadaten deutsch · Kommentar-Block vom beschriebenen Step getrennt · Beleg steht aus (Tap-Lauf) · Kriterium-Auslegung · Test-Schärfe, Denylist

## Verdikt

**Merge-blockierend:** nein für den Diff als solchen; F-1 und F-2 sind vor dem Slice-Abschluss zu klären (F-2 ist Scope-Entscheidung des Implementers/Architects, da außerhalb des Diffs).

**Übergabe:** Findings an den Implementer; Finding-Klassen in die Slice-Closure §7. DoD-Konformität (Smoke-Lauf grün, Akzeptanzkriterium 4) prüft der Verifier.
