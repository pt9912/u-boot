package confirm_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driven/confirm"
)

func TestConfirmer_AnswerYes(t *testing.T) {
	t.Parallel()
	for _, ans := range []string{"y\n", "Y\n", "yes\n", "Yes\n", "YES\n"} {
		ans := ans
		t.Run(strings.TrimSpace(ans), func(t *testing.T) {
			t.Parallel()
			in := strings.NewReader(ans)
			out := &bytes.Buffer{}
			c := confirm.New(in, out)
			got, err := c.ConfirmTreatAsExisting(context.Background(), "/tmp/x", []string{"README.md", "docs"})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got {
				t.Errorf("got false, want true")
			}
		})
	}
}

func TestConfirmer_AnswerNo(t *testing.T) {
	t.Parallel()
	for _, ans := range []string{"n\n", "N\n", "no\n", "\n", "", "anything else\n"} {
		ans := ans
		t.Run(strings.TrimSpace(ans), func(t *testing.T) {
			t.Parallel()
			in := strings.NewReader(ans)
			out := &bytes.Buffer{}
			c := confirm.New(in, out)
			got, err := c.ConfirmTreatAsExisting(context.Background(), "/tmp/x", []string{"README.md"})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got {
				t.Errorf("got true, want false (default-N policy)")
			}
		})
	}
}

func TestConfirmer_PromptShowsIndicators(t *testing.T) {
	t.Parallel()
	in := strings.NewReader("n\n")
	out := &bytes.Buffer{}
	c := confirm.New(in, out)
	indicators := []string{"README.md", "docs", "scripts"}
	if _, err := c.ConfirmTreatAsExisting(context.Background(), "/tmp/x", indicators); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	prompt := out.String()
	for _, ind := range indicators {
		if !strings.Contains(prompt, ind) {
			t.Errorf("prompt does not list indicator %q: %q", ind, prompt)
		}
	}
	if !strings.Contains(prompt, "/tmp/x") {
		t.Errorf("prompt does not include baseDir: %q", prompt)
	}
	if !strings.Contains(prompt, "[y/N]") {
		t.Errorf("prompt does not show default-N hint: %q", prompt)
	}
}

