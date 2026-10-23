# Welle welle-forge-cli-gh-glab: Forge-CLIs `gh` und `glab` als Katalog-Features im Devcontainer (Release `v0.15.0`)

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (eine `welle-<Kennung>-results.md` führt u-boot nicht, [`MR-003`](../../../harness/conventions/MR-003-roadmap-wellen.md)). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** Release `v0.15.0`, das nächste Minor nach `v0.14.0` (u-boots Wellen sind die Release-Versionen, [`MR-003`](../../../harness/conventions/MR-003-roadmap-wellen.md)); die Nummer ist im Auftrag an den Planner genannt und vom Projektinhaber am 2026-10-23 bestätigt (§5, Frage 1). Der Release-Cut ist Slice dieser Welle (§4).

**Verantwortlich:** pt9912. **Datum:** 2026-10-23.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Ein Nutzer, dessen Projekt auf GitHub oder GitLab liegt, aktiviert im Devcontainer die Forge-CLI seines Anbieters als Katalog-Feature — `gh` (GitHub) und `glab` (GitLab) — und findet sie nach dem Bau im Container, auch im Sandbox-Profil: mit den Hosts, die das Werkzeug braucht, in der Egress-Liste und ausgewiesen, mit einer Lieferkette, die vor der Zusage geprüft ist ([`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)), und **ohne Zugangsdaten in einer erzeugten Datei oder einem Image**: ein Token kommt nur als Umgebungsvariable des Hosts über `${localEnv:<NAME>}` in den Container ([`LH-FA-DEV-009`](../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)). Vertrag: [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Katalog), [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Hosts) und [`LH-FA-DEV-009`](../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) in der Fassung des Lastenheft-CR dieser Welle. Zugesagt wird, wie bei `ollama` und `claude-code`, nur das Gemessene; welche Quelle, welche Version und welche Hosts, steht erst nach der Messung fest.

Das *Mehr* gegenüber den Slice-DoDs: **beide** Werkzeuge folgen derselben Zusage-Form (Quelle geprüft, Hosts ausgewiesen, Token nur als Umgebungsvariable), Lastenheft, Spezifikation, Handbuch, READMEs, `CHANGELOG.md` und der Katalog im Code sagen über beide dasselbe, der **Nutzerpfad nach der Aktivierung** (Anmeldung im Container, Verhalten ohne gesetzte Variable, Verhalten hinter der Egress-Liste) ist Prüfgegenstand des Validators, und das Release `v0.15.0` ist geschnitten.

Anlass: Beauftragung durch den Projektinhaber am 2026-10-19 (Roadmap, *Nächste Wellen*); Eröffnung am 2026-10-23, nach dem Release-Cut `v0.14.0`.

**Stand des Codes (gelesen am 2026-10-23):** Der Katalog kennt weder `gh` noch `glab` — die Suche nach beiden Namen über `internal/`, `spec/lastenheft.md`, `spec/spezifikation.md`, `docs/user/` und `README.md` findet nichts. Neue Einträge entstehen in `featureCatalogue` (`internal/hexagon/application/devcontainer_features.go`; die Felder `source`, `defaultVersion`, `publisher`, `dockerfileBlock`, `bindsHostSocket`, `relaxations` stehen im Typ `featureCatalogueEntry`), die Hosts je Feature in `egressFeatureHosts` (`internal/hexagon/application/devcontainer_egress.go`; Zeilen für `claude-code` und `ollama` vorhanden). Ob ein Eintrag mit OCI-Quelle genügt oder ein Dockerfile-Block nötig ist, entscheidet die Messung, nicht dieser Plan.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- Der Release `v0.14.0` ist geschnitten (eingetreten: Tag `v0.14.0` auf `68f3c94f`, [`slice-v1-release-cut-v0.14.0`](done/slice-v1-release-cut-v0.14.0.md) liegt in `done/`; die vorige Welle liegt als [`welle-devcontainer-shell-und-zweites-up`](done/welle-devcontainer-shell-und-zweites-up.md) in `done/`).
- Beauftragung durch den Projektinhaber (eingetreten: 2026-10-19, Roadmap *Nächste Wellen*; der Planner-Auftrag vom 2026-10-23 eröffnet die Welle).
- Kein Slice liegt in `in-progress/` (WIP-Limit 1; gelesen am 2026-10-23: `in-progress/` trägt nur `carveouts.md`, `README.md` und `roadmap.md`).
- **Kein** Trigger der Welle: die Messung, die ADR und die Freigabe des Lastenheft-Wortlauts. Sie sind Ergebnis oder Start-Bedingung von Slices (§4), nicht der Welle.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

1. **Alle Slices der Tabelle §4 liegen in `done/`**, der Release-Cut zuletzt. Die ADR aus der Messung steht auf `Accepted` und ist seit ihrem Commit unverändert (`git diff <Commit>..HEAD` über die Datei leer; neben `make doc-immutable`, weil ADRs den Status inline führen, [`BEO-DOK/doc-immutable-greift-nicht-bei-inline-status`](observations/BEO-DOK/doc-immutable-greift-nicht-bei-inline-status/observation.md)).
2. **Jede Zusage hat ihre Messung:** Quelle, Version und Hosts von `gh` und `glab` stehen mit Umgebung und Datum in einem Messprotokoll unter `docs/messungen/`; Lastenheft, Spezifikation, Handbuch und READMEs sagen nicht mehr, als die Messung trägt ([`BEO-DOK/zusage-weiter-als-messung`](observations/BEO-DOK/zusage-weiter-als-messung/observation.md), [`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md)). Jede Bedingung in Spezifikation und Handbuch hat einen Test oder ist gestrichen (Beleg-Tabelle des Spezifikations-Slice).
3. **Keine Zugangsdaten in erzeugten Dateien:** Für beide Features belegt ein Test am erzeugten Ergebnis (Golden Case für `devcontainer.json` und Dockerfile, soweit erzeugt), dass kein Tokenwert und keine Datei mit Zugangsdaten entsteht; die Referenz ist nur `${localEnv:<NAME>}`. Das Messprotokoll enthält keine Tokenwerte ([`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md)); die Suche danach über `docs/messungen/` und `tools/messung/` ist leer.
4. **Nutzerpfad als Prüfgegenstand:** Der Validator fährt am ausgelieferten Binary auf dem zugesagten Host (Linux `amd64`, Docker): (i) Feature aktivieren, `generate devcontainer`, Container bauen — das Werkzeug liegt im Container und meldet seine Version; (ii) im Sandbox-Profil hinter der Egress-Liste: der Aufruf, den die Zusage nennt, gelingt mit den ausgewiesenen Hosts und scheitert mit einer lesbaren Meldung, wenn ein Host fehlt; (iii) Anmeldung mit einem Token aus der Umgebungsvariable des Hosts, und ohne gesetzte Variable: das Verhalten, das Handbuch und Meldung nennen; (iv) `doctor` meldet die Quelle der Zugangsdaten nach [`LH-FA-DEV-009`](../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), soweit der CR sie auf die neuen Features ausdehnt. Eine Abweichung ist ein Befund mit Ausgang (Folge-Slice oder Carveout) ([`BEO-DOK/vertrag-schweigt-zum-nutzerpfad`](observations/BEO-DOK/vertrag-schweigt-zum-nutzerpfad/observation.md)).
5. **Benannte Smokes in der CI:** die neuen Docker-Tests (Bau je Feature mit dem gemessenen Aufruf) sind im Lauf `integration-docker` sichtbar bestanden und nicht übersprungen; `ci` und `integration-docker` sind auf dem Push-Tip des Release-Cuts grün. Ein Host-Lauf ersetzt sie nicht ([`AGENTS.md`](../../../AGENTS.md) §3.1); ein rotes Gate erreicht `done/` nur mit dokumentiertem Carveout.
6. **Gates:** `make gates` grün; `make docs-check` 0 Befunde; `make doc-trace` ohne Waise; `make doc-immutable RANGE="<Commit der ADR>..HEAD"` 0 Befunde.
7. **Validator-Bericht liegt vor**, bevor der Release-Cut beginnt; seine Befunde sind Aufträge an den Cut-Slice.
8. **Release `v0.15.0` ist geschnitten** (Tag gesetzt; Push und Tag nur nach ausdrücklicher Freigabe des Projektinhabers im Gespräch); Kontrolle nach `docs/maintainer/releasing.md` im Cut-Slice.
9. **Closure-Notiz in §7 mit Lerneintrag** (geschärfte Regel · neuer Sensor · benannte Spec-Lücke) und Lese-Schritt des Beobachtungs-Registers: Jeder Eintrag über der Schwelle trägt einen Ausgang oder die benannte Grenze „Entscheidung des Projektinhabers steht aus“; ein Ausgang wird nicht eigenmächtig zugewiesen.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| [`slice-v1-forge-cli-quellen-messung`](done/slice-v1-forge-cli-quellen-messung.md) | Quelle und Prüfung der Lieferkette je Werkzeug, Hosts im Sandbox-Profil, Umgang mit Tokens; Messprotokoll; Folge-ADR; Entwurf des Lastenheft-Wortlauts. Kein Produktcode, kein Lastenheft | [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-009`](../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) |
| [`slice-cr-devcontainer-forge-cli`](in-progress/slice-cr-devcontainer-forge-cli.md) | Lastenheft-CR auf die nächste freie Version: Katalog um `gh` und `glab`, Zusage nur für Gemessenes, Aussage zu Tokens; Wortlaut-Vorschlag mit Freigabe des Projektinhabers | [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-009`](../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) |
| `slice-v1-devcontainer-feature-gh-glab` (Datei entsteht nach CR und ADR) | Code `gh` und `glab` in einem Slice: zwei Katalogeinträge vom Typ Dockerfile-Block, Template-Änderung für beide Profile, Hosts, `remoteEnv` und `doctor`, Golden Cases, Docker-Test, Handbuch, READMEs, `CHANGELOG.md` | [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-009`](../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) |
| `slice-v2-spec-forge-cli` (Datei entsteht nach dem Code) | Spezifikation und Sicht nach dem Code, mit Beleg-Tabelle (Bedingung → Test oder Streichung) | [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) |
| `slice-v1-release-cut-v0.15.0` (Datei erst beim Trigger) | Release-Cut `v0.15.0`: `CHANGELOG.md` zuordnen, Lead mit Plattform und Grenzen, Aufträge aus dem Validator-Bericht; Tag und Push nur nach ausdrücklicher Freigabe | [`ADR-0004`](../adr/0004-ci-system.md), [`ADR-0007`](../adr/0007-distributionswege-ghcr.md) |

Die Kennungen der Slices ab Zeile 3 sind **vorläufig** (Muster der Vorwellen); die Dateien entstehen, wenn ihr Inhalt feststeht, damit sie keine Zustandsmarker mit Ablaufpunkt tragen ([`AGENTS.md`](../../../AGENTS.md) §3.10). Die Folge-ADR gehört zum Messungs-Slice (Muster [`slice-v1-devcontainer-ollama-messung`](done/slice-v1-devcontainer-ollama-messung.md)), sie ist kein eigener Slice. `gh` und `glab` sind ein gemeinsamer Code-Slice (gleiche Quellenform, gleicher Eintragstyp, gleiche Template-Änderung; [`ADR-0035`](../adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) Freigabe-Punkt g, Standard aus der Empfehlung, vom Projektinhaber nicht gefragt).

**Nicht in der Tabelle, aber berührt:** [`slice-v1-katalog-features-quellen-und-freigabe`](open/slice-v1-katalog-features-quellen-und-freigabe.md) (ohne Welle; Abgrenzung in §5).

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- **Blockiert:** keine Welle. `welle-ollama-cloud-modelle` (Roadmap, *Nächste Wellen*) hängt weder an `gh` noch an `glab`; beide Wellen ändern denselben Katalog und dieselbe Egress-Tabelle (`featureCatalogue`, `egressFeatureHosts`), und das WIP-Limit 1 ordnet sie nacheinander. Nichts wird befördert.
- **Wird blockiert von:** keiner offenen Welle. Die Slices ohne Welle in `open/` teilen sich das WIP-Limit mit dieser Welle ([`slice-v1-katalog-features-quellen-und-freigabe`](open/slice-v1-katalog-features-quellen-und-freigabe.md), [`slice-v1-docker-cli-bau-auf-standard-basis-image`](open/slice-v1-docker-cli-bau-auf-standard-basis-image.md), [`slice-v1-shell-erkennung-fehlender-benutzer-eng`](open/slice-v1-shell-erkennung-fehlender-benutzer-eng.md), [`slice-cr-lastenheft-messstaende`](open/slice-cr-lastenheft-messstaende.md), [`slice-v1-devcontainer-start-colima-messung`](open/slice-v1-devcontainer-start-colima-messung.md)).
- **Abgrenzung zu [`slice-v1-katalog-features-quellen-und-freigabe`](open/slice-v1-katalog-features-quellen-und-freigabe.md) (N1 bis N3 von [`ADR-0032`](../adr/0032-katalog-features-socket-und-lockerungen.md)), gelesen am 2026-10-23:** Der Slice bleibt **ohne Welle** und ist **nicht Teil** dieser Welle. N1 (Quellen von `cpp` und `postgres-client`) berührt kein Forge-Werkzeug. N2 (Render-Schlüssel einer externen Quelle, `projectFeatureSourceVersion`, Vorlage `devcontainer.json.tmpl`) und N3 (`generate` prüft die Freigabeliste nicht) betreffen diese Welle **nur dann**, wenn die Messung für `gh` oder `glab` eine Quelle außerhalb des Katalogs ergibt (`source:`-Override oder ein Feature mit eigener Quelle, das der Katalog als externe Quelle führen müsste): Dann schreibt `generate` eine Quelle in einen Schlüssel, die Freigabeprüfung und der Schlüssel sind genau die Lücken von N3 und N2. **Bedingte Abhängigkeit:** Ergibt die Messung eine externe Quelle, liegen N3 (und N2) vor dem Code-Slice des betroffenen Werkzeugs in `done/`, oder der Weg wird ohne externe Quelle gewählt (Katalogeintrag mit OCI-Quelle im Katalog, Dockerfile-Block wie `ollama`); die Entscheidung liegt beim Projektinhaber (Frage 2). Ergibt sie keine externe Quelle, bleibt der Slice unabhängig. Die Messung der beiden Slices ist nicht dieselbe: Der Messungs-Slice dieser Welle misst Werkzeuge, jener die Auflösung eines Schlüssels durch die Dev-Container-CLI; Ergebnisse dürfen sich gegenseitig belegen, aber keiner ersetzt den anderen.

- **Aufträge des Lastenheft-CR an den Code-Slice `slice-v1-devcontainer-feature-gh-glab` (Stand: Lastenheft `0.22.0`):** Die Bedingungen des Lastenheft-Wortlauts brauchen je einen Test im Code-Slice oder werden gestrichen; die Beleg-Tabelle steht in [`slice-cr-devcontainer-forge-cli`](in-progress/slice-cr-devcontainer-forge-cli.md) §7. Ausdrücklich: die Wirkung in beiden Profilen (Standard und Sandbox) braucht einen pinnenden Test (je Profil ein erzeugtes Dockerfile mit dem Block von `gh` und `glab`).

### Fragen an den Projektinhaber (beantwortet am 2026-10-23)

1. **Release-Nummer:** `v0.15.0` (das nächste Minor nach `v0.14.0`). *Antwort:* bestätigt.
2. **Katalog-Slice (N1 bis N3):** abgegrenzt und ohne Welle. *Antwort:* durch die Quellenwahl erledigt, keine eigene Antwort: Der Projektinhaber hat den Dockerfile-Block mit Katalogeintrag ohne `source:` gewählt ([`ADR-0035`](../adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) Punkte a und b), der weder N2 noch N3 braucht; der Weg über `source:` baut in der erzeugten Form nicht (Messprotokoll B1). Die bedingte Abhängigkeit in der Abgrenzung oben tritt nicht ein.
3. **Reihenfolge der Wellen:** `welle-ollama-cloud-modelle` bleibt in *Nächste Wellen*. *Antwort:* bestätigt, diese Welle läuft zuerst (Forge vor Ollama).
4. **Quellenwahl und Wortlaut:** *Antwort:* Quellenwahl und die Punkte a bis d, f und h der ADR entschieden (Release-Archiv als Dockerfile-Block mit SHA-256 für `gh` und `glab`; `GH_TOKEN` und `GITLAB_TOKEN` getrennt; Wirkung in beiden Profilen; Zusage nur für das Gemessene, keine Nachmessung; `glab duo` und `glab orbit` nicht zugesagt). Die Punkte e (Telemetrie nicht beeinflussen, im Handbuch nennen) und g (ein Code-Slice) wurden nicht gefragt und gelten als Standard aus der Empfehlung der ADR. **Der Wortlaut-Entwurf ([`wortlaut-entwurf.md`](../../messungen/slice-v1-forge-cli-quellen-messung/wortlaut-entwurf.md)) ist nicht freigegeben;** der CR-Slice legt ihn dem Projektinhaber vor.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **Alles, was nicht `gh` oder `glab` ist:** weitere Forge-Werkzeuge (etwa `tea`, `bitbucket-cli`), eine Forge-Abstraktion in `u-boot` oder ein Befehl, der Remotes oder Repositories anlegt — die Welle bringt zwei Werkzeuge in den Katalog, keine Funktion von `u-boot`.
- **Anmeldung durch `u-boot`:** `u-boot` führt weder `gh auth login` noch `glab auth login` aus und speichert keinen Token; es reicht Umgebungsvariablen des Hosts durch ([`LH-FA-DEV-009`](../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)). Ein Credential-Helper, der Git-Operationen gegen die Forge authentifiziert, ist Gegenstand der Messung (Frage nach dem Umgang mit Tokens), aber nur als Zusage, wenn gemessen.
- **Selbst gehostete Instanzen** (GitHub Enterprise, GitLab Self-Managed) und Hosts außerhalb der gemessenen: nicht zugesagt; der Nutzer trägt weitere Hosts selbst in die Egress-Liste ein.
- **Plattformen über die Zusage-Menge hinaus** (Linux `amd64`, lokale Docker-Engine, wie [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten); ob der Bau auf `arm64` gemessen wird, entscheidet der Messungs-Slice): nicht zugesagt, soweit nicht gemessen.
- **N1 bis N3 von [`ADR-0032`](../adr/0032-katalog-features-socket-und-lockerungen.md)** ([`slice-v1-katalog-features-quellen-und-freigabe`](open/slice-v1-katalog-features-quellen-und-freigabe.md)) und **N4** ([`slice-v1-docker-cli-bau-auf-standard-basis-image`](open/slice-v1-docker-cli-bau-auf-standard-basis-image.md)): eigene Vorgänge, siehe §5.
- **Änderungen an Gates, Schwellen und Architekturregeln** ([`AGENTS.md`](../../../AGENTS.md) §3.6) und an Accepted-ADRs: Eine Änderung des Katalogwegs grundsätzlich wäre eine Folge-ADR mit `Supersedes`, nicht Teil dieser Welle.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

*(Wird bei der Closure gefüllt: Ergebnis (Zeiger auf den `done/`-Release-Cut-Slice als Detailquelle, [`MR-003`](../../../harness/conventions/MR-003-roadmap-wellen.md)) · Trigger je Punkt · Lerneintrag · Lese-Schritt des Beobachtungs-Registers.)*

Zähler: [`observations/`](observations/README.md) (Sichtung bei der Eröffnung: §8).

## 8. Register-Sichtung bei der Eröffnung

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md` (Eröffnungs-Schritt 2) — gelesen am 2026-10-23, Zähler = Dateien in `evidence/`. Die Sichtung betrifft die Sub-Areas dieser Welle: APP (Katalog, Egress-Tabelle), E2E (Docker-Tests), DOK (Lastenheft, Spezifikation, Handbuch, Messprotokoll), WZG (Messskript).

**Über der Schwelle von drei, ohne zugewiesenen Ausgang (Entscheidung des Projektinhabers steht aus; die Welle weist nichts zu):**

- [`BEO-DOK/zusage-weiter-als-messung`](observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (16): trifft diese Welle unmittelbar, denn sie sagt Werkzeug-Eigenschaften zu, die gemessen werden. Vorsorge im Plan: Closure-Punkt 2, im Messungs-Slice der Messsatz mit Umgebung, Datum und Zahl der Läufe.
- [`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md) (3): die Messung läuft mit der Form, die `u-boot` ausliefert (erzeugtes Dockerfile oder erzeugte `devcontainer.json`), nicht mit einer Hilfsvariante.
- [`BEO-DOK/vertrag-schweigt-zum-nutzerpfad`](observations/BEO-DOK/vertrag-schweigt-zum-nutzerpfad/observation.md) (3): Closure-Punkt 4 nennt den Nutzerpfad als Prüfgegenstand.

