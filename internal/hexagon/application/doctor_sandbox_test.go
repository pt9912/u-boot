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
	sandboxYAMLDocker = sandboxYAMLHead + "  profile: sandbox\n  sandbox:\n    nestedRuntime: docker\n"
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
		restrict bool
		want     domain.Severity
		contains string
	}{
		{"no devcontainer config", "schemaVersion: 1\nproject:\n  name: demo\n", true, false, false, domain.SeverityOK, "skipped"},
		{"sandbox without nested runtime", sandboxYAMLHead + "  profile: sandbox\n", true, false, false, domain.SeverityOK, "nothing to check"},
		{"podman without sandbox profile", sandboxYAMLHead + "  sandbox:\n    nestedRuntime: podman\n", true, true, false, domain.SeverityWarn, "no effect"},
		{"docker without sandbox profile", sandboxYAMLHead + "  sandbox:\n    nestedRuntime: docker\n", true, true, false, domain.SeverityWarn, "no effect"},
		{"docker with sandbox profile", sandboxYAMLDocker, true, false, false, domain.SeverityOK, "--privileged"},
		{"fuse available", sandboxYAMLPodman, true, true, false, domain.SeverityOK, "/dev/fuse is available"},
		{"fuse missing, warn policy", sandboxYAMLPodman, true, false, false, domain.SeverityWarn, "vfs"},
		{"fuse missing, fail policy", sandboxYAMLPodman + "    onUnavailable: fail\n", true, false, false, domain.SeverityError, "`fail`"},
		{"userns restricted by AppArmor sysctl", sandboxYAMLPodman, true, true, true, domain.SeverityError, "apparmor_restrict_unprivileged_userns"},
		{"non-linux host cannot tell", sandboxYAMLPodman + "    onUnavailable: fail\n", false, false, false, domain.SeverityOK, "cannot be determined"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			items := runSandboxDoctor(t, func(fs *fakeFS) {
				writeDoctorFile(t, fs, "u-boot.yaml", tc.yaml)
				if tc.linux {
					_ = fs.WriteFile("/proc/sys/kernel/osrelease", []byte("6.8\n"), 0o444)
				}
				if tc.restrict {
					_ = fs.WriteFile("/proc/sys/kernel/apparmor_restrict_unprivileged_userns", []byte("1\n"), 0o644)
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

// LH-FA-DEV-008: static consistency check.
func TestDoctor_SandboxEgress(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		yaml     string
		want     domain.Severity
		contains string
	}{
		{"off", sandboxYAMLHead, domain.SeverityOK, "not enabled"},
		{"on without sandbox", sandboxYAMLHead + "  sandbox:\n    egress:\n      enabled: true\n", domain.SeverityWarn, "no effect"},
		{"on with sandbox", sandboxYAMLHead + "  profile: sandbox\n  sandbox:\n    egress:\n      enabled: true\n", domain.SeverityOK, "verified when the container starts"},
		{"on with nested docker", sandboxYAMLDocker + "    egress:\n      enabled: true\n", domain.SeverityWarn, "bypass the restriction"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			items := runSandboxDoctor(t, func(fs *fakeFS) { writeDoctorFile(t, fs, "u-boot.yaml", tc.yaml) })
			d := findDiagnostic(t, items, "devcontainer.sandbox.egress")
			if d.Severity != tc.want || !strings.Contains(d.Message, tc.contains) {
				t.Errorf("got %v %q, want %v containing %q", d.Severity, d.Message, tc.want, tc.contains)
			}
		})
	}
}

// LH-FA-CLI-005 exception 2 / LH-FA-DEV-007: exactly the `vfs`
// fallback item of the runtime check carries DegradationFallback
// (`doctor --quiet` shows it next to the errors). The flag sits on the
// item, not on the ID: the other warn items of the same ID and of the
// egress and credentials checks stay unflagged, and so does the
// `fail`-policy error. Host-side, only the `vfs` fallback is
// determinable; the daemon and egress fallbacks arise at container
// start.
func TestDoctor_SandboxFallbackFlagIsOnTheItem(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		yaml    string
		id      string
		want    domain.Severity
		flagged bool
	}{
		{"runtime: vfs fallback (warn policy)", sandboxYAMLPodman, "devcontainer.sandbox.runtime", domain.SeverityWarn, true},
		{"runtime: fuse missing, fail policy is an error", sandboxYAMLPodman + "    onUnavailable: fail\n", "devcontainer.sandbox.runtime", domain.SeverityError, false},
		{"runtime: no effect without sandbox, same ID", sandboxYAMLHead + "  sandbox:\n    nestedRuntime: podman\n", "devcontainer.sandbox.runtime", domain.SeverityWarn, false},
		{"egress: no effect without sandbox", sandboxYAMLHead + "  sandbox:\n    egress:\n      enabled: true\n", "devcontainer.sandbox.egress", domain.SeverityWarn, false},
		{"egress: bypass with nested docker", sandboxYAMLDocker + "    egress:\n      enabled: true\n", "devcontainer.sandbox.egress", domain.SeverityWarn, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			items := runSandboxDoctor(t, func(fs *fakeFS) {
				writeDoctorFile(t, fs, "u-boot.yaml", tc.yaml)
				_ = fs.WriteFile("/proc/sys/kernel/osrelease", []byte("6.8\n"), 0o444)
			})
			d := findDiagnostic(t, items, tc.id)
			if d.Severity != tc.want {
				t.Fatalf("severity = %v, want %v (%q)", d.Severity, tc.want, d.Message)
			}
			if d.DegradationFallback != tc.flagged {
				t.Errorf("DegradationFallback = %v, want %v (%q)", d.DegradationFallback, tc.flagged, d.Message)
			}
		})
	}
}

// The plaintext-credential warn item is no degradation fallback
// (LH-FA-DEV-009), so `doctor --quiet` drops it.
func TestDoctor_SandboxCredentialsWarnIsNoFallback(t *testing.T) {
	t.Parallel()
	items := runSandboxDoctor(t, func(fs *fakeFS) {
		writeDoctorFile(t, fs, "u-boot.yaml", sandboxYAMLHead+"  profile: sandbox\n")
		writeDoctorFile(t, fs, ".env.example", "GIT_TOKEN=abcdef123456789xyz\n")
	})
	d := findDiagnostic(t, items, "devcontainer.sandbox.credentials")
	if d.Severity != domain.SeverityWarn || d.DegradationFallback {
		t.Errorf("got %v fallback=%v, want an unflagged warn item", d.Severity, d.DegradationFallback)
	}
}
