# Slice v1: e2e-Helfer `startLabeled` — Container-ID aus stdout, Image vorab ziehen (CI-Fehler `integration-docker`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** `v0.11.0` (Release-Versionen sind Wellen, [`MR-003`](../../../../harness/conventions/MR-003-roadmap-wellen.md)). **Blocker des Cuts** [`slice-v1-release-cut-v0.11.0`](../open/slice-v1-release-cut-v0.11.0.md): dessen DoD verlangt einen grünen CI-Lauf auf dem Tag-Commit, und `integration-docker` ist auf `7f9f2f5` rot. Der Slice läuft vor dem Cut.

**Bezug:** [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) und [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (die Datei `internal/e2e/sandbox_status_docker_test.go` trägt beide im Kopfkommentar), [`ADR-0022`](../../adr/0022-fallback-ausweisung-statusdateien-im-container.md) (Statusdateien, Vertrauensgrenze des Inspector-Adapters, die der betroffene Test gegen eine echte Engine prüft), [`ADR-0004`](../../adr/0004-ci-system.md) (CI-System; der Job `integration-docker` ist dessen Sensor). Der Slice **lockert kein Gate** ([`AGENTS.md`](../../../../AGENTS.md) §3.6): Er macht einen Test robust, der auf einem Runner ohne Image-Cache fälschlich rot wird; der Test prüft danach dasselbe wie vorher.

**Berührte Spec-Stellen:** — (nur Testcode; Lastenheft und Spezifikation bleiben unverändert).

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** `TestE2E_SandboxDevcontainer_StatusFiles_InspectorAgainstEngine` ist auf einem Runner grün, auf dem `docker.io/library/alpine:3` nicht gecacht ist; der Helfer `startLabeled` liest die Container-ID nur aus stdout und weist jede Ausgabe ab, die keine ID ist.

### Befund (gelesen am 2026-10-18)

- **Der Fehler:** Der Lauf `39103939800` (Workflow `integration-docker`) auf `7f9f2f5` meldet `--- FAIL: TestE2E_SandboxDevcontainer_StatusFiles_InspectorAgainstEngine (1.02s)` und `sandbox_status_docker_test.go:362: seed files: Error response from daemon: page not found` (gelesen mit `gh run view 39103939800 --log`, Zeilen 377 und 378 der Ausgabe). `7f9f2f5` ändert nur die Roadmap; der Fehler liegt im Bestand, nicht in dieser Änderung.
- **Die Ursache, an der Stelle im Code gelesen:**
  - `dockerOutput` (`internal/e2e/sandbox_devcontainer_docker_test.go`, Zeile 40) liefert `exec.CommandContext(…).CombinedOutput()`: stdout und stderr gemischt.
  - `startLabeled` (`internal/e2e/sandbox_status_docker_test.go`, Zeilen 64 bis 75) startet `docker run -d --label devcontainer.local_folder=… <image> sleep 1800`, nimmt `strings.TrimSpace(out)` der gemischten Ausgabe als Container-Namen und übergibt ihn an `docker rm` und an die Aufrufer.
  - Zwei Aufrufe starten ein Registry-Image ohne vorherigen Pull: `startLabeled(ctx, t, folder, "docker.io/library/alpine:3")` in Zeile 351 und Zeile 405 desselben Tests. Ist das Image nicht gecacht, schreibt `docker run -d` „Unable to find image … locally“ und den Pull-Fortschritt nach stderr, danach die ID nach stdout; die gemischte Ausgabe ist mehrzeilig. `docker exec "<mehrzeiliger Text>"` (Zeile 361) ergibt vom Daemon `404 page not found`, weil der Name als Pfad im API-Aufruf landet.
  - Lokal war das Image gecacht, `docker run -d` schrieb nur die ID; deshalb war der Test bei jedem lokalen Lauf grün. Das Verhalten von `docker run` auf dem frischen Runner ist aus dem Diagnose-Lauf abgeleitet (Ausgabe auf Kommandoebene nachgestellt); der Log nennt nur die Folge (`page not found`), nicht den Pull-Text. Der Implementer stellt es vor der Änderung selbst nach (§3, Schritt 1).
- **Weitere Verwender eines Registry-Images in `internal/e2e` (gelesen, Wortsuche nach `docker run`, `startLabeled`, `alpine`, `FROM`, `exec.Command`):**
  - **Gleicher Fall (`docker run -d`, Ausgabe als Name):** nur die beiden Aufrufe oben. Alle anderen Aufrufe von `startLabeled` (Zeilen 168, 185, 210, 235, 261, 281, 305, 315, 326) nutzen ein selbst gebautes Image (`buildProjectImage`, lokal vorhanden, kein Pull). Die `docker run -d --name <name> …`-Starts in `sandbox_devcontainer_docker_test.go` (Zeilen 225 und 573) nutzen ebenfalls ein gebautes Image und verwenden die Ausgabe nur in der Fehlermeldung, den Namen setzen sie selbst.
  - **Anderer Fall (Registry-Image im Container oder im Build):** `docker run --rm --network=host docker.io/library/alpine:3 echo NESTED_OK` im Container (Zeile 423; der innere Daemon zieht, das Ergebnis wird mit `strings.Contains(out, "NESTED_OK")` geprüft, Pull-Text stört nicht), `docker pull -q alpine:3` im Container (Zeile 688), `FROM docker.io/library/alpine:3` und `FROM alpine:3` in den Skripten, die im Container laufen (Zeilen 308 und 536), `image: alpine:3` im dortigen Compose (Zeile 537), und die generierten Dockerfiles (`FROM mcr.microsoft.com/devcontainers/base:debian`, Vorlage `Dockerfile.tmpl`, Zeile 3), die `docker build -q` baut (`buildProjectImage` und die Aufrufe in `sandbox_devcontainer_docker_test.go`): der Pull läuft im Build, die Ausgabe von `-q` ist die Image-ID, der Test liest sie nicht, sondern vergibt den Tag selbst. Kein Befund dort, solange kein Test eine Container-ID aus gemischter Ausgabe liest.
  - Die Akzeptanz-Tests (`postgres_*`, `keycloak_*`, `otel_*`, `logs_*`, `updown_*`, `down_volumes_*`) bauen über den Produktions-Adapter (Compose); sie lesen keine ID aus `docker run`.
- **`-q` und `--pull`:** `docker run --help` (gelesen, Host-CLI) nennt `-q, --quiet` („Suppress the pull output“) und `--pull`; dass `--pull=missing` der Default ist und deshalb nichts ändert, ist die Annahme aus der Dokumentation der CLI, nicht gelesen. Wirksam ist `-q` oder das getrennte Erfassen von stdout; der Implementer prüft beides an der CLI des Test-Images (`u-boot:test-docker-tools`), nicht an der Host-CLI.

### Nebenbefund (nur genannt, nicht Teil des Slice)

Das Image `u-boot:test-docker-tools` installiert `docker-ce-cli` ohne Version (`Dockerfile`, Stage `test-docker-tools` ab Zeile 97, Installation in Zeile 110: `docker-ce-cli docker-compose-plugin`); die CLI-Version des Tests ist damit nicht reproduzierbar, und die Pull-Ausgabe hängt vom CLI ab. Ein Schritt `docker version` vor `make test-docker` im Workflow `integration.yml` machte sie im Log sichtbar. Beides ist eine Gate- bzw. CI-Änderung ([`AGENTS.md`](../../../../AGENTS.md) §3.6) und braucht einen eigenen Vorgang; **hier nicht geändert**.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Änderung am Workflow `integration.yml`, am `Dockerfile` (Stage `test-docker-tools`) oder an `repo.mk`:** Gate- und CI-Änderungen verlangen einen Plan- oder ADR-Anker (§3.6) und gehören nicht zu einem Testfix; der Nebenbefund ist ein eigener Vorgang, falls der Projektinhaber ihn will.
- **Ein Test-Image-Cache oder ein Vorwärmen im CI-Job:** wäre eine Pipeline-Änderung, und der Test bliebe lokal rot, wenn das Image fehlt; der Fix gehört in den Test.
- **Die Aufrufe von `docker build -q` und die Skripte, die im Container laufen:** kein Befund (siehe oben), ihre Ausgabe wird nicht als Container-ID gelesen.
- **Der CI-Nachweis:** erst nach dem Push, und der Push braucht die ausdrückliche Freigabe des Projektinhabers je Vorgang; der Slice endet mit dem lokalen Nachweis und benennt den CI-Lauf als offenen Beleg (§6).
- **Die Beobachtung `BEO-E2E/ci-ueberspringt-nested-podman-build`:** betrifft andere Tests (Skip), nicht diesen Helfer.

## 2. Definition of Done

Zwei Liefer-Punkte (Helfer, Tests des Helfers); alles andere sind die Pflichten der Closure.

- [x] **Helfer:** `startLabeled` liest die Container-ID nur aus stdout (oder aus der letzten nichtleeren Zeile), prüft sie als 64-stelligen Hex-String und endet mit `t.Fatalf`, wenn die Ausgabe keine ID ist; die Pull-Ausgabe von `docker run -d` ist unterdrückt (`-q`) oder getrennt erfasst; `alpine:3` wird in `TestE2E_SandboxDevcontainer_StatusFiles_InspectorAgainstEngine` vorab mit `docker pull -q` gezogen.
- [x] **Tests des Helfers ohne Docker-Pflicht:** Die ID-Auswertung ist eine reine Funktion mit Tests auf Fake-Ausgaben (nur ID; Pull-Zeilen plus ID; leere Ausgabe; ID mit Fehlerzeile; Ausgabe ohne ID); Mutation: ohne die Hex-Prüfung oder mit gemischter Ausgabe schlägt mindestens ein Fall fehl (Nachweis in §7).
- [x] **Lokaler Nachweis:** Der Test lief auf einem Host ohne gecachtes `alpine:3` grün (`docker rmi docker.io/library/alpine:3`, danach der Lauf; Colima) **und** `make test-docker` ist grün; die Ausgaben stehen in der Closure-Notiz.
- [x] `make gates` grün, `make doc-trace` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: keine öffentliche Schnittstelle berührt; falls `harness/replay.md` oder `docs/maintainer/quality.md` den Helfer oder das Verhalten des Targets beschreiben, nachgezogen (gelesen vom Implementer, sonst in §7 vermerkt).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne Wellen-Betrieb hier geprüft.

## 3. Plan (vor Code)

Jede Aussage über Code steht mit Fundstelle ([`AGENTS.md`](../../../../AGENTS.md) §3.10); der Implementer gleicht die Tabelle nach der letzten Fix-Runde gegen `git diff --stat` ab, und eine Datei, die er als „unverändert“ führt, belegt er durch eine gelesene Aussage.

**Ansatz (Schritte):**

1. **Nachstellen:** `docker rmi docker.io/library/alpine:3`, dann `docker run -d alpine:3 sleep 5` mit getrenntem stdout und stderr (Colima): der Implementer hält fest, was auf stderr steht und was auf stdout. Der Test (ohne Änderung) muss dabei rot werden; sonst ist die Ursache nicht belegt und der Plan ist zu berichtigen.
2. **Helfer:** Eine reine Funktion `containerIDFromRunOutput(stdout string) (string, error)` nimmt die letzte nichtleere Zeile, verlangt `^[0-9a-f]{64}$` und nennt in der Fehlermeldung die gesamte Ausgabe. `startLabeled` ruft `docker run -d -q …` mit **getrennt erfasstem** stdout (stderr nur für die Fehlermeldung) und wertet nur stdout aus. Der Helfer `dockerOutput` selbst bleibt unverändert (viele Aufrufer; sein Vertrag „gemischte Ausgabe“ gilt weiter).
3. **Test des betroffenen Falls:** `docker pull -q docker.io/library/alpine:3` vor Zeile 351 (mit `t.Fatalf` bei Fehler). Der Helfer allein genügt für Korrektheit; der Pull macht den Test zusätzlich unabhängig vom Pull-Zeitpunkt (der `ctx` des Tests hat zehn Minuten, das Pull-Zeitbudget fließt in ihn ein).
4. **Tests des Helfers:** untagged (ohne `//go:build docker`), damit `make test` sie ohne Engine fährt; die reine Funktion und ihre Tests stehen in einer untagged Datei desselben Pakets `e2e_test`. Der Implementer prüft, ob `make test` das Paket `internal/e2e` überhaupt kompiliert und ausführt (nur `_test.go`-Dateien im Paket), und ob die Coverage-Berechnung dadurch eine Änderung erfährt; ist das nicht der Fall, steht der Befund in §7 und die Tests stehen unter dem Tag `docker` (laufen dann nur in `make test-docker`, weiter ohne Engine-Zugriff).
5. **Mutation:** die Hex-Prüfung entfernen und stdout und stderr mischen: mindestens ein Fall der Tests schlägt fehl; die Mutation wird danach zurückgenommen und der Befund in §7 festgehalten.
6. **Lokaler Nachweis:** `docker rmi docker.io/library/alpine:3` (schlägt das fehl, weil ein Container es nutzt, wird zuerst aufgeräumt), dann nur dieser Test mit der Rezeptur des Ziels `test-docker` (`docker run --rm --network=host … u-boot:test-docker-tools go test -tags docker -run TestE2E_SandboxDevcontainer_StatusFiles_InspectorAgainstEngine ./internal/e2e`, Docker-only, §3.1), danach `make test-docker` vollständig (rund zehn Minuten, `-timeout 30m`).

| Datei / Komponente | Änderungs-Art | Begründung (Fundstelle) |
|---|---|---|
| `internal/e2e/sandbox_status_docker_test.go` | update | `startLabeled` (Zeilen 64 bis 75): stdout getrennt, ID-Prüfung, `-q`; `docker pull -q docker.io/library/alpine:3` vor Zeile 351 in `TestE2E_SandboxDevcontainer_StatusFiles_InspectorAgainstEngine` (Zeile 340) |
| `internal/e2e/container_id_test.go` (untagged) | neu | reine Funktion `containerIDFromRunOutput` und ihre Tests auf Fake-Ausgaben; Fälle: nur ID, Pull-Zeilen plus ID, leer, Fehlerzeile plus ID, Ausgabe ohne ID — nach dem Befund oben |
| `internal/e2e/sandbox_devcontainer_docker_test.go` | unverändert | gelesen: `dockerOutput` (Zeile 40) liefert `CombinedOutput()` und hat viele Aufrufer, die gemischte Ausgabe erwarten (Fehlermeldungen, `strings.Contains`); die beiden `docker run -d --name`-Starts (Zeilen 225 und 573) setzen den Namen selbst und nutzen ein gebautes Image |
| `repo.mk`, `Dockerfile`, `.github/workflows/integration.yml` | unverändert | gelesen: `test-docker` führt `go test -tags docker -timeout 30m ./...` im Image `test-docker-tools` aus (`repo.mk`, Ziel `test-docker`); Gate- und CI-Änderungen sind ausgeschlossen (§1) |
| `harness/replay.md`, `docs/maintainer/quality.md` | unverändert | gelesen (Wortsuche nach `startLabeled`, `alpine`, `Image-Cache`, `dockerOutput`): beide Dateien nennen keinen dieser Begriffe, also beschreiben sie weder den Helfer noch Image-Cache-Voraussetzungen von `test-docker` |
| `docs/plan/planning/done/slice-v1-e2e-status-ci-helfer.md` | update | diese Tabelle (Dateiname der neuen Testdatei, Abgleich gegen `git diff --stat`) |
| `docs/plan/planning/in-progress/roadmap.md` | update | Sektion „Offene Wellen“: der Ruhe-Marker „Nichts in Arbeit.“ war durch den Verweis auf diesen Slice ersetzt, solange er in `in-progress/` lag (Befund von `make docs-check`, `planning-drift`); nach dem Move steht er wieder, die Zeile `v0.11.0` unter „Nächste Wellen“ nennt `slice-v1-podman-aktuelle-version` als Blocker des Cuts |
| `docs/plan/planning/open/README.md`, `docs/plan/planning/open/slice-v1-podman-aktuelle-version.md` | anderer Vorgang (`db4eb13`); in dieser Closure nur Zeilen und Verweise berichtigt | Commit `db4eb13` (Podman-Slice) liegt zwischen den Commits dieses Slice; `make slice-mv` zieht dessen Verweise auf diesen Slice nach `done/` nach, `open/README.md` trägt die berichtigten Zeilen (geliefert, startet jetzt, Cut wartet) |
| `docs/reviews/<Datum>-slice-v1-e2e-status-ci-helfer.md` | neu, durch den Reviewer | Review-Report |
| `docs/reviews/<Datum>-verification-slice-v1-e2e-status-ci-helfer.md` | neu, durch den Verifier | Verification-Report |
| `docs/plan/planning/observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/` | neu (Closure) | Beobachtung `observation.md`, `state.md`, `evidence/slice-v1-e2e-status-ci-helfer.md`; Zähler 1 |
| `docs/plan/planning/open/slice-v1-release-cut-v0.11.0.md`, `docs/plan/planning/open/README.md` | update (Closure) | Start-Trigger §4 und DoD um die beiden Blocker ergänzt; Zeilen des Podman-Slice und des Cuts berichtigt |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `next` → `in-progress`): Priorisierung durch den Projektinhaber; der Slice geht dem Cut [`slice-v1-release-cut-v0.11.0`](../open/slice-v1-release-cut-v0.11.0.md) voraus. WIP-Limit 1: kein anderer Slice liegt in `in-progress/`.

