# docs/plan/planning

Planning-Artefakte (Slice-, Tranchen- und Welle-Pläne) für u-boot.

## Lifecycle

Artefakte durchlaufen vier Verzeichnisse in dieser Reihenfolge
(Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine):

```
open/ → next/ → in-progress/ → done/
```

- Übergänge erfolgen per `git mv`, damit die Datei-Historie erhalten
  bleibt.
- Ein Artefakt darf nicht in mehreren Lifecycle-Verzeichnissen
  gleichzeitig liegen.
- Inhalte in `done/` dürfen nachträglich nur korrigierend (Tippfehler,
  Querverweise, Archiv-Hinweise) geändert werden. Substanzielle
  inhaltliche Änderungen erzeugen ein neues Artefakt in `open/` oder
  `next/` mit Verweis auf den vorhergehenden Stand.
- Vor dem Übergang nach `done/` muss der Slice eine
  Verification-Evidence nach
  [`../../../harness/verification.md`](../../../harness/verification.md)
  tragen oder auf ein eigenes Evidence-Artefakt verweisen.

## Ablage

Im Planning-Baum liegen nur Dokumente; wohin Skripte und Messprotokolle gehören, steht in
[`harness/conventions.md`](../../../harness/conventions.md#ablage-von-messprotokollen-und-skripten).

## Slices vs. Wellen — zwei Status-Mechanismen

Zwei bewusst getrennte Status-Mechanismen:

- **Slices** (und Tranchen) tragen ihren Status über das **Verzeichnis**
  (`open/` → `next/` → `in-progress/` → `done/`). Der aktive Durchlauf nimmt
  ausschließlich Slices/Tranchen auf.
- **Wellen** sind bei u-boot die **Release-Versionen**. Eine **offene** Welle liegt als flache
  Datei `welle-<kebab-slug>.md` direkt in `planning/` (kein Lifecycle-Ordner); bei Closure wandert
  sie per `git mv` nach `done/`. **Geplante** Wellen haben noch keine Datei: sie stehen in der
  [`in-progress/roadmap.md`](in-progress/roadmap.md) unter *Nächste Wellen*. Die Roadmap führt
  *Offene Wellen* (ein Zeiger je Welle-Datei), *Nächste Wellen*, Meilensteine und
  *Abgeschlossene Wellen*.

u-boot-Abweichungen vom Baseline-Default (dokumentiert in
[`harness/conventions.md`](../../../harness/conventions.md), [`MR-003`](../../../harness/conventions/MR-003-roadmap-wellen.md)/[`MR-002`](../../../harness/conventions/MR-002-carveout-inventar.md)):

- **Keine `welle-<Kennung>-results.md`** — die Closure-Notiz steht in §7 der Welle-Datei und
  verweist auf den `done/`-Release-Cut-Slice als Detailquelle (inkl. Tranchen-Hashes).
- **Keine Archivierung** der Zeitdokumente geschlossener Wellen (`make archive-welle` bleibt
  ungenutzt); `done/` behält den Volltext.
- **Aufgelöste Carveouts** wandern **nicht** nach `done/`, sondern bleiben im
  Master-Inventar [`in-progress/carveouts.md`](in-progress/carveouts.md) mit
  Audit-Trail des auflösenden Slice.

## Referenzsemantik

Planning-Artefakte folgen Baseline-Regelwerk `grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP) und
[`ADR-0013`](../adr/0013-dokumentationsreferenzmodell.md):
Slices referenzieren `LH-*` und aktive ADRs normativ. Referenzen auf
andere Slices, Carveouts oder Roadmap/Wellen sind nur Trigger-,
Owner-, Closure- oder Orchestrierungskontext und erzeugen keine
Spezifikation.

## Dateiname-Konvention

Drei verbindliche Formate, abhängig vom Artefakttyp:

- `slice-<phase>-<kebab-slug>.md` – Slice-Pläne pro Meilenstein-Phase.
  Beispiel: [`slice-m1-repo-skeleton.md`](done/slice-m1-repo-skeleton.md).
- `tranche-<nr>-<kebab-slug>.md` – Tranchen-Pläne innerhalb eines
  Slice. Beispiel: `tranche-01-init-flow.md`.

- `welle-<kebab-slug>.md` – Welle-Pläne, flach in `planning/` (siehe oben).
  Das Muster ist der Wellen-Name, nicht die Versionsnummer; die Release-Version steht im Titel.

Ein Artefakt verwendet genau eines der drei Formate. Slice-Pläne sind
für phasenübergreifende Vorhaben; Tranchen-Pläne zerlegen einen Slice in
inkrementell auslieferbare Stücke.

## Ausnahmen: Master-Dokumente

Zwei übergreifende Master-Dokumente folgen keinem der drei Formate
und liegen dauerhaft in `in-progress/`:

- [`in-progress/roadmap.md`](in-progress/roadmap.md) — Stand aller
  Slices und Tranchen (Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine).
- [`in-progress/carveouts.md`](in-progress/carveouts.md) — Inventar
  aller bewussten Carveouts mit Plan-Verweis (Baseline-Regelwerk `modul-07-carveouts.md`).

Das Beobachtungs-Register [`observations/`](observations/README.md) liegt dauerhaft unter `planning/`
und folgt ebenfalls keinem der Formate; Form und Kürzel stehen in dessen README und in
[`harness/conventions.md`](../../../harness/conventions.md).
