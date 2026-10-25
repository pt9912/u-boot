# Verification-Report: slice-v1-release-cut-v0.16.0

* Rolle: Verifier (`harness/verification.md`), frischer Kontext
* Gegenstand: `git diff 7956e157..HEAD` (HEAD `c2cf2202`; Move `2b5714c7`, Verweis-Nachzug `ea309e70`, Start `e8c42ef1`, Cut `39e4772d`, Umsetzung `f35f2e8c`, Gate-Eintrag `d1c1abc3`, Review-Report `1b28f793`/`b0ba2256`, Fix LOW-1 `c2cf2202`), Slice
  [`slice-v1-release-cut-v0.16.0`](../plan/planning/done/slice-v1-release-cut-v0.16.0.md) (Minor `v0.16.0`)
* Bezug: [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md); Review-Report [`2026-10-25-slice-v1-release-cut-v0.16.0`](2026-10-25-slice-v1-release-cut-v0.16.0.md); Validator-Bericht [`2026-10-24-welle-ollama-cloud-modelle-validierung`](2026-10-24-welle-ollama-cloud-modelle-validierung.md)
* Nicht Gegenstand: Push, Tag, `publish`, Tap (Hauptagent, M1 bis M6) und die Closure (M8)

## Urteil

**Bereit für M1 (Push), mit zwei Bedingungen des Hauptagenten, die kein Gate trägt:** (1) das Datum 2026-10-25 gilt nur, wenn der Tag noch am 2026-10-25 (Ortszeit) gesetzt wird; (2) der Lead im Wortlaut geht dem Projektinhaber vor M1 vor (Frage 1 des Slice). Kein Befund, der Produktcode, Pipeline oder Spec verlangt. Der Release-Text geht an keiner Stelle über das Gemessene hinaus, die Auflagen A-1 bis A-3 sind erfüllt, LOW-1 des Reviews ist behoben. Die DoD-Punkte, die erst nach dem Push entstehen (Punkt 3 mit M2, M5, M6; Closure), sind erwartungsgemäß offen und unten einzeln benannt.

## Sensoren (selbst gefahren, `HEAD` `c2cf2202`, Arbeitsbaum sauber)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | Exit 0 | Log im Scratchpad `ver-gates.log`: `baseline-verify` v6.18.1 OK 54 Dateien; `lint` 0 issues; `coverage-gate` total 92,1 % gegen 90 %; `docs-check` 881 Dateien, 0 Befunde; `ablage-pruefung` 26 Fälle, 0 Fehlschläge |
| `make doc-trace` | Exit 0 | 84 Anforderungen, 0 Waisen |
| `make doc-commits RANGE="origin/main..HEAD"` | Exit 0 | 25 Commits, 0 Befunde |
| `make doc-immutable RANGE="v0.15.1..HEAD"` | Exit 0 | 112 Commits, 0 Befunde (`spec/` und `docs/plan/adr/` ohne Diff) |
| `make doc-planning` | Exit 0 | 881 Dateien, 0 Befunde |
| `make build VERSION=0.16.0` | Exit 0 | `docker run --rm u-boot:latest --version` = `u-boot version 0.16.0`; OCI-Label `org.opencontainers.image.version` = `0.16.0` (lokal, nicht gepusht) |
| M-T (Tap-Trockenlauf) | nachgefahren, identisch | Skript `cut16/mt.sh` des Implementers im Container `u-boot:test`, `--network none`, Repo `:ro`; die Zeilen mit `exit=` und die Zahl der Clone-Aufrufe sind gleich zu `mt.out` (diff leer); siehe unten |
| `make test-docker`, `make ci`, `make fullbuild` | nicht gefahren | Auftrag (Speicher); der Diff berührt als Produktcode nur `var version`; Implementer-Belege `ci.log`/`fullbuild.log`/`docker-targeted.log` im Scratchpad `cut16/` gelesen, nicht reproduziert (siehe DoD Punkt 3) |
| `make verify-depguard` | nicht gelaufen | keine Importregel, kein Go-Import berührt |

