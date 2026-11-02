# CR-Vorschlag: Katalog-Schlüssel `java-toolchain` (Temurin, LTS 21) — Alt/Neu-Wortlaut je Stelle

**Status:** vorgelegt dem Projektinhaber (Freigabe steht aus; 2026-11-02).

**Freigabe-Kontext (Projektinhaber, 2026-11-02):** die vier Entscheidungen des
Hauptagenten-Gesprächs (§4 des Plans
[`slice-v1-katalog-feature-jdk-gepinnt`](../open/slice-v1-katalog-feature-jdk-gepinnt.md))
tragen diesen Entwurf: **Schlüsselname und Konflikt-Regel** — `java-toolchain`,
aktiviertes externes `java` neben dem Schlüssel ist ein fachlicher
Validierungsfehler (Exit-Code `10`), das externe `java` allein bleibt wählbar;
**Vertrieb** — Eclipse Temurin, LTS 21, Release-Archive aus der Veröffentlichung
der Temurin-Binaries; **Pin-Zahlen** — nicht im Vertragstext, im CHANGELOG und in
der Lastenheft-Historie als Chronik; **Prüfsummen-Quelle** — Same-Origin je
Release, die JSON-Schnittstelle `api.adoptium.net` bleibt aus. Dazu die
**Brew-Auflage** (2026-11-02): die Messung des Umsetzungs-Slices liest die
Brew-Form des Vertriebs und trägt nur den Befund — keine zweite
Installationsform im CR. Diese Setzungen sind der Freigabe-**Kontext**, nicht
der freigegebene **Wortlaut**; der *Neu*-Wortlaut der Stellen §§2 ist je Stelle
freizugeben, und `spec/lastenheft.md` wird erst damit geändert, wortgleich
(Muster [`slice-cr-devcontainer-forge-cli`](../done/slice-cr-devcontainer-forge-cli.md)
§1: Freigabe als eigener Schritt vor der Lastenheft-Änderung). Der *Neu*-Wortlaut
ist so formuliert, dass er nach dem Kopieren ohne weitere Bearbeitung
contract-konform ist; die Belege stehen im Kommentar je Stelle, nicht im
Vertrag. Die J-Marker (**J1**..**J11**) sind Entwurfs-Kennungen dieses Dokuments,
keine Vertragskennungen; sie entfallen mit dem Lastenheft-Commit.

