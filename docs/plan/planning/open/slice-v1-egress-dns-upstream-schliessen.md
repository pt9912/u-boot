# Slice v1: Egress — Pfad zum Upstream-Resolver für Prozesse ohne `sudo` schließen oder messen (`LH-FA-DEV-008`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Der Slice ist ein eigener Slice **nach** [`welle-sandbox-egress-ollama`](../done/welle-sandbox-egress-ollama.md) (Entscheidung des Projektinhabers, 2026-10-18) und hält deren Closure und das Release nicht auf; er hat keine Closure-Bedingung, die von seiner DoD verschieden wäre. Reihenfolge: nach der Closure jener Welle.

**Bezug:** [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Egress-Restriktion: „Allowlist für ausgehenden Netzwerkverkehr“; Zielbild, Code folgt), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil: Schadensbegrenzung); [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) (Mechanismus `dnsmasq` + nftables-Set, `Accepted`), [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md) (Messbefund M1 a, `Accepted`).
Anlass: Befund M1 (a) aus [`slice-v1-sandbox-egress-neuladen-entwurf`](../done/slice-v1-sandbox-egress-neuladen-entwurf.md); Entscheidung des Projektinhabers (2026-10-18): Code zum Ziel, eigener Folge-Slice, das Lastenheft wird nicht abgeschwächt.

**Berührte Spec-Stellen:** [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) *Grenze* (nur der Teil zur Schreibbarkeit der Statusdatei; dessen Nachführung gehört dem Spezifikations-Slice der Welle, nicht diesem Slice), [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) (Allowlist; der Architect liest, ob ein DNS-Satz dort steht). Der Slice **liest** sie und ändert keine; ist nach der Messung eine Spec-Änderung nötig, benennt der Architect die Stelle, und ein eigener Spezifikations-Slice folgt dem Code.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Prozess im Sandbox-Container ohne `sudo` bekommt über den Upstream-Resolver keine Adresse mehr für einen Namen außerhalb der Allowlist — oder die Messung zeigt, dass und warum das mit dem Mechanismus von [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) nicht geht, und das Ergebnis ist ausgewiesen. **Messung zuerst**, Umsetzung nur nach dem Messbefund.

### Gefundener Stand (gelesen am 2026-10-18)

- **Die Regeln lassen DNS zu jedem Upstream für alle Prozesse zu.** `internal/hexagon/application/templates/devcontainer/egress-init.sh.tmpl`: `upstreams` kommt aus den IPv4-`nameserver`-Zeilen der gesicherten `/etc/resolv.conf.u-boot-orig` (Zeile 36); die Schleife in den Zeilen 63 bis 66 legt je Upstream `ip daddr "$ns" udp dport 53 accept` und `… tcp dport 53 accept` in die Kette `output` — ohne Bedingung an Benutzer oder Prozess. Der lokale `dnsmasq` fragt nur die Upstreams für erlaubte Namen (`server=/$h/$ns`, Zeile 52; `no-resolv`, Zeile 47); der Container zeigt über `/etc/resolv.conf` auf `127.0.0.1` (Zeile 72). Ein Prozess, der den Upstream direkt anspricht, umgeht den lokalen Server.
- **Gemessen** (Messprotokoll [`rohdaten-02-m1.txt`](../observations/BEO-DOK/messprotokoll-ausserhalb-des-repos/evidence/messung-slice-v1-sandbox-egress-neuladen-entwurf/rohdaten-02-m1.txt), Abschnitt „DNS direkt an den Upstream-Resolver 192.168.5.1 (ohne sudo)“): `example.org` löst mit zwei A-Einträgen auf; die Verbindung zur Adresse läuft in den Timeout (`Connection timed out after 5006 milliseconds`). Je ein Lauf, Messstand Docker 29.5.2 auf Colima 0.10.3 (arm64). Ein Datenkanal über DNS ist nicht gemessen.
- **Beleg für verschachtelte Container** (Messung des Cut-Slice [`slice-v1-release-cut-v0.12.0`](../done/slice-v1-release-cut-v0.12.0.md), Abschnitt *Umsetzung (Implementer)*, *Befund zum roten Test*, 2026-10-19, Linux `amd64`, native Docker Engine 29.9.1, zwei Läufe): Der Docker-Test `TestE2E_SandboxDevcontainer_NestedPodman_BuildEgress` (`internal/e2e`) ist rot: In `STEP 3/3: RUN ! getent hosts example.com` löst der Name auf (nur IPv6), weil der verschachtelte Podman-Container als Rückfall `8.8.8.8` fragt, ein Upstream des Docker-Hosts, den die Regeln aus Zeilen 63 bis 66 für jeden Prozess zulassen; die TCP-Verbindung bleibt gesperrt; mit `--dns 9.9.9.9` ist der Test grün. Dieselbe Ursache wie oben, gelesen und gemessen, nicht für Colima oder die CI.
- **Was heute dazu gesagt wird:** [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) Alternative C2 („nicht erlaubte Namen sind nicht auflösbar“) gilt für den lokalen `dnsmasq`; das Handbuch (Abschnitt *Internet-Beschränkung (Egress)*, Absatz *Grenze der Namensauflösung*) weist die direkte Abfrage aus. [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) bleibt unverändert: Der Guardrail-Satz beschreibt die Aufhebung durch Capability, der offene DNS-Pfad ist eine Lücke des Codes gegenüber dem Zielbild, keine Abschwächung des Lastenhefts.

