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
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/pt9912/u-boot/internal/adapter/driven/clock"
	"github.com/pt9912/u-boot/internal/adapter/driven/confirm"
	"github.com/pt9912/u-boot/internal/adapter/driven/docker"
	"github.com/pt9912/u-boot/internal/adapter/driven/fs"
	"github.com/pt9912/u-boot/internal/adapter/driven/git"
	"github.com/pt9912/u-boot/internal/adapter/driven/logger"
	"github.com/pt9912/u-boot/internal/adapter/driven/netprobe"
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
	// The logger adapter (LH-QA-004) renders to stderr in text form
	// by default, at a level mutable via *slog.LevelVar so the
	// LH-FA-CLI-005 verbosity flags can flip it from
	// PersistentPreRunE. Default Info; --debug / --verbose raise to
	// Debug, --quiet lowers to Warn.
	logLevel := new(slog.LevelVar)
	logLevel.Set(slog.LevelInfo)
	logAdapter := logger.New(stderr, logger.FormatText, logLevel)
	// The docker probe (LH-FA-DIAG-002) shells out to `docker` for
	// read-only diagnostics (version + reachability + compose version).
	// Used by the doctor service; M6's DockerEngine port stays separate.
	dockerAdapter := docker.New()
	// The docker engine (LH-FA-UP-002) shells out to `docker compose`
	// for state-mutating operations: ComposeUp / ComposeDown /
	// ComposePs. Used by up/down services. Distinct from the probe
	// adapter — the slice plan §Sentinel-Schichtung pins separate
	// pre-probe classification (ErrDockerUnavailable vs.
	// ErrComposeRuntime) which the engine adapter enforces.
	dockerEngineAdapter := docker.NewEngine()
	// The netprobe adapter (M6-T3) provides TCP-reachability checks
	// for the UpService polling loop. Stateless; depguard rule
	// `application-no-net` enforces all net.* usage funnels through here.
	netprobeAdapter := netprobe.New()
	// The clock adapter (M6-T4-fund) wraps time.Now / time.Sleep
	// so the UpService polling-loop iteration timing is injectable
	// in tests.
	clockAdapter := clock.New()

	// Application services. The text-progress adapter renders
	// LH-FA-INIT-005 §609 / LH-FA-CLI-005A §262 affected-paths
	// events on stdout before any write happens; CLI-emitted post-
	// success messages land afterwards on the same stream; errors
	// go to stderr via the `fmt.Fprintf` below so the streams stay
	// distinct even when a caller pipes them together. Tests that
	// wrap stdout in a buffer must not interpose a separate flush.
	initSvc := application.NewInitProjectService(fsAdapter, yamlAdapter, gitAdapter, progressAdapter, confirmAdapter, logAdapter)
	doctorSvc := application.NewDoctorService(fsAdapter, yamlAdapter, gitAdapter, dockerAdapter, logAdapter)
	addSvc := application.NewAddServiceService(fsAdapter, yamlAdapter, logAdapter)
	upSvc := application.NewUpService(fsAdapter, yamlAdapter, dockerEngineAdapter, netprobeAdapter, clockAdapter, logAdapter)
	downSvc := application.NewDownService(fsAdapter, dockerEngineAdapter, confirmAdapter, logAdapter)
	generateSvc := application.NewGenerateService(fsAdapter, yamlAdapter, logAdapter)

	// Driving adapter (CLI).
	app := cli.New(version, initSvc, doctorSvc, addSvc, upSvc, downSvc, generateSvc, cli.WithLogLevel(logLevel))

	err := app.Execute(ctx, args, stdout, stderr)
	if err != nil {
		fmt.Fprintf(stderr, "u-boot: %v\n", err)
	}
	return cli.ExitCode(err)
}
