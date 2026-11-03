package application

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// The project environment values `devcontainer.containerEnv`
// (lastenheft 0.32.0, Absatz Projekt-Umgebungswerte): the generate
// run renders the map of `u-boot.yaml` into the `containerEnv`
// object of the generated `.devcontainer/devcontainer.json`, beside
// the cache values of the two Go keys, without duplicate entries and
// idempotent — a regeneration rewrites the entries and a removed
// entry disappears. The values belong to the `containerEnv` level
// alone: the Dockerfile block and the hand-written rest carry none
// of them.

// The cache entry names of the two Go keys (SPEC-054, Cache-Werte): the
// fixed pair a project value collides with (Zusage · Kopplung — the
// names bind the collision rejection to the render contract of the
// template's cache pair).
const (
	cacheGoName  = "GOCACHE"
	cacheModName = "GOMODCACHE"
)

// toolchainRenderState is the pre-write validation of the generate run
// around the pinned toolchain (SPEC-054) and the project environment
// values (lastenheft 0.32.0): the SPEC-054 conflicts and checksum reads
// and the collision of a project value with the cache names — each
// rejection is the exit-10 class and stands before the plan and write,
// so a rejected call changes nothing. The helper keeps the two
// conditions out of [GenerateService.generateDevcontainer]'s gocyclo
// budget.
func toolchainRenderState(
	ctx context.Context, dc *ubootYAMLDevcontainer, fetcher driven.Fetcher,
) (toolchainState, error) {
	state, err := resolveToolchain(ctx, dc, fetcher)
	if err != nil {
		return toolchainState{}, fmt.Errorf("generate devcontainer: %w: %s",
			driving.ErrConfigValueInvalid, err.Error())
	}
	if msg := containerEnvConflictMessage(dc, state); msg != "" {
		return toolchainState{}, fmt.Errorf("generate devcontainer: %w: %s",
			driving.ErrConfigValueInvalid, msg)
	}
	return state, nil
}

// containerEnvData is one rendered entry of the project environment
// map in the template data: the sorted entry order keeps the
// generated file byte-deterministic across runs, and both fields
// carry their JSON-escaped content — the map is written per entry
// through `config set`, but a hand edit can name anything, and a
// name or value with a quote or control character must not break
// the generated JSON.
type containerEnvData struct {
	Name  string
	Value string
}

// collectContainerEnvData renders the project environment map in
// sorted entry order with the JSON-escaped fields. nil when the map
// is absent or empty, so the pre-0.32.0 bytes stay.
func collectContainerEnvData(dc *ubootYAMLDevcontainer) []containerEnvData {
	if dc == nil || len(dc.ContainerEnv) == 0 {
		return nil
	}
	names := make([]string, 0, len(dc.ContainerEnv))
	for name := range dc.ContainerEnv {
		names = append(names, name)
	}
	sort.Strings(names)
	out := make([]containerEnvData, 0, len(names))
	for _, name := range names {
		out = append(out, containerEnvData{
			Name:  jsonEnvString(name),
			Value: jsonEnvString(dc.ContainerEnv[name]),
		})
	}
	return out
}

// jsonEnvString escapes one content for the JSON string literal of
// the generated devcontainer.json: quotes and backslashes carry
// their short escapes and every control character below 0x20 its
// `�xx` form, so each entry of the map renders a valid JSON string.
func jsonEnvString(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		default:
			if r < 0x20 {
				fmt.Fprintf(&b, `\u%04x`, r)
			} else {
				b.WriteRune(r)
			}
		}
	}
	return b.String()
}

// containerEnvConflictMessage is the rejection of a colliding pair
// (lastenheft 0.32.0): a project value with one of the cache
// directory names `GOCACHE`/`GOMODCACHE` beside an active
// `go-toolchain` or `golangci-lint` key puts two writers on one
// entry — the generator renders its cache value, the map renders the
// project value — so the run ends before the plan and write. The
// `java-toolchain` key carries no `containerEnv` value and never
// conflicts. Empty when the two go separate ways: without the map,
// without an active Go key, or without a cache name in the map.
func containerEnvConflictMessage(dc *ubootYAMLDevcontainer, state toolchainState) string {
	if dc == nil || len(dc.ContainerEnv) == 0 {
		return ""
	}
	if !state.goToolchain && !state.golangciLint {
		return ""
	}
	var collided, keys []string
	for _, name := range []string{cacheGoName, cacheModName} {
		if _, ok := dc.ContainerEnv[name]; ok {
			collided = append(collided, "devcontainer.containerEnv."+name)
		}
	}
	if len(collided) == 0 {
		return ""
	}
	for _, active := range []struct {
		on  bool
		key string
	}{
		{state.goToolchain, featureGoToolchain},
		{state.golangciLint, featureGolangciLint},
	} {
		if active.on {
			keys = append(keys, active.key)
		}
	}
	keyList := strings.Join(keys, " and ")
	var offs []string
	for _, k := range keys {
		offs = append(offs, "`u-boot config set devcontainer.features."+k+".enabled false`")
	}
	off := strings.Join(offs, " or ")
	return collided[0] + " sets a project value the active key " + keyList +
		" also renders under containerEnv (the cache directories of the Go keys); " +
		"two writers of one entry are not supported — remove the entry from the map in u-boot.yaml " +
		"or switch the key off with `" + off + "`"
}

// containerEnvWarnings is the disclosure of the generate run
// (lastenheft 0.32.0): one plain entry per rendered project value,
// naming the entry and that the value stands in the containerEnv
// level alone. Like the other notices of the run the entries carry
// neither [driving.WarningEntry.SecurityRelaxation] nor
// DegradationFallback, so `--quiet` drops them.
func containerEnvWarnings(dc *ubootYAMLDevcontainer) []driving.WarningEntry {
	if dc == nil || len(dc.ContainerEnv) == 0 {
		return nil
	}
	var out []driving.WarningEntry
	for _, e := range collectContainerEnvData(dc) {
		out = append(out, driving.WarningEntry{
			Code:    "LH-FA-DEV-001",
			Level:   "warn",
			Subject: ".devcontainer/devcontainer.json",
			Message: "project environment value " + e.Name + "=" + e.Value +
				" rendered into the containerEnv object of .devcontainer/devcontainer.json; " +
				"the container processes read the value, the Dockerfile block and the hand-written rest carry none of the values",
		})
	}
	return out
}
