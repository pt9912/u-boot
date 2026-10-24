#!/usr/bin/env bash
# Gegenprobe zu redact() in lib.sh: jede Beispielzeile traegt das Wort PROBEWERT in einer der Formen, in denen ein Token, ein Anmelde-Code,
# eine Anmelde-URL, ein Geraeteschluessel, ein Kontoname oder eine Modell-Antwort in einer Ausgabe stehen kann; nach redact() darf es nirgends
# mehr vorkommen. Neutrale Zeilen muessen unveraendert bleiben. Laeuft VOR dem ersten Messlauf (00-vorlauf.sh schreibt das Artefakt).
# Aufruf: tools/messung/ollama-cloud-konto/redact-probe.sh (Docker; liest und schreibt nichts im Repo).
# Zusage: die Formen der Liste unten; nicht zugesagt: ein Token ohne Praefix und ohne eine dieser Formen.
set -u
here=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
WORK=$(mktemp -d) ROH=$(mktemp -d)
export WORK ROH
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
key=PROBEWERT0123456789ABCDEF
ssh-ed25519 PROBEWERTAAAAC3NzaC1lZDI1NTE5AAAAIJ0123456789
PROBEWERTb3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAAB
You are now signed in as PROBEWERT
contact PROBEWERT@example.com for access
{"response":"PROBEWERT","done":true}
{"message":{"role":"assistant","content":"PROBEWERT"}}
{"type":"result","result":"PROBEWERT text with \"quotes\" and more"}
ANTWORT PROBEWERT'

aus=$(printf '%s\n' "$zeilen" | redact)
n=$(printf '%s\n' "$zeilen" | wc -l)
m=$(printf '%s\n' "$aus" | wc -l)
rc=0
[ "$n" -eq "$m" ] || { echo "FEHLER: $n Zeilen hinein, $m heraus"; rc=1; }
while IFS= read -r z; do
  case $z in *PROBEWERT*) echo "DURCHGELASSEN: $z"; rc=1 ;; esac
done <<< "$aus"

# Schluesselblock ueber mehrere Zeilen: danach darf weder ein Kopf noch ein Koerper noch ein Fuss stehen
block=$(printf '%s\n' '-----BEGIN OPENSSH PRIVATE KEY-----' 'PROBEWERTb3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQ' 'PROBEWERT2222222222222222222222222222222222222' '-----END OPENSSH PRIVATE KEY-----' 'danach neutral' | redact)
case $block in *PROBEWERT*|*BEGIN*|*END\ OPENSSH*) echo "DURCHGELASSEN (Schluesselblock): $block"; rc=1 ;; esac
case $block in *'danach neutral') ;; *) echo "FEHLER: Zeile nach dem Schluesselblock fehlt"; rc=1 ;; esac

# Gegenstueck: Zeilen ohne Token bleiben unveraendert
neutrale=0
for neutral in 'modus=discover benutzer=vscode HOME=/home/vscode' 'OLLAMA_API_KEY OLLAMA_CLOUD_BASE_URL OLLAMA_CONTEXT_LENGTH' 'container-exit=0  gesamtdauer=18.8s  token-anzahl' 'Error: could not connect to ollama server' 'pull model manifest: Get "https://registry.ollama.ai/v2/library/x/manifests/latest": dial tcp: lookup registry.ollama.ai on 127.0.0.1:53: no such host' 'ollama version is 0.40.6' 'You need to be signed in to Ollama to run Cloud models.' 'AUSSAGE C1-run-exit=0' 'AUSSAGE ollama-api-key-wirkung=keine' 'ausgabe: 1 Zeilen, 3 Bytes, nicht leer' '/home/vscode/.ollama/id_ed25519 387 /home/vscode/.ollama' 'Error: ollama cloud is disabled: remote model details are unavailable' '172.17.0.1 host.docker.internal' 'query[A] host.docker.internal from 127.0.0.1' 'docker client=29.9.2 server=29.9.2 os/arch=linux/amd64'; do
  neutrale=$((neutrale + 1))
  [ "$(printf '%s\n' "$neutral" | redact)" = "$neutral" ] || { echo "VERAENDERT: $neutral"; rc=1; }
done
[ "$rc" -eq 0 ] && echo "redact-probe: $n Formen und ein Schluesselblock geschwaerzt, $neutrale neutrale Zeilen unveraendert"
exit "$rc"
