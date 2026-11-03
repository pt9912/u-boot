package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The project environment key `devcontainer.containerEnv` at the user
// path (lastenheft 0.32.0, LH-FA-CONF-005, LH-FA-CLI-006): the entry
// `config set` takes is validated (name against the rule, value a
// string), the base path is no known path, the fachliche Fehler end
// with exit code 10 and nothing written, `generate devcontainer`
// discloses each rendered value, and a colliding cache value ends the
// generation with exit code 10.

func TestRun_ConfigSet_ContainerEnvValueAndNames(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	writeFeatureProject(t, "schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n")
	var out, errOut bytes.Buffer
	if code := run(context.Background(), []string{"config", "set", "devcontainer.containerEnv.BUILD_MODE", "local"}, &out, &errOut); code != 0 {
		t.Fatalf("config set BUILD_MODE: exit %d, stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	out.Reset()
	errOut.Reset()
	if code := run(context.Background(), []string{"config", "get", "devcontainer.containerEnv.BUILD_MODE"}, &out, &errOut); code != 0 {
		t.Fatalf("config get BUILD_MODE: exit %d, stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	if got := strings.TrimSpace(out.String()); got != "local" {
		t.Errorf("config get printed %q, want local", got)
	}

	// An invalid name ends with exit code 10 and writes nothing.
	body, err := os.ReadFile(filepath.Join(".", "u-boot.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	out.Reset()
	errOut.Reset()
	if code := run(context.Background(), []string{"config", "set", "devcontainer.containerEnv.BAD-NAME", "x"}, &out, &errOut); code != 10 {
		t.Errorf("config set BAD-NAME: exit %d, want 10; stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	after, err := os.ReadFile(filepath.Join(".", "u-boot.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(body) {
		t.Errorf("the rejected call wrote to u-boot.yaml:\n--- before ---\n%s\n--- after ---\n%s", body, after)
	}

	// The base path is no known config path — no whole-map write: both
	// read and write end with exit code 10.
	for _, args := range [][]string{
		{"config", "set", "devcontainer.containerEnv", "map"},
		{"config", "get", "devcontainer.containerEnv"},
	} {
		out.Reset()
		errOut.Reset()
		if code := run(context.Background(), args, &out, &errOut); code != 10 {
			t.Errorf("%v: exit %d, want 10; stdout=%q stderr=%q", args, code, out.String(), errOut.String())
		}
		if !strings.Contains(out.String()+errOut.String(), "not a known config path") {
			t.Errorf("%v: the message does not name the shape, got:\n%s\n%s", args, out.String(), errOut.String())
		}
	}
}

// The disclosure (LH-FA-DEV-001): `generate devcontainer` names every
// rendered project value in the command output and carries the entries
// into the generated `containerEnv` object; `--quiet` drops the
// disclosure (a plain warning, like the other notices of the run).
func TestRun_Generate_ContainerEnvNotices(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	writeFeatureProject(t, "schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n")
	for _, entry := range [][2]string{
		{"devcontainer.containerEnv.BUILD_MODE", "local"},
		{"devcontainer.containerEnv.GOTOOLCHAIN", "local"},
	} {
		var out, errOut bytes.Buffer
		if code := run(context.Background(), []string{"config", "set", entry[0], entry[1]}, &out, &errOut); code != 0 {
			t.Fatalf("config set %s: exit %d, stdout=%q stderr=%q", entry[0], code, out.String(), errOut.String())
		}
	}
	var out, errOut bytes.Buffer
	if code := run(context.Background(), []string{"generate", "devcontainer"}, &out, &errOut); code != 0 {
		t.Fatalf("generate devcontainer: exit %d; stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	for _, name := range []string{"BUILD_MODE", "GOTOOLCHAIN"} {
		want := "Warning: project environment value " + name + "="
		if n := strings.Count(out.String(), want); n != 1 {
			t.Errorf("stdout contains %q %d times, want 1:\n%s", want, n, out.String())
		}
	}
	jsonBody, err := os.ReadFile(filepath.Join(".devcontainer", "devcontainer.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`"BUILD_MODE": "local"`,
		`"GOTOOLCHAIN": "local"`,
	} {
		if !strings.Contains(string(jsonBody), want) {
			t.Errorf("devcontainer.json lacks %q:\n%s", want, jsonBody)
		}
	}

	out.Reset()
	errOut.Reset()
	if code := run(context.Background(), []string{"--quiet", "generate", "devcontainer"}, &out, &errOut); code != 0 {
		t.Fatalf("--quiet generate devcontainer: exit %d; stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	if got := strings.ToLower(out.String() + errOut.String()); strings.Contains(got, "project environment value") {
		t.Errorf("--quiet must drop the project-value disclosure:\nstdout=%q\nstderr=%q", out.String(), errOut.String())
	}
}

// LH-FA-DEV-001 (the collision of the cache names) at the user path:
// a project value `GOCACHE` beside the active `go-toolchain` key ends
// `generate devcontainer` with exit code 10 and nothing is written;
// with the value removed by hand the same project generates.
func TestRun_Generate_ContainerEnvCollisionExits10(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	writeFeatureProject(t, "schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n")
	for _, key := range []string{"go-toolchain"} {
		var out, errOut bytes.Buffer
		if code := run(context.Background(), []string{"config", "set", "devcontainer.features." + key + ".enabled", "true"}, &out, &errOut); code != 0 {
			t.Fatalf("config set %s: exit %d, stdout=%q stderr=%q", key, code, out.String(), errOut.String())
		}
	}
	var out, errOut bytes.Buffer
	if code := run(context.Background(), []string{"config", "set", "devcontainer.containerEnv.GOCACHE", "/tmp/go-cache"}, &out, &errOut); code != 0 {
		t.Fatalf("config set GOCACHE: exit %d, stdout=%q stderr=%q", code, out.String(), errOut.String())
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
	if !strings.Contains(out3.String(), "GOCACHE") || !strings.Contains(out3.String(), "go-toolchain") {
		t.Errorf("the rejection names the entry and the active key, got:\n%s", out3.String())
	}

	// Guard: the hand-edited config without GOCACHE generates.
	body, err := os.ReadFile(filepath.Join(".", "u-boot.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	cleaned := strings.NewReplacer("  containerEnv:\n    GOCACHE: /tmp/go-cache\n", "").Replace(string(body))
	if cleaned == string(body) {
		t.Fatalf("the fixture edit removed nothing:\n%s", body)
	}
	if err := os.WriteFile(filepath.Join(".", "u-boot.yaml"), []byte(cleaned), 0o644); err != nil {
		t.Fatal(err)
	}
	var out4, errOut4 bytes.Buffer
	if code := run(context.Background(), []string{"generate", "devcontainer"}, &out4, &errOut4); code != 0 {
		t.Fatalf("generate after the collision was removed: exit %d\n%s\n%s", code, out4.String(), errOut4.String())
	}
}
