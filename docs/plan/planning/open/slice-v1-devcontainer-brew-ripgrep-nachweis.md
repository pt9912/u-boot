# Slice V1: Nachweis von `rg` im Entwicklungs-Devcontainer und Lock-Schritt — Messung am Image und am Container durch den Projektinhaber, nach dem Release `v0.24.0`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Folge-Slice von [`slice-v1-devcontainer-brew-ripgrep-im-repo`](../done/slice-v1-devcontainer-brew-ripgrep-im-repo.md), das mit „`rg` nicht gemessen“ geschlossen wurde; er übernimmt dessen Rest- und Folgepunkte 1 (`rg`-Nachweis) und 2 (Lock-Schritt), §7 dort. Entscheidung des Projektinhabers 2026-11-06 (Sitzungsdialog): Der Nachweis ist ein eigener Slice; er kommt **vom Projektinhaber nach dem Release `v0.24.0`**, auf einem Docker-Host außerhalb des Agenten-Containers. Der Cut [`slice-v1-release-cut-v0.24.0`](../done/slice-v1-release-cut-v0.24.0.md) behauptet nichts über `rg` im Container und führt den Nachweis nur als Rest-Punkt. Eine Closure-Bedingung, die von der DoD dieses Slice verschieden ist, gibt es nicht.

