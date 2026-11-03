# Slice CR: Pin-Politik der gepinnten Werkzeugkette — der fünfte Ort `go.mod` (Gleich-Pin-Politik, `pinsync`-Gatter)

**Lifecycle:** Der Zustand dieser Datei ist das Verzeichnis; Wechsel nur durch `git mv` (Modul 5 §Lifecycle). **Bezug:** [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) (Erweiterung), [`LH-FA-CONF-002`](../../../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) (nur lesend); Muster [`slice-cr-devcontainer-java-toolchain`](../done/slice-cr-devcontainer-java-toolchain.md). **Anlass:** Projektinhaber 2026-11-03 („die Nummern müssen zusammenpassen", „ja, schneide den d-check-bump mit pinsync danach — und auch ein gate für die nötigen tools und versionen"). **Verantwortlich:** pt9912. **Autor:** pt9912 (Rolle Planner, im Auftrag). **Datum:** 2026-11-03.

## 1. Ziel und Abgrenzung

**Ziel:** die Pin-Polität der zwei Go-Schlüssel erhält den **fünften Ort `go.mod`** (`go`-Direktive): die Direktive wird auf den Pin minor gehoben (`1.26.0` → `1.27.0`, die Toolchain trägt `1.27.3`) und fortan zieht jede Pin-Hebung die Direktive nach — mit einem `pinsync`-Wächter (d-check-Modul, byte-gleicher Wertvergleich über die Pin-Stellen) als Gatter der Kopplung. **Nicht im CR:** die Umsetzung (Bump, Config, Test-Ausdehnung — Umsetzungs-Slice), der JDK-Schlüssel (nicht Teil der Kopplung, [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) trägt die zwei Go-Schlüssel), der d-check-Bump selbst (eigener Vorgang).

## 2. Definition of Done

- [x] CR-Wortlaut je Stelle freigegeben (Vorschlags-Datei mit Alt/Neu)
- [x] Übernahme gleichlautend je Freigabe (Der Wortlaut-Abgleich je Stelle trägt §7 — Hybrid-Eröffnungsformel je Verifier V-1 ausgewiesen): [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette)-Absatz, Historie 0.34.0/entsprechend; Abgleich je Stelle in §7
- [x] Gates, Reports, Closure je Muster der CR-Slices

## 3. Plan (vor Code)

