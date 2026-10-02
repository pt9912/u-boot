# AGENTS.md — Briefing für AI-Coding-Agenten

## 1. Was diese Datei ist

Onboarding-Briefing für jede AI-Session, die in diesem Repo Code oder
Dokumentation ändert. Sie verweist auf die kanonischen Quellen und
formuliert die Hard Rules, die der Implementer-Agent immer
einhalten muss.

Regeln dieser Datei: Baseline-Regelwerk `modul-09-implementierung.md`
§Ziel-Form: AGENTS.md — sie trägt Hard Rules und Pointer auf kanonische
Quellen, sie dupliziert deren Inhalt nicht; sonst entsteht Drift. Sie
ersetzt keine Spec, ADR oder Slice-Doku.

**Bei Konflikt zwischen dieser Datei und einer kanonischen Quelle gilt
die kanonische Quelle** (Source Precedence — siehe
[`harness/README.md`](harness/README.md)). Dann ist diese Datei anzupassen.

Strukturregeln (ID-Schemata, Verzeichniskonvention, Adaptionen ggü.
Baseline, Modus-Deklarationen pro Sub-Area, Zusatzklassen für
Sensors-Bindung) leben in
[`harness/conventions.md`](harness/conventions.md).

Das **Regelwerk der adoptierten Baseline** (AI-Harness-Kurs,
`pt9912/ai-harness-course`) ist die **präsente, nachschlagbare Vertiefung**
zu diesem Briefing: ein self-navigierbares **Modul-Bundle**
(`.harness/baseline/v6.13.0/regelwerk/`, `README.md` = Index), **committet
vendored** samt `SHA256SUMS`, netzlos auf jedem Checkout. Das
Integritäts-Manifest `.harness/baseline/v6.13.0/SHA256SUMS` ist offline
prüfbar per `tools/harness/fetch-baseline-cache.sh --verify`.

Ob der gepinnte Stand noch der aktuelle ist, beantwortet ein zweiter,
netzgebundener Lauf: `tools/harness/fetch-baseline-cache.sh --check-freshness`
(Exit `3` = Review-Bump fällig). Jeder Harness-/Baseline-Slice führt ihn aus und
hält das Ergebnis in seiner Evidence fest — auch den Negativbefund. Kadenz,
Zuständigkeit und Nicht-Ziele:
[`MR-004`](harness/conventions/MR-004-regelwerk-vendored.md) §Freshness-Audit.

Die verkörperte Form (dieses Briefing, die Konventionen, deine
ausgefüllten Artefakte) **führt**; das Regelwerk wird **pro Entscheidung
nachgeschlagen, deren operative Detailtiefe das Briefing nicht trägt**.
Dabei **nur den benötigten Abschnitt** laden (README ist der Index),
**nicht das ganze Regelwerk im Kontext halten**. Breiterer Pflicht-Blick
bleibt bei: Bootstrap, Änderung an
[`harness/conventions.md`](harness/conventions.md) (Adaptionen `MR-<NNN>`,
Source-Precedence, ID-Schema), Drift-Audit gegen die Baseline
(Baseline-Regelwerk `modul-02-harness-bootstrap.md` §Freshness-Audit der
vendored Baseline). Derivativ: bei Konflikt gelten die kanonischen Quellen.

Die **Skelett-Vorlagen** der Baseline liegen **vendored** unter
`.harness/baseline/v6.13.0/templates/` und tragen zwei Rollen: als
**Referenz-Form**, auf die das Regelwerk mit `../templates/…` als „Ziel-Form“
verweist, und als Vorlage, die beim Anlegen neuer Artefakte (ADR, Slice, Welle,
Carveout, Review-Report) **kopiert und ausgefüllt** wird statt frei zu
formulieren.

## 2. Kanonische Quellen (Source Precedence)

In dieser Reihenfolge lesen und auflösen:

