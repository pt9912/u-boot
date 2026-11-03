# ablage-pruefung.mk — Positivliste der Ablage (Fragment des Repos;
# Form-Bruder von tool-pruefung.mk / devcontainer.mk, per include aus repo.mk
# gebunden).
#
# Prueft die Ablage gegen [harness/conventions.md](../../harness/conventions.md)
# (Abschnitt Ablage von Messprotokollen und Skripten): im Planning-Baum nur
# Dokumente, unter observations/ nur die Beobachtungsform, unter docs/ kein
# Skript/Dockerfile. bash + git, kein Docker (Muster
# tools/harness/history-range-guard.sh). Das Skript und sein Waechter-Test
# liegen unter tools/harness/.

GATE_CHECKS += ablage-pruefung

.PHONY: ablage-pruefung
ablage-pruefung: ## Ablage im Planning-Baum und unter docs/ gegen die Positivliste pruefen.
	@bash tools/harness/ablage-pruefung.sh
	@bash tools/harness/ablage-pruefung-test.sh