# Review-Report: slice-v1-katalog-projekt-containerenv — 2026-11-03

**Review-Art:** Umsetzungs-Review (Code-Diff gegen Plan, Vertrags-Carrier und
Hard Rules). Nicht geprüft: DoD-Closure und Verification-Evidence — Verifier
([`harness/verification.md`](../../harness/verification.md)).

**Gegenstand:** [`slice-v1-katalog-projekt-containerenv`](../plan/planning/done/slice-v1-katalog-projekt-containerenv.md),
Commit `59ad23b4` auf `main` (einziger Umsetzungs-Commit; CR-Slice und
Vertrags-Übernahme liegen im Vorgänger-Vorgang
[`slice-cr-devcontainer-containerenv`](../plan/planning/done/slice-cr-devcontainer-containerenv.md)).
22 Dateien, +1558/−46.

**Vertrags-Carrier (gelesen, gegen die geprüft):**
[`spec/lastenheft.md`](../../spec/lastenheft.md) 0.32.0, Absatz
*Projekt-Umgebungswerte* (Zeile 387, innerhalb
[`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen)),
Historie-Zeile 0.32.0 (Zeile 1634);
[`spec/spezifikation.md`](../../spec/spezifikation.md)
[`SPEC-040`](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)
Unterabschnitt *Projekt-Umgebungswerte* samt Beleg-Tabelle (Zeilen 550–563) und
Historie-Zeilen 2026-11-03 (1627–1628).

**Skill:** `.harness/skills/reviewer.md` @ `v0.20.0`
**Modell:** glm-5.3-flash · **Datum:** 2026-11-03 · **Rolle:** Reviewer (Modul 10)

**Eingangs-Kontext** (gelesen, gegen die geprüft):

- Slice-Plan §1–§8 samt CR-Entscheidungen (§4), Risiken (§6) und §3-Tabelle am
  Berichtigungs-Stand des Umsetzungs-Diffs
- CR-Reviews der Vorgänger-Vorgänge:
  [`2026-11-03-slice-cr-devcontainer-containerenv.md`](2026-11-03-slice-cr-devcontainer-containerenv.md)
  (F-1…F-7, darunter F-5 zur Wertetabellen-Form des Basis-Pfads) und der
  Verify-Berichtigungs-Commit `45833353` (V-1–V-4)
- Code-Kontext: [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)-Werkzeugketten-Stand
  (`internal/hexagon/application/devcontainer_toolchain.go` —
  `toolchainState`, `resolveToolchain`),
  [`driven.ValidateDevcontainerEnvName`](../../internal/hexagon/port/driven/devcontainer_launcher.go)
  (`[A-Za-z_][A-Za-z0-9_]*`),
  [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)-Namens-Weg
  (`internal/hexagon/application/devcontainer_start.go:579`),
  Schema-Reader (`internal/hexagon/application/initproject.go`,
  `doctor.go:746`), YAML-Codec
  (`internal/adapter/driven/yaml/codec.go` — `scalarNode` schreibt Strings
  mit Tag `!!str`, Zahlen-Schreibweisen werden gequoted)

---

## Findings

### F-1 — HIGH · Kommentar von `listCandidatePaths` behauptet eine Verhaltens-Grenze, die der Code nicht trägt

- **Kategorie:** HIGH
- **Quelle:** [`AGENTS.md`](../../AGENTS.md) §3.7 (Ein Kommentar beschreibt,
  was da ist — Zusage/Kopplung/Abgrenzung/Rang-Zeiger/Grenze) ·
  Baseline-Regelwerk `grundlagen-harness-dateien.md` §Was ein Kommentar trägt
- **Fundstellen:**
  - [`internal/hexagon/application/config_many.go`](../../internal/hexagon/application/config_many.go) ·
    „whitelisted path and stays out of the list.“
  - Gegenstelle desselben Diffs: [`internal/hexagon/application/devcontainer_containerenv.go`](../../internal/hexagon/application/devcontainer_containerenv.go) ·
    „a hand edit can name anything“
- **Klasse:** `kommentar-falsche-verhaltensaussage`
- **Befund:** Der Kommentar sagt, ein durch Hand-Edit gesetzter Name
  außerhalb der env-name-Regel werde vom Domain-Konstruktor abgewiesen und
  bleibe deshalb aus der `config list`-Ausgabe heraus. Genau das Gegenteil
  gilt — und zwar absichtlich: der Konstruktor splitet bewusst **unvalidiert**
  ([`domain/configpath.go`](../../internal/hexagon/domain/configpath.go) ·
  „deliberately without a domain validator“, Test
  `TestNewConfigPath_ContainerEnv_EntryParsesTheName` pinnt „`A.B` … want
  the raw remainder as name“). Nur der leere Name wird verworfen. Ein
  hand-editierter Eintrag `BAD-NAME` bekommt also einen whitelisted path und
  erscheint in `config list`; die Aussage des Kommentars stimmt nur für den
  Namen `""`. Der Kommentar widerspricht damit auch dem zweiten Kommentar
  desselben Diffs („a hand edit can name anything“ — dort als Grund für das
  JSON-Escaping), der das Ist-Verhalten korrekt beschreibt.
- **Risiko:** Der nächste Änderer liest die Grenze als Sicherheit: Er kann auf
  der Prämisse „ungültige Namen kommen nie durch den Listen-/Render-Pfad“
  eine Escaping- oder Validierungs-Maßnahme abziehen oder list-seitige
  Schritte als Namens-Filter behandeln — kein Gate fängt den Kommentar
  (`make lint` prüft keine Aussageinhalte; die Tests pinnen das Gegenteil,
  lassen aber den Kommentar als Erklärung stehen).
- **Verifizierbar:** Review-only — Berichtigung per Kommentar-Diff;
  die Berichtigung ist ein Commit am Plan-Anker
  [`LH-FA-CONF-005`](../../spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern)
  (Konfig-Anzeige) mit Testbezug, kein Träger-Änderungstext.

### F-2 — MEDIUM · [`SPEC-040`](../../spec/spezifikation.md) Wertetabelle sagt „Lesen: ja“ für den Basis-Pfad — umgesetzt und gepinnt ist Exit `10` für dieselbe Lesung

- **Kategorie:** MEDIUM
- **Quelle:** [`SPEC-040`](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)
  (Wertetabelle gegen Beleg-Tabelle und Umsetzung) · MEDIUM-Anker
  [`harness/review.md`](../../harness/review.md) §Finding Categories
  (Doku/Spec driftet gegen implementiertes Verhalten)
- **Fundstellen:**
  - [`spec/spezifikation.md`](../../spec/spezifikation.md) Wertetabelle ·
    „| `devcontainer.containerEnv` | ja | ja (nur Unterpfade“
  - [`spec/spezifikation.md`](../../spec/spezifikation.md) Beleg-Tabelle ·
    „Basis-Pfad je Lesen und Schreiben mit Exit `10` am Nutzerpfad“
  - [`cmd/uboot/containerenv_test.go`](../../cmd/uboot/containerenv_test.go) ·
    `TestRun_ConfigSet_ContainerEnvValueAndNames` — `config get
    devcontainer.containerEnv` → Exit `10`, „not a known config path“
  - [`CHANGELOG.md`](../../CHANGELOG.md) · „ends with exit code `10` as an
    unknown path, also for `config get`“
- **Klasse:** `spec-drift-wertetabelle-lesen-basispfad`
- **Befund:** Die Umsetzung legt den Basis-Pfad als **keine Pfad-Form**
  aus (weder lesbar noch schreibbar; unbekannter Pfad, Exit `10`) und trägt
  diese Auslegung in Beleg-Tabelle, CHANGELOG und Domain-Test. Die Spalte
  Lesen der Wertetabellen-Zeile steht dagegen unverändert auf `ja` —
  Nachbarmuster der Tabelle ist die Lesung der Listenpfade
  (`devcontainer.featureSources.allow` liest die ganze Liste). Die
  CR-Review-Runde hatte genau diese Lücke als
  [F-5 `pfad-form-luecke-wertetabelle`](2026-11-03-slice-cr-devcontainer-containerenv.md)
  benannt („Klärung vor dem Umsetzungs-Diff … oder Auslegungs-Entscheidung
  des Umsetzungs-Slices mit Ausweis“); die Auslegung ist ausgewiesen, aber
  die Wertetabellen-Zelle wurde nicht auf den Ist-Stand gebracht — derselbe
  SPEC-Abschnitt trägt nun beide Aussagen.
- **Risiko:** Ein Nutzer (oder ein nächster Agentenlauf) liest die
  Wertetabelle, erwartet `config get devcontainer.containerEnv` als Map-Lesung
  (wie bei `featureSources.allow` die Listen-Lesung) und erhält Exit `10`;
  der Vertragsstand ist für die Les-Frage widersprüchlich, und kein Gate —
  `make docs-check` prüft Anker, nicht Semantik — kann das auflösen. Die
  Zusage-Ebene bleibt sauber ([`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen)
  verspricht die Basis-Pfad-Lesung nicht), nur die Technik-Zeile drifft.
