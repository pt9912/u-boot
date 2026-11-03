# Review-Report: slice-cr-devcontainer-containerenv — 2026-11-03

**Review-Art:** Contract-Review (CR-Diff gegen den freigegebenen Wortlaut, Plan und
Hard Rules). Nicht geprüft: DoD-Closure und Verification-Evidence — Verifier.

**Gegenstand:** [`slice-cr-devcontainer-containerenv`](../plan/planning/done/slice-cr-devcontainer-containerenv.md),
Range `f36d5f59..d62aad5e` — Schnitt `89b584e2` (Plan + [`cr-vorschlag-containerenv.md`](../plan/cr/cr-vorschlag-containerenv.md)),
reine Moves `046ea7ba`/`367046c6` (0-Diff, R100) mit Verweis-Nachzug je eigenem
Commit `2c7f4521`/`0be9cc54` (je 1 eingehend), Roadmap-Start `2830bb4c`,
Vertrags-Übernahme `d62aad5e` (`spec/lastenheft.md` 0.32.0, `spec/spezifikation.md`
[SPEC-040](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config),
Vorschlags-Statuszeile). Freigegebener Wortlaut: Vorschlags-Datei, Status
„freigegeben (2026-11-03, … „Freigabe wie vorgelegt“ …)“ — die Freigabe selbst
ist im Agenten-Gespräch gefallen und außerhalb der Repo-Evidenz.

**Skill:** `.harness/skills/reviewer.md` @ `v0.20.0`
**Modell:** glm-5.3-flash · **Datum:** 2026-11-03 · **Rolle:** Reviewer (Modul 10)

**Eingangs-Kontext** (gelesen, gegen die geprüft wurden):

- Slice-Plan §1–§8 samt Risiken und §3-Dateitabelle; Umsetzungs-Plan
  [`slice-v1-katalog-projekt-containerenv`](../plan/planning/done/slice-v1-katalog-projekt-containerenv.md)
  als Kontext (liegt in `open/`)
- Freigegebener Vorschlag Stellen (a)–(d) samt vier gebilligten Entscheidungen
- Lastenheft: Absatz *Projekt-Umgebungswerte* (Zeile 387, innerhalb
  [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen)
  nach dem Cache-Absatz), Historie-Zeile 0.32.0 (Zeile 1634), Versionsfeld
- [`SPEC-040`](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)
  am Ist-Stand (Zeilen 515–561) und am Alt-Stand `f36d5f59` (Übersteuerungs-Zeilen
  522/544 des JDK-Stands); [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)
  als Form-Muster
- Fundstelle der im Vorschlag genannten Validation: `driven.ValidateDevcontainerEnvName`
  lebt (`internal/hexagon/port/driven/devcontainer_launcher.go:72`)
- Vorgänger-Bestand-Zeilen der CR-Slices (`open/README.md` §Bestand, Zeilen 92–93)

---

## Findings

### F-1 — MEDIUM · abgelaufene Zustandsmarker bei der Übernahme nicht nachgezogen

- **Kategorie:** MEDIUM
- **Quelle:** [`AGENTS.md`](../../AGENTS.md) §3.10 (Zustandsmarker mit Ablaufpunkt wird beim Ablauf auf den Ist-Stand gesetzt) und §3.7
- **Fundstellen:**
  - [`docs/plan/planning/in-progress/roadmap.md`](../plan/planning/in-progress/roadmap.md) · „die Freigabe steht aus.“
  - [`docs/plan/planning/done/cr-vorschlag-containerenv.md`](../plan/cr/cr-vorschlag-containerenv.md) · „die Vertrags-Übernahme folgt im CR-Slice.“
