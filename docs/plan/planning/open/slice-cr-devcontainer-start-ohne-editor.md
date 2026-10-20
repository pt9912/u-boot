# Slice CR: Devcontainer ohne Editor starten — Lastenheft-CR (neue Anforderung, Befehlsübersicht, `--quiet`, `command`-Enum, Abgrenzung zum Sandbox-Profil)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Erster von vier Folge-Slices von [`slice-v1-devcontainer-start-ohne-cli`](../in-progress/slice-v1-devcontainer-start-ohne-cli.md) (Reihenfolge CR → Fundament → Befehl → Spezifikation, Muster [`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md)). Release-Zuordnung offen (Vorschlag des Planners: V2, nächstes Minor nach dem Befehl-Slice; entscheidet der Projektinhaber).

**Bezug:** [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Abgrenzung zum Hilfs-Container), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (`--quiet`, Ausnahme 1), [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) (Enum `command`), [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) (Befehlsübersicht), [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) und [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Nachbarn, unverändert bzw. Muster der Zusage-Menge); [`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) (Träger der Entscheidung, muss `Accepted` sein).
Anlass: Entscheidungen des Projektinhabers vom 2026-10-20 im Eltern-Slice; der Wortlaut liegt dort in §7 (*Lastenheft-Wortlaut-Vorschlag*, Stellen 1 bis 6) und wird hier unverändert übernommen, soweit der Projektinhaber ihn freigibt.

**Berührte Spec-Stellen:** `lastenheft.md` §4.3 (neue Anforderung nach [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle), §16 (Historie), Versionsfeld. `SPEC-*` und `ARC-*` berührt der CR nicht (siehe den Spezifikations-Slice). Der Verweis zeigt aufwärts: Das Lastenheft nennt diesen Slice nie.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Architect). **Datum:** 2026-10-20.

---

## 1. Ziel und Abgrenzung

**Ziel:** Das Lastenheft sagt den Befehl `u-boot devcontainer up` / `down` zu: eine neue Anforderung, die zwei Zeilen der Befehlsübersicht, die `--quiet`-Regel samt Ausweisung der Socket-Lockerung, den Wert `devcontainer` im Enum von `command` und die Abgrenzung des Hilfs-Containers vom erzeugten Ergebnis. Es ist das **Zielbild**: kein Messstand, kein Prozesszustand, keine Slice-Kennung. Die Zusage-Menge ist „unter Docker auf Linux und unter Colima (Linux-VM)“; ob sie trägt, entscheiden die Docker-Tests und die Einzelmessung des Befehl-Slice, nicht dieser CR.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Produktcode, Port, Befehl, Tests** — der CR schreibt den Vertrag; das Zielbild geht dem Code voraus ([`slice-v1-devcontainer-start-fundament`](../open/slice-v1-devcontainer-start-fundament.md) und [`slice-v1-cli-devcontainer-up-down`](../open/slice-v1-cli-devcontainer-up-down.md) liefern ihn).
- **`spec/spezifikation.md` und `spec/architecture.md`** — eine Spezifikation nach dem Code ([`slice-v2-spec-devcontainer-start`](../open/slice-v2-spec-devcontainer-start.md)); der Wertebereich von `command` in [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) gehört dorthin, nicht in den CR.
- **Handbuch und `CHANGELOG.md`** — es ändert sich kein Verhalten; der Befehl kommt mit dem Befehl-Slice.
- **Eine Zusage über das Belegte hinaus** — keine Aussage zu Docker Desktop, Podman als Host-Engine, entfernter Engine, Windows; der Wortlaut nennt sie als nicht zugesagt. Zurücknehmen, was der Befehl-Slice widerlegt, ist ein Folge-CR.
- **Priorität und Release-Nummer festlegen** — Entscheidung des Projektinhabers (Stelle 2 des Wortlauts).

## 2. Definition of Done

- [ ] **Wortlaut im Lastenheft** nach Stelle 1 bis 5 des Eltern-Slice (neue Anforderung `LH-FA-DEV-010` <!-- d-check:ignore (Kennung des Vorschlags) -->, Befehlsübersicht, `--quiet`, Enum, Abgrenzung), so freigegeben vom Projektinhaber; Zielbild ohne Prozesszustand und ohne Messstand.
- [ ] **Versionsfeld und Historie (§16)** auf die nächste freie Minor-Version, eine Zeile mit den Stellen; Priorität (Vorschlag `V2`) in der Anforderung und in der Befehlsübersicht gleich.
- [ ] **Freigabe des Projektinhabers** im Gespräch belegt und in §7 vermerkt (`BEO-DOK/vertragsaenderung-nach-review`: kein Wortlaut nach dem Review ohne Nachlauf-Review über den Diff).
- [ ] `make gates` grün, `make docs-check` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: Roadmap-Zeile „In Arbeit“ bei der Closure; Handbuch und `CHANGELOG.md` unverändert (kein Verhalten).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; sie wird nach der letzten Fix-Runde gegen `git diff --stat` abgeglichen (`AGENTS.md` §3.10).

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Stellen 1 bis 6 des Wortlauts. Gelesen am 2026-10-20: [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) ist die letzte Anforderung von §4.3 (die neue steht danach, Nummer 010); Befehlstabelle in [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) endet mit `egress list` (V2) und `template list` (V1); Versionsfeld `0.17.0`, letzte Historie-Zeile `0.17.0` |
| `docs/plan/planning/in-progress/roadmap.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Zeile „In Arbeit“ bei der Closure |
| `docs/plan/planning/open/README.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Bestand (Zeile dieses Slice) |
| `docs/reviews/` <!-- d-check:ignore (Pfad-Zitat) --> | neu (Reviewer, Verifier) | Review und Verifikation |
| `spec/spezifikation.md`, `spec/architecture.md`, `docs/user/benutzerhandbuch.md`, `CHANGELOG.md`, `internal/`, `cmd/` <!-- d-check:ignore (Pfad-Zitate) --> | **unverändert** | kein Code, kein Verhalten; [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) (Zeile `enum` von `command`) und [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) zieht der Spezifikations-Slice nach dem Code nach (Befund im [ADR](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 10) |

## 4. Trigger

**Start** (`open` → `in-progress`, ggf. über `next`): `Verantwortlich:` gesetzt; [`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) ist `Accepted`; der Projektinhaber hat den Wortlaut (Eltern-Slice §7), die Priorität und die Release-Zuordnung entschieden; kein anderer Slice in `in-progress/` (WIP-Limit 1). Der Slice ist unabhängig von [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md); wer zuerst läuft, bekommt die kleinere Versionsnummer.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Projektinhaber ändert die Zusage-Menge so, dass eine neue Messung nötig wäre (Entscheidung im ADR).
- `in-progress` → `open` (blockiert): die Freigabe des Wortlauts steht aus; die Evidence nennt den Grund.

