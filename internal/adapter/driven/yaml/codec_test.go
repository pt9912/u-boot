package yaml_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driven/yaml"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

type sample struct {
	Name    string `yaml:"name"`
	Version int    `yaml:"version"`
}

func TestCodec_MarshalUnmarshalRoundTrip(t *testing.T) {
	// Why: the round-trip is the actual semantic contract — Marshal +
	// Unmarshal must preserve the value. A separate test asserts the
	// presence of expected keys without pinning yaml.v3's formatting
	// (see TestCodec_MarshalEmitsExpectedKeys), so a future yaml.v3
	// bump that changes quoting/indent does not break the round-trip
	// test for the wrong reason.
	in := sample{Name: "demo", Version: 1}

	data, err := yaml.New().Marshal(in)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var back sample
	if err := yaml.New().Unmarshal(data, &back); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if back != in {
		t.Fatalf("Round-trip: got %+v, want %+v", back, in)
	}
}

func TestCodec_MarshalEmitsExpectedKeys(t *testing.T) {
	// Why: a smoke check that the YAML output mentions the struct's
	// `yaml:"…"` keys. Tolerant of yaml.v3 formatting (quoting,
	// indent, trailing newline) — only the key presence is pinned.
	in := sample{Name: "demo", Version: 1}

	data, err := yaml.New().Marshal(in)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	got := string(data)
	for _, key := range []string{"name:", "version:"} {
		if !strings.Contains(got, key) {
			t.Errorf("Marshal output missing key %q; got:\n%s", key, got)
		}
	}
}

func TestCodec_Unmarshal_InvalidYAMLReturnsError(t *testing.T) {
	var dst sample
	err := yaml.New().Unmarshal([]byte(":\n  not yaml"), &dst)
	if err == nil {
		t.Fatalf("Unmarshal(invalid): expected error, got nil")
	}
}

func TestCodec_PatchScalar_UpdatesExistingScalar(t *testing.T) {
	// Why: the primary M5-T4 use case — flip
	// services.postgres.enabled from false to true while leaving
	// the surrounding document (including comments) intact.
	input := []byte("# leading comment\n" +
		"schemaVersion: 1\n" +
		"project:\n" +
		"  name: demo\n" +
		"services:\n" +
		"  postgres:\n" +
		"    enabled: false\n")

	out, err := yaml.New().PatchScalar(input,
		[]string{"services", "postgres", "enabled"}, true)
	if err != nil {
		t.Fatalf("PatchScalar: %v", err)
	}

	got := string(out)
	for _, want := range []string{
		"# leading comment",
		"schemaVersion: 1",
		"name: demo",
		"enabled: true",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q; got:\n%s", want, got)
		}
	}
	if strings.Contains(got, "enabled: false") {
		t.Errorf("output still contains the pre-patch value; got:\n%s", got)
	}
}

func TestCodec_PatchScalar_InsertsMissingPath(t *testing.T) {
	// Why: u-boot init writes a minimal u-boot.yaml without a
	// services: block. The first `u-boot add postgres` must be able
	// to insert the full path services.postgres.enabled: true even
	// though every level along the way is missing.
	input := []byte("schemaVersion: 1\n" +
		"project:\n" +
		"  name: demo\n")

	out, err := yaml.New().PatchScalar(input,
		[]string{"services", "postgres", "enabled"}, true)
	if err != nil {
		t.Fatalf("PatchScalar: %v", err)
	}

	got := string(out)
	for _, want := range []string{
		"schemaVersion: 1",
		"name: demo",
		"services:",
		"postgres:",
		"enabled: true",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q; got:\n%s", want, got)
		}
	}
}

func TestCodec_PatchScalar_PreservesSiblingUnknownFields(t *testing.T) {
	// Why: LH-FA-CONF-002 announces V1-fields (devcontainer,
	// template, services.keycloak.persistence) that current Go
	// structs do not know about. A naive Marshal(struct) round-trip
	// would drop them; the Node-API based PatchScalar must keep
	// them in place.
	input := []byte("schemaVersion: 1\n" +
		"project:\n" +
		"  name: demo\n" +
		"devcontainer:\n" +
		"  enabled: false\n" +
		"services:\n" +
		"  keycloak:\n" +
		"    enabled: false\n" +
		"    persistence: embedded\n")

	out, err := yaml.New().PatchScalar(input,
		[]string{"services", "postgres", "enabled"}, true)
	if err != nil {
		t.Fatalf("PatchScalar: %v", err)
	}

	got := string(out)
	for _, want := range []string{
		"devcontainer:",
		"keycloak:",
		"persistence: embedded",
		"postgres:",
		"enabled: true",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q; got:\n%s", want, got)
		}
	}
}

