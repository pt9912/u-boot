# Messprotokoll: überlebt eine Hand-Ergänzung von `--memory` und `--cpus` ein erneutes `generate devcontainer`?

Messung zum Entwurf [`slice-v1-devcontainer-ressourcenlimits`](../../plan/planning/open/slice-v1-devcontainer-ressourcenlimits.md) (Anlass: Frage des Projektinhabers vom 2026-10-31, ob man Devcontainern schon eine Höchstmenge an CPU und Arbeitsspeicher mitgeben kann). Gemessen ist nur, was `u-boot generate devcontainer` mit der Datei `.devcontainer/devcontainer.json` tut. Es lief **kein Container**; ob eine Engine `--memory` und `--cpus` aus `runArgs` anwendet, ist nicht Gegenstand.

## Messstand

| Eigenschaft | Wert |
|---|---|
| Datum | 2026-10-31 |
| Binary | `u-boot version 0.16.1` (installiert, im `PATH`) |
| Host | Linux, Wegwerf-Ordner (`u-boot init --no-git --no-interactive demo`), keine Docker-Engine beteiligt |
| Läufe | **ein Lauf je Variante** (keine Streuung bestimmt) |
| Durchführung | Hauptagent, Erstlauf mit zwei Skripten aus dem Scratchpad der Sitzung (`messen-limits.sh`, `messen-limits2.sh`) |

## Skript und Aufruf

[`tools/messung/devcontainer-ressourcenlimits/messen.sh`](../../../tools/messung/devcontainer-ressourcenlimits/messen.sh), aus der Repo-Wurzel:

```bash
bash tools/messung/devcontainer-ressourcenlimits/messen.sh > docs/messungen/slice-v1-devcontainer-ressourcenlimits/roh-lauf.txt 2>&1
```

**Das abgelegte Skript ist nicht gelaufen.** Es ist die bereinigte Fassung der beiden Skripte des Erstlaufs (ein Arbeitsordner mit `mktemp` statt eines festen Pfads, ein leerer Platzhalter-Helfer entfernt, die Varianten in einer Datei, Reihenfolge A1, A2, B, A3). Die Rohdaten in diesem Verzeichnis stammen aus dem Erstlauf; `roh-lauf.txt` (der Aufruf oben) liegt nicht vor. Wer die Befunde belegen will, fährt das Skript und gleicht die Dateien ab. Von den Skripten des Erstlaufs wurde nur die Ausgabe von `generate` je Variante (`gen*.out`) und die erzeugte Datei vor der Änderung gesichert, nicht die Standardausgabe der Skripte; die Zeile mit dem Exit-Code von `generate` steht darum nur in der Rückmeldung des Hauptagenten, nicht in einer Rohdatei.

## Ergebnis je Variante

| Variante | Eingriff | Beobachtung | Rohdaten |
|---|---|---|---|
| Ausgangslage | `generate devcontainer`, Standardprofil | Datei angelegt (`created`); Block ohne `runArgs`-Zeile | [`roh-generate-1-standardprofil.txt`](roh-generate-1-standardprofil.txt), [`roh-devcontainer-json-standardprofil.txt`](roh-devcontainer-json-standardprofil.txt) |
| A1 | Zeile `"runArgs": ["--memory=4g", "--cpus=2"],` **innerhalb** des Managed Blocks (hinter `"name"`), Standardprofil, danach `generate devcontainer` | Die Zeile ist **entfernt**; Exit `0`, Log `updated-block`, Meldung „Updated devcontainer managed block“, **keine Warnung** | [`roh-generate-a1-innerhalb-des-blocks.txt`](roh-generate-a1-innerhalb-des-blocks.txt) |
| A2 | wie A1, mit `generate devcontainer --force` | `u-boot: unknown flag: --force`; das Flag existiert für diesen Befehl nicht (Exit-Code nach der Rückmeldung des Hauptagenten `2`) | [`roh-generate-a2-force.txt`](roh-generate-a2-force.txt) |
| B | Kommentarzeile `// HAND: --memory=4g` **hinter** `// END U-BOOT MANAGED BLOCK: init`, danach `generate devcontainer` | Die Zeile bleibt erhalten; „devcontainer already up to date; no changes.“ | [`roh-generate-b-ausserhalb-des-blocks.txt`](roh-generate-b-ausserhalb-des-blocks.txt) |
| A3 | Sandbox-Profil (`devcontainer.profile: sandbox`), ohne weitere Features | Der Generator rendert **keine** `runArgs`-Zeile; die Probe, ein vorhandenes `runArgs` zu ergänzen, fand keine Zeile. Das zweite `generate` meldete „already up to date“. **A3 sagt über das Überleben einer Ergänzung im Sandbox-Profil nichts.** | [`roh-generate-a3-sandbox-1.txt`](roh-generate-a3-sandbox-1.txt), [`roh-generate-a3-sandbox-2.txt`](roh-generate-a3-sandbox-2.txt), [`roh-devcontainer-json-sandboxprofil.txt`](roh-devcontainer-json-sandboxprofil.txt) |

**Verworfene Variante.** Die erste Fassung der Variante A im Erstlauf (`messen-limits.sh`) ersetzte per `awk` eine vorhandene `runArgs`-Zeile; im Standardprofil gibt es keine, die Datei blieb unverändert, das zweite `generate` meldete „already up to date“. Die Variante war wirkungslos, ihr Ergebnis wird nicht benutzt und ihr Skriptteil ist nicht übernommen. Gültig ist A1 (`messen-limits2.sh`), die Einfügung hinter der Zeile `"name"`.

## Lesehilfe

- Der Managed Block umschließt in beiden gemessenen Dateien das gesamte JSON-Objekt (siehe die beiden `roh-devcontainer-json-*.txt`). Daraus folgt, dass außerhalb des Blocks kein Schlüssel des Objekts stehen kann; das folgt aus der Form des Blocks und ist **nicht** mit der Dev-Container-CLI oder VS Code gemessen (ob sie eine Datei mit Text hinter dem Objekt lesen).
- Daraus für den Slice: Eine Begrenzung von Hand in der `devcontainer.json` überlebt ein erneutes `generate devcontainer` nicht; ein Weg über `u-boot` fehlt.

## Nicht gemessen

- Ob `--memory` und `--cpus` im Container ankommen, auch nach `u-boot devcontainer up` (Hilfs-Container mit Dev-Container-CLI): kein Container gestartet.
- Das Sandbox-Profil (`--privileged`, `nestedRuntime: docker`) mit Begrenzung; die Wirkung auf den verschachtelten Docker und dessen Container.
- Docker-Engines außer Linux `amd64` mit lokaler Engine; ein zu kleines Limit (OOM-Kill von Startscripten); zwei Instanzen mit Begrenzung.
- Das Überleben einer Ergänzung im Sandbox-Profil **mit** einer `runArgs`-Zeile (etwa mit `nestedRuntime: docker`, das `--privileged` rendert); gemessen ist nur das Profil ohne weitere Features.
- Mehr als ein Lauf je Variante; andere Versionen von `u-boot` als 0.16.1.
