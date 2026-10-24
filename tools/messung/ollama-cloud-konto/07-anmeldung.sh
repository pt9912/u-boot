#!/bin/bash
# Anmeldung des Ollama-Servers im Wegwerf-Container (Weg b, K4). NUR der Projektinhaber fuehrt dieses Skript aus, an seinem Terminal.
# Die Anmelde-URL und der Code erscheinen ausschliesslich auf diesem Terminal: das Skript schreibt keine Datei, leitet nichts um und bricht
# ab, wenn Eingabe oder Ausgabe kein Terminal ist (etwa unter `!` in einem Agenten-Gespraech, dessen Ausgabe im Verlauf stuende).
# Ablauf: (1) Container und Server vorbereiten (ruft 06-lauf-container.sh vorbereiten; dessen Rohdaten tragen keine URL),
#         (2) `ollama signin` im Container: URL im Browser des Projektinhabers bestaetigen,
#         (3) Rueckgabe an den Agenten. Danach liefert 06-lauf-container.sh die Messung (Hosts, Dateien, Laeufe).
# Rueckabwicklung (spaeter, 06-lauf-container.sh ende): signout im Container, Container und Volume entfernen; DANACH entfernt der Projektinhaber den
# Geraeteschluessel dieses Wegwerf-Containers in den Kontoeinstellungen von ollama.com.
# Aufruf: WORK=<messordner> tools/messung/ollama-cloud-konto/07-anmeldung.sh
set -u
if [ ! -t 0 ] || [ ! -t 1 ]; then
  echo "Abbruch: dieses Skript gehoert an ein echtes Terminal (Eingabe und Ausgabe), nicht in eine Umleitung und nicht in die Befehlszeile eines Agenten (Ausgabe stuende im Verlauf)."
  exit 2
fi
. "$(dirname "$0")/lib.sh"
vorlauf_pruefen || exit 2
B=$PFX-b

echo "Schritt 1: Container $B und Server vorbereiten (Rohdaten ohne URL unter $ROH) ..."
"$here/06-lauf-container.sh" vorbereiten > /dev/null 2>&1 || { echo "Vorbereitung fehlgeschlagen (06-lauf-container.sh vorbereiten)"; exit 1; }
docker exec "$B" true || { echo "Container $B laeuft nicht"; exit 1; }
echo
echo "Schritt 2: Anmeldung. Gleich erscheint eine URL (und gegebenenfalls ein Code)."
echo "  - Oeffne die URL in deinem Browser und bestaetige die Anmeldung dieses Wegwerf-Geraets."
echo "  - Gib URL und Code NICHT an den Agenten weiter."
echo "Weiter mit Enter."
read -r _
docker exec -it -u vscode "$B" ollama signin
echo
echo "Hast du im Browser bestaetigt? Wenn ja, Enter. (Bei Abbruch Strg-C: der Container bleibt, erneut aufrufen.)"
read -r _
echo
echo "Fertig. Sage dem Agenten: \"Anmeldung bestaetigt\". Der Container $B laeuft weiter (Wegwerf-Volume $PFX-b-modelle)."
echo "Nach der Messung: 06-lauf-container.sh ende, danach den Geraeteschluessel dieses Containers in den Kontoeinstellungen entfernen."
