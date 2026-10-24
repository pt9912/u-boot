# Changelog

All notable changes to **u-boot** are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Internal `u-boot generate changelog` ([`LH-FA-GEN-001`](spec/lastenheft.md#lh-fa-gen-001--generate-befehl)..[`LH-FA-GEN-005`](spec/lastenheft.md#lh-fa-gen-005--idempotenz), [`LH-AK-007`](spec/lastenheft.md#lh-ak-007--changelog-generator))
maintains a Keep-a-Changelog-formatted changelog for user projects;
this file is the same format applied to u-boot itself.

## [Unreleased]

### Added

- Catalogue feature `ollama-host` ([`LH-FA-DEV-003`](spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-006`](spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-CLI-005`](spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)): `u-boot config set devcontainer.features.ollama-host.enabled true`, then `u-boot generate devcontainer` writes one line into `runArgs` of `devcontainer.json`, `--add-host=host.docker.internal:host-gateway`, in the default and in the sandbox profile, and nothing else for the connection: no relay, no environment variable, no host in the egress list. It is no OCI feature and renders no key under `features`; `version` without `source` ends with exit code `10` like for `ollama`, `gh` and `glab`. `generate devcontainer` discloses the connection with one notice per run; in the sandbox profile it is a security relaxation and stays visible under `--quiet`, and while the egress restriction is active it adds that the connection does not work behind it (a warning, no error). Pinned by tests, including one run of the container started by `u-boot devcontainer up` that resolves `host.docker.internal`; that a service of the host is reachable is not tested.

### Changed

- In the default profile, `devcontainer.json` carries a `runArgs` key when the feature `ollama-host` is active; without it the file is byte-identical to before.
- The user manual and both READMEs no longer say that cloud models of Ollama are not included ([`LH-FA-DEV-003`](spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [ADR-0036](docs/plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md)). The manual describes the flow inside the container of the sandbox profile: start the server, sign in with `ollama signin`, then `ollama run`, `ollama launch claude` or `claude -p`; you sign in yourself; `u-boot` sets no variable for it and adds no `ollama.com` to the egress list. It names the limits of what is promised: only the measured setup (a container you start yourself, Linux `amd64`, local Docker engine, no terminal, two models, one to three runs, hand variants), `ANTHROPIC_API_KEY` must not be set for `claude -p`, and the flow in the container that `u-boot devcontainer up` starts and behind the egress restriction with a sign-in are not promised. It also states the measured effect of `OLLAMA_NO_CLOUD` and of model names with `[1m]`, and where the device key lives (the volume of the models). New manual section for `ollama-host`: the relay and the binding of the service stay with you, what a relay opens (every container in the network `bridge`, not checked with a foreign container), `claude -p` through the service of the host with a relay (the line set by hand, not the one `u-boot` generates), and that the connection does not work behind the egress restriction. The specification covers `ollama-host` and the cloud paragraph of `ollama`; its sandbox rule no longer says that `runArgs` comes only from `nestedRuntime` and the egress restriction.

## [0.15.1] - 2026-10-24

Patch release. **Two checks of the feature entries in `u-boot.yaml` now run where they were promised.** `u-boot generate devcontainer` rejects the source of an enabled feature that is not on the allowlist, and `u-boot config set devcontainer.features.<feature>.version <value>` rejects a `version` without `source` for `ollama`, `gh` and `glab`.
**Behaviour change for scripts:** `u-boot generate devcontainer` ends with exit code `10` instead of `0` and writes nothing when the source of an enabled feature is not allowed; the entry below names the two ways out (allow the source, or switch the feature off). `u-boot config set` ends with exit code `10` instead of `0` for a `version` without `source` at `ollama`, `gh` and `glab`.
The lead of `[0.15.0]` named the check in `config set` as a follow-up; this release delivers it. A file that already carries such a `version` entry stays rejected; the way out is to remove the entry by hand.
**Measured** at the built binary on one Linux host (`docker run` with a project mount) and by mutation tests of the verifiers. `config set`: probed at the binary for `gh`, `glab` and `java`; `ollama` is pinned by tests at the service level only, not probed at the binary. `generate`: probed at the binary with one entry from a source that is not allowed, in the default and in the sandbox profile, with the ways out (allow the source with `config set`, allow it with `--allow-external-feature-sources`, switch the feature off) and with a catalogue feature without `source` (`node`), which is accepted. **Not measured:** whitespace at the edge of a source in `u-boot.yaml`, `u-boot doctor` with a missing `enabled`, the `--json` output of the rejection in `config set`, macOS, Podman and other hosts.
No new command and no new option. Lastenheft 0.22.0.

### Fixed

- `u-boot generate devcontainer` now checks the allowlist of external feature sources ([`LH-FA-DEV-003`](spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-NFA-SEC-004`](spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)): an enabled feature whose `source` is not listed byte-for-byte in `devcontainer.featureSources.allow`, or given with `--allow-external-feature-sources` in the same call, ends with exit code `10` and code [`LH-FA-DEV-003`](spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), and writes nothing, also with `--dry-run`, in the default and in the sandbox profile. Before, the call exited `0` and wrote the unlisted source into `devcontainer.json`, while `u-boot doctor` reported the entry as an error and `u-boot config set …source` refused it. Projects that carry such an entry now need one of the two ways out the message names: allow the source (`u-boot config set devcontainer.featureSources.allow <source>` or the flag), or switch the feature off (`u-boot config set devcontainer.features.<name>.enabled false`).
- `u-boot config set devcontainer.features.<feature>.version <value>` (also with several path/value pairs in one call) now checks the value like `u-boot generate devcontainer` ([`LH-FA-CONF-005`](spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern)): for `ollama`, `gh` and `glab`, a `version` without `source` ends with exit code `10` and leaves `u-boot.yaml` unchanged, also with `--dry-run` and `--diff`. Before, the call reported success and wrote the entry, after which `generate`, `doctor` and every further `config set` rejected the project. A file that already carries such an entry is still rejected; remove the `version` entry by hand.

## [0.15.0] - 2026-10-23

Fifteenth release. **Forge CLIs in the devcontainer: the catalogue features `gh` (GitHub CLI) and `glab` (GitLab CLI).**
`u-boot config set devcontainer.features.gh.enabled true` (and `glab`), then `u-boot generate devcontainer`: the tool is in the container after the build, in the default and in the sandbox profile (entry below). A token reaches the container only as an environment variable of the host (`GH_TOKEN`, `GITLAB_TOKEN`); in the sandbox profile `devcontainer.json` carries a `${localEnv:…}` reference, never a value. The default profile carries no reference; the way there, `docker exec -e GH_TOKEN`, is in the user manual (measured once with a placeholder).
**Promised only on a Linux host (`amd64`) with a local Docker engine** (measured with Docker 29.9.2); **not promised:** macOS with Colima (whoever installs `u-boot` with Homebrew finds the features, but they are not promised there), Docker Desktop, rootless Docker, Podman as host engine, a remote engine, Linux `arm64`, Windows.
**No authenticated call of `gh` or `glab` is measured:** every call with a token variable ran with a placeholder and ended in 401 or `Bad credentials`. Measured is what works without a valid token: the build of the generated Dockerfile in both profiles (also once per profile without the layer cache, downloads from GitHub and GitLab, base image already local), the versions, the notices, the egress hosts and the behaviour with and without the variable. Not promised: `glab duo` and `glab orbit`, self-hosted instances, hosts after a sign-in, telemetry and update hosts (the user manual names them).
**Limit of `version:`:** for the features `gh`, `glab` and `ollama`, `version:` without `source:` ends `u-boot generate devcontainer` with exit code `10`; other catalogue features such as `java` take `version:` without `source:`. `u-boot config set devcontainer.features.gh.version <value>` (and `glab`) does not check this yet: it reports success and writes the entry, and then `generate`, `doctor` and every further `config set` reject the project (one run, Linux `amd64`); the way out is to remove the `version` entry from `u-boot.yaml` by hand. The check in `config set` is a follow-up, not part of this release.
Configurations without `gh` and `glab` render byte-identical. Lastenheft 0.22.0.

### Added

- Catalogue features `gh` (GitHub CLI) and `glab` (GitLab CLI) ([`LH-FA-DEV-003`](spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-009`](spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)). `u-boot config set devcontainer.features.gh.enabled true` (and `glab`) makes `u-boot generate devcontainer` write a Dockerfile block per tool in the default and in the sandbox profile: the pinned release archive (`gh` 2.103.0 of `cli/cli`, `glab` 1.124.0 of `gitlab-org/cli`), checked by its SHA-256 before the extraction, only the program extracted. No signature was found for either archive (by search pattern, not exhaustive); the GitHub attestation of `gh` is not verified by `u-boot`; `u-boot generate devcontainer` prints one notice per active feature, and `--quiet` drops it. For `gh` and `glab`, `version:` without `source:` ends `u-boot generate devcontainer` with exit code `10` (`u-boot config set` does not check it yet, see the lead above). With the egress restriction the hosts `api.github.com` and `github.com` (`gh`; both are in the base list already) and `gitlab.com` (`glab`) are allowed. In the sandbox profile `remoteEnv` passes `GH_TOKEN` and `GITLAB_TOKEN` as `${localEnv:…}` references next to `GIT_TOKEN` (never a value), and `u-boot doctor` checks the reference and plain-text tokens in project files. The default profile gets no token reference. Calls with a valid token, hosts after a sign-in, telemetry and update hosts, self-hosted instances, `glab duo` and `glab orbit` are not promised.

### Changed

- With `gh` or `glab` active, the generated Dockerfile of the default profile carries `USER root` ahead of the block (it stood there before only with the sandbox profile or a UID override). Configurations without the two features render byte-identical.

## [0.14.0] - 2026-10-23

Fourteenth release. **Open a shell in the devcontainer without an editor (`u-boot devcontainer shell`); `u-boot devcontainer up` judges by the state of the container.**
`u-boot devcontainer shell` starts the devcontainer like `u-boot devcontainer up` (without `--recreate`) and opens an interactive shell in it, as the user and in the workspace the start reports, with the credentials that the generated `devcontainer.json` passes on as `${localEnv:NAME}`; Docker (client and local engine) is enough. The helper container holds the engine socket of the host only for the duration of the start (the command says so before it starts the helper, also under `--quiet`); the shell runs without a helper container and without a socket that the command mounts (measured in one session: no container `u-boot-devcontainer-up-…` ran, and the value of `GIT_TOKEN` was not in `docker inspect` of the container; Linux `amd64`, Docker 29.9.2).
**Promised only on a Linux host (`amd64`) with a local Docker engine** (Docker tests, measured with Docker 29.9.2); **not promised:** macOS with Colima (whoever installs `u-boot` with Homebrew finds the command, but it is not promised there), Docker Desktop, rootless Docker, Podman as host engine, a remote engine, Linux `arm64`, Windows. The interactive session is measured under `script` (a pseudo-terminal); a terminal operated by a person is not measured.
**Behaviour change for scripts:** a repeated `u-boot devcontainer up` without `--recreate` against a container whose start script reported `failed` ends with exit code `11` instead of `0`, also when the Dev Container CLI itself succeeds; the message names `--recreate`, which helps once the cause is fixed (with the cause still present, the recreated `up` ended with `11` again in one run). This needs the `u-boot.yaml` of the project (entry below).
**Limit of a repeated `up`:** against a running container, a repeated `up` without `--recreate` does not run the `postCreateCommand` and the `postStartCommand` again (the `postAttachCommand` runs on every `up`; entry below). After exit code `12` (a `postStartCommand` ended with 3) and after Ctrl-C during a `postCreateCommand`, the devcontainer kept running, and the following `up` ended with exit code `0` and `Devcontainer already running` without running the command again (one run each, Linux `amd64`, Docker 29.9.2; after Ctrl-C the command returned after 10.5 seconds). Exit code `0` therefore does not mean that the lifecycle commands ran through; `shell` then opens the session. Other failures (for example a failed build) are not measured. `u-boot devcontainer up --recreate` creates the container anew.
New in the `--json` output of `u-boot devcontainer up`: `data.alreadyRunning` (always present, additive). Lastenheft 0.21.0.

### Added

- `u-boot devcontainer shell` starts the devcontainer like `u-boot devcontainer up` (without `--recreate`) and opens an interactive shell in it, as the user and in the workspace the start reports. The variables that the generated `devcontainer.json` passes on as `"NAME": "${localEnv:NAME}"` entries (in the sandbox profile `GIT_TOKEN` and the key of an agent) are set in the shell when they are set on the host. The helper container holds the engine socket only for the duration of the start, and the command says so before it starts the helper, also under `--quiet`; the shell runs without a helper container and without a socket that the command mounts. Everything `u-boot` prints goes to standard error, standard output belongs to the shell. Once the shell has begun the command ends with exit code `0` whatever the exit status of the shell was; a container whose start script reported `failed` opens no shell (exit code `11`); `--json`, `--dry-run`, `--diff` and arguments end with exit code `2`. End the session with `exit`; if `u-boot` ends in another way (for example `kill -9`), the shell in the container can keep running with its commands. Promised on a Linux host (`amd64`) with a local Docker engine, like `up`.

### Changed

- **Behaviour change for scripts:** `u-boot devcontainer up` judges by the state of the container. A repeated `up` without `--recreate` against a container whose start script reported `result=failed` in its status file now ends with exit code `11` instead of `0`, also when the Dev Container CLI itself succeeds. The message names the cause of the capability (the sentence `u-boot doctor` shows), the time of the report, that the container is running and that the call did not run the start script, and `u-boot devcontainer up --recreate`. This needs a `u-boot.yaml` in the project folder, as before: without it the `up` that aborts ends with `12` and a repeated `up` does not report the earlier abort and ends with `0`. A script that repeats `up` after a start script reported `failed` in its status file gets `11` (with a `u-boot.yaml`, as above); the way out is `u-boot devcontainer up --recreate`. Measured with the start script of the egress restriction (`onUnavailable: fail`, capability `NET_ADMIN` withdrawn): the first, the second and a third `up` end with `11`, and the same project without `u-boot.yaml` ends with `12` and then `0` (Docker tests, Linux `amd64`, Docker 29.9.2, one run); the start scripts for `nestedRuntime: podman` and `docker` are pinned only by unit tests with a simulated engine ([`LH-FA-DEV-010`](spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006--exit-codes)).
- `u-boot devcontainer up` says when the container was already running: if it ran before the call, the command ran without `--recreate` and the Dev Container CLI reports the same ID, the result line reads `Devcontainer already running: container …` instead of `Devcontainer started: …`, and the `--json` envelope carries `data.alreadyRunning: true` (the field is always present, `false` otherwise; `status` and `exitCode` do not change). The command asks the engine once for the container before it starts the helper container; a failed query is no failure of the command. Measured with a project whose `postCreateCommand`, `postStartCommand` and `postAttachCommand` each append a line to a file (Docker test, Linux `amd64`, Docker 29.9.2, Dev Container CLI 0.89.0, one run): against a running container the `postCreateCommand` and the `postStartCommand` do not run again; after `u-boot devcontainer down` and `up` the `postStartCommand` runs again in the same container and the `postCreateCommand` does not; the `postAttachCommand` runs on every `up` ([`LH-FA-DEV-010`](spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-CLI-007`](spec/lastenheft.md#lh-fa-cli-007--dry-run)).

## [0.13.1] - 2026-10-22

Patch release. **Security fix for the sandbox devcontainer profile.** The sandbox profile promises not to mount the container-engine socket of the host ([`LH-FA-DEV-006`](spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)). The built-in feature `docker-cli` (`docker-outside-of-docker`) mounts it according to its metadata, and from `v0.6.0` to `v0.13.0` both profiles rendered the feature without a warning.
`u-boot generate devcontainer` now ends with exit code `10` when the sandbox profile and an active `docker-cli` meet, and `u-boot doctor` reports that state as an error (exit code `11`); the standard profile keeps `docker-cli`. The sandbox profile also discloses the container options that the built-in feature `go` brings and warns about each feature from an external source. The statements about the socket in the disclosure of `u-boot devcontainer up`, in its `--help` and in the manual were wrong for both profiles and are corrected.
**An existing container keeps the socket until it is recreated; the way out is in the entry *Security*.** Whether a project is affected depends on whether its container was built with the feature; on the current default base image the build of `docker-cli` stopped in two runs (entry *Security*).
The Go toolchain moves from 1.27.2 to 1.27.3 because a HIGH stdlib advisory failed the image scan (entry *Security*).
No new command. **Behaviour change for scripts:** in the sandbox profile with an active `docker-cli`, `u-boot generate devcontainer` ends with exit code `10` instead of `0`, and `u-boot doctor` with `11`. Lastenheft 0.21.0; on `main` the Lastenheft also describes `u-boot devcontainer shell` and a further case of exit code `11` for a repeated `u-boot devcontainer up` (0.20.0) — **neither is part of this release.**

### Security

- **The sandbox profile promises not to mount the container-engine socket of the host, and the built-in feature `docker-cli` declares such a mount.** `docker-cli` installs the feature `docker-outside-of-docker`; its metadata declare a bind mount of `/var/run/docker.sock`, and the Dev Container CLI merges the mounts of a feature into the start of the container (read on 2026-10-22 in the metadata of `docker-outside-of-docker:1`, which was version 1.10.1 that day, and in the Dev Container CLI `0.89.0`). A process in the container that can use the socket has root-equivalent access to the engine of the host. From `v0.6.0` to `v0.13.0` the generator rendered the feature in both profiles without a warning (measured with `v0.13.0`, read in the code of the earlier versions). Rated High: CVSS 3.1, `CVSS:3.1/AV:L/AC:H/PR:L/UI:N/S:C/C:H/I:H/A:H`, base score 7.8; there is no CVE ([`LH-FA-DEV-006`](spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`ADR-0032`](docs/plan/adr/0032-katalog-features-socket-und-lockerungen.md)).
  **Measured:** in the standard profile one Docker test starts a project with `docker-cli` and finds the bind mount of `/var/run/docker.sock` (as `/var/run/docker-host.sock`) in the mounts of the started container, read with `docker inspect` (Linux `amd64`, Docker 29.9.2, two runs on 2026-10-22), with the option `moby: false` set by hand in the rendered key. **Not measured:** the sandbox profile before this release, and the mount with the default option `moby: true`.
  **Whether a project is affected depends on the base image.** On the current default base image of the generated Dockerfile (`mcr.microsoft.com/devcontainers/base:debian`, Debian 13 `trixie` on the day of the test) the build with the default option stopped in two runs on 2026-10-22, both in the standard profile (the Docker test and `u-boot devcontainer up` on a throw-away project; Linux `amd64`, Docker 29.9.2), with "The 'moby' option is not supported on debian 'trixie'", so no container started. A container that was built with the feature earlier, or on another base image, carries the mount according to the metadata (read, not measured there).
  **What to do** if a project with the sandbox profile used `docker-cli`: disable the feature with `u-boot config set devcontainer.features.docker-cli.enabled false`; whoever needs Docker inside the container also sets `u-boot config set devcontainer.sandbox.nestedRuntime docker` (the profile then gets `--privileged`); run `u-boot generate devcontainer` again; recreate the container with `u-boot devcontainer up --recreate`, or rebuild it in the editor. **An existing container keeps the socket until it is recreated** (read, not measured: the mounts of a container are fixed when it is created). Disabling the feature and generating again need no update: with `v0.13.0` the same `config set` and `generate devcontainer` end with exit code `0`, and the key of the feature is gone from `devcontainer.json` (run on 2026-10-22; the start of the container was not part of that run). **Not measured:** that a container recreated after these steps starts without the socket, and the step with `nestedRuntime docker` (promised on Docker under Linux and on Colima, not on Docker Desktop). From `v0.13.1` on, `u-boot doctor` reports the state of the project configuration (exit code `11`).
  In the standard profile `docker-cli` stays selectable; the disclosure of `u-boot devcontainer up` and the manual now say that the generated configuration mounts the socket when the feature is enabled.
- `build(deps): Go toolchain 1.27.2 → 1.27.3` — fixes one HIGH stdlib advisory (CVE-2026-111943, `crypto/x509`, untrusted intermediate certificates) that failed the Trivy image scan on 2026-10-22 (Trivy 0.70.0, runtime image built with Go 1.27.2; fixed in Go 1.27.3).

### Fixed

- The sandbox profile promises not to mount the container-engine socket of the host. The built-in feature `docker-cli` (`docker-outside-of-docker`) declares such a mount in its metadata, and both profiles rendered the feature. `u-boot generate devcontainer` now ends with exit code `10` when the sandbox profile (from `devcontainer.profile: sandbox` or `--sandbox`) and an active `docker-cli` meet: nothing is written, also with `--dry-run`, and `--sandbox` is not written to `u-boot.yaml`. `u-boot doctor` reports the same state as an error of `devcontainer.features.allowlist` (exit code `11`). The standard profile keeps `docker-cli`.
  The error message names the way out: disable the feature with `u-boot config set devcontainer.features.docker-cli.enabled false`; whoever needs Docker inside the container also sets `u-boot config set devcontainer.sandbox.nestedRuntime docker` (`nestedRuntime: docker` alone does not lift the error); then run `u-boot generate devcontainer` again and recreate an existing container with `u-boot devcontainer up --recreate`. `config set` stays possible in such a project.
  Measured with a started container in the standard profile (one Docker test, Linux `amd64`, Docker 29.9.2, two runs on 2026-10-22, with the option `moby: false` set by hand in the rendered key): the container carries a bind mount of `/var/run/docker.sock`.
  Affected: `v0.6.0` to `v0.13.0` ([`LH-FA-DEV-006`](spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006--exit-codes)).
- The disclosure of `u-boot devcontainer up`, its `--help` and the manual said that the devcontainer configuration that `u-boot` generates does not mount the socket, for both profiles. They now say that the generated configuration mounts it if a feature that does so is enabled, such as the built-in feature `docker-cli`, and that the sandbox profile rejects that feature. The text does not say "only": `u-boot` does not read the metadata of a feature from an external source ([`LH-FA-DEV-010`](spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)).

### Changed

- **Behaviour change for scripts:** a project with the sandbox profile and an active `docker-cli` ended `u-boot generate devcontainer` with exit code `0`; it now ends with `10`, and `u-boot doctor` ends with `11` (see above).
- `u-boot generate devcontainer` in the sandbox profile warns about each active built-in feature that adds container options (`go`: `--cap-add=SYS_PTRACE` and `--security-opt=seccomp=unconfined`, one warning each, read from the metadata of the feature) and about each enabled feature with a `source:` override (`u-boot` does not read its metadata, and the profile promises none of its properties for it). Both kinds of warning stay visible under `--quiet`; the exit code and the written files do not change. The Lastenheft allows the options that a built-in feature brings when the command output discloses them ([`LH-FA-DEV-006`](spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-CLI-005`](spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)). Lastenheft 0.21.0.

## [0.13.0] - 2026-10-21

Thirteenth release. **Start and stop the devcontainer without an editor (`u-boot devcontainer up` and `down`), DNS to the upstream resolvers only for the local resolver of the egress restriction, a repeated run of the egress start script in one transaction.**
`u-boot devcontainer up [--recreate]` and `u-boot devcontainer down` start and stop the devcontainer of a project without an editor, Node.js or the Dev Container CLI on the host; Docker (client and local engine) is enough.
A short-lived helper container holds the engine socket of the host for the duration of the call (root-equivalent access to the engine); the command says so before it starts the helper, also under `--quiet`.
**Promised only on a Linux host (`amd64`) with a local Docker engine** (Docker tests); **not promised:** macOS with Colima (whoever installs `u-boot` with Homebrew finds the command, but it is not promised there), Docker Desktop, rootless Docker, Podman as host engine, a remote engine, Linux `arm64`, Windows.
**The first call needs network access** and builds the helper image locally (single measurement: 1716 MiB, about 1.7 GiB); Trivy reported 597 findings (HIGH and CRITICAL) in the operating-system part of the Node base image of the helper and none in the Dev Container CLI (single measurement); the image is neither shipped nor scanned continuously.
**Further limits:** after Ctrl-C the command returns only after about ten seconds; the images `vsc-…` stay after `down` and `--recreate`; the masking of credentials in the output replaces only unchanged values of at least four bytes; the user ID in the container is not adapted to the owner of the project folder; the exit codes are `10`, `11`, `12`, `14` and `2` (entry below);
after `up` ended with exit code `12` because a `postStartCommand` failed, and after Ctrl-C during a `postCreateCommand`, the container kept running and the next `up` reported success without running the failed `postStartCommand` or the aborted `postCreateCommand` again (observed in one run each; other failures, such as a failed build, and exit code `11` from a start script not measured), and `up --recreate` creates the container anew.
**Behaviour changes:** the egress start script ends with exit code `11` and `failed`, also with `onUnavailable: warn`, when the user `nobody` is missing, the resolver does not start or does not run under that user, or when its `nft -f` transaction does not go through (a host on which `nft` cannot create the `ip6` table ends without any rule; read in the script, not run on such a host);
in the sandbox container `podman run --dns <upstream resolver>` resolves nothing any more, and an unlisted name fails in nested Podman after about 36 seconds instead of about 20 milliseconds; both take effect for an existing project after `u-boot generate devcontainer` and a rebuild of the container.
`u-boot generate devcontainer` warns when it replaces an existing `egress-init.sh` (with `--json` the envelope status is then `warn`, the exit code `0`).
The DNS change is measured on one native Linux x86_64 host with Docker 29.9.1 in two runs (**not measured:** Colima, Docker Desktop, `arm64`, CI); the transaction is pinned by a test, its timing is not measured (entries below).
**In CI** the Docker integration tests passed, and the four tests of nested Podman were skipped (the runner restricts unprivileged user namespaces), so CI gives no evidence for nested Podman. Lastenheft 0.18.0.

### Added

- `u-boot devcontainer up [--recreate]` and `u-boot devcontainer down` start and stop the devcontainer of a project without a devcontainer-capable editor, Node.js or the Dev Container CLI on the host; Docker (client and local engine) is enough.
  `up` runs the Dev Container CLI in a short-lived helper container and builds the helper image locally from pinned sources on the first call (network access needed). The helper container holds the engine socket of the host for the duration of the call, which is root-equivalent access to the engine;
  the command says so before it starts the helper, also under `--quiet`, and as a `diagnostics[]` entry with `--json` (so the envelope status of a successful `up` is `warn`, the exit code `0`). The devcontainer configuration that `u-boot` generates does not mount the socket, a hand-written `.devcontainer/devcontainer.json` can, and its `initializeCommand` runs with the access of the helper.
  The helper container receives only the variables that `devcontainer.json` references as `${localEnv:<NAME>}` and that are set on the host, by name. `remoteEnv` is the form meant for credentials: with `containerEnv` the resolved value stays in the metadata (`Config.Env`) of the started container as long as it exists, and with both forms the value is on a command line inside the helper container for the duration of a lifecycle command (`docker exec -e` for `remoteEnv`, `docker run -e` for `containerEnv`), visible to processes with access to the daemon.
  `--recreate` removes an existing container and creates it anew, named volumes stay. `down` stops the container and removes neither it nor a volume; without a container, or with a stopped one, it ends with exit code `0`.
  After the start `up` discloses the degradation fallbacks the container reports, if the project has a `u-boot.yaml`. Exit codes: `10` (no `.devcontainer/devcontainer.json`, a path with a comma, double quote or line break, a `${localEnv:…}` name that is not a plain variable name), `11` (no `docker` binary, the daemon does not answer, a remote engine, or a container start script that aborts and reports `failed`; the last needs a `u-boot.yaml`, otherwise the exit code is `12`), `12` (any other failure of the build or the start), `14` (a read error in the file system), `2` (bare `u-boot devcontainer`, `--dry-run`, `--diff`).
  The shipped container image carries no `docker` binary, so the commands end there with exit code `11`.
  The manual describes the way into the container without an editor (`docker exec -it -u <user> -w <workspace> <container> bash` with the values of the result line; `remoteEnv` does not apply in that shell, so `GIT_TOKEN` of the sandbox profile was not set there, and `docker exec -e <NAME>` passes a variable explicitly; one run each) and a repeated `up` after a failure or an abort (observed in one run each on Linux `amd64` with Docker 29.9.1: after exit code `12` because a `postStartCommand` failed, and after Ctrl-C during a `postCreateCommand`, the container kept running and the next `up` ended with exit code `0` without running the failed `postStartCommand` or the aborted `postCreateCommand` again, so this exit code `0` does not mean that the lifecycle commands completed; other failures, such as a failed build, and exit code `11` from a start script not measured; `up --recreate` creates the container anew, named volumes stay).
  Promised: a Linux host with a local Docker engine on `amd64` (Docker tests in `make test-docker`); not promised: Colima, Docker Desktop, rootless Docker, Podman as host engine, a remote engine, Linux `arm64`, Windows (the command itself is measured on none of them; on Colima only the measuring route, the Dev Container CLI in a throw-away container, has run). A catalogue feature (`node`) is installed in the test; other features are not tested
  ([`LH-FA-DEV-010`](spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-CLI-005`](spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-NFA-USE-004`](spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe)).
- `u-boot generate devcontainer` warns when it replaces an existing `.devcontainer/egress-init.sh`: an existing container keeps the egress list of the image it was built from and has to be rebuilt
  (`devcontainer up --remove-existing-container`). The warning does not appear when the file is written for the first time, when nothing changes or under `--dry-run`; `--quiet` drops it, and with
  `--json` it is one entry of `diagnostics[]` (envelope status `warn`, exit code `0`, as for every warning). The warning appears only when the managed block of an existing `egress-init.sh` is replaced: activating the restriction for the first time (the script is created) and switching it off (the script stays in place) give no warning
  ([`LH-FA-DEV-008`](spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-CLI-005`](spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)).

### Changed

- Running the egress start script (`egress-init.sh`) again in a running container replaces both nftables tables in a single `nft -f` transaction instead of deleting them and building them anew, stops the old local resolver with a bounded wait before the new one starts, and writes `/etc/resolv.conf` last. A transaction that does not go through ends the script with exit code `11` and status `failed`, also with `onUnavailable: warn`, and leaves the previous rules and resolver in place. A host on which `nft` cannot create the `ip6` table therefore ends without any rule and with exit code `11` (read in the script, not run on such a host); before, the IPv4 table stayed in place. The exit codes, the status file and the user of the resolver are otherwise unchanged. The structure is pinned by a test; a window of zero and an uninterrupted name resolution are **not** promised. Not measured for this version: `amd64`, Docker Desktop, CI, IPv6, load, `nestedRuntime`. The counts of the manual (attempts that passed the restriction during a repeated run) describe the earlier version
  ([`LH-FA-DEV-008`](spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)).
- The egress restriction of the sandbox devcontainer opens DNS (port 53, UDP and TCP) to the upstream resolvers only for the user of the local resolver (`nobody`, set with `user=` in the generated resolver configuration; the firewall rules take the uid from the same value). A process without `sudo` no longer gets an answer for a name outside the allowlist from the upstream resolvers.
  In a nested Podman container an unlisted name takes about 36 seconds to fail instead of about 20 milliseconds, so a `podman build` step that resolves one runs longer (43 instead of 7 seconds in the replayed build), and `podman run --dns <upstream resolver>` resolves nothing any more; the nested `podman build` behind the restriction, which failed on `amd64`, passes.
  The start script ends with exit code `11` and writes `failed` to the status file, also with `onUnavailable: warn`, when the user `nobody` does not exist in the image, the resolver does not start, or the resolver does not run under that user. Still open: everything with `sudo`, traffic to the addresses of allowed names, and queries for subdomains of allowed names, which the resolver forwards.
  Measured on a native Linux x86_64 host with Docker 29.9.1 in two runs; not measured: Colima, Docker Desktop, `arm64`, CI, `nestedRuntime: docker`. An existing container keeps its rules until `u-boot generate devcontainer` and a rebuild
  ([`LH-FA-DEV-008`](spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)).
- The `Notes:` of `u-boot egress list` and its `--help` say that a changed list takes effect after a rebuild of the container, and state the subdomain behaviour only as far as a test pins it
  (an entry covers its subdomains one and two levels below it; an entry for a subdomain does not open its parent name for a connection by name). The error for a wildcard entry reads
  "an allowed name already covers its subdomains, one and two levels below it" ([`LH-FA-DEV-008`](spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)).

## [0.12.0] - 2026-10-19

Twelfth release. **Egress list display (`u-boot egress list`), Ollama as a sandbox devcontainer feature, subdomain clarification for the egress allowlist.**
`u-boot egress list` shows the effective egress allowlist with the origin of each host, offline and from the project files; its `--json` output reports `command: "egress"`, which joins the
values of the `command` field. A wildcard entry (`*.example.org`) in `devcontainer.sandbox.egress.allow` is still rejected with exit code `10`, and the error now names the host to enter;
the manual states the measured behaviour for subdomains. The built-in catalogue gains `ollama` (local language models, **CPU only**, effective **only with the sandbox profile**): the generated Dockerfile
installs Ollama 0.40.6 from the pinned, checksum-verified release archive, and with the egress restriction the default allowlist gains the two hosts of the model pull.
**Not included: cloud models of Ollama** (names ending in `-cloud`). Such a model can be pulled, but it does not run: the host `ollama.com` is not on the allowlist, and `ollama signin` without an account reported `account unavailable` (signing in with an account was not examined).
**A changed egress allowlist takes effect only after a rebuild of the container**; there is no reload of a running container, and a restart keeps the old list.
`u-boot --help` links the user manual, and the `--dry-run` log lines of `config set` and `generate` read in the conditional. Lastenheft 0.17.0.
**Measured** with Docker 29.5.2 on Colima 0.10.3 (arm64) and, for the `ollama` feature, in two runs on **one** Linux x86_64 host with a native Docker Engine 29.9.1 (2026-10-19; no CI run);
**not measured:** Docker Desktop, Linux behind a VM, larger models, whether the second model-pull host stays the same.
**Known on `amd64`:** a nested `podman build` behind the egress restriction fails at a `RUN` step that expects an unlisted name not to resolve: `example.com` resolves (IPv6 addresses only), because the nested container asks the upstream resolver of the Docker host directly, and the restriction allows that; the connection to the address stays blocked. Measured in two runs on one Linux x86_64 host (native Docker Engine 29.9.1); with another upstream resolver (`--dns 9.9.9.9`) the build passed (one run); on Colima, arm64, and in CI not measured.

### Added

- Devcontainer features: the built-in catalogue gains `ollama` (Ollama, local language models, CPU only). `u-boot config set devcontainer.features.ollama.enabled true` makes the
  sandbox Dockerfile install Ollama 0.40.6 from the pinned release archive of `ollama/ollama` (SHA-256 checked before the extraction; the GPU libraries are left out; `zstd` is removed again)
  and mounts the named volume `<project>-ollama-${devcontainerId}` on `/home/vscode/.ollama`. It is no OCI feature and renders no key under `features`; it takes effect only with the sandbox
  profile (without it `generate devcontainer` says so and generates nothing). The upstream publishes no signature for the archive; `generate devcontainer` reports that with a warning that is no security relaxation,
  and it and the "no effect" notice go under `--quiet`. The image starts no server and sets neither a listen address nor a port forwarding: the server the user starts listens on loopback only, and a call
  through the container address fails. A `version` override without `source` is rejected with exit code `10`. With the egress restriction the default allowlist gains `registry.ollama.ai` and
  `dd20bb891979d25aebc8bec07b2b3bbc.r2.cloudflarestorage.com` (measured for two models at one time and place; whether the second name stays the same is **not measured**). Models survive removing the
  container, and they survived a rebuild of the image in the measurements below. Measured with Docker 29.5.2 on Colima 0.10.3 (arm64) and Ollama 0.40.6, and in **two runs on one host** on
  2026-10-19: Linux x86_64 (`amd64`), native Docker Engine 29.9.1 without a VM, kernel 6.8.0, `kernel.apparmor_restrict_unprivileged_userns=0`, Ollama 0.40.6. On that host the generated
  Dockerfile built (download of the archive 1.44 GB, image 1.01 GB against 0.95 GB without the block), the server listened on loopback only and reported `library=cpu`, `smollm2:135m` and
  `qwen2.5:0.5b` were pulled and run behind the egress restriction (in the first run through exactly the two hosts above, which were recorded there; the second run pulled successfully, its hosts were not recorded), `ollama.com`, `example.com` and `huggingface.co` stayed blocked, the models survived
  removing the container and a rebuild with Dev Container CLI 0.89.0 (`devcontainer up --remove-existing-container --build-no-cache`), and the combination with `nestedRuntime: podman` and with
  `nestedRuntime: docker` built, left the package list unchanged and ran a small nested container while the server and a pull ran. This is no CI run: the two Docker integration tests of the feature ran in each
  run with a visible pass on that host; a run in CI is **not measured**. **Not measured:** Docker Desktop, Linux behind a VM, larger models, whether the second host stays the same over time and
  region; not provided:
  GPU acceleration, a connection to an Ollama of the host, managing models
  ([`LH-FA-DEV-003`](spec/lastenheft.md#lh-fa-dev-003--devcontainer-features),
  [`LH-NFA-SEC-004`](spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte),
  [`LH-FA-DEV-008`](spec/lastenheft.md#lh-fa-dev-008--egress-restriktion),
  [`LH-FA-CLI-005`](spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging),
  [ADR-0024](docs/plan/adr/0024-ollama-im-sandbox-devcontainer.md)).
- `u-boot --help` links the user manual (German) next to the README: the manual of the installed release (`blob/v<version>`) and the current manual (`blob/main`).
  A build whose version is not a release version (`X.Y.Z` or `X.Y.Z-rc.N`), such as a development or self-built binary, shows only the link to `main`;
  the tag link of a release candidate can lead nowhere once that pre-release tag has been deleted
  ([`LH-FA-CLI-002`](spec/lastenheft.md#lh-fa-cli-002--hilfeausgabe)).
- `u-boot egress list` shows the egress allowlist and the origin of each host (`base`, `runtime`, `docker-apt`, `feature:<name>`, `clone`, `user`; a host that several origins
  contribute lists all of them). The list is the one `u-boot generate devcontainer` writes into the start script: both come from one computation, and a test compares them.
  The command reads only `u-boot.yaml` and `.git/config`, needs no network, no Docker and no running container, and changes nothing. It states whether the restriction is active (sandbox
  profile and `devcontainer.sandbox.egress.enabled`) and, if it is not, shows the list that would apply on activation. The output also says that a name includes its subdomains and that the list
  is computed from the project files, not read from a container; that note carries its measurement limits (Dev Container CLI 0.89.0 on Colima, arm64; the editor actions are not measured).
  `--json` emits `data.{active,sandboxProfile,egressEnabled,cloneSource,hosts}`; the output stays under `--quiet`;
  `--dry-run` and `--diff` are rejected with exit code `2`. The clone source is resolved as for `generate devcontainer --sandbox`, also while the restriction is not active: an `origin`
  URL that this resolution refuses (credentials, unsafe characters) is refused with exit code `10`
  ([`LH-FA-DEV-008`](spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-CLI-005`](spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)).

### Changed

- A wildcard entry in `devcontainer.sandbox.egress.allow` (`*.example.org`) is still rejected with exit code `10`; the error now says that subdomains of an allowed name are already
  included and names the host to enter (`example.org`). The user manual states the measured subdomain behaviour
  ([`LH-FA-DEV-008`](spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)).
- `--dry-run` no longer logs a write as done. The INFO lines on stderr that report a write of `config set` and `generate` now read in the conditional when nothing was
  written: `config set: would update`, `config set (many): would update`, `generate: would create`, `generate: would update block`,
  `generate: would repair Unreleased header`, `generate devcontainer: would create`, `generate devcontainer: would update block` and
  `generate devcontainer: would update allowlist`. Without `--dry-run` (also with `--diff` alone) the lines keep their past tense (`config set: updated`, `generate: created`, …); levels, fields,
  the stdout summary, `--json` and the exit codes are unchanged. A script that matched the old stderr wording of a dry run must match the new one;
  `--json` is the contract for machines ([`LH-FA-CLI-007`](spec/lastenheft.md#lh-fa-cli-007--dry-run)).

### Fixed

- The user manual no longer says that restarting the container applies a changed `devcontainer.sandbox.egress.allow` list. The list is part of the startup script that the image
  build copies into the image, so a restart runs the old script with the old list; rebuild the container after `u-boot generate devcontainer`. Measured with the Dev Container CLI 0.89.0
  (stop and `devcontainer up` kept the old list, `devcontainer up --remove-existing-container` applied the new one) on Docker 29.5.2 with Colima 0.10.3, arm64; **not measured:** the editor
  actions that correspond to a rebuild, Docker Desktop, Linux without a VM ([`LH-FA-DEV-008`](spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)).
- The user manual no longer says without qualification that an entry of `devcontainer.sandbox.egress.allow` leaves its parent name blocked. It separates what a test pins (the status of a connection by name: `example.org` and `amazonaws.com` cover the named subdomains, `www.example.org` does not open `example.org`; the test cannot tell a name that does not resolve from a dropped address) from what was measured in one run (with the entry `api.github.com`, `github.com` and `www.github.com` did not resolve; with the entry `example.org`, `example.org` and `www.example.org` had the same IP addresses). **Not measured:** the entry `www.example.org` itself; the statement about `example.org` is by analogy with `api.github.com`. **Not checked:** a connection to `example.org` by address when the two names share an address. Measured on Docker 29.5.2 with Colima 0.10.3, arm64; **not measured:** other environments. The Ollama section of the user manual says that a call through the container address fails instead of "is refused", because the test cannot tell a refusal from a drop ([`LH-FA-DEV-008`](spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)).
- The user manual states two measured limits of the egress restriction for a process without `sudo`: it can ask the upstream resolver directly and receives the address of a name outside the list
  (the connection to that address stays blocked; name resolution was measured, not the throughput of a data channel), and it can replace the status file, because the status directory belongs to
  `vscode`. Changing the rules still needs `sudo`. Measured once on Docker 29.5.2 with Colima 0.10.3, arm64; **not measured:** other environments ([`LH-FA-DEV-008`](spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)).
- The user manual no longer says that a nested run on `amd64` is not measured, and it states two limits of the egress restriction that two runs measured. On Linux x86_64 with a native
  Docker Engine 29.9.1 (2026-10-19, two runs, no CI run): a nested `podman run` and a nested `docker run` of a small image worked with and without the restriction. Running the startup script
  again in a container with an active restriction (two measurements of 30 runs) left the restriction off in 27 and 28 runs for a measured span of up to 48 and 50 ms (a lower bound; a connection attempt to the container's gateway on port 9 was refused by the gateway instead of dropped) and made name resolution of an
  allowed name fail for 30 to 76 ms in 30 and 27 runs; no connection by a blocked name succeeded. A nested Docker container reached `example.com` and a blocked address
  during the restriction, a nested Podman container did not (two runs; the manual already named nested traffic as outside the restriction). The name resolution of a nested Podman container is not restricted either (see the lead). Not changed: the restriction stays a guardrail
  ([`LH-FA-DEV-007`](spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)).
- The user manual says that a cloud model of Ollama (a model name ending in `-cloud`) cannot be run in the sandbox devcontainer: `ollama pull` succeeded and the model appeared in
  `ollama list`, but `ollama run` failed because the name `ollama.com` did not resolve (the host is not on the allowlist, and the feature does not add it), and `ollama signin` without an account reported
  `account unavailable`; signing in with an account, and whether a cloud model works with an account, was not examined. Measured once on Linux x86_64 with a native Docker Engine 29.9.1 and Ollama 0.40.6 (2026-10-19, behind
  the egress restriction); **not measured:** other environments ([`LH-FA-DEV-003`](spec/lastenheft.md#lh-fa-dev-003--devcontainer-features),
  [`LH-FA-DEV-008`](spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)).

## [0.11.0] - 2026-10-18

Eleventh release. **Podman 6.1.4 in the sandbox devcontainer, a missing `/dev/fuse` as a start error, host-side disclosure of the degradation fallbacks, "Would …" dry-run wording.**
`nestedRuntime: podman` now carries Podman 6.1.4 from the pinned, signature-checked `podman-static` archive; a missing `/dev/fuse` on the container host is a start error of the container
engine (like `/dev/net/tun`, so both devices are host prerequisites), and `u-boot doctor` reports it. `u-boot doctor` and `u-boot generate devcontainer` disclose the degradation fallbacks
(no Docker daemon, `vfs` storage, no egress restriction) from the host; the `--dry-run` summaries of `remove`, `generate` and `config set` say "Would …".
**Behaviour change for scripts that evaluate `doctor` without `--strict`:** on a Linux host outside a container, a missing `/dev/fuse` is now an error (exit code `11`) where `v0.10.0` reported
a warning (exit code `0`) with `onUnavailable: warn`; inside a container (runtime image, devcontainer) the `/dev/fuse` check is skipped (info, exit code `0`; the check of `kernel.apparmor_restrict_unprivileged_userns` before it still applies). **Existing projects** get Podman 6.1.4 and the
status files after the next `u-boot generate devcontainer` and a rebuild of the container. Measured with Docker 29.5.2 on Colima 0.10.3, arm64; **not measured:** Linux without a VM, `amd64`
(for Podman 6.1.4 only the download, checksum, signature and version query are checked; a nested run on `amd64` is not measured), Docker Desktop, VS Code. Podman 6 supports only cgroups v2 (read from its release notes, not measured), so hosts with cgroups v1 are affected. Lastenheft 0.13.0.

### Added

- Sandbox devcontainer: the degradation fallbacks (container without Docker daemon, `vfs` storage, no egress restriction) are now disclosed from the host. The container start scripts
  (`sandbox-init.sh`, `docker-init.sh`, `egress-init.sh`) write a small status file in the container (`/var/lib/u-boot/status/vfs`, `daemon`, `egress`; exit codes and stderr lines are unchanged),
  and `u-boot doctor` reads it through the container engine (`docker ps` and `docker cp`, never `docker exec`) and reports one item per configured capability:
  `devcontainer.sandbox.fallback.vfs`, `devcontainer.sandbox.fallback.daemon` and `devcontainer.sandbox.fallback.egress`. A reported `fallback` is a warning that stays visible under `--quiet`,
  `failed` (policy `fail`) an error (exit code `11`), a state that cannot be read an info line ("not determinable", never a fallback). After writing, `u-boot generate devcontainer`
  warns when an **existing** container reported `fallback` at its last start (also under `--quiet`; not with `--dry-run`; silent without engine, container or status file); `u-boot init` asks the
  same way after writing, but the `u-boot.yaml` it writes configures no nested runtime and no egress restriction, so it has nothing to ask for. The status file comes from a container that runs an agent: it is a hint, not proof, and "not reported" does not mean "not occurred".
  **Existing projects** get the status files after the next `u-boot generate devcontainer` and a rebuild of the container; until then the state is "not determinable". The generated Dockerfile creates
  `/var/lib/u-boot/status`; `runArgs` and mounts do not change. Lastenheft 0.13.0 is unchanged ([`LH-FA-DEV-007`](spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer),
  [`LH-FA-DEV-008`](spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-CLI-005`](spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [ADR-0022](docs/plan/adr/0022-fallback-ausweisung-statusdateien-im-container.md)).

### Changed

- Sandbox devcontainer (`nestedRuntime: podman`): the container now carries **Podman 6.1.4** instead of the Debian package (5.4.2), from the statically linked release archive of
  the community project `podman-static` (Podman, `crun`, `conmon`, `netavark`, `aardvark-dns`, `passt`, `fuse-overlayfs` under `/usr/local`). Version, SHA-256 per architecture and the primary
  fingerprint of the signing key are pinned in the generated Dockerfile; the image build stops before the extraction when the checksum, the key import or the signature does not match. Debian
  still provides `uidmap` and `ca-certificates` only. The archive is a binary package from a single community maintainer and pipeline, not a script; it is downloaded at image build time, so the
  default egress allowlist does not change. `registries.conf` searches Docker Hub only (`unqualified-search-registries = ['docker.io']`); short names other than `alpine` and `mycompany/img` are not
  measured, and fully qualified names stay the recommendation. The `vfs` fallback (container started without `/dev/fuse`) works with Podman 6: the start script writes a user drop-in
  (`~/.config/containers/storage.conf.d/99-u-boot-vfs.conf`) and `image_copy_tmp_dir` in the storage volume; `podman run`, `podman build` and `docker build` run under `vfs` (measured with the volume
  mounted), and the status file still reports `fallback`. Measured with Docker 29.5.2 on Colima 0.10.3, arm64; for `amd64` the download, checksum, signature and version query are checked under
  emulation, a nested run on `amd64` and Docker Desktop are **not measured**. Podman 6 supports only cgroups v2 (its release notes for `v6.0.0`; read, not measured), so **hosts with cgroups v1 are affected**:
  `u-boot doctor` does not check the cgroup version of the container host. **Existing projects** get the new version after the next `u-boot generate devcontainer` and
  a rebuild of the container. Lastenheft 0.13.0 is unchanged ([`LH-FA-DEV-007`](spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer),
  [ADR-0023](docs/plan/adr/0023-podman-aus-podman-static-im-sandbox-devcontainer.md)).
- Sandbox devcontainer (`nestedRuntime: podman`): a missing `/dev/fuse` on the container host is a start error of the container engine, like `/dev/net/tun`, and no longer a
  degradation step: the generated container carries `--device=/dev/fuse`, and where the device is missing the engine refuses to start it (measured with Docker 29.5.2 on
  Colima 0.10.3, arm64: `docker run` exit code `127`, `docker start` exit code `1`, `error gathering device information while adding custom device`); the start script never runs.
  The `vfs` storage fallback applies only to a container that was started without the device (for example after changed `runArgs`). `u-boot doctor` reports a missing `/dev/fuse` on a Linux
  host outside a container as an **error** (exit code `11`) under both `onUnavailable` values, without the fallback flag; it does not check `/dev/net/tun`. Inside a container (the runtime image,
  a devcontainer) the `/dev/fuse` check is skipped (info, exit code `0`; the check of `kernel.apparmor_restrict_unprivileged_userns` that precedes it in the same item applies unchanged, so a value of `1` is still an error, exit code `11`, there too): the device there is the container's own, not the container host's, and a devcontainer started without it shows in
  `devcontainer.sandbox.fallback.vfs` as far as its status file is readable (otherwise "not determinable"). Compared to `v0.10.0`, where a missing device inside a container was a warning ("fall back to vfs", exit code `0`)
  or, with `onUnavailable: fail`, an error (exit code `11`), the `/dev/fuse` check no longer reads the container's device at all (info, exit code `0` under both values, as long as the AppArmor check before it passes).
  **Regression for scripts that evaluate `doctor` without `--strict`, on a Linux host outside a container:**
  the check `devcontainer.sandbox.runtime` was a warning (exit code `0`) with `onUnavailable: warn` and now fails with exit code `11`; load the `fuse` kernel module on the container host (`modprobe fuse`)
  or set `nestedRuntime: none`. The generated files do not change. Lastenheft 0.13.0 ([`LH-FA-DEV-007`](spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer),
  [`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006--exit-codes), [ADR-0022](docs/plan/adr/0022-fallback-ausweisung-statusdateien-im-container.md)).
- `--dry-run`: the summary of `remove`, `generate` and `config set` (one pair or several) says what would happen, like `init` and `add` already did:
  `Would remove service "x".` with the list `Would change:`, `Would generate <artifact> (<paths>).`, `Would update <artifact> managed block (<paths>).`,
  `Would repair <artifact> structure (<paths>).` and `config: Would set <path> <old> → <new>.`. Without `--dry-run` (also with `--diff` alone) the lines keep their
  past form. The "already up to date; no changes." lines, the `--json` output and errors are unchanged; `--quiet --dry-run` still shows the
  summary (the plan is the result of the call), now with the `Would …` wording. A script that matched the old
  dry-run line must use `--json` ([`LH-FA-CLI-007`](spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-NFA-USE-004`](spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe)).

## [0.10.0] - 2026-10-17

Tenth release. **Devcontainer features `sshd` and `claude-code`, `podman build` in the sandbox devcontainer, a new `--quiet` contract.**
The built-in catalogue gains `sshd` and `claude-code`; `nestedRuntime: podman` builds images, which adds a sixth security relaxation
(`--device=/dev/net/tun`) that **existing `nestedRuntime: podman` projects get after the next `u-boot generate devcontainer`** (the container host needs the
device); `--quiet` now keeps four security and fallback disclosures visible and drops the success summaries of the subcommands that change something;
Docker Desktop is no longer promised for `nestedRuntime: docker`. Lastenheft 0.12.0.

### Added

- Devcontainer features: the built-in catalogue gains `sshd` (the official `sshd` feature, an SSH server on port 2222 with password
  and key login and root login permitted, the accounts stay locked until you set a password or add a key; no options, no egress host). It is allowed in the sandbox profile; `generate` discloses the active
  login service with a warning that stays visible under `--quiet`, and the port is not forwarded automatically
  ([`LH-FA-DEV-003`](spec/lastenheft.md#lh-fa-dev-003--devcontainer-features),
  [`LH-FA-DEV-006`](spec/lastenheft.md#lh-fa-dev-006--sandbox-profil),
  [`LH-FA-CLI-005`](spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging),
  [ADR-0020](docs/plan/adr/0020-drittanbieter-features-im-katalog.md)).
- Devcontainer features: the built-in catalogue gains `claude-code` (the Claude Code CLI from the Anthropic feature
  `ghcr.io/anthropics/devcontainer-features/claude-code:1`, a third-party feature that `generate` reports; its install script runs at image build; no options). It needs the
  `node` feature: on the base image the build fails without Node.js, and `generate` says so. In the sandbox profile it is allowed; the generated
  `remoteEnv` then carries the reference `${localEnv:ANTHROPIC_API_KEY}` next to `GIT_TOKEN`, never a value, and with the egress restriction the
  allowlist gains `registry.npmjs.org` and `api.anthropic.com` only
  (the sign-in hosts and the hosts of telemetry, error reports, MCP connectors, updates and plugins are not included; add them to `devcontainer.sandbox.egress.allow`)
  ([`LH-FA-DEV-003`](spec/lastenheft.md#lh-fa-dev-003--devcontainer-features),
  [`LH-FA-DEV-009`](spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer),
  [ADR-0020](docs/plan/adr/0020-drittanbieter-features-im-katalog.md)).

### Changed

- `devcontainer.sandbox.nestedRuntime: docker` is promised on Docker under Linux and on Colima (Linux VM) only.
  Docker Desktop is no longer promised: it is not measured, and the user manual lists it as not verified.
  Homebrew is documented for macOS only (Apple Silicon and Intel); the tap smoke now runs on `macos-14` and `macos-15-intel`.
- The descriptions of the Linux packages and of the Homebrew formula are English ([`LH-LESE-002`](spec/lastenheft.md#lh-lese-002--sprache)),
  and the formula has no `version` line: Homebrew reads the version from the URLs.
- `devcontainer.sandbox.nestedRuntime: podman` builds images: `podman build` (and `docker build` through the alias) builds a Dockerfile with
  `FROM`, `RUN` and `COPY` in the container, cold, from the layer cache and with `--no-cache`, with and without the egress restriction
  ([`LH-FA-DEV-007`](spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer),
  [ADR-0021](docs/plan/adr/0021-podman-build-im-sandbox-devcontainer.md)). A `RUN` needs the network of the nested container, which needs
  `/dev/net/tun`: the generated `runArgs` gain a sixth security relaxation, `--device=/dev/net/tun`, and `generate` reports it like the other five.
  Measured with Docker 29.5.2 on Colima 0.10.3 (arm64, Ubuntu 24.04 VM) with `kernel.apparmor_restrict_unprivileged_userns=0` in the VM and `/dev/net/tun`
  on the container host, with the base images `alpine:3` and `debian:stable-slim`; with the value `1` (the default of an Ubuntu 24.04 VM)
  nested Podman does not start (exit code `11`). **Not measured:** Docker Desktop, amd64 and CI. Two integration tests (`TestE2E_SandboxDevcontainer_NestedPodman_Build`, `…_BuildEgress`) cover it and skip on such hosts.
  Lastenheft 0.12.0 promises it on the measured hosts. `buildah`, `docker buildx`, `docker compose` and BuildKit frontends are not part of the container; use fully qualified image names.
  **Existing `nestedRuntime: podman` projects:** after the next `u-boot generate devcontainer` the container needs `/dev/net/tun` on the container host;
  where the device is missing the engine refuses to start the container (Docker: `error gathering device information while adding custom device`), there is no degradation step and `u-boot doctor` does not check it.
  Load the `tun` kernel module on the host (`modprobe tun`) or set `nestedRuntime: none`.
- The default egress allowlist with `nestedRuntime: podman` or `docker` gains `production.cloudfront.docker.com`: Docker Hub redirects layer downloads
  to it. Measured with `nestedRuntime: podman`: without the host a pull behind the restriction failed with a DNS error. With `nestedRuntime: docker`
  the host is on the same allowlist line; the pull through the Docker daemon is not measured
  ([`LH-FA-DEV-008`](spec/lastenheft.md#lh-fa-dev-008--egress-restriktion)).
- The contract for `--quiet` ([`LH-FA-CLI-005`](spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), Lastenheft 0.10.0) names four exceptions that stay visible: the
  disclosure of security relaxations in the sandbox profile, the disclosure of degradation fallbacks (container without daemon, `vfs`
  storage, dropped egress restriction), the per-service progress of `u-boot up` and the path summary that `--force` prints before
  overwriting. A script cannot silence them. Success and summary messages of the subcommands that change something disappear under
  `--quiet`; `u-boot doctor` lists only error items, degradation fallbacks and the warn items that concern a security property of the sandbox profile (see the sandbox profile entry below); the query commands (`config get`, `config list`, `logs`, `template list`) keep their output. The levels themselves are unchanged (`--quiet` shows only errors). `init`, `generate`, `add`, `remove` and `config set` no longer print their success summary under
  `--quiet` (`remove` also drops its purge warnings), and `u-boot doctor --quiet` no longer lists warn items other than a
  degradation fallback or a security warning of the sandbox profile (the exit code carries the result). `u-boot doctor --strict` with only warn items still exits with `11`
  and now names the warn items in its error message. `u-boot down --quiet` holds the Compose output back and prints it only if `down` fails. `--quiet --json` is unchanged. The output of `--dry-run` and `--diff`
  is the result of the call and stays under `--quiet`: `--dry-run --quiet` prints the plan (the same summary as without `--quiet`; notices other than the
  exceptions above stay dropped) and `--diff --quiet` prints the diff, for `init`, `generate`, `add`, `remove` and `config set`; the success message after an actual write
  is still dropped, also with `--diff` and without `--dry-run`. Of several verbosity options (`--quiet`, `--verbose`, `--debug`) the one given last on the command line wins, without rejection: `--quiet --verbose` shows the additional detail and lifts `--quiet` (also for summaries and warnings), `--verbose --quiet` shows errors only. `--debug --quiet` therefore no longer raises the level to debug; it is `--quiet`. The fixed order `--debug`, `--verbose`, `--quiet` is gone.

- The sandbox profile contract grows (Lastenheft 0.11.0): credentials of an agent (for example the key of a coding agent that an
  enabled devcontainer feature uses) never appear in the image, the workspace volume, `u-boot.yaml` or any generated file; they reach the
  container at run time through an environment variable, and `u-boot doctor` checks the source
  ([`LH-FA-DEV-009`](spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)). A login service that is active in the container
  (password and key login, root login permitted) is disclosed in the command output, and that disclosure stays visible under `--quiet`, as do the warn items of
  `u-boot doctor` that concern a security property of the sandbox profile (credential source and egress restriction). The notice of `u-boot generate devcontainer` that `egress.enabled` has no effect without the sandbox profile
  (no restriction was generated) stays visible under `--quiet` as well; the notice that `nestedRuntime` has no effect is a
  configuration message and is dropped. Only `generate devcontainer` produces these notices; the text of the messages and the `--json` output are unchanged ([`LH-FA-DEV-006`](spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-CLI-005`](spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)).
  `u-boot doctor --quiet` now lists those warn items: a plaintext credential in a project file, an https clone without a runtime
  token source, an active agent feature without the reference of its environment variable, the egress restriction together with
  `nestedRuntime: docker`, and `egress.enabled` without the sandbox profile
  (the restriction does not apply then). The warn item that `nestedRuntime` has no effect without the sandbox profile is a
  configuration message and stays hidden under `--quiet`; exit codes, `--strict` and `--json` are unchanged.
  The run-time handover of agent credentials is delivered with the `claude-code` feature (see Added): the existing check
  `devcontainer.sandbox.credentials` (no new check code) also finds an Anthropic key (prefix `sk-ant-`) or a value of
  `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN` or `CLAUDE_CODE_OAUTH_TOKEN` in a project file, and looks for the
  `${localEnv:ANTHROPIC_API_KEY}` reference while `claude-code` is active.

## [0.9.0] - 2026-10-15

Ninth release. **More devcontainer features.** The built-in catalogue gains `python`, `sdkman` and `homebrew`
(the latter a community feature, reported as third party by `u-boot generate devcontainer`), each with the egress hosts it needs.

### Added

- Devcontainer features: the built-in catalogue gains `python` (the official `python` feature; egress hosts `pypi.org` and
  `files.pythonhosted.org`), `sdkman` (an alias for the official `java` feature, which installs
  SDKMAN; `java` and `sdkman` together render it once) and `homebrew` (the community feature
  `ghcr.io/meaningful-ooo/devcontainer-features/homebrew:2`). `u-boot generate devcontainer` reports a third-party catalogue
  feature as a warning (its install script runs at image build). The egress allowlist gets the hosts of all three features; hosts of
  further SDKMAN candidates go into `devcontainer.sandbox.egress.allow`
  ([`LH-FA-DEV-003`](spec/lastenheft.md#lh-fa-dev-003--devcontainer-features),
  [`LH-NFA-SEC-004`](spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte),
  [ADR-0020](docs/plan/adr/0020-drittanbieter-features-im-katalog.md)).

## [0.8.1] - 2026-10-14

Patch release. **Current Docker in the sandbox devcontainer.** With `nestedRuntime: docker` the container now runs Docker's own
packages (Docker 29.9, Buildx 0.38, Compose 5.6) instead of the older Debian archive versions (26.1, 0.13, 2.26).

### Changed

- Sandbox devcontainer: with `nestedRuntime: docker` the Docker daemon, CLI, Buildx and Compose now come from the Docker APT
  repository (`docker-ce`, `containerd.io`, `docker-buildx-plugin`, `docker-compose-plugin`; the image build verifies the pinned
  signing-key fingerprint) instead of the Debian archive, so the container runs the current Docker (29.9, Buildx 0.38, Compose 5.6
  at the time of writing) rather than 26.1. `download.docker.com` joins the default egress allowlist for `docker`; further hosts
  go into `devcontainer.sandbox.egress.allow`
  ([`LH-FA-DEV-007`](spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer),
  [ADR-0019](docs/plan/adr/0019-docker-daemon-aus-docker-apt-repo.md)).

## [0.8.0] - 2026-10-14

Eighth release. **Own Docker daemon in the sandbox devcontainer.** `devcontainer.sandbox.nestedRuntime: docker`
gives the sandbox container an own Docker daemon with BuildKit, Buildx and Compose for repos whose gates need a real
Docker engine (`-o type=cacheonly`, `RUN --network=none`, `docker compose`, `--privileged` dind sidecars); it requires
`--privileged`, which the command output reports as a relaxation. Lastenheft 0.5.0.

### Added

- Sandbox devcontainer: `devcontainer.sandbox.nestedRuntime: docker` provides an own Docker daemon in the container
  (Debian packages `docker.io`, `docker-cli`, `docker-buildx`, `docker-compose`; BuildKit, Buildx and Compose) next to
  `podman` and `none`. The daemon starts at container start, its storage lives in a volume per instance
  (`<project>-docker-${devcontainerId}`), and the container user reaches it through the `docker` group; no host socket is
  mounted. `docker` requires `--privileged`; the command output reports it as a relaxation and notes that the profile no
  longer isolates the container from its host. A daemon that is not ready at container start degrades per
  `devcontainer.sandbox.onUnavailable` (warning without daemon, or exit code `11`); `u-boot doctor` warns for `docker`
  without the sandbox profile and with the egress restriction
  ([`LH-FA-DEV-007`](spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), Lastenheft 0.5.0,
  [ADR-0018](docs/plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md)).

### Security

- `build(deps): Go toolchain 1.27.1 → 1.27.2` — fixes three HIGH stdlib advisories (CVE-2026-78667 `net/http` Range
  headers, CVE-2026-78669 HTTP/2 SETTINGS frames, CVE-2026-97031 `crypto/tls` ECH) that failed the Trivy image scan.

## [0.7.0] - 2026-10-01

Seventh release. **Linux packages.** Every release now carries `.deb` and `.rpm` packages
(amd64, arm64) built with `nfpm`, next to the binaries and the Homebrew formula; the Keycloak
acceptance test is back in the default integration lane.

### Added

- Debian and RPM packages: every release now carries `.deb` and `.rpm` packages (amd64, arm64) built with
  `nfpm` and listed in `SHA256SUMS`; `make packages` builds them locally. Install with
  `apt install ./u-boot_<version>_amd64.deb` or `dnf install ./u-boot-<version>-1.x86_64.rpm`
  ([`LH-OPEN-002`](spec/lastenheft.md#lh-open-002--paketierung),
  [ADR-0017](docs/plan/adr/0017-linux-pakete-nfpm.md)).

### Changed

- The Keycloak acceptance test ([`LH-AK-003`](spec/lastenheft.md#lh-ak-003--keycloak-flow)) is back in the default `make test-docker` lane
  (build tag `docker` only; the `acceptance_extended` opt-in is gone). The original failure was a
  transient Quay.io outage (2026-06-01), not a defect
  ([`slice-v1-keycloak-ci-flake`](docs/plan/planning/done/slice-v1-keycloak-ci-flake.md)).

## [0.6.0] - 2026-10-01

Sixth release. **Devcontainer sandbox profile and a broad CLI cleanup.**
`u-boot init --devcontainer --sandbox` / `generate devcontainer --sandbox` render a
devcontainer for autonomous agents (named-volume workspace per instance, cloned from
`origin` or `devcontainer.sandbox.repository`; optional rootless Podman with every security
relaxation reported; opt-in DNS-driven egress allowlist; git credentials only at runtime;
adjustable container UID) — [`LH-FA-DEV-004`](spec/lastenheft.md#lh-fa-dev-004--benutzerrechte),
[`-006`](spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)..[`-009`](spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer),
[ADR-0012](docs/plan/adr/0012-devcontainer-egress-firewall.md),
[ADR-0014](docs/plan/adr/0014-nested-podman-sandbox-devcontainer.md),
[ADR-0015](docs/plan/adr/0015-sandbox-volumes-pro-instanz.md). The V1 cleanup backlog is
done: `config get/set` with several paths (atomic), `config list`, structured error hints,
`logs` with several services / format flags / time range, `up` recreate warnings and partial
snapshots, `down` volume names, a real `remove --purge`, and a rollback-aware
`generate devcontainer`. Homebrew joins the distribution paths
([ADR-0016](docs/plan/adr/0016-homebrew-distribution-per-tap.md)); the Go toolchain moves to 1.27.1
(closes eight HIGH stdlib advisories). Details below.

### Added

- Homebrew distribution ([ADR-0016](docs/plan/adr/0016-homebrew-distribution-per-tap.md), [`LH-OPEN-002`](spec/lastenheft.md#lh-open-002--paketierung)):
  the `publish` workflow now attaches `SHA256SUMS` and, for stable tags, a Homebrew
  formula (`u-boot.rb`, filled from those digests by `scripts/homebrew-formula-fill.sh`)
  to every release, and a `tap` job pulls it into `pt9912/homebrew-u-boot` once the
  tap and the `HOMEBREW_TAP_GITHUB_TOKEN` secret exist. Install:
  `brew tap pt9912/u-boot && brew trust pt9912/u-boot && brew install u-boot`
  (macOS/Linux, amd64/arm64; Homebrew 5 requires the explicit `brew trust`).
- `u-boot up` / `down` / `remove --purge` extensions: `up` warns before Compose
  recreates a container (`docker compose --dry-run up`), reports `ports[]`
  structurally and, when it fails mid-way, the partial service snapshot;
  `down --volumes` names the removed volumes; `remove --purge` now really removes
  the service's named volume (best effort, one warning per volume that is in
  use) ([`slice-v1-recreate-detection`](docs/plan/planning/done/slice-v1-recreate-detection.md),
  [`-multi-port-services`](docs/plan/planning/done/slice-v1-multi-port-services.md),
  [`-up-partial-snapshot-on-failure`](docs/plan/planning/done/slice-v1-up-partial-snapshot-on-failure.md),
  [`-down-volumes-named-list`](docs/plan/planning/done/slice-v1-down-volumes-named-list.md),
  [`-volume-auto-removal`](docs/plan/planning/done/slice-v1-volume-auto-removal.md);
  [`LH-FA-UP-003`](spec/lastenheft.md#lh-fa-up-003--startstatus-anzeigen), [`LH-FA-UP-004`](spec/lastenheft.md#lh-fa-up-004--umgebung-stoppen), [`LH-FA-ADD-007`](spec/lastenheft.md#lh-fa-add-007--service-entfernen)).
- `u-boot logs` extensions: several services (`logs api db`), `--no-log-prefix`,
  `--timestamps`, and a time range `--since` / `--until` (duration or
  timestamp; invalid → exit 2), passed through to `docker compose logs`
  ([`slice-v1-logs-multi-service-filter`](docs/plan/planning/done/slice-v1-logs-multi-service-filter.md),
  [`-format-flags`](docs/plan/planning/done/slice-v1-logs-format-flags.md),
  [`-time-range-filter`](docs/plan/planning/done/slice-v1-logs-time-range-filter.md);
  [`LH-FA-UP-005`](spec/lastenheft.md#lh-fa-up-005--logs-anzeigen)).
- `u-boot config` extensions: `config get <p1> <p2>…` and `--json-array`
  (`data.entries[]`), atomic `config set <p1> <v1> <p2> <v2>…`, new
  `config list` (all set paths, sorted), and a structured `data.hint
  {command, action, argument}` on the `config get`/`set` error envelopes
  ([`slice-v1-config-structured-hint`](docs/plan/planning/done/slice-v1-config-structured-hint.md),
  [`-multi-path-get`](docs/plan/planning/done/slice-v1-config-multi-path-get.md),
  [`-multi-path-set`](docs/plan/planning/done/slice-v1-config-multi-path-set.md),
  [`-list-subcommand`](docs/plan/planning/done/slice-v1-config-list-subcommand.md);
  [`LH-FA-CONF-001`](spec/lastenheft.md#lh-fa-conf-001--projektkonfiguration), [`LH-FA-CONF-005`](spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern)).
- Egress restriction for the sandbox devcontainer ([`LH-FA-DEV-008`](spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), V2, [ADR-0012](docs/plan/adr/0012-devcontainer-egress-firewall.md)):
  `devcontainer.sandbox.egress.enabled` / `.allow`; DNS-driven allowlist
  (local `dnsmasq` + nftables set), default hosts per stack, `NET_ADMIN`
  reported as a relaxation, degradation via `onUnavailable`.
- Sandbox profile for devcontainers (Lastenheft 0.3.1,
  [`LH-FA-DEV-006`](spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)): `u-boot init --devcontainer --sandbox` /
  `u-boot generate devcontainer --sandbox` render a named-volume workspace
  per devcontainer instance (`${devcontainerId}` in `mounts`, [ADR-0015](docs/plan/adr/0015-sandbox-volumes-pro-instanz.md))
  that is cloned from `devcontainer.sandbox.repository` or else `origin` (no host bind mount, no socket mount, no
  `--privileged`); no remote → no clone step and a warning; `origin` URLs
  with credentials are rejected (exit 10).
- `devcontainer.user.uid` ([`LH-FA-DEV-004`](spec/lastenheft.md#lh-fa-dev-004--benutzerrechte)): adjustable container-user
  UID via the `USER_UID` build argument.
- `devcontainer.sandbox.nestedRuntime: podman` and
  `devcontainer.sandbox.onUnavailable: warn | fail`
  ([`LH-FA-DEV-007`](spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [ADR-0014](docs/plan/adr/0014-nested-podman-sandbox-devcontainer.md)): rootless Podman in the
  sandbox devcontainer with each security relaxation reported; startup
  script with `vfs` fallback / exit 11.
- Git credentials at runtime only ([`LH-FA-DEV-009`](spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)): `GIT_TOKEN` pass-through and
  credential helper; new `u-boot doctor` checks
  `devcontainer.sandbox.runtime` and `devcontainer.sandbox.credentials`
  and `devcontainer.sandbox.egress` (16 checks in total). See
  [`docs/user/benutzerhandbuch.md#411-einen-sandbox-devcontainer-für-autonome-agenten-einrichten`](docs/user/benutzerhandbuch.md#411-einen-sandbox-devcontainer-für-autonome-agenten-einrichten).

### Changed

- `u-boot generate devcontainer` is now rollback-aware: if one of the multi-file
  writes (`devcontainer.json`, `Dockerfile`, helper scripts, `u-boot.yaml`) fails,
  the previous state is restored (new files and a newly created `.devcontainer/`
  are removed); a failing restore is reported as `rollback incomplete`
  ([`slice-v2-generate-devcontainer-rollback-aware-write`](docs/plan/planning/done/slice-v2-generate-devcontainer-rollback-aware-write.md)).
- `build(deps): Go toolchain 1.26.6 → 1.27.1`, `golangci-lint v2.12.2 → v2.14.0` and `govulncheck v1.1.4 → v1.8.0`
  (routine pins; the old govulncheck panics on Go 1.27, `go.mod` stays at `go 1.26.0`; no code changes needed).
- `build(deps): Go toolchain 1.26.5 → 1.26.6` — fixes eight HIGH stdlib
  advisories (e.g. CVE-2026-33818, `encoding/asn1` denial of service) that
  failed the Trivy image scan.
- Removed the internal `driving.AddPreviewMode` type-alias; the canonical
  `driving.PreviewMode` is the sole spelling of the preview-mode enum. No
  external consumers exist (self-contained CLI).

## [0.5.0] - 2026-07-25

Fünftes Release. **Sicherheits-Fix und lokale Templates.** Der
Go-Toolchain-Bump schließt CVE-2026-39822 im ausgelieferten Runtime-Image
(HIGH, `os.Root` Directory-Traversal) — das ist der eigentliche Anlass des
Cuts. Fachlich neu: `u-boot init --template ./pfad` löst eine Projektvorlage
aus dem echten Dateisystem auf, nicht mehr nur aus dem eingebauten Katalog,
mit plattformunabhängiger Name-vs-Pfad-Klassifikation, geteiltem
`template.yaml`-Parser, Symlink-Schutz und harmonisierten Exit-Codes
([`LH-FA-TPL-003`](spec/lastenheft.md#lh-fa-tpl-003--eigene-templates),
[ADR-0009](docs/plan/adr/0009-template-format-yaml-files.md)). Details unten.

### Security

- `build(deps): Go-Toolchain 1.26.4 → 1.26.5` — behebt **CVE-2026-39822**
  (`os.Root` Symlink-Following / Directory-Traversal in der Go-`stdlib`,
  Severity HIGH, von Trivy im Runtime-Image gemeldet). Routine-Pin-Bump nach
  [`LH-FA-BUILD-003`](spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik)
  (`GO_VERSION` in `Dockerfile` + `Makefile`), kein Verhaltens-Delta.

### Added

- `feat(template): u-boot init --template ./pfad` — lokale
  User-Templates ([`LH-FA-TPL-003`](spec/lastenheft.md#lh-fa-tpl-003--eigene-templates), [ADR-0009](docs/plan/adr/0009-template-format-yaml-files.md) §Entscheidung „Lokale
  User-Templates"). `--template` löst jetzt neben Katalog-Namen
  (`basic`) auch Dateisystem-Pfade auf (`./mein-tpl`, `/abs/tpl`,
  `~/tpl`). Die Klassifikation ist eine reine, plattformunabhängige
  `domain.ClassifyTemplateRef`-Regel (kein FS-Stat); ein Composite-
  Resolver delegiert an den eingebauten `embed.FS`-Katalog oder den
  neuen Filesystem-Resolver (`localtemplates`). Gleiche `template.yaml`-
  Validierung (apiVersion-Gate + Metadaten-Minimum) wie der Katalog
  über das geteilte `templateyaml`-Paket. Fehlerklassen: fehlender
  Pfad / kein Verzeichnis / fehlende `template.yaml` → Exit 10
  (`ErrTemplateNotFound`); malformed `template.yaml` → Exit 10
  (`ErrTemplateInvalid`); Symlink im Template-Baum → Exit 10
  (Pfad-Safety-Reject, kein Teil-Output); Render-Fehler → Exit 14.
  In einem lokalen `template.yaml` deklarierte `variables:` werden
  beim Render aktuell **ignoriert** (keine Substitution, kein
  Prompt); Variable-Auflösung (`--var key=value`) ist out-of-scope
  und folgt in einem eigenen Slice.

## [0.4.0] - 2026-06-08

Fourth release. Completes the V1 machine-readable-CLI milestone: the
full [`slice-v1-cli-json-dry-run`](docs/plan/planning/done/slice-v1-cli-json-dry-run.md) cluster delivers `--json` /
`--dry-run` / `--diff` for **all ten** spec-enum subcommands
([`LH-NFA-USE-004`](spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) Minimalkontrakt + [`LH-FA-CLI-007`](spec/lastenheft.md#lh-fa-cli-007--dry-run)/[`LH-FA-CLI-008`](spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) Voll-Schema),
followed by the `jsonArgsValidator`-consolidation for add/init/
generate. Plus the `u-boot logs` subcommand, devcontainer-features
support with a drift-doctor check, and the transitional reject-gate
teardown at cluster T_close. [ADR-0010](docs/plan/adr/0010-kein-http-driving-adapter.md)-Re-Eval-Trigger-2 (JSON-CLI
as the canonical machine interface) is satisfied. Details below.

### Added

- `feat(cli): u-boot config/config get/config set --json`
  ([`LH-FA-CLI-007`](spec/lastenheft.md#lh-fa-cli-007--dry-run)/[`LH-FA-CLI-008`](spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) / [`LH-NFA-USE-004`](spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) / [`LH-FA-CONF-001`](spec/lastenheft.md#lh-fa-conf-001--projektkonfiguration)..[`LH-FA-CONF-005`](spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern)) —
  achter Folge-Slice (8/9) des Cluster-Slice
  [`slice-v1-cli-json-dry-run`](docs/plan/planning/done/slice-v1-cli-json-dry-run.md). **Erster Read-only+Modifying-
  Hybrid**: drei Sub-Formen teilen `command: "config"` mit
  Pflicht-`subcommand` (`"show"`/`"get"`/`"set"`, §322). Bare +
  `get` sind read-only (nur `--json`, Minimal+Data-Carrier
  `{body}` bzw. `{path, value}`); `set` ist modifying mit
  `--dry-run`/`--diff` (Voll-Schema, `plannedFiles[]` =
  genau eine Zeile `u-boot.yaml`, `hunks[]` aus patched-vs-
  current Bytes). **NoOp**: `oldValue == newValue` → kein
  WriteFile, `plannedFiles: []` + `data.noOp: true` + leeres
  `diagnostics: []` (kein `level: "info"`, Spec §2.1).
  **`--dry-run`/`--diff`-Reject auf den Read-only-Formen** via
  neuem `cli.ErrDryRunNotApplicable` → Exit 2 (T0-(g),
  Pattern-Erbe logs' `ErrFollowJSONNotSupported`). **Drei-Klassen-
  Sentinel-Split** (T0-(m)): `ErrConfigValueInvalid` aufgeteilt in
  `driving.ErrConfigWriteRejected` (non-writable Pfad, Hint
  `u-boot add <svc>`) + `driving.ErrConfigPostPatchSanityFailed`
  (Post-Patch-Roundtrip) — beide Exit 10, damit JSON-Konsumenten
  per `code` disambiguieren statt per Message-Substring.
  **Mapper-Tabelle mit 10 Rows** (T0-(f), FS-first): Row 1
  `ErrConfigFileSystem`→[`LH-NFA-REL-003`](spec/lastenheft.md#lh-nfa-rel-003--abbruch-bei-kritischen-fehlern)/Exit 14, Rows 2-3
  Schema/Post-Patch-Sanity→[`LH-FA-CONF-002`](spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration), Rows 4-5-7
  [`LH-FA-CONF-005`](spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern) (Multi-Use: Path-Unknown/Write-Rejected/
  Value-Not-Set — Disambiguation per Message-Prefix), Row 6
  `ErrConfigValueInvalid`→[`LH-FA-CONF-001`](spec/lastenheft.md#lh-fa-conf-001--projektkonfiguration), Row 8
  `ErrProjectNotInitialized`→[`LH-FA-INIT-001`](spec/lastenheft.md#lh-fa-init-001--neues-projekt-initialisieren) (Environment-
  Operation Pattern-Erbe up/down/generate/logs), Row 9 Reject
  →Exit 2, Row 10 Default. **PreviewMode-Cluster**:
  `ConfigService.fsFactory` + `selectFS` + `NewConfigServiceWith
  Factory` + `cmd/uboot/main.go`-Wiring (fünfter Preview-Factory);
  `config set` routet seinen WriteFile über die mode-spezifische
  FS und surface't `ConfigSetResponse.PlannedFiles` aus dem
  Recorder für den `--diff`-Renderer. **SilenceLogger** silenced
  die fünf `s.logger.*`-Sites im JSON-Mode; **Orphan-Feature-WARN**
  ([`LH-FA-DEV-003`](spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)/`level: "warn"`) migriert in `diagnostics[]`
  (Dual-Emission). **Subcommand-bewusste `reportErrorSub`/
  `writeErrorEnvelopeSub`** (additiv; Single-Form-Caller
  unverändert) erfüllen die §322-Pflicht auch auf dem Error-Pfad.
  Allowlist-Reject-Liste schrumpft von 4 auf 1 (nur noch
  `template (bare)`). Drei Review-Runden (zwei HIGH + ein MED, alle
  gefixt). Doku in [`docs/user/cli-json-output.md §6.9`](docs/maintainer/cli-json-contract.md).
- `feat(cli): u-boot logs --json` ([`LH-FA-CLI-007`](spec/lastenheft.md#lh-fa-cli-007--dry-run) /
  [`LH-NFA-USE-004`](spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) / [`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006--exit-codes)) — siebter Folge-Slice
  (7/9) des Cluster-Slice [`slice-v1-cli-json-dry-run`](docs/plan/planning/done/slice-v1-cli-json-dry-run.md).
  **Read-only-Klasse** auf lokalem FS (analog up-down):
  weder `--dry-run` noch `--diff` — nur `--json` mit
  typisiertem Data-Carrier. **`logsStatusData.lines []string`
  ohne omitempty** (Empty-Array-Pin: `[]`, NICHT `null`,
  Pattern-Erbe up-down's `services []serviceStatus` ohne
  omitempty). **T0-(a) Single-Envelope + `--follow --json`
  Reject** (Option (A)): Spec-§1841-Konsens (Single-Envelope
  pro CLI-Call) wird honoriert; `--follow --json` ist
  inkompatibel und wird in `runLogs` Stage-1 (vor UC-Call)
  mit neuem `ErrFollowJSONNotSupported` → [`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006--exit-codes)/
  Exit 2 rejected — bounded `--tail=N`-Pfad ist die
  einzige Akquisitions-Form. **T0-(i) Validation-Order**:
  `--follow --json` schlägt `--tail=-1`; CLI-Stage-1
  reihenfolge pinnt Reject-Sentinel VOR Tail-Validation
  (`TestLogsJSON_ValidationOrder_FollowJSONBeatsInvalidTail`
  Pin). **Neuer FS-Sentinel** `driving.ErrLogsFileSystem`
  mit Read-spezifischer Message-Form (`"logs: filesystem
  read failed"`, Pattern-Erbe up-down). T3 wrapt zwei
  FS-Read-Stellen (`logsservice.go:117/137` —
  `checkProjectInitialized` + `checkComposeFile`) auf
  Multi-`%w` (`fmt.Errorf("logs service: Exists(%q): %w: %w",
  path, driving.ErrLogsFileSystem, err)`). **Mapper-Tabelle
  mit 9 Rows** (T0-(f)): Row 1 FS-Sentinel-first (LH-NFA-REL-
  003/Exit 14), Rows 2-3 shared `mapComposeRuntimeSentinel`-
  Helper (Docker/ComposeRuntime → [`LH-NFA-REL-003`](spec/lastenheft.md#lh-nfa-rel-003--abbruch-bei-kritischen-fehlern)), Row 4
  shared `ErrComposeFileMissing` ([`LH-FA-UP-001`](spec/lastenheft.md#lh-fa-up-001--umgebung-starten)/Exit 10),
  Row 5 cross-cutting `ErrProjectNotInitialized` (LH-FA-INIT-
  001/Exit 10, Pattern-Erbe generate als Environment-
  Operation), Row 6 domain-level `ErrInvalidServiceName`
  ([`LH-FA-INIT-006`](spec/lastenheft.md#lh-fa-init-006--projektnamen-validierung)/Exit 10), Row 7 logs-only
  `ErrFollowJSONNotSupported` ([`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006--exit-codes)/Exit 2 — T0-(a)
  Reject-Pfad), Row 8 logs-only `ErrInvalidLogsTail`
  ([`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006--exit-codes)/Exit 2), Row 9 Default [`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006--exit-codes)/Exit 1.
  **Switch-Order-FS-first-Defense** (Pattern-Erbe up-down
  R2-HIGH-2 / R3-HIGH-1) via `TestLogsJSON_MultiWrap_
  FSAndDocker_SwitchOrderFSFirst_ByDesign`-Pin — synthetische
  FS+Docker-Chain → diagnostics[0].code = [`LH-NFA-REL-003`](spec/lastenheft.md#lh-nfa-rel-003--abbruch-bei-kritischen-fehlern)
  (FS-Klasse, Mapper-FS-first), `exitCode = 11` (Docker-
  Sub-Klasse, ExitCode-Helper-Driven-first); **`(code,
  exitCode)`-Tupel-Disambiguation** per `cli-json-output.md
  §6.7` ist der Vertrag (Pattern-Erbe up-down T8).
  **`baseDirSanitizedError`-Wiederverwendung**: `runLogs`
  wrappt UC-Errors mit `sanitizeBaseDir(err, cwd)` vor
  `reportError` (Path-Leak-Defense, Pattern-Erbe up-down
  T5). **`runLogs(ctx, stdout, errOut io.Writer, args,
  flags, uc, getwd)`-Signatur** Cluster-konsistent mit
  up/down/remove (errOut für strukturierte Pfade reserviert,
  heute `_ = errOut`-Stub). **`logsFlags{Follow, Tail,
  Service, JSON, Quiet}`** mit neuen `JSON`/`Quiet`-Boolean-
  Feldern und `IsValid()`-Builder (Pattern-Erbe up's
  `upFlags`). **Cluster-Allowlist** erweitert
  (`jsonallowlist.go`): `"u-boot logs": true`; Reject-Liste
  von 5 → 4 (`config bare/get/set`, `template bare` bleiben).
  **`cli.isFilesystemError`** erweitert um
  `driving.ErrLogsFileSystem` → Exit 14;
  **`cli.isUsageError`** erweitert um
  `cli.ErrFollowJSONNotSupported` → Exit 2 (`cli.go`).
  **15 Acceptance-Pins** (`logs_acceptance_test.go`):
  T0-(a)/T0-(i)/T0-(j)(ii) verbatim, 9 Mapper-Coverage-
  Pins, Empty-Array-Pin, Trailing-Newline-Strip-Pin, Path-
  Leak-Sanitizer-Pin, FS+Docker-Switch-Order-Defense-Pin.
  Pre-T6-Review: HIGH=0, MED=4, LOW=6 (T7 fixte MED-1 +
  LOW-5; MED-3 → §6.8-Doku-Pflicht hier eingelöst).
  Pre-T8-Bestätigungsrunde: HIGH=0, MED=2, LOW=2 (MED-1
  Mapper-Kommentar-Drift + MED-2 Defense-Pin + LOW-1 Plan-
  Drift in `ba7d06f` gefixt; LOW-2 CRLF-Lücke in §6.8 als
  bekannte Limitation dokumentiert). **Vier neue
  open/-Stubs** (T6 R2-LOW): [`slice-v1-logs-format-flags`](docs/plan/planning/done/slice-v1-logs-format-flags.md),
  [`slice-v1-logs-multi-service-filter`](docs/plan/planning/done/slice-v1-logs-multi-service-filter.md),
  [`slice-v1-logs-time-range-filter`](docs/plan/planning/done/slice-v1-logs-time-range-filter.md).
  `[ba7d06f, b502cd5, 343e622, 69cfc0d,
  c21ba28, 0fe74e4]`.

- `feat(cli): u-boot up --json / u-boot down --json`
  ([`LH-FA-CLI-007`](spec/lastenheft.md#lh-fa-cli-007--dry-run) / [`LH-NFA-USE-004`](spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) / [`LH-FA-UP-001`](spec/lastenheft.md#lh-fa-up-001--umgebung-starten)/[`LH-FA-UP-003`](spec/lastenheft.md#lh-fa-up-003--startstatus-anzeigen)/[`LH-FA-UP-004`](spec/lastenheft.md#lh-fa-up-004--umgebung-stoppen) /
  [`LH-FA-CLI-005A`](spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung)) — sechster Folge-Slice (6/9) des Cluster-
  Slice [`slice-v1-cli-json-dry-run`](docs/plan/planning/done/slice-v1-cli-json-dry-run.md). **Read-only-Klasse** auf
  lokalem FS: weder `--dry-run` noch `--diff` (Cluster-Slice
  Z. 464-467) — nur `--json` mit typisierten Data-Carriern.
  `u-boot up` und `u-boot down` sind im selben Slice gebündelt
  weil beide den Compose-Status lesen und das Confirmer-Swap-
  Pattern teilen. **`upStatusData.services[]`** trägt
  `serviceStatus{name, state, port, healthcheck}` (plain Go-
  Strings mit `omitempty` für port/healthcheck — keine Three-
  State-Disambiguation nötig; [`LH-FA-UP-003`](spec/lastenheft.md#lh-fa-up-003--startstatus-anzeigen) Mindestangaben);
  `data.timeoutFireAndForget *bool omitempty` als Marker nur im
  `--timeout=0`-Pfad (Pattern-Erbe remove's `*bool`-Key-
  Absence-Disambiguation). **`downStatusData.removedVolumes
  bool`** ohne omitempty (`false` ist legitimer Success-Wert
  "nichts entfernt"). **`UpRequest.SilenceProgress bool` +
  `DownRequest.SilenceConfirmer bool`** symmetrisch zum
  remove-Pattern; CLI setzt sie auf `flags.JSON`. **Application-
  Layer-ProgressSink-Branch** (`UpService.Up`): `effective :=
  req.ProgressSink; if req.SilenceProgress { effective =
  io.Discard }` — Compose-Phase-Stream wird im JSON-Mode
  unterdrückt, nil-Default bleibt im DockerEngine-Adapter
  (`progressSinkOrDiscard`-Pattern). **Application-Layer-
  Confirmer-Branch** (`DownService.runConfirmationGate` Row 4):
  Request-time Gate-Branch ohne Field-Mutation (`confirmer :=
  s.confirmer; if req.SilenceConfirmer { confirmer =
  noopConfirmer{} }`) — kein neuer `downMu`-Mutex nötig, race-
  frei by construction. **Refuse-by-Default-Semantik** im
  JSON-Mode: bei `--volumes --json` OHNE `--yes` returnt
  `noopConfirmer.ConfirmRemoveVolumes` `(false, nil)` → fällt
  durch in `ErrConfirmationRequired`/Exit 10 (Symmetrie zum
  `--no-interactive`-Pfad). JSON-Konsumenten MUSSEN `--yes`
  explizit setzen für destructive `--volumes`. **Zwei neue
  FS-Sentinels** `driving.ErrUpFileSystem` /
  `driving.ErrDownFileSystem` mit Read-spezifischer Message-Form
  (`"<cmd>: filesystem read failed"`, NICHT `"mutation
  failed"` weil up/down read-only auf lokalem FS). T3 migriert
  fünf FS-Read-Wraps (`upservice.go:105/138/148` +
  `downservice.go:81/97`) auf Multi-`%w`. **Mapper-Tabelle mit
  verbindlicher Switch-Order** (T0-(e) R3-HIGH-1): zehn Rows
  mit Mapper-Heim-Spalte (`mapUp`/`mapDown`/`helper`/`beide`).
  Row 1 FS-Sentinel-first ([`LH-NFA-REL-003`](spec/lastenheft.md#lh-nfa-rel-003--abbruch-bei-kritischen-fehlern)/Exit 14), Rows 2-3
  shared `mapComposeRuntimeSentinel(err)`-Helper in neuem File
  `cli/composesentinel.go` ([`LH-NFA-REL-003`](spec/lastenheft.md#lh-nfa-rel-003--abbruch-bei-kritischen-fehlern) für `driven.ErrDocker
  Unavailable` und `driven.ErrComposeRuntime`), Row 4 up-only
  `ErrStabilizationTimeout`/[`LH-FA-UP-001`](spec/lastenheft.md#lh-fa-up-001--umgebung-starten)/Exit 12, Row 5 down-
  only `ErrConfirmationRequired`/[`LH-FA-INIT-005`](spec/lastenheft.md#lh-fa-init-005--überschreibschutz)/Exit 10
  (geteilt mit init/remove), Row 6 shared `ErrComposeFileMissing`
  /[`LH-FA-UP-001`](spec/lastenheft.md#lh-fa-up-001--umgebung-starten), Row 7 cross-cutting `ErrProjectNotInitialized`/
  **[`LH-FA-INIT-001`](spec/lastenheft.md#lh-fa-init-001--neues-projekt-initialisieren)** (Pattern-Erbe generate als Environment-
  Operation, NICHT add/remove [`LH-FA-ADD-001`](spec/lastenheft.md#lh-fa-add-001--add-on-befehl)), Row 8 up-only
  `ErrInvalidTimeout`/[`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006--exit-codes), Row 9 down-only
  `ErrConflictingModeFlags`/[`LH-FA-CLI-005A`](spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), Row 10 Default
  [`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006--exit-codes). **`(code, exitCode)`-Tupel-Disambiguation**
  (cli-json-output.md §6.7-Doku-Pin): bei Multi-`%w`-Wraps
  liefern Mapper (FS-first) und ExitCode-Helper (Driven-first
  per `cli.go:285-313`) zwei getrennte Klassifikationen —
  Konsument disambiguiert über das Tupel, NICHT über `code`
  allein. Pattern-Erbe remove's [`LH-FA-ADD-007`](spec/lastenheft.md#lh-fa-add-007--service-entfernen) Multi-Use.
  **`baseDirSanitizedError`-Helper-Extraktion**: aus
  `cli/remove.go:465-538` nach neuem File `cli/sanitize.go`
  (package-intern; `remove.go:299` nutzt unverändert weiter).
  `runUp`/`runDown` wrappen UC-Errors mit `sanitizeBaseDir(err,
  cwd)` vor `reportError` — 11 FS-Read- und Compose-Runtime-
  Wraps in upservice/downservice tunneln keinen absoluten
  Filesystem-Pfad mehr in `diagnostic.message`. **Allowlist-
  Migration**: `u-boot up` + `u-boot down` als zwei separate
  Einträge in `jsonAllowlist()` (Cobra-CommandPath-Form);
  Reject-Liste schrumpft von 7 auf 5 (logs, config bare/get/set,
  template bare). Tranchen-Reihenfolge: T1 + T4 entfallen
  (noopConfirmer/io.Discard existieren als Helper; Composition-
  Root braucht kein Wiring-Update bei Bool-Field-Pattern); T2
  Port-Types (`e966a83`), T3 Application-Layer (`86fb5b2`), T5
  CLI-RunE + Mapper + Sanitizer-Helper-Extraktion (`a5aaf9c`),
  T6 28 Acceptance-Tests (`2473988`), T7 Pre-T8-Adressierung
  von 11 Findings (`31f7238`), T8 Closure. **Acceptance-
  Coverage**: 28 CLI-Tests (18 up + 10 down) plus 2 Application-
  Layer-Tests in `downservice_test.go` für die `noopConfirmer`-
  Branch-Defense (`removeVolumesCalls == 0` + Contrast-Pin).
  **Out-of-Scope-V1 Carveouts mit open/-Stubs**: Recreate-
  Detection ([`slice-v1-recreate-detection`](docs/plan/planning/done/slice-v1-recreate-detection.md)), Volume-Named-Liste
  ([`slice-v1-down-volumes-named-list`](docs/plan/planning/done/slice-v1-down-volumes-named-list.md)), Partial-Snapshot bei
  Mid-ComposeUp-Failure
  ([`slice-v1-up-partial-snapshot-on-failure`](docs/plan/planning/done/slice-v1-up-partial-snapshot-on-failure.md)),
  strukturierte Multi-Port-Liste
  ([`slice-v1-multi-port-services`](docs/plan/planning/done/slice-v1-multi-port-services.md)).
  Coverage-Gate 91 %.
- `feat(cli): u-boot remove --json / --dry-run / --diff`
  ([`LH-FA-CLI-007`](spec/lastenheft.md#lh-fa-cli-007--dry-run)/[`LH-FA-CLI-008`](spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) / [`LH-NFA-USE-004`](spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) / [`LH-FA-ADD-007`](spec/lastenheft.md#lh-fa-add-007--service-entfernen) /
  [`LH-FA-CLI-005A`](spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung)) — fünfter Folge-Slice (5/9) des Cluster-Slice
  [`slice-v1-cli-json-dry-run`](docs/plan/planning/done/slice-v1-cli-json-dry-run.md). `u-boot remove` ist die **inverse
  Operation zu `add`** (strip managed-block aus `compose.yaml` +
  `.env.example`, flip `services.<name>.enabled` auf `false` in
  `u-boot.yaml`, optional Volume-Purge via `--purge`-Gate). Acht
  Flag-Kombinationen plus `--purge`-Dimension: drei JSON-Modi
  (Minimal+Data, Voll-Schema Dry-Run, Voll-Schema Preview-and-
  Apply) plus Human-Mode mit/ohne `--diff`/`--purge`. Pattern-
  Erbe von add/init/generate 1:1: `driving.PreviewMode` direkt
  (kein Service-Prefix-Alias), `RecordingFileSystem`-driven-
  Adapter, Pure-Go LCS-Diff-Renderer, `previewModeFromFlags`-
  Mapping, generalisierte Error-Emission-Helper. Remove-spezifisch:
  **Confirmer-Gate für `--purge`** ([`LH-FA-CLI-005A`](spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) §254) mit
  `noopConfirmer`-Swap im JSON-Mode (T0-(j) — neues Pattern,
  nicht aus init-Progress-Swap geerbt: Service-Field-Mutation
  mit defer-Restore innerhalb `removeMu`-Lock-Scope; lokale-
  Variable-Variante verworfen weil `runPurgeGate`-Signature-
  Refactor); **`data.volumesPurged`** als `*bool` im Envelope
  (`false` deferred-Status in v0.3.0 ist valider Success-Wert,
  Plain-`bool`+omitempty würde Error-Pfad-Zero und Success-Pfad-
  `false` identisch droppen); **WARN-Migration**: heutige
  `printRemoveSummary`-stderr-WARNING bei `--purge && !VolumesPurged`
  wandert im JSON-Mode in `diagnostics[]` mit
  code: "[`LH-FA-ADD-007`](spec/lastenheft.md#lh-fa-add-007--service-entfernen)", `level: "warn"` (Multi-Use des Codes
  für ERROR `ErrServiceUnregistered` UND WARN deferred-Volumes —
  Konsumenten disambiguieren über `(code, level)`-Tupel);
  **Custom-`Args`-Validator** `validateRemoveArgs(a *App)` als
  Cobra-PositionalArgs-Closure mit `*App`-Capture (R11/R12/R13-Pin),
  ersetzt `cobra.ExactArgs(1)` — emittiert den
  [`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006--exit-codes)-Envelope auf stdout BEVOR der Sentinel zu Cobra
  zurückgeht (Spec §1841-Symmetrie); **Voll-Schema bei `--dry-run`/
  `--diff`** auch im NoPositionalArg- und TooManyArgs-Pfad
  (Spec §1842, R13-HIGH-1 Validator-Flag-Awareness via
  `cmd.Flags().GetBool`); **`baseDirSanitizedError`-Wrapper** für
  `diagnostic.message`: FS-Wraps der Form
  `fmt.Errorf("... %s: %w: %w", absPath, ErrRemoveFileSystem, raw)`
  tunneln den absoluten Filesystem-Pfad in den User-facing Output —
  Sanitizer ersetzt `<baseDir>/foo` durch `foo` und bare `<baseDir>`
  durch `.`, an Word-Boundaries (`replaceBareBaseDir` ist robust
  gegen Substring-Kollisionen wie `<baseDir>-cache/lock`,
  R14-MED-1 + R15-LOW-1); `errors.Is`/`As` bleiben intakt via
  Unwrap-Chain. **Dry-Run-WARN-Suppression** in
  `printRemoveSummary`: Use-Case skippt `runPurgeGate` in
  `PreviewDryRun` (T0-(h)(a)) und führt keine Mutation aus — die
  WARN-Prosa wäre semantisch falsch ("ist-deferred" statt
  "würde-deferred"); Fix unterdrückt WARN-Block bei
  `previewMode == PreviewDryRun`, `PreviewAndApply` behält die
  WARN. Remove-spezifische LH-Code-Mapper-Tabelle in
  `mapRemoveErrorToDiagnostic` (FS-first Switch-Order):
  [`LH-NFA-REL-003`](spec/lastenheft.md#lh-nfa-rel-003--abbruch-bei-kritischen-fehlern)/14 für `ErrRemoveFileSystem`, [`LH-FA-CLI-005A`](spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung)/10
  für `ErrConfirmerUnavailable` (neuer Sentinel) UND
  `ErrConflictingModeFlags`/2, [`LH-FA-INIT-005`](spec/lastenheft.md#lh-fa-init-005--überschreibschutz)/10 für
  `ErrConfirmationRequired`, `LH-FA-ADD-{001,002,005,007}`/10 für
  fachliche Service-Sentinels, [`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006--exit-codes)/2 für
  `ErrServiceNameMissing` und Cobra-too-many-args. **Idempotenz-
  NoOp-Semantik**: nur `PriorState=Deactivated` qualifiziert
  (Single+Repeat-Call → `plannedFiles: []`, `changes: []`,
  `data.priorState=data.state="deactivated"`, Exit 0);
  `EnabledUnset` und `InconsistentBlock` sind state-transitioning
  (`Changed!=nil`, Voll-`plannedFiles[]`). **`delete`-Action im
  Recorder**: `RemoveAll`-Captures für `extraFiles` werden als
  `plannedFiles[].action: "delete"` gewired — `--diff` rendert
  Old-Content + leeren New-Content. Tranchen-Reihenfolge (T1
  entfällt — `noopConfirmer` lebt bereits in M4 Confirmer-Slice):
  T2 Port-Types (`d0c9c5d`), T3 Application-Layer mit 8 FS-Wrap-
  Stellen Multi-`%w` (`dbbf7b1`), T4 `newPreviewFSFactory`-Helper
  (`3b079dd`), T5 CLI-RunE-Rewrite (`3188e75`), T6+T6-A
  Acceptance-Tests inkl. `WithDataKeyAbsent`/`WithDataKeyPresent`-
  Helper-Erweiterung in `jsontestutil` (`9eae9ec`), T7 Pre-T8-
  Review-Fixes R13-HIGH-1 + R13-MED-1 + R14-HIGH-2 + R14-MED-1
  (`4fb3fea`), T8 Closure mit R15-LOW-1 Sanitizer-Substring-
  Robustheit + R15-LOW-2 Coverage-Pin. Acceptance-Coverage:
  23 Pin-Tests in `remove_acceptance_test.go` (drei JSON-Modi,
  Idempotenz-Repeat, EnabledUnset-Normalisierung, Mid-Write-
  Failure mit `plannedFiles[]`-Capture, ConfirmationRequired-
  Pfad, ServiceUnregistered ERROR vs. deferred-Volumes WARN,
  Concurrent-Invocations-Mutex, Dry-Run-WARN-Suppression mit
  PreviewAndApply-Kontrast, Sanitizer-Path-Leak-Defense plus
  Substring-Kollisions-Robustheit, TooManyArgs-Voll-Schema bei
  `--dry-run`, `delete`-Action-Hunk-Vertrag). Coverage-Gate
  91.10 %.
- `feat(cli): u-boot generate --json / --dry-run / --diff`
  ([`LH-FA-CLI-007`](spec/lastenheft.md#lh-fa-cli-007--dry-run)/[`LH-FA-CLI-008`](spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) / [`LH-NFA-USE-004`](spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) / [`LH-FA-GEN-001`](spec/lastenheft.md#lh-fa-gen-001--generate-befehl)..[`LH-FA-GEN-005`](spec/lastenheft.md#lh-fa-gen-005--idempotenz) /
  [`LH-FA-DEV-001`](spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen)/[`LH-FA-DEV-003`](spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)) — vierter Folge-Slice (4/9) des Cluster-
  Slice [`slice-v1-cli-json-dry-run`](docs/plan/planning/done/slice-v1-cli-json-dry-run.md). `u-boot generate` ist nach
  doctor/add/init der nächste modifying-Subcommand und der
  **erste**, der mehrere Artefakte (changelog/readme/env-example/
  devcontainer) über einen einzigen Subcommand bedient. Vier
  Flag-Kombinationen mit derselben Symmetrie wie add/init:
  `--json` (Minimal+Data-Envelope), `--dry-run --json`
  (Voll-Schema mit `plannedFiles[]`/`changes[]`, kein FS-Write),
  `--diff --json` (Voll-Schema mit Hunks, Preview-and-Apply),
  `--dry-run --diff --json` (Vorschau plus Hunks, kein Write).
  Human-Mode `generate --diff` rendert Unified-Diff am stdout.
  **Action-Klassifikation via `data.action`** (`created` /
  `updated-block` / `no-op` / `repaired-manual`):
  UpdatedBlock und RepairedManual sind FS-semantisch identisch
  (`plannedFiles[i].action: "modify"`), `data.action` ist der
  einzige Discriminator. **Multi-Artefakt-Envelope-Form**:
  `command="generate"`, kein `subcommand`-Feld (Cobra-Positional-
  Arg-Semantik), `data.artifact` trägt das Artefakt. Pattern-
  Erbe init→generate 1:1: `PreviewMode`-Carrier (direkt, kein
  Service-Prefix-Alias), `RecordingFileSystem`-driven-Adapter,
  Pure-Go LCS-Diff-Renderer, `previewModeFromFlags`-Mapping,
  generalisierte Error-Emission-Helper. Generate-spezifische
  Erweiterungen: **2 von 8 Recorder-Mutations-Methoden** im
  Capture-Set (`WriteFile`, `MkdirAll`); **kein** `GitClient`
  und **kein** `ProgressPort` (schmaler als init);
  **`generateMu sync.Mutex`** auf `GenerateService`;
  **per-Artefakt LH-Code-Tabelle** im neuen
  `mapGenerateErrorToDiagnostic(err, artifact)` (changelog→
  [`LH-FA-GEN-002`](spec/lastenheft.md#lh-fa-gen-002--changelog-erzeugen), readme→[`LH-FA-GEN-003`](spec/lastenheft.md#lh-fa-gen-003--readme-erzeugen), env-example→
  [`LH-FA-GEN-004`](spec/lastenheft.md#lh-fa-gen-004--beispiel-env-erzeugen), devcontainer→[`LH-FA-DEV-001`](spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen));
  **`ErrConfigValueInvalid`-Sentinel-Wrap** auf
  `validateAllowExternalFeatureSourcesEntries`/
  `applyAllowExternalFeatureSources` für den
  [`LH-FA-DEV-003`](spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)-URL-Reject-Pfad (Spec §720 fordert exakt
  Exit 10 für ungültige `--allow-external-feature-sources`-URLs;
  ohne den Sentinel-Wrap wäre der Pfad auf Default
  [`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006--exit-codes)/Exit 1 gefallen); **Multi-`%w`-Wrap auf
  ~17 FS-Wrap-Stellen** in `application/generate.go` (Switch-
  Order-Sicherheit analog init T6; ohne Multi-`%w` würde ein
  Multi-Wrap mit FS-Sentinel + fachlich-Sentinel auf Exit 10
  downgraden). **`cliJSONEnvelope.Data`-Feld + `newDataEnvelope`-
  Konstruktor** wurden aus dem Template-Slice 9/9 in generate
  vorgezogen (Generate ist der erste Multi-Artefakt-Konsument
  mit `data`-Bedarf); Template-Slice 9/9 erbt das Feld nur noch.
  **`writeErrorEnvelope`/`reportError` um `data any`-Trailing-
  Param erweitert** (init/add reichen `nil` durch — nicht-
  brechende Erweiterung). Acceptance-Coverage: 15 Tests in
  `generate_acceptance_test.go` (drei JSON-Modi, 4 ManualConflict-
  Codes als Sub-Tests, URL-Reject-[`LH-FA-DEV-003`](spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), ArtifactUnknown-
  Exit-2, ProjectNotInitialized-Exit-10, FS-Failure-Exit-14,
  Allow-External-Mutex, NoOp-Empty-Arrays, UpdatedBlock-vs-
  RepairedManual-Action-Discriminator, Human-Mode-Summary +
  Diff-Rendering). Devcontainer-Phase-1-Atomicity + Phase-2-
  Half-Write-Carveout (V2-Open-Slice
  [`slice-v2-generate-devcontainer-rollback-aware-write`](docs/plan/planning/done/slice-v2-generate-devcontainer-rollback-aware-write.md)) +
  Repeat-Idempotency leben in den Application-Layer-Tests.
  Coverage-Gate ≥ 91 %.

### Changed

- **BREAKING** `u-boot template list --json` Ausgabe-Format
  ([slice-v1-cli-json-dry-run-template](docs/plan/planning/done/slice-v1-cli-json-dry-run-template.md), neunter und letzter
  Folge-Slice des Cluster-Slice [`slice-v1-cli-json-dry-run`](docs/plan/planning/done/slice-v1-cli-json-dry-run.md)):
  der bisherige rohe, pretty-indented `[]templateJSON`-Array-
  Output wurde auf den [`LH-NFA-USE-004`](spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe)-Minimalkontrakt-Envelope
  migriert — `{"status":"ok","command":"template","subcommand":
  "list","diagnostics":[],"exitCode":0,"data":[…]}` (single-line,
  compact). Die Template-Liste lebt jetzt im `data`-Feld;
  Konsumenten, die das Top-Level-Array lasen, müssen auf `.data`
  umstellen. Damit ist das letzte nicht-spec-konforme `--json`-
  Surface geschlossen (alle neun Folge-Slices des Clusters tragen
  jetzt den Minimalkontrakt; der bewusste Doctor-Slice-Carveout ist
  aufgelöst). bare `u-boot template --json` bleibt Exit-2-Reject
  (§1838: `subcommand` verpflichtend für `command="template"`;
  Help-Parent ohne eigenes Datum). Ein Katalog-IO-Fehler mappt auf
  [`LH-NFA-REL-003`](spec/lastenheft.md#lh-nfa-rel-003--abbruch-bei-kritischen-fehlern)/Exit 14. Doku in
  [`docs/user/cli-json-output.md §6.2`](docs/maintainer/cli-json-contract.md).
- **Cluster-T_close** (Abschluss des [`slice-v1-cli-json-dry-run`](docs/plan/planning/done/slice-v1-cli-json-dry-run.md)-
  Clusters): die transitionale `--json`-Reject-Mechanik (Allowlist-
  Map + `applyJSONRejectGate` am Root-`PersistentPreRunE` +
  `ErrJSONNotImplemented`) wurde **entfernt** — alle Spec-Enum-Forms
  sind migriert, es gibt nichts mehr zu rejecten. Der bare-
  `u-boot template --json`-Reject ist von gate-getragen
  (`ErrJSONNotImplemented`) auf **RunE-getragen**
  (`ErrTemplateSubcommandRequired`, Exit 2, envelope-LOS) umgestellt,
  damit er den Gate-Abbau überlebt ohne Hilfetext zu leaken. Keine
  Verhaltensänderung für Konsumenten (Exit 2 bleibt); rein interner
  Mechanik-Abbau (netto −90 LOC). Doku in
  [`docs/user/cli-json-output.md §6.1`](docs/maintainer/cli-json-contract.md).

### Fixed

- **`add`/`init`/`generate` `--json` Envelope-Symmetrie + Path-Leak-
  Defense** ([slice-v1-cli-json-envelope-consolidation](docs/plan/planning/done/slice-v1-cli-json-envelope-consolidation.md), R15-Cross-
  Slice-1): ein Wrong-Arg-Aufruf unter `--json` (`u-boot --json add`
  ohne Service, `u-boot --json add a b` mit zu vielen) emittierte
  bisher nur eine nackte Cobra-stderr-Meldung und **kein** JSON auf
  stdout (Spec §1841-Verletzung); `--dry-run`/`--diff` wählten zudem
  nicht das Voll-Schema (§1842). Jetzt tragen alle drei Commands den
  Args-Envelope (Exit 2) über den geteilten `jsonArgsValidator` —
  konsolidiert mit `config`/`remove`, die das Muster schon hatten.
  Zusätzlich: absolute Filesystem-Pfade aus Use-Case-Fehler-Wraps
  werden nun via `sanitizeBaseDir` aus `diagnostic.message`
  entfernt (Path-Leak/Info-Disclosure-Defense, symmetrisch zu
  `config`/`remove`). Kein Verhaltenswechsel für korrekte Aufrufe.
- `fix(cli): mapAddErrorToDiagnostic Backup-Sentinels auf`
  [`LH-NFA-REL-003`](spec/lastenheft.md#lh-nfa-rel-003--abbruch-bei-kritischen-fehlern) (add↔init Diagnostic-Code-Harmonisierung) —
  `mapAddErrorToDiagnostic` mappte `ErrBackupSuffixExhausted`/
  `ErrBackupSourceMissing` auf [`LH-FA-INIT-005`](spec/lastenheft.md#lh-fa-init-005--überschreibschutz) (Validation-Klasse,
  Exit-10-Suggestion), während `isFilesystemError` sie ohnehin
  zu Exit 14 routet — Envelope-Code und Exit-Klasse waren also
  desynchron. `mapInitErrorToDiagnostic` macht es bereits korrekt
  ([`LH-NFA-REL-003`](spec/lastenheft.md#lh-nfa-rel-003--abbruch-bei-kritischen-fehlern) + Exit 14). Cleanup zieht add auf dieselbe
  Klassifikation nach. Die Branch ist im add-Pfad defensiv (heute
  ruft kein add-Use-Case `runBackup`/`BackupPath`), bleibt aber
  für zukünftige Catalog-Erweiterungen erhalten. Cross-Slice-
  Drift-Finding aus [`slice-v1-cli-json-dry-run-init`](docs/plan/planning/done/slice-v1-cli-json-dry-run-init.md) Review-Round-9
  (`d7f9e65`); Plan-Pfad
  [`slice-v1-cli-cleanup-add-backup-error-class`](docs/plan/planning/done/slice-v1-cli-cleanup-add-backup-error-class.md).
  Coverage-Gate unverändert grün (91.10 %).

### Added

- `feat(cli): u-boot init --json / --dry-run / --diff`
  ([`LH-FA-CLI-007`](spec/lastenheft.md#lh-fa-cli-007--dry-run)/[`LH-FA-CLI-008`](spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) / [`LH-NFA-USE-004`](spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe)) — zweiter modifying-Sub-
  command mit JSON-Envelope-Migration und wichtigster Onboarding-
  Use-Case (Cluster-Slice [`slice-v1-cli-json-dry-run`](docs/plan/planning/done/slice-v1-cli-json-dry-run.md), Folge-Slice
  3/9). Vier neue Flag-Kombinationen mit derselben Symmetrie wie
  `add`: `--json` (Minimalkontrakt-Envelope), `--dry-run --json`
  (Voll-Schema mit `plannedFiles[]`/`changes[]`, kein FS-Write),
  `--diff --json` (Voll-Schema mit `plannedFiles[].hunks[]`,
  Preview-and-Apply), `--dry-run --diff --json` (Vorschau plus
  Hunks, kein Write). Human-Mode `init --diff` rendert Unified-
  Diff am stdout. Pattern-Erbe von add 1:1: gemeinsame
  `PreviewMode`-Carrier-Types (`AddPreviewMode` als Alias auf
  kanonischen `driving.PreviewMode`), `RecordingFileSystem`-
  driven-Adapter, Pure-Go LCS-Diff-Renderer,
  `previewModeFromFlags`-Mapping, generalisierte Error-Emission-
  Helper (`reportError`/`writeErrorEnvelope`/`writeDiff`/
  `lastPlannedPath`) in `cli/erroremission.go`. Init-spezifische
  Erweiterungen: **sechs der acht** `driven.FileSystem`-
  Mutations-Methoden im Capture-Set (`MkdirAll`, `WriteFile`
  direkt plus `CopyExclusive`/`Mkdir`/`MkdirAll`/`Copy`/
  `RemoveAll` indirekt via `BackupPath`); **`initGit`-Skip im
  Dry-Run** (Composition-Root-`initFSFactory`-Closure liefert
  Recorder; `Init()` skippt den separaten `driven.GitClient`-Port
  bei `PreviewMode == PreviewDryRun`, weil git am Recorder vorbei
  auf die echte Disk schreibt); **ProgressPort-Silencing im
  JSON-Mode** (`req.SilenceProgress = flags.JSON` swappt
  `s.progress` auf Noop, weil add keinen stdout-bound Port hat
  und init der erste mit Progress-Events während des Use-Case-
  Laufs ist); **Service-Race-Safe via `sync.Mutex`** auf
  `InitProjectService.initMu` (analog `AddServiceService`).
  Sieben init-spezifische LH-Codes als Spec-Anker
  ([`LH-FA-INIT-001`](spec/lastenheft.md#lh-fa-init-001--neues-projekt-initialisieren)..[`LH-FA-INIT-007`](spec/lastenheft.md#lh-fa-init-007--git-repository-initialisierung)),
  drei mit dedizierten Sentinels in
  `mapInitErrorToDiagnostic` ([`LH-FA-INIT-004`](spec/lastenheft.md#lh-fa-init-004--bestehendes-projekt-erkennen) Marker-Kollision,
  [`LH-FA-INIT-005`](spec/lastenheft.md#lh-fa-init-005--überschreibschutz) Force/Backup-Usage, [`LH-FA-INIT-006`](spec/lastenheft.md#lh-fa-init-006--projektnamen-validierung) Name-
  Validierung); **neuer `driving.ErrInitFileSystem`-Sentinel**
  als Multi-`%w`-Wrap auf fünf FS-Wrap-Stellen (vier direkte in
  `initproject.go` — `MkdirAll` + drei `WriteFile`-Actions — plus
  ein Caller-side Wrap auf das `runBackup`-Ergebnis, der typed
  Backup-Sentinels durchreicht und nur rohe FS-Errors wrappt;
  R3-Wrap-Strategie konsolidiert die ursprünglich pro
  `backup.go`-Site geplanten Wraps auf den Aufruferand).
  Switch-Order-Pflicht: FS-first, weil Multi-`%w` sonst Exit-14
  auf Exit-10 downgraded. Template-Modus (`init --template <name>`) ist in
  V1 **mutex** zu `--dry-run`/`--diff` —
  `driving.ErrTemplateConflictsWithFlag` rejected die Kombination
  (Exit-Code 2, [`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006--exit-codes)); Template-Preview wandert in
  eigenen Folge-Slice. Context-Cancellation bleibt
  Status-quo-Carveout (Cross-Cutting Exit-130-Convention für alle
  modifying-Subcommands ist eigener Block). Acceptance-Coverage:
  ~17 Flag-Matrix-Tests inkl. Soft-Existing × `--devcontainer`,
  Planning-Phase-Force-Failure (Exit 10), Mid-Write-Failure mit
  Switch-Order-Pin (Exit 14), Concurrent-Init-Mutex-Pin auf zwei
  Goroutinen, Path-Anchor-Pin für positional `<name>` mit
  trailing-slash/dot-slash/abs-path, initGit-Skip-Pin (`.git/`
  fehlt + Spy-Counter 0), JSON-stdout-Cleanliness-Pin
  (`json.Decode → io.EOF`). Coverage-Gate ≥ 91 %.

- `feat(cli): u-boot add --json / --dry-run / --diff`
  ([`LH-FA-CLI-007`](spec/lastenheft.md#lh-fa-cli-007--dry-run)/[`LH-FA-CLI-008`](spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) / [`LH-NFA-USE-004`](spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe)) — erster modifying-Sub-
  command mit JSON-Envelope-Migration (Cluster-Slice
  [`slice-v1-cli-json-dry-run`](docs/plan/planning/done/slice-v1-cli-json-dry-run.md), Folge-Slice 2/9). Vier neue
  Flag-Kombinationen: `--json` (Minimalkontrakt-Envelope ohne
  Plan), `--dry-run --json` (Voll-Schema mit `plannedFiles[]`/
  `changes[]`, kein FS-Write), `--diff --json` (Voll-Schema
  Preview-and-Apply mit `plannedFiles[].hunks[]`),
  `--dry-run --diff --json` (Vorschau plus Hunks, kein Write).
  Human-Mode-`--diff` rendert Unified-Diff-String an stdout
  (`+`/`-`/space-Prefix plus `@@`-Header). Cluster-Infrastruktur
  als Pattern-Vorbild: neuer `RecordingFileSystem`-driven-Adapter
  (`internal/adapter/driven/recordingfs/`) implementiert alle 8
  Mutations-Methoden mit Passthrough-Schalter, modelliert den
  impliziten `MkdirAll`-Effekt auf Parent-Dirs; Pure-Go LCS-Diff-
  Renderer (`internal/adapter/driving/cli/diff/`); Composition-
  Root-`fsFactory(driving.AddPreviewMode)`-Closure in
  `cmd/uboot/main.go`. Diagnostic-Codes sind LH-Kennungen
  (`LH-FA-ADD-{001,002,005,006}`/`LH-FA-INIT-{004,006}`/
  [`LH-NFA-REL-003`](spec/lastenheft.md#lh-nfa-rel-003--abbruch-bei-kritischen-fehlern)). Mid-Write-Failure-UX: Voll-Schema-Envelope
  zeigt `plannedFiles[]` bis zur Failure-Stelle, `diagnostics[].file`
  markiert die Failure-Position, `exitCode: 14` ([`LH-NFA-REL-003`](spec/lastenheft.md#lh-nfa-rel-003--abbruch-bei-kritischen-fehlern)).
  `changes[].count`-Semantik: für `create` total Lines der neuen
  Datei, für `modify` `+`-Lines aus den Hunks (Spec §477 exakt),
  für `delete` `0`, für binary `CountBytesDiff`. Service-Race-
  Safe via `sync.Mutex` auf `AddServiceService`; recursive
  Dep-Installs (für künftige Catalogue-Erweiterungen) erben den
  Outer-`PreviewMode`. AssertFullEnvelope erweitert um
  `checkHunks`-Helper mit Field-Name-Drift-Pin
  (`offset` statt `oldStart` failt sofort).

- `feat(cli): u-boot doctor --json plus Root-PersistentFlag
  --json` ([`LH-NFA-USE-004`](spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) / [`LH-FA-CLI-007`](spec/lastenheft.md#lh-fa-cli-007--dry-run)) — Pattern-Vorbild für
  die maschinen-lesbare CLI-Surface (Cluster-Slice
  [`slice-v1-cli-json-dry-run`](docs/plan/planning/done/slice-v1-cli-json-dry-run.md), Folge-Slice 1/9). **Doctor**
  emittiert mit `--json` einen Spec-§1841-Minimalkontrakt-
  Envelope (`status`/`command`/`diagnostics`/`exitCode`);
  All-OK-Fall ergibt `diagnostics: []` gemäß Lastenheft-Beispiel
  §1846-1852. `SeverityOK` und `SeverityInfo` werden gefiltert
  (Spec §1834 erlaubt `level` nur `warn|error`). `--quiet --json`
  ist semantisch identisch zu reinem `--json`; `--strict --json`
  upgraded Warn auf Exit-Code 11 ohne `status`-Drift (Spec §1837
  koppelt `status` an höchsten `level`, nicht an `--strict`).
  Broken-Pipe-Resistenz: fachlicher Exit-Code 11 hat Vorrang vor
  Write-Fehlern. **Root-PersistentFlag `--json`** für alle 10
  Spec-Enum-Subcommands; nicht-migrierte Forms rejecten mit
  `ErrJSONNotImplemented` (Exit-Code 2, [`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006--exit-codes)-Klasse)
  plus Folge-Slice-Verweis. Zentrale Allowlist in
  `cli/root.go` via `PersistentPreRunE`; `--help` als
  read-only Escape-Hatch durchgelassen. **`u-boot template list
  --json`** wandert vom lokalen Flag aufs Root-Flag (beide
  Schreibweisen `template list --json` und `--json template list`
  identisches Output); Envelope-Migration folgt mit
  [`slice-v1-cli-json-dry-run-template`](docs/plan/planning/done/slice-v1-cli-json-dry-run-template.md), dokumentiert im
  `carveouts.md`-Temporär-Eintrag. **Common-Envelope
  `cliJSONEnvelope`** (`internal/adapter/driving/cli/jsonenvelope.go`)
  trägt Minimalkontrakt- und Voll-Schema-Felder über zwei
  Konstruktoren (`newMinimalEnvelope`/`newFullEnvelope`); Voll-
  Schema-Felder via Pointer-Wrapping (`*bool`/`*[]T`) als
  Anti-Drift gegen `omitempty`-Semantik-Refactor. **Schema-
  Helper-Sub-Package** `internal/adapter/driving/cli/jsontestutil/`
  mit zwei Modi `AssertMinimalEnvelope`/`AssertFullEnvelope`
  (Options-Pattern, kein neuer Dep) plus `DefaultAllowedCodes`-
  Registry für die 13 Doctor-Check-Codes. **Drei aktive Drift-
  Gates** schützen die Code-Registry: (1) `application.DoctorCheckIDs()`
  ↔ Map-Vollständigkeit, (2) Map ↔ Markdown-Roundtrip-Parser
  auf `docs/user/cli-json-output.md` §5.1 mit HTML-Marker-
  Sektion-Begrenzung, (3) Helper-Reject im Acceptance-Pfad
  für undokumentierte Codes. **Schema-Vertrag-Doku**
  ([`docs/user/cli-json-output.md`](docs/maintainer/cli-json-contract.md))
  zitiert Minimalkontrakt und Voll-Schema verbatim, dokumentiert
  Code-Registry und Per-Command-Migrations-Reihenfolge.
- `feat(logs): u-boot logs [service] [--follow] [--tail <n>]`
  ([`LH-FA-UP-005`](spec/lastenheft.md#lh-fa-up-005--logs-anzeigen)) — neuer Subcommand streamt Compose-Logs als
  V1-Erweiterung der `up`/`down`-Familie. Ohne Service-Argument
  läuft `docker compose logs` über alle Services aus
  `compose.yaml` (T0-(a) Compose-Facade-Semantik, kein
  u-boot.yaml-Filter); mit Service-Argument nur diesen einen
  (Format-Validation via `domain.NewServiceName`, Existenz-Check
  delegiert an Compose). `--follow` blockt bis Ctrl-C und
  beendet via SIGINT-Vertrag-Schicht-1+2+3 (Adapter gibt
  `ctx.Err()` unverdeckt zurück, Use-Case übersetzt zu
  `(LogsResponse{}, nil)`, CLI exit-code 0). `--tail <n>` mit
  Stage-1-Validation auf nicht-negative Ganzzahlen (Default
  leer → Use-Case-Normalisierung zu Compose-`"all"`). Exit-
  Code-Mapping analog `up`/`down`: 10 (User/Project-State),
  11 (`ErrDockerUnavailable`), 12 (`ErrComposeRuntime`),
  14 (FS), 2 (CLI-Usage). Docker-tag E2E-Tests gegen echten
  postgres-Stack pinnen `--tail`-Buffer-Content und
  `--follow`-SIGINT-Vertrag.
- `feat(devcontainer): Devcontainer-Features-Allowlist und Katalog`
  ([`LH-FA-DEV-003`](spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)) — 8 Built-in-Features (`git`, `docker-cli`,
  `node`, `java`, `go`, `cpp`, `kubectl-helm`, `postgres-client`)
  plus External-Source-Allowlist via
  `devcontainer.featureSources.allow`. CLI:
  `u-boot config set devcontainer.features.<name>.{enabled,source,version}`
  plus `--allow-external-feature-sources <url>[,<url>...]` auf
  den drei Spec-§714-717-Pfaden (`init --devcontainer`,
  `generate devcontainer`, `config set devcontainer.featureSources.allow`).
  Doctor-Check `devcontainer.features.allowlist` (Error bei
  Allowlist-Violation, Warn bei Orphan-Activation oder fehlendem
  `enabled:`-Key). User-Doku in
  [`docs/user/benutzerhandbuch.md#410-devcontainer-features-nutzen`](docs/user/benutzerhandbuch.md#410-devcontainer-features-nutzen).
- `feat(doctor): devcontainer.features.drift Check` — über-Spec
  Drift-Erkennung zwischen `u-boot.yaml`'s Features-Map und den
  Keys im gerenderten `.devcontainer/devcontainer.json`. Drei
  Warn-Cases: aktiviertes Feature fehlt im JSON (Case 1, inkl.
  Datei-fehlt-Disziplin), deaktiviertes Feature noch im JSON
  (Case 2a), JSON-Key ohne cfg-Pendant (Case 2b, Hand-Edit-Hint).
  Repair-Hint `u-boot generate devcontainer`. Doctor-Total
  steigt 12→13.

## [0.3.0] - 2026-06-01

Third release. Completes the V1 „Add-on Catalogue Expansion"
milestone (5/5): the catalogue now ships three integrated
service add-ons — Postgres (since MVP), Keycloak ([`LH-FA-ADD-003`](spec/lastenheft.md#lh-fa-add-003--keycloak-hinzufügen)
/ [`LH-AK-003`](spec/lastenheft.md#lh-ak-003--keycloak-flow)) and OpenTelemetry ([`LH-FA-ADD-004`](spec/lastenheft.md#lh-fa-add-004--opentelemetry-hinzufügen) / [`LH-AK-004`](spec/lastenheft.md#lh-ak-004--opentelemetry-flow)) —
plus the matching `u-boot remove <service>` mirror, the
[`LH-FA-ADD-006`](spec/lastenheft.md#lh-fa-add-006--add-on-abhängigkeiten) `--with-deps` dependency-resolution mechanism, and
a doku-only audit closure for three V1 spec-IDs. Architectural
side-effect: the per-service catalogue pattern grew from a flat
`(compose, env, volume)`-tuple in M5 to a declarative entry with
`requiredEnvKeys` / `volumeRefLiteral` / `volumeOptional` /
`healthcheckOptional` / `extraFiles` — any new add-on plugs in
by adding one catalogue row and three templates.

### Verified

- **Three V1 spec-IDs audit-closed** by
  [`slice-v1-audit-done`](docs/plan/planning/done/slice-v1-audit-done.md)
  — Doku-only verification that the existing code/doc state
  already satisfies the requirements:
  [`LH-FA-BUILD-006`](spec/lastenheft.md#lh-fa-build-006--aggregator-targets) (Aggregator-Targets `gates`/`ci`/`fullbuild` in
  the Makefile),
  [`LH-NFA-MAINT-004`](spec/lastenheft.md#lh-nfa-maint-004--dokumentierte-schnittstellen) (Add-on and template interfaces documented via
  [ADR-0008](docs/plan/adr/0008-plugin-system-statisch.md)/-0009 + driving/driven port doc-comments + slice docs),
  [`LH-NFA-PORT-003`](spec/lastenheft.md#lh-nfa-port-003--containerfreundlichkeit) (u-boot itself runs in container / devcontainer:
  GHCR distroless image + container-aware `doctor` since v0.2.0 +
  six-platform binary distribution + `init --devcontainer`-generated
  files).

### Added

- **`u-boot add otel`** — [`LH-FA-ADD-004`](spec/lastenheft.md#lh-fa-add-004--opentelemetry-hinzufügen) / [`LH-AK-004`](spec/lastenheft.md#lh-ak-004--opentelemetry-flow). Third and
  final add-on of the v0.3.0 milestone catalogue. Compose-Service
  mit Image-Pin `otel/opentelemetry-collector:0.108.0` (Stable),
  Port-Mappings `4317:4317` (OTLP/gRPC) + `4318:4318` (OTLP/HTTP),
  `command: --config=/etc/otel-collector-config.yaml`, Bind-Mount
  der gerenderten Config-Datei. Kein Healthcheck im Mindest-
  Setup ([`LH-AK-004`](spec/lastenheft.md#lh-ak-004--opentelemetry-flow) §2374 toleriert `running` ODER `healthy`).
  Mindest-Collector-Config in `otel-collector-config.yaml`:
  Receivers `otlp/grpc+http`, Processors `batch`, Exporters
  `debug` (stdout), Pipelines `logs`/`metrics`/`traces` — alle
  drei Signal-Typen aus [`LH-FA-ADD-004`](spec/lastenheft.md#lh-fa-add-004--opentelemetry-hinzufügen) §880.
  
  Internal: Catalogue-Pattern wächst um drei Felder pro Service —
  `extraFiles []extraFileEntry` für whole-file artefacts abseits
  von compose+env+volume (für OTel die Collector-Config-Datei),
  plus `envOptional` (implizit via leerem `envTmpl`) und
  `healthcheckOptional` für Services, die das Standard-Pattern
  legitim nicht brauchen. `executeAdd` schreibt extraFiles als
  vierten Slot nach yaml/compose/env; `executeRemove` löscht sie
  symmetrisch. `serviceComplete` skipt healthcheck-presence für
  `healthcheckOptional`; explicit `healthcheck.disable: true`
  bleibt hart abgelehnt. Acceptance-Helper-Reuse aus
  [`slice-v1-keycloak`](docs/plan/planning/done/slice-v1-keycloak.md) T3 (`acceptance_helpers.go`) — OTel-E2E
  bleibt ~30 Zeilen. **Makefile-Patch**: `test-docker`-Target
  mountet jetzt `/tmp` host-shared, damit Compose-Bind-Mount-
  Pfade vom Daemon (Host) aufgelöst werden können — sonst sieht
  der Daemon nur den Container-Pfad `t.TempDir()` nicht und
  erstellt einen leeren Verzeichnis-Mount, der den Collector
  beim Config-Read crasht. See
  [`slice-v1-otel`](docs/plan/planning/done/slice-v1-otel.md).
- **`u-boot add keycloak`** — [`LH-FA-ADD-003`](spec/lastenheft.md#lh-fa-add-003--keycloak-hinzufügen) / [`LH-AK-003`](spec/lastenheft.md#lh-ak-003--keycloak-flow). Second
  add-on in the catalogue after Postgres. Compose-Service mit
  Image-Pin `quay.io/keycloak/keycloak:26.0` (LTS), Port-Mapping
  `8080:8080`, `command: start-dev` für [`LH-AK-003`](spec/lastenheft.md#lh-ak-003--keycloak-flow)-Boot, Healthcheck
  via `/dev/tcp/localhost/9000` (bash-builtin, kein curl im
  Image) gegen `/health/ready`. Admin-Credentials via Placeholder-
  Env-Block (`KEYCLOAK_ADMIN=CHANGEME_KEYCLOAK_ADMIN` +
  `KEYCLOAK_ADMIN_PASSWORD=CHANGEME_KEYCLOAK_ADMIN_PASSWORD`).
  **Persistenz: flüchtige H2-In-Container-Datenbank** — kein
  Volume, nach `docker compose down` weg; [`LH-AK-003`](spec/lastenheft.md#lh-ak-003--keycloak-flow) verlangt nur
  Endpoint-200/302. Persistente externe Postgres-Anbindung
  ([`LH-FA-ADD-003`](spec/lastenheft.md#lh-fa-add-003--keycloak-hinzufügen) §857) bleibt als eigener Folge-Slice
  (`slice-v1-keycloak-external-postgres`, Trigger: Nutzer-Bedarf).  <!-- d-check:ignore (angekuendigter Folge-Slice, noch ohne Plan-Datei) -->
  Internal refactor: `renderPostgresTemplates` → generischer
  `renderServiceTemplates(svc)` über neue Service-Catalogue-
  Tabelle; `hasRequiredEnvKeys` / `contentScanState` /
  `inspectVolumeArtefact` / `patchTargetsFor` werden per-Service
  über `requiredEnvKeys` / `volumeRefLiteral` / `volumeOptional`
  parametrisiert, damit Keycloak's volume-loser Pfad nicht in
  den Postgres-Repair-Loop läuft. Test-Helper-Extraktion:
  `internal/e2e/acceptance_helpers.go` teilt die init+add+up-
  Pipeline mit dem [`LH-AK-002`](spec/lastenheft.md#lh-ak-002--postgresql-flow)-Postgres-Test (Boot-Zeit-Carveout
  für Keycloak: 4 min UpService-Timeout vs. 90 s Postgres). See
  [`slice-v1-keycloak`](docs/plan/planning/done/slice-v1-keycloak.md).
- **`u-boot add <service> --with-deps`** — [`LH-FA-ADD-006`](spec/lastenheft.md#lh-fa-add-006--add-on-abhängigkeiten) add-on
  dependency mechanism. New domain type `AddOnDependency` (path-
  conditional service dependency declaration) + per-service
  catalogue side-table `dependenciesFor(svc)` (Postgres has none
  today; first non-nil row lands with [`slice-v1-keycloak`](docs/plan/planning/done/slice-v1-keycloak.md)). When
  the requested add-on declares a dep that is not yet registered
  in `u-boot.yaml`, the four-mode dispatch decides what happens:
  `--with-deps` auto-installs the chain (recursive `Add` calls,
  flag inherited so transitive deps follow); `--yes` has the same
  effect; `--no-interactive` (without `--yes`/`--with-deps`)
  fails fast with the new `ErrDependenciesRequired` sentinel
  (exit 10); default-interactive prompts via the new
  `Confirmer.ConfirmAddDependency(ctx, svc, missing)` driven-port
  method (mirror of `ConfirmRemoveVolumes` from M6). Postgres-
  only flows are unchanged — the no-deps short-circuit keeps the
  load+resolve cost out of the MVP catalogue path. Breaking
  refactor in the application layer: `NewAddServiceService`
  now takes a `Confirmer` between `yaml` and `logger`; all eight
  callsites updated in lock-step. See
  [`slice-v1-addons-deps`](docs/plan/planning/done/slice-v1-addons-deps.md).
- **`u-boot remove <service> [--purge]`** — first slice of the
  v0.3.0 milestone ("Add-on Catalogue Expansion"). Mirror of
  `u-boot add`: detects the [`LH-FA-ADD-005`](spec/lastenheft.md#lh-fa-add-005--mehrfaches-hinzufügen-verhindern) service state, strips
  the `service.<name>` managed block from `compose.yaml` and
  `.env.example`, then sets `services.<name>.enabled: false` in
  `u-boot.yaml`. Idempotent: removing an already-disabled service
  is a no-op with a clear message. Inconsistent project state
  (orphan block, missing entry) surfaces as `ErrServiceInconsistent`
  with a manual-cleanup hint. New driving sentinel
  `ErrServiceUnregistered` (exit 10) distinguishes "service was
  never added" from "service name not in the catalogue"
  (`ErrServiceUnsupported`). [`LH-FA-ADD-007`](spec/lastenheft.md#lh-fa-add-007--service-entfernen) §"Volumes nur auf
  explizite Anforderung": `--purge` opts in destructively and
  triggers the [`LH-FA-CLI-005A`](spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) §254 confirmation gate (mirror of
  `u-boot down --volumes`); auto-removal of volumes is deferred
  to a follow-up slice — v0.3.0's `--purge` summary points at
  `docker volume rm <name>` for the manual cleanup. Internal:
  `detectServiceState` extracted from the M5 add path to a
  package-level function so both add and remove share it without
  duplication. See
  [`slice-v1-add-remove`](docs/plan/planning/done/slice-v1-add-remove.md).

## [0.2.0] - 2026-06-01

Second release. Adds the first two V1 template features
(`template list` + `init --template`), a cross-platform binary
distribution (six platforms as GitHub-Release assets), and a
container-aware `doctor` that no longer mis-reports a healthy host
as 4 errors when run from inside the distroless image. v0.1.1 was
originally planned as a patch-only tag for the doctor fix but is
skipped in favour of this minor bump — three features landed before
the tag-push and strict SemVer wants a MINOR bump for them.

### Added

- **`u-boot template list [--json]`** — first V1 template
  subcommand ([`LH-FA-TPL-004`](spec/lastenheft.md#lh-fa-tpl-004--templates-auflisten)). Enumerates the built-in project-
  template catalog with name, description, and version in a
  tabwriter-aligned table; `--json` emits a structured array
  with the full [`LH-FA-TPL-002`](spec/lastenheft.md#lh-fa-tpl-002--template-metadaten) metadata surface (`supportedAddOns`,
  `generatedFiles`, `requiredTools`, `variables`). Bootstrap
  built-in: `basic` (one template; further built-ins follow on
  demand per [ADR-0009](docs/plan/adr/0009-template-format-yaml-files.md) §Folgepunkte 4). Fully hexagonal:
  `domain.TemplateMetadata` + `Validate()` (kebab-case-name
  regex, `ErrInvalidTemplate` sentinel), driven port
  `TemplateCatalog`, embed.FS-backed `externaltemplates` adapter,
  application `TemplateListService` (multi-`%w` so the original
  `domain.ErrInvalidTemplate` chain survives), CLI
  `template list` rendering. Adapter directory consolidated to
  `internal/adapter/driven/externaltemplates/` (no hyphen) for
  consistency with the existing `driven/`-adapter naming; [ADR-0009](docs/plan/adr/0009-template-format-yaml-files.md)
  §Entscheidung updated to match. See
  [`slice-v1-template-list`](docs/plan/planning/done/slice-v1-template-list.md).
- **`u-boot init <name> --template <name>`** — second V1 template
  feature, the render path of [`LH-FA-TPL-001`](spec/lastenheft.md#lh-fa-tpl-001--projektvorlagen) / [`LH-FA-TPL-002`](spec/lastenheft.md#lh-fa-tpl-002--template-metadaten). The
  init service delegates file rendering to the new
  `TemplateInitService` when `--template` is set; project structure
  directories and `git init` stay with the InitProjectService so
  the user-observable flow is one command. Byte-identity
  guarantee: `u-boot init demo --template basic` produces a
  project byte-identical to `u-boot init demo` for the six default
  files (`u-boot.yaml`, `compose.yaml`, `README.md`,
  `CHANGELOG.md`, `.env.example`, `.gitignore`) — pinned by an
  E2E `diff -r` test against the production catalog. Render engine:
  Go `text/template` for `*.tmpl` files, 1:1 copy for non-`.tmpl`
  files (per [ADR-0009](docs/plan/adr/0009-template-format-yaml-files.md) §Entscheidung); `template.yaml` metadata is
  skipped. Two-phase render-then-write: a render error in any file
  short-circuits before the first disk write, so a buggy template
  no longer leaves a half-populated project. New
  `domain.TemplatePath` validator rejects `..` segments, absolute
  paths, Windows drive letters, backslashes, NUL bytes, and empty
  strings ([`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006--exit-codes) exit 10 via `ErrInvalidTemplatePath`).
  Mutex with `--devcontainer`/`--force`/`--backup`: surfaces as
  `ErrTemplateConflictsWithFlag` (exit 2) — v1 is fresh-init-only.
  Soft-existing-detection is skipped on the template path because
  `--template` resolves the "is this an existing project?"
  ambiguity by intent; the hard-existing check
  (`u-boot.yaml` present → `ErrProjectExists`) remains the
  safety net. Variable resolution + `--var key=value` deferred to
  a future slice (basic has no variables). See
  [`slice-v1-template-init`](docs/plan/planning/done/slice-v1-template-init.md).
- **Cross-platform binary distribution** for six platforms
  (Linux/macOS/Windows × amd64/arm64). `make build-binaries`
  cross-compiles every supported `GOOS`/`GOARCH` combination via
  the pinned `golang:$(GO_VERSION)` builder image (CGO disabled,
  `-ldflags "-s -w -X main.version=$(VERSION)"`, output to
  `bin/u-boot-<os>-<arch>[.exe]`). `.github/workflows/publish.yml`
  builds the same set after the GHCR push on every `v*` tag and
  attaches them as GitHub-Release assets via `gh release upload`.
  See
  [`slice-v2-binary-distribution`](docs/plan/planning/done/slice-v2-binary-distribution.md)
  — [ADR-0007](docs/plan/adr/0007-distributionswege-ghcr.md) §Folgepunkte 1 trigger pulled forward by the
  doctor-container-awareness feedback.
- Quickstart in `README.md` / `README.de.md` gets a host-native
  install block (`curl -sSL … | chmod +x` for Linux/macOS,
  `Invoke-WebRequest` for Windows) as the primary recommended path;
  the GHCR `docker run …` block is demoted to "alternative for
  container/CI workflows".
- `internal/hexagon/port/driven.RuntimeEnvironment` port plus
  `internal/adapter/driven/runtime.FileEnv` adapter: best-effort
  container detection via `/.dockerenv` (Docker Engine / Desktop)
  and `/run/.containerenv` (Podman / CRI-O / buildah). Drives the
  doctor-container-awareness change below.

### Changed

- `u-boot doctor` now skips the four host-prerequisite checks
  (`git.installed`, `docker.installed`, `docker.reachable`,
  `docker.compose.installed`) when running inside a container, with
  a `SeverityInfo` diagnostic and a hint that points at
  [`slice-v0.1.1-doctor-container-awareness`](docs/plan/planning/done/slice-v0.1.1-doctor-container-awareness.md)
  for the rationale. Effect: `docker run --rm
  ghcr.io/pt9912/u-boot:0.2.0 doctor` no longer mis-reports a
  healthy host as 4 errors; exit code on an otherwise-clean project
  goes from 11 to 0. This addresses real-world feedback from the
  first `v0.1.0` GHCR pull (2026-05-31) where the distroless image's
  lack of bundled `docker` / `git` binaries surfaced as false-
  positive errors.
- Host installations are unaffected — `runtime.FileEnv` returns
  `false` outside containers, so the existing
  [`LH-FA-DIAG-002`](spec/lastenheft.md#lh-fa-diag-002--lokale-voraussetzungen-prüfen)-classified errors / warnings remain.

### Notes

`releases/latest/download/u-boot-<os>-<arch>[.exe]` resolves to the
highest stable tag — since `v0.1.0` predates binary assets, the
`latest`-shortcut starts working with `v0.2.0` (or any later tag).

## [0.1.0] - 2026-05-31

First public release. Closes the MVP scope from
[`spec/lastenheft.md`](spec/lastenheft.md) MVP-priority IDs: all
`LH-AK-*`, `LH-FA-*` and `LH-SA-*` items are delivered (audit trail
in [`docs/plan/planning/in-progress/roadmap.md`](docs/plan/planning/in-progress/roadmap.md)
§MVP-Bilanz).

### Added — Subcommands

- `u-boot init [name] [--devcontainer]` — generate project skeleton
  (`u-boot.yaml`, `compose.yaml`, `README.md`, `CHANGELOG.md`,
  `.env.example`, `.gitignore`, `docker/`, `scripts/`, `docs/`)
  and run `git init` ([`LH-FA-INIT-001`](spec/lastenheft.md#lh-fa-init-001--neues-projekt-initialisieren)..[`LH-FA-INIT-007`](spec/lastenheft.md#lh-fa-init-007--git-repository-initialisierung), [`LH-AK-001`](spec/lastenheft.md#lh-ak-001--minimaler-init-flow), [`LH-AK-005`](spec/lastenheft.md#lh-ak-005--devcontainer-flow)).
  Mode flags `--yes` / `--no-interactive` / `--assume-existing`
  ([`LH-FA-CLI-005A`](spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung)); re-init with `--force` / `--backup` for
  managed-block edits vs. full overwrite ([`LH-FA-INIT-005`](spec/lastenheft.md#lh-fa-init-005--überschreibschutz)).
- `u-boot doctor [--strict]` — 11 diagnostic checks against the
  local environment + project, severity-classified
  (ok / warn / error), repair-hint output, exit-code 11 on errors
  (or warns under `--strict`) ([`LH-FA-DIAG-001`](spec/lastenheft.md#lh-fa-diag-001--doctor-befehl)..[`LH-FA-DIAG-004`](spec/lastenheft.md#lh-fa-diag-004--reparaturhinweise)).
- `u-boot add <service>` — idempotent state-machine for service
  add-ons; today's catalogue: `postgres` only ([`LH-FA-ADD-001`](spec/lastenheft.md#lh-fa-add-001--add-on-befehl)/[`LH-FA-ADD-002`](spec/lastenheft.md#lh-fa-add-002--postgresql-hinzufügen)/[`LH-FA-ADD-005`](spec/lastenheft.md#lh-fa-add-005--mehrfaches-hinzufügen-verhindern),
  [`LH-AK-002`](spec/lastenheft.md#lh-ak-002--postgresql-flow), [`LH-AK-006`](spec/lastenheft.md#lh-ak-006--idempotenz)). Keycloak ([`LH-AK-003`](spec/lastenheft.md#lh-ak-003--keycloak-flow)) and
  OpenTelemetry ([`LH-AK-004`](spec/lastenheft.md#lh-ak-004--opentelemetry-flow)) follow in V1.
- `u-boot up [--timeout <sec>]` and `u-boot down [--volumes]` —
  Compose wrapper with healthcheck polling and TCP port probes
  ([`LH-FA-UP-001`](spec/lastenheft.md#lh-fa-up-001--umgebung-starten)..[`LH-FA-UP-004`](spec/lastenheft.md#lh-fa-up-004--umgebung-stoppen)). `--timeout 0` is fire-and-forget.
- `u-boot generate <changelog|readme|env-example|devcontainer>` —
  idempotent block-replace via the `U-BOOT MANAGED BLOCK: init`
  marker; user content outside the managed region is preserved
  byte-identically. `changelog` carries the [`LH-AK-007`](spec/lastenheft.md#lh-ak-007--changelog-generator) pin
  (no destructive edits to existing entries). Exit codes
  `0` / `2` / `10` / `14` per [`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006--exit-codes) ([`LH-FA-GEN-001`](spec/lastenheft.md#lh-fa-gen-001--generate-befehl)..[`LH-FA-GEN-005`](spec/lastenheft.md#lh-fa-gen-005--idempotenz),
  [`LH-FA-DEV-001`](spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen)/[`LH-FA-DEV-004`](spec/lastenheft.md#lh-fa-dev-004--benutzerrechte)/[`LH-FA-DEV-005`](spec/lastenheft.md#lh-fa-dev-005--ports)).
- `u-boot config [get <path> | set <path> <value>]` — whitelist-
  scoped reads/writes with two-stage schema validation (struct
  unmarshal + per-path domain re-validation) before any
  `WriteFile`. `services.<svc>.enabled` is get-only; toggling
  happens through `add` / `remove` to keep the add-on state
  machine atomic ([`LH-FA-CONF-001`](spec/lastenheft.md#lh-fa-conf-001--projektkonfiguration)..[`LH-FA-CONF-005`](spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern)).

### Added — CI & release infrastructure

- GitHub Actions CI workflow `.github/workflows/ci.yml` with three
  PR-blocking jobs ([`LH-QA-003`](spec/lastenheft.md#lh-qa-003--ci-fähigkeit-github-actions)): `gates (lint + test +
  coverage-gate)`, `security-gates (govulncheck)`,
  `image-scan (trivy HIGH+CRITICAL)`. All actions SHA-pinned;
  Docker-only runner ([`LH-FA-BUILD-007`](spec/lastenheft.md#lh-fa-build-007--docker-only-workflow)); per-job minimal
  permissions.
- GitHub Actions release workflow
  `.github/workflows/publish.yml` triggered on `v*` tags. Strict
  SemVer-2.0 validation (rejects leading-zero numeric prereleases
  and build-metadata `+...` tags), GHCR image push to
  `ghcr.io/pt9912/u-boot:<version>` (plus `:latest` for stable
  tags), OCI label verification, and live `--version` smoke test
  against the tag-derived `VERSION`.
- Local outer/inner-loop parity: `make image-scan` reproduces the
  `image-scan` CI job using the same Trivy version
  (`TRIVY_VERSION ?= 0.70.0`) the action installs.
- Multi-stage distroless runtime image (`gcr.io/distroless/static-debian12:nonroot`)
  built via `make build`; CGO-disabled static binary; version
  injected at build time as `-X main.version=<UBOOT_VERSION>`
  and as the `org.opencontainers.image.version` label.

### Added — Architecture & documentation

- Hexagonal architecture ([`LH-FA-ARCH-001`](spec/lastenheft.md#lh-fa-arch-001--hexagonales-pattern)..[`LH-FA-ARCH-003`](spec/lastenheft.md#lh-fa-arch-003--import-regeln-und-enforcement), [ADR-0002](docs/plan/adr/0002-hexagonale-architektur.md)):
  `internal/hexagon/{domain,application,port/{driving,driven}}`
  + `internal/adapter/{driving,driven}`. `depguard` enforces
  layer rules in CI.
- 10 ADRs cover language (Go), architecture (hexagonal), lint
  profile (SOLID-near), CI system, CLI framework (Cobra), revive
  custom rules, distribution path (GHCR), plugin system (static —
  no plugins), template format (YAML + Go `text/template`), and
  the HTTP adapter (not built; CLI-only).
- User-facing setup docs:
  [`docs/user/quality.md`](docs/maintainer/quality.md) (quality-gates
  overview) and
  [`docs/user/branch-protection.md`](docs/maintainer/branch-protection.md)
  (one-time GitHub UI activation of required status checks).
- German `spec/lastenheft.md` (~3000 lines, 14 sections + 4 open
  points all decided) is the single source of truth; English
  `README.md` / German `README.de.md` are equivalent.

### Known limitations and deliberate carve-outs

- **Add-on catalogue is intentionally small:** only `postgres`
  ships in v0.1.0. Keycloak and OpenTelemetry are V1.
- **Templates implementation is V1.** Format is decided
  ([ADR-0009](docs/plan/adr/0009-template-format-yaml-files.md): YAML + `text/template`); the three implementation
  slices ([`slice-v1-template-list`](docs/plan/planning/done/slice-v1-template-list.md), [`slice-v1-template-init`](docs/plan/planning/done/slice-v1-template-init.md),
  [`slice-later-local-templates`](docs/plan/planning/done/slice-later-local-templates.md)) follow on demand.
- **JSON / machine-readable output is V1.** `--json` and
  `--dry-run` flags ([`LH-FA-CLI-007`](spec/lastenheft.md#lh-fa-cli-007--dry-run)/[`LH-FA-CLI-008`](spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe), [`LH-NFA-USE-004`](spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe)) are
  not yet shipped; [ADR-0010](docs/plan/adr/0010-kein-http-driving-adapter.md) (no HTTP adapter) explicitly relies
  on this V1 track landing.
- **Distribution is GHCR-only.** Binary, Homebrew, Debian/RPM
  paths are deferred with explicit trigger slices in [ADR-0007](docs/plan/adr/0007-distributionswege-ghcr.md).
  `npm` / `pip` are rejected (ecosystem mismatch).
- **No plugin loader.** Add-on system stays statically compiled
  into u-boot ([ADR-0008](docs/plan/adr/0008-plugin-system-statisch.md)). Four re-evaluation triggers documented
  in [ADR-0008](docs/plan/adr/0008-plugin-system-statisch.md) §Folgepunkte.
- **CLI-only.** No HTTP / daemon adapter ([ADR-0010](docs/plan/adr/0010-kein-http-driving-adapter.md)); programmatic
  consumers use subprocess + `--json` once V1 lands.
- **Inner-loop is Docker-only** ([`LH-FA-BUILD-007`](spec/lastenheft.md#lh-fa-build-007--docker-only-workflow)). GNU `make`
  remains the single non-Docker host dependency (permanent
  carve-out to [`LH-NFA-PORT-002`](spec/lastenheft.md#lh-nfa-port-002--keine-unnötigen-systemabhängigkeiten)).

### Setup — required one-time GitHub UI activation

Before merging external PRs against `main`, activate the three
required status checks in GitHub UI per
[`docs/user/branch-protection.md`](docs/maintainer/branch-protection.md):
the exact match strings are the workflow `name:` fields
(`gates (lint + test + coverage-gate)`,
`security-gates (govulncheck)`,
`image-scan (trivy HIGH+CRITICAL)`), not the shorter `jobs.<key>`
identifiers.

[Unreleased]: https://github.com/pt9912/u-boot/compare/v0.15.1...HEAD
[0.15.1]: https://github.com/pt9912/u-boot/compare/v0.15.0...v0.15.1
[0.15.0]: https://github.com/pt9912/u-boot/compare/v0.14.0...v0.15.0
[0.14.0]: https://github.com/pt9912/u-boot/compare/v0.13.1...v0.14.0
[0.13.1]: https://github.com/pt9912/u-boot/compare/v0.13.0...v0.13.1
[0.13.0]: https://github.com/pt9912/u-boot/compare/v0.12.0...v0.13.0
[0.12.0]: https://github.com/pt9912/u-boot/compare/v0.11.0...v0.12.0
[0.11.0]: https://github.com/pt9912/u-boot/compare/v0.10.0...v0.11.0
[0.10.0]: https://github.com/pt9912/u-boot/compare/v0.9.0...v0.10.0
[0.9.0]: https://github.com/pt9912/u-boot/compare/v0.8.1...v0.9.0
[0.8.1]: https://github.com/pt9912/u-boot/compare/v0.8.0...v0.8.1
[0.8.0]: https://github.com/pt9912/u-boot/compare/v0.7.0...v0.8.0
[0.7.0]: https://github.com/pt9912/u-boot/compare/v0.6.0...v0.7.0
[0.6.0]: https://github.com/pt9912/u-boot/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/pt9912/u-boot/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/pt9912/u-boot/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/pt9912/u-boot/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/pt9912/u-boot/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/pt9912/u-boot/releases/tag/v0.1.0