**Gegenstand:** die Lieferung zu DoD-Zeile 1 des CR-Slices
[`slice-cr-devcontainer-java-toolchain`](slice-cr-devcontainer-java-toolchain.md)
(„Wortlaut-Entwurf je Stelle und Freigabe" — die Stellen (a) bis (f) der Zeile).
Diese Datei ist kein Slice und trägt keinen Lebenszyklus; sie ändert nichts an
den Verträgen und diskutiert die Entscheidungen 2026-11-02 nicht neu
(CR-Disziplin §1 des CR-Slices). Ihr Zustand folgt dem CR-Slice.

**Welle:** ohne Welle; der CR-Slice schlägt Release `v0.20.0` vor (die
Feature-Katalog-Welle der Roadmap); die Versionsvergabe schneidet der
Release-Cut.

**Bezug:** [`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md)
(`Accepted`, 2026-11-02 — der Eintragstyp; wird verlinkt, nicht wiederholt),
[`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)
(der Eintragstyp mit den Sätzen S1 bis S11, seit Lastenheft 0.30.0 — der neue
Schlüssel erbt den Typ),
[`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette)
(Pin-Politik — bleibt bei den zwei Go-Schlüsseln, Regel 4 unten),
[`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)
(die Abweisungs-Fälle enden je mit Exit-Code `10`),
[`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-und-golangci-lint-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)
(der Eintragstyp in der Technik),
[`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen),
[`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion),
[`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config),
[`SPEC-050`](../../../../spec/spezifikation.md#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen).
Die Mess-Basis der Werkzeugkette
([`docs/messungen/slice-v1-katalog-gepinnte-toolchain-entwurf/README.md`](../../../messungen/slice-v1-katalog-gepinnte-toolchain-entwurf/README.md))
trägt die Mechanik (Same-Origin, Exit-Codes, Meldungs-Form); die JDK-Neuheiten
(Temurin-Archivform, `JAVA_HOME`/`java -version`, Brew-Befund) sind
Umsetzungs-Messung
([`slice-v1-katalog-feature-jdk-gepinnt`](../open/slice-v1-katalog-feature-jdk-gepinnt.md)
§2) — der CR trägt Zusage nur soweit ein Docker-Test sie belegt.

**Verantwortlich:** — (Freigabe: Projektinhaber).

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-11-02.

---

## 1. Regelrahmen

Fünf Regeln, an denen jede Stelle unten gemessen ist:

1. **Decken-Regel.** Der *Neu*-Wortlaut als künftiger Lastenheft-Text trägt
   keine ADR-, Slice-, Messungs- oder Planning-Kennungen ([`AGENTS.md`](../../../../AGENTS.md)
   §3.4; die Referenzmatrix verbietet die Kanten aus dem Vertrags-Stratum nach
   unten). Wo der Vertrag schweigt und der
   [`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md),
   die Messung oder eine der vier Entscheidungen die Einzelheit trägt, steht
   die Belegstelle im Kommentar — nicht im Blockquote.
2. **Keine Zusage über den Typ, die Messung und die Entscheidungen hinaus.**
   Der Schlüssel erbt die Zusagen des Eintragstyps (Lastenheft-Zeile 424; die
   Fundstellen stehen in den Belegen), nicht mehr. Die Temurin-Details sind vor
   der Umsetzungs-Messung Form, nicht Zusage; die Grenz-Sätze der Stellen
   folgen dem Muster „nur, soweit ein Test mit Build-Tag `docker` es belegt"
   ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md);
   die pinnenden Tests der Beleg-Tabelle trägt der Umsetzungs-Diff —
   [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)).
3. **Keine Zahlen im Vertragstext.** Die konkrete Temurin-Revision (das gepinnte
   Release von LTS 21) ist eine Messung des Umsetzungs-Slices; der Vertrag nennt
   die LTS-Version und die Veröffentlichungs-Form, die Zahlen stehen im
   CHANGELOG und in der Lastenheft-Historie als Chronik (Muster der Freigabe F4
   des CR-Vorschlags
   [`cr-vorschlag-gepinnte-werkzeugkette`](cr-vorschlag-gepinnte-werkzeugkette.md)
   §Stelle 7; im Lastenheft 0.30.0 trägt die Historie-Zeile 1627 die Werte als
   Chronik — gelesen am 2026-11-02).
4. **[`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette)
   bleibt bei den zwei Go-Schlüsseln.** Gelesen 2026-11-02 (Fundstelle
   [`spec/lastenheft.md`](../../../../spec/lastenheft.md) Zeilen 520 bis 522:
   „die Schlüssel `go-toolchain` und `golangci-lint`"). Eine Erweiterung wäre
   ein neuer Wortlaut-Anker und läge in keiner der vier Entscheidungen; die
   Pin-Stellen des JDK-Schlüssels sind die Template-`ARG` (Default-Pin) und die
   Übersteuerung des Projekts (CR-Slice §1; Vorschlag des Plans §1 — beide
   gleichlautend).
5. **Grabstein-frei.** Der Typ ist vertraglich lebendig (Lastenheft 0.30.0);
   der neue Schlüssel erzeugt keine neue Kennung im Lastenheft — er fügt der
   Beispiel-Liste eine Zeile und dem Abschnitt einen Absatz hinzu (Stellen 1
   und 2). Die Abweisungen sind dieselben des Typs, angewandt auf den
   dritten Schlüssel — keine zweite Konflikt-Kennung
   (Rückführungs-Trigger (a) des CR-Slices §4 würde greifen, wollte der
   Abgleich sie verlangen).

## 2. Stellen

### Stelle 1 — [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features): die Listen-Zeile

**Alt (gelesen 2026-11-02):** die Beispiel-Liste (Fundstelle
[`spec/lastenheft.md`](../../../../spec/lastenheft.md) Zeilen 402 bis 416)
trägt „Go-Toolchain aus gepinntem Release (`go-toolchain`)" (Zeile 414) und
„golangci-lint mit gepinntem Release (`golangci-lint`)" (Zeile 415) nebeneinander,
gefolgt von „Anschluss an einen Ollama-Dienst des Hosts (`ollama-host`)"
(Zeile 416). Ein Java-Schlüssel steht in der Liste nicht; „Java" (Zeile 407)
ist die Beispiel-Zeile des externen OCI-Features.

**Neu (Vorschlag):** eine Zeile einfügen zwischen Zeile 415 und Zeile 416 —
die Toolchain-Zeilen stehen beieinander, der neue Schlüssel schließt die
Familie an:

> - JDK aus gepinntem Release (`java-toolchain`)

**Beleg:** die Position folgt der Ordnung der zwei Go-Zeilen (gelesen Zeilen
414 bis 416); der Schlüssel-Name `java-toolchain` ist entschieden (2026-11-02,
§4 Frage 1 des Plans). Der Name der Listen-Zeile ist der gleiche, den
[`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)
(Zeilen 522, 544) nach der Erweiterung trägt (Stelle 10).

### Stelle 2 — [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features): der Zusage-Absatz des Schlüssels

**Alt (gelesen 2026-11-02):** nach dem `gh`/`glab`-Absatz (Zeile 422) trägt
der Absatz der zwei Go-Schlüssel (Zeile 424, Sätze S1 bis S11) die Typ-Zusage —
in beiden Profilen wählbar, Release-Archiv mit geprüfter SHA-256 je
Architektur, Default-Pin als Build-Argument des Katalog-Templates, `version:`
Übersteuerung nur mit der offiziellen Prüfsummen-Quelle (Lese-Prüfung bei der
Erzeugung, gelesene Prüfsummen im erzeugten Block), die Abweisungen je
Exit-Code `10`, die Ausweisung beim Erzeugen — und danach (Zeile 426) folgt
„Auf Wunsch des Nutzers erzeugt ...". Ein Absatz für den Java-Schlüssel steht
nicht.

**Neu (Vorschlag):** ein Absatz nach Zeile 424, vor „Auf Wunsch des Nutzers
erzeugt“ (Zeile 426):

> `java-toolchain` (JDK aus gepinntem Release) ist in beiden Profilen (**J1**)
> als eingebautes Feature wählbar (`u-boot config set
> devcontainer.features.java-toolchain.enabled true`). (**J2**) Der Schlüssel
> installiert im erzeugten Devcontainer-Dockerfile ein JDK der
> Eclipse-Temurin-Distribution der LTS-Version 21 aus einem gepinnten
> Release-Archiv der Veröffentlichung der Temurin-Binaries und prüft die
> SHA-256-Prüfsumme des Archivs je Architektur, bevor es entpackt wird; eine
> abweichende Prüfsumme bricht den Bau vor dem Entpacken ab. (**J3**) Die
> Version des Schlüssels steht als Build-Argument im erzeugten Dockerfile und
> stammt aus dem Katalog-Template (Default-Pin); sie wird nur mit einem
> u-boot-Release geändert. (**J4**) Ein Projekt kann die Version übersteuern
> (`u-boot config set
> devcontainer.features.java-toolchain.version <version>`). (**J5**) Die
> Übersteuerung wird nur angenommen, wenn der Eintrag seine offizielle
> Prüfsummen-Quelle angibt — für diesen Schlüssel die Prüfsummen-Datei
> derselben Veröffentlichung der Temurin-Binaries, aus der das Release-Archiv
> stammt, je gewählter Version. Die Quelle wird bei der Erzeugung gelesen,
> und die gelesenen Prüfsummen trägt der erzeugte Dockerfile-Block. Ist die Quelle nicht lesbar oder überdeckt sie die
> gewählte Version nicht, schlägt die Erzeugung der Übersteuerung mit einem
> fachlichen Validierungsfehler fehl
> ([`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes): Exit-Code `10`; die nicht
> lesbare Quelle ist kein Feature-Source-Fehler, die Eingabe der
> Übersteuerung selbst ist unbrauchbar). (**J6**) Ein Eintrag der gepinnten
> Werkzeugkette ohne Angabe seiner offiziellen Prüfsummen-Quelle weist eine
> `version:`-Angabe mit einem fachlichen Validierungsfehler ab
> ([`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes): Exit-Code `10`). (**J7**)
> Die Archive des Typs sind nicht signiert, und `u-boot` prüft keine
> Signatur; die Grenze der Prüfung weist `u-boot` beim Erzeugen aus. (**J8**)
> Die gleichzeitige Aktivierung des externen Features `java` und des
> Schlüssels `java-toolchain` ist ein fachlicher Validierungsfehler
> (Exit-Code `10`); das externe `java` allein bleibt unverändert wählbar.
> (**J9**) `u-boot` weist beim Erzeugen aus: den aktivierten Block des
> Schlüssels, ob die Version der Default-Pin oder eine Übersteuerung ist und,
> bei einer Übersteuerung, die Prüfsummen-Quelle. (**J10**) Zugesagt ist die
> Erzeugung und der Bau des Schlüssels auf einem Linux-Host (`amd64`) mit
> lokaler, nicht verschachtelter Docker-Engine und nur soweit, wie ein Test
> mit Build-Tag `docker` es belegt; für andere Hosts, Engines und
> Architekturen wird nichts zugesagt. (**J11**) Nicht zugesagt: das
> Arbeitsverhalten des JDK in einer per `version:` übersteuerten Version;
> das Überstehen eines Neubaus des Images durch die Installation, soweit
> kein Test es belegt; der Schlüssel trägt unter `containerEnv` kein
> Cache-Verzeichnis (Stelle 4).

**Beleg je Satz:** **J1** — das Muster des Typs (Zeile 424, erster Satz); die
Dockerfile-Block-Einträge wirken in beiden Profilen
([`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen)
Zeile 1555 „Anders als `ollama` wirken sie in beiden Profilen"); die
Aktivierung ist profil-unabhängig (Fundstelle `dockerfileBlockFeatureActive`,
`internal/hexagon/application/devcontainer_features.go` Zeilen 576 bis 589 —
kein Profil-Schalter). **J2** — der Vertrieb und die LTS-Version sind
entschieden (2026-11-02, §4 Frage 2); das Architektur-Muster trägt der
Go-Block ([`Dockerfile.tmpl`](../../../../internal/hexagon/application/templates/devcontainer/Dockerfile.tmpl)
Zeilen 138 bis 165: zwei Prüfsummen-`ARG` je Architektur Zeilen 149 und 150,
`sha256sum -c` vor dem Entpacken Zeile 161); die konkrete Archivform
(Namensmuster, Lizenzform) liest die Messung des Umsetzungs-Slices (Plan §2) —
der Satz sagt sie nicht (Regel 3); „eine abweichende Prüfsumme bricht den Bau
vor dem Entpacken ab" ist die Typ-Zusage (Lastenheft-Zeile 424, zweiter
Satz — und die Negativprobe des
[`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md)
§Fitness Function). **J3/J4** — das Muster der Sätze S3/S4 (Zeile 424) und
[`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md)
Entscheidung 2, erster Halbsatz. **J5** — die Bedingung ist die Typ-Regel
(Entscheidung 2: „trägt er sie, geht eine `version:`-Übersteuerung ohne
`source:` durch, und der Block prüft die Prüfsumme der nutzergewählten Version
gegen diese Quelle"); die Quelle je Schlüssel konkretisiert Stelle 3
(Same-Origin je Release, entschieden 2026-11-02 — Frage 4); die
Exit-Code-Zuordnung `10` ist der freigegebene Stand der Typ-Zusage
(Entscheidung der Vorlage F3, 2026-11-02: Exit-Code `10` je Fehlerfall — der
Wortlaut der zwei Schlüssel trägt ihn, Fundstelle Zeile 424, fünfter Satz) und
die nicht-lesbar/nicht-überdeckend-Form ist der Wortlaut derselben Stelle.
**J6** — das Muster von S6 (Zeile 424, sechster Satz); die Regel ist
schlüsselfrei im Code (Fundstelle `featureVersionNeedsSourceForCat`,
`devcontainer_features.go` Zeilen 844 bis 852; die Abweisung eines aktiven
Blocks ohne Quellen-Angabe `resolveToolchain`, `devcontainer_toolchain.go`
Zeilen 219 bis 225). **J7** — das Muster von S7 (Zeile 424, siebter Satz) und
die Grenz-Form der Block-Kommentare
([`Dockerfile.tmpl`](../../../../internal/hexagon/application/templates/devcontainer/Dockerfile.tmpl)
Zeilen 144 bis 146: „not signed ... does not protect against a release that
was harmful from the start"); [`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md)
Entscheidung 3. **J8** — entschieden (2026-11-02, Frage 1); das Muster der
Abweisung trägt `goToolchainConflictMessage`
(`internal/hexagon/application/devcontainer_toolchain.go` Zeilen 181 bis 192,
Text der Abweisung mit den zwei Auswegen) und `goFeatureActive` (Zeilen 168
bis 179); der Wortlaut ist der von S8, auf die zwei Java-Schlüssel angewandt —
eine zweite Konflikt-Regel entsteht nicht (Regel 5). **J9** — das Muster von
S9 (Zeile 424, achter Satz); [`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md)
Entscheidung 3 (letzter Satz). **J10/J11** — das Muster von S10/S11 (Zeile
424, neunter und zehnter Satz) und der Grenz-Satz-Muster der Historie-Zeilen
0.27.0/0.28.0 (Fundstelle Zeilen 1624 bis 1625); „nur soweit ein Test mit
Build-Tag `docker` es belegt" trägt
[`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md);
der `containerEnv`-Null-Beitrag ist die Abgrenzung zu der Form der zwei
Go-Schlüssel (Stelle 4), im selben Satz-Strang wie der S11-Muster der zwei
Schlüssel die dort nicht genannten Umgebungswerte nennt.

### Stelle 3 — die offizielle Prüfsummen-Quelle des Schlüssels

Der [ADR-0041](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md)
lässt die Quelle je Schlüssel in der Tabelle der Spec stehen; der Vorschlag,
gelesen am Same-Origin-Muster der zwei Schlüssel (die Prüfsummen kommen aus
derselben Veröffentlichung wie das Archiv — Kommentar des Templates,
[`Dockerfile.tmpl`](../../../../internal/hexagon/application/templates/devcontainer/Dockerfile.tmpl)
Zeilen 139 bis 146):

| Schlüssel | Offizielle Prüfsummen-Quelle (Vorschlag) | Begründung (gelesen) |
|---|---|---|
| `java-toolchain` | die Prüfsummen-Datei derselben Veröffentlichung der Temurin-Binaries, aus der das Release-Archiv stammt, je gewählter Version | entschieden (2026-11-02, §4 Frage 4 des Plans) — Same-Origin am Archiv wie die zwei Schlüssel ([`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-und-golangci-lint-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) Tabelle, Fundstelle [`spec/spezifikation.md`](../../../../spec/spezifikation.md) Zeilen 983 bis 986); die JSON-Schnittstelle `api.adoptium.net` bleibt **ausdrücklich** nicht Quelle: zusätzlicher Host ohne Same-Origin-Garantie (derselben Entscheidung) |

**Abweisungsfall:** der Schlüssel trägt seine Quelle (die Tabelle oben, Träger
Satz **J5**); `api.adoptium.net` ist kein Fall der Quelle, sondern ihr
ausdrücklich ausgeschlossener Ersatz — der Abgleich der Übersteuerung liest
nur die Quelle der Tabelle. Ob die Prüfsummen-Datei einer Veröffentlichung
eine beliebige nutzergewählte Version abdeckt, liest die Messung des
Umsetzungs-Slices; deckt sie eine gewählte Version nicht, bleibt dort die
Abweisung ([`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md)
Re-Evaluierungs-Trigger; Rückführung §4 (b) des CR-Slices). Der Meldungs-Wortlaut
der Quelle (die Konstante verdoppelt ihn als Meldungstext — Fundstelle
`devcontainer_features.go` Zeilen 215 bis 217 „The values double as the wording
of the generate notices") und die URL-Form des Abrufs je Version sind
Umsetzungs-Stellen (Muster der URL-Bauer der zwei Quellen,
`devcontainer_toolchain.go` Zeilen 261 bis 269, gelesen am Schalter).

**Was der Docker-Test belegt, was Messung bleibt:** belegt wird, dass der
erzeugte Block mit einer `version:`-Übersteuerung gegen die für die gewählte
Version gelesenen Prüfsummen baut und dass eine verfälschte Prüfsumme den Bau
vor dem Entpacken abbricht (Negativprobe, Muster
[`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-und-golangci-lint-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)
Beleg-Tabelle Zeile 1001) — die Prüf-Form, nicht die Quelle im Einzelnen; die
URL-Formen der Temurin-Veröffentlichung pinnt der Docker-Test auf den dritten
Schlüssel, für andere Umgebungen als die der Messung bleibt der Abruf
unbelegt (Muster der Grenze, ebenda Absatz zur URL-Form, Zeile 988).

### Stelle 4 — `containerEnv`-Form: kein Cache-Wert des JDK-Schlüssels

**Alt (gelesen 2026-11-02):** die zwei Go-Schlüssel tragen unter `containerEnv`
der erzeugten `devcontainer.json` die Cache-Verzeichnisse `GOCACHE` und
`GOMODCACHE` unterhalb von `${containerWorkspaceFolder}` (Lastenheft-Zeile 385;
[`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-und-golangci-lint-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)
Abschnitt *Cache-Werte*, Fundstelle [`spec/spezifikation.md`](../../../../spec/spezifikation.md)
Zeile 992). Die Render-Stelle trägt den Zustand je Schlüssel:
[`devcontainer.json.tmpl`](../../../../internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl)
Zeilen 42 bis 47 (der Block rendert bei `{{- if or .GoToolchain .GolangciLint}}`,
beide Werte untereinander).

**Neu (Vorschlag):** der Schlüssel `java-toolchain` trägt **keinen** Wert unter
`containerEnv`; der Block-Text in [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen)
wird um diesen Schlüssel **nicht** erweitert, und der Satz **J11** sagt den
Null-Beitrag als Grenze aus. `JAVA_HOME` und `PATH` stehen als `ENV`-Zeilen im
erzeugten Dockerfile-Block — Muster der `ENV PATH`-Zeile des Go-Blocks
([`Dockerfile.tmpl`](../../../../internal/hexagon/application/templates/devcontainer/Dockerfile.tmpl)
Zeile 165); der konkrete `JAVA_HOME`-Wert (Installationspfad der entpackten
Distribution) und die Form des Entpackens sind Umsetzungs-Mess-Stellen, die
der Docker-Test pinnt (`java -version` und `JAVA_HOME` im gebauten Container —
Plan [`slice-v1-katalog-feature-jdk-gepinnt`](../open/slice-v1-katalog-feature-jdk-gepinnt.md)
§2). Was die Messung entscheidet, ist also die **Form der Installations-Umgebung**
im Block, nicht ein `containerEnv`-Wert.

**Begründung:** `containerEnv` ist in diesem Eintragstyp der Träger der
Cache-Verzeichnisse, die unter dem Workspace-Volume liegen und das
Neuanlegen des Containers überstehen (die Messung M1 bis M3 der
Werkzeugkette, getragen seit Lastenheft 0.30.0); JAVA_HOME ist kein Cache,
sondern wie die `ENV PATH`-Zeile des Go-Blocks die Installations-Umgebung der
Block-Installation — sie übersteht den Image-Neubau mit dem Bild, nicht das
Recreate des Containers, und braucht keinen Eintrag in der erzeugten
Konfiguration. Die Caches von JVM-Bau-Werkzeugen (Gradle, Maven) liegen bei
deren Werkzeugen und sind nicht Bestandteil der JDK-Distribution; keine der
vier Entscheidungen 2026-11-02 trägt sie — der CR erfindet keinen Wert (Regel 2).
**Risiko des CR-Slices §6.2** (das JDK des Containers und das gepinnte JDK
weichen ab): nicht durch den `containerEnv`-Weg getragen; der Docker-Test
prüft `java -version` und `JAVA_HOME`, und die Konflikt-Frage ist die Stelle 5 —
anders als beim Go-Fall trägt die Basis (`mcr.microsoft.com/devcontainers/base:debian`,
[`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen)
Zeile 1551) kein JDK; ob weitere aktivierte Features ein JDK bringen, liest
der Umsetzungs-Diff als Beobachtung.

### Stelle 5 — Verhältnis externes `java` + `java-toolchain`

**Lage:** der Katalog führt `java` → `ghcr.io/devcontainers/features/java`,
Standardversion `1`, „Java + SDKMAN“
([`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen)
Zeile 1552; Fundstelle `devcontainer_features.go` Zeilen 273 bis 277 —
gelesen 2026-11-02) und die Konflikt-Wächter-Maschine des Typs prüft heute
allein `go` gegen `go-toolchain` (`goFeatureActive`, `devcontainer_toolchain.go`
Zeilen 168 bis 179; Meldung Zeilen 181 bis 192).

**Entschieden — Abweisung (Exit-Code `10`):** die Setzung vom 2026-11-02
(§4 Frage 1) trägt den Satz **J8**; die Gegenoption Warnung ist damit
abgewiesen, wie in der Werkzeugkette vor ihr (Abwägung in
[`cr-vorschlag-gepinnte-werkzeugkette`](cr-vorschlag-gepinnte-werkzeugkette.md)
§Stelle 5: die Warnung nimmt die nicht getestete Konfiguration an und lässt
die Wirkung der `PATH`-Reihenfolge der zwei JDK-Installationen — der Zustand,
den kein Test trägt). Das externe `java` allein bleibt unverändert wählbar
(Satz **J8**, zweiter Halbsatz); seine Zeile
[`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion)
(Zeile 1466) ändert sich nicht.

**Offen — `sdkman`:** der Katalog führt `sdkman` auf dieselbe Quelle mit
geteiltem Render-Schlüssel („teilt sich den Render-Schlüssel mit `java`",
[`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen)
Zeile 1552; Katalog-Eintrag mit der Beschreibung „SDKMAN (installed by the
official java feature)",
[`devcontainer_features.go`](../../../../internal/hexagon/application/devcontainer_features.go)
Zeilen 310 bis 314). Ob ein aktives `sdkman` allein ein JDK installiert (die
Aktivierung rendert das `java`-Feature) und deshalb unter denselben
Abweisungsfall wie das externe `java` fällt, liest der Implementer an den
Feature-Metadaten und pinnt es mit einem Test — Freigabe-Punkt **F7**;
Empfehlung: derselbe Abweisungsfall, sobald die gelesene Lage ihn trägt; ohne
sie bleibt der Absatz beim externen `java` und die Frage geht als
Beobachtung ins Register. In keiner der beiden Stände verspricht der CR, was
der Docker-Test nicht trägt (Regel 2).

### Stelle 6 — Default-Pin-Werte: Chronik, keine Vertragszahlen

**Vorschlag: die Werte stehen nicht im Vertragstext.** Der Vertrag sagt
„LTS-Version 21" und die Veröffentlichungs-Form (Sätze **J3**, **J5**); die
konkrete Revision — das gepinnte Temurin-Release von LTS 21 — ist die
Default-Pin-Messung des Umsetzungs-Slices (Entscheidung 2026-11-02, Frage 3).
Begründung aus der Hauspraxis: das Muster F4 des Vorlaufs — eine Zahl im
Lastenheft machte jede Pin-Hebung zu einem Vertrags-Commit neben dem
Code-Commit — und hier ohne die Vier-Stellen-Kopplung ([`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette)
bleibt bei den zwei Go-Schlüsseln, Regel 4), also ohne Test-Abgleich über
andere Stellen als die Template-`ARG` und die erzeugte Dockerfile-`ARG` (die
zweite entsteht mit der Übersteuerung, Muster der Pin-Zweige
[`Dockerfile.tmpl`](../../../../internal/hexagon/application/templates/devcontainer/Dockerfile.tmpl)
Zeilen 147 bis 155 und 176 bis 184 — gelesen am 2026-11-02). **Der gelesene
Stand** der zwei bestehenden Pins (Beleg): Go `1.27.3`, `golangci-lint`
`v2.14.0` in [`Dockerfile.tmpl`](../../../../internal/hexagon/application/templates/devcontainer/Dockerfile.tmpl)
Zeilen 152 bis 154 und 181 bis 183 — die Template-`ARG` ist die Quellen-Stelle der
Pin-Werte ([`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette));
für `java-toolchain` existiert die `ARG`-Stelle nicht, sie entsteht mit dem
Umsetzungs-Diff. Die Zahl steht im CHANGELOG-Eintrag der Stelle 12 und — nach
der Lastenheft-Änderung — als Chronik in der Historie-Zeile `0.31.0` (Muster
der Zeile 1627, gelesen). **Alternative** (ein Satz mit der Zahl in
[`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features))
bleibt möglich; sie ist durch die Entscheidungen 2026-11-02 (Frage 3
ausdrücklich: CHANGELOG und Historie, nicht Vertragstext) abweichend — kein
F-Punkt, der Wortlaut trägt die Setzung.

### Stelle 7 — Technik-Stelle: [`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-und-golangci-lint-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) erweitern oder eigene Stelle

**Alt (gelesen 2026-11-02):**
[`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-und-golangci-lint-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)
trägt den Eintragstyp an einer Stelle — Titel (Zeile 973), die Absätze
*Eintragstyp* (Zeile 977), *Default-Pin und Übersteuerung* (Zeile 979),
*Prüfsummen-Quellen* mit der zwei-Schlüssel-Tabelle (Zeilen 981 bis 986),
*Abweisungen* (Zeile 990), *Cache-Werte* (Zeile 992), *Meldungen*
(Zeile 994) und die Beleg-Tabelle (Zeilen 996 bis 1010); der Katalog-Abschnitt
[`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen)
listet die Dockerfile-Block-Keys Zeile für Zeile (Zeilen 1554 bis 1556 —
gelesen); dort bekommt der `java-toolchain`-Eintrag seine Zeile.

**Frage (Freigabe-Punkt F6):** die Stelle erweitern oder eine eigene
Technik-Stelle daneben schneiden. **Empfehlung: erweitern** — der Eintragstyp
steht an einer Stelle: die Absätze tragen den dritten Schlüssel, die
Quellen-Tabelle bekommt die Zeile der Stelle 3, die Beleg-Tabelle bekommt die
Java-Zeilen (je Bedingung der Beleg-Auftrag; die Test-Namen trägt der
Umsetzungs-Diff — [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)),
und der Titel nennt `java-toolchain`. In der Gegenoption (eine eigene Stelle
neben [`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-und-golangci-lint-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv))
stünden die Typ-Regeln (Abweisungen, Meldungen,
`version:`-Bedingung) an zwei Stellen; dieselbe Regel wäre doppelt
belegt und driftfähig — genau die Kopie-Form, die
[`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)
und [`AGENTS.md`](../../../../AGENTS.md) §3.10 vermeiden; nur die Java-Abweichungen (Vertrieb, JAVA_HOME)
wären neu — aber die Abweichungen sind Zeilen derselben Tabellen, keine
eigenen Absätze. **Contra** (der Erweiterung): der Titel verlässt seine
Zwei-Schlüssel-Fassung, und der Absatz *Cache-Werte* erweitert einen
ausdrücklich auf die zwei Go-Schlüssel begrenzten Zustand um die Abgrenzung
des dritten (der CR trägt die Form der Abgrenzung — Stelle 4). Die Wahl wird
beim Implementer als Teil der Freigabe-Vorlage geführt (CR-Slice §1,
DoD-Zeile 3); der *Neu*-Wortlaut der Spec-Stellen (nicht der Lastenheft-Stellen)
ist Teil derselben Freigabe wie der Lastenheft-Wortlaut.

### Stelle 8 — Meldung des Erzeugens (Warn-Form)

**Träger:** Sätze **J7**/**J9** (der
[`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md)
Entscheidung 3 letzter Satz: die Meldung weist den aktiven Block und, bei
einer Übersteuerung, deren Quelle aus; „Wortlaut und Ausweisung sind
CR-Wortlaut, wie bei den Vorlagen eine Warnung" — die zwei Schlüssel tragen
den freigegebenen Stand, Fundstelle Zeile 424, achter Satz).

**Form (Vorschlag):** je aktiviertem Schlüssel ein Meldungseintrag der
Ebene `warn` mit der Kennung [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)
als `code` (Regel [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe),
Muster `toolchainNotice` — Fundstelle `devcontainer_toolchain.go` Zeilen 322
bis 337, gelesen 2026-11-02: `Level: "warn"`, `Code` =
[`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), die
Message trägt die Quelle (als Wortlaut der Konstante), die Version (Default-Pin
per `ARG`-Namen oder Übersteuerung samt gelesener Quelle) und die Grenze („the
archive is not signed and u-boot checks no signature, so the check does not
protect against a release that was harmful from the start"); `Subject:
".devcontainer/devcontainer.json"`). Die Java-Meldung trägt daneben die
Vertriebs-Form (Eclipse-Temurin-Distribution der LTS-Version 21) — sie ersetzt das
Quellen-Formular nicht, sondern konkretisiert es. Die Meldung ist keine
Sicherheitslockerung und kein Degradations-Fallback (keine der vier
Ausnahmen von `--quiet`, [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)
Zeile 196: „Alle anderen Hinweise und Warnungen entfallen unter `--quiet`") —
der Typ ist seit 0.30.0 so vertraglich (keine fünfte Ausnahme); der CR
erfindet keine zusätzliche Ausweis-Pflicht.

### Stelle 9 — [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion): die Hosts-Zeile

**Alt (gelesen 2026-11-02):** die Allowlist-Tabelle (Fundstelle
[`spec/spezifikation.md`](../../../../spec/spezifikation.md) Zeilen 1459 bis
1476) trägt die Zeile `java` mit `repo.maven.apache.org`, `repo1.maven.org`
(Zeile 1466) und keine Zeile `java-toolchain`.

**Neu (Form-Vorschlag):** eine Tabellen-Zeile nach dem Muster der
Feature-Zeilen (Zeilen 1464 bis 1473), mit den drei Werten:

> Quelle: `java-toolchain` · Hosts: die Hosts der Veröffentlichung der
> Temurin-Binaries (Archiv-Download und die Prüfsummen-Datei derselben
> Veröffentlichung, ein Origin; je der gelesenen Lage, siehe F8) ·
> Herkunft: `feature:java-toolchain`

Die Herkunfts-Nennung folgt dem Muster der
Feature-Zeilen (Zeilen 1464 bis 1473). **Zeitpunkt (Freigabe-Punkt F8):**
die Host-Werte sind in dieser Rolle nicht gemessen — die Vertriebs-Messung
(Brew-Form, Archivform) ist Pflicht des Umsetzungs-Slices (Plan §2). Zwei
Aufstellungen: (a) die Zeile entsteht mit dem CR-Commit auf die **gelesene**
Vertriebs-Lage — der Implementer liest die Veröffentlichung am CR-Tag, nennt
Fundstelle und Datum in §7, und der Docker-Test des Umsetzungs-Slices pinnt
die Hosts; deckt die Umsetzungs-Messung die Zeile anders, ist sie eine
Plan-Berichtigung (Rückführung §4 (b) des CR-Slices deckt den Widerspruch) —
Empfehlung; (b) die Zeile entsteht mit dem Umsetzungs-Slice nach der Messung
— der DoD-Punkt
([`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion)-Zeile)
verschiebt sich dorthin, was eine
Plan-Berichtigung am CR-Slice (DoD-Zeile 3) voraussetzt. In keinem der beiden
Stände steht ein erfundener Host im Vertrag (Regel 2).

### Stelle 10 — [`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config): die zwei Schlüssel-Namen-Zeilen

**Alt (gelesen 2026-11-02):** die Werte-Zelle der `version`-Zeile (Fundstelle
[`spec/spezifikation.md`](../../../../spec/spezifikation.md), Zeile 522) und die
Schreibpipeline-Stufe 6 (Zeile 544) tragen je den Klammer-Zusatz mit den zwei
Namen: „(`go-toolchain`, `golangci-lint` — ihre Übersteuerung prüft
`generate devcontainer` gegen die Quellen-Angabe, [`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-und-golangci-lint-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv))"
(Zeile 522; Zeile 544 in der verkürzten Form „(`go-toolchain`, `golangci-lint`;
[`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-und-golangci-lint-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv))"
in derselben Regel-Weise).

**Neu (Vorschlag):** je Zeile der dritte Name in der Klammer, der Zusatz selbst
unverändert:

> ... oder einer der Schlüssel der gepinnten Werkzeugkette ist (`go-toolchain`,
> `golangci-lint`, `java-toolchain` — ihre Übersteuerung prüft
> `generate devcontainer` gegen die Quellen-Angabe,
> [`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-und-golangci-lint-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)) ...

> ... muss der Eintrag eine `source` tragen oder einer der Schlüssel der
> gepinnten Werkzeugkette sein (`go-toolchain`, `golangci-lint`,
> `java-toolchain`; [`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-und-golangci-lint-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)) ...

**Beleg:** die Regel steht schlüsselfrei im
Code (`featureVersionNeedsSourceForCat`, Fundstelle Zeilen 844 bis 852 —
der dritte Name verlangt keine zweite Regel-Stelle, nur die Konstanten der
Stelle 3); die Zwei-Zeilen-Fassung ist die Fundstellen-Aussage des CR-Slices
§3 („die zwei Schlüssel-Namen-Zeilen 522 und 544“ — gelesen am 2026-11-02).

### Stelle 11 — Handbuch-Katalog-Stellen und Namenslisten

**Alt (gelesen 2026-11-02):** die Nutzer-Doku trägt den Schlüssel nicht; die
Katalog-Tabelle (Fundstelle [`docs/user/benutzerhandbuch.md`](../../../../docs/user/benutzerhandbuch.md)
Zeilen 1129 bis 1149) kennt `java` (Zeile 1134) und die zwei
Go-Schlüssel-Zeilen (Zeilen 1147 und 1148); die Namenslisten der READMEs
tragen die zwei Go-Schlüssel mit „`gh` and `glab` (GitHub and GitLab CLI),
`go-toolchain` and `golangci-lint` (toolchains from pinned release archives,
both profiles)" (Fundstelle [`README.md`](../../../../README.md) Zeile 153 bzw.
[`README.de.md`](../../../../README.de.md) Zeile 156); das JVM-Beispiel
(Zeilen 1084 bis 1089, 1113) zeigt
`java.version 21` — die Berichtigung dieses Beispiels ist eine Mess-Frage des
Umsetzungs-Slices (Tags des `java`-Features, Plan §2 (c)), nicht dieses CR.

**Neu (Form-Vorschlag):** je Stelle — die Katalog-Tabellen-Zeile nach dem
Muster der zwei Go-Zeilen (Zeilen 1147 bis 1148: „Kein OCI-Feature: Das
erzeugte Dockerfile installiert ... beide Profile; die Version steht als
Build-Argument im Dockerfile und ist mit `version` übersteuerbar
(`u-boot` liest dafür die Prüfsummen-Datei derselben Veröffentlichung)"),
ohne Kennungen im Nutzertext ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6);
die Namenslisten der READMEs bekommen jeweils ein Paar `java-toolchain`
(Wortlaut der Paare je der Form der Freigabe). **Ist-Zustand-Frage:** die
Nutzer-Doku beschreibt den Ist-Zustand — die Schlüssel existieren dort noch
nicht; falls der Projektinhaber die Ist-Zustand-Regel auch hier zieht, gehen
die Diffs in den Umsetzungs-Slice (Rückführung §4 des CR-Slices; Muster der
Plan-Berichtigung in
[`slice-cr-devcontainer-gepinnte-werkzeugkette`](../done/slice-cr-devcontainer-gepinnte-werkzeugkette.md)
§7). Die Form der Stelle ist Teil der Freigabe; der CR-Disziplin-Abschnitt
(§1) entscheidet sie, nicht dieser Vorschlag.

### Stelle 12 — CHANGELOG-Form

**Alt (gelesen 2026-11-02):** der `## [Unreleased]`-Block (Fundstelle
[`CHANGELOG.md`](../../../../CHANGELOG.md) Zeile 12) ist leer; der
hinterliegende `0.19.0`-Block (Zeilen 14 bis 29) trägt die Form der
Werkzeugketten-Zusagen samt der Pin-Zahlen als Chronik (Zeile 23: „with this
release Go `1.27.3` and golangci-lint `v2.14.0`“).

**Neu (Vorschlag):** der Block bekommt einen `### Added`-Eintrag in derselben
Form: der Schlüssel `java-toolchain` als sechster Dockerfile-Block-Eintrag
neben `ollama`, `gh`, `glab` und den zwei Go-Schlüsseln, beide Profile; die
Vertriebs-Form (Eclipse-Temurin-Distribution, LTS 21, Release-Archiv aus der
Veröffentlichung der Temurin-Binaries mit geprüfter SHA-256 je Architektur vor
dem Entpacken); die `version:`-Übersteuerung mit Verifikation gegen die
Prüfsummen-Datei derselben Veröffentlichung und die drei Abweisungs-Fälle je
Exit-Code `10` (nicht lesbar/nicht überdeckend, ohne Quellen-Angabe,
`java` + `java-toolchain`); `JAVA_HOME` und `PATH` im erzeugten Block, keine
`containerEnv`-Werte des Schlüssels; die Warn-Meldung, die unter `--quiet`
entfällt; **die konkrete Temurin-Revision als Chronik** — sie steht im
Eintrag erst mit der Messung des Umsetzungs-Slices (Zeitpunkt: der
Release-Cut, Muster „with this release ...“ der Zeile 23); `Promised`/
`Not promised` je der Grenz-Sätze (**J10**/**J11**). Kennungen tragen die
Link-Klammern (wie Zeile 23), nicht der Nutzertext.

## 3. Risiken des CR-Slices (§6) — Ausgang je Stellen-Bezug

| Risiko (CR-Slice §6) | Ausgang in diesem Vorschlag |
|---|---|
| 1 — Wortlaut geht nicht wortgleich ins Lastenheft | bearbeitet: die Stellen §§2 sind Copy-Vorlagen mit Anker-Links innerhalb der Zieldatei; Zeitpunkt der Historie-Zeile folgt dem DoD-Zeile-2-Abgleich (Stellen 1, 2, 10). |
| 2 — Zusage weiter als ADR, Messung oder Freigabe | bearbeitet: Stelle 3 (Was der Docker-Test belegt) und die Grenz-Sätze **J10**/**J11**; die Brew-Form trägt keinen CR-Wortlaut (Freigabe-Kontext, Mess-Punkt F5). |
| 3 — Bedingung der Technik-Stelle ohne pinnenden Test | bearbeitet: die Beleg-Tabelle trägt der Umsetzungs-Diff je Zeile (Stelle 7); keine Belege im CR. |
| 4 — Handbuch beschreibt den Ist-Zustand | benannt: Stelle 11 — Form je Freigabe; falls die Ist-Zustand-Regel auch hier zieht, Diffs in den Umsetzungs-Slice (Rückführung CR-Slice §4). |
| 5 — Plan-Aussagen am Fundort nicht gelesen | bearbeitet: die Fundstellen dieses Vorschlags sind am 2026-11-02 gelesen und jeder Satz je Stelle trägt seine; der Abgleich läuft nach jeder Fix-Runde erneut ([`AGENTS.md`](../../../../AGENTS.md) §3.10). |

## 4. Punkte für die Freigabe

Die vier Entscheidungen und die Brew-Auflage vom 2026-11-02 sind **getragen**;
neu zu entscheiden ist der *Neu*-Wortlaut je Stelle (§§2) und der Punkt F6 bis
F8:

| Punkt | Frage | Empfehlung | Stand (2026-11-02) |
|---|---|---|---|
| **F1** | Schlüssel-Name und Verhältnis zum Katalog-Schlüssel `java` (Abweisung oder Warnung) | `java-toolchain`, Abweisung mit Exit-Code `10` | **entschieden** — `java-toolchain`, Abweisung mit Exit-Code `10`; externes `java` allein wählbar (Satz **J8**) |
| **F2** | JDK-Vertrieb und Archivform | Eclipse Temurin, LTS 21, Archive aus der Veröffentlichung der Temurin-Binaries | **entschieden** — Temurin, LTS 21; die konkrete Revision liest die Messung des Umsetzungs-Slices |
| **F3** | Default-Pin-Zahlen im Vertragstext | nein — Chronik im CHANGELOG und in der Lastenheft-Historie; im CR Muster F4 des Vorlaufs | **entschieden** — nein (Stelle 6) |
| **F4** | Prüfsummen-Quelle | Same-Origin je Release; `api.adoptium.net` aus | **entschieden** — Prüfsummen-Datei derselben Veröffentlichung, je gewählter Version (Stelle 3) |
| **F5** | Die Brew-Auflage (2026-11-02) — Vertriebs-Befund über Homebrew (Cask `temurin`/`temurin@21`, Formel-Verhältnis zum Feature `homebrew`) | Mess-Punkt des Umsetzungs-Slices; keine zweite Installationsform im CR | **Auflage** — Mess-Pflicht in [`slice-v1-katalog-feature-jdk-gepinnt`](../open/slice-v1-katalog-feature-jdk-gepinnt.md) §2 (d); der CR trägt sie nicht |
| **F6** | Technik-Stelle: [`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-und-golangci-lint-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) erweitern oder eigene Stelle daneben | erweitern — der Eintragstyp bleibt an einer Stelle (Stelle 7) | **offen** — entscheidet mit der Wortlaut-Freigabe |
| **F7** | Fall `sdkman` + `java-toolchain`: derselbe Abweisungsfall wie externes `java`? | derselbe Abweisungsfall, sobald die gelesene Feature-Lage ihn trägt; sonst Absatz beim externen `java` allein und Frage ins Beobachtungs-Register (Stelle 5) | **offen** — entscheidet mit der Wortlaut-Freigabe |
| **F8** | Zeitpunkt der Hosts-Zeile in [`SPEC-019`](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) | im CR-Commit auf die gelesene Vertriebs-Lage, gepinnt durch den Docker-Test des Umsetzungs-Slices (Stelle 9) | **offen** — entscheidet mit der Wortlaut-Freigabe |