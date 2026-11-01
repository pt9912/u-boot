package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ErrInvalidResourceLimit signals a value of `devcontainer.resources.memory`
// or `devcontainer.resources.cpus` outside its character set
// (LH-FA-CONF-002). The CLI maps it to exit code 10.
var ErrInvalidResourceLimit = errors.New("invalid resource limit")

// The values land verbatim in a JSON string of the generated
// devcontainer.json, so the character set is closed: digits, one `.` for
// `cpus` and one lowercase unit letter for `memory`.
var (
	memoryLimitRE = regexp.MustCompile(`^[0-9]+[bkmg]?$`)
	cpuLimitRE    = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?$`)
)

// MemoryLimit is a validated value of `devcontainer.resources.memory`.
type MemoryLimit string

// NewMemoryLimit accepts a positive integer with an optional unit `b`, `k`,
// `m` or `g` (lowercase). The bounds of the engine (minimum, size of the
// host) are not checked here.
func NewMemoryLimit(raw string) (MemoryLimit, error) {
	if memoryLimitRE.MatchString(raw) && !onlyZeros(strings.TrimRight(raw, "bkmg")) {
		return MemoryLimit(raw), nil
	}
	return "", fmt.Errorf("%w: devcontainer.resources.memory %q; allowed: a positive integer with an optional unit b, k, m or g (for example 512m, 3g)",
		ErrInvalidResourceLimit, raw)
}

// CPULimit is a validated value of `devcontainer.resources.cpus`.
type CPULimit string

// NewCPULimit accepts a positive decimal number with `.` as separator.
func NewCPULimit(raw string) (CPULimit, error) {
	if cpuLimitRE.MatchString(raw) && !onlyZeros(strings.ReplaceAll(raw, ".", "")) {
		return CPULimit(raw), nil
	}
	return "", fmt.Errorf("%w: devcontainer.resources.cpus %q; allowed: a positive decimal number with `.` (for example 2, 1.5)",
		ErrInvalidResourceLimit, raw)
}

func onlyZeros(digits string) bool {
	return strings.Trim(digits, "0") == ""
}
