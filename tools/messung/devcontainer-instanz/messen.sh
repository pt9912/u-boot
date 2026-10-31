#!/usr/bin/env bash
# messen.sh — Messung des Slice slice-v1-devcontainer-instanz-messung (LH-FA-DEV-010).
#
# FRAGE: Was tut die Dev-Container-CLI des Hilfs-Images, wenn N = 2, 3 oder 4 gleichzeitige
# Instanzen desselben Projektordners über ein zweites Id-Label (`--id-label u-boot.instance=<Name>`)
# unterschieden werden? Eigene Container, eigene Volumes (Workspace, Docker-Speicher, Modelle) oder
# geteilte? Was tut `up` ohne Flag neben einer Instanz? Welche Namen taugen als Label-Wert?
#
# AUFRUF (aus der Repo-Wurzel; Docker-Client + lokale Engine, das u-boot-Binary, Netz):
#   bash tools/messung/devcontainer-instanz/messen.sh <2|3|4|namen> bin/u-boot-linux-amd64 \
#     > docs/messungen/slice-v1-devcontainer-instanz-messung/rohdaten-stufe<N>-lauf<L>.txt 2>&1
#
#   Stufe 2|3|4: M0 bis M7 und M9 mit N gleichzeitigen Instanzen (Sandbox-Profil mit nestedRuntime
#                docker und Feature ollama; danach Standardprofil; bei N = 2 zusätzlich ein Lauf
#                ohne Limit als Bezug).
#   namen:       M8 (Label-Werte) im Standardprofil.
#
# MESSBEDINGUNG je Instanz: runArgs der devcontainer.json im Wegwerf-Projektordner um
# --memory=3g --memory-swap=3g --cpus=1 ergänzt (kein weiteres `generate` danach); die Wirkung wird
# am Container zurückgelesen (HostConfig). Der Aufruf der CLI spiegelt `runArgs` in
# internal/adapter/driven/docker/devcontainer_launcher.go (docker run --rm --name … --mount Socket
# --mount Projektordner -e GIT_TOKEN <Image> devcontainer up --workspace-folder <Ordner> …), hängt
# aber `--id-label` an; das Skript fährt NICHT `u-boot devcontainer up`.
#
# Lockerung (ausgewiesen): der Hilfs-Container bekommt den Docker-Socket des Hosts.
# Alles, was das Skript anlegt, trägt das Präfix ubmess-/vsc-ubmess-; nie ein Glob über fremde
# Namen. Vorprüfung: Abbruch mit Befund, nichts wird am Host geändert (kein sysctl).
# Aufräumen per trap EXIT/INT/TERM in der Reihenfolge Hilfs-Container, Messcontainer, Volumes,
# Images, Projektordner; danach Inventur nachher und Differenz zur Inventur vorher (M9).
set -u

stufe=${1:?usage: messen.sh <2|3|4|namen> <u-boot-binary>}
bin=$(readlink -f "${2:?usage: messen.sh <2|3|4|namen> <u-boot-binary>}")
here=$(cd "$(dirname "$0")" && pwd)
ctx="$here/../../../internal/adapter/driven/docker/devcontainercli"

HELPER_IMG=ubmess-devcontainer-cli:0.89.0
export GIT_TOKEN="platzhalter-kein-zugangsdatum-0123456789abcdef"

base=$(mktemp -d)
declare -A VN CID DUR
UP_LIMIT=300
STAGE_LIMIT=900
DEADLINE=0
RELAXED=0
MON_PID=""
FOREIGN=()
HSEQ=0
PSEQ=0
VOLS=""
DIRS=()
LAST_RC=0 LAST_DUR=0 LAST_CID="" PDIR="" PNAME=""
STATUS="vollstaendig"

section() { printf '\n######## %s\n' "$*"; }
now() { date +%s; }
memavail_kb() { awk '/^MemAvailable:/{print $2}' /proc/meminfo; }
hostmem() {
  echo "  [host] MemAvailable $(( $(memavail_kb) / 1024 )) MiB; loadavg $(cut -d' ' -f1-3 /proc/loadavg); $(free -m | awk '/^Mem:/{print "Mem used " $3 " MiB, free " $4 " MiB"} /^Swap:/{print "Swap used " $3 " MiB"}' | tr '\n' ' ')"
}

inventur() {
  docker ps -a --format '{{.Names}}	{{.Image}}' | sort > "$1.ps"
  docker volume ls -q | sort > "$1.vol"
  docker images --format '{{.Repository}}:{{.Tag}}	{{.ID}}' | sort > "$1.img"
}

# ---- Aufräumen ----------------------------------------------------------------------------

teardown_dir() {
  local dir=$1 ids
  ids=$(docker ps -aq --filter "label=devcontainer.local_folder=$dir")
  for id in $ids; do track_volumes "$id"; done
  # shellcheck disable=SC2086
  [ -n "$ids" ] && docker rm -f -v $ids > /dev/null 2>&1
}

