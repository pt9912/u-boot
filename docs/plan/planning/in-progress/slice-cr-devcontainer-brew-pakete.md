# Slice CR: Homebrew-Formeln im Devcontainer — Lastenheft: eine neue Anforderung (Formel-Liste im Image, Tap nur mit ausdrücklicher Freigabe, Formeln ungepinnt und ausgewiesen)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle, Release-Zuordnung offen. Erster von zwei Slices: ihm folgt [`slice-v1-devcontainer-homebrew-pakete`](../next/slice-v1-devcontainer-homebrew-pakete.md).

**Bezug:** [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (Anwendung: ein Tap ist eine externe Quelle; der Wortlaut dieser Anforderung ändert sie nicht), [`ADR-0046`](../../adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md) (Entscheidungsgrundlage, `Accepted`). Eine neue Anforderung im Abschnitt der Devcontainer-Anforderungen, nach [`LH-FA-DEV-012`](../../../../spec/lastenheft.md#lh-fa-dev-012--umgebungsdatei-des-devcontainers).
Anlass: Projektinhaber, 2026-11-05 (Entscheidungen 1 und 2 im Slice [`slice-v1-devcontainer-homebrew-pakete`](../next/slice-v1-devcontainer-homebrew-pakete.md), §6).

**Berührte Spec-Stellen:** `lastenheft.md`: eine neue Anforderung nach [`LH-FA-DEV-012`](../../../../spec/lastenheft.md#lh-fa-dev-012--umgebungsdatei-des-devcontainers) (Kennung vergibt die Übernahme), §16 (Historie), Versionsfeld. `spezifikation.md` und `architecture.md` berührt dieser Slice nicht; Schlüsselnamen, Konfigurationspfade, das Verfahren der Tap-Freigabe, die Konflikt-Regel mit dem Katalog-Feature `homebrew` sowie Bau-, Plattform- und Messgrenzen stehen im Umsetzungs-Slice in der Spezifikation.

**Verantwortlich:** —

**Autor:** Claude (Planner-Entwurf auf Wunsch des Projektinhabers). **Datum:** 2026-11-05.

---

## 1. Ziel und Abgrenzung

**Ziel:** Das Lastenheft enthält eine neue Anforderung: Das Produkt installiert Homebrew-Formeln, die der Benutzer in der Projektkonfiguration als Liste führt, im erzeugten Devcontainer-Image; eine Formel aus einem Tap braucht eine ausdrückliche Freigabe, ohne sie endet die Erzeugung mit einem fachlichen Fehler (Exit-Code `10`); die Formeln sind ungepinnt und werden beim Erzeugen als solche ausgewiesen. Der Wortlaut nennt Anforderungen und keine Schlüsselnamen, Plattform-, Bau- oder Messgrenzen.

Ausgangslage (gelesen im Bericht des Architects, Fundstellen in [`ADR-0046`](../../adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md)): Keine Anforderung des Lastenhefts beschreibt eine freie Formel-Liste im Image oder eine Freigabe für Taps; [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) nimmt Paketmanager-Pakete im Image-Bau aus, nennt aber nicht freigegebene Skripte.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Spezifikation und Code** — Folge-Slice [`slice-v1-devcontainer-homebrew-pakete`](../next/slice-v1-devcontainer-homebrew-pakete.md); die Spezifikation setzt den freigegebenen Vertrag voraus.
- **Pin-Zusage für Formeln** — Homebrew pinnt Formeln nicht verlässlich; der Wortlaut sagt „ungepinnt“, nicht mehr.
- **Plattform-, Bau- und Messgrenzen der Zusage** — keine Lastenheft-Aussage; sie stehen in der Spezifikation und in den Tests des Umsetzungs-Slices, mit ihrer Messgrundlage.

## 2. Definition of Done

- [x] **CR-Vorschlag und Freigabe:** `docs/plan/cr/cr-vorschlag-brew-pakete.md` (Form `cr-vorschlag.template.md`) trägt je Stelle Alt/Neu; der Projektinhaber gibt den Wortlaut frei und bestätigt [`ADR-0046`](../../adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md) (Schlüsselname, Listenpfad, Pin-Stufe).
- [x] **Übernahme:** `spec/lastenheft.md` trägt die neue Anforderung (Lastenheft-Version `0.36.0` oder die nächste freie, Historie-Zeile mit Datum); Abweichungen der Übernahme sind im CR-Vorschlag ausgewiesen.

Pflichten der Closure (zählen nicht mit):

- [ ] `make gates` grün (im Devcontainer-Container, `AGENTS.md` §3.1).
- [ ] Review durchgeführt, Report unter `docs/reviews/` (`.harness/skills/reviewer.md`), kein Self-Review.
- [ ] Verification-Report unter `docs/reviews/`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/cr/cr-vorschlag-brew-pakete.md` | update | Vorschlags-Datei nach `cr-vorschlag.template.md`: Freigabe, Übernahme, Abweichungen (bleibt dauerhaft im CR-Ordner) |
| `spec/lastenheft.md` | update | die neue Anforderung, Versionsfeld, §16 Historie |
| `docs/plan/planning/next/slice-v1-devcontainer-homebrew-pakete.md` | update | Bezug, Doku-Pflicht und Test-Anker auf die neue Anforderung (Review F-3) |
| `docs/plan/planning/in-progress/slice-cr-devcontainer-brew-pakete.md` | neu / update | dieser Plan selbst |
| `docs/plan/planning/open/README.md`, `docs/plan/planning/in-progress/roadmap.md` | update | Zeilen des Slice, Zustandsmarker |
| `docs/reviews/2026-11-05-slice-cr-devcontainer-brew-pakete.md` | neu | Review-Report |

Die Tabelle wird nach der letzten Fix-Runde gegen `git diff --stat` abgeglichen (`AGENTS.md` §3.10).

## 4. Trigger

**Start** (`open` → `next` → `in-progress`): Projektinhaber priorisiert; `in-progress/` ist leer (WIP-Limit 1).

**Rückführungen:** `in-progress` → `open`, wenn Review oder Verifikation eine Änderung am Wortlaut verlangen, die der Projektinhaber erneut freigeben muss.

## 5. Closure-Trigger

DoD vollständig, Wortlaut im Lastenheft, `make gates` grün, Reports liegen vor, Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- Der Wortlaut enthält Details, die in die Spezifikation gehören (Plattformen, Messgrenzen). — **Ausgang:** offen; der Reviewer prüft jeden Satz auf Anforderungsebene.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `spec/` und `docs/plan/` (DOK); beide stehen in `harness/conventions.md` §Modus-Deklaration.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Stichworte Zusage, Messung, Vertragsänderung, Freigabe); vier Treffer:

- [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (37 Evidence-Dateien, offen, kein Ausgang zugewiesen): der Neu-Wortlaut nennt Anforderungen; der Reviewer prüft jeden Satz gegen die ADR und die Entscheidungen des Projektinhabers (DoD 1, Risiko 1).
- [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md) (17 Evidence-Dateien, offen): ein Wortlaut, der nach dem Review entsteht, hat ihn nicht durchlaufen; die Übernahme ins Lastenheft ändert den freigegebenen Wortlaut nicht, und jede Abweichung steht im CR-Vorschlag (DoD 2).
- [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (41 Evidence-Dateien, verkörpert in `AGENTS.md` §3.10): jede Bedingung im Neu-Wortlaut („nur wenn“, „immer“) bekommt im Umsetzungs-Slice einen pinnenden Test oder entfällt.
- [`BEO-DOK/freigabe-vorlage-verdichtet-den-entwurf`](../observations/BEO-DOK/freigabe-vorlage-verdichtet-den-entwurf/observation.md) (1 Evidence-Datei, offen): die Vorlage an den Projektinhaber zeigt den Wortlaut Satz für Satz und nicht verdichtet.

**Modus-Begründungsblock:** alle berührten Sub-Areas sind Greenfield. Kein BF-/Hybrid-Block nötig.
