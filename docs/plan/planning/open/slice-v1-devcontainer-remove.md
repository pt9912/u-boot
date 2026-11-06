# Slice V1: `devcontainer remove` — der Befehl: Container entfernen, mit `--volumes` auch seine Volumes (Abfrage, `--yes`, `--instance`, `--json`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-devcontainer-remove`](../welle-devcontainer-remove.md), Release-Zuordnung offen. Letzter von vier Slices; ihm gehen [`slice-v2-spec-devcontainer-remove`](../in-progress/slice-v2-spec-devcontainer-remove.md) und davor der CR voraus. Reihe (vier Slices, je einzeln lieferbar): [`slice-v1-devcontainer-remove-messung`](../done/slice-v1-devcontainer-remove-messung.md) → [`slice-cr-devcontainer-remove`](../done/slice-cr-devcontainer-remove.md) → [`slice-v2-spec-devcontainer-remove`](../in-progress/slice-v2-spec-devcontainer-remove.md) → [`slice-v1-devcontainer-remove`](slice-v1-devcontainer-remove.md).

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Scope), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (Exit-Codes), [`ADR-0015`](../../adr/0015-sandbox-volumes-pro-instanz.md).

**Berührte Spec-Stellen:** [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (und die Stelle, die der Spec-Slice anlegt); `architecture.md` — keine.

**Verantwortlich:** —

**Autor:** Claude (Planner-Entwurf auf Wunsch des Projektinhabers). **Datum:** 2026-11-05.

---

## 1. Ziel und Abgrenzung

**Ziel:** `u-boot devcontainer remove [--instance <name>] [--volumes] [--yes]` entfernt den Container des Projekts und mit `--volumes` seine Volumes, nach dem Wortlaut des Lastenhefts und der Spezifikation.

Ausgangslage (gelesen): Der Befehl `devcontainer down` liegt im CLI-Adapter (`internal/adapter/driving/cli/devcontainer.go`: `newDevcontainerDownCommand`, `runDevcontainerDown`, `mapDevcontainerDownErrorToDiagnostic`) und im Dienst `DevcontainerStartService.Down` (`internal/hexagon/application/devcontainer_start.go`), der den Container mit `inspector.FindContainer` sucht und mit `launcher.Stop` stoppt. Die Bestätigung bei zerstörerischen Aktionen liefert der Port in `internal/hexagon/port/driven/confirm.go`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Änderung von `devcontainer down`** — behält seine Zusage (stoppen, nichts entfernen).
- **Vertrag und Spezifikation** — liegen in den zwei vorangehenden Slices; dieser Slice setzt sie um und ändert sie nicht.
- **Entfernen von Images und des Hilfs-Images** — nicht Gegenstand des Vertrags; es wäre ein eigener Vorgang mit eigener Abfrage.

## 2. Definition of Done

- [ ] **Befehl und Ports:** `remove` (Container, `--volumes`, `--instance`, `--yes`, `--json`) arbeitet nach dem Vertrag; die Ports in `internal/hexagon/port/` bekommen die nötigen Operationen (Entfernen von Container und Volumes, Auflistung der Volumes eines Containers), der Docker-Adapter setzt sie um; Exit-Codes und Diagnostic-Abbildung stehen in der zentralen Klassifikation **und** in der Abbildung des neuen Subkommandos (Dual-Classifier-Regel, `.harness/skills/reviewer.md`). Tests nach [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Happy, Boundary, Negative, ohne Container, laufender Container, Instanz, Abfrage abgelehnt; die Abfrage steht vor dem Stoppen, Exit-Code `12`, wenn Stoppen oder Entfernen scheitert, Exit-Code `0` für `--instance` ohne Container der Instanz, Exit-Code `0` und keine Abfrage für `remove --volumes` ohne `--yes` im nicht-interaktiven Modus ohne Container zu entfernen, Lastenheft `0.37.2`).
- [ ] **Docker-Test:** ein Test mit Build-Tag `docker` (`internal/e2e`) führt `u-boot devcontainer remove` des **gebauten Binary** gegen echte Container eines erzeugten Projekts aus, auf einem Linux-Host (`amd64`) mit lokaler Docker-Engine (die Zusage-Menge des Lastenhefts), und belegt je Fall mit Lesestand vor und nach dem Lauf ([`BEO-DOK/plan-annahme-ueber-fremdzustand-ohne-messstelle`](../observations/BEO-DOK/plan-annahme-ueber-fremdzustand-ohne-messstelle/observation.md)): (a) ein laufender Container ist nach `remove` entfernt, seine Volumes bleiben; (b) `remove --volumes --yes` entfernt genau die Volumes dieses Containers und lässt einen zweiten Container (Instanz) und seine Volumes unberührt; (c) `remove --volumes` ohne `--yes` ohne Terminal endet bei vorhandenem Container mit Exit-Code `10`, Container und Volumes bleiben, ein laufender Container läuft weiter; (d) ohne Container endet `remove` mit Exit-Code `0`, auch `remove --volumes` ohne `--yes` ohne Terminal (Lastenheft `0.37.2`: keine Abfrage, Exit-Code `0`, nicht `10`); (e) abgelehnte Abfrage von `remove --volumes` (Lastenheft `0.37.1`, F-3a): Exit-Code `10`, ein laufender Container läuft weiter, kein Volume ist entfernt; (f) `remove --instance a` ohne Container der Instanz `a`, bei vorhandenem anderen Container des Projekts (F-3b): Exit-Code `0`, der andere Container und seine Volumes sind unberührt. Der Beleg ist der Einzellauf des Tests oder der CI-Lauf `integration-docker` auf dem Push-Tip, nicht das Gesamtpaket ([`BEO-E2E/e2e-paket-am-zeitlimit`](../observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md)); ein übersprungener Test belegt nichts. Der benannte Smoke der Welle (Kriterium 4 in [`welle-devcontainer-remove`](../welle-devcontainer-remove.md) §3) liest diesen Test; er ist kein zweiter.
- [ ] **Grenzen des Satzes „die Volumes dieses Containers“ (Review des CR F-8):** je ein Test (Docker-Test oder Test des Dienstes mit Fake-Port, wie es die Spezifikation festlegt) oder eine in der Spezifikation und im Handbuch ausgewiesene Nicht-Zusage für: ein Volume, das ein anderer Container hält (die Messung belegt nur, dass das Entfernen scheitert, solange ein Container das Volume hält; das Verhalten von `remove --volumes` und der Exit-Code, Lastenheft `0.37.1`: `12` bei scheiterndem Entfernen); anonyme Volumes (`docker rm` lässt sie stehen, `docker rm -v` entfernt sie; die Wahl legt die Spezifikation fest); nicht gemessene Profile (Standardprofil, Katalog-Features außer `ollama`: Test oder Nicht-Zusage).
- [ ] **Bedingungen des Lastenhefts `0.37.2` und `0.37.3` (Nachlauf-Reviews des CR; (1) und (2) gehören zu `0.37.2`, (3) zu `0.37.3`):** (1) der Image-Satz der Zusage-Menge („Im ausgelieferten Container-Image … läuft er nicht, und `u-boot devcontainer remove` läuft dort ebenfalls nicht“) bekommt einen pinnenden Test oder steht in der Spezifikation und im Handbuch als Nicht-Zusage, sonst entfällt er im Lastenheft (`AGENTS.md` §3.10; für den Satz zu `up` ist heute kein pinnender Test gefunden, [`cr-vorschlag-devcontainer-remove`](../../cr/cr-vorschlag-devcontainer-remove.md) Prüfung der Folgeänderung `0.37.2`); (2) Fall (e) der abgelehnten Abfrage braucht ein Terminal: der Implementer liest vor dem Zählen, ob `internal/e2e` ein Testmittel dafür trägt (der Nachlauf-Review fand keinen Treffer für `creack/pty`, `openpty`, `ptmx`, `script -qc` in `internal/e2e`; `internal/adapter/driving/cli/devcontainer_shell_test.go` Zeilen 309 bis 317 nennt für den Pseudo-Terminal-Fall einen Docker-Test unter `script`), sonst benennt der Handoff die Lücke gegenüber dem Benannten Smoke der Welle (Fall (e)), statt den Fall still auf eine tiefere Ebene zu verlagern. (3) Fall (d) pinnt zugleich die Ausnahme in [`LH-FA-CLI-005A`](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) (Lastenheft `0.37.3`, Zeile 214): `remove --volumes` ohne `--yes` im nicht-interaktiven Modus endet ohne Container mit Exit-Code `0` und bei einem Container zu entfernen mit `10`; `u-boot down --volumes` und `u-boot remove --purge` bleiben bei `10`.
- [ ] **Bedingungen des Wortlauts ohne eigenen Fall (Verification des CR V-2; Lastenheft `0.37.3`, Spiegelstrich zu `remove`, Fehlerklassen und Nicht-Ziele in [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)):** je Fall ein Test oder eine in der Spezifikation und im Handbuch ausgewiesene Nicht-Zusage, sonst entfällt die Bedingung im Lastenheft (`AGENTS.md` §3.10). Ebene: Docker-Test mit Lesestand vor und nach dem Lauf, wo der Fall eine Engine braucht; Test des Dienstes oder der CLI mit Fake-Port, wo die Spezifikation es festlegt; ein nur auf der tieferen Ebene belegter Fall steht im Handoff als solcher. Die Sätze nennen die Zählung des Spiegelstrichs zu `remove` (Satz 1 bis 11). (g) **Auswahl ohne `--instance`** (Satz 8; die Messung M5 belegt sie für `down` und `up`, nicht für `remove`): (g1) mit einem Container ohne Instanz und dem einer Instanz im Projekt entfernt `remove` ohne `--instance` nur den Container ohne Instanz, der Container der Instanz und seine Volumes sind unberührt; (g2) gibt es nur Container von Instanzen (laufend, gestoppt, gemischt), endet er mit Exit-Code `10` und einem Hinweis auf `--instance`, kein Container ist entfernt und, auch mit `--volumes --yes`, kein Volume. (h) **Volume der Ollama-Modelle** (Satz 3; Messung M4 belegt für `docker rm`, dass das Volume den Container überlebt): in einem Projekt mit dem Feature `ollama` entfernt `remove --volumes --yes` auch dieses Volume, `remove` ohne `--volumes` lässt es stehen. (i) **Nennung der Volumes** (Satz 4 und 6): `remove --volumes` nennt vor der Abfrage die Volumes, die es entfernen würde (die Reihenfolge Nennung, Abfrage, Eingriff prüft der Test des Dienstes am Fake-Port), und nach dem Entfernen die entfernten Volumes; die genannten stimmen mit den Volumes überein, die der Lesestand nach dem Lauf als fehlend zeigt. (j) **Ungültiger Instanzname** (Fehlerklassen): `remove --instance <name>` mit der leeren Zeichenkette, einem Namen mit `=`, einem Großbuchstaben, einem Namen über 32 Zeichen und einem Namen, der mit einem Bindestrich beginnt oder endet, endet mit Exit-Code `10`, und nichts ist entfernt (Grenze: 32 Zeichen gelten, 33 nicht). (k) **Exit-Code `11`** (Fehlerklassen): ist Docker nicht erreichbar oder die Engine nicht lokal, endet `remove` mit Exit-Code `11`, und nichts ist entfernt; ein Docker-Test nur, wenn `internal/e2e` ein Testmittel dafür trägt (nicht gelesen), sonst die tiefere Ebene. (l) **Nicht-Ziel „Volumes eines nicht mehr vorhandenen Containers“** (Satz 2 und Nicht-Ziele): nach `remove` ohne `--volumes` bleiben die Volumes (Lesestand), und ein folgendes `remove --volumes --yes` ohne Container endet mit Exit-Code `0` und entfernt keines von ihnen.
- [ ] **Doku:** Benutzerhandbuch (Abgrenzung zu `u-boot remove <dienst>` und `devcontainer down`, Beispiel, Warnung „löscht Daten“; ohne Kennungen), CHANGELOG, Hilfetext; die Kopien der Befehlsliste, die `remove` nennen müssen (Plan-Tabelle §3): README-Tabellen, `docs/maintainer/zusagen-und-messstaende.md`, die Stellen des Benutzerhandbuchs.

Pflichten der Closure (zählen nicht mit):

- [ ] `make gates` grün (im Devcontainer-Container, `AGENTS.md` §3.1).
- [ ] Review durchgeführt, Report unter `docs/reviews/` (`.harness/skills/reviewer.md`), kein Self-Review.
- [ ] Verification-Report unter `docs/reviews/`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapter/driving/cli/devcontainer.go` | update | Subkommando `remove`; `mapDevcontainerRemoveErrorToDiagnostic`; zentrale Klassifikation |
| `internal/hexagon/application/devcontainer_start.go` | update | `Remove` neben `Down` |
| `internal/hexagon/port/driving/`, `internal/hexagon/port/driven/` | update | Request/Response und Port-Operationen (die Dateien nennt der Implementer nach Lesen) |
| `internal/adapter/driven/docker/` | update | Entfernen von Container und Volumes, Auflistung der Mounts |
| Tests dazu und `internal/e2e` (Build-Tag `docker`) | update / neu | [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) |
| `docs/user/benutzerhandbuch.md` | update | Doku-Pflicht (`AGENTS.md` §5; keine Kennungen, Regel 6); Stellen, die gelesen wurden und `remove` berühren (Stand `7f5110cc`): Zeilen 59 und 60 (Voraussetzungs-Tabelle: Docker-Zeile nennt `devcontainer up`, `down`, `shell`), Zeile 65 („`devcontainer down` und `devcontainer shell` brauchen einen erreichbaren Docker-Daemon“), Zeile 164 („laufen im Image nicht“), Zeile 324 (`--quiet`: Liste der ändernden Befehle), Zeilen 845 bis 849 (Befehlsliste zu `up`, `down`, `shell`; ein Beispiel zu `remove`), Zeilen 1042 und 1043 („räumt keine auf; entfernen Sie … mit `docker rm` und `docker volume rm`“: nennt jetzt `remove --instance`), Zeile 1057 ff. („Stoppen, Neuanlage und Volumes“: Abgrenzung `down` / `remove` / `up --recreate`), Zeilen 2259 bis 2261 (Fehlersuche `volume is in use`: `docker rm`, dann `docker volume rm`); der Implementer liest weitere Treffer von `docker rm` und `docker volume rm` |
| `CHANGELOG.md` | update | Doku-Pflicht (`AGENTS.md` §5): Eintrag unter [Unreleased] zu `u-boot devcontainer remove` |
| `README.md`, `README.de.md` | update | Befehlstabelle: `README.md` Zeile 214 (`devcontainer down`: „it removes neither the container nor a volume“, nennt den entfernenden Befehl nicht) und Zeile 215 (`shell`), `README.de.md` Zeilen 218 und 219; eine Zeile `devcontainer remove` nach `shell` in beiden. Gelesen, unverändert: `README.md` Zeile 225 / `README.de.md` Zeile 229 (Satz zum Start des Devcontainers mit `up` und `shell`: betrifft das Starten) und `README.md` Zeile 222 / `README.de.md` Zeile 226 (Podman ≥ 4.0 „drop-in supported“: trägt keinen Befehlsnamen). Keine Kennungen (`AGENTS.md` §5 Regel 6) |
| `docs/maintainer/zusagen-und-messstaende.md` | update | Zeile 58 („Für `devcontainer up`, `devcontainer down` und `devcontainer shell` ist Podman als Host-Engine nicht zugesagt“): `devcontainer remove` ergänzen, weil die Zusage-Menge auch `remove` trägt; die Zeilen 59 und 60 zitieren einen README-Wortlaut („drop-in supported; not promised for `devcontainer up`/`down`/`shell`“), den `README.md` Zeile 222 nicht trägt — der Implementer berichtigt die Fundstelle, die er anfasst. Der Messstand zu `remove` kommt aus dem Docker-Test, nicht aus [`slice-v1-devcontainer-remove-messung`](../done/slice-v1-devcontainer-remove-messung.md) (`remove` ist dort nicht gefahren) |
| `docs/maintainer/devcontainer-arbeiten.md` | gelesen, unverändert | Tabelle der Make-Targets Zeilen 12 bis 16 (`dev-up`, `dev-shell`, `dev-down`): `remove` bekommt kein Make-Target in diesem Slice |
| `docs/plan/planning/open/README.md`, `docs/plan/planning/in-progress/roadmap.md` | update | Zeilen des Slice, Zustandsmarker |

Der Implementer liest die Dateien, ergänzt die Tabelle und gleicht sie vor dem Handoff gegen `git diff --stat` ab (`AGENTS.md` §3.10).

## 4. Trigger

**Start** (`next` → `in-progress`): [`slice-v2-spec-devcontainer-remove`](../in-progress/slice-v2-spec-devcontainer-remove.md) liegt in `done/`; `in-progress/` ist leer.

**Rückführungen:** `in-progress` → `next`, wenn die Port-Änderung mehr als einen Slice trägt; → `open`, wenn die Messung eine Operation nicht trägt, die der Vertrag verlangt.

## 5. Closure-Trigger

DoD vollständig, `make gates` grün, `make test-docker` gelaufen oder begründet nicht, Reports liegen vor, Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- `--volumes` löscht Daten, die der Nutzer nicht erwartet (Ollama-Modelle, Workspace im Sandbox-Profil). — **Ausgang:** offen; die Abfrage nennt jedes Volume, `--json` ohne `--yes` wird abgelehnt.
- Die Volume-Erkennung über die Mounts erfasst ein Volume, das ein anderer Container teilt. — **Ausgang:** offen bis Messung M2; der Testauftrag zu den Grenzen steht in §2 (zweite Zeile der DoD zu den Volumes, Review des CR F-8).

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `internal/adapter/driving/cli` (CLI), `internal/hexagon/application` (APP), `internal/hexagon/port` (PORT), `internal/adapter/driven` (ADP), `internal/e2e` (E2E) und `docs/plan/` (DOK); alle stehen in `harness/conventions.md` §Modus-Deklaration.

**Vorgelagert — offene Beobachtungen sichten:** vor dem Übergang nach `next/` auszuführen (Register `docs/plan/planning/observations/`); der Entwurf trägt noch kein Ergebnis.

**Modus-Begründungsblock:** alle berührten Sub-Areas sind Greenfield. Kein BF-/Hybrid-Block nötig.
