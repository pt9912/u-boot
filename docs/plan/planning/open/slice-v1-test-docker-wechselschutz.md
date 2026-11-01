# Slice V1: Wechselschutz für `make test-docker` — der Start stellt fest, ob der Stand seit dem letzten grünen Lauf unverändert ist, und lehnt den Wiederholungslauf ab; CI läuft unbehindert

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle, Release-Zuordnung offen (Werkzeug-Änderung, kein Produktvertrag).

**Anlass (Projektinhaber, 2026-11-01, im Verifikationslauf zu [`slice-v1-devcontainer-host-git-uebernehmen`](../in-progress/slice-v1-devcontainer-host-git-uebernehmen.md)):** Der Verifier fuhr den `test-docker`-Gesamtlauf (~40 min) erneut, obwohl der Stand seit dem letzten grünen Lauf nur von der Fix-Runde berührt war und der DoD-Nachweis ohnehin der CI-Lauf ist. Eine Regel („Gesamt-Suiten nur in CI") reicht nicht: **Der Start selbst muss den Nicht-Bedarf feststellen.** Die Mechanik existiert als Muster: `record-gates` stempelt HEAD (`tools/harness/record-gates.sh`, Kopplung mit dem Stop-Hook) und den Inhalts-Hash (`tools/harness/working-tree-hash.sh`); der Stop-Hook verweigert das Turn-Ende bei nicht gedecktem Inhalt.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Implementer, Entwurf auf Wunsch des Projektinhabers). **Datum:** 2026-11-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** `make test-docker` prüft **vor** dem Start, ob der Arbeitsbaum-Stand
(HEAD plus Inhalts-Hash, gleiche Kopplung wie bei `record-gates`) seit dem
letzten grünen Lauf unverändert ist. Ist er es — und der Lauf steht nicht in
CI und ist nicht erzwungen — endet das Target rot mit einer Meldung, die
Nebenbedeutung trägt: der Stand ist gedeckt, der Gesamt-Nachweis trägt der
CI-Lauf nach dem Push; ein erzwungener Lauf ist `FORCE=1 make test-docker`.
Nach einem grünen Lauf schreibt das Target den Stempel
(`.harness/state/test-docker-passed.head`, `….diffsha`). In CI
(`GITHUB_ACTIONS` gesetzt) läuft das Target unbehindert und stempelt nichts.

### Gefundener Stand (gelesen am 2026-11-01)

- `tools/harness/record-gates.sh` (gelesen): schreibt `.harness/state/gates-passed.head` (aufgelöste SHA via `head_wert`) und `….diffsha` (via `tools/harness/working-tree-hash.sh`); der Stop-Hook (`.claude/hooks/stop-require-gates.sh`, nicht gelesen) vergleicht denselben Hash.
- `repo.mk` `test-docker:` (gelesen, Zeilen 186 bis 204): `$(DOCKER_BUILD) --target test-docker-tools -t $(IMAGE):test-docker-tools .` + `docker run --rm --network=host -v … $(IMAGE):test-docker-tools go test -v -tags docker -timeout 30m ./...`; im lokalen Modus die verschachtelte Engine mit derselben Rezeptur über `scripts/bau.sh`? (nicht gelesen — prüfen beim Bau). CI-Aufruf: `.github/workflows/integration.yml` Zeile 47 `make test-docker` (gelesen).
- Arbeitsverzeichnis-Hilfsmittel: `tools/harness/working-tree-hash.sh` (nicht gelesen — beim Bau lesen).

**Ausdrücklich NICHT in diesem Slice:**

- Keine Regelung für `make gates` (sein Stop-Hook-Stempel regelt die Deckung bereits), `make test` oder die Coverage-Targets; nur `test-docker`.
- Keine Aufteilung der Suite (nights/Job-Verzweigung: ein eigener Slice, wenn der Projektinhaber Option B der Kosten-Abwägung wählt) und keine Beschleunigung der Tests.
- Keine Änderung der CI-Evidenz-Klasse: der DoD-Nachweis bleibt der CI-Lauf; der Stempel ist kein Ersatzbeleg für Dokumente, die den CI-Lauf verlangen.

## 2. Definition of Done

- [ ] Das Target verweigert den Wiederholungslauf bei unverändertem Stand (Meldung mit dem Stand-Bezug und dem `FORCE=1`-Weg), läuft bei geändertem Stand, in CI und bei `FORCE=1`; nach jedem grünen Lauf trägt es den Stempel; ein abgebrochener Lauf schreibt keinen Stempel, und ein abgebrochener Lauf nach einem grünen lässt den alten Stépmer stehen (der geprüfte Zustand galt beim letzten grünen Lauf — die Prüfsumme des Inhalts steuert, nicht die Zeit).
- [ ] Skript-Tests je Fall (Muster `ablage-pruefung-test.sh`): unverändert+Stempel → Verweigerung; geänderter Inhalt (Mutation an einer Testdatei) → Lauf; `FORCE=1` → Lauf; CI-Umgebung → Lauf ohne Stempel; Stempel-Format wie bei `gates-passed`.
- [ ] `harness/README.md` Sensors-Zeile von `test-docker` trägt die Wächter-Form (Bindung der Spalte ergänzt); `harness/mk`-Fragmente (MR-012) unberührt; `make gates` bleibt grün.
- [ ] Review, Verifikation im frischen Kontext, Closure-Notiz, `make gates` grün auf dem Abschluss-Commit.

**Wächter und rot färbende Mutationen:**

| Wächter | Mutation |
|---|---|
| Skript-Test „unverändert+Stempel → Verweigerung" | die Vergleichslogge auf HEAD allein beschränken (Inhalts-Hash ignoriert) |
| „geänderter Inhalt → Lauf" | vor der Prüfung kein Inhalts-Hash, nur Zeitstempel (die Lücke: geänderter Inhalt, alter Zeitstempel) |

Ein Mutationstest je Wächter bei der Closure geführt (Projektinhaber-Praxis).

## 3. Plan (vor Code)

| Datei / Komponente | Art | Begründung / Fundstelle |
|---|---|---|
| `tools/harness/test-docker-guard.sh` | neu | Modus `check`/`record`; Wiederverwendung von `head_wert` (Kopie aus `record-gates.sh`, dessen Form beibehalten) und `working-tree-hash.sh` |
| `repo.mk` (`test-docker-Recipeur`, beide Mod-Zweige) | update | `check` vor dem Laufruf, `record` nach dem grünen Lauf; `FORCE` als Variable; CI-Bypass via `GITHUB_ACTIONS` |
| `tools/harness/test-docker-guard-test.sh` | neu | die drei Läufe des DoD |
| `.harness/state/` (neue Stempel) | neu | `.gitignore` prüfen: `gates-passed` ist verworfen? (prüfen) |
| `harness/README.md` (Sensors-Tabelle Zeile `test-docker`) | update | die Wächter-Zeile; Form wie die Zeilen der Werkzeuge |
| `docs/plan/planning/open/README.md`, diese Datei | update | Zeile im Bestand |

## 4. Trigger

**Start:** Projektinhaber gibt den Auftrag nach der Closure von [`slice-v1-devcontainer-host-git-uebernehmen`](../in-progress/slice-v1-devcontainer-host-git-uebernehmen.md); kein anderer Slice in `in-progress/`.

**Rückführung:** —

## 5. Closure-Trigger

DoD vollständig; `make gates` grün; Review und Verifikation; Closure-Notiz; `git mv` nach `done/`.

## 6. Risiken

1. **Der Stempel verleitet, CI zu übergehen:** Der Stempel bedeutet nur, dass der Stand gedeckt wurde — die Zusage-Klasse der DoDs bleibt CI. — Ausgang: bei Closure; die Meldung nennt es.
2. **FORCE-Bedingungen:** erzwungene Läufe ohne Bedarf sind genau das, was der Wächter verhindern will; das Feld ist bewusst unbequem (Variable statt Flag). — Ausgang: bei Closure.
3. **CI-Bypass verliert den Stempel in CI:** dort stempelt nichts; das Muster `gates-passed` wird nicht gebrochen. — Ausgang: bei Closure.
4. **Abbruch nach grünem Lauf:** der alte Stépmer bleibt; der nächste Lauf prüft Inhalts-Gleichheit, nicht Zeit. — Ausgang: bei Closure (DoD-Punkt 1).

## 7. Closure-Notiz

Noch nicht geschlossen; der Slice liegt in `open/` (Entwurf).

## 8. Sub-Area-Prüfungen

Berührt werden `tools/` und `scripts/` (WZG, Greenfield laut
[`harness/conventions.md`](../../harness/conventions.md#modus-deklaration-pro-sub-area)),
`harness/` (DOK, Greenfield) und `.github/workflows/` (lesend). Evidenz-/Diskrepanz-Risiko: niedrig (Skript-Läufe gegen den eigenen Stempel).