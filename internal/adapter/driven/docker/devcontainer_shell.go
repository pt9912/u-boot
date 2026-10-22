package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

const (
	// shellHeadLimit bounds the head of standard error, and without a
	// terminal of standard output, that the launcher keeps to recognise an
	// error of the docker client.
	shellHeadLimit = 4096

	// Messages of the docker client (Docker 29.9.2) that end `docker exec`
	// before the session begins; the first non-empty line of the stream
	// carries them.
	clientDaemonErrorPrefix = "Error response from daemon: "
	clientEngineDownPrefix  = "Cannot connect to the Docker daemon at "
	clientNoTerminalText    = "cannot attach stdin to a TTY-enabled container because stdin is not a terminal"
	clientOCIFailurePrefix  = "OCI runtime exec failed: "
)

// Shell implements [driven.DevcontainerLauncher.Shell].
func (l *Launcher) Shell(ctx context.Context, opts driven.DevcontainerShellOptions) error {
	if err := validateShellOptions(opts); err != nil {
		return rejectedInput(err)
	}
	envNames, err := l.passedEnvNames(opts.EnvNames)
	if err != nil {
		return err
	}
	if err := l.preflight(ctx); err != nil {
		return err
	}
	shell, err := l.probeShell(ctx, opts)
	if err != nil {
		return err
	}
	return l.runSession(ctx, shell, envNames, opts)
}

// validateShellOptions applies the checks of the port to the container ID
// (like [Launcher.Stop]), the user and the workspace.
func validateShellOptions(opts driven.DevcontainerShellOptions) error {
	if !containerIDPattern.MatchString(opts.ContainerID) {
		return fmt.Errorf("docker exec: %q is not a container ID", opts.ContainerID)
	}
	if err := driven.ValidateDevcontainerShellUser(opts.User); err != nil {
		return err
	}
	return driven.ValidateDevcontainerWorkspace(opts.Workspace)
}

// probeArgs is the argument list of the probe of one shell: no `-i`, no
// `-t`, no variables; the shell runs `true` as the user in the workspace.
func probeArgs(opts driven.DevcontainerShellOptions, shell string) []string {
	return []string{"exec", "-u", opts.User, "-w", opts.Workspace, opts.ContainerID, shell, "-c", "true"}
}

// shellArgs is the argument list of the session: standard input open, `-t`
// only with a terminal, user and workspace of the caller, the passed
// variables by name only, then the shell without an argument. No helper
// container, no mount, no privilege.
func shellArgs(opts driven.DevcontainerShellOptions, envNames []string, shell string) []string {
	args := []string{"exec", "-i"}
	if opts.Terminal {
		args = append(args, "-t")
	}
	args = append(args, "-u", opts.User, "-w", opts.Workspace)
	for _, name := range envNames {
		args = append(args, "-e", name)
	}
	return append(args, opts.ContainerID, shell)
}

// probeShell returns the first of `bash` and `sh` that runs as the user in
// the workspace. When none does, the error carries the end of the output of
// each attempt.
func (l *Launcher) probeShell(ctx context.Context, opts driven.DevcontainerShellOptions) (string, error) {
	var attempts []string
	for _, shell := range []string{"bash", "sh"} {
		out, err := exec.CommandContext(ctx, l.binary, probeArgs(opts, shell)...).CombinedOutput()
		if err == nil {
			return shell, nil
		}
		if ctxErr := ctx.Err(); ctxErr != nil {
			return "", ctxErr
		}
		detail := tailOf(out)
		if detail == "" {
			detail = err.Error()
		}
		attempts = append(attempts, shell+": "+detail)
	}
	return "", fmt.Errorf("devcontainer shell: no shell runs as user %q in %q of container %s (%s): %w",
		opts.User, opts.Workspace, opts.ContainerID, strings.Join(attempts, "; "), driven.ErrDevcontainerRuntime)
}

