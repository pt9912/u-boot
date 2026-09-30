package application_test

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

const (
	sandboxYAMLHead   = "schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n"
	sandboxYAMLPodman = sandboxYAMLHead + "  profile: sandbox\n  sandbox:\n    nestedRuntime: podman\n"
)

func runSandboxDoctor(t *testing.T, setup func(fs *fakeFS)) []domain.Diagnostic {
	t.Helper()
	svc, fs, _, _, _ := newDoctorService(t)
	setup(fs)
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	return resp.Report.Items
}

func writeDoctorFile(t *testing.T, fs *fakeFS, rel, body string) {
	t.Helper()
	if err := fs.WriteFile(filepath.Join(doctorBaseDir, rel), []byte(body), 0o644); err != nil {
		t.Fatalf("seed %s: %v", rel, err)
	}
}

// LH-FA-DEV-007: host-determinable prerequisites.
func TestDoctor_SandboxRuntime(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		yaml     string
		linux    bool
		fuse     bool
		want     domain.Severity
		contains string
	}{
		{"no devcontainer config", "schemaVersion: 1\nproject:\n  name: demo\n", true, false, domain.SeverityOK, "skipped"},
		{"sandbox without nested runtime", sandboxYAMLHead + "  profile: sandbox\n", true, false, domain.SeverityOK, "nothing to check"},
		{"podman without sandbox profile", sandboxYAMLHead + "  sandbox:\n    nestedRuntime: podman\n", true, true, domain.SeverityWarn, "no effect"},
		{"fuse available", sandboxYAMLPodman, true, true, domain.SeverityOK, "/dev/fuse is available"},
		{"fuse missing, warn policy", sandboxYAMLPodman, true, false, domain.SeverityWarn, "vfs"},
		{"fuse missing, fail policy", sandboxYAMLPodman + "    onUnavailable: fail\n", true, false, domain.SeverityError, "`fail`"},
		{"non-linux host cannot tell", sandboxYAMLPodman + "    onUnavailable: fail\n", false, false, domain.SeverityOK, "cannot be determined"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			items := runSandboxDoctor(t, func(fs *fakeFS) {
				writeDoctorFile(t, fs, "u-boot.yaml", tc.yaml)
				if tc.linux {
					_ = fs.WriteFile("/proc/sys/kernel/osrelease", []byte("6.8\n"), 0o444)
				}
				if tc.fuse {
					_ = fs.WriteFile("/dev/fuse", []byte{}, 0o666)
				}
			})
			d := findDiagnostic(t, items, "devcontainer.sandbox.runtime")
			if d.Severity != tc.want || !strings.Contains(d.Message, tc.contains) {
				t.Errorf("got %v %q, want %v containing %q", d.Severity, d.Message, tc.want, tc.contains)
			}
		})
	}
}

// LH-FA-DEV-009: plaintext credentials and the https token source.
func TestDoctor_SandboxCredentials(t *testing.T) {
	t.Parallel()
	gitCfg := func(url string) string { return "[remote \"origin\"]\n\turl = " + url + "\n" }
	withToken := `{"remoteEnv": {"GIT_TOKEN": "${localEnv:GIT_TOKEN}"}}`
	cases := []struct {
		name     string
		yaml     string
		files    map[string]string
		want     domain.Severity
		contains string
	}{
		{"profile default skips", sandboxYAMLHead, nil, domain.SeverityOK, "not active"},
		{"ssh origin ok", sandboxYAMLHead + "  profile: sandbox\n",
			map[string]string{".git/config": gitCfg("git@github.com:o/r.git")}, domain.SeverityOK, "no https clone"},
		{"https origin with token source", sandboxYAMLHead + "  profile: sandbox\n",
			map[string]string{".git/config": gitCfg("https://github.com/o/r.git"), ".devcontainer/devcontainer.json": withToken},
			domain.SeverityOK, "runtime token source is configured"},
		{"https origin without token source", sandboxYAMLHead + "  profile: sandbox\n",
			map[string]string{".git/config": gitCfg("https://github.com/o/r.git"), ".devcontainer/devcontainer.json": "{}"},
			domain.SeverityWarn, "no runtime token source"},
		{"token in devcontainer.json", sandboxYAMLHead + "  profile: sandbox\n",
			map[string]string{".devcontainer/devcontainer.json": `{"postCreateCommand": "git clone https://ghp_abcdefghijklmnopqrstuvwxyz0123456789@github.com/o/r.git"}`},
			domain.SeverityWarn, ".devcontainer/devcontainer.json"},
		{"token assignment in .env.example", sandboxYAMLHead + "  profile: sandbox\n",
			map[string]string{".env.example": "GIT_TOKEN=abcdef123456789xyz\n"},
			domain.SeverityWarn, ".env.example"},
		{"env reference is not a credential", sandboxYAMLHead + "  profile: sandbox\n",
			map[string]string{".devcontainer/devcontainer.json": withToken}, domain.SeverityOK, "no https clone"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			items := runSandboxDoctor(t, func(fs *fakeFS) {
				writeDoctorFile(t, fs, "u-boot.yaml", tc.yaml)
				for rel, body := range tc.files {
					writeDoctorFile(t, fs, rel, body)
				}
			})
			d := findDiagnostic(t, items, "devcontainer.sandbox.credentials")
			if d.Severity != tc.want || !strings.Contains(d.Message, tc.contains) {
				t.Errorf("got %v %q, want %v containing %q", d.Severity, d.Message, tc.want, tc.contains)
			}
		})
	}
}
