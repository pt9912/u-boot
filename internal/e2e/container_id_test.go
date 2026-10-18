package e2e_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

var containerIDPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)

// containerIDFromRunOutput returns the container ID that `docker run -d`
// wrote to stdout: the last non-empty line, which must be a 64-digit
// lowercase hex string. Any other output is an error naming the whole output.
func containerIDFromRunOutput(stdout string) (string, error) {
	lines := strings.Split(strings.ReplaceAll(stdout, "\r\n", "\n"), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		if !containerIDPattern.MatchString(line) {
			break
		}
		return line, nil
	}
	return "", fmt.Errorf("no container ID in the output of docker run -d: %q", stdout)
}

func TestContainerIDFromRunOutput(t *testing.T) {
	id := strings.Repeat("ab12", 16)
	cases := []struct {
		name    string
		out     string
		want    string
		wantErr bool
	}{
		{"only the ID", id + "\n", id, false},
		{"ID without newline", id, id, false},
		{"pull lines then ID", "Unable to find image 'alpine:3' locally\n3: Pulling from library/alpine\nStatus: Downloaded newer image for alpine:3\n" + id + "\n", id, false},
		{"CRLF and trailing blank lines", "pulling\r\n" + id + "\r\n\r\n", id, false},
		{"spaces and tab around the ID", "pulling\n  \t" + id + " \t\n", id, false},
		{"CRLF line ending only", id + "\r\n", id, false},
		{"two IDs: the last wins", strings.Repeat("cd34", 16) + "\n" + id + "\n", id, false},
		{"empty", "", "", true},
		{"blank only", "\n \n", "", true},
		{"no ID", "Unable to find image 'alpine:3' locally\n", "", true},
		{"short ID", id[:12] + "\n", "", true},
		{"ID followed by an error line", id + "\ndocker: Error response from daemon\n", "", true},
		{"upper-case hex", strings.ToUpper(id) + "\n", "", true},
		{"ID with a prefix on the line", "x" + id + "\n", "", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := containerIDFromRunOutput(c.out)
			if (err != nil) != c.wantErr || got != c.want {
				t.Errorf("containerIDFromRunOutput(%q) = %q, %v; want %q, error=%v", c.out, got, err, c.want, c.wantErr)
			}
		})
	}
}