- **Verifizierbar:** Review-only — Korrektur der Zelle (Lesen `nein`,
  Schreiben `ja (nur Unterpfade …)`) oder eine bewusste Les-Freigabe des
  Projektinhabers mit nachziehender Zeile; Anker der Berichtigung ist der
  Unterabschnitt *Projekt-Umgebungswerte* von
  [`SPEC-040`](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config).

### F-3 — LOW · Plan-§3-Zeile nennt „READMEs“, der Diff berührt kein README — zweitvorkommen der CR-Klasse

- **Kategorie:** LOW
- **Quelle:** [`AGENTS.md`](../../AGENTS.md) §3.10 (jeder Eintrag ohne Diff
  wird berichtigt) — zweites Vorkommen der Klasse
  [`plan-zeile-ohne-diff-und-ohne-begründung`](2026-11-03-slice-cr-devcontainer-containerenv.md)
  (CR-Review F-4; das CR-Vorkommen wurde in `45833353`/`df6c40ce`
  berichtigt). Steering-Loop-Zähler der Klasse: 2/3.
- **Fundstellen:**
  - [`docs/plan/planning/done/slice-v1-katalog-projekt-containerenv.md`](../plan/planning/done/slice-v1-katalog-projekt-containerenv.md)
    §3 · „READMEs, `CHANGELOG.md`“
  - `git show 59ad23b4 --name-status` — kein `README*.md`-Berührungs-Eintrag;
    [README.md](../../README.md)/[README.de.md](../../README.de.md) führen
    nur Beispiele (`config set project.name …`), keinen Schlüsselkatalog,
    sind also inhaltlich nachtragfrei korrekt berührt-los
