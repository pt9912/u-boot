#!/bin/bash
# M3: Zustellung vom Host in den Container (docker cp, docker exec) und Verhalten bei gestopptem Container.
# Aufruf: 04-m3-zustellung.sh > rohdaten-04-m3.txt 2>&1
# Voraussetzung: Docker (Colima); Image u-boot-messung-egress-neuladen-dig (README.md, Schritt 0).
# Der Container laeuft wie der Devcontainer (sleep infinity, --cap-add=NET_ADMIN); postStartCommand per `docker exec` als vscode.
set -u
IMG=u-boot-messung-egress-neuladen-dig
c=egress-m3
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"; docker rm -f $c >/dev/null 2>&1' EXIT
docker rm -f $c >/dev/null 2>&1
docker run -d --name $c --cap-add=NET_ADMIN $IMG sleep infinity >/dev/null
docker exec -u vscode $c sh -c 'sudo sh /usr/local/bin/u-boot-egress-init'; echo "postStart exit=$?"
own() { docker exec -u vscode $c ls -l "$@" | awk '{print "    " $1, $3, $4, $9}'; }
dns() { docker exec -u vscode $c dig +time=3 +tries=1 @127.0.0.1 "$1" A 2>&1 | grep -o 'status: [A-Z]*' | head -1; }

echo "######## 1: docker cp einer Liste in ein Root-eigenes Verzeichnis (/etc)"
printf 'api.github.com\nexample.org\n' > "$tmp/egress.allow"
ls -l "$tmp/egress.allow" | awk '{print "    Host-Datei: " $1, $3}'
docker cp "$tmp/egress.allow" $c:/etc/u-boot-egress.allow; echo "docker cp rc=$?"
own /etc/u-boot-egress.allow
if docker exec -u vscode $c sh -c ': >> /etc/u-boot-egress.allow' 2>/dev/null; then echo "    vscode schreibt in /etc/u-boot-egress.allow: ja"; else echo "    vscode schreibt in /etc/u-boot-egress.allow: nein"; fi
if docker exec -u vscode $c sh -c 'rm -f /etc/u-boot-egress.allow' 2>/dev/null; then echo "    vscode loescht /etc/u-boot-egress.allow: ja"; else echo "    vscode loescht /etc/u-boot-egress.allow: nein"; fi
docker exec -u vscode $c test -e /etc/u-boot-egress.allow && echo "    Datei noch da" || echo "    Datei weg"

echo "######## 2: docker cp einer Liste in ein Verzeichnis des Benutzers (/var/lib/u-boot/status, gehoert vscode)"
docker cp "$tmp/egress.allow" $c:/var/lib/u-boot/status/egress.allow; echo "docker cp rc=$?"
own /var/lib/u-boot/status/egress.allow
if docker exec -u vscode $c sh -c ': >> /var/lib/u-boot/status/egress.allow' 2>/dev/null; then echo "    vscode schreibt die Datei: ja"; else echo "    vscode schreibt die Datei: nein"; fi
docker exec -u vscode $c sh -c 'rm -f /var/lib/u-boot/status/egress.allow && printf "evil.example\n" > /var/lib/u-boot/status/egress.allow'; echo "    vscode loescht und legt neu an: rc=$?"
own /var/lib/u-boot/status/egress.allow

echo "######## 3: Reload per docker exec -u root mit dem gebackenen Script (Liste unveraendert)"
echo "example.org vorher: $(dns example.org)"
docker exec -u root $c sh /usr/local/bin/u-boot-egress-init; echo "docker exec -u root: exit=$?"
echo "api.github.com nachher: $(dns api.github.com)"

echo "######## 4: Host ersetzt das gebackene Script durch eine Kopie mit anderer Liste, danach der postStartCommand-Aufruf als vscode"
docker exec -u root $c sh -c 'sed -e "s|^hosts=.*|hosts=\"api.github.com example.org\"|" /usr/local/bin/u-boot-egress-init > /tmp/neu.sh'
docker cp $c:/tmp/neu.sh "$tmp/neu.sh"
docker cp "$tmp/neu.sh" $c:/usr/local/bin/u-boot-egress-init; echo "docker cp ueber das gebackene Script rc=$?"
own /usr/local/bin/u-boot-egress-init
echo "example.org vor dem Aufruf: $(dns example.org)"
docker exec -u vscode $c sh -c 'sudo sh /usr/local/bin/u-boot-egress-init'; echo "postStart exit=$?"
echo "example.org nach dem Aufruf: $(dns example.org)"

echo "######## 5: gestoppter Container"
docker stop -t 2 $c >/dev/null
docker exec -u root $c true 2>&1 | head -2 | sed 's/^/    exec: /'
docker cp "$tmp/egress.allow" $c:/etc/u-boot-egress.allow; echo "docker cp in den gestoppten Container: rc=$?"
docker cp $c:/etc/u-boot-egress.allow "$tmp/zurueck.allow"; echo "docker cp aus dem gestoppten Container: rc=$?"
docker inspect -f '    Running={{.State.Running}}' $c

echo "######## 6: nach docker start (ohne postStartCommand): Zustand der Beschraenkung"
docker start $c >/dev/null; sleep 1
docker exec -u root $c sh -c 'echo "    nft-Tabellen: [$(nft list tables | tr "\n" " ")]"; echo "    resolv.conf: $(grep nameserver /etc/resolv.conf | tr "\n" " ")"; echo "    dnsmasq: [$(ps -u nobody -o args= | grep ^dnsmasq)]"'
docker exec -u root $c sh -c 'curl -sS -m 5 -o /dev/null -w "    example.com ohne Beschraenkung: http=%{http_code}\n" https://example.com 2>&1 | tail -1'
docker exec -u vscode $c sh -c 'sudo sh /usr/local/bin/u-boot-egress-init'; echo "postStart exit=$?"
docker exec -u root $c sh -c 'curl -sS -m 5 -o /dev/null -w "    example.com nach postStart: http=%{http_code}\n" https://example.com 2>&1 | tail -1'

echo "######## 7: Eigentuemer der kopierten Datei: docker cp -a und Tar-Strom mit festem Eigentuemer (Tar im Wegwerf-Container erzeugt)"
docker exec -u root $c sh -c 'rm -f /etc/u-boot-egress.allow /etc/u-boot-egress-a.allow /etc/u-boot-egress-t0.allow /etc/u-boot-egress-t1000.allow'
docker cp -a "$tmp/egress.allow" $c:/etc/u-boot-egress-a.allow; echo "docker cp -a rc=$?"
own /etc/u-boot-egress-a.allow
for uid in 0 1000; do
  docker run --rm debian:stable-slim sh -c "mkdir /t && printf 'api.github.com\n' > /t/u-boot-egress-t$uid.allow && tar --owner=$uid --group=$uid -C /t -cf - u-boot-egress-t$uid.allow" | docker cp - $c:/etc/; echo "Tar-Strom mit owner=$uid: docker cp rc=$?"
  own /etc/u-boot-egress-t$uid.allow
  if docker exec -u vscode $c sh -c ": >> /etc/u-boot-egress-t$uid.allow" 2>/dev/null; then echo "    vscode schreibt die Datei (owner=$uid): ja"; else echo "    vscode schreibt die Datei (owner=$uid): nein"; fi
done
