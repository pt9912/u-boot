# Slice V1: Ein zweiter Name für die Konfigurationsdatei — `.u-boot.yml` neben `u-boot.yaml`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** Vorschlag Release `v0.17.0` (Roadmap *Nächste Wellen*): Der Slice ändert einen Vertrag des Produkts (Name der Konfigurationsdatei) und gehört damit in ein Release. Die Zuordnung entscheidet der Projektinhaber; bis dahin ist er ohne Welle geführt. Anlass: Wunsch des Projektinhabers vom 2026-10-31, die Datei wie `.d-check.yml` als `.u-boot.yml` schreiben zu können.

**Bezug:** [`LH-FA-CONF-002`](../../../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) (Konfigurationsdatei), [`LH-FA-INIT-004`](../../../../spec/lastenheft.md#lh-fa-init-004--bestehendes-projekt-erkennen) (Erkennung eines bestehenden Projekts, harte Marker), [`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz) (Überschreibschutz). Voraussetzung: [`slice-v1-konfigdatei-name-als-konstante`](../in-progress/slice-v1-konfigdatei-name-als-konstante.md) (die eine Stelle für den Namen).

**Berührte Spec-Stellen:** nicht gelesen — prüfen (die Spezifikation nennt `u-boot.yaml` an vielen Stellen; Suche beim Start).

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

**Ziel:** Ein Projekt darf seine Konfiguration als `.u-boot.yml` führen; alle Befehle finden und ändern sie unter diesem Namen, ohne dass ein bestehendes Projekt mit `u-boot.yaml` etwas ändern muss.

**Befund** (gelesen am 2026-10-31; [`AGENTS.md`](../../../../AGENTS.md) §3.10):

- **Der Name ist fest.** `u-boot` kennt nur `u-boot.yaml` (32 Pfad-Literale in 15 Dateien des Application-Layers, siehe [`slice-v1-konfigdatei-name-als-konstante`](../in-progress/slice-v1-konfigdatei-name-als-konstante.md)); es gibt keine Suche nach einem anderen Namen und keine Option dafür. d-check dagegen liest `.d-check.yml` aus der Wurzel und kennt `-config`.
- **Drei Stellen tragen mehr als das Lesen:** `collisionError` (`internal/hexagon/application/initproject.go`, Zeile 935) erkennt nur `u-boot.yaml`, `compose.yaml` und `.env.example` als harte Marker eines bestehenden Projekts ([`LH-FA-INIT-004`](../../../../spec/lastenheft.md#lh-fa-init-004--bestehendes-projekt-erkennen)); `credentialScanFiles` (`doctor_sandbox.go`, Zeilen 58 bis 64) nennt `u-boot.yaml` unter den Dateien, die `u-boot doctor` auf Klartext-Zugangsdaten prüft — eine Konfiguration unter anderem Namen bliebe ungeprüft; `init` legt `u-boot.yaml` an (`initproject.go`, Zeile 854).
- **Nutzertexte:** `README.md` und `README.de.md` nennen `u-boot.yaml` je siebenmal, das Benutzerhandbuch 42-mal; Meldungstexte im Nicht-Test-Code tragen den Namen 251-mal als Wort ([`slice-v1-konfigdatei-name-als-konstante`](../in-progress/slice-v1-konfigdatei-name-als-konstante.md), Befund).

**Wege** (Entscheidung des Projektinhabers vor dem Start; der Slice wählt nicht):

- **A, genau ein zusätzlicher Name:** `.u-boot.yml` neben `u-boot.yaml`. Kleinste Vertragsänderung.
- **B, die übliche Menge:** `u-boot.yaml`, `u-boot.yml`, `.u-boot.yaml` und `.u-boot.yml`. Mehr Komfort, mehr Fälle für Rangfolge und Tests.
- **C, die Datei über eine Option benennen** (zum Beispiel `--config <Datei>` wie bei d-check): größere Änderung an jedem Befehl, ändert die Aufrufe in Skripten.

**Entscheidungen, die vor dem Start fallen müssen** (der Plan nennt Empfehlungen, nicht Beschlüsse):

1. *Mehrere Namen vorhanden:* Ein Fehler mit Exit-Code nach [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) ist die Empfehlung (kein stiller Vorrang, der zwei Wahrheiten verdeckt); die Alternative ist ein fester Vorrang.
2. *Was `init` anlegt:* Standard bleibt `u-boot.yaml` (bestehende Projekte, Skripte und Dokumentation); ein Schalter für `.u-boot.yml` ist die Empfehlung.
3. *`config set` und die übrigen Schreiber:* Sie ändern die vorhandene Datei unter ihrem Namen und benennen nie um.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Migration von `u-boot.yaml` nach `.u-boot.yml`.** Der Slice macht den zweiten Namen möglich, nicht verpflichtend; ein Umbenennen bleibt Handarbeit.
- **Ersetzen des Textes `u-boot.yaml` in Meldungen und Doku durch „die Konfigurationsdatei“** über das nötige Maß hinaus. Texte, die den Namen der gefundenen Datei nennen müssen (Fehlermeldungen), folgen dem Fund; der übrige Text bleibt ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) gilt für Zusagen im Handbuch).
- **Weitere Projektdateien** (`compose.yaml`, `.env.example`): Die Frage nach ihren Namen stellt sich nicht.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Finden und Ändern unter dem zweiten Namen** ([`LH-FA-CONF-002`](../../../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration)): Alle Befehle, die die Konfiguration lesen oder schreiben, arbeiten mit `.u-boot.yml`; ohne diese Datei und mit `u-boot.yaml` ändert sich kein Verhalten (Bestandstests und Golden-/Replay-Cases unverändert grün); liegen mehrere Namen vor, gilt die entschiedene Regel mit Exit-Code und Meldung.
- [ ] **`init`, Erkennung und Sicherheit:** `init` erkennt eine vorhandene Konfiguration unter jedem Namen als harten Marker und überschreibt sie nie ([`LH-FA-INIT-004`](../../../../spec/lastenheft.md#lh-fa-init-004--bestehendes-projekt-erkennen), [`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz)); `u-boot doctor` prüft jede vorhandene Konfiguration auf Klartext-Zugangsdaten; jeweils mit Test und rot färbender Mutation.
- [ ] **Lastenheft-CR, Spezifikation, Handbuch, READMEs und `CHANGELOG.md` nach dem Code** mit Beleg-Tabelle (Bedingung → Test oder Streichung); Handbuch und READMEs ohne Kennungen ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6).
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
| `internal/hexagon/domain/config.go` | update | die Liste der Kandidaten (aus [`slice-v1-konfigdatei-name-als-konstante`](../in-progress/slice-v1-konfigdatei-name-als-konstante.md)) bekommt den zweiten Eintrag; die Rangfolge-Regel kommt in die Domäne, das Prüfen auf vorhandene Dateien in die Application |
| `internal/hexagon/application/` (die 15 Dateien mit dem Pfad-Literal) | update | jede Stelle liest über die Auflösung statt über den festen Namen; nach dem Vorgänger-Slice greifen sie auf die Konstante zu |
| `internal/hexagon/application/initproject.go` | update | `collisionError` (Zeile 935): harte Marker über alle Kandidaten; Anlegen der Datei (Zeile 854) und Laden (Zeile 519) |
| `internal/hexagon/application/doctor_sandbox.go` | update | `credentialScanFiles` (Zeilen 58 bis 64) prüft jeden vorhandenen Kandidaten |
| Tests der Application und der CLI, Docker-Tests | update / neu | Bestandstests unverändert; neue Fälle je Bedingung (zweiter Name, beide Namen, keiner) |
| `spec/lastenheft.md` | CR vor dem Code | [`LH-FA-CONF-002`](../../../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration), [`LH-FA-INIT-004`](../../../../spec/lastenheft.md#lh-fa-init-004--bestehendes-projekt-erkennen) und [`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz); Wortlaut-Freigabe des Projektinhabers |
| `spec/spezifikation.md`, `spec/architecture.md` | update (prüfen) | nicht gelesen |
| `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md`, `CHANGELOG.md` | update (prüfen) | 7, 7 und 42 Vorkommen des Namens; nicht auf die Aussage gelesen |
| `harness/replay.md` | prüfen | Golden Cases für `init` und `generate`; nicht gelesen |
| `docs/plan/adr/` | neu | eine ADR für Wahl und Rangfolge (Vorbild: [ADR-0031](../../adr/0031-devcontainer-shell-ohne-editor.md)) |
| diese Datei, `docs/plan/planning/open/README.md` | update | Zeile dieses Slice; Moves per `make slice-mv` |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [`slice-v1-konfigdatei-name-als-konstante`](../in-progress/slice-v1-konfigdatei-name-als-konstante.md) liegt in `done/`; `Verantwortlich:` gesetzt; der Projektinhaber hat den Weg (A, B oder C) und die drei Entscheidungen aus §1 getroffen und den Lastenheft-CR freigegeben; kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Weg C ändert die Aufrufe jedes Befehls — dann ein eigener Slice je Befehlsgruppe.
- `in-progress` → `open` (blockiert): Der Lastenheft-CR wird nicht freigegeben.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; `make gates` grün; Golden-/Replay-Cases für Projekte mit `u-boot.yaml` unverändert grün; Review-Report und Verification-Report liegen vor; die Closure-Notiz mit Steering-Loop-Lerneintrag ist geschrieben. Die Welle schließt mit dem Release (Closure-Bedingung der Welle, nicht dieser Slice).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Eine Konfiguration unter dem neuen Namen bleibt ungeprüft** (Klartext-Zugangsdaten, Überschreibschutz), weil eine der Stellen mit der Namensliste vergessen wird. — **Ausgang:** bei Closure; Vorsorge: die Liste steht an einer Stelle ([`slice-v1-konfigdatei-name-als-konstante`](../in-progress/slice-v1-konfigdatei-name-als-konstante.md)), und ein Test pro Stelle.
- **Zwei Konfigurationen in einem Projekt** (zwei Namen): stiller Vorrang verdeckt Abweichungen. — **Ausgang:** bei Closure (die Regel aus §1).
- **Skripte und Dokumentation nennen weiter `u-boot.yaml`:** Der Standard von `init` ändert sich nicht; Nutzer, die ihn ändern, müssen ihre Skripte nachziehen. — **Ausgang:** bei Closure.
- **Zusage weiter als die Messung** ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)): Das Handbuch nennt nur gemessene Fälle. — **Ausgang:** bei Closure.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `internal/hexagon/application` (APP), `internal/hexagon/domain` (DOM), `internal/adapter/driving/cli` (CLI) und `internal/e2e` (E2E) sowie `spec/`, `harness/` und `docs/plan/` (DOK). Alle sind in der Modus-Deklaration von [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) als Greenfield auditiert.

**Vorgelagert — offene Beobachtungen sichten:** Das Register ([`observations/`](../observations/)) ist nach den Titeln durchgesehen; keine Beobachtung betrifft den Namen der Konfigurationsdatei. Gelesen: keine. Nur nach Titel gesichtet, Inhalt beim Start zu lesen: [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