- **Klasse:** `zustandsmarker-abgelaufen-bei-uebernahme`
- **Befund:** `d62aad5e` setzt die Statuszeile des Vorschlags auf *freigegeben*
  und übernimmt den Vertrags-Wortlaut — die Freigabe ist also vor dem Commit
  erteilt und die Übernahme im selben Commit erfolgt. Die Roadmap-Zeile (in
  `2830bb4c` geschrieben, in `d62aad5e` nicht berührt) sagt am Range-Tip weiter
  „die Freigabe steht aus“, und die Statuszeile des Vorschlags sagt „die
  Vertrags-Übernahme folgt im CR-Slice“ — beides ist zum Commit-Zeitpunkt
  falsch.
- **Risiko:** Der nächste Leser (Verifier, Closure-Lauf) liest den Zustand und
  muss die Autorisierung der Vertrags-Zeilen aus der Commit-Message
  rekonstruieren; nimmt er die Zeile ernst, hält er die Übernahme für vorzeitig.
- **Verifizierbar:** Review-only (kein Gate prüft Zustandswahrheit).

### F-2 — MEDIUM · Stelle (c) weicht vom freigegebenen Wortlaut ab, Abweichung nicht ausgewiesen

- **Kategorie:** MEDIUM
- **Quelle:** Freigabe „wie vorgelegt“ (Wortlaut ist Norm) + [`AGENTS.md`](../../AGENTS.md) §3.10 (Abgleich nach jeder Änderung; eintragsgenauer Ausweis)
- **Fundstellen:**
  - [`cr-vorschlag-containerenv.md`](../plan/cr/cr-vorschlag-containerenv.md) Stelle (c) · „beide Zeilen tragen den Namen `devcontainer.containerEnv` neben die bestehenden“
  - [`spec/spezifikation.md`](../../spec/spezifikation.md) Wertetabelle · „| `devcontainer.containerEnv` | ja | ja | Map von Namen zu String-Werten;“
  - Commit `d62aad5e` · „Vertrags-Übernahme wortgleich“ / „[`SPEC-040`](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config) die Pfad-Zeile devcontainer.containerEnv“
- **Klasse:** `abweichung-freigabe-wortlaut-nicht-ausgewiesen`
- **Befund:** Der freigegebene Neu-Wortlaut von Stelle (c) bindet an die zwei
  Übersteuerungs-Zeilen des Alt-Satzes (JDK-Stand 522 und 544: die
  Wertetabellen-Zeile `devcontainer.features.<feature>.version` und
  Schreibpipeline-Schritt 6 — beide vorhanden, beide unberührt) und sagt, *beide*
  Zeilen tragen den Namen `devcontainer.containerEnv`. Übernommen ist stattdessen
  **eine neue, eigene Pfad-Zeile** in der Wertetabelle (zwischen
  `devcontainer.enabled` und `devcontainer.featureSources.allow`); die zwei
  Übersteuerungs-Zeilen tragen den Namen nicht. Die Commit-Message benennt die
  Stelle bereits in der neuen Lesart („die Pfad-Zeile“, Singular) und behauptet
  zugleich „wortgleich“; ein Ausweis dieser Abweichung je Stelle existiert im
  Repo nicht (§7 des Slices ist noch leer). Für die Freigabe-Form ist die
  konkrete Übertragung (eigene Zeile statt Nennung in den zwei Übersteuerungs-Zeilen)
  die sinnvolle Lesart — sie ist aber nicht die freigegebene Form.
- **Risiko:** Der §7-Abgleich des DoD („wortgleich … Abweichungen je ausgewiesen“)
  scheitert an Stelle (c) oder wird still auf die neue Lesart abgenommen; dann
  steht ein Vertragsstand ohne belegte Freigabe-Basis für genau diese Form, auf
  dem der Umsetzungs-Slice baut.
- **Verifizierbar:** Review-only — Auflösung je nach Auslegung: Rückfrage an den
  Projektinhaber (Übertragungs-Lesart als Freigabe bestätigen) oder Ausweis je
  Stelle in §7 des Slices.

### F-3 — LOW · Stelle (b) Historie-Zeile nicht wortgleich; Ausweis nur in der Commit-Message

