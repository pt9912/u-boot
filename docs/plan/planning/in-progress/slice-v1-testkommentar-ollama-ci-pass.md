# Slice v1: Docker-Test — der Kommentar vor `TestE2E_SandboxDevcontainer_Ollama` nennt die Kopplung an `-v` statt eines erledigten Zustands (`LH-FA-DEV-003`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle, Release-Zuordnung: keine (reiner Kommentar im Test; kein Eintrag im CHANGELOG nötig). Der Slice hat keine Closure-Bedingung, die von seiner DoD verschieden wäre. Er ist der Folge-Schritt (b) und der Verifier-Befund V-4 von [`slice-cr-lastenheft-messstaende`](../done/slice-cr-lastenheft-messstaende.md) (§7, *Folge-Schritte*; [`2026-10-25-slice-cr-lastenheft-messstaende-verify`](../../../reviews/2026-10-25-slice-cr-lastenheft-messstaende-verify.md)).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Katalog-Feature `ollama`; der Test ist sein Docker-Test), [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) (Entscheidung 10 nennt den CI-Aufruf ohne `-v`; die ADR ist `Accepted` und bleibt unverändert, [`AGENTS.md`](../../../../AGENTS.md) §3.5), [`AGENTS.md`](../../../../AGENTS.md) §3.7 und §3.9 (Kommentar beschreibt, was da ist).
Anlass: Der Kommentar vor `TestE2E_SandboxDevcontainer_Ollama` (`internal/e2e/sandbox_devcontainer_docker_test.go`, Zeilen 1246 bis 1253) schließt mit „the amd64 run is not a verification of the architecture claim until a CI run with a visible pass.“ Der CI-Lauf mit sichtbarem Pass liegt vor (Run 40164437737, 2026-10-25; [`slice-cr-lastenheft-messstaende`](../done/slice-cr-lastenheft-messstaende.md), Beleg M5); der Satz beschreibt einen Zustand, der erledigt ist.

