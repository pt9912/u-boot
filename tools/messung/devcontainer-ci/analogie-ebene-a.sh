# analogie-ebene-a.sh — M5 zu slice-v1-bau-nur-devcontainer: Container A.
# Läuft unprivilegiert (kein --privileged), nur dem Host-Docker-Socket gemountet —
# die Analogie zum CI-Job-Container auf einem GitHub-Runner (Läufe über docker.sock, nicht über
# einen privilegierten Container). Startet darauf zwei Formen der Ebene 2:
# a) Devcontainer-Container mit --privileged (das runArgs-Bild des Projekts) — DIND-Ebene 2;
# b) Devcontainer-Container ohne --privileged — die Frage, ob es unprivilegiert läuft.
# Aufruf siehe messen.sh Modus analogie.
set -uo pipefail
echo "== Container A: $(id -un) uid=$(id -u) unprivilegiert gestartet (kein --privileged)"
echo "CapEff=$(grep CapEff /proc/self/status | cut -d' ' -f2)"
echo "docker (Client) in A: $(docker --version) — Server ist der layer-1-Daemon über den Socket:"
docker version --format 'layer-1 server={{.Server.Version}}'
echo
echo "== (a) Ebene 2 mit --privileged: der Devcontainer-Container, DIND-Analogie"
# Bind-Quellen eines `-v` löst der layer-1-Daemon auf dem Host auf (nicht in A) —
# das Skript kommt deshalb über stdin in den Ebene-2-Container, nicht als Pfad-Mount.
docker run -i --rm --privileged \
  -v uboot-m1-docker-analogie-priv:/var/lib/docker \
  uboot-devc-m1 bash -s < /m/analogie-ebene2.sh
echo "ebene2-priv-exit=$?"
echo
echo "== (b) Ebene 2 ohne --privileged: was fehlt (Erwartung: der Daemon-Start fällt zurück)"
# Der Startscript fällt nach Policy (onUnavailable: warn) auf fallback zurück.
docker run -i --rm \
  -v uboot-m1-docker-analogie-nop:/var/lib/docker \
  uboot-devc-m1 bash -s < /m/analogie-ebene2-nop.sh
echo "ebene2-nop-exit=$?"