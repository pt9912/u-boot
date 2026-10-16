# harness/mk/hooks-install.mk — Fragment des Commit-Kennungs-Waechters (Herkunft:
# ai-harness-init, MR-011). EIN KOMMANDO, KEIN GATE.
#
# DER TRAEGER LIEGT VERSIONIERT UNTER .githooks/commit-msg UND REIST MIT DEM KLON;
# SEINE AKTIVIERUNG TUT DAS NICHT. `core.hooksPath` ist lokale Konfiguration,
# dieses Ziel ist der eine Schritt dazwischen. Ein Klon ohne diesen Aufruf ist
# ungeprueft, und `git commit --no-verify` umgeht den Traeger auch danach — ein
# Stolperdraht, keine Sandbox.
#
# DER TRAEGER GEHOERT DEM REPO. Der Name ist von git fixiert; die Aktivierung
# unten nimmt das Programm in .githooks/ in Betrieb. Die Pruefung daneben,
# tools/harness/commit-msg-traceability.sh, kann ein eigener Traeger aufrufen.
#
# WAS DER TRAEGER NICHT ERREICHT. Die zweite Haelfte der Traceability-Zusage —
# das Doku-Update bei beruehrtem oeffentlichem Vertrag — ist von einem
# Commit-Waechter nicht mechanisch pruefbar und bleibt Arbeit des Committenden.
# Geprueft wird die ANWESENHEIT einer Kennung, nicht ihre Wahrheit: eine Message,
# die zusaetzlich einen nicht aufloesbaren Verweis nennt, geht mit derselben
# Kennung durch. Beide Grenzen stehen auch im Kopf von
# tools/harness/commit-msg-traceability.sh.
#
# WAS ER MITNIMMT. Er haengt am Commit und sieht jede Klasse, die `git` erzeugt —
# auch die Commits der Repo-Werkzeuge. `make slice-mv`
# committet intern mit dem Slice-Namen; er ist eine Kennung aus der Menge, der
# Commit des Werkzeugs besteht den Waechter.
#
# ABHAENGIGKEIT. Dieses Rezept ruft `git` — die zugelassene Host-Abhaengigkeit
# dieses Aufbaus; der Traeger selbst laeuft im Commit-Pfad und setzt weder Docker
# noch Netz noch ein Gate-Bild voraus.
.PHONY: hooks-install

# Das Verzeichnis des Traegers ist das des Repos; git fixiert den Hook-Namen.
HOOKS_DIR ?= .githooks

# Fehlt der Traeger, sagt dieses Ziel das und bricht ab: `git config` haette sonst
# einen Pfad gesetzt, unter dem nichts liegt, und jeder Commit liefe ungeprueft
# durch, waehrend die Konfiguration einen Waechter behauptet.
# Das Ausfuehrrecht setzt das Rezept nach: git verwirft einen nicht ausfuehrbaren
# Hook still, und ein verlorenes Bit waere ein Waechter, der nur so aussieht.
hooks-install: ## den git-eigenen commit-msg-Traeger im Klon aktivieren (core.hooksPath) — KEIN Gate
	@test -f "$(HOOKS_DIR)/commit-msg" || { echo "hooks-install: $(HOOKS_DIR)/commit-msg liegt nicht — dieses Ziel aktiviert den commit-msg-Traeger dieses Verzeichnisses."; exit 2; }
	@test -x "$(HOOKS_DIR)/commit-msg" || chmod +x "$(HOOKS_DIR)/commit-msg"
	@git config core.hooksPath "$(HOOKS_DIR)"
	@printf '%s\n' "hooks-install: core.hooksPath=$$(git config --get core.hooksPath) — $(HOOKS_DIR)/commit-msg laeuft ab dem naechsten Commit (Umgehung: git commit --no-verify)."
