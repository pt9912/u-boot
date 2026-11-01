# Review-Report: slice-v1-release-cut-v0.17.0 — 2026-11-01

**Review-Art:** Code (Release-Cut: Diff gegen Plan, Hard Rules und `releasing.md`; Lead gegen Lastenheft und Handbuch)

**Gegenstand:** [`slice-v1-release-cut-v0.17.0`](../plan/planning/done/slice-v1-release-cut-v0.17.0.md), Diff `f58dc581..v0.17.0` (= `HEAD` `f89424e1`, Baum sauber): Cut `c5ed86dc`, Slice-Anlage `f54b3345`, vier Move-Commits des Lifecycle-Werkzeugs, Sensoren `f89424e1`.

**Skill:** `.harness/skills/reviewer.md` @ `f89424e1`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-11-01

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan [`slice-v1-release-cut-v0.17.0`](../plan/planning/done/slice-v1-release-cut-v0.17.0.md) und das Muster [`slice-v1-release-cut-v0.16.1`](../plan/planning/done/slice-v1-release-cut-v0.16.1.md)
- [`docs/maintainer/releasing.md`](../maintainer/releasing.md) §3 bis §7, `harness/review.md`
- Lastenheft 0.26.0, [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Zeile 504), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)
- [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md) (nur Versionsaussage), [`ADR-0038`](../plan/adr/0038-devcontainer-instanz-zweites-id-label.md), [`ADR-0039`](../plan/adr/0039-devcontainer-shell-dritter-aufruf-bei-fehlendem-benutzer.md) (Deprecated, nicht Inhalt des Release)
- `AGENTS.md` (§3.7, §3.9, §3.10, §5 Regel 6), Benutzerhandbuch, `.github/workflows/`