Aufgeräumt: das Image `u-boot:latest` (Bau dieser Prüfung) ist entfernt; es bestehen keine Container und Volumes dieser Prüfung. Die Bilder `u-boot:test`, `u-boot:coverage`, `u-boot:lint` stammen aus `make gates` und bleiben als Cache.

## DoD, Punkt für Punkt (Slice §2)

| # | DoD-Punkt | Urteil | Beleg |
|---|---|---|---|
| 1 | Release-Text und Versionsstellen | **erfüllt** | `CHANGELOG.md`: `## [0.16.0] - 2026-10-25` mit Lead, `## [Unreleased]` leer davor, Vergleichslinks `[Unreleased]: …/v0.16.0...HEAD` und `[0.16.0]: …/v0.15.1...v0.16.0` vor `[0.15.1]`. `0.16.0-dev` in `cmd/uboot/main.go` (`var version`), `repo.mk` (Kommentar und `VERSION ?=`), `Dockerfile` (Kommentare und `ARG UBOOT_VERSION`). Beide READMEs: Status `v0.16.0`, `VERSION=0.16.0`, Tabellenzeile `v0.16.0` nach `v0.15.1`. Handbuch: Kopf `v0.16.0`/`2026-10-25`, `u-boot version 0.16.0`, `download/v0.16.0`, `VERSION=0.16.0`, Gültigkeitsbereich. Beleg-Tabelle (Slice §3) trägt jede Aussage; Abgleich siehe unten. Keine Kennung in `README.md`, `README.de.md`, `docs/user/` (Suche nach `LH-`, `ADR-`, `SPEC-`, `ARC-`, `/plan/`, `spec/spezifikation`, `architecture.md`: 0 Treffer) |
| 2 | Handbuch-Auflagen | **erfüllt** | A-2 und A-3 je als Satz bzw. Absatz mit Reichweite, siehe Auflagen. L-1, L-3, I-1 sind in §3 *Folge-Kandidaten* als Folge ausgewiesen, nicht entschieden |
| 3 | Tag und Kontrolle | **offen, erwartungsgemäß** | Vor dem Tag fehlt: voller `make test-docker` (nur die zwei `TestE2E_OllamaHost_*` gezielt, 2 PASS, 0 FAIL, 0 SKIP laut Implementer-Log; der Slice weist das selbst als nicht erfüllt aus), M2 (CI auf dem Push-Tip, Zählung der SKIPs), M4 (`tap-check`), M5 (Tag), M6 (Kontrolle; alle fünf Unterpunkte). M-T ist festgehalten und nachgefahren (grün) |
| 4 | `make gates` grün nach dem letzten Commit | **erfüllt** für `c2cf2202` | dieser Lauf; nach der Closure erneut nötig |
| 5 | Review durchgeführt, Report liegt vor | **erfüllt** | [Review-Report](2026-10-25-slice-v1-release-cut-v0.16.0.md); LOW-1 nach dem Report behoben (`c2cf2202`); dieser Report ist der Verification-Report vor M1 |
| 6 | Befunde des Validators haben einen Ausgang | **erfüllt im Text, Register offen** | A-1 bis A-3 im Text mit Beleg; M-1 als Text im Lead, Register-Evidence ausdrücklich Closure; M-2 im Lead; M-3 als Satz (A-2); L-1 Folgepflicht; L-2 als Satz (A-3); L-3, I-1 Folge; I-2 ohne Aufgabe; I-3 im Lead |
| 7 | Doku-Update (CHANGELOG, READMEs, Handbuch, `releasing.md`, Roadmap, Zustandsmarker) | **teils, plan-konform** | CHANGELOG, READMEs, Handbuch, Roadmap-Zeiger, Welle-Datei und `open/README.md` tragen den Ist-Stand. `releasing.md` (Zeiger „zuletzt“) und die Roadmap-Zeilen mit Tag-Hash sind Closure/M8 (Plan-Tabelle); die Abweichung vom Auftrag des Hauptagenten (Zeiger schon im Cut) ist im Slice offen benannt und folgt dem Plan |
| 8 | Closure-Notiz mit Steering-Loop-Lerneintrag | offen | Closure, nach M6 |
| 9 | Beobachtungs-Register fortgeschrieben | offen | Closure (Evidence-Datei `welle-ollama-cloud-modelle.md`, Beleg M-1) |
| 10 | Jedes Risiko aus §6 trägt einen Ausgang | offen | §6 alle „offen bis zur Closure“ |
| 11 | Drei Paarungen getragen | offen | Closure; Folge-Kandidaten stehen ohne Datei (Entscheidung Frage 6) |