func TestCodec_PatchScalar_EmptyPathReturnsErr(t *testing.T) {
	_, err := yaml.New().PatchScalar([]byte("k: v\n"), nil, true)
	if !errors.Is(err, driven.ErrYAMLPathInvalid) {
		t.Fatalf("PatchScalar(empty path): want ErrYAMLPathInvalid, got %v", err)
	}
}

func TestCodec_PatchScalar_UnsupportedScalarTypeReturnsErr(t *testing.T) {
	// Why: PatchScalar's contract restricts value to scalar Go
	// kinds. Passing a slice (would otherwise become a YAML
	// sequence) must be rejected — destructive replacement of a
	// sub-tree is out of scope.
	_, err := yaml.New().PatchScalar([]byte("k: v\n"),
		[]string{"k"}, []string{"a", "b"})
	if !errors.Is(err, driven.ErrYAMLPathInvalid) {
		t.Fatalf("PatchScalar(slice value): want ErrYAMLPathInvalid, got %v", err)
	}
}

func TestCodec_PatchScalar_IntermediateScalarReturnsErr(t *testing.T) {
	// Why: if services is already a scalar (`services: foo`), the
	// patch cannot silently overwrite it with a mapping. Surface a
	// path-invalid error so the caller can decide how to handle it.
	input := []byte("services: foo\n")
	_, err := yaml.New().PatchScalar(input,
		[]string{"services", "postgres", "enabled"}, true)
	if !errors.Is(err, driven.ErrYAMLPathInvalid) {
		t.Fatalf("PatchScalar(scalar intermediate): want ErrYAMLPathInvalid, got %v", err)
	}
}

func TestCodec_PatchScalar_EmptyDocumentCreatesPath(t *testing.T) {
	// Why: u-boot.yaml is never empty in practice (init writes
	// schemaVersion + project), but the adapter should gracefully
	// handle an empty / whitespace-only input rather than panic on
	// the missing document root.
	out, err := yaml.New().PatchScalar([]byte(""),
		[]string{"services", "postgres", "enabled"}, true)
	if err != nil {
		t.Fatalf("PatchScalar(empty): %v", err)
	}
	got := string(out)
	for _, want := range []string{"services:", "postgres:", "enabled: true"} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q; got:\n%s", want, got)
		}
	}
}

func TestCodec_PatchScalar_NonMappingRootReturnsErr(t *testing.T) {
	// Why: a YAML document whose root is a sequence (e.g. starts
	// with `- item`) cannot host mapping keys. Reject explicitly
	// so a caller bug surfaces with a path-invalid sentinel rather
	// than a confusing in-place mutation.
	_, err := yaml.New().PatchScalar([]byte("- a\n- b\n"),
		[]string{"services"}, true)
	if !errors.Is(err, driven.ErrYAMLPathInvalid) {
		t.Fatalf("PatchScalar(sequence root): want ErrYAMLPathInvalid, got %v", err)
	}
}

func TestCodec_PatchScalar_PreservesValueLineComment(t *testing.T) {
	// Why: yaml.v3 attaches inline trailing comments
	// (`enabled: false  # CHANGEME …`) to the VALUE node, not the
	// key node. A naive `*valueNode = *scalar` copy silently drops
	// them — review finding #1 (M5-T3). This test pins the
	// preservation contract for the dominant LH-FA-CONF-002 idiom
	// where every CHANGEME-style scalar carries a guidance comment.
	input := []byte("schemaVersion: 1\n" +
		"services:\n" +
		"  postgres:\n" +
		"    enabled: false   # CHANGEME after first deploy\n")

	out, err := yaml.New().PatchScalar(input,
		[]string{"services", "postgres", "enabled"}, true)
	if err != nil {
		t.Fatalf("PatchScalar: %v", err)
	}

	got := string(out)
	if !strings.Contains(got, "enabled: true") {
		t.Errorf("output missing patched value; got:\n%s", got)
	}
	if !strings.Contains(got, "CHANGEME after first deploy") {
		t.Errorf("output lost the value-side LineComment; got:\n%s", got)
	}
}

