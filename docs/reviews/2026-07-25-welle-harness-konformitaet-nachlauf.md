# Review-Report: `welle-harness-konformitaet-nachlauf` — 2026-07-25

**Review-Art:** Code-Review (gegen Plan + Konventionen) mit Design-Review-Anteil
für `spec/architecture.md` (gegen Architektur und Bestand). Kein Plan-Review:
die vier Slice-Pläne sind hier Prüf-Maßstab, nicht Prüf-Gegenstand.

**Gegenstand:** `11fc7b1..HEAD` (= `c35249d`), fünf Commits:
`56c7ee2`, `d5d896a`, `be3d33c`, `c7fe437`, `c35249d`

**Skill:** `.harness/skills/reviewer.md` @ `d5d896a` · <!-- d-check:ignore (Adopter-spezifischer Skill-Pfad) -->
**Modell:** `claude-opus-5[1m]` · **Datum:** 2026-07-25

**Rollentrennung:** Der Lauf erfolgte in Frischkontext gegen den Diff; der
Reviewer hat weder geplant noch implementiert
([`AGENTS.md`](../../AGENTS.md) §Role Separation).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- `.harness/skills/reviewer.md` (Prüf-Maßstab **und** Prüf-Gegenstand),
  `.harness/skills/closure-note-reviewer.md`
- [`harness/review.md`](../../harness/review.md) (kanonische Prosa, gewinnt bei
  Abweichung zum Skill), [`harness/verification.md`](../../harness/verification.md),
  [`harness/roles.md`](../../harness/roles.md)
- [`AGENTS.md`](../../AGENTS.md) §Hard Rules, §Source Precedence, §Quality Gates
- [`harness/conventions.md`](../../harness/conventions.md) (`MR-000`..`MR-009`,
  §Freshness-Audit, §Modus-Deklaration)
- Slice-Pläne der Welle:
  [`slice-harness-reviewer-skills-und-review-ablage`](../plan/planning/done/slice-harness-reviewer-skills-und-review-ablage.md),
  [`slice-harness-sub-area-modus-audit`](../plan/planning/done/slice-harness-sub-area-modus-audit.md),
  [`slice-harness-baseline-freshness-audit`](../plan/planning/done/slice-harness-baseline-freshness-audit.md),
  [`slice-harness-architecture-template-konformitaet`](../plan/planning/done/slice-harness-architecture-template-konformitaet.md)
- Vorgabe FS-1..FS-4:
  [`slice-harness-regelwerk-adoption-v3.5.1`](../plan/planning/done/slice-harness-regelwerk-adoption-v3.5.1.md)
  §7
