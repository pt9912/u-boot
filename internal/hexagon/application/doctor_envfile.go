package application

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// checkDevcontainerEnvFile is the LH-FA-DEV-012 check: when
// `.devcontainer/devcontainer.json` carries the `--env-file` line,
// the file it names must exist, because `docker run` aborts with exit
// code 125 otherwise. A missing, unreadable or line-less
// devcontainer.json skips the check (ok with a message). The check asks
// only [driven.FileSystem.Exists] of the environment file: it neither
// reads the content nor asks Git. A missing file is a plain warn item,
// so `--quiet` drops it.
func (s *DoctorService) checkDevcontainerEnvFile(_ context.Context, baseDir string) domain.Diagnostic {
	body, err := s.fs.ReadFile(filepath.Join(baseDir, ".devcontainer", "devcontainer.json"))
	if err != nil || !strings.Contains(string(body), envFileRunArg) {
		return domain.Diagnostic{ID: checkIDDevcontainerEnvFile, Severity: domain.SeverityOK,
			Message: "devcontainer.json does not reference " + devcontainerEnvFileRel + "; environment file check skipped."}
	}
	exists, err := s.fs.Exists(filepath.Join(baseDir, devcontainerEnvFileRel))
	if err != nil {
		return domain.Diagnostic{ID: checkIDDevcontainerEnvFile, Severity: domain.SeverityOK,
			Message: "cannot check " + devcontainerEnvFileRel + ": " + err.Error() + "; environment file check skipped."}
	}
	if !exists {
		return domain.Diagnostic{ID: checkIDDevcontainerEnvFile, Severity: domain.SeverityWarn,
			Message: "devcontainer.json passes --env-file " + devcontainerEnvFileRel + ", but the file is missing; starting the container fails.",
			Hint:    "Run `u-boot generate devcontainer` to create " + devcontainerEnvFileRel + "."}
	}
	return domain.Diagnostic{ID: checkIDDevcontainerEnvFile, Severity: domain.SeverityOK,
		Message: devcontainerEnvFileRel + " exists."}
}