func TestCodec_PatchScalar_AliasIntermediateReturnsErr(t *testing.T) {
	// Why: yaml anchors+aliases (`&base` / `*base`) are valid syntax;
	// silently overwriting through an alias would mutate the anchor
	// target instead of just the patched branch. Reject explicitly
	// with a hint that the caller must inline the anchored target
	// first (review finding #3, M5-T3).
	input := []byte("base: &base\n" +
		"  enabled: false\n" +
		"services:\n" +
		"  postgres: *base\n")

	_, err := yaml.New().PatchScalar(input,
		[]string{"services", "postgres", "enabled"}, true)
	if !errors.Is(err, driven.ErrYAMLPathInvalid) {
		t.Fatalf("PatchScalar(alias intermediate): want ErrYAMLPathInvalid, got %v", err)
	}
	if !strings.Contains(err.Error(), "alias") {
		t.Errorf("error must mention alias to guide the caller; got: %v", err)
	}
}

// --- M5-T4b: PatchMappingEntryYAML + LocateMarkedEntry --------------

const samplePostgresFragment = "image: postgres:16-alpine\n" +
	"environment:\n" +
	"  POSTGRES_USER: postgres\n" +
	"  POSTGRES_PASSWORD: CHANGEME_POSTGRES_PASSWORD\n" +
	"  POSTGRES_DB: postgres\n" +
	"ports:\n" +
	"  - \"5432:5432\"\n" +
	"healthcheck:\n" +
	"  test: [\"CMD\", \"pg_isready\"]\n"

func TestPatchMappingEntryYAML_InsertIntoEmptyInline(t *testing.T) {
	// Why: M3 fresh-init compose has `services: {}` as inline empty.
	// First `u-boot add postgres` must rewrite the parent header to
	// block style and splice in the BEGIN/END-wrapped entry under it
	// without touching the surrounding init block or top-level
	// comments.
	input := []byte("# Compose stack for demo.\n" +
		"name: demo\n" +
		"\n" +
		"services: {}\n" +
		"\n" +
		"volumes: {}\n")
	out, err := yaml.New().PatchMappingEntryYAML(input,
		"services", "postgres", []byte(samplePostgresFragment),
		"service.postgres")
	if err != nil {
		t.Fatalf("PatchMappingEntryYAML: %v", err)
	}
	got := string(out)
	for _, want := range []string{
		"# Compose stack for demo.",
		"name: demo",
		"  # BEGIN U-BOOT MANAGED BLOCK: service.postgres",
		"  postgres:",
		"    image: postgres:16-alpine",
		"  # END U-BOOT MANAGED BLOCK: service.postgres",
		"\nvolumes: {}",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q; got:\n%s", want, got)
		}
	}
	// `services: {}` must be rewritten to block-style `services:`.
	if strings.Contains(got, "services: {}") {
		t.Errorf("services inline-empty header was not rewritten to block style; got:\n%s", got)
	}
}

func TestPatchMappingEntryYAML_ReplaceExisting(t *testing.T) {
	// Why: idempotent re-add — the second `u-boot add postgres` on
	// an already-managed compose.yaml must replace exactly the
	// existing block byte range.
	input := []byte("services:\n" +
		"  # BEGIN U-BOOT MANAGED BLOCK: service.postgres\n" +
		"  postgres:\n" +
		"    image: stale:1\n" +
		"  # END U-BOOT MANAGED BLOCK: service.postgres\n")
	out, err := yaml.New().PatchMappingEntryYAML(input,
		"services", "postgres", []byte(samplePostgresFragment),
		"service.postgres")
	if err != nil {
		t.Fatalf("PatchMappingEntryYAML: %v", err)
	}
	got := string(out)
	if strings.Contains(got, "stale:1") {
		t.Errorf("old block contents not replaced; got:\n%s", got)
	}
	if !strings.Contains(got, "image: postgres:16-alpine") {
		t.Errorf("new block contents missing; got:\n%s", got)
	}
	if strings.Count(got, "BEGIN U-BOOT MANAGED BLOCK: service.postgres") != 1 {
		t.Errorf("expected exactly 1 BEGIN marker, got:\n%s", got)
	}
}