cleanup() {
  trap - EXIT INT TERM
  [ -n "$MON_PID" ] && { kill "$MON_PID" 2> /dev/null; wait "$MON_PID" 2> /dev/null; MON_PID=""; monitor_summary; }
  section "Aufraeumen (trap)"
  local helpers v imgs
  helpers=$(docker ps -aq --filter "name=ubmess-helper-")
  # shellcheck disable=SC2086
  [ -n "$helpers" ] && docker rm -f $helpers > /dev/null 2>&1
  for d in "${DIRS[@]:-}"; do [ -n "$d" ] && teardown_dir "$d"; done
  for v in $VOLS; do
    case $v in ubmess-*) docker volume rm -f "$v" > /dev/null 2>&1 ;; esac
  done
  imgs=$(docker images -q --filter 'reference=vsc-ubmess-*')
  # shellcheck disable=SC2086
  [ -n "$imgs" ] && docker rmi -f $imgs > /dev/null 2>&1
  docker rmi -f "$HELPER_IMG" > /dev/null 2>&1
  rm -rf "$base" 2> /dev/null
  if [ -e "$base" ]; then
    echo "rm -rf $base scheiterte (Dateien anderer Besitzer); Fallback: Entfernen im Container"
    docker run --rm -v "$base":/d alpine:3 sh -c 'rm -rf /d/* /d/.[!.]* 2>/dev/null; true' > /dev/null 2>&1
    rm -rf "$base"
  fi
  [ -e "$base" ] && echo "BEFUND: Arbeitsordner $base bleibt stehen" || echo "Arbeitsordner entfernt"
  section "M9: Inventur nachher und Differenz"
  if [ -f "$INVDIR/vorher.ps" ]; then
    inventur "$INVDIR/nachher"
    for k in ps vol img; do
      echo "--- Differenz $k (vorher -> nachher; leer = keine)"
      diff "$INVDIR/vorher.$k" "$INVDIR/nachher.$k"
    done
    echo "--- Reste mit Praefix ubmess-/vsc-ubmess-:"
    { docker ps -a --format '{{.Names}}' ; docker volume ls -q; docker images --format '{{.Repository}}'; } | grep -E '^(ubmess-|vsc-ubmess-)' || echo "keine"
    echo "--- Platz nachher: $(df -h --output=avail / | tail -n 1 | tr -d ' ') frei auf / (Docker-Wurzel liegt auf /)"
  fi
  rm -rf "$INVDIR" 2> /dev/null
  section "Ende"
  echo "Stand der Stufe: $STATUS"
}
INVDIR=$(mktemp -d)
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

# Ueberwachung des Hosts (nur lesen): schreibt $base/mon.log, setzt $base/abort bei Abbruchbedingung.
monitor() {
  local hi=0 sf ma si so id
  while :; do
    sf=$(awk '/^SwapFree:/{print $2}' /proc/meminfo)
    ma=$(memavail_kb)
    read -r si so < <(vmstat 1 2 | tail -n 1 | awk '{print $7, $8}')
    echo "$(date +%T) MemAvailable_kB=$ma SwapFree_kB=$sf si=${si:-0} so=${so:-0}" >> "$base/mon.log"
    if [ $(( ${si:-0} + ${so:-0} )) -ge 5000 ]; then hi=$((hi + 1)); else hi=0; fi
    [ "$hi" -ge 5 ] && echo "si+so >= 5000 KiB/s in 5 Proben" > "$base/abort"
    for id in "${FOREIGN[@]:-}"; do
      [ -n "$id" ] || continue
      if [ "$(docker inspect --format '{{.State.Running}} {{.State.OOMKilled}}' "$id" 2>&1)" != "true false" ]; then
        echo "vorhandener Container ${id:0:12} laeuft nicht mehr oder OOMKilled" > "$base/abort"
      fi
    done
    sleep 1
  done
}

monitor_summary() {
  [ -f "$base/mon.log" ] || return 0
  section "Ueberwachung des Hosts (Zusammenfassung)"
  echo "Proben: $(wc -l < "$base/mon.log"); erste: $(head -n 1 "$base/mon.log"); letzte: $(tail -n 1 "$base/mon.log")"
  awk '{split($2,a,"=");split($3,b,"=");split($4,c,"=");split($5,d,"=");
        if(NR==1||a[2]<mi)mi=a[2]; if(NR==1){s0=b[2];smin=b[2]} if(b[2]<smin)smin=b[2]; if(c[2]+d[2]>mx)mx=c[2]+d[2]; if(c[2]+d[2]>=1)n++}
        END{printf "min MemAvailable %d MiB; SwapFree am Anfang %d MiB, minimal %d MiB (Swap-Wachstum %d MiB); max si+so %d KiB/s; Proben mit si+so>0: %d\n", mi/1024, s0/1024, smin/1024, (s0-smin)/1024, mx, n}' "$base/mon.log"
  echo "Swap-Druck (Host-Swap gewachsen um mehr als 256 MiB oder si/so dauerhaft): $(awk '{split($3,b,"=");if(NR==1)s0=b[2];if(b[2]<m||NR==1)m=b[2]}END{print ((s0-m)/1024>256)?"JA":"nein"}' "$base/mon.log")"
  [ -f "$base/abort" ] && echo "Abbruchbedingung ausgeloest: $(cat "$base/abort")"
  echo "dmesg (OOM-Zeilen, nur lesen): $(dmesg 2>&1 | grep -ci 'out of memory\|oom-kill' || true)"
  echo "Proben-Auszug (jede 10.):"; awk 'NR%10==1' "$base/mon.log"
}

# ---- Helfer -------------------------------------------------------------------------------

track_volumes() {
  local v
  for v in $(docker inspect --format '{{range .Mounts}}{{if eq .Type "volume"}}{{.Name}} {{end}}{{end}}' "$1" 2> /dev/null); do
    case " $VOLS " in *" $v "*) ;; *) VOLS="$VOLS $v" ;; esac
  done
}

