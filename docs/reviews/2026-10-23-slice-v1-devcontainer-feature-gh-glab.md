# Review-Report: slice-v1-devcontainer-feature-gh-glab — 2026-10-23

**Rolle:** Reviewer (Modul 10, [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md), [`harness/review.md`](../../harness/review.md)), frischer Kontext. Prüfung gegen Plan, [`ADR-0035`](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md), Lastenheft 0.22.0 und Hard Rules, nicht gegen die DoD (Verifier).

**Gegenstand:** Commit-Range `db989760..HEAD` (`a87749da`, `b7102f19`, `b930721a`, `19139e31`), Slice [`slice-v1-devcontainer-feature-gh-glab`](../plan/planning/in-progress/slice-v1-devcontainer-feature-gh-glab.md), Welle [`welle-forge-cli-gh-glab`](../plan/planning/welle-forge-cli-gh-glab.md), [`ADR-0035`](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md), [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer).

## Ergebnis

**Keine HIGH-, keine MEDIUM-Findings.** Zwei LOW, vier INFO. Der Dockerfile-Block ist gegen Messprotokoll und Handvarianten korrekt, die Zugangsdaten-Zusage hält, die Tests pinnen.

## Findings

### F-1 — LOW — Kommentar in `export_test.go` sagt „sandbox Dockerfile“, der Plan führt die Datei als „unverändert (gelesen)“