Nicht ausgeführt: `make test-docker` (Auftrag), `make verify-depguard` (keine Importregel berührt), jede Aktion mit Außenwirkung.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der Handbuch-Abschnitt "Mehrere Instanzen" (Unterabschnitt unter dem Sandbox-Profil) sagt "Parallele Instanzen desselben Projekts starten Sie aus **getrennten Ordnern**", "Derselbe Ordner bleibt eine Instanz" und "Das Dev-Containers-Werkzeug löst `devcontainerId` pro Projektordner auf". Der Abschnitt "Mehrere Devcontainer je Projekt (Instanzen)" derselben Handbuch-Fassung beschreibt `--instance`, mehrere Container je Ordner mit je eigenen Volumes. Der Release liefert beide Abschnitte aus; der Cut (Handbuch nur Kopf und Versionsstellen) hat sie nicht gegeneinander gelesen. | [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), `AGENTS.md` §5 Regel 5, §3.10 | `docs/user/benutzerhandbuch.md` · "Derselbe Ordner bleibt eine Instanz." (Zeile 2288) gegen "Mit `--instance <Name>` starten Sie neben dem" (Zeile 887) | nein — kein Gate liest Handbuchaussagen gegeneinander; `make docs-check` bestätigt es nicht | Handbuch widerspricht sich zu Instanzen (veraltete Aussage neben neuem Abschnitt) |
| F-2 | LOW | Die Liste "Not promised" im Lead ist kürzer als die Liste im Lastenheft: es fehlen "mit kaltem Cache oder bei kaltem Bau des Images" und "Dauer und Platzbedarf des Starts". "Promised" und die übrigen Punkte stimmen wörtlich überein (zwei Instanzen, Workspace-Volumes, Sandbox-Profil, direkt nach dem Start, Linux `amd64`, lokale Docker-Engine). Die Liste ist nicht als Auswahl gekennzeichnet. | [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), `AGENTS.md` §3.10, Beobachtung `BEO-DOK/zusage-weiter-als-messung` | `CHANGELOG.md` · "the behaviour under load, other hosts and engines, and the isolation in the standard profile" | nein — Review-only | Lead-Grenzliste kürzer als Vertragsliste |
| F-3 | LOW | Der Roadmap-Zeiger trägt "(Push und Tag stehen aus)". Der Tag `v0.17.0` liegt auf `HEAD`, `publish` ist grün; der Marker ist abgelaufen. Der Abbau steht als offener Punkt in der Checkliste "Nacharbeit nach dem Tag" (Zeile "Roadmap: Zeile `v0.17.0` nach *Abgeschlossene Wellen*"). | `AGENTS.md` §3.10 (Zustandsmarker mit Ablaufpunkt) | `docs/plan/planning/in-progress/roadmap.md` · "(Push und Tag stehen aus)" | nein — Review-only | Zustandsmarker nach Ablauf nicht abgebaut |
| F-4 | LOW | Der Entwurf [`slice-v1-devcontainer-ressourcenlimits`](../plan/planning/in-progress/slice-v1-devcontainer-ressourcenlimits.md) trägt "Vorschlag, nicht entschieden: Release `v0.17.0`"; dieselbe Angabe steht in seiner Zeile in `open/README.md` ("ohne Welle (Vorschlag `v0.17.0`)"). `v0.17.0` ist veröffentlicht und enthält den Slice nicht (Slice §1 nennt das als Nicht-Ziel, ändert aber den Marker nicht). | `AGENTS.md` §3.10 (Zustandsmarker mit Ablaufpunkt) | Entwurf-Slice, Feld `Welle:` · "Vorschlag, nicht entschieden: Release `v0.17.0`"; `docs/plan/planning/open/README.md` · "ohne Welle (Vorschlag `v0.17.0`)" | nein — Review-Only | Zustandsmarker nach Ablauf nicht abgebaut |
| F-5 | LOW | Der Entwurf [`slice-v1-konfigdatei-alternativer-name`](../plan/planning/open/slice-v1-konfigdatei-alternativer-name.md) trägt "Vorschlag Release `v0.17.0` (Roadmap *Nächste Wellen*)". Die Zeile `v0.17.0` in *Nächste Wellen* hat dieser Diff aus der Roadmap entfernt; der Zeiger zeigt auf nichts, und der Slice ist weder im Release noch im Slice §1 genannt. | `AGENTS.md` §3.10, Pflichtfeld Plan-Tabelle (jede berührte Aussage) | Entwurf-Slice, Feld `Welle:` · "Vorschlag Release `v0.17.0` (Roadmap *Nächste Wellen*)" | nein — Review-only | Zeiger auf vom Diff entfernte Roadmap-Zeile |
| F-6 | LOW | Slice §1 ("Weitere Slices im Release, ohne eigenen CHANGELOG-Eintrag") und der Roadmap-Zeiger ("Inhalt") nennen sieben Slices; alle sieben liegen in `done/`. Zwischen `v0.16.1` und dem Cut sind zwei weitere Slices mit Code im Diff geschlossen und nicht genannt: [`slice-v1-devcontainer-bau-lokal`](../plan/planning/done/slice-v1-devcontainer-bau-lokal.md) (`scripts/bau.sh`, `.devcontainer/`, `repo.mk`, `Dockerfile`; der Slice nennt es nur als "Build-Werkzeug" ohne Kennung) und [`slice-v1-konfigdatei-name-als-konstante`](../plan/planning/done/slice-v1-konfigdatei-name-als-konstante.md) (`domain.ConfigFileName`, nur in Slice §6 als Konstante erwähnt). | `AGENTS.md` §3.10 | `docs/plan/planning/in-progress/roadmap.md` · Zeiger "Inhalt:" unter *Offene Wellen* | nein — Review-only | Inhaltsliste des Release unvollständig |
| F-7 | INFO | Auffälligkeit des Implementers "Binary geändert": belegt gelesen. Im Produktcode `v0.16.1..v0.17.0` stehen `--instance` samt Namensprüfung und Inspector, die Referenz `CLAUDE_CODE_OAUTH_TOKEN` samt `doctor`-Regel und Hinweistext, `domain.ConfigFileName` (Wert `u-boot.yaml` unverändert, `config_filename_literal_test.go` pinnt die Literale) und eine Kommentar-Änderung in `devcontainer_shell.go`. Jede Nutzer-Wirkung steht im CHANGELOG; kein weiterer Effekt gefunden. Build-Werkzeug (`BUILD_MODE`, `scripts/bau.sh`, `.devcontainer/`, Dockerfile-Stages `compile`/`lint`/`test`/`coverage`) und Messskripte unter `tools/messung/` berühren die Stufen `build`/`runtime` nicht; ein CHANGELOG-Eintrag fehlt vertretbar, weil das Produkt davon unberührt ist (Image-Scan 0, `--version`-Gegenprobe im Slice). Der Lead sagt zu Recht nichts über einen unveränderten Stand. | Maintainability | `CHANGELOG.md` · "### Changed" | ja — `make gates` (Slice-Sensoren), CI `ci` | — |
| F-8 | INFO | Handbuch-Version `2.1` belassen: begründet. `v0.16.1` trug `2.0`; `2.1` stammt aus der Umstellung auf den Ist-Zustand (Commit `edbd2089`, 2026-10-31) und liegt damit in diesem Release; `releasing.md` §3 und §5 verlangen keine weitere Anhebung. `Stand` (2026-11-01) und `Software-Version` stimmen mit CHANGELOG und Tag-Datum überein. | `docs/maintainer/releasing.md` §3, §5 | `docs/user/benutzerhandbuch.md` · "Handbuch-Version: 2.1" | nein — Review-only | — |
| F-9 | INFO | CI-Nachweis auf dem Push-Tip `f89424e1` gegengelesen (REST-API, Job-Log): `integration-docker` 41290870231 `success`, im Log 3471 `--- PASS`, 4 `--- SKIP` (alle nested Podman), 0 `--- FAIL`, Paket `internal/e2e` `ok` nach 930,9 s; `ci` 41290870228, `tap-check` 41290880081 und `publish` 41292635297 (Jobs `publish`, `tap`, beide Paket-Smokes) `success`. Die Anzeige in `gh run list` ("0s", `completedAt` vor `startedAt`) ist ein Zeitstempel-Artefakt der Läufe; die Log-Zeiten (08:55:39 bis 09:12:07) tragen eine Dauer von etwa 17 Minuten. Zuständigkeit für die Verifikation der Läufe: Verifier. | `AGENTS.md` §3.1 | Slice §3 *Sensoren* · "`make test-docker` | **nicht gelaufen**" | ja — `gh run view`, Log-Zählung | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `CHANGELOG.md` Kopf `[0.17.0]`, Datum, "Seventeenth release." (Zählung der Einträge `0.1.0` bis `0.16.0`: fortlaufend, `0.16.0` "Sixteenth") | geprüft, ohne Befund |
| `CHANGELOG.md` Lead gegen Diff `v0.16.1..v0.17.0`: `--instance`, Exit-Code `10`, Referenz `CLAUDE_CODE_OAUTH_TOKEN`, "never a value", Neu-Erzeugung per `generate devcontainer`, "sets up no login", "keeps no list of instances" (Handbuch) | geprüft, ohne Befund |
| `CHANGELOG.md` "Promised" gegen Lastenheft Zeile 504 | geprüft, ohne Befund (Grenzliste siehe F-2) |
| `CHANGELOG.md` Fuß: `[Unreleased]: …v0.17.0...HEAD`, `[0.17.0]: …v0.16.1...v0.17.0`, `[0.16.1]` unverändert | geprüft, ohne Befund |
| Versionsstrings `cmd/uboot/main.go:43`, `repo.mk:61`, `Dockerfile:45` samt Kommentaren (`repo.mk:57,60`, `Dockerfile:41,43`) | geprüft, ohne Befund; kein weiterer Treffer von `0.17.0` außerhalb von `docs/`, READMEs und der Lastenheft-Historie |
| `0.16.1-dev` im Baum (`grep -rI`, ohne `.git`) | geprüft, ohne Befund: nur Planung (`done/`, Slice), `docs/reviews/`, `docs/messungen/` (AGENTS.md §3.9 b) |
| `0.16.1` außerhalb dieser Bereiche | geprüft, ohne Befund: [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md) (Accepted, Aussage über jenen Stand), Entwurf `ressourcenlimits` (Messstand), CHANGELOG-Abschnitt, Beobachtungs-Evidence |
| `README.md`, `README.de.md` | geprüft, ohne Befund: nur Status-Version und `VERSION=`-Beispiel; keine Kennung, keine Chronik, kein "released" (`AGENTS.md` §5 Regel 6) |
| Handbuch Kopf, Versionsbeispiele, Gültigkeitsbereich | geprüft, ohne Befund (Widerspruch im Inhalt: F-1) |
| Roadmap-Zeile *Nächste Wellen* `v0.17.0` entfernt, Inhalt-Slices (sieben) in `done/` | geprüft, ohne Befund (Lücke der Liste: F-6, Marker: F-3) |
| Slice §3 Plan-Tabelle gegen `git diff --stat f58dc581..v0.17.0` (10 Dateien: `CHANGELOG.md`, `Dockerfile`, beide READMEs, `main.go`, `roadmap.md`, Slice, `open/README.md`, Handbuch, `repo.mk`) | geprüft, ohne Befund: jede Zeile hat Plan-Eintrag und Diff; "kein Diff" für `releasing.md`, Workflows, `scripts/`, `packaging/`, `Makefile`, `go.mod`, `spec/`, `docs/plan/adr/` bestätigt |
| Slice §3 Zahlen "70 Dateien (3438 Einfügungen, 213 Löschungen)" | geprüft, ohne Befund (`git diff --stat v0.16.1..v0.17.0`, Produkt-Ausschnitt) |
| Slice *Sensoren*: `make test-docker` und `tap-check` als "nicht gelaufen" ausgewiesen, Trockenlauf mit benannter Grenze | geprüft, ohne Befund (Ergebnisse der übrigen Sensor-Zeilen nicht erneut ausgeführt: Docker-Läufe `make gates`/`ci`/`fullbuild` liegen beim Verifier) |
| Slice DoD und Nacharbeit §7: Kästchen offen, §7 "Folgt nach dem Tag" | geprüft, ohne Befund: der Zustand "Tag gesetzt, Nacharbeit offen" ist ehrlich ausgewiesen; die Checkliste (Roadmap, Datum, Job `tap`, Zeiger `releasing.md` "zuletzt" auf `v0.16.1`, Closure, `git mv` nach `done/`) ist vollständig und offen |
| `git mv`-Disziplin (§3.3): vier Move-Commits sind reine Moves mit getrenntem Verweis-Nachzug | geprüft, ohne Befund |
| Reviewer-HIGH-Liste (Dual-Classifier, Referenzmodell, Sprachvertrag, Architekturgrenze, Gate-Lockerung, Immutabilität, Carveout) | geprüft, ohne Befund: der Cut berührt keinen Produktcode, kein Gate, kein ADR, keine Spec |
| Produktcode `v0.16.1..v0.17.0` (Dual-Classifier der Instanz-Sentinels, Tests) | nicht geprüft: Gegenstand der Inhalts-Slices, die je ein Review haben; der Cut ändert nur drei Versionszeilen |
| `make test-docker`, `make verify-depguard` | nicht ausgeführt: Auftrag bzw. keine Importregel berührt |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 5 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Handbuch widerspricht sich zu Instanzen · Lead-Grenzliste kürzer als Vertragsliste · Zustandsmarker nach Ablauf nicht abgebaut (F-3, F-4) · Zeiger auf vom Diff entfernte Roadmap-Zeile · Inhaltsliste des Release unvollständig

