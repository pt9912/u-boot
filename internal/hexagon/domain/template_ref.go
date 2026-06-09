package domain

import "strings"

// TemplateRefKind classifies a raw `--template` argument as either a
// built-in catalog name or a local filesystem path
// (slice-later-local-templates, ADR-0009 §Entscheidung "Lokale
// User-Templates").
type TemplateRefKind int

const (
	// TemplateRefCatalog is a built-in catalog identifier (e.g.
	// `basic`) resolved against the embedded template set.
	TemplateRefCatalog TemplateRefKind = iota

	// TemplateRefPath is a local filesystem path (e.g. `./my-template`,
	// `/abs/tpl`, `~/tpl`) resolved against the real filesystem.
	TemplateRefPath
)

// ClassifyTemplateRef decides whether ref names a built-in catalog
// template or a local filesystem path.
//
// The rule is platform-independent and deterministic: it does NOT
// consult `filepath.Separator` or touch the filesystem, so the
// Linux/macOS/Windows binaries classify a given string identically and
// there is no TOCTOU window between classification and resolution
// (slice-later-local-templates T0-(a)).
//
// ref is a path when it:
//
//   - starts with `./`, `../`, or `/`;
//   - is exactly `~` or starts with `~/`;
//   - contains a forward slash `/` or backslash `\` anywhere;
//   - looks like a Windows drive designator (`C:…`).
//
// Otherwise it is a catalog name (validated separately against the
// kebab-case shape in [TemplateMetadata.Validate]).
//
// `~`-expansion is the resolver's job, not the domain's; `~user` is
// not treated as a home alias and — having no slash — falls through to
// the catalog branch, where it harmlessly fails name lookup.
func ClassifyTemplateRef(ref string) TemplateRefKind {
	switch {
	case strings.HasPrefix(ref, "./"),
		strings.HasPrefix(ref, "../"),
		strings.HasPrefix(ref, "/"):
		return TemplateRefPath
	case ref == "~", strings.HasPrefix(ref, "~/"):
		return TemplateRefPath
	case strings.ContainsAny(ref, `/\`):
		return TemplateRefPath
	case looksLikeWindowsDrive(ref):
		return TemplateRefPath
	default:
		return TemplateRefCatalog
	}
}

// looksLikeWindowsDrive reports whether ref begins with a Windows
// drive designator: an ASCII letter followed by a colon (`C:\tpl`,
// `D:/tpl`, or the drive-relative `C:tpl`). The separator-bearing
// forms are already caught by the `/`+`\` check in
// [ClassifyTemplateRef]; this helper additionally catches the bare
// `C:tpl` form, which can never be a valid kebab-case catalog name
// (no colon) and so is safely routed to the path branch.
func looksLikeWindowsDrive(ref string) bool {
	if len(ref) < 2 || ref[1] != ':' {
		return false
	}
	c := ref[0]
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}