### Messfragen (Arbeitsplan; Docker auf Colima, Container wie im Messprotokoll der Vorgänger-Messung; der Architect entscheidet den Weg)

| # | Frage | Warum |
|---|---|---|
| D1 | Läuft der `dnsmasq` unter einem Benutzer, den eine nftables-Bedingung (`meta skuid`) von `vscode` unterscheidet? Welche Benutzer-IDs tragen `dnsmasq` und der Nutzer im gemessenen Image? | Voraussetzung für eine Regel „DNS zum Upstream nur für den `dnsmasq`“; der Benutzername des `dnsmasq` ist hier nicht gelesen, nur gemessen werden kann er |
| D2 | Mit einer solchen Regel (Upstream-DNS nur für den `dnsmasq`, andere Prozesse verworfen): löst ein Prozess ohne `sudo` einen Namen außerhalb der Liste direkt noch auf? Lösen erlaubte Namen weiter auf (Wirkung auf `api.github.com` und die Default-Allowlist)? | Wirksamkeit und Nebenwirkung |
| D3 | Wege um die Regel: TCP statt UDP, ein anderer Upstream (IPv6 ist gesperrt, Zeilen 67 bis 69), DNS über eine erlaubte Adresse, eine andere Portnummer; was bleibt offen? | Grenze der Lösung ehrlich benennen |
| D4 | Bleibt das Startscript beim zweiten Lauf und beim Neuaufbau wie bisher lauffähig (Messvariante aus dem Vorgänger-Messprotokoll)? | keine Regression der Start-Eigenschaften |
| D5 | Offen für den Architect: ein Datenkanal über die erlaubten Namen (Anfragen an erlaubte Namen werden an den Upstream weitergeleitet) — vorhanden, ja oder nein, mit welchem Durchsatz? | nur wenn D2 den direkten Pfad schließt; sonst Grenze im Handbuch |

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Entziehen von `sudo`** (Option S aus [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md)) — ändert Benutzerrechte und alle Startscripte; mit `sudo` bleibt die Restriktion ein Guardrail ([`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)), das Ziel hier ist nur der Weg *ohne* `sudo`.
- **Neu-Laden der Liste im laufenden Container** — entschieden gegen (R0, [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md)).
- **Eine Restriktion außerhalb des Containers** (Proxy, Netz-/DNS-Ebene) — eigener Vorgang ([`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) Alternative D).
- **Die Nachführung von [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) *Grenze* (Schreibbarkeit der Statusdatei)** — Spezifikations-Slice der Welle ([`slice-v2-spec-sandbox-egress-ollama`](../done/slice-v2-spec-sandbox-egress-ollama.md)).
- **Änderung des Lastenhefts** — nach dem Stand oben nicht erwartet (Zielbild, Code folgt); zeigt D3, dass der Pfad grundsätzlich nicht schließbar ist, geht der Befund an den Projektinhaber, und das Handbuch bleibt bei der ausgewiesenen Grenze.
- **Docker Desktop, Linux-Host ohne VM, `amd64`** — nicht gemessen, nicht zugesagt; jede Aussage nennt den Messstand.

