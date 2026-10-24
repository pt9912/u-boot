#!/usr/bin/env bash
# Gegenprobe zu redact() in lib.sh: jede Beispielzeile traegt das Wort PROBEWERT in einer der Formen, in denen ein Token, ein Anmelde-Code
# oder eine Anmelde-URL in einer Ausgabe stehen kann; nach redact() darf es nirgends mehr vorkommen. Laeuft VOR dem ersten Messlauf.
# Aufruf: tools/messung/ollama-cloud/redact-probe.sh (Docker; liest und schreibt nichts im Repo).
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
X-Gitlab-Token: PROBEWERT
OLLAMA_API_KEY=PROBEWERT
ANTHROPIC_API_KEY=PROBEWERT
ANTHROPIC_AUTH_TOKEN=PROBEWERT
export GITHUB_TOKEN=PROBEWERT
{"token":"PROBEWERT"}
{"access_token": "PROBEWERT", "x":1}
{"password":"PROBEWERT"}
git clone https://user:PROBEWERT@gitlab.com/x/y.git
https://example.com/f?X-Amz-Signature=PROBEWERT
sk-ant-PROBEWERT1234
sk-PROBEWERT1234
ollama_PROBEWERT1234
ghp_PROBEWERT1234
You need to sign in to Ollama. Visit https://ollama.com/connect?name=PROBEWERT&key=PROBEWERT to sign in
Open this URL: https://ollama.com/connect?name=host&key=PROBEWERT
please open https://ollama.com/signin#PROBEWERT
your code: PROBEWERT0123456789ABCDEF
Public key: PROBEWERT0123456789ABCDEFGHIJ
key=PROBEWERT0123456789ABCDEF'

aus=$(printf '%s\n' "$zeilen" | redact)
n=$(printf '%s\n' "$zeilen" | wc -l)
m=$(printf '%s\n' "$aus" | wc -l)
rc=0
[ "$n" -eq "$m" ] || { echo "FEHLER: $n Zeilen hinein, $m heraus"; rc=1; }
while IFS= read -r z; do
  case $z in *PROBEWERT*) echo "DURCHGELASSEN: $z"; rc=1 ;; esac
done <<< "$aus"
# Gegenstueck: Zeilen ohne Token bleiben unveraendert
for neutral in 'modus=discover benutzer=vscode HOME=/home/vscode' 'OLLAMA_API_KEY OLLAMA_CLOUD_BASE_URL OLLAMA_CONTEXT_LENGTH' 'container-exit=0  gesamtdauer=18.8s  token-anzahl' 'Error: could not connect to ollama server' 'pull model manifest: Get "https://registry.ollama.ai/v2/library/x/manifests/latest": dial tcp: lookup registry.ollama.ai: no such host'; do
  [ "$(printf '%s\n' "$neutral" | redact)" = "$neutral" ] || { echo "VERAENDERT: $neutral"; rc=1; }
done
[ "$rc" -eq 0 ] && echo "redact-probe: $n Formen geschwaerzt, neutrale Zeilen unveraendert"
exit "$rc"
