#!/bin/bash
# runs inside a debian container as root; $1 = hosts for egress, rest = names to probe
set -u
export DEBIAN_FRONTEND=noninteractive
apt-get update -qq >/dev/null && apt-get install -y -qq --no-install-recommends dnsmasq-base nftables curl dnsutils ca-certificates >/dev/null 2>&1
dnsmasq --version | head -1
HOSTS="$1"
shift
sed -e 's|{{.OnUnavailable}}|fail|g' \
    -e "s|^hosts=.*|hosts=\"$HOSTS\"|" /tmpl/egress-init.sh.tmpl > /tmp/egress-init.sh
mkdir -p /var/lib/u-boot/status
sh /tmp/egress-init.sh; echo "init exit=$?"
echo "--- dnsmasq conf"; grep -E '^(server|nftset)' /etc/u-boot-egress-dnsmasq.conf
for n in "$@"; do
  full=$(dig +time=3 +tries=1 @127.0.0.1 "$n" A 2>&1)
  st=$(echo "$full" | grep -o 'status: [A-Z]*')
  a=$(echo "$full" | awk '/IN[ \t]+(A|CNAME)[ \t]/ && !/^;/ {printf "%s>%s ", $4, $5}')
  code=$(curl -sS -m 6 -o /dev/null -w '%{http_code}' "https://$n" 2>&1 | tail -c 100)
  echo "RESULT $n | dns[$st] $a| curl[$code]"
done
