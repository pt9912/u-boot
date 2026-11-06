# Review-Report: slice-v1-tap-nachzug-verfahren-doku — 2026-11-06

**Review-Art:** Code (Dokumentations-Diff) — geprüft gegen Plan, Entscheidungen und Hard Rules.

**Gegenstand:** [`slice-v1-tap-nachzug-verfahren-doku`](../plan/planning/in-progress/slice-v1-tap-nachzug-verfahren-doku.md), Commit `11fc3b01` gegen `0afd3272` (zwei Dateien:
`docs/maintainer/releasing.md`, `docs/maintainer/README.md`).

**Skill:** `.harness/skills/reviewer.md` @ `11fc3b01`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-11-06

**Eingangs-Kontext:**

- Slice-Plan [`slice-v1-tap-nachzug-verfahren-doku`](../plan/planning/in-progress/slice-v1-tap-nachzug-verfahren-doku.md) (DoD, §3 Plan-Tabelle, §4 Entscheidungen, §6 Risiken)
- [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) (Entscheidung 5, Smoke im Tap); [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung)
- `AGENTS.md` (Hard Rules, besonders §3.1, §3.7, §3.10), `harness/conventions.md`
- gelesene Fundstellen im Code: `packaging/homebrew-tap/.github/workflows/smoke.yml`, `scripts/tap-nachzug.sh`,
  `scripts/tap-check.sh`
