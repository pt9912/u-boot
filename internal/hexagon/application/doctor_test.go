package application_test

import (
	"context"
	"errors"
	iofs "io/fs"
	"path/filepath"
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
	// 11 checks total after M5-T7: 9 from M4 (write-permissions, git,
	// docker, docker-reachable, compose-installed, uboot.yaml,
	// compose.yaml, devcontainer.json, devcontainer.dockerfile) plus
	// the two new ones: services.enabled-key and
	// devcontainer.forwardPorts.consistency.
	if got := len(resp.Report.Items); got != 11 {
		t.Fatalf("Report.Items = %d, want 11", got)
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

func TestDoctor_UbootYaml_OKWithServicesBlock(t *testing.T) {
	t.Parallel()
	// Why: M5-T1 added a `services:` field to ubootYAMLConfig. A
	// u-boot.yaml that already carries the services-block (post-
	// `u-boot add`) must still parse and pass the doctor's validity
	// check. Roundtrip indicator for the schema extension — if
	// yaml.v3 can't decode `enabled: true` into *bool, this test
	// goes red before T2 lands.
	svc, fs, _, _, _ := newDoctorService(t)
	seedUbootYAML(t, fs, doctorBaseDir, `schemaVersion: 1
project:
  name: demo-service
services:
  postgres:
    enabled: true
  keycloak:
    enabled: false
`)

	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "uboot.yaml.valid")
	if d.Severity != domain.SeverityOK {
		t.Errorf("Severity = %v, want OK (services-block must not break validation)", d.Severity)
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

// ----------------------------------------------------------------------------
// T5 — compose.yaml validation
// ----------------------------------------------------------------------------

func seedComposeYAML(t *testing.T, fs *fakeFS, baseDir, body string) {
	t.Helper()
	if err := fs.WriteFile(baseDir+"/compose.yaml", []byte(body), 0o644); err != nil {
		t.Fatalf("seedComposeYAML: %v", err)
	}
}

func TestDoctor_ComposeYaml_WarnWhenMissing(t *testing.T) {
	t.Parallel()
	svc, _, _, _, _ := newDoctorService(t)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "compose.yaml.valid")
	if d.Severity != domain.SeverityWarn {
		t.Errorf("Severity = %v, want Warn (missing compose.yaml ≠ Error)", d.Severity)
	}
}

func TestDoctor_ComposeYaml_OKOnValidFile(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	seedComposeYAML(t, fs, doctorBaseDir,
		"services:\n  app:\n    image: nginx:latest\n  db:\n    image: postgres:16\n")

	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "compose.yaml.valid")
	if d.Severity != domain.SeverityOK {
		t.Errorf("Severity = %v, want OK", d.Severity)
	}
	if !strings.Contains(d.Message, "2 service(s)") {
		t.Errorf("Message does not surface service count: %q", d.Message)
	}
}

func TestDoctor_ComposeYaml_ErrorOnInvalidSyntax(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	// Unclosed bracket → fails yaml.v3 parse.
	seedComposeYAML(t, fs, doctorBaseDir, "services: [unclosed\n")

	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "compose.yaml.valid")
	if d.Severity != domain.SeverityError {
		t.Errorf("Severity = %v, want Error", d.Severity)
	}
	if !strings.Contains(d.Message, "not valid YAML") {
		t.Errorf("Message does not name the syntax problem: %q", d.Message)
	}
}

func TestDoctor_ComposeYaml_ErrorOnMissingServices(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	// Valid YAML but no services key — Compose without services is
	// not a meaningful Compose file.
	seedComposeYAML(t, fs, doctorBaseDir, "version: \"3.9\"\n")

	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "compose.yaml.valid")
	if d.Severity != domain.SeverityError {
		t.Errorf("Severity = %v, want Error", d.Severity)
	}
	if !strings.Contains(d.Message, "no `services:` entries") {
		t.Errorf("Message does not name the missing services: %q", d.Message)
	}
}

func TestDoctor_ComposeYaml_ErrorOnEmptyServices(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	// services-key present but empty mapping → still "no services".
	seedComposeYAML(t, fs, doctorBaseDir, "services:\n")

	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "compose.yaml.valid")
	if d.Severity != domain.SeverityError {
		t.Errorf("Severity = %v, want Error on empty services mapping", d.Severity)
	}
}

// ----------------------------------------------------------------------------
// T6 — devcontainer validation
// ----------------------------------------------------------------------------

