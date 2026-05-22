# Slice M3: `u-boot init`-Flow

> **Status:** In progress
> **DoD:** T1 ✅ `132d1a1` + `f5c784a` / T2 ✅ `aaf4d8d` + `39387b9` / T3 ✅ `937adb1` + `2b6582c` / T4a ✅ `5296671` + `ecb8379` / T4b ✅ `077c3e5` + `4d07542` / T4c ✅ `2925471` + Review-Fix-Commit folgt / T5 offen

## Auslöser

Bis M2d ist u-boot ein Bootstrap-Skelett (Build-Infrastruktur, leere
hexagonale Schichten, CLI mit `--help`/`--version`-Stub). M3 liefert
den ersten fachlichen Use-Case: `u-boot init` erzeugt eine Projekt-
struktur (`LH-FA-INIT-001..007`, `LH-FA-CONF-001..003`).

Mit M3 entstehen die ersten produktiven Pakete unter `internal/`,
wodurch zwei M2-Carveouts automatisch greifbar werden:

- `slice-m3-coverage-threshold-aktivieren.md` (Schwellwert von 0 auf
  80 heben).
- `slice-m3-depguard-aktivierung-verifizieren.md` (alle 8
  depguard-Regelblöcke real verifizieren).

## Lieferumfang (MVP-Pflicht-Set)

`u-boot init [<name>] [--devcontainer] [--no-git] [--backup] [--force]
[--assume-existing]`

Pflicht-Verhalten (Lastenheft-Verweis):

- `LH-FA-INIT-001` Befehl `u-boot init`.
- `LH-FA-INIT-002` Projektname: explizit oder aus Arbeitsverzeichnis
  abgeleitet + normalisiert.
- `LH-FA-INIT-003` Projektstruktur: `docker/`, `scripts/`, `docs/`,
  `README.md`, `CHANGELOG.md`, `compose.yaml`, `.env.example`,
  `u-boot.yaml`, `.gitignore`.
- `LH-FA-INIT-004` Bestehendes Projekt erkennen.
- `LH-FA-INIT-005` Überschreibschutz mit `--backup`/`--force`.
- `LH-FA-INIT-006` Projektnamen-Validierung
  (`^[a-z]([a-z0-9-]{0,61}[a-z0-9])?$`).
- `LH-FA-INIT-007` Git-Repository-Initialisierung (Default an,
  `--no-git`-Override).
- `LH-FA-CLI-005A` Interaktivität (`--yes`, `--no-interactive`,
  `--assume-existing`).
- `LH-FA-CLI-006` Exit-Codes (`0`/`2`/`10`/`11`/`14`).
- `LH-FA-CONF-001..003` u-boot.yaml mit `schemaVersion: 1`.

Devcontainer-Erzeugung (`--devcontainer`-Flag) folgt in M4 als eigener
Slice — M3 erzeugt die Devcontainer-Dateien **nicht**.

## Tranchen-Schnitt

Vorschlag (jede Tranche eigener Commit, je grün durch alle Gates):

1. **T1 — Domain + Driven Ports + minimale Driven Adapter.** ✅ Done
   (Commit `132d1a1` + Review-Fixes `f5c784a`)
   `internal/hexagon/domain/`: `Project`, `ProjectName` (mit Regex aus
   `LH-FA-INIT-006`). `internal/hexagon/port/driven/`: `FileSystem`,
   `YAMLCodec`, `Git`, `Clock`. `internal/adapter/driven/{fs,yaml,git,
   clock}/`: konkrete Implementierungen. Tests pro Schicht (Domain mit
   Property-Style, Driven-Adapter mit `t.TempDir`/`os/exec` echt).

2. **T2 — Application + Driving Port.** ✅ Done
   (Commit `aaf4d8d` + Review-Fixes folgen)
   `internal/hexagon/port/driving/InitProjectUseCase`.
   `internal/hexagon/application/InitProjectService` orchestriert
   die Driven-Ports; Tests mit Fakes für FileSystem/YAMLCodec/Git.
   *Bewusste Lücke:* `LH-FA-INIT-004` Soft-Existing-Detection
   (≥3 Strukturelemente + `--assume-existing`) liegt in
   [`open/slice-m4-soft-existing-detection.md`](../open/slice-m4-soft-existing-detection.md).

3. **T3 — Driving Adapter CLI + Wiring → erster lauffähiger
   `u-boot init`.** ✅ Done (dieser Commit)
   `internal/adapter/driving/cli/`: Cobra-basiertes `init`-Command
   plus Exit-Code-Mapping; CLI-Framework-Wahl in ADR-0005
   verankert (löst gleichzeitig `slice-m3-cli-framework-adr` auf).
   `cmd/uboot/main.go`: Wiring aller Schichten. End-to-End-Smoke
   verifiziert: `docker run … init demo-project --no-git` erzeugt
   die LH-FA-INIT-003-Mindeststruktur + `u-boot.yaml` (LH-FA-CONF-002).