# up <ordner> [CLI-Optionen…]: ein `devcontainer up` im Hilfs-Container, gemessen.
up() {
  local dir=$1 hn t0
  shift
  HSEQ=$((HSEQ + 1))
  hn="ubmess-helper-$HSEQ"
  t0=$(now)
  if [ -f "$base/abort" ]; then
    echo "BEFUND: Abbruch durch die Ueberwachung ($(cat "$base/abort")); kein weiterer up"
    STATUS="Abbruch durch Ueberwachung: $(cat "$base/abort")"
    return 1
  fi
  if [ "$t0" -ge "$DEADLINE" ]; then
    echo "BEFUND: Zeitlimit des Laufs (${STAGE_LIMIT} s) erreicht; kein weiterer up"
    STATUS="Abbruch durch Zeitlimit des Laufs (${STAGE_LIMIT} s)"
    return 1
  fi
  timeout -k 10 "$UP_LIMIT" docker run --rm --name "$hn" \
    --mount "type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock" \
    --mount "type=bind,source=$dir,target=$dir" -e GIT_TOKEN \
    "$HELPER_IMG" devcontainer up --workspace-folder "$dir" "$@" > "$base/$hn.out" 2> "$base/$hn.err"
  LAST_RC=$?
  LAST_DUR=$(( $(now) - t0 ))
  if [ "$LAST_RC" -eq 124 ] || [ "$LAST_RC" -eq 137 ]; then
    echo "BEFUND: Zeitlimit (${UP_LIMIT} s) fuer devcontainer up [$*] erreicht; Hilfs-Container $hn wird entfernt. Letzte Zeilen der CLI:"
    echo "  Prozesse im Hilfs-Container (docker top):"; docker top "$hn" 2>&1 | cut -c1-200 | sed 's/^/    /'
    docker rm -f "$hn" > /dev/null 2>&1
    tail -n 4 "$base/$hn.err" | cut -c1-200 | sed 's/^/    /'
    STATUS="Zeitlimit bei up [$*]"
    return 1
  fi
  LAST_CID=$(tail -n 1 "$base/$hn.out" | grep -o '"containerId":"[0-9a-f]*"' | cut -d'"' -f4)
  echo "up $hn [$*] -> exit=$LAST_RC dauer=${LAST_DUR} s containerId=${LAST_CID:0:12}"
  echo "  Ergebnis: $(tail -n 1 "$base/$hn.out" | cut -c1-260)"
  if [ "$LAST_RC" -ne 0 ] || [ -z "$LAST_CID" ]; then
    echo "  stderr (letzte 10 Zeilen):"
    tail -n 10 "$base/$hn.err" | cut -c1-300 | sed 's/^/    /'
    return 1
  fi
  [ -n "$LAST_CID" ] && track_volumes "$LAST_CID"
  return 0
}

wait_dockerd() { # <cid>: Sekunden bis der verschachtelte Docker antwortet
  local t0 i
  t0=$(now)
  for i in $(seq 1 45); do
    if docker exec -u root "$1" docker info > /dev/null 2>&1; then echo "$(( $(now) - t0 ))"; return 0; fi
    sleep 2
  done
  echo "nein nach $(( $(now) - t0 ))"
  return 1
}

limits() { # <cid>: zurückgelesene Grenze und Zustand
  docker inspect --format 'Memory={{.HostConfig.Memory}} MemorySwap={{.HostConfig.MemorySwap}} NanoCpus={{.HostConfig.NanoCpus}} Privileged={{.HostConfig.Privileged}} OOMKilled={{.State.OOMKilled}} Running={{.State.Running}} ExitCode={{.State.ExitCode}}' "$1"
}

new_project() { # <sandbox|standard> <limit|nolimit> -> PDIR, PNAME
  local profile=$1 lim=$2
  PSEQ=$((PSEQ + 1))
  PNAME="ubmess-$PSEQ"
  PDIR="$base/$PNAME"
  mkdir -p "$PDIR"
  DIRS+=("$PDIR")
  if [ "$profile" = sandbox ]; then
    ( cd "$PDIR" && "$bin" init --devcontainer --sandbox --no-git --yes > /dev/null 2>&1 )
    cat > "$PDIR/u-boot.yaml" <<EOF
schemaVersion: 1
project:
    name: $PNAME
devcontainer:
    enabled: true
    profile: sandbox
    sandbox:
        nestedRuntime: docker
    features:
        ollama:
            enabled: true
EOF
    ( cd "$PDIR" && "$bin" generate devcontainer --yes > /dev/null 2>&1 )
  else
    ( cd "$PDIR" && "$bin" init --devcontainer --no-git --yes > /dev/null 2>&1 )
  fi
  local f="$PDIR/.devcontainer/devcontainer.json"
  if [ "$lim" = limit ]; then
    if grep -q '"runArgs"' "$f"; then
      sed -i 's/"runArgs": \[/"runArgs": ["--memory=3g", "--memory-swap=3g", "--cpus=1", /' "$f"
    else
      sed -i 's/^  "name": \(.*\),$/  "name": \1,\n  "runArgs": ["--memory=3g", "--memory-swap=3g", "--cpus=1"],/' "$f"
    fi
  fi
  echo "Projekt $PNAME ($profile, $lim) in $PDIR"
  grep -n 'runArgs\|"mounts"\|source=' "$f" | sed 's/^/  json: /'
}

labels_of() { # <cid>
  docker inspect --format '  devcontainer.local_folder={{index .Config.Labels "devcontainer.local_folder"}}
  devcontainer.config_file={{index .Config.Labels "devcontainer.config_file"}}
  u-boot.instance={{index .Config.Labels "u-boot.instance"}}
  Zahl der Labels insgesamt: {{len .Config.Labels}}
  Image={{.Image}} Config.Image={{.Config.Image}}
  Einhaengungen mit docker.sock: {{range .Mounts}}{{if eq .Destination "/var/run/docker.sock"}}X{{end}}{{end}}' "$1"
  echo "  Label-Namen: $(docker inspect --format '{{range $k,$v := .Config.Labels}}{{$k}} {{end}}' "$1")"
}

