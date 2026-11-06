# tools/mk/devcontainer.mk — Bedien-Targets fuer den Devcontainer dieses Repos
# (.devcontainer/, u-boot.yaml). Sie pruefen nichts und haengen nicht an
# GATE_CHECKS; repo.mk bindet die Datei ein, .d-check.yml (targets.makefiles)
# liest sie fuer das Doku-Gate.
#
# Die Targets rufen `u-boot devcontainer` auf. Das Binary `u-boot` auf dem
# Host ist hier Bedienung, kein Gate-Vertrag: die Gates brauchen weiter nur
# Docker und GNU make. Fehlt das Binary, endet das Target mit Exit 2.
# Der Container laeuft im Sandbox-Profil mit verschachteltem Docker; sein
# Workspace ist ein Clone von origin (u-boot.yaml), nicht der Arbeitsbaum.

U_BOOT ?= u-boot

# Die erzeugte devcontainer.json verweist mit ${localEnv:GIT_TOKEN} und
# ${localEnv:GH_TOKEN} auf die Umgebung des Hosts, und u-boot gibt nur gesetzte
# Variablen weiter. dev-up und dev-shell setzen beide aus `gh auth token`, wenn
# die Variable fehlt; ein gesetzter Wert hat Vorrang. Ist gh nicht installiert
# oder nicht angemeldet, bleiben sie leer. Die Werte gehen ueber die Umgebung
# des Rezepts, nie ueber die Befehlszeile. Die anderen Targets setzen sie nicht.
dev-up dev-shell: export GIT_TOKEN ?= $(shell gh auth token 2>/dev/null)
dev-up dev-shell: export GH_TOKEN ?= $(shell gh auth token 2>/dev/null)

define dev_require
@command -v $(U_BOOT) >/dev/null 2>&1 || { echo "$(1): '$(U_BOOT)' not found on PATH" >&2; exit 2; }
endef

.PHONY: dev-up dev-down dev-shell dev-remove

dev-up: ## Devcontainer starten (RECREATE=1 legt ihn neu an, INSTANCE=<name> waehlt eine Instanz, GIT_TOKEN und GH_TOKEN aus gh auth token)
	$(call dev_require,dev-up)
	$(U_BOOT) devcontainer up $(if $(RECREATE),--recreate) $(if $(INSTANCE),--instance $(INSTANCE))

dev-down: ## Devcontainer stoppen (er bleibt bestehen; INSTANCE=<name> waehlt eine Instanz)
	$(call dev_require,dev-down)
	$(U_BOOT) devcontainer down $(if $(INSTANCE),--instance $(INSTANCE))

dev-shell: ## Devcontainer starten und eine Shell darin oeffnen (INSTANCE=<name> waehlt eine Instanz, GIT_TOKEN und GH_TOKEN aus gh auth token)
	$(call dev_require,dev-shell)
	$(U_BOOT) devcontainer shell $(if $(INSTANCE),--instance $(INSTANCE))

dev-remove: ## Devcontainer entfernen (INSTANCE=<name> waehlt eine Instanz)
	$(call dev_require,dev-remove)
	$(U_BOOT) devcontainer remove --volumes --yes $(if $(INSTANCE),--instance $(INSTANCE))
