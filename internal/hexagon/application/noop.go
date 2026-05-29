package application

import (
	"context"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// noopProgress is the nil-tolerant default for the
// [driven.ProgressPort] dependency of services that accept a nil
// progress port at construction time — does nothing on every method,
// so the call sites stay free of nil checks.
type noopProgress struct{}

func (noopProgress) AffectedFiles(_ string, _ []driven.AffectedFile) {}

// noopConfirmer is the nil-tolerant default for the
// [driven.Confirmer] dependency — always declines, so a service
// constructed without a wired Confirmer behaves like a strict
// non-interactive run.
type noopConfirmer struct{}

func (noopConfirmer) ConfirmTreatAsExisting(_ context.Context, _ string, _ []string) (bool, error) {
	return false, nil
}

func (noopConfirmer) ConfirmRemoveVolumes(_ context.Context, _ string) (bool, error) {
	return false, nil
}

// noopLogger is the nil-tolerant default for the [driven.Logger]
// dependency — every level discards. Keeps the application services'
// debug/info call sites free of nil checks.
type noopLogger struct{}

func (noopLogger) Debug(_ string, _ ...any) {}
func (noopLogger) Info(_ string, _ ...any)  {}
func (noopLogger) Warn(_ string, _ ...any)  {}
func (noopLogger) Error(_ string, _ ...any) {}