fam_of() { # <ziel>
  case $1 in
    /workspaces/*) echo workspace ;;
    /var/lib/docker) echo docker ;;
    /home/vscode/.ollama) echo ollama ;;
    */.local/share/containers) echo containers ;;
    *) echo "sonst:$1" ;;
  esac
}

fill_vn() { # <k> <cid>: VN[k,familie]=Volume-Name
  local line name dest
  while IFS='|' read -r name dest; do
    [ -n "$name" ] && VN["$1,$(fam_of "$dest")"]=$name
  done < <(docker inspect --format '{{range .Mounts}}{{if eq .Type "volume"}}{{.Name}}|{{.Destination}}{{"\n"}}{{end}}{{end}}' "$2")
}

dcid_of() { # <cid>: aufgeloester ${devcontainerId} aus dem Workspace-Volume
  local n
  n=$(docker inspect --format '{{range .Mounts}}{{if eq .Type "volume"}}{{if eq .Destination "/workspaces/'"$PNAME"'"}}{{.Name}}{{end}}{{end}}{{end}}' "$1")
  echo "${n#"$PNAME"-workspace-}"
}

df_free() { df -h --output=avail / | tail -n 1 | tr -d ' '; }

# ---- Phase: Sandbox-Profil ------------------------------------------------------------------

