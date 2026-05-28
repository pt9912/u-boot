# Slice M5: `u-boot add postgres`-Flow

> **Status:** Done
> **DoD:** T1 ✅ `995726a` / T2 ✅ `f054986` / T3 ✅ `4091ac9` / T4a ✅ `6052b71` / T4b ✅ `cca1254` / T4c ✅ `529be1c` / T5 ✅ `23209cd` / T6 ✅ `7e7bbcd` / T7 ✅ `075667d`

## Auslöser

Nach M3 (`u-boot init`) und M4 (`u-boot doctor`) ist das dritte MVP-
Subkommando dran: `u-boot add <service>`, mit PostgreSQL als erstem
konkreten Add-on.

Spec-Pflicht für M5 (alle MVP-Priorität):

- **`LH-FA-ADD-001`** Befehlsstruktur `u-boot add <service>`, nur in
  initialisiertem Projekt (`u-boot.yaml` vorhanden).
- **`LH-FA-ADD-002`** PostgreSQL hinzufügen: Compose-Service +
  Volume + `.env.example`-Einträge + Port + Healthcheck.
- **`LH-FA-ADD-005`** Doppel-Add-Verhinderung über die
  `services.<name>.enabled`-State-Machine in `u-boot.yaml`.

Out of Scope (V1):

- **`LH-FA-ADD-003`** Keycloak (V1).
- **`LH-FA-ADD-004`** OTel (V1).
- **`LH-FA-ADD-006`** Add-on-Abhängigkeiten + `--with-deps` (V1).
- **`LH-FA-ADD-007`** `u-boot remove <service>` (V1).

## State-Machine (LH-FA-ADD-005)

Pro Service-Name gibt es sechs beobachtbare Zustände beim Add-Versuch.
Malformed Managed-Blocks für den Ziel-Service sind ein Pre-Classification-
Abort: sobald `compose.yaml` einen kaputten `service.<name>`-Block enthält,
bricht `add` mit `ErrServiceInconsistent` ab, auch wenn der YAML-Eintrag
deaktiviert oder `enabled` unset ist. Die Tabelle beschreibt nur die
wohlgeformt-vorhanden/fehlend-Fälle.

| Zustand                | `services.<name>` in u-boot.yaml | `enabled` | Managed-Block in compose.yaml | Add-Aktion                                                                              |
| ---------------------- | -------------------------------- | --------- | ----------------------------- | --------------------------------------------------------------------------------------- |
| **unregistered**       | fehlt                            | —         | fehlt                         | Neu anlegen: services-Eintrag + Compose-Block + .env.example-Block (LH-FA-ADD-002).     |
| **active**             | vorhanden                        | `true`    | vorhanden                     | No-op, wenn PostgreSQL-Artefakte vollständig sind; sonst Artefakt-Repair ohne Service-Duplikat. |
| **deactivated**        | vorhanden                        | `false`   | vorhanden/fehlend, aber wohlgeformt | Re-Aktivierung: `enabled: true` + Compose-Block + .env-Block neu erzeugen.              |
| **enabled-key-fehlt**  | vorhanden                        | (unset)   | vorhanden/fehlend, aber wohlgeformt | (Doctor-warn-Pfad; Add interpretiert als deactivated und re-aktiviert wie oben.)        |
| **inconsistent-yaml**  | fehlt                            | —         | vorhanden (managed)           | Abort: Compose-Block ohne YAML-Anker. ErrServiceInconsistent → Code 10 + Repair-Hint.   |
| **inconsistent-block** | vorhanden                        | `true`    | fehlt                         | Compose-Block neu erzeugen (deterministisch); kein Abort.                               |

## Tranchen-Schnitt

1. **T1 — u-boot.yaml services-Schema + Domain-Types.**
   - `ubootYAMLConfig` um `Services map[string]ubootYAMLService` mit
     `Enabled *bool` (Pointer, um „unset" von `false` zu
     unterscheiden) erweitern. `omitempty`-Marshal-Tag, damit
     `u-boot init` (frischer Projektstart) keinen leeren
     `services:`-Block schreibt.
   - Domain-Type `ServiceName` (analog `ProjectName`, eigene
     Validierungs-Regex — Service-Namen müssen YAML-key-fähig +
     Compose-name-fähig sein).
   - Domain-Type `ServiceState` mit den 6 oben tabellierten Werten.
   - Tests: marshal/unmarshal-roundtrip mit + ohne services-Block,
     `Enabled`-Pointer-Semantik (nil vs &false vs &true),
     ServiceName-Validierung.

