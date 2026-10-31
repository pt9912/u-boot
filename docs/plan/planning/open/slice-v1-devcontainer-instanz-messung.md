# Slice V1: Mehrere Devcontainer je Projekt — Messung der Volume-Namen, des Id-Labels und des Verhaltens ohne Flag (zwei Instanzen, Dev-Container-CLI aus dem Hilfs-Image)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle (wie der Code-Slice; Release-Zuordnung beim Projektinhaber, die Roadmap ist nicht geändert). Erster von drei Slices (Messung → Lastenheft-CR → Code): ihm folgen [`slice-cr-devcontainer-instanz`](slice-cr-devcontainer-instanz.md) (Lastenheft-CR mit den Messwerten) und [`slice-v1-devcontainer-instanz-label`](slice-v1-devcontainer-instanz-label.md) (Code, Tests, Docker-Test, Doku). Der Schnitt folgt dem Vorbild [`slice-v1-devcontainer-ollama-messung`](../done/slice-v1-devcontainer-ollama-messung.md) → [`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md) → [`slice-v1-devcontainer-feature-ollama`](../done/slice-v1-devcontainer-feature-ollama.md): Die Zusage im Lastenheft nennt nur, was hier gemessen ist.

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (`up`, `down` und `shell`), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil: Workspace im benannten Volume), [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Speicher von Docker und Podman im benannten Volume); [ADR-0028](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) (Hilfs-Image mit der Dev-Container-CLI; aktive ADR, nur die Fundstellen zu `${devcontainerId}` gelesen).
Anlass: Projektinhaber, 2026-10-31 (Weg A des Entwurfs [`slice-v1-devcontainer-instanz-label`](slice-v1-devcontainer-instanz-label.md) entschieden; Anforderung: Instanzen teilen kein Volume).

**Berührte Spec-Stellen:** — (der Slice liest [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) und ändert keine Stelle; Lastenheft, Spezifikation und Handbuch folgen in den beiden Folge-Slices).

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-31.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Mit gemessenen Belegen steht fest, was die Dev-Container-CLI des Hilfs-Images tut, wenn zwei Instanzen desselben Projektordners über ein zweites Id-Label unterschieden werden: eigene Container, eigene Volumes (Workspace, Docker-Speicher, Podman-Speicher, Modelle) oder geteilte; was `up` ohne Flag neben einer vorhandenen Instanz tut; welche Namen für das Label und die Docker-Namen taugen. Ergebnis: ein Messprotokoll und eine Messwert-Tabelle, aus der der Lastenheft-CR ([`slice-cr-devcontainer-instanz`](slice-cr-devcontainer-instanz.md)) seine Platzhalter füllt. Zugesagt wird später nur, was hier gemessen ist.

**Warum ein eigener Mess-Slice (Prüfung des Schnitts, Planner, 2026-10-31).** Die Anforderung des Projektinhabers (Instanzen teilen kein Volume) macht aus der Volume-Frage eine Zusage; der Code-Entwurf [`slice-v1-devcontainer-instanz-label`](slice-v1-devcontainer-instanz-label.md) führte die Messung bisher als ersten Schritt *im* Code-Slice und das Lastenheft-CR als Bedingung seines Starts. Das kehrt die Reihenfolge der Abhängigkeit um: Der Wortlaut des CR hängt vom Messwert ab (Ausgang der Volume-Messung entscheidet, ob die Zusage „eigene Volumes“ lautet oder ob der Befehl bei fehlender Isolation mit Exit-Code `10` abbricht), und der Code hängt vom CR. Wie bei Ollama steht die Messung deshalb vor dem CR (Messung → CR → Code); drei Slices statt einem mit drei Phasen, weil jeder einzeln lieferbar ist und der CR eine Freigabe des Projektinhabers braucht, auf die der Code nicht warten darf, wenn die Messung schon etwas Eigenes geliefert hat.

### Gefundener Stand (gelesen am 2026-10-31; [`AGENTS.md`](../../../../AGENTS.md) §3.10)

