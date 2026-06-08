package application_test

import (
	"context"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// newConfigServiceWithLogger builds a ConfigService over a fresh
// fake FS + YAML and a capturing fakeLogger, seeded with the basic
// u-boot.yaml. Used by the T3 SilenceLogger / Warnings pins where
// the test needs to inspect the emitted log lines.
func newConfigServiceWithLogger(t *testing.T) (*application.ConfigService, *fakeLogger) {
	t.Helper()
	fs := newFakeFS()
	fs.markDirExists(configTestBaseDir)
	logger := &fakeLogger{}
	svc := application.NewConfigService(fs, &fakeYAML{}, logger)
	seedConfigUbootYAML(t, fs)
	return svc, logger
}

// TestConfigSet_SilenceLogger_SuppressesStderrLogs pins the T0-(n)
// JSON-mode logger silencing: with SilenceLogger=true, the
// "config set: updated" Info line (and every other Set-path log
// site) is routed to the no-op sink so the stderr stream stays
// clean for machine consumers.
func TestConfigSet_SilenceLogger_SuppressesStderrLogs(t *testing.T) {
	t.Parallel()
	svc, logger := newConfigServiceWithLogger(t)

	_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir:       configTestBaseDir,
		Path:          mustConfigPath(t, "project.name"),
		Value:         "renamed-project",
		SilenceLogger: true,
	})
	if err != nil {
		t.Fatalf("Set: %v", err)
	}
	if len(logger.entries) != 0 {
		t.Errorf("SilenceLogger=true should suppress all log sites; got %+v", logger.entries)
	}
}

// TestConfigSet_SilenceLoggerFalse_EmitsInfo pins the contrast: the
// default (SilenceLogger=false) keeps today's stderr Info line so
// the non-JSON CLI path is unchanged.
func TestConfigSet_SilenceLoggerFalse_EmitsInfo(t *testing.T) {
	t.Parallel()
	svc, logger := newConfigServiceWithLogger(t)

	_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "project.name"),
		Value:   "renamed-project",
	})
	if err != nil {
		t.Fatalf("Set: %v", err)
	}
	found := false
	for _, e := range logger.entries {
		if e.Level == "INFO" && e.Msg == "config set: updated" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("default SilenceLogger=false should emit `config set: updated` INFO; got %+v", logger.entries)
	}
}

// TestConfigSet_OrphanFeature_PopulatesWarnings pins the T0-(n)
// WARN-migration: activating a non-catalogued feature without a
// source override appends a [driving.WarningEntry] to the response
// so the CLI can map it into the JSON envelope's diagnostics[]
// without forcing the consumer to capture stderr. The Info log
// still fires too (dual emission).
func TestConfigSet_OrphanFeature_PopulatesWarnings(t *testing.T) {
	t.Parallel()
	svc, logger := newConfigServiceWithLogger(t)

	resp, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "devcontainer.features.unknown-thing.enabled"),
		Value:   "true",
	})
	if err != nil {
		t.Fatalf("Set: %v", err)
	}
	if len(resp.Warnings) != 1 {
		t.Fatalf("Warnings len = %d, want 1; got %+v", len(resp.Warnings), resp.Warnings)
	}
	w := resp.Warnings[0]
	if w.Code != "LH-FA-DEV-003" {
		t.Errorf("Warning.Code = %q, want LH-FA-DEV-003", w.Code)
	}
	if w.Level != "warn" {
		t.Errorf("Warning.Level = %q, want warn (Spec §1834 — warn|error only)", w.Level)
	}
	if w.Subject != "unknown-thing" {
		t.Errorf("Warning.Subject = %q, want feature name `unknown-thing`", w.Subject)
	}
	// Dual emission: stderr Info still fires when not silenced.
	infoFound := false
	for _, e := range logger.entries {
		if e.Level == "INFO" && strings.Contains(e.Msg, "orphan feature activation") {
			infoFound = true
			break
		}
	}
	if !infoFound {
		t.Errorf("expected orphan-feature INFO log alongside the WarningEntry; got %+v", logger.entries)
	}
}

// TestConfigSet_OrphanFeature_SilenceLogger_WarningsSurviveLogSuppression
// pins the key T0-(n) contract: SilenceLogger suppresses the stderr
// stream but MUST NOT drop the structured WARN — the consumer that
// asked for --json still sees the orphan warning via the envelope's
// diagnostics[], not via captured stderr.
func TestConfigSet_OrphanFeature_SilenceLogger_WarningsSurviveLogSuppression(t *testing.T) {
	t.Parallel()
	svc, logger := newConfigServiceWithLogger(t)

	resp, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir:       configTestBaseDir,
		Path:          mustConfigPath(t, "devcontainer.features.unknown-thing.enabled"),
		Value:         "true",
		SilenceLogger: true,
	})
	if err != nil {
		t.Fatalf("Set: %v", err)
	}
	if len(logger.entries) != 0 {
		t.Errorf("SilenceLogger=true should suppress the orphan Info log; got %+v", logger.entries)
	}
	if len(resp.Warnings) != 1 || resp.Warnings[0].Code != "LH-FA-DEV-003" {
		t.Errorf("WarningEntry must survive log suppression; got %+v", resp.Warnings)
	}
}

// TestConfigSet_NonOrphanScalar_NoWarnings pins that a routine
// scalar write (project.name) carries no Warnings — the field is
// nil on the happy path, so the CLI emits an empty diagnostics[].
func TestConfigSet_NonOrphanScalar_NoWarnings(t *testing.T) {
	t.Parallel()
	svc, _ := newConfigServiceWithLogger(t)

	resp, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir,
		Path:    mustConfigPath(t, "project.name"),
		Value:   "renamed-project",
	})
	if err != nil {
		t.Fatalf("Set: %v", err)
	}
	if resp.Warnings != nil {
		t.Errorf("Warnings = %+v, want nil on the non-orphan scalar path", resp.Warnings)
	}
}
