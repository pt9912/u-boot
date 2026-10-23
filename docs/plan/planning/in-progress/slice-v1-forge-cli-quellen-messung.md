# Slice v1: Forge-CLIs `gh` und `glab` im Devcontainer — Quelle und Lieferkette prüfen, Hosts und Token-Umgang messen, belegen (Katalog-Schlüssel `gh` und `glab`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-forge-cli-gh-glab`](../welle-forge-cli-gh-glab.md) (Release `v0.15.0`) — erster Slice der Welle (Messung, Quellenwahl, Folge-ADR, Wortlaut-Entwurf; ändert weder Lastenheft noch Produktcode). Ihm folgen der Lastenheft-CR, der Code je Werkzeug, die Spezifikation und der Release-Cut (Welle-Datei §4; Kennungen vorläufig).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Devcontainer-Features; Katalog), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (keine verdeckte Ausführung fremder Skripte; bindet die Wahl der Quelle), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Hosts), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zugangsdaten nur als Umgebungsvariable); [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) (Drittanbieter-Features im Katalog; Muster Messreihe mit der Dev-Container-CLI), [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) (Muster: Quelle wählen und messen, Folge-ADR), [`ADR-0032`](../../adr/0032-katalog-features-socket-und-lockerungen.md) (Metadaten eines Features, Sicherheitseigenschaften).

Anlass: Beauftragung durch den Projektinhaber am 2026-10-19 (Roadmap, *Nächste Wellen*); Welle eröffnet am 2026-10-23.

**Berührte Spec-Stellen:** [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) (Katalog) — der Slice **liest** sie und ändert keine; die Nachführung folgt im Spezifikations-Slice der Welle. Das Lastenheft ändert nur der CR. Der Verweis zeigt **aufwärts**: Lastenheft und Spezifikation nennen diesen Slice nie.

**Verantwortlich:** pt9912 (Rolle Implementer, Start 2026-10-23); Freigaben: Projektinhaber.

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-23.

---

## 1. Ziel und Abgrenzung


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Mit gemessenen Belegen steht fest, **wie** `gh` und `glab` in den Devcontainer kommen — mit welcher Quelle, welcher Prüfung der Lieferkette und welcher Version, welchen Hosts im laufenden Container (Sandbox-Profil, Egress-Liste) und welchem Umgang mit Tokens, ohne dass ein Zugangsdatum in einer erzeugten Datei oder einem Image steht. Ergebnis: ein Messprotokoll, eine Folge-ADR (Quelle und Version gepinnt **nach** der Messung) und der mit Messwerten gefüllte Wortlaut-Entwurf für den Lastenheft-CR. Zugesagt wird später nur, was hier gemessen ist.

### Gefundener Stand (gelesen am 2026-10-23)

