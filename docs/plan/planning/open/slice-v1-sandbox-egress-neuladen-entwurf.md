# Slice v1: Egress-Liste im laufenden Container neu laden — Entwurf, Messung und Sicherheits-Optionen (`LH-FA-DEV-008`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-sandbox-egress-ollama`](../welle-sandbox-egress-ollama.md) — vierter Slice der Welle (Entwurf und Messung; ändert weder Lastenheft noch Produktcode). Er legt am Ende die Folge-Slices als Dateien an, wenn der Projektinhaber eine Variante wählt (§1, Verzweigung); die Welle trägt sie dann in §4 ein (die Zahl der Slices ändert sich durch einen Eintrag dort, nicht stillschweigend).

**Bezug:** [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Egress-Restriktion, Guardrail-Aussage), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil: Schadensbegrenzung, kein Pfad vom Container zum Host), [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Startscripte, Degradation); [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) (Mechanismus; Alternative D: Restriktion außerhalb des Containers), [`ADR-0022`](../../adr/0022-fallback-ausweisung-statusdateien-im-container.md) (Host liest den Container ohne Programmaufruf).
Anlass: Projektinhaber, 2026-10-18, Vorschlag (b): Neu-Laden im laufenden Container ohne Neustart, vom Host aus; ein Agent im Container soll sich damit nicht selbst Hosts freischalten können.