- **Die Volume-Namen stehen mit `${devcontainerId}` in der erzeugten Konfiguration:** `volumeMount` (`internal/hexagon/application/devcontainer_sandbox.go`, Zeile 140) liefert `source=<Volume>-${devcontainerId},target=<Ziel>,type=volume`; die Tests nennen vier Familien: `workspace`, `containers`, `docker`, `ollama` (`internal/hexagon/application/generate_ollama_test.go`, Zeilen 27 bis 30). Dass die Dev-Container-CLI `${devcontainerId}` aus den Id-Labels ableitet und es deshalb mit dem Id-Label wechselt, ist aus der Spezifikation der CLI angenommen, **nicht gemessen** ([ADR-0028](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) führt die Frage nur für den Editor-Vergleich, Zeile 41).
- **Der Aufruf der CLI** (`runArgs`, `internal/adapter/driven/docker/devcontainer_launcher.go`, Zeilen 284 bis 303): `docker run --rm --name <Name> --mount <Socket> --mount <Projektordner> [-e <NAME>…] <Image> devcontainer up --workspace-folder <Ordner> [--remove-existing-container]`; kein `--id-label`.
- **Die Suche** (`FindContainer`, `internal/adapter/driven/docker/inspector.go`, Zeilen 75 bis 100): `docker ps -a --no-trunc --filter label=devcontainer.local_folder=<Ordner>` und der erste Treffer. Ohne Flag trifft dieser Filter auch einen Container, der zusätzlich ein Instanz-Label trägt — ob und wie `up` ohne Flag neben einer Instanz einen eigenen Container anlegt, ist Messfrage M6.
- **Es gibt schon einen Docker-Test mit zwei Instanzen, er misst die CLI nicht:** `TestE2E_SandboxDevcontainer_TwoInstancesIsolated` (`internal/e2e/sandbox_devcontainer_docker_test.go`, ab Zeile 557) startet zwei Container mit `docker run` und ersetzt `${devcontainerId}` von Hand durch `idaaaa` und `idbbbb`; er belegt die Isolation der Mounts bei verschiedenen Werten, nicht, dass die CLI den Wert beim zweiten Id-Label wechselt.
- **Version der CLI:** `@devcontainers/cli` `0.89.0` (`internal/adapter/driven/docker/devcontainercli/package.json`, Abhängigkeit); das Hilfs-Image trägt den Tag `<CLI-Version>-<12 Hex des Kontexts>` (`helperImageTag`, `devcontainer_launcher.go`, Zeile 370; Rumpf nicht gelesen).
- **Das Standardprofil hat kein Workspace-Volume:** Der Helfer des Tests oben (`generatedMounts`) verlangt `workspaceMount` leer für das Sandbox-Profil (Zeile 551 bis 553); für das Standardprofil nennt das Handbuch den Bind-Mount des Projektordners (`docs/user/benutzerhandbuch.md`, Zeile 1274: Zeile *Workspace*, Spalte Sandbox-Profil: benanntes Volume, Spalte Standardprofil: Bind-Mount). Zwei Instanzen im Standardprofil teilen den Projektordner; ob die Messung das bestätigt, ist M3.
- **Ältere Messprotokolle als Vorbild:** `docs/messungen/slice-v1-cli-devcontainer-up-down/README.md` (Aufruf eines Skripts mit dem ausgelieferten Binary, Rohdaten als `.txt`), Skripte unter `tools/messung/devcontainer-start-befehl/` (Namen gelesen, Inhalt nicht).

### Messfragen (Arbeitsplan; Linux `amd64`, lokale Engine, Hilfs-Image mit der Dev-Container-CLI `0.89.0`)

Aufbau: Das Skript legt einen Wegwerf-Projektordner an (`u-boot init --devcontainer --sandbox`, Sandbox-Profil mit `nestedRuntime: docker` und dem Feature `ollama`, damit alle vier Volume-Familien erscheinen; eine zweite Probe im Standardprofil) und ruft die CLI so auf, wie `runArgs` es tut (Spiegelpfad, Socket), einmal ohne und zweimal mit `--id-label`. Die Befehle stehen wörtlich im Skript. Die Messenden führen es als Docker-Aufrufe aus, ohne Host-Go und ohne Host-Node.

