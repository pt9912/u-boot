# Slice V1: Tap-Smoke auf `macos-15` (arm64) statt auf dem außer Betrieb genommenen `macos-14`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle (kleine Korrektur an der CI des Tap-Repos, kein Release-Cut; es gibt
keine Closure-Bedingung, die von der DoD dieses Slice verschieden ist).

**Anlass.** Der Lauf `smoke` im Tap-Repo `pt9912/homebrew-u-boot` zu `u-boot 0.24.0` (Run
`42191356009`, Event `push`, Kopf-Commit `b9f23331a5` „u-boot 0.24.0“, `run_attempt` 2, Stand
`cancelled`; gelesen am 2026-11-06 mit `gh api`) hat zwei Jobs: `install (macos-15-intel)` (Job
`127811577331`) ist `success`; `install (macos-14)` (Job `127811574012`) ist `cancelled`, trägt das
Label `macos-14` und hat keinen Runner bekommen (`runner_name` leer). Die Annotation des Jobs lautet
„The configuration 'macos-14-arm64-us-default' is not supported“. GitHub nennt in der Ankündigung
[„GitHub Actions: macOS 14 runner image retirement“](https://github.blog/changelog/2026-10-01-github-actions-macos-14-runner-image-retirement)
(gelesen am 2026-11-06) die Außerbetriebnahme des macOS-14-Images zum 2026-11-02 für die Labels
`macos-14`, `macos-14-large` und `macos-14-xlarge` und empfiehlt als arm64-Labels unter anderem
`macos-15` und `macos-latest`. Die README von `actions/runner-images` (gelesen am 2026-11-06) führt
`macos-15` als „macOS 15 Arm64“, `macos-15-intel` als „macOS 15“ für x64 und `macos-latest` als
„macOS 26 Arm64“. Dass die Außerbetriebnahme die Ursache dieses Laufs ist, folgt aus der
Annotation, dem fehlenden Runner und der Ankündigung; bewiesen ist es erst, wenn der Lauf mit dem
neuen Label grün ist. Das Label `macos-14` bekam bis zum 2026-11-04 noch Runner: die Läufe
`41636802191` (2026-11-03) und `41844822254` (2026-11-04) waren dort `success`. Die Matrix steht in
`packaging/homebrew-tap/.github/workflows/smoke.yml` Zeile 24 (`runner: [macos-14, macos-15-intel]`);
die Datei im Tap ist heute byte-gleich mit ihr (`cmp` gegen den Inhalt von
`gh api repos/pt9912/homebrew-u-boot/contents/.github/workflows/smoke.yml`, 2026-11-06). Der
Projektinhaber hat am 2026-11-06 entschieden, den Runner zu wechseln.

**Bezug:** [`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung) (Paketierung: Homebrew
gewählt, eigener Tap), [ADR-0016](../../adr/0016-homebrew-distribution-per-tap.md) (Entscheidung 5: Smoke im
Tap).

**Berührte Spec-Stellen:** — (die Smoke-Matrix steht in keiner Stelle von Lastenheft, Spezifikation oder
Sicht; [ADR-0016](../../adr/0016-homebrew-distribution-per-tap.md) Entscheidung 5 nennt `smoke.yml` und seine
Prüfungen, kein Runner-Label).

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner, im Auftrag des Hauptagenten). **Datum:** 2026-11-06.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Die Matrix des Smoke-Workflows für den Homebrew-Tap nennt für den arm64-Job das Label
`macos-15`; die Quelle `packaging/homebrew-tap/.github/workflows/smoke.yml` und die Datei im Tap-Repo
`pt9912/homebrew-u-boot` sind dieselbe Datei, und ein Lauf `smoke` im Tap zeigt `install (macos-15)`
und `install (macos-15-intel)` grün.

**Planentscheidung: `macos-15` (arm64), ausdrücklich, nicht `macos-latest`.** Die Wahl des Labels ist
ein Vorschlag des Planners; die Freigabe zum Umsetzen gab der Projektinhaber am 2026-11-06.

- **Grund für `macos-15`:** Beide Labels der Matrix nennen die macOS-Version ausdrücklich (arm64
  `macos-15`, Intel `macos-15-intel`), beide Jobs laufen auf derselben Hauptversion, und ein Wechsel der
  Version ist eine Änderung dieser einen Zeile, die ein Review sieht. Nach der README von
  `actions/runner-images` zeigt `macos-latest` heute auf macOS 26 (arm64) und wandert bei einem Wechsel
  über ein bis zwei Monate (README, Abschnitt zum `-latest`-Label); der arm64-Job liefe dann ohne Commit
  im Repo auf einer anderen Version als der Intel-Job, und ein Rot durch die neue Version
  (Homebrew-Audit, vorinstallierte Taps des Images) wäre an keiner Änderung festzumachen.
- **Gegenposition: `macos-latest` für den arm64-Job.** Das Label wird laut README erst nach Ankündigung
  und mit Vorlauf auf eine neue Version verschoben, es wird nicht stillgelegt: Der Anlass dieses Slice
  ist ein ausdrückliches Label, das GitHub stilllegt, und es ist der zweite Fall nach `macos-13` (Intel; [`slice-v2-homebrew-formula`](../done/slice-v2-homebrew-formula.md),
  Abschnitt „Stand nach `v0.9.0`“). Mit `macos-latest` wäre kein Eingriff in die Datei nötig, und der
  Smoke liefe weiter. Der Preis steht oben: die Version wandert ohne Commit, und die beiden Jobs
  laufen auf verschiedenen Versionen.
- **Neubewertung, wenn:** ein Label der Matrix ein drittes Mal durch Außerbetriebnahme rot wird, oder
  GitHub die Stilllegung von `macos-15` ankündigt. Dann prüft der Planner `macos-latest` für den
  arm64-Job gegen diese Wahl und legt das Ergebnis als eigenen Slice vor.
- **Keine ADR:** [ADR-0016](../../adr/0016-homebrew-distribution-per-tap.md) Entscheidung 5 nennt `smoke.yml`
  und seine Prüfungen, kein Label; keine Accepted-ADR wird berührt, und die Wahl ist mit einer Zeile
  umkehrbar.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Kein Wechsel von `macos-15-intel`.** Der Intel-Job ist im Lauf `42191356009` `success`; die README
  von `actions/runner-images` führt das Label ohne Hinweis auf eine Außerbetriebnahme (der gelesene
  Hinweis dort gilt macOS 14). Ob und wann `macos-15-intel` endet, ist nicht gelesen (Risiko in §6).
- **Keine Änderung an Schritten, Audit-Gate oder Formel.** `brew audit --strict --online` gatet den
  Lauf (`smoke.yml` Zeilen 58 bis 66) und wird nicht gelockert (`AGENTS.md` §3.6); ein neuer Befund
  unter `macos-15` wäre ein eigener Vorgang, kein Teil dieses Slice.
- **Kein Test und kein Gate, das die Labels pinnt.** Ob GitHub ein Label bedient, ist Fremdzustand und
  offline nicht lesbar; ein Test, der das Label aus der Datei liest, prüft die Datei gegen sich selbst.
  Das Wiederauftreten fängt das Beobachtungs-Register (§7).
- **Kein Nachzug-Mechanismus und keine Maintainer-Doku für `smoke.yml`.** `scripts/tap-nachzug.sh`
  schreibt nur `Formula/u-boot.rb` in den Tap; wie eine geänderte Vorlage dorthin kommt (Kopie,
  Vergleich), steht in keiner Maintainer-Doku (`packaging/homebrew-tap/README.md` Zeilen 19 bis 21 nennt
  nur das Anlegen des Tap-Repos). Ein Verfahren oder eine Automatik ist ein anderer Vorgang; es gibt dazu
  noch keine Slice-Datei, die Frage steht im Bericht an den Projektinhaber.
- **Kein CHANGELOG-Eintrag, kein Release, kein Tag.** Die CI des Tap-Repos ist nicht das ausgelieferte
  Produkt: Binary, Formel und Pakete ändern sich nicht. `CHANGELOG.md` Zeile 454 (Abschnitt `[0.10.0]`)
  nennt `macos-14` und `macos-15-intel` als Stand dieses Release und bleibt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — drei Liefer-Punkte; die Gate-Läufe und die Closure-Pflichten zählen nicht mit.

- [ ] **Quelle:** `packaging/homebrew-tap/.github/workflows/smoke.yml` nennt in der Matrix
      `runner: [macos-15, macos-15-intel]`; der Diff der Datei besteht aus dieser Zeile und einem
      Kommentar zur Label-Wahl (§3), sonst nichts.
- [ ] **Tap:** `.github/workflows/smoke.yml` in `pt9912/homebrew-u-boot` ist byte-gleich mit der Quelle:
      `gh api repos/pt9912/homebrew-u-boot/contents/.github/workflows/smoke.yml --jq .content | base64 -d | cmp - packaging/homebrew-tap/.github/workflows/smoke.yml`
      endet mit Exit 0; der Tap-Commit steht in §7. Den Push in den Tap macht der Hauptagent nach Review
      und Verification, nicht der Implementer; die Freigabe des Projektinhabers liegt für diesen Schritt vor
      (Auftrag vom 2026-11-06).
- [ ] **Smoke:** Ein Lauf `smoke` im Tap, gestartet mit `gh workflow run smoke --repo pt9912/homebrew-u-boot`
      (der Push-Trigger beobachtet nur `Formula/u-boot.rb`, `smoke.yml` Zeilen 10 bis 13; der Lauf
      `39008767100` in [`slice-v2-homebrew-smoke-haertung`](../done/slice-v2-homebrew-smoke-haertung.md) lief
      ebenso über `workflow_dispatch`), endet mit `success`, und die Jobs `install (macos-15)` und
      `install (macos-15-intel)` sind beide `success`. Run-ID und beide Job-IDs stehen in §7. Der Beleg, dass
      der erste Job ein arm64-Runner war, kommt aus den Job-Metadaten (`labels` ist `["macos-15"]`,
      `runner_name` nicht leer) und dem Image-Abschnitt des Schritts „Set up job“ im Log; beides wird mit
      `gh api` gelesen, nicht allein mit `gh run view --log`
      ([`BEO-DOK/ci-log-falsch-gelesen`](../observations/BEO-DOK/ci-log-falsch-gelesen/observation.md)). Ist
      das Log nicht lesbar, steht das in §7, und die Aussage „arm64“ steht dort nicht. Die Aussage „grün“
      nennt die Formel-Version, die im Lauf installiert wurde (`0.24.0`, solange kein neuer Tag im Tap liegt),
      und kein „macOS grün“ darüber hinaus.
- [ ] `make gates` grün (im Devcontainer-Container, `AGENTS.md` §3.1).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [ ] Verification-Report ([`harness/verification.md`](../../../../harness/verification.md)) liegt vor: vor dem
      Push in den Tap für die Quelle, die Gates und die Plan-Tabelle; die Punkte „Tap“ und „Smoke“ prüft der
      Verifier nach dem Push nach, nur lesend.
- [ ] Doku-Update: keine Datei berührt, weil kein öffentlicher Vertrag berührt ist; die Zeilen „unverändert“
      in §3 sind mit gelesenen Stellen belegt, bei der Closure gegen `git diff --stat` abgeglichen.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis
      `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**,
      er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7
      notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — hier geprüft, weil das Repo ohne
      Wellen-Betrieb für diesen Slice arbeitet.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md` §Was ist eine Sub-Area? — die Liste
liefert die Pfad-Kandidaten für §8. Jede Zeile „unverändert“ trägt die Stelle, die gelesen wurde
(`AGENTS.md` §3.10).

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `packaging/homebrew-tap/.github/workflows/smoke.yml` | update | Zeile 24: `macos-14` wird `macos-15`; davor ein Kommentar (Abgrenzung, Indikativ, ohne Chronik): beide Labels nennen die macOS-Version ausdrücklich (arm64 `macos-15`, amd64 `macos-15-intel`), ein Versionswechsel ist eine Änderung dieser Zeile in Vorlage und Tap, kein wanderndes `-latest`. Wortlaut ASCII wie die übrigen Kommentare der Datei (Zeilen 3 bis 8). |
| Tap-Repo `pt9912/homebrew-u-boot`, `.github/workflows/smoke.yml` | update (extern, nicht in `git diff --stat`) | Kopie der Quelle nach dem Commit der Quelle; Push durch den Hauptagent nach Review und Verification (DoD „Tap“). Der Tap-Branch `main` ist nicht geschützt (`gh api repos/pt9912/homebrew-u-boot/branches/main`, `protected` ist `false`, 2026-11-06). |
| `docs/plan/planning/in-progress/slice-v1-tap-smoke-runner-macos-15.md` | neu | dieser Plan; die Zustandswechsel sind reine `git mv` (open → next → in-progress → done) |
| `docs/plan/planning/open/README.md` | update | Bestandszeile des Plans (Tabelle „Bestand“); die Zeile folgt dem Slice bei jedem `git mv` |
| `docs/reviews/<Datum>-slice-v1-tap-smoke-runner-macos-15.md` | neu | Review-Report (`.harness/skills/reviewer.md`) |
| `docs/reviews/<Datum>-verification-slice-v1-tap-smoke-runner-macos-15.md` | neu | Verification-Report; Namensform wie die Berichte zu `slice-v2-homebrew-smoke-haertung` |
| `docs/plan/planning/observations/BEO-<KUERZEL>/<slug>/…` | neu oder update, bei der Closure | Beobachtung „Runner-Label stillgelegt“ mit den Vorgängen `slice-v2-homebrew-formula` (`macos-13`) und diesem Slice (`macos-14`), falls die Closure sie anlegt; das Kürzel ist nachzuschlagen (§8) |
| Tests | keine Änderung | Kein Test liest `smoke.yml`: `cmd/uboot/homebrew_scripts_test.go` liest `scripts/homebrew-formula.rb.tmpl` (Zeile 73) und ruft `scripts/homebrew-formula-fill.sh` und `scripts/tap-formula-check.sh` (Funktionen `runFill`, Zeilen 21 bis 35, und `runTapCheck`, Zeilen 120 bis 129); Zeile 70 nennt `slice-v2-homebrew-smoke-haertung` nur im Kommentar. Eine Suche nach `smoke.yml`, `macos-1` und `runner:` über `*.go`, `*.sh`, `*.mk`, `*.py`, `*.yml`, `*.yaml`, `*.rb` und `*.tmpl` fand außer der Datei selbst nur den Fehlertext `fake runner: …` in `internal/adapter/driven/docker/inspector_test.go` Zeile 47. |
| `.github/workflows/*.yml` (`ci`, `integration`, `publish`, `tap-check`) | unverändert | Keine Datei nennt ein `macos`-Label (Suche über `.github/workflows/`); der Smoke läuft im Tap-Repo, nicht im Hauptrepo. |
| `scripts/tap-nachzug.sh`, `scripts/tap-check.sh` | unverändert | `tap-nachzug.sh` klont den Tap und kopiert nur das Release-Asset nach `Formula/u-boot.rb` (Schritte `mkdir -p "$work/tap/Formula"` und `cp "$work/u-boot.rb" "$work/tap/Formula/u-boot.rb"`); `smoke.yml` kommt nie per Nachzug in den Tap. `tap-check.sh` probt den Push des Nachzug-Tokens per Klon und `git push --dry-run`; mit dem Token des Hauptagenten hat es nichts zu tun. |
| `packaging/homebrew-tap/README.md` | unverändert | Gelesen Zeilen 1 bis 21: nennt `.github/workflows/smoke.yml` (Zeile 21) als Datei, die beim Anlegen des Tap-Repos hineinkopiert wird, und kein Runner-Label. |
| `docs/maintainer/releasing.md` | unverändert | Gelesen §2 (Zeilen 50 bis 51: „README und Smoke-Workflow aus `packaging/homebrew-tap/`“) und §6 (Zeilen 162 bis 164: Schritte des Smoke, Audit gatet); kein Runner-Label, keine Plattform-Zusage des Smoke. |
| `CHANGELOG.md` | unverändert | Gelesen: `[Unreleased]` leer (Zeile 12, der nächste Abschnitt `[0.24.0]` beginnt in Zeile 14); Zeile 454 steht im Abschnitt `[0.10.0]` (Zeile 423) und beschreibt den Stand dieses Release. |
| `docs/plan/adr/0016-homebrew-distribution-per-tap.md` | unverändert | Accepted; Entscheidung 5 (Zeile 29) nennt `smoke.yml` und seine Prüfungen, kein Label. |

**Ansatz in der Reihenfolge der Rollen:** Implementer ändert die Quelle (eine Zeile und ein Kommentar) und
läuft `make gates` → Reviewer → Verifier (Quelle, Gates, Plan-Tabelle gegen `git diff --stat`) →
Hauptagent kopiert die Quelle in den Tap, committet, pusht, vergleicht mit `cmp`, startet den Lauf
(`gh workflow run smoke --repo pt9912/homebrew-u-boot`) und liest ihn → Verifier-Nachlauf nur lesend →
Closure. Der Tap-Push steht nach Review und Verification, weil ein Push in ein fremdes Repo zwar per
Folge-Commit umkehrbar ist, aber einen Lauf in dessen CI auslöst.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Die Freigabe zum Umsetzen liegt vor (Projektinhaber, 2026-11-06); der
Hauptagent beginnt mit dem Implementer-Lauf, sobald in `in-progress/` kein anderer Slice liegt
(WIP-Limit 1). Die Übergänge `open` → `next` → `in-progress` sind je ein reiner `git mv`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Braucht die Änderung mehr als die Matrix-Zeile
  und den Kommentar (etwa einen geänderten Schritt, weil er unter `macos-15` nicht läuft), ist der
  Zuschnitt falsch; der Mehrbedarf wird ein eigener Slice.
- `in-progress` → `open` (blockiert): Verweigert GitHub den Push von `.github/workflows/smoke.yml` in den
  Tap (fehlendes Recht), oder startet der Lauf nicht, kehrt der Slice zurück, bis der Projektinhaber das
  Recht gibt. Ein Carveout entsteht dabei nicht: kein gelockertes Gate, sondern ein fehlender Zugang.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln.

DoD vollständig (drei Liefer-Punkte plus Gate- und Closure-Pflichten); Review und Verification (mit dem
Nachlauf für Tap und Smoke); Closure-Notiz mit Lerneintrag; `git mv` nach `done/`. Die Closure liest den
Tap-Lauf selbst (`gh api`, nur lesend) und wiederholt den Vergleich Quelle gegen Tap nach der letzten
Fix-Runde (`AGENTS.md` §3.10); was sie selbst gelesen hat und was nur der Hauptagent berichtet, trennt
§7. Lerneintrag, soweit er ansteht: die Beobachtung „ein Runner-Label der Smoke-Matrix wird von GitHub
stillgelegt“ steht zweimal im Bestand (`macos-13` in
[`slice-v2-homebrew-formula`](../done/slice-v2-homebrew-formula.md), `macos-14` hier) und wird als
**gezählt, nicht verkörpert** eingetragen, solange sie unter drei bleibt.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst. Der Ausgang steht bei der Closure; vorab steht, woran er
sich entscheidet.

- **Die Ursache ist nicht bewiesen.** Annotation, fehlender Runner und Ankündigung passen zusammen
  (§Anlass), aber der Lauf mit dem neuen Label ist der Beweis. Bleibt auch `macos-15` rot mit einer
  Runner-Meldung, ist die Ursache eine andere. — **Ausgang:** bei Closure (entfallen, wenn der Lauf grün
  ist; sonst eingetreten und zurück nach `open/`).
- **Weiterer Runner-Wandel.** Die Matrix steht nach dem Wechsel auf `macos-15` und `macos-15-intel`; beide
  Versionen werden einmal stillgelegt, und die README von `actions/runner-images` nennt für
  `macos-15-intel` keinen Termin (gelesen, 2026-11-06; ein Termin außerhalb dieser Datei ist nicht
  gelesen). — **Ausgang:** weiter offen → Beobachtung im Register (§5, §7 der Closure); der
  Neubewertungs-Trigger der Label-Wahl steht in §1.
- **Die Stilllegung fällt erst beim nächsten Formel-Push auf.** `smoke.yml` läuft nur bei einem Push auf
  `Formula/u-boot.rb` und per `workflow_dispatch` (Zeilen 10 bis 14, kein `schedule`); der rote Lauf zu
  `0.24.0` kam vier Tage nach dem Stilllegungsdatum und nur, weil ein Release ihn auslöste. Ein
  regelmäßiger Lauf wäre eine Zusage über Runner-Minuten und eine Änderung am Verfahren, nicht an diesem
  Slice. — **Ausgang:** weiter offen → Frage an den Projektinhaber (Bericht), bei Bedarf eigener Slice.
- **Schreibzugriff auf den Tap.** Der Push einer Datei unter `.github/workflows/` braucht bei GitHub ein
  eigenes Recht neben dem Schreibrecht auf den Inhalt (aus dem Verhalten von GitHub bekannt, in diesem Repo
  nicht gemessen). Gelesen: das Token des Hauptagenten ist ein fine-grained PAT (`gh auth status`), das
  Repo meldet `push` und `admin` als wahr, `main` ist nicht geschützt; frühere Rollouts von `smoke.yml`
  in den Tap (Commits `a8498fd`, `76a7faa`, genannt in
  [`slice-v2-homebrew-formula`](../done/slice-v2-homebrew-formula.md) und
  [`slice-v2-homebrew-smoke-haertung`](../done/slice-v2-homebrew-smoke-haertung.md)) liefen, mit welchem
  Token, ist nicht gelesen. Die Messstelle ist der Push selbst. — **Ausgang:** bei Closure (entfallen oder
  eingetreten → Rückführung nach `open/`, §4).
- **Datei-Drift zwischen Quelle und Tap.** Heute byte-gleich (§Anlass); danach kann jede Seite allein
  geändert werden, und kein Skript schreibt `smoke.yml` in den Tap. Der Vergleich in der DoD belegt den
  Stand zur Closure, nicht danach. — **Ausgang:** bei Closure (der Vergleich nach der letzten Fix-Runde);
  das Verfahren für künftige Änderungen bleibt weiter offen → Frage an den Projektinhaber (Bericht).
- **Ein neuer Befund unter `macos-15`.** Ein anderes Image bringt eine andere Homebrew-Version und andere
  vorinstallierte Taps mit; `brew audit --strict --online` kann etwas melden, das unter `macos-14` nicht
  auftrat, und der Hinweis „Failed to import … vagrant.rb“ des Intel-Laufs (Angabe des Projektinhabers im
  Cut-Slice [`slice-v1-release-cut-v0.24.0`](../done/slice-v1-release-cut-v0.24.0.md) §7, nicht von mir
  gelesen) kann auf dem neuen Image andere Gestalt haben. Das Audit-Gate bleibt (`AGENTS.md` §3.6). —
  **Ausgang:** bei Closure (eingetreten → eigener Slice für die Formel oder den Schritt; dieser Slice
  schließt erst mit grünem Lauf).
- **Knappe macOS-Runner.** Ein Job kann lange `queued` stehen (Lauf `36894925908` in
  [`slice-v2-homebrew-formula`](../done/slice-v2-homebrew-formula.md)). Ein Hänger ist kein Rot der Datei. —
  **Ausgang:** bei Closure (Rerun des Jobs, die Lauf-ID bleibt in §7).

## 7. Closure-Notiz

Noch nicht geschlossen; der Slice liegt in `open/` (Entwurf). Die Closure trennt bei den Belegen, was der
Planner selbst gelesen hat (Quelle, Vergleich, Tap-Lauf per `gh api`) und was nur der Hauptagent berichtet.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

**Vorgelagert — Sub-Area-Wahl prüfen:** `packaging/homebrew-tap/` und `.github/workflows/` stehen in keiner
Zeile der Modus-Deklaration in [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area).
Der Slice berührt eine Datei in `packaging/`; sie hat keine eigene Konvention und keine eigenständige
Inventur-Linie — wie `cmd/uboot` eine Sub-Area-Aspirantin, nicht angelegt. Die Planungs-Artefakte liegen
in der Sub-Area DOK (Greenfield). Das Kürzel einer Beobachtung zu diesem Gegenstand wird nachgeschlagen,
nicht erfunden: keine Zeile passt wörtlich; WZG (`tools/`, `scripts/`) ist die nächste, die Entscheidung
trifft die Closure.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen am 2026-11-06 (alle Kürzel
aufgelistet, die Einträge unten gelesen); die Stände nennen die Zahl der Evidence-Dateien.

- [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md)
  (verkörpert in `AGENTS.md` §3.10): angewandt — jede Aussage in §3 und §Anlass steht mit gelesener
  Stelle; der Abgleich gegen `git diff --stat` wird nach der letzten Fix-Runde wiederholt.
- [`BEO-DOK/plan-annahme-ueber-fremdzustand-ohne-messstelle`](../observations/BEO-DOK/plan-annahme-ueber-fremdzustand-ohne-messstelle/observation.md)
  (3×, offen, Kandidat für den Planner-Schnitt): trifft — „GitHub bedient `macos-15` als arm64“ und die
  Ursache sind Annahmen über ein Fremdsystem. Die Messstelle steht im Plan: Job-Metadaten und Image-Abschnitt
  des Laufs (DoD „Smoke“), die Ursache entscheidet der grüne Lauf (§6). Ob die Closure das als Anwendung
  ohne Beleg vermerkt, hängt am Ergebnis.
- [`BEO-DOK/ci-log-falsch-gelesen`](../observations/BEO-DOK/ci-log-falsch-gelesen/observation.md) (1×, offen):
  trifft — die DoD liest den Lauf mit `gh api` (Jobs, Annotationen, Log), nicht allein mit
  `gh run view --log`.
- [`BEO-DOK/unumkehrbarer-schritt-vor-der-kontrolle`](../observations/BEO-DOK/unumkehrbarer-schritt-vor-der-kontrolle/observation.md)
  (1×, offen): verwandt — der Push in ein fremdes Repo steht nach Review und Verification (§3, Ansatz); er
  ist per Folge-Commit umkehrbar, löst aber einen Lauf in der CI des Tap aus. Keine Zählung.
- [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)
  (über der Schwelle, offen): trifft — die Aussage zum Lauf nennt Formel-Version und Job-IDs, kein
  „macOS grün“ darüber hinaus (DoD „Smoke“).
- [`BEO-WZG/kopplung-nicht-nachgezogen`](../observations/BEO-WZG/kopplung-nicht-nachgezogen/observation.md)
  (2×, offen): verwandt, nicht gleich — dort hängt ein Skript an der Form eines Erzeugers, hier steht
  dieselbe Datei an zwei Orten. Nicht gezählt.
- Eine Beobachtung zu einem stillgelegten Runner-Label gibt es nicht (die Suche nach `macos` und `runner`
  über das Register fand nur Treffer zu den Ubuntu-Runnern der CI und zur Plattform macOS in anderem Zusammenhang); §5 sieht sie vor.

**Modus-Begründungsblock:** entfällt — alle berührten Sub-Areas GF bzw. keine Sub-Area (Aspirantin).
