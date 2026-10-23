# Review: Slice `slice-v1-release-cut-v0.15.0` (Release-Cut `v0.15.0`)

* Rolle: Reviewer (`.harness/skills/reviewer.md`), frischer Kontext
* Gegenstand: `git diff 07353780..HEAD` (HEAD `7f630132`), Welle `welle-forge-cli-gh-glab`
* Maßstab: Plan ([`slice-v1-release-cut-v0.15.0`](../plan/planning/in-progress/slice-v1-release-cut-v0.15.0.md)), Muster des Cuts `v0.14.0`, `docs/maintainer/releasing.md`, `AGENTS.md` §3.7, §3.9, §3.10, §5 Regel 6; nicht die DoD (Verifier)
* Sensor: `make gates` selbst gefahren, Exit 0 (`GATES_EXIT=0`, `ablage-pruefung` 26 Fälle 0 Fehlschläge). `make test-docker` und `make verify-depguard` nicht gefahren: Produktcode im Diff ist nur `var version`; das lokale `test-docker`-Ergebnis des Implementers wird nur auf Darstellung geprüft, nicht reproduziert (Lauf über 30 Minuten, ersetzt ohnehin keinen CI-Lauf).

## Findings

### MEDIUM

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.10 (bedingung in Doku, die eine Eigenschaft des Codes behauptet); Skill-MEDIUM „Doku driftet gegen implementiertes Verhalten“
- `pfad`: `CHANGELOG.md` · „**Limit of `version:`:** `version:` without `source:` ends `u-boot generate devcontainer` with exit code `10`“; ebenso Tabellenzeile `v0.15.0` in `README.md`/`README.de.md` („`version` without `source` ends“)
- `klasse`: bedingung-breiter-als-code
- `befund`: Der Code lehnt `version` ohne `source` nur für Katalogeinträge mit `dockerfileBlock` ab (`validateDevcontainerFeatures`, `internal/hexagon/application/devcontainer_features.go`: `featureCatalogue()[k].dockerfileBlock`; gesetzt für drei Einträge, `gh`, `glab`, `ollama`). Der Lead nennt die Regel ohne Einschränkung und mit Platzhalter `<feature>`; für jedes andere Feature (z. B. `go`) ist `version` ohne `source` zulässig. Der Satz im Eintrag *Added* ist durch den Kontext (gh/glab) enger, der Lead und die Tabellenzeile sind es nicht.
- `risiko`: Ein Leser des Release-Lead hält `version:` bei Features ohne Dockerfile-Block für abgelehnt, oder der nächste Änderer übernimmt die zu breite Bedingung. Genau das Muster aus `BEO-DOK/spec-bedingung-ohne-code-entsprechung`; der Implementer hat den Pfad gelesen (Abweichung 3), die Reichweite der Regel aber nicht auf die drei Einträge eingegrenzt.
- `verifizierbar`: Review-only (Lesung des genannten Symbols); Gegenprobe: `u-boot config set devcontainer.features.go.version 1.22` plus `generate devcontainer` im Projekt, nicht gefahren.

### LOW

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.10 (Zustandsmarker mit Ablaufpunkt)
- `pfad`: `docs/plan/planning/in-progress/roadmap.md` · „(Push und Tag stehen aus)“; `docs/plan/planning/open/README.md` · „in Arbeit (Push und Tag stehen aus)“
- `klasse`: zustandsmarker-mit-ablauf
- `befund`: Zwei Zustandsmarker nennen einen Ablaufpunkt (Push, Tag); der Plan lässt die Closure sie auf den Ist-Stand setzen. Heute stimmen sie (`origin/main` = `07353780`, kein Tag `v0.15*`).
- `risiko`: Bleiben sie nach dem Push stehen, behaupten sie einen falschen Stand. Die Closure-Schritte (Roadmap, `open/README`) müssen beide Stellen anfassen.
- `verifizierbar`: Review-only

### INFO

- Verifier/Hauptagent: M2 (CI auf dem Push-Tip) steht aus und ist im Lead korrekt nicht behauptet. Die unveränderten Status-Block-Sätze „The Docker integration tests run in CI; the four tests of nested Podman are skipped there“ (beide READMEs) stehen unter Version `v0.15.0`; M2 muss Zahl `--- SKIP` und die vier Docker-Tests der Welle gegen diesen Satz prüfen, sonst bleibt er eine Aussage über den Vorgängerlauf.
- Die Aussage des Implementers (Selbstcheck b), `0.14.0` stehe außerhalb von Planung/Reviews/Messprotokoll nur in Tabellenzeilen, CHANGELOG-Abschnitt und `releasing.md`, ist ungenau: `spec/lastenheft.md` Zeilen 1590 und 1593 nennen Lastenheft-Version 0.14.0 (anderer Gegenstand). Kein Handlungsbedarf.
- Verifier: DoD-Abgleich, Tag-/Push-Reihenfolge M1 bis M8, `tap-check` (M4) bleiben dessen Frage.