func seedDevcontainerJSON(t *testing.T, fs *fakeFS, baseDir, body string) {
	t.Helper()
	if err := fs.WriteFile(baseDir+"/.devcontainer/devcontainer.json", []byte(body), 0o644); err != nil {
		t.Fatalf("seedDevcontainerJSON: %v", err)
	}
}

func seedDevcontainerDockerfile(t *testing.T, fs *fakeFS, baseDir, body string) {
	t.Helper()
	if err := fs.WriteFile(baseDir+"/.devcontainer/Dockerfile", []byte(body), 0o644); err != nil {
		t.Fatalf("seedDevcontainerDockerfile: %v", err)
	}
}

func TestDoctor_DevcontainerJSON_OKWhenAbsent(t *testing.T) {
	t.Parallel()
	svc, _, _, _, _ := newDoctorService(t)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.json.valid")
	if d.Severity != domain.SeverityOK {
		t.Errorf("Severity = %v, want OK (optional, not present)", d.Severity)
	}
}

func TestDoctor_DevcontainerJSON_OKOnValidFile(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	seedDevcontainerJSON(t, fs, doctorBaseDir,
		`{"name":"my-project","image":"ubuntu:22.04"}`)

	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.json.valid")
	if d.Severity != domain.SeverityOK {
		t.Errorf("Severity = %v, want OK", d.Severity)
	}
	if !strings.Contains(d.Message, "my-project") {
		t.Errorf("Message does not surface name: %q", d.Message)
	}
}

func TestDoctor_DevcontainerJSON_OKOnJSONCWithComments(t *testing.T) {
	t.Parallel()
	// Verifies the stripJSONC integration: line + block comments +
	// trailing commas in a realistic file shape.
	svc, fs, _, _, _ := newDoctorService(t)
	seedDevcontainerJSON(t, fs, doctorBaseDir, `{
  // The container
  "name": "demo",
  /* base image */
  "image": "ubuntu:22.04",
  "forwardPorts": [3000, 8080,],
}`)

	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.json.valid")
	if d.Severity != domain.SeverityOK {
		t.Errorf("Severity = %v, want OK on JSONC with comments", d.Severity)
	}
}

func TestDoctor_DevcontainerJSON_WarnOnInvalidJSON(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	seedDevcontainerJSON(t, fs, doctorBaseDir, `{"name": "demo", "image": [unclosed`)

	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.json.valid")
	if d.Severity != domain.SeverityWarn {
		t.Errorf("Severity = %v, want Warn (u-boot.yaml has no devcontainer.enabled schema yet)", d.Severity)
	}
	if !strings.Contains(d.Message, "not valid JSON") {
		t.Errorf("Message does not name JSON-syntax problem: %q", d.Message)
	}
}

func TestDoctor_DevcontainerJSON_WarnOnMissingName(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	seedDevcontainerJSON(t, fs, doctorBaseDir, `{"image":"ubuntu:22.04"}`)

	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.json.valid")
	if d.Severity != domain.SeverityWarn {
		t.Errorf("Severity = %v, want Warn", d.Severity)
	}
	if !strings.Contains(d.Message, "missing required `name`") {
		t.Errorf("Message does not name the missing field: %q", d.Message)
	}
}

func TestDoctor_DevcontainerJSON_WarnOnMissingImageAndBuild(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	// Name set but neither image nor build → fails VS Code Dev
	// Container minimum compatibility.
	seedDevcontainerJSON(t, fs, doctorBaseDir, `{"name":"demo"}`)

	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.json.valid")
	if d.Severity != domain.SeverityWarn {
		t.Errorf("Severity = %v, want Warn", d.Severity)
	}
	if !strings.Contains(d.Message, "`image` or `build`") {
		t.Errorf("Message does not name image/build requirement: %q", d.Message)
	}
}

func TestDoctor_DevcontainerJSON_OKWithBuildAsObject(t *testing.T) {
	t.Parallel()
	// `build` can be an object (dockerfile-context tuple) instead of
	// a string — both satisfy the minimum-compat check.
	svc, fs, _, _, _ := newDoctorService(t)
	seedDevcontainerJSON(t, fs, doctorBaseDir,
		`{"name":"demo","build":{"dockerfile":"Dockerfile","context":"."}}`)

	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.json.valid")
	if d.Severity != domain.SeverityOK {
		t.Errorf("Severity = %v, want OK with build-as-object", d.Severity)
	}
}

func TestDoctor_DevcontainerDockerfile_OKWhenAbsent(t *testing.T) {
	t.Parallel()
	svc, _, _, _, _ := newDoctorService(t)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.dockerfile.valid")
	if d.Severity != domain.SeverityOK {
		t.Errorf("Severity = %v, want OK (optional, not present)", d.Severity)
	}
}