**Berührte Spec-Stellen:** — (reiner Kommentar im Test; kein Satz der Spezifikation oder der Sicht ist berührt; [`SPEC-047`](../../../../spec/spezifikation.md#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen) verweist für die Zusage-Menge auf das Lastenheft und bleibt gültig).

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-25.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Kommentar vor `TestE2E_SandboxDevcontainer_Ollama` sagt, was gilt (Zusage des Tests, Kopplung an `-v` im Aufruf von `test-docker`), und nennt keinen erledigten Zustand.

### Gefundener Stand (gelesen am 2026-10-25)

- **Der Kommentar:** `internal/e2e/sandbox_devcontainer_docker_test.go`, Doc-Kommentar vor `TestE2E_SandboxDevcontainer_Ollama` (Zeilen 1246 bis 1253). Er trägt zwei Teile. Der erste (Zusage) nennt, was der Test belegt: der erzeugte Dockerfile baut; das Image trägt das gepinnte Ollama ohne GPU-Bibliotheken und `zstd`; der Server hört hinter der Egress-Restriktion nur auf der Loopback-Adresse (die Container-Adresse lehnt ab), läuft auf der CPU, zieht ein Modell über die zwei Hosts der Default-Allowlist, ein nicht gelisteter Host bleibt gesperrt; das Modell überlebt das Entfernen des Containers. Der zweite Teil „The test needs network access and about 1.5 GB of download“ ist eine Grenze (Voraussetzung des Laufs). Der Schlusssatz „the amd64 run is not a verification of the architecture claim until a CI run with a visible pass“ ist ein Zustandsmarker mit Ablaufpunkt: Er sagt „noch nicht“, und das „bis“ ist eingetreten.
- **Der Test belegt, was der Kommentar zusagt:** Fundstellen im Testrumpf: `ollamaImageScript`-Marker `OLLAMA_VERSION_PINNED`, `OLLAMA_LIB_PRESENT`, `NO_GPU_LIBRARIES`, `NO_ZSTD` (Zeile 1269), Server-Script mit `LOOPBACK_ANSWERS`, `CONTAINER_ADDRESS_REFUSED`, `LIBRARY_CPU`, `PULL_OK`, `FOREIGN_BLOCKED` (Zeile 1307), `MODEL_SURVIVES_CONTAINER_REMOVAL` nach `docker rm` (Zeile 1324). Der Test ruft nur `requireDocker(t)` auf; er hat keine Architektur-Bedingung, und der Pass ist derselbe Test auf `amd64` und `arm64`.
- **Der CI-Beleg:** [`docs/messungen/slice-cr-lastenheft-messstaende/README.md`](../../../messungen/slice-cr-lastenheft-messstaende/README.md): Run 40164437737, Job 121082471598, `--- PASS: TestE2E_SandboxDevcontainer_Ollama (42.97s)`, Paket `internal/e2e` `ok`; Aufruf `make test-docker` mit Ausgabe `-v`. Dort *Lesehilfen*: das Log trägt für den Pass nur die Zeile `--- PASS`.
- **Der Aufruf:** `repo.mk`, Ziel `test-docker`, Zeile 139: `go test -v -tags docker -timeout 30m ./...`; `.github/workflows/integration.yml` Zeile 47 ruft `make test-docker` auf. Ohne `-v` zeigt `go test` für ein grünes Paket nur die Paketzeile `ok`; ein Pass und ein Skip des einzelnen Tests sind dort nicht zu trennen.
- **[`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) Entscheidung 10** (Zeile 80): „der CI-Aufruf läuft heute ohne `-v`, ein Skip ist von einem Pass dort nicht zu trennen“ — das war der Stand beim Schreiben der ADR; heute trägt `repo.mk` das `-v`. Die ADR bleibt unverändert ([`AGENTS.md`](../../../../AGENTS.md) §3.5); sie ist hier gelesen, nicht Gegenstand.
- **Der alte Wortlaut im Baum** (`grep -rn "visible pass"`, gelesen am 2026-10-25): Treffer im Test (Zeile 1253), in zwei `done/`-Slices als Zitat ([`slice-v1-devcontainer-ollama-amd64-messung`](../done/slice-v1-devcontainer-ollama-amd64-messung.md) §1, [`slice-cr-lastenheft-messstaende`](../done/slice-cr-lastenheft-messstaende.md) §1 und §7), im Verifier-Report (V-4) und im CHANGELOG-Eintrag einer ausgelieferten Version (Zeile 190, anderer Wortlaut: „run with a visible pass on that host“, Messaussage zum Host, nicht der Kommentar). Ein Treffer in Handbuch, README, Spezifikation, `harness/` oder `AGENTS.md` besteht nicht. `done/`-Zitate, Reports und die ausgelieferte CHANGELOG-Zeile sind Historie und werden nicht berührt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Test-Logik, Assertions, Testhelfer, andere Tests derselben Datei** — der Slice ist ein reiner Kommentar; ein Verhaltens-Diff wäre ein anderer Vorgang. Ob andere Kommentare der Datei einen erledigten Zustand tragen, ist per `grep -n -i "CI run\|until a\|not a verification" internal/e2e/*.go` gelesen: kein weiterer Treffer; eine weitergehende Durchsicht der Datei ist nicht Teil des Slice.
- **`repo.mk`, `.github/workflows/`, `-v` im Aufruf** — Gate-Haltung ([`AGENTS.md`](../../../../AGENTS.md) §3.6); der Slice liest den Aufruf und nennt ihn im Kommentar, ändert ihn nicht.
- **[`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md)** — `Accepted`, immutable; Entscheidung 10 beschreibt den Stand der Entscheidung. Eine Folgeentscheidung ist nicht nötig: die Zusage ist im Lastenheft (Zusage-Menge, [`slice-cr-lastenheft-messstaende`](../done/slice-cr-lastenheft-messstaende.md)) bereits auf den CI-Beleg gestellt.
- **Lastenheft, Spezifikation, Handbuch, README, CHANGELOG** — tragen den Satz nicht (siehe *Alter Wortlaut im Baum*). Zeile 416 des Lastenhefts (Folge-Schritt (a) desselben Slice) ist ein anderer Vorgang (Lastenheft-CR).
- **`podman`-Tests ohne CI-Gate auf `amd64`** — Risiko-Ausgang des Vorgängers ([`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md)); der Ollama-Test betrifft es nicht.

## 2. Definition of Done

- [ ] **Der Kommentar vor `TestE2E_SandboxDevcontainer_Ollama` ist neu gefasst:** Er behält die Zusage des Tests und die Grenze (Netz, rund 1,5 GB Download), streicht den Schlusssatz „the amd64 run is not a verification … until a CI run with a visible pass“ und nennt stattdessen die Kopplung: der Pass erscheint in der CI-Ausgabe nur, solange der Aufruf von `test-docker` in `repo.mk` `-v` trägt. Er enthält keine Chronik und keinen erledigten Zustand („until“, „not yet“, „former“, „früher“ u. ä.; [`AGENTS.md`](../../../../AGENTS.md) §3.7, §3.9 (a)). Der Diff ändert nur Kommentarzeilen (`git diff` zeigt keine Zeile, die kein `//` trägt).
- [ ] **Alter Wortlaut im Baum geprüft** ([`AGENTS.md`](../../../../AGENTS.md) §3.9 (b)): `grep -rn "visible pass"` und `grep -rn "not a verification"` über das ganze Repo; Treffer außerhalb von `done/`, `docs/reviews/` und der ausgelieferten CHANGELOG-Zeile bestehen nicht mehr; ADR, Handbuch, README und Spezifikation sind gelesen, nicht geändert.
- [ ] `make gates` grün. Ein volles `make test-docker` läuft nicht (Kommentar-Diff; der Lauf kostet rund 15 Minuten und etwa 1,5 GB Download, ohne dass er etwas pinnt, was der Diff ändern kann). Die Datei trägt den Build-Tag `docker`; weder `.golangci.yml` noch `Makefile` setzen `-tags` (gelesen, kein Treffer), nur `repo.mk` Zeile 139 im Ziel `test-docker`. `make gates` kompiliert und lintet sie also nicht; der Diff-Selbstcheck trägt die Form, siehe §6 Risiko 1.

Die Closure-Pflichten (Review und Verification durch eigene Kontexte, Closure-Notiz mit Lerneintrag, Beobachtungs-Register, Risiko-Ausgänge, drei Paarungen) gelten wie in jedem Slice und stehen in §5.

## 3. Plan (vor Code)

Die Tabelle nennt jede Datei, die der Diff berührt ([`AGENTS.md`](../../../../AGENTS.md) §3.10); der Implementer gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab. Gelesen sind die Fundstellen in der letzten Spalte.

| Datei / Komponente | Änderungs-Art | Begründung / Fundstelle |
|---|---|---|
| `internal/e2e/sandbox_devcontainer_docker_test.go` | update (nur Kommentarzeilen 1246 bis 1253) | Doc-Kommentar vor `TestE2E_SandboxDevcontainer_Ollama`; [`AGENTS.md`](../../../../AGENTS.md) §3.7. Kein Test, der neu hinzukommt: der Test selbst ändert sich nicht, und ein Kommentar hat keinen Test (die Kopplung an `-v` ist eine Grenze des Aufrufs, keine Bedingung in Spec oder Handbuch, §3.10 letzter Punkt gilt nicht) |
| `docs/plan/planning/in-progress/slice-v1-testkommentar-ollama-ci-pass.md` | neu (dieser Plan), später update (DoD, §7), `git mv` nach `done/` | Lifecycle; Move getrennt vom Inhalt ([`AGENTS.md`](../../../../AGENTS.md) §3.3) |
| `docs/plan/planning/open/README.md` | update (Zeile im *Bestand*) | Tabelle führt jeden Slice; Zeile gelesen: Spalten *Datei* und *Gegenstand* |
| `docs/plan/planning/observations/BEO-CLI/kommentar-beschreibt-abwesendes/` | update bei der Closure (neue Datei `evidence/slice-v1-testkommentar-ollama-ci-pass.md`, `state.md`) | neunter Beleg; die Beobachtung ist verkörpert (`AGENTS.md §3.9`); bisher acht Dateien in `evidence/`, gezählt am 2026-10-25 |
| `docs/plan/planning/done/slice-cr-lastenheft-messstaende.md` | unverändert | §7 *Folge-Schritte* (b) benennt diesen Schritt; gelesen: Zeile 309; die Datei ist `done/` und trägt den Hinweis bereits, die Verweise auf diesen Slice zieht `make slice-mv` bei der Closure nach, falls nötig |
| `docs/plan/planning/in-progress/roadmap.md` | unverändert | Zeile 20 „Nichts in Arbeit.“; `open/` führt die Roadmap nicht (gelesen: Abschnitte *Offene Wellen*, *Nächste Wellen*); ein Slice ohne Welle erscheint dort nicht |
| `repo.mk`, `.github/workflows/integration.yml`, `Makefile`, `.golangci.yml`, `Dockerfile` | unverändert | `repo.mk` Zeile 139 (`go test -v -tags docker …`) und `integration.yml` Zeile 47 gelesen: der Kommentar nennt den Aufruf, der Diff ändert ihn nicht |
| `docs/plan/adr/0024-ollama-im-sandbox-devcontainer.md`, `spec/`, `docs/user/`, `README*.md`, `CHANGELOG.md`, `harness/`, `AGENTS.md` | unverändert | der alte Wortlaut steht dort nicht (`grep`, siehe *Gefundener Stand*); [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) Entscheidung 10 (Zeile 80) gelesen und bleibt |

## 4. Trigger

**Start** (`open` → `next` → `in-progress`): Priorisierung durch den Projektinhaber; keine Abhängigkeit, kein anderer Slice ist Voraussetzung. Der Slice ist klein genug, um zwischen zwei Wellen zu laufen.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): wenn das Lesen der Datei weitere Kommentare mit erledigtem Zustand findet, die mehr als eine Zeile ändern; dann nur die Fundstelle des Auftrags umsetzen und die weiteren als Folge-Kandidaten nennen.
- `in-progress` → `open` (blockiert): wenn `repo.mk` das `-v` nicht mehr trägt (Gate-Änderung beim Projektinhaber); dann ist die Kopplung im Kommentar neu zu lesen, nicht zu raten.

## 5. Closure-Trigger

DoD vollständig; `make gates` grün; Review (`.harness/skills/reviewer.md`) und Verification (`harness/verification.md`) in eigenen Kontexten, kein Self-Review; Closure-Notiz in §7 mit Lerneintrag; `git mv` nach `done/` als eigener Commit nach dem Inhalt ([`AGENTS.md`](../../../../AGENTS.md) §3.3 Fall 2).

Lerneintrag: Mögliche Antwort ist, dass eine Beobachtung gezählt wird (`BEO-CLI/kommentar-beschreibt-abwesendes`, neunter Beleg: Zustandsmarker mit Ablaufpunkt in einem Testkommentar, erst durch den Folge-Slice eines anderen Vorgangs gefunden, nicht durch den Selbstcheck §3.9) und dass nichts neu verkörpert wird. Ein neuer Sensor ist nicht zu erwarten (die Beobachtung ist ohne mechanischen Sensor entschieden); der Implementer prüft, ob die Muster aus §3.9 (a) einen Ablaufpunkt-Marker („until a …“) einschließen, und nennt eine Lücke, falls nicht.

## 6. Risiken und offene Punkte

- **1. Die Datei trägt den Build-Tag `docker`; `make gates` kompiliert und lintet sie nicht** (gelesen: kein `-tags` in `.golangci.yml` und `Makefile`). Ein Kommentar-Diff bricht den Build nicht, solange jede geänderte Zeile mit `//` beginnt; das prüft der Selbstcheck am `git diff`. Ein Formatfehler wäre nur mit dem Tag sichtbar; der Implementer nennt im Bericht, was lief (etwa `go vet`/`gofmt` im Test-Image oder ein auf diesen Test beschränktes `-run`) oder dass es nicht lief. — **Ausgang:** offen bis zur Umsetzung; der Bericht trägt die Antwort.
- **2. Die Kopplung an `-v` driftet:** Wird `-v` aus `repo.mk` entfernt, stimmt der Kommentar nicht mehr. Das ist die Eigenart einer Kopplung (Klasse des Kommentars: Kopplung, [`AGENTS.md`](../../../../AGENTS.md) §3.7); sie wird benannt und mit dem Ziel `test-docker` in `repo.mk` genannt, damit die Suche nach dem Namen den Kommentar findet (§3.9 (b)). — **Ausgang:** weiter offen: nicht als Beobachtung, sondern als Eigenschaft der Kopplung; der Retirement-Check gilt für `-v` in `repo.mk` ohnehin nicht (kein Herkunfts-Anker dort).
- **3. Der Kommentar behauptet zu viel über den Pass:** „Der Pass erscheint nur mit `-v`“ ist über `go test` belegt (Ausgabe grüner Pakete ohne `-v`: nur die Paketzeile), nicht über das Log von Run 40164437737 allein. Der Satz wird als Aussage über `go test` formuliert, nicht über das CI-Log. — **Ausgang:** entfällt mit der gelesenen Formulierung; der Reviewer prüft.
- **4. Der Verweis auf den CI-Beleg im Kommentar:** Ein Kommentar mit Lauf-ID und Datum wäre Chronik. Er trägt keine; der Beleg steht im Messprotokoll. — **Ausgang:** entfallen, wenn der Kommentar keine Run-ID nennt.

## 7. Closure-Notiz

Wird bei der Closure geschrieben (nach der Umsetzung; Review und Verification liegen vor). Vorbereitete Felder:

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —

Beleg-Tabelle, vom Planner vorgelegt (jede Zeile ist eine Fundstelle, die der Implementer vor dem Schreiben des Kommentars wieder liest; ändert sich der Wortlaut, wird die Zeile mit geändert):

| Satz im neuen Kommentar (Entwurf, Englisch) | Fundstelle |
|---|---|
| „the generated Dockerfile builds; the image carries the pinned Ollama without GPU libraries and `zstd`“ | Testrumpf: `ollamaImageScript`-Marker `OLLAMA_VERSION_PINNED`, `OLLAMA_LIB_PRESENT`, `NO_GPU_LIBRARIES`, `NO_ZSTD`; `buildProjectImage` |
| „behind the egress restriction the server the user starts listens on loopback only, runs on the CPU, pulls a model through the two default-allowlist hosts while an unlisted host stays blocked“ | `ollamaServerScript`-Marker `LOOPBACK_ANSWERS`, `CONTAINER_ADDRESS_REFUSED`, `LIBRARY_CPU`, `PULL_OK`, `FOREIGN_BLOCKED` |
| „the model survives the removal of the container“ | `ollamaRestartScript`, Marker `MODEL_SURVIVES_CONTAINER_REMOVAL`, nach `rm -f` des ersten Containers |
| „The test needs network access and about 1.5 GB of download.“ (unverändert) | Tests laden Archiv und Modell (Messprotokoll [`slice-v1-devcontainer-ollama-amd64-messung`](../done/slice-v1-devcontainer-ollama-amd64-messung.md): Archiv 1,44 GB); Timeout 25 Minuten im Rumpf |
| „The architecture is not asserted in the test; the run on a host shows the pass in the CI output only while `test-docker` in `repo.mk` calls `go test` with `-v`.“ (Richtung des neuen Satzes; Wortlaut legt der Implementer fest) | `repo.mk` Zeile 139 (`go test -v -tags docker -timeout 30m ./...`); `go test` ohne `-v` gibt für ein grünes Paket nur die Paketzeile aus; der Test selbst ruft nur `requireDocker(t)` |

Mögliche Fassung für den Schluss des Kommentars (nur Richtung; der Implementer prüft jede Aussage vor dem Eintrag selbst): „The test needs network access and about 1.5 GB of download. Its pass shows in the CI output only while `test-docker` in `repo.mk` runs `go test -v`.“

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird `internal/e2e` (E2E, Test-Infrastruktur; Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area)) und `docs/plan/` (DOK, Greenfield). Die Schwelle ≥ 2 von 3 Achsen ist für beide in der Modus-Deklaration entschieden; der Slice teilt keine Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (`docs/plan/planning/observations/`). Treffer:

- [`BEO-CLI/kommentar-beschreibt-abwesendes`](../observations/BEO-CLI/kommentar-beschreibt-abwesendes/observation.md) — verkörpert (`AGENTS.md §3.9`), acht Evidence-Dateien (gezählt am 2026-10-25); dieser Slice ist ein weiteres Auftreten (Zustandsmarker mit Ablaufpunkt in einem Testkommentar), das Muster (a) aus §3.9 deckt „früher“ und „no longer“, nicht ein „until“. Die Sub-Area der Beobachtung ist CLI; der Beleg liegt in E2E: der Implementer nennt das bei der Closure im Register, statt es zu glätten.
- [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) — verkörpert (`AGENTS.md §3.10`); relevant als Prüflinse für diese Plan-Tabelle (jede Fundstelle oben wurde gelesen); kein neuer Beleg erwartet.
- [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md) — gelesen, berührt den Slice nicht (Podman-Tests, nicht der Ollama-Test); kein Beleg.

Evidenz-/Diskrepanz-Risiko: niedrig. Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
