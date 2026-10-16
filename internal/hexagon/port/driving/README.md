# internal/hexagon/port/driving

Interfaces, über die u-boot von außen angesprochen wird
([`SPEC-012`](../../../../spec/spezifikation.md#spec-012--schichten-und-verzeichnislayout-unter-internal));
implementiert von `internal/hexagon/application/`, genutzt von
`internal/adapter/driving/`.

Schichtbeschreibung und Port-Bestand:
[`spec/architecture.md`](../../../../spec/architecture.md#23-hexagonportdriving).
Request-/Response-Typen und Sentinels stehen neben dem jeweiligen Interface;
ihre Zuordnung zu Exit-Codes beschreiben die
[Fehlermodelle](../../../../spec/architecture.md#7-fehlermodelle-und-resilienz).
Importregeln:
[`SPEC-013`](../../../../spec/spezifikation.md#spec-013--import-regel-tabelle-der-schichten).