| # | Frage | Warum |
|---|---|---|
| M0 | Messstand: Docker-Client und -Server, Engine-Name und -Treiber, Kernel, `uname -m`, Dev-Container-CLI-Version (aus dem Hilfs-Image: `devcontainer --version`), Tag des Hilfs-Images, `u-boot version` | Zusage-Menge des CR |
| M1 | Containerkennungen und Labels zweier gleichzeitig laufender Instanzen desselben Ordners (`docker inspect`: `devcontainer.local_folder`, `devcontainer.config_file`, die Id-Labels, `Config.Labels` insgesamt) mit `--id-label devcontainer.local_folder=<Ordner> --id-label u-boot.instance=<Name>` (Arbeitsnamen des Entwurfs) | trägt, ob `FindContainer` auf zwei Labels filtern kann |
| M2 | Wechselt `${devcontainerId}` mit dem zweiten Id-Label? Der aufgelöste Wert je Instanz (aus den Quellnamen der Mounts des Containers); derselbe Ordner ohne Flag als Bezug; eine zweite Messung mit nur dem Label `devcontainer.local_folder` explizit (Gleichheit mit dem Wert ohne Flag?) | **Kernfrage der Anforderung** (Instanzen teilen kein Volume) |
| M2b | *Nur wenn M2 ergibt, dass der Wert nicht wechselt:* Lösen `mounts` einen weiteren Platzhalter auf, der pro Instanz verschieden ist (etwa `${localEnv:NAME}` mit einem per `-e NAME` durchgereichten Wert, wie `runArgs` es für Variablen kann), und überschreibt oder ergänzt die CLI die Mounts auf anderem Weg? | Kandidaten für das Schließen der Lücke; keine Wahl, nur Messung |
| M3 | Volume-Namen und -Zahl je Instanz für Workspace, Docker-Speicher, Podman-Speicher (`nestedRuntime: podman` in einer eigenen Probe), Modelle (`ollama`): `docker volume ls` vor und nach; zwei Instanzen verschiedene Volumes? Inhaltsprobe: eine Datei im Workspace von Instanz 1 ist in Instanz 2 nicht sichtbar; ein `docker pull`-Image im verschachtelten Docker von Instanz 1 fehlt in Instanz 2. Im Standardprofil: Sehen beide Instanzen dieselbe Datei im Projektordner (Bind-Mount)? | Zusage „eigene Volumes“; Grenze des Standardprofils |
| M4 | Image: Wird das Image des Projekts für die zweite Instanz neu gebaut oder geteilt (Image-ID, `docker images`, Dauer des zweiten `up` bei warmem Cache)? | Platzbedarf und Dauer des zweiten Starts |
| M5 | Dauer (Sekunden) und Platzbedarf des zweiten Starts: `docker system df -v` und `docker volume ls` vor und nach, Größe der neuen Volumes; kalter Go-Cache im neuen Workspace-Volume nur benannt, nicht aufgefüllt | Ressourcen-Risiko des Entwurfs |
| M6 | Verhalten von `up` ohne Flag, wenn eine Instanz läuft: legt die CLI einen eigenen Container an oder übernimmt sie den der Instanz? Welche Labels trägt der Container ohne Flag, trifft `docker ps --filter label=devcontainer.local_folder=<Ordner>` beide? Umgekehrt: Instanz anlegen, wenn schon ein Container ohne Flag existiert (Bestandscontainer) | Standardverhalten ohne Flag, Bestandscontainer (Risiko im Code-Slice) |
| M7 | `--remove-existing-container` mit zweitem Id-Label: entfernt die CLI nur den Container der Instanz? `docker stop` des einen Containers lässt den anderen laufen? Bleiben die Volumes beider erhalten? | `--recreate`, `down` je Instanz |
| M8 | Namen: Welche Instanznamen nimmt die CLI als Label-Wert und was wird daraus (Länge 1, 32, 63, 64, 128; Großbuchstaben; `.`, `_`, `-` am Rand; Leerzeichen; Anführungszeichen; Nicht-ASCII; leere Zeichenkette)? Taugt der Name im Filter von `docker ps` (Zeichen, die der Filter anders liest)? | Zeichensatz und Länge als Vorschlag mit Begründung im CR |
| M9 | Aufräumen: Alles, was die Messung anlegt, trägt ein Präfix und wird am Ende entfernt (Container, Volumes, Images, Projektordner); `docker ps -a`, `docker volume ls` und `docker images` vorher und nachher, Differenz leer | Messung hinterlässt nichts (vgl. [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md)) |

