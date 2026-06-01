package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

func TestNewTemplatePath(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		raw      string
		want     string // canonical form when accepted; "" when rejected
		wantErr  bool
		wantFrag string // expected fragment in the error message
	}{
		// Accepted shapes.
		{name: "simple file", raw: "u-boot.yaml", want: "u-boot.yaml"},
		{name: "nested file", raw: "docker/Dockerfile", want: "docker/Dockerfile"},
		{name: "dot-prefixed file", raw: ".gitignore", want: ".gitignore"},
		{name: "leading ./ collapses", raw: "./compose.yaml", want: "compose.yaml"},
		{name: "double slash collapses", raw: "docker//compose.yaml", want: "docker/compose.yaml"},
		{name: "trailing slash collapses", raw: "docker/", want: "docker"},

		// Rejected shapes.
		{
			name: "empty", raw: "",
			wantErr: true, wantFrag: "empty path",
		},
		{
			name: "absolute unix", raw: "/etc/passwd",
			wantErr: true, wantFrag: "absolute",
		},
		{
			name: "absolute windows backslash", raw: `\etc\passwd`,
			wantErr: true, wantFrag: "backslash",
		},
		{
			name: "backslash + parent-dir bypass (review-followup F2)", raw: `docker\..\..\etc\passwd`,
			wantErr: true, wantFrag: "backslash",
		},
		{
			name: "NUL byte rejected (review-followup F4)", raw: "foo\x00bar",
			wantErr: true, wantFrag: "NUL",
		},
		{
			name: "windows drive letter", raw: "C:foo",
			wantErr: true, wantFrag: "drive letter",
		},
		{
			name: "windows drive letter with backslash (review-followup F2 catches backslash first)", raw: `D:\bar`,
			wantErr: true, wantFrag: "backslash",
		},
		{
			name: "leading parent dir", raw: "../escape",
			wantErr: true, wantFrag: "..",
		},
		{
			name: "parent dir mid-path (would clean-normalize)", raw: "foo/../bar",
			wantErr: true, wantFrag: "..",
		},
		{
			name: "nested parent dir", raw: "docker/../../../etc/passwd",
			wantErr: true, wantFrag: "..",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			p, err := domain.NewTemplatePath(tc.raw)
			if !tc.wantErr {
				if err != nil {
					t.Fatalf("NewTemplatePath(%q) = err %v, want nil", tc.raw, err)
				}
				if got := p.String(); got != tc.want {
					t.Errorf("canonical form = %q, want %q", got, tc.want)
				}
				return
			}
			if err == nil {
				t.Fatalf("NewTemplatePath(%q) = nil error, want error", tc.raw)
			}
			if !errors.Is(err, domain.ErrInvalidTemplatePath) {
				t.Errorf("err = %v, want wrap of domain.ErrInvalidTemplatePath", err)
			}
			if tc.wantFrag != "" && !strings.Contains(err.Error(), tc.wantFrag) {
				t.Errorf("err.Error() = %q, missing fragment %q", err.Error(), tc.wantFrag)
			}
		})
	}
}

func TestTemplatePath_StringRoundTrips(t *testing.T) {
	t.Parallel()
	// `NewTemplatePath(p.String())` returns an equal value for every
	// accepted input — the cleaning happens once in the first
	// constructor, the second call is a no-op on the cleaned form.
	inputs := []string{
		"u-boot.yaml",
		"docker/Dockerfile",
		"./compose.yaml",
		"docker//compose.yaml",
	}
	for _, in := range inputs {
		t.Run(in, func(t *testing.T) {
			t.Parallel()
			first, err := domain.NewTemplatePath(in)
			if err != nil {
				t.Fatalf("first NewTemplatePath: %v", err)
			}
			second, err := domain.NewTemplatePath(first.String())
			if err != nil {
				t.Fatalf("second NewTemplatePath: %v", err)
			}
			if first.String() != second.String() {
				t.Errorf("round-trip diverged: %q → %q", first.String(), second.String())
			}
		})
	}
}
