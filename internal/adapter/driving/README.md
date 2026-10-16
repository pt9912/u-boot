# internal/adapter/driving

Konkrete Driver — Einstiegspunkte aus der Außenwelt
([`SPEC-012`](../../../spec/spezifikation.md#spec-012--schichten-und-verzeichnislayout-unter-internal)).
Ein Cobra-Command je Subkommando; das Wiring liegt in `cmd/uboot/`.

Schichtbeschreibung, Exit-Code-Klassifikation und Importregeln:
[`spec/architecture.md`](../../../spec/architecture.md#25-adapterdriving),
[Fehlermodelle](../../../spec/architecture.md#7-fehlermodelle-und-resilienz) und
[`SPEC-013`](../../../spec/spezifikation.md#spec-013--import-regel-tabelle-der-schichten).
Welcher Befehl welche Anforderung erfüllt, folgt aus den Kennungen im Code
(`make doc-trace`).