phase_sandbox() { # <N> <limit|nolimit> <Kennung>
  local N=$1 lim=$2 tag=$3 k i j fam ok
  section "[$tag] Sandbox-Profil, $N Instanzen, $lim"
  new_project sandbox "$lim"
  local dir=$PDIR name=$PNAME
  CID=() DUR=() VN=()
  echo "Platz frei vorher: $(df_free); docker volume ls (ubmess-): $(docker volume ls -q | grep -c '^ubmess-')"
  docker system df 2>&1 | sed 's/^/  df: /'
  section "[$tag] M1/M4/M5: Start der Instanzen nacheinander"
  for k in $(seq 1 "$N"); do
    echo "-- Instanz $k"
    hostmem
    if ! up "$dir" --id-label "devcontainer.local_folder=$dir" --id-label "u-boot.instance=i$k"; then
      echo "BEFUND: Start der Instanz $k von $N scheitert (Stufe bricht ab; Messwerte bis hierher gelten)"
      STATUS="Abbruch bei Instanz $k von $N (Lauf $tag)"
      hostmem; docker ps -a --filter "label=devcontainer.local_folder=$dir" --format '  {{.ID}} {{.Status}}'
      return 1
    fi
    CID[$k]=$LAST_CID
    DUR[$k]=$LAST_DUR
    echo "  verschachtelter Docker antwortet nach: $(wait_dockerd "${CID[$k]}") s"
    echo "  Images vsc-ubmess-* auf dem Host: $(docker images --filter 'reference=vsc-ubmess-*' --format '{{.Repository}}:{{.Tag}}={{.ID}}' | tr '\n' ' ')"
    echo "  Grenze zurueckgelesen: $(limits "${CID[$k]}")"
    hostmem
  done
  echo "Dauer der Starts (s): $(for k in $(seq 1 "$N"); do printf 'i%s=%s ' "$k" "${DUR[$k]}"; done)"

  section "[$tag] M1: Container und Labels"
  for k in $(seq 1 "$N"); do echo "Instanz $k: ${CID[$k]}"; labels_of "${CID[$k]}"; done
  echo "Container-IDs paarweise verschieden: $(printf '%s\n' "${CID[@]}" | sort -u | wc -l) von $N"
  echo "docker ps -a --filter label=devcontainer.local_folder=<Ordner>:"
  docker ps -a --filter "label=devcontainer.local_folder=$dir" --format '  {{.ID}} instance={{.Label "u-boot.instance"}} {{.Status}}'

  section "[$tag] M2/M3: Volume-Namen (Instanz x Familie) und aufgeloestes devcontainerId"
  for k in $(seq 1 "$N"); do fill_vn "$k" "${CID[$k]}"; done
  for k in $(seq 1 "$N"); do
    echo "Instanz $k: devcontainerId=$(dcid_of "${CID[$k]}")"
    for fam in workspace docker ollama containers; do echo "    $fam -> ${VN[$k,$fam]:-<keine>}"; done
  done
  for fam in workspace docker ollama; do
    local all n
    all=$(for k in $(seq 1 "$N"); do echo "${VN[$k,$fam]:-}"; done)
    n=$(echo "$all" | sort -u | grep -c .)
    echo "Familie $fam: $n verschiedene Namen bei $N Instanzen -> paarweise verschieden: $([ "$n" -eq "$N" ] && echo ja || echo NEIN)"
  done
  echo "devcontainerId paarweise verschieden: $(for k in $(seq 1 "$N"); do dcid_of "${CID[$k]}"; done | sort -u | wc -l) von $N"
  echo "docker volume ls (ubmess-) nach dem Start:"
  docker volume ls -q | grep '^ubmess-' | sed 's/^/  /'
  echo "Zahl der Volumes mit Praefix ubmess-$PSEQ-: $(docker volume ls -q | grep -c "^$name-")  (erwartet: 3 Familien x $N = $((3 * N)))"

  section "[$tag] M3: Inhaltsprobe ueber alle Paare"
  for i in $(seq 1 "$N"); do
    docker exec -u root "${CID[$i]}" sh -c "echo ws-$i > /workspaces/$name/marker-$i && echo ol-$i > /home/vscode/.ollama/marker-$i" && echo "Instanz $i: Marker in Workspace und Modell-Volume geschrieben"
    docker exec -u root "${CID[$i]}" sh -c "docker pull -q busybox:1.36 > /dev/null && docker tag busybox:1.36 ubmess-probe-$i:1" > "$base/pull-$i.log" 2>&1 && echo "Instanz $i: Image ubmess-probe-$i:1 im verschachtelten Docker" || { echo "Instanz $i: pull/tag scheiterte:"; tail -n 3 "$base/pull-$i.log"; }
  done
  local leaks=0
  for i in $(seq 1 "$N"); do
    for j in $(seq 1 "$N"); do
      [ "$i" = "$j" ] && continue
      local w o m
      w=$(docker exec -u root "${CID[$j]}" sh -c "test -e /workspaces/$name/marker-$i && echo PRESENT || echo absent")
      o=$(docker exec -u root "${CID[$j]}" sh -c "test -e /home/vscode/.ollama/marker-$i && echo PRESENT || echo absent")
      m=$(docker exec -u root "${CID[$j]}" sh -c "docker image inspect ubmess-probe-$i:1 > /dev/null 2>&1 && echo PRESENT || echo absent")
      echo "Datei/Image von Instanz $i in Instanz $j: workspace=$w ollama=$o docker-image=$m"
      case "$w$o$m" in *PRESENT*) leaks=$((leaks + 1)) ;; esac
    done
  done
  for i in $(seq 1 "$N"); do
    echo "Eigene Daten Instanz $i sichtbar: workspace=$(docker exec -u root "${CID[$i]}" cat "/workspaces/$name/marker-$i" 2>&1) ollama=$(docker exec -u root "${CID[$i]}" cat "/home/vscode/.ollama/marker-$i" 2>&1) image=$(docker exec -u root "${CID[$i]}" sh -c "docker image inspect -f '{{.Id}}' ubmess-probe-$i:1 2>&1 | cut -c1-19")"
  done
  echo "Paare mit sichtbarem Fremdinhalt: $leaks"

  section "[$tag] M5: Ressourcen bei $N laufenden Instanzen"
  hostmem
  docker stats --no-stream --format '  {{.ID}} mem={{.MemUsage}} ({{.MemPerc}}) cpu={{.CPUPerc}} pids={{.PIDs}}' $(for k in $(seq 1 "$N"); do echo "${CID[$k]}"; done)
  for k in $(seq 1 "$N"); do echo "Instanz $k: $(limits "${CID[$k]}")"; done
  echo "Platz frei jetzt: $(df_free)"
  docker system df 2>&1 | sed 's/^/  df: /'
  echo "Groesse der neuen Volumes und Images (docker system df -v):"
  docker system df -v 2>&1 | grep -E 'ubmess|^VOLUME NAME|^REPOSITORY' | sed 's/^/  /'

  if [ "$lim" = limit ]; then
    section "[$tag] M6: up ohne Flag neben $N Instanzen"
    local avail need
    avail=$(memavail_kb)
    need=$((3 * 1048576 + 1572864))
    if [ "$avail" -ge "$need" ]; then
      echo "MemAvailable $avail kB >= $need kB (3 GiB + 1,5 GiB): M6 laeuft"
      if up "$dir"; then
        local newc=$LAST_CID
        echo "Neuer Container (ohne Flag): ${newc:0:12}; gleich einem der Instanz-Container: $(printf '%s\n' "${CID[@]}" | grep -c "^$newc$")"
        labels_of "$newc"
        echo "devcontainerId des Containers ohne Flag: $(dcid_of "$newc"); gleich dem einer Instanz: $(for k in $(seq 1 "$N"); do dcid_of "${CID[$k]}"; done | grep -c "^$(dcid_of "$newc")$")"
        echo "Volumes des Containers ohne Flag:"; docker inspect --format '{{range .Mounts}}{{if eq .Type "volume"}}    {{.Name}}{{"\n"}}{{end}}{{end}}' "$newc"
        echo "docker ps -a --filter label=devcontainer.local_folder=<Ordner> trifft:"
        docker ps -a --filter "label=devcontainer.local_folder=$dir" --format '  {{.ID}} instance={{.Label "u-boot.instance"}} {{.Status}}'
        echo "Treffer: $(docker ps -aq --filter "label=devcontainer.local_folder=$dir" | wc -l) (Instanzen $N + ohne Flag 1 = $((N + 1)) erwartet bei eigenem Container)"
        echo "FindContainer-Form (erster Treffer): $(docker ps -a --no-trunc --filter "label=devcontainer.local_folder=$dir" --format '{{.ID}}' | head -n 1 | cut -c1-12)"
        # Nur ein vom Aufruf ohne Flag NEU angelegter Container wird entfernt; ist es ein
        # Container der Instanzen, bleibt er (er gehoert M7).
        if printf '%s\n' "${CID[@]}" | grep -q "^$newc$"; then
          echo "Aufruf ohne Flag hat den Container einer Instanz uebernommen (Instanz $(docker inspect --format '{{index .Config.Labels "u-boot.instance"}}' "$newc")); nichts entfernt"
        else
          docker rm -f -v "$newc" > /dev/null 2>&1 && echo "neu angelegten Container ohne Flag entfernt"
        fi
      fi
    else
      echo "M6 vorwaerts: NICHT GEMESSEN (Host-Schutz): MemAvailable $avail kB < $need kB (3 GiB + 1,5 GiB fuer einen weiteren Container)"
    fi
  fi

  section "[$tag] M7: --remove-existing-container, docker stop"
  echo "vor: $(for k in $(seq 1 "$N"); do printf 'i%s=%s ' "$k" "${CID[$k]:0:12}"; done)"
  local old=${CID[1]}
  if up "$dir" --remove-existing-container --id-label "devcontainer.local_folder=$dir" --id-label "u-boot.instance=i1"; then
    CID[1]=$LAST_CID
    echo "Instanz 1 neu: ${CID[1]:0:12}; alter Container noch da: $(docker inspect --type container "$old" > /dev/null 2>&1 && echo ja || echo nein)"
    for k in $(seq 2 "$N"); do echo "Instanz $k noch derselbe Container und laeuft: $(docker inspect --format '{{.Id}} {{.State.Running}}' "${CID[$k]}" | cut -c1-12,65-)"; done
    echo "devcontainerId Instanz 1 nach --remove-existing-container: $(dcid_of "${CID[1]}") (vorher $(echo "${VN[1,workspace]#"$name"-workspace-}"))"
    wait_dockerd "${CID[1]}" > /dev/null
    echo "Daten von Instanz 1 nach Neuanlage: workspace=$(docker exec -u root "${CID[1]}" cat "/workspaces/$name/marker-1" 2>&1) ollama=$(docker exec -u root "${CID[1]}" cat "/home/vscode/.ollama/marker-1" 2>&1) image=$(docker exec -u root "${CID[1]}" sh -c "docker image inspect -f '{{.Id}}' ubmess-probe-1:1 2>&1 | cut -c1-19")"
  fi
  docker stop "${CID[2]}" > /dev/null 2>&1
  echo "nach docker stop der Instanz 2:"
  docker ps -a --filter "label=devcontainer.local_folder=$dir" --format '  {{.ID}} instance={{.Label "u-boot.instance"}} {{.Status}}'
  echo "Volumes aller Instanzen noch vorhanden: $(for k in $(seq 1 "$N"); do for fam in workspace docker ollama; do docker volume inspect "${VN[$k,$fam]}" > /dev/null 2>&1 && echo ok || echo FEHLT; done; done | sort | uniq -c | tr '\n' ' ')"

  teardown_dir "$dir"
  if [ "$lim" = limit ]; then phase_bezug "$dir" "$tag"; fi
  return 0
}