- `kategorie`: LOW · `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.9 (b), §3.10 · `klasse`: kommentar-alter-name-nach-diff
- `pfad`: `internal/hexagon/application/export_test.go` · „whose installation the sandbox Dockerfile carries.“
- `befund`: Das Feld `DockerfileBlock` der Test-Projektion gilt seit dem Diff für `gh` und `glab` in beiden Profilen; der gleichlautende Kommentar in `devcontainer_features.go` wurde berichtigt, dieser nicht. Die Plan-Tabelle führt die Datei als „unverändert“ mit „projizieren allgemein … gelesen“ — der Kommentar wurde nicht mitgelesen.
- `risiko`: Der nächste Änderer eines Dockerfile-Block-Features liest die Test-Projektion als Beleg, dass der Block nur im Sandbox-Profil wirkt.
- `verifizierbar`: Review-only (Suche nach „sandbox Dockerfile“ im Baum).

### F-2 — LOW — Handbuch sagt „keine Signatur gefunden“ ohne die Grenze der Messung

- `kategorie`: LOW · `quelle`: [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Wortlaut „nach dem Suchmuster, nicht abschließend“), [`AGENTS.md`](../../AGENTS.md) §3.10 (Aussage über Verhalten nennt nur Gemessenes) · `klasse`: zusage-weiter-als-messung
- `pfad`: `docs/user/benutzerhandbuch.md` · „keine Signatur gefunden; `u-boot generate devcontainer` weist das aus“ (Zeilen `gh` und `glab` der Katalogtabelle)
- `befund`: Lastenheft, CHANGELOG, Warnungstext und Block-Kommentar tragen den Zusatz „nach Suchmuster, nicht abschließend“ (Messung A12: Suchmuster, nicht abschließend); die beiden Handbuchzeilen lassen ihn weg.
- `risiko`: Gering. Ein Nutzer liest „keine Signatur gefunden“ als „es gibt keine“. Der Spec-Slice schreibt den Abschnitt mit den Grenzen ohnehin nach, dort lässt sich das mitziehen.
- `verifizierbar`: Review-only.

### F-3 — INFO — Mutationsprotokoll: Mutationen per Editor gesetzt, nicht per `git apply`

- `kategorie`: INFO · `quelle`: Plan §3 Mutationsprotokoll · `klasse`: protokoll-abweichung-offen-benannt
- `befund`: §7 nennt die Abweichung selbst („mit dem Editor gesetzt und als Patch gesichert“). Die Rücknahme per `git apply -R` und der saubere Baum sind belegt. Keine Wirkung auf die Aussagekraft.
- `verifizierbar`: Review-only.

### F-4 — INFO — Eigene Nachfahrt der Mutationen: sieben Läufe, sechs rot, eine ohne Wirkung

- `kategorie`: INFO · `quelle`: Maintainability · `klasse`: mutation-nachgefahren
- Nachgefahren mit `docker run … u-boot:test-docker-tools go test -count=1 ./internal/hexagon/application/ ./cmd/uboot/`, Rücknahme per `git apply -R`, Baum danach sauber:
  - B1 (a) `or .UID .Sandbox .Gh .Glab` → `or .UID .Sandbox`: rot (`…BlocksInBothProfiles`, `…UserRootOnlyWithReason`, `…ActivateAndDeactivateLater`, `TestRun_Generate_ForgeCLINotices`).
  - B6 (a) Zeile `glab` aus `egressFeatureHosts()`: rot (`…EgressHosts`, `…TogetherWithOllamaAndRuntime`, `TestRun_EgressList_ForgeCLIHostOrigins`).
  - X3 `&& f.Source == ""` entfernt: rot (`…SourceOverrideIsExternalFeature`, `…RemoteEnvReferences`, `TestGenerateDevcontainer_Ollama_NotActive`, `TestDoctor_SandboxCredentials_ForgeCLI`).
  - Zusätzlich, nicht im Plan: letzte Ziffer von `GLAB_SHA256_ARM64` geändert: rot (`…BlocksInBothProfiles`, `…TogetherWithOllamaAndRuntime`); `--no-same-owner` aus dem `glab`-`tar` entfernt: rot (dieselben zwei); `GITLAB_TOKEN` → `GITLAB_TOKEN2` in `agentCredentialEnv`: rot (`…RemoteEnvReferences`, beide `ForgeCLI`-Doctor-Tests).
  - Ohne Wirkung: `&& cd / \` → `&& cd /tmp \` (alle vier Blöcke der Datei, auch die fremden): beide Pakete grün. Der Rücksprung ins Wurzelverzeichnis vor `rm -rf "$work"` ist in keinem Golden Case gepinnt. Funktional folgenlos (`rm -rf` des aktuellen Verzeichnisses gelingt), deshalb INFO.
- `verifizierbar`: nachgefahren (siehe oben).

### F-5 — INFO — Einmaliges Host-python zum Schreiben von Quelltext

- `kategorie`: INFO · `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.1 · `klasse`: host-werkzeug-schreibt-quelltext
- `befund`: Bewertet an `templates.go`, `devcontainer_sandbox.go`, `devcontainer_features.go` und den übrigen Quelltext-Diffs: nur ASCII, kein CR, kein Whitespace am Zeilenende, im Dockerfile-Template genau eine entfernte Zeile (die `if`-Bedingung), also byte-identisch für bestehende Konfigurationen; `gofmt` und `make lint` grün. §3.1 verbietet die Host-Go-Toolchain als Nachweis, nicht das Schreiben; der Nachweis lief über `make gates`. Der Lauf ist in den Artefakten nicht belegt (nur gemeldet); ich konnte nur das Ergebnis prüfen, nicht den Vorgang. Keine Auffälligkeit.
- `verifizierbar`: `make gates`.

### F-6 — INFO — Grenzen der Docker-Tests (Zuständigkeit Verifier/Validator)

- `kategorie`: INFO · `quelle`: Maintainability · `klasse`: docker-test-netz-und-arch
- `befund`: `TestE2E_ForgeCLI_EgressHosts` ruft `gitlab.com` live an (Kopfkommentar benennt, dass ein Netzfehler rot färbt). Der `arm64`-Pin ist nicht gebaut ([`ADR-0035`](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) Konsequenzen nennen es; das Messprotokoll nennt das Nachlesen gegen die Prüfsummen-Dateien); die Negativprobe überschreibt beide Architektur-Werte und gilt damit auf beiden Hosts. Hermetik der Container-Prüfung: Auswahl nach Nadel (Nullprüfsumme im Befehl) statt nach Präfix, keine `t.Parallel`; mein Lauf von `PinMismatchBreaksBuild`, `TokenReferences…` und `StreamContains` (31,3 s) ließ die Container-Zahl vorher und nachher bei 4.
- `verifizierbar`: `make test-docker` gezielt.

