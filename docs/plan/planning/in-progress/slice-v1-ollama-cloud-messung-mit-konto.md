# Slice v1: Cloud-Modelle von Ollama mit Konto — Claude Code im Devcontainer über einen angemeldeten Ollama (Host-Server und Server im Container) messen (Katalog-Feature `ollama`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-ollama-cloud-modelle`](../welle-ollama-cloud-modelle.md) (Release `v0.16.0`, vom Projektinhaber bestätigt am 2026-10-24) — zweiter Slice der Welle: die Messung **mit Konto** (Ausgang A der Welle-Datei §1). Ihm geht [`slice-v1-ollama-cloud-messung`](../done/slice-v1-ollama-cloud-messung.md) voraus (Messung ohne Konto); ihm folgen die ADR, der Lastenheft-CR, der Code, die Spezifikation und der Release-Cut (Welle-Datei §4; Kennungen vorläufig).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Devcontainer-Features; Katalog-Feature `ollama`), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Hosts), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zugangsdaten), [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Zusage-Menge: Linux `amd64`, lokale Engine); [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) (Ollama im Sandbox-Devcontainer; Vorgabe „kein Anschluss an einen Ollama des Hosts“, siehe §1 *Spannung*), [`ADR-0035`](../../adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) (Muster: Zusage nur für das Gemessene, Token nur als Umgebungsvariable).

Anlass: Entscheidung des Projektinhabers vom 2026-10-24 (Rückfrage zu Frage 2 der Welle-Datei): Ausgang A, Messung mit einem angemeldeten Ollama auf dem Host. Der Plan entstand am selben Tag.