## Verdikt

**Merge-blockierend:** nein — der Tag ist gesetzt und veröffentlicht; er ist nicht zurückzunehmen ([`releasing.md`](../maintainer/releasing.md) *Zweck*), und kein Finding betrifft Binary, Versionsstrings oder Compare-Links. F-1 (MEDIUM) betrifft das Handbuch im Repo und lässt sich nach dem Tag per Folgecommit beheben; der Slice geht nicht nach `done/`, bevor der Hauptagent F-1 entschieden hat. F-3 und F-4 fallen in die offene Nacharbeit. Die Findings werden nicht herabgestuft, wenn der Implementer widerspricht; HIGH gibt es nicht, ein Konfliktpfad über den Architect entfällt.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die Closure-Notiz §7 des Slice und von dort in den Zähler. Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation (Verifier, anderer Kontext). Transiente Befunde: keine, die `docs-check` betreffen; der parallele Verifier-Report dieses Slice lag zum Zeitpunkt des Laufs ungetrackt im Baum und wurde nicht gelesen.

## Nachlauf Fix-Runde

Enger Nachlauf gegen `git diff f89424e1..HEAD` (Handbuch *Mehrere Instanzen*, CHANGELOG-Lead, `tools/mk/devcontainer.mk`, `harness/README.md`). Maßstab: Plan und Entscheidung [`ADR-0038`](../plan/adr/0038-devcontainer-instanz-zweites-id-label.md), [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`AGENTS.md`](../../AGENTS.md) §3.10.

