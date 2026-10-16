# harness/mk/slice-mv.mk — Fragment des Lifecycle-Wechsels (Herkunft:
# ai-harness-init, MR-011). EIN KOMMANDO, KEIN GATE.
#
# Es liegt im Fragment-Verzeichnis wie die Gate-Fragmente, haengt aber NICHTS an
# GATE_CHECKS und steht in keiner Prerequisite-Kette: ein Move prueft nichts, er
# bewegt.
#
# DER AUFRUF BEWEGT UND COMMITTET IM VERSIONIERTEN BAUM DIESES REPOS.
# Er verschiebt eine Planning-Datei per `git mv`, zieht die Verweise auf sie in
# beiden Richtungen nach und legt zwei Commits an: den reinen Move und, falls
# Verweise anfielen, deren Nachtrag getrennt davon. Einen unsauberen Arbeitsbaum
# bricht er ab, statt fremden Inhalt in einen der Commits zu nehmen. Er laeuft
# nur auf ausdruecklichen Aufruf, nie nebenbei.
.PHONY: slice-mv

# Der Ablageort des Werkzeugs ist das Repo (tools/harness/).
SLICE_MV ?= tools/harness/slice-mv.sh

# Die Pfade, die der Verweis-Nachzug ausnimmt — REPO-POLITIK, keine Mechanik des
# Werkzeugs (die Begruendung traegt dessen Kopf). Die Vorgabe ist aus dem
# mitgelieferten Regelwerk abgeleitet; ein Repo mit einer anderen Politik setzt
# die Variable. Der Ort dafuer ist JEDE Make-Quelle, die der Aggregator einbindet
# — `?=` nimmt einen frueher gesetzten Wert an, ein spaeter gesetzter gewinnt —,
# und das Rezept unten reicht den Wert als Umgebung an das Werkzeug durch: eine
# blosse Zuweisung genuegt damit, ohne `export` und ohne Kommandozeilen-Praefix.
SLICE_MV_AUSGENOMMENE_PFADE ?= :!.harness/baseline :!docs/plan/adr

# Fehlt das Werkzeug, sagt dieses Ziel das und bricht ab: `make` haette sonst ein
# Rezept ueber einem Namen, den keine Datei traegt, und der Aufruf endet mit der
# Meldung des Interpreters statt mit der des Fragments.
slice-mv: ## Lifecycle-Wechsel eines Slice inkl. Verweis-Nachzug (SLICE=slice-<Kennung> TO=<open|next|in-progress|done>) — KEIN Gate
	@test -f "$(SLICE_MV)" || { echo "slice-mv: $(SLICE_MV) liegt nicht — dieses Ziel ruft das Skript im Repo."; exit 2; }
	@SLICE_MV_AUSGENOMMENE_PFADE='$(SLICE_MV_AUSGENOMMENE_PFADE)' bash "$(SLICE_MV)" "$(SLICE)" "$(TO)"
