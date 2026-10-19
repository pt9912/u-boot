#!/bin/bash
# Reproduziert die Messung der Fälle a bis d (Egress, Subdomains und Wildcard).
# Aufruf: run.sh <repo-wurzel> [fall...]   mit Fall aus: a a2 b c d2 d3 cd-c cd-d
# Voraussetzung: Docker mit Zugriff auf `debian:stable` und `--cap-add=NET_ADMIN`.
# Der Container führt das Startscript `egress-init.sh.tmpl` mit ersetzter
# `hosts=`-Zeile aus (inner.sh) beziehungsweise die Messungen ohne Allowlist
# (inner-cd.sh). Die Ausgabe entspricht den Dateien rohdaten-1.txt bis -3.txt.
set -eu
root=$(cd "${1:?repo-wurzel}" && pwd)
shift
here=$(cd "$(dirname "$0")" && pwd)
tmpl="$root/internal/hexagon/application/templates/devcontainer/egress-init.sh.tmpl"
cases=${*:-a a2 b c d2 d3 cd-c cd-d}

run() { # run <script> <args...>
  docker run --rm --cap-add=NET_ADMIN -v "$tmpl":/tmpl/egress-init.sh.tmpl:ro \
    -v "$here":/m:ro debian:stable bash "/m/$1" "${@:2}"
}

for c in $cases; do
  case "$c" in
    a)  echo "== a: allow example.org"
        run inner.sh "example.org" example.org www.example.org a.b.example.org example.com ;;
    a2) echo "== a2: allow amazonaws.com"
        run inner.sh "amazonaws.com" s3.amazonaws.com s3.us-east-1.amazonaws.com sts.us-east-1.amazonaws.com example.com ;;
    b)  echo "== b: allow api.github.com"
        run inner.sh "api.github.com" api.github.com github.com www.github.com uploads.api.github.com ;;
    c)  echo "== c: allow www.mozilla.org files.pythonhosted.org deb.debian.org"
        run inner.sh "www.mozilla.org files.pythonhosted.org deb.debian.org" \
          www.mozilla.org www-mozilla.fastly-edge.com files.pythonhosted.org dualstack.python.map.fastly.net deb.debian.org ;;
    d2) echo "== d2: allow *.example.org (full script, nftset + server)"
        run inner.sh "*.example.org" example.org www.example.org a.b.example.org example.com ;;
    d3) echo "== d3: allow .example.org"
        run inner.sh ".example.org" example.org www.example.org example.com ;;
    cd-c) run inner-cd.sh c ;;
    cd-d) run inner-cd.sh d ;;
    *)  echo "unbekannter Fall: $c" >&2; exit 2 ;;
  esac
done
