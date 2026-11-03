# Slice V1: Katalog-Feature für ein gepinntes JDK — der Eintragstyp „gepinnte Werkzeugkette" ([`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md)), Dockerfile-Block mit Release-Archiv und SHA-256, und die Korrektur des Handbuchbeispiels `java.version 21`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** Vorschlag Release `v0.20.0` (Roadmap *Nächste Wellen*, Feature-Katalog); bis zur Entscheidung des Projektinhabers ist der Slice ohne Welle geführt und in der Roadmap nicht eingetragen (die Welle-Zeile der Roadmap trägt die Nachfolge-Slices der gepinnten Werkzeugkette, nicht diesen Slice). Anlass: Wunsch des Projektinhabers vom 2026-10-31, Kotlin/Gradle-Projekte im Devcontainer mit gepinnten Werkzeugen zu bauen. Der Eintragstyp „gepinnte Werkzeugkette" steht seit [`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) (`Accepted`, 2026-11-02) und Lastenheft 0.30.0 als Katalog-Typ fest; dieser Slice erbt ihn.

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Katalog der Features; der Eintragstyp samt Hybrid-Versionierung, S1 bis S11, steht dort seit 0.30.0 — ob der JDK-Schlüssel als neuer Eintrag einen Lastenheft-CR braucht, entscheidet der Projektinhaber, siehe §4), [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) (Pin-Politik; trägt die zwei Go-Schlüssel, [`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) Entscheidung 4), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (Ausweisung dessen, was beim Bau läuft); [`ADR-0020`](../../adr/0020-drittanbieter-features-im-katalog.md) (Katalog und Drittanbieter; die Namensraum-Regel gilt laut [`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) §Verhältnis nur für Einträge mit OCI-Quelle), [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) (Muster Release-Archiv mit Prüfsumme im Dockerfile-Block), [`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) (der Typ).

**Berührte Spec-Stellen** (gelesen am 2026-11-02): [`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) (der Eintragstyp, gelesen), [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) (Katalog und Einträge, Abschnitt gelesen), [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) (Hosts-Zeilen (das Ollama-Konto gilt für die Externen; die dritte Zeile kommt mit [`SPEC-054`](../../../../spec/spezifikation.md)). Die [`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)-Zeilen der drei Namen sind im CR geliefert; eine eigene Technik-Stelle neben der [`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) ist nicht geplant (Freigabe-Punkt F6 entschieden). [`SPEC-050`](../../../../spec/spezifikation.md#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen) (Katalog im Sandbox-Profil; die Stellen sind nicht im Einzelnen gelesen, die Suche steht beim Start). Der Verweis zeigt **aufwärts**: Lastenheft und Spezifikation nennen diesen Slice nie.

**Verantwortlich:** pt9912 (Implementer; Umsetzung des freigegebenen CR-Wortlauts 0.31.0, Start beim freien WIP-Limit)

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-11-02 (Erstfassung 2026-10-31, vor [`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md)).

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Wer einen Devcontainer für ein JVM-Projekt erzeugt, bekommt ein JDK als Katalog-Eintragstyp „gepinnte Werkzeugkette" ([`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)): dessen Version und Release-Archiv pinnt `u-boot`, und der Bau prüft die SHA-256-Prüfsumme der offiziellen Prüfsummen-Quelle vor dem Entpacken; der Default-Pin steht im Template und ist per `version:` übersteuerbar, die Übersteuerung verifiziert gegen die offizielle Prüfsummen-Quelle — und das Handbuch nennt für die Java-Version nur eine Angabe, die es gibt.

**Befund** (gelesen am 2026-11-02, nach der Umsetzung [`slice-v1-katalog-gepinnte-toolchain`](../done/slice-v1-katalog-gepinnte-toolchain.md); was nicht gelesen oder gemessen ist, steht je Punkt; [`AGENTS.md`](../../../../AGENTS.md) §3.10):

- **Der Typ steht im Vertrag seit 0.30.0.** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) trägt die Sätze S1 bis S11 (Zeile 424): Dockerfile-Block mit Release-Archiv und SHA-256 je Architektur, Default-Pin als Build-Argument (nur mit einem u-boot-Release geändert), `version:`-Übersteuerung nur mit offizieller Prüfsummen-Quelle und Lese-Prüfung bei der Erzeugung, Abweisungen je Exit-Code `10`, Ausweisung beim Erzeugen. [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) (Zeile 522) trägt die Pin-Politik der vier Stellen — an die zwei Go-Schlüssel gebunden, gepinnte Werte Go `1.27.3`, `golangci-lint` `v2.14.0`; die Historie-Zeile 0.30.0 ([`spec/lastenheft.md`](../../../../spec/lastenheft.md), Zeile 1627) trägt sie als Chronik. Die Umsetzung liegt in [`slice-v1-katalog-gepinnte-toolchain`](../done/slice-v1-katalog-gepinnte-toolchain.md), der CR-Slice war [`slice-cr-devcontainer-gepinnte-werkzeugkette`](../done/slice-cr-devcontainer-gepinnte-werkzeugkette.md).
- **Die Typ-Maschine im Code ist generalisierbar, ihre Schalter sind geschlossen.** Der Katalog-Eintrag trägt `dockerfileBlock: true` plus die Angabe `checksumSource` (`internal/hexagon/application/devcontainer_features.go`Zeilen 173 bis 182; Muster-Einträge `go-toolchain`/`golangci-lint` Zeilen 338 bis 347, Konstanten Zeilen 206 bis 221). Die Regel-Brücke „ein Eintrag ohne OCI-Quelle nimmt kein `version:` ohne `source:`, ein Eintrag mit Prüfsummen-Quelle schon" ist allgemein und für synthetische Einträge testbar (`featureVersionNeedsSourceForCat`, Zeilen 844 bis 852). **Aber** die Typ-Grenze `toolchainChecksumSourceSupported` schaltet auf eine geschlossene Liste der zwei Konstanten (`internal/hexagon/application/devcontainer_toolchain.go` Zeilen 159 bis 166), ebenso der Version-Normalisierer (Zeilen 237 bis 255) und der Quellen-Leser `readToolchainChecksums` (Zeilen 261 bis 269) — ein dritter Schlüssel erweitert jede dieser Schaltstellen um seine Quelle und ihre URL-Formen. Die Render-Zeile der Blöcke (`resolveToolchain`, Abweiser und Pin-Slots Zeilen 194 bis 235; Template-Daten in `generate.go` Zeilen 848 bis 852) liest den Zustand je Schlüssel.
- **Der Schlüssel `java` ist belegt.** `featureCatalogue` führt `java` → `ghcr.io/devcontainers/features/java` mit Standardversion `1` und „Java + SDKMAN“ (Zeilen 273 bis 277); `sdkman` trägt dieselbe Quelle und teilt mit ihm den Render-Schlüssel, bei abweichender Version bleiben beide (Kommentar `dedupeFeatureData`, Zeilen 457 bis 460). Die Hosts-Zeile `feature:java` der Default-Allowlist trägt `repo.maven.apache.org` und `repo1.maven.org` ([`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion); Handbuch-Zeile 1682). **Nicht gelesen in dieser Rolle:** das Feature selbst, welche JDK-Version es installiert und ob diese gepinnte ist.
- **Das Handbuchbeispiel steht noch ungeprüft da.** `docs/user/benutzerhandbuch.md` zeigt `u-boot config set devcontainer.features.java.version 21` und sagt, der Eintrag heiße dann `"ghcr.io/devcontainers/features/java:21": {}` (Zeilen 1087 bis 1088 und 1113; wieder bei 1304 bis 1305). Eine Messung in der Sitzung vom 2026-10-31 ergab, ein Tag `21` gebe es nicht; **in dieser Rolle weder gelesen noch wiederholt** und kein Messprotokoll im Repo (`docs/messungen/` trägt keines zu `java`). Die Messung steht am Anfang des Slice (§2); die Stellen sind jetzt im Befund gelesen, nicht an den alten Zahlen vom 2026-10-31 geführt.
- **Die Weg-Frage ist durch [`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) beantwortet.** Die Drei-Weg-Frage dieses Plans (A fest gepinnt, B Tabelle, C Handbuch) hat die Entscheidung von 2026-11-02 gelöst — Hybrid: Default-Pin im Template, übersteuerbar per `version:`, Verifikation gegen die offizielle Prüfsummen-Quelle; A und B sind als Alleinlösungen verworfen (Abwägung), D „nichts tun" ebenfalls — und der CR für den Typ (wählbare Version mit Quelle, Abweisung ohne sie, Exit-Code `10`) ist mit dem Lastenheft 0.30.0 geschrieben. Der JDK-Schlüssel stellt die Weg-Frage nicht mehr; was noch zu entscheiden ist, steht in §4.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Gradle als Katalogeintrag** — ein eigener Vorgang: andere Quelle, anderer Pin (Distribution mit Prüfsumme) und die Frage, ob ein Projekt mit Wrapper ihn braucht; der Projektinhaber hat ihn nicht gefordert.
- **Kotlin-Compiler, detekt, ktlint** — kommen über die Plugins des Projekts beim Gradle-Lauf; `u-boot` pinnt sie nicht und kann es nicht (Projektsache, Versionskatalog und Sperrdatei des Projekts).
- **Der lokale Modus des Makefiles (`BUILD_MODE`) im Kotlin-Projekt** — Sache des jeweiligen Projekts; `u-boot` ändert fremde Repositories nicht und schreibt kein Makefile ([`ADR-0037`](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md) gilt für `u-boot` selbst).
- **Die Anhebung der JDK-Version eines Projekts** — Projektänderung (Toolchain-Angabe, Basis-Image, Gradle-Version), kein Gegenstand von `u-boot`.
- **Die Vier-Stellen-Pin-Politik für den JDK-Schlüssel** — [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) trägt die zwei Go-Schlüssel und ist durch die CI-Pins des u-boot-Projekts begründet (Root-`Dockerfile`, `repo.mk`; [`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) Entscheidung 4); ein JDK hat diese Kopplung nicht, seine Pin-Stelle ist die Template-`ARG` (Default-Pin) und die Übersteuerung des Projekts. Ob die Politik sich auf den neuen Schlüssel erweitert, ist CR-Wortlaut; Vorschlag: nein.
- **Der Lastenheft-CR selbst** — ein eigener CR-Slice nach dem Muster von [`slice-cr-devcontainer-gepinnte-werkzeugkette`](../done/slice-cr-devcontainer-gepinnte-werkzeugkette.md), noch nicht angelegt und kleiner als jener: Typ, Hybrid-Versionsfrage und Abweisungen stehen im Vertrag 0.30.0; der CR trägt nur Schlüssel-Name und Listen-Zeile der Beispiele in [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), einen Absatz der Zusage (Muster der Zeile 424), die Schlüssel-Namen in [`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config) (Zeilen 522 und 544) und die Konflikt-Regel der §4-Frage 1 — falls der Projektinhaber ihn verlangt, liegt er vor diesem Slice (§4).

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

- [ ] **Katalogeintrag JDK des Typs „gepinnte Werkzeugkette"** (nach [`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) Entscheidung 1; Name `java-toolchain`, Konflikt-Regel externes `java` + `java-toolchain` je Exit-Code `10` — entschieden 2026-11-02, §4): dritte Quellen-Konstante und ihre URL-Form in `devcontainer_toolchain.go` (Typ-Grenze Zeilen 159 bis 166, Normalisierer und URL-Leser), Katalog-Eintrag mit `checksumSource` neben Zeilen 338 bis 347, Dockerfile-Block in `Dockerfile.tmpl` nach dem Muster `.GoToolchain`/`.GolangciLint` (Zeilen 138 bis 194: `ARG` für Version und beide Prüfsummen je Form, `sha256sum -c` vor dem Entpacken, `JAVA_HOME` und `PATH`), Golden Cases (Default-Pin und `version:`-Übersteuerung mit gelesenen Prüfsummen), Docker-Test: der Bau läuft und `java -version` meldet den gerenderten Pin, Negativprobe mit verfälschter Prüfsumme bricht vor dem Entpacken ab. *Beleg:* Tabellentest des Katalogs, Golden Case, Docker-Test im CI-Lauf `integration-docker`; die Prüfsummen sind am Tag der Messung gegen die offizielle Prüfsummen-Quelle gelesen, nicht aus dem Gedächtnis ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)).
- [ ] **Handbuchbeispiel `java.version`** und die Messung, die es trägt: erst die Messung, welche Tags das `java`-OCI-Feature hat (Datum und Dev-Container-CLI-Version im Protokoll unter `docs/messungen/<kennung>/`, [`harness/conventions.md`](../../../../harness/conventions.md#ablage-von-messprotokollen-und-skripten)), dann nennt das Handbuch nur eine Angabe, die es gibt; das Beispiel mit `21` entfällt oder wird berichtigt (Stellen in §3). **Mess-Pflicht gespalten** (Reuse gegen Neu-Messung; das Muster der Typ-Tests trägt einen Teil): **getragen von den Tests der Umsetzung** [`slice-v1-katalog-gepinnte-toolchain`](../done/slice-v1-katalog-gepinnte-toolchain.md) — Regel-Brücke und synthetic-Einträge (`TestFeatureCatalogue_ToolchainChecksumSources`, `featureVersionNeedsSourceForCat`), die E2E-Muster Bau/Negativprobe/Cache (`TestE2E_Toolchain_*`), das URL-Formular-Pin (`TestToolchainChecksumParsing`, Abschnitt *url forms*); **neu gemessen** werden nur (a) die Prüfsummen-Quelle des JDK-Vertriebs (die Prüfsummen-Datei derselben Veröffentlichung, entschieden §4; Same-Origin wie [`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv): Archiv und Prüfsumme vom selben Release), (b) die Archivform je Architektur und die Ausgaben `java -version`/`JAVA_HOME` im gebauten Container, (c) die Tags des `java`-OCI-Features für die Handbuch-Berichtigung, (d) **die Brew-Form des Vertriebs** (Auflage des Projektinhabers, 2026-11-02, §4 Frage 5): Existenz und Form des Temurin-Angebots über Homebrew (Cask `temurin`/`temurin@21` bzw. Formel-Verhältnis zum Katalog-Feature `homebrew`), im Devcontainer-Kontext (Linux), dokumentiert als Vertriebs-Vergleich neben der Archiv-Form — der Befund trägt nur den Vermerk, ob und wie brew anbietet, keine zweite Installationsform im CR (die Mess-Pflicht der Umsetzung bleibt) — nichts im Repo trägt sie.
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

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt (gelesen am 2026-11-02, wo eine Stelle steht). Eine Zeile „prüfen” nennt eine Datei, die der Planner nicht gelesen hat. Sie wird nach jeder Fix-Runde gegen `git diff --stat` abgeglichen (`AGENTS.md` §3.10); die Zeilennummern verschieben sich, der Implementer liest sie am Ist-Stand nach.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/application/devcontainer_features.go` | update | Katalog-Eintrag des neuen Schlüssels mit `dockerfileBlock: true` und eigener `checksumSource`-Konstante neben den Einträgen Zeilen 338 bis 347, Konstante neben Zeilen 206 bis 221; die Warnungen (Meldungen des Erzeugens, Muster `toolchainWarnings`) führen den neuen Schlüssel; die Regel-Brücke (`featureVersionNeedsSourceForCat`, Zeilen 844 bis 852) ist schon allgemein |
| `internal/hexagon/application/devcontainer_toolchain.go` | update | Typ-Grenze `toolchainChecksumSourceSupported` (Zeilen 159 bis 166) um die dritte Quellen-Konstante erweitern, Quellen-Schalter im URL-Leser (Zeilen 257 bis 269), Version-Form im Normalisierer (Zeilen 237 bis 255), Pin-Slot und ggf. Konflikt-Meldung `java` + Schlüssel nach dem Muster `goToolchainConflictMessage` (Zeilen 181 bis 192, je §4) |
| `internal/hexagon/application/generate.go` (Templates-Daten), `templates.go` | update (prüfen) | die Template-Felder stehen in `templates.go` Zeilen 86 bis 95 (`GoToolchain`, `GolangciLint` und je Pin-Feld), gereicht in `generate.go` Zeilen 848 bis 852; die Öffnung des dritten Blocks ist nicht im Einzelnen gelesen |
| `internal/hexagon/application/templates/devcontainer/Dockerfile.tmpl` | update | Block nach dem Muster `.GoToolchain`/`.GolangciLint` (Zeilen 138 bis 194): doppelte `ARG`-Form (Übersteuerung mit `.JdkPin`-Feld, sonst literaler Default-Pin — die gleiche zweite Stelle, die der Pin-Abgleich des Musters liest), `sha256sum -c` vor dem Entpacken, Entpacken, `JAVA_HOME` und `PATH` |
| Tests in `internal/hexagon/application/` (Muster `generate_toolchain_test.go`, `devcontainer_features_test.go`), `export_test.go` | update / neu | Tabellentest des Katalogs samt Quellen-Angabe, synthetic-Einträge über die Regel-Brücke, Golden Case je Profil und Übersteuerung; die Namen wählt der Implementer, die Beleg-Tabelle der Spec trägt sie |
| `internal/e2e/toolchain_docker_test.go` oder ein neuer Docker-Test | update / neu | Bau mit dem Block und `java -version` im Container, Negativprobe der Prüfsumme (Muster `TestE2E_Toolchain_PinMismatchBreaksBuild`); Aufräumen wie die Nachbartests |
| `tools/messung/<name>/`, `docs/messungen/<kennung>/` | neu | Messung der Prüfsummen-Quelle des JDK-Vertriebs, der Release-Archiv-Form und der Tags des `java`-Features (Datum, CLI-Version, Zahl der Läufe) |
| `spec/lastenheft.md` | update (eigener CR-Slice, vor diesem) | Listen-Zeile in [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Beispiele, Zeilen 402 bis 416), Absatz der Zusage (Muster Zeile 424), Konflikt-Regel; [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) bleibt bei den zwei Go-Schlüsseln, Vorschlag §1 |
| `spec/spezifikation.md` | update | **geliefert im CR** (kein erneuter Eingriff an diesen Stellen): die Schlüssel-Namen in [`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config) (Zeilen 522 und 544) und die [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion)-Hosts-Zeile sind mit dem CR-Diff berichtigt, [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) ist unverändert (gelesen am CR-Stand, CR-Slice §3); **hier bleibt**: die Java-Zeilen der [`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)-Beleg-Tabelle nach dem Code (F6 entschieden: erweitert, keine eigene Stelle daneben) |
| `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md`, `harness/replay.md` | update (prüfen) | Handbuch-Beispiel `java.version 21` (Zeilen 1087 bis 1088, 1113, 1304 bis 1305, gelesen), Tabellenzeilen der Features (Zeilen 1134 und 1140), Egress-Zeile (Zeile 1682); die Namenslisten der READMEs und die Golden-Case-Ablage in `harness/replay.md` sind nicht gelesen; ohne Kennungen im Nutzertext |
| `spec/architecture.md`, `docs/plan/adr/`, `Makefile`, `.golangci.yml`, `.d-check.yml` | **unverändert** | Accepted-ADRs sind nicht änderbar; [`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) trägt den Typ, ein weiterer Eingriff braucht eine Folge-ADR; Gates und Schwellen bleiben ([`AGENTS.md`](../../../../AGENTS.md) §3.6) |
| diese Datei, `docs/plan/planning/open/README.md`, `docs/plan/planning/in-progress/roadmap.md` (State-Berichtigung nach dem Move, Commit `767c1122`), `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl` (die Abgrenzung: kein `containerEnv`-Wert des dritten Blocks), `cmd/uboot/toolchain_test.go` (der Nutzerpfad-Test zum Konflikt) | update | Zeile dieses Slice; Moves per `make slice-mv`; **Nachtrag nach Review F-3** — die drei Zeilen fehlten in der Erst-Tabelle |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Entscheidungen des Projektinhabers** (Hauptagenten-Gespräch, 2026-11-02 — eingegetragen je Frage mit dem Ausgang; der Implementer wählt nicht). Die frühere Frage 1 (Weg A: feste Version im Template / Weg B: wählbar aus gepinnter Tabelle / Weg C: nur Handbuch) ist durch [`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) §§Entscheidung und §§Konsequenzen beantwortet (Hybrid; Entscheidung 2026-11-02; CR in Lastenheft 0.30.0, Historie-Zeile), siehe Befund §1.

1. **Schlüsselname und Verhältnis zum Katalog-Schlüssel `java` — entschieden: `java-toolchain`, Abweisung mit Exit-Code `10`.** Der Name `java` ist im Katalog vergeben: ein OCI-Feature (`ghcr.io/devcontainers/features/java`, Standardversion `1`, „Java + SDKMAN“; Fundstelle `internal/hexagon/application/devcontainer_features.go` Zeilen 273 bis 277), dessen Hosts `repo.maven.apache.org` und `repo1.maven.org` in der Default-Allowlist stehen ([`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), von einer Messung der Installations-Hosts getragen, die dieser Planner nicht wiederholt hat). Der neue Schlüssel trägt den Familien-Namen der zwei Go-Schlüssel: `java-toolchain`. **Konflikt-Regel (Muster S8):** aktiviertes externes `java` neben `java-toolchain` ist ein fachlicher Validierungsfehler mit Exit-Code `10` (Muster `goToolchainConflictMessage`, `internal/hexagon/application/devcontainer_toolchain.go` Zeilen 181 bis 192 — welche der zwei Installationen die Befehle finden, entscheidet die `PATH`-Reihenfolge); das externe `java` allein bleibt unverändert wählbar. Träger: CR.
2. **JDK-Vertrieb und Archivform — entschieden: Eclipse Temurin, LTS 21.** Die Release-Archive liegen je Architektur in der Veröffentlichung der Temurin-Binaries (Archivform und Lizenzform liest die Messung am gewählten Stand); die konkrete Versions-Zahl ist die Default-Pin-Messung (Frage 3), nicht eine Vertragszahl.
3. **Default-Pin-Werte — Vertrieb entschieden (Frage 2), konkrete Version liest die Messung.** Das JDK der Basis-Images und `jvmToolchain` eines Projekts müssen zusammenpassen. **Die Zahlen stehen im CHANGELOG (als Chronik), nicht im Vertragstext** — Muster der Freigabe F4 des CR-Vorschlags [`cr-vorschlag-gepinnte-werkzeugkette`](../done/cr-vorschlag-gepinnte-werkzeugkette.md) §Stelle 7 (Stand als Chronik in CHANGELOG und Lastenheft-Historie, Muster Ollama: „mit der gepinnten Version von Ollama” ohne Zahl).
4. **Prüfsummen-Quelle — entschieden: Same-Origin je Release.** Die offizielle Prüfsummen-Quelle des Schlüssels ist die Prüfsummen-Datei derselben Veröffentlichung, in der das Release-Archiv liegt (Muster der zwei Go-Seiten, [`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)); die JSON-Schnittstelle `api.adoptium.net` bleibt aus — zusätzlicher Host ohne Same-Origin-Garantie. Ob die gewählte Form eine beliebige nutzergewählte Version abdeckt, liest die Messung; deckt sie sie nicht, bleibt dort die Abweisung ([`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) Re-Evaluierungs-Trigger).
5. **Auflage: die Brew-Form des Vertriebs prüfen.** Die Mess-Pflicht (§2) bekommt einen Zusatzpunkt: die Messung liest die Brew-Form des Weges — Existenz und Form des Temurin-Angebots über Homebrew (Cask `temurin`/`temurin@21` beziehungsweise das Formel-Verhältnis im Katalog-Feature `homebrew`), im Devcontainer-Kontext (Linux) — und dokumentiert den Befund als Vertriebs-Vergleich neben der Archiv-Form. Keine zweite Installationsform im CR ohne die Messung: der Brew-Befund trägt nur den Vermerk, ob und wie brew den Vertrieb anbietet; er begründet keine zweite Zusage.

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; der CR-Slice [`slice-cr-devcontainer-java-toolchain`](../done/slice-cr-devcontainer-java-toolchain.md) (Listen-Zeile, Absatz der Zusage, Konflikt-Regel, [`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)-Namen; Muster [`slice-cr-devcontainer-gepinnte-werkzeugkette`](../done/slice-cr-devcontainer-gepinnte-werkzeugkette.md)) liegt in `done/` (beobachtbar: die Datei liegt dort) und das Lastenheft nennt den Schlüssel; kein anderer Slice in `in-progress/` (WIP-Limit 1, gelesen beim Start).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): trägt die Messung der Prüfsummen-Form (Frage 4) Abweichungen je Architektur, die eigene Docker-Tests je Form brauchen, oder zeigt der Brew-Befund (Frage 5), dass der CR eine zweite Vertriebsform braucht, geht der Teil in einen eigenen Slice; die Teilung entscheidet der Projektinhaber.
- `in-progress` → `open` (blockiert): der gewählte Vertrieb liefert für eine der beiden Architekturen kein Archiv mit veröffentlichter Prüfsumme in derselben Quelle (Same-Origin) oder deckt die gewählte Version mit ihr nicht ab — an der Stelle bleibt bis zum Folge-Entscheid die Abweisung ([`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) Re-Evaluierungs-Trigger).

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

1. **Prüfsummen aus dem Gedächtnis oder von einer Spiegelseite.** — **Ausgang:** offen bis zur Closure. Vorsorge: Prüfsummen werden am Tag der Messung von der offiziellen Prüfsummen-Quelle gelesen und im Protokoll datiert; ein Reviewer vergleicht sie im frischen Kontext erneut.
2. **Das JDK des Containers und das gepinnte JDK weichen ab** (das Basis-Image oder ein aktiviertes Feature bringt ein eigenes JDK, die `PATH`-Reihenfolge entscheidet). — **Ausgang:** offen bis zur Closure. Vorsorge: der Docker-Test prüft `java -version` und `JAVA_HOME` im gebauten Container; die Wechselwirkung mit dem Eintrag `java` ist die §4-Frage 1 (Muster `go` + `go-toolchain`: eine `PATH`-Frage, die kein Test trägt).
3. **Eine Zusage reicht weiter als die Messung** (ein JDK-Archiv gemessen auf einer Architektur oder einem Host). — **Ausgang:** offen bis zur Closure. Vorsorge: Messsatz mit Umgebung, Datum und Zahl der Läufe; keine Aussage über nicht gemessene Hosts ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)).
4. **Eine Bedingung in Spezifikation oder Handbuch hat keinen Test** (zum Beispiel „in beiden Profilen” oder die Schlüssel-Namen-Zeilen). — **Ausgang:** offen bis zur Closure. Vorsorge: Beleg-Tabelle nach dem Code (Muster [`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) unten), jeder Teilsatz mit eigenem Pin ([`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)).
5. **Der Docker-Test braucht Netz** (Download des Archivs) und hängt am Cache des Hosts; der Block lädt aus dem Host des gewählten Vertriebs (ein neuer Host neben der Basis), und eine Negativprobe hinterlässt den Zwischencontainer des fehlgeschlagenen Baus. — **Ausgang:** offen bis zur Closure. Vorsorge: Aufräumen wie die Nachbartests, Nachweis ist der CI-Lauf ([`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md), [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md)).
6. **Plan-Aussagen über Code, die am Fundort nicht gelesen sind** — im Befund und in der Tabelle von §3 Zahlen und Verhalten über den Umsetzungs-Diff der gepinnten Werkzeugkette hinaus — sind die bekannteste Fundstellen-Klasse dieses Plans ([`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md); die Umsetzungs-Closure [`slice-v1-katalog-gepinnte-toolchain`](../done/slice-v1-katalog-gepinnte-toolchain.md) §7 nennt je Berichtigungs-Commit wiederholte Fundstellen-Korrekturen, weil die Korrektur selbst die Zeilen-Zahlen verschiebt). — **Ausgang:** offen bis zur Closure. Vorsorge: jede Code-Aussage dieses Plans trägt die Fundstelle dieses Tages (Zeilen am Stand 2026-11-02 gelesen, hier genannt); der §3.10-Abgleich (Tabelle gegen `git diff --stat`) läuft nach jeder späteren Berichtigungs-Runde erneut, nicht nur vor dem ersten Commit ([`AGENTS.md`](../../../../AGENTS.md) §3.10).

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

**Vorgelagert — offene Beobachtungen sichten** (nur die in diesem Plan zitierten Einträge sind gelesen, ihre Zähler-Stände nicht neu erhoben — beim Start durchgehen): [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (Risiko 3), [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (Risiko 4), [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (Risiko 6, gelesen in der Umsetzungs-Closure [`slice-v1-katalog-gepinnte-toolchain`](../done/slice-v1-katalog-gepinnte-toolchain.md)), [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md) und [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md) (Risiko 5). Evidenz-/Diskrepanz-Risiko insgesamt: mittel (die Tags des `java`-Features, die Form der JDK-Prüfsummen-Quelle und das Zusammenliegen von `java` und dem neuen Schlüssel sind nicht gemessen).

**Modus-Begründungsblock:** alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
