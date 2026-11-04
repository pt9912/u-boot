# Slice V1: `.devcontainer/devcontainer.env` als `--env-file` des Devcontainers — `generate devcontainer` legt die Datei und eine Vorlage an, `runArgs` bindet sie ein, `doctor` prüft die Ignorierung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle, Release-Zuordnung offen (es gibt keine Closure-Bedingung,
die von der DoD dieses Slice verschieden ist).

**Anlass.** Anfrage des Projektinhabers 2026-11-04: `.devcontainer/devcontainer.env`
soll unterstützt werden. Entscheidungen des Projektinhabers 2026-11-04 (Rückfrage
in der Sitzung): (1) Zweck ist Docker `--env-file`, direkt als Container-Umgebung,
ohne Umweg über `${localEnv:<NAME>}`; (2) `generate devcontainer` erzeugt zusätzlich
eine Vorlage mit Platzhaltern; (3) die echte Datei legt `generate devcontainer`
mit an, wenn sie fehlt, damit die `--env-file`-Zeile nie auf eine fehlende Datei
zeigt (`docker run --env-file` bricht bei fehlender Datei ab).

**Bezug:** [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) (Scope),
[`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) und
[`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)
(beide sagen heute: Zugangsdaten stehen in keiner erzeugten Datei — die Zusage wird
berührt, siehe §6), [`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz)
(Überschreibschutz der angelegten Datei).
Die ADR-Nummer vergibt der Architekt-Schritt über den ADR-Index (nächste freie
Nummer; der Stand `0043` ist der letzte vergebene).

**Berührte Spec-Stellen:** [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)
(Umgebung des Starts) · [`SPEC-053`](../../../../spec/spezifikation.md#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs)
(Reihenfolge der `runArgs`) · `spezifikation.md` §Devcontainer-Erzeugung (neuer
`SPEC-<NNN>` für die Datei; Kennung vergibt der Spec-Schritt nach dem ID-Schema).

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-11-04.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Ein Projekt führt Umgebungsvariablen des Devcontainers in
`.devcontainer/devcontainer.env` (`NAME=WERT` je Zeile). Der erzeugte
`devcontainer.json` trägt `--env-file=${localWorkspaceFolder}/.devcontainer/devcontainer.env`
in `runArgs`; Docker setzt die Variablen beim Start des Containers. Dazu:

- `u-boot generate devcontainer` legt `devcontainer.env` an, wenn sie fehlt
  (Kommentar-Kopf, keine Variablen), und fasst eine vorhandene Datei **nie** an —
  die Datei gehört dem Benutzer und ist kein Managed Block.
- `u-boot generate devcontainer` erzeugt `devcontainer.env.example` (verwaltet,
  Platzhalter wie `CHANGEME_…`, nie reale Secrets, [`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz)).
- Die `.gitignore` des Projekts ignoriert `.devcontainer/devcontainer.env`; die
  `.example`-Datei bleibt getrackt.
