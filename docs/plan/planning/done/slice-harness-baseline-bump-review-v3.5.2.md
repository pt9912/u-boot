# Slice Harness: Review-Bump der Baseline v3.5.1 → v3.5.2

**Status:** open → next → in-progress → done (Datei wird durch die
Verzeichnisse bewegt).

**Welle:** noch keiner Welle zugeordnet; als Wartungs-Kandidat in [`roadmap.md`](../in-progress/roadmap.md) §Nächste Wellen geführt.

**Bezug:** Ausgelöst vom **ersten** Freshness-Audit-Lauf
(`tools/harness/fetch-baseline-cache.sh --check-freshness`, 2026-07-25,
Exit 3): Der Kurs `pt9912/ai-harness-course` trägt den Release-Tag `v3.5.2`,
der lokale Pin steht auf `v3.5.1`. Kein `LH`-/`ADR`-Bezug — reine
Baseline-Wartung nach der `MR-004`-Bump-Prozedur
([`harness/conventions.md`](../../../../harness/conventions.md)).

**Autor:** pt9912. **Datum:** 2026-07-25.

---

## 1. Ziel

Den Versions-Bump der vendorten Baseline als **eine Einheit** ausführen —
inklusive der Lese-Leistung, die ihn von einem Auto-Update unterscheidet:
Was hat sich zwischen `v3.5.1` und `v3.5.2` geändert, und macht eine dieser
Änderungen eine u-boot-Adaption (`MR-000`..`MR-009`) ungültig?

## 2. Definition of Done

- [x] **Delta gelesen, nicht nur gezogen:** Die Änderungen zwischen den beiden
  Tags sind gesichtet und in einer Kurz-Liste festgehalten (welche Module,
  welche Templates, welche Regel-Änderungen).
- [x] **Adaptions-Ledger gegengeprüft:** Für jede `MR-*`-Adaption ist
  festgehalten, ob sie unter `v3.5.2` unverändert gilt, angepasst werden muss
  oder entfällt. Ein Bump ohne diesen Abgleich wäre genau das Auto-Update, das
  die Routine verhindern soll.
- [x] **Bump als Einheit ausgeführt** (`MR-004`): (1) `**Stand:**`-Pin,
  (2) Vendor-Pfad `.harness/baseline/v3.5.2/` per Skript-Lauf,
  (3) `AGENTS.md`-Pointer, (4) `harness/README.md`-Guides-Zeile.
- [x] **Alt-Stand behandelt:** entweder `.harness/baseline/v3.5.1/` entfernt
  (Vendor-Pfad ist tag-versioniert, der Glob `.harness/baseline/**` bleibt
  tag-agnostisch) oder bewusst als Referenz behalten — mit Begründung.
- [x] `--verify` gegen den neuen Stand grün (42+ Dateien, Manifest-Deckung).
- [x] `--check-freshness` danach Exit 0.
- [x] `make docs-check` grün.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.harness/baseline/v3.5.2/` | neu (Skript-Lauf) | vendorter Stand |
| `.harness/baseline/v3.5.1/` | entfernen oder begründet behalten | ein Stand pro Repo ist die Regel |
| [`harness/conventions.md`](../../../../harness/conventions.md) | update | Pin, Adoptionsdatum, `MR-*`-Abgleich |
| [`AGENTS.md`](../../../../AGENTS.md), [`harness/README.md`](../../../../harness/README.md) | update | T1-/T2-Pointer |

## 4. Trigger

Bereits gefeuert (Freshness-Audit 2026-07-25). Einplanung ist eine
Priorisierungs-Entscheidung, kein weiterer Trigger.

## 5. Closure-Trigger

Bump vollständig, `--verify` und `--check-freshness` grün, `MR-*`-Abgleich
dokumentiert, Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- **Stiller Regel-Wechsel:** Eine Regeländerung in `v3.5.2` kann eine
  bestehende Adaption entwerten, ohne dass ein Gate ausschlägt. Genau dagegen
  steht der `MR-*`-Abgleich in der DoD.
- **Template-Drift:** Geänderte Templates wirken auf künftige Artefakte, nicht
  rückwirkend. Bestehende Artefakte werden **nicht** migriert; falls das
  Regelwerk das fordert, ist es ein eigener Slice (und bei Vertrags-Berührung
  ein Change Request, vgl. `MR-008`).
- **Bundle-Layout:** Vor dem Entpack-Lauf gilt weiterhin Schritt 0 —
  Asset-Name und innerer Tree gegen den echten Release prüfen.

## 7. Closure-Notiz (nach `done/`)

### Delta-Lektüre v3.5.1 → v3.5.2

36 der 42 vendorten Dateien unterscheiden sich — **aber nur drei substanziell**.
Der Rest ist der Tag in den Kurs-Links (`…/blob/v3.5.1/…` → `…/blob/v3.5.2/…`).
Methode: beide Bäume kopiert, `v3.5.[12]` normalisiert, erneut verglichen.

| Datei | Änderung |
|---|---|
| `regelwerk/README.md` | Stand: Kurs-Welle 33 (2026-07-23) → **Welle 34** (2026-07-24) |
| `regelwerk/grundlagen-konventionen.md` | **neuer Absatz** zum Change Request (s. u.) |
| `templates/spec/lastenheft.template.md` | Kommentar-Block über `## 7. Historie`, der denselben Punkt als Ausfüll-Hinweis spiegelt |

