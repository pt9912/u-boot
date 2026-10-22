# Review-Report: [`slice-v1-release-cut-v0.13.1`](../plan/planning/in-progress/slice-v1-release-cut-v0.13.1.md) — Nachlauf-Review der Fix-Runde nach dem Review — 2026-10-22

**Review-Art:** Code/Doku, Nachlauf (klein) — die Release-Texte und der Advisory-Entwurf der Fix-Runde gegen Plan, Entscheidungen 27 bis 32, [`releasing.md`](../maintainer/releasing.md) §3 und die Hard Rules (`AGENTS.md` §3.7, §3.9, §3.10, §5 Regel 6). Wartbarkeit und Wahrheit des Textes, keine DoD-Closure (Verifier). Das erste Review ([`2026-10-22-slice-v1-release-cut-v0.13.1`](2026-10-22-slice-v1-release-cut-v0.13.1.md): F-1 bis F-8) ist Kontext und nicht erneut geprüft; geprüft ist, ob F-1 bis F-5 im Diff behoben sind und ob der Diff Neues einführt. Keine Einschätzung des Implementers ist ungeprüft übernommen.

**Gegenstand:** `git diff 44a65e31..68e6b83a -- CHANGELOG.md README.md README.de.md docs/user/ docs/maintainer/` (3 Dateien; `docs/user/` und `docs/maintainer/` ohne Diff), der Plan-Teil derselben Range (`git diff -U0`, 72 Zeilen: Liefer-Punkt *Security und Advisory*, Plan-Tabelle, B20, B23 bis B26, M2, M5b, M6, Absatz *Advisory-Entwurf*, Abschnitt *Fix-Runde nach dem Review*, §6, Entscheidungen 27 bis 32) und der Advisory-Entwurf im Scratchpad (`advisory-v0.13.1/`: `payload.json`, `payload-other.json`, `description.md`, `README.txt`; nicht im Repo). Zwei Commits: `d4f73682` (Review-Report, unverändert aufgenommen) und `68e6b83a` (Fix-Runde). `HEAD` = `68e6b83a`, Arbeitsbaum zu Beginn sauber, nichts gepusht, kein Tag.

**Skill:** `.harness/skills/reviewer.md` @ `68e6b83a` · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-22

**Eingangs-Kontext:**

- Slice [`slice-v1-release-cut-v0.13.1`](../plan/planning/in-progress/slice-v1-release-cut-v0.13.1.md) §3 *Fix-Runde nach dem Review*, §6, §7 Entscheidungen 27 bis 32
- [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes); [`ADR-0032`](../plan/adr/0032-katalog-features-socket-und-lockerungen.md) (Entscheidung 6, `Accepted`, ohne Diff)
- [`releasing.md`](../maintainer/releasing.md) §3, §6, §8; `AGENTS.md` §3.7, §3.9, §3.10, §5 Regel 6
- Vorheriges Review am Bereich: [`2026-10-22-slice-v1-release-cut-v0.13.1`](2026-10-22-slice-v1-release-cut-v0.13.1.md)

**Selbst gelesen oder nachgefahren (nicht aus dem Slice übernommen):**

