#!/usr/bin/env bash
# messen-leicht.sh — leichte Messungen des Slice slice-v1-devcontainer-host-git-uebernehmen
# (LH-FA-DEV-009, LH-FA-DEV-010): Git-Identität aus Umgebungsvariablen (Messung 1),
# `gh auth token` auf dem Host (Messung 3), `git credential fill` mit und ohne Helper
# (Messung 4, Beleg für den Verzicht auf Quelle 3), `git -C <Ordner> config`-Verhalten
# (Messung 5) und die Fake-Quellen-Probe für den späteren Docker-Test (Messung 7, Probe).
#
# AUFRUF (aus der Repo-Wurzel):
#   bash tools/messung/devcontainer-host-git-uebernehmen/messen-leicht.sh \
#     > docs/messungen/slice-v1-devcontainer-host-git-uebernehmen/rohdaten-m135-leicht-lauf<L>.txt 2>&1
#
# SICHERHEIT: Der echte Token wird nirgends ausgegeben und in keine Datei geschrieben;
# als Beleg dienen nur Länge und SHA-256-Präfix der Ausgabe. Namen und E-Mails der
# Messung sind Platzhalter. Alle Verzeichnisse der Messung liegen unter $TMP (mktemp -d)
# und werden per trap EXIT/INT/TERM entfernt. Das Skript ändert sonst nichts am Host.
set -u

WIEDERHOLUNGEN="${WIEDERHOLUNGEN:-2}"
TMP="$(mktemp -d)"
LIMIT=20
export TMP

aufraeumen() { rm -rf "$TMP"; }
trap aufraeumen EXIT INT TERM

kopf() {
  echo "# Messung leicht (M1, M3, M4, M5, M7-Probe) — $(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "# Host: $(uname -srm)"
  echo "# git: $(git --version)"
  echo "# gh: $(command -v gh >/dev/null && gh --version | head -1 || echo fehlt)"
  echo "# Wiederholungen je Zelle: $WIEDERHOLUNGEN"
  echo
}

# zelle <Kennung> <Beschreibung> <Befehl-String>
zelle() {
  local id="$1" beschr="$2" befehl="$3" n
  for n in $(seq 1 "$WIEDERHOLUNGEN"); do
    echo "### Zelle $id, Wiederholung $n: $beschr"
    echo "befehl: $befehl"
    timeout -k 5 "$LIMIT" bash -c "$befehl" </dev/null >"$TMP/out" 2>"$TMP/err"
    echo "exit-status: $?"
    echo "--- stdout"
    cat "$TMP/out"
    echo "--- stderr"
    cat "$TMP/err"
    echo
  done
}

kopf

echo "## M1: Git-Identität aus GIT_AUTHOR_* / GIT_COMMITTER_* (ohne Konfigurationsdatei)"
echo "## (GIT_CONFIG_GLOBAL und GIT_CONFIG_SYSTEM auf /dev/null, Wegwerf-Repository)"
zelle M1-A "alle vier Variablen gesetzt; Commit, Autor und Committer gelesen" \
  "export GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null; T=\$(mktemp -d \"\$TMP/initXXXX\"); cd \"\$T\" && git init -q . && GIT_AUTHOR_NAME=h-name GIT_AUTHOR_EMAIL=h-email GIT_COMMITTER_NAME=c-name GIT_COMMITTER_EMAIL=c-email git commit -q --allow-empty -m m && git log --format='%an|%ae|%cn|%ce' -1"
zelle M1-B "nur GIT_AUTHOR_NAME und GIT_AUTHOR_EMAIL gesetzt" \
  "export GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null; T=\$(mktemp -d \"\$TMP/initXXXX\"); cd \"\$T\" && git init -q . && GIT_AUTHOR_NAME=h-name GIT_AUTHOR_EMAIL=h-email git commit -q --allow-empty -m m && git log --format='%an|%ae|%cn|%ce' -1"
zelle M1-C "nur GIT_AUTHOR_NAME (keine E-Mail)" \
  "export GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null; T=\$(mktemp -d \"\$TMP/initXXXX\"); cd \"\$T\" && git init -q . && GIT_AUTHOR_NAME=h-name git commit -q --allow-empty -m m && git log --format='%an|%ae|%cn|%ce' -1"
