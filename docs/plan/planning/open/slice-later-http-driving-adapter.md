# Slice Later: HTTP-Driving-Adapter (Daemon-Variante)

## Auslöser

`spec/architecture.md` §7 nennt als „geplante Erweiterung":

> *„HTTP-Driving-Adapter, falls u-boot perspektivisch eine
> Daemon-Variante bekommen soll."*

Keine Spec-Anforderung, kein konkreter Use-Case heute — aber im
Architektur-Dokument als Roadmap-Andeutung verankert. Damit ist es
nach `LH-FA-PROJDOCS-005` ein temporärer Carveout (prospektive
Doku-Phrase) und braucht einen Slice-Plan.

## Aufhebungsbedingung

Mindestens einer der folgenden Trigger:

1. Konkreter Use-Case für u-boot als langlaufender Daemon (z. B.
   Multi-Projekt-Orchestrierung, Web-Dashboard, IDE-Integration
   über HTTP-API).
2. Externer Bedarf an einer Maschinen-Schnittstelle, die über die
   JSON-CLI-Ausgabe (`LH-NFA-USE-004`) hinausgeht.

In dem Fall: ADR mit Begründung + neuer Driving-Adapter unter
`internal/adapter/driving/http/`. Layer-Regeln in
`spec/architecture.md` §3 müssen den neuen Adapter nicht
erweitern — er fällt unter die bestehende `adapter/driving`-Kategorie.

## Akzeptanzkriterien

- ADR existiert mit konkreter Use-Case-Beschreibung und
  Implementierungs-Plan.
- Falls implementiert: `internal/adapter/driving/http/` mit
  HTTP-Server-Wireup, der die Driving-Ports per HTTP exponiert.
- Falls nur entschieden „wird nicht gebaut": entsprechender
  Status-Vermerk + Entfernung aus `spec/architecture.md` §7.
- Eintrag in `carveouts.md` (HTTP-Driving-Adapter) wird entfernt.

## Bezug

- Auslösende Doku: `spec/architecture.md` §7.
- Inventar-Eintrag: [`carveouts.md`](../in-progress/carveouts.md) →
  HTTP-Driving-Adapter.
- Hängt von: konkretem externen Trigger (heute keiner).
- Phase: Later — bewusst niedrigste Priorität, keine MVP- oder
  V1-Pflicht.
