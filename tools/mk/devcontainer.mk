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

define dev_require
@command -v $(U_BOOT) >/dev/null 2>&1 || { echo "$(1): '$(U_BOOT)' not found on PATH" >&2; exit 2; }
endef

.PHONY: dev-up dev-down dev-shell

dev-up: ## Devcontainer starten (RECREATE=1 legt ihn neu an)
	$(call dev_require,dev-up)
	$(U_BOOT) devcontainer up $(if $(RECREATE),--recreate)

dev-down: ## Devcontainer stoppen (er bleibt bestehen)
	$(call dev_require,dev-down)
	$(U_BOOT) devcontainer down

dev-shell: ## Devcontainer starten und eine Shell darin oeffnen
	$(call dev_require,dev-shell)
	$(U_BOOT) devcontainer shell