**Rückführungen — vorab benannt:**

- `in-progress` → `open` (Ursache nicht belegt): Das Nachstellen (Schritt 1) reproduziert den Fehler nicht; dann ist der Befund des Diagnose-Agenten unvollständig und der Plan braucht eine neue Diagnose, ehe Code geändert wird.
- `in-progress` → `next` (zu groß): Die Untersuchung weiterer Verwender (§1) ergibt Befunde außerhalb des Helfers; dann gehören sie in eigene Slices.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7). Der CI-Beleg (`integration-docker` grün auf dem Commit nach dem Push) ist **kein** Closure-Kriterium dieses Slice, weil der Push die Freigabe des Projektinhabers braucht; er steht als „weiter offen“ in §6 und ist Teil der DoD des Cuts.

## 6. Risiken und offene Punkte

Jedes Risiko bekommt bei der Closure genau einen Ausgang (eingetreten / entfallen / weiter offen).

- **Ursache nur abgeleitet:** Der Log nennt `page not found`, nicht den Pull-Text; die Kette „Pull-Fortschritt auf stderr → mehrzeiliger Name → 404“ stammt aus der Diagnose und der Nachstellung auf Kommandoebene. — **Ausgang:** offen; Mitigation: Schritt 1 stellt sie vor jeder Änderung selbst nach. — **Ausgang:** entfallen: belegt (Review: Vorher-Stand ohne `alpine:3` rot mit `seed files: … page not found`, Kommandoebene mit der CLI des Test-Images: Pull-Text nur auf stderr, ID auf stdout; Nachher-Stand grün).
- **Weitere Tests mit gemischter Ausgabe:** Ein Test, der eine ID, einen Namen oder einen Wert aus `dockerOutput` liest und nicht nur auf `Contains` prüft, kann denselben Fehler tragen. — **Ausgang:** offen; Mitigation: Wortsuche nach Verwendungen des Rückgabewerts in §1 (gelesen: die Ausgabe von `dockerOutput` fließt in Fehlermeldungen, in `strings.Contains` und in zwei Gleichheitsvergleiche, Zeile 253 in `sandbox_devcontainer_docker_test.go` (`docker exec … cat`) und Zeile 410 in `sandbox_status_docker_test.go` (`docker inspect`); beide betreffen bestehende Container ohne Pull), der Implementer wiederholt sie nach dem Fix. — **Ausgang:** entfallen: der Review fand im Paket `internal/e2e` und in den Tests von `internal/adapter/driven/docker` keine weitere Stelle, die eine Container-ID oder einen Namen aus gemischter Ausgabe eines Registry-Image-Starts liest.
- **Lokaler Nachweis ohne CI-Beleg:** Ein Colima-Lauf ohne Cache bildet den Runner nur nach (Netz, Docker-Version, Registry-Mirror, Rate-Limit des Docker Hub für anonyme Pulls). — **Ausgang:** offen; der CI-Lauf nach dem Push ist der Beleg und braucht die Freigabe des Projektinhabers; ein Rate-Limit-Fehler beim Pull wäre ein neuer Befund und kein Beleg gegen den Fix. — **Ausgang:** weiter offen: der CI-Lauf `integration-docker` nach dem Push ist Teil der DoD des Cuts [`slice-v1-release-cut-v0.11.0`](../open/slice-v1-release-cut-v0.11.0.md).
- **Unversionierte `docker-ce-cli`:** Das Verhalten der Pull-Ausgabe kann sich zwischen CLI-Versionen ändern; der Fix hängt nicht von ihrem Wortlaut ab (letzte Zeile, Hex-Prüfung), die Reproduzierbarkeit bleibt offen. — **Ausgang:** offen; Nebenbefund (§1), eigener Vorgang auf Entscheidung des Projektinhabers. — **Ausgang:** weiter offen (unverändert; der Hinweis auf einen Schritt `docker version` vor `make test-docker` steht in §7).
- **Pull-Fehler des Tests:** Scheitert `docker pull -q alpine:3` (Netz, Rate-Limit), endet der Test mit `t.Fatalf`; das ist so vorgesehen und kein Fehler des Helfers. — **Ausgang:** weiter offen (Review F-3, nur genannt); ein gemeldeter Lauf endete nach 30,05 s ohne gesicherte Meldung, Ursache unbelegt (§7, Flake-Beobachtung); nicht als eingetreten belegt.
- **Hygiene-Tests unter dem Tag:** Liegen die Tests des Helfers unter `docker`, laufen sie nicht in `make gates`. — **Ausgang:** offen; Schritt 4 entscheidet am Befund, die Entscheidung steht in §7. — **Ausgang:** entfallen: die Tests des Helfers sind untagged und laufen in `make test` (§7).

