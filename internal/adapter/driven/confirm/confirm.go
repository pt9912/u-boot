// Package confirm is the CLI-side implementation of
// [driven.Confirmer]. It prompts the user on stdout/stderr and reads
// the answer from stdin, returning a yes/no decision to the
// application layer (LH-FA-INIT-004 soft-existing-detection).
//
// Layer rules (LH-FA-ARCH-003): driven adapter, may import
// `hexagon/port/driven` and the Go standard library; must not import
// `hexagon/application` or `adapter/driving`.
package confirm

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
)

// Confirmer prompts the user via the configured input/output streams.
// One Confirmer instance is intended per process; the wiring layer
// (cmd/uboot) constructs it with os.Stdin / os.Stderr. The prompt
// renders to `out` (stderr by convention, so stdout stays usable for
// machine-readable summaries) and reads one line of stdin for the
// answer.
type Confirmer struct {
	in  io.Reader
	out io.Writer
}

// New constructs a Confirmer. Pass real os.Stdin / os.Stderr in
// production; tests pass `*bytes.Buffer` (or `strings.NewReader`)
// pairs so the prompt can be inspected and the answer pre-seeded.
func New(in io.Reader, out io.Writer) *Confirmer {
	return &Confirmer{in: in, out: out}
}

// ConfirmTreatAsExisting renders the LH-FA-INIT-004 soft-existing-
// detection prompt with the matched indicators in a deterministic
// list, then reads one line from stdin. Accepts `y` / `yes` (case
// insensitive) as confirmation; anything else (including EOF / empty
// line) is a "no". A read error other than EOF is returned to the
// caller.
//
// The default-shown answer in the prompt is `N` (capitalized), so an
// empty response defaults to "no, proceed with fresh init"; this is
// the safer default because confirming "yes" aborts the use case.
func (c *Confirmer) ConfirmTreatAsExisting(_ context.Context, baseDir string, indicators []string) (bool, error) {
	fmt.Fprintf(c.out, "Detected %d project-structure element(s) in %s:\n", len(indicators), baseDir)
	for _, ind := range indicators {
		fmt.Fprintf(c.out, "  - %s\n", ind)
	}
	fmt.Fprint(c.out, "Treat as an existing u-boot project? Re-init requires --backup or --force. [y/N] ")

	scanner := bufio.NewScanner(c.in)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, fmt.Errorf("read confirmation: %w", err)
		}
		// EOF without input → take the default (N).
		return false, nil
	}
	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	return answer == "y" || answer == "yes", nil
}

// ConfirmRemoveVolumes renders the LH-FA-CLI-005A §254 destructive-
// confirmation prompt for `u-boot down --volumes`. Defaults to `N`
// because confirming "yes" deletes named Compose volumes (typically
// persistent service data — Postgres tables, Redis snapshots).
//
// Mirrors [ConfirmTreatAsExisting]'s parsing rules: accepts `y` /
// `yes` (case-insensitive) as confirmation; anything else
// (including EOF, empty line, "no") is "no". A read error other
// than EOF surfaces to the caller.
func (c *Confirmer) ConfirmRemoveVolumes(_ context.Context, baseDir string) (bool, error) {
	fmt.Fprintf(c.out, "About to remove all named Compose volumes in %s.\n", baseDir)
	fmt.Fprint(c.out, "Data in these volumes will be PERMANENTLY DELETED. Proceed? [y/N] ")

	scanner := bufio.NewScanner(c.in)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, fmt.Errorf("read confirmation: %w", err)
		}
		// EOF without input → take the default (N).
		return false, nil
	}
	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	return answer == "y" || answer == "yes", nil
}

// ConfirmAddDependency renders the LH-FA-ADD-006 interactive prompt
// when `u-boot add <svc>` would auto-install missing add-on
// dependencies. Defaults to `N` — pulling in extra services is a
// non-trivial side effect, so the safer-default is "no, abort".
//
// Mirrors [ConfirmTreatAsExisting]'s parsing rules: accepts `y` /
// `yes` (case-insensitive) as confirmation; anything else
// (including EOF, empty line, "no") is "no". A read error other
// than EOF surfaces to the caller.
func (c *Confirmer) ConfirmAddDependency(_ context.Context, svc string, missing []string) (bool, error) {
	fmt.Fprintf(c.out, "Service %q requires the following missing add-ons: %s.\n", svc, strings.Join(missing, ", "))
	fmt.Fprint(c.out, "Install them now? [y/N] ")

	scanner := bufio.NewScanner(c.in)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, fmt.Errorf("read confirmation: %w", err)
		}
		return false, nil
	}
	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	return answer == "y" || answer == "yes", nil
}
