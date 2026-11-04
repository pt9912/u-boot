# Review-Report: slice-v1-release-cut-v0.22.0 — 2026-11-04

**Review-Art:** Code (Release-Cut gegen Plan, Release-Prozess und Hard Rules)

**Gegenstand:** Cut-Commit `f1711f3a` (Tag `v0.22.0`, veröffentlicht) und `git diff v0.21.0..v0.22.0`

**Skill:** `.harness/skills/reviewer.md` @ 3ec45359 ·
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-11-04

**Eingangs-Kontext:**

- [`slice-v1-release-cut-v0.22.0`](../plan/planning/done/slice-v1-release-cut-v0.22.0.md) (Stand: in-progress, §7 leer und DoD offen sind bewusst, kein Finding)
- [`ADR-0043`](../plan/adr/0043-nur-devcontainer-bau-pfad.md), [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette)
- `docs/maintainer/releasing.md` §3 und §4 (Maßstab; nur die Versions- und Datumsstellen gelesen)
- `AGENTS.md` §3.7, §3.9, §3.10, §5 Regel 6
- gelesene Diffs: `git show f1711f3a`, `git diff v0.21.0..v0.22.0 -- cmd internal ':!*_test.go' go.mod`, `git diff v0.21.0..v0.22.0 --stat`

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der Lead sagt „no non-test file under `cmd/` or `internal/` differs from 0.21.0“. Der Vergleich `v0.21.0..v0.22.0` zeigt genau eine solche Datei: `cmd/uboot/main.go` (`var version = "0.21.0-dev"` → `"0.22.0-dev"`), geändert im Cut-Commit selbst. Die Messung des Slice („leer“, `v0.21.0..HEAD` am 2026-11-04) lag vor dem Cut-Commit. | `AGENTS.md` §3.10, [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) (Release-Aussage) | `CHANGELOG.md` · „no non-test file under `cmd/` or `internal/` differs from 0.21.0“ | ja — `git diff v0.21.0..v0.22.0 --stat -- cmd internal ':!*_test.go'` listet `cmd/uboot/main.go` | Aussage über Code gilt für den Stand vor dem Cut-Commit, nicht für den getaggten |
| F-2 | LOW | Lead und zweiter Eintrag nennen von `go.mod` nur die Anhebung der `go`-Direktive. Der Diff `v0.21.0..v0.22.0` ordnet zusätzlich `github.com/spf13/cobra` und `github.com/spf13/pflag` von `// indirect` zu direkten Anforderungen um und mischt die `require`-Blöcke neu; Versionen bleiben gleich. | `AGENTS.md` §3.10 | `go.mod` · `github.com/spf13/cobra v1.10.2` | ja — `git diff v0.21.0..v0.22.0 -- go.mod` | Beleg-Menge kleiner als der Diff |
| F-3 | LOW | Plan-Tabelle §3 führt `docs/plan/planning/open/README.md` als „update“ und `go.mod` als „kein Diff“. Der Cut-Commit berührt keine der beiden Dateien (Bestand-Zeile stand schon vorher dort); `go.mod` unterscheidet sich gegenüber `v0.21.0` (F-2), der Slice nennt das selbst in §1. Die Tabelle nennt keinen Bezugspunkt für „kein Diff“. Die Moves des Slice selbst (über das Make-Target des Lifecycle-Wechsels) stehen als Zeile, sind aber nicht in `git show --stat` des Cut-Commits. | `AGENTS.md` §3.10 (Tabelle gegen `git diff --stat`) | [`slice-v1-release-cut-v0.22.0`](../plan/planning/done/slice-v1-release-cut-v0.22.0.md) §3 · „`spec/`, ADRs, `.github/`, `scripts/`, `packaging/`, `Makefile`, `go.mod` | kein Diff“ | ja — Abgleich Tabelle gegen `git show f1711f3a --stat` | Plan-Tabelle ohne Diff-Bezugspunkt |
| F-4 | LOW | Die Roadmap-Zeile unter *Offene Wellen* trägt „Push und Tag stehen aus (Freigabe des Projektinhabers)“. Der Tag `v0.22.0` steht auf `origin` (`refs/tags/v0.22.0` → `f1711f3a`). Der Marker mit Ablaufpunkt ist abgelaufen; die Closure löst ihn laut §6 Risiko 5 ein. | `AGENTS.md` §3.10 (Zustandsmarker), §3.7 (Zustandsfeld) | `docs/plan/planning/in-progress/roadmap.md` · „Push und Tag stehen aus“ | ja — `git ls-remote --tags origin v0.22.0` | Zustandsmarker über Ablaufpunkt hinaus |
| F-5 | INFO | Die Versionsnummer-Wahl (Minor) liegt beim Projektinhaber (Slice §6 Risiko 6). Der Lead nennt zwei Änderungen an Bau- und Gate-Kette ohne Produktänderung; das trägt auch ein Patch-Release. Kein Urteil. | Slice §1 | `CHANGELOG.md` · „Twenty-second release.“ | Review-only | — |
| F-6 | INFO | `make test-docker` lief lokal nicht (cgroup, benannt im Slice). Der Nachweis für `integration-docker` am Push-Tip und für `publish`/`tap` ist Aufgabe des Verifiers; dieser Lauf hat CI-Läufe nicht gelesen. | `harness/verification.md` | — | Review-only (Zuständigkeit Verifier) | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Versionsstellen: `cmd/uboot/main.go`, `repo.mk` (Wert und Kommentar), `Dockerfile` (`ARG UBOOT_VERSION` und beide Kommentare), `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md` (Kopf `Software-Version`/`Stand`, `u-boot version`, `download/v…`, `VERSION=`, Gültigkeitsbereich) | geprüft, ohne Befund: alle auf `0.22.0`/`v0.22.0` bzw. `0.22.0-dev`; Suche nach `0.21.0`/`0.21` im Baum (ohne `.git`, Baseline, Reviews, `done/`, Messungen): Treffer nur CHANGELOG-Chronik, Lastenheft-/Spec-Änderungshistorie (Lastenheft-Version 0.21.0, anderer Zähler), Slice/Roadmap/Open-README als Bezug auf den Vorgänger und ein Testwert `doctor_test.go` |
| `CHANGELOG.md` Compare-Links | geprüft, ohne Befund: `[Unreleased]` → `v0.22.0...HEAD`, `[0.22.0]` → `v0.21.0...v0.22.0`, `[0.21.0]` unverändert; Block `## [0.22.0] - 2026-11-04` steht unter leerem `[Unreleased]` |
| Datum | geprüft, ohne Befund: CHANGELOG-Datum, Handbuch-`Stand` und Commit-Datum des Tags (2026-11-04) stimmen überein |
| Lead: Aussagen über `go`-Direktive `1.27.0` und Toolchain-Pin `1.27.3` | geprüft, ohne Befund: `go.mod` `go 1.27.0`; `GO_VERSION` `1.27.3` in `Dockerfile`, `repo.mk`, `.devcontainer/Dockerfile` |
| Lead: Anzahl und Inhalt der Einträge | geprüft, ohne Befund: zwei Einträge *Changed* (Bau-Pfad; Pin-Inventur), Lead nennt beide ohne Zusatzaussage zum Produkt |
| Nutzer-Doku (`README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md`) im Diff | geprüft, ohne Befund: keine Anforderungs-, ADR-, Spec- oder Planungsverweise hinzugefügt (`AGENTS.md` §5 Regel 6); Zeilen sind reine Versionswechsel |
| Kommentar-Selbstcheck (`AGENTS.md` §3.9) an den Kommentaren in `repo.mk`, `Dockerfile`, `main.go` | geprüft, ohne Befund: nur Versionswerte getauscht, keine Chronik-Muster, kein Symbol im Diff umbenannt |
| Gate-/Harness-Lockerung, Importregeln, Sprachvertrag, Lifecycle | geprüft, ohne Befund: Cut-Commit berührt keine Gate-Datei, keinen Produktcode außer der Versionszeile, kein ADR |
| Lifecycle-Moves des Slice (`open/` → `next/` → `in-progress/`) | geprüft, ohne Befund: Commit-Folge zeigt reine Move-Commits mit getrennten Verweis-Commits |
| Plan-Tabelle §3 gegen `git show f1711f3a --stat` | geprüft, Befund F-3; die übrigen Zeilen (CHANGELOG, `main.go`, `repo.mk`, `Dockerfile`, drei Doku-Dateien, Roadmap) decken den Diff vollständig |
| `docs/maintainer/releasing.md` §3/§4 vollständig, `make docs-check`, `make gates`, CI-Läufe, Tap-Trockenlauf | nicht geprüft: Report liest Diffs, führt keine Gates aus; Gate-Nachweis ist Aufgabe des Verifiers |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 3 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Aussage über Code gilt für den Stand vor dem Cut-Commit · Beleg-Menge kleiner als der Diff · Plan-Tabelle ohne Diff-Bezugspunkt · Zustandsmarker über Ablaufpunkt hinaus

## Verdikt

**Merge-blockierend:** nein für den Diff als Ganzes (kein HIGH); F-1 ist MEDIUM. Der Tag steht und ist veröffentlicht, ein Release-Text ist danach nicht mehr stillschweigend änderbar; Fortgang ist Sache des Implementers und des Projektinhabers (Korrektur im nächsten Release-Text oder Begründung in der Closure-Notiz). Die Herabstufung bleibt dem Reviewer verwehrt, auch wenn der Tag gesetzt ist.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7. F-1 bestätigt Risiko 1 des Slice („Der Lead sagt mehr als die Beleg-Menge“). Kein Konflikt-Pfad, da kein HIGH und kein Rollen-Widerspruch. DoD-Konformität prüft der Verifier separat.
