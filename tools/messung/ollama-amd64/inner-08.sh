#!/bin/bash
# Laeuft IM Container (Image oll) als vscode; aufgerufen von 08-wiederholungslauf.sh.
# Aufruf im Container: inner-08.sh <laeufe>
# Misst, ob und wie lange waehrend eines Wiederholungslaufs des erzeugten Startscripts die Restriktion fehlt (Fenster) und ob erlaubte
# Aufrufe ausfallen. Sonden (Zeilen "<t0-ns> <t1-ns> <name> <exit>" in /tmp/p-*.log):
#   B1..B8  TCP-Verbindung zum Gateway des Containers, Port 9, mit 50 ms Frist (bash /dev/tcp): unter der Restriktion verworfen
#           (Exit 124), ohne Restriktion abgelehnt (Exit 1, sofort). Ein Exit ungleich 124 ist ein Treffer im Fenster.
#   N       curl auf einen gesperrten NAMEN (example.com): Exit 6 (Name nicht aufgeloest) oder 28 (verworfen); 0 waere ein Treffer.
#   A       curl auf einen erlaubten Namen (github.com): Exit 0 erwartet; ungleich 0 ist ein Ausfall.
#   D       getent hosts github.com (Aufloesung ueber den Resolver des Skripts): ungleich 0 ist ein Ausfall der Aufloesung.
set -u
runs=${1:-30}
gw=$(ip route | awk '/^default/ {print $3; exit}')
ns() { date +%s%N; }

echo "-- Gateway vorhanden: $([ -n "$gw" ] && echo ja || echo nein)"
echo "-- Basislinie ohne Restriktion (vor dem ersten Start des Skripts): B-Sonde"
timeout 0.5 bash -c "echo > /dev/tcp/$gw/9" 2> /dev/null; echo "B ohne Restriktion exit=$? (1 = abgelehnt, 124 = Frist)"
echo "AUSSAGE basislinie-ohne-restriktion-b-exit=$(timeout 0.5 bash -c "echo > /dev/tcp/$gw/9" 2> /dev/null; echo $?)"
echo "-- erster Start des Skripts"
s=$(ns); sudo sh /usr/local/bin/u-boot-egress-init; rc=$?; e=$(ns)
echo "egress-init exit=$rc dauer=$(( (e - s) / 1000000 )) ms"
echo "-- Basislinie mit Restriktion: B-Sonde (Exit 124 erwartet), N-Sonde, A-Sonde, D-Sonde"
timeout 0.05 bash -c "echo > /dev/tcp/$gw/9" 2> /dev/null; echo "B mit Restriktion exit=$?"
curl -sS -m 2 -o /dev/null http://example.com 2> /dev/null; echo "N exit=$?"
curl -sS -m 5 -o /dev/null https://github.com 2> /dev/null; echo "A exit=$?"
getent hosts github.com > /dev/null; echo "D exit=$?"
echo "AUSSAGE basislinie-mit-restriktion-b-exit=$(timeout 0.05 bash -c "echo > /dev/tcp/$gw/9" 2> /dev/null; echo $?)"

rm -f /tmp/p-*.log /tmp/marker.log
pids=()
for k in 1 2 3 4 5 6 7 8; do
  ( sleep "0.00$k"; while :; do t0=$(ns); timeout 0.05 bash -c "echo > /dev/tcp/$gw/9" 2> /dev/null; r=$?; echo "$t0 $(ns) B$k $r"; done ) >> /tmp/p-B.log &
  pids+=($!)
