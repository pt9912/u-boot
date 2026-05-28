package application_test

import (
	"context"
	"errors"
	iofs "io/fs"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/application/managedblock"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// --- M5-T4c: Templates + executeAdd + Active-Repair + End-to-End ----

func TestT4c_TemplateNames_IncludesServices(t *testing.T) {
	names, err := application.TemplateNamesForTest()
	if err != nil {
		t.Fatalf("templateNames: %v", err)
	}
	for _, want := range []string{
		"services/postgres.compose.tmpl",
		"services/postgres.volume.tmpl",
		"services/postgres.env.tmpl",
	} {
		found := false
		for _, n := range names {
			if n == want {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("template %q missing from templateNames() = %v", want, names)
		}
	}
}

func TestT4c_RenderPostgresTemplates(t *testing.T) {
	for _, name := range []string{
		"services/postgres.compose.tmpl",
		"services/postgres.volume.tmpl",
		"services/postgres.env.tmpl",
	} {
		t.Run(name, func(t *testing.T) {
			body, err := application.RenderTemplateForTest(name, "postgres")
			if err != nil {
				t.Fatalf("RenderTemplate(%s): %v", name, err)
			}
			if len(body) == 0 {
				t.Errorf("template %s rendered empty", name)
			}
		})
	}
}

// --- executeAdd happy-paths --------------------------------------

func TestAdd_T4c_Register_WritesAllThreeFiles(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs, "schemaVersion: 1\nproject:\n  name: demo\n")
	// no compose.yaml, no .env.example → bootstrap + create paths

	resp, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if resp.State != domain.ServiceStateActive {
		t.Errorf("State = %s, want active", resp.State.String())
	}

	// All three files written.
	for _, p := range []string{"u-boot.yaml", "compose.yaml", ".env.example"} {
		body, err := fs.ReadFile(filepath.Join(addTestBaseDir, p))
		if err != nil {
			t.Errorf("expected %s written: %v", p, err)
			continue
		}
		if len(body) == 0 {
			t.Errorf("%s written empty", p)
		}
	}

	// u-boot.yaml enabled flipped.
	body, _ := fs.ReadFile(filepath.Join(addTestBaseDir, "u-boot.yaml"))
	if !strings.Contains(string(body), "enabled: true") {
		t.Errorf("u-boot.yaml missing enabled: true; got:\n%s", body)
	}

	// compose.yaml has both managed blocks under the right anchors.
	composeBody, _ := fs.ReadFile(filepath.Join(addTestBaseDir, "compose.yaml"))
	if !strings.Contains(string(composeBody), "BEGIN U-BOOT MANAGED BLOCK: service.postgres") {
		t.Errorf("service block missing; got:\n%s", composeBody)
	}
	if !strings.Contains(string(composeBody), "BEGIN U-BOOT MANAGED BLOCK: volume.postgres") {
		t.Errorf("volume block missing; got:\n%s", composeBody)
	}

	// .env.example has the wrapped POSTGRES_* keys.
	envBody, _ := fs.ReadFile(filepath.Join(addTestBaseDir, ".env.example"))
	for _, want := range []string{
		"BEGIN U-BOOT MANAGED BLOCK: service.postgres",
		"POSTGRES_USER=postgres",
		"POSTGRES_PASSWORD=CHANGEME_POSTGRES_PASSWORD",
		"POSTGRES_DB=postgres",
	} {
		if !strings.Contains(string(envBody), want) {
			t.Errorf(".env.example missing %q; got:\n%s", want, envBody)
		}
	}
}

func TestAdd_T4c_RebuildBlock_SkipsUBootYAMLWrite(t *testing.T) {
	// InconsistentBlock: enabled: true already; only compose + env
	// should be written.
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs,
		"schemaVersion: 1\nproject:\n  name: demo\n"+
			"services:\n  postgres:\n    enabled: true\n")

	resp, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	for _, p := range resp.Changed {
		if p == "u-boot.yaml" {
			t.Errorf("u-boot.yaml in Changed for rebuild path; got %v", resp.Changed)
		}
	}
}

// --- Active-Repair --------------------------------------------------

