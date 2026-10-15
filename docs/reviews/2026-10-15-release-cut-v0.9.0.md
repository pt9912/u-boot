# Review-Report: [`slice-v1-release-cut-v0.9.0`](../plan/planning/in-progress/slice-v1-release-cut-v0.9.0.md) — 2026-10-15

**Review-Art:** Code/Doku — Release-Cut-Diff gegen den Ablauf in `docs/maintainer/releasing.md` §3 Nr. 1–5 und §5, den Slice-Plan und die Slice-Vorlage der Baseline.

**Gegenstand:** Commit `4a39e7e` (Release-Cut v0.9.0, zehn Dateien) und [`slice-v1-release-cut-v0.9.0`](../plan/planning/in-progress/slice-v1-release-cut-v0.9.0.md) (liegt bis zur Closure in `in-progress/`, danach `done/`).

**Skill:** `.harness/skills/reviewer.md` @ 4a39e7e · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-15

**Eingangs-Kontext:**

- [`slice-v1-release-cut-v0.9.0`](../plan/planning/in-progress/slice-v1-release-cut-v0.9.0.md); Vorläuferreport `2026-10-14-release-cut-v0.8.1.md` (dessen Altlasten — unsortierte Releases-Tabelle, Roadmap-Meilensteine bis v0.5.0, Roadmap-Kopfdatum — werden nicht erneut gemeldet)
- [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [ADR-0020](../plan/adr/0020-drittanbieter-features-im-katalog.md) (Status `Accepted`)
- [`slice-v1-devcontainer-features-homebrew-sdkman`](../plan/planning/done/slice-v1-devcontainer-features-homebrew-sdkman.md), [`slice-v1-devcontainer-feature-python`](../plan/planning/done/slice-v1-devcontainer-feature-python.md), [`slice-v1-sandbox-docker-nachholmessung-desktop`](../plan/planning/open/slice-v1-sandbox-docker-nachholmessung-desktop.md), [`slice-v2-homebrew-formula`](../plan/planning/in-progress/slice-v2-homebrew-formula.md)
- `docs/maintainer/releasing.md` §3, §5; Slice-Vorlage `.harness/baseline/v6.13.0/templates/docs/plan/planning/slice.template.md`
- `AGENTS.md` (Hard Rules), `harness/review.md`

Kein Docker-Lauf außer `make docs-check`, keine Änderung an bestehenden Dateien. Der Tag `v0.9.0` existiert lokal nicht (`git tag` endet bei `v0.8.1`).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Der CHANGELOG-Eintrag unter 0.9.0 schreibt „The egress allowlist gets the hosts of both features", nennt davor aber drei Katalog-Schlüssel (`python`, `sdkman`, `homebrew`). Der Code trägt Hosts für alle drei (`devcontainer_egress.go:77-79`); „both" ist dadurch mehrdeutig (zwei oder drei Features). | [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), Maintainability | `CHANGELOG.md:25` | ja — Gegenlesen | Zahlwort im Release-Eintrag passt nicht zur Aufzählung |
| F-2 | INFO | „released 2026-10-15" in `README.md` und `README.de.md` sowie die Ziele `releases/tag/v0.9.0`, `download/v0.9.0` und der Compare-Link `v0.8.1...v0.9.0` lösen erst nach dem Tag-Push auf; das Datum steht an sechs Stellen fest und stimmt nur bei Tag-Push am selben Tag. | `docs/maintainer/releasing.md` §3/§5/§7 | `README.md:10`, `README.md:206`, `README.de.md:11`, `README.de.md:210`, `CHANGELOG.md:14`, `docs/user/benutzerhandbuch.md:5` | ja — nach dem Push: `releasing.md` §6 | — |
| F-3 | INFO | Die Slice-DoD (§2) und die Closure-Notiz (§7) sind offen — vor dem Tag der Sollzustand. Mit diesem Report ist der Punkt „Review durchgeführt" erfüllbar; `make fullbuild`, `make test-docker`, Tag, `publish`-Lauf, Closure-Notiz und der `git mv` nach `done/` sind Maintainer-/Verifier-Aufgabe. Die Verschiebung der Welle nach „Abgeschlossene Wellen" folgt nach dem Tag (`releasing.md` §3 Nr. 4). | `harness/verification.md`, `docs/maintainer/releasing.md` §4, §7 | Cut-Slice in `in-progress/`, Zeilen 33-42, 71-79 | nein — Verifier | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| CHANGELOG (§3 Nr. 1): `[Unreleased]` als leerer Anker (`CHANGELOG.md:12`), `[0.9.0] - 2026-10-15` mit Lead-Absatz (`:14-18`), Compare-Links `[Unreleased]: …v0.9.0...HEAD` und `[0.9.0]: …v0.8.1...v0.9.0` am Fuß, Folgezeile `[0.8.1]` unverändert | geprüft, ohne Befund (Zahlwort siehe F-1) |
| CHANGELOG-Inhalt: `### Added` deckt sich mit der Umsetzung — `python` mit `pypi.org`/`files.pythonhosted.org`, `sdkman` als Alias des offiziellen `java`-Features (Einmal-Rendering `devcontainer_features.go:311`), `homebrew` als Community-Feature `ghcr.io/meaningful-ooo/devcontainer-features/homebrew:2` (`:223-226`), Drittanbieter-Warnung bei `generate devcontainer` (`:325-352`), Hosts je Feature in `devcontainer_egress.go:77-79`; Anker [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) und [ADR-0020](../plan/adr/0020-drittanbieter-features-im-katalog.md) lösen auf | geprüft, ohne Befund |
| Versionsstrings (§3 Nr. 2): `cmd/uboot/main.go:43`, `Makefile:27`, `Dockerfile:39` je `0.9.0-dev`; Kommentar-Erwähnungen (Dockerfile `:35`, `:37`; Makefile `:23`, `:26`) mitgezogen; `grep` nach `0.8.1-dev` findet nur die Beschreibung des Bumps im Cut-Slice (Zeile 49) sowie frühere Slices und Reports | geprüft, ohne Befund |
| READMEs (§3 Nr. 3): Status-Block, Zeile `v0.9.0` in der Releases-Tabelle und `VERSION=0.9.0` in EN und DE konsistent zueinander, zum CHANGELOG und zu [ADR-0020](../plan/adr/0020-drittanbieter-features-im-katalog.md); Zeile steht an der Bestandsposition (Fehlordnung der Tabelle als Altlast bekannt) | geprüft, ohne Befund |
| Benutzerhandbuch: `Software-Version: v0.9.0`, `Stand: 2026-10-15`, Beispielausgabe `u-boot version 0.9.0`, `download/v0.9.0`, `VERSION=0.9.0`, Gültigkeitsbereich konsistent | geprüft, ohne Befund |
| Claims zu Docker Desktop: README, Releases-Zeilen, CHANGELOG und Cut-Slice (§1) sagen Docker Desktop nicht zu; das Handbuch nennt es nur als Beispiel für die VM (`:1005`) und unter „Nicht geprüft" (`:1100`); Folge-Slice [`slice-v1-sandbox-docker-nachholmessung-desktop`](../plan/planning/open/slice-v1-sandbox-docker-nachholmessung-desktop.md) existiert in `open/` | geprüft, ohne Befund |
| Roadmap (§3 Nr. 4): „Aktuelle Welle" führt v0.9.0 mit „Tag-Push steht aus", Link auf den Cut-Slice löst auf, offener macOS-Smoke von [`slice-v2-homebrew-formula`](../plan/planning/in-progress/slice-v2-homebrew-formula.md) bleibt benannt; „Abgeschlossene Wellen" trägt die Zeilen für [`slice-v1-devcontainer-features-homebrew-sdkman`](../plan/planning/done/slice-v1-devcontainer-features-homebrew-sdkman.md) (`roadmap.md:85`) und [`slice-v1-devcontainer-feature-python`](../plan/planning/done/slice-v1-devcontainer-feature-python.md) (`:86`) | geprüft, ohne Befund |
| Slice gegen Vorlage (§3 Nr. 5): Kopf (Lifecycle, Welle, Bezug mit aktiven ADRs, Berührte Spec-Stellen „—" mit Begründung, Verantwortlich `pt9912`, Autor/Datum), §1–§8 in Vorlagenreihenfolge; §1 Ausschlüsse mit Begründung und existierenden Folge-Slices; §3 Plan-Tabelle deckt alle zehn Dateien des Commits; §4 Start-Trigger erfüllt (beide Feature-Slices in `done/`), Rückführungen vorab benannt; §6 beide Risiken tragen einen Ausgang, das Zustimmungsmodell-Risiko ist als „weiter offen: Entscheidung des Projektinhabers" mit Verweis auf den Review-Befund geführt | geprüft, ohne Befund |
| Slice-Lifecycle: Slice liegt nur in `in-progress/`, Eintrag in `in-progress/README.md` vorhanden; Verweise auf `done/`- und `open/`-Slices lösen auf; ADR/Lastenheft verweisen nicht abwärts; Cut-Slice in `in-progress/` bis zur Closure ist kein Befund | geprüft, ohne Befund |
| Hard Rules: kein `//nolint`, kein `git mv`, keine Gate-/Schwellen- oder ADR-Änderung, `spec/architecture.md` unberührt; Kommentare der Versionsstrings beschreiben den Ist-Zustand (3.7); Commit-Message nennt [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) und [ADR-0020](../plan/adr/0020-drittanbieter-features-im-katalog.md) | geprüft, ohne Befund |
| `make docs-check` | ausgeführt: 207 Dateien geprüft, 0 Befunde (vor Anlage dieses Reports) |

Nicht geprüft:

- `make gates`, `make ci`, `make test-docker`, `make fullbuild` — Auftrag verbietet Docker-Läufe außer `make docs-check`; Verifier führt sie vor dem Tag aus (`docs/maintainer/releasing.md` §4).
- `publish`-Lauf, Tag-Hash, GHCR-/Release-Assets, Tap-Job — der Tag ist noch nicht gesetzt (`releasing.md` §6).
- Inhaltliche Richtigkeit der Feature-Umsetzung (Tests, Goldens) — abgedeckt durch die Reports `2026-10-15-devcontainer-features-homebrew-sdkman` und `2026-10-15-devcontainer-feature-python`, nicht Gegenstand dieses Laufs.
- Entscheidung zum Zustimmungsmodell für Drittanbieter-Features — Zuständigkeit des Projektinhabers, im Slice (§6) als offen geführt.
- `.github/workflows/` und `nfpm`-Konfiguration — im Diff nicht berührt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Zahlwort im Release-Eintrag passt nicht zur Aufzählung

## Verdikt

**Merge-blockierend:** nein — kein HIGH/MEDIUM; das LOW-Finding ist Doku-Präzision und blockiert den Tag nicht. Der Cut ist mit `releasing.md` §3 Nr. 1–5 konsistent; offen bleiben die Maintainer-/Verifier-Schritte (§4 Gates, §5 Tag, §6 Kontrolle, Closure und `git mv` nach `done/`).

**Übergabe:** Findings gehen an den Implementer; F-1 lässt sich vor dem Tag im CHANGELOG korrigieren, F-2 ist nach dem Tag-Push gegenzulesen. Der Report ersetzt keine Verification.
