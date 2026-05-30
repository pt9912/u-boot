package application_test

import (
	"bytes"
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/application/managedblock"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

const generateTestBaseDir = "/tmp/u-boot-generate-test/demo"

// newGenerateService constructs the service against the shared
// in-memory FS plus the yaml.v3-backed fake codec. The BaseDir is
// pre-registered as an existing directory so the Exists() check on
// u-boot.yaml fails on file-absence, not on directory-absence.
func newGenerateService(t *testing.T) (*application.GenerateService, *fakeFS) {
	t.Helper()
	fs := newFakeFS()
	fs.markDirExists(generateTestBaseDir)
	y := &fakeYAML{}
	svc := application.NewGenerateService(fs, y, nil)
	return svc, fs
}

// seedGenerateUbootYAML writes a minimal u-boot.yaml under the test
// BaseDir so the project-initialized gate passes.
func seedGenerateUbootYAML(t *testing.T, fs *fakeFS) {
	t.Helper()
	body := "schemaVersion: 1\nproject:\n  name: t-uboot-gen\n"
	if err := fs.WriteFile(filepath.Join(generateTestBaseDir, "u-boot.yaml"),
		[]byte(body), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
}

func TestGenerate_BaseDirEmpty_NonSentinelError(t *testing.T) {
	t.Parallel()
	svc, _ := newGenerateService(t)
	_, err := svc.Generate(context.Background(), driving.GenerateRequest{
		BaseDir:  "",
		Artifact: domain.ArtifactChangelog,
	})
	if err == nil {
		t.Fatal("expected non-nil error for empty BaseDir")
	}
	// BaseDir-empty is a programmer-/CLI-bug, not a fachliche
	// condition — must NOT map to ErrProjectNotInitialized.
	if errors.Is(err, driving.ErrProjectNotInitialized) {
		t.Errorf("empty BaseDir leaked ErrProjectNotInitialized: %v", err)
	}
}

func TestGenerate_NoUbootYAML_ReturnsErrProjectNotInitialized(t *testing.T) {
	t.Parallel()
	svc, _ := newGenerateService(t)
	for _, art := range []domain.Artifact{
		domain.ArtifactChangelog,
		domain.ArtifactReadme,
		domain.ArtifactEnvExample,
		domain.ArtifactDevcontainer,
	} {
		_, err := svc.Generate(context.Background(), driving.GenerateRequest{
			BaseDir:  generateTestBaseDir,
			Artifact: art,
		})
		if !errors.Is(err, driving.ErrProjectNotInitialized) {
			t.Errorf("artifact=%s: err = %v, want wrap of ErrProjectNotInitialized", art, err)
		}
	}
}

// TestGenerate_StubHandlers pins the remaining unimplemented handlers.
// The pin reduces by one with each of T3..T5 as those tranches replace
// handlers; T5 removes the test when the last stub is gone
// (slice-m7-generate.md T5 DoD).
//
// Pin tracks remaining stubs: 4 in T1 → 3 in T2 → 2 in T3 → 1 in T4 →
// 0 in T5 (test deleted). Updated for T2: env-example is now real,
// so the catalogue covers only the three remaining artefacts.
func TestGenerate_StubHandlers_RemainingThreeReturnErrStubHandler(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)

	cases := []domain.Artifact{
		domain.ArtifactChangelog,
		domain.ArtifactReadme,
		domain.ArtifactDevcontainer,
	}
	for _, art := range cases {
		_, err := svc.Generate(context.Background(), driving.GenerateRequest{
			BaseDir:  generateTestBaseDir,
			Artifact: art,
		})
		if err == nil {
			t.Errorf("artifact=%s: expected stub-handler error, got nil", art)
			continue
		}
		if !errors.Is(err, application.ErrStubHandlerForTest) {
			t.Errorf("artifact=%s: err = %v, want wrap of errStubHandler", art, err)
		}
	}
}

// TestGenerate_UnknownArtifactValue_DefensiveBranch pins the
// out-of-enum-range guard. The CLI validates via [domain.NewArtifact]
// before constructing the request, so this branch is unreachable in
// production; the test protects against a future enum addition that
// forgets to update the dispatch switch.
func TestGenerate_UnknownArtifactValue_DefensiveBranch(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)

	_, err := svc.Generate(context.Background(), driving.GenerateRequest{
		BaseDir:  generateTestBaseDir,
		Artifact: domain.Artifact(99),
	})
	if !errors.Is(err, domain.ErrInvalidArtifact) {
		t.Errorf("err = %v, want wrap of domain.ErrInvalidArtifact", err)
	}
}

