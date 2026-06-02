# Review Harness

Diese Datei standardisiert Reviews fuer `u-boot`. Review beantwortet:
Welche Risiken, Vertragsbrueche oder Wartbarkeitsprobleme enthaelt ein
Diff, bevor Verification die DoD-Closure prueft?

Review ist eine Entscheidungsvorlage. Reviewer implementieren nicht,
verifizieren nicht vollstaendig gegen DoD und validieren nicht den
realen Nutzerbedarf.

## Input Context

Ein reproduzierbarer Review braucht mindestens:

| Quelle | Zweck |
| --- | --- |
| Diff oder PR-Patch | Was wurde geaendert |
| Aktiver Slice oder Plan | Scope, DoD, explizite Nicht-Ziele |
| Betroffene `LH-*`-IDs | Produktvertrag und Akzeptanzkriterien |
| Betroffene ADRs | Architektur- und Toolentscheidungen |
| [`AGENTS.md`](../AGENTS.md) | Hard Rules |
| [`harness/roles.md`](roles.md) | Rollenabgrenzung |
| [`harness/verification.md`](verification.md) | Abgrenzung zu Verification |
| Relevante Tests/Gates | Welche Sensoren koennen Findings bestaetigen |

Ohne Plan- oder Spec-Kontext ist ein Review nur Codekritik, kein
Harness-Review.

## Finding Categories

Findings werden absteigend sortiert: HIGH vor MEDIUM vor LOW vor INFO.

| Kategorie | Bedeutung | Blockiert |
| --- | --- | --- |
| HIGH | Produktvertrag, Sicherheit, Datenintegritaet, Architekturgrenze oder CI-Gate kann brechen | ja |
| MEDIUM | Wahrscheinliche Regression, fehlender Test, unklare Fehlerklassifikation oder relevantes Drift-Risiko | normalerweise vor Merge klaeren |
| LOW | Wartbarkeit, Doku-Praezision oder kleine Konsistenzluecke ohne unmittelbaren Vertragsbruch | nein |
| INFO | Beobachtung oder Kontext ohne Aenderungspflicht | nein |

### HIGH-Anker fuer `u-boot`

Ein Finding ist HIGH, wenn es eines dieser Muster trifft:

- `LH-*`-Akzeptanzkriterium oder Exit-Code-Vertrag kann verletzt werden.
- Hexagonale Importregel, Port-/Adapter-Grenze oder ADR wird gebrochen.
- User-Dateien koennen ohne Spec-konforme Bestaetigung, Backup oder
  managed-block-Schutz ueberschrieben werden.
- Docker-only-Harness, CI-Gate oder Security-Gate wird gelockert oder
  umgangen.
- Generierte Artefakte koennen kaputt oder inkonsistent werden
  (`compose.yaml`, `.env.example`, devcontainer, README/CHANGELOG).
- Neuer temporaerer Carveout entsteht ohne Inventar-Eintrag und
  Plan-Anker.

### MEDIUM-Anker fuer `u-boot`

Ein Finding ist MEDIUM, wenn es eines dieser Muster trifft:

- Neuer oeffentlicher CLI-Pfad hat keinen negativen Test oder keine
  Sentinel-/Exit-Code-Abdeckung.
- Fehlerbehandlung ist plausibel, aber nicht klar einer Kategorie
  `2/10/11/12/14` zugeordnet.
- Doku/README/CHANGELOG driftet gegen Spec oder implementiertes
  Verhalten.
- Tests belegen Verhalten, aber nicht die relevante `LH-*`-ID.
- Ein LOW-Muster wiederholt sich und wird zum Drift-Signal.

## Review Lenses

Reviewer pruefen diese Linsen explizit und berichten auch
Negativbefunde fuer relevante Linsen.

| Linse | Prueffrage |
| --- | --- |
| Spec / Traceability | Sind `LH-*`-/ADR-Anker korrekt, vollstaendig und testbar? |
| Exit Codes | Sind Fehlerpfade nach `LH-FA-CLI-006` klassifiziert und getestet? |
| Architecture | Bleiben hexagonale Grenzen, Ports, Adapter und Wiring sauber? |
| Docker-only / Gates | Bleibt der reproduzierbare `make`-/Docker-Harness intakt? |
| File Safety | Sind managed blocks, Backups, Two-Phase-Planung und Confirmers respektiert? |
| Tests | Gibt es Happy-, Boundary- und Negative-Pins fuer neue Vertraege? |
| Replay / Golden Sets | Sind Generator-Goldens fuer Fresh-State, Idempotenz und Safety-Pfade vorhanden? |
| Carveouts | Sind temporaere Ausnahmen dokumentiert und mit Plan-Anker versehen? |
| Docs / Release | Sind README, `docs/user/`, ADR-Index, Roadmap und CHANGELOG konsistent? |

## Output Schema

Jedes Finding verwendet dieses Schema:

```text
<CATEGORY> <path>:<line> - <kurzer Titel>
Quelle: <LH-*|ADR-*|Hard Rule|Maintainability>
Befund: <1-2 beobachtbare Saetze>
Risiko: <warum das relevant ist>
Verifizierbar: <Sensor/Test/Review-only>
```

Am Ende jedes Reviews:

```text
Geprueft ohne Befund:
- <Linse oder Pfad>

Nicht geprueft:
- <Linse oder Pfad> — <Grund>
```

## Non-Goals

- Keine Implementierungsvorschlaege als Ersatz fuer Findings.
- Keine Refactors ausserhalb des Diff-Scopes.
- Keine DoD-Closure. Das ist Verifier-Aufgabe.
- Keine Validation gegen Nutzerbedarf. Das ist Validator-Aufgabe.
- Keine Abwertung von Findings, weil ihre Behebung unbequem ist.

## Steering Loop

Wenn dasselbe Finding dreimal auftritt:

1. Klassifikation in dieser Datei schaerfen.
2. Pruefen, ob `AGENTS.md`, ADR oder Spec eine Hard Rule braucht.
3. Pruefen, ob ein computational Sensor moeglich ist
   (`make lint`, `make docs-check`, `make verify-depguard`, Tests).
4. Falls es nur temporaer toleriert wird: Carveout-Inventar und
   Plan-Anker aktualisieren.
