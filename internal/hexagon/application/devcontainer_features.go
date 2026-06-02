package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// ErrInvalidFeatureSource is the domain sentinel for LH-FA-DEV-003
// source-format violations, re-exported here so existing call sites
// in this package keep working after the slice-v1-devcontainer-
// features review-followup R1 moved the sentinel to the domain
// layer. New code should reference [domain.ErrInvalidFeatureSource]
// directly; this alias stays for backward compatibility with the
// T1-T4 test surface and existing wrap-sites.
var ErrInvalidFeatureSource = domain.ErrInvalidFeatureSource

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

// stringSet converts a key slice into a lookup-only set. Used in
// this package by the drift detector ([classifyDriftCase1]) and by
// other set-membership-style comparisons. Slice-followup-
// devcontainer-features-drift-doctor Review-Followup S4 unified
// this with the previous `keysAsSet` doctor-local helper.
func stringSet(keys []string) map[string]struct{} {
	out := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		out[k] = struct{}{}
	}
	return out
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

// featureCatalogueEntry is the per-feature configuration shared by
// the render pipeline (T3: source + defaultVersion → JSON key form
// `<source>:<version>`) and the catalogue lookup (T2: name →
// source/version resolution when [ubootYAMLDevcontainerFeature.Source]
// is empty). Slice-v1-devcontainer-features T0-(c) defines the
// shape; adding a new built-in feature means: drop a new entry into
// [featureCatalogue], done — no other code change.
//
// `source` is the canonical OCI ref **without** URL scheme (e.g.
// `ghcr.io/devcontainers/features/node`) because devcontainer.json
// `features:` keys use the OCI-ref form, not the URL form. The
// Allowlist (user-provided) carries URL-form entries with scheme;
// the T4 enforcement reconciles the two when a user supplies a
// `source:` override.
type featureCatalogueEntry struct {
	source         string
	defaultVersion string
	shortDesc      string
}

// featureCatalogue lists the built-in devcontainer-feature catalogue
// per spec/lastenheft.md:692-707 + T0-Outcomes (c). Built-in entries
// are aktivierbar without an Allowlist entry — slice-v1-devcontainer-
// features §AK "Statischer Katalog". External features that don't
// match any catalogue key must declare a `source:` override and have
// that source in [ubootYAMLFeatureSources.Allow] (T4 enforcement).
//
// The default-version slug is `1` across the board today; once
// upstream tags drift apart the per-entry pin can be adjusted
// without changing call sites.
func featureCatalogue() map[string]featureCatalogueEntry {
	return map[string]featureCatalogueEntry{
		"git": {
			source:         "ghcr.io/devcontainers/features/git",
			defaultVersion: "1",
			shortDesc:      "Git CLI",
		},
		"docker-cli": {
			source:         "ghcr.io/devcontainers/features/docker-outside-of-docker",
			defaultVersion: "1",
			shortDesc:      "Docker CLI (outside-of-docker)",
		},
		"node": {
			source:         "ghcr.io/devcontainers/features/node",
			defaultVersion: "1",
			shortDesc:      "Node.js",
		},
		"java": {
			source:         "ghcr.io/devcontainers/features/java",
			defaultVersion: "1",
			shortDesc:      "Java + SDKMAN",
		},
		"go": {
			source:         "ghcr.io/devcontainers/features/go",
			defaultVersion: "1",
			shortDesc:      "Go toolchain",
		},
		"cpp": {
			source:         "ghcr.io/devcontainers/features/cpp",
			defaultVersion: "1",
			shortDesc:      "C++ toolchain",
		},
		"kubectl-helm": {
			source:         "ghcr.io/devcontainers/features/kubectl-helm-minikube",
			defaultVersion: "1",
			shortDesc:      "kubectl + helm + minikube",
		},
		"postgres-client": {
			source:         "ghcr.io/devcontainers/features/postgresql-client",
			defaultVersion: "1",
			shortDesc:      "PostgreSQL client",
		},
	}
}

// featureFor returns the catalogue entry for the given feature name.
// The boolean second return mirrors map-lookup convention so callers
// can branch on "unknown feature" without panicking — used by T3's
// renderer to decide between catalogue-lookup and source-override.
func featureFor(name domain.FeatureName) (featureCatalogueEntry, bool) {
	entry, ok := featureCatalogue()[name.String()]
	return entry, ok
}

// devcontainerFeatureData is the per-feature projection that
// [templateData.Features] carries into the devcontainer.json
// renderer. Source + Version compose the JSONC feature key as
// `"<Source>:<Version>": {}`. Slice-v1-devcontainer-features T3.
type devcontainerFeatureData struct {
	Source  string
	Version string
}