// --- T2: generate env-example ----------------------------------------

const envExampleProjectName = "t-uboot-gen"

func envExamplePath() string {
	return filepath.Join(generateTestBaseDir, ".env.example")
}

// renderedEnvExample returns the rendered env.example.tmpl body for
// the test project name. The handler under test produces this exact
// byte sequence for the absent state.
func renderedEnvExample(t *testing.T) []byte {
	t.Helper()
	body, err := application.RenderTemplateForTest("env.example.tmpl", envExampleProjectName)
	if err != nil {
		t.Fatalf("render template: %v", err)
	}
	return body
}

// renderedEnvExampleBlock returns just the BEGIN..END region of the
// rendered template. The handler splices this region into existing
// files when the present-with-block state is detected.
func renderedEnvExampleBlock(t *testing.T) []byte {
	t.Helper()
	block, err := application.RenderManagedBlockOnlyForTest(renderedEnvExample(t), managedblock.InitName)
	if err != nil {
		t.Fatalf("extract block: %v", err)
	}
	return block
}

func seedEnvExample(t *testing.T, fs *fakeFS, body []byte) {
	t.Helper()
	if err := fs.WriteFile(envExamplePath(), body, 0o644); err != nil {
		t.Fatalf("seed .env.example: %v", err)
	}
}

func generateEnvExample(t *testing.T, svc *application.GenerateService) (driving.GenerateResponse, error) {
	t.Helper()
	return svc.Generate(context.Background(), driving.GenerateRequest{
		BaseDir:  generateTestBaseDir,
		Artifact: domain.ArtifactEnvExample,
	})
}

func TestGenerateEnvExample_Absent_ReturnsCreated(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)

	resp, err := generateEnvExample(t, svc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Action != driving.GenerateActionCreated {
		t.Errorf("Action = %v, want Created", resp.Action)
	}
	if got, want := resp.Changed, []string{".env.example"}; !equalStrings(got, want) {
		t.Errorf("Changed = %v, want %v", got, want)
	}

	got, err := fs.ReadFile(envExamplePath())
	if err != nil {
		t.Fatalf("read written .env.example: %v", err)
	}
	if want := renderedEnvExample(t); !bytes.Equal(got, want) {
		t.Errorf("written body differs from rendered template:\n got=%q\nwant=%q", got, want)
	}
}

func TestGenerateEnvExample_PresentWithStaleBlock_ReturnsUpdatedBlock(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)
	// Seed a stale init block — different body from the current
	// template so the handler must rerender it.
	stale := []byte("# BEGIN U-BOOT MANAGED BLOCK: init\n# stale body — superseded by the template\n# END U-BOOT MANAGED BLOCK: init\n")
	seedEnvExample(t, fs, stale)

	resp, err := generateEnvExample(t, svc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Action != driving.GenerateActionUpdatedBlock {
		t.Errorf("Action = %v, want UpdatedBlock", resp.Action)
	}

	got, err := fs.ReadFile(envExamplePath())
	if err != nil {
		t.Fatalf("read updated .env.example: %v", err)
	}
	// The init block region should now be the rendered template's
	// init block bytes (no content outside the block in this fixture).
	if want := renderedEnvExampleBlock(t); !bytes.Equal(got, want) {
		t.Errorf("updated body differs from rendered block:\n got=%q\nwant=%q", got, want)
	}
}

// TestGenerateEnvExample_DoubleRun_SecondCallNoOp pins LH-FA-GEN-005
// idempotency on two axes: the Action field reports NoOp **and** the
// FileSystem fake records zero WriteFile calls during the second
// invocation. Asserting only the Action would let a regression slip
// past where a handler returns NoOp but still touches the file.
func TestGenerateEnvExample_DoubleRun_SecondCallNoOp(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)

	// First run: writes the artefact (Created).
	if _, err := generateEnvExample(t, svc); err != nil {
		t.Fatalf("first run: %v", err)
	}
	writesAfterFirst := len(fs.writtenPaths())

	// Second run: NoOp expected; no new WriteFile invocations.
	resp, err := generateEnvExample(t, svc)
	if err != nil {
		t.Fatalf("second run: %v", err)
	}
	if resp.Action != driving.GenerateActionNoOp {
		t.Errorf("second run Action = %v, want NoOp", resp.Action)
	}
	if len(resp.Changed) != 0 {
		t.Errorf("second run Changed = %v, want empty", resp.Changed)
	}
	if delta := len(fs.writtenPaths()) - writesAfterFirst; delta != 0 {
		t.Errorf("second run produced %d WriteFile call(s), want 0; writes = %v",
			delta, fs.writtenPaths())
	}
}