- `u-boot doctor` meldet `warn`, wenn `devcontainer.env` getrackt oder nicht
  ignoriert ist (Befund der Quelle von Zugangsdaten, wie die Prüfung in
  [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Kein Einlesen der Datei durch `u-boot` selbst** (kein Fallback für
  `${localEnv:NAME}`): Der Projektinhaber hat `--env-file` gewählt; ein
  `localEnv`-Fallback wäre ein zweiter Mechanismus mit eigener Rangfolge.
- **Kein Auflösen von Werten in Ausgaben:** `u-boot` liest den Inhalt der Datei nicht
  und gibt keinen Wert aus; `doctor` prüft nur Existenz, Git-Status und Ignorierung.
- **Keine Änderung der bestehenden `remoteEnv`-Referenzen** (`GIT_TOKEN`, Agent-Schlüssel):
  ihr Pfad über `${localEnv:<NAME>}` bleibt; die neue Datei ist ein zusätzlicher,
  benutzergeführter Weg.
- **Kein Unterstützen anderer Dateinamen oder mehrerer `--env-file`:** ein Name, ein Ort.
- **Kein Verhalten für die Dev-Container-CLI-Variante `envFile`:** die CLI kennt
  keinen solchen Schlüssel; der Weg ist `runArgs` (Fundstelle zu verifizieren, §6).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — ≤ 3 Liefer-Punkte.

- [ ] **Vertrag.** Neue ADR (`Accepted`, Index nachgezogen) zur Entscheidung
      `--env-file` mit angelegter Datei; Lastenheft-CR mit Änderungsvermerk und
      Versionszeile (Zusage „Zugangsdaten stehen in keiner erzeugten Datei“ in
      [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)
      und [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)
      auf die von `u-boot` erzeugten Dateien begrenzt; die benutzergeführte,
      ignorierte `devcontainer.env` ist ausgenommen); neuer `SPEC-<NNN>`. Der CR
      braucht die Bestätigung des Projektinhabers vor dem Code.
- [ ] **Erzeugung und Prüfung.** Golden-/Replay-Cases für `generate devcontainer`
      (Default- und Sandbox-Profil, mit und ohne vorhandene Datei, Datei bleibt
      byte-gleich, `--dry-run` schreibt nichts); `runArgs`-Reihenfolge gepinnt;
      `doctor`-Check mit Test; Docker-Test (Build-Tag `docker`): Variable aus der
      Datei ist im gestarteten Container gesetzt, auch über `u-boot devcontainer up`.
- [ ] **Doku.** Nutzer-Doku (`docs/user/`, beide READMEs) beschreibt Datei,
      Vorlage und `.gitignore` in Worten ohne Kennungen ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6);
      CHANGELOG nennt die Kennungen; Maintainer-Doku (`docs/maintainer/devcontainer-arbeiten.md`)
      folgt.
- [ ] `make gates` grün im Devcontainer-Container. *(nicht gezählt — Gate-Lauf)*
- [ ] Review, Verifikation, Closure-Notiz, Register, Risiko-Ausgänge,
      Drei-Paarungen. *(nicht gezählt — Closure-Pflichten)*

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — Pfad-Kandidaten für §8. Fundstellen sind am 2026-11-04
gelesen; der Implementer ergänzt die Tabelle im ersten Lauf (u. a. gegen
`git diff --stat`, [`AGENTS.md`](../../../../AGENTS.md) §3.10).

| Datei / Komponente | Art | Begründung / Fundstelle |
|---|---|---|
| `docs/plan/adr/00NN-…md`, `docs/plan/adr/README.md` | neu / update | Entscheidung `--env-file`, Datei-Anlage, Ignorierung; Index-Zeile |
| `spec/lastenheft.md` | update | CR: Ausnahme der Zusage in [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) und `-010`, Änderungsvermerk, Versionszeile |
| `spec/spezifikation.md` | update | neuer `SPEC-<NNN>`; [`SPEC-053`](../../../../spec/spezifikation.md#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs) (Reihenfolge `runArgs`) |
| `internal/hexagon/application/devcontainer_sandbox.go` | update | `devcontainerTemplateData`: beide Zweige (`!sandbox` Zeile 133, Sandbox danach) setzen `RunArgs`; die `--env-file`-Zeile kommt in beiden Profilen hinzu, Reihenfolge festlegen |
| `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl` | update (nur falls nötig) | `runArgs` wird aus `.RunArgs` gerendert (Zeilen 14 und 39); im Default-Profil ohne Resources entsteht heute kein `runArgs` — die Zeile ändert die Ausgabe **jedes** Projekts |
| `internal/hexagon/application/generate.go` | update | `planDevcontainerFiles` (Zeile 1053): Spezifikation der Dateien; neue Zeile für die Vorlage (verwaltet) und die Anlage der echten Datei (nur wenn fehlend, nie ersetzen); Phase-1/2-Atomarität bleibt |
| `internal/hexagon/application/templates/devcontainer/devcontainer.env.example.tmpl` | neu | Vorlage mit Platzhaltern |
| `internal/hexagon/application/templates/gitignore.tmpl` | update | Eintrag `.devcontainer/devcontainer.env`; `.gitignore` ist ganz verwaltet (`fileTemplates`, `templates.go`) |
| `internal/hexagon/application/devcontainer_start.go` | update (prüfen) | `containerNameOverlay` (Zeile 455) mischt `--name` in `runArgs` der Hilfs-Container-Kopie; `${localWorkspaceFolder}` muss im Hilfs-Container zum Host-Pfad auflösen — messen |
| `internal/hexagon/application/doctor.go` (neuer Check) | update | `warn` bei getrackter/nicht ignorierter Datei; Aufbau nach `doctor_sandbox.go` |
| `.devcontainer/devcontainer.json`, `.devcontainer/devcontainer.env.example`, `.gitignore` (dieses Repo) | update / neu | eigenes Erzeugnis von `u-boot generate devcontainer`; `runArgs` heute `["--privileged"]` |
| Tests: `internal/hexagon/application/generate_*_test.go`, `doctor_*_test.go`, `internal/e2e/devcontainer_*_docker_test.go`, Golden-Cases | neu / update | Happy/Boundary/Negativ, Replay nach [`harness/replay.md`](../../../../harness/replay.md) |
| `docs/user/*.md`, `README.md`, `README.de.md`, `CHANGELOG.md`, `docs/maintainer/devcontainer-arbeiten.md` | update | Doku; Nutzer-Doku ohne Kennungen |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Projektinhaber priorisiert den Slice und
bestätigt den Lastenheft-CR (Wortlaut der Ausnahme); kein anderer Slice in
`in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß): Vertrag, Generator und Doku sind drei
  Liefer-Punkte; wächst der Generator über den Docker-Test hinaus (z. B. durch
  eine Änderung der Hilfs-Container-Auflösung), teilt sich der Slice in
  CR-Slice und Code-Slice.
- `in-progress` → `open` (blockiert): Löst `${localWorkspaceFolder}` im
  Hilfs-Container von `u-boot devcontainer up` nicht zum Host-Pfad auf und
  trägt keine andere Form, kehrt der Slice mit der Messung als Beleg zurück.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln.

DoD vollständig (drei Liefer-Punkte plus Gate- und Closure-Pflichten); Review
und Verifikation; Closure-Notiz mit Lerneintrag; `git mv` nach `done/`.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Zusage-Konflikt:** [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)
  und [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)
  sagen, Zugangsdaten stünden in keiner erzeugten Datei. Die Ausnahme ist eine
  Vertragsänderung (Lockerung einer Sicherheitszusage) und braucht ADR und CR
  (AGENTS.md §3.6). — **Ausgang:** offen bis zum CR.
- **Ausgabe jedes Projekts ändert sich:** Das Default-Profil erzeugt heute ohne
  Resources/`ollama-host` kein `runArgs`; die neue Zeile ändert den Output aller
  Projekte (das Template-Datum nennt „byte-identisch“ für das UID-Default,
  `devcontainer_sandbox.go`, Doc-Kommentar von `devcontainerTemplateData`). Folge:
  `doctor` und Replay-Cases aller Profile. — **Ausgang:** offen.
- **Fehlende Datei nach Löschen:** Löscht der Benutzer `devcontainer.env`, bricht
  `docker run` ab. `doctor` meldet die fehlende Datei als `warn` mit Hinweis
  `u-boot generate devcontainer`. — **Ausgang:** offen.
- **Pfadauflösung im Hilfs-Container (`devcontainer up`) und im Editor:**
  `${localWorkspaceFolder}` ist zu messen (Hilfs-Container, VS Code, Colima).
  — **Ausgang:** offen → Messung im Slice, Beleg unter `docs/messungen/`.
- **Syntax der Datei:** `docker --env-file` kennt weder Anführungszeichen noch
  `export`; die vorgesehene Vorlage nennt das. — **Ausgang:** bei Closure.

## 7. Closure-Notiz

Noch nicht geschlossen; der Slice liegt in `open/` (Entwurf).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `spec/`, `docs/plan/` (DOK),
`internal/hexagon/application` (APP), `internal/e2e` (E2E) und `tools/` nicht; alle
sind deklarierte Sub-Areas ([`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area)),
keine neue Kandidatin.

**Vorgelagert — offene Beobachtungen sichten:** Register
[`observations/`](../observations/README.md) noch nicht durchgegangen — der
Implementer sichtet sie im ersten Lauf; bis dahin steht hier kein Befund.

**Modus-Begründungsblock:** entfällt — alle berührten Sub-Areas GF.
