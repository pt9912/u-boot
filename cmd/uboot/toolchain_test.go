package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// LH-FA-DEV-003 (pinned toolchain, SPEC-054), LH-NFA-SEC-004,
// LH-FA-CLI-005 at the user path (real adapters): `generate devcontainer`
// gives one supply-chain notice per active key of the pinned toolchain, the
// Dockerfile carries the block pin, and `--quiet` drops the notices (no
// SecurityRelaxation). The go-toolchain key needs no source read (default
// pin), so the notices ask no network.
func TestRun_Generate_ToolchainNotices(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	const head = "schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n"
	for name, profile := range map[string]string{"default profile": "", "sandbox profile": "  profile: sandbox\n"} {
		t.Run(name, func(t *testing.T) {
			writeFeatureProject(t, head+profile)
			for _, key := range []string{"go-toolchain", "golangci-lint", "java-toolchain"} {
				var out, errOut bytes.Buffer
				if code := run(context.Background(), []string{"config", "set", "devcontainer.features." + key + ".enabled", "true"}, &out, &errOut); code != 0 {
					t.Fatalf("config set %s: exit %d, stdout=%q stderr=%q", key, code, out.String(), errOut.String())
				}
			}
			var out, errOut bytes.Buffer
			if code := run(context.Background(), []string{"generate", "devcontainer"}, &out, &errOut); code != 0 {
				t.Fatalf("generate devcontainer: exit %d; stdout=%q stderr=%q", code, out.String(), errOut.String())
			}
			for _, want := range []string{
				"Warning: devcontainer feature go-toolchain installs",
				"Warning: devcontainer feature golangci-lint installs",
				"Warning: devcontainer feature java-toolchain installs",
			} {
				if n := strings.Count(out.String(), want); n != 1 {
					t.Errorf("stdout contains %q %d times, want 1:\n%s", want, n, out.String())
				}
			}
			df, err := os.ReadFile(filepath.Join(".devcontainer", "Dockerfile"))
			if err != nil {
				t.Fatal(err)
			}
			for _, want := range []string{"ARG GO_VERSION=", "ARG GOLANGCI_LINT_VERSION=", "ARG JDK_VERSION="} {
				if !strings.Contains(string(df), want) {
					t.Errorf("Dockerfile lacks %q", want)
				}
			}

			out.Reset()
			errOut.Reset()
			if code := run(context.Background(), []string{"--quiet", "generate", "devcontainer"}, &out, &errOut); code != 0 {
				t.Fatalf("--quiet generate devcontainer: exit %d; stdout=%q stderr=%q", code, out.String(), errOut.String())
			}
			if got := strings.ToLower(out.String() + errOut.String()); strings.Contains(got, "feature go-toolchain") || strings.Contains(got, "feature golangci-lint") || strings.Contains(got, "feature java-toolchain") {
				t.Errorf("--quiet must drop the notices of the keys:\nstdout=%q\nstderr=%q", out.String(), errOut.String())
			}
		})
	}
}

