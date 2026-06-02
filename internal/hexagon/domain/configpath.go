package domain

import (
	"errors"
	"fmt"
	"strings"
)

// ConfigPathKind classifies the kind of `u-boot.yaml` field a
// [ConfigPath] addresses. Used by [application.ConfigService] to
// switch on the read / write semantics per path family without
// re-parsing the dotted string at every call site.
//
// Closed-set enum: only paths the M8 slice explicitly whitelists
// land here. V1 fields (`services.<svc>.persistence`,
// `devcontainer.featureSources.allow`, …) get new kinds in
// their respective V1 slices.
type ConfigPathKind int

const (
	// ConfigProjectName addresses `project.name`. Read + write are
	// both allowed; writes validate the new value through
	// [NewProjectName] (LH-FA-INIT-006).
	ConfigProjectName ConfigPathKind = iota

	// ConfigServiceEnabled addresses `services.<svc>.enabled`. Read
	// is allowed; **write is rejected** by [application.ConfigService]
	// because flipping the boolean alone bypasses the
	// LH-FA-ADD-005 state machine (compose-block / env-block /
	// volume-block atomicity). Users go through `u-boot add` /
	// `u-boot remove` for service toggling. The [ConfigPath]
	// constructor parses the path; the Set-side rejection lives in
	// the use case so the domain layer stays free of lifecycle
	// rules.
	ConfigServiceEnabled

	// ConfigDevcontainerEnabled addresses `devcontainer.enabled`.
	// Read + write both allowed; writes parse the value as
	// bool.
	ConfigDevcontainerEnabled

	// ConfigDevcontainerFeatureEnabled addresses
	// `devcontainer.features.<feature>.enabled` (slice-v1-
	// devcontainer-features T4). Read + write both allowed; writes
	// parse as bool. Unlike services.<svc>.enabled, no
	// state-machine gates the write — feature activation is a pure
	// YAML edit that the renderer (T3) and doctor (T5) consume on
	// next run.
	ConfigDevcontainerFeatureEnabled

	// ConfigDevcontainerFeatureSource addresses
	// `devcontainer.features.<feature>.source` (T4). Read + write
	// both allowed; writes parse as string AND enforce the
	// LH-FA-DEV-003 allowlist rule (the new value must appear in
	// `devcontainer.featureSources.allow`).
	ConfigDevcontainerFeatureSource

	// ConfigDevcontainerFeatureVersion addresses
	// `devcontainer.features.<feature>.version` (T4). Read + write
	// both allowed; writes parse as string with non-empty check
	// only (per-feature pin override).
	ConfigDevcontainerFeatureVersion

	// ConfigDevcontainerFeatureSourcesAllow addresses
	// `devcontainer.featureSources.allow` (T4, spec/lastenheft.md:
	// 717). The path is a LIST, not a scalar — writes append new
	// sources via the special list-path code-route in
	// [application.ConfigService] (PatchScalar handles scalars only).
	// `--allow-external-feature-sources` is the canonical write
	// vector per Spec §712-717.
	ConfigDevcontainerFeatureSourcesAllow
)

// ConfigPath is a typed, whitelisted reference to a leaf in
// `u-boot.yaml` reachable by `u-boot config get/set`. The
// constructor [NewConfigPath] is the only producer; the zero value
// is invalid.
//
// The `WriteAllowed` flag mirrors the M8-T1 §D1 Get/Set table
// rather than a separate constructor per direction. The
// application service inspects `WriteAllowed` before forwarding a
// Set call; Get ignores the flag and always proceeds.
type ConfigPath struct {
	// Kind classifies which u-boot.yaml field the path addresses.
	Kind ConfigPathKind

	// Service is populated only when Kind == ConfigServiceEnabled.
	// The wildcard `<svc>` segment of `services.<svc>.enabled` is
	// format-validated through [NewServiceName]; the catalogue
	// membership check (is `<svc>` a service this u-boot release
	// knows how to add) lives in the application layer because the
	// catalogue is MVP-/V1-phased.
	Service ServiceName

	// Feature is populated only when Kind is one of the
	// ConfigDevcontainerFeature* kinds. The wildcard `<feature>`
	// segment is format-validated through [NewFeatureName]; the
	// catalogue-vs-external decision is the renderer's concern
	// (T3 silently skips unknown names; T4 enforces source-in-
	// allowlist when a Source override is set).
	Feature FeatureName

	// WriteAllowed reports whether `u-boot config set <path>` is
	// permitted for this path. false for ConfigServiceEnabled;
	// true for every other kind. Decoupled from Kind so a future
	// change (e.g. opening services.<svc>.enabled to write once a
	// `u-boot config force-set` flag lands) stays a flag flip
	// instead of an enum split.
	WriteAllowed bool
}

// ErrInvalidConfigPath signals that a raw dotted path does not
// match any entry in the M8 whitelist. The CLI wraps this through
// [driving.ErrConfigPathUnknown] to surface exit code 10
// (LH-FA-CLI-006).
var ErrInvalidConfigPath = errors.New("invalid config path")