func TestAdd_T4c_Active_StaleServiceMissingImage_Repairs(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs,
		"schemaVersion: 1\nproject:\n  name: demo\n"+
			"services:\n  postgres:\n    enabled: true\n")
	// Service block exists but has no image: line.
	seedCompose(t, fs,
		"services:\n"+
			"  # BEGIN U-BOOT MANAGED BLOCK: service.postgres\n"+
			"  postgres:\n"+
			"    environment:\n"+
			"      POSTGRES_USER: a\n"+
			"      POSTGRES_PASSWORD: b\n"+
			"      POSTGRES_DB: c\n"+
			"    volumes:\n"+
			"      - postgres-data:/x\n"+
			"    ports:\n"+
			"      - \"5432:5432\"\n"+
			"    healthcheck:\n"+
			"      test: foo\n"+
			"  # END U-BOOT MANAGED BLOCK: service.postgres\n"+
			"\n"+
			"volumes:\n"+
			"  # BEGIN U-BOOT MANAGED BLOCK: volume.postgres\n"+
			"  postgres-data: {}\n"+
			"  # END U-BOOT MANAGED BLOCK: volume.postgres\n")
	seedEnv(t, fs, envBlockComplete("postgres"))

	resp, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if resp.PriorState != domain.ServiceStateActive {
		t.Errorf("PriorState = %s, want active", resp.PriorState.String())
	}
	containsCompose := false
	for _, p := range resp.Changed {
		if p == "compose.yaml" {
			containsCompose = true
		}
		if p == "u-boot.yaml" {
			t.Errorf("u-boot.yaml unexpectedly in Changed = %v", resp.Changed)
		}
	}
	if !containsCompose {
		t.Errorf("compose.yaml not repaired; Changed = %v", resp.Changed)
	}
}

func TestAdd_T4c_Active_MissingVolume_RepairsCompose(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs,
		"schemaVersion: 1\nproject:\n  name: demo\n"+
			"services:\n  postgres:\n    enabled: true\n")
	// Service complete, volume block missing entirely.
	seedCompose(t, fs,
		strings.ReplaceAll(composeBlockComplete("postgres"),
			"\nvolumes:\n  # BEGIN U-BOOT MANAGED BLOCK: volume.postgres\n  postgres-data: {}\n  # END U-BOOT MANAGED BLOCK: volume.postgres\n",
			""))
	seedEnv(t, fs, envBlockComplete("postgres"))

	resp, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if len(resp.Changed) != 1 || resp.Changed[0] != "compose.yaml" {
		t.Errorf("expected Changed=[compose.yaml], got %v", resp.Changed)
	}
	composeBody, _ := fs.ReadFile(filepath.Join(addTestBaseDir, "compose.yaml"))
	if !strings.Contains(string(composeBody), "BEGIN U-BOOT MANAGED BLOCK: volume.postgres") {
		t.Errorf("volume block missing after repair; got:\n%s", composeBody)
	}
}

func TestAdd_T4c_Active_MissingEnvBlock_RepairsEnvOnly(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs,
		"schemaVersion: 1\nproject:\n  name: demo\n"+
			"services:\n  postgres:\n    enabled: true\n")
	seedCompose(t, fs, composeBlockComplete("postgres"))
	// .env.example missing entirely.

	resp, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if len(resp.Changed) != 1 || resp.Changed[0] != ".env.example" {
		t.Errorf("expected Changed=[.env.example], got %v", resp.Changed)
	}
}

func TestAdd_T4c_Active_UserCustomPortStaysNoOp(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs,
		"schemaVersion: 1\nproject:\n  name: demo\n"+
			"services:\n  postgres:\n    enabled: true\n")
	// Customised port — all required fields still present.
	seedCompose(t, fs,
		strings.Replace(composeBlockComplete("postgres"),
			"\"5432:5432\"",
			"\"5433:5432\"", 1))
	seedEnv(t, fs, envBlockComplete("postgres"))

	resp, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if resp.Changed != nil {
		t.Errorf("expected nil Changed (no-op despite custom port), got %v", resp.Changed)
	}
}

// --- False-positive guards -----------------------------------------

