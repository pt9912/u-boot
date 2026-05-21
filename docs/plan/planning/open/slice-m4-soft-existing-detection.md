# Slice M4: Soft-Existing-Detection für `u-boot init`

## Auslöser

M3-T2 implementiert nur die **Hard**-Marker-Erkennung aus
`LH-FA-INIT-004` (`u-boot.yaml`/`compose.yaml`/`.env.example`
präsent → `ErrProjectExists`).

Die **Soft**-Erkennung fehlt komplett — Spec verlangt:

> *„Liegt keine Projektsteuerdatei vor, gilt das Verzeichnis nur als
> wahrscheinliches bestehendes Projekt, wenn mindestens drei Elemente
> aus dem Mindestumfang der Projektstruktur bereits vorhanden sind.
> In diesem Fall muss `u-boot init` im interaktiven Modus explizit
> nachfragen, ob das Verzeichnis als bestehendes Projekt behandelt
> werden soll. Im nicht-interaktiven Modus ist die automatische
> Erkennung nur dann aktiv, wenn `--assume-existing` gesetzt ist."*

Bewusste M3-Lücke (Carveout, `LH-FA-PROJDOCS-005`); muss vor
MVP-Closure erfüllt sein.

## Aufhebungsbedingung

Soft-Detection ist im `InitProjectService` aktiv:

1. Wenn kein Hard-Marker präsent: prüfe `README.md`,
   `CHANGELOG.md`, `docs/`, `scripts/`, `docker/`,
   `.devcontainer/devcontainer.json` als Soft-Indikatoren.
2. ≥ 3 Soft-Indikatoren präsent → „wahrscheinliches bestehendes
   Projekt".
3. Im nicht-interaktiven Modus ohne `--assume-existing` → Exit-Code
   `10` (`LH-FA-CLI-005A`).
4. Im interaktiven Modus → Rückfrage über den Driving-Port (neuer
   Confirmation-Port oder Erweiterung des bestehenden CLI-Adapters).
5. Mit `--assume-existing` → wie Hard-Marker.

## Akzeptanzkriterien

- `LH-FA-INIT-004` Soft-Detection ist im `InitProjectService`
  implementiert und durch Tests abgedeckt (alle drei Pfade:
  unter-3-Elemente, ≥3-Elemente ohne `--assume-existing`,
  ≥3-Elemente mit `--assume-existing`).
- CLI-Adapter bietet `--assume-existing` als Flag und reicht es
  durch (`LH-FA-CLI-005A`).
- Interaktive Rückfrage über einen Driven-Port (z. B.
  `port/driven.Confirmer` mit Fake im Test), damit die
  Application-Schicht nicht direkt `os.Stdin` liest.
- `make gates` grün.
- Eintrag in `carveouts.md` als gelöst markieren.

## Out of Scope

- `--backup`/`--force`/`--no-interactive` (gehören zu
  `slice-m3-init-flow.md` Tranche T4).
- Devcontainer-Erkennung als Soft-Marker (nur als File-Pfad
  geprüft, kein semantischer Parse).

## Bezug

- Auslösende Spec: `LH-FA-INIT-004`, `LH-FA-CLI-005A`.
- Inventar-Eintrag: [`carveouts.md`](../in-progress/carveouts.md) →
  Soft-Existing-Detection fehlt.
- Hängt von: M3-T3 (CLI-Adapter mit Flag-Parsing), M3-T4
  (`--no-interactive`/`--yes`-Verhalten).
- Slice gehört thematisch zu M4 (Doctor + Init-Härtung), nicht zu
  M3 (das ist der erste lauffähige Init-Flow).
