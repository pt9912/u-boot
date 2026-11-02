package driven

import "context"

// Fetcher abstracts one HTTP GET of a public text resource at project
// generation time. The production implementation lives in
// `internal/adapter/driven/fetch` and is the only place `net/http`
// runs in the project.
//
// Layer rules: driven port; the application layer reaches it through
// this interface so the hexagon stays free of external I/O
// (LH-FA-ARCH-003). Generated devcontainer configurations read the
// official checksum source of a pinned-toolchain entry when a project
// overrides its version with `version:` (LH-FA-DEV-003, SPEC-054);
// the default pin needs no read, because the checksums stand in the
// template beside the version.
//
// Contract of the single method:
//
//   - `url` is the absolute HTTPS URL of the resource; the adapter
//     performs a plain GET (no redirects beyond the client default,
//     no headers beyond the Go defaults) and returns the body bytes.
//   - A non-2xx status, a transport failure or a context cancellation
//     is a non-nil error; the caller classifies it (the unreadable
//     source is a fachlicher Validierungsfehler of the override, not
//     an infrastructure error — the override itself is unusable).
//   - The body is read in full and bounded by the adapter's size cap,
//     so a huge or unbounded response cannot exhaust the memory of a
//     `generate` run; a cap breach is a non-nil error.
//   - The adapter retains no state across calls; every call performs
//     a new request.
type Fetcher interface {
	Get(ctx context.Context, url string) ([]byte, error)
}
