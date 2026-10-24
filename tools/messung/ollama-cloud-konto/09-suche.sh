#!/bin/bash
# Suche nach Zugangsdaten ueber Rohdaten und Skripte (DoD 1, iii): Tokenformen, Schluesselbloecke, Authorization-Header, Zugangsdaten in URLs,
# Zuweisungen an Namen mit Zugangswort, Anmelde-URLs mit Query, Schluesselparameter. Treffer werden mit Datei und Zeile ausgewiesen (die Zeile
# gekuerzt), Exit 1 bei einem Treffer ausserhalb der Ausnahmeliste. Treffer sind von Hand zu lesen; die Suche ersetzt das Lesen nicht.
#
# AUSNAHMELISTE (benannt, je mit Grund):
#  Zeichenketten (werden aus der Zeile entfernt, bevor sie erneut geprueft wird; ein echter Wert in derselben Zeile bleibt ein Treffer):
#    ANTHROPIC_AUTH_TOKEN=ollama      Platzhalter, den Claude Code gegen einen Ollama-Server braucht; kein Geheimnis (Angabe des Projektinhabers 2026-10-24)
#    OLLAMA_API_KEY=platzhalter       Platzhalter fuer die Wirkungsprobe von OLLAMA_API_KEY; kein gueltiger Schluessel (Frage 4 des Slice)
#  Pfade (Treffer werden gezaehlt und ausgewiesen, zaehlen nicht als Treffer; nur bei --alle, ueber ganz docs/messungen und tools/messung):
#    docs/messungen/slice-v1-devcontainer-ollama-amd64-messung/   Testnamen mit Ollama_…, GPG-Fingerabdruecke, Podman-Signing-Key-ID (oeffentlich)
#    docs/messungen/slice-v1-devcontainer-ollama-messung/         oeffentliche Geraeteschluessel, die `ollama serve` ausgibt (Public key; der private Teil wurde nie ausgegeben)
#    docs/messungen/slice-v1-forge-cli-quellen-messung/           Prosa ueber die Formen der Schwaerzung, GPG-Fingerabdruck der Quelle
#    tools/messung/forge-cli/inner-07.sh                          Zuweisungen mit $DUMMY (Platzhalter)
#    tools/messung/devcontainer-start-befehl/01-einzelmessung.sh  Platzhalter "platzhalter-kein-zugangsdatum-…"
#  Dateien, die die Formen selbst tragen: lib.sh, redact-probe.sh, 06-suche.sh / 09-suche.sh der Messskript-Verzeichnisse.
#
# Aufruf: tools/messung/ollama-cloud-konto/09-suche.sh [--alle] [verzeichnis...]
#   Vorgabe: das Verzeichnis dieser Messung und tools/messung/ollama-cloud-konto/ (muss leer sein); --alle: docs/messungen und tools/messung mit der Pfad-Ausnahmeliste
set -u
here=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
ROOT=$(cd "$here/../../.." && pwd)
alle=0; [ "${1:-}" = --alle ] && { alle=1; shift; }
if [ $# -eq 0 ]; then
  if [ $alle = 1 ]; then set -- "$ROOT/docs/messungen" "$ROOT/tools/messung"
  else set -- "$ROOT/docs/messungen/slice-v1-ollama-cloud-messung-mit-konto" "$ROOT/tools/messung/ollama-cloud-konto"; fi
fi
pfad_ausnahmen=(
  'docs/messungen/slice-v1-devcontainer-ollama-amd64-messung/'
  'docs/messungen/slice-v1-devcontainer-ollama-messung/'
  'docs/messungen/slice-v1-forge-cli-quellen-messung/'
  'tools/messung/forge-cli/inner-07.sh'
  'tools/messung/devcontainer-start-befehl/01-einzelmessung.sh'
)
selbst='/(ollama-cloud|ollama-cloud-konto)/(lib|redact-probe|06-suche|09-suche)\.sh:'
rc=0
shopt -s nocasematch
formen=(
  'ghp_[A-Za-z0-9]{10,}|gho_[A-Za-z0-9]{10,}|github_pat_[A-Za-z0-9_]{10,}|glpat-[A-Za-z0-9_-]{10,}|sk-ant-[A-Za-z0-9_-]{10,}|sk-[A-Za-z0-9]{20,}|ollama_[A-Za-z0-9]{16,}'
  'BEGIN [A-Z ]*PRIVATE KEY|ssh-ed25519 [A-Za-z0-9+/]{20,}|ssh-rsa [A-Za-z0-9+/]{20,}'
  '(Authorization|Bearer) *:? *[A-Za-z0-9._~+/=-]{12,}'
  '://[^/ :@]+:[^/ @]+@'
  '(token|secret|passw(or)?d|api_?key)[A-Za-z0-9_]*=[^ $%[]'
  'https?://[^ ]*(connect|signin|device)[^ ]*\?[^ ]+'
  'key=[A-Za-z0-9+/_-]{16,}'
  '[A-Za-z0-9._%+-]+@[A-Za-z0-9-]+\.[A-Za-z0-9.-]*[A-Za-z]{2,}'
)
ausn_pfad=0; ausn_str=0
for f in "${formen[@]}"; do
  treffer=""
  while IFS= read -r l; do
    [ -n "$l" ] || continue
    [[ $l =~ \[entfernt\]|\[Query\ entfernt|PROBEWERT|\[E-Mail\ entfernt\] ]] && continue
    [[ $l =~ $selbst ]] && continue
    p=${l%%:*}; rel=${p#"$ROOT"/}
    ok=0
    for a in "${pfad_ausnahmen[@]}"; do [[ $alle = 1 && $rel == "$a"* ]] && { ok=1; ausn_pfad=$((ausn_pfad + 1)); break; }; done
    [ $ok = 1 ] && continue
    rest=${l#*:}; rest=${rest#*:}
    r2=${rest//ANTHROPIC_AUTH_TOKEN=ollama/}; r2=${r2//OLLAMA_API_KEY=platzhalter/}
    if [ "$r2" != "$rest" ] && ! [[ $r2 =~ $f ]]; then ausn_str=$((ausn_str + 1)); continue; fi
    treffer+="$(printf '%s' "$l" | cut -c1-200)"$'\n'
  done < <(grep -r -n -I -E -i -e "$f" "$@" 2> /dev/null)
  [ -z "$treffer" ] || { echo "TREFFER fuer Muster: $f"; printf '%s' "$treffer"; rc=1; }
done
echo "ausnahmen: $ausn_pfad Treffer in den Pfaden der Ausnahmeliste (nur mit --alle), $ausn_str Zeilen nur mit einem benannten Platzhalter"
[ $rc = 0 ] && echo "suche: keine Treffer in $* (Formen: Praefix-Token, Schluessel, Header, URL-Zugang, Zuweisung, Anmelde-URL, Schluesselparameter, E-Mail)"
exit $rc
