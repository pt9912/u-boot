package application

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"time"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// sandboxStatusDir is the directory in the devcontainer where the three
// start scripts write their status file (LH-FA-DEV-007); the generated
// Dockerfile creates it.
const sandboxStatusDir = "/var/lib/u-boot/status"

const (
	checkIDSandboxFallbackVFS    = "devcontainer.sandbox.fallback.vfs"
	checkIDSandboxFallbackDaemon = "devcontainer.sandbox.fallback.daemon"
	checkIDSandboxFallbackEgress = "devcontainer.sandbox.fallback.egress"
)

// fallbackProbe describes one degradation fallback of LH-FA-DEV-007 /
// LH-FA-DEV-008 whose start script writes a status file.
type fallbackProbe struct {
	// checkID is the `doctor` item the state is reported under.
	checkID string
	// code is the requirement a warning of `generate`/`init` names.
	code string
	// file is the name of the status file below [sandboxStatusDir].
	file string
	// what names the capability in a sentence.
	what string
	// fallback is the sentence for `result=fallback`.
	fallback string
	// failed is the sentence for `result=failed`.
	failed string
	// hint is the repair hint of the warning and error items.
	hint string
}

func vfsProbe() fallbackProbe {
	return fallbackProbe{
		checkID:  checkIDSandboxFallbackVFS,
		code:     "LH-FA-DEV-007",
		file:     "vfs",
		what:     "the nested Podman storage",
		fallback: "nested Podman uses vfs storage because /dev/fuse was not available in the container",
		failed:   "the start script ended with exit code 11 because /dev/fuse was not available in the container (onUnavailable: fail)",
		hint:     "Start the container with the generated runArgs (`--device=/dev/fuse`); run `u-boot generate devcontainer` and rebuild the container if devcontainer.json was edited.",
	}
}

func daemonProbe() fallbackProbe {
	return fallbackProbe{
		checkID:  checkIDSandboxFallbackDaemon,
		code:     "LH-FA-DEV-007",
		file:     "daemon",
		what:     "the nested Docker daemon",
		fallback: "the container has no Docker daemon because it did not become ready",
		failed:   "the start script ended with exit code 11 because the nested Docker daemon did not become ready (onUnavailable: fail)",
		hint:     "Start the container with the generated runArgs (`--privileged`); the daemon log is /var/log/u-boot-dockerd.log in the container.",
	}
}

func egressProbe() fallbackProbe {
	return fallbackProbe{
		checkID:  checkIDSandboxFallbackEgress,
		code:     "LH-FA-DEV-008",
		file:     "egress",
		what:     "the egress restriction",
		fallback: "the egress restriction is not active because the NET_ADMIN capability was not granted",
		failed:   "the start script ended with exit code 11: the NET_ADMIN capability was not granted (onUnavailable: fail), the user 'nobody' for the local resolver does not exist, the nftables rules could not be applied, or the local resolver did not start or does not run as 'nobody' (the last three under any onUnavailable)",
		hint:     "Start the container with the generated runArgs (`--cap-add=NET_ADMIN`), or restrict egress at network or DNS level. After exit code 11 with the capability granted, read the message of the start script on standard error in the container log: it names a missing user 'nobody', nftables rules that could not be applied, or a resolver that did not start or runs under another user.",
	}
}

// configuredFallbackProbes returns the fallbacks whose capability the
// configuration asks for, in the order vfs, daemon, egress. Without the
// sandbox profile no start script exists and nothing is returned.
func configuredFallbackProbes(dc *ubootYAMLDevcontainer, sandbox bool) []fallbackProbe {
	if !sandbox {
		return nil
	}
	var out []fallbackProbe
	if nestedPodman(dc) {
		out = append(out, vfsProbe())
	}
	if nestedDocker(dc) {
		out = append(out, daemonProbe())
	}
	if egressEnabled(dc) {
		out = append(out, egressProbe())
	}
	return out
}

