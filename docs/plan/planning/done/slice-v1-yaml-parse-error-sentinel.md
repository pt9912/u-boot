# Slice V1: `driven.ErrYAMLParse`-Sentinel für Parse-Error-Klassifikation

> **Status:** Done
> **DoD:** Implementation ✅ `1008326` (Sentinel + Adapter-Wrap + 4 Codec-Methods + Fake-Mirror + Generate-Code-10-Pfad + Doctor-Anti-Drift-Kommentar + 6 neue Tests)

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
Code **14** (technisch), obwohl [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006-exit-codes) Code **10**
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
   ([`slice-v1-plugin-system-entscheidung.md`](../done/slice-v1-plugin-system-entscheidung.md)),
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
- **Diskriminator-Regel im YAML-Adapter** (Antwort auf Review M1):
  yaml.v3 exportiert nur `*yaml.TypeError` als typisierten Fehler;
  Syntax-/Tokenizer-Fehler kommen als generischer `error` mit
  `yaml: …`-Präfix. Der Adapter klassifiziert deshalb wie folgt:
  - jeder Fehler, der aus einem yaml.v3-Aufruf (`yaml.Unmarshal`,
    `(*yaml.Decoder).Decode`, `node.Decode`) zurückkommt, gilt als
    Parse-Fehler und wird mit `ErrYAMLParse` gewrappt;
  - Fehler aus dem `bytes`/IO-Pfad davor (heute keine — content
    wird vom Caller übergeben) bleiben ungewrappt.
  Heuristik auf den `yaml:`-Präfix wird **nicht** verwendet, weil
  der Klassifikator damit yaml.v3-interne Strings stabil halten
  müsste. Type-Assertion auf `*yaml.TypeError` ist optional als
  zusätzlicher Sanity-Check, semantisch aber nicht erforderlich.
- **Sentinel-Geltungsbereich quer durch alle Codec-Methoden**
  (Antwort auf Review M2): jeder Methodenpfad, der content
  parsed, wrappt Parse-Fehler mit `ErrYAMLParse`. Konkret:
  - `Unmarshal` (heute `codec.go:36`),
  - `PatchScalar` content-Parse (heute `codec.go:61` —
    Fehlertext `parse yaml: %w` heute non-sentinel),
  - `PatchMappingEntryYAML` content-Parse,
  - `LocateMarkedEntry` content-Parse.
  Damit kann _jeder_ Application-Caller `errors.Is(err,
  driven.ErrYAMLParse)` einheitlich prüfen, statt sich pro
  Codec-Methode zu fragen, ob die Klassifikation greift.
  `valueYAML`-Parse-Fehler in `PatchMappingEntryYAML` bleiben
  bewusst `ErrYAMLFragmentInvalid` (die kommen vom Caller-
  konstruierten Fragment, nicht von User-Content) — der
  Geltungsbereich gilt für _content_, nicht für strukturelle
  Fragment-Validierung.
- YAML-Adapter `internal/adapter/driven/yaml/codec.go` wendet
  obige Regel an. Read-Fehler bleiben unverändert (kein doppeltes
  Wrappen).
- **Repair-Hint-Normalisierung im Adapter** (Antwort auf Review N5):
  beim Wrap mit `ErrYAMLParse` strippt der Adapter einen
  führenden `yaml: `-Präfix aus der yaml.v3-Fehlermeldung, damit
  Caller-seitiges `%v` keine doppelten `yaml: yaml: …`-Hints
  produziert. Helper `stripYAMLPrefix(err) string` privat im
  Adapter-Package, in `codec_test.go` gepinnt.
- Fake-YAML-Codec in `internal/hexagon/application/fakes_test.go`
  spiegelt den Produktions-Vertrag _automatisch_ (Antwort auf
  Review N2 — Option a): `fakeYAML.Unmarshal` /
  `fakeYAML.PatchScalar` etc. wrappen yaml.v3-Parse-Fehler mit
  `ErrYAMLParse` genau wie der Produktions-Adapter, ohne dass der
  Test-Code einen Konfig-Hook setzen muss. Damit testet jeder
  Generate-Test, der einen kaputten YAML-Bytestream einseedet,
  automatisch den Sentinel-Pfad. Anti-Boilerplate.

### M7-T5-Callsite-Migration

- `collectDevcontainerForwardPorts` (`generate.go:676-695`)
  klassifiziert **zwischen** Helper-Return und dem heute
  bestehenden `ErrGenerateFileSystem`-Wrap (Antwort auf Review
  M3 — der heutige `%v`-Tail in `generate.go:691` würde die
  Sentinel-Chain sonst schlucken):

  ```go
  ports, err := collectActiveServicePorts(s.fs, s.yaml, baseDir, services)
  if err != nil {
      if errors.Is(err, driven.ErrYAMLParse) {
          return nil, fmt.Errorf(
              "%w: compose.yaml is unparseable (%v); repair the YAML manually",
              driving.ErrGenerateManualConflict, err)
      }
      return nil, fmt.Errorf("%w: collectActiveServicePorts: %v",
          driving.ErrGenerateFileSystem, err)
  }
  ```

  - Parse-Fehler ⇒ `driving.ErrGenerateManualConflict`
    → Exit-Code 10.
  - Andere Fehler bleiben `ErrGenerateFileSystem` → Exit-Code 14.
