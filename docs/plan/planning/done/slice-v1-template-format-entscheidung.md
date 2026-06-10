# Slice V1: Template-Format-Entscheidung ([`LH-OPEN-004`](../../../../spec/lastenheft.md#lh-open-004--template-format-entschieden))

## Auslöser

[`LH-OPEN-004`](../../../../spec/lastenheft.md#lh-open-004--template-format-entschieden) (Template-Format) ist in `spec/lastenheft.md` §14 offen:

> *„Das genaue Format für Templates ist noch festzulegen. Mögliche
> Optionen: YAML-Metadaten plus Dateivorlagen / Cookiecutter-kompatible
> Templates / eigenes Template-System / OCI-basierte Template-Pakete."*

Mit M3-T2 ist intern bereits `text/template + embed.FS` für die
init-Templates entstanden — das ist ein **applikatives** Template-
System für die u-boot-Codebase selbst, nicht das **externe**
Template-System aus `LH-FA-TPL-*` (für Zielprojekte). Das externe
Format ist weiterhin offen.

## Aufhebungsbedingung

[ADR-0008](../../adr/0008-plugin-system-statisch.md) (oder gleichwertig nummeriert) trifft eine Entscheidung
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
  Pflegeaufwand), Sicherheits-Implikationen ([`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)) und
  konkreter Entscheidung.
- [`LH-OPEN-004`](../../../../spec/lastenheft.md#lh-open-004--template-format-entschieden) in `spec/lastenheft.md` § 14 wird als entschieden
  markiert.
- [`LH-FA-TPL-001`](../../../../spec/lastenheft.md#lh-fa-tpl-001--projektvorlagen)..[`LH-FA-TPL-004`](../../../../spec/lastenheft.md#lh-fa-tpl-004--templates-auflisten) werden auf das gewählte Format konkretisiert.
- Eintrag in `carveouts.md` ([`LH-OPEN-004`](../../../../spec/lastenheft.md#lh-open-004--template-format-entschieden)) wird entfernt.
- Folge-Slice für die Implementierung des Template-Systems wird in
  `open/` angelegt.

## Bezug

- Auslösende Spec: [`LH-OPEN-004`](../../../../spec/lastenheft.md#lh-open-004--template-format-entschieden), [`LH-FA-TPL-001`](../../../../spec/lastenheft.md#lh-fa-tpl-001--projektvorlagen)..[`LH-FA-TPL-004`](../../../../spec/lastenheft.md#lh-fa-tpl-004--templates-auflisten).
- Inventar-Eintrag: [`carveouts.md`](../in-progress/carveouts.md) →
  [`LH-OPEN-004`](../../../../spec/lastenheft.md#lh-open-004--template-format-entschieden) (mit Slice-Closure entfernt).
- Hängt von: erstem konkreten Template-Bedarf (vermutlich nach
  MVP-Closure, weil MVP nur das `basic`-Template per Default
  liefert).
- Phase: V1, weil das Template-System bis MVP-Closure nicht
  ausgerollt wird.
- **Schließung 2026-05-31:** Entscheidung in
  [ADR-0009](../../adr/0009-template-format-yaml-files.md) —
  YAML-Metadaten + `text/template`-Files (Option 1). [`LH-OPEN-004`](../../../../spec/lastenheft.md#lh-open-004--template-format-entschieden)
  in `spec/lastenheft.md` §14 als entschieden markiert; vier
  Implementierungs-Slices (template-list / template-init /
  local-templates plus mindestens ein Built-in-Template) sind
  in [ADR-0009](../../adr/0009-template-format-yaml-files.md) §Folgepunkte verbindlich genannt und werden bei
  konkretem Bedarf in `open/` angelegt.
