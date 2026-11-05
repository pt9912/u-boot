# Slice V2: `devcontainer remove` — Spezifikation: Aufruf, Suche, Abfrage, Ausgabe und Exit-Codes von `remove`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle, Release-Zuordnung offen. Dritter von vier Slices; ihm geht [`slice-cr-devcontainer-remove`](slice-cr-devcontainer-remove.md) voraus. Reihe (vier Slices, je einzeln lieferbar): [`slice-v1-devcontainer-remove-messung`](../in-progress/slice-v1-devcontainer-remove-messung.md) → [`slice-cr-devcontainer-remove`](slice-cr-devcontainer-remove.md) → [`slice-v2-spec-devcontainer-remove`](slice-v2-spec-devcontainer-remove.md) → [`slice-v1-devcontainer-remove`](slice-v1-devcontainer-remove.md).

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Scope), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (Exit-Codes), [`ADR-0015`](../../adr/0015-sandbox-volumes-pro-instanz.md).

**Berührte Spec-Stellen:** [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (Aufruf, Suche, Ausgabe und Exit-Codes von `up` und `down`; `remove` erweitert sie) und der Abschnitt zur `--json`-Ausgabe, in dem `devcontainer` als Wert von `command` steht; ob eine neue `SPEC-<NNN>` entsteht oder [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) wächst, entscheidet der Architect.

**Verantwortlich:** —

**Autor:** Claude (Planner-Entwurf auf Wunsch des Projektinhabers). **Datum:** 2026-11-05.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Spezifikation legt fest, wie `u-boot devcontainer remove` den Container sucht (wie `down`: Label `devcontainer.local_folder`, Instanz-Label), welche Volumes `--volumes` erfasst, wie die Abfrage und `--yes` wirken, was die Ausgabe und die `--json`-Ausgabe tragen und welche Exit-Codes gelten.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Lastenheft** — liegt in [`slice-cr-devcontainer-remove`](slice-cr-devcontainer-remove.md); die Spezifikation präzisiert es, erweitert es nie.
- **Code** — Folge-Slice [`slice-v1-devcontainer-remove`](slice-v1-devcontainer-remove.md).

## 2. Definition of Done

- [ ] **Spezifikation:** Aufruf, Suche, Volume-Auswahl, Abfrage, Ausgabe (Text und `--json`), Exit-Codes und Grenzen stehen in `spec/spezifikation.md`; jede Bedingung, die den Code beschreibt („nur wenn“, „immer“), bekommt im Code-Slice einen pinnenden Test oder entfällt (`AGENTS.md` §3.10).

Pflichten der Closure (zählen nicht mit):

- [ ] `make gates` grün (im Devcontainer-Container, `AGENTS.md` §3.1).
- [ ] Review durchgeführt, Report unter `docs/reviews/` (`.harness/skills/reviewer.md`), kein Self-Review.
- [ ] Verification-Report unter `docs/reviews/`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/spezifikation.md` | update | [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) und die `--json`-Stelle bzw. eine neue `SPEC-<NNN>` |
| `docs/plan/planning/open/README.md`, `docs/plan/planning/in-progress/roadmap.md` | update | Zeilen des Slice, Zustandsmarker |

## 4. Trigger

**Start** (`next` → `in-progress`): [`slice-cr-devcontainer-remove`](slice-cr-devcontainer-remove.md) liegt in `done/`; `in-progress/` ist leer.

**Rückführungen:** `in-progress` → `next`, wenn der CR-Wortlaut eine Spec-Stelle nicht trägt.

## 5. Closure-Trigger

DoD vollständig, `make gates` grün, Reports liegen vor, Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- Die Spezifikation beschreibt Verhalten, das der Code noch nicht hat (Zielbild). — **Ausgang:** offen; der Code-Slice pinnt jede Aussage mit einem Test.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `spec/` und `docs/plan/` (DOK); alle stehen in `harness/conventions.md` §Modus-Deklaration.

**Vorgelagert — offene Beobachtungen sichten:** vor dem Übergang nach `next/` auszuführen (Register `docs/plan/planning/observations/`); der Entwurf trägt noch kein Ergebnis.

**Modus-Begründungsblock:** alle berührten Sub-Areas sind Greenfield. Kein BF-/Hybrid-Block nötig.
