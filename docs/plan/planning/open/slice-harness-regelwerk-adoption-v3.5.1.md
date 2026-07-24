# Slice Harness: Erst-Adoption Betriebsregelwerk v3.5.1

**Status:** open → next → in-progress → done (Datei wird durch die
Verzeichnisse bewegt).

**Phase:** Harness-/Baseline-Wartung (kein Produkt-Meilenstein, keine Welle).

**Bezug:** Erst-Adoption der externen Harness-Konvention (AI-Harness-Kurs,
`pt9912/ai-harness-course`) als committet-vendored Baseline, Stand `v3.5.1`.
Kein `LH`-/`ADR`-Neubezug auf Anforderungsebene: reine Harness-Konformität und
Provenienz, keine funktionale Anforderung und keine Architekturentscheidung am
Produkt. Berührt die bestehende Referenzmodell-Disziplin
([`LH-FA-PROJDOCS-006`](../../../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell),
[`ADR-0013`](../../adr/0013-dokumentationsreferenzmodell.md)) nur als
Leser/Nutznießer, nicht ändernd.

**Autor:** pt9912. **Datum:** 2026-07-24.

---

## 1. Ziel

u-boot **erstmalig** auf das extern gepflegte Betriebsregelwerk (AI-Harness-Kurs)
formal adoptieren und den Stand auf `v3.5.1` (Kurs-Welle 33) pinnen. u-boot ist
bereits weitgehend harness-geformt (`AGENTS.md`, `harness/*`, `spec/*`,
`docs/plan/{adr,planning}`, d-check-Gate) und wird vom Regelwerk selbst als
Referenz-Implementierung zitiert (`tools/check_refs.py` — inzwischen deprecated
und auf d-check migriert, Commit `841ff7d`; Zwei-Straten-Klassifikation;
bootstrap-aware `coverage-gate`). Deshalb ist dies **kein
Greenfield-Bootstrap**, sondern das **Nachziehen der bislang fehlenden
Adoptionsschicht**: das Regelwerk netzlos vendoren, die repo-lokalen
Abweichungen in einer neuen `harness/conventions.md` als Adaptions-Ledger
festhalten und die Verweis-/Ignore-Kanten schließen — ohne Änderung an Code-
oder Produktverhalten.

Umfang bewusst **Kern-Adoption**; drei bekannte Konformitätslücken (Reviewer-
Skill-Datei, `docs/reviews/`, vollständiges Sub-Area-Modus-Audit des
Bestandscodes) werden als **dokumentierte Folge-Slices** ausgelagert
(§7), nach dem Muster eines abgegrenzten „reviewter Umfang / offener Punkt".

## 2. Definition of Done

<!-- Prüfbare Kriterien für den Übergang nach done/. -->

- [ ] **Fetch-/Vendor-Skript** `tools/harness/fetch-baseline-cache.sh` angelegt
  (self-contained, `set -euo pipefail`): zieht das Release-Asset
  `lab-regelwerk.zip` von `pt9912/ai-harness-course`. Tag-Quelle: ohne Argument
  der `**Stand:**`-Pin in `harness/conventions.md` (**Skript-Eingabe**, nicht
  vollumfänglicher SSoT — s. `MR-004` Bump-Prozedur, `MR-007` Ortswahl); **mit**
  Argument ein expliziter Tag. Zwei Modi: Default = re-vendor (Netz),
  `--verify` = offline Integritätsprüfung.
- [ ] **Bootstrap-Reihenfolge als Barriere** (Henne-Ei aufgelöst): Der **erste**
  Vendor-Lauf übergibt den Tag **explizit** (`fetch-baseline-cache.sh v3.5.1`),
  weil `harness/conventions.md` im selben Slice erst entsteht. Reihenfolge:
  (1) `conventions.md` mit `**Stand:** v3.5.1`-Pin schreiben → (2) Skript mit
  explizitem Tag laufen lassen → ab dann liest das Skript den Pin argumentlos.
