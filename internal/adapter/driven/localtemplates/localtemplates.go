// Package localtemplates resolves `u-boot init --template <path>`
// against the real filesystem (ADR-0009 §Entscheidung "Lokale
// User-Templates", slice-later-local-templates). It is the filesystem-
// backed sibling of the embed.FS catalog adapter `externaltemplates`:
// both satisfy [driven.TemplateFiles], and both parse `template.yaml`
// via the shared `templateyaml` package so the apiVersion gate and
// metadata validation stay identical.
//
// This package imports neither `externaltemplates` nor
// `adapter/driven/fs`; it is a self-contained concrete adapter over
// the `os` / `io/fs` stdlib (LH-FA-ARCH-003 — adapters wire via
// `cmd/uboot`, not to each other).
package localtemplates

import (
	"context"
	"errors"
	"fmt"
	iofs "io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/pt9912/u-boot/internal/adapter/driven/templateyaml"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// Resolver is the filesystem-backed TemplateFiles adapter. It resolves
// a user-supplied directory path (with optional `~` / `~/…` home
// expansion) to an [iofs.FS] rooted at that directory, after gating on
// a present and valid `template.yaml`. The zero value is usable.
type Resolver struct{}

// Static check: Resolver satisfies the TemplateFiles port.
var _ driven.TemplateFiles = (*Resolver)(nil)

// New returns a filesystem TemplateFiles resolver.
func New() *Resolver { return &Resolver{} }

// Open resolves name as a filesystem path and returns an [iofs.FS]
// rooted at the template directory.
//
// Error classification (slice-later-local-templates T0-(d)):
//
//   - absent path / not a directory / no `template.yaml` →
//     [driven.ErrTemplateNotFound] (exit 10 — the user fixes the path).
//   - `template.yaml` present but malformed / unsupported apiVersion /
//     failing the metadata minimum → [driven.ErrTemplateInvalid]
//     (exit 10 — the user fixes the metadata).
//   - a technical stat failure (e.g. permission denied) → a plain
//     wrapped error the service maps to the technical class (exit 14).
//
// Symlink policy is NOT enforced here (T0-(e)): the resolver hands back
// the rooted, unfollowed [iofs.FS]; the application render-loop rejects
// any symlink entry during its walk. The root itself may be an
// absolute path, a `..`-bearing path, or a symlink — it is the user's
// explicit choice (T0-(f)); only paths *inside* the template are
// guarded.
//
// Relative refs (e.g. `./tpl`) resolve against the OS process working
// directory via [os.Stat] / [os.DirFS], NOT against any injected
// getwd seam the CLI uses for the project base dir. In production the
// two coincide (the CLI's getwd defaults to [os.Getwd]); a caller that
// injects a divergent getwd (a future daemon, a test) must pass an
// absolute `--template` path to stay unambiguous.
func (*Resolver) Open(ctx context.Context, name string) (iofs.FS, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if name == "" {
		return nil, fmt.Errorf("%w: empty template path", driven.ErrTemplateNotFound)
	}

	dir, err := expandHome(name)
	if err != nil {
		// Environment failure (e.g. $HOME unset for a `~/…` ref). Not a
		// user-fixable template error; surface a plain wrapped error so
		// the service routes it to the technical class.
		return nil, fmt.Errorf("resolve template path %q: %w", name, err)
	}

	info, statErr := os.Stat(dir)
	switch {
	case errors.Is(statErr, iofs.ErrNotExist):
		return nil, fmt.Errorf("%w: %q", driven.ErrTemplateNotFound, name)
	case statErr != nil:
		return nil, fmt.Errorf("stat template path %q: %w", name, statErr)
	case !info.IsDir():
		return nil, fmt.Errorf("%w: %q is not a directory", driven.ErrTemplateNotFound, name)
	}

	root := os.DirFS(dir)
	if _, err := templateyaml.Read(root, "."); err != nil {
		// templateyaml.Read returns a wrapped iofs.ErrNotExist when
		// template.yaml is absent, and domain.ErrInvalidTemplate (or a
		// plain parse error) when it is present but malformed. Split the
		// two so a missing file maps to not-found and a bad file to
		// invalid — both exit 10, but distinct user guidance.
		if errors.Is(err, iofs.ErrNotExist) {
			return nil, fmt.Errorf("%w: %q has no template.yaml", driven.ErrTemplateNotFound, name)
		}
		return nil, fmt.Errorf("%w: %q: %w", driven.ErrTemplateInvalid, name, err)
	}
	return root, nil
}

// expandHome expands a leading `~` or `~/…` to the user's home
// directory via [os.UserHomeDir]. A bare `~user` form is NOT expanded
// (ADR-0009 / T0-(a)) — it has no `/`, so the caller's classifier
// never routes it here anyway; if it did, it would fall through
// unchanged and fail the directory stat. All other paths pass through
// untouched.
func expandHome(p string) (string, error) {
	if p != "~" && !strings.HasPrefix(p, "~/") {
		return p, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if p == "~" {
		return home, nil
	}
	return filepath.Join(home, p[len("~/"):]), nil
}
