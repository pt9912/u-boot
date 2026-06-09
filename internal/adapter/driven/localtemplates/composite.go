package localtemplates

import (
	"context"
	iofs "io/fs"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// Composite is a [driven.TemplateFiles] that dispatches by reference
// shape: a catalog name goes to the embedded-catalog resolver, a
// filesystem path to the local resolver. The split is the pure
// [domain.ClassifyTemplateRef] rule — no filesystem probe — so the
// dispatch is deterministic and platform-independent
// (slice-later-local-templates T0-(a)/(a2)).
//
// Living in the adapter layer is deliberate: only here may both
// resolvers be known at once (LH-FA-ARCH-003). The CLI passes the raw
// `--template` string straight through as the template ref and does no
// dispatch of its own; `cmd/uboot` wires one Composite as the single
// [driven.TemplateFiles] the [application.TemplateInitService]
// consumes.
type Composite struct {
	catalog driven.TemplateFiles
	local   driven.TemplateFiles
}

// Static check: Composite satisfies the TemplateFiles port.
var _ driven.TemplateFiles = (*Composite)(nil)

// NewComposite wires the catalog (name-resolving) and local (path-
// resolving) resolvers behind one port. Both are mandatory; a nil
// argument is a wiring bug that would panic on first dispatch.
func NewComposite(catalog, local driven.TemplateFiles) *Composite {
	return &Composite{catalog: catalog, local: local}
}

// Open classifies ref and delegates. A [domain.TemplateRefPath] ref
// (e.g. `./tpl`, `/abs`, `~/tpl`) goes to the local resolver; every
// other ref is treated as a catalog name.
func (c *Composite) Open(ctx context.Context, ref string) (iofs.FS, error) {
	if domain.ClassifyTemplateRef(ref) == domain.TemplateRefPath {
		return c.local.Open(ctx, ref)
	}
	return c.catalog.Open(ctx, ref)
}
