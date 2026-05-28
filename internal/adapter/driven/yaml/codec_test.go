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
