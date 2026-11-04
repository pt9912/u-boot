# Slice V1: `.devcontainer/devcontainer.env` als `--env-file` des Devcontainers — `generate devcontainer` legt die Datei und eine Vorlage an, `runArgs` bindet sie ein, `doctor` meldet, wenn die Datei fehlt

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
ADR: [`ADR-0044`](../../adr/0044-devcontainer-umgebungsdatei-env-file.md).

**Berührte Spec-Stellen:** [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)
(Umgebung des Starts) · [`SPEC-053`](../../../../spec/spezifikation.md#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs)
(Reihenfolge der `runArgs`) · [`SPEC-055`](../../../../spec/spezifikation.md#spec-055--umgebungsdatei-des-devcontainers-pfad-inhalt-reihenfolge-in-runargs-anlage-und-prüfung)
(Umgebungsdatei des Devcontainers).

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
- `u-boot doctor` meldet `warn`, wenn `devcontainer.env` fehlt, die
  `devcontainer.json` die Zeile aber trägt (Quelle: [`SPEC-055`](../../../../spec/spezifikation.md#spec-055--umgebungsdatei-des-devcontainers-pfad-inhalt-reihenfolge-in-runargs-anlage-und-prüfung)
  und [`ADR-0044`](../../adr/0044-devcontainer-umgebungsdatei-env-file.md)).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Kein Einlesen der Datei durch `u-boot` selbst** (kein Fallback für
  `${localEnv:NAME}`): Der Projektinhaber hat `--env-file` gewählt; ein
  `localEnv`-Fallback wäre ein zweiter Mechanismus mit eigener Rangfolge.
- **Kein Auflösen von Werten in Ausgaben:** `u-boot` liest den Inhalt der Datei nicht
  und gibt keinen Wert aus; `doctor` prüft nur die Existenz (Quelle: [`SPEC-055`](../../../../spec/spezifikation.md#spec-055--umgebungsdatei-des-devcontainers-pfad-inhalt-reihenfolge-in-runargs-anlage-und-prüfung)).
- **Keine Änderung der bestehenden `remoteEnv`-Referenzen** (`GIT_TOKEN`, Agent-Schlüssel):
  ihr Pfad über `${localEnv:<NAME>}` bleibt; die neue Datei ist ein zusätzlicher,
  benutzergeführter Weg.
- **Kein Unterstützen anderer Dateinamen oder mehrerer `--env-file`:** ein Name, ein Ort.
- **Kein Verhalten für die Dev-Container-CLI-Variante `envFile`:** die CLI kennt
  keinen solchen Schlüssel; der Weg ist `runArgs` (Fundstelle zu verifizieren, §6).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — ≤ 3 Liefer-Punkte.

- [x] **Vertrag.** Neue ADR (`Accepted`, Index nachgezogen) zur Entscheidung
      `--env-file` mit angelegter Datei; Lastenheft-CR mit Änderungsvermerk und
      Versionszeile (Zusage „Zugangsdaten stehen in keiner erzeugten Datei“ in
      [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)
      und [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)
      auf die von `u-boot` erzeugten Dateien begrenzt; die benutzergeführte,
      ignorierte `devcontainer.env` ist ausgenommen); neuer `SPEC-<NNN>`. Der CR
      braucht die Bestätigung des Projektinhabers vor dem Code.
- [x] **Erzeugung und Prüfung.** Golden-/Replay-Cases für `generate devcontainer`
      (Default- und Sandbox-Profil, mit und ohne vorhandene Datei, Datei bleibt
      byte-gleich, `--dry-run` schreibt nichts); `runArgs`-Reihenfolge gepinnt;
      `doctor`-Check mit Test; Docker-Test (Build-Tag `docker`): Variable aus der
      Datei ist im gestarteten Container gesetzt, auch über `u-boot devcontainer up`.
- [x] **Doku.** Nutzer-Doku (`docs/user/`, beide READMEs) beschreibt Datei,
      Vorlage und `.gitignore` in Worten ohne Kennungen ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6);
      CHANGELOG nennt die Kennungen; Maintainer-Doku (`docs/maintainer/devcontainer-arbeiten.md`)
      folgt.
- [x] `make gates` grün im Devcontainer-Container. *(nicht gezählt — Gate-Lauf)*
- [x] Review ([Report](../../../reviews/2026-11-04-slice-v1-devcontainer-env-datei.md)), Verifikation ([Report](../../../reviews/verification-2026-11-04-slice-v1-devcontainer-env-datei.md)), Closure-Notiz, Register, Risiko-Ausgänge,
      Drei-Paarungen. *(nicht gezählt — Closure-Pflichten)*

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — Pfad-Kandidaten für §8. Fundstellen sind am 2026-11-04
gelesen; der Implementer ergänzt die Tabelle im ersten Lauf (u. a. gegen
`git diff --stat`, [`AGENTS.md`](../../../../AGENTS.md) §3.10).

| Datei / Komponente | Art | Begründung / Fundstelle |
|---|---|---|
| `docs/plan/adr/0044-devcontainer-umgebungsdatei-env-file.md`, `docs/plan/adr/README.md` | neu / update (vor dem Code übernommen) | Entscheidung `--env-file`, Datei-Anlage, Ignorierung; Index-Zeile |
| `spec/lastenheft.md`, `spec/spezifikation.md` | update (vor dem Code übernommen) | CR 0.35.0 ([`LH-FA-DEV-012`](../../../../spec/lastenheft.md#lh-fa-dev-012--umgebungsdatei-des-devcontainers), Ausnahmen in `-009`/`-010`); [`SPEC-055`](../../../../spec/spezifikation.md#spec-055--umgebungsdatei-des-devcontainers-pfad-inhalt-reihenfolge-in-runargs-anlage-und-prüfung) samt Beleg-Tabelle, [`SPEC-053`](../../../../spec/spezifikation.md#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs) und weitere Stellen nachgezogen |
| `internal/hexagon/application/devcontainer_sandbox.go` | update | Konstanten `devcontainerEnvFileRel`, `devcontainerEnvFileExampleRel`, `envFileRunArg`; `devcontainerTemplateData`: beide Zweige beginnen `RunArgs` mit `envFileRunArg`; Doc-Kommentar berichtigt |
| `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl` | unverändert | gelesen: `runArgs` wird aus `.RunArgs` gerendert (Zeilen 14 und 39), nur bei nicht leerer Liste; die Liste ist jetzt nie leer, die Vorlage braucht keine Änderung |
| `internal/hexagon/application/generate.go` | update | `planDevcontainerFiles` nimmt die Beispieldatei in die Spezifikation auf und hängt den Plan von `planDevcontainerEnvFile` an; neue Funktion `planDevcontainerEnvFile` (nur `Exists`, nie `ReadFile`) liefert `devcontainerFilePlan` mit Aktion `Write` oder `NoOp`; `executeDevcontainerPlans` unverändert |
| `internal/hexagon/application/templates/devcontainer/devcontainer.env.tmpl`, `…/devcontainer.env.example.tmpl` | neu | Kommentar-Kopf der Datei (kein Block); verwaltete Vorlage mit `CHANGEME_`-Werten |
| `internal/hexagon/application/templates/gitignore.tmpl`, `internal/adapter/driven/externaltemplates/templates/basic/.gitignore.tmpl` | update | Zeile `.devcontainer/devcontainer.env`; die Vorlage `basic` ist byte-gleich zum Standardpfad (`TestInitProjectService_FromTemplate_BasicByteIdenticalToDefaultPath`) |
| `internal/hexagon/application/templates.go`, `internal/hexagon/application/initproject.go` | update | `fileTemplate.UserOwned`; `actionKeepExisting`: `planFile` lässt eine vorhandene benutzergeführte Datei stehen (kein Lesen, kein `.bak`, kein Abbruch), `executeTemplatedFiles` überspringt sie; `devcontainerFileTemplates` führt Beispieldatei und Umgebungsdatei |
| `internal/hexagon/application/doctor.go`, `internal/hexagon/application/doctor_envfile.go` | update / neu | `checkIDDevcontainerEnvFile` in `DoctorCheckIDs` und `Check`; Prüfung nach dem Muster von `doctor_sandbox.go` (liest nur `devcontainer.json`, von der Umgebungsdatei nur `Exists`); `credentialScanFiles` unverändert (`doctor_sandbox.go` gelesen, Zeile 60: die Datei steht nicht darin) |
| `internal/adapter/driving/cli/jsontestutil/coderegistry.go` | update | Code `devcontainer.envfile` in `DefaultAllowedCodes` |
| `internal/hexagon/application/devcontainer_start.go` | unverändert | gelesen: `containerNameOverlay` (Zeile 461) hängt nur `--name=` an `runArgs` an und lässt die übrigen Einträge stehen; die Auflösung von `${localWorkspaceFolder}` im Hilfs-Container belegt `TestE2E_EnvFile_VariableReachesContainerViaDevcontainerUp` |
| Tests `internal/hexagon/application`: `generate_envfile_test.go`, `generate_envfile_helpers_test.go`, `initproject_envfile_test.go`, `doctor_envfile_test.go` | neu | Belege der Beleg-Tabelle von [`SPEC-055`](../../../../spec/spezifikation.md#spec-055--umgebungsdatei-des-devcontainers-pfad-inhalt-reihenfolge-in-runargs-anlage-und-prüfung) (Happy/Boundary/Negativ) |
| Tests `internal/adapter/driving/cli/doctor_envfile_test.go` | neu | `TestDoctor_EnvFile_QuietDropsWarning` |
| Tests `generate_test.go`, `generate_sandbox_test.go`, `generate_host_socket_test.go`, `generate_ollama_host_test.go`, `generate_ollama_test.go`, `generate_containerenv_test.go`, `generate_toolchain_test.go`, `dryrun_log_test.go`, `doctor_test.go`, `initproject_test.go`, `template_init_test.go` | update | `runArgs`-Pins, Goldens, Dateilisten und Zähler an die neue Zeile und die zwei neuen Dateien angepasst |
| `internal/e2e/envfile_docker_test.go`, `internal/e2e/sandbox_devcontainer_docker_test.go` | neu / update | Docker-Tests (Build-Tag `docker`); `generatedRunArgs` löst `${localWorkspaceFolder}` zum Projektordner auf, die Egress-Erwartung zählt die Zeile mit |
| `.devcontainer/devcontainer.json`, `.devcontainer/devcontainer.env.example`, `.gitignore` (dieses Repo) | update / neu | eigenes Erzeugnis von `u-boot generate devcontainer` (Binary aus `make build`, zweiter Lauf NoOp); `runArgs` trägt die Zeile vor `--privileged`; die Zeile `.devcontainer/devcontainer.env` in der `.gitignore` von Hand; `.devcontainer/devcontainer.env` selbst ist ignoriert und nicht getrackt (kurzzeitig getrackt mit `c7fae976`, auf Anweisung des Projektinhabers wieder aus dem Index genommen) |
| `harness/replay.md` | update | neue Zeile für die Umgebungsdatei; Zeilen mit festen `runArgs` nennen die Zeile der Umgebungsdatei |
| `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md` | update | Doku; Nutzer-Doku ohne Kennungen |
| `docs/maintainer/devcontainer-arbeiten.md` | update | Hinweis zum frischen Klon: die Datei ist ignoriert und fehlt im Klon; einmal `u-boot generate devcontainer` vor `make dev-up`/`dev-shell`, für dieses Repo wie für Nutzerprojekte |
| `docs/plan/cr/cr-vorschlag-devcontainer-env-datei.md` | neu / update | CR-Vorschlag; Status „freigegeben wie vorgelegt“, Block „Abweichungen der Übernahme je Stelle“ (Stellen (c), (d)) |
| `docs/messungen/slice-v1-devcontainer-env-datei/`, `docs/messungen/README.md`, `tools/messung/devcontainer-env-datei/messen.sh` | neu / update | Messung der `--env-file`-Eigenschaften (Protokoll, Rohdaten, Skript, Bestand-Zeile) |
| `docs/plan/planning/in-progress/roadmap.md`, `docs/plan/planning/open/README.md` | update | Roadmap-Zeile statt Ruhe-Marker (Start); Bestand-Zeile des Slice an den Ist-Stand |
| `docs/reviews/2026-11-04-slice-v1-devcontainer-env-datei.md`, `docs/reviews/verification-2026-11-04-slice-v1-devcontainer-env-datei.md` | neu | Review- und Verification-Report |
| `docs/plan/planning/observations/BEO-E2E/docker-tests-pinnen-vertrag-ausserhalb-des-gates/`, `…/BEO-DOK/zusage-weiter-als-messung/`, `…/BEO-DOK/adr-verweist-abwaerts-auf-plan/` | neu / update | Register: eine neue Beobachtung, zwei weitere Belege |
| `internal/e2e/sandbox_devcontainer_docker_test.go` (Nachzug nach dem Volllauf) | update | `TestE2E_SandboxDevcontainer_NestedDocker` pinnte `runArgs = [--privileged]`; erwartet jetzt die Zeile der Umgebungsdatei vor `--privileged` |

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
  (AGENTS.md §3.6). — **Ausgang:** entfallen: der CR ist übernommen (Lastenheft 0.35.0, Freigabe-Zelle der Historie); die benutzergeführte Datei ist ausgenommen, `u-boot` liest sie nicht (gepinnt: `TestDoctor_EnvFile_NeverReadsContent`, `TestGenerateDevcontainer_EnvFile_ExistingNeverTouched`).
- **Ausgabe jedes Projekts ändert sich:** Das Default-Profil erzeugt heute ohne
  Resources/`ollama-host` kein `runArgs`; die neue Zeile ändert den Output aller
  Projekte (das Template-Datum nennt „byte-identisch“ für das UID-Default,
  `devcontainer_sandbox.go`, Doc-Kommentar von `devcontainerTemplateData`). Folge:
  `doctor` und Replay-Cases aller Profile. — **Ausgang:** eingetreten und benannt ([`ADR-0044`](../../adr/0044-devcontainer-umgebungsdatei-env-file.md) Konsequenz (a), [`SPEC-055`](../../../../spec/spezifikation.md#spec-055--umgebungsdatei-des-devcontainers-pfad-inhalt-reihenfolge-in-runargs-anlage-und-prüfung), CHANGELOG, `harness/replay.md`); die Goldens und rund 28 Tests sind angepasst. Ein Docker-Test (`TestE2E_SandboxDevcontainer_NestedDocker`) pinnte die alte Liste und blieb bis zum Volllauf von `make test-docker` unbemerkt, weil `make gates` ihn nicht fährt — nachgezogen, geführt in [`BEO-E2E/docker-tests-pinnen-vertrag-ausserhalb-des-gates`](../observations/BEO-E2E/docker-tests-pinnen-vertrag-ausserhalb-des-gates/observation.md).
- **Fehlende Datei nach Löschen:** Löscht der Benutzer `devcontainer.env`, bricht
  `docker run` ab. `doctor` meldet die fehlende Datei als `warn` mit Hinweis
  `u-boot generate devcontainer`. — **Ausgang:** entfallen: `generate` legt eine gelöschte Datei neu an (`TestGenerateDevcontainer_EnvFile_DeletedIsRecreated`), `doctor` meldet `warn` (`TestDoctor_EnvFile_WarnsWhenMissing`).
- **Pfadauflösung im Hilfs-Container (`devcontainer up`) und im Editor:**
  `${localWorkspaceFolder}` ist zu messen (Hilfs-Container, VS Code, Colima).
  — **Ausgang:** entfallen für den Hilfs-Container: `TestE2E_EnvFile_VariableReachesContainerViaDevcontainerUp` ist im Volllauf von `make test-docker` am 2026-11-04 `--- PASS` (4,80 s, verschachtelte Engine, ein Lauf); für den Editor nicht automatisierbar, im Lastenheft als nicht zugesagt geführt; der Lauf auf einer nicht verschachtelten Engine steht aus (siehe nächstes Risiko).
- **Syntax der Datei:** `docker --env-file` kennt weder Anführungszeichen noch
  `export`; die vorgesehene Vorlage nennt das. — **Ausgang:** entfallen: gemessen ([Protokoll](../../../messungen/slice-v1-devcontainer-env-datei/README.md): Anführungszeichen sind Teil des Werts, `export` bricht ab), der Kommentar-Kopf der angelegten Datei nennt es.
- **Die Zusage zum Start über `u-boot devcontainer up` stützt sich auf zwei Docker-Tests in der verschachtelten Engine** (Review F-3, Verifikation V-5). — **Ausgang:** weiter offen → [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md); der Lauf `integration-docker` auf einer nicht verschachtelten Engine kommt nach dem Push.
- **Frischer Klon ohne die Datei** (Review F-1): `docker run --env-file` bricht ab (Exit `125`, gemessen). — **Ausgang:** eingetreten, als Folge akzeptiert (Entscheidung des Projektinhabers 2026-11-04): [`ADR-0044`](../../adr/0044-devcontainer-umgebungsdatei-env-file.md) Konsequenz (f), Hinweis im Handbuch und in `docs/maintainer/devcontainer-arbeiten.md`, auch für dieses Repo (die Datei ist ignoriert und nicht getrackt; eine kurzzeitig getrackte Fassung, `c7fae976`, ist auf Anweisung des Projektinhabers wieder aus dem Index genommen).
- **Eine getrackte Datei wäre von der `.gitignore`-Zeile nicht geschützt** (Verifikation V-4). — **Ausgang:** entfallen: die Datei dieses Repos ist nicht getrackt, die Zeile greift.

## 7. Closure-Notiz

- **Was hat funktioniert:** Die Folge Vertrag → Entscheidung → Umsetzung hielt: Lastenheft 0.35.0, [`ADR-0044`](../../adr/0044-devcontainer-umgebungsdatei-env-file.md) und [`SPEC-055`](../../../../spec/spezifikation.md#spec-055--umgebungsdatei-des-devcontainers-pfad-inhalt-reihenfolge-in-runargs-anlage-und-prüfung) standen vor dem Code; alle 22 Belege der Beleg-Tabelle haben einen existierenden Test (Verifikation), Mutationsproben des Implementers wurden rot. `make gates` grün (Verifier-Lauf auf `615e8b3b`: Coverage 92,10 %); die zwei Docker-Tests zur Umgebungsdatei sind im Volllauf von `make test-docker` `--- PASS`.
- **Was ging anders als geplant:** (1) Der Review fand zwei HIGH-Befunde: der frische Klon (F-1, vom Projektinhaber als akzeptierte Folge entschieden: Hinweis im Handbuch, auch für dieses Repo) und Verweise abwärts in der ADR (F-2, in place berichtigt, weil die ADR nie veröffentlicht war). (2) Die Zusage zum Start über `up` war im Lastenheft zunächst „nicht zugesagt“; auf Anweisung des Projektinhabers („Lastenheft ist Zielbild“) trägt [`LH-FA-DEV-012`](../../../../spec/lastenheft.md#lh-fa-dev-012--umgebungsdatei-des-devcontainers) jetzt die Formel von [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Zusage-Menge, nur soweit ein Docker-Test es belegt); das Lastenheft 0.35.0 ist dabei in place geändert worden, weil es nicht veröffentlicht war. (3) Der Volllauf von `make test-docker` fand einen Docker-Test, der die alte `runArgs`-Liste pinnte (`TestE2E_SandboxDevcontainer_NestedDocker`); er ist nachgezogen und gezielt grün (64,8 s), ein erneuter Volllauf steht aus. (4) Der Plan-Abgleich nach der Fix-Runde (Verifikation V-1) hat vier Zeilen berichtigt. (5) Die Datei `.devcontainer/devcontainer.env` dieses Repos war kurzzeitig getrackt (`c7fae976`, `git add -f` gegen die Ignore-Zeile); der Projektinhaber stellte klar, dass sie ignoriert sein soll — sie ist wieder aus dem Index genommen, der Klon-Hinweis gilt auch für dieses Repo.
- **Gate-Nachweis des Docker-Tests:** Volllauf am 2026-11-04 mit 3738 `PASS` und einem `FAIL` (der oben genannte Test) auf dem Stand vor dem Nachzug; danach nur der gezielte Lauf des einen Tests. Die zwei `TestE2E_Resources_*` aus dem ersten Lauf dieses Tages waren dort grün, weil die cgroup dieses Containers zuvor von Hand umgebaut war (siehe [`slice-v1-sandbox-nested-docker-cgroup`](../done/slice-v1-sandbox-nested-docker-cgroup.md)).
- **Steering-Loop-Eintrag:** Ein Vertragswert (`runArgs`), den Docker-getaggte Tests pinnen, wird von `make gates` nicht geprüft; eine Änderung des Vertragswerts lässt solche Tests unbemerkt rot, bis der Volllauf läuft — gezählt, nicht verkörpert (1×, [`BEO-E2E/docker-tests-pinnen-vertrag-ausserhalb-des-gates`](../observations/BEO-E2E/docker-tests-pinnen-vertrag-ausserhalb-des-gates/observation.md)).
- **Beobachtungs-Register (`../observations/`):** `BEO-E2E/docker-tests-pinnen-vertrag-ausserhalb-des-gates/` neu (Beleg `evidence/slice-v1-devcontainer-env-datei.md`, 1×); `evidence/slice-v1-devcontainer-env-datei.md` in `BEO-DOK/zusage-weiter-als-messung/` ergänzt (F-3, F-4 der Lead-Klasse) und in `BEO-DOK/adr-verweist-abwaerts-auf-plan/` (F-2).
- **Folge-Slices:** [`slice-v1-sandbox-nested-docker-cgroup`](../done/slice-v1-sandbox-nested-docker-cgroup.md) ist eine Datei in `open/`.
- **Risiken aus §6:** je ein Ausgang, siehe §6.
- **Drei Paarungen:** Anker — kein Herkunfts-Anker gesetzt (nichts verkörpert); Folge-Slice — die Datei in `open/` liegt vor; Register — die Evidence-Dateien liegen vor.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `spec/`, `docs/plan/` (DOK),
`internal/hexagon/application` (APP), `internal/e2e` (E2E) und `tools/` nicht; alle
sind deklarierte Sub-Areas ([`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area)),
keine neue Kandidatin.

**Vorgelagert — offene Beobachtungen sichten:** Register
[`observations/`](../observations/README.md) durchgegangen am 2026-11-04 (nach Titeln): Treffer [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (Review F-3 und F-4), [`BEO-DOK/adr-verweist-abwaerts-auf-plan`](../observations/BEO-DOK/adr-verweist-abwaerts-auf-plan/observation.md) (F-2), [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (Plan-Abgleich nach der Fix-Runde, V-1); sonst keine Treffer.

**Modus-Begründungsblock:** entfällt — alle berührten Sub-Areas GF.
