package application

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/pt9912/u-boot/internal/hexagon/application/managedblock"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// detectActiveArtifacts is the pure classifier described in the
// M5-T4c slice plan: for an Active LH-FA-ADD-005 state, decide which
// LH-FA-ADD-002 artefacts (Compose service block, volume block,
// .env.example block) are missing or stale. Returns an
// activeArtifactsStatus with the per-artefact flags.
//
// Abort-class conditions (malformed marker, wrong anchor, user-
// manual entry without marker, parse error) bubble up as
// ErrServiceInconsistent — the classifier never silently overwrites
// such state. Read-only: makes no FS writes and accepts no plan
// data.
func (s *AddServiceService) detectActiveArtifacts(baseDir string, svc domain.ServiceName) (activeArtifactsStatus, error) {
	composePath := filepath.Join(baseDir, "compose.yaml")
	composeBody, _, exists, err := s.loadForPatch(composePath)
	if err != nil {
		// loadForPatch returns ErrBackupUnsupportedKind for symlinks
		// and non-regular files. Propagating it here gives the
		// classifier the same security contract as the patch-phase
		// load (no foreign-target reads for the LH-FA-ADD-002
		// content check).
		return activeArtifactsStatus{}, err
	}
	if !exists {
		// Active means the compose-block marker existed during
		// state-detection. A vanished compose.yaml here is a TOCTOU
		// race; treat as inconsistent so the user sees the repair
		// hint instead of a silent rewrite.
		return activeArtifactsStatus{}, fmt.Errorf(
			"%w: compose.yaml vanished between detectServiceState and detectActiveArtifacts",
			driving.ErrServiceInconsistent)
	}

	serviceStale, err := s.inspectServiceArtefact(composeBody, svc)
	if err != nil {
		return activeArtifactsStatus{}, err
	}

	volumeMissing, err := s.inspectVolumeArtefact(composeBody, svc)
	if err != nil {
		return activeArtifactsStatus{}, err
	}

	envMissingOrStale, err := s.inspectEnvArtefact(baseDir, svc)
	if err != nil {
		return activeArtifactsStatus{}, err
	}

	return activeArtifactsStatus{
		ServiceStale:      serviceStale,
		VolumeMissing:     volumeMissing,
		EnvMissingOrStale: envMissingOrStale,
	}, nil
}

// inspectServiceArtefact runs LocateMarkedEntry for the
// service.<svc>-block under services.<svc> and applies the symmetric
// anchor + content checks.
func (s *AddServiceService) inspectServiceArtefact(composeBody []byte, svc domain.ServiceName) (bool, error) {
	res, err := s.yaml.LocateMarkedEntry(composeBody, "services", svc.String(),
		serviceMarkerName(svc))
	if err != nil {
		return false, fmt.Errorf("%w: malformed service block for %q: %v",
			driving.ErrServiceInconsistent, svc.String(), err)
	}
	if res.MarkerSomewhereElse {
		return false, fmt.Errorf("%w: service.%s marker is not under services.%s",
			driving.ErrServiceInconsistent, svc.String(), svc.String())
	}
	if res.EntryExists && !res.MarkerInEntry {
		return false, fmt.Errorf("%w: services.%s exists but is not u-boot-managed",
			driving.ErrServiceInconsistent, svc.String())
	}
	if !res.MarkerInEntry {
		// Active state guarantees the marker exists; reaching here
		// would be a state-detection / detect-artefact disagreement.
		return false, fmt.Errorf("%w: service.%s marker disappeared between phases",
			driving.ErrServiceInconsistent, svc.String())
	}
	return !hasRequiredServiceFieldsFor(svc, res.BlockBody), nil
}

