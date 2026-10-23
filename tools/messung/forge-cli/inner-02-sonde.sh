#!/bin/bash
# Laeuft IM Werkzeug-Container; aufgerufen von 02-oci-quellen.sh. Sondiert Namen unter bekannten Praefixen in ghcr.io.
set -u
for ref in devcontainers/features/gitlab-cli devcontainers/features/glab devcontainers-extra/features/gitlab-cli devcontainers-extra/features/glab devcontainers-contrib/features/gh-cli devcontainers-contrib/features/gitlab-cli; do
  t=$(curl -s "https://ghcr.io/token?scope=repository:$ref:pull" | jq -r .token)
  printf '%s http=' "$ref"
  curl -s -o /dev/null -w '%{http_code}\n' -H "Authorization: Bearer $t" "https://ghcr.io/v2/$ref/tags/list"
done
