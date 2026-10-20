# Helper image of `u-boot devcontainer`: the pinned Dev Container CLI with a
# Docker client, run as a short-lived container that talks to the host
# engine socket. The Go adapter embeds this file with package.json and
# package-lock.json as the build context; the image tag derives from the
# CLI version and the hash of the three files.
#
# Pins: both base images by index digest (linux/amd64 and linux/arm64), the
# CLI by exact version in package.json and by integrity hash in
# package-lock.json. Installation is `npm ci --ignore-scripts`: no
# lifecycle script of a package runs.
FROM docker:29-cli@sha256:6fae2214c3121a28d5aa639d09032f61310b7c9c75ffbad1025f261a1b621800 AS dockercli
FROM node:22-bookworm@sha256:0e5f906573693feaa1e21057ebdcfdb5bd5021f050b2dc7c9deceb629c7da2a8
COPY --from=dockercli /usr/local/bin/docker /usr/local/bin/docker
COPY --from=dockercli /usr/local/libexec/docker/cli-plugins/docker-buildx /usr/local/libexec/docker/cli-plugins/docker-buildx
WORKDIR /opt/devcontainer-cli
COPY package.json package-lock.json ./
RUN npm ci --ignore-scripts
ENV PATH="/opt/devcontainer-cli/node_modules/.bin:${PATH}"
