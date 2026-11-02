package application

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// The catalogue entry type "gepinnte Werkzeugkette" (LH-FA-DEV-003, lastenheft
// 0.30.0; SPEC-054): the keys `go-toolchain` and `golangci-lint` carry the
// toolchains of the u-boot build itself as Dockerfile blocks — release-
// archive installation with a SHA-256 check before the extraction, in both
// profiles. The version stands as a build argument in the catalogue template
// (the one pin place, LH-FA-DEV-011); a project overrides it per key with
// `version:` and the generate run reads the official checksum source of the
// key for the chosen version before anything is written (SPEC-054, S5).
// The conflict with an enabled external `go` feature and an override whose
// source is unreadable or does not cover the version end the run with a
// fachlicher Validierungsfehler (exit code 10, nothing written); callers
// wrap this file's errors into [driving.ErrConfigValueInvalid].

// goReleaseListingURL is the URL of the download listing of go.dev/dl with
// ALL versions (`include=all`): a `version:` override may name an older,
// unsupported release, and the listing carries its checksums with it. Same
// origin as the archive download of the block.
const goReleaseListingURL = "https://go.dev/dl/?mode=json&include=all"

// toolchainVersionCharset is the character set of a `version:` value of the
// two keys: release-version slugs ("1.27.3", "v2.14.0", "1.27rc1"). The value
// enters a download URL and the ARG block of the generated Dockerfile, so
// characters outside the set are unusable input (SPEC-054, S5: the input of
// the override is useless, exit code 10) and not the path of a shell escape.
var toolchainVersionCharset = regexp.MustCompile(`^[A-Za-z0-9._+-]+$`)

// goChecksumsFor reads the linux tar.gz checksums of `version` from the
// download listing of go.dev/dl. The ARG convention of the key carries the
// version without the leading `go` (root Dockerfile, `GO_VERSION`), so the
// read accepts both spellings and normalises to the bare form; the archive
// file names are `go<version>.linux-<arch>.tar.gz`. A version the listing
// does not name for linux/amd64 and linux/arm64 is the not-covered rejection.
func goChecksumsFor(listing []byte, version string) (toolchainPin, error) {
	version = strings.TrimPrefix(version, "go")
	if version == "" {
		return toolchainPin{}, fmt.Errorf("the version of features.go-toolchain is empty after the `go` prefix is trimmed")
	}
	amd64Hash, err := goListingChecksum(listing, version, "amd64")
	if err != nil {
		return toolchainPin{}, err
	}
	arm64Hash, err := goListingChecksum(listing, version, "arm64")
	if err != nil {
		return toolchainPin{}, err
	}
	return toolchainPin{version: version, sha256AMD64: amd64Hash, sha256ARM64: arm64Hash}, nil
}

// goListingChecksum extracts one archive checksum from the JSON listing; the
// filename and the kind/os/arch guards pin the row, not the position.
func goListingChecksum(listing []byte, version, arch string) (string, error) {
	var releases []goDownloadRelease
	if err := json.Unmarshal(listing, &releases); err != nil {
		return "", fmt.Errorf("the checksum source %s is not the version listing u-boot reads: %v", checksumSourceGo, err)
	}
	want := "go" + version + ".linux-" + arch + ".tar.gz"
	for _, release := range releases {
		for _, file := range release.Files {
			if file.Filename == want && file.Os == "linux" && file.Kind == "archive" {
				return file.SHA256, nil
			}
		}
	}
	return "", fmt.Errorf("the checksum source %s does not cover the version %s: no archive %s", checksumSourceGo, version, want)
}

// goDownloadListing is the shape of the `?mode=json` answer of go.dev/dl.
type goDownloadRelease struct {
	Version string           `json:"version"`
	Files   []goDownloadFile `json:"files"`
}

type goDownloadFile struct {
	Filename string `json:"filename"`
	Os       string `json:"os"`
	Arch     string `json:"arch"`
	SHA256   string `json:"sha256"`
	Kind     string `json:"kind"`
}