- [ ] **Bundle-Layout v3.5.1 korrekt behandelt:** Archiv-Wurzel trägt
  `regelwerk/` **und** `templates/` als Geschwister. Das Skript vendored
  **beide** (bestätigte Adoptions-Entscheidung: Upstream-Default, damit die
  `../templates/…`-Modul-Links netzlos auflösen und u-boot keine co-located
  Blank-Templates vorhalten muss) nach `.harness/baseline/v3.5.1/regelwerk/`
  bzw. `.harness/baseline/v3.5.1/templates/`.
- [ ] **Manifest-Vollständigkeit als tragende Barriere:** `SHA256SUMS` wird über
  den **tatsächlichen** Dateibestand (`find`, nicht Top-Level-Glob) beider
  Bäume erzeugt; `--verify` prüft zusätzlich Datei-Anzahl == Manifest-Zeilen
  (nicht nur `sha256sum -c`, das auf einer Teilmenge grün bliebe). Der Leer-Fall
  scheitert ohnehin laut (`sha256sum -c` exit 1).
- [ ] **v3.5.1 sauber vendored** unter `.harness/baseline/v3.5.1/`
  (`regelwerk/` = 21 Dateien: `README.md`-Index + 3 `grundlagen-*` + 17
  `modul-*`; `templates/` = Referenz-Vorlagen; `SHA256SUMS`);
  `fetch-baseline-cache.sh --verify` offline grün. Kein Doppel-Nesting
  `regelwerk/regelwerk/`. (Datei-Zählung in §6 Schritt 0 gegen den echten Tree
  bestätigen — maßgeblich ist der `find`-Bestand, nicht diese Schätzung.)
- [ ] **`harness/conventions.md` neu angelegt** (Pflicht-Artefakt laut Regelwerk,
  Existenz Pflicht/Form Wahl) mit den sieben kanonischen Abschnitten: `Purpose`,
  `Baseline` (mit `**Stand:** v3.5.1`-Pin + `SHA256SUMS`-Bezug), `Adoptierte
  Konventions-Quellen` (extern Kurs + in-Repo vendored + Gate-Config),
  `Adaptions-Block` (MR-Ledger, s. u.), `Zusatzklassen-Deklaration für
  Sensors-Bindung`, `Modus-Deklaration pro Sub-Area`, `Glossar` (optional/
  Verweis auf Lastenheft). Dupliziert **keinen** Baseline-Text — nur verweisen.
