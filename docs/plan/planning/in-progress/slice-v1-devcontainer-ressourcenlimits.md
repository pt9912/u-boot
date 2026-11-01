# Slice V1: Devcontainer begrenzen — Arbeitsspeicher und CPU aus `u-boot.yaml` als `--memory` und `--cpus` in `runArgs`, Weg A gewählt (Projektinhaber, 2026-11-01)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. `v0.17.0` ist veröffentlicht und enthält diesen Slice nicht; die Zuordnung zu einem Release ist offen und die Priorisierung entscheidet der Projektinhaber, die Roadmap ist nicht geändert. Eine Closure-Bedingung, die von der DoD dieses Slice verschieden ist, gibt es nicht. Anlass: Frage des Projektinhabers am 2026-10-31, ob man den Devcontainern schon eine Höchstmenge an CPU und Arbeitsspeicher mitgeben kann; er hat gesagt, den Gegenstand zu klären und bei Bedarf einen Slice anzulegen.

**Entscheidung des Projektinhabers (2026-11-01):** Weg A (zwei Felder `devcontainer.resources.memory` und `devcontainer.resources.cpus` mit Validierung, als `--memory` und `--cpus` in `runArgs`, in beiden Profilen). Der Slice ist ein Entwurf: nichts ist umgesetzt, kein Lastenheft-, Spezifikations- oder Codetext ist geändert. Die Umsetzung startet der Projektinhaber in einer eigenen Devcontainer-Instanz (`u-boot devcontainer up --instance`), ohne Termin. Was noch beim Projektinhaber liegt, steht am Anfang von §6.

