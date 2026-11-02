// Package fetch is the net/http-backed implementation of the
// `port/driven.Fetcher` interface (SPEC-054). It performs the HTTP
// GET behind the read of an official checksum source of a pinned
// toolchain entry (`go-toolchain`, `golangci-lint`) when a project
// overrides the pin with `version:`; the default pin is not read,
// its checksums stand beside the version in the catalogue template.
//
// Layer rules (LH-FA-ARCH-003): driven adapter — imports
// `hexagon/port/driven` + stdlib only; the application layer is
// blocked from importing `net/http` directly by the depguard rule
// `application-no-net`, so all HTTP traffic of a `generate` run
// funnels through here.
//
// `u-boot` reads the source only for a `version:` override, never at
// another call; a failed read is classified by the caller as a
// fachlicher Validierungsfehler of the override (SPEC-054), not as
// an infrastructure error of `generate`.
package fetch

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// Response cap of one Get call: the two checksum sources are a JSON
// listing (go.dev/dl, with `include=all` a few megabytes) and a
// checksums file of a release (kilobytes); 64 MiB bound both with
// room for a listing that adds versions without breaking the cap.
const responseLimit int64 = 64 << 20

// Client is the production Fetcher adapter. Construct with [New].
type Client struct {
	http    *http.Client
	timeout time.Duration
}

// Static check: Client satisfies the Fetcher port.
var _ driven.Fetcher = (*Client)(nil)

// New constructs the adapter: a short-lived GET client with the
// given timeout per call (and no cookie jar, proxy from the
// environment, as the package default of `http.Client` keeps it).
func New(timeout time.Duration) *Client {
	return &Client{
		http:    &http.Client{Timeout: timeout},
		timeout: timeout,
	}
}

// Get performs one GET of url and returns the body bytes. A non-2xx
// status, a transport failure, a cancellation or a body over
// [responseLimit] is a non-nil error; the caller of the use case
// classifies it (SPEC-054: the unreadable source ends the override,
// not the command).
func (c *Client) Get(ctx context.Context, url string) ([]byte, error) {
	reqCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("Get %q: %w", url, err)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Get %q: %w", url, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("Get %q: unexpected status %s: %.4096s", url, resp.Status, body)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, responseLimit+1))
	if err != nil {
		return nil, fmt.Errorf("Get %q: %w", url, err)
	}
	if int64(len(body)) > responseLimit {
		return nil, fmt.Errorf("Get %q: response over the %d-byte limit", url, responseLimit)
	}
	return body, nil
}
