package domain

// ConfigFileName is the file name `init` creates for the project
// configuration. It is the only place in non-test code that spells
// the name as a path literal; every path to the configuration file
// takes it from here.
const ConfigFileName = "u-boot.yaml"

// ConfigFileCandidates lists the file names the project
// configuration may carry. The list holds
// [ConfigFileName] only. Each call returns a fresh slice, so callers
// may modify the result.
func ConfigFileCandidates() []string {
	return []string{ConfigFileName}
}