// collectDevcontainerFeatures projects the enabled entries from
// `cfg.Devcontainer.Features` into the renderer's sorted feature
// list. Per slice-v1-devcontainer-features T3:
//
//   - Skip entries with `Enabled == nil` (T5-doctor-Warn, not a
//     load-time reject) or `*Enabled == false` (registered but
//     deactivated).
//   - When an entry's `Source` override is non-empty, render it
//     verbatim (T4 enforces allowlist membership — T3 trusts the
//     value here). When `Source` is empty, look up the canonical
//     OCI ref via [featureFor]; unknown names without a source
//     override are silently skipped here so T4 can surface them
//     with the proper Exit-Code-10 path.
//   - When an entry's `Version` override is non-empty, use it;
//     otherwise fall back to the catalogue's `defaultVersion`. For
//     a Source-override without Version override, default to "1"
//     (the upstream devcontainers/features convention) so external
//     features render as `"<source>:1": {}` rather than dangling
//     on a missing colon.
//   - Sort the result alphabetically by Source so the rendered
//     JSON is deterministic across map-iteration shuffles.
//
// Returns nil when no features are enabled — the template skips
// the `"features": { … }` key entirely in that case (preserves
// byte-equality with pre-T3 devcontainer.json files).
func collectDevcontainerFeatures(cfg ubootYAMLConfig) []devcontainerFeatureData {
	if cfg.Devcontainer == nil || len(cfg.Devcontainer.Features) == 0 {
		return nil
	}
	out := make([]devcontainerFeatureData, 0, len(cfg.Devcontainer.Features))
	for name, entry := range cfg.Devcontainer.Features {
		if entry.Enabled == nil || !*entry.Enabled {
			continue
		}
		data, ok := projectFeatureEntry(name, entry)
		if !ok {
			continue
		}
		out = append(out, data)
	}
	// Audit-Followup A2: sort by the full render-key (Source ":" Version)
	// rather than Source alone. Two enabled features with the same
	// Source but different Versions otherwise inherit Go's
	// map-iteration randomness and break the LH-FA-DEV-005 idempotency
	// contract (`generate devcontainer` flipped UpdatedBlock vs NoOp
	// across runs).
	sort.Slice(out, func(i, j int) bool {
		if out[i].Source != out[j].Source {
			return out[i].Source < out[j].Source
		}
		return out[i].Version < out[j].Version
	})
	return out
}

// projectFeatureEntry resolves one feature entry into the
// renderer's projection. Returns ok=false when the entry has no
// `Source:` override AND the name is not a built-in catalogue key —
// the silent-skip behaviour T3 needs so T4 can surface the failure
// with the proper Exit-Code-10 path.
//
// Important: this helper does NOT inspect `entry.Enabled` — the
// caller is responsible for the enabled-filter.
// [collectDevcontainerFeatures] filters before calling;
// [projectAllFeatureEntries] (drift-detector) deliberately keeps
// disabled entries projected. Both rely on the shared
// [projectFeatureSourceVersion] core so a future "skip-on-disabled"
// addition would have to be made there explicitly — not by accident
// in this wrapper.
func projectFeatureEntry(name string, entry ubootYAMLDevcontainerFeature) (devcontainerFeatureData, bool) {
	source, version, ok := projectFeatureSourceVersion(name, entry.Source, entry.Version)
	if !ok {
		return devcontainerFeatureData{}, false
	}
	return devcontainerFeatureData{Source: source, Version: version}, true
}

