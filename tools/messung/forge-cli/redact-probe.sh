#!/usr/bin/env bash
# Gegenprobe zu redact() in lib.sh: jede Beispielzeile traegt das Wort PROBEWERT in einer der Formen, in denen ein Token in einer Ausgabe
# stehen kann; nach redact() darf es nirgends mehr vorkommen. Aufruf: tools/messung/forge-cli/redact-probe.sh (Docker; liest und schreibt nichts im Repo).
# Zusage: die Formen der Liste unten; nicht zugesagt: ein Token ohne Praefix und ohne eine dieser Formen.
set -u
here=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
WORK=$(mktemp -d) OUT=$(mktemp -d)
export WORK OUT
. "$here/lib.sh"

zeilen='Authorization: Bearer PROBEWERT
authorization: token PROBEWERT
Authorization: PROBEWERT
PRIVATE-TOKEN: PROBEWERT
private-token: PROBEWERT
JOB-TOKEN: PROBEWERT
X-Gitlab-Token: PROBEWERT
GITLAB_TOKEN=PROBEWERT
GH_TOKEN=PROBEWERT
export GITHUB_TOKEN=PROBEWERT
{"token":"PROBEWERT"}
{"access_token": "PROBEWERT", "x":1}
{"password":"PROBEWERT"}
git clone https://user:PROBEWERT@gitlab.com/x/y.git
https://x-access-token:PROBEWERT@github.com/x/y
x-access-token:PROBEWERT@github.com
https://example.com/f?X-Amz-Signature=PROBEWERT
ghp_PROBEWERT1234
glpat-PROBEWERT1234
github_pat_PROBEWERT1234'

aus=$(printf '%s\n' "$zeilen" | redact)
n=$(printf '%s\n' "$zeilen" | wc -l)
m=$(printf '%s\n' "$aus" | wc -l)
rc=0
[ "$n" -eq "$m" ] || { echo "FEHLER: $n Zeilen hinein, $m heraus"; rc=1; }
while IFS= read -r z; do
  case $z in *PROBEWERT*) echo "DURCHGELASSEN: $z"; rc=1 ;; esac
done <<< "$aus"
# Gegenstueck: eine Zeile ohne Token bleibt unveraendert
neutral='container-exit=0  gesamtdauer=18.8s  token-anzahl'
[ "$(printf '%s\n' "$neutral" | redact)" = "$neutral" ] || { echo "VERAENDERT: $neutral"; rc=1; }
[ "$rc" -eq 0 ] && echo "redact-probe: $n Formen geschwaerzt, neutrale Zeile unveraendert"
exit "$rc"