**Grenzen vorab:** Linux `amd64`, lokale Engine, zwei Instanzen, ein Host; die CLI-Version aus M0. Ein Urteil über drei oder mehr Instanzen, über andere Hosts (Colima, Docker Desktop, rootless, `arm64`) und über die Last ist nicht Gegenstand. Mindestens zwei Läufe, datiert, mit Umgebung; der zweite Lauf auf frischem Docker-Stand (Cache des ersten benannt).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Produktcode und Tests** — [`slice-v1-devcontainer-instanz-label`](slice-v1-devcontainer-instanz-label.md); ein Slice, der misst, setzt nicht um (Rollentrennung, `harness/roles.md`).
- **Lastenheft, Spezifikation, Handbuch, `CHANGELOG.md`** — [`slice-cr-devcontainer-instanz`](slice-cr-devcontainer-instanz.md) (Lastenheft, mit Freigabe des Projektinhabers) und der Code-Slice (Spezifikation, Handbuch, `CHANGELOG.md`); das Protokoll nennt keine Zusage.
- **Das Schließen einer Lücke bei M2** — Ergibt M2, dass die Volumes geteilt würden, ist das Schließen (etwa über das Schema der Mounts) ein Vorgang mit Wirkung auf den Generator und seine Golden Cases; die Messung benennt die Kandidaten (M2b), der Fall geht als Rückführung an den Planner (§4) und wird ein eigener Slice vor dem CR.
- **Die Wahl des Weges** — mit Weg A entschieden (Projektinhaber, 2026-10-31); dieser Slice prüft nur, ob die Annahmen von Weg A tragen.
- **Mehr als zwei Instanzen, andere Hosts, Last** — ein Schluss über die Messung hinaus ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)).

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Messung** ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)): M0 bis M9 sind gemessen oder als „nicht gemessen“ mit Grund ausgewiesen (M2b nur, wenn M2 es verlangt); mindestens zwei Läufe, datiert, mit Messstand (Docker, Engine, Kernel, Architektur, CLI-Version, Tag des Hilfs-Images). *Beleg:* Messprotokoll `docs/messungen/slice-v1-devcontainer-instanz-messung/` (`README.md` mit Messstand, Aufruf, Lesehilfen; Rohdaten als `rohdaten-*.txt`), Skript und Hilfsdateien unter `tools/messung/devcontainer-instanz/`, Zeile im Bestand von `docs/messungen/README.md` ([`harness/conventions.md`](../../../../harness/conventions.md#ablage-von-messprotokollen-und-skripten)).
- [ ] **Auswertung und Übergabe:** Eine Messwert-Tabelle in §7 (Frage → Wert → Rohdaten-Datei → Platzhalter des CR, der davon abhängt) und ein Urteil zur Kernfrage in genau einer von drei Fassungen: (1) `${devcontainerId}` wechselt, alle gemessenen Volume-Familien sind je Instanz getrennt; (2) wechselt nicht oder nur teilweise — dann benennt die Tabelle die betroffenen Familien und die Kandidaten aus M2b mit Messwert; (3) nicht entscheidbar mit Grund. Die Platzhalter in [`slice-cr-devcontainer-instanz`](slice-cr-devcontainer-instanz.md) §1 sind gegen die Tabelle geprüft und, wo der Messwert abweicht, im Wortlaut-Entwurf berichtigt (nur gelesene Aussagen).
- [ ] **Annahmen des Code-Slice geprüft:** [`slice-v1-devcontainer-instanz-label`](slice-v1-devcontainer-instanz-label.md) §3 und §6 sind gegen das Ergebnis berichtigt (Test-Aufträge für die Docker-Tests: welche Volume-Namen, welches Verhalten ohne Flag, welche Namen gültig; Risiken *Volumes und Image je Instanz*, *Mehrdeutige Suche*, *Bestandscontainer*, *Ressourcen* mit dem Messwert).
- [ ] `make gates` grün (im Modus der Umgebung).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `docs/messungen/README.md` (Bestand); kein Handbuch, kein `CHANGELOG.md` (kein Produktverhalten).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — insbesondere [`BEO-APP/container-identitaet-nur-ueber-label`](../observations/BEO-APP/container-identitaet-nur-ueber-label/observation.md) (M1 und M6 sind die „Messung mit zwei Containern am selben Label“, sein Trigger); **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

Gelesen am 2026-10-31; der Implementer prüft die Zeilen am Ist-Stand ([`AGENTS.md`](../../../../AGENTS.md) §3.10). Kein Produktcode und kein Test: `internal/**`, `spec/**`, `docs/user/`, `README.md`, `README.de.md`, `CHANGELOG.md` und `Makefile` bleiben ohne Diff.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/messung/devcontainer-instanz/` (Skript, Hilfsdateien) | neu | Aufbau und Messung M0 bis M9; Aufruf wie das Vorbild `tools/messung/devcontainer-start-befehl/` (Inhalt nicht gelesen) |
| `docs/messungen/slice-v1-devcontainer-instanz-messung/` (`README.md`, `rohdaten-*.txt`) | neu | Messprotokoll mit Messstand und Rohdaten ([`docs/messungen/README.md`](../../../messungen/README.md): Ablage) |
| `docs/messungen/README.md` | update | Zeile im Bestand (Verzeichnis → Skriptpfad) |
| `docs/plan/planning/open/slice-cr-devcontainer-instanz.md` | update | Platzhalter des Wortlaut-Entwurfs gegen die Messwerte füllen oder berichtigen |
| `docs/plan/planning/open/slice-v1-devcontainer-instanz-label.md` | update | §3 und §6 gegen das Ergebnis berichtigt (siehe DoD) |
| `docs/plan/planning/observations/BEO-APP/container-identitaet-nur-ueber-label/` (`evidence/slice-v1-devcontainer-instanz-messung.md` neu, `state.md`) | neu / update | Beleg der Messung mit zwei Containern am selben Label (Closure) |
| `docs/reviews/` (Review-Report, Verification Evidence) | neu | Review und Verifikation nach `harness/review.md` und `harness/verification.md` |
| diese Datei, `docs/plan/planning/open/README.md` | update | Messung (§7), Closure; Zeile des Slice; Moves per `make slice-mv` mit Verweis-Nachzug |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; kein anderer Slice in `in-progress/` (WIP-Limit 1); der Projektinhaber hat den Slice priorisiert. Kein Vorgänger: Weg A ist entschieden (Projektinhaber, 2026-10-31).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Messung verlangt mehr als zwei Instanzen oder mehr als einen Host, um die Kernfrage M2 zu beantworten.
- `in-progress` → `open` (blockiert): Zwei gleichzeitige Instanzen lassen sich nicht herstellen (Engine, Speicher, Hilfs-Image nicht baubar).
- **Ausgang (2) der Kernfrage** (Volumes geteilt oder teilweise): Der Slice schließt mit dem Befund, nicht mit einer Lösung; der Planner schneidet vor dem CR einen Slice für das Schließen der Lücke (Schema der Mounts, Generator, Golden Cases) oder legt den Abbruch mit Exit-Code `10` als Wortlaut des CR vor — beides Entscheidung des Projektinhabers (Anforderung 2026-10-31, [`slice-v1-devcontainer-instanz-label`](slice-v1-devcontainer-instanz-label.md) §6).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; `make gates` grün; Messprotokoll und Messwert-Tabelle liegen vor; die Platzhalter des CR sind gefüllt oder als „nicht gemessen“ gestrichen; Review-Report und Verification-Report liegen vor; die Closure-Notiz mit Steering-Loop-Lerneintrag ist geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Zusage weiter als die Messung** ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)): zwei Instanzen auf einem Host und einer CLI-Version tragen keine Zusage für viele Instanzen oder andere Hosts; der CR übernimmt Messbedingungen und Grenzen Satz für Satz. — **Ausgang:** bei Closure.
- **Die Messung ahmt die CLI-Aufrufe von `u-boot` nach, statt `u-boot devcontainer up` zu fahren** (das Flag gibt es noch nicht). — **Ausgang:** bei Closure (die Messung nennt, dass `runArgs` nachgebaut ist; der Docker-Test des Code-Slice fährt den Befehl selbst).
- **Die Messung hinterlässt Container, Volumes oder Images** (M9). — **Ausgang:** bei Closure.
- **Der Hilfs-Container läuft mit dem Docker-Socket des Hosts** (Lockerung, ausgewiesen im Befehl): Das Skript fährt dieselbe Lockerung; sie gehört in den Messstand des Protokolls. — **Ausgang:** bei Closure.
- **Die Kernfrage M2 ergibt nicht eindeutig „wechselt“ oder „wechselt nicht“** (etwa nur für einen Teil der Familien). — **Ausgang:** bei Closure (Fassung (3) der DoD).

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt: Was hat funktioniert · Was ging anders als geplant · Messwert-Tabelle (Frage → Wert → Rohdaten → Platzhalter des CR) · Urteil zur Kernfrage · Steering-Loop-Eintrag · Beobachtungs-Register · Folge-Slices · Risiken aus §6 · Drei Paarungen.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `tools/` (WZG, Messskript) sowie `harness/`, `docs/plan/` und `docs/messungen/` (DOK); beide sind in der Modus-Deklaration von [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) als Greenfield auditiert. Kein Produktcode, keine Code-Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Das Register ([`observations/`](../observations/)) ist nach den Titeln durchgesehen. Gelesen (Titel und Kopf des Entwurfs, nicht der Volltext): [`BEO-APP/container-identitaet-nur-ueber-label`](../observations/BEO-APP/container-identitaet-nur-ueber-label/observation.md) (ein Beleg, sein Trigger „Messung mit zwei Containern am selben Label“ tritt mit diesem Slice ein). Nur nach Titel gesichtet, Inhalt beim Start zu lesen: [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md) (Aufräumen, M9), [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md), [`BEO-DOK/messprotokoll-ausserhalb-des-repos`](../observations/BEO-DOK/messprotokoll-ausserhalb-des-repos/observation.md) (Ablage; nicht gelesen).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
