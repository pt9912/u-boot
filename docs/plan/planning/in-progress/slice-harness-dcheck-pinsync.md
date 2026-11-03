# Slice CR: Pin-Politik der gepinnten Werkzeugkette — der fünfte Ort `go.mod` (Gleich-Pin-Politik, `pinsync`-Gatter)

**Lifecycle:** Der Zustand dieser Datei ist das Verzeichnis; Wechsel nur durch `git mv` (Modul 5 §Lifecycle). **Bezug:** [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) (Erweiterung), [`LH-FA-CONF-002`](../../../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) (nur lesend); Muster [`slice-cr-devcontainer-java-toolchain`](../done/slice-cr-devcontainer-java-toolchain.md). **Anlass:** Projektinhaber 2026-11-03 („die Nummern müssen zusammenpassen", „ja, schneide den d-check-bump mit pinsync danach — und auch ein gate für die nötigen tools und versionen"). **Verantwortlich:** pt9912. **Autor:** pt9912 (Rolle Planner, im Auftrag). **Datum:** 2026-11-03.

## 1. Ziel und Abgrenzung

**Ziel:** die Pin-Polität der zwei Go-Schlüssel erhält den **fünften Ort `go.mod`** (`go`-Direktive): die Direktive wird auf den Pin minor gehoben (`1.26.0` → `1.27.0`, die Toolchain trägt `1.27.3`) und fortan zieht jede Pin-Hebung die Direktive nach — mit einem `pinsync`-Wächter (d-check-Modul, byte-gleicher Wertvergleich über die Pin-Stellen) als Gatter der Kopplung. **Nicht im CR:** die Umsetzung (Bump, Config, Test-Ausdehnung — Umsetzungs-Slice), der JDK-Schlüssel (nicht Teil der Kopplung, [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) trägt die zwei Go-Schlüssel), der d-check-Bump selbst (eigener Vorgang).

## 2. Definition of Done

- [ ] CR-Wortlaut je Stelle freigegeben (Vorschlags-Datei mit Alt/Neu)
- [ ] Übernahme wortgleich: [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette)-Absatz, Historie 0.34.0/entsprechend; Abgleich je Stelle in §7
- [ ] Gates, Reports, Closure je Muster der CR-Slices

## 3. Plan (vor Code)

| Datei | Änderungs-Art | Begründung |
|---|---|---|
| `cr-vorschlag-lh-dev-011-pin-inventur.md` | neu (geliefert 2026-11-03) | Inventar je Werkzeugkette samt GNU-make-Erweiterung; Status freigegeben wie vorgelegt |
| `spec/lastenheft.md` | update (geliefert 2026-11-03) | Pin-Inventar-Absatz unter [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette), Versionsfeld 0.34.0, Historie-Zeile |
| `d-check.mk` | update (Umsetzungs-Diff) | Bump v0.91.0 samt verifiziertem Digest, Bump-Datums-Kommentar, die sechs `--disable pinsync`-Zusätze je `--print-mk` der Version |
| `.d-check.yml` | update (Umsetzungs-Diff) | Modul `pinsync` in `modules`, drei Pin-Gruppen (`go` fünf Quellen, `golangci-lint` vier, `devcontainer-cli` zwei) |
| `Makefile` | update (Umsetzungs-Diff) | GNU-make-Guard am Kopf, Mindest-Version 4.3 |
| `go.mod` | update (Umsetzungs-Diff) | Direktive auf 1.27.0 (Pin-Minor), Require-Layout-Normierung des Tidy-Laufs |
| `spec/spezifikation.md` | update (Umsetzungs-Diff) | [`SPEC-006`](../../../../spec/spezifikation.md#spec-006--build-args-des-dockerfiles-und-pin-politik) Abschnitte Pin-Inventar/Pin-Gleichheit samt Beleg-Tabelle, [`SPEC-016`](../../../../spec/spezifikation.md#spec-016--go-toolchain-mindestversion-und-dockerfile-pin) Direktive, Historie |
| `docs/maintainer/quality.md` | update (Umsetzungs-Diff) | Abschnitt 7 (Werkzeug-Pins und GNU-make-Mindest-Version, Gate-Verhalten) |
| `harness/README.md` | update (Umsetzungs-Diff) | docs-check-Zeile des Gate-Index nennt die Pin-Gleichheit |
| `CHANGELOG.md` | update (Umsetzungs-Diff) | Eintrag unter `[Unreleased]` |
| `docs/messungen/slice-harness-dcheck-pinsync/` | neu (Umsetzungs-Diff) | Beleg-Lauf je Gatter mit Mismatch-Gegenproben und make-Guard-Negativprobe |
| dieser Slice (diese Datei), Roadmap | update | Plan-Tabelle und Umsetzungs-Zeilen; Lifecycle je `make slice-mv` bei Closure |
| Tests der Vier-Stellen-Kopplung | kein Diff in diesem Slice | `TestToolchainPinPlacesEqual` bleibt (liest die vier Code-Stellen, [`SPEC-054`](../../../../spec/spezifikation.md)); die Pin-Gruppen im Gatter prüfen die Quellen-Zeilen je Pin |

## 4. Trigger und Rückführungen

**Start:** `Verantwortlich:` gesetzt; kein anderer Slice in `in-progress/`. **Rückführungen:** `next` (zu groß — etwa wenn die JDK-Stelle auch in die Kopplung soll); `open` (Freigabe verbleibt aus).

## 5–8

Closure-Trigger: Wortlaut wortgleich je Stelle + Gates + Reports + Notiz. Risiken: Wortlaut-Abgleiche (Muster der CR-Slices), Zusage-Messbasis (pinsync prüft Gleichheit, keine Aktualität — Grenze im Wortlaut), Beleg-Aufträge, §3.10. §7 bei Closure gefüllt. Sub-Areas: DOK Greenfield.

**Der freizugebende Wortlaut (Stelle b, freigegeben 2026-11-03 wie vorgelegt samt GNU-make-Erweiterung, übernommen ins Lastenheft 0.34.0 — der Minor-Vorbehalt der Direktive trägt die freigegebene Hebung `1.27.0` gegen den Pin `1.27.3` als Definition des „denselben Pin-Werts“):** „Die `go`-Direktive der `go.mod` dieses Repos ist der fünfte Ort der Pin-Kopplung: die Direktive trägt denselben Wert wie die `GO_VERSION`-Stellen (der Minor gleich, die Patch-Stelle der Toolchain bleibt in den `GO_VERSION`-Stellen); jede Hebung des Pins zieht die Direktive nach. Ein Wächter (`pinsync`, gleiches Byte-Verhalten) prüft die Gleichheit je Gatter-Lauf."

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*