**Bezug:** [`LH-FA-DEV-013`](../../../../spec/lastenheft.md#lh-fa-dev-013--homebrew-formeln-im-devcontainer) (Formel-Liste im Devcontainer), [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) (der eigene Devcontainer ist das erzeugte Projekt), [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) (die acht Werkzeuge der Prüfung, die Homebrew nicht überdecken darf); [`ADR-0046`](../../adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md) (Re-Evaluierungs-Trigger bei einem Befund am Block).

**Berührte Spec-Stellen:** [`SPEC-057`](../../../../spec/spezifikation.md#spec-057--homebrew-formeln-im-devcontainer-katalog-schlüssel-brew-pfade-freigabe-der-taps-konflikt-mit-homebrew-pin-meldungen-hosts-und-grenzen) (Absatz *Nicht zugesagt*, der Satz „gebaut ist sie nicht“, `spec/spezifikation.md` Zeile 1321): gelesen, in diesem Slice nicht geändert; zum Handbuch siehe §1.

**Verantwortlich:** — *(die Messung führt der Projektinhaber aus; die Rolle des Implementers, die das Protokoll in das Repo bringt und die Lock-Datei prüft, wird beim Übergang `open` → `next` gesetzt)*

**Autor:** pt9912 (Rolle Planner, Entwurf auf Entscheidung des Projektinhabers). **Datum:** 2026-11-06.

---

## 1. Ziel und Abgrenzung

**Ziel:** Mit Befehlsausgaben im Messprotokoll ist belegt, ob `rg` im `PATH` eines aus `.devcontainer/` gebauten Images und Containers steht (Login- und Nicht-Login-Shell), ob das Verzeichnis `bin` von Homebrew ein Werkzeug der Gates überdeckt, und was der Neubau mit `.devcontainer/devcontainer-lock.json` tut; der Rest steht als Nachtrag in der Evidence-Datei der Kombination und als Aussage zum Satz „gebaut ist sie nicht“.

**Herkunft.** [`slice-v1-devcontainer-brew-ripgrep-im-repo`](../done/slice-v1-devcontainer-brew-ripgrep-im-repo.md) belegt den Bau in der CI (`brew install --formula 'ripgrep'` im Log, `make gates` im fertigen Image grün), nicht `command -v rg` im fertigen Image und nicht die Container-Ebene mit den OCI-Features (E2, dort §1). Den Nachweis kann in der Sitzung des Agenten niemand erbringen: Der Container, in dem der Agent arbeitet, trägt Homebrew aus dem OCI-Feature und kein `rg`, und ein Neubau zerstört seinen Zustand (dort DoD Liefer-Punkt 3). Die Befehle stehen dort in DoD Liefer-Punkt 3 und in §7 (Rest- und Folgepunkte 1 und 2); dieser Slice übernimmt sie unverändert (§2).

**Wer was tut.** Der Projektinhaber führt die Befehle auf seinem Docker-Host aus und legt die Ausgaben vor; der Implementer bringt sie als Messprotokoll unter `docs/messungen/` in das Repo, liest die Lock-Datei und schreibt den Nachtrag. Der Implementer löst keinen Neubau des Entwicklungs-Containers aus (`docs/maintainer/devcontainer-arbeiten.md`, Abschnitt *Voraussetzungen und Umgebung*).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Kein Code, kein Template, keine Änderung an `u-boot.yaml` oder `.devcontainer/`.** Der Slice misst. Zeigt der Nachweis ein fehlendes `rg` oder ein überdecktes Werkzeug, ist das ein Vorgang für den Architect und den Projektinhaber ([`ADR-0046`](../../adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md) Re-Evaluierungs-Trigger); ein Fehler-Slice entsteht nur auf Entscheidung des Projektinhabers.
- **Keine Handänderung an `.devcontainer/devcontainer-lock.json`** (E1 des Vorgängers). Die Datei gehört dem Editor bzw. der Dev-Container-CLI; der Slice committet höchstens, was das Werkzeug geschrieben hat, und bearbeitet es nicht nach.
- **Keine Änderung von [`SPEC-057`](../../../../spec/spezifikation.md#spec-057--homebrew-formeln-im-devcontainer-katalog-schlüssel-brew-pfade-freigabe-der-taps-konflikt-mit-homebrew-pin-meldungen-hosts-und-grenzen) und des Handbuchs.** Der Satz „gebaut ist sie nicht“ steht in `spec/spezifikation.md` (Zeile 1321) und, mit anderem Wortlaut („gebaut ist die Kombination nicht“), in `docs/user/benutzerhandbuch.md` (Abschnitt *Grenzen der Messung*, Zeilen 1365 bis 1373). Eine Änderung ist ein Vorgang des Architect und braucht den gemessenen Beleg dieses Slice (E5 des Vorgängers); eine Handmessung an einem Host ist kein Test und trägt keine Zusage (Abschnitt zu Bedingungen in `AGENTS.md` §3.10).
- **Keine Pin-Gruppe für `HOMEBREW_REVISION`, kein Drift-Test des eigenen Projekts, kein `rg` in der Anwesenheits-Prüfung von `tool-pruefung`** (F1 und F2 des Vorgängers, dort §7 Rest- und Folgepunkte 4 und 5): Sie verschärfen ein Gate oder ändern das Lastenheft und bleiben bei Architect und Projektinhaber.
- **Kein `arm64`** (Rest- und Folgepunkt 8 des Vorgängers): Die Messung läuft auf dem Host des Projektinhabers; der Slice sagt nichts über eine Plattform, die er nicht misst.
- **Nicht die Messung der cgroup nach der Neuanlage** ([`slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung`](slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung.md)): ein anderer Gegenstand mit eigener DoD, auch wenn beide an einer Neuanlage durch den Projektinhaber hängen (Frage Q1 in §6).

## 2. Definition of Done

Drei Liefer-Punkte; Gate-Läufe und Closure-Pflichten zählen nicht mit.

- [ ] **Liefer-Punkt 1 — `rg`-Nachweis als Messprotokoll.** Unter `docs/messungen/slice-v1-devcontainer-brew-ripgrep-nachweis/` liegt ein Protokoll (`README.md` und Rohdaten als `.txt` oder `.log`, Ablage nach [`harness/conventions.md`](../../../../harness/conventions.md#ablage-von-messprotokollen-und-skripten); die Befehle laufen von Hand, ein Skript unter `tools/messung/` entsteht nur auf Wunsch des Projektinhabers). Es nennt je Zelle Befehl, Exit-Code und Ausgabe, dazu die Umgebung: Host, Docker-Version, Datum, `git rev-parse HEAD` des Klons und `git status` (sauber oder mit Abweichung). Die Form wählt der Projektinhaber (eine oder beide); die nicht gewählte steht im Protokoll als „nicht gemessen“, und das Protokoll sagt, was die gewählte Form nicht belegt (die Image-Ebene baut nur das Dockerfile ohne die OCI-Features, [`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](../observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md)).
  - *Image-Ebene* (das Dockerfile wie die CI, ohne OCI-Features): `docker build --no-cache -t u-boot-devcontainer:rg -f .devcontainer/Dockerfile .devcontainer`; dann `docker run --rm u-boot-devcontainer:rg bash -c 'command -v rg && rg --version'` und dasselbe mit `bash -lc` und mit `sh -lc`; dazu `docker run --rm u-boot-devcontainer:rg ls /home/linuxbrew/.linuxbrew/bin`. Kein Name in der `ls`-Ausgabe darf einer der acht Werkzeug-Namen von `tool-pruefung` sein (`go`, `golangci-lint`, `make`, `bash`, `curl`, `gh`, `git`, `docker`); der Abgleich steht im Protokoll, ein Treffer ist ein Befund. `--no-cache` hält einen warmen Layer-Cache aus dem Beleg ([`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md)).
  - *Container-Ebene* (mit den OCI-Features, der Ort, an dem `rg` gebraucht wird): Vorab ist alles Ungepushte im Container zu pushen; Neuanlage mit `make dev-up RECREATE=1`; in `make dev-shell` `rg --version`, `brew --version` und die Befehle der übrigen Features und Blöcke (`node`, `claude`, `gh`, `go`, `golangci-lint`, `docker info`); das Protokoll hält fest, ob die Neuanlage dasselbe Workspace-Volume wieder gebunden hat (R4 des Vorgängers).
- [ ] **Liefer-Punkt 2 — Lock-Schritt (nach dem Neubau auf der Container-Ebene).** `git diff -- .devcontainer/devcontainer-lock.json` im Klon des Projektinhabers; der Wortlaut des Unterschieds und **der Weg des Neubaus** (Editor oder `make dev-up`) stehen im Protokoll. Der Weg gehört dazu, weil `u-boot devcontainer up` die Dev-Container-CLI mit `--workspace-folder`, ggf. `--override-config`, `--remove-existing-container` und `--id-label` aufruft und kein Lock-Datei-Argument übergibt (`internal/adapter/driven/docker/devcontainer_launcher.go`, Funktion `runArgs`, Zeilen 542 bis 553); ob die CLI die Datei ohne ein solches Argument schreibt, ist nicht gelesen. Danach:
  - **Zeigt der Diff eine Änderung:** ein eigener Commit nur dieser Datei (Meldung mit Slice-Kennung und [`LH-FA-DEV-013`](../../../../spec/lastenheft.md#lh-fa-dev-013--homebrew-formeln-im-devcontainer)). Prüfung vor dem Commit: gültiges JSON; die Schlüssel der Datei sind genau die Schlüssel von `features` in `.devcontainer/devcontainer.json` (heute `claude-code:1.0.5` und `node:2.1.1`, Zeilen 31 und 32); die Einträge `claude-code:1.0.5` und `node:2.1.1` sind byte-gleich oder ihre Änderung ist im Protokoll begründet. Der Push mit diesem Commit ist ein neuer Tip: `ci` (Job `gates`) und `integration-docker` werden auf ihm gelesen, Run-IDs im Protokoll, Log über `gh api repos/pt9912/u-boot/actions/jobs/<id>/logs` ([`BEO-DOK/ci-log-falsch-gelesen`](../observations/BEO-DOK/ci-log-falsch-gelesen/observation.md)).
  - **Zeigt der Diff keine Änderung** (der Eintrag `homebrew:2.0.5`, Zeilen 13 bis 17 der Datei, bleibt): kein Handentfernen; das Protokoll nennt den Befund samt Weg des Neubaus, und der Ausgang geht an Architect und Projektinhaber.
  - **Wurde nur die Image-Ebene gemessen,** ist der Lock-Schritt nicht ausgeführt (die Image-Ebene berührt die Datei nicht) und bleibt als offen im Nachtrag stehen; der Slice schließt dann nicht.
- [ ] **Liefer-Punkt 3 — Nachtrag und Aussage zum Satz „gebaut ist sie nicht“.** Die Zeile *Nachtrag (offen)* in `docs/plan/planning/observations/BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf/evidence/slice-v1-devcontainer-brew-ripgrep-im-repo.md` ist durch das Ergebnis ersetzt: je Form gemessen oder nicht, Ergebnis je Zelle mit Zeiger auf das Protokoll, Wortlaut des Lock-Diffs, Run-IDs; „offen“ steht nur, was offen bleibt. Der Slice trägt außerdem eine Aussage zum Satz in [`SPEC-057`](../../../../spec/spezifikation.md#spec-057--homebrew-formeln-im-devcontainer-katalog-schlüssel-brew-pfade-freigabe-der-taps-konflikt-mit-homebrew-pin-meldungen-hosts-und-grenzen) und im Handbuch: „unverändert“ mit Grund, oder die Übergabe an den Architect mit den Fundstellen des Protokolls als Beleg; die Textänderung selbst nicht in diesem Slice (§1). Fehlt `rg` oder wird ein Werkzeug überdeckt, nennt der Nachtrag den Befund und übergibt ihn (§1, erster Punkt).
- [ ] `make gates` grün im Devcontainer-Container, am Stand nach der letzten Fix-Runde (`AGENTS.md` §3.1). *(nicht gezählt — Gate-Lauf)*
- [ ] Review (Report unter `docs/reviews/`, kein Self-Review) und Verifikation durch andere Rollen; Closure-Notiz mit Lerneintrag; Beobachtungs-Register fortgeschrieben; jedes Risiko aus §6 trägt einen Ausgang; die drei Paarungen sind getragen. *(nicht gezählt — Closure-Pflichten)*

## 3. Plan (vor Code)

Fundstellen am 2026-11-06 gelesen; der Implementer ergänzt die Tabelle im ersten Lauf und gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab ([`AGENTS.md`](../../../../AGENTS.md) §3.10).

| Datei / Komponente | Änderungs-Art | Begründung / Fundstelle |
|---|---|---|
| `docs/messungen/slice-v1-devcontainer-brew-ripgrep-nachweis/` (`README.md`, Rohdaten) | neu | Messprotokoll, Liefer-Punkt 1 und 2; nur `.md`, `.txt`, `.log` (Positivliste `make ablage-pruefung`) |
| `docs/messungen/README.md` | update | Zeile im Bestand mit Spalte *Skripte* „keine (Befehle von Hand im Protokoll)“, nach dem Muster der Zeile zu `slice-cr-lastenheft-messstaende` |
| `docs/plan/planning/observations/BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf/evidence/slice-v1-devcontainer-brew-ripgrep-im-repo.md` | update | der Nachtrag, Liefer-Punkt 3; die Datei trägt die Zeile *Nachtrag (offen)* mit dem Zeiger auf diesen Slice |
| `docs/plan/planning/observations/BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf/state.md` | update | Lese-Schritt der Closure: Zähler und Ausgang (Entscheidung des Projektinhabers steht aus) |
| `.devcontainer/devcontainer-lock.json` | update, bedingt | nur wenn der Diff nach dem Neubau eine Änderung zeigt; eigener Commit, nicht von Hand bearbeitet (Liefer-Punkt 2). Heute: Zeilen 3 bis 12 `claude-code:1.0.5` und `node:2.1.1`, Zeilen 13 bis 17 `homebrew:2.0.5`; letzter Commit `12fe9778` |
| `spec/spezifikation.md` ([`SPEC-057`](../../../../spec/spezifikation.md#spec-057--homebrew-formeln-im-devcontainer-katalog-schlüssel-brew-pfade-freigabe-der-taps-konflikt-mit-homebrew-pin-meldungen-hosts-und-grenzen)), `docs/user/benutzerhandbuch.md` | unverändert (gelesen) | Zeile 1321 und Zeilen 1365 bis 1373 tragen den Satz; die Änderung wäre ein Vorgang des Architect (§1) |
| `.devcontainer/Dockerfile`, `.devcontainer/devcontainer.json`, `u-boot.yaml` | unverändert | die Messung gilt dem Stand von `v0.24.0`; ein Befund geht an den Architect |
| `internal/adapter/driven/docker/devcontainer_launcher.go` | unverändert (gelesen) | `runArgs` (Zeile 525): kein Lock-Datei-Argument an die CLI |
| `docs/maintainer/devcontainer-arbeiten.md` | prüfen | die Werkzeugliste nennt `ripgrep` mit dem Satz „nur was gemessen ist“ (Liefer-Punkt 2 des Vorgängers); ein Befund zu `rg` kann den Satz ändern |
| `docs/plan/planning/open/README.md`, `docs/plan/planning/in-progress/roadmap.md` | update | Bestand-Zeile dieses Slice (Plan-Commit; Closure setzt den Zustand); Roadmap-Marker nach dem Lifecycle-Wechsel, solange der Slice in `in-progress/` liegt (Sensor `planning-drift`, wie beim Vorgänger) |
| `docs/reviews/<Datum>-slice-v1-devcontainer-brew-ripgrep-nachweis.md`, `docs/reviews/verification-<Datum>-slice-v1-devcontainer-brew-ripgrep-nachweis.md` | neu | Review- und Verification-Report (Plan-Tabelle führt sie, [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md)) |
| dieser Slice | neu / Moves | `make slice-mv` je Lifecycle-Schritt; Closure-Notiz vor dem `git mv` nach `done/` |

Kein Test: Der Slice ändert keinen Code. Die Befehle sind ein Handlauf des Projektinhabers (Wiederholbarkeit: Befehle im Protokoll).

## 4. Trigger

**Start** (`open` → `next` → `in-progress`): Der Release `v0.24.0` ist geschnitten (Tag gesetzt, `publish` kontrolliert; [`slice-v1-release-cut-v0.24.0`](../done/slice-v1-release-cut-v0.24.0.md) liegt in `done/`), **und** der Projektinhaber erbringt den Nachweis (Image- und/oder Container-Ebene) und legt die Ausgaben vor; kein anderer Slice in `in-progress/` (WIP-Limit 1). Ohne beides bleibt der Slice in `open/`.

**Rückführungen — vorab benannt:**

- `in-progress` → `open` (blockiert): Die Ausgaben fehlen, der Bau oder die Neuanlage des Projektinhabers scheitert, oder der Neubau ist aus einem Grund außerhalb des Repos nicht möglich; die Messung ruht, bis der Projektinhaber sie wieder aufnimmt.
- `in-progress` → `next` (zu groß): Der Nachweis zeigt einen Befund am Block, und die Behebung (Template, Spec, `.devcontainer/`) ist Arbeit eines eigenen Slice; dieser Slice schließt mit dem Befund und der Übergabe.

## 5. Closure-Trigger

DoD vollständig; das Protokoll liegt im Repo mit Befehl, Exit-Code und Ausgabe je Zelle; der Nachtrag steht in der Evidence-Datei; der Lock-Schritt ist ausgeführt oder sein Ausbleiben steht mit dem Weg des Neubaus als Befund im Protokoll; Review und Verifikation durch andere Rollen; Closure-Notiz mit Lerneintrag. Wurde nur die Image-Ebene gemessen, schließt der Slice nicht: Der Lock-Schritt bleibt offen, und der Slice bleibt in `open/` oder geht dorthin zurück.

## 6. Risiken und offene Punkte

Jedes Risiko trägt bei der Closure genau einen Ausgang (eingetreten / entfallen / weiter offen); bis dahin steht die Messstelle.

1. **`rg` fehlt oder steht nicht in allen Shell-Formen im `PATH`.** Der Block setzt `ENV PATH` und schreibt `/etc/profile.d/u-boot-brew.sh` (`.devcontainer/Dockerfile` Zeilen 124 und 125); gelesen ist der CI-Bau, nicht das fertige Image. — **Messstelle:** die drei Zellen `bash -c`, `bash -lc`, `sh -lc` (Liefer-Punkt 1). — **Ausgang:** bei der Closure; bei Eintreten ein Vorgang für Architect und Projektinhaber (§1).
2. **Ein Werkzeug von Homebrew überdeckt eines der Gates** (R3 des Vorgängers; `ENV PATH` stellt `/home/linuxbrew/.linuxbrew/bin` vor den übrigen `PATH`, `.devcontainer/Dockerfile` Zeile 125). `tool-pruefung` prüft bei `make`, `bash`, `curl`, `gh`, `git`, `docker` Floor oder Anwesenheit, nicht den Ort (`tools/harness/tool-pruefung.sh`). — **Messstelle:** die `ls`-Zelle und der Namensabgleich (Liefer-Punkt 1). — **Ausgang:** bei der Closure.
3. **Der Lock-Schritt zeigt keine Änderung, weil der Weg des Neubaus die Datei nicht schreibt.** `u-boot devcontainer up` übergibt der CLI kein Lock-Datei-Argument (Liefer-Punkt 2); ein unveränderter Diff nach `make dev-up RECREATE=1` ist dann kein Befund über den Editor-Weg und umgekehrt. — **Messstelle:** das Protokoll nennt den Weg; bei „keine Änderung“ ist die Frage, ob ein Neubau über den Editor die Datei schreibt, ausdrücklich offen. — **Ausgang:** bei der Closure; weiter offen → Architect und Projektinhaber, falls nur ein Weg gemessen wurde.
4. **Die Neuanlage kostet Zustand** (R4 des Vorgängers): Nicht Gepushtes und von Hand per `brew install` Hinzugefügtes sind danach weg; ob die Neuanlage dasselbe Workspace-Volume wieder bindet, ist nicht gelesen ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) Absatz *Nicht-Ziele*). — **Messstelle:** Push vor der Neuanlage (Liefer-Punkt 1); das Protokoll hält das Volume fest. — **Ausgang:** bei der Closure.
5. **Eine Aussage über `rg` gelangt weiter als die Messung.** Die Messung ist ein Host des Projektinhabers, `amd64` oder was das Protokoll nennt, ein Lauf je Zelle; die Maintainer-Doku und der Nachtrag sagen nur das. — **Messstelle:** der Reviewer liest jede Aussage über `rg` im Diff gegen eine Zelle des Protokolls ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)). — **Ausgang:** bei der Closure.
6. **Der Nachweis kommt nicht oder spät.** Er hängt am Projektinhaber (R8 des Vorgängers); bis dahin gilt „nicht gemessen“, und nichts im Repo behauptet `rg` im Container. — **Messstelle:** der Slice bleibt in `open/`; der Cut `v0.24.0` und die Maintainer-Doku bleiben ohne Aussage über `rg` im Container. — **Ausgang:** weiter offen, solange der Projektinhaber nicht liefert (Start-Trigger, §4).
7. **Die Landlock-Warnung des Bau-Schritts** (`Warning: Landlock ABI 10 or later is required …`, Rest- und Folgepunkt 7 des Vorgängers) tritt im Bau des Projektinhabers auf. — **Messstelle:** das Protokoll gibt die Ausgabe des Baus wieder. — **Ausgang:** bei Wiederholung auf einem anderen Host eine Beobachtung im Register (kein Carveout: nichts ist gelockert).

**Offene Fragen an den Projektinhaber (vor dem Start):**

- **Q1 — Eine Neuanlage für zwei Slices?** [`slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung`](slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung.md) hängt an derselben Neuanlage (`make dev-up RECREATE=1`, nach dem Release). Beide Slices haben eigene DoD und eigene Protokolle; das WIP-Limit 1 verlangt, dass nur einer in `in-progress/` liegt. Reihenfolge und Zuschnitt sind zu entscheiden.
- **Q2 — Form des Nachweises:** Image-Ebene, Container-Ebene oder beide? Nur die Container-Ebene führt den Lock-Schritt aus (Liefer-Punkt 2).
- **Q3 — Skript:** Genügen die Befehle von Hand, oder soll die Messung als Skript unter `tools/messung/` wiederholbar sein?

## 7. Closure-Notiz

Noch nicht geschlossen; der Slice liegt in `open/` (Entwurf, Start-Trigger in §4).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `docs/messungen/` und `docs/plan/` (DOK, Greenfield), bedingt `.devcontainer/devcontainer-lock.json` (Repo-Konfiguration, keine Sub-Area in [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area)); Code unter `internal/**` und `cmd/` wird nicht berührt, keine neue Kandidatin.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen am 2026-11-06 (der Zähler ist die Zahl der Evidence-Dateien). Treffer: [`BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf`](../observations/BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf/observation.md) (7; der Nachtrag dieses Slice, Liefer-Punkt 3) · [`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](../observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md) (10; Image-Ebene ohne OCI-Features ist nicht die Auslieferung des Devcontainers, Liefer-Punkt 1) · [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md) (1; `--no-cache`) · [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (40; Risiko 5 und der Satz in [`SPEC-057`](../../../../spec/spezifikation.md#spec-057--homebrew-formeln-im-devcontainer-katalog-schlüssel-brew-pfade-freigabe-der-taps-konflikt-mit-homebrew-pin-meldungen-hosts-und-grenzen)) · [`BEO-DOK/ci-log-falsch-gelesen`](../observations/BEO-DOK/ci-log-falsch-gelesen/observation.md) (1; CI-Läufe auf dem Tip nach dem Lock-Commit) · [`BEO-WZG/messbeleg-fehlt-oder-kommt-nach-dem-lauf`](../observations/BEO-WZG/messbeleg-fehlt-oder-kommt-nach-dem-lauf/observation.md) (2; die Ausgaben stehen als Protokoll im Repo, nicht nur in der Sitzung) · [`BEO-WZG/kopplung-nicht-nachgezogen`](../observations/BEO-WZG/kopplung-nicht-nachgezogen/observation.md) (2; Lock-Datei als bewusst benannter Nachzug-Ort) · [`BEO-APP/podman-static-pin-alterung`](../observations/BEO-APP/podman-static-pin-alterung/observation.md) (6; `HOMEBREW_REVISION`, nicht Gegenstand, §1) · [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (89; dieser Plan nennt Zeilen und Namen nur nach Lesen, §3). Sonst keine Treffer.

**Modus-Begründungsblock:** entfällt — alle berührten Sub-Areas GF.