func TestPatchMappingEntryYAML_MissingParentAppends(t *testing.T) {
	// Why: a user-pflichtFiles compose.yaml may lack `services:` entirely.
	// Patch must append a fresh top-level `services:` block at the
	// end of the file separated by one blank line.
	input := []byte("name: demo\n" +
		"networks:\n" +
		"  default: {}\n")
	out, err := yaml.New().PatchMappingEntryYAML(input,
		"services", "postgres", []byte(samplePostgresFragment),
		"service.postgres")
	if err != nil {
		t.Fatalf("PatchMappingEntryYAML: %v", err)
	}
	got := string(out)
	for _, want := range []string{
		"name: demo",
		"networks:",
		"\nservices:",
		"BEGIN U-BOOT MANAGED BLOCK: service.postgres",
		"image: postgres:16-alpine",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q; got:\n%s", want, got)
		}
	}
	// Existing top-level keys must remain byte-identical at the start.
	if !strings.HasPrefix(got, "name: demo\nnetworks:\n  default: {}\n") {
		t.Errorf("prefix mutated; got:\n%s", got)
	}
}

func TestPatchMappingEntryYAML_IndentedMarkerRecognized(t *testing.T) {
	// Why: production writes the marker indented under the parent
	// (column 2 under `services:`). The scanner must locate it
	// regardless of indent.
	input := []byte("services:\n" +
		"  mywebapp:\n" +
		"    image: nginx\n" +
		"  # BEGIN U-BOOT MANAGED BLOCK: service.postgres\n" +
		"  postgres:\n" +
		"    image: stale\n" +
		"  # END U-BOOT MANAGED BLOCK: service.postgres\n")
	out, err := yaml.New().PatchMappingEntryYAML(input,
		"services", "postgres", []byte(samplePostgresFragment),
		"service.postgres")
	if err != nil {
		t.Fatalf("PatchMappingEntryYAML: %v", err)
	}
	got := string(out)
	if !strings.Contains(got, "mywebapp:") {
		t.Errorf("user-custom service was clobbered; got:\n%s", got)
	}
	if strings.Contains(got, "image: stale") {
		t.Errorf("old block not replaced; got:\n%s", got)
	}
	if !strings.Contains(got, "image: postgres:16-alpine") {
		t.Errorf("new block missing; got:\n%s", got)
	}
}

func TestPatchMappingEntryYAML_MalformedMarkerReturnsErr(t *testing.T) {
	input := []byte("services:\n" +
		"  # BEGIN U-BOOT MANAGED BLOCK: service.postgres\n" +
		"  postgres:\n" +
		"    image: foo\n")
	// no END
	_, err := yaml.New().PatchMappingEntryYAML(input,
		"services", "postgres", []byte(samplePostgresFragment),
		"service.postgres")
	if err == nil {
		t.Fatalf("expected malformed-block error, got nil")
	}
	if !strings.Contains(err.Error(), "malformed") {
		t.Errorf("err must mention malformed; got: %v", err)
	}
}

func TestPatchMappingEntryYAML_DuplicateBeginReturnsErr(t *testing.T) {
	input := []byte("services:\n" +
		"  # BEGIN U-BOOT MANAGED BLOCK: service.postgres\n" +
		"  postgres:\n" +
		"  # BEGIN U-BOOT MANAGED BLOCK: service.postgres\n" +
		"  # END U-BOOT MANAGED BLOCK: service.postgres\n")
	_, err := yaml.New().PatchMappingEntryYAML(input,
		"services", "postgres", []byte(samplePostgresFragment),
		"service.postgres")
	if err == nil {
		t.Fatalf("expected malformed-block (duplicate BEGIN) error")
	}
	if !strings.Contains(err.Error(), "malformed") {
		t.Errorf("err must mention malformed; got: %v", err)
	}
}

func TestPatchMappingEntryYAML_TopLevelDuplicateReturnsErr(t *testing.T) {
	input := []byte("name: a\n" +
		"name: b\n" +
		"services: {}\n")
	_, err := yaml.New().PatchMappingEntryYAML(input,
		"services", "postgres", []byte(samplePostgresFragment),
		"service.postgres")
	if !errors.Is(err, driven.ErrYAMLFragmentInvalid) {
		t.Fatalf("err = %v, want ErrYAMLFragmentInvalid", err)
	}
}

