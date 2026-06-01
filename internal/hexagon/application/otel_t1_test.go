package application_test

import (
	"bytes"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/pt9912/u-boot/internal/hexagon/application"
)

// slice-v1-otel T1: Templates + extraFiles-Catalogue-Erweiterung.
// Diese Tests pinnen den T1-Stand und die T2-Voraussetzung
// (isSupportedService("otel") bleibt false bis T2 executeAdd/Remove
// extraFiles-aware macht).

func TestOtelT1_TemplatesEmbedded(t *testing.T) {
	t.Parallel()
	names, err := application.TemplateNamesForTest()
	if err != nil {
		t.Fatalf("templateNames: %v", err)
	}
	wantPresent := []string{
		"services/otel.compose.tmpl",
		"services/otel.config.tmpl",
	}
	for _, want := range wantPresent {
		if !containsName(names, want) {
			t.Errorf("template %q missing from templateNames() = %v", want, names)
		}
	}
	if containsName(names, "services/otel.env.tmpl") {
		t.Errorf("services/otel.env.tmpl unexpectedly present — OTel-Catalogue setzt envTmpl=\"\" statt eines leeren Template-Files")
	}
	if containsName(names, "services/otel.volume.tmpl") {
		t.Errorf("services/otel.volume.tmpl unexpectedly present — OTel ist volumeOptional")
	}
}

func TestOtelT1_ServiceCatalogue_OtelEntry(t *testing.T) {
	t.Parallel()
	cat := application.ServiceCatalogueForTest()
	entry, ok := cat["otel"]
	if !ok {
		t.Fatal("serviceCatalogue() missing otel entry — Render-Catalogue ist T1-Voraussetzung")
	}
	if entry.ComposeTmpl != "services/otel.compose.tmpl" {
		t.Errorf("otel ComposeTmpl = %q, want %q", entry.ComposeTmpl, "services/otel.compose.tmpl")
	}
	if entry.EnvTmpl != "" {
		t.Errorf("otel EnvTmpl = %q, want \"\" (OTel hat keinen .env.example-Block)", entry.EnvTmpl)
	}
	if entry.VolumeTmpl != "" {
		t.Errorf("otel VolumeTmpl = %q, want \"\" (OTel ist volumeOptional)", entry.VolumeTmpl)
	}
	if len(entry.ExtraFiles) != 1 {
		t.Fatalf("otel ExtraFiles len = %d, want 1 (config-file)", len(entry.ExtraFiles))
	}
	if entry.ExtraFiles[0].Path != "otel-collector-config.yaml" {
		t.Errorf("otel ExtraFiles[0].Path = %q, want %q", entry.ExtraFiles[0].Path, "otel-collector-config.yaml")
	}
	if entry.ExtraFiles[0].Tmpl != "services/otel.config.tmpl" {
		t.Errorf("otel ExtraFiles[0].Tmpl = %q, want %q", entry.ExtraFiles[0].Tmpl, "services/otel.config.tmpl")
	}
}

func TestOtelT1_ServiceCatalogue_PostgresAndKeycloak_NoExtraFiles(t *testing.T) {
	t.Parallel()
	cat := application.ServiceCatalogueForTest()
	if len(cat["postgres"].ExtraFiles) != 0 {
		t.Errorf("postgres ExtraFiles = %v, want nil (T1 darf Postgres-Catalogue nicht ändern)", cat["postgres"].ExtraFiles)
	}
	if len(cat["keycloak"].ExtraFiles) != 0 {
		t.Errorf("keycloak ExtraFiles = %v, want nil (T1 darf Keycloak-Catalogue nicht ändern)", cat["keycloak"].ExtraFiles)
	}
}

func TestOtelT1_RenderOtel_HasComposeAndExtraFile_NoEnvNoVolume(t *testing.T) {
	t.Parallel()
	composeFrag, volumeFrag, envVars, extraFiles, err := application.RenderServiceTemplatesForTest(mustNewServiceName(t, "otel"))
	if err != nil {
		t.Fatalf("RenderServiceTemplates(otel): %v", err)
	}
	if len(composeFrag) == 0 {
		t.Error("otel ServiceFragment leer")
	}
	if !bytes.Contains(composeFrag, []byte("otel/opentelemetry-collector:0.108.0")) {
		t.Errorf("otel ServiceFragment ohne Stable-Image-Pin; got:\n%s", composeFrag)
	}
	if !bytes.Contains(composeFrag, []byte("4317:4317")) {
		t.Errorf("otel ServiceFragment ohne OTLP/gRPC-Port; got:\n%s", composeFrag)
	}
	if !bytes.Contains(composeFrag, []byte("4318:4318")) {
		t.Errorf("otel ServiceFragment ohne OTLP/HTTP-Port; got:\n%s", composeFrag)
	}
	if volumeFrag != nil {
		t.Errorf("otel VolumeFragment = %q, want nil (volumeOptional)", volumeFrag)
	}
	if envVars != nil {
		t.Errorf("otel EnvVariables = %q, want nil (kein .env-Block)", envVars)
	}
	if len(extraFiles) != 1 {
		t.Fatalf("otel ExtraFiles len = %d, want 1", len(extraFiles))
	}
	xf := extraFiles[0]
	if xf.Path != "otel-collector-config.yaml" {
		t.Errorf("ExtraFiles[0].Path = %q, want %q", xf.Path, "otel-collector-config.yaml")
	}
	// Content muss syntaktisch valides YAML sein — Round-Trip durch
	// yaml.v3-Unmarshal pinnt das.
	var parsed map[string]any
	if err := yaml.Unmarshal(xf.Content, &parsed); err != nil {
		t.Errorf("otel-collector-config.yaml ist kein gültiges YAML: %v\nContent:\n%s", err, xf.Content)
	}
	for _, want := range []string{"receivers", "processors", "exporters", "service"} {
		if _, ok := parsed[want]; !ok {
			t.Errorf("otel-collector-config.yaml fehlt Top-Level-Key %q; parsed keys = %v", want, parsed)
		}
	}
}

func TestOtelT1_RenderPostgresAndKeycloak_NoExtraFiles(t *testing.T) {
	t.Parallel()
	// Postgres + Keycloak werden NICHT mit extra files belastet —
	// T1 darf das nicht aus Versehen aktivieren.
	for _, svc := range []string{"postgres", "keycloak"} {
		_, _, _, extraFiles, err := application.RenderServiceTemplatesForTest(mustNewServiceName(t, svc))
		if err != nil {
			t.Fatalf("RenderServiceTemplates(%s): %v", svc, err)
		}
		if len(extraFiles) != 0 {
			t.Errorf("%s ExtraFiles = %v, want nil (Catalogue-Eintrag ohne extra files)", svc, extraFiles)
		}
	}
}

func TestOtelT2_IsSupportedService_AllThreeTrue(t *testing.T) {
	t.Parallel()
	// T2 hat die Catalogue erweitert: Postgres + Keycloak + OTel
	// sind jetzt alle supported.
	for _, svc := range []string{"postgres", "keycloak", "otel"} {
		if !application.IsSupportedServiceForTest(mustNewServiceName(t, svc)) {
			t.Errorf("isSupportedService(%s) muss nach T2 true sein", svc)
		}
	}
	if application.IsSupportedServiceForTest(mustNewServiceName(t, "ghost-service")) {
		t.Error("isSupportedService(ghost-service) muss false sein — Catalogue ist whitelist")
	}
}
