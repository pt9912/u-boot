package application

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// The catalogue key `brew` (LH-FA-DEV-013): a Dockerfile block that checks out
// Homebrew at a commit pinned in the template and installs the formulas of
// `devcontainer.features.brew.packages` as the container user. The OCI feature
// `homebrew` is a different key and a different installation; both active in
// one project is the exit-10 rejection of [brewConflictMessage].
const (
	featureBrew     = "brew"
	featureHomebrew = "homebrew"
)

// brewBlockActive reports whether the `brew` entry is enabled without a
// `source:` override. Only then the Dockerfile carries the block and the
// formula list is asked for its tap approvals; with `enabled: false`, without
// `enabled` or with a `source:` (an external OCI feature) the list is only
// checked for its form.
func brewBlockActive(dc *ubootYAMLDevcontainer) bool {
	return dockerfileBlockFeatureActive(dc, featureBrew)
}

// homebrewFeatureActive reports whether the OCI feature `homebrew` is enabled
// in dc, with or without a `source:` override: any enabled entry at the key
// installs a Homebrew, and which of two installations the commands find is a
// path-order question no test carries (= [goFeatureActive] for the key
// `homebrew`).
func homebrewFeatureActive(dc *ubootYAMLDevcontainer) bool {
	if dc == nil {
		return false
	}
	f, ok := dc.Features[featureHomebrew]
	return ok && f.Enabled != nil && *f.Enabled
}

// brewConflictMessage is the rejection of the active `brew` block next to the
// enabled OCI feature `homebrew`. Empty when the two go separate ways. It names
// the other key, because the two spellings are easy to mix up.
func brewConflictMessage(dc *ubootYAMLDevcontainer) string {
	if !brewBlockActive(dc) || !homebrewFeatureActive(dc) {
		return ""
	}
	return "devcontainer.features.brew and devcontainer.features.homebrew install two Homebrews in one image; " +
		"which one the commands find is an order of PATH and is not supported — " +
		"switch `brew` off with `u-boot config set devcontainer.features.brew.enabled false` " +
		"or the feature `homebrew` with `u-boot config set devcontainer.features.homebrew.enabled false`"
}

// brewPackagesPlacementError is the single statement of the rule "only the key
// brew without a source takes a package list". [validateBrewEntries] (read
// path) and [planBrewPackages] (write path of config set) both ask it.
func brewPackagesPlacementError(feature, source string) error {
	switch {
	case feature != featureBrew:
		return fmt.Errorf("devcontainer.features.%s.packages: only the key brew takes a package list (the Homebrew block of u-boot); the key %s has none", feature, feature)
	case source != "":
		return fmt.Errorf("devcontainer.features.brew.packages: the entry has a source and is an external feature without the Homebrew block of u-boot, so it takes no package list; remove the source or the list")
	}
	return nil
}

// validateBrewEntries checks the form of the two brew lists in dc: every entry
// of `devcontainer.brewTaps.allow` is `<user>/<tap>`, and a package list stands
// only at the key `brew` without `source:` with entries of the form
// `<formula>` or `<user>/<tap>/<formula>`. It does not look at `enabled`: a
// switched-off entry keeps its list valid. The approval of the taps
// ([unlistedBrewTaps]) and the conflict with `homebrew` ([brewConflictMessage])
// are rules of the generate run and of the write path, not of the form.
func validateBrewEntries(dc *ubootYAMLDevcontainer) error {
	if dc == nil {
		return nil
	}
	if dc.BrewTaps != nil {
		for i, raw := range dc.BrewTaps.Allow {
			if _, err := domain.NewBrewTap(raw); err != nil {
				return fmt.Errorf("brewTaps.allow[%d]: %w", i, err)
			}
		}
	}
	for _, name := range sortedFeatureNames(dc.Features) {
		entry := dc.Features[name]
		if len(entry.Packages) == 0 {
			continue
		}
		if err := brewPackagesPlacementError(name, entry.Source); err != nil {
			return err
		}
		for i, raw := range entry.Packages {
			if _, err := domain.NewBrewPackage(raw); err != nil {
				return fmt.Errorf("features.%s.packages[%d]: %w", name, i, err)
			}
		}
	}
	return nil
}