func TestAdd_T4c_Active_CommentedPostgresPasswordDoesNotCount(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs,
		"schemaVersion: 1\nproject:\n  name: demo\n"+
			"services:\n  postgres:\n    enabled: true\n")
	seedCompose(t, fs, composeBlockComplete("postgres"))
	// .env.example has commented POSTGRES_PASSWORD only.
	seedEnv(t, fs,
		"# BEGIN U-BOOT MANAGED BLOCK: service.postgres\n"+
			"POSTGRES_USER=a\n"+
			"# POSTGRES_PASSWORD=secret\n"+
			"POSTGRES_DB=c\n"+
			"# END U-BOOT MANAGED BLOCK: service.postgres\n")

	resp, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if len(resp.Changed) != 1 || resp.Changed[0] != ".env.example" {
		t.Errorf("expected env repair, got Changed=%v", resp.Changed)
	}
}

func TestAdd_T4c_Active_PostgresUserUnderLabelsDoesNotCount(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs,
		"schemaVersion: 1\nproject:\n  name: demo\n"+
			"services:\n  postgres:\n    enabled: true\n")
	// POSTGRES_USER lives under labels:, not environment:.
	seedCompose(t, fs,
		"services:\n"+
			"  # BEGIN U-BOOT MANAGED BLOCK: service.postgres\n"+
			"  postgres:\n"+
			"    image: postgres\n"+
			"    labels:\n"+
			"      POSTGRES_USER: tricky\n"+
			"    environment:\n"+
			"      POSTGRES_PASSWORD: p\n"+
			"      POSTGRES_DB: d\n"+
			"    volumes:\n"+
			"      - postgres-data:/x\n"+
			"    ports:\n"+
			"      - \"5432:5432\"\n"+
			"    healthcheck:\n"+
			"      test: foo\n"+
			"  # END U-BOOT MANAGED BLOCK: service.postgres\n"+
			"\n"+
			"volumes:\n"+
			"  # BEGIN U-BOOT MANAGED BLOCK: volume.postgres\n"+
			"  postgres-data: {}\n"+
			"  # END U-BOOT MANAGED BLOCK: volume.postgres\n")
	seedEnv(t, fs, envBlockComplete("postgres"))

	resp, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	composeRepair := false
	for _, p := range resp.Changed {
		if p == "compose.yaml" {
			composeRepair = true
		}
	}
	if !composeRepair {
		t.Errorf("expected compose repair when POSTGRES_USER is in labels:, got Changed=%v", resp.Changed)
	}
}

func TestAdd_T4c_Active_ImageEmptyValueTriggersRepair(t *testing.T) {
	// Why: `image:` with no value (or `image: ""`) is treated as no
	// image at all — pins the trimmed-non-empty rule of the scanner.
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs,
		"schemaVersion: 1\nproject:\n  name: demo\n"+
			"services:\n  postgres:\n    enabled: true\n")
	seedCompose(t, fs,
		strings.Replace(composeBlockComplete("postgres"),
			"image: postgres:16-alpine", "image: \"\"", 1))
	seedEnv(t, fs, envBlockComplete("postgres"))

	resp, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	composeRepair := false
	for _, p := range resp.Changed {
		if p == "compose.yaml" {
			composeRepair = true
		}
	}
	if !composeRepair {
		t.Errorf("expected compose repair for empty image:, got Changed=%v", resp.Changed)
	}
}

func TestAdd_T4c_Active_CommentedImageTriggersRepair(t *testing.T) {
	// Why: a commented-out `# image: postgres:16` is no image —
	// comment-stripping must drop the # POSTGRES_PASSWORD style
	// matches.
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs,
		"schemaVersion: 1\nproject:\n  name: demo\n"+
			"services:\n  postgres:\n    enabled: true\n")
	seedCompose(t, fs,
		strings.Replace(composeBlockComplete("postgres"),
			"image: postgres:16-alpine", "# image: postgres:16-alpine", 1))
	seedEnv(t, fs, envBlockComplete("postgres"))

	resp, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	composeRepair := false
	for _, p := range resp.Changed {
		if p == "compose.yaml" {
			composeRepair = true
		}
	}
	if !composeRepair {
		t.Errorf("expected compose repair for commented image:, got Changed=%v", resp.Changed)
	}
}

