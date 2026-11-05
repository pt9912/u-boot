#!/bin/bash
# M1 Bau: das erzeugte Dockerfile mit dem Katalog-Schluessel `brew` (Formeln ripgrep und die Tap-Formel pt9912/u-boot/u-boot) baut im
# Standardprofil und im Sandbox-Profil (mit eingeschalteter Egress-Restriktion ohne zusaetzliche Eintraege; die Restriktion greift beim
# Containerstart, nicht beim Bau). Gemessen: Dauer, Image-Groesse gegen das Basis-Image, Versionen im Image, Benutzer, Revision des
# Homebrew-Verzeichnisses gegen den Pin, Homebrew allein (leere Liste), und die Negativprobe: ein Pin, den es nicht gibt, bricht den Bau ab.
# Alle Baeue laufen mit --no-cache (die Dauer ist die eines kalten Baus; das Basis-Image liegt lokal).
# Aufruf: UBOOT=<binary> WORK=<messordner> ./01-bau.sh      -> rohdaten-m1.txt
set -u
. "$(dirname "$0")/lib.sh"
rohdaten rohdaten-m1.txt

echo "== M1 Bau, $(now)"
echo "host: $(uname -m) docker=$(docker version --format '{{.Server.Version}} {{.Server.Os}}/{{.Server.Arch}}')"
echo "base-image bytes: $(docker image inspect --format '{{.Size}}' "$BASE" 2>/dev/null || echo nicht-lokal)"
echo "u-boot: $("$UBOOT" --version 2>&1 | head -1)"

TAP="devcontainer.brewTaps.allow=pt9912/u-boot"
PKG="devcontainer.features.brew.packages=ripgrep,pt9912/u-boot/u-boot"
ON="devcontainer.features.brew.enabled=true"

mark "Projekte"
projekt def 0 "$TAP" "$ON" "$PKG"
projekt sbx 1 "$TAP" "$ON" "$PKG" "devcontainer.sandbox.egress.enabled=true"
projekt alone 0 "$ON"
projekt egr 1 "$ON" "devcontainer.sandbox.egress.enabled=true"

for v in def sbx alone egr; do
  echo "-- Dockerfile $v: Zeilen=$(wc -l < "$(vdir $v)/.devcontainer/Dockerfile") HOMEBREW_REVISION=$(grep -m1 '^ARG HOMEBREW_REVISION=' "$(vdir $v)/.devcontainer/Dockerfile")"
done
echo "-- brew-Zeile im Dockerfile def:"; grep -n "brew install" "$(vdir def)/.devcontainer/Dockerfile"

probe='set -u
echo "user=$(id -un) uid=$(id -u)"
echo "brew-pfad=$(command -v brew)"
echo "brew-version=$(brew --version 2>&1 | head -1)"
echo "revision=$(git -C /home/linuxbrew/.linuxbrew/Homebrew rev-parse HEAD)"
echo "eigentuemer=$(stat -c %U:%G /home/linuxbrew/.linuxbrew /home/linuxbrew/.linuxbrew/Homebrew)"
echo "formeln=$(brew list --formula 2>&1 | tr "\n" " ")"
echo "rg=$(rg --version 2>&1 | head -1)"
echo "u-boot-aus-tap=$(command -v u-boot) $(u-boot --version 2>&1 | head -1)"
echo "git-origin=$(git -C /home/linuxbrew/.linuxbrew/Homebrew remote get-url origin)"
echo "shallow=$(git -C /home/linuxbrew/.linuxbrew/Homebrew rev-parse --is-shallow-repository)"
echo "env-homebrew=$(env | grep -c "^HOMEBREW")"'

for v in def sbx; do
  mark "Bau $v"
  if bild "$v" --no-cache; then
    docker run --rm "$(vimg "$v")" bash -c "$probe"
    aussage "m1-$v-bau=ok"
  else
    aussage "m1-$v-bau=FEHLER"
  fi
done

mark "Homebrew allein (leere Liste)"
if bild alone --no-cache; then
  docker run --rm "$(vimg alone)" bash -c 'echo "brew-version=$(brew --version 2>&1 | head -1)"; echo "formeln=$(brew list --formula 2>&1 | tr "\n" " ")"; echo "revision=$(git -C /home/linuxbrew/.linuxbrew/Homebrew rev-parse HEAD)"'
  aussage "m1-alone-bau=ok"
fi

mark "Images fuer 04-egress.sh (Sandbox-Profil, Restriktion an, Homebrew allein)"
bild egr --no-cache

mark "Negativprobe: ein Pin, den es nicht gibt, bricht den Bau ab (kein Image)"
if bild def --no-cache --build-arg HOMEBREW_REVISION=0000000000000000000000000000000000000001; then
  aussage "m1-pin-falsch=BAU-GELUNGEN(unerwartet)"
else
  aussage "m1-pin-falsch=bau-abgebrochen"
fi
echo "== ende $(now)"
