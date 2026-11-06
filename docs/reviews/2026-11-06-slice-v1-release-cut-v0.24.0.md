# Review-Report: slice-v1-release-cut-v0.24.0 — 2026-11-06

**Review-Art:** Release-Cut-Review (Text und Versionsstellen) gegen Plan, Spec, Hard Rules und `releasing.md` §3/§4 (Modul 10). Keine DoD-Prüfung — das ist die Frage des Verifiers.

**Gegenstand:** Cut-Commit `58411aaf` (Basis `e38d6461`; 8 Dateien, 28 Einfügungen, 19 Löschungen) und, soweit der Lead Aussagen über den Release trägt, der Stand von `main` gegen den Tag `v0.23.0` (`git diff v0.23.0..HEAD --stat -- cmd internal ':!*_test.go'`: 27 Dateien, wie im Plan gelesen). Slice [`slice-v1-release-cut-v0.24.0`](../plan/planning/done/slice-v1-release-cut-v0.24.0.md); der Slice ist bewusst nicht geschlossen (§7 leer, DoD-Haken offen, kein Tag) — kein Finding.

**Skill:** `.harness/skills/reviewer.md`
**Modell:** claude (Agentenlauf) · **Datum:** 2026-11-06

**Eingangs-Kontext:**

- Slice-Plan (§1 Ziel, Abgrenzung, M-Start (a)/(b); §2 DoD; §3 Plan-Tabelle; §6 Risiken 1, 2, 7, 10).
- [`docs/maintainer/releasing.md`](../maintainer/releasing.md) §3 und §4.
- [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-FA-DEV-013`](../../spec/lastenheft.md#lh-fa-dev-013--homebrew-formeln-im-devcontainer), [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette), [`ADR-0043`](../plan/adr/0043-nur-devcontainer-bau-pfad.md), [`ADR-0046`](../plan/adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md), [`ADR-0047`](../plan/adr/0047-brew-formeln-lastenheft-stellen-und-verteilung-lastenheft-spezifikation.md); Lastenheft 0.37.3 (`spec/lastenheft.md` Version im Kopf).
- `AGENTS.md` §3.2, §3.7, §3.9, §3.10, §5 Regel 6; [`harness/review.md`](../../harness/review.md).
- Quellen der *Changed*-Aussage gelesen: `tools/mk/tool-pruefung.mk`, `tools/harness/gate-reihenfolge.sh`, `tools/harness/gate-reihenfolge-test.sh`.
- Vorgänger-Report derselben Klasse: [`2026-11-04-slice-v1-release-cut-v0.23.0`](2026-11-04-slice-v1-release-cut-v0.23.0.md).

**Sensors:**

- `make docs-check`: nach dem Schreiben dieses Reports gefahren (Ergebnis im Handback).
- `make gates`, `make ci`, `make build VERSION=0.24.0`, Tap-Trockenlauf, `make test-docker`, `make verify-depguard`: nicht gefahren — Review-Rolle liest, wiederholt keine Gate-Läufe (Vorab-Prüfungen sind Teil von DoD-Punkt 2, also Verifier-Aufgabe); keine Importregel berührt (Diff im Produktcode: allein die Versionszeile in `cmd/uboot/main.go`).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der Lead nennt den Exit-Code `10` für `devcontainer remove --volumes` ohne `--yes` im nicht-interaktiven Modus und unter `--json` ohne Bedingung. Das Lastenheft 0.37.3 und der *Added*-Eintrag sagen: der Fall gilt nur, wenn es einen Container zu entfernen gibt; ohne Container endet derselbe Aufruf mit Exit-Code `0` und ohne Abfrage. | [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Plan §6 Risiko 1 und 7, `AGENTS.md` §3.10 | `CHANGELOG.md` · „(exit code `10` without `--yes` in non-interactive mode and under `--json`)“ | nein — Review-only (Abgleich Lead gegen *Added*-Eintrag und Lastenheft-Zeile „gilt das nur, wenn es einen Container zu entfernen gibt“) | Lead-Aussage ohne Bedingung des Eintrags |
| F-2 | LOW | Der Lead sagt „`tool-pruefung` stands before all other checks of `make gates`“ ohne die Einschränkung „serial“; der *Changed*-Eintrag führt sie („so a serial `make gates` …“, „Not promised: the order of the checks under `make -j`“), der Kommentar in `tools/mk/tool-pruefung.mk` ebenso. | `AGENTS.md` §3.10, Plan §1 (b) („Zusage und Grenze unter `make -j`“) | `CHANGELOG.md` · „`tool-pruefung` stands before all other checks of `make gates`“ | nein — Review-only | Lead-Aussage ohne Bedingung des Eintrags |

### Zu F-1 (Failure-Szenario)

Ein Leser, der nur den Lead liest (der erste Absatz des Blocks, die Zusammenfassung des Release), baut ein Skript oder eine CI-Stufe auf „`remove --volumes` ohne `--yes` endet nicht-interaktiv mit `10`“ und behandelt Exit `0` als Erfolg einer Löschung. Gibt es keinen Container, endet der Aufruf mit `0`, obwohl nichts gelöscht wurde und keine Freigabe geprüft wurde; das Handbuch (`docs/user/benutzerhandbuch.md`, Abschnitt „Kein Container“ und „Exit-Codes von `devcontainer remove`“) und der *Added*-Eintrag tragen den Fall, der Lead nicht. Das ist die Klasse, die Plan §6 Risiko 1 für den Lead benennt (Cut `v0.23.0`: einem Satz fehlte die Bedingung); die Messstelle „Lead gegen Einträge lesen“ greift hier. Die Spannung zwischen Exit-Code `2` und `10` bleibt wie im Plan (Risiko 7) beim Projektinhaber und ist nicht Gegenstand dieses Findings; es geht allein um die fehlende Bedingung „es gibt einen Container“.

### Zu F-2 (Einordnung)

Der Lead ordnet die Änderung ausdrücklich als „how the repository itself is built, not what the product does“ ein; das begrenzt den Schaden. Die Klasse ist dieselbe wie bei F-1, deshalb im selben Lauf gezählt, aber ohne Vertragsfolge für den Nutzer des Produkts: LOW mit Anker (§3.10, Plan §1 (b)).

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `CHANGELOG.md` *Changed*-Eintrag zur Gate-Reihenfolge gegen `tools/mk/tool-pruefung.mk` (`GATE_CHECKS := tool-pruefung $(GATE_CHECKS)`), `tools/harness/gate-reihenfolge.sh` (liest `record-gates` aus `make -pn gates`), `tools/harness/gate-reihenfolge-test.sh` (grün auf Kopie, rot bei `+=`) und `tool-pruefung`-Rezept (beide Skripte laufen darin) | geprüft, ohne Befund; *Promised*/*Not promised* decken sich mit dem gelesenen Code; kein Wort zu `BUILD_MODE` oder zum eigenen Devcontainer (Plan §1 (b), Zeilen 2 und 3) |
| `CHANGELOG.md` Lead: `devcontainer remove` (Stopp zuerst, Volumes bleiben, `--volumes` nach `[y/N]` auf Standard-Error), `brew` (gepinnter Commit ohne Installer-Skript, Formeln/Taps nicht gepinnt und ausgewiesen, Freigabe `devcontainer.brewTaps.allow`, Exit `10` vor dem Schreiben), „`arm64` is not measured for `brew`“, „not promised inside the released container image“, „Lastenheft 0.36.0 to 0.37.3“ | geprüft, ohne Befund (außer F-1 an einer Teilaussage); Aussagen decken sich mit den *Added*-Einträgen und dem Lastenheft-Kopf (0.37.3); nennt alle drei Einträge |
| Risiko 2 / `rg` im Container: `ripgrep`/`rg` in den hinzugefügten Zeilen des Cut-Commits | geprüft, ohne Befund; keine hinzugefügte Zeile nennt `rg` oder `ripgrep`; der Wortlaut im unverändert gebliebenen *Added*-Eintrag ist die Testformel der `brew`-Zeile |
| Versionsstellen: `cmd/uboot/main.go`, `repo.mk` (Kommentar und `VERSION ?=`), `Dockerfile` (Kommentare und `ARG UBOOT_VERSION=`) | geprüft, ohne Befund; alle drei Stellen samt Kommentaren auf `0.24.0-dev`; kein Rest `0.23.0` in `README*`, `docs/user`, `docs/maintainer`, `packaging`, `scripts`, `.github`, `tools`, `Makefile`, `repo.mk`, `Dockerfile`, `cmd`, `internal` (Rest nur in Messprotokollen unter `docs/messungen/` und im CHANGELOG, dort zutreffend als frühere Fassung) |
| `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md` (Kopf, `u-boot version`, `download/v…`, `VERSION=`, Gültigkeitsbereich; `Stand: 2026-11-06` gleich dem Datum des Laufs) | geprüft, ohne Befund; keine Anforderungs-, ADR- oder Spec-Kennung hinzugekommen (`AGENTS.md` §5 Regel 6) |
| `CHANGELOG.md` Struktur: `## [Unreleased]` leer, `## [0.24.0] - 2026-11-06`, Compare-Links `[Unreleased]: …v0.24.0...HEAD` und `[0.24.0]: …v0.23.0...v0.24.0`, Reihenfolge *Added* vor *Changed* | geprüft, ohne Befund |
| `docs/maintainer/releasing.md` §4 (Kommentarzeile zu `make gates` nennt `tool-pruefung`) | geprüft, ohne Befund; deckt sich mit der Gate-Zeile in `harness/README.md` |
| Scope-Disziplin: Cut-Commit berührt weder `spec/`, ADRs, `.github/`, `scripts/`, `packaging/`, `Makefile`, `d-check.mk`, `.golangci.yml`, `go.mod` noch Produktcode außer der Versionszeile | geprüft, ohne Befund (Gate-Lockerung nach `AGENTS.md` §3.6 und Suppression nach §3.2 nicht berührt) |
| Kommentar-Selbstcheck `AGENTS.md` §3.9 über die geänderten Kommentare (`repo.mk`, `Dockerfile`, `cmd/uboot/main.go`) | geprüft, ohne Befund; Kommentare nennen die geltende Zusage, keine Chronik; alter Wert `0.23.0-dev` im Baum nicht mehr referenziert |
| Traceability der Commit-Message | geprüft, ohne Befund; nennt den Slice und [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-DEV-013`](../../spec/lastenheft.md#lh-fa-dev-013--homebrew-formeln-im-devcontainer), [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) |
| Lifecycle / Roadmap | geprüft, ohne Befund; der Cut-Commit berührt `roadmap.md` nicht (die Zeile unter *Offene Wellen* kam mit `e38d6461`; die Zeile unter *Abgeschlossene Wellen* gehört nach dem Tag, Plan §2 Punkt 3) |
| Carveouts, Dateisicherheit, Exit-Code-Klassifikation, Sprachvertrag, Architekturgrenzen, Replay | nicht berührt (kein Code- oder Generator-Diff) |

Nicht geprüft:

- `make gates`, `make ci`, `make build VERSION=0.24.0` samt Label-Gegenprobe, Tap-Trockenlauf nach `releasing.md` §4 — Verifier-Aufgabe (DoD-Punkt 2); ihre Belege stehen nach Plan §6 Risiko 8 im Stand-Bericht und im Verification-Report.
- `ci`, `integration-docker` und `tap-check` auf dem Push-Tip — der Push steht aus (Plan §1, erste Freigabe); der Review läuft nach Plan §3 Schritt 4 vor dem Push.
- Dass die Nummer `0.24.0` als Minor gewählt ist — vom Projektinhaber am 2026-11-06 bestätigt (Plan, Kopf); kein Gegenstand des Reviews.
- Wahrheit der Messaussagen im unveränderten *Added*-Eintrag zu `brew` (Messwerte `+292 MB`, `35 s`, Engine 29.9.3) — im Cut-Commit nicht geändert, getragen von den Messprotokollen der trägenden Slices.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Lead-Aussage ohne Bedingung des Eintrags (2 Treffer: F-1, F-2)

## Verdikt

**Merge-blockierend:** ja — für den Push (nicht für die Arbeit am Slice): F-1 ist ein MEDIUM mit Failure-Szenario und wird vor dem Push geklärt, damit der Lead mit dem Tag nicht eine Zusage veröffentlicht, die enger oder weiter ist als Eintrag und Lastenheft. F-2 blockiert nicht.

**Übergabe:** Findings gehen an den Implementer; der Lead-Text ist Teil des Cut-Commits, jede Änderung daran braucht einen weiteren Cut-Commit vor dem Push und eine erneute Lesung gegen den Diff (M-L). Die Finding-Klasse geht in die Slice-Closure §7 und von dort in den Zähler; die Beobachtung [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md) führt die verwandte Klasse. Dieser Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier.
