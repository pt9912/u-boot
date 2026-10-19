#!/usr/bin/env bash
# ablage-pruefung-test — Waechter-Test fuer ablage-pruefung.sh.
#
# Je Fall entsteht ein frisches git-Repo in einem Wegwerf-Verzeichnis. Der Grundzustand
# (Beobachtungsform, Platzhalter, ein Protokoll unter docs/messungen/, ein Skript unter
# tools/messung/) muss gruen sein; jede Mutation (eine Datei ausserhalb der Positivliste)
# muss das Skript rot faerben und den Pfad der Datei nennen. Ein Fall, der nicht greift,
# bricht den Test.
#
# Kein Docker, kein Netz: bash + git.
set -euo pipefail

skript="$(cd "$(dirname "$0")" && pwd)/ablage-pruefung.sh"
arbeit="$(mktemp -d)"
trap 'rm -rf "$arbeit"' EXIT

versagt=0
faelle=0

# neues_repo — leert das Arbeitsverzeichnis und legt den gruenen Grundzustand an.
neues_repo() {
  rm -rf "${arbeit:?}/r"
  mkdir -p "$arbeit/r"
  cd "$arbeit/r"
  git init -q
  local d
  for d in open next done; do
    mkdir -p "docs/plan/planning/$d"
    : >"docs/plan/planning/$d/.gitkeep"
  done
  mkdir -p docs/plan/planning/in-progress docs/plan/planning/observations/BEO-DOK/x/evidence \
    docs/plan/planning/observations/BEO-E2E/y/evidence docs/messungen/slice-a tools/messung/a
  : >docs/plan/planning/in-progress/slice-a.md
  : >docs/plan/planning/observations/README.md
  : >docs/plan/planning/observations/BEO-DOK/x/observation.md
  : >docs/plan/planning/observations/BEO-DOK/x/state.md
  : >docs/plan/planning/observations/BEO-DOK/x/evidence/slice-a.md
  : >docs/plan/planning/observations/BEO-E2E/y/observation.md
  : >docs/plan/planning/observations/BEO-E2E/y/state.md
  : >docs/plan/planning/observations/BEO-E2E/y/evidence/slice-a.md
  : >docs/messungen/slice-a/README.md
  : >docs/messungen/slice-a/rohdaten-1.txt
  : >docs/messungen/slice-a/rohdaten-2.log
  : >tools/messung/a/run.sh
  : >tools/messung/a/Dockerfile.cli
}

# erwarte <name> <exit> [<pfad-in-meldung>] — fuehrt das Skript im aktuellen Repo aus.
erwarte() {
  local name="$1" soll="$2" pfad="${3:-}" ist=0 aus
  faelle=$((faelle + 1))
  aus="$(bash "$skript" 2>&1)" || ist=$?
  if [ "$ist" -ne "$soll" ]; then
    echo "FAIL  $name: Exit $ist, erwartet $soll" >&2
    echo "$aus" | sed 's/^/      /' >&2
    versagt=$((versagt + 1))
    return
  fi
  if [ -n "$pfad" ] && ! grep -qF -- "$pfad" <<<"$aus"; then
    echo "FAIL  $name: Meldung nennt '$pfad' nicht" >&2
    echo "$aus" | sed 's/^/      /' >&2
    versagt=$((versagt + 1))
    return
  fi
  echo "ok    $name"
}

# mutation <name> <pfad> — legt die Datei an, stagt sie und erwartet Rot mit dem Pfad.
mutation() {
  local name="$1" pfad="$2" inhalt="${3:-}" art="${4:-datei}"
  neues_repo
  mkdir -p "$(dirname "$pfad")"
  case "$art" in
    symlink) ln -s ziel "$pfad" ;;
    *) printf '%s' "$inhalt" >"$pfad" ;;
  esac
  git add -A
  [ "$art" = "exec" ] && git update-index --chmod=+x -- "$pfad"
  erwarte "$name" 1 "$pfad"
}

# gruen <name> <pfad> — legt die Datei an, stagt sie und erwartet Gruen (Gegenprobe zu zu engen Regeln).
gruen() {
  local name="$1" pfad="$2"
  neues_repo
  mkdir -p "$(dirname "$pfad")"
  : >"$pfad"
  git add -A
  erwarte "$name" 0 "OK"
}

neues_repo
git add -A
erwarte "Grundzustand ist gruen" 0 "OK"

mutation "Skript unter docs/messungen/" docs/messungen/slice-a/run.sh
mutation "Dockerfile unter docs/messungen/" docs/messungen/slice-a/Dockerfile.cli
mutation "Perl-Sonde unter docs/" docs/messungen/slice-a/probe.pl
mutation "Skript im Planning-Baum" docs/plan/planning/in-progress/hilfe.sh
mutation "Skript im Register" docs/plan/planning/observations/BEO-DOK/x/evidence/run.sh
mutation "Rohdaten im Register" docs/plan/planning/observations/BEO-DOK/x/evidence/slice-a/rohdaten-1.txt
mutation "Beobachtung mit fremdem Dateinamen" docs/plan/planning/observations/BEO-DOK/x/notizen.md
mutation "Datei neben der Beobachtung" docs/plan/planning/observations/BEO-DOK/loses.md
mutation "Nicht-md-Datei im Planning-Baum" docs/plan/planning/in-progress/grafik.png
mutation ".gitkeep in einem vierten Verzeichnis" docs/plan/planning/in-progress/.gitkeep
mutation "Skript ohne Endung mit Ausfuehrungsbit" docs/messungen/slice-a/run "echo x" exec
mutation "Skript ohne Endung, Shebang, ausserhalb von messungen" docs/user/run $'#!/bin/sh\necho x\n'
mutation "Grossgeschriebenes Skript" docs/messungen/slice-a/RUN.SH
mutation "Grossgeschriebenes Skript ausserhalb von messungen" docs/user/Run.SH
mutation "dockerfile kleingeschrieben" docs/messungen/slice-a/dockerfile
mutation "PowerShell-Skript" docs/messungen/slice-a/x.ps1
mutation "json unter messungen" docs/messungen/slice-a/x.json
mutation "png unter messungen" docs/messungen/slice-a/x.png
mutation "grossgeschriebene Endung TXT" docs/messungen/slice-a/rohdaten.TXT
mutation "versteckte Datei unter messungen" docs/messungen/slice-a/.rohdaten.txt
mutation "Datei ohne Endung unter messungen" docs/messungen/slice-a/rohdaten
mutation "Symlink unter messungen" docs/messungen/slice-a/link.md "" symlink
mutation "md mit Ausfuehrungsbit" docs/messungen/slice-a/notiz.md "" exec
gruen "png ausserhalb von messungen und planning ist erlaubt" docs/user/bild.png

# Eine Datei, die git nicht verfolgt, ist nicht Gegenstand (Grenze: git ls-files, nicht der Arbeitsbaum).
neues_repo
git add -A
: >docs/messungen/slice-a/ungetrackt.sh
erwarte "ungetrackte Datei zaehlt nicht" 0 "OK"

echo "ablage-pruefung-test: $faelle Faelle, $versagt Fehlschlaege"
[ "$versagt" -eq 0 ]
