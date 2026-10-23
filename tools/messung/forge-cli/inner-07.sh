#!/bin/bash
# Laeuft IM Container (Image einer Variante) als Benutzer vscode, wie im Devcontainer; aufgerufen von 07-laufzeit.sh.
# Aufruf im Container: inner-07.sh <modus> <werkzeug>
#   modus      discover    keine Restriktion; protokollierender dnsmasq reicht alle Namen weiter (zeigt, welche Hosts ein Aufruf braucht)
#              restricted  das erzeugte Startscript unveraendert (`sudo sh /usr/local/bin/u-boot-egress-init`); danach ersetzt ein
#                          protokollierender dnsmasq mit derselben Konfiguration den Resolver des Skripts (Verfahren von ollama-amd64)
#   werkzeug   gh | glab
# Es wird kein echtes Zugangsdatum verwendet: $DUMMY (Umgebung) ist ein Platzhalter. GIT_TOKEN ist, wenn gesetzt, derselbe Platzhalter.
set -u
mode=$1; tool=$2
log=/tmp/dns.log
DUMMY=${DUMMY:-MESSUNG-KEIN-ECHTES-TOKEN}
now() { date -u +%T.%N | cut -c1-12; }
mark() { echo "---- $(now) $*"; }
sudo touch "$log"; sudo chmod 666 "$log"
up=$(awk '/^nameserver/ && $2 !~ /:/ {print $2}' /etc/resolv.conf | head -1)

if [ "$mode" = discover ]; then
  printf 'listen-address=127.0.0.1\nbind-interfaces\nno-resolv\nno-hosts\nserver=%s\nlog-queries\nlog-facility=%s\n' "$up" "$log" | sudo tee /tmp/dns-discover.conf > /dev/null
  sudo dnsmasq --conf-file=/tmp/dns-discover.conf
  echo "nameserver 127.0.0.1" | sudo tee /etc/resolv.conf > /dev/null
else
  sudo sh /usr/local/bin/u-boot-egress-init; echo "egress-init exit=$?"
  echo "-- Status-Datei des Skripts"; cat /var/lib/u-boot/status/egress
  grep -E '^server=' /etc/u-boot-egress-dnsmasq.conf | sed 's/^/conf: /'
  sudo kill "$(cat /run/u-boot-egress-dnsmasq.pid)"; sleep 1
  { cat /etc/u-boot-egress-dnsmasq.conf; echo "log-queries"; echo "log-facility=$log"; } | sudo tee /tmp/dns-restricted.conf > /dev/null
  sudo dnsmasq --conf-file=/tmp/dns-restricted.conf
fi
sleep 1
hostsline=$(grep '^hosts=' /usr/local/bin/u-boot-egress-init 2>/dev/null | cut -d'"' -f2)
echo "modus=$mode werkzeug=$tool home=$HOME user=$(id -un)"
[ "$mode" != discover ] && echo "Hosts der Allowlist im Startscript: $hostsline"