// brewPackages returns the formula list of the active `brew` block in the
// order of `u-boot.yaml`; entries outside the form are skipped (the read path
// rejected them before the render).
func brewPackages(dc *ubootYAMLDevcontainer) []string {
	if !brewBlockActive(dc) {
		return nil
	}
	var out []string
	for _, raw := range dc.Features[featureBrew].Packages {
		if p, err := domain.NewBrewPackage(raw); err == nil {
			out = append(out, p.String())
		}
	}
	return out
}

// brewTapViolation names one tap in the formula list of the active `brew`
// block that is outside `devcontainer.brewTaps.allow`, with the formulas that
// come from it.
type brewTapViolation struct {
	tap      string
	packages []string
}

// unlistedBrewTaps returns the taps of the active `brew` block that are not in
// `devcontainer.brewTaps.allow`, sorted by tap. The comparison is byte-equal,
// like the one of the feature sources ([featureSourceInAllow]); the list of
// `devcontainer.featureSources.allow` and `--yes` free no tap. A formula
// without a tap needs no approval.
func unlistedBrewTaps(dc *ubootYAMLDevcontainer) []brewTapViolation {
	if !brewBlockActive(dc) {
		return nil
	}
	allowed := map[string]struct{}{}
	if dc.BrewTaps != nil {
		allowed = stringSet(dc.BrewTaps.Allow)
	}
	byTap := map[string][]string{}
	for _, raw := range dc.Features[featureBrew].Packages {
		p, err := domain.NewBrewPackage(raw)
		if err != nil || p.Tap() == "" {
			continue
		}
		if _, ok := allowed[p.Tap()]; !ok {
			byTap[p.Tap()] = append(byTap[p.Tap()], p.String())
		}
	}
	taps := make([]string, 0, len(byTap))
	for tap := range byTap {
		taps = append(taps, tap)
	}
	sort.Strings(taps)
	out := make([]brewTapViolation, 0, len(taps))
	for _, tap := range taps {
		out = append(out, brewTapViolation{tap: tap, packages: byTap[tap]})
	}
	return out
}

// brewTapViolationMessage is the text of the rejection of the tap approval,
// shared by `generate devcontainer`, `config set` and `doctor`: per tap the
// tap and its formulas, then the ways out in this order (allow the tap, remove
// the entry, switch the feature off) and the rule of the comparison.
func brewTapViolationMessage(violations []brewTapViolation) string {
	parts := make([]string, 0, len(violations))
	for _, v := range violations {
		quoted := make([]string, 0, len(v.packages))
		for _, p := range v.packages {
			quoted = append(quoted, fmt.Sprintf("%q", p))
		}
		parts = append(parts, fmt.Sprintf(
			"devcontainer.features.brew.packages: the tap %q of %s is not in devcontainer.brewTaps.allow; allow it with `u-boot config set devcontainer.brewTaps.allow %s`, remove the entry from the list in u-boot.yaml, or switch the feature off with `u-boot config set devcontainer.features.brew.enabled false`",
			v.tap, strings.Join(quoted, ", "), v.tap))
	}
	return strings.Join(parts, "; ") +
		". The approval is separate from devcontainer.featureSources.allow, the comparison is byte-equal and --yes does not substitute"
}

// brewTapFinding is the finding of the doctor check
// `devcontainer.features.allowlist` for the approval of the taps: the text of
// [brewTapViolationMessage] and empty without a violation. The conflict with
// `homebrew` is no finding of the check.
func brewTapFinding(dc *ubootYAMLDevcontainer) string {
	if violations := unlistedBrewTaps(dc); len(violations) > 0 {
		return brewTapViolationMessage(violations)
	}
	return ""
}

