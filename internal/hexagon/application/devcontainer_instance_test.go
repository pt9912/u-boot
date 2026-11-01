package application_test

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// instanceContainer is one container of the project folder in
// [instanceInspector]; instance is the value of its `u-boot.instance` label,
// empty for the container without an instance.
type instanceContainer struct {
	id       string
	running  bool
	instance string
}

// instanceInspector is the [driven.DevcontainerInspector] test double with the
// search of the port: containers newest first, an instance selects the
// container with that label, no instance selects the newest without the label
// and fails with [driven.ErrOnlyInstanceContainers] when all carry it.
type instanceInspector struct {
	containers []instanceContainer
	asked      []string
}

func (f *instanceInspector) FindContainer(_ context.Context, _, instance string) (driven.DevcontainerContainer, bool, error) {
	f.asked = append(f.asked, instance)
	for _, c := range f.containers {
		if c.instance == instance {
			return driven.DevcontainerContainer{ID: c.id, Running: c.running}, true, nil
		}
	}
	if instance == "" && len(f.containers) > 0 {
		return driven.DevcontainerContainer{}, false, driven.ErrOnlyInstanceContainers
	}
	return driven.DevcontainerContainer{}, false, nil
}

func (f *instanceInspector) ReadFile(context.Context, string, string, int) ([]byte, error) {
	return nil, driven.ErrContainerFileNotFound
}

const (
	idOfExisting = fallbackContainerID
	idOfInstance = "1111111111111111111111111111111111111111111111111111111111111111"
	idOfOther    = "2222222222222222222222222222222222222222222222222222222222222222"
)

// newInstanceService builds the service on the launcher double and an
// [instanceInspector]; the project has a plain devcontainer.json.
func newInstanceService(t *testing.T, l *fakeLauncher, insp *instanceInspector) *application.DevcontainerStartService {
	t.Helper()
	fs := newFakeFS()
	fs.markDirExists(devcontainerTestBaseDir)
	if err := fs.WriteFile(filepath.Join(devcontainerTestBaseDir, ".devcontainer", "devcontainer.json"), []byte(plainConfig), 0o644); err != nil {
		t.Fatal(err)
	}
	return application.NewDevcontainerStartService(fs, &fakeYAML{}, l, insp, nil)
}

// instanceService returns the three calls of the service of
// [newInstanceService].
func instanceService(t *testing.T, l *fakeLauncher, insp *instanceInspector) (
	up func(instance string, recreate bool) (driving.DevcontainerUpResponse, error),
	down func(instance string) (driving.DevcontainerDownResponse, error),
	shell func(instance string) error,
) {
	t.Helper()
	svc := newInstanceService(t, l, insp)
	up = func(instance string, recreate bool) (driving.DevcontainerUpResponse, error) {
		return svc.Up(context.Background(), driving.DevcontainerUpRequest{BaseDir: devcontainerTestBaseDir, Instance: instance, Recreate: recreate})
	}
	down = func(instance string) (driving.DevcontainerDownResponse, error) {
		return svc.Down(context.Background(), driving.DevcontainerDownRequest{BaseDir: devcontainerTestBaseDir, Instance: instance})
	}
	shell = func(instance string) error {
		return svc.Shell(context.Background(), driving.DevcontainerShellRequest{BaseDir: devcontainerTestBaseDir, Instance: instance})
	}
	return up, down, shell
}

