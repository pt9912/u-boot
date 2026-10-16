# Welle schließen (Harness)

Argument: $ARGUMENTS

Dieser Command führt die **Planner**-Rolle für die **Wellen-Closure** (Modul 6). Eine Welle schließt
**nicht** durch einen einzelnen Slice-Übergang, sondern durch einen geordneten Ablauf, der alle ihre
Slices bündelt — **fünf Schritte, jeder hinterlässt einen Beleg, keiner ein Datum.** Erst wenn alle
fünf Belege vorliegen, ist die Welle *auditierbar* geschlossen.

**Die Welle-Plan-Datei wandert bei Closure per `git mv` nach `done/`** — der Zustand ist die
**Verzeichnis-Position, kein `Status:`-Feld** (wie beim Slice). Die Closure-Notiz steht in §7 der
Welle-Datei; der `done/`-Release-Cut-Slice ist die Detailquelle. Eine separate Results-Datei gibt es nicht.

Kanonische Quellen (vendored Regelwerk, `.harness/baseline/<tag>/regelwerk/`): Modul 6 (Roadmap /
Wellen-Closure), Modul 7 (Carveouts), Modul 5 (Lifecycle). Bei Konflikt gilt der Kurs.

## Repo-lokale Adaptionen, die du beachten MUSST

- **Docker-only + Gate-Nachweis/Stop-Hook.** Nur `make`-Targets; `make gates` endet mit
  `record-gates`. Ein neuer Commit mit Inhalt, den kein grüner Lauf deckt, lässt den Stop-Hook rot
  (streng auch ohne Commit, wenn `.harness/stop-gate-streng` liegt oder `STOP_GATE_STRENG=1` gesetzt
  ist) → nach dem Wave-Self-Close-Commit `make gates` grün bestätigen.
- **Strenges Doc-Gate.** `LH-`/`ADR-`/`MR-`-Kennungen in gescannten `.md` als klickbare Anker-Links —
  die Welle-Datei und die Roadmap werden gescannt.
- **u-boots Wellen sind die Release-Versionen** ([`MR-003`](../../harness/conventions/MR-003-roadmap-wellen.md)):
  Eine Welle bekommt ihre Datei `welle-<kebab-slug>.md` flach unter `docs/plan/planning/` erst bei der
  **Eröffnung**; geplante Wellen stehen nur in der Roadmap unter *Nächste Wellen*. Es gibt **keine**
  `welle-<Kennung>-results.md` — die Closure-Notiz steht in §7 der Welle-Datei und verweist auf den
  `done/`-Release-Cut-Slice als Detailquelle — und **keine Archivierung** (`make archive-welle` bleibt ungenutzt).
- **Beobachtungs-Register:** `docs/plan/planning/observations/`, Kennung `BEO-<KUERZEL>/<slug>`; das
  Kürzel steht in der Spalte *Kürzel* der Modus-Deklaration in `harness/conventions.md`.
- **Roadmap-Marker:** `.d-check.yml` (`planning`) erwartet im Block „Offene Wellen“ den Ruhe-Marker genau
  dann, wenn `in-progress/` keinen Slice trägt; `make doc-planning` prüft das.
- **Commit via Message-Datei** (`git commit -F <datei>`).

## Vorbedingung: Kontext lesen

1. `CLAUDE.md` (falls vorhanden), `harness/README.md`, `AGENTS.md`, `harness/conventions.md`,
   **Modul 6** on-demand und die Welle-Plan-Datei (`docs/plan/planning/<welle-id>.md`, v. a. §3
   Closure-Kriterien) lesen.

## Die fünf Schritte (Modul 6 — jeder mit Beleg, keiner mit Datum)

2. **Schritt 1 — Trigger prüfen.** Alle Slices der Welle liegen in `done/`; `make gates` grün; die
   welle-spezifischen Closure-Kriterien aus der Welle-Datei §3 sind erfüllt (z. B. ein benannter Smoke).
   Das ist die **beobachtbare** Bedingung, nicht der Kalendertag. Fehlt ein Beleg (ein Slice nicht
   `done`, ein Gate rot), **schließt die Welle nicht** — kein halbfertiges `done/`. Erzeuge die Belege
   **real** (Gate-Ausgabe, Smoke-Lauf), behaupte sie nicht.
3. **Schritt 2 — Trigger-Audit der Welle.** Drei Artefaktklassen tragen einen Trigger, alle drei
   werden geprüft: **Carveout** (Modul 7) → aufgelöst · verlängert (mit Folge-Slice) · permanent ·
   **bootstrap-aware Gate** (Modul 13) → Stufe hochschalten, oder Carveout eröffnen, wenn die neue
   Schwelle rot ist · **Entscheidung/ADR** (Modul 4) → Re-Evaluierungs-Trigger bestätigen oder
   Folge-Entscheidung mit `supersedes`. Die Welle darf *mit* dokumentiertem Carveout schließen —
   **nie** mit einem stillen roten Gate, einer stehengebliebenen Reifestufe oder einer Entscheidung,
   deren Re-Evaluierungs-Bedingung längst eintrat. Keine fälligen Trigger → je Klasse eine belegte
   „0 offen"-Feststellung, kein Auslassen. **Ein Trigger ohne Wächter ist eine Absichtserklärung mit
   Verfallsdatum.**
