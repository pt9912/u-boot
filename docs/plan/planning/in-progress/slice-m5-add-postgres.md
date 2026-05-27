# Slice M5: `u-boot add postgres`-Flow

> **Status:** In progress
> **DoD:** T1 offen / T2 offen / T3 offen / T4 offen / T5 offen / T6 offen / T7 offen

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

Pro Service-Name gibt es vier beobachtbare Zustände beim Add-Versuch:

| Zustand                | `services.<name>` in u-boot.yaml | `enabled` | Managed-Block in compose.yaml | Add-Aktion                                                                              |
| ---------------------- | -------------------------------- | --------- | ----------------------------- | --------------------------------------------------------------------------------------- |
| **unregistered**       | fehlt                            | —         | fehlt                         | Neu anlegen: services-Eintrag + Compose-Block + .env.example-Block (LH-FA-ADD-002).     |
| **active**             | vorhanden                        | `true`    | vorhanden                     | No-op (idempotent); deutlicher Hinweis dass Service schon aktiv ist.                    |
| **deactivated**        | vorhanden                        | `false`   | irrelevant                    | Re-Aktivierung: `enabled: true` + Compose-Block + .env-Block neu erzeugen.              |
| **enabled-key-fehlt**  | vorhanden                        | (unset)   | irrelevant                    | (Doctor-warn-Pfad; Add interpretiert als deactivated und re-aktiviert wie oben.)        |
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
   - `AddServiceRequest` (BaseDir, ServiceName domain.ServiceName,
     plus die persistenten Mode-Flags: AssumeExisting/NoInteractive/
     Force? — Force vermutlich nicht für add).
   - `AddServiceResponse` (resulting state, geänderte Pfade).
   - Sentinels:
     - `ErrServiceUnsupported` (Service-Name nicht im
       built-in-Katalog, heute nur „postgres").
     - `ErrServiceAlreadyActive` (no-op-Pfad mit Hinweis).
     - `ErrServiceInconsistent` (LH-FA-ADD-005-inconsistent-yaml-Fall).
     - `ErrProjectNotInitialized` (kein `u-boot.yaml` → LH-FA-ADD-001).
     Mapping zu LH-FA-CLI-006-Exit-Codes (vermutlich 10 für
     validation, ggf. ein eigener Code 13 für „project-state").

3. **T3 — Application-Service-Skeleton + State-Detection.**
   - `AddServiceService` orchestriert
     `FileSystem`/`YAMLCodec`/`Logger`/(`managedblock`).
   - `detectServiceState(baseDir, name)` liest `u-boot.yaml` +
     `compose.yaml` (managed-block-Marker), klassifiziert den
     State.
   - `Add(ctx, req)`: dispatcht je nach State auf no-op / re-aktivieren /
     neu-anlegen / inconsistent-abort. Plan-and-Execute-Split
     analog M3-T4b: erst alle geplanten File-Edits sammeln, dann
     ausführen — Plan-Fehler verhindert jeden Side-Effect.
   - Tests für die State-Detection (jeder der 6 States hat einen
     Fake-FS-Setup).

4. **T4 — PostgreSQL-Templates + Write-Pfad.**
   - `application/templates/services/postgres.compose.tmpl`:
     Service-Block mit
     `image: postgres:16-alpine`, `environment` (POSTGRES_USER /
     POSTGRES_PASSWORD / POSTGRES_DB), `volumes` (named-volume),
     `ports: ["5432:5432"]`, `healthcheck` (`pg_isready`).
   - `application/templates/services/postgres.env.tmpl`:
     `POSTGRES_USER=postgres`,
     `POSTGRES_PASSWORD=CHANGEME_POSTGRES_PASSWORD`,
     `POSTGRES_DB=postgres`. (Sicherheits-Convention: explizit
     `CHANGEME_*` als Default; nie reale Defaults.)
   - Template-Loading via `embed.FS` (analog M3-T4b).
   - Managed-Block-Marker `BEGIN/END U-BOOT MANAGED BLOCK:
     service.postgres` (analog `init`-Marker aus LH-SA-FILE-002).
   - u-boot.yaml-Patch: `services.postgres.enabled: true` einfügen
     (managed-block-marker-style, weil u-boot.yaml als
     whole-file-managed gilt — möglicherweise braucht es eine
     Schema-präzise YAML-Manipulation statt eines string-managed-
     blocks; Entscheidung im T4-Slice).

5. **T5 — LH-FA-ADD-005-State-Machine-Tests.**
   - End-to-end-Tests für jede State-Transition (mit fake FS +
     fake yaml-codec):
     - unregistered → active (neu-anlegen).
     - active → ErrServiceAlreadyActive (no-op).
     - deactivated → active (re-aktivieren, Compose-Block neu).
     - inconsistent-yaml → ErrServiceInconsistent (Abort).
     - inconsistent-block → active (Compose-Block-Rebuild).
     - enabled-key-fehlt → treated as deactivated (Add re-aktiviert).
   - Idempotenz: zweimal `u-boot add postgres` produziert
     identischen finalen Zustand (zweite Invocation =
     ErrServiceAlreadyActive).

6. **T6 — CLI-Subcommand `u-boot add <service>`.**
   - Cobra-Sub-Subcommand-Struktur: `add` als Parent-Command (für
     spätere Add-Ons in V1), heute mit nur `postgres`-Argument.
     Alternative: `add <service>`-Positional-Arg statt
     Sub-Subcommand-Struktur, ServiceName-Validierung an
     domain.ServiceName.
   - Persistente Flags `--yes`/`--no-interactive` lesen (heute
     NoOp für postgres-MVP, weil keine Abhängigkeits-Prompts).
   - Exit-Code-Mapping über bestehende ExitCode()-Funktion
     (Validation-Sentinels 10, ErrServiceInconsistent ggf. 13?).
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