// brewViolationMessage is the pre-write rule of the active `brew` block, in
// order: the conflict with the OCI feature `homebrew`, then the approval of the
// taps. Empty when the block is inactive or both hold. The generate run and
// the write path of `config set` ask it, so the rules stand once.
func brewViolationMessage(dc *ubootYAMLDevcontainer) string {
	if msg := brewConflictMessage(dc); msg != "" {
		return msg
	}
	return brewTapFinding(dc)
}

// brewEnablingViolation is the rule of the write path of `config set` for the
// `enabled` leaf of the keys `brew` and `homebrew`: the enabling asks
// [brewViolationMessage] on the patched document (the conflict of the two keys
// and the approval of the taps of the formula list). Switching off is always
// possible, at either key; the written value decides, not the state of the
// sibling key. Empty for any other key.
func brewEnablingViolation(dc *ubootYAMLDevcontainer, feature string, enabled bool) string {
	if !enabled || (feature != featureBrew && feature != featureHomebrew) {
		return ""
	}
	return brewViolationMessage(dc)
}

// rejectBrewViolations is the brew part of the pre-write gate of
// [GenerateService.generateDevcontainer]: the exit-10 rejection of
// [brewViolationMessage], before any plan and write, under `--dry-run` too.
func rejectBrewViolations(dc *ubootYAMLDevcontainer) error {
	if msg := brewViolationMessage(dc); msg != "" {
		return fmt.Errorf("generate devcontainer: %w: %s", driving.ErrConfigValueInvalid, msg)
	}
	return nil
}

// brewWarnings returns the notices `generate devcontainer` gives for an active
// `brew` block, in both profiles: one for the block (the source of Homebrew,
// where its pin stands, the formula list as not pinned and not checked) and,
// sorted by tap, one per approved tap that the formula list uses (a source of
// its own). They carry neither [driving.WarningEntry.SecurityRelaxation] nor
// [driving.WarningEntry.DegradationFallback], so `--quiet` drops them. The
// text names neither the pinned commit nor a version: the pin stands once, in
// the Dockerfile.
func brewWarnings(dc *ubootYAMLDevcontainer) []driving.WarningEntry {
	if !brewBlockActive(dc) {
		return nil
	}
	msg := "devcontainer feature brew installs Homebrew from github.com/Homebrew/brew (commit pinned as HOMEBREW_REVISION in .devcontainer/Dockerfile) without an installer script when the image is built; "
	if packages := brewPackages(dc); len(packages) > 0 {
		msg += "it installs the formulas " + strings.Join(packages, ", ") +
			", which are not pinned and not checked, so a later build can install other contents than an earlier one"
	} else {
		msg += "no formula is listed"
	}
	out := []driving.WarningEntry{{
		Code:    "LH-NFA-SEC-004",
		Level:   "warn",
		Message: msg,
		Subject: ".devcontainer/devcontainer.json",
	}}
	for _, tap := range usedBrewTaps(dc) {
		out = append(out, driving.WarningEntry{
			Code:  "LH-NFA-SEC-004",
			Level: "warn",
			Message: "devcontainer feature brew installs formulas from the tap " + tap +
				" (approved in devcontainer.brewTaps.allow) when the image is built; the tap is a source of its own, " +
				"its formulas are code that Homebrew loads from the repository of the tap and runs, and it is not pinned and not checked",
			Subject: ".devcontainer/devcontainer.json",
		})
	}
	return out
}

// usedBrewTaps returns the sorted taps of the formula list of the active
// `brew` block, without duplicates.
func usedBrewTaps(dc *ubootYAMLDevcontainer) []string {
	seen := map[string]struct{}{}
	var taps []string
	for _, raw := range brewPackages(dc) {
		p, err := domain.NewBrewPackage(raw)
		if err != nil || p.Tap() == "" {
			continue
		}
		if _, ok := seen[p.Tap()]; !ok {
			seen[p.Tap()] = struct{}{}
			taps = append(taps, p.Tap())
		}
	}
	sort.Strings(taps)
	return taps
}
