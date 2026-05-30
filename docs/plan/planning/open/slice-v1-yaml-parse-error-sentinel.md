# Slice V1: `driven.ErrYAMLParse`-Sentinel für Parse-Error-Klassifikation

## Auslöser

M7-T5 (Code-Review-Followup N2, Commits `27de9c5` + `51d8f6d`)
identifizierte eine Klassifikations-Lücke in
`collectDevcontainerForwardPorts`
(`internal/hexagon/application/generate.go`):

`collectActiveServicePorts` (Doctor-Helper) macht zwei Schritte auf
`compose.yaml`:

1. `fs.ReadFile(...)` — technische Fehler (Disk weg, Permission,
   File vanished).
2. `yamlCodec.Unmarshal(...)` — fachliche Fehler (User hat
   ungültiges YAML geschrieben).

Beide bubbeln als generic `error` durch. Der Application-Code kann
sie nicht unterscheiden, weil die
`depguard.application-no-yaml`-Regel
(`spec/architecture.md` §4 / `.golangci.yml`) den Import von
`gopkg.in/yaml.v3` in der Application-Schicht verbietet. Damit
landet ein Parse-Fehler heute als `ErrGenerateFileSystem` → Exit-
Code **14** (technisch), obwohl `LH-FA-CLI-006` Code **10**
(fachlich) verlangt.

Die Lücke ist real-existierend in M7-T5, aber strukturell breiter:
auch M4 doctor, M5 detectServiceState und M6 compose-Reads parsen
YAML, ohne IO- von Parse-Fehlern unterscheiden zu können — sie
fallen heute nur deshalb nicht auf, weil sie ihre Wrap-Sentinels
inline definieren (M5 `ErrProjectNotInitialized` für u-boot.yaml-
Parse-Fehler, M4 produziert Diagnostics statt Exit-Codes). M7-T5
ist der erste Caller, dessen Spec-Mandat eine harte Trennung
verlangt und der das nicht inline lösen kann.

## Aufhebungsbedingung

Mindestens einer von zwei Triggern feuert:

1. **Realer User-Bug-Report**: ein User meldet, dass `u-boot
   generate devcontainer` auf einer kaputten `compose.yaml` mit
   Exit 14 statt 10 endet und die Diagnose-Erwartung damit reißt.
2. **Nächste Stelle, die den `YAMLCodec`-Port aus anderem Grund
   anfasst** — z. B. der V1-Plugin-Slice
   ([`slice-v1-plugin-system-entscheidung.md`](slice-v1-plugin-system-entscheidung.md)),
   wenn er Plugin-Manifeste als YAML lädt und dort eigene
   Klassifikation braucht. Dieser Slice trägt den Sentinel dann
   als Side-Effekt mit und migriert die M7-T5-Stelle im selben PR.

Bis dahin bleibt die Lücke explizit im Top-Kommentar von
`collectDevcontainerForwardPorts` dokumentiert (Stand `27de9c5`).

## Akzeptanzkriterien

### Sentinel-Einführung

- Neuer Sentinel `driven.ErrYAMLParse` in
  `internal/hexagon/port/driven/yamlcodec.go`. Doc-Kommentar
  erklärt: signalisiert Parse-Fehler (nicht IO-Fehler) und ist
  von Application-Callern via `errors.Is` abgreifbar, ohne den
  YAML-Adapter zu importieren.
- YAML-Adapter `internal/adapter/driven/yaml/codec.go` wrappt
  yaml.v3-Parse-Fehler (TypeError, SyntaxError) mit
  `ErrYAMLParse`. Read-Fehler bleiben unverändert (kein
  doppeltes Wrappen).
- Fake-YAML-Codec in `internal/hexagon/application/fakes_test.go`
  spiegelt den Vertrag: ein Test, der einen Parse-Fehler
  simulieren will, kann den Fake so konfigurieren, dass er
  `ErrYAMLParse`-wrapped zurückgibt.

### M7-T5-Callsite-Migration

- `collectDevcontainerForwardPorts` (`generate.go`) prüft
  `errors.Is(err, driven.ErrYAMLParse)` nach dem
  `collectActiveServicePorts`-Aufruf:
  - Parse-Fehler ⇒ Wrap in `driving.ErrGenerateManualConflict`
    mit Repair-Hint („compose.yaml is unparseable: <detail>;
    repair the YAML manually") → Exit-Code 10.
  - Andere Fehler bleiben `ErrGenerateFileSystem` → Exit-Code 14.
- N2-Doc-Kommentar in `generate.go` wird auf
  `Resolved in <slice-hash>` aktualisiert oder entfernt.

### Test-Pin

- Neuer Test `TestGenerateDevcontainer_CorruptComposeYAML_Code10`
  in `generate_test.go`: seed eine syntaktisch kaputte
  `compose.yaml` (z. B. unmatched bracket); assert
  `errors.Is(err, driving.ErrGenerateManualConflict)` und
  `cli.ExitCode(err) == 10`. Anti-Drift-Pin gegen ein
  versehentliches Zurück-Drift auf Code 14.

### Carveout-Cleanup

- Eintrag in [`carveouts.md`](../in-progress/carveouts.md) wird
  entfernt.
- Roadmap-Eintrag dieses Slices wird auf Done gesetzt.

## Out of Scope (V2+)

- **Migration anderer Caller** (M4 doctor, M5 detectServiceState,
  M6 compose-Reads) auf den Sentinel. Jede Stelle hat ihre eigene
  Klassifikations-Semantik — M5 mappt u-boot.yaml-Parse-Fehler
  bewusst auf `ErrProjectNotInitialized` (Projekt-State-Kohärenz,
  nicht reine YAML-Pathologie), M4 produziert Doctor-Diagnostics,
  M6 hat noch keine Spec-Mandate für die Trennung. Eine
  opt-in-Evolution ist sauberer als ein Big-Bang-Refactor.
- **Feinere Parse-Kategorien** (TypeError vs SyntaxError vs
  unknown-key). M7-T5 braucht nur die binäre Parse-vs-nicht-Parse-
  Unterscheidung; tieferes Kategorisieren wartet auf konkreten
  User-Bedarf.
- **Andere Driven-Adapter** (z. B. JSONC im Devcontainer-Pfad).
  `stripJSONC` + `encoding/json` ist heute nur Read-only-Pfad und
  hat noch keine Parse-Fehler-Klassifikations-Pflicht; folgt
  separat, wenn ein analoger Caller auftaucht.

## Bezug

- Auslösender Review: M7 Post-Merge-Review-Followup N2
  (`27de9c5` Fix-Commit, `51d8f6d` Slice-Plan-Eintrag).
- Spec: `LH-FA-CLI-006` (Exit-Code-Partitionierung 10 vs 14).
- Architektur: `spec/architecture.md` §4 (`depguard.application-
  no-yaml`-Regel, die die heutige inline-Lösung blockt).
- Hängt von: nichts; Sentinel-Einführung ist eine reine
  Port-Erweiterung.
- Trigger-Beziehung: könnte als Folie eines V1-Plugin-Slices
  laufen, wenn der zuerst landet — dann verbraucht der Plugin-
  Slice den Sentinel mit und dieser Slice wird auf done/
  geschoben mit Verweis.
- Phase: V1 (post-MVP, nicht release-blocking).