2. **T2 — Driving-Port `AddServiceUseCase` + Sentinels.**
   - `AddServiceRequest` (BaseDir, ServiceName domain.ServiceName).
     Persistent-Mode-Flags bleiben bewusst draußen: `add postgres` ist
     MVP nicht-interaktiv; spätere Add-ons bekommen eigene Request-
     Felder statt Init-Flags zu erben.
   - `AddServiceResponse` (ServiceName, PriorState, State,
     geänderte Pfade). Vollständig aktiver Zustand ist ein
     nil-error No-op: `PriorState=Active`, `State=Active`,
     `Changed=nil`. T4 erweitert den Vertrag für aktive, aber
     unvollständige PostgreSQL-Artefakte: `PriorState=Active`,
     `State=Active`, `Changed!=nil` nach deterministischem
     Volume-/.env-Repair.
   - Sentinels:
     - `ErrServiceUnsupported` (Service-Name nicht im
       built-in-Katalog, heute nur „postgres").
     - `ErrServiceInconsistent` (LH-FA-ADD-005-inconsistent-yaml-Fall
       sowie malformed managed compose-block).
     - `ErrProjectNotInitialized` (kein `u-boot.yaml` → LH-FA-ADD-001).
     Mapping zu LH-FA-CLI-006-Exit-Code 10 (validation); kein eigener
     Code 13 für Add-Projektzustand in M5.

3. **T3 — Application-Service-Skeleton + State-Detection.**

   **Scope-Schnitt T3↔T4:** T3 liefert Skeleton + `detectServiceState`
   + Top-Level-Dispatch + Plan-Datenstruktur. T4 liefert die
   Execute-Bodies (Template-Render + File-Writes). So bleibt T3
   testbar ohne Template-Assets; T4 ist reine Render+Write-Tranche.
   `executeAdd` ist in T3 ein Stub, der `errors.New("addservice
   execute: not yet implemented (M5-T4)")` zurückgibt — kein eigener
   Carveout-Eintrag nötig, weil T4 in derselben Slice-Datei als
   nächste Tranche dokumentiert ist (Plan-Pflicht aus
   `LH-FA-PROJDOCS-005` erfüllt) und der CLI-Subcommand erst T6 ist,
   d. h. der Stub ist bis dahin nur über direkte Service-Aufrufe in
   Tests erreichbar.

   **Datei-Layout:**
   - `internal/hexagon/application/addservice.go` — `AddServiceService`
     analog `InitProjectService`, DI nur über die für add nötigen
     Ports: `fs driven.FileSystem`, `yaml driven.YAMLCodec`,
     `logger driven.Logger` (nil-tolerant via lokalem `noopLogger`).
     Kein `Confirmer`/`ProgressPort` — add ist MVP nicht-interaktiv
     und produziert nur den Summary-Bericht aus T6.
   - `internal/hexagon/application/addservice_test.go` — Tests in
     `_test`-Package (testpackage-Konvention, siehe carveouts).
   - `internal/hexagon/application/export_test.go` — schmale
     Test-Brücke für unexported Helpers (`detectServiceState`,
     `planAdd`, `servicePlan`/`serviceAction`-Projektion), damit die
     Tests im externen `application_test`-Package bleiben können,
     ohne Implementierungsdetails produktiv zu exportieren.
   - Marker-Name-Helfer `serviceMarkerName(svc domain.ServiceName)
     string` (= `"service." + svc.String()`) lokal in
     `addservice.go`; bewusst **nicht** in `managedblock`, weil
     `managedblock` keine `domain`-Imports hat (Schicht-Hygiene).

   **YAMLCodec-Port-Erweiterung (Variante B):**
   - `internal/hexagon/port/driven/yamlcodec.go` erweitert
     `YAMLCodec` um:
     `PatchScalar(content []byte, path []string, value any) ([]byte, error)`.
   - `internal/adapter/driven/yaml/codec.go` implementiert den Patch
     mit der `yaml.Node`-API von `gopkg.in/yaml.v3`: Mapping-Pfad
     gehen, fehlende Mapping-Knoten anlegen, Scalar am Ziel setzen,
     danach re-marshaln. Kommentare und unbekannte Felder bleiben so
     weit erhalten, wie `yaml.v3` sie am Node führt.
   - `internal/hexagon/application/fakes_test.go` zieht die Methode
     im Fake mit. T3 testet mindestens Update eines vorhandenen
     Scalars und Insert eines fehlenden `services.postgres.enabled`.
   - Die Application-Schicht importiert **nicht** `gopkg.in/yaml.v3`;
     YAML-Bibliotheksdetails bleiben im Driven-Adapter
     (LH-FA-ARCH-003).

   **`detectServiceState(baseDir, name)`-Algorithmus:**
   1. `s.fs.Exists(baseDir/u-boot.yaml)` → `(false, nil)` ⇒
      `ErrProjectNotInitialized`. `Exists`-Fehler werden als
      technischer Wrap ohne Add-Sentinel zurückgegeben; ein
      Permission-/I/O-Problem darf nicht als Projektzustand
      fehlklassifiziert werden.
   2. `s.fs.ReadFile + s.yaml.Unmarshal` in die existierende
      `ubootYAMLConfig`-Struktur (T1-Schema). `ReadFile`-Fehler
      werden getrennt behandelt: `fs.ErrNotExist` ⇒
      `ErrProjectNotInitialized` (TOCTOU/missing config), sonstiger
      I/O-/Permission-Fehler ⇒ technischer Wrap ohne Add-Sentinel.
      Parse-Fehler ⇒ `ErrProjectNotInitialized` mit gewrapptem
      Parse-Detail.
   3. `services[name]`-Lookup mit vier beobachtbaren YAML-Subzuständen:
      - Eintrag fehlt
      - vorhanden, `Enabled == nil` (Pointer-Semantik T1: unset)
      - vorhanden, `*Enabled == true`
      - vorhanden, `*Enabled == false`
   4. Compose-Marker-Check: `s.fs.Exists(baseDir/compose.yaml)` →
      bei `(false, nil)` gilt block-absent (kein read). `Exists`-
      Fehler werden als technischer Wrap ohne Add-Sentinel
      zurückgegeben. Sonst
      `s.fs.ReadFile(baseDir/compose.yaml)` lesen. `ReadFile`-Fehler
      werden nicht als block-absent interpretiert: `fs.ErrNotExist`
      ⇒ block-absent (TOCTOU/missing compose), sonstiger
      I/O-/Permission-Fehler ⇒ technischer Wrap ohne Add-Sentinel.
      Danach `managedblock.Find(content, Marker{StyleHash,
      serviceMarkerName(name)})` auf dem gelesenen Content:
      - nil-error ⇒ block-present
      - `ErrBlockNotFound` ⇒ block-absent
      - `ErrBlockMalformed` ⇒ return
        `fmt.Errorf("%w: malformed managed compose-block for service %q",
        ErrServiceInconsistent, name)`; fachliche Inkonsistenz,
        CLI-Code 10, kein neuer Sentinel, kein siebter `ServiceState`.
        Diese Prüfung passiert vor der Kombinations-Klassifikation und
        gilt deshalb auch für `enabled: false` und `enabled` unset.
      - sonstige Fehler ⇒ wrap + return (generic/technical).
   5. Kombinations-Klassifikation gemäß T1-Tabelle:

      | YAML-Eintrag | `Enabled` | Compose-Block | → State                  |
      | ------------ | --------- | ------------- | ------------------------ |
      | fehlt        | —         | fehlt         | `Unregistered`           |
      | fehlt        | —         | vorhanden     | `InconsistentYAML`       |
      | vorhanden    | `true`    | vorhanden     | `Active`                 |
      | vorhanden    | `true`    | fehlt         | `InconsistentBlock`      |
      | vorhanden    | `false`   | *             | `Deactivated`            |
      | vorhanden    | `nil`     | *             | `EnabledUnset`           |

   **`Add(ctx, req)`-Top-Level-Dispatch:**
   1. Validate `req.BaseDir != ""` (analog Init; nicht-sentinel-Error).
   2. Catalogue-Check vor State-Detection:
      `supportedServices()` liefert heute nur `"postgres"`.
      Match-Miss ⇒ `ErrServiceUnsupported`. Catalogue-Fehler ist
      projekt-state-unabhängig — daher zuerst.
   3. `detectServiceState` ausführen.
      Wenn `detectServiceState` wegen malformed Compose-Block
      `ErrServiceInconsistent` returnt, bricht `Add` direkt ab; dieser
      Pfad wird nicht über einen `ServiceState` modelliert.
   4. Switch:
      - `Active` ⇒ `AddServiceResponse{PriorState: Active, State:
        Active, Changed: nil}, nil` (core-state no-op,
        idempotent). **T3 voll implementiert.** T4 erweitert diesen
        Pfad um die PostgreSQL-Artefaktprüfung: der endgültige
        CLI-Flow ist nur dann ein No-op, wenn neben
        `services.postgres.enabled: true` und dem `service.postgres`-
        Compose-Block auch der `volume.postgres`-Compose-Block und
        der `.env.example`-Block `service.postgres` konsistent
        vorhanden sind.
      - `InconsistentYAML` ⇒ `fmt.Errorf("%w: managed compose-block
        for service %q has no matching u-boot.yaml anchor; remove
        the block manually or restore the anchor", ErrServiceInconsistent,
        name)` (Repair-Hint per Spec §895). **T3 voll implementiert.**
      - `Unregistered`/`Deactivated`/`EnabledUnset`/`InconsistentBlock`
        ⇒ `planAdd(req, state)` → `executeAdd(plan)`. **T3:
        `planAdd` voll implementiert; `executeAdd` ist Stub.**

   **Plan-Datenstruktur:**
   ```go
   type serviceAction int

   const (
       actionRegister     serviceAction = iota // Unregistered → Active
       actionReactivate                        // Deactivated/EnabledUnset → Active
       actionRebuildBlock                      // InconsistentBlock → Active
       actionRepairArtifacts                   // Active core-state, missing LH-FA-ADD-002 artefacts
   )

   type servicePlan struct {
       Service    domain.ServiceName
       PriorState domain.ServiceState
       Action     serviceAction
       // T4 erweitert: YAMLPatch, ComposeBlockBody, EnvBlockBody,
       // sowie die per-File-Pläne (analog filePlan aus initproject.go).
   }
   ```

   **Tests (T3-Scope, alle in `_test`-Package mit Fake-FS aus
   `fakes_test.go`):**
   - `detectServiceState`: 7 Fixtures — 6 State-Fixtures aus der
     Tabelle plus ein malformed-managed-block-Fixture. Die 6 State-
     Fixtures assertieren auf `domain.ServiceState`; das malformed-
     Fixture assertiert `errors.Is(err, ErrServiceInconsistent)`.
     Jeder Fixture-Setup nur die zwei nötigen Files (`u-boot.yaml`,
     ggf. `compose.yaml`).
     Zusätzlich pinnt ein Test, dass ein malformed `service.postgres`-
     Block auch bei `services.postgres.enabled: false` nicht als
     deactivated durchrutscht, sondern vor der State-Klassifikation mit
     `ErrServiceInconsistent` abbricht.
   - `Add`-Top-Level:
     - `Active` → no-op-Response; zusätzlich Assertion dass kein
       `WriteFile`/`Mkdir` am Fake-FS aufgerufen wurde.
     - `InconsistentYAML` → `ErrServiceInconsistent` (errors.Is).
     - Unsupported name (z. B. `"redis"`) → `ErrServiceUnsupported`
       (errors.Is); **vor** State-Detection — Test prüft dass kein
       `Exists`/`ReadFile` aufgerufen wurde.
     - `Project-not-initialized` (kein `u-boot.yaml`) →
       `ErrProjectNotInitialized` (errors.Is).
     - `u-boot.yaml`-Read-Fehler nach erfolgreichem `Exists` →
       non-nil error, aber **kein** `ErrProjectNotInitialized`
       (Permission/I/O ist technisch, nicht Projektzustand).
     - `compose.yaml`-Read-Fehler nach erfolgreichem `Exists` →
       non-nil error, aber **kein** `ErrServiceInconsistent` und
       kein block-absent-Fallback (Permission/I/O ist technisch).
     - `u-boot.yaml`- oder `compose.yaml`-`Exists`-Fehler →
       non-nil error, aber **kein** `ErrProjectNotInitialized`,
       **kein** `ErrServiceInconsistent` und kein Missing-/Absent-
       Fallback (Permission/I/O ist technisch).
     - `BaseDir == ""` → non-nil error (kein Sentinel; analog
       `InitProjectService.Init`).
   - Plan-Bildung: für jeden der 4 mutierenden States ein Test, der
     den Stub-Execute-Error abfängt und den Plan via
     `export_test.go` beobachtbar macht — Assertions auf
     `plan.Action` + `plan.PriorState`.
     Kein FS-Write erwartet (Stub-Execute returnt vor jeder
     Mutation).

   **DoD T3:**
   - `YAMLCodec.PatchScalar` im Port, im yaml-Adapter und im Fake
     implementiert + getestet.
   - `detectServiceState`-Tests grün für alle 6 States plus malformed-
     block-Inkonsistenz.
   - Top-Level-`Add` returnt korrekte Sentinels für die fachlichen
     Error-Pfade, technische Non-Sentinel-Fehler für `Exists`- und
     `ReadFile`-Failures und no-op für `Active`.
   - Plan-Bildung asserted für die 4 mutierenden States
     (`Unregistered`, `Deactivated`, `EnabledUnset`,
     `InconsistentBlock`).
   - `make gates` grün.
   - DoD-Line in der Slice-Datei: `T3 ✅ <commit-hash>` (Konvention
     [[feedback-done-slice-dod-hash]]).

4. **T4 — PostgreSQL-Templates + Write-Pfad.**

   T4 ist die größte Tranche des Slices und wird analog zu M3-T4
   in drei Sub-Tranchen geschnitten, damit Reviews jeweils einen
   isolierten Konzern bekommen und Commit-Hashes pro Konzern in die
   DoD-Line wandern können (Konvention
   [[feedback-done-slice-dod-hash]]).

   - **T4a** macht das M3-Compose-Scaffold add-on-fähig
     (init-Block trennt sich von den add-on-veränderlichen
     Top-Level-Maps; reine M3-Refactor, kein Add-Code).
   - **T4b** liefert die `PatchMappingEntryYAML`-Port-Erweiterung
     mit dem yaml.v3-Adapter und seinem internen Marker-Scanner
     (driven-only, kein Application-Touch).
   - **T4c** verschmilzt T4a + T4b mit den PostgreSQL-Templates,
     `executeAdd`-Implementation, Active-Repair-Pfad,
     `.env.example`-Strategie und End-to-end-Tests.

   Reihenfolge ist hart: T4c verbraucht beides; T4a und T4b sind
   in beliebiger Reihenfolge implementierbar, aber T4a-zuerst macht
   die T4b-Adapter-Tests realistischer (Fixtures aus dem neuen
   Scaffold).

   ---

   **T4a — Compose-Scaffold-Restrukturierung + `renderManagedBlockOnly`.**

   *Auslöser:* Die aktuelle M3-`compose.yaml.tmpl` enthält
   `services: {}` innerhalb des `BEGIN/END ... init`-Blocks.
   Service-/Volume-Add-on-Blöcke dürfen dort **nicht** eingefügt
   werden, weil ein späteres `u-boot init --force` den kompletten
   `init`-Block ersetzt und sonst alle Add-on-Blöcke löscht. T4a
   schneidet das Scaffold so, dass nur die init-eigene
   Basiskonfiguration im `init`-Block liegt; die add-on-veränderlichen
   Top-Level-Maps `services:` und `volumes:` liegen außerhalb. Diese
   Tranche enthält bewusst keinen Add-Code — sie ist reines M3-Refactor
   plus Regression-Pinning.

   *Datei-Layout:*
   - `internal/hexagon/application/templates/compose.yaml.tmpl` —
     `init`-Block enthält nur init-eigene Basis (z. B. `name: …`,
     `networks: {}` falls nötig); `services:` und `volumes:` als
     leere Top-Level-Maps **außerhalb** des Blocks.
   - `internal/hexagon/application/initproject.go` — neuer Helper
     `renderManagedBlockOnly(rendered []byte, marker
     managedblock.Marker) ([]byte, error)` extrahiert den
     `BEGIN/END ... <marker.Name>`-Bereich aus dem voll gerenderten
     Template. `executeReplaceBlock` ruft den Helper bei
     `actionReplaceBlock` auf und splict damit nur den
     Block-Bereich in das existierende File. Templates, die außerhalb
     des Blocks Content führen (heute nur `compose.yaml.tmpl`), bleiben
     dadurch im Re-Init-Pfad nicht-destruktiv.
   - **Ensure-Scaffold-Pass nach Block-Replace.** Splice allein
     reicht nicht: die heutige M3-`compose.yaml.tmpl` hat
     `name: …`, `services: {}` und `networks: default: …` *innerhalb*
     des `init`-Blocks (kein `volumes:`). Wenn T4a den neuen
     Init-Block (mit `name:` + `networks:`, aber ohne `services:`)
     einsplict, wäre `services:` nach `u-boot init --force` aus der
     Datei verschwunden — und `volumes:` hat noch nie existiert,
     muss aber für T4c-Add-Patches verfügbar sein. Damit `u-boot
     add` danach zuverlässig patchen kann, ergänzt
     `executeReplaceBlock` für `compose.yaml` nach dem
     Block-Splice einen Ensure-Pass:
     prüfe, ob die Top-Level-Keys `services:` und `volumes:` außerhalb
     jedes Managed-Blocks existieren; falls nicht, hänge `services: {}`
     bzw. `volumes: {}` als neue Top-Level-Maps an (mit genau einer
     trennenden Leerzeile vom vorhergehenden Top-Level-Block).
     `networks:` lebt heute im `init`-Block (mit User-relevanten
     Defaults wie `default.name`) und bleibt dort; der Ensure-Pass
     fasst Top-Level-`networks:` nicht an, weil das ein User-Custom-
     Block sein könnte. Manuelle Top-Level-Keys bleiben generell
     unangetastet. Der Ensure-Pass läuft nur für `compose.yaml`;
     andere Templates (`README.md` etc.) haben keinen analogen Bedarf.
   - `internal/hexagon/application/initproject_test.go` — neue
     Regression-Test-Fälle (s. u.).

   *Helper-Vertrag (`renderManagedBlockOnly`):*
   - Findet `marker.Begin()` und `marker.End()` im gerenderten
     Template (gleiche Suchsemantik wie `managedblock.Find`).
   - Rückgabe: der vollständige Bytebereich von der `BEGIN`- bis
     einschließlich der `END`-Zeile (also den Marker mit Inhalt),
     bereit zum splicen in den bestehenden File-Block-Bereich.
   - Fehler: `managedblock.ErrBlockNotFound` / `ErrBlockMalformed`
     wenn das Template selbst keinen oder einen kaputten Block
     enthält — Production-Templates dürfen nicht kaputt sein, also
     ist das ein Programming-Error, kein User-Fehler.

   *Verhältnis zu `actionWrite`:* T4a ändert **nichts** am Fresh-Init-
   Pfad. `actionWrite` schreibt weiter das komplette gerenderte
   Template (init-Block + leere Top-Level-Maps); nur
   `actionReplaceBlock` wird auf den Helper umgestellt.

   *Tests (T4a-Scope):*
   - Fresh-Init: `u-boot init` → `compose.yaml` enthält genau einen
     `init`-Block + genau ein Top-Level-`services:` (leer) + genau
     ein Top-Level-`volumes:` (leer), Marker-Bytebereich endet vor
     dem ersten Top-Level-Key außerhalb.
   - Re-Init mit Add-on-Blöcken: Test seed'et eine `compose.yaml`
     mit (a) altem Init-Block-Body, (b) zusätzlichem
     `service.postgres`-Marker unter `services:`, (c) zusätzlichem
     `volume.postgres`-Marker unter `volumes:`. `u-boot init --force`
     muss den Init-Block neu rendern, beide Add-on-Marker erhalten
     und keinen Top-Level-Duplikat erzeugen.
   - **Migration-Test Alt-M3-Compose:** Test seed'et eine
     `compose.yaml` mit der echten heutigen M3-Form: `init`-Block
     enthält `name: <project>`, `services: {}` und `networks:
     default: name: <project>-default`; **kein** `volumes:` (das
     fügt T4a/T4c neu hinzu, da M3 es nie hatte).
     `u-boot init --force` muss danach produzieren:
     - genau einen `init`-Block mit `name: …` und
       `networks: default: …` (ohne `services: {}` darin),
     - genau einen Top-Level-`services:` (leer) außerhalb des Blocks,
       angelegt durch Ensure-Scaffold,
     - genau einen Top-Level-`volumes:` (leer) außerhalb des Blocks,
       auch via Ensure-Scaffold (neuer Key, in Alt-M3 nicht vorhanden).
     Pinnt den Ensure-Scaffold-Pfad und verhindert, dass ein bloßer
     Re-Init die leeren Add-on-Hosts aus alten Projekten löscht.
   - **Migration + User-Custom-Top-Level-Key:** wie oben, aber die
     Alt-Compose enthält zusätzlich einen User-Top-Level-Key
     außerhalb des `init`-Blocks (z. B. `x-user-config:` mit
     User-Inhalt; `networks:` testet das nicht ehrlich, weil das
     bereits im init-Block-Body lebt). Der Ensure-Pass darf den
     User-Block nicht anfassen oder duplizieren.
   - Helper-Direkttests: `renderManagedBlockOnly` für ein Template
     ohne den geforderten Marker (ProgrammingError-Pfad) und mit
     malformed-Block; beide assertieren auf das `managedblock`-
     Sentinel.
   - Backup-Verhalten unverändert pro existierender `planFile`-
     Priorisierung (`internal/hexagon/application/initproject.go`):
     - **`--force` alleine + Block vorhanden** ⇒ `actionReplaceBlock`,
       Backup=false; nur der `init`-Block wird gerendert-extrahiert
       und gesplict.
     - **`--force --backup` + Block vorhanden** ⇒ `actionReplaceBlock`
       mit Backup=true; gleiche Splice-Mechanik plus Backup.
     - **`--backup` alleine** ⇒ `actionOverwriteFull`; schreibt den
       kompletten gerenderten Output (init-Block + leere
       Top-Level-Maps). **Akzeptiert destruktiv für den
       Live-State**: vorhandene Add-on-Blöcke
       (`service.postgres`/`volume.postgres`) verschwinden aus der
       Live-`compose.yaml`. Die `.bak`-Datei sichert den
       Recovery-Pfad — der User kann sie zurückkopieren oder
       `u-boot add postgres` erneut laufen lassen, um die Add-on-
       Blöcke wiederherzustellen. T4a ändert das Verhalten bewusst
       nicht: `--backup`-Semantik ist „Reset auf Template, Backup
       für Rollback"; wer Add-on-Blöcke beim Re-Init bewahren will,
       nutzt `--force --backup`. M3-Kompatibilität bleibt
       gewahrt; die `init --backup`-Konvention ist seit M3-T4a so
       dokumentiert.
       Test pinnt das explizit als documented-behavior: eine
       Add-on-bestückte `compose.yaml` + `--backup` ⇒ Live-Datei
       enthält keine Add-on-Blöcke mehr, `.bak`-Datei enthält
       sie noch.
     T4a testet alle drei Pfade getrennt mit einer Add-on-vorbelegten
     `compose.yaml`.

   *DoD T4a:*
   - Neue Tests grün; bestehende Init-Tests (insbesondere
     `--force` mit managed-block) unverändert grün.
   - Helper `renderManagedBlockOnly` 100 % Coverage.
   - `make gates` grün.
   - DoD-Line: `T4a ✅ <commit-hash>` (+ Review-Fix-Commit falls
     anfallend, analog M3-T4a).

   ---

   **T4b — `PatchMappingEntryYAML`-Port + yaml-Adapter + Fake.**

   *Auslöser:* T4c braucht einen byte-erhaltenden, strukturiert
   validierten YAML-Mapping-Patch, um managed Compose-Blocks unter
   `services:`/`volumes:` einzufügen, ohne die Datei komplett zu
   re-marshalen (das würde Kommentare und manuelle Einträge zerstören;
   `PatchScalar` aus T3 ist nur für skalare Leaves konzipiert). T4b
   liefert diesen Port plus Adapter plus Fake — und nichts sonst. Die
   Application-Schicht wird in T4b nicht angefasst; der Port hat
   keinen Caller bis T4c.

   *Datei-Layout:*
   - `internal/hexagon/port/driven/yamlcodec.go` — `YAMLCodec` um
     **zwei** Methoden erweitern:
     1. `PatchMappingEntryYAML(content []byte, parentKey string,
        entryKey string, valueYAML []byte, markerName string)
        ([]byte, error)` — der mutierende Patch.
     2. `LocateMarkedEntry(content []byte, parentKey string,
        entryKey string, markerName string) (LocateResult, error)`
        — read-only-Inspektion, damit Application-Code Anker-Checks
        und Inhalts-Prüfungen machen kann, **ohne** yaml.v3
        importieren zu müssen (depguard-Regel `application-no-yaml`).

     ```go
     type LocateResult struct {
         EntryExists         bool   // parentKey.entryKey existiert als Mapping-Eintrag
         MarkerInEntry       bool   // Block markerName hängt direkt als Sub-Element von entryKey
         MarkerSomewhereElse bool   // Block markerName existiert, aber NICHT unter entryKey
         BlockBody           []byte // wenn MarkerInEntry: rohe Bytes zwischen BEGIN/END (ohne die Marker-Zeilen)
     }
     ```

     Beide Methoden nutzen denselben adapter-lokalen Hash-Marker-
     Scanner; `LocateMarkedEntry` ist im Wesentlichen Phase 1+2 des
     `PatchMappingEntryYAML`-Drei-Phasen-Algorithmus ohne den
     Splice-Schritt.

     Plus neue Port-Sentinels:
     `ErrYAMLFragmentInvalid` (Scalar/Sequence in `valueYAML`,
     Top-Level-Key-Duplikat, Parent-Shape ≠ mapping/null/{}) und
     `ErrYAMLAnchorMismatch` (Block mit `markerName` existiert,
     aber außerhalb des `parentKey.entryKey`-Bereichs).

     **Eindeutige Fehlerverteilung zwischen den beiden Methoden** —
     damit die Application keinen technischen Driven-Port-Error im
     Pre-Check-Pfad fangen muss:
     - `LocateMarkedEntry` returnt für jeden **wohlgeformten** Input
       ein `LocateResult` ohne Fehler — auch für Wrong-Anchor
       (User-Manual-Entry oder Marker-am-falschen-Ort). Der Caller
       branch't auf die `LocateResult`-Flags und übersetzt
       fachlich in `ErrServiceInconsistent`. Nicht-nil-Errors
       kommen nur aus echten Adapter-Fehlern: Parse-Fehler im
       `content`, malformed Managed-Block-Marker (`BEGIN` ohne
       `END` / duplicate `BEGIN`).
     - `PatchMappingEntryYAML` gibt `ErrYAMLAnchorMismatch` zurück,
       wenn der vorgelagerte Locate-Check unterlaufen wird (z. B.
       von Future-Callern außerhalb von `AddServiceService`). Das
       ist reine Defense-in-Depth; der primäre Application-Check
       läuft über `LocateMarkedEntry`-Flags und feuert den
       Sentinel nie.
     `parentKey` ist bewusst ein einzelner String (kein `path []string`):
     M5 patcht ausschließlich Top-Level-Maps in `compose.yaml`
     (`services`, `volumes`); ein generischer Pfad-Parameter würde
     Anlege-Semantik für tiefere Ebenen suggerieren, die der
     Adapter-Scanner gar nicht implementiert. Bei realem Bedarf für
     nested patches kommt ein eigener Slice.
   - **Byte-Erhaltungs-Vertrag (präzise):** der Adapter darf an
     `content` **nur** drei Bytebereiche verändern und keinen
     darüber hinaus:
     1. den Bytebereich des Ziel-Managed-Blocks (vom `BEGIN`-
        Marker bis einschließlich `END`-Marker plus terminierender
        Zeilenende-Sequenz);
     2. wenn der Block neu angelegt wird: einen Insertion-Range an
        der Stelle des `entryKey`-Mapping-Eintrags unter dem
        Parent-Key — alle anderen Mapping-Entries unter demselben
        Parent (insbesondere User-Custom-Services / -Volumes
        zwischen den Markern oder als Siblings des Marker-Blocks)
        bleiben byte-identisch;
     3. wenn der Parent-Key fehlt oder in der einzeiligen
        `{}` / `null`-Form steht: die einzeilige Parent-Header-
        Zeile darf in eine mehrzeilige Form umgeschrieben werden
        (`services: {}` → `services:\n  …`), und/oder ein neuer
        Parent-Top-Level-Block darf am Datei-Ende angelegt werden
        (mit genau einer trennenden Leerzeile vor dem ersten
        Eintrag); jeder andere bestehende Top-Level-Block bleibt
        byte-identisch.
     Konsequenzen: User-Kommentare zwischen Marker-Blöcken,
     User-Custom-Services außerhalb der Marker, Top-Level-Comments
     vor/nach jedem Parent-Block — alle byte-identisch. Der
     Adapter rewrite't **nicht** den ganzen Parent-Map-Bereich,
     selbst wenn es bequemer wäre.
   - `internal/adapter/driven/yaml/codec.go` — Adapter-Implementation
     mit drei Phasen:
     1. **Validate** mit yaml.v3: `content` parsen, `valueYAML`
        parsen; sicherstellen dass `valueYAML` ein Mapping-Root ist
        (Scalar/Sequence ⇒ Fehler), Top-Level-Key-Duplikate prüfen,
        Parent-Key-Shape prüfen — **erlaubt**:
        - **block-style Mapping** (`services:` gefolgt von eingerückten
          Sub-Keys auf eigener Zeile);
        - **leeres Flow-Mapping** (`services: {}`) — wird in
          block-style umgeschrieben;
        - **null / fehlend** — wird als block-style Mapping angelegt.

        **Nicht erlaubt** (technischer Reject mit
        `ErrYAMLFragmentInvalid`):
        - **nicht-leeres Flow-Mapping** (`services: { mywebapp: {...} }`)
          — ein block-style Managed-Entry kann nicht ohne komplettes
          Umschreiben der Parent-Map eingefügt werden, was den
          Byte-Erhaltungs-Vertrag verletzt. Reject ist konservativer
          als stiller Style-Wechsel; ein V1-Folge-Slice könnte einen
          expliziten Flow→Block-Migration-Pfad anbieten, MVP nicht.
        - **Scalar/Sequence** als Parent-Wert (überschneidet sich mit
          `ErrYAMLPathInvalid`-Pfad in `PatchScalar`; hier eigener
          Sentinel für PatchMappingEntryYAML).
     2. **Locate** mit adapter-lokaler Hash-Marker-Scanner-Logik:
        Bytebereich der Parent-Map (`<key>:` bis nächste
        gleiche-Einrückung-Top-Level-Zeile) und des Managed-Blocks
        (`# BEGIN U-BOOT MANAGED BLOCK: <markerName>` bis
        `# END U-BOOT MANAGED BLOCK: <markerName>`) bestimmen.
        **Adapter-seitige Anker-Validation (symmetrisch):** Defense-
        in-Depth gegen Caller-Bugs außerhalb von `AddServiceService`.
        Beide Richtungen werden vom Adapter geprüft:
        1. **Marker → Eintrag**: Block mit `markerName` existiert
           außerhalb des Bytebereichs von `parentKey.entryKey` ⇒
           `ErrYAMLAnchorMismatch`.
        2. **Eintrag → Marker**: Mapping-Eintrag `parentKey.entryKey`
           existiert, enthält aber keinen Sub-Element-Marker mit
           `markerName` ⇒ `ErrYAMLAnchorMismatch`. Verhindert, dass
           der Adapter eine User-manuelle Definition unter unserem
           Schlüssel überschreibt oder dass ein zweiter
           `entryKey`-Sub-Key in derselben Parent-Map landet
           (duplicate-key).
        Application-Code (T4c-Pre-Patch-Anker-Check) liefert den
        fachlichen `ErrServiceInconsistent` *vor* diesem technischen
        Sentinel; beide bleiben unabhängig und ein Future-Caller
        außerhalb von `AddServiceService` ist trotzdem sicher.
     3. **Splice** als Byte-Range-Edit; alles außerhalb der berührten
        Bereiche unverändert.
   - `internal/adapter/driven/yaml/codec_test.go` — Adapter-Tests
     (s. u.).
   - `internal/hexagon/application/fakes_test.go` — `fakeYAML`-Fake
     bekommt **beide** neuen Methoden (`PatchMappingEntryYAML` +
     `LocateMarkedEntry`). Anders als beim Production-Adapter
     darf der Fake `managedblock` importieren (er ist `_test.go`),
     was die Fake-Implementation drastisch vereinfacht; trotzdem
     muss er den Byte-Erhaltungsvertrag dieselben asserts honorieren
     und denselben `LocateResult`-Vertrag liefern (Anker-State-Drift
     zwischen Fake und Production wäre ein Test-Leak), damit
     Application-Tests in T4c repräsentative Outputs sehen.

   *Schicht-Hygiene:* Der Adapter importiert **nicht**
   `internal/hexagon/application/managedblock`, weil die
   Depguard-Regel `adapter-no-application`
   Adapter→Application-Imports verbietet (`spec/architecture.md` §4,
   `.golangci.yml` §depguard). Die Scanner-Logik muss die
   Hash-Marker-Literale `# BEGIN U-BOOT MANAGED BLOCK: <name>` /
   `# END U-BOOT MANAGED BLOCK: <name>` deshalb redundant
   implementieren. Trade-off: zwei Implementations derselben
   Marker-Konvention. Eine spätere Konsolidierung (Marker-Konstanten
   in ein neutrales `internal/markers`-Package, das beide Schichten
   importieren dürfen) ist bewusst kein T4-Carveout — wird relevant,
   sobald V1 einen dritten Marker-Consumer bekommt; aktuell kein
   sichtbarer Schaden.

   *Patch-Verhalten (deterministisch):*
   - `valueYAML` muss ein Mapping-Root sein. Scalar oder Sequence
     ⇒ `ErrYAMLFragmentInvalid` (neuer Sentinel im Port).
   - Doppelte Top-Level-Keys in `content` ⇒ technischer
     Patch-Fehler vor jedem Write (kein silent-merge).
   - Existierender Parent-Key, der weder Mapping noch leeres
     `{}`/`null` ist (z. B. ein Scalar) ⇒ technischer Fehler;
     PatchMappingEntryYAML ersetzt nichts.
   - Fehlt der Parent-Key komplett, wird er als neuer
     Top-Level-Key **außerhalb** jedes Managed-Blocks angelegt
     (am Datei-Ende, getrennt durch genau eine Leerzeile vom letzten
     Top-Level-Block).
   - Ist der Parent-Key als `{}` oder `null` vorhanden, wird die
     einzeilige Form in eine mehrzeilige Mapping-Form aufgelöst.
   - Existiert der Ziel-Managed-Block, wird nur sein Bytebereich per
     adapter-lokalem Byte-Splice ersetzt; fehlt er, wird ein neuer,
     eingerückter Hash-Managed-Block unter dem Parent-Key
     eingefügt.
   - Bytes außerhalb des angelegten/ersetzten Parent- oder
     Managed-Block-Bereichs bleiben unverändert; Kommentare und
     manuelle Einträge außerhalb dieses Bereichs dürfen nicht durch
     YAML-Reformatting wandern.

   *Adapter-Tests (T4b-Scope):*
   - Happy-Path-Insert: `compose.yaml` mit leerer Top-Level-Map
     `services: {}` und kein Marker → Patch fügt mehrzeilige
     `services: postgres: …` mit BEGIN/END ein; Kommentare oberhalb
     bleiben byte-identisch.
   - Happy-Path-Replace: existierender wohlgeformter
     `service.postgres`-Block → Patch ersetzt genau diesen Bereich.
   - Missing-Parent-Key: kein `services:` im File → Patch legt
     Top-Level-Key am Datei-Ende an.
   - Indented-Marker: existierende `services:`-Map mit dem
     Marker eingerückt unter dem Key (wie produktiv geschrieben);
     Scanner muss die gleiche Einrückungsebene erkennen.
   - Malformed-Marker: BEGIN ohne END / duplicate BEGIN →
     Fehler (analog zu `managedblock.ErrBlockMalformed`); kein
     Splice.
   - Top-Level-Key-Duplikate in `content` → technischer Fehler.
   - Non-Mapping-Fragment (Scalar/Sequence in `valueYAML`) → Fehler.
   - **Nicht-leere Flow-Style-Parent-Map:** seed'e `services: {
     mywebapp: { image: nginx } }` und versuche `entryKey="postgres"`
     einzufügen → `ErrYAMLFragmentInvalid` mit Hint auf
     block-style-Erfordernis. Leeres Flow-Mapping (`services: {}`)
     bleibt erlaubt und wird in block-style umgeschrieben — eigener
     Happy-Path-Test.
   - **`LocateMarkedEntry`-Tests** (read-only Inspektion, gleicher
     Scanner wie Patch; **alle wohlgeformten Fälle returnen
     `(LocateResult, nil)`** — kein Sentinel-Error für Wrong-Anchor):
     - Clean: kein `services:`, kein Marker → `LocateResult{}` (alle
       Flags false), `err=nil`.
     - Managed: Block korrekt unter `services.postgres` →
       `EntryExists=true, MarkerInEntry=true, BlockBody=` (Bytes
       zwischen BEGIN/END ohne Marker-Zeilen), `err=nil`.
     - User-Manual-Entry-ohne-Marker: `services.postgres: {image: ...}`
       ohne Marker → `EntryExists=true, MarkerInEntry=false,
       MarkerSomewhereElse=false`, `err=nil`. **Kein
       `ErrYAMLAnchorMismatch`** — der Application-Caller muss
       fachlich darauf branchen.
     - Marker am falschen Ort: `service.postgres`-Marker unter
       `volumes:` → `EntryExists=false, MarkerInEntry=false,
       MarkerSomewhereElse=true`, `err=nil`. Ebenfalls kein
       Sentinel.
     - Malformed Marker (BEGIN ohne END, duplicate BEGIN) → echter
       Adapter-Fehler analog `managedblock.ErrBlockMalformed`.
     - Parse-Fehler im `content` → wrapped Parse-Error.
     - Konsistenz-Test: für ein gegebenes `content` muss
       `LocateMarkedEntry` denselben Anker-State sehen wie
       `PatchMappingEntryYAML` (gemeinsamer Scanner, kein Drift).
   - **Anker-Mismatch (adapter-eigene Defense-in-Depth, symmetrisch):**
     - **Marker → Eintrag:** seed'e `content` so, dass der Block mit
       `markerName` existiert, aber unter einem anderen Parent-Key
       (`volumes:` statt `services:`) oder unter einem anderen
       Eintrag (`services.other` statt `services.<entryKey>`).
       Patch-Call mit `parentKey="services"`, `entryKey="postgres"`
       → `ErrYAMLAnchorMismatch`; kein Splice.
     - **Eintrag → Marker:** seed'e `content` mit einer User-
       manuellen `services.postgres:`-Definition (z. B. `image:
       postgres:14`) **ohne** den `# BEGIN ... service.postgres`-
       Marker als Sub-Element. Patch-Call mit `parentKey="services"`,
       `entryKey="postgres"`, `markerName="service.postgres"` →
       `ErrYAMLAnchorMismatch`; kein Splice und kein
       duplicate-key.
   - Byte-Erhaltungstest (gemäß präzisem Vertrag oben): `content`
     mit (a) Kommentaren oberhalb und zwischen Top-Level-Blöcken,
     (b) einem User-Custom-Service unter `services:` außerhalb
     jedes Markers, (c) einer User-Kommentar-Zeile zwischen dem
     User-Service und dem `service.postgres`-Marker. Nach Patch
     gilt byte-Identität für: alle Kommentare, den User-Service,
     alle anderen Top-Level-Blöcke. Geändert sein dürfen nur
     der Marker-Bereich selbst und ggf. die Parent-Header-Zeile
     (wenn `services: {}` aufgelöst wurde) oder der neu angelegte
     Parent-Block (wenn `services:` gefehlt hat).
   - Depguard-Verifikation: `make verify-depguard` ist bereits
     etabliert und prüft die `adapter-no-application`-Regel; T4b
     fügt keine neue Verifikation, vertraut auf das bestehende Gate.

   *DoD T4b:*
   - Adapter implementiert + alle obigen Tests grün.
   - Fake-Variante mit semantisch identischem Verhalten;
     `make gates` bleibt grün (Application-Code touched nicht).
   - DoD-Line: `T4b ✅ <commit-hash>` (+ Review-Fix-Commit
     falls anfallend).

   ---

   **T4c — postgres-Templates + `executeAdd` + Active-Repair + e2e.**

   Setzt T4a (Scaffold) und T4b (`PatchMappingEntryYAML`) voraus.
   Hier verschmilzt alles: Templates landen, `executeAdd` wird
   real, Active-Repair (`actionRepairArtifacts`-Pfad aus T3) wird
   ausgewertet, `.env.example`-Strategie wird implementiert,
   `u-boot.yaml` bekommt seinen `PatchScalar`-Patch, End-to-end-Tests
   pinnen den kompletten Add-Flow.

   *Templates:*
   - `internal/hexagon/application/templates/services/postgres.compose.tmpl`
     — YAML-Fragment für den Wert unter `services.postgres` mit
     `image: postgres:16-alpine`, `environment` (POSTGRES_USER /
     POSTGRES_PASSWORD / POSTGRES_DB), `volumes` (named-volume),
     `ports: ["5432:5432"]`, `healthcheck` (`pg_isready`). Das
     Template enthält **keinen** Top-Level-`services:`-Key — die
     Einfügung passiert via `PatchMappingEntryYAML` strukturiert
     unter der Top-Level-Map.
     Die `environment`-Werte nutzen Compose-Default-Syntax mit
     Fallback auf die `.env.example`-Defaults:
     `POSTGRES_USER: ${POSTGRES_USER:-postgres}`,
     `POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:-CHANGEME_POSTGRES_PASSWORD}`,
     `POSTGRES_DB: ${POSTGRES_DB:-postgres}`.
     **Begründung:** Docker Compose liest automatisch nur `.env`
     (nicht `.env.example`); ohne Defaults im Compose-File würde
     der Acceptance-Flow LH-AK-002 (`u-boot init && u-boot add
     postgres && u-boot up`, `spec/lastenheft.md` §LH-AK-002) ohne
     manuelles `cp .env.example .env` scheitern. `u-boot up` ist
     M6 und kann den Acceptance-Bug nicht für M5 schließen, also
     muss M5 selbst out-of-the-box laufbar sein. `.env.example`
     bleibt die kanonische Override-Vorlage; sobald der User
     Werte ändern will, kopiert er die Datei nach `.env` wie im
     Template-Kommentar dokumentiert. Die `CHANGEME_*`-Convention
     für `.env.example` bleibt unverändert (Sicherheits-Hinweis,
     dass Production-Werte ersetzt werden müssen).
   - `…/services/postgres.volume.tmpl` — YAML-Fragment für den Wert
     unter `volumes.postgres-data`; MVP reicht ein leeres Mapping
     `{}`. Fragment enthält keinen Top-Level-`volumes:`-Key.
   - `…/services/postgres.env.tmpl` — **reiner Variablen-Body**
     (drei Zeilen `POSTGRES_USER=postgres`,
     `POSTGRES_PASSWORD=CHANGEME_POSTGRES_PASSWORD`,
     `POSTGRES_DB=postgres`), **ohne** `# BEGIN/END U-BOOT MANAGED
     BLOCK: service.postgres`-Marker. Die Marker werden in T4c über
     einen Wrap-Schritt ergänzt; das Template selbst bleibt
     marker-frei, damit dieselben drei Zeilen auch für die
     Compose-`environment:`-Referenz (LH-AK-002-Default-Werte)
     wiederverwendbar bleiben. Sicherheits-Convention: explizit
     `CHANGEME_*`, nie reale Defaults.

   *Env-Block-Wrap-Helper:* T4c führt einen kleinen Helper
   `renderEnvManagedBlock(svc domain.ServiceName, varsBody []byte)
   ([]byte, error)` ein, der den reinen Variablen-Body in einen
   wohlgeformten Hash-Managed-Block einrahmt:

   ```
   # BEGIN U-BOOT MANAGED BLOCK: service.<svc>
   <varsBody>
   # END U-BOOT MANAGED BLOCK: service.<svc>
   ```

   Das ist der Vertrag für `managedblock.Replace`: das Replacement
   **muss** die BEGIN/END-Marker enthalten, sonst würde Replace
   beim nächsten `u-boot add postgres` den Block nicht mehr finden
   und einen zweiten Block ans Datei-Ende anhängen
   (Idempotenz-Bruch). Der Wrap-Helper macht diesen Vertrag
   explizit und verhindert, dass create/append/replace-Pfade
   unterschiedlich verpacken.

   Der Helper produziert für `domain.ServiceName("postgres")` exakt:
   `[]byte("# BEGIN U-BOOT MANAGED BLOCK: service.postgres\n…\n# END
   U-BOOT MANAGED BLOCK: service.postgres\n")`. Wird in allen drei
   Env-Pfaden (create, append, replace) verwendet, sodass die
   produzierte Block-Form garantiert byte-identisch ist und ein
   späterer `managedblock.Find` ihn wiederfindet.
   - Embed-Directive: `//go:embed templates/*.tmpl
     templates/services/*.tmpl`. `renderTemplate("services/postgres.compose.tmpl",
     …)` bleibt mit dem bestehenden `"templates/"+name`-Pfadmodell
     kompatibel. **`templateNames()` muss in T4c angepasst werden:**
     die heutige Implementation in
     `internal/hexagon/application/templates.go` macht ein flaches
     `templateFS.ReadDir("templates")` und würde `services` als
     Directory-Eintrag zurückgeben — nicht die Service-Templates
     darunter. T4c stellt `templateNames()` auf `fs.WalkDir` um und
     liefert relative Pfade (`"compose.yaml.tmpl"`,
     `"services/postgres.compose.tmpl"` etc.); der Integrity-Self-Test
     in `templates_test.go` wird auf die neue Namens-Form aktualisiert.

   *`executeAdd`-Architektur (Plan-and-Execute, klar getrennt):*

   T4c führt eine klare Drei-Schicht-Aufteilung ein, um den
   T3-Stub-Vertrag konsistent zu erweitern und die Doppelrolle
   „Detection vs. Plan-Bau" zu eliminieren:

   1. **`detectActiveArtifacts` ist ein pure classifier**
      (keine FS-Writes, keine Plan-Datenstruktur), läuft in `Add`
      nur für `state == Active` und liefert Flags:
      ```go
      type activeArtifactsStatus struct {
          ServiceStale       bool // service.postgres-Block existiert, Pflichtfeld fehlt
          VolumeMissing      bool // volume.postgres-Marker absent
          VolumeStale        bool // volume.postgres-Marker present, Inhalt verletzt Pflicht
          EnvMissingOrStale  bool // .env.example fehlt, Block absent, oder Pflicht-Key fehlt
      }
      ```
      Abort-Fälle (malformed, Wrong-Anchor, User-Manual-Entry) returnen
      `ErrServiceInconsistent` direkt aus dem Classifier; Flags werden
      dann nicht aufgesetzt.

   2. **`servicePlan` (aus T3) bekommt einen `RepairFlags`-Slot**:
      ```go
      type servicePlan struct {
          Service     domain.ServiceName
          PriorState  domain.ServiceState
          Action      serviceAction
          RepairFlags activeArtifactsStatus // nur für actionRepairArtifacts gefüllt
      }
      ```
      `planAdd` aus T3 setzt `RepairFlags` zero für die vier
      pre-existing mutierenden Actions; T4c-`Add` ergänzt
      `RepairFlags` aus `detectActiveArtifacts` wenn Active mit
      Repair-Bedarf.

   3. **`executeAdd(ctx, baseDir, plan)` ist die einzige Stelle, die
      die Plan-Phase und die Execute-Phase macht** — analog zum M3-
      `InitProjectService.Init`-Pattern (`initproject.go:174`).
      Innerhalb baut es ein internes `executePlan` mit den drei
      Datei-Slots auf (Plan-Phase a–h unten), validiert vollständig
      vor jedem Write, und schreibt nur die nicht-nil Slots
      (Execute-Phase). `executePlan` ist **kein Input** und **kein
      Rückgabewert** — rein internes Bauwerk, leichter testbar via
      `export_test.go`-Brücke wenn Tests die Plan-Bildung in
      Isolation prüfen wollen.

   **Signatur-Anpassung der T3-Stub-Signatur:** T3 lieferte
   `executeAdd(ctx context.Context, plan servicePlan)`. T4c stellt
   die Signatur auf `executeAdd(ctx context.Context, baseDir string,
   plan servicePlan)` um — `baseDir` wird neu durchgereicht, `plan`
   bleibt der T3-erweiterte `servicePlan` (jetzt mit `RepairFlags`).
   Der Call-Site in `Add` (`addservice.go:201`) wird entsprechend
   angepasst.

   **Internes `executePlan`**: drei optionale Datei-Slots,
   vor-validierte Zielinhalte. Lebt nur innerhalb von
   `executeAdd` (Plan-Phase baut auf, Execute-Phase liest):

   ```go
   type executePlan struct {
       // baseDir wird über den executeAdd-Parameter durchgereicht,
       // bewusst nicht im Struct dupliziert.
       UBootYAML  *fileWrite // nil bei actionRebuildBlock / actionRepairArtifacts
       Compose    *fileWrite // nil wenn kein Compose-Patch nötig (z. B. Env-only-Repair)
       EnvExample *fileWrite // nil wenn .env.example bereits passt
   }
   type fileWrite struct {
       Path string         // relativ zu baseDir (executeAdd-Parameter)
       Body []byte         // schon gerenderter/gepatchter Body
       Mode iofs.FileMode  // mode-preserved vom Read; defaultFileMode falls neu
   }
   ```

   1. **Plan-Phase** (jeder Schritt darf mit non-write Fehler
      abbrechen; Plan-Fehler garantiert: keine geschriebene Datei).
      Reihenfolge ist hart, weil spätere Schritte auf früheren
      aufbauen:

      a. **Render** alle drei Templates in Bytes (Production-Fail
         ⇒ Programming-Error).

      b. **Load** über einen gemeinsamen Helper `loadForPatch(path)
         (body []byte, mode iofs.FileMode, exists bool, err error)`
         analog zum M3-T4b-`planFile`-Pattern
         (`internal/hexagon/application/initproject.go`):
         1. `s.fs.Lstat(path)` zuerst. `iofs.ErrNotExist` ⇒
            `exists=false, mode=defaultFileMode (0o644), body=nil` —
            saubere Trennung von „Datei fehlt" vs. „Datei kaputt".
         2. Symlink-Reject: wenn `info.Mode()&iofs.ModeSymlink != 0`,
            Abort mit `driving.ErrBackupUnsupportedKind` (gleicher
            Sentinel wie M3-T4b: ein `.env.example -> /etc/passwd`-
            Symlink darf nie verfolgt werden, sonst würde Add das
            Linkziel lesen und überschreiben).
         3. Non-regular-Reject: wenn die Mode-Bits weder regular
            noch Symlink anzeigen (z. B. device file, named pipe),
            ebenfalls Abort — Add hat nichts mit special files
            zu tun.
         4. `mode = info.Mode().Perm()` einfrieren für den späteren
            `fileWrite.Mode`-Slot; neu angelegte Dateien (`exists=false`)
            bekommen `defaultFileMode`.
         5. `s.fs.ReadFile(path)` für `body`. TOCTOU zwischen Lstat
            und ReadFile (file verschwindet) wird als `exists=false`
            behandelt — gleiche Semantik wie der Lstat-Pfad.

         Per Helper-Aufruf:
         - `u-boot.yaml` (**mandatory**, sonst wäre `detectServiceState`
           schon mit `ErrProjectNotInitialized` abgebrochen). Der Caller
           prüft nach `loadForPatch` explizit `exists==true`; ein
           `exists==false` ist hier ein TOCTOU-Race (Datei zwischen
           State-Detection und Load verschwunden) und führt zu
           `fmt.Errorf("%w: u-boot.yaml vanished between
           detectServiceState and loadForPatch",
           driving.ErrProjectNotInitialized)`. **Nicht** in den
           create-Pfad fallen — sonst würde `PatchScalar` auf einem
           leeren Body eine halbierte `u-boot.yaml` mit nur
           `services.postgres.enabled: true` ohne `schemaVersion` /
           `project.name` schreiben (Datenverlust). Symlink hier
           wäre überraschend, wird vom Helper aber identisch
           abgelehnt.
         - `compose.yaml` (optional, `exists=false` triggert
           Bootstrap weiter unten).
         - `.env.example` (optional, `exists=false` triggert
           create-Pfad in der Env-Strategie weiter unten).

      b'. **Parse `u-boot.yaml`** in der Plan-Phase (Pre-Write-Fail-
         pflichtig). `detectServiceState` parsed den Body nur lokal
         für die State-Klassifikation und verwirft das Ergebnis;
         T4c braucht für den Bootstrap-Block den `project.name` und
         könnte ihn theoretisch aus dem Pfad ableiten, aber die
         Single-Source-of-Truth ist `u-boot.yaml`. Daher: Plan-Phase
         macht `s.yaml.Unmarshal(uBootYAMLBody, &cfg ubootYAMLConfig)`;
         der entstehende `cfg` liefert (i) `cfg.Project.Name` für den
         Bootstrap und (ii) den existierenden `services`-Subtree für
         den `enabled`-Vergleich aus Schritt e unten. Parse-Fehler ⇒
         `fmt.Errorf("plan: parse u-boot.yaml: %w", err)` als
         technischer Pre-Write-Fail (nicht `ErrProjectNotInitialized`,
         weil die Datei existiert — sie ist zwischen
         `detectServiceState` und Plan-Phase kaputt geworden, was ein
         seltenes TOCTOU-Race ist, aber semantisch nicht „Projekt
         existiert nicht" bedeutet).

      c. **Compose-Bootstrap** (action-agnostisch, läuft VOR jedem
         Compose-Patch und VOR dem Anker-Check) ist eng definiert:
         er greift **nur**, wenn `compose.yaml` komplett fehlt
         oder der gelesene Body nur Whitespace enthält. In dem Fall
         ersetzt T4c den Body in-memory durch einen minimalen
         Split-Block-Compose (init-Block mit `name: <cfg.Project.Name>`
         aus dem geparsten `u-boot.yaml` + leere Top-Level-Maps
         `services:`/`volumes:` außerhalb des Blocks). Der Bootstrap-Body lebt nur in der
         Plan-Phase; geschrieben wird er erst über die normale
         Compose-Slot-Mechanik unten.
         **Nicht-leere Compose ohne `init`-Marker** ist explizit
         **kein** Bootstrap-Trigger — ein User mit manuell
         gepflegter `compose.yaml` (Migration aus einem fremden
         Projekt, Setup vor `u-boot init` etc.) würde sonst seinen
         kompletten Compose-Inhalt verlieren. Stattdessen patcht
         T4c direkt in die vorhandene Datei; `PatchMappingEntryYAML`
         (T4b) legt fehlende Top-Level-Keys `services:`/`volumes:`
         deterministisch an und respektiert den T4b-Byte-Erhaltungs-
         Vertrag für alles andere. Das `init`-Block-Konzept ist
         eine M3-Konvention für `u-boot init`-Scaffolding, kein
         Add-Vorbedingung — `u-boot add` kommt mit jeder validen
         Compose klar.

      d. **Pre-Patch-Anker-Check** (action-agnostisch, vor jedem
         `PatchMappingEntryYAML`-Aufruf, **symmetrisch in beide
         Richtungen** — sonst entstehen Duplikate oder Überschreiben).
         Implementiert über den **`LocateMarkedEntry`-Port** (T4b);
         die Application hat **kein** eigenes YAML-AST und **kein**
         yaml.v3-Import. Pro Ziel-Marker (`service.postgres` unter
         `services.postgres`, `volume.postgres` unter
         `volumes.postgres-data`) ein `LocateMarkedEntry`-Call;
         die Branching-Logik:

         | `EntryExists` | `MarkerInEntry` | `MarkerSomewhereElse` | Bedeutung | Aktion |
         | --- | --- | --- | --- | --- |
         | false | false | false | clean: nichts da, frei für Anlegen | weiter zu Schritt f |
         | true  | true  | false | managed: bestehender u-boot-Block | weiter zu Schritt f (Replace) |
         | true  | false | false | **User-Manual-Entry ohne Marker** | `ErrServiceInconsistent` mit Repair-Hint „services.postgres exists but is not u-boot-managed; remove or rename the manual entry, or run `u-boot add postgres` after setting `services.postgres.enabled: false` in u-boot.yaml" |
         | false | false | true  | **Marker am falschen Ort** | `ErrServiceInconsistent` mit Repair-Hint „service.postgres marker exists outside services.postgres; please remove the orphan marker" |
         | true  | true  | true  | (technisch impossible — Locate-Adapter würde das als duplicate detecten und `ErrYAMLAnchorMismatch` zurückgeben) | Fehler propagieren |

         Beide Richtungen sind action-agnostisch und laufen für jeden
         mutierenden Pfad. Ein expliziter „Takeover"-Pfad (User-Entry
         in u-boot-Verwaltung übernehmen) ist V1 (siehe Diskussion
         in der Re-Emission-Semantik unten); MVP rejected explizit
         und gibt eine konkrete Repair-Anweisung.

         Defense-in-Depth: der Adapter (`PatchMappingEntryYAML` im
         Splice-Schritt) gibt sicherheitshalber auch
         `ErrYAMLAnchorMismatch` zurück, wenn der vorgelagerte
         Locate-Check unterlaufen wird (z. B. von Future-Callern
         außerhalb von `AddServiceService`).

      e. **Plane `UBootYAML`-Slot** nur wenn die Action einen
         YAML-Anker-Patch braucht — d. h. der aktuelle Zustand des
         `services.<name>.enabled`-Schlüssels nicht bereits `true`
         ist:
         - `actionRegister` (Unregistered): Slot setzen, weil der
           `services.<name>`-Eintrag komplett fehlt.
         - `actionReactivate` (Deactivated / EnabledUnset): Slot
           setzen, weil `enabled` `false` oder unset ist.
         - `actionRebuildBlock` (InconsistentBlock): Slot **nicht**
           setzen — `enabled: true` steht schon, der Compose-Block
           fehlt.
         - `actionRepairArtifacts` (Active mit fehlenden Artefakten):
           Slot **nicht** setzen — `enabled: true` steht schon.
         Wenn der Slot gesetzt wird: `s.yaml.PatchScalar(yamlBody,
         []string{"services", name.String(), "enabled"}, true)`
         liefert den `UBootYAML.Body`.

      f. **Plane `Compose`-Slot** nur wenn ein Compose-Patch nötig
         ist:
         - `actionRegister` / `actionReactivate` /
           `actionRebuildBlock`: beide Patches
           (`service.postgres` + `volume.postgres`) nacheinander
           über `PatchMappingEntryYAML` auf den (ggf. via
           Bootstrap erzeugten) Compose-Body anwenden, das Ergebnis
           als `Compose.Body` setzen.
         - `actionRepairArtifacts`: gepatcht werden alle in
           `detectActiveArtifacts` als **fehlend oder stale**
           markierten Compose-Artefakte. Das umfasst zwei Fälle:
           - **strukturell fehlend**: der `volume.postgres`-Marker
             ist im Compose-Body gar nicht vorhanden → Volume-Patch.
             (Der Service-Marker ist per Definition vorhanden, sonst
             wäre der T3-State `InconsistentBlock` statt `Active` und
             wir liefen in `actionRebuildBlock`.)
           - **stale content** (Inhalts-Check aus
             `detectActiveArtifacts`): der `service.postgres`-Marker
             ist vorhanden, aber sein Block-Body verletzt die
             Pflichtfeld-Präsenz (fehlendes `image`, `environment.POSTGRES_*`,
             `volumes`-Referenz, `ports`, `healthcheck`) →
             Service-Patch (Re-Emission). Analog für stale
             `volume.postgres`-Inhalt; die T3-Active-Detection
             liefert pro Artefakt ein „needs repair"-Flag, das
             T4c hier in Patch-Aufrufe übersetzt.

           Die nicht-betroffenen Artefakte bleiben byte-identisch
           (Garantie über den T4b-Byte-Erhaltungs-Vertrag). Wenn nur
           der Env-Block stale/missing ist und beide Compose-Marker
           OK sind, bleibt `Compose` nil.

      g. **Plane `EnvExample`-Slot** nach der .env.example-Strategie
         unten — Slot bleibt nil, wenn der Block bereits identisch
         vorhanden ist und kein Patch nötig wäre.

      h. **Plan-Konsistenz-Check:** mindestens ein Slot muss gesetzt
         sein, sonst hätte `Add` bereits den No-op-Pfad genommen.
         Defensive Assertion; verletzte Annahme ⇒
         Programming-Error.

   2. **Execute-Phase**: für jeden nicht-nil-Slot
      `s.fs.WriteFile(BaseDir/Path, Body, Mode)` aufrufen. Reihenfolge
      ist deterministisch (`UBootYAML` → `Compose` → `EnvExample`),
      damit ein partieller Write-Fehler immer den gleichen Sub-State
      hinterlässt und debuggbar bleibt. Mode-Preservation aus dem
      geladenen File analog M3-T4b.

   3. **Response**: `Changed` enthält genau die Pfade der
      nicht-nil-Slots (nicht alle drei pauschal). `PriorState` aus
      dem Plan, `State = Active`. Bei `actionRepairArtifacts` mit
      nur Env-Repair: `Changed = [".env.example"]`. Bei
      `actionRebuildBlock`: `Changed = ["compose.yaml"]`. Bei
      `actionRegister`: `Changed = ["u-boot.yaml", "compose.yaml",
      ".env.example"]`.

   *Active-Repair-Pfad (`detectActiveArtifacts` + `actionRepairArtifacts`):*
   Der T3-Active-Code-Pfad wird in T4c erweitert. Statt sofort no-op
   zurückzugeben, ruft `Add` (für `state == Active`) den pure
   classifier `detectActiveArtifacts(baseDir, svc, fs, yaml) ⇒
   (activeArtifactsStatus, error)` auf. Er macht keine FS-Writes
   und bekommt keine Plan-Datenstruktur — nur Flags zurück. Ein
   Marker allein reicht bewusst **nicht** als „komplett" —
   LH-FA-ADD-002 verlangt konkreten Inhalt (Compose-Service +
   Volume + .env.example-Einträge + Port + Healthcheck).

   Der classifier nutzt den **`LocateMarkedEntry`-Port** für die
   strukturellen Checks (kein yaml.v3-Import in Application) und
   inspiziert den zurückgegebenen `BlockBody` mit zeilenbasiertem
   Pattern-Matching für den Inhalts-Check. Das ist pragmatisch
   (deklariert als bewusster MVP-Trade-off; eine yaml.v3-genaue
   Parse-API käme bei Bedarf später), aber **scharf normalisiert**,
   damit triviale False-Positives (kommentierte Defaults,
   Healthchecks die disable'd sind, Keys außerhalb des
   `environment`-Sub-Blocks) nicht als „vollständig" durchgehen.

   **Scan-Regeln (für jeden Pflichtfeld-Check):**
   1. **Comment-Stripping vor Match**: jede Zeile wird `TrimSpace`'d;
      Zeilen, die nach Trim mit `#` beginnen, werden ignoriert.
      `# POSTGRES_PASSWORD: …` oder Inline-Kommentare am Zeilenende
      werden vor dem Match entfernt. Eine Zeile wie
      `POSTGRES_PASSWORD: x  # CHANGEME` zählt als Treffer
      (Pre-Kommentar-Teil enthält den Key), aber `# POSTGRES_PASSWORD: x`
      nicht.
   2. **Block-Kontext-Tracking**: ein einfacher Stack über
      Einrückungstiefe verfolgt, welcher Top-Level-Key (`environment:`,
      `volumes:`, `ports:`, `healthcheck:`) gerade aktiv ist.
      Ein `POSTGRES_USER:` außerhalb von `environment:` (z. B. unter
      `labels:` oder am Service-Root) zählt **nicht** als Treffer für
      die `environment`-Pflichtfelder. Indent-Decreases pop'en den
      Stack.
   3. **Healthcheck-`disable: true`-Ausnahme**: ein `healthcheck:`-
      Mapping mit `disable: true` zählt **nicht** als gültiger
      Healthcheck (Compose-Semantik: `disable: true` schaltet den
      Healthcheck aus, was LH-AK-002 verletzt). Der Stale-Detector
      muss diese Form explizit als „kein Healthcheck" werten.
   4. **Trimmed, non-empty value**: `image:` ohne Wert (`image:`
      gefolgt von Leerzeile/Kommentar) zählt nicht; `image: foo`
      zählt; `image: ""` zählt nicht.

   Das ist immer noch kein voller YAML-Parser, aber präzise genug
   für LH-FA-ADD-002/LH-AK-002. Negative Tests pinnen die obigen
   Edge-Cases explizit (s. Test-Liste unten).

   - **Service-Block** (`service.postgres` unter `services.postgres`)
     via `yaml.LocateMarkedEntry(composeBody, "services", "postgres",
     "service.postgres")`:
     - `EntryExists=false && MarkerSomewhereElse=false` ⇒ eigentlich
       `InconsistentBlock`-State; State-Detection hätte das bereits
       abgefangen, hier defensive Vollständigkeit.
     - `MarkerSomewhereElse=true` oder `EntryExists=true &&
       MarkerInEntry=false` ⇒ Abort mit `ErrServiceInconsistent`
       (Wrong-Anchor / User-Manual-Entry; identisch zum
       Pre-Patch-Anker-Check).
     - `MarkerInEntry=true` ⇒ **Inhalts-Check** auf `BlockBody`
       via zeilenbasiertem Substring-Matching. Pflichtfelder gemäß
       LH-FA-ADD-002 und LH-AK-002 — der Check prüft **Präsenz**,
       nicht Werte (User-Tuning ist erlaubt: Port-Mapping ändern,
       Healthcheck-Intervall anpassen, Postgres-Image-Version
       pinnen):
       - Zeile mit `image:` (irgendein scalar, nicht leer).
       - Zeilen mit `POSTGRES_USER:`, `POSTGRES_PASSWORD:`,
         `POSTGRES_DB:` (innerhalb eines `environment:`-Sub-Blocks).
       - Zeile mit `volumes:` und mindestens eine Folge-Zeile die
         `postgres-data` enthält.
       - Zeile mit `ports:` und mindestens eine Folge-Zeile mit
         `- ` (Sequenz-Eintrag) — LH-AK-002 verlangt explizit „der
         konfigurierte Port … ist auf localhost erreichbar".
       - Zeile mit `healthcheck:` und mindestens eine eingerückte
         Folge-Zeile (Mapping mit Sub-Key) — LH-AK-002 verlangt
         explizit „Container … erreicht den Healthcheck-Status
         `healthy`".
       Fehlt eine Pflicht-Komponente ⇒ `status.ServiceStale = true`.
       Der Volume-Body bleibt vom Inhalts-Check frei (MVP-Template
       ist leeres Mapping; User darf `driver:` etc. customizen).
       Malformed Block ⇒ `LocateMarkedEntry` returnt
       `ErrYAMLAnchorMismatch` oder `managedblock.ErrBlockMalformed`-
       Äquivalent → classifier macht Abort mit
       `ErrServiceInconsistent`.
   - **Volume-Block** (`volume.postgres` unter
     `volumes.postgres-data`) via `yaml.LocateMarkedEntry(composeBody,
     "volumes", "postgres-data", "volume.postgres")`:
     - `EntryExists=false && MarkerSomewhereElse=false` ⇒
       `status.VolumeMissing = true`.
     - Wrong-Anchor / User-Manual-Entry ⇒ Abort mit
       `ErrServiceInconsistent`.
     - `MarkerInEntry=true` mit Inhalts-Check: MVP-Template ist
       leeres Mapping; nicht-leeres User-Mapping erlaubt. Es gibt
       hier keine echten Pflichtfelder; `status.VolumeStale` bleibt
       false außer der Body ist syntaktisch kaputt (dann Abort).
   - **`.env.example`-Block** (`service.postgres` als
     Hash-Managed-Block; `.env.example` ist Single-Block-File, nicht
     mapping-strukturiert → kein `LocateMarkedEntry`-Aufruf, sondern
     direkt `managedblock.Find` + Inhalts-Scan auf Block-Body):
     - Datei oder Block fehlt ⇒ `status.EnvMissingOrStale = true`.
     - Block malformed ⇒ Abort mit `ErrServiceInconsistent`.
     - **Inhalts-Check** über den Block-Body (zeilenweise scan):
       muss eine Zeile `POSTGRES_USER=…`, eine `POSTGRES_PASSWORD=…`
       und eine `POSTGRES_DB=…` enthalten (Wert egal — User darf
       `CHANGEME_*` durch Production-Werte ersetzen). Fehlt eine
       Variable ⇒ `status.EnvMissingOrStale = true`.
   - **Alle Flags false** ⇒ `Add` returnt echten No-op, alle drei
     Slots in `executeAdd` bleiben nil, `Changed=nil`.
   - **Mindestens ein Flag true** ⇒ `Add` setzt
     `servicePlan.Action = actionRepairArtifacts` und
     `servicePlan.RepairFlags = status`, übergibt an `executeAdd`,
     das nur die geflaggten Artefakte patcht (Compose-Slot wenn
     ServiceStale/VolumeMissing/VolumeStale; EnvExample-Slot wenn
     EnvMissingOrStale; UBootYAML-Slot bleibt nil weil `enabled:
     true` schon steht).

   *Semantik der Re-Emission:* Marker = u-boot-managed. User-
   Modifikationen *innerhalb* eines Markers gelten als
   überschreibbar — der nächste `u-boot add postgres` darf den Block
   neu emittieren, wenn der Inhalts-Check eine Pflicht-Komponente
   vermisst. Wer Werte einzelner Pflichtfelder customizen will (z. B.
   Port, Healthcheck-Intervall, Image-Tag), behält die User-Edits,
   solange die Pflichtfeld-Präsenz nicht verletzt ist; das ist der
   einzige robuste Customization-Pfad.

   **Kein Marker-Entfernen als Schutzpfad:** der naheliegende
   Gedanke „User entfernt die `# BEGIN/END … service.postgres`-
   Marker, um den Block aus u-boot-Verwaltung zu nehmen" funktioniert
   **nicht**. Ein `services.postgres`-Eintrag in `u-boot.yaml` mit
   `enabled: true` plus fehlender Compose-Marker klassifiziert die
   T3-State-Detection als `InconsistentBlock`, und der nächste Add
   läuft in `actionRebuildBlock` und schreibt den Marker (mit dem
   Default-Template-Body) wieder rein — der User-Custom-Block
   außerhalb des Markers würde dann mit dem managed Block koexistieren
   und Docker Compose würde den Key duplizieren bzw. überschreiben.
   Wer einen Service komplett aus u-boot-Verwaltung nehmen will,
   setzt `services.postgres.enabled: false` in `u-boot.yaml` und
   benennt seinen Custom-Service auf einen anderen Compose-Key (z. B.
   `postgres-custom`). Eine explizite „Service aus u-boot-Verwaltung
   nehmen, aber unter gleichem Namen weiterpflegen"-Funktion ist V1
   (LH-FA-ADD-007 `u-boot remove` + Folge-Slice).

   *Tests für stale/incomplete (T4c-Scope):*
   - `Active` + `service.postgres`-Block ohne `image:` ⇒
     `actionRepairArtifacts`, Service-Block neu geschrieben.
   - `Active` + `service.postgres`-Block ohne `environment.POSTGRES_PASSWORD` ⇒
     `actionRepairArtifacts`, Service-Block neu.
   - `Active` + `service.postgres`-Block ohne `ports:` ⇒
     `actionRepairArtifacts`, Service-Block neu (LH-AK-002 verlangt
     Port-Erreichbarkeit).
   - `Active` + `service.postgres`-Block ohne `healthcheck:` ⇒
     `actionRepairArtifacts`, Service-Block neu (LH-AK-002 verlangt
     Healthcheck-Status `healthy`).
   - `Active` + `service.postgres`-Block mit allen Pflichtfeldern
     plus User-Custom-Port (`ports: ["5433:5432"]`) ⇒ echter No-op,
     User-Customization bleibt erhalten.
   - `Active` + `service.postgres`-Block mit allen Pflichtfeldern
     plus User-Custom-Healthcheck (`healthcheck: {interval: 5s, ...}`) ⇒
     echter No-op, User-Tuning bleibt erhalten.
   - `Active` + `.env.example`-Block ohne `POSTGRES_PASSWORD=` ⇒
     `actionRepairArtifacts`, Env-Block neu, `Changed=[".env.example"]`.
   - `Active` + `.env.example`-Block mit User-überschriebenem
     `POSTGRES_PASSWORD=actual-prod-secret` ⇒ echter No-op,
     User-Secret bleibt erhalten (Inhalts-Check prüft nur Existenz
     der Key-Zeile, nicht den Wert).

   **False-Positive-Negativ-Tests** (pinnen die Scan-Regeln aus dem
   Active-Repair-Abschnitt — würden alle ohne Comment-Stripping /
   Block-Kontext-Tracking als „vollständig" durchgehen):
   - `Active` + `.env.example`-Block mit **nur kommentiertem**
     `# POSTGRES_PASSWORD=…` ⇒ `actionRepairArtifacts`, Env-Block
     neu (Comment-Stripping greift, nicht-kommentierte Zeile fehlt).
   - `Active` + `service.postgres`-Block mit `POSTGRES_USER:` unter
     `labels:` statt `environment:` ⇒ `actionRepairArtifacts`,
     Service-Block neu (Block-Kontext-Tracking erkennt falschen
     Parent).
   - `Active` + `service.postgres`-Block mit `healthcheck:` aber nur
     `disable: true` darunter ⇒ `actionRepairArtifacts`, Service-Block
     neu (disable-Ausnahme greift; LH-AK-002 erwartet aktiven
     Healthcheck).
   - `Active` + `service.postgres`-Block mit `image:` ohne Wert
     (z. B. `image:` gefolgt von Leerzeile oder
     `image: ""`) ⇒ `actionRepairArtifacts`, Service-Block neu.
   - `Active` + `service.postgres`-Block mit kommentierter
     `# image: postgres:16` und keine andere image-Zeile ⇒
     `actionRepairArtifacts`, Service-Block neu.
   - **Positiver Gegen-Test:** `Active` + `service.postgres`-Block
     mit allen Pflichtfeldern + Inline-Trailing-Kommentar
     (`POSTGRES_PASSWORD: x  # CHANGEME after first deploy`) ⇒
     echter No-op (Pre-Kommentar-Teil enthält den Key; Match
     erfolgreich).

   *Architektur-Hinweis:* Diese Prüfung wird bewusst **nicht** als
   siebter `ServiceState` modelliert; LH-FA-ADD-005 beschreibt nur
   die Doppel-Add-State-Machine. Fehlende/incomplete Volume-/Env-/
   Service-Artefakte sind PostgreSQL-spezifische LH-FA-ADD-002-
   Repair-Fälle und gehören damit zum Service-Plan, nicht zur
   State-Klassifikation.

   *Response-Vertrag-Update:*
   - T3 garantierte (Code) bisher `Changed=nil` für jeden
     Active-Pfad. T4c relaxiert das: `PriorState=Active` mit
     `Changed!=nil` ist beim Artefakt-Repair erlaubt.
   - Der Doc-Comment in `port/driving/addservice.go` ist bereits auf
     den T4-Vertrag aktualisiert (Review-Fix vom T3-Review). T4c
     entfernt nur den `Changed=nil`-only-Test
     (`TestAdd_ActiveIsNilErrorNoOp`) bzw. spaltet ihn in:
     `Active-fully-consistent → no-op` und
     `Active-with-missing-artifact → repair`.

   *`.env.example`-Strategie:* Hash-Managed-Block mit Marker
   `service.postgres`. Fälle:
   - **create**: Datei fehlt → mit dem Block neu anlegen.
   - **append**: Datei vorhanden, kein `service.postgres`-Block →
     Block ans Dateiende anhängen, mit genau einer trennenden
     Leerzeile falls die Datei nicht leer ist.
   - **replace**: Block existiert wohlgeformt → per
     `managedblock.Replace` deterministisch ersetzen.
   - **malformed-abort**: Block existiert malformed → Add aborten
     mit `ErrServiceInconsistent`; **keine** der drei Zieldateien
     darf geschrieben werden (Plan-Fail-Pflicht).

   *Compose-Patches:*
   - `s.yaml.PatchMappingEntryYAML(composeBody, "services",
     "postgres", composeBlockBody, "service.postgres")` für den
     Service-Block.
   - `s.yaml.PatchMappingEntryYAML(composeBody, "volumes",
     "postgres-data", volumeBlockBody, "volume.postgres")` für den
     Volume-Block.
   - LH-FA-ADD-005-State-Detection bleibt am
     `service.postgres`-Marker; der Volume-Marker ist Teil des
     LH-FA-ADD-002-Write-Pfads und wird bei jedem mutierenden Add
     deterministisch mitgeschrieben.

   *Plan-and-Execute-Garantie:* T4c hält den M3-Split ein. Alle
   pre-write-Fehler (Render, Parse, malformed Marker, unsupported
   fragment shape, Patch-Validation) feuern in der Plan-Phase und
   garantieren: keine der drei Zieldateien wird geschrieben.
   Write-Fehler während der Execute-Phase können mit dem heutigen
   `driven.FileSystem.WriteFile`-Port nicht transaktional über
   mehrere Dateien zurückgerollt werden; T4c pinnt deshalb explizit
   mit Tests, dass vor dem ersten Write alle Zielinhalte validiert
   sind. Weitergehende Multi-File-Atomicity ist keine M5-Anforderung
   und bekommt keinen Carveout; falls produktfachlich gefordert,
   startet sie als neuer Slice statt als M5-Restschuld.

   *Tests (T4c-Scope):*
   - Templates: Render-Roundtrip für alle drei (`postgres.compose`,
     `postgres.volume`, `postgres.env`); `templateNames()`-Test
     mit dem erweiterten Embed-Set.
   - `executeAdd` für `actionRegister` (Unregistered): alle drei
     Zieldateien geschrieben; `compose.yaml` enthält genau ein
     `service.postgres`- und genau ein `volume.postgres`-Block,
     beide außerhalb des `init`-Blocks; `.env.example` enthält
     den `service.postgres`-Block; `u-boot.yaml` enthält
     `services.postgres.enabled: true` und keine Kommentar-
     Verluste an den Top-Level-Feldern.
   - `executeAdd` für `actionReactivate` (Deactivated + EnabledUnset):
     `enabled` wird via `PatchScalar` umgeschaltet/gesetzt; Compose-
     Block wird neu erzeugt; `.env.example` wird ggf. ergänzt.
   - `executeAdd` für `actionRebuildBlock` (InconsistentBlock):
     Compose-Block wird via `PatchMappingEntryYAML` neu eingefügt;
     mit und ohne vorhandene `compose.yaml` (Bootstrap-Pfad).
   - **Symlink-Reject-Tests** (pinnt `loadForPatch`-Sicherheitskante,
     analog M3-T4b):
     - `compose.yaml` ist ein Symlink ⇒ Abort mit
       `driving.ErrBackupUnsupportedKind`; **kein** Write an irgendeine
       der drei Zieldateien.
     - `.env.example` ist ein Symlink ⇒ identisch.
     - `u-boot.yaml` ist ein Symlink ⇒ identisch (defensive; in der
       Praxis selten, aber Vertrag muss gelten).
     - Non-regular file (Mode-Bits zeigen weder regular noch
       Symlink; im Test ein neuer `fakeFS`-Helper analog zu
       `markSymlink`) ⇒ identisch.
   - **Mode-Preservation-Test:** seed'e `compose.yaml` mit Mode
     `0o600`. Nach `executeAdd` muss die geschriebene Datei wieder
     `0o600` haben — nicht `defaultFileMode`. Symmetrisch für
     `.env.example`.
   - **u-boot.yaml-TOCTOU-Test:** simuliere die Race über
     `fakeFS.failLstatOn` für den Plan-Phase-`Lstat`-Aufruf mit
     `iofs.ErrNotExist`. `executeAdd` muss mit
     `driving.ErrProjectNotInitialized` aborten und **nicht** eine
     halbierte `u-boot.yaml` schreiben (`fakeFS.writtenPaths()`
     enthält keinen u-boot.yaml-Eintrag).
   - **Bootstrap-Trigger-Tests** (pinnt die enge Bootstrap-Definition):
     - `compose.yaml` fehlt komplett ⇒ Bootstrap erzeugt Split-Block-
       Compose; Patches darauf; resultierende Datei hat init-Block +
       `service.postgres` + `volume.postgres`.
     - `compose.yaml` existiert mit nur Whitespace ⇒ identisch zu
       fehlend (Bootstrap zieht).
     - **User-Compose ohne `init`-Marker bleibt erhalten (block-style):**
       seed'e `compose.yaml` mit User-Inhalt als block-style mapping:
       ```yaml
       services:
         mywebapp:
           image: nginx
       networks:
         my-net: {}
       ```
       (plus Kommentare, kein `# BEGIN ... init`-Marker). `executeAdd`
       darf **nicht** das Scaffold drüberkopieren; stattdessen werden
       die u-boot-Marker via `PatchMappingEntryYAML` in die existierende
       Datei eingefügt, der User-Service `mywebapp` und der
       User-`networks`-Block bleiben byte-identisch.
     - **User-Compose mit nicht-leerem Flow-Style-`services:` ⇒
       fachlicher Fehler:** seed'e `compose.yaml` mit
       `services: { mywebapp: {image: nginx} }` (nicht-leeres
       Flow-Mapping). T4c muss mit `ErrServiceInconsistent` aborten
       (Repair-Hint: „convert services: to block style") und keine
       Datei schreiben — der T4b-Adapter feuert hier
       `ErrYAMLFragmentInvalid`, T4c übersetzt zu fachlichem
       Sentinel.
     - User-Compose ohne `init`-Marker UND ohne `services:`-Key:
       `PatchMappingEntryYAML` legt `services:` als neuen Top-Level-
       Key am Datei-Ende an; User-Inhalt davor unverändert.
   - `executeAdd` für `actionRepairArtifacts`: nur das fehlende
     Artefakt wird geschrieben (Volume oder Env, oder beide).
     `Changed` enthält genau die geschriebenen Pfade.
   - Active + alle Artefakte komplett → echter No-op, `Changed=nil`,
     kein Write am Fake-FS.
   - Active + malformed Volume-Block → `ErrServiceInconsistent`.
   - Active + malformed Env-Block → `ErrServiceInconsistent`.
   - Active + `service.postgres`-Marker am falschen Ort (unter
     `volumes:` statt `services.postgres`) → `ErrServiceInconsistent`.
   - Active + `volume.postgres`-Marker am falschen Ort (unter
     `services:` statt `volumes.postgres-data`) →
     `ErrServiceInconsistent`. Symmetrisch zum Service-Wrong-Anchor-
     Test; pinnt die strukturelle Anker-Prüfung für beide Marker
     unabhängig voneinander.
   - **User-Manual-Entry-ohne-Marker (action-agnostisch):** seed'e
     `compose.yaml` mit einer User-manuellen `services.postgres:`-
     Definition (z. B. ein `image: my-fork/postgres:custom`, kein
     u-boot-Marker als Sub-Element) und `u-boot.yaml` mit
     `services.postgres.enabled: true`. T3 klassifiziert das als
     `InconsistentBlock` (entry da, Marker fehlt). T4c muss vor
     jedem `PatchMappingEntryYAML`-Aufruf via Pre-Patch-Anker-Check
     (Richtung Eintrag → Marker) mit `ErrServiceInconsistent`
     aborten — **nicht** den User-Block überschreiben und **nicht**
     einen zweiten `postgres:`-Key im Mapping anlegen. Symmetrische
     Tests für `volumes.postgres-data:` mit User-Manual-Inhalt
     plus alle drei mutierenden States (`Unregistered`,
     `Deactivated`, `EnabledUnset`).
   - **Wrong-Anchor in non-Active-States** (pinnt den
     action-agnostischen Pre-Patch-Anker-Check): seed'e
     `Deactivated` (`enabled: false`) mit einem wohlgeformten
     `service.postgres`-Marker unter `volumes:`. `u-boot add postgres`
     muss mit `ErrServiceInconsistent` aborten — nicht reaktivieren
     und nicht ein zweites `services.postgres`-Mapping anlegen.
     Analoge Tests für `EnabledUnset` und `InconsistentBlock`
     (jeweils ein Test mit einem fehl-verankerten Marker pro
     State); jeder Pfad muss vor jedem Write abbrechen.
   - **Env-Wrap-Vertrag-Test**: `renderEnvManagedBlock(postgres,
     varsBody)`-Output muss exakt mit `# BEGIN U-BOOT MANAGED
     BLOCK: service.postgres\n` beginnen und mit `# END U-BOOT
     MANAGED BLOCK: service.postgres\n` enden;
     `managedblock.Find(output, …)` muss `varsBody` als Block-Body
     zurückgeben.
   - **Replace-Idempotenz mit Marker-Erhaltung**: zweimal
     `executeAdd` mit `actionRepairArtifacts`-Env-Repair auf
     derselben `.env.example`. Nach Run 1: genau ein Block mit
     BEGIN+END. Nach Run 2: weiterhin genau ein Block (kein
     zweiter Block am Datei-Ende), Datei byte-identisch zu Run 1.
     Pinnt dass `managedblock.Replace` durch den Wrap-Helper die
     Marker erhält.
   - `.env.example`-Pfade: create / append / replace / malformed-
     abort jeweils explizit.
   - Compose-Schreiben byte-erhalten (gemäß T4b-Byte-Erhaltungs-
     Vertrag): seed'e `compose.yaml` mit (a) Kommentaren oberhalb
     der Top-Level-Maps, (b) einem User-Custom-Service unter
     `services:` außerhalb jedes u-boot-Markers, (c) einem
     User-Top-Level-Block `networks:` mit Custom-Inhalt.
     Assertiere nach `executeAdd`: User-Service / -Networks-Block /
     User-Kommentare bleiben byte-identisch; nur der
     `service.postgres`- / `volume.postgres`-Markerbereich und ggf.
     die Parent-Header-Form sind verändert.
   - Pre-Write-Fail-Garantie: erzwinge einen Render-Fehler im
     Template-Mock oder einen Patch-Fehler im Fake; assertiere dass
     `fakeFS.writtenPaths()` keine der drei Zieldateien enthält.
   - Yaml-Adapter-End-to-End: ein Integrationstest mit dem **echten**
     yaml-Adapter (statt Fake) prüft, dass `managedblock.Find(...,
     "service.postgres")` und `managedblock.Find(..., "volume.postgres")`
     auf dem produzierten Output erfolgreich sind. Parsebarkeit
     allein reicht nicht — die T3-State-Detection findet später
     exakt diese Marker wieder.

   *Carveouts nach T4c:*
   - Der T3-Carveout-Eintrag (`executeAdd`-Stub +
     `actionRepairArtifacts`-Enum ohne Caller) wird entfernt — beide
     sind jetzt voll implementiert.
   - Keine neuen temporären Carveouts erwartet. Falls die
     adapter-lokale Marker-Scanner-Logik (Trade-off aus T4b) in der
     Praxis stört, kommt sie als V1-Slice; T4c fügt dafür keinen
     Vorrat-Eintrag.

   *DoD T4c:*
   - Alle obigen Tests grün; `make gates` grün.
   - Coverage-Threshold (90 %) gehalten — `executeAdd` /
     `detectActiveArtifacts` getestet, Templates über Render-
     Roundtrip abgedeckt.
   - DoD-Line: `T4c ✅ <commit-hash>` (+ Review-Fix-Commit
     falls anfallend).
   - Slice-Plan §T3 wird **nicht** mehr angefasst (Historie
     korrekt durch Sub-Tranchen-Hashes); §T4c kann den
     Carveout-Entfernungspunkt referenzieren.

5. **T5 — LH-FA-ADD-005-State-Machine-Tests.**
   - End-to-end-Tests für jede State-Transition (mit fake FS +
     fake yaml-codec):
     - unregistered → active (neu-anlegen).
     - active → active (nil-error no-op,
       `PriorState=Active`, `State=Active`, `Changed=nil`) nur bei
       vollständig vorhandenem Service-, Volume- und Env-Block.
     - deactivated → active (re-aktivieren, Compose-Block neu).
     - inconsistent-yaml → ErrServiceInconsistent (Abort).
     - malformed managed compose-block → ErrServiceInconsistent
       (Abort, Code 10 über CLI-Mapping).
     - inconsistent-block → active (Compose-Block-Rebuild).
     - enabled-key-fehlt → treated as deactivated (Add re-aktiviert).
     - active + fehlender `volume.postgres`-Block → active
       (`actionRepairArtifacts`, Volume-Rebuild, kein Service-Duplikat).
     - active + fehlender `.env.example`-Block `service.postgres` →
       active (`actionRepairArtifacts`, Env-Rebuild).
     - active + malformed `volume.postgres`- oder Env-Block →
       `ErrServiceInconsistent` (Abort).
     - disabled/unset + malformed `service.postgres`-Block →
       `ErrServiceInconsistent` (Abort vor Reaktivierung).
   - Idempotenz: zweimal `u-boot add postgres` produziert
     identischen finalen Zustand (zweite Invocation =
     nil-error no-op mit `Changed=nil`).
   - Re-Init-Regression: `u-boot init`-Scaffold → `u-boot add
     postgres` → `u-boot init --force` darf die `service.postgres`-
     und `volume.postgres`-Blöcke in `compose.yaml` nicht entfernen.
     Der Test pinnt damit die T4a-Entscheidung, Add-on-Blöcke außerhalb
     des `init`-Blocks zu platzieren.
   - Missing-file-Recovery: vorhandene `u-boot.yaml` mit
     `services.postgres.enabled: true`, aber fehlende `compose.yaml`,
     muss nach `u-boot add postgres` wieder eine parsebare
     `compose.yaml` mit Service- und Volume-Block ergeben.
   - Env-Idempotenz: zweimal `u-boot add postgres` darf in
     `.env.example` genau einen `service.postgres`-Managed-Block
     hinterlassen; bestehende manuelle Zeilen außerhalb des Blocks
     bleiben erhalten.

6. **T6 — CLI-Subcommand `u-boot add <service>`.**
   - Cobra-Struktur: `add <service>` als Positional-Arg auf einem
     `add`-Command, **nicht** `add postgres` als fest verdrahtetes
     Sub-Subcommand. Dadurch läuft auch ein syntaktisch valides, aber
     nicht unterstütztes `u-boot add redis` durch
     `domain.ServiceName` + Katalogcheck und endet mit
     `ErrServiceUnsupported` → Exit-Code 10 statt Cobra-Usage-Code 2.
   - Persistente Flags `--yes`/`--no-interactive` lesen (heute
     NoOp für postgres-MVP, weil keine Abhängigkeits-Prompts).
   - Exit-Code-Mapping über bestehende ExitCode()-Funktion
     (Validation-Sentinels inklusive `ErrServiceInconsistent` → 10).
   - Output: knapper Summary-Bericht
     (`Created services.postgres.enabled: true`,
     `Wrote managed block in compose.yaml`, `.env.example`).

7. **T7 — Closure: doctor-Integration + Slice-Move.**
   - Doctor-Check `services.<name>.enabled` muss explizit gesetzt
     sein (LH-FA-ADD-005 §893): neuer Check `services.enabled-key`
     mit warn-Severity bei unset.
   - Doctor-Check für die in M4-T6 deferred
     `forwardPorts`-Konsistenz: jetzt mit dem neuen services-Schema
     umsetzbar.
   - Doctor-Check für die in M4-T6 deferred
     `devcontainer.enabled`-Severity-Eskalation: braucht
     `devcontainer:`-Block in u-boot.yaml-Schema — entweder in
     diesem T7 nachziehen oder als Carveout für eigenen Slice
     verschieben.
   - End-to-end-Smoke (`docker run u-boot init demo && u-boot add
     postgres && u-boot doctor`) muss grün laufen.
   - slice-m5-add-postgres.md nach `done/`. Roadmap M5 → Done.

> **Archiv-Hinweis zum Smoke-Punkt (M5-Done-Review):** Im distroless
> Runtime-Image (`make build`, distroless/static-nonroot) fehlen
> `git` und `docker` per Design (LH-NFA-PORT-002 — minimale
> Host-Deps). `u-boot doctor` meldet diese Werkzeuge dann als
> Error (`git.installed`, `docker.installed`). Die `u-boot init`-
> und `u-boot add postgres`-Schritte des Smoke laufen im Runtime-
> Image, ebenso jeder Projektcheck (`uboot.yaml.valid`,
> `compose.yaml.valid`, `services.enabled-key`, …). Der Smoke gilt
> als grün, wenn die Projektchecks grün sind; die Werkzeug-Checks
> sind erwartet rot im distroless-Container und werden auf einem
> Dev-System mit installiertem git+docker grün. Eine voll grüne
> Smoke-Variante im Runtime-Image würde einen FAT-Container brauchen
> (z. B. mit `docker:cli` Multi-Stage), was dem distroless-Ziel
> widerspricht — daher bewusst out of scope.

## Akzeptanzkriterien (Slice-Level)

- `LH-FA-ADD-001`, `LH-FA-ADD-002`, `LH-FA-ADD-005` abgehakt.
- `make gates` grün.
- M4-T6-deferred-Carveouts (forwardPorts + devcontainer.enabled)
  in T7 entweder geschlossen oder explizit auf eigene Folge-
  Slices verschoben.

## Out of Scope

- **Keycloak / OTel-Add-Ons** (V1).
- **`--with-deps` + Add-on-Abhängigkeiten** (V1).
- **`u-boot remove <service>`** (V1).
- **Custom-Services** (eigene Templates per `u-boot template`-Slice).

## Bezug

- Auslösende Spec: `LH-FA-ADD-001..002`, `LH-FA-ADD-005`
  (`spec/lastenheft.md` §4.5).
- Vorgänger: [`slice-m4-doctor`](../done/slice-m4-doctor.md) hat
  zwei explizit-deferred Concerns (forwardPorts + devcontainer-
  Severity), die mit T7 hier zur Auflösung kommen können.
- Nachfolger: MVP-Closure-Slice (LH-AK-001..002) ist der
  Acceptance-Demo-Pfad `mkdir demo && cd demo && u-boot init &&
  u-boot add postgres && u-boot doctor`.
