# Review-Report: adr-0038-devcontainer-instanz — 2026-10-31

> Der Report entstand per `cp` aus der Vorlage (`diff -q` identisch) und wurde mit `Write` gefüllt; der Reviewer-Agent hat kein Edit-Werkzeug.

**Review-Art:** Design (ADR im Status `Proposed` gegen Spec, Entscheidungslage, Messprotokoll und Code-Ist-Zustand)

**Gegenstand:** uncommittete Änderung auf `78c5404a`: [`ADR-0038`](../plan/adr/0038-devcontainer-instanz-zweites-id-label.md) (neu), `docs/plan/adr/README.md` (Indexzeile), Plan-Tabelle von [`slice-v1-devcontainer-instanz-label`](../plan/planning/next/slice-v1-devcontainer-instanz-label.md) (`git diff HEAD`)

**Skill:** `.harness/skills/reviewer.md` @ `78c5404a`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-31

**Eingangs-Kontext:**

- [`slice-v1-devcontainer-instanz-label`](../plan/planning/next/slice-v1-devcontainer-instanz-label.md) (Plan-Tabelle, §3)
- [`ADR-0015`](../plan/adr/0015-sandbox-volumes-pro-instanz.md), [`ADR-0028`](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md), [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md), [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md) (Formvergleich), [`ADR-0013`](../plan/adr/0013-dokumentationsreferenzmodell.md)
- [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Lastenheft 0.25.0), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)
- Messprotokoll [`README.md`](../messungen/slice-v1-devcontainer-instanz-messung/README.md) samt Rohdaten, Closure-Notiz [`slice-v1-devcontainer-instanz-messung`](../plan/planning/done/slice-v1-devcontainer-instanz-messung.md) §7
- Code: `internal/adapter/driven/docker/inspector.go`, `devcontainer_launcher.go`, `internal/hexagon/application/devcontainer_start.go`, `devcontainer_fallback.go`, `internal/adapter/driving/cli/devcontainer.go`
- `AGENTS.md` (§3.4, §3.5, §3.6, §3.7, §3.10), `harness/conventions/MR-008-adr-form.md`, ADR-Vorlage der Baseline `v6.18.1`

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Die ADR schärft nur den Satz der Variante D („zwei Instanzen aus **demselben** Ordner bleiben eine Instanz“) und erklärt, sie löse keine Entscheidung von [ADR-0015](../plan/adr/0015-sandbox-volumes-pro-instanz.md) ab. Der tragende Satz von [ADR-0015](../plan/adr/0015-sandbox-volumes-pro-instanz.md) steht in dessen Entscheidung: „Eine Instanz ist ein Ordner … Parallele Instanzen desselben Projekts laufen aus getrennten Ordnern“; das gilt mit `--instance` nicht mehr. Die Indexzeile zitiert „eine Instanz ist ein Ordner“, die ADR selbst nennt den anderen Satz. | `AGENTS.md` §3.5, [`ADR-0015`](../plan/adr/0015-sandbox-volumes-pro-instanz.md) | `docs/plan/adr/0038-devcontainer-instanz-zweites-id-label.md` · „Er gilt weiter für den Container ohne Instanz“; `docs/plan/adr/README.md` · „den Satz „eine Instanz ist ein Ordner“ im Indexeintrag | nein — Review-only | Schärft-Reichweite enger als die Wirkung, Zitat weicht zwischen ADR und Index ab |
| F-2 | LOW | Die ADR führt als „nicht gemessen“: `up` ohne Flag bei Container ohne Instanz und neuerer Instanz. Die Rohdaten zeigen es für die CLI im Standardprofil: nach Bestandscontainer und Instanz `i1` (läuft) trifft „Erneut ohne Flag“ den Bestandscontainer (`rohdaten-stufe2-lauf1.txt`, Abschnitt `[T2-standard] M6 umgekehrt`). Das Protokoll-README sagt es in M6 umgekehrt, widerspricht sich aber mit dem ersten M6-Absatz („Standardprofil neben Instanzen nicht gemessen“); die ADR folgt dem README. | [`ADR-0028`](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) (Messgrundlage), `AGENTS.md` §3.10 | `docs/plan/adr/0038-devcontainer-instanz-zweites-id-label.md` · „gemessen ist nur der umgekehrte Weg der Anlage“; Messprotokoll `README.md` (M6 umgekehrt) · „ein erneutes `up` ohne Flag trifft weiter den Bestandscontainer“ | ja — Rohdaten lesen; Docker-Test des Code-Slice | Messaussage der ADR liegt unter dem Rohdatenstand; Protokoll-README in sich uneinheitlich (Richtung sicher, nichts überbehauptet) |
| F-3 | LOW | Vier Aussagen stehen ohne Kennzeichnung als weder gelesen noch gemessen: (a) Entscheidung 4: dass `docker ps --format` das Label je Treffer ausgibt (der Filter auf fehlendes Label wird als ungemessen markiert, das Ausweichmittel nicht); (b) Entscheidung 2: „Ihr Hash bestimmt die `devcontainerId`“ (gemessen ist nur, dass sie mit dem zweiten Label wechselt, M2); (c) Alternative D: Compose-Projektname wirke „auf Dienste, nicht auf den Devcontainer und seine Volumes“; (d) Alternative E: Bestandscontainer würden „beim nächsten `up --recreate`“ verwaist (aus M2 abgeleitet; ein `up` ohne `--recreate` fände den Bestand dann schon nicht mehr). | `AGENTS.md` §3.10, `AGENTS.md` §3.7 | `docs/plan/adr/0038-devcontainer-instanz-zweites-id-label.md` · „liest darum das Label je Treffer aus der Ausgabe (`docker ps --format`)“; · „Ihr Hash bestimmt die `devcontainerId`“; · „nicht auf den Devcontainer und seine Volumes“; · „beim nächsten `up --recreate` verwaist“ | nein — Review-only; (a) und (d) fängt der Docker-Test des Code-Slice | Ableitung als Befund geführt, Mess-/Lesegrenze nicht an der Aussage |
| F-4 | LOW | „Weg A“, „Weg B“, „Weg C“ und „Fassung (i)“ sind in der ADR nicht definiert; sie stammen aus dem Slice, den eine ADR nicht als Quelle verlinken darf, und aus einem „Gespräch am 2026-10-31“. Die Entscheidung des Projektinhabers hat in der ADR keinen auflösbaren Beleg ([ADR-0031](../plan/adr/0031-devcontainer-shell-ohne-editor.md) führt dafür einen Abschnitt „Entscheidungen des Projektinhabers“). | [`ADR-0013`](../plan/adr/0013-dokumentationsreferenzmodell.md), `harness/conventions/MR-008-adr-form.md` | `docs/plan/adr/0038-devcontainer-instanz-zweites-id-label.md` · „Vom Projektinhaber am 2026-10-31 entschieden: **Weg A**“ | nein — Review-only | Entscheidungsherkunft in der ADR nicht auflösbar |
| F-5 | LOW | Die Fitness-Zeile „depguard: die Namensregel liegt im Port-Paket `driven`“ nennt `make lint` als Sensor. depguard prüft Importrichtungen (Adapter importiert `domain` nicht), nicht den Ort einer Funktion; eine Namensregel im falschen Paket fiele nicht auf. | `AGENTS.md` §3.10 (Bedingung ohne Test), [`ADR-0003`](../plan/adr/0003-solid-nahes-lint-profil.md) | `docs/plan/adr/0038-devcontainer-instanz-zweites-id-label.md` · „die Namensregel liegt im Port-Paket `driven`; Application und Adapter importieren sie von dort“ | nein — `make lint` bestätigt nur den Importteil | Sensor-Zusage weiter als der Sensor |
| F-6 | LOW | Die neue Plan-Zeile „Test-Doubles und Tests der Signatur“ nennt `devcontainer_shell_test.go` und `export_test.go` (Paket `application`). Beide enthalten kein `FindContainer`; die Berührung ist nicht belegt (die Zeile benennt als Fundstelle nur die Signatur). | `AGENTS.md` §3.10 | Plan-Tabelle des Slice Slice-Plan (Bezug oben) · „`devcontainer_start_test.go`, `devcontainer_shell_test.go`, `export_test.go`“ | ja — `grep FindContainer` in beiden Dateien: 0 Treffer | Plan-Eintrag ohne belegte Berührung |
| F-7 | INFO | Vom Projektinhaber zu bestätigen mit `Accepted`: der Labelname `u-boot.instance` (steht nur in Entscheidung 2, nicht als eigener Punkt; ein Wechsel lässt bestehende Instanzen verwaisen; Docker-Konvention wäre Reverse-DNS, `u-boot.instance` ist ein Kurzname ohne Domäne), die Namensregel (Kleinbuchstaben, Ziffern, `-`, 1 bis 32; strenger als gemessen, M8 nahm 23 von 25 Werten, keine Grenze unter 128), die Mechanik (Entscheidungen 3 bis 9), und ob die offene Messung (Container ohne Instanz plus neuere Instanz, Sandbox-Profil und `down`/`shell`) vor oder nach `Accepted` läuft. | `AGENTS.md` §3.6 | `docs/plan/adr/0038-devcontainer-instanz-zweites-id-label.md` · „mit `Accepted` ist er bestätigt“ | nein — Entscheidung des Projektinhabers | Entscheidungsliste |
| F-8 | INFO | Nicht gemessen und in der ADR nicht genannt: wiederholtes `up --instance <name>` ohne `--recreate` (M7 misst nur `--remove-existing-container`); `down --instance <unbekannt>` (Verhalten „nichts zu stoppen“, Exit-Code `0`, nicht ausgesprochen); Instanz-Container tragen kein `devcontainer.config_file` (M1), Werkzeuge, die darüber suchen, finden sie nicht. Zuständig: Docker-Test und Spec-Slice. | Maintainability | Messprotokoll `README.md` (M1) · „`devcontainer.config_file` fehlt“ | nein — `make test-docker` des Code-Slice | Lücke im Konsequenzen-Abschnitt (nicht überbehauptet) |
| F-9 | INFO | Der Link auf das Messprotokoll trägt eine Slice-Kennung im Pfad des Verzeichnisses. Das ist der Pfad einer Messung (Rang Messprotokoll, nicht Slice), [ADR-0028](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) und [ADR-0036](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) verlinken ebenso; `make docs-check` meldet 0 Befunde. Kein Verstoß gegen die Decken-Regel. | [`ADR-0013`](../plan/adr/0013-dokumentationsreferenzmodell.md) | `docs/plan/adr/0038-devcontainer-instanz-zweites-id-label.md` · „Das Messprotokoll“ (Link auf `README.md` des Messverzeichnisses) | ja — `make docs-check` | Hinweis, kein Befund |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Form (MADR, MR-008): Kopffelder, Abschnittsfolge, Titel `# ADR 0038:`, Status `Proposed`, Geschichte, Indexzeile, Anker-Links ([`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes), `LH-*`, ADR-Links) | geprüft, ohne Befund; `make docs-check`: 0 Befunde |
| Alternativen: sieben Optionen (nichts tun, Worktrees, Ordnerpfad, Compose, Label auch für den Standard, neuester Container, gewählte Variante) mit Pro/Contra; Konsequenzen nennen Race, Unsichtbarkeit für `doctor`/`generate`/`init`, Label-Wechsel, offene Messung | geprüft, ohne Befund (Einzelaussagen siehe F-3) |
| Fundstellen im Code: `runArgs` ohne `--id-label` (`devcontainer_launcher.go`), `FindContainer` mit `ps -a --filter label=devcontainer.local_folder=…` und erstem Treffer (`inspector.go`), die drei Aufrufer (`runningContainerID`, `Down`, `fallbackReader.read`), `startOutcome` über `fallbackReader`, `doctor`/`generate`/`init` über `fallbackReader`, `mapDevcontainerDownErrorToDiagnostic` ohne Klasse-`10`-Sentinel, `isDevcontainerValidationError` von `isValidationError` eingebunden, `ValidateDevcontainerProjectFolder`, `DevcontainerUpOptions`, die drei Request-Typen | geprüft, ohne Befund, alle Symbole existieren und tun, was die ADR sagt |
| Messaussagen (1) bis (5) des Kontexts gegen das Protokoll (M1, M2, M6, M6 umgekehrt, M8) | geprüft, ohne Befund; ADR verlinkt statt Messwerte zu wiederholen (Einschränkung F-2) |
| Decken-Regel: kein Slice, keine Welle, kein Carveout als Bezug oder Trigger in der ADR; die Rückschau der ADR geht nur auf Lastenheft, Spec und ADRs | geprüft, ohne Befund |
| `AGENTS.md` §3.5: die Accepted-ADRs 0015, 0028, 0031 bleiben unverändert (`git status`: nur die neue ADR, Index, Slice berührt); „Schärft“ statt „Supersedes“ ist für den Satz der Variante D wahr (Reichweite siehe F-1) | geprüft, ohne Befund bis auf F-1 |
| Übereinstimmung mit Lastenheft 0.25.0 ([`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)): Flag, Namensregel, Exit-Code `10` für Namen und „nur Instanzen“, Verhalten ohne Flag für `up`, `down`, `shell`, Zusage-Menge als Verweis | geprüft, ohne Befund |
| Dual-Classifier-Abdeckung der Entscheidung 8 (zentral plus Diagnostic von `up` und `down`, `shell` nur Klasse `2`) | geprüft, ohne Befund; im Code noch nicht vorhanden (Plan-Stand), Pin-Tests sind verlangt |
| Mess-/Prozesssprache im Entscheidungstext (Indikativ, keine Chronik der Verwerfung) | geprüft, ohne Befund |
| Plan-Tabelle des Slice: Dateien und Symbole der neuen Zeilen (`driving/devcontainer.go` Zeilen 11/56/89/114, `driven/devcontainer_launcher.go` 28, `noop.go` 51, `cli.go` 386, `devcontainer.go` 548, `fakeInspector.FindContainer`, `inspector_test.go`, E2E- und Launcher-Tests, `cli-json-contract.md` §6.11, `zusagen-und-messstaende.md`, `harness/replay.md` „nennt `up`/`down`/`shell` nicht“) | geprüft, ohne Befund bis auf F-6 |
| Code im Repo (Implementierung) | nicht geprüft — der Diff enthält keinen Code; der Code-Slice startet erst mit `Accepted` |
| `make test-docker`, `make verify-depguard` | nicht gelaufen — kein Docker-Pfad und keine Importregel berührt; keine Container gestartet |
| Handbuch und Spec-Sprache | nicht geprüft — nicht Teil des Diffs |

**Sensors.** `make docs-check`: 0 Befunde (949 Dateien). `make ablage-pruefung`: ok (26 Fälle, 0 Fehlschläge). Beide nur lesend.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 6 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Schärft-Reichweite enger als die Wirkung · Messaussage unter dem Rohdatenstand · Ableitung ohne Kennzeichnung als Befund geführt · Entscheidungsherkunft in der ADR nicht auflösbar · Sensor-Zusage weiter als der Sensor · Plan-Eintrag ohne belegte Berührung

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM. Die sechs LOW sind Wortlaut- und Präzisionsfragen, die vor `Accepted` zu bereinigen sind, weil eine Accepted-ADR nicht mehr umgeschrieben wird (`AGENTS.md` §3.5). Offen für den Projektinhaber: F-7.

**Übergabe:** Findings gehen an den Architect (ADR ist `Proposed`, noch mutable); die Finding-Klassen gehen in die Closure des Code-Slice §7. DoD-Konformität prüft der Verifier.