- **Klasse:** `plan-zeile-ohne-diff-und-ohne-begründung`
- **Befund:** Die Doku-Zeile der §3-Tabelle zählt READMEs zur
  Änderungsfamilie; im Umsetzungs-Diff wird kein README berührt. Die Zeile
  wurde im selben Commit für `repo.mk` und `.devcontainer` berichtigt, die
  READMEs-Unternennung blieb stehen. Der DoD-Satz „READMEs falls betroffen“
  (§2) trägt die sachliche Seite; die §3-Zeile ist der unberichtet gebliebene
  Abgleich-Rest.
- **Risiko:** Der §3.10-Abgleich bei der Closure („jeder Eintrag ohne Diff
  wird berichtigt“, Stand nach der letzten Fix-Runde) läuft gegen einen
  Eintrag ohne Beleg; als wiederkehrendes Muster (2. Vorkommen in drei
  Vorgängen) nähert sich die Klasse der 3×-Schwelle.
- **Verifizierbar:** Review-only — Streichung der Unternennung in §3 (oder
  Ausweis „falls betroffen: nicht betroffen“) bei der Fix-Runde.

### F-4 — INFO · Die Ausweisung zeigt den JSON-escaped Wert, nicht den Datei-Wert

- **Kategorie:** INFO
- **Quelle:** Maintainability (Ausweisungs-Treue)
- **Fundstelle:** [`internal/hexagon/application/devcontainer_containerenv.go`](../../internal/hexagon/application/devcontainer_containerenv.go) ·
  `+"=" + e.Value`
