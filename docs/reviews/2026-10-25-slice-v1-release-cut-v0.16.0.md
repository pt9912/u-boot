# Review-Report: slice-v1-release-cut-v0.16.0

* Rolle: Reviewer (`.harness/skills/reviewer.md`), frischer Kontext
* Gegenstand: `git diff 7956e157..HEAD` (HEAD `d1c1abc3`; Move `2b5714c7`, Verweis-Nachzug `ea309e70`, Start `e8c42ef1`, Cut `39e4772d`, Umsetzung `f35f2e8c`, Gate-Eintrag `d1c1abc3`), Slice
  [`slice-v1-release-cut-v0.16.0`](../plan/planning/in-progress/slice-v1-release-cut-v0.16.0.md)
  (Minor `v0.16.0`, Einstufung *Added*/*Changed*, kein Advisory); Muster: Cuts `v0.15.0` und `v0.15.1`
* Bezug: [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md), Validator-Bericht [`2026-10-24-welle-ollama-cloud-modelle-validierung`](2026-10-24-welle-ollama-cloud-modelle-validierung.md) (Auflagen A-1 bis A-3)
* Nicht Gegenstand: DoD-Prüfung (Verifier)

## Urteil

Kein HIGH, kein MEDIUM. Zwei LOW (abgelaufene Zustandsmarker in der Welle-Datei; Lead und Handbuch-Kopf gegen das Tag-Datum), drei INFO.
Der Lead trägt jede geprüfte Aussage an einer Messung, einer Zusage-Zeile oder einem Test; A-1 ist vollständig, A-2 und A-3 sind im Handbuch mit ausgewiesener Reichweite umgesetzt.

## Sensoren (selbst gefahren)

| Sensor | Ergebnis |
|---|---|
| `make gates` | Exit 0 (Log im Scratchpad `rev16-gates.log`; `ablage-pruefung` 26 Fälle, 0 Fehlschläge als letzter Block) |
| `make doc-trace` | Exit 0, 84 Anforderungen, 0 Waisen |
| `make doc-commits RANGE="origin/main..HEAD"` | Exit 0, 22 Commits, 0 Befunde (der Implementer meldete 20; zwei Commits sind seither hinzugekommen) |
| `make test-docker`, `make ci`, `make fullbuild` | nicht selbst gefahren: Auftrag (Speicher); der Diff berührt nur `var version` als Produktcode. Die Docker-Tests der Welle (zwei `TestE2E_OllamaHost_*`) las ich im Log des Implementers (`docker-targeted.log`: 2 PASS, 0 FAIL, 0 SKIP), reproduzierte sie nicht |
| `make build VERSION=0.16.0` | nicht gefahren (Speicher); die Versionsstellen las ich statt dessen (siehe Linse 2) |
| `make verify-depguard` | nicht gelaufen: keine Importregel berührt |

## Findings

### LOW-1

* `kategorie`: LOW
* `quelle`: `AGENTS.md` §3.10 (Zustandsmarker mit Ablaufpunkt werden beim Ablauf berichtigt)
* `pfad`: `docs/plan/planning/welle-ollama-cloud-modelle.md` · „Aufträge an den Release-Cut“ (§7, Absatz mit „`Stand` (heute das Schreibdatum `2026-10-24`) und die `Software-Version` (heute `v0.15.1`)“ sowie „unter `[Unreleased]` / *Changed* `v0.16.0` zuordnen“)
* `klasse`: abgelaufener-zustandsmarker-in-welle-datei
* `befund`: Der Cut-Diff berührt diesen Absatz nur im Link, nicht im Inhalt. Nach dem Cut-Commit steht das Handbuch auf `v0.16.0` / `2026-10-25` und der Eintrag unter `[0.16.0]`; „heute“ und „unter `[Unreleased]`“ beschreiben den Zustand vor dem Cut.
* `risiko`: Der nächste Leser (Welle-Closure) liest „heute v0.15.1“ als Ist-Zustand; kein Gate fängt es. Der Start-Commit hat die übrigen Marker der Datei berichtigt, diesen nicht (Plan-Tabelle führte nur „Abgelaufene Zustandsmarker“ der Zeilen 135/148).
* `verifizierbar`: Review-only

### LOW-2

* `kategorie`: LOW
* `quelle`: `docs/maintainer/releasing.md` §5 (Datum gegen das Tag-Datum halten), Slice §3 *M0*
* `pfad`: `CHANGELOG.md` · „## [0.16.0] - 2026-10-25“; gleichlautend Tabellenzeile beider READMEs und Handbuch „Stand: 2026-10-25“
* `klasse`: datum-gegen-tag-datum
* `befund`: Das Datum ist `date +%F` am Tag der Prüfung (2026-10-25, Ortszeit) und stimmt. Es gilt nur, wenn der Tag noch am 2026-10-25 gesetzt wird; der Slice benennt das selbst (M0).
* `risiko`: Wird der Tag später gesetzt, stehen vier Stellen auf dem falschen Datum und es braucht einen weiteren Cut-Commit vor dem Tag. Kein Gate prüft es. Als LOW geführt, weil der Plan die Bedingung trägt; die Prüfung liegt beim Hauptagenten vor M1/M5.
* `verifizierbar`: Review-only (`date +%F` vor dem Tag)

### INFO-1

* `kategorie`: INFO
* `quelle`: Maintainability
* `pfad`: `CHANGELOG.md` · „It also states the measured effect of `OLLAMA_NO_CLOUD` and of model names with `[1m]`, and where the device key lives (the volume of the models).“
* `klasse`: changelog-nennt-neue-handbuchsaetze-nicht
* `befund`: Die Bullets *Added*/*Changed* blieben im Wortlaut (Plan). Die im Cut hinzugekommenen Handbuch-Sätze (Rückbau des Geräteschlüssels, Durchreichung von `ANTHROPIC_API_KEY`, Eingrenzung des Satzes zur ungesetzten Host-Variable) stehen nicht im CHANGELOG. Nichts Falsches, nur nicht genannt.
* `risiko`: Keines für den Vertrag; ein Nutzer, der nur den CHANGELOG liest, erfährt vom Rückbau nicht. Der Implementer entscheidet.
* `verifizierbar`: Review-only

### INFO-2

* `kategorie`: INFO (Zuständigkeit Hauptagent/Verifier)
* `quelle`: `BEO-DOK/messvariante-belegt-nicht-die-auslieferung`
* `pfad`: `CHANGELOG.md` · „**Measured at the binary** (one run, Linux `amd64`)“; Slice §3 *M-V*
* `klasse`: messvariante-binary-aus-arbeitsbaum
* `befund`: Gemessen wurde an einem Binary aus dem Arbeitsbaum (`VERSION=0.16.0-cut16`), nicht am ausgelieferten Artefakt. Der Unterschied liegt nur im Versionsstring (kein Produktcode im Diff); der Lead sagt „at the binary“ ohne Zusatz. Im Vorbild `v0.15.1` gleiche Form.
* `risiko`: Gering; die Reichweite ist durch „one run“ und „The path with a `source` was not run“ begrenzt. Die Prüfung `--version` = Tag liegt im Workflow `publish`.
* `verifizierbar`: Review-only

### INFO-3

* `kategorie`: INFO (Zuständigkeit Verifier)
* `quelle`: Slice §2 DoD, `harness/verification.md`
* `pfad`: Slice · „`make test-docker` ist gelaufen (Zahl der übersprungenen Tests)“
* `klasse`: dod-punkt-offen-test-docker
* `befund`: Der Implementer nennt das Kriterium selbst als nicht erfüllt (gezielter Lauf statt vollem `make test-docker`, M2 steht aus). Roadmap-Zeile „(Push und Tag stehen aus)“ und der Kopf „Stand vor M1: nichts gepusht, kein Tag“ sind Marker mit Ablaufpunkt (Closure M8).
* `risiko`: Kein Befund am Diff; Hinweis, dass die Closure beide Marker entfernen muss.
* `verifizierbar`: `make test-docker` / CI auf dem Push-Tip

## Negativbefunde (geprüft, ohne Befund)

* **CHANGELOG `[0.16.0]` Lead (Wortlaut):** jede Aussage gegen Validator-Bericht, [ADR-0036](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md), Slice-Beleg-Tabelle und M-V gelesen. Zusage-Menge Linux `amd64`, lokale Engine, Docker 29.9.2 (Validator Messumgebung); Nicht-Zusagen identisch zur Reihe `[0.15.0]`; Weg über `devcontainer up` nicht gemessen und nicht zugesagt, mit dem, was ohne Konto am Pfad gemessen ist (z6: Schlüsselpaar im Volume; a2: `ollama.com` löst nach Eintrag und `up --recreate`; Meldungswechsel und Anmeldung ausdrücklich nicht gemessen); Neuanlage mit demselben Volume nicht zugesagt (Lastenheft-Zeile laut Slice); Host-Anschluss liefert den Namen, nicht den Dienst, hinter der Egress-Restriktion „trägt nicht“ ([ADR-0036](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) Entscheidung 4, gemessen: `claude -p` Exit 124), `generate` warnt ohne Adresse (Docker-Test D2: Allowlist ohne Host); Exit `10` für `ollama-host.version`, `ollama.version` und mit `--dry-run`, `u-boot.yaml` unverändert (M-V-Log `i3.log`); zwei Zusage-Mengen als Grenze (Carveout 2). **A-1 ist vollständig** (Linux `amd64`, Weg über `up` nicht gemessen, Relay in einem Satz, Grenze der zwei Mengen). Keine Aussage geht über das Gemessene hinaus. Einstufung *Added*/*Changed*, kein Advisory, plan-konform; „Sixteenth release“ stimmt zur Reihe.
* **Vergleichslinks, `[Unreleased]`:** `[Unreleased]` leer, `…/compare/v0.16.0...HEAD`, neue Zeile `[0.16.0]: …/compare/v0.15.1...v0.16.0` vor `[0.15.1]`.
* **Version `0.16.0` an allen Stellen:** `cmd/uboot/main.go` (`var version = "0.16.0-dev"`), `repo.mk` (Kommentar und `VERSION ?=`), `Dockerfile` (Kommentare und `ARG UBOOT_VERSION`), beide READMEs (Status, `VERSION=0.16.0`, Tabellenzeile nach `v0.15.1`), Handbuch (Kopf, `u-boot version 0.16.0`, `download/v0.16.0`, `VERSION=0.16.0`, Gültigkeitsbereich). Wortsuche `0.15.1` über das Repo außerhalb von Planung, Reviews, Messprotokollen: nur Tabellenzeilen `v0.15.1`, CHANGELOG-Abschnitt und Links, `releasing.md` (Zeiger „zuletzt“, plan-konform bis zur Closure M8). `.github/`, `scripts/`, `packaging/`, `Makefile`, `go.mod` im Diff nicht vorhanden. Kein weiterer 0.15.x-Rest.
* **Handbuch A-2:** Satz zu `ANTHROPIC_API_KEY` nennt nur b1/b2 als gemessen (je ein Lauf `u-boot devcontainer shell`), die Durchreichung ausdrücklich als „abgeleitet aus der Referenz, nicht mit einem Schlüssel gemessen“, `unset`/`env -u` als nicht gemessen und nicht empfohlen; Code-Fundstelle (`agentCredentialEnv`, `remoteEnv`) mit Slice-Zeile belegt.
* **Handbuch A-3:** Volume-Name `<projekt>-ollama-<Kennung>` stimmt zu `volumeMount(name+"-ollama", …)` in `devcontainer_sandbox.go` und zum gemessenen `vcut16-ollama-<devcontainerId>`; Reihenfolge `docker rm` vor `docker volume rm` und „volume is in use“ gegen `rm1.log`/`rm2.log` gelesen, „(not removed)“ gegen `down.log`; kein Satz zu `ollama signout`; „ob die Seite den Schlüssel anzeigt, nicht gemessen“ ausgewiesen.
* **Handbuch V-1:** Der Satz ist auf die Pfade eingegrenzt, die tragen: ungesetzt gegen b1 (`u-boot devcontainer shell`, Linux `amd64`), „Not logged in“ gegen die Messung in [`slice-v1-devcontainer-feature-claude-code`](../plan/planning/done/slice-v1-devcontainer-feature-claude-code.md) (Dev-Container-CLI, Colima `arm64`), mit dem Zusatz „mit `u-boot devcontainer shell` nicht gefahren“. Keine neue Zusage.
* **Handbuch-Kopf, Nutzer-Doku-Regel (`AGENTS.md` §5 Regel 6):** `Software-Version: v0.16.0`, `Stand: 2026-10-25`; im Diff von `README.md`, `README.de.md` und `docs/user/` keine `LH-`/`ADR-`/`SPEC-`-Kennung und kein Verweis auf Plan- oder Entwicklerdokumente.
* **READMEs, Klammer-Änderung `config set`:** „`gh`, `glab`, `java`, `ollama` und `ollama-host`“ (englisch und deutsch gleichlautend) trägt das Binary-Log (Exit `10` für beide `ollama*`); der Eintrag `[0.15.1]` und die Tabellenzeile `v0.15.1` bleiben als Geschichte unverändert. Katalogzeile und Tabellenzeile sagen Zusage-Menge, Nicht-Zusagen, Weg über `up` nicht gemessen, Relay, ohne Datum im Status-Block und ohne „released“.
* **M-T (Tap-Trockenlauf):** Skript und Ausgabe (`mt.sh`, `mt.out`) gegen `releasing.md` §4 gelesen: Fill Exit 0 mit vier `v0.16.0`-URLs, Check Exit 0, mit Vorgänger-Tag Exit 1, gemischte Formel Exit 1; `tap-nachzug` zeigt `clone`, `diff --quiet`, `add`, `commit -m u-boot 0.16.0`, `push`; mit `v0.15.1` bricht es vor `clone` ab (0 Clone-Aufrufe). Die Grenzen (`is_stable`, echter `curl`/`push`, Token) benennt der Slice selbst. Plausibel, deckt sich mit dem erwarteten Protokoll.
* **Plan-Tabelle gegen `git diff --stat 7956e157..HEAD`:** alle elf Dateien haben einen Plan-Eintrag und umgekehrt; „kein Diff“-Zeilen (`releasing.md`, `carveouts.md`, `spec/`, `docs/plan/adr/`, Workflows, `internal/`) sind im Diff tatsächlich nicht vorhanden und im Slice durch gelesene Fundstellen belegt. `releasing.md` bleibt bis zur Closure bei `v0.15.1`: plan-konform (Plan-Tabelle und Umsetzungs-Punkt 6, Abweichung vom Auftrag des Hauptagenten offen benannt).
* **Moves und Lifecycle:** `2b5714c7` und `ea309e70` sind getrennte Commits (reiner Move, dann Verweis-Nachzug); der Slice liegt nur in `in-progress/`; `doc-commits` ohne Befund.
* **Kommentare im Diff (§3.7/§3.9):** die Kommentare in `repo.mk` und `Dockerfile` nennen den geltenden Wert; kein Kommentar beschreibt Entferntes.
* **Immutabilität/Spec-Straten:** `spec/` und `docs/plan/adr/` ohne Diff; die Spec-Straten verweisen im Diff nicht abwärts.

## Nicht geprüft

* Volles `make test-docker`, `make ci`, `make fullbuild`, `make build VERSION=0.16.0`: Auftrag/Speicher (siehe Sensoren).
* Inhalt der Rohlogs der Validator-Läufe b1/b2/z6: nur über den Validator-Bericht gelesen, nicht erneut gefahren.
* Der Lastenheft-Wortlaut zur Neuanlage mit demselben Volume (Zeile 414): über die Zitierung im Slice gelesen, nicht im Lastenheft nachgeschlagen.
* Die Wirkung von Push, Tag, `publish` und Tap: liegt nach dem Review (Hauptagent, M1 bis M6).
