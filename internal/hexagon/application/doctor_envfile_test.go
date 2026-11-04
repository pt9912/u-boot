package application_test

import (
	"context"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// Cases of the `devcontainer.envfile` check (LH-FA-DEV-012): existence
// of the environment file that `runArgs` names, nothing else.

const envRefJSON = `{"runArgs": ["` + envFileOption + `"]}`

func envFileDiagnostic(t *testing.T, setup func(fs *fakeFS)) domain.Diagnostic {
	t.Helper()
	return findDiagnostic(t, runSandboxDoctor(t, setup), "devcontainer.envfile")
}

// Happy: the line is there and the file is missing: warn with the path and
// the repair command; no flag that would survive `--quiet`.
func TestDoctor_EnvFile_WarnsWhenMissing(t *testing.T) {
	t.Parallel()
	d := envFileDiagnostic(t, func(fs *fakeFS) {
		writeDoctorFile(t, fs, ".devcontainer/devcontainer.json", envRefJSON)
	})
	if d.Severity != domain.SeverityWarn {
		t.Fatalf("severity = %v, want warn: %+v", d.Severity, d)
	}
	if !strings.Contains(d.Message, ".devcontainer/devcontainer.env") || !strings.Contains(d.Hint, "u-boot generate devcontainer") {
		t.Errorf("message/hint = %q / %q", d.Message, d.Hint)
	}
	if d.SecurityRelaxation || d.DegradationFallback {
		t.Errorf("the finding must carry neither flag: %+v", d)
	}
}

// Happy: the file exists, whatever it holds.
func TestDoctor_EnvFile_OKWhenPresent(t *testing.T) {
	t.Parallel()
	d := envFileDiagnostic(t, func(fs *fakeFS) {
		writeDoctorFile(t, fs, ".devcontainer/devcontainer.json", envRefJSON)
		writeDoctorFile(t, fs, ".devcontainer/devcontainer.env", "")
	})
	if d.Severity != domain.SeverityOK {
		t.Errorf("severity = %v, want ok: %+v", d.Severity, d)
	}
}

// Boundary: without the line or without devcontainer.json the check is
// skipped (ok), also when the file is missing.
func TestDoctor_EnvFile_SkippedWithoutReference(t *testing.T) {
	t.Parallel()
	cases := map[string]func(fs *fakeFS){
		"no devcontainer.json": func(*fakeFS) {},
		"json without line": func(fs *fakeFS) {
			writeDoctorFile(t, fs, ".devcontainer/devcontainer.json", `{"runArgs": ["--privileged"]}`)
		},
		"json without runArgs": func(fs *fakeFS) {
			writeDoctorFile(t, fs, ".devcontainer/devcontainer.json", `{"name": "x"}`)
		},
	}
	for name, setup := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			d := envFileDiagnostic(t, setup)
			if d.Severity != domain.SeverityOK || !strings.Contains(d.Message, "skipped") {
				t.Errorf("diagnostic = %+v, want ok and skipped", d)
			}
		})
	}
}

// Negative: the check never reads the environment file.
func TestDoctor_EnvFile_NeverReadsContent(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	writeDoctorFile(t, fs, ".devcontainer/devcontainer.json", envRefJSON)
	writeDoctorFile(t, fs, ".devcontainer/devcontainer.env", "SECRET=value\n")
	if _, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir}); err != nil {
		t.Fatalf("Check: %v", err)
	}
	if n := fs.readFileCallCount(doctorBaseDir + "/.devcontainer/devcontainer.env"); n != 0 {
		t.Errorf("the environment file was read %d times", n)
	}
}

// Negative: the plaintext credential check of the sandbox profile does not
// scan the environment file, even when it holds a token pattern.
func TestDoctor_SandboxCredentials_IgnoresEnvFile(t *testing.T) {
	t.Parallel()
	items := runSandboxDoctor(t, func(fs *fakeFS) {
		writeDoctorFile(t, fs, "u-boot.yaml", sandboxYAMLHead+"  profile: sandbox\n")
		writeDoctorFile(t, fs, ".devcontainer/devcontainer.json",
			`{"remoteEnv": {"GIT_TOKEN": "${localEnv:GIT_TOKEN}"}, "runArgs": ["`+envFileOption+`"]}`)
		writeDoctorFile(t, fs, ".devcontainer/devcontainer.env", "GIT_TOKEN=abcdef123456789xyz\n")
	})
	d := findDiagnostic(t, items, "devcontainer.sandbox.credentials")
	if d.Severity != domain.SeverityOK {
		t.Errorf("credentials = %+v, want ok (the environment file is not scanned)", d)
	}
}
