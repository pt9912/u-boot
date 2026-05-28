# Slice M5: `u-boot add postgres`-Flow

> **Status:** In progress
> **DoD:** T1 ✅ `995726a` / T2 ✅ `f054986` / T3 offen / T4 offen / T5 offen / T6 offen / T7 offen

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
   - **Compose-Init-Block-Kompatibilität zuerst lösen:** Die aktuelle
     M3-`compose.yaml.tmpl` enthält `services: {}` innerhalb des
     `BEGIN/END ... init`-Blocks. Service-/Volume-Add-on-Blöcke dürfen
     dort **nicht** eingefügt werden, weil ein späteres
     `u-boot init --force` den kompletten `init`-Block ersetzt und
     sonst alle Add-on-Blöcke löscht.
     T4 ändert deshalb die Compose-Scaffold-Struktur so, dass nur die
     init-eigene Basiskonfiguration im `init`-Block liegt; die
     add-on-veränderlichen Top-Level-Maps `services:` und `volumes:`
     liegen außerhalb des `init`-Blocks.
   - Der bestehende Init-Replace-Pfad muss dafür angepasst werden:
     `InitProjectService.executeReplaceBlock` darf bei Templates mit
     Content außerhalb des `init`-Blocks nicht mehr das komplette
     gerenderte Template in den alten Block splicen. T4 führt einen
     Helper ein (z. B. `renderManagedBlockOnly` /
     `extractManagedBlock`) und ersetzt bei `actionReplaceBlock` nur
     den gerenderten `BEGIN/END ... init`-Bereich. Dadurch kann
     `compose.yaml.tmpl` außerhalb des `init`-Blocks Top-Level-Maps
     enthalten, ohne dass `u-boot init --force` `services:`/`volumes:`
     dupliziert oder in den `init`-Block verschiebt.
     Regressionstest: bestehende Split-Compose-Datei mit
     `service.postgres`/`volume.postgres` → `u-boot init --force` →
     genau ein `init`-Block, genau ein Top-Level-`services:`, genau ein
     Top-Level-`volumes:`, Add-on-Marker bleiben außerhalb des
     `init`-Blocks.
   - Für bestehende M3-Projekte mit der alten Compose-Scaffold führt
     der Add-Write-Pfad vor dem Service-Patch eine einmalige
     Normalisierung aus: den vorhandenen `init`-Block durch die neue
     init-Block-Form ersetzen (ohne `services: {}` darin), danach
     fehlende Top-Level-Maps `services:`/`volumes:` außerhalb jedes
     Managed-Blocks strukturiert anlegen. Erst danach werden
     `service.postgres` und `volume.postgres` gepatcht. T4 testet
     explizit, dass beide Add-on-Marker außerhalb des Bytebereichs des
     `init`-Blocks liegen.
   - Fehlt `compose.yaml` trotz vorhandener `u-boot.yaml`, bootstrapped
     `executeAdd` eine minimale Compose-Datei in der neuen
     Split-Block-Form (init-Block + leere `services:`/`volumes:`-Maps)
     und führt anschließend dieselben strukturierten Patches aus.
     Damit ist der `enabled: true`, Compose-Block-fehlt-Recovery-Pfad
     aus `InconsistentBlock` auch bei komplett fehlender Datei
     deterministisch.
   - `application/templates/services/postgres.compose.tmpl`:
     YAML-Fragment für den Wert unter `services.postgres` mit
     `image: postgres:16-alpine`, `environment` (POSTGRES_USER /
     POSTGRES_PASSWORD / POSTGRES_DB), `volumes` (named-volume),
     `ports: ["5432:5432"]`, `healthcheck` (`pg_isready`).
     Das Template enthält **keinen** Top-Level-`services:`-Key.
     Grund: die Einfügung passiert strukturiert unter der außerhalb
     des `init`-Blocks liegenden Top-Level-Map `services:`; ein
     angehängter zweiter Top-Level-Key wäre ungültig bzw.
     parserabhängig überschreibend.
     Die `environment`-Werte referenzieren die `.env.example`-Keys
     über Compose-Interpolation, z. B.
     `POSTGRES_USER: ${POSTGRES_USER}`,
     `POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}`,
     `POSTGRES_DB: ${POSTGRES_DB}`. Damit sind die `.env.example`-
     Einträge nicht nur Dokumentation, sondern die Quelle für die
     lokale Laufzeitkonfiguration.
   - `application/templates/services/postgres.volume.tmpl`:
     YAML-Fragment für den Wert unter `volumes.postgres-data`; für das
     MVP reicht ein leeres Mapping `{}`. Das Fragment enthält keinen
     Top-Level-`volumes:`-Key und wird unter `volumes:` mit dem Marker
     `volume.postgres` eingefügt.
   - `application/templates/services/postgres.env.tmpl`:
     `POSTGRES_USER=postgres`,
     `POSTGRES_PASSWORD=CHANGEME_POSTGRES_PASSWORD`,
     `POSTGRES_DB=postgres`. (Sicherheits-Convention: explizit
     `CHANGEME_*` als Default; nie reale Defaults.)
   - Template-Loading via `embed.FS` (analog M3-T4b). Weil die Service-
     Templates in einem Unterverzeichnis liegen, muss die Embed-
     Directive erweitert werden:
     `//go:embed templates/*.tmpl templates/services/*.tmpl`.
     `renderTemplate("services/postgres.compose.tmpl", ...)` bleibt
     mit dem bestehenden `"templates/"+name`-Pfadmodell kompatibel.
   - Compose-Schreibstrategie: kein blindes String-Append ans
     Dateiende, keine doppelten Top-Level-Keys und kein vollständiges
     Re-Marshal der bestehenden `compose.yaml`. T4 erweitert den
     YAML-Port um einen byte-erhaltenden, strukturiert validierten
     Mapping-Patch für YAML-Fragmente:
     `PatchMappingEntryYAML(content []byte, path []string, key string,
     valueYAML []byte, markerName string) ([]byte, error)`. Die
     Application rendert `postgres.compose.tmpl` als YAML-Fragment und
     übergibt die Bytes; sie importiert weiterhin kein `yaml.v3`.
     Der yaml-Adapter parst `content` und `valueYAML` mit `yaml.v3`
     zur Validierung; die Bytebereiche für Parent-Map und Managed-Block
     werden über einen begrenzten YAML-Zeilenscanner bestimmt. Der
     Adapter importiert **nicht** `internal/hexagon/application/managedblock`,
     weil die Depguard-Regel `adapter-no-application`
     Adapter→Application-Imports verbietet. Die Scanner-Logik muss die
     bestehenden Hash-Marker-Literale kompatibel zu `managedblock`
     behandeln und bekommt eigene Adapter-Tests für not-found,
     malformed, duplicate-begin und indented-marker-Fälle. Geschrieben
     wird nur der betroffene Parent-Map-Bereich als Byte-Splice:
     - `valueYAML` muss ein Mapping-Root sein (Scalar/Sequence ⇒
       Fehler).
     - Doppelte Top-Level-Keys oder ein vorhandener Parent-Key, der
       weder Mapping noch leeres `{}`/`null` ist, führen zu einem
       technischen Patch-Fehler vor jedem Write.
     - Fehlt der Parent-Key (`services`/`volumes`), wird genau dieser
       Top-Level-Key außerhalb jedes Managed-Blocks angelegt.
     - Ist der Parent-Key als `{}` oder `null` vorhanden, wird nur diese
       Zeile in eine mehrzeilige Mapping-Form umgeschrieben.
     - Existiert der Ziel-Managed-Block, wird nur dessen Bytebereich per
       adapter-lokalem Byte-Splice ersetzt; fehlt er, wird ein neuer,
       eingerückter Hash-Managed-Block unter dem Parent-Key eingefügt.
     - Bytes außerhalb des neu angelegten/ersetzten Parent- oder
       Managed-Block-Bereichs bleiben unverändert; Kommentare und
       manuelle Einträge außerhalb dieses Bereichs dürfen nicht durch
       YAML-Reformatting wandern.
     T4 testet diese Byte-Erhaltung mit einer `compose.yaml`, die
     manuelle Kommentare und Einträge außerhalb der Zielbereiche enthält.
     Zusätzlich pinnt ein Depguard-orientierter Test/Review-Check, dass
     `internal/adapter/driven/yaml` keinen Import auf
     `internal/hexagon/application` enthält.
   - Für PostgreSQL werden zwei Compose-Patches ausgeführt:
     `services.postgres` mit Marker `service.postgres` und
     `volumes.postgres-data` mit Marker `volume.postgres`. Die
     LH-FA-ADD-005-State-Erkennung bleibt bewusst am
     `service.postgres`-Marker; der Volume-Marker ist Teil des
     LH-FA-ADD-002-Write-Pfads und wird bei jedem mutierenden Add
     deterministisch mitgeschrieben.
   - **Aktiver, aber unvollständiger Add-on-Zustand:** Der
     `ServiceStateActive` aus T3 bedeutet nur: YAML-Anker
     `services.postgres.enabled: true` + `service.postgres`-
     Compose-Block sind vorhanden. T4 ergänzt vor dem endgültigen
     No-op eine PostgreSQL-Artefaktprüfung:
     - `volume.postgres`-Marker in `compose.yaml` fehlt ⇒
       `actionRepairArtifacts`, schreibt den Volume-Block neu.
     - `.env.example` fehlt oder enthält keinen `service.postgres`-
       Block ⇒ `actionRepairArtifacts`, erzeugt/ergänzt den Env-Block.
     - `.env.example` enthält einen malformed `service.postgres`-
       Block oder `compose.yaml` enthält einen malformed
       `volume.postgres`-Block ⇒ Abort mit `ErrServiceInconsistent`.
     - Sind Service-, Volume- und Env-Block vorhanden und wohlgeformt
       **und** strukturell am erwarteten Ort verankert
       (`services.postgres` bzw. `volumes.postgres-data`) ⇒ echter
       No-op mit `Changed=nil`.
     - Wohlgeformte Marker am falschen YAML-Ort oder ohne passenden
       Mapping-Eintrag gelten als fachlich inkonsistent und aborten mit
       `ErrServiceInconsistent`; Marker-Existenz allein reicht für
       LH-FA-ADD-002 nicht aus.
     Diese Prüfung wird bewusst nicht als siebter `ServiceState`
     modelliert, weil `LH-FA-ADD-005` nur die Doppel-Add-State-
     Machine beschreibt; fehlende Volume-/Env-Artefakte sind
     PostgreSQL-spezifische `LH-FA-ADD-002`-Repair-Fälle.
     T4 passt die `AddServiceUseCase`-Kommentare und Response-Tests an:
     `Changed=nil` bedeutet echter No-op; `PriorState=Active` allein
     garantiert nicht mehr, dass keine Datei geschrieben wurde.
   - Managed-Block-Marker stehen innerhalb der jeweiligen YAML-Map:
     `# BEGIN/END U-BOOT MANAGED BLOCK: service.postgres` unter
     `services:` und `# BEGIN/END U-BOOT MANAGED BLOCK:
     volume.postgres` unter `volumes:` (analog LH-SA-FILE-002).
     T4 testet explizit, dass die finale `compose.yaml` genau einen
     Top-Level-`services:`-Key und genau einen Top-Level-`volumes:`-
     Key enthält und von `yaml.v3` parsebar ist.
     Zusätzlich testet T4 mit dem echten yaml-Adapter-Output, dass
     `managedblock.Find(..., service.postgres)` und
     `managedblock.Find(..., volume.postgres)` auf der finalen Datei
     erfolgreich sind. Parsebarkeit allein reicht nicht, weil die
     State-Detection später exakt diese Marker wiederfinden muss.
     Ein weiterer Test legt einen wohlgeformten Marker am falschen Ort
     an und erwartet `ErrServiceInconsistent`, damit ein Marker ohne
     tatsächlichen `services.postgres`-/`volumes.postgres-data`-Eintrag
     nicht als vollständig aktiv durchrutscht.
   - `.env.example`-Schreibstrategie: das PostgreSQL-Env-Template wird
     als eigener Hash-Managed-Block mit Marker `service.postgres`
     geschrieben. Existiert `.env.example` nicht, wird sie mit diesem
     Block neu angelegt. Existiert sie ohne PostgreSQL-Block, wird der
     Block ans Dateiende angefügt (mit genau einer trennenden Leerzeile,
     falls die Datei nicht leer ist). Existiert der Block bereits, wird
     er per `managedblock.Replace` deterministisch ersetzt. Ist der
     Block malformed, bricht Add mit `ErrServiceInconsistent` ab; keine
     der drei Zieldateien darf dann geschrieben werden.
     T4 testet alle vier Env-Pfade: create, append, replace,
     malformed-abort.
   - u-boot.yaml-Patch: `services.postgres.enabled: true` einfügen
     via `s.yaml.PatchScalar(yamlBody, []string{"services",
     name.String(), "enabled"}, true)`. Keine Managed-Block-Marker in
     `u-boot.yaml`; keine naive Full-File-Re-Marshal-Strategie. Ziel:
     Kommentare und unbekannte/V1-Felder so weit erhalten, wie die
     yaml.v3-Node-API sie erhält.
   - T4 hält den Plan-and-Execute-Split aus M3 ein: erst alle
     Zielinhalte für `u-boot.yaml`, `compose.yaml` und `.env.example`
     vollständig berechnen, inklusive malformed-Marker-Checks; erst
     wenn alle Berechnungen erfolgreich sind, werden die Dateien
     geschrieben. Diese Garantie gilt für alle Pre-Write-Fehler
     (Parse, Render, malformed Marker, unsupported fragment shape).
     Write-Fehler während der Ausführung können mit dem heutigen
     `driven.FileSystem.WriteFile`-Port nicht transaktional über
     mehrere Dateien zurückgerollt werden; T4 testet deshalb
     explizit, dass vor dem ersten Write alle Zielinhalte validiert
     sind. Weitergehende Multi-File-Atomicity ist keine M5-
     Anforderung und bekommt keinen Carveout; falls sie später
     produktfachlich gefordert wird, startet sie als neue Anforderung
     mit eigenem Slice statt als M5-Restschuld.
   - **T4-DoD für den geänderten Response-Vertrag:** Die
     `AddServiceUseCase`-Kommentare und Response-Tests müssen pinnen,
     dass `Changed=nil` echter No-op bedeutet und `PriorState=Active`
     mit `Changed!=nil` beim Artefakt-Repair erlaubt ist. Der
     bestehende T2-Kommentar "`Changed` leer genau bei
     `PriorState=Active`" darf nach T4 nicht mehr existieren.

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
     Der Test pinnt damit die T4-Entscheidung, Add-on-Blöcke außerhalb
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
