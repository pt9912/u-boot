package application

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// Remove implements [driving.DevcontainerUseCase.Remove]
// (`u-boot devcontainer remove`, LH-FA-DEV-010). The steps run in this order,
// and nothing changes before the confirmation has been given:
//
//  1. the instance name is checked, before the engine is asked;
//  2. the container is searched like `down` searches it;
//  3. with Volumes the volumes of the container are read once and the
//     confirmation is asked, naming them;
//  4. a running container is stopped;
//  5. the container is removed, with `-v` under Volumes;
//  6. under Volumes each volume of the set that still exists is removed, and
//     the volumes are read again to decide which of them are gone.
//
// The command reads no project file and starts no helper container.
func (s *DevcontainerStartService) Remove(ctx context.Context, req driving.DevcontainerRemoveRequest) (driving.DevcontainerRemoveResponse, error) {
	if err := s.wired("remove", req.BaseDir); err != nil {
		return driving.DevcontainerRemoveResponse{}, err
	}
	if err := validInstance("remove", req.Instance); err != nil {
		return driving.DevcontainerRemoveResponse{}, err
	}
	container, found, err := s.findContainer(ctx, "remove", req.BaseDir, req.Instance)
	if err != nil {
		return driving.DevcontainerRemoveResponse{}, err
	}
	if !found {
		return driving.DevcontainerRemoveResponse{RemovedVolumeNames: []string{}}, nil
	}
	var volumes []string
	if req.Volumes {
		if volumes, err = s.containerVolumes(ctx, container.ID); err != nil {
			return driving.DevcontainerRemoveResponse{}, err
		}
		if err := s.approveRemoval(ctx, req, container.ID, volumes); err != nil {
			return driving.DevcontainerRemoveResponse{}, err
		}
	}
	s.logger.Debug("devcontainer remove", "instance", req.Instance, "volumes", req.Volumes, "running", container.Running, "volumeCount", len(volumes))
	resp := driving.DevcontainerRemoveResponse{ContainerID: container.ID, RemovedVolumeNames: []string{}}
	if container.Running {
		if err := s.launcher.Stop(ctx, container.ID); err != nil {
			return driving.DevcontainerRemoveResponse{}, fmt.Errorf(
				"devcontainer remove: docker stop of container %s failed: %w; the container still exists and may still run, no volume was touched",
				shortContainerID(container.ID), err)
		}
		resp.Stopped = true
	}
	if err := s.launcher.RemoveContainer(ctx, container.ID, req.Volumes); err != nil {
		return driving.DevcontainerRemoveResponse{}, fmt.Errorf(
			"devcontainer remove: docker rm of container %s failed: %w; the container still exists%s, no volume was touched",
			shortContainerID(container.ID), err, stoppedNote(resp.Stopped))
	}
	resp.Removed = true
	if !req.Volumes {
		return resp, nil
	}
	removed, err := s.removeVolumes(ctx, container.ID, volumes)
	if err != nil {
		return driving.DevcontainerRemoveResponse{}, err
	}
	resp.RemovedVolumeNames = removed
	return resp, nil
}

// stoppedNote is the part of the error text that tells a stopped container.
func stoppedNote(stopped bool) string {
	if stopped {
		return " (stopped)"
	}
	return ""
}

// shortContainerID is the first 12 characters of a container ID, the form
// `docker ps` shows.
func shortContainerID(id string) string {
	if len(id) > 12 {
		return id[:12]
	}
	return id
}

// containerVolumes reads the volumes the container mounts, sorted by name and
// without duplicates. A failure of the engine ends the call with
// [driven.ErrDockerUnavailable]: the question could not be asked.
func (s *DevcontainerStartService) containerVolumes(ctx context.Context, containerID string) ([]string, error) {
	names, err := s.inspector.ContainerVolumes(ctx, containerID)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, fmt.Errorf("devcontainer remove: the volumes of container %s could not be read: %w (%s)",
			shortContainerID(containerID), driven.ErrDockerUnavailable, err)
	}
	return sortedDistinct(names), nil
}

// sortedDistinct returns the names sorted and without repetition; it is never
// nil.
func sortedDistinct(names []string) []string {
	out := make([]string, 0, len(names))
	seen := map[string]bool{}
	for _, n := range names {
		if !seen[n] {
			seen[n] = true
			out = append(out, n)
		}
	}
	sort.Strings(out)
	return out
}

