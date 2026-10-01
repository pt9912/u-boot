//go:build docker

// Real-Docker checks of the Up/Down cluster: down --volumes names,
// recreate warnings (compose --dry-run) and `remove --purge` volume
// removal (slice-v1-down-volumes-named-list, -recreate-detection,
// -volume-auto-removal).

package e2e_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pt9912/u-boot/internal/adapter/driven/clock"
	confirmadapter "github.com/pt9912/u-boot/internal/adapter/driven/confirm"
	dockeradapter "github.com/pt9912/u-boot/internal/adapter/driven/docker"
	fsadapter "github.com/pt9912/u-boot/internal/adapter/driven/fs"
	"github.com/pt9912/u-boot/internal/adapter/driven/netprobe"
	yamladapter "github.com/pt9912/u-boot/internal/adapter/driven/yaml"
	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

const clusterComposeTmpl = `name: %s
services:
  web:
    image: busybox:1.36
    command: ["sleep", "%s"]
    volumes:
      - data:/data
volumes:
  data:
`

func writeClusterProject(t *testing.T, dir, project, sleep string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "u-boot.yaml"), []byte("schemaVersion: 1\nproject:\n  name: "+project+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "compose.yaml"), []byte(sprintfCompose(project, sleep)), 0o644); err != nil {
		t.Fatal(err)
	}
}

func sprintfCompose(project, sleep string) string {
	return strings.Replace(strings.Replace(clusterComposeTmpl, "%s", project, 1), "%s", sleep, 1)
}

func TestE2E_UpDownCluster_NamesRecreateAndDown(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker CLI not on PATH: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Minute)
	defer cancel()
	dir := t.TempDir()
	project := "uboot-e2e-cluster-" + strings.ToLower(filepath.Base(dir))
	project = strings.NewReplacer("_", "-", ".", "-").Replace(project)
	if len(project) > 40 {
		project = project[:40]
	}
	writeClusterProject(t, dir, project, "3600")

	engine := dockeradapter.NewEngine()
	t.Cleanup(func() {
		_ = engine.ComposeDown(context.Background(), dir, driven.ComposeDownOptions{RemoveVolumes: true})
	})
	fsys := fsadapter.New()
	upSvc := application.NewUpService(fsys, yamladapter.New(), engine, netprobe.New(), clock.New(), nil)
	downSvc := application.NewDownService(fsys, engine, confirmadapter.New(strings.NewReader(""), os.Stderr), nil)

	// First up: fresh stack → no recreate warning.
	resp, err := upSvc.Up(ctx, driving.UpRequest{BaseDir: dir, Timeout: 60 * time.Second, SilenceProgress: true})
	if err != nil {
		t.Fatalf("first up: %v", err)
	}
	for _, w := range resp.Warnings {
		t.Errorf("fresh stack must not warn: %+v", w)
	}

	// Changed command → Compose recreates the container: warning.
	writeClusterProject(t, dir, project, "3601")
	resp, err = upSvc.Up(ctx, driving.UpRequest{BaseDir: dir, Timeout: 60 * time.Second, SilenceProgress: true})
	if err != nil {
		t.Fatalf("second up: %v", err)
	}
	found := false
	for _, w := range resp.Warnings {
		found = found || (w.Code == "LH-FA-UP-003" && strings.Contains(w.Subject, "web"))
	}
	if !found {
		t.Errorf("recreate not announced; warnings = %+v", resp.Warnings)
	}

	// down --volumes lists the removed volume by its real name.
	down, err := downSvc.Down(ctx, driving.DownRequest{BaseDir: dir, RemoveVolumes: true, AssumeYes: true})
	if err != nil {
		t.Fatalf("down: %v", err)
	}
	if len(down.RemovedVolumeNames) != 1 || down.RemovedVolumeNames[0] != project+"_data" {
		t.Errorf("RemovedVolumeNames = %v, want [%s_data]", down.RemovedVolumeNames, project)
	}
}