| Datei | Änderungs-Art | Begründung |
|---|---|---|
| `cr-vorschlag-lh-dev-011-pin-inventur.md` | neu (geliefert 2026-11-03) | Inventar je Werkzeugkette samt GNU-make-Erweiterung; Status freigegeben wie vorgelegt |
| `spec/lastenheft.md` | update (geliefert 2026-11-03) | Pin-Inventar-Absatz unter [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette), Versionsfeld 0.34.0, Historie-Zeile |
| `d-check.mk` | update (Umsetzungs-Diff) | Bump v0.91.0 samt verifiziertem Digest, Kopf-Kommentar in Zusage-Form (F-5: ohne Chronik-Wort; das Datum im Commit), die sechs `--disable pinsync`-Zusätze je `--print-mk` der Version |
| `.d-check.yml` | update (Umsetzungs-Diff) | Modul `pinsync` in `modules`, drei Pin-Gruppen (`go` fünf Quellen, `golangci-lint` vier, `devcontainer-cli` zwei) |
| `Makefile` | update (Umsetzungs-Diff) | GNU-make-Guard am Kopf, Mindest-Version 4.3 |
| `go.mod` | update (Umsetzungs-Diff) | Direktive auf 1.27.0 (Pin-Minor), Require-Layout-Normierung des Tidy-Laufs |
| `spec/spezifikation.md` | update (Umsetzungs-Diff) | [`SPEC-006`](../../../../spec/spezifikation.md#spec-006--build-args-des-dockerfiles-und-pin-politik) Abschnitte Pin-Inventar/Pin-Gleichheit samt Beleg-Tabelle, [`SPEC-016`](../../../../spec/spezifikation.md#spec-016--go-toolchain-mindestversion-und-dockerfile-pin) Direktive, Historie |
| `docs/maintainer/quality.md` | update (Umsetzungs-Diff) | Abschnitt 7 (Werkzeug-Pins und GNU-make-Mindest-Version, Gate-Verhalten) |
| `harness/README.md` | update (Umsetzungs-Diff) | docs-check-Zeile des Gate-Index nennt die Pin-Gleichheit |
| `CHANGELOG.md` | update (Umsetzungs-Diff) | Eintrag unter `[Unreleased]` |
| `docs/messungen/slice-harness-dcheck-pinsync/` | neu (Umsetzungs-Diff) | Beleg-Lauf je Gatter mit Mismatch-Gegenproben und make-Guard-Negativprobe |
| dieser Slice (diese Datei), Roadmap | update | Plan-Tabelle und Umsetzungs-Zeilen; Lifecycle je `make slice-mv`; **nach §3.10 berichtigt (Review):** von der Roadmap trägt dieser Diff **keine Zeile** — die Bestand-/Marker-Züge laufen je nachbarkommenden Zügen; die §3.10-Abgleiche laufen je Fix-und-Closure-Runde erneut | bei Closure |
| Tests der Vier-Stellen-Kopplung | kein Diff in diesem Slice | `TestToolchainPinPlacesEqual` bleibt (liest die vier Code-Stellen, [`SPEC-054`](../../../../spec/spezifikation.md)); die Pin-Gruppen im Gatter prüfen die Quellen-Zeilen je Pin |
| die drei beruehrten Test-Kommentar-Dateien (internal/adapter/driven/docker/engine_test.go, docker_test.go, git_test.go) | update | die Chronik-Nennungen der Basis-Images auf den Ist-Stand (F-7/INFO je Fix; die Kommentar-Zeilen je gelesen) |
| docs/reviews/ — [Review-Report](../../../reviews/2026-11-03-slice-harness-dcheck-pinsync.md) und [Verify-Report](../../../reviews/2026-11-03-slice-harness-dcheck-pinsync-verify.md) | neu | die Berichte der zwei Rollen (Nachtrag je Verifier V-2 — die Berichts-Dateien-Familie) |

## 4. Trigger und Rückführungen

**Start:** `Verantwortlich:` gesetzt; kein anderer Slice in `in-progress/`. **Rückführungen:** `next` (zu groß — etwa wenn die JDK-Stelle auch in die Kopplung soll); `open` (Freigabe verbleibt aus).

## 5–8

Closure-Trigger: Wortlaut wortgleich je Stelle + Gates + Reports + Notiz. Risiken: Wortlaut-Abgleiche (Muster der CR-Slices), Zusage-Messbasis (pinsync prüft Gleichheit, keine Aktualität — Grenze im Wortlaut), Beleg-Aufträge, §3.10. §7 bei Closure gefüllt. Sub-Areas: DOK Greenfield.

**Der freizugebende Wortlaut (Stelle b, freigegeben 2026-11-03 wie vorgelegt samt GNU-make-Erweiterung, übernommen ins Lastenheft 0.34.0 — der Minor-Vorbehalt der Direktive trägt die freigegebene Hebung `1.27.0` gegen den Pin `1.27.3` als Definition des „denselben Pin-Werts“):** „Die `go`-Direktive der `go.mod` dieses Repos ist der fünfte Ort der Pin-Kopplung: die Direktive trägt denselben Wert wie die `GO_VERSION`-Stellen (der Minor gleich, die Patch-Stelle der Toolchain bleibt in den `GO_VERSION`-Stellen); jede Hebung des Pins zieht die Direktive nach. Ein Wächter (`pinsync`, gleiches Byte-Verhalten) prüft die Gleichheit je Gatter-Lauf."

## 7. Closure-Notiz

### Closure (Planner, 2026-11-03, nach Umsetzung , Review-Fix  und [Verification](../../../reviews/2026-11-03-slice-harness-dcheck-pinsync-verify.md))

**Evidence:** Review-Report [2026-11-03-slice-harness-dcheck-pinsync](../../../reviews/2026-11-03-slice-harness-dcheck-pinsync.md) (0 HIGH/2 MEDIUM/4 LOW/2 INFO je berichtigt) — Digest je eigener Verifikation, Positiv- und Mismatch-Läufe je eigener Proben · Verify-Report […-verify](../../../reviews/2026-11-03-slice-harness-dcheck-pinsync-verify.md) (DoD-Zeilen 1 bis 6 belegt erfüllt; V-1/V-2/V-3 je in der Berichtigung vor der Closure getragen; zwei unabhängige Mismatch-Belegspuren).

**Wortlaut-Abgleich je Stelle (V-1):** die Inventur-Aufsätze des Lastenhefts sind eine **Hybrid-Eröffnungsformel** gegenüber dem CR-Zitat — nicht wortgleich, inhaltsgleich je Verifier-Nachzug; beide freigegebenen Berichtigungen (Mindest-Toolchain 1.27, Minor-Klausel) sind je ihrer Freigabe getragen; alle weiteren Stellen der Übernahme je Vorschlags-Alt/Neu gelesen.

**Läufe:** baseline-verify: v6.18.1 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)
docker run --rm --network none -v "/Development/u-boot:/repo:ro" ghcr.io/pt9912/d-check@sha256:5d2f650d102c1d8224bb26930ec96d237506c07612fb8e0dbcdc42291702c57f
docker build  --build-arg GO_VERSION=1.27.3 --build-arg GOLANGCI_LINT_VERSION=v2.14.0 --build-arg UBOOT_VERSION=0.21.0-dev --no-cache-filter lint --target lint -t u-boot:lint .
docker build  --build-arg GO_VERSION=1.27.3 --build-arg GOLANGCI_LINT_VERSION=v2.14.0 --build-arg UBOOT_VERSION=0.21.0-dev --no-cache-filter test --target test -t u-boot:test .
docker build  --build-arg GO_VERSION=1.27.3 --build-arg GOLANGCI_LINT_VERSION=v2.14.0 --build-arg UBOOT_VERSION=0.21.0-dev --no-cache-filter coverage     --target coverage     --build-arg COVERAGE_THRESHOLD=90     -t u-boot:coverage .
ablage-pruefung: OK
ok    Grundzustand ist gruen
ok    Skript unter docs/messungen/
ok    Dockerfile unter docs/messungen/
ok    Perl-Sonde unter docs/
ok    Skript im Planning-Baum
ok    Skript im Register
ok    Rohdaten im Register
ok    Beobachtung mit fremdem Dateinamen
ok    Datei neben der Beobachtung
ok    Nicht-md-Datei im Planning-Baum
ok    .gitkeep in einem vierten Verzeichnis
ok    Skript ohne Endung mit Ausfuehrungsbit
ok    Skript ohne Endung, Shebang, ausserhalb von messungen
ok    Grossgeschriebenes Skript
ok    Grossgeschriebenes Skript ausserhalb von messungen
ok    dockerfile kleingeschrieben
ok    PowerShell-Skript
ok    json unter messungen
ok    png unter messungen
ok    grossgeschriebene Endung TXT
ok    versteckte Datei unter messungen
ok    Datei ohne Endung unter messungen
ok    Symlink unter messungen
ok    md mit Ausfuehrungsbit
ok    png ausserhalb von messungen und planning ist erlaubt
ok    ungetrackte Datei zaehlt nicht
ablage-pruefung-test: 26 Faelle, 0 Fehlschlaege Exit 0 je vor- und nach-commit (Stempel je HEAD;  =  am Schluss-Stand je Verifier-Nachrechnung des ), docker run --rm --network none -v "/Development/u-boot:/repo:ro" ghcr.io/pt9912/d-check@sha256:5d2f650d102c1d8224bb26930ec96d237506c07612fb8e0dbcdc42291702c57f 0 Befunde je Runde (1084/1085 Dateien) inkl. dem live pinsync-Wächter („3 Pin(s) geprüft, 11 Quelle(n) gelesen"), modus-pruefung: OK (8 targets in both modes; build, run, image-scan, packages, docs-check in docker; BUILD_MODE, PKG and RUN rules) OK (vorsorglich, 8 Targets), Guard-Negativprobe je Rohdaten (MAKE_VERSION=4.2 → (error)), Mismatch-Gegenproben je zwei unabhängige Belegspuren (Implement-/Reviewer-Proben); nicht gelaufen: verify-depguard/test-docker (keine Import-/Bau-Pfad-Berührung; die go.mod-Hebung deckt der compile-Stage des make test).

**Risiken je Ausgang:** Wortlaut-Abgleich — der Hybrid-Aufsatz je V-1 ausgeweisen (oben); Zusage-Messbasis — die Minor-Fang-Auslegung je CR-Parenthese und freigegebener Fix-Klausel (kein Aktualitäts-Versprechen je Wortlaut „Gleichheit, keine Aktualität"); Beleg-Aufträge — die Test-Namen je Zeile geliefert (Vier-Stellen-Kopplungstest bleibt beim Muster); §3.10 — die drei Test-Dateien und beide Berichts-Dateien je Plan-Zeile nachgetragen (V-2).

**Was ging anders als geplant:** die drei Chronik-Nennungen der Test-Kommentare (golang:1.26.3, v2.12.2) traten beim Review auf (F-7/INFO) und trugen die Ist-Stand-Berichtigung; die (-alpine)-Lese-Variante des Kommentars überstand die Erst-Berichtigung und kam erst vom Verifier (V-3 — die Kante des halb berichtigten Kommentars); der Vertrags-Aufsatz ist Hybrid statt wortgleich (V-1).

**Steering-Loop-Eintrag** (gezählt, nicht verkörpert; kein Eintrag trägt das Feld liegt in): Lehre — der pinsync-Wächter deckt die Markdown-Dateien, die Go-Tests die Code-Stellen; die zwei Prüffamilien überlappen an keiner Stelle, der Abgleich „doppelte Deckung" ist je Verweis getragen, kein Sensor-Zuwachs nötig. Kein Feld liegt in. Benannte Spec-Lücke: keine.

**Beobachtungs-Register**: keine Beobachtung angefallen in diesem Vorgang; der Beleg der Prüflinsen ([BEO-DOK/plan-aussage-ueber-code](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) und [BEO-DOK/zusage-weiter-als-messung](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)) bleibt je „kein Fund" — die zustandsmarker-Instanz der Roadmap-Zeile (V-4 INFO, „d-check-Bump folgt je Umsetzung" abgelaufen) ist in der Berichtigung getragen und trägt keine neue Zählung (die Zeile ist je Marker-Regel der Roadmap, nicht des Plans).

**Carveouts**: keiner. **Folge-Slice:** der Übergang [slice-v1-bau-nur-devcontainer](../open/slice-v1-bau-nur-devcontainer.md) liegt in open/ (drei Entscheidungen bejaht 2026-11-03, Start je WIP); der Registry-Tarball-Weg bleibt Grenze. **Weiter offen:** nichts aus der DoD.

**Drei Paarungen:** Anker — kein Eintrag mit Zielort-Feld; nichts verkörpert (die Gatter-Läufe tragen die Regeln, die der Inventur-Absatz der Lastenheft nennt) · Folge-Slice — der devcontainer-Übergang ist eine Datei in open/ (Verzeichnis-Lesung) · Register — keine neue Kennung („keine Beobachtung", die Lesung oben).

**Verifikation** (Läufe nach Closure-Inhalt, Move und Marker; Docker-Modus): make gates Exit 0 am 2026-11-03 (record-gates-Stempel = HEAD nach dem Nachzug). Der Lauf gilt für den Stand vor dieser Zeile; die Zeile ändert nur Prosa.

*(Wird bei der Closure gefüllt.)*