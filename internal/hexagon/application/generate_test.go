package application_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"reflect"
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

// T1's stub-pin test is removed in T5 — all four handlers are real
// now. errStubHandler and ErrStubHandlerForTest are removed along
// with this test (slice-m7-generate.md T5 DoD).

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

// --- T3: generate readme ---------------------------------------------
//
// The full state-machine is exercised by the T2 env-example tests
// against the shared generateManagedFile helper. T3 only adds the
// readme-specific shape pins:
//
//   * Absent → Created writes README.md (not .env.example), proving
//     the wrapper routes the correct relPath/template/style.
//   * User-content-after-block survives the splice, the Markdown
//     variant of the T2 add-on-preservation test.
//   * NoOp idempotency on the readme path, since the wrapper indirects
//     through the helper.

func readmePath() string {
	return filepath.Join(generateTestBaseDir, "README.md")
}

func renderedReadme(t *testing.T) []byte {
	t.Helper()
	body, err := application.RenderTemplateForTest("readme.md.tmpl", envExampleProjectName)
	if err != nil {
		t.Fatalf("render readme template: %v", err)
	}
	return body
}

func generateReadme(t *testing.T, svc *application.GenerateService) (driving.GenerateResponse, error) {
	t.Helper()
	return svc.Generate(context.Background(), driving.GenerateRequest{
		BaseDir:  generateTestBaseDir,
		Artifact: domain.ArtifactReadme,
	})
}

func TestGenerateReadme_Absent_ReturnsCreated(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)

	resp, err := generateReadme(t, svc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Action != driving.GenerateActionCreated {
		t.Errorf("Action = %v, want Created", resp.Action)
	}
	if got, want := resp.Changed, []string{"README.md"}; !equalStrings(got, want) {
		t.Errorf("Changed = %v, want %v", got, want)
	}

	got, err := fs.ReadFile(readmePath())
	if err != nil {
		t.Fatalf("read written README.md: %v", err)
	}
	if want := renderedReadme(t); !bytes.Equal(got, want) {
		t.Errorf("written body differs from rendered template:\n got=%q\nwant=%q", got, want)
	}
}

// TestGenerateReadme_UserContentAfterBlock_Preserved pins the
// Markdown variant of the T2 add-on-preservation test: a user-curated
// `## Custom section` after the `<!-- END U-BOOT MANAGED BLOCK: init -->`
// marker must survive an UpdatedBlock splice byte-identically. This is
// the realistic shape of a user-maintained README (init block is the
// scaffold, custom sections live below).
func TestGenerateReadme_UserContentAfterBlock_Preserved(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)
	userTail := "\n## Custom section\n\nUser-curated Markdown that must survive the splice.\n"
	// Stale init block + the user-curated tail.
	staleBlock := "<!-- BEGIN U-BOOT MANAGED BLOCK: init -->\n# stale heading — superseded by template\n<!-- END U-BOOT MANAGED BLOCK: init -->\n"
	if err := fs.WriteFile(readmePath(), []byte(staleBlock+userTail), 0o644); err != nil {
		t.Fatalf("seed README.md: %v", err)
	}

	resp, err := generateReadme(t, svc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Action != driving.GenerateActionUpdatedBlock {
		t.Errorf("Action = %v, want UpdatedBlock", resp.Action)
	}

	got, err := fs.ReadFile(readmePath())
	if err != nil {
		t.Fatalf("read updated README.md: %v", err)
	}
	if !strings.HasSuffix(string(got), userTail) {
		t.Errorf("user tail not preserved after splice; got tail = %q, want suffix = %q",
			lastN(got, len(userTail)+50), userTail)
	}
}

func TestGenerateReadme_DoubleRun_SecondCallNoOp(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)

	if _, err := generateReadme(t, svc); err != nil {
		t.Fatalf("first run: %v", err)
	}
	writesAfterFirst := len(fs.writtenPaths())

	resp, err := generateReadme(t, svc)
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

// lastN returns the last n bytes of data (or all if len(data) <= n);
// used to truncate diff-noise in test error messages.
func lastN(data []byte, n int) []byte {
	if len(data) <= n {
		return data
	}
	return data[len(data)-n:]
}

// --- T5: generate devcontainer ---------------------------------------

const (
	devcontainerJSONRel = ".devcontainer/devcontainer.json"
	dockerfileRel       = ".devcontainer/Dockerfile"
)

func devcontainerJSONPath() string { return filepath.Join(generateTestBaseDir, devcontainerJSONRel) }
func dockerfilePath() string       { return filepath.Join(generateTestBaseDir, dockerfileRel) }