func TestPatchMappingEntryYAML_ScalarFragmentReturnsErr(t *testing.T) {
	input := []byte("services: {}\n")
	_, err := yaml.New().PatchMappingEntryYAML(input,
		"services", "postgres", []byte("just-a-scalar\n"),
		"service.postgres")
	if !errors.Is(err, driven.ErrYAMLFragmentInvalid) {
		t.Fatalf("err = %v, want ErrYAMLFragmentInvalid", err)
	}
}

func TestPatchMappingEntryYAML_SequenceFragmentReturnsErr(t *testing.T) {
	input := []byte("services: {}\n")
	_, err := yaml.New().PatchMappingEntryYAML(input,
		"services", "postgres", []byte("- a\n- b\n"),
		"service.postgres")
	if !errors.Is(err, driven.ErrYAMLFragmentInvalid) {
		t.Fatalf("err = %v, want ErrYAMLFragmentInvalid", err)
	}
}

func TestPatchMappingEntryYAML_NonEmptyFlowParentReturnsErr(t *testing.T) {
	// Why: a non-empty flow-style parent cannot host a block-style
	// managed entry without rewriting the whole flow, which would
	// break byte preservation.
	input := []byte("services: { mywebapp: { image: nginx } }\n")
	_, err := yaml.New().PatchMappingEntryYAML(input,
		"services", "postgres", []byte(samplePostgresFragment),
		"service.postgres")
	if !errors.Is(err, driven.ErrYAMLFragmentInvalid) {
		t.Fatalf("err = %v, want ErrYAMLFragmentInvalid", err)
	}
}

func TestPatchMappingEntryYAML_MarkerSomewhereElseRejected(t *testing.T) {
	// Why: defensive — even though the application layer pre-checks
	// via LocateMarkedEntry, the adapter must reject if a marker for
	// the requested service lives outside the requested parent.
	input := []byte("services: {}\n" +
		"\n" +
		"volumes:\n" +
		"  # BEGIN U-BOOT MANAGED BLOCK: service.postgres\n" +
		"  postgres: {}\n" +
		"  # END U-BOOT MANAGED BLOCK: service.postgres\n")
	_, err := yaml.New().PatchMappingEntryYAML(input,
		"services", "postgres", []byte(samplePostgresFragment),
		"service.postgres")
	if !errors.Is(err, driven.ErrYAMLAnchorMismatch) {
		t.Fatalf("err = %v, want ErrYAMLAnchorMismatch", err)
	}
}

func TestPatchMappingEntryYAML_UserEntryWithoutMarkerRejected(t *testing.T) {
	// Why: a user-pflichtFiles `services.postgres:` without our marker is
	// a fachlicher conflict; the adapter must not overwrite it and
	// must not insert a duplicate.
	input := []byte("services:\n" +
		"  postgres:\n" +
		"    image: my-fork/postgres:custom\n")
	_, err := yaml.New().PatchMappingEntryYAML(input,
		"services", "postgres", []byte(samplePostgresFragment),
		"service.postgres")
	if !errors.Is(err, driven.ErrYAMLAnchorMismatch) {
		t.Fatalf("err = %v, want ErrYAMLAnchorMismatch", err)
	}
}

func TestPatchMappingEntryYAML_BytePreservation(t *testing.T) {
	// Why: pins the precise byte-preservation contract — user-custom
	// sibling entries under services: and top-level networks block
	// outside services: must survive a fresh insert byte-identical.
	input := []byte("# Top-level comment.\n" +
		"name: demo\n" +
		"\n" +
		"services:\n" +
		"  mywebapp:\n" +
		"    image: nginx\n" +
		"  # mywebapp belongs to me, do not touch.\n" +
		"\n" +
		"networks:\n" +
		"  default:\n" +
		"    name: demo-default\n")
	out, err := yaml.New().PatchMappingEntryYAML(input,
		"services", "postgres", []byte(samplePostgresFragment),
		"service.postgres")
	if err != nil {
		t.Fatalf("PatchMappingEntryYAML: %v", err)
	}
	got := string(out)
	for _, want := range []string{
		"# Top-level comment.",
		"name: demo",
		"mywebapp:",
		"image: nginx",
		"# mywebapp belongs to me, do not touch.",
		"networks:",
		"name: demo-default",
		"BEGIN U-BOOT MANAGED BLOCK: service.postgres",
		"image: postgres:16-alpine",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q after patch; got:\n%s", want, got)
		}
	}
}

// ---- LocateMarkedEntry ---------------------------------------------

