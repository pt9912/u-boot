#!/bin/bash
# O5 (Modellspeicher, Volume) und O6 (Start): Probe-Image aus 04-build-probe.sh, Container wie der Devcontainer (sleep infinity,
# NET_ADMIN, Befehle per `docker exec` als vscode = postStartCommand).
# Aufruf: 10-volume-und-start.sh > rohdaten-10.txt 2>&1
set -u
P=u-boot-messung-ollama-probe
D="api.github.com codeload.github.com deb.debian.org github.com objects.githubusercontent.com raw.githubusercontent.com security.debian.org registry.ollama.ai dd20bb891979d25aebc8bec07b2b3bbc.r2.cloudflarestorage.com"
strip() { sed -e 's/\x1b\[[0-9;?]*[a-zA-Z]//g' | tr '\r' '\n' | grep -v -E '^\s*$' | awk '!seen[$0]++'; }
cleanup() { docker rm -f ollama-m1 ollama-m2 ollama-m3 ollama-m4 >/dev/null 2>&1; docker volume rm -f u-boot-messung-ollama-v1 u-boot-messung-ollama-v2 >/dev/null 2>&1; }
cleanup

echo "######## V0: ohne Volume (Container-Schicht): Modell nach stop/start und nach Loeschen des Containers"
docker run -d --name ollama-m1 $P sleep infinity >/dev/null
docker exec ollama-m1 bash -c '(ollama serve > /tmp/s.log 2>&1 &); sleep 3; ollama pull smollm2:135m 2>&1 | tr "\r" "\n" | tail -1; ollama list | tail -n +2'
docker stop -t 2 ollama-m1 >/dev/null; docker start ollama-m1 >/dev/null; sleep 1
echo "-- nach stop/start desselben Containers:"
docker exec ollama-m1 bash -c 'ls /home/vscode/.ollama/models/manifests >/dev/null 2>&1 && echo "Manifeste noch da"; du -sh /home/vscode/.ollama | cut -f1'
docker rm -f ollama-m1 >/dev/null
docker run -d --name ollama-m1 $P sleep infinity >/dev/null
echo "-- neuer Container aus demselben Image:"
docker exec ollama-m1 bash -c 'ls /home/vscode/.ollama 2>&1 | head -2'
docker rm -f ollama-m1 >/dev/null
echo

echo "######## V1: benanntes Volume auf /home/vscode/.ollama, Verzeichnis im Image NICHT vorhanden"
docker run -d --name ollama-m1 -v u-boot-messung-ollama-v1:/home/vscode/.ollama $P sleep infinity >/dev/null
docker exec ollama-m1 bash -c 'ls -ld /home/vscode/.ollama; (ollama serve > /tmp/s.log 2>&1 &) ; sleep 3; head -5 /tmp/s.log | cut -c1-200; curl -sS -m 2 http://127.0.0.1:11434/ || true'
echo

echo "######## V2: Image mit vorbereitetem Verzeichnis (mkdir + chown vscode), benanntes Volume"
docker build -q -t u-boot-messung-ollama-probe-v2 - <<DF
FROM $P
USER root
RUN mkdir -p /home/vscode/.ollama && chown vscode:vscode /home/vscode/.ollama
USER vscode
DF
docker run -d --name ollama-m2 --cap-add=NET_ADMIN -v u-boot-messung-ollama-v2:/home/vscode/.ollama u-boot-messung-ollama-probe-v2 sleep infinity >/dev/null
docker exec ollama-m2 bash -c 'ls -ld /home/vscode/.ollama'
echo "-- postStart-Kette: Egress-Skript, dann serve (Hosts: Standard + registry.ollama.ai + Blob-Host)"
docker exec ollama-m2 bash -c "sed -e 's|^hosts=.*|hosts=\"$D\"|' /usr/local/bin/u-boot-egress-init > /tmp/e.sh; sudo sh /tmp/e.sh && (nohup ollama serve > /tmp/s.log 2>&1 &); echo chain-exit=\$?; sleep 3; ss -ltnp | grep 11434 | cut -c1-120; ps -o pid,user,args -C ollama"
docker exec ollama-m2 bash -c 'ollama pull qwen2.5:0.5b 2>&1 | tr "\r" "\n" | tail -2; du -sh /home/vscode/.ollama; ls -la /home/vscode/.ollama | tail -n +2' | strip
echo "-- docker stop / start desselben Containers: laeuft serve danach?"
docker stop -t 2 ollama-m2 >/dev/null; docker start ollama-m2 >/dev/null; sleep 1
docker exec ollama-m2 bash -c 'ps -o pid,user,args -C ollama || echo "kein ollama-Prozess nach dem Neustart"; ls /home/vscode/.ollama/models/manifests >/dev/null && echo "Modell-Manifeste noch da"; echo "-- resolv.conf nach Neustart:"; cat /etc/resolv.conf | grep nameserver; nft list tables 2>&1 | head -3'
echo "-- postStart erneut (wie bei jedem Start durch die Dev-Container-CLI)"
docker exec ollama-m2 bash -c "sudo sh /tmp/e.sh 2>&1 | tail -2; (nohup ollama serve > /tmp/s.log 2>&1 &); sleep 3; ollama list; ss -ltn | grep 11434"
echo "-- Container loeschen, neuen Container mit demselben Volume: Modell ohne Netz da?"
docker rm -f ollama-m2 >/dev/null
docker run -d --name ollama-m3 --cap-add=NET_ADMIN -v u-boot-messung-ollama-v2:/home/vscode/.ollama u-boot-messung-ollama-probe-v2 sleep infinity >/dev/null
docker exec ollama-m3 bash -c "sed -e 's|^hosts=.*|hosts=\"api.github.com\"|' /usr/local/bin/u-boot-egress-init > /tmp/e.sh; sudo sh /tmp/e.sh && (nohup ollama serve > /tmp/s.log 2>&1 &); sleep 3; ollama list; ollama run qwen2.5:0.5b 'Reply with the single word OK.' 2>&1 | tr '\r' '\n' | tail -1; echo; ollama ps | tail -1"
docker rm -f ollama-m3 >/dev/null
echo

echo "######## S: Egress-Skript nicht verfuegbar (ohne NET_ADMIN): Kette sudo egress-init && serve"
for pol in warn fail; do
  docker rm -f ollama-m4 >/dev/null 2>&1; docker run -d --name ollama-m4 u-boot-messung-ollama-probe-v2 sleep infinity >/dev/null
  docker exec ollama-m4 bash -c "sed -e 's|^on_unavailable=.*|on_unavailable=\"$pol\"|' /usr/local/bin/u-boot-egress-init > /tmp/e-$pol.sh; sudo sh /tmp/e-$pol.sh > /tmp/o.txt 2>&1 && (nohup ollama serve > /tmp/s.log 2>&1 &); rc=\$?; sleep 2; cut -c1-140 /tmp/o.txt; echo \"onUnavailable=$pol kette-exit=\$rc\"; ps -o pid,user,args -C ollama || echo 'serve nicht gestartet'"
done
cleanup
docker rmi -f u-boot-messung-ollama-probe-v2 >/dev/null 2>&1
