# Review-Report: [`slice-v1-release-cut-v0.8.1`](../plan/planning/done/slice-v1-release-cut-v0.8.1.md) — 2026-10-14

**Review-Art:** Code/Doku — Release-Cut-Diff gegen den Ablauf in `docs/maintainer/releasing.md` §3 Nr. 1–5 und §5, den Slice-Plan und die Slice-Vorlage der Baseline.

**Gegenstand:** Commit `290ce07` (Release-Cut v0.8.1, zwölf Dateien) und [`slice-v1-release-cut-v0.8.1`](../plan/planning/done/slice-v1-release-cut-v0.8.1.md) (liegt bis zur Closure in `in-progress/`, danach `done/`).

**Skill:** `.harness/skills/reviewer.md` @ 290ce07 · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-14

**Eingangs-Kontext:**

- [`slice-v1-release-cut-v0.8.1`](../plan/planning/done/slice-v1-release-cut-v0.8.1.md), Muster [`slice-v1-release-cut-v0.8.0`](../plan/planning/done/slice-v1-release-cut-v0.8.0.md); Vorläuferreport `2026-10-14-release-cut-v0.8.0.md` (dessen Altlasten F-1/F-2 werden nicht erneut gemeldet)
- [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [ADR-0019](../plan/adr/0019-docker-daemon-aus-docker-apt-repo.md) (Status `Accepted`), [ADR-0018](../plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md)
- [`slice-v1-sandbox-docker-ce`](../plan/planning/done/slice-v1-sandbox-docker-ce.md), [`slice-v1-sandbox-docker-nachholmessung-desktop`](../plan/planning/done/slice-v1-sandbox-docker-nachholmessung-desktop.md), [`slice-v2-homebrew-formula`](../plan/planning/done/slice-v2-homebrew-formula.md)
- `docs/maintainer/releasing.md` §3, §5; Slice-Vorlage `.harness/baseline/v6.13.0/templates/docs/plan/planning/slice.template.md`
- `AGENTS.md` (Hard Rules), `harness/review.md`

Kein Docker-Lauf, keine Änderung an bestehenden Dateien. Der Tag `v0.8.1` existiert lokal nicht (`git tag` endet bei `v0.8.0`).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Die Roadmap führt die beiden Slices der Umsetzung nirgends mehr: Der Cut ersetzt in „Aktuelle Welle" den Verweis auf [`slice-v1-sandbox-docker-ce`](../plan/planning/done/slice-v1-sandbox-docker-ce.md) durch den Cut-Slice, und unter „Abgeschlossene Wellen" gibt es keine Zeile für die Umsetzung von [ADR-0019](../plan/adr/0019-docker-daemon-aus-docker-apt-repo.md). Der Vorgänger-Cut legte die Zeile der Umsetzungs-Slices mit an (Zeile 81, Docker als verschachtelte Runtime); nach `releasing.md` §3 Nr. 4 bleibt nur der Wellen-Eintrag des Release bis zum Tag offen. Folge: Der geschlossene Slice ist von der Roadmap aus nicht auffindbar. | `docs/maintainer/releasing.md` §3 Nr. 4, `harness/conventions.md` (Roadmap-Form) | `docs/plan/planning/in-progress/roadmap.md:16`, `:81-82` | ja — Gegenlesen | Roadmap-Zeile für geschlossenen Slice fehlt |
| F-2 | LOW | Die Plan-Tabelle des Slice (§3) nennt vier Zeilen; der Commit ändert zusätzlich `docs/plan/planning/in-progress/README.md` (Slice-Zeile getauscht) und die beiden Reports `docs/reviews/2026-10-14-sandbox-docker-ce.md` und `-r2.md` (Slice-Link von `in-progress/` auf `done/`). Die Link-Korrektur ist richtig (Ziel `done/`-Datei des Slice (existiert) existiert), steht aber nicht im Plan. | Slice-Vorlage §3 (Plan vor Code) | Cut-Slice in `in-progress/`, Zeilen 46-51 | ja — `git show --stat 290ce07` | Plan-Tabelle unvollständig |
| F-3 | INFO | Die Kopfzeile der Roadmap trägt „Letzte Änderung: 2026-09-29", obwohl der Cut (und die Wellen davor) sie ändern. Bestand, nicht vom Cut verursacht; `releasing.md` §3 verlangt die Pflege nicht. | Maintainability | `docs/plan/planning/in-progress/roadmap.md:3` | ja — Gegenlesen | — |
| F-4 | INFO | „released 2026-10-14" in `README.md` und `README.de.md` sowie die Ziele `releases/tag/v0.8.1`, `download/v0.8.1` und der Compare-Link `v0.8.0...v0.8.1` lösen erst nach dem Tag-Push auf; vor dem Push ist die Aussage vorläufig (Praxis früherer Cuts). Das Datum `2026-10-14` steht an fünf Stellen fest (CHANGELOG, beide READMEs je Status und Tabelle, Handbuch „Stand") und stimmt nur bei Tag-Push am selben Tag. | `docs/maintainer/releasing.md` §3/§5/§7 | `README.md:10`, `README.md:206`, `README.de.md:11`, `README.de.md:210`, `CHANGELOG.md:14`, `docs/user/benutzerhandbuch.md:5` | ja — nach dem Push: `releasing.md` §6 | — |
| F-5 | INFO | Die Slice-DoD (§2) und die Closure-Notiz (§7) sind offen — vor dem Tag der Sollzustand. Mit diesem Report ist der Punkt „Review durchgeführt" erfüllbar; Tag, `publish`-Lauf, `make fullbuild`, `make test-docker`, Closure-Notiz und der `git mv` nach `done/` sind Maintainer-/Verifier-Aufgabe. Die Roadmap-Verschiebung der Welle nach „Abgeschlossene Wellen" folgt laut Plan nach dem Tag (`releasing.md` §7). | `harness/verification.md`, `docs/maintainer/releasing.md` §4, §7 | Cut-Slice in `in-progress/`, Zeilen 31-38, `:73-81` | nein — Verifier | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| CHANGELOG (§3 Nr. 1): `[Unreleased]` als leerer Anker (`CHANGELOG.md:12`), `[0.8.1] - 2026-10-14` mit Lead-Absatz (`:14-18`), Compare-Links `[Unreleased]: …v0.8.1...HEAD` und `[0.8.1]: …v0.8.0...v0.8.1` am Fuß | geprüft, ohne Befund |
| CHANGELOG-Inhalt: `### Changed` unter 0.8.1 stimmt mit der Umsetzung überein — Docker 29.9, Buildx 0.38, Compose 5.6 und die Altstände 26.1/0.13/2.26 decken sich mit [ADR-0019](../plan/adr/0019-docker-daemon-aus-docker-apt-repo.md) (Messung und Folgen); die vier Pakete (`docker-ce`, `containerd.io`, `docker-buildx-plugin`, `docker-compose-plugin`) und der Fingerprint-Pin stehen im Dockerfile-Template und im Test `generate_sandbox_test.go`; `download.docker.com` in der Default-Allowlist nur bei `docker` (`devcontainer_egress.go:69`, `spezifikation.md:819`, Handbuch `:1035`); der Block `[0.8.0]` nennt weiter `docker.io` als damaligen Stand (historisch korrekt) | geprüft, ohne Befund |
| Versionsstrings (§3 Nr. 2): `cmd/uboot/main.go:43`, `Makefile:27`, `Dockerfile:39` je `0.8.1-dev`; alle Kommentar-Erwähnungen (Dockerfile `:35`, `:37`; Makefile `:23`, `:26`) mitgezogen; `grep` nach `0.8.0-dev` findet nur die Beschreibung des Bumps im Slice-Plan (Cut-Slice, Zeile 49) und frühere Reports; keine weitere Quelle (Tests, Skripte, Workflows, Pakete) trägt `0.8.0` als aktuellen Wert | geprüft, ohne Befund |
| READMEs (§3 Nr. 3): Status-Block, Zeile `v0.8.1` in der Releases-Tabelle und `VERSION=0.8.1` in EN und DE konsistent zueinander, zum CHANGELOG und zu [LH-FA-DEV-007](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) und [ADR-0019](../plan/adr/0019-docker-daemon-aus-docker-apt-repo.md) (Anker lösen auf); die Zeile `v0.8.1` steht wie bei `v0.8.0` an der Bestandsposition (Fehlordnung der Tabelle als Altlast bekannt, nicht erneut gemeldet) | geprüft, ohne Befund |
| Benutzerhandbuch: `Software-Version: v0.8.1`, `Stand`, Beispielausgabe `u-boot version 0.8.1`, `download/v0.8.1`, `VERSION=0.8.1`, Gültigkeitsbereich konsistent; Beschreibung von `nestedRuntime: docker` (`:991`, `:1035`) deckt sich mit dem Code | geprüft, ohne Befund |
| Claims zu Docker Desktop: README, Releases-Zeilen, CHANGELOG und Slice-Abgrenzung sagen Docker Desktop nicht zu; das Handbuch nennt es unter „Nicht geprüft" (`:1089`) und nur als Beispiel für die VM bei einem Ausbruch (`:1002`); die Zusage des Lastenhefts (`spec/lastenheft.md:446`) ist als Risiko im Slice (§6) mit Folge-Slice [`slice-v1-sandbox-docker-nachholmessung-desktop`](../plan/planning/done/slice-v1-sandbox-docker-nachholmessung-desktop.md) getragen | geprüft, ohne Befund |
| Roadmap (§3 Nr. 4): „Aktuelle Welle" führt `v0.8.1` mit „Tag-Push steht aus", Link auf den Cut-Slice löst auf; offene Homebrew-Verifikation bleibt benannt (Befund F-1 zur Zeile der Umsetzung gesondert) | geprüft, ohne Befund |
| Slice gegen Vorlage (§3 Nr. 5): Kopf (Lifecycle, Welle, Bezug mit aktiven ADRs, Berührte Spec-Stellen „—" mit Begründung, Verantwortlich `pt9912`, Autor/Datum), §1–§8 in Vorlagenreihenfolge; §1 Ausschlüsse mit Begründung und existierenden Folge-Slices; §4 Rückführungen vorab benannt, Start-Trigger erfüllt (Slice [`slice-v1-sandbox-docker-ce`](../plan/planning/done/slice-v1-sandbox-docker-ce.md) liegt in `done/`); §6 beide Risiken tragen einen Ausgang; acht DoD-Punkte (Anzahl der Liefer-Punkte nicht beanstandet) | geprüft, ohne Befund |
| Slice-Lifecycle: Slice liegt nur in `in-progress/`, Eintrag in `in-progress/README.md` vorhanden, kein verwaister Verweis auf einen `in-progress/`-Pfad des geschlossenen Docker-CE-Slice mehr in `docs/plan`; ADR/Lastenheft verweisen nicht abwärts | geprüft, ohne Befund |
| Hard Rules: kein `//nolint`, kein `git mv`, keine Gate-/Schwellen- oder ADR-Änderung, `spec/architecture.md` unberührt; Kommentare der Versionsstrings beschreiben den Ist-Zustand (3.7); Commit-Message nennt [LH-FA-DEV-007](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) und [ADR-0019](../plan/adr/0019-docker-daemon-aus-docker-apt-repo.md) | geprüft, ohne Befund |

Nicht geprüft:

- `make gates`, `make ci`, `make test-docker`, `make fullbuild` — Auftrag verbietet Docker-Läufe; Verifier führt sie vor dem Tag aus (`docs/maintainer/releasing.md` §4). `make docs-check` wurde für diesen Report ausgeführt (Ergebnis siehe Übergabe).
- `publish`-Lauf, Tag-Hash, GHCR-/Release-Assets, Tap-Job — der Tag ist noch nicht gesetzt (`releasing.md` §6).
- Inhaltliche Richtigkeit der Feature-Aussagen gegen das erzeugte Dockerfile im Container — abgedeckt durch die Reports `2026-10-14-sandbox-docker-ce` und `-r2`, nicht Gegenstand dieses Laufs.
- `.github/workflows/` und `nfpm`-Konfiguration — im Diff nicht berührt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Roadmap-Zeile für geschlossenen Slice fehlt · Plan-Tabelle unvollständig

## Verdikt

**Merge-blockierend:** nein — kein HIGH/MEDIUM; die beiden LOW-Findings sind Doku-Nachzug und blockieren den Tag nicht. Der Cut ist mit `releasing.md` §3 Nr. 1–5 konsistent; offen bleiben die Maintainer-/Verifier-Schritte (§4 Gates, §5 Tag, §6 Kontrolle, Closure und `git mv` nach `done/`).

**Übergabe:** Findings gehen an den Implementer; F-1 lässt sich in der Nacharbeit nach dem Tag (`releasing.md` §7) mit der Wellen-Zeile erledigen; F-4 ist nach dem Tag-Push gegenzulesen. Der Report ersetzt keine Verification.
