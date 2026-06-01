package application_test

import (
	"context"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// slice-v1-keycloak T2: Per-Service-Probe-Mechanismus + Catalogue-
// Erweiterung. Diese Tests pinnen den T2-Stand:
//
//  - hasRequiredEnvKeysFor / hasRequiredServiceFieldsFor lesen die
//    requiredEnvKeys aus dem Service-Catalogue.
//  - inspectVolumeArtefact + patchTargetsFor skippen für
//    volumeOptional=true (Keycloak).
//  - isSupportedService("keycloak") ist jetzt true.

// --- hasRequiredEnvKeysFor pro Service ----------------------------------

func TestKeycloakT2_HasRequiredEnvKeysFor_PostgresAllThree(t *testing.T) {
	t.Parallel()
	body := []byte("POSTGRES_USER=postgres\nPOSTGRES_PASSWORD=secret\nPOSTGRES_DB=mydb\n")
	if !application.HasRequiredEnvKeysForTest(mustNewServiceName(t, "postgres"), body) {
		t.Error("postgres env-Block mit allen drei Keys → erwartet true")
	}
}

func TestKeycloakT2_HasRequiredEnvKeysFor_PostgresMissingOne(t *testing.T) {
	t.Parallel()
	body := []byte("POSTGRES_USER=postgres\nPOSTGRES_DB=mydb\n")
	if application.HasRequiredEnvKeysForTest(mustNewServiceName(t, "postgres"), body) {
		t.Error("postgres env-Block ohne POSTGRES_PASSWORD → erwartet false")
	}
}

func TestKeycloakT2_HasRequiredEnvKeysFor_KeycloakBoth(t *testing.T) {
	t.Parallel()
	body := []byte("KEYCLOAK_ADMIN=admin\nKEYCLOAK_ADMIN_PASSWORD=secret\n")
	if !application.HasRequiredEnvKeysForTest(mustNewServiceName(t, "keycloak"), body) {
		t.Error("keycloak env-Block mit beiden Keys → erwartet true")
	}
}

func TestKeycloakT2_HasRequiredEnvKeysFor_KeycloakMissingPassword(t *testing.T) {
	t.Parallel()
	body := []byte("KEYCLOAK_ADMIN=admin\n")
	if application.HasRequiredEnvKeysForTest(mustNewServiceName(t, "keycloak"), body) {
		t.Error("keycloak env-Block ohne KEYCLOAK_ADMIN_PASSWORD → erwartet false")
	}
}

func TestKeycloakT2_HasRequiredEnvKeysFor_UnknownService(t *testing.T) {
	t.Parallel()
	body := []byte("FOO=1\nBAR=2\n")
	if application.HasRequiredEnvKeysForTest(mustNewServiceName(t, "ghost-service"), body) {
		t.Error("ghost-service (kein Catalogue-Eintrag) → erwartet false")
	}
}

// --- hasRequiredServiceFieldsFor pro Service ----------------------------

func TestKeycloakT2_HasRequiredServiceFieldsFor_PostgresComplete(t *testing.T) {
	t.Parallel()
	// Vollständiger Postgres-Service-Block inkl. Volume-Ref.
	body := []byte("" +
		"image: postgres:16-alpine\n" +
		"environment:\n" +
		"  POSTGRES_USER: postgres\n" +
		"  POSTGRES_PASSWORD: secret\n" +
		"  POSTGRES_DB: mydb\n" +
		"volumes:\n" +
		"  - postgres-data:/var/lib/postgresql/data\n" +
		"ports:\n" +
		"  - \"5432:5432\"\n" +
		"healthcheck:\n" +
		"  test: [\"CMD-SHELL\", \"pg_isready\"]\n")
	if !application.HasRequiredServiceFieldsForTest(mustNewServiceName(t, "postgres"), body) {
		t.Error("vollständiger Postgres-Block → erwartet true")
	}
}

func TestKeycloakT2_HasRequiredServiceFieldsFor_PostgresMissingVolume(t *testing.T) {
	t.Parallel()
	body := []byte("" +
		"image: postgres:16-alpine\n" +
		"environment:\n" +
		"  POSTGRES_USER: postgres\n" +
		"  POSTGRES_PASSWORD: secret\n" +
		"  POSTGRES_DB: mydb\n" +
		"ports:\n" +
		"  - \"5432:5432\"\n" +
		"healthcheck:\n" +
		"  test: [\"CMD-SHELL\", \"pg_isready\"]\n")
	if application.HasRequiredServiceFieldsForTest(mustNewServiceName(t, "postgres"), body) {
		t.Error("Postgres-Block ohne Volume-Ref → erwartet false (Postgres ist nicht volumeOptional)")
	}
}

func TestKeycloakT2_HasRequiredServiceFieldsFor_KeycloakNoVolume(t *testing.T) {
	t.Parallel()
	// Vollständiger Keycloak-Block — KEIN Volume.
	body := []byte("" +
		"image: quay.io/keycloak/keycloak:26.0\n" +
		"environment:\n" +
		"  KEYCLOAK_ADMIN: admin\n" +
		"  KEYCLOAK_ADMIN_PASSWORD: secret\n" +
		"ports:\n" +
		"  - \"8080:8080\"\n" +
		"healthcheck:\n" +
		"  test: [\"CMD-SHELL\", \"true\"]\n")
	if !application.HasRequiredServiceFieldsForTest(mustNewServiceName(t, "keycloak"), body) {
		t.Error("Keycloak-Block ohne Volume → erwartet true (volumeOptional=true)")
	}
}

func TestKeycloakT2_HasRequiredServiceFieldsFor_KeycloakMissingAdminPassword(t *testing.T) {
	t.Parallel()
	body := []byte("" +
		"image: quay.io/keycloak/keycloak:26.0\n" +
		"environment:\n" +
		"  KEYCLOAK_ADMIN: admin\n" +
		"ports:\n" +
		"  - \"8080:8080\"\n" +
		"healthcheck:\n" +
		"  test: [\"CMD-SHELL\", \"true\"]\n")
	if application.HasRequiredServiceFieldsForTest(mustNewServiceName(t, "keycloak"), body) {
		t.Error("Keycloak-Block ohne KEYCLOAK_ADMIN_PASSWORD → erwartet false")
	}
}

func TestKeycloakT2_HasRequiredServiceFieldsFor_KeycloakHealthcheckDisabled(t *testing.T) {
	t.Parallel()
	body := []byte("" +
		"image: quay.io/keycloak/keycloak:26.0\n" +
		"environment:\n" +
		"  KEYCLOAK_ADMIN: admin\n" +
		"  KEYCLOAK_ADMIN_PASSWORD: secret\n" +
		"ports:\n" +
		"  - \"8080:8080\"\n" +
		"healthcheck:\n" +
		"  disable: true\n")
	if application.HasRequiredServiceFieldsForTest(mustNewServiceName(t, "keycloak"), body) {
		t.Error("Keycloak-Block mit healthcheck.disable=true → erwartet false")
	}
}

// --- Integration: detectServiceState + Active-Repair pinnt Endlos-Loop ---

func TestKeycloakT2_AddTwice_NoRepairLoop(t *testing.T) {
	t.Parallel()
	// Regression-Pin gegen F1: nach erstem Add ist Keycloak Active;
	// zweiter Add darf NICHT in actionRepairArtifacts laufen, weil
	// volumes.keycloak-data per Design nie existiert.
	fs := newFakeFS()
	if err := fs.WriteFile("/proj/u-boot.yaml", []byte("schemaVersion: 1\nproject:\n  name: demo\n"), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	svc := application.NewAddServiceService(fs, &fakeYAML{}, nil, nil)

	req := driving.AddServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustNewServiceName(t, "keycloak"),
	}
	resp1, err := svc.Add(context.Background(), req)
	if err != nil {
		t.Fatalf("first Add: %v", err)
	}
	if resp1.State != domain.ServiceStateActive {
		t.Fatalf("first Add State = %s, want Active", resp1.State)
	}
	if len(resp1.Changed) == 0 {
		t.Error("first Add Changed empty — expected at least u-boot.yaml")
	}

	resp2, err := svc.Add(context.Background(), req)
	if err != nil {
		t.Fatalf("second Add: %v", err)
	}
	if resp2.State != domain.ServiceStateActive {
		t.Errorf("second Add State = %s, want Active (idempotent)", resp2.State)
	}
	if len(resp2.Changed) != 0 {
		t.Errorf("second Add Changed = %v, want nil (no-op, no repair-loop)", resp2.Changed)
	}
}
