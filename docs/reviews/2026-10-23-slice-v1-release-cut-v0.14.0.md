# Review-Report: slice-v1-release-cut-v0.14.0 — 2026-10-23

**Rolle:** Reviewer (Modul 10, [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md), [`harness/review.md`](../../harness/review.md)), frischer Kontext. Prüfung gegen Plan, Validierungsbericht, Release-Prozess und Hard Rules, nicht gegen die DoD (Verifier).

**Gegenstand:** Commit-Range `f9680634..08babc88` (Release-Cut `v0.14.0`, Welle `welle-devcontainer-shell-und-zweites-up`, Slice [`slice-v1-release-cut-v0.14.0`](../plan/planning/in-progress/slice-v1-release-cut-v0.14.0.md), Vertrag [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)): Start- und Move-Commits `ff351704`, `8bb06e9f`, `d0cd12c5`, `51a1efe4`, Cut-Commit `9c6fedf1`, Umsetzung §3 `08babc88`.

**Eingang gelesen:** Plan (§1 bis §3, Beleg-Tabelle, Umsetzung), [`2026-10-22-welle-devcontainer-shell-und-zweites-up-validierung`](2026-10-22-welle-devcontainer-shell-und-zweites-up-validierung.md) (Auflagen A-1 bis A-3), [`releasing.md`](../maintainer/releasing.md) §3 bis §5, [`AGENTS.md`](../../AGENTS.md) §3.7, §3.9, §3.10, §5 Regel 6, Lastenheft-Tabelle zu den Exit-Codes und Spezifikation zu den Grenzen der Erkennung bei `shell`.

## Ergebnis

**GO für M1 (Push).** Keine HIGH-, keine MEDIUM-Findings. Ein LOW, drei INFO. Der Push-Tip muss der Commit sein, der diesen Report trägt (Hash in der Übergabe an den Hauptagenten); M2 und der Tag beziehen sich auf ihn.

## Findings

### F-1 — INFO — Grenze (e) und die unbedingte Formulierung in Handbuch 4.9 *Ende*

- `kategorie`: INFO · `quelle`: Maintainability, [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`AGENTS.md`](../../AGENTS.md) §3.10 · `klasse`: doku-unbedingt-gegen-spec-grenze
- `pfad`: `docs/user/benutzerhandbuch.md` · „**Ende.** Hat die Shell begonnen, endet der Befehl mit Exit-Code 0, was immer der Exit-Status der“; ebenso `CHANGELOG.md` *Added* „Once the shell has begun the command ends with exit code `0` whatever the exit status of the shell was“ (beide unverändert, nicht im Diff dieses Cuts).
- `befund`: Das Lastenheft führt die Ausnahme in der Exit-Code-Tabelle („Ausnahme: die Grenzen der Erkennung“), die Spezifikation nennt die Fälle; Handbuch und Changelog-Eintrag sagen den Satz ohne Ausnahme. Der Lead des Cuts selbst macht die Aussage nicht (er nennt nur `0`/`11`/`12` für `up`), er kollidiert also nicht mit (e). Die Handbuch-Zeile ist durch die Bedingung „Hat die Shell begonnen“ für den Regelfall wahr; (e) beschreibt den Fall, in dem eine Meldung der Sitzung selbst als Fehler des Clients gelesen wird (`12` statt `0`).
- `risiko`: gering. Das Szenario setzt eine Sitzung voraus, deren eigene Ausgabe die Meldung zum fehlenden Benutzer erzeugt; Entscheidung (m) des Projektinhabers lässt (e) bewusst aus dem Cut. Ein Befund wäre es erst, wenn der Lead die Aussage unbedingt träfe, und das tut er nicht. Der nächste Änderer des Handbuchs liest den unbedingten Satz als vollständige Zusage.
- `verifizierbar`: Review-only. Der Folge-Slice [`slice-v1-shell-erkennung-fehlender-benutzer-eng`](../plan/planning/open/slice-v1-shell-erkennung-fehlender-benutzer-eng.md) muss Handbuch und Changelog-Eintrag mitziehen, welchen Weg die Entscheidung auch wählt.