- lesende Gegenprobe gegen den Tap `pt9912/homebrew-u-boot` (`gh api`, nur GET; kein Dispatch, kein Push)

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der Absatz „Gelaufen und nicht gelaufen“ behauptet im Indikativ, Schritt 3 sei „in dieser Form … gelaufen“ und habe einen Tap-Commit erzeugt; der Implementer hat keinen Push ausgeführt, die Herkunft (Bericht des Hauptagenten) steht im Absatz nicht. Zugleich sagt der Absatz, die `B64`-Zeile sei „nicht Teil dieses Laufs“, so dass „diese Form“ nicht der abgedruckte Block ist; „dieses Laufs“ und „ist ausgeführt“ tragen keinen Bezugspunkt (kein Datum, keine Revision). Ich habe nur gelesen, dass der Tap-Commit `662c5fb5` `smoke.yml` auf `main` geändert hat, nicht den Wortlaut des Aufrufs. | `AGENTS.md` §3.10 (Aussage über Verhalten nennt nur, was gemessen oder gelesen wurde; Zustandsmarker ohne Ablaufpunkt) | `docs/maintainer/releasing.md` · „in dieser Form (mit vorher gesetzten Variablen, ohne `branch`) gelaufen“ | Review-only (der Wortlaut des Push-Aufrufs ist nicht prüfbar; Plan §4 Entscheidung 4 verlangt „ungeprobt“, wenn er nicht vorliegt) | Beleg-Zustand einer Probe ohne Herkunft und Ablaufpunkt in der Handlungsanleitung |
| F-2 | LOW | Schritt 3 liest die Quelle mit `base64 -w0` aus dem Arbeitsbaum; `-w0` ist eine GNU-Option, und der Kommentar „Base64 der Quelle ohne Zeilenumbruch“ sagt es als Eigenschaft des Befehls. Die Grenze nennt die Plattform-Abhängigkeit nicht, obwohl Plan §6 Risiko 2 den Ausgang „weiter offen → Grenze (d)“ vorsieht. Im Devcontainer (Linux) lief die Zeile hier. | `AGENTS.md` §3.10; Plan §6 Risiko 2 | `docs/maintainer/releasing.md` · „B64=$(base64 -w0 < packaging/homebrew-tap/$P)“ | Review-only | Plattform-Annahme eines ungeprobten Befehls nicht in der Grenze ausgewiesen |
| F-3 | LOW | Schritt 3 und 4 lesen beide `packaging/homebrew-tap/$P` aus dem Arbeitsbaum, nicht aus dem Commit auf `main`, den Schritt 2 als geprüft vorsieht. Bei lokal abweichendem Arbeitsbaum oder anderem Branch pusht Schritt 3 ungeprüften Inhalt, und der `cmp` in Schritt 4 vergleicht gegen dieselbe Quelle und merkt es nicht. | Plan DoD „Verfahren“ (c) Schritt 2 („die Quelle führt, der Tap folgt“); `BEO-DOK/unumkehrbarer-schritt-vor-der-kontrolle` | `docs/maintainer/releasing.md` · „Die Befehle setzen `P` (der Pfad der Datei) und lesen die Quelle aus dem Arbeitsbaum“ | Review-only | Kontrolle (Review und Verification) und gepushter Inhalt nicht auf denselben Stand gebunden |
| F-4 | INFO | Plan DoD (c) Schritt 3 nennt `branch` als Parameter; der Abschnitt lässt ihn weg und sagt es im Absatz „Gelaufen und nicht gelaufen“. Das ist offengelegt; der Push geht dann auf den Standard-Branch des Tap. | Plan DoD „Verfahren“ (c) | `docs/maintainer/releasing.md` · „-f message='<Commit-Text>' -f content="$B64" -f sha="$SHA"“ | Review-only | Abweichung vom Plan-Wortlaut, offengelegt |
| F-5 | INFO | Schritt 5 lässt Annotationen lesen; ein grüner Lauf trägt hier trotzdem eine Annotation der Stufe `failure` („Failed to import: …/hashicorp/homebrew-tap/Formula/vagrant.rb“) und zwei Warnungen aus dem Runner-Image (gelesen am Lauf `42212837961`, Job `success`). Der Abschnitt nennt nur die Annotation „not supported“ als Erkennungszeichen und nicht, dass Annotationen allein kein Rot belegen. | `BEO-DOK/ci-log-falsch-gelesen` (Klasse nahe) | `docs/maintainer/releasing.md` · „Ein Job ohne Runner trägt die Annotation“ | ja — `gh api repos/pt9912/homebrew-u-boot/check-runs/<job-id>/annotations` | Annotations-Rauschen als Rot-Signal lesbar |
| F-6 | INFO | Zeilenangaben in den Lesebefunden (`smoke.yml` Zeilen 10 bis 14; `tap-nachzug.sh` Zeilen 25 bis 27 und 34; `tap-check.sh` Zeile 34) sind gegen den Stand `11fc3b01` richtig (gelesen), gelten aber für Dateien, die ohne Test geändert werden dürfen; die erste ist die Datei, die der Abschnitt selbst zu ändern beschreibt. Der Plan sieht Zeilen als Fundstelle vor. | Plan §6 Risiko 4; `AGENTS.md` §3.10 | `docs/maintainer/releasing.md` · „(`smoke.yml` Zeilen 10 bis 14:“ | Review-only | Zeilen-Fundstelle ohne Sensor |
| F-7 | INFO | Plan §4 Entscheidung 2 macht den Sichtschritt in §4 von der Entscheidung des Projektinhabers abhängig; der Satz steht im Diff, die Entscheidung liegt dem Reviewer nicht als Artefakt vor. Zuständigkeit Verifier/Closure. | Plan §4 Entscheidung 2 | `docs/maintainer/releasing.md` · „Sichtschritt vor einem stabilen Tag, kein Gate“ | Review-only | Entscheidung ohne Artefakt |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Diff-Umfang: nur `docs/maintainer/releasing.md` und `docs/maintainer/README.md`; nichts unter `packaging/`, `scripts/`, `.github/` | geprüft, ohne Befund (`git diff --stat 0afd3272 HEAD`) |
| Plan-Tabelle §3 gegen Diff (beide Dateien als „update“ geführt) | geprüft, ohne Befund |
| Auffindbarkeit: §2, §4, §8 (Tabellenzeile) und Index je höchstens ein Satz, alle zeigen auf §9; Nummerierung §1 bis §8 unverändert | geprüft, ohne Befund |
| Fundstellen-Aussagen über `smoke.yml` (Trigger `push` auf `Formula/u-boot.rb`, `workflow_dispatch`, Name `smoke`), `scripts/tap-nachzug.sh` (Klon, `cp`, `add` nur `Formula/u-boot.rb`) und `scripts/tap-check.sh` (`git push --dry-run`) | gelesen, stimmen mit dem Text überein |
| Aussage „Der Push startet keinen Lauf“ | gegen die Lauf-Liste des Tap gelesen: auf den Commit `662c5fb5` (nur `smoke.yml`) liegt kein `push`-Lauf, nur ein `workflow_dispatch`; ohne Befund |
| Lesende Befehle der Schritte 1, 4, 5 (sha-Vergleich, `cmp`, Lauf-Liste, Jobs mit `labels`/`runner_name`/`conclusion`, Annotationen, Log mit `--allow-escape-sequences` und `grep`) | gegen den Tap ausgeführt (nur GET): sha und `cmp` für `README.md` und `smoke.yml` gleich der Quelle (Exit 0); Befehle liefern die beschriebenen Spalten, der `grep` trifft `Image:`, `Cellar/u-boot/0.24.0:` und `u-boot version 0.24.0`; ohne Befund |
| Erkennung des neuen Laufs über `created_at`/`head_sha` und ein älterer Dispatch-Lauf ohne Job | gegen die Lauf-Liste gelesen (Lauf `42196563230` `queued`, Lauf `42212837961` `success`, gleicher `head_sha`); ohne Befund |
| Aussage „Kein Skript, kein Gate und kein Test hält Quelle und Tap gleich“ | per Suche nach `homebrew-tap` und `smoke.yml` außerhalb `docs/` gegengelesen: einzige Fundstelle `packaging/homebrew-tap/README.md`; ohne Befund (Stand des Laufs, ohne Sensor) |
| Rechte-Hinweis: Zitat mit Quelle, keine Aussage über ein bestimmtes Token (§9 sagt ausdrücklich nichts über `HOMEBREW_TAP_GITHUB_TOKEN`) | geprüft, ohne Befund (Plan §6 Risiko 3) |
| Keine Run-, Commit- oder Job-IDs im Abschnitt | geprüft, ohne Befund |
| Kommentar-Selbstcheck §3.9 (a) bis (c) | Diff ist Prosa, kein Code-Kommentar; (a): Absatz „Gelaufen und nicht gelaufen“ siehe F-1; sonst ohne Befund |
| Hard Rules §3.1 (Docker-only), §3.2, §3.3, §3.4, §3.5, §3.6; Referenzrichtung; Nutzer-Doku-Regel (§5 Nr. 6) | geprüft, nicht berührt, ohne Befund |
| Sprachvertrag ([`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache)) | nicht berührt (Maintainer-Doku, Deutsch) |
| `make test-docker`, `make verify-depguard` | nicht gelaufen: kein Docker-Pfad und keine Importregel berührt |
| Inhalt des Zitats „Create or update file contents“ der GitHub-REST-Dokumentation | nicht geprüft: kein Zugriff auf die Seite aus dieser Rolle; Zitat nach Plan (gelesen 2026-11-06 vom Planner) übernommen, nicht von mir bestätigt |
| `make gates` / Wirkung der Befehle `PUT` (Schritt 3), Dispatch (Schritt 5), Rückweg (Schritt 6) | nicht gelaufen: Schreibzugriffe auf den Tap sind vom Auftrag ausgeschlossen; Gate-Nachweis ist Verifier-Aufgabe |
| `make docs-check` für diesen Report | gelaufen im Devcontainer: 1211 Dateien, 0 Befunde (nach Verlinkung der Kennungen im Report; der erste Lauf meldete 5 `id-unlinked`) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Beleg-Zustand einer Probe ohne Herkunft und Ablaufpunkt in der Handlungsanleitung ·
Plattform-Annahme eines ungeprobten Befehls nicht in der Grenze ausgewiesen · Kontrolle und gepushter Inhalt nicht auf
denselben Stand gebunden · Annotations-Rauschen als Rot-Signal lesbar

## Verdikt

**Merge-blockierend:** ja, wegen F-1 (MEDIUM). F-2 und F-3 sind LOW; F-4 bis F-7 INFO ohne Änderungspflicht.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen zusätzlich in die Slice-Closure §7. F-7 und
die Lesebefunde im Negativ-Block gehen als Hinweis an den Verifier (er führt die lesenden Befehle selbst aus und
gleicht die Plan-Tabelle gegen `git diff --stat` ab). Der Report ersetzt keine Verifikation.