func TestAdd_T4c_RebuildBlock_WithCompleteEnv_OnlyComposeChanged(t *testing.T) {
	// Why: InconsistentBlock + complete .env.example must yield
	// Changed=[compose.yaml] only — the env block is byte-identical
	// after Replace, so the slot must NOT be populated.
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs,
		"schemaVersion: 1\nproject:\n  name: demo\n"+
			"services:\n  postgres:\n    enabled: true\n")
	// no compose.yaml ⇒ Bootstrap; but we also seed a complete env.
	seedEnv(t, fs, envBlockComplete("postgres"))

	resp, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	for _, p := range resp.Changed {
		if p == ".env.example" {
			t.Errorf(".env.example was rewritten despite being complete; Changed=%v", resp.Changed)
		}
	}
	if len(resp.Changed) == 0 {
		t.Errorf("expected at least compose.yaml in Changed; got %v", resp.Changed)
	}
}

func TestAdd_T4c_Active_HealthcheckDisableTrueTriggersRepair(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs,
		"schemaVersion: 1\nproject:\n  name: demo\n"+
			"services:\n  postgres:\n    enabled: true\n")
	// healthcheck: disable: true is the documented "skip" form.
	seedCompose(t, fs,
		strings.Replace(composeBlockComplete("postgres"),
			"healthcheck:\n      test: [\"CMD\", \"pg_isready\"]\n",
			"healthcheck:\n      disable: true\n", 1))
	seedEnv(t, fs, envBlockComplete("postgres"))

	resp, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	composeRepair := false
	for _, p := range resp.Changed {
		if p == "compose.yaml" {
			composeRepair = true
		}
	}
	if !composeRepair {
		t.Errorf("expected compose repair with healthcheck.disable=true, got Changed=%v", resp.Changed)
	}
}

// --- Wrong-anchor in non-Active states ------------------------------

func TestAdd_T4c_Deactivated_WrongAnchorAborts(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs,
		"schemaVersion: 1\nproject:\n  name: demo\n"+
			"services:\n  postgres:\n    enabled: false\n")
	// service.postgres marker placed under volumes:, not services:
	seedCompose(t, fs,
		"services: {}\n"+
			"\n"+
			"volumes:\n"+
			"  # BEGIN U-BOOT MANAGED BLOCK: service.postgres\n"+
			"  orphan: {}\n"+
			"  # END U-BOOT MANAGED BLOCK: service.postgres\n")

	_, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if !errors.Is(err, driving.ErrServiceInconsistent) {
		t.Fatalf("err = %v, want ErrServiceInconsistent", err)
	}
}

// --- Symlink / non-regular reject ----------------------------------

func TestAdd_T4c_ComposeSymlinkRejected(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs, "schemaVersion: 1\nproject:\n  name: demo\n")
	composePath := filepath.Join(addTestBaseDir, "compose.yaml")
	fs.markSymlink(composePath)

	writesBefore := len(fs.writtenPaths())
	_, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if !errors.Is(err, driving.ErrBackupUnsupportedKind) {
		t.Fatalf("err = %v, want ErrBackupUnsupportedKind", err)
	}
	if got := len(fs.writtenPaths()) - writesBefore; got != 0 {
		t.Errorf("writes after symlink reject = %d, want 0", got)
	}
}

func TestAdd_T4c_EnvNonRegularRejected(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs, "schemaVersion: 1\nproject:\n  name: demo\n")
	fs.markIrregular(filepath.Join(addTestBaseDir, ".env.example"))

	_, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if !errors.Is(err, driving.ErrBackupUnsupportedKind) {
		t.Fatalf("err = %v, want ErrBackupUnsupportedKind", err)
	}
}

// --- u-boot.yaml-TOCTOU ---------------------------------------------

func TestAdd_T4c_UBootYAMLTOCTOU_IsProjectNotInitialized(t *testing.T) {
	svc, fs, _ := newAddService(t)
	yamlPath := filepath.Join(addTestBaseDir, "u-boot.yaml")
	seedUBootYAML(t, fs, "schemaVersion: 1\nproject:\n  name: demo\n")
	fs.failLstatOn = yamlPath
	fs.failLstatErr = iofs.ErrNotExist
	writesBefore := len(fs.writtenPaths())

	_, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if !errors.Is(err, driving.ErrProjectNotInitialized) {
		t.Fatalf("err = %v, want ErrProjectNotInitialized", err)
	}
	// No writes after the failing call.
	if delta := len(fs.writtenPaths()) - writesBefore; delta != 0 {
		t.Errorf("writes after TOCTOU race = %d, want 0 (new writes only); writes=%v",
			delta, fs.writtenPaths())
	}
}

