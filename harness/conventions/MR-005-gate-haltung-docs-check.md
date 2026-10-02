# MR-005 — Gate-Haltung: `docs-check` via direktem Container-Lauf; `scan.ignore` erweitert

- **Datum:** 2026-07-24
- **Geltungsbereich:** [`.d-check.yml`](../../.d-check.yml), [`Makefile`](../../Makefile)
  (`docs-check`).
- **Ersetzt-Baseline-Regel:** [`modul-02-harness-bootstrap.md` §Gate-Fragment `d-check.mk` (Schritt 2)](../../.harness/baseline/v6.13.0/regelwerk/modul-02-harness-bootstrap.md#gate-fragment-d-checkmk-schritt-2)
- **Adaption (neu gefasst 2026-07-25):** u-boot bindet das **tool-generierte
  Fragment** `d-check.mk` per `include` ein (`--print-mk`) und haelt `docs-check`
  als duennen Alias auf dessen `doc-check`. Der Digest-Pin lebt als
  `DCHECK_DIGEST` im `Makefile`, **nicht** im Fragment - ein Re-Generieren
  ueberschreibt ihn damit nicht. Aktive Module:
  `[links, anchors, ids, matrix, planning, targets]`.
  **Arbeitsfluss-Regel aus `planning` (2026-07-25):** Das Modul koppelt "die
  Roadmap benennt eine aktive Welle" an "ein `slice-*` liegt in `in-progress/`".
  Eine deklarierte Welle ohne Slice in Arbeit ist damit ein **Befund**, kein
  Zwischenzustand - wer einen Slice schliesst, zieht im selben Commit den
  naechsten nach oder schliesst die Welle. Bewusst uebernommen statt
  wegkonfiguriert: Eine Roadmap, die Aktivitaet behauptet, die nicht
  stattfindet, ist genau die Drift, die vorher nur per Aufmerksamkeit auffiel.
  **`targets`-Abgrenzung:** Autoritaet ist die Gate-Tabelle in
  [`harness/README.md`](../README.md) §Sensors; sie fuehrt Harness-**Sensoren**,
  nicht jede Makefile-Regel. Acht Build-/Utility-Regeln stehen einzeln benannt
  in `exempt-targets` - eine Bereichsabgrenzung, **kein** Carveout (es wird
  keine Pruefung ausgesetzt).
  **`ids`-Linkpolitik (seit 2026-07-25):** Alle vier Muster laufen mit
  `link-policy: always` - Kennungen sind auch **innerhalb von Code-Spans**
  linkpflichtig, nicht nur im nackten Fliesstext. **Kein `exempt-paths`,
  nirgends.** Ein Verzeichnis-Glob haette den Bestand *und* jedes kuenftige
  Dokument ausgenommen; stattdessen drei zeilengenaue Mechanismen:
  (1) verlinken, wo die Kennung eine echte Referenz ist - in `done/` ist das
  ausdruecklich zulaessig, weil
  [`LH-FA-PROJDOCS-003`](../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle)
  "Querverweise" als erlaubte nachtraegliche Korrektur nennt;
  (2) `d-check:ignore` **je Zeile**, wo die Kennung ein *Beleg* ist und kein
  Verweis (JSON-Payload-Beispiel, `pfad:`-Feld eines Review-Findings, das den
  Fundort zum Pruefzeitpunkt festhaelt) - die Begruendung steht an Ort und
  Stelle; (3) Bereichs-Schreibweisen (`LH-FA-INIT-001..007`) werden zu  <!-- d-check:ignore (Notations-Beispiel, keine Referenz) -->
  verlinkten Paaren aufgeloest, wie in
  [`spec/architecture.md`](../../spec/architecture.md) laengst ueblich.
  **Immutabilitaets-Sensor (seit 2026-07-25):** Das Modul `vcs` schuetzt
  **Accepted-ADRs** ueber eine Commit-Range (`make doc-immutable RANGE=…`,
  `STAGED=1` lokal) - kein Default-Modul, weil es eine Range braucht.
  `immutable-when: '^Accepted$'` trifft u-boots Form (Status als Zeile unter
  `## Status`, nicht als Inline-Feld); die Status-Trennung laeuft ueber
  `exclude-sections: [Status, Geschichte]`, weil `status-line` nur Kopf-Felder
  strippt und der zulaessige Uebergang nach `Superseded by <NNNN>-<slug>` sonst
  als Drift meldet (gemessen). **`done/`-Slices sind bewusst nicht erfasst:**
  Ihre Regel "nur korrigierend aenderbar" ist semantisch; ein Diff-Vergleich
  wuerde jede erlaubte Querverweis-Korrektur als Drift melden. Ein Sensor, der
  die eigene Regel bricht, ist schlechter als keiner.
  Das `slice`-Muster traegt zusaetzlich ein Versions-Suffix
  (`(?:\.[0-9]+)*`), damit ein Name wie `slice-...-v3.5.2` vollstaendig
  matcht statt am Punkt abzubrechen. Bewusst **kein** `MR-<NNN>`-ID-Pattern - die
  Adaptions-IDs dieses Ledgers bleiben linkfrei. `.harness/baseline/**` liegt im
  `scan.ignore` (tag-agnostischer Glob `**`), damit die repo-relativen Links der
  vendorten Regelwerk-/Template-Dateien nicht gewertet werden.
  **RTM-Trace (seit 2026-09-29):** `trace.slices` (dir `docs/plan/planning`,
  Pattern `^(slice-.+)\.md$`) schliesst u-boots Kennungsform an. Die Belegsicht
  bleibt auf die Liefer-Familien `LH-FA-*`/`LH-QA-*` (Default
  `requirements.id-pattern`) beschaenkt — Lesehinweise, Abgrenzungen,
  Zielbestimmung, Risiken u. a. sind strukturell belegfrei und wuerden als
  Waisen nur Rauschen erzeugen. `ids`-Muster um `ÄÖÜ` ergaenzt
  ([`LH-PÜ-001`](../../spec/lastenheft.md#lh-pü-001--grundfunktion)/[`LH-PÜ-002`](../../spec/lastenheft.md#lh-pü-002--hauptmodule)). `--require-complete` bleibt **aus**: advisory, weil
  ein Spec-CR neue Kennungen ohne Slice-Coverage gebiert — ein rotes Gate
  wird abgeschaltet statt befolgt. Begründung und Lückenbewertung im Slice
  [`slice-gate-rtm-traceability`](../../docs/plan/planning/done/slice-gate-rtm-traceability.md).
- **Begruendung:** Bis `v0.51.1` lief `docs-check` als handgeschriebener
  `docker run`-Aufruf, und das Fragment wurde gegen den `0.2.0`-Stand
  abgelehnt. Gegen `v0.51.1` kehrt sich die Abwaegung um: Das Fragment bringt
  `--network none` an jedem Target (ein Doku-Gate braucht kein Netz), fertige
  Targets fuer die opt-in-Module samt der jeweils rund achtzehn Glieder langen
  `--disable`-Ketten - die wachsen mit jedem neuen d-check-Modul und waeren von
  Hand eine Drift-Quelle ohne Sensor - und den Pin an einer Stelle. Der
  Target-Name bleibt `docs-check`, weil er in [`AGENTS.md`](../../AGENTS.md),
  [`harness/verification.md`](../verification.md), den CI-Workflows und dutzenden
  `done/`-Closures steht; ein Alias kostet eine Zeile, ein Umbenennen einen
  Doku-Sweep durch unveraenderliche Artefakte. Die Modul-Auswahl deckt den
  bestehenden Doku-Referenz-Vertrag. Der `scan.ignore`-Glob verengt **nicht** auf einen Tag,
  damit kuenftige vendored Staende automatisch erfasst sind; `.harness/skills/`
  (`MR-009`) bleibt ausserhalb des Baseline-Globs und damit pruefbar.
- **Fragment-Re-Generierung:** Bei einem Image-Bump wird `d-check.mk` neu
  erzeugt (`--print-mk` aus dem **neuen** Image, Ausgabe nach `d-check.mk`);
  der Digest im `Makefile` wird separat gesetzt. Das Fragment ist generiert -
  Handaenderungen daran waeren stille Drift und gehoeren stattdessen ins
  `Makefile` (eigene Targets) oder in `.d-check.yml` (Konfiguration).
- **Gate-Image-Stand:** `v0.51.1` (digest-gepinnt, Bump 2026-07-25 von `0.2.0`).
  Der Pin war 50 Releases alt geworden - es gibt fuer ihn **keine**
  Aktualitaets-Routine (anders als fuer die Regelwerk-Baseline, s. Abschnitt
  Freshness-Audit). Bump-Prozedur heute: `DCHECK_DIGEST` im `Makefile` auf den
  Digest der Zielversion, Trockenlauf gegen die **unveraenderte**
  `.d-check.yml` (Regression), dann diesen Stand nachziehen.
- **Aufloesungs-Trigger:** Der Teil "Modul-Auswahl bei d-check-Upgrade
  re-evaluieren" ist mit dem Bump auf `v0.51.1` **eingeloest**: Delta gesichtet,
  Kandidaten benannt, Auswahl bewusst vertagt (Protokoll im Slice
  [`slice-harness-dcheck-image-bump`](../../docs/plan/planning/done/slice-harness-dcheck-image-bump.md)).
  **Offen und ausdruecklich neu zu bewerten:** die Ablehnung des
  `--print-mk`-Fragments oben. Sie fiel gegen den `0.2.0`-Stand; gegen
  `v0.51.1` liefert das Fragment `--network none`, fertige Targets fuer die
  opt-in-Module samt `--disable`-Ketten und den Pin an einer Stelle. Dagegen
  steht der Namens-Bruch `doc-check` vs. u-boots `docs-check`. Eine
  Adaptions-Entscheidung altert mit ihrer Grundlage - deshalb steht der Stand
  jetzt oben im Block.
