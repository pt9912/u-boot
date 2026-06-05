package application_test

import (
	"errors"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// TestInitProjectService_RunBackup_RawFSErrorWrapsErrInitFileSystem
// pinnt das T0-(f) Switch-Order-Pendant der Wrap-Strategie aus
// initproject.go:1118-1132 (Review-Round-9 #3): wenn BackupPath
// einen NICHT-typisierten FS-Fehler hochwirft (z.B. EIO im Lstat,
// permission-denied beim Copy), MUSS runBackup ihn mit
// ErrInitFileSystem wrappen — nicht mit einer Backup-*-Sentinel.
// Sonst klassifiziert der CLI-Mapper den Fehler fälschlich als
// LH-FA-INIT-005 (Backup-Kategorie) statt LH-NFA-REL-003 (FS-Klasse,
// Exit 14).
func TestInitProjectService_RunBackup_RawFSErrorWrapsErrInitFileSystem(t *testing.T) {
	t.Parallel()
	fakeFS := newFakeFS()
	fakeFS.dirs["/proj"] = true
	fakeFS.files["/proj/u-boot.yaml"] = []byte("schemaVersion: 1\n")
	rawErr := errors.New("EIO: simulated disk failure")
	fakeFS.failLstatOn = "/proj/u-boot.yaml"
	fakeFS.failLstatErr = rawErr

	svc := application.NewInitProjectService(fakeFS, &fakeYAML{}, &fakeGit{}, nil, nil, nil)

	_, err := svc.RunBackupForTest("/proj", "u-boot.yaml")
	if err == nil {
		t.Fatal("runBackup: want raw-FS-error wrapped with ErrInitFileSystem, got nil")
	}
	if !errors.Is(err, driving.ErrInitFileSystem) {
		t.Errorf("err must wrap ErrInitFileSystem (FS-class), got: %v", err)
	}
	// Anti-Mis-Classification: darf NICHT als eine der typed Backup-
	// Sentinels durchgehen (sonst LH-FA-INIT-005-Mapping statt
	// LH-NFA-REL-003).
	for _, typed := range []error{
		driving.ErrBackupSuffixExhausted,
		driving.ErrBackupSourceMissing,
		driving.ErrBackupUnsupportedKind,
	} {
		if errors.Is(err, typed) {
			t.Errorf("raw FS-error must NOT classify as %v; got: %v", typed, err)
		}
	}
	// Der rohe FS-Fehler muss durch %w erreichbar bleiben (Diagnose-
	// Information für den User).
	if !errors.Is(err, rawErr) {
		t.Errorf("raw FS-error must remain unwrappable, got: %v", err)
	}
}

// TestInitProjectService_RunBackup_TypedSentinelsPassThrough pinnt
// die zweite Seite des Wrap-Strategie-Splits: die typed Backup-
// Sentinels werden DURCHGEREICHT (nur mit `backup %s: %w`-Präfix),
// NICHT mit ErrInitFileSystem überschrieben. Sonst wäre die CLI-
// Mapper-Klassifikation eindeutig (LH-NFA-REL-003), aber die
// Diagnose-Information welche Backup-Phase scheiterte ginge
// verloren.
func TestInitProjectService_RunBackup_TypedSentinelsPassThrough(t *testing.T) {
	t.Parallel()
	fakeFS := newFakeFS()
	fakeFS.dirs["/proj"] = true
	// Lstat returns ErrNotExist (kein File da) → BackupPath returnt
	// ErrBackupSourceMissing.
	// Keine fileSetup nötig — Lstat-Default ist ErrNotExist.

	svc := application.NewInitProjectService(fakeFS, &fakeYAML{}, &fakeGit{}, nil, nil, nil)

	_, err := svc.RunBackupForTest("/proj", "u-boot.yaml")
	if err == nil {
		t.Fatal("runBackup: want ErrBackupSourceMissing, got nil")
	}
	if !errors.Is(err, driving.ErrBackupSourceMissing) {
		t.Errorf("typed sentinel must pass through, got: %v", err)
	}
	// Typed-Pfad darf NICHT zusätzlich mit ErrInitFileSystem gewrappt
	// sein — sonst klassifiziert der CLI-Mapper LH-NFA-REL-003 statt
	// der präziseren Backup-Diagnose. (Hier wäre beides FS-Klasse,
	// aber das Prinzip muss eindeutig sein.)
	if errors.Is(err, driving.ErrInitFileSystem) {
		t.Errorf("typed sentinel must NOT be additionally wrapped with ErrInitFileSystem; got: %v", err)
	}
}