- **Kategorie:** LOW
- **Quelle:** [`AGENTS.md`](../../AGENTS.md) §3.4 (Decken-Regel) + Freigabe „wie vorgelegt“
- **Fundstellen:**
  - [`cr-vorschlag-containerenv.md`](../plan/cr/cr-vorschlag-containerenv.md) Stelle (b) · „`| 0.32.0 | Projekt-Umgebungswerte: `devcontainer.containerEnv` als Map des erzeugten DevContainers (Umsetzungs-Slice nach diesem CR) |`"
  - [`spec/lastenheft.md`](../../spec/lastenheft.md) Historie · „| 0.32.0 | 2026-11-03 | [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen): neuer Absatz der Projekt-Umgebungswerte — …“
- **Klasse:** `abweichung-freigabe-wortlaut-nicht-ausgewiesen` (zweites Vorkommen — Drift-Signal)
- **Befund:** Die übernommene Historie-Zeile trägt vier Zellen (Version, Datum,
  Zusammenfassung, Freigabe-Aussage) statt der Vorschlags-Dreizellen-Form und
  fasst den Absatz zusammen statt den Vorschlagstext zu wiederholen; der
  Down-Verweis „(Umsetzungs-Slice nach diesem CR)“ fehlt zurecht (Decken-Regel).
  Die Abweichung ist in der Commit-Message benannt („Decken-Regel der
  Historie-Zeile berichtigt“), aber die Commit-Message ist Chronik, nicht der
  je-Stelle-Ausweis des DoD.
- **Risiko:** Der §7-Abgleich findet Abweichungen ohne Repo-internen Ausweis;
  als wiederkehrendes Muster (F-2 und F-3) driftet der Wortgleichheits-Anspruch
  der CR-Übernahme vom belegten Zustand.
- **Verifizierbar:** Review-only — Ausweis in §7 bei der Closure (Risiko 1 des
  Plans sieht genau das vor).

### F-4 — LOW · Plan-§3-Zeile „open/README.md (Bestand-Zeile)“ ohne Diff und ohne Begründung

- **Kategorie:** LOW
- **Quelle:** [`AGENTS.md`](../../AGENTS.md) §3.10 (jeder Eintrag ohne Diff wird berichtigt)
- **Fundstellen:**
  - [`slice-cr-devcontainer-containerenv.md`](../plan/planning/done/slice-cr-devcontainer-containerenv.md) §3 · „`roadmap.md` (die Zeile beim Start), `open/README.md` (Bestand-Zeile)“
  - [`docs/plan/planning/open/README.md`](../plan/planning/open/README.md) §Bestand — keine Zeile für diesen Slice; `git diff --name-status f36d5f59..d62aad5e` berührt die Datei nicht
- **Klasse:** `plan-zeile-ohne-diff-und-ohne-begründung`
- **Befund:** Die Plan-§3-Tabelle führt die Bestand-Zeile in `open/README.md`
  als Änderung; im Range wurde sie nie geschrieben (auch nicht beim Schnitt
  `89b584e2`). Die Vorgänger-CR-Slices tragen in derselben Tabelle Zeilen
  ([`slice-cr-devcontainer-gepinnte-werkzeugkette`](../plan/planning/done/slice-cr-devcontainer-gepinnte-werkzeugkette.md), Zeile 92;
  [`slice-cr-devcontainer-java-toolchain`](../plan/planning/done/slice-cr-devcontainer-java-toolchain.md), Zeile 93) — der Bestand-Muster-Satz
  ist also real.
- **Risiko:** Der Bestand des Planning-Baums drifft vom etablierten Muster ab,
  und der §3.10-Abgleich bei der Closure stolpert über einen Eintrag ohne
  Diff — unklar, ob die Zeile noch nachgeholt oder der Plan-Eintrag gestrichen
  werden soll.
- **Verifizierbar:** Review-only — Abgleich Plan §3 gegen `git diff --stat` bei
  der Closure (AGENTS.md §3.10).

### F-5 — LOW · Schreibform des Basis-Pfades `devcontainer.containerEnv` im Vertrag nicht spezifiziert