// golangciChecksumsURL is the URL of the checksums file of the release with
// the given tag: same origin as the release archive (the tag carries the
// leading `v`, the asset names drop it — golangci-lint-<version>-*.tar.gz).
func golangciChecksumsURL(tag string) string {
	base := strings.TrimPrefix(tag, "v")
	return "https://github.com/golangci/golangci-lint/releases/download/" + tag +
		"/golangci-lint-" + base + "-checksums.txt"
}

// golangciChecksumsFor reads the tar.gz checksums of both architectures from
// the CHECKSUMS file of the release with the tag. The lines of the file pair
// a SHA-256 and the file name in whitespace form; the read matches the exact
// file name of the archive, not the line order.
func golangciChecksumsFor(body []byte, tag string) (toolchainPin, error) {
	base := strings.TrimPrefix(tag, "v")
	amd64Hash, err := checksumsFileChecksum(body, checksumSourceGolangci, "golangci-lint-"+base+"-linux-amd64.tar.gz")
	if err != nil {
		return toolchainPin{}, err
	}
	arm64Hash, err := checksumsFileChecksum(body, checksumSourceGolangci, "golangci-lint-"+base+"-linux-arm64.tar.gz")
	if err != nil {
		return toolchainPin{}, err
	}
	return toolchainPin{version: tag, sha256AMD64: amd64Hash, sha256ARM64: arm64Hash}, nil
}

// checksumsFileChecksum extracts the checksum of one file from a checksums
// file in the `<sha256>  <filename>` line form; the source names the error
// (each key carries its own file form).
func checksumsFileChecksum(body []byte, source, filename string) (string, error) {
	for _, line := range strings.Split(string(body), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 || fields[1] != filename {
			continue
		}
		return fields[0], nil
	}
	return "", fmt.Errorf("the checksum source %s does not cover the file %s", source, filename)
}

// javaChecksumsURL is the URL of the Temurin checksums file of one
// architecture for the given version: the release tag carries the leading
// `jdk-` (the `+` of the version is URL-encoded in the path) and the
// checksums file is an asset of the same release the archive comes from —
// same origin (SPEC-054). The asset naming follows the Temurin convention
// `OpenJDK21U-jdk_<slug>_linux_hotspot_<version with `_` for `+`>.tar.gz`;
// the slug is `x64` (amd64) or `aarch64` (arm64).
func javaChecksumsURL(version, slug string) string {
	v := strings.TrimPrefix(version, "jdk-")
	tag := "jdk-" + strings.ReplaceAll(v, "+", "%2B")
	return "https://github.com/adoptium/temurin21-binaries/releases/download/" + tag +
		"/OpenJDK21U-jdk_" + slug + "_linux_hotspot_" + strings.ReplaceAll(v, "+", "_") + ".tar.gz.sha256.txt"
}

// javaChecksumsFor reads the checksums of both architectures from the
// checksums files of the Temurin release; the files pair a SHA-256 and the
// archive file name in the `<sha256>  <filename>` line form (same helper as
// the golangci key), matched on the exact archive file name.
func javaChecksumsFor(ctx context.Context, fetcher driven.Fetcher, version string) (toolchainPin, error) {
	v := strings.TrimPrefix(version, "jdk-")
	slug := strings.ReplaceAll(v, "+", "_")
	x64Hash, err := fetchJdkChecksum(ctx, fetcher, version, "x64", "OpenJDK21U-jdk_x64_linux_hotspot_"+slug+".tar.gz")
	if err != nil {
		return toolchainPin{}, err
	}
	arm64Hash, err := fetchJdkChecksum(ctx, fetcher, version, "aarch64", "OpenJDK21U-jdk_aarch64_linux_hotspot_"+slug+".tar.gz")
	if err != nil {
		return toolchainPin{}, err
	}
	return toolchainPin{version: v, sha256AMD64: x64Hash, sha256ARM64: arm64Hash}, nil
}

// fetchJdkChecksum fetches one checksums file of the release and extracts the
// checksum of the named archive.
func fetchJdkChecksum(ctx context.Context, fetcher driven.Fetcher, version, slug, archive string) (string, error) {
	body, err := fetchChecksums(ctx, fetcher, javaChecksumsURL(version, slug), checksumSourceJava)
	if err != nil {
		return "", err
	}
	return checksumsFileChecksum(body, checksumSourceJava, archive)
}