done
( while :; do t0=$(ns); curl -sS -m 1 -o /dev/null http://example.com 2> /dev/null; r=$?; echo "$t0 $(ns) N $r"; done ) >> /tmp/p-N.log &
pids+=($!)
( while :; do t0=$(ns); curl -sS -m 3 -o /dev/null https://github.com 2> /dev/null; r=$?; echo "$t0 $(ns) A $r"; sleep 0.05; done ) >> /tmp/p-A.log &
pids+=($!)
( while :; do t0=$(ns); getent hosts github.com > /dev/null 2>&1; r=$?; echo "$t0 $(ns) D $r"; sleep 0.01; done ) >> /tmp/p-D.log &
pids+=($!)
sleep 2

for k in $(seq 1 "$runs"); do
  s=$(ns); sudo sh /usr/local/bin/u-boot-egress-init; rc=$?; e=$(ns)
  echo "$k $s $e $rc" >> /tmp/marker.log
  sleep 2
done
kill "${pids[@]}" 2> /dev/null; sleep 1
kill "${pids[@]}" 2> /dev/null

cat /tmp/p-*.log | sort -n > /tmp/alle.log
echo
echo "== Auswertung ($runs Laeufe des Skripts; Zeitfenster je Lauf: von Start des Skripts bis 1,5 s nach seinem Ende)"
echo "Sonden insgesamt: B=$(grep -c ' B[0-9] ' /tmp/alle.log) N=$(grep -c ' N ' /tmp/alle.log) A=$(grep -c ' A ' /tmp/alle.log) D=$(grep -c ' D ' /tmp/alle.log)"
awk -v runs="$runs" '
  FILENAME == "/tmp/marker.log" { rs[$1]=$2; re[$1]=$3; rrc[$1]=$4; dauer[$1]=($3-$2)/1e6; next }
  { t0=$1; t1=$2; n=$3; rc=$4
    k=0
    for (i=1;i<=runs;i++) if (t1 >= rs[i] && t0 <= re[i] + 1.5e9) { k=i; break }
    if (n ~ /^B/ && rc != 124) { hit_b[k]++; if (!(k in f0) || t0 < f0[k]) f0[k]=t0; if (t1 > l1[k]) l1[k]=t1 }
    if (n == "N" && rc == 0) hit_n[k]++
    if (n == "A" && rc != 0) { out_a[k]++; if (!(k in af) || t0 < af[k]) af[k]=t0; if (t1 > al[k]) al[k]=t1 }
    if (n == "D" && rc != 0) { out_d[k]++; if (!(k in df) || t0 < df[k]) df[k]=t0; if (t1 > dl[k]) dl[k]=t1 }
  }
  END {
    printf "%-4s %-8s %-9s %-10s %-13s %-8s %-10s %-8s %-14s\n", "lauf", "init-ms", "B-treffer", "fenster-ms", "fenster-ab-ms", "N-treffer", "A-ausfall", "A-ms", "D-ausfall (ms)"
    for (i=1;i<=runs;i++) {
      w = (hit_b[i] > 0) ? (l1[i]-f0[i])/1e6 : 0
      ab = (hit_b[i] > 0) ? (f0[i]-rs[i])/1e6 : 0
      am = (out_a[i] > 0) ? (al[i]-af[i])/1e6 : 0
      dm = (out_d[i] > 0) ? (dl[i]-df[i])/1e6 : 0
      printf "%-4d %-8.0f %-9d %-10.0f %-13.0f %-8d %-10d %-8.0f %d (%.0f)\n", i, dauer[i], hit_b[i], w, ab, hit_n[i], out_a[i], am, out_d[i], dm
      sb += (hit_b[i] > 0); sa += (out_a[i] > 0); sd += (out_d[i] > 0); sn += (hit_n[i] > 0); if (rrc[i] != 0) bad++
      if (hit_b[i] > 0) { sumw += w; if (w > maxw) maxw = w }
    }
    printf "AUSSAGE laeufe=%d skript-exit-ungleich-0=%d\n", runs, bad
    printf "AUSSAGE laeufe-mit-fenster-B=%d laeufe-mit-N-treffer=%d laeufe-mit-A-ausfall=%d laeufe-mit-D-ausfall=%d\n", sb, sn, sa, sd
    printf "fenster-B: max %.0f ms, Summe %.0f ms\n", maxw, sumw
    printf "AUSSAGE sonden-ausserhalb-der-laeufe: B-treffer=%d A-ausfaelle=%d D-ausfaelle=%d\n", hit_b[0], out_a[0], out_d[0]
  }' /tmp/marker.log /tmp/alle.log
echo
echo "-- Ereignisse: B-Treffer (Exit ungleich 124), A- und D-Ausfaelle, N-Treffer (Zeit in ms seit Start des Laufs; Lauf 0 = ausserhalb)"
awk -v runs="$runs" '
  FILENAME == "/tmp/marker.log" { rs[$1]=$2; re[$1]=$3; next }
  { t0=$1; t1=$2; nm=$3; rc=$4; k=0
    for (i=1;i<=runs;i++) if (t1 >= rs[i] && t0 <= re[i] + 1.5e9) { k=i; break }
    ev = (nm ~ /^B/ && rc != 124) || (nm == "A" && rc != 0) || (nm == "D" && rc != 0) || (nm == "N" && rc == 0)
    if (ev) printf "lauf=%d sonde=%s exit=%d t0=+%.0fms t1=+%.0fms\n", k, nm, rc, (t0 - (k ? rs[k] : t0))/1e6, (t1 - (k ? rs[k] : t1))/1e6
  }' /tmp/marker.log /tmp/alle.log | head -400
