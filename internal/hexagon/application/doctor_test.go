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

// goodDockerProbe returns a DockerProbe configured for the happy
// path: docker 24.0.7 client, daemon reachable, compose 2.20.0
// plugin. Tests that want to isolate one failing check override the
// returned probe's fields.
func goodDockerProbe() *fakeDockerProbe {
	return &fakeDockerProbe{
		version:        "24.0.7",
		composeVersion: "2.20.0",
	}
}

// goodGit returns a Git fake configured for the happy path.
func goodGit() *fakeGit {
	return &fakeGit{version: "2.43.0"}
}

func newDoctorService(t *testing.T) (*application.DoctorService, *fakeFS, *fakeGit, *fakeDockerProbe, *fakeLogger) {
	t.Helper()
	fs := newFakeFS()
	fs.markDirExists(doctorBaseDir)
	git := goodGit()
	docker := goodDockerProbe()
	logger := &fakeLogger{}
	return application.NewDoctorService(fs, &fakeYAML{}, git, docker, logger), fs, git, docker, logger
}

func TestDoctor_RequiresBaseDir(t *testing.T) {
	t.Parallel()
	svc := application.NewDoctorService(newFakeFS(), &fakeYAML{}, goodGit(), goodDockerProbe(), nil)
	_, err := svc.Check(context.Background(), driving.DoctorRequest{})
	if err == nil {
		t.Fatal("expected error for empty BaseDir, got nil")
	}
	if !errors.Is(err, driving.ErrBaseDirMissing) {
		t.Errorf("err = %v, want wrapped ErrBaseDirMissing (shared sentinel with init)", err)
	}
}

func TestDoctor_WritePermissions_OKOnWritableDir(t *testing.T) {
	t.Parallel()
	svc, _, _, _, _ := newDoctorService(t)

	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	// 6 checks total after T4: write-permissions, git, docker,
	// docker-reachable, compose, uboot.yaml.
	if got := len(resp.Report.Items); got != 6 {
		t.Fatalf("Report.Items = %d, want 6", got)
	}
	d := findDiagnostic(t, resp.Report.Items, "fs.write-permissions")
	if d.Severity != domain.SeverityOK {
		t.Errorf("write-permissions Severity = %v, want OK", d.Severity)
	}
	// Happy path's u-boot.yaml is the absent-warn case (newDoctorService
	// does not seed one); the report is allowed to have warnings on
	// the happy path, but no errors.
	if resp.Report.HasErrors() {
		t.Errorf("HasErrors() = true on happy path")
	}
}

// findDiagnostic returns the single Diagnostic with the given ID
// from items. Calls t.Fatal when there is no match or more than one
// — both indicate a broken test, not the system under test.
func findDiagnostic(t *testing.T, items []domain.Diagnostic, id string) domain.Diagnostic {
	t.Helper()
	var found *domain.Diagnostic
	for i := range items {
		if items[i].ID == id {
			if found != nil {
				t.Fatalf("findDiagnostic(%q): duplicate match in items", id)
			}
			found = &items[i]
		}
	}
	if found == nil {
		t.Fatalf("findDiagnostic(%q): not found in items %+v", id, items)
	}
	return *found
}

func TestDoctor_WritePermissions_ErrorOnDeniedWrite(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
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
	d := findDiagnostic(t, resp.Report.Items, "fs.write-permissions")
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
	svc, fs, _, _, _ := newDoctorService(t)
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
	d := findDiagnostic(t, resp.Report.Items, "fs.write-permissions")
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
	svc, _, _, _, logger := newDoctorService(t)

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

// ----------------------------------------------------------------------------
// T3 — external-binary probes (git, docker, compose)
// ----------------------------------------------------------------------------

func TestDoctor_Git_OKWhenAvailable(t *testing.T) {
	t.Parallel()
	svc, _, _, _, _ := newDoctorService(t)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "git.installed")
	if d.Severity != domain.SeverityOK {
		t.Errorf("Severity = %v, want OK", d.Severity)
	}
	if !strings.Contains(d.Message, "2.43.0") {
		t.Errorf("Message does not surface version: %q", d.Message)
	}
}

