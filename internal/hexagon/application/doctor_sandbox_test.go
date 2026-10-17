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

// LH-FA-CLI-005 exceptions 1 and 2: each warn item of the sandbox
// checks carries exactly one of the two `--quiet` flags or none. The
// flags sit on the item, not on the ID: `devcontainer.sandbox.runtime`
// and `devcontainer.sandbox.egress` also emit items of the other class
// or of none. `vfs` is a degradation fallback (exception 2); the
// credential and egress warn items are security relaxations (exception
// 1); the `nestedRuntime` "no effect" item is a configuration message
// and carries neither flag; the `fail`-policy error carries none (its
// Severity shows it). Host-side, only the `vfs` fallback is
// determinable; the daemon and egress fallbacks arise at container
// start.
func TestDoctor_SandboxFallbackFlagIsOnTheItem(t *testing.T) {
	t.Parallel()
	gitCfg := func(url string) string { return "[remote \"origin\"]\n\turl = " + url + "\n" }
	cases := []struct {
		name       string
		yaml       string
		files      map[string]string
		id         string
		want       domain.Severity
		fallback   bool
		relaxation bool
	}{
		{"runtime: vfs fallback (warn policy)", sandboxYAMLPodman, nil, "devcontainer.sandbox.runtime", domain.SeverityWarn, true, false},
		{"runtime: fuse missing, fail policy is an error", sandboxYAMLPodman + "    onUnavailable: fail\n", nil, "devcontainer.sandbox.runtime", domain.SeverityError, false, false},
		{"runtime: no effect without sandbox, same ID", sandboxYAMLHead + "  sandbox:\n    nestedRuntime: podman\n", nil, "devcontainer.sandbox.runtime", domain.SeverityWarn, false, false},
		{"runtime: nested docker is OK and stays unflagged", sandboxYAMLDocker, nil, "devcontainer.sandbox.runtime", domain.SeverityOK, false, false},
		{"egress: no effect without sandbox", sandboxYAMLHead + "  sandbox:\n    egress:\n      enabled: true\n", nil, "devcontainer.sandbox.egress", domain.SeverityWarn, false, true},
		{"egress: bypass with nested docker", sandboxYAMLDocker + "    egress:\n      enabled: true\n", nil, "devcontainer.sandbox.egress", domain.SeverityWarn, false, true},
		{"egress: configured is OK and stays unflagged", sandboxYAMLHead + "  profile: sandbox\n  sandbox:\n    egress:\n      enabled: true\n", nil, "devcontainer.sandbox.egress", domain.SeverityOK, false, false},
		{"credentials: plaintext key", sandboxYAMLHead + "  profile: sandbox\n",
			map[string]string{".env.example": "GIT_TOKEN=abcdef123456789xyz\n"}, "devcontainer.sandbox.credentials", domain.SeverityWarn, false, true},
		{"credentials: https clone without token source", sandboxYAMLHead + "  profile: sandbox\n",
			map[string]string{".git/config": gitCfg("https://github.com/o/r.git"), ".devcontainer/devcontainer.json": "{}"},
			"devcontainer.sandbox.credentials", domain.SeverityWarn, false, true},
		{"credentials: https clone with token source is OK and stays unflagged", sandboxYAMLHead + "  profile: sandbox\n",
			map[string]string{".git/config": gitCfg("https://github.com/o/r.git"), ".devcontainer/devcontainer.json": `{"remoteEnv": {"GIT_TOKEN": "${localEnv:GIT_TOKEN}"}}`},
			"devcontainer.sandbox.credentials", domain.SeverityOK, false, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			items := runSandboxDoctor(t, func(fs *fakeFS) {
				writeDoctorFile(t, fs, "u-boot.yaml", tc.yaml)
				for rel, body := range tc.files {
					writeDoctorFile(t, fs, rel, body)
				}
				_ = fs.WriteFile("/proc/sys/kernel/osrelease", []byte("6.8\n"), 0o444)
			})
			d := findDiagnostic(t, items, tc.id)
			if d.Severity != tc.want {
				t.Fatalf("severity = %v, want %v (%q)", d.Severity, tc.want, d.Message)
			}
			if d.DegradationFallback != tc.fallback {
				t.Errorf("DegradationFallback = %v, want %v (%q)", d.DegradationFallback, tc.fallback, d.Message)
			}
			if d.SecurityRelaxation != tc.relaxation {
				t.Errorf("SecurityRelaxation = %v, want %v (%q)", d.SecurityRelaxation, tc.relaxation, d.Message)
			}
		})
	}
}