func TestLocateMarkedEntry_CleanReturnsZero(t *testing.T) {
	input := []byte("name: demo\n")
	res, err := yaml.New().LocateMarkedEntry(input,
		"services", "postgres", "service.postgres")
	if err != nil {
		t.Fatalf("LocateMarkedEntry: %v", err)
	}
	if res.EntryExists || res.MarkerInEntry || res.MarkerSomewhereElse {
		t.Errorf("expected zero LocateResult, got %+v", res)
	}
}

func TestLocateMarkedEntry_ManagedReturnsBlockBody(t *testing.T) {
	input := []byte("services:\n" +
		"  # BEGIN U-BOOT MANAGED BLOCK: service.postgres\n" +
		"  postgres:\n" +
		"    image: postgres:16-alpine\n" +
		"  # END U-BOOT MANAGED BLOCK: service.postgres\n")
	res, err := yaml.New().LocateMarkedEntry(input,
		"services", "postgres", "service.postgres")
	if err != nil {
		t.Fatalf("LocateMarkedEntry: %v", err)
	}
	if !res.EntryExists || !res.MarkerInEntry || res.MarkerSomewhereElse {
		t.Errorf("flags wrong: %+v", res)
	}
	if !strings.Contains(string(res.BlockBody), "image: postgres:16-alpine") {
		t.Errorf("BlockBody missing entry; got %q", res.BlockBody)
	}
}

func TestLocateMarkedEntry_UserManualEntryWithoutMarker(t *testing.T) {
	input := []byte("services:\n" +
		"  postgres:\n" +
		"    image: my-fork/postgres:custom\n")
	res, err := yaml.New().LocateMarkedEntry(input,
		"services", "postgres", "service.postgres")
	if err != nil {
		t.Fatalf("LocateMarkedEntry: %v", err)
	}
	if !res.EntryExists {
		t.Errorf("EntryExists should be true; got %+v", res)
	}
	if res.MarkerInEntry || res.MarkerSomewhereElse {
		t.Errorf("no marker present, expected both marker flags false; got %+v", res)
	}
	if res.BlockBody != nil {
		t.Errorf("BlockBody must be nil when MarkerInEntry false; got %q", res.BlockBody)
	}
}

func TestLocateMarkedEntry_MarkerAtWrongAnchor(t *testing.T) {
	input := []byte("services: {}\n" +
		"\n" +
		"volumes:\n" +
		"  # BEGIN U-BOOT MANAGED BLOCK: service.postgres\n" +
		"  something: {}\n" +
		"  # END U-BOOT MANAGED BLOCK: service.postgres\n")
	res, err := yaml.New().LocateMarkedEntry(input,
		"services", "postgres", "service.postgres")
	if err != nil {
		t.Fatalf("LocateMarkedEntry: %v", err)
	}
	if res.EntryExists || res.MarkerInEntry || !res.MarkerSomewhereElse {
		t.Errorf("expected MarkerSomewhereElse only; got %+v", res)
	}
}

func TestLocateMarkedEntry_MalformedReturnsErr(t *testing.T) {
	input := []byte("services:\n" +
		"  # BEGIN U-BOOT MANAGED BLOCK: service.postgres\n" +
		"  postgres: {}\n")
	_, err := yaml.New().LocateMarkedEntry(input,
		"services", "postgres", "service.postgres")
	if err == nil {
		t.Fatalf("expected malformed-block error, got nil")
	}
	if !strings.Contains(err.Error(), "malformed") {
		t.Errorf("err must mention malformed; got: %v", err)
	}
}

func TestLocateMarkedEntry_ParseErrorWrapped(t *testing.T) {
	_, err := yaml.New().LocateMarkedEntry([]byte(":\n  bad yaml"),
		"services", "postgres", "service.postgres")
	if err == nil {
		t.Fatalf("expected parse error, got nil")
	}
	// Must not be one of the structural sentinels.
	for _, sentinel := range []error{driven.ErrYAMLAnchorMismatch, driven.ErrYAMLFragmentInvalid} {
		if errors.Is(err, sentinel) {
			t.Errorf("err %v must not wrap structural sentinel %v", err, sentinel)
		}
	}
}