# M2: Bezug ohne Flag und mit nur dem Label devcontainer.local_folder (Projekt ohne Instanzen).
phase_bezug() {
  local dir=$1 tag=$2 a b
  section "[$tag] M2 Bezug: ohne Flag und nur mit devcontainer.local_folder"
  if up "$dir"; then
    a=$(dcid_of "$LAST_CID"); echo "devcontainerId ohne Flag: $a"; labels_of "$LAST_CID"
    docker rm -f -v "$LAST_CID" > /dev/null 2>&1
  fi
  if up "$dir" --id-label "devcontainer.local_folder=$dir"; then
    b=$(dcid_of "$LAST_CID"); echo "devcontainerId mit --id-label devcontainer.local_folder: $b"; labels_of "$LAST_CID"
    docker rm -f -v "$LAST_CID" > /dev/null 2>&1
  fi
  echo "Gleich: $([ -n "${a:-}" ] && [ "${a:-}" = "${b:-}" ] && echo ja || echo nein)"
  echo "Gleich dem Wert einer Instanz (i1): $([ "${a:-x}" = "$(echo "${VN[1,workspace]:-}" | sed "s/^$PNAME-workspace-//")" ] && echo ja || echo nein)"
}

# ---- Phase: Standardprofil ------------------------------------------------------------------

phase_standard() { # <N> <Kennung>
  local N=$1 tag=$2 k i j
  section "[$tag] Standardprofil, $N Instanzen, limit"
  new_project standard limit
  local dir=$PDIR name=$PNAME
  local -A SC
  for k in $(seq 1 "$N"); do
    echo "-- Instanz $k"
    if ! up "$dir" --id-label "devcontainer.local_folder=$dir" --id-label "u-boot.instance=i$k"; then
      echo "BEFUND: Standardprofil Instanz $k scheitert"; STATUS="Abbruch im Standardprofil bei Instanz $k (Lauf $tag)"; return 1
    fi
    SC[$k]=$LAST_CID
    echo "  Grenze zurueckgelesen: $(limits "${SC[$k]}")"
  done
  echo "Mounts der Instanz 1:"; docker inspect --format '{{range .Mounts}}  {{.Type}} {{.Source}} -> {{.Destination}}{{"\n"}}{{end}}' "${SC[1]}"
  echo "Volumes der Instanzen (leer = Bind-Mount): $(for k in $(seq 1 "$N"); do docker inspect --format '{{range .Mounts}}{{if eq .Type "volume"}}{{.Name}} {{end}}{{end}}' "${SC[$k]}"; done)"
  echo "-- M3 Standardprofil: Datei im Projektordner"
  docker exec "${SC[1]}" sh -c "echo std-1 > /workspaces/$name/ubmess-marker-1" 2>&1 | head -n 2
  for j in $(seq 1 "$N"); do echo "Instanz $j sieht ubmess-marker-1: $(docker exec "${SC[$j]}" sh -c "cat /workspaces/$name/ubmess-marker-1 2>&1")"; done
  echo "Host sieht die Datei im Projektordner: $(cat "$dir/ubmess-marker-1" 2>&1)"
  docker exec -u root "${SC[1]}" rm -f "/workspaces/$name/ubmess-marker-1"
  teardown_dir "$dir"

  section "[$tag] M6 umgekehrt: Bestandscontainer ohne Flag, danach Instanz"
  if up "$dir"; then
    local first=$LAST_CID
    labels_of "$first"
    if up "$dir" --id-label "devcontainer.local_folder=$dir" --id-label "u-boot.instance=i1"; then
      echo "Instanz i1 uebernimmt den Bestandscontainer: $([ "$LAST_CID" = "$first" ] && echo ja || echo "nein (neuer Container ${LAST_CID:0:12})")"
    fi
    echo "Container danach (Filter local_folder):"
    docker ps -a --filter "label=devcontainer.local_folder=$dir" --format '  {{.ID}} instance={{.Label "u-boot.instance"}} {{.Status}}'
    echo "Erneut ohne Flag:"
    if up "$dir"; then
      echo "ohne Flag trifft: ${LAST_CID:0:12} (Bestandscontainer ${first:0:12}); Container gesamt: $(docker ps -aq --filter "label=devcontainer.local_folder=$dir" | wc -l)"
    fi
  fi
  teardown_dir "$dir"
}

