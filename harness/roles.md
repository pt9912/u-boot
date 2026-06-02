# Agent Roles

Diese Datei trennt Rollen fuer AI-gestuetzte Arbeit an `u-boot`.
Rollen sind Kontextgrenzen, keine Personen. Eine Person oder ein Agent
kann mehrere Rollen nacheinander ausfuehren, aber nicht mit demselben
Eingabe-Kontext und nicht ohne Uebergabe-Artefakt.

## Slice Sequence

```text
Planner -> Architect -> Implementation -> Reviewer -> Verifier -> Validator -> Planner
```

Jeder Rollenwechsel braucht ein sichtbares Artefakt: Plan, ADR-Bezug,
Diff, Findings, Verification-Evidence, Validation-Evidence oder
Closure-Notiz. Ohne Artefakt ist es kein Rollenwechsel, sondern nur ein
Kontextwechsel ohne Pruefbarkeit.

## Role Contracts

| Rolle | Primaere Frage | Eingabe-Kontext | Output |
| --- | --- | --- | --- |
| Planner | Was wird als naechstes klein genug geliefert? | Roadmap, aktiver Slice, Spec-IDs, Carveouts | Slice-/Tranche-Plan, Lifecycle-Entscheidung, Closure-Notiz |
| Architect | Passt die Loesung zu ADRs und Architektur? | Spec, Architektur, ADRs, Slice-Plan | bestaetigte ADR-Bezuege, Folge-ADR-Vorschlag oder Architektur-Finding |
| Implementation | Wie wird der Slice minimal und korrekt umgesetzt? | aktiver Slice, relevante Spec/ADR, `AGENTS.md`, engste Codepfade | Code-/Doku-Diff, lokale Sensor-Evidence, offene Risiken |
| Reviewer | Welche Risiken oder Vertragsbrueche enthaelt der Diff? | Plan, ADRs, Spec-Anker, Diff, relevante Tests | Findings mit HIGH/MEDIUM/LOW/INFO und Datei-/Zeilenbezug |
| Verifier | Erfuellt das Ergebnis DoD, Spec und Gates? | Slice-DoD, Diff, Tests, Make-Target-Ausgaben, Traceability | Verification-Evidence, fehlende Sensoren, DoD-Abweichungen |
| Validator | Loest das Ergebnis den realen Nutzer-/Release-Bedarf? | Nutzerpfad, README/Quickstart, Release-Ziel, generierte Artefakte | Validation-Evidence oder Rueckgabe an Planner |

## Boundaries

### Planner

Der Planner bewegt Planning-Artefakte durch
`open/ -> next/ -> in-progress/ -> done/`, schneidet zu grosse Arbeit
in Tranchen und pflegt Roadmap sowie Carveout-Inventar.

Der Planner implementiert nicht und stuft Review-Findings nicht ohne
Architektur- oder Verification-Artefakt herunter.

### Architect

Der Architect prueft ADR-Konformitaet, hexagonale Schichtgrenzen,
Port-/Adapter-Schnittstellen, depguard-Vertraege, Gate-Politik und
Release-/Distributionsentscheidungen.

Der Architect kann eine Folge-ADR verlangen oder vorschlagen. Accepted
ADRs werden nicht stillschweigend ueberschrieben.

### Implementation

Implementation setzt nur den aktiven Scope um. Sie nutzt die engsten
sinnvollen Sensoren frueh und `make gates` als normalen Code-Handoff,
wenn Docker verfuegbar ist.

Implementation darf ADR- oder Spec-Konflikte nicht pragmatisch
uebergehen. Bei Konflikt entsteht ein Uebergabe-Artefakt an Architect
oder Planner.

### Reviewer

Reviewer priorisiert Bugs, Vertragsbrueche, Architekturdrift,
Exit-Code-Fehler, managed-block-Sicherheitsrisiken, fehlende Tests und
Carveout-Verstoesse.

Finding-Kategorien:

| Kategorie | Bedeutung |
| --- | --- |
| HIGH | Kann Produktvertrag, Sicherheit, Datenintegritaet, Architekturgrenze oder CI-Gate brechen |
| MEDIUM | Wahrscheinliche Regression, fehlender Test, unklare Fehlerklassifikation oder Drift-Risiko |
| LOW | Wartbarkeit, Doku-Praezision, kleine Konsistenzluecke |
| INFO | Beobachtung ohne Aenderungspflicht |

Reviewer verifiziert nicht die komplette DoD-Closure. Das ist Aufgabe
des Verifier.

### Verifier

Verifier prueft "built the thing right": DoD gegen Diff, Tests gegen
Spec-ID, Make-Targets gegen Handoff, Docs gegen Links, Carveouts gegen
Inventar.

Verifier darf nicht behaupten, ein Gate sei gruen, wenn es nicht
ausgefuehrt wurde. Nicht ausgefuehrte Sensoren werden mit Grund
gelistet.

### Validator

Validator prueft "built the right thing": Passt das Ergebnis zum
Nutzerpfad, zur Release-Kommunikation, zu README/Quickstart und zu den
generierten Artefakten?

Fuer `u-boot` sind typische Validation-Fragen:

- Versteht ein neuer Nutzer den CLI-Flow aus README und `--help`?
- Erzeugen `init`, `add`, `generate`, `up`, `down` oder `logs` die
  erwartbaren Projektartefakte?
- Passt die Release-Story zu CHANGELOG, README und Spec-Prioritaeten?

## Conflict Handling

| Konflikt | Entscheidungspfad |
| --- | --- |
| Reviewer meldet ADR-Verstoss, Implementation widerspricht | Architect prueft ADR-Aktualitaet und Code. Output: bestaetigtes Finding oder Folge-ADR. |
| Verifier findet DoD-Luecke, Implementation haelt Scope fuer erledigt | Planner entscheidet: Slice zurueck, DoD anpassen oder Carveout mit Plan-Anker. |
| Validation ist rot, Verification ist gruen | Planner entscheidet, ob der Slice fachlich falsch geschnitten war oder ein Folge-Slice reicht. |
| Gate-Lockerung waere bequem | Architect/Planner brauchen ADR- oder Carveout-Anker; Implementation aendert nicht still. |

## Minimal Handoff Fields

Jede Rolle uebergibt knapp:

```text
Role:
Input context:
Changed artefacts:
Evidence:
Open risks:
Next role:
```

Diese Felder koennen als kurzer Abschnitt in Slice-Closure, PR-Text,
Review-Kommentar oder finalem Agent-Handoff stehen.
