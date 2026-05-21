# Slice V2: revive Custom-Rules-Profil

## Auslöser

ADR-0003 (SOLID-nahes Lint-Profil) hat als offenen Folgepunkt:

> *„Erweiterung um `revive`-Custom-Rules in einem Folge-ADR, falls
> die default-Konfiguration zu schwach wird."*

Heute läuft `revive` im Default-Profil (mit einem
`unused-parameter`-Excludes für Tests). Das ist akzeptabel, solange
das Default-Profil keine konkreten Pain-Points hat. Sobald
Pattern-Drift oder Stil-Findings dokumentiert sind, wird ein
projekt-eigenes `revive`-Profil fällig — V2-Material (post-MVP).

## Aufhebungsbedingung

Mindestens einer der folgenden Trigger:

1. Wiederholte `revive`-Findings, die das Default-Profil nicht
   greift, aber im Review aufschlagen (z. B. Naming-Konventionen,
   Receiver-Naming, Cyclomatic-Drift jenseits cyclop).
2. Neuer Code-Style-Beschluss, der eine spezifische `revive`-Regel
   nötig macht (z. B. `package-comments`-Pflicht für alle Pakete).

In dem Fall: ADR-0006 mit konkreten Rules, plus Aktivierung im
`.golangci.yml` unter `linters.settings.revive.rules`.

## Akzeptanzkriterien

- `docs/plan/adr/0006-revive-custom-rules.md` (oder gleichwertig
  nummeriert) existiert.
- `.golangci.yml` enthält den `revive.rules`-Block mit konkreten
  Custom-Rules + `Why:`-Kommentar pro Regel.
- `make lint` läuft grün auf dem aktuellen Code-Stand.
- Eintrag in `carveouts.md` (`ADR-0003 Folgepunkt: revive-Custom-Rules`)
  wird entfernt oder als gelöst markiert.

## Bezug

- Auslösende ADR: `0003-solid-nahes-lint-profil.md` Folgepunkte.
- Inventar-Eintrag: [`carveouts.md`](../in-progress/carveouts.md) →
  ADR-0003 revive-Custom-Rules.
- Hängt von: keinem konkreten anderen Slice — kann jederzeit
  angegangen werden, wenn der Trigger eintritt.
- Phase: V2, weil heute kein akuter Bedarf.
