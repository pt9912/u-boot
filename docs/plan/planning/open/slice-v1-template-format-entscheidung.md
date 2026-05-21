# Slice V1: Template-Format-Entscheidung (`LH-OPEN-004`)

## Auslöser

`LH-OPEN-004` (Template-Format) ist in `spec/lastenheft.md` §14 offen:

> *„Das genaue Format für Templates ist noch festzulegen. Mögliche
> Optionen: YAML-Metadaten plus Dateivorlagen / Cookiecutter-kompatible
> Templates / eigenes Template-System / OCI-basierte Template-Pakete."*

Mit M3-T2 ist intern bereits `text/template + embed.FS` für die
init-Templates entstanden — das ist ein **applikatives** Template-
System für die u-boot-Codebase selbst, nicht das **externe**
Template-System aus `LH-FA-TPL-*` (für Zielprojekte). Das externe
Format ist weiterhin offen.

## Aufhebungsbedingung

ADR-0008 (oder gleichwertig nummeriert) trifft eine Entscheidung
zwischen mindestens diesen Optionen:

1. **YAML-Metadaten + Datei-Templates** — Template-Verzeichnis mit
   `template.yaml` als Metadaten + `text/template`-Files (entspricht
   dem aktuellen Pattern aus M3-T2, nur erweitert).
2. **Cookiecutter-kompatibel** — Standard-Format aus dem
   Python-Ökosystem; Jinja2-Templates statt `text/template`.
3. **Eigenes Template-System** — projektspezifisch, höchste
   Flexibilität, höchster Pflegeaufwand.
4. **OCI-basierte Template-Pakete** — Templates als OCI-Artefakte in
   einer Registry; Versionierung über OCI-Tags.

## Akzeptanzkriterien

- ADR mit Vergleich, Trade-offs (Ökosystem-Reichweite, Cobra-Drift,
  Pflegeaufwand), Sicherheits-Implikationen (`LH-NFA-SEC-004`) und
  konkreter Entscheidung.
- `LH-OPEN-004` in `spec/lastenheft.md` § 14 wird als entschieden
  markiert.
- `LH-FA-TPL-001..004` werden auf das gewählte Format konkretisiert.
- Eintrag in `carveouts.md` (`LH-OPEN-004`) wird entfernt.
- Folge-Slice für die Implementierung des Template-Systems wird in
  `open/` angelegt.

## Bezug

- Auslösende Spec: `LH-OPEN-004`, `LH-FA-TPL-001..004`.
- Inventar-Eintrag: [`carveouts.md`](../in-progress/carveouts.md) →
  `LH-OPEN-004`.
- Hängt von: erstem konkreten Template-Bedarf (vermutlich nach
  MVP-Closure, weil MVP nur das `basic`-Template per Default
  liefert).
- Phase: V1, weil das Template-System bis MVP-Closure nicht
  ausgerollt wird.