// ADR-0038 decision 6 (LH-FA-DEV-010): an instance name the rule rejects ends
// `up`, `down` and `shell` with the sentinel of exit code 10 before the
// inspector and the launcher.
func TestDevcontainerInstance_InvalidNameEndsBeforeAnyPort(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"Abc", "a=b", "-a", "a-", "a_b", "a b", strings.Repeat("a", 33)} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			l := newLauncher()
			insp := &instanceInspector{containers: []instanceContainer{{id: idOfExisting, running: true}}}
			up, down, shell := instanceService(t, l, insp)

			if _, err := up(name, false); !errors.Is(err, driving.ErrDevcontainerInstanceInvalid) {
				t.Errorf("up: err = %v, want ErrDevcontainerInstanceInvalid", err)
			}
			if _, err := down(name); !errors.Is(err, driving.ErrDevcontainerInstanceInvalid) {
				t.Errorf("down: err = %v, want ErrDevcontainerInstanceInvalid", err)
			}
			if err := shell(name); !errors.Is(err, driving.ErrDevcontainerInstanceInvalid) {
				t.Errorf("shell: err = %v, want ErrDevcontainerInstanceInvalid", err)
			}
			if len(l.upOpts) != 0 || len(l.stopIDs) != 0 || len(l.shellOpts) != 0 || len(insp.asked) != 0 {
				t.Errorf("a port was reached: Up %d, Stop %d, Shell %d, FindContainer %d", len(l.upOpts), len(l.stopIDs), len(l.shellOpts), len(insp.asked))
			}
		})
	}
}

// ADR-0038 decision 5 (LH-FA-DEV-010): a project with containers of instances
// only ends `up` (also with `--recreate`), `down` and `shell` without an
// instance with the sentinel of exit code 10 and a hint at `--instance`; no
// helper container, no `docker stop` and no session follow.
func TestDevcontainerInstance_OnlyInstancesEndsBeforeTheLauncher(t *testing.T) {
	t.Parallel()
	l := newLauncher()
	insp := &instanceInspector{containers: []instanceContainer{
		{id: idOfInstance, running: true, instance: "i2"},
		{id: idOfOther, running: false, instance: "i1"},
	}}
	up, down, shell := instanceService(t, l, insp)

	check := func(what string, err error) {
		t.Helper()
		if !errors.Is(err, driving.ErrDevcontainerOnlyInstances) {
			t.Errorf("%s: err = %v, want ErrDevcontainerOnlyInstances", what, err)
			return
		}
		if !strings.Contains(err.Error(), "--instance") {
			t.Errorf("%s: message %q does not point at --instance", what, err)
		}
	}
	_, err := up("", false)
	check("up", err)
	_, err = up("", true)
	check("up --recreate", err)
	_, err = down("")
	check("down", err)
	check("shell", shell(""))
	if len(l.upOpts) != 0 || len(l.stopIDs) != 0 || len(l.shellOpts) != 0 {
		t.Errorf("the launcher was reached: Up %d, Stop %d, Shell %d", len(l.upOpts), len(l.stopIDs), len(l.shellOpts))
	}
}

// ADR-0038 decision 5 (SPEC-049, paragraph *Instanzen*): the check for
// "only instances" of `up` and `shell` follows the disclosure of the socket,
// which the call raises before it asks the engine; the invalid name does not
// reach the disclosure.
func TestDevcontainerInstance_DisclosureComesBeforeTheOnlyInstancesCheck(t *testing.T) {
	t.Parallel()
	insp := &instanceInspector{containers: []instanceContainer{{id: idOfInstance, running: true, instance: "i1"}}}
	disclosed := 0
	onWarning := func(driving.WarningEntry) { disclosed++ }
	svc := newInstanceService(t, newLauncher(), insp)

	_, err := svc.Up(context.Background(), driving.DevcontainerUpRequest{BaseDir: devcontainerTestBaseDir, OnWarning: onWarning})
	if !errors.Is(err, driving.ErrDevcontainerOnlyInstances) || disclosed != 1 {
		t.Errorf("up: err = %v, disclosures = %d, want ErrDevcontainerOnlyInstances after one disclosure", err, disclosed)
	}
	disclosed = 0
	err = svc.Shell(context.Background(), driving.DevcontainerShellRequest{BaseDir: devcontainerTestBaseDir, OnWarning: onWarning})
	if !errors.Is(err, driving.ErrDevcontainerOnlyInstances) || disclosed != 1 {
		t.Errorf("shell: err = %v, disclosures = %d, want ErrDevcontainerOnlyInstances after one disclosure", err, disclosed)
	}
	disclosed = 0
	_, err = svc.Up(context.Background(), driving.DevcontainerUpRequest{BaseDir: devcontainerTestBaseDir, Instance: "Bad", OnWarning: onWarning})
	if !errors.Is(err, driving.ErrDevcontainerInstanceInvalid) || disclosed != 0 {
		t.Errorf("up with an invalid name: err = %v, disclosures = %d, want ErrDevcontainerInstanceInvalid and none", err, disclosed)
	}
}

