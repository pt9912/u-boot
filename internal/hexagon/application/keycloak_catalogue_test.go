package application_test

import (
	"bytes"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
)

// slice-v1-keycloak T1: Service-Catalogue + renderServiceTemplates +
// Keycloak-Templates landen im embed.FS. Diese Tests pinnen den
// T1-Stand und die T2-Voraussetzung („Keycloak ist noch nicht in
// der Catalogue").

func TestKeycloakT1_TemplatesEmbedded(t *testing.T) {
	t.Parallel()
	names, err := application.TemplateNamesForTest()
	if err != nil {
		t.Fatalf("templateNames: %v", err)
	}
	wantPresent := []string{
		"services/keycloak.compose.tmpl",
		"services/keycloak.env.tmpl",
	}
	for _, want := range wantPresent {
		if !containsName(names, want) {
			t.Errorf("template %q missing from templateNames() = %v", want, names)
		}
	}
	if containsName(names, "services/keycloak.volume.tmpl") {
		t.Errorf("services/keycloak.volume.tmpl unexpectedly present — Keycloak default-Persistenz ist H2-In-Container, kein Volume")
	}
}

func TestKeycloakT1_PostgresTemplatesUnchanged(t *testing.T) {
	t.Parallel()
	// Postgres bleibt mit allen drei Templates im embed.FS — T1
	// darf das M5-Pattern nicht aufweichen.
	names, err := application.TemplateNamesForTest()
	if err != nil {
		t.Fatalf("templateNames: %v", err)
	}
	for _, want := range []string{
		"services/postgres.compose.tmpl",
		"services/postgres.env.tmpl",
		"services/postgres.volume.tmpl",
	} {
		if !containsName(names, want) {
			t.Errorf("postgres template %q missing — T1 darf M5 nicht regressen", want)
		}
	}
}

func TestKeycloakT1_ServiceCatalogue_PostgresEntry(t *testing.T) {
	t.Parallel()
	cat := application.ServiceCatalogueForTest()
	entry, ok := cat["postgres"]
	if !ok {
		t.Fatal("serviceCatalogue() missing postgres entry")
	}
	if entry.ComposeTmpl != "services/postgres.compose.tmpl" {
		t.Errorf("postgres ComposeTmpl = %q, want %q", entry.ComposeTmpl, "services/postgres.compose.tmpl")
	}
	if entry.EnvTmpl != "services/postgres.env.tmpl" {
		t.Errorf("postgres EnvTmpl = %q, want %q", entry.EnvTmpl, "services/postgres.env.tmpl")
	}
	if entry.VolumeTmpl != "services/postgres.volume.tmpl" {
		t.Errorf("postgres VolumeTmpl = %q, want %q (Postgres ist volume-pflichtig)", entry.VolumeTmpl, "services/postgres.volume.tmpl")
	}
}

func TestKeycloakT1_ServiceCatalogue_KeycloakEntry(t *testing.T) {
	t.Parallel()
	cat := application.ServiceCatalogueForTest()
	entry, ok := cat["keycloak"]
	if !ok {
		t.Fatal("serviceCatalogue() missing keycloak entry — Keycloak-Render-Pfad ist T1-Voraussetzung")
	}
	if entry.ComposeTmpl != "services/keycloak.compose.tmpl" {
		t.Errorf("keycloak ComposeTmpl = %q, want %q", entry.ComposeTmpl, "services/keycloak.compose.tmpl")
	}
	if entry.EnvTmpl != "services/keycloak.env.tmpl" {
		t.Errorf("keycloak EnvTmpl = %q, want %q", entry.EnvTmpl, "services/keycloak.env.tmpl")
	}
	if entry.VolumeTmpl != "" {
		t.Errorf("keycloak VolumeTmpl = %q, want \"\" (Default-Persistenz embedded H2, kein Volume)", entry.VolumeTmpl)
	}
}