**Berührte Spec-Stellen:** [`SPEC-047`](../../../../spec/spezifikation.md#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen) (Feature `ollama`) — der Slice **liest** sie und ändert keine; die Nachführung folgt im Spezifikations-Slice der Welle. Das Lastenheft ändert nur der CR. Der Verweis zeigt **aufwärts**: Lastenheft und Spezifikation nennen diesen Slice nie.

**Verantwortlich:** pt9912 (Rolle Implementer)

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-24.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Mit gemessenen Belegen **mit Konto** steht fest, (a) ob und auf welchem Weg ein Claude Code im Devcontainer einen **angemeldeten Ollama-Server des Hosts** erreicht und mit den Cloud-Modellen `glm-5.3-flash:cloud[1m]` und `gpt-oss:120b-cloud` läuft, mit welchen Hosts, welchen Egress- und Sandbox-Konsequenzen und welcher Abhängigkeit von der Docker-Engine, und (b) was ein **im Container angemeldeter Ollama** an Hosts und Dateien braucht und ob dessen Geräteschlüssel ein Zugangsdatum im Sinne von [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) ist; jeweils im Standardprofil **und** im Sandbox-Profil (Nutzerpfad). Dazu die offenen Messpunkte des ersten Slice: Hosts nach der Anmeldung, Dateien, `OLLAMA_NO_CLOUD` mit Konto, Wirkung von `OLLAMA_API_KEY`, Verhalten von `ollama launch claude` in `restricted` gegen `discover`. Ergebnis: ein Messprotokoll mit Rohdaten und die Eingabe für die Folge-ADR. Zugesagt wird später nur, was gemessen ist, bedingungsgleich: „mit dem Konto des Projektinhabers auf einem Host gemessen“, nicht „Cloud-Modelle laufen“.

### Gefundener Stand (gelesen am 2026-10-24)

- **Ergebnis ohne Konto.** [Messprotokoll](../../../messungen/slice-v1-ollama-cloud-messung/README.md): Ollama `0.40.6` kennt `launch` und das Ziel `claude`; ohne Konto endet `ollama run` und `ollama launch claude` in `discover`/`allow` mit „sign in“-Meldungen, in `restricted` startet `launch` das installierte `claude` (Ursache nicht gemessen); `ollama launch claude` übergibt `--model` und setzt `ANTHROPIC_BASE_URL=http://127.0.0.1:11434` (Handvariante mit Platzhalter); `~/.ollama/id_ed25519` entsteht beim ersten Start des Servers, im Modell-Volume. Die Tabelle *Nicht gemessen* des Protokolls ist die Eingabe dieses Slice.
- **Nutzung nach Angabe des Projektinhabers (2026-10-24):** `export ANTHROPIC_AUTH_TOKEN=ollama; export ANTHROPIC_BASE_URL="http://localhost:11434"; claude --model glm-5.3-flash:cloud[1m]`. Der Wert `ollama` ist ein Platzhalter, kein Geheimnis; Claude Code spricht direkt mit dem Ollama-Server, die Anmeldung zum Cloud-Dienst liegt am Server. Die Bedeutung von `[1m]` im Modellnamen ist **ungeprüft** (Messfall K6).
- **Host-Server (Angabe des Projektinhabers, nicht am Host gelesen):** ein Ollama-Dienst auf `127.0.0.1:11434`, das Ollama-Binary liegt auf dem Host, der Projektinhaber ist dort angemeldet. Der Slice liest und ruft den Host-Ollama **nicht** auf (kein Host-Werkzeug, kein Zugriff auf `~/.ollama`); er erreicht den Server nur über das Netz aus einem Container.
- **Handbuch.** `docs/user/benutzerhandbuch.md`, Abschnitt *Ollama (optional)*: Zeile 2023 „Cloud-Modelle sind nicht enthalten“ (nur ohne Konto gemessen); Zeilen 2043 und 2044: „GPU-Beschleunigung, ein Anschluss an einen Ollama-Dienst des Hosts und die Verwaltung von Modellen durch `u-boot` gibt es nicht“ (Aussage über `u-boot`, nicht über Docker).
- **Spannung zu früheren Vorgaben.** [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) Zeile 19 nennt als Vorgabe des Projektinhabers vom 2026-10-18 „keinen Anschluss an einen Ollama des Hosts“; [`welle-sandbox-egress-ollama`](../done/welle-sandbox-egress-ollama.md) §6 schloss ihn aus (Begründung dort: widerspräche [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil); die Aufzählung der Anforderung (Zeilen 448 bis 456 des Lastenhefts) nennt keinen Netzanschluss an den Host, aber Schadensbegrenzung als Zweck). Die Entscheidung vom 2026-10-24 macht den Host-Server zum **Messgegenstand**; sie ändert weder die ADR (Accepted, unveränderlich) noch das Lastenheft. Ob eine Zusage daraus wird, entscheidet die Folge-ADR (Frage 5 in §6).
- **Egress-Mechanismus.** `egress-init.sh.tmpl` (Zeilen 60 bis 96): dnsmasq mit `no-hosts` und `no-resolv`, je erlaubtem Namen `server=/<name>/<upstream>` und `nftset=`; nftables-Kette `output` mit `policy drop`, erlaubt sind `oifname lo`, aufgebaute Verbindungen und `@allowed` (nur Adressen, die dnsmasq für erlaubte Namen beantwortet hat). Eine Adresse des Docker-Gateways ist damit nur über einen **Namen** erlaubbar; ob `host.docker.internal` (ein Eintrag in `/etc/hosts` bei `--add-host`) durch dnsmasq beantwortet wird, ist Messfall K3. `NewEgressHost` (`internal/hexagon/domain/sandbox.go`, `egressHostPattern`, Zeile 122) akzeptiert Namen aus Ziffern, Buchstaben, `-` und mindestens einem Punkt, also auch `host.docker.internal`; eine IP-Adresse passt syntaktisch ebenfalls (ob ein solcher Eintrag über `server=/…/` und `nftset=` etwas bewirkt, ist Messfall K3).
- **Kein `--add-host` in der Auslieferung.** Die Suche nach `add-host`, `host-gateway` und `host.docker.internal` über `internal/` (Go, Vorlagen, Skripte) hatte keinen Treffer; `RunArgs` (`internal/hexagon/application/templates.go`, Zeile 89) stammt aus den Lockerungen des Sandbox-Profils (`devcontainer_sandbox.go`, `egressRelaxation()` Zeile 161: `--cap-add=NET_ADMIN`). Ein Weg über `host.docker.internal` ist deshalb eine **Handvariante** des Messprojekts (Zusatz in `runArgs` durch das Messskript), ausgewiesen.
- **Zugangsvariablen.** `agentCredentialEnv` (`devcontainer_sandbox.go`, Zeile 441): `ANTHROPIC_API_KEY`, `GH_TOKEN`, `GITLAB_TOKEN`; für `ollama` keine. Das Binary enthält die Namen `OLLAMA_API_KEY`, `OLLAMA_CLOUD_BASE_URL`, `OLLAMA_NO_CLOUD`, `OLLAMA_LAUNCH_DSH_API_KEY` (Protokoll, Abschnitt E).
- **Messmuster.** `tools/messung/ollama-cloud/` (gelesen: `lib.sh` mit `redact`, `umgebung_pruefen`, `bestand`, `praefix_entfernen`; `06-suche.sh`). `umgebung_pruefen` sperrt unter anderem `ANTHROPIC_*`, `OLLAMA*`, `*TOKEN*`, `*API_KEY*`, `*AUTH*`; `lib.sh` startet jedes Skript mit `exec env -i … bash` neu, wobei die Werte der Positivliste auf der Kommandozeile von `env` stehen. `redact` ist formgebunden. `06-suche.sh` ohne Argument trifft 36 Zeilen älterer, abgeschlossener Protokolle (kein Zugangsdatum, Verification-Report des ersten Slice, DoD 3).
- **Bericht des Verifiers (Übergabe).** [Verification-Report](../../../reviews/2026-10-24-slice-v1-ollama-cloud-messung-verify.md), Abschnitt *Übergabe an den Folge-Slice*: sechs Punkte (Schwärzung vor dem ersten Lauf mit Artefakt, Zugang nur als Umgebungsvariable mit benannter Ausnahme, offene Fragen aus den Rohdaten, Messhygiene, Suche mit Ausnahmeliste, Zusage-Menge). Sie sind in §2 und §3 aufgenommen.

**Stand der Umsetzung (Implementer, 2026-10-24):** Weg (a) und Weg (b) sind gemessen ([Messprotokoll](../../../messungen/slice-v1-ollama-cloud-messung-mit-konto/README.md), 40 von 40 Prompts); den Wartepunkt der Anmeldung im Container hat der Projektinhaber am eigenen Terminal mit `tools/messung/ollama-cloud-konto/07-anmeldung.sh` aufgelöst. **Kontozustand (Projektinhaber, 2026-10-24):** bestätigt ohne Eintrag aus der Messung (nur ein Eintrag, nichts von heute). Nach Beobachtung des Projektinhabers (eigene Wiederholung der Anmeldung, nicht Teil der Messung, nicht in den Rohdaten) trägt die Anmeldung den Geräteschlüssel automatisch im Konto ein; nach der Messung stand dort kein Eintrag von heute. Ob `ollama signout` ihn entfernt hat, ist **nicht gemessen** (naheliegende, ungemessene Erklärung; andere nicht ausgeschlossen). Die 28 Anfragen von heute auf der Kontoseite sind eine Beobachtung ohne Schluss (Herkunft nicht ablesbar); siehe Abschnitt *Beobachtung des Projektinhabers* des Messprotokolls. Entscheidungen des Projektinhabers vom 2026-10-24 zu den Fragen 1 bis 5 in §6: höchstens 40 Prompts (20 Weg a, 14 Weg b, je 2 für `OLLAMA_NO_CLOUD` und den Platzhalter von `OLLAMA_API_KEY`, 2 Reserve), Relay erlaubt, Anmeldung im Container durch den Projektinhaber selbst, `OLLAMA_API_KEY` nur als Platzhalter, die Vorgabe „kein Anschluss an einen Ollama des Hosts“ (siehe *Spannung*) gilt für diese Messung nicht; Handbuch und Lastenheft bleiben unverändert.

**Zugangsregeln dieses Slice (Vorgabe des Projektinhabers, 2026-10-24):**

1. Keine Schlüsseldatei des Hosts wird gelesen, kopiert oder in einen Container eingebunden; `~/.ollama/id_ed25519` **nie**. Der Zugang zum angemeldeten Ollama geschieht nur über den laufenden Host-Server (Netz), nie über Dateien.
2. Nichts davon in Dateien, Protokollen oder im Chat: keine Anmelde-URL, kein Code, kein Token, kein Antwort-Inhalt.
3. `redact()` und `umgebung_pruefen()` werden **vor dem ersten Lauf** mit einer Gegenprobe geprüft, die ein Artefakt mit Zeitstempel hinterlässt (Verification-Report V-3); der erste Lauf bricht ab, wenn das Artefakt fehlt.

**Messfälle** (Modelle: `glm-5.3-flash:cloud[1m]` und `gpt-oss:120b-cloud`, Linux `amd64`, lokale Engine, erzeugte Projekte mit `ollama`, `claude-code` und `node`; ein Prompt ist ein Satz, die erbetene Antwort ein Wort; „Lauf“ = Container neu, Volume des Weges neu, falls nicht anders genannt):

*Weg (a) — angemeldeter Ollama-Server des Hosts, Claude Code im Container:*

- **K1 Erreichbarkeit.** Aus dem Container: löst `host.docker.internal` auf (`getent hosts`), erreicht der Container `host.docker.internal:11434` ohne Hilfe (Hypothese: nein, weil der Dienst auf `127.0.0.1` lauscht; zu messen, nicht anzunehmen), und mit einem **Relay** (K1b: ein Hilfs-Container mit `--network host`, der auf der Adresse, die `host-gateway` im Messnetz ergibt, den Port 11434 an `127.0.0.1:11434` des Hosts weiterreicht, nur für die Dauer des Falls). Ausgabe von `ollama --version` im Container bei erreichbarem Server (ob und welche Version des **Servers** sie meldet).
- **K2 Standardprofil (nur `claude-code` und `node`, kein Ollama im Container):** `claude --model <modell> -p "<prompt>"` mit `ANTHROPIC_BASE_URL` und dem Platzhalter `ANTHROPIC_AUTH_TOKEN=ollama`, gesetzt vom Messskript im Aufruf (nicht in der Umgebung des Hosts). Nutzerpfad: Schritte vom `config set` bis zur Antwort und die Meldungen dazwischen.
- **K3 Sandbox-Profil, `discover` und `restricted`:** dieselben Aufrufe (K2) plus `ollama run <modell> "<prompt>"` mit `OLLAMA_HOST` des Container-Clients auf den Host-Server, plus `ollama launch claude --model <modell> --yes` (headless; Eingabe über stdin, falls `launch` sie annimmt). In `restricted`: die erzeugte Liste allein (K3a), mit dem Eintrag `host.docker.internal` in `devcontainer.sandbox.egress.allow` (K3b, erzeugt über `u-boot config set`); ob dnsmasq den Namen beantwortet (`no-hosts`) und nftables die Adresse zulässt, steht in den Rohdaten (Resolver-Log, `nft list set`). Greift K3b nicht, gilt K3c als **Handvariante**: eine zusätzliche `address=/host.docker.internal/<adresse>`-Zeile im Messcontainer.
- **K3d Hosts im Container:** Namen im Resolver-Log (`discover`) während K2/K3: was Claude Code und `ollama` im Container selbst anfragen; der Verkehr des Host-Servers zum Cloud-Dienst liegt außerhalb des Containers und ist **nicht** sichtbar.

*Weg (b) — angemeldeter Ollama-Server im Container (Sandbox-Profil):*

- **K4 Anmeldung.** `ollama signin` im Container-Server durch den Projektinhaber **an seinem Terminal** (`07-anmeldung.sh`): Ablauf (Ausgang des Aufrufs, Warten, Folgeaufruf), Hosts des Anmeldeablaufs aus dem Resolver-Log, Dateien als Pfad, Größe, Rechte (`stat`, kein Inhalt), Zugehörigkeit zum Modell-Volume gegen Workspace-Volume und Wurzel-Dateisystem. Die Anmelde-URL erscheint nur auf dem Terminal des Projektinhabers.
- **K5 Lauf.** `ollama run`, `claude` (direkt, `ANTHROPIC_BASE_URL=http://127.0.0.1:11434`) und `ollama launch claude` mit beiden Modellen in `discover`, dann `restricted` mit der Liste und den dann nötigen Einträgen (`ollama.com`, weitere aus dem Resolver-Log); Hosts nach der Anmeldung; Neuanlage des Containers mit erhaltenem Volume: überlebt die Anmeldung (K5b)?
- **K6 Modellname.** `glm-5.3-flash:cloud[1m]` gegen `glm-5.3-flash:cloud`: ein Aufruf je Name, um die Bedeutung von `[1m]` zu klären (Anthropic-Konvention oder Teil des Ollama-Namens).
- **K7 Variablen.** `OLLAMA_NO_CLOUD=1` mit Konto (`run`, `launch`; Server im Container und Client); `OLLAMA_API_KEY` mit dem **Platzhalter** `platzhalter` (gesetzt im Aufruf): ändert sich Meldung oder Verhalten gegenüber der Anmeldung, und wohin gehen die Anfragen (Resolver-Log)? Ein gültiger Schlüssel gehört zur Frage 4 (§6); ohne Antwort bleibt die Wirkung mit gültigem Schlüssel ungemessen.
- **K8 `launch` in `restricted` gegen `discover`:** warum `launch claude` ohne Konto in `restricted` das `claude` startete und in `discover` abbrach, jetzt mit Konto (Verification-Report, Übergabe Punkt 3).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Lastenheft, Spezifikation, Handbuch, Produktcode, ADR** — eigene Slices der Welle (Welle-Datei §4). Die Messung zeigt, was geht; sie schreibt keinen Wortlaut für das Lastenheft und hebt weder [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) noch [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) auf. Ein Weg über `--add-host` oder ein Relay ist keine Zusage und kein Erzeugnis von `u-boot`.
- **Anbindung an den Host durch `u-boot`** (`--add-host`, Relay, `--network host`, Änderung der Bindung des Host-Servers) — nicht Produkt, nur Handvariante des Messprojekts; der Slice ändert weder `OLLAMA_HOST` noch die Konfiguration des Host-Dienstes.
- **Jeder Aufruf des Host-Ollama und jedes Lesen von Dateien des Hosts** (`~/.ollama`, `~/.claude`, Konfigurationsdateien eines Agenten) — Zugangsregel 1. Der Host-Server wird ausschließlich aus Containern des Slice angesprochen.
- **Lasttests, Kontingent-Ausschöpfung, lokale Modelle, Qualität der Antworten** — die Zahl der Prompts ist durch die Obergrenze aus Frage 1 begrenzt; Cloud-Modelle haben keine lokale Last.
- **Ein gültiger `OLLAMA_API_KEY`** — nur mit ausdrücklicher Antwort auf Frage 4; sonst Platzhalter (K7).
- **Andere Plattformen und Engines** (`arm64`, Docker Desktop, Colima, Podman als Host-Engine, rootless Docker) — Zusage-Menge ist Linux `amd64` mit lokaler Engine ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)); die Abhängigkeit des Weges (a) von der Engine (`host-gateway` unter Linux) wird am gemessenen Host ausgewiesen, andere Hosts stehen als „nicht gemessen“.
- **Pin-Hebung von Ollama oder Claude Code** — eigener Slice nach [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) Entscheidung 9; Claude Code ist beim Bau nicht gepinnt, die Version steht im Protokoll.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **Sicherung vor dem ersten Lauf** ([`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)): (i) `redact()` ist an einer Gegenprobe geprüft, die für die Formen des Zugangs mit Konto (Anmelde-URL und -Code, Token von Ollama, Geräteschlüssel, Antwort-Zeilen) Platzhalter in diesen Formen schwärzt und neutrale Zeilen unverändert lässt, **mit Artefakt** (Rohdatei mit Zeitstempel, vor dem ersten Lauf geschrieben; der Lauf bricht ohne das Artefakt ab); (ii) `umgebung_pruefen()` bleibt an, ihr Ergebnis steht mit Zeitstempel in einer Rohdatei **vor** dem ersten Lauf; wird für `OLLAMA_API_KEY` (Frage 4) eine Ausnahme gebraucht, ist sie in `lib.sh` benannt (genau ein Name) und der Wert erreicht kein Skript auf einer Kommandozeile; (iii) die Suche (`06-suche.sh`) hat eine **benannte Ausnahmeliste** (der Platzhalter `ANTHROPIC_AUTH_TOKEN=ollama` und die Treffer der älteren, abgeschlossenen Protokolle, je mit Grund) und ist über das Verzeichnis dieser Messung und `tools/messung/ollama-cloud-konto/` leer; (iv) Rohdaten der Läufe mit Konto sind **von Hand gelesen**, Zeile für Zeile, bevor sie ins Repo kommen (Stufe 1 schreibt außerhalb des Repos, Stufe 2 kopiert nach dem Lesen); (v) die Rohdaten tragen **keinen** Antwort-Inhalt eines Modells und keine Anmelde-URL. Dies ist der dritte Beleg für [`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md), sobald ein Lauf mit Zugang stattfand; die Closure trägt ihn ein.
- [x] **Weg (a), Host-Server** (K1 bis K3, K6 an der Stelle, an der der Weg funktioniert): im Standardprofil und im Sandbox-Profil (`discover`, `restricted` ohne und mit Eintrag) gemessen, je Aussage Umgebung (Docker-Version, Host-Engine), Datum, Läufe (zwei, mit frischem Volume und neuem Container); die Abhängigkeit von der Engine ist ausgewiesen (`host-gateway`, Bindung des Dienstes, Relay); jede Handvariante steht in einer Tabelle; **bedingungsgleich formuliert** („mit dem Konto des Projektinhabers auf einem Linux-`amd64`-Host, Docker x, ein Host-Server, zwei Modelle“). *Beleg:* Skripte unter `tools/messung/ollama-cloud-konto/`, Protokoll mit Rohdaten unter `docs/messungen/slice-v1-ollama-cloud-messung-mit-konto/` ([`harness/conventions.md`](../../../../harness/conventions.md#ablage-von-messprotokollen-und-skripten)).
- [x] **Weg (b), Server im Container** (K4 bis K8): Anmeldung durch den Projektinhaber am Terminal, Hosts nach der Anmeldung, Dateien (Pfad, Größe, Rechte, kein Inhalt; ob im Modell-Volume, im Workspace-Volume oder im Wurzel-Dateisystem), ob der Geräteschlüssel ein Zugangsdatum im Sinne von [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) ist (die Messung liefert Ort und Überleben des Containers, die **Wertung** trifft die Folge-ADR), `OLLAMA_NO_CLOUD` mit Konto, Wirkung von `OLLAMA_API_KEY` (Platzhalter), `launch` in `restricted` gegen `discover`; die Rückabwicklung der Anmeldung (`ollama signout`, Schlüssel in den Kontoeinstellungen, Volume entfernt) ist im Protokoll ausgewiesen. Dazu die **Liste „nicht gemessen“** und die Nutzerpfad-Tabelle je Profil (Schritt, Meldung, Beleg) als Eingabe für ADR und Validator.
- [ ] `make gates` grün; `make docs-check` 0 Befunde; `make ablage-pruefung` grün; `make doc-commits RANGE="origin/main..HEAD"` 0 Befunde.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: Welle-Datei und Roadmap (Zeilen der Folge-Slices), `docs/messungen/README.md` (Bestand); kein öffentlicher Vertrag berührt (Lastenheft, Spezifikation, Handbuch, `CHANGELOG.md` bleiben unverändert).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/messung/ollama-cloud-konto/` | neu | Messskripte des Weges (a) und (b): `00-vorlauf.sh` (Gegenprobe `redact`, `umgebung_pruefen`, Artefakt mit Zeitstempel), `01-messstand.sh`, `02-projekte.sh` (Varianten `std`, `res`, `hgw` = `res` mit `--add-host=host.docker.internal:host-gateway` in `runArgs`, durch das Skript gesetzt), `03-lauf-vorbereitung.sh`, `04-relay.sh` (Relay-Container, nur für die Dauer eines Falls, Adresse aus `getent` im Container), `05-lauf-host.sh` (Weg a), `06-lauf-container.sh` (Weg b ohne Anmeldung), `07-anmeldung.sh` (nur der Projektinhaber, Ausgabe auf sein Terminal, nicht in Rohdaten), `08-aufraeumen.sh`, `09-suche.sh` (mit benannter Ausnahmeliste), `inner.sh`, `claude-stub.sh` falls nötig; Muster `tools/messung/ollama-cloud/` (gelesen: `lib.sh`, `06-suche.sh`); **eigenes `lib.sh`**, damit die abgeschlossene Messung unverändert bleibt; kein Skript liest oder setzt einen Zugang, keine Datei des Hosts wird eingebunden |
| `docs/messungen/slice-v1-ollama-cloud-messung-mit-konto/` | neu | `README.md` (Messstand, Fälle K1 bis K8, Handvarianten, Nutzerpfad je Profil, Liste „nicht gemessen“, Zugangsdaten und Rohdaten) und Rohdaten (`.txt`, `.log`), nach Handlesung |
| `docs/messungen/README.md` | update | Zeile im Bestand |
| `docs/reviews/2026-10-…-slice-v1-ollama-cloud-messung-mit-konto.md`, `…-verify.md` | neu | Review- und Verification-Report ([`MR-009`](../../../../harness/conventions/MR-009-skills-und-review-ablage.md)) |
| `docs/plan/planning/observations/BEO-DOK/…`, `BEO-WZG/…` (`state.md`, `evidence/slice-v1-ollama-cloud-messung-mit-konto.md`) | update / neu | Register-Fortschreibung bei der Closure; kein Zähler wird gesetzt |
| diese Datei, `../welle-ollama-cloud-modelle.md`, `open/README.md`, `in-progress/roadmap.md` | update | diese Datei: `Verantwortlich:`, Gefundener Stand (gemessen), DoD-Häkchen; Welle-Datei: Zeile des Slice, Frage 2 und Out-of-Scope (Host-Server als Messgegenstand); `open/README.md` und Roadmap: Zeiger |
| `internal/`, `cmd/`, `spec/`, `docs/user/`, `README.md`, `README.de.md`, `CHANGELOG.md` | **unverändert** | Produktcode, Verträge und Nutzertext folgen in den Folge-Slices der Welle; belegt durch die gelesenen Stellen in §1 (*Gefundener Stand*: `egress-init.sh.tmpl` Zeilen 60 bis 96, `devcontainer_egress.go` Zeilen 78 bis 93, `devcontainer_sandbox.go` Zeilen 161 und 441, `templates.go` Zeile 89, `sandbox.go` Zeile 122, Handbuch Zeilen 2023 und 2043, `lastenheft.md` Zeilen 441 bis 503; [`SPEC-047`](../../../../spec/spezifikation.md#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen) ohne Aussage zu Cloud-Modellen und zum Host-Server: Suche nach „cloud“ und „host-gateway“ über `spezifikation.md` ohne Treffer außer den Docker-Hosts `cloudflare`/`cloudfront`) |

**Abgleich mit `git diff --stat origin/main..HEAD` (Stand nach der Fix-Runde, 2026-10-24):** Gegenüber der Tabelle gilt im Ist:

- Variante `hgw` entfällt; die Handvariante `--add-host=host.docker.internal:host-gateway` steht in `runArgs` **aller** Projekte (`std`, `res`, `alw`); `alw` ist `res` plus `devcontainer.sandbox.egress.allow=host.docker.internal`.
- Unter `tools/messung/ollama-cloud-konto/` liegen zusätzlich `10-uebernehmen.sh` (Stufe 2 der Rohdaten), `redact-probe.sh`, `Dockerfile.cli` und `lib.sh`; `claude-stub.sh` entstand nicht. Der Relay-Befehl `probe` und die Fälle `nftadd`, `curl`, `versionen` sind Ergänzungen von `04-relay.sh` und `05-lauf-host.sh`.
- Zusätzlich berührt, nicht in der Tabelle: `docs/plan/planning/done/slice-v1-ollama-cloud-messung.md` (zwölf Zeilen Verweiskorrektur aus `make slice-mv`), die `state.md` von zwei Beobachtungen (Link-Nachzug aus `make slice-mv`), `docs/plan/planning/welle-ollama-cloud-modelle.md` (Link-Nachzug), `in-progress/roadmap.md` (Zeiger auf den Slice, nötig für `docs-check`) und `docs/reviews/2026-10-24-slice-v1-ollama-cloud-messung-mit-konto.md` (Review-Report).
- Noch nicht berührt und bei der Closure fällig: Register-Evidence, Welle-Datei (Zeile des Slice), Closure-Notiz, Verification-Report.

Ansatz:

- **Reihenfolge:** `00-vorlauf.sh` (Artefakt) → `01` bis `03` → Weg (a) mit **einem** Modell zuerst (`glm-5.3-flash:cloud[1m]`), K1 bis K3; erst wenn der Weg trägt oder sein Hindernis feststeht, das zweite Modell → Weg (b) mit Anmeldung durch den Projektinhaber → `08-aufraeumen.sh`. Bei einem Hindernis, das mehr Anfragen nicht klären, wird nicht wiederholt (keine Schleifen).
- **Umgebung:** Bereinigung nicht mit `env -i NAME=<wert>`, sondern im selben Prozess durch `unset` aller nicht erlaubten Namen (Positivliste), damit auch eine benannte Ausnahme (Frage 4) nie als Wert auf einer Kommandozeile steht; Durchreichen in den Container als `-e NAME` ohne Wert. Der Platzhalter `ANTHROPIC_AUTH_TOKEN=ollama` und `ANTHROPIC_BASE_URL` werden im Aufruf im Container gesetzt, nie in der Umgebung des Hosts.
- **Speicher:** ein Messcontainer zur Zeit, kein Parallellauf, kein lokales Modell (Cloud-Modelle sind Verweise von wenigen hundert Byte, Protokoll des ersten Slice, Abschnitt F), Image einmal bauen und im zweiten Lauf wiederverwenden, Aufräumen nach jedem Weg; Docker-Speicher des Hosts und das Nachbar-Repo bleiben unberührt (Präfix `u-boot-messung-ocloud-konto`, Löschen nur nach Namen).
- **Messhygiene aus dem ersten Slice:** `timeout -k 5` bei jedem Aufruf im Container; die Überschrift der Rohdaten für `allow` trägt den Modus `allow` (das Skript übergibt Modus und Bezeichnung getrennt); Lauf 2 mit frischem Volume für Weg (a); Rohdaten werden nicht doppelt geschwärzt.
- **Rückabwicklung (b):** nach dem Fall `ollama signout` im Container, Entfernen des Volumes (`praefix_entfernen`), und der Projektinhaber bestätigt den Kontozustand nach der Messung (Schritt des Protokolls): bestätigt wird der Kontozustand (kein Eintrag aus der Messung), nicht das Entfernen eines bestimmten Schlüssels. Die Anmeldung trägt den Schlüssel nach Beobachtung des Projektinhabers im Konto ein; die Wirkung von `signout` dort ist **nicht gemessen**, der Rückbau im Konto durch den Nutzer bleibt zu bestätigen.
- **Relay (K1b):** Image mit Digest gepinnt (gelesen beim ersten Pull, ins Protokoll), `--network host`, bindet nur die Adresse, die `host-gateway` im Container ergibt, nie `0.0.0.0`; läuft nur während eines Falls; `docker ps` nach dem Fall belegt, dass er weg ist. Solange er läuft, kann jeder Container im Netz der Adresse den angemeldeten Server des Hosts nutzen — die Messung läuft deshalb ohne fremde Container im Messnetz, soweit der Skript es prüfen kann (Anzahl der Container vorher und nachher).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; kein anderer Slice in `in-progress/` (WIP-Limit 1, gelesen beim Start); die Antworten des Projektinhabers auf die Fragen 1 bis 5 in §6 liegen vor (Frage 1 Obergrenze der Prompts, Frage 2 Relay, Frage 3 Anmeldung im Container und Rückabwicklung sind Bedingungen für den jeweiligen Weg); der Host-Server des Projektinhabers läuft und ist angemeldet (zu Beginn aus einem Container geprüft, nicht am Host); Netz und Docker auf einem Linux-`amd64`-Host mit lokaler Engine; die Umgebung des Messlaufs trägt keine Zugangsvariable (Prüfung nach Namen, Artefakt des Vorlaufs); der Slice [`slice-v1-ollama-cloud-messung`](../done/slice-v1-ollama-cloud-messung.md) liegt in `done/`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Weg (a) und Weg (b) tragen zusammen mehr als drei Liefer-Punkte (etwa weil K1 ergibt, dass der Weg (a) einen eigenen Unterbau braucht); die Teilung entscheidet der Planner, Weg (b) kann dann ein eigener Slice werden.
- `in-progress` → `open` (blockiert): (a) der Host-Server ist nicht angemeldet oder nicht erreichbar, (b) der Projektinhaber beantwortet Frage 2 oder 3 mit „nein“ und der betroffene Weg entfällt (dann misst der Slice den anderen Weg und weist den entfallenen als „nicht gemessen“ aus), (c) die Obergrenze der Prompts ist vor Abschluss eines Weges erreicht (der Rest steht als „nicht gemessen“).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; `make gates` und `make docs-check` sind grün; das Messprotokoll liegt im Repo, trägt je Aussage Umgebung, Datum und Läufe, führt die ungemessenen Punkte und keinen Tokenwert (Suche leer über das Verzeichnis dieser Messung und `tools/messung/ollama-cloud-konto/`, benannte Ausnahmeliste im Skript, Treffer außerhalb davon von Hand gelesen); das Artefakt der Gegenprobe liegt **vor** dem ersten Lauf im Verlauf; Review-Report und Verification-Report liegen vor; die Closure-Notiz mit Steering-Loop-Lerneintrag ist geschrieben; die Zahl der gestellten Prompts steht im Protokoll und liegt unter der Obergrenze aus Frage 1. Die Wertung (Zusage, Lastenheft, Aufnahme von Hosts, Geräteschlüssel als Zugangsdatum) ist keine Bedingung dieses Slice: sie gehört der Folge-ADR.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht. Die Ausgänge werden bei der Closure gesetzt; hier steht je Risiko die Vorsorge.

- **Zugangsdaten gelangen in Rohdaten, Skripte, Chat oder Verlauf** (Anmelde-URL und -Code, Token, Geräteschlüssel, Antwort-Inhalt; ein Lauf mit Zugang ist der dritte Beleg). — **Vorsorge:** DoD 1 (Gegenprobe mit Artefakt vor dem ersten Lauf, Ausnahmeliste, Handlesung jeder Rohzeile in zwei Stufen, kein Antwort-Inhalt, Anmelde-URL nur auf dem Terminal des Projektinhabers), Zugangsregeln in §1, Bereinigung der Umgebung im selben Prozess; Rohdaten werden **von Hand gelesen** ([`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md); `redact` ist formgebunden). **Ausgang:** bei der Closure.
- **Eine Aussage reicht weiter als die Messung** (ein Host-Server, ein Konto, zwei Modelle, eine Claude-Code-Version; „läuft“ ist ein Lauf, nicht eine Zusage). — **Vorsorge:** Aussagen bedingungsgleich („mit dem Konto des Projektinhabers auf einem Host“), die Liste „nicht gemessen“, kein Wortlaut für das Lastenheft; Wertung nur in der Folge-ADR ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)). **Ausgang:** bei der Closure.
- **Die Messvariante ist nicht die Auslieferung** (`--add-host` in `runArgs`, Relay, Platzhalter-Variablen im Aufruf, `docker exec` statt Sitzung der CLI, Weg (b) nur im Sandbox-Profil, weil nur dort Ollama im Container liegt). — **Vorsorge:** Tabelle der Handvarianten im Protokoll mit „keine begründet eine Zusage“; K2 misst den Standardpfad ohne Ollama im Container ([`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](../observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md)). **Ausgang:** bei der Closure.
- **Der Host-Server wird über das Netz für andere Container nutzbar** (Relay auf der Gateway-Adresse; angemeldeter Dienst). — **Vorsorge:** Frage 2; Relay nur während eines Falls, nur die aus dem Container ermittelte Adresse, nie `0.0.0.0`, Aufräumen belegt; keine fremden Container im Messnetz (Anzahl vorher und nachher). **Ausgang:** bei der Closure.
- **Kontingent und Kosten des Cloud-Dienstes** (Claude Code sendet mit Systemprompt und Werkzeugliste größere Anfragen als ein Prompt von einem Satz; die Größe ist nicht gemessen). — **Vorsorge:** Frage 1 (Obergrenze), ein Modell zuerst, keine Wiederholungsschleifen, die Zahl der gestellten Prompts steht im Protokoll. **Ausgang:** bei der Closure.
- **Die Anmeldung im Container hinterlässt einen Eintrag im Konto des Projektinhabers** (Annahme des Plans; nach Beobachtung des Projektinhabers eingetreten, der Eintrag stand zum Prüfzeitpunkt nach der Messung nicht mehr im Konto). — **Vorsorge:** Frage 3, Wegwerf-Volume, `ollama signout` (Wirkung im Konto nicht gemessen), Schritt „Kontozustand prüfen und bestätigen (kein Eintrag aus der Messung)“ im Protokoll; vom Projektinhaber am 2026-10-24 bestätigt (nur ein Eintrag, nichts von heute). **Ausgang:** bei der Closure.
- **Host-Werkzeuge außerhalb von Docker** (Host-`sed`, Host-`ollama`, Heredoc-Schreiben; im ersten Slice trat Host-`sed` auf). — **Vorsorge:** Skripte rufen nur `docker`, `make`, `bash` und Coreutils; Auswertungen in Containern; kein Aufruf des Host-Ollama ([`BEO-WZG/host-werkzeug-ausserhalb-docker`](../observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md)). **Ausgang:** bei der Closure.
- **Speicher des Hosts** (Nachbar-Repo wurde durch OOM beendet). — **Vorsorge:** Ansatz „Speicher“ in §3; ein Container zur Zeit, Image einmal bauen. **Ausgang:** bei der Closure.
- **Die Messung eines früher ausgeschlossenen Anschlusses wird als Zusage gelesen** ([`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md), Handbuch Zeile 2043). — **Vorsorge:** §1 *Spannung*, Frage 5, Protokoll ohne Wortlaut für das Lastenheft; Handbuch und ADR bleiben unverändert, bis die Folge-ADR entschieden hat. **Ausgang:** bei der Closure.
- **Plan-Aussagen über Code, die der Code nicht trägt** (Egress-Mechanismus, `no-hosts`, Katalog, Handbuch). — **Vorsorge:** §1 nennt Datei und Zeile je Aussage; der Reviewer gleicht §3 gegen `git diff --stat` ab ([`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md), verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.10). **Ausgang:** bei der Closure.
- **Claude Code ist beim Bau nicht gepinnt, Modellnamen der Cloud ändern sich** (die Messung ist nicht wiederholbar). — **Vorsorge:** Versionen und Datum im Messstand; Wiederholung ist kein Ziel. **Ausgang:** bei der Closure.

### Fragen an den Projektinhaber (beantwortet, 2026-10-24)

1. **Kontingent:** Darf die Messung Modelle der Cloud über den Host-Server verbrauchen, und bis zu welcher **Obergrenze**? *Vorschlag des Planners:* höchstens **40 Prompts** (je ein Satz Eingabe, ein Wort erbeten), davon bis zu 20 im Weg (a) (Standardprofil 4; Sandbox `discover` 12; `restricted` 4), bis zu 14 im Weg (b), je 2 für `OLLAMA_NO_CLOUD` und für den Platzhalter von `OLLAMA_API_KEY`, 2 Reserve; Minimalvariante: nur `glm-5.3-flash:cloud[1m]`, dann etwa 20. Ein Prompt über Claude Code kann mehrere Anfragen auslösen (nicht gemessen). *Antwort (2026-10-24):* höchstens 40 Prompts (20 Weg a, 14 Weg b, je 2 für `OLLAMA_NO_CLOUD` und den Platzhalter von `OLLAMA_API_KEY`, 2 Reserve); beide Modelle; nie überschreiten, Zähler im Protokoll.
2. **Relay für den Weg (a):** Der Host-Server lauscht laut Angabe auf `127.0.0.1:11434`; ein Container erreicht diese Adresse über das Gateway nicht ohne Weiteres (Hypothese). *Vorschlag:* ein Hilfs-Container mit `--network host`, der nur für die Dauer eines Falls und nur auf der Gateway-Adresse weiterreicht (§3). *Alternative:* der Projektinhaber ändert die Bindung seines Dienstes selbst (der Slice tut es nicht), oder der Weg (a) entfällt als „nicht erreichbar gemessen“. *Antwort (2026-10-24):* Relay erlaubt (Hilfs-Container `--network host`, nur auf der Gateway-Adresse, nur während eines Falls, danach entfernt).
3. **Anmeldung im Container (Weg b):** Darf der Projektinhaber einen Wegwerf-Container anmelden (er führt `07-anmeldung.sh` an seinem Terminal aus, bestätigt die Anmelde-URL im Browser, entfernt danach den Schlüssel in den Kontoeinstellungen)? *Antwort (2026-10-24):* ja; der Projektinhaber meldet selbst an seinem Terminal an (`07-anmeldung.sh`) und entfernt danach den Geräteschlüssel in den Kontoeinstellungen.
4. **`OLLAMA_API_KEY`:** Stellt der Projektinhaber einen gültigen Schlüssel als Umgebungsvariable des Hosts bereit (nur für diesen Fall, nie in Datei, Skript, Protokoll oder Chat)? *Vorschlag:* nein; die Wirkung wird mit einem Platzhalter gemessen, die Wirkung mit gültigem Schlüssel bleibt „nicht gemessen“. *Antwort (2026-10-24):* nein; nur der Platzhalter, die Wirkung eines gültigen Schlüssels bleibt ungemessen.
5. **Vorgabe vom 2026-10-18** („kein Anschluss an einen Ollama des Hosts“, [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md)): Hebt die Entscheidung vom 2026-10-24 sie nur für die **Messung** auf (die Folge-ADR bewertet das Ergebnis), oder soll die Folge-ADR den Host-Server als Option zusagen können (dann `Supersedes` oder `Schärft:` und eine Änderung von [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) im CR)? *Vorschlag:* nur für die Messung; die Wertung nach dem Ergebnis. *Antwort (2026-10-24):* die vorgenannte Vorgabe gilt für diese Messung nicht; die Folge-ADR bewertet das Ergebnis; Handbuch und Lastenheft bleiben in diesem Slice unverändert.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

*(Wird bei der Closure gefüllt: Was hat funktioniert · Was ging anders als geplant · Steering-Loop-Eintrag · Beobachtungs-Register · Folge-Slices · Risiken aus §6 mit Ausgang · Drei Paarungen. Stand: Messung abgeschlossen, Review liegt vor; die Closure-Notiz folgt.)*

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind `tools/` (Kürzel WZG, Messskripte unter `tools/messung/`) und `docs/messungen/`, `docs/plan/` (Kürzel DOK); beide stehen als eigene Sub-Areas in der Modus-Deklaration ([`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area)). `internal/`, `spec/` und `docs/user/` werden gelesen, nicht geändert (§3).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen am 2026-10-24 (Zähler = Dateien in `evidence/`; Stand vor der Closure des ersten Slice). Treffer: [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (21), [`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](../observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md) (4) und [`BEO-DOK/vertrag-schweigt-zum-nutzerpfad`](../observations/BEO-DOK/vertrag-schweigt-zum-nutzerpfad/observation.md) (4), alle über der Schwelle, ohne zugewiesenen Ausgang (Entscheidung des Projektinhabers steht aus; der Slice weist nichts zu), Vorsorge in §2 und §6, der Nutzerpfad je Profil als Tabelle des Protokolls (DoD 3); [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (51, verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.10, §1 nennt die Fundstellen); [`BEO-WZG/host-werkzeug-ausserhalb-docker`](../observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md) (15; §6); [`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md) (2; ein Lauf mit Zugang wäre der dritte Beleg, DoD 1 und §6); [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md) (15; kein Lastenheft-Text in diesem Slice); [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md) (2; Aufräumen nach Präfix und Namen, §3).

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

**Alle berührten Sub-Areas GF** (WZG und DOK, Modus-Deklaration in `harness/conventions.md`): kein Modus-Begründungsblock.
