#!/bin/bash
# D5: Datenkanal ueber erlaubte Namen: werden Anfragen an Namen unter einem erlaubten Namen an den Upstream weitergeleitet, und mit welchem
# Durchsatz? Der Antwort-Container (inner-resp.sh) ist der Upstream des Sandbox-Containers (docker run --dns R) und zaehlt, was bei ihm ankommt.
# Aufruf: WORK=<messordner> OUT=<ausgabeverzeichnis> 08-d5.sh -> $OUT/rohdaten-08-d5.txt
# Voraussetzung: 02-projekte.sh. Faelle: ist, skuid (Messvariante). Der Durchsatz ist ein Messwert (zeitabhaengig, nicht Teil von vergleich.sh).
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-08-d5.txt
r=$PFX-resp
c=$PFX-d5
docker rm -f -v $r > /dev/null 2>&1
docker run -d --init --name $r -u root --entrypoint sleep $IMG_TOOL infinity > /dev/null
docker cp "$here/inner-resp.sh" $r:/tmp/inner-resp.sh
docker exec -u root $r sh /tmp/inner-resp.sh
R=$(docker inspect -f '{{.NetworkSettings.Networks.bridge.IPAddress}}' $r)
echo "Antwort-Container (= Upstream des Sandbox-Containers) R=$R"
for mode in ist skuid; do
  echo; echo "######## D5, Fall $mode, $(now)"
  docker exec -u root $r sh -c ': > /tmp/dns53.log'
  start $c $IMG_TOOL --dns "$R"
  egress $c $mode
  docker exec -u vscode $c bash /m/inner-d5.sh $mode "$R"
  echo "-- am Antwort-Container angekommen (Log dns53.log)"
  echo "Anfragen gesamt: $(docker exec $r grep -c 'query\[A\]' /tmp/dns53.log)"
  echo "davon unter api.github.com: $(docker exec $r grep -c 'query\[A\] .*\.api\.github\.com from' /tmp/dns53.log)"
  echo "davon unter example.com: $(docker exec $r grep -c 'query\[A\] .*\.example\.com from' /tmp/dns53.log)"
  echo "Absender (Adressen, Anzahl): $(docker exec $r sh -c "grep 'query\[A\]' /tmp/dns53.log | sed 's/.* from //' | sort | uniq -c | tr -s ' ' | tr '\n' ';'")"
  echo "Beispielzeile (gekuerzt): $(docker exec $r sh -c "grep 'query\[A\] kanal1.api' /tmp/dns53.log | head -1")"
  aussage "$mode d5 am-upstream-angekommen-kanal1-api.github.com=$([ "$(docker exec $r grep -c 'query\[A\] kanal1.api.github.com from' /tmp/dns53.log)" -gt 0 ] && echo ja || echo nein)"
  aussage "$mode d5 am-upstream-angekommen-kanal1-example.com=$([ "$(docker exec $r grep -c 'query\[A\] kanal1.example.com from' /tmp/dns53.log)" -gt 0 ] && echo ja || echo nein)"
  docker rm -f -v $c > /dev/null
done
docker rm -f -v $r > /dev/null
