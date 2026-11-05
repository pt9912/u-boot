package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// LH-FA-DEV-013: an entry of `devcontainer.features.brew.packages` is a
// formula or a formula of a tap.
func TestNewBrewPackage_Accepts(t *testing.T) {
	t.Parallel()
	cases := []struct {
		raw, wantTap string
	}{
		{"ripgrep", ""},
		{"node@22", ""},
		{"gcc@14", ""},
		{"c++filt", ""},
		{"a", ""},
		{"0ad", ""},
		{"python-tk@3.13", ""},
		{"pt9912/u-boot/u-boot", "pt9912/u-boot"},
		{"Homebrew/core/ripgrep", "Homebrew/core"},
		{"user.name/my_tap/formula+1", "user.name/my_tap"},
	}
	for _, tc := range cases {
		got, err := domain.NewBrewPackage(tc.raw)
		if err != nil {
			t.Errorf("NewBrewPackage(%q): unexpected error %v", tc.raw, err)
			continue
		}
		if got.String() != tc.raw || got.Tap() != tc.wantTap {
			t.Errorf("NewBrewPackage(%q) = (%q, tap %q), want (%q, tap %q)", tc.raw, got.String(), got.Tap(), tc.raw, tc.wantTap)
		}
	}
}

// LH-FA-DEV-013, boundary and negative: two or four segments, a leading `-`,
// whitespace, shell characters, empty segments and URL forms are rejected;
// the rejection wraps [domain.ErrInvalidBrewPackage].
func TestNewBrewPackage_Rejects(t *testing.T) {
	t.Parallel()
	cases := map[string]string{
		"empty":                 "",
		"two segments":          "pt9912/u-boot",
		"four segments":         "a/b/c/d",
		"leading dash":          "-ripgrep",
		"option in segment":     "a/b/--head",
		"space":                 "rip grep",
		"leading space":         " ripgrep",
		"trailing space":        "ripgrep ",
		"newline":               "ripgrep\nrm",
		"semicolon":             "ripgrep;id",
		"dollar":                "$(id)",
		"backtick":              "`id`",
		"ampersand":             "a&b",
		"pipe":                  "a|b",
		"quote":                 "a'b",
		"double quote":          `a"b`,
		"backslash":             `a\b`,
		"glob":                  "rip*",
		"empty tap segment":     "pt9912//formula",
		"empty formula segment": "pt9912/u-boot/",
		"leading slash":         "/pt9912/u-boot",
		"url":                   "https://github.com/pt9912/homebrew-u-boot",
		"at start":              "@22",
		"dot start":             ".hidden",
		"utf8":                  "ripégrep",
	}
	for name, raw := range cases {
		if _, err := domain.NewBrewPackage(raw); !errors.Is(err, domain.ErrInvalidBrewPackage) {
			t.Errorf("%s: NewBrewPackage(%q) error = %v, want a wrap of ErrInvalidBrewPackage", name, raw, err)
		}
	}
}

// The rejection names the entry so the user finds it in the list.
func TestNewBrewPackage_ErrorNamesTheEntry(t *testing.T) {
	t.Parallel()
	_, err := domain.NewBrewPackage("rip grep")
	if err == nil || !strings.Contains(err.Error(), `"rip grep"`) {
		t.Errorf("error = %v, want one that names the entry", err)
	}
}

// LH-FA-DEV-013: an entry of `devcontainer.brewTaps.allow` is `<user>/<tap>`.
func TestNewBrewTap(t *testing.T) {
	t.Parallel()
	for _, ok := range []string{"pt9912/u-boot", "Homebrew/core", "a/b", "user.name/my_tap"} {
		got, err := domain.NewBrewTap(ok)
		if err != nil || got != ok {
			t.Errorf("NewBrewTap(%q) = (%q, %v), want it accepted unchanged", ok, got, err)
		}
	}
	for _, bad := range []string{"", "pt9912", "pt9912/u-boot/u-boot", "pt9912/", "/u-boot", "pt9912/u boot", "-a/b", "a/-b", "a/b;id", " a/b"} {
		if _, err := domain.NewBrewTap(bad); !errors.Is(err, domain.ErrInvalidBrewPackage) {
			t.Errorf("NewBrewTap(%q) error = %v, want a wrap of ErrInvalidBrewPackage", bad, err)
		}
	}
}