- **Kategorie:** LOW
- **Quelle:** Maintainability (Vertragskongruenz der Wertetabelle-Form)
- **Fundstellen:**
  - [`spec/spezifikation.md`](../../spec/spezifikation.md) Wertetabelle · „| `devcontainer.containerEnv` | ja | ja | Map von Namen zu String-Werten;“ (Schreiben `ja` ohne Wildcard)
  - derselbe Unterabschnitt · „`u-boot config set devcontainer.containerEnv.<name> <wert>` validiert den Namen“
- **Klasse:** `pfad-form-luecke-wertetabelle`
- **Befund:** Die Wertetabelle setzt den Basis-Pfad als schreibbar (`ja`), der
  Unterabschnitt definiert aber nur den Schreibzugriff über den untergeordneten
  Pfad `<name>`. Andere Familienpfade der Tabelle tragen die Wildcard ausdrücklich
  (`devcontainer.features.<feature>.enabled`, `services.<svc>.enabled`); ob der
  Basis-Pfad selbst schreibbar ist (ganze Map in einem Aufruf?) und welche
  Pipeline-Regeln dann gelten (No-Op-Schritt 3, Listenpfad-Semantik, Schritt 6),
  sagt der Vertrag nicht.
- **Risiko:** Der Umsetzungs-Slice muss die Semantik des Basis-Pfad-Schreibens
  erfinden oder den Wert „ja“ klein auslegen — beides ist eine
  Contract-Decision, die der CR eigentlich trägt; ohne Klarung landet ein
  implementiertes Verhalten im Code ohne Vertragstext (Muster `BEO-DOK/spec-bedingung-ohne-code-entsprechung`).
- **Verifizierbar:** Review-only — Klärung im CR (Kurzsatz im Unterabschnitt)
  oder Auslegungs-Entscheidung des Umsetzungs-Slices mit Ausweis.

### F-6 — INFO · Tippfehler in der freigegebenen Vorschlags-Datei

- **Kategorie:** INFO
- **Quelle:** Maintainability
- **Fundstelle:** [`cr-vorschlag-containerenv.md`](../plan/cr/cr-vorschlag-containerenv.md) · „die vii Entscheidungen des Umsetzungs-Plans“
- **Klasse:** `zahlangabe-typo`
- **Befund:** „vii“ statt „vier“ — als römische Zahl gelesen zählt der Satz
  sieben Entscheidungen; gebilligt sind vier. Reiner Schreibfehler ohne
  Vertragswirkung (die vier Stellen sind separat benannt).
- **Risiko:** —
- **Verifizierbar:** Review-only.

### F-7 — INFO · Move der Vorschlags-Datei bei der Closure nicht im Plan geführt

- **Kategorie:** INFO
- **Quelle:** Maintainability (Vorgänger-Muster)
- **Fundstellen:**
  - Plan §3 führt nur die Neuanlage der Vorschlags-Datei · „`cr-vorschlag-containerenv.md` (neu)“
  - Vorgänger: [`cr-vorschlag-java-toolchain.md`](../plan/cr/cr-vorschlag-java-toolchain.md) und [`cr-vorschlag-gepinnte-werkzeugkette.md`](../plan/cr/cr-vorschlag-gepinnte-werkzeugkette.md) liegen in `done/`; hier liegt die übernommene Datei in [`open/`](../plan/cr/cr-vorschlag-containerenv.md) mit Verweis `../in-progress/`
- **Klasse:** `move-nachzug-vorschlagsdatei`
- **Befund:** Die freigegebene und übernommene Vorschlags-Datei liegt im
  Lifecycle-Verzeichnis `open/`, während ihr Slice in `in-progress/` liegt; beim
  Slice-Move nach `done/` muss sie nach Vorgänger-Muster mitwandern (reiner
  Move, Verweis-Nachzug — `make slice-mv` <!-- d-check:ignore (Make-Target, keine Slice-Kennung) --> zieht den eingehenden Verweis nach).
- **Risiko:** Keines unmittelbar; Kontext für die Closure.
- **Verifizierbar:** Review-only.