// toolchainPin is the render state of one `version:` override: the version in
// the ARG convention of its key (the Go key without, the golangci-lint key
// with the leading `v`, the java-toolchain key as the bare Temurin version
// with `+`), the two checksums read from the official source and that source
// as the message wording. The zero value means default pin.
type toolchainPin struct {
	version     string
	sha256AMD64 string
	sha256ARM64 string
}

// The exported accessor methods are the template form of the fields: the
// devcontainer templates read the pin of an override, and text/template
// cannot read unexported fields.
func (p toolchainPin) Version() string     { return p.version }
func (p toolchainPin) SHA256AMD64() string { return p.sha256AMD64 }
func (p toolchainPin) SHA256ARM64() string { return p.sha256ARM64 }

// toolchainState tells the three templates which blocks to render and which
// of them carries an override instead of the pin of the template.
type toolchainState struct {
	goToolchain   bool
	golangciLint  bool
	javaToolchain bool
	goPin         toolchainPin
	lintPin       toolchainPin
	javaPin       toolchainPin
}

// toolchainChecksumSourceSupported reports whether the catalogue entry is of
// the pinned-toolchain kind: a Dockerfile block whose entry carries one of
// the three official checksum sources. It is the type boundary of the
// `version:` override (SPEC-054): an entry of the type without a checksum
// source keeps the plain rejection, and so does every entry outside the type.
func toolchainChecksumSourceSupported(cat featureCatalogueEntry) bool {
	return cat.dockerfileBlock && (cat.checksumSource == checksumSourceGo || cat.checksumSource == checksumSourceGolangci || cat.checksumSource == checksumSourceJava)
}

// goFeatureActive reports whether the external `go` feature is enabled in dc,
// with or without a `source:` override: any enabled entry at the key installs
// a Go toolchain (the catalogue one or an external one), and which of two
// installations the commands find is a path-order question no test carries —
// the combination is the S8 rejection (SPEC-054).
func goFeatureActive(dc *ubootYAMLDevcontainer) bool {
	if dc == nil {
		return false
	}
	f, ok := dc.Features["go"]
	return ok && f.Enabled != nil && *f.Enabled
}

// goToolchainConflictMessage is the S8 rejection of an enabled external `go`
// feature next to the active `go-toolchain` key. Empty when the two go
// separate ways.
func goToolchainConflictMessage(dc *ubootYAMLDevcontainer) string {
	if !goFeatureActive(dc) || !dockerfileBlockFeatureActive(dc, featureGoToolchain) {
		return ""
	}
	return "devcontainer.features.go and devcontainer.features.go-toolchain install two Go toolchains in one image; " +
		"which one the commands find is an order of PATH and is not supported — " +
		"switch `go-toolchain` off with `u-boot config set devcontainer.features.go-toolchain.enabled false` " +
		"or the external `go` feature with `u-boot config set devcontainer.features.go.enabled false`"
}

// javaFeatureActive reports whether the external `java` feature is enabled in
// dc (= [goFeatureActive] for the java key): any enabled entry at the key
// installs a JDK, and which of two installations the commands find is a
// path-order question no test carries — the combination with the active
// `java-toolchain` key is the exit-10 rejection of SPEC-054 (Muster Go-Fall).
func javaFeatureActive(dc *ubootYAMLDevcontainer) bool {
	if dc == nil {
		return false
	}
	f, ok := dc.Features["java"]
	return ok && f.Enabled != nil && *f.Enabled
}

// javaToolchainConflictMessage is the exit-10 rejection of an enabled
// external `java` feature next to the active `java-toolchain` key, wording
// like the go case (SPEC-054). Empty when the two go separate ways.
func javaToolchainConflictMessage(dc *ubootYAMLDevcontainer) string {
	if !javaFeatureActive(dc) || !dockerfileBlockFeatureActive(dc, featureJavaToolchain) {
		return ""
	}
	return "devcontainer.features.java and devcontainer.features.java-toolchain install two JDKs in one image; " +
		"which one the commands find is an order of PATH and is not supported — " +
		"switch `java-toolchain` off with `u-boot config set devcontainer.features.java-toolchain.enabled false` " +
		"or the external `java` feature with `u-boot config set devcontainer.features.java.enabled false`"
}

