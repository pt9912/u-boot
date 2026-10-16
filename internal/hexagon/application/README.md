# internal/hexagon/application

Anwendungslogik (Use-Cases). Orchestriert Domäne und Ports, enthält
keine externe I/O ([`SPEC-012`](../../../spec/spezifikation.md#spec-012--schichten-und-verzeichnislayout-unter-internal)).

Schichtbeschreibung, Komponenten und Importregeln:
[`spec/architecture.md`](../../../spec/architecture.md#22-hexagonapplication)
und [`SPEC-013`](../../../spec/spezifikation.md#spec-013--import-regel-tabelle-der-schichten).

Festlegungen, die im Code dieser Schicht umgesetzt sind:

- Stabilisierung von `up`:
  [`SPEC-036`](../../../spec/spezifikation.md#spec-036--stabilisierung-von-u-boot-up-takt-zustandsklassen-und-entscheidungsmatrix),
  [`LH-FA-UP-001.a`](../../../spec/spezifikation.md#lh-fa-up-001a--erkennung-probebarer-ports-in-composeyaml)
- Bestätigungspfad von `down --volumes`:
  [`SPEC-038`](../../../spec/spezifikation.md#spec-038--bestätigungspfad-von-u-boot-down---volumes)
- Generatoren: [`SPEC-039`](../../../spec/spezifikation.md#spec-039--zustandsmaschinen-der-generatoren-u-boot-generate)
- Konfiguration: [`SPEC-040`](../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)
- Diagnose: [`SPEC-041`](../../../spec/spezifikation.md#spec-041--schweregrade-der-prüfung-composeyamlvalid)

Welche Anforderung welcher Service erfüllt, folgt aus den Kennungen im
Code (`make doc-trace`), nicht aus einer Liste an dieser Stelle.