Die inhaltliche Neuerung in einem Satz: **„Change Request" ist bewusst kein
Harness-Konstrukt** — kein `CR-*`-ID-Schema, keine eigene Datei, kein Gate,
sondern der *externe* Vorgang der Vertragsvereinbarung. Im Repo hinterlässt ein
angenommener CR nur einen **Fußabdruck**: Version-Bump des Lastenhefts, eine
Zeile in dessen `## Historie` mit Verweis auf den externen Vorgang, und die
geänderten `LH-*`. Dazu die Hard Rule: **weder ADR noch Slice dürfen `LH-*` je
ändern** — sie referenzieren nur.

Keine Änderung an Verzeichniskonvention, Lifecycle, ID-Schemata,
Slice-/Roadmap-/ADR-/Review-Report-Templates oder Modus-/Sub-Area-Regeln.

### `MR-*`-Gegenprobe

| Adaption | Gilt unter `v3.5.2` | Aktion |
|---|---|---|
| `MR-000` Baseline-Aussage | unverändert | keine |
| `MR-001` zwei Spec-Straten ohne Technik-Stratum | unverändert | keine |
| `MR-002` Carveout-Inventar an fester Stelle | unverändert | keine |
| `MR-003` Roadmap folgt Wellen-Template | unverändert — Template inhaltlich identisch | Formulierung entpinnt (»der Baseline« statt »v3.5.1-«), damit sie den nächsten Bump überlebt |
| `MR-004` committet vendored, beide Bäume | unverändert | Pin im Titel und Bump-Historie nachgezogen |
| `MR-005` Gate-Haltung, `scan.ignore` | unverändert — Glob war bereits tag-agnostisch (`.harness/baseline/**`), hat den Bump ohne Anfassen überstanden | keine |
| `MR-006` Modus-Deklaration je Sub-Area | unverändert | keine |
| `MR-007` Ortswahl `.harness/` | unverändert | keine |
| `MR-008` ADR-Form per CR | **inhaltlich vereinbar, Fußabdruck fehlt** | Nachtrag im Block + Folge-Slice (s. u.) |
| `MR-009` Skills/Reviews-Ablage | unverändert | keine |