// approveRemoval is the confirmation gate of `remove --volumes`: `--yes` gives
// it, `--no-interactive` and `--json` refuse without asking, and otherwise the
// [driven.Confirmer] asks with the container and the volumes. A refusal ends
// the call with [driving.ErrConfirmationRequired], an unreadable answer with
// [driving.ErrConfirmerUnavailable]; both say that nothing was stopped or
// removed and name the container and the volumes.
func (s *DevcontainerStartService) approveRemoval(ctx context.Context, req driving.DevcontainerRemoveRequest, containerID string, volumes []string) error {
	if req.AssumeYes {
		return nil
	}
	held := "no volume"
	if len(volumes) > 0 {
		held = "volumes " + strings.Join(volumes, ", ")
	}
	nothing := fmt.Sprintf("container %s (%s) was neither stopped nor removed", shortContainerID(containerID), held)
	switch {
	case req.NonInteractive:
		return fmt.Errorf("devcontainer remove: %w: --volumes needs a confirmation and --no-interactive forbids asking; pass --yes; %s",
			driving.ErrConfirmationRequired, nothing)
	case req.SilenceConfirmer:
		return fmt.Errorf("devcontainer remove: %w: --volumes needs a confirmation and --json never asks; pass --yes; %s",
			driving.ErrConfirmationRequired, nothing)
	}
	confirmed, err := s.confirmer.ConfirmRemoveDevcontainer(ctx, containerID, volumes)
	if err != nil {
		return fmt.Errorf("devcontainer remove: %w: the answer could not be read (%s); %s",
			driving.ErrConfirmerUnavailable, err, nothing)
	}
	if !confirmed {
		return fmt.Errorf("devcontainer remove: %w: the confirmation was declined; %s", driving.ErrConfirmationRequired, nothing)
	}
	return nil
}

// removeVolumes removes each volume of the set that exists after the container
// is gone and returns the volumes of the set that no longer exist, sorted. It
// tries every volume even when one fails, then reads the volumes again. A
// failed removal, a volume that remains and a failed read end the call with
// [driven.ErrDevcontainerRuntime]; the error names the volumes that are gone
// and the ones that remain.
func (s *DevcontainerStartService) removeVolumes(ctx context.Context, containerID string, set []string) ([]string, error) {
	if len(set) == 0 {
		return []string{}, nil
	}
	short := shortContainerID(containerID)
	existing, err := s.inspector.ExistingVolumes(ctx, set)
	if err != nil {
		return nil, volumeFailure(ctx, fmt.Errorf(
			"devcontainer remove: container %s is removed, but the engine could not list its volumes (%s); the state of the volumes %s is not established; check with `docker volume ls`",
			short, err, strings.Join(set, ", ")))
	}
	var failures []string
	for _, name := range existing {
		if err := s.launcher.RemoveVolume(ctx, name); err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return nil, ctxErr
			}
			failures = append(failures, name+": "+err.Error())
		}
	}
	remaining, err := s.inspector.ExistingVolumes(ctx, set)
	if err != nil {
		return nil, volumeFailure(ctx, fmt.Errorf(
			"devcontainer remove: container %s is removed, but the engine could not list the volumes after their removal (%s); the state of the volumes %s is not established; check with `docker volume ls`%s",
			short, err, strings.Join(set, ", "), failureNote(failures)))
	}
	removed := difference(set, remaining)
	if len(remaining) > 0 || len(failures) > 0 {
		return nil, volumeFailure(ctx, fmt.Errorf(
			"devcontainer remove: container %s is removed, but %s: volumes removed: %s; volumes that remain: %s%s; remove them with `docker volume rm <name>` once no container holds them",
			short, remainingCause(len(failures) > 0), listOrNone(removed), listOrNone(remaining), failureNote(failures)))
	}
	return removed, nil
}

// volumeFailure gives a failure of the volume step the runtime identity
// (exit 12), unless the context ended the call.
func volumeFailure(ctx context.Context, err error) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}
	return fmt.Errorf("%w: %w", err, driven.ErrDevcontainerRuntime)
}

func remainingCause(failed bool) string {
	if failed {
		return "docker volume rm failed"
	}
	return "volumes of the set still exist"
}

func failureNote(failures []string) string {
	if len(failures) == 0 {
		return ""
	}
	return "; docker volume rm: " + strings.Join(failures, "; ")
}

func listOrNone(names []string) string {
	if len(names) == 0 {
		return "none"
	}
	return strings.Join(names, ", ")
}

// difference returns the names of set that are not in other, in the order of
// set; it is never nil.
func difference(set, other []string) []string {
	skip := make(map[string]bool, len(other))
	for _, n := range other {
		skip[n] = true
	}
	out := make([]string, 0, len(set))
	for _, n := range set {
		if !skip[n] {
			out = append(out, n)
		}
	}
	return out
}
