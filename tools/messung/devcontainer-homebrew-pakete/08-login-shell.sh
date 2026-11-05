#!/bin/bash
# M7 Login-Shell: steht `brew` im PATH des Image-Benutzers in Login- und Nicht-Login-Shells, wenn der Block die Datei
# /etc/profile.d/u-boot-brew.sh anlegt? Gemessen in beiden Profilen (Standard und Sandbox; Projekte und Formeln wie 01-bau.sh, beide Images
# werden hier mit --no-cache gebaut): die Datei (Eigentuemer, Rechte, Inhalt); bash, sh und zsh als Kommando (-c), als Login-Shell (-lc) und
# interaktiv (-ic, -lic); die Sitzungen von `docker exec` (ohne Login, mit `bash -lc` und `sh -lc`) und die Login-Sitzung, die `su -l vscode`
# startet (die Art von Sitzung, die ein SSH-Dienst startet; der Dienst selbst ist nicht gemessen). Ausserdem der PATH der Login-Shell (die
# Verzeichnisse von Homebrew stehen vorn, einmal, und andere Werkzeuge bleiben auffindbar), ob die Formeln der Liste dort laufen, und die
# Kontrolle: ohne die Datei findet dieselbe Login-Sitzung das Kommando nicht (die Messung misst die Wirkung der Datei).
# Aufruf: UBOOT=<binary> WORK=<messordner> ./08-login-shell.sh      -> rohdaten-m7.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-m7.txt

TAP="devcontainer.brewTaps.allow=pt9912/u-boot"
PKG="devcontainer.features.brew.packages=ripgrep,pt9912/u-boot/u-boot"
ON="devcontainer.features.brew.enabled=true"