## 5. Closure-Trigger

DoD vollständig, Review-Report und Verification-Evidence liegen vor, Freigabe des Projektinhabers belegt, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **Der Wortlaut sagt mehr zu als die Tests tragen** (Colima ohne CI, Standardprofil mit UID ≠ 1000). — **Ausgang:** weiter offen; der Befehl-Slice trägt die Tests und die Einzelmessung, ein widerlegender Befund geht zurück an den Architect und führt zu einem Folge-CR ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)).
- **Ein Wort trägt eine Wirkung an anderer Stelle** („weist aus“ bleibt unter `--quiet`). — **Ausgang:** weiter offen; der Reviewer prüft jedes Verb des Wortlauts gegen Stelle 3.
- **Der Wortlaut ändert sich nach dem Review.** — **Ausgang:** weiter offen; Nachlauf-Review über den Diff jeder Änderung ([`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md)).
- **Versionsnummer kollidiert mit [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md).** — **Ausgang:** weiter offen; die Nummer ist die nächste freie zum Zeitpunkt des Schreibens, nicht vorab vergeben.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird `spec/` (DOK, Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area)); die Schwelle aus zwei von drei Achsen ist dort erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md) (Vertragstext), [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (Zusage-Menge, „weist aus“ unter `--quiet`), [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md) (ein neuer Befehl gehört in die Befehlsübersicht, die Stichwort-Suche findet sie nicht; deshalb nennt die DoD die Befehlsübersicht ausdrücklich). Evidenz-/Diskrepanz-Risiko: mittel.

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
