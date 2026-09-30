package application

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// sandboxCoercers maps each Lastenheft 0.3.x config kind to its
// Stage-1 coercion: raw user string → (typed YAML scalar, canonical
// string). The list kind `egress.allow` is not here: it takes the
// list-path route ([ConfigService.setEgressAllow]).
func sandboxCoercers() map[domain.ConfigPathKind]func(string) (any, string, error) {
	str := func(v string, err error) (any, string, error) { return v, v, err }
	return map[domain.ConfigPathKind]func(string) (any, string, error){
		domain.ConfigDevcontainerUserUID: func(raw string) (any, string, error) {
			uid, err := domain.ParseContainerUID(raw)
			return uid, strconv.Itoa(uid), err
		},
		domain.ConfigDevcontainerProfile: func(raw string) (any, string, error) {
			v, err := domain.NewDevcontainerProfile(raw)
			return str(string(v), err)
		},
		domain.ConfigDevcontainerSandboxNestedRuntime: func(raw string) (any, string, error) {
			v, err := domain.NewNestedRuntime(raw)
			return str(string(v), err)
		},
		domain.ConfigDevcontainerSandboxOnUnavailable: func(raw string) (any, string, error) {
			v, err := domain.NewOnUnavailable(raw)
			return str(string(v), err)
		},
		domain.ConfigDevcontainerSandboxRepository: func(raw string) (any, string, error) {
			repo := strings.TrimSpace(raw)
			return str(repo, validateCloneURL(repo))
		},
		domain.ConfigDevcontainerSandboxEgressEnabled: func(raw string) (any, string, error) {
			b, err := strconv.ParseBool(raw)
			if err != nil {
				return nil, "", fmt.Errorf("%w: expects a bool (true/false/1/0), got %q", domain.ErrInvalidSandboxSetting, raw)
			}
			return b, strconv.FormatBool(b), nil
		},
	}
}

// coerceSandboxConfigValue is the Stage-1 coercion of the scalar
// Lastenheft 0.3.x config kinds. ok=false means the kind is not one of
// them (caller continues its own dispatch).
func coerceSandboxConfigValue(path domain.ConfigPath, raw string) (coerced any, formatted string, ok bool, err error) {
	coerce, found := sandboxCoercers()[path.Kind]
	if !found {
		return nil, "", false, nil
	}
	coerced, formatted, verr := coerce(raw)
	if verr != nil {
		return nil, "", true, fmt.Errorf("%w: %s: %w", driving.ErrConfigValueInvalid, path, verr)
	}
	return coerced, formatted, true, nil
}

// sandboxYAMLPaths maps the scalar kinds to their YAML paths.
func sandboxYAMLPaths() map[domain.ConfigPathKind][]string {
	return map[domain.ConfigPathKind][]string{
		domain.ConfigDevcontainerUserUID:              {"devcontainer", "user", "uid"},
		domain.ConfigDevcontainerProfile:              {"devcontainer", "profile"},
		domain.ConfigDevcontainerSandboxNestedRuntime: {"devcontainer", "sandbox", "nestedRuntime"},
		domain.ConfigDevcontainerSandboxOnUnavailable: {"devcontainer", "sandbox", "onUnavailable"},
		domain.ConfigDevcontainerSandboxRepository:    {"devcontainer", "sandbox", "repository"},
		domain.ConfigDevcontainerSandboxEgressEnabled: {"devcontainer", "sandbox", "egress", "enabled"},
	}
}

// sandboxConfigYAMLPath returns the YAML path of a scalar kind.
func sandboxConfigYAMLPath(path domain.ConfigPath) ([]string, bool) {
	p, ok := sandboxYAMLPaths()[path.Kind]
	return p, ok
}

// sandboxConfigValue returns the stored string form of one of the
// Lastenheft 0.3.x kinds (including the `egress.allow` list, comma-
// joined); "" means unset.
func sandboxConfigValue(cfg ubootYAMLConfig, path domain.ConfigPath) (string, bool) {
	var dc ubootYAMLDevcontainer
	if cfg.Devcontainer != nil {
		dc = *cfg.Devcontainer
	}
	var sb ubootYAMLDevcontainerSandbox
	if dc.Sandbox != nil {
		sb = *dc.Sandbox
	}
	var eg ubootYAMLSandboxEgress
	if sb.Egress != nil {
		eg = *sb.Egress
	}
	uid, enabled := "", ""
	if dc.User != nil && dc.User.UID != nil {
		uid = strconv.Itoa(*dc.User.UID)
	}
	if eg.Enabled != nil {
		enabled = strconv.FormatBool(*eg.Enabled)
	}
	values := map[domain.ConfigPathKind]string{
		domain.ConfigDevcontainerUserUID:              uid,
		domain.ConfigDevcontainerProfile:              dc.Profile,
		domain.ConfigDevcontainerSandboxNestedRuntime: sb.NestedRuntime,
		domain.ConfigDevcontainerSandboxOnUnavailable: sb.OnUnavailable,
		domain.ConfigDevcontainerSandboxRepository:    sb.Repository,
		domain.ConfigDevcontainerSandboxEgressEnabled: enabled,
		domain.ConfigDevcontainerSandboxEgressAllow:   strings.Join(eg.Allow, ","),
	}
	v, ok := values[path.Kind]
	return v, ok
}