# pfad_messung <variante>  Misst im Image der Variante; Aussagen m7-<variante>-<modus>=<Ergebnis>.
pfad_messung() {
  local v=$1 img; img=$(vimg "$v")
  echo "-- Datei im Image $v"
  docker run --rm "$img" bash -c 'stat -c "%U:%G %a %n" /etc/profile.d/u-boot-brew.sh; cat /etc/profile.d/u-boot-brew.sh; ls /etc/profile.d/'
  aussage "m7-$v-datei=$(docker run --rm "$img" bash -c 'stat -c "%U:%G %a" /etc/profile.d/u-boot-brew.sh 2>&1')"
  echo "-- Shell-Arten (nicht interaktiv ohne TTY; -i erzwingt eine interaktive Shell)"
  local mode r
  for mode in "bash -c" "sh -c" "zsh -c" "bash -lc" "sh -lc" "zsh -lc" "bash -ic" "bash -lic"; do
    r=$(docker run --rm "$img" $mode 'command -v brew || echo MISSING' 2>&1 | tail -1)
    printf '%-10s %s\n' "$mode" "$r"
    aussage "m7-$v-$(echo "$mode" | tr ' ' '_')=$r"
  done
  echo "-- PATH der Login-Shell: bash -lc und sh -lc, Zahl der Eintraege von Homebrew, erster Eintrag, git und Formeln der Liste"
  for mode in "bash -lc" "sh -lc"; do
    docker run --rm "$img" $mode 'echo "PATH=$PATH"'
    aussage "m7-$v-$(echo "$mode" | tr ' ' '_')-eintraege-bin=$(docker run --rm "$img" $mode 'echo "$PATH" | tr ":" "\n" | grep -c "^/home/linuxbrew/.linuxbrew/bin$"')"
    aussage "m7-$v-$(echo "$mode" | tr ' ' '_')-erster=$(docker run --rm "$img" $mode 'echo "$PATH" | cut -d: -f1')"
    aussage "m7-$v-$(echo "$mode" | tr ' ' '_')-git=$(docker run --rm "$img" $mode 'command -v git || echo MISSING')"
    aussage "m7-$v-$(echo "$mode" | tr ' ' '_')-rg=$(docker run --rm "$img" $mode 'rg --version | head -1' 2>&1 | tail -1)"
    aussage "m7-$v-$(echo "$mode" | tr ' ' '_')-formeln=$(docker run --rm "$img" $mode 'brew list --formula | tr "\n" " "' 2>&1 | tail -1)"
  done
  echo "-- Sitzungen in einem laufenden Container"
  local c="$PFX-login-$v"
  docker rm -f "$c" > /dev/null 2>&1
  docker run -d --name "$c" "$img" sleep 600 > /dev/null
  r=$(docker exec "$c" bash -c 'command -v brew || echo MISSING' 2>&1 | tail -1); echo "docker exec bash -c: $r"; aussage "m7-$v-exec-bash_-c=$r"
  r=$(docker exec "$c" sh -c 'command -v brew || echo MISSING' 2>&1 | tail -1); echo "docker exec sh -c: $r"; aussage "m7-$v-exec-sh_-c=$r"
  r=$(docker exec "$c" bash -lc 'command -v brew || echo MISSING' 2>&1 | tail -1); echo "docker exec bash -lc: $r"; aussage "m7-$v-exec-bash_-lc=$r"
  r=$(docker exec "$c" sh -lc 'command -v brew || echo MISSING' 2>&1 | tail -1); echo "docker exec sh -lc: $r"; aussage "m7-$v-exec-sh_-lc=$r"
  r=$(docker exec -u root "$c" su -l vscode -c 'command -v brew || echo MISSING' 2>&1 | tail -1); echo "su -l vscode -c: $r"; aussage "m7-$v-su-l=$r"
  r=$(docker exec -u root "$c" su -l vscode -s /bin/sh -c 'command -v brew || echo MISSING' 2>&1 | tail -1); echo "su -l vscode -s /bin/sh -c: $r"; aussage "m7-$v-su-l-sh=$r"
  echo "-- Kontrolle: ohne die Datei (im Container als root entfernt) findet die Login-Sitzung das Kommando nicht"
  docker exec -u root "$c" rm /etc/profile.d/u-boot-brew.sh
  r=$(docker exec "$c" bash -lc 'command -v brew || echo MISSING' 2>&1 | tail -1); echo "docker exec bash -lc ohne Datei: $r"; aussage "m7-$v-kontrolle-bash_-lc=$r"
  r=$(docker exec -u root "$c" su -l vscode -c 'command -v brew || echo MISSING' 2>&1 | tail -1); echo "su -l vscode -c ohne Datei: $r"; aussage "m7-$v-kontrolle-su-l=$r"
  r=$(docker exec "$c" bash -c 'command -v brew || echo MISSING' 2>&1 | tail -1); echo "docker exec bash -c ohne Datei (ENV PATH): $r"; aussage "m7-$v-kontrolle-bash_-c=$r"
  docker rm -f "$c" > /dev/null
}

echo "== M7 Login-Shell, $(now)"
echo "host: $(uname -m) docker=$(docker version --format '{{.Server.Version}} {{.Server.Os}}/{{.Server.Arch}}')"
echo "u-boot: $("$UBOOT" --version 2>&1 | head -1)"

mark "Projekte"
projekt def 0 "$TAP" "$ON" "$PKG"
projekt sbx 1 "$TAP" "$ON" "$PKG" "devcontainer.sandbox.egress.enabled=true"
for v in def sbx; do
  echo "-- Dockerfile $v: $(grep -c 'u-boot-brew.sh' "$(vdir $v)/.devcontainer/Dockerfile") Zeile(n) mit u-boot-brew.sh"
  grep -n 'u-boot-brew.sh' "$(vdir $v)/.devcontainer/Dockerfile"
done

for v in def sbx; do
  mark "Bau $v"
  if bild "$v" --no-cache; then
    aussage "m7-$v-bau=ok"
    aussage "m7-$v-revision=$(docker run --rm "$(vimg "$v")" git -C /home/linuxbrew/.linuxbrew/Homebrew rev-parse HEAD)"
    mark "Messung $v"
    pfad_messung "$v"
  else
    aussage "m7-$v-bau=FEHLER"
  fi
done
echo "== ende $(now)"
