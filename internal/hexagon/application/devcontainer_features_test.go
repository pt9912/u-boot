package application_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// TestValidateFeatureSource pins the T1 failure-table from
// slice-v1-devcontainer-features for `devcontainer.featureSources.allow`
// entries.
func TestValidateFeatureSource(t *testing.T) {
	t.Parallel()
	t.Run("Accepts", func(t *testing.T) {
		t.Parallel()
		cases := []string{
			"https://ghcr.io/devcontainers/features/node",
			"http://example.test/features/x",
			"oci://registry.local/features/custom",
			"HTTPS://UPPER.CASE/path", // scheme is case-insensitive per RFC 3986
			"  https://leading.trailing/whitespace  ", // trimmed before parse
		}
		for _, in := range cases {
			in := in
			t.Run(in, func(t *testing.T) {
				t.Parallel()
				if err := application.ValidateFeatureSourceForTest(in); err != nil {
					t.Errorf("ValidateFeatureSource(%q): unexpected error: %v", in, err)
				}
			})
		}
	})

	t.Run("Rejects", func(t *testing.T) {
		t.Parallel()
		cases := []struct {
			name   string
			input  string
			expect string // substring expected in error message
		}{
			{"empty", "", "empty source string"},
			{"whitespace only", "   ", "empty source string"},
			{"no scheme", "ghcr.io/devcontainers/features/node", "no URL scheme"},
			{"unsupported scheme", "ftp://example.test/feature", "unsupported scheme"},
			{"file scheme rejected", "file:///local/path", "unsupported scheme"},
			{"no host (https://)", "https://", "no host component"},
			{"no host (https:///path)", "https:///path", "no host component"},
		}
		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				err := application.ValidateFeatureSourceForTest(tc.input)
				if err == nil {
					t.Fatalf("ValidateFeatureSource(%q): expected error, got nil", tc.input)
				}
				if !errors.Is(err, application.ErrInvalidFeatureSource) {
					t.Errorf("err = %v, want wrap of ErrInvalidFeatureSource", err)
				}
				if !strings.Contains(err.Error(), tc.expect) {
					t.Errorf("err = %v, want substring %q", err, tc.expect)
				}
			})
		}
	})
}

