# Slice v1: Harness-Hygiene — gofmt-Bestand und Formatter-Gate, `verify-depguard` Fall `domain-isoliert`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — kein Release-Bezug festgelegt; der Slice läuft **nach** dem Cut [`slice-v1-release-cut-v0.11.0`](../done/slice-v1-release-cut-v0.11.0.md) und ist nicht Teil von `v0.11.0` (Entscheidung des Projektinhabers 2026-10-18, „nach belieben“). Herkunft: Nebenbefund des Reviews zu [`slice-cr-sandbox-fuse-startfehler`](../done/slice-cr-sandbox-fuse-startfehler.md) (F-4, Closure-Notiz §7: viele Bestandsdateien nicht `gofmt`-konform, kein Formatter im Gate) und Befund F-7 des Reviews zu [`slice-v1-sandbox-fallback-host-ausweisung`](../done/slice-v1-sandbox-fallback-host-ausweisung.md) (`make verify-depguard`, Fall `domain-isoliert`).

**Bezug:** [`LH-QA-004`](../../../../spec/lastenheft.md#lh-qa-004--linting-solid-nahes-lint-profil) (Lint-Profil),
[`ADR-0003`](../../adr/0003-solid-nahes-lint-profil.md) (Lint-Profil, `depguard`).
Eine **Verschärfung eines Gates, kein Lockern**: [`AGENTS.md`](../../../../AGENTS.md) §3.6 verlangt für eine Abschwächung einen Plan- oder ADR-Anker; für eine Verschärfung ist dieser Slice der Plan-Anker.

**Berührte Spec-Stellen:** keine geplant (`docs/maintainer/quality.md` und `harness/sensors/verify-depguard.md` sind Maintainer-Doku bzw. Sensor-Vertrag, keine Spec-Straten).

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Go-Quellbestand ist formatiert und das Lint-Gate hält ihn formatiert (Teil A); `make verify-depguard` beweist alle neun `depguard`-Regeln, auch `domain-isoliert` (Teil B).

### Entscheidung des Projektinhabers (2026-10-18)

„Nach belieben“: Beide Teile laufen als ein Harness-Slice nach `v0.11.0`; **nicht** Teil von `v0.11.0`.

### Gefundener Stand (Ist-Stand, gelesen am 2026-10-18)

**Teil A — Formatierung.**

- `.golangci.yml` führt `version: "2"` und aktiviert unter `linters.enable` fünf Default-Linter und 24 SOLID-nahe Linter; **kein** `formatters`-Abschnitt (gelesen, ganze Datei). Das Lint-Ziel ist `golangci-lint run ./...` in der Docker-Stage `lint` (`Dockerfile`, `FROM golangci/golangci-lint:${GOLANGCI_LINT_VERSION} AS lint`, `RUN golangci-lint run ./...`; Ziel `lint` in `repo.mk`).
- Der Reviewer fand Bestandsdateien, die nicht `gofmt`-konform sind (genannt: `internal/adapter/driving/cli/cli.go`, `addservice.go`, `git_test.go`; Closure-Notiz von [`slice-cr-sandbox-fuse-startfehler`](../done/slice-cr-sandbox-fuse-startfehler.md) §7). **Gemessen im Docker** (`golangci-lint fmt --diff ./...`, Image `u-boot:lint`, [`AGENTS.md`](../../../../AGENTS.md) §3.1): 29 Dateien, Liste im Anhang A.
- `docs/maintainer/quality.md` nennt „5 Defaults + 24 SOLID-nahe Linter“ und „Damit sind 29 Linter aktiv“ (§1, vor der reservierten §1.1) und führt die SOLID-nahen Linter in §1.2; ein Formatter ist dort nicht aufgeführt.
- Geprüft an der gepinnten Version `v2.14.0` (Lauf im Docker): Das Schema `version: "2"` aktiviert Formatter unter `formatters.enable`; `golangci-lint fmt ./...` formatiert den Bestand, `golangci-lint run` meldet unformatierte Dateien.

**Teil B — `verify-depguard`.**

- Die Tabelle `CASES` in `scripts/verify-depguard.sh` führt neun Fälle; der erste ist `domain-isoliert|internal/hexagon/domain|github.com/pt9912/u-boot/internal/hexagon/port/driven|domain must not depend on port`. Das Skript wählt das verbotene Ziel so, dass kein Go-Importzyklus entsteht, weil `depguard` erst nach dem Typecheck läuft (Kommentar im Skript über `CASES`).
- Der Fall `domain-isoliert` scheitert: `make lint` endet mit „import cycle not allowed … port/driven → domain → port/driven (typecheck)“, nicht mit „domain must not depend on port“ (Review F-7, am Basiscommit `3b740aa` selbst reproduziert; im Bestand). Ursache: `port/driven` importiert `domain` (etwa `template_catalog.go`), und jedes verbotene Ziel der Regel importiert `domain` ebenfalls: gelesen sind `internal/adapter/driven/clock/clock.go:9` (importiert `port/driven`) und `internal/hexagon/port/driving/config.go` (importiert `domain`); `domain` selbst importiert nichts aus `internal/` außerhalb von `domain` (Suche nach `internal/` in `internal/hexagon/domain`, ohne Treffer außerhalb des eigenen Pakets).
- Folge: Der Sensor kann die Regel `domain-isoliert` nicht beweisen und endet nie grün; `make lint` bleibt für die Regel selbst grün (`harness/README.md` §Sensors, Zeile `make verify-depguard`; `harness/sensors/verify-depguard.md` beschreibt den Vertrag „alle neun Regeln“).
- Das Skript verlangt einen sauberen Arbeitsbaum in den Schicht-Verzeichnissen (Exit `2` sonst) und braucht je Regel einen vollen `make lint`-Lauf (im Kopfkommentar „3–5 Minuten“ geschätzt).

### Lösungsraum Teil B — der Implementer wählt und misst

Ein Import, der **keinen Zyklus** bildet, setzt ein verbotenes Ziel voraus, das `domain` nicht (transitiv) importiert. Unter den gelesenen Paketen erfüllt keines die Bedingung; eine Möglichkeit ist ein vom Skript angelegtes und wieder entferntes Wegwerf-Paket ohne Importe unter `internal/adapter/` (das das Präfix der Regel trifft), eine andere ein Paket, das der Implementer beim Lesen findet. **Nicht gemessen, ob diese Lösung trägt;** der Implementer prüft sie mit `make verify-depguard` (alle neun Fälle `PASS`) und hält im Plan-Eintrag fest, welche Lösung er gewählt hat und warum.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Lockerung oder Entfernung einer Regel, eines Linters oder einer Schwelle:** [`AGENTS.md`](../../../../AGENTS.md) §3.6; der Slice verschärft (Formatter) und beweist (Sensor), er lockert nichts.
- **Eine Änderung der `depguard`-Regeln selbst** (`.golangci.yml`, Abschnitt `depguard`): Die Regeln stimmen; der Befund liegt im Beweis-Skript.
- **Ein Formatter über `gofmt` hinaus** (`gofumpt`, `goimports`, `gci`): Der Auftrag nennt `gofmt`; jeder weitere Formatter wäre eine zusätzliche Verschärfung mit eigener Begründung.
- **Eine Änderung am Produktverhalten:** Das Formatieren ändert kein Verhalten (Test-Lauf grün als Beleg); der Slice berührt Code nur als Formatierung.
- **Ein Make-Target für den Formatierungslauf, falls der Implementer eines vorschlägt:** ein neues Target erscheint im Gate-Index (`harness/README.md`, d-check `targets`) und ist eine eigene Entscheidung; der Plan-Eintrag des Implementers nennt es, der Projektinhaber entscheidet.

## 2. Definition of Done

- [ ] **Teil A — Formatierung und Gate:** Der Bestand ist `gofmt`-konform (der Lauf des Formatters im Docker ändert danach nichts mehr); `.golangci.yml` aktiviert `gofmt` als Formatter (Verschärfung, kein Lockern); `make lint` schlägt bei einer nicht formatierten Datei an (am Beispiel einer absichtlich unformatierten Datei gesehen und wieder entfernt); der Formatierungs-Commit ist von der Gate-Änderung getrennt, damit der Diff des Bestands keine Logik enthält (§3.3 von [`AGENTS.md`](../../../../AGENTS.md)) ([`LH-QA-004`](../../../../spec/lastenheft.md#lh-qa-004--linting-solid-nahes-lint-profil), [`ADR-0003`](../../adr/0003-solid-nahes-lint-profil.md)).
- [ ] **Teil B — Sensor:** `scripts/verify-depguard.sh` beweist den Fall `domain-isoliert` mit einem Import, der keinen Zyklus bildet; `make verify-depguard` endet mit allen neun Regeln `PASS` und Exit `0`; der Vertrag in `harness/sensors/verify-depguard.md` stimmt mit dem Skript überein.
- [ ] Doku-Update: `docs/maintainer/quality.md` nennt den Formatter (Zahl der aktiven Prüfungen in §1, Tabelle der Linter in §1.2) und `.golangci.yml` trägt ihn im Kopfkommentar.
- [ ] `make gates` grün (inklusive des neuen Formatter-Gates); `make verify-depguard` grün; `make doc-trace` ohne neue Waise.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne Wellen-Betrieb hier geprüft.

## 3. Plan (vor Code)

Jede Aussage über Code steht mit Fundstelle (gelesen am 2026-10-18). Der Implementer erweitert die Tabelle um jede Datei, die der Diff berührt (der Formatierungs-Commit berührt viele Dateien; er nennt sie als Gruppe mit der gemessenen Zahl), gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab ([`AGENTS.md`](../../../../AGENTS.md) §3.10) und belegt eine Datei, die er als „unverändert“ führt, durch eine gelesene Aussage.

| Datei / Komponente | Änderungs-Art | Begründung (Fundstelle) |
|---|---|---|
| Formatierung: 29 Go-Dateien (Liste im Anhang A) | update (nur Formatierung) | Commit `d94da25`, vor der Gate-Änderung; `git diff --stat`: 170 Einfügungen, 177 Entfernungen, nur Whitespace, Leerzeilen und ein Umbruch (`git diff -w`: 6 Einfügungen, 13 Entfernungen, davon Leerzeilen, ein umgebrochenes Closure-Literal in `internal/e2e/sandbox_devcontainer_docker_test.go` und die vorhandene `//nolint:unparam`-Direktive am Doc-Kommentar von `dependenciesFor` in `internal/hexagon/application/addservice.go`, die `gofmt` ans Ende des Kommentarblocks hinter eine `//`-Zeile setzt; der Satzrest der Begründung wird im Nachzug-Commit der Review-Runde zurückgeholt, siehe nächste Zeile); Beleg: `make test` und `go vet -tags docker ./...` (Image `u-boot:test-docker-tools`) grün, `make gates` grün |
| `internal/hexagon/application/addservice.go` | update (nur Kommentar) | Review F-1: Begründung und Direktive `//nolint:unparam` von `dependenciesFor` stehen in einer Zeile; Wirkung der Direktive unverändert (`make lint` grün) |
| `docs/reviews/2026-10-18-slice-v1-harness-gofmt-und-depguard.md` | neu | Review-Report, mit diesem Nachzug committet |
| `.golangci.yml` | update | Commit `2265909`: Abschnitt `formatters.enable: [gofmt]` (Schema v2 der gepinnten Version `v2.14.0`, `Dockerfile` `ARG GOLANGCI_LINT_VERSION`; `golangci-lint fmt` und der Lauf von `run` kennen den Formatter, gelesen am Lauf); Kopfkommentar um den Formatter ergänzt; kein bestehender Eintrag entfernt oder gelockert |
| `docs/maintainer/quality.md` | update | §1, Absatz nach der Linter-Tabelle: nennt den Formatter, die Fehlermeldung und den Docker-Aufruf `golangci-lint fmt ./...` im Image `u-boot:lint` (probiert: Exit 0, keine Änderung am sauberen Baum); die Zahl 29 bleibt, der Formatter ist kein Linter; §1.2 unverändert (Tabelle der 24 SOLID-nahen Linter, ein Formatter gehört nicht hinein) |
| `scripts/verify-depguard.sh` | update | Commit `a6deaaf`: `CASES`-Zeile `domain-isoliert` mit optionalem fünften Feld (Wegwerf-Paket `internal/hexagon/port/verifydepguardprobe`, keine Importe); Anlegen vor, Entfernen nach dem Lauf, Aufräumen im `trap`; Kommentar über `CASES` nennt die Begründung |
| `harness/sensors/verify-depguard.md` | update | Absatz unter „Vertrag“ zum Wegwerf-Paket; Exit `2` nennt das vorhandene Probe-Verzeichnis |
| `harness/README.md` | update | Zeile `make lint` im Sensors-Index nennt den Formatter `gofmt`; kein Target ist hinzugekommen oder entfallen |
| `.git-blame-ignore-revs` | neu | enthält `d94da2561469a7bdbec2d0f8b05fdee3315d37af`, damit `git blame` den Format-Commit überspringt (GitHub liest die Datei von selbst, lokal `git config blame.ignoreRevsFile .git-blame-ignore-revs`); Hash erst nach dem Format-Commit bekannt, daher eigener Folge-Commit |
| `docs/plan/planning/in-progress/roadmap.md` | update | Ruhe-Marker „Nichts in Arbeit.“ durch die Arbeitszeile mit Slice-Link ersetzt |
| dieser Slice | update | Plan §3 und Anhang A |
| `repo.mk`, `Dockerfile` | unverändert | gelesen: `lint` ruft `docker build --target lint`, die Stage ruft `golangci-lint run ./...`; der Formatter in `.golangci.yml` greift ohne weitere Option (`make lint` am Beispiel einer unformatierten Datei rot gesehen) |

**Gewählter Weg Teil A:** `golangci-lint fmt ./...` im Image `u-boot:lint` (die Version `v2.14.0`, die auch das Gate `make lint` fährt); damit sind Formatierer und Prüfer derselbe Binär-Stand und der Lauf Docker-only. Ein separates `gofmt -w` im Tools-Image wäre ein zweiter Stand ohne Bindung an den Pin. Der Formatter `gofmt` des Linters läuft mit den Defaults von `golangci-lint` (inklusive `simplify`). Ein Make-Target für den Lauf wird nicht vorgeschlagen; der Aufruf steht in `docs/maintainer/quality.md` §1. Beobachtung: Eine Datei, die zugleich einen `unused`-Befund trägt, zeigte in der ersten Probe nur diesen, nicht den Formatter-Befund; die Probe mit nur einem Formatierungsfehler zeigt `File is not properly formatted (gofmt)` und Exit 2.

**Nachzug aus dem Review (F-1, F-3, F-4):** Der Kommentarblock von `dependenciesFor` in `internal/hexagon/application/addservice.go` trägt Begründung und `//nolint:unparam`-Direktive wieder in einer Zeile (nur Kommentar). Das Probe-Paket des Skripts trägt einen Paketkommentar, damit `revive` nicht anschlägt; im Lauf des Falls `domain-isoliert` ist die `depguard`-Meldung der einzige Befund. `docs/maintainer/quality.md` §1 nennt Nutzerangabe im Container und `simplify` (gelesen am Lauf: `for i, _ := range` wird zu `for i := range`).

**Hinweis aus dem Review, nicht Gegenstand dieses Slice (F-2, Bestand vor `38f17a5`):** Im Quellbaum stehen neun Inline-`nolint`-Kommentare, obwohl [`AGENTS.md`](../../../../AGENTS.md) §3.2 sie ausschließt, und `nolintlint` ist nicht aktiviert. Wirksam (`//nolint:`): `dependenciesFor` in `internal/hexagon/application/addservice.go` (`unparam`); `pickHostPort` in `internal/hexagon/application/upservice_portprobe_docker_test.go` (`gosec`); das Feld `lastCtx` in `fakeCatalog` in `internal/hexagon/application/template_list_test.go` (`containedctx`); je ein `dupl` an `TestUpJSON_FSError_SanitizesPathInDiagnosticMessage` (`internal/adapter/driving/cli/up_acceptance_test.go`), `TestDownJSON_FSError_SanitizesPathInDiagnosticMessage` (`down_acceptance_test.go`) und `TestLogsJSON_FSError_SanitizesPathInDiagnosticMessage` (`logs_acceptance_test.go`). Wirkungslos, weil mit Leerzeichen (`// nolint:dupl`): `mapUpErrorToDiagnostic` in `up.go`, `mapDownErrorToDiagnostic` in `down.go`, `mapLogsErrorToDiagnostic` in `logs.go` (alle in `internal/adapter/driving/cli/`). Der Planner führt es bei der Closure ins Register oder als offenen Posten.

**Gewählter Weg Teil B:** Wegwerf-Paket unter dem verbotenen Präfix `internal/hexagon/port` ohne Importe. Messung: vorher `make verify-depguard` 8 von 9 (`FAIL domain-isoliert: expected desc not found`, Exit 2 von `make` bei Skript-Exit 1), nachher 9 von 9, Exit 0. Aufräumen bei Abbruch: das Skript mit `SIGTERM` nach vier Sekunden beendet, danach fehlen Probe-Paket und Verletzungsdatei, `git status` sauber. Keine `depguard`-Regel geändert.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): der Cut [`slice-v1-release-cut-v0.11.0`](../done/slice-v1-release-cut-v0.11.0.md) liegt in `done/` (Entscheidung des Projektinhabers 2026-10-18: nach `v0.11.0`); kein anderer Slice liegt in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): zwei Liefer-Punkte sind die Obergrenze; zeigt die Messung des Bestands einen sehr großen Formatierungs-Diff oder braucht Teil B einen größeren Umbau des Skripts, wird in zwei Slices geschnitten (A, B).
- `in-progress` → `open` (blockiert): Der Formatter der gepinnten Version von `golangci-lint` verhält sich anders als angenommen (Schema, Lauf im Docker); dann klärt der Implementer es an der Dokumentation, und wenn es eine Pin-Anhebung braucht, ist das ein eigener Vorgang.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7). Wird mit diesem Slice eine Regel verkörpert (etwa ein Satz in `AGENTS.md` zum Formatter), trägt die Closure-Notiz das Feld `liegt in` und das Ziel den Anker `seit slice-v1-harness-gofmt-und-depguard`.

