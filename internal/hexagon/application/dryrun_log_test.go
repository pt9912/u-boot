package application_test

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// pastTense liefert die Wörter, die ein eingetretenes Ergebnis melden; im
// Dry-Run darf keine Info-Zeile eines Schreibpfads eines davon tragen
// (LH-FA-CLI-007, slice-v1-cli-dry-run-log-would).
func pastTense() []string { return []string{"updated", "created", "repaired"} }

// dryRunLogCase fährt einen Schreibpfad mit dem gegebenen Vorschau-Modus
// und liefert den Logger, der alle Zeilen mitgeschnitten hat.
type dryRunLogCase struct {
	name  string
	done  string // Info-Zeile ohne Dry-Run (Vergangenheitsform)
	would string // Info-Zeile im Dry-Run
	args  string // Schlüssel und Werte beider Zeilen, als fmt.Sprint der Argumente
	run   func(t *testing.T, mode driving.PreviewMode) (*fakeLogger, error)
}

func infoMessages(l *fakeLogger) []string {
	var out []string
	for _, e := range l.entries {
		if e.Level == "INFO" {
			out = append(out, e.Msg)
		}
	}
	return out
}

// infoArgs liefert Schlüssel und Werte der ersten INFO-Zeile mit der Meldung.
func infoArgs(l *fakeLogger, msg string) (string, bool) {
	for _, e := range l.entries {
		if e.Level == "INFO" && e.Msg == msg {
			return fmt.Sprint(e.Args), true
		}
	}
	return "", false
}

func runConfigSet(t *testing.T, seed, path, value string, mode driving.PreviewMode) (*fakeLogger, error) {
	t.Helper()
	prod := newFakeFS()
	prod.markDirExists(configTestBaseDir)
	seedConfigUbootYAMLWithDevcontainer(t, prod, seed)
	logger := &fakeLogger{}
	svc := application.NewConfigServiceWithFactory(configFactoryForProd(prod), &fakeYAML{}, logger)
	_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
		BaseDir: configTestBaseDir, Path: mustConfigPath(t, path), Value: value, PreviewMode: mode,
	})
	return logger, err
}