// runSession runs `docker exec` with the streams of the caller and returns
// when the client ends. Without a terminal a cancelled context interrupts the
// client; with a terminal the context is not observed, because the client
// owns the terminal. The exit status of the client is the exit status of the
// shell and is no error, except where the client reports an error before the
// session began ([Launcher.clientErrorBeforeSession]).
func (l *Launcher) runSession(ctx context.Context, shell string, envNames []string, opts driven.DevcontainerShellOptions) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}
	args := shellArgs(opts, envNames, shell)
	var cmd *exec.Cmd
	if opts.Terminal {
		cmd = exec.CommandContext(context.WithoutCancel(ctx), l.binary, args...)
	} else {
		cmd = exec.CommandContext(ctx, l.binary, args...)
		cmd.Cancel = func() error { return cmd.Process.Signal(os.Interrupt) }
	}
	cmd.WaitDelay = l.abortGrace
	cmd.Stdin = opts.Stdin
	stdout, stderr := &headBuffer{limit: shellHeadLimit}, &headBuffer{limit: shellHeadLimit}
	if opts.Terminal {
		// The client finds the terminal on its standard output; a pipe in
		// between would take it away.
		cmd.Stdout = opts.Stdout
	} else {
		cmd.Stdout = &teeWriter{dst: opts.Stdout, head: stdout}
	}
	cmd.Stderr = &teeWriter{dst: opts.Stderr, head: stderr}

	runErr := cmd.Run()
	if !opts.Terminal {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
	}
	exit := 0
	if runErr != nil && !errors.Is(runErr, exec.ErrWaitDelay) {
		var exitErr *exec.ExitError
		if !errors.As(runErr, &exitErr) {
			return fmt.Errorf("devcontainer shell: docker exec did not run (%s): %w", runErr.Error(), driven.ErrDevcontainerRuntime)
		}
		exit = exitErr.ExitCode()
	}
	if text, ok := l.clientErrorBeforeSession(ctx, opts, exit, stdout, stderr); ok {
		return fmt.Errorf("devcontainer shell: docker exec ended before the session began (exit status %d: %s): %w",
			exit, text, driven.ErrDevcontainerRuntime)
	}
	return nil
}

// clientErrorBeforeSession reports the message of the docker client when
// `docker exec` ended with an error of the client before the session began.
// The exit status of the shell cannot tell the two apart (status 1 is also
// the status of a shell whose last command failed); the message of the
// client can, where the client is the only writer to the stream
// ([Launcher.clientMessageOnStderr], [ociFailureOnStdout]).
func (l *Launcher) clientErrorBeforeSession(ctx context.Context, opts driven.DevcontainerShellOptions, exit int, stdout, stderr *headBuffer) (string, bool) {
	if exit == 0 {
		return "", false
	}
	if text, ok := l.clientMessageOnStderr(ctx, opts, stdout, stderr); ok {
		return text, true
	}
	return ociFailureOnStdout(opts, exit, stdout, stderr)
}

// clientMessageOnStderr recognises the messages the client prints on
// standard error:
//
//   - With a terminal, the shell writes through the terminal, so standard
//     error holds client messages only; without one, standard output must be
//     empty and the first line of standard error must name this container or
//     this user (a `docker` command inside the session prints the same
//     prefix about other containers).
//   - The engine message is taken only if the engine is not reachable now,
//     because a `docker` command inside the session prints the same text.
//   - A request for a terminal on standard input that is none.
func (l *Launcher) clientMessageOnStderr(ctx context.Context, opts driven.DevcontainerShellOptions, stdout, stderr *headBuffer) (string, bool) {
	if !opts.Terminal && stdout.total != 0 {
		return "", false
	}
	first := stderr.firstLine()
	switch {
	case strings.HasPrefix(first, clientDaemonErrorPrefix):
		named := strings.Contains(first, opts.ContainerID) || strings.Contains(first, "unable to find user "+opts.User)
		return first, named
	case strings.HasPrefix(first, clientEngineDownPrefix):
		return first, l.engineUnreachable(ctx)
	case first == clientNoTerminalText:
		return first, opts.Terminal
	}
	return "", false
}

// ociFailureOnStdout recognises `OCI runtime exec failed`, which the client
// delivers on the stream of the session with status 126 or 127 (a shell has
// these statuses too). Without a terminal the line is the only output; with a
// terminal the client owns standard output and the launcher does not see it.
func ociFailureOnStdout(opts driven.DevcontainerShellOptions, exit int, stdout, stderr *headBuffer) (string, bool) {
	line := stdout.firstLine()
	found := !opts.Terminal && (exit == 126 || exit == 127) && stderr.total == 0 && strings.HasPrefix(line, clientOCIFailurePrefix)
	return line, found
}

// engineUnreachable asks the engine again, independent of the cancellation of
// ctx.
func (l *Launcher) engineUnreachable(ctx context.Context) bool {
	askCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), inspectorTimeout)
	defer cancel()
	return l.probe.Info(askCtx) != nil
}

// headBuffer keeps the first limit bytes written to it and counts all
// written bytes.
type headBuffer struct {
	data  []byte
	limit int
	total int
}

func (b *headBuffer) add(p []byte) {
	b.total += len(p)
	if room := b.limit - len(b.data); room > 0 {
		b.data = append(b.data, p[:min(room, len(p))]...)
	}
}

// firstLine is the first non-empty line of the kept head, trimmed.
func (b *headBuffer) firstLine() string {
	for _, line := range strings.Split(string(b.data), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			return line
		}
	}
	return ""
}

// teeWriter records what passes in head and forwards it to dst; a nil dst
// discards.
type teeWriter struct {
	dst  io.Writer
	head *headBuffer
}

func (w *teeWriter) Write(p []byte) (int, error) {
	w.head.add(p)
	if w.dst == nil {
		return len(p), nil
	}
	return w.dst.Write(p)
}
