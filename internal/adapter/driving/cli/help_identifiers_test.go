package cli_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
)

// identifierPatterns match the identifiers of requirements,
// specification, decisions, slices and tranches. A user reading the
// help of the binary has no use for them: the text says in words what
// the identifier named.
func identifierPatterns() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(`\bLH-[A-Z]+(-[A-Z0-9]+)*-?[0-9]+[A-Z]?\b`),
		regexp.MustCompile(`\bLH-AK-[0-9]+\b`),
		regexp.MustCompile(`\b(SPEC|ADR|ARC)-[0-9]+\b`),
		regexp.MustCompile(`\bslice-[a-z0-9][a-z0-9.-]*`),
		regexp.MustCompile(`§`),
		regexp.MustCompile(`\bT[0-9]+-\(`),
	}
}

// germanWords are German words that do not belong in English help text.
var germanWords = regexp.MustCompile(`(?i)\b(voll-schema|fachlich\w*)\b`)

// announcementPatterns match announcements and history. Help text and
// messages describe what is; the list is narrow so that a description
// of the present state ("planned changes", "later calls", "no longer
// needed") stays green.
func announcementPatterns() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(`(?i)\bcoming (soon|in)\b`),
		regexp.MustCompile(`\bComing\b`),
		regexp.MustCompile(`will follow`),
		regexp.MustCompile(`\bjoin in\b`),
		regexp.MustCompile(`(?i)\bnot yet\b`),
		regexp.MustCompile(`\bNOT\b.*\byet\b`),
		regexp.MustCompile(`\bdeferred\b`),
		regexp.MustCompile(`\bplanned for\b`),
		regexp.MustCompile(`\bTODO\b`),
		regexp.MustCompile(`\bV[0-9]+\b`),
		regexp.MustCompile(`\bv[0-9]+\.[0-9]+(\.[0-9]+)?\b`),
		regexp.MustCompile(`\btoday\b`),
		regexp.MustCompile(`\bMVP\b`),
		regexp.MustCompile(`\bpreviously\b`),
		regexp.MustCompile(`\bformerly\b`),
		regexp.MustCompile(`\bused to\b`),
		regexp.MustCompile(`\bwas (removed|renamed|replaced)\b`),
		regexp.MustCompile(`deprecated since`),
		regexp.MustCompile(`\bsince v?[0-9]`),
		regexp.MustCompile(`\bchanged in\b`),
		regexp.MustCompile(`\bin v[0-9]`),
		regexp.MustCompile(`\b(früher|bisher|vorher|seit)\b`),
	}
}

// userTextFindings returns the patterns that match text, line by line.
// A line that links the manual of a release tag (`blob/v`) is exempt
// from the version pattern.
func userTextFindings(text string) []string {
	all := append(identifierPatterns(), germanWords)
	all = append(all, announcementPatterns()...)
	var found []string
	for _, line := range strings.Split(text, "\n") {
		if strings.Contains(line, "blob/v") {
			continue
		}
		for _, re := range all {
			if m := re.FindString(line); m != "" {
				found = append(found, m)
			}
		}
	}
	return found
}

func TestUserTextFindings_Patterns(t *testing.T) {
	red := []string{
		"supports only postgres (LH-FA-ADD-002)",
		"Exit codes (LH-FA-CLI-006):",
		"see LH-AK-005",
		"per ADR-0012",
		"defined in SPEC-031",
		"see slice-m8-config.md",
		"§D1 of the plan",
		"per T0-(b)",
		"voll-schema",
		"Coming in later slices:",
		"Keycloak will follow in their own V1 slices",
		"v0.3.0 does NOT auto-remove volumes",
		"volume removal is deferred (v0.3.0)",
		"today MVP supports only postgres",
	}
	for _, s := range red {
		if len(userTextFindings(s)) == 0 {
			t.Errorf("no pattern matches %q", s)
		}
	}
	green := []string{
		"Exit codes:",
		"after confirming the data is no longer needed",
		"planned changes without writing",
		"a later u-boot devcontainer up starts the same container again",
		"You will need to start a new shell",
		"User manual (German), this version: https://x/blob/v0.17.0/docs/m.md",
		"LHS is not an identifier",
	}
	for _, s := range green {
		if f := userTextFindings(s); len(f) != 0 {
			t.Errorf("%q must stay green, matched %v", s, f)
		}
	}
}

// TestHelpTreeCarriesNoIdentifiersOrHistory (seit slice-v1-hilfetexte-ohne-kennungen)
// walks the whole command
// tree and reports every help text (Short, Long, Example, option
// usage) with an identifier, a German word, an announcement or a
// history remark.
func TestHelpTreeCarriesNoIdentifiersOrHistory(t *testing.T) {
	app := cli.New("0.17.0", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	texts := app.HelpTextsForTest()
	if len(texts) < 50 {
		t.Fatalf("only %d help texts collected — tree walk broken", len(texts))
	}
	for _, h := range texts {
		for _, f := range userTextFindings(h.Text) {
			t.Errorf("%s [%s]: %q", h.Path, h.Field, f)
		}
	}
}

// TestOptionUsageCarriesNoBacktickWord pins that no option's usage text
// contains a backtick: Cobra reads the first backticked word of a usage
// text as the name of the option's value and prints it as a type name
// (`--quiet up` instead of `--quiet`).
func TestOptionUsageCarriesNoBacktickWord(t *testing.T) {
	app := cli.New("0.17.0", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	for _, h := range app.HelpTextsForTest() {
		if strings.HasPrefix(h.Field, "--") && strings.Contains(h.Text, "`") {
			t.Errorf("%s [%s]: backtick in usage text %q", h.Path, h.Field, h.Text)
		}
	}
}