// statusOutcome classifies what reading a status file yielded.
type statusOutcome int

const (
	// statusReported: a valid document was read.
	statusReported statusOutcome = iota
	// statusUndetermined: nothing could be read (no container, no
	// engine, no file); never a finding.
	statusUndetermined
	// statusInvalid: a file exists but is not a valid document.
	statusInvalid
)

// statusReading is the result of reading one status file.
type statusReading struct {
	outcome statusOutcome
	status  domain.SandboxStatus
	// stopped is true when the container the document comes from is
	// stopped (the state of its last start).
	stopped bool
	// reason explains statusUndetermined.
	reason string
}

func undetermined(reason string) statusReading {
	return statusReading{outcome: statusUndetermined, reason: reason}
}

func invalidStatus() statusReading { return statusReading{outcome: statusInvalid} }

// fallbackReader reads the status files of a project's devcontainer. From
// the host it asks the [driven.DevcontainerInspector] once per run for the
// container and reads each file through it; inside the devcontainer
// (`inContainer`) it reads the files directly through the FileSystem
// port, because `docker` is not reachable there.
type fallbackReader struct {
	inspector   driven.DevcontainerInspector
	fs          driven.FileSystem
	inContainer bool
	baseDir     string

	looked    bool
	container driven.DevcontainerContainer
	found     bool
	findErr   error
}

func (r *fallbackReader) read(ctx context.Context, p fallbackProbe) statusReading {
	file := path.Join(sandboxStatusDir, p.file)
	if r.inContainer {
		return r.readLocal(file)
	}
	if !r.looked {
		r.container, r.found, r.findErr = r.inspector.FindContainer(ctx, r.baseDir)
		r.looked = true
	}
	switch {
	case errors.Is(r.findErr, driven.ErrContainerEngineNotLocal):
		return undetermined("the container engine is remote, so a container there cannot be tied to this project folder")
	case r.findErr != nil:
		return undetermined("the container engine could not be queried")
	case !r.found:
		return undetermined("no devcontainer container was found for this project folder")
	}
	raw, err := r.inspector.ReadFile(ctx, r.container.ID, file, domain.SandboxStatusMaxBytes)
	switch {
	case errors.Is(err, driven.ErrContainerFileNotFound):
		return undetermined("the container has no status file (it was created before the status files existed, has not been started yet, or the start script ended before it wrote)")
	case errors.Is(err, driven.ErrContainerFileInvalid):
		return invalidStatus()
	case err != nil:
		return undetermined("the status file could not be read from the container")
	}
	reading := parseStatus(raw)
	reading.stopped = !r.container.Running
	return reading
}

// readLocal reads a status file inside the devcontainer. The entry must
// be a regular file within the size limit before it is read.
func (r *fallbackReader) readLocal(file string) statusReading {
	if exists, err := r.fs.Exists(sandboxStatusDir); err != nil || !exists {
		return undetermined("this u-boot runs in a container without the status directory " + sandboxStatusDir)
	}
	info, err := r.fs.Lstat(file)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return undetermined("the container has no status file (it was created before the status files existed, has not been started yet, or the start script ended before it wrote)")
		}
		return undetermined("the status file could not be read")
	}
	if !info.Mode().IsRegular() || info.Size() > domain.SandboxStatusMaxBytes {
		return invalidStatus()
	}
	raw, err := r.fs.ReadFile(file)
	if err != nil {
		return undetermined("the status file could not be read")
	}
	return parseStatus(raw)
}

func parseStatus(raw []byte) statusReading {
	status, err := domain.ParseSandboxStatus(raw)
	if err != nil {
		return invalidStatus()
	}
	return statusReading{outcome: statusReported, status: status}
}

// reportedSuffix names the time of the report and, for a stopped
// container, that it is the state of its last start.
func reportedSuffix(rd statusReading) string {
	s := " (reported at " + rd.status.At().Format(time.RFC3339)
	if rd.stopped {
		s += ", container stopped"
	}
	return s + ")"
}