**Bezug:** [`LH-FA-CONF-002`](../../../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) (Inhalt von `u-boot.yaml`, Exit-Code `10` bei ungültigen Werten), [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) (Erzeugung der Devcontainer-Konfiguration), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil: Ausweisung jeder Sicherheitslockerung), [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (`nestedRuntime: docker` mit `--privileged`), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (Exit-Codes); je nach Weg [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (`up`, Zusage nur Linux-Host `amd64` mit lokaler Engine). Eine ADR zu Ressourcenlimits gibt es nicht (Suche nach `memory`, `--cpus` und `ressourcen` in `docs/plan/adr/`: ein Treffer in [ADR-0023](../../adr/0023-podman-aus-podman-static-im-sandbox-devcontainer.md) in einer Messzeile ohne Bezug).

**Berührte Spec-Stellen:** [`SPEC-005`](../../../../spec/spezifikation.md#spec-005--schema-der-projektkonfiguration-u-bootyaml) (Schema von `u-boot.yaml`) · [`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config) (Tabelle der erlaubten Pfade und Schreibpipeline) · je nach Weg [`SPEC-052`](../../../../spec/spezifikation.md#spec-052--feature-ollama-host-im-devcontainer-eintragstyp-container-option-meldung-und-grenzen) (`runArgs`-Eintrag von `ollama-host`; nicht gelesen, prüfen).

**Verantwortlich:** pt9912

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-31.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Wer einem Devcontainer eine Höchstmenge an Arbeitsspeicher und CPU mitgeben will, setzt sie in `u-boot.yaml` (Arbeitsnamen `devcontainer.resources.memory` und `devcontainer.resources.cpus`), und `u-boot generate devcontainer` rendert sie als `--memory` und `--cpus` in `runArgs`, in beiden Profilen; ohne Angabe bleibt die erzeugte Datei wie heute — oder der Slice sagt nach der Messung, was nicht zugesagt werden kann.

**Befund** (gelesen und gemessen am 2026-10-31, Fundstellen je Punkt; [`AGENTS.md`](../../../../AGENTS.md) §3.10):

- **Heute gibt es keine Option dafür.** Die Suche nach `hostRequirements`, `--memory`, `--cpus`, `memory_limit` und `cpuLimit` in `spec/`, `internal/` (darin die Vorlagen) und `docs/user/` hat keinen Treffer. `runArgs` der erzeugten `devcontainer.json` rendert die Vorlage `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl` aus `.RunArgs` (Zeilen 14 bis 15 im Sandbox-Profil, Zeilen 39 bis 40 im Standardprofil); gefüllt wird die Liste in `devcontainerTemplateData` (`internal/hexagon/application/devcontainer_sandbox.go`, Zeilen 89 bis 136) nur aus festen Quellen: `hostGatewayRunArgs` (`devcontainer_features.go`, Zeile 573; der Eintrag des Features `ollama-host`, in beiden Profilen, Zeilen 97 und 132), `podmanRelaxations` (Zeile 281, sechs Optionen), `dockerRelaxation` (Zeile 163, `--privileged`) und `egressRelaxation` (Zeile 166, `--cap-add=NET_ADMIN`).
- **`u-boot.yaml` hat kein Feld.** `ubootYAMLDevcontainer` (`internal/hexagon/application/initproject.go`, Zeilen 69 bis 79) trägt `enabled`, `featureSources`, `features`, `user`, `profile` und `sandbox`; [`LH-FA-CONF-002`](../../../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) (Zeile 829 im Lastenheft) zählt die optionalen Schlüssel einzeln auf, ein Schlüssel für Ressourcen ist nicht darunter.
- **Von Hand geht es, hält aber nicht.** Messung am 2026-10-31 (u-boot 0.16.1, Linux, Wegwerf-Ordner; Protokoll [`slice-v1-devcontainer-ressourcenlimits`](../../../messungen/slice-v1-devcontainer-ressourcenlimits/README.md), Einzelmessung je Variante, kein Container gestartet): Eine Zeile `"runArgs": ["--memory=4g", "--cpus=2"]` **innerhalb** des Managed Blocks von `.devcontainer/devcontainer.json` (Standardprofil) entfernt das nächste `u-boot generate devcontainer` ohne Warnung (Exit `0`, Log `updated-block`, Meldung „Updated devcontainer managed block“). Eine Kommentarzeile **hinter** `// END U-BOOT MANAGED BLOCK: init` bleibt erhalten; der Block umschließt aber das gesamte JSON-Objekt, ein Schlüssel kann außerhalb nicht stehen (aus der Form des Blocks, mit der Dev-Container-CLI nicht gemessen). Das Flag `--force` hat der Befehl nicht. Im Sandbox-Profil ohne weitere Features rendert der Generator keine `runArgs`-Zeile; für dieses Profil sagt die Messung über das Überleben einer Ergänzung nichts. Ein Weg über `u-boot` fehlt also.
- **Die Prüfung des Werts sitzt an einer Stelle, die dieses Feld nutzen kann.** `validateDevcontainer` (`devcontainer_sandbox.go`, Zeile 19) ist der gemeinsame Einstieg von `config`, `generate` und `doctor` (so sein Doc-Kommentar; die Aufrufstellen sind nicht einzeln gelesen); `validateDevcontainerSandbox` (Zeile 32) wrappt Fehler in `domain.ErrInvalidSandboxSetting`, das die CLI auf Exit-Code `10` abbildet (`internal/adapter/driving/cli/generate.go`, Zeile 295; weitere Abbildungen in `init.go` Zeile 362, `egress.go` Zeile 239, `cli.go` Zeile 479). Der Diagnose-Code dort ist [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) und der Doc-Kommentar des Sentinels (`internal/hexagon/domain/sandbox.go`, Zeilen 11 bis 21) nennt Sandbox-Werte; ein Feld für beide Profile passt dazu nicht ohne Entscheidung (siehe §6).
- **Der Wert landet ungeschützt in JSON.** Die Vorlagen laufen mit `text/template` (`templates.go`, Import; `devcontainer.json.tmpl`, Zeile 15: `"{{$a}}"`), ohne Escaping; ein Wert mit Anführungszeichen oder Rückstrich bräche die JSON-Datei. Das Vorbild der Zeichenmenge ist `cloneURLSafeRE` (`devcontainer_sandbox.go`, Zeile 74).
- **Zusage heute:** [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) verlangt, dass das Profil kein `--privileged` und keine zusätzlichen Capabilities setzt, außer sie sind durch [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) oder [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) verlangt oder von einem eingebauten Feature mitgebracht, und dass jede Sicherheitslockerung in der Befehlsausgabe ausgewiesen wird (Zeilen 458 und 460). `--memory` und `--cpus` verlangen keine Lockerung; ob die Ausweisung sie trotzdem fassen soll, ist nicht festgelegt (Punkt 3 in §6).

**Wege** (Entscheidung des Projektinhabers, 2026-11-01: **A gewählt**; B und C sind nicht gewählt und bleiben als Begründung stehen):

- **Gewählt: A (Entscheidung Projektinhaber 2026-11-01), zwei Felder `devcontainer.resources.memory` und `devcontainer.resources.cpus` mit Validierung.** *Nutzen:* Der häufige Bedarf (zwei Höchstwerte) ist gedeckt, in beiden Profilen, mit einer Prüfung beim Schreiben (`config set`) und beim Lesen (`generate`, `doctor`); die Zeichenmenge des Werts ist geschlossen, das Rendern in JSON dadurch sicher. Das Sandbox-Profil bleibt der einzige Ort der Lockerungen: ein Limit ist keine. *Aufwand:* mittel. Zwei Felder in `ubootYAMLDevcontainer`, zwei Pfad-Arten in `internal/hexagon/domain/configpath.go` (Muster: `ConfigDevcontainerUserUID`, Zeile 77, und `ConfigDevcontainerSandboxRepository`, Zeile 96) samt Tabellen in `devcontainer_sandbox_config.go` (Zeilen 19, 35, 67 bis 71, 106 bis 110), Wertprüfung, Rendern in beiden Zweigen von `devcontainerTemplateData`, Lastenheft-CR, Spec, Handbuch. *Haken:* Der öffentliche Vertrag von `u-boot.yaml` wächst (Lastenheft-CR); Grammatik der Werte der Engine ist nicht gemessen (Punkt 5 in §6); weitere Grenzen (`--memory-swap`, `--pids-limit`, `--ulimit`) sind damit nicht abgedeckt und brauchen einen eigenen Vorgang.
- **B, nicht gewählt (Grund: öffnet `--privileged`, `--cap-add` und weitere Optionen am Sandbox-Profil vorbei und braucht eine Einordnung freier Optionen, für die kein Bedarf über Arbeitsspeicher und CPU hinaus vorliegt), eine allgemeine Liste `devcontainer.runArgs` in `u-boot.yaml`.** *Nutzen:* Sie deckt jede weitere `docker run`-Option ohne neuen Slice. *Aufwand:* hoch. Die Liste öffnet `--privileged`, `--cap-add`, `--device`, `--security-opt` und `--mount` am Sandbox-Profil vorbei: [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) verlangt die Ausweisung jeder Lockerung, `sandboxRelaxationWarnings` (`devcontainer_sandbox.go`, Zeile 308) erzeugt sie bisher aus festen Quellen (`podmanRelaxations`, `dockerRelaxation`, `egressRelaxation`, `featureRelaxationWarnings`), nicht aus freien Zeichenketten. Ein Weg B müsste jede Option einordnen (zulassen, ausweisen, abweisen) oder die Ausweisung auf eine Liste ausdehnen; im Standardprofil entstünde ein zweiter Weg an der Prüfung von Features vorbei (`validateDevcontainerFeatures`, `devcontainer_features.go`, Zeile 703, prüft Quellen, nicht `runArgs`). *Haken:* größerer Lastenheft-CR, eine Klassifikation freier Optionen, die der Slice nicht gemessen hat; ein Test je Option ist nicht abgeschlossen.
- **C, nicht gewählt (Grund: die Grenze bliebe außerhalb von `u-boot` und hielte nicht über `up --recreate` hinaus), nur Handbuch: `docker update` nach dem Start, nicht über `u-boot`.** *Nutzen:* kein Produktcode, kein Lastenheft-CR. *Aufwand:* gering (ein Handbuch-Absatz mit Messung). *Haken:* Der Befehl gilt, soweit angenommen, nur für den laufenden Container und nicht für den nächsten `u-boot devcontainer up --recreate` ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Zeile 501); ob er mit `--privileged` und dem verschachtelten Docker wirkt, ist nicht gemessen, und er läuft am Aufruf von `u-boot` vorbei. Das Problem bleibt für jeden, der die Grenze im Projekt festhalten will.

**Begründung des Planners für A (vom Projektinhaber übernommen):** Weg A ist der kleinste Weg, der das Problem für beide Profile schließt, die Werte lassen sich auf eine geschlossene Zeichenmenge prüfen, und er lässt die Ausweisung der Lockerungen unberührt. Weg C bleibt der Rückfall, wenn die Messung zeigt, dass die Grenze im Sandbox-Profil nicht trägt (siehe die offenen Messungen unten); der Wechsel dorthin wäre eine neue Entscheidung des Projektinhabers (Rückführung in §4).

**Verhältnis zu `--instance`** (Annahme, nicht gemessen; Messung 5 unten): Die Felder stehen in `u-boot.yaml` und wirken über die erzeugte `.devcontainer/devcontainer.json` des Projektordners; jede Instanz desselben Ordners ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Zeilen 497 ff.) liefe damit mit derselben Grenze **je Container**; eine Grenze je Instanz oder je Projekt gibt es nicht, und die Summe über Instanzen ist nicht die Grenze des Hosts (Risiko 7 in §6).

**Werteformat** (Vorschlag zur Freigabe, Grammatik der Engine nicht gemessen; Risiko 5 in §6): `memory` eine positive ganze Zahl mit optionaler Einheit aus `b`, `k`, `m`, `g` (Beispiel `3g`, `512m`); `cpus` eine positive Dezimalzahl mit Punkt (Beispiel `1.5`, `2`). Beide als Zeichenkette in `u-boot.yaml`; geschlossene Zeichenmenge `0-9`, `.` und die vier Einheiten, kein Leerzeichen, kein führendes `-`; Null und leere Zeichenkette sind ungültig (leer = nicht gesetzt). Ob Großbuchstaben, `t` und `--memory=0` zugelassen werden, entscheidet die Messung der Engine (Messung 6).

**Offene Messungen vor dem Code** (nicht gemessen, nicht behauptet; die erste Aufgabe des Slice; sie legen fest, wie weit die Zusage von Weg A reicht, und nur ein Befund „trägt im Sandbox-Profil nicht“ führt zurück zur Entscheidung des Projektinhabers):

1. Greifen `--memory` und `--cpus` im Sandbox-Profil (`--privileged`, `nestedRuntime: docker`) und im Standardprofil, gemessen im Container, den `u-boot devcontainer up` aus der erzeugten Konfiguration startet?
2. Wirken die Grenzen auf den verschachtelten Docker und seine Container (Vererbung der cgroup)?
3. Verhalten auf Docker-Engines außer Linux `amd64` mit lokaler Engine: nicht zugesagt ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)); die Messung nennt, was sie nicht abdeckt.
4. Was passiert bei einem zu kleinen Limit — Abbruch eines Startscripts durch den OOM-Killer (`u-boot-docker-init` und `u-boot-egress-init` stehen in `postStartCommand`, `devcontainer_sandbox.go`, Zeilen 119 und 130) und welcher Exit-Code und welche Ausweisung folgen ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Zeile 506: Exit-Code `11` bei einem Startscript, das wegen eines Umgebungsproblems abbricht)?
5. Zusammenspiel mit mehreren Instanzen ([`slice-v1-devcontainer-instanz-label`](../done/slice-v1-devcontainer-instanz-label.md) §6, Risiko *Ressourcen*; dort entschieden: Weg A, jede Instanz soll isoliert laufen): Gilt die Grenze je Container und was folgt für die Summe?
6. Welche Werte nimmt die Engine an (Einheiten `b`, `k`, `m`, `g`, Großbuchstaben, Dezimalstellen bei `--cpus`, Untergrenze von `--memory`, `0`) und wie lautet ihre Fehlermeldung bei einem abgewiesenen Wert; daraus folgt die Grammatik im Abschnitt *Werteformat*.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **GPU-Zuteilung.** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) sagt GPU-Beschleunigung ausdrücklich nicht zu (Zeile 412); eine Grenze dafür wäre eine neue Zusage ohne Gegenstand im Produkt.
- **Plattenquoten und Größe der Volumes.** Der Docker-Speicher des Sandbox-Profils liegt im Volume `<name>-docker` (`devcontainer_sandbox.go`, Zeile 117); eine Quote braucht eine Option des Speichertreibers der Engine, nicht `runArgs`, und ist ein anderer Vorgang.
- **Grenzen des Compose-Stacks** ([`LH-FA-UP-001`](../../../../spec/lastenheft.md#lh-fa-up-001--umgebung-starten)). Der Stack ist eine eigene Umgebung mit eigener Konfiguration; `u-boot devcontainer up` startet ihn nicht ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Zeile 508).
- **Grenzen für den Hilfs-Container von `u-boot devcontainer up`.** Er ist kurzlebig und gehört nicht zum erzeugten Ergebnis ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Zeile 504); er baut und startet nur.
- **Grenzen des Hosts oder der VM** (Colima, Docker Desktop). Sie liegen außerhalb von `u-boot` und der Engine, die `u-boot` anspricht; die Wirkung einer Grenze im Container kann dort nicht größer sein als die VM.
- **Weitere `docker run`-Optionen** (`--memory-swap`, `--pids-limit`, `--ulimit`, `--shm-size`). Sie sind Gegenstand von Weg B; der Slice liefert sie nur, wenn der Projektinhaber B wählt.
- **Eine Änderung bei einem laufenden Container.** Die Grenze greift, soweit angenommen, erst bei einem neuen Container (`up --recreate`); ob `generate devcontainer` darauf hinweist (wie bei der Egress-Liste, `harness/replay.md` Zeile 55), ist Punkt 4 in §6 und nicht gemessen.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Messung und Freigabe des Wortlauts** ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)): Die sechs offenen Messungen aus §1 sind gefahren (Linux `amd64`, lokale Engine, mindestens zwei Läufe je Zelle, datiert, mit Umgebung; beide Profile, mit und ohne `nestedRuntime: docker`; ein zu kleines Limit; zwei Instanzen; Werte der Engine). Der Weg ist entschieden (A, 2026-11-01); der Projektinhaber hat den Wortlaut des Lastenheft-CR freigegeben (Alt/Neu-Vorschlag in §6, Frage 3). *Beleg:* Messprotokoll unter `docs/messungen/slice-v1-devcontainer-ressourcenlimits/` (die Messung der Hand-Ergänzung liegt dort, die Skripte unter `tools/messung/devcontainer-ressourcenlimits/`; [`harness/conventions.md`](../../../../harness/conventions.md#ablage-von-messprotokollen-und-skripten)); keine Zugangsdaten in den Rohdaten.
- [ ] **Weg A:** `devcontainer.resources.memory` und `.cpus` sind in `u-boot.yaml` lesbar und über `u-boot config set` schreibbar, ungültige Werte enden mit Exit-Code `10` ([`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), Klasse „Struktur-, Namens- und Konfigurationsvalidierung“), `u-boot generate devcontainer` rendert `--memory=<Wert>` und `--cpus=<Wert>` in `runArgs` beider Profile, ohne Angabe ist die erzeugte Datei byte-identisch zu heute (Bestandstest), die Grenzen erzeugen keinen Eintrag der Sicherheitslockerungen (oder die Entscheidung zu Punkt 3 in §6 ist umgesetzt); Tests je Bedingung und je Mutation; ein Golden Case in `harness/replay.md`; Reihenfolge von `--memory` und `--cpus` in `runArgs` relativ zu den Optionen von `nestedRuntime`, Egress und `ollama-host` ist festgelegt und von einem Test gepinnt; `doctor` meldet einen ungültigen Wert über den bestehenden Befund zu `u-boot.yaml` (Fehler der Schema-Prüfung), ohne neuen Check.
- [ ] **Docker-Test und Zusage-Menge:** Ein Test mit Build-Tag `docker` belegt, dass das gesetzte Limit im gestarteten Container ankommt (nach dem Ergebnis der Messung im Sandbox-Profil mit `nestedRuntime: docker` und im Standardprofil); Spezifikation, Lastenheft (CR), Handbuch und `CHANGELOG.md` nennen nur das Gemessene (Linux-Host `amd64`, lokale Engine) mit Beleg-Tabelle (Bedingung → Test oder Streichung, [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)); Handbuch und READMEs ohne Kennungen ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6).
- [ ] `make gates` grün (im Modus der Umgebung).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update für die öffentliche Schnittstelle (Schlüssel in `u-boot.yaml`, Handbuch, README, `CHANGELOG.md`), soweit ein öffentlicher Vertrag berührt ist.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

Gelesen am 2026-10-31, die Stellen von Weg A am 2026-11-01 am Stand `HEAD` nachgelesen (Zeilen unten sind die vom 2026-11-01; sie verschieben sich, der Implementer liest sie am Ist-Stand nach, [`AGENTS.md`](../../../../AGENTS.md) §3.10). Eine Zeile „prüfen“ nennt eine Datei, die der Planner nicht oder nicht an der berührten Stelle gelesen hat. Die Tabelle gilt für Weg A. Die Datei `internal/hexagon/application/config.go` ändert kein anderer Slice mehr: [`slice-v1-konfigdatei-name-als-konstante`](../done/slice-v1-konfigdatei-name-als-konstante.md) liegt in `done/`, `in-progress/` trägt keinen Slice.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/messung/devcontainer-ressourcenlimits/messen.sh`, `docs/messungen/slice-v1-devcontainer-ressourcenlimits/` | neu (mit diesem Entwurf angelegt), Messungen folgen | Skript (bereinigt, nicht gelaufen) und Protokoll der Messung der Hand-Ergänzung; die fünf offenen Messungen aus §1 kommen als weitere Skripte und Rohdaten dazu |
| `docs/messungen/README.md` | update | Zeile in der Tabelle *Bestand* |
| `internal/hexagon/application/initproject.go` | update | `ubootYAMLDevcontainer` (Zeilen 69 bis 79) bekommt ein Feld `Resources` mit einem neuen Typ (Muster: `ubootYAMLDevcontainerUser`, Zeilen 84 bis 86, mit `omitempty` und Zeigern, damit „nicht gesetzt“ unterscheidbar bleibt) |
| `internal/hexagon/domain/configpath.go`, `configpath_test.go` | update | zwei neue Pfad-Arten (Muster: `ConfigDevcontainerUserUID`, Zeile 77, und `ConfigDevcontainerSandboxRepository`, Zeile 96) samt Parser und `String()` (Zeile 265, Round-Trip-Eigenschaft); Test: Happy/Boundary/Negative |
| `internal/hexagon/domain/` (neue Datei für die Wertprüfung, Name beim Bau) | neu | Wertobjekte für Speichergrenze und CPU-Anzahl mit geschlossener Zeichenmenge (Vorbild `NewContainerUID`, `NewNestedRuntime` in `sandbox.go`); der Sentinel bei Fehlern ist eine Entscheidung (Punkt 1 in §6) |
| `internal/hexagon/domain/sandbox.go` | update, falls der Sentinel `ErrInvalidSandboxSetting` mitgenutzt wird | Doc-Kommentar (Zeilen 11 bis 21) nennt heute nur Sandbox-Werte |
| `internal/hexagon/application/devcontainer_sandbox_config.go`, `config_sandbox_test.go` | update | gelesen: `sandboxCoercers` (Zeile 16, ein Eintrag je Pfad-Art), `sandboxYAMLPaths` (Zeile 65) und die Lesetabelle `values` in `sandboxConfigValue` (Zeile 105) bekommen die zwei Pfade |
| `internal/hexagon/application/config.go` | prüfen (Aufrufstellen gelesen) | gelesen: `coerceConfigValue` (Zeile 386) ruft zuerst `coerceSandboxConfigValue`; `validateDevcontainer` ruft `config.go` in Zeile 690 und 757. Ob die zwei neuen Pfade dort ohne Änderung durchlaufen, belegt der Test beim Bau |
| `internal/hexagon/application/devcontainer_sandbox.go` | update | `validateDevcontainer` (Zeile 19) ruft die Prüfung der Felder, `devcontainerTemplateData` (Zeilen 89 bis 136) hängt die Optionen in **beiden** Zweigen an `RunArgs` (Standardprofil: Zeile 97, vor dem `return`; Sandbox: nach Zeile 132), der Doc-Kommentar (Zeilen 81 bis 88) nennt die Quellen von `RunArgs` |
| `internal/hexagon/application/templates.go` | update (nur Kommentar) | gelesen: Doc-Kommentar zu `RunArgs` (Zeilen 88 bis 90: „options of the sandbox profile, then the `ollama-host` option in either profile“) nennt die Quellen; er folgt dem Diff ([`AGENTS.md`](../../../../AGENTS.md) §3.9 (b)) |
| `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl` | **unverändert** | gelesen: Zeilen 14 bis 15 (Sandbox) und 39 bis 40 (Standardprofil) rendern jede Zeichenkette in `.RunArgs` als `"{{$a}}"` in einer Zeile; zwei weitere Einträge brauchen keine Vorlagenänderung |
| `internal/hexagon/application/doctor.go`, `doctor_sandbox.go` | `doctor.go` unverändert (Verhalten), `doctor_sandbox.go` nicht gelesen | gelesen: `checkUbootYaml` ruft `validateDevcontainer(cfg.Devcontainer)` (`doctor.go`, Zeile 617) und meldet einen Fehler der Schema-Prüfung als Fehler-Befund „u-boot.yaml devcontainer schema invalid“; ein ungültiger Wert der neuen Felder fällt damit ohne neuen Check in diesen Befund, sobald `validateDevcontainerSandbox` ihn prüft. Der Hinweis dort (Zeile 622) nennt „URL scheme, host, or feature name“ und wird mit dem Slice [`slice-v1-hilfetexte-ohne-kennungen`](../open/slice-v1-hilfetexte-ohne-kennungen.md) berührt (gleiche Zeile, Reihenfolge klären, Frage 6) |
| `internal/adapter/driving/cli/generate.go`, `init.go`, `egress.go`, `cli.go` | prüfen (Abbildungen gelesen) | gelesen: `domain.ErrInvalidSandboxSetting` wird auf Exit-Code `10` abgebildet in `generate.go` Zeile 295 (Diagnose-Code [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)), `init.go` Zeile 362, `egress.go` Zeile 239 und `cli.go` Zeile 479; bei Feldern für beide Profile ist der Code zu klären (Frage 1). Die Liste „Writable paths“ im Hilfetext von `config` (`config.go`, Zeile 168) nennt nur `project.name` und `devcontainer.enabled`; die neuen Pfade stehen dort nicht, ein Eingriff ist nicht nötig |
| `internal/hexagon/application/generate_resources_test.go` (Name beim Bau) | neu | Happy (beide Profile), Boundary (Zeichenmenge, Einheiten), Negative (Exit `10`), Reihenfolge mit `ollama-host` und `nestedRuntime`, byte-identisch ohne Angabe; Vorbild: `generate_ollama_host_test.go` (Hilfsfunktion `runArgsOf`, Zeile 65, Konstante `hostLineJSON`, Zeile 25) |
| `internal/hexagon/application/generate_sandbox_test.go`, `generate_ollama_host_test.go`, `generate_host_socket_test.go` | prüfen | Erwartungswerte von `runArgs` (`generate_host_socket_test.go`, Zeilen 129 bis 159 erwartet `[--privileged]`) bleiben ohne Angabe unverändert; bei Änderung gehört sie in den Diff |
| `internal/e2e/` (neue Datei mit Build-Tag `docker`, Name beim Bau) | neu | Docker-Test, der ein Limit im gestarteten Container belegt (Vorbild: `devcontainer_start_docker_test.go`, `TestDevcontainerUp_Docker_S1_S9_StandardProfile`, Zeile 383, und `sandbox_devcontainer_docker_test.go`; beide nicht gelesen) |
| `harness/replay.md` | update | neue Zeile in der Tabelle der Golden Cases (Zeilen 49 bis 66 führen die Fälle von `generate devcontainer`) |
| `spec/lastenheft.md` | **CR vor dem Code** (Wortlaut freigegeben vom Projektinhaber; nicht in diesem Entwurf geschrieben) | [`LH-FA-CONF-002`](../../../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) (Zeile 832 am 2026-11-01) zählt die optionalen Schlüssel einzeln auf; ein neuer Schlüssel im öffentlichen Vertrag von `u-boot.yaml` ändert die Zusage. Wortlaut, Ort (Ergänzung dort, in [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) oder eine neue Anforderung) und Zeile im Änderungsverlauf (Tabelle, Zeilen 1579 bis 1607) entscheidet der Projektinhaber |
| `spec/spezifikation.md` | update (nach dem Code) | gelesen: [`SPEC-005`](../../../../spec/spezifikation.md#spec-005--schema-der-projektkonfiguration-u-bootyaml) (Zeile 1119: Beispiel-Schema mit auskommentierten optionalen Feldern, `devcontainer.user.uid` bis `egress.allow`) und [`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config) (Zeile 509: Tabelle „Erlaubte Pfade“ mit Spalten Pfad, Lesen, Schreiben, Werte; zwei neue Zeilen); Regeln für das Rendern und die Prüfung als Spec-Absatz mit Struktur-ID nach dem Schema in `harness/conventions.md` (nicht ad hoc) |
| `spec/architecture.md` | prüfen | ob die Konfigurationsschicht dort beschrieben ist; nicht gelesen. Der Arbeitsbaum ist am 2026-11-01 sauber (kein anderer Slice ändert die Datei) |
| `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md` | update | gelesen (Treffer-Zeilen): Abschnitt 4.9 „Einen Devcontainer einrichten“ (Zeile 701, mit `u-boot config set devcontainer.user.uid 501` in Zeile 727), 4.10 (Zeile 1010), 4.11 (Zeile 1262) und die Tabelle der Schlüssel von `u-boot.yaml` (Zeile 2454: `devcontainer.user.uid` · Bedeutung · Werte · schreibbar). Zwei neue Zeilen dort und ein Absatz in 4.9; ohne Kennungen ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6); Wortlaut nur über das Gemessene (Messung 1 bis 6) |
| `docs/maintainer/cli-json-contract.md` | unverändert | gelesen (Suche nach `devcontainer.user`, `devcontainer.sandbox`): kein Treffer; die neuen Felder berühren den JSON-Vertrag nicht, die Diagnose-Codes folgen Frage 1 |
| `docs/plan/adr/` und `docs/plan/adr/README.md` | neu (Frage 7: ob eine ADR nötig ist) | eine ADR für Weg A und die Zusage-Menge (Vorbild: [ADR-0037](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md)), mit Eintrag im ADR-Index; die Wahl von A ist eine Entscheidung des Projektinhabers, die Frage ist, ob sie als ADR oder im CR festgehalten wird |
| `internal/adapter/driven/docker/devcontainer_launcher.go` | unverändert | gelesen: `runArgs` (Zeile 292) ist die Argumentliste des **Hilfs-Containers** von `u-boot devcontainer up`, nicht das `runArgs` der `devcontainer.json`; die Grenze wirkt über die Dev-Container-CLI aus der erzeugten Datei, nicht über diese Funktion (Hilfs-Container ausgeschlossen, §1) |
| diese Datei, `docs/plan/planning/open/README.md` | update | Zeile dieses Slice; Moves per `make slice-mv` |

**Wächter und rot färbende Mutationen** (je neuem Wächter; der Implementer führt sie aus und hält sie fest, Weg A):

| Wächter | Mutation, die ihn rot färbt |
|---|---|
| Test „Standardprofil rendert `--memory` und `--cpus`“ | in `devcontainerTemplateData` das `return` vor dem Anhang der Optionen (Zeile 97) setzen |
| Test „Sandbox-Profil rendert beide“ | den Anhang nach Zeile 132 entfernen |
| Bestandstest „ohne Angabe byte-identisch“ (Standard- und Sandbox-Profil, auch mit `ollama-host` und `nestedRuntime`) | immer `--memory=` und `--cpus=` anhängen, auch bei leerem Wert |
| Test „Zeichenmenge“ (Anführungszeichen, Rückstrich, Leerzeichen, `;`, Zeilenwechsel, führendes `-`) | Zeichenmenge um Anführungszeichen erweitern |
| Test „ungültiger Wert endet mit Exit-Code `10`“ (`generate` und `config set`) | die Abbildung des Fehlers auf Exit-Code `10` in der CLI entfernen |
| Test „Grenzen sind keine Sicherheitslockerung“ (Entscheidung zu Punkt 3 in §6) | die Grenze als `SecurityRelaxation`-Eintrag melden |
| Docker-Test „Limit im Container“ | die Option beim Rendern auslassen; der Container läuft dann ohne Grenze, die Abfrage der Grenzen am gestarteten Container (Feldnamen beim Bau der Probe messen) schlägt fehl |

**Abgleich mit `git diff --stat` (Stand nach dem Commit der Doku, Implementer):** Abweichungen von der Tabelle: `internal/hexagon/application/doctor_features_test.go` kommt hinzu (ein Test pinnt den Befund zu `devcontainer.resources.memory`, `doctor.go` und `doctor_sandbox.go` sind unverändert); `internal/adapter/driving/cli/generate.go`, `init.go` und `egress.go` sind unverändert, weil die Fehler dort als `%v` in `ErrGenerateManualConflict` oder `ErrConfigSchemaInvalid` eingewickelt sind und Exit-Code `10` schon über diese Sentinels entsteht (nur `cli.go` nimmt `ErrInvalidResourceLimit` in `isConfigValidationError` auf); `README.md` und `README.de.md` sind unverändert (sie nennen keine Schlüssel von `u-boot.yaml`; gesucht, nicht Zeile für Zeile gelesen); keine ADR (Entscheidung des Projektinhabers); `spec/architecture.md` ist unverändert und nicht gelesen. Die Messskripte `messen-engine.sh` und das Protokoll `roh-engine.txt` sind neu.

Die Wege B und C sind nicht gewählt; ihre Plan-Zeilen (Liste `[]string` mit Einordnung je Option, Handbuch-Absatz) entfallen. Käme es zu einem Wechsel, ist der Plan neu zu schreiben (§4).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `next` → `in-progress`): Start in einer eigenen Devcontainer-Instanz (`u-boot devcontainer up --instance`), Projektinhaber 2026-11-01; kein Termin. Dazu: der Projektinhaber hat den Slice einem Release zugeordnet (Frage 8); `Verantwortlich:` gesetzt; der Weg ist entschieden (A, 2026-11-01); der Wortlaut des Lastenheft-CR ist freigegeben, bevor Code entsteht (Frage 3); kein anderer Slice in `in-progress/` (WIP-Limit 1; am 2026-11-01 liegt dort keiner). Die Messungen aus §1 sind der erste Schritt und legen die Reichweite der Zusage fest.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Messung zeigt, dass die Grenze im Sandbox-Profil oder auf den verschachtelten Docker nicht trägt und ein eigener Bau nötig wird (neue Entscheidung des Projektinhabers: Zusage nur Standardprofil, oder Weg C); oder der Lastenheft-CR wird ein eigener Slice vor dem Code (Frage 3).
- `in-progress` → `open` (blockiert): Der Projektinhaber entscheidet den Weg nicht, oder die Messung kann keinen Container mit Grenze herstellen (Engine nicht verfügbar).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; `make gates` grün; die Zusage in Lastenheft, Spezifikation und Handbuch nennt nur das Gemessene (Linux-Host `amd64`, lokale Engine, die gemessenen Profile); Review-Report und Verification-Report liegen vor; die Closure-Notiz mit Steering-Loop-Lerneintrag ist geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Entscheidungen des Projektinhabers vom 2026-11-01** (Weg A; Punkte 1, 3, 4, 5, 7 und 8 entschieden, 2 und 6 offen):

- **Zu 1:** eigener Sentinel mit eigenem Diagnose-Code, Exit-Code `10`; `ErrInvalidSandboxSetting` bleibt auf Sandbox-Werte beschränkt.
- **Zu 3:** der Lastenheft-CR entsteht in diesem Slice als erster Commit vor dem Code; der Alt/Neu-Vorschlag unten ist freigegeben, das Werteformat folgt nach Messung 6.
- **Zu 4 und 5:** keine Ausweisung als Sicherheitslockerung (ein Test pinnt das), kein Hinweis bei bestehendem Container.
- **Zu 7 und 8:** keine ADR; Release: nächster Minor, `CHANGELOG.md` *Added*.
- **Zu 6:** steht mit dem Slice [`slice-v1-hilfetexte-ohne-kennungen`](../open/slice-v1-hilfetexte-ohne-kennungen.md) unverändert offen (kein Widerspruch, gleiche Zeilen).

**Verbleibende Fragen und ursprüngliche Fassung:**

1. **Sentinel und Diagnose-Code** der neuen Felder (Risiko 1 unten): `ErrInvalidSandboxSetting` mit Code `LH-FA-DEV-006` mitnutzen (Doc-Kommentar in `internal/hexagon/domain/sandbox.go`, Zeilen 11 bis 21, um die Felder ergänzen) oder eigener Sentinel mit eigenem Code (Exit-Code beide `10`); der Code ist Vertrag (`cli-json-contract.md`). <!-- d-check:ignore (Code-Wert als Zitat) -->
2. **Werteformat** (Vorschlag im Abschnitt *Werteformat* in §1): Einheiten, Großschreibung, Null, Untergrenze; die Messung 6 liefert die Grundlage, die Freigabe gehört dem Projektinhaber.
3. **Lastenheft-CR: Wortlaut und Gefäß.** Ein Lastenheft-CR ist nötig, bevor Code entsteht ([`LH-FA-CONF-002`](../../../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) zählt die optionalen Schlüssel einzeln auf). Frage: eigener Slice nach dem Muster [`slice-cr-devcontainer-instanz-zusage-getestet`](../done/slice-cr-devcontainer-instanz-zusage-getestet.md) (CR zuerst, Code danach, so wie bei `--instance`) oder im selben Slice. *Alt/Neu-Vorschlag, nur Entwurf, nicht ins Lastenheft geschrieben:*
   - *Fähigkeit, `LH-FA-CONF-002`.* Alt (Ende der Aufzählung, Zeile 832): „… sowie `devcontainer.sandbox.nestedRuntime` (`none` | `podman`), `devcontainer.sandbox.onUnavailable` (`warn` | `fail`), `devcontainer.sandbox.repository` und `devcontainer.sandbox.egress.enabled` / `.allow`.“ Neu: „… und `devcontainer.sandbox.egress.enabled` / `.allow` sowie `devcontainer.resources.memory` (positive ganze Zahl mit optionaler Einheit `b`, `k`, `m` oder `g`) und `devcontainer.resources.cpus` (positive Dezimalzahl); ein ungültiger Wert endet mit Exit-Code `10`.“ (Format vorbehaltlich Messung 6 und Frage 2.) <!-- d-check:ignore (Wortlaut-Entwurf, nicht im Lastenheft) -->
   - *Zusage-Menge, `LH-FA-DEV-001` oder neue Anforderung (Ort offen).* Neu: „Sind `devcontainer.resources.memory` oder `devcontainer.resources.cpus` gesetzt, trägt die erzeugte `.devcontainer/devcontainer.json` in `runArgs` `--memory=<Wert>` beziehungsweise `--cpus=<Wert>`, in beiden Profilen; ohne Angabe bleibt die erzeugte Datei unverändert. Zugesagt ist, dass das Limit im Container ankommt, den `u-boot devcontainer up` aus dieser Datei startet, **nur** auf einem Linux-Host `amd64` mit lokaler Engine und nur in den Profilen und Kombinationen, die ein Test mit Build-Tag `docker` belegt; für andere Hosts und Engines, für die Wirkung auf den verschachtelten Docker und auf die Summe mehrerer Instanzen wird nichts zugesagt, soweit die Messung nichts anderes zeigt.“ <!-- d-check:ignore (Wortlaut-Entwurf, nicht im Lastenheft) -->
   - *Änderungsverlauf:* Zeile der nächsten freien Version (nach `0.26.0`, Zeile 1610).
4. **Ausweisung** (Risiko 3): Die Grenzen sind keine Sicherheitslockerung; ob die Befehlsausgabe sie dennoch nennt, oder genau nicht (Test pinnt die Entscheidung).
5. **Hinweis bei bestehendem Container** (Risiko 4): Ob `generate devcontainer` darauf hinweist, dass die Grenze erst bei einem neuen Container (`up --recreate`) greift.
6. **Reihenfolge mit [`slice-v1-hilfetexte-ohne-kennungen`](../open/slice-v1-hilfetexte-ohne-kennungen.md):** beide berühren Meldungen in `devcontainer_sandbox.go`/`doctor.go`; der später gestartete Slice liest die Texte des früheren am Ist-Stand (kein inhaltlicher Widerspruch, nur dieselben Zeilen).
7. **ADR:** ob Weg A und die Zusage-Menge eine ADR brauchen oder im CR stehen (Plan-Zeile `docs/plan/adr/`).
8. **Release-Zuordnung** (Patch oder Minor; Wortlaut im `CHANGELOG.md`: *Added*).

**Grenzen der Zusage** (Festlegung dieses Entwurfs, vom Projektinhaber zu bestätigen): Zugesagt wird nur, was ein Test pinnt. Für die Aussagen „beide Profile“, „ohne Angabe byte-identisch“ und „das Limit kommt im Container an“ gilt: Unit-Tests pinnen das Rendern; ein **Docker-Test hinter dem Build-Tag `docker`** pinnt das Ankommen im gestarteten Container; der Nachweis dafür ist der CI-Lauf des Jobs `integration-docker` (der volle Lauf von `make test-docker` ist Sache der CI, Entwurf [`slice-v1-test-docker-in-ci`](../open/slice-v1-test-docker-in-ci.md)); Host: Linux `amd64`, lokale Engine. Eine Bedingung in Spec oder Handbuch ohne pinnenden Test wird gestrichen ([`AGENTS.md`](../../../../AGENTS.md) §3.10, letzter Punkt).

**Risiken:**

1. **Sentinel und Diagnose-Code eines Felds für beide Profile:** `domain.ErrInvalidSandboxSetting` ist für Sandbox-Werte dokumentiert und die CLI meldet ihn unter dem Code [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (`generate.go`, Zeile 295); ein Feld, das im Standardprofil wirkt, braucht entweder einen anderen Code oder einen eigenen Sentinel, beide mit Exit-Code `10`. — **Ausgang:** bei Closure.
2. **Wirkung im Sandbox-Profil und im verschachtelten Docker:** `--privileged` und der eigene Daemon im Container sind der Grund für die Messung; trägt die Grenze dort nicht oder nur für den äußeren Container, ist die Zusage entsprechend klein (nur Standardprofil, oder nur der äußere Container), und die Dokumentation sagt es. — **Ausgang:** bei Closure.
3. **Ausweisung:** `--memory` und `--cpus` sind keine Sicherheitslockerung nach [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil); ob die Befehlsausgabe sie trotzdem nennt, entscheidet der Projektinhaber im CR. Ein Test pinnt die Entscheidung (kein `SecurityRelaxation`-Eintrag durch die Grenzen, oder genau einer). — **Ausgang:** bei Closure.
4. **Bestehender Container:** Die Grenze greift erst bei einem neuen Container (angenommen, nicht gemessen); ein erneutes `generate devcontainer` ändert `devcontainer.json`, aber nicht den laufenden Container. Ob und wie der Befehl darauf hinweist (wie bei der Egress-Liste), ist offen. — **Ausgang:** bei Closure.
5. **Grammatik der Werte:** Die Zeichenmenge ist geschlossen zu halten (`text/template` ohne Escaping); was die Engine als Wert annimmt (Einheiten, Dezimalstellen, Untergrenze von `--memory`), ist nicht gemessen, und Docker, Podman und andere Engines können sich unterscheiden. Eine Prüfung, die enger ist als die Engine, ist besser als eine, die einen Wert durchlässt, den die Engine beim Start abweist. — **Ausgang:** bei Closure.
6. **Zu kleines Limit:** Ein Startscript (`u-boot-docker-init`, `u-boot-egress-init`) oder der Dienst im Container kann vom OOM-Killer beendet werden; der Exit-Code von `up` ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten): `11`, wenn ein Startscript wegen eines Umgebungsproblems abbricht, `12` sonst) und die Meldung dazu sind nicht gemessen. — **Ausgang:** bei Closure.
7. **Mehrere Instanzen** (Weg A: dieselbe Grenze je Instanz, aus `u-boot.yaml`; Annahme, Messung 5): Eine Grenze gilt je Container, nicht je Projekt; bei zwei oder mehr Instanzen ([`slice-v1-devcontainer-instanz-label`](../done/slice-v1-devcontainer-instanz-label.md), dort Weg A entschieden) ist die Summe der Grenzen nicht die Grenze des Hosts. Beide Slices berühren sich nach der Plan-Tabelle jenes Entwurfs in keiner Quelldatei des Produkts (er ändert den Launcher, die CLI-Unterbefehle und `FindContainer`, dieser die Konfiguration und die Erzeugung); gemeinsam sind Lastenheft, Spezifikation, Handbuch, `CHANGELOG.md` und die Indizes, deren Reihenfolge der Merge klärt. Die Reihenfolge der Release-Zuordnung entscheidet der Projektinhaber. — **Ausgang:** bei Closure.
8. **Zusage weiter als die Messung** ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)): Die Messung läuft auf einem Linux-Host `amd64` mit lokaler Engine; Docker Desktop, Colima, rootless Docker, Podman als Host-Engine, entfernte Engines, Linux-`arm64` und Windows sind nicht zugesagt ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)), und die Grenze der VM unter Colima oder Docker Desktop begrenzt den Container zusätzlich. — **Ausgang:** bei Closure.
9. **Bedingung in Spec oder Handbuch ohne Test** ([`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)): „in beiden Profilen“, „ohne Angabe unverändert“ und „nur auf dem gemessenen Host“ sind Aussagen über den Code und brauchen je einen pinnenden Test oder entfallen ([`AGENTS.md`](../../../../AGENTS.md) §3.10, letzter Punkt). — **Ausgang:** bei Closure.
10. **Bestehende Golden Cases:** Die Zeilen von `harness/replay.md` für `ollama-host`, `nestedRuntime` und Egress pinnen `runArgs`; ohne Angabe der neuen Felder dürfen sie sich nicht ändern. — **Ausgang:** bei Closure.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt: Was hat funktioniert · Was ging anders als geplant · Steering-Loop-Eintrag · Beobachtungs-Register · Folge-Slices · Risiken aus §6 · Drei Paarungen.)*

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

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden je nach Weg `internal/hexagon/domain` (DOM, Pfad-Arten und Wertprüfung), `internal/hexagon/application` (APP, Schema, Rendern, Prüfung), `internal/adapter/driving/cli` (CLI, nur Abbildung des Fehlers, prüfen), `internal/e2e` (E2E, Docker-Test), `tools/` (WZG, Messskript) sowie `spec/`, `harness/`, `docs/plan/` und `docs/messungen/` (DOK). Alle sind in der Modus-Deklaration von [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) als Greenfield auditiert; eine zu grobe Sub-Area ist nicht dabei.

**Vorgelagert — offene Beobachtungen sichten:** Das Register ([`observations/`](../observations/)) ist nach den Titeln durchgesehen; von den passenden ist der Anfang der `observation.md` gelesen, der Zähler ist die Zahl der Dateien unter `evidence/` am 2026-10-31. [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (DOK, 28 Dateien, über der Schwelle von drei) und [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (DOK, 26) betreffen die Zusage-Menge und die Bedingungen in Spezifikation und Handbuch (Risiken 8 und 9); [`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](../observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md) (DOK, 9) betrifft die Messung der Hand-Ergänzung, die an einer Datei im Wegwerf-Ordner lief und nicht an der Umsetzung des Feldes: die Messung nach dem Bau läuft an der Umsetzung; [`BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf`](../observations/BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf/observation.md) (E2E, 3) betrifft Kombinationen mit `nestedRuntime` ohne Lauftest (der Docker-Test dieses Slice nennt seine Kombination und lässt die übrigen aus). Nur nach Titel gesichtet, Inhalt beim Start zu lesen: [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (61, die Regel steht in [`AGENTS.md`](../../../../AGENTS.md) §3.10).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
