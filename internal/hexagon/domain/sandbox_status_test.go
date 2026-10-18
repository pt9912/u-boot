package domain_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

const validStatusAt = "2026-10-17T18:52:54Z"

func statusDoc(result string) string {
	return "u-boot-status 1\nresult=" + result + "\nat=" + validStatusAt + "\n"
}

// LH-FA-DEV-007: the closed status format of the container start scripts.
func TestParseSandboxStatus_Valid(t *testing.T) {
	t.Parallel()
	cases := map[string]domain.SandboxStatusResult{
		"ok":       domain.SandboxStatusOK,
		"fallback": domain.SandboxStatusFallback,
		"failed":   domain.SandboxStatusFailed,
	}
	for raw, want := range cases {
		for _, doc := range []string{statusDoc(raw), strings.TrimSuffix(statusDoc(raw), "\n")} {
			got, err := domain.ParseSandboxStatus([]byte(doc))
			if err != nil {
				t.Fatalf("result %q: %v", raw, err)
			}
			if got.Result() != want {
				t.Errorf("result = %q, want %q", got.Result(), want)
			}
			if wantAt := time.Date(2026, 10, 17, 18, 52, 54, 0, time.UTC); !got.At().Equal(wantAt) {
				t.Errorf("at = %v, want %v", got.At(), wantAt)
			}
		}
	}
}

func TestParseSandboxStatus_Invalid(t *testing.T) {
	t.Parallel()
	exactLimit := "u-boot-status 1\nresult=ok\nat=" + validStatusAt + "\n"
	cases := map[string]string{
		"empty":                    "",
		"wrong version":            "u-boot-status 2\nresult=ok\nat=" + validStatusAt + "\n",
		"missing header":           "result=ok\nat=" + validStatusAt + "\n",
		"unknown result":           statusDoc("degraded"),
		"upper-case result":        statusDoc("OK"),
		"empty result":             statusDoc(""),
		"free text line":           statusDoc("ok") + "note=everything is fine\n",
		"free text instead of at":  "u-boot-status 1\nresult=ok\nall good\n",
		"swapped fields":           "u-boot-status 1\nat=" + validStatusAt + "\nresult=ok\n",
		"two lines":                "u-boot-status 1\nresult=ok\n",
		"blank line inside":        "u-boot-status 1\n\nresult=ok\nat=" + validStatusAt + "\n",
		"two trailing newlines":    statusDoc("ok") + "\n",
		"carriage return":          "u-boot-status 1\r\nresult=ok\r\nat=" + validStatusAt + "\r\n",
		"escape sequence":          "u-boot-status 1\nresult=ok\nat=" + validStatusAt + "\x1b[31m\n",
		"nul byte":                 "u-boot-status 1\nresult=ok\x00\nat=" + validStatusAt + "\n",
		"non-ascii":                "u-boot-status 1\nresult=ok\nat=2026-10-17T18:52:54Zä\n",
		"timestamp with offset":    "u-boot-status 1\nresult=ok\nat=2026-10-17T18:52:54+02:00\n",
		"timestamp fraction":       "u-boot-status 1\nresult=ok\nat=2026-10-17T18:52:54.5Z\n",
		"timestamp without Z":      "u-boot-status 1\nresult=ok\nat=2026-10-17T18:52:54\n",
		"timestamp date only":      "u-boot-status 1\nresult=ok\nat=2026-10-17\n",
		"timestamp out of range":   "u-boot-status 1\nresult=ok\nat=2026-13-17T18:52:54Z\n",
		"timestamp short fields":   "u-boot-status 1\nresult=ok\nat=2026-1-7T8:2:4Z\n",
		"over the size limit":      exactLimit + strings.Repeat("x", domain.SandboxStatusMaxBytes),
		"huge":                     strings.Repeat("a", 1<<20),
		"result with trailing gap": "u-boot-status 1\nresult=ok \nat=" + validStatusAt + "\n",
	}
	for name, doc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if _, err := domain.ParseSandboxStatus([]byte(doc)); !errors.Is(err, domain.ErrInvalidSandboxStatus) {
				t.Errorf("err = %v, want ErrInvalidSandboxStatus", err)
			}
		})
	}
}

// The size limit is exact: 512 bytes parse, 513 do not.
func TestParseSandboxStatus_SizeBoundary(t *testing.T) {
	t.Parallel()
	base := statusDoc("ok")
	// Pad the result value to reach the boundary; the padded document is
	// invalid for its value, so the boundary is read from the error text.
	atLimit := strings.Replace(base, "result=ok", "result=ok"+strings.Repeat("x", domain.SandboxStatusMaxBytes-len(base)), 1)
	if len(atLimit) != domain.SandboxStatusMaxBytes {
		t.Fatalf("setup: len = %d", len(atLimit))
	}
	_, errAt := domain.ParseSandboxStatus([]byte(atLimit))
	_, errOver := domain.ParseSandboxStatus([]byte(atLimit + "x"))
	if errAt == nil || strings.Contains(errAt.Error(), "larger than") {
		t.Errorf("512 bytes: err = %v, want a value error, not the size error", errAt)
	}
	if errOver == nil || !strings.Contains(errOver.Error(), "larger than") {
		t.Errorf("513 bytes: err = %v, want the size error", errOver)
	}
}