- OpenAPI-Beschreibung der GitHub-REST-API (`info.version` 1.1.4, `github/rest-api-description`, `descriptions/api.github.com/api.github.com.yaml`, lesend per `gh api` GET am 2026-10-22 in ein neues Scratchpad-Verzeichnis, nur mit `grep` und `sed` gelesen): Schemata `repository-advisory-create`, `repository-advisory-update`, `security-advisory-ecosystems` und die Operation `PATCH /repos/{owner}/{repo}/security-advisories/{ghsa_id}` mit ihren Beispielen und Antworten (siehe NL-1).
- Entwurf in Docker (`alpine:3` mit `jq` 1.8.2): beide JSON-Dateien gültig; `summary` 121 Zeichen, `description` 4021 Zeichen; Schlüssel `cvss_vector_string`, `cwe_ids`, `description`, `summary`, `vulnerabilities`; kein `severity`, kein `cve_id`; `cwe_ids` `["CWE-668"]`; `.description` von `payload.json` gleich `description.md` (Diff leer); `payload-other.json` unterscheidet sich nur in `package.ecosystem` (`other`). Satzweiser Vergleich von `description.md` mit den Zeilen 24 bis 28 von `CHANGELOG.md`: die Abweichungen sind Gliederung, die Spanne in der Schreibweise der API und „before the fix“ statt „before this release“; kein Satz weiter als der Eintrag.
- `gh run list --workflow publish` und `gh run view 39567098112 --json conclusion,jobs` (lesend, Lauf von `v0.13.0`): die vier Jobnamen aus `README.txt` und Plan stimmen wörtlich mit `.github/workflows/publish.yml` (`name:` Zeilen 32, 268, 288, 313) und mit dem Lauf überein; `gh release view v0.13.0 --json assets,isPrerelease,isDraft`: 12 Assets, kein Prerelease, kein Draft. Die Läufe von `publish` zu `v0.11.0`, `v0.12.0` und `v0.13.0` dauerten etwa zwei Minuten (`startedAt` bis `updatedAt`).
- `scripts/tap-nachzug.sh` (Meldungen „updated to“ und „already at“), `Dockerfile.tmpl` Zeile 3 (`FROM … base:debian` ohne Profilbedingung), `internal/e2e/devcontainer_start_docker_test.go` (Signatur `dcGenerateProject` Zeile 191, `Sandbox: sandbox` Zeile 200, Aufruf mit `false` Zeile 464, Kommentar Zeilen 446 bis 460), README `nestedRuntime: docker` (Status-Block), Fix-Slice §3 *Umsetzung* 1, Handbuch Zeilen 738, 779, 975, 1183, 2226.
- Alle hinzugefügten Zeilen der READMEs und `git log 44a65e31..68e6b83a --format=%B` durchsucht (siehe Prüfpunkt 6).

**Sensor- und Abfrage-Läufe des Reviewers** (Docker-only; Exit-Code des `make`-Prozesses):

| Lauf | Ergebnis |
|---|---|
| `make docs-check` mit diesem Report | Exit 0, 651 Dateien, 0 Befunde (Abschnitt *Abschluss*) |
| Gelesen, nicht gefahren | `make gates`, `make test-docker`, `make govulncheck`, `make image-scan`: die Runde ändert nur Texte (Abschnitt *Nicht geprüft*) |

---

## Findings