func TestKeycloakT1_RenderPostgres_ByteIdentity(t *testing.T) {
	t.Parallel()
	// Render-Refaktor darf das Postgres-Output Byte-Identity nicht
	// brechen. Vergleich gegen direktes renderTemplate für jedes
	// der drei Postgres-Files.
	composeFrag, volumeFrag, envVars, _, err := application.RenderServiceTemplatesForTest(mustNewServiceName(t,"postgres"))
	if err != nil {
		t.Fatalf("RenderServiceTemplates(postgres): %v", err)
	}
	wantCompose, err := application.RenderTemplateForTest("services/postgres.compose.tmpl", "")
	if err != nil {
		t.Fatalf("RenderTemplate(postgres.compose): %v", err)
	}
	if !bytes.Equal(composeFrag, wantCompose) {
		t.Errorf("postgres ServiceFragment Byte-Identity gebrochen")
	}
	wantVolume, err := application.RenderTemplateForTest("services/postgres.volume.tmpl", "")
	if err != nil {
		t.Fatalf("RenderTemplate(postgres.volume): %v", err)
	}
	if !bytes.Equal(volumeFrag, wantVolume) {
		t.Errorf("postgres VolumeFragment Byte-Identity gebrochen")
	}
	wantEnv, err := application.RenderTemplateForTest("services/postgres.env.tmpl", "")
	if err != nil {
		t.Fatalf("RenderTemplate(postgres.env): %v", err)
	}
	if !bytes.Equal(envVars, wantEnv) {
		t.Errorf("postgres EnvVariables Byte-Identity gebrochen")
	}
}

func TestKeycloakT1_RenderKeycloak_ProducesContentWithoutVolume(t *testing.T) {
	t.Parallel()
	composeFrag, volumeFrag, envVars, _, err := application.RenderServiceTemplatesForTest(mustNewServiceName(t,"keycloak"))
	if err != nil {
		t.Fatalf("RenderServiceTemplates(keycloak): %v", err)
	}
	if len(composeFrag) == 0 {
		t.Error("keycloak ServiceFragment leer")
	}
	if !bytes.Contains(composeFrag, []byte("quay.io/keycloak/keycloak:26.0")) {
		t.Errorf("keycloak ServiceFragment ohne LTS-Image-Pin; got:\n%s", composeFrag)
	}
	if !bytes.Contains(composeFrag, []byte("8080:8080")) {
		t.Errorf("keycloak ServiceFragment ohne Port-Mapping 8080; got:\n%s", composeFrag)
	}
	if volumeFrag != nil {
		t.Errorf("keycloak VolumeFragment = %q, want nil (kein Volume per T1-Catalogue-Eintrag)", volumeFrag)
	}
	if len(envVars) == 0 {
		t.Error("keycloak EnvVariables leer")
	}
	if !bytes.Contains(envVars, []byte("KEYCLOAK_ADMIN=CHANGEME_KEYCLOAK_ADMIN")) {
		t.Errorf("keycloak EnvVariables ohne KEYCLOAK_ADMIN-Placeholder; got:\n%s", envVars)
	}
	if !bytes.Contains(envVars, []byte("KEYCLOAK_ADMIN_PASSWORD=CHANGEME_KEYCLOAK_ADMIN_PASSWORD")) {
		t.Errorf("keycloak EnvVariables ohne KEYCLOAK_ADMIN_PASSWORD-Placeholder; got:\n%s", envVars)
	}
}

func TestKeycloakT2_IsSupportedService_BothTrue(t *testing.T) {
	t.Parallel()
	// T2 erweitert die Catalogue nach der Detect-Generalisierung —
	// jetzt sind Postgres UND Keycloak supported.
	if !application.IsSupportedServiceForTest(mustNewServiceName(t, "keycloak")) {
		t.Error("isSupportedService(keycloak) muss nach T2 true sein")
	}
	if !application.IsSupportedServiceForTest(mustNewServiceName(t, "postgres")) {
		t.Error("isSupportedService(postgres) regressed — T2 darf das nicht brechen")
	}
	// Eine nicht-katalogisierte Service-Name bleibt strikt rejected.
	if application.IsSupportedServiceForTest(mustNewServiceName(t, "ghost-service")) {
		t.Error("isSupportedService(ghost-service) muss false sein — Catalogue ist whitelist")
	}
}

func TestKeycloakT1_RenderServiceTemplates_UnknownService(t *testing.T) {
	t.Parallel()
	// Ein Service-Name, der nicht in der Catalogue steht, muss
	// einen klaren Fehler liefern statt zu crashen.
	_, _, _, _, err := application.RenderServiceTemplatesForTest(mustNewServiceName(t,"ghost-service"))
	if err == nil {
		t.Fatal("RenderServiceTemplates(ghost-service): expected error, got nil")
	}
}

// --- helpers --------------------------------------------------------------

func containsName(names []string, want string) bool {
	for _, n := range names {
		if n == want {
			return true
		}
	}
	return false
}