func TestDoctor_Git_ErrorWhenMissing(t *testing.T) {
	t.Parallel()
	fs := newFakeFS()
	fs.markDirExists(doctorBaseDir)
	git := &fakeGit{versionErr: errors.New("exec: \"git\": executable file not found")}
	svc := application.NewDoctorService(fs, &fakeYAML{}, git, goodDockerProbe(), nil)

	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "git.installed")
	if d.Severity != domain.SeverityError {
		t.Errorf("Severity = %v, want Error", d.Severity)
	}
	if !strings.Contains(d.Hint, "Install git") {
		t.Errorf("Hint missing install suggestion: %q", d.Hint)
	}
}

func TestDoctor_Docker_OKAtOrAboveMinimum(t *testing.T) {
	t.Parallel()
	cases := []string{"24.0.0", "24.0.7", "24.5.1", "25.0.0", "30.99.0"}
	for _, v := range cases {
		v := v
		t.Run(v, func(t *testing.T) {
			t.Parallel()
			fs := newFakeFS()
			fs.markDirExists(doctorBaseDir)
			docker := goodDockerProbe()
			docker.version = v
			svc := application.NewDoctorService(fs, &fakeYAML{}, goodGit(), docker, nil)
			resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
			if err != nil {
				t.Fatalf("Check: %v", err)
			}
			d := findDiagnostic(t, resp.Report.Items, "docker.installed")
			if d.Severity != domain.SeverityOK {
				t.Errorf("Severity = %v, want OK (version %s)", d.Severity, v)
			}
		})
	}
}

func TestDoctor_Docker_ErrorBelowMinimum(t *testing.T) {
	t.Parallel()
	cases := []string{"23.0.0", "20.10.21", "1.0.0"}
	for _, v := range cases {
		v := v
		t.Run(v, func(t *testing.T) {
			t.Parallel()
			fs := newFakeFS()
			fs.markDirExists(doctorBaseDir)
			docker := goodDockerProbe()
			docker.version = v
			svc := application.NewDoctorService(fs, &fakeYAML{}, goodGit(), docker, nil)
			resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
			if err != nil {
				t.Fatalf("Check: %v", err)
			}
			d := findDiagnostic(t, resp.Report.Items, "docker.installed")
			if d.Severity != domain.SeverityError {
				t.Errorf("Severity = %v, want Error (version %s)", d.Severity, v)
			}
			if !strings.Contains(d.Message, "below the LH-FA-DIAG-002 minimum") {
				t.Errorf("Message missing minimum hint: %q", d.Message)
			}
		})
	}
}

func TestDoctor_Docker_WarnOnUnparseableVersion(t *testing.T) {
	t.Parallel()
	fs := newFakeFS()
	fs.markDirExists(doctorBaseDir)
	docker := goodDockerProbe()
	docker.version = "garbage"
	svc := application.NewDoctorService(fs, &fakeYAML{}, goodGit(), docker, nil)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "docker.installed")
	if d.Severity != domain.SeverityWarn {
		t.Errorf("Severity = %v, want Warn on unparseable version", d.Severity)
	}
}

func TestDoctor_Docker_ErrorWhenMissing(t *testing.T) {
	t.Parallel()
	fs := newFakeFS()
	fs.markDirExists(doctorBaseDir)
	docker := goodDockerProbe()
	docker.versionErr = errors.New("exec: \"docker\": executable file not found")
	svc := application.NewDoctorService(fs, &fakeYAML{}, goodGit(), docker, nil)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "docker.installed")
	if d.Severity != domain.SeverityError {
		t.Errorf("Severity = %v, want Error", d.Severity)
	}
	if !strings.Contains(d.Hint, "Install Docker") {
		t.Errorf("Hint missing install pointer: %q", d.Hint)
	}
}