## 2. Definition of Done

- [ ] **Messung D1 bis D4** (D5 falls D2 den Pfad schließt) gemessen, mit Messstand, Skripten und Rohdaten im Repo (Ablage wie das Messprotokoll des Vorgängers); je Frage ein Ergebnis oder die Angabe „nicht gemessen“ mit Grund.
- [ ] **Entscheidungs-ADR und Umsetzung:** Eine Folge-ADR zu [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) (`Accepted` bleibt unverändert) trägt den Messbefund und den gewählten Weg; bei einem Weg, der schließt: Änderung am Startscript-Template, Golden Case und Docker-Test (`//go:build docker`), der den direkten Upstream-Pfad ohne `sudo` abweist und die erlaubten Namen weiter auflöst; Handbuch (Absatz *Grenze der Namensauflösung*) und `CHANGELOG.md` auf den Ist-Stand. Bei einem Weg, der nicht schließt: die Grenze steht ausgewiesen im Handbuch, und der Projektinhaber hat entschieden.
- [ ] **Der Test `TestE2E_SandboxDevcontainer_NestedPodman_BuildEgress` ist abgenommen:** Er wird mit dieser Änderung grün oder wird mit Begründung angepasst (Ausgang in §7); bis dahin bleibt `make test-docker` auf einem Host, dessen Upstream-Resolver ein Rückfall-Nameserver von Podman ist, rot.
- [ ] **Die Spec-Stelle**, die nach dem Befund nachzuführen ist (falls eine), ist benannt; die Nachführung selbst entsteht als eigener Spezifikations-Slice nach dem Code.
- [ ] `make gates` grün; `make test-docker` grün, wenn Code geändert wurde.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: Handbuch, ADR-Index, `CHANGELOG.md`, soweit der öffentliche Vertrag berührt ist.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo mit Wellen von der nächsten Welle-Closure geprüft.

## 3. Plan (vor Code)

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; der Implementer ergänzt sie in seinem ersten Lauf und gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab (`AGENTS.md` §3.10). Jede Zeile ist eine Erwartung, keine gelesene Aussage über unveränderte Dateien.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Messprotokoll unter `docs/plan/planning/observations/BEO-DOK/messprotokoll-ausserhalb-des-repos/evidence/messung-slice-v1-egress-dns-upstream-schliessen/` (Skripte, Rohdaten, README mit Messstand und Voraussetzungen) | neu | D1 bis D5; im Repo statt außerhalb ([`BEO-DOK/messprotokoll-ausserhalb-des-repos`](../observations/BEO-DOK/messprotokoll-ausserhalb-des-repos/observation.md)) |
| `docs/plan/adr/NNNN-….md` und `docs/plan/adr/README.md` | neu / update | Folge-ADR zu [`ADR-0012`](../../adr/0012-devcontainer-egress-firewall.md) mit Messbefund und Weg (nächste freie Nummer) |
| `internal/hexagon/application/templates/devcontainer/egress-init.sh.tmpl` | update (nur wenn die Messung einen Weg trägt) | die Regeln zu Zeilen 63 bis 66 und ihre Voraussetzung (D1) |
| Golden Case und Test des erzeugten Scripts (`TestGenerateDevcontainer_SandboxEgress` und die Golden-Dateien; der Architect liest die Stellen) | update | das erzeugte Script ändert sich; Replay-Evidence nach `harness/replay.md` |
| `internal/e2e` (Docker-Test, Build-Tag `docker`) | neu / update | belegt die Zusage: direkter Upstream-Pfad ohne `sudo` abgewiesen, erlaubte Namen lösen auf; dazu `TestE2E_SandboxDevcontainer_NestedPodman_BuildEgress` (roter Test, Abnahme in der DoD) |
| `docs/user/benutzerhandbuch.md`, `CHANGELOG.md` | update | Absatz *Grenze der Namensauflösung* und die Aussage zum lokalen DNS-Server auf den Ist-Stand |
| diese Slice-Datei, `open/README.md`, Roadmap bei Start und Closure | update; Moves | Lifecycle |
| `docs/reviews/` | neu (Reviewer, Verifier) | Review und Verifikation |