zelle M1-D "nur GIT_AUTHOR_EMAIL (kein Name)" \
  "export GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null; T=\$(mktemp -d \"\$TMP/initXXXX\"); cd \"\$T\" && git init -q . && GIT_AUTHOR_EMAIL=h-email git commit -q --allow-empty -m m && git log --format='%an|%ae|%cn|%ce' -1"
zelle M1-E "keine Variable, keine Konfiguration (negativer Referenzfall)" \
  "export GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null; T=\$(mktemp -d \"\$TMP/initXXXX\"); cd \"\$T\" && git init -q . && git commit -q --allow-empty -m m"
zelle M1-F "git var GIT_AUTHOR_IDENT mit allen vier Variablen" \
  "export GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null; T=\$(mktemp -d \"\$TMP/initXXXX\"); cd \"\$T\" && git init -q . && GIT_AUTHOR_NAME=h-name GIT_AUTHOR_EMAIL=h-email GIT_COMMITTER_NAME=c-name GIT_COMMITTER_EMAIL=c-email git var GIT_AUTHOR_IDENT; echo \"exit=\$?\""
zelle M1-G "git var GIT_AUTHOR_IDENT ohne Identität (negativer Referenzfall)" \
  "export GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null; T=\$(mktemp -d \"\$TMP/initXXXX\"); cd \"\$T\" && git init -q . && git var GIT_AUTHOR_IDENT; echo \"exit=\$?\""

echo
echo "## M3: gh auth token auf dem Host"
echo "## (Länge und SHA-256-Präfix statt des Werts; der Token selbst wird nie ausgegeben)"
zelle M3-A "unveränderte Umgebung; Token-Länge, Nichtleere, SHA-256-Präfix, Exit-Code" \
  "T=\$(mktemp -d \"\$TMP/ghXXXX\"); if O=\$(gh auth token 2>\"\$T/err\"); then L=\${#O}; P=\$(printf '%s' \"\$O\" | sha256sum | cut -c1-16); echo \"exit=0 laenge=\$L leere-ausgabe=\$([ \$L -eq 0 ] && echo ja || echo nein) sha256-praefix=\$P\"; else echo \"exit=\$? (Fehler)\"; cat \"\$T/err\"; fi"
zelle M3-B "gh nicht im PATH (PATH auf ein leeres Verzeichnis)" \
  "T=\$(mktemp -d \"\$TMP/nopathXXXX\"); mkdir \"\$T/b\"; PATH=\"\$T/b\" gh auth token; echo \"exit=\$?\""
zelle M3-C "nicht angemeldet: leere gh-Konfiguration (GH_CONFIG_DIR auf Wegwerf-Verzeichnis)" \
  "T=\$(mktemp -d \"\$TMP/nologinXXXX\"); mkdir \"\$T/cfg\"; env -u GH_TOKEN GH_CONFIG_DIR=\"\$T/cfg\" gh auth token; echo \"exit=\$?\""
zelle M3-D "mit --hostname github.com (expliziter Host)" \
  "gh auth token --hostname github.com >/dev/null 2>&1; echo \"exit=\$? (Token-Inhalt nicht gelesen)\""
zelle M3-E "GH_TOKEN in der Umgebung gesetzt (Vorgangsweise: Vorrang der Umgebungsvariable)" \
  "O=\$(GH_TOKEN=platzhalter-gh-token-0123456789 gh auth token); echo \"exit=\$? laenge=\${#O} sha256-praefix=\$(printf '%s' \"\$O\" | sha256sum | cut -c1-16)\""

echo
echo "## M4: git credential fill ohne und mit Helper (Beleg für den Verzicht auf Quelle 3)"
echo "## (GIT_TERMINAL_PROMPT=0, Zeitlimit im Zellen-Aufruf; kein Netz nötig, host.example ist ein"
echo "##  Beispiel-Host)"
CREDS_STDIN='protocol=https\nhost=host.example\n\n'
zelle M4-A "ohne konfigurierten Helper, mit GIT_TERMINAL_PROMPT=0" \
  "export GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null; export GIT_TERMINAL_PROMPT=0; unset DISPLAY GIT_ASKPASS SSH_ASKPASS; printf '$CREDS_STDIN' | git credential fill; echo \"exit=\$?\""
