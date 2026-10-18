#!/bin/bash
# Frage c (CNAME-Kette) und d (nur Subdomains). Aufruf: inner-cd.sh c | d
set -u
export DEBIAN_FRONTEND=noninteractive
apt-get update -qq >/dev/null && apt-get install -y -qq --no-install-recommends dnsmasq-base nftables curl dnsutils ca-certificates >/dev/null 2>&1
dnsmasq --version | head -1
case "$1" in
c)
  echo "--- CNAME-Kandidaten (unrestringiert)"
  for n in pypi.org files.pythonhosted.org www.python.org registry.npmjs.org ghcr.io docs.github.com www.mozilla.org www.debian.org deb.debian.org dl.google.com cdn.jsdelivr.net www.docker.com registry-1.docker.io production.cloudflare.docker.com; do
    echo "$n: $(dig +short "$n" CNAME | tr '\n' ' ')"
  done
  ;;
d)
  for pat in '*.example.org' '.example.org'; do
    echo "=== pattern $pat"
    pkill dnsmasq 2>/dev/null; sleep 1
    ns=$(awk '/^nameserver/ && $2 !~ /:/ {print $2; exit}' /etc/resolv.conf)
    cat > /tmp/d.conf <<EOF
listen-address=127.0.0.1
bind-interfaces
no-resolv
no-hosts
server=/$pat/$ns
EOF
    if dnsmasq --test --conf-file=/tmp/d.conf 2>&1; then
      dnsmasq --conf-file=/tmp/d.conf --pid-file=/tmp/d.pid
      sleep 1
      for n in example.org www.example.org a.b.example.org example.com; do
        echo "$n: $(dig +time=3 +tries=1 @127.0.0.1 "$n" A 2>&1 | grep -oE 'status: [A-Z]+') $(dig +short +time=3 +tries=1 @127.0.0.1 "$n" A 2>&1 | head -1)"
      done
    fi
  done
  ;;
esac