## 4. Trigger

**Start** (`open` → `in-progress`, ggf. über `next`): `Verantwortlich:` gesetzt; [`welle-sandbox-egress-ollama`](../done/welle-sandbox-egress-ollama.md) ist geschlossen (Reihenfolge-Entscheidung des Projektinhabers); kein anderer Slice in `in-progress/` (WIP-Limit 1). Der Start hängt nicht am Release.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): D3 zeigt mehrere Umgehungswege, deren Schließung jeweils eine eigene Messreihe braucht.
- `in-progress` → `open` (blockiert): die Messumgebung (Docker auf Colima mit `NET_ADMIN`) steht nicht zur Verfügung; die Evidence nennt den Grund (kein stilles Rot).

## 5. Closure-Trigger

DoD vollständig, Review-Report und Verification-Evidence liegen vor, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7); die Entscheidung des Projektinhabers zum gewählten Weg (schließen oder Grenze ausweisen) ist im Gespräch belegt und in §7 vermerkt.

## 6. Risiken und offene Punkte

- **Die Regel gegen direkte DNS-Abfragen trifft auch den `dnsmasq` oder erlaubte Namen** (Auflösung bricht). — **Ausgang:** weiter offen bis zur Messung D2.
- **Eine Regel nach Benutzer trägt nur, solange der Nutzer ohne `sudo` bleibt;** mit `sudo` ist jede Regel aufhebbar. — **Ausgang:** weiter offen; Grenze im Handbuch und in der Folge-ADR zu benennen (Guardrail, [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)).
- **Ein Datenkanal über erlaubte Namen bleibt möglich** (nicht gemessen, D5). — **Ausgang:** weiter offen bis zur Messung.
- **Die Messung gilt nur für Colima/arm64.** — **Ausgang:** weiter offen; Grenze in Handbuch und Protokoll benennen.
- **Der Slice ändert das Startscript, dessen Startrennen aus [`slice-v1-egress-wildcard-pruefung`](../done/slice-v1-egress-wildcard-pruefung.md) nicht ausgeschlossen ist.** — **Ausgang:** weiter offen; D4 misst die Startrennen-Läufe wie der Vorgänger.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** voraussichtlich berührt werden `internal/hexagon/application` (APP, Greenfield; das Template liegt unter `templates/devcontainer/`), `internal/e2e` (E2E, Greenfield; ein Docker-Test) und `docs/` (DOK, Greenfield) laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); jede erfüllt die Schwelle aus zwei von drei Achsen laut Modus-Deklaration.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md) (offen; ein neuer Docker-Test mit Image-Bau ist ein weiterer Kandidat), [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md) (offen; ein Skip ist im CI-Aufruf ohne `-v` nicht von einem Pass zu trennen), [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (verkörpert; die Aussagen zum Template oben sind an den genannten Zeilen gelesen, die Benutzer-Frage des `dnsmasq` ist als ungelesen markiert) und [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (verkörpert; eine Handbuch-Bedingung zum DNS-Pfad braucht einen Test oder entfällt). Evidenz-/Diskrepanz-Risiko vor der Messung: mittel.

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