// projectFeatureSourceVersion is the core source-/version-projection
// shared by [projectFeatureEntry] (renderer-side) and
// [projectAllFeatureEntries] (drift-detector-side). It deliberately
// takes scalar source/version inputs instead of a
// ubootYAMLDevcontainerFeature, so callers cannot accidentally drag
// in an Enabled-dependent branch.
//
// Returns the resolved (source, version, ok=true) when either:
//
//   - the caller supplied a non-empty Source override (Allowlist-
//     conformance is enforced elsewhere — this helper just resolves);
//     defaultVersion "1" applies when Version is empty.
//   - the name matches a built-in catalogue entry; the catalogue's
//     defaultVersion applies when Version is empty.
//
// Returns ("", "", false) for the orphan case (no Source override
// and the name is not in the catalogue) — caller skips silently.
// Slice-followup-devcontainer-features-drift-doctor Review-Followup
// S1 extracted this helper to remove the synthetic-`Enabled=&true`-
// probe in [projectAllFeatureEntries].
func projectFeatureSourceVersion(name, source, version string) (resolvedSource, resolvedVersion string, ok bool) {
	if source != "" {
		v := version
		if v == "" {
			v = "1"
		}
		return source, v, true
	}
	featureName, err := domain.NewFeatureName(name)
	if err != nil {
		return "", "", false
	}
	cat, found := featureFor(featureName)
	if !found {
		return "", "", false
	}
	v := version
	if v == "" {
		v = cat.defaultVersion
	}
	return cat.source, v, true
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

// ErrDevcontainerJSONUnparsable signals that the devcontainer.json
// (after stripJSONC) is not valid JSON. The drift check returns
// it so the caller can skip with a clear repair-hint while the
// primary file-validity diagnostic (`devcontainer.json.valid`,
// from M5-T7's checkDevcontainerJSON) handles the severity side.
// Slice-followup-devcontainer-features-drift-doctor T1.
var ErrDevcontainerJSONUnparsable = errors.New("devcontainer.json unparsable after stripJSONC")

// driftJSONFeatureKeys reads `.devcontainer/devcontainer.json` and
// returns the keys of its `features:` map. The three outcomes:
//
//   - File absent → (nil, nil, nil): drift detector treats this
//     as "no JSON keys"; Case 1 still fires when cfg has enabled
//     entries.
//   - File present + valid JSONC → (keys, present=true, nil) where
//     `keys` may be empty if the JSON has no `features:` section
//     (also a legitimate "no JSON keys" state for Case 1).
//   - File present + unparsable → (nil, true,
//     ErrDevcontainerJSONUnparsable): caller short-circuits to
//     skip; the validity diagnostic
//     (`devcontainer.json.valid`) reports the parse failure
//     separately.
//
// fs is the driven port (no direct os.* calls in the application
// layer). baseDir is the project root.
func driftJSONFeatureKeys(fs driven.FileSystem, baseDir string) (keys []string, present bool, err error) {
	path := filepath.Join(baseDir, ".devcontainer", "devcontainer.json")
	exists, ferr := fs.Exists(path)
	if ferr != nil {
		// FS errors are caller-classified — bubble up so the
		// dispatcher decides between skip and a doctor-level
		// error.
		return nil, false, ferr
	}
	if !exists {
		return nil, false, nil
	}
	body, ferr := fs.ReadFile(path)
	if ferr != nil {
		return nil, true, ferr
	}
	stripped := stripJSONC(body)
	var shape struct {
		Features map[string]json.RawMessage `json:"features"`
	}
	if jerr := json.Unmarshal(stripped, &shape); jerr != nil {
		return nil, true, fmt.Errorf("%w: %v", ErrDevcontainerJSONUnparsable, jerr)
	}
	keys = make([]string, 0, len(shape.Features))
	for k := range shape.Features {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys, true, nil
}

// featureDriftProjection holds the two sets the drift detector
// builds from `cfg.Devcontainer.Features`:
//
//   - expectedKeys: render-keys of all entries with Enabled = &true.
//     These are the keys that MUST appear in the JSON `features:`
//     map for the project to be drift-free (Case 1 = expected \
//     jsonKeys).
//   - knownProjectableKeys: render-keys of ALL projectable entries
//     (enabled, disabled, or unset). Lets the detector distinguish
//     Case 2a (disabled feature still in JSON — render-key is
//     in known but not in expected) from Case 2b (JSON key
//     completely foreign — not in known at all).
//
// "Projectable" excludes entries that projectFeatureEntry would
// silently skip: orphan activations without a `source:` override.
// Those land in [DoctorService.checkDevcontainerFeaturesAllowlist]
// (Teil A), not in the drift check.
type featureDriftProjection struct {
	expectedKeys         map[string]struct{}
	knownProjectableKeys map[string]struct{}
}

// projectAllFeatureEntries builds the [featureDriftProjection] from
// every entry in `cfg.Devcontainer.Features` — regardless of
// `Enabled`. Reuses the same [projectFeatureEntry] helper that the
// T3 generator uses, so the cfg-side render-key is byte-identical
// to whatever the generator would emit (S2 finding from the
// slice-plan review).
//
// The detector then performs three set-differences:
//
//   - Case 1: expectedKeys \ jsonKeys
//   - Case 2a: (jsonKeys ∩ knownProjectableKeys) \ expectedKeys
//   - Case 2b: jsonKeys \ knownProjectableKeys
//
// Caller responsibility: nil-cfg / nil-Devcontainer handling lives
// at the doctor-check call site so the skip-disziplin (precise
// nil-vs-empty-vs-populated decision) stays observable from one
// place.
func projectAllFeatureEntries(cfg ubootYAMLConfig) featureDriftProjection {
	out := featureDriftProjection{
		expectedKeys:         map[string]struct{}{},
		knownProjectableKeys: map[string]struct{}{},
	}
	if cfg.Devcontainer == nil {
		return out
	}
	for name, entry := range cfg.Devcontainer.Features {
		// Review-Followup S1: direct call to
		// [projectFeatureSourceVersion] (Enabled-agnostic core)
		// instead of the previous synthetic `Enabled=&true`-probe
		// against [projectFeatureEntry]. Both pathways share one
		// definition of "projectable", so a future change to
		// Enabled-handling cannot accidentally regress this
		// detector. The renderer-skip case (orphan: no catalogue
		// + no source override) still returns ok=false, so
		// knownProjectableKeys is never polluted with names the
		// renderer would silently drop.
		source, version, ok := projectFeatureSourceVersion(name, entry.Source, entry.Version)
		if !ok {
			continue
		}
		renderKey := source + ":" + version
		out.knownProjectableKeys[renderKey] = struct{}{}
		if entry.Enabled != nil && *entry.Enabled {
			out.expectedKeys[renderKey] = struct{}{}
		}
	}
	return out
}