// NewConfigPath parses raw and returns the matching [ConfigPath].
// Whitelist:
//
//   - `project.name`                                → ConfigProjectName, write-OK
//   - `devcontainer.enabled`                        → ConfigDevcontainerEnabled, write-OK
//   - `services.<svc>.enabled`                      → ConfigServiceEnabled, write-rejected
//   - `devcontainer.featureSources.allow`           → ConfigDevcontainerFeatureSourcesAllow, write-OK
//   - `devcontainer.features.<feature>.enabled`     → ConfigDevcontainerFeatureEnabled, write-OK
//   - `devcontainer.features.<feature>.source`      → ConfigDevcontainerFeatureSource, write-OK
//   - `devcontainer.features.<feature>.version`     → ConfigDevcontainerFeatureVersion, write-OK
//
// Wildcard segments (`<svc>`, `<feature>`) are parsed through their
// respective domain validators; format failures wrap into
// [ErrInvalidConfigPath] with the cause-error chain preserved.
// Any other dotted path returns [ErrInvalidConfigPath] with the
// unknown segment in the message.
func NewConfigPath(raw string) (ConfigPath, error) {
	switch raw {
	case "project.name":
		return ConfigPath{Kind: ConfigProjectName, WriteAllowed: true}, nil
	case "devcontainer.enabled":
		return ConfigPath{Kind: ConfigDevcontainerEnabled, WriteAllowed: true}, nil
	case "devcontainer.featureSources.allow":
		return ConfigPath{Kind: ConfigDevcontainerFeatureSourcesAllow, WriteAllowed: true}, nil
	}

	// `services.<svc>.enabled` requires a 3-segment split with the
	// middle segment a valid ServiceName.
	if svc, ok := strings.CutPrefix(raw, "services."); ok {
		if name, ok := strings.CutSuffix(svc, ".enabled"); ok {
			service, err := NewServiceName(name)
			if err != nil {
				// Double-wrap so callers can branch on either
				// ErrInvalidConfigPath (generic config-validation
				// path) or ErrInvalidServiceName (specific cause
				// for service-name-format-errors). Requires Go
				// 1.20+ multi-%w semantics.
				return ConfigPath{}, fmt.Errorf("%w: services.<svc>.enabled: %w",
					ErrInvalidConfigPath, err)
			}
			return ConfigPath{
				Kind:         ConfigServiceEnabled,
				Service:      service,
				WriteAllowed: false,
			}, nil
		}
	}

	// `devcontainer.features.<feature>.{enabled,source,version}` —
	// 4-segment split where the third segment is a valid FeatureName.
	if feat, ok := strings.CutPrefix(raw, "devcontainer.features."); ok {
		for _, suffix := range []struct {
			leaf string
			kind ConfigPathKind
		}{
			{".enabled", ConfigDevcontainerFeatureEnabled},
			{".source", ConfigDevcontainerFeatureSource},
			{".version", ConfigDevcontainerFeatureVersion},
		} {
			if name, ok := strings.CutSuffix(feat, suffix.leaf); ok {
				feature, err := NewFeatureName(name)
				if err != nil {
					return ConfigPath{}, fmt.Errorf("%w: devcontainer.features.<feature>%s: %w",
						ErrInvalidConfigPath, suffix.leaf, err)
				}
				return ConfigPath{
					Kind:         suffix.kind,
					Feature:      feature,
					WriteAllowed: true,
				}, nil
			}
		}
	}

	return ConfigPath{}, fmt.Errorf("%w: %q is not a known config path; allowed: project.name, devcontainer.enabled, devcontainer.featureSources.allow, services.<svc>.enabled, devcontainer.features.<feature>.{enabled,source,version}",
		ErrInvalidConfigPath, raw)
}

// String returns the canonical dotted representation of the path.
// Round-trip: `NewConfigPath(p.String())` returns an equal value
// (the [ConfigPath] equality compares all four fields). Used by
// the application service for log lines and error messages.
func (p ConfigPath) String() string {
	switch p.Kind {
	case ConfigProjectName:
		return "project.name"
	case ConfigDevcontainerEnabled:
		return "devcontainer.enabled"
	case ConfigServiceEnabled:
		return "services." + p.Service.String() + ".enabled"
	case ConfigDevcontainerFeatureSourcesAllow:
		return "devcontainer.featureSources.allow"
	case ConfigDevcontainerFeatureEnabled:
		return "devcontainer.features." + p.Feature.String() + ".enabled"
	case ConfigDevcontainerFeatureSource:
		return "devcontainer.features." + p.Feature.String() + ".source"
	case ConfigDevcontainerFeatureVersion:
		return "devcontainer.features." + p.Feature.String() + ".version"
	default:
		return fmt.Sprintf("ConfigPath(kind=%d)", int(p.Kind))
	}
}