- [ ] **Adaptions-Ledger** in `harness/conventions.md` deckt mindestens ab
  (Nummerierung `MR-000` aufwärts; jeder Block trägt Datum, Geltungsbereich,
  Adaption, Begründung, Auflösungs-Trigger):
  - `MR-000` — Baseline-Aussage: adoptiert AI-Harness-Kurs `v3.5.1`; etablierte
    ID-Schemata (`ADR-<NNNN>`, `LH-*`, `slice-<phase>-<slug>`,
    `PH`/`TC`/`CO`-Familien) und Lifecycle `open → next → in-progress → done`
    unverändert vom Baseline-Default.
  - `MR-001` — Source-Precedence-Rangordnung: u-boots 9-Rang-Ordnung mit **zwei**
    Spec-Straten (`contract_spec` = Lastenheft, `view_spec` = Architektur) und
    **ohne** Technik-Stratum (`spec/spezifikation.md`) — konform, da das
    Technik-Stratum laut Regelwerk optional ist. Repo-Klasse: Tooling/Referenz.
  - `MR-002` — Carveout-Inventar-Ort: `docs/plan/planning/in-progress/carveouts.md`
    statt `docs/plan/carveouts/CO-<NNN>-*.md` (etablierte Abweichung, bereits in
    [`AGENTS.md`](../../../../AGENTS.md) und der d-check-`matrix` verankert).
  - `MR-003` — Roadmap-Ort `in-progress/roadmap.md` (nicht
    `docs/plan/planning/roadmap.md`; löst die Regelwerk-interne Pfad-Doppelnennung
    auf); Wellen leben inline in der Roadmap, keine eigenständigen Wellen-Pläne.
  - `MR-004` — **Regelwerk-Lese-Form committet vendored + Baseline-Pin v3.5.1**
    (Kern-MR): Lese-Form ist die vendored `.harness/baseline/v3.5.1/`-Ablage
    (netzlos, offline verifizierbar); `regelwerk/` **und** `templates/` werden
    vendored (Upstream-Default, bewusst **kein** Consumer-Ausschluss der
    Templates). Die vendorten Templates tragen laut `AGENTS.template` **zwei
    Rollen**: Referenz-Form (Ziel der `../templates/…`-Verweise, netzlos weil
    parallel zu `regelwerk/`) **und** Kopiervorlage — neue Artefakte (ADR,
    Slice, Welle, Carveout, Review-Report) werden aus
    `.harness/baseline/v3.5.1/templates/` **kopiert und ausgefüllt**, nicht frei
    formuliert. **Auflösungs-Trigger / Bump-Prozedur:** ein Versions-Bump ist
    **keine** Ein-Zeilen-Änderung — er fasst mindestens vier Stellen an und ist
    als Einheit auszuführen: (1) `**Stand:**`-Pin in `harness/conventions.md`,
    (2) Vendor-Pfad `.harness/baseline/<tag>/` (Skript-Lauf), (3) `AGENTS.md`-
    Pointer, (4) `harness/README.md`-Guides-Zeile. Der Pin ist nur die
    **Skript-Eingabe**; Pfad und Verweise sind abgeleitet und mitzuziehen.
  - `MR-007` — **Ortswahl `.harness/` (dot-prefixed, committet) vs. `harness/`:**
    Vendored Baseline und künftige generierte Harness-Artefakte (`.harness/skills/`,
    s. FS-1) liegen im **dot-prefixed, aber getrackten** `.harness/`; die
    autoren-gepflegte Harness-Prosa (`README.md`, `roles.md`, `review.md`,
    `verification.md`, `replay.md`, `conventions.md`) bleibt im getrackten
    `harness/`. Bewusste Trennung: `.harness/` = maschinen-materialisierte /
    vendorte Form (Kurs-Ökosystem-Konvention, wie `d-check`/belief-agent), die
    **trotz** Dot-Prefix committet ist (Ausnahme zur „Dot = ignorieren"-Lesart —
    explizit dokumentiert, damit der Split kein Verwechslungs-Fallstrick ist);
    `harness/` = handgeschriebene Vertragsdoku. Der Dot-Prefix hält den
    generierten/vendorten Bestand optisch vom autoren-Bestand getrennt.
  - `MR-005` — Gate-Haltung: `docs-check` läuft via direktem
    `docker run $(D_CHECK_IMAGE)` (kein `d-check.mk`-`--print-mk`-Fragment);
    aktive Module `[links, anchors, ids, matrix]`; bewusst **kein** `MR-<NNN>`-
    ID-Pattern (Adaptions-IDs bleiben linkfrei); `.harness/baseline/**` neu im
    `scan.ignore`.
  - `MR-006` — Modus-Deklaration pro Sub-Area (s. §9): Doku-Sub-Areas GF;
    Bestandscode-Sub-Areas mit ehrlicher GF/BF-Einordnung + Graduation-Bedingung
    für jede BF-Markierung.
  - `MR-008` — **Verweis** (kein Vertrags-Edit): die Angleichung der ADR-Form an
    das vendored MADR-Template kollidiert mit `LH-FA-PROJDOCS-002` (Vertrags-
    Stratum) und ist ein Change Request, kein conventions-MR. Dieser Block
    dokumentiert nur den Stand und zeigt auf den CR-Slice
    [`slice-cr-adr-format-madr`](slice-cr-adr-format-madr.md) (FS-5); Bestand
    lean+grandfathered, neue ADRs MADR.
- [ ] **`AGENTS.md` §1/Source-Precedence-Kopf ergänzt** um den Pointer auf die
  vendored Baseline (Sync-Trigger **T2**): Regelwerk vendored unter
  `.harness/baseline/v3.5.1/regelwerk/`, `README.md` = Index; Lesemodell **„pro
  Entscheidung den benötigten Abschnitt nachschlagen, nicht das ganze Regelwerk
  im Kontext halten"**; breiter Pflicht-Blick nur bei Bootstrap / Änderung an
  `harness/conventions.md` / Drift-Audit. Verweis auf `harness/conventions.md`
  als Ort der Strukturregeln. **Zusätzlich** die vendorten **Templates**
  (`.harness/baseline/v3.5.1/templates/`) mit ihrer Doppelrolle nennen
  (Referenz-Form + Kopiervorlage): neue Artefakte per Kopie aus der vendorten
  Vorlage anlegen, nicht frei formulieren (`MR-004`).