func TestLocateMarkedEntry_FlowNonEmptyParentTreatedAsAbsent(t *testing.T) {
	// Why: flow-style non-empty parents are rejected by Patch with
	// ErrYAMLFragmentInvalid; Locate reports EntryExists=false so
	// the application layer has a single source of truth and never
	// silently overwrites a flow user-entry.
	input := []byte("services: { mywebapp: { image: nginx } }\n")
	res, err := yaml.New().LocateMarkedEntry(input,
		"services", "postgres", "service.postgres")
	if err != nil {
		t.Fatalf("LocateMarkedEntry: %v", err)
	}
	if res.EntryExists || res.MarkerInEntry || res.MarkerSomewhereElse {
		t.Errorf("expected zero LocateResult for flow-non-empty parent; got %+v", res)
	}
}

func TestPatchMappingEntryYAML_EmptyValueRendersInlineMapping(t *testing.T) {
	// Why: the volume template renders as an empty mapping (`{}`);
	// the renderer must produce `<entryKey>: {}` not a bare colon.
	input := []byte("volumes: {}\n")
	out, err := yaml.New().PatchMappingEntryYAML(input,
		"volumes", "postgres-data", []byte(""), "volume.postgres")
	if err != nil {
		t.Fatalf("Patch: %v", err)
	}
	got := string(out)
	if !strings.Contains(got, "postgres-data: {}") {
		t.Errorf("empty-fragment value not rendered as `entryKey: {}`; got:\n%s", got)
	}
}

func TestPatchMappingEntryYAML_ParseErrorWrapped(t *testing.T) {
	_, err := yaml.New().PatchMappingEntryYAML([]byte(":\n  bad yaml"),
		"services", "postgres", []byte("image: x\n"), "service.postgres")
	if err == nil {
		t.Fatalf("expected parse error, got nil")
	}
}

func TestPatchMappingEntryYAML_FragmentParseErrorWrapped(t *testing.T) {
	_, err := yaml.New().PatchMappingEntryYAML([]byte("services: {}\n"),
		"services", "postgres", []byte("- not\n  : yaml"), "service.postgres")
	if !errors.Is(err, driven.ErrYAMLFragmentInvalid) {
		t.Fatalf("err = %v, want ErrYAMLFragmentInvalid", err)
	}
}

func TestPatchMappingEntryYAML_NullParentTreatedAsEmpty(t *testing.T) {
	// Why: `services: null` and `services: ~` both encode an empty
	// mapping value; the patch must convert to block-style and insert.
	input := []byte("services: null\n")
	out, err := yaml.New().PatchMappingEntryYAML(input,
		"services", "postgres", []byte("image: x\n"), "service.postgres")
	if err != nil {
		t.Fatalf("Patch: %v", err)
	}
	got := string(out)
	if strings.Contains(got, "services: null") {
		t.Errorf("services: null should have been rewritten to block style; got:\n%s", got)
	}
	if !strings.Contains(got, "BEGIN U-BOOT MANAGED BLOCK: service.postgres") {
		t.Errorf("marker missing; got:\n%s", got)
	}
}

func TestLocateMarkedEntry_PatchConsistency(t *testing.T) {
	// Why: shared scanner — Locate and Patch must classify the same
	// content the same way. We test that for every well-formed input
	// where Locate reports MarkerInEntry=true, Patch performs a
	// replace (no error); when Locate reports both marker flags false
	// AND EntryExists=false, Patch inserts (no error).
	managed := []byte("services:\n" +
		"  # BEGIN U-BOOT MANAGED BLOCK: service.postgres\n" +
		"  postgres: {}\n" +
		"  # END U-BOOT MANAGED BLOCK: service.postgres\n")
	res, err := yaml.New().LocateMarkedEntry(managed, "services", "postgres", "service.postgres")
	if err != nil || !res.MarkerInEntry {
		t.Fatalf("setup: Locate must report MarkerInEntry, got %+v err=%v", res, err)
	}
	if _, err := yaml.New().PatchMappingEntryYAML(managed,
		"services", "postgres", []byte("image: x\n"), "service.postgres"); err != nil {
		t.Errorf("Patch on managed content failed: %v", err)
	}

	clean := []byte("services: {}\n")
	res, err = yaml.New().LocateMarkedEntry(clean, "services", "postgres", "service.postgres")
	if err != nil {
		t.Fatalf("clean Locate: %v", err)
	}
	if res.EntryExists || res.MarkerInEntry || res.MarkerSomewhereElse {
		t.Fatalf("clean Locate must be zero, got %+v", res)
	}
	if _, err := yaml.New().PatchMappingEntryYAML(clean,
		"services", "postgres", []byte("image: x\n"), "service.postgres"); err != nil {
		t.Errorf("Patch on clean content failed: %v", err)
	}
}

