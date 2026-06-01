package application

import (
	"context"
	"errors"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// RemoveServiceService implements [driving.RemoveServiceUseCase] for
// the `u-boot remove <service>` flow (LH-FA-ADD-007). It is the
// mirror of [AddServiceService]: same dependencies (FileSystem,
// YAMLCodec, Confirmer, Logger), same state-machine layout
// (detect → execute), opposite direction.
//
// T1-skeleton scope: this file currently provides only the
// constructor, the static port-check, and a stub [Remove] that
// returns `errors.New("not yet implemented")`. T2 fills in the
// detect / execute pipeline against the existing
// [domain.ServiceState]-constants; T3 layers the `--purge`
// confirmation gate. The split keeps the T1 surface small enough
// for the CLI wiring slice (T4) to be reviewed independently.
type RemoveServiceService struct {
	fs        driven.FileSystem
	yaml      driven.YAMLCodec
	confirmer driven.Confirmer
	logger    driven.Logger
}

// Static check: RemoveServiceService satisfies the driving port.
var _ driving.RemoveServiceUseCase = (*RemoveServiceService)(nil)

// NewRemoveServiceService constructs the service with the driven
// adapters injected by the wiring layer. fs and yaml are mandatory
// (the constructor trusts the wiring layer to provide non-nil
// instances, matching the existing [NewAddServiceService] /
// [NewConfigService] pattern); confirmer and logger are nil-
// tolerant and fall back to the package-internal no-op
// implementations so callers (tests, scripts that do not care
// about prompts) need not wire a stub.
func NewRemoveServiceService(fs driven.FileSystem, yaml driven.YAMLCodec, confirmer driven.Confirmer, logger driven.Logger) *RemoveServiceService {
	if confirmer == nil {
		confirmer = noopConfirmer{}
	}
	if logger == nil {
		logger = noopLogger{}
	}
	return &RemoveServiceService{fs: fs, yaml: yaml, confirmer: confirmer, logger: logger}
}

// Remove is the T1 stub of [driving.RemoveServiceUseCase.Remove].
// It validates BaseDir and returns a "not-implemented" error for
// the actual state-machine path; T2 replaces the stub with the
// detect/execute flow against [domain.ServiceState] constants and
// the managed-block helpers used by M5 add.
//
// The validation here is the same shape T2 will keep — surfacing
// BaseDir-empty before the use case touches the filesystem so
// wiring errors are caught with a clear message regardless of
// fakeFS behaviour.
func (*RemoveServiceService) Remove(_ context.Context, req driving.RemoveServiceRequest) (driving.RemoveServiceResponse, error) {
	if req.BaseDir == "" {
		return driving.RemoveServiceResponse{}, errors.New("BaseDir is required")
	}
	return driving.RemoveServiceResponse{}, errors.New("remove: not yet implemented (T2 lands the state machine)")
}
