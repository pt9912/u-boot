# Verification-Report: slice-v1-devcontainer-container-name — 2026-11-02

**Rolle:** Verifier (frischer Kontext; Baseline-Regelwerk `modul-11`)
**Gegenstand:** [`slice-v1-devcontainer-container-name`](../plan/planning/done/slice-v1-devcontainer-container-name.md) — Slice-Commits `8912f7fe` (Messprotokoll/Messskript) und `9b349879`, `0f6a219e`, `95b5c864`, `b5be7aaa`, `ff56fd50`, `b919d4b4`, `94564e5b`, `0ff780ba`, `be98839f`; Abgleichstand **nach der letzten Fix-Runde = `be98839f`**, Repo-Tip `71cf2f23` (sauber).
**Ausgeschlossene Commits** (getrennter Release-Cut-Slice, nicht Gegenstand): `e0259b97`, `e68786eb`, `8cd0976b`, `71cf2f23`.

**Eingangs-Kontext:**

- Slice-Plan §1 (Ziel), §2 (DoD), §3 (Plan-Tabelle + Wächter-Tabelle), §6 (Fragen 1–9 entschieden, Messungen M1–M6, Risiken), §7 (Closure offen)
- Lastenheft **0.29.0**: [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (CR-Wortlaut), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-FA-INIT-002`](../../spec/lastenheft.md#lh-fa-init-002--projektname), [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache), [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe), [`LH-FA-ARCH-002`](../../spec/lastenheft.md#lh-fa-arch-002--schichten-und-verzeichnislayout)/[`003`](../../spec/lastenheft.md#lh-fa-arch-003--import-regeln-und-enforcement); [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)
- [ADR-0038](../plan/adr/0038-devcontainer-instanz-zweites-id-label.md) (`Accepted`, im Diff nicht berührt), Messprotokoll [`docs/messungen/slice-v1-devcontainer-container-name/README.md`](../messungen/slice-v1-devcontainer-container-name/README.md)
- Review-Report [`docs/reviews/2026-11-01-slice-v1-devcontainer-container-name.md`](2026-11-01-slice-v1-devcontainer-container-name.md) mit F-1..F-7
- [`harness/verification.md`](../../harness/verification.md), [`AGENTS.md`](../../AGENTS.md) (§3.1, §3.10), [`harness/roles.md`](../../harness/roles.md)

---

## Scope

- Slice: [`slice-v1-devcontainer-container-name`](../plan/planning/done/slice-v1-devcontainer-container-name.md)
- IDs: [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Lastenheft 0.29.0, CR), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-FA-INIT-002`](../../spec/lastenheft.md#lh-fa-init-002--projektname), [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe), [`LH-FA-ARCH-002`](../../spec/lastenheft.md#lh-fa-arch-002--schichten-und-verzeichnislayout)/[`003`](../../spec/lastenheft.md#lh-fa-arch-003--import-regeln-und-enforcement), [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache), [`ADR-0038`](../plan/adr/0038-devcontainer-instanz-zweites-id-label.md) (unverändert), [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)
- Artefakte: `internal/hexagon/domain/containername.go` (+ Test), `internal/hexagon/port/driven|driving/devcontainer{,_launcher}.go` (+ Tests), `internal/hexagon/application/devcontainer_start.go` (+ Test), `internal/adapter/driven/docker/devcontainer_launcher.go` (+ Tests), `internal/adapter/driving/cli/devcontainer.go` (+ Tests), `internal/e2e/` (fünf docker-getaggte Dateien, `devcontainer_name_docker_test.go` neu), Messprotokoll + `tools/messung/devcontainer-container-name/messen.sh`, `spec/lastenheft.md`, `spec/spezifikation.md`, `docs/user/benutzerhandbuch.md`, `CHANGELOG.md`

## DoD-Abgleich

DoD-Quelle: Plan §2. Gates und Closure-Pflichten zählen laut Plan-Kopf nicht zu den Liefer-Punkten; die vier DoD-Zeilen sind hier je einzeln belegt.

| # | DoD-Punkt | Zustand | Evidence |
|---|---|---|---|
| 1 | **Messung und Freigabe** | **erbracht** | Protokoll [`docs/messungen/slice-v1-devcontainer-container-name/README.md`](../messungen/slice-v1-devcontainer-container-name/README.md) (Commit `8912f7fe`): datiert 2026-11-01 (Läufe 17:27–17:33 UTC), Umgebung vollständig (Linux `amd64`, Docker 29.9.2, CLI 0.89.0, `u-boot` 0.17.0-dev), je Messung zwei Läufe (`rohdaten-<modus>-lauf{1,2}.txt`, je Modus je zwei), M1–M6 samt Zusatzprobe, verworfene Vorläufe von M4/M6 ausgewiesen (ehrliche Messdisziplin). Grenzen benannt (Länge oberhalb 257 nicht gesucht; Kollisions-Wortlaut hinter der 260-Zeichen-Kürzung). Skript `tools/messung/devcontainer-container-name/messen.sh` (538 Zeilen) committet, Ablage konform zur [Konvention](../../harness/conventions.md#ablage-von-messprotokollen-und-skripten). Freigabe: §6 — Namensform **A**, Weg **S1'**, Kollisions**abbruch**, Ausgabe **Name mit ausgeben**, CR-Wortlaut **freigegeben** (Fragen 1–5, 7–9, Projektinhaber 2026-11-01); der CR im Lastenheft trägt den freigegebenen Wortlaut im Ist-Zustand („Ein Container ohne einen solchen Namen behält …“, ohne die Chronik-Formulierung „vor dieser Fassung“). |
| 2 | **Der Name** | **erbracht** — mit einer begrenzten Restnotiz (unten) | **Name am Container, mit/ohne Instanz, `--recreate`:** e2e `TestDevcontainerUp_Docker_ContainerName` (`internal/e2e/devcontainer_name_docker_test.go`, Build-Tag `docker`) — Name am Container (`docker inspect`), Übernahme bei gleichem Ordner, `--recreate` nennt ihn neu, `devcontainer.json` byte-identisch; grün im CI-Lauf **41340432215**. Unit-Pins: `TestDevcontainerUp_ComposesContainerName`/`…WithInstance`/`…FromFolder` (Bildungsregel `u-boot-<projekt>[-<instanz>]`, Ordner-Basisname ohne `u-boot.yaml`), `TestComposeDevcontainerContainerName` (+ leeres Projekt nicht zusammensetzbar). **`down`/`shell` finden ihn weiter über die Labels:** `inspector_test.go` (`TestInspector_FindContainer*`, Argv-Filter auf Labels, unverändert im Diff) plus der grüne e2e-Standard (`TestDevcontainerShell_Docker_Session` und die Down-Instanz-Tests laufen **gegen namentragenden Containern** — jeder heute gestartete Container trägt den Namen) im selben CI-Lauf; `FindContainer` ist im Diff unberührt. **`devcontainerId`/Volumes gleich:** Messung M3 und M6 (paarweise identische Werte in vier Zellen, [ADR-0038](../plan/adr/0038-devcontainer-instanz-zweites-id-label.md) Entscheidung 3 gemessen widerlegt) — DoD erlaubt „Test oder Messung“; die byte-identische erzeugte Datei zusätzlich e2e-gepinnt. **Ungültiger/zulanger Bestandteil → Exit 10:** `TestDevcontainerUp_RejectsUnusableContainerName` (Sentinel `ErrDevcontainerContainerNameInvalid`), Pin-Tabelle `TestExecute_Devcontainer_SentinelsPerSubcommand` (Zeilen „container name“/„container name conflict“, Exit `10` je Pin je Subkommando `up` — F-1-Fix). Keine Längen-Obergrenze zugesagt: Projekteiteil ≤ 63 ([`LH-FA-INIT-002`](../../spec/lastenheft.md#lh-fa-init-002--projektname), `projectname.go`), Instanzteil ≤ 32 (`ValidateDevcontainerInstance`, unverändert), Präfix + zwei Bindestriche ⇒ Zusammensetzung ≤ 103 Zeichen, unter den bis 257 gemessenen akzeptierten Namen; die Längengrenze nach oben ist im Protokoll als nicht gesucht ausgewiesen und die Zusage in Spec/CR entsprechend nur auf das Gemessene begrenzt. **Kollision nach Frage 5:** e2e `TestDevcontainerUp_Docker_ConflictingContainerName` (Exit `10` über die CLI, haltender Container läuft weiter), Unit `TestLauncher_Up_ConflictingContainerName` (fremd/gleich/ohne Label — F-3-Fix/scheiternde Befragung/kein Name fragt die Engine nicht). |
| 3 | **Doku und Zusage-Menge** | **erbracht** | Lastenheft 0.29.0: CR im Absatz zu `up` in [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Fehlerklassen-Erweiterung in [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) („Bestandteil des Containernamen verletzt die Engine-Regel“, „Container eines anderen Projektordners hält den Containernamen“ ⇒ `10`), Historie 0.29.0. [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes): drei neue Blöcke (Containername, Weg in die Konfiguration, Namenskollision vor dem Start), Ausgabetabelle `10` erweitert, `up`-Zeile `[name <Containername>]`, `shell`-Ausgabe/`--quiet` nachgezogen. Handbuch `docs/user/benutzerhandbuch.md`: Abschnitt „Der Name des Containers“ — **0 Treffer** auf `LH-`/`SPEC-`/`ADR-`/`spec/`/`slice-`/`welle-`/`roadmap`/`in-progress` über die ganze Datei (geprüft per Suche; [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache)-Konvention §5 Regel 6). `CHANGELOG.md`: Eintrag mit Promised/Not-promised-Begrenzung (Linux-Host `amd64`, lokale Engine) und Release-Zuordnung Minor/*Changed* (Frage 8). Bedingungen: die Code-zugespitzten Bedingungen der Spec tragen pinnende Tests (Kollision ohne `local_folder`-Label → F-3-Test; „keine mergebare Struktur ⇒ kein Name“ → `TestDevcontainerUp_HandsNoNameForAnUnreadableConfiguration`; Name in der Zeile von `up` und `shell` → CLI-Tests samt `--quiet` und `--json`-Unverändertheit „no field by the name“, `devcontainer_test.go` Zeile 574). Restnotiz: zwei Spezial-Bedingungen sind **mess- statt testgetragen** — siehe Rest-Risiken. |
| 4 | `make gates` grün; Review; Verification; Closure-Notiz mit Lerneintrag; Register fortgeschrieben; Risiko-Ausgänge; drei Paarungen | **teils erbracht — die übrigen Punkte sind Closure-Pflichten, nicht Lücken des Baues** | `make gates`: **grün**, EXIT=0, Docker-Modus am Host (lauf am 2026-11-02); Stempel `.harness/state/gates-passed.{head,diffsha}` = `71cf2f23…` = HEAD (d-check 1007 Dateien / 0 Befunde, baseline-verify v6.18.1 OK, coverage 92,20 % ≥ 90 %, ablage-pruefung OK, 26 Faelle/0 Fehlschlaege). Review: Report 2026-11-01 vorliegend; F-1..F-4 behoben (Details unten). Verification: **dieser Report**. **Offen bis Closure:** §7 Closure-Notiz mit Lerneintrag, Beobachtungs-Register-Fortschreibung, Risiko-Ausgänge in der Closure-Notiz, drei Paarungen (Anker · Folge-Slice · Register, Baseline `modul-05` §Closure) — alles Schritte des Closure-Schritts selbst, der mit den ersten drei DoD-Zeilen und diesem Report sein Belegset hat; der Slice liegt weiter in `in-progress/` (korrekt, nicht widerrechtlich geschlossen). |

**Vorab-Ausgänge der §6-Risiken für die Closure-Notiz** (Beleg je Zeile; die Ausgänge selbst trägt die Closure-Notiz):

1. CLI kennt keine Namensoption — Ausgang: Weg **S1'** über `--override-config` (M6 gemessen, umgesetzt).
2. Name ändert Identität — Ausgang: M3/M6 widerlegt, zusätzlich byte-identische erzeugte Datei e2e-gepinnt.
3. Kollision — Ausgang: Abbruch Exit `10` mit klarer Meldung vor dem Hilfs-Container (e2e + Unit).
4. Zusage weiter als die Messung — Ausgang: CR/CHANGELOG auf Linux-Host `amd64`, lokale Engine begrenzt.
5. Bedingung ohne Test — Ausgang: mehrheitlich gepinnt; zwei mess-getragene Bedingungen als Rest benannt (unten).
6. Docker-Test nicht lokal — Ausgang: CI-Lauf **41340432215** grün.

## Review-Findings F-1..F-7 (Report 2026-11-01)

| ID | Klasse | Zustand | Beleg |
|---|---|---|---|
| F-1 | HIGH | **behoben** | `internal/adapter/driving/cli/devcontainer_test.go` Zeilen 444–448: Pin-Tabelle trägt „container name“ (`ErrDevcontainerContainerNameInvalid`) und „container name conflict“ (`…Conflict`) mit Exit `10`, Diagnose-Quelle [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Subkommando `up` — beide Sentinels je Pin in beiden Klassifikatoren geprüft. Der `-Invalid`-Sentinel ist damit auf Exit-Code-Ebene gepinnt, der `-Conflict`-Sentinel zusätzlich im e2e (CI 41340432215). |
| F-2 | MEDIUM | **behoben** | Plan-Zeile 21 trägt `docs/plan/planning/in-progress/roadmap.md` (Arbeitszeile `95b5c864`, Ruhe-Marker bei Closure) und die open/README-Bestand-Zeile; Plan-Geändert-Diff `b5be7aaa..HEAD` zeigt die Zeile. |
| F-3 | MEDIUM | **behoben** | `TestLauncher_Up_ConflictingContainerName`/„holder without a local_folder label is foreign“ mit fakeDocker-`ps-output` `"abc123…\t\n"` (kein Label) ⇒ `ErrDevcontainerContainerNameConflict`, kein Hilfs-Container — die [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)-Bedingung ist gepinnt. |
| F-4 | LOW | **behoben** | Plan-Diff zeigt `harness/replay.md` auf „**unverändert** (S2 verworfen, S1')“ und `docs/plan/adr/` auf „**unverändert** — kein Folge-ADR“ gestellt; beide Zeilen haben im Diff keinen Eintrag (unten, Plan-vs-Code). |
| F-5 | INFO | nicht behoben, nicht gefordert | Der Doppelprüfpfad im Launcher (zweite `ValidateDockerContainerName`, Klasse `rejectedInput` ⇒ `11` für einen künftigen Direktaufrufer) bleibt dokumentiert in Spec („der Launcher wendet sie … noch einmal an“) — kein DoD-Gegenstand. |
| F-6 | INFO | nicht behoben (mess-getragen), **als Restpunkt benannt** | Siehe Rest-Risiken. |
| F-7 | INFO | nicht behoben, nicht gefordert | Der Kommentar von `mapDevcontainerDownErrorToDiagnostic` (`devcontainer.go` 646–649) zählt weiter nur die Klassen, die `down` ohnehin nicht trifft; `down` kann keine Namen-Sentinels erzeugen — nicht falsch, nur enger als der Stand. |

## Plan-vs-Code (§3 gegen die Slice-Commits, Stand nach `be98839f`)

Dateiunion der Slice-Commits (je `git show --stat`, union über `9b349879`…`be98839f` plus Vorläufer `8912f7fe`):

| Plan-Zeile (§3) | Diff | Abgleich |
|---|---|---|
| `tools/messung/…/messen.sh`, [`docs/messungen/slice-v1-devcontainer-container-name/`](../messungen/slice-v1-devcontainer-container-name/README.md), `docs/messungen/README.md` | `8912f7fe`, `0f6a219e`, `b5be7aaa` | ✓ deckt (Rohdaten je Modus zwei Läufe) |
| `driven/docker/devcontainer_launcher{,_test,_docker_test}.go`, `export_test.go` | `b5be7aaa`, `ff56fd50` | ✓ deckt |
| `port/driven/devcontainer_launcher.go` | `b5be7aaa`, `ff56fd50` | ✓ deckt |
| `port/driven/devcontainer_launcher_test.go` | `b5be7aaa` | ✓ deckt |
| `application/devcontainer_start.go`, `_test.go` | `b5be7aaa` | ✓ deckt (Test-Datei trägt je Bedingung) |
| `domain/containername{,_test}.go` | `b5be7aaa` | ✓ deckt (neu) |
| `port/driving/devcontainer.go` | `b5be7aaa` | ✓ deckt |
| e2e: `devcontainer_name_docker_test.go` (neu), `devcontainer_start/instance/shell/forge_cli_docker_test.go` | `b5be7aaa`, `b919d4b4`, `0ff780ba`, `be98839f` | ✓ deckt (Zeile nach `94564e5b` erweitert) |
| `driving/cli/devcontainer.go`, `_test.go`, `_shell_test.go`, `_instance_test.go` | `b5be7aaa`, `ff56fd50` | ✓ deckt |
| `spec/lastenheft.md`, `spec/spezifikation.md` | `b5be7aaa` | ✓ deckt |
| `docs/user/benutzerhandbuch.md`, `CHANGELOG.md` | `b5be7aaa` | ✓ deckt |
| diese Datei, `open/README.md`, `roadmap.md` | `9b349879`, `0f6a219e`, `95b5c864`, `ff56fd50`, `94564e5b` | ✓ deckt (nach F-2) |
| **unverändert**-Zeilen: `devcontainer_sandbox.go`, `devcontainer.json.tmpl`, `generate_*_test.go`, `harness/replay.md`, `inspector.go`, `ValidateDevcontainerInstance`, `doctor*.go`, `docs/plan/adr/` | — | ✓ **alle im Diff abwesend** (geprüft gegen die Union; die gelesenen Fundstellen des Plans sind deckungsgleich: `FindContainer` sucht über Labels, `doctor` über den Fallback-Reader) |

**Every diff file with a plan entry: erfüllt** bis auf eine Ausnahme:

| Diff-Datei ohne Plan-Eintrag | Einordnung |
|---|---|
| [`docs/reviews/2026-11-01-slice-v1-devcontainer-container-name.md`](2026-11-01-slice-v1-devcontainer-container-name.md) (angelegt `ff56fd50`, geändert `b919d4b4`) | Review-Artefakt, nicht Produkt-Diff; die Plan-Tabelle §3 des Implementers listet die berührten Produkt-/Doku-/Test-/Mess-Dateien. Hinweis zur Sauberkeit, keine DoD-Verletzung; bei künftigen §3.10-Abgleichen als Ausnahme-Klasse sichtbar machen, falls das Repo sie nicht ohnehin ausnimmt. |

## Sensors

| Sensor | Ergebnis | Evidence |
|---|---|---|
| `make gates` (Docker-Modus, Host) | **pass** | EXIT=0 am 2026-11-02; `baseline-verify: v6.18.1 OK — 54 Dateien`, `d-check: 1007 Datei(en) geprüft, 0 Befund(e)`, lint+test im Docker-Test-Stage grün, `coverage-gate: OK — coverage 92.20% meets threshold 90%`, `ablage-pruefung: OK` (26 Faelle, 0 Fehlschlaege), Stempel `.harness/state/gates-passed.head` = `71cf2f23…` = HEAD |
| `make test-docker` | **not run lokal — belegt durch CI** | Projektinhaber-Entscheidung 2026-11-01: kein lokaler Lauf; Nachweis ist der CI-Lauf **integration-docker 41340432215** (`//go:build docker`, 16m35s, grün, Push-Trigger auf dem Tip) |
| `make ci` | **not run lokal — belegt durch CI** | CI-Lauf **ci 41340432264** auf dem Tip: `gates (lint + test + coverage-gate)`, `security-gates (govulncheck)`, `image-scan (trivy HIGH+CRITICAL)` — alle drei Jobs grün (2026-11-01) |
| `make verify-depguard` | **not run, begründet** | Keine Importregel berührt: der Diff enthält keine Änderung am Layer-Set oder an der depguard-Konfiguration; die neuen Importe folgen bestehenden Regeln (`application`→`domain`, `domain` blattfrei `fmt`, `driven/docker`→`port/driven`, `cli`→`driving`); Lint mit depguard in `gates` grün; Negativbefund „Architekturgrenzen“ des Review-Reports |
| `make modus-pruefung` | **not run, begründet** | Docker-only-Werkzeuge (`Makefile`, `Dockerfile`, `d-check.mk`, `repo.mk`) sind im Diff nicht berührt |
| `make ablage-pruefung` | pass (in `gates` enthalten) | s. o., inklusive Wächter-Test |

## Traceability

| ID | Beleg |
|---|---|
| [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) | e2e `TestDevcontainerUp_Docker_{ContainerName,ConflictingContainerName}` (CI 41340432215), Unit-Pins (Compose, Rejects, Conflicting, RemovesTemporary, HandsNoName), CR im Lastenheft 0.29.0, Messprotokoll M1–M6 |
| [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) | Lastenheft-Fehlerklassen (zwei neuen `10`-Fällen im `up`-Absatz), [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) Exit-Tabelle, Pin-Tabelle `TestExecute_Devcontainer_SentinelsPerSubcommand` (F-1-Fix), e2e-Konflikt mit Exit `10` |
| [`LH-FA-INIT-002`](../../spec/lastenheft.md#lh-fa-init-002--projektname) | Projektname-Quelle: Spec-Block „Containername“ verweist auf die Regel; `TestDevcontainerUp_ComposesContainerNameFromFolder` (Ordner-Basisname, Normalisierung), `domain/projectname.go` Grenzen unverändert |
| [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) | Pin-Tabelle mit Diagnose-Klasse je Sentinel; „`--json` unchanged by the name“ (`devcontainer_test.go` Zeile 574 ff.) |
| [`LH-FA-ARCH-002`](../../spec/lastenheft.md#lh-fa-arch-002--schichten-und-verzeichnislayout)/[`003`](../../spec/lastenheft.md#lh-fa-arch-003--import-regeln-und-enforcement) | `domain` blattfrei (`fmt`), Engine-Regel einmalig im Port `driven`, Overlay-Merge in `application`, I/O im `driven`-Adapter, Ausgabe in `cli`; Lint/depguard grün in `gates`; Negativbefund des Reviewers |
| [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache) | CLI-Texte englisch („Devcontainer started: container …, name …“), Handbuch-Sprachzustand unverändert, kein Identifikator im Handbuch (0 Treffer, geprüft) |
| [`ADR-0038`](../plan/adr/0038-devcontainer-instanz-zweites-id-label.md) | Unberührt (Datei im Diff abwesend); Kopplung gemessen (M3/M6: `devcontainerId`, Volumes, Labels je Zelle identisch) und im e2e gegenprobt (Label bleibt am namentragenden Container) |
| [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) | Drei neue Blöcke + Ausgabetabelle; Bedingungen je an Test oder benannte Messung gebunden |

## Carveouts

- Neu: **keine** — `docs/plan/carveouts/README.md` ist im Diff nicht berührt; keine Lockierung von Gates, Schwellen oder Architekturregeln.
- Gelöst: —
- Unverändert: —

## Rest-Risiken / offene Punkte vor der Closure

1. **Zwei Spezial-Bedingungen der Spec sind mess- statt testgetragen** (DoD 3 formuliert „pinnender Test oder entfällt“; beides ist über Messung belegt, F-6 vom Review als implizit testgetragen eingestuft — der Verifier hält dies für tragfähig, zählt es aber als Rest, weil kein eigener Testpin existiert):
   - „Der haltende Container wird nicht entfernt oder umbenannt, auch nicht mit `--recreate`“ — M5 gemessen; der Abbruchpfad liegt vollständig vor `runUp`, ein `--recreate` erreicht den haltenden Container nicht (der Konflikt-e2e-Test setzt `Recreate` nicht, prüft aber, dass der Holder weiter läuft).
   - „Die CLI löst die relativen Pfade der Konfiguration unabhängig vom Ablageort der Kopie auf“ — M6-Sonde (fremdes Verzeichnis, exit=0); implizit zusätzlich dadurch getragen, dass **jeder** docker-getaggte Name-Test die Kopie in `/tmp` übergibt und die relativen Pfade im Projekt auflöst.
2. **Länge oberhalb 257 Zeichen nicht gemessen** (Protokoll-Grenze): Die Zusage ist in Spec/CR/CHANGELOG bewusst auf das Gemessene begrenzt und durch die Zusammensetzungs-Grenze (≤ 103 Zeichen) geschlossen; der Engine-Ablehnpfad oberhalb 257 ist nicht als Test gepinnt — konsistent zur Grenze, nicht eine Verletzung.
3. **Closure-Offenes (DoD-Zeile 4):** Closure-Notiz §7 mit Lerneintrag, Register-Fortschreibung, Risiko-Ausgänge und drei Paarungen (Anker · Folge-Slice · Register) — der reguläre Closure-Schritt; Belegset dazu liefert dieser Report.
4. INFO-Findings F-5/F-7 bleiben dokumentiert (Review-Report); kein Handlungsbedarf für diese Closure.

## Gesamturteil

**Closure-fähig.** Alle drei Liefer-DoD-Zeilen (Messung/Freigabe, Name inkl. Tests je Bedingung und Mutation, Doku/Zusage-Menge) sind erbracht; `make gates` ist am aktuellen Stand (`be98839f`-Stand, Stempel = Tip `71cf2f23`) grün, die docker-getaggten Zusagen sind durch den grünen CI-Lauf 41340432215 auf dem Tip belegt, die Review-Findings F-1..F-4 sind im Code/Plan nachgeprüft behoben. Kein Befund, der die Closure von [`slice-v1-devcontainer-container-name`](../plan/planning/done/slice-v1-devcontainer-container-name.md) blockiert; die Closure-Notiz mit den Risiko-Ausgängen, der Register-Fortschreibung und den drei Paarungen ist der noch offene Schritt, der diesen Report als Beleg liest.

## Commit / Artefakt

- Slice: `8912f7fe`, `9b349879`, `0f6a219e`, `95b5c864`, `b5be7aaa`, `ff56fd50`, `b919d4b4`, `94564e5b`, `0ff780ba`, `be98839f`
- HEAD/Stempel: `71cf2f2311b5244eb4bc92fb2167e9d3488a311f`
- CI: `ci`-Lauf `41340432264` (grün), `integration-docker`-Lauf `41340432215` (grün), jeweils 2026-11-01 auf dem Push-Tip
- Lokal: `make gates` EXIT=0 am 2026-11-02, `.harness/state/gates-passed.*`