- [ ] **`harness/README.md` §Guides ergänzt** (Sync-Trigger **T1**): Zeile auf
  das vendored Regelwerk `.harness/baseline/v3.5.1/regelwerk/README.md`, die
  vendorten Templates und auf `harness/conventions.md`.
- [ ] **Sync-Trigger T1/T2 im Repo definiert** (L3): Die Begriffe stammen aus
  dem Regelwerk (`modul-02`, „T1/T2 als Sync-Trigger") und sind für einen
  Repo-Leser sonst opak. `harness/conventions.md` definiert sie knapp
  (T1 = Pointer in `harness/README.md`, T2 = Pointer in `AGENTS.md`
  Source-Precedence) oder verweist auf die Baseline-Fundstelle
  `.harness/baseline/v3.5.1/regelwerk/modul-02-harness-bootstrap.md`.
- [ ] **`.d-check.yml` `scan.ignore` erweitert** um `.harness/baseline/**`
  (tag-**agnostischer** Glob `**`, nicht auf `v3.5.1` einengen), damit `links`/
  `anchors`/`ids` die vendorten Regelwerk-/Template-Dateien nicht werten
  (repo-relative `../…`-Links, `../templates/…`). Kopfkommentar der `.d-check.yml`
  entsprechend nachziehen.
- [ ] **`.gitignore` verifiziert:** `.harness/baseline/**` ist **committet**
  (nicht ignoriert). Der `.harness/cache/`-Ignore ist **defensiv / aktuell
  No-op** (L4): der Vendor-Prozess nutzt `mktemp`, nicht `.harness/cache/` —
  der Eintrag spiegelt nur die Kurs-Ökosystem-Konvention und schützt vorab gegen
  ein versehentliches Commit ephemerer Nebenprodukte. Entweder so
  gekennzeichnet ergänzen **oder** weglassen, bis ein Skript real dorthin
  schreibt.
- [ ] **`make docs-check` grün** (`links` / `anchors` / `ids` / `matrix`) mit den
  neuen Dateien im Baum; kein neuer Befund durch das Vendoring (dank
  `scan.ignore`).
- [ ] **`make gates` grün** (`lint` / `test` / `coverage-gate` / `docs-check`) —
  reine Doku-/Harness-Änderung, kein Code-Delta, daher erwartet grün ohne
  Anfassen von `lint`/`test`/`coverage`.
- [ ] **Review + Verification-Evidence** nach
  [`harness/review.md`](../../../../harness/review.md) bzw.
  [`harness/verification.md`](../../../../harness/verification.md); Slice trägt
  Closure-Notiz (§8) vor dem `git mv` nach `done/`.
- [ ] **Roadmap auf v3.5.1-Wellen-Template umgebaut** (`roadmap.template.md`):
  [`roadmap.md`](../in-progress/roadmap.md) trägt Aktuelle Welle / Nächste
  Wellen / Meilensteine / Abhängigkeitsgraph / Abgeschlossene Wellen /
  Historische Trigger-Verschiebungen; Release-Versionen = Wellen (`MR-003`);
  Harness-Wartung als aktive Welle sichtbar.
- [ ] **Template-Konformitäts-Retrofits (klein, kollisionsfrei, eingefaltet):**
  `harness/README.md` erhält `## Purpose`-Heading (Regelwerk-Pflicht, Modul 9);
  [`docs/plan/planning/README.md`](../README.md) erhält Abschnitt „Slices vs.
  Wellen — zwei Status-Mechanismen" (an `MR-002`/`MR-003` angepasst). Der
  Template-Konformitäts-Audit ist damit für die S-Posten abgearbeitet;
  `spec/architecture.md` (§Sequenz/§Fehlermodelle, M-Aufwand) bleibt bewusst
  eigenem Slice vorbehalten, der ADR-Format-CR ist FS-5.

## 3. Plan (vor Code)

<!-- Datei-/Komponenten-Ebene; Implementation-Agent erweitert im ersten Lauf. -->

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/fetch-baseline-cache.sh` | neu | Vendor-/Verify-Skript; Tag aus `**Stand:**`-Pin; vendored `regelwerk/`+`templates/`; Manifest per `find` (Vollständigkeit) |
| `.harness/baseline/v3.5.1/regelwerk/` | neu | 21 Dateien (`README.md`-Index + 3 `grundlagen-*` + 17 `modul-*`); netzlos committet |
| `.harness/baseline/v3.5.1/templates/` | neu | Referenz-Vorlagen (Upstream-Default vendored), damit `../templates/…`-Links auflösen |
| `.harness/baseline/v3.5.1/SHA256SUMS` | neu | Integritäts-/Provenienz-Pin über beide Bäume |
| `harness/conventions.md` | neu | Pflicht-Artefakt: 7 Abschnitte, Baseline-Pin, Adaptions-Ledger `MR-000..MR-008`, T1/T2-Definition |
| `AGENTS.md` | update | §1/Kopf: Pointer auf vendored Regelwerk (T2), Lesemodell „pro Entscheidung nachschlagen", Verweis auf `conventions.md` |
| `harness/README.md` | update | §Guides: Zeile auf vendored Regelwerk + `conventions.md` (T1) |
| `.d-check.yml` | update | `scan.ignore` += `.harness/baseline/**`; Kopfkommentar |
| `.gitignore` | update | `.harness/cache/` ignorieren; Baseline committet (No-op bestätigen) |
| `docs/plan/planning/in-progress/roadmap.md` | update | Harness-Wartungs-Eintrag |
| `Makefile` | prüfen | optional: `docs-check` läuft schon; keine Pflicht-Änderung (Fetch-Skript bleibt manuell/Wartung) |

## 4. Trigger

<!-- next → in-progress. -->

Manuell, sobald dieser Plan reviewt und freigegeben ist. Kein Welle-Vorgänger;
reiner Wartungs-Slice. Sinnvoll vor dem nächsten funktionalen Slice, damit alle
Sessions auf einer explizit gepinnten, netzlos lesbaren Baseline arbeiten.

## 5. Closure-Trigger

<!-- Wann ist der Slice done? -->

DoD vollständig + `make gates` grün + `fetch-baseline-cache.sh --verify` offline
grün + Code-/Doku-Review nach [`harness/review.md`](../../../../harness/review.md)
+ Verification-Evidence nach
[`harness/verification.md`](../../../../harness/verification.md) + Closure-Notiz
geschrieben. Danach `git mv` `open/` → (`next/` →) `in-progress/` → `done/`.

## 6. Risiken und offene Punkte

<!-- Was könnte schief gehen? Welche Carveouts entstehen ggf.? -->

- **Vendorte Modul-Links vs. `docs-check`:** Die Regelwerk-Module nutzen
  repo-relative Links (`../templates/…`, Nachbar-Module). Tragende Absicherung
  ist der `scan.ignore`-Glob `.harness/baseline/**` — er muss **vor** dem ersten
  `docs-check`-Lauf greifen, sonst meldet `links`/`anchors` die vendorten Dateien.
  Da wir **beide** Bäume vendoren, lösen die `../templates/…`-Links zusätzlich
  real auf (doppelte Absicherung).
- **Bundle-Layout-Annahme:** v3.5.1 legt `regelwerk/`+`templates/` als Geschwister
  ab (empirisch in der belief-agent-Vendoring-Kopie bestätigt — dort allerdings
  nur `regelwerk/` gezogen). Vor Anfassen des Entpack-Pfads Schritt 0: Asset-Name
  (`lab-regelwerk.zip`) + innerer Tree gegen den echten Release verifizieren.
- **Kein Carveout erwartet:** Reine additive Doku-/Harness-Arbeit; keine rote
  Gate-Lockerung, kein Code-Delta. Sollte wider Erwarten ein Gate rot werden, gilt
  Diagnose-vor-Carveout (lokal mit CI-Setup reproduzieren, erst dann Plan-Carveout).
- **Nicht im Scope (bewusst, s. §7):** Reviewer-Skill-Datei, `docs/reviews/`,
  vollständiges Sub-Area-Modus-Audit des Bestandscodes, Durchsetzungsschicht
  (`.claude/hooks/` — `.claude/` ist zudem git-ignoriert). Je eigener Folge-Slice.

## 7. Folge-Slices (dokumentierte offene Punkte)

<!-- Kern-Adoption jetzt; diese Punkte bewusst ausgelagert, nicht vergessen. -->

- **FS-1 — Reviewer-/Closure-Skill-Dateien:** `.harness/skills/reviewer.md` und
  `.harness/skills/closure-note-reviewer.md` mit ≥ 2 repo-spezifischen
  HIGH-Regeln, abgeleitet aus der bestehenden Prosa
  [`harness/review.md`](../../../../harness/review.md) /
  [`harness/verification.md`](../../../../harness/verification.md). Grund:
  Regelwerk fordert per-Repo Skill-Datei, nicht nur Prosa.
- **FS-2 — `docs/reviews/`-Ablage:** Verzeichnis + Konvention „ein Report pro
  Lauf, Folgeläufe als neue Datei" (Negativbefund-Disziplin). Entsteht spätestens
  mit dem ersten formalen Reviewer-Skill-Lauf (FS-1).
- **FS-3 — Sub-Area-Modus-Audit Bestandscode:** vollständige Drei-Achsen-
  Qualifikation und GF/BF-Diagnose je Code-Sub-Area (`hexagon/`, `cmd/`,
  `internal/`), falls `MR-006` (§9) nur eine pragmatische Erst-Einordnung liefert.
- **FS-4 — Freshness-Audit-Routine (optional):** dokumentierte Prüf-Kadenz der
  Kurs-Release-**Liste** (nicht nur Asset-Hash), die einen Review-Bump auslöst
  statt Auto-Update.
- **FS-5 — ADR-Format-CR:** eigener Change-Request-Slice
  [`slice-cr-adr-format-madr`](slice-cr-adr-format-madr.md), der
  `LH-FA-PROJDOCS-002` (Vertrags-Stratum) auf die vendored MADR-ADR-Template-Form
  hebt und den Accepted-Bestand grandfathered. Bewusst **nicht** in diesem
  Harness-Slice, weil er eine Rang-1-Anforderung ändert; `MR-008` verweist nur.

## 8. Closure-Notiz (nach `done/`)

### Verification Evidence

Scope:
- Slice: `slice-harness-regelwerk-adoption-v3.5.1`
- IDs: **keine** `LH-*`-/`ADR-*`-Anforderung geändert (reine Harness-/
  Baseline-Adoption); [`LH-FA-PROJDOCS-006`](../../../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell)
  und [`ADR-0013`](../../adr/0013-dokumentationsreferenzmodell.md) nur lesend
  berührt. Die Vertrags-Änderung `LH-FA-PROJDOCS-002` ist als CR ausgelagert
  (FS-5, [`slice-cr-adr-format-madr`](slice-cr-adr-format-madr.md)).
- Artefakte: `tools/harness/fetch-baseline-cache.sh`,
  `.harness/baseline/v3.5.1/{regelwerk,templates}/` (+`SHA256SUMS`, 42 Dateien),
  `harness/conventions.md` (`MR-000..MR-008`), `AGENTS.md`, `harness/README.md`,
  [`docs/plan/planning/README.md`](../README.md), `.d-check.yml`, `.gitignore`,
  [`roadmap.md`](../in-progress/roadmap.md).

DoD-Abgleich: alle Punkte erfüllt (Fetch-Skript + Bootstrap-Barriere; beide
Bäume vendored + `--verify` offline grün; `conventions.md` 7 Abschnitte +
Ledger; AGENTS.md/README T1/T2 + Templates-Doppelrolle; `scan.ignore`; Roadmap
auf Wellen-Template; Retrofits `harness/README` §Purpose + `planning/README`
Slices-vs-Wellen). Offene Punkte bewusst als FS-1..FS-5 vertagt (§7).

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | `lint + test + coverage-gate + docs-check green` (Coverage 91.40% ≥ 90%) |
| `make docs-check` | pass | 120 Dateien, 0 Befunde (`.harness/baseline/**` ignoriert) |
| `fetch-baseline-cache.sh --verify` | pass | 42 Dateien, vollständig, offline (Quelle-vs-vendored + Manifest-Deckung) |

Traceability:
| ID / Pflicht | Beleg |
| --- | --- |
| Regelwerk `modul-02` (Adoptions-Pflichtstruktur) | `.harness/baseline/v3.5.1/` vendored + `harness/conventions.md` (7 Abschnitte); offene Pflichten als FS-1/FS-2/FS-3 dokumentiert |
| [`LH-FA-PROJDOCS-006`](../../../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell) / [`ADR-0013`](../../adr/0013-dokumentationsreferenzmodell.md) | `make docs-check` `matrix`/`links`/`anchors` grün — Referenzmodell unverletzt |
| Integrität vendored Regelwerk | `SHA256SUMS`; regelwerk-Hashes byte-identisch zum Schwester-Repo belief-agent |

Carveouts: Neu: none. Gelöst: none. Unverändert: none.

Nicht ausgeführt:
- `make ci` / `make govulncheck` / `make image-scan` — reine Doku-/Harness-
  Änderung, kein Code-/Image-Delta; `make gates` ist der passende Closure-Sensor.
- `make test-docker` — kein Docker-/Compose-Verhalten berührt.

Independent Review (Frischkontext, Rollentrennung): kein HIGH; ein **MEDIUM**
(regelwerk-Kopie `maxdepth-1`, Under-Copy-Barriere deckte nur post-copy) +
LOWs — **abgearbeitet**: rekursive Kopie + echte Quelle-vs-vendored-Zählung,
Kommentar präzisiert, `$root` statt `${OLDPWD}`, `MR-004`-Abschnittsname. Re-
vendor + `--verify` nach Fix grün, regelwerk unverändert byte-identisch.

Commit / Artefakt: **pending** (Closing-Commit; Hash hier eintragen) + `git mv`
`open/` → `done/`.

### Steering-Loop-Lerneintrag

- Ist-Zustand vor Adoption schlägt Ableitung: u-boot war bereits Referenz-
  Implementierung im Regelwerk — die Adoption war Re-Vendoring, kein Greenfield.
- Template-Konformität bestehender Docs ist ein **eigener** Workstream: nur die
  kollisionsfreien S-Retrofits eingefaltet; strukturelle Umbauten (Roadmap) auf
  Ansage, Vertrags-Kollisionen (ADR-Format) als CR ausgelagert.

## 9. Sub-Area-Modus-Begründung

<!-- Pflicht laut Slice-Form: berührte Sub-Areas + Modus. -->

Dieser Slice berührt schreibend nur **Doku-/Harness-Sub-Areas**; der Bestandscode
(`hexagon/`, `cmd/`, `internal/`) wird **nicht** angefasst.

### Sub-Area: harness / Konventionen / Baseline

- **Modus:** GF (Doku-führt-Sub-Area; Konventionen werden hier erstmals als
  eigenes Artefakt niedergeschrieben, nicht aus Bestandscode rekonstruiert).
- **Konventions-Dichte:** hoch (neuer Adaptions-Ledger `MR-000..MR-008`).
- **Evidenz-/Diskrepanz-Risiko:** niedrig — Doku-führt, keine Code-Inventur; das
  Risiko ist Link-/Anker-Drift, durch `docs-check` abgedeckt.
- **Graduation-Bedingung:** n/a (GF).

### Sub-Area: spec / docs/plan (Referenz, nur lesend berührt)

- **Modus:** GF (spec-führt-Sub-Area; hier nur als Verweisziel referenziert).
- **Graduation-Bedingung:** n/a (GF).

> **Bestandscode-Einordnung (Erst-Pass, verfeinerbar über FS-3):** `hexagon/`,
> `cmd/`, `internal/` sind spec-/architektur-getrieben entstanden (Doku vor Code,
> hexagonale Regeln aus `spec/architecture.md`) und werden als **GF** geführt.
> Diese Erst-Einordnung wird in `harness/conventions.md` §Modus-Deklaration
> festgehalten; ein vollständiges Drei-Achsen-Audit je Code-Sub-Area — falls eine
> Struktur doch BF-Züge trägt (z. B. retrofit-verdächtige Bereiche) — ist als
> FS-3 ausgelagert und bekäme dort seine Graduation-Bedingung.
