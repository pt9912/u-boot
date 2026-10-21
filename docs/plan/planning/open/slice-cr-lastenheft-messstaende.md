# Slice CR: Lastenheft ohne Messstände — Messumgebungen in [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) und [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) als Zusage-Grenze statt als „gemessen: Colima, arm64“ (Lastenheft, nächste freie Version nach 0.19.0)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle, Release-Zuordnung: nach dem Release-Cut `v0.12.0`. Der Slice hält den Cut nicht auf und hat keine Closure-Bedingung, die von seiner DoD verschieden wäre. Er geht aus dem Review-Befund F-15 und der Auflage A-6 des Verification-Reports von [`slice-v1-devcontainer-ollama-amd64-messung`](../done/slice-v1-devcontainer-ollama-amd64-messung.md) hervor und setzt die Regel „Specs sind Zielbild“ fort, die [`slice-cr-lastenheft-zielbild`](../done/slice-cr-lastenheft-zielbild.md) (Prozesszustand aus der Historie) und [`slice-v2-spec-zielbild`](../done/slice-v2-spec-zielbild.md) (Messstände aus der Spezifikation) für ihre Stratum begonnen haben.

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Katalog-Feature `ollama`), [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Container-Runtime im Sandbox-Devcontainer), Lastenheft [§16 Historie](../../../../spec/lastenheft.md#16-historie); [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) (Entscheidung 10: `amd64` nicht verifiziert, bleibt unverändert).
Anlass: Review-Finding F-15 ([`2026-10-19-slice-v1-devcontainer-ollama-amd64-messung`](../../../reviews/2026-10-19-slice-v1-devcontainer-ollama-amd64-messung.md)) und Verification-Auflage A-6; die Vorgabe des Projektinhabers (2026-10-19): das Lastenheft trägt keine Messstände, das Handbuch trägt den Ist-Zustand.

**Berührte Spec-Stellen:** — (Lastenheft ist Rang 1 und ändert Anforderungstexte und seine Historie; `SPEC-*` und `ARC-*` bleiben unverändert, gelesen wird in §3). Der Verweis zeigt aufwärts: Das Lastenheft nennt diesen Slice nie.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Das Lastenheft nennt in den Anforderungstexten von [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) und [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) keine Messumgebung („Docker 29.5.2, Colima 0.10.3, arm64“, „Stand der Messung“, „gemessen: Colima, arm64“) mehr als Beleg der Zusage, sondern nur noch die Zusage und ihre Grenzen; wo die Messumgebung heute die Zusage-Menge bestimmt („die gemessenen Hosts“, „andere Architekturen und Container-Engines als gemessen“), steht die Menge im Wortlaut. **Der Wortlaut braucht die Freigabe des Projektinhabers**; dieser Slice ist nur der Plan, kein Text.

**Zusage-Grenzen bleiben.** Ein Satz „nicht zugesagt: …“ ist Vertrag (Zielbild der Grenze), kein Messstand, und bleibt im Sinn: „nicht zugesagt sind GPU-Beschleunigung, ein Anschluss an einen Ollama-Dienst des Hosts, andere Architekturen …, das Überstehen des Neubaus des Images sowie Auswahl oder Verwaltung von Modellen durch `u-boot`“; „Docker Desktop und Podman als Host-Engine sind für `docker` nicht zugesagt“; „BuildKit-Erweiterungen … und `docker compose` sind nicht zugesagt“. Sie nennen, was **nicht** zugesagt ist; die Messumgebung gehört ins Handbuch (Ist-Zustand, dort bereits geführt) und in die ADR.

### Gefundener Stand (gelesen am 2026-10-19; Suche nach „gemessen“, „Colima“, „arm64“, „amd64“, „nicht verifiziert“, „nicht gemessen“, „Messstand“, „Docker Desktop“, „29.5“, „29.9“, „0.40“, „0.10.3“, „Ubuntu“, „Messung“ über `spec/lastenheft.md`; die Treffer-Absätze gelesen)

Die Fundstellen mit Klasse (Klassen wie in [`slice-v2-spec-zielbild`](../done/slice-v2-spec-zielbild.md): **A** Messstand/Messumgebung, wird gestrichen oder in eine Zusage-Menge übersetzt; **B** Zusage-Grenze, bleibt; **C** Messbedingung einer Anforderung, bleibt; **D** Historie-Chronik):

| Fundstelle (`spec/lastenheft.md`) | Aussage (gelesen) | Klasse | Einordnung |
|---|---|---|---|
| Zeile 409, [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Absatz `ollama`) | „Gemessen ist das auf Docker 29.5.2, Colima 0.10.3, arm64, Ollama 0.40.6; nicht zugesagt sind GPU-Beschleunigung, ein Anschluss an einen Ollama-Dienst des Hosts, andere Architekturen und Container-Engines als gemessen, das Überstehen des Neubaus des Images sowie Auswahl oder Verwaltung von Modellen durch `u-boot`.“ | **A** (erste Satzhälfte) und **B** (Aufzählung), mit Verknüpfung | Die Messumgebung bestimmt hier die Zusage-Menge („andere Architekturen … als gemessen“); der neue Wortlaut nennt die Menge ausdrücklich. **Entscheidung des Projektinhabers (unten, Frage 1):** ob `amd64` und welche Engines zugesagt sind; [`slice-v1-devcontainer-ollama-amd64-messung`](../done/slice-v1-devcontainer-ollama-amd64-messung.md) hat ausdrücklich keine Zusage für `amd64` gemacht (die Messung ist eine Einzelmessung ohne CI-Lauf) |
| Zeile 463, [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Punkt `podman`, `podman build`) | „(gemessen mit `docker.io/library/alpine:3` und `debian:stable-slim`); das gilt auf den gemessenen Hosts (Stand der Messung: Docker 29.5.2 auf Colima 0.10.3, arm64, VM Ubuntu 24.04 mit `kernel.apparmor_restrict_unprivileged_userns=0`; andere Architekturen und Engines nicht gemessen).“ | **A** | Die Zusage-Menge „gemessene Hosts“ wird im Wortlaut genannt (Frage 1); die Messbilder (`alpine:3`, `debian:stable-slim`) sind Messstand, kein Vertrag; der Satz „BuildKit-Erweiterungen … nicht zugesagt“ bleibt (**B**) |
| Zeile 465, [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Punkt `podman`, Start unter Docker) | „… auf Hosts, die verschachtelte User-Namespaces erlauben (gemessen: Colima, arm64, mit `kernel.apparmor_restrict_unprivileged_userns=0` in der VM als Voraussetzung) und wenn die Devices `/dev/fuse` und `/dev/net/tun` … vorhanden sind; mit dem Wert `1` … endet das Startscript mit Exit-Code `11` …“ | **A** (Klammer „gemessen: Colima, arm64“) und **C** (Voraussetzung `…=0`, Exit-Code `11`) | Die Voraussetzung (Wert `0` des Kernel-Parameters) bleibt als Bedingung der Zusage; nur der Zusatz „gemessen: Colima, arm64“ entfällt oder wird zur Menge |
| Zeile 465, derselbe Absatz, Satz zu `docker` | „Bei `docker` startet es unter Docker auf Linux und unter Colima (Linux-VM); Docker Desktop und Podman als Host-Engine sind für `docker` nicht zugesagt.“ | **B** | Zusage-Menge im Wortlaut (aus [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) Version 0.7.0); bleibt |
| Zeile 1563, Historie 0.7.0 | „Zusage für `nestedRuntime: docker` auf die gemessenen Hosts eingegrenzt: Docker auf Linux und Colima (Linux-VM). Docker Desktop ist nicht mehr zugesagt, weil es nicht gemessen ist …“ | **D** | Chronik einer Änderung; ob Historie-Zeilen früherer Versionen Messangaben tragen dürfen, entscheidet der Projektinhaber (Frage 2) |
| Zeile 1568, Historie 0.12.0 | „… gemessen auf Docker 29.5.2 unter Colima 0.10.3 (arm64, VM Ubuntu 24.04 mit `kernel.apparmor_restrict_unprivileged_userns=0`) …“ | **D** | wie 0.7.0 |
| Zeile 1569, Historie 0.13.0 | „… (gemessen: Docker 29.5.2 auf Colima 0.10.3, `docker run` Exit-Code `127`, `docker start` Exit-Code `1`) …“ | **D** | wie 0.7.0 |
| Zeile 1571, Historie 0.15.0 | „Katalog-Feature `ollama` … für die gemessenen Hosts und Versionen; nicht zugesagt: GPU, Anschluss an einen Ollama des Hosts, Modellverwaltung.“ | **D** | wie 0.7.0 |
| Zeile 1079, [`LH-NFA-PERF-001`](../../../../spec/lastenheft.md#lh-nfa-perf-001--schnelle-cli-antwort) | „(gemessen ohne Docker-Kommunikation, Kaltstart)“ | **C** | Messbedingung der Anforderung (Leistungsziel), keine Messumgebung eines Hosts; bleibt |
| Zeile 424, [`LH-FA-DEV-004`](../../../../spec/lastenheft.md#lh-fa-dev-004--benutzerrechte) (Benutzerrechte), Absatz zur UID | „(z. B. `501` unter macOS mit Colima)“ | — | Beispiel, kein Messstand; bleibt |
| Zeilen 1484 und 1559, [`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung), Historie 0.3.5 | `.deb`/`.rpm` für amd64 und arm64 als Release-Assets | **B** | Zusage-Menge der Paketierung, keine Messung; bleibt |

### Offene Fragen an den Projektinhaber (vor dem Wortlaut)

1. **Zusage-Menge.** Sagt das Lastenheft für `ollama` und für `podman build` die Architektur `arm64` unter Colima zu, oder auch `amd64` auf Linux (gemessen in einer Einzelmessung, ohne CI-Lauf; [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) Entscheidung 10 bleibt „nicht verifiziert“)? Ohne diese Entscheidung kann der Wortlaut die Messumgebung nicht durch die Menge ersetzen.
2. **Historie.** Tragen die Zeilen 0.7.0, 0.12.0, 0.13.0 und 0.15.0 ihre Messangaben weiter (Chronik einer Version), oder bekommt eine neue Zeile der nächsten freien Version (nach 0.19.0) die Aussage „die Messangaben der Anforderungstexte entfallen“, ohne die älteren Zeilen zu ändern (Vorbild: Zeile 0.17.0 von [`slice-cr-lastenheft-zielbild`](../done/slice-cr-lastenheft-zielbild.md))?
3. **Der Wortlaut selbst.** Der CR legt einen Alt/Neu-Vorschlag je Stelle vor; die Freigabe erfolgt je Stelle durch den Projektinhaber (Muster: Wortlaut-Vorschlag in [`slice-cr-lastenheft-zielbild`](../done/slice-cr-lastenheft-zielbild.md)).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Wortlaut selbst** — braucht die Freigabe des Projektinhabers; der Plan schlägt keinen Text vor, der als freigegeben gelesen werden kann.
- **Eine Zusage für `amd64`** — ist die Frage 1; eine Zusage ohne CI-Lauf mit sichtbarem Pass widerspräche der Linie von [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) Entscheidung 10; ändern würde sie eine Folge-ADR.
- **Handbuch, CHANGELOG, Spezifikation, ADRs** — das Handbuch trägt die Messumgebung als Ist-Zustand bereits; die Spezifikation trägt nach [`slice-v2-spec-zielbild`](../done/slice-v2-spec-zielbild.md) keine; die ADRs sind `Accepted` und bleiben unverändert.
- **Zusage-Grenzen** („nicht zugesagt …“, „nicht gemessen“ als Grenze einer Zusage) — Vertrag, nicht Messstand (siehe Klasse **B**).
- **Messstände anderer Dokumente im Planning-Baum** (Slices, Wellen, ADRs) — Zeitdokumente, bleiben.

## 2. Definition of Done

- [ ] **Wortlaut freigegeben:** Der Projektinhaber hat die Fragen 1 bis 3 entschieden; die Stellen der Klasse **A** in `spec/lastenheft.md` tragen keine Messumgebung mehr; jede Zusage-Menge steht im Wortlaut; die Klassen **B** und **C** sind unverändert; die Version (nächste freie Minor-Version zum Zeitpunkt des Commits; zuletzt `0.19.0`) und die neue Historie-Zeile nennen die Änderung ohne Messangabe.
- [ ] **Abgleich nach unten:** Handbuch (Abschnitte *Ollama (optional)* und *Container-Runtime*) und CHANGELOG sagen über die Zusage-Menge dasselbe wie der neue Wortlaut; `make doc-trace` und `make docs-check` melden keinen Befund; jede Aussage des Plans über das Lastenheft ist am Absatz gelesen (`AGENTS.md` §3.10).
- [ ] `make gates` grün; Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`); Closure-Notiz mit Steering-Loop-Lerneintrag; Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist; jedes Risiko aus §6 trägt einen Ausgang; die drei Paarungen sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` (Absätze der Zeilen 409, 463 und 465; Versionsfeld im Kopf; §16 Historie) | update | Streichen der Klasse **A**, Zusage-Menge im Wortlaut, neue Historie-Zeile; erst nach Freigabe |
| `docs/user/benutzerhandbuch.md`, `CHANGELOG.md` | update, falls die Zusage-Menge im Wortlaut vom heutigen Handbuch-Text abweicht | Abgleich nach unten; Handbuch kennungsfrei (`AGENTS.md` §5 Regel 6) |
| diese Slice-Datei, `open/README.md`, Register, Roadmap bei Start und Closure; `docs/reviews/` | update; neu | Lifecycle, Review und Verifikation |
| `spec/spezifikation.md`, `spec/architecture.md`, `docs/plan/adr/`, `internal/**`, `Makefile`, `.github/workflows/` | unverändert | Auftrag; das Lastenheft ist Rang 1, die Spezifikation trägt nach [`slice-v2-spec-zielbild`](../done/slice-v2-spec-zielbild.md) keinen Messstand |

## 4. Trigger

**Start** (`open` → `next` → `in-progress`): `Verantwortlich:` gesetzt; der Release-Cut `v0.12.0` ist geschnitten; der Projektinhaber hat die Fragen 1 und 2 entschieden (Wortlaut-Freigabe folgt je Stelle im Slice); kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Zusage-Menge verlangt Änderungen über [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) und [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) hinaus (weitere Anforderungen mit Messbezug); dann trennt der Architect die Anforderungen.
- `in-progress` → `open` (blockiert): Der Projektinhaber gibt den Wortlaut nicht frei oder verlangt eine Zusage für `amd64`, die eine Folge-ADR braucht; Befund in §7.

## 5. Closure-Trigger

DoD vollständig, Review-Report und Verification-Report liegen vor, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **Das Streichen der Messumgebung weitet die Zusage unbemerkt aus** („die gemessenen Hosts“ wird zu „alle Hosts“). — **Ausgang:** weiter offen bis zum Wortlaut; jede Zusage-Menge steht im Wortlaut ausdrücklich, die Grenzen bleiben ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md), 3×).
- **Eine Vertragsänderung nach dem Review wird nicht mehr gelesen** ([`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md)). — **Ausgang:** weiter offen; Nachlauf-Review über den Diff nach der letzten Fix-Runde.
- **Das Streichen widerspricht der Historie** (Zeilen früherer Versionen nennen die Messumgebung weiter). — **Ausgang:** weiter offen bis zur Frage 2.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `spec/` (Lastenheft, DOK), `docs/user/` und `CHANGELOG.md` (DOK) laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); alle Greenfield.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (3×, Schwelle erreicht, Ausgang beim Projektinhaber), [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md) (offen, vier Belege), [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (verkörpert in `AGENTS.md` §3.10). Evidenz-/Diskrepanz-Risiko: mittel (Zusage-Menge).

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
