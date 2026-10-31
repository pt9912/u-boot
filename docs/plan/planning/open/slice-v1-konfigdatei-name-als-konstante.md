# Slice V1: Der Dateiname der Konfiguration `u-boot.yaml` als eine Konstante im Domain-Layer statt 33 Literale im Application-Layer

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Es gibt keine Closure-Bedingung, die von der DoD dieses Slice verschieden ist. Anlass: Die Frage des Projektinhabers, ob die Konfigurationsdatei auch `.u-boot.yaml` heißen könnte (2026-10-31), zeigte, dass der Name an vielen Stellen von Hand wiederholt wird.

**Bezug:** [`LH-FA-CONF-002`](../../../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) (Inhalt der Konfigurationsdatei; der Name der Datei ist Teil des Vertrags und ändert sich in diesem Slice nicht).

**Berührte Spec-Stellen:** —. Das Verhalten bleibt gleich. Die Konvention „der Name steht an einer Stelle“ geht in die Sicht-Spec (`spec/architecture.md`, siehe §3), weil sie im Code entsteht ([`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area), Modus der Sub-Area APP).

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-31.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der Name der Konfigurationsdatei steht im Nicht-Test-Code an genau einer Stelle, in der neuen Datei `internal/hexagon/domain/config.go` (Wunsch des Projektinhabers), alle Pfad-Literale greifen darauf zu, das Verhalten bleibt unverändert, und ein neues Literal fällt auf. Die Stelle trägt den Namen als Liste von Kandidaten mit genau einem Eintrag (`u-boot.yaml`), damit ein zweiter akzeptierter Name wie `.u-boot.yml` später ein Eintrag der Liste ist und keine Änderung an 33 Stellen.

**Befund** (gelesen und gezählt am 2026-10-31; [`AGENTS.md`](../../../../AGENTS.md) §3.10):

- **33 Pfad-Literale `"u-boot.yaml"` in 15 Dateien**, alle unter `internal/hexagon/application/`: `config.go` (6), `generate.go` (6), `doctor.go` (4), `initproject.go` (4), `addservice_execute.go` (2), `removeservice.go` (2) und je eines in `devcontainer_start.go`, `devcontainer_egress.go`, `doctor_sandbox.go`, `config_many.go`, `upservice.go`, `logsservice.go`, `egress_list.go`, `downservice.go`, `addservice.go` (Zählung mit `grep -rn '"u-boot\.yaml"' --include=*.go internal cmd`, Testdateien ausgenommen). Adapter und `cmd/` tragen das Literal nicht.
- **Eine Konstante gibt es nicht** (Suche nach `ConfigFile`, `configFileName` und ähnlichen Namen ohne Treffer). `internal/hexagon/domain/configpath.go` liegt im Domain-Layer, nennt den Namen nur in Doc-Kommentaren und trägt `ConfigPathKind` und `ConfigPath`, nicht den Dateinamen.
- **Die Schichtregel erlaubt den Ort.** `.golangci.yml`, Regel `domain-isoliert` (Zeilen 173 bis 184): Der Domain-Layer darf `application`, `port` und `adapter` nicht importieren; die Regel `application-no-adapter` verbietet der Application nur Adapter. Application importiert Domain heute schon.
- **Weitere Vorkommen des Textes `u-boot.yaml`:** insgesamt 251 Zeilen im Nicht-Test-Go (Literale, Meldungstexte und Kommentare) und 506 Zeilen in 78 Testdateien. Die Meldungstexte sind Produkttext an den Nutzer und tragen den Namen als Wort.

**Wege** (Entscheidung des Projektinhabers vor dem Start; der Slice wählt nicht):

- **A, Konstante, Ersetzung und Schutz:** Eine exportierte Konstante im Domain-Layer (Arbeitsname `ConfigFileName`), die 33 Stellen greifen darauf zu, und ein Test oder eine Lint-Regel schlägt an, wenn im Nicht-Test-Code wieder ein Pfad-Literal `"u-boot.yaml"` entsteht. Die Lint-Regel würde das Profil in `.golangci.yml` ergänzen ([ADR-0003](../../adr/0003-solid-nahes-lint-profil.md) zum Profil; nicht gelesen — prüfen), ein Test liegt im Application-Paket und liest Quelltext.
- **B, Konstante und Ersetzung ohne Schutz:** Kleiner, aber das nächste Literal schleicht sich unbemerkt wieder ein (dieselbe Klasse wie jeder Wert, der von Hand gleich gehalten wird).
- **C, nichts.** Eine Namensänderung bleibt eine Suche über 33 Stellen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein zweiter akzeptierter Name** (`.u-boot.yml`, Wunsch des Projektinhabers vom 2026-10-31, außerdem denkbar `u-boot.yml` und `.u-boot.yaml`). Das ändert den Vertrag: welche Datei `init` anlegt, wie die Marker beim Init eine vorhandene Konfiguration erkennen (`collisionError` in `initproject.go`, Zeile 935, listet `u-boot.yaml`, `compose.yaml`, `.env.example`), die Rangfolge, wenn mehrere Namen vorliegen, und die Rückwärtskompatibilität alter Projekte; es braucht einen Lastenheft-CR. Der Entwurf [`slice-v1-konfigdatei-alternativer-name`](slice-v1-konfigdatei-alternativer-name.md) trägt es; dieser Slice schafft nur die Voraussetzung (die Liste von Kandidaten mit einem Eintrag).
- **Meldungstexte, Kommentare und Testliterale** (insgesamt 251 und 506 Zeilen mit dem Text). Der Name ist dort Wort im Satz oder Erwartungswert eines Tests; ihn aus einer Konstante zusammenzusetzen macht Texte schwerer lesbar und die Tests weniger unabhängig vom Code, den sie prüfen.
- **Der Name in erzeugten Dateien und Vorlagen** (Templates, README-Texte, Handbuch). Sie sind Ausgabe, keine Pfade, und Gegenstand der Replay-/Golden-Cases ([`harness/replay.md`](../../../../harness/replay.md)); eine Änderung dort wäre ein Verhaltensvorgang.
- **Weitere Dateinamen-Literale** (`compose.yaml`, `.env.example`, `.devcontainer/devcontainer.json`). Sie sind eine eigene Frage, die der Befund nicht gezählt hat; der Slice misst nur `u-boot.yaml`.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Konstante und Ersetzung** ([`LH-FA-CONF-002`](../../../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration)): Im Nicht-Test-Code steht das Pfad-Literal `"u-boot.yaml"` nur noch in der Konstante (Zählung: eine Stelle statt 33); das Verhalten ist unverändert (alle Tests und die Golden-/Replay-Cases unverändert grün).
- [ ] **Schutz** (Weg A): Ein Test oder eine Lint-Regel färbt sich rot, wenn im Nicht-Test-Code wieder ein Pfad-Literal `"u-boot.yaml"` außerhalb der Konstante entsteht; die rot färbende Mutation ist benannt und gemessen.
- [ ] **Konvention in der Sicht-Spec:** `spec/architecture.md` nennt die Konstante als die eine Stelle des Dateinamens, ohne Wellen, Slices oder ADRs ([`AGENTS.md`](../../../../AGENTS.md) §3.4).
- [ ] `make gates` grün (im Modus der Umgebung).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

Gelesen am 2026-10-31; die Zeilennummern verschieben sich, der Implementer liest sie am Ist-Stand nach ([`AGENTS.md`](../../../../AGENTS.md) §3.10). Eine Zeile „prüfen“ nennt eine Datei, die der Planner nicht gelesen hat. Welche Zeilen der Weg berührt, entscheidet der Projektinhaber; die Tabelle nennt Weg A.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/domain/config.go` | neu | die eine Stelle für den Dateinamen: die Liste der Kandidaten (ein Eintrag, `u-boot.yaml`) und der Name, den `init` anlegt (Arbeitsnamen `ConfigFileCandidates` und `ConfigFileName`); im Paket `domain` gibt es `config.go` bisher nicht, `configpath.go` bleibt unverändert und nennt den Namen nur in Doc-Kommentaren (gelesen: Zeilen 1 bis 25) |
| `internal/hexagon/application/config.go`, `generate.go`, `doctor.go`, `initproject.go`, `addservice_execute.go`, `removeservice.go`, `devcontainer_start.go`, `devcontainer_egress.go`, `doctor_sandbox.go`, `config_many.go`, `upservice.go`, `logsservice.go`, `egress_list.go`, `downservice.go`, `addservice.go` | update | 33 Pfad-Literale durch die Konstante ersetzen (Fundstellen im Befund; die Zeilennummern liest der Implementer am Ist-Stand) |
| Schutz-Test im Paket `application` oder `.golangci.yml` | neu / update (prüfen) | Weg A; ein Test, der den Quelltext liest, oder ein Eintrag im Lint-Profil; die Wahl und das Verhältnis zu [ADR-0003](../../adr/0003-solid-nahes-lint-profil.md) sind offen und nicht gelesen |
| `spec/architecture.md` | update | die Konvention als Satz im Abschnitt der Domain-Schicht; die genaue Stelle ist nicht gelesen — prüfen |
| `harness/replay.md` | unverändert, prüfen | Der Slice ändert kein Generator-Verhalten; nicht gelesen |
| diese Datei, `docs/plan/planning/open/README.md` | update | Zeile dieses Slice; Moves per `make slice-mv` |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; der Projektinhaber hat den Weg (A, B oder C) entschieden; kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Schutz verlangt eine neue Lint-Regel, die das Profil ändert — dann ist das ein Vorgang mit Entscheidung nach [ADR-0003](../../adr/0003-solid-nahes-lint-profil.md) und nicht Teil dieses Slice.
- `in-progress` → `open` (blockiert): Eine der 33 Stellen lässt sich nicht ohne Verhaltensänderung ersetzen (etwa ein Literal, das als Teil eines längeren Pfads zusammengesetzt wird).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; `make gates` grün; die Zählung nennt eine Stelle statt 33; Golden-/Replay-Cases unverändert grün; Review-Report und Verification-Report liegen vor; die Closure-Notiz mit Steering-Loop-Lerneintrag ist geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Verhaltensänderung bei einer Ersetzung** (ein Literal ist Teil eines zusammengesetzten Pfads oder einer Meldung). — **Ausgang:** bei Closure; Vorsorge: Bestandstests und Golden-/Replay-Cases laufen unverändert.
- **Der Schutz schlägt falsch an** (Testdateien tragen 506 Zeilen mit dem Namen; Meldungstexte 251 im Nicht-Test-Code). — **Ausgang:** bei Closure; der Schutz erfasst nur das Pfad-Literal `"u-boot.yaml"` als eigenständigen String im Nicht-Test-Code.
- **Ein neuer Test oder eine Lint-Regel ohne rot färbende Mutation** ([`AGENTS.md`](../../../../AGENTS.md) §3.6). — **Ausgang:** bei Closure.
- **Die Zählung ist unvollständig:** `"u-boot.yaml"` kann in einer Form stehen, die das Muster `"u-boot\.yaml"` nicht trifft (etwa in einer Raw-String oder als Teil eines Pfades mit Schrägstrich). — **Ausgang:** bei Closure.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt: Was hat funktioniert · Was ging anders als geplant · Steering-Loop-Eintrag · Beobachtungs-Register · Folge-Slices · Risiken aus §6 · Drei Paarungen.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `internal/hexagon/application` (APP) und `internal/hexagon/domain` (DOM) sowie `spec/` (DOK). Alle sind in der Modus-Deklaration von [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) als Greenfield auditiert; für APP steht dort ausdrücklich, dass eine im Code entstandene Konvention in die Sicht-Spec nachgezogen wird (Punkt 3 der DoD).

**Vorgelagert — offene Beobachtungen sichten:** Das Register ([`observations/`](../observations/)) ist nach den Titeln der Verzeichnisse durchgesehen; unter `BEO-APP` betrifft keine Beobachtung Dateinamen-Literale oder Konstanten, und für die Domain (DOM) gibt es im Register kein Verzeichnis. Der Befund dieses Slice ist nach den Regeln des Registers ein erster Beleg und kann bei der Closure als neue Beobachtung angelegt werden.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
