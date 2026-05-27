# Slice V1: Streaming Backup-Copy (256 MiB-Cap aufheben)

> **Status:** Done
> **DoD:** Commit `5715f4f`

## Auslöser

Die T4a-Review (Finding #6) hat die Speicherspitze der Backup-Mechanik
aus `LH-FA-INIT-005` adressiert: `BackupPath` rief `fs.ReadFile`+
`fs.WriteFile` auf, lud also den vollständigen Dateiinhalt in den
Prozess-Heap. Bei einem multi-GB-Asset im Backup-Scope (`docs/`,
`docker/`) hätte das den Prozess gekippt. Als MVP-Kompromiss war ein
harter Cap eingezogen (`maxBackupFileSize = 256 << 20`); Dateien
darüber wurden mit `ErrBackupTooLarge` (Exit-Code 14) abgelehnt.

## Aufhebung

`FileSystem`-Driven-Port um zwei Streaming-Primitive erweitert
(`Copy` und `CopyExclusive`); Backup-Pfade benutzen `io.Copy` statt
`ReadFile`+`WriteFile`. Memory-Footprint ist jetzt durch die
io.Copy-interne Buffergröße (typisch 32 KiB) begrenzt, unabhängig
von der Dateigröße. Der Cap und die `ErrBackupTooLarge`-Sentinel
sind entfernt.

## Geliefert

- **`port/driven.FileSystem`** zwei neue Methoden:
  - `Copy(src, dst string, mode fs.FileMode) error` — non-exclusive,
    truncate-overwrite, für nested-File-Copies in einem bereits
    reservierten Tree.
  - `CopyExclusive(src, dst string, mode fs.FileMode) error` —
    O_CREATE|O_EXCL, retourniert wrapped `fs.ErrExist` bei
    Collision; für den top-level Backup-Slot.
- **`adapter/driven/fs`** Implementierungen via `os.Open` + `os.OpenFile` +
  `io.Copy`. Gemeinsamer `streamCopy`-Helper, der den
  Open-Flag-Unterschied parametrisiert; `errors.Join` für
  Copy+Close-Fehler-Pairs.
- **`application/backup.go`** refactored:
  - `maxBackupFileSize` const entfernt.
  - Size-Check in `BackupPath` und `copyTreeNestedFile` entfernt.
  - `createBackupFile` ruft `fs.CopyExclusive` (statt
    `ReadFile`+`WriteFileExclusive`).
  - `copyTreeNestedFile` ruft `fs.Copy` (statt
    `ReadFile`+`WriteFile`).
- **`driving/initproject.go`**: `ErrBackupTooLarge` entfernt.
- **`cli/cli.go`**: `isFilesystemError` ohne `ErrBackupTooLarge`;
  Doc-Block-Code-Tabelle aktualisiert.
- **`cli/cli_test.go`**: ExitCode-Tabellen-Test ohne
  `ErrBackupTooLarge`.
- **`application/backup_test.go`**: zwei Tests gelöscht
  (`TestBackupPath_FileTooLarge_Rejected`,
  `TestBackupPath_NestedFileTooLarge_Rejected`).
- **`application/fakes_test.go`**: `sizeOverride`-Feld entfernt
  (samt init + Lstat-Override-Read). `Copy`/`CopyExclusive`-Fake-
  Methoden hinzugefügt, die `failReadOn`/`failOn` transparent ehren,
  damit die bestehenden Error-Injection-Tests
  (`TestBackupPath_TreeCopyFailure_RollsBack`,
  `TestBackupPath_RollbackFailure_JoinsErrors`,
  `TestBackupPath_ReadFileErrorPropagates`,
  `TestBackupPath_NestedReadFileErrorPropagates`) ohne Anpassung
  weiterlaufen.
- **`adapter/driven/fs/fs_test.go`**: 3 neue Adapter-Tests
  (`TestFS_Copy`, `TestFS_CopyExclusive_FailsOnExisting`,
  `TestFS_Copy_LargeFile_DoesNotOverallocate`). Der Large-File-Test
  erzeugt eine 1-GiB-Sparse-Datei via `os.Truncate(1 << 30)` und
  prüft, dass `Copy` sie streamen kann (Größe nach Stat: 1 GiB).
- **`.golangci.yml`** `interfacebloat`-Exclude für
  `internal/hexagon/port/driven/filesystem.go`: das Port hat jetzt
  12 Methoden (vorher 10) — bewusste Aufweichung für diese eine
  zentrale FS-Abstraktion, mit Begründung im Exclude-Kommentar.

## Akzeptanzkriterien (alle erfüllt)

- 1-GiB-Sparse-Datei wird erfolgreich gesichert
  (`TestFS_Copy_LargeFile_DoesNotOverallocate`).
- `make gates` grün; alle bestehenden Backup-Tests laufen weiter.
- Carveout-Eintrag in `carveouts.md` entfernt.
- Roadmap-Tabelle `Carveout-Auflösungs-Slices` auf `Done` gesetzt.

## Out of Scope

- Resumable-Backup nach Prozessabbruch (keine Partial-Wiederaufnahme).
- Backup-Integritätsprüfung (Hash/Checksum) — eigener Slice, falls
  jemals gefordert.
- Parallele Copies mehrerer Dateien — sequenziell bleibt die MVP-
  Charakteristik.

## Bezug

- Auslöser: T4a-Review Finding #6 (Commit `5296671` + Review-Fix
  `ecb8379` für M3-T4a).
- Spec-Bezug: `LH-FA-INIT-005` (Überschreibschutz), implizit
  `LH-NFA-USE-*` (kein Memory-Footgun).
- Aufhebung dokumentiert in: [`carveouts.md`](../in-progress/carveouts.md)
  und [`roadmap.md`](../in-progress/roadmap.md).
- `interfacebloat`-Exclude-Carveout in der **permanenten**
  Carveout-Tabelle dokumentiert (bewusste Schicht-Auflockerung).