- **Klasse:** `ausweisung-escaped-wert`
- **Befund:** `containerEnvWarnings` iteriert `collectContainerEnvData` und
  zeigt damit die **escaped** Form (Quote, Backslash, `\n`), nicht den
  YAML-Wert; für alle Werte ohne Steuerzeichen sind beide Formen identisch.
  [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen)
  verlangt nur, dass jede erzeugte Zeile genannt wird — die Escaping-Form
  benennt dieselbe Zeile.
- **Risiko:** Keine Vertragslücke; bei einem Wert mit Steuerzeichen unterscheidet
  sich die Meldung optisch vom Dateiinhalt.
- **Verifizierbar:** Review-only.

### F-5 — INFO · Hand-defekte Namen: `generate` rendert unverändert (escaped), die Abweisung sitzt an der [`SPEC-049`](../../spec/spezifikation.md)-Namens-Stelle

- **Kategorie:** INFO
- **Quelle:** [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen)
  (Enforcement-Satz nennt nur `config set`) ·
  [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../plan/planning/observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)
- **Fundstellen:**
  - [`spec/lastenheft.md`](../../spec/lastenheft.md) 0.32.0 · „Der Name jedes
    Wertes ist ein gültiger Dev-Container-Umgebungsname … gegen beide setzt
    `u-boot config set` je fachlichen Fehler ab“
  - [`internal/hexagon/application/devcontainer_start.go:584`](../../internal/hexagon/application/devcontainer_start.go)
  · „driven.ValidateDevcontainerEnvName(name) != nil“ (Namen der
    erzeugten `containerEnv`, [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)-Hand-over-Stelle)
  - Pins des Diffs: `TestContainerEnv_LoadPathParsesTheMap`
    (Schema-Lesung nimmt hand-editierte Map), `TestGenerateDevcontainer_ContainerEnv_EscapesJSONContent`
    (Escape rendert, JSON bleibt gültig)
- **Klasse:** `grenze-ausnahme-ohne-ausweis`
- **Befund:** Ein hand-editierter Name außerhalb der Namensregel wird durch
  `generate` **nicht** abgewiesen — er rendert escaped in die erzeugte Datei
  (JSON bleibt gültig) und scheitert erst bei `devcontainer up` an der
  bestehenden Namens-Norm der Hilfs-Container-Übergabe. Der Vertrags-Satz
  bindet die Namens-Abweisung wörtlich nur an `config set`; die Grenze
  „`generate` weist nicht ab, er rendert“ trägt kein Nutzertext-Satz und
  kein Contract-Regel-Wortlaut.
- **Risiko:** Das Grenz-Satz-Fehlen ist nach [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../plan/planning/observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)
  kein Vertragsbruch (die Enforce-Stelle ist benannt, der Test trägt die
  Grenze), aber die Kette „`config list`/`generate` akzeptieren, `up`
 Exit `10`“ ist eine Grenze ohne Wortlaut — ein künftiger CR zum
  Namensraum muss die Kette kennen.
- **Verifizierbar:** Review-only — Grenz-Satz in
  [`SPEC-040`](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)
  Unterabschnitt (oder bewusste Nicht-Zusage in der „nicht zugesagt“-Zeile)
  bei nächster Vertrags-Berührung.

---

## Geprüft ohne Befund

