# Slice V1: Leichter Isolationstest für Instanzen in der Suite, der schwere Vier-Instanzen-Test entfällt ersatzlos (`LH-FA-DEV-010`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Geschlossen (Entscheidungen des Projektinhabers 2026-11-01, §1). Erster von zwei Entwürfen: Der zweite, [`slice-cr-devcontainer-instanz-zusage-getestet`](../open/slice-cr-devcontainer-instanz-zusage-getestet.md), passt die Zusage im Lastenheft an das an, was die Suite danach belegt; er folgt diesem Slice.

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Isolation je Instanz; Scope: wie die Zusage belegt wird, nicht ihr Wortlaut), [`ADR-0038`](../../adr/0038-devcontainer-instanz-zweites-id-label.md) (zweites Id-Label).
Anlass: Entscheidung des Projektinhabers, 2026-11-01. `TestDevcontainerInstance_Docker_FourInstancesAreIsolated` ist nicht praxistauglich: Er scheiterte im CI-Lauf `integration-docker` auf dem GitHub-Runner an der Vorprüfung (`SwapFree` 3071 MiB, „FINDING“ nach 0,03 s) und braucht einen großen Host.

**Berührte Spec-Stellen:** keine zu ändernde. Gelesen: [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (Absatz „Instanzen“, `spezifikation.md` Zeile 837) nennt weder Zahl noch Grenze und keinen Test. Der Lastenheft-Absatz mit Zahl und Grenze ist Gegenstand des zweiten Entwurfs.

**Verantwortlich:** pt9912 (Lauf instanz-test-leicht, Hauptbaum; der Lauf zu `slice-v1-shell-erkennung-fehlender-benutzer-eng` ist ein anderer, ein eigener Zweig besteht nicht)

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-11-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Suite (`make test-docker`, `integration-docker`) enthält einen Isolationstest für zwei Instanzen, der auf dem GitHub-Runner läuft; der schwere Test mit vier Instanzen, `nestedRuntime: docker`, `ollama` und Grenzen entfällt ersatzlos (kein Messskript). Die Messung in `docs/messungen/slice-v1-devcontainer-instanz-messung/` bleibt unverändert.

### Gelesener Stand (2026-11-01)

- Der schwere Test steht in `internal/e2e/devcontainer_instance_docker_test.go`, Funktion `TestDevcontainerInstance_Docker_FourInstancesAreIsolated` (Zeilen 212 bis 413). Er startet den Container ohne Instanz und vier Instanzen (`dcInstanceCount = 4`) im Sandbox-Profil mit `devcontainer.sandbox.nestedRuntime=docker`, `…onUnavailable=warn` und `devcontainer.features.ollama.enabled=true`; `dcLimitRunArgs` setzt `--memory=3g --memory-swap=3g --cpus=1` in `runArgs`. Die Vorprüfung `dcRequireInstanceHost` (Zeilen 81 bis 104) lässt den Test mit „FINDING“ scheitern, wenn die Engine keine Memory-/Swap-Limits erzwingt, weniger CPUs als Instanzen da sind, `MemAvailable` unter `dcInstanceMinMemAvailable` (9 GiB) oder `SwapFree` unter `dcInstanceMinSwapFree` (8 GiB) liegt, oder unter 40 GiB im Docker-Root frei sind (`dcInstanceMinFreeBytes`).
- Er prüft in einem Lauf: (1) Container-, Folder- und Instanz-Label, kein `devcontainer.config_file` an Instanzen; (2) paarweise verschiedene Container-IDs und Volumes der drei Familien (Workspace, Docker-Speicher, Modelle); (3) die zurückgelesene Grenze; (4) eine Marker-Datei je Volume, in keinem anderen Container sichtbar; (5) `up`/`shell` ohne `--instance` treffen den Container ohne Instanz, zweites `up --instance` trifft denselben Container; (6) `down` ohne und mit `--instance`, `down --instance nope`; (7) `--recreate` einer Instanz lässt Volumes und Datei stehen; (8) ohne Container ohne Instanz enden `up`, `up --recreate`, `down`, `shell` mit Exit-Code 10, ohne Hilfs-Container, ohne Änderung.
- Die übrigen zwei Tests der Datei laufen schon leicht: `TestDevcontainerInstance_Docker_InvalidNamesEndWithExit10` (kein Container) und `TestDevcontainerInstance_Docker_StandardProfileSharesTheProjectFolder` (Standardprofil, drei Container, Bind-Mount).
- Das Sandbox-Profil ohne `nestedRuntime` und ohne `ollama` trägt genau ein Volume: `volumeMount(name+"-workspace", …)` mit `-${devcontainerId}` (`internal/hexagon/application/devcontainer_sandbox.go`, Zeile 104 und 139 bis 141). Die Mounts `-docker` (Zeile 117) und `-ollama` (Zeile 123) entstehen über dieselbe Funktion `volumeMount`; Text-Pins der drei Mounts stehen in `internal/hexagon/application/generate_ollama_test.go` (Zeilen 27 bis 30) und `generate_sandbox_test.go` (Zeile 67).
- Die Beleg-Stellen des schweren Tests: `docs/maintainer/zusagen-und-messstaende.md` Zeile 178 (Abschnitt „Mehrere Devcontainer je Projekt“, Fundstelle). Der Test trägt in Kommentar und Docstring die Bedingung „fails with a finding when the host does not carry the stage“ (Zeilen 11 bis 14, 76 bis 80).

### Was die Suite danach belegt und was nicht

Ohne die schweren Bestandteile belegt ein Sandbox-Test mit Standard-Einstellungen: Container- und Label-Identität, ein Workspace-Volume je Instanz (paarweise verschieden), Marker-Isolation über alle geordneten Paare, die Container-Wahl ohne `--instance`, `down`, `--recreate` (Volume und Datei bleiben), Exit-Code 10 ohne Container ohne Instanz. Er belegt **nicht**: die Volumes für Docker-Speicher und Modelle (nur durch Text-Pins der Mounts belegt, siehe oben), die Grenze von 1 CPU und 3 GiB je Instanz, drei und vier Instanzen, den Zustand mit laufendem verschachtelten Docker-Daemon.

**Fragen und Wege — entschieden** (Entscheidung Projektinhaber 2026-11-01):

1. **Form des leichten Tests: Weg A.** Sandbox-Profil mit Standard-Einstellungen (`dcGenerateProject(ctx, t, true, nil)`), ohne `nestedRuntime`, `ollama`, Grenzen und Speicher-Vorprüfung; ein Workspace-Volume je Instanz, Marker in jedem Container. Die Wege B und C entfallen. Begründung des Projektinhabers: „Workspace-Isolation bei zwei Instanzen reicht“.
2. **Anzahl der Instanzen: zwei** (statt drei) plus der Container ohne Instanz. Der Rückfall des Entwurfs gilt von Anfang an.
3. **Messskript: entfällt.** Der schwere Test `TestDevcontainerInstance_Docker_FourInstancesAreIsolated` wird ersatzlos entfernt; er „machen wir nicht mehr“ (er scheiterte im Lauf `integration-docker` an `SwapFree` 3071 MiB < 8192 MiB). Es entsteht kein Skript unter `tools/messung/`; die Wege A bis C der Frage entfallen.
4. **Messstand:** `docs/messungen/slice-v1-devcontainer-instanz-messung/` bleibt unverändert (gelesen: der Ordner nennt den Test nicht als Beleg, Suche `FourInstances` ohne Treffer dort); es gibt keinen neuen Lauf.

**Test und Wächter.** Name: `TestDevcontainerInstance_Docker_TwoInstancesAreIsolated`. Der Test ist **nicht gelaufen** (kein `make test-docker`, Entscheidung des Projektinhabers; der Lauf `integration-docker` auf `main` entscheidet). Rot färbende Mutationen, aus dem Quelltext gelesen: (1) `up --instance` ohne das Instanz-Label (oder mit dem Label einer anderen Instanz) — die Prüfung `label u-boot.instance` schlägt an; (2) das Workspace-Volume der Instanzen auf einen gemeinsamen Namen setzen (Produktcode `devcontainer_sandbox.go`, `-${devcontainerId}` entfernt) — `dcRequireDistinct "workspace volume"` und die Marker-Paarprobe (`code == 0` gegen `i == j`) schlagen an; (3) `up` ohne `--instance` wählt die neueste Instanz — `again.ContainerID != existing.ContainerID` und der Hostname der `shell`; (4) `--recreate` entfernt das Volume — `volumeAt(...) != workspaces[1]` und die gelöschte Marker-Datei.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Wortlaut der Zusage in [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)** — das ist eine Vertragsänderung mit Freigabe des Projektinhabers: [`slice-cr-devcontainer-instanz-zusage-getestet`](../open/slice-cr-devcontainer-instanz-zusage-getestet.md). Bis dahin bleibt das Lastenheft weiter als der Test; Risiko in §6.
- **Änderung von Produktcode** — der Slice belegt Bestandsverhalten, er ändert keines; ein Fund geht als Fehler an einen eigenen Vorgang.
- **Die Entscheidung, `make test-docker` nur in der CI zu fahren** — [`slice-v1-test-docker-in-ci`](../open/slice-v1-test-docker-in-ci.md).
- **Eine Lockerung des Gates** — kein Gate wird gesenkt; der Test wird aus der Suite genommen und durch einen leichteren ersetzt, die Mutationsprobe in der DoD belegt, dass der leichte Test noch rot werden kann (`AGENTS.md` §3.6).
- **Neue Messungen der Instanz-Zusage auf anderen Hosts** — die Messung bleibt die vom 2026-10-31.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **Leichter Test:** `TestDevcontainerInstance_Docker_TwoInstancesAreIsolated` liegt in `internal/e2e/devcontainer_instance_docker_test.go`, trägt keine Vorprüfung auf `MemAvailable`/`SwapFree`/Docker-Root, keine Grenzen und kein `nestedRuntime`/`ollama` (Weg A, §1 Frage 1), und läuft im Lauf `integration-docker` auf dem GitHub-Runner grün (Beleg: der Lauf auf `main` oder Branch; ein lokaler `make test-docker` ist nicht gefordert). Beleg: Lauf `integration-docker` 41285105546 auf `main` (Commit `a6264402`), Test `--- PASS` in 16,32 s (§7).
- [x] **Schwerer Test entfernt:** der schwere Test ist ersatzlos aus `internal/e2e` entfernt (kein Messskript, §1 Frage 3); die Konstanten `dcInstanceMinMemAvailable`, `dcInstanceMinSwapFree`, `dcInstanceMinFreeBytes`, `dcInstanceMemory`, `dcInstanceNanoCPU`, `dcInstanceCount`, `dcDockerDest`, `dcOllamaDest` und die Helfer `dcRequireInstanceHost`, `dcLimitRunArgs`, `dcMeminfoBytes` sind entfernt; `make lint` meldet keine ungenutzte Konstante.
- [ ] **Wächter und Doku-Verweise** (nicht vollständig: (a) ist nicht erfüllt, (b) bis (d) sind erfüllt; siehe §7): (a) Mutationsprobe am leichten Test, im Bericht des Implementers (etwa: Marker-Pfad eines Volumes teilen, `--instance` ohne Label, `up` ohne `--instance` wählt die neueste Instanz) — jede Mutation färbt den Test rot; (b) `docs/maintainer/zusagen-und-messstaende.md` (Fundstelle Zeile 178) nennt den leichten Test statt des schweren Tests; (c) Kopfkommentar der Test-Datei (Zeilen 3 bis 14) und der Docstring beschreiben den Ist-Stand (`AGENTS.md` §3.7); (d) Slice-/Register-Verweise auf den schweren Test sind nachgezogen (Fundstellen per Suche über den ganzen Baum, §3).
- [x] `make gates` grün (Verifier auf `172387b5`; Schluss-Lauf der Closure im Bericht).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: `docs/maintainer/zusagen-und-messstaende.md`; `CHANGELOG.md` Zeile 16 (Satz „Pinned by unit tests and by a Docker test with four instances in the sandbox profile (`nestedRuntime: docker`, feature `ollama`, one CPU and 3 GiB each) …“) nennt den schweren Test der Sache nach und ist auf den Ist-Stand gesetzt; `docs/user/` nennt den Test nicht (Handbuch `benutzerhandbuch.md` Zeilen 896 bis 925 gelesen: weder Zahl gleichzeitiger Instanzen noch Grenze noch Test; unverändert).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/e2e/devcontainer_instance_docker_test.go` | update | schwerer Test (Zeilen 212 bis 413) entfernt, leichter Test ergänzt; Helfer und Konstanten des schweren Tests (Zeilen 33 bis 54, 56 bis 124) entfernt, soweit unbenutzt (gelesen: `dcUpInstance` und `dcRequireDistinct` nutzt auch `StandardProfileSharesTheProjectFolder`, `dcWatchNoHelper` auch `InvalidNamesEndWithExit10`; `dcFolderContainers`, `dcRunning` und `dcMarkerPaths` nutzt nur der schwere Test, der leichte Test braucht sie wieder; alle Helfer haben Aufrufer nur in dieser Datei); Kopfkommentar Zeilen 3 bis 14 neu; Akzeptanz [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) |
| `docs/messungen/slice-v1-devcontainer-instanz-messung/README.md` | unverändert | Suche nach `Test` ohne Treffer: nennt den Test nicht; Rohdaten-Dateien unverändert; kein Messskript (§1 Frage 3) |
| `docs/maintainer/zusagen-und-messstaende.md` | update | Fundstelle (Zeilen 177 bis 183) nennt den leichten Test und was die Suite nicht läuft; Abschnitt „Mehrere Devcontainer je Projekt“ (Zeilen 168 bis 186) |
| `CHANGELOG.md` | update | Zeile 16: die Aussage zum Docker-Test mit vier Instanzen, einer CPU und 3 GiB auf den leichten Test gesetzt (Abschnitt `[Unreleased]`) |
| `docs/plan/planning/open/README.md`, `docs/plan/planning/in-progress/roadmap.md` | update | Zeile des Slice aus der `open/`-Liste entfernt; Roadmap-Zeiger unter „Offene Wellen“ ergänzt (der Zeiger des anderen Slice in `in-progress/` bleibt) |
| `docs/plan/planning/open/slice-cr-devcontainer-instanz-zusage-getestet.md` | update | Entwurf auf die Entscheidung gesetzt: zwei Instanzen plus Container ohne Instanz, kein Messskript, Variante 1 nur „zwei Instanzen, Workspace-Volumes“ (Review F-1); Lastenheft unverändert |
| `internal/e2e/devcontainer_start_docker_test.go` | update | `dcInspectDoc.HostConfig`: die Felder `Memory`, `MemorySwap`, `NanoCpus` ohne Leser entfernt (Suche über `internal/`; Leser von `HostConfig` sind `Binds` und `Mounts`, Review F-4) |
| `docs/plan/planning/observations/BEO-E2E/…` | update, neu | Closure (§7): Evidence-Dateien; Sichtung in §8 |
| `docs/plan/planning/done/slice-v1-devcontainer-instanz-label.md`, `docs/reviews/2026-11-01-slice-v1-devcontainer-instanz-label*.md`, `docs/plan/planning/observations/BEO-APP/container-identitaet-nur-ueber-label/evidence/slice-v1-devcontainer-instanz-label.md` | unverändert (Chronik) | nennen den Testnamen (Suche `FourInstances`); abgeschlossene Slices und Reports werden nicht umgeschrieben; der Implementer liest die Stellen und trägt die Fundstelle ein (`AGENTS.md` §3.10) |
| `spec/lastenheft.md`, `spec/spezifikation.md`, `spec/architecture.md`, `docs/user/`, `README.md`, `README.de.md`, `.github/workflows/integration.yml`, `Makefile` | unverändert | Zusage-Wortlaut: zweiter Entwurf; die Spec nennt keinen Test (Suche `FourInstances` ohne Treffer in `spec/`, `docs/user/`, `CHANGELOG.md`; der Treffer in `CHANGELOG.md` ist ein Satz ohne Testnamen, siehe dessen Zeile); der Workflow ruft `make test-docker` ohne Filter (gelesen, `integration.yml` Zeilen 46 bis 47) |

Der Implementer ergänzt die Tabelle im ersten Lauf und gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab (`AGENTS.md` §3.10).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; der Projektinhaber hat die Fragen 1 bis 4 aus §1 beantwortet (Vermerk in §7); WIP-Limit 1 pro Lauf (Baseline-Regelwerk `modul-05-planning-harness.md` §Trigger je Lifecycle-Übergang und WIP-Limit): der Slice [`slice-v1-shell-erkennung-fehlender-benutzer-eng`](../done/slice-v1-shell-erkennung-fehlender-benutzer-eng.md) liegt in `in-progress/` und gehört einem anderen Lauf; dieser Lauf hält nur diesen Slice (Anweisung des Projektinhabers, 2026-11-01).

**Rückführungen — vorab benannt:**

- `in-progress` → `open` (blockiert): Der leichte Test läuft auf dem Runner nicht grün und der Grund liegt im Produkt (dann ein eigener Vorgang für den Fehler) oder im Runner (dann ein eigener Vorgang; weiter verkleinern geht nicht unter zwei Instanzen).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, der Lauf `integration-docker` mit dem leichten Test ist grün (oder ein roter Lauf ist als Carveout dokumentiert, nie als stilles Rot), Closure-Notiz mit Lerneintrag in §7.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Das Lastenheft sagt mehr, als die Suite belegt:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) sagt vier Instanzen mit `nestedRuntime: docker`, `ollama` und Grenze zu; nach diesem Slice belegt die Suite dafür nur Teile (siehe §1), die Zusage selbst belegt dann ein manuelles Skript (`BEO-DOK/zusage-weiter-als-messung`). — **Ausgang:** weiter offen → [`slice-cr-devcontainer-instanz-zusage-getestet`](../open/slice-cr-devcontainer-instanz-zusage-getestet.md).
- **Der leichte Test belegt die Docker-Speicher- und Modell-Volumes nicht** (nur Text-Pins der Mounts, `generate_ollama_test.go` Zeilen 27 bis 30). — **Ausgang:** weiter offen → `BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf` (Evidence `evidence/slice-v1-devcontainer-instanz-test-leicht.md`, §7).
- **Der Runner trägt drei Sandbox-Container gleichzeitig nicht** (Container ohne Instanz plus zwei Instanzen; Speicher, Zeit gegen `dcTestBudget` von 20 Minuten in `devcontainer_start_docker_test.go` Zeile 49 und `timeout-minutes: 30` des Jobs). Der neue Test ist nicht gelaufen. — **Ausgang:** entfallen: Der Lauf `integration-docker` 41285105546 auf `main` ist grün, der Test steht dort als `--- PASS` (16,32 s), nicht als `SKIP`; der Rücksprung nach `open/` (§4) trat nicht ein.
- **Zwei Annahmen des Tests sind erst durch den Lauf belegt** (Review F-6): (a) das Sandbox-Profil mit Standard-Einstellungen (`dcGenerateProject(ctx, t, true, nil)`, ohne `egress.enabled`/`onUnavailable`, die andere Sandbox-Tests setzen) startet auf dem Runner; (b) die Marker im Workspace-Volume sind per `exec -u root` les- und schreibbar. — **Ausgang:** entfallen: Derselbe Lauf ist grün; (a) das Standard-Sandbox-Profil startete auf dem Runner, (b) die Marker waren per `exec -u root` les- und schreibbar (sonst wäre die Marker-Paarprobe rot geworden).
- **Der Wächter für `--recreate` trägt nur die Marker-Prüfung** (Review F-8): der Volume-Name folgt aus den Labels und bleibt gleich, `volumeAt != workspaces[1]` fängt ein entferntes Volume nicht; es trägt `test -e` auf die Marker-Datei nach `--recreate`. — **Ausgang:** weiter offen, akzeptiert: Die Mutation „`--recreate` entfernt das Volume“ ist nur aus dem Quelltext bewertet, nicht gefahren (Verify-Report V-2, [`…-verify`](../../../reviews/2026-11-01-slice-v1-devcontainer-instanz-test-leicht-verify.md)); der Test trägt sie mit der Marker-Prüfung allein.
- **Die Handbuch-Bedingung „soweit aktiviert“** (Docker-Speicher und Modelle, `benutzerhandbuch.md` Zeilen 903 bis 906) hat keinen Lauf, nur Text-Pins der Mounts (Review F-5); der CI-Lauf belegt sie nicht. — **Ausgang:** weiter offen → [`slice-cr-devcontainer-instanz-zusage-getestet`](../open/slice-cr-devcontainer-instanz-zusage-getestet.md).
- **Der Messweg entfällt:** es gibt kein Skript (Entscheidung Projektinhaber 2026-11-01); die Zusage für vier Instanzen mit `nestedRuntime: docker`, `ollama` und Grenze belegt dann nur die Messung vom 2026-10-31. — **Ausgang:** entfallen (kein Messweg); die Diskrepanz zur Zusage trägt Risiko 1.
- **Der Test-Name steht in abgeschlossenen Dokumenten** (Slice, Reports, Evidence). — **Ausgang:** entfallen: Chronik, wird nicht geändert; Verweise stehen nur in Chronik und im Slice selbst (Verify-Report, Negativbefund „Entfernte Konstanten und Helfer“).
- **Der Kopfkommentar der Test-Datei beschreibt weiter den schweren Test** (Zeilen 3 bis 14). — **Ausgang:** entfallen: DoD-Punkt (c) ist erfüllt (Verify-Report, DoD-Abgleich 3(c)).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren) · `grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln.

Vermerk: Die Fragen 1 bis 4 aus §1 sind entschieden (Entscheidung Projektinhaber 2026-11-01): Weg A, zwei Instanzen, kein Messskript, Messstand unverändert.

- **Was hat funktioniert:** Der Ersatz läuft auf dem GitHub-Runner. Beleg (nur Gelesenes): Lauf `integration-docker` **41285105546** auf `main`, Commit `a6264402`, gestartet 2026-11-01T07:31:16Z, Dauer 15 min 14 s, Ergebnis `success` (`gh run list --branch main`, `gh run view 41285105546`). Im Job-Log (`gh api repos/pt9912/u-boot/actions/jobs/124784712075/logs`, Befehl `go test -v -tags docker -timeout 30m ./...`) steht `--- PASS: TestDevcontainerInstance_Docker_TwoInstancesAreIsolated (16.32s)` (Ausgabezeit 07:46:27Z, das Paket gibt am Ende aus); für diesen Test steht kein `SKIP` und kein `FAIL` im Log, die beiden Nachbarn derselben Datei stehen ebenfalls auf `PASS` (1,52 s und 6,54 s). Review (F-1 bis F-8) und Verifikation (V-1 bis V-4) lagen vor; die Fix-Runde `32cecdf6` behob F-1 bis F-5.
- **Was ging anders als geplant:** (1) Der Test war bei Review und Verifikation nicht gelaufen (Entscheidung des Projektinhabers, kein `make test-docker`); der Verify-Befund V-1 (MEDIUM, „Test nicht gelaufen“) ist durch den obigen Lauf aufgelöst. (2) **DoD 3(a) ist nicht erfüllt wie gefordert** (V-2, LOW): Es gibt kein Mutationsprotokoll; die vier Mutationen (Instanz-Label, gemeinsames Workspace-Volume, `up` wählt die neueste Instanz, `--recreate` entfernt das Volume) hat der Verifier **aus dem Quelltext bewertet, nicht gefahren**; die Wirksamkeit der Wächter ist gelesen, nicht gemessen. (3) V-3 (LOW): Der Tag `docker` baut kein lokales Gate (`make lint`, `make test`); die Übersetzbarkeit war nur gelesen und ist erst durch den CI-Lauf belegt. Der Anker dafür ist [`slice-v1-test-docker-in-ci`](../open/slice-v1-test-docker-in-ci.md).
- **Steering-Loop-Eintrag:** gezählt, nicht verkörpert (kein `liegt in`). Lehre: Ein Docker-Test, den kein lokales Gate baut und der auf Entscheidung nicht lokal läuft, hat als einzigen Beleg den CI-Lauf auf `main`; Review und Verifikation konnten ihn nur lesen (V-1, V-3), der Closure-Trigger „grüner Lauf `integration-docker`“ wartete auf den Push. Finding-Klassen des Verifiers: Test ohne Lauf · Wächter ohne Mutationslauf · Test-Tag nicht von einem lokalen Gate gebaut. Anker der Entscheidung: [`slice-v1-test-docker-in-ci`](../open/slice-v1-test-docker-in-ci.md); kein Register-Eintrag angelegt, keine Schwelle behauptet.
- **Beobachtungs-Register (`../observations/`):** keine neue Beobachtung angelegt. `evidence/slice-v1-devcontainer-instanz-test-leicht.md` in `BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf/` ergänzt (der einzige Docker-Test, der `nestedRuntime: docker` mit `ollama` startete, ist aus der Suite entfernt), `state.md` dort nachgetragen; Zähler und Ausgang bleiben unverändert (nicht gesetzt). Die übrigen in §8 gesichteten Beobachtungen (`BEO-E2E/test-docker-im-verschachtelten-daemon-nicht-gleichwertig`, `BEO-E2E/test-zaehlt-container-paralleler-pakete`, `BEO-E2E/docker-test-hinterlaesst-zwischencontainer`): kein Auftreten gelesen, keine Evidence-Datei, `state.md` unverändert; ob der neue Test Container oder Volumes zurücklässt, ist nicht gemessen und zählt nicht als Beleg.
- **Folge-Slices:** [`slice-cr-devcontainer-instanz-zusage-getestet`](../open/slice-cr-devcontainer-instanz-zusage-getestet.md) — ist eine Datei in `open/`
- **Risiken aus §6:** neun Risiken mit Ausgang: entfallen (Runner trägt drei Container, zwei Annahmen des Tests, Messweg, Testname in Chronik, Kopfkommentar), weiter offen (Zusage weiter als die Suite → Entwurf `slice-cr-devcontainer-instanz-zusage-getestet`; Docker-Speicher- und Modell-Volumes → `BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf`; `--recreate`-Wächter, akzeptiert, V-2; Handbuch-Bedingung „soweit aktiviert“ → Entwurf).
- **Drei Paarungen:** Anker: kein `liegt in` in diesem Slice, nichts zu prüfen. Folge-Slice: [`slice-cr-devcontainer-instanz-zusage-getestet`](../open/slice-cr-devcontainer-instanz-zusage-getestet.md) liegt in `open/`. Register: `BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf` trägt die neue Evidence-Datei. Haken folgt nach dem Move.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `internal/e2e` (Kürzel E2E), `tools/` (WZG), `spec/`, `harness/`, `docs/plan/`, `docs/messungen/` (DOK; hier nur `docs/maintainer/` und Planung) — jede erfüllt laut Modus-Deklaration (`harness/conventions.md`) die Schwelle; keine zu grobe.

**Vorgelagert — offene Beobachtungen sichten:** Register `BEO-E2E` und `BEO-DOK` durchgegangen (Verzeichnisnamen gelesen, Treffer-Texte gelesen für die genannten):
- `BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf` — Treffer: die Kombination `nestedRuntime` plus `ollama` hat keinen Lauf der Suite; der Umbau nimmt den einzigen Test, der sie startet, aus der Suite (Risiko 2). Zähler-Stand: aus den Dateien unter `evidence/` zu lesen (Implementer).
- `BEO-E2E/test-docker-im-verschachtelten-daemon-nicht-gleichwertig` und `BEO-E2E/test-zaehlt-container-paralleler-pakete` — Treffer: der leichte Test startet mehrere Container; `dcFolderContainers` filtert über das Label des Projektordners (gelesen, Zeile 155 bis 162), nicht über alle Container. Zähler-Stand: vom Implementer zu lesen.
- `BEO-E2E/docker-test-hinterlaesst-zwischencontainer` — Treffer: Aufräumen der Container und Volumes des neuen Tests (der schwere Test räumt Volumes des Container ohne Instanz selbst auf, Zeilen 381 bis 388). Zähler-Stand: vom Implementer zu lesen.
- `BEO-DOK/zusage-weiter-als-messung` — Treffer: siehe Risiko 1, Gegenstand des zweiten Entwurfs.

### Sub-Area: `internal/e2e`

- **Modus:** GF
- **Konventionen-Dichte:** Build-Tag `//go:build docker` und Fake-Clock-Pflicht in der Architektur-Sicht (`harness/conventions.md`, Modus-Deklaration E2E).
- **Phase-Reife:** Phase 4 (Bestandstest wird ersetzt, Gate `integration-docker` besteht).
- **Evidenz-/Diskrepanz-Risiko:** niedrig; die Diskrepanz zwischen Zusage und Suite ist Gegenstand des zweiten Entwurfs.
- **Reconciliation-Aufwand:** gering; Folge-Slice: `slice-cr-devcontainer-instanz-zusage-getestet`.

*Alle übrigen berührten Sub-Areas (`tools/`, `spec/`/`harness/`/`docs/plan/`): GF.*
