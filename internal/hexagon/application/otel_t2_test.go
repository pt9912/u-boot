package application_test

import (
	"context"
	"errors"
	iofs "io/fs"
	"sort"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// slice-v1-otel T2: executeAdd / executeRemove sind extraFiles-aware,
// Catalogue um otel erweitert. Diese Tests pinnen den T2-Stand:
//
//  - Add schreibt compose.yaml + u-boot.yaml + otel-collector-config.
//    yaml — kein .env.example (envTmpl="").
//  - Zweiter Add ist idempotent (kein Repair-Loop, F1-Schutz analog
//    Keycloak T2).
//  - Remove löscht den Compose-Block und die Collector-Config-Datei,
//    setzt enabled=false.

func TestOtelT2_Add_WritesComposeAndConfig_NoEnv(t *testing.T) {
	t.Parallel()
	fs := newFakeFS()
	if err := fs.WriteFile("/proj/u-boot.yaml", []byte("schemaVersion: 1\nproject:\n  name: demo\n"), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	svc := application.NewAddServiceService(fs, &fakeYAML{}, nil, nil)

	resp, err := svc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustNewServiceName(t, "otel"),
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if resp.State != domain.ServiceStateActive {
		t.Fatalf("Add State = %s, want Active", resp.State)
	}

	got := append([]string{}, resp.Changed...)
	sort.Strings(got)
	want := []string{"compose.yaml", "otel-collector-config.yaml", "u-boot.yaml"}
	sort.Strings(want)
	if !equalStrings(got, want) {
		t.Errorf("Changed = %v, want %v", got, want)
	}
	if containsString(resp.Changed, ".env.example") {
		t.Errorf("Changed enthält .env.example, sollte aber nicht (OTel-Catalogue envTmpl=\"\")")
	}

	// Verify config file ended up on disk with non-empty body.
	cfg, err := fs.ReadFile("/proj/otel-collector-config.yaml")
	if err != nil {
		t.Fatalf("read otel-collector-config.yaml: %v", err)
	}
	if !strings.Contains(string(cfg), "receivers:") || !strings.Contains(string(cfg), "exporters:") {
		t.Errorf("otel-collector-config.yaml fehlen Receivers/Exporters; got:\n%s", cfg)
	}
}

func TestOtelT2_AddTwice_NoRepairLoop(t *testing.T) {
	t.Parallel()
	fs := newFakeFS()
	if err := fs.WriteFile("/proj/u-boot.yaml", []byte("schemaVersion: 1\nproject:\n  name: demo\n"), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	svc := application.NewAddServiceService(fs, &fakeYAML{}, nil, nil)
	req := driving.AddServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustNewServiceName(t, "otel"),
	}

	if _, err := svc.Add(context.Background(), req); err != nil {
		t.Fatalf("first Add: %v", err)
	}
	resp2, err := svc.Add(context.Background(), req)
	if err != nil {
		t.Fatalf("second Add: %v", err)
	}
	if resp2.State != domain.ServiceStateActive {
		t.Errorf("second Add State = %s, want Active (idempotent)", resp2.State)
	}
	if len(resp2.Changed) != 0 {
		t.Errorf("second Add Changed = %v, want nil (no-op, kein Repair-Loop)", resp2.Changed)
	}
}

func TestOtelT2_Remove_DeletesComposeBlockAndConfig(t *testing.T) {
	t.Parallel()
	fs := newFakeFS()
	if err := fs.WriteFile("/proj/u-boot.yaml", []byte("schemaVersion: 1\nproject:\n  name: demo\n"), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	addSvc := application.NewAddServiceService(fs, &fakeYAML{}, nil, nil)
	if _, err := addSvc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustNewServiceName(t, "otel"),
	}); err != nil {
		t.Fatalf("Add: %v", err)
	}

	// Pre-condition: config-file exists.
	if _, err := fs.ReadFile("/proj/otel-collector-config.yaml"); err != nil {
		t.Fatalf("pre-condition: otel-collector-config.yaml fehlt nach Add: %v", err)
	}

	removeSvc := application.NewRemoveServiceService(fs, &fakeYAML{}, nil, nil)
	resp, err := removeSvc.Remove(context.Background(), driving.RemoveServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustNewServiceName(t, "otel"),
	})
	if err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if resp.State != domain.ServiceStateDeactivated {
		t.Errorf("Remove State = %s, want Deactivated", resp.State)
	}
	if !containsString(resp.Changed, "otel-collector-config.yaml") {
		t.Errorf("Changed = %v, want enthält otel-collector-config.yaml", resp.Changed)
	}
	if !containsString(resp.Changed, "compose.yaml") {
		t.Errorf("Changed = %v, want enthält compose.yaml", resp.Changed)
	}

	// Post-condition: config-file deleted.
	if _, err := fs.ReadFile("/proj/otel-collector-config.yaml"); err == nil {
		t.Error("otel-collector-config.yaml existiert noch nach Remove")
	} else if !isFakeNotExist(err) {
		t.Errorf("unexpected error reading deleted file: %v", err)
	}
}

func TestOtelT2_RemoveTwice_Idempotent(t *testing.T) {
	t.Parallel()
	fs := newFakeFS()
	if err := fs.WriteFile("/proj/u-boot.yaml", []byte("schemaVersion: 1\nproject:\n  name: demo\n"), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	addSvc := application.NewAddServiceService(fs, &fakeYAML{}, nil, nil)
	if _, err := addSvc.Add(context.Background(), driving.AddServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustNewServiceName(t, "otel"),
	}); err != nil {
		t.Fatalf("Add: %v", err)
	}

	removeSvc := application.NewRemoveServiceService(fs, &fakeYAML{}, nil, nil)
	req := driving.RemoveServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustNewServiceName(t, "otel"),
	}
	if _, err := removeSvc.Remove(context.Background(), req); err != nil {
		t.Fatalf("first Remove: %v", err)
	}
	// Zweiter Remove: Service ist deaktiviert; idempotent „nothing to do".
	resp, err := removeSvc.Remove(context.Background(), req)
	if err != nil {
		t.Fatalf("second Remove: %v", err)
	}
	if len(resp.Changed) != 0 {
		t.Errorf("second Remove Changed = %v, want nil (Service ist deactivated; extraFile ist schon weg)", resp.Changed)
	}
}

// --- helpers --------------------------------------------------------------

func containsString(s []string, want string) bool {
	for _, v := range s {
		if v == want {
			return true
		}
	}
	return false
}

func isFakeNotExist(err error) bool {
	return errors.Is(err, iofs.ErrNotExist)
}