- **Katalog.** `featureCatalogue` (`internal/hexagon/application/devcontainer_features.go`) führt weder `gh` noch `glab`; die Suche nach beiden Namen über `internal/`, `spec/lastenheft.md`, `spec/spezifikation.md`, `docs/user/` und `README.md` ist leer. Der Typ `featureCatalogueEntry` trägt `source`, `defaultVersion`, `shortDesc`, `publisher` (Drittanbieter), `dockerfileBlock` (Installation im Sandbox-Dockerfile statt OCI-Feature, nur `ollama`), `bindsHostSocket` und `relaxations`; der Doc-Kommentar von `featureCatalogue` verlangt, bei einem neuen Eintrag die Metadaten (`mounts`, `privileged`, `capAdd`, `securityOpt`, `entrypoint`, `dependsOn`) zu lesen und die Felder zu füllen.
- **Hosts.** `egressFeatureHosts` (`internal/hexagon/application/devcontainer_egress.go`) bildet den Schlüssel eines Features auf seine Laufzeit-Hosts ab; `claude-code` und `ollama` stehen mit gemessenen Hosts darin, der Doc-Kommentar der Funktion unterscheidet Laufzeit- und Bau-Hosts. `gh` und `glab` fehlen.
- **Zugangsdaten.** [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Lastenheft, gelesen): Zugangsdaten im Sandbox-Profil weder im Image, noch im Workspace-Volume, noch in `u-boot.yaml` oder einer anderen erzeugten Datei; Übergabe zur Laufzeit per Umgebungsvariable (`${localEnv:<NAME>}`), alternativ Secret-Mount; `doctor` prüft die Quelle. Die Umsetzung für `claude-code` ist gelesen ([`SPEC-044`](../../../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) *Zugangsdaten des Agenten*; `agentCredentialEnv` in `devcontainer_sandbox.go`, Feld `AgentEnv` in `templates.go`, Schleife `remoteEnv` in `devcontainer.json.tmpl`, `checkSandboxCredentials` in `doctor_sandbox.go`): `remoteEnv` trägt mit Sandbox-Profil `GIT_TOKEN` und die Namen aus `agentCredentialEnv`, jeweils als `${localEnv:<NAME>}`; ohne Sandbox-Profil entsteht kein `remoteEnv`; `doctor` prüft die Referenz. Das ist das Muster des Wortlaut-Entwurfs.
- **Die Werkzeuge selbst: gemessen** (Messstand und Aussagen im [Messprotokoll](../../../messungen/slice-v1-forge-cli-quellen-messung/README.md)); was nicht gemessen ist, steht dort in Abschnitt F.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Lastenheft-CR und der Katalogeintrag im Code** — eigene Slices der Welle (Welle-Datei §4); die Zusage beruht auf der Messung, der CR schreibt sie danach, der Code folgt dem Lastenheft.
- **Ein Produktcode-Eingriff** (Katalog, Egress-Tabelle, Vorlagen) — der Slice liefert Messskript, Protokoll, ADR und Entwurf, wie [`slice-v1-devcontainer-ollama-messung`](../done/slice-v1-devcontainer-ollama-messung.md).
- **N1 bis N3 von [`ADR-0032`](../../adr/0032-katalog-features-socket-und-lockerungen.md)** ([`slice-v1-katalog-features-quellen-und-freigabe`](../open/slice-v1-katalog-features-quellen-und-freigabe.md)) — anderer Gegenstand (Quellen von `cpp` und `postgres-client`, Render-Schlüssel, Freigabeprüfung in `generate`). Berührt die Messung eine externe Quelle, **nennt** das Protokoll diese Abhängigkeit als Ergebnis, der Slice behebt sie nicht (Welle-Datei §5).
- **N4, der Bau von `docker-cli` auf dem Standard-Basis-Image** ([`slice-v1-docker-cli-bau-auf-standard-basis-image`](../open/slice-v1-docker-cli-bau-auf-standard-basis-image.md)) — andere Ursache; zeigt die Messung dort einen gleichartigen Bau-Fehler bei `gh` oder `glab`, ist das ein Befund des Protokolls mit Ausgang, kein Umbau hier.
- **Selbst gehostete Instanzen, weitere Forge-Werkzeuge, Anmeldung durch `u-boot`** — Out-of-Scope der Welle (Welle-Datei §6); das Protokoll misst `github.com` und `gitlab.com` und sagt nichts über andere.