### F-2 — LOW — Beleg für „transient, ohne Code-Bezug“ steht auf einer Wortsuche, die Dateiliste trägt ihn

- `kategorie`: LOW · `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.10 (Aussage über Code per Lesung, nicht per Suchwort) · `klasse`: unveraenderter-beleg-per-suche
- `pfad`: [`slice-v1-release-cut-v0.14.0`](../plan/planning/in-progress/slice-v1-release-cut-v0.14.0.md) · „keine Datei mit `download.docker.com` hat sich seit `v0.13.1` geändert“
- `befund`: Die Begründung nennt die Suche nach dem Host-Namen. Gelesen habe ich stattdessen `git diff v0.13.1..HEAD --stat -- internal/e2e`: Geändert sind nur `devcontainer_shell_docker_test.go` (neu) und `devcontainer_start_docker_test.go`; `sandbox_devcontainer_docker_test.go` (dort liegt `TestE2E_SandboxDevcontainer_NestedDocker_EgressAptUpdate`) und alle Generator-/Template-Pfade außerhalb von `internal/adapter`, `internal/hexagon`, `cmd/uboot` sind seit `v0.13.1` unverändert; das `Dockerfile` hat nur die Versionszeilen geändert. Der Code-Bezug ist damit ausgeschlossen, aber über die Dateiliste, nicht über die genannte Suche. Zwei von drei Läufen scheiterten am selben Netzfehler (`unexpected eof`), einer bestand: das stützt „Netzweg“, belegt aber keine Rate.
- `risiko`: Der Satz „transient“ steht im Plan als Befund; ein späterer Leser übernimmt ihn ohne die Dateiliste. Kein Einfluss auf die Release-Entscheidung.
- `verifizierbar`: Review-only (Belege: `git diff v0.13.1..HEAD --stat -- internal/e2e`); Gegenprobe ist M2 (siehe INFO-Hinweis zu M2).

### F-3 — INFO — Datumswechsel: Ortszeit 23., UTC noch 22.

- `kategorie`: INFO · `quelle`: [`releasing.md`](../maintainer/releasing.md) §5 (Datum gegen Tag-Datum halten) · `klasse`: datum-utc-gegen-ortszeit
- `pfad`: `CHANGELOG.md` · „## [0.14.0] - 2026-10-23“; Plan §3 *Umsetzung* M0.
- `befund`: Beim Review ist `date +%F` 2026-10-23 (CEST, 00:48), `date -u` 2026-10-22. CHANGELOG, beide Tabellenzeilen und Handbuch-`Stand` tragen konsistent den 23.; der Hinweis des Implementers im Plan ist zutreffend und genügt. `git log -1 --format=%cs` liest die Ortszeit des Committers (`08babc88`: 2026-10-23), der Tag-Commit trägt damit den 23.
- `risiko`: Wird der Tag vor 02:00 CEST gepusht, steht der `publish`-Lauf und das GitHub-Release in UTC auf dem 22. Das ist ein Darstellungsunterschied, kein Fehler im Text, solange `%cs` des Tag-Commits der Quelle des Datums folgt (§7 von `releasing.md`).
- `verifizierbar`: Hauptagent vor M1 und M5 (`date +%F`, `git log -1 --format=%cs`).

### F-4 — INFO — Nicht entfernbare Container aus `make test-docker`

- `kategorie`: INFO (Zuständigkeit Hauptagent/Maintainer) · `quelle`: Maintainability · `klasse`: testrest-ausserhalb-arbeitsbaum
- `pfad`: Plan §3 *Umsetzung* · „zwei beendete Container des fehlgeschlagenen Baus (`funny_dhawan`, `awesome_shockley`“
- `befund`: Der lokale Lauf hinterließ zwei beendete Container, die der Implementer nicht entfernen durfte. Der Arbeitsbaum ist sauber (`git status` leer).
- `risiko`: Nur Host-Hygiene; kein Einfluss auf den Diff.
- `verifizierbar`: Review-only.

## Prüfung nach Schwerpunkten

1. **Versionsstellen.** `grep` nach `0.13.1`, `0.13.0`, `0.14.0` über den Baum außer `docs/plan`, `docs/reviews`, `.harness/baseline`, CHANGELOG. `cmd/uboot/main.go:43`, `repo.mk:51/54/55` und `Dockerfile:35/37/39` tragen übereinstimmend `0.14.0-dev` (Muster der Vorgänger: `<neu>-dev`, `releasing.md` §3 Punkt 2). Handbuch: Kopf, `:82`, `:108`, `:135`, Gültigkeitsbereich `v0.14.0`; READMEs: Status, `VERSION=`, Tabellenzeile. Verbleibend `0.13.x`: Tabellenzeilen `v0.13.0`/`v0.13.1` (Historie), `releasing.md:57` (Zeiger „zuletzt“, Closure), Handbuch `:955` (Messstand „`u-boot` 0.13.0, Docker 29.9.1“, richtig als alter Messstand belassen). Spec-Treffer in `lastenheft.md`/`spezifikation.md` sind Dokumentversionen, nicht die Release-Nummer. Vergleichslinks am Fuß: `[Unreleased]` → `v0.14.0...HEAD`, `[0.14.0]` → `v0.13.1...v0.14.0`, vor `[0.13.1]`; Format wie die Vorgänger. Das Tag `v0.14.0` existiert nicht (nur `v0.13.1` und älter).
2. **Aussagen über Verhalten (§3.10).** Lead gegen Validierungsbericht und Plan-Beleg-Tabelle gelesen: Exit `12` (ein `postStartCommand` endet mit 3) und Strg-C (Rückkehr nach 10,5 s) mit folgendem `up` Exit `0`/`already running`: Validator Punkt (ii), im Lead als „one run each, Linux `amd64`, Docker 29.9.2“ gekennzeichnet; „Andere Fehler (etwa ein gescheiterter Bau) sind nicht gemessen“ steht dabei. Zähler der Lebenszyklus-Befehle: Docker-Test `TestDevcontainerUp_Docker_AlreadyRunningLifecycleCounters` (in der lokalen Zählung PASS). `--recreate` bei fortbestehender Ursache wieder `11`: Validator, ein Lauf, so im Text. Zusage-Menge (Linux `amd64`, lokale Engine, Docker 29.9.2) und alle sieben Nicht-Zusagen (Colima auch bei Homebrew, Docker Desktop, rootless, Podman als Host-Engine, entfernte Engine, `arm64`, Windows) stehen im Lead und in der Tabellenzeile; das Terminal ist als „unter `script`“ gekennzeichnet. Keine Aussage ohne Beleg gefunden; der CI-Satz fehlt bewusst im Lead (Abweichung 1 im Plan), die READMEs behalten den bestehenden Satz. Grenze (e): siehe F-1, der Lead trifft keine Aussage, die (e) widerlegt.
3. **Kennungsfreiheit.** `grep` auf `LH-`, `ADR-`, `SPEC-`, `MR-`, `ARC-`, `slice-`, `welle-`, `roadmap` in den hinzugefügten Zeilen von `README.md`, `README.de.md` und `docs/user/`: 0 Treffer. Im CHANGELOG nur „Lastenheft 0.21.0“.
4. **Datum.** 2026-10-23 an allen vier Stellen konsistent (F-3).
5. **Sensoren und M-T.** `make gates` auf dem Stand `08babc88` von mir gelaufen: Exit 0 (`lint`, `test`, `coverage-gate`, `docs-check`, `baseline-verify`, `ablage-pruefung` 26 Fälle, 0 Fehlschläge). Nicht von mir wiederholt: `make ci`, `make fullbuild`, `make build VERSION=0.14.0`, `make test-docker` (Zeit; die Berichte des Implementers sind in Form und Zahlen plausibel: `--version` mit und ohne `-dev`, Trivy ohne Fund, M-T rot gesehen bei `v0.13.1`). `make verify-depguard` entfällt zu Recht (keine Importregel berührt, einziger Go-Diff ist `var version`). Zum `test-docker`-Befund: Das Paket-Zeitlimit mit Exit 2 entspricht den Cuts `v0.13.0`/`v0.13.1` ([`BEO-E2E/e2e-paket-am-zeitlimit`](../plan/planning/observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md)); die neun neuen Tests liefen vor dem Limit und bestanden. Der eine `--- FAIL` ist als Netzstörung plausibel (F-2: die Dateiliste schließt Code-Bezug aus). **Für Verifier und M2:** (a) `integration-docker` auf dem Push-Tip muss `--- FAIL` 0 und `EgressAptUpdate` als `PASS` zeigen und darf kein `panic: test timed out` tragen; (b) die Tests nach dem Ollama-Test im Paket `internal/e2e` liefen lokal nicht bis zum Ergebnis, ihre Abdeckung kommt nur aus der CI; (c) lokal liefen die Nested-Podman-Tests (0 SKIP), die CI erwartet vier SKIP.
6. **Kommentare, Zustandsmarker, Plan-Tabelle, Push-Reife.** §3.9: Die Kommentare in `repo.mk` und `Dockerfile` nennen den geltenden Wert, keiner beschreibt Entferntes; der alte Wert `0.13.1-dev` steht außerhalb von Planung/Reviews/Changelog nicht mehr. Plan-Tabelle gegen `git diff --stat f9680634..HEAD` (12 Dateien) abgeglichen: jede Datei hat einen Eintrag (CHANGELOG, `main.go`, `repo.mk`, `Dockerfile`, beide READMEs, Handbuch, Roadmap, Welle-Datei, `open/README.md`, der Slice), der Validierungsbericht liegt vor dem Plan-Commit; `releasing.md`, `spec/`, `internal/` haben keinen Diff, wie der Plan sagt. Zustandsmarker: Roadmap „(Push und Tag stehen aus)“ und die `Stand`-Zeile der Welle-Datei nennen den Ist-Stand; der Marker läuft mit M8 ab und ist dort vorgesehen. Push-Reife: `8bb06e9f` (reiner Move) ist allein docs-check-rot und wird von `d0cd12c5` repariert; die Spitze ist grün. Die Commits gehören in einem Push.

## Negativbefund (geprüft, ohne Befund)

- `CHANGELOG.md` `[0.14.0]` (Lead, Vergleichslinks) · `cmd/uboot/main.go`, `repo.mk`, `Dockerfile` (Versionsstellen, Kommentare) · `README.md`, `README.de.md` (Status, Tabellenzeile, `VERSION=`) · `docs/user/benutzerhandbuch.md` (Versionsstellen, A-3 an beiden Stellen) · Kennungsfreiheit der Nutzerdoku · Datumskonsistenz · Roadmap-Zeiger, Welle-Datei, `open/README.md` · Plan-Tabellenabgleich · Auflagen A-1, A-2, A-3 des Validators (alle drei im Text umgesetzt) · Commit-Messages (Slice-Kennung und [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), keine `SPEC-`/`ARC-`) · Lifecycle (Move als reiner Rename-Commit, Verweis-Nachzug getrennt, [`AGENTS.md`](../../AGENTS.md) §3.3).

## Nicht geprüft

- `make ci`, `make fullbuild`, `make build VERSION=0.14.0`, `make test-docker`, M-T und `make verify-depguard`: nicht wiederholt (Laufzeit, Docker-Engine-Last, Netzabhängigkeit; Implementer-Protokoll gelesen, Wiederholung ist Sache des Verifiers und von M2). `make verify-depguard` entfällt mit Begründung (siehe oben).
- Das Verhalten des Binaries, die CI-Läufe auf dem Push-Tip (M2), `tap-check` (M4): Sache von Verifier und Hauptagent.
- Die tatsächliche Gültigkeit der Messungen des Validators: Eingang, nicht Ergebnis dieses Reviews.

## Übergabe

GO für M1. Vor dem Push `date +%F` (F-3) und `make docs-check`; nach dem Push M2 mit den Prüfpunkten aus Abschnitt 5. Der Tag gehört auf den in M2 geprüften Push-Tip, das ist der Commit mit diesem Report.
