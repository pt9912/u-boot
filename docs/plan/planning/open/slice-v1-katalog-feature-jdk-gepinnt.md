# Slice V1: Katalog-Feature für ein gepinntes JDK — Dockerfile-Block mit Release-Archiv und SHA-256, und die Korrektur des Handbuchbeispiels `java.version 21`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** Vorschlag Release `v0.18.0` (Roadmap *Nächste Wellen*, Feature-Katalog); bis zur Entscheidung des Projektinhabers ist der Slice ohne Welle geführt und in der Roadmap nicht eingetragen. Anlass: Wunsch des Projektinhabers vom 2026-10-31, Kotlin/Gradle-Projekte im Devcontainer mit gepinnten Werkzeugen zu bauen.

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Katalog der Features; ob der neue Schlüssel einen Lastenheft-CR braucht, entscheidet der Projektinhaber, siehe §4), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (Ausweisung dessen, was beim Bau läuft); [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) (Katalog und Drittanbieter), [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) (Muster Release-Archiv mit Prüfsumme im Dockerfile-Block).

**Berührte Spec-Stellen:** [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) (Katalogtabelle), voraussichtlich [`SPEC-050`](../../../../spec/spezifikation.md#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen) (Katalog im Sandbox-Profil); die Stellen sind nicht im Einzelnen gelesen, die Suche steht beim Start. Der Verweis zeigt **aufwärts**: Lastenheft und Spezifikation nennen diesen Slice nie.

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

**Ziel:** Wer einen Devcontainer für ein JVM-Projekt erzeugt, bekommt ein JDK, dessen Version und Release-Archiv `u-boot` pinnt und beim Bau gegen eine SHA-256-Prüfsumme prüft — und das Handbuch nennt für die Java-Version nur eine Angabe, die es gibt.

**Befund** (gelesen am 2026-10-31; was nicht gelesen oder gemessen ist, steht je Punkt; [`AGENTS.md`](../../../../AGENTS.md) §3.10):

- **Muster der Dockerfile-Blöcke.** `featureCatalogueEntry.dockerfileBlock` (`internal/hexagon/application/devcontainer_features.go`, Zeilen 165 bis 171) kennzeichnet einen Eintrag ohne OCI-Quelle, dessen Installation das Dockerfile trägt; „the pin stays in the template“. Es gibt drei Einträge: `ollama`, `gh`, `glab` (Zeilen 305 bis 316). Version und Prüfsummen stehen als `ARG` in `templates/devcontainer/Dockerfile.tmpl` (`GH_VERSION` Zeile 104, `GLAB_VERSION` Zeile 125, `OLLAMA_VERSION` Zeile 77), je mit `SHA256` für `amd64` und `arm64` und `sha256sum -c` vor dem Entpacken.
- **Der Nutzer wählt die Version heute nicht.** `validateDevcontainerFeatures` weist `version:` ohne `source:` bei einem Dockerfile-Block ab (Doc-Kommentar Zeilen 685 bis 688, Prüfung über `featureVersionNeedsSource`; Exit-Code `10`). Eine nutzergewählte JDK-Version ist damit kein Fall des bestehenden Musters, sondern eine neue Zusage (§4, Weg B).
- **Der Eintrag `java` ist ein OCI-Feature.** `featureCatalogue` führt `java` → `ghcr.io/devcontainers/features/java` mit der Standardversion `1` (Zeilen 252 bis 256); der Katalog nennt es „Java + SDKMAN“. **Nicht gelesen in dieser Rolle:** das Feature selbst, welche JDK-Version es installiert und ob diese gepinnt ist.
- **Das Handbuchbeispiel.** `docs/user/benutzerhandbuch.md` zeigt `u-boot config set devcontainer.features.java.version 21` und sagt, der Eintrag heiße dann `…/java:21` (Zeilen 1003 bis 1004 und 1029; wieder bei 1218 bis 1219). `version` ist der OCI-Tag des Features. Eine Messung in der Sitzung vom 2026-10-31 ergab, ein Tag `21` gebe es nicht; **in dieser Rolle weder gelesen noch wiederholt** und kein Messprotokoll im Repo gefunden. Die Messung steht am Anfang des Slice (§2).
- **Der Katalog steht im Lastenheft.** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) führt die Features als Liste (Zeilen 400 bis 408); `ollama` kam mit Lastenheft 0.15.0, `gh` und `glab` mit 0.22.0 (Historie, Zeilen 1596 und 1603). Die Liste ist eine Beispielliste, keine abschließende: `python`, `homebrew` und `sdkman` kamen ohne Lastenheft-Änderung ([`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md) §1, dort gelesen), `ollama`, `gh` und `glab` mit einer Zusage samt Prüfsumme und Messbedingung per Lastenheft-CR. Ob ein gepinntes JDK mit Prüfsumme in die zweite Klasse fällt, ist eine Entscheidung des Projektinhabers (§4); **nicht gelesen** ist, ob `gh` und `glab` den CR wegen der Prüfsumme oder wegen der Forge-Anmeldung brauchten.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Gradle als Katalogeintrag** — ein eigener Vorgang: andere Quelle, anderer Pin (Distribution mit Prüfsumme) und die Frage, ob ein Projekt mit Wrapper ihn braucht; der Projektinhaber hat ihn nicht gefordert.
- **Kotlin-Compiler, detekt, ktlint** — kommen über die Plugins des Projekts beim Gradle-Lauf; `u-boot` pinnt sie nicht und kann es nicht (Projektsache, Versionskatalog und Sperrdatei des Projekts).
- **Der lokale Modus des Makefiles (`BUILD_MODE`) im Kotlin-Projekt** — Sache des jeweiligen Projekts; `u-boot` ändert fremde Repositories nicht und schreibt kein Makefile ([`ADR-0037`](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md) gilt für `u-boot` selbst).
- **Die Anhebung der JDK-Version eines Projekts** — Projektänderung (Toolchain-Angabe, Basis-Image, Gradle-Version), kein Gegenstand von `u-boot`.
- **Der Lastenheft-CR selbst** — ein eigener CR-Slice nach dem Muster von [`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md), noch nicht angelegt; falls der Projektinhaber einen CR verlangt, liegt er vor diesem Slice (§4).

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Katalogeintrag JDK** (nach dem Weg, den der Projektinhaber in §4 wählt): Dockerfile-Block in `Dockerfile.tmpl` mit Version, SHA-256 für `amd64` und `arm64` aus der Quelle des Herstellers und `sha256sum -c` vor dem Entpacken, Katalogeintrag, Golden Case, ein Docker-Test, in dem der Bau läuft und `java -version` die gepinnte Version meldet. *Beleg:* Tabellentest des Katalogs, Golden Case, Docker-Test im CI-Lauf `integration-docker`; die Prüfsummen sind am Tag der Messung gegen die Herstellerseite gelesen, nicht aus dem Gedächtnis ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)).
- [ ] **Handbuchbeispiel `java.version`**: erst die Messung, welche Tags das `java`-Feature hat (Datum und Dev-Container-CLI-Version im Protokoll unter `docs/messungen/<kennung>/`, [`harness/conventions.md`](../../../../harness/conventions.md#ablage-von-messprotokollen-und-skripten)), dann nennt das Handbuch nur eine Angabe, die es gibt; das Beispiel mit `21` entfällt oder wird berichtigt (Stellen in §3).
- [ ] `make gates` grün; `make docs-check` 0 Befunde; `make test-docker` (der Docker-Pfad ist berührt; Nachweis ist der CI-Lauf).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: Spezifikation nach dem Code mit Beleg-Tabelle (Bedingung → Test oder Streichung, [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)), Handbuch, READMEs und `CHANGELOG.md` ohne Kennungen im Nutzertext ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6), `harness/replay.md`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt (gelesen am 2026-10-31, wo eine Stelle steht). Eine Zeile „prüfen“ nennt eine Datei, die der Planner nicht gelesen hat. Sie wird nach der letzten Fix-Runde gegen `git diff --stat` abgeglichen (`AGENTS.md` §3.10); die Zeilennummern verschieben sich, der Implementer liest sie am Ist-Stand nach.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/application/devcontainer_features.go` | update | Katalogeintrag und Konstante des neuen Schlüssels neben `featureOllama`, `featureGh`, `featureGlab` (Zeilen 196 bis 199), Eintrag mit `dockerfileBlock: true`; `dockerfileBlockFeatureActive` (Zeile 549) und die Warnungen (ab Zeile 616) lesen den Schlüssel; bei Weg B zusätzlich die Prüfung `featureVersionNeedsSource` (Zeile 687) |
| `internal/hexagon/application/templates/devcontainer/Dockerfile.tmpl` | update | Block `{{- if .Jdk}}` nach dem Muster von `.Gh` (ab Zeile 96): `ARG` für Version und beide Prüfsummen, `sha256sum -c`, Entpacken nach `/usr/local`, `JAVA_HOME` und `PATH` |
| `internal/hexagon/application/` (Renderer der Vorlagen-Daten), `config.go` | update (prüfen) | wo `.Gh` und `.Glab` an die Vorlage gereicht werden und wo `config.go` (Zeile 554) den Schlüssel kennt, ist nicht gelesen |
| `internal/hexagon/application/export_test.go`, `devcontainer_features_test.go`, `generate_features_test.go` | update / neu | Tabellentest des Katalogs, Golden Case je Profil; die Namen wählt der Implementer |
| `internal/e2e/devcontainer_start_docker_test.go` oder ein neuer Docker-Test | update / neu | Bau mit dem Block, `java -version` im Container; Aufräumen wie die Nachbartests |
| `tools/messung/<name>/`, `docs/messungen/<kennung>/` | neu | Messung der Tags des `java`-Features und der Herstellerprüfsummen (Datum, CLI-Version, Zahl der Läufe) |
| `spec/lastenheft.md` | update (eigener CR-Slice, vor diesem) | Liste in [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Zeilen 400 bis 408) und Historie |
| `spec/spezifikation.md` | update | [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) (Katalogtabelle, Zeile 1453 nennt die Schlüssel), [`SPEC-050`](../../../../spec/spezifikation.md#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen) (prüfen); nach dem Code, mit Beleg-Tabelle |
| `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md`, `harness/replay.md` | update (prüfen) | Katalogtabelle und Abschnitt zu Java (Handbuch Zeilen 1003 bis 1004, 1029, 1050, 1218 bis 1219 gelesen; die Tabellenzeile des neuen Schlüssels und die Namenslisten der READMEs sind nicht gelesen); ohne Kennungen im Nutzertext |
| `spec/architecture.md`, `docs/plan/adr/`, `Makefile`, `.golangci.yml`, `.d-check.yml` | **unverändert** | Accepted-ADRs sind nicht änderbar; ein anderer Weg als das Muster von [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) braucht eine Folge-ADR; Gates und Schwellen bleiben ([`AGENTS.md`](../../../../AGENTS.md) §3.6) |
| diese Datei, `docs/plan/planning/open/README.md` | update | Zeile dieses Slice; Moves per `make slice-mv` |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Entscheidungen des Projektinhabers vor dem Start** (der Implementer wählt nicht):

1. **Weg.** **A:** ein Schlüssel mit einer im Template fest gepinnten Version (wie `gh`); die Hebung ist eine Änderung des Templates, der Nutzer wählt keine Version, ein Projekt mit anderer JDK-Version trägt das JDK weiter von Hand in den handgeschriebenen Teil des Dockerfiles ein. **B:** die Version ist wählbar aus einer Tabelle gepinnter Versionen im Template (je Version beide Prüfsummen); das ändert die Regel, dass `version:` bei einem Dockerfile-Block abgewiesen wird, und ist eine zusätzliche Zusage im Lastenheft. **C:** kein neuer Schlüssel; das Handbuch zeigt, wie das JDK im handgeschriebenen Teil gepinnt wird (nur die Korrektur des Beispiels `java`).
2. **Hersteller und Archivform** des JDK (die Wahl bestimmt Quelle, Prüfsummen und Lizenzaussage).
3. **Version** bei Weg A (oder die Menge bei Weg B), abhängig von dem, was die Projekte brauchen; das JDK der Basis-Images und `jvmToolchain` eines Projekts müssen zusammenpassen.

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; der Projektinhaber hat Weg, Hersteller und Version entschieden und den Slice einem Release zugeordnet; verlangt er einen Lastenheft-CR (Weg B ja, weil die Zusage „wählbare Version“ neu ist; Weg A nach seiner Entscheidung), liegt der CR-Slice in `done/` (beobachtbar: die Datei liegt dort) und das Lastenheft nennt den Schlüssel; kein anderer Slice in `in-progress/` (WIP-Limit 1, gelesen beim Start). Bei Weg C entfällt der CR.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Weg B zeigt, dass die Versionstabelle mehr als ein JDK-Format oder mehr als einen Hersteller braucht; die Teilung entscheidet der Projektinhaber.
- `in-progress` → `open` (blockiert): (a) die Herstellerseite liefert für eine der beiden Architekturen kein Archiv mit veröffentlichter Prüfsumme; (b) der Projektinhaber entscheidet den Weg nicht.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; `make gates` und `make docs-check` sind grün; der Docker-Test läuft im Lauf `integration-docker` des CI sichtbar bestanden (oder die Closure benennt den Nachweis, der an dessen Stelle trat, mit Grund); jede Bedingung der Beleg-Tabelle ist durch einen Test belegt oder gestrichen; Review-Report, ein Nachlauf-Review über jeden Diff von Vertragstext und der Verification-Report liegen vor; die Closure-Notiz mit Steering-Loop-Lerneintrag ist geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

Die Ausgänge stehen bei der Closure; bis dahin nennt jede Zeile die Vorsorge.

1. **Prüfsummen aus dem Gedächtnis oder von einer Spiegelseite.** — **Ausgang:** offen bis zur Closure. Vorsorge: Prüfsummen werden am Tag der Messung von der Herstellerseite gelesen und im Protokoll datiert; ein Reviewer vergleicht sie im frischen Kontext erneut.
2. **Das JDK der Basis-Images und das gepinnte JDK weichen ab** (das `java`-Feature oder ein Basis-Image bringt ein eigenes JDK, `PATH`-Reihenfolge entscheidet). — **Ausgang:** offen bis zur Closure. Vorsorge: der Docker-Test prüft `java -version` und `JAVA_HOME` im gebauten Container; die Wechselwirkung mit dem Eintrag `java` ist nicht gelesen und wird beim Start gemessen.
3. **Eine Zusage reicht weiter als die Messung** (ein JDK-Archiv gemessen auf einer Architektur oder einem Host). — **Ausgang:** offen bis zur Closure. Vorsorge: Messsatz mit Umgebung, Datum und Zahl der Läufe; keine Aussage über nicht gemessene Hosts ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)).
4. **Eine Bedingung in Spezifikation oder Handbuch hat keinen Test** (zum Beispiel „nur im Sandbox-Profil“ oder „in beiden Profilen“). — **Ausgang:** offen bis zur Closure. Vorsorge: Beleg-Tabelle nach dem Code, jeder Teilsatz mit eigenem Pin ([`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)).
5. **Der Docker-Test braucht Netz** (Download des Archivs) und hängt am Cache des Hosts. — **Ausgang:** offen bis zur Closure. Vorsorge: Aufräumen wie die Nachbartests, Nachweis ist der CI-Lauf ([`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md), [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md)).

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

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `internal/hexagon/application` (APP), `internal/e2e` (E2E), `tools/` (WZG, Messskript) sowie `spec/`, `harness/`, `docs/plan/` und `docs/messungen/` (DOK); `docs/user/`, die READMEs und `CHANGELOG.md` zählen wie in den übrigen Slices zu DOK. Alle sind in der Modus-Deklaration von [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) als Greenfield auditiert; eine zu grobe Sub-Area ist nicht dabei.

**Vorgelagert — offene Beobachtungen sichten** (Register am 2026-10-31 nicht neu durchgegangen; nur die in diesem Plan zitierten Einträge sind gelesen, ihre Zähler-Stände nicht neu erhoben — beim Start durchgehen): [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (Risiko 3), [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (Risiko 4), [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md) und [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md) (Risiko 5). Evidenz-/Diskrepanz-Risiko insgesamt: mittel (die Tags des `java`-Features und die Wechselwirkung mit dem JDK der Basis-Images sind nicht gemessen).

**Modus-Begründungsblock:** alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