1. [`spec/lastenheft.md`](spec/lastenheft.md) — Vertrag: normative Anforderungen, Akzeptanzkriterien, Exit-Code- und Sprachverträge.
2. [`spec/spezifikation.md`](spec/spezifikation.md) — Technik: Algorithmen, Schemata, Defaults, Fehler-Codes, externe Verträge; präzisiert das Lastenheft, erweitert es nie.
3. [`spec/architecture.md`](spec/architecture.md) — Sicht: Schichten, Komponenten, Importregeln, Podman-/Docker-Annahmen.
4. [`docs/plan/adr/`](docs/plan/adr/) — ADR-Verzeichnis und -Index.
5. Aktiver Slice in [`docs/plan/planning/in-progress/`](docs/plan/planning/in-progress/) oder [`next/`](docs/plan/planning/next/) samt [`roadmap.md`](docs/plan/planning/in-progress/roadmap.md) — Wellen-Sequenz, Tranchen, DoD und Closure-Bedingungen.
6. Ausführbare Harness-Verträge: [`Makefile`](Makefile), [`Dockerfile`](Dockerfile), [`.golangci.yml`](.golangci.yml) und [`.github/workflows/`](.github/workflows/).
7. Nutzer-Doku unter [`docs/user/`](docs/user/) und Maintainer-Doku unter [`docs/maintainer/`](docs/maintainer/), besonders [`docs/maintainer/quality.md`](docs/maintainer/quality.md).
8. [`README.md`](README.md), [`README.de.md`](README.de.md) und [`CHANGELOG.md`](CHANGELOG.md).
9. **AGENTS.md (diese Datei)** — Agent-Briefing und Hard Rules.
10. [`harness/README.md`](harness/README.md) — Harness-Einstieg.

## 3. Harte Regeln

### 3.1 Docker-only

`u-boot` hat keinen Host-Go-Toolchain-Vertrag. Build, Lint, Tests, Coverage
und Security-Gates laufen über `make`, das Docker nutzt. Host braucht nur
Docker und GNU `make`.

**Falsch:** `go test ./...` oder `golangci-lint run` auf dem Host als Nachweis für einen Handoff.
**Richtig:** `make test`, `make lint`, `make gates`.

**Begründung:** Toolchain-Reproduzierbarkeit + Supply-Chain-Defense. Lokale
Host-Toolchain-Befehle dürfen nicht als alleiniger Nachweis für einen Handoff
dienen.

### 3.2 Suppression-Verbot

Inline-`//nolint` bricht das Lint-Gate und ist nicht die Standardlösung.
Dauerhafte Ausnahmen leben zentral in [`.golangci.yml`](.golangci.yml) mit
Begründung und, wenn temporär, zusätzlich im Carveout-Inventar.

### 3.3 git mv + Inhaltsänderung = zwei Commits

Wenn eine Datei verschoben **und** der Inhalt umgeschrieben wird, sind das
zwei Commits — der Move-Commit bleibt rein (Git erkennt R-Rename). Welcher
zuerst kommt, sagt der Vorgang:

1. Regelfall: `git mv source target` → eigener Commit, dann Inhalt umschreiben.
2. Lifecycle-Übergang nach `done/`: erst der Inhalt (DoD-Häkchen,
   Closure-Notiz), dann der reine `git mv` — die Notiz ist die Bedingung für
   `done/`, nicht ihre Folge.

**Begründung:** Sonst fällt die Rename-Detection unter die 50%-
Similarity-Schwelle und `git log --follow` wird unzuverlässig.

### 3.4 Architektur ist sprach- und meilensteinfrei

`spec/architecture.md` darf Pfade zu **Code-Modulen** referenzieren
(`internal/hexagon/`), aber **keine** Wellen, Slices, Commit-Hashes oder
Closure-Daten. Die zeitliche Schicht lebt in
`docs/plan/planning/` und den späteren Closure-Notizen. Auch **keine
ADR-Bezüge**: Die Sicht steht im Stabilitäts-Rang über der ADR; welche ADR
eine Aussage verbindlich macht, deklariert die ADR in ihrem `Schärft:`-Feld.
Dasselbe gilt für Lastenheft und Spezifikation: Die Decken-Regel verbietet
Verweise nach unten (ADR, Slice, Carveout, Roadmap).

Diese Regel ist *verkörpert*, nicht hier entschieden — sie folgt aus dem
Sicht-Stratum (Baseline-Regelwerk `modul-03-spec.md`
§Ziel-Form: Architektur-Sicht).

### 3.5 ADRs sind nach `Accepted` immutable

Accepted ADRs werden nicht inhaltlich umgeschrieben. Neue oder geänderte
Entscheidungen entstehen als neue ADR oder als klar dokumentierte
Folgeentscheidung mit Verweis auf die alte ADR (`Supersedes`/`Superseded by`).

### 3.6 Gates dürfen nicht ohne ADR gelockert werden

Gates, Coverage-Schwellen und Architekturregeln dürfen nicht still gelockert
werden. Eine Abschwächung braucht einen expliziten Plan- oder ADR-Anker; eine
befristete Ausnahme für einen Teil (einen Layer, einen Pfad) ist keine
Senkung, sondern ein Carveout mit Trigger und Folge-Slice; die Schwelle selbst
bleibt.

