# Review-Report: slice-v1-release-cut-v0.16.1

* Rolle: Reviewer (frischer Kontext), Datum 2026-10-25
* Gegenstand: Diff `f9cc5e87..HEAD` (Move `c238bc71`, Nachzug `cb62453c`, Start `1d092bd0`, Cut `bde533fc`, Umsetzungsvermerk `34809217`), Slice [`slice-v1-release-cut-v0.16.1`](../plan/planning/in-progress/slice-v1-release-cut-v0.16.1.md)
* Maßstab: Plan, Beleg-Tabelle (§3), `AGENTS.md` §3.7, §3.9, §3.10, §5 Regel 6, `docs/maintainer/releasing.md`, Lastenheft 0.24.0, Messprotokoll `docs/messungen/slice-cr-lastenheft-messstaende/`. Keine DoD-Prüfung (Verifier).

## Zusammenfassung

0 HIGH, 2 MEDIUM, 1 LOW, 3 INFO. Die Versionsstellen, der Binary-Diff, die Vergleichslinks, das Datum und die Sensoren tragen. Zwei Aussagen stehen breiter als ihre Fundstelle: der `podman`-Zusagesatz im Lead und die `ollama`-Messklammer in `README.md`/`README.de.md`.

## Findings

### F-1 MEDIUM
* `kategorie`: MEDIUM
* `quelle`: Lastenheft 0.24.0 `LH-FA-DEV-007` (Zeile 474), Beobachtung `BEO-DOK/zusage-weiter-als-messung`
* `pfad`: `CHANGELOG.md` · "`podman build` (for images from Docker Hub) and the start of `podman` under Docker are promised on the same two."
* `klasse`: zusage-weiter-als-fundstelle (Bedingung fehlt)
* `befund`: Das Lastenheft sagt den Start von `podman` nur auf Hosts zu, die verschachtelte User-Namespaces erlauben (`kernel.apparmor_restrict_unprivileged_userns=0`, in der VM bei Colima) und auf denen `/dev/fuse` und `/dev/net/tun` vorhanden sind; bei Wert `1` endet das Startscript mit Exit-Code `11` (Zeile 474), und `podman build` gilt nur "auf den Hosts, auf denen `podman` ... startet" (Zeile 472). Der Zusagesatz im Lead nennt "the same two" ohne diese Voraussetzung; der Wert `0` erscheint erst im Messsatz darunter, dort als Messumgebung und nicht als Bedingung der Zusage. Die Tabellenzeile in `README.md` (Spalte `v0.16.1`) und `README.de.md` fasst die Zusage ebenso ohne Voraussetzung.
* `risiko`: Ein Nutzer auf Ubuntu 24.04 (Standardwert `1`) liest "promised on a Linux host (amd64)" als erfüllt und bekommt Exit-Code 11; der Lead sagt genau für diesen Fall nichts und dementiert damit den Vertragstext, den er wiedergibt. Beleg-Tabelle §3 (Zeile "`podman build` und der Start ...") verlangt "wie oben", also dieselbe Breite wie die Fundstelle.
* `verifizierbar`: Review-only

### F-2 MEDIUM
* `kategorie`: MEDIUM
* `quelle`: `BEO-DOK/messvariante-belegt-nicht-die-auslieferung`, `AGENTS.md` §3.10, Beleg-Tabelle §3 (Zeile Colima), Hinweis des Implementers in §3 *Ist gegen Plan*
* `pfad`: `README.md` · "measured on Colima, arm64, in two runs on one Linux x86_64 host and in one CI run on `amd64`"; `README.de.md` · "gemessen auf Colima, arm64, in zwei Läufen auf einem Linux-x86_64-Host und in einem CI-Lauf auf `amd64`"
* `klasse`: zusage-weiter-als-messung
* `befund`: Die Klammer ist grammatisch eine Aufzählung "Colima, arm64 / zwei Läufe auf einem Host / ein CI-Lauf" und liest sich als "auf Colima, arm64 gemessen". Der Lead und das Handbuch (`docs/user/benutzerhandbuch.md`, *Grenzen der Messung*: "Unter Colima (`arm64`) ist ein von Hand nachgeschriebenes Probe-Image gemessen, der Bau des erzeugten Dockerfiles nicht") sagen für Colima das Gegenteil des Erzeugnisses. Der Cut hat die Zeile gelesen und bewusst nicht geändert (Plan: "kein Diff"), obwohl der Status-Block mit `v0.16.1` die Version trägt und der Lead als Release-Aussage daneben steht.
* `risiko`: Das README, die erste Seite des Nutzers, behauptet für Colima eine Messung des `ollama`-Features; belegt ist nur ein Probe-Image. Lead und README widersprechen sich im Release.
* `verifizierbar`: Review-only (kein Sensor prüft Aussage gegen Messung)

### F-3 LOW
* `kategorie`: LOW
* `quelle`: `AGENTS.md` §3.10 (Zustandsmarker mit Ablaufpunkt)
* `pfad`: `docs/plan/planning/in-progress/roadmap.md` · "(Push und Tag stehen aus)"; Spalte Slice der Zeile `v0.16.1`: "Push und Tag je nach Freigabe des Projektinhabers"
* `klasse`: ablaufpunkt-marker
* `befund`: Zwei Zustandsmarker laufen mit M1/M5 ab. Der Plan führt sie bis zur Closure (M8, Roadmap-Zeile nach *Abgeschlossene Wellen*), sie sind also heute wahr und plan-konform.
* `risiko`: Bleibt einer nach dem Push stehen, behauptet die Roadmap einen Zustand, der nicht mehr vorliegt. Reiner Merkposten für die Closure.
* `verifizierbar`: Review-only