**Unter der Schwelle, mit Treffer:** [`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md) (2; Messung einer Anmeldung ist die Gelegenheit zum dritten Beleg: Closure-Punkt 3); [`BEO-DOK/plan-aussage-ueber-code`](observations/BEO-DOK/plan-aussage-ueber-code/observation.md) und [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (beide in [`AGENTS.md`](../../../AGENTS.md) §3.10 verkörpert); [`BEO-DOK/abhilfe-ohne-wirkungsprobe`](observations/BEO-DOK/abhilfe-ohne-wirkungsprobe/observation.md) (1; Meldung bei fehlendem Token nennt eine Abhilfe, die der Implementer ausführt); [`BEO-E2E/docker-test-scheitert-am-externen-netz`](observations/BEO-E2E/docker-test-scheitert-am-externen-netz/observation.md) (1), [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md) (1) und [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md) (2): jeder Bau-Test eines Features aus dem Netz trifft alle drei; [`BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf`](observations/BEO-E2E/kombination-ollama-nested-runtime-ohne-lauf/observation.md): Kombinationen mit `nestedRuntime` sind Prüfgegenstand der Code-Slices, soweit die Messung sie nahelegt.

Die Einträge unter der Schwelle gehören in §8 der jeweiligen Slice-Pläne; erreicht einer mit dieser Welle drei, braucht er einen eigenen Slice.
