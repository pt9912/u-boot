# Slice Gate: Traceability-Matrix (`--trace`) für u-boot konfigurieren

**Status:** **abgeschlossen** (2026-09-29, **Delivery-Hash: `3851589`**).
Ergebnis: **79 Anforderungen, 0 Waisen** (advisory; `--require-complete`
bleibt aus).

**Welle:** noch keiner Welle zugeordnet (Wartungs-Kandidat in
[`roadmap.md`](../in-progress/roadmap.md) §Nächste Wellen) — Befund aus dem
Gate-Ausbau, bewusst **nach** `welle-gate-ausbau-v0.51` eingeplant.

**Bezug:** [`LH-FA-PROJDOCS-006`](../../../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell)
(Referenzmodell) und die handgepflegte Traceability-Matrix in
`spec/lastenheft.md` §13. **Achtung:** §13 liegt im **Vertrags-Stratum** — jede
Änderung daran ist ein Change Request mit Version-Bump und Historie-Zeile
([`LH-FA-PROJDOCS-002`](../../../../spec/lastenheft.md#lh-fa-projdocs-002--adr-format)-Muster).

**Autor:** pt9912. **Datum:** 2026-07-25.

---

## 1. Ziel

Die Frage „welche Anforderung ist durch kein Artefakt belegt?" maschinell
beantwortbar machen. Heute ist sie nur von Hand zu beantworten — beim
Durchsehen der offenen Punkte musste genau das gemacht werden.

`--trace` erzeugt die Matrix bereits; sie ist nur nie konfiguriert worden.
**Vorab gemessen (2026-07-25):** Zwei Zeilen Konfiguration genügen, um aus
79 „Waisen" **12** zu machen:

```yaml
trace:
  slices:
    dir: docs/plan/planning
    file-pattern: '^(slice-.+)\.md$'
```

Der Default erwartet Slice-Dateinamen der Form `slice-<NNN>-…`; u-boot nutzt
`slice-<phase>-<slug>`. Ein reiner Formatunterschied — die ADR-Spalte
funktionierte von Anfang an.

**Die verbleibenden 12 sind Traceability-Lücken, keine
Implementierungslücken:** [[`LH-FA-CLI-001`](../../../../spec/lastenheft.md#lh-fa-cli-001--cli-aufruf)](../../../../spec/lastenheft.md#lh-fa-cli-001--cli-aufruf) (CLI-Aufruf), `-002` (Hilfeausgabe),
`-003` (Versionsausgabe) sind demonstrierbar geliefert; sie stammen aus frühen
Skeleton-Slices, die ihre Kennungen nie benannt haben. Dasselbe Muster bei
[[`LH-FA-DOC-001`](../../../../spec/lastenheft.md#lh-fa-doc-001--compose-datei-erzeugen)](../../../../spec/lastenheft.md#lh-fa-doc-001--compose-datei-erzeugen)/`-003`/`-004`, [[`LH-QA-001`](../../../../spec/lastenheft.md#lh-qa-001--automatisierte-tests)](../../../../spec/lastenheft.md#lh-qa-001--automatisierte-tests)/`-002`,
[[`LH-FA-PROJDOCS-004`](../../../../spec/lastenheft.md#lh-fa-projdocs-004--archivierung)](../../../../spec/lastenheft.md#lh-fa-projdocs-004--archivierung),
[`LH-FA-BUILD-003`](../../../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik),
[`LH-FA-DEV-002`](../../../../spec/lastenheft.md#lh-fa-dev-002--vs-code-kompatibilität).

**Messkorrektur bei Ausführung (2026-09-29):** Nach gut einem Vierteljahr
lautet die Waisenliste nicht mehr 12, sondern **6** (CLI-002/003/004,
DOC-003/004, QA-002) — seither geschlossene Slices haben vier der zwölf
genuine belegt; [`LH-FA-CLI-001`](../../../../spec/lastenheft.md#lh-fa-cli-001--cli-aufruf), [`LH-FA-DOC-001`](../../../../spec/lastenheft.md#lh-fa-doc-001--compose-datei-erzeugen), [`LH-QA-001`](../../../../spec/lastenheft.md#lh-qa-001--automatisierte-tests) und
[`LH-FA-PROJDOCS-004`](../../../../spec/lastenheft.md#lh-fa-projdocs-004--archivierung) waren nur durch die **Selbstreferenz** dieses Plans
hier in §1 gedeckt, was als Evaluations-Auftrag gezählt hat. Die
Einzelfall-Bewertung unten behandelt deshalb alle zwölf Ursprungslücken.

## 2. Definition of Done

- [x] **Vollständiger Kennungs-Umfang entschieden und konfiguriert.** Die RTM
  sieht nur `LH-FA-*` und `LH-QA-*` — **79 von 139** Kennungen; das ist die
  getroffene Entscheidung, nicht ein Zufall des Default-Musters: `NFA`, `SA`,
  `DA`, `AK`, `ZB`, `PE`, `PÜ`, `LESE`, `ABG`, `RISK`, `MVP`, `OPEN` sind
  Kontext-, Hinweis- oder Abgrenzungs-Familien, die strukturell nie einen
  Slice belegen. Die Begründung lebt als Kommentar in
  [`.d-check.yml`](../../../../.d-check.yml) §trace und in `MR-005`.
- [x] **Umlaut-Fall gelöst:** das `ids`-Muster kennt jetzt `ÄÖÜ`
  (`LH(?:-[A-ZÄÖÜ0-9]+)+-\d{3}[A-Z]?`); die beiden §13-Tabellenzeilen zu
  [`LH-PÜ-001`](../../../../spec/lastenheft.md#lh-pü-001--grundfunktion)/`-002` sind verlinkt.
- [x] **Advisory vor Gate.** `--require-complete` bleibt **aus**. Nach
  Abarbeitung der Lückenliste wäre ein Gate jetzt zwar grün — aber
  jeder Spec-CR (z. B. [`slice-cr-sandbox-devcontainer`](../open/slice-cr-sandbox-devcontainer.md))
  gebiert neue Kennungen ohne Slice-Coverage; ein rotes Gate würde
  abgeschaltet statt befolgt. `make doc-trace` bleibt das informative Target.
- [x] **Die 12 Lücken einzeln bewertet** — Tabelle:

| Kennung | Antwort | Evidenz |
|---|---|---|
| [[`LH-FA-CLI-001`](../../../../spec/lastenheft.md#lh-fa-cli-001--cli-aufruf)](../../../../spec/lastenheft.md#lh-fa-cli-001--cli-aufruf) | (a) | `slice-m1-repo-skeleton` (CLI-Stub `main.go` mit `--help`/`--version`) |
| [`LH-FA-CLI-002`](../../../../spec/lastenheft.md#lh-fa-cli-002--hilfeausgabe) | (a) | `slice-m1-repo-skeleton` (`--help`) |
| [`LH-FA-CLI-003`](../../../../spec/lastenheft.md#lh-fa-cli-003--versionsausgabe) | (a) | `slice-m1-repo-skeleton` (`--version`) |
| [`LH-FA-CLI-004`](../../../../spec/lastenheft.md#lh-fa-cli-004--fehlerausgabe) | (a) | `slice-v1-cli-json-dry-run` (strukturiertes `diagnostics`/`exitCode`-Fehlermodell); Vorstufe Exit-Code-Mapping in `slice-m3-init-flow` |
| [[`LH-FA-DOC-001`](../../../../spec/lastenheft.md#lh-fa-doc-001--compose-datei-erzeugen)](../../../../spec/lastenheft.md#lh-fa-doc-001--compose-datei-erzeugen) | (a) | `slice-m3-init-flow` (compose.yaml im generierten Template-Satz) |
| [`LH-FA-DOC-003`](../../../../spec/lastenheft.md#lh-fa-doc-003--netzwerk) | (a) | `slice-m5-add-postgres` (`networks:` in der erzeugten Compose) |
| [`LH-FA-DOC-004`](../../../../spec/lastenheft.md#lh-fa-doc-004--volumes) | (a) | `slice-m5-add-postgres` (Top-Level-`volumes:` für stateful Services) |
| [[`LH-QA-001`](../../../../spec/lastenheft.md#lh-qa-001--automatisierte-tests)](../../../../spec/lastenheft.md#lh-qa-001--automatisierte-tests) | (a) | `slice-m6-docker-integrationstests` (Test-Infrastruktur) |
| [`LH-QA-002`](../../../../spec/lastenheft.md#lh-qa-002--testbare-akzeptanzkriterien) | (a) | `slice-m6-docker-integrationstests` (Acceptance-Test-Vehikel, [`LH-AK-002`](../../../../spec/lastenheft.md#lh-ak-002--postgresql-flow)-Flow) |
| [[`LH-FA-PROJDOCS-004`](../../../../spec/lastenheft.md#lh-fa-projdocs-004--archivierung)](../../../../spec/lastenheft.md#lh-fa-projdocs-004--archivierung) | (c) | `docs/archive/` existiert seit `5efcd93` (2026-06-02) — Vorphase ohne Slice-Anker; Beleg ist der Commit, die Deklaration dieser Slice |
| [`LH-FA-BUILD-003`](../../../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik) | (a) | `slice-m3-build-polish` (`GO_VERSION`/`GOLANGCI_LINT_VERSION`-ARGs, Digest-Pinning) |
| [`LH-FA-DEV-002`](../../../../spec/lastenheft.md#lh-fa-dev-002--vs-code-kompatibilität) | (a) | `slice-m3-init-flow` (Devcontainer-Erzeugung, VS-Code-kompatibel gebaut) |

  Alle (a)-Nachträge sind in den Zielslices als markierte Querverweis-
  Korrektur (mit Bezugs-Marker auf diesen Slice) eingetragen; die Slices
  beschreiben die gelieferte Fähigkeit bereits in ihrem eigenen Text — die
  Kennung verlinkt einen bestehenden Beleg, sie behauptet keinen neuen.
- [x] **§13-Kollision entschieden, nicht übergangen.** §13 bleibt unangetastet
  (Prioritäts-Index, wie in der Plan-Empfehlung vorgesehen); die generierte
  RTM ist die Belegsicht. Kein CR nötig — §13 trägt die Prioritäten (MVP/V1),
  die die RTM gerade nicht führt; beide Sichten haben einen eigenen Job.
- [x] **`MR-005` nachgezogen** (Trace-Konfiguration als Teil der Gate-Haltung).
- [x] `make docs-check` grün; `make doc-trace` liefert **79 Anforderungen,
  0 Waisen**.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`.d-check.yml`](../../../../.d-check.yml) `trace` | neu | `slices`-Block + `requirements.id-pattern` |
| [`.d-check.yml`](../../../../.d-check.yml) `ids` | update | Umlaut-Fall im Kennungs-Muster |
| `docs/plan/planning/done/*.md` | update falls (a) | nachträgliche Kennungs-Nennung, Querverweis-Korrektur |
| [`harness/conventions.md`](../../../../harness/conventions.md) `MR-005` | update | Trace-Konfiguration dokumentieren |

## 4. Trigger

Gefeuert: Befund beim Durchsehen der offenen Lastenheft-Punkte (2026-07-25).
Bewusst **nach** `welle-gate-ausbau-v0.51` und nach dem anstehenden
Release-Cut eingeplant — der Sicherheits-Fix im Runtime-Image drängt stärker.

## 5. Closure-Trigger

Umfang entschieden, 12 Lücken einzeln bewertet, §13-Frage beantwortet,
`make doc-trace` mit dokumentierter Waisenzahl, `make docs-check` grün,
Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- **Ein weites Muster erzeugt Rauschen statt Signal.** Alle 139 Kennungen in
  die RTM zu nehmen produziert Dutzende „Waisen", die keine sind
  (Lesehinweise, Abgrenzungen, Risiken). Die Familien-Auswahl ist die
  eigentliche inhaltliche Arbeit dieses Slice, nicht die Konfiguration.
- **Nachträgliche Kennungs-Nennung in `done/`-Slices ist eine Grenzfrage.**
  Einen fehlenden Link zu setzen ist klar eine Querverweis-Korrektur. Eine
  Kennung *neu einzufügen*, die dort nie stand, ist mehr als das — es
  behauptet rückwirkend einen Bezug. Je Fall zu prüfen, im Zweifel Variante
  (c) statt (a).
- **Zwei Matrizen sind eine Drift-Quelle.** Genau das Muster, das `targets`
  und `planning` bei anderen Doppelpflegen aufgedeckt haben. Deshalb ist die
  §13-Frage ein DoD-Punkt und keine Fußnote.
- **Kein Carveout erwartet.**

## 7. Closure-Notiz (nach `done/`)

- **Was hat funktioniert:** Die Vorab-Messung des Plans (zwei Zeilen
  Konfiguration) hat exakt gestimmt — der Formatunterschied war der einzige
  Konfigurations-Aufwand; die eigentliche Arbeit war wie vorgesehen die
  Familien-Entscheidung und die Einzelfall-Bewertung.
- **Was ging anders als geplant:** Die Lückenliste war von 12 auf 6
  geschrumpft, und vier der inzwischen gedeckten Kennungen waren nur durch
  die Selbstreferenz dieses Plans abgedeckt — ein Scanner-Artefakt, das die
  Auswertung zusätzlich erzwingt. Alle (a)-Fälle (elf von zwölf) waren echte
  Querverweis-Korrekturen auf bereits im Ziel-Slice beschriebene Lieferungen;
  nur [`LH-FA-PROJDOCS-004`](../../../../spec/lastenheft.md#lh-fa-projdocs-004--archivierung) fiel auf (c) (Vorphase-Commit ohne Slice-Anker).
- **Steering-Loop-Lerneintrag:** Die RTM deckt ein Scanner-Artefakt auf, das
  die `ids`-/`matrix`-Module betreffen kann: eine Kennung in einem Plan zählt
  als Coverage, auch wenn der Plan sie nur als Lücke benennt. Für
  Lücken-Dokumentation in Plänen ist das akzeptabel, solange die Schließung
  per Einzelfall-Bewertung nachkommt — ein Gate auf `--require-complete`
  würde diese Selbstreferenz-Schleife zementieren; das bleibt ein
  weiteres Argument für advisory-only.
- **Folg-Slices:** keine. `--require-complete`-Entscheidung ist bewusst
  offen gelassen und an einen künftigen Anlass (stabile Kennungs-Basis,
  automatisierte Spec-Freshness-Prüfung) gebunden.

## 8. Sub-Area-Modus-Begründung

Berührte Sub-Areas: *harness / Konventionen* und *spec / docs* — beide **GF**
nach [`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration. Sollte die §13-Frage zu einem CR führen, ist zusätzlich das
Vertrags-Stratum berührt — dann greift der Fußabdruck-Weg aus
[`slice-harness-lastenheft-historie-cr-fussabdruck`](../done/slice-harness-lastenheft-historie-cr-fussabdruck.md).