### F-4 INFO
* `quelle`: Plan §1/§7 Frage 2
* `pfad`: `CHANGELOG.md` · "Patch release."
* `befund`: Einstufung Patch/*Changed* ist Vorgabe des Projektinhabers; der Binary-Diff stützt "Patch" (siehe Negativbefund). Die Gegenposition (Zusage-Menge wächst um Linux `amd64` für `ollama` und `podman`) steht im Plan benannt. Keine Änderung verlangt. Zuständigkeit: Projektinhaber.

### F-5 INFO
* `quelle`: `AGENTS.md` §3.7/§3.9 (a)
* `pfad`: `CHANGELOG.md` · "in words and no longer as the environment of a measurement"
* `befund`: "no longer" ist Release-Prosa über die Form des Vertragstexts, kein Code-/Konfigurations-Kommentar; §3.7 gilt für Kommentare und Zustandsfelder. Der Implementer hat die Stelle ausgewiesen. Kein Finding.

### F-6 INFO
* `quelle`: `harness/review.md` (Sensoren)
* `pfad`: Slice §3 *Sensoren*
* `befund`: `make test-docker` und `make fullbuild` sind laut Vermerk nicht gelaufen (Begründung: Produktcode-Diff eine Zeile; `ci` und `build` gelaufen). Ich habe `make test-docker` auftragsgemäß ebenfalls nicht gefahren. Die Docker-Pfade belegt M2 auf dem Push-Tip; Zuständigkeit Hauptagent/Verifier. `make verify-depguard` entfällt, keine Importregel berührt.

## Geprüft, ohne Befund

* **Binary-Diff**: `git diff v0.16.0..HEAD -- ':!docs' ':!spec' ':!*.md'` besteht aus der Testdatei (nur `//`-Zeilen, die Kommentarzeilen vor `TestE2E_SandboxDevcontainer_Ollama`) und den Versionszeilen in `Dockerfile` (3), `cmd/uboot/main.go` (1), `repo.mk` (3, inklusive Kommentare). "The code of the binary is unchanged; only its version string differs" trägt; kein "byte-identical".
* **Versionsstellen**: `grep -rn "0.16.0-dev"` trifft nur Planungstext des Slice (Plan/Umsetzungsvermerk), nicht Code, Doku oder Workflows. `main.go`, `repo.mk`, `Dockerfile` = `0.16.1-dev`; READMEs (Status, `VERSION=`, Tabellenzeile), Handbuch (Kopf `:4`, `:82`, `:108`, `:135`, `:3207`) = `0.16.1`; `Stand` = 2026-10-25.
* **Datum / M0**: `date +%F` = 2026-10-25 = `CHANGELOG`, Tabellenzeilen, Handbuch-`Stand`.
* **CHANGELOG-Struktur**: `[Unreleased]` leer, `[0.16.1]` mit Lead und unverändertem `### Changed`-Eintrag; Vergleichslinks `[Unreleased]: …v0.16.1...HEAD`, `[0.16.1]: …v0.16.0...v0.16.1` korrekt; kein Advisory, kein *Security*, kein Wort "security"/"vulnerability" im Lead.
* **Lead gegen Lastenheft/Messung** (außer F-1): `ollama` auf Linux `amd64` und Colima (Zeile 412); Cloud-Modelle und `ollama-host` nur Linux `amd64` (Zeilen 414, 418); zwei Mengen getrennt; "two runs on one host" (Handbuch, Lauf 2026-10-19) und ein CI-Pass auf Ubuntu 24.04 (Messprotokoll, Job 121082471598, `--- PASS` für `…_Ollama`); `podman` mit Wert `0` gemessen, CI-Runner Wert `1`, vier `SKIP` (Messprotokoll); Colima nur Probe-Image (Handbuch). Kein Run-Identifier, keine CI-Aussage zum Tip, kein "unified/resolved".
* **Nutzer-Doku ohne Kennungen** (`AGENTS.md` §5 Regel 6): der Diff in `README.md`, `README.de.md` und `docs/user/` trägt keine `LH-`/`ADR-`/`SPEC-`/Slice-Kennung und keinen Planverweis.
* **Pläne gegen `git diff --stat f9cc5e87..HEAD`**: zehn Dateien, jede mit Plan-Eintrag; `releasing.md` unverändert, der Zeiger "zuletzt" (`v0.16.0`) bleibt bis M8 plan-konform (Plan §3, Zeile `releasing.md`); Workflows, `Makefile`, `go.mod`, `scripts/`, `spec/` ohne Diff.
* **M-T-Protokoll**: plausibel; die rote Gegenprobe (2b, 2c, 3b) ist benannt, die Grenze (`is_stable`, echter `curl`/`git push`, Token) ausgewiesen. Nicht nachgefahren (Verifier-Aufgabe).
* **Sensoren, selbst gefahren**: `make gates` Exit 0 (`ablage-pruefung` 26 Fälle, 0 Fehlschläge); `make doc-trace` 84 Anforderungen, 0 Waisen; `make doc-commits RANGE="origin/main..HEAD"` 6 Commits, 0 Befunde.

## Nicht geprüft

* `make test-docker` (gezielt ausgelassen, Auftrag), `make ci`/`make fullbuild`/`make verify-depguard` (kein Produktcode, keine Importregel; Ergebnisse des Implementers nicht reproduziert).
* CI auf dem Push-Tip (M2), Tag, `publish`, Tap: liegen nach dem Push; Hauptagent.
* M-T nicht nachgefahren; DoD-Abgleich: Verifier.
