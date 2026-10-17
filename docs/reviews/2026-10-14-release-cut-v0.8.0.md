# Review-Report: [`slice-v1-release-cut-v0.8.0`](../plan/planning/done/slice-v1-release-cut-v0.8.0.md) — 2026-10-14

**Review-Art:** Code/Doku — Release-Cut-Diff gegen den Ablauf in `docs/maintainer/releasing.md` §3–§5, den Slice-Plan und die Slice-Vorlage der Baseline.

**Gegenstand:** Commit `fd40f78` (Release-Cut v0.8.0, neun Dateien) und [`slice-v1-release-cut-v0.8.0`](../plan/planning/done/slice-v1-release-cut-v0.8.0.md).

**Skill:** `.harness/skills/reviewer.md` @ fd40f78 · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-14

**Eingangs-Kontext:**

- [`slice-v1-release-cut-v0.8.0`](../plan/planning/done/slice-v1-release-cut-v0.8.0.md), Muster [`slice-v1-release-cut-v0.7.0`](../plan/planning/done/slice-v1-release-cut-v0.7.0.md)
- [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [ADR-0018](../plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md) (Status `Accepted`), Lastenheft 0.5.0
- [`slice-v1-sandbox-docker-nachholmessung-desktop`](../plan/planning/done/slice-v1-sandbox-docker-nachholmessung-desktop.md), [`slice-v2-homebrew-formula`](../plan/planning/done/slice-v2-homebrew-formula.md)
- `docs/maintainer/releasing.md` §3 Vorbereitung, §4 Vorab-Prüfung, §5 Tag setzen
- `AGENTS.md` (Hard Rules), `harness/review.md`

Kein Docker-Lauf, keine Änderung an bestehenden Dateien. Der Tag `v0.8.0` existiert lokal nicht (`git tag` endet bei `v0.7.0`).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Die Releases-Tabelle ist nicht nach Version sortiert: Die Zeile `v0.8.0` folgt auf `v0.7.0` und steht damit vor `v0.5.0` und `v0.4.0` (Reihenfolge 0.1, 0.2, 0.3, 0.6, 0.7, 0.8, 0.5, 0.4). Die Fehlordnung stammt aus dem Bestand; der Cut schreibt sie fort, in `README.md` und `README.de.md`. | Maintainability, `docs/maintainer/releasing.md` §3 Nr. 3 | `README.md:205`, `README.de.md:209` | ja — Gegenlesen | Releases-Tabelle ungeordnet |
| F-2 | LOW | Die Roadmap führt weder in der Meilenstein-Tabelle (endet bei `v0.5.0`, Zeilen 33–37) noch im Abhängigkeitsgraphen (Zeile 50) Einträge für `v0.6.0`, `v0.7.0`, `v0.8.0`; unter „Abgeschlossene Wellen" fehlt eine Zeile `v0.7.0` (vorhanden: `v0.5.0` Zeile 71, `v0.6.0` Zeile 80). `releasing.md` §3 Nr. 4 verlangt nur die Aktive-Welle-Führung bis zum Tag; die Lücken liegen im Bestand. | `docs/maintainer/releasing.md` §3 Nr. 4, Maintainability | `docs/plan/planning/in-progress/roadmap.md:33-37`, `:50`, `:80` | ja — Gegenlesen | Roadmap-Meilensteine hinter Release-Stand |
| F-3 | INFO | Ein v0.8.0-Eintrag unter „Abgeschlossene Wellen" fehlt bewusst: Der Slice-Plan (§3) und `releasing.md` §3 Nr. 4 verlangen die Verschiebung erst nach dem Tag; „Aktuelle Welle" nennt `v0.8.0` mit „Tag-Push steht aus" (`roadmap.md:16`). Die neue Zeile in Abgeschlossene Wellen (`roadmap.md:81`) gehört der Spec-/Umsetzungs-Arbeit, nicht dem Release, und ist konsistent (Hashes `6d726e0` und `e1c2c3a` existieren, die vier Slices liegen unter `done/`). | `docs/maintainer/releasing.md` §3 Nr. 4 | `docs/plan/planning/in-progress/roadmap.md:16`, `:81` | ja — Gegenlesen | — |
| F-4 | INFO | `README.md` und `README.de.md` sagen „released 2026-10-14", die Releases-Tabelle und das Handbuch verlinken `releases/tag/v0.8.0` und `download/v0.8.0`, die CHANGELOG-Compare-Links zeigen auf `v0.7.0...v0.8.0`; der Tag ist noch nicht gepusht, die Ziele lösen bis zum `publish`-Lauf nicht auf. Das entspricht der Praxis früherer Cuts (siehe [`slice-v1-release-cut-v0.7.0`](../plan/planning/done/slice-v1-release-cut-v0.7.0.md)), die Aussage ist vor dem Tag-Push vorläufig. | `docs/maintainer/releasing.md` §3/§5 | `README.md:10`, `README.md:205`, `README.de.md:11`, `docs/user/benutzerhandbuch.md:106` | ja — nach dem Push: `releasing.md` §6 | — |
| F-5 | INFO | Das Datum `2026-10-14` steht an fünf Stellen als Release-Datum fest (CHANGELOG `[0.8.0]`, beide READMEs Status und Tabelle, Handbuch „Stand"); es stimmt nur, wenn der Tag am selben Tag gesetzt wird. Der Cut-Commit trägt das Datum 2026-10-14. Wird am Folgetag getaggt, driften Tag-Datum und alle fünf Stellen. | `docs/maintainer/releasing.md` §3 Nr. 1 | `CHANGELOG.md:14`, `README.md:10`, `README.de.md:11`, `docs/user/benutzerhandbuch.md:5` | nein — Review-only (Tag-Datum nach Push gegenlesen) | — |
| F-6 | INFO | Die Slice-DoD (§2) ist vollständig offen (alle Häkchen leer, §7 Closure-Notiz „—"); das ist vor dem Tag der Sollzustand. Offene DoD-Punkte, die diesen Review betreffen: „Review durchgeführt, Report unter `docs/reviews/`" wird mit diesem Report erfüllbar; Tag, `publish`-Lauf, `make fullbuild` und Closure-Notiz sind Verifier-/Maintainer-Aufgabe. Der Commit nennt keine Evidence für `make ci`/`make test-docker`/`make fullbuild` (der Cut-Slice v0.7.0 trug sie). | `harness/verification.md`, `docs/maintainer/releasing.md` §4 | [`slice-v1-release-cut-v0.8.0`](../plan/planning/done/slice-v1-release-cut-v0.8.0.md) (§2, Zeilen 33–42) | nein — Verifier (`make fullbuild`, `gh run list`) | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| CHANGELOG: `[Unreleased]` als leerer Anker (`CHANGELOG.md:12`), `[0.8.0] - 2026-10-14` mit Lead-Absatz (`:14-19`), Compare-Links `[Unreleased]: …v0.8.0...HEAD` und `[0.8.0]: …v0.7.0...v0.8.0` am Fuß | geprüft, ohne Befund |
| CHANGELOG-Inhalt: Lead „Eighth release" (0.1 bis 0.8 stimmt), „Lastenheft 0.5.0" stimmt mit `spec/lastenheft.md` Version 0.5.0; [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)`-Anker und [ADR-0018](../plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md) (Status `Accepted`) lösen auf; Abschnitte `### Added` und `### Security` unverändert in den [0.8.0]-Block gewandert | geprüft, ohne Befund |
| Versionsstrings: `cmd/uboot/main.go:43`, `Makefile` (`VERSION ?=`), `Dockerfile` (`ARG UBOOT_VERSION=`) jeweils `0.8.0-dev`; alle vier Kommentar-Erwähnungen (Dockerfile zwei, Makefile zwei) mitgezogen; `grep` nach `0.7.0-dev` findet außerhalb von Slices (historisch) und Slice-Plan (Beschreibung des Bumps) keinen Rest | geprüft, ohne Befund |
| Versionsstrings an genau drei Stellen: keine weitere Quelle (Tests, Skripte, Workflows) trägt `0.7.0` als aktuellen Wert | geprüft, ohne Befund |
| READMEs: Status-Block, Releases-Zeile `v0.8.0` und `VERSION=0.8.0` in beiden Sprachen konsistent zueinander, zum CHANGELOG und zu [LH-FA-DEV-007](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) und [ADR-0018](../plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md); Docker Desktop wird in README, Releases-Zeile und CHANGELOG nicht zugesagt | geprüft, ohne Befund |
| Benutzerhandbuch: `Software-Version: v0.8.0`, `Stand: 2026-10-14`, Beispielausgabe, `download/v0.8.0`, `VERSION=0.8.0`, Gültigkeitsbereich konsistent; die Docker-Desktop-Erwähnungen (`:1000` Isolationswarnung, `:1086` „Nicht geprüft: Docker Desktop") sagen keine Unterstützung zu | geprüft, ohne Befund |
| Slice gegen Vorlage: Kopf (Lifecycle, Welle, Bezug, Berührte Spec-Stellen „—" mit Begründung, Verantwortlich, Autor/Datum), §1–§8 in Vorlagenreihenfolge vorhanden; §1 Ausschlüsse mit Begründung und existierenden Folge-Slices; §4 Rückführungen vorab benannt; §6 Risiken tragen je einen Ausgang | geprüft, ohne Befund |
| Slice-Lifecycle und Referenzmodell: Slice liegt nur in `in-progress/` (nicht doppelt), Eintrag im `in-progress/README.md` vorhanden, die vier Sandbox-Docker-Zeilen entfernt (Slices bereits unter `done/`); Ziele der Slice-Links ([slice-v1-sandbox-docker-nachholmessung-desktop](../plan/planning/done/slice-v1-sandbox-docker-nachholmessung-desktop.md), [slice-v2-homebrew-formula](../plan/planning/done/slice-v2-homebrew-formula.md)) existieren; ADR/Lastenheft verweisen nicht abwärts | geprüft, ohne Befund |
| Hard Rules: kein `//nolint`, kein `git mv` im Commit, keine Gate-/Schwellenänderung, keine ADR-Änderung, `spec/architecture.md` unberührt | geprüft, ohne Befund |
| Commit-Traceability: Message nennt [LH-FA-DEV-007](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) und [ADR-0018](../plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md) | geprüft, ohne Befund |

Nicht geprüft:

- `make gates`, `make ci`, `make test-docker`, `make fullbuild`, `make docs-check` — Auftrag verbietet Docker-Läufe; Verifier führt sie vor dem Tag aus (`docs/maintainer/releasing.md` §4).
- `publish`-Lauf, Tag-Hash, GHCR-/Release-Assets, Tap-Job — der Tag ist noch nicht gesetzt (`releasing.md` §6).
- Inhaltliche Richtigkeit der Feature-Aussagen gegen den Code (`nestedRuntime: docker`) — abgedeckt durch die Reports `2026-10-14-sandbox-docker-runtime` und `-r2`, nicht Gegenstand dieses Laufs.
- `.github/workflows/` und `nfpm`-Konfiguration — im Diff nicht berührt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Releases-Tabelle ungeordnet · Roadmap-Meilensteine hinter Release-Stand

## Verdikt

**Merge-blockierend:** nein — kein HIGH/MEDIUM; die beiden LOW-Findings stammen aus dem Bestand und blockieren den Tag nicht. Der Cut ist mit `releasing.md` §3 Nr. 1–4 konsistent; offen bleiben die Maintainer-/Verifier-Schritte (§4 Gates, §5 Tag, §6 Kontrolle, Slice-Closure und `git mv` nach `done/`).

**Übergabe:** Findings gehen an den Implementer; F-4 und F-5 sind nach dem Tag-Push gegenzulesen. Der Report ersetzt keine Verification.