// The plaintext-credential warn item is no degradation fallback but a
// security relaxation (LH-FA-DEV-009, LH-FA-CLI-005 exception 1): it is
// classified, so `doctor --quiet` shows it, yet it is not a fallback.
func TestDoctor_SandboxCredentialsWarnIsNoFallback(t *testing.T) {
	t.Parallel()
	items := runSandboxDoctor(t, func(fs *fakeFS) {
		writeDoctorFile(t, fs, "u-boot.yaml", sandboxYAMLHead+"  profile: sandbox\n")
		writeDoctorFile(t, fs, ".env.example", "GIT_TOKEN=abcdef123456789xyz\n")
	})
	d := findDiagnostic(t, items, "devcontainer.sandbox.credentials")
	if d.Severity != domain.SeverityWarn || d.DegradationFallback || !d.SecurityRelaxation {
		t.Errorf("got %v fallback=%v relaxation=%v, want a warn item that is a relaxation and no fallback",
			d.Severity, d.DegradationFallback, d.SecurityRelaxation)
	}
}

// Negative against over-flagging: the warn items of the feature checks
// (allowlist: no catalogue entry and no source; drift: active feature
// missing in devcontainer.json) concern configuration and provenance,
// not a security property of the sandbox profile, and carry no flag
// (LH-FA-CLI-005 exception 1), even when the sandbox profile is active.
func TestDoctor_FeatureWarnItemsCarryNoQuietFlag(t *testing.T) {
	t.Parallel()
	svc, fs, _, _, _ := newDoctorService(t)
	seedDoctorUbootYAMLFeatures(t, fs, `schemaVersion: 1
project:
  name: demo
devcontainer:
  enabled: true
  profile: sandbox
  features:
    ghost:
      enabled: true
    node:
      enabled: true
`)
	writeDoctorFile(t, fs, ".devcontainer/devcontainer.json", "{\"features\": {}}")
	resp, err := svc.Check(context.Background(), driving.DoctorRequest{BaseDir: doctorBaseDir})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	for _, id := range []string{"devcontainer.features.allowlist", "devcontainer.features.drift"} {
		d := findDiagnostic(t, resp.Report.Items, id)
		if d.Severity != domain.SeverityWarn {
			t.Fatalf("%s: severity = %v, want warn (%q)", id, d.Severity, d.Message)
		}
		if d.SecurityRelaxation || d.DegradationFallback {
			t.Errorf("%s: carries a --quiet flag: %+v", id, d)
		}
	}
}

// Every OK item of the sandbox checks stays unflagged: the flag marks
// warn items only (a flag on an OK item would not change the
// Severity but would surface it under `--quiet`).
func TestDoctor_SandboxOKItemsCarryNoQuietFlag(t *testing.T) {
	t.Parallel()
	items := runSandboxDoctor(t, func(fs *fakeFS) {
		writeDoctorFile(t, fs, "u-boot.yaml", sandboxYAMLDocker+"    egress:\n      enabled: true\n")
	})
	for _, d := range items {
		if d.Severity == domain.SeverityOK && (d.SecurityRelaxation || d.DegradationFallback) {
			t.Errorf("OK item %s carries a --quiet flag: %q", d.ID, d.Message)
		}
	}
}