- **Vertragskongruenz je Satz (Lastenheft 0.32.0, alle Sätze):** Map-Typ
  und Lesung → `ubootYAMLDevcontainer.ContainerEnv`
  (`internal/hexagon/application/initproject.go`) · Render ins verwaltete
  `containerEnv`-Objekt → Template
  `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl`
  (innerhalb des `BEGIN U-BOOT MANAGED BLOCK: init`) · „neben den
  Generator-Werten … ohne Duplikate“ → Merge hinter den Cache-Zeilen +
  `assertCacheValuesOnce` · idempotent (Regeneration schreibt neu,
  entfernter Wert verschwindet) → `TestGenerateDevcontainer_ContainerEnv_RemovedEntryDisappears`
  · Namens-Norm + Abweisung je Exit `10` je `config set` →
  `validateContainerEnvName` via [`driven.ValidateDevcontainerEnvName`](../../internal/hexagon/port/driven/devcontainer_launcher.go),
  `TestConfigSet_ContainerEnv_InvalidNameExits10`,
  `TestRun_ConfigSet_ContainerEnvValueAndNames` · Wert ist String (auch
  leer) → `TestConfigSet_ContainerEnv_HappyPathAndNoOp` (leerer Wert schreibt
  als `"FOO: """`-Eintrag, Get liest verbatim) · Kollision je Exit `10`,
  nichts geschrieben, ohne Schlüssel erlaubt → `toolchainRenderState` +
  `TestGenerateDevcontainer_ContainerEnv_CollisionExits10` (je
  Schlüsselkombination inkl. `--sandbox`, `java-toolchain` und
  `java`-Feature-Gegenprobe) · „nur in der `containerEnv`-Ebene, Dockerfile
  und handgeschriebener Rest tragen sie nicht“ →
  `TestGenerateDevcontainer_ContainerEnv_InBothProfiles` (Dockerfile-Namen)
  · Ausweisung je erzeugter Zeile → `TestGenerateDevcontainer_ContainerEnv_Notices`
  + `TestRun_Generate_ContainerEnvNotices`.
- **Beleg-Tabelle ↔ Baum:** alle 16 genannten Test-Namen existieren (ggrep
  über `internal/` und `cmd/`), inklusive der Subtest-Namen (`…/golden
  default profile`, golangci-only-`--sandbox`-Fall je genau ein Cache-Eintrag,
  java-toolchain ohne Cache-Raum); die [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)-Präzedenz-Frage
  „Testnamen gehören dem Umsetzungs-Diff“ ist erfüllt — keine Zeile bleibt
  als Auftrag stehen.
- **Ebenen-/Import-Regeln:** Domain importiert nichts Neues — die
  Namens-Prüffunktion bewusst **nicht** in der Domain dupliziert, sondern in
  der Anwendung via `port/driven` (`validateContainerEnvName`);
  `depguard`-Kontrakt bleibt intakt (`verify-depguard` 9/9, unten), keines
  der neun Muster bricht. Der Kommentar-Bezug
  `[driven.ValidateDevcontainerEnvName]` in `domain/configpath.go` ist
  Doc-Verweis ohne Import — zulässig.
