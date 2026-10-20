# Slice v2: Spezifikation und Sicht zum Devcontainer-Start ohne Editor (nach dem Code)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-devcontainer-start-ohne-editor`](../welle-devcontainer-start-ohne-editor.md). Vierter und letzter von vier Folge-Slices von [`slice-v1-devcontainer-start-ohne-cli`](../in-progress/slice-v1-devcontainer-start-ohne-cli.md); nach [`slice-v1-cli-devcontainer-up-down`](../open/slice-v1-cli-devcontainer-up-down.md) (Spezifikation nach dem Code, Muster [`slice-v2-spec-sandbox-egress-ollama`](../done/slice-v2-spec-sandbox-egress-ollama.md)). Release-Zuordnung offen (Vorschlag des Planners: im Release des Befehls; entscheidet der Projektinhaber).

**Bezug:** die neue Anforderung des Lastenhefts aus [`slice-cr-devcontainer-start-ohne-editor`](../open/slice-cr-devcontainer-start-ohne-editor.md), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer); [`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) (muss `Accepted` sein; Entscheidung 10 nennt die Stellen).

**Berührte Spec-Stellen:** `spezifikation.md`: eine neue Stelle (nächste freie Kennung nach [`SPEC-048`](../../../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf)), [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) (Ausnahme 1; außerdem die Liste der ändernden Subkommandos, Ausnahme 2 mit „in der Warnungsausgabe von `generate devcontainer` und `init`“, der Klammersatz „(nur dieser Befehl erzeugt sie)“ in Ausnahme 1 und die Aufzählung „Gekennzeichnet sind in der Warnungsausgabe von `generate devcontainer` …“, alle gelesen im Review vom 2026-10-20), [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) (Enum `command` und die Prosa „Der Wert `egress` gehört zum Feld `command` …“), [`SPEC-024`](../../../../spec/spezifikation.md#spec-024--beispielaufrufe-zu-den-anforderungen) (Beispielliste), [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) (Aufzählung der lesenden Befehle; Label), [`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten) (Aufzählung der ausweisenden Befehle); `architecture.md` §2.3, §2.4, §2.6, §3 (`ARC-008`), §7.1. [`SPEC-042`](../../../../spec/spezifikation.md#spec-042--aufruf-der-compose-engine-und-preflight) bleibt unverändert (sie gilt „vor jedem Compose-Aufruf“; der Launcher ruft Compose nicht auf). Der Verweis zeigt aufwärts.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Architect). **Datum:** 2026-10-20.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Spezifikation beschreibt den Befehl so, wie Code und Tests ihn tragen: Aufrufform (Hilfs-Image, `docker run`-Argumente, Spiegelpfad), Vorab-Prüfung (zwei plus eine, nicht die drei Schritte von [`SPEC-042`](../../../../spec/spezifikation.md#spec-042--aufruf-der-compose-engine-und-preflight)), Durchreichen der Umgebung, Ausgabe unter `--quiet` und `--json`, Exit-Codes, Fallback-Ausweisung und Grenzen. Die Sicht (`architecture.md`) führt Port, Adapter, die externe Abhängigkeit und die Fehlerklasse nach. **Jede Bedingung, die eine Eigenschaft des Codes behauptet, nennt ihren Test oder entfällt** (`AGENTS.md` §3.10).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Code, Tests, Handbuch, `CHANGELOG.md`** — liegen in [`slice-v1-devcontainer-start-fundament`](../open/slice-v1-devcontainer-start-fundament.md) und [`slice-v1-cli-devcontainer-up-down`](../open/slice-v1-cli-devcontainer-up-down.md).
- **`spec/lastenheft.md`** — der CR liegt davor ([`slice-cr-devcontainer-start-ohne-editor`](../open/slice-cr-devcontainer-start-ohne-editor.md)); die Spezifikation präzisiert das Lastenheft, erweitert es nie.
- **Pins, Digests, Versionen im Wortlaut der Spezifikation** — sie stehen in den eingebetteten Dateien und im Golden Case des Fundament-Slice (ein Stand, kein Zustand); die Spezifikation nennt die Pin-Form, nicht die Zahlen.
- **Messstände** — sie stehen im Handbuch und im Messprotokoll; die Spezifikation führt Zusagen und Grenzen als Zielbild (Muster [`slice-v2-spec-zielbild`](../done/slice-v2-spec-zielbild.md)).
- **Nachführung von [`SPEC-042`](../../../../spec/spezifikation.md#spec-042--aufruf-der-compose-engine-und-preflight)** — keine Änderung nötig; ändert der Code die Aufrufe von `up`/`down` (Compose), ist das ein anderer Vorgang.

## 2. Definition of Done

Drei Liefer-Punkte; Gate-Läufe und Closure-Pflichten darunter zählen nicht mit.

- [ ] **Spezifikation:** eine neue Stelle zum Befehl (Aufrufform, Vorab-Prüfung, Pfad-Abbildung, Umgebung, Ausgabe und `--json`-Daten, Exit-Codes, Fallback-Ausweisung, Grenzen, Pin-Form ohne Zahlen) und die Nachführung in [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) (Ausnahme 1; außerdem die Liste der ändernden Subkommandos, Ausnahme 2 mit „in der Warnungsausgabe von `generate devcontainer` und `init`“, der Klammersatz „(nur dieser Befehl erzeugt sie)“ in Ausnahme 1 und die Aufzählung „Gekennzeichnet sind in der Warnungsausgabe von `generate devcontainer` …“, alle gelesen im Review vom 2026-10-20), [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) (Enum `command` und die Prosa „Der Wert `egress` gehört zum Feld `command` …“), [`SPEC-024`](../../../../spec/spezifikation.md#spec-024--beispielaufrufe-zu-den-anforderungen) (Beispielliste), [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) und [`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten) (Aufzählung der Befehle, die Fallbacks ausweisen) sowie [`SPEC-048`](../../../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf) (Abbruch: der Verweis auf `devcontainer up` und Exit `11`), jeweils als Verfeinerung des Lastenhefts.
- [ ] **Sicht (`architecture.md`):** Port `DevcontainerLauncher` (§2.4), Adapter und eingebetteter Kontext im Paket `docker` (§2.6), Driving-Port (§2.3), `ARC-008` (Hilfs-Container als Rolle der Container-Runtime), Fehlerklasse in §7.1 (11 und 12 für den neuen Pfad, zwei Klassifikatoren); sprach- und meilensteinfrei, ohne ADR-, Slice- und Wellenbezüge (`AGENTS.md` §3.4).
- [ ] **Beleg-Tabelle:** je Bedingung der neuen Stelle der Test aus den Slices davor (Datei und Name) oder die Streichung; die Tabelle steht in §7 der Closure-Notiz.
- [ ] `make gates` grün, `make docs-check` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: Roadmap-Zeile „In Arbeit“ bei der Closure; Handbuch und `CHANGELOG.md` unverändert.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; sie wird nach der letzten Fix-Runde gegen `git diff --stat` abgeglichen (`AGENTS.md` §3.10).

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/spezifikation.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), [`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten) und [`SPEC-048`](../../../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf) nennen `devcontainer up` als Befehl, der die Statusdatei liest und `failed` als Exit `11` abbildet ([`ADR-0029`](../../adr/0029-devcontainer-start-exit-11-und-socket-haltung.md); Grenzen: keine Datei, altes `failed`). Vor dem Schreiben sucht der Autor nach jeder Aufzählung der Befehle und nach `egress list` in der ganzen Datei und trägt jede Fundstelle hier ein; die Liste in dieser Zeile ist kein Vollständigkeitsnachweis. Exit-Klasse eines Startscript-Abbruchs: Die Stelle zum Befehl sagt ausdrücklich, ob er auf `11` (Umgebungsproblem, wie [`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten) und `doctor`) oder auf `12` abgebildet wird; die Entscheidung steht in [`ADR-0029`](../../adr/0029-devcontainer-start-exit-11-und-socket-haltung.md) Entscheidung 1 (Exit `11`, mit den Grenzen: keine Statusdatei ohne `u-boot.yaml`, altes `failed`). Neue Stelle nach [`SPEC-048`](../../../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf) (gelesen am 2026-10-20: `SPEC-048` ist die letzte Stelle); [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) Ausnahme (1) zählt die Ausweisungen auf (Sandbox-Profil, Login-Dienst, Warn-Einträge von `doctor`) und kennt die Socket-Lockerung nicht; [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) Abschnitt *Ermittlung durch den Host* sagt „`doctor`, `generate devcontainer` und `init` lesen die Dateien“ und „`u-boot` setzt kein eigenes Label“; [`LH-FA-DEV-007.a`](../../../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten) nennt `doctor`, `generate devcontainer` und `init` als ausweisende Befehle; [`SPEC-042`](../../../../spec/spezifikation.md#spec-042--aufruf-der-compose-engine-und-preflight) Preflight mit drei Schritten „vor jedem Compose-Aufruf“ (bleibt); die Zeile `enum` von `command` in [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) (laut Eltern-Slice) |
| `spec/architecture.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | gelesen am 2026-10-20: §2.4 listet die Driven-Ports (`DockerProbe`, `DevcontainerInspector`, `DockerEngine`, …), §2.6 beschreibt „drei Adapter im selben Paket“ `docker/`, §3 `ARC-008` nennt „Diagnose … und Lifecycle-Operationen der erzeugten Umgebung“, §7.1 die Zeilen „Container-Runtime nicht erreichbar“ (11) und „Laufzeitfehler der Runtime“ (12) |
| `docs/plan/planning/in-progress/roadmap.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Zeile „In Arbeit“ bei der Closure |
| `docs/plan/planning/open/README.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Bestand |
| `docs/reviews/` <!-- d-check:ignore (Pfad-Zitat) --> | neu (Reviewer, Verifier) | Review und Verifikation |
| `spec/lastenheft.md`, `docs/user/benutzerhandbuch.md`, `CHANGELOG.md`, `internal/`, `cmd/` <!-- d-check:ignore (Pfad-Zitate) --> | **unverändert** | Lastenheft im CR-Slice; Handbuch und `CHANGELOG.md` im Befehl-Slice; kein Code |

## 4. Trigger

**Start** (`open` → `in-progress`, ggf. über `next`): `Verantwortlich:` gesetzt; [`slice-v1-cli-devcontainer-up-down`](../open/slice-v1-cli-devcontainer-up-down.md) ist geschlossen (Spezifikation nach dem Code); kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die neue Stelle braucht mehr Bedingungen, als die Tests der Vorgänger-Slices tragen, und es entstünde ein Test-Slice.
- `in-progress` → `open` (blockiert): die Vorgänger-Slices sind nicht geschlossen oder haben die Zusage zurückgenommen; die Evidence nennt den Grund.

## 5. Closure-Trigger

DoD vollständig, Review-Report und Verification-Evidence liegen vor, die Beleg-Tabelle in §7 hat für jede Bedingung einen Test oder eine Streichung, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Die Spezifikation sagt mehr, als die Tests tragen** (Bedingung zu Engines, Pfaden, Plattformen). — **Ausgang:** weiter offen; die Beleg-Tabelle ist das Gegenmittel ([`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md), [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)).
- **Die Sicht trägt einen Meilenstein-, Slice- oder ADR-Bezug.** — **Ausgang:** weiter offen; `AGENTS.md` §3.4 und der Reviewer; die Sicht beschreibt Zustand, nicht Verlauf.
- **Eine weitere Stelle zählt die Befehle mit Fallback-Ausweisung auf**, die der Architect nicht gelesen hat. — **Ausgang:** weiter offen; Suche über die Befehlsnamen und Lesen jedes Treffers (Muster [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md)).

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird `spec/` (DOK, Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area)); die Schwelle aus zwei von drei Achsen ist dort erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md), [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md), [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md), [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md). Evidenz-/Diskrepanz-Risiko: mittel.

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
