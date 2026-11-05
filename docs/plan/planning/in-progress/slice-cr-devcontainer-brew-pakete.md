# Slice CR: Homebrew-Formeln im Devcontainer — Lastenheft-Änderung zu `LH-FA-DEV-003` und `LH-FA-CONF-002`: Schlüssel `brew` mit `packages`, Tap-Freigabeliste

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle, Release-Zuordnung offen. Erster von zwei Slices: ihm folgt [`slice-v1-devcontainer-homebrew-pakete`](../next/slice-v1-devcontainer-homebrew-pakete.md).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Änderung), [`LH-FA-CONF-002`](../../../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) (zwei neue Schlüssel in `u-boot.yaml`), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (Auslegung: ein Tap ist eine vom Nutzer genannte Quelle; nur Anwendung, falls der Wortlaut sie nicht ändert), [`ADR-0046`](../../adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md) (Entscheidungsgrundlage, `Proposed`).
Anlass: Projektinhaber, 2026-11-05 (Entscheidungen 1 und 2 im Slice [`slice-v1-devcontainer-homebrew-pakete`](../next/slice-v1-devcontainer-homebrew-pakete.md), §6).

**Berührte Spec-Stellen:** `lastenheft.md`, Abschnitt zu [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Katalog: Schlüssel `brew`, Feld `packages`), Abschnitt zu [`LH-FA-CONF-002`](../../../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) (Schlüssel `devcontainer.brewTaps.allow`), ggf. [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), §16 (Historie), Versionsfeld. `spezifikation.md` und `architecture.md` berührt dieser Slice nicht.

**Verantwortlich:** —

**Autor:** Claude (Planner-Entwurf auf Wunsch des Projektinhabers). **Datum:** 2026-11-05.

---

## 1. Ziel und Abgrenzung

**Ziel:** Das Lastenheft sagt zu, dass der Katalogschlüssel `brew` eine Liste von Homebrew-Formeln (`devcontainer.features.brew.packages`) im Devcontainer-Dockerfile installiert, dass Formeln aus einem Tap eine Eintragung in `devcontainer.brewTaps.allow` brauchen und dass die Formeln ungepinnt sind; der Wortlaut trägt nur, was die Messung in [`slice-v1-devcontainer-homebrew-pakete`](../next/slice-v1-devcontainer-homebrew-pakete.md) (DoD 3) und die ADR tragen.

Ausgangslage (gelesen im Bericht des Architects, Fundstellen in [`ADR-0046`](../../adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md)): Keine Zusage des Lastenhefts beschreibt einen Katalogschlüssel mit freier Paketliste oder eine Freigabeliste für Taps; [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) nimmt Paketmanager-Pakete im Image-Bau aus, nennt aber nicht freigegebene Skripte.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Spezifikation und Code** — Folge-Slice [`slice-v1-devcontainer-homebrew-pakete`](../next/slice-v1-devcontainer-homebrew-pakete.md); die Spezifikation setzt den freigegebenen Vertrag voraus.
- **Pin-Zusage für Formeln** — Homebrew pinnt Formeln nicht verlässlich; der Wortlaut sagt „ungepinnt“, nicht mehr.
- **Plattform-Zusage über Linux-`amd64` hinaus** — sie folgt der Messung, nicht der Erwartung.

## 2. Definition of Done

- [ ] **CR-Vorschlag und Freigabe:** `docs/plan/cr/cr-vorschlag-brew-pakete.md` (Form `cr-vorschlag.template.md`) trägt je Stelle Alt/Neu; der Projektinhaber gibt den Wortlaut frei und bestätigt [`ADR-0046`](../../adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md) (Schlüsselname, Listenpfad, Pin-Stufe).
- [ ] **Übernahme:** `spec/lastenheft.md` trägt den Wortlaut (Version `0.36.0` oder die nächste freie, Historie-Zeile mit Datum); Abweichungen der Übernahme sind im CR-Vorschlag ausgewiesen.

Pflichten der Closure (zählen nicht mit):

- [ ] `make gates` grün (im Devcontainer-Container, `AGENTS.md` §3.1).
- [ ] Review durchgeführt, Report unter `docs/reviews/` (`.harness/skills/reviewer.md`), kein Self-Review.
- [ ] Verification-Report unter `docs/reviews/`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/cr/cr-vorschlag-brew-pakete.md` | neu | Vorschlags-Datei nach `cr-vorschlag.template.md` (bleibt dauerhaft im CR-Ordner) |
| `spec/lastenheft.md` | update | die genannten Abschnitte, Versionsfeld, §16 Historie |
| `docs/plan/planning/open/README.md`, `docs/plan/planning/in-progress/roadmap.md` | update | Zeilen des Slice, Zustandsmarker |

Der Implementer liest `lastenheft.md` an den Stellen und ergänzt die Tabelle (`AGENTS.md` §3.10).

## 4. Trigger

**Start** (`open` → `next` → `in-progress`): Projektinhaber priorisiert; `in-progress/` ist leer (WIP-Limit 1). Die Messung aus [`slice-v1-devcontainer-homebrew-pakete`](../next/slice-v1-devcontainer-homebrew-pakete.md) (DoD 3) ist die Grundlage für „nur Gemessenes“; ohne sie bleibt der Wortlaut ein Zielbild mit ausgewiesenen Grenzen.

**Rückführungen:** `in-progress` → `open`, solange die Freigabe des Projektinhabers aussteht.

## 5. Closure-Trigger

DoD vollständig, Wortlaut im Lastenheft, `make gates` grün, Reports liegen vor, Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- Der Wortlaut sagt mehr zu, als die Messung trägt (Zusage weiter als Messung). — **Ausgang:** offen; der Reviewer prüft jeden Satz gegen das Protokoll.
- Messung und CR liegen in umgekehrter Reihenfolge der Slices (die Messung steht in DoD 3 des Folge-Slice). — **Ausgang:** offen; der CR trägt Zielbild-Sätze mit ausgewiesenen Grenzen, bis das Protokoll vorliegt.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `spec/` und `docs/plan/` (DOK); beide stehen in `harness/conventions.md` §Modus-Deklaration.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Stichworte Zusage, Messung, Vertragsänderung, Freigabe); vier Treffer:

- [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (37 Evidence-Dateien, offen, kein Ausgang zugewiesen): der Neu-Wortlaut trägt nur Gemessenes oder ein ausgewiesenes Zielbild; der Reviewer prüft jeden Satz gegen die Grundlage (DoD 1, Risiko 1).
- [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md) (17 Evidence-Dateien, offen): ein Wortlaut, der nach dem Review entsteht, hat ihn nicht durchlaufen; die Übernahme ins Lastenheft ändert den freigegebenen Wortlaut nicht, und jede Abweichung steht im CR-Vorschlag (DoD 2).
- [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (41 Evidence-Dateien, verkörpert in `AGENTS.md` §3.10): jede Bedingung im Neu-Wortlaut („nur wenn“, „immer“) bekommt im Umsetzungs-Slice einen pinnenden Test oder entfällt.
- [`BEO-DOK/freigabe-vorlage-verdichtet-den-entwurf`](../observations/BEO-DOK/freigabe-vorlage-verdichtet-den-entwurf/observation.md) (1 Evidence-Datei, offen): die Vorlage an den Projektinhaber zeigt den Wortlaut Satz für Satz und nicht verdichtet.

**Modus-Begründungsblock:** alle berührten Sub-Areas sind Greenfield. Kein BF-/Hybrid-Block nötig.