Zuerst der Befund zu einer Unwahrheit (im Plan und im Entwurfs-`README.txt`, nicht im Release-Text); **eine Übertreibung in CHANGELOG, READMEs oder Advisory-Entwurf habe ich nicht gefunden**: jede Kurzform sagt, was der Eintrag *Security* auch sagt (Prüfpunkt 1). Die drei INFO sind Beobachtungen ohne Änderungspflicht.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| NL-1 | LOW | Plan und Entwurf führen als Lesung, die Regel `published` → `draft` stehe in der OpenAPI-Beschreibung 1.1.4 nicht („nicht zu finden“, „keine Regel zum Wechsel“). Sie steht dort: Die Operation `PATCH …/security-advisories/{ghsa_id}` führt das Beispiel `invalid_state_transition` („Example of an invalid state transition, from `published` to `draft`.“) und die 422-Antwort „Invalid state transition from `published` to `draft`.“; „nur `closed`“ steht dort nicht. | [`AGENTS.md`](../../AGENTS.md) §3.10 (Aussage über eine gelesene Quelle); Plan Entscheidung 27 und B26 | Plan · „die Regel `published` → `draft` ist in der OpenAPI-Beschreibung 1.1.4 nicht zu finden“ (B26); Plan · „der Implementer fand die Regel in der OpenAPI-Beschreibung 1.1.4 nicht“ (§6); Plan · „steht in der OpenAPI-Beschreibung 1.1.4 nicht (gelesen“ (Fix-Runde, Zeile F-3); Entwurf `README.txt` · „keine Regel zum Wechsel“ | Review-only (`gh api` GET der Beschreibung, Suche nach `invalid_state_transition`) | Negativaussage über eine gelesene Quelle unvollständig gelesen |
| NL-2 | INFO | Betroffenheit in den Kurzformen: (a) *Fixed* trägt weiter „Affected: `v0.6.0` to `v0.13.0`“ ohne den Vorbehalt, den README-Tabellenzeile und Advisory-Entwurf in dieser Runde bekommen haben; (b) „whether a project is affected depends on the base image“ nennt einen von zwei Faktoren, denn mit von Hand gesetztem `moby: false` baut `docker-cli` auch auf Debian 13 und trägt den Mount (gemessen), die Standard-Option steht nur im Eintrag *Security* („the build with the default option“). | Plan Entscheidung 28; [`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](../plan/planning/observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md) | `CHANGELOG.md` · „Affected: `v0.6.0` to `v0.13.0`“; `README.md` · „whether a project is affected depends on the base image“ | Review-only | Kurzform nennt einen von zwei Faktoren der Betroffenheit |
| NL-3 | INFO | Das Handbuch sagt für das Standardprofil „bindet … ein“ (Zeilen 779, 975; „wenn ein Feature es tut“ Zeilen 738, 1183, 2226) und ist damit fester als „declares … according to its metadata“ der Kurzformen. Die Aussage ist durch den Docker-Test im Standardprofil gedeckt (Mount an einem gestarteten Container gemessen, Variante `moby: false`); kein Widerspruch (Bewertung unter Prüfpunkt 1). | [`AGENTS.md`](../../AGENTS.md) §3.10 (Bedingung in Spec oder Handbuch braucht einen Test); Plan N4 („das Handbuch bekommt keine Bedingung über einen beweglichen Tag ohne Test“) | `docs/user/benutzerhandbuch.md` · „Das Feature `docker-cli` bindet im Standardprofil das Docker-Socket des Hosts in den“ | Review-only; `make test-docker` (Test `TestDevcontainerUp_Docker_StandardProfileDockerCLIMountsTheSocket`) | Handbuch fester als Kurzform, durch Test gedeckt |
| NL-4 | INFO | Das Risiko „Advisory öffentlich, `publish` scheitert“ sagt, der Eintrag *Security* „nennt die Abhilfe, die ohne Update wirkt“. Gemessen ist, dass die Schritte Abschalten und Erzeugen mit `v0.13.0` enden (Exit 0, Schlüssel entfernt); „wirkt“ reicht bis zum neu angelegten Container ohne Socket, und der ist laut *Security* nicht gemessen. | Plan B25 (b); [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md) | Plan · „nennt die Abhilfe, die ohne Update wirkt“ (§6) | Review-only | Plan-Satz weiter als die Messung |

## Zu den Findings (Risiko, Failure-Szenario)

**NL-1 (LOW).** *Szenario:* Der Verifier liest B26 („nicht zu finden, gelesen am 2026-10-22“) gegen die Quelle und findet die Regel; der Hauptagent liest in `README.txt`, die Regel sei eine Annahme ohne Quelle. *Einordnung:* Die Aussage ist nirgends als Tatsache gesetzt, der Release-Text und der Advisory-Text (`payload*.json`, `description.md`) kommen darauf nicht zu sprechen (Suche nach „draft“ in `CHANGELOG.md`, READMEs und `description.md`: 0 Treffer); die Vorsorge (M5b nach den Kernprüfungen) bleibt richtig und ist besser gestützt als geschrieben. Falsch ist die Lesung, nicht die Entscheidung. Das erste Review hatte die Regel richtig zitiert („führt `published` → `draft` als ungültigen Zustandswechsel“); die Fix-Runde hat es als „Grenze“ übernommen, ohne die Beispiele der Operation zu lesen. Unbelegt bleibt „nur `closed`“ (die Beschreibung nennt den Wechsel nach `closed` nicht). Berichtigung betrifft Plan und `README.txt`, nicht CHANGELOG, READMEs oder den Advisory-Text; ein weiteres Nachlauf-Review nach Plan §6 löst sie nicht aus.

**NL-2 (INFO).** Kein Fehler im Diff: Der Eintrag *Security* trägt dieselbe Überschrift-Zeile und die Einzelheiten darunter; die READMEs verweisen auf ihn („what depends on the base image“). Zu (a): Der erste Review hat die Zeile als Versionsspanne des Fehlers für wahr befunden; die Aussage reicht nicht weiter als „From `v0.6.0` to `v0.13.0` the generator rendered the feature in both profiles without a warning“. Zu (b): `Dockerfile.tmpl` Zeile 3 trägt `FROM … base:debian` ohne Profilbedingung, die Abhängigkeit gilt für beide Profile; gemessen ist sie nur im Standardprofil, und das steht im Eintrag.

**NL-3 (INFO), Bewertung und Empfehlung zu Prüfpunkt 1.** *Kein Widerspruch:* (1) Die Handbuchzeilen 779 und 1183 sprechen vom Standardprofil, in dem der Docker-Test den Mount an einem gestarteten Container misst (Linux `amd64`, Docker 29.9.2; gelesen am Test: `dcGenerateProject(ctx, t, false, …)`, `docker inspect`, Variante `moby: false` von Hand); das ist die Test-Deckung, die [`AGENTS.md`](../../AGENTS.md) §3.10 für eine Bedingung im Handbuch verlangt. (2) Für das Sandbox-Profil sagt das Handbuch „wird abgewiesen“ (Zeilen 781, 1183), nicht „bindet“; Zeile 975 beschreibt das Feature allgemein und nennt im selben Satz „Im Sandbox-Profil nicht wählbar“. (3) Dass das Handbuch die Basis-Image-Bedingung nicht trägt, ist Plan-Entscheidung (N4, offener Slice). *Rest:* Ein Leser im Standardprofil auf dem Standard-Basis-Image sieht den Bauabbruch, nicht den Mount; „bindet … ein“ beschreibt das Feature und ist dort nicht falsch, aber nicht beobachtbar. *Empfehlung:* nicht ändern. Eine Änderung wäre ein weiterer Textcommit an einem gedeckten Nutzertext (erneut M0, M2, Nachlauf), und die Bedingung gehört in den Vorgang von N4, sobald der Projektinhaber den Weg entschieden hat. Der Plan führt dieselbe Zuordnung in seinem Selbstcheck („nicht Gegenstand dieser Runde“); die Begründung „vom Review nicht beanstandet“ dort ersetzt keine Lesung — hier ist sie nachgeholt.

**NL-4 (INFO).** Nur im Plan (§6), nicht im Release-Text; dort steht richtig „need no update“ und „Not measured: that a container recreated after these steps starts without the socket“. Keine Änderungspflicht.

## Zu den sechs Prüfpunkten

**1. F-1, Kurzformen gegen den Eintrag *Security*.**

| Stelle | Kurzform | Eintrag *Security* | Ergebnis |
|---|---|---|---|
| `README.md` Status | „the built-in feature `docker-cli` declares such a mount according to its metadata“ | „declares such a mount“; Lead „mounts it according to its metadata“ | gleich |
| `README.de.md` Status | „deklariert laut seinen Metadaten einen solchen Mount“ | wie oben | gleich, inhaltsgleich zu en |
| Tabellenzeile `v0.13.1` (en, de) | „declares … (`v0.6.0` to `v0.13.0` rendered it in both profiles without a warning; whether a project is affected depends on the base image)“ | „From `v0.6.0` to `v0.13.0` the generator rendered the feature in both profiles without a warning (measured with `v0.13.0`, read in the code of the earlier versions)“; „Whether a project is affected depends on the base image“ | nicht firmer; der Klammerzusatz „measured … read …“ fehlt, die Zeile behauptet nichts darüber hinaus |
| *Fixed* | „declares such a mount in its metadata“ | wie oben | gleich |
| Advisory `summary` | „declares a mount of the container-engine socket of the host“ | wie oben | gleich |
| Advisory `description` | satzweise gegen die Zeilen 24 bis 28 verglichen | — | kein Satz weiter; „Affected“ trägt den Vorbehalt sogar ausdrücklich |
| „(read, not measured)“ beim vorhandenen Container (Status en/de, Tabellenzeile en/de) | „keeps the socket until it is recreated (read, not measured)“ | „(read, not measured: the mounts of a container are fixed when it is created)“ | korrekt; Quelle [`ADR-0032`](../plan/adr/0032-katalog-features-socket-und-lockerungen.md) Entscheidung 6 (gelesen, nicht gemessen) |

„whether a project is affected depends on the base image“: **wahr** und mit Datum und Umgebung belegt — im Eintrag *Security* (2026-10-22, Linux `amd64`, Docker 29.9.2, Debian 13 `trixie` am Tag des Tests, zwei Läufe im Standardprofil, Meldung „The 'moby' option is not supported on debian 'trixie'“, Fix-Slice §3 *Umsetzung* 1 und Plan §3 *Gegenprobe zu N4*); die Kurzformen tragen Datum und Umgebung nicht und verweisen auf den Eintrag. Die andere Seite der Abhängigkeit (ein anderes Basis-Image trägt den Mount) ist im Eintrag als „read, not measured there“ ausgewiesen. Restpunkt: NL-2. Handbuch: NL-3.

**2. F-2, Status-Block.** Der erste Absatz beider READMEs steht im Präsens und nennt Zustand und Zusage („promises“, „declares“, „therefore rejects“, „keeps“); er enthält weder „released“ noch „The release is“ noch „now“/„nun“ noch ein Datum noch eine Kennung (Suche über alle hinzugefügten README-Zeilen: Treffer nur in den Tabellenzeilen — Datum `2026-10-22` und „now“/„nun“, beides dort zulässig). „created with the feature earlier“ / „früher mit dem Feature angelegter Container“ beschreibt einen Zustand (vorhandene Container), keine Chronik des Releases. en und de sind inhaltsgleich (Reihenfolge der Aussagen und Klammern gleich). [`releasing.md`](../maintainer/releasing.md) §3 Nr. 3 erfüllt.

**3. F-3 und F-5.** M-Tabelle: M6 steht vor M5b, die Kernprüfungen (Lauf `publish` `success` mit den vier Jobs, `tap` wirklich gelaufen, Release ohne Prerelease/Draft mit zwölf Assets, `--version` des Images) stehen in M6 und bedingen M5b; DoD-Punkt und Entscheidung 27 sind gleich. Das §6-Risiko „Advisory öffentlich, `publish` scheitert“ trägt Szenario, Vorsorge (M5b nach den Kernprüfungen), Ausgang (`releasing.md` §8: Lauf neu starten, `--clobber`, Tag nicht verschieben, oder nächster Patch-Release) und Meldung an den Projektinhaber; das Fenster-Risiko nennt das längere Fenster. `README.txt`: Bedingungen (a) bis (d) stimmen mit M6 überein; die Jobnamen stimmen mit `publish.yml` überein, die Log-Zeilen mit `scripts/tap-nachzug.sh` (Zeilen 30 und 38) und [`releasing.md`](../maintainer/releasing.md) §6; die Abfragen (`gh run list --workflow publish`, `gh run view … --json conclusion,jobs --jq`, `gh release view … --json assets,isPrerelease,isDraft`) habe ich lesend gegen den Lauf von `v0.13.0` gefahren: gleiche Felder, vier Jobs `success`, zwölf Assets. „Etwa drei bis fünf Minuten“ ist eine Schätzung ohne Beleg-Zeile; die letzten Läufe von `publish` dauerten etwa zwei Minuten, M5b hängt an den Bedingungen, nicht an der Zeit. „`published` → `draft` ungültig“ steht als Annahme (B26, §6, M5b, Entscheidung 27, `README.txt`) und nirgends als Tatsache — die Begründung dazu ist falsch gelesen (NL-1). Messgrenzen: (a) die zwei N4-Läufe im Standardprofil — wahr: der Docker-Test ruft `dcGenerateProject(ctx, t, false, …)` (`sandbox` = `false`, Zeile 464, Signatur Zeile 191), der erste Lauf ist der Test selbst (Fix-Slice §3 *Umsetzung* 1: „erzeugtes Projekt (`dcGenerateProject`, Feature `docker-cli`)“), der zweite `init` ohne `--sandbox` (Plan §3 *Gegenprobe zu N4*: „kein `profile`“); (b) Ende der Abhilfe nicht gemessen — wahr: die Wirkungsprobe endet nach `generate` (Plan §3, Grenze); (c) `nestedRuntime docker` ungeprobt, auf Docker Desktop nicht zugesagt — wahr: README-Status „`nestedRuntime: docker` is promised on Docker under Linux and on Colima, not on Docker Desktop“. Alle drei stehen in *Security* und im Advisory-Entwurf, (a) inline, (b) und (c) unter **Not measured**.

**4. F-4, B20.** B20 stimmt mit dem Reviewer-Lauf und mit der *Abweichung* 1 überein (`govulncheck` mit Go 1.27.2: 0 Funde in importierten Paketen, 1 in benötigten Modulen, „doesn't appear to call“; die Ausgabe benennt den Fund ohne `-show verbose` nicht, also ist er nicht als `CVE-2026-111943` belegt). Der Spiegelstrich in *Security* sagt „fixes one HIGH stdlib advisory (`CVE-2026-111943`, …) that failed the Trivy image scan“ und nichts zur Erreichbarkeit; die Tabellenzeile der READMEs ebenso („a HIGH stdlib advisory failed the image scan“); das Advisory nennt den Pin nicht (Suche nach „1.27“ in `description.md`: 0 Treffer). Die Zeile `make ci` (`5040b1f1`) in der Sensor-Tabelle des Plans führt nur „No vulnerabilities found“, die Zeile vollständig steht in B20 und *Abweichung* 1 — gekürzt, nicht widersprüchlich.

**5. Advisory-Entwurf.** Siehe *Selbst gelesen*: gültiges JSON (Docker), `summary` 121 von 1024 Zeichen, `description` 4021 von 65535, kein `severity` neben `cvss_vector_string`, kein `cve_id`, `cwe_ids` `["CWE-668"]`, alle Schlüssel im Schema definiert (`additionalProperties: false`), `ecosystem` `go` und `other` in der Aufzählung, `package.name` `github.com/pt9912/u-boot` = `go.mod` Zeile 1. CVSS `CVSS:3.1/AV:L/AC:H/PR:L/UI:N/S:C/C:H/I:H/A:H` = 7,8 nachgerechnet (Exploitability 1,1498, Impact 6,0474, Scope geändert). Text folgt *Security*, kein Satz weiter. Unverändert gegenüber dem ersten Review: der schreibende Aufruf ist nicht erprobt.

**6. Kennungen und Messages.** Hinzugefügte README-Zeilen: keine `LH-`, `ADR-`, `SPEC-`, `ARC-`, Slice-, Wellen-, `BEO-`, `CVE-`, `GHSA`-, `CWE`-Kennung, kein Verweis auf Planung oder Spezifikation. `git log 44a65e31..68e6b83a --format=%B`: 0 Treffer für `SPEC-` und `ARC-`; 0 Treffer (ohne Beachtung der Schreibweise) für Security, Advisory, CWE, GHSA, CVE, `docker-cli`, Socket, Mount, Sandbox, `moby`. Die Messages nennen den Fehler nicht in Worten; „Betroffenheit haengt vom Basis-Image ab“ und „Veroeffentlichung am Tag“ stehen nicht auf der Liste der Entscheidung 31 und benennen weder Bauteil noch Profil noch Socket — im Rahmen der Entscheidung, das Fenster ist ohnehin der Push.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `CHANGELOG.md` (Diff: *Security* Zeilen 26 und 27, *Fixed* Zeile 33) | geprüft; Befund nur NL-2 (a), sonst ohne Befund (Kurzform gleich Eintrag, „both in the standard profile“ und **Not measured** wahr und belegt, B23 und B25) |
| `README.md`, `README.de.md` (Status-Block, Tabellenzeile `v0.13.1`) | geprüft, ohne Befund (F-1 und F-2 behoben; en und de inhaltsgleich; keine Kennung) |
| `docs/user/`, `docs/maintainer/` | ohne Diff; Handbuchzeilen zu `docker-cli` gelesen — Bewertung NL-3 |
| Advisory-Entwurf (`payload.json`, `payload-other.json`, `description.md`) | geprüft, ohne Befund |
| Advisory-Entwurf `README.txt` | geprüft; Befund NL-1. Hinweis ohne Finding: der letzte Absatz („Dieser Entwurf aendert sich nicht mehr, wenn der CHANGELOG-Eintrag sich aendert - dann wird er … nachgezogen“) widerspricht sich im Wortlaut (gemeint: nicht von selbst); keine Folge für die Aufrufe |
| Plan: M-Tabelle (M2, M5b, M6), DoD-Punkt *Security und Advisory*, §6 (zwei Risiken), Entscheidungen 27 bis 32, Plan-Tabelle (14 Dateien gegen `git diff --stat 393fdb87..68e6b83a`, `cmd/internal` 1 Datei) | geprüft; Befunde NL-1 und NL-4; Plan-Tabelle gleich dem Diff |
| Beleg-Tabelle B20, B23 bis B26 | geprüft; B20, B23, B24, B25 ohne Befund; B26 siehe NL-1 |
| Commit-Messages `d4f73682`, `68e6b83a` | geprüft, ohne Befund (Prüfpunkt 6) |
| Kommentar-Selbstcheck [`AGENTS.md`](../../AGENTS.md) §3.9: Diff berührt keinen Code und keinen Kommentar | geprüft, ohne Befund |
| Dual-Classifier, Exit-Codes, Architekturgrenzen, Sprachvertrag, Dateisicherheit, ADR-Immutabilität (`docs/plan/adr/` ohne Diff), Referenzmodell | nicht berührt: Texte und Plan |

**Nicht geprüft:**

- `make gates`, `make test`, `make test-docker`, `make ci`, `make fullbuild`, `make verify-depguard` — die Runde ändert weder Code noch Pins noch Tests; die Läufe des Implementers (`gates` Exit 0 auf dem Stand der Runde) habe ich nicht wiederholt, ihre Wiederholung auf der Spitze gehört dem Hauptagenten und dem Verifier.
- Der schreibende Aufruf an die GitHub-API (Anlegen und Veröffentlichen des Advisory) — nicht ausgeführt, ohne Freigabe; ebenso die Wirkung von `state=published` und das Verhalten von `published` → `closed`.
- `govulncheck` mit Go 1.27.2 und `-show verbose`, `make image-scan` — Kontext des ersten Reviews, nicht wiederholt.
- Plan-Aussagen außerhalb des Diffs der Runde — nicht erneut gelesen.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Negativaussage über eine gelesene Quelle unvollständig gelesen · Kurzform nennt einen von zwei Faktoren der Betroffenheit · Handbuch fester als Kurzform, durch Test gedeckt · Plan-Satz weiter als die Messung

## Verdikt

**Merge-blockierend:** nein. Kein HIGH, kein MEDIUM; F-1 bis F-5 des ersten Reviews sind im Text behoben, keine Kurzform sagt mehr als der Eintrag *Security*, der Entwurf ist gültig und folgt ihm. **GO für den Verifier** auf `68e6b83a`.

**Übergabe:** NL-1 geht an den Implementer und den Hauptagenten: Die Berichtigung betrifft den Plan (B26, §6, Zeile F-3 der Fix-Runde) und `README.txt` des Entwurfs, nicht CHANGELOG, READMEs oder den Advisory-Text; sie sollte vor dem Verifier liegen, weil der Verifier B26 gegen die Quelle liest. NL-2 bis NL-4 sind Beobachtungen ohne Änderungspflicht; zu NL-3 empfehle ich, den Handbuchtext nicht zu ändern. Ein Finding ist nicht herabgestuft worden, weil der Implementer es anders sah; die gesetzten Entscheidungen 27 bis 32 sind unverändert übernommen. Die Finding-Klassen gehen in die Closure und von dort in den Zähler. Dieser Report ersetzt keine Verification — DoD und Wirkung prüft der Verifier getrennt.

## Abschluss

`make docs-check` mit diesem Report: Exit 0, 651 Dateien, 0 Befunde. Der Report ist nicht committet; `git status` zeigt außer ihm keine Änderung.
