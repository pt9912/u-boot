package application

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	iofs "io/fs"
	"path/filepath"

	"github.com/pt9912/u-boot/internal/hexagon/application/managedblock"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// executePlan is the internal plan-phase output. Slots are nil when
// the corresponding file does not need to be written. baseDir lives
// outside the struct so the plan-data stays request-shape-free.
type executePlan struct {
	UBootYAML  *fileWrite
	Compose    *fileWrite
	EnvExample *fileWrite
}

// fileWrite is one pending write: the relative path, the
// pre-validated bytes, and the mode to write under.
type fileWrite struct {
	Path string
	Body []byte
	Mode iofs.FileMode
}

// executeAdd owns the full plan-phase and execute-phase for a single
// `u-boot add` invocation. T3's stub signature
// (`executeAdd(ctx, plan)`) was widened in T4c to accept baseDir
// directly so the plan-phase can run loadForPatch / PatchScalar /
// PatchMappingEntryYAML against the project tree without each
// helper receiving its own copy.
//
// Errors fall into two buckets:
//   - pre-write fail: any plan-phase error means no file is written.
//   - per-file write fail: bubbles up after partial writes; the
//     deterministic order (u-boot.yaml → compose.yaml →
//     .env.example) keeps the failure mode debuggable.
func (s *AddServiceService) executeAdd(_ context.Context, baseDir string, plan servicePlan) (driving.AddServiceResponse, error) {
	tmpls, err := s.renderServiceTemplates(plan.Service)
	if err != nil {
		return driving.AddServiceResponse{}, err
	}

	yamlBody, yamlMode, yamlConfig, err := s.loadAndParseUBootYAML(baseDir)
	if err != nil {
		return driving.AddServiceResponse{}, err
	}

	composeBody, composeMode, composeExists, err := s.loadOptional(baseDir, "compose.yaml")
	if err != nil {
		return driving.AddServiceResponse{}, err
	}
	envBody, envMode, envExists, err := s.loadOptional(baseDir, ".env.example")
	if err != nil {
		return driving.AddServiceResponse{}, err
	}

	composeBody, composeBootstrapped := bootstrapComposeIfEmpty(composeBody, yamlConfig.Project.Name)

	ep, err := s.buildExecutePlan(plan, tmpls,
		yamlBody, yamlMode,
		composeBody, composeMode, composeExists || composeBootstrapped,
		envBody, envMode, envExists,
	)
	if err != nil {
		return driving.AddServiceResponse{}, err
	}

	return s.runExecutePlan(baseDir, plan, ep)
}

// loadAndParseUBootYAML loads + parses u-boot.yaml. Missing file
// here is a TOCTOU race (state-detection already verified it
// existed) and maps to ErrProjectNotInitialized. Parse errors are
// reported as non-sentinel pre-write fails.
func (s *AddServiceService) loadAndParseUBootYAML(baseDir string) ([]byte, iofs.FileMode, ubootYAMLConfig, error) {
	yamlPath := filepath.Join(baseDir, "u-boot.yaml")
	body, mode, exists, err := s.loadForPatch(yamlPath)
	if err != nil {
		return nil, 0, ubootYAMLConfig{}, err
	}
	if !exists {
		return nil, 0, ubootYAMLConfig{}, fmt.Errorf(
			"%w: u-boot.yaml vanished between detectServiceState and loadForPatch",
			driving.ErrProjectNotInitialized)
	}
	var cfg ubootYAMLConfig
	if err := s.yaml.Unmarshal(body, &cfg); err != nil {
		return nil, 0, ubootYAMLConfig{}, fmt.Errorf("plan: parse u-boot.yaml: %w", err)
	}
	return body, mode, cfg, nil
}

// loadOptional loads a file that may legitimately not exist. Returns
// exists=false with mode=defaultFileMode for the missing case so the
// caller can plan a fresh write.
func (s *AddServiceService) loadOptional(baseDir, name string) ([]byte, iofs.FileMode, bool, error) {
	body, mode, exists, err := s.loadForPatch(filepath.Join(baseDir, name))
	if err != nil {
		return nil, 0, false, err
	}
	return body, mode, exists, nil
}

