# Slice V1: Plugin-System-Entscheidung (`LH-OPEN-003`)

## Auslöser

`LH-OPEN-003` (Plugin-System) ist in `spec/lastenheft.md` §14 offen:

> *„Es ist zu klären, ob Add-ons langfristig fest eingebaut oder als
> Plugins nachladbar sein sollen."*

`spec/architecture.md` §7 nennt das Plugin-System als „geplante
Erweiterung" (Driven-Port `PluginRegistry`). Beide Aussagen sind
prospektiv ohne konkrete Entscheidung.

## Aufhebungsbedingung

ADR-0007 (oder gleichwertig nummeriert) trifft eine Entscheidung
zwischen mindestens diesen Optionen:

1. **Statisch eingebaute Add-ons** — neue Services
   (Keycloak, OTel, …) werden im u-boot-Binary mitgeliefert; neue
   Add-ons brauchen u-boot-Release.
2. **Plugin-System über Driven-Port** — `PluginRegistry`-Port lädt
   externe Plugin-Binaries oder OCI-Bundles zur Laufzeit;
   Sicherheits- und Versionierungs-Modell ist Teil des ADR.
3. **Hybrid** — Kern-Add-ons statisch, exotische via Plugin.

## Akzeptanzkriterien

- ADR existiert mit Vergleich der Optionen, Trade-offs,
  Sicherheits-Implikationen (`LH-NFA-SEC-004` ist relevant) und
  konkreter Entscheidung.
- `LH-OPEN-003` in `spec/lastenheft.md` § 14 wird als entschieden
  markiert (analog `LH-OPEN-001` Go).
- `spec/architecture.md` §7 wird auf das ADR umgebogen.
- Eintrag in `carveouts.md` (`LH-OPEN-003`) wird entfernt.
- Falls die Entscheidung "Plugin-System" lautet, folgt ein eigener
  Implementierungs-Slice in `open/`.

## Bezug

- Auslösende Spec: `LH-OPEN-003`, `spec/architecture.md` §7.
- Inventar-Eintrag: [`carveouts.md`](../in-progress/carveouts.md) →
  `LH-OPEN-003`.
- Hängt von: erster Wunsch nach drittem Add-on oder externem Service,
  der nicht zum Kern gehört (vermutlich nach MVP-Closure).
- Phase: V1, weil das Add-on-System bis MVP-Closure statisch bleibt.