### 3.7 Ein Kommentar beschreibt, was da ist

Gilt für Code, Konfiguration und Skripte — und für Zustandsfelder (unten).
Ein Kommentar trägt eine dieser Klassen — **Zusage · Kopplung · Abgrenzung ·
Rang-Zeiger · Grenze** — und schreibt an den, der die Stelle *ändert*, nicht an den, der die
Entscheidung *trifft*. Regeln dieser Sektion: Baseline-Regelwerk
`grundlagen-harness-dateien.md` §Was ein Kommentar trägt.

**Falsch:** „Ohne dieses Feld behauptete die Ausgabe einen Zustand, der nicht
vorlag“ — Konjunktiv über die verworfene Alternative.
**Richtig:** „`dryRun` ist wahr, wenn nichts geschrieben wurde“ — Indikativ über den Zustand.

**Falsch:** „die frühere Fassung prüfte nur die Länge“ — beschreibt
abwesenden Text.
**Richtig:** die geltende Zusage nennen; die vorige hält `git`.

**Zustandsfelder ebenso:** Eine `Stand`-/`Status`-Zelle in Roadmap,
Beobachtungs-Register oder Meilenstein-Tabelle nennt den Zustand und den Beleg
als auflösbaren Anker, nicht die Chronik; das Drift-Log der Roadmap trägt nur
Umplanungen, keine Schließungen und keine erreichten Meilensteine.

**Begründung:** Die Abwägung gehört in die ADR, die Historie in `git`, die
Herkunft in **ein** auflösbares Feld (`LH-*`, `ADR-*`). Was daneben steht, liest
jeder Lauf mit und bezahlt es mit Kontext.

### 3.8 Rollen und Review

Rollen sind Kontextgrenzen. Nutze [`harness/roles.md`](harness/roles.md) für
Planner-, Architect-, Implementation-, Reviewer-, Verifier- und
Validator-Verträge. Wer geplant oder implementiert hat, reviewt oder
verifiziert nicht mit demselben Eingabe-Kontext. Jeder Rollenwechsel braucht
ein Übergabe-Artefakt: Plan, ADR-Bezug, Diff, Findings, Verification-Evidence,
Validation-Evidence oder Closure-Notiz.

Reviews folgen [`harness/review.md`](harness/review.md): Findings werden als
HIGH/MEDIUM/LOW/INFO klassifiziert und mit Quelle, Risiko und Verifizierbarkeit
dokumentiert. Reviewer implementieren nicht und ersetzen keine Verification.

### 3.9 Spec-Traceability

