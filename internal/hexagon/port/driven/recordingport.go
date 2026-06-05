package driven

// RecorderPort is the optional twin of [FileSystem] that surfaces
// the captured mutation log of the recording FS adapter
// (slice-v1-cli-json-dry-run-add T0-(i)).
//
// Production FS adapters (`internal/adapter/driven/fs/`) do NOT
// implement this interface; only `internal/adapter/driven/recordingfs/`'s
// [RecordingFileSystem] does. The Composition-Root in `cmd/uboot/main.go`
// constructs the recorder per [driving.PreviewMode] and hands the
// same instance to the use case both as [FileSystem] and as
// RecorderPort — the use case takes a nil RecorderPort when no
// recording is wanted ([driving.PreviewNone]) and a non-nil one
// otherwise.
//
// LH-FA-ARCH-002/-003 (hexagonal layering): the port lives in
// `port/driven` because it is a capability the application layer
// consumes via interface; only the recording adapter realises it.
// The driving CLI adapter consumes the captured records indirectly
// through [driving.AddServiceResponse.PlannedFiles], never via this
// interface — keeping the CLI's depguard `adapter-driving-no-driven`
// rule clean.
type RecorderPort interface {
	// Captured returns the mutations the wrapped FS has been asked
	// to perform since the recorder was constructed, in call order.
	// The slice is a defensive copy; callers may mutate it without
	// affecting future Captured calls. An empty slice means no
	// mutations have been recorded yet (or were attempted on the
	// underlying FS).
	Captured() []FileMutationRecord
}

// FileMutationRecord carries one FS-mutation event captured by the
// recording adapter. The fields are wire-neutral; the CLI adapter
// maps them into [driving.PlannedFile] for the LH-FA-CLI-007 §326
// JSON envelope.
//
// NewContent is the body argument the use case passed to WriteFile /
// WriteFileExclusive — empty for delete-style actions. OldContent is
// the pre-mutation snapshot the recorder fetched (via the underlying
// FS's ReadFile) before applying the action; nil when the target did
// not exist beforehand (Action then is "create").
//
// Action holds one of the LH-FA-CLI-007 §354 enum values:
// "create", "modify", "delete". The recorder resolves the value
// from the OldContent presence (nil → create) and the called
// method (RemoveAll → delete).
type FileMutationRecord struct {
	Path       string
	Action     string
	NewContent []byte
	OldContent []byte
}