---

## Geprüft ohne Befund

- **Wortgleichheit Stelle (a):** Zeichen-level-Abgleich (programmatisch,
  `difflib.SequenceMatcher`) des Lastenheft-Absatzes gegen das freigegebene
  Blockquote — **identisch**; einziger Unterschied ist das Absatz-Label
  `**Projekt-Umgebungswerte.**`, das im Ist-Stand als Absatz-Anfang in derselben
  Form steht. Negativbefund.
- **Decken-Regel** ([`AGENTS.md`](../../AGENTS.md) §3.4): Neue Stellen
  Lastenheft-Absatz, Lastenheft-Historie 0.32.0, [`SPEC-040`](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)-Unterabschnitt
  samt Beleg-Tabelle und Spec-Historie 2026-11-03 — **kein** Verweis abwärts
  (Slice, Carveout, Roadmap, Vorschlag, ADR); die Historie-Zeile trägt nur
  `LH-*`-Links. Die Formulierungen „Pflicht der Umsetzung“ / „Umsetzungs-Diff“
  folgen dem [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)-Präzedenz
  (bereits in der 0.31.0-Zeile verankert). Der im Commit-Message zugesagte
  Decken-Regel-Berichtiger (Down-Verweis aus dem Vorschlags-Wortlaut entfernen)
  ist bestätigt — keine Reste.
- **Lastenheft-Nutzertext:** Keine `SPEC-*`-/ADR-/Slice-Kennungen im neuen
  Absatz; Abweisungen mit Exit-Code `10` und Grenz-Sätze (idempotent, ohne
  Duplikate, ohne den Schlüssel erlaubt) in der Satzform der Nachbartexte.
