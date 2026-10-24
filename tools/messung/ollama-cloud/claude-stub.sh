#!/bin/bash
# Platzhalter fuer `claude` (HANDVARIANTE, nicht die Auslieferung): zeigt, was `ollama launch claude` an den Prozess uebergibt.
# Ausgewiesen werden die Argumente, die NAMEN der Umgebungsvariablen mit ANTHROPIC, CLAUDE, OLLAMA, API, TOKEN oder KEY und der Wert
# von ANTHROPIC_BASE_URL nur, wenn er auf 127.0.0.1 oder localhost zeigt (sonst steht "anderer Wert"). Kein anderer Wert wird ausgegeben.
echo "STUB claude gestartet, Argumente: $*"
echo "STUB Namen: $(env | cut -d= -f1 | grep -E 'ANTHROPIC|CLAUDE|OLLAMA|API|TOKEN|KEY' | sort | tr '\n' ' ')"
v=${ANTHROPIC_BASE_URL:-}
case "$v" in
  http://127.0.0.1*|http://localhost*) echo "STUB ANTHROPIC_BASE_URL=$v" ;;
  "") echo "STUB ANTHROPIC_BASE_URL ungesetzt" ;;
  *) echo "STUB ANTHROPIC_BASE_URL anderer Wert" ;;
esac
exit 0