## Prüfpunkte der Anfrage

1. CHANGELOG `[0.15.0]`: Zusage-Menge (Linux `amd64`, lokale Engine, Docker 29.9.2), „kein authentifizierter Aufruf gemessen“, Nicht-Zusagen und `${localEnv:…}`-Aussage decken sich mit dem Validierungsbericht (Host: Zeile 15; „alle Aufrufe mit Variable … 401 / Bad credentials“). Kalt-Bau-Aussage mit Grenze (Basis-Image lokal, ein Lauf je Profil) ist ehrlich und in Abweichung 2 belegt; kein CI-Satz im Diff. M-2-Satz am Code gelesen: Schreiben gelingt (`revalidateFeatureEntry` prüft nur Nicht-Leere, `config.go`, Fall `ConfigDevcontainerFeatureVersion`); jedes Laden (`readUbootYAMLBody` → `validateDevcontainer` → `validateDevcontainerFeatures`, `config.go`) lehnt ab, daher scheitern `generate`, `doctor` und weitere `config set` (Exit 10 über `ErrConfigSchemaInvalid`); `SetMany` liest über denselben Pfad. Sachlich richtig, Reichweite siehe MEDIUM. Vergleichslinks `[Unreleased]…v0.15.0...HEAD`, `[0.15.0]…v0.14.0...v0.15.0` richtig; „Fifteenth release“ im Stil von „Fourteenth“; Lastenheft 0.22.0 = `spec/lastenheft.md` Zeile 9.
2. Versionsstellen: `grep` über den Baum. `main.go`, `repo.mk` (3 Stellen), `Dockerfile` (3 Stellen), beide READMEs (Status, `VERSION=`), Handbuch (Kopf, 3 Beispiele, Gültigkeitsbereich) stehen auf `0.15.0`; verbliebene `0.14.0` nur Tabellenzeilen, CHANGELOG-Abschnitt/Link, `releasing.md` (richtig bis zur Closure), Lastenheft-Historie, Messprotokoll, Planung. Nichts Vergessenes, nichts Falsches ersetzt. README, README.de und Handbuch: im `+`-Diff keine `LH-/ADR-/SPEC-/ARC-/slice-/welle-`-Kennung und keine neue CI-Behauptung.
3. Datum: `date +%F` = `2026-10-23`, identisch in CHANGELOG, beiden Tabellenzeilen, Handbuch-`Stand`.
4. M-T-Protokoll: plausibel (2b/3b rote Gegenproben, Grenze benannt: Stellvertreter für `curl`/`git`, `is_stable`, Token nicht geprobt); nicht reproduziert.
5. `test-docker`: Exit 2 am Paket-Zeitlimit, 1 FAIL (`NestedDocker_Restart`, Netz `download.docker.com`), Einzellauf PASS, vier Welle-Tests PASS, „ersetzt den CI-Lauf nicht“, M2 als ausstehend: ehrlich dargestellt; „Netzstörung, keine Codeänderung“ ist durch Einzellauf PASS belegt, nicht durch Wortsuche.
6. Plan-Tabelle vs. `git diff --stat`: 11 Dateien, jede hat einen Plan-Eintrag; `releasing.md`, `spec/`, ADR, Workflows kein Diff, wie geplant. Kommentare in `repo.mk`/`Dockerfile` nennen den geltenden Wert (§3.7/§3.9 a bis c ohne Befund).
7. Commits: `9f15adf5` ist reiner Move (0 Zeilen), `138cdc41` Verweis-Nachzug als eigener Commit, danach Inhalt; die Spitze `7f630132` ist inhaltlich, nicht der Zwischenzustand. `main` ist 6 vor `origin/main` (vor diesem Report); alle gehen als ein Push.
8. Keine Tag-/Push-Aktion: `git tag -l 'v0.15*'` leer, `origin/main` = `07353780`.

## Negativbefund

- Geprüft, ohne Befund: `Dockerfile`, `repo.mk`, `cmd/uboot/main.go` (Versionsstellen, Kommentare); README/README.de (Kennungsfreiheit, keine CI-Neubehauptung); Handbuch (nur Versionsstellen, `:2060`-Abhilfe deckt sich mit Lead); Welle-Datei (Zustandsmarker abgelöst); `open/README` Zeile der Nachbesserung auf Ist-Stand; Lifecycle (Move rein, Datei nur in `in-progress/`).
- Nicht geprüft: `make test-docker`, `make ci`/`image-scan`, `make verify-depguard` (keine Importregel berührt), `tap-check` und Reproduktion von M-T (Hauptagent, M4); Wortlaute der Messprotokolle jenseits der zitierten Validatorzeilen.
