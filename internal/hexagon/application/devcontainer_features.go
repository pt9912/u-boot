package application

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// ErrInvalidFeatureSource signals that a string in the
// `devcontainer.featureSources.allow` list does not satisfy the
// LH-FA-DEV-003 source-format rules. Sentinel-comparable via
// [errors.Is]; wrapped with details (which rule failed) at the call
// site. Maps to exit-code 10 (LH-FA-DEV-003) via the use-case-error-
// to-exit-code wiring.
var ErrInvalidFeatureSource = errors.New("invalid feature source")

// isAllowedFeatureSourceScheme reports whether the given lowercased
// scheme is one of the LH-FA-DEV-003 supported schemes:
//
//   - `http`/`https`: canonical devcontainers/features OCI-via-HTTPS.
//   - `oci`: lower-level OCI scheme some private registries surface.
//
// Implemented as a switch rather than a map to keep the application
// package free of package-level mutable globals (gochecknoglobals).
func isAllowedFeatureSourceScheme(scheme string) bool {
	switch scheme {
	case "http", "https", "oci":
		return true
	}
	return false
}

// validateFeatureSource checks a single allowlist entry against the
// LH-FA-DEV-003 failure-table:
//
//   - empty source string
//   - missing URL scheme
//   - scheme not in [isAllowedFeatureSourceScheme]
//   - missing host component
//
// The function returns nil on success and an [ErrInvalidFeatureSource]-
// wrapped error otherwise. Whitespace-only strings are treated as
// empty (after [strings.TrimSpace]); callers that want a stricter
// "no leading/trailing whitespace" rule should compare before
// trimming.
//
// Implementation deliberately avoids `net/url` so the application
// layer stays free of the net/* stdlib family per the depguard
// `application-no-net` rule (LH-FA-ARCH-003). The parser splits on
// `://` for scheme and `/` for host — sufficient for the Spec's
// non-empty-scheme + non-empty-host shape.
func validateFeatureSource(raw string) error {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return fmt.Errorf("%w: empty source string", ErrInvalidFeatureSource)
	}
	schemeSep := strings.Index(trimmed, "://")
	if schemeSep <= 0 {
		return fmt.Errorf("%w: %q has no URL scheme (expected http://, https://, or oci://)",
			ErrInvalidFeatureSource, trimmed)
	}
	scheme := strings.ToLower(trimmed[:schemeSep])
	if !isAllowedFeatureSourceScheme(scheme) {
		return fmt.Errorf("%w: %q has unsupported scheme %q (expected http, https, or oci)",
			ErrInvalidFeatureSource, trimmed, trimmed[:schemeSep])
	}
	authorityAndPath := trimmed[schemeSep+len("://"):]
	host := authorityAndPath
	if slash := strings.IndexByte(host, '/'); slash >= 0 {
		host = host[:slash]
	}
	if host == "" {
		return fmt.Errorf("%w: %q has no host component",
			ErrInvalidFeatureSource, trimmed)
	}
	return nil
}

// dedupeFeatureSources returns a copy of in with duplicate entries
// removed, preserving the first-occurrence order. Whitespace around
// entries is trimmed before comparison; per
// `spec/lastenheft.md:1352` the dedupe is silent (no error on
// duplicates, the second occurrence is dropped).
//
// The function does NOT validate the entries — callers should run
// [validateFeatureSource] on each entry first if they want format
// rejection. Pairing the two in [normaliseFeatureSources] avoids
// repeating the loop at every call site.
func dedupeFeatureSources(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, raw := range in {
		trimmed := strings.TrimSpace(raw)
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		out = append(out, trimmed)
	}
	return out
}

// normaliseFeatureSources validates every entry in in via
// [validateFeatureSource], then deduplicates via
// [dedupeFeatureSources]. On the first validation failure the
// function returns nil and the wrapped error so callers can map to
// exit-code 10. On success the returned slice carries the trimmed,
// deduplicated entries in first-occurrence order — the canonical
// shape to write back to `u-boot.yaml` and to compare against
// `devcontainer.features.<name>.source` values in T4.
func normaliseFeatureSources(in []string) ([]string, error) {
	for i, raw := range in {
		if err := validateFeatureSource(raw); err != nil {
			return nil, fmt.Errorf("featureSources.allow[%d]: %w", i, err)
		}
	}
	return dedupeFeatureSources(in), nil
}

// validateDevcontainerFeatures runs the LH-FA-DEV-003 schema checks
// on the `devcontainer:` sub-tree that T1 owns:
//
//   - Every key in [ubootYAMLDevcontainer.Features] is a valid
//     [domain.FeatureName].
//   - Every non-empty `Source:` value passes [validateFeatureSource]
//     (URL-format check only — allowlist-membership enforcement is
//     T4's job).
//   - Every entry in [ubootYAMLFeatureSources.Allow] passes
//     [validateFeatureSource] (T4 will additionally silent-dedupe at
//     write time; here we only reject invalid entries).
//
// Returns nil when the sub-tree is empty (nil-Devcontainer is the
// pre-LH-FA-DEV-003 default for projects that never opted in). The
// function is read-only — it neither rewrites nor deduplicates;
// callers that want the canonical shape run [normaliseFeatureSources]
// separately.
//
// `Enabled fehlt` is *not* an error here (Plan: Doctor-Warn, not a
// load-time reject) — that check lives in T5's doctor pipeline.
// `Source set but not in Allow` is also deferred (T4 enforcement).
func validateDevcontainerFeatures(dc *ubootYAMLDevcontainer) error {
	if dc == nil {
		return nil
	}
	if dc.FeatureSources != nil {
		for i, raw := range dc.FeatureSources.Allow {
			if err := validateFeatureSource(raw); err != nil {
				return fmt.Errorf("featureSources.allow[%d]: %w", i, err)
			}
		}
	}
	// Iterate the features map in sorted key order so the *first*
	// error a user sees is deterministic across map-iteration shuffles.
	keys := make([]string, 0, len(dc.Features))
	for k := range dc.Features {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if _, err := domain.NewFeatureName(k); err != nil {
			return fmt.Errorf("features key %q: %w", k, err)
		}
		entry := dc.Features[k]
		if entry.Source != "" {
			if err := validateFeatureSource(entry.Source); err != nil {
				return fmt.Errorf("features.%s.source: %w", k, err)
			}
		}
	}
	return nil
}
