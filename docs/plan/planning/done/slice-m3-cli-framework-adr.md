# Slice M3: ADR-0005 CLI-Framework Cobra

> **Status:** Done
> **DoD:** Commit `937adb1` (M3-T3 + ADR-0005)

## Auslöser

ADR-0001 (Implementierungssprache Go) hat als offenen Folgepunkt:

> *„Wahl des CLI-Frameworks (`flag` aus stdlib reicht für MVP-Stub;
> Cobra wird mit `add`/`generate`/`config`-Subkommandos
> wahrscheinlich nötig)."*

In der M3-Slice-Planung wurde Cobra per User-Entscheidung implizit
festgelegt (siehe [`slice-m3-init-flow.md`](slice-m3-init-flow.md)
„Tech-Stack-Entscheidungen für die Tranchen"), aber ohne eigenen
ADR. Damit ist die Wahl im Repo nirgends formal dokumentiert —
LH-FA-PROJDOCS-005 / -002 verlangt einen Slice-Plan plus ADR.

## Aufhebungsbedingung

M3-T3 (CLI-Adapter) bringt Cobra in `go.mod`. Mit T3 oder unmittelbar
danach wird ein eigener ADR-0005 angelegt, der die Wahl dokumentiert
und ADR-0001s Folgepunkt schließt.

## Akzeptanzkriterien

- `docs/plan/adr/0005-cli-framework-cobra.md` existiert und folgt
  dem ADR-Format aus `LH-FA-PROJDOCS-002` (Status, Datum, Kontext,
  Entscheidung, Konsequenzen).
- Inhalt nennt mindestens: Vergleich mit Alternativen (`flag` stdlib,
  `urfave/cli`), Trade-offs (Lernkurve, Dep-Größe), konkrete
  Cobra-Versions-Pin in `go.mod`.
- ADR-0001 „Offene Folgepunkte" wird um den geschlossenen
  CLI-Framework-Eintrag bereinigt (Verweis auf ADR-0005).
- Eintrag in `carveouts.md` (`ADR-0001 Folgepunkt: CLI-Framework`)
  wird entfernt oder als gelöst markiert.

## Bezug

- Auslösende ADR: `0001-implementierungssprache-go.md` Folgepunkte.
- Inventar-Eintrag: [`carveouts.md`](../in-progress/carveouts.md) →
  ADR-0001 CLI-Framework offen.
- Hängt von: M3-T3 (CLI-Adapter mit Cobra).
