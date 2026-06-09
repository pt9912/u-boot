package domain_test

import (
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

func TestClassifyTemplateRef(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		ref  string
		want domain.TemplateRefKind
	}{
		// Catalog names — kebab-case identifiers, no path markers.
		{"plain name", "basic", domain.TemplateRefCatalog},
		{"kebab name", "micronaut-sveltekit", domain.TemplateRefCatalog},
		{"single char", "a", domain.TemplateRefCatalog},
		{"digits", "svc123", domain.TemplateRefCatalog},
		{"tilde-user not a home alias", "~user", domain.TemplateRefCatalog},
		{"empty stays catalog", "", domain.TemplateRefCatalog},

		// Explicit relative / absolute prefixes.
		{"dot-slash", "./my-template", domain.TemplateRefPath},
		{"dot-dot-slash", "../tpl", domain.TemplateRefPath},
		{"absolute", "/abs/tpl", domain.TemplateRefPath},

		// Home expansion forms.
		{"bare tilde", "~", domain.TemplateRefPath},
		{"tilde-slash", "~/tpl", domain.TemplateRefPath},

		// Embedded separators anywhere.
		{"nested forward slash", "foo/bar", domain.TemplateRefPath},
		{"backslash", `foo\bar`, domain.TemplateRefPath},

		// Windows drive designators (platform-independent classification).
		{"drive backslash", `C:\tpl`, domain.TemplateRefPath},
		{"drive forward slash", "D:/tpl", domain.TemplateRefPath},
		{"drive relative bare", "C:tpl", domain.TemplateRefPath},
		{"lowercase drive", "c:tpl", domain.TemplateRefPath},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := domain.ClassifyTemplateRef(tc.ref); got != tc.want {
				t.Errorf("ClassifyTemplateRef(%q) = %v, want %v", tc.ref, got, tc.want)
			}
		})
	}
}
