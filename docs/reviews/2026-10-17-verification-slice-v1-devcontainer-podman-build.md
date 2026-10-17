# Verification-Bericht: slice-v1-devcontainer-podman-build — 2026-10-17 <!-- d-check:ignore (Slice-Kennung, Slice liegt in in-progress/ und wandert) -->

**Rolle:** Verifier (Modul 8/11) — Frage „Bauen wir es richtig?“ gegen Plan und DoD. Kein Review, keine Validierung, keine Reparatur.

**Gegenstand:** Commit-Range `5f372bc..HEAD` (17 Commits, Spitze `0588895`), Slice `slice-v1-devcontainer-podman-build` <!-- d-check:ignore (Slice-Kennung, Slice liegt in in-progress/) --> (Welle `welle-devcontainer-sshd-claude-code`, siebter Slice).

**Eingang:** DoD-Bestätigung und Sensor-Belege des Implementers (Slice §3 Messprotokoll Lauf 1 und 2, Plan-Tabelle), Review-Report [`2026-10-17-slice-v1-devcontainer-podman-build.md`](2026-10-17-slice-v1-devcontainer-podman-build.md).

## Verification Evidence

Scope:

- Slice: `slice-v1-devcontainer-podman-build` <!-- d-check:ignore (Slice-Kennung, Slice liegt in in-progress/) -->
- IDs: [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-007.b`](../../spec/spezifikation.md#lh-fa-dev-007b--ausstattung-bei-nestedruntime-podman), [`SPEC-019`](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), [ADR-0021](../plan/adr/0021-podman-build-im-sandbox-devcontainer.md), [ADR-0014](../plan/adr/0014-nested-podman-sandbox-devcontainer.md), [ADR-0012](../plan/adr/0012-devcontainer-egress-firewall.md)
- Artefakte: `internal/hexagon/application/devcontainer_sandbox.go` (`podmanRelaxations`), `devcontainer_egress.go` (`egressRegistryHosts`), `internal/e2e/sandbox_devcontainer_docker_test.go`, Lastenheft 0.12.0, Spezifikation, Handbuch, CHANGELOG, `harness/replay.md`, [`ADR-0021`](../plan/adr/0021-podman-build-im-sandbox-devcontainer.md)

### Sensoren (vom Verifier selbst gefahren)

| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass (Exit 0) | `coverage-gate: OK — coverage 91.20% meets threshold 90%`, lint, test, docs-check, baseline-verify grün |
| `make doc-trace` | pass | `83 Anforderung(en), 0 Waise(n).` |
| `make doc-immutable RANGE="5f372bc..HEAD"` | pass | `history-range-guard … 17 Commit(s) — OK`, `d-check: 310 Datei(en) geprüft, 0 Befund(e)`; `git diff --stat` nennt [`ADR-0014`](../plan/adr/0014-nested-podman-sandbox-devcontainer.md) nicht |
| `make test-docker` | pass (Exit 0) | Paket `internal/e2e` `ok 394.431s`; Colima-VM `kernel.apparmor_restrict_unprivileged_userns = 0` (per `colima ssh` geprüft) |
| Gezielter Nachlauf `-v` der drei NestedPodman-Tests (derselbe Container und Aufruf wie das Target, nur `-run`/`-v`) | pass, **kein Skip** | `…NestedPodman_Build` PASS 22,77 s, `…NestedPodman_BuildEgress` PASS 62,31 s, `…NestedPodman` PASS 18,72 s |
| `make docs-check` | siehe Abschnitt „Abschluss“ | nach Anlage dieser Datei gefahren |
| `make verify-depguard` | nicht gelaufen | keine Importregel und kein `.golangci.yml` berührt (Diff-Stat) |

Hinweis zum `make test-docker`-Beleg: das Target läuft ohne `-v`; ein Skip wäre dort unsichtbar. Deshalb der `-v`-Nachlauf (beide Build-Tests liefen mit realer Laufzeit, kein Skip).

## DoD-Abgleich

| DoD-Zeile | Befund | Evidence |
| --- | --- | --- |
| Liefer-Punkt 1 — Messprotokoll M0–M6, Matrix durchlaufen, Folge-ADR, [`ADR-0014`](../plan/adr/0014-nested-podman-sandbox-devcontainer.md) byte-gleich | erfüllt | Messprotokoll Lauf 1/2 im Slice §3 (Host, Engine, Arch, Kernel, `podman`-Version, Paketliste, sysctl, Digest); gewählte Matrix-Zeile benannt („fehlende Lockerung“, §3 „Entscheidung nach Matrix“); [`ADR-0021`](../plan/adr/0021-podman-build-im-sandbox-devcontainer.md) MADR, `Schärft:` [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), im ADR-Index; `doc-immutable` 0 Befunde. Lücken siehe D-1 (M4/M6-Teilmessung) und D-2 |
| Liefer-Punkt 2 — Umsetzung und Beleg | erfüllt | `podmanRelaxations()` mit sechs Optionen; `egressRegistryHosts()` mit `production.cloudfront.docker.com`; `Dockerfile.tmpl` laut Diff-Stat unberührt (Soll laut Plan: kein weiteres Paket); Tests `…_Build` (kalt, Cache, `--no-cache`, Alias) und `…_BuildEgress` laufen grün ohne Skip; Golden/Unit angepasst (6 bzw. 7 mit Egress); `harness/replay.md` Zeile `podman` nachgezogen; Spezifikation gegen Code gelesen |
| Liefer-Punkt 3 — Lastenheft 0.12.0, Handbuch, README, CHANGELOG | erfüllt | Kopf 0.12.0, §16-Historienzeile mit Messstand und Quelle, Zusage in [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer); Handbuch (Abschnitt „Verschachteltes Podman“, Egress-Tabelle, Grenzen); CHANGELOG zwei Einträge; README.md/README.de.md ohne Build-Aussage, unberührt (Plan: gelesen) |
| `make gates` grün, `make doc-trace` ohne Waise, `make test-docker` grün auf dem Mess-Host | erfüllt | siehe Sensoren |
| Review durchgeführt, Report unter `docs/reviews/` | erfüllt | Report vorhanden (8 Findings: 0 HIGH, 1 MEDIUM, 2 LOW, 5 INFO); Behebung siehe nächster Abschnitt |
| Doku-Update öffentlicher Vertrag | erfüllt | s. Liefer-Punkt 3 |
| Closure-Notiz mit Steering-Loop-Lerneintrag | **offen** | Slice §7 enthält nur „—“ (Closure-Vorgang, vor `git mv` nach `done/`) |
| Beobachtungs-Register fortgeschrieben oder „keine Beobachtung“ notiert | **offen** | §7 leer; Register unverändert (Diff-Stat) |
| Jedes Risiko aus §6 trägt Ausgang | teilweise | alle Punkte tragen einen Ausgang-Text, mehrere „offen“ (siehe Auflagen); Endstand erst bei Closure |
| Drei Paarungen (Anker · Folge-Slice · Register) | offen | Closure/Welle-Closure |
| DoD-Checkboxen im Slice §2 | **alle `[ ]`** | werden mit der Closure gesetzt; diese Datei setzt sie nicht |

## Review-Findings am Code und Text geprüft (nicht an der Prosa des Berichts)

| Finding | Soll | Befund am Artefakt |
| --- | --- | --- |
| F-1 MEDIUM `/dev/net/tun` Host-Vorbedingung | Vorbedingung im Lastenheft-Vertrag, Spezifikation, Handbuch, CHANGELOG | behoben: Lastenheft Host-Zusage („wenn das Device `/dev/net/tun` auf dem Container-Host vorhanden ist … ohne das Device startet der Container nicht“) und Degradations-Absatz („keine Degradationsstufe … `u-boot doctor` prüft das Device nicht“) im Text; Spezifikation [`LH-FA-DEV-007.b`](../../spec/spezifikation.md#lh-fa-dev-007b--ausstattung-bei-nestedruntime-podman) mit Fehlerbild; Handbuch mit Fehlerbild, Abhilfe (`modprobe tun`, als nicht gemessen gekennzeichnet) und Hinweis für Bestandsprojekte; CHANGELOG „Existing `nestedRuntime: podman` projects“. Code unverändert: keine `doctor`-Prüfung und keine Degradationsstufe — die Lücke ist ausgewiesen, nicht geschlossen (Slice §1 schließt es als neue Zusage aus; Slice §6) |
| F-2 LOW Egress-Zusage hängt am Blob-Host | Wortlaut „gemessen mit“ | behoben im Text: Lastenheft „(gemessen mit `docker.io/library/alpine:3` und `debian:stable-slim`)“, Handbuch „der Blob-Host von Docker Hub kann sich ändern“, [`ADR-0021`](../plan/adr/0021-podman-build-im-sandbox-devcontainer.md) Re-Evaluierungs-Trigger nennt den Wechsel. Dauer-Sensor fehlt weiterhin (Restrisiko, Slice §6 „offen“) |
| F-3 LOW Doc-Kommentar `podmanRelaxations` | Zeiger auf beide ADRs | behoben im Code: „(measured, see ADR 0014 and ADR 0021 for `/dev/net/tun`)“; `sandboxRelaxationWarnings`-Kommentar „six options“ |
| F-4 INFO CI-Skip auf `ubuntu-latest` | im Slice vermerkt | vermerkt (Slice §6 „CI und Freigabe“); Workflow unberührt; siehe Auflage A-4 |
| F-5 INFO Marker im Echo | nicht tautologisch | bestätigt: Beleg ist Exit-Code der `RUN`-Zeilen (`set -e`, `getent … &&`, `! getent … &&`) plus `EGRESS_BUILD_OK`; Mutation des Implementers (Host entfernt, Exit 125) zeigt Wirkung |
| F-6 INFO cloudfront-Host bei `docker` | Aussage „gemessen“ begrenzen | behoben im Text: Handbuch „für `nestedRuntime: docker` nicht gemessen (gemessen ist der Pull mit `podman`)“; CHANGELOG-Aussage unscharf, siehe D-3 |
| F-7 INFO Freigabe nur als Prosa | im Slice als Entscheidung ausweisen | **nur teilweise** — siehe D-2 und Auflage A-1 |
| F-8 INFO Kurznamen verallgemeinert | auf Messung begrenzen | behoben im Text: Handbuch „gemessen: `alpine` löst auf, `mycompany/img` nicht; andere Kurznamen sind nicht gemessen“; Spezifikation nennt dieselbe Begrenzung; Lastenheft sagt nur „voll qualifizierte Image-Namen“ |

## Lastenheft-Wortlaut gegen Messung

Zusagen in [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (0.12.0) und ihre Messdeckung:

| Zusage | Messdeckung | Urteil |
| --- | --- | --- |
| `podman build` und `docker build` (Alias) bauen `FROM`/`RUN`/`COPY` ohne Host-Socket | Lauf 2 M2/M3, Test `…_Build` (grün, selbst gefahren) | gedeckt |
| auch hinter Egress mit Default-Allowlist (alpine:3, debian:stable-slim) | Lauf 2 M5; Test `…_BuildEgress` fährt nur `debian:stable-slim` automatisch, `alpine:3` unter Egress nur manuell gemessen | gedeckt, Test deckt die Hälfte (D-4) |
| Messstand Docker 29.5.2 / Colima 0.10.3 / arm64 / VM Ubuntu 24.04 / Schalter 0; andere Architekturen und Engines nicht gemessen | Lauf 2 M0 | gedeckt; amd64, Docker Desktop, Podman-Host, `buildx`/`compose` ausdrücklich ausgenommen |
| Host-Zusage: Docker auf Hosts mit erlaubten User-Namespaces und vorhandenem `/dev/net/tun`; Wert 1 → Exit 11 | Lauf 1 M2 (Exit 11 gemessen), Lauf 2 | gedeckt |
| fehlendes `/dev/net/tun` → Start abgelehnt | Fehlerbild im Slice §6 als gemessen genannt (`docker run` Exit 127); [`ADR-0021`](../plan/adr/0021-podman-build-im-sandbox-devcontainer.md) enthält diese Messung **nicht** in der Messtabelle | Aussage plausibel, Beleg nicht in der ADR (D-5) |

Der Wortlaut sagt nicht mehr zu als gemessen. Die Rücknahme der früheren Zusage „engine-neutral, startet unter Docker (inkl. Colima) und Podman“ für `podman` ist in der Historienzeile 0.12.0 benannt und stützt sich auf die Messung (Exit 11 mit Standardwert 1; Podman-Host nie gemessen).

## Plan-vs-Code-Diff

Plan-Tabelle §3 gegen `git diff --stat 5f372bc..HEAD` (15 Dateien): jede berührte Datei steht in der Tabelle, jeder „berührt“-Eintrag hat einen Diff, jeder „nicht berührt“-Eintrag (`Dockerfile.tmpl`, `sandbox-init.sh.tmpl`, `devcontainer.json.tmpl`, READMEs, Welle-Datei, `open/README.md`, `integration.yml`, `doctor_sandbox.go`) fehlt erwartungsgemäß im Diff. Keine Abweichung. Ein Hinweis zur Tabelle: Der Eintrag zu `lastenheft.md` beschreibt noch den Plan-Zustand („Wortlaut wird dem Projektinhaber vor dem Commit vorgelegt, §4“), der Text ist inzwischen durch §1 „freigegeben“ überholt — kosmetisch.

Plan-Aussagen gegen Code: `podmanRelaxations()` (sechs Optionen, Reihenfolge wie im Test gepinnt), `egressRegistryHosts()` (fünf Hosts, Zeile in [`SPEC-019`](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) und Handbuch identisch), Warnungsform je Option unverändert. Keine Plan-Aussage ohne Fundstelle gefunden.

## Diskrepanzen

- **D-1 (LOW)** Messplan-Abdeckung: M4 „Größe des Volumes `…/containers`“ gemessen als `du` der Rootless-Storage (142 MB) statt als Volume-Größe; im Handbuch als „rund 140 MB“ formuliert, korrekt, aber die Plan-Zeile ist nicht wörtlich erfüllt. M3 `docker build`/`FROM alpine` im Build: Lauf 2 nachgeholt (gedeckt).
- **D-2 (MEDIUM, F-7 offen)** Die Abweichung von der Matrix-Zeile „Rückführung nach `next/`“ ist nicht als eigene datierte Entscheidung ausgewiesen. Der Slice nennt „Nach Vorab-Entscheidung des Projektinhabers“ (§3), „ok für Lastenheft und ADR“ (§1) und die Freigabe in §6, **ohne** den Wortlaut der Entscheidung vom 2026-10-17 („Lockerung zusätzlich ausweisen“ statt Rückführung nach `next/`) und ohne Datum an der Abweichungsstelle. §7 ist leer. [`ADR-0021`](../plan/adr/0021-podman-build-im-sandbox-devcontainer.md) trägt nur „`Accepted` durch den Projektinhaber“. Ich kann die Freigaben nicht prüfen (kein Beleg außerhalb von Prosa); sie sind nur dann tragfähig, wenn die Closure sie als Entscheidung festschreibt (A-1).
- **D-3 (INFO)** CHANGELOG-Eintrag zum cloudfront-Host sagt „without it a pull behind the restriction failed with a DNS error“ ohne Einschränkung auf den Podman-Pfad; Handbuch und ADR grenzen sauber ab. Kein Widerspruch zum Lastenheft (der nennt Hosts nicht), aber unschärfer als der Rest.
- **D-4 (INFO)** `…_BuildEgress` fährt nur `debian:stable-slim`; die Lastenheft-Aussage „alpine:3 unter Egress“ ist Lauf-2-Messung ohne Test (Zusage gemessen, nicht wiederholbar abgesichert).
- **D-5 (INFO)** Das Fehlerbild „`/dev/net/tun` fehlt → `docker run` Exit 127“ steht in Slice §6, Spezifikation und Handbuch als gemessen, aber nicht im Messprotokoll Lauf 1/2 des Slice §3 und nicht in der Messtabelle von [`ADR-0021`](../plan/adr/0021-podman-build-im-sandbox-devcontainer.md). Die Herkunft ist nicht auflösbar, ohne den Implementer zu fragen.
- **D-6 (INFO)** Mutationsbeleg des Reviewers für `make test-docker` fehlt (Sandbox ohne Netz); der Implementer hat fünf Mutationen rot gesehen (§3 „Test“). Der Verifier hat sie nicht wiederholt (Sensor `make test-docker` grün, kein Mutationslauf).
- Kein Verstoß gegen Hard Rules im Diff erkennbar (kein `//nolint`, Importe unverändert, Sprachvertrag eingehalten, [`ADR-0014`](../plan/adr/0014-nested-podman-sandbox-devcontainer.md) unverändert).

## Carveouts

- Neu: none · Gelöst: none · Unverändert: Master-Inventar nicht im Diff.

## Nicht ausgeführt

- `make verify-depguard` — keine Importregel berührt.
- `make ci` (`govulncheck`, `image-scan`) — kein Security-/Release-Pfad berührt.
- amd64, Docker Desktop, Podman als Host-Engine, `buildx`/`compose` — nicht gemessen und nicht zugesagt.

## Verdikt

**DoD-Liefer-Punkte 1–3 und die Sensor-Zeilen: erfüllt.** Offen sind ausschließlich die Closure-Zeilen (Closure-Notiz, Register, Risiko-Ausgänge, Paarungen, DoD-Häkchen) — das ist der Vorgang nach Verifier und Validator. Review-Findings F-1, F-2, F-3, F-6, F-8 sind am Code bzw. Text behoben; F-4 und F-7 sind nur teilweise gelöst (siehe Auflagen).

## Auflagen für die Closure

- **A-1 (F-7, vor `git mv`):** In Slice §7 die Freigaben des Projektinhabers (2026-10-17) als Entscheidung festschreiben: „Lockerung zusätzlich ausweisen“ statt Rückführung nach `next/` (Abweichung von der Matrix-Zeile, Bezug `AGENTS.md` §3.6) und „ok für Lastenheft und ADR“ — mit Wortlaut, Datum und Fundstelle des Beleges. Ohne diesen Beleg bleibt D-2 offen.
- **A-2 Roadmap-Ruhe-Marker:** Der Marker „Nichts in Arbeit.“ ist entfernt, solange der Slice in `in-progress/` liegt. Mit dem `git mv` nach `done/` den Marker wiederherstellen, sofern kein weiterer Slice in `in-progress/` liegt (d-check `planning`); die Welle-Zeile unter *Offene Wellen* wandert erst mit der Welle-Closure.
- **A-3 Register und §7:** Beobachtungs-Register fortschreiben oder in §7 „keine Beobachtung angefallen“ notieren. Als Kandidaten (Finding-Klassen des Reviews, nicht Verifier-Beschluss): „Neue Host-Vorbedingung ohne Vertrags- und doctor-Abdeckung“, „Zusage hängt an externer Host-Topologie ohne Dauer-Sensor“, „Matrix-Abweichung belegt nur durch Prosa“; die Klassen „Kommentar-Zeiger“ und „Messaussage auf zweiten Pfad“ ordnen sich BEO-CLI/kommentar-beschreibt-abwesendes bzw. BEO-DOK/plan-aussage-ueber-code zu. Ein `liegt in` steht nur bei verkörperter Regel (nichts verkörpert beobachtet).
- **A-4 amd64-/CI-Grenze:** `integration.yml` läuft auf `ubuntu-latest`; ob dort der Schalter 1 ist, ist ungemessen, und beide Build-Tests überspringen dann. Die Zusage hat auf amd64 und in CI **keine** Maschinenevidenz, nur die manuelle arm64-Messung. §6 muss diese Risiken in der Closure mit Ausgang „weiter offen → Register“ (nicht „entfallen“) schließen; die Folgepflicht aus [`ADR-0021`](../plan/adr/0021-podman-build-im-sandbox-devcontainer.md) (amd64-Messung über CI, Gate-Änderung ist Entscheidung des Projektinhabers) benennt einen Folge-Slice oder eine Beobachtung. Kein Docs- oder Lastenheft-Text darf amd64/CI als belegt darstellen (heute nicht der Fall).
- **A-5 Risiken aus §6:** je Risiko Ausgang eintragen. Erwartet „weiter offen“: `/dev/net/tun` ohne Degradationsstufe/`doctor` (Folge-Slice nur auf Zusage), Egress-Blob-Host ohne Dauer-Sensor, nicht pinnbares `podman`/Base-Image, CI-Skip, Speicher großer Images. Erwartet „entfallen“: Alias/BuildKit-Abgrenzung, Wortlaut-Freigabe, Plan-Tabellen-Abgleich (§7 nennt das Ergebnis: 15 Dateien, keine Abweichung).
- **A-6 Messbeleg D-5:** das `/dev/net/tun`-Fehlerbild (Exit 127) im Messprotokoll oder in [`ADR-0021`](../plan/adr/0021-podman-build-im-sandbox-devcontainer.md) als Messung nachweisen, oder die Zeile „gemessen“ in Handbuch/Spezifikation/Slice §6 zurücknehmen. [`ADR-0021`](../plan/adr/0021-podman-build-im-sandbox-devcontainer.md) ist `Accepted`; eine Änderung ist nur als Folge-ADR oder im Slice-Protokoll möglich, nicht durch Umschreiben (`AGENTS.md` §3.5).
- **A-7 Welle-Closure-Voraussetzungen:** Welle-Datei Kriterium „sieben Slices“ erfüllt erst mit diesem Slice in `done/`; Lastenheft, Spezifikation, Handbuch und CHANGELOG sagen über `nestedRuntime: podman` dasselbe (Verifier-Stichprobe: konsistent bei sechs Optionen, Host `…cloudfront…`, Kurzname-Begrenzung, Exit 11/Tun-Vorbedingung); zu entscheiden bleiben die Beobachtungs-Register-Einträge der Welle und der Release-Cut-Hinweis zu den `--quiet`-Zeilen (Welle-Datei §7). Die Validator-Prüfung (Nutzerpfad `podman build` im Container) steht noch aus.
- **A-8 Nachlauf:** Verification-Evidence (dieser Bericht) und `make gates` nach den Closure-Commits erneut; `make doc-immutable` über die um die Closure-Commits erweiterte Range.

## Abschluss

`make docs-check` nach Anlage dieser Datei: Ergebnis im Übergabetext an den Planner.