## Review-Report: Behandlung der Findings

| Finding | Stand | Beleg |
|---|---|---|
| LOW-1 (abgelaufene Marker in den Aufträgen an den Cut, Welle-Datei) | **behoben** | `c2cf2202` ersetzt „unter `[Unreleased]` / *Changed* zuordnen“, „Stand (heute …)“, „Software-Version (heute `v0.15.1`)“ durch den Ist-Stand (`Stand` `2026-10-25`, `Software-Version` `v0.16.0`). Die Welle-Datei führt sonst nur Marker mit Ablaufpunkt Cut/Closure („der Release-Cut steht aus“, „steht aus“, Zeilen 141, 149) |
| LOW-2 (Datum gegen Tag-Datum) | **offen, Bedingung für den Hauptagenten** | siehe Datum |
| INFO-1 (CHANGELOG nennt die neuen Handbuch-Sätze zu A-2/A-3 nicht) | offen, nicht blockierend | Der Lead sagt Geräteschlüssel im Volume und das Changed-Bullet nennt „where the device key lives“; Rückbau und Durchreichung von `ANTHROPIC_API_KEY` fehlen im CHANGELOG. Weder Plan noch Validator verlangen es (der Plan lässt die Bullets im Wortlaut); Entscheidung beim Hauptagenten/Projektinhaber, wenn der Lead im Wortlaut vorgelegt wird |
| INFO-2 (Messung am Arbeitsbaum-Binary `0.16.0-cut16`) | zur Kenntnis | Unterschied nur der Versionsstring; kein Produktcode im Diff. Der Lead sagt „at the binary (one run)“, „The path with a `source` was not run“ |
| INFO-3 (Marker „Push und Tag stehen aus“, „Stand vor M1“) | zur Kenntnis | Roadmap-Zeiger und Kopf der Umsetzung in §3: Marker mit Ablaufpunkt, die Closure (M8) muss beide entfernen |

## Auflagen des Validators