4. **T4 — Überschreibschutz + nicht-interaktive Modi.**
   Wird in drei Sub-Tranchen geliefert, weil die Spec-Anforderungen
   (`LH-FA-INIT-005` Backup-Konvention + Managed-Block-Logik nach
   `LH-SA-FILE-002`, plus `LH-FA-CLI-005A` Modi-Flags) sonst einen
   einzelnen Commit überfrachten:

   - **T4a — Backup-Mechanik.** ✅ Done (`5296671` + Review-Fix-Commit folgt)
     `FileSystem`-Driven-Port um `Lstat` / `Mkdir` /
     `WriteFileExclusive` / `RemoveAll` erweitert (Review hat
     `IsDir` durch das Lstat-basierte Trio ersetzt — Mode-Preservation
     + Symlink-Detection + TOCTOU-Schutz in einem Schwung).
     `application/backup.go` mit `BackupPath` — kleinster-freier-
     Suffix-Algorithmus (`.bak`, `.bak.1`, …) nach `LH-FA-INIT-005`
     §607/608, Files und Verzeichnisse (rekursiv), Rollback bei
     Tree-Backup-Fehler mit `errors.Join` für Sekundär-Fehler;
     Symlink-Rejection (`ErrBackupUnsupportedKind`), 256 MiB-Size-Cap
     (`ErrBackupTooLarge` — temporärer Carveout, Aufhebung in
     [`slice-v1-backup-streaming-copy`](../open/slice-v1-backup-streaming-copy.md)).
     TOCTOU-sichere Top-Level-Reservierung via
     `Mkdir`/`WriteFileExclusive` + Race-Retry-Loop. Sentinels
     `ErrBackupSourceMissing`/`ErrBackupSuffixExhausted`/
     `ErrBackupUnsupportedKind`/`ErrBackupTooLarge` in der
     Driving-Port verankert, `cli.ExitCode` mit Code 14
     (Filesystem-technisch) erweitert. Fakes vollständig überarbeitet
     (Lstat + Modus-Preservation + Symlink-Modellierung +
     Ancestor-Recording bei `WriteFile`/`MkdirAll`). 24 Tests in
     `backup_test.go`. Coverage 92.9 %.

   - **T4b — Managed-Block-Parser + Force/Backup-Flow.** ✅ Done (`077c3e5` + Review-Fix-Commit folgt)
     `application/managedblock/` mit Marker-Parser pro Dateityp nach
     `LH-SA-FILE-002` (YAML/.env `#`, Markdown `<!-- -->`, JSONC `//`):
     `Find` / `Has` / `Replace` plus `ErrBlockNotFound` /
     `ErrBlockMalformed`. Regex-basierter Matcher mit `(?m)`-Multiline,
     leading-whitespace-Toleranz (für indentierte Blocks) und
     CRLF-Toleranz.
     `InitProjectRequest` um `Force`/`Backup` erweitert,
     `driving.ErrForceRequiresBackup` (Code 10) neu. Service-Logik:
     literale Spec-Lesart (LH-FA-INIT-005 §605/§609/§617/§619) — plan-
     and-execute split, `planFile` entscheidet pro Datei
     `Write`/`ReplaceBlock`/`OverwriteFull`/`Abort*`, Execute-Phase
     dispatcht. Plan-Fehler verhindern jeden Side-Effect.
     Templates README/CHANGELOG/compose/.env.example mit
     `BEGIN/END U-BOOT MANAGED BLOCK: init`-Wrapper; .gitignore +
     u-boot.yaml whole-file-managed (§611-Liste excl).
     `Response.Backups []BackupAction` mit Original-/Backup-Pfad-
     Tupel. Summary-Output via `io.Writer` am Service-Konstruktor
     (`progress`), `nil` tolerant zu `io.Discard`, in `cmd/uboot`
     auf `os.Stdout` verdrahtet. 19 Parser-Tests + 12 neue Service-
     Tests, Coverage 91.6 %. End-to-End-Smoke: `docker run u-boot
     init smoke-demo --no-git` schreibt managed-block-Marker
     korrekt; Re-Init ohne Flags → Exit-Code 10.

   - **T4c — Modi-Flags + Exit-Codes.** ✅ Done (Commit folgt + Review-Fix-Commit)
     Cobra-Flags `--force` / `--backup` / `--assume-existing` lokal
     am init-Command, `--yes` / `--no-interactive` als persistente
     Root-Flags (LH-FA-CLI-005A — gelten auch für künftige Befehle
     wie `add`/`remove`/`config set`). Konflikt-Check `--yes` +
     `--no-interactive` → `cli.ErrConflictingModeFlags` →
     `isUsageError`-Pfad → Exit-Code 2 nach `LH-FA-CLI-005A` §235.
     `--assume-existing` wird angenommen + validiert (init-only,
     spec §238), als `AssumeExisting bool` in
     `driving.InitProjectRequest` durchgereicht; M3 hat keine
     Soft-Detection (§247 — auf deterministischen Pfaden NoOp), die
     Hard-Marker-Logik aus T2/T4b schützt bereits gegen vorhandene
     u-boot.yaml/compose.yaml/.env.example. Sentinel
     `ErrConflictingModeFlags` lebt in `cli` (nicht in `driving`),
     weil die Modi-Flags die Application-Schicht nie erreichen.
     `printInitSummary` ergänzt um Backups-Sektion. End-to-End-Smoke
     verifiziert: fresh init / Konflikt-Check / `--force` ohne
     `--backup` (Code 10) / `--force --backup` (full re-init mit
     summary + 6 .bak-Dateien). 6 neue CLI-Tests (Flag-Pass-Through,
     Konflikt, --yes/--no-interactive NoOp, Backup-Summary,
     leerer Backup-Fall, init-only-Geltung von --assume-existing).
     ExitCode-Tabellentest um ErrConflictingModeFlags erweitert.

