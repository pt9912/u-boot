# Slice CR: `egress` im Wertebereich des Feldes `command` der `--json`-Ausgaben (`LH-NFA-USE-004`, `LH-FA-CLI-007`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-sandbox-egress-ollama`](../welle-sandbox-egress-ollama.md) — zusätzlicher CR-Slice der Welle (kein Zählplatz in der Reihe der übrigen Slices; Eintrag in §4 der Welle-Datei), ergibt sich aus Review und Verification von [`slice-v1-cli-egress-anzeige`](../done/slice-v1-cli-egress-anzeige.md) (Review F-3, Verification Auflage 2). Er folgt [`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md) und geht [`slice-v2-spec-sandbox-egress-ollama`](../open/slice-v2-spec-sandbox-egress-ollama.md) voraus.

**Bezug:** [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) (bindet alle `--json`-Ausgaben an das Enum von `command`), [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) (definiert das Enum). Verwandt: [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (der Befehl `u-boot egress list`), [`LH-SA-CLI-001`](../../../../spec/lastenheft.md#lh-sa-cli-001--befehlsstruktur).
Anlass: Review-Report `2026-10-18-slice-v1-cli-egress-anzeige` F-3 und Verification-Report Auflage 2 (beide unter `docs/reviews/`); der Projektinhaber hat angeordnet, bis zum Release ohne Rückfragen weiterzumachen.

**Berührte Spec-Stellen:** `lastenheft.md` §3 ([`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run)), §16 (Historie), Versionsfeld. [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) (Zeile `enum` von `command`) berührt der CR **nicht**: Sie gehört [`slice-v2-spec-sandbox-egress-ollama`](../open/slice-v2-spec-sandbox-egress-ollama.md) („Spezifikation nach dem Code“, nach dem Lastenheft). `architecture.md` berührt der CR nicht.

**Verantwortlich:** pt9912

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Wert `egress` steht im Wertebereich des Feldes `command` im Lastenheft, damit die ausgelieferte Ausgabe `command: "egress"` von `u-boot egress list --json` innerhalb des Vertrags liegt (Änderung auf die **nächste freie Minor-Version**, `0.16.0`; Begründung unten). Der Wortlaut ist **aus dem freigegebenen Befehl `u-boot egress list` abgeleitet, nicht wortgleich freigegeben**; die Bestätigung des Projektinhabers ist ein **Gate vor dem Release** (nicht vor `done/` dieses Slice), nach dem Muster der Stelle 4 im geschlossenen Slice [`slice-cr-sandbox-egress-anzeige`](../done/slice-cr-sandbox-egress-anzeige.md) §7.

**Warum eine Lastenheft-Änderung:** [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) verlangt für alle `--json`-Ausgaben den Minimalkontrakt mit `command` „als in [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) definiertes Enum“; dieses Enum nennt zehn Werte, `egress` nicht. Der Code liefert `command: "egress"` (`runEgressList` in `internal/adapter/driving/cli/egress.go`, `newDataEnvelope("egress", "list", …)`). Kein Gate prüft das Enum (`jsontestutil` verlangt `subcommand` nur für `template` und `config`, gelesen laut Review F-3); die Abweichung vom Lastenheft bliebe ohne diesen CR bestehen.

### Gefundener Stand (gelesen am 2026-10-18; Fundstellen nach Datei und Absatz, ohne Zeilennummern)

- `spec/lastenheft.md`, [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run), Absatz „Für `--dry-run --json` ist die Ausgabe als maschinenlesbares JSON …“: `command` ist einer der Werte `init`, `add`, `remove`, `up`, `down`, `doctor`, `logs`, `generate`, `config` oder `template`. Der Satz steht im Absatz zur Vorschau (`--dry-run --json`); `egress list` lehnt `--dry-run` und `--diff` ab (`ErrEgressPreviewNotApplicable`, Exit 2), das Enum bindet den Befehl nur über [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe).
- Dieselbe Anforderung, Absatz „Bei gruppierten Befehlen wie `command == "template"` oder `command == "config"` muss das Feld `subcommand` gesetzt sein (z. B. `list`, `get`, `set`).“: ein Beispiel mit „wie“, keine Vollständigkeitsliste; `egress list` setzt `subcommand: "list"`.
- `spec/lastenheft.md`, [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe), Spiegelstrich „`command` (Hauptbefehl als in [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) definiertes Enum)“ — der Verweis ist die einzige Quelle des Wertebereichs für normale `--json`-Ausgaben; und Absatz „Für `--json`-Antworten gilt zusätzlich: … bei `template` und `config` ist `subcommand` verpflichtend …“ (eine Aufzählung ohne „wie“; `egress` ist darin nicht genannt, setzt das Feld aber trotzdem).
- `spec/lastenheft.md`, [`LH-SA-CLI-001`](../../../../spec/lastenheft.md#lh-sa-cli-001--befehlsstruktur): „`subcommand` ist für kommandospezifische Unterbefehle reserviert (z. B. `template`, `config`).“ — Beispiel, keine Liste.
- `spec/lastenheft.md` §16 Historie: letzte Zeile `0.15.0` (Katalog-Feature `ollama` aus [`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md)); davor `0.14.0` (Anzeige der Egress-Liste; trägt die Klausel „Bestätigung des Projektinhabers vor dem Release ausstehend“ zur Befehlsübersicht); Versionsfeld im Kopf: `0.15.0`.
- `docs/maintainer/cli-json-contract.md` §6.10 (Vermerk „Offene Vertragslücke“) nennt die Lücke und verweist auf diesen Slice.

### Wortlaut-Vorschlag (alt/neu; aus dem freigegebenen Befehl abgeleitet — **Bestätigung des Projektinhabers vor dem Release ausstehend**)

Das Lastenheft nennt keine Slice- oder Wellen-Kennung; im Lastenheft stehen die Verweise als Anker `#…` der Datei.

**Stelle 1 — [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run), Wertebereich von `command`** (nur der Eintrag `egress`; sonst bleibt der Absatz unverändert).

Alt:

> … `command` ist einer der Werte `init`, `add`, `remove`, `up`, `down`, `doctor`, `logs`, `generate`, `config` oder `template`. …

Neu:

> … `command` ist einer der Werte `init`, `add`, `remove`, `up`, `down`, `doctor`, `logs`, `generate`, `config`, `egress` oder `template`. …

Platzierung: die Änderung steht im Enum selbst, weil [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) auf ihn verweist; eine zweite Wertemenge in [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) wäre eine Kopie, und Kopien driften. [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) bleibt unverändert.

**Stelle 2 — Versionsfeld und Historie (§16).**

Alt: `| Version          | 0.15.0 |`

Neu: `| Version          | 0.16.0 |`.

Neue Historie-Zeile (Datum: das des Lastenheft-Commits):

> | `<Version>` | `<Datum des Lastenheft-Commits>` | [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run): Der Wert `egress` gehört zum Wertebereich des Feldes `command` (nach [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) für alle `--json`-Ausgaben; aus dem freigegebenen Befehl `u-boot egress list` abgeleitet; Bestätigung des Projektinhabers vor dem Release ausstehend). | Vereinbarung mit dem Projektinhaber |

**Gelesen, bewusst nicht geändert (Entscheidung des Projektinhabers, nicht dieses CR):** ob `egress` auch in die Aufzählung „bei `template` und `config` ist `subcommand` verpflichtend“ ([`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe)) gehört. Der Code setzt `subcommand: "list"` immer; das Lastenheft verbietet es nicht (das Feld ist optional). Der Auftrag lautet, nur den Eintrag `egress` zu ergänzen.

### Lastenheft-Version: getrennt (0.16.0), nicht gebündelt mit `slice-cr-devcontainer-ollama`

Bündeln hieße: beide Änderungen im selben Lastenheft-Commit mit derselben Version `0.15.0`. Entschieden ist **getrennt**, der Enum-CR trägt `0.16.0` (oder `0.15.0`, wenn er vor dem Ollama-CR landet; die Version ist die nächste freie zum Zeitpunkt des Commits):

1. **Verschiedener Freigabestand.** Der Ollama-CR trägt einen ausdrücklich freigegebenen Wortlaut (Gespräch, „Zusage nur für Gemessenes“) und die Messwerte; dieser CR trägt einen abgeleiteten, nicht freigegebenen Satz mit Gate vor dem Release. In einer Version und einer Historie-Zeile stünde das Freigegebene neben dem Unbestätigten; die Klausel „Bestätigung ausstehend“ gälte für die ganze Zeile oder für eine halbe.
2. **Verschiedener Review.** Der Review des Ollama-CR läuft über die freigegebene Fassung (Lehre der Vorgänger-Welle); mit gebündeltem Enum-Satz bräuchte dieser Review auch den Nachlauf nach der Bestätigung (Register `BEO-DOK/vertragsaenderung-nach-review`, 2×). Getrennte Commits halten den Nachlauf auf den Diff dieses CR.
3. **Verschiedene Rücknahme.** Lehnt der Projektinhaber den Enum-Wert ab oder ändert ihn, ist das ein CR-Nachlauf nur dieser Zeile; bei einer Bündelung läge er in der Version des Ollama-Features.

Die Reihenfolge der Welle ist verbindlich: erst [`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md) (`0.15.0`), dann dieser Slice (`0.16.0`). Die Versionsangabe in [`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md) (`0.15.0` bleibt dort) ist dort mitgeführt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **[`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) (Zeile `enum` von `command`), [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) und `architecture.md`** — gehören [`slice-v2-spec-sandbox-egress-ollama`](../open/slice-v2-spec-sandbox-egress-ollama.md) (Übergabe dort eingetragen); die Spezifikation folgt dem Lastenheft und dem Code.
- **Code und Tests** — der Code liefert `command: "egress"` bereits ([`slice-v1-cli-egress-anzeige`](../done/slice-v1-cli-egress-anzeige.md)); ein Gate für das Enum (etwa in `jsontestutil`) wäre ein neuer Sensor und ein eigener Slice, hier nicht angelegt: Das Lastenheft ändert nur die Aussage.
- **Die Aufzählung der Pflicht von `subcommand`** ([`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe)) — ein anderer Satz, eine eigene Entscheidung des Projektinhabers (siehe oben).
- **Die Befehlsübersicht und ihre Bestätigung** ([`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle)) — Gate des Anzeige-CR ([`slice-cr-sandbox-egress-anzeige`](../done/slice-cr-sandbox-egress-anzeige.md) §7), nicht dieses.

## 2. Definition of Done

- [x] **Lastenheft 0.16.0** (Commit `42fc809`): Stellen 1 und 2 (Enum-Eintrag `egress`, Versionsfeld, Historie-Zeile mit Datum, Verweis „Vereinbarung mit dem Projektinhaber“) im Wortlaut dieses Vorschlags; die Historie-Zeile trägt die Klausel „Bestätigung des Projektinhabers vor dem Release ausstehend“, solange das Gate offen ist; die Bestätigung ist ein **Gate vor dem Release** (Welle-Datei §3), nicht vor `done/` dieses Slice.
- [x] **Vorab-Suche und Nachzug:** Die Suche nach den Werten des Enums (`template`, `generate`, `command`) und nach `egress` im Lastenheft ist festgehalten (Treffer unter „Gefundener Stand“ gelesen); der Vermerk „Offene Vertragslücke“ in `docs/maintainer/cli-json-contract.md` §6.10 und die Aussage „listet zehn Subcommands“ in §6 derselben Datei sind nach dem Lastenheft-Commit gegen den Ist-Stand gelesen und nachgeführt.
- [x] `make gates` grün; `make doc-trace` ohne neue Waise. Läufe und Exit-Codes in §7.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8); der Review läuft über den **vorgelegten** Wortlaut und über jede spätere Änderung nach der Bestätigung (Report-Commit `84e2239`, Fix-Runde `903c73f` und `2e9660d`, Verification-Report-Commit `227baa1`).
- [x] Doku-Update: Benutzerhandbuch (Tabelle der Pflichtfelder der `--json`-Ausgabe: `command` und `subcommand`) und `cli-json-contract.md`; kein CHANGELOG (kein Produktverhalten; der Code liefert den Wert bereits).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo mit Wellen von der Welle-Closure geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | Stellen 1 und 2 ([`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) Enum; Versionsfeld; §16 Historie); erst nach dem Landen von [`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md) (Version) |
| `docs/maintainer/cli-json-contract.md` | update | §6.10 Vermerk „Offene Vertragslücke“ auf den Ist-Stand setzen (Lücke geschlossen, Bestätigung ausstehend) und §6 „listet zehn Subcommands“ prüfen; je Aussage gelesen |
| `docs/plan/planning/done/slice-cr-cli-json-command-enum.md` (diese Datei) | update | Freigabe bzw. Gate-Stand, Vorab-Suche, Closure (§7) |
| `docs/reviews/` (Review-Report, Verification Evidence) | neu | Review und Verifikation nach `harness/review.md` und `harness/verification.md` |
| `docs/user/benutzerhandbuch.md` | update | Tabelle „Immer vorhanden“: `command` führt `egress`, `subcommand` führt `egress` (`list`); gelesen gegen `runEgressList` in `internal/adapter/driving/cli/egress.go` |
| `docs/plan/planning/in-progress/roadmap.md` | update | Arbeitszeile (bei der Closure zurück auf „Nichts in Arbeit.“); Verweis-Nachzug des Moves |
| `docs/plan/planning/done/slice-cr-devcontainer-ollama.md`, `done/slice-v1-cli-egress-anzeige.md`, `open/README.md`, `welle-sandbox-egress-ollama.md`, `docs/reviews/2026-10-18-verification-slice-cr-devcontainer-ollama.md`, `open/slice-v2-spec-sandbox-egress-ollama.md` | update | reiner Verweis-Nachzug der Moves (`make slice-mv`), kein Inhalt |
| `docs/plan/planning/welle-sandbox-egress-ollama.md` | update | Zustandsmarker „Zweites Release-Gate“ auf den Ist-Stand (im Lastenheft `0.16.0` gelandet, Bestätigung offen; Verification V-3); das Gate steht dort bereits in §3 |
| `docs/plan/planning/open/slice-v2-spec-sandbox-egress-ollama.md` | update | Übergabe: Aufzählung „Implementiert sind …“ in `architecture.md` (Plan-Zeile `spec/architecture.md`, Review F-5) und DoD-Zeile „Sicht und Historie“ (Verification V-5a) |
| `docs/plan/planning/observations/BEO-DOK/` (`plan-aussage-ueber-code`: `state.md` und Evidence; `vorab-suche-uebersieht-befehlsuebersicht`: `state.md` und Evidence) | update, neu | Register bei der Closure (§7) |
| `CHANGELOG.md` | unverändert | der CR ändert kein Produktverhalten; ein Eintrag zu `egress list --json` steht bereits im Code-Slice (zu lesen beim Start) |

Die Vorab-Suche (Pflicht vor jeder Vertragsänderung, Muster [`slice-cr-sandbox-egress-anzeige`](../done/slice-cr-sandbox-egress-anzeige.md)): über `spec/lastenheft.md` nach den Mustern `command`, `Enum`, `template` und `egress` (in `alpine:3` per `docker run`, nicht auf dem Host); zusätzlich nach dem Namen eines benachbarten Befehls (hier `generate`) und nach der Befehlsübersicht, nicht nur nach Stichworten der Anforderung (Lehre aus [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md)). Gelesen am 2026-10-18 (Planner, Docker-`grep` mit den Mustern `command` und `Enum`, danach die Absätze gelesen): der Wertebereich steht an genau einer Stelle ([`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run)); [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) verweist darauf; [`LH-SA-CLI-001`](../../../../spec/lastenheft.md#lh-sa-cli-001--befehlsstruktur) nennt nur Beispiele. Der Implementer hat sie vor dem Lastenheft-Commit wiederholt (nur über das Lastenheft; Handbuch, Spezifikation und Sicht fand erst das Review).

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; [`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md) liegt in `done/` (Lastenheft-Version `0.15.0` ist vergeben); [`slice-v1-cli-egress-anzeige`](../done/slice-v1-cli-egress-anzeige.md) liegt in `done/` (der Code liefert `command: "egress"`). Eine Freigabe des Wortlauts ist **keine** Start-Bedingung: Der Projektinhaber hat angeordnet, bis zum Release ohne Rückfragen weiterzumachen; der Satz geht als abgeleitet mit Gate vor dem Release ins Lastenheft.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): der Wortlaut verlangt mehr als eine Anforderung, etwa weil die Pflicht von `subcommand` mitgeändert werden soll.
- `in-progress` → `open` (blockiert): der Projektinhaber lehnt den Wert ab oder verlangt einen anderen Wortlaut vor dem Commit; das Lastenheft bleibt unverändert.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Verification-Evidence und Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7). Der Lastenheft-Commit führt den Satz als abgeleitet mit der Klausel „Bestätigung ausstehend“; §7 vermerkt das Gate vor dem Release. Das Gate steht in den Closure-Kriterien der Welle (§3 der Welle-Datei).

## 6. Risiken und offene Punkte

- **Der Projektinhaber bestätigt den Enum-Wert nicht oder ändert den Wortlaut.** — **Ausgang:** weiter offen → Release-Gate (Welle-Datei §3, zweites Release-Gate; Verification V-1): Eine Ablehnung ist ein CR-Nachlauf (Review des Diffs der Änderung, nicht der ganzen Datei); `command: "egress"` müsste dann angepasst oder der Befehl umbenannt werden (ein Code-Slice).
- **Die Reihenfolge der Lastenheft-Versionen kollidiert mit [`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md).** — **Ausgang:** eingetreten wie erwartet, ohne Kollision: `0.16.0` nach `0.15.0` des Ollama-CR (Lastenheft-Commit `42fc809`); getrennte Versionen sind entschieden (§1).
- **Die Spezifikation führt [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) nicht nach.** — **Ausgang:** weiter offen, Übergabe getragen: [`slice-v2-spec-sandbox-egress-ollama`](../open/slice-v2-spec-sandbox-egress-ollama.md) trägt [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) und die Aufzählung in `architecture.md` (dort §2 und §3 eingetragen) und schließt erst nach diesem CR.
- **Ein Review läuft vor der Bestätigung und deckt die spätere Fassung nicht ab** (Register [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md)). — **Ausgang:** weiter offen: der Review lief über den vorgelegten Wortlaut (`42fc809`), die Fix-Runde änderte keinen Lastenheft-Text; ein Nachlauf-Review über den Diff jeder Änderung nach der Bestätigung bleibt Auftrag des Release-Gates.

## 7. Closure-Notiz

**Wortlaut und Freigabe.** Der Wortlaut ist aus dem freigegebenen Befehl `u-boot egress list` abgeleitet, nicht wortgleich freigegeben (Anordnung des Projektinhabers vom 2026-10-18, bis zum Release ohne Rückfragen weiterzumachen). Das Lastenheft `0.16.0` (Commit `42fc809`) trägt den Satz mit der Klausel „Bestätigung des Projektinhabers vor dem Release ausstehend“. Der Verifier hat den Lastenheft-Diff (drei Hunks, 3 Einfügungen, 2 Löschungen) Satz für Satz mit §1 verglichen: wortgleich bis auf die Link-Form.

**Vorab-Suche (2026-10-18 durch den Planner, wiederholt vor dem Lastenheft-Commit durch den Implementer).** Umfang: Planner und Implementer suchten in `spec/lastenheft.md` (Docker-`grep` nach `command`, `Enum`; danach die Absätze gelesen): der Wertebereich steht an einer Stelle ([`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run)), [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) verweist darauf, [`LH-SA-CLI-001`](../../../../spec/lastenheft.md#lh-sa-cli-001--befehlsstruktur) nennt Beispiele. **Lücke:** Die Suche lief nur über das Lastenheft. Die Kopie des Enums im Benutzerhandbuch (Review F-1, MEDIUM) und die Aufzählung „Implementiert sind …“ in `architecture.md` (Review F-5) fand erst der Review; der Plan hatte „kein Handbuch“ mit der Begründung „kein Produktverhalten“ geführt. Behoben in `903c73f` (Handbuch) beziehungsweise übergeben (Sicht, siehe Folge-Slices).

**Belege.** Review-Report [`2026-10-19-slice-cr-cli-json-command-enum`](../../../reviews/2026-10-19-slice-cr-cli-json-command-enum.md) (Commit `84e2239`; 0 HIGH, 1 MEDIUM, 3 LOW, 2 INFO; F-1 bis F-4 behoben in `903c73f` und `2e9660d`, F-5 übergeben, F-6 durch diese Notiz getragen) und Verification-Report [`2026-10-19-verification-slice-cr-cli-json-command-enum`](../../../reviews/2026-10-19-verification-slice-cr-cli-json-command-enum.md) (Commit `227baa1`; GO mit Auflagen). Auflagen des Verifiers: V-2 und V-3 nachgezogen (Plan-Tabelle §3, Zustandsmarker „Zweites Release-Gate“ in der Welle-Datei), V-4 durch die Roadmap-Rückstellung, V-5a durch die DoD-Zeile im Spezifikations-Slice; **V-5b** (das Handbuch führt für `config` die Subcommands `get`, `set`, `list`, `config.go` liefert zusätzlich `show`) besteht vor diesem CR und liegt außerhalb seines Gegenstands: nicht behoben, hier als Hinweis an den Handbuch-Pfleger festgehalten, kein Register-Eintrag (ein Beleg einer neuen Klasse ohne Wiederauftreten).

**Stand des Gates (Release-Gate, nicht Gate vor `done/`).** Die Bestätigung des Projektinhabers für den abgeleiteten Wortlaut steht aus; das Lastenheft trägt die Klausel, bis er bestätigt. Das Gate steht als zweites Release-Gate in den Closure-Kriterien der Welle [`welle-sandbox-egress-ollama`](../welle-sandbox-egress-ollama.md) §3 (Zustand dort: im Lastenheft `0.16.0` gelandet, Bestätigung offen); der Release-Cut der Welle schneidet nicht, solange es offen ist. Bei Bestätigung: Änderung genau der Klausel, Vermerk mit Datum hier und in der Welle-Datei. Bei Änderung oder Ablehnung: CR-Nachlauf mit Review des Diffs (Register [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md)); `command: "egress"` müsste dann angepasst werden (Code-Slice).

**Berichtigung der Zustandsmarker (Review F-3, Verification V-3).** In diesem Slice (Gefundener Stand, Stelle 2, Ziel, Risiko 2, Vorab-Suche) und in der Welle-Datei (§3, zweites Release-Gate) auf den Ist-Stand gesetzt. Gesucht im Baum nach „noch nicht gelandet“, „erwartet“, „wiederholt“ und `0.16.0` (ohne `docs/reviews/`): kein weiterer abgelaufener Marker zu diesem CR; die Treffer in `done/slice-cr-devcontainer-ollama.md` und `done/slice-v1-cli-egress-anzeige.md` benennen den Enum-CR ohne Zustandsbehauptung (Fundstellen: §6 „eingetreten wie erwartet“, Übergabe-Absatz, Folge-Slices). Offen bleibt die Bestätigung selbst.

**Was hat funktioniert:** Die Trennung der Version (`0.16.0` neben `0.15.0`) hielt den Review und den Nachlauf-Bedarf auf den Diff dieses CR. Der Review fand die Handbuch-Kopie des Enums (F-1), die keine Vorab-Suche über das Lastenheft finden konnte; der Verifier trennte die Frage nach der Wortgleichheit (erfüllt) von der Bestätigung (Gate) und zeigte die zwei Diff-Gruppen ohne Plan-Zeile.

**Was ging anders als geplant:** Die Vorab-Suche war zu eng (nur Lastenheft), „kein Handbuch“ im Plan stimmte nicht. Die Plan-Tabelle deckte die Verweis-Nachzüge von `make slice-mv` und die Übergabe im Spezifikations-Slice nicht (V-2), und ein Zustandsmarker in der Welle-Datei blieb stehen (V-3), obwohl die Datei im Diff stand.

**Steering-Loop-Eintrag:** gezählt, nicht verkörpert; mit diesem Slice wurde keine Regel und kein Sensor verkörpert, kein Eintrag erreicht erstmals die Schwelle von drei. Lerneintrag als **geschärfte Regel für die Vorab-Suche** (Kandidat, zweites Auftreten): Bei einer Lastenheft-Änderung, die einen Wertebereich oder eine Befehlsliste ergänzt, sucht die Vorab-Suche die Aufzählung (die Werte der Liste, nicht nur das Stichwort der Anforderung) in allen Orten, an denen eine Kopie stehen kann: Lastenheft, Spezifikation, Sicht, Benutzerhandbuch, `docs/maintainer/` und README; „kein Produktverhalten“ begründet nicht „kein Handbuch“, wenn die Kopie in der Nutzer-Doku steht — Register [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md), 2×. Der Rückfall in [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (verkörpert in `AGENTS.md` §3.10) ist der einundzwanzigste Vorgang; eine weitere Schärfung ist nicht entschieden.

**Beobachtungs-Register (`../observations/`):** [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md): `evidence/slice-cr-cli-json-command-enum.md` ergänzt (zweiter Beleg, `state.md` nachgeführt, offen, unter der Schwelle). [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md): `evidence/slice-cr-cli-json-command-enum.md` ergänzt (F-2, F-3, V-2, V-3), `state.md` nachgeführt. [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md): kein neuer Beleg, weil die Fix-Runde keinen Lastenheft-Text nach dem Review geändert hat (Verification, Abschnitt 3); der Nachlauf-Bedarf nach der Bestätigung ist Risiko 4 und würde bei einer Änderung dort als dritter Beleg zählen. Kein Eintrag erreicht mit diesem Vorgang die Schwelle von drei.

**Folge-Slices:** [`slice-v2-spec-sandbox-egress-ollama`](../open/slice-v2-spec-sandbox-egress-ollama.md) trägt [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) (Zeile `enum` von `command`) und die Aufzählung „Implementiert sind …“ der Driving-Adapter in `architecture.md` (Plan-Zeile `spec/architecture.md` und DoD „Sicht und Historie“, gelesen vor dem Eintrag); das Zitat des Schemas in `cli-json-contract.md` §3 bleibt bis dahin wortgleich zu [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json). Ein Gate für das Enum (etwa in `jsontestutil`) ist kein Slice dieser Welle; neue Folge-Slices entstehen nicht.

**Risiken aus §6:** jedes trägt in §6 genau einen Ausgang: weiter offen (Bestätigung: Release-Gate; Spec-Übergabe; Nachlauf-Review nach der Bestätigung), eingetreten wie erwartet (Versionsreihenfolge `0.16.0` nach `0.15.0`).

**Grenzen:** (1) Der Wortlaut ist nicht freigegeben, nur abgeleitet; die Bestätigung ist Release-Gate. (2) Die Wortgleichheit ist von Hand geprüft, nicht mit einem Zeichen-Diff-Werkzeug. (3) `jsontestutil` (`checkSubcommand`) hat nur der Review gelesen. (4) `make test-docker`, `make verify-depguard` und `make doc-immutable` liefen nicht (kein Produktcode, keine Importregel, keine ADR im Diff). (5) V-5b ist nicht behoben.

**Drei Paarungen:** Anker: kein Feld `liegt in` in diesem Slice, nichts zu paaren. Folge-Slice: [`slice-v2-spec-sandbox-egress-ollama`](../open/slice-v2-spec-sandbox-egress-ollama.md) liegt in `open/`. Register: `vorab-suche-uebersieht-befehlsuebersicht`, `plan-aussage-ueber-code` und `vertragsaenderung-nach-review` sind Verzeichnisse im Register und tragen je ein nicht leeres `evidence/`.

**Abgleich der Plan-Tabelle und Sensor-Läufe der Closure:** `git diff --stat d20be6d` nach dem letzten Inhalt vor diesem Absatz: 17 Einträge — Lastenheft, `cli-json-contract.md`, Benutzerhandbuch, dieser Slice, Roadmap, Review-Report, Verification-Report, sechs Dateien mit Verweis-Nachzug von `make slice-mv` (davon `open/slice-v2-spec-sandbox-egress-ollama.md` und `welle-sandbox-egress-ollama.md` zusätzlich mit Inhalt) und vier Dateien des Registers (zwei `state.md`, zwei Evidence). Jeder Eintrag hat eine Zeile in der Plan-Tabelle §3, jede Zeile mit Datei-Eintrag einen Diff; `CHANGELOG.md` steht nicht im Diff. Der Move nach `done/` und die Rückstellung der Roadmap-Arbeitszeile folgen als eigene Commits (Roadmap und Verweis-Nachzug sind in §3 geführt). Läufe (Exit-Code des `make`-Prozesses selbst): `make gates` Exit 0 (Coverage 91,40 % gegen 90 %; `docs-check` 454 Dateien, 0 Befunde), `make doc-trace` Exit 0 (83 Anforderungen, 0 Waisen). Host-Werkzeuge: `git`, `grep`, `sed`, `python3` für Text- und Diff-Arbeit, nicht für Gates.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird allein `spec/` (DOK, Greenfield) laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); `docs/maintainer/cli-json-contract.md` zieht den Vermerk nach. Die Schwelle (zwei von drei Achsen) ist laut Modus-Deklaration erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md) (offen, 2 Belege: Review über die freigegebene Fassung plus Nachlauf), [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md) (offen, 1 Beleg: Vorab-Suche nach dem Namen eines benachbarten Befehls und der Befehlsübersicht), [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (verkörpert in `AGENTS.md` §3.10: jede Aussage dieses Plans über Code ist an Datei und Symbol gelesen). Evidenz-/Diskrepanz-Risiko: niedrig.

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