// fallbackDiagnostic maps one reading to the `doctor` item
// (LH-FA-DEV-007): `fallback` is a warning that carries the fallback
// flag, `failed` an error, `ok` OK, not determinable an info, an invalid
// document a warning without the flag. The text names resolved fields
// only, never the content of the file.
func fallbackDiagnostic(p fallbackProbe, rd statusReading) domain.Diagnostic {
	switch rd.outcome {
	case statusUndetermined:
		return domain.Diagnostic{ID: p.checkID, Severity: domain.SeverityInfo,
			Message: fmt.Sprintf("The state of %s at the container start is not determinable: %s.", p.what, rd.reason),
			Hint:    "A container created before the status files existed needs `u-boot generate devcontainer` and a rebuild; no report does not mean that no fallback occurred."}
	case statusInvalid:
		return domain.Diagnostic{ID: p.checkID, Severity: domain.SeverityWarn,
			Message: fmt.Sprintf("The status file for %s is invalid; the state at the container start is not determinable.", p.what),
			Hint:    "The container start script writes the file; rebuild the container or start it again. An edited, oversized or non-regular file is ignored."}
	}
	suffix := reportedSuffix(rd)
	switch rd.status.Result() {
	case domain.SandboxStatusFallback:
		return domain.Diagnostic{ID: p.checkID, Severity: domain.SeverityWarn, DegradationFallback: true,
			Message: "Fallback: " + p.fallback + suffix + ".", Hint: p.hint}
	case domain.SandboxStatusFailed:
		return domain.Diagnostic{ID: p.checkID, Severity: domain.SeverityError,
			Message: "Start failed: " + p.failed + suffix + ".", Hint: p.hint}
	}
	return domain.Diagnostic{ID: p.checkID, Severity: domain.SeverityOK,
		Message: fmt.Sprintf("The container start script reported %s as available%s; the container reports this about itself.", p.what, suffix)}
}

// fallbackWording is the part of a fallback warning that depends on the
// command: how the state of the container is introduced and what the
// regenerated files need.
type fallbackWording struct {
	// subject introduces the reported time: "state of <subject> at its last
	// start (<time>): ".
	subject string
	// tail ends the warning text; empty for none.
	tail string
}

// generateFallbackWording is the text of `generate` and `init`: the
// container exists before the command runs and keeps its image.
func generateFallbackWording() fallbackWording {
	return fallbackWording{
		subject: "the existing devcontainer container",
		tail:    "; the regenerated files take effect after the container is rebuilt",
	}
}

// startFallbackWording is the text of `devcontainer up` and `devcontainer
// shell`: the command started (or found running) the container and regenerates
// no file.
func startFallbackWording() fallbackWording {
	return fallbackWording{subject: "the devcontainer container"}
}

// fallbackWarnings returns one warning per probe whose existing container
// reports `result=fallback` (LH-FA-DEV-007, LH-FA-CLI-005 exception 2):
// the warning carries the fallback flag and states that it is the state of
// the container's last start, in the wording of the command. Every other
// outcome (ok, failed, not determinable, invalid) is silent.
func fallbackWarnings(ctx context.Context, r *fallbackReader, probes []fallbackProbe, w fallbackWording) []driving.WarningEntry {
	var out []driving.WarningEntry
	for _, p := range probes {
		rd := r.read(ctx, p)
		if rd.outcome != statusReported || rd.status.Result() != domain.SandboxStatusFallback {
			continue
		}
		out = append(out, driving.WarningEntry{
			Code:  p.code,
			Level: "warn",
			Message: "state of " + w.subject + " at its last start (" +
				rd.status.At().Format(time.RFC3339) + "): " + p.fallback + w.tail,
			Subject:             ".devcontainer/devcontainer.json",
			DegradationFallback: true,
		})
	}
	return out
}
