# Slice v1: Cloud-Modelle von Ollama mit Claude Code — Hosts, Fehlertext und Verhalten bis zur Anmeldung ohne Konto messen (Katalog-Feature `ollama`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-ollama-cloud-modelle`](../welle-ollama-cloud-modelle.md) (Release-Nummer offen, Vorschlag `v0.16.0`) — erster Slice der Welle (Messung ohne Konto; ändert weder Lastenheft noch Produktcode). Ihm folgen die bedingte Messung mit Konto, die ADR, der Lastenheft-CR, der Code, die Spezifikation und der Release-Cut (Welle-Datei §4; Kennungen vorläufig).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Devcontainer-Features; Katalog-Feature `ollama`), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Hosts), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zugangsdaten nur als Umgebungsvariable); [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) (Ollama im Sandbox-Devcontainer; Muster der Messung; Re-Evaluierungs-Trigger „ein Modell braucht weitere Hosts“), [`ADR-0035`](../../adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) (Muster: Zusage nur für das Gemessene, Token nur als Umgebungsvariable).

Anlass: Beauftragung durch den Projektinhaber am 2026-10-19 (Roadmap, *Nächste Wellen*); Welle eröffnet am 2026-10-24.

**Berührte Spec-Stellen:** [`SPEC-047`](../../../../spec/spezifikation.md#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen) (Feature `ollama`), [`SPEC-044`](../../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) (Zugangsdaten des Agenten) — der Slice **liest** sie und ändert keine; die Nachführung folgt im Spezifikations-Slice der Welle. Das Lastenheft ändert nur der CR. Der Verweis zeigt **aufwärts**: Lastenheft und Spezifikation nennen diesen Slice nie.

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

**Ziel:** Mit gemessenen Belegen **ohne Konto** steht fest, was ein Cloud-Modell von Ollama im Sandbox-Devcontainer an Hosts braucht, welchen Fehlertext der Nutzer hinter der Egress-Liste sieht, ob ein Host-Eintrag `ollama.com` diesen Fehler behebt oder zum nächsten Hindernis führt, und wie sich `ollama launch claude` und `ollama signin` bis zur Anmeldung verhalten. Ergebnis: ein Messprotokoll (Methode wie O4 in [`slice-v1-devcontainer-ollama-messung`](../done/slice-v1-devcontainer-ollama-messung.md): protokollierender Resolver, Standard-Allowlist plus Kandidaten-Hosts) und ein Messplan für den Lauf mit Konto. **Was ohne Konto nicht messbar ist, steht im Protokoll als „nicht gemessen“ und im Plan als Entscheidung und Zugang des Projektinhabers** (Welle-Datei §5, Frage 2): die Anmeldung, der Lauf eines Cloud-Modells, das Zusammenspiel mit Claude Code und die Hosts nach der Anmeldung. Zugesagt wird später nur, was gemessen ist.

### Gefundener Stand (gelesen am 2026-10-24)

- **Handbuch.** `docs/user/benutzerhandbuch.md`, Abschnitt *Ollama (optional)*, Zeilen 2023 bis 2029: Der Satz „Cloud-Modelle sind nicht enthalten“ nennt als gemessen `ollama pull …-cloud` (`success`, in `ollama list`), den Fehler von `ollama run` (`ollama.com` nicht auflösbar) und `ollama signin` ohne Konto (`account unavailable`), Ollama `0.40.6`, Linux `amd64`, hinter der Egress-Beschränkung; Anmeldung mit Konto und Lauf mit Konto stehen als nicht gemessen. Das Protokoll dieses Slice **wiederholt** diese Messung nicht als Neuigkeit, sondern erweitert sie um die Fragen unten und nennt, was von der früheren Messung bestätigt ist.
- **Hosts.** `egressFeatureHosts` (`internal/hexagon/application/devcontainer_egress.go`, Zeilen 92 und 93): `claude-code` führt `registry.npmjs.org` und `api.anthropic.com`, `ollama` führt `registry.ollama.ai` und den Blob-Host; der Doc-Kommentar (Zeilen 78 bis 80) hält fest, dass `ollama.com` für lokale Modelle nicht nötig ist (gemessen in O4, zwei bis vier Anfragen beim Start des Servers).
- **Katalog und Meldungen.** `featureCatalogue` (`internal/hexagon/application/devcontainer_features.go`): `ollama` Zeile 279 (`dockerfileBlock: true`), `claude-code` Zeile 273 (OCI-Quelle); `ollamaWarnings` (Zeile 535): ohne Sandbox-Profil „has no effect“, mit Profil die Ausweisung der Prüfsummen-Grenze. Die Version von Ollama steht als `OLLAMA_VERSION` im erzeugten Dockerfile (Pin laut [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md), Handbuch: `0.40.6`).
- **Zugangsdaten.** `agentCredentialEnv` (`internal/hexagon/application/devcontainer_sandbox.go`, Zeile 441): für `ollama` keine Variable. Ob eine existiert, ist eine Messfrage (C4), keine Annahme.
- **Messmuster.** `tools/messung/ollama-amd64/` (gelesen: `lib.sh` mit `redact` und `rohdaten`, `05-laufzeit.sh`, `inner-05.sh` mit den Modi `discover` und `restricted`): der protokollierende `dnsmasq` (`log-queries`) vor dem Aufruf, die Allowlist des erzeugten Containers im Modus `restricted`. Die `redact`-Schwärzung ist formgebunden (vorsignierte URLs, `Authorization`/`Bearer`, `X-Amz-…`); sie ist **keine** Schwärzung für ein Token in beliebiger Form ([`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md)).

- **Gemessen (2026-10-24, Implementer).** Protokoll [`README.md`](../../../messungen/slice-v1-ollama-cloud-messung/README.md), Skripte `tools/messung/ollama-cloud/`: C0 bis C7 in zwei Läufen je Modus (`discover`, `restricted`, `allow`), ohne Konto. Die gepinnte Version `0.40.6` kennt `launch` und das Ziel `claude` (keine Pin-Hebung). Im Protokoll stehen Handvarianten (`devcontainer up` der Dev-Container-CLI statt `u-boot devcontainer up`, Platzhalter für `claude`, Lauf 2 mit vorhandenem Volume) und die Liste „nicht gemessen“ mit dem Messplan für den Lauf mit Konto.

**Messfälle (alle ohne Konto, ohne Token, Linux `amd64`, lokale Docker-Engine, erzeugtes Projekt mit `ollama`, `claude-code` und `node`):**

- **C0 Unterbefehl:** Kennt die gepinnte Ollama-Version `launch` und das Ziel `claude` (`ollama launch --help`, `ollama --version`)? Wenn nein, ist das Ergebnis, dass die Welle eine Pin-Hebung braucht ([`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) Entscheidung 9: Hebung per Slice); die Messung läuft dann nicht mit einer Version, die `u-boot` nicht ausliefert, sondern wird mit der benannten Version als Variante ausgewiesen ([`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](../observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md)).
- **C1 Hosts (Entdeckung):** `ollama pull`, `ollama run` und `ollama launch claude` mit einem Cloud-Modell im Modus `discover`; Namen aus dem Resolver-Log, mit Zahl der Anfragen und ohne Antwort-Inhalte.
- **C2 Fehlertext (Egress):** dieselben Aufrufe im Modus `restricted` mit der erzeugten Liste: Wortlaut der Meldung, Exit-Code, Dauer. **C2b:** mit dem Eintrag `ollama.com` in `devcontainer.sandbox.egress.allow` (erzeugt über `u-boot`, nicht von Hand im Container): ändert sich der Fehler, und wie lautet der nächste? „Behebt den Fehler“ heißt hier nur: die Meldung ändert sich; ein Lauf wird ohne Konto nicht erwartet und nicht behauptet.
- **C3 `ollama launch claude`:** Was der Aufruf bis zum Hindernis tut — lädt oder startet er selbst etwas (Hosts, Dateien, Ausgabe), findet er das installierte `claude`, verlangt er Anmeldung vor dem Start?
- **C4 `ollama signin` und Zugangs-Variable:** Ausgabe bis zum Hindernis, Hosts, Dateien; ob `ollama --help`, `ollama launch --help` oder die Ausgabe der Aufrufe eine Umgebungsvariable für den Zugang nennen (gelesen aus dem Binary des Containers, nicht aus dem Netz).
- **C5 `OLLAMA_NO_CLOUD`:** Wirkung auf `pull`, `run` und `launch` ohne Konto (Meldung, Hosts).
- **C6 Dateien:** welche Pfade die Aufrufe im Container anlegen oder ändern (Pfad, Größe, **kein Inhalt**), und ob einer davon im Modell-Volume liegt ([`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer): Zugangsdaten nicht im Workspace-Volume).
- **C7 Standardprofil:** `generate devcontainer` mit `ollama` und `claude-code` ohne Sandbox-Profil: erzeugte Meldung und erzeugte Dateien (kein Container-Lauf nötig, da `u-boot` dort kein Ollama installiert); damit der Nutzerpfad im Standardprofil für den Validator beschrieben ist.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Anmeldung mit Konto, Lauf eines Cloud-Modells, Zusammenspiel mit Claude Code nach der Anmeldung, Hosts nach der Anmeldung** — ohne Konto nicht messbar; Entscheidung und Zugang des Projektinhabers (Welle-Datei §5, Frage 2), Folge-Slice `slice-v1-ollama-cloud-messung-mit-konto` (bedingt, Kennung vorläufig, noch keine Datei). Das Protokoll führt sie als „nicht gemessen“, der Plan des Folge-Slice hält den Messplan.
- **Jeder Zugang, jedes Token und jede Umgebungsvariable mit Zugangsdaten** — der Slice verwendet keinen; Token gehören ausschließlich in die Umgebung des Hosts beim Lauf des Folge-Slice, nie in Datei, Skript, Protokoll oder Chat ([`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md)).
- **Lastenheft, Spezifikation, Handbuch, Produktcode, ADR** — eigene Slices der Welle (Welle-Datei §4); die Zusage beruht auf der Messung, die Folge-ADR entscheidet, der CR schreibt, der Code folgt dem Lastenheft. Der Slice liefert Skripte und Protokoll.
- **Pin-Hebung von Ollama** — falls C0 sie verlangt: eigener Slice nach [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) Entscheidung 9; der Slice weist es als Befund mit Ausgang aus und hebt den Pin nicht.
- **Andere Plattformen und Engines** (`arm64`, Docker Desktop, Podman als Host-Engine) — Zusage-Menge ist Linux `amd64` mit lokaler Engine; ein anderer Host steht im Protokoll und begründet keine Zusage.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **Messung ohne Konto** ([`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)): C0 bis C7 gemessen **am erzeugten Projekt** (Ergebnis von `generate devcontainer` mit dem ausgelieferten Binary, nicht an einer Handvariante; weicht die Version von Ollama ab, steht sie als Variante im Protokoll). Je Aussage Umgebung, Datum und Zahl der Läufe; **bedingungsgleich formuliert:** „ohne Konto gemessen, Erfolgspfad nicht gemessen“, nicht „Cloud-Modelle laufen/laufen nicht“. *Beleg:* Skripte unter `tools/messung/ollama-cloud/`, Protokoll mit Rohdaten unter `docs/messungen/slice-v1-ollama-cloud-messung/` ([`harness/conventions.md`](../../../../harness/conventions.md#ablage-von-messprotokollen-und-skripten)).
- [x] **Messplan für den Lauf mit Konto und Liste „nicht gemessen“:** Das Protokoll führt die ungemessenen Punkte (Anmeldung, Lauf eines Cloud-Modells, Zusammenspiel mit Claude Code, Hosts nach der Anmeldung, Name einer Zugangs-Variable, falls C4 keine nennt) je mit dem Aufruf, den der Folge-Slice fährt, und der Bedingung, unter der er läuft (Zugang des Projektinhabers als Umgebungsvariable des Hosts). Das ist die Eingabe für die Entscheidung zu Frage 2 und für die Folge-ADR; ein Wortlaut für das Lastenheft entsteht hier nicht.
- [x] **Keine Zugangsdaten:** Es wurde kein Zugang verwendet; vor dem ersten Lauf ist `redact` an einer Gegenprobe geprüft (Platzhalter in den Formen, die das Werkzeug ausgeben kann); die Umgebung des Messlaufs trägt keine der Variablen `ANTHROPIC_API_KEY`, `GH_TOKEN`, `GITLAB_TOKEN`, `GIT_TOKEN` und keine Variable mit Ollama-Zugang (geprüft nach **Namen**, nicht nach Werten, im Protokoll belegt); die Suche nach Zugangsdaten (`06-suche.sh`) ist über das Verzeichnis dieser Messung und `tools/messung/ollama-cloud/` leer; über ganz `docs/messungen/` und `tools/messung/` trifft sie ältere Protokolle (öffentliche Hostschlüssel, GPG-Fingerabdrücke, Platzhalter, Testnamen), von Hand als keine Zugangsdaten gelesen (Protokoll, Abschnitt *Zugangsdaten und Rohdaten*).
  **Planänderung gegenüber der Fassung vom Plan-Commit `e60956c3` (Entscheidung: gelten lassen, Verifier im [Verification-Report](../../../reviews/2026-10-24-slice-v1-ollama-cloud-messung-verify.md), Abschnitt *Entscheidung zur Planänderung DoD 3*; vom Planner bei der Closure bestätigt):** Die Fassung vom Plan-Commit lautete „die Suche nach Zugangsdaten über `docs/messungen/` und `tools/messung/` ist leer“. Der Implementer hat sie in diesem Slice umformuliert (Suche leer nur über das Verzeichnis dieser Messung und `tools/messung/ollama-cloud/`; Treffer über den Gesamtbestand von Hand gelesen). Begründung: ältere, abgeschlossene Protokolle tragen öffentliche Schlüsselzeilen, GPG-Kennungen und Platzhalter, sodass die ursprüngliche Zeile über den Bestand nicht erfüllbar ist, ohne fremde Protokolle zu ändern; die Treffer sind von Hand gelesen (Reviewer und Verifier) und keine Zugangsdaten. Der Zweck der Zeile (dieser Slice bringt kein Zugangsdatum ins Repo) ist über den eigenen Bestand getragen. Die Änderung ist eine Abschwächung des Maßstabs, die der Implementer vornahm und kennzeichnete; wirksam ist sie durch die Verifier-Entscheidung und die Bestätigung des Planners.
- [x] `make gates` grün; `make docs-check` 0 Befunde; `make ablage-pruefung` grün (Skripte unter `tools/messung/`, Protokoll unter `docs/messungen/`); `make doc-commits RANGE="origin/main..HEAD"` 0 Befunde.
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
| `tools/messung/ollama-cloud/` | neu | Messskripte (Aufruf `00-alles.sh`; `01` bis `06`, `inner.sh`, `claude-stub.sh`, `Dockerfile.cli`), `lib.sh` mit `redact` und `redact-probe.sh` als Gegenprobe vor dem ersten Lauf, Resolver-Modi `discover` und `restricted`; Muster `tools/messung/ollama-amd64/` (gelesen: `lib.sh`, `05-laufzeit.sh`, `inner-05.sh`) und `tools/messung/forge-cli/` (Aufruf-Form); kein Skript liest oder setzt einen Zugang |
| `docs/messungen/slice-v1-ollama-cloud-messung/` | neu | `README.md` (Messprotokoll: Messstand, C0 bis C7, Liste „nicht gemessen“ mit Messplan) und Rohdaten (`.txt`, `.log`); Umgebung, Datum, Läufe je Aussage; keine Tokenwerte, kein Inhalt von Zugangs- oder Konfigurationsdateien |
| `docs/messungen/README.md` | update | Zeile im Bestand (Verzeichnis und Skripte) |
| `docs/reviews/2026-10-…-slice-v1-ollama-cloud-messung.md`, `…-verify.md` (und, falls nötig, `…-nachlauf.md`) | neu | Review- und Verification-Report der Rollenläufe (DoD; Ablage nach [`MR-009`](../../../../harness/conventions/MR-009-skills-und-review-ablage.md)) |
| `docs/plan/planning/observations/BEO-DOK/…`, `BEO-WZG/…` (`state.md`, `evidence/slice-v1-ollama-cloud-messung.md`) | update / neu | Register-Fortschreibung bei der Closure (§7); kein Zähler wird gesetzt |
| diese Datei, `../welle-ollama-cloud-modelle.md`, `open/README.md`, `in-progress/roadmap.md` | update | diese Datei: `Verantwortlich:`, Gefundener Stand (gemessen), DoD-Häkchen, Plan-Tabelle; Welle-Datei: Zeile des Slice und die Zeilen der Folge-Slices nach dem Ergebnis; `open/README.md`: reiner Verweis-Nachzug von `make slice-mv`; Roadmap: Zeiger auf den Slice während der Arbeit, Ruhe-Marker danach |
| `internal/`, `cmd/`, `spec/`, `docs/user/`, `README.md`, `README.de.md`, `CHANGELOG.md` | **unverändert** | Produktcode, Verträge und Nutzertext folgen in den Folge-Slices der Welle; belegt durch die gelesenen Stellen in §1 (*Gefundener Stand*: Handbuch Zeilen 2023 bis 2029, `devcontainer_egress.go` Zeilen 78 bis 93, `devcontainer_features.go` Zeilen 273, 279 und 535, `devcontainer_sandbox.go` Zeile 441, `lastenheft.md` Zeile 411 und `spezifikation.md` [`SPEC-047`](../../../../spec/spezifikation.md#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen) ohne Aussage zu Cloud-Modellen) |

Ansatz: erst C0 (Unterbefehl und Version), dann das erzeugte Projekt (C7 zuerst, ohne Container), dann Bau und Lauf im Sandbox-Profil (C1 bis C6) im Modus `discover`, danach `restricted` ohne und mit dem Eintrag `ollama.com`. Der Messlauf startet mit bereinigter Umgebung (`env -i`: nur `PATH`, `HOME`, `USER`, `DOCKER_HOST` falls gesetzt, `WORK`, `OUT`, `ROOT`, `LAUF`; die Prüfung nach Namen steht im Protokoll), weil `${localEnv:<NAME>}` in der erzeugten `devcontainer.json` Variablen des Hosts an den Container reicht, sobald die Dev-Container-CLI sie liest. Ein Messcontainer erhält nie ein Host-Verzeichnis mit Konfigurationsdateien eines Agenten.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; kein anderer Slice in `in-progress/` (WIP-Limit 1, gelesen beim Start); Netz und Docker auf dem Messhost verfügbar (Messung auf einem Linux-`amd64`-Host mit lokaler Engine, wie die Zusage-Menge von [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten); ein anderer Host steht im Protokoll und begründet keine Zusage); die Umgebung des Messlaufs trägt keine Zugangsvariable (Prüfung nach Namen). Die offenen Fragen der Welle-Datei (§5) sind dem Projektinhaber vorgelegt; die Messung ohne Konto setzt keine Antwort voraus.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): C0 ergibt, dass `launch` erst mit einer anderen Ollama-Version existiert und die Messung deshalb in zwei Teile zerfällt (Pin-Hebung zuerst, Messung danach); die Teilung entscheidet der Planner.
- `in-progress` → `open` (blockiert): (a) der Messhost erreicht `ollama.com` oder `registry.ollama.ai` nicht (Netz), (b) der Projektinhaber verlangt vor der Messung eine Entscheidung zu Frage 2, die er noch nicht getroffen hat (nicht erwartet: die Messung setzt keine Antwort voraus).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; `make gates` und `make docs-check` sind grün; das Messprotokoll liegt im Repo, trägt je Aussage Umgebung, Datum und Läufe, führt die ungemessenen Punkte mit Messplan und keinen Tokenwert (Suche leer über das Verzeichnis dieser Messung und `tools/messung/ollama-cloud/`; die Treffer über den Gesamtbestand sind von Hand gelesen und keine Zugangsdaten, DoD 3); Review-Report und Verification-Report liegen vor; die Closure-Notiz mit Steering-Loop-Lerneintrag ist geschrieben, und die Frage 2 der Welle-Datei ist dem Projektinhaber mit dem Messplan vorgelegt (die Antwort ist keine Bedingung dieses Slice).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht. Die Ausgänge werden bei der Closure gesetzt; hier steht je Risiko die Vorsorge.

- **Eine Aussage reicht weiter als die Messung** („Cloud-Modelle sind nicht nutzbar“ oder „`ollama.com` genügt“, gemessen ohne Konto). — **Vorsorge:** Aussagen bedingungsgleich zur Messbedingung („ohne Konto“), die Liste „nicht gemessen“ im Protokoll, kein Wortlaut für das Lastenheft ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)). **Ausgang:** bei der Closure.
- **Die Messvariante ist nicht die Auslieferung** (Ollama in anderer Version, Hilfs-Dockerfile statt des erzeugten). — **Vorsorge:** Messung am Ergebnis von `generate devcontainer`; eine Abweichung steht im Protokoll als Variante ([`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](../observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md)). **Ausgang:** bei der Closure.
- **Zugangsdaten gelangen in Rohdaten oder Umgebung des Messlaufs** (Host-Variable `ANTHROPIC_API_KEY`, die die Dev-Container-CLI über `${localEnv:…}` in den Container reicht; Ausgabe eines Werkzeugs mit Anmelde-URL oder Code). — **Vorsorge:** bereinigte Umgebung, Prüfung nach Namen, `redact` mit Gegenprobe, Suche vor dem Commit; kein Inhalt von Konfigurationsdateien im Protokoll; der Git-Verlauf behält, was einmal committet ist ([`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md), 2×). **Ausgang:** bei der Closure.
- **`ollama launch` oder `claude` als Ziel existiert in der gepinnten Version nicht** (C0). — **Vorsorge:** C0 zuerst; Rückführung in §4; Pin-Hebung nur als eigener Slice. **Ausgang:** bei der Closure.
- **Die Messung ruft den Cloud-Dienst von Ollama oder Anthropic auf** (Anfragen ohne Konto, Telemetrie, Kontingent). — **Vorsorge:** Aufrufe bis zum ersten Hindernis, keine Wiederholungsschleifen, Zahl der Läufe im Protokoll; kein Lasttest. **Ausgang:** bei der Closure.
- **Die Messung hängt am Netz und am Cache des Hosts** (Registries, Release-Archive). — **Vorsorge:** kalter Bau, Datum und Stand der Quelle im Protokoll, Aufräumen der Zwischencontainer im Skript ([`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md), [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md)). **Ausgang:** bei der Closure.

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

*(Wird bei der Closure gefüllt: Was hat funktioniert · Was ging anders als geplant · Steering-Loop-Eintrag · Beobachtungs-Register · Folge-Slices · Risiken aus §6 mit Ausgang · Drei Paarungen. Der Slice ist noch nicht gestartet.)*

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind `tools/` (Kürzel WZG, Messskripte unter `tools/messung/`) und `docs/messungen/`, `docs/plan/` (Kürzel DOK); beide stehen als eigene Sub-Areas in der Modus-Deklaration ([`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area)), ausdifferenziert genug. `internal/`, `spec/` und `docs/user/` werden gelesen, nicht geändert (§3).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen am 2026-10-24 (Zähler = Dateien in `evidence/`). Treffer: [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (21), [`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](../observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md) (4) und [`BEO-DOK/vertrag-schweigt-zum-nutzerpfad`](../observations/BEO-DOK/vertrag-schweigt-zum-nutzerpfad/observation.md) (4), alle über der Schwelle, ohne zugewiesenen Ausgang (Entscheidung des Projektinhabers steht aus; der Slice weist nichts zu), Vorsorge in §6 und §2; [`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md) (2, eine Messung mit Zugang wäre der dritte Beleg; dieser Slice nutzt keinen, DoD 3); [`BEO-WZG/host-werkzeug-ausserhalb-docker`](../observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md) (15; Skripte rufen `docker` und `make`, keine Host-Werkzeuge außer `bash` und `docker`); [`BEO-DOK/messprotokoll-ausserhalb-des-repos`](../observations/BEO-DOK/messprotokoll-ausserhalb-des-repos/observation.md) (1; Skripte und Rohdaten von Anfang an im Repo, §3); [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md) (1) und [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md) (2), §6.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

**Alle berührten Sub-Areas GF** (WZG und DOK, Modus-Deklaration in `harness/conventions.md`): kein Modus-Begründungsblock.