## Negativbefunde (geprüft, ohne Befund)

- `Dockerfile.tmpl`: Pins `GH_VERSION=2.103.0`, `GLAB_VERSION=1.124.0` und alle vier SHA-256 stimmen zeichengleich mit `tools/messung/forge-cli/block-ghtgz.dockerfile`, `block-glabtgz.dockerfile` und `rohdaten-03.txt`; Block-Körper gleich den Handvarianten; Reihenfolge Download → `sha256sum -c -` → `tar` → Aufräumen; `--strip-components=1` mit Member-Pfad nur bei `gh`; kein `ENV`/`auth login`. Bedingung `or .UID .Sandbox .Gh .Glab`; die Blöcke liegen innerhalb des `USER root`-Zweigs vor `USER vscode`; Binaries gehören `root`.
- Zugangsdaten ([`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`ADR-0035`](../plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) Entscheidung 5): `AgentEnv` wird in `devcontainerTemplateData` erst nach dem frühen `return` des Standardprofils gesetzt, also nur im Sandbox-Profil; `remoteEnv` rendert `${localEnv:…}`; `agentCredentialEnv` über `dockerfileBlockFeatureActive`; Reihenfolge fest; Doctor-Abhilfe und Hinweis nennen die neuen Variablen; kein Tokenwert in Template, Block, Test (nur Kanarienwerte).
- Ausweisungstexte (`forgeCLIWarnings`) gegen Messung A10/A12 und Lastenheft 0.22.0: Text von `gh` nennt Attestierung (existiert, Prüflogik nicht gemessen) und SHA-256-Grenze, Text von `glab` nennt Suchmuster und keine Attestierung; keine `SecurityRelaxation`/`DegradationFallback`.
- `source:`-Überschreibung (X3) und Hosts am Namen: `egressFeatureHosts` hat `gh` und `glab`; `github.com`/`api.github.com` stehen in der Basisliste (gelesen).
- Architektur: keine neuen Importe in `internal/hexagon/application`, kein neuer Port, kein neuer Sentinel; Dual-Classifier nicht berührt; keine `//nolint`; `spec/`, Accepted-ADR und `.golangci.yml` nicht im Diff.
- Doku: `docs/user` und READMEs kennungsfrei ([`AGENTS.md`](../../AGENTS.md) §5 Regel 6); README-Zahl 16 stimmt mit `featureCatalogue()` (16 Einträge); CHANGELOG, `harness/replay.md`, Roadmap, `open/README.md`, Welle-Datei und Lifecycle-Moves (reine `git mv` mit Verweis-Nachzug, `a87749da`/`b7102f19`) konsistent.
- Plan-Tabelle gegen `git diff --stat`: alle 23 Dateien haben einen Eintrag, jeder Eintrag „update/neu“ einen Diff; Ausnahme siehe F-1 („unverändert“ ohne gelesenen Kommentar).
- Sensors: `make gates` selbst gefahren, Exit `0` (Coverage 92,00 % gegen 90 %); `make test-docker` nur gezielt (siehe F-6).

## Nicht geprüft

- `make verify-depguard`: nicht gelaufen, weil keine Importregel berührt ist.
- Vollständiger `make test-docker`: nicht gelaufen (Paket-Zeitlimit bekannt); `BuildsInBothProfiles` und `EgressHosts` nicht selbst gefahren, nur die Implementer-Läufe in §7 gelesen. Das Handbuch- und Spec-Thema `GH_TOKEN`/`GITLAB_TOKEN` (B10) liegt beim Spec-Slice.
