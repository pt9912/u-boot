#!/bin/bash
# Suche nach Zugangsdaten ueber Rohdaten und Skripte (DoD "Keine Zugangsdaten"): Tokenformen, Zuweisungen an Namen mit Zugangswort,
# Authorization-Header, Zugangsdaten in URLs, Schluesselbloecke, Anmelde-URLs mit Query. Treffer werden mit Datei und Zeile ausgewiesen
# (die Zeile gekuerzt), Exit 1 bei einem Treffer. Bereits geschwaerzte Zeilen ([entfernt]) und die Dateien, die die Formen selbst
# tragen (lib.sh, redact-probe.sh, diese Datei), zaehlen nicht. Treffer sind von Hand zu lesen; die Suche ersetzt das Lesen nicht.
# Aufruf: tools/messung/ollama-cloud/06-suche.sh [verzeichnis...]  (Vorgabe: docs/messungen/ und tools/messung/)
set -u
here=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
ROOT=$(cd "$here/../../.." && pwd)
[ $# -gt 0 ] || set -- "$ROOT/docs/messungen" "$ROOT/tools/messung"
rc=0
formen=(
  'ghp_[A-Za-z0-9]{10,}|gho_[A-Za-z0-9]{10,}|github_pat_[A-Za-z0-9_]{10,}|glpat-[A-Za-z0-9_-]{10,}|sk-ant-[A-Za-z0-9_-]{10,}|sk-[A-Za-z0-9]{20,}|ollama_[A-Za-z0-9]{16,}'
  'BEGIN [A-Z ]*PRIVATE KEY|ssh-ed25519 [A-Za-z0-9+/]{20,}|ssh-rsa [A-Za-z0-9+/]{20,}'
  '(Authorization|Bearer) *:? *[A-Za-z0-9._~+/=-]{12,}'
  '://[^/ :@]+:[^/ @]+@'
  '(token|secret|passw(or)?d|api_?key)[A-Za-z0-9_]*=[^ $%[]'
  'https?://[^ ]*(connect|signin|device)[^ ]*\?[^ ]+'
  'key=[A-Za-z0-9+/_-]{16,}'
)
for f in "${formen[@]}"; do
  out=$(grep -r -n -I -E -i -e "$f" "$@" 2> /dev/null | cut -c1-200 | grep -v -E '\[entfernt\]|\[Query entfernt|PROBEWERT|06-suche\.sh|/lib\.sh|redact-probe\.sh')
  [ -z "$out" ] || { echo "TREFFER fuer Muster: $f"; echo "$out"; rc=1; }
done
[ $rc = 0 ] && echo "suche: keine Treffer in $* (Formen: Praefix-Token, Schluessel, Header, URL-Zugang, Zuweisung, Anmelde-URL)"
exit $rc