func TestE2E_RemovePurgeRemovesVolume(t *testing.T) {
	res := runAcceptanceFlow(t, acceptanceFlow{
		projectName: "t-uboot-e2e-purge",
		serviceName: "postgres",
		envKeys:     []string{"POSTGRES_USER"},
		upTimeout:   90 * time.Second,
		ctxTimeout:  4 * time.Minute,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	engine := dockeradapter.NewEngine()

	before, err := engine.ListVolumeNames(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var target string
	for _, n := range before {
		if strings.HasSuffix(n, "_postgres-data") && strings.Contains(n, "t-uboot-e2e-purge") {
			target = n
		}
	}
	if target == "" {
		t.Fatalf("expected a postgres-data volume of the project; volumes: %v", before)
	}

	// The removed service's container stays behind as an orphan (its
	// compose block is gone) and would keep port 5432 busy for the next
	// test: clean it up explicitly.
	t.Cleanup(func() { forceRemoveProject("t-uboot-e2e-purge") })
	// Volume in use (stack still up): per-volume warning, not an error.
	fsys := fsadapter.New()
	removeSvc := application.NewRemoveServiceService(fsys, yamladapter.New(), confirmadapter.New(strings.NewReader(""), os.Stderr), nil).WithDockerEngine(engine)
	pg, _ := domain.NewServiceName("postgres")
	inUse, err := removeSvc.Remove(ctx, driving.RemoveServiceRequest{BaseDir: res.dir, ServiceName: pg, Purge: true, Yes: true})
	if err != nil {
		t.Fatalf("Remove (in use): %v", err)
	}
	if inUse.VolumesPurged || !inUse.PurgeAttempted || len(inUse.Warnings) == 0 {
		t.Errorf("in-use volume: purged=%v attempted=%v warnings=%+v, want a per-volume warning", inUse.VolumesPurged, inUse.PurgeAttempted, inUse.Warnings)
	}
}

func TestE2E_RemovePurgeRemovesVolume_StackDown(t *testing.T) {
	res := runAcceptanceFlow(t, acceptanceFlow{
		projectName: "t-uboot-e2e-purge-down",
		serviceName: "postgres",
		envKeys:     []string{"POSTGRES_USER"},
		upTimeout:   90 * time.Second,
		ctxTimeout:  4 * time.Minute,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	engine := dockeradapter.NewEngine()
	if err := engine.ComposeDown(ctx, res.dir, driven.ComposeDownOptions{}); err != nil {
		t.Fatalf("down (keep volumes): %v", err)
	}
	removeSvc := application.NewRemoveServiceService(fsadapter.New(), yamladapter.New(), confirmadapter.New(strings.NewReader(""), os.Stderr), nil).WithDockerEngine(engine)
	pg, _ := domain.NewServiceName("postgres")
	resp, err := removeSvc.Remove(ctx, driving.RemoveServiceRequest{BaseDir: res.dir, ServiceName: pg, Purge: true, Yes: true})
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if !resp.VolumesPurged || len(resp.PurgedVolumes) != 1 || !strings.HasSuffix(resp.PurgedVolumes[0], "_postgres-data") {
		t.Fatalf("purge result = purged:%v volumes:%v warnings:%+v", resp.VolumesPurged, resp.PurgedVolumes, resp.Warnings)
	}
	after, err := engine.ListVolumeNames(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, n := range after {
		if n == resp.PurgedVolumes[0] {
			t.Errorf("volume %s still exists after --purge", n)
		}
	}
}

// forceRemoveProject removes every container, volume and network of a
// Compose project by label (an orphaned stack cannot be `down`ed once
// its compose file declares no services).
func forceRemoveProject(project string) {
	filter := "label=com.docker.compose.project=" + project
	for _, kind := range [][]string{{"ps", "-aq"}, {"volume", "ls", "-q"}, {"network", "ls", "-q"}} {
		out, _ := exec.Command("docker", append(kind, "--filter", filter)...).Output()
		ids := strings.Fields(string(out))
		if len(ids) == 0 {
			continue
		}
		rm := map[string][]string{"ps": {"rm", "-f"}, "volume": {"volume", "rm", "-f"}, "network": {"network", "rm"}}[kind[0]]
		_ = exec.Command("docker", append(rm, ids...)...).Run()
	}
}