// inspectVolumeArtefact checks that the volume.<svc> marker hangs
// under volumes.<svc>-data. Missing is a repair flag; wrong anchor
// / user-manual / malformed is an abort.
//
// Services whose catalogue entry declares `volumeOptional: true`
// (slice-v1-keycloak T2 — Keycloak's flüchtige H2-In-Container-
// Persistenz) skip the probe entirely and report no-repair. Without
// this skip the post-T1 add-flow would loop in
// `actionRepairArtifacts` forever, since `volumes.keycloak-data` is
// by design never written.
func (s *AddServiceService) inspectVolumeArtefact(composeBody []byte, svc domain.ServiceName) (bool, error) {
	if entry, ok := catalogueFor(svc); ok && entry.volumeOptional {
		return false, nil
	}
	res, err := s.yaml.LocateMarkedEntry(composeBody, "volumes", volumeEntryKey(svc),
		volumeMarkerName(svc))
	if err != nil {
		return false, fmt.Errorf("%w: malformed volume block for %q: %v",
			driving.ErrServiceInconsistent, svc.String(), err)
	}
	if res.MarkerSomewhereElse {
		return false, fmt.Errorf("%w: volume.%s marker is not under volumes.%s",
			driving.ErrServiceInconsistent, svc.String(), volumeEntryKey(svc))
	}
	if res.EntryExists && !res.MarkerInEntry {
		return false, fmt.Errorf("%w: volumes.%s exists but is not u-boot-managed",
			driving.ErrServiceInconsistent, volumeEntryKey(svc))
	}
	return !res.MarkerInEntry, nil
}

// inspectEnvArtefact opens .env.example (if any) and reports whether
// the service block exists and contains the required env keys for
// the service (slice-v1-keycloak T2: per-service required-keys list
// looked up via [catalogueFor] instead of POSTGRES_* hardcoding).
// Missing file / missing block / missing required key all translate
// to needs-repair; malformed is an abort.
func (s *AddServiceService) inspectEnvArtefact(baseDir string, svc domain.ServiceName) (bool, error) {
	envPath := filepath.Join(baseDir, ".env.example")
	envBody, _, exists, err := s.loadForPatch(envPath)
	if err != nil {
		return false, err
	}
	if !exists {
		return true, nil
	}
	marker := managedblock.Marker{Style: managedblock.StyleHash, Name: serviceMarkerName(svc)}
	start, end, fErr := managedblock.Find(envBody, marker)
	switch {
	case errors.Is(fErr, managedblock.ErrBlockNotFound):
		return true, nil
	case fErr != nil:
		return false, fmt.Errorf("%w: malformed .env.example block for %q: %v",
			driving.ErrServiceInconsistent, svc.String(), fErr)
	}
	// Strip the BEGIN line (we already know start = BEGIN line start)
	// — scan only the body between the markers for the required keys.
	blockBody := extractEnvBlockBody(envBody, start, end)
	return !hasRequiredEnvKeysFor(svc, blockBody), nil
}

// extractEnvBlockBody returns the bytes between BEGIN and END marker
// lines (exclusive), without the marker lines themselves.
// managedblock.Find returns the byte range that includes both marker
// lines plus the terminator newline of END.
func extractEnvBlockBody(envBody []byte, start, end int) []byte {
	// advance past the BEGIN line
	bodyStart := start
	if nl := bytes.IndexByte(envBody[bodyStart:end], '\n'); nl != -1 {
		bodyStart += nl + 1
	}
	// retreat from end past the END line's trailing newline + the END line itself
	bodyEnd := end
	if bodyEnd > 0 && envBody[bodyEnd-1] == '\n' {
		bodyEnd--
	}
	// step back to the start of the END line
	for bodyEnd > bodyStart && envBody[bodyEnd-1] != '\n' {
		bodyEnd--
	}
	return envBody[bodyStart:bodyEnd]
}

// hasRequiredServiceFieldsFor runs the M5-T4c content-presence check
// on a service block body. The scan rules are unchanged (comment
// stripping, indent-stack block context, healthcheck.disable: true
// exception, trimmed-non-empty values); the required env-keys and
// volume-ref-literal come from the service's catalogue entry
// (slice-v1-keycloak T2). Unknown services collapse to „not
// complete", which on the Add path is dominated by isSupportedService.
func hasRequiredServiceFieldsFor(svc domain.ServiceName, blockBody []byte) bool {
	entry, ok := catalogueFor(svc)
	if !ok {
		return false
	}
	state := newContentScanState(entry)
	for _, raw := range bytes.Split(blockBody, []byte("\n")) {
		state.feedServiceLine(raw)
	}
	return state.serviceComplete()
}