## 2. Definition of Done


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **Messung je Werkzeug** ([`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)): für `gh` und `glab` die in Frage kommenden Quellen (offizielles oder Drittanbieter-Feature, Paketquelle, Release-Archiv, Dockerfile-Block), ihre Lieferkette (wer liefert, was läuft beim Bau, Prüfung durch Signatur, Prüfsumme oder Pin), die Version, die Metadaten des Features (`mounts`, `privileged`, `capAdd`, `securityOpt`), der Bau auf dem Standard-Basis-Image in der Form, die `u-boot` ausliefert, und die Hosts im laufenden Container hinter der Egress-Liste (Laufzeit getrennt von Bau). *Beleg:* Skript unter `tools/messung/<name>/`, Messprotokoll mit Rohdaten unter `docs/messungen/slice-v1-forge-cli-quellen-messung/` ([`harness/conventions.md`](../../../../harness/conventions.md#ablage-von-messprotokollen-und-skripten)), je Aussage Umgebung, Datum und Zahl der Läufe.
- [x] **Umgang mit Tokens gemessen** ([`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)): welche Umgebungsvariable je Werkzeug die Anmeldung trägt, ob ein Aufruf mit gesetzter Variable ohne `auth login` gelingt, ob und wo das Werkzeug dabei eine Datei oder einen Cache schreibt (Pfad, Inhalt ohne Tokenwert), wie es ohne gesetzte Variable endet, und ob Git-Operationen gegen die Forge einen Credential-Helper brauchen. *Beleg:* Messprotokoll; die Rohdaten tragen **keinen Tokenwert** (geschwärzt, vor dem Commit durchsucht; [`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md)); für die Messung dient ein eigens angelegtes, kurzlebiges Token mit minimalen Rechten oder ein Aufruf ohne Token, das Protokoll nennt welches.
- [ ] **Folge-ADR und Wortlaut-Entwurf:** ADR mit `Accepted` (Quelle und Version gepinnt nach der Messung, Verglichene Alternativen, `Schärft:`-Feld nach oben; Muster [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md)), im ADR-Index eingetragen; Wortlaut-Vorschlag für den CR (alt/neu) mit der Liste „nicht gemessen“ je Satz, den der Projektinhaber freigibt (kein Eintrag im Lastenheft in diesem Slice). Die Quellenwahl, die vom Muster abweicht (Drittanbieter, Dockerfile-Block, externe Quelle), entscheidet der Projektinhaber, nicht der Implementer.
- [x] `make gates` grün; `make docs-check` 0 Befunde; `make ablage-pruefung` grün (Skripte unter `tools/messung/`, Protokoll unter `docs/messungen/`); `make doc-commits RANGE="origin/main..HEAD"` 0 Befunde.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: ADR-Index, Welle-Datei und Roadmap (Zeilen der Folge-Slices), `docs/messungen/README.md`, falls dort Verzeichnisse geführt werden; kein öffentlicher Vertrag berührt (Lastenheft, Spezifikation, Handbuch, `CHANGELOG.md` bleiben unverändert).
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
| `tools/messung/forge-cli/` | neu | Messskripte `00-alles.sh` bis `08-aufraeumen.sh`, `lib.sh`, `inner-02*.sh`, `inner-03*.sh`, `inner-07.sh`, `Dockerfile.cli`, vier `block-*.dockerfile`; Muster `tools/messung/ollama-amd64/` (gelesen: `lib.sh`, `02-projekte.sh`, `03-bau.sh`, `05-laufzeit.sh`, `inner-05.sh`, `10-devcontainer-cli.sh`, `11-aufraeumen.sh`, `Dockerfile.cli`) |
| `docs/messungen/slice-v1-forge-cli-quellen-messung/` | neu | `README.md` (Messprotokoll), `wortlaut-entwurf.md` (Wortlaut-Entwurf für den CR) und Rohdaten (`.txt`, `.log`); Umgebung, Datum, Läufe je Aussage; keine Tokenwerte |
| `docs/plan/adr/0035-forge-clis-gh-und-glab-im-devcontainer.md`, `docs/plan/adr/README.md` | neu / update | Folge-ADR nach MADR ([`MR-008`](../../../../harness/conventions/MR-008-adr-form.md)), Status `Proposed`; Index-Zeile (letzte ADR-Datei am 2026-10-23: `0034-…`) |
| `docs/messungen/README.md` | update | Zeile im Bestand (Verzeichnis und Skripte) |
| diese Datei, `../welle-forge-cli-gh-glab.md`, `open/README.md`, `in-progress/roadmap.md` | update | diese Datei: `Verantwortlich:`, Gefundener Stand (gelesen/gemessen), DoD-Häkchen, Plan-Tabelle; Welle-Datei: Absatz nach §4-Tabelle und Frage 4 in §5; `open/README.md`: reiner Verweis-Nachzug von `make slice-mv`; Roadmap: Zeiger unter *Offene Wellen* statt „Nichts in Arbeit.“ |
| `internal/`, `cmd/`, `spec/`, `docs/user/`, `README.md`, `README.de.md`, `CHANGELOG.md` | **unverändert** | Produktcode, Verträge und Nutzertext folgen in den Folge-Slices der Welle; belegt durch die Lesung der Katalog-, Egress- und Vorlagen-Stellen (§1 Gefundener Stand; zusätzlich gelesen am 2026-10-23: `Dockerfile.tmpl` (`USER root` unter `or .UID .Sandbox`, Ollama-Block darin, Credential-Helper), `devcontainer.json.tmpl` (`remoteEnv` nur im Sandbox-Zweig), `devcontainerTemplateData` in `devcontainer_sandbox.go` (Rückkehr vor den Schaltern, wenn kein Sandbox-Profil), `egressBaseHosts` in `devcontainer_egress.go`) |

Ansatz: erst `gh`, dann `glab` je nach demselben Raster (Quelle → Lieferkette → Bau in der ausgelieferten Form → Hosts → Token); ein Skript, ein Protokoll, Rohdaten je Lauf. Die Reihenfolge der Folge-Slices richtet sich nach dem Ergebnis.

## 4. Trigger


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; kein anderer Slice in `in-progress/` (WIP-Limit 1, gelesen beim Start); Netz und Docker auf dem Messhost verfügbar (Messung auf einem Linux-`amd64`-Host mit lokaler Engine, wie die Zusage-Menge von [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten); ein anderer Host steht im Protokoll und begründet keine Zusage). Die offenen Fragen der Welle-Datei (§5) sind dem Projektinhaber vorgelegt; die Messung selbst setzt keine Antwort voraus.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Messung zeigt, dass jedes Werkzeug einen eigenen Weg braucht und die Teilung in zwei Messungs-Slices sinnvoller ist; die Teilung entscheidet der Planner.
- `in-progress` → `open` (blockiert): (a) eine Quelle verlangt eine externe Quelle, deren Auflösung und Freigabeprüfung ([`slice-v1-katalog-features-quellen-und-freigabe`](../open/slice-v1-katalog-features-quellen-und-freigabe.md), N2 und N3) offen sind, und der Projektinhaber entscheidet über die Reihenfolge nicht; (b) der Projektinhaber entscheidet eine Quellenwahl nicht, die vom Muster abweicht.

## 5. Closure-Trigger


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; `make gates` und `make docs-check` sind grün; die Folge-ADR steht auf `Accepted` und im Index; das Messprotokoll liegt im Repo, trägt je Aussage Umgebung, Datum und Läufe und keinen Tokenwert (Suche über `docs/messungen/` und `tools/messung/` leer); der Wortlaut-Vorschlag liegt vor; Review-Report und Verification-Report liegen vor; die Closure-Notiz mit Steering-Loop-Lerneintrag ist geschrieben.

## 6. Risiken und offene Punkte


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

Die Ausgänge stehen bei der Closure; bis dahin nennt jede Zeile die Vorsorge.

1. **Eine Zusage reicht weiter als die Messung** (eine Quelle „funktioniert“, gemessen mit einer Version, einem Host, einem Lauf). — **Ausgang:** offen bis zur Closure. Vorsorge: Messsatz mit Umgebung, Datum und Zahl der Läufe; der Wortlaut-Entwurf führt die Liste „nicht gemessen“ je Satz ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md), 16×, offen).
2. **Die Messvariante ist nicht die Auslieferung** (ein Werkzeug wird in einem Hilfs-Dockerfile gemessen, `u-boot` erzeugt eine andere Form). — **Ausgang:** offen bis zur Closure. Vorsorge: der Bau läuft mit dem Ergebnis von `generate devcontainer` (Sandbox-Profil und ohne), nicht mit einer Hand-Variante; was anders läuft, steht im Protokoll ([`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](../observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md)).
3. **Zugangsdaten in den Rohdaten** (Token in Umgebung, Konfigurationsdatei oder Ausgabe des Werkzeugs). — **Ausgang:** offen bis zur Closure. Vorsorge: kurzlebiges Token mit minimalen Rechten, Schwärzen im Skript, Durchsuchen von Rohdaten und Verlauf vor dem Commit; ein Fund wäre der dritte Beleg für [`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md) (2×).
4. **Ein Werkzeug braucht eine externe Quelle oder einen Drittanbieter**, deren Auflösung und Freigabeprüfung im Katalog noch Lücken haben (N2, N3 des Katalog-Slice) oder deren Lieferkette der Projektinhaber nicht trägt. — **Ausgang:** offen bis zur Closure. Vorsorge: das Protokoll nennt die Abhängigkeit als Ergebnis; §4 (Rückführung nach `open`) und Frage 2 der Welle-Datei.
5. **Die Messung braucht das Netz und hängt am Cache des Hosts** (Registries, Paketquellen, die Forge). — **Ausgang:** offen bis zur Closure. Vorsorge: Läufe ohne lokalen Cache, Datum und Stand der Quelle im Protokoll, Aufräumen der Zwischencontainer im Skript ([`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md), [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md)).

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

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `tools/` (WZG, Messskript) sowie `docs/plan/`, `docs/messungen/` und `harness/` (DOK: ADR, Protokoll, Planung). Beide sind in der Modus-Deklaration von [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) als Greenfield auditiert; eine zu grobe Sub-Area ist nicht dabei. Code-Sub-Areas (APP, E2E, CLI) berührt der Slice nicht.

**Vorgelagert — offene Beobachtungen sichten** (Register gelesen am 2026-10-23; Zähler = Dateien in `evidence/`): [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (16, offen über der Schwelle) und [`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](../observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md) (3, offen über der Schwelle): Risiken 1 und 2; [`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md) (2): Risiko 3; [`BEO-DOK/messprotokoll-ausserhalb-des-repos`](../observations/BEO-DOK/messprotokoll-ausserhalb-des-repos/observation.md) (1): Skripte und Rohdaten liegen von Anfang an im Repo (DoD); [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (in [`AGENTS.md`](../../../../AGENTS.md) §3.10 verkörpert): Prüfpflicht für jede Aussage über Code; [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md) (1) und [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md) (2): Risiko 5. Evidenz-/Diskrepanz-Risiko insgesamt: hoch (über die Werkzeuge ist nichts gelesen oder gemessen, §1).

**Modus-Begründungsblock:** alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
