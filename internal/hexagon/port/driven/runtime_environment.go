package driven

// RuntimeEnvironment exposes coarse facts about the runtime u-boot
// itself is executing in. Today the only fact is "are we inside a
// container?" — the answer drives `doctor`-skip semantics for the
// host-prerequisite checks (`docker.*`, `git.installed`) per
// `slice-v0.1.1-doctor-container-awareness` (v0.1.0 GHCR-distroless
// container has no docker / git binary on $PATH; without the skip
// the user gets 4 false-positive errors on a healthy host).
//
// The interface is deliberately I/O-free in its declared signature:
// detection is best-effort, and a missing detection file means
// "not in container", not "I/O error". Adapters that look at
// `/.dockerenv` / `/run/.containerenv` / cgroup-v1 markers absorb
// the OS-level errors and just return the boolean.
type RuntimeEnvironment interface {
	// InContainer reports whether u-boot detects itself running in
	// a container. Implementations must default to `false` when
	// detection is ambiguous — false positives would skip checks
	// the user actually wants.
	InContainer() bool
}