// loadForPatch is the M5-T4c read helper: Lstat first to reject
// symlinks and non-regular kinds with ErrBackupUnsupportedKind, then
// ReadFile. Missing files report exists=false with default mode so
// callers can route into create paths cleanly.
func (s *AddServiceService) loadForPatch(path string) ([]byte, iofs.FileMode, bool, error) {
	info, err := s.fs.Lstat(path)
	if err != nil {
		if errors.Is(err, iofs.ErrNotExist) {
			return nil, defaultFileMode, false, nil
		}
		return nil, 0, false, fmt.Errorf("lstat %s: %w", path, err)
	}
	if info.Mode()&iofs.ModeSymlink != 0 {
		return nil, 0, false, fmt.Errorf("%w: %s is a symlink",
			driving.ErrBackupUnsupportedKind, path)
	}
	if !info.Mode().IsRegular() {
		return nil, 0, false, fmt.Errorf("%w: %s is not a regular file",
			driving.ErrBackupUnsupportedKind, path)
	}
	body, err := s.fs.ReadFile(path)
	if err != nil {
		if errors.Is(err, iofs.ErrNotExist) {
			// TOCTOU: vanished between Lstat and ReadFile — treat as
			// missing.
			return nil, defaultFileMode, false, nil
		}
		return nil, 0, false, fmt.Errorf("read %s: %w", path, err)
	}
	return body, info.Mode().Perm(), true, nil
}

// renderedTemplates is the rendered template bytes used across the
// plan phase. Volume template is also a fragment; the env template is
// the raw variable body (the wrap helper adds the BEGIN/END markers
// during the env-edit step).
type renderedTemplates struct {
	ServiceFragment []byte
	VolumeFragment  []byte
	EnvVariables    []byte
}

// serviceCatalogueEntry is the per-service configuration shared by
// the render pipeline (slice-v1-keycloak T1: composeTmpl, envTmpl,
// volumeTmpl) and the detect/probe pipeline (slice-v1-keycloak T2:
// requiredEnvKeys, volumeRefLiteral, volumeOptional). Adding a new
// add-on means: drop a new entry here, drop the three templates into
// templates/services/, done — the rest of addservice_detect.go +
// removeservice.go consumes the entry declaratively.
type serviceCatalogueEntry struct {
	composeTmpl string
	envTmpl     string
	// volumeTmpl is the embedded volume-template path. Empty string
	// means: no volume managed block; volumeOptional must be true.
	volumeTmpl string
	// requiredEnvKeys are the .env.example keys the active-repair
	// detection looks for (LH-FA-ADD-002 / spec-equivalent for the
	// add-on). The service block's `environment:` sub-block scan
	// reuses the same list (lowercase / colon-suffixed variant).
	// Order is insignificant.
	requiredEnvKeys []string
	// volumeRefLiteral is the canonical compose volumes-map ref the
	// service block carries (e.g. "postgres-data" → `- postgres-
	// data:/var/lib/postgresql/data`). Empty when volumeOptional is
	// true.
	volumeRefLiteral string
	// volumeOptional disables the volume-artefact probe and the
	// volume-managed-block patch for services that do not need a
	// persistent named volume. Keycloak's default H2-In-Container-
	// Persistenz is the prototype.
	volumeOptional bool
}

// serviceCatalogue lists the per-service render + detect/probe
// configuration. Every service in `supportedServices()` must have
// an entry; the inverse direction (entry without supported-flag) is
// the slice-v1-keycloak T1-only-state preserved here for reference
// purposes (Keycloak landed in the catalogue with T1 to let
// renderServiceTemplates be reachable from tests, but
// `isSupportedService("keycloak")` only flipped to true with T2 once
// the detect pipeline became service-aware).
func serviceCatalogue() map[string]serviceCatalogueEntry {
	return map[string]serviceCatalogueEntry{
		"postgres": {
			composeTmpl:      "services/postgres.compose.tmpl",
			envTmpl:          "services/postgres.env.tmpl",
			volumeTmpl:       "services/postgres.volume.tmpl",
			requiredEnvKeys:  []string{"POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB"},
			volumeRefLiteral: "postgres-data",
			volumeOptional:   false,
		},
		"keycloak": {
			composeTmpl:      "services/keycloak.compose.tmpl",
			envTmpl:          "services/keycloak.env.tmpl",
			volumeTmpl:       "",
			requiredEnvKeys:  []string{"KEYCLOAK_ADMIN", "KEYCLOAK_ADMIN_PASSWORD"},
			volumeRefLiteral: "",
			volumeOptional:   true,
		},
	}
}

