# Review-Report: slice-v1-hilfetexte-ohne-kennungen — 2026-11-01

**Review-Art:** Code, Diff `1805e54a..HEAD` (Umsetzung `43544e0a`, davor zwei reine Lifecycle-Commits).
**Reviewer:** Rolle Reviewer in eigenem Kontext (`.harness/skills/reviewer.md`); der Reviewer durfte keine Datei schreiben, der Implementer legt den Report ab.
**Entscheidungen des Projektinhabers (2026-11-01):** Frage 1 Weg A (Warnzeile `Warning: <Text>`), Frage 3 Vorlagen im selben Slice, Frage 4 Backtick-Wörter mit; Frage 2 (Release) offen bis zur Closure.
**Sensoren:** `make gates` grün (Coverage 92,20 %); `make test-docker` lief nicht (die CI trägt den Nachweis; Wortlaut in `internal/e2e/devcontainer_shell_docker_test.go`).

## Findings und Ausgang

| Nr. | Stufe | Befund | Ausgang |
|---|---|---|---|
| H1 | HIGH | Kommentare in `remove.go` und `removeservice.go` beschrieben „v0.3.0 / deferred“, obwohl `--purge` das Volume entfernt ([`AGENTS.md`](../../AGENTS.md) §3.7, §3.9 (b)) | behoben: sieben Stellen auf den Ist-Zustand geschrieben |
| M1 | MEDIUM | Plan-Tabelle §3 nicht gegen `git diff --stat` abgeglichen; Zustandsmarker veraltet ([`AGENTS.md`](../../AGENTS.md) §3.10) | behoben: Tabelle, §6, §7 und Zeile in `open/README.md` nachgezogen |
| M2 | MEDIUM | Warnzeilen-Pins in `devcontainer_test.go`, `devcontainer_shell_test.go`, `updown_cluster_test.go` prüften nur noch `Warning:` | behoben: Pins tragen den Text der Warnung als Anker |
| L1 | LOW | `.devcontainer/Dockerfile` und `.devcontainer/docker-init.sh` des Repos trugen die alten Vorlagen-Kommentare | behoben: Kommentare wie die Vorlagen; Plan-Zeile ergänzt |
| L2 | LOW | Zeile in `doctor.go` nach dem Streichen unförmig | behoben: neu umbrochen |
| L3 | LOW | Namensregel in `config`-Hilfe und Doctor-Hinweis unvollständig (letztes Zeichen) | behoben: „ending with a letter or digit“ |
| L4 | LOW | Kommentar „have crept into“ in `help_identifiers_test.go` erzählt Herkunft | behoben: Zusage statt Herkunft |
| I1 | INFO | Grenzen der Wächter | siehe unten; geht in den Lerneintrag der Closure-Notiz |
| I2 | INFO | `--purge`-Hilfe stimmt mit `applyVolumePurge` überein; „the named volume“ ist für keycloak und otel (ohne Volume) unscharf, nicht falsch | keine Änderung |
| I3 | INFO | Release-Zuordnung und Wirkung der Warnzeile auf Skripte | beim Projektinhaber (Frage 2) |

## Grenzen der Wächter (I1)

- **Wächter A (Hilfebaum):** sammelt `Short`, `Long`, `Example` und lokale Optionen; `Use`, `Aliases` und die Cobra-Texte von `help` und `completion` fehlen. Die Ankündigungsliste ist eng („currently“, „for now“, „future“ sind nicht erfasst).
- **Wächter B (Meldungs-Literale):** erfasst Literale in `errors.New`, `fmt.*`, `warn`/`relax` (zweites Argument) und den Feldern `Message`, `Hint`, `MissingHint`, `Detail` sowie in Namen wie `hint`, `message`, `msg`, `prose`, `reason`. Nicht erfasst: `:=`-Zuweisungen, über `%s` oder Hilfsfunktionen eingesetzte Kennungen, Pakete außerhalb der Wurzeln (`cmd/uboot`), `.tmpl`-Dateien.
- **Vorlagen:** Die Kommentare in den Vorlagen sind durch keinen Wächter gedeckt; eine neue Kennung dort bleibt grün.

## Negativbefund

Kein `//nolint`, keine Änderung an `.golangci.yml`, `Makefile` oder `repo.mk`, keine Senkung der Schwelle, kein Accepted-ADR berührt, keine Abwärts-Verweise aus der Spezifikation. Neue Texte sind Englisch. Aussagen im neuen Text sind gegen den Code gelesen: 63 Zeichen (`projectNameMaxLen`), Katalog (`servicename.go`), `--purge` (`removeservice.go`, `remove.go`), `init --template` (`composite.go`), „at least three“ (`softExistingThreshold`).

## Nicht geprüft

`make test-docker`; die Golden-Wortlaut-Pins der Vorlagen einzeln (gedeckt durch grüne `make gates`).