## 7. Closure-Notiz

**Belege:** Review-Report [`2026-10-18-slice-v1-e2e-status-ci-helfer`](../../../reviews/2026-10-18-slice-v1-e2e-status-ci-helfer.md) (0 HIGH, 0 MEDIUM, 1 LOW, 4 INFO; F-1 bis F-5 behoben oder als Risiko geführt) und Verification-Report [`2026-10-18-verification-slice-v1-e2e-status-ci-helfer`](../../../reviews/2026-10-18-verification-slice-v1-e2e-status-ci-helfer.md) (Liefer-Punkte der DoD erfüllt, Auflagen 1 bis 6; in dieser Closure erfüllt). Commits: `814072a` (Fix), `b440a1c` (Plan-Tabelle, Roadmap-Verweis), `fe4e5f5` und `342ac02` (Review-Fixes F-2, F-4, F-5), `cd6b0a2` (F-1, F-3, Review-Report) und der Closure-Commit; `db4eb13` gehört zum Podman-Slice, nicht zu diesem Vorgang.

- **Reproduktion (Review):** Vorher-Stand (`git archive 575281f`) mit entferntem `alpine:3`: **rot**, `sandbox_status_docker_test.go:362: seed files: Error response from daemon: page not found` (3,28 s), derselbe Fehler wie im CI-Lauf `39103939800`. Mit der CLI des Test-Images (Docker 29.9.1): `docker run -d alpine:3 sleep 5` schreibt auf stdout nur die ID, auf stderr den Pull-Text; mit `-d -q` stdout nur die ID, stderr leer. Nachher-Stand: ungecachtes `alpine:3` **grün** (7,51 s), gecacht grün (5,48 s). Die Ursachenkette aus §1 ist damit belegt.
- **Läufe der Verification:** 28 Einzelläufe von `…StatusFiles_InspectorAgainstEngine` ohne Fehlschlag (10 mit `-run StatusFiles`, 10 nur dieser Test, 8 unter vier gleichzeitigen Prozessen mit CPU-Last; 12 davon nach `docker rmi` des Images). Vollständiges `make test-docker`: `ok internal/e2e 635.727s` (Gesamtlauf 10 min 42 s, `-timeout 30m`), grün. `make gates` grün (Exit 0; `coverage-gate` 91,00 % bei Schwelle 90 %, `d-check` 0 Befunde), `make doc-trace` grün (83 Anforderungen, 0 Waisen), `make doc-immutable RANGE="575281f..HEAD"` grün.
- **Mutationsnachweis (Review, `make test` mit Mutation an `container_id_test.go`, danach zurückgenommen):** drei Mutationen erkannt (Hex-Länge 64 → 63: 4 Fälle rot; letzte → erste nichtleere Zeile: 3 Fälle rot; Großbuchstaben erlaubt: `upper-case hex` rot), zwei überlebten (ohne `strings.TrimSpace`; ohne CRLF-`ReplaceAll`). Behebung von F-2 in `fe4e5f5` und `342ac02`: Die Fälle „spaces and tab“, „CRLF line ending only“ und „two IDs“ sind ergänzt, die `TrimSpace`-Mutation ist damit rot; der CRLF-`ReplaceAll` war redundant (`TrimSpace` deckt `\r`, die Mutante war äquivalent) und ist entfernt. Der Verifier hat die Mutationen nicht wiederholt; `TestContainerIDFromRunOutput` hat 14 Fälle.
- **Fundstellen (`AGENTS.md` §3.10):** `internal/e2e/sandbox_status_docker_test.go`: `startLabeled` (`docker run -d -q --name <generiert>`, stdout und stderr in getrennten `strings.Builder`, `t.Cleanup` mit `rm -f -v` vor dem Start, nur stdout wird ausgewertet) und `TestE2E_SandboxDevcontainer_StatusFiles_InspectorAgainstEngine` (Vorab-Pull `docker pull -q docker.io/library/alpine:3` vor dem ersten `startLabeled`); `internal/e2e/container_id_test.go`: `containerIDFromRunOutput` (letzte nichtleere Zeile, `^[0-9a-f]{64}$`, Fehler mit der ganzen Ausgabe) und `TestContainerIDFromRunOutput`. `git diff --stat 575281f..HEAD` gleicht gegen §3 ab: jede berührte Datei steht in der Tabelle; die „unverändert“-Zeilen (`sandbox_devcontainer_docker_test.go`, `repo.mk`, `Dockerfile`, `integration.yml`, `replay.md`, `quality.md`) fehlen im Diff (vom Verifier gelesen).
- **Entscheidung „Tests des Helfers untagged“:** `make test` kompiliert und fährt `internal/e2e` ohne Engine (`ok … internal/e2e 0.006s`, Review), die Coverage-Schwelle bleibt erfüllt (91,00 %); die Tests stehen deshalb ohne Build-Tag in `container_id_test.go` und laufen in `make gates`. Risiko „Hygiene-Tests unter dem Tag“ entfällt.
- **Doku nicht berührt:** `harness/replay.md` und `docs/maintainer/quality.md` nennen weder `startLabeled`, `alpine`, `Image-Cache` noch `dockerOutput` (je 0 Treffer, vom Review selbst gesucht); keine öffentliche Schnittstelle berührt, kein Eintrag in `CHANGELOG.md`.
- **Flake-Beobachtung (Vorbehalt):** In einem gemeldeten Lauf endete `…InspectorAgainstEngine` nach 30,05 s ohne gesicherte Meldung. **Ursache unbelegt, nicht reproduzierbar** (28 Läufe danach ohne Fehlschlag); die Vermutung (Registry-Client-Zeitlimit des Daemons von 30 s beim neuen Pull, kurzzeitige Erreichbarkeit von Docker Hub) ist durch keine Meldung gestützt. Der Flake gilt **nicht als behoben**. **Beim nächsten Auftreten die volle Ausgabe sichern** (`-v`, Datei): `docker pull alpine:3: …` bedeutet Netz oder Registry, `seed files`, `FindContainer` oder `ReadFile` einen Fehler in Test oder Adapter. Er steht in [`observation.md`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md) unter „Benannt, nicht gezählt“.
- **Offene Punkte:** (1) **CI-Beleg `integration-docker` grün** auf dem Commit nach dem Push: braucht die Freigabe des Projektinhabers, kein Closure-Kriterium dieses Slice, **Teil der DoD des Cuts** [`slice-v1-release-cut-v0.11.0`](../open/slice-v1-release-cut-v0.11.0.md). (2) **Unversionierte `docker-ce-cli`** im Image `test-docker-tools` (`Dockerfile`, Stage `test-docker-tools`) und ein Schritt `docker version` vor `make test-docker` in `integration.yml`: Hinweis, Gate-/CI-Änderung (`AGENTS.md` §3.6), eigener Vorgang auf Entscheidung des Projektinhabers, nicht angelegt. (3) **Pull-Fehler des Tests:** scheitert `docker pull -q` (Netz, Rate-Limit), endet der Test mit `t.Fatalf` statt `t.Skip` (Review F-3); ein zweiter Versuch oder ein Skip wäre ein eigener Vorgang. (4) Acht fremde beendete Container auf dem Host (`friendly_pasteur`, `cool_colden`, `heuristic_robinson`, `boring_wilbur`, `schema-guard-dev-pg`, `cdc-examples-*`) blieben unangetastet; nur genannt.
- **Was hat funktioniert:** Der Review stellte den Fehler auf dem Vorher-Stand nach (das Image entfernen, den Test laufen lassen) und prüfte die Kommandoebene mit der CLI des Test-Images statt der des Hosts; damit war die abgeleitete Ursachenkette des Plans belegt. Die Mutationen fanden zwei Testlücken, die kein Gate gefangen hätte.
- **Was ging anders als geplant:** Die Plan-Tabelle führte `roadmap.md` und den fremden Vorgang `db4eb13` zunächst nicht (F-1); der Abgleich gegen `git diff --stat` gelang erst nach der Berichtigung. Zwei Mutanten überlebten die erste Testfassung (F-2); der Cleanup des Containers war erst nach der ID-Prüfung registriert (F-4).
- **Finding-Klassen des Reviews:** Plan-Tabelle gegen `diff --stat` unvollständig (F-1) · Mutant überlebt, redundante Normalisierung (F-2) · Netzabhängigkeit durch Vorab-Pull (F-3) · Cleanup erst nach ID-Prüfung registriert (F-4) · Kommentar ungenau zum Ist-Zustand (F-5).
- **Steering-Loop-Eintrag (nichts verkörpert, gezählt):** (1) Klasse „Test grün nur mit lokalem Cache“ (der Fehler zeigte sich erst auf dem Runner): erster Vorgang, neue Beobachtung [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md), Zähler 1, Stand offen. (2) „Plan-Tabelle gegen `diff --stat` unvollständig“ (F-1) ist die Klasse von [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (verkörpert, `AGENTS.md` §3.10); der Vorgang ist dort nicht als eigene Evidence angelegt, weil die Lücke vor der Closure berichtigt wurde und ein fremder Vorgang (`db4eb13`) im Range lag, kein Plan-Fehler über Code. **Entscheidung nicht getroffen:** ob dort ein weiteres Mal gezählt wird, ist offen. (3) Die übrigen Klassen (F-2 bis F-5) einmal, ohne Register-Eintrag. Kein `liegt in`: mit diesem Slice wird nichts verkörpert. `BEO-E2E/ci-ueberspringt-nested-podman-build` ist **unverändert offen** (Trigger: erster CI-Lauf nach dem Push; dieser Slice berührt die Build-Tests nicht).
- **Folge-Slices:** keine aus diesem Slice. [`slice-v1-podman-aktuelle-version`](../open/slice-v1-podman-aktuelle-version.md) startet jetzt und ist wie dieser Slice Blocker des Cuts.
- **Roadmap und Cut:** Der Ruhe-Marker „Nichts in Arbeit.“ steht nach dem Move wieder unter „Offene Wellen“. Die Zeile `v0.11.0` unter „Nächste Wellen“ nennt beide Blocker; der Start-Trigger und die DoD von `slice-v1-release-cut-v0.11.0` tragen sie (dieser Slice erfüllt, der Podman-Slice offen, dazu der grüne CI-Beleg `integration-docker` nach dem Push).
- **Drei Paarungen:** Anker (kein `liegt in`-Feld, nichts verkörpert) · Folge-Slice (keiner) · Register (ein Pfad, `evidence/` nicht leer) — nach dem `git mv` nach `done/` geprüft.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist `internal/e2e` (E2E, Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006)); `docs/plan/` (DOK, Greenfield) für Plan und Closure.

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register ([`../observations/`](../observations/README.md)) durchgegangen — [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md) betrifft andere Tests desselben Pakets (Skip auf dem Runner), nicht diesen Helfer; [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (verkörpert, `AGENTS.md` §3.10): jede Aussage dieses Plans über Code steht mit Fundstelle in §1 und §3, die Ursachenkette ist als „abgeleitet“ gekennzeichnet; keine weitere Beobachtung trifft zu.

Alle berührten Sub-Areas GF.
