#!/usr/bin/env bash
# messen.sh — Messung des Slice slice-v1-devcontainer-ressourcenlimits.
#
# FRAGE: Überlebt eine von Hand in .devcontainer/devcontainer.json gesetzte
# Zeile "runArgs": ["--memory=4g", "--cpus=2"] ein erneutes
# `u-boot generate devcontainer`?
#
# WAS ES MISST (nur Dateien in einem Wegwerf-Ordner, kein Docker, kein Container):
#   A1. Standardprofil: runArgs INNERHALB des Managed Blocks, danach generate.
#   A2. dasselbe mit `generate devcontainer --force` (das Flag existiert nicht).
#   A3. Sandbox-Profil: Rendert der Generator ohne weitere Features eine
#       runArgs-Zeile? Wenn nicht, ist dort nichts zu ergänzen.
#   B.  Standardprofil: eine Kommentarzeile AUSSERHALB des Blocks (hinter
#       END U-BOOT MANAGED BLOCK), danach generate.
#
# AUFRUF (aus der Repo-Wurzel, `u-boot` im PATH):
#   bash tools/messung/devcontainer-ressourcenlimits/messen.sh \
#     > docs/messungen/slice-v1-devcontainer-ressourcenlimits/roh-lauf.txt 2>&1
#
# SCHREIBT NICHT in das Repo, außer nach stdout; der Arbeitsordner entsteht mit
# mktemp und wird am Ende entfernt.
set -u

command -v u-boot >/dev/null || { echo "messen.sh: u-boot nicht im PATH" >&2; exit 2; }

base="$(mktemp -d)"
trap 'rm -rf "$base"' EXIT
cd "$base" || exit 1

echo "# Datum: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
echo "# Host: $(uname -sm)"
u-boot --version

u-boot init --no-git --no-interactive demo >init.out 2>&1; echo "init exit=$?"
u-boot config set devcontainer.enabled true >/dev/null 2>&1
echo "=== Ausgangslage: Standardprofil"
u-boot generate devcontainer >gen1.out 2>&1; echo "generate#1 exit=$?"; head -5 gen1.out
cat .devcontainer/devcontainer.json
cp .devcontainer/devcontainer.json vorher.json

echo "=== A1 Standardprofil: Zeile runArgs innerhalb des Blocks einfügen"
awk '{ print } /"name": "demo",/ { print "  \"runArgs\": [\"--memory=4g\", \"--cpus=2\"]," }' vorher.json >.devcontainer/devcontainer.json
sed -n 1,6p .devcontainer/devcontainer.json
u-boot generate devcontainer >gen4.out 2>&1; echo "generate exit=$?"; head -8 gen4.out
grep -n "runArgs" .devcontainer/devcontainer.json || echo "runArgs nach generate: nicht mehr vorhanden"
ls -a .devcontainer

echo "=== A2 dasselbe mit --force"
awk '{ print } /"name": "demo",/ { print "  \"runArgs\": [\"--memory=4g\", \"--cpus=2\"]," }' vorher.json >.devcontainer/devcontainer.json
u-boot generate devcontainer --force >gen5.out 2>&1; echo "generate --force exit=$?"; head -8 gen5.out
grep -n "runArgs" .devcontainer/devcontainer.json || echo "runArgs nach generate --force: nicht mehr vorhanden"
ls -a .devcontainer

echo "=== B Standardprofil: Kommentarzeile AUSSERHALB des Blocks (nach END)"
cp vorher.json .devcontainer/devcontainer.json
printf '\n// HAND: --memory=4g\n' >>.devcontainer/devcontainer.json
u-boot generate devcontainer >gen3.out 2>&1; echo "generate exit=$?"; head -8 gen3.out
tail -4 .devcontainer/devcontainer.json
ls -a .devcontainer

echo "=== A3 Sandbox-Profil (rendert runArgs selbst?)"
u-boot config set devcontainer.profile sandbox >/dev/null 2>&1; echo "profile exit=$?"
u-boot generate devcontainer >gen6.out 2>&1; echo "generate sandbox exit=$?"; head -8 gen6.out
grep -n "runArgs" .devcontainer/devcontainer.json || echo "Sandbox-Profil: keine runArgs-Zeile gerendert"
cp .devcontainer/devcontainer.json sandbox-vorher.json
# Ohne runArgs-Zeile ändert awk nichts: die Probe zeigt dann nur, dass ein
# unverändertes generate "up to date" meldet, nicht das Überleben einer Ergänzung.
awk '{ if ($0 ~ /"runArgs"/) { sub(/\]/, ", \"--memory=4g\", \"--cpus=2\"]") } print }' sandbox-vorher.json >.devcontainer/devcontainer.json
grep -n "runArgs" .devcontainer/devcontainer.json || echo "keine runArgs-Zeile: Ergänzung nicht möglich"
u-boot generate devcontainer >gen7.out 2>&1; echo "generate#2 sandbox exit=$?"; head -8 gen7.out