| Auflage | Urteil | Beleg |
|---|---|---|
| A-1 (Lead) | **erfüllt** | Der Lead nennt Linux `amd64`, lokale Docker-Engine (Docker 29.9.2 gemessen: Validator *Messumgebung* und Protokoll mit Konto, Messstand), die Nicht-Zusagen (macOS mit Colima auch für Homebrew-Nutzer, Docker Desktop, rootless Docker, Podman als Host-Engine, entfernte Engine, `arm64`, Windows); der Weg über `u-boot devcontainer up` ist „not measured and not promised“, mit dem, was dort ohne Konto gemessen ist (Schlüsselpaar im Volume: Validator z6; Auflösung von `ollama.com` nach Eintrag, `generate` und `up --recreate`: a2), und ausdrücklich nicht dem Meldungswechsel von `ollama run`, keiner Anmeldung, keinem Cloud-Aufruf; der Relay in einem Satz (Name, nicht Dienst; Erreichbarkeit nicht getestet; dahinter geht es nicht, gemessen als Z12 in [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md)); Grenze der zwei Zusage-Mengen (Carveout 2). Neuanlage mit demselben Volume „not promised“ (Lastenheft Zeile 414, Slice benennt die Fundstelle) |
| A-2 (Satz zu `ANTHROPIC_API_KEY`) | **erfüllt** | Handbuch, Spiegelstrich unter *Ollama (optional)*: Durchreichung „abgeleitet aus der Referenz, nicht mit einem Schlüssel gemessen“; gemessen b1 (ungesetzt kommt ungesetzt an) und b2 (leer gesetzt kommt als leere, gesetzte Variable an), je ein Lauf `u-boot devcontainer shell`, Linux `amd64` (Validator-Prüfpunkte b1, b2 gelesen und gegen den Wortlaut gelegt); `unset`/`env -u` „nicht gemessen“, nicht empfohlen. Keine Wirkung versprochen |
| A-3 (Satz zu Rückbau des Geräteschlüssels) | **erfüllt** | Absatz *Rückbau*: Seite des Kontos ohne Adresse; `id_ed25519.pub` im Volume (z6, 81 Byte), Anzeige auf der Seite „nicht gemessen“; `down` entfernt weder Container noch Volume (`down.log`: „not removed“, Volume bleibt); `docker volume rm` abgelehnt mit „volume is in use“ bei gestopptem Container (`rm1.log`, Exit 1), nach `docker rm` entfernt (`rm2.log`, Exit 0); kein Satz zu `ollama signout`. „Mit dem Volume gehen die geladenen Modelle verloren“ ist aus dem Handbuch-Punkt zum Modell-Speicher gelesen, nicht eigens gemessen (Slice §3 Punkt 4 sagt es so). Nicht gefahren und im Satz nicht genannt: der Weg über `up --recreate` |
| V-1 (Widerspruch im Satz zur ungesetzten Host-Variable) | **entschieden: eingegrenzt** | Handbuch *Claude Code*, Zugangsdaten: „nicht gesetzt“ (b1, `u-boot devcontainer shell`) und „Not logged in“ (Dev-Container-CLI, Colima `arm64`, Messung im Slice [`slice-v1-devcontainer-feature-claude-code`](../plan/planning/done/slice-v1-devcontainer-feature-claude-code.md)), mit „mit `u-boot devcontainer shell` nicht gefahren“; keine neue Zusage |
| I-3 (Exit `10` am Binary im Lead) | **erfüllt** | Lead: ein Lauf, Linux `amd64`, `ollama-host.version` und `ollama.version` ohne `source` Exit `10`, `u-boot.yaml` unverändert, `--dry-run` ebenfalls; Rohlog `cut16/mv/i3.log` gelesen: Exit 10 viermal, „u-boot.yaml unveraendert“. Der Pfad mit `source` ist „not run“. Eintrag `[0.15.1]` unverändert |

## Prüfungen im Einzelnen

