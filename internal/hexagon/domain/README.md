# internal/hexagon/domain

Reine Datentypen und invariantes Verhalten der u-boot-Domäne. Kein I/O,
keine externen Libraries
([`SPEC-012`](../../../spec/spezifikation.md#spec-012--schichten-und-verzeichnislayout-unter-internal)).

Schichtbeschreibung und Inhalt:
[`spec/architecture.md`](../../../spec/architecture.md#21-hexagondomain).
Importregeln:
[`SPEC-013`](../../../spec/spezifikation.md#spec-013--import-regel-tabelle-der-schichten),
[`SPEC-034`](../../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln).
Festlegungen, die Domänentypen tragen:
[`SPEC-036`](../../../spec/spezifikation.md#spec-036--stabilisierung-von-u-boot-up-takt-zustandsklassen-und-entscheidungsmatrix)
(Container-Zustand, Restart-Schwelle) und
[`SPEC-040`](../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config)
(Konfigurationspfade).
