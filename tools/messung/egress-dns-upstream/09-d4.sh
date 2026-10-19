#!/bin/bash
# D4: Bleibt das Startscript beim zweiten Lauf, nach einem Neustart des Containers, beim Neuaufbau (neuer Container aus demselben Image) und in
# 50 erneuten Laeufen hintereinander lauffaehig? Faelle ist (erzeugtes Startscript) und skuid (Messvariante); der Fall ist ist die Gegenprobe.
# Aufruf: WORK=<messordner> OUT=<ausgabeverzeichnis> 09-d4.sh -> $OUT/rohdaten-09-d4.txt
# Voraussetzung: 02-projekte.sh. Die 50 Laeufe folgen der Messung B des Vorgaengers (slice-v1-sandbox-egress-neuladen-entwurf), ohne Last.
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-09-d4.txt
c=$PFX-d4
for mode in ist skuid; do
  echo; echo "######## D4, Fall $mode, $(now)"
  start $c $IMG_TOOL
  echo "-- erster Lauf (postStartCommand)"
  egress $c $mode; docker exec -u root $c bash /m/inner-d4.sh $mode erster-lauf
  echo "-- zweiter Lauf im selben Container"
  egress $c $mode; docker exec -u root $c bash /m/inner-d4.sh $mode zweiter-lauf
  echo "-- docker restart, danach der Lauf des postStartCommand"
  docker restart -t 2 $c > /dev/null
  echo "nach dem Neustart, vor dem Lauf: resolv.conf=$(docker exec $c grep -c '^nameserver 127.0.0.1$' /etc/resolv.conf) (0 = von Docker neu geschrieben), dnsmasq-prozesse=$(docker exec $c pgrep -x dnsmasq | wc -l), Original-Datei vorhanden=$(docker exec $c sh -c '[ -f /etc/resolv.conf.u-boot-orig ] && echo ja || echo nein')"
  egress $c $mode; docker exec -u root $c bash /m/inner-d4.sh $mode nach-neustart
  echo "-- Neuaufbau: neuer Container aus demselben Image"
  docker rm -f -v $c > /dev/null
  start $c $IMG_TOOL
  egress $c $mode; docker exec -u root $c bash /m/inner-d4.sh $mode neuaufbau
  echo "-- 50 erneute Laeufe hintereinander (Rennen um den Port des dnsmasq); je Lauf Zustand pruefen, Auffaelliges wird ausgegeben"
  bad=0; msg=0
  for i in $(seq 1 50); do
    out=$(docker exec -u root $c sh -c "if [ $mode = skuid ]; then sh /m/variante.sh skuid; else sh /usr/local/bin/u-boot-egress-init; fi" 2>&1); rc=$?
    if [ $rc != 0 ] || printf '%s' "$out" | grep -q -i -E 'address already in use|failed to create|error'; then msg=$((msg+1)); echo "lauf $i: exit=$rc Meldung: $(printf '%s' "$out" | grep -i -E 'address already in use|failed to create|error' | head -1)"; fi
    z=$(docker exec -u root $c bash /m/inner-d4.sh $mode "lauf-$i" | head -1)
    if [ $mode = skuid ]; then erw="direkt-example.com=keine"; erws="skuid-regeln=4"; else erw="direkt-example.com=antwort"; erws="skuid-regeln=0"; fi
    case "$z" in *"dnsmasq-prozesse=1 uid=65534, pid-lebt=ja"*"$erws"*"lokal-api.github.com='status: NOERROR'"*"$erw"*) ;; *) bad=$((bad+1)); echo "lauf $i: AUFFAELLIG $z" ;; esac
    if [ "$i" = 1 ] || [ "$i" = 50 ]; then echo "$z"; fi
  done
  echo "Zusammenfassung: 50 Laeufe, $bad mit auffaelligem Zustand danach, $msg mit Fehlermeldung oder Exit ungleich 0"
  aussage "$mode d4 50-laeufe-auffaellig=$bad meldungen=$msg"
  echo "-- Regeln nach dem 50. Lauf"
  docker exec -u root $c nft list chain ip u_boot_egress output | sed 's/^[[:space:]]*/    /'
  docker rm -f -v $c > /dev/null
done
