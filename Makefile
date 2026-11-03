# Makefile — generiert von ai-harness-init (Aggregator). Die Gate-Belange leben als
# Fragmente unter harness/mk/*.mk; jedes haengt seine Checks an GATE_CHECKS. Der
# Gate-Nachweis (record-gates) laeuft strikt ZULETZT via Ordnungskante auf GATE_CHECKS
# — waehrend make -j die Checks parallelisiert; .NOTPARALLEL ist bewusst NICHT gewaehlt
# (das serialisierte das ganze Makefile). Sprach-agnostisch: ohne --lang matchen nur
# baseline/doc-gate/enforce, mit --lang zusaetzlich das Code-Gate-Fragment.

# Floor of the GNU make: the host and the CI runners carry at least this (LH-FA-DEV-011).
MAKE_MIN_MAJOR := 4
MAKE_MIN_MINOR := 3
MAKE_MIN_VERSION := $(MAKE_MIN_MAJOR).$(MAKE_MIN_MINOR)
MAKE_MAJOR := $(firstword $(subst ., ,$(MAKE_VERSION)))
MAKE_MINOR := $(word 2,$(subst ., ,$(MAKE_VERSION)))
# Mindest-Version der GNU make (LH-FA-DEV-011): der Lauf bricht ab, wenn die
# laufende make-Version (MAKE_VERSION, "major.minor") aelter als die
# Mindest-Version ist — make-eigener Guard, kein Stellenvergleich; er greift an
# der Pars-Zeit jedes Laufs, vor jeder anderen Logik. Stand des Messwerts:
# Host und Devcontainer-Image tragen 4.4.1, die CI-Runner (ubuntu-24.04) 4.3,
# die Mindest-Version ist 4.3.
ifeq (,$(shell if [ $(MAKE_MAJOR) -gt $(MAKE_MIN_MAJOR) ] || { [ $(MAKE_MAJOR) -eq $(MAKE_MIN_MAJOR) ] && [ $(MAKE_MINOR) -ge $(MAKE_MIN_MINOR) ]; }; then echo ok; fi))
$(error GNU make $(MAKE_MIN_VERSION) or newer is required, the running make reports $(MAKE_VERSION))
endif

GATE_CHECKS :=

.PHONY: gates help

# Gate-Fragmente je Belang (baseline/doc-gate/enforce + Sprach-Code-Gates) einbinden.
# Alphabetisch (baseline < doc-gate < enforce < <lang>); die Ordnungskante unten steht
# NACH beiden Includes und sieht GATE_CHECKS damit vollstaendig.
include harness/mk/*.mk

# Die Targets des Repos: repo.mk gehoert dem Repo, der Bootstrap legt sie
# nur an, wo sie fehlt. Sie steht NACH den Fragmenten, damit ihr GATE_CHECKS += die
# Ordnungskante unten erreicht und ihr = die ?=-Vorgaben der Fragmente ueberschreibt;
# -include, damit eine geloeschte repo.mk make nicht bricht.
-include repo.mk

help: ## Diese Hilfe
	@grep -hE '^[a-z0-9-]+:.*##' $(MAKEFILE_LIST) | sort | awk 'BEGIN{FS=":.*##"}{printf "  %-14s %s\n",$$1,$$2}'

# gates haengt allein an record-gates; record-gates haengt an ALLEN akkumulierten
# Checks — der Nachweis laeuft strikt nach den Checks (Ordnungskante), waehrend make
# -j die Checks parallel faehrt. Das record-gates-Rezept liefert harness/mk/enforce.mk.
gates: record-gates ## Alle Gates (Checks parallel, Nachweis zuletzt)
record-gates: $(GATE_CHECKS)