func runGenerate(t *testing.T, seedYAML string, seedFiles map[string]string, req driving.GenerateRequest) (*fakeLogger, error) {
	t.Helper()
	prod := newFakeFS()
	prod.markDirExists(generateTestBaseDir)
	if err := prod.WriteFile(filepath.Join(generateTestBaseDir, "u-boot.yaml"), []byte(seedYAML), 0o644); err != nil {
		t.Fatal(err)
	}
	for rel, body := range seedFiles {
		if err := prod.WriteFile(filepath.Join(generateTestBaseDir, rel), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	logger := &fakeLogger{}
	svc := application.NewGenerateServiceWithFactory(configFactoryForProd(prod), &fakeYAML{}, logger)
	req.BaseDir = generateTestBaseDir
	_, err := svc.Generate(context.Background(), req)
	return logger, err
}

const dryRunGenSeed = "schemaVersion: 1\nproject:\n  name: t-uboot-gen\n"

func dryRunLogCases() []dryRunLogCase {
	changelogNoUnreleased := "<!-- BEGIN U-BOOT MANAGED BLOCK: init -->\n# Changelog\n\nUser-curated intro.\n" +
		"<!-- END U-BOOT MANAGED BLOCK: init -->\n\n## [0.1.0] - 2026-01-01\n\n### Added\n\n- initial release\n"
	staleEnv := "# BEGIN U-BOOT MANAGED BLOCK: init\n# stale\n# END U-BOOT MANAGED BLOCK: init\n"
	gen := func(art domain.Artifact, files map[string]string, extra func(*driving.GenerateRequest)) func(*testing.T, driving.PreviewMode) (*fakeLogger, error) {
		return func(t *testing.T, mode driving.PreviewMode) (*fakeLogger, error) {
			req := driving.GenerateRequest{Artifact: art, PreviewMode: mode}
			if extra != nil {
				extra(&req)
			}
			return runGenerate(t, dryRunGenSeed, files, req)
		}
	}
	return []dryRunLogCase{
		{"config set scalar", "config set: updated", "config set: would update",
			"[path project.name old t-uboot-config new renamed]",
			func(t *testing.T, m driving.PreviewMode) (*fakeLogger, error) {
				return runConfigSet(t, manyFixture, "project.name", "renamed", m)
			}},
		{"config set featureSources.allow", "config set: updated", "config set: would update",
			"[path devcontainer.featureSources.allow old https://ghcr.io/orgX/features/custom-rust new https://ghcr.io/orgX/features/custom-rust,https://example.test/x]",
			func(t *testing.T, m driving.PreviewMode) (*fakeLogger, error) {
				return runConfigSet(t, fixtureUBootYAMLDevcontainerSeed, "devcontainer.featureSources.allow", "https://example.test/x", m)
			}},
		{"config set egress.allow", "config set: updated", "config set: would update",
			"[path devcontainer.sandbox.egress.allow old  new a.example.com]",
			func(t *testing.T, m driving.PreviewMode) (*fakeLogger, error) {
				return runConfigSet(t, manyFixture, "devcontainer.sandbox.egress.allow", "a.example.com", m)
			}},
		{"config set many", "config set (many): updated", "config set (many): would update",
			"[items 2]",
			func(t *testing.T, m driving.PreviewMode) (*fakeLogger, error) {
				prod := newFakeFS()
				prod.markDirExists(configTestBaseDir)
				seedMany(t, prod, manyFixture)
				logger := &fakeLogger{}
				svc := application.NewConfigServiceWithFactory(configFactoryForProd(prod), &fakeYAML{}, logger)
				_, err := svc.SetMany(context.Background(), driving.ConfigSetManyRequest{BaseDir: configTestBaseDir,
					PreviewMode: m,
					Items: []driving.ConfigSetItem{
						{Path: mustConfigPath(t, "devcontainer.user.uid"), Value: "501"},
						{Path: mustConfigPath(t, "project.name"), Value: "renamed"},
					}})
				return logger, err
			}},
		{"generate changelog created", "generate: created", "generate: would create",
			"[artifact changelog path CHANGELOG.md project t-uboot-gen]",
			gen(domain.ArtifactChangelog, nil, nil)},
		{"generate readme created", "generate: created", "generate: would create",
			"[artifact readme path README.md project t-uboot-gen]",
			gen(domain.ArtifactReadme, nil, nil)},
		{"generate changelog repaired header", "generate: repaired Unreleased header", "generate: would repair Unreleased header",
			"[artifact changelog path CHANGELOG.md project t-uboot-gen]",
			gen(domain.ArtifactChangelog, map[string]string{"CHANGELOG.md": changelogNoUnreleased}, nil)},
		{"generate env-example updated block", "generate: updated block", "generate: would update block",
			"[artifact env-example path .env.example project t-uboot-gen]",
			gen(domain.ArtifactEnvExample, map[string]string{".env.example": staleEnv}, nil)},
		{"generate devcontainer created", "generate devcontainer: created", "generate devcontainer: would create",
			"[project t-uboot-gen forwardPorts [] features 0 changed [.devcontainer/devcontainer.json .devcontainer/Dockerfile]]",
			gen(domain.ArtifactDevcontainer, nil, nil)},
		{"generate devcontainer allowlist", "generate devcontainer: allowlist updated", "generate devcontainer: would update allowlist",
			"[added 1 total 1]",
			func(t *testing.T, m driving.PreviewMode) (*fakeLogger, error) {
				return gen(domain.ArtifactDevcontainer, nil, func(r *driving.GenerateRequest) {
					r.AllowExternalFeatureSources = []string{"https://example.test/features/x"}
				})(t, m)
			}},
	}
}

// Happy + Boundary: im Dry-Run steht die Zeile im Konjunktiv-Vokabular auf
// INFO und keine Info-Zeile trägt Vergangenheitsform; ohne Dry-Run und bei
// `--diff` allein (PreviewAndApply schreibt wirklich) bleibt die Zeile
// unverändert.
func TestDryRunLog_WriteLines(t *testing.T) {
	t.Parallel()
	for _, tc := range dryRunLogCases() {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			logger, err := tc.run(t, driving.PreviewDryRun)
			if err != nil {
				t.Fatalf("dry-run: %v", err)
			}
			msgs := infoMessages(logger)
			if !containsString(msgs, tc.would) {
				t.Errorf("dry-run INFO lines = %q, want %q", msgs, tc.would)
			}
			if containsString(msgs, tc.done) {
				t.Errorf("dry-run INFO carries past tense line %q: %q", tc.done, msgs)
			}
			if got, _ := infoArgs(logger, tc.would); got != tc.args {
				t.Errorf("dry-run %q fields = %s, want %s", tc.would, got, tc.args)
			}
			for _, m := range msgs {
				for _, w := range pastTense() {
					if strings.Contains(m, w) {
						t.Errorf("dry-run INFO line %q contains past tense %q", m, w)
					}
				}
			}

			for _, mode := range []driving.PreviewMode{driving.PreviewNone, driving.PreviewAndApply} {
				logger, err := tc.run(t, mode)
				if err != nil {
					t.Fatalf("mode %v: %v", mode, err)
				}
				msgs := infoMessages(logger)
				if !containsString(msgs, tc.done) {
					t.Errorf("mode %v INFO lines = %q, want %q", mode, msgs, tc.done)
				}
				if got, _ := infoArgs(logger, tc.done); got != tc.args {
					t.Errorf("mode %v %q fields = %s, want %s", mode, tc.done, got, tc.args)
				}
				if containsString(msgs, tc.would) {
					t.Errorf("mode %v INFO carries would line %q: %q", mode, tc.would, msgs)
				}
			}
		})
	}
}