// catalogueFor returns the entry for the given service name. The
// boolean second return mirrors map-lookup convention so callers can
// branch on „unknown service" without panicking.
func catalogueFor(svc domain.ServiceName) (serviceCatalogueEntry, bool) {
	entry, ok := serviceCatalogue()[svc.String()]
	return entry, ok
}

// renderServiceTemplates evaluates the per-service templates via the
// [serviceCatalogue] lookup. The current templates are static (every
// service-specific name is hardcoded — see
// templates/services/<svc>.*.tmpl). templateData.Name is therefore
// passed empty, even though the init-tier templates require it; the
// T1-Sub-Decision in slice-v1-keycloak is to keep templates inline-
// hardcoded for now and defer `{{ .Name }}`-Parametrisierung +
// image-tag-overrides (e.g. `services.keycloak.imageTag`) to a
// follow-up slice.
//
// VolumeFragment is left nil when the service's catalogue entry has
// an empty `volumeTmpl`; callers must respect a nil VolumeFragment
// in the volume-patch step (slice-v1-keycloak T2 extends
// `patchTargetsFor` to skip the volume slot for volume-less
// services).
func (*AddServiceService) renderServiceTemplates(svc domain.ServiceName) (renderedTemplates, error) {
	entry, ok := serviceCatalogue()[svc.String()]
	if !ok {
		return renderedTemplates{}, fmt.Errorf("service catalogue: unknown service %q", svc.String())
	}
	data := templateData{}
	composeFrag, err := renderTemplate(entry.composeTmpl, data)
	if err != nil {
		return renderedTemplates{}, fmt.Errorf("plan: render %s: %w", entry.composeTmpl, err)
	}
	var volumeFrag []byte
	if entry.volumeTmpl != "" {
		volumeFrag, err = renderTemplate(entry.volumeTmpl, data)
		if err != nil {
			return renderedTemplates{}, fmt.Errorf("plan: render %s: %w", entry.volumeTmpl, err)
		}
	}
	envVars, err := renderTemplate(entry.envTmpl, data)
	if err != nil {
		return renderedTemplates{}, fmt.Errorf("plan: render %s: %w", entry.envTmpl, err)
	}
	return renderedTemplates{
		ServiceFragment: composeFrag,
		VolumeFragment:  volumeFrag,
		EnvVariables:    envVars,
	}, nil
}

// bootstrapComposeIfEmpty rewrites an empty or missing compose body
// into a minimal split-block scaffold so the patch-phase always sees
// a well-formed parent.
//
// "Empty" is strictly missing-file or all-whitespace. A user-edited
// compose without our init block is NOT a bootstrap trigger — the
// patch phase will splice into the existing structure instead.
func bootstrapComposeIfEmpty(composeBody []byte, projectName string) ([]byte, bool) {
	if len(bytes.TrimSpace(composeBody)) > 0 {
		return composeBody, false
	}
	scaffold := fmt.Sprintf(
		"# BEGIN U-BOOT MANAGED BLOCK: init\n"+
			"# Compose stack for %s.\n"+
			"\n"+
			"name: %s\n"+
			"\n"+
			"networks:\n"+
			"  default:\n"+
			"    name: %s-default\n"+
			"# END U-BOOT MANAGED BLOCK: init\n"+
			"\n"+
			"services: {}\n"+
			"\n"+
			"volumes: {}\n",
		projectName, projectName, projectName)
	return []byte(scaffold), true
}

