package application

import (
	"fmt"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// validateDevcontainer runs every schema check on the `devcontainer:`
// sub-tree: the LH-FA-DEV-003 feature checks plus the Lastenheft
// 0.3.0 checks of [validateDevcontainerSandbox]. Load-time callers
// (config, generate, doctor) use this single entry point so a new
// key family cannot be forgotten at one of the call sites.
func validateDevcontainer(dc *ubootYAMLDevcontainer) error {
	if err := validateDevcontainerFeatures(dc); err != nil {
		return err
	}
	return validateDevcontainerSandbox(dc)
}

// validateDevcontainerSandbox checks the closed value sets of
// `devcontainer.user.uid` (LH-FA-DEV-004), `devcontainer.profile`
// (LH-FA-DEV-006) and `devcontainer.sandbox.*` (LH-FA-DEV-007).
// Empty / absent values are valid (defaults apply). Errors wrap
// [domain.ErrInvalidSandboxSetting]; callers classify them as
// LH-FA-CLI-006 exit 10.
func validateDevcontainerSandbox(dc *ubootYAMLDevcontainer) error {
	if dc == nil {
		return nil
	}
	if dc.User != nil && dc.User.UID != nil {
		if _, err := domain.NewContainerUID(*dc.User.UID); err != nil {
			return err
		}
	}
	if dc.Profile != "" {
		if _, err := domain.NewDevcontainerProfile(dc.Profile); err != nil {
			return err
		}
	}
	if dc.Sandbox == nil {
		return nil
	}
	if dc.Sandbox.NestedRuntime != "" {
		if _, err := domain.NewNestedRuntime(dc.Sandbox.NestedRuntime); err != nil {
			return err
		}
	}
	if dc.Sandbox.OnUnavailable != "" {
		if _, err := domain.NewOnUnavailable(dc.Sandbox.OnUnavailable); err != nil {
			return err
		}
	}
	return nil
}

// coerceSandboxConfigValue is the Stage-1 coercion of the four
// Lastenheft 0.3.0 config kinds. ok=false means the kind is not one
// of them (caller continues its own dispatch).
func coerceSandboxConfigValue(path domain.ConfigPath, raw string) (coerced any, formatted string, ok bool, err error) {
	var verr error
	switch path.Kind {
	case domain.ConfigDevcontainerUserUID:
		var uid int
		uid, verr = domain.ParseContainerUID(raw)
		coerced, formatted = uid, fmt.Sprint(uid)
	case domain.ConfigDevcontainerProfile:
		var v domain.DevcontainerProfile
		v, verr = domain.NewDevcontainerProfile(raw)
		coerced, formatted = string(v), string(v)
	case domain.ConfigDevcontainerSandboxNestedRuntime:
		var v domain.NestedRuntime
		v, verr = domain.NewNestedRuntime(raw)
		coerced, formatted = string(v), string(v)
	case domain.ConfigDevcontainerSandboxOnUnavailable:
		var v domain.OnUnavailable
		v, verr = domain.NewOnUnavailable(raw)
		coerced, formatted = string(v), string(v)
	default:
		return nil, "", false, nil
	}
	if verr != nil {
		return nil, "", true, fmt.Errorf("%w: %s: %w", driving.ErrConfigValueInvalid, path, verr)
	}
	return coerced, formatted, true, nil
}

// sandboxConfigYAMLPath maps the four kinds to their YAML paths.
func sandboxConfigYAMLPath(path domain.ConfigPath) ([]string, bool) {
	switch path.Kind {
	case domain.ConfigDevcontainerUserUID:
		return []string{"devcontainer", "user", "uid"}, true
	case domain.ConfigDevcontainerProfile:
		return []string{"devcontainer", "profile"}, true
	case domain.ConfigDevcontainerSandboxNestedRuntime:
		return []string{"devcontainer", "sandbox", "nestedRuntime"}, true
	case domain.ConfigDevcontainerSandboxOnUnavailable:
		return []string{"devcontainer", "sandbox", "onUnavailable"}, true
	}
	return nil, false
}

// sandboxConfigValue returns the stored string form of one of the
// four kinds; "" means unset.
func sandboxConfigValue(cfg ubootYAMLConfig, path domain.ConfigPath) (string, bool) {
	dc := cfg.Devcontainer
	switch path.Kind {
	case domain.ConfigDevcontainerUserUID:
		if dc == nil || dc.User == nil || dc.User.UID == nil {
			return "", true
		}
		return fmt.Sprint(*dc.User.UID), true
	case domain.ConfigDevcontainerProfile:
		if dc == nil {
			return "", true
		}
		return dc.Profile, true
	case domain.ConfigDevcontainerSandboxNestedRuntime:
		if dc == nil || dc.Sandbox == nil {
			return "", true
		}
		return dc.Sandbox.NestedRuntime, true
	case domain.ConfigDevcontainerSandboxOnUnavailable:
		if dc == nil || dc.Sandbox == nil {
			return "", true
		}
		return dc.Sandbox.OnUnavailable, true
	}
	return "", false
}