func TestGenerateEnvExample_PresentNoBlock_ReturnsErrManualConflict(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)
	// User-curated .env.example without the init managed block.
	seedEnvExample(t, fs, []byte("FOO=bar\nBAZ=qux\n"))

	writesBefore := len(fs.writtenPaths())
	_, err := generateEnvExample(t, svc)
	if !errors.Is(err, driving.ErrGenerateManualConflict) {
		t.Fatalf("err = %v, want wrap of ErrGenerateManualConflict", err)
	}
	// No write side-effect on the manual-conflict path.
	if delta := len(fs.writtenPaths()) - writesBefore; delta != 0 {
		t.Errorf("manual-conflict path produced %d WriteFile call(s), want 0; writes = %v",
			delta, fs.writtenPaths())
	}
}

func TestGenerateEnvExample_PresentMalformedBlock_ReturnsErrManualConflict(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)
	// BEGIN marker without a matching END — the malformed-block
	// branch of managedblock.Find.
	malformed := []byte("# BEGIN U-BOOT MANAGED BLOCK: init\n# orphan body — missing END marker\nFOO=bar\n")
	seedEnvExample(t, fs, malformed)

	writesBefore := len(fs.writtenPaths())
	_, err := generateEnvExample(t, svc)
	if !errors.Is(err, driving.ErrGenerateManualConflict) {
		t.Fatalf("err = %v, want wrap of ErrGenerateManualConflict", err)
	}
	// The detail message must point at the malformed-block cause so
	// the user can tell the two ErrGenerateManualConflict surfaces
	// apart (missing marker vs. malformed marker).
	if msg := err.Error(); !strings.Contains(strings.ToLower(msg), "malformed") {
		t.Errorf("error message %q lacks 'malformed' detail", msg)
	}
	// No write side-effect on the malformed-block path either.
	if delta := len(fs.writtenPaths()) - writesBefore; delta != 0 {
		t.Errorf("malformed-block path produced %d WriteFile call(s), want 0; writes = %v",
			delta, fs.writtenPaths())
	}
}

// TestGenerateEnvExample_AddOnBlockSurvives_BlockReplaceSplice pins
// the managedblock.Replace contract: content outside the init block —
// in this fixture a `service.postgres` block written by `u-boot add` —
// must remain byte-identical after the init-block re-splice.
func TestGenerateEnvExample_AddOnBlockSurvives_BlockReplaceSplice(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)
	addOnBlock := "# BEGIN U-BOOT MANAGED BLOCK: service.postgres\nPOSTGRES_USER=postgres\nPOSTGRES_PASSWORD=CHANGEME_POSTGRES_PASSWORD\nPOSTGRES_DB=postgres\n# END U-BOOT MANAGED BLOCK: service.postgres\n"
	// Stale init block + the add-on block below it.
	seed := []byte("# BEGIN U-BOOT MANAGED BLOCK: init\n# stale body\n# END U-BOOT MANAGED BLOCK: init\n" + addOnBlock)
	seedEnvExample(t, fs, seed)

	resp, err := generateEnvExample(t, svc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Action != driving.GenerateActionUpdatedBlock {
		t.Errorf("Action = %v, want UpdatedBlock", resp.Action)
	}

	got, err := fs.ReadFile(envExamplePath())
	if err != nil {
		t.Fatalf("read updated .env.example: %v", err)
	}
	// The add-on block must appear byte-identically in the new content.
	if !bytes.Contains(got, []byte(addOnBlock)) {
		t.Errorf("add-on block was lost or modified after init splice; got:\n%q", got)
	}
}

// equalStrings is a tiny helper kept local to this file to avoid
// reaching for reflect.DeepEqual / slices.Equal across Go versions.
func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