// buildExecutePlan is the per-action core of the plan phase. It runs
// the symmetric anchor check (Pre-Patch-Anker-Check), fills the
// three optional fileWrite slots, and asserts at least one slot was
// set.
func (s *AddServiceService) buildExecutePlan(
	plan servicePlan,
	tmpls renderedTemplates,
	yamlBody []byte, yamlMode iofs.FileMode,
	composeBody []byte, composeMode iofs.FileMode, composeAvailable bool,
	envBody []byte, envMode iofs.FileMode, envExists bool,
) (executePlan, error) {
	ep := executePlan{}

	if err := s.preCheckComposeAnchors(plan, composeBody, composeAvailable); err != nil {
		return executePlan{}, err
	}

	if needsUBootYAMLSlot(plan) {
		patched, err := s.yaml.PatchScalar(yamlBody,
			[]string{"services", plan.Service.String(), "enabled"}, true)
		if err != nil {
			return executePlan{}, fmt.Errorf("plan: patch u-boot.yaml: %w", err)
		}
		ep.UBootYAML = &fileWrite{Path: "u-boot.yaml", Body: patched, Mode: yamlMode}
	}

	composeOut, composeChanged, err := s.planComposePatches(plan, composeBody, tmpls)
	if err != nil {
		return executePlan{}, err
	}
	if composeChanged {
		ep.Compose = &fileWrite{Path: "compose.yaml", Body: composeOut, Mode: composeMode}
	}

	envOut, envChanged, err := s.planEnvEdit(plan, envBody, envExists, tmpls)
	if err != nil {
		return executePlan{}, err
	}
	if envChanged {
		ep.EnvExample = &fileWrite{Path: ".env.example", Body: envOut, Mode: envMode}
	}

	if ep.UBootYAML == nil && ep.Compose == nil && ep.EnvExample == nil {
		return executePlan{}, fmt.Errorf(
			"plan: no slot populated for action %s on %s — programmer error",
			plan.Action.String(), plan.Service.String())
	}
	return ep, nil
}

// preCheckComposeAnchors runs the Pre-Patch-Anker-Check via
// LocateMarkedEntry for both compose markers (service.<svc> and
// volume.<svc>) regardless of which slots will be populated. The
// adapter's anchor validation also runs at patch time as defence in
// depth; this pre-check produces the fachlich correct
// ErrServiceInconsistent instead of the technical
// ErrYAMLAnchorMismatch.
func (s *AddServiceService) preCheckComposeAnchors(plan servicePlan, composeBody []byte, composeAvailable bool) error {
	if !composeAvailable {
		return nil
	}
	if err := s.preCheckAnchor(plan.Service, composeBody, "services", plan.Service.String(),
		serviceMarkerName(plan.Service)); err != nil {
		return err
	}
	return s.preCheckAnchor(plan.Service, composeBody, "volumes", volumeEntryKey(plan.Service),
		volumeMarkerName(plan.Service))
}

// preCheckAnchor branches the symmetric LocateResult table. Wrong-
// anchor / entry-without-marker both abort with
// ErrServiceInconsistent before any write.
func (s *AddServiceService) preCheckAnchor(svc domain.ServiceName, composeBody []byte, parentKey, entryKey, markerName string) error {
	res, err := s.yaml.LocateMarkedEntry(composeBody, parentKey, entryKey, markerName)
	if err != nil {
		return fmt.Errorf("%w: malformed %s block for %q: %v",
			driving.ErrServiceInconsistent, markerName, svc.String(), err)
	}
	if res.MarkerSomewhereElse {
		return fmt.Errorf(
			"%w: %s marker exists outside %s.%s; remove the orphan marker first",
			driving.ErrServiceInconsistent, markerName, parentKey, entryKey)
	}
	if res.EntryExists && !res.MarkerInEntry {
		return fmt.Errorf(
			"%w: %s.%s exists but is not u-boot-managed; remove or rename the manual entry "+
				"or set services.%s.enabled: false in u-boot.yaml",
			driving.ErrServiceInconsistent, parentKey, entryKey, svc.String())
	}
	return nil
}

