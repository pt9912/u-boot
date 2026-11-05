package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ErrInvalidBrewPackage signals that a string in
// `devcontainer.features.brew.packages` or in `devcontainer.brewTaps.allow`
// does not satisfy the form of a Homebrew formula or tap (LH-FA-DEV-013).
// It lives in the domain layer so the application layer wraps it into the
// exit-code-10 sentinels of the read and write paths.
var ErrInvalidBrewPackage = errors.New("invalid homebrew package")

// brewSegmentPattern is one segment of a formula entry or of a tap: it starts
// with a letter or a digit, so it never begins with `-` (an option of
// `brew install`), and carries no whitespace and no shell character. `@` is
// part of formula names such as `node@22`; `u-boot` offers no version pin
// through it.
var brewSegmentPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9@._+-]*$`)

// BrewPackage is a validated entry of `devcontainer.features.brew.packages`:
// a formula (`<formel>`) or a formula of a tap (`<user>/<tap>/<formel>`).
// The zero value is invalid; use [NewBrewPackage] to construct.
type BrewPackage struct {
	raw string
	tap string
}

// NewBrewPackage validates one entry. Two segments, more than three segments
// and a segment outside [brewSegmentPattern] are rejected; the entry is taken
// byte-exact, so surrounding whitespace is a violation and not trimmed.
func NewBrewPackage(raw string) (BrewPackage, error) {
	segments := strings.Split(raw, "/")
	switch len(segments) {
	case 1, 3:
	default:
		return BrewPackage{}, fmt.Errorf("%w: %q has %d segments; allowed: <formula> or <user>/<tap>/<formula>",
			ErrInvalidBrewPackage, raw, len(segments))
	}
	for _, s := range segments {
		if !brewSegmentPattern.MatchString(s) {
			return BrewPackage{}, fmt.Errorf("%w: %q has the segment %q outside %s (letters, digits, `@ . _ + -`; no whitespace, no leading `-`)",
				ErrInvalidBrewPackage, raw, s, brewSegmentPattern)
		}
	}
	p := BrewPackage{raw: raw}
	if len(segments) == 3 {
		p.tap = segments[0] + "/" + segments[1]
	}
	return p, nil
}

// String returns the entry as written.
func (p BrewPackage) String() string { return p.raw }

// Tap returns `<user>/<tap>` for a formula of a tap and the empty string for a
// formula of Homebrew itself.
func (p BrewPackage) Tap() string { return p.tap }

// NewBrewTap validates one entry of `devcontainer.brewTaps.allow`:
// `<user>/<tap>`, each segment in the character set of [NewBrewPackage]. The
// entry is taken byte-exact.
func NewBrewTap(raw string) (string, error) {
	segments := strings.Split(raw, "/")
	if len(segments) != 2 {
		return "", fmt.Errorf("%w: tap %q has %d segments; allowed: <user>/<tap>",
			ErrInvalidBrewPackage, raw, len(segments))
	}
	for _, s := range segments {
		if !brewSegmentPattern.MatchString(s) {
			return "", fmt.Errorf("%w: tap %q has the segment %q outside %s (letters, digits, `@ . _ + -`; no whitespace, no leading `-`)",
				ErrInvalidBrewPackage, raw, s, brewSegmentPattern)
		}
	}
	return raw, nil
}