# ---- M8: Namen -----------------------------------------------------------------------------

phase_namen() {
  section "M8: Label-Werte (Standardprofil)"
  new_project standard limit
  local dir=$PDIR nm label cnt
  local long32 long63 long64 long128
  long32=$(printf 'a%.0s' $(seq 1 32)); long63=$(printf 'a%.0s' $(seq 1 63))
  long64=$(printf 'a%.0s' $(seq 1 64)); long128=$(printf 'a%.0s' $(seq 1 128))
  local names=("a" "$long32" "$long63" "$long64" "$long128" "Abc" "abc-1" "abc_1" "abc.1" ".abc" "_abc" "-abc" "abc." "abc_" "abc-" "a b" 'a"b' "a'b" "ä-ö" "a,b" "a=b" "a*b" "a/b" "a:b" "")
  for nm in "${names[@]}"; do
    printf '\n-- Name (Laenge %s): %q\n' "${#nm}" "${nm:0:40}"
    if up "$dir" --id-label "devcontainer.local_folder=$dir" --id-label "u-boot.instance=$nm"; then
      label=$(docker inspect --format '{{index .Config.Labels "u-boot.instance"}}' "$LAST_CID")
      echo "  Label-Wert am Container: laenge=${#label} gleich dem Eingabewert: $([ "$label" = "$nm" ] && echo ja || echo nein) Wert=$(printf '%q' "${label:0:40}")"
      cnt=$(docker ps -aq --filter "label=u-boot.instance=$nm" | wc -l)
      echo "  docker ps --filter label=u-boot.instance=<Name> trifft: $cnt Container (erwartet 1; Filter mit Wert des Labels)"
      echo "  Treffer mit dem ganzen Namen im Filter 'label=u-boot.instance' (nur Schluessel): $(docker ps -aq --filter "label=u-boot.instance" | wc -l)"
      docker rm -f -v "$LAST_CID" > /dev/null 2>&1
    else
      echo "  BEFUND Name (Laenge ${#nm}) $(printf '%q' "${nm:0:40}"): kein Ergebnis von devcontainer up (exit=$LAST_RC); Container dieses Ordners, die er hinterliess, werden entfernt"
      teardown_dir "$dir"
    fi
  done
  teardown_dir "$dir"
}

# ---- Ablauf -------------------------------------------------------------------------------

section "M0: Messstand"
date -u +"Zeit (UTC): %Y-%m-%dT%H:%M:%SZ"
echo "Stufe/Modus: $stufe"
echo "Host: $(uname -srm); uname -m: $(uname -m)"
echo "u-boot: $("$bin" --version 2>&1 | head -n 1) ($bin)"
docker version --format 'Docker Client {{.Client.Version}} ({{.Client.Os}}/{{.Client.Arch}}), Server {{.Server.Version}} ({{.Server.Os}}/{{.Server.Arch}})' 2>&1
docker info --format 'Engine: {{.OperatingSystem}}; Name {{.Name}}; Treiber {{.Driver}}; Cgroup {{.CgroupVersion}} ({{.CgroupDriver}}); MemoryLimit {{.MemoryLimit}} SwapLimit {{.SwapLimit}}; CPUs {{.NCPU}}; MemTotal {{.MemTotal}}; Kernel {{.KernelVersion}}; Architektur {{.Architecture}}; Wurzel {{.DockerRootDir}}; Context-Endpunkt: lokal' 2>&1
echo "docker context: $(docker context show 2>&1); DOCKER_HOST=${DOCKER_HOST:-<ungesetzt>}"
free -m
echo "Plattenplatz /: $(df -h / | tail -n 1)"
echo "kernel.apparmor_restrict_unprivileged_userns = $(sysctl -n kernel.apparmor_restrict_unprivileged_userns 2>&1)"
echo "Hilfs-Image vorhanden (Bestand): $(docker images --format '{{.Repository}}:{{.Tag}} {{.ID}}' | grep '^u-boot-devcontainer-cli:' | tr '\n' ' ')"