func TestDoctor_DevcontainerDockerfile_OKWithFromDirective(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	seedDevcontainerDockerfile(t, fs, doctorBaseDir,
		"# Comment\n\nFROM ubuntu:22.04\nRUN apt-get update\n")

	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.dockerfile.valid")
	if d.Severity != domain.SeverityOK {
		t.Errorf("Severity = %v, want OK", d.Severity)
	}
}

func TestDoctor_DevcontainerDockerfile_WarnWithoutFrom(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	// Dockerfile with only comments — no FROM directive.
	seedDevcontainerDockerfile(t, fs, doctorBaseDir, "# just a comment\n\n")

	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.dockerfile.valid")
	if d.Severity != domain.SeverityWarn {
		t.Errorf("Severity = %v, want Warn", d.Severity)
	}
	if !strings.Contains(d.Message, "no `FROM` directive") {
		t.Errorf("Message does not name the missing directive: %q", d.Message)
	}
}

func TestDoctor_DevcontainerDockerfile_LowercaseFromAccepted(t *testing.T) {
	t.Parallel()
	// Docker's parser is case-insensitive; the doctor should be too.
	svc, fs, _, _, _ := newDoctorService(t)
	seedDevcontainerDockerfile(t, fs, doctorBaseDir, "from ubuntu:22.04\n")

	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.dockerfile.valid")
	if d.Severity != domain.SeverityOK {
		t.Errorf("Severity = %v, want OK with lowercase `from`", d.Severity)
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

// --- M5-T7: devcontainer.enabled severity escalation ---------------

func TestDoctor_T7_DevcontainerSeverity_EscalatesWhenEnabledTrue(t *testing.T) {
	svc, fs, _, _, _ := newDoctorService(t)
	// u-boot.yaml flips devcontainer.enabled=true → all devcontainer
	// checks must now be Error per LH-FA-DIAG-002 §1073.
	if err := fs.WriteFile(filepath.Join(doctorBaseDir, "u-boot.yaml"),
		[]byte("schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n"), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	// Seed a devcontainer.json that fails the LH-FA-DIAG-002 minimum-
	// compat check (missing both `image` and `build`) so the
	// devcontainer-validation path returns the escalated severity.
	if err := fs.WriteFile(filepath.Join(doctorBaseDir, ".devcontainer", "devcontainer.json"),
		[]byte(`{"name":"demo"}`), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.json.valid")
	if d.Severity != domain.SeverityError {
		t.Errorf("devcontainer.json.valid Severity = %v, want Error when devcontainer.enabled=true",
			d.Severity)
	}
}

func TestDoctor_T7_DevcontainerSeverity_WarnWhenEnabledFalse(t *testing.T) {
	svc, fs, _, _, _ := newDoctorService(t)
	if err := fs.WriteFile(filepath.Join(doctorBaseDir, "u-boot.yaml"),
		[]byte("schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: false\n"), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if err := fs.WriteFile(filepath.Join(doctorBaseDir, ".devcontainer", "devcontainer.json"),
		[]byte(`{"name":"demo"}`), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.json.valid")
	if d.Severity != domain.SeverityWarn {
		t.Errorf("devcontainer.json.valid Severity = %v, want Warn when devcontainer.enabled=false",
			d.Severity)
	}
}

// --- M5-T7: services.enabled-key ------------------------------------

func TestDoctor_T7_ServicesEnabledKey_OKWhenAllExplicit(t *testing.T) {
	svc, fs, _, _, _ := newDoctorService(t)
	if err := fs.WriteFile(filepath.Join(doctorBaseDir, "u-boot.yaml"),
		[]byte("schemaVersion: 1\nproject:\n  name: demo\nservices:\n  postgres:\n    enabled: true\n  keycloak:\n    enabled: false\n"),
		0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "services.enabled-key")
	if d.Severity != domain.SeverityOK {
		t.Errorf("services.enabled-key Severity = %v, want OK", d.Severity)
	}
}

func TestDoctor_T7_ServicesEnabledKey_WarnWhenMissing(t *testing.T) {
	svc, fs, _, _, _ := newDoctorService(t)
	// postgres has no enabled-key; keycloak does.
	if err := fs.WriteFile(filepath.Join(doctorBaseDir, "u-boot.yaml"),
		[]byte("schemaVersion: 1\nproject:\n  name: demo\nservices:\n  postgres: {}\n  keycloak:\n    enabled: true\n"),
		0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "services.enabled-key")
	if d.Severity != domain.SeverityWarn {
		t.Errorf("services.enabled-key Severity = %v, want Warn", d.Severity)
	}
	if !strings.Contains(d.Message, "postgres") {
		t.Errorf("warn message should list `postgres`; got %q", d.Message)
	}
}

// --- M5-T7: devcontainer.forwardPorts.consistency -------------------

func TestDoctor_T7_ForwardPorts_OKWhenComplete(t *testing.T) {
	svc, fs, _, _, _ := newDoctorService(t)
	if err := fs.WriteFile(filepath.Join(doctorBaseDir, "u-boot.yaml"),
		[]byte("schemaVersion: 1\nproject:\n  name: demo\nservices:\n  postgres:\n    enabled: true\n"),
		0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if err := fs.WriteFile(filepath.Join(doctorBaseDir, "compose.yaml"),
		[]byte("services:\n  postgres:\n    image: postgres\n    ports:\n      - \"5432:5432\"\n"),
		0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if err := fs.WriteFile(filepath.Join(doctorBaseDir, ".devcontainer", "devcontainer.json"),
		[]byte("{\"name\":\"demo\",\"image\":\"x\",\"forwardPorts\":[5432]}"),
		0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.forwardPorts.consistency")
	if d.Severity != domain.SeverityOK {
		t.Errorf("forwardPorts Severity = %v, want OK; msg=%q", d.Severity, d.Message)
	}
}

func TestDoctor_T7_ForwardPorts_WarnWhenMissing(t *testing.T) {
	svc, fs, _, _, _ := newDoctorService(t)
	if err := fs.WriteFile(filepath.Join(doctorBaseDir, "u-boot.yaml"),
		[]byte("schemaVersion: 1\nproject:\n  name: demo\nservices:\n  postgres:\n    enabled: true\n"),
		0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if err := fs.WriteFile(filepath.Join(doctorBaseDir, "compose.yaml"),
		[]byte("services:\n  postgres:\n    image: postgres\n    ports:\n      - \"5432:5432\"\n"),
		0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if err := fs.WriteFile(filepath.Join(doctorBaseDir, ".devcontainer", "devcontainer.json"),
		[]byte("{\"name\":\"demo\",\"image\":\"x\",\"forwardPorts\":[3000]}"),
		0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.forwardPorts.consistency")
	if d.Severity != domain.SeverityWarn {
		t.Errorf("forwardPorts Severity = %v, want Warn; msg=%q", d.Severity, d.Message)
	}
	if !strings.Contains(d.Message, "5432") {
		t.Errorf("warn message should mention 5432; got %q", d.Message)
	}
}

func TestDoctor_T7_ForwardPorts_SkippedWhenDevcontainerMissing(t *testing.T) {
	svc, fs, _, _, _ := newDoctorService(t)
	if err := fs.WriteFile(filepath.Join(doctorBaseDir, "u-boot.yaml"),
		[]byte("schemaVersion: 1\nproject:\n  name: demo\nservices:\n  postgres:\n    enabled: true\n"),
		0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if err := fs.WriteFile(filepath.Join(doctorBaseDir, "compose.yaml"),
		[]byte("services:\n  postgres:\n    image: postgres\n    ports:\n      - \"5432:5432\"\n"),
		0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.forwardPorts.consistency")
	if d.Severity != domain.SeverityOK {
		t.Errorf("forwardPorts Severity = %v, want OK (skipped); msg=%q", d.Severity, d.Message)
	}
}

func TestDoctor_T7_ForwardPorts_HostContainerMappingExtractsContainerPort(t *testing.T) {
	svc, fs, _, _, _ := newDoctorService(t)
	if err := fs.WriteFile(filepath.Join(doctorBaseDir, "u-boot.yaml"),
		[]byte("schemaVersion: 1\nproject:\n  name: demo\nservices:\n  api:\n    enabled: true\n"),
		0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	// Compose host:container, devcontainer forwards the container port.
	if err := fs.WriteFile(filepath.Join(doctorBaseDir, "compose.yaml"),
		[]byte("services:\n  api:\n    image: foo\n    ports:\n      - \"8080:80\"\n"),
		0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if err := fs.WriteFile(filepath.Join(doctorBaseDir, ".devcontainer", "devcontainer.json"),
		[]byte("{\"name\":\"demo\",\"image\":\"x\",\"forwardPorts\":[80]}"),
		0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	d := findDiagnostic(t, resp.Report.Items, "devcontainer.forwardPorts.consistency")
	if d.Severity != domain.SeverityOK {
		t.Errorf("expected container-port extraction OK, got %v; msg=%q", d.Severity, d.Message)
	}
}