5. **T5 — Cleanup: Carveout-Auflösung.**
   depguard-Verifikation pro Schicht (siehe
   [`slice-m3-depguard-aktivierung-verifizieren.md`](slice-m3-depguard-aktivierung-verifizieren.md)),
   `carveouts.md` aktualisieren, Roadmap M3 = Done.

   *Vorgezogen erledigt:* Coverage-Schwellwert ist nach M3-T1
   direkt auf 90 % gehoben (siehe
   [`../done/slice-m3-coverage-threshold-aktivieren.md`](../done/slice-m3-coverage-threshold-aktivieren.md));
   bleibt aus dem Carveout-Inventar entfernt.

## Akzeptanzkriterien (Slice-Level)

- `LH-AK-001` Minimaler Init-Flow läuft grün (`mkdir demo && cd demo
  && u-boot init && u-boot doctor`). `doctor` ist noch nicht
  implementiert; dieser AK wird mit M4 vollständig erfüllt. M3
  liefert: `u-boot init` läuft und erzeugt die Pflichtstruktur.
- `u-boot init my-service` und `u-boot init` (Name aus Verzeichnis)
  funktionieren beide.
- `u-boot init` zweimal hintereinander ohne Flag → Exit-Code `10`
  (Überschreibschutz).
- `u-boot init --force` ohne `--backup` und ohne managed block →
  Exit-Code `10` mit Hinweis auf `--backup`.
- `u-boot init --backup --force` → schreibt Backup, überschreibt,
  Exit-Code `0`.
- `u-boot init --no-git` → Repo wird nicht initialisiert.
- `make gates` grün; alle 8 depguard-Regeln verifiziert (eine
  Lint-grüne Variante pro Schicht).
- `make coverage-gate` grün gegen den jetzt aktiven Default
  `THRESHOLD=90`.
- `carveouts.md` Eintrag für depguard-leer-Match als gelöst
  markiert (der Coverage-Bootstrap-Eintrag ist bereits in
  M3-T1-Review-Folgecommit aufgehoben).

## Out of Scope

- Devcontainer-Erzeugung (`--devcontainer`-Flag): eigener Slice
  M4 (`LH-FA-DEV-001..005`).
- `u-boot doctor`: eigener Slice M4 (`LH-FA-DIAG-*`).
- Service-Add-ons (`u-boot add postgres` etc.): eigener Slice M5+.
- `--dry-run`/`--diff`-Flags (`LH-FA-CLI-007/008`, V1).
- JSON-Output (`LH-NFA-USE-004`, V1).
- Template-System (`LH-FA-TPL-*`, V1).

## Bezug

- Auslösende Spec: `LH-FA-INIT-001..007`, `LH-FA-CONF-001..003`,
  `LH-FA-CLI-005A`, `LH-FA-CLI-006`.
- Hängt von: M2d (Carveout-Disziplin etabliert).
- Löst auf: zwei M3-Carveouts in
  [`carveouts.md`](carveouts.md).
- Wird ggf. auslösen: `slice-v1-gomodguard-rules.md` (sobald Cobra
  und yaml.v3 in `go.mod` landen).