// ADR-0038 decisions 4 and 9, the untested case of the ADR (existing
// container and a newer instance): without an instance `up`, `down` and
// `shell` meet the container without an instance, with an instance the
// container of that instance; the instance reaches the launcher and every
// question to the inspector.
func TestDevcontainerInstance_ExistingContainerAndNewerInstance(t *testing.T) {
	t.Parallel()
	newInsp := func() *instanceInspector {
		return &instanceInspector{containers: []instanceContainer{
			{id: idOfInstance, running: true, instance: "i1"},
			{id: idOfExisting, running: true},
		}}
	}

	t.Run("up without an instance", func(t *testing.T) {
		t.Parallel()
		l := newLauncher()
		insp := newInsp()
		up, _, _ := instanceService(t, l, insp)
		resp, err := up("", false)
		if err != nil || !resp.AlreadyRunning || resp.ContainerID != idOfExisting {
			t.Fatalf("resp = %+v, err = %v, want the container without an instance, already running", resp, err)
		}
		if len(l.upOpts) != 1 || l.upOpts[0].Instance != "" {
			t.Errorf("launcher options = %+v, want no instance", l.upOpts)
		}
	})
	t.Run("up with the instance", func(t *testing.T) {
		t.Parallel()
		l := newLauncher()
		l.upResult.ContainerID = idOfInstance
		insp := newInsp()
		up, _, _ := instanceService(t, l, insp)
		resp, err := up("i1", false)
		if err != nil || !resp.AlreadyRunning {
			t.Fatalf("resp = %+v, err = %v, want the instance, already running", resp, err)
		}
		if len(l.upOpts) != 1 || l.upOpts[0].Instance != "i1" {
			t.Errorf("launcher options = %+v, want the instance i1", l.upOpts)
		}
		for _, asked := range insp.asked {
			if asked != "i1" {
				t.Errorf("the inspector was asked for the instance %q, want i1", asked)
			}
		}
	})
	t.Run("up with an instance that does not exist yet", func(t *testing.T) {
		t.Parallel()
		l := newLauncher()
		up, _, _ := instanceService(t, l, newInsp())
		resp, err := up("i9", false)
		if err != nil || resp.AlreadyRunning {
			t.Fatalf("resp = %+v, err = %v, want a start of a new instance", resp, err)
		}
		if len(l.upOpts) != 1 || l.upOpts[0].Instance != "i9" {
			t.Errorf("launcher options = %+v, want the instance i9", l.upOpts)
		}
	})
	t.Run("down without an instance", func(t *testing.T) {
		t.Parallel()
		l := newLauncher()
		_, down, _ := instanceService(t, l, newInsp())
		resp, err := down("")
		if err != nil || !resp.Stopped || resp.ContainerID != idOfExisting {
			t.Fatalf("resp = %+v, err = %v", resp, err)
		}
		if len(l.stopIDs) != 1 || l.stopIDs[0] != idOfExisting {
			t.Errorf("Stop calls = %v, want only the container without an instance", l.stopIDs)
		}
	})
	t.Run("down with the instance", func(t *testing.T) {
		t.Parallel()
		l := newLauncher()
		_, down, _ := instanceService(t, l, newInsp())
		resp, err := down("i1")
		if err != nil || !resp.Stopped || resp.ContainerID != idOfInstance {
			t.Fatalf("resp = %+v, err = %v", resp, err)
		}
		if len(l.stopIDs) != 1 || l.stopIDs[0] != idOfInstance {
			t.Errorf("Stop calls = %v, want only the instance", l.stopIDs)
		}
	})
	t.Run("shell without an instance", func(t *testing.T) {
		t.Parallel()
		l := newLauncher()
		insp := newInsp()
		_, _, shell := instanceService(t, l, insp)
		if err := shell(""); err != nil {
			t.Fatalf("Shell: %v", err)
		}
		if len(l.upOpts) != 1 || l.upOpts[0].Instance != "" || len(l.shellOpts) != 1 {
			t.Errorf("Up options %+v, Shell options %d, want one start without an instance and one session", l.upOpts, len(l.shellOpts))
		}
	})
	t.Run("shell with the instance", func(t *testing.T) {
		t.Parallel()
		l := newLauncher()
		_, _, shell := instanceService(t, l, newInsp())
		if err := shell("i1"); err != nil {
			t.Fatalf("Shell: %v", err)
		}
		if len(l.upOpts) != 1 || l.upOpts[0].Instance != "i1" {
			t.Errorf("Up options %+v, want the instance i1", l.upOpts)
		}
	})
}

