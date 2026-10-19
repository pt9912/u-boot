#!/usr/bin/env bash
# ablage-pruefung — Positivliste der Ablage im Planning-Baum und unter docs/.
#
# ZUSAGE (harness/conventions.md, Abschnitt "Ablage von Messprotokollen und Skripten"):
#   1. Im Planning-Baum (docs/plan/planning/) liegen nur *.md und die drei Platzhalter
#      open/.gitkeep, next/.gitkeep, done/.gitkeep.
#   2. Unter docs/plan/planning/observations/ liegt nur die Beobachtungsform:
#      README.md und BEO-<KUERZEL>/<slug>/{observation.md,state.md,evidence/<vorgangs-id>.md}.
#   3. Unter docs/ liegt kein Skript und kein Dockerfile (*.sh, *.bash, *.pl, *.py, *.awk,
#      *.mk, Dockerfile*).
#
# GEPRUEFT wird die Menge der von git verfolgten Dateien (`git ls-files`, Index eingeschlossen),
# nicht der Arbeitsbaum. Eine Datei ausserhalb der Positivliste wird umgezogen (Skripte nach
# tools/messung/<name>/, Protokolle nach docs/messungen/<kennung>/); die Liste hier zu erweitern
# ist eine bewusste Aenderung mit Begruendung, keine Ausnahme im Aufrufer.
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
skript='(\.(sh|bash|pl|py|awk|mk)$|(^|/)Dockerfile[^/]*$)'

# bewerte <pfad> — gibt den Grund aus, wenn der Pfad ausserhalb der Positivliste liegt;
# gibt nichts aus, wenn er drin liegt.
bewerte() {
  local f="$1"
  case "$f" in
    "$observations"*)
      [[ "$f" =~ $beobachtung ]] ||
        echo "unter observations/ nur Beobachtungsform (observation.md, state.md, evidence/<vorgangs-id>.md)"
      ;;
    "$planning"*)
      if [[ "$f" == *.md ]] || [[ "$f" =~ $platzhalter ]]; then
        return 0
      fi
      echo "im Planning-Baum nur *.md und open/next/done/.gitkeep"
      ;;
    docs/*)
      if [[ "$f" =~ $skript ]]; then
        echo "unter docs/ kein Skript und kein Dockerfile (Skripte nach tools/messung/<name>/)"
      fi
      ;;
  esac
}

fehler=0
while IFS= read -r -d '' f; do
  grund="$(bewerte "$f")"
  if [ -n "$grund" ]; then
    echo "ablage-pruefung: $f -- $grund" >&2
    fehler=$((fehler + 1))
  fi
done < <(git ls-files -z -- docs)

if [ "$fehler" -ne 0 ]; then
  echo "ablage-pruefung: $fehler Datei(en) ausserhalb der Positivliste." >&2
  exit 1
fi
echo "ablage-pruefung: OK"