zelle M4-A2 "mit Zeitmessung (Sekunden, date +%sN)" \
  "export GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null; export GIT_TERMINAL_PROMPT=0; unset DISPLAY GIT_ASKPASS SSH_ASKPASS; S=\$(date +%s%N); printf '$CREDS_STDIN' | git credential fill >/dev/null 2>\"\$TMP/fillerr\"; RC=\$?; E=\$(date +%s%N); echo \"exit=\$RC dauer=\$(( (E-S)/1000000 )) ms\"; cat \"\$TMP/fillerr\""
zelle M4-B "mit Helper (store, Wegwerf-Datei, Platzhalter-Token)" \
  "T=\$(mktemp -d \"\$TMP/helperXXXX\"); printf 'https://x-access-token:platzhalter-cred-fill-0123456789@host.example\n' > \"\$T/store\"; printf 'protocol=https\nhost=host.example\n\n' | git -c credential.helper=\"store --file=\$T/store\" credential fill; echo \"exit=\$?\""
zelle M4-C "Helper liefert nichts (leere Store-Datei)" \
  "export GIT_CONFIG_SYSTEM=/dev/null; T=\$(mktemp -d \"\$TMP/helperXXXX\"); : > \"\$T/creds\"; printf 'protocol=https\nhost=host.example\n\n' | git -c credential.helper=\"store --file=\$T/creds\" credential fill; echo \"exit=\$?\""

echo
echo "## M5: git -C <Ordner> config user.name / user.email (wirksamer Wert)"
CONF_A="$TMP/syscfg"; CONF_B="$TMP/globcfg"
echo "## (GIT_CONFIG_SYSTEM=$CONF_A, GIT_CONFIG_GLOBAL=$CONF_B, Wegwerf-Repository; Werte Platzhalter)"
zelle M5-A "Wert nur in global (kein local): der wirksame Wert ist der globale" \
  "printf '[user]\n\tname = glob-name\n\temail = glob-email\n' > '$CONF_B'; : > '$CONF_A'; T=\$(mktemp -d \"\$TMP/repoXXXX\"); git -C \"\$T\" init -q .; GIT_CONFIG_GLOBAL='$CONF_B' GIT_CONFIG_SYSTEM='$CONF_A' git -C \"\$T\" config user.name; echo \"exit=\$?\"; GIT_CONFIG_GLOBAL='$CONF_B' GIT_CONFIG_SYSTEM='$CONF_A' git -C \"\$T\" config user.email; echo \"exit=\$?\""
zelle M5-B "Wert in local und global: local gewinnt" \
  "printf '[user]\n\tname = glob-name\n' > '$CONF_B'; : > '$CONF_A'; T=\$(mktemp -d \"\$TMP/repoXXXX\"); git -C \"\$T\" init -q .; git -C \"\$T\" config user.name local-name; GIT_CONFIG_GLOBAL='$CONF_B' GIT_CONFIG_SYSTEM='$CONF_A' git -C \"\$T\" config user.name; echo \"exit=\$?\""
zelle M5-C "Wert in system und global: global vor system" \
  "printf '[user]\n\tname = sys-name\n' > '$CONF_A'; printf '[user]\n\tname = glob-name\n' > '$CONF_B'; T=\$(mktemp -d \"\$TMP/repoXXXX\"); git -C \"\$T\" init -q .; GIT_CONFIG_GLOBAL='$CONF_B' GIT_CONFIG_SYSTEM='$CONF_A' git -C \"\$T\" config user.name; echo \"exit=\$?\""
zelle M5-D "Schlüssel fehlt überall: Exit-Code und leere Ausgabe" \
  ": > '$CONF_A'; : > '$CONF_B'; T=\$(mktemp -d \"\$TMP/repoXXXX\"); git -C \"\$T\" init -q .; GIT_CONFIG_GLOBAL='$CONF_B' GIT_CONFIG_SYSTEM='$CONF_A' git -C \"\$T\" config user.name; echo \"exit=\$?\""
