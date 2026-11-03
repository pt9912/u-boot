# Review-Report: slice-harness-dcheck-pinsync — 2026-11-03

**Review-Art:** Code-Review (Umsetzungs-Diff) — geprüft wird gegen Plan, Vertrags-Carrier und Hard Rules (Maintainability); DoD/Closure ist Verifier-Aufgabe und hier nicht geprüft.

**Gegenstand:** Commit `a75bdb25` („Umsetzung der Pin-Inventur", 16 Dateien) des Vorgangs [`slice-harness-dcheck-pinsync`](../plan/planning/done/slice-harness-dcheck-pinsync.md); Vertrags-Norm [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) 0.34.0, freigegebener CR `cr-vorschlag-lh-dev-011-pin-inventur` („wie vorgelegt samt GNU-make-Erweiterung").

**Skill:** `.harness/skills/reviewer.md` @ Repo-Stand a75bdb25 ·
**Modell:** GLM (Claude-Agenten-Harness) · **Datum:** 2026-11-03

**Eingangs-Kontext:**

- Slice-Plan [`slice-harness-dcheck-pinsync`](../plan/planning/done/slice-harness-dcheck-pinsync.md) (§3 Plan-Tabelle nach dem Bump des Commits)
- CR `cr-vorschlag-lh-dev-011-pin-inventur` (freigegeben 2026-11-03) und die Slice-Stelle „freizugebender Wortlaut"
- [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) (Pin-Politik, Pin-Inventur-Absatz, Historie 0.34.0), [`LH-OPEN-001`](../../spec/lastenheft.md#lh-open-001--implementierungssprache-entschieden), [`SPEC-006`](../../spec/spezifikation.md#spec-006--build-args-des-dockerfiles-und-pin-politik), [`SPEC-016`](../../spec/spezifikation.md#spec-016--go-toolchain-mindestversion-und-dockerfile-pin), [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)
- `AGENTS.md` §3.3, §3.7, §3.9, §3.10, §5; `harness/conventions.md` (`MR-000`, `MR-012`, `MR-007`/`MR-009`-Form); `.harness/baseline/v6.18.1` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill
- d-check-Modul-Handbuch §4.27 „Modul `pinsync`" (gelesen aus dem gepullten v0.91.0-Image, hermetisch); `d-check --print-mk` von v0.91.0

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | [`LH-OPEN-001`](../../spec/lastenheft.md#lh-open-001--implementierungssprache-entschieden) (Entscheidungs-Record, Rank 1) nennt „Mindest-Toolchain: Go 1.26 oder neuer"; der Diff hebt die `go`-Direktive auf `1.27.0` und formuliert [`SPEC-016`](../../spec/spezifikation.md#spec-016--go-toolchain-mindestversion-und-dockerfile-pin) (Rank 2) mit „Go 1.27 oder neuer". Ein Go-1.26-Werkzeugkasten mit `GOTOOLCHAIN=off` bricht am Bau; die Rank-1-Aussage zur Mindestversion wird damit über ohne Berührung falsch, während Rank 2 sie verschärft (Spec „erweitert" hier statt nur zu präzisieren). Keine Code-Falschheit, sondern eine Vertrags-Spanne, die nur über eine freigegebene Berichtigung des Rank-1-Textes (oder eine Bindung der Mindestversion an den Pin-Minor) aufzulösen ist. | Source Precedence (`MR-001`); [`LH-OPEN-001`](../../spec/lastenheft.md#lh-open-001--implementierungssprache-entschieden) | `spec/lastenheft.md` · „Mindest-Toolchain: Go 1.26 oder neuer." | Review-only | Mindestversion-Vertrags-Spanne |
| F-2 | MEDIUM | Der Pin-Inventur-Absatz sagt, die Direktive „trägt denselben Pin-Wert wie die `GO_VERSION`-Stellen"; der Umsetzungs-Stand trägt `1.27.0` neben den Pins `1.27.3` — gleich nur je Minor. Die Definition des „denselben Pin-Werts" als Minor-Gleichheit (Minor-Vorbehalt) steht im freigegebenen Wortlaut des Slices und im CR (Gruppen-Tabelle, „Offen zur Freigabe"), nicht aber im übernommenen Rank-1-Text; der Slice behauptet, sie sei „übernommen ins Lastenheft 0.34.0" — sie ist es dort wortlautlich nicht. Eine nur-lastenheftlesende Prüfung liest einen Vertragsbruch ins aktuelle Stand (oder „berichtigt" `1.27.0` → `1.27.3` gegen die freigegebene Form). | [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) 0.34.0; `MR-001` | `spec/lastenheft.md` · „die Direktive trägt denselben Pin-Wert wie die `GO_VERSION`-Stellen" | Review-only | Pin-Wert-Definition ohne Rank-1-Träger |
| F-3 | LOW | Der neu eingefügte GNU-make-Guard (Definitionen plus `$(error)`-Abbruch) steht zwischen dem Kopf-Dok-Kommentar der Makefile (Aggregator/Gate-Mechanik) und `GATE_CHECKS :=`, dem von ihm beschriebenen Subjekt — genau das §3.9(c)-Muster, das der Kommentar-Selbstcheck vor dem Commit abfängt. | `AGENTS.md` §3.9 (c) | `Makefile` · „GATE_CHECKS :=" | Review-only | Helfer zwischen Doc-Kommentar und Deklaration |
| F-4 | LOW | Die neue Grenze „Shell-Werkzeuge … ihre Versionen kommen mit dem Basis-Image des Pins" setzt je Werkzeugkette ein gepinntes Basis-Image voraus; die Dev-Container-Basis `mcr.microsoft.com/devcontainers/base:debian` trägt keinen Pin (beweglicher Distributions-Tag, Pin-Entscheidung ausgesetzt beim Projektinhaber), und der CR führte die Basis-Image-Gruppe ausdrücklich nur über die `golang:*`-Tags. Für den Entwicklungs-Verlauf (Shell-Werkzeuge im erzeugten und im eigenen Dev-Container) behauptet die Grenze so eine Bestimmtheit, die der Stand nicht trägt. | `AGENTS.md` §3.10 (Spec-Bedingung ohne Test-Beleg); CR-Scope | `spec/lastenheft.md` · „ihre Versionen kommen mit dem Basis-Image des Pins" | Review-only | Grenze trifft ungepinnte Basis nicht |
| F-5 | LOW | Der Bump-Kommentar in `d-check.mk` trägt Chronik neben der tragenden Zusage: „Image-Bump 2026-11-03 auf v0.91.0" und „neu ist das Modul pinsync" erzählen den Wechsel; die Zusage der Datei (Ein-Modul-Targets tragen `--disable pinsync` je `--print-mk` der Version) steht daneben in einem Block, deren Datums- und Neu-Teil beim nächsten Bump veraltet statt berichtigt erscheint. | `AGENTS.md` §3.7 | `d-check.mk` · „Image-Bump 2026-11-03 auf v0.91.0" | Review-only | Chronik-Kommentar |
| F-6 | LOW | Die Plan-Zeile „dieser Slice (diese Datei), Roadmap | update" verspricht eine Roadmap-Berührung; der Diff berührt nur den Slice-Plan, und die Vorgangs-Zeile der Roadmap sagt „d-check-Bump folgt je Umsetzung", obwohl der Bump im selben Commit umgesetzt ist. Der im Commit behauptete §3.10-Abgleich hat eine Plan-Zelle offen gelassen (Eintrag ohne Diff) und die Zeile steht auf einem inzwischen abgelaufenen Zustandsstand. | `AGENTS.md` §3.10 | [`slice-harness-dcheck-pinsync`](../plan/planning/done/slice-harness-dcheck-pinsync.md) · „dieser Slice (diese Datei), Roadmap | update" | Review-only | Plan-Zeile ohne Diff |
| F-7 | INFO | Drei Testkommentare behaupten den Test-Image-Stand `golang:1.26.3` beziehungsweise `golangci/golangci-lint:v2.12.2-alpine` als Jetzt-Stand (`make test`-Pfad, CI, Gate-Image); der Pin trägt `1.27.3` beziehungsweise `v2.14.0` seit 53041bdb. Vor-bestehend, vom Diff nicht berührt — Nachfolge-Pflege, keine Verstoß-Zurechnung zu diesem Commit. | `AGENTS.md` §3.7 | `internal/adapter/driven/git/git_test.go` · „green on `golang:1.26.3`"; `internal/adapter/driven/docker/engine_test.go` · „CI runs inside `golang:1.26.3`"; `internal/adapter/driven/docker/docker_test.go` · „gates image runs in golang:1.26.3" | Review-only | veraltete Test-Umgebungs-Kommentare |
| F-8 | INFO | Die Abgrenzung „make-eigener Guard, kein Stellenvergleich" (Kommentar, [`SPEC-006`](../../spec/spezifikation.md#spec-006--build-args-des-dockerfiles-und-pin-politik), quality.md §7) ist nur mit dem CR-Kontext lesbar, wo die Stelle „kein pinsync-Stellenvergleich" heißt; der CR trägt den Vergleich als Guard je Mindest-Version, der Ausschluss gilt dem Pin-Stellen-Gleichlauf, nicht dem Zahlvergleich am Kopf der Datei. | `AGENTS.md` §3.7 (Abgrenzung lesbar tragen) | `Makefile` · „make-eigener Guard, kein Stellenvergleich" | Review-only | Abgrenzung nur mit CR-Kontext lesbar |

## Offene Punkte des Implementers — Bewertung

| Punkt | Fundstelle | Verdict |
|---|---|---|
| (a) Vertrags-Spanne [`LH-OPEN-001`](../../spec/lastenheft.md#lh-open-001--implementierungssprache-entschieden) („Go 1.26 oder neuer") | `spec/lastenheft.md` · „Mindest-Toolchain: Go 1.26 oder neuer." | **F-1, MEDIUM.** Kein wörtlicher Code-Gegenbeweis (1.27 erfüllt „1.26 oder neuer" als Untergrenze der Entscheidung), aber der Diff verschärft die Mindestversion über Rank 2 und lässt Rank 1 falsch stehen — CR-/Freigabe-Fälligkeit, nicht still im Code zu lösen. Die „Pin-Hebung … Routine"-Klausel der Stelle trägt den ARG-Pin, nicht die Sprach-Direktive. |
| (b) Chronik-Nennungen `golang:1.26.3` | drei Test-Dateien (nicht zwei) | **F-7, INFO.** §3.7-Verstoß bejaht, aber vor-bestehend (seit 53041bdb, Bump auf 1.27.3) und vom Diff nicht berührt; der §3.9(b)-Trigger („der Diff ändert oder umbenennt") greift für diese Werte nicht. |
| (c) Dev-Container-Basis untagged | `mcr.microsoft.com/devcontainers/base:debian` | **F-4, LOW.** Die Grenze trägt die Form nicht: der CR-Fallen-Bereich „Basis-Image-Gruppe" umfasste nur `golang:*`-Tags; der verallgemeinerte Rank-1-Satz „des Pins" trifft die ungepinnte Basis nicht (und `gcr.io/distroless/static-debian12:nonroot` der Laufzeit-Stufe ebenso wenig); die Pin-Entscheidung steht als N4 beim Projektinhaber offen (Trail `docs/reviews/2026-10-22-*`). |
| (d) tidy-Normierung außerhalb der Direktiven-Hebung | `go.mod` · Require-Layout | **Ausgewisene.** Mechanische Folge des `go mod tidy`-Laufs am Docker-Stage, im Plan (Zeile `go.mod`) und in der Spec-Historie deklariert, Byte-Stabilität belegt (`roh-05`); die Alt-Form markierte direkte Abhängigkeiten (cobra, pflag, yaml) als `// indirect` — die Normierung korrigiert sie als direkte. |
| Minor-Fangauslegung der `go`-Gruppe | `.d-check.yml` · `^go (\d+\.\d+)\.` gegen `(\d+\.\d+)\.\d+$` je `GO_VERSION`-Stelle | **Ausgewisene.** Getragen vom freigegebenen CR (Gruppen-Tabelle: „`go.mod` auf `1.27.0` gehoben (Minor = Pin `1.27.3`)"; „Offen zur Freigabe: … Hebung auf `1.27.0` (Minor-Gleichheit)"), vom CR-Status („freigegeben wie vorgelegt"), vom Modul-Handbuch („Textgleichheit … byte-genau") und vom gemessenen Wächter-Verhalten; die Rand-1-Lücke der Definition ist F-2. |
| Mindest-Version 4.3 (statt 4.4.1) | `Makefile` · `MAKE_MIN_VERSION` | **Ausgewisene.** Freigabe verlangte den Guard ohne Zahl; der gemessene Träger-Basis: CI-Runner ubuntu-24.04 tragen 4.3 (Distribution), Host und Dev-Container-Image 4.4.1 — 4.3 ist die höchste von allen realen Trägern getragene Stufe, belegt in Kommentar, [`SPEC-006`](../../spec/spezifikation.md#spec-006--build-args-des-dockerfiles-und-pin-politik), quality.md §7 und Negativprobe (`roh-04`, eigener Gegenlauf bestätigt). |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Digest des d-check-Bumps (`sha256:5d2f650d…c57f`) | geprüft, ohne Befund — Release-Body des Tags v0.91.0 (GitHub-API, eigener Abruf) und `RepoDigests` des gezogenen Images tragen ihn; `d-check.mk` zitiert ihn |
| Sechs `--disable pinsync`-Zusätze in `d-check.mk` gegen `--print-mk` des v0.91.0-Images | geprüft, ohne Befund — die sechs Ein-Modul-Rezepte (vcs, commits, planning, tracked, targets, structure) stimmen je Disable-Liste wörtlich mit der emitted Form überein |
| `.d-check.yml`-Gruppenkonfiguration gegen das Modul-Handbuch (§4.27) | geprüft, ohne Befund — je Pin ≥ 2 Quellen, je Ausdruck genau eine (unbenannte) Fanggruppe, RE2, zeilenankernd; `devcontainer-cli` 2 Quellen, `golangci-lint` 4, `go` 5 |
| Positiver Gate-Lauf (geprüfter Baum) | geprüft, ohne Befund — „3 Pin(s) geprüft, 11 Quelle(n) gelesen", 0 Befunde, Exit 0 |
| Eigene Mismatch-Gegenprobe (Kopie, `ARG GO_VERSION` 1.27.3 → 1.26.9) | geprüft, ohne Befund — vier `pinsync-mismatch`-Zeilen (je abweichende Quelle), Meldung nennt alle Quellen-Werte, Exit 1; bestätigt `roh-02`/`roh-03` unabhängig |
| GNU-make-Guard (eigenes Probe-Lauf) | geprüft, ohne Befund — `make -n help` parst grün, `make MAKE_VERSION=4.2 -n help` bricht mit `$(error)` und Exit 2 vor jeder Logik ab; bestätigt `roh-04` |
| make-Guard-Arithmetik | geprüft, ohne Befund — Zahlvergleich je Major/Minor („4.10" würde string-wise fälschlich unter „4.3" liegen), `4.3.1`/`4.4.1` richtig gelesen; Kommandozeilen-Setzung sticht (Probe-Fahrbahn für CI) |
| Hilfs-Image-Tag-Ableitung des Spec/Config-Kommentars („keine eigene Text-Stelle") | geprüft, ohne Befund — Code-Stelle gelesen (`devcontainercli/Dockerfile` „the image tag derives from the CLI version and the hash of the three files"; Launcher `helperImageTag`); damit trägt die Ersetzung der CR-Quelle „CLI-Image-Tag-Nennung" durch die Resolved-URL im Lock |
| CHANGELOG-IDs und Wortlaut | geprüft, ohne Befund — `AGENTS.md` §5 Regel 6 setzt Kennungen ausdrücklich im CHANGELOG; Form (Zusagen/„Promised/Not promised") folgt dem Muster der Vorgänger-Zeilen |
| quality.md §7 (Nummerierung nach §6, Pfade der Verweise, Träger-Trennung: „Die Pin-Werte selbst stehen an ihren Stellen, nicht in dieser Datei") | geprüft, ohne Befund |
| Beleg-Tabelle [`SPEC-006`](../../spec/spezifikation.md#spec-006--build-args-des-dockerfiles-und-pin-politik) (Referenzen auf `TestToolchainPinPlacesEqual` in der genannten Test-Datei, Gate-Läufe, Guard-Probe) | geprüft, ohne Befund — Test existiert (`generate_toolchain_test.go`), Pin-Stellen-Werte (`1.27.3`, `v2.14.0`, `0.89.0`) am Code gelesen |
| Plan §3 gegen `git show --stat` (16 Dateien) | eine Abweichung — die Roadmap-Zelle (F-6); alle übrigen Diff-Dateien haben Plan-Zeilen, keine Plan-Datei fehlt im Diff |
| §3.9(a) Kommentare („früher/no longer"-Muster) im Gesamtdiff | geprüft, ohne Befund |
| `spec/spezifikation.md` Referenz-Richtung (Spec verweist auf Lastenheft aufwärts; kein Slice-/Roadmap-Verweis neu) und Decken-Regel | geprüft, ohne Befund |

**Nicht geprüft:**

- `make gates` in voller Breite — Handoff-Nachweis der [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md)-Form ist Implementer-/Verifier-Sache; der Reviewer fährt den engsten Sensor (docs-check samt Wächter- und Guard-Proben, hier eigenständig gefahren).
- `make test` / `make test-docker` — kein Produktcode im Diff (nur Harness-Konfiguration und Doku).
- `make verify-depguard` — keine Importregeln berührt (Sensor nicht anwendbar).
- `make image-scan` / `make govulncheck` — keine Runtime-Image- oder Go-Abhängigkeits-Änderung im Diff.
- Voller Inhalt von `d-check` CHANGELOG §[0.91.0] — nur die Digest-Zeile des Release-Bodies gelesen (die tragende Aussage ist die Pin-Gleichheit am eigenen Gate-Lauf, nicht die Release-Prosa).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 4 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Mindestversion-Vertrags-Spanne · Pin-Wert-Definition ohne Rank-1-Träger · Helfer zwischen Doc-Kommentar und Deklaration · Grenze trifft ungepinnte Basis nicht · Chronik-Kommentar · Plan-Zeile ohne Diff · veraltete Test-Umgebungs-Kommentare · Abgrenzung nur mit CR-Kontext lesbar

## Verdikt

**Merge-blockierend:** ja — für die Slice-Closure (der Commit selbst ist die Umsetzung; die zwei MEDIUM sind Rank-1-Vertrags-Spannen: [`LH-OPEN-001`](../../spec/lastenheft.md#lh-open-001--implementierungssprache-entschieden) Mindestversion und die ungetragene Definition des „denselben Pin-Werts" je [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette); beide sind nur über eine freigegebene Berichtigung des Lastenheft-Textes zu lösen und gehören vor §7/Closure berichtigt). Kein Rollen-Widerspruch: der Implementer hat beide Spannen selbst offen gelegt, die Bewertung bestätigt sie; die übrigen offenen Punkte sind Ausgewisene (Minor-Fangauslegung, 4.3-Floor, tidy-Normierung, CHANGELOG-IDs, Digest, `--print-mk`-Konform, Hilfs-Image-Ableitung). Die LOW-Findings sind nicht-blockierende Nachträge an den gleichen Trägern.

**Übergabe:** Findings gehen an den Implementer (F-1/F-2 berühren Wortlaut, der eine Freigabe braucht — Rückkante Review → Plan/CR, nicht stille Berichtigung); die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler.Dieser Report ist ein **Lauf-Beleg** (Audit: dieser Commit, dieser Skill, dieses Modell, dieses Verdikt); DoD-/Spec-Konformität und Verification-Evidence prüft der Verifier separat (Modul 11).