- **Sentinel-Erweiterung statt neuer Sentinel** (Antwort auf
  Review N1): `ErrGenerateManualConflict` wird hier bewusst über
  managed-block-Konflikte hinaus auf „User muss YAML manuell
  reparieren" erweitert. Beide Fälle teilen die Semantik
  „[`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006-exit-codes) Code 10: fachlich, manuelles Eingreifen", und
  der Exit-Code ist gleich — eine separate
  `ErrGenerateInvalidCompose`-Sentinel würde ohne aktuellen
  Mehrwert die Sentinel-Tabelle aufblähen. Diese Erweiterung
  wird im Doc-Kommentar von `ErrGenerateManualConflict`
  (`driving/errors.go`) festgehalten, damit der nächste Reviewer
  sieht, dass „Manual Conflict" weiter gefasst ist als
  managed-block.
- N2-Doc-Kommentar in `generate.go:663-675` wird **entfernt**
  (Antwort auf Review N4 — Anti-Drift sitzt im Test, nicht im
  Kommentar; „Resolved in <hash>"-Cruft im Code wird vermieden).

### Co-Touch in `doctor.go`

- `collectActiveServicePorts` (`doctor.go:899-929`) reicht
  `yamlCodec.Unmarshal`-Fehler heute nackt durch — Sentinel-
  Chain bleibt damit intakt (Antwort auf Review M4). Anti-Drift-
  Pin: ein Kommentar direkt über dem `yamlCodec.Unmarshal`-Aufruf
  hält fest, dass der Return-Pfad bewusst nicht gewrappt wird,
  damit `driven.ErrYAMLParse` für die Generate-Klassifikation
  via `errors.Is` greifbar bleibt. Wer später einen
  „read compose.yaml: %v"-Wrap einzieht, bricht den Code-10-
  Pfad in `generate devcontainer` und muss `%w` benutzen.
- **Doctor-seitig keine Verhaltensänderung**: Doctor produziert
  Diagnostics, keine Exit-Codes — die Sentinel-Klassifikation
  ändert dort nichts (Out of Scope V1+, vgl. unten).

### Test-Pin

- Neuer Test `TestGenerateDevcontainer_CorruptComposeYAML_Code10`
  in `generate_test.go`: seed eine syntaktisch kaputte
  `compose.yaml` (z. B. unmatched bracket) **in den FS-Fake**,
  rufe `Generate(ArtifactDevcontainer)`; assert
  `errors.Is(err, driving.ErrGenerateManualConflict)` und
  `cli.ExitCode(err) == 10`. Test nutzt den realen
  `fakeYAML`-Codec (der laut N2-Regel automatisch
  `ErrYAMLParse` wrappt) — kein expliziter Sentinel-Inject im
  Test-Setup. Anti-Drift-Pin gegen ein versehentliches
  Zurück-Drift auf Code 14.
- Zusätzlicher Test im Adapter:
  `TestCodec_Unmarshal_WrapsParseError_AsErrYAMLParse` in
  `internal/adapter/driven/yaml/codec_test.go`. Sichert die
  Adapter-Regel direkt, unabhängig vom Application-Test-Pfad.

### Carveout-Cleanup

- Eintrag in [`carveouts.md`](../in-progress/carveouts.md) wird
  entfernt.
- Slice-Datei wird nach `done/` verschoben; **die DoD-Line in
  diesem File trägt direkt den Implementations-Commit-Hash**
  (kein `git log --grep`-Platzhalter — Memory
  `feedback_done_slice_dod_hash`, Antwort auf Review N3).
- Roadmap-Eintrag dieses Slices wird auf Done gesetzt und
  referenziert denselben Commit-Hash.

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
- Spec: [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006-exit-codes) (Exit-Code-Partitionierung 10 vs 14).
- Architektur: `spec/architecture.md` §4 (`depguard.application-
  no-yaml`-Regel, die die heutige inline-Lösung blockt).
- Hängt von: nichts. Sentinel-Einführung ist eine reine
  Port-Erweiterung. Co-Touch in `doctor.go` (M4-Code) ist nur
  ein dokumentarischer Anti-Drift-Kommentar, keine
  Verhaltensänderung — Doctor produziert Diagnostics, keine
  Exit-Codes, und der `collectActiveServicePorts`-Helper reicht
  Sentinel-Errors heute schon nackt durch.
- Trigger-Beziehung: könnte als Folie eines V1-Plugin-Slices
  laufen, wenn der zuerst landet — dann verbraucht der Plugin-
  Slice den Sentinel mit und dieser Slice wird auf done/
  geschoben mit Verweis.
- Phase: V1 (post-MVP, nicht release-blocking).
