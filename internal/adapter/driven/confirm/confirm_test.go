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

const confirmRemoveID = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

// SPEC-058, answers: y and yes in any case agree; n, a longer word, the empty
// line and the end of the input decline.
func TestConfirmer_RemoveDevcontainer_Answers(t *testing.T) {
	t.Parallel()
	cases := map[string]bool{
		"y\n": true, "Y\n": true, "yes\n": true, "YES\n": true, " yes \n": true,
		"n\n": false, "no\n": false, "yep\n": false, "\n": false, "": false, "ja\n": false,
	}
	for answer, want := range cases {
		t.Run(strings.TrimSpace(answer), func(t *testing.T) {
			t.Parallel()
			c := confirm.New(strings.NewReader(answer), &bytes.Buffer{})
			got, err := c.ConfirmRemoveDevcontainer(context.Background(), confirmRemoveID, []string{"v1"})
			if err != nil || got != want {
				t.Errorf("answer %q: got %t, err %v, want %t", answer, got, err, want)
			}
		})
	}
}

// SPEC-058, prompt: the first twelve characters of the container ID, every
// volume on a line of its own, the warning and `[y/N]` with the default N; the
// container without a volume says so and has no warning about volumes.
func TestConfirmer_RemoveDevcontainer_Prompt(t *testing.T) {
	t.Parallel()
	out := &bytes.Buffer{}
	c := confirm.New(strings.NewReader("n\n"), out)
	if _, err := c.ConfirmRemoveDevcontainer(context.Background(), confirmRemoveID, []string{"demo-docker-x", "demo-workspace-x"}); err != nil {
		t.Fatal(err)
	}
	want := "About to remove the devcontainer container 0123456789ab of this project and its volumes:\n" +
		"  - demo-docker-x\n  - demo-workspace-x\n" +
		"Data in these volumes will be PERMANENTLY DELETED. Proceed? [y/N] "
	if out.String() != want {
		t.Errorf("prompt:\n%q\nwant\n%q", out.String(), want)
	}
	out.Reset()
	c = confirm.New(strings.NewReader("n\n"), out)
	if _, err := c.ConfirmRemoveDevcontainer(context.Background(), confirmRemoveID, nil); err != nil {
		t.Fatal(err)
	}
	want = "About to remove the devcontainer container 0123456789ab of this project; it holds no volume.\nProceed? [y/N] "
	if out.String() != want {
		t.Errorf("prompt without volumes:\n%q\nwant\n%q", out.String(), want)
	}
}

// SPEC-058: an answer that cannot be read is an error, not a "no".
func TestConfirmer_RemoveDevcontainer_ReadErrorPropagates(t *testing.T) {
	t.Parallel()
	c := confirm.New(&erroringReader{err: errSimulated}, &bytes.Buffer{})
	if _, err := c.ConfirmRemoveDevcontainer(context.Background(), confirmRemoveID, nil); err == nil {
		t.Fatal("expected error, got nil")
	}
}
