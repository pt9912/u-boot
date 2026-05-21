// Command u-boot is the developer environment bootloader for Docker-based
// projects. See spec/lastenheft.md for the full functional specification.
//
// This is the MVP bootstrap entry point. Subcommands like `init`, `add`,
// `up`, `down`, `doctor`, `generate`, `config` are stubbed and will be
// wired in via internal/ packages in later slices.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// version is overridable at build time via -ldflags.
//
//	go build -ldflags="-X main.version=v0.1.0" ./cmd/uboot
var version = "0.1.0-dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("u-boot", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var (
		showHelp    bool
		showVersion bool
	)
	fs.BoolVar(&showHelp, "help", false, "Show usage and exit")
	fs.BoolVar(&showHelp, "h", false, "Show usage and exit (shorthand)")
	fs.BoolVar(&showVersion, "version", false, "Print the version and exit")
	fs.Usage = func() { printHelp(stderr) }

	if err := fs.Parse(args); err != nil {
		return 2
	}

	if showVersion {
		fmt.Fprintln(stdout, version)
		return 0
	}

	if showHelp || fs.NArg() == 0 {
		printHelp(stdout)
		return 0
	}

	cmd := fs.Arg(0)
	fmt.Fprintf(stderr, "u-boot: command %q is not implemented yet in this MVP bootstrap\n", cmd)
	fmt.Fprintln(stderr, "Run 'u-boot --help' for the list of planned commands.")
	return 2
}

func printHelp(w io.Writer) {
	fmt.Fprintln(w, "u-boot — a developer environment bootloader for Docker-based projects")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  u-boot [--version] [--help] <command> [args]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Planned commands (see spec/lastenheft.md):")
	fmt.Fprintln(w, "  init              Initialize a new project")
	fmt.Fprintln(w, "  add <service>     Add a service add-on (postgres, keycloak, otel)")
	fmt.Fprintln(w, "  remove <service>  Remove a service add-on")
	fmt.Fprintln(w, "  up                Start the development environment")
	fmt.Fprintln(w, "  down              Stop the development environment")
	fmt.Fprintln(w, "  doctor            Check local prerequisites")
	fmt.Fprintln(w, "  logs              Show service logs")
	fmt.Fprintln(w, "  generate <kind>   Generate artifacts (changelog, readme, env-example, devcontainer)")
	fmt.Fprintln(w, "  config            Show or update project configuration")
	fmt.Fprintln(w, "  template          Manage project templates")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Flags:")
	fmt.Fprintln(w, "  --help, -h        Show this help")
	fmt.Fprintln(w, "  --version         Print the version")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Prerequisites (LH-FA-DIAG-002):")
	fmt.Fprintln(w, "  - Docker Engine >= 24.0.0")
	fmt.Fprintln(w, "  - Docker Compose >= 2.20.0")
	fmt.Fprintln(w, "  - Git")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "This MVP bootstrap exposes --help and --version only. Subcommands")
	fmt.Fprintln(w, "follow in later slices (see docs/plan/planning/).")
}