## 6. Risiken und offene Punkte

Jedes Risiko bekommt bei der Closure genau einen Ausgang (eingetreten / entfallen / weiter offen).

- **Der Formatierungs-Diff ist groß und verdeckt eine Logikänderung:** Der Bestand ist nicht gemessen. — **Ausgang:** offen; Mitigation: eigener Commit nur für die Formatierung, `make test` vor und nach dem Commit, der Reviewer liest den Diff auf reine Formatierung.
- **Der Formatter ändert Importgruppen oder Kommentare so, dass ein Test oder ein Sensor bricht** (etwa ein Test, der Quelltext liest): nicht gelesen. — **Ausgang:** offen; Mitigation: `make gates` nach der Formatierung.
- **Das Formatter-Gate bricht Zweig-Arbeit anderer Slices** (nicht formatierte neue Dateien): eine Verschärfung wirkt sofort. — **Ausgang:** offen; Mitigation: der Slice läuft allein (WIP-Limit 1), die Gate-Änderung steht nach dem Formatierungs-Commit.
- **Der Wegwerf-Import für `domain-isoliert` bildet doch einen Zyklus oder trifft die Regel nicht:** nicht gemessen. — **Ausgang:** offen; Mitigation: `make verify-depguard` zeigt es, andernfalls wählt der Implementer eine andere Lösung (§1).
- **Das Skript hinterlässt ein Wegwerf-Paket bei Abbruch:** Das Skript hat ein `trap` für die Verletzungsdatei (`cleanup`); ein neues Paket braucht denselben Schutz. — **Ausgang:** offen; Mitigation: Aufräumen im `trap`, Test durch Abbruch eines Laufs.