func seedGenerateComposeYAML(t *testing.T, fs *fakeFS, body string) {
	t.Helper()
	if err := fs.WriteFile(filepath.Join(generateTestBaseDir, "compose.yaml"),
		[]byte(body), 0o644); err != nil {
		t.Fatalf("seed compose.yaml: %v", err)
	}
}

func seedUBootYAMLPostgres(t *testing.T, fs *fakeFS) {
	t.Helper()
	body := "schemaVersion: 1\nproject:\n  name: t-uboot-gen\nservices:\n  postgres:\n    enabled: true\n"
	if err := fs.WriteFile(filepath.Join(generateTestBaseDir, "u-boot.yaml"),
		[]byte(body), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
}

func seedUBootYAMLPostgresAndDummy(t *testing.T, fs *fakeFS) {
	t.Helper()
	body := "schemaVersion: 1\nproject:\n  name: t-uboot-gen\nservices:\n  postgres:\n    enabled: true\n  dummy:\n    enabled: true\n"
	if err := fs.WriteFile(filepath.Join(generateTestBaseDir, "u-boot.yaml"),
		[]byte(body), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
}

const composeYAMLPostgres = `services:
  postgres:
    image: postgres:16-alpine
    ports:
      - "5432:5432"
`

const composeYAMLPostgresAndDummy = `services:
  postgres:
    image: postgres:16-alpine
    ports:
      - "5432:5432"
  dummy:
    image: nginx:latest
    ports:
      - "8080:80"
`

func generateDevcontainer(t *testing.T, svc *application.GenerateService) (driving.GenerateResponse, error) {
	t.Helper()
	return svc.Generate(context.Background(), driving.GenerateRequest{
		BaseDir:  generateTestBaseDir,
		Artifact: domain.ArtifactDevcontainer,
	})
}

// parseDevcontainerForwardPorts extracts the forwardPorts integer
// list from a rendered devcontainer.json (after stripping JSONC
// comments). Returns nil when the key is absent — matching the
// LH-FA-DEV-005 "darf fehlen" contract.
func parseDevcontainerForwardPorts(t *testing.T, body []byte) []int {
	t.Helper()
	stripped := application.StripJSONCForTest(body)
	if !json.Valid(stripped) {
		t.Fatalf("rendered devcontainer.json is not valid JSON after stripJSONC:\n%s",
			stripped)
	}
	var shape struct {
		ForwardPorts []int `json:"forwardPorts"`
	}
	if err := json.Unmarshal(stripped, &shape); err != nil {
		t.Fatalf("unmarshal devcontainer.json: %v", err)
	}
	return shape.ForwardPorts
}

func TestGenerateDevcontainer_Absent_BothFilesCreated_MinimumFields(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLPostgres(t, fs)
	seedGenerateComposeYAML(t, fs, composeYAMLPostgres)

	resp, err := generateDevcontainer(t, svc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Action != driving.GenerateActionCreated {
		t.Errorf("Action = %v, want Created", resp.Action)
	}
	wantChanged := []string{devcontainerJSONRel, dockerfileRel}
	if !equalStrings(resp.Changed, wantChanged) {
		t.Errorf("Changed = %v, want %v", resp.Changed, wantChanged)
	}

	// devcontainer.json: valid JSONC, name + build + remoteUser + forwardPorts.
	jsonBody, err := fs.ReadFile(devcontainerJSONPath())
	if err != nil {
		t.Fatalf("read devcontainer.json: %v", err)
	}
	stripped := application.StripJSONCForTest(jsonBody)
	if !json.Valid(stripped) {
		t.Fatalf("devcontainer.json not valid after stripJSONC:\n%s", stripped)
	}
	var dc struct {
		Name         string `json:"name"`
		Build        struct {
			Dockerfile string `json:"dockerfile"`
			Context    string `json:"context"`
		} `json:"build"`
		ForwardPorts []int  `json:"forwardPorts"`
		RemoteUser   string `json:"remoteUser"`
	}
	if err := json.Unmarshal(stripped, &dc); err != nil {
		t.Fatalf("unmarshal devcontainer.json: %v", err)
	}
	if dc.Name != envExampleProjectName {
		t.Errorf("name = %q, want %q", dc.Name, envExampleProjectName)
	}
	if dc.Build.Dockerfile != "./Dockerfile" || dc.Build.Context != "." {
		t.Errorf("build = %+v, want dockerfile=./Dockerfile, context=.", dc.Build)
	}
	if dc.RemoteUser != "vscode" {
		t.Errorf("remoteUser = %q, want %q (LH-FA-DEV-004)", dc.RemoteUser, "vscode")
	}
	if !reflect.DeepEqual(dc.ForwardPorts, []int{5432}) {
		t.Errorf("forwardPorts = %v, want [5432]", dc.ForwardPorts)
	}

	// Dockerfile: present, USER vscode line.
	dockerBody, err := fs.ReadFile(dockerfilePath())
	if err != nil {
		t.Fatalf("read Dockerfile: %v", err)
	}
	if !bytes.Contains(dockerBody, []byte("USER vscode")) {
		t.Errorf("Dockerfile missing `USER vscode` line (LH-FA-DEV-004):\n%s", dockerBody)
	}
}

func TestGenerateDevcontainer_NoActiveServices_ForwardPortsAbsent(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs) // u-boot.yaml without services

	if _, err := generateDevcontainer(t, svc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	body, err := fs.ReadFile(devcontainerJSONPath())
	if err != nil {
		t.Fatalf("read devcontainer.json: %v", err)
	}
	stripped := application.StripJSONCForTest(body)
	// LH-FA-DEV-005: forwardPorts "darf fehlen" — assert the key is
	// not even present in the rendered JSON (not just an empty array).
	if bytes.Contains(stripped, []byte("forwardPorts")) {
		t.Errorf("forwardPorts key emitted when no active services; body:\n%s", stripped)
	}
}

func TestGenerateDevcontainer_MultipleServices_SortedContainerPorts(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLPostgresAndDummy(t, fs)
	seedGenerateComposeYAML(t, fs, composeYAMLPostgresAndDummy)

	if _, err := generateDevcontainer(t, svc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	body, err := fs.ReadFile(devcontainerJSONPath())
	if err != nil {
		t.Fatalf("read devcontainer.json: %v", err)
	}
	got := parseDevcontainerForwardPorts(t, body)
	// Container-side ports, sorted, deduplicated.
	want := []int{80, 5432}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("forwardPorts = %v, want %v (container-side, sorted)", got, want)
	}
}

// TestGenerateDevcontainer_AntiDriftAgainstDoctorPortHelper pins the
// invariant that the generator and the doctor
// `devcontainer.forwardPorts.consistency` check share a single
// source of truth (slice plan §"Wiederverwendung der Compose-Port-
// Detektion"): both must read from collectActiveServicePorts. A
// refactor that points the generator at a different port source
// (e.g. a u-boot.yaml-based tree) would silently flip this test.
func TestGenerateDevcontainer_AntiDriftAgainstDoctorPortHelper(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLPostgresAndDummy(t, fs)
	seedGenerateComposeYAML(t, fs, composeYAMLPostgresAndDummy)

	// Expected: what the doctor helper would compute on the same
	// fixture. Service names must be sorted ascending to match the
	// activeServiceNames() contract.
	doctorPorts, err := application.CollectActiveServicePortsForTest(
		fs, &fakeYAML{}, generateTestBaseDir, []string{"dummy", "postgres"})
	if err != nil {
		t.Fatalf("doctor port helper: %v", err)
	}
	// Fixture-sanity guard (review-followup N3): the fixture must
	// produce a non-empty port list, otherwise a future regression
	// that breaks port detection in both call sites — generator
	// returning nil AND doctor returning nil — would silently
	// DeepEqual-true and the anti-drift pin would lose its bite.
	if len(doctorPorts) == 0 {
		t.Fatalf("fixture sanity: expected non-empty doctor port list, got nil")
	}

	if _, err := generateDevcontainer(t, svc); err != nil {
		t.Fatalf("generate devcontainer: %v", err)
	}
	body, err := fs.ReadFile(devcontainerJSONPath())
	if err != nil {
		t.Fatalf("read devcontainer.json: %v", err)
	}
	generatorPorts := parseDevcontainerForwardPorts(t, body)

	if !reflect.DeepEqual(generatorPorts, doctorPorts) {
		t.Errorf("generator/doctor forwardPorts drift:\n  generator = %v\n  doctor    = %v",
			generatorPorts, doctorPorts)
	}
}

func TestGenerateDevcontainer_DoubleRun_SecondCallNoOp(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLPostgres(t, fs)
	seedGenerateComposeYAML(t, fs, composeYAMLPostgres)

	if _, err := generateDevcontainer(t, svc); err != nil {
		t.Fatalf("first run: %v", err)
	}
	writesAfterFirst := len(fs.writtenPaths())

	resp, err := generateDevcontainer(t, svc)
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

// TestGenerateDevcontainer_AtomicConflict_NoPartialWrite pins the
// atomic plan-and-execute contract: if one file is present with a
// missing init block, the call must fail with
// ErrGenerateManualConflict and write **nothing** — including not
// touching the absent counterpart. Otherwise the next call would
// see a half-written state and re-classify it as a conflict on a
// freshly-written file.
func TestGenerateDevcontainer_AtomicConflict_NoPartialWrite(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLPostgres(t, fs)
	seedGenerateComposeYAML(t, fs, composeYAMLPostgres)
	// Seed devcontainer.json with content but no init block.
	if err := fs.WriteFile(devcontainerJSONPath(),
		[]byte(`{"name": "user-managed"}`), 0o644); err != nil {
		t.Fatalf("seed devcontainer.json: %v", err)
	}

	writesBefore := len(fs.writtenPaths())
	_, err := generateDevcontainer(t, svc)
	if !errors.Is(err, driving.ErrGenerateManualConflict) {
		t.Fatalf("err = %v, want wrap of ErrGenerateManualConflict", err)
	}
	if delta := len(fs.writtenPaths()) - writesBefore; delta != 0 {
		t.Errorf("conflict produced %d WriteFile call(s), want 0; writes = %v",
			delta, fs.writtenPaths())
	}
	// The Dockerfile must not have been created.
	exists, err := fs.Exists(dockerfilePath())
	if err != nil {
		t.Fatalf("Exists(Dockerfile): %v", err)
	}
	if exists {
		t.Errorf("Dockerfile was created despite atomic conflict on devcontainer.json")
	}
}

// TestGenerateDevcontainer_MalformedBlock_AtomicNoWrite is the
// sibling of the conflict test for the malformed-block branch.
func TestGenerateDevcontainer_MalformedBlock_AtomicNoWrite(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLPostgres(t, fs)
	seedGenerateComposeYAML(t, fs, composeYAMLPostgres)
	// Seed Dockerfile with BEGIN-without-END.
	if err := fs.WriteFile(dockerfilePath(),
		[]byte("# BEGIN U-BOOT MANAGED BLOCK: init\nFROM scratch\n"), 0o644); err != nil {
		t.Fatalf("seed Dockerfile: %v", err)
	}

	writesBefore := len(fs.writtenPaths())
	_, err := generateDevcontainer(t, svc)
	if !errors.Is(err, driving.ErrGenerateManualConflict) {
		t.Fatalf("err = %v, want wrap of ErrGenerateManualConflict", err)
	}
	if msg := err.Error(); !strings.Contains(strings.ToLower(msg), "malformed") {
		t.Errorf("error %q lacks 'malformed' detail", msg)
	}
	if delta := len(fs.writtenPaths()) - writesBefore; delta != 0 {
		t.Errorf("malformed-block path produced %d WriteFile call(s), want 0", delta)
	}
}

// TestGenerateDevcontainer_PartialClean_FreshDockerfileAbsent pins
// the slice plan §"partial-clean" row: devcontainer.json present
// with a fresh block, Dockerfile absent. Aggregate action is
// Created (no UpdatedBlock because the present file is already
// fresh; the absent file is newly written).
func TestGenerateDevcontainer_PartialClean_FreshDockerfileAbsent_ReturnsCreated(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedUBootYAMLPostgres(t, fs)
	seedGenerateComposeYAML(t, fs, composeYAMLPostgres)

	// First run creates both files. We then delete the Dockerfile
	// to simulate the partial-clean state.
	if _, err := generateDevcontainer(t, svc); err != nil {
		t.Fatalf("first run: %v", err)
	}
	delete(fs.files, dockerfilePath())

	resp, err := generateDevcontainer(t, svc)
	if err != nil {
		t.Fatalf("second run: %v", err)
	}
	if resp.Action != driving.GenerateActionCreated {
		t.Errorf("Action = %v, want Created (only Dockerfile was newly written; devcontainer.json was fresh)",
			resp.Action)
	}
	if !equalStrings(resp.Changed, []string{dockerfileRel}) {
		t.Errorf("Changed = %v, want [%q]", resp.Changed, dockerfileRel)
	}
}

// --- T4: generate changelog ------------------------------------------

func changelogPath() string {
	return filepath.Join(generateTestBaseDir, "CHANGELOG.md")
}

func renderedChangelog(t *testing.T) []byte {
	t.Helper()
	body, err := application.RenderTemplateForTest("changelog.md.tmpl", envExampleProjectName)
	if err != nil {
		t.Fatalf("render changelog template: %v", err)
	}
	return body
}

func seedChangelog(t *testing.T, fs *fakeFS, body []byte) {
	t.Helper()
	if err := fs.WriteFile(changelogPath(), body, 0o644); err != nil {
		t.Fatalf("seed CHANGELOG.md: %v", err)
	}
}

func generateChangelog(t *testing.T, svc *application.GenerateService) (driving.GenerateResponse, error) {
	t.Helper()
	return svc.Generate(context.Background(), driving.GenerateRequest{
		BaseDir:  generateTestBaseDir,
		Artifact: domain.ArtifactChangelog,
	})
}

func TestGenerateChangelog_Absent_ReturnsCreated(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)

	resp, err := generateChangelog(t, svc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Action != driving.GenerateActionCreated {
		t.Errorf("Action = %v, want Created", resp.Action)
	}
	if got, want := resp.Changed, []string{"CHANGELOG.md"}; !equalStrings(got, want) {
		t.Errorf("Changed = %v, want %v", got, want)
	}

	got, err := fs.ReadFile(changelogPath())
	if err != nil {
		t.Fatalf("read written CHANGELOG.md: %v", err)
	}
	if want := renderedChangelog(t); !bytes.Equal(got, want) {
		t.Errorf("written body differs from rendered template:\n got=%q\nwant=%q", got, want)
	}
}

func TestGenerateChangelog_DoubleRun_SecondCallNoOp(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)

	if _, err := generateChangelog(t, svc); err != nil {
		t.Fatalf("first run: %v", err)
	}
	writesAfterFirst := len(fs.writtenPaths())

	resp, err := generateChangelog(t, svc)
	if err != nil {
		t.Fatalf("second run: %v", err)
	}
	if resp.Action != driving.GenerateActionNoOp {
		t.Errorf("second run Action = %v, want NoOp", resp.Action)
	}
	if delta := len(fs.writtenPaths()) - writesAfterFirst; delta != 0 {
		t.Errorf("second run produced %d WriteFile call(s), want 0; writes = %v",
			delta, fs.writtenPaths())
	}
}

func TestGenerateChangelog_PresentNoBlock_ReturnsErrManualConflict(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)
	seedChangelog(t, fs, []byte("# Custom changelog\n\nNo init block here.\n"))

	writesBefore := len(fs.writtenPaths())
	_, err := generateChangelog(t, svc)
	if !errors.Is(err, driving.ErrGenerateManualConflict) {
		t.Fatalf("err = %v, want wrap of ErrGenerateManualConflict", err)
	}
	if delta := len(fs.writtenPaths()) - writesBefore; delta != 0 {
		t.Errorf("manual-conflict path produced %d WriteFile call(s), want 0", delta)
	}
}

func TestGenerateChangelog_PresentMalformedBlock_ReturnsErrManualConflict(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)
	seedChangelog(t, fs, []byte("<!-- BEGIN U-BOOT MANAGED BLOCK: init -->\n# orphan body\n"))

	writesBefore := len(fs.writtenPaths())
	_, err := generateChangelog(t, svc)
	if !errors.Is(err, driving.ErrGenerateManualConflict) {
		t.Fatalf("err = %v, want wrap of ErrGenerateManualConflict", err)
	}
	if msg := err.Error(); !strings.Contains(strings.ToLower(msg), "malformed") {
		t.Errorf("error message %q lacks 'malformed' detail", msg)
	}
	if delta := len(fs.writtenPaths()) - writesBefore; delta != 0 {
		t.Errorf("malformed-block path produced %d WriteFile call(s), want 0", delta)
	}
}

// TestGenerateChangelog_UserEditedBlock_WithUnreleased_NoOp pins the
// idempotency contract for a user-curated file where the init block
// has been touched (so bytes.Equal(existing, rendered) is false) but
// a `## [Unreleased]` section is still present. The handler must not
// re-render the block (LH-AK-007 "vorhandene Inhalte werden nicht
// zerstört") and must not insert a redundant Unreleased header.
func TestGenerateChangelog_UserEditedBlock_WithUnreleased_NoOp(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)
	// Stale init block (different from the freshly-rendered template),
	// but the user added their own Unreleased entries below — so the
	// `## [Unreleased]` header is preserved.
	seed := []byte("<!-- BEGIN U-BOOT MANAGED BLOCK: init -->\n# Changelog\n\nUser-curated intro.\n\n## [Unreleased]\n\n### Added\n\n- user entry\n<!-- END U-BOOT MANAGED BLOCK: init -->\n")
	seedChangelog(t, fs, seed)

	writesBefore := len(fs.writtenPaths())
	resp, err := generateChangelog(t, svc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Action != driving.GenerateActionNoOp {
		t.Errorf("Action = %v, want NoOp", resp.Action)
	}
	if len(resp.Changed) != 0 {
		t.Errorf("Changed = %v, want empty", resp.Changed)
	}
	if delta := len(fs.writtenPaths()) - writesBefore; delta != 0 {
		t.Errorf("NoOp path produced %d WriteFile call(s), want 0", delta)
	}
}

// TestGenerateChangelog_UserEditedBlock_MissingUnreleased_RepairedManual
// pins the RepairedManual path from the slice plan §T4-table: the
// user has cut a release (moved Unreleased entries to `## [0.1.0]`)
// and forgot to seed a fresh Unreleased. The handler inserts an
// Unreleased stub before the first release section, outside the
// init block.
func TestGenerateChangelog_UserEditedBlock_MissingUnreleased_RepairedManual(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)
	// User-edited block (no Unreleased inside) + a release section
	// after the END marker. This matches the slice plan's
	// RepairedManual fixture description.
	seed := []byte(
		"<!-- BEGIN U-BOOT MANAGED BLOCK: init -->\n" +
			"# Changelog\n\nUser-curated intro.\n" +
			"<!-- END U-BOOT MANAGED BLOCK: init -->\n\n" +
			"## [0.1.0] - 2026-01-01\n\n" +
			"### Added\n\n- initial release\n",
	)
	seedChangelog(t, fs, seed)

	resp, err := generateChangelog(t, svc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Action != driving.GenerateActionRepairedManual {
		t.Errorf("Action = %v, want RepairedManual", resp.Action)
	}
	if got, want := resp.Changed, []string{"CHANGELOG.md"}; !equalStrings(got, want) {
		t.Errorf("Changed = %v, want %v", got, want)
	}

	got, err := fs.ReadFile(changelogPath())
	if err != nil {
		t.Fatalf("read repaired CHANGELOG.md: %v", err)
	}
	// The Unreleased stub must appear before the release section and
	// after the END marker. Find both anchors and assert the order.
	unreleasedIdx := bytes.Index(got, []byte("## [Unreleased]"))
	releaseIdx := bytes.Index(got, []byte("## [0.1.0]"))
	endMarkerIdx := bytes.Index(got, []byte("<!-- END U-BOOT MANAGED BLOCK: init -->"))
	if unreleasedIdx < 0 {
		t.Fatal("repaired body missing ## [Unreleased] header")
	}
	if releaseIdx < 0 {
		t.Fatal("repaired body lost ## [0.1.0] header")
	}
	if endMarkerIdx < 0 || unreleasedIdx < endMarkerIdx {
		t.Errorf("Unreleased stub was inserted inside or before the init block (END idx=%d, Unreleased idx=%d); want outside the block",
			endMarkerIdx, unreleasedIdx)
	}
	if unreleasedIdx >= releaseIdx {
		t.Errorf("Unreleased stub at offset %d is not before release section at offset %d", unreleasedIdx, releaseIdx)
	}
	// Idempotency: a second invocation now sees Unreleased present
	// and must NoOp.
	writesAfterRepair := len(fs.writtenPaths())
	resp2, err := generateChangelog(t, svc)
	if err != nil {
		t.Fatalf("second run after repair: %v", err)
	}
	if resp2.Action != driving.GenerateActionNoOp {
		t.Errorf("second run Action = %v, want NoOp", resp2.Action)
	}
	if delta := len(fs.writtenPaths()) - writesAfterRepair; delta != 0 {
		t.Errorf("second run produced %d WriteFile call(s), want 0", delta)
	}
}

// TestGenerateChangelog_UserEditedBlock_FencedReleaseOnly_NoOp pins
// the review-followup S1 fix: a `## [1.2.3]` header that only
// appears inside a ```` ```md ```` fenced code block must NOT count
// as a real release section. Without the fix the handler would
// splice the Unreleased stub *inside* the fence and corrupt the
// user's Markdown.
func TestGenerateChangelog_UserEditedBlock_FencedReleaseOnly_NoOp(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)
	seed := []byte(
		"<!-- BEGIN U-BOOT MANAGED BLOCK: init -->\n" +
			"# Changelog\n\nUser-curated intro.\n" +
			"<!-- END U-BOOT MANAGED BLOCK: init -->\n\n" +
			"Documentation example:\n\n" +
			"```md\n## [1.2.3] - 2026-01-01\n\n### Added\n\n- example\n```\n",
	)
	seedChangelog(t, fs, seed)

	writesBefore := len(fs.writtenPaths())
	resp, err := generateChangelog(t, svc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Action != driving.GenerateActionNoOp {
		t.Errorf("Action = %v, want NoOp (the fenced [1.2.3] must not anchor a splice)", resp.Action)
	}
	if delta := len(fs.writtenPaths()) - writesBefore; delta != 0 {
		t.Errorf("fenced-only path produced %d WriteFile call(s), want 0", delta)
	}
}

// TestGenerateChangelog_UserEditedBlock_FencedReleaseBeforeReal_SpliceAtReal
// pins the harder S1 case: a fenced `## [9.9.9]` example precedes a
// real `## [0.1.0]` section. The Unreleased stub must be spliced
// before the *real* one, not the fenced one — otherwise the fence
// is broken.
func TestGenerateChangelog_UserEditedBlock_FencedReleaseBeforeReal_SpliceAtReal(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)
	seed := []byte(
		"<!-- BEGIN U-BOOT MANAGED BLOCK: init -->\n" +
			"# Changelog\n\nUser-curated intro.\n" +
			"<!-- END U-BOOT MANAGED BLOCK: init -->\n\n" +
			"Example format:\n\n" +
			"```md\n## [9.9.9]\n```\n\n" +
			"## [0.1.0] - 2026-01-01\n\n### Added\n\n- initial release\n",
	)
	seedChangelog(t, fs, seed)

	resp, err := generateChangelog(t, svc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Action != driving.GenerateActionRepairedManual {
		t.Errorf("Action = %v, want RepairedManual", resp.Action)
	}

	got, err := fs.ReadFile(changelogPath())
	if err != nil {
		t.Fatalf("read repaired CHANGELOG.md: %v", err)
	}
	// The fenced [9.9.9] must remain intact (still inside the fence).
	if !bytes.Contains(got, []byte("```md\n## [9.9.9]\n```")) {
		t.Errorf("fenced [9.9.9] example was corrupted; got:\n%s", got)
	}
	// The Unreleased stub must sit between the fence and [0.1.0].
	fenceCloseIdx := bytes.Index(got, []byte("```md\n## [9.9.9]\n```"))
	unreleasedIdx := bytes.Index(got, []byte("## [Unreleased]"))
	realReleaseIdx := bytes.Index(got, []byte("## [0.1.0]"))
	if fenceCloseIdx >= unreleasedIdx || unreleasedIdx >= realReleaseIdx {
		t.Errorf("Unreleased stub at offset %d not between fence (%d) and real release (%d)",
			unreleasedIdx, fenceCloseIdx, realReleaseIdx)
	}
}

// TestGenerateChangelog_UserEditedBlock_FencedUnreleased_DoesNotCount
// pins the symmetric S1 case for [hasChangelogUnreleased]: a fenced
// `## [Unreleased]` example must NOT make the handler believe a real
// Unreleased section exists. With a real release section also
// present, the handler should hit the RepairedManual path.
func TestGenerateChangelog_UserEditedBlock_FencedUnreleased_DoesNotCount(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)
	seed := []byte(
		"<!-- BEGIN U-BOOT MANAGED BLOCK: init -->\n" +
			"# Changelog\n\nUser-curated intro.\n" +
			"<!-- END U-BOOT MANAGED BLOCK: init -->\n\n" +
			"Convention reference:\n\n" +
			"```md\n## [Unreleased]\n```\n\n" +
			"## [0.1.0] - 2026-01-01\n",
	)
	seedChangelog(t, fs, seed)

	resp, err := generateChangelog(t, svc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Action != driving.GenerateActionRepairedManual {
		t.Errorf("Action = %v, want RepairedManual (the fenced Unreleased must not count)", resp.Action)
	}
}

// TestGenerateChangelog_CRLFFreshBlock_NoOp pins the review-followup
// S2 fix: a CHANGELOG.md saved with CRLF line endings that otherwise
// matches the rendered template byte-for-byte (LF) must still
// register as fresh — the LF-normalised bytes.Equal must accept the
// match instead of routing to the user-edited branch.
func TestGenerateChangelog_CRLFFreshBlock_NoOp(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)
	rendered := renderedChangelog(t)
	crlf := bytes.ReplaceAll(rendered, []byte{'\n'}, []byte{'\r', '\n'})
	seedChangelog(t, fs, crlf)

	writesBefore := len(fs.writtenPaths())
	resp, err := generateChangelog(t, svc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Action != driving.GenerateActionNoOp {
		t.Errorf("Action = %v, want NoOp (CRLF-edited fresh block must register as fresh)", resp.Action)
	}
	if delta := len(fs.writtenPaths()) - writesBefore; delta != 0 {
		t.Errorf("CRLF NoOp path produced %d WriteFile call(s), want 0", delta)
	}
}

// TestGenerateEnvExample_CRLFFreshBlock_NoOp covers the same S2 fix
// applied to the generateManagedFile path (T2/T3 helper) — a CRLF-
// edited fresh env-example registers as NoOp instead of silently
// rewriting CRLF→LF.
func TestGenerateEnvExample_CRLFFreshBlock_NoOp(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)
	rendered := renderedEnvExample(t)
	crlf := bytes.ReplaceAll(rendered, []byte{'\n'}, []byte{'\r', '\n'})
	seedEnvExample(t, fs, crlf)

	writesBefore := len(fs.writtenPaths())
	resp, err := generateEnvExample(t, svc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Action != driving.GenerateActionNoOp {
		t.Errorf("Action = %v, want NoOp (CRLF-edited fresh block must register as fresh)", resp.Action)
	}
	if delta := len(fs.writtenPaths()) - writesBefore; delta != 0 {
		t.Errorf("CRLF NoOp path produced %d WriteFile call(s), want 0", delta)
	}
}

// TestGenerateChangelog_UserEditedBlock_NoUnreleased_NoVersion_NoOp
// covers the conservative branch where the user has neither an
// Unreleased section nor a recognisable release section — the
// handler refuses to invent structure and leaves the file alone.
func TestGenerateChangelog_UserEditedBlock_NoUnreleased_NoVersion_NoOp(t *testing.T) {
	t.Parallel()
	svc, fs := newGenerateService(t)
	seedGenerateUbootYAML(t, fs)
	seed := []byte(
		"<!-- BEGIN U-BOOT MANAGED BLOCK: init -->\n" +
			"# Changelog\n\nFree-form notes.\n" +
			"<!-- END U-BOOT MANAGED BLOCK: init -->\n\n" +
			"Just a paragraph, no Keep-a-Changelog headers.\n",
	)
	seedChangelog(t, fs, seed)

	writesBefore := len(fs.writtenPaths())
	resp, err := generateChangelog(t, svc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Action != driving.GenerateActionNoOp {
		t.Errorf("Action = %v, want NoOp", resp.Action)
	}
	if delta := len(fs.writtenPaths()) - writesBefore; delta != 0 {
		t.Errorf("conservative NoOp produced %d WriteFile call(s), want 0", delta)
	}
}

// TestGenerateChangelog_LHAK007_FlowEndToEnd pins LH-AK-007 verbatim:
// `u-boot init && u-boot generate changelog` produces a CHANGELOG.md
// that exists, has not destroyed any pre-existing content (vacuously
// true for the fresh-init path), and contains the expected
// `## [Unreleased]` section. The test exercises both application
// services side-by-side rather than going through the CLI adapter.
func TestGenerateChangelog_LHAK007_FlowEndToEnd(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	fs := newFakeFS()
	fs.markDirExists(dir)
	y := &fakeYAML{}
	git := &fakeGit{}
	prog := &fakeProgress{}
	conf := &fakeConfirmer{}

	initSvc := application.NewInitProjectService(fs, y, git, prog, conf, nil)
	genSvc := application.NewGenerateService(fs, y, nil)

	if _, err := initSvc.Init(context.Background(), driving.InitProjectRequest{
		BaseDir: dir,
		Name:    "ak007demo",
		SkipGit: true,
	}); err != nil {
		t.Fatalf("init: %v", err)
	}

	resp, err := genSvc.Generate(context.Background(), driving.GenerateRequest{
		BaseDir:  dir,
		Artifact: domain.ArtifactChangelog,
	})
	if err != nil {
		t.Fatalf("generate changelog: %v", err)
	}
	// Init already writes CHANGELOG.md, so generate sees a present
	// fresh file and returns NoOp. That still satisfies LH-AK-007 —
	// the file exists, pre-existing content is preserved, the
	// Unreleased section is present.
	if resp.Action != driving.GenerateActionNoOp {
		t.Errorf("Action = %v, want NoOp (init already wrote a fresh CHANGELOG.md)", resp.Action)
	}

	body, err := fs.ReadFile(filepath.Join(dir, "CHANGELOG.md"))
	if err != nil {
		t.Fatalf("read CHANGELOG.md: %v", err)
	}
	if !bytes.Contains(body, []byte("## [Unreleased]")) {
		t.Errorf("CHANGELOG.md missing ## [Unreleased] section; got:\n%s", body)
	}
	if !bytes.Contains(body, []byte("ak007demo")) {
		t.Errorf("CHANGELOG.md does not reference project name; got:\n%s", body)
	}
}