**Der eine echte Befund (`MR-008`):** u-boots Praxis widerspricht der neuen
Regel *nicht* — die Entscheidung zur MADR-Umstellung fiel außerhalb des Repos,
[`slice-cr-adr-format-madr`](slice-cr-adr-format-madr.md) war nur das
Ausführungs-Vehikel, und `MR-008` hält bereits fest, dass `conventions.md`
selbst nichts am Vertrag ändern darf. Was **fehlt**, ist der geforderte
Fußabdruck: [`spec/lastenheft.md`](../../../../spec/lastenheft.md) hat keinen
`## Historie`-Abschnitt und steht unverändert auf Version `0.1.0`, obwohl
[`LH-FA-PROJDOCS-002`](../../../../spec/lastenheft.md#lh-fa-projdocs-002--adr-format)
am 2026-07-24 geändert wurde. Die Vertragsänderung ist damit am Vertrag selbst
nicht ablesbar. Bewusst **nicht** in diesem Slice behoben: Der Fußabdruck
braucht Entscheidungen (Versions-Sprung, Status-Feld, Form des Verweises unter
der `matrix`-Regel), die über einen Baseline-Bump hinausgehen — ausgelagert als
[`slice-harness-lastenheft-historie-cr-fussabdruck`](slice-harness-lastenheft-historie-cr-fussabdruck.md).

### Verification Evidence

Scope:
- Slice: `slice-harness-baseline-bump-review-v3.5.2`
- IDs: **keine** Anforderung geändert.
  [`LH-FA-PROJDOCS-002`](../../../../spec/lastenheft.md#lh-fa-projdocs-002--adr-format)
  nur lesend (Befund oben, Behebung im Folge-Slice).
- Artefakte: `.harness/baseline/v3.5.2/{regelwerk,templates}/` + `SHA256SUMS`
  (42 Dateien), Entfernung von `.harness/baseline/v3.5.1/`,
  [`harness/conventions.md`](../../../../harness/conventions.md) (Pin,
  Adoptions-/Bump-Zeile, `MR-003`/`MR-004`/`MR-008`),
  [`AGENTS.md`](../../../../AGENTS.md), [`harness/README.md`](../../../../harness/README.md),
  `.harness/skills/reviewer.md`, [`docs/reviews/README.md`](../../../reviews/README.md),
  [`roadmap.md`](../in-progress/roadmap.md), `tools/harness/fetch-baseline-cache.sh`.

DoD-Abgleich: alle Punkte erfüllt. Der Alt-Stand `v3.5.1` wurde **entfernt**
(nicht als Referenz behalten): Es gibt genau einen gültigen Stand pro Repo, die
Historie liegt in Git, und ein zweiter Baum lüde zum versehentlichen Lesen der
falschen Fassung ein.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `fetch-baseline-cache.sh v3.5.2` (re-vendor) | pass | 42 Dateien, Under-Copy-Barriere grün (Quelle = vendored) |
| `fetch-baseline-cache.sh --verify` | pass | 42 Dateien, vollständig, offline gegen das neue `SHA256SUMS` |
| `fetch-baseline-cache.sh --check-freshness` | **Exit 0** | „Pin `v3.5.2` ist der neueste Release-Tag" — der Sensor, der diesen Slice ausgelöst hat, ist wieder grün |
| `make docs-check` | pass | 132 Dateien / 0 Befunde |

Traceability:
| ID / Pflicht | Beleg |
| --- | --- |
| `MR-004` Bump-Prozedur (vier Stellen als Einheit) | (1) `**Stand:**`-Pin, (2) Vendor-Pfad `.harness/baseline/v3.5.2/`, (3) `AGENTS.md`-Pointer, (4) `harness/README.md`-Guides-Zeile — plus zwei weitere Pointer, die seit der Prozedur-Formulierung entstanden sind (Skill-Datei, `docs/reviews/README.md`) |
| §Freshness-Audit (Kadenz) | ereignisgetriebener Lauf ausgeführt, Ergebnis hier festgehalten (Exit 0 nach dem Bump) |
| Integrität des neuen Stands | `SHA256SUMS` über beide Bäume, `--verify` offline grün |

Carveouts: Neu: none. Gelöst: none. Unverändert: none.

Nicht ausgeführt:
- `make gates` / `make ci` — kein Go-Delta; die Baseline ist Doku.
- Migration bestehender Artefakte auf geänderte Templates — **entfällt**: kein
  Template hat sich inhaltlich geändert (nur Link-Tags), das Lastenheft-Template
  nur um einen Ausfüll-Kommentar.

Independent Review: **nicht durchgeführt.** Begründung statt Auslassung: Der
Bump ist mechanisch (Skript-Lauf plus Pointer-Nachzug), die inhaltliche
Prüfleistung — Delta-Lektüre und `MR-*`-Gegenprobe — ist oben vollständig
offengelegt und am vendorten Bestand nachvollziehbar. Der eine substanzielle
Befund ist als Folge-Slice ausgewiesen statt still behoben. Bei einem Bump mit
geänderten Regel- oder Template-Inhalten wäre ein Frischkontext-Review
angezeigt.

Commit / Artefakt: `162141f` (Bump als Einheit: Vendor-Baum v3.5.2, Entfernung v3.5.1, alle sechs Pointer, `MR-003`/`MR-004`/`MR-008`, Folge-Plan).

### Steering-Loop-Lerneintrag

- **Der teure Teil eines Bumps ist das Lesen, nicht das Ziehen.** 36 von 42
  Dateien „geändert" sah nach einem großen Delta aus; nach Normalisierung der
  Versions-Strings blieben drei Dateien und **eine** Regel übrig. Wer nur
  `diff -rq` liest, überschätzt den Aufwand — und wer nur das Skript laufen
  lässt, übersieht die eine Regel. Für den nächsten Bump: Tag-Normalisierung ist
  der erste Schritt der Delta-Lektüre, nicht der letzte.
- **Die neue Regel traf genau die Stelle, die u-boot schon einmal schwierig
  fand.** `MR-008` war bei der Erst-Adoption der Grenzfall („darf eine
  conventions-MR den Vertrag ändern?" — nein). Die `v3.5.2`-Klarstellung
  bestätigt die damalige Auflösung und legt zugleich eine Lücke offen, die im
  Repo niemand gesehen hatte: Der Vertrag trägt keine Spur seiner eigenen
  Änderung. Ein Bump ist damit nicht nur Pflege, sondern ein Audit von außen.
- **Tag-gepinnte Formulierungen altern mit.** `MR-003` sprach von der
  „v3.5.1-`roadmap.template.md`-Struktur" und war nach dem Bump falsch, obwohl
  sich am Template nichts geändert hatte. Konsequenz: Adaptions-Texte
  referenzieren „die Baseline", die Version steht **einmal** im Pin.
- **Zwei Pointer mehr als die Prozedur kennt.** `MR-004` nennt vier
  Bump-Stellen; inzwischen existieren sechs (Skill-Datei und
  `docs/reviews/README.md` kamen mit `MR-009` dazu). Die Zahl „vier" in der
  Prozedur ist beim nächsten Bump zu prüfen, statt ihr zu vertrauen.
- **Folge-Slices:**
  [`slice-harness-lastenheft-historie-cr-fussabdruck`](slice-harness-lastenheft-historie-cr-fussabdruck.md).

## 8. Sub-Area-Modus-Begründung

Berührte Sub-Areas: *harness / Konventionen / Baseline* und
*Harness-Tooling* — beide **GF** nach
[`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration.