// LH-FA-DEV-003 (SPEC-054, S8) at the user path: the enabled external `go`
// feature next to the active `go-toolchain` key is a fachlicher
// Validierungsfehler (exit 10) and nothing is written; the same project with
// the conflict removed generates.
func TestRun_Generate_ToolchainConflictExits10(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	writeFeatureProject(t, "schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n  features:\n    go:\n      enabled: true\n")
	for _, key := range []string{"go-toolchain", "golangci-lint"} {
		var out, errOut bytes.Buffer
		if code := run(context.Background(), []string{"config", "set", "devcontainer.features." + key + ".enabled", "true"}, &out, &errOut); code != 0 {
			t.Fatalf("config set %s: exit %d, stdout=%q stderr=%q", key, code, out.String(), errOut.String())
		}
	}
	var out, errOut bytes.Buffer
	if code := run(context.Background(), []string{"generate", "devcontainer"}, &out, &errOut); code != 10 {
		t.Fatalf("generate devcontainer: exit %d, want 10; stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	if _, err := os.Stat(".devcontainer"); !os.IsNotExist(err) {
		t.Errorf("the rejected run must write nothing; .devcontainer exists (%v)", err)
	}
	var out2, errOut2 bytes.Buffer
	if code := run(context.Background(), []string{"generate", "devcontainer", "--dry-run", "--json"}, &out2, &errOut2); code != 10 {
		t.Fatalf("generate --dry-run --json: exit %d, want 10; stdout=%q stderr=%q", code, out2.String(), errOut2.String())
	}
	if !strings.Contains(out2.String(), "LH-FA-DEV-003") {
		t.Errorf("the rejection envelope carries the code LH-FA-DEV-003, got:\n%s", out2.String())
	}
	if !strings.Contains(out2.String(), "install two Go toolchains") {
		t.Errorf("the rejection names the conflict, got:\n%s", out2.String())
	}

	// Guard: with go-toolchain switched off the same project generates.
	var out3, errOut3 bytes.Buffer
	if code := run(context.Background(), []string{"config", "set", "devcontainer.features.go-toolchain.enabled", "false"}, &out3, &errOut3); code != 0 {
		t.Fatalf("config set go-toolchain.enabled false: exit %d\n%s\n%s", code, out3.String(), errOut3.String())
	}
	var out4, errOut4 bytes.Buffer
	if code := run(context.Background(), []string{"generate", "devcontainer"}, &out4, &errOut4); code != 0 {
		t.Fatalf("generate after the conflict was removed: exit %d\n%s\n%s", code, out4.String(), errOut4.String())
	}
}

// LH-FA-DEV-003 (SPEC-054, Muster Go-Fall) at the user path: the enabled
// external `java` feature next to the active `java-toolchain` key is a
// fachlicher Validierungsfehler (exit 10) and nothing is written; the same
// project with the conflict removed generates.
func TestRun_Generate_JavaToolchainConflictExits10(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	writeFeatureProject(t, "schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n  features:\n    java:\n      enabled: true\n")
	var out, errOut bytes.Buffer
	if code := run(context.Background(), []string{"config", "set", "devcontainer.features.java-toolchain.enabled", "true"}, &out, &errOut); code != 0 {
		t.Fatalf("config set java-toolchain: exit %d, stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	var out2, errOut2 bytes.Buffer
	if code := run(context.Background(), []string{"generate", "devcontainer"}, &out2, &errOut2); code != 10 {
		t.Fatalf("generate devcontainer: exit %d, want 10; stdout=%q stderr=%q", code, out2.String(), errOut2.String())
	}
	if _, err := os.Stat(".devcontainer"); !os.IsNotExist(err) {
		t.Errorf("the rejected run must write nothing; .devcontainer exists (%v)", err)
	}
	var out3, errOut3 bytes.Buffer
	if code := run(context.Background(), []string{"generate", "devcontainer", "--dry-run", "--json"}, &out3, &errOut3); code != 10 {
		t.Fatalf("generate --dry-run --json: exit %d, want 10; stdout=%q stderr=%q", code, out3.String(), errOut3.String())
	}
	if !strings.Contains(out3.String(), "LH-FA-DEV-003") {
		t.Errorf("the rejection envelope carries the code LH-FA-DEV-003, got:\n%s", out3.String())
	}
	if !strings.Contains(out3.String(), "install two JDKs") {
		t.Errorf("the rejection names the conflict, got:\n%s", out3.String())
	}

	// Guard: with java-toolchain switched off the same project generates.
	var out4, errOut4 bytes.Buffer
	if code := run(context.Background(), []string{"config", "set", "devcontainer.features.java-toolchain.enabled", "false"}, &out4, &errOut4); code != 0 {
		t.Fatalf("config set java-toolchain.enabled false: exit %d\n%s\n%s", code, out4.String(), errOut4.String())
	}
	var out5, errOut5 bytes.Buffer
	if code := run(context.Background(), []string{"generate", "devcontainer"}, &out5, &errOut5); code != 0 {
		t.Fatalf("generate after the conflict was removed: exit %d\n%s\n%s", code, out5.String(), errOut5.String())
	}
}