section "Vorpruefung"
ok=1
docker info > /dev/null 2>&1 && echo "(a) docker info antwortet: ja" || { echo "(a) docker info antwortet: NEIN"; ok=0; }
[ "$(docker info --format '{{.MemoryLimit}} {{.SwapLimit}} {{.CgroupVersion}}' 2>&1)" = "true true 2" ] && echo "(a) MemoryLimit/SwapLimit wahr, cgroup v2: ja" || { echo "(a) MemoryLimit/SwapLimit/cgroup v2: NEIN"; ok=0; }
echo "(b) Inventur vorher (Rohdaten):"
inventur "$INVDIR/vorher"
echo "--- docker ps -a"; docker ps -a --format '{{.Names}}	{{.Image}}	{{.Status}}'
echo "--- docker volume ls"; cat "$INVDIR/vorher.vol"
echo "--- docker images"; cat "$INVDIR/vorher.img"
pre=$( { cat "$INVDIR/vorher.ps" "$INVDIR/vorher.vol" "$INVDIR/vorher.img"; } | grep -cE '(^|[[:space:]])(vsc-)?ubmess-')
[ "$pre" = 0 ] && echo "(c) kein Objekt mit Praefix ubmess-/vsc-ubmess- vorhanden: ja" || { echo "(c) Objekte mit Praefix vorhanden: $pre -> NEIN"; ok=0; }
if [ "$stufe" != namen ]; then
  need=$((stufe * 3 * 1048576 + 1572864))
  have=$(memavail_kb)
  echo "(d) MemAvailable $have kB; benoetigt $need kB (N x 3 GiB + 1,5 GiB, N=$stufe; 3 GiB wie --memory=3g)"
  if [ "$have" -ge "$need" ]; then echo "(d) Speicher: erfuellt"
  elif [ "$stufe" = 4 ] && [ "$have" -ge 9437184 ] && [ "$(awk '/^SwapFree:/{print $2}' /proc/meminfo)" -ge 8388608 ]; then
    RELAXED=1
    echo "(d) Speicher: NICHT nach Slice-Schwelle erfuellt; ABWEICHUNG VON DER VORPRUEFUNG DES SLICE (Projektinhaber, 2026-10-31, nur Stufe 4): Start erlaubt, weil MemAvailable >= 9 GiB (drei Instanzen) und SwapFree >= 8 GiB ($(awk '/^SwapFree:/{print $2}' /proc/meminfo) kB); Instanzgrenze bleibt --memory=3g --memory-swap=3g; Ueberwachung von Host-Swap und vmstat si/so; Swap-Druck heisst: der Host-Swap waechst um mehr als 256 MiB oder si+so bleibt dauerhaft hoch"
  else echo "(d) Speicher: NICHT erfuellt"; ok=0; fi
  avail_gb=$(df -BG --output=avail / | tail -n 1 | tr -dc 0-9)
  echo "(d) freier Platz der Docker-Wurzel (Dateisystem /): ${avail_gb} GB (benoetigt 40 GB)"
  [ "$avail_gb" -ge 40 ] && echo "(d) Platz: erfuellt" || { echo "(d) Platz: NICHT erfuellt"; ok=0; }
fi
echo "(e) kernel.apparmor_restrict_unprivileged_userns = $(sysctl -n kernel.apparmor_restrict_unprivileged_userns 2>&1): $([ "$(sysctl -n kernel.apparmor_restrict_unprivileged_userns 2>&1)" = 0 ] && echo 'Podman-Probe moeglich' || echo 'Podman-Familie NICHT MESSBAR auf diesem Host (Host wird nicht umgestellt)')"
if [ "$ok" != 1 ]; then
  echo "BEFUND: Vorpruefung nicht erfuellt; Stufe $stufe wird nicht gemessen, nichts wurde angelegt oder geaendert."
  STATUS="nicht gemessen, Vorpruefung nicht erfuellt"
  exit 3
fi

[ "$stufe" = namen ] && UP_LIMIT=90
DEADLINE=$(( $(now) + STAGE_LIMIT ))
echo "Zeitlimits: je devcontainer up ${UP_LIMIT} s, je Lauf ${STAGE_LIMIT} s (Werte aus den Laeufen der Stufen 2 und 3: ein up dauert bis 7 s, ein Lauf 67 bis 83 s; die Limits liegen darueber um mehr als das Zehnfache bzw. fuer ein kaltes Image im Rahmen eines Baus)"
if [ "$stufe" != namen ]; then
  mapfile -t FOREIGN < <(docker ps -q)
  monitor > /dev/null 2>&1 &
  MON_PID=$!
  echo "Ueberwachung gestartet (alle ~2 s: MemAvailable, SwapFree, vmstat si/so, Lauf der ${#FOREIGN[@]} vorhandenen Container); Abbruchschwelle: si+so >= 5000 KiB/s in 5 aufeinanderfolgenden Proben (~10 s) oder ein vorhandener Container laeuft nicht mehr bzw. OOMKilled. Begruendung: ein ruhender Host zeigt 0; dauerhaftes Auslagern im Megabyte-Takt ist Swap-Druck, kurze Spitzen beim Start sind es nicht. Relaxed (Swap-Lockerung): ${RELAXED:-0}"
fi

section "Hilfs-Image bauen (ubmess-devcontainer-cli)"
t0=$(now)
docker build -q -t "$HELPER_IMG" "$ctx" > /dev/null 2> "$base/build.err" || { echo "BEFUND: Bau des Hilfs-Images scheitert"; tail -n 10 "$base/build.err"; STATUS="Abbruch: Hilfs-Image nicht baubar"; exit 4; }
echo "Bau: $(( $(now) - t0 )) s; Image-ID $HELPER_IMG: $(docker images -q "$HELPER_IMG")"
echo "Image-ID u-boot-devcontainer-cli (Bestand): $(docker images --format '{{.Repository}}:{{.Tag}}={{.ID}}' | grep '^u-boot-devcontainer-cli:' | tr '\n' ' ')"
echo "Dev-Container-CLI: $(docker run --rm "$HELPER_IMG" devcontainer --version 2>&1)"
echo "Tag des Mess-Images: ${HELPER_IMG#*:} (Kontextdateien: $(cd "$ctx" && sha256sum Dockerfile package.json package-lock.json | cut -c1-12 | tr '\n' ' '))"

if [ "$stufe" = namen ]; then
  phase_namen
else
  N=$stufe
  phase_sandbox "$N" limit "S$N-limit" || true
  phase_standard "$N" "T$N-standard" || true
  if [ "$N" = 2 ]; then
    phase_sandbox 2 nolimit "S2-bezug-ohne-limit" || true
  fi
fi