// resolveToolchain validates the rejections of SPEC-054 and reads the
// official checksum source of every active key with a `version:` override.
// Without an override the fetcher is not asked: the default pin and its
// checksums stand in the template (one place, LH-FA-DEV-011). The error
// messages are the texts of the exit-10 rejections; the caller wraps them
// into [driving.ErrConfigValueInvalid] before anything is written.
func resolveToolchain(ctx context.Context, dc *ubootYAMLDevcontainer, fetcher driven.Fetcher) (toolchainState, error) {
	var state toolchainState
	if dc == nil {
		return state, nil
	}
	if msg := goToolchainConflictMessage(dc); msg != "" {
		return state, fmt.Errorf("%s", msg)
	}
	if msg := javaToolchainConflictMessage(dc); msg != "" {
		return state, fmt.Errorf("%s", msg)
	}
	state.goToolchain = dockerfileBlockFeatureActive(dc, featureGoToolchain)
	state.golangciLint = dockerfileBlockFeatureActive(dc, featureGolangciLint)
	state.javaToolchain = dockerfileBlockFeatureActive(dc, featureJavaToolchain)
	type toolchainKey struct {
		key     string
		cat     featureCatalogueEntry
		pinSlot *toolchainPin
	}
	for _, tk := range []toolchainKey{
		{key: featureGoToolchain, cat: featureCatalogue()[featureGoToolchain], pinSlot: &state.goPin},
		{key: featureGolangciLint, cat: featureCatalogue()[featureGolangciLint], pinSlot: &state.lintPin},
		{key: featureJavaToolchain, cat: featureCatalogue()[featureJavaToolchain], pinSlot: &state.javaPin},
	} {
		version := dc.Features[tk.key].Version
		if version == "" || !dockerfileBlockFeatureActive(dc, tk.key) {
			continue
		}
		if !toolchainChecksumSourceSupported(tk.cat) {
			return toolchainState{}, fmt.Errorf("devcontainer.features.%s.version: the entry has no official checksum source; an entry of the pinned toolchain takes a version override only with its source", tk.key)
		}
		value, err := normalizeToolchainVersion(tk.key, version, tk.cat.checksumSource)
		if err != nil {
			return toolchainState{}, err
		}
		if err := readToolchainChecksums(ctx, fetcher, tk.cat.checksumSource, tk.pinSlot, value); err != nil {
			return toolchainState{}, err
		}
	}
	return state, nil
}

// normalizeToolchainVersion validates the character set and normalises the
// value to the ARG convention of the key (the Go key without the leading
// `go`, the golangci-lint key with the leading `v`); the convention keeps the
// download URLs of the template block correct for both spellings.
func normalizeToolchainVersion(key, value, source string) (string, error) {
	if !toolchainVersionCharset.MatchString(value) {
		return "", fmt.Errorf("devcontainer.features.%s.version: the value %q is not a release version u-boot can read from the checksum source %s", key, value, source)
	}
	switch key {
	case featureGoToolchain:
		value = strings.TrimPrefix(value, "go")
	case featureGolangciLint:
		value = "v" + strings.TrimPrefix(value, "v")
	case featureJavaToolchain:
		// ARG form of the java-toolchain key: the bare Temurin version
		// with `+` as the build number separator ("21.0.13+9", the Temurin
		// tag without the prefix `jdk-`); both spellings read the same
		// checksums files (Docker-Test pin, SPEC-054).
		value = strings.TrimPrefix(value, "jdk-")
	}
	if value == "" {
		return "", fmt.Errorf("devcontainer.features.%s.version: the value is empty after the version prefix is trimmed", key)
	}
	return value, nil
}