func TestDoctor_DockerReachable_ErrorOnDaemonDown(t *testing.T) {
	t.Parallel()
	fs := newFakeFS()
	fs.markDirExists(doctorBaseDir)
	docker := goodDockerProbe()
	docker.infoErr = errors.New("Cannot connect to the Docker daemon at unix:///var/run/docker.sock")
	svc := application.NewDoctorService(fs, &fakeYAML{}, goodGit(), docker, nil)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "docker.reachable")
	if d.Severity != domain.SeverityError {
		t.Errorf("Severity = %v, want Error", d.Severity)
	}
	if !strings.Contains(d.Hint, "Start Docker") {
		t.Errorf("Hint missing start-docker pointer: %q", d.Hint)
	}
}

func TestDoctor_Compose_OKAtOrAboveMinimum(t *testing.T) {
	t.Parallel()
	cases := []string{"2.20.0", "2.21.5", "2.99.0", "3.0.0"}
	for _, v := range cases {
		v := v
		t.Run(v, func(t *testing.T) {
			t.Parallel()
			fs := newFakeFS()
			fs.markDirExists(doctorBaseDir)
			docker := goodDockerProbe()
			docker.composeVersion = v
			svc := application.NewDoctorService(fs, &fakeYAML{}, goodGit(), docker, nil)
			resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
			if err != nil {
				t.Fatalf("Check: %v", err)
			}
			d := findDiagnostic(t, resp.Report.Items, "docker.compose.installed")
			if d.Severity != domain.SeverityOK {
				t.Errorf("Severity = %v, want OK (version %s)", d.Severity, v)
			}
		})
	}
}

func TestDoctor_Compose_ErrorBelowMinimum(t *testing.T) {
	t.Parallel()
	cases := []string{"2.19.0", "2.0.0", "1.29.2"}
	for _, v := range cases {
		v := v
		t.Run(v, func(t *testing.T) {
			t.Parallel()
			fs := newFakeFS()
			fs.markDirExists(doctorBaseDir)
			docker := goodDockerProbe()
			docker.composeVersion = v
			svc := application.NewDoctorService(fs, &fakeYAML{}, goodGit(), docker, nil)
			resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
			if err != nil {
				t.Fatalf("Check: %v", err)
			}
			d := findDiagnostic(t, resp.Report.Items, "docker.compose.installed")
			if d.Severity != domain.SeverityError {
				t.Errorf("Severity = %v, want Error (version %s)", d.Severity, v)
			}
		})
	}
}

func TestDoctor_Compose_ErrorWhenMissing(t *testing.T) {
	t.Parallel()
	fs := newFakeFS()
	fs.markDirExists(doctorBaseDir)
	docker := goodDockerProbe()
	docker.composeVersionErr = errors.New("docker: 'compose' is not a docker command")
	svc := application.NewDoctorService(fs, &fakeYAML{}, goodGit(), docker, nil)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "docker.compose.installed")
	if d.Severity != domain.SeverityError {
		t.Errorf("Severity = %v, want Error", d.Severity)
	}
}

// ----------------------------------------------------------------------------
// T4 — u-boot.yaml validation
// ----------------------------------------------------------------------------

// seedUbootYAML writes a u-boot.yaml at baseDir with the given body.
// Tests use this to dial in the validation paths (valid file,
// malformed YAML, wrong schemaVersion, invalid project name, etc.).
func seedUbootYAML(t *testing.T, fs *fakeFS, baseDir, body string) {
	t.Helper()
	if err := fs.WriteFile(baseDir+"/u-boot.yaml", []byte(body), 0o644); err != nil {
		t.Fatalf("seedUbootYAML: %v", err)
	}
}