// `generate devcontainer` meldet den Block-Austausch im Dry-Run ebenfalls im
// Konjunktiv: erst echt erzeugen, dann die Konfiguration ändern.
func TestDryRunLog_GenerateDevcontainerUpdatedBlock(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		mode driving.PreviewMode
		want string
	}{
		{driving.PreviewDryRun, "generate devcontainer: would update block"},
		{driving.PreviewNone, "generate devcontainer: updated-block"},
	} {
		prod := newFakeFS()
		prod.markDirExists(generateTestBaseDir)
		seedUBootYAMLPostgres(t, prod)
		seedGenerateComposeYAML(t, prod, composeYAMLPostgres)
		logger := &fakeLogger{}
		svc := application.NewGenerateServiceWithFactory(configFactoryForProd(prod), &fakeYAML{}, logger)
		req := driving.GenerateRequest{BaseDir: generateTestBaseDir, Artifact: domain.ArtifactDevcontainer}
		if _, err := svc.Generate(context.Background(), req); err != nil {
			t.Fatalf("seed run: %v", err)
		}
		seedUBootYAMLPostgresAndDummy(t, prod)
		seedGenerateComposeYAML(t, prod, composeYAMLPostgresAndDummy)
		logger.entries = nil
		req.PreviewMode = tc.mode
		if _, err := svc.Generate(context.Background(), req); err != nil {
			t.Fatalf("mode %v: %v", tc.mode, err)
		}
		if msgs := infoMessages(logger); !containsString(msgs, tc.want) {
			t.Errorf("mode %v INFO lines = %q, want %q", tc.mode, msgs, tc.want)
		}
	}
}

// Negative: der Dry-Run schreibt nichts, und No-op sowie Fehler melden
// weiterhin nichts auf INFO.
func TestDryRunLog_DryRunWritesNothingAndNoOpStaysQuiet(t *testing.T) {
	t.Parallel()
	prod := newFakeFS()
	prod.markDirExists(configTestBaseDir)
	seedConfigUbootYAMLWithDevcontainer(t, prod, manyFixture)
	before, _ := prod.ReadFile(filepath.Join(configTestBaseDir, "u-boot.yaml"))
	logger := &fakeLogger{}
	svc := application.NewConfigServiceWithFactory(configFactoryForProd(prod), &fakeYAML{}, logger)
	set := func(value string) error {
		_, err := svc.Set(context.Background(), driving.ConfigSetRequest{
			BaseDir: configTestBaseDir, Path: mustConfigPath(t, "project.name"), Value: value,
			PreviewMode: driving.PreviewDryRun,
		})
		return err
	}
	if err := set("renamed"); err != nil {
		t.Fatal(err)
	}
	after, _ := prod.ReadFile(filepath.Join(configTestBaseDir, "u-boot.yaml"))
	if string(before) != string(after) {
		t.Error("dry-run changed u-boot.yaml")
	}
	logger.entries = nil
	if err := set("t-uboot-config"); err != nil { // gleicher Wert: No-op
		t.Fatal(err)
	}
	if msgs := infoMessages(logger); len(msgs) != 0 {
		t.Errorf("no-op dry-run INFO lines = %q, want none", msgs)
	}
	if err := set("   "); err == nil { // Fehler: keine Erfolgszeile
		t.Fatal("expected an error for an invalid name")
	}
	if msgs := infoMessages(logger); len(msgs) != 0 {
		t.Errorf("failed dry-run INFO lines = %q, want none", msgs)
	}
}