// readToolchainChecksums reads the source for the version and fills the pin
// slot; a source that is unreadable (or that does not cover the version) is
// the rejection of SPEC-054 — the same reason, whether no checksum source is
// read at all because no fetcher is wired or a read fails.
func readToolchainChecksums(ctx context.Context, fetcher driven.Fetcher, source string, slot *toolchainPin, version string) error {
	switch source {
	case checksumSourceGo:
		body, err := fetchChecksums(ctx, fetcher, goReleaseListingURL, source)
		if err != nil {
			return err
		}
		pin, err := goChecksumsFor(body, version)
		if err != nil {
			return err
		}
		*slot = pin
	case checksumSourceGolangci:
		url := golangciChecksumsURL(version)
		body, err := fetchChecksums(ctx, fetcher, url, source)
		if err != nil {
			return err
		}
		pin, err := golangciChecksumsFor(body, version)
		if err != nil {
			return err
		}
		*slot = pin
	case checksumSourceJava:
		pin, err := javaChecksumsFor(ctx, fetcher, version)
		if err != nil {
			return err
		}
		*slot = pin
	default:
		return fmt.Errorf("devcontainer catalogue: the checksum source %q of the pinned toolchain takes no URL of an official source", source)
	}
	return nil
}

// fetchChecksums reads one source URL through the driven port; the error
// names the source and the URL so the message of the rejection points at
// what was read.
func fetchChecksums(ctx context.Context, fetcher driven.Fetcher, url, source string) ([]byte, error) {
	if fetcher == nil {
		return nil, fmt.Errorf("the checksum source %s (%s) is not readable: no source reader is wired", source, url)
	}
	body, err := fetcher.Get(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("the checksum source %s (%s) is not readable: %v", source, url, err)
	}
	return body, nil
}

// toolchainWarnings returns the one notice `generate devcontainer` gives for
// each active pinned-toolchain key (SPEC-054, S7/S9): the block, the version
// as the default pin or as an override (with its checksum source) and the
// boundary of the check — the archives of the type are not signed and u-boot
// verifies no signature. Like the notices of the forge-CLIs the entries carry
// neither [driving.WarningEntry.SecurityRelaxation] nor
// [driving.WarningEntry.DegradationFallback], so `--quiet` drops them.
func toolchainWarnings(state toolchainState) []driving.WarningEntry {
	var out []driving.WarningEntry
	if state.goToolchain {
		out = append(out, toolchainNotice(featureGoToolchain, "GO_VERSION", state.goPin, checksumSourceGo, ""))
	}
	if state.golangciLint {
		out = append(out, toolchainNotice(featureGolangciLint, "GOLANGCI_LINT_VERSION", state.lintPin, checksumSourceGolangci, ""))
	}
	if state.javaToolchain {
		out = append(out, toolchainNotice(featureJavaToolchain, "JDK_VERSION", state.javaPin, checksumSourceJava,
			"the JDK of the Eclipse Temurin distribution of the LTS version 21"))
	}
	return out
}

// toolchainNotice is the single notice of one active key; the distribution
// clause (empty for the two Go keys) carries the Vertriebs-Form of the
// java-toolchain key (SPEC-054, Meldungen: sie konkretisiert die Quelle, sie
// ersetzt sie nicht).
func toolchainNotice(key, arg string, pin toolchainPin, source, distribution string) driving.WarningEntry {
	version := "the version pinned as " + arg + " in .devcontainer/Dockerfile is the default pin of the catalogue"
	if pin.version != "" {
		version = "the version " + pin.version + " stands as " + arg + " in .devcontainer/Dockerfile (override of the default pin; " +
			"the checksums were read from " + source + ")"
	}
	what := "devcontainer feature " + key + " installs its release archive from " + source + " when the image is built; "
	if distribution != "" {
		what = "devcontainer feature " + key + " installs " + distribution + " from " + source + " when the image is built; "
	}
	return driving.WarningEntry{
		Code:  "LH-NFA-SEC-004",
		Level: "warn",
		Message: what +
			version + "; it is checked by its SHA-256 only: the archive is not signed and u-boot checks no signature, " +
			"so the check does not protect against a release that was harmful from the start",
		Subject: ".devcontainer/devcontainer.json",
	}
}

// noopFetcher is the nil-tolerant default of the [driven.Fetcher]
// dependency of [GenerateService] — its errors make every `version:`
// override of a pinned-toolchain key the unreadable-source rejection
// (SPEC-054, exit code 10), because a generate run without a wired
// reader cannot verify an override.
type noopFetcher struct{}

func (noopFetcher) Get(_ context.Context, _ string) ([]byte, error) {
	return nil, fmt.Errorf("no source reader is wired")
}