// TestDedupeFeatureSources pins the silent-dedupe contract from
// `spec/lastenheft.md:1352`: duplicates are dropped silently,
// first-occurrence order is preserved, whitespace is trimmed before
// comparison.
func TestDedupeFeatureSources(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		in   []string
		want []string
	}{
		{
			name: "no duplicates",
			in:   []string{"https://a.test/x", "https://b.test/y"},
			want: []string{"https://a.test/x", "https://b.test/y"},
		},
		{
			name: "duplicate dropped (second)",
			in:   []string{"https://a.test/x", "https://b.test/y", "https://a.test/x"},
			want: []string{"https://a.test/x", "https://b.test/y"},
		},
		{
			name: "whitespace difference treated as duplicate",
			in:   []string{"https://a.test/x", "  https://a.test/x  "},
			want: []string{"https://a.test/x"},
		},
		{
			name: "empty input",
			in:   nil,
			want: []string{},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := application.DedupeFeatureSourcesForTest(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Dedupe(%v) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

// TestNormaliseFeatureSources pins the combined validate+dedupe
// pipeline: first invalid entry short-circuits with the index
// surfaced in the wrapped message; valid input is returned in
// trimmed + deduped form.
func TestNormaliseFeatureSources(t *testing.T) {
	t.Parallel()

	t.Run("happy path trims and dedupes", func(t *testing.T) {
		t.Parallel()
		got, err := application.NormaliseFeatureSourcesForTest([]string{
			"  https://a.test/x  ",
			"https://b.test/y",
			"https://a.test/x",
		})
		if err != nil {
			t.Fatalf("Normalise: unexpected error %v", err)
		}
		want := []string{"https://a.test/x", "https://b.test/y"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("first invalid entry short-circuits with index", func(t *testing.T) {
		t.Parallel()
		_, err := application.NormaliseFeatureSourcesForTest([]string{
			"https://a.test/x",
			"", // index 1 is bad
			"https://b.test/y",
		})
		if err == nil {
			t.Fatalf("Normalise: expected error, got nil")
		}
		if !errors.Is(err, application.ErrInvalidFeatureSource) {
			t.Errorf("err = %v, want wrap of ErrInvalidFeatureSource", err)
		}
		if !strings.Contains(err.Error(), "featureSources.allow[1]") {
			t.Errorf("err = %v, want index marker featureSources.allow[1]", err)
		}
	})
}

// TestValidateDevcontainerFeatures pins the T1 schema-validation
// contract for `devcontainer.features.<name>` map keys and source
// fields, plus the allowlist-source format checks driven by the
// integrated validateFeatureSource call.
func TestValidateDevcontainerFeatures(t *testing.T) {
	t.Parallel()

	t.Run("nil devcontainer is accepted", func(t *testing.T) {
		t.Parallel()
		// Project without a devcontainer subtree: validator must be
		// a no-op so pre-LH-FA-DEV-003 projects keep loading.
		yaml := []byte("schemaVersion: 1\nproject:\n  name: demo\n")
		if err := application.ValidateDevcontainerFeaturesForTest(t, yaml); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("empty devcontainer is accepted", func(t *testing.T) {
		t.Parallel()
		yaml := []byte("schemaVersion: 1\nproject:\n  name: demo\ndevcontainer:\n  enabled: true\n")
		if err := application.ValidateDevcontainerFeaturesForTest(t, yaml); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("valid features map and allowlist", func(t *testing.T) {
		t.Parallel()
		yaml := []byte(`schemaVersion: 1
project:
  name: demo
devcontainer:
  enabled: true
  featureSources:
    allow:
      - https://ghcr.io/orgX/features/custom-rust
  features:
    node:
      enabled: true
    java:
      enabled: true
      version: "21"
    custom-rust:
      enabled: true
      source: https://ghcr.io/orgX/features/custom-rust
`)
		if err := application.ValidateDevcontainerFeaturesForTest(t, yaml); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("invalid feature name rejected with domain wrap", func(t *testing.T) {
		t.Parallel()
		yaml := []byte(`schemaVersion: 1
project:
  name: demo
devcontainer:
  features:
    Bad_Name:
      enabled: true
`)
		err := application.ValidateDevcontainerFeaturesForTest(t, yaml)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, domain.ErrInvalidFeatureName) {
			t.Errorf("err = %v, want wrap of domain.ErrInvalidFeatureName", err)
		}
	})

	t.Run("invalid allowlist entry rejected with source wrap", func(t *testing.T) {
		t.Parallel()
		yaml := []byte(`schemaVersion: 1
project:
  name: demo
devcontainer:
  featureSources:
    allow:
      - https://a.test/ok
      - not-a-url
`)
		err := application.ValidateDevcontainerFeaturesForTest(t, yaml)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, application.ErrInvalidFeatureSource) {
			t.Errorf("err = %v, want wrap of application.ErrInvalidFeatureSource", err)
		}
		if !strings.Contains(err.Error(), "featureSources.allow[1]") {
			t.Errorf("err = %v, want index marker [1]", err)
		}
	})

	t.Run("invalid features.<name>.source rejected with source wrap", func(t *testing.T) {
		t.Parallel()
		yaml := []byte(`schemaVersion: 1
project:
  name: demo
devcontainer:
  features:
    custom:
      enabled: true
      source: ftp://wrong.scheme/feature
`)
		err := application.ValidateDevcontainerFeaturesForTest(t, yaml)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, application.ErrInvalidFeatureSource) {
			t.Errorf("err = %v, want wrap of ErrInvalidFeatureSource", err)
		}
		if !strings.Contains(err.Error(), "features.custom.source") {
			t.Errorf("err = %v, want path marker features.custom.source", err)
		}
	})

	t.Run("deterministic first-error across map iterations", func(t *testing.T) {
		t.Parallel()
		// Two invalid feature names: assert the *sorted-key-first*
		// one is the surfaced error. Map iteration is randomised in
		// Go, so this would be flaky without the explicit sort in
		// validateDevcontainerFeatures.
		yaml := []byte(`schemaVersion: 1
project:
  name: demo
devcontainer:
  features:
    Zzz_bad:
      enabled: true
    Aaa_bad:
      enabled: true
`)
		// Run multiple times so a non-deterministic implementation
		// would surface different errors across iterations.
		for i := 0; i < 5; i++ {
			err := application.ValidateDevcontainerFeaturesForTest(t, yaml)
			if err == nil {
				t.Fatalf("iter %d: expected error, got nil", i)
			}
			if !strings.Contains(err.Error(), "Aaa_bad") {
				t.Errorf("iter %d: err = %v, want first-by-sort-order key Aaa_bad", i, err)
			}
		}
	})
}