## 7. Closure-Notiz

Wird bei der Closure gefüllt.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `tools/`, `scripts/` (WZG), `harness/`, `docs/plan/` (DOK) und mit der Formatierung alle Go-Sub-Areas (DOM, APP, PORT, CLI, ADP, E2E; nur Formatierung) — alle Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006), jede erfüllt die Schwelle ≥ 2 von 3 Achsen; keine Ausdifferenzierung nötig. `cmd/uboot` ist keine Sub-Area (Aspirantin).

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register ([`../observations/`](../observations/README.md)) durchgegangen — `BEO-WZG/kopplung-nicht-nachgezogen` (`offen`, Zähler 1: ein Skript folgt einer geänderten Form nicht) ist verwandt mit Teil B (das Beweis-Skript folgt der Import-Lage nicht); der Slice trägt dazu bei Bedarf einen Beleg ein. `BEO-DOK/plan-aussage-ueber-code` (verkörpert, `AGENTS.md` §3.10) gilt für die Plan-Tabelle.

Alle berührten Sub-Areas GF.

## Anhang A: Formatierte Dateien (Commit `d94da25`, 29)

- `internal/adapter/driven/git/git_test.go`
- `internal/adapter/driven/runtime/runtime_test.go`
- `internal/adapter/driven/yaml/mapping_patch.go`
- `internal/adapter/driving/cli/cli.go`
- `internal/adapter/driving/cli/erroremission.go`
- `internal/adapter/driving/cli/jsontestutil/coderegistry.go`
- `internal/adapter/driving/cli/jsontestutil/jsontestutil.go`
- `internal/e2e/sandbox_devcontainer_docker_test.go`
- `internal/hexagon/application/acceptance_test.go`
- `internal/hexagon/application/addservice.go`
- `internal/hexagon/application/addservice_factory_test.go`
- `internal/hexagon/application/addservice_t5_test.go`
- `internal/hexagon/application/addservice_test.go`
- `internal/hexagon/application/config.go`
- `internal/hexagon/application/config_test.go`
- `internal/hexagon/application/devcontainer_features_test.go`
- `internal/hexagon/application/doctor.go`
- `internal/hexagon/application/generate.go`
- `internal/hexagon/application/generate_sandbox_test.go`
- `internal/hexagon/application/generate_test.go`
- `internal/hexagon/application/initproject.go`
- `internal/hexagon/application/initproject_runbackup_test.go`
- `internal/hexagon/application/initproject_test.go`
- `internal/hexagon/application/keycloak_catalogue_test.go`
- `internal/hexagon/application/noop_test.go`
- `internal/hexagon/application/template_init_test.go`
- `internal/hexagon/application/template_list_test.go`
- `internal/hexagon/domain/projectname_test.go`
- `internal/hexagon/domain/servicename_test.go`