// needsUBootYAMLSlot reports whether the action requires flipping
// the enabled scalar in u-boot.yaml. Register/Reactivate do;
// RebuildBlock/RepairArtifacts skip the slot because enabled is
// already true.
func needsUBootYAMLSlot(plan servicePlan) bool {
	switch plan.Action {
	case actionRegister, actionReactivate:
		return true
	case actionRebuildBlock, actionRepairArtifacts:
		return false
	}
	return false
}

// planComposePatches runs the per-action compose patches via
// PatchMappingEntryYAML. Returns the resulting body and a flag
// indicating whether anything changed (the Compose-slot stays nil
// when nothing changed).
func (s *AddServiceService) planComposePatches(plan servicePlan, composeBody []byte, tmpls renderedTemplates) ([]byte, bool, error) {
	patchService, patchVolume := patchTargetsFor(plan)
	if !patchService && !patchVolume {
		return composeBody, false, nil
	}

	out := composeBody
	if patchService {
		patched, err := s.yaml.PatchMappingEntryYAML(out, "services", plan.Service.String(),
			tmpls.ServiceFragment, serviceMarkerName(plan.Service))
		if err != nil {
			return nil, false, translatePatchErr(err, plan.Service, "services."+plan.Service.String())
		}
		out = patched
	}
	if patchVolume {
		patched, err := s.yaml.PatchMappingEntryYAML(out, "volumes", volumeEntryKey(plan.Service),
			tmpls.VolumeFragment, volumeMarkerName(plan.Service))
		if err != nil {
			return nil, false, translatePatchErr(err, plan.Service, "volumes."+volumeEntryKey(plan.Service))
		}
		out = patched
	}
	return out, true, nil
}

// patchTargetsFor returns whether the action needs to (re)write the
// service block and/or the volume block. Services whose catalogue
// entry declares `volumeOptional: true` (slice-v1-keycloak T2) never
// touch the volume slot — `renderServiceTemplates` returns a nil
// VolumeFragment for them anyway and a volume patch would either no-op
// or assert against nil.
func patchTargetsFor(plan servicePlan) (service, volume bool) {
	entry, ok := catalogueFor(plan.Service)
	volumeSkipped := ok && entry.volumeOptional
	switch plan.Action {
	case actionRegister, actionReactivate, actionRebuildBlock:
		return true, !volumeSkipped
	case actionRepairArtifacts:
		return plan.RepairFlags.ServiceStale, plan.RepairFlags.VolumeMissing && !volumeSkipped
	}
	return false, false
}

// translatePatchErr converts adapter-level patch sentinels into the
// fachlich-correct ErrServiceInconsistent / ErrYAMLFragmentInvalid
// surface so the CLI exit-code mapping stays clean.
func translatePatchErr(err error, svc domain.ServiceName, where string) error {
	switch {
	case errors.Is(err, driven.ErrYAMLAnchorMismatch):
		return fmt.Errorf("%w: patch of %s for %q rejected by adapter anchor check: %v",
			driving.ErrServiceInconsistent, where, svc.String(), err)
	case errors.Is(err, driven.ErrYAMLFragmentInvalid):
		return fmt.Errorf("%w: patch of %s for %q rejected by adapter fragment check: %v",
			driving.ErrServiceInconsistent, where, svc.String(), err)
	}
	return fmt.Errorf("plan: patch %s for %q: %w", where, svc.String(), err)
}