- **Struktur der Contract-Straten:** Absatz unter
  [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen)
  (Überschrift Zeile 377, nach dem Cache-Absatz, vor *Mindestdateien*);
  Spezifikation: Pfad-Zeile in Wertetabellen-Form (Spalten Lesen/Schreiben/Werte),
  Unterabschnitt nach der Schreibpipeline am Ende von [`SPEC-040`](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config),
  vor der Überschrift [`SPEC-041`](../../spec/spezifikation.md#spec-041--schweregrade-der-prüfung-composeyamlvalid);
  Beleg-Tabelle in Auftrags-Form — **keine erfundenen Testnamen** (Muster
  [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)).
- **Historie-Form:** Beide Spec-Straten tragen Datum und
  Freigabe-/Entscheidungs-Spalte gemäß Nachbartexten; Versionen konsistent
  (Versionsfeld 0.32.0 = Historie-Zeile 0.32.0 = Spec-Verweis „zum Lastenheft
  0.32.0“).
- **Plan-Aussagen über Code** (§3.10): Die im Vorschlag und Plan genannte
  Validation `driven.ValidateDevcontainerEnvName` existiert
  (`internal/hexagon/port/driven/devcontainer_launcher.go:72`); die Fundstellen
  der zwei Übersteuerungs-Zeilen („522 und 544 des JDK-Stands“) stimmen am
  Alt-Stand `f36d5f59`.
- **Exit-Code-Vertrag** ([`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)):
  Neue Abweisungen je Exit-Code `10`, nichts geschrieben — einheitlich über
  beide Stellen (Validation, Kollision) und beide Befehle (`config set`,
  `generate devcontainer`).
- **Lifecycle-Moves** ([`AGENTS.md`](../../AGENTS.md) §3.3): Beide Moves reine
  `git mv` (0-Diff), Verweis-Nachzug je eigener Commit (§2c7f4521/§0be9cc54,
  je 1 eingehender Verweis); der Vorschlags-Verweis zeigt am Range-Tip korrekt
  auf `in-progress/`. Keine Referenz zeigt auf die alten Orte.
- **§3.9 Kommentar-Selbstcheck:** Keine Code-, Config- oder Skript-Dateien im
  Range; die neuen Spec-Stellen tragen keine Kommentare der verbotenen Formen
  (kein „früher“, „nicht mehr“, „originally“; keine Verweise auf geänderte
  Symbole).
- **Kein Produktcode, kein Nutzertext:** `git diff --name-status` über den
  Range — genau fünf Dateien, alle aus Vertrags- und Planning-Straten; README,
  `docs/user/`, Handbuch und `CHANGELOG.md` sind nicht berührt (Muster der
  CR-Slices; Träger folgen der Implementierung im Umsetzungs-Slice).
- **Kein Gate berührt:** `Makefile`, `.golangci.yml`, `.d-check.yml`, ADR und
  `.github/workflows/` fehlen im Range.
- **Wiedergefundene Klassen:** Bericht-Dateien-Familie nicht fällig —
  `docs/reviews/` ist im Range nicht berührt; die zwei Unterabschnitts- und
  Pfad-Zeilen-Muster von [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)
  werden sauber fortgeführt.
- **Anker-Ziele:** [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) existieren;
  alle Spezifikations-Inter-Links ([`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv), [`SPEC-023`](../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen), [`SPEC-004`](../../spec/spezifikation.md#spec-004--projektnamen-muster), …) und
  Planning-Links lösen — bestätigt durch `make docs-check` (0 Befunde).

## Nicht geprüft

- **DoD-Closure und Verification-Evidence** — Verifier (Abgrenzung Reviewer).
- **`make gates` / `make lint` / `make test`** — kein Code- oder Config-Diff im
  Range; Vertrags-Slice ohne produktberührte Sub-Area, der Lauf trägt keinen
  Befund zur geprüften Frage (Wortlaut-Kongruenz). Begründung wie im
  Vorgänger-Review des java-toolchain-CR.
- **Der Freigabe-Akt selbst** — im Agenten-Gespräch gefallen, außerhalb der
  Repo-Evidenz; die Statuszeile und die Freigabe-Spalten der Historie-Zeilen
  sind als Beleg angenommen (Auftragskontext).

## Sensors

| Sensor | Ergebnis |
|---|---|
| `make docs-check` | **grün** — 1060 Dateien geprüft, 0 Befunde |
| `make ablage-pruefung` | **grün** — Wächter-Test 26 Fälle, 0 Fehlschläge; Positivliste des Planning-Baums (`.md`-Neuzugänge konform) |
| `make verify-depguard` | nicht gelaufen — keine Import-/Pfad-Änderung im Range (kein Code-Diff); Grund: gegenstandslos zu den geprüften Linsen |
| `make test-docker` | nicht gelaufen — kein Docker-Pfad und kein Produktcode im Range (Begründung wie oben) |

## Verdikt

**Mit Findings fortführen, vor der Closure klären.** Kein Blocking-Finding an
der übernommenen Wortlaute selbst — Stelle (a) ist wortgleich, die Decken-Regel
ist an allen neuen Stellen erfüllt, der CR-Diff bleibt sauber in den
Vertrags-Straten. Die zwei MEDIUM-Findings betreffen die Belegbarkeit der
Übernahme, nicht deren Substanz:

- **F-1** (abgelaufene Zustandsmarker in Roadmap und Vorschlags-Statuszeile) ist
  mit einem kleinen Berichtigungs-Commit zu lösen, der keine Vertragszeile
  ändert.
- **F-2** braucht eine Entscheidung: bestätigt der Projektinhaber die
  Übertragungs-Lesart von Stelle (c) (eigene Pfad-Zeile statt Nennung in den
  zwei Übersteuerungs-Zeilen) — oder wird die Abweichung je Stelle in §7
  ausgewiesen. **Der Verifier darf die Closure nicht schließen, solange F-2
  ohne Ausweis oder Bestätigung steht** (DoD: „wortgleich … Abweichungen je
  ausgewiesen“).
- **F-3, F-4, F-5** sind Closure- bzw. Umsetzungs-Anker (§7-Ausweis,
  §3.10-Abgleich, Klärung vor dem Umsetzungs-Diff) und blockieren den CR nicht
  rückwirkend.