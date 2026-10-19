#!/bin/bash
# D3, Zusatz: benutzerdefiniertes Docker-Netz mit IPv6 (Netz u-boot-messung-dns-net, Unique-Local-Praefix). Sondiert IPv6 gegen die Tabelle ip6
# des Startscripts und den Docker-eigenen Resolver 127.0.0.11, der in einem solchen Netz der einzige "Upstream" des Containers ist.
# Aufruf: WORK=<messordner> OUT=<ausgabeverzeichnis> 07-d3-netz.sh -> $OUT/rohdaten-07-d3-netz.txt
# Voraussetzung: 02-projekte.sh. Faelle: ohne (Kontrolle), ist, skuid (Messvariante).
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-07-d3-netz.txt
net=$PFX-net
r=$PFX-resp6
c=$PFX-d3netz
docker network rm $net > /dev/null 2>&1
docker network create --ipv6 --subnet fd00:6d73:6e00:7::/64 $net > /dev/null; echo "netz angelegt: exit=$?"
docker rm -f -v $r > /dev/null 2>&1
docker run -d --init --name $r --network $net -u root --entrypoint sleep $IMG_TOOL infinity > /dev/null
docker cp "$here/inner-resp.sh" $r:/tmp/inner-resp.sh
docker exec -u root $r sh /tmp/inner-resp.sh
R6=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.GlobalIPv6Address}}{{end}}' $r)
echo "Antwort-Container (IPv6) R6=$R6"
for mode in ohne ist skuid; do
  echo; echo "######## D3 Netz, Fall $mode, $(now)"
  start $c $IMG_TOOL --network $net
  case "$mode" in
    ohne) echo "(kein Startscript: Kontrolle)" ;;
    *)    egress $c $mode ;;
  esac
  docker exec -u vscode $c bash /m/inner-d3-netz.sh "$mode" "$R6"
  docker rm -f -v $c > /dev/null
done
docker rm -f -v $r > /dev/null
docker network rm $net > /dev/null