// --- Mode preservation ---------------------------------------------

func TestAdd_T4c_ModePreservedFromExisting(t *testing.T) {
	svc, fs, _ := newAddService(t)
	composePath := filepath.Join(addTestBaseDir, "compose.yaml")
	seedUBootYAML(t, fs,
		"schemaVersion: 1\nproject:\n  name: demo\n"+
			"services:\n  postgres:\n    enabled: true\n")
	// Existing compose has mode 0o600.
	if err := fs.WriteFile(composePath, []byte("name: demo\n"), 0o600); err != nil {
		t.Fatalf("seed: %v", err)
	}

	if _, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	}); err != nil {
		t.Fatalf("Add: %v", err)
	}

	if got := fs.fileModes[composePath]; got != 0o600 {
		t.Errorf("compose.yaml mode = %o, want 0o600 (preserved)", got)
	}
}

// --- Bootstrap ------------------------------------------------------

func TestAdd_T4c_Bootstrap_OnMissingCompose(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs, "schemaVersion: 1\nproject:\n  name: demo\n")

	if _, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	}); err != nil {
		t.Fatalf("Add: %v", err)
	}
	body, err := fs.ReadFile(filepath.Join(addTestBaseDir, "compose.yaml"))
	if err != nil {
		t.Fatalf("compose.yaml not written: %v", err)
	}
	for _, want := range []string{
		"BEGIN U-BOOT MANAGED BLOCK: init",
		"name: demo",
		"BEGIN U-BOOT MANAGED BLOCK: service.postgres",
		"BEGIN U-BOOT MANAGED BLOCK: volume.postgres",
	} {
		if !strings.Contains(string(body), want) {
			t.Errorf("bootstrap missing %q; got:\n%s", want, body)
		}
	}
}

func TestAdd_T4c_NoBootstrap_UserComposeWithoutInitMarkerPreserved(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs, "schemaVersion: 1\nproject:\n  name: demo\n")
	user := "# Custom compose from upstream\n" +
		"services:\n" +
		"  mywebapp:\n" +
		"    image: nginx\n"
	seedCompose(t, fs, user)

	if _, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	}); err != nil {
		t.Fatalf("Add: %v", err)
	}
	body, _ := fs.ReadFile(filepath.Join(addTestBaseDir, "compose.yaml"))
	got := string(body)
	if !strings.Contains(got, "# Custom compose from upstream") {
		t.Errorf("user comment lost; got:\n%s", got)
	}
	if !strings.Contains(got, "mywebapp:") {
		t.Errorf("user service lost; got:\n%s", got)
	}
	if !strings.Contains(got, "BEGIN U-BOOT MANAGED BLOCK: service.postgres") {
		t.Errorf("postgres marker missing; got:\n%s", got)
	}
}

// --- Env paths ------------------------------------------------------

func TestAdd_T4c_EnvBlock_AppendsWithSeparator(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs, "schemaVersion: 1\nproject:\n  name: demo\n")
	seedEnv(t, fs, "FOO=bar\nBAZ=qux\n")

	if _, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	}); err != nil {
		t.Fatalf("Add: %v", err)
	}
	body, _ := fs.ReadFile(filepath.Join(addTestBaseDir, ".env.example"))
	got := string(body)
	// User keys preserved.
	for _, want := range []string{"FOO=bar", "BAZ=qux", "BEGIN U-BOOT MANAGED BLOCK: service.postgres"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q; got:\n%s", want, got)
		}
	}
	// Separator: User content ends with `\n`, then exactly one blank
	// line, then BEGIN.
	if !strings.Contains(got, "BAZ=qux\n\n# BEGIN") {
		t.Errorf("missing separating blank line; got:\n%s", got)
	}
}

func TestAdd_T4c_EnvBlock_MalformedAborts(t *testing.T) {
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs, "schemaVersion: 1\nproject:\n  name: demo\n")
	// BEGIN without END
	seedEnv(t, fs,
		"# BEGIN U-BOOT MANAGED BLOCK: service.postgres\n"+
			"POSTGRES_USER=a\n")

	_, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if !errors.Is(err, driving.ErrServiceInconsistent) {
		t.Fatalf("err = %v, want ErrServiceInconsistent", err)
	}
}

