// Command u-boot is the developer environment bootloader for Docker-
// based projects. See spec/lastenheft.md for the full functional
// specification.
//
// This is the wiring layer (LH-FA-ARCH-002, LH-FA-BUILD-009): the
// only place that imports both `internal/hexagon/application` and
// `internal/adapter/driven/*`. The CLI adapter
// (`internal/adapter/driving/cli`) receives fully-constructed
// driving-port implementations via its constructor.
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/pt9912/u-boot/internal/adapter/driven/confirm"
	"github.com/pt9912/u-boot/internal/adapter/driven/fs"
	"github.com/pt9912/u-boot/internal/adapter/driven/git"
	"github.com/pt9912/u-boot/internal/adapter/driven/progress"
	"github.com/pt9912/u-boot/internal/adapter/driven/yaml"
	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/hexagon/application"
)

// version is overridable at build time via -ldflags.
//
//	go build -ldflags="-X main.version=v0.1.0" ./cmd/uboot
var version = "0.1.0-dev"

func main() {
	// Signal-aware context: Ctrl-C / SIGTERM cancel the use-case
	// instead of killing the binary. For short operations like `init`
	// this is unobservable; for long-running subcommands (`up`,
	// `logs`, `doctor` against external systems) it lets the
	// application layer abort cleanly.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	os.Exit(run(ctx, os.Args[1:], os.Stdout, os.Stderr))
}

// run wires up the dependency graph and dispatches to the CLI app.
// Split from main so tests can exercise the wiring without spawning
// a process.
func run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	// Driven adapters.
	fsAdapter := fs.New()
	yamlAdapter := yaml.New()
	gitAdapter := git.New()
	progressAdapter := progress.NewText(stdout)
	// The confirm adapter renders to stderr (the prompt is operator-
	// facing UI, not machine-readable output) and reads stdin for the
	// answer — LH-FA-INIT-004 soft-existing-detection.
	confirmAdapter := confirm.New(os.Stdin, stderr)

	// Application services. The text-progress adapter renders
	// LH-FA-INIT-005 §609 / LH-FA-CLI-005A §262 affected-paths
	// events on stdout before any write happens; CLI-emitted post-
	// success messages land afterwards on the same stream; errors
	// go to stderr via the `fmt.Fprintf` below so the streams stay
	// distinct even when a caller pipes them together. Tests that
	// wrap stdout in a buffer must not interpose a separate flush.
	initSvc := application.NewInitProjectService(fsAdapter, yamlAdapter, gitAdapter, progressAdapter, confirmAdapter)

	// Driving adapter (CLI).
	app := cli.New(version, initSvc)

	err := app.Execute(ctx, args, stdout, stderr)
	if err != nil {
		fmt.Fprintf(stderr, "u-boot: %v\n", err)
	}
	return cli.ExitCode(err)
}
