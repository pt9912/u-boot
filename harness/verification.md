# Verification Evidence

Diese Datei definiert, wie `u-boot` Slice-Ergebnisse verifiziert.
Verification beantwortet: Wurde der Slice richtig gebaut, gemessen an
Spec, ADR, DoD und ausgefuehrten Sensors?

Validation ist getrennt: Sie fragt, ob das Ergebnis den realen Nutzer-
oder Release-Bedarf trifft. Siehe [`harness/roles.md`](roles.md).

## Pflicht bei Slice-Closure

Jeder Slice, der nach `docs/plan/planning/done/` bewegt wird, braucht
eine sichtbare Verification-Evidence. Sie kann direkt im Slice stehen.
Nur bei sehr grosser Evidence entsteht ein eigenes Artefakt, das im
Slice verlinkt wird.

Minimum:

| Feld | Pflicht | Inhalt |
| --- | --- | --- |
| Scope | ja | Slice-ID, betroffene `LH-*`-/`ADR-*`-IDs, relevante Dateien |
| DoD-Abgleich | ja | Welche DoD-Punkte sind erfuellt, welche nicht |
| Sensors | ja | Ausgefuehrte Make-Targets oder Tests mit Ergebnis |
| Traceability | ja | Welche Tests/Gates belegen welche Spec- oder ADR-ID |
| Carveouts | ja | Neue, geloeste oder unveraenderte Carveouts |
| Nicht ausgefuehrt | ja | Ausgelassene Sensors mit Grund |
| Commit / Artefakt | wenn vorhanden | Commit-Hash, Release-Asset, Image-Tag oder generiertes Artefakt |

## Evidence Block

Standardblock fuer Slice-Closure:

```markdown
## Verification Evidence

Scope:
- Slice: `<slice-id>`
- IDs: `<LH-...>`, `<ADR-...>`
- Artefakte: `<Dateien/Pakete/Commands>`

DoD-Abgleich:
- [x] `<DoD-Punkt>` — Evidence: `<Test/Gate/Diff>`
- [ ] `<DoD-Punkt>` — offen: `<Grund/Folge-Slice/Carveout>`

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make test` | pass/fail/not run | `<kurzer Beleg>` |
| `make docs-check` | pass/fail/not run | `<kurzer Beleg>` |
| `make gates` | pass/fail/not run | `<kurzer Beleg>` |

Traceability:
| ID | Beleg |
| --- | --- |
| `<LH-...>` | `<Test/Gate/Doku>` |
| `<ADR-...>` | `<depguard/Test/Review>` |

Carveouts:
- Neu: `<none|Eintrag + Plan-Anker>`
- Geloest: `<none|Eintrag>`
- Unveraendert: `<none|Eintrag>`

Nicht ausgefuehrt:
- `<Sensor>` — `<Grund>`

Commit / Artefakt:
- `<hash|image-tag|release-asset|n/a>`
```

## Sensor-Auswahl

Der Verifier waehlt den engsten sinnvollen Sensor, muss aber begruenden,
wenn ein naheliegender Sensor nicht gelaufen ist.

| Aenderung | Mindest-Sensor | Normaler Closure-Sensor |
| --- | --- | --- |
| Nur Markdown/Doku | `make docs-check` | `make docs-check` plus gezielter Linkcheck |
| Go-Code ohne Docker-Runtime | `make test` | `make gates` |
| Architektur-/Importregeln | `make lint` | `make gates`; bei depguard-Aenderung zusaetzlich `make verify-depguard` |
| Coverage-relevanter Code | `make coverage-gate` | `make gates` |
| Docker-/Compose-/E2E-Verhalten | engster Docker-Test oder `make test-docker` | `make test-docker` plus `make gates` |
| Security-/Release-Pfad | `make govulncheck` oder `make image-scan` | `make ci` oder `make fullbuild` |

## Harte Regeln

- Ein gruenes `make gates` ersetzt nicht den DoD-Abgleich.
- Ein einzelner Unit-Test ersetzt nicht den Link auf die betroffene
  `LH-*`- oder `ADR-*`-ID.
- Nicht ausgefuehrte Sensors sind erlaubt, aber nur mit Grund.
- Neue temporaere Carveouts brauchen parallel einen Eintrag in
  [`docs/plan/planning/in-progress/carveouts.md`](../docs/plan/planning/in-progress/carveouts.md)
  und einen Plan-Anker.
- Reviewer-Findings sind keine Verification-Evidence. Sie koennen
  Evidence ausloesen, aber der Verifier prueft DoD, Spec und Sensors
  separat.
