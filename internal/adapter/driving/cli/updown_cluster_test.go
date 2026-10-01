package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

func clusterGetwd() (string, error) { return "/tmp/x/demo", nil }

// slice-v1-multi-port-services: `ports` is an additive structured array;
// `port` stays; no mapping → `[]`.
func TestUpJSON_StructuredPorts(t *testing.T) {
	uc := &fakeUpUseCase{resp: driving.UpResponse{Result: domain.UpResult{Stabilized: true, Services: []domain.ServiceStatus{
		{Name: "postgres", ContainerStatus: domain.StateRunning, Port: "5432:5432, 127.0.0.1:9091:9091", Healthcheck: "healthy"},
		{Name: "worker", ContainerStatus: domain.StateRunning},
	}}}}
	var stdout, stderr bytes.Buffer
	if err := newAppWithUp(uc, cli.WithGetwd(clusterGetwd)).Execute(context.Background(), []string{"--json", "up"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	var env struct {
		Data struct {
			Services []struct {
				Name  string   `json:"name"`
				Port  string   `json:"port"`
				Ports []string `json:"ports"`
			} `json:"services"`
		} `json:"data"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &env); err != nil {
		t.Fatalf("envelope: %v\n%s", err, stdout.String())
	}
	pg, w := env.Data.Services[0], env.Data.Services[1]
	if pg.Port != "5432:5432, 127.0.0.1:9091:9091" || len(pg.Ports) != 2 || pg.Ports[1] != "127.0.0.1:9091:9091" {
		t.Errorf("postgres = %+v", pg)
	}
	if w.Ports == nil || len(w.Ports) != 0 {
		t.Errorf("worker ports = %#v, want []", w.Ports)
	}
	if !strings.Contains(stdout.String(), `"ports":[]`) {
		t.Errorf("empty ports must serialize as [], got:\n%s", stdout.String())
	}
}

// slice-v1-recreate-detection: warnings reach diagnostics[] (JSON) and
// the human output; exit stays 0.
func TestUp_RecreateWarningsRendered(t *testing.T) {
	warn := driving.WarningEntry{Code: "LH-FA-UP-003", Level: "warn", Message: `container "x" will be recreated`, Subject: "x"}
	uc := &fakeUpUseCase{resp: driving.UpResponse{Warnings: []driving.WarningEntry{warn},
		Result: domain.UpResult{Stabilized: true, Services: []domain.ServiceStatus{{Name: "x", ContainerStatus: domain.StateRunning}}}}}
	var stdout, stderr bytes.Buffer
	if err := newAppWithUp(uc, cli.WithGetwd(clusterGetwd)).Execute(context.Background(), []string{"up"}, &stdout, &stderr); err != nil {
		t.Fatalf("human: %v", err)
	}
	if !strings.Contains(stdout.String(), "Warning (LH-FA-UP-003)") {
		t.Errorf("human output lacks the recreate warning:\n%s", stdout.String())
	}
	stdout.Reset()
	if err := newAppWithUp(uc, cli.WithGetwd(clusterGetwd)).Execute(context.Background(), []string{"--json", "up"}, &stdout, &stderr); err != nil {
		t.Fatalf("json: %v", err)
	}
	if !strings.Contains(stdout.String(), `"code":"LH-FA-UP-003"`) || !strings.Contains(stdout.String(), `"level":"warn"`) {
		t.Errorf("JSON lacks the warn diagnostic:\n%s", stdout.String())
	}
}

// slice-v1-up-partial-snapshot-on-failure: the error envelope carries
// data.services.
func TestUpJSON_PartialSnapshotOnError(t *testing.T) {
	uc := &fakeUpUseCase{
		resp: driving.UpResponse{PartialServices: []domain.ServiceStatus{{Name: "postgres", ContainerStatus: domain.StateDead}}},
		err:  errors.New("up service: stabilization timeout: " + driven.ErrComposeRuntime.Error()),
	}
	var stdout, stderr bytes.Buffer
	_ = newAppWithUp(uc, cli.WithGetwd(clusterGetwd)).Execute(context.Background(), []string{"--json", "up"}, &stdout, &stderr)
	var env struct {
		Status string `json:"status"`
		Data   struct {
			Services []struct{ Name, State string } `json:"services"`
		} `json:"data"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &env); err != nil {
		t.Fatalf("envelope: %v\n%s", err, stdout.String())
	}
	if env.Status != "error" || len(env.Data.Services) != 1 || env.Data.Services[0].Name != "postgres" || env.Data.Services[0].State != "dead" {
		t.Errorf("envelope = %+v\n%s", env, stdout.String())
	}
}

// slice-v1-down-volumes-named-list: removedVolumeNames is additive and
// always an array.
func TestDown_RemovedVolumeNamesRendered(t *testing.T) {
	uc := &fakeDownUseCase{resp: driving.DownResponse{RemovedVolumes: true, RemovedVolumeNames: []string{"demo_pg", "demo_x"}}}
	var stdout, stderr bytes.Buffer
	if err := newAppWithDown(uc, cli.WithGetwd(clusterGetwd)).Execute(context.Background(), []string{"down", "--volumes", "--yes"}, &stdout, &stderr); err != nil {
		t.Fatalf("human: %v", err)
	}
	if !strings.Contains(stdout.String(), "Removed volumes: demo_pg, demo_x") {
		t.Errorf("human output lacks the volume list:\n%s", stdout.String())
	}
	stdout.Reset()
	if err := newAppWithDown(uc, cli.WithGetwd(clusterGetwd)).Execute(context.Background(), []string{"--json", "down", "--volumes", "--yes"}, &stdout, &stderr); err != nil {
		t.Fatalf("json: %v", err)
	}
	if !strings.Contains(stdout.String(), `"removedVolumes":true`) || !strings.Contains(stdout.String(), `"removedVolumeNames":["demo_pg","demo_x"]`) {
		t.Errorf("JSON = %s", stdout.String())
	}
	stdout.Reset()
	empty := &fakeDownUseCase{resp: driving.DownResponse{}}
	if err := newAppWithDown(empty, cli.WithGetwd(clusterGetwd)).Execute(context.Background(), []string{"--json", "down"}, &stdout, &stderr); err != nil {
		t.Fatalf("json empty: %v", err)
	}
	if !strings.Contains(stdout.String(), `"removedVolumeNames":[]`) {
		t.Errorf("empty list must serialize as [], got %s", stdout.String())
	}
}

// slice-v1-volume-auto-removal: the purge result is rendered (human +
// JSON); the legacy "not yet automated" prose disappears once the real
// removal ran, failures show as warnings.
func TestRemove_PurgeResultRendered(t *testing.T) {
	pg, _ := domain.NewServiceName("postgres")
	ok := &fakeRemoveServiceUseCase{resp: driving.RemoveServiceResponse{
		ServiceName: pg, PriorState: domain.ServiceStateActive, State: domain.ServiceStateDeactivated,
		Changed: []string{"compose.yaml"}, VolumesPurged: true, PurgeAttempted: true, PurgedVolumes: []string{"demo_postgres-data"}}}
	var stdout, stderr bytes.Buffer
	if err := newAppWithRemove(ok, cli.WithGetwd(clusterGetwd)).Execute(context.Background(), []string{"remove", "postgres", "--purge", "--yes"}, &stdout, &stderr); err != nil {
		t.Fatalf("human: %v", err)
	}
	if !strings.Contains(stdout.String(), "Removed volumes: demo_postgres-data") || strings.Contains(stderr.String(), "NOT yet automated") {
		t.Errorf("stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
	stdout.Reset()
	if err := newAppWithRemove(ok, cli.WithGetwd(clusterGetwd)).Execute(context.Background(), []string{"--json", "remove", "postgres", "--purge", "--yes"}, &stdout, &stderr); err != nil {
		t.Fatalf("json: %v", err)
	}
	if !strings.Contains(stdout.String(), `"volumesPurged":true`) || !strings.Contains(stdout.String(), `"purgedVolumes":["demo_postgres-data"]`) {
		t.Errorf("JSON = %s", stdout.String())
	}

	warn := &fakeRemoveServiceUseCase{resp: driving.RemoveServiceResponse{
		ServiceName: pg, State: domain.ServiceStateDeactivated, Changed: []string{"compose.yaml"}, PurgeAttempted: true,
		Warnings: []driving.WarningEntry{{Code: "LH-FA-ADD-007", Level: "warn", Message: `--purge: volume "demo_postgres-data" could not be removed (in use)`}}}}
	stdout.Reset()
	stderr.Reset()
	if err := newAppWithRemove(warn, cli.WithGetwd(clusterGetwd)).Execute(context.Background(), []string{"remove", "postgres", "--purge", "--yes"}, &stdout, &stderr); err != nil {
		t.Fatalf("warn: %v", err)
	}
	if !strings.Contains(stderr.String(), "could not be removed (in use)") || strings.Contains(stderr.String(), "NOT yet automated") {
		t.Errorf("stderr = %q", stderr.String())
	}
}
