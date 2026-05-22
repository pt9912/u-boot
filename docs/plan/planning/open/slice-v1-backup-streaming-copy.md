# Slice V1: Streaming Backup-Copy (256 MiB-Cap aufheben)

## Motivation

Die T4a-Review (Finding #6) hat die Speicherspitze der Backup-
Mechanik aus `LH-FA-INIT-005` adressiert: `BackupPath` ruft heute
`fs.ReadFile`+`fs.WriteFile` auf, lädt also den vollständigen
Dateiinhalt in den Prozess-Heap. Bei einem multi-GB-Asset im
Backup-Scope (`docs/`, `docker/`, später `.devcontainer/`) würde das
den Prozess kippen.

Als MVP-Kompromiss wurde ein harter Cap eingezogen
(`maxBackupFileSize = 256 << 20` in
`internal/hexagon/application/backup.go`); überschreitende Dateien
werden mit `driving.ErrBackupTooLarge` abgelehnt (Exit-Code 14).
Dieser Slice hebt den Cap auf, sobald die `FileSystem`-Driven-Port
ein Streaming-Primitive bekommt.

## Scope

- Neue Methode auf `internal/hexagon/port/driven.FileSystem`:
  `Copy(src, dst string, mode fs.FileMode) error` (oder
  `CopyExclusive`-Variante für die TOCTOU-sichere Top-Level-
  Variante), Implementierung in `internal/adapter/driven/fs` via
  `os.Open`+`os.Create`+`io.Copy`.
- `application.BackupPath`/`copyTreeNestedFile`/`createBackupFile`
  rüsten von `ReadFile`+`WriteFile`/`WriteFileExclusive` auf das
  Streaming-Primitive um.
- `maxBackupFileSize`-Cap entfernen; `ErrBackupTooLarge` aus
  `driving/initproject.go` entfernen; `cli.isFilesystemError` /
  `ExitCode`-Tests anpassen.
- FakeFS in `application/fakes_test.go`: `Copy`-Methode hinzufügen,
  `sizeOverride`-Feld entfernen, die zugehörigen Tests
  (`TestBackupPath_FileTooLarge_Rejected`,
  `TestBackupPath_NestedFileTooLarge_Rejected`) löschen.

## Akzeptanzkriterien

- Eine 1 GiB-Test-Datei (real, via `os.Truncate` als sparse file)
  wird vom realen FS-Adapter erfolgreich gesichert, ohne dass die
  Prozess-RSS oberhalb der streaming-Block-Größe (typisch
  32–64 KiB) wächst (manuell beobachtet oder via `testing.AllocsPerRun`-
  Approximation).
- `make gates` grün; bestehende Backup-Tests laufen weiter.
- Carveout-Eintrag in `carveouts.md` entfernt; Eintrag in der
  Roadmap-Tabelle `Carveout-Auflösungs-Slices` auf `Done` gesetzt.

## Out of Scope

- Resumable-Backup nach Prozessabbruch (keine Partial-Wiederaufnahme).
- Backup-Integritätsprüfung (Hash/Checksum) — eigener Slice falls je
  gefordert.
- Parallele Copies mehrerer Dateien — sequenziell bleibt die MVP-
  Charakteristik.

## Bezug

- Auslöser: T4a-Review Finding #6 (Commit `5296671` + Review-Fix-
  Commit für M3-T4a).
- Spec-Bezug: `LH-FA-INIT-005` (Überschreibschutz), implizit
  `LH-NFA-USE-*` (kein Memory-Footgun).
- Hängt von: keinem anderen Slice; jederzeit ziehbar, sobald V1-
  Kapazität frei ist.
- Löst auf: Carveout `maxBackupFileSize 256 MiB-Cap in
  application/backup.go`.