- **Dual-Classifier (repo-spezifisch #1):** Kein neuer Sentinel; beide
  Abweisungen nutzen den bestehenden
  [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)-Klassifizierungs-Weg
  (`ErrConfigValueInvalid` → Exit `10`), der je Subkommando bereits
  abgebildet ist; Nutzerpfad-Pins: `config set` Exit `10`,
  `generate` Exit `10` mit `--dry-run --json`-Code-Kennung — die
  zentrale Klassifikation und die Abilden der zwei berührten Subkommandos
  tragen dieselbe Einträge-Form wie vor dem Diff.
- **Template-Arithmetik (hand-nachgezeichnet):** Komma-Platzierung liegt **vor**
  je Projekt-Eintrag (`{{if $goKeys}},{{else if $i}},{{end}}` →
  Cache-Fall: Komma vor erstem und jedem weitere Eintrag; ohne Cache: Komma
  erst ab dem zweiten) — Byte-Golden-Mischfall `mixedCacheJSONGolden` und
  `assertCacheValuesOnce` bestätigen; ohne Map und ohne Go-Schlüssel
  entsteht kein `containerEnv`-Objekt, mit Cache-Werten und leerer Map bleibt
  `toolchainJSONGolden` byte-identisch
  (`TestGenerateDevcontainer_ContainerEnv_WithoutTheMapByteIdentical`); der
  `{{- end}}`-Kommentar der Java-Grenze bleibt, erweitert um die
  Projekt-Werte-Zeile —
  [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)-Java-Grenze
  wahr („java-toolchain takes no cache space“).
- **Sortierung und JSON-Escape:** Sortiert je Schlüssel unterhalb der
  Cache-Zeilen (Positions-Abgleich im Test); Escaping trägt Quote,
  Backslash und `< 0x20` inklusive `\n/\r/\t`; der Zahlen-Fall
  (`config set … 123`) landet dank `scalarNode`-Tag `!!str` gequoted —
  kein Schema-Bruch gegen „jeder Wert ein String“.
- **Exit-Codes und Nicht-Schreiben:** Alle drei Abweisungsklassen laufen
  **vor** dem Schreiben (Validation vor Patch, Kollision vor Plan,
  SetMany-Items vor dem einen Write); Nothing-Written-Pins:
  `assertNothingWritten` je Kollisions-Fall, byte-identischer
  `u-boot.yaml`-Vergleich je Name-Fall,
  `.devcontainer`-Existenz-Negativprobe am Nutzerpfad. Exit-Codes nach
  [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) und
  [§ SPEC-043](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet):
  Ausweisung als gewöhnliche Warnung (kein `SecurityRelaxation`/`DegradationFallback`,
  gepinnt je Test), unter `--quiet` entfallend.
- **File Safety:** PatchScalar bleibt kommentarbewahrend (Happy-Path-Test
  prüft `# project comment`), keine User-Datei ohne Managed-Block getoucht;
  `.devcontainer/devcontainer.json`-Regeneration nur im verwalteten Objekt.
- **Docker-only / Gates:** Kein Makefile-, `.golangci.yml`- oder
  `.d-check.yml`-Diff (name-status); kein Inline-`//nolint`; der
  modus-bezogene `repo.mk`-Kommentar wurde berichtigt, der Ist-Stand
  (Träger: `containerEnv`-Objekt) mit `make -n compile` nachgeprüft:
  `BUILD_MODE=local` im Prozess-Umfeld → `bash scripts/bau.sh compile`,
  ohne den Wert → `docker build … --target compile`.
- **Referenzrichtung (Decken-Regel):** [`spec/spezifikation.md`](../../spec/spezifikation.md)
  verweist nur aufwärts ([`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen),
  [`SPEC-054`](../../spec/spezifikation.md)); Lastenheft und Nutzertext nennen
  Slice/ADR/Roadmap nicht; „Umsetzung“-Formulierungen folgen der
  [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)-Form
  (bereits verankerte Präzedenz, kein Down-Verweis).
- **Nutzertext-Kennungen:** [`docs/user/benutzerhandbuch.md`](../../docs/user/benutzerhandbuch.md)
  trägt 0 `LH-`/`SPEC-`/`ADR-`/`slice-`-Treffer, READMEs unberührt; die
  [`CHANGELOG.md`](../../CHANGELOG.md)-Zeile trägt Kennungen auf
  zulässiger Ebene und stützt ihre Promised/Not-promised-Weite auf die
  Messgrenzen (Container-Prozess ja, Image und Dockerfile nein —
  `docker run`-Gegenprobe) auf einem benannten Prüf-Umfang
  (Linux/`amd64`, lokale Engine).
- **Ebene-Messung (CR-Risiko 1 — grenz-Sätze gegen die Messgrenzen):**
  `TestE2E_ContainerEnv_ValuesReachTheContainerProcess` **gelaufen, PASS**
  (`docker exec`-Probe trägt beide Werte, Image-Negativprobe `docker run`
  ohne beide, erzeugtes Dockerfile ohne die Namen); die
  [`CHANGELOG.md`](../../CHANGELOG.md)-Weite „as far as the tests of this
  entry show“ und der Handbau-Satz „belegen Docker-Tests“ bleiben innerhalb
  der Messgrenze — kein [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md)-Befund.
  Die `make -n`-Kette für `BUILD_MODE=local` habe ich selbst nachgemessen
  (s. o.).
- **Drei gemeldete Grenzfälle (Bewertung):** (a) `config list` zeigt den
  leer-wertigen Eintrag nicht — bestehender gemeinsamer „paths with a
  value“-Filter, in `TestConfigGet_ContainerEnv_Extraction` mit Kommentar
  gepinnt; normativ verspricht nur `u-boot config` (bare) die gesamte
  Konfiguration ([`LH-FA-CONF-005`](../../spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern))
  — keine Zusage-Lücke, Get liest den leeren Wert verbatim. (b)
  Hand-defekte Namen: siehe F-5 (INFO). (c) Ebene-Messung stückweise:
  siehe voriger Punkt — getragen.
- **Eigenes Devcontainer:** drei Werte in [`u-boot.yaml`](../../u-boot.yaml),
  ENV-Rest hinter dem Managed-Block in `.devcontainer/Dockerfile` entfallend,
  erzeugte `containerEnv`-Objekt trägt die fünf Werte
  (`TestOwnDevcontainer_ProjectValuesUnderContainerEnv`); der
  `repo.mk`-Kommentar beschreibt den Ist-Stand, kein „früher“-Kommentar
  (§3.7-Prüfung über den ganzen Diff: keine Chronik-Kommentare, keine
  abwesenden Texte, kein Helfer zwischen Doc-Kommentar und Funktion — §3.9
  (a)–(c) sauber; die einzige verletzte Kommentar-Klasse des Diffs ist F-1).
- **§3.10-Plan-Abgleich:** 22 Dateien je `--name-status` gegen §3 — alle
  zugeordnet (Lastenheft-Zeile ist CR-umscoped „(CR-Slice zuerst)“;
  `spec/architecture.md`, ADRs, `Makefile`, `.golangci.yml`, `.d-check.yml`
  wie geändert-unverändert); Rest: F-3 (READMEs-Unternennung). Zustandsmarker:
  Roadmap-Zeile des Slices und §7 am Ist-Stand, kein abgelaufener Marker
  (die CR-Klasse [`zustandsmarker-abgelaufen-bei-uebernahme`](2026-11-03-slice-cr-devcontainer-containerenv.md)
  trat hier nicht wieder auf).
- **Wiedergefundene Klassen:** kein 3×-Vorfall; die zwei wiederkehrenden
  Zeilen der Klasse `plan-zeile-ohne-diff-und-ohne-begründung` zählen 2/3
  (F-3). Die [`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md)-Klasse
  trat nicht auf: die Plan-Aussagen über die Fundstellen (Name `u-boot.yaml`
  statt `.devcontainer/u-boot.yaml`, `repo.mk:36`) standen am Diff-Nachlass
  und stimmten.
- **Keine ADR-Berührung, kein Carveout** (name-status; CO-Index unberührt).

## Nicht geprüft

- **DoD-Closure, Verification-Evidence, `make gates`-Vollstand** — Verifier;
  der Umsetzungs-Commit ist der Berichts-Gegenstand, die Closure-Bedingungen
  (§5 des Plans) sind ein anderes Fragezeichen.
- **`make lint`** — nicht gelaufen (Docker-Lint-Stage); Grund: der
  Implementer berichtet die Gate-Runde grün, und die Review-Frage
  (Kommentar-/Vertrags-/Ebenen-Linsen) wird von keinem Lint-Lauf mitgetragen;
  ein Staticcheck-/Cyclop-Befund wäre ein Verification-Befund, kein
  Review-Finding-Beleg.
- **Voller `make test-docker`-Lauf** — nicht gelaufen; Grund: Zeitbudget —
  gelaufen ist die **zugeschnittene Ebene-Messung** des neuen Docker-Pfads
  (`TestE2E_ContainerEnv_ValuesReachTheContainerProcess`, PASS, 11.4 s) in
  derselben `test-docker-tools`-Umgebung; der Voller-Lauf gehört der
  Verification-/CI-Abdeckung (`integration-docker`).
- **Validation gegen realen Nutzerbedarf** — Validator-Rolle.
- **Der CR-Wortlaut-Freigabeakt** — außerhalb der Repo-Evidenz; die
  Übernahme-Form ist im Vorgänger-Review belegt (`45833353`).

## Sensors

| Sensor | Ergebnis |
|---|---|
| `make docs-check` | **grün** — 1065 Dateien, 0 Befunde (erneut gelaufen nach Ablage dieser Datei; der erste Lauf mit 1064 Dateien stieß auf den falschen [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet)-Anker und wurde berichtigt) |
| `make verify-depguard` | **grün** — alle 9 Regeln verifiziert (PASS-Zeilen im Log) |
| `make test` | **grün** — Docker-Test-Stage, Exit 0 |
| Ebene-Messung (zugeschnitten) | **grün** — `test-docker-tools`-Container, `-run TestE2E_ContainerEnv_ValuesReachTheContainerProcess`, `--- PASS (11.40s)` |
| `make -n compile` mit/ohne `BUILD_MODE=local` | **gemessen** — Prozess-Wert wählt `scripts/bau.sh`, Default wählt `docker build --target compile` (bestätigt den `repo.mk`-Kommentar) |
| `make lint` | nicht gelaufen — siehe „Nicht geprüft“ |
| `make test-docker` (Voller-Lauf), `make gates` | nicht gelaufen — Verifier-Scope, zugeschnittener e2e-Lauf als Gegenprobe (siehe oben) |
| `make ablage-pruefung` | nicht gelaufen — kein Planning-Baum-Eintrag im Review-Vorgang; die Berichts-Datei liegt in der erlaubten `docs/reviews/`-Ablage |

## Verdikt

**Mit Findings fortführen.** Der Umsetzungs-Diff trägt die Vertrags-Zusagen
in allen sechs bestätigten Entscheidungen (Unterpfade schreibbar,
Basis-Pfad nicht schreibbar, nur `containerEnv`-Ebene, verwaltet/idempotent,
Kollisions-Abweisung bei aktivem Toolchain-Schlüssel, Ausweisung je
Wertzeile), die Ebene-Messung läuft, kein Gate ist gelockert, kein
Carveout und kein ADR-Bedarf entsteht.

- **F-1 (HIGH, Kommentar)** ist vor der Closure mit einem kleinen
  Kommentar-Berichtigen zu lösen — kein Träger-Änderungstext, kein
  Gate-Diff.
- **F-2 (MEDIUM, Spec-Zelle)** braucht vor der Closure eine
  Entscheidung: Wertetabellen-Zelle von
  [`SPEC-040`](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)
  auf die umgesetzte Lesung (Lesen `nein`) bringen oder die Les-Frage als
  bewusste Abweichung beim Projektinhaber hochziehen.
- **F-3 (LOW)** ist ein §3.10-Rest bei der Closure; die Klasse zählt 2/3 —
  bei einem dritten Vorkommen läuft der Steering-Loop über den
  Schwellen-Übergang.
- **F-4/F-5 (INFO)** sind Beobachtungs-Anker ohne Änderungspflicht.

Der Verifier darf die Closure nicht schließen, solange F-1 und F-2 offen
sind; F-3 gehört in denselben Fix-Commit mit §3-Abgleich.