#!/usr/bin/env bash
# ablage-pruefung — Positivliste der Ablage unter docs/ (harness/conventions.md,
# Abschnitt "Ablage von Messprotokollen und Skripten").
#
# GEPRUEFT wird die Menge der von git verfolgten Dateien unter docs/ (`git ls-files -s`,
# Index eingeschlossen), nicht der Arbeitsbaum. Eine Datei ausserhalb der Liste wird
# umgezogen (Skripte nach tools/messung/<name>/, Protokolle nach docs/messungen/<kennung>/);
# die Liste hier zu erweitern ist eine bewusste Aenderung mit Begruendung.
#
# Regeln, in dieser Reihenfolge je Datei:
#   1. Modus: unter docs/ nur 100644 (kein Ausfuehrungsbit, kein Symlink).
#   2. Planning-Baum (docs/plan/planning/): nur Namen auf `.md` (kleingeschrieben) und
#      open/.gitkeep, next/.gitkeep, done/.gitkeep; unter observations/ nur README.md und
#      BEO-<KUERZEL>/<slug>/{observation.md,state.md,evidence/<vorgangs-id>.md}.
#   3. docs/messungen/: nur nicht-versteckte Dateien auf `.md`, `.txt`, `.log` (kleingeschrieben).
#   4. Im uebrigen docs/: kein Name auf .sh .bash .pl .py .awk .mk .ps1 oder Dockerfile*
#      (Gross-/Kleinschreibung gleich).
#   5. Unter docs/ beginnt keine Datei mit `#!` (Shebang).
#
# Nicht geprueft: Inhalt jenseits der ersten zwei Bytes, Dateien ausserhalb von docs/.
#
# Kein Docker, kein Netz: bash + git (Muster history-range-guard.sh). Exit 0: alles in der Liste.
# Exit 1: mindestens eine Datei ausserhalb. Exit 2: kein git-Arbeitsbaum.
set -euo pipefail

root="$(git rev-parse --show-toplevel 2>/dev/null)" || {
  echo "ablage-pruefung: kein git-Arbeitsbaum." >&2
  exit 2
}
cd "$root"

planning="docs/plan/planning/"
observations="docs/plan/planning/observations/"
beobachtung='^docs/plan/planning/observations/(README\.md|BEO-[A-Z0-9]+/[^/]+/(observation\.md|state\.md|evidence/[^/]+\.md))$'
platzhalter='^docs/plan/planning/(open|next|done)/\.gitkeep$'
protokoll='^docs/messungen/(.*/)?[^./][^/]*\.(md|txt|log)$'
skript='(\.(sh|bash|pl|py|awk|mk|ps1)$|(^|/)dockerfile[^/]*$)'

# bewerte <modus> <pfad> — gibt den Grund aus, wenn der Pfad ausserhalb der Positivliste liegt;
# gibt nichts aus, wenn er drin liegt.
bewerte() {
  local modus="$1" f="$2" klein
  if [ "$modus" != "100644" ]; then
    echo "unter docs/ nur normale Dateien ohne Ausfuehrungsbit (Modus $modus)"
    return 0
  fi
  case "$f" in
    "$observations"*)
      [[ "$f" =~ $beobachtung ]] ||
        echo "unter observations/ nur Beobachtungsform (observation.md, state.md, evidence/<vorgangs-id>.md)"
      ;;
    "$planning"*)
      if [[ "$f" == *.md ]] || [[ "$f" =~ $platzhalter ]]; then
        :
      else
        echo "im Planning-Baum nur *.md und open/next/done/.gitkeep"
      fi
      ;;
    docs/messungen/*)
      [[ "$f" =~ $protokoll ]] ||
        echo "unter docs/messungen/ nur .md, .txt, .log (Skripte nach tools/messung/<name>/)"
      ;;
    *)
      klein="$(printf '%s' "$f" | tr '[:upper:]' '[:lower:]')"
      if [[ "$klein" =~ $skript ]]; then
        echo "unter docs/ kein Skript und kein Dockerfile (Skripte nach tools/messung/<name>/)"
      fi
      ;;
  esac
  if [ -f "$f" ] && [ "$(head -c 2 -- "$f")" = "#!" ]; then
    echo "unter docs/ keine Datei mit Shebang (Skripte nach tools/messung/<name>/)"
  fi
}

fehler=0
while IFS= read -r -d '' eintrag; do
  meta="${eintrag%%$'\t'*}"
  f="${eintrag#*$'\t'}"
  modus="${meta%% *}"
  grund="$(bewerte "$modus" "$f")"
  if [ -n "$grund" ]; then
    echo "ablage-pruefung: $f -- $grund" >&2
    fehler=$((fehler + 1))
  fi
done < <(git ls-files -s -z -- docs)

if [ "$fehler" -ne 0 ]; then
  echo "ablage-pruefung: $fehler Datei(en) ausserhalb der Positivliste." >&2
  exit 1
fi
echo "ablage-pruefung: OK"