### Findings

- `kategorie`: LOW
  `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.10, [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)
  `pfad`: `docs/user/benutzerhandbuch.md` · „Jede Instanz hat damit eigene Volumes, der Container"
  `klasse`: Handbuch-Aussage stärker als Zusage und Test
  `befund`: Der Satz nennt alle Volumes der Namensliste davor (Workspace, Podman, Docker, `ollama`) je Instanz eigen. Gemessen und im Docker-Test `TestDevcontainerInstance_Docker_TwoInstancesAreIsolated` gepinnt sind nur die Workspace-Volumes (`dcRequireDistinct(t, "workspace volume", …)`); die Zusage in der Spec nimmt Docker-Speicher und Modelle ausdrücklich aus, und der Abschnitt *Mehrere Devcontainer je Projekt (Instanzen)* formuliert dort vorsichtiger („Volumes, deren Name von der Kennung des Devcontainers abhängt“).
  `risiko`: Der Absatz liest sich für Docker-Speicher und Modelle wie eine Zusage; ein Änderer des Handbuchs oder ein Nutzer mit `nestedRuntime: docker` verlässt sich auf eine Isolation, die kein Test belegt. Der Satz ist als Ist-Zustand aus [`ADR-0038`](../plan/adr/0038-devcontainer-instanz-zweites-id-label.md) (Volume-Namen tragen `${devcontainerId}`, die Id wechselt mit dem Id-Label) herleitbar, daher kein MEDIUM.
  `verifizierbar`: Review-only (Gegenlesen mit `internal/e2e/devcontainer_instance_docker_test.go`)

- `kategorie`: INFO
  `quelle`: Maintainability
  `pfad`: `tools/mk/devcontainer.mk` · `$(if $(INSTANCE),--instance $(INSTANCE))`
  `klasse`: Make-Variable unquotiert ins Shell-Rezept
  `befund`: `INSTANCE` wird ungequotet eingesetzt; ein Wert mit Shell-Metazeichen wirkt als Shell-Syntax. Der Name wird erst von `u-boot` geprüft (Exit-Code `10`), das Rezept läuft vorher durch die Shell. Der Aufrufer ist der Nutzer selbst, keine Vertragsverletzung.
  `risiko`: gering; nur lokaler Komfort-Pfad, kein Gate.
  `verifizierbar`: Review-only

### Geprüfte Fragen

1. Handbuch-Absatz. Fundstellen gelesen: `runArgs` in `internal/adapter/driven/docker/devcontainer_launcher.go` (mit Instanz zwei `--id-label`, ohne keines; Test in `devcontainer_launcher_test.go` pinnt beide Listen), `volumeMount` in `internal/hexagon/application/devcontainer_sandbox.go` (`<volume>-${devcontainerId}`), `generate_ollama_test.go` (`dockerMount`, `ollamaMount`, `workspaceMountLine`), Docker-Test (Container ohne Instanz trägt kein Instanz-Label, Instanzen eigene Workspace-Volumes und Container). „Container ohne Instanz behält seine“ ist durch den Argumente-Test und [`ADR-0038`](../plan/adr/0038-devcontainer-instanz-zweites-id-label.md) Entscheidung 3 getragen. „Das Werkzeug bildet `devcontainerId` aus dem Projektordner“ und „fließt zusätzlich der Name ein“ ist das Verhalten der externen CLI; Beleg ist die Messung der ADR (Punkt 1), kein Test im Repo außer dem Docker-Test über das Ergebnis. Keine Zusage-/Mess-Sprache im Absatz. Widerspruch zum Abschnitt *Mehrere Devcontainer je Projekt (Instanzen)* (Zeile 887 ff.): kein Widerspruch in der Sache, nur der Stärkeunterschied des Findings oben. Widerspruch zu [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) und Lastenheft: keiner über das Finding hinaus.
2. Reste im Baum zu „getrennte Ordner“, „Derselbe Ordner bleibt eine Instanz“, „pro Projektordner“, „je Projektordner“ in `docs/user`, `README*`, `docs/maintainer`, `CHANGELOG.md`: keine Treffer. Die Treffer zu „Worktree“ (Zeile 1339) betreffen die Git-Remote-Erkennung und gehören nicht zum Thema.
3. CHANGELOG-Lead *Not promised*: deckungsgleich mit der Nicht-Zusage-Liste in [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (drei und mehr Instanzen, Docker-Speicher, Modelle, Podman-Speicher, CPU-/Speicherbegrenzung, Last, kalter Cache oder kalter Bau, Dauer und Platzbedarf des Starts, andere Hosts und Engines, Standardprofil); die Lead-Kurzform lässt die Beispiele (Go-Build, Modell-Pull) weg, ohne Inhalt zu ändern. Ohne Befund.
4. `tools/mk/devcontainer.mk` und `harness/README.md`: `make -n dev-up INSTANCE=slice-a RECREATE=1` ergibt `u-boot devcontainer up --recreate --instance slice-a`; ohne `INSTANCE` bleibt das Rezept das alte (nur ein Leerzeichen am Ende). `make help` und die drei Tabellenzeilen nennen `INSTANCE=<name>` einheitlich. Ohne Befund.

### Nicht geprüft

`make test-docker` (laut Auftrag nicht ausgeführt; das Verhalten der externen CLI zur `devcontainerId` ist hier nur über Messprotokoll und ADR belegt). Verifikation der DoD: Verifier.

### Summary Nachlauf

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 1 |

**Verdikt Nachlauf:** nicht blockierend; das LOW geht an den Implementer (Klasse: Handbuch-Aussage stärker als Zusage und Test). Transiente Befunde: keine.
