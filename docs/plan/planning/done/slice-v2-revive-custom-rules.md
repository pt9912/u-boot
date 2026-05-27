# Slice V2: revive Custom-Rules-Profil

> **Status:** Done
> **DoD:** Commit `98e5c9c`

## Auslöser

ADR-0003 (SOLID-nahes Lint-Profil) hatte als offenen Folgepunkt die
revive-Custom-Rules-Erweiterung. Bis zu dieser Sitzung lief revive
ohne expliziten `rules`-Block und nutzte damit die Default-Regeln aus
golangci-lint — implizit, nicht versioniert.

**Vorgezogen ohne den ursprünglich definierten Trigger** (wiederholte
Reviewer-Findings / neuer Style-Beschluss). In der M3-Anker-Triage-
Sitzung 2026-05-27 hat der User entschieden, den Carveout proaktiv
aufzuheben — die Code-Basis stabilisiert sich gerade auf MVP-Schwelle,
und ein expliziter Regel-Block macht zukünftige Lint-Profil-Änderungen
zur Policy-Entscheidung statt zum implizit-Mitwachsen.

## Aufhebung

`linters.settings.revive.rules` in `.golangci.yml` enthält jetzt
alle 24 Default-Regeln explizit aufgezählt plus eine projekt-
spezifische Erweiterung (`unused-receiver`). golangci-lint's Schema-
Verhalten (`rules:` ersetzt die Defaults vollständig) macht die
Enumeration zur Pflicht; sie ist in ADR-0006 ausführlich begründet.

## Geliefert

- **`docs/plan/adr/0006-revive-custom-rules.md`**: Kontext (defaults-
  vs-explicit-Mechanik), Entscheidung (24 Defaults + `unused-receiver`),
  Konsequenzen, Verworfenes.
- **`.golangci.yml`** `linters.settings.revive.rules` mit 25 Regeln
  in zwei klar getrennten Sektionen („revive default rule set
  (preserved)" + „project-specific extras (ADR-0006)").
- **`.golangci.yml`** zweiter Test-Exclude für `revive`
  (`^unused-receiver` in `_test.go`), gleicher Grund wie der
  bestehende `^unused-parameter`-Exclude (stateless Test-Fakes).
- **Refactoring-Beifang**: `resolveProjectName` in
  `internal/hexagon/application/initproject.go` von Methode auf
  `InitProjectService` zu Free-Function umgebaut (Service-Receiver
  wurde nicht referenziert; jetzt klar als stateless markiert).
- **Carveouts.md**: `ADR-0003 revive-Custom-Rules`-Zeile entfernt.
- **Roadmap**: Slice → Done, Phase „V2-vorgezogen".
- **READMEs**: Carveout-Count 8 → 7.
- **`make lint` grün** mit dem neuen Regel-Block; 6 initiale
  `unused-receiver`-Findings (5 in Test-Fakes, 1 in Production)
  durch den Test-Exclude + Refactoring aufgelöst.

## Out of Scope

- Weitere Custom-Rules wie `early-return`, `confusing-naming`,
  `cognitive-complexity`: ohne konkreten Trigger schwer zu
  rechtfertigen. Bei nächstem Bedarf in einer ADR-0006-Folgesektion
  ergänzen.
- Pro-Schicht-revive-Profile (z. B. striktere Regeln in `domain/`):
  Overengineering für den aktuellen Stand.
- Anpassungen an den `unused-parameter`-Test-Exclude: bleibt wie
  gehabt (M2b).

## Bezug

- Auslösende ADR: `0003-solid-nahes-lint-profil.md` Folgepunkte.
- Resultierende ADR: [`0006-revive-custom-rules.md`](../../adr/0006-revive-custom-rules.md).
- Aufhebung dokumentiert in: [`carveouts.md`](../in-progress/carveouts.md)
  (Zeile entfernt) und [`roadmap.md`](../in-progress/roadmap.md)
  (Carveout-Auflösungs-Slice-Tabelle).
- Hängt von: nichts; vorgezogen ohne Trigger nach User-Entscheidung.