# probe <etikett> <befehl...>: fuehrt aus, zeigt Ausgabe (gekuerzt), Exit-Code, DNS-Namen waehrend des Aufrufs (+ Antwort ausserhalb der Liste)
# und neue/geaenderte Dateien unter $HOME. Umgebungsvariablen vor dem Aufruf gelten nur fuer ihn (env VAR=… probe …).
probe() {
  local label=$1; shift
  local n0 m out rc
  n0=$(wc -l < "$log"); m=$(mktemp); sleep 0.2
  mark "$label"
  echo "\$ $*"
  "$@" > "$m.out" 2>&1 < /dev/null; rc=$?
  sed -e 's/\x1b\[[0-9;?]*[a-zA-Z]//g' "$m.out" | tr '\r' '\n' | grep -v -E '^\s*$' | head -"${LINES_MAX:-10}" | cut -c1-200
  rm -f "$m.out"
  echo "exit=$rc"
  echo "AUSSAGE $mode-$tool-$label-exit=$rc"
  echo "-- DNS-Anfragen (A/AAAA)"
  tail -n +$((n0 + 1)) "$log" | grep -E 'query\[(A|AAAA)\]' | awk '{print $6}' | sort | uniq -c
  if [ "$mode" != discover ]; then
    echo "-- Antworten ausserhalb der Liste (REFUSED/NXDOMAIN/NODATA/config)"
    tail -n +$((n0 + 1)) "$log" | grep -E 'REFUSED|NXDOMAIN|NODATA|config ' | awk '{print $4, $5, $6}' | sort -u | head -5
  fi
  echo "-- geaenderte Dateien unter HOME (Pfad, Bytes)"
  find "$HOME" -newer "$m" -type f -printf '%p %s\n' 2>/dev/null | sort | sed "s#^$HOME#~#"
  rm -f "$m"
}
# Dateien, die unter HOME existieren, mit Pruefung auf den Platzhalter (Anzahl der Treffer, nie der Inhalt)
platzhalter_in_dateien() {
  echo "-- git config --global --get-regexp credential (Eintraege, die ein Werkzeug geschrieben haette)"; git config --global --get-regexp "^credential" 2>&1 | cut -c1-200 || true
  echo "-- Dateien unter HOME ausserhalb der Image-Grundausstattung (Pfad, Bytes, Zeilen mit dem Platzhalter)"
  find "$HOME" -type f \( -path "$HOME/.config/gh/*" -o -path "$HOME/.config/glab-cli/*" -o -path "$HOME/.local/state/*" -o -path "$HOME/.cache/*" -o -path "$HOME/.local/share/*" -o -path "$HOME/.gitconfig" \) 2>/dev/null | sort | while read -r f; do
    echo "$(echo "$f" | sed "s#^$HOME#~#") $(stat -c %s "$f") platzhalter-treffer=$(grep -c -F "$DUMMY" "$f")"
  done
}
mkdir -p /tmp/dl
echo "-- Git-Konfiguration im Image (system, vom Dockerfile der Vorlage) und im Home (global)"
git config --system --get-regexp "^credential" | cut -c1-200
git config --global --list 2>&1 | cut -c1-200
if [ "$tool" = gh ]; then
  probe version gh --version
  LINES_MAX=14 probe help-environment gh help environment
  probe auth-status-ohne-variable gh auth status
  probe api-ohne-variable gh api rate_limit
  GH_TOKEN="$DUMMY" probe auth-status-gh_token gh auth status
  GH_TOKEN="$DUMMY" probe api-gh_token gh api rate_limit
  GITHUB_TOKEN="$DUMMY" probe api-github_token gh api rate_limit
  GH_TOKEN="$DUMMY" probe repo-view-gh_token gh repo view cli/cli
  GH_TOKEN="$DUMMY" probe release-download-gh_token gh release download v2.103.0 -R cli/cli -p 'gh_*_checksums.txt' -D /tmp/dl
  probe git-ls-remote-ohne-git_token git ls-remote https://github.com/cli/cli.git HEAD
  GIT_TOKEN="$DUMMY" probe git-ls-remote-git_token git ls-remote https://github.com/cli/cli.git HEAD
  probe config-list gh config list
  LINES_MAX=12 probe env-hilfe-telemetrie-update sh -c "gh help environment | grep -i -E \"telemetry|update|NO_\""
  GH_TOKEN="$DUMMY" probe repo-clone-gh_token gh repo clone cli/cli /tmp/clone-gh -- --depth 1
  probe auth-git-credential-hilfe gh auth setup-git --help
  GH_TOKEN="$DUMMY" probe auth-setup-git gh auth setup-git --hostname github.com
  platzhalter_in_dateien
else
  probe version glab --version
  probe auth-status-ohne-variable glab auth status
  probe api-ohne-variable glab api projects/gitlab-org%2Fcli --silent
  probe repo-view-ohne-variable glab repo view gitlab-org/cli
  probe release-list-ohne-variable glab release list -R gitlab-org/cli
  probe release-download-ohne-variable glab release download v1.124.0 -R gitlab-org/cli --asset-name checksums.txt -D /tmp/dl
  GITLAB_TOKEN="$DUMMY" probe auth-status-gitlab_token glab auth status
  GITLAB_TOKEN="$DUMMY" probe api-gitlab_token glab api projects/gitlab-org%2Fcli --silent
  GITLAB_ACCESS_TOKEN="$DUMMY" probe api-gitlab_access_token glab api projects/gitlab-org%2Fcli --silent
  OAUTH_TOKEN="$DUMMY" probe api-oauth_token glab api projects/gitlab-org%2Fcli --silent
  GITLAB_TOKEN="$DUMMY" probe repo-view-gitlab_token glab repo view gitlab-org/cli
  probe git-ls-remote-ohne-git_token git ls-remote https://gitlab.com/gitlab-org/cli.git HEAD
  GIT_TOKEN="$DUMMY" probe git-ls-remote-git_token git ls-remote https://gitlab.com/gitlab-org/cli.git HEAD
  probe config-git-protocol glab config get git_protocol
  probe config-telemetrie-update-duo sh -c "for k in telemetry check_update duo_cli_auto_download duo_cli_auto_run show_whats_new; do printf \"%s=\" \$k; glab config get \$k; done"
  probe auth-login-hilfe-git-credential-helper sh -c "glab auth login --help | grep -i -E \"credential|stdin\""
  probe repo-clone-ssh-standard glab repo clone gitlab-org/cli /tmp/clone-ssh -- --depth 1
  probe config-git-protocol-https glab config set git_protocol https --global
  probe repo-clone-https glab repo clone gitlab-org/cli /tmp/clone-https -- --depth 1
  GIT_TOKEN="$DUMMY" probe repo-clone-https-git_token glab repo clone gitlab-org/cli /tmp/clone-https2 -- --depth 1
  platzhalter_in_dateien
  echo "-- Schluessel (nicht Werte) in ~/.config/glab-cli/config.yml"; awk -F: "/^[a-z_]+:/ {print \$1}" "$HOME/.config/glab-cli/config.yml" 2>/dev/null | tr "\n" " "; echo
fi
echo
mark "alle DNS-Anfragen der Sitzung (eindeutig, A/AAAA)"
grep -E 'query\[(A|AAAA)\]' "$log" | awk '{print $6}' | sort | uniq -c
