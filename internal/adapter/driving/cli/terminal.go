package cli

import (
	"io"
	"os"
)

// isTerminal reports whether r is a terminal: a file that is a character
// device and not the null device. `devcontainer shell` asks it of the standard
// input and allocates a terminal for the session when it is one.
//
// Limit: the check uses the file mode only, so a character device that is no
// terminal (`/dev/zero`) counts as one; `docker exec -t` then refuses it and
// the launcher ends the command with exit code 12. A reader that is no file
// (a buffer, a pipe wrapped in another reader) is no terminal.
func isTerminal(r io.Reader) bool {
	f, ok := r.(*os.File)
	if !ok {
		return false
	}
	info, err := f.Stat()
	if err != nil || info.Mode()&os.ModeCharDevice == 0 {
		return false
	}
	null, err := os.Stat(os.DevNull)
	return err != nil || !os.SameFile(info, null)
}