**Berührte Spec-Stellen:** [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) (Ermittlung durch den Host: „ohne ein Programm im Container auszuführen“, `docker cp`, kein `docker exec`), [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) (Liste), `lastenheft.md` §3 [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (nur als Wortlaut-Vorschlag unten). Der Slice **liest** sie und ändert keine; die Nachführung folgt in den Folge-Slices und im Spezifikations-Slice der Welle.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Architect klärt **zuerst**, was ein Prozess im Sandbox-Container ohne `sudo` und mit `sudo` an der Egress-Beschränkung tun kann, misst den Weg „Liste im laufenden Container ersetzen“ auf Colima und legt dem Projektinhaber die **Optionen** mit ihren Sicherheitsfolgen vor — **ohne zu fragen und ohne zu entscheiden**; die Sicherheitsentscheidung gehört dem Projektinhaber. Ergebnis: eine Folge-ADR, die Optionen und Messung trägt, und eine der Verzweigungen unten.

### Gefundener Stand im Code und in der Doku (gelesen am 2026-10-18)

- **Die Liste ist ins Image gebacken.** `generate.go` rendert bei Egress `.devcontainer/egress-init.sh` (Kommentar „egress restriction script baked via COPY“); `Dockerfile.tmpl` hat `COPY egress-init.sh /usr/local/bin/u-boot-egress-init`; `egress-init.sh.tmpl` trägt die Hosts als `hosts="…"`. Eine geänderte Liste braucht damit `u-boot generate devcontainer` **und** ein neu gebautes Image. Das Handbuch (Abschnitt *Internet-Beschränkung (Egress)*) sagt: „Nach einer Änderung führen Sie `u-boot generate devcontainer` aus und starten den Container neu.“ — **ob „neu starten“ genügt oder „neu aufbauen“ gemeint ist, ist ungemessen** (Messfrage M4); trifft die Befürchtung zu, ist das ein Doku-Fehler unabhängig vom Neu-Laden.
- **Wie das Script läuft:** `devcontainerTemplateData` (`internal/hexagon/application/devcontainer_sandbox.go`) setzt `postStartCommand` auf `sudo sh /usr/local/bin/u-boot-egress-init` (vor dem Docker-Start-Script). Das Script ist laut seinem Kommentar „re-run safe“: es beendet einen vorigen `dnsmasq` über die PID-Datei, löscht die Tabellen `u_boot_egress` und `u_boot_egress6` (`nft delete table`), schreibt `/etc/u-boot-egress-dnsmasq.conf` neu, legt die Tabellen neu an, startet `dnsmasq` und setzt `/etc/resolv.conf` auf `127.0.0.1`; am Ende schreibt es die Statusdatei (`write_status`). Zwischen `nft delete table` und `nft add chain … policy drop` gibt es kein Filter-Regelwerk: **ein Fenster ohne Beschränkung** (Dauer ungemessen, Messfrage M2).
- **Wer was schreiben darf (aus Dockerfile und Doku gelesen, nicht gemessen):** Das Dockerfile setzt im Sandbox-Zweig `USER root` und am Ende `USER vscode`; `COPY` legt das Script als Root ab. `/var/lib/u-boot/status` gehört `vscode` (`chown vscode:vscode`), die Statusdatei ist deshalb von jedem Prozess des Benutzers beschreibbar ([`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), Absatz *Grenze*). Die Regeln selbst brauchen `NET_ADMIN` (`egressRelaxation`: `--cap-add=NET_ADMIN`) und damit Root im Container; **dass ein Prozess ohne `sudo` sie nicht ändern kann, ist nicht gemessen.**
- **Was das Lastenheft und die Doku heute sagen:** [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion): „Die Restriktion ist ein Guardrail und keine Sandbox-Grenze; Prozesse mit der nötigen Capability können sie aufheben.“ [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) (Kontext) und das Handbuch nennen zusätzlich `sudo` des Base-Images als den Weg; im Lastenheft steht `sudo` nicht. Der `postStartCommand` ruft `sudo sh` ohne Passwortabfrage auf, das setzt passwortloses `sudo` für `vscode` voraus (Folgerung des Planners; M1 misst es).
- **Der Host führt heute nie ein Programm im Container aus.** `DevcontainerInspector` (`internal/hexagon/port/driven/devcontainer_inspector.go`) hat genau `FindContainer` und `ReadFile`; sein Kommentar sagt „it never runs a program in the container“; [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) legt das als Eigenschaft fest („ohne ein Programm im Container auszuführen“, `docker cp` statt `docker exec`). Ein Neu-Laden vom Host aus wäre der **erste** Aufruf eines Programms im Container und braucht einen neuen Driven-Port (oder eine Erweiterung) samt Spezifikationsänderung. Der Container wird über das Label `devcontainer.local_folder` gefunden (`FindContainer`; entfernte Engine: `ErrContainerEngineNotLocal`).
- **Kein Pfad vom Container zum Host.** [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) verbietet Host-Socket und Host-Geheimnisse im Container; ein vom Host gestarteter Aufruf (Host → Container) berührt das nicht, ein Auslöser, den der Container anstößt und der auf dem Host wirkt, widerspräche ihm.

### Messfragen (Arbeitsplan; Docker auf Colima, Sandbox-Container mit Egress aus dem erzeugten Dockerfile, wie [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) und [`ADR-0021`](../../adr/0021-podman-build-im-sandbox-devcontainer.md))

| # | Frage | Warum |
|---|---|---|
| M1 | Als `vscode` **ohne** `sudo`: `nft list tables`, Schreibversuche auf `/usr/local/bin/u-boot-egress-init`, `/etc/u-boot-egress-dnsmasq.conf`, `/etc/resolv.conf`, `/run/u-boot-egress-dnsmasq.pid`, Signal an den `dnsmasq`; Verbindung zu einem nicht erlaubten Host. Danach dasselbe mit `sudo`. Ist `sudo` ohne Passwort verfügbar (`sudo -n true`)? | belegt, was ein Agent ohne `sudo` an der Liste nicht ändern kann, und was mit `sudo` doch; Grundlage für den Wortlaut zu „Agent schaltet sich selbst Hosts frei“ |
| M2 | Erneuter Lauf des Scripts im laufenden Container (`sudo sh /usr/local/bin/u-boot-egress-init`, mit anderer Liste): Dauer des Fensters ohne Filter, Schicksal laufender Verbindungen (ein laufender Download), DNS-Auflösung während des Neustarts von `dnsmasq`, Status `ok`. Variante: neue Tabelle atomar anlegen und die alte danach entfernen. | misst, ob „Script erneut ausführen“ als Neu-Laden trägt und ob es ein offenes Fenster hat |
| M3 | Zustellung vom Host: `docker cp` einer Liste in ein Root-eigenes Verzeichnis (Eigentümer und Rechte nach dem Kopieren), `docker exec -u root` mit dem Reload-Aufruf; Verhalten bei gestopptem Container und bei entfernter Engine. | Optionen R1 und R2 unten |
| M4 | Liste ändern, `generate devcontainer`, Container **neu starten** ohne Neubau: gilt danach die alte oder die neue Liste? | belegt oder widerlegt den Satz im Handbuch |

### Optionen (Sicherheitsentscheidung des Projektinhabers; hier dargelegt, nicht entschieden — die Spalten zu M1 füllt die Messung)

| Opt. | Weg | Auslöser | Was ein Prozess im Container ohne `sudo` damit kann | Was mit `sudo` | Eingriff |
|---|---|---|---|---|---|
| R0 | **Kein Neu-Laden**; die Doku nennt den richtigen Weg (Neubau des Images) | — | nichts neues | nichts neues | nur Handbuch (und M4) |
| R1 | **Vom Host**: ein Befehl (Arbeitsname `u-boot egress reload`) berechnet die Liste aus `u-boot.yaml` (dieselbe Berechnung wie [`slice-v1-cli-egress-anzeige`](slice-v1-cli-egress-anzeige.md)), stellt sie per Container-Engine als Root zu und führt das Reload-Script aus (`docker exec -u root`) | Nutzer auf dem Host | nichts: der Auslöser liegt auf dem Host, die Liste in einer Root-eigenen Datei | wie heute (`sudo` kann die Regeln ohnehin aufheben) | neuer Driven-Port (`Exec`-artig), Reload-Script im Image (einmal Neubau), Lastenheft-Zusage, Spec-Änderung in [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) |
| R2 | **Zwei Schritte**: Host legt die Liste per `docker cp` in ein Root-eigenes Verzeichnis, der Nutzer führt im Container `sudo sh /usr/local/bin/u-boot-egress-reload` aus | Nutzer im Container, mit `sudo` | nichts (braucht `sudo`) | alles, was `sudo` heute schon kann | kleinerer Port (nur Schreiben), aber eine Aktion im Container; Doku |
| R3 | **Im Container ohne `sudo`** auslösbar (Beobachter der Datei `/workspaces/…/egress.allow` oder eine `sudoers`-Regel für einen Reload-Befehl mit nutzerbestimmter Liste) | Prozess im Container | **kann sich selbst Hosts freischalten** | wie heute | **verworfen als Vorschlag**: widerspricht dem Wunsch des Projektinhabers; hier nur genannt, damit die Ablehnung begründet ist |
| R4 | **Außerhalb des Containers** (Proxy oder Netz-/DNS-Ebene, [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) Alternative D) | Host | nichts, auch nicht mit `sudo` | nichts | eigener Vorgang, Host-Konfiguration ist Nicht-Ziel von `u-boot` |
| S | **Zusatz zu R1/R2**: `sudo` für `vscode` entfernen (dann auch kein `sudo sh` im `postStartCommand`; die Start-Scripte brauchen einen anderen Weg zu Root) | — | wie R1 | **nichts** an der Egress-Regel | ändert die Benutzerrechte aus [`LH-FA-DEV-004`](../../../../spec/lastenheft.md#lh-fa-dev-004--benutzerrechte) und alle drei Start-Scripte; eigener Slice, nicht Teil dieser Welle |

**Einordnung des Planners (nicht entscheidend):** R1 ist der Weg, der dem Wunsch („vom Host, ohne dass der Agent sich selbst freischaltet“) am nächsten kommt, ohne die Guardrail-Aussage zu ändern: Solange `vscode` passwortloses `sudo` hat (M1), kann ein Agent die Regeln heute schon aufheben; R1 fügt keine Fähigkeit hinzu, **löst aber das Versprechen „kann sich nicht selbst freischalten“ nicht ein** — dafür braucht es S oder R4. Der Wortlaut des Lastenhefts muss das ausweisen (Vorschlag unten).

### Wortlaut-Vorschlag (alt/neu; **nur** wenn R1 oder R2 gewählt wird und M1 es trägt; Freigabe des Projektinhabers vor dem Lastenheft-Commit im Folge-CR-Slice)

Version: die nächste freie Nummer zum Zeitpunkt des Lastenheft-Commits (nach [`slice-cr-sandbox-egress-anzeige`](slice-cr-sandbox-egress-anzeige.md) also `0.15.0` oder höher, je nach Reihenfolge des Landens mit dem Ollama-CR).

**Stelle 1 — [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), Spiegelstrich zur Grenze der Restriktion.**

Alt:

> - Die Restriktion ist ein Guardrail und keine Sandbox-Grenze; Prozesse mit der nötigen Capability können sie aufheben. Die Dokumentation muss das ausdrücklich sagen.

Neu (Satzteil nach dem Semikolon gilt nur, wenn M1 bestätigt, dass `sudo` des Base-Images genügt):

> - Die Restriktion ist ein Guardrail und keine Sandbox-Grenze; Prozesse mit der nötigen Capability können sie aufheben, im Container des Base-Images genügt dafür `sudo`, und ein Agent mit `sudo` kann sich daher auch Hosts selbst freischalten. Die Dokumentation muss das ausdrücklich sagen.

**Stelle 2 — [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), neuer Spiegelstrich (Variante R1).**

> - Eine geänderte Allowlist lässt sich im laufenden Container ohne Neubau des Images und ohne Neustart des Containers wirksam machen, mit einem Befehl auf dem Host (`u-boot egress reload`). Der Befehl führt dazu ein Programm im Container aus; ein Auslöser im Container, den ein Prozess ohne Root-Rechte nutzen könnte, gibt es nicht. Er greift nur auf den Container des Projektordners zu und nicht auf eine entfernte Container-Engine; ohne laufenden Container endet er mit einem Umgebungsfehler (Exit-Code `11`, [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)).

(Variante R2 ersetzt „mit einem Befehl auf dem Host“ durch „mit einem Befehl auf dem Host, der die Liste bereitstellt, und einem Aufruf im Container“.) Die Zeile zu [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Abfragebefehl, `--quiet`) entfällt: `reload` ändert etwas und gibt unter `--quiet` nur Fehler aus (Absatz *Unter `--quiet`* in [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging); der Folge-CR nimmt den Befehl dort in die Aufzählung der ändernden Subkommandos auf).

**Nachzug außerhalb des Lastenhefts (nicht Teil dieses Slice):** [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) sagt, `doctor`, `generate devcontainer` und `init` lesen „ohne ein Programm im Container auszuführen“; das gilt für diese drei Befehle weiter, `reload` ist der erste Befehl mit einem Programmaufruf im Container — die Spezifikation grenzt den Satz entsprechend ein (Spezifikations-Slice der Welle).

### Verzweigung nach der Messung (vorab benannt)

- **V1 — R0:** Der Projektinhaber wählt kein Neu-Laden. Der Slice liefert die Handbuch-Korrektur (M4) und die Entscheidung als ADR; kein Folge-Slice. Die Welle schließt ohne Neu-Laden.
- **V2 — R1 oder R2:** Der Slice legt zwei Dateien in `open/` an — `slice-cr-sandbox-egress-neuladen` (Lastenheft-CR mit dem freigegebenen Wortlaut) und `slice-v1-sandbox-egress-neuladen` (Code, Test, Doku) — und trägt sie in [`welle-sandbox-egress-ollama`](../welle-sandbox-egress-ollama.md) §4 und `open/README.md` ein. Der Neu-Laden-Code folgt dem CR.
- **V3 — R4 oder S:** Eigener Vorgang außerhalb der Welle; der Slice trägt die Vorlage (ADR) und einen Eintrag in der Roadmap (*Nächste Wellen*) bzw. beim Projektinhaber.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Umsetzung (Befehl, Port, Reload-Script) und der Lastenheft-Commit** — Folge-Slices nach der Entscheidung (V2); ein Slice, der misst und entscheidet, setzt nicht um (Rollentrennung, `harness/roles.md`).
- **Das Entziehen von `sudo`** (Option S) — ändert Benutzerrechte und alle Start-Scripte; eigener Vorgang mit eigener Messung.
- **Die Anzeige der wirksamen Liste** — [`slice-v1-cli-egress-anzeige`](slice-v1-cli-egress-anzeige.md); `reload` nutzt später dieselbe Berechnung, aber der Slice hängt nicht an ihr.
- **Wildcards** — [`slice-v1-egress-wildcard-pruefung`](../in-progress/slice-v1-egress-wildcard-pruefung.md).
- **Docker Desktop, Linux-Host, amd64** — nicht gemessen, nicht zugesagt; jede Aussage nennt den Messstand (Colima, arm64).

## 2. Definition of Done

- [ ] **Messung:** M1 bis M4 sind gemessen und mit Messstand (Docker-, Colima-Version, Architektur, Base-Image-Digest) in der Evidence und in §3 festgehalten; die Tabelle der Optionen trägt die gemessenen Zellen („ohne `sudo`“, „mit `sudo`“), und die Grenzen sind benannt.
- [ ] **Entscheidungsvorlage und Folge-ADR:** Eine Folge-ADR zu [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) (nächste freie Nummer, Index `docs/plan/adr/README.md`; `Proposed` bis zur Entscheidung des Projektinhabers, danach `Accepted`) legt Optionen, Messung und Einordnung dar; [ADR-0012](../../adr/0012-devcontainer-egress-firewall.md) bleibt unverändert (`AGENTS.md` §3.5).
- [ ] **Verzweigung ausgeführt:** V1, V2 oder V3 ist in §7 und in der Welle ([`welle-sandbox-egress-ollama`](../welle-sandbox-egress-ollama.md) §4, §5) vollzogen; bei V2 existieren die beiden Folge-Slices als Dateien in `open/` mit dem freigegebenen Wortlaut; die Handbuch-Aussage „neu starten“ ist nach M4 richtig gestellt.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `docs/user/benutzerhandbuch.md` (Satz zum Neustart), ADR-Index; `CHANGELOG.md` nur, wenn M4 eine nutzersichtbare Korrektur ergibt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo mit Wellen von der Welle-Closure geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/adr/00NN-egress-liste-im-laufenden-container.md` (nächste freie Nummer) und `docs/plan/adr/README.md` | neu / update | Folge-ADR zu [ADR-0012](../../adr/0012-devcontainer-egress-firewall.md) mit Optionen, Messung, Einordnung |
| `docs/user/benutzerhandbuch.md` | update | Abschnitt *Internet-Beschränkung (Egress)*, Satz „… starten den Container neu“ nach M4; bei V2 erst im Code-Folge-Slice der Befehl |
| `CHANGELOG.md` | update (bedingt) | nur bei nutzersichtbarer Korrektur nach M4 |
| `docs/plan/planning/open/slice-cr-sandbox-egress-neuladen.md`, `slice-v1-sandbox-egress-neuladen.md`, `docs/plan/planning/open/README.md`, `docs/plan/planning/welle-sandbox-egress-ollama.md` | neu / update (bedingt, V2) | Folge-Slices und Welle-Eintrag |
| `docs/plan/planning/in-progress/roadmap.md` | update (bedingt, V3) | Eintrag in *Nächste Wellen* |
| `docs/reviews/` (Review-Report, Verification Evidence) | neu | Review und Verifikation |
| `docs/plan/planning/observations/…/evidence/` | neu (bedingt) | Beleg, falls eine Klasse wiederkehrt |
| diese Slice-Datei | update | Messung (§3), Closure (§7) |

Kein Produktcode, kein Test, kein Lastenheft in diesem Slice. Abgleich der Tabelle gegen `git diff --stat` nach der **letzten** Fix-Runde (`AGENTS.md` §3.10).

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; kein anderer Slice in `in-progress/` (WIP-Limit 1). Er hängt von keinem Slice der Welle ab und kann nach dem Wildcard-Slice oder parallel zur Planung der Anzeige starten; er läuft nach [`slice-v1-cli-egress-anzeige`](slice-v1-cli-egress-anzeige.md), wenn die Welle in der Reihenfolge der Tabelle arbeitet.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Messung zeigt, dass R1 ein Reload-Script mit atomarem Tabellenwechsel braucht, dessen Entwurf eine eigene Messreihe verlangt.
- `in-progress` → `open` (blockiert): die Messumgebung (Docker auf Colima mit `NET_ADMIN`) steht nicht zur Verfügung; die Evidence nennt den Grund (kein stilles Rot).

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Verification-Evidence und Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7); die Entscheidung des Projektinhabers (V1, V2 oder V3) ist im Gespräch belegt und in §7 vermerkt, die Folge-ADR trägt ihren Status.

## 6. Risiken und offene Punkte

- **Der Projektinhaber verbindet mit „Neu-Laden vom Host“ das Versprechen „der Agent kann sich nicht selbst freischalten“**, das mit passwortlosem `sudo` nicht einlösbar ist (siehe Einordnung). — **Ausgang:** offen bis zur Entscheidung; die Vorlage nennt es ausdrücklich; der Wortlaut (Stelle 1) weist es im Lastenheft aus.
- **Das offene Fenster beim erneuten Script-Lauf (M2)** könnte Verkehr ohne Filter durchlassen. — **Ausgang:** offen bis zur Messung; erwartet: Fenster in Millisekunden, mit atomarem Tabellenwechsel vermeidbar; die Folge-ADR nennt den Messwert.
- **Die Messung unter Colima gilt nur dort** (Docker 29.5.2 / Colima 0.10.3 / arm64 sind der Messstand der Vorgänger-Messungen; hier neu zu nennen). — **Ausgang:** weiter offen; Grenze in Handbuch und ADR benannt.
- **Die Handbuch-Aussage „neu starten“ ist falsch** (M4). — **Ausgang:** offen bis zur Messung; eingetreten: Handbuch-Korrektur in diesem Slice (kein eigener Slice).
- **Der erste Programmaufruf im Container bricht die Eigenschaft aus [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)** („ohne ein Programm im Container“) für die Lesebefehle nicht, aber ändert den Port-Zuschnitt. — **Ausgang:** offen; die Spezifikation grenzt den Satz im Spezifikations-Slice ein (siehe Wortlaut-Nachzug).

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt: Messergebnis M1–M4, Entscheidung des Projektinhabers (V1/V2/V3), Folge-ADR, Folge-Slices, Steering-Loop-Eintrag, Register, Risiken.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `docs/plan/` und `docs/user/` (DOK, Greenfield) laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); gemessen wird in Wegwerf-Containern, nicht in einer Sub-Area. Die Schwelle (zwei von drei Achsen) ist laut Modus-Deklaration erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md) (offen, ein Beleg; eine Messung auf Colima ohne CI-Gegenprobe ist dieselbe Lage — die Grenze wird ausgewiesen) und [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (verkörpert; die Aussagen zu Dateirechten oben sind aus Dockerfile und Doku gelesen und als „nicht gemessen“ markiert, nicht als Tatsache). Evidenz-/Diskrepanz-Risiko: mittel (mehrere Aussagen warten auf M1).

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
