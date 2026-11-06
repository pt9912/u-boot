# Review-Report: slice-v1-devcontainer-brew-ripgrep-im-repo — 2026-11-06

**Review-Art:** Code-Review gegen Plan und Konventionen: Diff gegen Slice-Plan, Lastenheft, Spezifikation, ADR und Hard Rules. Der Diff enthält Konfiguration, erzeugte Dateien und Maintainer-Doku, keinen Go-Code.

**Gegenstand:** Commit `9d255cbb` (HEAD, Basis `7fbad506`), vier Dateien, 39 Einfügungen, 5 Löschungen: `u-boot.yaml`, `.devcontainer/Dockerfile`, `.devcontainer/devcontainer.json`, `docs/maintainer/devcontainer-arbeiten.md`. Zum Abgleich der Plan-Tabelle zusätzlich die Spanne ab `80beb49b~1` (sieben Dateien, Slice-Datei als reiner Move).

**Skill:** `.harness/skills/reviewer.md`. **Modell:** claude-sonnet-5-5. **Datum:** 2026-11-06.

## Eingangs-Kontext und Sensors

**Gelesen:**

- Slice [`slice-v1-devcontainer-brew-ripgrep-im-repo`](../plan/planning/in-progress/slice-v1-devcontainer-brew-ripgrep-im-repo.md): §1 (Ziel, Entscheidungen E1 bis E3, gefundener Stand, Nicht-Ziele), §2 (DoD), §3 (Plan-Tabelle), §4 und §5 (Trigger).
- Lastenheft: [`LH-FA-DEV-013`](../../spec/lastenheft.md#lh-fa-dev-013--homebrew-formeln-im-devcontainer), [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) (Pin-Inventur, im Plan als nicht berührt geführt).
- Spezifikation: [`SPEC-057`](../../spec/spezifikation.md#spec-057--homebrew-formeln-im-devcontainer-katalog-schlüssel-brew-pfade-freigabe-der-taps-konflikt-mit-homebrew-pin-meldungen-hosts-und-grenzen) (Dockerfile-Block, Pin, Meldungen, Gemessen, Nicht zugesagt), [`SPEC-031`](../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) (Bau-Pfad).
- [`ADR-0046`](../plan/adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md) (Kontext, gelesene und nicht gelesene Punkte), [`ADR-0043`](../plan/adr/0043-nur-devcontainer-bau-pfad.md) (die CI baut das Image), Messprotokoll [`slice-v1-devcontainer-homebrew-pakete`](../messungen/slice-v1-devcontainer-homebrew-pakete/README.md) (M1, M2).
- `AGENTS.md` §3.1, §3.2, §3.7, §3.9, §3.10, §5; `harness/review.md`; `harness/conventions.md`.

**Sensors:**

- Eigene Reproduktion der Erzeugung, nur lesend: Git-Kopie von HEAD unter dem Scratchpad, `/tmp/u-boot generate devcontainer` (Binary vom 2026-11-06 07:19, vor dem Commit aus dem Baum gebaut). Ergebnis: `Dockerfile` byte-gleich zu HEAD, `devcontainer.json` unterscheidet sich nur in der Zeile `postCreateCommand` (die Kopie hat einen lokalen Pfad als `origin`; im Klon mit `origin` ist sie gleich, wie der Plan es ausweist). `u-boot doctor` in der Kopie: `devcontainer.features.allowlist` und `devcontainer.features.drift` je `ok`.
- Pin: `git ls-remote https://github.com/Homebrew/brew refs/tags/7.0.12` liefert `bf464f57275cae0be106b43d190c98564ac59b0e`, gleich dem `ARG HOMEBREW_REVISION` im Dockerfile und dem Kommentar „release 7.0.12“.
- Suche nach Lesern von `homebrew`/`brew` in `scripts/`, `Makefile`, `repo.mk`, `tools/mk/`, `tools/harness/`, `.github/workflows/`: kein Skript, kein Target und kein Workflow ruft `brew` des Entwicklungs-Containers auf. `.devcontainer/devcontainer-lock.json` bleibt im Diff unberührt und trägt weiter den Eintrag `homebrew:2.0.5` (im Plan als Zwischenstand benannt).
- `make docs-check`: nach Anlage dieses Reports gefahren, Ergebnis siehe Abschnitt *Sensor-Lauf* am Ende.
- `make gates`, `make test-docker`, `make verify-depguard`: nicht gefahren. Der Diff enthält keinen Go-Code und keine Importregel; `make gates` im Devcontainer-Container ist Closure-Pflicht des Implementers und Nachweis des Verifiers. `make test-docker` baut `.devcontainer/` nicht (Plan-Tabelle, Zeile `integration.yml`). Ein Bau des Images (`docker build`) lief nicht: er ist der CI-Lauf auf dem Push-Tip und der Nachweis des Projektinhabers (E2), nicht dieser Lauf.

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Die Plan-Tabelle in §3 führt `docs/plan/planning/open/README.md`, nennt aber `docs/plan/planning/in-progress/roadmap.md` nicht. Die Slice-Spanne berührt die Roadmap (Commit `7fbad506`: der Absatz „Nichts in Arbeit.“ entfällt, zwei Zeilen gelöscht). | `AGENTS.md` §3.10 (Plan-Tabelle listet jede Datei, die der Diff berührt; Abgleich gegen `git diff --stat`) | [`slice-v1-devcontainer-brew-ripgrep-im-repo`](../plan/planning/in-progress/slice-v1-devcontainer-brew-ripgrep-im-repo.md) §3 · „`docs/plan/planning/open/README.md` \| update“ (Plan-Tabelle) | nein — Review-only (Abgleich der Tabelle gegen `git diff --stat 80beb49b~1 HEAD`) | Berührte Datei fehlt in der Plan-Tabelle |
| F-2 | INFO | Der Satz in der Maintainer-Doku sagt, das Dockerfile installiere `ripgrep` beim Bau; das ist am Dockerfile gelesen und von der Erzeugung bestätigt, ein Bau in der CI liegt zum Zeitpunkt des Commits nicht vor. Der Plan (DoD Liefer-Punkt 2) bindet den Satz an „CI-Bau grün, Run-ID in der Closure“. | `AGENTS.md` §3.10; Slice DoD Liefer-Punkt 2 | `docs/maintainer/devcontainer-arbeiten.md` · „installiert das Dockerfile beim Bau des Images `ripgrep` als Homebrew-Formel“ | ja — CI-Lauf auf dem Push-Tip (Schritt *Build the devcontainer image*); Zuständigkeit Verifier | Doku-Aussage über Verhalten vor ihrem Beleg |
| F-3 | INFO | Die Kombination des Blocks `brew` mit `ollama`, `gh`, den Werkzeugketten und `nestedRuntime` in einem Image ist nach [`SPEC-057`](../../spec/spezifikation.md#spec-057--homebrew-formeln-im-devcontainer-katalog-schlüssel-brew-pfade-freigabe-der-taps-konflikt-mit-homebrew-pin-meldungen-hosts-und-grenzen) „nicht gebaut“, das Messprotokoll baut `brew` mit Basis-Image und Formeln; dieser Commit ist der erste Bau der Kombination, und `arm64` ist nicht gemessen. Der Plan benennt das (E4, E5, Liefer-Punkt 3). | [`SPEC-057`](../../spec/spezifikation.md#spec-057--homebrew-formeln-im-devcontainer-katalog-schlüssel-brew-pfade-freigabe-der-taps-konflikt-mit-homebrew-pin-meldungen-hosts-und-grenzen) Absatz *Nicht zugesagt* | `.devcontainer/Dockerfile` · „RUN HOMEBREW_NO_AUTO_UPDATE=1 HOMEBREW_NO_ANALYTICS=1 brew install --formula 'ripgrep'“ | ja — CI-Lauf auf dem Push-Tip, `rg`-Nachweis des Projektinhabers; Zuständigkeit Verifier | Erster Bau einer ungemessenen Kombination |
| F-4 | INFO | `ENV PATH` stellt `/home/linuxbrew/.linuxbrew/bin` vor den bisherigen Wert; die Abhängigkeiten von `ripgrep` (laut Messprotokoll `bzip2`, `pcre2`, `zlib-ng-compat`) können Programme im selben Verzeichnis verlinken, die gleichnamige Systemprogramme überdecken. `tool-pruefung` prüft Version oder Anwesenheit von `go`, `golangci-lint`, `make`, `bash`, `curl`, `gh`, `git`, `docker`, nicht den Ort. Der Plan deckt das mit dem `ls`-Befehl des Projektinhabers (Liefer-Punkt 3). | Slice DoD Liefer-Punkt 3 | `.devcontainer/Dockerfile` · „ENV PATH="/home/linuxbrew/.linuxbrew/bin:/home/linuxbrew/.linuxbrew/sbin:${PATH}"“ | ja — `docker run … ls /home/linuxbrew/.linuxbrew/bin` des Projektinhabers; Zuständigkeit Verifier | PATH-Vorrang eines Paketmanager-Verzeichnisses |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `u-boot.yaml` | geprüft, ohne Befund: nur der Austausch `homebrew` (mit `version: 2.0.5`) gegen `brew` mit `packages: [ripgrep]`; `enabled: true` gesetzt, kein Rest von `homebrew`, keine Tap-Freigabe nötig (Formel ohne Tap, [`SPEC-057`](../../spec/spezifikation.md#spec-057--homebrew-formeln-im-devcontainer-katalog-schlüssel-brew-pfade-freigabe-der-taps-konflikt-mit-homebrew-pin-meldungen-hosts-und-grenzen) *Freigabe der Taps*); übrige Schlüssel unverändert |
| `.devcontainer/Dockerfile` | geprüft, ohne Befund: 32 Einfügungen, keine Löschung, wie in der Plan-Vorschau; Erzeugung reproduziert byte-gleich; der Block folgt den fünf Schritten von [`SPEC-057`](../../spec/spezifikation.md#spec-057--homebrew-formeln-im-devcontainer-katalog-schlüssel-brew-pfade-freigabe-der-taps-konflikt-mit-homebrew-pin-meldungen-hosts-und-grenzen) *Dockerfile-Block* (Build-Abhängigkeiten, `ENV PATH`, Checkout mit Revisionsprüfung als `vscode`, `brew install --formula`, `USER root`); die Umgebungsvariablen stehen nur am `RUN` des Installationsschritts, kein `ENV HOMEBREW_*`; die Zeilen `ARG GO_VERSION=` und `ARG GOLANGCI_LINT_VERSION=` (Quellen der `pinsync`-Gruppen) sind unverändert; Pin gegen den Release-Tag bestätigt |
| `.devcontainer/devcontainer.json` | geprüft, ohne Befund: Eintrag `homebrew:2.0.5` entfällt, Komma am `node`-Eintrag entfernt, JSON bleibt gültig; `brew` erzeugt keinen Schlüssel unter `features`; `doctor` meldet Allowlist und Drift mit `ok` |
| `.devcontainer/devcontainer-lock.json` | geprüft, ohne Befund: im Diff unberührt (E1); der verwaiste Eintrag ist als Zwischenstand im Plan benannt, kein Gate liest die Datei |
| `docs/maintainer/devcontainer-arbeiten.md` | geprüft, ohne Befund gegen die Hard Rules: „die Formel ist nicht gepinnt“ stimmt mit dem Dockerfile-Kommentar und [`SPEC-057`](../../spec/spezifikation.md#spec-057--homebrew-formeln-im-devcontainer-katalog-schlüssel-brew-pfade-freigabe-der-taps-konflikt-mit-homebrew-pin-meldungen-hosts-und-grenzen) überein; keine Aussage über Bauzeit oder Image-Größe; Maintainer-Doku, daher kein Verstoß gegen die Kennungsregel der Nutzer-Doku (`AGENTS.md` §5 Regel 6). Vorbehalt: F-2 |
| Docker-only und Gates (`AGENTS.md` §3.1, §3.6) | geprüft, ohne Befund: keine Änderung an `ci.yml`, `integration.yml`, `repo.mk`, `Makefile`, `scripts/bau.sh`, `.d-check.yml`, `.golangci.yml`; keine Lockerung, kein `//nolint`, keine Schwelle berührt |
| Kommentare im Diff (`AGENTS.md` §3.7, §3.9) | geprüft, ohne Befund: die Kommentare im Dockerfile-Block stehen im Indikativ und tragen Zusage, Grenze und Kopplung (Pin und Prüfung, Umfang des Pins, `PATH`-Datei); keine Verweise auf Symbole, die der Diff umbenennt; die Kommentare stammen aus dem unveränderten Template; `u-boot.yaml` trägt keinen Kommentar |
| Commit-Nachricht (Traceability) | geprüft, ohne Befund: Slice-Kennung, [`LH-FA-DEV-013`](../../spec/lastenheft.md#lh-fa-dev-013--homebrew-formeln-im-devcontainer) und [`ADR-0046`](../plan/adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md) genannt; Attribution-Zeile vorhanden |
| Plan-Abgleich (`AGENTS.md` §3.10) | geprüft, bis auf F-1 ohne Befund: `git diff --stat` des Commits stimmt mit den erwarteten Zahlen (Dockerfile +32, `devcontainer.json` +1/−2, `u-boot.yaml` +3/−2) und den Tabellenzeilen überein; die als „unverändert“ geführten Dateien sind im Diff nicht enthalten |
| Referenzmodell, ADR-Immutabilität, Lifecycle (`git mv`) | geprüft, ohne Befund: keine ADR, Spec oder Architektur-Sicht berührt; der Slice liegt nur in `in-progress/`, die Moves sind reine Move-Commits (`adcda41a`, `9ff0613f`) |
| Sprachvertrag, Exit-Codes, Dual-Classifier, Dateisicherheit, Architekturgrenzen | geprüft, ohne Befund: kein Go-Code im Diff, die erzeugten Dateien sind Englisch ([`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache)) |
| Konflikt mit [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) | geprüft, ohne Befund: der Block führt kein Installer-Skript aus; die Warnung der Erzeugung nennt Quelle, Pin und die Ungepinntheit der Formel |

**Nicht geprüft:**

- Bau des Images und `rg` im `PATH` eines neu gebauten Containers — ein Bau braucht das Netz des CI oder den Docker-Host des Projektinhabers; DoD Liefer-Punkt 3 und Bedingung E2, Sache des Verifiers.
- Bauzeit, Dauer und Image-Größe im CI (E3) — erst nach dem Lauf auf dem Push-Tip messbar.
- Neuschreibung der Lock-Datei beim Neubau — liegt nach dem Slice-Commit und beim Projektinhaber (Lock-Schritt).
- Die DoD-Häkchen und die Closure-Pflichten des Slice — Verifier-Aufgabe, nicht Gegenstand des Reviews.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Berührte Datei fehlt in der Plan-Tabelle · Doku-Aussage über Verhalten vor ihrem Beleg · Erster Bau einer ungemessenen Kombination · PATH-Vorrang eines Paketmanager-Verzeichnisses

## Verdikt

**Merge-blockierend:** nein. Kein HIGH und kein MEDIUM. Der Diff ist die reproduzierbare Erzeugung aus `u-boot.yaml`, der Pin stimmt mit dem Release-Tag überein, und kein Gate, Skript oder Workflow wird berührt. F-1 ist eine Berichtigung der Plan-Tabelle; F-2 bis F-4 sind Übergaben an den Verifier und den Projektinhaber (CI-Lauf auf dem Push-Tip, `rg`-Nachweis).

**Übergabe:** F-1 an den Implementer (Plan-Tabelle, vor der Closure nach `AGENTS.md` §3.10 gegen `git diff --stat` des Endstands abzugleichen). F-2 bis F-4 an den Verifier. Die Finding-Klassen gehen in die Closure-Notiz (§7) des Slice.

## Sensor-Lauf

`make docs-check` nach Anlage dieses Reports: `1191 Datei(en) geprüft, 0 Befund(e)` (der erste Lauf mit dem Report fand eine unverlinkte Slice-Kennung im Pfad-Feld von F-1; berichtigt).
