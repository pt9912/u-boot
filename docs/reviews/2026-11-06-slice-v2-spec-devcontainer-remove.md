# Review-Report: slice-v2-spec-devcontainer-remove — 2026-11-06

**Review-Art:** Plan-/Design-Review der Spezifikation: Diff gegen Slice-Plan, Lastenheft, Messprotokoll und Hard Rules (kein Code im Diff).

**Gegenstand:** Commit `f64294a0` (HEAD, Basis `638996ca`), vier Dateien, 179 Einfügungen, 27 Löschungen: [`spec/spezifikation.md`](../../spec/spezifikation.md) (neue [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes), Nachzug in [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes), `--quiet`-Liste, Feld `command`, Verweisliste, Historie), [`spec/architecture.md`](../../spec/architecture.md), [`docs/maintainer/cli-json-contract.md`](../maintainer/cli-json-contract.md) und der Slice [`slice-v2-spec-devcontainer-remove`](../plan/planning/done/slice-v2-spec-devcontainer-remove.md) (Plan-Tabelle, erste DoD-Zeile).

**Skill:** `.harness/skills/reviewer.md`. **Modell:** claude-sonnet-5-5. **Datum:** 2026-11-06.

## Eingangs-Kontext und Sensors

**Gelesen:**

- Slice [`slice-v2-spec-devcontainer-remove`](../plan/planning/done/slice-v2-spec-devcontainer-remove.md) (Kopf, §1 bis §6, §8) und der Folge-Slice [`slice-v1-devcontainer-remove`](../plan/planning/open/slice-v1-devcontainer-remove.md) (DoD, Plan-Tabelle, Nicht-Ziele), Welle [`welle-devcontainer-remove`](../plan/planning/welle-devcontainer-remove.md) nur über die Verweise.
- Lastenheft: [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (alle Spiegelstriche), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes).
- Spezifikation: [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (ganz), [`SPEC-038`](../../spec/spezifikation.md#spec-038--bestätigungspfad-von-u-boot-down---volumes), `--quiet`-Abschnitt, Kopf und Historie.
- Messprotokoll [`slice-v1-devcontainer-remove-messung`](../messungen/slice-v1-devcontainer-remove-messung/README.md) (Kurzfassung, M1 bis M5, Grenzen), gegen jede `Messung M…`-Fundstelle von [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes) gelesen.
- Code als Beleg der Aussagen über Bestand: `internal/hexagon/application/downservice.go` (`Down`, `runConfirmationGate`: Freigabe vor dem Engine-Zugriff), `internal/adapter/driven/confirm/confirm.go` (`ConfirmRemoveVolumes`: `y`/`yes`, Ende der Eingabe als `N`, Lesefehler als Fehler), `internal/hexagon/port/driving/removeservice.go` (`ErrConfirmerUnavailable`), `Dockerfile` (Stage `runtime`: `gcr.io/distroless/static-debian12:nonroot`, kein `docker`).
- [`ADR-0015`](../plan/adr/0015-sandbox-volumes-pro-instanz.md) (Bezug des Slice, im Diff nicht berührt), `AGENTS.md` §3.4, §3.7, §3.10, §5, `harness/review.md`.

**Sensors:**

- `make docs-check`: gefahren auf HEAD `f64294a0` vor Anlage dieses Reports, `1175 Datei(en) geprüft, 0 Befund(e)`; nach Anlage dieses Reports erneut gefahren: `1176 Datei(en) geprüft, 0 Befund(e)` (der erste Lauf mit dem Report fand 22 unverlinkte Kennungen und ein totes Linkziel im Report selbst; berichtigt).
- Eigene Proben, nur lesend: `git grep -lF` auf den Anker von [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (93 Dateien, stimmt mit der Plan-Aussage überein); Suche in `spec/` und `docs/maintainer/` nach weiteren Aufzählungen von `up`/`down`/`shell` ohne `remove` (Treffer nur in `docs/maintainer/devcontainer-arbeiten.md` und `docs/maintainer/zusagen-und-messstaende.md`, beide vom Folge-Slice getragen); `git diff --stat` gegen die Plan-Tabelle (vier Dateien, vier Plan-Zeilen mit „update“, zwei „gelesen, unverändert“ ohne Diff).
- `make gates`, `make test-docker`, `make verify-depguard`: nicht gefahren. Der Diff enthält nur Doku, keinen Code, keinen Docker-Pfad und keine Importregel; ein `make gates`-Nachweis ist Pflicht der Closure und Sache des Verifiers.

## Findings

### MEDIUM `spec/spezifikation.md` · „Im Standardprofil ohne Katalog-Feature hält der Container kein Volume (Messung M1, `std` und `anon`)“ — Messbeleg trägt die Aussage nicht, `anon` hält ein Volume

**Klasse:** Aussage nennt eine Messung, die sie nicht trägt (Spec-Aussage über Verhalten ohne gelesenen Beleg)
**Quelle:** `AGENTS.md` §3.10 (eine Aussage über Verhalten nennt nur, was gemessen oder gelesen wurde; Bedingung in der Spec bekommt einen Test oder wird gestrichen); [`slice-v1-devcontainer-remove-messung`](../messungen/slice-v1-devcontainer-remove-messung/README.md) M1
**Befund:** M1 führt `anon` als Standardprofil mit einem von Hand angehängten `VOLUME` und **einem anonymen Volume** in `Mounts`; nur `std` hat 0 Volumes. [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes) belegt „kein Volume, die Menge ist dann leer“ mit `std` **und** `anon`, und führt im selben Abschnitt *Menge der Volumes* die anonymen Volumes (64 Hexzeichen, „gemessen in M1 und M2e“) als Teil der Menge.
**Risiko:** Der Implementer des Folge-Slice baut den Docker-Test „leere Menge im Standardprofil“ (Beleg-Tabelle, Zeile zur Menge der Volumes) nach dem Wortlaut auch mit der Variante `anon`: Sie liefert ein Volume, die Aussage und der Test widersprechen der Regel „alle Einhängungen vom Typ `volume`“ in derselben Stelle, und kein Gate zeigt es vor dem Lauf des Docker-Tests.
**Verifizierbar:** Review-only (Lesen von M1 gegen die Stelle; `make docs-check` prüft keine Messaussagen)

### MEDIUM `spec/spezifikation.md` · „Die Umsetzung des Befehls belegt jede Aussage dieser Stelle; eine Aussage ohne Beleg entfällt.“ — Beleg-Tabelle deckt vier Aussagen der Stelle nicht ab

**Klasse:** Bedingung in der Spec ohne pinnenden Beleg (Beleg-Tabelle unvollständig gegen den Text)
**Quelle:** `AGENTS.md` §3.10 (letzter Spiegelstrich: Bedingung in Spec oder Handbuch bekommt einen Test, der sie pinnt, oder wird gestrichen); DoD-Zeile 1 des Slice; die Regel der Stelle selbst („eine Bedingung, die dort keinen Beleg bekommt, entfällt“)
**Befund:** Keine Zeile der Tabelle *Belege* nennt: (a) „Der Befehl liest keine Projektdatei (weder `u-boot.yaml` noch `.devcontainer/devcontainer.json`), braucht keinen Hilfs-Container und bindet keinen Container-Runtime-Socket ein“; (b) „Der Exit-Code `14` kommt bei `remove` nicht vor“; (c) den Absatz *Nach dem Entfernen* („Ein folgendes `u-boot devcontainer up` legt einen neuen Container an“, Wiedereinhängen und Neuanlage der Volumes); (d) den Inhalt der Fehlermeldung bei Exit-Code `12` („nennt den gescheiterten Schritt, die Meldung der Engine, die Volumes, die entfernt sind, und die, die bestehen“) und der Meldung zu Exit-Code `10` ohne Freigabe. Zeile 8 der Tabelle belegt nur die Exit-Codes `2` und `10` des Aufrufs, Zeile 5 nur Zustand und Exit-Code `12`.
**Risiko:** Nach der eigenen Regel der Stelle entfallen diese Aussagen, oder der Code-Slice pinnt sie nicht, weil seine Arbeitsliste die Tabelle ist: Ein späterer Änderer lässt `remove` `u-boot.yaml` lesen oder den Socket einbinden, die Spec bleibt wahr im Text und falsch im Verhalten, und kein Test oder Gate schlägt an.
**Verifizierbar:** Review-only (Tabelle *Belege* gegen die Absätze *Aufruf*, *Pfade und Exit-Codes*, *Nach dem Entfernen*, *Fehler beim Entfernen*)

### MEDIUM `docs/plan/planning/open/slice-v1-devcontainer-remove.md` · „dieser Slice setzt sie um und ändert sie nicht“ — Übergabe Spec-Slice → Code-Slice hat keinen Träger für die Zielbild-Pflichten

**Klasse:** Übergabe ohne benennbares Artefakt beim Empfänger (Pflicht im Absender-Slice, nicht im Empfänger-Slice)
**Quelle:** `AGENTS.md` §3.10 (Zustandsmarker mit Ablaufpunkt wird beim Ablauf entfernt oder auf den Ist-Stand gesetzt); `harness/roles.md` (Übergabe-Artefakte); Slice [`slice-v2-spec-devcontainer-remove`](../plan/planning/done/slice-v2-spec-devcontainer-remove.md) §6
**Befund:** §6 des Spec-Slice setzt als Ausgang „offen; der Code-Slice pinnt jede Aussage mit einem Test“, und [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes) trägt den Marker „Die Stelle beschreibt das Zielbild des Befehls“. Der Folge-Slice nennt [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes) nirgends (sein Kopf sagt noch „und die Stelle, die der Spec-Slice anlegt“), seine DoD verweist nicht auf die Tabelle *Belege*, seine Plan-Tabelle führt `spec/spezifikation.md` nicht, und sein Nicht-Ziel sagt „dieser Slice setzt sie um und ändert sie nicht“; `architecture.md` führt er als „keine“, obwohl der Spec-Slice die Port-Zeilen dort ausdrücklich unverändert lässt, „weil die Namen der neuen Port-Operationen der Code-Slice festlegt; er trägt sie dort nach“.
**Risiko:** Der Implementer des Code-Slice folgt seiner DoD und seiner Plan-Tabelle: Er pinnt nur, was dort steht, entfernt den Marker „Zielbild“ nicht (er darf die Spec nach seinem Nicht-Ziel nicht ändern) und lässt die Port-Zeilen der Architektur stehen. Die Spec behauptet danach dauerhaft ein Zielbild für Verhalten, das der Code hat, und `architecture.md` nennt `DevcontainerLauncher (Up/Stop/Shell)` neben `entfernen`; kein Gate sieht beides.
**Verifizierbar:** Review-only (Lesen der DoD und der Plan-Tabelle des Folge-Slice gegen §6 des Spec-Slice und gegen [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes))

### LOW `spec/spezifikation.md` · „Ohne `--volumes` haben `--yes` und `--no-interactive` keine Wirkung“ — Satz widerspricht dem Satz davor

**Klasse:** Zwei Aussagen derselben Stelle schließen sich aus (Spec-intern)
**Quelle:** `AGENTS.md` §3.10 (Bedingung in der Spec muss pinnbar sein); [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) (`--yes` und `--no-interactive` sind exklusiv)
**Befund:** Derselbe Absatz *Aufruf* sagt, dass `--yes` zusammen mit `--no-interactive` „auch ohne `--volumes`“ mit Exit-Code `2` endet; der letzte Satz sagt, ohne `--volumes` hätten beide Optionen keine Wirkung. Die Tabelle der Freigabe führt in den ersten beiden Zeilen beide Spalten mit „beliebig“, was die Kombination einschließt.
**Risiko:** Ein Test, der den Satz „keine Wirkung“ pinnt (`remove --yes --no-interactive` ohne `--volumes` endet mit `0`), und ein Test für Exit-Code `2` derselben Kombination sind nicht beide zu schreiben; der Implementer entscheidet still.
**Verifizierbar:** Review-only

### LOW `spec/spezifikation.md` · „mit `--volumes` entfernt sie `docker rm -v` (Messung M2e)“ — gemessen ist `docker rm -f -v`, der Ablauf nutzt `docker rm -v` nach `docker stop`

**Klasse:** Messung auf andere Aufrufform übertragen (Aussage über Verhalten ohne gemessene Form)
**Quelle:** `AGENTS.md` §3.10 (eine Aussage über Verhalten nennt nur, was gemessen wurde); Messprotokoll M2e, M4
**Befund:** M2e und M4 fahren `docker rm -f -v`; M2c fährt `docker rm` ohne `-f` auf einem gestoppten Container, ohne `-v`. Schritt 8 von [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes) ruft `docker rm -v` ohne `-f` auf einem zuvor gestoppten Container und nennt dafür „Messung M2e“ (anonyme Volumes) beziehungsweise „Messung M2a“ (Schritt 7/8).
**Risiko:** Die Reichweite-Liste der Stelle führt die nicht gemessenen Fälle einzeln auf, diese Aufrufform steht nicht darunter; der Docker-Test des Folge-Slice ist die einzige Stelle, die sie belegt, und der Leser hält sie für gemessen.
**Verifizierbar:** Review-only (Messprotokoll gegen Schritt 8)

### LOW `docs/maintainer/cli-json-contract.md` · „liefert den Minimalkontrakt mit `subcommand: "remove"`“ — Zielbild in der Gegenwartsform, ohne den Marker der Spec

**Klasse:** Zielbild-Aussage ohne Marker in einem Dokument, das gelesenen Bestand beschreibt
**Quelle:** `AGENTS.md` §3.10 (Aussage über Verhalten nennt nur, was gelesen wurde; Zustandsmarker); [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes) selbst (trägt „Zielbild“)
**Befund:** Der neue Absatz in §6.11 und die neue Zeile der Matrix stehen in der Gegenwartsform („liefert“, „Es gibt keine Abfrage“) in einem Dokument, dessen übrige Absätze Symbole und Tests des vorhandenen Codes nennen (`devcontainerUpData`, `mapDevcontainerDownErrorToDiagnostic`, `TestDevcontainerUp_AlreadyRunning`); `remove` hat im Code keinen Befehl. Die Spec trägt den Marker, dieses Dokument nicht.
**Risiko:** Ein Maintainer, der die Matrix oder den Abschnitt liest, hält `remove --json` für ausgeliefert; die Tabelle der Sentinels führt (zu Recht) nur `up` und `down`, und kein Gate vergleicht die Matrix mit dem Cobra-Baum.
**Verifizierbar:** Review-only

### INFO `spec/spezifikation.md` · „stützen sich auf das Messprotokoll“ — Slice-Kennung als Linktext in der Spec

**Klasse:** Slice-Kennung als Linktext im Technik-Stratum
**Quelle:** Kopf von `spec/spezifikation.md` („verweist nur aufwärts … nie auf … Slices“); `AGENTS.md` §3.4 (Decken-Regel)
**Befund:** Das Ziel ist ein Messprotokoll unter `docs/messungen/` (Sub-Area DOK), kein Slice; der Linktext ist die Slice-Kennung. `make docs-check` ist grün; Präzedenz: [`SPEC-057`](../../spec/spezifikation.md#spec-057--homebrew-formeln-im-devcontainer-katalog-schlüssel-brew-pfade-freigabe-der-taps-konflikt-mit-homebrew-pin-meldungen-hosts-und-grenzen) verlinkt ein Messprotokoll mit dem Text „Messprotokoll“.
**Risiko:** Keine Änderungspflicht; ein Leser kann den Linktext als Verweis auf den Slice lesen.
**Verifizierbar:** `make docs-check` (grün)

### INFO `spec/spezifikation.md` · „Der Beleg ist der Aufruf ohne `docker` im `PATH`, nicht das Image selbst.“ — Image-Satz des Lastenhefts hat einen Ersatzbeleg

**Klasse:** Zusage-Menge-Satz mit Ersatzbeleg (Zuständigkeit Verifier)
**Quelle:** [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Zusage-Menge, Image-Satz); DoD-Punkt (1) des Folge-Slice
**Befund:** Der Image-Satz („`u-boot devcontainer remove` läuft dort ebenfalls nicht“) wird in [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes) auf den Aufruf ohne `docker` im `PATH` abgebildet; die Aussage „enthält kein `docker`“ ist am `Dockerfile` (Stage `runtime`, distroless) gelesen und stimmt. Ob der Ersatzbeleg die DoD-Alternative „pinnender Test oder Nicht-Zusage in Spec und Handbuch“ erfüllt, entscheidet der Verifier des Code-Slice.
**Risiko:** Keine Änderungspflicht an diesem Diff.
**Verifizierbar:** Review-only (Verifier)

### INFO `docs/maintainer/cli-json-contract.md` · „Ein Fehler-Umschlag trägt keine `data`“ — Teilzustand nach Exit-Code `12` nur im Meldungstext

**Klasse:** Fehler-Umschlag ohne `data` bei teilweise ausgeführter destruktiver Operation
**Quelle:** [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe); [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes) (Absatz *Fehler beim Entfernen*, Absatz *Ausgabe*)
**Befund:** Nach einem Fehler in Schritt 9 ist der Container entfernt und ein Teil der Volumes ebenfalls; der Umschlag trägt weder `removed` noch `removedVolumeNames`, die Fehlermeldung nennt die entfernten und die bestehenden Volumes als Text. Die Wahl entspricht `up` (dort fehlt `data` im Fehlerfall ebenfalls); für `remove` ist sie die einzige Ausgabe über den Zustand einer Löschung.
**Risiko:** Kein Vertragsbruch; ein Skript liest den Zustand aus dem Meldungstext. Zur Kenntnis für den Architect und den Validator.
**Verifizierbar:** Review-only

### INFO `spec/architecture.md` · „Einzelheiten [`SPEC-049`]“ — Zeiger der Beschreibung von `DevcontainerStartService` nennt [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes) nicht; Zeilenangaben der Plan-Tabelle gelten für den Stand vor dem Diff

**Klasse:** Rang-Zeiger deckt nicht alle genannten Operationen; Plan-Zeilenangaben gegen den Ist-Stand
**Quelle:** `AGENTS.md` §3.10 (Abgleich der Tabelle für den Stand nach der letzten Fix-Runde)
**Befund:** Die Beschreibung nennt jetzt „entfernen“ und verweist für Einzelheiten weiter nur auf [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes), dessen Absatz unter der Überschrift auf [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes) zeigt (ein Sprung mehr). Die Plan-Tabelle des Slice führt Zeilennummern (`Zeilen 299 bis 314`, `Zeile 619`, `Zeile 656`), die nach dem Diff verschoben sind; nur die erste Zeile trägt den Stand `7f5110cc`.
**Risiko:** Keine Änderungspflicht; der Abgleich der Tabelle gegen `git diff --stat` bei der Closure gilt für den Stand nach der letzten Fix-Runde.
**Verifizierbar:** Review-only

## Negativbefunde

**Geprüft ohne Befund:**

- Lastenheft gegen [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes): Aufruf, Suche (Instanz, ohne Instanz neben Instanzen, nur Instanzen mit Exit-Code `10`, kein Container mit Exit-Code `0`), Reihenfolge Abfrage vor Stoppen und Entfernen, Freigabe-Tabelle (sechs Fälle), Exit-Codes `10`/`11`/`12`/`2`/`0`, `--json` ohne Abfrage, Nennung der Volumes vor der Abfrage und nach dem Entfernen, Nicht-Ziele (keine Aufräum-Liste, Volumes eines nicht mehr vorhandenen Containers) — kein Widerspruch zum Wortlaut von [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) und [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung). Die Spannung zwischen Exit-Code `2` (Zeile 212) und `10` (Zeile 214) im Lastenheft benennt die Stelle und der CR-Vorschlag führt sie als offen (nicht doppelt gemeldet).
- Messungen: `M2a` bis `M2d`, `M3` (Projekt `std`, Projektordner unberührt), `M4` (Podman-Entfernen nicht gemessen), `M5` (Fälle B bis E), Umgebung (Docker Engine 29.9.3, `overlayfs`, höchstens drei Container) stimmen mit dem Protokoll überein.
- Aussagen über Bestand: Freigabe vor dem Engine-Zugriff bei `down --volumes` (`runConfirmationGate` vor `projectVolumeSnapshot` und `ComposeDown`), `y`/`yes` ohne Beachtung der Schreibung und Ende der Eingabe als Ablehnung (`ConfirmRemoveVolumes`), Image ohne `docker` (`Dockerfile`, Stage `runtime`), 93 verlinkende Dateien des Ankers von [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes).
- Dual-Classifier-Regel: kein neuer Sentinel im Diff; die Beleg-Tabelle verlangt Pins in der zentralen Klassifikation **und** in der Abbildung des Befehls, die Tabelle der Fehlerklassen in `architecture.md` ordnet Stoppen und Entfernen der Klasse `12` zu, `architecture.md` nennt vier Diagnostic-Abbildungen.
- Referenzmodell und Hard Rules: `architecture.md` ohne ADR-, Slice- oder Meilenstein-Bezug (§3.4); [`ADR-0015`](../plan/adr/0015-sandbox-volumes-pro-instanz.md) und alle `Accepted`-ADRs unberührt (§3.5); Planning-Artefakt liegt nur in `in-progress/` (§3.3, kein Move im Commit); keine Chronik oder Konjunktiv-Kommentare in den neuen Absätzen (§3.7, §3.9 (a) bis (c)); Sprachvertrag: die Texte der Abfrage und der Zeilen in [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes) sind Englisch ([`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache)).
- Plan-Tabelle gegen `git diff --stat`: vier geänderte Dateien, jede mit Plan-Zeile; die zwei „gelesen, unverändert“-Zeilen (`open/README.md`, `roadmap.md`) tragen keinen Diff.
- Kopien der Befehlsliste: `--quiet`-Liste, `--quiet`-Absatz, Feld `command`, Verweisliste zu [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (Aufruf, Instanzen, Name, Suche, Instanzen-Absätze, Vorab-Prüfung, Reichweite), `architecture.md` (CLI-Adapter, `ARC-008`, Fehlerklassen), Matrix und §6.11 des JSON-Vertrags nennen `remove`; die Sentinel-Tabelle des JSON-Vertrags führt bewusst `up` und `down`.
- Commit-Message: nennt Slice, [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`ADR-0015`](../plan/adr/0015-sandbox-volumes-pro-instanz.md); keine Struktur-IDs ([`AGENTS.md`](../../AGENTS.md) §5 Regel 1).

**Nicht geprüft:**

- Code, Tests und Docker-Pfade — nicht im Diff; die Umsetzung ist Gegenstand von [`slice-v1-devcontainer-remove`](../plan/planning/open/slice-v1-devcontainer-remove.md).
- Nutzerdokumentation (`README.md`, `README.de.md`, `docs/user/`, `CHANGELOG.md`) — vom Folge-Slice getragen; im Diff nicht berührt (§5 Regel 5 betrifft dort den Zeitpunkt der Lieferung des Befehls).
- `make gates`, `make test-docker`, `make verify-depguard` — Begründung oben (nur Doku; DoD-Nachweis ist Sache des Verifiers).
- DoD-Abgleich der Closure-Pflichten (grüne Gates, Verification-Report, Closure-Notiz) — Verifier-Aufgabe, kein Gegenstand des Reviews.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 3 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Aussage nennt eine Messung, die sie nicht trägt · Bedingung in der Spec ohne pinnenden Beleg · Übergabe ohne benennbares Artefakt beim Empfänger · Zwei Aussagen derselben Stelle schließen sich aus · Messung auf andere Aufrufform übertragen · Zielbild-Aussage ohne Marker · Slice-Kennung als Linktext im Technik-Stratum · Zusage-Menge-Satz mit Ersatzbeleg · Fehler-Umschlag ohne `data` bei teilweise ausgeführter destruktiver Operation · Rang-Zeiger deckt nicht alle genannten Operationen

## Verdikt

**Merge-blockierend:** ja — für die drei MEDIUM-Findings (Berichtigung der Messangabe, Ergänzung oder Streichung der vier Aussagen ohne Beleg, Träger der Übergabe im Folge-Slice oder eine begründete Ablehnung); die LOW- und INFO-Findings blockieren nicht. Kein HIGH: Der Diff ändert keinen Vertrag der Ränge 1 bis 3 gegen die Decken-Regel, keinen Code und keine Gate-Schwelle.

**Übergabe:** Findings gehen an den Implementer des Slice [`slice-v2-spec-devcontainer-remove`](../plan/planning/done/slice-v2-spec-devcontainer-remove.md); die Finding-Klassen gehen in die Closure-Notiz (§7) und von dort in den Zähler. Der Report ersetzt weder `make gates` noch den Verification-Report; die DoD-Konformität prüft der Verifier.
