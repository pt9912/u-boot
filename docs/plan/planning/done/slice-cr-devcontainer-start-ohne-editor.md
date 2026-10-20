# Slice CR: Devcontainer ohne Editor starten — Lastenheft-CR (neue Anforderung, Befehlsübersicht, `--quiet`, `command`-Enum, Abgrenzung zum Sandbox-Profil)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-devcontainer-start-ohne-editor`](../welle-devcontainer-start-ohne-editor.md). Erster von vier Folge-Slices von [`slice-v1-devcontainer-start-ohne-cli`](../done/slice-v1-devcontainer-start-ohne-cli.md) (Reihenfolge CR → Fundament → Befehl → Spezifikation, Muster [`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md)). Priorität V2; Release: das nächste Minor nach `v0.12.0` (die Nummer vergibt der Release-Cut).

**Bezug:** [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Abgrenzung zum Hilfs-Container), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (`--quiet`, Ausnahme 1), [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) (Enum `command`), [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) (Befehlsübersicht), [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) und [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Nachbarn, unverändert bzw. Muster der Zusage-Menge); [`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) (Träger der Entscheidung, muss `Accepted` sein).
Anlass: Entscheidungen des Projektinhabers vom 2026-10-20 im Eltern-Slice; der Wortlaut liegt dort in §7 (*Lastenheft-Wortlaut-Vorschlag*, Stellen 1 bis 6) und wird hier unverändert übernommen (vom Projektinhaber freigegeben).

**Berührte Spec-Stellen:** `lastenheft.md` §4.3 (neue Anforderung nach [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle), §16 (Historie), Versionsfeld. `SPEC-*` und `ARC-*` berührt der CR nicht (siehe den Spezifikations-Slice). Der Verweis zeigt aufwärts: Das Lastenheft nennt diesen Slice nie.

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Rolle Architect). **Datum:** 2026-10-20.

---

## 1. Ziel und Abgrenzung

**Ziel:** Das Lastenheft sagt den Befehl `u-boot devcontainer up` / `down` zu: eine neue Anforderung, die zwei Zeilen der Befehlsübersicht, die `--quiet`-Regel samt Ausweisung der Socket-Lockerung, den Wert `devcontainer` im Enum von `command` und die Abgrenzung des Hilfs-Containers vom erzeugten Ergebnis. Es ist das **Zielbild**: kein Messstand, kein Prozesszustand, keine Slice-Kennung. Die Zusage-Menge ist „unter Docker auf einem Linux-Host mit lokaler Engine (`amd64`)“; Colima steht nicht darin und kommt durch einen Folge-CR nach der Einzelmessung hinzu. Ob die Linux-Zusage trägt, entscheiden die Docker-Tests des Befehl-Slice, nicht dieser CR.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Produktcode, Port, Befehl, Tests** — der CR schreibt den Vertrag; das Zielbild geht dem Code voraus ([`slice-v1-devcontainer-start-fundament`](../done/slice-v1-devcontainer-start-fundament.md) und [`slice-v1-cli-devcontainer-up-down`](../open/slice-v1-cli-devcontainer-up-down.md) liefern ihn).
- **`spec/spezifikation.md` und `spec/architecture.md`** — eine Spezifikation nach dem Code ([`slice-v2-spec-devcontainer-start`](../open/slice-v2-spec-devcontainer-start.md)); der Wertebereich von `command` in [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) gehört dorthin, nicht in den CR.
- **Handbuch und `CHANGELOG.md`** — es ändert sich kein Verhalten; der Befehl kommt mit dem Befehl-Slice.
- **Eine Zusage über das Belegte hinaus** — keine Aussage zu Docker Desktop, Podman als Host-Engine, entfernter Engine, Windows; der Wortlaut nennt sie als nicht zugesagt. Zurücknehmen, was der Befehl-Slice widerlegt, ist ein Folge-CR.
- **Die Release-Nummer festlegen** — das Release ist das nächste Minor nach `v0.12.0`, die Nummer vergibt der Release-Cut; die Priorität V2 steht im Wortlaut (Stelle 2).

## 2. Definition of Done

- [x] **Wortlaut im Lastenheft** nach Stelle 1 bis 5 des Eltern-Slice (neue Anforderung [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Befehlsübersicht, `--quiet`, Enum, Abgrenzung), so freigegeben vom Projektinhaber; Zielbild ohne Prozesszustand und ohne Messstand. *Beleg:* Verification §3 und §4 (wortgleich, kein ADR-, Slice- oder Prozessverweis im Zusatz); Fundstellen in §7.
- [x] **Versionsfeld und Historie (§16)** auf die nächste freie Minor-Version, eine Zeile mit den Stellen; Priorität (`V2`) in der Anforderung und in der Befehlsübersicht gleich. *Beleg:* Versionsfeld `0.18.0`, Historie-Zeile `0.18.0`, beide Zeilen der Befehlsübersicht `V2` (Verification §2 Zeile 2).
- [x] **Freigabe des Projektinhabers** im Gespräch belegt und in §7 vermerkt (`BEO-DOK/vertragsaenderung-nach-review`: kein Wortlaut nach dem Review ohne Nachlauf-Review über den Diff). *Beleg:* §7, *Freigabe-Vermerk*.
- [x] `make gates` grün, `make docs-check` grün. *Beleg:* §7, *Gates*.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor: [`2026-10-20-slice-cr-devcontainer-start-ohne-editor`](../../../reviews/2026-10-20-slice-cr-devcontainer-start-ohne-editor.md), Verification [`2026-10-20-verify-slice-cr-devcontainer-start-ohne-editor`](../../../reviews/2026-10-20-verify-slice-cr-devcontainer-start-ohne-editor.md)
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: Roadmap-Zeile „In Arbeit“ bei der Closure; Handbuch und `CHANGELOG.md` unverändert (kein Verhalten).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben: zwei Evidence-Dateien (§7, *Beobachtungs-Register*).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen) (§6, §7 *Risiken aus §6*).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der Slice-Closure geprüft (§7). Closure-Commit `91a1fb3`.

## 3. Plan (vor Code)

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; sie wird nach der letzten Fix-Runde gegen `git diff --stat` abgeglichen (`AGENTS.md` §3.10).

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Stellen 1 bis 6 des Wortlauts. Gelesen am 2026-10-20: [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) ist die letzte Anforderung von §4.3 (die neue steht danach, Nummer 010); Befehlstabelle in [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) endet mit `egress list` (V2) und `template list` (V1); Versionsfeld `0.17.0`, letzte Historie-Zeile `0.17.0` |
| `docs/plan/planning/in-progress/roadmap.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Zeile „In Arbeit“ bei Start und bei der Closure |
| `docs/plan/planning/open/README.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Bestand (Zeile dieses Slice) |
| `docs/reviews/` <!-- d-check:ignore (Pfad-Zitat) --> | neu (Reviewer, Verifier); update (nur Verweise) | Review und Verifikation dieses Slice (zwei neue Dateien); die zwei Reports des Eltern-Slice bleiben inhaltlich, `make slice-mv` zieht je eine Zeile Verweise nach |
| `docs/plan/planning/done/slice-cr-devcontainer-start-ohne-editor.md` (diese Datei) | update | `Verantwortlich:`, Plan-Tabelle, Closure; Move per `make slice-mv` |
| `docs/plan/planning/done/slice-v1-devcontainer-start-ohne-cli.md`, `done/slice-v1-devcontainer-start-fundament.md`, `open/slice-v2-spec-devcontainer-start.md` <!-- d-check:ignore (Pfad-Zitate) --> | update (nur Verweise) | `make slice-mv` zieht den Pfad dieses Slice nach (`open/` → `in-progress/`, bei der Closure `in-progress/` → `done/`) |
| `docs/plan/planning/open/slice-v1-cli-devcontainer-up-down.md` <!-- d-check:ignore (Pfad-Zitat) --> | update (Verweis und Inhalt) | Verweis-Nachzug durch `make slice-mv`; **inhaltlich** in der Fix-Runde (Review LOW 2): Test S4 prüft den Marker im benannten Volume nach `down`, einem folgenden `up` und `--recreate` (gelesen im Diff `9e6933b..HEAD`) |
| `docs/plan/planning/welle-devcontainer-start-ohne-editor.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Verweise (`make slice-mv`); Versionsnummer `0.18.0` vergeben (Start); bei der Closure der Zustand des CR in der Reihenfolge (§5) und die Versionsnummer in den Freigaben |
| `docs/plan/planning/observations/BEO-DOK/{plan-aussage-ueber-code,spec-bedingung-ohne-code-entsprechung}/` <!-- d-check:ignore (Pfad-Zitat) --> | update (Closure) | je eine Evidence-Datei `slice-cr-devcontainer-start-ohne-editor.md` und ein Absatz in `state.md` |
| `docs/plan/planning/open/slice-cr-lastenheft-messstaende.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Versionsangaben „nächste freie Version nach `0.18.0`“ (Titel, Historie-Frage, Stand `0.18.0`), weil dieser Slice `0.18.0` belegt |
| `spec/spezifikation.md`, `spec/architecture.md`, `docs/user/benutzerhandbuch.md`, `CHANGELOG.md`, `internal/`, `cmd/` <!-- d-check:ignore (Pfad-Zitate) --> | **unverändert** | kein Code, kein Verhalten; [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) (Zeile `enum` von `command`) und [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) zieht der Spezifikations-Slice nach dem Code nach (Befund im [ADR](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 10) |

## 4. Trigger

**Start** (`open` → `in-progress`, ggf. über `next`): `Verantwortlich:` gesetzt; [`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) ist `Accepted`; der Projektinhaber hat den Wortlaut (Eltern-Slice §7), die Priorität und die Release-Zuordnung entschieden; kein anderer Slice in `in-progress/` (WIP-Limit 1). Der Slice ist unabhängig von [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md); wer zuerst läuft, bekommt die kleinere Versionsnummer.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Projektinhaber ändert die Zusage-Menge so, dass eine neue Messung nötig wäre (Entscheidung im ADR).
- `in-progress` → `open` (blockiert): die Freigabe des Wortlauts steht aus; die Evidence nennt den Grund.