func TestAdd_T4c_EnvBlock_ReplaceIdempotent(t *testing.T) {
	// Two Active-repair runs produce the same env file byte-identical.
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs,
		"schemaVersion: 1\nproject:\n  name: demo\n"+
			"services:\n  postgres:\n    enabled: true\n")
	seedCompose(t, fs, composeBlockComplete("postgres"))
	// env block exists but lacks POSTGRES_DB → repair.
	seedEnv(t, fs,
		"# BEGIN U-BOOT MANAGED BLOCK: service.postgres\n"+
			"POSTGRES_USER=a\n"+
			"POSTGRES_PASSWORD=b\n"+
			"# END U-BOOT MANAGED BLOCK: service.postgres\n")

	if _, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	}); err != nil {
		t.Fatalf("Add 1: %v", err)
	}
	body1, _ := fs.ReadFile(filepath.Join(addTestBaseDir, ".env.example"))
	// Second run: now content is complete, must be no-op.
	resp2, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	})
	if err != nil {
		t.Fatalf("Add 2: %v", err)
	}
	if resp2.Changed != nil {
		t.Errorf("second Add expected no-op, got Changed=%v", resp2.Changed)
	}
	body2, _ := fs.ReadFile(filepath.Join(addTestBaseDir, ".env.example"))
	if string(body1) != string(body2) {
		t.Errorf("env file changed between runs:\nbody1:\n%s\nbody2:\n%s", body1, body2)
	}
	// Exactly one BEGIN marker.
	if c := strings.Count(string(body1), "BEGIN U-BOOT MANAGED BLOCK: service.postgres"); c != 1 {
		t.Errorf("expected exactly one BEGIN after repair, got %d", c)
	}
}

// --- renderEnvManagedBlock wrap contract ---------------------------

func TestRenderEnvManagedBlock_WrapsAndFindsBack(t *testing.T) {
	wrapped := application.RenderEnvManagedBlockForTest("postgres",
		[]byte("POSTGRES_USER=a\nPOSTGRES_PASSWORD=b\nPOSTGRES_DB=c\n"))
	got := string(wrapped)
	if !strings.HasPrefix(got, "# BEGIN U-BOOT MANAGED BLOCK: service.postgres\n") {
		t.Errorf("missing BEGIN prefix; got:\n%s", got)
	}
	if !strings.HasSuffix(got, "# END U-BOOT MANAGED BLOCK: service.postgres\n") {
		t.Errorf("missing END suffix; got:\n%s", got)
	}
	// managedblock.Find must recognise the wrap.
	marker := managedblock.Marker{Style: managedblock.StyleHash, Name: "service.postgres"}
	if !managedblock.Has(wrapped, marker) {
		t.Errorf("managedblock.Has rejected the wrapped block; got:\n%s", got)
	}
}

// --- Adapter-end-to-end (real yaml.Codec via fake delegation) -------

func TestAdd_T4c_AdapterE2E_ProducedComposeIsManagedblockFindable(t *testing.T) {
	// Why: the fake delegates to the production adapter, so the
	// compose bytes produced by Add() must round-trip through
	// managedblock.Find for both markers. Pins that the M5-T6 CLI
	// (which reads the file back) will find them.
	svc, fs, _ := newAddService(t)
	seedUBootYAML(t, fs, "schemaVersion: 1\nproject:\n  name: demo\n")

	if _, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     addTestBaseDir,
		ServiceName: postgresName(t),
	}); err != nil {
		t.Fatalf("Add: %v", err)
	}
	body, _ := fs.ReadFile(filepath.Join(addTestBaseDir, "compose.yaml"))
	for _, marker := range []managedblock.Marker{
		{Style: managedblock.StyleHash, Name: "service.postgres"},
		{Style: managedblock.StyleHash, Name: "volume.postgres"},
	} {
		if _, _, err := managedblock.Find(body, marker); err != nil {
			t.Errorf("managedblock.Find(%s) failed: %v; body:\n%s",
				marker.Name, err, body)
		}
	}
}