Code-, Test- und Doku-Änderungen müssen die betroffenen `LH-*`, `ADR-*` oder
Slice-IDs kennen. Neue öffentliche CLI-Verträge brauchen mindestens einen Spec-
oder ADR-Anker und einen Test- oder Gate-Nachweis. CLI-Ausgaben,
Fehlermeldungen und generierte Dateien bleiben Englisch
([`LH-LESE-002`](spec/lastenheft.md#lh-lese-002--sprache)), auch wenn Plan- und
Spec-Dokumente deutsch sind.

### 3.10 Dokumentationsreferenzen

Referenzen zwischen Lastenheft, Spezifikation, Architektur, ADRs, Slices,
Carveouts und Roadmap/Wellen folgen
[`LH-FA-PROJDOCS-006`](spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell)
und [`ADR-0013`](docs/plan/adr/0013-dokumentationsreferenzmodell.md): Normative
Kraft existiert nur auf aufwärtsgerichteten Inter-Layer-Kanten plus
ADR-interner Lineage. Alles Richtung Slice, Carveout oder Roadmap ist Kontext,
Traceability oder Buchführung, keine Spezifikation.

### 3.11 Verification Evidence

Slice-Closure braucht Verification-Evidence nach
[`harness/verification.md`](harness/verification.md). Gates allein reichen
nicht: Die Evidence muss DoD, Spec-/ADR-IDs, ausgeführte Sensors, nicht
ausgeführte Sensors und Carveouts sichtbar verbinden.

### 3.12 Replay / Golden Sets

Generator-Änderungen folgen [`harness/replay.md`](harness/replay.md). Neue oder
geänderte CLI-Generatoren brauchen Golden Cases für Fresh-State, Idempotenz und
relevante Safety-Pfade. Intentional geänderter Output muss in Slice, Test oder
Commit begründet werden.

### 3.13 Exit-Code-Verträge

Die Klassifikation aus
[`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006--exit-codes) ist ein
Produktvertrag. Neue Subcommands müssen ihre Fehlerpfade auf die bestehenden
Exit-Code-Kategorien abbilden und Tests für relevante Sentinels pinnen.

### 3.14 Managed-Block- und Dateisicherheit

Generatoren und Re-Init-Pfade dürfen User-Dateien nicht opportunistisch
überschreiben. Nutze die vorhandenen managed-block-, Plan-and-Execute-,
Backup- und Two-Phase-Patterns. Destruktive Operationen brauchen die im
Spec/Slice verlangte Bestätigungslogik.

### 3.15 Planning-Lifecycle und Carveouts

Planning-Artefakte folgen `open/ → next/ → in-progress/ → done/`;
Lifecycle-Bewegungen erfolgen per `git mv` (siehe 3.3). Substanzielle
Änderungen an `done/`-Artefakten erzeugen einen neuen Slice statt die alte
Closure umzuschreiben. Jeder neue temporäre Carveout bekommt parallel einen
Eintrag in
[`docs/plan/planning/in-progress/carveouts.md`](docs/plan/planning/in-progress/carveouts.md)
und einen Plan-Anker.

## 4. Quality Gates

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-harness-dateien.md`
§harness/README.md als Einstiegspunkt. Der Gate-Index steht **einmal**, in
[`harness/README.md`](harness/README.md) §Sensors — dort steht auch die
*Bindung* jedes Targets. Diese Datei führt die Liste nicht.

Vor Handoff mindestens den engsten sinnvollen Sensor ausführen; für
Codeänderungen ist `make gates` der normale Abschluss, sofern die Umgebung Docker
zulässt. Kein Target nennen, das im Makefile nicht existiert — auch nicht in Prosa.

## 5. Dokumentations-Regeln

| # | Regel | Datei |
|---|---|---|
| 1 | **Anforderungs-IDs und ADR-Nummern** müssen in PRs/Commits referenziert sein — sie sagen, welche Zusage oder Entscheidung berührt ist. Struktur-IDs (`SPEC-<NNN>`, `ARC-<NNN>`) adressieren *innerhalb* der Spec und gehören nicht in die Commit-Message. | — |
| 2 | Vergeben werden IDs beim Spec-/ADR-Schreiben nach dem in `harness/conventions.md` deklarierten ID-Schema (`LH-*` im Lastenheft, `SPEC-<NNN>` in der Spezifikation, `ARC-<NNN>` in der Sicht, ADR-Nummern über den ADR-Index) — nie ad hoc im PR. | — |
| 3 | Neue ADRs müssen den ADR-Index aktualisieren. | — |
| 4 | Roadmap/Status-Geschichte lebt in `docs/plan/planning/`, nicht in `spec/architecture.md`. | — |
| 5 | Öffentliche Verträge (README, `docs/user/`, ADR-Index, Roadmap, Slice, CHANGELOG) werden nachgezogen, wenn ein öffentlicher Vertrag berührt ist. | — |

## 6. Minimal Agent Workflow

Pro Slice:

1. [`harness/README.md`](harness/README.md) lesen und die Rolle aus [`harness/roles.md`](harness/roles.md) bestimmen.
2. Relevante kanonische Quelle lesen (Source Precedence beachten); bei Review-Rolle [`harness/review.md`](harness/review.md) anwenden.
3. Betroffene Requirement-/ADR-/Slice-IDs identifizieren.
4. Kleinste sinnvolle Änderung planen.
5. Engsten nützlichen Sensor laufen lassen.
6. Repo-weiten Gate-Lauf vor Handoff (`make gates`).
7. Doku/Indizes aktualisieren, falls ein öffentlicher Vertrag berührt; Verification-Evidence nach [`harness/verification.md`](harness/verification.md), wenn ein Slice geschlossen oder ein öffentlicher Vertrag berührt wird.
8. Ausgeführte Sensors, nicht ausgeführte Sensors und verbleibende Risiken berichten — keine Erfolgsmeldung ohne Gate-Ausführung.

Dieser Workflow deckt ausschließlich die Implementer-Rolle ab. Schritt 8
ist der Rollenwechsel, kein Abschluss: Bericht → Handoff an Reviewer
(`.harness/skills/reviewer.md`, siehe `harness/README.md` §Guides) →
Verifier. Kein Self-Review — anderer Kontext findet andere Findings,
derselbe Kontext dieselben blinden Flecken (Baseline-Regelwerk
`modul-08-agentenrollen.md`).