zelle M5-E "Ordner ohne Repository (kein .git)" \
  ": > '$CONF_A'; : > '$CONF_B'; T=\$(mktemp -d \"\$TMP/plainXXXX\"); GIT_CONFIG_GLOBAL='$CONF_B' GIT_CONFIG_SYSTEM='$CONF_A' git -C \"\$T\" config user.name; echo \"exit=\$?\""
zelle M5-F "leerer Wert (user.name als Leerezeichenkette gesetzt)" \
  "printf '[user]\n\tname =\n' > '$CONF_B'; : > '$CONF_A'; T=\$(mktemp -d \"\$TMP/repoXXXX\"); git -C \"\$T\" init -q .; O=\$(GIT_CONFIG_GLOBAL='$CONF_B' GIT_CONFIG_SYSTEM='$CONF_A' git -C \"\$T\" config user.name); echo \"exit=\$? ausgabe-laenge=\${#O}\""
zelle M5-G "includeIf: global-Datei mit includeIf auf das Wegwerf-Verzeichnis" \
  ": > '$CONF_A'; printf '[user]\n\tname = glob-name\n' > '$CONF_B'; printf '[includeIf \"gitdir:%s/\"]\n\tpath = inc-cfg\n' \"$TMP\" >> '$CONF_B'; printf '[user]\n\tname = inc-name\n' > '$TMP/inc-cfg'; T=\$(mktemp -d \"\$TMP/repoXXXX\"); git -C \"\$T\" init -q .; GIT_CONFIG_GLOBAL='$CONF_B' GIT_CONFIG_SYSTEM='$CONF_A' git -C \"\$T\" config user.name; echo \"exit=\$?\""
zelle M5-H "keine Konfigurationsdatei überhaupt (negativer Referenzfall; Ordner mit Repository)" \
  "export GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null; T=\$(mktemp -d \"\$TMP/repoXXXX\"); git -C \"\$T\" init -q .; git -C \"\$T\" config user.name; echo \"exit=\$?\""

echo
echo "## M7-Probe: Fake-Quelle für den Docker-Test (Skript als gh im PATH, Wegwerf-Git-Konfiguration)"
zelle M7-A "Fake-gh (Skript) im PATH liefert den Platzhalter" \
  "T=\$(mktemp -d \"\$TMP/fakeXXXX\"); mkdir \"\$T/bin\"; printf '#!/bin/sh\necho platzhalter-fake-gh-token-0123456789\n' > \"\$T/bin/gh\"; chmod +x \"\$T/bin/gh\"; PATH=\"\$T/bin:\$PATH\" gh auth token; echo \"exit=\$?\""
zelle M7-B "Fake-gh, der nichts liefert (leere Ausgabe)" \
  "T=\$(mktemp -d \"\$TMP/fakeXXXX\"); mkdir \"\$T/bin\"; printf '#!/bin/sh\n' > \"\$T/bin/gh\"; chmod +x \"\$T/bin/gh\"; PATH=\"\$T/bin:\$PATH\" gh auth token; echo \"exit=\$?\""
zelle M7-C "Fake-gh mit Fehlerausgabe und Exit 1" \
  "T=\$(mktemp -d \"\$TMP/fakeXXXX\"); mkdir \"\$T/bin\"; printf '#!/bin/sh\necho fake-gh-error 1>&2\nexit 1\n' > \"\$T/bin/gh\"; chmod +x \"\$T/bin/gh\"; PATH=\"\$T/bin:\$PATH\" gh auth token; echo \"exit=\$?\""

echo
echo "## Kamera: Werte, die die Zellen nicht zeigen (Kontrolle)"
echo "## (Prüfung nach dem Lauf: keine der Rohdaten-Ausgaben enthält den echten Token; die"
echo "##  Skriptlogik erlaubt das nicht — der Token wird nur in einer Shell-Variable des Laufrags"
echo "##  gehalten und nur als Länge und SHA-256-Präfix belegt)"
echo "# Ende $(date -u +%Y-%m-%dT%H:%M:%SZ)"