// Package domain holds u-boot's pure domain types and invariants.
// It must not depend on application, port, adapter, or any I/O library
// (LH-FA-ARCH-002, LH-FA-ARCH-003).
package domain

import "fmt"

// DevcontainerContainerNamePrefix starts every composed container name of the
// devcontainer the launcher starts (`u-boot devcontainer up`, LH-FA-DEV-010),
// so the name is recognizable beside the engine's random names.
const DevcontainerContainerNamePrefix = "u-boot"

// ComposeDevcontainerContainerName joins the container name of the started
// devcontainer from the parts: `u-boot-<project>` and, with a non-empty
// instance, `u-boot-<project>-<instance>`. It composes only; whether the
// engine accepts the name is not decided here — the admissible character set
// of the engine (`--name`, measured with Docker Engine 29.9.2) lives at the
// driven port next to the call, and an invalid part ends the call there as a
// validation error, not as a silent skip.
//
// A project part is the project name from u-boot.yaml or, without the file,
// the normalized basename of the project folder; an instance part is an
// admissible instance name. An empty project part is not composable and ends
// with an error.
func ComposeDevcontainerContainerName(project, instance string) (string, error) {
	name := DevcontainerContainerNamePrefix
	if project == "" {
		return "", fmt.Errorf("container name %q: empty project name part", name)
	}
	name += "-" + project
	if instance != "" {
		name += "-" + instance
	}
	return name, nil
}
