package main

import (
	"os/exec"
	"strings"
	"testing"
)

// publish.yml package smoke: the asset name is selected from the release's asset
// list, not derived from the tag. A pre-release version (0.10.0-rc.1) is stored by
// GitHub as 0.10.0.rc.1, so a name built from the version 404s.
func runAssetName(t *testing.T, format, assets string) (string, error) {
	t.Helper()
	if _, err := exec.LookPath("bash"); err != nil {
		t.Skip("bash not available")
	}
	cmd := exec.Command("bash", "../../scripts/package-asset-name.sh", format)
	cmd.Stdin = strings.NewReader(assets)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

const assetsPrerelease = `SHA256SUMS
u-boot-0.10.0.rc.1-1.aarch64.rpm
u-boot-0.10.0.rc.1-1.x86_64.rpm
u-boot-darwin-amd64
u-boot-linux-amd64
u-boot_0.10.0.rc.1_amd64.deb
u-boot_0.10.0.rc.1_arm64.deb
`

const assetsStable = `SHA256SUMS
u-boot-0.9.0-1.aarch64.rpm
u-boot-0.9.0-1.x86_64.rpm
u-boot-linux-amd64
u-boot_0.9.0_amd64.deb
u-boot_0.9.0_arm64.deb
u-boot.rb
`

func TestPackageAssetName_Selects(t *testing.T) {
	for _, tc := range []struct{ format, assets, want string }{
		{"deb", assetsPrerelease, "u-boot_0.10.0.rc.1_amd64.deb"},
		{"rpm", assetsPrerelease, "u-boot-0.10.0.rc.1-1.x86_64.rpm"},
		{"deb", assetsStable, "u-boot_0.9.0_amd64.deb"},
		{"rpm", assetsStable, "u-boot-0.9.0-1.x86_64.rpm"},
	} {
		got, err := runAssetName(t, tc.format, tc.assets)
		if err != nil || got != tc.want {
			t.Errorf("%s: got %q (%v), want %q", tc.format, got, err, tc.want)
		}
	}
}

func TestPackageAssetName_FailsClosed(t *testing.T) {
	if _, err := runAssetName(t, "deb", "SHA256SUMS\nu-boot-linux-amd64\n"); err == nil {
		t.Errorf("no deb asset must fail")
	}
	if _, err := runAssetName(t, "deb", "u-boot_1_amd64.deb\nu-boot_2_amd64.deb\n"); err == nil {
		t.Errorf("two deb assets must fail")
	}
	if _, err := runAssetName(t, "zip", assetsStable); err == nil {
		t.Errorf("unknown format must fail")
	}
}