// hasRequiredEnvKeysFor checks that the env block body contains
// non-commented `<KEY>=...` assignments for every key the service's
// catalogue entry declares as required (slice-v1-keycloak T2).
// Values may be anything. An unknown service (no catalogue entry)
// returns false — caller treats that as needs-repair, which on the
// Add path is dominated by the earlier `isSupportedService` reject;
// on remove the env probe is unused.
func hasRequiredEnvKeysFor(svc domain.ServiceName, blockBody []byte) bool {
	entry, ok := catalogueFor(svc)
	if !ok {
		return false
	}
	seen := make(map[string]bool, len(entry.requiredEnvKeys))
	for _, raw := range bytes.Split(blockBody, []byte("\n")) {
		trimmed := bytes.TrimSpace(raw)
		if len(trimmed) == 0 || trimmed[0] == '#' {
			continue
		}
		cleaned := stripInlineComment(trimmed)
		for _, key := range entry.requiredEnvKeys {
			if bytes.HasPrefix(cleaned, append([]byte(key), '=')) {
				seen[key] = true
			}
		}
	}
	for _, key := range entry.requiredEnvKeys {
		if !seen[key] {
			return false
		}
	}
	return true
}

// stripInlineComment removes any `# ...` tail starting at an
// unescaped `#` outside of quotes. Sufficient for the LH-FA-CONF /
// .env.example shapes the M5-T4c repair check looks at.
func stripInlineComment(line []byte) []byte {
	for i := 0; i < len(line); i++ {
		if line[i] == '#' {
			return bytes.TrimSpace(line[:i])
		}
	}
	return line
}

// contentScanState tracks the active sub-block (`environment`,
// `volumes`, `ports`, `healthcheck`) and per-field hits while walking
// a service block line by line. Parametrised by the service's
// catalogue entry (slice-v1-keycloak T2) — `requiredEnvKeys` drives
// the env-completeness check, `volumeRefLiteral` + `volumeOptional`
// drive the volume sub-block check.
type contentScanState struct {
	entry serviceCatalogueEntry

	hasImage            bool
	envSeen             map[string]bool
	hasVolumeRef        bool
	hasPortEntry        bool
	hasHealthcheckSub   bool
	healthcheckDisabled bool

	subBlock            string // "environment", "volumes", "ports", "healthcheck", or ""
	subBlockIndent      int
	subBlockEntryIndent int
}

func newContentScanState(entry serviceCatalogueEntry) *contentScanState {
	return &contentScanState{
		entry:               entry,
		envSeen:             make(map[string]bool, len(entry.requiredEnvKeys)),
		subBlockEntryIndent: -1,
	}
}

// feedServiceLine processes one line of the service block body and
// updates the state.
func (s *contentScanState) feedServiceLine(raw []byte) {
	indent := leadingSpaces(raw)
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || trimmed[0] == '#' {
		return
	}
	cleaned := stripInlineComment(trimmed)

	// Pop sub-block if indent dropped back to the service block's
	// own indent (or shallower).
	if s.subBlock != "" && indent <= s.subBlockIndent {
		s.subBlock = ""
		s.subBlockEntryIndent = -1
	}

	if s.subBlock == "" {
		s.enterPossibleSubBlock(cleaned, indent)
		return
	}
	s.feedSubBlockEntry(cleaned, indent)
}

// enterPossibleSubBlock recognises top-level service keys (image,
// environment, volumes, ports, healthcheck) and pushes sub-block
// context where needed.
func (s *contentScanState) enterPossibleSubBlock(cleaned []byte, indent int) {
	switch {
	case bytes.HasPrefix(cleaned, []byte("image:")):
		rest := bytes.TrimSpace(cleaned[len("image:"):])
		if len(rest) > 0 && !bytes.Equal(rest, []byte("\"\"")) && !bytes.Equal(rest, []byte("''")) {
			s.hasImage = true
		}
	case bytes.HasPrefix(cleaned, []byte("environment:")):
		s.subBlock = "environment"
		s.subBlockIndent = indent
	case bytes.HasPrefix(cleaned, []byte("volumes:")):
		s.subBlock = "volumes"
		s.subBlockIndent = indent
	case bytes.HasPrefix(cleaned, []byte("ports:")):
		s.subBlock = "ports"
		s.subBlockIndent = indent
	case bytes.HasPrefix(cleaned, []byte("healthcheck:")):
		s.subBlock = "healthcheck"
		s.subBlockIndent = indent
	}
}