* **Lead gegen Zusage-Zeilen und Messung:** Jeder Satz des Leads gegen Validator-Bericht, Messprotokoll mit Konto (Messstand: Ubuntu 24.04.5, Docker 29.9.2, lokale Engine; „two models, one to three runs, hand variants“), [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) (Z11, Z12) und Lastenheft Zeile 414 gelegt. Nichts geht weiter als das Gemessene; kein CI-Satz (M2 steht aus); kein Satz zu `claude -p` im Lead.
* **Einstufung:** *Added* (ein Eintrag) und *Changed* (zwei Einträge) unverändert im Wortlaut; kein `### Security`, kein Advisory (M5b entfällt), wie im Plan und im Beschluss des Starts (Frage 5). `[Unreleased]` ist leer. „Sixteenth release“ passt zur Reihe (`[0.15.0]` trägt „Fifteenth release“).
* **Version an allen Stellen:** Suche nach `0.15.0` und `0.15.1` außerhalb von `plan/`, `reviews/`, `messungen/`, `baseline/` und CHANGELOG: Treffer nur in den Tabellenzeilen `v0.15.0`/`v0.15.1` beider READMEs (Geschichte), `docs/maintainer/releasing.md:57` (Zeiger „zuletzt“, Closure M8) und zwei Zeilen im Änderungsverlauf von Lastenheft und Spezifikation (Stand 0.15.0 des Lastenhefts, nicht berührt). Kein `0.15.1-dev` im Baum.
* **Datum (Regel M0):** `date +%F` zum Prüfzeitpunkt = `2026-10-25` (Ortszeit, 00:30 CEST); CHANGELOG, Tabellenzeile beider READMEs und Handbuch-`Stand` tragen `2026-10-25`. **Keine Abweichung jetzt.** Der Cut begann am 2026-10-24 und lief über Mitternacht; die Datumsstellen sind erst vor dem Cut-Commit umgestellt, die Läufe `make gates` (erster Lauf) und M-V liegen davor, geändert sind nur Datumszeichen. **Bedingung:** Push (M1) und Tag (M5) müssen am 2026-10-25 geschehen; ab 2026-10-26 ist ein weiterer Cut-Commit mit Datum `date +%F` nötig und M2 auf dem neuen Push-Tip erneut.
* **M-T:** Nachgefahren (siehe Sensoren) mit identischem Ergebnis: Fill Exit 0 mit vier `v0.16.0`-URLs; `tap-formula-check.sh` Exit 0 mit `v0.16.0`, Exit 1 mit `v0.15.1` und bei gemischter Formel (3 von 4 Zeilen); `tap-nachzug.sh` mit Stellvertretern für `curl`/`git` protokolliert für `v0.16.0` `curl`, `clone`, `diff --quiet`, `add`, `commit -m u-boot 0.16.0`, `push`, mit `v0.15.1` Exit 1 nach dem `curl`-Aufruf, 0 Clone-Aufrufe. **Beweist nicht:** `is_stable`, echter `curl`, echter `git push`, Token; `tap-check` (M4) bleibt Sache des Hauptagenten.
* **M-V:** Rohlogs `up.log`, `down.log`, `rm1.log`, `rm2.log` gelesen (Reihenfolge: `up`, `down` „not removed“, `docker volume rm` Exit 1 „in use“, `docker rm`, `docker volume rm` Exit 0). Aufräumen belegt: `img-before`/`img-after` und `vol-before`/`vol-after` gleich (diff leer); aktuell 0 Volumes `vcut16-*`.
* **Plan-Tabelle gegen `git diff --stat 7956e157..HEAD`:** 12 Dateien: `CHANGELOG.md`, `Dockerfile`, `README.md`, `README.de.md`, `cmd/uboot/main.go`, `repo.mk`, `docs/user/benutzerhandbuch.md`, Roadmap, `open/README.md`, Welle-Datei, der Slice, der Review-Report. Jede Datei hat einen Plan-Eintrag; der Review-Report ist als `docs/reviews/<Datum>-slice-…` geplant; der Verification-Report (diese Datei) ebenfalls. „Kein Diff“-Zeilen (`releasing.md`, `carveouts.md`, `spec/`, ADRs, Workflows, `scripts/`, `packaging/`, `Makefile`, `go.mod`, `tools/`, `internal/`) im Diff nicht vorhanden. Der Fix `c2cf2202` berührt nur die schon geplante Welle-Datei. Nach dieser Runde kommt dieser Report hinzu.
* **Zustandsmarker (`AGENTS.md` §3.10):** Offen sind nur Marker mit benanntem Ablaufpunkt: Roadmap „(Push und Tag stehen aus)“, Slice-Kopf „Stand vor M1: nichts gepusht, kein Tag“, Welle-Datei „der Release-Cut steht aus“. Bei M8/Closure zu entfernen oder auf den Ist-Stand zu setzen. Die Aussage in §3 „`origin/main` unverändert“ gilt bis M1.
* **Kommentare im Diff (§3.7, §3.9):** `repo.mk` und `Dockerfile` nennen den geltenden Wert; kein Kommentar beschreibt Entferntes.

## Was für M1 bis M8 bleibt (Hauptagent und Projektinhaber)