## 5. Closure-Trigger

DoD vollständig, Review-Report und Verification-Evidence liegen vor, Freigabe des Projektinhabers belegt, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **Der Wortlaut sagt mehr zu als die Tests tragen** (Standardprofil mit UID ≠ 1000). — **Ausgang:** weiter offen; der Befehl-Slice [`slice-v1-cli-devcontainer-up-down`](../open/slice-v1-cli-devcontainer-up-down.md) trägt die Tests und die Einzelmessung, ein widerlegender Befund geht zurück an den Architect und führt zu einem Folge-CR ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)). Der Review ordnete jede Bedingung des Wortlauts einem Test zu; die eine ohne Test („`down` entfernt keine Volumes“) trägt nach der Fix-Runde S4.
- **Ein Wort trägt eine Wirkung an anderer Stelle** („weist aus“ bleibt unter `--quiet`). — **Ausgang:** entfallen; der Review prüfte jedes Verb des Wortlauts gegen Stelle 3: die Socket-Ausweisung steht in Ausnahme 1, die Fallback-Ausweisung in Ausnahme 2 von [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), beide bleiben unter `--quiet` (Review, *Prüfung nach Auftrag* 3; Verification §3).
- **Der Wortlaut ändert sich nach dem Review.** — **Ausgang:** entfallen (nicht eingetreten); `git diff 69280df..HEAD -- spec/` ist leer, der Review las den Stand `69280df`; ein Nachlauf-Review ist nicht ausgelöst ([`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md)).
- **Versionsnummer kollidiert mit [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md).** — **Ausgang:** entfallen; `0.18.0` war frei (letzte Historie-Zeile `0.17.0`) und ist vergeben, der Messstände-CR nennt „nächste freie Version nach `0.18.0`“ (Review, *Prüfung nach Auftrag* 4).

## 7. Closure-Notiz

**Geliefert.** Der Lastenheft-Wortlaut Stellen 1 bis 6 des Eltern-Slice ([`slice-v1-devcontainer-start-ohne-cli`](../done/slice-v1-devcontainer-start-ohne-cli.md) §7), unverändert übernommen, im Lastenheft `0.18.0` (Commit `69280df`). Fundstellen in `spec/lastenheft.md`:

| Stelle | Fundstelle |
|---|---|
| 1 neue Anforderung | [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) in §4.3, nach [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), `*Priorität: V2.*` |
| 2 Befehlsübersicht | [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle), Zeilen `u-boot devcontainer up` und `u-boot devcontainer down`, beide `V2` |
| 3 `--quiet` | [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging): Liste der ändernden Subkommandos und Ausnahme 1 (Einbindung des Container-Runtime-Sockets in den Hilfs-Container) |
| 4 Enum `command` | [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run): Wert `devcontainer` zwischen `egress` und `template` |
| 5 Abgrenzung | [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil): Spiegelstrich zum Container-Runtime-Socket |
| 6 Version und Historie | Versionsfeld `0.18.0`, §16 Zeile `0.18.0` |

Kein Produktcode, keine Spezifikation, kein Handbuch, kein `CHANGELOG.md`.

**Freigabe-Vermerk (Projektinhaber, im Gespräch am 2026-10-20, vermerkt im Eltern-Slice [`slice-v1-devcontainer-start-ohne-cli`](../done/slice-v1-devcontainer-start-ohne-cli.md) §7, *Entscheidungsvermerk* 1 bis 4).** (1) Der Wortlaut, Stellen 1 bis 6, ist freigegeben, unverändert; das Lastenheft trägt ihn wortgleich (Verification §3). (2) Priorität **V2**; Release das nächste Minor nach `v0.12.0`, die Nummer vergibt der Release-Cut. (3) Colima steht im Wortlaut nicht; ein Folge-CR nimmt es nach der Einzelmessung des Befehl-Slice auf. (4) [`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) und [`ADR-0029`](../../adr/0029-devcontainer-start-exit-11-und-socket-haltung.md) sind `Accepted`. Die Versionsnummer `0.18.0` ist die nächste freie zum Zeitpunkt des CR.

**Closure-Trigger (§5), je Punkt.**

| Punkt | Beleg |
|---|---|
| DoD vollständig | §2; offen bleiben allein die Zeilen, die erst nach dem Move gelten: Doku-Update (Roadmap-Zeile „In Arbeit“) und Closure-Commit |
| Review-Report | [`2026-10-20-slice-cr-devcontainer-start-ohne-editor`](../../../reviews/2026-10-20-slice-cr-devcontainer-start-ohne-editor.md) (kein HIGH, kein MEDIUM, zwei LOW, vier INFO) |
| Verification-Evidence | [`2026-10-20-verify-slice-cr-devcontainer-start-ohne-editor`](../../../reviews/2026-10-20-verify-slice-cr-devcontainer-start-ohne-editor.md) (GO; V-1 LOW und V-2 LOW behoben, V-3 bis V-5 INFO) |
| Freigabe des Projektinhabers | *Freigabe-Vermerk* oben |
| Closure-Notiz mit Steering-Loop-Lerneintrag | dieser Abschnitt |

**Gates.** `make docs-check` und `make gates` grün auf dem Stand dieser Closure-Inhalte (Lauf des Planners, Ergebnis im Handback). Verification auf `2981c43`: `make gates` Exit 0, `make doc-trace` 84 Anforderungen, 0 Waisen (Zeile [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) `ok`), `make doc-immutable RANGE="9e6933b..HEAD"` ohne Befund. Nicht gelaufen: `make test-docker` und `make verify-depguard` (kein Code-, Docker- oder Importregel-Pfad im Diff).

- **Was hat funktioniert:** Die Reihenfolge Freigabe → Wortlaut → Review → Verification, ohne dass der Wortlaut nach dem Review noch einmal geändert wurde: Der Review las den Stand `69280df` (Wortlaut Stelle für Stelle gegen den freigegebenen Vorschlag, Konsistenz im Lastenheft, jedes Verb gegen Stelle 3, jede Bedingung gegen einen Test des Befehl- und Fundament-Slice), die Fix-Runde `2981c43` berührte nur Planungstext, und `git diff 69280df..HEAD -- spec/` blieb leer; ein Nachlauf-Review war damit nicht ausgelöst. Der Review fand die einzige Bedingung ohne Test („`down` entfernt weder den Container noch seine Volumes“) und ließ sie im Testplan des Befehl-Slice (S4) schließen, bevor der Code entsteht.
- **Was ging anders als geplant:** (1) Die Plan-Tabelle (§3) nannte beim Review vier Zeilen für zwölf berührte Dateien, und nach der Fix-Runde führte sie `open/slice-v1-cli-devcontainer-up-down.md` als „nur Verweise“, obwohl dort Test S4 inhaltlich geändert wurde (Review LOW 1, Verification V-1); berichtigt in dieser Closure. (2) Der Review-Report trug einen zerschossenen Link-Text, den `make docs-check` nicht meldet (Verification V-2, im Commit `7f74f8c` behoben). (3) Die ADR-Spalte von [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) in `make doc-trace` bleibt leer, weil die Kennung bei `Accepted` von [`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) nicht bestand und die ADR nicht umgeschrieben wird; die Zeile ist über den Slice-Beleg `ok` (Review INFO, Verification V-4).
- **Abgleich der Plan-Tabelle (§3) gegen `git diff --stat`** (Stand HEAD `7f74f8c` plus die Dateien dieser Closure, nach der letzten Fix-Runde; `git diff --stat 9e6933b` ohne Move): 16 Dateien im Diff gegen die Tabelle, dazu zwei neue Evidence-Dateien. Jede Zeile hat einen Eintrag: `spec/lastenheft.md`; `roadmap.md` (Start; die Closure-Zeile „In Arbeit“ folgt nach dem Move); `open/README.md`; die vier Dateien unter `docs/reviews/` (zwei neue, zwei Linknachzüge); dieser Slice; `done/slice-v1-devcontainer-start-ohne-cli.md`, Fundament- und Spezifikations-Slice (nur Verweise); `open/slice-v1-cli-devcontainer-up-down.md` (Verweis und Test S4, Zeile berichtigt); `open/slice-cr-lastenheft-messstaende.md`; die Welle; die zwei `state.md` und die zwei Evidence-Dateien des Registers. Ohne Diff: `spec/spezifikation.md`, `spec/architecture.md`, Handbuch, `CHANGELOG.md`, `internal/`, `cmd/` (wie in der Tabelle „unverändert“; im Slice-Diff `9e6933b..HEAD` unberührt, Verification §5). Der Eltern-Slice trägt in §7 *Lastenheft-Wortlaut-Vorschlag* weiter die Aussage, `0.18.0` sei für den Messstände-CR vorgesehen; sie ist die Chronik des Vorschlags in einem geschlossenen Slice und bleibt.
- **Steering-Loop-Eintrag:** Lerneintrag als **gezählte Beobachtungen, keine geschärfte Regel verkörpert, kein neuer Sensor, keine Spec-Lücke.** (1) [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md): dreißigster Vorgang (Plan-Tabelle ohne Zeilen für berührte Dateien, Zeile „nur Verweise“ trotz inhaltlicher Änderung nach der Fix-Runde; Review LOW 1, Verification V-1). (2) [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md): neunter Vorgang (Halbsatz des geschriebenen Lastenheft-Textes ohne Testzuordnung nach `down`; Review LOW 2, im Testplan des Befehl-Slice behoben). (3) [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md): **kein Beleg, kein Auftreten.** Der Wortlaut entstand vor dem Review und blieb danach unverändert (`git diff 69280df..HEAD -- spec/` leer); die Beobachtung zählt Vorgänge, in denen er sich nach dem Review änderte, und bleibt bei fünf Belegen. Ob ein Vorgang ohne Auftreten im Register zu führen ist, kennt die Beobachtungsform nicht; es ist keine Evidence-Datei angelegt. (4) [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md): kein neuer Beleg; der Review fand keine Aussage über die Zusage-Menge hinaus, das Risiko trägt der Befehl-Slice. (5) [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md): kein neuer Beleg; die DoD nannte die Befehlsübersicht ausdrücklich, und beide Zeilen stehen in [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle). Verkörpert wurde mit diesem Slice nichts; die Einträge sind gezählt, nicht verkörpert.
- **Beobachtungs-Register (`../observations/`):** `evidence/slice-cr-devcontainer-start-ohne-editor.md` in `BEO-DOK/plan-aussage-ueber-code/` (Zähler steht damit bei 30×) und in `BEO-DOK/spec-bedingung-ohne-code-entsprechung/` (9×) ergänzt, je mit einem Absatz in `state.md`. Keine neue Beobachtung angelegt; `vertragsaenderung-nach-review`, `zusage-weiter-als-messung` und `vorab-suche-uebersieht-befehlsuebersicht` unverändert (kein Auftreten).
- **Folge-Slices:** [`slice-v1-devcontainer-start-fundament`](../done/slice-v1-devcontainer-start-fundament.md), [`slice-v1-cli-devcontainer-up-down`](../open/slice-v1-cli-devcontainer-up-down.md) und [`slice-v2-spec-devcontainer-start`](../open/slice-v2-spec-devcontainer-start.md) sind Dateien in `open/`. Neuer Slice: keiner. Der Folge-CR für Colima ist keine Datei; er entsteht nach der Einzelmessung des Befehl-Slice.
- **Risiken aus §6:** alle vier tragen einen Ausgang. Entfallen: „Wort trägt eine Wirkung an anderer Stelle“ (Review und Verification haben jedes Verb gegen Stelle 3 gelesen), „Wortlaut ändert sich nach dem Review“ (nicht eingetreten, `spec/` seit `69280df` unverändert), „Versionsnummer kollidiert“ (`0.18.0` vergeben, Messstände-CR nachgezogen). Weiter offen mit Anker: „Wortlaut sagt mehr zu als die Tests tragen“ im Befehl-Slice [`slice-v1-cli-devcontainer-up-down`](../open/slice-v1-cli-devcontainer-up-down.md) (Docker-Tests S1, S3, S4, S8, S9 und Einzelmessung Colima).
- **Zur Entscheidung vorgelegt (Projektinhaber), nicht entschieden:** (a) Teilsensor „Dateiliste der Plan-Tabelle gegen `git diff --stat`“ zu `BEO-DOK/plan-aussage-ueber-code` (30 Vorgänge, wiederholt vorgelegt). (b) Rückfall von `BEO-DOK/spec-bedingung-ohne-code-entsprechung` (9 Vorgänge): ob [`AGENTS.md`](../../../../AGENTS.md) §3.10 um den Lastenheft-Text im Plan erweitert wird (so schon im Eltern-Slice vorgelegt). (c) Ob ein Vorgang, der eine offene Beobachtung nicht wiederholt (hier `vertragsaenderung-nach-review`: Wortlaut vor dem Review, danach unverändert), im Register geführt wird und wie; die Beobachtungsform kennt nur Auftreten. Ausgänge werden hier nicht zugewiesen.
- **Drei Paarungen:** Anker: kein Feld `liegt in` in diesem Slice (nichts verkörpert), nichts zu paaren. Folge-Slice: die drei Dateien in `open/` existieren. Register: die beiden Verzeichnisse oben haben je eine neue Evidence-Datei; die übrigen Verzeichnisse sind unverändert.
- **Closure-Commits:** Inhalt (dieser Abschnitt, DoD, Plan-Tabelle, Welle, `open/README.md`, Register), danach der reine `git mv` (`make slice-mv`), danach der Nachzug: Roadmap-Zeile „In Arbeit“ (Ruhe-Marker „Nichts in Arbeit.“), die Zeile „Doku-Update“ der DoD und der Closure-Commit-Hash in der letzten DoD-Zeile (`91a1fb3`).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird `spec/` (DOK, Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area)); die Schwelle aus zwei von drei Achsen ist dort erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md) (Vertragstext), [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (Zusage-Menge, „weist aus“ unter `--quiet`), [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md) (ein neuer Befehl gehört in die Befehlsübersicht, die Stichwort-Suche findet sie nicht; deshalb nennt die DoD die Befehlsübersicht ausdrücklich). Evidenz-/Diskrepanz-Risiko: mittel.

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
