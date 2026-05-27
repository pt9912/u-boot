package application_test

import (
	"context"
	"errors"
	iofs "io/fs"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

const doctorBaseDir = "/tmp/u-boot-doctor/proj"

func newDoctorService(t *testing.T) (*application.DoctorService, *fakeFS, *fakeLogger) {
	t.Helper()
	fs := newFakeFS()
	fs.markDirExists(doctorBaseDir)
	logger := &fakeLogger{}
	return application.NewDoctorService(fs, logger), fs, logger
}

func TestDoctor_RequiresBaseDir(t *testing.T) {
	t.Parallel()
	svc := application.NewDoctorService(newFakeFS(), nil)
	_, err := svc.Check(context.Background(), driving.DoctorRequest{})
	if err == nil {
		t.Fatal("expected error for empty BaseDir, got nil")
	}
}

func TestDoctor_WritePermissions_OKOnWritableDir(t *testing.T) {
	t.Parallel()
	svc, _, _ := newDoctorService(t)

	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if len(resp.Report.Items) != 1 {
		t.Fatalf("Report.Items = %d, want 1", len(resp.Report.Items))
	}
	d := resp.Report.Items[0]
	if d.ID != "fs.write-permissions" {
		t.Errorf("ID = %q, want fs.write-permissions", d.ID)
	}
	if d.Severity != domain.SeverityOK {
		t.Errorf("Severity = %v, want OK", d.Severity)
	}
	if resp.Report.HasErrors() {
		t.Errorf("HasErrors() = true on writable dir")
	}
}

func TestDoctor_WritePermissions_ErrorOnDeniedWrite(t *testing.T) {
	t.Parallel()
	svc, fs, _ := newDoctorService(t)
	denied := errors.New("permission denied")
	// failOn on the exclusive-write of the sentinel makes
	// WriteFileExclusive fail with the configured error.
	fs.failOn = doctorBaseDir + "/.u-boot-doctor-probe"
	fs.failErr = denied

	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if !resp.Report.HasErrors() {
		t.Fatalf("HasErrors() = false, want true on denied write")
	}
	d := resp.Report.Items[0]
	if d.Severity != domain.SeverityError {
		t.Errorf("Severity = %v, want Error", d.Severity)
	}
	if !strings.Contains(d.Message, "permission denied") {
		t.Errorf("Message does not surface underlying error: %q", d.Message)
	}
	if !strings.Contains(d.Hint, "chmod") {
		t.Errorf("Hint missing chmod suggestion: %q", d.Hint)
	}
}

func TestDoctor_WritePermissions_ErrorOnExistingSentinel(t *testing.T) {
	t.Parallel()
	svc, fs, _ := newDoctorService(t)
	// Pre-existing sentinel → WriteFileExclusive returns ErrExist.
	if err := fs.WriteFile(doctorBaseDir+"/.u-boot-doctor-probe", []byte("stale"), 0o600); err != nil {
		t.Fatalf("setup: %v", err)
	}

	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if !resp.Report.HasErrors() {
		t.Fatal("HasErrors() = false, want true on stale sentinel")
	}
	d := resp.Report.Items[0]
	if d.Severity != domain.SeverityError {
		t.Errorf("Severity = %v, want Error", d.Severity)
	}
	// Hint must point at the cleanup path, not at chmod.
	if !strings.Contains(d.Hint, "Remove") || strings.Contains(d.Hint, "chmod") {
		t.Errorf("Hint does not match stale-sentinel case: %q", d.Hint)
	}
}

func TestDoctor_EmitsLoggerEvents(t *testing.T) {
	t.Parallel()
	svc, _, logger := newDoctorService(t)

	if _, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir}); err != nil {
		t.Fatalf("Check: %v", err)
	}

	// Expect at least one Debug "starting" + one Info "complete".
	var sawStart, sawComplete bool
	for _, e := range logger.entries {
		if e.Level == "DEBUG" && strings.Contains(e.Msg, "starting") {
			sawStart = true
		}
		if e.Level == "INFO" && strings.Contains(e.Msg, "complete") {
			sawComplete = true
		}
	}
	if !sawStart || !sawComplete {
		t.Errorf("logger entries = %v; want at least one start + one complete", logger.entries)
	}
}

func TestDoctor_SentinelCleanedUpOnSuccess(t *testing.T) {
	t.Parallel()
	svc, fs, _ := newDoctorService(t)
	if _, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir}); err != nil {
		t.Fatalf("Check: %v", err)
	}
	exists, err := fs.Exists(doctorBaseDir + "/.u-boot-doctor-probe")
	if err != nil && !errors.Is(err, iofs.ErrNotExist) {
		t.Fatalf("Exists check: %v", err)
	}
	if exists {
		t.Errorf("sentinel still exists after successful check — not cleaned up")
	}
}
