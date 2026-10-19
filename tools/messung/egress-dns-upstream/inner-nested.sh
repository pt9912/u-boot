#!/bin/bash
# Laeuft im Container (Image pod, runArgs wie devcontainer.json) als vscode; $1 = Modus (ist|skuid). D2, verschachtelter Podman. Siehe 05-d2-nested.sh.
# Der Ablauf folgt nestedPodmanEgressBuildScript aus internal/e2e/sandbox_devcontainer_docker_test.go (egress-init, sandbox-init, Dockerfile mit
# drei Schritten, `podman build`); dazu Zaehler der Messvariante und Einzelproben mit `podman run`.
set -u
mode=${1:?ist|skuid}
echo "id: $(id)"
echo "-- Startscript des Egress ($mode, mit Zaehlern)"
sudo sh /m/variante.sh "$mode" zaehler; echo "egress-init exit=$?"
echo "-- weitere Zaehler (Messinstrument, nicht terminierend): Benutzer-ID 1000 (vscode), allein und mit Ziel 8.8.8.8 oder 8.8.4.4, an Port 53"
for sel in "meta skuid 1000" "meta skuid 1000 ip daddr 8.8.8.8" "meta skuid 1000 ip daddr 8.8.4.4"; do
  # shellcheck disable=SC2086
  sudo nft insert rule ip u_boot_egress output oifname != lo meta l4proto '{ tcp, udp }' th dport 53 $sel counter
done
echo "resolv.conf des Sandbox-Containers: $(grep -v '^#' /etc/resolv.conf | tr '\n' ' ')"
echo "-- sandbox-init"
sh /usr/local/bin/u-boot-sandbox-init; echo "sandbox-init exit=$?"
mkdir -p "$HOME/ctx" && cd "$HOME/ctx"
cat > Dockerfile <<'DF'
FROM docker.io/library/debian:stable-slim
RUN cat /etc/resolv.conf
RUN getent hosts deb.debian.org > /dev/null && echo ALLOWED_RESOLVES
RUN ! getent hosts example.com && echo UNLISTED_BLOCKED
DF
echo "-- podman build (Dockerfile: Schritt 1 resolv.conf anzeigen, Schritt 2 getent deb.debian.org, Schritt 3 ! getent example.com)"
s=$(date +%s)
podman build -t probe-egress . 2>&1 | sed -E 's/^([0-9a-f]{12,})$/<id>/' | tail -40
rc=${PIPESTATUS[0]}
echo "podman build exit=$rc dauer=$(( $(date +%s) - s ))s"
echo "AUSSAGE $mode nested-build-exit=$rc"
echo "-- Zaehler nach dem Bau (Versuche anderer Benutzer als nobody: Pakete an Port 53, nicht ueber lo; 'drop' = bis zum Ende nicht angenommen)"
sudo -n nft list chain ip u_boot_egress output | grep -E 'counter' | sed 's/^[[:space:]]*/    /'

if true; then
  p() { # <beschriftung> <podman-run-argumente...> : getent example.com / deb.debian.org / api.github.com im verschachtelten Container
    local label=$1; shift
    echo "-- podman run $label"
    timeout 60 podman run --rm "$@" docker.io/library/debian:stable-slim sh -c '
      echo "resolv.conf: $(grep -v "^#" /etc/resolv.conf | tr "\n" " ")"
      for h in example.com deb.debian.org api.github.com; do
        s=$(date +%s%N); a=$(getent hosts "$h" | head -1 | cut -d" " -f1); e=$(date +%s%N)
        echo "$h: ${a:-keine Adresse} ($(( (e - s) / 1000000 )) ms)"
      done' > /tmp/p.out 2>&1
    tail -6 /tmp/p.out
    while read -r h a rest; do
      case "$h" in example.com:|deb.debian.org:|api.github.com:) echo "AUSSAGE $mode nested-run [$label] ${h%:}=$([ "$a" = keine ] && echo keine || echo adresse)" ;; esac
    done < /tmp/p.out
  }
  p "(Vorgabe: pasta, Standard-DNS)"
  p "--dns 8.8.8.8 (Upstream des Docker-Hosts)" --dns 8.8.8.8
  p "--dns 9.9.9.9 (kein Upstream des Docker-Hosts)" --dns 9.9.9.9
  p "--network=host" --network=host
  echo "-- Zaehler nach den Einzelproben"
  sudo -n nft list chain ip u_boot_egress output | grep -E 'counter' | sed 's/^[[:space:]]*/    /'
  dropped=$(sudo -n nft list chain ip u_boot_egress output | sed -n 's/.*counter packets \([0-9]*\) bytes [0-9]* drop/\1/p')
  echo "AUSSAGE $mode nested-dns-verworfen=$([ "${dropped:-0}" -gt 0 ] && echo ja || echo nein)"
fi