func TestDoctor_UbootYaml_WarnWhenMissing(t *testing.T) {
	t.Parallel()
	svc, _, _, _, _ := newDoctorService(t)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "uboot.yaml.valid")
	if d.Severity != domain.SeverityWarn {
		t.Errorf("Severity = %v, want Warn (missing u-boot.yaml ≠ Error)", d.Severity)
	}
	if !strings.Contains(d.Hint, "u-boot init") {
		t.Errorf("Hint missing init suggestion: %q", d.Hint)
	}
}

func TestDoctor_UbootYaml_OKOnValidFile(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	seedUbootYAML(t, fs, doctorBaseDir, "schemaVersion: 1\nproject:\n  name: demo-service\n")

	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "uboot.yaml.valid")
	if d.Severity != domain.SeverityOK {
		t.Errorf("Severity = %v, want OK", d.Severity)
	}
	if !strings.Contains(d.Message, "demo-service") {
		t.Errorf("Message does not surface project name: %q", d.Message)
	}
}

func TestDoctor_UbootYaml_ErrorOnInvalidSyntax(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	// Unclosed bracket → fails yaml.v3 parse.
	seedUbootYAML(t, fs, doctorBaseDir, "schemaVersion: 1\nproject: [unclosed\n")

	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "uboot.yaml.valid")
	if d.Severity != domain.SeverityError {
		t.Errorf("Severity = %v, want Error", d.Severity)
	}
	if !strings.Contains(d.Message, "not valid YAML") {
		t.Errorf("Message does not name the syntax problem: %q", d.Message)
	}
}

func TestDoctor_UbootYaml_ErrorOnWrongSchemaVersion(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	seedUbootYAML(t, fs, doctorBaseDir, "schemaVersion: 2\nproject:\n  name: demo\n")

	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "uboot.yaml.valid")
	if d.Severity != domain.SeverityError {
		t.Errorf("Severity = %v, want Error", d.Severity)
	}
	if !strings.Contains(d.Message, "schemaVersion is 2") {
		t.Errorf("Message does not name the wrong version: %q", d.Message)
	}
}

func TestDoctor_UbootYaml_ErrorOnMissingProjectName(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	seedUbootYAML(t, fs, doctorBaseDir, "schemaVersion: 1\nproject:\n  name: \"\"\n")

	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "uboot.yaml.valid")
	if d.Severity != domain.SeverityError {
		t.Errorf("Severity = %v, want Error", d.Severity)
	}
	if !strings.Contains(d.Message, "missing required") {
		t.Errorf("Message does not name the missing field: %q", d.Message)
	}
}

func TestDoctor_UbootYaml_ErrorOnInvalidProjectName(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	// Uppercase violates LH-FA-INIT-006 (lowercase-only regex).
	seedUbootYAML(t, fs, doctorBaseDir, "schemaVersion: 1\nproject:\n  name: DemoService\n")

	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "uboot.yaml.valid")
	if d.Severity != domain.SeverityError {
		t.Errorf("Severity = %v, want Error", d.Severity)
	}
	if !strings.Contains(d.Message, "DemoService") {
		t.Errorf("Message does not surface offending name: %q", d.Message)
	}
}

func TestDoctor_SemverPreReleaseHandledAsMajorMinor(t *testing.T) {
	t.Parallel()
	// `2.20.0-rc1` must parse as 2.20 (≥ 2.20 minimum → OK).
	fs := newFakeFS()
	fs.markDirExists(doctorBaseDir)
	docker := goodDockerProbe()
	docker.composeVersion = "2.20.0-rc1"
	svc := application.NewDoctorService(fs, &fakeYAML{}, goodGit(), docker, nil)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "docker.compose.installed")
	if d.Severity != domain.SeverityOK {
		t.Errorf("Severity = %v, want OK on 2.20.0-rc1", d.Severity)
	}
}

func TestDoctor_SentinelCleanedUpOnSuccess(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
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