// --- V1 yaml-parse-error sentinel ----------------------------------

// TestCodec_Unmarshal_WrapsParseError_AsErrYAMLParse pins the
// contract introduced by slice-v1-yaml-parse-error-sentinel: every
// content-parse failure surfaces as a [driven.ErrYAMLParse]-wrapped
// error so application callers can branch via `errors.Is`. Anti-
// drift pin against a future refactor that swaps the wrap for a
// nakedly returned `yaml.Unmarshal` error.
func TestCodec_Unmarshal_WrapsParseError_AsErrYAMLParse(t *testing.T) {
	// Corrupt YAML — leading `:` without a key.
	err := yaml.New().Unmarshal([]byte(":\n  bad\n"), new(map[string]any))
	if err == nil {
		t.Fatalf("expected parse error, got nil")
	}
	if !errors.Is(err, driven.ErrYAMLParse) {
		t.Errorf("err = %v, want wrap of driven.ErrYAMLParse", err)
	}
}

// TestCodec_PatchScalar_WrapsParseError_AsErrYAMLParse extends the
// sentinel cover to the second content-parse site so a content
// argument that fails to parse always routes through the sentinel.
func TestCodec_PatchScalar_WrapsParseError_AsErrYAMLParse(t *testing.T) {
	_, err := yaml.New().PatchScalar([]byte(":\n  bad\n"), []string{"a"}, "v")
	if err == nil {
		t.Fatalf("expected parse error, got nil")
	}
	if !errors.Is(err, driven.ErrYAMLParse) {
		t.Errorf("err = %v, want wrap of driven.ErrYAMLParse", err)
	}
}

// TestCodec_LocateMarkedEntry_WrapsParseError_AsErrYAMLParse covers
// the LocateMarkedEntry content-parse site. Existing
// TestLocateMarkedEntry_ParseErrorWrapped asserted only "an error";
// this pin makes the sentinel explicit.
func TestCodec_LocateMarkedEntry_WrapsParseError_AsErrYAMLParse(t *testing.T) {
	_, err := yaml.New().LocateMarkedEntry([]byte(":\n  bad\n"),
		"services", "postgres", "service.postgres")
	if err == nil {
		t.Fatalf("expected parse error, got nil")
	}
	if !errors.Is(err, driven.ErrYAMLParse) {
		t.Errorf("err = %v, want wrap of driven.ErrYAMLParse", err)
	}
}

// TestCodec_PatchMappingEntryYAML_WrapsParseError_AsErrYAMLParse
// covers the PatchMappingEntryYAML content-parse site (via the
// `assertNoTopLevelDuplicate` helper). Sibling to the existing
// TestPatchMappingEntryYAML_ParseErrorWrapped, which only asserted
// that an error was returned.
func TestCodec_PatchMappingEntryYAML_WrapsParseError_AsErrYAMLParse(t *testing.T) {
	_, err := yaml.New().PatchMappingEntryYAML([]byte(":\n  bad yaml\n"),
		"services", "postgres", []byte("image: x\n"), "service.postgres")
	if err == nil {
		t.Fatalf("expected parse error, got nil")
	}
	if !errors.Is(err, driven.ErrYAMLParse) {
		t.Errorf("err = %v, want wrap of driven.ErrYAMLParse", err)
	}
}

// TestCodec_StripYAMLPrefix pins the M1-review-followup behaviour:
// yaml.v3 messages frequently carry a leading `yaml: ` prefix that
// would surface doubled (`yaml: yaml: ...`) once wrapped via
// `%v`. The stripper removes it so user-facing messages stay clean.
// Indirect test — we trigger a parse error and assert the resulting
// message does NOT carry a doubled prefix.
func TestCodec_StripYAMLPrefix_NoDoubledPrefix(t *testing.T) {
	err := yaml.New().Unmarshal([]byte(":\n  bad\n"), new(map[string]any))
	if err == nil {
		t.Fatalf("expected parse error, got nil")
	}
	if strings.Contains(err.Error(), "yaml: yaml:") {
		t.Errorf("error message has doubled `yaml: yaml:` prefix: %q", err.Error())
	}
}