// ADR-0038 consequences (not measured, pinned here): `down --instance <name>`
// with a name that has no container is "nothing to stop", exit code 0, and
// calls no `docker stop`; so is `down` without an instance in a project
// without any container.
func TestDevcontainerInstance_DownWithoutAMatchingContainer(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name       string
		containers []instanceContainer
		instance   string
	}{
		{"unknown instance", []instanceContainer{{id: idOfExisting, running: true}, {id: idOfInstance, running: true, instance: "i1"}}, "i2"},
		{"instance in a project without a container", nil, "i1"},
		{"no instance in a project without a container", nil, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l := newLauncher()
			_, down, _ := instanceService(t, l, &instanceInspector{containers: tc.containers})
			resp, err := down(tc.instance)
			if err != nil || resp != (driving.DevcontainerDownResponse{}) || len(l.stopIDs) != 0 {
				t.Errorf("resp = %+v, err = %v, Stop calls = %v, want nothing to stop", resp, err, l.stopIDs)
			}
		})
	}
}

// ADR-0038 decision 9: after the start the status files are read from the
// container of the instance, not from the one without an instance.
func TestDevcontainerInstance_StatusReadingNamesTheInstance(t *testing.T) {
	t.Parallel()
	for _, instance := range []string{"", "i1"} {
		t.Run("instance "+instance, func(t *testing.T) {
			t.Parallel()
			insp := newFakeInspector(true)
			insp.setStatus("vfs", "ok")
			svc, _ := devcontainerSvc(t, plainConfig, yamlPodman, newLauncher(), insp)
			if _, err := svc.Up(context.Background(), driving.DevcontainerUpRequest{BaseDir: devcontainerTestBaseDir, Instance: instance}); err != nil {
				t.Fatalf("Up: %v", err)
			}
			if len(insp.readPaths) == 0 {
				t.Fatal("no status file was read")
			}
			for _, asked := range insp.findInstances {
				if asked != instance {
					t.Errorf("the inspector was asked for the instance %q, want %q", asked, instance)
				}
			}
		})
	}
}

// A project with containers of instances only is "not determinable" for the
// status reading, never a finding: when the check before the launcher passes
// and the question after the start meets that state (a container created in
// between), the status file of the other container is not read.
func TestDevcontainerInstance_OnlyInstancesIsNotDeterminableForTheStatusReading(t *testing.T) {
	t.Parallel()
	insp := newFakeInspector(true)
	insp.findAnswers = []findAnswer{{container: driven.DevcontainerContainer{ID: fallbackContainerID, Running: true}, found: true}, {err: driven.ErrOnlyInstanceContainers}}
	insp.setStatus("vfs", "failed")
	svc, _ := devcontainerSvc(t, plainConfig, yamlPodman, newLauncher(), insp)
	if _, err := svc.Up(context.Background(), driving.DevcontainerUpRequest{BaseDir: devcontainerTestBaseDir}); err != nil {
		t.Errorf("Up: %v, want a success: the status file of the container that was not found must not be read", err)
	}
	if len(insp.readPaths) != 0 {
		t.Errorf("status files read: %v", insp.readPaths)
	}
}