// planEnvEdit prepares the .env.example slot per the create / append
// / replace / malformed-abort strategy. Calls managedblock.Find
// exactly once per invocation so the malformed-detect path is
// unambiguous.
//
// Existing well-formed blocks whose body already contains the three
// required POSTGRES_* keys are left untouched (slot nil) so
// user-customised values survive. The same content-completeness
// rule the Active-Repair classifier uses (hasRequiredEnvKeys) is
// reused here — RebuildBlock with a complete env therefore returns
// Changed=["compose.yaml"] only, mirroring the slice-spec contract.
func (*AddServiceService) planEnvEdit(plan servicePlan, envBody []byte, envExists bool, tmpls renderedTemplates) ([]byte, bool, error) {
	if !envEditNeeded(plan) {
		return envBody, false, nil
	}
	block := renderEnvManagedBlock(plan.Service, tmpls.EnvVariables)
	if !envExists || len(bytes.TrimSpace(envBody)) == 0 {
		return block, true, nil
	}
	marker := managedblock.Marker{Style: managedblock.StyleHash, Name: serviceMarkerName(plan.Service)}
	start, end, fErr := managedblock.Find(envBody, marker)
	switch {
	case fErr == nil:
		existing := extractEnvBlockBody(envBody, start, end)
		if hasRequiredEnvKeysFor(plan.Service, existing) {
			// User-customised values survive — slot stays nil.
			return envBody, false, nil
		}
		updated, err := managedblock.Replace(envBody, marker, block)
		if err != nil {
			return nil, false, fmt.Errorf("%w: replace .env.example block for %q: %v",
				driving.ErrServiceInconsistent, plan.Service.String(), err)
		}
		return updated, true, nil
	case errors.Is(fErr, managedblock.ErrBlockNotFound):
		return appendEnvBlock(envBody, block), true, nil
	default:
		return nil, false, fmt.Errorf("%w: malformed .env.example block for %q: %v",
			driving.ErrServiceInconsistent, plan.Service.String(), fErr)
	}
}

// envEditNeeded reports whether the action touches .env.example.
// RepairArtifacts only touches it when the env flag is set; the
// others always do.
func envEditNeeded(plan servicePlan) bool {
	switch plan.Action {
	case actionRegister, actionReactivate, actionRebuildBlock:
		return true
	case actionRepairArtifacts:
		return plan.RepairFlags.EnvMissingOrStale
	}
	return false
}

// appendEnvBlock appends a managed env block to a non-empty .env.example
// body with exactly one separating blank line.
func appendEnvBlock(envBody, block []byte) []byte {
	var b bytes.Buffer
	b.Write(envBody)
	if len(envBody) > 0 && envBody[len(envBody)-1] != '\n' {
		b.WriteByte('\n')
	}
	b.WriteByte('\n')
	b.Write(block)
	return b.Bytes()
}

// renderEnvManagedBlock wraps a marker-free env variables body in
// BEGIN/END managed-block lines. All three env paths (create, append,
// replace) go through this wrapper so the produced block is
// byte-identical across paths; managedblock.Replace stays idempotent.
func renderEnvManagedBlock(svc domain.ServiceName, varsBody []byte) []byte {
	var b bytes.Buffer
	b.WriteString("# BEGIN U-BOOT MANAGED BLOCK: ")
	b.WriteString(serviceMarkerName(svc))
	b.WriteByte('\n')
	body := bytes.TrimRight(varsBody, "\n")
	b.Write(body)
	if len(body) > 0 {
		b.WriteByte('\n')
	}
	b.WriteString("# END U-BOOT MANAGED BLOCK: ")
	b.WriteString(serviceMarkerName(svc))
	b.WriteByte('\n')
	return b.Bytes()
}

// runExecutePlan writes the populated slots in the deterministic
// order u-boot.yaml → compose.yaml → .env.example.
func (s *AddServiceService) runExecutePlan(baseDir string, plan servicePlan, ep executePlan) (driving.AddServiceResponse, error) {
	var changed []string
	for _, w := range []*fileWrite{ep.UBootYAML, ep.Compose, ep.EnvExample} {
		if w == nil {
			continue
		}
		mode := w.Mode
		if mode == 0 {
			mode = defaultFileMode
		}
		if err := s.fs.WriteFile(filepath.Join(baseDir, w.Path), w.Body, mode); err != nil {
			return driving.AddServiceResponse{}, fmt.Errorf("write %s: %w", w.Path, err)
		}
		changed = append(changed, w.Path)
	}
	return driving.AddServiceResponse{
		ServiceName: plan.Service,
		PriorState:  plan.PriorState,
		State:       domain.ServiceStateActive,
		Changed:     changed,
	}, nil
}