func TestConfirmer_ReadErrorPropagates(t *testing.T) {
	t.Parallel()
	// Simulate an unrecoverable read error via a reader that always
	// fails. EOF (no bytes) → "no" (default); a different error →
	// propagated.
	in := &erroringReader{err: errSimulated}
	out := &bytes.Buffer{}
	c := confirm.New(in, out)
	_, err := c.ConfirmTreatAsExisting(context.Background(), "/tmp/x", []string{"README.md"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

type erroringReader struct {
	err error
}

func (r *erroringReader) Read(_ []byte) (int, error) {
	return 0, r.err
}

var errSimulated = &simulatedErr{}

type simulatedErr struct{}

func (*simulatedErr) Error() string { return "simulated read failure" }

func TestConfirmer_RemoveVolumes_AnswerYes(t *testing.T) {
	t.Parallel()
	for _, ans := range []string{"y\n", "Y\n", "yes\n", "Yes\n", "YES\n"} {
		ans := ans
		t.Run(strings.TrimSpace(ans), func(t *testing.T) {
			t.Parallel()
			in := strings.NewReader(ans)
			out := &bytes.Buffer{}
			c := confirm.New(in, out)
			got, err := c.ConfirmRemoveVolumes(context.Background(), "/tmp/proj")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got {
				t.Errorf("got false, want true")
			}
		})
	}
}

func TestConfirmer_RemoveVolumes_AnswerNo(t *testing.T) {
	t.Parallel()
	// Same default-N policy as ConfirmTreatAsExisting — empty
	// response, EOF, and anything-not-y/yes all answer no.
	for _, ans := range []string{"n\n", "N\n", "no\n", "\n", "", "garbage\n"} {
		ans := ans
		t.Run(strings.TrimSpace(ans), func(t *testing.T) {
			t.Parallel()
			in := strings.NewReader(ans)
			out := &bytes.Buffer{}
			c := confirm.New(in, out)
			got, err := c.ConfirmRemoveVolumes(context.Background(), "/tmp/proj")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got {
				t.Errorf("got true, want false (default-N for destructive op)")
			}
		})
	}
}

func TestConfirmer_RemoveVolumes_PromptShowsBaseDirAndWarning(t *testing.T) {
	t.Parallel()
	in := strings.NewReader("n\n")
	out := &bytes.Buffer{}
	c := confirm.New(in, out)
	if _, err := c.ConfirmRemoveVolumes(context.Background(), "/tmp/proj"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	prompt := out.String()
	if !strings.Contains(prompt, "/tmp/proj") {
		t.Errorf("prompt does not include baseDir: %q", prompt)
	}
	if !strings.Contains(prompt, "PERMANENTLY DELETED") {
		t.Errorf("prompt does not warn about data loss: %q", prompt)
	}
	if !strings.Contains(prompt, "[y/N]") {
		t.Errorf("prompt does not show default-N hint: %q", prompt)
	}
}

func TestConfirmer_RemoveVolumes_ReadErrorPropagates(t *testing.T) {
	t.Parallel()
	in := &erroringReader{err: errSimulated}
	out := &bytes.Buffer{}
	c := confirm.New(in, out)
	_, err := c.ConfirmRemoveVolumes(context.Background(), "/tmp/proj")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestConfirmer_AddDependency_AnswerYes(t *testing.T) {
	t.Parallel()
	for _, ans := range []string{"y\n", "Y\n", "yes\n", "Yes\n", "YES\n"} {
		ans := ans
		t.Run(strings.TrimSpace(ans), func(t *testing.T) {
			t.Parallel()
			in := strings.NewReader(ans)
			out := &bytes.Buffer{}
			c := confirm.New(in, out)
			got, err := c.ConfirmAddDependency(context.Background(), "keycloak", []string{"postgres"})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got {
				t.Errorf("got false, want true")
			}
		})
	}
}

func TestConfirmer_AddDependency_AnswerNo(t *testing.T) {
	t.Parallel()
	// Same default-N policy: empty, EOF, n/no/garbage all decline.
	for _, ans := range []string{"n\n", "N\n", "no\n", "\n", "", "garbage\n"} {
		ans := ans
		t.Run(strings.TrimSpace(ans), func(t *testing.T) {
			t.Parallel()
			in := strings.NewReader(ans)
			out := &bytes.Buffer{}
			c := confirm.New(in, out)
			got, err := c.ConfirmAddDependency(context.Background(), "keycloak", []string{"postgres"})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got {
				t.Errorf("got true, want false (default-N)")
			}
		})
	}
}

func TestConfirmer_AddDependency_PromptShowsServiceAndMissing(t *testing.T) {
	t.Parallel()
	in := strings.NewReader("n\n")
	out := &bytes.Buffer{}
	c := confirm.New(in, out)
	if _, err := c.ConfirmAddDependency(context.Background(), "keycloak", []string{"postgres", "redis"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	prompt := out.String()
	if !strings.Contains(prompt, "keycloak") {
		t.Errorf("prompt does not include service name: %q", prompt)
	}
	if !strings.Contains(prompt, "postgres") || !strings.Contains(prompt, "redis") {
		t.Errorf("prompt does not list missing add-ons: %q", prompt)
	}
	if !strings.Contains(prompt, "[y/N]") {
		t.Errorf("prompt does not show default-N hint: %q", prompt)
	}
}

func TestConfirmer_AddDependency_ReadErrorPropagates(t *testing.T) {
	t.Parallel()
	in := &erroringReader{err: errSimulated}
	out := &bytes.Buffer{}
	c := confirm.New(in, out)
	_, err := c.ConfirmAddDependency(context.Background(), "keycloak", []string{"postgres"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
