package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// ErrInvalidSandboxStatus signals that a sandbox status document
// (LH-FA-DEV-007: the file a container start script writes per
// capability) is outside its closed format. Every deviation is the
// same class: the document is a hint from a container that runs an
// agent, so nothing outside the closed set is interpreted.
var ErrInvalidSandboxStatus = errors.New("invalid sandbox status document")

// SandboxStatusMaxBytes is the size limit of a status document.
const SandboxStatusMaxBytes = 512

const (
	sandboxStatusHeader     = "u-boot-status 1"
	sandboxStatusResultKey  = "result="
	sandboxStatusAtKey      = "at="
	sandboxStatusTimeLayout = "2006-01-02T15:04:05Z"
)

var sandboxStatusTimeShape = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)

// SandboxStatusResult is the closed value set of `result` in a status
// document.
type SandboxStatusResult string

const (
	// SandboxStatusOK means the capability was available at the start.
	SandboxStatusOK SandboxStatusResult = "ok"
	// SandboxStatusFallback means the start script degraded (policy
	// `warn`): `vfs` storage, no daemon, or no egress restriction.
	SandboxStatusFallback SandboxStatusResult = "fallback"
	// SandboxStatusFailed means the start script ended with exit code
	// 11 (policy `fail`).
	SandboxStatusFailed SandboxStatusResult = "failed"
)

// SandboxStatus is the parsed status document of one capability. It
// carries only the resolved fields; the raw text is never kept.
type SandboxStatus struct {
	result SandboxStatusResult
	at     time.Time
}

// Result returns the reported outcome.
func (s SandboxStatus) Result() SandboxStatusResult { return s.result }

// At returns the time the start script reported (UTC).
func (s SandboxStatus) At() time.Time { return s.at }

// ParseSandboxStatus parses a status document strictly: at most
// [SandboxStatusMaxBytes] bytes of printable ASCII and newlines, exactly
// the three lines `u-boot-status 1`, `result=ok|fallback|failed` and
// `at=<YYYY-MM-DDTHH:MM:SSZ>` (one optional final newline). Anything
// else, including free text, wraps [ErrInvalidSandboxStatus].
func ParseSandboxStatus(raw []byte) (SandboxStatus, error) {
	if len(raw) > SandboxStatusMaxBytes {
		return SandboxStatus{}, fmt.Errorf("%w: larger than %d bytes", ErrInvalidSandboxStatus, SandboxStatusMaxBytes)
	}
	for _, b := range raw {
		if b != '\n' && (b < 0x20 || b > 0x7e) {
			return SandboxStatus{}, fmt.Errorf("%w: control or non-ASCII byte", ErrInvalidSandboxStatus)
		}
	}
	lines := strings.Split(strings.TrimSuffix(string(raw), "\n"), "\n")
	if len(lines) != 3 || lines[0] != sandboxStatusHeader {
		return SandboxStatus{}, fmt.Errorf("%w: expected the header %q and two fields", ErrInvalidSandboxStatus, sandboxStatusHeader)
	}
	resultValue, ok := strings.CutPrefix(lines[1], sandboxStatusResultKey)
	if !ok {
		return SandboxStatus{}, fmt.Errorf("%w: second line is not `result=`", ErrInvalidSandboxStatus)
	}
	result := SandboxStatusResult(resultValue)
	switch result {
	case SandboxStatusOK, SandboxStatusFallback, SandboxStatusFailed:
	default:
		return SandboxStatus{}, fmt.Errorf("%w: result outside {ok, fallback, failed}", ErrInvalidSandboxStatus)
	}
	atValue, ok := strings.CutPrefix(lines[2], sandboxStatusAtKey)
	if !ok {
		return SandboxStatus{}, fmt.Errorf("%w: third line is not `at=`", ErrInvalidSandboxStatus)
	}
	// time.Parse also takes a fractional second the layout does not name;
	// the shape check keeps the timestamp strict.
	at, err := time.Parse(sandboxStatusTimeLayout, atValue)
	if err != nil || !sandboxStatusTimeShape.MatchString(atValue) {
		return SandboxStatus{}, fmt.Errorf("%w: at is not a UTC timestamp of the form %s", ErrInvalidSandboxStatus, sandboxStatusTimeLayout)
	}
	return SandboxStatus{result: result, at: at}, nil
}
