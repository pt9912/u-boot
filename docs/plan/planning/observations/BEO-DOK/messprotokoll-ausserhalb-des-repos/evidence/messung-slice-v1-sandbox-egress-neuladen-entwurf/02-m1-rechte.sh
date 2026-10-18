#!/bin/bash
# M1: Was kann ein Prozess im Sandbox-Container ohne und mit `sudo` an der Egress-Beschraenkung aendern?
# Aufruf: 02-m1-rechte.sh > rohdaten-02-m1.txt 2>&1
# Voraussetzung: Docker (Colima); Image u-boot-messung-egress-neuladen-dig (siehe README.md, Schritt 0).
# Der Container laeuft wie der Devcontainer (sleep infinity, --cap-add=NET_ADMIN); der postStartCommand
# aus devcontainer.json (`sudo sh /usr/local/bin/u-boot-egress-init`) laeuft per `docker exec` als vscode.
set -u
here=$(cd "$(dirname "$0")" && pwd)
IMG=u-boot-messung-egress-neuladen-dig
c=egress-m1
docker rm -f $c >/dev/null 2>&1
docker run -d --name $c --cap-add=NET_ADMIN $IMG sleep infinity >/dev/null
echo "######## postStartCommand (als vscode)"
docker exec -u vscode $c sh -c 'sudo sh /usr/local/bin/u-boot-egress-init'; echo "postStart exit=$?"
echo "######## Phase ohne sudo (als vscode)"
docker exec -i -u vscode $c bash -s ohne < "$here/inner-m1.sh"
echo "######## Phase mit sudo (als vscode)"
docker exec -i -u vscode $c bash -s mit < "$here/inner-m1.sh"
docker rm -f $c >/dev/null 2>&1