| # | Aktion | Zu prüfen / Hinweis |
|---|---|---|
| vor M1 | `date +%F` = `2026-10-25`; Lead im Wortlaut dem Projektinhaber vorlegen (Frage 1; INFO-1 als Gelegenheit); `make docs-check` lokal; Arbeitsbaum sauber (nach diesem Report-Commit) | Freigabe des Projektinhabers im Gespräch; ein Push mit allen Commits des Cuts, beide Move-Commits davor |
| M1 | `git push origin main` | Hauptagent, nach Freigabe; kein `--tags` |
| M2 | `ci` und `integration-docker` auf dem Push-Tip | Job-Log per REST-API und Check-Run-Annotationen lesen. `ci` `success` (Jobs `gates`, `security-gates`, `image-scan`), keine Annotation „Unexpected input(s)“. `integration-docker`: `--- FAIL` 0, kein Timeout; **`--- SKIP` zählen, erwartet die vier Nested-Podman-Tests** (`TestE2E_SandboxDevcontainer_NestedPodman_Build`, `…_NestedPodman_BuildEgress`, `…_NestedPodman`, `…_StatusFiles_VFS`); die zwei Tests `TestE2E_OllamaHost_DefaultProfileResolvesHostDockerInternal` und `TestE2E_OllamaHost_SandboxWithEgressResolvesAndAddsNoHost` **sichtbar `--- PASS`**, nicht übersprungen; Dauer gegen `timeout-minutes`; Abschnitt „Operating System“ der Logs. Das ersetzt den vollen `make test-docker` des DoD Punkt 3 |
| M3 | Vorab-Tag | entfällt (kein `-rc`) |
| M4 | `gh workflow run tap-check` unmittelbar vor M5 | grün = Token darf auf das Tap schreiben |
| M5 | `git tag v0.16.0` auf dem in M2 geprüften Commit, `git push origin v0.16.0` | nur nach ausdrücklicher Freigabe im Gespräch (eigene Rückfrage, nicht die von M1); `date +%F` erneut; der Tag wird nicht verschoben |
| M5b | Advisory | entfällt |
| M6 | `gh run watch`; `gh release view v0.16.0 --json assets,isPrerelease,isDraft` (zwölf Assets, kein Prerelease, kein Draft); `docker run --rm ghcr.io/pt9912/u-boot:0.16.0 --version` = `0.16.0`; Digest `:0.16.0` = `:latest`; Job `tap` (Log „updated to v0.16.0“ oder „already at“, keine Warnung zum Secret, Annotationen); Paket-Smoke rpm/fedora und deb/ubuntu `success`; Release-Notes verweisen auf `CHANGELOG.md` §[0.16.0] | Lauf-ID in §7 des Slice |
| M7 | Mac: `brew tap pt9912/u-boot`, `brew trust pt9912/u-boot`, `brew install u-boot`, `u-boot --version` = `0.16.0` | Projektinhaber; steht seit `v0.13.1` unbestätigt aus; ersetzt keine Zusage für macOS |
| M8 | Nacharbeit nach [`releasing.md`](../maintainer/releasing.md) §7: Roadmap-Zeilen mit Tag-Hash und Datum, Zeiger „zuletzt“ in `releasing.md`, Datum gegen `git log -1 --format=%cs v0.16.0`; Marker entfernen; Closure-Notiz, Register-Evidence (M-1), Entscheidung über Folge-Slices (L-1, L-3, I-1), `git mv` nach `done/`; danach `close-welle` | Closure erst nach M6; `make gates` danach erneut |

## Verification Evidence (Block nach `harness/verification.md`)

Scope:
- Slice: [`slice-v1-release-cut-v0.16.0`](../plan/planning/done/slice-v1-release-cut-v0.16.0.md)
- IDs: [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md)
- Artefakte: `CHANGELOG.md`, `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md`, `cmd/uboot/main.go`, `repo.mk`, `Dockerfile`

DoD-Abgleich: siehe Tabelle oben; Punkte 1, 2, 4, 5 erfüllt, Punkt 3 (Tag, M2, M4, M5, M6) und Punkte 8 bis 11 offen mit Grund (Reihenfolge der Maintainer-Aktionen, Closure).

Carveouts:
- Neu: none
- Gelöst: none
- Unverändert: Carveout 2 (zwei Zusage-Mengen für `ollama`, Anker [`slice-cr-lastenheft-messstaende`](../plan/planning/open/slice-cr-lastenheft-messstaende.md)), im Lead als Grenze benannt

Nicht ausgeführt:
- `make test-docker`, `make ci`, `make fullbuild` — Auftrag/Speicher; Produktcode im Diff nur `var version`; Implementer-Belege gelesen. Ersatz: M2 (CI auf dem Push-Tip) mit Zählung der SKIPs
- `make verify-depguard` — keine Importregel berührt
- Messung mit Konto am Pfad `u-boot devcontainer up` — ausdrücklich nicht Teil des Cuts (Frage 2)

Commit / Artefakt:
- Geprüfter Stand `c2cf2202`; Tag `v0.16.0` nicht gesetzt