- Berührte `LH-*`-IDs (lesend):
  [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes),
  [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung),
  [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe),
  [`LH-FA-ARCH-001`](../../spec/lastenheft.md#lh-fa-arch-001--hexagonales-pattern)..[`LH-FA-ARCH-003`](../../spec/lastenheft.md#lh-fa-arch-003--import-regeln-und-enforcement),
  [`LH-FA-PROJDOCS-001`](../../spec/lastenheft.md#lh-fa-projdocs-001--mindeststruktur),
  [`LH-FA-PROJDOCS-003`](../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle),
  [`LH-FA-PROJDOCS-005`](../../spec/lastenheft.md#lh-fa-projdocs-005--carveout-disziplin),
  [`LH-FA-PROJDOCS-006`](../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell),
  [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache)
- [`spec/architecture.md`](../../spec/architecture.md) (Sicht) und der Bestand
  unter `internal/` als Gegenprobe
- Vorherige Findings am selben Bereich: **keine** — dies ist der erste Report in
  diesem Verzeichnis (FS-2-Entstehungszeitpunkt)
- Ausgeführte Sensoren: `make docs-check` (130 Dateien / 0 Befunde auf HEAD;
  131 / 0 mit diesem Report),
  `bash -n tools/harness/fetch-baseline-cache.sh` (ok),
  `tools/harness/fetch-baseline-cache.sh --check-freshness` (Exit 3, Arbeitsbaum
  danach unverändert)

---

## Findings

Sortiert absteigend nach Kategorie. Jedes Finding folgt dem §Output-Schema des
Reviewer-Skills.

### F-1 — Modus-Aussage `internal/adapter/driving/cli` trägt den eigenen Diff nicht mehr

- `kategorie`: MEDIUM
- `quelle`: `MR-006` ([`harness/conventions.md`](../../harness/conventions.md)),
  Reviewer-Skill §MEDIUM „Modus-Drift"
- `pfad`: `harness/conventions.md:346` (Modus-Tabelle), `harness/conventions.md:229`
  (`MR-006` Adaptions-Text)
- `befund`: Die Modus-Tabelle führt `internal/adapter/driving/cli` als **Hybrid**
  mit der Graduation-Bedingung „Beide Regeln stehen in der Sicht-Spec (Abschnitt
  Fehlermodelle), geliefert von [`slice-harness-architecture-template-konformitaet`](../plan/planning/done/slice-harness-architecture-template-konformitaet.md)
  → dann GF". Diese Bedingung ist im selben Diff eingetreten:
  `spec/architecture.md` §7.2 (`c35249d`) enthält Sentinel-Schichtung und
  Dual-Classifier-Regel; die Commit-Message von `c35249d` stellt das selbst fest
  („Damit ist die Graduation-Bedingung der Hybrid-Sub-Area
  internal/adapter/driving/cli aus dem FS-3-Audit erfuellt"). Tabelle und
  `MR-006` sind unverändert geblieben.
- `risiko`: Die Modus-Deklaration ist Prüf-Kontext für nachfolgende Slices
  (der FS-4-Slice §8 beruft sich ausdrücklich auf „die auditierte Fassung"). Eine
  Hybrid-Markierung, deren Graduation bereits stattgefunden hat, hält eine
  erledigte Ausnahme künstlich offen und entwertet das Signal „BF/Hybrid = hier
  läuft Code → Doku".
- `verifizierbar`: Review-only (kein Sensor prüft Modus-Aussagen gegen den Diff).

### F-2 — „an zwei Stellen" ist am Code nicht belegt: die Envelope-Seite ist neun kommandolokale Mapper

- `kategorie`: MEDIUM
- `quelle`: [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) /
  [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe);
  DoD-Punkt „Bestand unverletzt — jede Aussage ist am Code belegbar"
  ([`slice-harness-architecture-template-konformitaet`](../plan/planning/done/slice-harness-architecture-template-konformitaet.md) §2)
- `pfad`: `spec/architecture.md:437` (§7.2 Dual-Classifier-Regel);
  gleichlautend `.harness/skills/reviewer.md:49`-`58` (HIGH-Anker #1)
- `befund`: Beide Texte formulieren die Regel als „ein Driving-Sentinel wird an
  **zwei** Stellen des Driving-Adapters geführt". Im Bestand steht der
  Exit-Code-Klassifikation (`internal/adapter/driving/cli/cli.go:274` `ExitCode`)
  keine *eine* Envelope-Abbildung gegenüber, sondern **neun** kommandolokale
  Mapper (`mapInitErrorToDiagnostic`, `mapAddErrorToDiagnostic`,
  `mapUpErrorToDiagnostic`, `mapDownErrorToDiagnostic`,
  `mapLogsErrorToDiagnostic`, `mapRemoveErrorToDiagnostic`,
  `mapConfigErrorToDiagnostic`, `mapGenerateErrorToDiagnostic`,
  `mapTemplateErrorToDiagnostic`). `driving.ErrProjectNotInitialized` wird
  entsprechend in `cli.go` **plus sieben** dieser Mapper geführt,
  `ErrComposeFileMissing` in `cli.go` plus drei, `ErrConfirmationRequired` in
  `cli.go` plus drei.
- `risiko`: Die Regel existiert, um genau die Drift beim Sentinel-Split zu
  verhindern. Wer sie wörtlich anwendet, hält einen Split für vollständig,
  sobald `ExitCode` und *ein* Mapper gepflegt sind — die übrigen Kommandos
  liefern dann weiter den alten Diagnostic-Code bei neuem Exit-Code. Der Fehler
  träfe zuerst die querliegenden Sentinels, also die mit der größten Fan-out.
- `verifizierbar`: Review-only für den Text; die Lücke selbst wäre durch einen
  Test fangbar, der je Sentinel Exit-Klasse und Envelope-Code über alle
  Kommandos vergleicht — ein solcher Sensor existiert heute nicht
  (`make test` pinnt Exit-Codes, nicht die Mapper-Abdeckung).

### F-3 — Init-Sequenz kehrt die im Code bewusst gewählte Reihenfolge um

- `kategorie`: MEDIUM
- `quelle`: DoD-Punkt „Bestand unverletzt — jede Aussage ist am Code belegbar"
  ([`slice-harness-architecture-template-konformitaet`](../plan/planning/done/slice-harness-architecture-template-konformitaet.md) §2);
  [`LH-FA-INIT-004`](../../spec/lastenheft.md#lh-fa-init-004--bestehendes-projekt-erkennen)
- `pfad`: `spec/architecture.md:330`-`334` (Sequenz „Projekt initialisieren")
- `befund`: Das Diagramm ordnet `Plan je Datei` (Zeile 330) →
  `betroffene Pfade melden` (331) →
  `Bestaetigung anfordern (nur bei Soft-Existing)` (333) →
  `Backup + Schreiben` (335). Im
  Bestand läuft die Bestätigung **vor** dem Plan:
  `internal/hexagon/application/initproject.go:385` ruft `checkSoftExisting`
  (und darin `confirmer.ConfirmTreatAsExisting`) auf, erst danach folgen
  `planTemplatedFiles` (Zeile 396) und `emitSummary` (Zeile 406). Der Code
  begründet die Reihenfolge explizit: „runs before the per-file plan so the user
  gets a single, project-level message instead of a per-file collision cascade".
- `risiko`: Die Sequenz ist die einzige Doku-Stelle, an der der Re-Init-Fluss als
  Ganzes steht; sie ist im Slice ausdrücklich als „der eigentliche Lehrfall"
  eingeführt. Eine invertierte Reihenfolge kehrt eine begründete Design-
  Entscheidung um und lädt dazu ein, sie im Code „an die Spec anzupassen".
- `verifizierbar`: Review-only (Diagramm-Inhalt ist von keinem Gate erfasst).

### F-4 — Up-Sequenz zeigt einen Use-Case-Schritt, den es nicht gibt

- `kategorie`: MEDIUM
- `quelle`: DoD-Punkt „Bestand unverletzt — jede Aussage ist am Code belegbar";
  [`LH-FA-UP-001`](../../spec/lastenheft.md#lh-fa-up-001--umgebung-starten)
- `pfad`: `spec/architecture.md:388`-`389` (Sequenz „Umgebung starten"),
  Lane-Definition `spec/architecture.md:384`
- `befund`: Das Diagramm zeigt `App->>PN: Compose-Datei validieren` (Zeile 388)
  gefolgt von `PN->>Ad: Runtime-Aufruf (read-only)` (389) als eigenen
  Use-Case-Schritt vor dem Start. Im Bestand prüft `UpService.Up`
  (`internal/hexagon/application/upservice.go:62`-`107`) über
  `checkProjectInitialized` (`FileSystem.Exists`) und `readComposeFile`
  (`FileSystem` + `YAMLCodec`) — also über Datei-Ports, nicht über die Runtime —
  und ruft danach direkt `engine.ComposeUp` auf. Die read-only Runtime-Probe
  liegt eine Ebene tiefer im Driven-Adapter (`Engine.preflight`,
  `internal/adapter/driven/docker/engine.go:68`, aufgerufen aus `ComposeUp`,
  `ComposeDown`, `ComposePs`). Ergänzend nennt die Lane
  `adapter/driven (Runtime, Clock)` (Zeile 384) weder `FileSystem`/`YAMLCodec` noch die von
  `UpService` genutzte `NetProbe`.
- `risiko`: Die Sequenz beschreibt damit eine Architektur, die von der
  implementierten abweicht: Sie legt nahe, die 11-vs-12-Klassifikation entstehe
  im Use-Case, während sie tatsächlich eine Adapter-interne Preflight-Eigenschaft
  ist. Wer die Klassifikation ändert, sucht sie an der falschen Stelle.
- `verifizierbar`: Review-only.

### F-5 — §7.1 nennt die Exit-Klassen `13`/`15` „reserviert" — der Vertrag reserviert `3`–`9` und `16`–`19`

- `kategorie`: MEDIUM
- `quelle`: [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)
- `pfad`: `spec/architecture.md:425`-`426`
- `befund`: Die Sicht-Spec schreibt: „Die Klassen `13` und `15` des
  Exit-Code-Vertrags sind reserviert und aktuell nicht belegt; `3`–`9` und
  `16`–`19` sind laut Vertrag nicht zu verwenden."
  [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) führt
  `13` und `15` dagegen als *definierte* technische Klassen („technischer
  Infrastruktur-/Feature-Source-Fehler", „technischer Ausführungsfehler außerhalb
  der fachlichen Domäne") und erlaubt sie ausdrücklich: „`13` bis `15` dürfen
  zusätzlich verwendet werden, wenn deren Bedeutung für den aufrufenden Kontext
  explizit dokumentiert ist." Das Wort „reserviert (nicht verwenden)" gehört im
  Vertrag allein zu `3`–`9` und `16`–`19`. Unbelegt sind `13`/`15` tatsächlich
  (kein `return 13`/`return 15` im Bestand) — falsch ist die Klassifikation als
  reserviert.
- `risiko`: Rang 2 widerspricht Rang 1 in genau der Tabelle, die als
  Exit-Klassen-Überblick gedacht ist. Ein Slice, der eine technische Fehlerklasse
  braucht, liest hier ein Verbot, wo der Vertrag eine dokumentationspflichtige
  Erlaubnis vorsieht — die Fehlklassifikation landet dann auf `1`.
- `verifizierbar`: Review-only (`make docs-check` prüft Referenzen, keine
  Aussagen-Konsistenz zwischen den Straten).

### F-6 — §6 beruft sich auf ein Gate, das Sequenz-Drift nicht sieht

- `kategorie`: MEDIUM
- `quelle`: Hard Rule „Verification Evidence" / [`harness/review.md`](../../harness/review.md)
  §Linse Docker-only-Gates; Reviewer-Skill §MEDIUM „Doku driftet gegen
  implementiertes Verhalten"
- `pfad`: `spec/architecture.md:304`-`307`
- `befund`: Die Präambel von §6 behauptet: „Lanes sind Architektur-Schichten,
  nicht Funktionen — diese Granularität ist per `depguard` abgesichert und
  driftet nicht, ohne dass ein Gate ausschlägt." `depguard`
  (`.golangci.yml`, Regelblöcke `application-no-adapter`,
  `port-driving-no-driven`, …) prüft ausschließlich **Import-Richtungen**. Weder
  Aufruf-Reihenfolge noch Lane-Zuordnung noch die Frage, ob ein gezeigter Schritt
  überhaupt existiert, sind davon erfasst; F-3 und F-4 sind zwei Drifts, die
  bei grünem `make lint` bestehen.
- `risiko`: Die Aussage stammt aus der Risiko-Gegenmaßnahme des Slice („Sicht
  driftet gegen Code … diese Ebene ist durch `depguard` maschinell abgesichert").
  Sie ersetzt eine reale Kontrolle durch eine behauptete und nimmt künftigen
  Lesern den Anlass, die Diagramme gegen den Code zu prüfen.
- `verifizierbar`: `make lint` (belegt, dass depguard grün ist, während F-3/F-4
  bestehen) — für die Textaussage selbst Review-only.

### F-7 — Ergebnisse werden als erledigt deklariert, während ihre Slices in `open/` liegen

- `kategorie`: MEDIUM
- `quelle`: [`AGENTS.md`](../../AGENTS.md) §Source Precedence Rang 4 und
  §Planning-Lifecycle;
  [`LH-FA-PROJDOCS-003`](../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle)
- `pfad`: `docs/plan/planning/open/slice-harness-*.md` (vier Dateien),  <!-- d-check:ignore (Fundort zum Pruefzeitpunkt) -->
  `docs/plan/planning/in-progress/roadmap.md:28`-`31`,
  `harness/conventions.md:228`-`236` (`MR-006`), `harness/conventions.md:334`
- `befund`: Die Inhalte der vier Wellen-Slices sind geliefert
  (`d5d896a`..`c35249d`), die Slice-Dateien liegen unverändert in `open/`, ihre
  DoD-Haken sind ungesetzt und ihre §7-Closure-Notizen leer; die Roadmap führt
  alle vier mit Status „open". Gleichzeitig deklariert
  `harness/conventions.md` die Ergebnisse bereits als abgeschlossen: `MR-006`
  „**Audit ausgefuehrt (2026-07-25)**" und „Aufloesungs-Trigger: Audit
  erledigt", die Modus-Tabelle „Stand: auditiert am 2026-07-25", `MR-004`
  „seit 2026-07-25 ausfuehrbar statt nur zugesagt", `MR-009` als `Accepted`-Ort
  für bereits existierende Dateien. `AGENTS.md` verortet den aktiven Slice in
  `next/` oder `in-progress/`.
- `risiko`: Wer `harness/conventions.md` liest, kann nicht erkennen, dass das
  zitierte Audit-Protokoll (`…/open/slice-harness-sub-area-modus-audit.md` §9)  <!-- d-check:ignore (Fundort zum Pruefzeitpunkt) -->
  aus einem ungeschlossenen Slice stammt. Das Lifecycle-Verzeichnis verliert
  seine Aussagekraft als Statusquelle, und die Verzeichnis-Angabe ist die
  einzige, die maschinell sichtbar wäre.
- `verifizierbar`: Review-only. **Zuständigkeitshinweis:** Der DoD-Abgleich und
  die Closure-Entscheidung gehören dem Verifier bzw. Planner
  ([`harness/verification.md`](../../harness/verification.md)); dieses Finding
  betrifft nur den beobachtbaren Widerspruch zwischen Deklaration und Ablage.

### F-8 — §10 datiert die Architektur auf einen Stand, den der neue Kopf widerlegt

- `kategorie`: LOW
- `quelle`: Kopf-Hard-Rule `spec/architecture.md:7`-`11` („keine … Closure-Daten");
  Maintainability
- `pfad`: `spec/architecture.md:506` (§10 Evolution)
- `befund`: Der neue Kopf trägt „**Letzte Änderung:** 2026-07-25"; §10 beginnt
  unverändert mit „Diese Architektur ist der Stand vom 2026-05-22." Die beiden
  Datumsangaben widersprechen sich, und die zweite ist genau die zeitliche
  Schicht, die der neue Kopf aus der Datei verweist.
- `risiko`: Gering, aber selbstwidersprüchlich in derselben Datei; die
  Kopf-Hard-Rule verliert Autorität, wenn ihr erster Verstoß drei Abschnitte
  weiter unten steht.
- `verifizierbar`: Review-only.

### F-9 — Bestätigungs-Gate flach auf `10` gesetzt, ohne die `2`-Verzweigung des Vertrags

- `kategorie`: LOW
- `quelle`: [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung),
  [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)
- `pfad`: `spec/architecture.md:418` (§7.1, Zeile „Bestätigungs-/Interaktions-Gate")
- `befund`: Die neue Tabelle ordnet das Bestätigungs-/Interaktions-Gate
  ausnahmslos der Exit-Klasse `10` zu.
  [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung)
  kennt beide Zweige: `10` für die Abweisung destruktiver Operationen und die
  implizite Bestandserkennung, aber `2` für „mit `--no-interactive` bricht der
  Aufruf bei jeder offenen Bestätigungsfrage … ab". Der Bestand hat die
  Spannung bereits aufgelöst (Use-Case-Refusal → `10`, Flag-Exklusivität → `2`),
  aber nur als Go-Doc-Kommentar in
  `internal/hexagon/port/driving/down.go:121`-`132`.
- `risiko`: Die Sicht-Spec stellt eine Vertrags-Verzweigung als eindeutig dar.
  Der Slice-Plan verlangt für genau diesen Fall ein Finding mit CR-Bedarf statt
  eines stillen Angleichs („Vertrags-Kollision", §6) — die Auflösung bleibt
  damit unterhalb der Doku-Ebene sichtbar.
- `verifizierbar`: Review-only; die Code-Seite ist durch die
  Exit-Code-Pin-Tests (`make test`) belegt.

### F-10 — Drei neue Slices erklären sich „ohne Welle", stehen aber als Wellen-Kopf in der Roadmap

- `kategorie`: LOW
- `quelle`: `MR-003` ([`harness/conventions.md`](../../harness/conventions.md)),
  Maintainability
- `pfad`: `docs/plan/planning/open/slice-harness-baseline-bump-review-v3.5.2.md:5`,  <!-- d-check:ignore (Fundort zum Pruefzeitpunkt) -->
  `docs/plan/planning/open/slice-harness-architecture-bestandsabgleich.md:5`,  <!-- d-check:ignore (Fundort zum Pruefzeitpunkt) -->
  `docs/plan/planning/open/slice-harness-internal-readme-kennungs-retrofit.md:5`;  <!-- d-check:ignore (Fundort zum Pruefzeitpunkt) -->
  Gegenstück `docs/plan/planning/in-progress/roadmap.md:41`-`43`
- `befund`: Alle drei Dateien tragen im Kopf „**Welle:** ohne Welle
  (Harness-/Spec-Wartung)". Die Roadmap führt dieselben drei unter „Nächste
  Wellen" als jeweils „wichtigsten Slice" einer benannten Welle
  („Spec-Wartung: `architecture.md` §2 …", „Harness-Wartung: Baseline-Review-Bump
  …", „Harness-Wartung: `internal/`-READMEs …").
- `risiko`: Gering; erzeugt aber zwei sich widersprechende Antworten auf „gehört
  dieser Slice zu einer Welle?" und damit Reibung beim nächsten Wellen-Schnitt.
- `verifizierbar`: Review-only.

### F-11 — `--check-freshness` belegt Exit `3`, den `LH-FA-CLI-006` für das Produkt sperrt

- `kategorie`: INFO
- `quelle`: [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)
- `pfad`: `tools/harness/fetch-baseline-cache.sh:133`
- `befund`: Der neue Modus meldet „Review-Bump fällig" mit Exit `3`. `3`–`9` sind
  im Exit-Code-Vertrag als „reserviert (nicht verwenden)" geführt. Der Vertrag
  bindet das Produkt (`u-boot`), nicht die Harness-Skripte unter `tools/` — es
  liegt **kein** Vertragsbruch vor. Die Kollision der Zahlenräume innerhalb eines
  Repos ist dennoch bewusst zu tragen.
- `risiko`: Keine unmittelbare; relevant erst, falls ein Harness-Skript je in
  einen Produkt-Pfad wandert oder ein Gate beide Exit-Räume gemeinsam auswertet.
- `verifizierbar`: Review-only.

### F-12 — Skript-Ausgaben des neuen Modus sind deutsch

- `kategorie`: INFO
- `quelle`: [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache)
- `pfad`: `tools/harness/fetch-baseline-cache.sh:100`-`132`
- `befund`: Die Meldungen von `--check-freshness` („nicht gefunden", „faellig",
  „manuell pruefen") sind deutsch, konsistent mit den bestehenden Modi
  `--verify`/re-vendor derselben Datei. [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache) bindet CLI-Ausgaben,
  Fehlermeldungen und generierte Dateien des **Produkts**; ein Harness-Skript
  ist keine Produkt-Ausgabe. Grenzfall bewusst geprüft, kein Finding höherer
  Kategorie.
- `risiko`: Keine.
- `verifizierbar`: Review-only.

### F-13 — §2.3-Sentinel-Liste ist unvollständig; der neue §2.5-Verweis zeigt darauf

- `kategorie`: INFO
- `quelle`: [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)
- `pfad`: `spec/architecture.md:87` (§2.3 „Code 10"), referenziert von
  `spec/architecture.md:119` (§2.5)
- `befund`: Die Code-10-Liste in §2.3 nennt neun Sentinels; `isValidationError`
  (`internal/adapter/driving/cli/cli.go:308`-`330` samt Helfern) führt darüber
  hinaus u. a. `ErrComposeFileMissing`, `ErrConfirmationRequired`,
  `ErrConfirmerUnavailable`, `ErrGenerateManualConflict`,
  `ErrServiceUnregistered`, `ErrDependenciesRequired` sowie die
  Template-Sentinels. Der in diesem Diff neu gefasste §2.5-Satz („Prädikat-Helper
  mappen die in §2.3 gelisteten Driving-Sentinels") verweist damit auf eine
  unvollständige Liste. Der Befund ist bekannt und als
  [`slice-harness-architecture-bestandsabgleich`](../plan/planning/done/slice-harness-architecture-bestandsabgleich.md)
  eingeplant; der Diff hat ihn nicht erzeugt, sondern nur einen Verweis darauf
  gelegt.
- `risiko`: Gering, solange der Folge-Slice läuft; ohne ihn wächst die Liste
  weiter auseinander.
- `verifizierbar`: Review-only.

---

## Negativbefunde

### Geprüft, ohne Befund

- geprüft, ohne Befund: **FS-1 „≥ 2 repo-spezifische HIGH-Regeln"** —
  `.harness/skills/reviewer.md` markiert zwei HIGH-Anker als repo-spezifisch;
  beide sind es substanziell und nicht generisch: Der Dual-Classifier-Bruch
  beschreibt eine u-boot-eigene Adapter-Konstruktion (`ExitCode` + Envelope-
  Mapper, code-seitig belegt in `cli.go:269`-`302` und den `map*ToDiagnostic`-
  Funktionen, entstanden aus u-boots eigener Review-Historie); der
  Referenzmodell-Verstoß hängt an
  [`LH-FA-PROJDOCS-006`](../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell)/[`ADR-0013`](../plan/adr/0013-dokumentationsreferenzmodell.md),
  einer repo-eigenen Erfindung mit eigener d-check-`matrix`-Klasse. Ein
  generischer Reviewer prüft weder das eine noch das andere. Zur *Formulierung*
  des ersten Ankers siehe F-2 — die Zahl der Regeln erfüllt FS-1, die Präzision
  einer davon nicht.
- geprüft, ohne Befund: **`.harness/skills/closure-note-reviewer.md`** — trägt
  die vom Slice geforderte explizite Abweichung („u-boot hat dieses Gate
  nicht"), leitet die sieben Struktur-Pflichtfelder korrekt aus
  [`harness/verification.md`](../../harness/verification.md) ab und benennt
  zwei repo-spezifische HIGHs (Hash-Platzhalter, Sensor ohne Beleg), die beide
  aus u-boots eigener Praxis stammen.
- geprüft, ohne Befund: **Skill vs. kanonische Prosa** — die Skills duplizieren
  keine Tabellen aus [`harness/review.md`](../../harness/review.md) /
  [`harness/verification.md`](../../harness/verification.md), verweisen aufwärts
  und deklarieren den Vorrang der Prosa. Die HIGH-/MEDIUM-Anker der Skill-Datei
  widersprechen den Ankern der Prosa an keiner Stelle; die Skill-Liste ist
  echte Verfeinerung, keine Abweichung.
- geprüft, ohne Befund: **`docs/reviews/README.md`** — Zweck, Dateiname-Konvention,
  Ein-Report-pro-Lauf-Regel, Negativbefund-Pflicht, Verweis auf die vendorte
  Kopiervorlage ohne Zweitkopie im Repo, Abgrenzungstabelle gegen Verification
  Evidence und Closure-Notiz. Deckt FS-2 vollständig.
- geprüft, ohne Befund: **`MR-009`** — Geltungsbereich, Adaption, beide
  Abweichungen vom Vorlagentext, Begründung gegen
  [`LH-FA-PROJDOCS-001`](../../spec/lastenheft.md#lh-fa-projdocs-001--mindeststruktur)
  (Mindest- statt Maximalstruktur, additiv, kein CR) und Auflösungs-Trigger sind
  vollständig; die Ortswahl ist mit `MR-005`/`MR-007` konsistent.
- geprüft, ohne Befund: **FS-3 „echtes Audit oder Bestätigungs-Audit"** — das
  Protokoll ([`slice-harness-sub-area-modus-audit`](../plan/planning/done/slice-harness-sub-area-modus-audit.md) §9) widerlegt den Erst-Pass an
  drei Stellen statt ihn zu bestätigen: `cmd/uboot` wird auf Aspirantin
  zurückgestuft (1/3 Achsen), `internal/**/README.md` kommt als vorher
  unsichtbare Sub-Area hinzu, `internal/adapter/driving/cli` wird von GF auf
  Hybrid korrigiert. Es findet zusätzlich eine echte Regelverletzung außerhalb
  seines Auftrags (`scan.ignore`-Glob ohne Carveout-Eintrag). Jede GF-Aussage
  trägt einen positiven Beleg (Spec-/Architektur-Stelle), nicht nur ein
  fehlendes Gegenanzeichen — Stichprobe `internal/hexagon/port` gegen die
  depguard-Blöcke `port-driving-no-driven`/`port-driven-no-driving` in
  `.golangci.yml`: bestätigt.
- geprüft, ohne Befund: **Graduation-Bedingungen der Nicht-GF-Aussagen** — beide
  Nicht-GF-Sub-Areas tragen eine benannte, überprüfbare Bedingung mit
  Folge-Slice (Hybrid → [`slice-harness-architecture-template-konformitaet`](../plan/planning/done/slice-harness-architecture-template-konformitaet.md);
  BF → [`slice-harness-internal-readme-kennungs-retrofit`](../plan/planning/done/slice-harness-internal-readme-kennungs-retrofit.md)). Dass die erste
  bereits erfüllt und nicht nachgezogen ist, steht als F-1.
- geprüft, ohne Befund: **Carveout-Disziplin
  ([`LH-FA-PROJDOCS-005`](../../spec/lastenheft.md#lh-fa-projdocs-005--carveout-disziplin))** —
  der neue Eintrag in `carveouts.md` benennt Gegenstand, Begründung, Herkunft
  des Befunds, Klasse „temporär", Folge-Slice und Trigger. Der Diff schafft
  keinen neuen Carveout, sondern inventarisiert einen bestehenden nach.
- geprüft, ohne Befund: **`--check-freshness` Read-only-Zusage** — statisch: der
  Freshness-Zweig (`fetch-baseline-cache.sh:136`-`140`) liegt vor jedem
  Schreibpfad, `check_freshness` enthält kein `mkdir`/`cp`/`rm`/`>`-Ziel unter
  `.harness/`; dynamisch: Lauf ausgeführt, `git status --porcelain` danach leer.
- geprüft, ohne Befund: **`--check-freshness` Fail-loud-Semantik** —
  `set -euo pipefail` aktiv; das `|| true` an der Tag-Pipeline entschärft nur die
  Zuweisung, der unmittelbar folgende `[ -n "$tags" ]`-Gate wandelt jeden
  Netz-, API-, Werkzeug- oder Formatfehler in `exit 1` (kein stilles „alles
  aktuell"); die Sonderfälle sind abgedeckt: leere Antwort → `exit 1`, Pin nicht
  in der Liste → `exit 1` mit eigener Meldung (die Bedingung `newer` leer ⇔ Pin
  nicht in der sortierten Liste trägt), fehlendes `curl` → `exit 1`. Der
  Rückgabewert wird über `check_freshness || rc=$?; exit "$rc"` korrekt
  propagiert; die dadurch innerhalb der Funktion ausgesetzte `errexit`-Wirkung
  ist unschädlich, weil alle Fehlerpfade explizit `exit 1` aufrufen.
  Gegenprobe: `bash -n` ok, realer Lauf → Exit `3` mit korrekter Tag-Liste.
- geprüft, ohne Befund: **`--check-freshness` gegen `--verify`** — Abgrenzung
  „Integrität ≠ Aktualität" steht im Skript-Kopf, in
  [`AGENTS.md`](../../AGENTS.md) und in
  [`harness/conventions.md`](../../harness/conventions.md) §Freshness-Audit
  wortgleich; `--verify` bleibt offline und unverändert.
- geprüft, ohne Befund: **Gate-Politik** — `--check-freshness` ist bewusst nicht
  in `make gates`/`make ci` eingetragen; das ist als Nicht-Ziel begründet
  dokumentiert und lockert kein bestehendes Gate. `Makefile`, `.golangci.yml`,
  `.github/workflows/` und Coverage-Schwelle sind im Diff unberührt.
- geprüft, ohne Befund: **Kadenz- und Zuständigkeitsdoku (FS-4)** —
  §Freshness-Audit nennt Sensor, Exit-Codes, doppelte Kadenz (ereignisgetrieben
  tragend, kalendarisch als Netz), Zuständigkeit, Auslöser (Review-Bump nach
  `MR-004`) und Nicht-Ziele; `MR-004` löst seinen bisher nur zugesagten Satz
  darauf ein. `AGENTS.md` trägt den Pointer.
- geprüft, ohne Befund: **Regelkonforme Reaktion auf den ersten Befund** — der
  Lauf meldete `v3.5.2` gegen den Pin `v3.5.1`; statt eines stillen Bumps liegt
  [`slice-harness-baseline-bump-review-v3.5.2`](../plan/planning/done/slice-harness-baseline-bump-review-v3.5.2.md) als Plan in `open/`, dessen DoD
  Delta-Lektüre und `MR-*`-Gegenprobe **vor** dem Bump verlangt. Kein
  Auto-Update, keine Pin-Mutation im Diff. Ich habe den Lauf reproduziert.
- geprüft, ohne Befund: **Referenzmodell
  ([`LH-FA-PROJDOCS-006`](../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell))** —
  `spec/architecture.md` enthält nach dem Umbau keine Verweise auf Slice,
  Carveout, Roadmap oder ADR (`grep` über `ADR`, `slice-`, `welle`, `Tranche`,
  `M<n>`: nur die Hard-Rule-Zeile selbst); die neuen Abschnitte verweisen
  ausschließlich aufwärts auf `lastenheft.md` und seitwärts auf `.golangci.yml`.
  `spec/lastenheft.md` und die ADRs sind im Diff unberührt.
- geprüft, ohne Befund: **Kopf-Hard-Rule „keine Wellen/Slices/Hashes/
  Meilenstein-Tags"** — der Bestand der Datei ist frei von Phase-Tags,
  Slice-Referenzen und LOC-/Commit-Angaben. Einzige verbleibende zeitliche
  Aussage ist das Datum in §10 (F-8).
- geprüft, ohne Befund: **Abschnitts-Umnummerierung** — die drei internen
  Prosa-Verweise sind mitgezogen (§5 depguard, §7 Fehlermodelle, §10 Evolution,
  „Abschnitt 4" im depguard-Kopf); `make docs-check` `anchors` grün; §2.1–§2.6
  blieben stabil, die externen README-Verweise auf §2.4 sind damit intakt.
- geprüft, ohne Befund: **§3 Externe Abhängigkeiten gegen den Bestand** —
  Stichprobe: `gopkg.in/yaml.v3` erscheint in Produktivcode ausschließlich unter
  `internal/adapter/driven/{yaml,templateyaml}`; Cobra ausschließlich unter
  `internal/adapter/driving/cli` (13 Dateien, keine im Wiring oder im Hexagon);
  `internal/hexagon/domain` importiert keine externe Bibliothek; von den vier
  genannten stdlib-Fähigkeiten erscheint in `hexagon/application` nur das
  Template-Rendering (`embed` + `text/template` in `templates.go`/
  `template_init.go`, kein `os`/`os/exec`). Die vier Driven-Ports `DockerProbe`,
  `DockerEngine`, `YAMLCodec`, `Git` existieren wie beschrieben.
- geprüft, ohne Befund: **§7.1 Exit-Klassen-Zuordnung gegen `ExitCode`** — jede
  Tabellenzeile ist am Code belegt: Usage → `2` (`isUsageError`), fachliche
  Validierung → `10` (`isValidationError` inkl. der `domain`-Sentinels),
  Doctor-Strict → `11` (`ErrDoctorFailures`, Schwelle tatsächlich im Adapter),
  Runtime nicht erreichbar → `11` (`driven.ErrDockerUnavailable`, im
  Driven-Adapter erzeugt), Compose-Laufzeit/Stabilisierung → `12`
  (`driven.ErrComposeRuntime`, `driving.ErrStabilizationTimeout`), technische
  Persistenz → `14` (`isFilesystemError`), Rest → `1`. Abweichungen nur bei den
  reservierten Klassen (F-5) und beim Bestätigungs-Gate (F-9).
- geprüft, ohne Befund: **§7.2 Sentinel-Schichtung** — die Aussage „Driven-
  Sentinels zuerst, danach Driving-/Application-Sentinels" deckt sich exakt mit
  `ExitCode` (`cli.go:278`-`300`) und mit dem dortigen Doc-Kommentar; die
  Begründung (äußere Ursache würde von der fachlichen Hülle verdeckt) ist
  code-konsistent.
- geprüft, ohne Befund: **§7.3 Resilienz-Muster** — „Plan vor Execute" belegt
  (`initproject.go` `planFile`/`filePlan`/`executeTemplatedFiles`),
  „Nil-tolerante Ports" belegt (`application/noop.go`, `noopLogger`/
  `noopConfirmer`/`noopProgress`), „Context an den blockierenden Rändern" belegt
  (`DockerEngine`/`NetProbe` mit `ctx`, `Clock` ohne), „Kein stilles
  Überschreiben" belegt (Backup-Slot-Reservierung, managed blocks). Damit ist
  auch die im FS-3-Audit §9.2 benannte Rückwärts-Lücke der Sub-Area
  `hexagon/application` eingelöst.
- geprüft, ohne Befund: **Add-on-Sequenz (§6)** — Reihenfolge „Ist-Zustand aus
  Konfiguration *und* Compose-Datei lesen, Zustand klassifizieren, dann
  schreiben" deckt sich mit `addservice_detect.go`/`addservice_execute.go`; die
  Idempotenz-Aussage passt zu
  [`LH-AK-006`](../../spec/lastenheft.md#lh-ak-006--idempotenz).
- geprüft, ohne Befund: **Planning-Lifecycle-Mechanik
  ([`LH-FA-PROJDOCS-003`](../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle))** —
  kein Planning-Artefakt liegt in zwei Lifecycle-Verzeichnissen (`uniq -d` über
  alle Dateinamen: nur die je eigenen `README.md`), keine Copy-Delete-Bewegung
  (`git diff --diff-filter=R` leer, es gibt schlicht keine Bewegung). Die
  Statusfrage steht getrennt als F-7.
- geprüft, ohne Befund: **Roadmap-Buchführung** — Aktuelle Welle eröffnet,
  Slice-Tabelle vollständig, FS-5 korrekt als in der Vorwelle geschlossen
  ausgewiesen, Abhängigkeitsgraph um `K` ergänzt, „Nächste Wellen" um die drei
  neuen Wartungs-Zeilen erweitert. Format-Regel „Wellen, keine Termine"
  (`MR-003`) gewahrt.
- geprüft, ohne Befund: **Commit-Messages gegen die Diffs** — alle fünf
  beschreiben, was sie tun; die Stichproben stimmen: `d5d896a` „127 Dateien"
  und `c35249d` „129 Dateien" passen zur Wachstumskurve (mein Lauf auf HEAD:
  130 Dateien; 131 nach Anlage dieses Reports), `c7fe437` beschreibt Exit-Semantik
  und ersten Befund korrekt, `be3d33c` listet die drei Korrekturen des Audits
  vollständig. Eine Message ist präziser als der Bestand, nicht ungenauer:
  `c35249d` erklärt die Graduation der Hybrid-Sub-Area für erfüllt — der Ledger
  zieht nicht nach (F-1).
- geprüft, ohne Befund: **`AGENTS.md`-Ergänzung** — sieben Zeilen, inhaltlich
  deckungsgleich mit §Freshness-Audit, ohne Konflikt zu einer kanonischen
  Quelle; die Ergänzung steht im Betriebsregelwerk-Block, nicht in den Hard
  Rules, also ohne neue Pflicht ohne Anker.
- geprüft, ohne Befund: **`make docs-check`** — 130 Dateien, 0 Befunde
  (`links`, `anchors`, `ids`, `matrix`), inklusive der neuen Verzeichnisse
  `.harness/skills/` und `docs/reviews/`. Bestätigt zugleich die `MR-005`-Zusage,
  dass `.harness/skills/` außerhalb des `scan.ignore`-Globs
  `.harness/baseline/**` liegt und tatsächlich geprüft wird.
- geprüft, ohne Befund: **`harness/README.md`** — trägt keine Skript-Modi-Zeile,
  nur den T1-Pointer auf `.harness/baseline/v3.5.1/regelwerk/README.md`. Der
  FS-4-Plan-Eintrag „update falls nötig" löst sich damit korrekt zu „nicht
  nötig" auf; kein stiller Drift.

### Nicht geprüft

- nicht geprüft: **`make lint`, `make test`, `make coverage-gate`, `make gates`,
  `make govulncheck`, `make image-scan`** — der Diff enthält kein Go-Delta
  (`git diff --stat`: keine `*.go`-Datei berührt); der engste sinnvolle Sensor
  ist `make docs-check`. Der Code unter `internal/` wurde ausschließlich
  **lesend** als Beleg-Gegenprobe herangezogen. F-6 stützt sich auf die
  Eigenschaft von depguard, nicht auf einen Lauf.
- nicht geprüft: **DoD-Abgleich der vier Slices und Verification Evidence** —
  ausdrücklich Verifier-Aufgabe
  ([`harness/verification.md`](../../harness/verification.md),
  [`harness/roles.md`](../../harness/roles.md) §Verifier). F-7 berichtet nur den
  beobachtbaren Widerspruch zwischen Deklaration und Ablage, nicht die
  Closure-Reife.
- nicht geprüft: **Closure-Notizen** — es existiert in dieser Welle keine; der
  Schwester-Skill `closure-note-reviewer.md` gilt nur für `done/`
  und wurde hier nur als *Gegenstand* gelesen, nicht angewendet.
- nicht geprüft: **Validation gegen realen Nutzerbedarf** (README/Quickstart,
  Release-Story) — Validator-Aufgabe; der Diff berührt keinen Nutzerpfad.
- nicht geprüft: **Vendorter Baseline-Bestand `.harness/baseline/v3.5.1/**`** —
  im Diff unberührt und per `scan.ignore` bewusst außerhalb der Referenz-Prüfung
  (`MR-005`). Die Integrität wurde nicht neu verifiziert; `--verify` ist im
  Diff unverändert.
- nicht geprüft: **Inhaltlicher Abgleich `v3.5.1` gegen `v3.5.2`** — das ist der
  Auftrag des ausgelösten [`slice-harness-baseline-bump-review-v3.5.2`](../plan/planning/done/slice-harness-baseline-bump-review-v3.5.2.md) und nicht
  Gegenstand dieses Diffs.
- nicht geprüft: **`internal/**/README.md` gegen den Code-Bestand** — die
  BF-Diagnose des Audits wurde stichprobenfrei übernommen; der Retrofit liegt
  als eigener Slice vor. Eine eigene Inventur hätte den Diff-Scope verlassen.
- nicht geprüft: **Vorherige Findings am selben Bereich** — es gibt keine;
  dieser Report ist der erste unter `docs/reviews/`. Der Steering-Loop
  („dasselbe Finding zum dritten Mal") ist damit noch nicht anwendbar.

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 7 |
| LOW | 3 |
| INFO | 3 |

Verteilung nach Artefakt: `spec/architecture.md` 6 (F-2 bis F-6, F-8, F-9,
F-13 anteilig), `harness/conventions.md` 2 (F-1, F-7),
`.harness/skills/reviewer.md` 1 (F-2, geteilt), Planning/Roadmap 2 (F-7, F-10),
`tools/harness/fetch-baseline-cache.sh` 2 (F-11, F-12, beide INFO).

## Verdikt

**Merge-blockierend:** ja — für F-1 bis F-7 (MEDIUM). Nach
[`harness/review.md`](../../harness/review.md) §Finding Categories sind MEDIUM
„normalerweise vor Merge zu klären"; hier kommt hinzu, dass fünf der sieben
(F-2 bis F-6) denselben, im Slice selbst formulierten DoD-Punkt verfehlen —
„Bestand unverletzt: jede Aussage ist am Code belegbar; abweichende Funde werden
zu Findings, nicht stillschweigend weggeschrieben". Der Slice hat den richtigen
Maßstab gesetzt; der neue §6/§7 hält ihn stellenweise nicht.

**Kein HIGH.** Es gibt keinen Vertragsbruch, keine Architekturverletzung, keine
Gate-Lockerung und kein Dateisicherheits-Risiko: Der Diff enthält kein Go-Delta,
`docs-check` ist grün, der neue Skript-Modus ist nachweislich read-only und
fail-loud, das FS-3-Audit ist substanziell (es widerlegt den Erst-Pass an drei
Stellen und deckt eine echte Carveout-Lücke auf), und FS-1 erfüllt die
Zwei-Regeln-Vorgabe inhaltlich. Der Schwerpunkt dieses Laufs liegt deshalb im
Negativbefund-Block.

**Übergabe:** Findings gehen an die Implementation. F-7 hat zusätzlich eine
Planner-Kante (Lifecycle-Bewegung der vier Slices), F-9 eine mögliche
Architect-/CR-Kante. Der Report ersetzt keine Verifikation — DoD- und
Spec-Konformität prüft der Verifier separat
([`harness/verification.md`](../../harness/verification.md)); er ersetzt auch
keinen Gate-Lauf: ausgeführt wurden nur `make docs-check`, `bash -n` und der
Freshness-Modus.