// feedSubBlockEntry processes an indented line inside one of the
// recognised sub-blocks.
func (s *contentScanState) feedSubBlockEntry(cleaned []byte, indent int) {
	if s.subBlockEntryIndent == -1 {
		s.subBlockEntryIndent = indent
	}
	switch s.subBlock {
	case "environment":
		s.feedEnvironmentEntry(cleaned)
	case "volumes":
		if s.entry.volumeRefLiteral != "" && bytes.Contains(cleaned, []byte(s.entry.volumeRefLiteral)) {
			s.hasVolumeRef = true
		}
	case "ports":
		if bytes.HasPrefix(cleaned, []byte("- ")) || bytes.Equal(cleaned, []byte("-")) {
			s.hasPortEntry = true
		}
	case "healthcheck":
		s.feedHealthcheckEntry(cleaned)
	}
}

// feedEnvironmentEntry recognises any `<KEY>:` line where KEY is in
// the service's catalogue `requiredEnvKeys` and marks the matching
// flag. The value is irrelevant — user customisation of values is
// explicitly allowed.
func (s *contentScanState) feedEnvironmentEntry(cleaned []byte) {
	for _, key := range s.entry.requiredEnvKeys {
		if bytes.HasPrefix(cleaned, append([]byte(key), ':')) {
			s.envSeen[key] = true
			return
		}
	}
}

// feedHealthcheckEntry counts any sub-key as a valid healthcheck
// presence — except `disable: true`, which per Compose semantics
// turns the healthcheck off and therefore violates LH-AK-002.
func (s *contentScanState) feedHealthcheckEntry(cleaned []byte) {
	if bytes.HasPrefix(cleaned, []byte("disable:")) {
		rest := bytes.TrimSpace(cleaned[len("disable:"):])
		if bytes.Equal(rest, []byte("true")) {
			s.healthcheckDisabled = true
		}
		return
	}
	s.hasHealthcheckSub = true
}

// serviceComplete returns whether every LH-FA-ADD-002 / LH-AK-002
// required field is present. Env-key requirement is the catalogue
// `requiredEnvKeys` list (slice-v1-keycloak T2); the volume-ref
// requirement is skipped when the service is `volumeOptional`
// (Keycloak's flüchtige H2-In-Container-Persistenz).
func (s *contentScanState) serviceComplete() bool {
	if !s.hasImage {
		return false
	}
	for _, key := range s.entry.requiredEnvKeys {
		if !s.envSeen[key] {
			return false
		}
	}
	if !s.entry.volumeOptional && !s.hasVolumeRef {
		return false
	}
	if !s.hasPortEntry {
		return false
	}
	if !s.hasHealthcheckSub || s.healthcheckDisabled {
		return false
	}
	return true
}

// leadingSpaces counts the number of leading space (' ') or tab
// characters on the line.
func leadingSpaces(line []byte) int {
	for i := 0; i < len(line); i++ {
		if line[i] != ' ' && line[i] != '\t' {
			return i
		}
	}
	return len(line)
}

// volumeEntryKey returns the canonical compose volumes-map entry key
// for a service (postgres → "postgres-data"). Kept as a helper so
// future LH-FA-ADD-003/-004 services can register their own naming.
func volumeEntryKey(svc domain.ServiceName) string {
	return svc.String() + "-data"
}

// volumeMarkerName returns the canonical managed-block name for the
// volume artefact of a service ("volume.postgres" for postgres).
func volumeMarkerName(svc domain.ServiceName) string {
	return "volume." + svc.String()
}
