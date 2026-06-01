package application_test

import (
	"context"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

func TestRemoveServiceService_New(t *testing.T) {
	t.Parallel()
	// fs + yaml mandatory but not nil-checked at construction
	// (matches NewAddServiceService / NewConfigService pattern); the
	// wiring layer is trusted to provide them. confirmer and logger
	// are nil-tolerant.
	svc := application.NewRemoveServiceService(newFakeFS(), &fakeYAML{}, nil, nil)
	if svc == nil {
		t.Fatal("NewRemoveServiceService returned nil")
	}
}

func TestRemoveServiceService_Remove_EmptyBaseDirRejected(t *testing.T) {
	t.Parallel()
	svc := application.NewRemoveServiceService(newFakeFS(), &fakeYAML{}, nil, nil)
	name := mustServiceName(t, "postgres")

	_, err := svc.Remove(context.Background(), driving.RemoveServiceRequest{
		BaseDir:     "",
		ServiceName: name,
	})
	if err == nil {
		t.Fatal("Remove: want error for empty BaseDir, got nil")
	}
}

func TestRemoveServiceService_Remove_StubReturnsNotYetImplemented(t *testing.T) {
	t.Parallel()
	// T1-skeleton pin: the stub path returns a non-nil error so the
	// CLI wiring slice (T4) sees a clear failure if it lands before
	// T2 fills in the state machine. The exact error wording is not
	// part of the contract — pinning only the non-nil property.
	svc := application.NewRemoveServiceService(newFakeFS(), &fakeYAML{}, nil, nil)
	name := mustServiceName(t, "postgres")

	resp, err := svc.Remove(context.Background(), driving.RemoveServiceRequest{
		BaseDir:     "/proj",
		ServiceName: name,
	})
	if err == nil {
		t.Fatal("Remove: T1-skeleton must return an error until T2 lands")
	}
	if len(resp.Changed) != 0 || resp.VolumesPurged {
		t.Errorf("Remove: response carries side-effect state on error: %+v", resp)
	}
}

func mustServiceName(t *testing.T, raw string) domain.ServiceName {
	t.Helper()
	name, err := domain.NewServiceName(raw)
	if err != nil {
		t.Fatalf("NewServiceName(%q): %v", raw, err)
	}
	return name
}