4. **Schritt 3 — Closure-Notiz in §7 der Welle-Datei schreiben** (Zeiger auf den `done/`-Release-Cut-Slice als
   Detailquelle; Steering-Loop-Einträge und Register-Zeiger stehen in dessen Closure-Notiz). Hält fest, *was gelernt wurde*: geliefert · was
   funktionierte · was anders lief · **Steering-Loop-Einträge** (geschärfte Regel / neuer Sensor /
   benannte Spec-Lücke) · Zeiger aufs Beobachtungs-Register · Folge-Slices · Verifikation (die Belege
   aus Schritt 1). **Ohne Lerneintrag ist die Welle nicht „fertig", nur „weg".**
   **Der Lese-Schritt des Beobachtungs-Registers gehört hierher** (`docs/plan/planning/observations/README.md`,
   Modul 6): jeder Eintrag, dessen `evidence/` **≥ 3** Dateien führt, wandert in die
   Steering-Loop-Einträge und wird zur **verkörperten Regel** mit Herkunfts-Anker
   (`seit welle-<Kennung>`). Der Ausgang steht danach in seiner `state.md`, das Verzeichnis bleibt liegen;
   *gestrichen* trägt dieselbe `state.md`, mit Begründung — still löschen macht die Beobachtung
   ununterscheidbar von einer, die es nie gab.
   Erreicht kein Eintrag 3×, ist *„kein Eintrag über der Schwelle"* die Feststellung, die in die
   Closure-Notiz gehört — Auslassen ist keine Antwort. Was **unter** 3× steht, liest diese Closure
   **nicht**; dafür ist der Sichtungs-Schritt der Slice-Planung zuständig (`/plan-welle`).
   **Zugleich gehört die Welle-Plan-Datei per `git mv` nach `done/`** — wegen der repo-lokalen Hard
   Rule 3.3 (Move ≠ Inhalt) als **eigener reiner Move-Commit** (s. Schritt 4). Der Move bricht die
   Inbound-Links (Roadmap + die Welle-Verweise der Slices) **und** die eigenen `../`-Links der Datei
   (jetzt eine Ebene tiefer) → im selben Zug reconcilen, bis `docs-check` grün ist.
   **Zum Schluss die drei Paarungen prüfen** — erst jetzt, weil sie die gerade entstandenen Einträge
   prüfen: (a) *Anker* — wo ein Steering-Loop-Eintrag das Feld `liegt in <Zielort>` trägt, existiert
   der Zielort und trägt `seit welle-<Kennung>` bzw. `seit slice-<Kennung>`; (b) *Folge-Slice* — jeder genannte
   Folge-Slice existiert als Datei irgendwo im Planning-Lifecycle, nicht nur in `open/`;
   (c) *Register* — jede genannte Kennung `BEO-<KUERZEL>/<slug>` existiert als Verzeichnis im
   Register, und jedes Verzeichnis trägt ein nicht leeres `evidence/`. Rot heißt in allen drei
   Fällen: etwas wurde versprochen und nicht angelegt.
5. **Schritt 4 — Wave-Self-Close-Commit + Move.** Der **Self-Close-Commit** (Inhalt) trägt: die
   Welle-Datei §7 (Closure-Notiz mit Zeiger auf den Release-Cut-Slice; **kein `Status:`-Feld** — der
   Zustand ist die Position) + die Roadmap-Fortschreibung (Schritt 5). **Danach** der reine
   **`git mv`-Commit** der Welle-Plan-Datei nach `done/` und der **Link-Reconciliation-Commit** (Schritt 3):
   Hard Rule 3.3 trennt Move und Inhalt, daher mehrere Commits statt des einen Baseline-Self-Close-Commits
   — die Bewegung bleibt beobachtbar (ein zusammenhängender Zug). Commits via `-F`.
6. **Schritt 5 — Roadmap fortschreiben** (`in-progress/roadmap.md`, im selben Commit): die Welle
   bekommt ihre Zeile in *Abgeschlossene Wellen* (mit Zeiger auf den Release-Cut-Slice), ihr Zeiger
   verlässt *Offene Wellen*. **Befördert wird niemand** — welche Wellen offen sind, sagen die flachen
   Welle-Dateien; die Zeilen unter *Nächste Wellen* bleiben stehen, bis ihre Welle eröffnet wird. Den
   zugehörigen Meilenstein auf *erreicht* setzen, falls die Welle ihn erfüllt; löste ein Trigger eine
   Umplanung aus, bekommt *Historische Trigger-Verschiebungen* ihren Eintrag — eine Schließung ist
   keine Umplanung und gehört nicht dorthin.

## Abschluss

7. `make gates` grün nach dem Commit bestätigen (nach einem neuen Commit gibt der Stop-Hook erst
   frei, wenn ein grüner Lauf dessen Inhalt deckt). Erst wenn **alle fünf Belege** vorliegen — Trigger · Carveout-Audit · Closure-Notiz ·
   Self-Close-Commit · fortgeschriebene Roadmap — ist die Welle auditierbar geschlossen.

**Merke (Modul 6):** Datum ist *Output*, nie *Trigger*. Wer die Welle am Kalendertag schließt, kappt
halbfertige Slices und produziert genau die Auditierbarkeits-Lücke, die der Harness verhindert.

Gates nicht überspringen. Keine Erfolgsmeldung ohne Command-Ausgabe.
