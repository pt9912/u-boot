package cli

import (
	"log/slog"
	"strconv"
)

// verbosity is the level selected by the LH-FA-CLI-005 verbosity
// options. The option given last on the command line decides; the
// zero value is the default level.
type verbosity int

const (
	verbosityDefault verbosity = iota
	verbosityQuiet
	verbosityVerbose
	verbosityDebug
)

// slogLevel maps the verbosity to the logger level: --quiet is Error,
// --verbose and --debug are both Debug, the default is Info.
func (v verbosity) slogLevel() slog.Level {
	switch v {
	case verbosityQuiet:
		return slog.LevelError
	case verbosityVerbose, verbosityDebug:
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}

// verbosityFlag is the flag value behind --quiet, --verbose and
// --debug. The three flags share the App's verbosity field, so every
// occurrence on the command line overwrites the one before it; Cobra
// parses the arguments in command-line order and calls Set once per
// occurrence. It satisfies the pflag.Value interface.
type verbosityFlag struct {
	app  *App
	mode verbosity
}

// Set applies one occurrence: `--flag` or `--flag=true` selects the
// flag's mode, `--flag=false` returns to the default level when this
// flag's mode is the current one.
func (f *verbosityFlag) Set(s string) error {
	on, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	switch {
	case on:
		f.app.verbosity = f.mode
	case f.app.verbosity == f.mode:
		f.app.verbosity = verbosityDefault
	}
	return nil
}

// String reports whether the flag's mode is the selected one.
func (f *verbosityFlag) String() string {
	return strconv.FormatBool(f.app.verbosity == f.mode)
}

// Type reports "bool" so that the help output shows the option
// without a value placeholder.
func (*verbosityFlag) Type() string { return "bool" }

// quiet reports whether --quiet is the verbosity option given last.
// The output filters of the subcommands read this value, so a
// trailing --verbose or --debug lifts --quiet for them as well.
func (a *App) quiet() bool { return a.verbosity == verbosityQuiet }
