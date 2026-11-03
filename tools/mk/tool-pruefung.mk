# tool-pruefung.mk — das Laufzeit-Gate der Pin-Inventur (Fragment des Repos;
# Form-Bruder von tools/mk/devcontainer.mk, per include aus repo.mk gebunden).
#
# Zwei Schichten decken die Pins ([LH-FA-DEV-011](../../../spec/lastenheft.md),
# Pin-Inventur-Absatz): pinsync (d-Check-Modul, im docs-check-Gatter) prueft die
# WERTE an den Stellen; tool-pruefung prueft die INSTALLIERTEN Werkzeuge des
# laufenden Kontexts gegen dieselben Stellen (gelesen, nicht kopiert). Ein Lauf,
# dessen Werkzeugkette nicht die Pins traegt, endet rot — der Lauf-Ort ist
# egal, die Werkzeuge tragen ([ADR-0043](../../docs/plan/adr/0043-nur-devcontainer-bau-pfad.md)).
# Das Skript und sein Rot-Beweis stehen unter tools/harness/tool-pruefung.sh;
# der Lauf im Devcontainer-Container (die Standard-Umgebung) und am bare Host
# (ohne Werkzeugkette) trugen beide Proben des Slices.

GATE_CHECKS += tool-pruefung

.PHONY: tool-pruefung
tool-pruefung: ## Installed tools of the running context match the pins of the inventory (LH-FA-DEV-011, ADR-0043).
	bash tools/harness/tool-pruefung.sh